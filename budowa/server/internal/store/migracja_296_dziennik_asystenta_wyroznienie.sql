-- Migracja 296 dodaje kolumny wyróżnienia oraz powodu wyróżnienia wpisu
-- dziennika asystenta wraz z indeksem wspierającym wykaz.

ALTER TABLE wpis_dziennika_asystenta ADD COLUMN wazny INTEGER NOT NULL DEFAULT 0;
ALTER TABLE wpis_dziennika_asystenta ADD COLUMN notatka_wyroznienia TEXT;

CREATE INDEX idx_wpis_dziennika_asystenta_wazny
    ON wpis_dziennika_asystenta(okno_kod, wazny, utworzono DESC);
