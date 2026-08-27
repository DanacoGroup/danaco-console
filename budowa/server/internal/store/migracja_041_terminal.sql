-- Migracja zakłada trwałość modułu terminala: kartę powłoki odtwarzalną po restarcie
-- rdzenia oraz dziennik procesów zakończonych.

-- Kolumna kodu okna niesie identyfikator okna komunikacji w postaci kontraktu, bez klucza
-- obcego, bo okno bywa bytem samej pamięci nadzorcy.
CREATE TABLE terminal_karta (
    id              INTEGER PRIMARY KEY AUTOINCREMENT,
    kod             TEXT    NOT NULL UNIQUE,
    okno_kod        TEXT    NOT NULL,
    powloka         TEXT    NOT NULL
                            CHECK(powloka IN ('powershell','cmd','bash','node','python','ssh')),
    tytul           TEXT,
    katalog_roboczy TEXT,
    stan            TEXT    NOT NULL DEFAULT 'running' CHECK(stan IN ('running','exited')),
    utworzono       TEXT    NOT NULL DEFAULT (strftime('%Y-%m-%dT%H:%M:%fZ','now'))
);
CREATE INDEX idx_terminal_karta_okno ON terminal_karta(okno_kod, utworzono DESC);

-- Tabela niesie ten sam wiersz co stan żywy po to, żeby podgląd procesów pokazał
-- również procesy już zakończone.
CREATE TABLE terminal_proces (
    id            INTEGER PRIMARY KEY AUTOINCREMENT,
    kod           TEXT    NOT NULL UNIQUE,
    karta_kod     TEXT,
    okno_kod      TEXT    NOT NULL,
    pid           INTEGER,
    pid_nadrzedny INTEGER,
    polecenie     TEXT    NOT NULL,
    inicjator     TEXT    NOT NULL DEFAULT 'operator'
                          CHECK(inicjator IN ('operator','model')),
    stan          TEXT    NOT NULL DEFAULT 'running'
                          CHECK(stan IN ('running','finished','failed','stopped')),
    kod_wyjscia   INTEGER,
    uruchomiono   TEXT    NOT NULL DEFAULT (strftime('%Y-%m-%dT%H:%M:%fZ','now')),
    zakonczono    TEXT
);
CREATE INDEX idx_terminal_proces_okno ON terminal_proces(okno_kod, uruchomiono DESC);
CREATE INDEX idx_terminal_proces_karta ON terminal_proces(karta_kod, uruchomiono DESC);
CREATE INDEX idx_terminal_proces_stan ON terminal_proces(stan, uruchomiono DESC);
