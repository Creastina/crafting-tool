use sea_orm_migration::{prelude::*, schema::*};

#[derive(DeriveMigrationName)]
pub struct Migration;

#[async_trait::async_trait]
impl MigrationTrait for Migration {
    async fn up(&self, manager: &SchemaManager) -> Result<(), DbErr> {
        manager
            .create_table(
                Table::create()
                    .table(Todos::Table)
                    .if_not_exists()
                    .col(pk_auto(Todos::Id))
                    .col(string(Todos::Title))
                    .col(string(Todos::Status))
                    .col(string(Todos::Kind))
                    .col(string(Todos::Material))
                    .col(boolean(Todos::IsDone).default(false))
                    .col(boolean(Todos::IsPartsMissing).default(false))
                    .col(string(Todos::Notes))
                    .to_owned(),
            )
            .await
    }

    async fn down(&self, manager: &SchemaManager) -> Result<(), DbErr> {
        manager
            .drop_table(Table::drop().table(Todos::Table).to_owned())
            .await
    }
}

#[derive(DeriveIden)]
enum Todos {
    Table,
    Id,
    Title,
    Status,
    Kind,
    Material,
    IsDone,
    IsPartsMissing,
    Notes,
}
