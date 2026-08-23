-- Migracja 216 — zadania projektu modułu Workspace.
--
-- Zadanie jest jednym bytem trzech widoków huba planowania: pozycją listy,
-- kartą tablicy i słupkiem osi czasu. Trzech tabel nie ma, bo trzy tabele
-- rozjechałyby się przy pierwszej zmianie stanu wykonanej z innego widoku.
--
-- Etykiety i lista kontrolna leżą w kolumnach tego samego wiersza: etykiety
-- rozdzielone znakiem nowego wiersza, lista kontrolna zapisem JSON. Osobne
-- tabele wiążące dawałyby tu wyłącznie koszt złączeń — żadne okno nie pyta
-- o etykietę bez zadania ani o krok bez zadania.
--
-- Klucz porządkowy karty jest NAPISEM, nie liczbą: wstawienie karty między dwie
-- sąsiednie ma dopisać klucz pośredni, a nie przepisać całą kolumnę.

CREATE TABLE zadanie_projektu (
    id                       INTEGER PRIMARY KEY AUTOINCREMENT,
    identyfikator_zewnetrzny TEXT    NOT NULL UNIQUE,
    projekt_id               INTEGER NOT NULL REFERENCES projekt(id) ON DELETE CASCADE,
    tytul                    TEXT    NOT NULL,
    opis                     TEXT    NOT NULL DEFAULT '',
    stan                     TEXT    NOT NULL DEFAULT 'todo'
                                     CHECK(stan IN ('todo','inProgress','inReview',
                                                    'done','blocked','cancelled')),
    waga                     TEXT    NOT NULL DEFAULT 'normal'
                                     CHECK(waga IN ('low','normal','high','urgent')),
    -- Rodzaj wykonawcy pusty znaczy zadanie jeszcze nieprzypisane.
    rodzaj_wykonawcy         TEXT    NOT NULL DEFAULT ''
                                     CHECK(rodzaj_wykonawcy IN ('','operator','agent')),
    wykonawca                TEXT    NOT NULL DEFAULT '',
    zadanie_nadrzedne        TEXT    NOT NULL DEFAULT '',
    kolumna_tablicy          TEXT    NOT NULL DEFAULT '',
    klucz_porzadkowy         TEXT    NOT NULL DEFAULT '',
    -- Granice czasu w milisekundach epoki; zero znaczy brak granicy.
    poczatek_ms              INTEGER NOT NULL DEFAULT 0,
    termin_ms                INTEGER NOT NULL DEFAULT 0,
    ukonczono_ms             INTEGER NOT NULL DEFAULT 0,
    szacunek_minut           INTEGER NOT NULL DEFAULT 0,
    spedzono_minut           INTEGER NOT NULL DEFAULT 0,
    postep_procent           INTEGER NOT NULL DEFAULT 0,
    kamien_milowy            INTEGER NOT NULL DEFAULT 0,
    regula_powtarzalnosci    TEXT    NOT NULL DEFAULT '',
    -- Etykiety rozdzielone znakiem nowego wiersza.
    etykiety                 TEXT    NOT NULL DEFAULT '',
    -- Lista kontrolna zapisem JSON (tablica pozycji).
    lista_kontrolna          TEXT    NOT NULL DEFAULT '',
    utworzono                TEXT    NOT NULL DEFAULT (strftime('%Y-%m-%dT%H:%M:%fZ','now')),
    zaktualizowano           TEXT    NOT NULL DEFAULT (strftime('%Y-%m-%dT%H:%M:%fZ','now'))
);
CREATE INDEX idx_zadanie_projektu_projekt ON zadanie_projektu(projekt_id, stan, id DESC);
CREATE INDEX idx_zadanie_projektu_termin ON zadanie_projektu(projekt_id, termin_ms);
CREATE INDEX idx_zadanie_projektu_kolumna ON zadanie_projektu(projekt_id, kolumna_tablicy, klucz_porzadkowy);
