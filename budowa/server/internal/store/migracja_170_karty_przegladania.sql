-- Migracja 170 — karty przeglądania, grupy kart i przestrzenie robocze modułu
-- Browser.
--
-- Opracowanie modułu (rozdz. 2.1) stawia kartę, grupę kart i przestrzeń roboczą
-- obok siebie jako trzy różne byty jednego okna przeglądarki: karta niesie
-- adres, grupa niesie nazwę i barwę zwijanego zestawu, a przestrzeń robocza jest
-- zapisanym kompletem kart przełączanym bez utraty stanu. Trzy byty, trzy
-- tabele — jedna tabela z kolumną „rodzaj" kazałaby każdemu odczytowi odsiewać
-- dwie trzecie wierszy.
--
-- Karta zamknięta zostaje w tabeli. `browser.tab.close` oddaje `closed:true`,
-- a nie „karty nigdy nie było": karta zamknięta jest częścią historii okna
-- i przestrzeni roboczej, która ją zapamiętała. Kolumna `zamknieta` odsiewa ją
-- z wykazu `browser.tab.list`, zamiast kasować wiersz, do którego odwołuje się
-- zapisany zestaw.
--
-- Grupa i przestrzeń nie niosą składu w kolumnie tekstowej. Skład jest po
-- stronie karty (`grupa`, `przestrzen`), bo karta należy do jednej grupy
-- i jednej przestrzeni naraz. Wykaz identyfikatorów przepisany do kolumny JSON
-- byłby drugą prawdą o tej samej przynależności — i rozjechałby się przy
-- pierwszym zamknięciu karty.
--
-- Okno jest kolumną tekstową, nie więzem obcym — tak samo jak w migracji 047:
-- `windowId` modułu Browser jest oknem operacyjnym, nie oknem komunikacji.

-- ── Karta przeglądania ────────────────────────────────────────────────────────
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

-- ── Grupa kart ────────────────────────────────────────────────────────────────
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

-- ── Przestrzeń robocza przeglądania ───────────────────────────────────────────
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
