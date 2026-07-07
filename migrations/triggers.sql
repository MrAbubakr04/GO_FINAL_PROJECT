
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

----------------- set_normalize_phone_num_to_phones ------------------
CREATE or replace TRIGGER phone_normalize_tr
BEFORE INSERT OR UPDATE OF phone_num
ON phones
FOR EACH ROW
EXECUTE FUNCTION normilize_phone_num();

----------------- ACCOUNTS -----------------------
---------------- set_dates_to_accounts -----------
CREATE TRIGGER set_dates_to_accounts
BEFORE UPDATE ON accounts
FOR EACH ROW
EXECUTE FUNCTION set_dt_updated();

----------------- set_normalize_phone_num_to_accounts -----------
CREATE or replace TRIGGER accounts_phone_num_normalize_tr
BEFORE INSERT OR UPDATE OF phone_num
ON accounts
FOR EACH ROW
EXECUTE FUNCTION normilize_phone_num();