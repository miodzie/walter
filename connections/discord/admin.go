// Copyright 2022-present miodzie. All rights reserved.
// This source code is licensed under the Apache 2.0 license found
// in the LICENSE file in the root directory of this source tree.

package discord

import "time"

func (con *Connection) IsAdmin(userId string) bool {
	for _, a := range con.Admins {
		if a == userId {
			return true
		}
	}
	return false
}

func (con *Connection) TimeoutUser(channel string, user string, until time.Time) error {
	return con.session.GuildMemberTimeout(channel, user, &until)
}
