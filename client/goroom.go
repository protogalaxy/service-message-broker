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

package client

import (
	"errors"
	"fmt"
	"net/http"
	"strings"

	"github.com/arjantop/cuirass"
	"github.com/arjantop/saola/httpservice"
	"github.com/protogalaxy/common/serviceerror"
	"golang.org/x/net/context"
)

type GoRoom interface {
	JoinRoom(ctx context.Context, roomID, userID string) error
}

var _ GoRoom = (*GoRoomClient)(nil)

type GoRoomClient struct {
	Client   *httpservice.Client
	Executor cuirass.Executor
}

func (c *GoRoomClient) JoinRoom(ctx context.Context, roomID, userID string) error {
	cmd := cuirass.NewCommand("JoinRoom", func(ctx context.Context) (interface{}, error) {
		requestJson := fmt.Sprintf(`{"user_id": "%s"}`, userID)
		url := fmt.Sprintf("http://localhost:10200/%s/join", roomID)
		req, err := http.NewRequest("POST", url, strings.NewReader(requestJson))
		if err != nil {
			return nil, errors.New("Problem creating request: " + err.Error())
		}
		req.Header.Set("Content-Type", "application/json; charset=utf-8")

		res, err := c.Client.Do(ctx, req)
		defer res.Body.Close()
		if err != nil {
			return nil, err
		} else if res.StatusCode == http.StatusOK {
			return nil, nil
		} else {
			return nil, serviceerror.Decode(res.Body)
		}
	}).Build()

	_, err := c.Executor.Exec(ctx, cmd)
	return err
}
