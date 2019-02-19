import "time"

type Game struct {
    Id          int    `json: "id"`
    Title       string `json: "title"`
    StartedAt   time   `json: "started_at"`
    GameYear    time   `json: "game_year"`
    Phase       int    `json: "phase"`
    PhaseEnd    time   `json: "phase_end"`
    OrdersInterval int `json: "orders_interval"`
}
