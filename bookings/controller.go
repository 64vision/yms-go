package bookings

import (
	"encoding/json"
	u "gollux/utils"
	"net/http"
)

func OpenDate(w http.ResponseWriter, r *http.Request) {
	(w).Header().Set("Access-Control-Allow-Origin", "*")
	slot := &BookingSlot{}
	var resp map[string]interface{}
	err := json.NewDecoder(r.Body).Decode(slot) //decode the request body into struct and failed if any error occur
	if err != nil {
		//panic(err)
		u.Respond(w, u.Message(false, "Invalid request"))
		return
	}

	slot.CreatedBy = r.Context().Value("user").(string)
	resp = slot.OpenBookingDate()
	u.Respond(w, resp)
}
