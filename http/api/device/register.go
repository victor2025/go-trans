package device

import "github.com/gin-gonic/gin"

/*
*

	@author: victor2022
	@since: 2025/1/27
*/
func RegisterRouter(router *gin.RouterGroup) {
	router.GET("ping", ping)
}
