-- Migracja 006 — identyfikator zewnętrzny sesji, okna komunikacji i wiadomości.
--
-- Rdzeń posługuje się identyfikatorem tekstowym nadanym w pamięci przez pakiet
-- `session` (`sesja-…`, `okno-…`, `wiad-…`) i ten identyfikator wychodzi
-- kontraktem do klienta. Wiersz bazy ma własny klucz główny INTEGER. Bez trwałego
-- odwzorowania jednego na drugie po restarcie rdzenia nie da się połączyć okna
-- wskazanego przez klienta z jego historią w tabeli `wiadomosc`.
--
-- Kolumna jest dodatkiem: NULL oznacza wiersz założony wprost w bazie, bez
-- odpowiednika w pamięci rdzenia. Indeks jest częściowy, więc wiersze bez
-- identyfikatora zewnętrznego nie kolidują ze sobą.

ALTER TABLE sesja            ADD COLUMN identyfikator_zewnetrzny TEXT;
ALTER TABLE okno_komunikacji ADD COLUMN identyfikator_zewnetrzny TEXT;
ALTER TABLE wiadomosc        ADD COLUMN identyfikator_zewnetrzny TEXT;

CREATE UNIQUE INDEX idx_sesja_identyfikator_zewnetrzny
    ON sesja(identyfikator_zewnetrzny)
    WHERE identyfikator_zewnetrzny IS NOT NULL;

CREATE UNIQUE INDEX idx_okno_identyfikator_zewnetrzny
    ON okno_komunikacji(identyfikator_zewnetrzny)
    WHERE identyfikator_zewnetrzny IS NOT NULL;

CREATE UNIQUE INDEX idx_wiadomosc_identyfikator_zewnetrzny
    ON wiadomosc(identyfikator_zewnetrzny)
    WHERE identyfikator_zewnetrzny IS NOT NULL;
