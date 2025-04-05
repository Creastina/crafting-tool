#[cfg(not(target_arch = "wasm32"))]
use sea_orm::entity::prelude::*;
use serde::{Deserialize, Serialize};

#[derive(Clone, Debug, PartialEq, Eq, Serialize, Deserialize, Default)]
#[cfg_attr(
    not(target_arch = "wasm32"),
    derive(DeriveEntityModel),
    sea_orm(table_name = "todos")
)]
pub struct Model {
    #[cfg_attr(
        not(target_arch = "wasm32"),
        sea_orm(primary_key, auto_increment = false)
    )]
    pub id: i32,
    pub title: String,
    pub status: String,
    pub kind: String,
    pub material: String,
    #[serde(default)]
    pub is_done: bool,
    #[serde(default)]
    pub is_parts_missing: bool,
    pub notes: String,
    #[cfg_attr(not(target_arch = "wasm32"), sea_orm(ignore))]
    #[serde(default)]
    pub is_new: bool,
}

#[cfg(not(target_arch = "wasm32"))]
#[derive(Copy, Clone, Debug, DeriveRelation, EnumIter)]
pub enum Relation {}

#[cfg(not(target_arch = "wasm32"))]
impl ActiveModelBehavior for ActiveModel {}

pub use Model as Todo;