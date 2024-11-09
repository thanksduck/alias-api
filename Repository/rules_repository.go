package repository

import (
	"context"
	"fmt"

	"github.com/jackc/pgx/v5"
	db "github.com/thanksduck/alias-api/Database"
	models "github.com/thanksduck/alias-api/Models"
)

func FindRuleByID(id uint32) (*models.Rule, error) {
	pool := db.GetPool()
	var rule models.Rule
	err := pool.QueryRow(context.Background(),
		`SELECT id, user_id, username, alias_email, destination_email, active, comment,name FROM rules WHERE id = $1`, id).Scan(
		&rule.ID, &rule.UserID, &rule.Username, &rule.AliasEmail, &rule.DestinationEmail, &rule.Active, &rule.Comment, &rule.Name)
	if err != nil {
		if err == pgx.ErrNoRows {
			return nil, fmt.Errorf("rule with id %d not found", id)
		}
		return nil, fmt.Errorf("error querying rule: %w", err)
	}
	return &rule, nil
}

func FindRulesByUserID(userID uint32) ([]models.Rule, error) {
	pool := db.GetPool()
	rows, err := pool.Query(context.Background(),
		`SELECT id, user_id, username, alias_email, destination_email, active, comment,name FROM rules WHERE user_id = $1`, userID)
	if err != nil {
		return nil, fmt.Errorf("error querying rules: %w", err)
	}
	defer rows.Close()

	var rules []models.Rule
	for rows.Next() {
		var rule models.Rule
		err := rows.Scan(&rule.ID, &rule.UserID, &rule.Username, &rule.AliasEmail, &rule.DestinationEmail, &rule.Active, &rule.Comment, &rule.Name)
		if err != nil {
			return nil, fmt.Errorf("error scanning rule: %w", err)
		}
		rules = append(rules, rule)
	}
	return rules, nil
}

func FindRulesByUsername(username string) ([]models.Rule, error) {
	pool := db.GetPool()
	rows, err := pool.Query(context.Background(),
		`SELECT id, user_id, username, alias_email, destination_email, active, comment,name FROM rules WHERE username = $1`, username)
	if err != nil {
		return nil, fmt.Errorf("error querying rules: %w", err)
	}
	defer rows.Close()

	var rules []models.Rule
	for rows.Next() {
		var rule models.Rule
		err := rows.Scan(&rule.ID, &rule.UserID, &rule.Username, &rule.AliasEmail, &rule.DestinationEmail, &rule.Active, &rule.Comment, &rule.Name)
		if err != nil {
			return nil, fmt.Errorf("error scanning rule: %w", err)
		}
		rules = append(rules, rule)
	}
	return rules, nil
}

func FindRuleByAliasEmail(aliasEmail string) (*models.Rule, error) {
	pool := db.GetPool()
	var rule models.Rule
	err := pool.QueryRow(context.Background(),
		`SELECT id, user_id, username, alias_email, destination_email, active, comment FROM rules WHERE alias_email = $1`, aliasEmail).Scan(
		&rule.ID, &rule.UserID, &rule.Username, &rule.AliasEmail, &rule.DestinationEmail, &rule.Active, &rule.Comment)
	if err != nil {
		return nil, err
	}
	return &rule, nil
}

func CreateNewRule(rule *models.Rule) (*models.Rule, error) {
	pool := db.GetPool()
	tx, err := pool.Begin(context.Background())
	if err != nil {
		return nil, fmt.Errorf("failed to begin transaction: %w", err)
	}
	defer tx.Rollback(context.Background())

	err = tx.QueryRow(context.Background(),
		`INSERT INTO rules (user_id, username, alias_email, destination_email, comment, name)
		 VALUES ($1, $2, $3, $4, $5, $6) RETURNING id, user_id, username, alias_email, destination_email, active, comment, name`,
		rule.UserID, rule.Username, rule.AliasEmail, rule.DestinationEmail, rule.Comment, rule.Name).Scan(
		&rule.ID, &rule.UserID, &rule.Username, &rule.AliasEmail, &rule.DestinationEmail, &rule.Active, &rule.Comment, &rule.Name)
	if err != nil {
		fmt.Println(err)
		return nil, err
	}

	_, err = tx.Exec(context.Background(),
		`UPDATE users SET alias_count = alias_count + 1 WHERE id = $1`, rule.UserID)
	if err != nil {
		fmt.Println(err)
		return nil, err
	}

	err = tx.Commit(context.Background())
	if err != nil {
		return nil, fmt.Errorf("failed to commit transaction: %w", err)
	}

	return rule, nil
}

func UpdateRuleByID(id uint32, rule *models.Rule) (*models.Rule, error) {
	pool := db.GetPool()
	err := pool.QueryRow(context.Background(),
		`UPDATE rules SET alias_email = $1, destination_email = $2, comment = $3, name = $4, active = $5 WHERE id = $6 RETURNING id, user_id, username, alias_email, destination_email, active, comment, name`,
		rule.AliasEmail, rule.DestinationEmail, rule.Comment, rule.Name, rule.Active, id).Scan(
		&rule.ID, &rule.UserID, &rule.Username, &rule.AliasEmail, &rule.DestinationEmail, &rule.Active, &rule.Comment, &rule.Name)
	if err != nil {
		fmt.Println(err)
		return nil, err
	}
	return rule, nil
}

func ToggleRuleByID(id uint32) error {
	pool := db.GetPool()
	_, err := pool.Exec(context.Background(),
		`UPDATE rules SET active = NOT active WHERE id = $1`, id)
	if err != nil {
		fmt.Println(err)
		return err
	}
	return nil
}

func DeleteRuleByID(id uint32, userID uint32) error {
	pool := db.GetPool()
	tx, err := pool.Begin(context.Background())
	if err != nil {
		return fmt.Errorf("failed to begin transaction: %w", err)
	}
	defer tx.Rollback(context.Background())

	_, err = tx.Exec(context.Background(),
		`DELETE FROM rules WHERE id = $1`, id)
	if err != nil {
		fmt.Println(err)
		return err
	}

	_, err = tx.Exec(context.Background(),
		`UPDATE users SET alias_count = alias_count - 1 WHERE id = $1`, userID)
	if err != nil {
		fmt.Println(err)
		return err
	}

	err = tx.Commit(context.Background())
	if err != nil {
		return fmt.Errorf("failed to commit transaction: %w", err)
	}

	return nil
}
