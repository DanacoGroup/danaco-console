-- Migracja 079 — wybudzenie biegu orkiestracji: jeden mechanizm dla obu kombajnów.
--
-- Wybudzenie to nie pauza. Pauza czeka na czas albo na Operatora — wiadomo,
-- kiedy ruszy. Wybudzenie czeka na świat: na reakcję, która może przyjść za
-- godzinę, za trzy dni albo nigdy. Bieg zawieszony musi przeżyć restart rdzenia:
-- bieg trzymany w pamięci procesu ginie z procesem i nie jest wybudzeniem, tylko
-- czekaniem.
--
-- Tabela `oczekiwanie_biegu` już istnieje — z terminem, polityką po terminie
-- i powodem wybudzenia — a `core/handlers_automatyka_petla.go` prowadzi silnik
-- z zegarem. Ta migracja nie zakłada drugiego kompletu, tylko rozluźnia nośnik
-- istniejącego.
--
-- `przebieg_id` nie jest NOT NULL, bo bieg orkiestracji nie ma przebiegu
-- automatyki ani jego `automatyka_id`. Oczekiwanie wisi albo na przebiegu
-- automatyki, albo na biegu orkiestracji — nigdy na obu i nigdy na żadnym.
--
-- Wyzwalacz uruchamia automatykę od początku, a wybudzenie wznawia konkretny
-- bieg w konkretnym miejscu. To dwa różne byty i dlatego `oczekiwanie_biegu`
-- nie jest drugim wyzwalaczem. Kolumna `sygnal` została bez CHECK-a właśnie
-- dlatego, że świata nie da się zamknąć w słowniku: kod sygnału ustala ten, kto
-- sygnał wyśle.

-- ── Przebudowa `oczekiwanie_biegu` — dwa nośniki ────────────────────────────
--
-- SQLite nie zdejmuje NOT NULL, więc tabela idzie przez przebudowę: nowa tabela,
-- przepisanie wierszy, podmiana nazwy. Na `oczekiwanie_biegu` nie wskazuje żaden
-- klucz obcy, więc przebudowa obejmuje jedną tabelę.
--
-- Więz „dokładnie jeden nośnik" liczy wyrażenia logiczne jako 0/1 — suma równa 1
-- wyraża „dokładnie jeden" bez wyzwalacza. Oczekiwanie bez nośnika byłoby
-- czekaniem niczyim, z dwoma — czekaniem dwóch biegów na jeden sygnał.

CREATE TABLE oczekiwanie_biegu_nowa (
    id               INTEGER PRIMARY KEY AUTOINCREMENT,
    -- Nośnik pierwszy: przebieg automatyki.
    przebieg_id      INTEGER REFERENCES przebieg_automatyki(id) ON DELETE CASCADE,
    -- Nośnik drugi: bieg orkiestracji.
    bieg_id          INTEGER REFERENCES bieg_orkiestracji(id) ON DELETE CASCADE,
    -- Miejsce zatrzymania. Dla automatyki to identyfikator zewnętrzny kroku;
    -- dla orkiestracji — numer miejsca obsady zapisany tekstem. Jedna kolumna,
    -- bo oba są tym samym faktem: „gdzie bieg stanął".
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

-- ── Indeksy odtworzone i rozszerzone o drugi nośnik ─────────────────────────
--
-- Jedno czynne oczekiwanie na nośnik. Bieg stoi w jednym miejscu, więc nie może
-- czekać na dwie rzeczy naraz. Dwa indeksy częściowe zamiast jednego, bo NULL
-- nie równa się NULL i wspólny indeks przepuściłby dublet dla drugiego nośnika.

CREATE UNIQUE INDEX idx_oczekiwanie_biegu_czynne
    ON oczekiwanie_biegu(przebieg_id)
    WHERE wybudzono IS NULL AND przebieg_id IS NOT NULL;

CREATE UNIQUE INDEX idx_oczekiwanie_biegu_czynne_orkiestracja
    ON oczekiwanie_biegu(bieg_id)
    WHERE wybudzono IS NULL AND bieg_id IS NOT NULL;

-- Doręczenie sygnału: „które biegi czekają na ten sygnał". Wspólny dla obu
-- nośników, bo sygnał ze świata nie wie, czyj bieg budzi.
CREATE INDEX idx_oczekiwanie_biegu_sygnal
    ON oczekiwanie_biegu(sygnal)
    WHERE wybudzono IS NULL;

-- Budzik terminów. Oczekiwanie bez terminu zostaje poza indeksem — czekanie bez
-- zegara jest wyborem Operatora, nie zaległością.
CREATE INDEX idx_oczekiwanie_biegu_termin
    ON oczekiwanie_biegu(termin)
    WHERE wybudzono IS NULL AND termin IS NOT NULL;
