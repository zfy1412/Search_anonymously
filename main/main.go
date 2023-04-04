package main

import (
	"database/sql"
	"fmt"
	"github.com/gin-gonic/gin"
	"goweb/compute"
	"goweb/send"
	"net/http"
)

func main() {
	//encryption.Test_homomorphicCrypto()
	ans := []string{}
	r := gin.Default()
	r.LoadHTMLGlob(
		"views/HTML/*",
	)
	r.StaticFS("/views", http.Dir("./views"))
	r.GET("/login", func(c *gin.Context) {
		c.HTML(http.StatusOK, "Login.html", nil)
	})
	r.GET("/", func(c *gin.Context) {
		c.Redirect(http.StatusMovedPermanently, "/login")
	})
	r.POST("/login", func(c *gin.Context) {
		id := c.PostForm("idnumber")
		username := c.PostForm("username")
		password := c.PostForm("password")
		user := &send.User{
			Id:       id,
			Name:     username,
			Password: password,
		}
		flag := send.Ask(id, username, password)
		if flag == 1 {
			key, _ := send.GenerateToken(*user)
			fmt.Println(key)
			//c.JSON(http.StatusOK, gin.H{"token": key})
			c.Redirect(http.StatusMovedPermanently, "/map")
		} else {
			//send.Insert(id,username,password)
			//c.JSON(http.StatusOK, gin.H{"token": nil})
			c.Redirect(http.StatusMovedPermanently, "/login")
		}
	})
	r.GET("/map", func(c *gin.Context) {
		c.HTML(http.StatusOK, "map.html", nil)
	})
	r.POST("/map", func(c *gin.Context) {
		var now send.Date
		err := c.BindJSON(&now)
		fmt.Println(err)
		fmt.Println(now)
		fmt.Println(now.Etime)
		fmt.Println(len(now.Node))
		ans = compute.Count(now)
		//c.HTML(http.StatusOK, "state_red.html", nil)
		c.JSON(302, gin.H{"location": "/sate3"})

		//c.JSON(302, gin.H{"location": "http://1.15.146.175:7995/sate1"})
		//c.Redirect(302, "http://10.0.4.15:7995/sate1")
		//if err != nil {
		//	//log.Info(err)
		//	c.JSON(200, gin.H{"errcode": 400, "description": "Post Data Err"})
		//	return
		//} else {
		//	//fmt.Println(reqInfo.Data)
		//}
	})
	r.GET("/json", func(c *gin.Context) {
		var db *sql.DB
		var err error
		db, err = sql.Open("mysql", "root:1592933843zzz@tcp(127.0.0.1:3306)/sql_find")

		if err != nil {
			fmt.Printf("connect to db 127.0.0.1:3306 error: %v\n", err)
		}

		length := len(ans)
		i := 0
		var d1 send.Driver
		var d2 send.Driver
		var d3 send.Driver
		for i < length {
			if i == 0 {
				d1.Uid = ans[0]
				sqlStr := "select id, name, phone from driver where id = ?"
				rows, err := db.Query(sqlStr, d1.Uid)
				if err != nil {
					fmt.Printf("query failed, err:%v\n", err)
				}
				defer rows.Close()
				for rows.Next() {
					err := rows.Scan(&d1.Uid, &d1.Name, &d1.Phone)
					if err != nil {
						fmt.Printf("scan failed, err:%v\n", err)
					}
				}

				i++
			} else if i == 1 {
				d2.Uid = ans[1]
				sqlStr := "select id, name, phone from driver where id = ?"
				rows, err := db.Query(sqlStr, d2.Uid)
				if err != nil {
					fmt.Printf("query failed, err:%v\n", err)
				}
				defer rows.Close()
				for rows.Next() {
					err := rows.Scan(&d2.Uid, &d2.Name, &d2.Phone)
					if err != nil {
						fmt.Printf("scan failed, err:%v\n", err)
					}
				}

				i++
			} else if i == 2 {
				d3.Uid = ans[2]
				sqlStr := "select id, name, phone from driver where id = ?"
				rows, err := db.Query(sqlStr, d3.Uid)
				if err != nil {
					fmt.Printf("query failed, err:%v\n", err)
				}
				defer rows.Close()
				for rows.Next() {
					err := rows.Scan(&d3.Uid, &d3.Name, &d3.Phone)
					if err != nil {
						fmt.Printf("scan failed, err:%v\n", err)
					}
				}

				i++
			}
		}

		data1 := gin.H{
			"taxiDriver": d1.Name,
			"phone":      d1.Phone,
			"ID":         d1.Uid,
		}
		data2 := gin.H{
			"taxiDriver": d2.Name,
			"phone":      d2.Phone,
			"ID":         d2.Uid,
		}
		data3 := gin.H{
			"taxiDriver": d3.Name,
			"phone":      d3.Phone,
			"ID":         d3.Uid,
		}
		if i == 0 {
			data := gin.H{
				"length": length,
				"data1":  data1,
				"data2":  nil,
				"data3":  nil,
			}
			c.JSON(http.StatusOK, data)
		} else if i == 1 {
			data := gin.H{
				"length": length,
				"data1":  data1,
				"data2":  data2,
				"data3":  nil,
			}
			c.JSON(http.StatusOK, data)
		} else {
			data := gin.H{
				"length": length,
				"data1":  data1,
				"data2":  data2,
				"data3":  data3,
			}
			c.JSON(http.StatusOK, data)
		}

	})
	r.GET("/sate3", func(c *gin.Context) {
		//c.JSON(200, gin.H{"data": "1"})
		c.HTML(http.StatusOK, "state_green.html", nil)

	})
	r.GET("/register", func(c *gin.Context) {
		c.HTML(http.StatusOK, "register.html", nil)

	})
	r.POST("/register", func(c *gin.Context) {
		id := c.PostForm("idnumber")
		username := c.PostForm("username")
		password := c.PostForm("password")
		user := &send.User{
			Id:       id,
			Name:     username,
			Password: password,
		}
		db, err := sql.Open("mysql", "root:1592933843zzz@tcp(127.0.0.1:3306)/sql_find")

		if err != nil {
			fmt.Printf("connect to db 127.0.0.1:3306 error: %v\n", err)
		}
		sqlStr := "insert into user(id, password,name) values(?,?,?)"
		_, errr := db.Exec(sqlStr, user.Id, user.Password, user.Name)
		if errr != nil {
			fmt.Println("插入数据错误", err)
		}
		c.Redirect(http.StatusMovedPermanently, "/login")
	})
	//r.POST("/register", func(c *gin.Context) {
	//	id := c.PostForm("id")
	//	username := c.PostForm("username")
	//	password := c.PostForm("password")
	//	mark := send.Insert(id, username, password)
	//	if mark == 1 {
	//		c.Redirect(http.StatusMovedPermanently, "http://127.0.0.1:8080/login")
	//		//c.Request.URL.Path = "/login"
	//		//r.HandleContext(c)
	//	}
	//})
	r.Run(":80")
}
