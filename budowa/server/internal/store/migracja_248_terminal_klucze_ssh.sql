-- Migracja 248 wprowadza wykaz kluczy SSH znanych rdzeniowi: wiersz niesie ścieżkę klucza
-- prywatnego na maszynie rdzenia, treść klucza publicznego, odcisk oraz znacznik ochrony
-- hasłem, bez materiału tajnego.

CREATE TABLE terminal_klucz (
    id          INTEGER PRIMARY KEY AUTOINCREMENT,
    kod         TEXT    NOT NULL UNIQUE,
    nazwa       TEXT    NOT NULL,
    rodzaj      TEXT    NOT NULL CHECK(rodzaj IN ('ed25519','rsa','ecdsa')),
    -- Odcisk klucza publicznego w postaci używanej przez OpenSSH.
    odcisk      TEXT    NOT NULL DEFAULT '',
    -- Treść klucza publicznego do przeniesienia na host docelowy.
    klucz_jawny TEXT    NOT NULL DEFAULT '',
    -- Ścieżka klucza prywatnego na maszynie rdzenia.
    sciezka     TEXT    NOT NULL UNIQUE,
    -- Czy klucz jest chroniony hasłem; samo hasło nie ma tu kolumny.
    haslo       INTEGER NOT NULL DEFAULT 0 CHECK(haslo IN (0,1)),
    utworzono   TEXT    NOT NULL DEFAULT (strftime('%Y-%m-%dT%H:%M:%fZ','now'))
);
