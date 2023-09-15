package main

import (
	"context"
	"flag"
	"fmt"
	"log"
	"os"
	"time"

	"fyne.io/systray"
	"go.uber.org/zap"

	"nv-pr-app/pkg/dbmigrate"
	github "nv-pr-app/pkg/github/service"
	"nv-pr-app/pkg/persistence"

	"nv-pr-app/internal/config"
	"nv-pr-app/internal/pull_request/repository"
	"nv-pr-app/internal/pull_request/service"
	"nv-pr-app/internal/tray"
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

	logger, _ := zap.NewProduction()

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

	// err = pullRequestService.NotifyPullRequests(ctx)
	// if err != nil {
	// 	panic(err)
	// }

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

	//Run icon tray
	app, err := tray.New(sc.Projects)
	if err != nil {
		panic(err)
	}
	go jobToCreateMenus(context.Background(), logger, app, sc.Projects, pullRequestService)

	systray.Run(app.OnReady, app.OnExit)
}

func jobToCreateMenus(ctx context.Context, logger *zap.Logger, app *tray.App, projects []config.Project, pullResquestService *service.PullRequestService) {

	for {
		logger.Info("Starting new Github request")

		logger.Info("Reseting menu")
		app.Reset()
		for _, project := range projects {
			getPullRequestsByProject(ctx, logger, app, project, pullResquestService)
		}
		logger.Info("Waiting new loop")

		time.Sleep(60 * time.Second)
	}
}

func getPullRequestsByProject(ctx context.Context, logger *zap.Logger, app *tray.App, project config.Project, pullResquestService *service.PullRequestService) {
	for _, repo := range project.Repos {
		go getPullRequestsByOwner(ctx, logger, app, project.Owner, repo, pullResquestService)
	}
}

func getPullRequestsByOwner(ctx context.Context, logger *zap.Logger, app *tray.App, owner, repo string, pullResquestService *service.PullRequestService) {
	logger.Info(fmt.Sprintf("Get %s:%s", owner, repo))
	prs, err := pullResquestService.ListPullRequests(ctx, owner, repo)
	if err != nil {
		logger.Error(err.Error())
	}

	if len(prs) > 0 {
		logger.Info(fmt.Sprintf("%d PRs to %s", len(prs), repo))
		for _, pr := range prs {
			repo := pr.Head.Repo.Name
			app.Show(repo)
			app.AddMenuItem(repo, pr.Title, pr.HTMLURL)
		}
	}

}

// func onReady() {

// 	mBrowser := systray.AddMenuItem("Open", "Open Browser")
// 	go func() {
// 		<-mBrowser.ClickedCh
// 		fmt.Println("Oepn")
// 		openBrowser("http://www.google.com")
// 		fmt.Println("Close")
// 	}()
// 	mQuit := systray.AddMenuItem("Quit", "Quit the whole app")
// 	// Sets the icon of a menu item. Only available on Mac and Windows.
// 	mQuit.SetIcon(wellzIcon)
// 	go func() {
// 		<-mQuit.ClickedCh
// 		fmt.Println("Requesting quit")
// 		systray.Quit()
// 		fmt.Println("Finished quitting")
// 	}()
// }

// func onExit() {
// 	// clean up here
// }
