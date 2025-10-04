package database

import "encoding/json"

func GetInventoryItems(boxId int) ([]InventoryItemWithProjectCount, error) {
	items, err := Select[InventoryItemWithProjectCount]("select * from inventory_item_with_count where box_id = $1 order by name", boxId)
	if err != nil {
		return nil, err
	}

	for i, item := range items {
		item.fillProperties()
		items[i] = item
	}

	return items, nil
}

func GetInventoryItem(id, boxId int) (*InventoryItemWithProjectCount, error) {
	item, err := SelectOne[InventoryItemWithProjectCount]("select * from inventory_item_with_count where box_id = $1 and id = $2", boxId, id)
	if err != nil {
		return nil, err
	}

	item.fillProperties()

	return &item, nil
}

func CreateInventoryItem(inventoryItem InventoryItem) (*InventoryItemWithProjectCount, error) {
	tx, err := dbMap.Begin()
	if err != nil {
		return nil, err
	}

	err = tx.Insert(&inventoryItem)
	if err != nil {
		_ = tx.Rollback()
		return nil, err
	}

	for key, value := range inventoryItem.Properties {
		prop := InventoryItemProperty{
			InventoryItemId: inventoryItem.Id,
			Name:            key,
			Value:           value,
		}
		err = tx.Insert(&prop)
		if err != nil {
			_ = tx.Rollback()
			return nil, err
		}
	}

	err = tx.Commit()
	if err != nil {
		return nil, err
	}

	return GetInventoryItem(inventoryItem.Id, inventoryItem.BoxId)
}

func UpdateInventoryItem(inventoryItem InventoryItem) error {
	tx, err := dbMap.Begin()
	if err != nil {
		return err
	}

	_, err = tx.Update(&inventoryItem)
	if err != nil {
		_ = tx.Rollback()
		return err
	}

	if inventoryItem.Properties != nil && len(inventoryItem.Properties) != 0 {
		props, err := json.Marshal(inventoryItem.Properties)
		if err != nil {
			_ = tx.Rollback()
			return err
		}

		//language=postgresql
		_, err = tx.Exec(`
with incoming(data) as (values ($2::jsonb)),
     parsed as (select $1::int as inventory_item_id,
                       key,
                       value
                from incoming, jsonb_each_text(data)),
     upserted as (
         insert into inventory_item_property (inventory_item_id, name, value)
             select inventory_item_id, key as name, value from parsed
             on conflict (inventory_item_id, name) do update
                 set value = excluded.value
             returning inventory_item_id, name)
delete
from inventory_item_property
where inventory_item_id = $1
  and name not in (select key from parsed)
`, inventoryItem.Id, props)
	} else if len(inventoryItem.Properties) == 0 {
		_, err = tx.Exec("delete from inventory_item_property where inventory_item_id = $1", inventoryItem.Id)
		if err != nil {
			_ = tx.Rollback()
			return err
		}
	}

	return tx.Commit()
}

func DeleteInventoryItem(id, boxId int) error {
	_, err := dbMap.Exec("delete from inventory_item where id = $1 and box_id = $2", id, boxId)
	return err
}
