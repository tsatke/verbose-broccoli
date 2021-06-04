package app

import (
	"database/sql"
	"database/sql/driver"
	"fmt"
	"time"
)

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

func (i *PostgresDocumentRepo) Update(header DocumentHeader, acl ACL) error {
	return tx(i.db, func(tx *sql.Tx) error {
		docHeaderUpdate, err := tx.Prepare(`UPDATE au_document_headers SET (name, owner, created, updated) = ($1, $2, $3, $4) WHERE doc_id = $5`)
		if err != nil {
			return fmt.Errorf("prepare header insert: %w", err)
		}
		defer func() {
			_ = docHeaderUpdate.Close()
		}()

		_, err = docHeaderUpdate.Exec(header.Name, header.Owner, header.Created, header.Updated, header.ID)
		if err != nil {
			return fmt.Errorf("update header: %w", err)
		}

		docACLInsert, err := tx.Prepare(`UPDATE au_document_acls SET (username, read, write, delete, share) = ($1, $2, $3, $4, $5) WHERE doc_id = $6`)
		if err != nil {
			return fmt.Errorf("prepare acl insert: %w", err)
		}
		defer func() {
			_ = docACLInsert.Close()
		}()

		for _, perm := range acl.Permissions {
			_, err = docACLInsert.Exec(perm.Username, perm.Read, perm.Write, perm.Delete, perm.Share, header.ID)
			if err != nil {
				return fmt.Errorf("update ACL: %w", err)
			}
		}

		return nil
	})
}

func (i *PostgresDocumentRepo) Get(id DocID) (DocumentHeader, error) {
	row := i.db.QueryRow(`SELECT doc_id, name, owner, created, updated FROM au_document_headers WHERE doc_id = $1`, id)

	var h DocumentHeader
	var nt nullableTime
	if err := row.Scan(&h.ID, &h.Name, &h.Owner, &h.Created, &nt); err != nil {
		return DocumentHeader{}, fmt.Errorf("scan: %w", err)
	}
	if nt.Valid {
		h.Updated = nt.Time
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

type nullableTime struct {
	Time  time.Time
	Valid bool // Valid is true if Time is not NULL
}

// Scan implements the Scanner interface.
func (nt *nullableTime) Scan(value interface{}) error {
	nt.Time, nt.Valid = value.(time.Time)
	return nil
}

// Value implements the driver Valuer interface.
func (nt nullableTime) Value() (driver.Value, error) {
	if !nt.Valid {
		return nil, nil
	}
	return nt.Time, nil
}
