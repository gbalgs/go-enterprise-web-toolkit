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
	migrationLogTableName = "gewt_migration_log"
)

type Migrator struct {
	db             *sqlx.DB
	migratableDB   Migratable
	modules        []string
	upMigrations   map[string]db.MigrationSlice
	downMigrations map[string]db.MigrationSlice
	upSeeders      map[string]db.MigrationSlice
	downSeeders    map[string]db.MigrationSlice
	migrationLogs  map[string]*db.MigrationLog
	seedLogs       map[string]*db.MigrationLog
}

func New(dir string, dbConfig db.DBConfig) *Migrator {
	m := Migrator{}
	dbUrl := fmt.Sprintf("%s:%s@tcp(%s:%d)/%s", dbConfig.User, dbConfig.Password,
		dbConfig.Host, dbConfig.Port, dbConfig.Name)

	var err error
	m.db, err = sqlx.Open(dbConfig.Type, dbUrl)
	if err != nil {
		panic(err)
	}
	m.migratableDB = MysqlMigratable{}

	m.upMigrations = make(map[string]db.MigrationSlice)
	m.downMigrations = make(map[string]db.MigrationSlice)
	m.upSeeders = make(map[string]db.MigrationSlice)
	m.downSeeders = make(map[string]db.MigrationSlice)
	m.migrationLogs = make(map[string]*db.MigrationLog)
	m.seedLogs = make(map[string]*db.MigrationLog)

	//collect all migration scripts
	m.collectModulesAndMigrations(dir)

	//find db migration status
	m.findMigrationStatus()
	return &m
}

func (m *Migrator) Migrate() {
	//migrate shcema
	for _, module := range m.modules {
		lastVersion := m.migrationLogs[module].Version
		var newVersions []*db.Migration
		for _, migration := range m.upMigrations[module] {
			if migration.Version > lastVersion {
				newVersions = append(newVersions, migration)
			}
		}
		tx := m.db.MustBegin()
		var latestVersion *db.Migration
		for _, version := range newVersions {
			sql, err := ioutil.ReadFile(version.Path)
			if err != nil {
				log.Printf("Read %s migration script failed: %s, %v", module, version.Path, err)
				panic(err)
			}
			latestVersion = version
			commands := m.migratableDB.GetMigrationCommands(string(sql))
			//do the migration for the module
			for _, cmd := range commands {
				tx.MustExec(cmd)
			}
		}
		tx.MustExec(m.migratableDB.MigrationLogInsertSql(), module,
			latestVersion.Version, latestVersion.Type, latestVersion.Operation, time.Now())
		//insert migration log for the module
		err := tx.Commit()
		if err != nil {
			log.Printf("Apply %s migration failed: %v", module, err)
			panic(err)
		}
	}

	//insert seed data
}

func (m *Migrator) Rollback(timeToRollback time.Time) {
	//delete seed data
	//rollback schema changes
}

func (m *Migrator) findMigrationStatus() {
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
	for _, module := range m.modules {
		var mlog db.MigrationLog
		err = m.db.Get(&mlog, m.migratableDB.GetModuleMigrationStatusSql(), module)
		if err != nil {
			log.Printf("Get migration status failed: for module: %s: %v", module, err)
			panic(err)
		}
		m.migrationLogs[module] = &mlog
	}
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

			//scriptType=0 for migration
			moduleMigrations := path.Join(dir, module, "migrations")
			m.collectMigrationScipts(module, moduleMigrations, 0)
			sort.Sort(m.upMigrations[module])
			sort.Sort(m.downMigrations[module])
			//scriptType=1 for seeder
			moduleSeeders := path.Join(dir, module, "seeders")
			sort.Sort(m.upSeeders[module])
			sort.Sort(m.downSeeders[module])
			m.collectMigrationScipts(module, moduleSeeders, 1)
		}
	}
}

func (m *Migrator) collectMigrationScipts(module string, moduleMigrationPath string, scriptType int) {

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
				migration := db.Migration{module, name, scriptPath, scriptType, op, version}
				//for migration
				if scriptType == 0 {
					if op == 1 {
						//up migration
						m.upMigrations[module] = append(m.upMigrations[module], &migration)
					} else if op == 0 {
						//down migration
						m.downMigrations[module] = append(m.downMigrations[module], &migration)
					}
				} else if scriptType == 1 { //for seeder
					if op == 1 {
						//up migration
						m.upSeeders[module] = append(m.upSeeders[module], &migration)
					} else if op == 0 {
						//down migration
						m.downSeeders[module] = append(m.downSeeders[module], &migration)
					}
				}
			}
		}
	}
}
