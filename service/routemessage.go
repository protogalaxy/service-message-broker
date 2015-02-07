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
package service

import (
	"io/ioutil"
	"net/http"

	"github.com/arjantop/saola/httpservice"
	"github.com/protogalaxy/common/serviceerror"
	"github.com/protogalaxy/service-message-broker/router"
	"golang.org/x/net/context"
)

type RouteMessage struct {
	Router router.MessageRouter
}

// DoHTTP implements saola.HttpService.
func (h *RouteMessage) DoHTTP(ctx context.Context, w http.ResponseWriter, r *http.Request) error {
	msg, err := ioutil.ReadAll(r.Body)
	if err != nil {
		return err
	}
	if len(msg) == 0 {
		return serviceerror.BadRequest("invalid_request", "Empty request body")
	}

	resp, err := h.Router.Route(ctx, msg)
	if err != nil {
		return serviceerror.InternalServerError("server_error", "Unable to route the message", err)
	}

	w.WriteHeader(http.StatusOK)
	w.Header().Set("Content-Type", "application/octet-stream")
	_, err = w.Write(resp)
	if err != nil {
		return err
	}
	return nil
}

// Do implements saola.Service.
func (h *RouteMessage) Do(ctx context.Context) error {
	return httpservice.Do(h, ctx)
}

// Name implements saola.Service.
func (h *RouteMessage) Name() string {
	return "routemessage"
}
