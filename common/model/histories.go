package model

type HistoricalData struct {
	High   float64 `json:"high"`
	Low    float64 `json:"low"`
	Open   float64 `json:"open"`
	Close  float64 `json:"close"`
	Time   int64   `json:"time"`
	Change float64 `json:"change"`
}
