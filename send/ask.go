package send

import (
	"database/sql"
	"fmt"
	_ "github.com/go-sql-driver/mysql"
)

func Ask(xx string, yy string, zz string) int {
	var db *sql.DB
	var err error
	db, err = sql.Open("mysql", "root:1592933843zzz@tcp(127.0.0.1:3306)/sql_find")

	if err != nil {
		fmt.Printf("connect to db 127.0.0.1:3306 error: %v\n", err)
		return 0
	}
	sqlStr := "select id, name, password from user where id > ?"
	rows, err := db.Query(sqlStr, 0)
	if err != nil {
		fmt.Printf("query failed, err:%v\n", err)
		return 0
	}
	defer rows.Close()
	var flag int
	flag = 0
	// 循环读取结果集中的数据
	for rows.Next() {
		var u User
		err := rows.Scan(&u.Id, &u.Name, &u.Password)
		if err != nil {
			fmt.Printf("scan failed, err:%v\n", err)
			return 0
		}
		fmt.Println(u.Id)
		//x:=strconv.Itoa(u.id)
		//x := string(u.id)
		fmt.Println(xx)
		fmt.Println(u.Name)
		fmt.Println(u.Password)
		if xx == u.Id && u.Name == yy && zz == u.Password {

			flag = 1
		}

	}
	return flag
}
