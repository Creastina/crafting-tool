use serde::{Deserialize, Serialize};
use std::error::Error;
use std::fmt::{Display, Formatter};
use std::str::FromStr;

#[derive(Serialize, Deserialize, Debug, Eq, PartialEq, Ord, PartialOrd, Clone, Copy)]
#[serde(rename_all = "camelCase")]
pub enum CraftingErrorCode {
    Database,
    NotFound,
    Unknown,
}

impl Default for CraftingErrorCode {
    fn default() -> Self {
        Self::Unknown
    }
}

impl FromStr for CraftingErrorCode {
    type Err = ();

    fn from_str(s: &str) -> Result<Self, Self::Err> {
        Ok(s.into())
    }
}

impl From<&str> for CraftingErrorCode {
    fn from(s: &str) -> Self {
        match s {
            "Database" => CraftingErrorCode::Database,
            "NotFound" => CraftingErrorCode::NotFound,
            _ => CraftingErrorCode::Unknown,
        }
    }
}

#[cfg(target_arch = "wasm32")]
impl Display for CraftingErrorCode {
    fn fmt(&self, f: &mut Formatter<'_>) -> std::fmt::Result {
        f.write_str(format!("{:?}", self).as_str())
    }
}

#[cfg(not(target_arch = "wasm32"))]
impl Display for CraftingErrorCode {
    fn fmt(&self, f: &mut Formatter<'_>) -> std::fmt::Result {
        f.write_str(serde_json::to_string(self).unwrap().as_str())
    }
}

#[derive(Ord, PartialOrd, Eq, PartialEq, Clone, Debug, Serialize, Deserialize, Default)]
#[serde(rename_all = "camelCase")]
pub struct CraftingError {
    pub error_type: CraftingErrorCode,
    pub message: String,
}

#[cfg(target_arch = "wasm32")]
impl Display for CraftingError {
    fn fmt(&self, f: &mut Formatter<'_>) -> std::fmt::Result {
        f.write_str(format!("{:?}", self).as_str())
    }
}

#[cfg(not(target_arch = "wasm32"))]
impl Display for CraftingError {
    fn fmt(&self, f: &mut Formatter<'_>) -> std::fmt::Result {
        f.write_str(serde_json::to_string(self).unwrap().as_str())
    }
}

impl Error for CraftingError {}

impl CraftingError {
    fn new(
        message: impl Into<String>,
        error_type: CraftingErrorCode,
    ) -> Self {
        Self {
            message: message.into(),
            error_type,
        }
    }

    pub fn database(message: impl Into<String>) -> Self {
        Self::new(message, CraftingErrorCode::Database)
    }

    pub fn not_found(message: impl Into<String>) -> Self {
        Self::new(message, CraftingErrorCode::NotFound)
    }

    pub fn unknown(message: impl Into<String>) -> Self {
        Self::new(message, CraftingErrorCode::Unknown)
    }
}

pub type CraftingErrorResult = Result<(), CraftingError>;

pub type CraftingResult<T> = Result<T, CraftingError>;
