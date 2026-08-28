-- Migracja 265 wprowadza bibliotekę szablonów przepływów: szablon jest bytem odrębnym od
-- automatyki źródłowej, więc kroki leżą jako migawka, bez klucza obcego do tabeli automatyk.
CREATE TABLE szablon_automatyki (
    id                       INTEGER PRIMARY KEY AUTOINCREMENT,
    identyfikator_zewnetrzny TEXT    NOT NULL UNIQUE,
    nazwa                    TEXT    NOT NULL,
    opis                     TEXT,
    kroki                    TEXT    NOT NULL,
    zaktualizowano           TEXT    NOT NULL DEFAULT (strftime('%Y-%m-%dT%H:%M:%fZ','now'))
);
-- Indeks porządkuje szablony przepływów alfabetycznie po nazwie, w kolejności, w jakiej
-- biblioteka szablonów jest przedstawiana przy przeglądaniu.
CREATE INDEX idx_szablon_automatyki_nazwa ON szablon_automatyki(nazwa, id);
