DROP TRIGGER IF EXISTS trg_update_rents_timestamp ON rents;
DROP TRIGGER IF EXISTS trg_update_lockers_timestamp ON lockers;

DROP TABLE IF EXISTS rents;
DROP TABLE IF EXISTS lockers;
DROP TABLE IF EXISTS bloqs;

DROP FUNCTION IF EXISTS update_timestamp();

DROP TYPE IF EXISTS rent_status;
DROP TYPE IF EXISTS rent_size;
DROP TYPE IF EXISTS locker_status;