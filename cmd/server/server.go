package server

import (
	"context"
	"fmt"
	"log"
	"os"

	"github.com/gin-gonic/gin"
	"github.com/spf13/cobra"

	"github.com/ccrayz/sandbox-api/internal/handlers"

	"github.com/prometheus/client_golang/prometheus/promhttp"
)

var ServerCmd = &cobra.Command{
	Use:   "server",
	Short: "Server is a CLI application",
	Run: func(cmd *cobra.Command, args []string) {
		r := gin.Default()
		r.Use(contextMiddleware())
		r.Use(prometheusMiddleware())

		r.GET("/metrics", gin.WrapH(promhttp.Handler()))
		r.GET("/health", handlers.HealthCheck)
		r.GET("/error", handlers.ChangeErrorFlag)

		r.POST("/block", handlers.GetLatesBlocks)
		if err := r.Run(); err != nil {
			fmt.Println("Failed to start server:", err)
		}
	},
}

// TODO: add opentelemetry middleware
func contextMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		ctx := context.WithValue(c.Request.Context(), "requestID", "sandbox-api")
		c.Request = c.Request.WithContext(ctx)
		c.Next()
	}
}

func init() {
	log.Println("Show Server Config")
	log.Println("- ETH_L2_ENDPOINT:", os.Getenv("ETH_L2_ENDPOINT"))
}
