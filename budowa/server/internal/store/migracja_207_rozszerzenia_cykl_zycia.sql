-- Migracja 207 tworzy cztery miejsca cyklu życia pozycji katalogu: kolekcje
-- kuratorskie, dziennik zdarzeń, wersje z przypięciem oraz paczki przesłane
-- instalacją personal.

-- Tabela kolekcja_rozszerzen przechowuje nazwany zestaw pozycji katalogu,
-- wymieniany w całości przy każdym zapisie, nie etykietę pojedynczej pozycji.
CREATE TABLE kolekcja_rozszerzen (
    id                       INTEGER PRIMARY KEY AUTOINCREMENT,
    identyfikator_zewnetrzny TEXT    NOT NULL UNIQUE,
    nazwa                    TEXT    NOT NULL,
    opis                     TEXT,
    -- Oznaczenie barwne kolekcji podane przez operatora; rdzeń go nie interpretuje.
    oznaczenie_barwne        TEXT,
    zaktualizowano           INTEGER NOT NULL DEFAULT 0
);

CREATE TABLE pozycja_kolekcji_rozszerzen (
    kolekcja_id      INTEGER NOT NULL REFERENCES kolekcja_rozszerzen(id) ON DELETE CASCADE,
    -- Wartość danych, nie więz obcy: kolekcja może wskazywać pozycję już odinstalowaną.
    rozszerzenie_kod TEXT    NOT NULL,
    kolejnosc        INTEGER NOT NULL DEFAULT 0,
    PRIMARY KEY (kolekcja_id, rozszerzenie_kod)
);

-- Tabela historia_rozszerzenia jest dziennikiem zdarzeń cyklu życia pozycji,
-- nie stanem: każdy wiersz niesie czynność, wersję przed i po oraz czas.
CREATE TABLE historia_rozszerzenia (
    id                       INTEGER PRIMARY KEY AUTOINCREMENT,
    identyfikator_zewnetrzny TEXT    NOT NULL UNIQUE,
    rozszerzenie_kod         TEXT    NOT NULL,
    czynnosc                 TEXT    NOT NULL
                                     CHECK(czynnosc IN ('installed','updated','enabled','disabled',
                                                        'rolledBack','uninstalled','configured')),
    wersja_przed             TEXT,
    wersja_po                TEXT,
    szczegol                 TEXT,
    -- Milisekundy epoki — `ExtensionHistoryEntry.occurredAt` niesie je wprost.
    zaszlo                   INTEGER NOT NULL
);
CREATE INDEX idx_historia_rozszerzenia ON historia_rozszerzenia(rozszerzenie_kod, zaszlo DESC, id DESC);

-- Tabela wersja_rozszerzenia trzyma numer wersji, dziennik zmian i odwołanie
-- do paczki, żeby cofnięcie było przywróceniem, a nie przepisaniem napisu.
CREATE TABLE wersja_rozszerzenia (
    id               INTEGER PRIMARY KEY AUTOINCREMENT,
    rozszerzenie_kod TEXT    NOT NULL,
    wersja           TEXT    NOT NULL,
    dziennik_zmian   TEXT,
    -- Odwołanie do paczki, z której wersja powstała; puste dla wersji
    -- zarejestrowanej bez przesyłki.
    paczka_odwolanie TEXT,
    utworzono        INTEGER NOT NULL DEFAULT 0,
    UNIQUE (rozszerzenie_kod, wersja)
);

-- Kolumna wersja_przypieta w tabeli rozszerzenie niesie jedną przypiętą wersję
-- pozycji; pusta wartość znaczy brak przypięcia.
ALTER TABLE rozszerzenie ADD COLUMN wersja_przypieta TEXT;

-- Tabela paczka_rozszerzenia przechowuje paczkę przesłaną instalacją personal
-- wraz z odwołaniem do pliku na dysku i sumą kontrolną.
CREATE TABLE paczka_rozszerzenia (
    id                       INTEGER PRIMARY KEY AUTOINCREMENT,
    -- uploadRef oddawany operatorowi i przyjmowany z powrotem przy instalacji.
    identyfikator_zewnetrzny TEXT    NOT NULL UNIQUE,
    nazwa_pliku              TEXT    NOT NULL,
    -- Odwołanie względne do pliku paczki w magazynie treści rdzenia.
    sciezka                  TEXT    NOT NULL,
    rozmiar                  INTEGER NOT NULL DEFAULT 0,
    suma_kontrolna           TEXT    NOT NULL,
    utworzono                INTEGER NOT NULL DEFAULT 0
);
