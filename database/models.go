package database

type InventoryBox struct {
	Id   int    `db:"id,primarykey,autoincrement" json:"id"`
	Name string `db:"name,unique" json:"name"`
}

type InventoryItem struct {
	BoxId      int               `db:"box_id" json:"-"`
	Id         int               `db:"id,primarykey,autoincrement" json:"id"`
	Name       string            `db:"name,notnull" json:"name"`
	Note       string            `db:"note" json:"note"`
	Count      int               `db:"count,notnull" json:"count"`
	Unit       string            `db:"unit,notnull" json:"unit"`
	Properties map[string]string `db:"-" json:"properties"`
}

type InventoryItemWithProjectCount struct {
	InventoryItem
	ProjectCount int `db:"project_count" json:"projectCount"`
}

type InventoryItemWithBoxName struct {
	InventoryItem
	BoxName string `db:"box_name" json:"boxName"`
}

func (i *InventoryItem) fillProperties() {
	i.Properties = make(map[string]string)

	props, err := Select[InventoryItemProperty]("select * from inventory_item_property where inventory_item_id = $1", i.Id)
	if err != nil {
		return
	}

	for _, prop := range props {
		i.Properties[prop.Name] = prop.Value
	}
}

type InventoryItemProperty struct {
	InventoryItemId int    `db:"inventory_item_id" json:"-"`
	Id              int    `db:"id,primarykey,autoincrement" json:"id"`
	Name            string `db:"name,notnull" json:"name"`
	Value           string `db:"value,notnull" json:"value"`
}

type ProjectCategory struct {
	Id   int    `db:"id,primarykey,autoincrement" json:"id"`
	Name string `db:"name,unique" json:"name"`
}

type Project struct {
	CategoryId     int             `db:"category_id" json:"-"`
	Id             int             `db:"id,primarykey,autoincrement" json:"id"`
	Name           string          `db:"name,notnull" json:"name"`
	Note           string          `db:"note" json:"note"`
	IsArchived     bool            `db:"is_archived" json:"isArchived"`
	InventoryItems []InventoryItem `db:"-" json:"inventoryItems"`
}

func (p *Project) fillInventoryItems() {
	inventoryItems, err := Select[InventoryItem](`
select ii.*
from inventory_item ii
         inner join public.project_inventory_item pii on ii.id = pii.inventory_item_id
where pii.project_id = $1
`, p.Id)
	if err != nil {
		return
	}

	p.InventoryItems = inventoryItems
}

type ProjectInventoryItem struct {
	ProjectId       int `db:"project_id" json:"-"`
	InventoryItemId int `db:"inventory_item_id" json:"-"`
}

type Instruction struct {
	Id   int    `db:"id,primarykey,autoincrement" json:"id"`
	Name string `db:"name,notnull" json:"name"`
	Note string `db:"note" json:"note"`
}

type InstructionWithStepCount struct {
	Id             int    `db:"id,primarykey,autoincrement" json:"id"`
	Name           string `db:"name,notnull" json:"name"`
	Note           string `db:"note" json:"note"`
	DoneStepCount  int    `db:"done_step_count" json:"doneStepCount"`
	TotalStepCount int    `db:"total_step_count" json:"totalStepCount"`
}

type InstructionStep struct {
	InstructionId int    `db:"instruction_id" json:"-"`
	Id            int    `db:"id,primarykey,autoincrement" json:"id"`
	Description   string `db:"description,notnull" json:"description"`
	Done          bool   `db:"done" json:"done"`
	Position      int    `db:"position" json:"-"`
}
