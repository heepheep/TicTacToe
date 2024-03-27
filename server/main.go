package main

import (
	"math/rand"
	"net/http"
	"strconv"
	"time"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

type GameBoard struct {
	GameId      int          `json:"id"`
	Player1     string       `json:"player1"`
	Player2     string       `json:"player2"`
	Turn        int          `json:"turn"`
	CurrMove    string       `json:"currMove"`
	Board       [3][3]string `json:"board"`
	Ans         string       `json:"ans"`
	IsFinnished bool         `json:"isFinnished"`
}
type GameSearch struct {
	SearchId int    `json:"id"`
	Player1  string `json:"player1"`
	IsBegun  bool   `json:"isBegun"`
}

var Games = []GameBoard{
	{GameId: 0, CurrMove: "X", Player1: "teh7v", Player2: "heepheep", Turn: 1, Board: [3][3]string{{"", "", ""}, {"", "", ""}, {"", "", ""}}, IsFinnished: false},
}
var Searches = []GameSearch{
	{SearchId: 0, Player1: "teh7v", IsBegun: true},
}

func getGameBoards(c *gin.Context) {
	c.IndentedJSON(http.StatusOK, Games)

}
func getGameById(c *gin.Context) {
	var gameId, _ = strconv.Atoi(c.Param("id"))
	for _, v := range Games {
		if v.GameId == gameId {
			c.IndentedJSON(http.StatusOK, Games[gameId])
			return
		}
	}
	c.IndentedJSON(220, gin.H{"INFO": "Game not found"})
}
func makeMove(c *gin.Context) {
	type Move struct {
		GameId     int    `json:"id"`
		PlayerName string `json:"player"`
		X          int    `json:"x"`
		Y          int    `json:"y"`
	}
	var myMove Move
	if err := c.ShouldBindJSON(&myMove); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"ERROR": err.Error()})
		return
	}
	for index, v := range Games {
		if v.GameId == myMove.GameId {
			if v.IsFinnished {
				c.JSON(240, gin.H{"INFO": "Игра закончена"})
				return
			}
			var moveType string
			if v.Player1 == myMove.PlayerName {
				if v.Turn == 0 {
					moveType = v.CurrMove
					if v.CurrMove == "X" {
						v.CurrMove = "O"
					} else {
						v.CurrMove = "X"
					}
				} else {
					c.JSON(230, gin.H{"INFO": "Не твой ход"})
					return
				}
			} else if v.Player2 == myMove.PlayerName {
				if v.Turn == 1 {
					moveType = v.CurrMove
					if v.CurrMove == "X" {
						v.CurrMove = "O"
					} else {
						v.CurrMove = "X"
					}
				} else {
					c.JSON(230, gin.H{"INFO": "Не твой ход"})
					return
				}
			} else {
				c.JSON(240, gin.H{"ERROR": "Отказано в доступе"})
				return
			}
			if v.Board[myMove.Y][myMove.X] == "" {
				v.Board[myMove.Y][myMove.X] = moveType
			} else {
				c.JSON(230, gin.H{"INFO": "Ячейка уже занята"})
				return
			}
			if v.Turn == 1 {
				v.Turn = 0
			} else {
				v.Turn = 1
			}

			type Answer struct {
				Board GameBoard `json:"board"`
				Ans   string    `json:"ans"`
			}
			var myAnswer Answer
			myAnswer.Ans = checkWinner(v.Board)
			if myAnswer.Ans == "" {
				if v.Turn == 1 {
					v.Ans = "Ход " + v.Player2
				} else if v.Turn == 0 {
					v.Ans = "Ход " + v.Player1
				}
			} else if myAnswer.Ans == "X" {
				v.Ans = "Победили крестики"
			} else if myAnswer.Ans == "O" {
				v.Ans = "Победили нолики"
			} else if myAnswer.Ans == "Ничья" {
				v.Ans = "Ничья"
			} else {
				v.Ans = "CheckApi"
			}
			if myAnswer.Ans != "" {
				v.IsFinnished = true
			}
			Games[index] = v
			myAnswer.Board = Games[index]
			c.IndentedJSON(http.StatusOK, myAnswer)
			return
		}
	}

}
func checkWinner(board [3][3]string) string {
	lines := [][3]int{
		{0, 1, 2}, {3, 4, 5}, {6, 7, 8}, // rows
		{0, 3, 6}, {1, 4, 7}, {2, 5, 8}, // columns
		{0, 4, 8}, {2, 4, 6}, // diagonals
	}

	for _, line := range lines {
		a, b, c := line[0], line[1], line[2]
		if board[a/3][a%3] != "" && board[a/3][a%3] == board[b/3][b%3] && board[a/3][a%3] == board[c/3][c%3] {
			return board[a/3][a%3]
		}
	}

	for _, row := range board {
		for _, cell := range row {
			if cell == "" {
				return "" // Continue game, empty cell found
			}
		}
	}

	return "Ничья"
}
func makeGame(c *gin.Context) {
	var player struct {
		PlayerName string `json:"playerName" binding:"required"`
	}
	if err := c.ShouldBindJSON(&player); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"ERROR": err.Error()})
	}
	for index, v := range Searches {
		if !v.IsBegun {
			Searches[index].IsBegun = true
			var newGame GameBoard
			newGame.Board = [3][3]string{{"", "", ""}, {"", "", ""}, {"", "", ""}}
			newGame.Player1 = v.Player1
			newGame.Player2 = player.PlayerName
			newGame.GameId = len(Games)
			newGame.IsFinnished = false
			newGame.Turn = rand.Intn(2)
			newGame.CurrMove = "X"
			if newGame.Turn == 1 {
				newGame.Ans = "Ход " + newGame.Player2
			} else if newGame.Turn == 0 {
				newGame.Ans = "Ход " + newGame.Player1
			}
			Games = append(Games, newGame)
			c.IndentedJSON(http.StatusOK, newGame)
			return
		}
	}
	print(Searches)
	var Game GameSearch
	Game.IsBegun = false
	Game.Player1 = player.PlayerName
	Game.SearchId = len(Searches)
	Searches = append(Searches, Game)
	c.IndentedJSON(http.StatusCreated, Game)
}

func main() {
	router := gin.Default()
	config := cors.DefaultConfig()
	config.AllowAllOrigins = true
	config.AllowMethods = []string{"POST", "GET", "PUT", "DELETE", "PATCH"}
	config.AllowHeaders = []string{"Origin", "Content-Type", "Authorization", "Accept", "User-Agent", "Cache-Control", "Pragma"}
	config.ExposeHeaders = []string{"Content-Length"}
	config.AllowCredentials = true
	config.MaxAge = 12 * time.Hour

	router.Use(cors.New(config))
	router.GET("/get-all-games", getGameBoards)
	router.GET("/get-game/:id", getGameById)
	router.POST("/make-game", makeGame)
	router.POST("/make-move", makeMove)
	router.Run("192.168.1.10:8080")

}
