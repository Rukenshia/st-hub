package scraper

import (
	"sthub/lib"
	"sthub/lib/battle"
)

// ModBattleInfo represents information passed from the python game mod to scraper
type ModBattleInfo struct {
	Status     string
	Timestamp  string
	ShipID     uint64
	InDivision bool

	// Battle end info
	Win      bool
	Damage   uint64
	Survived bool
	Kills    uint64
}

// IsBattleStart returns whether the mod battle information hints to a started battle
func (m *ModBattleInfo) IsBattleStart() bool {
	return m.Status == "active"
}

// IsBattleEnd returns whether the mod battle information hints to a completed battle
func (m *ModBattleInfo) IsBattleEnd() bool {
	return m.Status == "finished"
}

// IsBattleQuit returns whether the mod battle information hints to an abandoned battle
func (m *ModBattleInfo) IsBattleQuit() bool {
	return m.Status == "abandoned"
}

// ToBattleStartRequest transforms the ModBattleInfo into a Battle start request
func (m *ModBattleInfo) ToBattleStartRequest() *lib.StartBattleRequest {
	return &lib.StartBattleRequest{
		ShipID:     m.ShipID,
		InDivision: m.InDivision,
		Timestamp:  m.Timestamp,
	}
}

// GetStatistics returns BattleStatistics from the given ModBattleInfo
func (m *ModBattleInfo) GetStatistics() battle.Statistics {
	return battle.Statistics{
		Damage: battle.CorrectableUInt{
			Value: m.Damage,
		},
		InDivision: battle.CorrectableBool{
			Value: m.InDivision,
		},
		Kills: battle.CorrectableUInt{
			Value: m.Kills,
		},
		Survived: m.Survived,
		Win:      m.Win,
	}
}
