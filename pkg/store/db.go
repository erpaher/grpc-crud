package store

import (
	"context"
	"database/sql"
	"log"

	"github.com/jackc/pgx/v4"
	"github.com/jackc/pgx/v4/stdlib"
	migrate "github.com/rubenv/sql-migrate"
)

type db struct {
	conn *sql.DB
}

func Database(url string) *db {
	connConfig, _ := pgx.ParseConfig(url)
	connStr := stdlib.RegisterConnConfig(connConfig)
	conn, err := sql.Open("pgx", connStr)
	if err != nil {
		log.Fatalf("Unable to connect to database: %v\n", err)
	}

	migrations := &migrate.FileMigrationSource{
		Dir: "/home/paher/go/src/github.com/erpaher/grpc_crud/pkg/store/migrations",
	}

	_, err = migrate.Exec(conn, "postgres", migrations, migrate.Up)
	if err != nil {
		log.Fatalf("Unable to migrate database: %v\n", err)
	}

	return &db{conn}
}

func (db *db) CreateUser(ctx context.Context, user User) (*User, error) {
	conn, err := stdlib.AcquireConn(db.conn)
	if err != nil {
		return nil, err
	}
	defer stdlib.ReleaseConn(db.conn, conn)
	err = conn.BeginFunc(ctx, func(tx pgx.Tx) error {
		err := tx.QueryRow(ctx, "INSERT INTO users (name, age, type) VALUES ($1, $2, $3) RETURNING id, created_at, updated_at", user.Name, user.Age, user.UserType).Scan(&user.ID, &user.CreatedAt, &user.UpdatedAt)
		if err != nil {
			return err
		}
		for i, item := range user.Items {
			err := tx.QueryRow(ctx, "INSERT INTO items (name, user_id) VALUES ($1, $2) RETURNING id, user_id, created_at, updated_at", item.Name, user.ID).Scan(&user.Items[i].ID, &user.Items[i].UserID, &user.Items[i].CreatedAt, &user.Items[i].UpdatedAt)
			if err != nil {
				return err
			}
		}
		return nil
	})
	if err != nil {
		return nil, err
	}
	return &user, nil
}

func (db *db) UpdateUser(ctx context.Context, user User) (*User, error) {
	conn, err := stdlib.AcquireConn(db.conn)
	if err != nil {
		return nil, err
	}
	defer stdlib.ReleaseConn(db.conn, conn)
	err = conn.BeginFunc(ctx, func(tx pgx.Tx) error {
		err := tx.QueryRow(ctx, "UPDATE users SET name=$1, age=$2, type=$3, updated_at=NOW() WHERE id=$4 RETURNING created_at, updated_at", user.Name, user.Age, user.UserType, user.ID).Scan(&user.CreatedAt, &user.UpdatedAt)
		if err != nil {
			return err
		}
		for i, item := range user.Items {
			err := tx.QueryRow(ctx, "UPDATE items SET name=$1, updated_at=NOW() WHERE id=$2 AND user_id=$3 RETURNING created_at, updated_at", item.Name, item.ID, user.ID).Scan(&user.Items[i].CreatedAt, &user.Items[i].UpdatedAt)
			if err != nil {
				return err
			}
		}
		return nil
	})
	if err != nil {
		return nil, err
	}
	return &user, nil
}

func (db *db) DeleteUser(ctx context.Context, userID uint32) error {
	conn, err := stdlib.AcquireConn(db.conn)
	if err != nil {
		return err
	}
	defer stdlib.ReleaseConn(db.conn, conn)
	err = conn.BeginFunc(ctx, func(tx pgx.Tx) error {
		_, err := tx.Exec(ctx, "DELETE FROM users CASCADE WHERE id=$1", userID)
		return err
	})
	return err
}

func (db *db) GetUser(ctx context.Context, userID uint32) (*User, error) {
	conn, err := stdlib.AcquireConn(db.conn)
	if err != nil {
		return nil, err
	}
	defer stdlib.ReleaseConn(db.conn, conn)
	user := User{}
	err = conn.BeginFunc(ctx, func(tx pgx.Tx) error {
		rows, err := tx.Query(ctx, "SELECT users.id, users.name, users.age, users.type, users.created_at, users.updated_at,  items.id, items.name, items.created_at, items.updated_at FROM users LEFT JOIN items ON (items.user_id=users.id) WHERE users.id=$1", userID)
		if err != nil {
			return err
		}
		defer rows.Close()

		err = sql.ErrNoRows
		for rows.Next() {
			var (
				itemID                       sql.NullInt64
				itemName                     sql.NullString
				itemUpdatedAt, itemCreatedAt sql.NullTime
			)
			err = rows.Scan(&user.ID, &user.Name, &user.Age, &user.UserType, &user.CreatedAt, &user.UpdatedAt, &itemID, &itemName, &itemCreatedAt, &itemUpdatedAt)
			if err != nil {
				return err
			}
			if itemID.Valid {
				user.Items = append(user.Items, Item{uint32(itemID.Int64), itemName.String, user.ID, &itemCreatedAt.Time, &itemUpdatedAt.Time})
			}
		}
		return err
	})
	if err != nil {
		return nil, err
	}
	return &user, nil
}

func (db *db) ListUser(ctx context.Context, page, limit uint32) (users []*User, err error) {
	conn, err := stdlib.AcquireConn(db.conn)
	if err != nil {
		return nil, err
	}
	defer stdlib.ReleaseConn(db.conn, conn)
	err = conn.BeginFunc(ctx, func(tx pgx.Tx) error {
		rows, err := tx.Query(ctx, "SELECT id FROM users OFFSET $1 LIMIT $2", page*limit, limit)
		if err != nil {
			return err
		}
		defer rows.Close()
		err = sql.ErrNoRows
		for rows.Next() {
			var userID uint32
			err = rows.Scan(&userID)
			if err != nil {
				return err
			}
			user, err := db.GetUser(ctx, userID)
			if err != nil {
				return err
			}
			users = append(users, user)
		}
		return err
	})
	if err != nil {
		return nil, err
	}
	return users, nil
}

func (db *db) CreateItem(ctx context.Context, name string, userID uint32) (*Item, error) {
	conn, err := stdlib.AcquireConn(db.conn)
	if err != nil {
		return nil, err
	}
	defer stdlib.ReleaseConn(db.conn, conn)
	var item Item
	err = conn.BeginFunc(ctx, func(tx pgx.Tx) error {
		return tx.QueryRow(ctx, "INSERT INTO items (name, user_id) VALUES ($1, $2) RETURNING id, user_id, created_at, updated_at", name, userID).Scan(&item.ID, &item.UserID, &item.CreatedAt, &item.UpdatedAt)
	})
	if err != nil {
		return nil, err
	}
	return &item, nil
}

func (db *db) UpdateItem(ctx context.Context, itemID uint32, name string) (*Item, error) {
	conn, err := stdlib.AcquireConn(db.conn)
	if err != nil {
		return nil, err
	}
	defer stdlib.ReleaseConn(db.conn, conn)
	var item = Item{ID: itemID, Name: name}

	err = conn.BeginFunc(ctx, func(tx pgx.Tx) error {
		return tx.QueryRow(ctx, "UPDATE items SET name=$1, updated_at=NOW() WHERE id=$2 RETURNING user_id, created_at, updated_at", item.Name, item.ID).Scan(&item.UserID, &item.CreatedAt, &item.UpdatedAt)
	})
	if err != nil {
		return nil, err
	}
	return &item, nil
}
