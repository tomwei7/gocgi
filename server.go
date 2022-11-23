package gocgi

import (
	"strings"
	"context"
	"log"
	"net/http"
	"net/http/cgi"
	"os"
)

type muxRoute struct{
	cgiHandler http.Handler
	staticHandler map[string]http.Handler
}

func(m *muxRoute) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	for pattern, handler := range m.staticHandler {
		if strings.HasPrefix(r.URL.Path, pattern) {
			handler.ServeHTTP(w, r)
			return
		}
	}
	m.cgiHandler.ServeHTTP(w, r)
}

func(m *muxRoute) HandleStatic(pattern, dir string) {
	m.staticHandler[pattern] = http.StripPrefix(pattern, http.FileServer(http.Dir(dir)))
}

func newMuxRoute(cgiHandler http.Handler) *muxRoute{
	return &muxRoute{cgiHandler: cgiHandler, staticHandler: make(map[string]http.Handler)}
}


type Server struct {
	opts   *GoCGIOptions
	logger *log.Logger

	cgiLogFile *os.File

	httpServer *http.Server
}

func (s *Server) ListenAndServe() error {
	s.logger.Printf("listen and serve cgi-bin %s at %s, static map %v", s.opts.Path, s.opts.Addr, s.opts.StaticMap)
	return s.httpServer.ListenAndServe()
}

func (s *Server) Shutdown(ctx context.Context) error {
	err := s.httpServer.Shutdown(ctx)
	if s.cgiLogFile != nil {
		s.cgiLogFile.Close()
	}
	return err
}

func New(logger *log.Logger, opts *GoCGIOptions) (*Server, error) {
	var cgiLogFile *os.File
	if opts.Stderr != "" {
		var err error
		if cgiLogFile, err = os.OpenFile(opts.Stderr, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0644); err != nil {
			return nil, err
		}
	}

	cgiHandler := &cgi.Handler{
		Path:       opts.Path,
		Root:       opts.Root,
		Dir:        opts.Dir,
		Env:        opts.Env,
		InheritEnv: opts.InheritEnv,
		Stderr:     cgiLogFile,
	}

	mux := newMuxRoute(cgiHandler)

	for dir, pattern := range opts.StaticMap {
		mux.HandleStatic(pattern, dir)
	}

	handler := WithBasicAuth(opts.Users, mux)
	handler = WithLogger(logger, handler)

	httpServer := &http.Server{
		Addr:    opts.Addr,
		Handler: handler,
	}

	server := &Server{
		opts:       opts,
		logger:     logger,
		cgiLogFile: cgiLogFile,
		httpServer: httpServer,
	}
	return server, nil
}
