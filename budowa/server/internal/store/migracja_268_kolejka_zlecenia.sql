-- Migracja 268 wprowadza zlecenia kolejki jako byt odrębny od pozycji kolejki: niesie ładunek
-- strukturalny, priorytet, termin wykonania, klucz idempotencji i warunek przetworzenia.
CREATE TABLE zlecenie_kolejki (
    id                       INTEGER PRIMARY KEY AUTOINCREMENT,
    identyfikator_zewnetrzny TEXT    NOT NULL UNIQUE,
    kolejka_id               INTEGER NOT NULL REFERENCES kolejka(id) ON DELETE CASCADE,
    stan                     TEXT    NOT NULL DEFAULT 'oczekuje'
                                     CHECK(stan IN ('oczekuje','odlozone','przetwarzane',
                                                    'zakonczone','bledne','martwe','zdjete')),
    ladunek                  TEXT,
    priorytet                INTEGER NOT NULL DEFAULT 0,
    proby                    INTEGER NOT NULL DEFAULT 0,
    warunek                  TEXT,
    klucz_idempotencji       TEXT,
    przebieg_id              INTEGER REFERENCES przebieg_automatyki(id) ON DELETE SET NULL,
    ekspert_docelowy         TEXT,
    zlecenie_zrodlowe_id     INTEGER REFERENCES zlecenie_kolejki(id) ON DELETE SET NULL,
    termin                   TEXT,
    utworzono                TEXT    NOT NULL DEFAULT (strftime('%Y-%m-%dT%H:%M:%fZ','now')),
    zaktualizowano           TEXT    NOT NULL DEFAULT (strftime('%Y-%m-%dT%H:%M:%fZ','now'))
);
-- Indeks porządkuje zlecenia jednej kolejki po priorytecie i kolejności zapisu, w porządku,
-- w jakim zlecenia trafiają do przetwarzania.
CREATE INDEX idx_zlecenie_kolejki_porzadek
    ON zlecenie_kolejki(kolejka_id, priorytet, id);
-- Indeks porządkuje zlecenia po stanie i czasie ostatniej zmiany, co pozwala odczytać zlecenia
-- trwale nieudane wszystkich kolejek naraz.
CREATE INDEX idx_zlecenie_kolejki_stan ON zlecenie_kolejki(stan, zaktualizowano DESC);
-- Indeks porządkuje zlecenia po czasie utworzenia, co pozwala policzyć zlecenia oczekujące
-- w wybranych odcinkach czasu.
CREATE INDEX idx_zlecenie_kolejki_czas ON zlecenie_kolejki(utworzono);
CREATE UNIQUE INDEX idx_zlecenie_kolejki_idempotencja
    ON zlecenie_kolejki(kolejka_id, klucz_idempotencji)
    WHERE klucz_idempotencji IS NOT NULL;
