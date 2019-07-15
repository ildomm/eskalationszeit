package models

import (
	"github.com/go-openapi/swag"
)

// Currency currency
type Currency struct {
	// name
	Name string `json:"name,omitempty"`

	// pair
	Pair string `json:"pair,omitempty"`

	// symbol
	Symbol string `json:"symbol,omitempty"`

	// value
	Value *float32 `json:"value,omitempty"`
}

const (

	// CurrencySymbolUSD captures enum value "USD"
	CurrencySymbolUSD string = "USD"

	// CurrencySymbolBTC captures enum value "BTC"
	CurrencySymbolBTC string = "BTC"

	// CurrencySymbolNZD captures enum value "NZD"
	CurrencySymbolNZD string = "NZD"
)

// MarshalBinary interface implementation
func (m *Currency) MarshalBinary() ([]byte, error) {
	if m == nil {
		return nil, nil
	}
	return swag.WriteJSON(m)
}

// UnmarshalBinary interface implementation
func (m *Currency) UnmarshalBinary(b []byte) error {
	var res Currency
	if err := swag.ReadJSON(b, &res); err != nil {
		return err
	}
	*m = res
	return nil
}
