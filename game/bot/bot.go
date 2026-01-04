package bot

import (
	"errors"

	"github.com/mukesh-1608/four-in-row/game"
)

// GetBestMove returns the column index for the best move according to deterministic rules.
//
// Priority:
// 1. Win immediately
// 2. Block opponent win
// 3. Center preference (3, 2/4, 1/5, 0/6)
func GetBestMove(g *game.Game, botColor int) (int, error) {
	board := g.Board

	opponentColor := 1
	if botColor == 1 {
		opponentColor = 2
	}

	validMoves := getValidMoves(board)
	if len(validMoves) == 0 {
		return -1, errors.New("no valid moves")
	}

	// 1. Winning move
	for _, col := range validMoves {
		if canWin(board, col, botColor) {
			return col, nil
		}
	}

	// 2. Block opponent
	for _, col := range validMoves {
		if canWin(board, col, opponentColor) {
			return col, nil
		}
	}

	// 3. Center priority
	centerPriority := []int{3, 2, 4, 1, 5, 0, 6}
	for _, col := range centerPriority {
		if isValidMove(board, col) {
			return col, nil
		}
	}

	return validMoves[0], nil
}

func getValidMoves(board [6][7]int) []int {
	moves := []int{}
	for c := 0; c < 7; c++ {
		if board[0][c] == 0 {
			moves = append(moves, c)
		}
	}
	return moves
}

func isValidMove(board [6][7]int, col int) bool {
	return col >= 0 && col < 7 && board[0][col] == 0
}

// canWin simulates a move and checks if it results in a win
func canWin(board [6][7]int, col int, color int) bool {
	row := -1
	for r := 5; r >= 0; r-- {
		if board[r][col] == 0 {
			row = r
			break
		}
	}

	if row == -1 {
		return false
	}

	board[row][col] = color
	return checkWinHelper(board, row, col, color)
}

// checkWinHelper duplicates win detection logic to keep bot deterministic
func checkWinHelper(board [6][7]int, lastRow, lastCol, playerNum int) bool {
	directions := [][2]int{
		{0, 1},  // Horizontal
		{1, 0},  // Vertical
		{1, 1},  // Diagonal \
		{1, -1}, // Diagonal /
	}

	for _, d := range directions {
		dr, dc := d[0], d[1]
		count := 1

		for i := 1; i < 4; i++ {
			r, c := lastRow+dr*i, lastCol+dc*i
			if r < 0 || r >= 6 || c < 0 || c >= 7 || board[r][c] != playerNum {
				break
			}
			count++
		}

		for i := 1; i < 4; i++ {
			r, c := lastRow-dr*i, lastCol-dc*i
			if r < 0 || r >= 6 || c < 0 || c >= 7 || board[r][c] != playerNum {
				break
			}
			count++
		}

		if count >= 4 {
			return true
		}
	}

	return false
}
