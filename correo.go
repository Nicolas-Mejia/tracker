package main

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"time"

	_ "github.com/joho/godotenv/autoload"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/mongo/readpref"
)

func scriptTrackingNational(id string) {
	mongoClient := getClient()
	ShipmentColl := mongoClient.Database("correo").Collection("shipment")

	//format the id with the uri env vars and get the URI ready
	uri := formatUri(id)

	response, err := http.Get(uri)

	if err != nil {
		fmt.Println("Error in the get request")
	}

	defer response.Body.Close()

	var tu TrackingUpdate
	if json.NewDecoder(response.Body).Decode(&tu); err != nil {
		fmt.Println("Error in the json decoder")
	}

	if tu.Data.Cantidad == 0 {
		fmt.Println("The package doesn't exist/wrong shipment ID.")
		return
	}

	shipment := getShipmentFromDB(id, ctx, ShipmentColl)
	//if I don't have that shipment ID in my DB:
	if shipment.TrackingID == "" {
		//create a new shipment
		shipment = Shipment{
			TrackingID: id,
			Quantity:   tu.Data.Cantidad,
		}

		//insert the new shipment in the DB
		insertResult, err := ShipmentColl.InsertOne(ctx, shipment)
		if err != nil {
			fmt.Println("Error in the insert")
		}

		fmt.Println("Inserted a single document: ", insertResult.InsertedID)

	} else {
		//if I do have that shipment ID in my DB:
		if shipment.Quantity != tu.Data.Cantidad {
			//if I have that shipment ID in my DB:
			//update the shipment in the DB
			filter := bson.M{"trackingID": id}
			update := bson.M{
				"$set": bson.M{
					"quantity": tu.Data.Cantidad,
				},
			}

			//change to inactive when the shipment is delivered

			updateResult, err := ShipmentColl.UpdateOne(ctx, filter, update)
			if err != nil {
				fmt.Println("Error in the update")
			}

			fmt.Printf("Matched %v documents and updated %v documents.\n", updateResult.MatchedCount, updateResult.ModifiedCount)

			//send update via whatsapp of the last event
			//index of the last event
			eventos := *tu.Data.Eventos
			req := makeWppMessageReq(eventos[0], id)

			client := http.Client{
				Timeout: 30 * time.Second,
			}

			_, err = client.Do(req)
			if err != nil {
				fmt.Printf("client: error making http request: %s\n", err)
				os.Exit(1)
			}

		} else {
			fmt.Println("The shipment is up to date.")
		}
	}

	//e := *tu.Data.Eventos
	//fmt.Println(e[0])

	return
}

func getShipmentFromDB(trackingID string, ctx context.Context, ShipmentColl *mongo.Collection) Shipment {
	var shipment Shipment

	filter := bson.M{"trackingID": trackingID}

	err := ShipmentColl.FindOne(ctx, filter).Decode(&shipment)

	if err != nil {
		if err.Error() == "mongo: no documents in result" {
			return Shipment{}
		} else {
			fmt.Println(err)
		}
	}

	return shipment
}

func formatUri(id string) string {
	l := len(id)
	prefix := id[:2]
	code := id[2 : l-2]
	uri1 := os.Getenv("correoUri1")
	uri2 := os.Getenv("correoUri2")
	uri3 := os.Getenv("correoUri3")

	uri := uri1 + prefix + uri2 + code + uri3

	return uri
}

func getClient() *mongo.Client {
	dbUser := os.Getenv("userTracker")
	dbPass := os.Getenv("passTracker")
	var uri string = fmt.Sprintf("mongodb+srv://%v:%v@cluster0.aqqsk.mongodb.net/?retryWrites=true&w=majority", dbUser, dbPass)

	client, err := mongo.Connect(ctx, options.Client().ApplyURI(uri))

	if err != nil {
		panic(err)
	}

	// Ping the primary
	if err := client.Ping(ctx, readpref.Primary()); err != nil {
		panic(err)
	}

	return client
}

func makeWppMessageReq(evento Evento, id string) *http.Request {
	//send message via whatsapp

	wppUri1 := os.Getenv("wppUri1")
	wppUri2 := os.Getenv("wppUri2")
	wppVersion := os.Getenv("wppVersion")
	wppSenderPhoneNumber := os.Getenv("wppSenderPhoneNumber")
	wppRecipientPhoneNumber := os.Getenv("wppRecipientPhoneNumber")
	wppAccessToken := os.Getenv("wppAccessToken")

	uri := wppUri1 + wppVersion + "/" + wppSenderPhoneNumber + wppUri2

	//create the message
	message := WppMessage{
		MessagingProduct: "whatsapp",
		RecipientType:    "individual",
		To:               wppRecipientPhoneNumber,
		Type:             "template",
		Template: Template{
			Name: "shipment_update",
			Language: Language{
				Code: "es_AR",
			},
			Components: []Component{
				{
					Type: "body",
					Parameters: []Parameter{
						{
							Type: "text",
							Text: id,
						},
						{
							Type: "text",
							Text: formatEventDetails(evento),
						},
					},
				},
			},
		},
	}

	//create the request
	b, err := json.Marshal(message)
	req, err := http.NewRequest("POST", uri, bytes.NewBuffer(b))
	if err != nil {
		fmt.Println("Error in the request")
	}

	//add the headers
	req.Header.Add("Content-Type", "application/json")
	req.Header.Add("Authorization", "Bearer "+wppAccessToken)

	return req

}

func formatEventDetails(evento Evento) string {
	var details string

	details = "Evento: " + evento.Name + "\\n" +
		"Fecha: " + evento.Date + "\\n" +
		"Planta: " + evento.Location + "\\n" +
		"Estado: " + evento.Status + "\\n" +
		"Origen: " + evento.Country + "\\n" +
		"Motivo no entrega: " + evento.NotDeliveredBecause
	return details
}

type Shipment struct {
	ID         primitive.ObjectID `json:"_id,omitempty" bson:"_id,omitempty"`
	TrackingID string             `json:"trackingID,omitempty" bson:"trackingID,omitempty"`
	Quantity   int                `json:"quantity,omitempty" bson:"quantity,omitempty"`
}

type WppMessage struct {
	MessagingProduct string   `json:"messaging_product" bson:"messaging_product"`
	RecipientType    string   `json:"recipient_type" bson:"recipient_type"`
	To               string   `json:"to" bson:"to"`
	Type             string   `json:"type" bson:"type"`
	Template         Template `json:"template" bson:"template"`
}

type Template struct {
	Name       string      `json:"name" bson:"name"`
	Language   Language    `json:"language" bson:"language"`
	Components []Component `json:"components" bson:"components"`
}

type Language struct {
	Code string `json:"code" bson:"code"`
}

type Component struct {
	Type       string      `json:"type" bson:"type"`
	Parameters []Parameter `json:"parameters" bson:"parameters"`
}

type Parameter struct {
	Type string `json:"type" bson:"type"`
	Text string `json:"text" bson:"text"`
}

//ES AL PEDO GUARDARME LOS EVENTOS, SACAR
