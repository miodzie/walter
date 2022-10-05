package irc

import (
	"crypto/tls"
	"errors"
	"github.com/miodzie/seras/connections/irc/plugin"
	"github.com/miodzie/seras/log"
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
	// TODO: Update to use internal logger or NilLogger
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

func (con *Connection) Send(msg seras.Message) error {
	// An \n cuts off an IRC message, therefor split and send it as multiple messages.
	if strings.Contains(msg.Content, "\n") {
		split := strings.Split(msg.Content, "\n")
		var anyErr error
		for _, s := range split {
			newMsg := msg
			newMsg.Content = s
			if err := con.Send(newMsg); err != nil {
				anyErr = err
			}
		}
		// Leave early.
		return anyErr
	}
	con.irc.Privmsg(msg.Target, msg.Content)
	log.Debugf("[%s]: %+v\n", con.Name(), msg)
	return nil
}

func (con *Connection) Reply(msg seras.Message, content string) error {
	reply := seras.Message{Content: content, Target: msg.Target}
	if isPM(msg) {
		reply.Target = msg.Author.Nick
	}
	return con.Send(reply)
}

func isPM(msg seras.Message) bool {
	return !strings.Contains(msg.Target, "#")
}

func (con *Connection) Mods() []seras.Module {
	return con.mods
}
func (con *Connection) AddMods(mods []seras.Module) {
	con.mods = append(con.mods, mods...)
}

func (con *Connection) ModList() map[string]interface{} {
	mods := make(map[string]interface{})
	for _, m := range con.config.Mods {
		mods[m] = struct{}{}
	}

	return mods
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
