-- Migracja 142 — kolekcje zapytań okna API Client modułu Developer.
--
-- Kolekcja jest zestawem zapytań HTTP wraz ze środowiskami, w których je się
-- uruchamia. Zapytanie wpisane raz i zgubione po zamknięciu okna nie jest
-- klientem API, tylko polem tekstowym — dlatego kolekcja ma miejsce w bazie.
--
-- `zapytania` i `srodowiska` trzymamy jako tekst JSON, a nie tabelami
-- podrzędnymi. Kontrakt niesie je jako surowy JSON (`json.RawMessage`) i nikt
-- nie pyta o pojedyncze zapytanie kolekcji osobnym żądaniem: kolekcja przychodzi
-- i odchodzi w całości. Rozbicie jej na wiersze byłoby rozkładaniem i składaniem
-- tej samej struktury po obu stronach bez jednego odbiorcy tej pracy.
--
-- Import OpenAPI zapisuje się tą samą tabelą — kolekcja z importu nie różni się
-- niczym od kolekcji ułożonej ręcznie poza tym, skąd wzięła zapytania.
CREATE TABLE developer_kolekcja_api (
    id         INTEGER PRIMARY KEY AUTOINCREMENT,
    kod        TEXT    NOT NULL UNIQUE,
    okno_kod   TEXT    NOT NULL,
    nazwa      TEXT    NOT NULL,
    zapytania  TEXT    NOT NULL DEFAULT '[]',
    srodowiska TEXT,
    zmieniono  TEXT    NOT NULL DEFAULT (strftime('%Y-%m-%dT%H:%M:%fZ','now'))
);
CREATE INDEX idx_developer_kolekcja_okno ON developer_kolekcja_api(okno_kod, zmieniono DESC);
