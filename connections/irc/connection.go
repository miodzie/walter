package irc

import (
	"crypto/tls"
	"github.com/miodzie/seras/connections/irc/plugin"
	"strings"
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
	// TODO: Remove me.
	name string
}

func (con *Connection) Name() string {
	return con.name
}

func (con *Connection) SetName(s string) {
	con.name = s
}

func New(conf Config) (*Connection, error) {
	ircCon := irc.IRC(conf.Nick, conf.Username)
	ircCon.UseTLS = true
	ircCon.TLSConfig = &tls.Config{InsecureSkipVerify: true}
	ircCon.UseSASL = conf.SASL
	// TODO: Update to use internal logger or NilledLogger
	//ircCon.Log = nil
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
		var channel string
		var args = event.Arguments
		if event.Code == "PRIVMSG" {
			args = strings.Split(event.Arguments[1], " ")
			channel = event.Arguments[0]
		}
		stream <- seras.Message{
			Content:   event.Message(),
			Arguments: args,
			Target:    channel,
			Author: seras.Author{
				Id:      event.Host,
				Nick:    event.Nick,
				Mention: event.Nick,
			},
			Code:           event.Code,
			ConnectionName: con.Name(),
			Raw:            event.Raw,
			Timestamp:      time.Now(),
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
