-- Migracja 175 — wytwory sesji przeglądania i zrzuty stron
-- (`browser.artifact.add`, `browser.screenshot.capture`,
-- `browser.snapshot.screenshot.get`).
--
-- Wytwór jest wynikiem czynności Operatora na stronie: zrzutem, archiwum,
-- wyodrębnionymi danymi, adnotacją albo rejestrem sieciowym. Bajty leżą
-- w magazynie treści pod sumą kontrolną, a tabela trzyma odwołanie — tak samo
-- jak treść biblioteki i zasoby modułu Design. Wiersz bez odwołania nie ma
-- prawa powstać: wytwór, za którym nie ma ani jednego bajtu, jest dokładnie tą
-- szkodą, przed którą stoją sprawdziany skutku tego produktu.
--
-- Zrzut ma własną tabelę obok wytworu, bo niesie pomiary, których wytwór nie
-- zna: tryb (widok, cała strona, obszar, element), format pliku oraz wymiary
-- w pikselach. Wciśnięcie ich w kolumnę JSON wytworu odebrałoby możliwość
-- odczytu zrzutu po odwołaniu (`browser.snapshot.screenshot.get` pyta
-- `screenshotRef` albo `snapshotId`, a nie identyfikatorem wytworu).
CREATE TABLE wytwor_przegladania (
    id                       INTEGER PRIMARY KEY AUTOINCREMENT,
    identyfikator_zewnetrzny TEXT    NOT NULL UNIQUE,
    okno                     TEXT    NOT NULL,
    rodzaj                   TEXT    NOT NULL,
    tytul                    TEXT,
    tresc_odwolanie          TEXT    NOT NULL,
    typ_mime                 TEXT,
    rozmiar_bajtow           INTEGER,
    url_zrodla               TEXT,
    migawka_zewnetrzna_id    TEXT,
    utworzono                TEXT    NOT NULL DEFAULT (strftime('%Y-%m-%dT%H:%M:%fZ','now'))
);
CREATE INDEX idx_wytwor_przegladania_okno ON wytwor_przegladania(okno, utworzono DESC, id);

CREATE TABLE zrzut_przegladania (
    id                       INTEGER PRIMARY KEY AUTOINCREMENT,
    identyfikator_zewnetrzny TEXT    NOT NULL UNIQUE,
    okno                     TEXT    NOT NULL,
    migawka_zewnetrzna_id    TEXT,
    tresc_odwolanie          TEXT    NOT NULL,
    tryb                     TEXT    NOT NULL,
    format                   TEXT    NOT NULL,
    szerokosc                INTEGER NOT NULL DEFAULT 0,
    wysokosc                 INTEGER NOT NULL DEFAULT 0,
    rozmiar_bajtow           INTEGER,
    utworzono                TEXT    NOT NULL DEFAULT (strftime('%Y-%m-%dT%H:%M:%fZ','now'))
);
CREATE INDEX idx_zrzut_przegladania_okno ON zrzut_przegladania(okno, utworzono DESC, id);
CREATE INDEX idx_zrzut_przegladania_odwolanie ON zrzut_przegladania(tresc_odwolanie);
