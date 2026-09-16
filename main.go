package main

// ================================
// 交叉编译说明（本文件仅保留注释，实际入口在 web.go）
// ================================
// 设置交叉编译参数，用于生成不同平台下的可执行程序。
//
//go env -w CGO_ENABLED=0
//go env -w GOOS=linux
//go env -w GOARCH=amd64
//go env -w GOARM=6
//
// Linux 平台编译：
//go build -o shortplay web.go
//
// Windows 平台编译：
//go env -w GOOS=windows
//go build -o shortplay.exe web.go
//
// 切换回 Windows 本地开发：
//go env -w GOOS=windows
