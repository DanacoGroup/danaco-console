-- Migracja 290 dodaje tabele sond kondycji oraz serii ich wyników,
-- opisujących zmierzony stan produktu w czasie.

CREATE TABLE sonda_kondycji (
    id                       INTEGER PRIMARY KEY AUTOINCREMENT,
    identyfikator_zewnetrzny TEXT    NOT NULL UNIQUE,
    nazwa                    TEXT    NOT NULL,
    rodzaj                   TEXT    NOT NULL,
    cel                      TEXT    NOT NULL,
    komponent_kod            TEXT,
    odstep_ms                INTEGER NOT NULL,
    limit_czasu_ms           INTEGER,
    oczekiwany_status        INTEGER,
    tresc_wysylana           TEXT,
    cel_dostepnosci          REAL,
    czynna                   INTEGER NOT NULL DEFAULT 1,
    ostatni_stan             TEXT,
    ostatni_przebieg         INTEGER,
    utworzono                INTEGER NOT NULL,
    zaktualizowano           INTEGER
);
CREATE INDEX idx_sonda_kondycji_rodzaj ON sonda_kondycji(rodzaj, id);
CREATE INDEX idx_sonda_kondycji_komponent ON sonda_kondycji(komponent_kod, id);

-- Wiersz serii jest niezmienny: pomiar raz wykonany nie zmienia się później.
-- Usunięcie sondy zabiera serię ze sobą (`ON DELETE CASCADE`) — seria bez
-- definicji byłaby liczbami bez pytania, na które odpowiadają.
CREATE TABLE wynik_sondy_kondycji (
    id                       INTEGER PRIMARY KEY AUTOINCREMENT,
    identyfikator_zewnetrzny TEXT    NOT NULL UNIQUE,
    sonda_kod                TEXT    NOT NULL
                                     REFERENCES sonda_kondycji(identyfikator_zewnetrzny) ON DELETE CASCADE,
    stan                     TEXT    NOT NULL,
    wykonano                 INTEGER NOT NULL,
    czas_odpowiedzi_ms       INTEGER,
    status_http              INTEGER,
    szczegol                 TEXT,
    blad_kod                 TEXT
);
CREATE INDEX idx_wynik_sondy_kondycji_sonda ON wynik_sondy_kondycji(sonda_kod, wykonano DESC, id DESC);
CREATE INDEX idx_wynik_sondy_kondycji_czas ON wynik_sondy_kondycji(wykonano DESC, id DESC);
