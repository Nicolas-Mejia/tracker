package main

import "go.mongodb.org/mongo-driver/bson/primitive"

type Order struct {
	Id          primitive.ObjectID `bson:"_id" json:"_id,omitempty"`
	User        primitive.ObjectID `bson:"user,omitempty"`
	TrackingNac string             `bson:"trackingNac,omitempty"` // opcional
	TrackingInt string             `bson:"trackingInt,omitempty"` // opcional
	//	Books       []Book `bson:"books,omitempty"`      later
	IsInt    bool      `bson:"isInt,omitempty"` // opcional
	isActive bool      `bson:"isActive,omitempty"`
	History  *[]Evento `bson:"history,omitempty"`
}

type TrackingUpdate struct {
	Id   primitive.ObjectID `bson:"_id" json:"_id,omitempty"`
	Rta  string             `bson:"rta,omitempty"`
	Code int                `bson:"code,omitempty"`
	Data *PackageData       `bson:"data" json:"data"`
}

type Evento struct {
	Name                string `bson:"codigoEvento" json:"codigoEvento,omitempty"`
	Date                string `bson:"fechaEvento" json:"fechaEvento"`
	Location            string `bson:"planta" json:"planta"`
	Status              string `bson:"estadoEntrega" json:"estadoEntrega"`
	NotDeliveredBecause string `bson:"motivoNoEntrega" json:"motivoNoEntrega"`
	Country             string `bson:"nombrePais" json:"nombrePais"`
}

type PackageData struct {
	Eventos        *[]Evento `bson:"eventos" json:"eventos"`
	Id             int       `bson:"id" json:"id"`
	CodigoProducto string    `bson:"codigoProducto" json:"codigoProducto"`
	CodigoPais     string    `bson:"codigoPais" json:"codigoPais"`
	Cantidad       int       `bson:"cantidad" json:"cantidad"`
}

type UserPackages struct {
	Packages       PackageData `bson:"eventos" json:"eventos"`
	Id             int         `bson:"id" json:"id"`
	CodigoProducto string      `bson:"codigoProducto" json:"codigoProducto"`
	CodigoPais     string      `bson:"codigoPais" json:"codigoPais"`
	Cantidad       int         `bson:"cantidad" json:"cantidad"`
}
