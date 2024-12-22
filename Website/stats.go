package main

func loadStats() StatsResponse {
	// Here, you'd typically query a DB or metrics aggregator to get real data.
	response := StatsResponse{
		AvgLatencyMs:  123,
		ProjectsCount: 42,
		TimeSaved:     217,
	}
	return response
}
