package bot

import (
	"connect4/game"
	"errors"
)

// GetBestMove returns the column index for the best move according to the deterministic rules.
// Rules priority:
// 1. Win immediately
// 2. Block opponent win
// 3. Center preference (3, 2/4, 1/5, 0/6)
func GetBestMove(g *game.Game, botColor int) (int, error) {
	// Create a copy of the board to simulate moves
	board := g.Board
	opponentColor := 1
	if botColor == 1 {
		opponentColor = 2
	}

	// Helper to find valid moves
	validMoves := getValidMoves(board)
	if len(validMoves) == 0 {
		return -1, errors.New("no valid moves")
	}

	// 1. Check for winning move
	for _, col := range validMoves {
		if canWin(board, col, botColor) {
			return col, nil
		}
	}

	// 2. Check for blocking opponent win
	for _, col := range validMoves {
		if canWin(board, col, opponentColor) {
			return col, nil
		}
	}

	// 3. Prefer center columns
	// Priority order: 3, 2, 4, 1, 5, 0, 6
	centerPriority := []int{3, 2, 4, 1, 5, 0, 6}
	for _, col := range centerPriority {
		if isValidMove(board, col) {
			return col, nil
		}
	}

	// Fallback (should be covered by center priority if there are valid moves)
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

// canWin simulates placing a piece and checks if it results in a win
func canWin(board [6][7]int, col int, color int) bool {
	// Simulate gravity
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

	// Place piece temporarily
	board[row][col] = color

	// Check win (reusing the logic from game package would be ideal, 
	// but to avoid circular imports or refactoring for now, we duplicate the simple checkWin logic locally 
	// or we export the checkWin from game package publicly.
	// Since game/logic.go has checkWin unexported, let's copy the pure logic helper here for independence/determinism.)
	return checkWinHelper(board, row, col, color)
}

// checkWinHelper is a copy of the win detection logic
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

		// Forward
		for i := 1; i < 4; i++ {
			r, c := lastRow+dr*i, lastCol+dc*i
			if r < 0 || r >= 6 || c < 0 || c >= 7 || board[r][c] != playerNum {
				break
			}
			count++
		}

		// Backward
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
