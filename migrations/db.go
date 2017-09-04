package migrations

import "strings"

type Migratable interface {
	SelectMigrationTableSql() string
	CreateMigrationTableSql() string
	GetModuleMigrationStatusSql() string
	MigrationLogInsertSql() string
	MigrationLogDeleteSql() string
	GetMigrationCommands(string) []string
}

type MysqlMigratable struct{}

func (m MysqlMigratable) SelectMigrationTableSql() string {
	return "SELECT table_name FROM information_schema.tables WHERE table_name = ? AND table_schema = (SELECT DATABASE())"
}

func (m MysqlMigratable) CreateMigrationTableSql() string {
	return `CREATE TABLE gewt_migration_log(id INT NOT NULL AUTO_INCREMENT PRIMARY KEY,
	module VARCHAR(255) NOT NULL, version INT, type INT, operation INT, apply_date TIMESTAMP NOT NULL)`
}

func (m MysqlMigratable) GetModuleMigrationStatusSql() string {
	return `SELECT id, module, version, type, operation, apply_date
		FROM gewt_migration_log WHERE module = ? AND apply_date=(SELECT max(date)`
}

func (m MysqlMigratable) MigrationLogInsertSql() string {
	return "INSERT INTO gewt_migration_log(module, version, type, operation, apply_date) values (?,?,?,?,?)"
}

func (m MysqlMigratable) MigrationLogDeleteSql() string {
	return "DELETE FROM gewt_migration_log WHERE migration_id = ?"
}

func (m MysqlMigratable) GetMigrationCommands(sql string) []string {
	count := strings.Count(sql, ";")
	commands := strings.SplitN(string(sql), ";", count)
	return commands
}
