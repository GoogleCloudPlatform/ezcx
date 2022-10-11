// Copyright 2022 Google LLC
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     https://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package ezcx

import (
	"context"
	"log"
	"net"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"
)

var (
	ServerDefaultSignals []os.Signal = []os.Signal{
		syscall.SIGINT,
		syscall.SIGTERM,
		syscall.SIGHUP,
	}
)

//Handler interface: Defines the contract.
type Handler interface {
	http.Handler
	Handle(res *WebhookResponse, req *WebhookRequest) error
}

// Functional Adapter: HandlerFunc is an adapter.
// HandlerFunc satisfies the Handler interface
type HandlerFunc func(*WebhookResponse, *WebhookRequest) error

// Seems  redundant; may serve a purpose, though, for structural handlers.
// (ie Need to implement for functional handler to satisfy Handle which would
// require implementation for structural handlers.)
func (h HandlerFunc) Handle(res *WebhookResponse, req *WebhookRequest) error {
	return h(res, req)
}

// yaquino@2022-10-07: http.Request's context is flowd down to the WebhookRequest
// via WebhookRequestFromRequest (requests.go)
func (h HandlerFunc) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()
	req, err := WebhookRequestFromRequest(r)
	if err != nil {
		log.Println("Error during WebhookRequestFromRequest")
		log.Println(err)
		return
	}
	req.ctx = r.Context // flowing down the requests's Context added.
	res := req.InitializeResponse()
	err = h(res, req)
	if err != nil {
		log.Println("Error during HandlerFunc execution")
		log.Println(err)
		return
	}
	err = res.WriteResponse(w)
	if err != nil {
		log.Println("Error during WebhookResponse.WriteResponse")
		return
	}
}

type Server struct {
	signals []os.Signal
	signal  chan os.Signal
	errs    chan error
	server  *http.Server
	mux     *http.ServeMux
	lg      *log.Logger
}

func NewServer(ctx context.Context, addr string, lg *log.Logger, signals ...os.Signal) *Server {
	return new(Server).Init(ctx, addr, lg, signals...)
}

func (s *Server) Init(ctx context.Context, addr string, lg *log.Logger, signals ...os.Signal) *Server {
	if len(signals) == 0 {
		s.signals = ServerDefaultSignals
	} else {
		// rethink this later on.  We need to make sure there at least
		// the right group of signals!
		s.signals = signals
	}
	s.signal = make(chan os.Signal, 1)
	signal.Notify(s.signal, s.signals...)

	if lg == nil {
		lg = log.Default()
	}
	s.lg = lg

	s.errs = make(chan error)
	s.mux = http.NewServeMux()
	s.server = &http.Server{
		Addr:        addr,
		Handler:     s.mux,
		BaseContext: func(l net.Listener) context.Context { return ctx },
	}
	return s
}

func (s *Server) SetHandler(h http.Handler) {
	s.server.Handler = h
	if s.isMux(h) {
		s.mux = h.(*http.ServeMux)
	} else {
		s.mux = nil
	}
}

func (s *Server) ServeMux() *http.ServeMux {
	return s.mux
}

func (s *Server) isMux(h http.Handler) bool {
	_, ok := h.(*http.ServeMux)
	return ok
}

func (s *Server) HandleCx(pattern string, handler HandlerFunc) {
	s.mux.Handle(pattern, handler)
}

// yaquino@2022-09-21: I have concerns that checking the parent context will not work as desired.
func (s *Server) ListenAndServe(ctx context.Context) {
	defer func() {
		close(s.errs)
		close(s.signal)
	}()
	// Run ListenAndServe on a separate goroutine.
	s.lg.Printf("EZCX server listening and serving on %s\n", s.server.Addr)
	go func() {
		err := s.server.ListenAndServe()
		if err != nil {
			s.lg.Println(err)
			s.errs <- err
			close(s.errs)
		}
	}()

	for {
		select {
		// If the context is done, we need to return.
		case <-ctx.Done():
			s.lg.Println("EZCX server context is done")
			err := ctx.Err()
			if err != nil {
				s.lg.Print("EZCX server context error...")
				s.lg.Println(err)
			}
			return
		// If there's a non-nil error, we need to return
		case err := <-s.errs:
			if err != nil {
				s.lg.Print("EZCX server non-nil error...")
				s.lg.Println(err)
				return
			}
		case sig := <-s.signal:
			s.lg.Printf("EZCX server signal %s received...", sig)
			switch sig {
			case syscall.SIGHUP:
				s.lg.Println("EZCX reconfigure", sig)
				err := s.Reconfigure()
				if err != nil {
					s.errs <- err
				}
			default:
				s.lg.Printf("EZCX graceful shutdown initiated...")
				err := s.Shutdown(ctx)
				if err != nil {
					s.lg.Println(err)
				}
				s.lg.Println("EZCX shutdown SUCCESS")
				return
			}
		}
	}
}

// Omitted for now.
func (s *Server) Reconfigure() error {
	return nil
}

func (s *Server) Shutdown(ctx context.Context) error {
	timeout, cancel := context.WithTimeout(ctx, time.Second*5)
	defer cancel()
	err := s.server.Shutdown(timeout)
	if err != nil {
		return err
	}
	return nil
}
