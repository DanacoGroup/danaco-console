-- Migracja 180 — moduł Library: cykl życia zasobu, jego miejsce w strukturze
-- repozytorium oraz opis w schemacie Dublin Core wraz z polami niestandardowymi.
--
-- Trzy braki naraz, bo wszystkie trzy dotyczą jednego bytu — zasobu:
--
--   1. `stan` rozdziela wykaz czynny od archiwum. Kosz repozytorium
--      (`library.file.archive` / `library.file.restore`) jest przeniesieniem
--      między stanami, nie usunięciem wiersza: zasób zarchiwizowany zachowuje
--      wersje, etykiety i kolekcje, więc przywrócenie oddaje go w całości.
--      Usunięcie trwałe (`library.file.delete`) zdejmuje wiersz i wtedy dopiero
--      kaskada zabiera wersje.
--   2. `sciezka_repozytorium` jest drogą WEWNĄTRZ biblioteki
--      (`LibraryFile.path`), rozłączną z kolumną `sciezka`, która niesie
--      ścieżkę źródłową z maszyny Operatora i z rdzenia nie wychodzi
--      (`dane/library.go`). Bez osobnej kolumny `library.file.move` nie miałby
--      dokąd przenieść zasobu, a wypełnienie pola kontraktu ścieżką źródłową
--      wyniosłoby do klienta układ cudzego dysku.
--   3. Opis Dublin Core mieszka w tabeli towarzyszącej, nie w kolumnach
--      `plik_biblioteki`: piętnaście pól opisowych obciążałoby każdy odczyt
--      wykazu, a wykaz opisu nie pokazuje. Jeden wiersz opisu na jeden zasób —
--      klucz główny jest kluczem obcym.
--
-- Pola niestandardowe stoją dwutorowo, bo są dwiema różnymi rzeczami:
-- DEFINICJA pola należy do repozytorium (`pole_schematu_biblioteki`,
-- `library.schema.set`), a WARTOŚĆ pola do zasobu (kolumna
-- `pola_niestandardowe` opisu, mapa kod→wartość w zapisie JSON). Rozdział ten
-- ma skutek wprost w kontrakcie: zdjęcie definicji nie kasuje wartości
-- zapisanych przy zasobach i wartości wracają, gdy pole zostanie założone
-- ponownie.

-- ── Cykl życia zasobu i jego miejsce w strukturze ──────────────────────────
ALTER TABLE plik_biblioteki
    ADD COLUMN stan TEXT NOT NULL DEFAULT 'aktywny'
        CHECK(stan IN ('aktywny','zarchiwizowany'));

ALTER TABLE plik_biblioteki ADD COLUMN sciezka_repozytorium TEXT;

-- Wykaz domyślny pokazuje zasoby czynne, więc stan wchodzi do indeksu przed
-- porządkiem czasu.
CREATE INDEX idx_plik_biblioteki_stan ON plik_biblioteki(stan, zaktualizowano DESC);
CREATE INDEX idx_plik_biblioteki_sciezka_repozytorium
    ON plik_biblioteki(sciezka_repozytorium);
-- Rozpoznanie duplikatów dokładnych idzie po sumie kontrolnej — bez indeksu
-- byłoby przejściem po całym repozytorium przy każdym skanowaniu.
CREATE INDEX idx_plik_biblioteki_suma ON plik_biblioteki(suma_kontrolna);

-- ── Opis zasobu — Dublin Core i pola niestandardowe ────────────────────────
CREATE TABLE opis_zasobu_biblioteki (
    plik_id              INTEGER PRIMARY KEY
                         REFERENCES plik_biblioteki(id) ON DELETE CASCADE,
    tytul                TEXT,
    tworca               TEXT,
    temat                TEXT,
    opis                 TEXT,
    wydawca              TEXT,
    wspoltworca          TEXT,
    data_zasobu          TEXT,
    rodzaj               TEXT,
    format               TEXT,
    identyfikator        TEXT,
    zrodlo               TEXT,
    jezyk                TEXT,
    powiazanie           TEXT,
    zakres               TEXT,
    prawa                TEXT,
    -- Mapa kod pola → wartość w zapisie JSON. Kolumna na pole byłaby schematem
    -- zmienianym migracją przy każdym polu założonym przez Operatora.
    pola_niestandardowe  TEXT,
    zaktualizowano       TEXT NOT NULL DEFAULT (strftime('%Y-%m-%dT%H:%M:%fZ','now'))
);

-- ── Definicja pola niestandardowego schematu metadanych ────────────────────
CREATE TABLE pole_schematu_biblioteki (
    id             INTEGER PRIMARY KEY AUTOINCREMENT,
    kod            TEXT    NOT NULL UNIQUE,
    etykieta       TEXT    NOT NULL,
    rodzaj         TEXT    NOT NULL
                   CHECK(rodzaj IN ('tekst','liczba','data','logiczna','lista')),
    wymagane       INTEGER NOT NULL DEFAULT 0 CHECK(wymagane IN (0,1)),
    -- Zawężenia stosowalności pola: rodzaj treści i kolekcja. Puste znaczy
    -- „wszystkie" — brak wartości jest tu wartością domyślną, nie brakiem.
    mime_type      TEXT,
    kolekcja_kod   TEXT,
    -- Słownik dopuszczalnych wartości pola rodzaju `lista`, zapis JSON.
    opcje          TEXT,
    utworzono      TEXT    NOT NULL DEFAULT (strftime('%Y-%m-%dT%H:%M:%fZ','now'))
);
CREATE INDEX idx_pole_schematu_biblioteki_zakres
    ON pole_schematu_biblioteki(mime_type, kolekcja_kod);
