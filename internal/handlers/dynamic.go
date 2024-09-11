package handlers

import (
	"fmt"
	"log"
	"net/http"

	httpclient "github.com/ccrayz/sandbox-api/internal/http"
	"github.com/gin-gonic/gin"
)

type DynamicResponse struct {
	Result string `json:"result"`
}

func CallMultipleServices(ctx *gin.Context, endpoints []string, method string, path string, params []interface{}) {

	ch := make(chan DynamicResponse, len(endpoints))
	defer close(ch)

	for _, endpoint := range endpoints {
		go func(endpoint string) {
			client := httpclient.NewClient(endpoint)

			req, err := client.NewRequest(method, path, params)
			if err != nil {
				ch <- DynamicResponse{
					Result: "error",
				}
			}

			var data DynamicResponse
			log.Println(req.URL)
			_, err = client.Do(ctx, req, &data)
			if err != nil {
				ch <- DynamicResponse{
					Result: "error",
				}
				return
			}

			ch <- data
		}(endpoint)
	}

	for response := range ch {
		fmt.Println(response)
	}

	ctx.JSON(http.StatusOK, gin.H{
		"status": "ok",
	})
}

func CallService(ctx *gin.Context, endpoint string, method string, path string, params []interface{}) {
	client := httpclient.NewClient(endpoint)

	req, err := client.NewRequest(method, path, params)
	if err != nil {
		ctx.JSON(500, gin.H{
			"status": "error",
		})
	}

	var data DynamicResponse
	log.Println(req.URL)
	_, err = client.Do(ctx, req, &data)
	if err != nil {
		ctx.JSON(500, gin.H{
			"status": "error",
		})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"status": data.Result,
	})
}

func DynamicHandler(ctx *gin.Context, apiName string) {
	response := DynamicResponse{
		Result: fmt.Sprintf("I am a %s", apiName),
	}

	ctx.JSON(http.StatusOK, response)
}
