-- Migracja 040 — trwałość modułu Automations, część druga: harmonogram
-- (Scheduler) i przebieg automatyki (Execution Monitor).
--
-- Harmonogram jeden na automatykę. Wiersz inwentarza okna Scheduler mówi:
-- „Harmonogram obowiązuje po powiązaniu z workflow”, a kontrakt ma jedną komendę
-- `automation.schedule.set` przyjmującą `workflowId`. Więz UNIQUE na kolumnie
-- `automatyka_id` zamyka drogę do dwóch sprzecznych harmonogramów tej samej
-- automatyki — Operator nie zobaczyłby wtedy, który obowiązuje.
--
-- Cykliczność jest wyzwalaczem, a nie jego przeciwieństwem. Kontrakt niesie
-- notację cron dwa razy: polem `AutomationSchedule.cron` i wyzwalaczem rodzaju
-- `cron`. Kolumna `cron` jest cyklicznością podstawową harmonogramu, a tabela
-- wyzwalaczy — pozostałymi zdarzeniami wyzwalającymi (webhook, plik, warunek)
-- oraz cyklicznościami dodatkowymi. Najbliższe uruchomienie liczy rdzeń
-- z obu źródeł, więc kolumna `nastepne_uruchomienie` jest wynikiem, nie
-- deklaracją Operatora.
--
-- Przebieg jest zapisem wykonania, nie drugą kolejką. Kroki wykonuje jeden
-- silnik kolejek (migracja 003). Tabela `przebieg_automatyki` wiąże
-- uruchomienie automatyki z kolejką, która je realizuje, i przechowuje to,
-- czego kolejka nie wie: którą automatykę wykonuje, ile miała etapów i dlaczego
-- się nie powiodła. Etap bieżący i liczba etapów pochodzą z pozycji kolejki —
-- kolumny są ich utrwaleniem na potrzeby Execution Monitora po restarcie
-- rdzenia, a nie drugim licznikiem postępu.

-- ── Harmonogram automatyki — Scheduler ────────────────────────────────────────
CREATE TABLE harmonogram_automatyki (
    id                       INTEGER PRIMARY KEY AUTOINCREMENT,
    identyfikator_zewnetrzny TEXT    NOT NULL UNIQUE,
    automatyka_id            INTEGER NOT NULL UNIQUE
                                     REFERENCES automatyka(id) ON DELETE CASCADE,
    cron                     TEXT,
    strefa_czasowa           TEXT,
    czynny                   INTEGER NOT NULL DEFAULT 1 CHECK(czynny IN (0,1)),
    nastepne_uruchomienie    TEXT,
    utworzono                TEXT    NOT NULL DEFAULT (strftime('%Y-%m-%dT%H:%M:%fZ','now')),
    zaktualizowano           TEXT    NOT NULL DEFAULT (strftime('%Y-%m-%dT%H:%M:%fZ','now'))
);

-- ── Wyzwalacz poza cyklicznością podstawową ───────────────────────────────────
CREATE TABLE wyzwalacz_automatyki (
    id                       INTEGER PRIMARY KEY AUTOINCREMENT,
    harmonogram_id           INTEGER NOT NULL
                                     REFERENCES harmonogram_automatyki(id) ON DELETE CASCADE,
    identyfikator_zewnetrzny TEXT    NOT NULL UNIQUE,
    rodzaj                   TEXT    NOT NULL
                                     CHECK(rodzaj IN ('cron','webhook','file','condition')),
    wyrazenie                TEXT    NOT NULL,
    czynny                   INTEGER NOT NULL DEFAULT 1 CHECK(czynny IN (0,1)),
    kolejnosc                INTEGER NOT NULL DEFAULT 0,
    utworzono                TEXT    NOT NULL DEFAULT (strftime('%Y-%m-%dT%H:%M:%fZ','now'))
);
CREATE INDEX idx_wyzwalacz_automatyki_harmonogram
    ON wyzwalacz_automatyki(harmonogram_id, kolejnosc, id);

-- ── Przebieg automatyki — Execution Monitor ───────────────────────────────────
-- Kolumna `proba` niesie numer biegu naprawczego. Bieg naprawczy nie ma limitu
-- obiegów, więc kolumna nie ma ani górnej granicy, ani więzu, który
-- by ją narzucił.
CREATE TABLE przebieg_automatyki (
    id                       INTEGER PRIMARY KEY AUTOINCREMENT,
    identyfikator_zewnetrzny TEXT    NOT NULL UNIQUE,
    automatyka_id            INTEGER NOT NULL REFERENCES automatyka(id) ON DELETE CASCADE,
    kolejka_id               INTEGER REFERENCES kolejka(id) ON DELETE SET NULL,
    stan                     TEXT    NOT NULL DEFAULT 'pending'
                                     CHECK(stan IN ('pending','running','paused',
                                                    'succeeded','failed','stopped')),
    etap_biezacy             INTEGER NOT NULL DEFAULT 0,
    etapow                   INTEGER NOT NULL DEFAULT 0,
    proba                    INTEGER NOT NULL DEFAULT 0,
    komunikat_bledu          TEXT,
    rozpoczeto               TEXT    NOT NULL DEFAULT (strftime('%Y-%m-%dT%H:%M:%fZ','now')),
    zakonczono               TEXT
);
CREATE INDEX idx_przebieg_automatyki_automatyka
    ON przebieg_automatyki(automatyka_id, id DESC);
CREATE UNIQUE INDEX idx_przebieg_automatyki_kolejka
    ON przebieg_automatyki(kolejka_id) WHERE kolejka_id IS NOT NULL;
