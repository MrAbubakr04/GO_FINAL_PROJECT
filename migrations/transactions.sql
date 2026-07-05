CREATE TABLE transactions (
    id BIGSERIAL,
    from_acc_id INT8 NOT NULL,
    to_acc_id INT8 NOT NULL,
    amount INT8 NOT NULL,
    dt_created TIMESTAMP NOT NULL DEFAULT NOW(),

    CONSTRAINT transactions_id_pk PRIMARY KEY (id),
    CONSTRAINT transactions_amount_chk CHECK (amount > 0),

    CONSTRAINT transactions_from_acc_id_fk FOREIGN KEY (from_acc_id)
        REFERENCES accounts(id),

    CONSTRAINT transactions_to_acc_id_fk FOREIGN KEY (to_acc_id)
        REFERENCES accounts(id)
);