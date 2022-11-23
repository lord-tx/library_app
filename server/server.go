package main

import (
	"fmt"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
)

type Message struct {
	SenderID   uint   `json:"s_id"`
	Message    string `json:"message"`
	ReceiverID uint   `json:"r_id"`
	GroupID    uint   `json:"g_id"`
	Reply      bool   `json:"reply"`
	ReplyID    uint   `json:"rep_id"`
}

var wsupgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
}

var messageStore []Message
var newMessage Message
var connections map[int]*websocket.Conn

func wshandler(w http.ResponseWriter, r *http.Request, user_id int) {

	connections = make(map[int]*websocket.Conn)

	conn, err := wsupgrader.Upgrade(w, r, nil)
	if err != nil {
		fmt.Println("Failed to set websocket upgrade: ", err)
		return
	}

	connections[user_id] = conn
	// fmt.Print("The Connection created: ")
	// fmt.Println(connections[user_id])

	defer conn.Close()

	var data Message

	for {
		err := conn.ReadJSON(&data)
		if err != nil {
			break
		}

		/// Add new message to message store
		messageStore = append(messageStore, data)

		newMessage = data

		// fmt.Print("The Connection to be Called: ")
		// fmt.Println(data.ReceiverID)
		// fmt.Println(connections[int(data.ReceiverID)])

		/// Reply specific connection on its channel
		connection := connections[int(data.ReceiverID)]
		connection.WriteJSON(newMessage)

		// conn.WriteJSON(messageStore)
	}
}

func main() {

	r := gin.Default()

	r.LoadHTMLFiles("index.html")

	r.GET("/", func(c *gin.Context) {
		c.HTML(200, "index.html", nil)
	})

	r.GET("/ws/:user_id", func(c *gin.Context) {
		val, err := strconv.Atoi(c.Param("user_id"))

		if err != nil {
			return
		}
		wshandler(c.Writer, c.Request, val)
	})

	r.Run("localhost:8080")
}
