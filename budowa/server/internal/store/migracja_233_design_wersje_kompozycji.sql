-- Migracja 233 — nazwane wersje kompozycji Design Board.
--
-- `design.board.update` zapisuje układ BIEŻĄCY i zastępuje poprzedni: warstwy
-- są usuwane i wstawiane od nowa w jednej transakcji (`ZapiszKompozycje`).
-- Ciągu postaci tablicy nie ma więc dziś skąd wziąć — wersja jest migawką, którą
-- ten zapis omija.
--
-- Warstwy wersji leżą w osobnej tabeli, nie w jednym polu zapisu strukturalnego.
-- Przywrócenie wersji (`design.board.version.restore`) przepisuje je z powrotem
-- do `warstwa_kompozycji_design` kolumna w kolumnę, więc ten sam kształt po obu
-- stronach oszczędza przekład, który mógłby zgubić pole. Odczyt wykazu wersji
-- warstw NIE czyta (kontrakt: `versions` bez `layers`), stąd `liczba_warstw`
-- utrwalona w wierszu wersji — wykaz nie ma ważyć tyle, co cały zapis.
--
-- Identyfikator warstwy wersji jest kopią identyfikatora warstwy żywej, nie
-- odwołaniem do niej: warstwa, z której migawkę zdjęto, bywa już usunięta
-- z kompozycji, a wersja ma ją odtworzyć pod tym samym identyfikatorem.
-- Nie jest więc UNIQUE — ta sama warstwa występuje w każdej wersji, w której ją
-- zapisano.

CREATE TABLE wersja_kompozycji_design (
    id                       INTEGER PRIMARY KEY AUTOINCREMENT,
    identyfikator_zewnetrzny TEXT    NOT NULL UNIQUE,
    kompozycja_id            INTEGER NOT NULL REFERENCES kompozycja_design(id) ON DELETE CASCADE,
    nazwa                    TEXT,
    uzasadnienie             TEXT,
    liczba_warstw            INTEGER NOT NULL DEFAULT 0,
    utworzono                TEXT    NOT NULL DEFAULT (strftime('%Y-%m-%dT%H:%M:%fZ','now'))
);
-- design.board.version.list czyta wersje kompozycji, od najświeższej.
CREATE INDEX idx_wersja_kompozycji_design_kompozycja
    ON wersja_kompozycji_design(kompozycja_id, utworzono DESC, id DESC);

CREATE TABLE warstwa_wersji_kompozycji_design (
    id                    INTEGER PRIMARY KEY AUTOINCREMENT,
    wersja_id             INTEGER NOT NULL REFERENCES wersja_kompozycji_design(id) ON DELETE CASCADE,
    identyfikator_warstwy TEXT    NOT NULL,
    zasob_id              TEXT,
    x                     REAL,
    y                     REAL,
    szerokosc             REAL,
    wysokosc              REAL,
    kolejnosc             INTEGER NOT NULL DEFAULT 0,
    zablokowana           INTEGER NOT NULL DEFAULT 0 CHECK(zablokowana IN (0,1)),
    adnotacja             TEXT
);
CREATE INDEX idx_warstwa_wersji_kompozycji_design_wersja
    ON warstwa_wersji_kompozycji_design(wersja_id, kolejnosc, id);
