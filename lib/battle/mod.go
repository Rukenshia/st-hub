package battle

import "time"

// Battle is a match played in the game
type Battle struct {
	ID         string // xid of the battle
	FinishedAt *time.Time
	StartedAt  *time.Time
	Timestamp  string // Timestamp from the game mod
	Status     string
	ShipID     uint64
	ShipName   string
	Statistics Statistics
	Results    *Results `json:",omitempty"`
}

// Statistics contains basic statistics used by the frontend
type Statistics struct {
	InDivision CorrectableBool
	Win        bool
	Survived   bool
	Damage     CorrectableUInt
	Kills      CorrectableUInt
}

// AmmunitionResult contains information on damage, shots and hits for a specific type of ingame ammo
type AmmunitionResult struct {
	Damage float64
	Shots  float64
	Hits   float64
}

// Results contains detailed information from the WoWS battle result screen
type Results struct {
	Timestamp    string
	TeamID       uint64
	WinnerTeamID uint64
	BattleType   string
	Duration     float64
	PlaceInTeam  uint64

	Damage struct {
		Sum      uint64
		Fire     uint64
		Flooding uint64
		Ramming  uint64
	}
	Ammo struct {
		Torpedo        AmmunitionResult
		PlaneBomb      AmmunitionResult
		PlaneRocket    AmmunitionResult
		MainBatteryAP  AmmunitionResult
		MainBatterySAP AmmunitionResult
		MainBatteryHE  AmmunitionResult
		SecondaryAP    AmmunitionResult
		SecondarySAP   AmmunitionResult
		SecondaryHE    AmmunitionResult
	}
	FloodsCaused    uint64
	ShipsDetected   uint64
	LifeTime        float64
	PlanesKilled    uint64
	DistanceCovered float64

	Economics struct {
		Credits uint64
		BaseExp uint64
	}
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
