package models

import "time"

type Stat struct {
	Base     string    `json:"base"`
	Quote    string    `json:"quote"`
	AskPrice float64   `json:"askPrice"`
	BidPrice float64   `json:"bidPrice"`
	Source   string    `json:"source"`
	Timedump time.Time `json:"timedump" swaggerignore:"true"`
}

type Exchange interface {
	GetStat(basecoin string, quotecoin string) (Stat, error) //Get information from market
}
