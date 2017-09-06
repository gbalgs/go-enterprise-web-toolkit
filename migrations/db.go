package migrations

import "strings"

type Migratable interface {
	SelectMigrationTableSql() string
	CreateMigrationTableSql() string
	GetLatestMigrationLogSql() string
	GetMigrationLogsByDateSql() string
	InsertMigrationLogSql() string
	DeleteMigrationLogSql() string
	GetMigrationCommands(string) []string
}

type MysqlMigratable struct{}

func (m MysqlMigratable) SelectMigrationTableSql() string {
	return "SELECT table_name FROM information_schema.tables WHERE table_name = ? AND table_schema = (SELECT DATABASE())"
}

func (m MysqlMigratable) CreateMigrationTableSql() string {
	return `CREATE TABLE migration_logs(id INT NOT NULL AUTO_INCREMENT PRIMARY KEY,
	module VARCHAR(255) NOT NULL, version INT, type INT, operation INT, apply_date TIMESTAMP NOT NULL)`
}

func (m MysqlMigratable) GetLatestMigrationLogSql() string {
	return `SELECT module, version, type, operation, apply_date FROM migration_logs WHERE module=? AND type=? AND apply_date=(SELECT max(apply_date))`
}

func (m MysqlMigratable) GetMigrationLogsByDateSql() string {
	return `SELECT module, version, type, operation, apply_date FROM migration_logs WHERE module=? AND type=? AND apply_date>?`
}

func (m MysqlMigratable) InsertMigrationLogSql() string {
	return "INSERT INTO migration_logs(module, version, type, operation, apply_date) values (?,?,?,?,?)"
}

func (m MysqlMigratable) DeleteMigrationLogSql() string {
	return "DELETE FROM migration_logs WHERE module=? AND version=? AND type=?"
}

func (m MysqlMigratable) GetMigrationCommands(sql string) []string {
	count := strings.Count(sql, ";")
	commands := strings.SplitN(string(sql), ";", count)
	return commands
}
