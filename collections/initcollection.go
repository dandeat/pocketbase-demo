package collections

import (
	"github.com/pocketbase/pocketbase/tools/types"

	"github.com/pocketbase/pocketbase"
	models "github.com/pocketbase/pocketbase/models"
	"github.com/pocketbase/pocketbase/models/schema"
)

func InitCollection(app *pocketbase.PocketBase) error {
	// Account
	accountCol := &models.Collection{
		Name:       "account",
		Type:       models.CollectionTypeBase,
		ListRule:   nil,
		ViewRule:   types.Pointer("@request.auth.id != ''"),
		CreateRule: types.Pointer(""),
		UpdateRule: types.Pointer("@request.auth.id != ''"),
		DeleteRule: nil,
		Schema: schema.NewSchema(
			&schema.SchemaField{
				Name:     "name",
				Type:     schema.FieldTypeText,
				Required: true,
			},
			&schema.SchemaField{
				Name:     "balance",
				Type:     schema.FieldTypeNumber,
				Required: true,
			},
			&schema.SchemaField{
				Name:     "created_at",
				Type:     schema.FieldTypeDate,
				Required: true,
			},
			&schema.SchemaField{
				Name:     "updated_at",
				Type:     schema.FieldTypeDate,
				Required: true,
			},
		),
	}

	// Transaction
	transactionCol := &models.Collection{
		Name:       "transaction",
		Type:       models.CollectionTypeBase,
		ListRule:   nil,
		ViewRule:   types.Pointer("@request.auth.id != ''"),
		CreateRule: types.Pointer(""),
		UpdateRule: types.Pointer("@request.auth.id != ''"),
		DeleteRule: nil,
		Schema: schema.NewSchema(
			&schema.SchemaField{
				Name:     "account_id",
				Type:     schema.FieldTypeText,
				Required: true,
			},
			&schema.SchemaField{
				Name:     "amount",
				Type:     schema.FieldTypeNumber,
				Required: true,
			},
			&schema.SchemaField{
				Name:     "transaction_type",
				Type:     schema.FieldTypeText,
				Required: true,
			},
			&schema.SchemaField{
				Name:     "transaction_number",
				Type:     schema.FieldTypeText,
				Required: true,
			},
			&schema.SchemaField{
				Name:     "description",
				Type:     schema.FieldTypeText,
				Required: true,
			},
			&schema.SchemaField{
				Name:     "created_at",
				Type:     schema.FieldTypeDate,
				Required: true,
			},
			&schema.SchemaField{
				Name:     "updated_at",
				Type:     schema.FieldTypeDate,
				Required: true,
			},
		),
	}

	// Register collections
	err := app.Dao().SaveCollection(accountCol)
	if err != nil {
		return err
	}

	err = app.Dao().SaveCollection(transactionCol)
	if err != nil {
		return err
	}
	return nil
}
