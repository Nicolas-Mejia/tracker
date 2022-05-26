package main

import (
	"encoding/json"
	"errors"
	"net/http"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

func writeStatusOk(res http.ResponseWriter, orderJson []byte) {
	res.WriteHeader(http.StatusOK)
	res.Write(orderJson)
	res.Header().Set("Content-Type", "application/json")
	return
}

func writeStatusBadRequest(res http.ResponseWriter, message string) {
	res.WriteHeader(http.StatusBadRequest)
	res.Write([]byte(`{ "message": "` + message + `" }`))
}

func writeStatusForbidden(res http.ResponseWriter, message string) {
	res.WriteHeader(http.StatusForbidden)
	res.Write([]byte(`{ "message": "` + message + `" }`))
}

func writeStatusNotFound(res http.ResponseWriter, message string) {
	res.WriteHeader(http.StatusNotFound)
	res.Write([]byte(`{ "message": "` + message + `" }`))
}

func writeStatusConflict(res http.ResponseWriter, message string) {
	res.WriteHeader(http.StatusConflict)
	res.Write([]byte(`{ "message": "` + message + `" }`))
}

func writeInternalServerError(res http.ResponseWriter, err error) {
	res.WriteHeader(http.StatusInternalServerError)
	res.Write([]byte(`{ "message": "` + err.Error() + `" }`))
	return
}

func orderDecoder(res http.ResponseWriter, req *http.Request) (bool, Order) {
	var newOrder Order

	if err := json.NewDecoder(req.Body).Decode(&newOrder); err != nil {
		writeInternalServerError(res, err)
		return false, newOrder
	}

	return true, newOrder

}

func orderCheckNotExisting(res http.ResponseWriter, filter bson.D) bool {
	var order2Compare Order
	cursor := OrdersColl.FindOne(ctx, filter)

	if err := cursor.Decode(&order2Compare); err != nil {
		// check if the error is ErrNoDocuments, that means there isnt an order with that tracking number.
		if !errors.Is(err, mongo.ErrNoDocuments) {
			writeInternalServerError(res, err)
			return false
		}
	} else {
		// if no ErrNoDocuments, then the call did return a user, so that username is already in the system.
		writeStatusConflict(res, "There is an order created for that tracking id already.")
		return false
	}

	return true
}

func orderFindOneAndDecode(res http.ResponseWriter, filter bson.D) (bool, Order) {
	var orderResult Order
	cursor := OrdersColl.FindOne(ctx, filter)

	if err := cursor.Decode(&orderResult); err != nil {
		// check if the error is ErrNoDocuments, that means there isnt an order with that tracking number.
		if errors.Is(err, mongo.ErrNoDocuments) {
			// if err == ErrNoDocuments, then we couldn't find an order with that ID.
			writeStatusConflict(res, "There are no orders with that ID.")
			return false, orderResult
		} else {
			writeInternalServerError(res, err)
			return false, orderResult
		}
	}

	return true, orderResult
}

func orderChecks(order Order, res http.ResponseWriter) bool {
	// pedido nacional o internacional.
	ok := true
	written := false

	if order.TrackingNac == "" && order.TrackingInt == "" {
		res.WriteHeader(http.StatusBadRequest)
		res.Write([]byte(`{ "message": "No tracking added to the request." }`))
		ok = false
		written = true
	}

	if order.IsInt && order.TrackingNac != "" && written == false {
		res.WriteHeader(http.StatusBadRequest)
		res.Write([]byte(`{ "message": "An international order shouldn't have a national tracking id." }`))
		ok = false
		written = true
	}

	if !order.IsInt && order.TrackingInt != "" && written == false {
		res.WriteHeader(http.StatusBadRequest)
		res.Write([]byte(`{ "message": "A national order shouldn't have an international tracking id." }`))
		ok = false
		written = true
	}

	if order.User == primitive.NilObjectID && written == false {
		res.WriteHeader(http.StatusBadRequest)
		res.Write([]byte(`{ "message": "An order needs to have a user associated." }`))
		ok = false
		written = true
	}

	return ok
}
