use serde::{Serialize, Deserialize};

pub mod actor;

#[derive(Serialize, Deserialize)]
pub struct Battle {
    id: String,
    time: i64,
    statistics: BattleStatistics,
}

#[derive(Serialize, Deserialize)]
pub struct BattleStatistics {
    pub win: bool,
    pub damage: u32,
    pub kills: u8,
    pub hits: u32,
    pub division: bool,
}