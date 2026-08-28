-- Dodaje do tabeli podagent kolumny znacznika uruchomienia rdzenia, oznaki życia i powodu zakończenia oraz indeks częściowy odnajdujący podagentów osieroconych po awarii rdzenia.

-- Zapisuje w kolumnie uruchomienie_rdzenia znacznik rdzenia prowadzącego pracę podagenta, przy czym wartość pusta oznacza brak odnotowanego prowadzenia.
ALTER TABLE podagent ADD COLUMN uruchomienie_rdzenia TEXT NOT NULL DEFAULT '';

-- Dodaje kolumnę oznaka_zycia zapisującą ostatnie dotknięcie wiersza przez rdzeń prowadzący, używaną jako dowód czasu w meldunku o zamkniętym podagencie.
ALTER TABLE podagent ADD COLUMN oznaka_zycia TEXT;

-- Dodaje kolumnę powod_zakonczenia rozróżniającą, czy wiersz w stanie końcowym stanął wskutek błędu modelu, zatrzymania czy osierocenia po awarii rdzenia.
ALTER TABLE podagent ADD COLUMN powod_zakonczenia TEXT
    CHECK(powod_zakonczenia IS NULL OR powod_zakonczenia IN
          ('ukonczony','blad','zatrzymany','osierocony'));

-- Tworzy indeks częściowy na kolumnach stan i uruchomienie_rdzenia dla wierszy niezakończonych, przyspieszający odnajdywanie podagentów osieroconych przy starcie rdzenia.
CREATE INDEX idx_podagent_osieroceni
    ON podagent(stan, uruchomienie_rdzenia)
    WHERE stan IN ('pending','running');
