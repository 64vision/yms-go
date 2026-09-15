package bookings

import (
	"fmt"
	u "gollux/utils"
	"time"
)

type Booking struct {
	ID            int         `json:"id"`
	Reference     string      `json:"reference"`
	ShippingID    int         `json:"shipping_id"`
	ContainerNo   string      `json:"container_no"`
	BookingType   string      `json:"booking_type"`
	YardID        int         `json:"yard_id"`
	BookingDate   string      `json:"booking_date"`
	SlotID        int         `json:"slot_id"`
	SlotTime      string      `json:"slot_time"`
	ClientID      int         `json:"client_id"`
	ClientName    string      `json:"client_name"`
	ContainerType string      `json:"container_type"`
	ContainerLoad string      `json:"container_load"`
	ContainerSize string      `json:"container_size"`
	Truck         interface{} `json:"truck"`
	Affiliation   string      `json:"affiliation"`
	CreatedAt     time.Time   `json:"created_at"`
	Status        string      `json:"status"`
	Notes         string      `json:"notes"`
	BookingFee    float64     `json:"booking_fee"`
	TotalAmount   float64     `json:"total_amount"`
	DocsFee       float64     `json:"docs_fee"`
	ServiceFee    float64     `json:"service_fee"`
	PaymentStatus string      `json:"payment_status"`
	DailyRate     float64     `json:"daily_rate"`

	WithdrawalDate    string    `json:"withdrawal_date"`
	NumberOfDays      int       `json:"number_of_days"`
	WithdrawalAmount  float64   `json:"withdrawal_amount"`
	WithdrawalPayment string    `json:"withdrawal_payment"`
	Remarks           string    `json:"remarks"`
	WithdrawnAt       time.Time `json:"withdrawn_at"`
	WithdrawSlotID    int       `json:"withdraw_slot_id"`
	WithdrawSlotTime  string    `json:"withdraw_slot_time"`
	Zone              string    `json:"zone"`
	Slot              string    `json:"slot"`
	Index             string    `json:"index"`
}

func (b *Booking) GetBookingByID() map[string]interface{} {

	var booking Booking
	_, err := DBM.Query(&booking, `SELECT * FROM bookings WHERE id = ?`, b.ID)
	if err != nil {
		//panic(err)
		return u.Message(false, err.Error())
	}
	response := u.Message(true, "oK")
	response["booking"] = booking
	return response
}

func (b *Booking) PaidBilling() {
	billing := GetBillingByRefID(b.ID, "withdrawal")
	billing.Status = "Paid"
	billing.UpdatedAt = time.Now()
	billing.TotalAmount = b.WithdrawalAmount
	errdb := DBM.Update(billing)
	if errdb != nil {
		panic(errdb)
		fmt.Println("Error on updating billing!")
	}
}

func (b *Booking) Withdraw() map[string]interface{} {

	res, errdb := DBM.Exec(`UPDATE bookings
	SET withdrawal_date = ?,
		withdrawal_amount = ?,
		withdrawal_payment = ?,
		number_of_days = ?,
		remarks = ?,
		withdraw_slot_time=?,
		withdraw_slot_id=?
	WHERE id = ?`, b.WithdrawalDate, b.WithdrawalAmount, "Paid", b.NumberOfDays, "For Withdrawal", b.WithdrawSlotTime, b.WithdrawSlotID, b.ID)
	if errdb != nil {
		//panic(errdb)
		return u.Message(false, errdb.Error())
	}
	fmt.Println(res.Model())
	if res.RowsAffected() == 0 {
		return u.Message(false, "No rows affected!")
	}
	b.PaidBilling()
	bookedWithdrawSlot(b.WithdrawSlotID, b.WithdrawSlotTime)
	response := u.Message(true, "Save")
	return response
}

func (b *Booking) Cancel() map[string]interface{} {

	//var booking Booking
	res, err := DBM.Exec(`update bookings set status = 'Cancelled' where id = ?`, b.ID)
	if err != nil {
		panic(err)
		return u.Message(false, err.Error())
	}
	if res.RowsAffected() > 0 {
		//update slot booked count
		returnSlot(b.SlotID, b.SlotTime)
	}

	response := u.Message(true, "Booking has been successfully cancelled!")

	return response
}
func (b *Booking) AddBooking() map[string]interface{} {
	if b.ValidateBookingDate() != "Available" {
		return u.Message(false, "Slot is not available or already full!")
	}
	b.Status = "Active"
	b.PaymentStatus = "Credits"
	b.CreatedAt = time.Now()
	_, errdb := DBM.Model(b).Insert()
	if errdb != nil {
		panic(errdb)
		return u.Message(false, "Failed to create account, connection error")
	}

	bookedSlot(b.SlotID, b.SlotTime)

	billing := &Billing{}
	billing.Type = "booking"
	billing.RefID = b.ID
	billing.Description = "Booking and docs fee"
	billing.TotalAmount = b.TotalAmount
	billing.Discount = 0
	billing.Status = "Paid"
	billing.Client = b.ClientID
	billing.GenerateBilling()

	withdrawbilling := &Billing{}
	withdrawbilling.Type = "withdrawal"
	withdrawbilling.RefID = b.ID
	withdrawbilling.Description = "Container withdrawal Payment"
	withdrawbilling.Discount = 0
	withdrawbilling.Status = "Pending"
	withdrawbilling.Client = b.ClientID
	withdrawbilling.GenerateBilling()

	response := u.Message(true, "Slot has beed successfully booked!")
	response["booking"] = b
	response["booking_id"] = b.ID
	return response
}

func bookedWithdrawSlot(slotID int, slotTime string) {

	res, errdb := DBM.Exec(`UPDATE booking_slots
					SET slots = (
						SELECT jsonb_agg(
							CASE
								WHEN slot->>'time' = ? THEN
									jsonb_set(
										slot,
										'{booked_withdraw}',
										to_jsonb((slot->>'booked_withdraw')::int + 1)
									)
								ELSE
									slot
							END
						)
						FROM jsonb_array_elements(slots) AS slot
					)
					WHERE id = ?`, slotTime, slotID)
	if errdb != nil {
		panic(errdb)
	}
	fmt.Println(res.Model())
	fmt.Println("RowsAffected", res.RowsAffected())
}

func returnWithdrawSlot(slotID int, slotTime string) {

	res, errdb := DBM.Exec(`UPDATE booking_slots
					SET slots = (
						SELECT jsonb_agg(
							CASE
								WHEN slot->>'time' = ? THEN
									jsonb_set(
										slot,
										'{booked_withdraw}',
										to_jsonb((slot->>'booked_withdraw')::int - 1)
									)
								ELSE
									slot
							END
						)
						FROM jsonb_array_elements(slots) AS slot
					)
					WHERE id = ?`, slotTime, slotID)
	if errdb != nil {
		panic(errdb)
	}
	fmt.Println(res.Model())
	fmt.Println("RowsAffected", res.RowsAffected())
}
func returnSlot(slotID int, slotTime string) {

	res, errdb := DBM.Exec(`UPDATE booking_slots
					SET slots = (
						SELECT jsonb_agg(
							CASE
								WHEN slot->>'time' = ? THEN
									jsonb_set(
										slot,
										'{booked_drop}',
										to_jsonb((slot->>'booked_drop')::int - 1)
									)
								ELSE
									slot
							END
						)
						FROM jsonb_array_elements(slots) AS slot
					)
					WHERE id = ?`, slotTime, slotID)
	if errdb != nil {
		panic(errdb)
	}
	fmt.Println(res.Model())
	fmt.Println("RowsAffected", res.RowsAffected())
}
func bookedSlot(slotID int, slotTime string) {

	res, errdb := DBM.Exec(`UPDATE booking_slots
					SET slots = (
						SELECT jsonb_agg(
							CASE
								WHEN slot->>'time' = ? THEN
									jsonb_set(
										slot,
										'{booked_drop}',
										to_jsonb((slot->>'booked_drop')::int + 1)
									)
								ELSE
									slot
							END
						)
						FROM jsonb_array_elements(slots) AS slot
					)
					WHERE id = ?`, slotTime, slotID)
	if errdb != nil {
		panic(errdb)
	}
	fmt.Println(res.Model())
	fmt.Println("RowsAffected", res.RowsAffected())
}

func (b *Booking) ValidateBookingDate() string {
	fmt.Println("SlotID", b.SlotID)
	fmt.Println("SlotTime", b.SlotTime)
	var slot Slot
	res, err := DBM.Query(&slot, `SELECT  slot->>'time' as time, slot->>'status' as status, slot->>'drop' as drop, slot->>'booked_drop' as booked_drop
FROM booking_slots,
jsonb_array_elements(slots) AS slot
WHERE id = ?
  AND slot->>'time' = ?`, b.SlotID, b.SlotTime)
	if err != nil {
		panic(err)
		return err.Error()
	}
	if res.RowsReturned() == 0 {
		fmt.Println("Slot not found!")
		return "Slot not found!"
	}

	if slot.BookedDrop >= slot.Drop {
		return "Slot is full!"
	}
	return "Available"
}
