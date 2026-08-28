-- Migracja 260 — wersje definicji automatyki (Workflow Builder, panel „Wersje”).
--
-- Kolumna `automatyka.wersja` mówi wyłącznie, ile razy definicję zapisano.
-- Nie mówi, JAK wyglądała za którymś razem, więc porównanie dwóch wersji
-- i przywrócenie wcześniejszej nie mają z czego powstać. Ten krok daje im
-- źródło: migawkę kroków zapisaną przy każdym zapisie definicji.
--
-- Kroki leżą jako zapis strukturalny, nie jako wiersze. Wersja jest migawką
-- martwą — nikt jej nie edytuje po fakcie, nikt nie pyta o pojedynczy krok
-- wersji siódmej. Rozbicie migawki na tabelę kroków dołożyłoby drugą prawdę
-- o kroku obok `krok_automatyki` i kazałoby ją utrzymywać przy każdej zmianie
-- schematu kroku.
CREATE TABLE wersja_automatyki (
    id             INTEGER PRIMARY KEY AUTOINCREMENT,
    automatyka_id  INTEGER NOT NULL REFERENCES automatyka(id) ON DELETE CASCADE,
    wersja         INTEGER NOT NULL,
    kroki          TEXT    NOT NULL,
    opublikowana   INTEGER NOT NULL DEFAULT 0 CHECK(opublikowana IN (0,1)),
    autor          TEXT,
    utworzono      TEXT    NOT NULL DEFAULT (strftime('%Y-%m-%dT%H:%M:%fZ','now')),
    UNIQUE(automatyka_id, wersja)
);
-- automation.workflow.version.list czyta wersje od najnowszej.
CREATE INDEX idx_wersja_automatyki_czas
    ON wersja_automatyki(automatyka_id, wersja DESC);
