package network

import (
	"github.com/gin-gonic/gin"
	"websocket/types"
)

func response(c *gin.Context, s int, res interface{}, data ...string) {
	c.JSON(s, types.NewRes(s, res, data...))
}
