package main

import (
	"database/sql"
	"encoding/json"
	"log"
	"net/http"
)

type eventRequest struct {
	Path     string `json:"path"`
	Referrer string `json:"referrer"`
}

func handleEvent(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req eventRequest
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			http.Error(w, "bad request", http.StatusBadRequest)
			return
		}

		_, err := db.ExecContext(r.Context(),
			`INSERT INTO events (path, referrer, user_agent, ip) VALUES ($1, $2, $3, $4)`,
			req.Path, req.Referrer, r.Header.Get("User-Agent"), r.RemoteAddr,
		)
		if err != nil {
			log.Printf("insert event: %v", err)
			http.Error(w, "internal error", http.StatusInternalServerError)
			return
		}

		w.WriteHeader(http.StatusNoContent)
	}
}

type statsResponse struct {
	Total  int            `json:"total"`
	ByPath map[string]int `json:"by_path"`
	ByDay  map[string]int `json:"by_day"`
}

func handleStats(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var total int
		db.QueryRowContext(r.Context(), `SELECT COUNT(*) FROM events`).Scan(&total)

		byPath := make(map[string]int)
		rows, err := db.QueryContext(r.Context(),
			`SELECT path, COUNT(*) FROM events GROUP BY path ORDER BY COUNT(*) DESC`)
		if err == nil {
			defer rows.Close()
			for rows.Next() {
				var path string
				var count int
				rows.Scan(&path, &count)
				byPath[path] = count
			}
		}

		byDay := make(map[string]int)
		rows2, err := db.QueryContext(r.Context(),
			`SELECT DATE(created_at)::text, COUNT(*) FROM events GROUP BY DATE(created_at) ORDER BY 1 DESC`)
		if err == nil {
			defer rows2.Close()
			for rows2.Next() {
				var day string
				var count int
				rows2.Scan(&day, &count)
				byDay[day] = count
			}
		}

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(statsResponse{total, byPath, byDay})
	}
}

func corsMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", "https://x0lie.com")
		w.Header().Set("Access-Control-Allow-Methods", "POST, GET, OPTIONS")
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type")

		if r.Method == http.MethodOptions {
			w.WriteHeader(http.StatusNoContent)
			return
		}

		next.ServeHTTP(w, r)
	})
}
