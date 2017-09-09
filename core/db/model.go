package db

import (
	"github.com/satori/go.uuid"
	"strings"
	"time"
)

type Model struct {
	Id        string     `db:"size=32" json:"id"`
	CreatedAt time.Time  `db:"created_at,null" json:"created_at"`
	UpdatedAt time.Time  `db:"updated_at,null" json:"updated_at"`
	DeletedAt *time.Time `db:"deleted_at,null" json:"deleted_at"`
}

func GenerateModelId() string {
	u := uuid.NewV4()
	str := u.String()
	id := strings.Replace(str, "-", "", -1)
	return id
}
