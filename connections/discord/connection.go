package discord

import (
	"fmt"
	"strings"
	"sync"
	"time"

	"github.com/bwmarrin/discordgo"
	"github.com/miodzie/seras"
)

type Connection struct {
	session *discordgo.Session
	stream  chan seras.Message
	sync.Mutex
}

func New(token string) *Connection {
	session, err := discordgo.New("Bot " + token)
	if err != nil {
		panic(err)
	}
	disc := &Connection{session: session}

	return disc
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

func (con *Connection) onMessageCreate(s *discordgo.Session, e *discordgo.MessageCreate) {
	if e.Author.Bot {
		return
	}
	msg := seras.Message{
		Content:   e.Content,
		Channel:   e.ChannelID,
		Arguments: strings.Split(e.Content, " "),
		Author:    e.Author.Username,
		AuthorId:  e.Author.ID,
	}
	fmt.Printf("Discord:  [%s]: %s\n", msg.Author, msg.Content)
	con.stream <- msg
}

func (con *Connection) Send(msg seras.Message) error {
	_, err := con.session.ChannelMessageSend(msg.Channel, msg.Content)
	return err
}

func (con *Connection) TimeoutUser(channel string, user string, until time.Time) error {
	return con.session.GuildMemberTimeout(channel, user, &until)
}
