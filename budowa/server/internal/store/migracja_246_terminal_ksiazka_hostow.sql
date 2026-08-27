-- Migracja 246 zakłada tabelę hostów SSH modułu Terminal, przechowującą trwałe
-- wpisy książki adresowej niezawierające żadnego materiału tajnego.

CREATE TABLE terminal_host (
    id                INTEGER PRIMARY KEY AUTOINCREMENT,
    kod               TEXT    NOT NULL UNIQUE,
    -- Nazwa widoczna w wykazie; klucz rozpoznania dla operatora.
    nazwa             TEXT    NOT NULL,
    -- Adres celu: użytkownik@host albo alias konfiguracji OpenSSH maszyny rdzenia.
    cel               TEXT    NOT NULL,
    -- Port połączenia; pusty przyjmuje port domyślny protokołu.
    port              INTEGER CHECK(port IS NULL OR port BETWEEN 1 AND 65535),
    -- Folder porządkujący wykaz; pusty oznacza brak folderu.
    grupa             TEXT    NOT NULL DEFAULT '',
    -- Katalog roboczy karty zakładanej z tego wpisu.
    katalog_roboczy   TEXT    NOT NULL DEFAULT '',
    -- Wpis wykazu kluczy SSH; pusty przyjmuje klucz domyślny konfiguracji maszyny.
    klucz_kod         TEXT,
    -- Wpis hosta pośredniego, przez który prowadzi połączenie.
    host_posredni_kod TEXT REFERENCES terminal_host(kod) ON DELETE SET NULL,
    -- Notatka operatora albo ślad pochodzenia wpisu.
    notatka           TEXT    NOT NULL DEFAULT '',
    utworzono         TEXT    NOT NULL DEFAULT (strftime('%Y-%m-%dT%H:%M:%fZ','now')),
    zaktualizowano    TEXT    NOT NULL DEFAULT (strftime('%Y-%m-%dT%H:%M:%fZ','now'))
);
CREATE INDEX idx_terminal_host_grupa ON terminal_host(grupa, nazwa);
