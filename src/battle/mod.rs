pub mod actor;

pub struct Battle {
    id: String,
    time: i64,
    statistics: BattleStatistics,
}

pub struct BattleStatistics {
    win: bool,
    damage: u32,
    kills: u8,
    hits: u32,
    division: bool,
}