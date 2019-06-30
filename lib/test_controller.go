package lib

import (
	"fmt"
	"sthub/lib/battle"
	"time"

	"github.com/labstack/echo/v4"
)

type ActiveBattle struct {
	file   *TestIterationFile
	battle *battle.Battle
}

// TestController is a controller that holds info about all testing iterations
type TestController struct {
	files map[string]*TestIterationFile

	activeBattle *ActiveBattle
}

// NewTestController creates a new TestController
func NewTestController() *TestController {
	return &TestController{files: make(map[string]*TestIterationFile, 0)}
}

func (t *TestController) getFileFromContext(c echo.Context) (*TestIterationFile, error) {
	iterationName := fmt.Sprintf("%s-%s", c.Param("wowsVersion"), c.Param("iteration"))
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
	g.GET("/:wowsVersion/:iteration", t.GetIteration)
	g.GET("/:wowsVersion/:iteration/battles", t.GetBattles)
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
	ShipName   string
	InDivision bool
}

func (t *TestController) StartBattle(c echo.Context) error {
	file, err := t.getFileFromContext(c)
	if err != nil {
		return err
	}

	req := new(StartBattleRequest)
	if err := c.Bind(req); err != nil {
		return err
	}

	now := time.Now()

	battle := &battle.Battle{
		StartedAt: &now,
		Ship:      req.ShipID,
		ShipName:  req.ShipName,
		Status:    "active",
		Statistics: battle.BattleStatistics{
			InDivision: battle.CorrectableBool{Value: req.InDivision},
		},
	}

	t.activeBattle = &ActiveBattle{
		file:   file,
		battle: battle,
	}

	file.Battles = append(file.Battles, battle)

	c.JSON(200, battle)
	return nil
}
