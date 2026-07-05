
CREATE TABLE phones (
    id BIGSERIAL,
    phone_number VARCHAR(15) NOT NULL,
    client_id INT8 NULL,
    dt_created TIMESTAMP NOT NULL DEFAULT NOW(),
    dt_updated TIMESTAMP,
    active_to DATE,

    CONSTRAINT phones_id_pk PRIMARY KEY (id),
    CONSTRAINT phones_client_id_fk FOREIGN KEY (client_id) REFERENCES clients(id)
);

CREATE UNIQUE INDEX phones_phone_number_un
ON phones (phone_number)
WHERE active_to IS NULL;

CREATE UNIQUE INDEX phones_phone_num_hist_un
ON phones (phone_number, active_to)
WHERE active_to IS NOT NULL;