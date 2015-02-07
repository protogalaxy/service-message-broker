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
package service_test

import (
	"errors"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/protogalaxy/common/serviceerror"
	"github.com/protogalaxy/service-message-broker/router/routertest"
	"github.com/protogalaxy/service-message-broker/service"
	"golang.org/x/net/context"
)

func TestRouteMessage(t *testing.T) {
	s := &service.RouteMessage{
		Router: &routertest.MessageRouterMock{
			OnRoute: func(ctx context.Context, msg []byte) ([]byte, error) {
				if string(msg) != "msg" {
					t.Errorf("Unexpected routing message: %s", string(msg))
				}
				return []byte("msgresp"), nil
			},
		},
	}
	w := httptest.NewRecorder()
	req, _ := http.NewRequest("POST", "", strings.NewReader("msg"))
	err := s.DoHTTP(context.Background(), w, req)
	if err != nil {
		t.Fatalf("Unexpected error: %s", err)
	}
	if w.Code != http.StatusOK {
		t.Errorf("Invalid response code: %d", w.Code)
	}
	if ct := w.Header().Get("Content-Type"); ct != "application/octet-stream" {
		t.Errorf("Invalid response content type: %s", ct)
	}
	if msg := w.Body.String(); msg != "msgresp" {
		t.Errorf("Invalid response body: %s", msg)
	}
}

func TestRouteMessageRoutingError(t *testing.T) {
	expectedError := errors.New("err")
	s := &service.RouteMessage{
		Router: &routertest.MessageRouterMock{
			OnRoute: func(ctx context.Context, data []byte) ([]byte, error) {
				return nil, expectedError
			},
		},
	}
	w := httptest.NewRecorder()
	req, _ := http.NewRequest("POST", "", strings.NewReader("msg"))
	err := s.DoHTTP(context.Background(), w, req)

	re, ok := err.(serviceerror.ErrorResponse)
	if !ok {
		t.Fatalf("The returned error should be of type ErrorResponse: %v", err)
	}
	if re.StatusCode != http.StatusInternalServerError {
		t.Fatalf("Unexpected http status code: %d", re.StatusCode)
	}
	if re.ErrorCode != "server_error" {
		t.Fatalf("Unexpected response error code: %s", re.ErrorCode)
	}
	if re.Cause != expectedError {
		t.Fatalf("Unexpected response error cause: %v", re.Cause)
	}
}

type ErrorReader struct{}

func (e ErrorReader) Read(p []byte) (int, error) {
	return 0, errors.New("read")
}

func TestRouteMessageReadRequestError(t *testing.T) {
	s := &service.RouteMessage{}
	w := httptest.NewRecorder()

	req, _ := http.NewRequest("POST", "", ErrorReader{})
	err := s.DoHTTP(context.Background(), w, req)
	if err == nil || err.Error() != "read" {
		t.Fatalf("Unexpected error: %v", err)
	}
}

func TestRouteMessageEmptyRequestBodyError(t *testing.T) {
	s := &service.RouteMessage{}
	w := httptest.NewRecorder()
	req, _ := http.NewRequest("POST", "", strings.NewReader(""))
	err := s.DoHTTP(context.Background(), w, req)

	re, ok := err.(serviceerror.ErrorResponse)
	if !ok {
		t.Fatalf("The returned error should be of type ErrorResponse: %v", err)
	}
	if re.StatusCode != http.StatusBadRequest {
		t.Fatalf("Unexpected http status code: %d", re.StatusCode)
	}
	if re.ErrorCode != "invalid_request" {
		t.Fatalf("Unexpected response error code: %s", re.ErrorCode)
	}
	if re.Cause != nil {
		t.Fatalf("Unexpected response error cause: %v", re.Cause)
	}
}

type ErrorResponseWriter struct {
	*httptest.ResponseRecorder
}

func (e ErrorResponseWriter) Write(p []byte) (int, error) {
	return 0, errors.New("write")
}

func TestRouteMessageWriteResponseError(t *testing.T) {
	s := &service.RouteMessage{
		Router: &routertest.MessageRouterMock{
			OnRoute: func(ctx context.Context, msg []byte) ([]byte, error) {
				return []byte("resp"), nil
			},
		},
	}
	w := ErrorResponseWriter{httptest.NewRecorder()}
	req, _ := http.NewRequest("POST", "", strings.NewReader("msg"))
	err := s.DoHTTP(context.Background(), w, req)
	if err == nil || err.Error() != "write" {
		t.Fatalf("Unexpected error: %v", err)
	}
}
