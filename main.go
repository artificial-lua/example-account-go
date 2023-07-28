package main

import (
	"fmt"
	"log"

	"github.com/artificial-lua/example-account-go/auth"
	"github.com/artificial-lua/example-account-go/dbconnector"
	"github.com/labstack/echo"
)

const savedFileName string = "pages.csv"
const provideFileName string = "page.csv"

func handleHome(c echo.Context) error {
	log.Println("handleHome", c)
	return c.File("html/index.html")
}

func handleJoin(c echo.Context) error {
	newAccount, err := auth.MakeAccountObject(
		c.FormValue(("email")),
		c.FormValue(("password")),
		"",
		"",
		c.FormValue(("name")),
		c.FormValue(("birth")),
		c.FormValue(("gender")),
	)
	if err != nil {
		log.Println(err)
		return c.JSON(404, map[string]string{
			"message": "join failed",
			"reason":  err.Error(),
		})
	}
	newAccount.CryptoPassword("joinsalt")
	result, err := auth.CreateAcccount(newAccount)
	if err != nil {
		log.Println(err)
		return c.JSON(404, map[string]string{
			"message": "join failed",
			"reason":  err.Error(),
		})
	}
	fmt.Println(result)

	return c.Redirect(302, "/joinSuccess")
}

func handleJoinSuccess(c echo.Context) error {
	log.Println("handleJoinSuccess", c)
	return c.JSON(200, map[string]string{
		"message": "join success",
	})
}

func handleJoinFailed(c echo.Context) error {
	return c.JSON(200, map[string]string{
		"message": "join failed",
	})
}

func main() {
	db, err := dbconnector.NewPostgreSQLConnector()
	if err != nil {
		log.Fatal(err)
	}
	dbconnector.DBStartupTask(db)
	e := echo.New()
	e.GET("/", handleHome)
	e.POST("/join", handleJoin)
	e.GET("/joinSuccess", handleJoinSuccess)
	e.GET("/joinFailed", handleJoinFailed)
	e.Logger.Fatal(e.Start(":1323"))
}
