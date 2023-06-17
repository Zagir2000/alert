BEGIN TRANSACTION;

DO $$
BEGIN
    IF EXISTS (
        SELECT 
            * 
        FROM information_schema.tables 
        WHERE table_name = __metrics  AND table_schema = 'public'
    )

    THEN
        ALTER TABLE __metrics RENAME TO metrics;
    ELSE
        CREATE TABLE IF NOT EXISTS metrics (
            ID TEXT UNIQUE,
            MTYPE TEXT,
            DELTA BIGINT,
            VALUE DOUBLE PRECISION
        );
    END IF;
    CREATE INDEX IF NOT EXISTS metric_id ON metrics USING hash(ID);
END $$;
--
--
COMMIT TRANSACTION;