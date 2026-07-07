
-------------------- SET_DT_UPDATED -------------------
CREATE OR REPLACE FUNCTION set_dt_updated()
RETURNS TRIGGER AS $$
BEGIN
    IF NEW IS DISTINCT FROM OLD THEN
        NEW.dt_updated := NOW();
    END IF;
    RETURN NEW;
END;
$$ LANGUAGE plpgsql;


----------- UPPER_TRIM_FIO -----------------------------
CREATE OR REPLACE FUNCTION upper_tim_fio()
RETURNS TRIGGER AS $$
BEGIN
    NEW.name := upper(trim(coalesce(NEW.name, '')));
    NEW.surname := upper(trim(coalesce(NEW.surname, '')));
    NEW.fathername := upper(trim(coalesce(NEW.fathername, '')));
    RETURN NEW;
END;
$$ LANGUAGE plpgsql;

----------- NORMALIZE_PHONE_NUM -----------------------------
create or replace function normilize_phone_num()
returns trigger as $$
begin
	IF NEW.phone_num IS NOT NULL THEN
        -- Оставляем только цифры
        NEW.phone_num := regexp_replace(NEW.phone_num, '[^0-9]', '', 'g');
    END IF;
	return new;
end;
$$ language plpgsql;