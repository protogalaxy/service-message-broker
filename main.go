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
package main

import (
	"flag"
	"math/rand"
	"net"
	_ "net/http/pprof"
	"time"

	"github.com/protogalaxy/service-message-broker/Godeps/_workspace/src/github.com/golang/glog"
	"github.com/protogalaxy/service-message-broker/Godeps/_workspace/src/google.golang.org/grpc"
	"github.com/protogalaxy/service-message-broker/messagebroker"
	"github.com/protogalaxy/service-message-broker/router"
	"github.com/protogalaxy/service-message-broker/tictactoe"
)

func main() {
	flag.Parse()
	rand.Seed(time.Now().UnixNano())

	conn, err := grpc.Dial("localhost:9093")
	if err != nil {
		glog.Fatalf("could not connect: %v", err)
	}
	defer conn.Close()
	tttc := tictactoe.NewGameManagerClient(conn)

	socket, err := net.Listen("tcp", ":9090")
	if err != nil {
		glog.Fatalf("failed to listen: %v", err)
	}

	router := &router.MainRouter{
		Client: tttc,
	}

	broker := &messagebroker.Broker{
		Router: router,
	}

	grpcServer := grpc.NewServer()
	messagebroker.RegisterBrokerServer(grpcServer, broker)
	grpcServer.Serve(socket)
}
