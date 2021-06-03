DROP TABLE IF EXISTS "au_document_acls";
DROP TABLE IF EXISTS "au_document_headers";
DROP TABLE IF EXISTS "au_folders";

CREATE TABLE "au_folders"
(
    "id"    bigserial primary key,
    "path"  text         not null,
    "owner" varchar(255) not null
);

CREATE TABLE "au_document_headers"
(
    "id"        bigserial primary key,        -- the database document ID
    "doc_id"    varchar(255) not null unique, -- the document ID used by the application
    "name"      text         not null,
    "owner"     varchar(255) not null,
    "created"   timestamptz  not null,
    "updated"   timestamptz,                  -- null when there's no content stored yet
    "folder_id" bigserial    not null,

    CONSTRAINT fk_folder_id
        FOREIGN KEY (folder_id)
            REFERENCES au_folders (id)
);

CREATE TABLE "au_document_acls"
(
    "id"       bigserial primary key,
    "doc_id"   bigserial    not null,
    "username" varchar(255) not null,
    "read"     bool         not null,
    "write"    bool         not null,
    "delete"   bool         not null,
    "share"    bool         not null,

    CONSTRAINT fk_doc_id
        FOREIGN KEY (doc_id)
            REFERENCES au_document_headers (id)
);