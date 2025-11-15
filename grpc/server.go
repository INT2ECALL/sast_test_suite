package main

import (
	"net"
	"net/rpc"
	"os/exec"
)

// 定义RPC服务结构
type CommandService struct{}

// 危险方法：接收用户输入的命令并直接执行（无过滤）
// 输入：cmdStr（用户传入的命令字符串）
// 输出：命令执行结果
func (c *CommandService) ExecCommand(cmdStr string, reply *string) error {
	// 直接执行用户传入的命令（高危！）
	output, err := exec.Command("bash", "-c", cmdStr).CombinedOutput()
	if err != nil {
		*reply = "执行错误: " + err.Error()
		return nil
	}
	*reply = string(output)
	return nil
}

func main() {
	// 注册RPC服务
	_ = rpc.Register(new(CommandService))
	// 监听端口
	listener, err := net.Listen("tcp", ":8080")
	if err != nil {
		panic(err)
	}
	defer listener.Close()
	// 启动服务
	for {
		conn, err := listener.Accept()
		if err != nil {
			continue
		}
		go rpc.ServeConn(conn)
	}
}
