package store

func (s *Store) LogAction(userID int, action, target, ip string) error {
	_, err := s.db.Exec(
		`INSERT INTO actions_log (user_id, action, target, ip)
		 VALUES ($1, $2, $3, $4)`,
		userID, action, target, ip,
	)
	return err
}
