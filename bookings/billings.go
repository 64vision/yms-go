package bookings

import (
	"fmt"
	"time"
)

type Billing struct {
	ID          int       `json:"id" pg:"id,pk"`
	Type        string    `json:"type" pg:"type"`
	RefID       int       `json:"ref_id" pg:"ref_id"`
	Description string    `json:"description" pg:"description"`
	Discount    float64   `json:"discount" pg:"discount"`
	TotalAmount float64   `json:"total_amount" pg:"total_amount "`
	Status      string    `json:"status" pg:"status"`
	Client      int       `json:"client" pg:"client"`
	CreatedAt   time.Time `json:"created_at" pg:"created_at"`
	UpdatedAt   time.Time `json:"updated_at" pg:"updated_at"`
}

func (b *Billing) GenerateBilling() {
	b.CreatedAt = time.Now()
	_, errdb := DBM.Model(b).Insert()
	if errdb != nil {
		panic(errdb)
		fmt.Println("Error on generating billing!")
	}
}

func GetBillingByRefID(refID int, _type string) *Billing {

	var billing Billing
	err := DBM.Model(&billing).Where("ref_id = ? AND type = ?", refID, _type).Select()
	if err != nil {
		panic(err)
	}
	return &billing
}
