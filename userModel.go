package main

type User struct {
	// uso el mail como usuario
	Username string    `bson:"username,omitempty"`
	Password string    `bson:"password,omitempty"`
	Pedidos  *[]Pedido `bson:"pedidos,omitempty"`
	Mails    string    `bson:"mails,omitempty"`
}
