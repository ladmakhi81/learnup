package websocket

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
	"github.com/ladmakhi81/learnup/pkg/contracts"
	"net/http"
)

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
}

func Handler(manager *WsManager, tokenSvc contracts.Token) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		claim, err := tokenSvc.DecodeToken(ctx.GetHeader("authorization"))
		if err != nil {
			ctx.JSON(http.StatusInternalServerError, gin.H{})
			return
		}
		if claim == nil {
			ctx.JSON(http.StatusUnauthorized, gin.H{})
			return
		}
		conn, err := upgrader.Upgrade(ctx.Writer, ctx.Request, nil)
		if err != nil {
			ctx.JSON(http.StatusInternalServerError, gin.H{})
			return
		}
		manager.AddConnection(claim.UserID, conn)
		defer func() {
			manager.RemoveConnection(1)
			conn.Close()
		}()

		for {
			var message WsMessagePayload
			if err := conn.ReadJSON(&message); err != nil {
				fmt.Println(err)
				return
			}
			fmt.Printf("received : %+v\n", message)
		}
	}
}
