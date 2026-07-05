----------- ACCOUNT_STATUSES ------------------------------

CREATE TABLE account_statuses (
    id SERIAL,
    code TEXT NOT NULL,
    description TEXT NOT NULL,

    CONSTRAINT account_statuses_id_pk PRIMARY KEY (id),
    CONSTRAINT account_statuses_code_uk UNIQUE (code)
);

------------ ACCOUNTS -----------------------------------------------
CREATE TABLE accounts (
    id BIGSERIAL,
    phone_num VARCHAR(15) NOT NULL,
    dt_created TIMESTAMP NOT NULL DEFAULT NOW(),
    dt_updated TIMESTAMP,
    pin TEXT NOT NULL,

    balance_tj INT8 NOT NULL,
    balance_ru INT8 NOT NULL,
    balance_en INT8 NOT NULL,

    device TEXT NOT NULL,
    is_active BOOLEAN DEFAULT TRUE,
    status_id INT4 NOT NULL,
    active_to DATE,

    CONSTRAINT accounts_id_pk PRIMARY KEY (id),
    CONSTRAINT accounts_status_id_fk FOREIGN KEY (status_id)
        REFERENCES account_statuses(id),

    CONSTRAINT accounts_balance_tj_chk CHECK (balance_tj >= 0),
    CONSTRAINT accounts_balance_ru_chk CHECK (balance_ru >= 0),
    CONSTRAINT accounts_balance_en_chk CHECK (balance_en >= 0)
);

CREATE UNIQUE INDEX accounts_phone_num_un
ON accounts (phone_num)
WHERE active_to IS NULL;

CREATE UNIQUE INDEX accounts_phone_num_hist_un
ON accounts (phone_num, active_to)
WHERE active_to IS NOT NULL;