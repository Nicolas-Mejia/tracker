package models

type Pedido struct {
	Id                    int
	Titulo                string // opcional
	User                  *User
	TrackingNacional      string  // opcional
	TrackingInternacional string  // opcional
	Libros                []Libro // opcional

}
