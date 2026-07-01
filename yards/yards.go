package yards

type Yard struct {
	ID       int      `json:"id"`
	Name     string   `json:"name"`
	Address  string   `json:"address"`
	Owner    int      `json:"owner"`
	Capacity Capacity `json:"capacity"`
	Remarks  string   `json:"remarks"`
	Status   string   `json:"status"`
}
type Capacity struct {
	WithdrawalPerHour string `json:"withdrawal_per_hour"`
	ReturnPerHour     string `json:"return_per_hour"`
	StoragePerHour    string `json:"storage_per_hour"`
	RateDuePerHour    string `json:"rate_due_per_hour"`
	BookingRate       string `json:"booking_rate"`
	DailyStorageRate  string `json:"daily_storage_rate"`
	DocumentFee       string `json:"document_fee"`
}
type ShippingLine struct {
	ID      int    `json:"id"`
	Name    string `json:"name"`
	Owner   int    `json:"owner"`
	Status  string `json:"status"`
	Address string `json:"address"`
}
