package battle

import "time"

// Battle is a match played in the game
type Battle struct {
	FinishedAt *time.Time
	StartedAt  *time.Time
	Status     string
	Ship       uint64
	ShipName   string
	Statistics BattleStatistics
}

// BattleStatistics contains basic statistics used by the frontend
type BattleStatistics struct {
	InDivision CorrectableBool
	Win        bool
	Survived   bool
	Damage     CorrectableUInt
	Kills      CorrectableUInt
}

// CorrectableUInt is a container for a value that can be manually corrected
type CorrectableUInt struct {
	Value     uint64
	Corrected *uint64
}

// CorrectableBool is a container for a value that can be manually corrected
type CorrectableBool struct {
	Value     bool
	Corrected *bool
}
