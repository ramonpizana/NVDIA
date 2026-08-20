package vcard

type Vcard struct {
	Name  string
	Store string
	Price int // Mexican pesos (MXN), including the price shown by the retailer.
	Link  string
	Stock bool
}

func New(name, store string, price int, link string, stock bool) *Vcard {
	return &Vcard{
		Name:  name,
		Store: store,
		Price: price,
		Link:  link,
		Stock: stock,
	}
}
