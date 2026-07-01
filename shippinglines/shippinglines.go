package shippinglines

type ShippingLine struct {
	ID      int    `json:"id"`
	Name    string `json:"name"`
	Owner   int    `json:"owner"`
	Status  string `json:"status"`
	Address string `json:"address"`
}
