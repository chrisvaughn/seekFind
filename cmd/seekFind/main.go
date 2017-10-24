package main

import (
	"fmt"
	"math/rand"
	"time"

	"github.com/chrisvaughn/seekFind/pkg/game"
)

func main() {
	rand.Seed(time.Now().UTC().UnixNano())
	wordList, err := game.ReadWordList("wordlist.txt")
	if err != nil {
		fmt.Println("Error:", err.Error())
		return
	}
	fmt.Println(wordList)
	boardSize := 25
	board := game.BuildBoard(boardSize, wordList)
	if board != nil {
		game.PrintBoard(board)
		game.HTMLBoard(board, wordList, "output.html")
		fmt.Println("output.html saved with board size:", boardSize)
	} else {
		fmt.Println("Error: couldn't fit words")
	}
}
