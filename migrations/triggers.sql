
---------------- CLIENTS -------------------------
--------------- set_dates_to_clients -------------
CREATE TRIGGER set_dates_to_clients
BEFORE UPDATE ON clients
FOR EACH ROW
EXECUTE FUNCTION set_dt_updated();

--------------- set_normalize_fio_to_clients -----
CREATE TRIGGER set_normalize_fio_to_clients
BEFORE INSERT OR UPDATE ON clients
FOR EACH ROW
EXECUTE FUNCTION upper_tim_fio();

--------------------- PHONES ---------------------------
----------------- set_dates_to_phones ------------------
CREATE TRIGGER set_dates_to_phones
BEFORE UPDATE ON phones
FOR EACH ROW
EXECUTE FUNCTION set_dt_updated();

----------------- ACCOUNTS -----------------------
---------------- set_dates_to_accounts -----------
CREATE TRIGGER set_dates_to_accounts
BEFORE UPDATE ON accounts
FOR EACH ROW
EXECUTE FUNCTION set_dt_updated();