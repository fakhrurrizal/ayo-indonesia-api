package reqres

type MatchReportResponse struct {
    Match              MatchResponse   `json:"match"`
    MatchResult        string          `json:"match_result"`
    TopScorer          *PlayerResponse `json:"top_scorer"`
    HomeTeamWins       int64           `json:"home_team_total_wins"`
    AwayTeamWins       int64           `json:"away_team_total_wins"`
}