package dao

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"

	"github.com/gg-mike/ci-cd-app/backend/internal/controller/util"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

type DAO[T, TShort any] struct {
	DB *gorm.DB

	PKCond func(id uuid.UUID) T
	Filter func(ctx *gin.Context) (map[string]any, error)
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

func getRecord[T any](db *gorm.DB, m *T) (int, error) {
	err := db.Model(new(T)).Preload(clause.Associations).Where("deleted = ?", false).First(m).Error
	switch err {
	case nil:
		return http.StatusOK, nil
	case gorm.ErrRecordNotFound:
		return http.StatusNotFound, errors.New("no record found")
	default:
		return http.StatusInternalServerError, fmt.Errorf("database error: %v", err)
	}
}

func getBody[T any](ctx *gin.Context, raw *map[string]any, m *T) error {
	if err := ctx.BindJSON(raw); err != nil {
		return err
	}
	buf := new(bytes.Buffer)
	if err := json.NewEncoder(buf).Encode(raw); err != nil {
		return err
	}
	if err := json.NewDecoder(buf).Decode(m); err != nil {
		return err
	}
	return nil
}

func (dao *DAO[T, TShort]) GetOne(ctx *gin.Context) {
	id, err := getID(ctx.Params)
	if err != nil {
		util.MessageResponse(ctx, http.StatusBadRequest, "error in params: %v", err)
		return
	}

	m := dao.PKCond(id)
	code, err := getRecord[T](dao.DB, &m)
	if err != nil {
		util.MessageResponse(ctx, code, err.Error())
	} else {
		ctx.JSON(code, m)
	}
}

func (dao *DAO[T, TShort]) GetMany(ctx *gin.Context) {
	filters, err := dao.Filter(ctx)
	if err != nil {
		util.MessageResponse(ctx, http.StatusBadRequest, "error parsing query: %v", err)
	}

	offset, limit, order, err := Paginate(ctx)
	if err != nil {
		util.MessageResponse(ctx, http.StatusBadRequest, "error in query: %v", err)
	}

	o := new([]TShort)
	_db := dao.DB.Model(new(T))
	for key, value := range filters {
		_db = _db.Where(key, value)
	}
	err = _db.Offset(offset).Limit(limit).Order(order).Find(o).Error
	switch err {
	case nil:
		ctx.JSON(http.StatusOK, o)
	case gorm.ErrRecordNotFound:
		util.MessageResponse(ctx, http.StatusNotFound, "no record found")
	default:
		util.MessageResponse(ctx, http.StatusInternalServerError, "database error: %v", err)
	}
}

func (dao *DAO[T, TShort]) Create(ctx *gin.Context) {
	raw := map[string]any{}
	m := *new(T)
	if err := getBody[T](ctx, &raw, &m); err != nil {
		util.MessageResponse(ctx, http.StatusBadRequest, "error parsing body: %v", err)
		return
	}

	err := dao.DB.Model(new(T)).InstanceSet("obj", raw).Create(&m).Error
	switch err {
	case nil:
		ctx.JSON(http.StatusCreated, m)
	case gorm.ErrForeignKeyViolated:
		util.MessageResponse(ctx, http.StatusConflict, "incorrect foreign key")
	default:
		util.MessageResponse(ctx, http.StatusInternalServerError, "database error: %v", err)
	}
}

func (dao *DAO[T, TShort]) Update(ctx *gin.Context) {
	id, err := getID(ctx.Params)
	if err != nil {
		util.MessageResponse(ctx, http.StatusBadRequest, "error in params: %v", err)
		return
	}

	prev := dao.PKCond(id)
	code, err := getRecord[T](dao.DB, &prev)
	if err != nil {
		util.MessageResponse(ctx, code, err.Error())
	}
	m := prev
	raw := map[string]any{}
	if err := getBody[T](ctx, &raw, &m); err != nil {
		util.MessageResponse(ctx, http.StatusBadRequest, "error parsing body: %v", err)
		return
	}

	err = dao.DB.Model(&prev).InstanceSet("obj", raw).InstanceSet("prev", prev).Updates(m).Error
	switch err {
	case nil:
		ctx.JSON(http.StatusOK, m)
	case gorm.ErrForeignKeyViolated:
		util.MessageResponse(ctx, http.StatusConflict, "incorrect foreign key")
	default:
		util.MessageResponse(ctx, http.StatusInternalServerError, "database error: %v", err)
	}
}

func (dao *DAO[T, TShort]) Delete(ctx *gin.Context) {
	id, err := getID(ctx.Params)
	if err != nil {
		util.MessageResponse(ctx, http.StatusBadRequest, "error in params: %v", err)
		return
	}

	m := dao.PKCond(id)
	code, err := getRecord[T](dao.DB, &m)
	if err != nil {
		util.MessageResponse(ctx, code, err.Error())
		return
	}

	_, isForce := ctx.GetQuery("force")
	if err = dao.DB.Model(new(T)).InstanceSet("force", isForce).Delete(&m).Error; err != nil {
		util.MessageResponse(ctx, http.StatusInternalServerError, "database error: %v", err)
		return
	}

	util.MessageResponse(ctx, http.StatusOK, "Deleted record successfully")
}
