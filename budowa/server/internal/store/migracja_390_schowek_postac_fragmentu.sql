-- Migracja 390 dodaje kolumnę postaci fragmentu do wpisu historii schowka,
-- aby wklejenie mogło zachować postać źródła.

ALTER TABLE wpis_schowka ADD COLUMN postac_json TEXT;
