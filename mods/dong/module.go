package dong

import (
	"fmt"
	"github.com/jinzhu/gorm"
	"github.com/miodzie/dong"
	"github.com/miodzie/dong/impl"
	"github.com/miodzie/dong/interactors"
	"github.com/miodzie/seras"
	"os"
	"os/user"
	"path"
)

const WORKDIR = ".dong"

var workDir string

type Mod struct {
	running    bool
	repository dong.Repository
	fetcher    dong.Fetcher
}

func New() *Mod {
	createWorkDir()

	return &Mod{repository: initDatabase()}
}

func (mod *Mod) Name() string {
	return "dong"
}

func (mod *Mod) Start(stream seras.Stream, actions seras.Actions) error {
	mod.running = true
	for mod.running {
		msg := <-stream
		if msg.IsCommand("dong") {
			rando := interactors.NewRandomDongInteractor(mod.repository)
			var request interactors.RandomDongReq
			if len(msg.Arguments) > 1 {
				request.Category = msg.Arguments[1]
			}
			response := rando.Handle(request)
			// TODO: Handle error.
			actions.Reply(msg, response.Emoji)
		}

	}

	return nil
}

func (mod *Mod) Stop() {
	mod.running = false
}

func createWorkDir() {
	usr, err := user.Current()
	if err != nil {
		panic(err)
	}
	workDir = path.Join(usr.HomeDir, WORKDIR)

	if _, err := os.Stat(workDir); os.IsNotExist(err) {
		err := os.Mkdir(workDir, os.ModePerm)
		if err != nil {
			panic(err)
		}
	}
}

func initDatabase() *impl.GormRepository {
	db, err := gorm.Open("sqlite3", path.Join(workDir, "dongs.db"))
	if err != nil {
		fmt.Println(err)
		panic("failed to connect database")
	}

	return impl.NewGormRepository(db)
}
