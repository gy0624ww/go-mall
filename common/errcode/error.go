package errcode

import (
	"encoding/json"
	"fmt"
	"net/http"
	"path"
	"runtime"
)

type AppError struct {
	code     int    `json:"code"`
	msg      string `json:"msg"`
	cause    error  `json:"cause"`
	occurred string `json:"occurred"` // 保存由底层错误导致AppErr发生时的位置
}

func (e *AppError) Error() string {
	if e == nil {
		return ""
	}
	errBytes, err := json.Marshal(e.toStructuredError())
	if err != nil {
		return fmt.Sprintf("Error() is error: json marshal error: %v", err)
	}
	return string(errBytes)
}

func (e *AppError) String() string {
	return e.Error()
}

func (e *AppError) Code() int {
	return e.code
}

func (e *AppError) Msg() string {
	return e.msg
}

func (e *AppError) HttpStatusCode() int {
	switch e.Code() {
	case Success.Code():
		return http.StatusOK
	case ErrServer.Code(), ErrPanic.Code():
		return http.StatusInternalServerError
	case ErrParams.Code():
		return http.StatusBadRequest
	case ErrNotFound.Code():
		return http.StatusNotFound
	case ErrTooManyRequests.Code():
		return http.StatusTooManyRequests
	case ErrToken.Code():
		return http.StatusUnauthorized
	case ErrForbidden.Code():
		return http.StatusForbidden
	default:
		return http.StatusInternalServerError
	}
}

// WithCause 在逻辑执行中出现错误, 比如dao层返回的数据库查询错误
// 可以在领域层返回预定义的错误前附加上导致错误的基础错误。
// 如果业务模块预定义的错误码比较详细, 可以使用这个方法, 反之错误码定义的比较笼统建议使用Wrap方法包装底层错误生成项目自定义Error
// 并将其记录到日志后再使用预定义错误码返回接口响应
func (e *AppError) WithCause(err error) *AppError {
	newErr := e.Clone()
	newErr.cause = err
	newErr.occurred = getAppErrOccurredInfo()
	return newErr
}

func (e *AppError) Clone() *AppError {
	return &AppError{
		code:     e.code,
		msg:      e.msg,
		cause:    e.cause,
		occurred: e.occurred,
	}
}

func newError(code int, msg string) *AppError {
	if code > -1 {
		if _, duplicated := codes[code]; duplicated {
			panic(fmt.Sprintf("预定义错误码 %d 不能重复，请检查后更换", code))
		}
	}
	return &AppError{code: code, msg: msg}
}

// getAppErrOccurredInfo 获取项目中调用Wrap或者WithCause方法时的程序位置, 方便排查问题
func getAppErrOccurredInfo() string {
	pc, file, line, ok := runtime.Caller(2)
	if !ok {
		return ""
	}
	file = path.Base(file)
	funcName := runtime.FuncForPC(pc).Name()
	triggerInfo := fmt.Sprintf("func: %s, file: %s, line: %d", funcName, file, line)
	return triggerInfo
}

// AppendMsg 在Code不变的情况下, 在预定义Msg的基础上追加错误信息
func (e *AppError) AppendMsg(msg string) *AppError {
	n := e.Clone()
	n.msg = fmt.Sprintf("%s, %s", e.msg, msg)
	return n
}

// SetMsg 在Code不变的情况下, 重新设置错误信息, 覆盖预定义的Msg
func (e *AppError) SetMsg(msg string) *AppError {
	n := e.Clone()
	n.msg = msg
	return n
}

type formattedErr struct {
	Code     int         `json:"code"`
	Msg      string      `json:"msg"`
	Cause    interface{} `json:"cause"`
	Occurred string      `json:"occurred"`
}

// toStructuredError 在JSON Encode 前把Error进行格式化
func (e *AppError) toStructuredError() *formattedErr {
	fe := new(formattedErr)
	fe.Code = e.Code()
	fe.Msg = e.Msg()
	fe.Occurred = e.occurred
	if e.cause != nil {
		if appErr, ok := e.cause.(*AppError); ok {
			fe.Cause = appErr.toStructuredError()
		} else {
			fe.Cause = e.cause.Error()
		}
	}
	return fe
}

// Wrap 用于逻辑中包装底层函数返回的error 和 WithCause 一样都是为了记录错误链条
func Wrap(msg string, err error) *AppError {
	if err == nil {
		return nil
	}
	appErr := &AppError{code: -1, msg: msg, cause: err}
	appErr.occurred = getAppErrOccurredInfo()
	return appErr
}

func (e *AppError) UnWrap() error {
	return e.cause
}

// Is 与上面的UnWrap一起让 *AppError 支持 errors.Is(err, target)
func (e *AppError) Is(target error) bool {
	targetErr, ok := target.(*AppError)
	if !ok {
		return false
	}
	return targetErr.Code() == e.Code()
}
