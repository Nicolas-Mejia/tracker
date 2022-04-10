package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"

	"go.mongodb.org/mongo-driver/mongo"
	"golang.org/x/crypto/bcrypt"
)

func userRegister(res http.ResponseWriter, req *http.Request) {
	// inicializo variable para guardar los datos decodificados
	var newUser User

	// si hay un error decodificando la req devuelvo status 500
	if err := json.NewDecoder(req.Body).Decode(&newUser); err != nil {
		res.WriteHeader(http.StatusInternalServerError)
		res.Write([]byte(`{ "message": "` + err.Error() + `" }`))
		return
	}

	if newUser.Username == "" {
		res.WriteHeader(http.StatusBadRequest)
		res.Write([]byte(`{ "message": "No username key." }`))
		return
	}

	filter := User{Username: newUser.Username}

	var user2Check User
	cursor1 := UsersColl.FindOne(ctx, filter)

	if err := cursor1.Decode(&user2Check); err != nil {
		// verifico que el error sea que no hay documentos, es decir, el usuario no existe aun.
		if !errors.Is(err, mongo.ErrNoDocuments) {
			panic(err)
		}
	}

	// el usuario ya existe, devuelvo status
	if user2Check.Username != "" {
		res.WriteHeader(http.StatusConflict)
		res.Write([]byte(`{ "message": "There is an account registered to that email." }`))
		return
	}

	if len(newUser.Password) < 8 {
		res.WriteHeader(http.StatusForbidden)
		res.Write([]byte(`{ "message": "The password is too short. It should be at least 8 characters long." }`))
		return
	}

	// todo ok, registro el usuario

	hashedPass, err := bcrypt.GenerateFromPassword([]byte(newUser.Password), 12)

	newUser.Password = string(hashedPass)

	insertResult, err := UsersColl.InsertOne(ctx, newUser)
	if err != nil {
		fmt.Println("out 2")
		panic(err)
	}

	res.WriteHeader(http.StatusOK)
	res.Write([]byte(`{ "message": "The user was created successfully." }`))

	fmt.Println("DSADSADSADSSAD", insertResult.InsertedID)

}

func userLogin(res http.ResponseWriter, req *http.Request) {
	var user User

	// si hay un error decodificando la req devuelvo status 500
	if err := json.NewDecoder(req.Body).Decode(&user); err != nil {
		res.WriteHeader(http.StatusInternalServerError)
		res.Write([]byte(`{ "message": "` + err.Error() + `" }`))
		return
	}

	// creo un usuario con el nombre de usuario que me ingresaron, para usarlo en la consulta
	filter := User{Username: user.Username}

	var user2Compare User
	cursor1 := UsersColl.FindOne(ctx, filter)

	if err := cursor1.Decode(&user2Compare); err != nil {
		// verifico si el error es que no hay documentos, es decir, el usuario no existe.
		if errors.Is(err, mongo.ErrNoDocuments) {
			res.WriteHeader(http.StatusBadRequest)
			res.Write([]byte(`{ "message": "There isn't a user registered to that email." }`))
			return
		} else {
			panic(err)
		}
	}

	if err := bcrypt.CompareHashAndPassword([]byte(user2Compare.Password), []byte(user.Password)); err != nil {
		res.WriteHeader(http.StatusUnauthorized)
		res.Write([]byte(`{ "message": "The password is wrong." }`))
		return
	}

	res.WriteHeader(http.StatusOK)

	// aca meto el jwt y devuelvo el token

	userJson, err := json.Marshal(user2Compare)
	if err != nil {
		panic(err)
	}
	res.Write(userJson)

}
