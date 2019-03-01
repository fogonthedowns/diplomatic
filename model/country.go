package model

import (
	"database/sql/driver"
	"encoding/json"
	"errors"
	"regexp"
)

const (
	ENGLAND         = Country("England")
	FRANCE          = Country("France")
	RUSSIA          = Country("Russia")
	GERMANY         = Country("Germany")
	AUSTRIA_HUNGARY = Country("Austria-Hungary")
	ITALY           = Country("Italy")
	TURKEY          = Country("Turkey")
	NONE            = Country("None")
)

type Country string

var reCountry = regexp.MustCompile(`^(|England|France|Russia|Germany|Austria-Hungary|Italy|Turkey|None)$`)

func (d *Country) validate(s string) error {
	if matched := reCountry.MatchString(s); matched == false {
		return errors.New("Invalid value for Country")
	}
	return nil
}

func (d *Country) assign(s string) {
	*d = Country(s)
}

func (d *Country) UnmarshalJSON(b []byte) (err error) {
	var s string
	if err := json.Unmarshal(b, &s); err == nil {
		if err = d.validate(s); err == nil {
			d.assign(s)
		}
		return err
	}
	return err
}

func (d Country) Value() (driver.Value, error) {
	return string(d), nil
}

func (z *Country) Scan(s interface{}) (err error) {
	if z == nil {
		return errors.New("Country: Scan on nil pointer")
	}
	*z = Country(string(s.([]uint8)))
	return nil
}
