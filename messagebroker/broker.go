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

package messagebroker

import (
	"errors"
	"fmt"

	"github.com/protogalaxy/service-message-broker/Godeps/_workspace/src/golang.org/x/net/context"
	"github.com/protogalaxy/service-message-broker/router"
)

type Broker struct {
	Router router.MessageRouter
}

func (b *Broker) Route(ctx context.Context, req *RouteRequest) (*RouteReply, error) {
	if len(req.Data) == 0 {
		return nil, errors.New("empty message")
	}

	_, err := b.Router.Route(ctx, req.Data)
	if err != nil {
		return nil, fmt.Errorf("routing message: %s", err)
	}

	var rep RouteReply
	return &rep, nil
}
