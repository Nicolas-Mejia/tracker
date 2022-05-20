package main

type Book struct {
	ISBN   int    `bson:"isbn,omitempty"`
	Title   string `bson:"title,omitempty"`
	Author string `bson:"author,omitempty"`
	Image  string `bson:"image,omitempty"`
}
