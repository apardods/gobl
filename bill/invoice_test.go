package bill_test

import (
	"context"
	"testing"

	_ "github.com/invopop/gobl" // load regions
	"github.com/invopop/gobl/bill"
	"github.com/invopop/gobl/cal"
	"github.com/invopop/gobl/currency"
	"github.com/invopop/gobl/internal"
	"github.com/invopop/gobl/l10n"
	"github.com/invopop/gobl/num"
	"github.com/invopop/gobl/org"
	"github.com/invopop/gobl/pay"
	"github.com/invopop/gobl/regimes/common"
	"github.com/invopop/gobl/tax"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestRemoveIncludedTax(t *testing.T) {
	i := &bill.Invoice{
		Code: "123TEST",
		Tax: &bill.Tax{
			PricesInclude: common.TaxCategoryVAT,
		},
		Supplier: &org.Party{
			TaxID: &tax.Identity{
				Country: l10n.ES,
				Code:    "B98602642",
			},
		},
		Customer: &org.Party{
			TaxID: &tax.Identity{
				Country: l10n.ES,
				Code:    "54387763P",
			},
		},
		IssueDate: cal.MakeDate(2022, 6, 13),
		Lines: []*bill.Line{
			{
				Quantity: num.MakeAmount(1, 0),
				Item: &org.Item{
					Name:  "Test Item",
					Price: num.MakeAmount(100000, 2),
				},
				Taxes: tax.Set{
					{
						Category: "VAT",
						Percent:  num.NewPercentage(21, 2),
					},
				},
				Discounts: []*bill.LineDiscount{
					{
						Reason:  "Testing",
						Percent: num.NewPercentage(10, 2),
					},
				},
			},
		},
	}

	require.NoError(t, i.Calculate())

	i2 := i.RemoveIncludedTaxes()
	require.NoError(t, i2.Calculate())

	assert.Equal(t, "1000.00", i.Lines[0].Item.Price.String())

	assert.Empty(t, i2.Tax.PricesInclude)
	l0 := i2.Lines[0]
	assert.Equal(t, "826.4463", l0.Item.Price.String())
	assert.Equal(t, "826.4463", l0.Sum.String())
	assert.Equal(t, "826.4463", l0.Sum.String())
	assert.Equal(t, "82.6446", l0.Discounts[0].Amount.String())
	assert.Equal(t, "743.8017", l0.Total.String())

	assert.Equal(t, "743.80", i2.Totals.Total.String())
	assert.Equal(t, "900.00", i2.Totals.Payable.String())
}

func TestRemoveIncludedTax2(t *testing.T) {
	i := &bill.Invoice{
		Code: "123TEST",
		Tax: &bill.Tax{
			PricesInclude: common.TaxCategoryVAT,
		},
		Supplier: &org.Party{
			TaxID: &tax.Identity{
				Country: l10n.ES,
				Code:    "B98602642",
			},
		},
		Customer: &org.Party{
			TaxID: &tax.Identity{
				Country: l10n.ES,
				Code:    "54387763P",
			},
		},
		IssueDate: cal.MakeDate(2022, 6, 13),
		Lines: []*bill.Line{
			{
				Quantity: num.MakeAmount(1, 0),
				Item: &org.Item{
					Name:  "Item",
					Price: num.MakeAmount(4320, 2),
				},
				Taxes: tax.Set{
					{
						Category: "VAT",
						Percent:  num.NewPercentage(6, 2),
					},
				},
			},
			{
				Quantity: num.MakeAmount(1, 0),
				Item: &org.Item{
					Name:  "Item 2",
					Price: num.MakeAmount(259, 2),
				},
			},
			{
				Quantity: num.MakeAmount(1, 0),
				Item: &org.Item{
					Name:  "Item 3",
					Price: num.MakeAmount(300, 2),
				},
			},
		},
	}

	require.NoError(t, i.Calculate())

	i2 := i.RemoveIncludedTaxes()
	require.NoError(t, i2.Calculate())
	l0 := i2.Lines[0]
	assert.Equal(t, "40.7547", l0.Item.Price.String())
	assert.Equal(t, "40.7547", l0.Total.String())

	assert.Equal(t, "46.34", i2.Totals.Total.String())
	assert.Equal(t, "2.45", i2.Totals.Tax.String())
	assert.Equal(t, "48.79", i2.Totals.Payable.String())
}

func TestRemoveIncludedTax3(t *testing.T) {
	i := &bill.Invoice{
		Code: "123TEST",
		Tax: &bill.Tax{
			PricesInclude: common.TaxCategoryVAT,
		},
		Supplier: &org.Party{
			TaxID: &tax.Identity{
				Country: l10n.ES,
				Code:    "B98602642",
			},
		},
		Customer: &org.Party{
			TaxID: &tax.Identity{
				Country: l10n.ES,
				Code:    "54387763P",
			},
		},
		IssueDate: cal.MakeDate(2022, 6, 13),
		Lines: []*bill.Line{
			{
				Quantity: num.MakeAmount(1, 0),
				Item: &org.Item{
					Name:  "Item",
					Price: num.MakeAmount(23666, 2),
				},
				Taxes: tax.Set{
					{
						Category: "VAT",
						Percent:  num.NewPercentage(6, 2),
					},
				},
			},
			{
				Quantity: num.MakeAmount(2, 0),
				Item: &org.Item{
					Name:  "Item 2",
					Price: num.MakeAmount(23667, 2),
				},
				Taxes: tax.Set{
					{
						Category: "VAT",
						Percent:  num.NewPercentage(6, 2),
					},
				},
			},
			{
				Quantity: num.MakeAmount(12, 0),
				Item: &org.Item{
					Name:  "Item 3",
					Price: num.MakeAmount(1000, 2),
				},
				Taxes: tax.Set{
					{
						Category: "VAT",
						Percent:  num.NewPercentage(13, 2),
					},
				},
			},
			{
				Quantity: num.MakeAmount(18, 0),
				Item: &org.Item{
					Name:  "Local tax",
					Price: num.MakeAmount(150, 2),
				},
			},
		},
	}

	require.NoError(t, i.Calculate())

	i2 := i.RemoveIncludedTaxes()
	require.NoError(t, i2.Calculate())
	assert.Equal(t, "223.2642", i2.Lines[0].Total.String())
	assert.Equal(t, "106.19472", i2.Lines[2].Total.String()) // more accuracy

	assert.Equal(t, "803.00", i2.Totals.Sum.String())
	assert.Equal(t, "803.00", i2.Totals.Total.String())
	assert.Equal(t, "54.00", i2.Totals.Tax.String())
	assert.Equal(t, "857.00", i2.Totals.Payable.String())
}

func TestRemoveIncludedTaxQuantity(t *testing.T) {
	i := &bill.Invoice{
		Code: "123TEST",
		Tax: &bill.Tax{
			PricesInclude: common.TaxCategoryVAT,
		},
		Supplier: &org.Party{
			TaxID: &tax.Identity{
				Country: l10n.ES,
				Code:    "B98602642",
			},
		},
		Customer: &org.Party{
			TaxID: &tax.Identity{
				Country: l10n.ES,
				Code:    "54387763P",
			},
		},
		IssueDate: cal.MakeDate(2022, 6, 13),
		Lines: []*bill.Line{
			{
				Quantity: num.MakeAmount(100, 0),
				Item: &org.Item{
					Name:  "Test Item",
					Price: num.MakeAmount(1000, 2),
				},
				Taxes: tax.Set{
					{
						Category: "VAT",
						Percent:  num.NewPercentage(21, 2),
					},
				},
				Discounts: []*bill.LineDiscount{
					{
						Reason:  "Testing",
						Percent: num.NewPercentage(10, 2),
					},
				},
			},
		},
	}

	require.NoError(t, i.Calculate())

	i2 := i.RemoveIncludedTaxes()
	require.NoError(t, i2.Calculate())

	assert.Empty(t, i2.Tax.PricesInclude)
	l0 := i2.Lines[0]
	assert.Equal(t, "8.264463", l0.Item.Price.String())
	assert.Equal(t, "826.446300", l0.Sum.String())
	assert.Equal(t, "82.644630", l0.Discounts[0].Amount.String())
	assert.Equal(t, "743.801670", l0.Total.String())
	assert.Equal(t, "10.00", i.Lines[0].Item.Price.String())

	assert.Equal(t, "743.80", i2.Totals.Total.String())
	assert.Equal(t, "900.00", i2.Totals.Payable.String())
}

func TestRemoveIncludedTaxDeep(t *testing.T) {
	i := &bill.Invoice{
		Code: "123TEST",
		Tax: &bill.Tax{
			PricesInclude: common.TaxCategoryVAT,
		},
		Supplier: &org.Party{
			TaxID: &tax.Identity{
				Country: l10n.ES,
				Code:    "B98602642",
			},
		},
		Customer: &org.Party{
			TaxID: &tax.Identity{
				Country: l10n.ES,
				Code:    "54387763P",
			},
		},
		IssueDate: cal.MakeDate(2022, 6, 13),
		Lines: []*bill.Line{
			{
				Quantity: num.MakeAmount(364, 0),
				Item: &org.Item{
					Name:  "Test Item",
					Price: num.MakeAmount(5178, 2),
				},
				Taxes: tax.Set{
					{
						Category: "VAT",
						Percent:  num.NewPercentage(6, 2),
					},
				},
			},
			{
				Quantity: num.MakeAmount(1, 0),
				Item: &org.Item{
					Name:  "Test Item 2",
					Price: num.MakeAmount(5208, 2),
				},
				Taxes: tax.Set{
					{
						Category: "VAT",
						Percent:  num.NewPercentage(6, 2),
					},
				},
			},
		},
	}

	require.NoError(t, i.Calculate())

	i2 := i.RemoveIncludedTaxes()

	require.NoError(t, i2.Calculate())

	//data, _ := json.MarshalIndent(i2, "", "  ")
	//t.Log(string(data))

	assert.Empty(t, i2.Tax.PricesInclude)
	l0 := i2.Lines[0]
	assert.Equal(t, "48.849057", l0.Item.Price.String()) // note extra digit!
	assert.Equal(t, "17781.056748", l0.Sum.String())
	l1 := i2.Lines[1]
	assert.Equal(t, "49.1321", l1.Item.Price.String())
	assert.Equal(t, "49.1321", l1.Sum.String())

	assert.Equal(t, "17830.19", i2.Totals.Total.String())
	assert.Equal(t, "18900.00", i2.Totals.Payable.String())
}

func TestCalculate(t *testing.T) {
	i := &bill.Invoice{
		Code: "123TEST",
		Tax: &bill.Tax{
			PricesInclude: common.TaxCategoryVAT,
		},
		Supplier: &org.Party{
			TaxID: &tax.Identity{
				Country: l10n.ES,
				Code:    "B98602642",
			},
		},
		Customer: &org.Party{
			TaxID: &tax.Identity{
				Country: l10n.ES,
				Code:    "54387763P",
			},
		},
		IssueDate: cal.MakeDate(2022, 6, 13),
		Lines: []*bill.Line{
			{
				Quantity: num.MakeAmount(10, 0),
				Item: &org.Item{
					Name:  "Test Item",
					Price: num.MakeAmount(10000, 2),
				},
				Taxes: tax.Set{
					{
						Category: "VAT",
						Rate:     "standard",
					},
				},
				Discounts: []*bill.LineDiscount{
					{
						Reason:  "Testing",
						Percent: num.NewPercentage(10, 2),
					},
				},
				Charges: []*bill.LineCharge{
					{
						Reason:  "Testing Charge",
						Percent: num.NewPercentage(5, 2),
					},
				},
			},
		},
		Outlays: []*bill.Outlay{
			{
				Description: "Something paid in advance",
				Amount:      num.MakeAmount(1000, 2),
			},
		},
		Payment: &bill.Payment{
			Advances: []*pay.Advance{
				{
					Description: "Test Advance",
					Percent:     num.NewPercentage(30, 2), // 30%
				},
			},
		},
	}

	require.NoError(t, i.Calculate())
	assert.Equal(t, i.Totals.Sum.String(), "950.00")
	assert.Equal(t, i.Totals.Total.String(), "785.12")
	assert.Equal(t, i.Totals.Tax.String(), "164.88")
	assert.Equal(t, i.Totals.TotalWithTax.String(), "950.00")
	assert.Equal(t, i.Payment.Advances[0].Amount.String(), "285.00")
	assert.Equal(t, i.Totals.Advances.String(), "285.00")
	assert.Equal(t, i.Totals.Payable.String(), "960.00")
	assert.Equal(t, i.Totals.Due.String(), "675.00")
}

func TestValidation(t *testing.T) {
	inv := &bill.Invoice{
		Currency:  currency.EUR,
		IssueDate: cal.MakeDate(2022, 6, 13),
		Tax: &bill.Tax{
			PricesInclude: common.TaxCategoryVAT,
		},
		Supplier: &org.Party{
			Name: "Test Supplier",
			TaxID: &tax.Identity{
				Country: l10n.ES,
				Code:    "B98602642",
			},
		},
		Customer: &org.Party{
			Name: "Test Customer",
			TaxID: &tax.Identity{
				Country: l10n.ES,
				Code:    "54387763P",
			},
		},
		Lines: []*bill.Line{
			{
				Quantity: num.MakeAmount(10, 0),
				Item: &org.Item{
					Name:  "Test Item",
					Price: num.MakeAmount(10000, 2),
				},
				Taxes: tax.Set{
					{
						Category: "VAT",
						Rate:     "standard",
					},
				},
			},
		},
	}
	require.NoError(t, inv.Calculate())
	ctx := context.Background()
	err := inv.ValidateWithContext(ctx)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "code: cannot be blank")
	ctx = context.WithValue(ctx, internal.KeyDraft, true)
	assert.NoError(t, inv.ValidateWithContext(ctx))
}
