use crafting_entities::{Column, Entity, Todo};
use crafting_error::{CraftingError, CraftingErrorResult, CraftingResult};
use sea_orm::prelude::Expr;
use sea_orm::{
    ActiveModelTrait, ColumnTrait, ConnectOptions, Database, DatabaseConnection, DbErr,
    EntityTrait, IntoActiveModel, NotSet, QueryFilter, QueryOrder,
};

pub async fn get_database() -> Result<DatabaseConnection, DbErr> {
    let mut opts = ConnectOptions::new(std::env::var("DATABASE_URL").expect("Needs DATABASE_URL"));
    opts.sqlx_logging(true)
        .sqlx_logging_level(log::LevelFilter::Debug);

    Database::connect(opts).await
}

pub async fn get_todos(db: &DatabaseConnection) -> CraftingResult<Vec<Todo>> {
    Entity::find()
        .order_by_asc(Column::Title)
        .all(db)
        .await
        .map_err(|err| {
            log::error!("Failed to load todos {err}");
            CraftingError::database("Failed to load todos")
        })
}

pub async fn create_todo(todo: Todo, db: &DatabaseConnection) -> CraftingErrorResult {
    let mut todo = todo.into_active_model();
    todo.id = NotSet;
    todo.insert(db)
        .await
        .map_err(|err| {
            log::error!("Failed to create todo {err}");
            CraftingError::database("Failed to create todo")
        })
        .map(|_| ())
}

pub async fn update_todo(todo: Todo, db: &DatabaseConnection) -> CraftingErrorResult {
    Entity::update_many()
        .col_expr(Column::Status, Expr::value(todo.status))
        .col_expr(Column::Kind, Expr::value(todo.kind))
        .col_expr(Column::Material, Expr::value(todo.material))
        .col_expr(Column::IsDone, Expr::value(todo.is_done))
        .col_expr(Column::IsPartsMissing, Expr::value(todo.is_parts_missing))
        .col_expr(Column::Notes, Expr::value(todo.notes))
        .filter(Column::Id.eq(todo.id))
        .exec(db)
        .await
        .map_err(|err| {
            log::error!("Failed to update todo {err}");
            CraftingError::database("Failed to update todo")
        })
        .map(|_| ())
}

pub async fn delete_todo(id: i32, db: &DatabaseConnection) -> CraftingErrorResult {
    Entity::delete_by_id(id)
        .exec(db)
        .await
        .map_err(|err| {
            log::error!("Failed to delete todo {err}");
            CraftingError::database("Failed to delete todo")
        })
        .map(|_| ())
}
