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
	config  *Config
	sync.Mutex
}

func New(token string) (*Connection, error) {
	var disc Connection
	session, err := discordgo.New("Bot " + token)
	if err != nil {
		return &disc, err
	}
	disc.session = session

	return &disc, nil
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
		Content:       e.Content,
		Channel:       e.ChannelID,
		Arguments:     strings.Split(e.Content, " "),
		AuthorId:      e.Author.ID,
		AuthorNick:    e.Author.Username,
		AuthorMention: "<@" + e.Author.ID + ">",
	}
	fmt.Printf("Discord:  [%s]: %s\n", msg.AuthorNick, msg.Content)
	con.stream <- msg
}

func (con *Connection) Send(msg seras.Message) error {
	_, err := con.session.ChannelMessageSend(msg.Channel, msg.Content)
	return err
}

func (con *Connection) Reply(msg seras.Message, content string) error {
	reply := seras.Message{Content: content, Channel: msg.Channel}
	return con.Send(reply)
}
func (con *Connection) IsAdmin(userId string) bool {
	_, ok := con.config.Admins[userId]

	return ok
}

func (con *Connection) TimeoutUser(channel string, user string, until time.Time) error {
	return con.session.GuildMemberTimeout(channel, user, &until)
}
