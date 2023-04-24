package entities

type Balance struct {
	Current   float64 `json:"current"`
	Withdrawn int     `json:"withdrawn"`
}

type Withdraw struct {
	Order string `json:"order"`
	Sum   int    `json:"sum"`
}

type WithdrawWithTime struct {
	Withdraw
	ProcessedAt string `json:"processed_at"`
}
