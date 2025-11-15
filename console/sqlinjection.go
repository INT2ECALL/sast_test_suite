package main

import (
	"database/sql"
	"fmt"
	"log"

	_ "github.com/go-sql-driver/mysql"
)

// 模拟数据库连接（请替换为自己的数据库信息）
func getDB() *sql.DB {
	dsn := "root:password@tcp(127.0.0.1:3306)/testdb?charset=utf8mb4"
	db, err := sql.Open("mysql", dsn)
	if err != nil {
		log.Fatal("数据库连接失败:", err)
	}
	return db
}

// 危险方法：直接拼接用户输入到SQL语句（存在注入漏洞）
func getUserByUsername(username string) (string, error) {
	db := getDB()
	defer db.Close()

	// 漏洞点：直接拼接用户输入，未使用参数化查询
	sqlStr := fmt.Sprintf("SELECT password FROM users WHERE username = '%s'", username)
	fmt.Println("执行的SQL:", sqlStr) // 打印SQL便于观察

	var password string
	// 执行查询
	err := db.QueryRow(sqlStr).Scan(&password)
	if err != nil {
		return "", err
	}
	return password, nil
}

func main() {
	// 正常输入（预期行为）
	// username := "alice"

	// 恶意输入（SQL注入攻击）：通过单引号闭合原查询，添加OR条件使查询恒真
	// 结果会返回表中第一个用户的密码（或所有用户，取决于表结构）
	username := "alice' OR '1'='1"

	password, err := getUserByUsername(username)
	if err != nil {
		log.Println("查询错误:", err)
		return
	}
	fmt.Printf("查询到的密码: %s\n", password)
}
