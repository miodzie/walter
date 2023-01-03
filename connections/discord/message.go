// Copyright 2022-present miodzie. All rights reserved.
// This source code is licensed under the Apache 2.0 license found
// in the LICENSE file in the root directory of this source tree.

package discord

import (
	"fmt"
	"github.com/bwmarrin/discordgo"
	"strings"
	"walter"
	"walter/log"
)

func (con *Connection) Send(msg walter.Message) error {
	_, err := con.session.ChannelMessageSend(msg.Target, msg.Content)
	log.Debugf("[%s]: %+v\n", con.Name(), msg)
	return err
}

func (con *Connection) Reply(msg walter.Message, content string) error {
	reply := walter.Message{Content: content, Target: msg.Target}
	return con.Send(reply)
}

func (con *Connection) onMessageCreate(s *discordgo.Session, e *discordgo.MessageCreate) {
	if e.Author.Bot {
		return
	}
	msg := walter.Message{
		Content:   e.Content,
		Target:    e.ChannelID,
		Arguments: strings.Split(e.Content, " "),
		Author: walter.Author{
			Id:      e.Author.ID,
			Nick:    e.Author.Username,
			Mention: "<@" + e.Author.ID + ">",
		},
		ConnectionName: con.Name(),
		Raw:            e.Content,
		Timestamp:      e.Timestamp,
	}
	con.stream <- msg
}

// walter.MessageFormatter

func (con *Connection) Bold(str string) string {
	return fmt.Sprintf("**%s**", str)
}

func (con *Connection) Italicize(str string) string {
	return fmt.Sprintf("**%s**", str)
}
