package controller

import (
	"github.com/gin-gonic/gin"
	"log"
	"net/http"

	"gohmx/model"
)

func GetTankStatus(c *gin.Context) {
	storeid_s := c.DefaultQuery("storeid", "1000")
	tankid_s := c.Request.URL.Query().Get("tankid")

	data, err := model.GetTankStatus(storeid_s, tankid_s)
	if err != nil {
		log.Fatalln(err)
		return
	}

	// 打印返回值
	log.Println(data)

	c.JSON(http.StatusOK, gin.H{
		"code": 0,
		"data": data,
	})
}
