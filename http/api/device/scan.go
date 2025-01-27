package device

import (
	"github.com/gin-gonic/gin"
)

/*
*

	@author: victor2022
	@since: 2025/1/27
*/
func pong(c *gin.Context) {
	c.PostForm("senderId")
	c.PostForm("senderIp")
	c.PostForm("senderPort")

}
