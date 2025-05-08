package handlers

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"real_time_forum/backend/db/controllers"
)

type Message struct {
	Content  string `json:"message"`
	Receiver string `json:"receiver"`
	Sender   string `json:"sender"`
	Name 	string `json:"sendername"`
	Seen     bool   `json:"seen"`
	Created  string `json:"created"`
}


func HandlePrivateMessage(w http.ResponseWriter, r *http.Request) {
	var (
		newMessage    Message
		ok            bool
		err           error
		messageToSave string
		threadId      string
		messageSlice  []Message
	)

	req, err := io.ReadAll(r.Body)
	if err != nil {
		http.Error(w, "Failed to read request body", http.StatusBadRequest)
		log.Println("Error reading request body:", err)
		return
	}
	r.Body.Close()

	if err = json.Unmarshal(req, &newMessage); err != nil {
		http.Error(w, "Invalid JSON", http.StatusBadRequest)
		log.Println("JSON Unmarshal error:", err)
		return
	}
	// log.Println("Received Private Message: ", newMessage)

	if threadId, ok, err = controllers.IsThread(newMessage.Sender, newMessage.Receiver); err != nil {
		http.Error(w, "Database error", http.StatusInternalServerError)
		log.Println("IsThread error:", err)
		return
	}

	if !ok {

		newThreadId, err := controllers.CreateNewMessageTable(newMessage.Sender, newMessage.Receiver)

		if err != nil {
			log.Println(err.Error())
			return
		}

		messageToSave = fmt.Sprintf(`
		[{"message": "%s", "sender": "%s", "receiver": "%s","name" : "%s", "created": "%s", "seen": %t}]
		`, newMessage.Content, newMessage.Sender, newMessage.Receiver, newMessage.Name, newMessage.Created, newMessage.Seen)

		if err = controllers.CreateNewThread(messageToSave, newThreadId); err != nil {
			log.Println("Thread Does Not Exist: ", err)
		}

		res, _ := json.Marshal(newMessage)
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(map[string]string{"message" : "Messages updated successfully", "data" : string(res) })

	} else {
		messages, err := controllers.GetMessages(threadId)
		if err != nil {
			log.Println(err.Error())
			return
		}

		json.Unmarshal([]byte(messages), &messageSlice)

		messageSlice = append(messageSlice, newMessage)

		updatedMessages, _ := json.Marshal(messageSlice)

		controllers.UpdateMessageThread(threadId, string(updatedMessages))

		res, _ := json.Marshal(newMessage)
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(map[string]string{"message" : "Messages updated successfully", "data" : string(res)})

	}
}

func GetUserMessages(w http.ResponseWriter, r *http.Request) {
	userId := r.URL.Query().Get("userId")

	if userId == "" {
		return
	}

	threadIds, err := controllers.GetUserThreadIds(userId)

	if err != nil {
		log.Println(err)
		return
	}

	var userMessages []string

	for _, threadId := range threadIds {
		message, err := controllers.GetMessages(threadId)
		if err != nil {
			log.Println(err)
			return
		}

		userMessages = append(userMessages, message)
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]interface{}{"data" : userMessages})
}