-- Kosz sesji: odwracalne usunięcie zamiast natychmiastowego DELETE.
--
-- `session.delete` stawia znacznik `usunieto_o`, a czyszczenie trwałe wykonuje
-- rdzeń przy starcie, wyłącznie dla sesji leżących w koszu dłużej niż termin
-- (`core/trwalosc_kosza.go`). Utrata danych następuje nadal wyłącznie wskutek
-- tej komendy: kosz i czyszczenie po terminie są dwiema fazami tej samej
-- decyzji, nie nową drogą utraty. Odzysk w oknie terminu daje `session.restore`
-- — przywrócenie czyści znacznik i oddaje sesję historii bieżącej wraz z całym
-- zapisem.
--
-- Sesja w koszu jest niewidoczna: wykaz `Sesje.Lista` (`dane/sesje.go`) pyta
-- o `usunieto_o IS NULL`, więc sesje z kosza nie wracają ani w `session.list`,
-- ani w `session.archive.list`, ani w wynikach szukania. Odnajduje je wyłącznie
-- repozytorium kosza (`dane/sesje_kosz.go`).
--
-- Kosz nie jest wartością słownika `sesja.stan`: stan opisuje położenie sesji
-- w historii, nie jej byt. Osobna kolumna pozwala sesji wrócić z kosza dokładnie
-- w tym stanie, w którym ją usunięto.
--
-- Format znacznika — ISO 8601 UTC, jak `utworzono` i `zaktualizowano` tej samej
-- tabeli; NULL znaczy „żywa".

ALTER TABLE sesja ADD COLUMN usunieto_o TEXT;

-- Czyszczenie przy starcie pyta wyłącznie o wiersze ze znacznikiem.
CREATE INDEX idx_sesja_kosz ON sesja (usunieto_o) WHERE usunieto_o IS NOT NULL;
