package main

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"
)

func main() {
	r := gin.Default()

	r.GET("/", func(ctx *gin.Context) {
		ctx.String(200, "Hello World")
	})

	// 将 request body 绑定到不同的结构体中
	r.GET("/someJSON", func(ctx *gin.Context) {
		data := map[string]interface{}{
			"name": "china 你好呀",
		}
		ctx.AsciiJSON(http.StatusOK, data)
	})

	r.POST("/someJSON", SomeHandler)

	r.POST("/default", DefaultHandler)

	r.GET("bind-uri/:name/:age", BindUriHandler)

	r.POST("/login", LoginHandler)

	err := r.Run(":8080")
	if err != nil {
		panic(err)
	}
}

// ================= 将 request body 绑定到不同的结构体中 =================
type formA struct {
	Foo string `form:"foo" json:"foo"`
}

type formB struct {
	Bar string `form:"bar" json:"bar"`
}

// 要想多次绑定，可以使用 c.ShouldBindBodyWith.
func SomeHandler(c *gin.Context) {
	a := formA{}
	b := formB{}

	// 绑定JSON
	if errA := c.ShouldBindBodyWith(&a, binding.JSON); errA != nil {
		c.JSON(http.StatusOK, "A err:"+errA.Error())
		return
	} else if errB := c.ShouldBindBodyWith(&b, binding.JSON); errB != nil {
		c.JSON(http.StatusOK, "B err:"+errB.Error())
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"a": a,
		"b": b,
	})
}

// ================= 绑定表单字段的默认值 =================
type Person struct {
	Name   string   `form:"name,default=张三" json:"name" uri:"name" binding:"required"`
	Age    int      `form:"age,default=18" uri:"age"`
	Frinds []string `form:"frinds,default=张三,lisi"`
}

func DefaultHandler(c *gin.Context) {
	var person Person
	if err := c.ShouldBind(&person); err != nil {
		c.JSON(http.StatusOK, gin.H{
			"err": err.Error(),
		})
		return
	}
	c.JSON(http.StatusOK, person)
}

// ================= 绑定URL参数=================
func BindUriHandler(c *gin.Context) {
	var person Person
	if err := c.ShouldBindUri(&person); err != nil {
		c.JSON(http.StatusOK, gin.H{
			"err": err.Error(),
		})
		return
	}
	c.JSON(http.StatusOK, person)
}

// ================= 模型绑定和验证=================
type Login struct {
	Username string `form:"username" json:"username" binding:"required"`
	Password string `form:"password" json:"password" binding:"required"`
}

func LoginHandler(c *gin.Context) {
	login := Login{}
	if err := c.ShouldBind(&login); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"err": err.Error(),
		})
		return
	}

	if login.Username != "admin" || login.Password != "admin" {
		c.JSON(http.StatusUnauthorized, gin.H{
			"msg": "用户名或者密码错误",
		})
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"msg": "登录成功",
	})
}
