package cast

import (
	"fmt"
	"testing"
)

func TestSliceRemoveItem(t *testing.T) {
	list := []string{"a", "b", "d", "c", "d"}

	// fmt.Println(SliceRemoveItem(list, "a"))
	// fmt.Println(SliceRemoveItem(list, "c"))
	fmt.Println(RemoveSliceItem(list, "d"))

	l := []string{"a", "a"}
	fmt.Println(RemoveSliceItem(l, "a"))
}
