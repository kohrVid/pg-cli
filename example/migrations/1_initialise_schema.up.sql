CREATE TABLE IF NOT EXISTS table0 (
  field VARCHAR(16)
);

CREATE OR REPLACE FUNCTION modify_updated_at_column()
  RETURNS TRIGGER AS $$
  BEGIN
    NEW.updated_at = now();
    RETURN NEW;
  END;
$$ language 'plpgsql';

DROP TABLE table0;
