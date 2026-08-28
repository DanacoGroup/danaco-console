-- Migracja 234 — adnotacje kompozycji Design Board wraz z wątkami.
--
-- Pole `warstwa_kompozycji_design.adnotacja` niesie JEDNO zdanie bez autora,
-- bez czasu i bez wątku, a przy każdym `design.board.update` jedzie razem
-- z całym układem warstw i wraca przepisane od nowa. Uwaga zostawiona przez
-- jedną osobę znikałaby więc przy pierwszym przesunięciu warstwy przez drugą.
-- Adnotacja ma własny wiersz, własny czas i własnego autora — i przeżywa każdy
-- zapis układu.
--
-- Wątek powstaje przez `nadrzedna_id` wskazujące inną adnotację tej samej
-- kompozycji. Wskazanie jest identyfikatorem zewnętrznym, nie więzem obcym:
-- odpowiedź w wątku zakłada się zaraz po adnotacji nadrzędnej i oba wiersze
-- bywają wstawiane w jednym przebiegu okna, a więz obcy narzuciłby porządek,
-- którego kontrakt nie zna.
--
-- Adnotacja przypięta do kompozycji, a nie do warstwy, ma `warstwa_id` puste —
-- kontrakt czyni `layerId` niewymaganym właśnie po to.

CREATE TABLE adnotacja_kompozycji_design (
    id                       INTEGER PRIMARY KEY AUTOINCREMENT,
    identyfikator_zewnetrzny TEXT    NOT NULL UNIQUE,
    kompozycja_id            INTEGER NOT NULL REFERENCES kompozycja_design(id) ON DELETE CASCADE,
    warstwa_id               TEXT,
    nadrzedna_id             TEXT,
    autor                    TEXT,
    tresc                    TEXT    NOT NULL,
    zamknieta                INTEGER NOT NULL DEFAULT 0 CHECK(zamknieta IN (0,1)),
    utworzono                TEXT    NOT NULL DEFAULT (strftime('%Y-%m-%dT%H:%M:%fZ','now'))
);
-- design.annotation.list czyta adnotacje kompozycji, od najstarszej, bo wątek
-- czyta się w kolejności powstawania; zawężenie do niezamkniętych idzie tą samą
-- ścieżką.
CREATE INDEX idx_adnotacja_kompozycji_design_kompozycja
    ON adnotacja_kompozycji_design(kompozycja_id, zamknieta, id);
