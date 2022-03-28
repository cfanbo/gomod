package core

import (
	"fmt"
	"github.com/cfanbo/gomod/pkg"
	"os"
	"strings"
	"sync"

	"github.com/gdamore/tcell/v2"
	"github.com/gdamore/tcell/v2/encoding"
)

var onceRender sync.Once

var render *Render

type Render struct {
	screen tcell.Screen
	repos  *Repos
	pager  *pkg.Pager

	headerHeight int
}

func NewRender() *Render {
	onceRender.Do(func() {
		screen, e := tcell.NewScreen()
		if e != nil {
			fmt.Fprintf(os.Stderr, "%v\n", e)
			os.Exit(1)
		}

		if e := screen.Init(); e != nil {
			fmt.Fprintf(os.Stderr, "%v\n", e)
			os.Exit(1)
		}

		render = &Render{
			screen: screen,
		}
	})

	return render
}

func (r *Render) SetHeaderHeight(h int) {
	r.headerHeight = h
}

func (r *Render) SetRepos(repos *Repos) {
	r.repos = repos
}

func (r *Render) SetPager(pager *pkg.Pager) {
	r.pager = pager
}

func (r *Render) render() {
	var selectedStyle tcell.Style
	selectedStyle = tcell.StyleDefault.Background(tcell.ColorWhite).Foreground(tcell.ColorBlack)

	screen := r.screen
	screen.Clear()

	// top text
	table := NewTable(screen, NewPos(0, 0))
	table.AddRow().SetContent("Press `ESC` to quit.  Press the [s|f|m|g] key to sort").Print()

	// statistics
	var success int
	for _, v := range r.pager.Data() {
		if repo, ok := v.(*Repo); ok && repo.Done {
			success++
		}
	}
	table.AddRow().SetContent(fmt.Sprintf("Modules: %d total, %d success, %d failed.  Display range %d ~ %d",
		len(r.pager.Data()), success, len(r.pager.Data())-success, r.pager.StartIndex()+1, r.pager.EndIndex())).Print()
	table.AddRow().Print()

	// header
	row := table.AddRow().
		AddColWithWidth("STAR", 7).
		AddColWithWidth("FORK", 7).
		AddColWithWidth("SHARE", 7).
		AddColWithWidth("MODULE", 40)

	//AddColWithWidth("GitHub", 50).Print()
	pos := NewPos(row.getColX(50), row.pos.GetY())
	row.AddColl("GITHUB", WithWidth(50), WithPos(pos))
	row.Print()

	table.AddRow().SetContent(strings.Repeat("-", 111)).Print()

	pageRepos := r.pager.Result()
	for i, v := range pageRepos {
		repo, ok := v.(*Repo)
		if !ok {
			continue
		}

		row := table.AddRow()
		if i+r.pager.StartIndex() == r.pager.SelectedIndex() {
			row.SetStyle(selectedStyle)
		}

		row.AddColWithWidth(repo.GetStar(), 7).
			AddColWithWidth(repo.GetFork(), 7).
			AddColWithWidth(repo.GetShared(), 7).
			AddColWithWidth(repo.Mod, 40).
			AddColWithWidth(repo.RepoUrl, 60).
			Print()
	}

	screen.Show()
}

func (r *Render) LoadPagerData() {
	var dataList []interface{}
	dataList = getRepoList(r.repos)

	r.pager.SetData(dataList)
}

func (r *Render) Run() {
	r.SetHeaderHeight(5)
	r.repos.sort()

	encoding.Register()

	var dataList []interface{}
	dataList = getRepoList(r.repos)
	pager := pkg.NewPager(dataList)
	r.pager = pager

	var screenHeight int
	_, screenHeight = r.screen.Size()
	pager.PageSize(screenHeight - r.headerHeight)

	screen := r.screen
	for {
		switch ev := screen.PollEvent().(type) {
		case *tcell.EventResize:
			screen.Sync()
			_, screenHeight = screen.Size()
			pager.PageSize(screenHeight - r.headerHeight)
			r.render()
		case *tcell.EventKey:
			switch ev.Key() {
			case tcell.KeyDown:
				screen.Sync()
				pager.Next()
				r.render()
			case tcell.KeyUp:
				screen.Sync()
				pager.Prev()
				r.render()
			case tcell.KeyEsc:
				screen.Fini()
				os.Exit(0)
			case tcell.KeyEnter:
				go openWindow(pager, "github")
			case tcell.KeyPgUp:
				pager.PgUp()
				r.render()
			case tcell.KeyPgDn:
				pager.PgDown()
				r.render()
			default:
				if ev.Rune() == 's' || ev.Rune() == 'f' || ev.Rune() == 'm' || ev.Rune() == 'g' {
					screen.Sync()
					r.repos.detectKey(ev.Rune())

					dataList = getRepoList(r.repos)
					pager.SetData(dataList)

					r.render()
				} else if ev.Rune() == 'S' || ev.Rune() == 'F' || ev.Rune() == 'M' || ev.Rune() == 'G' {
					screen.Sync()
					r.repos.detectKey(ev.Rune() + 32)

					dataList = getRepoList(r.repos)
					pager.SetData(dataList)
					r.render()
				} else if ev.Rune() == ' ' {
					openWindow(pager, "pkggodev")
				}
			}
		}
	}
}
