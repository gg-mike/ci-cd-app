package util

import (
	"errors"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

func Nothing[T any](ctx *gin.Context, db *gorm.DB, m T) error { return nil }
func NotImplemented[T any](ctx *gin.Context, db *gorm.DB, m T) error { return errors.New("not implemented") }
