package emails

type OrderPlaced struct {
	Name    string
	ID      string
	Address string
	Items   []OrderPlacedItem
}

type OrderPlacedItem struct {
	Name     string
	Quantity uint
	Price    float32
}
