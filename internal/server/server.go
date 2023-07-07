package server

import (
	"context"
	"net/http"
	"time"

	"github.com/AnisaForWork/user_orders/internal/config"

	log "github.com/sirupsen/logrus"
	"golang.org/x/sync/errgroup"
)

type Server struct {
	httpServer *http.Server
	timeOutSec int
}

// function called to create server and configure it
func NewServer(cfg *config.Server, handler http.Handler) *Server {
	server := http.Server{
		Addr:              cfg.HTTP.Host + ":" + cfg.HTTP.Port,
		Handler:           handler,
		MaxHeaderBytes:    cfg.HTTP.MaxHeaderBytes,
		ReadTimeout:       cfg.HTTP.ReadTimeout,
		WriteTimeout:      cfg.HTTP.WriteTimeout,
		ReadHeaderTimeout: cfg.HTTP.ReadHeaderTimeout,
	}
	return &Server{
		httpServer: &server,
		timeOutSec: cfg.TimeOutSec,
	}
}

// Run starts configured server
// recives context for cancelation
// returns either error occurred during run, or during shutdow, or nothing
func (s *Server) Run(ctx context.Context) error {
	g, _ := errgroup.WithContext(ctx)
	g.Go(func() error {
		log.Info("REST running...")
		return s.httpServer.ListenAndServe()
	})

	g.Go(func() error {
		<-ctx.Done()
		return s.Shutdown()
	})

	return g.Wait()
}

func (s *Server) Shutdown() error {
	ctxTimeout, cancel := context.WithTimeout(context.Background(), time.Duration(s.timeOutSec)*time.Second)
	defer cancel()

	return s.httpServer.Shutdown(ctxTimeout)
}
