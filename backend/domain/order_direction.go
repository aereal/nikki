package domain

import "fmt"

const (
	OrderDirectionAsc OrderDirection = iota
	OrderDirectionDesc
)

type OrderDirection int

func (d OrderDirection) String() string {
	switch d {
	case OrderDirectionAsc:
		return "ASC"
	case OrderDirectionDesc:
		return "DESC"
	default:
		return fmt.Sprintf("OrderDirection(%d)", d)
	}
}
