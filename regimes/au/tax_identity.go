package au

import (
	"errors"
	"strconv"

	"github.com/invopop/gobl/cbc"
	"github.com/invopop/gobl/tax"
	"github.com/invopop/validation"
)

// Weights for ABN checksum
var taxWeightTableABN = [11]int{10, 1, 3, 5, 7, 9, 11, 13, 15, 17, 19}

func validateTaxIdentity(tID *tax.Identity) error {
	return validation.ValidateStruct(tID,
		validation.Field(&tID.Code,
			validation.By(validateTaxCode),
			validation.Skip,
		),
	)
}

func validateTaxCode(value interface{}) error {
	code, ok := value.(cbc.Code)
	if !ok || code == "" {
		return errors.New("invalid format")
	}
	val := code.String()
	return businessNumberCheck(val)
}

// Source: https://abr.business.gov.au/Help/AbnFormat
func businessNumberCheck(val string) error {
	if z, _ := strconv.Atoi(val); z == 0 {
		return errors.New("invalid format")
	}
	if len(val) != 11 {
		return errors.New("invalid format")
	}

	firstDigit, err := strconv.Atoi(string(val[0]))
	if err != nil {
		return errors.New("invalid format")
	}
	firstDigit--
	modifiedABN := strconv.Itoa(firstDigit) + val[1:]
	sum := 0
	for i := 0; i < 11; i++ {
		digit, err := strconv.Atoi(string(modifiedABN[i]))
		if err != nil {
			return errors.New("invalid format")
		}
		sum += digit * taxWeightTableABN[i]
	}
	if sum%89 == 0 {
		return nil
	}
	return errors.New("checksum mismatch")
}
