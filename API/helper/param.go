package helper

import (
	"strconv"

	"github.com/gin-gonic/gin"
)

func GetParamID(c *gin.Context) (int, bool) {
	paramID := c.Param("id")

	id, err := strconv.Atoi(paramID)
	if err != nil {
		return 0, false
	}

	return id, true
}
