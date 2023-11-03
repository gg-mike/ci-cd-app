package model

import "gorm.io/gorm"

func GetFromRaw[T any](tx *gorm.DB, field string) (T, bool) {
	raw, ok := tx.InstanceGet("raw")
	if !ok {
		return *new(T), false
	}
	val, ok := raw.(map[string]any)[field]
	if !ok {
		return *new(T), false
	}
	return val.(T), ok
}

func IsForce(tx *gorm.DB) bool {
	force, ok := tx.InstanceGet("force")
	return ok && force.(bool)
}
