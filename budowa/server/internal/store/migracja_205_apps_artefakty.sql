-- Migracja 205 zakłada tabelę artefaktów budowania modułu Apps, niosącą odwołanie do archiwum spakowanego przez silnik wykonania wdrożenia.

CREATE TABLE artefakt_apps (
    id                       INTEGER PRIMARY KEY AUTOINCREMENT,
    identyfikator_zewnetrzny TEXT    NOT NULL UNIQUE,
    okno                     TEXT    NOT NULL,
    wdrozenie_kod            TEXT,
    -- Wartości kontraktu (AppArtifactKind) wprost, bez tłumaczenia.
    rodzaj                   TEXT    NOT NULL
                                     CHECK(rodzaj IN ('bundle','binary','container','archive','sourceMap')),
    -- Odwołanie względne magazynu treści rdzenia.
    sciezka                  TEXT    NOT NULL,
    rozmiar                  INTEGER,
    suma_kontrolna           TEXT,
    utworzono                TEXT    NOT NULL DEFAULT (strftime('%Y-%m-%dT%H:%M:%fZ','now'))
);
CREATE INDEX idx_artefakt_apps_okno ON artefakt_apps(okno, id DESC);
CREATE INDEX idx_artefakt_apps_wdrozenie ON artefakt_apps(wdrozenie_kod, id DESC);
