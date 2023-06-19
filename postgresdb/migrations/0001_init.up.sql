BEGIN TRANSACTION;

DO $$
BEGIN
   
    CREATE TABLE IF NOT EXISTS metrics (
             id INT GENERATED ALWAYS AS IDENTITY,
            mname TEXT NOT NULL,
            mtype TEXT NOT NULL,
            delta BIGINT,
            value DOUBLE PRECISION,
            PRIMARY KEY(id),
            UNIQUE(mname, mtype)
    );

    CREATE INDEX IF NOT EXISTS metric_id ON metrics USING hash(id);
END $$;
--
--
COMMIT TRANSACTION;