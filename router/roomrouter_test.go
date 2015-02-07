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

package router_test

import (
	"errors"
	"testing"

	"github.com/protogalaxy/service-message-broker/router"
	"golang.org/x/net/context"
)

type GoRoomClientMock struct {
	OnJoinRoom func(ctx context.Context, roomID, userID string) error
}

func (m *GoRoomClientMock) JoinRoom(ctx context.Context, roomID, userID string) error {
	return m.OnJoinRoom(ctx, roomID, userID)
}

func TestRoomRouterJoinRoom(t *testing.T) {
	r := router.RoomRouter{
		GoRoomClient: &GoRoomClientMock{
			OnJoinRoom: func(ctx context.Context, roomID, userID string) error {
				if roomID != "room123" {
					t.Errorf("Unexpected room id: %s", roomID)
				}
				if userID != "user1" {
					t.Errorf("Unexpected user id: %s", userID)
				}
				return nil
			},
		},
	}
	resp, err := r.Route(nil, []byte(`{"room_id":"room123", "user_id":"user1"}`))
	if err != nil {
		t.Fatalf("Unexpected error: %#v", err)
	}
	if string(resp) != "{}" {
		t.Errorf("Invalid response: %#v", resp)
	}
}

func TestRoomRouterJoinRoomClientError(t *testing.T) {
	r := router.RoomRouter{
		GoRoomClient: &GoRoomClientMock{
			OnJoinRoom: func(ctx context.Context, roomID, userID string) error {
				return errors.New("error")
			},
		},
	}
	_, err := r.Route(nil, []byte(`{"room_id":"room123", "user_id":"user1"}`))
	if err == nil || err.Error() != "error" {
		t.Fatalf("Invalid error: %#v", err)
	}
}

func TestRoomRouterjoinRoomErrorDecodingInputMessage(t *testing.T) {
	r := router.RoomRouter{}
	_, err := r.Route(nil, []byte(`{`))
	if err == nil {
		t.Fatal("Expecting error but got none")
	}
}
