package main

import (
	"context"

	"go.mongodb.org/mongo-driver/mongo"
)

var client *mongo.Client
var ctx context.Context
var UsersColl *mongo.Collection
var OrdersColl *mongo.Collection

func main() {
	/*

		//dbUser := os.Getenv("userTracker")
		//dbPass := os.Getenv("passTracker")
		//var uri string = fmt.Sprintf("mongodb+srv://%v:%v@cluster0.aqqsk.mongodb.net/?retryWrites=true&w=majority", dbUser, dbPass)
		var uri string = "mongodb+srv://nico:bnHUj392nOGMaVg3@cluster0.aqqsk.mongodb.net/?retryWrites=true&w=majority"


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

		ctx = context.TODO()

		UsersColl = client.Database("tracker").Collection("users")
		OrdersColl = client.Database("tracker").Collection("orders")

		// inicializo router
		router := mux.NewRouter().StrictSlash(true)

		// endpoints
		router.HandleFunc("/users/register", userRegister).Methods("POST")
		router.HandleFunc("/users/login", userLogin).Methods("POST")
		router.HandleFunc("/orders", createOrder).Methods("POST")
		router.HandleFunc("/orders", getUserOrders).Methods("GET")
		router.HandleFunc("/orders/{id}", getOrderDetails).Methods("GET")
		router.HandleFunc("/orders/{id}", updateOrder).Methods("PUT")
		router.HandleFunc("/orders/{id}", deleteOrder).Methods("DELETE")
		http.Handle("/", router)

		log.Fatal(http.ListenAndServe(":8080", router))



		//getBookData("atomic habits")
	*/

	/*
		var sId string

		fmt.Print("Type a Shipment ID: ")
		fmt.Scan(&sId)
	*/

	scriptTrackingNational("RE102719588AR")
	//"RE102719588AR"
}
