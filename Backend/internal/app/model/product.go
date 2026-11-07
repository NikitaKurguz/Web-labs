package model

import "strconv"

type Product struct {
	ID       int
	Name     string
	Img      string
	Type     string
	Color    string
	Category string
	Price    int
}

func (p *Product) GetPriceRange() string {
	switch {
	case p.Price == 0:
		return "Бесплатно"
	case p.Price <= 500:
		return "До 500"
	case p.Price <= 2000:
		return "500-2000"
	default:
		return "От 2000"
	}
}

func (p *Product) PriceLabel() string {
	if p.Price == 0 {
		return "Бесплатно"
	}
	return strconv.Itoa(p.Price) + " Руб."
}
