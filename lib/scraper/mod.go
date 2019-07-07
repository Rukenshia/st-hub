package scraper

// ModBattleInfo represents information passed from the python game mod to scraper
type ModBattleInfo struct {
	Status string
	Timestamp string
	ShipID string
	InDivision bool

	// Battle end info
	Win    bool
	Damage uint64
	Survived bool
	Kills uint64
}

// IsBattleQuit returns whether the mod battle information hints to a started battle
func (m *ModBattleInfo) IsBattleStart() bool {
	return m.Status === "active"
}

// IsBattleQuit returns whether the mod battle information hints to a completed battle
func (m *ModBattleInfo) IsBattleEnd() bool {
	return m.Status === "finished"
}

// IsBattleQuit returns whether the mod battle information hints to an abandoned battle
func (m *ModBattleInfo) IsBattleQuit() bool {
	return m.Status === "abandoned"
}

// ToBattleStartRequest transforms the ModBattleInfo into a Battle start request
func (m *ModBattleInfo) ToBattleStartRequest() *lib.StartBattleRequest {
	return &lib.StartBattleRequest{
		ShipID: m.ShipID,
		InDivision: m.InDivision,
	}
}

// EnhanceBattle modifies the given battle statistics with information from the ModBattleInfo
// This is only done if the battle and mod info have the same timestamp.
func (m *ModBattleInfo) EnhanceBattle(b *battle.Battle) (*battle.Battle, error) {
	if b.Timestamp != m.Timestamp {
		return nil, fmt.Errorf("Timestamp mismatch: mod has %s, battle has %s", m.Timestamp, b.Timestamp)
	}

	b.InDivision.Value = m.InDivision
	b.Statistics.Win = m.Win
	b.Statistics.Survived = m.Survived
	b.Statistics.Damage.Value = m.Damage
	b.Statistics.Kills.Value = m.Kills

	return b, nil
}