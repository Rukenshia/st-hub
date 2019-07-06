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

func (t *TestController) RegisterRoutes(g *echo.Group) {
	g.GET("/:clientVersion/:iteration", t.GetIteration)
	g.GET("/:clientVersion/:iteration/battles", t.GetBattles)
	g.PUT("/:clientVersion/:iteration/battles/:battle", t.UpdateBattle)

	g.GET("/current", t.GetCurrentIteration)
	g.POST("/current/battles", t.StartBattle)
	g.POST("/current/battles/active", t.UpdateActiveBattle)
}

func (t *TestController) GetCurrentIteration(c echo.Context) error {
	if t.currentIteration == nil {
		c.String(404, "No active iteration")
		return nil
	}
	c.JSON(200, t.currentIteration)

	return nil
}

func (t *TestController) GetIteration(c echo.Context) error {
	file, err := t.getFileFromContext(c)
	if err != nil {
		return err
	}

	c.JSON(200, file)

	return nil
}

func (t *TestController) GetBattles(c echo.Context) error {
	file, err := t.getFileFromContext(c)
	if err != nil {
		return err
	}

	c.JSON(200, file.Battles)
	return nil
}

type StartBattleRequest struct {
	ShipID     uint64
	InDivision bool
}

func (t *TestController) StartBattle(c echo.Context) error {
	if t.currentIteration == nil {
		c.String(400, "No current iteration")
		return nil
	}

	if t.activeBattle != nil {
		c.Logger().Debug("StartBattle: cannot start because of an active battle")
		c.String(400, "There is already an active battle")
		return nil
	}

	req := new(StartBattleRequest)
	if err := c.Bind(req); err != nil {
		return err
	}

	if !t.currentIteration.HasShip(req.ShipID) {
		c.String(400, "Ship not part of testing iteration")
		return nil
	}

	ship := t.currentIteration.GetShip(req.ShipID)

	now := time.Now()

	battle := &battle.Battle{
		ID:        xid.New().String(),
		StartedAt: &now,
		ShipID:    ship.ID,
		ShipName:  ship.Name,
		Status:    "active",
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

	log.Printf(">>> %v", req.Statistics.InDivision.Value)

	// Check that unchangeable fields have not changed
	if req.StartedAt.String() != file.Battles[index].StartedAt.String() {
		c.String(400, "Can not change battle start time")
		return nil
	}

	file.Battles[index] = req

	return c.JSON(200, req)
}

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
		c.String(400, "Can not change battle start time")
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

	if req.Status == "finished" {
		now := time.Now()
		t.activeBattle.battle.FinishedAt = &now
		t.activeBattle = nil
	}

	return c.JSON(200, req)
}
