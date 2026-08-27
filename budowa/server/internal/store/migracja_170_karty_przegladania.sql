-- Migracja 170 zakłada tabele kart przeglądania, grup kart i przestrzeni roboczych modułu Browser jako trzech odrębnych bytów jednego okna.

-- Zakłada tabelę karta_przegladania niosącą adres, stan i przynależność karty do grupy oraz przestrzeni roboczej okna.
CREATE TABLE karta_przegladania (
    id                       INTEGER PRIMARY KEY AUTOINCREMENT,
    identyfikator_zewnetrzny TEXT    NOT NULL UNIQUE,
    okno                     TEXT    NOT NULL,
    url                      TEXT,
    tytul                    TEXT,
    stan                     TEXT    NOT NULL DEFAULT 'inactive',
    przypieta                INTEGER NOT NULL DEFAULT 0 CHECK(przypieta IN (0,1)),
    zamknieta                INTEGER NOT NULL DEFAULT 0 CHECK(zamknieta IN (0,1)),
    grupa                    TEXT,
    przestrzen               TEXT,
    karta_otwierajaca        TEXT,
    kolejnosc                INTEGER NOT NULL DEFAULT 0,
    ostatnio_czynna          TEXT,
    utworzono                TEXT    NOT NULL DEFAULT (strftime('%Y-%m-%dT%H:%M:%fZ','now'))
);
CREATE INDEX idx_karta_przegladania_okno ON karta_przegladania(okno, zamknieta, kolejnosc, id);

-- Zakłada tabelę grupa_kart_przegladania niosącą nazwę i barwę zwijanego zestawu kart w obrębie okna przeglądarki.
CREATE TABLE grupa_kart_przegladania (
    id                       INTEGER PRIMARY KEY AUTOINCREMENT,
    identyfikator_zewnetrzny TEXT    NOT NULL UNIQUE,
    okno                     TEXT    NOT NULL,
    nazwa                    TEXT    NOT NULL,
    barwa                    TEXT,
    zwinieta                 INTEGER NOT NULL DEFAULT 0 CHECK(zwinieta IN (0,1)),
    utworzono                TEXT    NOT NULL DEFAULT (strftime('%Y-%m-%dT%H:%M:%fZ','now'))
);
CREATE INDEX idx_grupa_kart_przegladania_okno ON grupa_kart_przegladania(okno, utworzono DESC, id);

-- Zakłada tabelę przestrzen_przegladania niosącą zapisany komplet kart przełączany bez utraty stanu, z profilem i oknem operacyjnym.
CREATE TABLE przestrzen_przegladania (
    id                       INTEGER PRIMARY KEY AUTOINCREMENT,
    identyfikator_zewnetrzny TEXT    NOT NULL UNIQUE,
    okno                     TEXT,
    nazwa                    TEXT    NOT NULL,
    profil                   TEXT,
    utworzono                TEXT    NOT NULL DEFAULT (strftime('%Y-%m-%dT%H:%M:%fZ','now')),
    zaktualizowano           TEXT    NOT NULL DEFAULT (strftime('%Y-%m-%dT%H:%M:%fZ','now'))
);
CREATE INDEX idx_przestrzen_przegladania_okno ON przestrzen_przegladania(okno, utworzono DESC, id);
