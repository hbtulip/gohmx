package controller

import (
	"github.com/gin-gonic/gin"
	"net/http"
)

func InitRouter() {
	//r := controller.InitRouter()
	r := gin.Default()

	r.GET("/ping", func(c *gin.Context) {
		c.JSON(200, gin.H{ //c.JSON c.XML c.YAML
			"message": "pong",
		})
	})

	//静态页面
	r.Static("/js", "./static/js")
	r.LoadHTMLFiles("./static/index.html") //, "./static/js/main.js")
	r.GET("/index", func(c *gin.Context) {
		c.HTML(200, "index.html", nil)

	})

	v1 := r.Group("/v1")
	{

		v1.GET("/hello", HelloPage)
		v1.GET("/hello/:name", func(c *gin.Context) {
			name := c.Param("name")
			c.String(http.StatusOK, "Hello %s", name)
		})
		v1.GET("/hellohmx", func(c *gin.Context) {
			c.Request.URL.Path = "/v1/hello/hmx"
			r.HandleContext(c)
		})

	}

	//油机联动业务监控接口
	oil := r.Group("/oilmachine")
	{
		// 传统请求URL: /getTankStatus?storeid=1000&tankid=2
		oil.GET("/getTankStatus", GetTankStatus)

	}

	//定义默认路由
	r.NoRoute(func(c *gin.Context) {
		c.JSON(http.StatusNotFound, gin.H{
			"status": 404,
			"error":  "404, 页面未找到",
		})
	})

	r.Run(":8080")
}
