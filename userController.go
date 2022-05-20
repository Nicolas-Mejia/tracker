package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"golang.org/x/crypto/bcrypt"
)

func userRegister(res http.ResponseWriter, req *http.Request) {
	// inicializo variable para guardar los datos decodificados
	var newUser User

	// si hay un error decodificando la req devuelvo status 500
	if err := json.NewDecoder(req.Body).Decode(&newUser); err != nil {
		writeInternalServerError(res, err)
		return
	}

	if newUser.Username == "" {
		writeStatusBadRequest(res, "No username key.")
		return
	}

	filter := bson.D{primitive.E{Key: "username", Value: newUser.Username}}

	var user2Check User
	cursor1 := UsersColl.FindOne(ctx, filter)

	if err := cursor1.Decode(&user2Check); err != nil {
		// verifico que el error sea que no hay documentos, es decir, el usuario no existe aun.
		if !errors.Is(err, mongo.ErrNoDocuments) {
			writeInternalServerError(res, err)
			return
		}
	} else {
		// if no ErrNoDocuments, then there already is a user with that username
		writeStatusConflict(res, "There is an account registered to that email.")
		return
	}

	if len(newUser.Password) < 8 {
		writeStatusForbidden(res, "The password is too short. It should be at least 8 characters long.")
		return
	}

	// todo ok, registro el usuario

	hashedPass, err := bcrypt.GenerateFromPassword([]byte(newUser.Password), 12)

	newUser.Id = primitive.NewObjectIDFromTimestamp(time.Now())
	newUser.Password = string(hashedPass)

	insertResult, err := UsersColl.InsertOne(ctx, newUser)
	if err != nil {
		writeInternalServerError(res, err)
		return
	}

	writeStatusOk(res, []byte(`{ "message": "The user was created successfully." }`))

	fmt.Println("Inserted: ", insertResult.InsertedID)
}

func userLogin(res http.ResponseWriter, req *http.Request) {
	var user User

	// si hay un error decodificando la req devuelvo status 500
	if err := json.NewDecoder(req.Body).Decode(&user); err != nil {
		writeInternalServerError(res, err)
		return
	}

	// creo un usuario con el nombre de usuario que me ingresaron, para usarlo en la consulta
	filter := bson.D{primitive.E{Key: "username", Value: user.Username}}

	var user2Compare User
	cursor1 := UsersColl.FindOne(ctx, filter)

	if err := cursor1.Decode(&user2Compare); err != nil {
		// check if the error is ErrNoDocuments, that means there isnt a user with that username.
		if errors.Is(err, mongo.ErrNoDocuments) {
			res.WriteHeader(http.StatusBadRequest)
			res.Write([]byte(`{ "message": "There isn't a user registered to that email." }`))
			return
		} else {
			writeInternalServerError(res, err)
			return
		}
	}

	fmt.Println(user2Compare.Id)

	if err := bcrypt.CompareHashAndPassword([]byte(user2Compare.Password), []byte(user.Password)); err != nil {
		res.WriteHeader(http.StatusUnauthorized)
		res.Write([]byte(`{ "message": "The password is wrong." }`))
		return
	}
	user2Compare.Password = ""

	// aca meto el jwt y devuelvo el token

	userJson, err := json.Marshal(user2Compare)
	if err != nil {
		writeInternalServerError(res, err)
		return
	}

	writeStatusOk(res, userJson)
}
