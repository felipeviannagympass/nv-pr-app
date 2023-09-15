package tray

import (
	"fmt"
	"log"
	"os/exec"
	"runtime"

	"fyne.io/systray"

	"nv-pr-app/internal/assets"
	"nv-pr-app/internal/config"
)

const (
	TITLE   = "Wellz PRs"
	TOOLTIP = "See all PRs to wellz, be nice and review"
)

type App struct {
	title     string
	tooltip   string
	trayIcon  []byte
	projects  []config.Project
	menuItems map[string]*systray.MenuItem
}

func New(projects []config.Project) (*App, error) {
	wellzIcon, err := assets.New(assets.WELLZ_ICON).Get()
	if err != nil {
		return nil, err
	}

	return &App{
		title:     TITLE,
		tooltip:   TOOLTIP,
		trayIcon:  wellzIcon,
		projects:  projects,
		menuItems: make(map[string]*systray.MenuItem),
	}, nil
}

func (a *App) OnReady() {
	systray.SetIcon(a.trayIcon)
	systray.SetTitle(a.title)
	systray.SetTooltip(a.tooltip)

	// a.Reset()
}

func (a *App) startMenus() {
	a.createMenus()
	// go func() {
	// 	<-mChange.ClickedCh
	// 	mChange.SetTitle("I've Changed")
	// }()
	systray.AddSeparator()
	a.createQuitMenuItem()
}

func (a *App) createMenus() {
	for _, project := range a.projects {
		for _, repo := range project.Repos {
			m := systray.AddMenuItem(repo, repo)
			m.Hide()
			a.menuItems[repo] = m
		}
	}
}
func (a *App) Show(menu string) {
	menuItem := a.menuItems[menu]
	menuItem.Show()
}
func (a *App) AddMenuItem(menu, title, url string) {
	menuItem := a.menuItems[menu]
	prMenuItem := menuItem.AddSubMenuItem(title, title)
	// prMenuItem.Disabled()
	go func() {
		<-prMenuItem.ClickedCh
		a.openBrowser(url)
	}()
}

func (a *App) createQuitMenuItem() {
	mQuit := systray.AddMenuItem("Quit", "Quit")
	// Sets the icon of a menu item. Only available on Mac and Windows.
	mQuit.SetIcon(a.trayIcon)
	go func() {
		<-mQuit.ClickedCh
		fmt.Println("Requesting quit")
		systray.Quit()
		fmt.Println("Finished quitting")
	}()
}

func (a *App) OnExit() {}

func (a *App) Reset() {
	systray.ResetMenu()
	a.startMenus()
}

func (a *App) openBrowser(targetURL string) {
	var err error

	switch runtime.GOOS {
	case "linux":
		err = exec.Command("xdg-open", targetURL).Start()
		// TODO: "Windows Subsytem for Linux" is also recognized as "linux", but then we need
		// err = exec.Command("rundll32.exe", "url.dll,FileProtocolHandler", targetURL).Start()
	case "windows":
		err = exec.Command("rundll32.exe", "url.dll,FileProtocolHandler", targetURL).Start()
	case "darwin":
		err = exec.Command("open", targetURL).Start()
	default:
		err = fmt.Errorf("unsupported platform %v", runtime.GOOS)
	}
	if err != nil {
		log.Fatal(err)
	}

}
