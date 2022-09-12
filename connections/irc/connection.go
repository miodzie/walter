package irc

import (
	"crypto/tls"
	"errors"
	"fmt"
	"github.com/miodzie/seras/connections/irc/plugin"
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
	ircCon := irc.IRC(conf.Nick, conf.Username)
	ircCon.UseTLS = true
	ircCon.TLSConfig = &tls.Config{InsecureSkipVerify: true}
	ircCon.UseSASL = conf.SASL
	ircCon.SASLLogin = conf.SASLUsername
	ircCon.SASLPassword = conf.SASLPassword
	con := &Connection{
		irc:    ircCon,
		config: &conf,
		mods:   []seras.Module{plugin.New(ircCon)},
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
		// TODO: Remove me.
		args := event.Arguments
		if event.Code == "PRIVMSG" {
			args = event.Arguments[1:]
		}
		fmt.Println(event.Raw)
		stream <- seras.Message{
			Content:   event.Message(),
			Arguments: args,
			Target:    event.Arguments[0],
			Author: seras.Author{
				Id:      event.Host,
				Nick:    event.Nick,
				Mention: "",
			},
			Code: event.Code,
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
	con.irc.ClearCallback("*")

	return nil
}

func (con *Connection) Send(msg seras.Message) error {
	con.irc.Privmsg(msg.Target, msg.Content)
	return nil
}
func (con *Connection) Reply(msg seras.Message, content string) error {
	reply := seras.Message{Content: content, Target: msg.Target}
	return con.Send(reply)
}

func (con *Connection) Mods() []seras.Module {
	return con.mods
}
func (con *Connection) AddMods(mods []seras.Module) {
	con.mods = append(con.mods, mods...)
}

func (con *Connection) IsAdmin(userId string) bool {
	for _, a := range con.config.Admins {
		if a == userId {
			return true
		}
	}
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
