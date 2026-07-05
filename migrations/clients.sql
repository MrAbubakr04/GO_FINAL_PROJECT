CREATE TABLE clients (
    id BIGSERIAL,
    dt_created TIMESTAMP NOT NULL DEFAULT NOW(),
    dt_updated TIMESTAMP,
    name TEXT,
    surname TEXT,
    fathername TEXT,
    doc_num TEXT NOT NULL,
    tin TEXT NOT NULL,
    birth_date DATE NOT NULL,
    gender VARCHAR(10) NOT NULL,
    address TEXT NOT NULL,
    active_to DATE,

    CONSTRAINT clients_id_pk PRIMARY KEY (id),

    CONSTRAINT clients_fio_len_min_2_chk CHECK (
        (length(coalesce(name, '')) > 0)::int +
        (length(coalesce(surname, '')) > 0)::int +
        (length(coalesce(fathername, '')) > 0)::int >= 2
    )
);

CREATE UNIQUE INDEX clients_tin_un
ON clients (tin)
WHERE active_to IS NULL;

CREATE UNIQUE INDEX clients_tin_history_un
ON clients (tin, active_to)
WHERE active_to IS NOT NULL;

CREATE INDEX clients_doc_num_idx ON clients (doc_num);
CREATE INDEX clients_name_idx ON clients (name);
CREATE INDEX clients_surname_idx ON clients (surname);
CREATE INDEX clients_fathername_idx ON clients (fathername);