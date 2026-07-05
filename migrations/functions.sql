
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