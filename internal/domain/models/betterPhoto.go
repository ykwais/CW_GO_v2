package models

type BetterPhoto struct {
	Data      []byte
	VehicleId int64
	Brand     string
	Model     string
	TotalCost float64
	Url       string
}
