package main

import (
	"database/sql"
	"fmt"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	_ "github.com/go-sql-driver/mysql"
)

// 数据库连接（替换为实际配置）
func initDB() *sql.DB {
	dsn := "root:password@tcp(127.0.0.1:3306)/testdb?charset=utf8mb4"
	db, err := sql.Open("mysql", dsn)
	if err != nil {
		panic("数据库连接失败: " + err.Error())
	}
	return db
}

var db = initDB()

// 风险点：间接拼接用户输入（可能规避SAST简单规则）
func buildQuery(username string) string {
	// 模拟“看似处理”实际未过滤的场景：去除空格（但单引号等特殊字符未处理）
	cleaned := strings.TrimSpace(username)
	// 拼接SQL（核心漏洞：用户输入直接进入SQL语句）
	return fmt.Sprintf("SELECT password FROM users WHERE username = '%s'", cleaned)
}

// Gin接口：用户登录查询（存在注入漏洞）
func loginHandler(c *gin.Context) {
	username := c.Query("username") // 接收用户输入的用户名
	if username == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "用户名不能为空"})
		return
	}

	// 构建SQL查询（间接拼接，可能绕过SAST基础检测）
	sqlStr := buildQuery(username)
	fmt.Println("执行的SQL:", sqlStr) // 打印SQL便于观察

	var password string
	// 执行查询（未使用参数化）
	err := db.QueryRow(sqlStr).Scan(&password)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "查询失败"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"msg": "查询成功", "password": password})
}

func main() {
	r := gin.Default()
	r.GET("/login", loginHandler) // 注册存在漏洞的接口
	r.Run(":8080")                // 启动服务
}
