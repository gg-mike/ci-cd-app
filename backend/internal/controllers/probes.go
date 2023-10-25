package controllers

import (
	"net/http"
	"sync/atomic"

	"github.com/gin-gonic/gin"
)

type ReadyProbe struct {
  state *atomic.Bool
}

func (rp *ReadyProbe) Init() {
  rp.state = &atomic.Bool{}
  rp.state.Store(false)
}

func (rp *ReadyProbe) Ready() {
  rp.state.Store(true)
}

func (rp *ReadyProbe) Readyz(ctx *gin.Context) {
  if rp.state == nil || !rp.state.Load() {
    ctx.Status(http.StatusServiceUnavailable)
    return
  }
  ctx.Status(http.StatusOK)
}

func Healthz(ctx *gin.Context) {
  ctx.Status(http.StatusOK)
}
