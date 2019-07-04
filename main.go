package main

import (
	"log"
	"net/http"
	"sthub/lib"
	"time"

	"github.com/levigross/grequests"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	//	"gopkg.in/toast.v1"
)

func main() {
	// Find current test iteration
	res, err := grequests.Get("https://hv59yay1u3.execute-api.eu-central-1.amazonaws.com/live/iteration/current", nil)
	if err != nil {
		//notification := toast.Notification{
		//AppID:   "{1AC14E77-02E7-4E5D-B744-2EB1AE5198B7}\\WindowsPowerShell\\v1.0\\powershell.exe",
		//Title:   "StHub Error (RECV_CURRENT_ITER)",
		//Message: "Could not get information on the current test iteration. Please contact Rukenshia",
		//}
		//notification.Push()
		log.Fatalln(err)
	}

	currentIteration := new(lib.TestIteration)
	if err := res.JSON(currentIteration); err != nil {
		//notification := toast.Notification{
		//AppID:   "{1AC14E77-02E7-4E5D-B744-2EB1AE5198B7}\\WindowsPowerShell\\v1.0\\powershell.exe",
		//Title:   "StHub Error (PARSE_CURRENT_ITER)",
		//Message: "Could not get information on the current test iteration. Please contact Rukenshia",
		//}
		//notification.Push()
		log.Fatalln(err)
	}

	testController, err := lib.NewTestController(currentIteration)
	if err != nil {
		//notification := toast.Notification{
		//AppID:   "{1AC14E77-02E7-4E5D-B744-2EB1AE5198B7}\\WindowsPowerShell\\v1.0\\powershell.exe",
		//Title:   "StHub Error (LC_CURRENT_ITER)",
		//Message: "Could not get information on the current test iteration. Please contact Rukenshia",
		//}
		//notification.Push()
		log.Fatalln(err)
	}

	// Echo instance
	e := echo.New()

	// Middleware
	e.Use(middleware.Recover())
	e.Use(middleware.Logger())
	e.Use(middleware.CORSWithConfig(middleware.CORSConfig{
		AllowOrigins: []string{"http://localhost:5000", "http://100.115.92.205:5000"},
		AllowMethods: []string{http.MethodGet, http.MethodPut, http.MethodPost, http.MethodDelete},
	}))

	// Frontend
	e.Static("/frontend", "frontend/public")

	testController.RegisterRoutes(e.Group("/iterations"))

	go func() {
		// Wait and open the browser
		timer := time.NewTimer(100 * time.Millisecond)

		<-timer.C

		// if err := exec.Command("rundll32", "url.dll,FileProtocolHandler", "http://localhost:1323/frontend/index.html").Start(); err != nil {
		// 	log.Fatal("Could not open browser")
		// }
	}()

	// Start server
	e.Logger.Fatal(e.Start(":1323"))
}
