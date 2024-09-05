package handlers

import (
	"log"
	"os"

	httpclient "github.com/ccrayz/sandbox-api/internal/http"

	"github.com/gin-gonic/gin"
)

var (
	L2_ENDPOINT = os.Getenv("ETH_L2_ENDPOINT")
	ERROR_PLUG  = false
)

type JsonResponse struct {
	JSONRPC string `json:"jsonrpc"`
	ID      int    `json:"id"`
	Result  string `json:"result"`
}

func HealthCheck(ctx *gin.Context) {
	ctx.JSON(200, gin.H{
		"status": "healthy",
	})
}

func ChangeErrorFlag(ctx *gin.Context) {
	if ERROR_PLUG {
		ERROR_PLUG = false
	} else {
		ERROR_PLUG = true
	}

	ctx.JSON(200, gin.H{
		"message": "flag changed",
	})
}

func GetLatesBlocks(ctx *gin.Context) {
	if ERROR_PLUG {
		ctx.JSON(500, gin.H{
			"status": "error",
		})
	} else {
		client := httpclient.NewClient(L2_ENDPOINT)

		req, err := client.NewRequest("POST", "", map[string]interface{}{
			"jsonrpc": "2.0",
			"method":  "eth_blockNumber",
			"params":  []interface{}{},
			"id":      1,
		})
		if err != nil {
			log.Fatalf("Failed to create request: %v", err)
		}

		var data JsonResponse
		_, err = client.Do(ctx, req, &data)
		if err != nil {
			log.Fatalf("Failed to send request: %v", err)
		}

		ctx.JSON(200, gin.H{
			"status": data.Result,
		})
	}
}
