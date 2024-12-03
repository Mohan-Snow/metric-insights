CREATE TABLE IF NOT EXISTS data (
    id BIGSERIAL PRIMARY KEY,
    test_data VARCHAR(255)
);

INSERT INTO data (test_data) SELECT 'SOME GENERATED DATA' FROM generate_series(1, 1000000);

select count(*) from data;

CREATE OR REPLACE PROCEDURE delay() AS $$
BEGIN
    -- sleeps 10 sec
    PERFORM pg_sleep(10);
END;
$$ LANGUAGE plpgsql;

CALL delay();