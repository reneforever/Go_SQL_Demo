package main

import (
	"database/sql"
	"fmt"
	"strings"

	_ "github.com/go-sql-driver/mysql"
)

type user struct {
	id   int
	name string
	age  int
}

var db *sql.DB // 连接池对象

func initDB() (err error) {
	dsn := "root:123456@tcp(127.0.0.1:3306)/sql_test"
	db, err = sql.Open("mysql", dsn)
	if err != nil {
		return err
	}
	err = db.Ping()
	if err != nil {
		return err
	}
	fmt.Println("database connected")
	return nil
}

func queryOne(id int) (err error) {
	// 单条查询记录
	var u1 user
	var sqlStr = `select id, name, age from user where id=?;`
	// scan自带关闭数据库连接能力
	queryErr := db.QueryRow(sqlStr, id).Scan(&u1.id, &u1.name, &u1.age)
	if queryErr != nil {
		return queryErr
	}
	return nil
}

func queryMore() (userList []interface{}, err error) {
	var sqlStr = `select id, name, age from user where id>?;`
	rows, err := db.Query(sqlStr, 0)
	if err != nil {
		fmt.Println("query failed", err)
	}
	defer rows.Close()
	for rows.Next() {
		var u user
		err := rows.Scan(&u.id, &u.name, &u.age)
		if err != nil {
			fmt.Println("scan failed", err)
			return nil, err
		}
		userList = append(userList, u)
	}
	return userList, nil
}

func insert(name string, age int8) (err error) {
	var build strings.Builder
	build.WriteString("INSERT INTO `user`(`name`, `age`) VALUES(")
	build.WriteString("\"")
	build.WriteString(name)
	build.WriteString("\"")
	build.WriteString(",")
	build.WriteString(fmt.Sprintf("%d", age))
	build.WriteString(")")
	sqlStr := build.String()
	fmt.Println("sqlStr", sqlStr)
	result, err := db.Exec(sqlStr)
	if err != nil {
		return err
	}
	fmt.Println("insert result is", result)
	return nil
}

func main() {
	initErr := initDB()
	if initErr != nil {
		fmt.Println("init DB connect failed", initErr)
		return
	}
	// 关闭数据库
	defer func() {
		err := db.Close()
		if err != nil {
			fmt.Println("database connect close failed", err)
			return
		}
		fmt.Println("database connect closed")
	}()
	fmt.Println("database open sccess")
	insertErr := insert("cc", 24)
	if insertErr != nil {
		fmt.Println("insert failed", insertErr)
	}
	queryErr := queryOne(1)
	if queryErr != nil {
		fmt.Println("query failed", queryErr)
		return
	}
	userList, queryMoreErr := queryMore()
	if queryMoreErr != nil {
		fmt.Println("query more failed")
		return
	}
	fmt.Println(userList)
}
