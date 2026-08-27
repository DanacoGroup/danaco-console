-- Migracja 102 zakłada historię wersji tożsamości eksperta oraz archiwizację
-- zamiast usunięcia w module Agents.

-- ── Archiwum eksperta ─────────────────────────────────────────────────────────
-- NULL znaczy ekspert czynny; format znacznika ISO 8601 UTC jest tym samym
-- co w pozostałych znacznikach czasu platformy.
ALTER TABLE agent ADD COLUMN zarchiwizowano_o TEXT;

-- Stan czynności sprzed archiwizacji. Przywrócenie oddaje eksperta takim,
-- jakim był, zamiast zgadywać, że każdy zarchiwizowany był czynny.
ALTER TABLE agent ADD COLUMN aktywny_przed_archiwum INTEGER;

CREATE INDEX idx_agent_archiwum ON agent (zarchiwizowano_o)
    WHERE zarchiwizowano_o IS NOT NULL;

-- ── Historia wersji tożsamości ────────────────────────────────────────────────
-- Historia wersji tożsamości powstaje wyzwalaczem bazy przy każdym założeniu
-- i każdej zmianie eksperta.
CREATE TABLE agent_wersja (
    id                   INTEGER PRIMARY KEY AUTOINCREMENT,
    agent_id             INTEGER NOT NULL REFERENCES agent(id) ON DELETE CASCADE,
    -- Numer odpowiada kolumnie `agent.wersja` z chwili zapisu migawki.
    numer                INTEGER NOT NULL,
    nazwa                TEXT    NOT NULL,
    opis                 TEXT    NOT NULL DEFAULT '',
    instrukcje_systemowe TEXT    NOT NULL DEFAULT '',
    kanal_kod            TEXT,
    model                TEXT,
    transport            TEXT,
    parametry_json       TEXT    NOT NULL DEFAULT '{}',
    imie_wlasne          TEXT    NOT NULL DEFAULT '',
    favikon              TEXT    NOT NULL DEFAULT '',
    ustawienia_json      TEXT    NOT NULL DEFAULT '',
    tryb_nakladki        TEXT    NOT NULL DEFAULT '',
    aktywny              INTEGER NOT NULL DEFAULT 1 CHECK (aktywny IN (0, 1)),
    -- Kto zapisał tę wersję; 'operator' albo 'restore' (przywrócenie).
    autor                TEXT    NOT NULL DEFAULT 'operator',
    -- Powód powstania wersji; przy przywróceniu niesie numer wersji źródłowej.
    powod                TEXT    NOT NULL DEFAULT '',
    zapisano             TEXT    NOT NULL DEFAULT (strftime('%Y-%m-%dT%H:%M:%fZ', 'now')),
    UNIQUE (agent_id, numer)
);

-- Historia czytana jest zawsze dla jednego eksperta, od najnowszej wersji,
-- więc indeks porządkuje wiersze malejąco po numerze wersji.
CREATE INDEX idx_agent_wersja_historia ON agent_wersja (agent_id, numer DESC);

-- ── Migawka stanu zastanego ───────────────────────────────────────────────────
-- Eksperci założeni przed tą migracją dostają jedną migawkę stanu bieżącego
-- jako początek historii wersji.
INSERT INTO agent_wersja
    (agent_id, numer, nazwa, opis, instrukcje_systemowe, kanal_kod, model, transport,
     parametry_json, imie_wlasne, favikon, ustawienia_json, tryb_nakladki, aktywny,
     autor, powod, zapisano)
SELECT id, wersja, nazwa, opis, instrukcje_systemowe, kanal_kod, model, transport,
       parametry_json, imie_wlasne, favikon, ustawienia_json, tryb_nakladki, aktywny,
       'operator', 'stan zastany przy migracji 102', zaktualizowano
FROM agent;

-- ── Wyzwalacze zakładające migawki ────────────────────────────────────────────
-- Wyzwalacz widzi każdą drogę zapisu tożsamości eksperta i każdą przyszłą,
-- więc historia nie zależy od tego, który kod dokonał zmiany.
CREATE TRIGGER agent_wersja_po_zalozeniu
AFTER INSERT ON agent
BEGIN
    INSERT INTO agent_wersja
        (agent_id, numer, nazwa, opis, instrukcje_systemowe, kanal_kod, model, transport,
         parametry_json, imie_wlasne, favikon, ustawienia_json, tryb_nakladki, aktywny, powod)
    VALUES
        (new.id, new.wersja, new.nazwa, new.opis, new.instrukcje_systemowe, new.kanal_kod,
         new.model, new.transport, new.parametry_json, new.imie_wlasne, new.favikon,
         new.ustawienia_json, new.tryb_nakladki, new.aktywny, 'zalozenie eksperta');
END;

CREATE TRIGGER agent_wersja_po_zmianie
AFTER UPDATE ON agent
WHEN new.nazwa            IS NOT old.nazwa
  OR new.opis             IS NOT old.opis
  OR new.instrukcje_systemowe IS NOT old.instrukcje_systemowe
  OR new.kanal_kod        IS NOT old.kanal_kod
  OR new.model            IS NOT old.model
  OR new.transport        IS NOT old.transport
  OR new.parametry_json   IS NOT old.parametry_json
  OR new.imie_wlasne      IS NOT old.imie_wlasne
  OR new.favikon          IS NOT old.favikon
  OR new.ustawienia_json  IS NOT old.ustawienia_json
  OR new.tryb_nakladki    IS NOT old.tryb_nakladki
  OR new.aktywny          IS NOT old.aktywny
  OR new.wersja           IS NOT old.wersja
BEGIN
    INSERT INTO agent_wersja
        (agent_id, numer, nazwa, opis, instrukcje_systemowe, kanal_kod, model, transport,
         parametry_json, imie_wlasne, favikon, ustawienia_json, tryb_nakladki, aktywny, powod)
    VALUES
        (new.id, new.wersja, new.nazwa, new.opis, new.instrukcje_systemowe, new.kanal_kod,
         new.model, new.transport, new.parametry_json, new.imie_wlasne, new.favikon,
         new.ustawienia_json, new.tryb_nakladki, new.aktywny, '')
    ON CONFLICT (agent_id, numer) DO UPDATE SET
        nazwa                = excluded.nazwa,
        opis                 = excluded.opis,
        instrukcje_systemowe = excluded.instrukcje_systemowe,
        kanal_kod            = excluded.kanal_kod,
        model                = excluded.model,
        transport            = excluded.transport,
        parametry_json       = excluded.parametry_json,
        imie_wlasne          = excluded.imie_wlasne,
        favikon              = excluded.favikon,
        ustawienia_json      = excluded.ustawienia_json,
        tryb_nakladki        = excluded.tryb_nakladki,
        aktywny              = excluded.aktywny,
        zapisano             = strftime('%Y-%m-%dT%H:%M:%fZ', 'now');
END;
