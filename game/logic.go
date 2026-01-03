package game

import (
        "errors"
)

const (
        Rows = 6
        Cols = 7
)

// ApplyMove attempts to place a disc in the specified column for the given player.
// It validates the move, updates the board, checks for win/draw conditions,
// and updates the game status and turn.
func ApplyMove(g *Game, playerID string, column int) error {
        // 1. Validate Game Status
        if g.Status != "playing" {
                return errors.New("game is not active")
        }

        // 2. Validate Turn
        if g.CurrentTurn != playerID {
                return errors.New("not your turn")
        }

        // 3. Validate Column Range
        if column < 0 || column >= Cols {
                return errors.New("invalid column index")
        }

        // 4. Determine Player Number (1 or 2)
        // We need a deterministic way to map playerID to 1 or 2.
        // In a real match, we'd store this mapping.
        // For this exercise, we'll iterate the map (order undefined) or check against specific assignments.
        // To be safe, let's assume Players map keys are player IDs.
        // We need to know which int corresponds to the player.
        // Let's assume we add a PlayerNumber field to Player or derived it.
        // For now, let's iterate to find the opponent to switch turn.
        var playerNum int
        var opponentID string

        // This is a bit inefficient but works for the skeleton.
        // Ideally, Game struct has P1 and P2 explicitly or Player struct has the number.
        // Let's modify models.go to include PlayerNum or infer it here.
        // In the absence of that, let's just use 1 for CurrentTurn and 2 for the other.
        // Actually, we need to know if 'playerID' corresponds to 1 or 2 on the board.
        // A simple way is to assign them when the game starts.
        // For this step (logic only), let's assume we can determine it.
        // Let's add a helper `GetPlayerNumber` or update the model.
        // Updating the model is cleaner. I will add `Color int` to Player (1 or 2).

        playerNum = g.Players[playerID].Color
        if playerNum == 0 {
                return errors.New("player color not assigned")
        }

        // Find opponent ID for turn switching
        for id := range g.Players {
                if id != playerID {
                        opponentID = id
                        break
                }
        }

        // 5. Find lowest empty row (Gravity)
        row := -1
        for r := Rows - 1; r >= 0; r-- {
                if g.Board[r][column] == 0 {
                        row = r
                        break
                }
        }

        if row == -1 {
                return errors.New("column is full")
        }

        // 6. Apply Move
        g.Board[row][column] = playerNum

        // 7. Check Win
        if checkWin(g.Board, row, column, playerNum) {
                g.Status = "finished"
                g.Winner = playerID
                return nil
        }

        // 8. Check Draw
        if checkDraw(g.Board) {
                g.Status = "finished"
                g.Winner = "draw"
                return nil
        }

        // 9. Switch Turn
        if opponentID != "" {
                g.CurrentTurn = opponentID
        } else {
                // Single player testing or weird state
                // keep same turn? or just leave it.
        }

        return nil
}

// checkWin checks for 4 connected discs including the position (lastRow, lastCol)
func checkWin(board [6][7]int, lastRow, lastCol, playerNum int) bool {
        // Directions: Horizontal, Vertical, Diagonal /, Diagonal \
        directions := [][2]int{
                {0, 1},  // Horizontal
                {1, 0},  // Vertical
                {1, 1},  // Diagonal \ (down-right)
                {1, -1}, // Diagonal / (down-left)
        }

        for _, d := range directions {
                dr, dc := d[0], d[1]
                count := 1 // Count the piece just placed

                // Check forward
                for i := 1; i < 4; i++ {
                        r, c := lastRow+dr*i, lastCol+dc*i
                        if r < 0 || r >= Rows || c < 0 || c >= Cols || board[r][c] != playerNum {
                                break
                        }
                        count++
                }

                // Check backward
                for i := 1; i < 4; i++ {
                        r, c := lastRow-dr*i, lastCol-dc*i
                        if r < 0 || r >= Rows || c < 0 || c >= Cols || board[r][c] != playerNum {
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

// checkDraw checks if the board is full
func checkDraw(board [6][7]int) bool {
        for c := 0; c < Cols; c++ {
                // If the top row has any empty spot, it's not a draw yet
                // (Gravity ensures lower rows are filled if top is filled, usually,
                // but checking top row for 0 is sufficient for full columns)
                if board[0][c] == 0 {
                        return false
                }
        }
        return true
}
