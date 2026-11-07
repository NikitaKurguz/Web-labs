package store

import (
	"strconv"
	"time"
)

type ActionRow struct {
	ID        int       `json:"id"`
	UserID    int       `json:"userId"`
	Action    string    `json:"action"`
	Target    string    `json:"target"`
	IP        string    `json:"ip"`
	CreatedAt time.Time `json:"createdAt"`
}

func (s *Store) GetActions(filterAction, sortBy, order string, limit int) ([]ActionRow, error) {
	allowedSort := map[string]bool{
		"created_at": true,
		"action":     true,
		"user_id":    true,
	}
	if !allowedSort[sortBy] {
		sortBy = "created_at"
	}
	if order != "asc" && order != "desc" {
		order = "desc"
	}
	if limit <= 0 || limit > 1000 {
		limit = 200
	}

	query := `
		SELECT id, user_id, action, target, ip, created_at
		FROM actions_log
	`
	args := []any{}
	argN := 1

	if filterAction != "" {
		query += ` WHERE action = $1`
		args = append(args, filterAction)
		argN = 2
	}

	query += ` ORDER BY ` + sortBy + ` ` + order
	query += ` LIMIT $` + strconv.Itoa(argN)
	args = append(args, limit)

	rows, err := s.db.Query(query, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var out []ActionRow
	for rows.Next() {
		var r ActionRow
		if err := rows.Scan(&r.ID, &r.UserID, &r.Action, &r.Target, &r.IP, &r.CreatedAt); err != nil {
			return nil, err
		}
		out = append(out, r)
	}
	return out, rows.Err()
}

type StatsRow struct {
	Period time.Time `json:"period"`
	Count  int       `json:"count"`
}

func (s *Store) GetStats(groupBy, filterAction, sortBy, order string, limit int) ([]StatsRow, error) {
	allowedGroupBy := map[string]bool{
		"hour": true, "day": true, "week": true, "month": true,
	}
	if !allowedGroupBy[groupBy] {
		groupBy = "day"
	}

	allowedSortBy := map[string]bool{"period": true, "count": true}
	if !allowedSortBy[sortBy] {
		sortBy = "period"
	}

	if order != "asc" && order != "desc" {
		order = "desc"
	}

	query := `
		SELECT
			date_trunc($1, created_at) AS period,
			COUNT(*) AS count
		FROM actions_log
	`
	args := []any{groupBy}

	if filterAction != "" {
		query += ` WHERE action = $2`
		args = append(args, filterAction)
	}

	query += ` GROUP BY period ORDER BY ` + sortBy + ` ` + order

	if limit > 0 {
		query += ` LIMIT $` + strconv.Itoa(len(args)+1)
		args = append(args, limit)
	}

	rows, err := s.db.Query(query, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var out []StatsRow
	for rows.Next() {
		var r StatsRow
		if err := rows.Scan(&r.Period, &r.Count); err != nil {
			return nil, err
		}
		out = append(out, r)
	}
	return out, rows.Err()
}
