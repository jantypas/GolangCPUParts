package Pipes

import (
	"fmt"
	"strconv"
	"testing"
)

func TestCList_Extend(t *testing.T) {
	cl := NewCList()
	cl.WriteNBytes([]byte{1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 11, 12, 13, 14, 15, 16})
	n, b := cl.ReadNBytes(10)
	n, b = cl.ReadNBytes(10)
	fmt.Println("Read " + strconv.Itoa(n) + " chars")
	fmt.Println(b)
}
