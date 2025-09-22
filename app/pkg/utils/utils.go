package utils

import (
	"fmt"
	"math/rand"
	"net/url"
	"strings"
	"time"

	"github.com/go-kratos/kratos/v2/log"
	"github.com/go-kratos/kratos/v2/transport"
	"github.com/gogf/gf/v2/errors/gerror"
	"github.com/jinzhu/copier"
	"github.com/sirupsen/logrus"
	"google.golang.org/protobuf/types/known/timestamppb"
)

func GetShortName(email string) string {
	return strings.Split(email, "@")[0]
}

// 生成32位随机密钥（包含大小写字母和数字）
func GetRandomKey32() string {
	const charset = "0123456789abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ"
	bytes := []byte(charset)
	result := make([]byte, 32) // 预分配32字节空间

	r := rand.New(rand.NewSource(time.Now().UnixNano()))
	for i := range result {
		result[i] = bytes[r.Intn(len(bytes))]
	}
	return string(result)
}

// CalculateMaxMinAvg 计算最大值、最小值、平均值
func CalculateMaxMinAvg[T int](values []T) (max, min T, avg float64) {
	if len(values) <= 0 {
		return
	}

	max = -1
	min = -1
	var sum T

	for _, value := range values {
		value_ := value
		//if value_ == -316 {
		//	value_ = 40
		//}
		//if value_ == -301 {
		//	value_ = 300
		//}
		if value_ < 0 || value_ == 101 || value_ == -201 || value_ == -401 {
			continue
		}
		if max == -1 {
			max = value_
		}
		if min == -1 {
			min = value_
		}
		if value_ > max {
			max = value_
		}
		if value_ < min {
			min = value_
		}
		sum += value_
	}

	avg = float64(sum) / float64(len(values))
	return
}

func Int64ArrayToIntArray(int64Array []int64) []int {
	intArray := []int{}
	for _, v := range int64Array {
		intArray = append(intArray, int(v))
	}
	return intArray
}

func GetNumFromMap(key string, data map[string]interface{}) (float64, bool) {
	if val, ok := data[key]; ok {
		if intVal, ok := val.(float64); ok {
			return intVal, true
		}
	}
	return 0, false
}

// GetClientIP 从 HTTP 请求中获取客户端 IP
func GetClientIP(header transport.Header) string {
	if ip := header.Get("X-Forwarded-For"); ip != "" {
		return strings.Split(ip, ",")[0]
	}
	if ip := header.Get("X-Real-IP"); ip != "" {
		return ip
	}
	return ""
}

// FilePathToUrl 文件路径转换成url
func FilePathToUrl(s3Url, token, webUrl string) (url2 string) {
	uri, err := url.Parse(s3Url)
	if err != nil {
		logrus.Error(err)
		return
	}
	filePath := uri.Host + uri.Path
	var params []string
	params = append(params, "filename="+filePath)
	if token != "" {
		params = append(params, "token="+token)
	}
	queryStr := ""
	if len(params) > 0 {
		queryStr = "?" + strings.Join(params, "&")
	}

	url2 = fmt.Sprintf("%s/%s%s", webUrl, "file", queryStr)
	return
}

// UrlToFilePath url转成文件路径
func UrlToFilePath(url2 string) (filePath string) {
	if url2 == "" {
		return
	}
	uri, err := url.Parse(url2)
	if err != nil {
		logrus.Error(err)
		return
	}
	query := uri.Query()
	fileName := query.Get("filename")
	filePath = fileName
	return
}

func IntArrayToInt64Array(intArray []int) []int64 {
	int64Array := []int64{}
	for _, v := range intArray {
		int64Array = append(int64Array, int64(v))
	}
	return int64Array
}

// IsTimeInRange 判断时间是否在指定时间区间内
func IsTimeInRange(timestamp time.Time, timeStartStr, timeEndStr string) bool {
	// 将开始和结束时间字符串解析为时间
	startTime, err := time.Parse("15:04", timeStartStr)
	if err != nil {
		logrus.Errorln(err)
		return false
	}

	endTime, err := time.Parse("15:04", timeEndStr)
	if err != nil {
		logrus.Errorln(err)
		return false
	}
	sHour := startTime.Hour()
	sMinute := startTime.Minute()
	eHour := endTime.Hour()
	eMinute := endTime.Minute()
	mHour := timestamp.Hour()
	mMinute := timestamp.Minute()

	flag1 := sHour <= mHour && sMinute <= mMinute
	flag2 := eHour > mHour || (eHour == mHour && eMinute >= mMinute)
	return flag1 && flag2

	//// 将 timestamp 调整到同一天，以便与开始和结束时间进行比较
	//timestamp = time.Date(time.Now().Year(), time.Now().Month(), time.Now().Day(), timestamp.Hour(), timestamp.Minute(), 0, 0, timestamp.Location())
	//
	//return timestamp.After(startTime) && timestamp.Before(endTime)
}

func GetTimeZoneOffsetByTimeZoneName(timeZoneName string) int {
	var timeZoneOffset int
	timeZone := strings.TrimSpace(timeZoneName)
	if timeZone != "" {
		loc, err := time.LoadLocation(timeZone)
		if err == nil && loc != nil {
			_, timeZoneOffset = time.Now().In(loc).Zone()
		}
	}
	return timeZoneOffset
}

// 获取时区
func GetTimeZone(timeZoneName string) *time.Location {
	timeZone := strings.TrimSpace(timeZoneName)
	loc, err := time.LoadLocation(timeZone)
	if err != nil {
		log.Error("load time zone failed, err: %+v", gerror.Wrap(err, ""))
		return time.UTC
	}
	return loc
}

// 获取时区和偏移量字符串,偏移量格式为"UTC+08:00"这样的字符串
func GetTimeZoneAndOffsetStr(timeZoneName string) (*time.Location, string) {
	timeZone := strings.TrimSpace(timeZoneName)
	loc, err := time.LoadLocation(timeZone)
	if err != nil {
		log.Error("load time zone failed, err: %+v", gerror.Wrap(err, ""))
		return time.UTC, timeZoneName
	}
	// 获取当前时间并将其转换为目标时区
	currentTime := time.Now().In(loc)

	// 获取时区偏移量（秒）
	_, offset := currentTime.Zone()
	hours := offset / 3600
	minutes := (offset % 3600) / 60

	// 处理负数分钟的问题，确保分钟数始终为正数
	if minutes < 0 {
		minutes = -minutes
	}

	// 手动添加符号和前导零
	sign := "+"
	if hours < 0 {
		sign = "-"
		hours = -hours
	}

	// 格式化为 UTC+08:00(Asia/Shanghai) 或 UTC-03:30(America/St_Johns) 或 UTC+05:30(Asia/Kolkata) 的形式
	formattedOffset := fmt.Sprintf("UTC%s%02d:%02d", sign, hours, minutes)

	return loc, formattedOffset
}

// time.Time类型与*timestamppb.Timestamp类型互相转换器
var CopierConverter = copier.Option{
	IgnoreEmpty: true,
	DeepCopy:    true,
	Converters: []copier.TypeConverter{
		{
			SrcType: time.Time{},
			DstType: &timestamppb.Timestamp{},
			Fn: func(src interface{}) (interface{}, error) {
				s, ok := src.(time.Time)
				if !ok {
					return nil, nil
				}
				return timestamppb.New(s), nil
			},
		},
		{
			SrcType: &time.Time{},
			DstType: &timestamppb.Timestamp{},
			Fn: func(src interface{}) (interface{}, error) {
				s, ok := src.(*time.Time)
				if !ok || s == nil { // 如果IgnoreEmpty设置为了false，s可能为nil
					return nil, nil
				}
				return timestamppb.New(*s), nil
			},
		},
		// 反向转换器
		{
			SrcType: &timestamppb.Timestamp{},
			DstType: time.Time{},
			Fn: func(src interface{}) (interface{}, error) {
				s, ok := src.(*timestamppb.Timestamp)
				if !ok || s == nil {
					return time.Time{}, nil
				}
				return s.AsTime(), nil
			},
		},
		{
			SrcType: &timestamppb.Timestamp{},
			DstType: &time.Time{},
			Fn: func(src interface{}) (interface{}, error) {
				s, ok := src.(*timestamppb.Timestamp)
				if !ok || s == nil {
					return nil, nil
				}
				return Ptr(s.AsTime()), nil
			},
		},
		// 时间转毫秒时间戳
		{
			SrcType: time.Time{},
			DstType: int64(0),
			Fn: func(src interface{}) (interface{}, error) {
				s, ok := src.(time.Time)
				if !ok {
					return nil, nil
				}
				return s.UnixMilli(), nil // 直接返回Unix时间戳
			},
		},
		{
			SrcType: &time.Time{},
			DstType: int64(0),
			Fn: func(src interface{}) (interface{}, error) {
				s, ok := src.(*time.Time)
				if !ok || s == nil {
					return nil, nil
				}
				return s.UnixMilli(), nil // 直接返回Unix时间戳
			},
		},
	},
}

// 可以处理time.Time类型与*timestamppb.Timestamp类型互相转换的copy
func Copy(toValue interface{}, fromValue interface{}) {
	copier.CopyWithOption(toValue, fromValue, CopierConverter)
}

// 判断字符串是否为有效的URL
// str: 待检查的字符串
// 返回值: 如果是有效URL返回true，否则返回false
func IsURL(str string) bool {
	// 解析URL
	u, err := url.Parse(str)
	if err != nil {
		// 解析失败，不是有效URL
		return false
	}
	// 有效的URL需要包含协议和主机名
	return u.Scheme != "" && u.Host != ""
}
