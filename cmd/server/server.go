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

var serverName string

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
		r.GET("/block", handlers.GetLatesBlocks)

		addDynamicHandlers(r)

		if err := r.Run(); err != nil {
			fmt.Println("Failed to start server:", err)
		}
	},
}

// TODO: add opentelemetry middleware
func contextMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		ctx := context.WithValue(c.Request.Context(), "requestID", serverName)
		c.Request = c.Request.WithContext(ctx)
		c.Next()
	}
}

func addDynamicHandlers(r *gin.Engine) {
	log.Println("Add dynamic path =============================================")
	r.GET("/dynamic", func(c *gin.Context) {
		handlers.DynamicHandler(c, serverName)
	})
	r.GET("/call", func(c *gin.Context) {
		handlers.CallService(c, os.Getenv("CALL_SERVICE"), "GET", "dynamic", []interface{}{})
	})
	r.GET("/call-multiple", func(c *gin.Context) {
		handlers.CallMultipleServices(c, []string{os.Getenv("CALL_SERVICES")}, "GET", "dynamic", []interface{}{})
	})
}

func init() {
	log.Println("Show Server Config")
	serverName = os.Getenv("SERVER_NAME")
	if serverName == "" {
		serverName = "sandbox-api"
	}

	log.Println("- SERVER_NAME:", serverName)
	log.Println("- CALL_SERVICE:", os.Getenv("CALL_SERVICE"))
	log.Println("- CALL_SERVICE:", os.Getenv("CALL_SERVICES"))
	log.Println("- ETH_L2_ENDPOINT:", os.Getenv("ETH_L2_ENDPOINT"))
}
