BEGIN TRANSACTION;

DO $$
BEGIN
   
    CREATE TABLE IF NOT EXISTS metrics (
        ID TEXT UNIQUE,
        MTYPE TEXT,
        DELTA BIGINT,
        VALUE DOUBLE PRECISION
    );

    CREATE INDEX IF NOT EXISTS metric_id ON metrics USING hash(ID);
END $$;
--
--
COMMIT TRANSACTION;