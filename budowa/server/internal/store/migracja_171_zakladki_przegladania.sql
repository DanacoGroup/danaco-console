-- Migracja 171 — zakładki okna przeglądarki (`browser.bookmark.*`).
--
-- Zakładka nie jest źródłem przeglądania. Źródło (`zrodlo_przegladania`,
-- migracja 047) narasta samo w miarę odwiedzania stron i opisuje przebieg
-- sesji; zakładkę Operator zakłada świadomie i ma ona trwać dłużej niż sesja,
-- razem z folderem, etykietami i własną notatką. Wspólna tabela kazałaby
-- odsiewać jedno od drugiego kolumną rodzaju przy każdym odczycie obu wykazów.
--
-- Etykiety idą kolumną JSON, bo kontrakt niesie je jako `tags: string[]`
-- wewnątrz `BrowserBookmark` — nie jest to byt samodzielny, który ktokolwiek
-- wyszukuje niezależnie od zakładki. Osobna tabela wiązań dawałaby złączenie
-- przy każdym odczycie i ani jednego nowego pytania, na które umiałaby
-- odpowiedzieć.
CREATE TABLE zakladka_przegladania (
    id                       INTEGER PRIMARY KEY AUTOINCREMENT,
    identyfikator_zewnetrzny TEXT    NOT NULL UNIQUE,
    okno                     TEXT    NOT NULL,
    url                      TEXT    NOT NULL,
    tytul                    TEXT,
    folder                   TEXT,
    etykiety_json            TEXT,
    notatka                  TEXT,
    utworzono                TEXT    NOT NULL DEFAULT (strftime('%Y-%m-%dT%H:%M:%fZ','now'))
);
CREATE INDEX idx_zakladka_przegladania_okno ON zakladka_przegladania(okno, utworzono DESC, id);
