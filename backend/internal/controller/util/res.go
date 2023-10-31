package util

import (
	"fmt"

	"github.com/gin-gonic/gin"
)

type Message struct {
	Message string `json:"message"`
}

func MessageResponse(ctx *gin.Context, status int, message string, a ...any) {
	ctx.JSON(status, Message{
		Message: fmt.Sprintf(message, a...),
	})
}
