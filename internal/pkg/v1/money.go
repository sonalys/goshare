package v1

import (
	"fmt"
	"math"
)

type Money struct {
	Value     int32
	Precision uint8
	Currency  string
}

func NewMoney(value int32, precision uint8, currency string) Money {
	return Money{
		Value:     value,
		Precision: precision,
		Currency:  currency,
	}
}

func (m Money) String() string {
	dec, cents := splitValue(m.Value, m.Precision)
	return fmt.Sprintf("%s%d.%."+fmt.Sprint(m.Precision)+"d", m.Currency, dec, cents)
}

func splitValue(value int32, precision uint8) (wholePart, fractionalPart int32) {
	factor := int32(math.Pow10(int(precision)))
	wholePart = value / factor
	fractionalPart = value % factor
	return
}
