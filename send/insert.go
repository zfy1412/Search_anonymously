package send
import (
	"database/sql"
	"fmt"
	_ "github.com/go-sql-driver/mysql"
)

func Insert( xx string, yy string, zz string)(int) {
	var err error
	var db *sql.DB
	db, err = sql.Open("mysql", "root:123@tcp(127.0.0.1:3306)/sql_user")

	//defer db.Close()token.go

	if err != nil {
		fmt.Printf("connect to db 127.0.0.1:3306 error: %v\n", err)
		return 0
	}
	sqlStr := "insert into user(id, name, password) values (?,?,?)"
	ret, err := db.Exec(sqlStr,xx, yy, zz)
	if err != nil {
		fmt.Printf("insert failed, err:%v\n", err)
		return 0
	}
	theID, err := ret.LastInsertId() // 新插入数据的id
	if err != nil {
		fmt.Printf("get lastinsert ID failed, err:%v\n", err)
		return 0
	}
	fmt.Printf("insert success, the id is %d.\n", theID)
	return 1
}