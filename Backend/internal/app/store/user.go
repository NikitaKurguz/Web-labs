package store

import "primarte/internal/app/model"

func (s *Store) CreateUser(u *model.User) error {
	return s.db.QueryRow(`INSERT INTO users (user_login, user_password, user_fullname, user_email, user_avatar) VALUES ($1, $2, $3, $4, $5) RETURNING user_id`,
		u.Login, u.Password, u.FullName, u.Email, u.Avatar).Scan(&u.Id)
}

func (s *Store) GetUserByLogin(login string) (*model.User, error) {
	u := &model.User{}
	err := s.db.QueryRow(
		`SELECT user_id, user_login, user_password, user_fullname, user_email, user_avatar, user_hash, is_admin
		FROM users
		WHERE user_login = $1`,
		login,
	).Scan(&u.Id, &u.Login, &u.Password, &u.FullName,
		&u.Email, &u.Avatar, &u.Hash, &u.IsAdmin)

	if err != nil {
		return nil, err
	}

	return u, nil
}

func (s *Store) GetUserById(Id int) (*model.User, error) {
	u := &model.User{}
	err := s.db.QueryRow(
		`SELECT user_id, user_login, user_password, user_fullname, user_email, user_avatar, user_hash, is_admin
		FROM users
		WHERE user_id = $1`,
		Id,
	).Scan(&u.Id, &u.Login, &u.Password, &u.FullName,
		&u.Email, &u.Avatar, &u.Hash, &u.IsAdmin)

	if err != nil {
		return nil, err
	}

	return u, nil
}

func (s *Store) UpdateHash(id int, hash string) error {
	_, err := s.db.Exec(
		`UPDATE users SET user_hash = $1
		WHERE user_id = $2`,
		hash,
		id,
	)

	return err
}

func (s *Store) UpdateUser(u *model.User) error {
	_, err := s.db.Exec(
		`UPDATE users
		 SET user_fullname = $1, user_email = $2, user_avatar = $3
		 WHERE user_id = $4`,
		u.FullName, u.Email, u.Avatar, u.Id,
	)
	return err
}

func (s *Store) UpdatePassword(id int, passwordHash string) error {
	_, err := s.db.Exec(
		`UPDATE users SET user_password = $1 WHERE user_id = $2`,
		passwordHash, id,
	)
	return err
}
