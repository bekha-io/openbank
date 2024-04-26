package types

import "time"

type Installment struct {
	Amount       Money
	Date         time.Time
	RemainAmount Money
}
