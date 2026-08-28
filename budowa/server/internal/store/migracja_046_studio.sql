-- Migracja zakłada trwałość modułu edytora dokumentów: dokument otwarty, historię jego
-- wersji i propozycję zmiany z operacji kontekstowej.

-- Wartości kolumny formatu są wartościami kontraktu wprost, nie ich tłumaczeniem na
-- osobny słownik bazy.
CREATE TABLE dokument_studio (
    id                       INTEGER PRIMARY KEY AUTOINCREMENT,
    identyfikator_zewnetrzny TEXT    NOT NULL UNIQUE,
    okno                     TEXT    NOT NULL,
    tytul                    TEXT,
    format                   TEXT    NOT NULL DEFAULT 'txt'
                                     CHECK(format IN ('pdf','docx','txt','markdown')),
    tresc                    TEXT,
    tresc_odwolanie          TEXT,
    plik_repozytorium_id     TEXT,
    -- Wskazuje ostatnią wersję identyfikatorem zewnętrznym, bo dokument bywa jeszcze bez wersji.
    wersja_biezaca_id        TEXT,
    utworzono                TEXT    NOT NULL DEFAULT (strftime('%Y-%m-%dT%H:%M:%fZ','now')),
    zaktualizowano           TEXT    NOT NULL DEFAULT (strftime('%Y-%m-%dT%H:%M:%fZ','now'))
);
CREATE INDEX idx_dokument_studio_okno ON dokument_studio(okno, id);

-- Wersja dokumentu jest wpisem repozytorium sesji, do którego zapis z żądaniem
-- wersjonowania dopisuje kolejny wiersz.
CREATE TABLE wersja_dokumentu_studio (
    id                       INTEGER PRIMARY KEY AUTOINCREMENT,
    identyfikator_zewnetrzny TEXT    NOT NULL UNIQUE,
    dokument_id              INTEGER NOT NULL REFERENCES dokument_studio(id) ON DELETE CASCADE,
    etykieta                 TEXT,
    podsumowanie             TEXT,
    skrot_tresci             TEXT,
    tresc                    TEXT,
    tresc_odwolanie          TEXT,
    utworzono                TEXT    NOT NULL DEFAULT (strftime('%Y-%m-%dT%H:%M:%fZ','now'))
);
CREATE INDEX idx_wersja_dokumentu_studio_dokument ON wersja_dokumentu_studio(dokument_id, utworzono DESC, id DESC);

-- Kolumna akcji trzyma pozycję rejestru akcji, która wygenerowała propozycję,
-- odróżniając ją od zwykłej wersji.
CREATE TABLE propozycja_zmiany_studio (
    id                       INTEGER PRIMARY KEY AUTOINCREMENT,
    identyfikator_zewnetrzny TEXT    NOT NULL UNIQUE,
    dokument_id              INTEGER NOT NULL REFERENCES dokument_studio(id) ON DELETE CASCADE,
    akcja_id                 TEXT    NOT NULL,
    wiadomosc_id             TEXT,
    tresc_wyniku             TEXT,
    tresc                    TEXT,
    tresc_odwolanie          TEXT,
    utworzono                TEXT    NOT NULL DEFAULT (strftime('%Y-%m-%dT%H:%M:%fZ','now'))
);
CREATE INDEX idx_propozycja_zmiany_studio_dokument ON propozycja_zmiany_studio(dokument_id, utworzono DESC);
