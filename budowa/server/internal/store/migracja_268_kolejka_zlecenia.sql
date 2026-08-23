-- Migracja 268 — zlecenia kolejki widziane przez Queue Managera
-- (rodzina `queue.item.*`, `queue.dead.list`, `queue.depth.get`).
--
-- Tabela `pozycja_kolejki` (migracja 003) zostaje nietknięta. Opisuje ona etap
-- pętli koordynator–wykonawca: tytuł, treść zlecenia, werdykt weryfikacji,
-- licznik obiegów. Kontraktowy `QueueItem` jest czym innym — niesie ładunek
-- strukturalny, priorytet, termin wykonania, klucz idempotencji i warunek
-- przetworzenia. Wtłoczenie jednego w drugie kazałoby kolumnie `tytul` nieść
-- ładunek, a `werdykt_weryfikacji` — stan zlecenia o innym słowniku.
--
-- Stan idzie słownikiem bazy (kolumna `baza` wyliczenia `QueueItemStatus`),
-- tak samo jak stan kolejki i stan pozycji.
--
-- Klucz idempotencji jest unikatowy W OBRĘBIE KOLEJKI, nie globalnie: ten sam
-- klucz w dwóch kolejkach opisuje dwa różne zlecenia dwóch różnych torów.
-- Zlecenie martwe zachowuje kolejkę źródłową — `queue.dead.list` bez wskazania
-- kolejki oddaje zadania martwe wszystkich kolejek, więc rozdzielenie ich na
-- osobną tabelę odebrałoby im pochodzenie.
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
-- queue.item.list czyta zlecenia kolejki w porządku przetwarzania.
CREATE INDEX idx_zlecenie_kolejki_porzadek
    ON zlecenie_kolejki(kolejka_id, priorytet, id);
-- queue.dead.list czyta zlecenia trwale nieudane wszystkich kolejek naraz.
CREATE INDEX idx_zlecenie_kolejki_stan ON zlecenie_kolejki(stan, zaktualizowano DESC);
-- queue.depth.get liczy zlecenia oczekujące w odcinkach czasu.
CREATE INDEX idx_zlecenie_kolejki_czas ON zlecenie_kolejki(utworzono);
CREATE UNIQUE INDEX idx_zlecenie_kolejki_idempotencja
    ON zlecenie_kolejki(kolejka_id, klucz_idempotencji)
    WHERE klucz_idempotencji IS NOT NULL;
