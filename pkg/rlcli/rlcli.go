package rlcli

import (
	"context"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/sanguohot/rlcli/pkg/common/log"
	"github.com/ulule/limiter/v3"
	mgin "github.com/ulule/limiter/v3/drivers/middleware/gin"
	"github.com/ulule/limiter/v3/drivers/store/memory"
	"net/http"
	"os"
	"os/signal"
	"time"
)

type Rlcli struct {
	limit string
	addr string
}

func New(limit, addr string) *Rlcli {
	return &Rlcli{
		addr: addr,
		limit: limit,
	}
}

func (s *Rlcli) defaultHandler(c *gin.Context) {
	type message struct {
		Message string `json:"message"`
		Timestamp int64 `json:"timestamp"`
	}
	resp := message{Message: "ok", Timestamp: time.Now().Unix()}
	c.JSON(http.StatusOK, resp)
}

func (s *Rlcli) rateLimitMiddleware() gin.HandlerFunc {
	rate, err := limiter.NewRateFromFormatted(s.limit)
	if err != nil {
		log.Logger.Fatal(err.Error())
	}
	store := memory.NewStore()
	return mgin.NewMiddleware(limiter.New(store, rate))
}

func (s *Rlcli) startServer() {
	gin.SetMode(gin.DebugMode)
	r := gin.New()
	r.Use(gin.Recovery())
	// 默认设置logger，但启用logger会导致吞吐量大幅度降低
	if os.Getenv("GIN_LOG") != "off" {
		r.Use(gin.Logger())
	}
	r.MaxMultipartMemory = 10 << 20 // 10 MB
	r.NoRoute(s.rateLimitMiddleware(), s.defaultHandler)
	//r.Any("/", s.rateLimitMiddleware(), s.defaultHandler)
	server := &http.Server{
		Addr:           fmt.Sprintf("%s", s.addr),
		Handler:        r,
		ReadTimeout:    10 * time.Second,
		WriteTimeout:   10 * time.Second,
		MaxHeaderBytes: 1 << 20, // 1 MB
	}
	go func() {
		if err := server.ListenAndServe(); err != nil {
			log.Logger.Fatal(err.Error())
		}
	}()
	log.Sugar.Infof("[http] listening => %s, limit => %v", server.Addr, s.limit)
	// apiserver发生错误后延时五秒钟，优雅关闭
	quit := make(chan os.Signal)
	signal.Notify(quit, os.Interrupt)
	<-quit
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	if err := server.Shutdown(ctx); err != nil {
		log.Logger.Fatal(err.Error())
	}
	log.Sugar.Infof("stop server => %s, limit => %v", server.Addr, s.limit)
}

func (s *Rlcli) Serve() {
	s.startServer()
}
