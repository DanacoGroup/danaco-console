-- Migracja 265 — biblioteka szablonów przepływów
-- (`automation.template.save`, `.list`, `.apply`).
--
-- Szablon jest bytem odrębnym od automatyki, z której powstał: automatyka
-- źródłowa bywa później zmieniona albo usunięta, a szablon ma dalej zakładać
-- automatyki o kształcie z chwili zapisu. Dlatego kroki leżą tu jako migawka,
-- a klucza obcego do `automatyka` nie ma.
CREATE TABLE szablon_automatyki (
    id                       INTEGER PRIMARY KEY AUTOINCREMENT,
    identyfikator_zewnetrzny TEXT    NOT NULL UNIQUE,
    nazwa                    TEXT    NOT NULL,
    opis                     TEXT,
    kroki                    TEXT    NOT NULL,
    zaktualizowano           TEXT    NOT NULL DEFAULT (strftime('%Y-%m-%dT%H:%M:%fZ','now'))
);
-- automation.template.list czyta bibliotekę w kolejności nazw.
CREATE INDEX idx_szablon_automatyki_nazwa ON szablon_automatyki(nazwa, id);
