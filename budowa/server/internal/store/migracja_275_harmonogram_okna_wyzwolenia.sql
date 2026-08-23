-- Migracja 275 — okna wykonania, nadzór obecności uruchomień, wyzwalacz
-- webhook i historia wyzwoleń (rodzina `schedule.*`).
--
-- Wszystko tu ma przeżyć restart rdzenia i taki jest powód istnienia tego
-- kroku. Okno wykonania trzymane w pamięci procesu przestałoby ograniczać
-- budzik po pierwszym ponownym złożeniu rdzenia, a Scheduler pokazywałby
-- ograniczenie, które już niczego nie ogranicza. To samo dotyczy tolerancji
-- nadzoru i klucza podpisu webhooka.
--
-- Okno wykonania jest wierszem, bo jest ich wiele na harmonogram: godziny
-- robocze poniedziałku bywają inne niż piątku. Dni tygodnia leżą zapisem
-- strukturalnym, bo kontrakt niesie je wykazem liczb ustalanym w całości.
CREATE TABLE okno_wykonania_harmonogramu (
    id             INTEGER PRIMARY KEY AUTOINCREMENT,
    harmonogram_id INTEGER NOT NULL REFERENCES harmonogram_automatyki(id) ON DELETE CASCADE,
    dni_tygodnia   TEXT    NOT NULL DEFAULT '[]',
    minuta_od      INTEGER NOT NULL DEFAULT 0,
    minuta_do      INTEGER NOT NULL DEFAULT 1440,
    strefa_czasowa TEXT,
    kolejnosc      INTEGER NOT NULL DEFAULT 0
);
CREATE INDEX idx_okno_wykonania_harmonogramu_wlasciciel
    ON okno_wykonania_harmonogramu(harmonogram_id, kolejnosc, id);

-- Nadzór obecności uruchomień i adres wejściowy webhooka są polami pojedynczymi
-- harmonogramu — po jednym na harmonogram, bez własnego cyklu życia.
-- `tolerancja_sekundy` zero wyłącza nadzór, tak samo jak mówi kontrakt.
-- Kolumna webhooka niesie ODWOŁANIE do sejfu, nigdy wartość klucza podpisu.
ALTER TABLE harmonogram_automatyki ADD COLUMN tolerancja_sekundy INTEGER NOT NULL DEFAULT 0;
ALTER TABLE harmonogram_automatyki ADD COLUMN regula_nadzoru TEXT;
ALTER TABLE harmonogram_automatyki ADD COLUMN odwolanie_podpisu TEXT;
ALTER TABLE harmonogram_automatyki ADD COLUMN okno_deduplikacji_sekundy INTEGER NOT NULL DEFAULT 300;

-- Historia wyzwoleń zapisuje moment RZECZYWISTY wraz z przyczyną. Bez niej
-- Scheduler pokazywałby wyłącznie terminy wyliczone naprzód, a Operator nie
-- miałby jak sprawdzić, czy harmonogram kiedykolwiek odpalił.
CREATE TABLE wyzwolenie_automatyki (
    id             INTEGER PRIMARY KEY AUTOINCREMENT,
    automatyka_id  INTEGER NOT NULL REFERENCES automatyka(id) ON DELETE CASCADE,
    harmonogram_id INTEGER REFERENCES harmonogram_automatyki(id) ON DELETE SET NULL,
    przyczyna      TEXT    NOT NULL DEFAULT 'manual'
                           CHECK(przyczyna IN ('schedule','event','manual','chain','backfill')),
    przebieg_id    INTEGER REFERENCES przebieg_automatyki(id) ON DELETE SET NULL,
    chwila         TEXT    NOT NULL DEFAULT (strftime('%Y-%m-%dT%H:%M:%fZ','now'))
);
-- schedule.trigger.history czyta wyzwolenia od najnowszego.
CREATE INDEX idx_wyzwolenie_automatyki_czas
    ON wyzwolenie_automatyki(automatyka_id, chwila DESC, id DESC);
CREATE INDEX idx_wyzwolenie_automatyki_harmonogram
    ON wyzwolenie_automatyki(harmonogram_id, chwila DESC, id DESC);
