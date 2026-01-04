package game

import "errors"

func ApplyMove(g *Game, playerID string, col int) error {
	if g.Status == "finished" {
		return errors.New("game already finished")
	}

	if g.CurrentTurn != playerID {
		return errors.New("not your turn")
	}

	if col < 0 || col > 6 {
		return errors.New("invalid column")
	}

	row := -1
	for r := 5; r >= 0; r-- {
		if g.Board[r][col] == 0 {
			row = r
			break
		}
	}

	if row == -1 {
		return errors.New("column full")
	}

	player := g.Players[playerID]
	g.Board[row][col] = player.Color

	if checkWin(g.Board, row, col, player.Color) {
		g.Status = "finished"
		g.Winner = playerID
		return nil
	}

	if checkDraw(g.Board) {
		g.Status = "finished"
		g.Winner = "draw"
		return nil
	}

	for id := range g.Players {
		if id != playerID {
			g.CurrentTurn = id
			break
		}
	}

	return nil
}

func checkDraw(board [6][7]int) bool {
	for c := 0; c < 7; c++ {
		if board[0][c] == 0 {
			return false
		}
	}
	return true
}

func checkWin(board [6][7]int, r, c, color int) bool {
	dirs := [][2]int{{0, 1}, {1, 0}, {1, 1}, {1, -1}}

	for _, d := range dirs {
		count := 1
		for i := 1; i < 4; i++ {
			nr, nc := r+d[0]*i, c+d[1]*i
			if nr < 0 || nr >= 6 || nc < 0 || nc >= 7 || board[nr][nc] != color {
				break
			}
			count++
		}
		for i := 1; i < 4; i++ {
			nr, nc := r-d[0]*i, c-d[1]*i
			if nr < 0 || nr >= 6 || nc < 0 || nc >= 7 || board[nr][nc] != color {
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
