package apiserver

import (
	"context"
	"fmt"
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/spf13/viper"
	"go.uber.org/fx"
	"os"
	"time"
)

func NewServer(cfg *viper.Viper) *gin.Engine {
	r := gin.Default()
	if !cfg.GetBool("app.debug") {
		gin.SetMode(gin.ReleaseMode)
		r = gin.New()
	}
	r.Use(cors.New(cors.Config{
		AllowMethods:     []string{"GET", "POST", "PUT", "PATCH", "DELETE", "HEAD", "OPTIONS"},
		AllowHeaders:     []string{"Origin", "Content-Length", "Content-Type", "Authorization"},
		AllowCredentials: true,
		AllowAllOrigins:  true,
		MaxAge:           12 * time.Hour,
	}))
	return r
}

func Start(lc fx.Lifecycle, g *gin.Engine) {
	lc.Append(fx.Hook{
		OnStart: func(ctx context.Context) error {
			var port = "8080"
			if os.Getenv("PORT") != "" {
				port = os.Getenv("PORT")
			}
			go g.Run(fmt.Sprintf("0.0.0.0:%s", port))
			return nil
		},
		OnStop: func(ctx context.Context) error {
			fmt.Println("bye :) ")
			return nil
		},
	})
}
