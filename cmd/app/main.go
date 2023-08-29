package main

import (
	"context"
	"flag"
	"log"
	"os"

	"nv-pr-app/pkg/dbmigrate"
	github "nv-pr-app/pkg/github/service"
	"nv-pr-app/pkg/persistence"

	"nv-pr-app/internal/config"
	"nv-pr-app/internal/pull_request/repository"
	"nv-pr-app/internal/pull_request/service"
)

func main() {
	var configFile string

	// -c options set configuration file path, but can be overwritten by CONFIG_FILE environment variable
	flag.StringVar(&configFile, "c", "configs/envs.json", "config file path")
	flag.Parse()

	// If you specify an option by using environment variables, it overrides any value loaded from the configuration file
	path := os.Getenv("CONFIG_FILE")
	if path != "" {
		configFile = path
	}

	// Load configuration yaml file using -c location/CONFIG_FILE and merging environments variables with higher precedence
	sc, err := config.LoadServiceConfig(configFile)
	if err != nil {
		log.Fatalf("main: could not load service configuration [%v]", err)
	}

	db, err := persistence.New(sc.Database.DbPath)
	if err != nil {
		panic(err)
	}

	m, err := dbmigrate.New(db.DB, sc.Database.MigrationsPath)
	if err != nil {
		panic(err)
	}
	err = m.Up()
	if err != nil {
		panic(err)
	}

	githubServide := github.New(&sc.Github)

	pullRequestService := service.New(repository.New(db), sc.Projects, githubServide)

	ctx := context.Background()
	err = pullRequestService.AddNonexistentPullRequests(ctx)
	if err != nil {
		panic(err)
	}

	err = pullRequestService.NotifyPullRequests(ctx)
	if err != nil {
		panic(err)
	}

	// fmt.Println("teste")
	// err = beeep.Beep(15, 2)
	// if err != nil {
	// 	panic(err)
	// }

	// err = beeep.Notify("Title", "Message body", "ic_launcher_foreground.png")
	// if err != nil {
	// 	panic(err)
	// }

	// err = beeep.Alert("Title", "Message body", "assets/warning.png")
	// if err != nil {
	// 	panic(err)
	// }
}
