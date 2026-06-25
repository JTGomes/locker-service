CREATE EXTENSION IF NOT EXISTS pgcrypto;

CREATE TABLE IF NOT EXISTS bloqs (
    id         UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    title      VARCHAR(255) NOT NULL,
    address    VARCHAR(500) NOT NULL
);

CREATE TYPE locker_status AS ENUM (
  'open',
  'closed'
);

CREATE TABLE IF NOT EXISTS lockers (
    id          UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    bloq_id     UUID NOT NULL REFERENCES bloqs(id) ON DELETE RESTRICT,
    status      locker_status NOT NULL DEFAULT 'closed',
    is_occupied BOOLEAN NOT NULL DEFAULT false,
    created_at  TIMESTAMPTZ NOT NULL DEFAULT now(),
    updated_at  TIMESTAMPTZ NOT NULL DEFAULT now()
);

CREATE INDEX idx_lockers_bloq_id ON lockers(bloq_id);

CREATE TYPE rent_status AS ENUM (
  'created', 
  'waiting_dropoff', 
  'waiting_pickup', 
  'delivered'
);

CREATE TYPE rent_size AS ENUM (
  'XS',
  'S',
  'M',
  'L',
  'XL'
);

CREATE TABLE IF NOT EXISTS rents (
    id              UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    locker_id       UUID NULL REFERENCES lockers(id) ON DELETE RESTRICT,
    weight          DOUBLE PRECISION NOT NULL CHECK (weight > 0),
    size            rent_size NOT NULL,
    status          rent_status NOT NULL DEFAULT 'created',
    created_at      TIMESTAMPTZ NOT NULL DEFAULT now(),
    updated_at      TIMESTAMPTZ NOT NULL DEFAULT now(),
    dropped_off_at  TIMESTAMPTZ NULL,
	  picked_up_at    TIMESTAMPTZ NULL
);

CREATE INDEX idx_rents_locker_id ON rents(locker_id);
CREATE UNIQUE INDEX uq_active_rent_per_locker ON rents (locker_id)
  WHERE status IN ('waiting_dropoff', 'waiting_pickup');


CREATE OR REPLACE FUNCTION update_timestamp()
RETURNS TRIGGER AS $$
BEGIN
    NEW.updated_at = NOW();
    RETURN NEW;
END;
$$ LANGUAGE plpgsql;

CREATE TRIGGER trg_update_lockers_timestamp
    BEFORE UPDATE ON lockers
    FOR EACH ROW
    EXECUTE FUNCTION update_timestamp();

CREATE TRIGGER trg_update_rents_timestamp
    BEFORE UPDATE ON rents
    FOR EACH ROW
    EXECUTE FUNCTION update_timestamp();
