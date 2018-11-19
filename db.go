package main

import (
	"database/sql"
	_ "github.com/go-sql-driver/mysql"
	"log"
)

func GetEmailListBySql(dbconfig string) []string {
	// 连接数据库
	db, err := sql.Open("mysql", dbconfig)
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	// 检查是否可连接
	err = db.Ping()
	if err != nil {
		log.Fatal(err)
	}

	rows, err := db.Query("SELECT IFNULL(account, 0) FROM users;")
	if err != nil {
		log.Fatal(err)
	}

	var arr []string
	var item string
	// 逐行读取
	for rows.Next() {
		rows.Columns()
		err := rows.Scan(&item)
		if err != nil {
			log.Fatal(err)
		}
		arr = append(arr, item)
		//fmt.Println(item)
	}
	return arr
}
