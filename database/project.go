package database

import "encoding/json"

func GetProjects(categoryId int, isArchived bool) ([]Project, error) {
	projects, err := Select[Project]("select * from project where category_id = $1 and is_archived = $2", categoryId, isArchived)
	if err != nil {
		return nil, err
	}

	for i, project := range projects {
		project.fillInventoryItems()
		projects[i] = project
	}

	return projects, nil
}

func GetProject(projectId, categoryId int) (*Project, error) {
	project, err := SelectOne[Project]("select * from project where id = $1 and category_id = $2", projectId, categoryId)
	if err != nil {
		return nil, err
	}

	project.fillInventoryItems()

	return &project, nil
}

func CreateProject(project Project, inventoryItems []int) (*Project, error) {
	tx, err := dbMap.Begin()
	if err != nil {
		return nil, err
	}

	err = tx.Insert(&project)
	if err != nil {
		_ = tx.Rollback()
		return nil, err
	}

	for _, inventoryItem := range inventoryItems {
		projectInventoryItem := ProjectInventoryItem{
			ProjectId:       project.Id,
			InventoryItemId: inventoryItem,
		}
		err = tx.Insert(&projectInventoryItem)
		if err != nil {
			_ = tx.Rollback()
			return nil, err
		}
	}

	err = tx.Commit()
	if err != nil {
		return nil, err
	}

	return GetProject(project.Id, project.CategoryId)
}

func UpdateProject(project Project, inventoryItems []int) error {
	tx, err := dbMap.Begin()
	if err != nil {
		return err
	}

	_, err = tx.Update(&project)
	if err != nil {
		_ = tx.Rollback()
	}

	if inventoryItems != nil && len(inventoryItems) != 0 {
		items, err := json.Marshal(inventoryItems)
		if err != nil {
			_ = tx.Rollback()
			return err
		}

		//language=postgresql
		_, err = tx.Exec(`
with incoming(data) as (values ($2::jsonb)),
     parsed as (select $1::int as project_id, (jsonb_array_elements_text(data))::int as inventory_item_id
                from incoming),
     upserted as (
         insert into project_inventory_item (project_id, inventory_item_id)
             select project_id, inventory_item_id from parsed
             on conflict (inventory_item_id, project_id) do nothing
             returning inventory_item_id, project_id)
delete
from project_inventory_item
where project_id = $1
  and inventory_item_id not in (select inventory_item_id from parsed)
`, project.Id, items)
	} else if len(inventoryItems) == 0 {
		_, err = tx.Exec("delete from project_inventory_item where project_id = $1", project.Id)
		if err != nil {
			_ = tx.Rollback()
			return err
		}
	}

	return tx.Commit()
}

func DeleteProject(projectId, categoryId int) error {
	_, err := dbMap.Exec("delete from project where id = $1 and category_id = $2", projectId, categoryId)
	return err
}
