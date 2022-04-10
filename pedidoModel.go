package main

type Pedido struct {
	Id          int     `bson:"Id"`
	Titulo      string  `bson:"Titulo"` // opcional
	User        *User   `bson:"User"`
	TrackingNac string  `bson:"TrackingNac"` // opcional
	TrackingInt string  `bson:"TrackingInt"` // opcional
	Libros      []Libro `bson:"Libros"`      // opcional

}

type allPedidos []Pedido
