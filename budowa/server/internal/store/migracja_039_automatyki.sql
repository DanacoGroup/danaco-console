-- Migracja 039 — trwałość modułu Automations, część pierwsza: definicja
-- automatyki (Workflow Builder) i układ zależności między jej krokami
-- (Orchestrator).
--
-- Automatyka jest komponentem własnym, nie bytem sesji. Dlatego tabela nie ma
-- kolumny sesji ani okna: ta sama automatyka jest widoczna ze strony głównej
-- niezależnie od karty sesji, w której Operator akurat pracuje.
--
-- Silnika kolejek tu nie ma. Wykonaniem kroków zajmuje się jeden silnik
-- kolejek: krok automatyki staje się pozycją kolejki dopiero w chwili
-- uruchomienia. Tabele poniżej opisują wyłącznie definicję — to, co Operator
-- zbudował, zanim cokolwiek ruszyło.
--
-- Zależność ma jedno miejsce zapisu. Kontrakt niesie zależność dwa razy: raz
-- polem `AutomationStep.dependsOn`, raz strukturą `AutomationDependency`.
-- Dwie tabele znaczyłyby dwie prawdy o tym samym łuku grafu, więc prawda jest
-- jedna — tabela `zaleznosc_kroku_automatyki` — a `dependsOn` powstaje z niej
-- przy odczycie.

-- ── Definicja automatyki ──────────────────────────────────────────────────────
CREATE TABLE automatyka (
    id                       INTEGER PRIMARY KEY AUTOINCREMENT,
    identyfikator_zewnetrzny TEXT    NOT NULL UNIQUE,
    nazwa                    TEXT    NOT NULL,
    opis                     TEXT,
    czynna                   INTEGER NOT NULL DEFAULT 1 CHECK(czynna IN (0,1)),
    -- Wersja rośnie przy każdym zapisie definicji. Panel akcji Workflow Buildera
    -- ma pozycję „Wersje”, a bez licznika nie da się powiedzieć, którą wersję
    -- Operator właśnie ogląda.
    wersja                   INTEGER NOT NULL DEFAULT 1,
    utworzono                TEXT    NOT NULL DEFAULT (strftime('%Y-%m-%dT%H:%M:%fZ','now')),
    zaktualizowano           TEXT    NOT NULL DEFAULT (strftime('%Y-%m-%dT%H:%M:%fZ','now'))
);
CREATE INDEX idx_automatyka_nazwa ON automatyka(nazwa, id);

-- ── Krok automatyki ───────────────────────────────────────────────────────────
-- Wartości kolumny `rodzaj` są wartościami kontraktu (AutomationStepKind), nie
-- ich tłumaczeniem. Dzięki temu przekład wiersza na krok kontraktu nie
-- potrzebuje słownika pośredniego.
CREATE TABLE krok_automatyki (
    id                       INTEGER PRIMARY KEY AUTOINCREMENT,
    automatyka_id            INTEGER NOT NULL REFERENCES automatyka(id) ON DELETE CASCADE,
    identyfikator_zewnetrzny TEXT    NOT NULL,
    nazwa                    TEXT,
    rodzaj                   TEXT    NOT NULL DEFAULT 'command'
                                     CHECK(rodzaj IN ('command','model','condition','branch','wait')),
    komenda                  TEXT,
    parametry                TEXT,
    warunek                  TEXT,
    kolejnosc                INTEGER NOT NULL DEFAULT 0,
    utworzono                TEXT    NOT NULL DEFAULT (strftime('%Y-%m-%dT%H:%M:%fZ','now')),
    zaktualizowano           TEXT    NOT NULL DEFAULT (strftime('%Y-%m-%dT%H:%M:%fZ','now')),
    UNIQUE (automatyka_id, identyfikator_zewnetrzny)
);
CREATE INDEX idx_krok_automatyki_kolejnosc ON krok_automatyki(automatyka_id, kolejnosc, id);

-- ── Zależność między krokami — Orchestrator ───────────────────────────────────
-- Więz pierwotny obejmuje parę kroków, więc ten sam łuk nie powstanie dwa razy.
-- Pętli własnej (krok zależny od siebie) schemat nie dopuszcza wprost; cykl
-- dłuższy wykrywa walidacja układu, bo SQLite nie ma na to więzu.
CREATE TABLE zaleznosc_kroku_automatyki (
    automatyka_id  INTEGER NOT NULL REFERENCES automatyka(id) ON DELETE CASCADE,
    krok_z         TEXT    NOT NULL,
    krok_do        TEXT    NOT NULL,
    rodzaj         TEXT    NOT NULL DEFAULT 'sequential'
                           CHECK(rodzaj IN ('sequential','parallel','conditional')),
    warunek        TEXT,
    utworzono      TEXT    NOT NULL DEFAULT (strftime('%Y-%m-%dT%H:%M:%fZ','now')),
    PRIMARY KEY (automatyka_id, krok_z, krok_do),
    CHECK (krok_z <> krok_do)
);
CREATE INDEX idx_zaleznosc_kroku_do ON zaleznosc_kroku_automatyki(automatyka_id, krok_do);
