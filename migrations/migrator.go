package migrations

import (
	"database/sql"
	"fmt"
	_ "github.com/go-sql-driver/mysql"
	"github.com/jmoiron/sqlx"
	"github.com/wen-bing/go-enterprise-web-toolkit/core/db"
	"io/ioutil"
	"log"
	"os"
	"path"
	"sort"
	"strconv"
	"strings"
	"time"
)

const (
	migrationLogTableName = "migration_logs"
)

type Migrator struct {
	db                   *sqlx.DB
	migratableDB         Migratable
	modules              []string
	schemaUpMigrations   map[string]db.MigrationSlice
	schemaDownMigrations map[string]db.MigrationSlice
	seederUpMigrations   map[string]db.MigrationSlice
	seederDownMigrations map[string]db.MigrationSlice
	schemaMigrationLogs  map[string]*db.MigrationLog
	seederMigrationLogs  map[string]*db.MigrationLog
}

func New(dir string, dbConfig db.DBConfig) *Migrator {
	m := Migrator{}
	dbUrl := fmt.Sprintf("%s:%s@tcp(%s:%d)/%s?charset=utf8&parseTime=True&loc=Local",
		dbConfig.User, dbConfig.Password, dbConfig.Host, dbConfig.Port, dbConfig.Name)

	var err error
	m.db, err = sqlx.Open(dbConfig.Type, dbUrl)
	if err != nil {
		panic(err)
	}
	m.migratableDB = MysqlMigratable{}

	m.schemaUpMigrations = make(map[string]db.MigrationSlice)
	m.schemaDownMigrations = make(map[string]db.MigrationSlice)
	m.seederUpMigrations = make(map[string]db.MigrationSlice)
	m.seederDownMigrations = make(map[string]db.MigrationSlice)
	m.schemaMigrationLogs = make(map[string]*db.MigrationLog)
	m.seederMigrationLogs = make(map[string]*db.MigrationLog)

	//collect all migration scripts
	m.collectModulesAndMigrations(dir)

	//find db migration status
	m.collectModuleLatestVersionInfo()
	return &m
}

func (m *Migrator) Done() {
	m.db.Close()
}

func (m *Migrator) Migrate() {
	m.migrateSchema()
	log.Println("Migration completed")
	m.insertSeederData()
	log.Println("Seeder completed")
}

func (m *Migrator) Rollback(timeToRollback time.Time) {
	m.deleteSeederData(timeToRollback)
	log.Println("Delete data completed")
	m.rollbackSchema(timeToRollback)
	log.Println("Rollback schema completed")
}

func (m *Migrator) migrateSchema() {
	//migrate schema
	for _, module := range m.modules {
		//collection migrations will be executed
		lastVersion := m.schemaMigrationLogs[module].Version
		var newVersions []*db.Migration
		for _, migration := range m.schemaUpMigrations[module] {
			if migration.Version > lastVersion {
				newVersions = append(newVersions, migration)
			}
		}
		if len(newVersions) == 0 {
			log.Printf("module %s schema is the latest version.\r\n", module)
			continue
		}
		m.executeMigrationScript(newVersions, module, db.OperationUp)
	}
}

func (m *Migrator) insertSeederData() {
	//insert seed data
	for _, module := range m.modules {
		lastVersion := m.seederMigrationLogs[module].Version
		var newVersions []*db.Migration
		for _, seeder := range m.seederUpMigrations[module] {
			if seeder.Version > lastVersion {
				newVersions = append(newVersions, seeder)
			}
		}
		if len(newVersions) == 0 {
			log.Printf("module %s seeder is the latest version.\r\n", module)
			continue
		}
		m.executeMigrationScript(newVersions, module, db.OperationUp)
	}
}

func (m *Migrator) deleteSeederData(timeToRollback time.Time) {
	//delete seed data
	for _, module := range m.modules {
		//migration type =1, means seeder
		migrationLogs := m.findMigrationLogsByDate(module, timeToRollback, db.TypeSeeder)
		if len(migrationLogs) == 0 {
			log.Printf("module %s seeeder logs is empty nothing to do.", module)
		} else {
			var revertVersions db.MigrationSlice
			for _, seeder := range m.seederDownMigrations[module] {
				for _, mlog := range migrationLogs {
					if mlog.Version == seeder.Version && mlog.Module == seeder.Module {
						revertVersions = append(revertVersions, seeder)
					}
				}
			}
			sort.Sort(sort.Reverse(revertVersions))

			m.executeMigrationScript(revertVersions, module, db.OperationDown)
		}
	}
}

func (m *Migrator) rollbackSchema(timeToRollback time.Time) {
	//rollback schema changes
	for _, module := range m.modules {
		//migration type =0, means schema
		migrationLogs := m.findMigrationLogsByDate(module, timeToRollback, db.TypeSchema)

		if len(migrationLogs) == 0 {
			log.Printf("module %s migration logs is empty nothing to do.", module)
		} else {
			var revertVersions db.MigrationSlice
			for _, migration := range m.schemaDownMigrations[module] {
				for _, mlog := range migrationLogs {
					if mlog.Version == migration.Version && mlog.Module == migration.Module {
						revertVersions = append(revertVersions, migration)
					}
				}
			}
			sort.Sort(sort.Reverse(revertVersions))
			m.executeMigrationScript(revertVersions, module, db.OperationDown)
		}
	}
}

func (m *Migrator) executeMigrationScript(newVersions []*db.Migration, module string, operation int) {
	tx := m.db.MustBegin()
	for _, version := range newVersions {
		sql, err := ioutil.ReadFile(version.Path)
		if err != nil {
			log.Printf("Read %s migration script failed: %s, %v", module, version.Path, err)
			panic(err)
		}
		commands := m.migratableDB.GetMigrationCommands(string(sql))
		//do the migration for the module
		for _, cmd := range commands {
			_, err := tx.Exec(cmd)
			if err != nil {
				tx.Rollback()
			}
		}

		if operation == db.OperationUp {
			log.Printf("Applied migration: %s", version.FileName)
			//insert migration log for the module
			r, err := tx.Exec(m.migratableDB.InsertMigrationLogSql(), module, version.Version, version.Type, version.Operation, time.Now())
			if err != nil {
				tx.Rollback()
			}
			rowAffect, _ := r.RowsAffected()
			log.Printf("Insert migration log, Rows affected: %d", rowAffect)
		}
		if operation == db.OperationDown {
			log.Printf("Revert migration: %s", version.FileName)
			r, err := tx.Exec(m.migratableDB.DeleteMigrationLogSql(), module, version.Version, version.Type)
			if err != nil {
				tx.Rollback()
			}
			rowAffect, _ := r.RowsAffected()
			log.Printf("Delete migration log, Rows affected: %d", rowAffect)
		}
	}
	err := tx.Commit()
	if err != nil {
		log.Printf("Commit %s migration failed: %v", module, err)
		tx.Rollback()
		panic(err)
	}
}

/**
Find each modules' migration version info
*/
func (m *Migrator) collectModuleLatestVersionInfo() {
	// Create the migrations table if it doesn't exist.
	tableExists, err := m.migrationTableExists()
	if err != nil {
		panic(err)
	}
	if !tableExists {
		if err := m.createMigrationsTable(); err != nil {
			panic(err)
		}
	}
	for _, moduleName := range m.modules {
		mlog := m.findLatestMigrationLog(moduleName, db.TypeSchema)
		m.schemaMigrationLogs[moduleName] = &mlog

		slog := m.findLatestMigrationLog(moduleName, db.TypeSeeder)
		m.seederMigrationLogs[moduleName] = &slog
	}
}

func (m *Migrator) findLatestMigrationLog(module string, migrationType int) db.MigrationLog {
	var mlog db.MigrationLog
	sqlStr := m.migratableDB.GetLatestMigrationLogSql()

	//find schema status
	rows, err := m.db.Queryx(sqlStr, module, db.TypeSchema)
	if err != nil {
		log.Printf("Get migration status failed: for module: %s: %v", module, err)
		panic(err)
	}
	for rows.Next() {
		err = rows.StructScan(&mlog)
		if err != nil {
			log.Printf("Get migration status failed: for module: %s: %v", module, err)
			panic(err)
		}
	}
	return mlog
}

func (m *Migrator) findMigrationLogsByDate(module string, timeToRollback time.Time, migrationType int) db.MigrationLogSlice {
	sqlStr := m.migratableDB.GetMigrationLogsByDateSql()
	//find seeder logs after the timeToRollback, those logs are need be rollbacked
	rows, err := m.db.Queryx(sqlStr, module, migrationType, timeToRollback)
	if err != nil {
		log.Printf("Rollback failed to find module: %s migration logs: %v", module, err)
		panic(err)
	}
	var migrationLogs db.MigrationLogSlice
	for rows.Next() {
		var l db.MigrationLog
		err = rows.StructScan(&l)
		if err != nil {
			log.Printf("Rollback failed to parse migration log: %v", err)
			panic(err)
		}
		migrationLogs = append(migrationLogs, &l)
	}
	return migrationLogs
}

// Returns true if the migration table already exists.
func (m *Migrator) migrationTableExists() (bool, error) {
	row := m.db.QueryRow(m.migratableDB.SelectMigrationTableSql(), migrationLogTableName)
	var tableName string
	err := row.Scan(&tableName)
	if err == sql.ErrNoRows {
		log.Println("Migrations table not found")
		return false, nil
	}
	if err != nil {
		log.Printf("Error checking for migration table: %v", err)
		return false, err
	}
	log.Println("Migrations table found")
	return true, nil
}

// Creates the migrations table if it doesn't exist.
func (m *Migrator) createMigrationsTable() error {
	_, err := m.db.Exec(m.migratableDB.CreateMigrationTableSql())
	if err != nil {
		log.Fatalf("Error creating migrations table: %v", err)
	}
	log.Printf("Created migrations table: %s", migrationLogTableName)
	return nil
}

func (m *Migrator) collectModulesAndMigrations(dir string) {
	if _, err := os.Stat(dir); os.IsNotExist(err) {
		panic(err)
	}
	subFolders, err := ioutil.ReadDir(dir)
	if err != nil {
		panic(err)
	}

	for _, folder := range subFolders {
		if folder.IsDir() {
			module := folder.Name()
			m.modules = append(m.modules, module)

			//scriptType=0 for schema
			moduleMigrations := path.Join(dir, module, "migrations")
			m.collectMigrationScripts(module, moduleMigrations, db.TypeSchema)
			sort.Sort(m.schemaUpMigrations[module])
			sort.Sort(m.schemaDownMigrations[module])
			//scriptType=1 for seeder
			moduleSeeders := path.Join(dir, module, "seeders")
			sort.Sort(m.seederUpMigrations[module])
			sort.Sort(m.seederDownMigrations[module])
			m.collectMigrationScripts(module, moduleSeeders, db.TypeSeeder)
		}
	}
}

func (m *Migrator) collectMigrationScripts(module string, moduleMigrationPath string, scriptType int) {
	if _, err := os.Stat(moduleMigrationPath); !os.IsNotExist(err) {
		scrips, err := ioutil.ReadDir(moduleMigrationPath)
		if err != nil {
			panic(err)
		}
		for _, f := range scrips {
			if !f.IsDir() {
				name := f.Name()
				//1-0-name, 1 for version, 0 for operation type
				token := strings.Split(name, "-")
				version, err := strconv.Atoi(token[0])
				if err != nil {
					panic("migration file name format error, should have version on first section")
				}
				op, err := strconv.Atoi(token[1])
				if err != nil {
					panic("migration file name format error, should have operation on second section")
				}
				scriptPath := path.Join(moduleMigrationPath, name)
				migration := db.Migration{module, scriptType, op, version, name, scriptPath}
				//for migration
				if scriptType == db.TypeSchema {
					if op == db.OperationUp {
						//up migration
						m.schemaUpMigrations[module] = append(m.schemaUpMigrations[module], &migration)
					} else if op == db.OperationDown {
						//down migration
						m.schemaDownMigrations[module] = append(m.schemaDownMigrations[module], &migration)
					}
				} else if scriptType == db.TypeSeeder { //for seeder
					if op == db.OperationUp {
						//up migration
						m.seederUpMigrations[module] = append(m.seederUpMigrations[module], &migration)
					} else if op == db.OperationDown {
						//down migration
						m.seederDownMigrations[module] = append(m.seederDownMigrations[module], &migration)
					}
				}
			}
		}
	}
}
