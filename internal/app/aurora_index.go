package app

import (
	"database/sql"
	_ "embed"
	"fmt"

	_ "github.com/lib/pq"
	"github.com/tsatke/verbose-broccoli/internal/app/config"
)

//go:embed init.sql
var init_sql string

type AuroraIndex struct {
	db *sql.DB
}

func NewAuroraIndex(cfg config.Config) (*AuroraIndex, error) {
	endpoint := cfg.GetString(config.AWSAuroraEndpoint)
	port := cfg.GetString(config.AWSAuroraPort)
	user := cfg.GetString(config.AWSAuroraUsername)
	pass := cfg.GetString(config.AWSAuroraPassword)
	database := cfg.GetString(config.AWSAuroraDatabase)

	dsn := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s", endpoint, port, user, pass, database)

	db, err := sql.Open("postgres", dsn)
	if err != nil {
		return nil, fmt.Errorf("sql open: %w", err)
	}

	i := &AuroraIndex{
		db: db,
	}
	if err := i.initDB(); err != nil {
		return nil, fmt.Errorf("init: %w", err)
	}

	return i, nil
}

func (i *AuroraIndex) tx(fn func(tx *sql.Tx) error) error {
	tx, err := i.db.Begin()
	if err != nil {
		return fmt.Errorf("begin tx: %w", err)
	}
	defer func() {
		_ = tx.Rollback() // will be ignored if tx was already committed
	}()

	if err := fn(tx); err != nil {
		return err
	}

	if err := tx.Commit(); err != nil {
		return fmt.Errorf("commit tx: %w", err)
	}

	return nil
}

func (i *AuroraIndex) initDB() error {
	return i.tx(func(tx *sql.Tx) error {
		_, err := tx.Exec(init_sql)
		if err != nil {
			return fmt.Errorf("exec init: %w", err)
		}
		return nil
	})
}

func (i *AuroraIndex) Create(header DocumentHeader, acl ACL) error {
	return i.tx(func(tx *sql.Tx) error {
		docHeaderInsert, err := tx.Prepare(`INSERT INTO au_document_headers (doc_id, name, size) VALUES ($1, $2, $3)`)
		if err != nil {
			return fmt.Errorf("prepare header insert: %w", err)
		}
		defer func() {
			_ = docHeaderInsert.Close()
		}()

		_, err = docHeaderInsert.Exec(header.ID, header.Name, header.Size)
		if err != nil {
			return err
		}

		docACLInsert, err := tx.Prepare(`INSERT INTO au_document_acls (doc_id, username, read, write, delete, share) VALUES ($1, $2, $3, $4, $5, $6)`)
		if err != nil {
			return fmt.Errorf("prepare acl insert: %w", err)
		}
		defer func() {
			_ = docACLInsert.Close()
		}()

		for _, perm := range acl.Permissions {
			_, err = docACLInsert.Exec(header.ID, perm.Username, perm.Read, perm.Write, perm.Delete, perm.Share)
			if err != nil {
				return err
			}
		}

		return nil
	})
}

func (i *AuroraIndex) GetByID(id string) (DocumentHeader, error) {
	row := i.db.QueryRow(`SELECT doc_id, name, size FROM au_document_headers WHERE doc_id = $1`, id)

	var h DocumentHeader
	if err := row.Scan(&h.ID, &h.Name, &h.Size); err != nil {
		return DocumentHeader{}, fmt.Errorf("scan: %w", err)
	}
	return h, nil
}

func (i *AuroraIndex) Delete(id string) error {
	panic("implement me")
}

func (i *AuroraIndex) ACL(id string) (ACL, error) {
	rows, err := i.db.Query(`SELECT username, read, write, delete, share FROM au_document_acls WHERE doc_id = $1`, id)
	if err != nil {
		return ACL{}, fmt.Errorf("get ACL: %w", err)
	}

	acl := ACL{
		Permissions: map[string]Permission{},
	}
	for rows.Next() {
		var p Permission
		if err := rows.Scan(&p.Username, &p.Read, &p.Write, &p.Delete, &p.Share); err != nil {
			return ACL{}, fmt.Errorf("scan: %w", err)
		}
		acl.Permissions[p.Username] = p
	}
	return acl, nil
}

func (i *AuroraIndex) Close() error {
	return i.db.Close()
}
