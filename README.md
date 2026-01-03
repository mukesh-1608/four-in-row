# Connect Four - Backend Engineering Assessment

## 1. Project Overview

This is a real-time multiplayer implementation of the classic "Connect Four" game, built as a backend engineering assessment. The project emphasizes correctness, clean architecture, and pure logic implementation without reliance on heavy frameworks.

**Key Technologies:**
- **Language:** Pure Go (Golang) 1.20+
- **Network:** `net/http` standard library + `gorilla/websocket` for real-time events
- **Database:** PostgreSQL (via `database/sql` + `lib/pq`) for historical data
- **Analytics:** Stubbed Kafka abstraction for event streaming

## 2. Architecture Overview

### HTTP Server
- A lightweight HTTP server listens on port 5000.
- Serves static files for a minimal test frontend.
- Exposes REST endpoints (`GET /leaderboard`) and a WebSocket upgrade endpoint (`/ws`).

### WebSocket Flow
1.  **Connection:** Clients connect to `/ws`.
2.  **State Management:** The server maintains an in-memory store (`game.Store`) of *active* games to ensure low-latency access.
3.  **Event Loop:** A dedicated goroutine per connection handles incoming messages (Move) and outgoing updates (Board State).
4.  **Broadcast:** State changes are broadcast to all connected players in the specific game.

### Game Lifecycle
- **Waiting:** Game created, waiting for players (simulated in this demo).
- **Playing:** Players take turns dropping discs. Bot responds immediately to human moves.
- **Finished:** Game ends on Win, Draw, or Forfeit. State is persisted to DB, and analytics events are emitted.

## 3. Game Logic

### Board Representation
- The board is a `[6][7]int` grid.
- `0`: Empty
- `1`: Player 1 (Red)
- `2`: Player 2 (Yellow/Bot)

### Logic Rules
- **Gravity:** Discs always fall to the lowest available row in the selected column.
- **Turn Enforcement:** The server rigorously validates that moves are only accepted from the current player.
- **Win Detection:** Checks 4 directions (Horizontal, Vertical, Diagonal /, Diagonal \) after every move.
- **Draw Detection:** Checks if the top row is full (implying the board is full).

## 4. Bot Strategy

The Bot is deterministic and follows a strict priority hierarchy:
1.  **Win:** If a winning move exists, take it immediately.
2.  **Block:** If the opponent has a winning move next turn, block it.
3.  **Center:** Prefer the center column (3), then radiating outwards (2,4 -> 1,5 -> 0,6). This maximizes strategic connection possibilities.

The bot simulates moves on a copy of the board to evaluate outcomes without affecting the live game state.

## 5. Reconnect Handling

- **Disconnect Detection:** WebSocket connection closure triggers a handler.
- **Grace Period:** A 30-second timer starts upon disconnect.
- **Reconnect:** If the player reconnects with the same ID within 30s, the timer is cancelled, and the connection is restored seamlessly.
- **Forfeit:** If the timer expires, the game is forfeited, and the opponent (or CPU) is declared the winner.

## 6. Persistence & Analytics

### Persistence
- **Strategy:** Only *completed* games are persisted to PostgreSQL to minimize write load and complexity.
- **Schema:** `games` table stores IDs, players, winner, and timestamps.
- **Reliability:** Database failures are logged but do not crash the game server (soft failure).

### Analytics
- **Abstraction:** A `AnalyticsProducer` interface abstracts the event stream.
- **Implementation:** Currently uses a `StubProducer` that logs events (`game_started`, `move_played`, `game_completed`) to stdout.
- **Future-Proofing:** The interface allows swapping in a real Kafka/Redpanda producer without changing game logic.

## 7. Trade-offs & Future Improvements

- **Matchmaking:** Currently stubbed with a single test game instance. A real lobby system/queue would be the next major feature.
- **Concurrency:** Game state is protected by mutexes. For massive scale, sharding game rooms across server instances (with Redis pub/sub) would be required.
- **Security:** No authentication is implemented. Production would require JWT/Session validation.
- **Testing:** Unit tests coverage could be expanded, particularly for edge cases in win detection.

## 8. How to Run Locally

The environment is pre-configured.

1.  **Start the Server:**
    ```bash
    go run main.go
    ```
2.  **Access the Game:**
    Open the Webview or navigate to `http://localhost:5000`.
3.  **Test Bot:** Click "Connect as Human" to play against the Bot.
4.  **Test Reconnect:** Click "Disconnect", wait 5s, and "Connect" again.
5.  **View Leaderboard:** `GET /leaderboard` (requires DB).
