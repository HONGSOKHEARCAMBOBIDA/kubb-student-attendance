package helper

import "github.com/gin-gonic/gin"

func GetParamID(c *gin.Context) (int, bool) {
	paramID, ok := c.Get("id")
	if !ok {
		return 0, false
	}
	id, ok := paramID.(float64)
	if !ok {
		return 0, false
	}
	return int(id), true
}
