package main

import "fmt"

func main() {
	dev()
}

func dev() {
	arr := [...]int32{1,2,3}
	fmt.Println(arr)

	slice := []int32{1,2,3,4,5}
	fmt.Println(slice)
	slice = append(slice, 7)
	fmt.Println(slice)

	sliceTest := make([]int, 3)
	fmt.Println(sliceTest)
}