package model

import (
	"database/sql"
	"database/sql/driver"
	"fmt"
	"strconv"
	"strings"
)

type Coordinates struct {
	Longitude float64 `json:"longitude"`
	Latitude  float64 `json:"latitude"`
}

func (p Coordinates) Value() (driver.Value, error) {
	return []byte(fmt.Sprintf("POINT(%.8f %.8f)", p.Longitude, p.Latitude)), nil
}

func (p *Coordinates) Scan(value interface{}) error {
	if value == nil {
		*p = Coordinates{}
		return nil
	}

	v, ok := value.(string)
	if !ok {
		return fmt.Errorf("Coordinates.Scan: expected string, got %T (%v)", value, value)
	}

	coords := strings.TrimLeft(strings.TrimRight(v, ")"), "POINT(")

	longitude, err := strconv.ParseFloat(strings.Split(coords, " ")[0], 64)
	if err != nil {
		return fmt.Errorf("Coordinates.Scan: cannot parse longitude")
	}

	latitude, err := strconv.ParseFloat(strings.Split(coords, " ")[1], 64)
	if err != nil {
		return fmt.Errorf("Coordinates.Scan: cannot parse latitude")
	}

	*p = Coordinates{
		Longitude: longitude,
		Latitude:  latitude,
	}

	return nil
}

// check interfaces
var (
	_ driver.Valuer = Coordinates{}
	_ sql.Scanner   = &Coordinates{}
)
