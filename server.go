package main

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"sync"
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type User struct {
	ID       primitive.ObjectID `json:"id" bson:"id"`
	Name     string             `json:"name" bson:"name"`
	Email    string             `json:"email" bson:"email"`
	Password string             `json:"password" bson:"password"`
}

type usersHandler struct {
	sync.Mutex
	store map[string]User
}

func (h *usersHandler) userspost(response http.ResponseWriter, request *http.Request) {
	if request.Method == "POST" {
		response.Header().Add("content-type", "application/json")
		var user usersHandler
		_ = json.NewDecoder(request.Body).Decode(&user)
		collection := client.Database("instadb").Collection("users")
		ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		defer cancel()
		result, _ := collection.InsertOne(ctx, user)
		fmt.Println(json.NewEncoder(response).Encode(result))
	}
}

/*func (h *usersHandler) get(w http.ResponseWriter, r *http.Request) {

	users := make([]User, len(h.store))

	h.Lock()
	i := 0
	for _, use := range h.store {
		users[i] = use
		i++
	}

	h.Unlock()
	jsonBytes, err := json.Marshal(users)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(err.Error()))
		return
	}

	w.Header().Add("content-type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(jsonBytes)
}

func (h *usersHandler) getuser(w http.ResponseWriter, r *http.Request) {

	h.Lock()
	i := 0
	for _, use := range h.store {
		users[i] = use
		i++
	}

	h.Unlock()
	jsonBytes, err := json.Marshal(users)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(err.Error()))
		return
	}

	w.Header().Add("content-type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(jsonBytes)
}

func (h *usersHandler) post(w http.ResponseWriter, r *http.Request) {

	bodyBytes, err := ioutil.ReadAll(r.Body)
	defer r.Body.Close()
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(err.Error()))
		return
	}

	ct := r.Header.Get("content-type")
	if ct != "application/json" {
		w.WriteHeader(http.StatusUnsupportedMediaType)
		w.Write([]byte(fmt.Sprintf("need content-type 'application/json', but got '%s'", ct)))
		return
	}

	var users User
	json.Unmarshal(bodyBytes, &users)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte(err.Error()))
	}

	h.Lock()
	h.store[users.ID] = users
	defer h.Unlock()
}*/

func newUserHandler() *usersHandler {
	return &usersHandler{
		store: map[string]User{},
	}
}

var client *mongo.Client

func main() {

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	client, err := mongo.Connect(ctx, options.Client().ApplyURI("mongodb://localhost:27017"))

	usersHandler := newUserHandler()
	http.HandleFunc("/users", usersHandler.userspost)
	//http.HandleFunc("/users/", usersHandler.getuser)
	http.ListenAndServe(":8080", nil)

	defer func() {
		if err = client.Disconnect(ctx); err != nil {
			panic(err)
		}
	}()

}
