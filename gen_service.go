package main

// 生成api的proto文件
////go:generate kratos proto add api/proto/chat.proto

// 生成 server 源码 (手动改下文件后缀 _service.go)
//go:generate kratos proto server ./api/proto/chat.proto -t app/internal/service
