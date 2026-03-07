package models

import "time"

type Stat struct
{
	Name string 'json:"symbol"'
	Price float64 'json:"price"'
	Source string 'json:"source"'
	Timedump time.Time 'json:"timedump"'
}
