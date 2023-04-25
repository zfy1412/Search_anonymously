package main

import (
	"database/sql"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
	"goweb/compute"
	"goweb/send"
	"net/http"
	"strconv"
	"time"
)

var upgrader = websocket.Upgrader{}

type Message struct {
	Name    string `json:"name"`
	Content string `json:"content"`
}

var clients = make(map[*websocket.Conn]bool)
var broadcast = make(chan Message)

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
		username := c.PostForm("username")
		password := c.PostForm("password")
		phone := c.PostForm("phoneNumber")
		flag := send.Ask(username, password)
		if flag == 1 {
			//key, _ := send.GenerateToken(*user)
			//fmt.Println(key)
			var db *sql.DB
			var err error
			db, err = sql.Open("mysql", "root:1592933843zzz@tcp(127.0.0.1:3306)/sql_find")
			if err != nil {
				fmt.Printf("connect to db 127.0.0.1:3306 error: %v\n", err)
			}
			sqlStr := "select uid from user where name = ?"
			rows, err := db.Query(sqlStr, username)
			var idd int
			for rows.Next() {
				err := rows.Scan(&idd)
				if err != nil {
					fmt.Printf("scan failed, err:%v\n", err)
				}
			}
			str := strconv.Itoa(idd)
			c.SetCookie("name", str, 24*60*60, "/", "localhost", false, false)
			c.SetCookie("username", username, 24*60*60, "/", "localhost", false, false)
			c.SetCookie("phone", phone, 24*60*60, "/", "localhost", false, false)
			c.Redirect(http.StatusMovedPermanently, "/userMain")
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
		val, _ := c.Cookie("name")
		for i < length {
			if i == 0 {
				d1.Uid = ans[0]
				sqlStr := "select carid, name, phone from driver where did = ?"
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
				sqlStr := "select carid, name, phone from driver where did = ?"
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
				sqlStr := "select carid, name, phone from driver where did = ?"
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
		timeStr := time.Now().Format("2006-01-02 15:04:05")
		if i == 1 {
			data := gin.H{
				"length": length,
				"data1":  data1,
				"data2":  nil,
				"data3":  nil,
			}
			sqlStr := "insert into history(uid,time ,did) values(?,?,?)"
			_, errr := db.Exec(sqlStr, val, timeStr, ans[0])
			if errr != nil {
				fmt.Println("插入数据错误", err)
			}
			c.JSON(http.StatusOK, data)
		} else if i == 2 {
			data := gin.H{
				"length": length,
				"data1":  data1,
				"data2":  data2,
				"data3":  nil,
			}

			sqlStr := "insert into history(uid,time ,did) values(?,?,?)"
			_, errr := db.Exec(sqlStr, val, timeStr, ans[0])
			if errr != nil {
				fmt.Println("插入数据错误", err)
			}
			sqlStr = "insert into history(uid,time ,did) values(?,?,?)"
			_, errr = db.Exec(sqlStr, val, timeStr, ans[1])
			if errr != nil {
				fmt.Println("插入数据错误", err)
			}
			c.JSON(http.StatusOK, data)
		} else {
			sqlStr := "insert into history(uid,time ,did) values(?,?,?)"
			_, errr := db.Exec(sqlStr, val, timeStr, ans[0])
			if errr != nil {
				fmt.Println("插入数据错误", err)
			}
			sqlStr = "insert into history(uid,time ,did) values(?,?,?)"
			_, errr = db.Exec(sqlStr, val, timeStr, ans[1])
			if errr != nil {
				fmt.Println("插入数据错误", err)
			}
			sqlStr = "insert into history(uid,time ,did) values(?,?,?)"
			_, errr = db.Exec(sqlStr, val, timeStr, ans[2])
			if errr != nil {
				fmt.Println("插入数据错误", err)
			}
			data := gin.H{
				"length": length,
				"data1":  data1,
				"data2":  data2,
				"data3":  data3,
			}
			c.JSON(http.StatusOK, data)
		}

	})
	r.GET("/historydata", func(c *gin.Context) {
		var db *sql.DB
		var err error
		db, err = sql.Open("mysql", "root:1592933843zzz@tcp(127.0.0.1:3306)/sql_find")

		if err != nil {
			fmt.Printf("connect to db 127.0.0.1:3306 error: %v\n", err)
		}
		var h send.History
		var j send.Jsonhistory
		data := []send.Jsonhistory{}
		val, _ := c.Cookie("name")
		fmt.Printf(val)
		sqlStr := "select uid,did,time  from history where uid = ?"
		rows, err := db.Query(sqlStr, val)
		for rows.Next() {
			fmt.Printf("123456489")
			err := rows.Scan(&h.Uid, &h.Did, &h.Time)
			if err != nil {
				fmt.Printf("scan failed, err:%v\n", err)
			}
			j.Time = h.Time
			var d send.Driver
			sqlStr := "select carid,name,phone  from driver where did = ?"
			rrows, errr := db.Query(sqlStr, h.Did)
			if errr != nil {
				fmt.Printf("scan failed, err:%v\n", err)
			}
			for rrows.Next() {
				errr := rrows.Scan(&d.Uid, &d.Name, &d.Phone)
				if errr != nil {
					fmt.Printf("scan failed, err:%v\n", errr)
				}
			}
			j.Info = d
			fmt.Printf(d.Phone)
			data = append(data, j)
		}

		c.JSON(http.StatusOK, data)

	})
	r.GET("/sate3", func(c *gin.Context) {
		//c.JSON(200, gin.H{"data": "1"})
		c.HTML(http.StatusOK, "result.html", nil)

	})
	r.GET("/register", func(c *gin.Context) {
		c.HTML(http.StatusOK, "register.html", nil)

	})
	r.GET("/userMain", func(c *gin.Context) {
		c.HTML(http.StatusOK, "user.html", nil)

	})
	r.GET("/messageHistory", func(c *gin.Context) {
		c.HTML(http.StatusOK, "messageHistory.html", nil)

	})
	r.GET("/result_h", func(c *gin.Context) {
		c.HTML(http.StatusOK, "result_h.html", nil)

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
		sqlStr := "insert into user(phone, password,name) values(?,?,?)"
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
	r.GET("/chat", func(c *gin.Context) {
		c.HTML(http.StatusOK, "chat.html", nil)
	})

	r.GET("/ws", func(c *gin.Context) {
		conn, err := upgrader.Upgrade(c.Writer, c.Request, nil)
		if err != nil {
			panic(err)
		}

		clients[conn] = true

		for {
			var message Message
			err := conn.ReadJSON(&message)
			if err != nil {
				delete(clients, conn)
				break
			}

			broadcast <- message
		}
	})

	go func() {
		for {
			message := <-broadcast

			for client := range clients {
				err := client.WriteJSON(message)
				if err != nil {
					client.Close()
					delete(clients, client)
				}
			}
		}
	}()

	r.Run(":8080")
}
