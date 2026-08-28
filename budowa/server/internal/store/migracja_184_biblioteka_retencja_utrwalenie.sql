-- Migracja 184 — moduł Library: polityki przechowywania i utrwalenie archiwalne.
--
-- Polityka retencji nie usuwa niczego sama. Przechowuje trzy rzeczy: poziom
-- zasięgu, na którym obowiązuje, liczbę dni przechowywania liczoną od ostatniej
-- zmiany zasobu oraz czynność po upływie okresu. Czynność `usuniecie` znaczy
-- „zgłoś do usunięcia", nie „usuń": trwałe usunięcie ma własną komendę
-- i własne potwierdzenie (`library.file.delete`). Bez tego rozdziału polityka
-- zapisana pomyłkowo zabierałaby zasoby bez śladu decyzji człowieka.
--
-- Poziom zasięgu zapisuje się kodem z katalogu `poziom_zasiegu` (kontrakt:
-- `ConfigScope`), więc warunku CHECK tu nie ma — wykaz poziomów należy do
-- tabeli katalogu, a jego powielenie w warunku byłoby drugą prawdą o zasięgu.
--
-- Zadanie utrwalenia jest zapisem SKUTKU, nie zleceniem do wykonania. Wiersz
-- powstaje po pracy: niesie wynik walidacji i jej zapis, żeby Operator mógł
-- wrócić do pytania „czy ten dokument naprawdę przeszedł normalizację" bez
-- powtarzania utrwalenia. Zasób wytworzony wskazywany jest kodem, bo bywa go
-- brak — profil PREMIS/METS opisuje zasób, nie wytwarza nowego pliku.

CREATE TABLE polityka_retencji_biblioteki (
    id                       INTEGER PRIMARY KEY AUTOINCREMENT,
    identyfikator_zewnetrzny TEXT    NOT NULL UNIQUE,
    zasieg                   TEXT    NOT NULL,
    zasieg_id                TEXT,
    dni_przechowywania       INTEGER NOT NULL,
    czynnosc                 TEXT    NOT NULL
                             CHECK(czynnosc IN ('przeglad','archiwizacja','usuniecie')),
    utworzono                TEXT    NOT NULL DEFAULT (strftime('%Y-%m-%dT%H:%M:%fZ','now'))
);
CREATE INDEX idx_polityka_retencji_biblioteki_zasieg
    ON polityka_retencji_biblioteki(zasieg, zasieg_id);

CREATE TABLE zadanie_utrwalenia_biblioteki (
    id                       INTEGER PRIMARY KEY AUTOINCREMENT,
    identyfikator_zewnetrzny TEXT    NOT NULL UNIQUE,
    plik_kod                 TEXT    NOT NULL,
    rodzaj                   TEXT    NOT NULL CHECK(rodzaj IN ('pdfa','bagit','premis')),
    plik_wynikowy_kod        TEXT,
    profil                   TEXT,
    poprawne                 INTEGER NOT NULL DEFAULT 0 CHECK(poprawne IN (0,1)),
    -- Zapis walidacji: co sprawdzono i z jakim wynikiem. Wynik bez zapisu byłby
    -- werdyktem bez uzasadnienia.
    raport                   TEXT    NOT NULL DEFAULT '',
    utworzono                TEXT    NOT NULL DEFAULT (strftime('%Y-%m-%dT%H:%M:%fZ','now'))
);
CREATE INDEX idx_zadanie_utrwalenia_biblioteki_plik
    ON zadanie_utrwalenia_biblioteki(plik_kod, utworzono DESC);
