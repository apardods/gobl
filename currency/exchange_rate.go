package currency

import (
	"fmt"

	"github.com/invopop/gobl/num"
	"github.com/invopop/validation"
)

// ExchangeRate contains data on the rate to be used when converting amounts from
// one currency into another.
//
// For reference, naming here is based on the following english grammar examples:
// - Exchange from USD to EUR.
// - Convert from USD into EUR.
//
// It should be possible to take any amount in the matching currency and multiply it
// by the amount defined in the exchange rate to determine the value.
//
// For example, our document is in EUR and some amounts are defined in USD. Our
// ExchangeRate instance may be defined and used as:
//
//	  rate := &currency.ExchangeRate{
//		From: currency.USD,
//		To: currency.EUR,
//		Amount: "0.875967",
//	  }
//
//	  val := MakeAmount(100, 2) // 100.00 USD
//	  val.Multiply(rate.Amount) // 87.60 EUR
type ExchangeRate struct {
	// Currency code this will be converted from.
	From Code `json:"from" jsonschema:"title=From"`
	// Currency code this exchange rate will convert into.
	To Code `json:"to" jsonschema:"title=To"`
	// How much is 1 of the "from" currency worth in the "to" currency.
	Amount num.Amount `json:"amount" jsonschema:"title=Amount"`
}

// Validate ensures the content of the exchange rate looks good.
func (er *ExchangeRate) Validate() error {
	return validation.ValidateStruct(er,
		validation.Field(&er.From, validation.Required),
		validation.Field(&er.To, validation.Required),
		validation.Field(&er.Amount, num.Positive),
	)
}

// Convert performs the currency conversion defined by the exchange rate.
func (er *ExchangeRate) Convert(amount num.Amount) num.Amount {
	a := amount.Multiply(er.Amount)
	z := er.To.Def().Zero()
	return a.Rescale(z.Exp()) // ensure scale always matches destination currency
}

// MatchExchangeRate will attempt to find the matching exchange rate that
// will convert from one currency into another. Will return nil if no
// match is found or the currencies are the same.
func MatchExchangeRate(rates []*ExchangeRate, from, to Code) *ExchangeRate {
	if from == to {
		return nil
	}
	for _, rate := range rates {
		if rate.From == from && rate.To == to {
			return rate
		}
	}
	return nil
}

// Convert will convert the provided amount from one currency into another or return
// nil if no match can be found. If the currencies are the same, the original
// amount will be returned.
func Convert(rates []*ExchangeRate, from, to Code, amount num.Amount) *num.Amount {
	if from == to {
		return &amount
	}
	if rate := MatchExchangeRate(rates, from, to); rate != nil {
		a := rate.Convert(amount)
		return &a
	}
	return nil
}

type exchangeRateValidation struct {
	rates []*ExchangeRate
	to    Code
}

// Validate performs validation on the provided value to see if it
// is present in the exchange rates.
func (erv *exchangeRateValidation) Validate(val any) error {
	cur, ok := val.(Code)
	if !ok || cur == CodeEmpty {
		return nil
	}
	if cur == erv.to {
		return nil
	}
	for _, r := range erv.rates {
		if r.From == cur && r.To == erv.to {
			return nil
		}
	}
	return fmt.Errorf("no exchange rate defined for '%v' to '%v'", cur, erv.to)
}

// CanExchangeTo will check to see if the currency to be validated can
// be converted into one of the provided rates.
func CanExchangeTo(rates []*ExchangeRate, to Code) validation.Rule {
	return &exchangeRateValidation{
		rates: rates,
		to:    to,
	}
}
