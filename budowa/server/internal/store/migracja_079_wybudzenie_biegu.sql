-- Migracja zakłada wybudzenie biegu orkiestracji jako jeden mechanizm dla obu kombajnów,
-- rozluźniając nośnik istniejącej tabeli oczekiwania.

-- Tabela przechodzi przebudowę, bo baza nie zdejmuje warunku obowiązkowości kolumny.

CREATE TABLE oczekiwanie_biegu_nowa (
    id               INTEGER PRIMARY KEY AUTOINCREMENT,
    -- Nośnik pierwszy: przebieg automatyki.
    przebieg_id      INTEGER REFERENCES przebieg_automatyki(id) ON DELETE CASCADE,
    -- Nośnik drugi: bieg orkiestracji.
    bieg_id          INTEGER REFERENCES bieg_orkiestracji(id) ON DELETE CASCADE,
    -- Miejsce zatrzymania biegu, jedna kolumna: oba nośniki niosą ten sam fakt, gdzie bieg stanął.
    krok_zewnetrzny  TEXT    NOT NULL,
    sygnal           TEXT    NOT NULL,
    termin           TEXT,
    po_terminie      TEXT    NOT NULL DEFAULT 'wznow'
                             CHECK(po_terminie IN ('wznow','ponow','przerwij')),
    zalozono         TEXT    NOT NULL DEFAULT (strftime('%Y-%m-%dT%H:%M:%fZ','now')),
    wybudzono        TEXT,
    powod_wybudzenia TEXT    CHECK(powod_wybudzenia IS NULL
                                   OR powod_wybudzenia IN ('sygnal','termin','operator')),
    tresc_sygnalu    TEXT,
    CHECK (wybudzono IS NULL OR powod_wybudzenia IS NOT NULL),
    CHECK ((przebieg_id IS NOT NULL) + (bieg_id IS NOT NULL) = 1)
);

INSERT INTO oczekiwanie_biegu_nowa
    (id, przebieg_id, krok_zewnetrzny, sygnal, termin, po_terminie,
     zalozono, wybudzono, powod_wybudzenia, tresc_sygnalu)
SELECT
     id, przebieg_id, krok_zewnetrzny, sygnal, termin, po_terminie,
     zalozono, wybudzono, powod_wybudzenia, tresc_sygnalu
  FROM oczekiwanie_biegu;

DROP TABLE oczekiwanie_biegu;

ALTER TABLE oczekiwanie_biegu_nowa RENAME TO oczekiwanie_biegu;

-- Jedno czynne oczekiwanie na nośnik.

CREATE UNIQUE INDEX idx_oczekiwanie_biegu_czynne
    ON oczekiwanie_biegu(przebieg_id)
    WHERE wybudzono IS NULL AND przebieg_id IS NOT NULL;

CREATE UNIQUE INDEX idx_oczekiwanie_biegu_czynne_orkiestracja
    ON oczekiwanie_biegu(bieg_id)
    WHERE wybudzono IS NULL AND bieg_id IS NOT NULL;

-- Doręczenie sygnału obu nośnikom naraz, niezależnie od tego, który z nich sygnał ten wybudzi ostatecznie.
CREATE INDEX idx_oczekiwanie_biegu_sygnal
    ON oczekiwanie_biegu(sygnal)
    WHERE wybudzono IS NULL;

-- Budzik terminów; oczekiwanie bez wyznaczonego terminu zostaje poza tym indeksem czasu, bo czekanie bez zegara jest wyborem.
CREATE INDEX idx_oczekiwanie_biegu_termin
    ON oczekiwanie_biegu(termin)
    WHERE wybudzono IS NULL AND termin IS NOT NULL;
