package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"time"

	"github.com/gorilla/mux"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

func createOrder(res http.ResponseWriter, req *http.Request) {
	ok, newOrder := orderDecoder(res, req)
	if !ok {
		return
	}

	if ok = orderChecks(newOrder, res); !ok {
		return
	}

	//check if an order with that tracking id was already created
	var filter bson.D
	if newOrder.IsInt {
		filter = bson.D{primitive.E{Key: "trackingInt", Value: newOrder.TrackingInt}}
	} else {
		filter = bson.D{primitive.E{Key: "trackingNac", Value: newOrder.TrackingNac}}
	}

	ok = orderCheckNotExisting(res, filter)
	if !ok {
		return
	}

	newOrder.Id = primitive.NewObjectIDFromTimestamp(time.Now())

	insertResult, err := OrdersColl.InsertOne(ctx, newOrder)
	if err != nil {
		writeInternalServerError(res, err)
	}

	orderJson, err := json.Marshal(newOrder)
	if err != nil {
		writeInternalServerError(res, err)
		return
	}

	writeStatusOk(res, orderJson)

	fmt.Println("Inserted: ", insertResult.InsertedID)
}

func updateOrder(res http.ResponseWriter, req *http.Request) {
	vars := mux.Vars(req)
	id, ok := vars["id"]
	if !ok {
		writeStatusConflict(res, "No id sent in the URL params.")
		return
	}

	ok, order := orderDecoder(res, req)
	if !ok {
		return
	}
	if ok = orderChecks(order, res); !ok {
		return
	}

	orderId, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		writeInternalServerError(res, err)
		return
	} else {
		order.Id = orderId
	}

	filter := bson.D{primitive.E{Key: "_id", Value: order.Id}}

	var result mongo.UpdateResult
	if result, err := OrdersColl.ReplaceOne(ctx, filter, order); err != nil {
		writeInternalServerError(res, err)
		return
	} else if result.ModifiedCount == 0 {
		writeStatusConflict(res, "We couldn't find an order with that ID.")
		return
	}

	writeStatusOk(res, []byte(`{ "message": "Order updated successfully." }`))

	fmt.Println("Updated: ", result.UpsertedID)
}

func getOrderDetails(res http.ResponseWriter, req *http.Request) {
	vars := mux.Vars(req)
	id, ok := vars["id"]
	if !ok {
		writeStatusConflict(res, "No id sent in the URL params.")
		return
	}
	orderId, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		writeInternalServerError(res, err)
		return
	}

	filter := bson.D{primitive.E{Key: "_id", Value: orderId}}

	ok, orderResult := orderFindOneAndDecode(res, filter)
	if !ok {
		return
	}

	if orderJson, err := json.Marshal(orderResult); err != nil {
		writeInternalServerError(res, err)
		return
	} else {
		writeStatusOk(res, orderJson)
	}
}

func getUserOrders(res http.ResponseWriter, req *http.Request) {
	var user User
	fmt.Println("asd",req.Body)

	err := json.NewDecoder(req.Body).Decode(&user)
	if err != nil {

		writeInternalServerError(res, err)
		return
	}

	fmt.Println(user)

	filter := bson.D{primitive.E{Key: "user", Value: user.Id}}

	ok, orderResult := orderFindAllAndDecode(res, filter)
	if !ok {
		return
	}

	if orderJson, err := json.Marshal(orderResult); err != nil {
		writeInternalServerError(res, err)
		return
	} else {
		writeStatusOk(res, orderJson)
	}
}

func deleteOrder(res http.ResponseWriter, req *http.Request) {
	vars := mux.Vars(req)
	id, ok := vars["id"]
	if !ok {
		writeStatusConflict(res, "No id sent in the URL params.")
		return
	}

	ok, order := orderDecoder(res, req)
	if !ok {
		return
	}

	orderId, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		writeInternalServerError(res, err)
		return
	} else {
		order.Id = orderId
	}

	filter := bson.D{primitive.E{Key: "_id", Value: order.Id}}

	var result mongo.DeleteResult
	if result, err := OrdersColl.DeleteOne(ctx, filter); err != nil {
		writeInternalServerError(res, err)
		return
	} else if result.DeletedCount == 0 {
		writeStatusConflict(res, "We couldn't find an order with that ID.")
		return
	}

	writeStatusOk(res, []byte(`{ "message": "Order deleted successfully." }`))

	fmt.Println("Deleted: ", result)

}

func getPackageHistory(res http.ResponseWriter, trackingId string) (bool, []Evento) {
	productCode := trackingId[0:2]
	idShipping := trackingId[2 : len(trackingId)-2]
	uri1 := os.Getenv("uriTracking1")
	uri2 := os.Getenv("uriTracking2")
	uri3 := os.Getenv("uriTracking3")
	uri := uri1 + productCode + uri2 + idShipping + uri3

	response, err := http.Get(uri)

	if err != nil {
		writeInternalServerError(res, err)
		return false, nil
	}

	var tu TrackingUpdate
	if json.NewDecoder(response.Body).Decode(&tu); err != nil {
		writeInternalServerError(res, err)
		return false, nil
	}

	if tu.Data.Cantidad == 0 {
		writeStatusNotFound(res, "El paquete no tiene movimientos aún")
		return false, nil
	}

	e := *tu.Data.Eventos

	return true, e
}

func getLatestPackageHistoryEvent(res http.ResponseWriter, pId string) (bool, Evento) {
	ok, p := getPackageHistory(res, pId)
	if !ok {
		return false, Evento{}
	} else {
		return true, p[0]
	}
}

func getUserPackages(res http.ResponseWriter, req *http.Request, pageNumber int, userId primitive.ObjectID, onlyActive bool) {

}
