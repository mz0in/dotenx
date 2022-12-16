package marketplaceStore

import (
	"context"
	"database/sql"
	"encoding/json"
	"errors"
	"fmt"

	"github.com/dotenx/dotenx/ao-api/db"
	"github.com/dotenx/dotenx/ao-api/models"
	"github.com/sirupsen/logrus"
)

func (store *marketplaceStore) ListItems(ctx context.Context, accountId, category, itemType string, enabled bool) ([]models.MarketplaceItem, error) {

	listItems := "SELECT * FROM marketplace_items WHERE "

	args := []interface{}{}
	i := 1
	if accountId != "" {
		listItems += fmt.Sprintf("creator_account_id = $%d and ", i)
		i++
		args = append(args, accountId)
	}
	if category != "" {
		listItems += fmt.Sprintf("category = $%d and ", i)
		i++
		args = append(args, category)
	}
	if itemType != "" {
		listItems += fmt.Sprintf("item_type = $%d and ", i)
		i++
		args = append(args, itemType)
	}
	if enabled {
		listItems += "enabled = true"
	} else {
		listItems += "enabled = false"
	}

	var stmt string
	switch store.db.Driver {
	case db.Postgres:
		stmt = listItems
	default:
		return []models.MarketplaceItem{}, fmt.Errorf("driver not supported")
	}

	var items []models.MarketplaceItem
	rows, err := store.db.Connection.Queryx(stmt, args...)
	if err != nil {
		if err == sql.ErrNoRows {
			err = errors.New("item not found")
		}
		logrus.Error(err.Error())
		return []models.MarketplaceItem{}, err
	}

	for rows.Next() {

		var item models.MarketplaceItem
		err := rows.StructScan(&item)
		if err != nil {
			return []models.MarketplaceItem{}, err
		}

		var f []models.MarketplaceItemFeature
		err = json.Unmarshal(item.FeaturesForDb, &f)
		if err != nil {
			return []models.MarketplaceItem{}, err
		}
		item.Features = f

		if item.UpdatedAtForDb.Valid {
			item.UpdatedAt = item.UpdatedAtForDb.Time
		}

		items = append(items, item)
	}
	return items, nil
}
