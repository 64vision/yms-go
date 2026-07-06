package bookings

import (
	"encoding/json"
	u "gollux/utils"
	"net/http"
)

func NewBooking_use(w http.ResponseWriter, r *http.Request) {
	(w).Header().Set("Access-Control-Allow-Origin", "*")
	entry := &Booking{}
	var resp map[string]interface{}
	err := json.NewDecoder(r.Body).Decode(entry) //decode the request body into struct and failed if any error occur
	if err != nil {
		//panic(err)
		u.Respond(w, u.Message(false, "Invalid request"))
		return
	}

	entry.ClientID = r.Context().Value("user").(int)
	resp = entry.AddBooking()
	u.Respond(w, resp)
}
