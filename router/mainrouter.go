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
	"errors"

	"github.com/protogalaxy/service-message-broker/Godeps/_workspace/src/golang.org/x/net/context"
	"github.com/protogalaxy/service-message-broker/tictactoe"
)

type MessageRouter interface {
	Route(ctx context.Context, data []byte) ([]byte, error)
}

var _ MessageRouter = (*MainRouter)(nil)

type MainRouter struct {
	Client tictactoe.GameManagerClient
}

type MessageType struct {
	Type string `json:"type"`
}

func (r *MainRouter) Route(ctx context.Context, data []byte) ([]byte, error) {
	var t MessageType
	err := json.Unmarshal(data, &t)
	if err != nil {
		return nil, err
	}
	return r.routeByType(ctx, t.Type, data)
}

type CreateGame struct {
	WithUser string `json:"with_user"`
}

type Turn struct {
	GameID string `json:"game_id"`
	MoveID int64  `json:"move_id"`
	Move   struct {
		X int `json:"x"`
		Y int `json:"y"`
	} `json:"move"`
}

func (r *MainRouter) routeByType(ctx context.Context, t string, data []byte) ([]byte, error) {
	switch t {
	case "create":
		var createGame CreateGame
		err := json.Unmarshal(data, &createGame)
		if err != nil {
			return nil, err
		}

		_, err = r.Client.CreateGame(ctx, &tictactoe.CreateRequest{
			UserIds: []string{"", createGame.WithUser},
		})
		if err != nil {
			return nil, err
		}
	case "turn":
		var turn Turn
		err := json.Unmarshal(data, &turn)
		if err != nil {
			return nil, err
		}

		_, err = r.Client.PlayTurn(ctx, &tictactoe.TurnRequest{
			UserId: "",
			GameId: turn.GameID,
			MoveId: turn.MoveID,
			Move: &tictactoe.TurnRequest_Square{
				X: int32(turn.Move.X),
				Y: int32(turn.Move.Y),
			},
		})
	}
	return nil, errors.New("unknown message type")
}
