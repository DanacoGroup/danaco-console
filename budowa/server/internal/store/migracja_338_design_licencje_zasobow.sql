-- Migracja 338 — licencje zasobów wciągniętych z katalogów zewnętrznych
-- (`design.stock.import`).
--
-- Licencja stoi w tabeli obok, a nie kolumną w `zasob_design`: dotyczy
-- WYŁĄCZNIE zasobów przyszłych z zewnątrz, a te są mniejszością wśród zasobów
-- okna. Kolumna w tabeli zasobów byłaby pusta przy każdym zasobie
-- wygenerowanym i wniesionym, a każdy odczyt zasobu ciągnąłby ją bez potrzeby.
--
-- Zapis licencji jest częścią wciągnięcia, nie dodatkiem po nim. Zasób
-- z katalogu zewnętrznego bez zapisanej licencji to materiał, o którym za pół
-- roku nikt nie powie, czy wolno go było użyć w kampanii — a pyta o to dopiero
-- pismo od właściciela praw.
--
-- Więz obcy z kasowaniem kaskadowym: licencja bez zasobu nie opisuje niczego.

CREATE TABLE licencja_zasobu_design (
    id                       INTEGER PRIMARY KEY AUTOINCREMENT,
    zasob_id                 INTEGER NOT NULL UNIQUE
                                     REFERENCES zasob_design(id) ON DELETE CASCADE,
    dostawca                 TEXT    NOT NULL,
    identyfikator_u_dostawcy TEXT    NOT NULL,
    licencja                 TEXT,
    autor                    TEXT,
    odsylacz                 TEXT,
    utworzono                TEXT    NOT NULL DEFAULT (strftime('%Y-%m-%dT%H:%M:%fZ','now'))
);

-- Odczyt idzie po zasobie; wyszukanie po dostawcy służy sprawdzeniu, czy ten
-- sam materiał nie został wciągnięty drugi raz.
CREATE INDEX idx_licencja_zasobu_design_dostawca
    ON licencja_zasobu_design(dostawca, identyfikator_u_dostawcy);
