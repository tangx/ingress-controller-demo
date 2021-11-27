package main

import (
	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
)

func main() {
	r := gin.Default()

	r.Any("/ping", func(c *gin.Context) {
		for k, v := range c.Request.Header {
			logrus.Printf("%s => %v", k, v)
		}
		c.String(200, "pong")
	})

	_ = r.Run(":8099")
}
