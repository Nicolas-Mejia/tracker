package main

import (
	"context"
	"fmt"

	"github.com/Nicolas-Mejia/tracker/models"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/mongo/readpref"
)

func test(){
	/*
	fmt.Println("asd")
	var password string = "test123"
	var bytesPass, err = bcrypt.GenerateFromPassword([]byte(password), 12)
	if err == nil{

		fmt.Println(bcrypt.CompareHashAndPassword(bytesPass, []byte(password)))
		fmt.Println(bcrypt.CompareHashAndPassword(bytesPass, []byte("test321")))
	}
	*/

	// conexión mongodb
	var uri string = "mongodb+srv://user:pass@cluster0.aqqsk.mongodb.net/db?retryWrites=true&w=majority"

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

	coll := client.Database("tracker").Collection("users")
	data := models.User{User:"nicomejia", Password:"test123"}

	result, err := coll.InsertOne(context.TODO(), data)

	fmt.Printf("Inserted document with _id: %v\n", result.InsertedID)



/*
	// inicializo router
	router := mux.NewRouter()

	// endpoints
	router.HandleFunc("/users", userHandler)
	http.Handle("/", router)

	http.ListenAndServe(":8080", router)

*/



}