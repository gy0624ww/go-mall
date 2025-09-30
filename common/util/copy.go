package util

import (
	"errors"
	"regexp"
	"time"

	"github.com/go-study-lab/go-mall/common/enum"
	"github.com/jinzhu/copier"
)

// CopyPropetrties 把属性从src 复制到dst
// 参数请传 pointer 类型
func CopyPropetrties(dst, src interface{}) error {
	err := copier.CopyWithOption(dst, src, copier.Option{
		IgnoreEmpty: true,
		DeepCopy:    true,
		Converters: []copier.TypeConverter{
			{ // time.Time 转换成字符串
				SrcType: time.Time{},
				DstType: copier.String,
				Fn: func(src interface{}) (dst interface{}, err error) {
					s, ok := src.(time.Time)
					if !ok {
						return nil, errors.New("src type is not time.Time")
					}
					return s.Format(enum.TimeFormatHyphenedYMDHIS), nil
				},
			},
			{ // 字符串转成time.Time
				SrcType: copier.String,
				DstType: time.Time{},
				Fn: func(src interface{}) (dst interface{}, err error) {
					s, ok := src.(string)
					if !ok {
						return nil, errors.New("src type is not time format string")
					}
					pattern := `^\d{4}-\d{2}-\d{2} \d{2}:\d{2}:\d{2}$` // YYYY-MM-DD HH:MM:SS
					matched, _ := regexp.MatchString(pattern, s)
					if matched {
						return time.Parse(enum.TimeFormatHyphenedYMDHIS, s)
					}
					return nil, errors.New("src type is not time format string")
				},
			},
		},
	})
	return err
}
