package game

import (
	"bufio"
	"fmt"
	"math/rand"
	"os"
	"strings"
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
	letterBytes     = "ABCDEFGHIJKLMNOPQRSTUVWXYZ"
	fitWordAttempts = 10000
)

// Board holds the state of the word search
type Board [][]string

type placement struct {
	x, y   int
	letter string
	style  int
}

// ReadWordList ingests a list of words from a word file
func ReadWordList(path string) (lines []string, err error) {
	file, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		lines = append(lines, strings.ToUpper(scanner.Text()))
	}
	return lines, scanner.Err()
}

func makeBoard(x int, y int) *Board {
	board := make(Board, y)
	for i := range board {
		board[i] = make([]string, x)
	}
	return &board
}

func copyBoard(board *Board) *Board {
	cpy := make(Board, len((*board)))
	for i := range *board {
		cpy[i] = make([]string, len((*board)[i]))
		for j := range (*board)[i] {
			cpy[i][j] = (*board)[i][j]
		}
	}
	return &cpy
}

func reverse(s string) string {
	chars := []rune(s)
	for i, j := 0, len(chars)-1; i < j; i, j = i+1, j-1 {
		chars[i], chars[j] = chars[j], chars[i]
	}
	return string(chars)
}

func fitHorizontal(word string, maxX, maxY int) []placement {
	placements := make([]placement, len(word))
	split := strings.Split(word, "")
	found := false
	count := 0
	for !found && count < fitWordAttempts {
		count++
		ix := rand.Intn(maxX)
		iy := rand.Intn(maxY)
		if iy+len(word) >= maxY {
			continue
		}
		for j, c := range split {
			placements[j] = placement{
				ix,
				iy + j,
				c,
				horizontal,
			}
		}
		found = true
	}
	if count >= fitWordAttempts {
		placements = nil
	}
	return placements
}

func fitVertical(word string, maxX, maxY int) []placement {
	placements := make([]placement, len(word))
	split := strings.Split(word, "")
	found := false
	count := 0
	for !found && count < fitWordAttempts {
		count++
		ix := rand.Intn(maxX)
		iy := rand.Intn(maxY)
		if ix+len(word) >= maxX {
			continue
		}
		for j, c := range split {
			placements[j] = placement{
				ix + j,
				iy,
				c,
				vertical,
			}
		}
		found = true
	}
	if count >= fitWordAttempts {
		placements = nil
	}
	return placements
}

func fitDiagonalLTR(word string, maxX, maxY int) []placement {
	placements := make([]placement, len(word))
	split := strings.Split(word, "")
	found := false
	count := 0
	for !found && count < fitWordAttempts {
		count++
		ix := rand.Intn(maxX)
		iy := rand.Intn(maxY)
		if ix+len(word) >= maxX || iy+len(word) >= maxY {
			continue
		}
		for j, c := range split {
			placements[j] = placement{
				ix + j,
				iy + j,
				c,
				diagonalLTR,
			}
		}
		found = true
	}
	if count >= fitWordAttempts {
		placements = nil
	}
	return placements
}

func fitDiagonalRTL(word string, maxX, maxY int) []placement {
	placements := make([]placement, len(word))
	split := strings.Split(word, "")
	found := false
	count := 0
	for !found && count < fitWordAttempts {
		count++
		ix := rand.Intn(maxX)
		iy := rand.Intn(maxY)
		if ix+len(word) >= maxX || iy-len(word) < 0 {
			continue
		}
		for j, c := range split {
			placements[j] = placement{
				ix + j,
				iy - j,
				c,
				diagonalRTL,
			}
		}
		found = true
	}
	if count >= fitWordAttempts {
		placements = nil
	}
	return placements
}

func fitWord(board *Board, word string) bool {
	wordDirection := rand.Intn(diagonalRTLReversed + 1)
	w := strings.Replace(word, " ", "", -1)
	var placements []placement
	x := len((*board))
	y := len((*board)[0])
	switch wordDirection {
	case horizontal:
		placements = fitHorizontal(w, x, y)
	case horizontalReversed:
		placements = fitHorizontal(reverse(w), x, y)
	case vertical:
		placements = fitVertical(w, x, y)
	case verticalReversed:
		placements = fitVertical(reverse(w), x, y)
	case diagonalLTR:
		placements = fitDiagonalLTR(w, x, y)
	case diagonalLTRReversed:
		placements = fitDiagonalLTR(reverse(w), x, y)
	case diagonalRTL:
		placements = fitDiagonalRTL(w, x, y)
	case diagonalRTLReversed:
		placements = fitDiagonalRTL(reverse(w), x, y)
	}
	if placements == nil {
		return false
	}
	for _, p := range placements {
		if (*board)[p.x][p.y] == "" || (*board)[p.x][p.y] == p.letter {
			(*board)[p.x][p.y] = p.letter
		} else {
			return false
		}
	}
	return true
}

func getRandomLetter() string {
	return string(letterBytes[rand.Intn(len(letterBytes))])
	// return "."
}

func fillBoard(board *Board) {
	for i := range *board {
		for j := range (*board)[i] {
			if (*board)[i][j] == "" {
				(*board)[i][j] = getRandomLetter()
			}
		}
	}
}

// BuildBoard takes a size and a slice of words and returns a Board
func BuildBoard(size int, words []string) (board *Board) {
	for _, word := range words {
		count := 0
		fit := false
		var stateBeforeWord *Board

		if board == nil {
			board = makeBoard(size, size)
		}

		for !fit && count < fitWordAttempts {
			count++
			stateBeforeWord = copyBoard(board)
			fit = fitWord(stateBeforeWord, word)
		}
		if fit {
			board = stateBeforeWord
		} else {
			return nil
		}
	}
	fillBoard(board)
	return board
}

// PrintBoard prints an ascii version of the board - good for debugging
func PrintBoard(board *Board) {
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

//HTMLBoard generates an HTML file ready to print
func HTMLBoard(board *Board, words []string, output string) {
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
