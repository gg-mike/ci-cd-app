package dao

import (
	"encoding/json"
	"errors"
	"io"
	"net/http"

	"github.com/gg-mike/ci-cd-app/backend/internal/controller/util"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

type DAO[T, TCore, TCreate, TShort any] struct {
	DB           *gorm.DB

	BeforeCreate func(ctx *gin.Context, model TCreate) error
	AfterCreate  func(ctx *gin.Context, model TCreate) error
	BeforeUpdate func(ctx *gin.Context, model TCreate) error
	AfterUpdate  func(ctx *gin.Context, model TCreate) error
	BeforeDelete func(ctx *gin.Context, model T) error
	AfterDelete  func(ctx *gin.Context, model T) error

	PKCond       func(id uuid.UUID) T
}

func getID(params gin.Params) (uuid.UUID, error) {
	_id, ok := params.Get("id")
	if !ok {
		return uuid.UUID{}, errors.New(`missing param "id"`)
	}
	id, err := uuid.Parse(_id)
	if err != nil {
		return uuid.UUID{}, errors.New(`error parsing param "id"`)
	}
	return id, nil
}

func (dao *DAO[T, TCore, TCreate, TShort]) GetOne(ctx *gin.Context) {
	id, err := getID(ctx.Params)
	if err != nil {
		util.MessageResponse(ctx, http.StatusBadRequest, "error in params: %v", err)
		return
	}

	m := dao.PKCond(id)
	o := new(T)
	err = dao.DB.Model(*new(T)).Preload(clause.Associations).First(o, &m).Error
	switch err {
	case gorm.ErrRecordNotFound:
		util.MessageResponse(ctx, http.StatusNotFound, "no record found")
	case nil:
		ctx.JSON(http.StatusOK, o)
	default:
		util.MessageResponse(ctx, http.StatusInternalServerError, "database error: %v", err)
	}
}

func (dao *DAO[T, TCore, TCreate, TShort]) GetMany(ctx *gin.Context) {
	body, err := io.ReadAll(ctx.Request.Body)
	if err != nil {
		util.MessageResponse(ctx, http.StatusBadRequest, "error reading body: %v", err)
	}
	filter := new(T)

	if err = json.Unmarshal(body, filter); len(body) != 0 && err != nil {
		util.MessageResponse(ctx, http.StatusBadRequest, "error parsing body: %v", err)
	}

	offset, limit, order, err := Paginate(ctx)
	if err != nil {
		util.MessageResponse(ctx, http.StatusBadRequest, "error in query: %v", err)
	}

	o := new([]TShort)
	err = dao.DB.Model(new(T)).Offset(offset).Limit(limit).Order(order).Find(o, filter).Error
	switch err {
	case gorm.ErrRecordNotFound:
		util.MessageResponse(ctx, http.StatusNotFound, "no record found")
	case nil:
		ctx.JSON(http.StatusOK, o)
	default:
		util.MessageResponse(ctx, http.StatusInternalServerError, "database error: %v", err)
	}
}

func (dao *DAO[T, TCore, TCreate, TShort]) Create(ctx *gin.Context) {
	m := *new(TCreate)
	if err := dao.BeforeCreate(ctx, m); err != nil {
		util.MessageResponse(ctx, http.StatusInternalServerError, err.Error())
		return
	}

	// dao.DB.Model(new(T)).Create()
	
	if err := dao.AfterCreate(ctx, m); err != nil {
		util.MessageResponse(ctx, http.StatusInternalServerError, err.Error())
		return
	}
	util.MessageResponse(ctx, http.StatusNotImplemented, "")
}

func (dao *DAO[T, TCore, TCreate, TShort]) Update(ctx *gin.Context) {
	m := *new(TCreate)
	if err := dao.BeforeUpdate(ctx, m); err != nil {
		util.MessageResponse(ctx, http.StatusInternalServerError, err.Error())
		return
	}

	// dao.DB.Model(new(T)).Updates()
	
	if err := dao.AfterUpdate(ctx, m); err != nil {
		util.MessageResponse(ctx, http.StatusInternalServerError, err.Error())
		return
	}
	util.MessageResponse(ctx, http.StatusNotImplemented, "")
}

func (dao *DAO[T, TCore, TCreate, TShort]) Delete(ctx *gin.Context) {
	m := *new(T)
	if err := dao.BeforeDelete(ctx, m); err != nil {
		util.MessageResponse(ctx, http.StatusInternalServerError, err.Error())
		return
	}

	// dao.DB.Model(new(T)).Delete()
	
	if err := dao.AfterDelete(ctx, m); err != nil {
		util.MessageResponse(ctx, http.StatusInternalServerError, err.Error())
		return
	}
	util.MessageResponse(ctx, http.StatusNotImplemented, "")
}
