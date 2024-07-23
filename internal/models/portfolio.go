package models

type PortfolioAsset struct {
	Ticker     string  `json:"ticker"`
	Name       string  `json:"name"`
	Quantity   float64 `json:"quantity"`
	TotalValue float64 `json:"totalValue"`
	Return     float64 `json:"return"`
}

type Portfolio struct {
	TotalValue    float64          `json:"totalValue"`
	AvailableCash float64          `json:"availableCash"`
	Assets        []PortfolioAsset `json:"assets"`
}
