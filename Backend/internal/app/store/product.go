package store

import "primarte/internal/app/model"

func (s *Store) GetAllProducts() ([]*model.Product, error) {
	rows, err := s.db.Query(
		`SELECT id, name, img, type, color, category, price 
		FROM products 
		ORDER BY id`,
	)
	if err != nil {
		return nil, err
	}

	defer rows.Close()

	var products []*model.Product

	for rows.Next() {
		pr := &model.Product{}
		err := rows.Scan(&pr.ID, &pr.Name, &pr.Img,
			&pr.Type, &pr.Color, &pr.Category, &pr.Price)

		if err != nil {
			return nil, err
		}

		products = append(products, pr)
	}
	return products, err
}

func (s *Store) GetProductById(id int) (*model.Product, error) {
	pr := &model.Product{}
	err := s.db.QueryRow(
		`SELECT id, name, img, type, color, category, price 
		FROM products 
		WHERE id = $1`,
		id,
	).Scan(&pr.ID, &pr.Name, &pr.Img, &pr.Type, &pr.Color,
		&pr.Category, &pr.Price)

	if err != nil {
		return nil, err
	}

	return pr, nil
}

func (s *Store) CreateProduct(pr *model.Product) error {
	err := s.db.QueryRow(
		`INSERT INTO products (name, img, type, color, category, price)
		VALUES ($1, $2, $3, $4, $5, $6)
		RETURNING id`,
		pr.Name, pr.Img, pr.Type, pr.Color, pr.Category, pr.Price,
	).Scan(&pr.ID)

	return err
}

func (s *Store) UpdateProduct(pr *model.Product) error {
	_, err := s.db.Exec(
		`UPDATE products 
		SET name = $1, img = $2, type = $3, 
		color = $4, category = $5, price = $6 
		WHERE id = $7`,
		pr.Name, pr.Img, pr.Type, pr.Color, pr.Category, pr.Price, pr.ID,
	)
	return err
}

func (s *Store) DeleteProductById(id int) error {
	_, err := s.db.Exec(`DELETE FROM products WHERE id = $1`, id)
	return err
}
