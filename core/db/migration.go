package db

type Migration struct {
	Module    string `json: "module" db:"module"`
	FileName  string `json:"file_name" db:"file_name"`
	Path      string `json:"path" db:"ignore"`
	Type      int    `json:"type" db:"type"`           //0 for migration, 1 for seeder
	Operation int    `json:"operation" db:"operation"` //1 for up, 0 for down
	Version   int    `json:"version" db:"version"`
}

type MigrationSlice []*Migration

type MigrationLog struct {
	Id uint64 `json:"id" db:"id"`
	Migration
}

func (ms MigrationSlice) Len() int {
	return len(ms)
}

func (ms MigrationSlice) Swap(i, j int) {
	ms[i], ms[j] = ms[j], ms[i]
}

func (ms MigrationSlice) Less(i, j int) bool {
	return ms[i].Version < ms[j].Version
}
