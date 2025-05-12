package socket

import (
	"encoding/json"
	"log"
	"net/http"
	"sync"

	"github.com/gorilla/websocket"
)

var (
	upgrader = websocket.Upgrader{
		CheckOrigin: func(r *http.Request) bool { return true },
	}

	users       []User
	usersonline = make(map[*websocket.Conn]string)
	mu          sync.Mutex
)

type User struct {
	Type string `json:"type"`
	Name string `json:"name"`
	Id   string `json:"id"`
}

type Message struct {
	Content  string `json:"message"`
	Receiver string `json:"receiver"`
	Sender   string `json:"sender"`
	Name     string `json:"sendername"`
	Seen     bool   `json:"seen"`
	Created  string `json:"created"`
}

// Upgrade connection to WebSocket
func Upgrade(w http.ResponseWriter, r *http.Request) *websocket.Conn {
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Println("Upgrade error:", err)
		return nil
	}
	return conn
}

// Handle WebSocket Connection
func HandleRequest(w http.ResponseWriter, r *http.Request) {
	conn := Upgrade(w, r)
	if conn == nil {
		http.Error(w, "WebSocket Upgrade Failed", http.StatusInternalServerError)
		return
	}
	defer conn.Close()

	for {
		// log.Println("Users found: ", users)
		var receiveinter map[string]interface{}
		err := conn.ReadJSON(&receiveinter)
		// log.Println("received: ", receiveinter)

		if err != nil {
			log.Println("❌ ReadJSON error:", err)
			mu.Lock()
			if id, exists := usersonline[conn]; exists {
				for i, user := range users {
					if user.Id == id {
						users[i] = users[len(users)-1]
						users = users[:len(users)-1]
						break
					}
				}
			} else {
				log.Println("⚠️ No user found for this connection in usersonline")
			}
			delete(usersonline, conn)
			mu.Unlock()

			broadcastUsers()
			return
		}

		// log.Println("✅ Received Through WebSocket:", receiveinter)

		// Handle based on "type" field in received JSON
		msgType, ok := receiveinter["type"].(string)
		if !ok {
			log.Println("⚠️ No valid 'type' field in message")
			continue
		}

		var user User

		switch msgType {
		case "user":
			name, nameOk := receiveinter["name"].(string)
			id, idOk := receiveinter["id"].(string)
			if !nameOk || !idOk {
				log.Println("⚠️ Missing 'name' or 'id' in user message")
				continue
			}
			user.Name = name
			user.Id = id
			user.Type = "user_list"
			handleUser(conn, user)
			// log.Println("👤 User connected:", user)

		case "new_post":
			data := struct {
				Type string                 `json:"type"`
				Data map[string]interface{} `json:"data"`
			}{
				Type: "new_post",
				Data: receiveinter,
			}

			for conn := range usersonline {
				if err := conn.WriteJSON(data); err != nil {
					log.Println("WriteJSON error:", err)
					mu.Lock()
					conn.Close()
					delete(usersonline, conn)
					mu.Unlock()
				}
			}
			// log.Println("📝 New post received:", receiveinter)

		case "reaction":
			data := struct {
				Type string                 `json:"type"`
				Data map[string]interface{} `json:"data"`
			}{
				Type: "reaction",
				Data: receiveinter,
			}

			for conn := range usersonline {
				if err := conn.WriteJSON(data); err != nil {
					log.Println("WriteJSON error:", err)
					mu.Lock()
					conn.Close()
					delete(usersonline, conn)
					mu.Unlock()
				}
			}
			// log.Println("Reactions: ", receiveinter)

		case "comment":
			data := struct {
				Type string                 `json:"type"`
				Data map[string]interface{} `json:"data"`
			}{
				Type: "comment",
				Data: receiveinter,
			}

			for conn := range usersonline {
				if err := conn.WriteJSON(data); err != nil {
					log.Println("WriteJSON error:", err)
					mu.Lock()
					conn.Close()
					delete(usersonline, conn)
					mu.Unlock()
				}
			}

		case "message":
			messageInterface, ok := receiveinter["data"]
			if !ok || messageInterface == nil {
				log.Println("No 'data' field found or it's nil")
				return
			}

			message, ok := messageInterface.(string)
			if !ok {
				log.Println("Expected 'data' to be a string but got:", messageInterface)
				return
			}

			var data Message
			if err := json.Unmarshal([]byte(message), &data); err != nil {
				log.Println("Unmarshal error:", err)
				return
			}
			
			sender := data.Sender
			receiver := data.Receiver

			rdata := struct {
				Type string      `json:"type"`
				Data interface{} `json:"data"`
			}{
				Type: "message",
				Data: data,
			}

			if sender != "" && receiver != "" {
				for conn := range usersonline {
					if usersonline[conn] == sender || usersonline[conn] == receiver {
						if err := conn.WriteJSON(rdata); err != nil {
							log.Println("WriteJSON error:", err)
							mu.Lock()
							conn.Close()
							delete(usersonline, conn)
							mu.Unlock()
						}
					}
				}
			}
			// log.Println("Message Received in the ws: ", data)

		case "typing":
			receiverid := receiveinter["receiverid"]
			senderid := receiveinter["senderid"].(string)
			sendername := receiveinter["sendername"].(string)
			data := struct {
				Type string      `json:"type"`
				SenderId string `json:"senderid"`
				SenderName string `json:"sendername"`
			}{
				Type: "typing",
				SenderId: senderid,
				SenderName: sendername,
			}

			for conn := range usersonline {
				if usersonline[conn] == receiverid {
					if err := conn.WriteJSON(data); err != nil {
						log.Println("WriteJSON error:", err)
						mu.Lock()
						conn.Close()
						delete(usersonline, conn)
						mu.Unlock()
					}
					break
				}
			}

		default:
			log.Println("⚠️ Unknown message type:", msgType)
		}
	}
}

func handleUser(conn *websocket.Conn, user User) {
	mu.Lock()
	if _, ok := usersonline[conn]; !ok  {
		usersonline[conn] = user.Id
		users = append(users, user)
	}
	mu.Unlock()
	broadcastUsers()
}

// Broadcast the updated users list to all connected clients
func broadcastUsers() {
	mu.Lock()
	defer mu.Unlock()
	data := struct {
		Type string `json:"type"`
		Data []User `json:"data"`
	}{
		Type: "user_list",
		Data: users,
	}
	for conn := range usersonline {
		if err := conn.WriteJSON(data); err != nil {
			log.Println("WriteJSON error:", err)
			conn.Close()
			delete(usersonline, conn)
		}
	}
}
