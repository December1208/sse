package webutil

import (
	"database/sql/driver"
	"encoding/json"
	"time"

	"github.com/google/uuid"
	"github.com/jinzhu/gorm"
)

type JSONB map[string]interface{}

func (j JSONB) Value() (driver.Value, error) {
	valueString, err := json.Marshal(j)
	return string(valueString), err
}

func (j *JSONB) Scan(value interface{}) error {
	if err := json.Unmarshal(value.([]byte), &j); err != nil {
		return err
	}
	return nil
}

type BaseModel struct {
	UUID      uuid.UUID `gorm:"type:uuid;primary_key;" json:"uuid"`
	CreatedAt int64     `sql:"index" json:"created_at"`
	UpdatedAt int64     `sql:"index" json:"updated_at"`
	IsDeleted bool      `sql:"index" json:"is_deleted"`
}

func (base *BaseModel) BeforeUpdate(scope *gorm.Scope) error {
	return scope.SetColumn("UpdatedAt", time.Now().Unix())
}

func (base *BaseModel) BeforeCreate(scope *gorm.Scope) error {
	modelUUID, err := uuid.NewUUID()
	if err != nil {
		return err
	}
	if err := scope.SetColumn("CreatedAt", time.Now().Unix()); err != nil {
		return err
	}
	if err := scope.SetColumn("UpdatedAt", time.Now().Unix()); err != nil {
		return err
	}
	if err := scope.SetColumn("IsDeleted", false); err != nil {
		return err
	}
	return scope.SetColumn("UUID", modelUUID)
}

func (base *BaseModel) GetInfo(result map[string]interface{}) map[string]interface{} {
	result["uuid"] = base.UUID
	result["created_at"] = base.CreatedAt
	return result
}
