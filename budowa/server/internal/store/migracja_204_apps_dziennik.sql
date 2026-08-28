-- Migracja 204 zakłada wspólną tabelę wierszy dziennika usług i wdrożeń modułu Apps, zapisywaną przez silnik wykonania wdrożenia i serwer podglądu.

CREATE TABLE wiersz_dziennika_apps (
    id            INTEGER PRIMARY KEY AUTOINCREMENT,
    okno          TEXT    NOT NULL,
    wdrozenie_kod TEXT,
    komponent_kod TEXT,
    chwila        INTEGER NOT NULL,
    tresc         TEXT    NOT NULL
);
CREATE INDEX idx_wiersz_dziennika_apps_okno ON wiersz_dziennika_apps(okno, chwila, id);
CREATE INDEX idx_wiersz_dziennika_apps_wdrozenie ON wiersz_dziennika_apps(wdrozenie_kod, chwila, id);
