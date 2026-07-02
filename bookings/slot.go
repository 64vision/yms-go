package bookings

import (
	u "gollux/utils"
	"time"
)

type BookingSlot struct {
	ID          int       `json:"id"`
	BookingDate string    `json:"booking_date"`
	Slots       []Slot    `pg:",array" json:"slots"`
	Status      string    `json:"status"`
	YardID      int       `json:"yard_id"`
	CreatedAt   time.Time `json:"created_at"`
	CreatedBy   string    `json:"created_by"`
	Type        string    `json:"type"`
}
type Slot struct {
	Time       string   `json:"time"`
	Capacity   int      `json:"capacity"`
	Status     string   `json:"status"`
	BookingIds []string `json:"booking_ids"`
}

func (bs *BookingSlot) OpenBookingDate() map[string]interface{} {
	if bs.ValidateBookingDate() != "Available" {
		return u.Message(false, "Booking date already exist!")
	}
	bs.Status = "Open"
	bs.CreatedAt = time.Now()
	_, errdb := DBM.Model(bs).Insert()
	if errdb != nil {
		panic(errdb)
		return u.Message(false, "Failed to create account, connection error")
	}
	response := u.Message(true, "Booking date opened successfully!")
	return response
}

func (bs *BookingSlot) ValidateBookingDate() string {
	var slot BookingSlot
	_, err := DBM.Query(&slot, `SELECT * FROM booking_slots where booking_date=?`, bs.BookingDate)
	if err != nil {
		panic(err)
		return err.Error()
	}
	if slot.ID == 0 {
		return "Available"
	}
	return "Already Exist"
}
