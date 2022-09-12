package irc

import (
	"errors"
	"fmt"
	"time"

	"github.com/miodzie/seras"
	irc "github.com/thoj/go-ircevent"
	"sync"
)

type Connection struct {
	irc    *irc.Connection
	config *Config
	mods   []seras.Module
	sync.Mutex
}

func New(conf Config) (*Connection, error) {
	con := &Connection{
		irc:    irc.IRC(conf.Nick, conf.Username),
		config: &conf,
	}

	return con, nil
}

func (con *Connection) Connect() (seras.Stream, error) {
	con.Lock()
	defer con.Unlock()
	err := con.irc.Connect(con.config.Server)
	if err != nil {
		return nil, err
	}
	stream := make(chan seras.Message)

	con.irc.AddCallback("*", func(event *irc.Event) {
		fmt.Println(event.Raw)
		stream <- seras.Message{
			Content:   event.Message(),
			Arguments: event.Arguments,
		}
	})

	go func() {
		con.irc.Loop()
	}()

	return stream, nil
}

func (con *Connection) Close() error {
	con.Lock()
	defer con.Unlock()
	con.irc.Disconnect()
	fmt.Println("why hang")
	con.irc.ClearCallback("*")

	return nil
}

func (con *Connection) Send(msg seras.Message) error {
	return errors.New("not implemented")
}
func (con *Connection) Reply(msg seras.Message, content string) error {
	return errors.New("not implemented")
}

func (con *Connection) Mods() []seras.Module {
	return con.mods
}
func (con *Connection) AddMods(mods []seras.Module) {
	con.mods = append(con.mods, mods...)
}

func (con *Connection) IsAdmin(userId string) bool {
	return false
}
func (con *Connection) TimeoutUser(channel string, user string, until time.Time) error {
	return errors.New("not implemented")
}

func (con *Connection) Bold(msg string) string {
	return msg
}
func (con *Connection) Italicize(msg string) string {
	return msg
}
