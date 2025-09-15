package encoder

import (
	"encoding/json"
	"html"
	"regexp"

	log "github.com/sirupsen/logrus"

	"reflect"
	"time"
	"unsafe"

	jsoniter "github.com/json-iterator/go"
	"github.com/modern-go/reflect2"
	"google.golang.org/protobuf/encoding/protojson"
	"google.golang.org/protobuf/proto"
	"google.golang.org/protobuf/types/known/timestamppb"
)

var (
	// UnmarshalOptions is a configurable JSON format parser.
	UnmarshalOptions = protojson.UnmarshalOptions{
		DiscardUnknown: true,
	}
	// htmlTagRegex 用于匹配 HTML 标签
	htmlTagRegex = regexp.MustCompile(`<[^>]*>`)
)

type JsonCodec struct{}

func (JsonCodec) Marshal(v any) ([]byte, error) {
	return jsoniter.Marshal(v)
}

func (JsonCodec) Unmarshal(data []byte, v any) error {
	// 记录原始数据长度
	originalLength := len(data)
	// 移除 HTML 标签
	cleanedData := htmlTagRegex.ReplaceAll(data, []byte(""))
	// 反转义 HTML 实体
	cleanedStr := html.UnescapeString(string(cleanedData))
	cleanedData = []byte(cleanedStr)
	// 记录清洗后的数据长度
	cleanedLength := len(cleanedData)

	// 记录日志
	if originalLength != cleanedLength {
		log.Errorf("XSS tags removed during Unmarshal. Original data length: %d, Cleaned data length: %d", originalLength, cleanedLength)
	}

	switch m := v.(type) {
	case json.Unmarshaler:
		return m.UnmarshalJSON(cleanedData)
	case proto.Message:
		return UnmarshalOptions.Unmarshal(cleanedData, m)
	default:
		rv := reflect.ValueOf(v)
		for rv := rv; rv.Kind() == reflect.Ptr; {
			if rv.IsNil() {
				rv.Set(reflect.New(rv.Type().Elem()))
			}
			rv = rv.Elem()
		}
		if m, ok := reflect.Indirect(rv).Interface().(proto.Message); ok {
			return UnmarshalOptions.Unmarshal(cleanedData, m)
		}
		return jsoniter.Unmarshal(cleanedData, m)
	}
}

func (JsonCodec) Name() string {
	return "json"
}

func init() {
	// json-iterator忽略struct“omitempty”标签
	jsoniter.RegisterExtension(&ignoreOmitEmptyTagExtension{})
	registerTimestampSerializer()
}

// 注册 Timestamp 的自定义序列化器
func registerTimestampSerializer() {
	jsoniter.RegisterTypeEncoderFunc(
		reflect.TypeOf((*timestamppb.Timestamp)(nil)).Elem().String(), // 确保类型名正确
		func(ptr unsafe.Pointer, stream *jsoniter.Stream) {
			// 先将 unsafe.Pointer 转换为 *timestamppb.Timestamp
			ts := (*timestamppb.Timestamp)(ptr)
			if ts == nil {
				stream.WriteNil()
				return
			}
			// 使用 RFC3339Nano 格式序列化
			stream.WriteString(ts.AsTime().Format(time.RFC3339Nano))
		},
		func(ptr unsafe.Pointer) bool {
			ts := (*timestamppb.Timestamp)(ptr)
			return ts == nil || (ts.Seconds == 0 && ts.Nanos == 0)
		},
	)
}

// json-iterator忽略struct“omitempty”标签
type ignoreOmitEmptyTagExtension struct {
	jsoniter.DummyExtension
}

type ignoreOmitEmptyTagEncoder struct {
	originDecoder jsoniter.ValEncoder
}

func (p *ignoreOmitEmptyTagEncoder) IsEmpty(ptr unsafe.Pointer) bool { //关键逻辑
	return false
}

func (p *ignoreOmitEmptyTagEncoder) Encode(ptr unsafe.Pointer, stream *jsoniter.Stream) {
	p.originDecoder.Encode(ptr, stream)
}

func (e *ignoreOmitEmptyTagExtension) DecorateEncoder(typ reflect2.Type, encoder jsoniter.ValEncoder) jsoniter.ValEncoder {
	return &ignoreOmitEmptyTagEncoder{encoder}
}
