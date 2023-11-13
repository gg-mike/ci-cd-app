package dao

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"

	"github.com/gg-mike/ci-cd-app/backend/internal/db"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

type DAO[T, TShort any] struct {
	Preload func(tx *gorm.DB) *gorm.DB
	Filter  func(ctx *gin.Context) (map[string]any, error)
}

type IDAO interface {
	GetOne(ctx *gin.Context)
	GetMany(ctx *gin.Context)
	Create(ctx *gin.Context)
	Update(ctx *gin.Context)
	Delete(ctx *gin.Context)
}

func (dao DAO[T, TShort]) GetOne(ctx *gin.Context) {
	id, err := getID(ctx.Params)
	if err != nil {
		dao.MessageResponse(ctx, http.StatusBadRequest, "error in params: %v", err)
		return
	}

	m, code, err := getRecordFromID[T](dao, id)
	if err != nil {
		dao.MessageResponse(ctx, code, err.Error())
	} else {
		ctx.JSON(code, m)
	}
}

func (dao DAO[T, TShort]) GetMany(ctx *gin.Context) {
	filters, err := dao.Filter(ctx)
	if err != nil {
		dao.MessageResponse(ctx, http.StatusBadRequest, "error parsing query: %v", err)
	}

	offset, limit, order, err := Paginate(ctx)
	if err != nil {
		dao.MessageResponse(ctx, http.StatusBadRequest, "error in query: %v", err)
	}

	o := new([]TShort)
	_db := db.Get().Model(new(T))
	for key, value := range filters {
		_db = _db.Where(key, value)
	}
	err = _db.Offset(offset).Limit(limit).Order(order).Find(o).Error
	switch err {
	case nil:
		ctx.JSON(http.StatusOK, o)
	case gorm.ErrRecordNotFound:
		dao.MessageResponse(ctx, http.StatusNotFound, "no record found")
	default:
		dao.MessageResponse(ctx, http.StatusInternalServerError, "database error: %v", err)
	}
}

func (dao DAO[T, TShort]) Create(ctx *gin.Context) {
	raw := map[string]any{}
	m := *new(T)
	if err := getBody[T](ctx, &raw, &m); err != nil {
		dao.MessageResponse(ctx, http.StatusBadRequest, "error parsing body: %v", err)
		return
	}

	err := db.Get().InstanceSet("raw", raw).Create(&m).Error
	switch err {
	case nil:
		var mMap map[string]interface{}
		jsonBytes, err := json.Marshal(m)
		if err != nil {
			dao.MessageResponse(ctx, http.StatusInternalServerError, err.Error())
			return
		}
		if err := json.Unmarshal(jsonBytes, &mMap); err != nil {
			dao.MessageResponse(ctx, http.StatusInternalServerError, err.Error())
			return
		}
		id, ok := mMap["id"]
		if !ok {
			dao.MessageResponse(ctx, http.StatusInternalServerError, "missing 'id' in map")
			return
		}

		created, code, err := getRecordFromID[T](dao, uuid.MustParse(id.(string)))
		if err != nil {
			dao.MessageResponse(ctx, code, err.Error())
			return
		}
		ctx.JSON(http.StatusCreated, created)
	case gorm.ErrForeignKeyViolated:
		dao.MessageResponse(ctx, http.StatusConflict, "incorrect foreign key")
	default:
		dao.MessageResponse(ctx, http.StatusInternalServerError, "database error: %v", err)
	}
}

func (dao DAO[T, TShort]) Update(ctx *gin.Context) {
	id, err := getID(ctx.Params)
	if err != nil {
		dao.MessageResponse(ctx, http.StatusBadRequest, "error in params: %v", err)
		return
	}

	prev, code, err := getRecordFromID[T](dao, id)
	if err != nil {
		dao.MessageResponse(ctx, code, err.Error())
	}
	m := prev
	raw := map[string]any{}
	if err := getBody[T](ctx, &raw, &m); err != nil {
		dao.MessageResponse(ctx, http.StatusBadRequest, "error parsing body: %v", err)
		return
	}

	err = db.Get().Model(&prev).InstanceSet("raw", raw).InstanceSet("prev", prev).Updates(m).Error
	switch err {
	case nil:
		updated, code, err := getRecordFromID[T](dao, id)
		if err != nil {
			dao.MessageResponse(ctx, code, err.Error())
		}
		ctx.JSON(http.StatusOK, updated)
	case gorm.ErrForeignKeyViolated:
		dao.MessageResponse(ctx, http.StatusConflict, "incorrect foreign key")
	default:
		dao.MessageResponse(ctx, http.StatusInternalServerError, "database error: %v", err)
	}
}

func (dao DAO[T, TShort]) Delete(ctx *gin.Context) {
	id, err := getID(ctx.Params)
	if err != nil {
		dao.MessageResponse(ctx, http.StatusBadRequest, "error in params: %v", err)
		return
	}

	m, code, err := getRecordFromID[T](dao, id)
	if err != nil {
		dao.MessageResponse(ctx, code, err.Error())
		return
	}

	_, isForce := ctx.GetQuery("force")
	if err = db.Get().Model(new(T)).InstanceSet("force", isForce).Delete(&m).Error; err != nil {
		dao.MessageResponse(ctx, http.StatusInternalServerError, "database error: %v", err)
		return
	}

	dao.MessageResponse(ctx, http.StatusOK, "Deleted record successfully")
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

func getRecordFromID[T, TShort any](dao DAO[T, TShort], id uuid.UUID) (T, int, error) {
	var err error
	m := new(T)
	if dao.Preload == nil {
		err = db.Get().Preload(clause.Associations).First(m, "id = ?", id).Error
	} else {
		err = dao.Preload(db.Get()).First(m, "id = ?", id).Error
	}
	switch err {
	case nil:
		return *m, http.StatusOK, nil
	case gorm.ErrRecordNotFound:
		return *m, http.StatusNotFound, errors.New("no record found")
	default:
		return *m, http.StatusInternalServerError, fmt.Errorf("database error: %v", err)
	}
}
