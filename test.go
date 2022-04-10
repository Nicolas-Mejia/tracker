package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	"github.com/gorilla/mux"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/mongo/readpref"
)

type task struct {
	ID      int    `json:"ID"`
	Name    string `json:"Name"`
	Content string `json:"Content"`
}

type allTasks []task

var tasks = allTasks{
	{
		ID:      1,
		Name:    "Task One",
		Content: "Some Content",
	},
}

func getTasks(w http.ResponseWriter, r *http.Request) {
	json.NewEncoder(w).Encode(tasks)
}

func test() {
	/*
		var password string = "test123"
		var bytesPass, err = bcrypt.GenerateFromPassword([]byte(password), 12)
		if err == nil{

			fmt.Println(bcrypt.CompareHashAndPassword(bytesPass, []byte("test123")))
			fmt.Println(bcrypt.CompareHashAndPassword(bytesPass, []byte("test321")))
		}
	*/

	// conexión mongodb
	var uri string = "mongodb+srv://nico:9pFtfrx7YnG4uJ3z@cluster0.aqqsk.mongodb.net/tracker?retryWrites=true&w=majority"

	client, err := mongo.Connect(context.TODO(), options.Client().ApplyURI(uri))

	if err != nil {
		panic(err)
	}
	defer func() {
		if err = client.Disconnect(context.TODO()); err != nil {
			panic(err)
		}
	}()

	// Ping the primary
	if err := client.Ping(context.TODO(), readpref.Primary()); err != nil {
		panic(err)
	}

	//var bytesPass []byte
	//bytesPass, err = bcrypt.GenerateFromPassword([]byte("test321"), 12)

	ctx := context.TODO()

	UsersColl := client.Database("tracker").Collection("users")
	//data := User{User:"nicomejia"}
	filter := User{Username: "papa", Password: "test123"}

	cursor1 := UsersColl.FindOne(ctx, filter)

	var user User
	cursor1.Decode(&user)

	if user.Username == "" {
		fmt.Println("asd")
	}

	cursor, err := UsersColl.Find(ctx, filter)

	var usersFil []User
	if err = cursor.All(ctx, &usersFil); err != nil {
		panic(err)
	}
	for _, item := range usersFil {
		fmt.Println(item.Username)
	}

	// inicializo router
	router := mux.NewRouter().StrictSlash(true)

	// endpoints
	router.HandleFunc("/users", userRegister)
	http.Handle("/", router)

	log.Fatal(http.ListenAndServe(":8080", router))

}
