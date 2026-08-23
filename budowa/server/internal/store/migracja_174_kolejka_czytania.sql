-- Migracja 174 — kolejka czytania (`browser.readlist.*`).
--
-- Kolejka czytania jest odłożeniem strony na później wraz z przypomnieniem
-- (opracowanie modułu, rozdz. 2.2). Pozycja przeczytana nie znika z tabeli:
-- `browser.readlist.remove` z polem `markRead` oddaje pozycję, a nie sam fakt
-- usunięcia — kolejka ma pamiętać, co już przeczytano, żeby ta sama strona nie
-- wracała jako nowa.
--
-- Przypomnienie jest chwilą, nie flagą: kontrakt niesie `remindAt` jako czas
-- w milisekundach epoki, więc kolumna trzyma znacznik ISO 8601 tej chwili,
-- a jego brak znaczy „bez przypomnienia".
CREATE TABLE pozycja_czytania_przegladania (
    id                       INTEGER PRIMARY KEY AUTOINCREMENT,
    identyfikator_zewnetrzny TEXT    NOT NULL UNIQUE,
    okno                     TEXT    NOT NULL,
    url                      TEXT    NOT NULL,
    tytul                    TEXT,
    notatka                  TEXT,
    przeczytana              INTEGER NOT NULL DEFAULT 0 CHECK(przeczytana IN (0,1)),
    przypomnienie            TEXT,
    utworzono                TEXT    NOT NULL DEFAULT (strftime('%Y-%m-%dT%H:%M:%fZ','now'))
);
CREATE INDEX idx_pozycja_czytania_okno ON pozycja_czytania_przegladania(okno, utworzono DESC, id);
