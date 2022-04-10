package main

type Libro struct {
	ISBN   int    `bson:"ISBN,omitempty"`
	Nombre string `bson:"nombre,omitempty"`
	Autor  string `bson:"Autor,omitempty"`
	Imagen string `bson:"Imagen,omitempty"`
}
