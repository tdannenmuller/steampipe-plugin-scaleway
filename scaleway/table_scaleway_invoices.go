package scaleway

import (
	"context"
	"fmt"
	"time"

	billing "github.com/scaleway/scaleway-sdk-go/api/billing/v2beta1"
	"github.com/scaleway/scaleway-sdk-go/scw"
	"github.com/turbot/steampipe-plugin-sdk/v5/grpc/proto"
	"github.com/turbot/steampipe-plugin-sdk/v5/plugin"
	"github.com/turbot/steampipe-plugin-sdk/v5/plugin/transform"
)

func tableScalewayInvoice(ctx context.Context) *plugin.Table {
	plugin.Logger(ctx).Debug("Initializing Scaleway Invoice table")
	return &plugin.Table{
		Name:        "scaleway_invoice",
		Description: "Invoices in your Scaleway account.",
		List: &plugin.ListConfig{
			Hydrate: listScalewayInvoices,
			KeyColumns: []*plugin.KeyColumn{
				{Name: "organization_id", Require: plugin.Optional},
			},
		},
		Columns: []*plugin.Column{
			{Name: "id", Type: proto.ColumnType_STRING, Description: "The unique identifier of the invoice.", Transform: transform.FromField("ID")},
			{Name: "organization_id", Type: proto.ColumnType_STRING, Description: "The organization ID associated with the invoice.", Transform: transform.FromField("OrganizationID")},
			{Name: "organization_name", Type: proto.ColumnType_STRING, Description: "The organization name associated with the invoice."},
			{Name: "start_date", Type: proto.ColumnType_TIMESTAMP, Description: "The start date of the billing period."},
			{Name: "stop_date", Type: proto.ColumnType_TIMESTAMP, Description: "The end date of the billing period."},
			{Name: "billing_period", Type: proto.ColumnType_TIMESTAMP, Description: "The billing period for the invoice."},
			{Name: "issued_date", Type: proto.ColumnType_TIMESTAMP, Description: "The date when the invoice was issued."},
			{Name: "due_date", Type: proto.ColumnType_TIMESTAMP, Description: "The due date for the invoice payment."},
			{Name: "total_untaxed", Type: proto.ColumnType_JSON, Description: "The total amount before tax."},
			{Name: "total_taxed", Type: proto.ColumnType_JSON, Description: "The total amount including tax."},
			{Name: "total_tax", Type: proto.ColumnType_JSON, Description: "The total tax amount."},
			{Name: "total_discount", Type: proto.ColumnType_JSON, Description: "The total discount amount."},
			{Name: "total_undiscount", Type: proto.ColumnType_JSON, Description: "The total amount before discount."},
			{Name: "type", Type: proto.ColumnType_STRING, Description: "The type of the invoice."},
			{Name: "state", Type: proto.ColumnType_STRING, Description: "The current state of the invoice."},
			{Name: "number", Type: proto.ColumnType_INT, Description: "The invoice number."},
			{Name: "seller_name", Type: proto.ColumnType_STRING, Description: "The name of the seller."},
		},
	}
}

type InvoiceItem struct {
	ID               string              `json:"id"`
	OrganizationID   string              `json:"organization_id"`
	OrganizationName string              `json:"organization_name"`
	StartDate        *time.Time          `json:"start_date"`
	StopDate         *time.Time          `json:"stop_date"`
	IssuedDate       *time.Time          `json:"issued_date"`
	DueDate          *time.Time          `json:"due_date"`
	TotalUntaxed     *scw.Money          `json:"total_untaxed"`
	TotalTaxed       *scw.Money          `json:"total_taxed"`
	TotalDiscount    *scw.Money          `json:"total_discount"`
	TotalUndiscount  *scw.Money          `json:"total_undiscount"`
	TotalTax         *scw.Money          `json:"total_tax"`
	Type             billing.InvoiceType `json:"type"`
	Number           int32               `json:"number"`
	State            string              `json:"state"`
	BillingPeriod    *time.Time          `json:"billing_period"`
	SellerName       string              `json:"seller_name"`
}

func listScalewayInvoices(ctx context.Context, d *plugin.QueryData, _ *plugin.HydrateData) (interface{}, error) {
	// Get client configuration
	client, err := getSessionConfig(ctx, d)
	if err != nil {
		plugin.Logger(ctx).Error("scaleway_invoice.listScalewayInvoices", "connection_error", err)
		return nil, err
	}

	// Check if the client is properly configured
	if client == nil {
		return nil, fmt.Errorf("scaleway client is not properly configured")
	}

	billingAPI := billing.NewAPI(client)

	// Prepare the request
	req := &billing.ListInvoicesRequest{}

	// Get the organization_id from the config
	scalewayConfig := GetConfig(d.Connection)
	var organizationID string
	if scalewayConfig.OrganizationID != nil {
		organizationID = *scalewayConfig.OrganizationID
	}

	// Check if organization_id is specified in the query
	if d.EqualsQualString("organization_id") != "" {
		organizationID = d.EqualsQualString("organization_id")
	}

	// Set the organization_id in the request if it's available
	if organizationID != "" {
		req.OrganizationID = &organizationID
	}

	// Make the API request to list invoices
	resp, err := billingAPI.ListInvoices(req)
	if err != nil {
		plugin.Logger(ctx).Error("scaleway_invoice.listScalewayInvoices", "api_error", err)
		return nil, err
	}

	for _, invoice := range resp.Invoices {
		item := InvoiceItem{
			ID:               invoice.ID,
			OrganizationID:   invoice.OrganizationID,
			OrganizationName: invoice.OrganizationName,
			StartDate:        invoice.StartDate,
			StopDate:         invoice.StopDate,
			IssuedDate:       invoice.IssuedDate,
			DueDate:          invoice.DueDate,
			TotalUntaxed:     invoice.TotalUntaxed,
			TotalTaxed:       invoice.TotalTaxed,
			TotalDiscount:    invoice.TotalDiscount,
			TotalUndiscount:  invoice.TotalUndiscount,
			TotalTax:         invoice.TotalTax,
			Type:             invoice.Type,
			Number:           invoice.Number,
			State:            invoice.State,
			BillingPeriod:    invoice.BillingPeriod,
			SellerName:       invoice.SellerName,
		}

		plugin.Logger(ctx).Debug("scaleway_invoice.listScalewayInvoices",
			"invoice_item", item,
		)

		d.StreamListItem(ctx, item)
	}

	return nil, nil
}
