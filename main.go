package main

import (
	"flag"
	"fmt"
	"github.com/wen-bing/go-enterprise-web-toolkit/migrations"
	"github.com/wen-bing/go-enterprise-web-toolkit/server"
	"os"
)

func main() {

	serverCommand := flag.NewFlagSet("server", flag.ExitOnError)
	serverEnvFlag := serverCommand.String("e", "development", "server running environment")
	serverConfig := serverCommand.String("c", ".", "")

	migrationCommand := flag.NewFlagSet("migration", flag.ExitOnError)
	migrationsEnvFlag := migrationCommand.String("e", "development", "server running environment")
	migrationsConfigFlag := migrationCommand.String("c", ".", "")
	migrationDirFlag := migrationCommand.String("d", "./migrations", "migration scripts path")

	if len(os.Args) == 1 {
		fmt.Println("Usage: gewt [command] [flags]")
		fmt.Println("sub commands are: ")
		fmt.Println("server		start the server application")
		fmt.Println("server -e {development/test/production} 	server runing environment")
		fmt.Println("server -c dir	server configuration files dir")
		fmt.Println("migration 	start migration tool to migrate database schema")
		fmt.Println("migration -d {migration scripts path}")
		fmt.Println("migration -e {development/test/production} 	migration runing environment")
		fmt.Println("migration -c dir	migration configuration files dir")
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
		migrator.Migrate()
	}

}
