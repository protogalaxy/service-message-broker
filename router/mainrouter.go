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
	"strings"

	"golang.org/x/net/context"
)

type MessageRouter interface {
	Route(ctx context.Context, data []byte) ([]byte, error)
}

var _ MessageRouter = (*MainRouter)(nil)

type MainRouter struct {
	RoomRouter MessageRouter
}

type Message struct {
	Type string          `json:"type"`
	ID   string          `json:"id"`
	Data json.RawMessage `json:"data"`
}

func (r *MainRouter) Route(ctx context.Context, data []byte) ([]byte, error) {
	var msg Message
	err := json.Unmarshal(data, &msg)
	if err != nil {
		return nil, err
	}
	return r.routeType(ctx, &msg)
}

func (r *MainRouter) routeType(ctx context.Context, msg *Message) ([]byte, error) {
	switch strings.SplitN(msg.Type, ".", 2)[0] {
	case "room":
		return r.RoomRouter.Route(ctx, msg.Data)
	default:
		panic("unknown message")
	}
}
