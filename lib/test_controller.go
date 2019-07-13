package lib

import (
	"fmt"
	"log"
	"sthub/lib/battle"
	"time"

	"github.com/labstack/echo/v4"
	"github.com/rs/xid"
)

// ActiveBattle represents an ongoing battle in a TestIterationFile
type ActiveBattle struct {
	file   *TestIterationFile
	battle *battle.Battle
}

// TestController is a controller that holds info about all testing iterations
type TestController struct {
	files map[string]*TestIterationFile

	activeBattle     *ActiveBattle
	currentIteration *TestIterationFile
}

// NewTestController creates a new TestController
func NewTestController(currentIteration *TestIteration) (*TestController, error) {
	// Load or create file
	file, err := LoadOrCreateIterationFile(currentIteration)
	if err != nil {
		return nil, err
	}
	iterationName := fmt.Sprintf("%s-%s", currentIteration.ClientVersion, currentIteration.IterationName)

	var activeBattle *ActiveBattle
	for _, b := range file.Battles {
		if b.Status == "active" {
			activeBattle = &ActiveBattle{
				battle: b,
				file:   file,
			}
		}
	}

	return &TestController{files: map[string]*TestIterationFile{
		iterationName: file,
	}, currentIteration: file, activeBattle: activeBattle}, nil
}

// getFileFromContext loads or returns the TestIterationFile from the parameters in the request
func (t *TestController) getFileFromContext(c echo.Context) (*TestIterationFile, error) {
	iterationName := fmt.Sprintf("%s-%s", c.Param("clientVersion"), c.Param("iteration"))
	file, ok := t.files[iterationName]
	if !ok {
		lf, err := LoadTestIterationFile(iterationName)
		if err != nil {
			return nil, err
		}

		t.files[iterationName] = lf
		file = lf

		countActive := 0
		for _, b := range file.Battles {
			if b.Status == "active" {
				countActive++

				// Always take the last active battle
				t.activeBattle = &ActiveBattle{
					file:   lf,
					battle: b,
				}
			}
		}

		if countActive > 0 {
			c.JSON(400, map[string]string{
				"message": "Multiple active battles",
			})
		}
	}

	return file, nil
}

// RegisterRoutes is used to register routes of the testcontroller to echo
func (t *TestController) RegisterRoutes(g *echo.Group) {
	g.GET("/:clientVersion/:iteration", t.GetIteration)
	g.GET("/:clientVersion/:iteration/battles", t.GetBattles)
	g.PUT("/:clientVersion/:iteration/battles/:battle", t.UpdateBattle)

	// Called by game modification
	g.GET("/current", t.GetCurrentIteration)
	g.POST("/current/battles", t.StartBattle)
	g.POST("/current/battles/active", t.UpdateActiveBattle)
}

// GetCurrentIteration returns information about the active iteration
func (t *TestController) GetCurrentIteration(c echo.Context) error {
	if t.currentIteration == nil {
		c.String(404, "No active iteration")
		return nil
	}
	c.JSON(200, t.currentIteration)

	return nil
}

// GetIteration returns information on the given iteration
func (t *TestController) GetIteration(c echo.Context) error {
	file, err := t.getFileFromContext(c)
	if err != nil {
		return err
	}

	c.JSON(200, file)

	return nil
}

// GetBattles returns a list of battle in the given iteration
func (t *TestController) GetBattles(c echo.Context) error {
	file, err := t.getFileFromContext(c)
	if err != nil {
		return err
	}

	c.JSON(200, file.Battles)
	return nil
}

// StartBattleRequest represents the body parameters used to start a new active battle
type StartBattleRequest struct {
	ShipID     uint64
	InDivision bool
	Timestamp  string
}

// StartBattle starts a new active battle, usually called when the game entered a battle and the timer hit 0.
// Any active battle will be marked as "abandoned" to prevent out-of-sync state with the game modification.
// Abandoned battles are usually an indicator of the user restarting their client before a game was marked as
// ended or quit.
func (t *TestController) StartBattle(c echo.Context) error {
	if t.currentIteration == nil {
		c.String(400, "No current iteration")
		return nil
	}
	now := time.Now()

	if t.activeBattle != nil {
		// Abandon current battle
		t.activeBattle.battle.FinishedAt = &now
		t.activeBattle.battle.Status = "abandoned"
		t.activeBattle = nil
	}

	req := new(StartBattleRequest)
	if err := c.Bind(req); err != nil {
		return err
	}

	if !t.currentIteration.HasShip(req.ShipID) {
		c.String(400, "ERR_NOT_IN_TESTING")
		return nil
	}

	ship := t.currentIteration.GetShip(req.ShipID)

	battle := &battle.Battle{
		ID:        xid.New().String(),
		StartedAt: &now,
		ShipID:    ship.ID,
		ShipName:  ship.Name,
		Status:    "active",
		Timestamp: req.Timestamp,
		Statistics: battle.Statistics{
			InDivision: battle.CorrectableBool{Value: req.InDivision},
		},
	}

	t.activeBattle = &ActiveBattle{
		file:   t.currentIteration,
		battle: battle,
	}

	t.currentIteration.Battles = append(t.currentIteration.Battles, battle)
	t.currentIteration.Save()

	c.JSON(200, battle)
	return nil
}

// UpdateBattle updates a battle with the given body. Minimal checks are done
// to prevent updating fields that should not be updated.
func (t *TestController) UpdateBattle(c echo.Context) error {
	file, err := t.getFileFromContext(c)
	if err != nil {
		return err
	}

	req := new(battle.Battle)
	if err := c.Bind(req); err != nil {
		return err
	}

	// find battle
	index := -1
	for i, b := range file.Battles {
		if b.ID == req.ID {
			index = i
			break
		}
	}

	if index == -1 {
		c.String(404, "Not found")
		return nil
	}

	// Check that unchangeable fields have not changed
	if req.StartedAt.Format(time.UnixDate) != file.Battles[index].StartedAt.Format(time.UnixDate) {
		c.String(400, "ERR_CHANGE_START_TIME")
		return nil
	}
	if req.ShipID != file.Battles[index].ShipID || req.ShipName != file.Battles[index].ShipName {
		c.String(400, "ERR_CHANGE_SHIP")
		return nil
	}

	file.Battles[index] = req

	return c.JSON(200, req)
}

// UpdateActiveBattle updates the active battle with the request body. Minimal checks are done
// to verify that no static fields are changed.
func (t *TestController) UpdateActiveBattle(c echo.Context) error {
	if t.activeBattle == nil {
		c.String(404, "No active battle")
		return nil
	}

	req := new(battle.Battle)
	if err := c.Bind(req); err != nil {
		return err
	}

	if req.ID != t.activeBattle.battle.ID {
		c.String(400, "Can only modify active battle")
		return nil
	}

	// Check that unchangeable fields have not changed
	if req.StartedAt.Format(time.UnixDate) != t.activeBattle.battle.StartedAt.Format(time.UnixDate) {
		c.String(400, "ERR_CHANGE_START_TIME")
		return nil
	}
	if req.ShipID != t.activeBattle.battle.ShipID || req.ShipName != t.activeBattle.battle.ShipName {
		c.String(400, "ERR_CHANGE_SHIP")
		return nil
	}

	index := -1
	for i, b := range t.activeBattle.file.Battles {
		if b == t.activeBattle.battle {
			index = i
			break
		}
	}

	if index == -1 {
		c.String(500, "Internal error getting active battle")
		return nil
	}

	*t.activeBattle.battle = *req
	t.activeBattle.file.Save()

	if req.Status != "active" {
		log.Printf("removing active battle")
		now := time.Now()
		t.activeBattle.battle.FinishedAt = &now
		t.activeBattle = nil
	}

	return c.JSON(200, req)
}

// GetActiveBattle returns the current active battle
func (t *TestController) GetActiveBattle() *battle.Battle {
	if t.activeBattle == nil {
		return nil
	}
	return t.activeBattle.battle
}
