package main

import (
	"bufio"
	"fmt"
	"math/rand"
	"os"
	"strings"
	"time"
)

const (
	horizontal int = iota
	horizontalReversed
	vertical
	verticalReversed
	diagonalLTR
	diagonalLTRReversed
	diagonalRTL
	diagonalRTLReversed
)

const (
	letterBytes     = "abcdefghijklmnopqrstuvwxyz"
	fitWordAttempts = 1000
	boardAttempts   = 5000
)

type gameBoard [][]string

func readWordList(path string) (lines []string, err error) {
	file, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		lines = append(lines, strings.ToLower(scanner.Text()))
	}
	return lines, scanner.Err()
}

func makeBoard(x int, y int) (board gameBoard) {
	board = make(gameBoard, y)
	for i := range board {
		board[i] = make([]string, x)
	}
	return
}

func reverse(s string) string {
	chars := []rune(s)
	for i, j := 0, len(chars)-1; i < j; i, j = i+1, j-1 {
		chars[i], chars[j] = chars[j], chars[i]
	}
	return string(chars)
}

func fitHorizontal(board *gameBoard, word string) bool {
	x := len(*board)
	y := len((*board)[0])
	var ix, iy int
	found := false
	attempts := 0
	for !found {
		ix = rand.Intn(x)
		iy = rand.Intn(y)
		found = iy+len(word) < y
		attempts++
		if attempts > fitWordAttempts {
			return false
		}
	}

	for j, c := range strings.Split(word, "") {
		if (*board)[ix][iy+j] == "" {
			(*board)[ix][iy+j] = c
		} else {
			return false
		}
	}
	return true
}

func fitHorizontalReversed(board *gameBoard, word string) bool {
	return fitHorizontal(board, reverse(word))
}

func fitVertical(board *gameBoard, word string) bool {
	x := len(*board)
	y := len((*board)[0])
	var ix, iy int
	found := false
	attempts := 0
	for !found {
		ix = rand.Intn(x)
		iy = rand.Intn(y)
		found = ix+len(word) < x
		attempts++
		if attempts > fitWordAttempts {
			return false
		}
	}

	for j, c := range strings.Split(word, "") {
		if (*board)[ix+j][iy] == "" {
			(*board)[ix+j][iy] = c
		} else {
			return false
		}
	}
	return true
}

func fitVerticalReversed(board *gameBoard, word string) bool {
	return fitVertical(board, reverse(word))
}

func fitDiagonalLTR(board *gameBoard, word string) bool {
	x := len(*board)
	y := len((*board)[0])
	var ix, iy int
	found := false
	attempts := 0
	for !found {
		ix = rand.Intn(x)
		iy = rand.Intn(y)
		found = ix+len(word) < x && iy+len(word) < y
		attempts++
		if attempts > fitWordAttempts {
			return false
		}
	}

	for j, c := range strings.Split(word, "") {
		if (*board)[ix+j][iy+j] == "" {
			(*board)[ix+j][iy+j] = c
		} else {
			return false
		}
	}
	return true
}

func fitDiagonalLTRReversed(board *gameBoard, word string) bool {
	return fitDiagonalLTR(board, reverse(word))
}

func fitDiagonalRTL(board *gameBoard, word string) bool {
	x := len(*board)
	y := len((*board)[0])
	var ix, iy int
	found := false
	attempts := 0
	for !found {
		ix = rand.Intn(x)
		iy = len(word) + rand.Intn(y-len(word))
		found = ix+len(word) < x && iy-len(word) >= 0
		attempts++
		if attempts > fitWordAttempts {
			return false
		}
	}

	for j, c := range strings.Split(word, "") {
		if (*board)[ix+j][iy-j] == "" {
			(*board)[ix+j][iy-j] = c
		} else {
			return false
		}
	}
	return true
}

func fitDiagonalRTLReversed(board *gameBoard, word string) bool {
	return fitDiagonalRTL(board, reverse(word))
}

func fitWord(board *gameBoard, word string) bool {
	wordDirection := rand.Intn(diagonalRTLReversed + 1)
	w := strings.Replace(word, " ", "", -1)
	switch wordDirection {
	case horizontal:
		return fitHorizontal(board, w)
	case horizontalReversed:
		return fitHorizontalReversed(board, w)
	case vertical:
		return fitVertical(board, w)
	case verticalReversed:
		return fitVerticalReversed(board, w)
	case diagonalLTR:
		return fitDiagonalLTR(board, w)
	case diagonalLTRReversed:
		return fitDiagonalLTRReversed(board, w)
	case diagonalRTL:
		return fitDiagonalRTL(board, w)
	case diagonalRTLReversed:
		return fitDiagonalRTLReversed(board, w)
	}
	return false
}

func buildBoard(board *gameBoard, words []string) bool {
	for _, word := range words {
		if !fitWord(board, word) {
			return false
		}
	}
	return true
}

func getRandomLetter() string {
	return string(letterBytes[rand.Intn(len(letterBytes))])
	// return "."
}

func fillBoard(board *gameBoard) {
	for i := range *board {
		for j := range (*board)[i] {
			if (*board)[i][j] == "" {
				(*board)[i][j] = getRandomLetter()
			}
		}
	}
}

func printBoard(board *gameBoard) {
	for i := range *board {
		for j := range (*board)[i] {
			c := (*board)[i][j]
			if c == "" {
				c = " "
			}
			fmt.Print(c)
		}
		fmt.Print("\n")
	}
}

func htmlBoard(board *gameBoard, words []string, output string) {
	f, err := os.Create(output)
	if err != nil {
		fmt.Println("Error:", err.Error())
		return
	}
	defer f.Close()
	fmt.Fprintln(f, `<div style="float:left">`)
	fmt.Fprintln(f, `<table style="font-family:'Lucida Console', monospace">`)
	for i := range *board {
		fmt.Fprintln(f, `<tr style="padding: 0; margin: 0">`)
		for j := range (*board)[i] {
			c := (*board)[i][j]
			fmt.Fprintf(f, `<td style="padding: 0; margin: 0; height: 20px; width: 20px">%s</td>`, c)
		}
		fmt.Fprintln(f, "</tr>")
	}
	fmt.Fprintln(f, "</table>")
	fmt.Fprintln(f, "</div>")
	fmt.Fprintln(f, `<div style="float:left">`)
	fmt.Fprintln(f, `<ul style="list-style: none;">`)
	for _, word := range words {
		fmt.Fprintf(f, "<li>%s</li>\n", word)
	}
	fmt.Fprintln(f, "</ul>")
	fmt.Fprintln(f, "</div>")
}

func main() {
	rand.Seed(time.Now().UTC().UnixNano())
	wordList, err := readWordList("wordlist.txt")
	if err != nil {
		fmt.Println("Error:", err.Error())
		return
	}
	fmt.Println(wordList)
	boardSize := 15
	maxBoardSize := 35
	for i := 0; i < boardAttempts; i++ {
		if (i%1000) == 0 && boardSize < maxBoardSize {
			boardSize += 5
		}
		board := makeBoard(boardSize, boardSize)
		built := buildBoard(&board, wordList)
		if built {
			fillBoard(&board)
			printBoard(&board)
			htmlBoard(&board, wordList, "output.html")
			fmt.Println("output.html saved with board size:", boardSize)
			return
		}
	}
	fmt.Println("Error: couldn't fit words")
}
