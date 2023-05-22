package main

import (
	"math/rand"
	"net/http"
	"rps_api/data"
	"time"

	"github.com/gin-gonic/gin"
)

// creating 3 data structures: slices
var winList = []string{
	"You should buy a lottery ticket.",
	"Great!",
	"Nice!",
}
var loseList = []string{
	"Too bad!",
	"Boohoo",
	"Opps",
}
var drawList = []string{
	"Great minds think alike",
	"Uh oh. Try again.",
	"Nobody wins, but you can try again.",
}

var config Config
var theRandom *rand.Rand

func start(c *gin.Context) {
	c.Data(http.StatusOK, "application/text", []byte("Tjena"))
}

func enableCors(c *gin.Context) {
	(*c).Header("Access-Control-Allow-Origin", "*")
}

func apiStats(c *gin.Context) {
	enableCors(c)
	totalGames, wins := data.Stats()
	c.JSON(http.StatusOK, gin.H{"totalGames": totalGames, "wins": wins})
}

func apiPlay(c *gin.Context) {
	enableCors(c)
	yourSelection := c.Query("yourSelection")
	computerSelection := randomizeSelection()
	winner := "Tie"
	val := theRandom.Intn(3) + 1
	message := drawList[val]
	if yourSelection == "STONE" && computerSelection == "SCISSOR" {
		winner = "You"
	}
	if yourSelection == "STONE" && computerSelection == "BAG" {
		winner = "Computer"
	}
	if yourSelection == "SCISSOR" && computerSelection == "BAG" {
		winner = "You"
	}
	if yourSelection == "SCISSOR" && computerSelection == "STONE" {
		winner = "Computer"
	}
	if yourSelection == "BAG" && computerSelection == "SCISSOR" {
		winner = "Computer"
	}
	if computerSelection == "BAG" && yourSelection == "STONE" {
		winner = "You"
	}

	if winner == "Computer" {
		message = winList[val]
	}
	if winner == "You" {
		message = loseList[val]
	}

	data.SaveGame(yourSelection, computerSelection, winner, message)
	c.JSON(http.StatusOK, gin.H{"winner": winner, "yourSelection": yourSelection, "computerSelection": computerSelection})
}

func randomizeSelection() string {
	val := theRandom.Intn(3) + 1
	if val == 1 {
		return "STONE"
	}
	if val == 2 {
		return "SCISSOR"
	}
	if val == 3 {
		return "BAG"
	}
	return "ERROR"

}

func main() {
	theRandom = rand.New(rand.NewSource(time.Now().UnixNano()))
	readConfig(&config)

	data.InitDatabase(config.Database.File,
		config.Database.Server,
		config.Database.Database,
		config.Database.Username,
		config.Database.Password,
		config.Database.Port)

	router := gin.Default()
	router.GET("/", start)
	router.GET("/api/play", apiPlay)
	router.GET("/api/stats", apiStats)
	// router.GET("/api/employee/:id", apiEmployeeById)
	// router.PUT("/api/employee/:id", apiEmployeeUpdateById)
	// router.DELETE("/api/employee/:id", apiEmployeeDeleteById)
	// router.POST("/api/employee", apiEmployeeAdd)

	// router.GET("/api/employees", employeesJson)
	// router.GET("/api/addemployee", addEmployee)
	// router.GET("/api/addmanyemployees", addManyEmployees)
	
	router.Run(":8080")

}
