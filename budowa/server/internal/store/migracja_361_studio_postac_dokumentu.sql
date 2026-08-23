-- Migracja 361 — postać dokumentu modułu Studio: arkusz stylów i sekcje.
--
-- ── Dlaczego postać dokumentu jest bytem trwałym ─────────────────────────────
-- Do tej pory rdzeń znał z dokumentu `documentId`, `content` i `title`. Model
-- był wobec dokumentu ślepy na jego postać: widział tekst, nie widział kroju,
-- wcięcia, tabeli ani obrazu — więc nie mógł ich ani przeczytać, ani zmienić.
-- Postać przestaje być stanem klienckim i staje się wierszem w bazie; zapis
-- dokumentu przenosi ją razem z treścią, zamiast ją gubić.
--
-- ── Dlaczego bloki i tabele idą jednym zapisem JSON, a style nie ─────────────
-- Drzewo dokumentu (sekcje → bloki → fragmenty o jednolitej postaci znaku,
-- z tabelami w środku) jest strukturą zagnieżdżoną, którą czyta się i zapisuje
-- CAŁĄ: każda czynność na postaci przelicza sąsiedztwo — zmiana wcięcia rusza
-- łamanie, scalenie komórki rusza szerokości. Rozłożenie tego na wiersze
-- kazałoby przy każdej czynności składać drzewo z kilkuset wierszy i pilnować
-- ich kolejności, a żadne zapytanie po pojedynczym fragmencie nie jest do
-- niczego potrzebne. Dlatego `postac_json` niesie drzewo jednym zapisem.
--
-- Styl nazwany i sekcja to co innego i dlatego mają wiersze. Po stylu się
-- PYTA: „ile miejsc go używa”, „które style dziedziczą po tym”, „usuń styl
-- i przenieś jego miejsca użycia”. Sprawdzian odbioru mierzy wprost, że zmiana
-- stylu przestawiła WSZYSTKIE miejsca użycia — a to jest pytanie do bazy, nie
-- do drzewa. Sekcja niesie własne nastawy strony i własne nagłówki, po których
-- pyta podgląd wydruku i wydanie do PDF.
--
-- ── Numer porządkowy postaci ─────────────────────────────────────────────────
-- `wersja_postaci` rośnie z każdym zapisem. Służy dwóm rzeczom: oknu, żeby
-- wiedziało, czy trzyma stan świeży, i dziennikowi czynności, żeby cofnięcie
-- wiedziało, na jakim stanie postaci czynność stała.

CREATE TABLE postac_dokumentu_studio (
    id                       INTEGER PRIMARY KEY AUTOINCREMENT,
    dokument_id              INTEGER NOT NULL UNIQUE
                                     REFERENCES dokument_studio(id) ON DELETE CASCADE,
    postac_json              TEXT    NOT NULL DEFAULT '{}',
    nastawy_strony_json      TEXT,
    wersja_postaci           INTEGER NOT NULL DEFAULT 1,
    utworzono                TEXT    NOT NULL DEFAULT (strftime('%Y-%m-%dT%H:%M:%fZ','now')),
    zaktualizowano           TEXT    NOT NULL DEFAULT (strftime('%Y-%m-%dT%H:%M:%fZ','now'))
);

-- Styl nazwany. Nazwa jest jego jedynym identyfikatorem w obrębie dokumentu —
-- tak samo jak w pakiecie biurowym — dlatego warunek UNIQUE stoi na parze
-- (dokument, nazwa), a nie na osobnym kluczu.
--
-- `styl_nadrzedny` trzymamy nazwą, nie kluczem wiersza: dziedziczenie ma
-- przetrwać przejęcie arkusza stylów z szablonu, gdzie klucze wierszy są inne,
-- a nazwy te same.
CREATE TABLE styl_nazwany_studio (
    id                       INTEGER PRIMARY KEY AUTOINCREMENT,
    dokument_id              INTEGER NOT NULL REFERENCES dokument_studio(id) ON DELETE CASCADE,
    nazwa                    TEXT    NOT NULL,
    nazwa_widoczna           TEXT,
    rodzaj                   TEXT    NOT NULL DEFAULT 'paragraph'
                                     CHECK(rodzaj IN ('paragraph','character','table','list')),
    styl_nadrzedny           TEXT,
    styl_nastepny            TEXT,
    postac_znaku_json        TEXT,
    postac_akapitu_json      TEXT,
    fabryczny                INTEGER NOT NULL DEFAULT 0,
    utworzono                TEXT    NOT NULL DEFAULT (strftime('%Y-%m-%dT%H:%M:%fZ','now')),
    zaktualizowano           TEXT    NOT NULL DEFAULT (strftime('%Y-%m-%dT%H:%M:%fZ','now')),
    UNIQUE(dokument_id, nazwa)
);
CREATE INDEX idx_styl_nazwany_studio_dokument ON styl_nazwany_studio(dokument_id, nazwa);
CREATE INDEX idx_styl_nazwany_studio_nadrzedny ON styl_nazwany_studio(dokument_id, styl_nadrzedny);

-- Sekcja o własnych nastawach. Pismo z załącznikiem w orientacji poziomej ma
-- być JEDNYM dokumentem, nie dwoma — stąd nastawy strony wiszą przy sekcji,
-- a nie tylko przy dokumencie.
--
-- Nagłówki i stopki idą jednym zapisem JSON, bo są wykazem najwyżej trzech
-- pozycji (strony zwykłe, pierwsza strona, strony parzyste) i nikt nie pyta
-- o nie osobno — pyta o nie sekcja, w całości.
CREATE TABLE sekcja_dokumentu_studio (
    id                       INTEGER PRIMARY KEY AUTOINCREMENT,
    identyfikator_zewnetrzny TEXT    NOT NULL UNIQUE,
    dokument_id              INTEGER NOT NULL REFERENCES dokument_studio(id) ON DELETE CASCADE,
    kolejnosc                INTEGER NOT NULL DEFAULT 0,
    tytul                    TEXT,
    zakres_od                INTEGER NOT NULL DEFAULT 0,
    zakres_do                INTEGER NOT NULL DEFAULT 0,
    rozpoczecie              TEXT    NOT NULL DEFAULT 'continuous'
                                     CHECK(rozpoczecie IN ('continuous','newPage','evenPage',
                                                           'oddPage','newColumn')),
    nastawy_strony_json      TEXT,
    naglowki_json            TEXT,
    numeracja_json           TEXT,
    znak_wodny_json          TEXT,
    utworzono                TEXT    NOT NULL DEFAULT (strftime('%Y-%m-%dT%H:%M:%fZ','now')),
    zaktualizowano           TEXT    NOT NULL DEFAULT (strftime('%Y-%m-%dT%H:%M:%fZ','now'))
);
CREATE INDEX idx_sekcja_dokumentu_studio_dokument
    ON sekcja_dokumentu_studio(dokument_id, kolejnosc, id);
