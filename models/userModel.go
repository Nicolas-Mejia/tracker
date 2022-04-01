package models

type User struct {
	// uso el mail como usuario
	User     string
	Password string
	Pedidos  *[]Pedido
	Mails    string
}
