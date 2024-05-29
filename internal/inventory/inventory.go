package inventory

import "database/sql"

type Service struct {
	db *sql.DB
}

type Item struct {
	ID       int
	Quantity int
}

func NewService(db *sql.DB) *Service {
	return &Service{db: db}
}

func (s *Service) AddItem(userID, itemID, quantity int) error {
	query := `
		INSERT INTO inventory (user_id, item_id, quantity) 
		VALUES ($1, $2, $3) 
		ON CONFLICT (user_id, item_id) 
		DO UPDATE SET quantity = inventory.quantity + $3`
	_, err := s.db.Exec(query, userID, itemID, quantity)
	return err
}

func (s *Service) RemoveItem(userID int, itemID int, quantity int) error {
	tx, err := s.db.Begin()
	if err != nil {
		return err
	}

	query := `
		UPDATE inventory 
		SET quantity = GREATEST(0, quantity - $3) 
		WHERE user_id = $1 AND item_id = $2 AND quantity >= $3`
	result, err := tx.Exec(query, userID, itemID, quantity)
	if err != nil {
		tx.Rollback()
		return err
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		tx.Rollback()
		return err
	}

	if rowsAffected > 0 {
		query := `
			DELETE FROM inventory 
       		WHERE user_id = $1 AND item_id = $2 AND quantity = 0`
		_, err = tx.Exec(query, userID, itemID)
		if err != nil {
			tx.Rollback()
			return err
		}
	}

	return tx.Commit()
}

func (s *Service) GetUserItems(userID int) ([]Item, error) {
	query := `
		SELECT item_id, quantity 
		FROM inventory 
		WHERE user_id = $1`
	rows, err := s.db.Query(query, userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	items := []Item{}
	for rows.Next() {
		var item Item
		if err := rows.Scan(&item.ID, &item.Quantity); err != nil {
			return nil, err
		}
	}
	return items, nil
}
