package server

import (
	_ "net/http/pprof"

	"github.com/go-kratos/kratos/v2/transport"
	"github.com/go-kratos/kratos/v2/transport/http"

	"github.com/tencat-dev/go-base/internal/conf"
)

type PprofServer transport.Server

// NewPprofServer new a Pprof server.
func NewPprofServer(c *conf.PprofServer) PprofServer {
	var opts []http.ServerOption
	if c.Addr != "" {
		opts = append(opts, http.Address(c.Addr))
	}
	srv := http.NewServer(opts...)
	return srv
}
