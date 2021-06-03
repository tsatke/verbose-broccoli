package app

import (
	"database/sql"
	_ "embed"
	"fmt"

	_ "github.com/lib/pq"
)

//go:embed init.sql
var init_sql string

var _ DocumentRepo = (*PostgresDocumentRepo)(nil)

type PostgresDocumentRepo struct {
	db *sql.DB
}

func NewPostgresDocumentRepo(p *PostgresDatabaseProvider) *PostgresDocumentRepo {
	return &PostgresDocumentRepo{
		db: p.DB,
	}
}

func (i *PostgresDocumentRepo) Create(header DocumentHeader, acl ACL) error {
	return tx(i.db, func(tx *sql.Tx) error {
		docHeaderInsert, err := tx.Prepare(`INSERT INTO au_document_headers (doc_id, name, owner, created) VALUES ($1, $2, $3, $4)`)
		if err != nil {
			return fmt.Errorf("prepare header insert: %w", err)
		}
		defer func() {
			_ = docHeaderInsert.Close()
		}()

		_, err = docHeaderInsert.Exec(header.ID, header.Name, header.Owner, header.Created)
		if err != nil {
			return fmt.Errorf("insert header: %w", err)
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
				return fmt.Errorf("insert ACL: %w", err)
			}
		}

		return nil
	})
}

func (i *PostgresDocumentRepo) Get(id DocID) (DocumentHeader, error) {
	row := i.db.QueryRow(`SELECT doc_id, name FROM au_document_headers WHERE doc_id = $1`, id)

	var h DocumentHeader
	if err := row.Scan(&h.ID, &h.Name); err != nil {
		return DocumentHeader{}, fmt.Errorf("scan: %w", err)
	}
	return h, nil
}

func (i *PostgresDocumentRepo) Delete(id DocID) error {
	panic("implement me")
}

func (i *PostgresDocumentRepo) ACL(id DocID) (ACL, error) {
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
