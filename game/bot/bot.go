package bot

import (
	"errors"

	"fourinrow/game"
)

func GetBestMove(g *game.Game, botColor int) (int, error) {
	board := g.Board
	enemy := 1
	if botColor == 1 {
		enemy = 2
	}

	valid := validMoves(board)
	if len(valid) == 0 {
		return -1, errors.New("no valid moves")
	}

	for _, c := range valid {
		if canWin(board, c, botColor) {
			return c, nil
		}
	}

	for _, c := range valid {
		if canWin(board, c, enemy) {
			return c, nil
		}
	}

	order := []int{3, 2, 4, 1, 5, 0, 6}
	for _, c := range order {
		if board[0][c] == 0 {
			return c, nil
		}
	}

	return valid[0], nil
}

func validMoves(board [6][7]int) []int {
	var m []int
	for c := 0; c < 7; c++ {
		if board[0][c] == 0 {
			m = append(m, c)
		}
	}
	return m
}

func canWin(board [6][7]int, col, color int) bool {
	r := -1
	for i := 5; i >= 0; i-- {
		if board[i][col] == 0 {
			r = i
			break
		}
	}
	if r == -1 {
		return false
	}
	board[r][col] = color
	return check(board, r, col, color)
}

func check(b [6][7]int, r, c, p int) bool {
	d := [][2]int{{0, 1}, {1, 0}, {1, 1}, {1, -1}}
	for _, x := range d {
		n := 1
		for i := 1; i < 4; i++ {
			rr, cc := r+x[0]*i, c+x[1]*i
			if rr < 0 || rr >= 6 || cc < 0 || cc >= 7 || b[rr][cc] != p {
				break
			}
			n++
		}
		for i := 1; i < 4; i++ {
			rr, cc := r-x[0]*i, c-x[1]*i
			if rr < 0 || rr >= 6 || cc < 0 || cc >= 7 || b[rr][cc] != p {
				break
			}
			n++
		}
		if n >= 4 {
			return true
		}
	}
	return false
}
