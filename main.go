package main

import (
	"flag"
	"fmt"
	"github.com/wen-bing/go-enterprise-web-toolkit/migrations"
	"github.com/wen-bing/go-enterprise-web-toolkit/server"
	"log"
	"os"
	"time"
)

func main() {
	serverCommand := flag.NewFlagSet("server", flag.ExitOnError)
	serverEnvFlag := serverCommand.String("e", "development", "server running environment")
	serverConfig := serverCommand.String("c", ".", "")

	migrationCommand := flag.NewFlagSet("migration", flag.ExitOnError)
	migrationsEnvFlag := migrationCommand.String("e", "development", "server running environment")
	migrationsConfigFlag := migrationCommand.String("c", ".", "")
	migrationDirFlag := migrationCommand.String("d", "./migrations", "migration scripts path")
	rollbackFlag := migrationCommand.String("r", "", "rollback schema to the date")

	if len(os.Args) == 1 {
		printUsage()
		return
	}
	switch os.Args[1] {
	case "server":
		serverCommand.Parse(os.Args[2:])
	case "migration":
		migrationCommand.Parse(os.Args[2:])
	default:
		fmt.Printf("%q is not valid command\n", os.Args[1])
		os.Exit(2)
	}

	if serverCommand.Parsed() {
		appConfig := LoadConfig(*serverEnvFlag, *serverConfig)
		app := server.New(appConfig)
		app.Start()
		return
	}

	if migrationCommand.Parsed() {
		migrationConfig := LoadConfig(*migrationsEnvFlag, *migrationsConfigFlag)
		migrator := migrations.New(*migrationDirFlag, migrationConfig.DB)

		if *rollbackFlag == "" {
			migrator.Migrate()
		} else {
			t, err := time.Parse("2006-01-02", *rollbackFlag)
			if err != nil {
				log.Fatal("rollback parameter format error: %s, should be 2016-01-02", *rollbackFlag)
			}
			log.Printf("Try to rollback to: %v", t)
			migrator.Rollback(t)
		}

		migrator.Done()
	}

}
func printUsage() {
	fmt.Println("Usage: gewt [command] [flags]")
	fmt.Println("sub commands are: ")
	fmt.Println("server		start the server application")
	fmt.Println("server -e {development/test/production} 	server runing environment")
	fmt.Println("server -c dir	server configuration files dir")
	fmt.Println("migration 	start migration tool to migrate database schema")
	fmt.Println("migration -d {migration scripts path}")
	fmt.Println("migration -e {development/test/production} 	migration runing environment")
	fmt.Println("migration -c dir	migration configuration files dir")
	fmt.Println("migration -r date rollback schema and date to the specified date")
}
