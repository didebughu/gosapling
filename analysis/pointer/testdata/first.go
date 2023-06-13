package testdata

import "fmt"

func test() {
	var ptr *int
	number := *ptr
	fmt.Println(number, *ptr)
}
