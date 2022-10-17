package discord

import (
	"github.com/bwmarrin/discordgo"
	"github.com/miodzie/seras"
	"sync"
)

type Connection struct {
	session *discordgo.Session
	stream  chan seras.Message
	config  *Config
	mods    []seras.Module
	sync.Mutex
	name string
}

func (con *Connection) Name() string {
	return con.name
}

func (con *Connection) SetName(s string) {
	con.name = s
}

func New(config Config) (*Connection, error) {
	disc := &Connection{config: &config}
	session, err := discordgo.New("Bot " + config.Token)
	if err != nil {
		return disc, err
	}
	disc.session = session

	return disc, nil
}

func (con *Connection) Connect() (seras.Stream, error) {
	con.Lock()
	defer con.Unlock()

	con.session.AddHandler(con.onMessageCreate)
	con.stream = make(chan seras.Message)
	con.session.Identify.Intents = discordgo.IntentsGuildMessages | discordgo.IntentGuildMessageTyping | discordgo.IntentsDirectMessages

	return con.stream, con.session.Open()
}

func (con *Connection) Close() error {
	con.Lock()
	defer con.Unlock()

	close(con.stream)

	return con.session.Close()
}
