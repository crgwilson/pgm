package main

import (
	"flag"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"

	"github.com/crgwilson/pgm/pkg/logger"
	"github.com/crgwilson/pgm/pkg/migrate"
	"github.com/crgwilson/pgm/pkg/pg"
)

const usageText = `pgm: PostgreSQL schema migrator

Usage:
    pgm [flags] <command>

Commands:
    init                   Perform all first time setup necessary for this tool to work
    up                     Run all available sql scripts until the highest available version is reached
    down                   Run all available sql scripts to completely revert all schema changes back to the first version
    version                Print the current schema version

`

func usage() {
	fmt.Print(usageText)
	flag.PrintDefaults()
	os.Exit(1)
}

func main() {
	// Init CLI flags
	verbose := flag.Bool("v", false, "Log more verbosely")
	sqlDir := flag.String("d", "./", "The directory containing SQL migration scripts")
	dbHost := flag.String("H", "localhost", "Host address of the PostgreSQL database")
	dbPort := flag.Int("p", 5432, "Host port of the PostgreSQL database")
	dbUser := flag.String("u", "postgres", "Login user for the PostgreSQL database")
	dbPassword := flag.String("P", "", "Login password for the PostgreSQL database")
	dbName := flag.String("D", "postgres", "The name of the database to connect to")
	dbSslMode := flag.String("s", "verify-full", "The 'sslmode' to set in the PostgreSQL connection URI")

	flag.Usage = usage
	flag.Parse()

	// Init logger
	var logLevel logger.LogLevel
	if *verbose {
		logLevel = logger.DebugLogLevel()
	} else {
		logLevel = logger.InfoLogLevel()
	}
	cliLogger := logger.NewCliLogger(logLevel)

	// Configure postgres connection
	pgConfig := pg.PostgresConfig{
		Address:  *dbHost,
		Port:     *dbPort,
		User:     *dbUser,
		Password: *dbPassword,
		Database: *dbName,
		SslMode:  *dbSslMode,
	}

	db, err := pg.OpenDb(pgConfig)
	if err != nil {
		errorLog := fmt.Sprintf("%v", err)
		cliLogger.Error(errorLog)
		os.Exit(2)
	}

	migrationStore := migrate.NewSchemaMigrationStore(db)
	migrator := migrate.NewMigrationManager(migrationStore, cliLogger)

	// Register all provided sql files
	files, err := ioutil.ReadDir(*sqlDir)
	if err != nil {
		errorLog := fmt.Sprintf("%v", err)
		cliLogger.Error(errorLog)
		os.Exit(3)
	}

	for _, file := range files {
		// If we find a non-sql file, we ignore it
		if filepath.Ext(file.Name()) != ".sql" {
			continue
		}

		sqlFileName := file.Name()
		sqlFilePath := *sqlDir + "/" + sqlFileName
		sqlFileContent, err := ioutil.ReadFile(sqlFilePath)
		if err != nil {
			errorLog := fmt.Sprintf("%v", err)
			cliLogger.Error(errorLog)
			os.Exit(4)
		}

		parsedSqlFile, err := migrate.ParseSqlFile(sqlFileName, sqlFileContent)
		if err != nil {
			errorLog := fmt.Sprintf("%v", err)
			cliLogger.Error(errorLog)
			os.Exit(5)
		}
		migrator.RegisterMigrationPath(parsedSqlFile)
	}

	// After all the flags we expect to find a subcommand of some sort
	switch os.Args[len(os.Args)-1] {
	case "init":
		err = migrator.InitDb()
		if err != nil {
			cliLogger.Error(fmt.Sprintf("%v", err))
			os.Exit(6)
		}
	case "up":
		// Upgrade DB schema using the `up.sql` files we know about
		highest := migrator.HighestAvailableVersion()
		err := migrator.Up(highest)
		if err != nil {
			cliLogger.Error(fmt.Sprintf("%v", err))
			os.Exit(7)
		}
	case "down":
		// Downgrade DB schema using the `down.sql` files we know about
		lowest := migrator.LowestAvailableVersion()
		err := migrator.Down(lowest)
		if err != nil {
			cliLogger.Error(fmt.Sprintf("%v", err))
			os.Exit(8)
		}
	case "version":
		// Get the current version of DB schema we have deployed
		version, err := migrator.CurrentVersion()
		if err != nil {
			cliLogger.Error(fmt.Sprintf("%v", err))
			os.Exit(9)
		}

		cliLogger.Info(version)
	default:
		// If we don't find a subcommand of some sort just print out the help info
		usage()
	}
}
