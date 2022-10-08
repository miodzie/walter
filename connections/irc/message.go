package irc

import (
	"github.com/miodzie/seras"
	"github.com/miodzie/seras/log"
	"strings"
)

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

// seras.MessageFormatter
func (con *Connection) Bold(msg string) string {
	return msg
}
func (con *Connection) Italicize(msg string) string {
	return msg
}
