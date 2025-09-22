package utils

func Ptr[T any](v T) *T {
	return &v
}

func PtrSlice[T any](vs []T) []*T {
	ps := make([]*T, len(vs))
	for i, v := range vs {
		vv := v
		ps[i] = &vv
	}
	return ps
}

func PtrMap[K comparable, V any](vs map[K]V) map[K]*V {
	ps := make(map[K]*V, len(vs))
	for k, v := range vs {
		vv := v
		ps[k] = &vv
	}
	return ps
}

func IsNilOrZero[T comparable](v *T) bool {
	var zero T
	return v == nil || *v == zero
}

// 判断指针指向的值是否和给定值相等
func IsPtrEqualValue[T comparable](ptr *T, value T) bool {
	if IsNilOrZero(ptr) {
		var zero T
		return value == zero
	}
	return *ptr == value
}

// 判断两个指针是否指向同一个值
func IsPtrEqual[T comparable](ptr1, ptr2 *T) bool {
	if IsNilOrZero(ptr1) && IsNilOrZero(ptr2) {
		return true
	}
	if ptr1 == nil || ptr2 == nil {
		return false
	}
	return *ptr1 == *ptr2
}

// 将字节数组指针转换为字符串指针
func ToStringPtr(v []byte) *string {
	s := string(v)
	return &s
}
