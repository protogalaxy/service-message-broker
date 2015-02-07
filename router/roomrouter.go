// Copyright (C) 2015 The Protogalaxy Project
//
// This program is free software: you can redistribute it and/or modify
// it under the terms of the GNU Affero General Public License as published by
// the Free Software Foundation, either version 3 of the License, or
// (at your option) any later version.
//
// This program is distributed in the hope that it will be useful,
// but WITHOUT ANY WARRANTY; without even the implied warranty of
// MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
// GNU Affero General Public License for more details.
//
// You should have received a copy of the GNU Affero General Public License
// along with this program.  If not, see <http://www.gnu.org/licenses/>.

package router

import (
	"encoding/json"

	"github.com/protogalaxy/service-message-broker/client"
	"golang.org/x/net/context"
)

var _ MessageRouter = (*RoomRouter)(nil)

type RoomRouter struct {
	GoRoomClient client.GoRoom
}

type JoinRoomMessage struct {
	RoomID string `json:"room_id"`
	UserID string `json:"user_id"`
}

func (r *RoomRouter) Route(ctx context.Context, data []byte) ([]byte, error) {
	var msg JoinRoomMessage
	err := json.Unmarshal(data, &msg)
	if err != nil {
		return nil, err
	}
	err = r.GoRoomClient.JoinRoom(ctx, msg.RoomID, msg.UserID)
	if err != nil {
		return nil, err
	}
	return []byte("{}"), nil
}
