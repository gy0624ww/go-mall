package httptool

import (
	"bytes"
	"context"
	"errors"
	"fmt"
	"io/ioutil"
	"net"
	"net/http"
	"strings"
	"sync"
	"time"

	"github.com/go-study-lab/go-mall/common/errcode"
	"github.com/go-study-lab/go-mall/common/logger"
	"github.com/go-study-lab/go-mall/common/util"
)

var (
	_Client *http.Client
	once    sync.Once
)

const (
	maxLogContentSize = 1024 * 1024 // 1MB
)

// formatLogContent 格式化日志内容，超过大小限制或文件上传时返回提示信息
func formatLogContent(content []byte, headers map[string]string) string {
	// 检查是否是文件上传
	if contentType, exists := headers["Content-Type"]; exists {
		if strings.Contains(contentType, "multipart/form-data") {
			return "File upload data, skip logging"
		}
	}
	// 检查内容大小
	if len(content) > maxLogContentSize {
		return "Data too long, skip logging"
	}
	return string(content)
}

func getHttpClient() *http.Client {
	if _Client != nil {
		// 因为Unit test里要把Client换掉, 所以虽然用了once.Do 但是还是先判断一下_Client有没有实例化
		// 不然在单测里, Mock API时, gock没办法拦截对外部API的请求
		return _Client
	}
	once.Do(func() {
		// MaxIdleConnsPerHost：决定了对于单个Host需要维持的连接池大小。该值应该根据性能测试的结果调整。
		// MaxIdleConns：全局的最大空闲连接，不要比MaxIdleConnsPerHost小，嫌麻烦的话建议不设置或者设置为0 --- 即不限制。
		// MaxConnsPerHost：对于单个Host允许的最大连接数，包含IdleConns，所以一般大于等于MaxIdleConnsPerHost。设置为等于MaxIdleConnsPerHost，也就是尽可能复用连接池中的连接。另外设置过小，可能会导致并发下降，它的默认值是不做限制。
		tr := &http.Transport{
			//Proxy: http.ProxyFromEnvironment,
			DialContext: (&net.Dialer{
				Timeout:   30 * time.Second,
				KeepAlive: 30 * time.Second,
			}).DialContext,
			MaxIdleConns:          100,
			IdleConnTimeout:       90 * time.Second,
			MaxIdleConnsPerHost:   50,
			MaxConnsPerHost:       50,
			ForceAttemptHTTP2:     true,
			TLSHandshakeTimeout:   10 * time.Second,
			ExpectContinueTimeout: 1 * time.Second,
		}
		_Client = &http.Client{Transport: tr}
	})
	return _Client
}

// SetUTHttpClient 让单元测试能把httpClient覆盖成具有Mock拦截设置的HttpClient
func SetUTHttpClient(client *http.Client) {
	_Client = client
}

func Request(method string, url string, options ...Option) (httpStatusCode int, respBody []byte, err error) {
	start := time.Now()
	reqOpts := defaultRequestOptions() // 默认的请求选项
	for _, opt := range options {      // 在reqOpts上应用通过options设置的选项
		err = opt.apply(reqOpts)
		if err != nil {
			return
		}
	}
	defer func() {
		if err != nil {
			logger.Error(reqOpts.ctx, "HTTP_REQUEST_ERROR_LOG",
				"method", method,
				"url", url,
				"body", formatLogContent(reqOpts.data, reqOpts.headers),
				"reply", formatLogContent(respBody, reqOpts.headers),
				"status", httpStatusCode,
				"err", err)
		}
	}()
	// 创建请求对象
	req, err := http.NewRequest(method, url, bytes.NewReader(reqOpts.data))
	if err != nil {
		return
	}
	reqOpts.ctx, _ = context.WithTimeout(reqOpts.ctx, reqOpts.timeout) // 给 Request 设置Timeout
	req = req.WithContext(reqOpts.ctx)
	defer req.Body.Close()

	// 在Header中添加追踪信息 把内部服务串起来
	traceId, spanId, _ := util.GetTraceInfoFromCtx(reqOpts.ctx)
	reqOpts.headers["traceid"] = traceId
	reqOpts.headers["spanid"] = spanId
	if len(reqOpts.headers) != 0 { // 设置请求头
		for key, value := range reqOpts.headers {
			req.Header.Add(key, value)
		}
	}
	// 发起请求
	client := getHttpClient()
	resp, err := client.Do(req)
	if err != nil {
		return
	}
	defer resp.Body.Close()
	// 记录请求日志
	dur := time.Since(start).Milliseconds()
	defer func() {
		if dur >= 3000 { // 超过 3s 返回, 记一条 Warn 日志
			logger.Warn(reqOpts.ctx, "HTTP_REQUEST_SLOW_LOG", "method", method, "url", url, "body", formatLogContent(reqOpts.data, reqOpts.headers), "reply", formatLogContent(respBody, reqOpts.headers), "err", err, "dur/ms", dur)
		} else {
			logger.Debug(reqOpts.ctx, "HTTP_REQUEST_DEBUG_LOG", "method", method, "url", url, "body", formatLogContent(reqOpts.data, reqOpts.headers), "reply", formatLogContent(respBody, reqOpts.headers), "err", err, "dur/ms", dur)
		}
	}()

	// 先读取响应体, 再进行状态码检查, 避免非HTTP 200 状态码时, 日志中响应体记录不到
	respBody, _ = ioutil.ReadAll(resp.Body)

	httpStatusCode = resp.StatusCode
	if httpStatusCode != http.StatusOK {
		// 返回非 200 时Go的 http 库不回返回error, 这里处理成error 调用方好判断
		err = errcode.Wrap("request api error", errors.New(fmt.Sprintf("non 200 response, response code: %d", httpStatusCode)))
		return
	}

	return
}

func Get(ctx context.Context, url string, options ...Option) (httpStatusCode int, respBody []byte, err error) {
	options = append(options, WithContext(ctx))
	return Request("GET", url, options...)
}

// Post 发起POST请求
func Post(ctx context.Context, url string, data []byte, options ...Option) (httpStatusCode int, respBody []byte, err error) {
	// 默认自带Header Content-Type: application/json 可通过 传递 WithHeaders 增加或者覆盖Header信息
	defaultHeader := map[string]string{"Content-Type": "application/json"}
	var newOptions []Option
	newOptions = append(newOptions, WithHeaders(defaultHeader), WithData(data), WithContext(ctx))
	newOptions = append(newOptions, options...)

	httpStatusCode, respBody, err = Request("POST", url, newOptions...)
	return
}

// 针对可选的HTTP请求配置项，模仿gRPC使用的Options设计模式实现
type requestOption struct {
	ctx     context.Context
	timeout time.Duration
	data    []byte
	headers map[string]string
}

func defaultRequestOptions() *requestOption {
	return &requestOption{
		ctx:     context.Background(),
		timeout: 5 * time.Second,
		data:    nil,
		headers: map[string]string{},
	}
}

type Option interface {
	apply(option *requestOption) error
}

type optionFunc func(option *requestOption) error

func (f optionFunc) apply(opts *requestOption) error {
	return f(opts)
}

func WithContext(ctx context.Context) Option {
	return optionFunc(func(opts *requestOption) (err error) {
		opts.ctx = ctx
		return
	})
}

func WithTimeout(timeout time.Duration) Option {
	return optionFunc(func(opts *requestOption) (err error) {
		opts.timeout, err = timeout, nil
		return
	})
}

func WithHeaders(headers map[string]string) Option {
	return optionFunc(func(opts *requestOption) (err error) {
		for k, v := range headers {
			opts.headers[k] = v
		}
		return
	})
}
func WithData(data []byte) Option {
	return optionFunc(func(opts *requestOption) (err error) {
		opts.data, err = data, nil
		return
	})
}
