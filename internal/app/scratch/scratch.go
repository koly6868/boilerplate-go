package scratch

import (
	"context"
	"net/http"

	"github.com/opentracing/opentracing-go"
	log "gitlab.com/and07/boilerplate-go/internal/pkg/logger"
)

// Scratch ...
type Scratch struct {
	portPublicHTTP  int
	portPrivateHTTP int
}

// Option ...
type Option struct {
	PortPublicHTTP  int
	PortPrivateHTTP int
}

// New ...
func New(ctx context.Context, opt Option) *Scratch {
	span, _ := opentracing.StartSpanFromContext(ctx, "NewScratch")
	defer span.Finish()
	s := &Scratch{
		portPublicHTTP:  opt.PortPublicHTTP,
		portPrivateHTTP: opt.PortPrivateHTTP,
	}
	return s
}

// Run ...
func (s *Scratch) Run(ctx context.Context, handles *http.ServeMux) {
	span, _ := opentracing.StartSpanFromContext(ctx, "Scratch.Run")
	defer span.Finish()

	srvPublic := s.runPublicHTTP(handles)
	srvPrivate := s.runPrivateHTTP()

	idleSrvConnsClosed := graceful(srvPrivate, srvPublic)
	go func() {
		log.Info("http.Private start")
		if errSrvPrivateListenAndServe := srvPrivate.ListenAndServe(); errSrvPrivateListenAndServe != http.ErrServerClosed {
			log.Errorf("HTTP server ListenAndServe: %v", errSrvPrivateListenAndServe)
		}
	}()
	go func() {
		log.Info("http.Public start")
		if errSrvPublicListenAndServe := srvPublic.ListenAndServe(); errSrvPublicListenAndServe != http.ErrServerClosed {
			log.Errorf("HTTP server ListenAndServe: %v", errSrvPublicListenAndServe)
		}
	}()
	<-idleSrvConnsClosed
}
