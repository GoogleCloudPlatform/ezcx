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

// HandlerFunc is an adapter that converts a given ezcx.HandlerFunc into an http.Handler.
type HandlerFunc func(*WebhookResponse, *WebhookRequest) error

// Implementing ServeHTTP allows the ezcx.HandlerFunc to satisfy the http.Handler interface.
//
// Error handling is an area of future improvement.  For instance, if a required parameter
// is missing, it should be up to the developer to handle that i.e.: return an HTTP error (400, 500)
// or return a ResponseMessage indicating something went wrong...
func (h HandlerFunc) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()
	if r.Method != http.MethodPost {
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}
	req, err := WebhookRequestFromRequest(r)
	if err != nil {
		log.Println("Error during WebhookRequestFromRequest")
		log.Println(err)
		return
	}
	req.ctx = r.Context // flowing down the requests's Context added..
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
	ctx = context.WithValue(ctx, Logger, lg)
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

// SetHandler allows the user to set a custom mux or handler.
func (s *Server) SetHandler(h http.Handler) {
	s.server.Handler = h
	if s.isMux(h) {
		s.mux = h.(*http.ServeMux)
	} else {
		s.mux = nil
	}
}

// ServeMux returns a copy of the currently set mux.
func (s *Server) ServeMux() *http.ServeMux {
	return s.mux
}

func (s *Server) isMux(h http.Handler) bool {
	_, ok := h.(*http.ServeMux)
	return ok
}

// HandleCx registers the handler for the given pattern.  While the HandleCx method itself
// isn't safe for concurrent usage, the underlying method it wraps (*ServeMux).Handle IS guarded
// by a mutex.
func (s *Server) HandleCx(pattern string, handler HandlerFunc) {
	s.mux.Handle(pattern, handler)
}

// ListenAndServe listens on the TCP network address srv.Addr and then calls Serve
// to handle requests on incoming connections. ListenAndServe is responsible for handling signals
// and managing graceful shutdown(s) whenever the right signals are intercepted.
func (s *Server) ListenAndServe(ctx context.Context) {
	defer func() {
		close(s.errs)
		close(s.signal)
	}()
	// Run ListenAndServe on a separate goroutine.
	s.lg.Printf("EZCX server listening and serving on %s\n", s.server.Addr)
	go func() {
		err := s.server.ListenAndServe()
		if err != nil && err != http.ErrServerClosed {
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
				} else {
					s.lg.Println("EZCX shutdown SUCCESS")
				}
				return
			}
		}
	}
}

func (s *Server) ListenAndServeTLS(ctx context.Context, certFile, keyFile string) {
	defer func() {
		close(s.errs)
		close(s.signal)
	}()
	// Run ListenAndServe on a separate goroutine.
	s.lg.Printf("EZCX server listening and serving on %s\n", s.server.Addr)
	go func() {
		err := s.server.ListenAndServeTLS(certFile, keyFile)
		if err != nil && err != http.ErrServerClosed {
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
				} else {
					s.lg.Println("EZCX shutdown SUCCESS")
				}
				return
			}
		}
	}
}

// Omitted for now.
func (s *Server) Reconfigure() error {
	return nil
}

// Shutdown provides graceful shutdown for the entire ezcx Server
func (s *Server) Shutdown(ctx context.Context) error {
	timeout, cancel := context.WithTimeout(ctx, time.Second*5)
	defer cancel()
	err := s.server.Shutdown(timeout)
	if err != nil {
		return err
	}
	return nil
}
