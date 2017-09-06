package db

import "time"

const (
	TypeSchema    = 0
	TypeSeeder    = 1
	OperationUp   = 0
	OperationDown = 1
)

type Migration struct {
	Module    string `db:"module"`
	Type      int    `db:"type"`      //0 for schema, 1 for seeder
	Operation int    `db:"operation"` //0 for up, 1 for down
	Version   int    `db:"version"`
	FileName  string `db:"-" sql:"-"`
	Path      string `db:"-" sql:"-"`
}

type MigrationSlice []*Migration

func (ms MigrationSlice) Len() int {
	return len(ms)
}

func (ms MigrationSlice) Swap(i, j int) {
	ms[i], ms[j] = ms[j], ms[i]
}

func (ms MigrationSlice) Less(i, j int) bool {
	return ms[i].Version < ms[j].Version
}

type MigrationLog struct {
	Migration
	Id        uint64    `json:"id" db:"id"`
	ApplyDate time.Time `db:"apply_date"`
}

type MigrationLogSlice []*MigrationLog

func (ms MigrationLogSlice) Len() int {
	return len(ms)
}

func (ms MigrationLogSlice) Swap(i, j int) {
	ms[i], ms[j] = ms[j], ms[i]
}

func (ms MigrationLogSlice) Less(i, j int) bool {
	return ms[i].Version < ms[j].Version
}
