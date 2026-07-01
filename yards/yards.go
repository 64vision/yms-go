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
	WithdrawalPerHour int     `json:"withdrawal_per_hour"`
	ReturnPerHour     int     `json:"return_per_hour"`
	StoragePerHour    int     `json:"storage_per_hour"`
	RateDuePerHour    float64 `json:"rate_due_per_hour"`
	BookingRate       float64 `json:"booking_rate"`
	DailyStorageRate  float64 `json:"daily_storage_rate"`
	DocumentFee       float64 `json:"document_fee"`
}
