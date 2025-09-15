package cast

import (
	"encoding/json"
	"time"
)

// remove slice item
func RemoveSliceItem(list []string, item string) []string {
	if len(list) == 0 {
		return make([]string, 0)
	}
	j := 0
	for i := 0; i < len(list); i++ {
		if list[i] != item {
			list[j] = list[i]
			j++
		} else if i != len(list)-1 {
			list[j] = list[i+1]
		}
	}
	return list[:j]
}

// RemoveRepeat Remove Repeat
func RemoveRepeat(list []string) []string {
	return RemoveRepeatSlice(list)
}

// RemoveRepeatInt Remove Repeat
func RemoveRepeatInt(list []int) []int {
	return RemoveRepeatSlice(list)
}

// RemoveRepeatSlice 去重，泛型方法
func RemoveRepeatSlice[T comparable](list []T) []T {
	if len(list) == 0 {
		return nil
	}
	m := make(map[T]struct{}, len(list))
	result := make([]T, 0, len(list))
	for _, v := range list {
		// new element
		if _, ok := m[v]; !ok {
			m[v] = struct{}{}
			result = append(result, v)
		}
	}
	return result
}

// CompareSlice 比较两个切片的内容是否一致
func CompareSlice[T comparable](arr1, arr2 []T) bool {
	if len(arr1) != len(arr2) {
		return false
	}
	m := make(map[T]struct{}, len(arr1))
	for _, v := range arr1 {
		m[v] = struct{}{}
	}
	for _, v := range arr2 {
		if _, ok := m[v]; !ok {
			return false
		}
	}
	return true
}

// IsSliceLess 判断第一个切片是否在元素数量或元素内容上 "小于" 第二个切片。
// 若第一个切片长度大于第二个切片，返回 true；
// 若第一个切片存在第二个切片没有的元素，返回 true；
// 同时返回第一个切片中存在但第二个切片中没有的元素。
func IsSliceLess[T comparable](arr1, arr2 []T) (bool, []T) {
	// 若 arr1 长度大于 arr2，直接返回 true
	if len(arr1) > len(arr2) {
		return true, arr1
	}

	// 记录 arr1 中元素的存在情况
	elemSet := make(map[T]struct{}, len(arr1))
	for _, v := range arr1 {
		elemSet[v] = struct{}{}
	}

	// 移除 arr2 中存在于 arr1 的元素
	for _, v := range arr2 {
		delete(elemSet, v)
	}

	// 收集第一个切片中存在但第二个切片中没有的元素
	var removedElements []T
	for v := range elemSet {
		removedElements = append(removedElements, v)
	}

	// 若 elemSet 不为空，说明 arr1 存在 arr2 没有的元素
	return len(elemSet) > 0, removedElements
}

// SliceContain 切片中是否包含某元素
func SliceContain[T comparable](list []T, item T) bool {
	for _, v := range list {
		if v == item {
			return true
		}
	}
	return false
}

// MapToSlice map to slice
func MapToSlice(m map[string]struct{}) []string {
	l := len(m)
	if l == 0 {
		return nil
	}
	list := make([]string, 0, l)
	for k, _ := range m {
		list = append(list, k)
	}
	return list
}

// ToBool casts an interface to a bool type.
func ToBool(i interface{}) bool {
	v, _ := ToBoolE(i)
	return v
}

// ToTime casts an interface to a time.Time type.
func ToTime(i interface{}) time.Time {
	v, _ := ToTimeE(i)
	return v
}

// ToDuration casts an interface to a time.Duration type.
func ToDuration(i interface{}) time.Duration {
	v, _ := ToDurationE(i)
	return v
}

// ToFloat64 casts an interface to a float64 type.
func ToFloat64(i interface{}) float64 {
	v, _ := ToFloat64E(i)
	return v
}

// ToFloat32 casts an interface to a float32 type.
func ToFloat32(i interface{}) float32 {
	v, _ := ToFloat32E(i)
	return v
}

// ToInt64 casts an interface to an int64 type.
func ToInt64(i interface{}) int64 {
	v, _ := ToInt64E(i)
	return v
}

// ToInt32 casts an interface to an int32 type.
func ToInt32(i interface{}) int32 {
	v, _ := ToInt32E(i)
	return v
}

// ToInt16 casts an interface to an int16 type.
func ToInt16(i interface{}) int16 {
	v, _ := ToInt16E(i)
	return v
}

// ToInt8 casts an interface to an int8 type.
func ToInt8(i interface{}) int8 {
	v, _ := ToInt8E(i)
	return v
}

// ToInt casts an interface to an int type.
func ToInt(i interface{}) int {
	v, _ := ToIntE(i)
	return v
}

// ToUint casts an interface to a uint type.
func ToUint(i interface{}) uint {
	v, _ := ToUintE(i)
	return v
}

// ToUint64 casts an interface to a uint64 type.
func ToUint64(i interface{}) uint64 {
	v, _ := ToUint64E(i)
	return v
}

// ToUint32 casts an interface to a uint32 type.
func ToUint32(i interface{}) uint32 {
	v, _ := ToUint32E(i)
	return v
}

// ToUint16 casts an interface to a uint16 type.
func ToUint16(i interface{}) uint16 {
	v, _ := ToUint16E(i)
	return v
}

// ToUint8 casts an interface to a uint8 type.
func ToUint8(i interface{}) uint8 {
	v, _ := ToUint8E(i)
	return v
}

// ToString casts an interface to a string type.
func ToString(i interface{}) string {
	v, _ := ToStringE(i)
	return v
}

// ToStringMapString casts an interface to a map[string]string type.
func ToStringMapString(i interface{}) map[string]string {
	v, _ := ToStringMapStringE(i)
	return v
}

// ToStringMapStringSlice casts an interface to a map[string][]string type.
func ToStringMapStringSlice(i interface{}) map[string][]string {
	v, _ := ToStringMapStringSliceE(i)
	return v
}

// ToStringMapBool casts an interface to a map[string]bool type.
func ToStringMapBool(i interface{}) map[string]bool {
	v, _ := ToStringMapBoolE(i)
	return v
}

// ToStringMapInt casts an interface to a map[string]int type.
func ToStringMapInt(i interface{}) map[string]int {
	v, _ := ToStringMapIntE(i)
	return v
}

// ToStringMapInt64 casts an interface to a map[string]int64 type.
func ToStringMapInt64(i interface{}) map[string]int64 {
	v, _ := ToStringMapInt64E(i)
	return v
}

// ToStringMap casts an interface to a map[string]interface{} type.
func ToStringMap(i interface{}) map[string]interface{} {
	v, _ := ToStringMapE(i)
	return v
}

// ToSlice casts an interface to a []interface{} type.
func ToSlice(i interface{}) []interface{} {
	v, _ := ToSliceE(i)
	return v
}

// ToBoolSlice casts an interface to a []bool type.
func ToBoolSlice(i interface{}) []bool {
	v, _ := ToBoolSliceE(i)
	return v
}

// ToStringSlice casts an interface to a []string type.
func ToStringSlice(i interface{}) []string {
	v, _ := ToStringSliceE(i)
	return v
}

// ToIntSlice casts an interface to a []int type.
func ToIntSlice(i interface{}) []int {
	v, _ := ToIntSliceE(i)
	return v
}

// ToDurationSlice casts an interface to a []time.Duration type.
func ToDurationSlice(i interface{}) []time.Duration {
	v, _ := ToDurationSliceE(i)
	return v
}

func TransferByJson[Src any, Tar any](src Src) (Tar, error) {
	var target Tar
	data, err := json.Marshal(src)
	if err != nil {
		return target, err
	}

	if err := json.Unmarshal(data, &target); err != nil {
		return target, err
	}

	return target, nil
}

func JsonString(data any) (string, error) {
	d, err := json.Marshal(data)
	return string(d), err
}
