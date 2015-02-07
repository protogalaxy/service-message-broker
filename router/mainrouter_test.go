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
	"testing"

	"github.com/protogalaxy/service-message-broker/router"
	"github.com/protogalaxy/service-message-broker/router/routertest"
	"golang.org/x/net/context"
)

func TestMainRouterJoinRoom(t *testing.T) {
	r := &router.MainRouter{
		RoomRouter: &routertest.MessageRouterMock{
			OnRoute: func(ctx context.Context, data []byte) ([]byte, error) {
				if string(data) != `"test"` {
					t.Errorf("Unexpected routing data: %#v", data)
				}
				return []byte("test"), nil
			},
		},
	}
	resp, err := r.Route(nil, []byte(`{"type":"room.join","data":"test"}`))
	if err != nil {
		t.Fatalf("Unexpected error: %v", err)
	}
	if string(resp) != "test" {
		t.Errorf("Invalid response data: %#v", resp)
	}
}

func TestMainRouterErrorDecodingInputMessage(t *testing.T) {
	r := &router.MainRouter{}
	_, err := r.Route(nil, []byte(`{`))
	if err == nil {
		t.Fatal("Expecting error but got none")
	}
}
