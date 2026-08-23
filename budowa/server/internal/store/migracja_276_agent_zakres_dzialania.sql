-- Migracja 276 — zakres działania eksperta: moduły zastosowania i Subagent
-- Network (Permissions Center, grupy „Dostęp do modułów i zasobów"
-- oraz „MultitaskingAI").
--
-- Dlaczego moduły zastosowania mają własną tabelę, a nie kolumnę z listą.
-- Kontrakt `agent.modules.set` mówi, że LISTA PUSTA znaczy brak ograniczenia,
-- a kod modułu pochodzi z rodziny `module.*`. Lista trzymana napisem
-- rozdzielonym przecinkiem nie da się zapytać „którzy eksperci wolno działają
-- w module X" bez przeszukiwania napisów, a to pytanie zadaje rdzeń przy
-- każdym nałożeniu eksperta na okno. Wiersz na moduł odpowiada na nie indeksem.
--
-- Brak wierszy dla eksperta znaczy dokładnie to samo, co lista pusta
-- w kontrakcie: BRAK OGRANICZENIA. Nie ma tu trzeciego stanu i nie ma potrzeby
-- odróżniania „nie ustawiono" od „ustawiono na puste" — stan wyjściowy
-- platformy jest pełnym dostępem, więc obie odpowiedzi są tą samą odpowiedzią.
--
-- `limit_podagentow` siedzi kolumną na ekspercie, a nie tabelą: wartość jest
-- jedna, liczbowa i czytana przy każdym powołaniu podagentów, więc złączenie
-- byłoby kosztem bez zysku. Zero znaczy Subagent Network wyłączony — jedyny
-- zapis wyłączenia, dokładnie jak mówi opis pola `Agent.subagentLimit`.
-- Wartość wyjściowa 15 jest maksimum technicznym platformy, tym samym, które
-- zna `subagent.spawn`.

CREATE TABLE agent_modul_zastosowania (
    agent_id   INTEGER NOT NULL REFERENCES agent(id) ON DELETE CASCADE,
    -- Kod modułu z rodziny `module.*` kontraktu. Bez więzu obcego: katalog
    -- modułów jest kontraktem, nie tabelą, a rdzeń sprawdza wskazanie przy
    -- zapisie.
    kod_modulu TEXT    NOT NULL,
    zapisano   TEXT    NOT NULL DEFAULT (strftime('%Y-%m-%dT%H:%M:%fZ','now')),
    PRIMARY KEY (agent_id, kod_modulu)
);

-- Pytanie zadawane przy nakładaniu eksperta na okno brzmi „czy TEN ekspert ma
-- wpis na TEN moduł", więc indeks klucza głównego wystarcza. Drugi indeks
-- odpowiada na pytanie odwrotne — którzy eksperci są związani z modułem —
-- zadawane przez wykaz biblioteki zawężony modułem.
CREATE INDEX idx_agent_modul_zastosowania_modul
    ON agent_modul_zastosowania(kod_modulu, agent_id);

ALTER TABLE agent ADD COLUMN limit_podagentow INTEGER NOT NULL DEFAULT 15
    CHECK (limit_podagentow BETWEEN 0 AND 15);
