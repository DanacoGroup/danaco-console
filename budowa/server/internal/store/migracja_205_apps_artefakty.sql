-- Migracja 205 — moduł Apps: artefakty budowania.
--
-- Artefakt jest plikiem na dysku, nie wpisem o pliku. Wiersz powstaje wtedy,
-- gdy silnik wykonania wdrożenia spakuje przestrzeń roboczą okna do archiwum
-- w magazynie treści rdzenia — kolumna `sciezka` trzyma odwołanie względne
-- magazynu (ten sam wzorzec co `zasob_designu.sciezka`), a `rozmiar`
-- i `suma_kontrolna` opisują bajty, które tam naprawdę leżą. Wiersz bez pliku
-- byłby dokładnie tym wzorcem szkody, który w tym produkcie już wystąpił:
-- wykazem zasobów, za którymi nie ma ani jednego bajtu.
--
-- `wdrozenie_kod` jest kodem zewnętrznym, nie więzem obcym — artefakt przeżywa
-- przebieg, z którego powstał, bo to on idzie potem do pakowania
-- (`apps.package.build` bierze „artefakt ostatniego wdrożenia udanego").

CREATE TABLE artefakt_apps (
    id                       INTEGER PRIMARY KEY AUTOINCREMENT,
    identyfikator_zewnetrzny TEXT    NOT NULL UNIQUE,
    okno                     TEXT    NOT NULL,
    wdrozenie_kod            TEXT,
    -- Wartości kontraktu (AppArtifactKind) wprost, bez tłumaczenia.
    rodzaj                   TEXT    NOT NULL
                                     CHECK(rodzaj IN ('bundle','binary','container','archive','sourceMap')),
    -- Odwołanie względne magazynu treści rdzenia.
    sciezka                  TEXT    NOT NULL,
    rozmiar                  INTEGER,
    suma_kontrolna           TEXT,
    utworzono                TEXT    NOT NULL DEFAULT (strftime('%Y-%m-%dT%H:%M:%fZ','now'))
);
CREATE INDEX idx_artefakt_apps_okno ON artefakt_apps(okno, id DESC);
CREATE INDEX idx_artefakt_apps_wdrozenie ON artefakt_apps(wdrozenie_kod, id DESC);
