package metrics

import (
	"github.com/chenjiandongx/ginprom"
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

func NewServer(reg *prometheus.Registry) {
	server := gin.New()
	server.Use(cors.Default())
	server.GET("/metrics", ginprom.PromHandler(promhttp.HandlerFor(reg, promhttp.HandlerOpts{Registry: reg})))
	go server.Run("0.0.0.0:8585")
}
