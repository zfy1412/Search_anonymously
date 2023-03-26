package main

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"goweb/send"
	"net/http"
)

func main() {

	r := gin.Default()
	r.LoadHTMLGlob(
		"views/HTML/*",
	)
	r.StaticFS("/views", http.Dir("./views"))
	r.GET("/login", func(c *gin.Context) {
		c.HTML(http.StatusOK, "Login.html", nil)
	})
	r.GET("/", func(c *gin.Context) {
		c.Redirect(http.StatusMovedPermanently, "http://127.0.0.1:7995/login")
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
			c.Redirect(http.StatusMovedPermanently, "http://127.0.0.1:7995/map")
		} else {
			//send.Insert(id,username,password)
			//c.JSON(http.StatusOK, gin.H{"token": nil})
			c.Redirect(http.StatusMovedPermanently, "http://127.0.0.1:7995/login")
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
		//ca := compute.Count(now)
		//c.HTML(http.StatusOK, "state_red.html", nil)
		c.JSON(302, gin.H{"location": "http://127.0.0.1:7995/sate3"})

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
		data1 := gin.H{
			"taxiDriver": "libro",
			"phone":      "18336218165",
			"ID":         "辽A1FD651",
		}
		data2 := gin.H{
			"taxiDriver": "hhh",
			"phone":      "11111111",
			"ID":         "辽A1Fd1we",
		}
		data := gin.H{
			"length": "3",
			"data1":  data1,
			"data2":  data2,
			"data3":  nil,
		}
		//这里要返回json格式的数据，所以用c.JSON,这样，数据就返回给请求方了
		c.JSON(http.StatusOK, data)
	})
	r.GET("/sate3", func(c *gin.Context) {
		//c.JSON(200, gin.H{"data": "1"})
		c.HTML(http.StatusOK, "state_green.html", nil)

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
	r.Run("127.0.0.1:7995")
}
