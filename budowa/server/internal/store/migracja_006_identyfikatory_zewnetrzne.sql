-- Migracja dokłada kolumnę identyfikatora zewnętrznego do sesji, okna komunikacji i
-- wiadomości, wiążącą wiersz z pamięcią rdzenia.

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
