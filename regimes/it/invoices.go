package it

import (
	"regexp"

	"github.com/invopop/gobl/bill"
	"github.com/invopop/gobl/org"
	"github.com/invopop/gobl/tax"
	"github.com/invopop/validation"
)

// invoiceValidator adds validation checks to invoices which are relevant
// for the region.
type invoiceValidator struct {
	inv *bill.Invoice
}

// normalizeInvoice is used to ensure the invoice data is correct.
func normalizeInvoice(inv *bill.Invoice) {
	normalizeSupplier(inv.Supplier)
	normalizeCustomer(inv.Customer)
	for _, line := range inv.Lines {
		normalizeLine(line)
	}
}

func normalizeSupplier(party *org.Party) {
	if party == nil {
		return
	}
	if party.Ext == nil || party.Ext[ExtKeySDIFiscalRegime] == "" {
		if party.Ext == nil {
			party.Ext = make(tax.Extensions)
		}
		party.Ext[ExtKeySDIFiscalRegime] = "RF01" // Ordinary regime is default
	}
}

func normalizeCustomer(party *org.Party) {
	if party == nil {
		return
	}
	if !isItalianParty(party) {
		return
	}
	// If the party is an individual, move the fiscal code to the identities.
	if party.TaxID.Type == "individual" { //nolint:staticcheck
		id := &org.Identity{
			Key:  IdentityKeyFiscalCode,
			Code: party.TaxID.Code,
		}
		party.TaxID.Code = ""
		party.TaxID.Type = "" //nolint:staticcheck
		party.Identities = org.AddIdentity(party.Identities, id)
	}
}

func normalizeLine(line *bill.Line) {
	for _, tax := range line.Taxes {
		if tax.Ext == nil {
			continue
		}
		if tax.Ext.Has("it-sdi-retained-tax") {
			tax.Ext[ExtKeySDIRetained] = tax.Ext["it-sdi-retained-tax"]
			delete(tax.Ext, "it-sdi-retained-tax")
		}
		if tax.Ext.Has("it-sdi-nature") {
			tax.Ext[ExtKeySDIExempt] = tax.Ext["it-sdi-nature"]
			delete(tax.Ext, "it-sdi-nature")
		}
	}
}

func validateInvoice(inv *bill.Invoice) error {
	v := &invoiceValidator{inv: inv}
	return v.validate()
}

func (v *invoiceValidator) validate() error {
	inv := v.inv
	return validation.ValidateStruct(inv,
		validation.Field(&inv.Tax,
			validation.By(v.tax),
			validation.Skip,
		),
		validation.Field(&inv.Supplier,
			validation.By(v.supplier),
			validation.Skip,
		),
		validation.Field(&inv.Customer,
			validation.By(v.customer),
			validation.Skip,
		),
		validation.Field(&inv.Lines,
			validation.Each(
				bill.RequireLineTaxCategory(tax.CategoryVAT),
				validation.Skip,
			),
			validation.Skip,
		),
	)
}

func (v *invoiceValidator) tax(value any) error {
	obj, _ := value.(*bill.Tax)
	if obj == nil {
		return nil
	}
	return validation.ValidateStruct(obj,
		validation.Field(&obj.Ext,
			tax.ExtensionsHas(
				ExtKeySDIFormat,
				ExtKeySDIDocumentType,
			),
			validation.Skip,
		),
	)
}

func (v *invoiceValidator) supplier(value interface{}) error {
	supplier, ok := value.(*org.Party)
	if !ok {
		return nil
	}

	return validation.ValidateStruct(supplier,
		validation.Field(&supplier.TaxID,
			validation.Required,
			tax.RequireIdentityCode,
			validation.Skip,
		),
		validation.Field(&supplier.Addresses,
			validation.Required,
			validation.Each(validation.By(validateAddress)),
			validation.Skip,
		),
		validation.Field(&supplier.Registration,
			validation.By(validateRegistration),
			validation.Skip,
		),
		validation.Field(&supplier.Ext,
			tax.ExtensionsRequires(ExtKeySDIFiscalRegime),
			validation.Skip,
		),
	)
}

func (v *invoiceValidator) customer(value interface{}) error {
	customer, _ := value.(*org.Party)
	if customer == nil {
		return nil
	}

	// Customers must have either a Tax ID (PartitaIVA)
	// or fiscal identity (codice fiscale)
	return validation.ValidateStruct(customer,
		validation.Field(&customer.TaxID,
			validation.Required,
			validation.When(
				isItalianParty(customer) && !hasFiscalCode(customer),
				tax.RequireIdentityCode,
			),
			validation.Skip,
		),
		validation.Field(&customer.Addresses,
			validation.When(
				isItalianParty(customer),
				// TODO: address not required for simplified invoices
				validation.Each(validation.By(validateAddress)),
			),
			validation.Skip,
		),
		validation.Field(&customer.Identities,
			validation.When(
				isItalianParty(customer) && !hasTaxIDCode(customer),
				org.RequireIdentityKey(IdentityKeyFiscalCode),
			),
			validation.Skip,
		),
	)
}

func hasTaxIDCode(party *org.Party) bool {
	return party != nil && party.TaxID != nil && party.TaxID.Code != ""
}

func hasFiscalCode(party *org.Party) bool {
	if party == nil || party.TaxID == nil {
		return false
	}
	return org.IdentityForKey(party.Identities, IdentityKeyFiscalCode) != nil

}

func isItalianParty(party *org.Party) bool {
	if party == nil || party.TaxID == nil {
		return false
	}
	return party.TaxID.Country.In("IT")
}

func validateTaxCombo(c *tax.Combo) error {
	switch c.Category {
	case tax.CategoryVAT:
		return validation.ValidateStruct(c,
			validation.Field(&c.Ext,
				validation.When(
					c.Percent == nil,
					tax.ExtensionsRequires(ExtKeySDIExempt),
				),
				validation.Skip,
			),
		)
	case TaxCategoryIRPEF, TaxCategoryIRES, TaxCategoryINPS, TaxCategoryENASARCO, TaxCategoryENPAM:
		return validation.ValidateStruct(c,
			validation.Field(&c.Ext,
				tax.ExtensionsRequires(ExtKeySDIRetained),
				validation.Skip,
			),
		)
	}
	return nil
}

func validateAddress(value interface{}) error {
	v, ok := value.(*org.Address)
	if v == nil || !ok {
		return nil
	}
	// Post code and street in addition to the locality are required in Italian invoices.
	return validation.ValidateStruct(v,
		validation.Field(&v.Street, validation.Required),
		validation.Field(&v.Code,
			validation.Required,
			validation.Match(regexp.MustCompile(`^\d{5}$`)),
		),
	)
}

func validateRegistration(value interface{}) error {
	v, ok := value.(*org.Registration)
	if v == nil || !ok {
		return nil
	}
	return validation.ValidateStruct(v,
		validation.Field(&v.Entry, validation.Required),
		validation.Field(&v.Office, validation.Required),
	)
}
