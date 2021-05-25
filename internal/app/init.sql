DROP TABLE IF EXISTS "au_document_acls";
DROP TABLE IF EXISTS "au_document_headers";

CREATE TABLE "au_document_headers"
(
    "id"     bigserial primary key,
    "doc_id" varchar(255) not null unique,
    "name"   text         not null,
    "size"   bigint       not null
);

CREATE TABLE "au_document_acls"
(
    "id"       bigserial primary key,
    "doc_id"   varchar(255) not null,
    "username" varchar(255) not null,
    "read"     bool         not null,
    "write"    bool         not null,
    "delete"   bool         not null,
    "share"    bool         not null,

    CONSTRAINT fk_doc_id
        FOREIGN KEY (doc_id)
            REFERENCES au_document_headers (doc_id)
);