// Copyright 2022-present miodzie. All rights reserved.
// This source code is licensed under the Apache 2.0 license found
// in the LICENSE file in the root directory of this source tree.

package dong

import (
	"github.com/jinzhu/gorm"
	"github.com/miodzie/dong"
	"github.com/miodzie/dong/impl"
	"github.com/miodzie/dong/usecases"
	"github.com/miodzie/walter"
	"github.com/miodzie/walter/log"
	"math/rand"
	"os"
	"os/user"
	"path"
	"time"
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

func (mod *Mod) Start(stream walter.Stream, actions walter.Actions) error {
	mod.running = true
	for mod.running {
		msg := <-stream
		if msg.IsCommand("8ball") || msg.IsCommand("8") {
			_ = actions.Reply(msg, get8BallResponse())
		}
		if msg.IsCommand("dong") {
			rando := usecases.NewRandomDongUseCase(mod.repository)
			var request usecases.RandomDongReq
			if len(msg.Arguments) > 1 {
				request.Category = msg.Arguments[1]
			}
			response := rando.Handle(request)
			if err := actions.Reply(msg, response.Emoji); err != nil {
				log.Error(err)
			}
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
		log.Error(err)
		panic("failed to connect database")
	}

	return impl.NewGormRepository(db)
}

type ModFactory struct {
}

func (m ModFactory) Create(interface{}) (walter.Module, error) {
	return New(), nil
}

func get8BallResponse() string {
	responses := []string{
		"It is certain.",
		"It is decidedly so.",
		"Without a doubt.",
		"Yes - definitely.",
		"You may rely on it.",
		"As I see it, yes.",
		"Most likely.",
		"Outlook good.",
		"Yes.",
		"Signs point to yes.",
		"Reply hazy, try again.",
		"Better not tell you now.",
		"Cannot predict now.",
		"Concentrate and ask again.",
		"Don't count on it.",
		"My reply is no.",
		"My sources say no.",
		"Outlook not so good.",
		"Very doubtful.",
	}

	rand.Seed(time.Now().UnixNano())
	return responses[rand.Intn(len(responses))]
}
