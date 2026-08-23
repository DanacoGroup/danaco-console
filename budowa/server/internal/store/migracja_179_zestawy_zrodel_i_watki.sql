-- Migracja 179 — zestawy tematyczne źródeł, wątki notatek oraz brakujące
-- własności źródła i notatki (`browser.source.group.*`, `browser.note.thread.*`,
-- `browser.note.update`).
--
-- Kontrakt niesie w `BrowserSource` pole `groupId`, a w `BrowserNote` —
-- `classification`, `threadId` i `pinned`. Migracja 047 zakładała te tabele,
-- zanim rodzina komend obejmowała grupowanie i klasyfikację, więc kolumn tych
-- w nich nie ma. Dopisanie ich tutaj jest jedynym sposobem, żeby okno mogło
-- oddać notatkę tak oznaczoną, jak ją Operator oznaczył — oznaczenie żyjące
-- wyłącznie w kliencie ginie przy przeładowaniu karty.
--
-- Przynależność stoi po stronie źródła i notatki, nie w kolumnie zestawu:
-- źródło należy do jednego zestawu, a notatka do jednego wątku, więc skład
-- zestawu jest zapytaniem po kolumnie, a nie drugą listą do utrzymania.
-- Kontrakt oddaje `sourceIds` i `noteIds` — rdzeń wylicza je z tej samej
-- kolumny, którą zapisuje.
ALTER TABLE zrodlo_przegladania ADD COLUMN grupa TEXT;
ALTER TABLE notatka_przegladania ADD COLUMN klasyfikacja TEXT;
ALTER TABLE notatka_przegladania ADD COLUMN watek TEXT;
ALTER TABLE notatka_przegladania ADD COLUMN przypieta INTEGER NOT NULL DEFAULT 0 CHECK(przypieta IN (0,1));

CREATE TABLE zestaw_zrodel_przegladania (
    id                       INTEGER PRIMARY KEY AUTOINCREMENT,
    identyfikator_zewnetrzny TEXT    NOT NULL UNIQUE,
    okno                     TEXT    NOT NULL,
    nazwa                    TEXT    NOT NULL,
    utworzono                TEXT    NOT NULL DEFAULT (strftime('%Y-%m-%dT%H:%M:%fZ','now')),
    zaktualizowano           TEXT    NOT NULL DEFAULT (strftime('%Y-%m-%dT%H:%M:%fZ','now'))
);
CREATE INDEX idx_zestaw_zrodel_okno ON zestaw_zrodel_przegladania(okno, utworzono DESC, id);

CREATE TABLE watek_notatek_przegladania (
    id                       INTEGER PRIMARY KEY AUTOINCREMENT,
    identyfikator_zewnetrzny TEXT    NOT NULL UNIQUE,
    okno                     TEXT    NOT NULL,
    nazwa                    TEXT    NOT NULL,
    utworzono                TEXT    NOT NULL DEFAULT (strftime('%Y-%m-%dT%H:%M:%fZ','now')),
    zaktualizowano           TEXT    NOT NULL DEFAULT (strftime('%Y-%m-%dT%H:%M:%fZ','now'))
);
CREATE INDEX idx_watek_notatek_okno ON watek_notatek_przegladania(okno, utworzono DESC, id);
