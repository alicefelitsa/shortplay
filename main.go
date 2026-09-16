package main

import "fmt"

//go env -w CGO_ENABLED=0
//go env -w GOOS=linux
//go env -w GOOS=windows
//go env -w GOPROXY=https://goproxy.cn,direct

func main() {
	fmt.Println("Hello, world!")
}
