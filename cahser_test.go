package main

import (
	"fmt"
	"testing"
	"unsafe"
)

func removeLastValue(slice []string, value string) []string {
	// 从切片末尾开始查找指定值元素的索引
	for i := len(slice) - 1; i >= 0; i-- {
		if slice[i] == value {
			// 将找到的指定值元素从切片中删除
			copy(slice[i:], slice[i+1:])
			// 调整切片的长度
			slice = slice[:len(slice)-1]
			break
		}
	}
	return slice
}
func TestRemove(t *testing.T) {
	a := []string{"a", "b", "c"}
	fmt.Println(unsafe.Pointer(&a))
	fmt.Println(a)
	a = removeLastValue(a, "b")
	fmt.Println(unsafe.Pointer(&a))
	fmt.Println(a)
}
