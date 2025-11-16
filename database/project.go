package database

import "encoding/json"

func GetProjects(categoryId int) ([]Project, error) {
	projects, err := Select[Project]("select * from project where category_id = $1", categoryId)
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

	return project, nil
}

func CreateProject(project Project, inventoryItems map[int]int) (*Project, error) {
	tx, err := dbMap.Begin()
	if err != nil {
		return nil, err
	}

	err = tx.Insert(&project)
	if err != nil {
		_ = tx.Rollback()
		return nil, err
	}

	for inventoryItem, count := range inventoryItems {
		projectInventoryItem := ProjectInventoryItem{
			ProjectId:       project.Id,
			InventoryItemId: inventoryItem,
			Count:           count,
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

func UpdateProject(project Project, inventoryItems map[int]int) error {
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
     parsed as (
         select 
             (key)::int   as inventory_item_id,
             (value)::int as count
         from incoming, jsonb_each_text(data) as kv(key, value)
     ),
     upserted as (
         insert into project_inventory_item (project_id, inventory_item_id, count)
             select $1, inventory_item_id, count from parsed
             on conflict (inventory_item_id, project_id) do 
             update set count = excluded.count
             returning inventory_item_id)
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

func ArchiveCategory(categoryId int) error {
	_, err := dbMap.Exec("update project set is_archived = true where category_id = $1", categoryId)
	if err != nil {
		return err
	}

	//language=postgresql
	_, err = dbMap.Exec(`
with inventory as (
	select ii.id, sum(pii.count) as count
	from inventory_item ii
		join project_inventory_item pii on pii.inventory_item_id = ii.id
		join project p on p.id = pii.project_id
	where ii.count > 0 and p.category_id = 1
	group by ii.id
)
update inventory_item ii
set count = min(ii.count - inventory.count, 0)
from inventory
where ii.id = inventory.id`, categoryId)

	return err
}

func ArchiveProject(projectId, categoryId int) error {
	_, err := dbMap.Exec("update project set is_archived = true where id = $1 and category_id = $2", projectId, categoryId)
	if err != nil {
		return err
	}

	//language=postgresql
	_, err = dbMap.Exec(`
with inventory as (
    select ii.id, sum(pii.count) as count
	from inventory_item ii
    	join project_inventory_item pii on pii.inventory_item_id = ii.id
    where ii.count > 0 and pii.project_id = $1
    group by ii.id
)
update inventory_item ii
set count = ii.count - inventory.count
from inventory
where ii.id = inventory.id`, projectId)

	return err
}
