-- Migracja zakłada tabele harmonogramu automatyki, jej wyzwalaczy dodatkowych oraz
-- przebiegu wykonania danej automatyki.

-- Harmonogram jest jeden na automatykę; więz unikalności na kolumnie automatyki zamyka
-- drogę do dwóch sprzecznych harmonogramów.
CREATE TABLE harmonogram_automatyki (
    id                       INTEGER PRIMARY KEY AUTOINCREMENT,
    identyfikator_zewnetrzny TEXT    NOT NULL UNIQUE,
    automatyka_id            INTEGER NOT NULL UNIQUE
                                     REFERENCES automatyka(id) ON DELETE CASCADE,
    cron                     TEXT,
    strefa_czasowa           TEXT,
    czynny                   INTEGER NOT NULL DEFAULT 1 CHECK(czynny IN (0,1)),
    nastepne_uruchomienie    TEXT,
    utworzono                TEXT    NOT NULL DEFAULT (strftime('%Y-%m-%dT%H:%M:%fZ','now')),
    zaktualizowano           TEXT    NOT NULL DEFAULT (strftime('%Y-%m-%dT%H:%M:%fZ','now'))
);

-- Wyzwalacz niesie zdarzenia wyzwalające poza cyklicznością podstawową: webhook, plik,
-- warunek oraz cykliczności dodatkowe.
CREATE TABLE wyzwalacz_automatyki (
    id                       INTEGER PRIMARY KEY AUTOINCREMENT,
    harmonogram_id           INTEGER NOT NULL
                                     REFERENCES harmonogram_automatyki(id) ON DELETE CASCADE,
    identyfikator_zewnetrzny TEXT    NOT NULL UNIQUE,
    rodzaj                   TEXT    NOT NULL
                                     CHECK(rodzaj IN ('cron','webhook','file','condition')),
    wyrazenie                TEXT    NOT NULL,
    czynny                   INTEGER NOT NULL DEFAULT 1 CHECK(czynny IN (0,1)),
    kolejnosc                INTEGER NOT NULL DEFAULT 0,
    utworzono                TEXT    NOT NULL DEFAULT (strftime('%Y-%m-%dT%H:%M:%fZ','now'))
);
CREATE INDEX idx_wyzwalacz_automatyki_harmonogram
    ON wyzwalacz_automatyki(harmonogram_id, kolejnosc, id);

-- Kolumna próby niesie numer biegu naprawczego; bieg naprawczy nie ma limitu obiegów,
-- więc kolumna nie ma ani górnej granicy, ani więzu.
CREATE TABLE przebieg_automatyki (
    id                       INTEGER PRIMARY KEY AUTOINCREMENT,
    identyfikator_zewnetrzny TEXT    NOT NULL UNIQUE,
    automatyka_id            INTEGER NOT NULL REFERENCES automatyka(id) ON DELETE CASCADE,
    kolejka_id               INTEGER REFERENCES kolejka(id) ON DELETE SET NULL,
    stan                     TEXT    NOT NULL DEFAULT 'pending'
                                     CHECK(stan IN ('pending','running','paused',
                                                    'succeeded','failed','stopped')),
    etap_biezacy             INTEGER NOT NULL DEFAULT 0,
    etapow                   INTEGER NOT NULL DEFAULT 0,
    proba                    INTEGER NOT NULL DEFAULT 0,
    komunikat_bledu          TEXT,
    rozpoczeto               TEXT    NOT NULL DEFAULT (strftime('%Y-%m-%dT%H:%M:%fZ','now')),
    zakonczono               TEXT
);
CREATE INDEX idx_przebieg_automatyki_automatyka
    ON przebieg_automatyki(automatyka_id, id DESC);
CREATE UNIQUE INDEX idx_przebieg_automatyki_kolejka
    ON przebieg_automatyki(kolejka_id) WHERE kolejka_id IS NOT NULL;
