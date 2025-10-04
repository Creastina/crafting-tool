package database

import (
	"context"
	"crafting/config"

	"github.com/DerKnerd/gorp"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/jackc/pgx/v5/stdlib"
)

var dbMap *gorp.DbMap

func GetDbMap() *gorp.DbMap {
	return dbMap
}

func SetupDatabase() {
	if dbMap == nil {
		pool, err := pgxpool.New(context.Background(), config.LoadedConfiguration.DatabaseUrl)
		if err != nil {
			panic(err)
		}

		conn := stdlib.OpenDBFromPool(pool)

		dialect := gorp.PostgresDialect{}

		dbMap = &gorp.DbMap{Db: conn, Dialect: dialect}
		AddTableWithName[InventoryBox]("inventory_box")
		AddTableWithName[InventoryItem]("inventory_item").
			SetUniqueTogether("box_id", "name")
		AddTableWithName[InventoryItemProperty]("inventory_item_property").
			SetUniqueTogether("inventory_item_id", "name")

		AddTableWithName[ProjectCategory]("project_category")
		AddTableWithName[Project]("project").
			SetUniqueTogether("category_id", "name")
		AddTableWithName[ProjectInventoryItem]("project_inventory_item")

		AddTableWithName[Instruction]("instruction")
		AddTableWithName[InstructionStep]("instruction_step")

		err = GetDbMap().CreateTablesIfNotExists()
		if err != nil {
			panic(err)
		}

		AddTableWithName[InventoryItemWithProjectCount]("inventory_item_with_count")

		_, err = conn.Exec(`
create or replace function add_foreign_key_if_not_exists(from_table text, from_column text, to_table text, to_column text)
returns void language plpgsql as
$$
declare 
   fk_exists boolean;
begin
    fk_exists := case when exists (select true
	from information_schema.table_constraints tc
		inner join information_schema.constraint_column_usage ccu
			using (constraint_catalog, constraint_schema, constraint_name)
		inner join information_schema.key_column_usage kcu
			using (constraint_catalog, constraint_schema, constraint_name)
	where constraint_type = 'FOREIGN KEY'
	  and ccu.table_name = to_table
	  and ccu.column_name = to_column
	  and tc.table_name = from_table
	  and kcu.column_name = from_column) then true else false end;
	if not fk_exists then
		execute format('alter table %s add constraint %s_%s_fkey foreign key (%s) references %s(%s) on delete cascade', from_table, from_table, to_table, from_column, to_table, to_column);
	end if;
end
$$;
`)
		if err != nil {
			panic(err)
		}

		_, err = conn.Exec(`
select add_foreign_key_if_not_exists('inventory_item', 'box_id', 'inventory_box', 'id');
select add_foreign_key_if_not_exists('inventory_item_property', 'inventory_item_id', 'inventory_item', 'id');
select add_foreign_key_if_not_exists('project', 'category_id', 'project_category', 'id');
select add_foreign_key_if_not_exists('project_inventory_item', 'project_id', 'project', 'id');
select add_foreign_key_if_not_exists('project_inventory_item', 'inventory_item_id', 'inventory_item', 'id');
select add_foreign_key_if_not_exists('instruction_step', 'instruction_id', 'instruction', 'id')`)
		if err != nil {
			panic(err)
		}

		_, err = conn.Exec(`
drop table if exists todos;
drop table if exists seaql_migrations`)
		if err != nil {
			panic(err)
		}

		_, err = conn.Exec(`
create or replace view inventory_item_with_count as
select *,
       (select count(p.*)
        from project p
                 inner join project_inventory_item pii on p.id = pii.project_id
        where pii.inventory_item_id = i.id) as project_count
from inventory_item i;
`)
		if err != nil {
			panic(err)
		}
	}
}
