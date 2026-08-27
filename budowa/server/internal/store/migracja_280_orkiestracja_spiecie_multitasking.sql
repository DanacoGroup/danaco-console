-- Migracja 280 dodaje spięcie kolejek automatyki z kolejkami środowiska
-- wieloagentowego, wskazujące koordynatora i rodzaj kolejki.

CREATE TABLE orkiestracja_spiecie_multitasking (
    automatyka_id INTEGER PRIMARY KEY REFERENCES automatyka(id) ON DELETE CASCADE,
    -- Puste wskazanie roli znaczy spięcie bez koordynatora; kolejki idą
    -- rodzajem wieloagentowym.
    okno_roli_id  INTEGER REFERENCES okno_komunikacji(id) ON DELETE SET NULL,
    zapisano      TEXT    NOT NULL DEFAULT (strftime('%Y-%m-%dT%H:%M:%fZ','now'))
);
