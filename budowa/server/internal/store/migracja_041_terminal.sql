-- Migracja 041 — trwałość modułu Terminal: karta powłoki (okno Terminal Tabs)
-- i dziennik procesów rejestru rdzenia (okno Process Monitor).
--
-- Dlaczego karta ma wiersz, skoro istnieje w pamięci rdzenia. Kontrakt nie ma
-- komendy `terminal.session.list`, więc klient po ponownym połączeniu zna
-- identyfikatory kart wyłącznie z własnej pamięci. Gdyby karta żyła tylko
-- w pamięci procesu rdzenia, restart serwera unieważniłby te identyfikatory
-- i `terminal.command.exec` odpowiadałby `not_found` na kartę, którą Operator
-- widzi na ekranie. Wiersz karty pozwala rdzeniowi ją odtworzyć.
--
-- Czego w karcie nie ma.
--   * Zmienne środowiska karty nie mają kolumny. Środowisko procesu
--     bywa nośnikiem poświadczeń, a baza nie jest sejfem; zmienne żyją wyłącznie
--     w pamięci rdzenia i po restarcie karta wraca bez nich.
--   * PID i kod wyjścia karty nie mają kolumny, bo karta nie jest procesem.
--     Kartą jest profil powłoki (rodzaj, katalog, środowisko), a procesem — każde
--     wykonane w niej polecenie; PID i kod wyjścia należą do `terminal_proces`.
--     Pola `TerminalSession.pid` i `.exitCode` kontraktu są opcjonalne i zostają
--     puste.
--
-- Dziennik procesów jest dziennikiem, nie stanem żywym. Proces czynny prowadzi
-- rejestr w pamięci rdzenia (to on ma uchwyt do drzewa potomstwa i tylko on
-- potrafi proces ubić). Tabela niesie ten sam wiersz po to, żeby Process Monitor
-- pokazał także procesy zakończone — kontrakt pozwala zawęzić wykaz do stanów
-- `finished`, `failed` i `stopped`, a stan żywy takich wierszy nie trzyma.

-- ── Karta terminala — profil powłoki jednego okna (Terminal Tabs) ─────────────
-- Kolumna `okno_kod` niesie identyfikator okna komunikacji w postaci kontraktu,
-- ten sam, którym posługuje się `TerminalSession.windowId`. Klucz obcy do
-- `okno_komunikacji` nie wchodzi: okno bywa bytem pamięci nadzorcy, który nie
-- ma jeszcze wiersza, a karta nie ma prawa nie powstać z tego powodu.
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

-- ── Dziennik procesów rejestru rdzenia (Process Monitor) ─────────────────────
-- `inicjator` jest kolumną, a nie wnioskiem z okoliczności: trzeba odróżnić
-- proces uruchomiony przez Operatora od uruchomionego poleceniem modelu,
-- a po zakończeniu procesu nie ma już skąd tego odczytać.
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
