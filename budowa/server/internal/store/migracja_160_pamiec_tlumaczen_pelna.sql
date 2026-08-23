-- Migracja 160 — pamięć tłumaczeń modułu Translate w kształcie, którego żąda
-- kontrakt rodziny `translate.memory.*`, oraz polityka pamięci okna.
--
-- Tabela `pamiec_tlumaczen` (migracja 054) powstała pod jedną komendę
-- (`memory.suggest`) i pod jedną drogę zapisu: parę segmentów zdjętą
-- z zatwierdzonego panelu. Stąd `panel_id NOT NULL`. Rodzina `memory.*` żąda
-- czego innego: `memory.set` zakłada parę BEZ panelu (Operator wpisuje ją
-- wprost), `memory.import` wnosi parę z pliku wymiany, a kontrakt
-- (`TranslationMemoryEntry`) niesie projekt, klienta, autora i kontekst sąsiedni.
-- Kolumny da się dołożyć poleceniem ALTER; zdjęcia warunku NOT NULL z kolumny
-- `panel_id` już nie — SQLite tego nie umie. Dlatego tabela powstaje od nowa,
-- a wiersze zastane przechodzą do niej przepisaniem: pary zebrane dotąd
-- z paneli zostają parami z panelem, reszta pól zostaje pusta, bo nikt jej
-- nigdy nie podał.
--
-- Zasięg pary (`zasieg`) niesie wartości `TranslationMemoryScope` kontraktu
-- wprost. Domyślną jest `card` — pamięć własna karcie sesji; pary szersze
-- (projekt, zespół) powstają, gdy Operator wskaże zasięg sam.
--
-- Polityka pamięci jest osobną tabelą, nie kolumnami okna: `memory.policy.get`
-- odpowiada wtedy „polityki nie ustawiono” brakiem wiersza, zamiast czterema
-- kolumnami NULL, których nie da się odróżnić od polityki wyzerowanej.

-- ── Pamięć tłumaczeń — para segmentów (Translation Memory) ──────────────────
CREATE TABLE pamiec_tlumaczen_nowa (
    id                       INTEGER PRIMARY KEY AUTOINCREMENT,
    identyfikator_zewnetrzny TEXT    NOT NULL UNIQUE,
    -- Panel, z którego para została zdjęta. NULL znaczy „para wniesiona wprost”
    -- (memory.set) albo „para z pliku wymiany” (memory.import) — nie brak danych.
    panel_id                 INTEGER          REFERENCES panel_tlumaczenia(id) ON DELETE SET NULL,
    jezyk                    TEXT    NOT NULL,
    segment_zrodlowy         TEXT    NOT NULL,
    segment_docelowy         TEXT    NOT NULL,
    projekt                  TEXT,
    klient                   TEXT,
    autor                    TEXT,
    -- Kontekst sąsiedni: zdanie przed i po. Dopasowanie kontekstowe
    -- (`TranslationMemoryPolicy.contextMatch`) podnosi wynik pary, której
    -- sąsiedztwo zgadza się z sąsiedztwem segmentu tłumaczonego.
    kontekst_poprzedni       TEXT,
    kontekst_nastepny        TEXT,
    zasieg                   TEXT    NOT NULL DEFAULT 'card'
                                     CHECK(zasieg IN ('card','project','team')),
    utworzono                INTEGER NOT NULL,
    zaktualizowano           INTEGER NOT NULL
);

INSERT INTO pamiec_tlumaczen_nowa
    (identyfikator_zewnetrzny, panel_id, jezyk, segment_zrodlowy, segment_docelowy,
     zasieg, utworzono, zaktualizowano)
SELECT identyfikator_zewnetrzny, panel_id, jezyk, segment_zrodlowy, segment_docelowy,
       'card', utworzono, utworzono
  FROM pamiec_tlumaczen;

DROP TABLE pamiec_tlumaczen;
ALTER TABLE pamiec_tlumaczen_nowa RENAME TO pamiec_tlumaczen;

CREATE INDEX idx_pamiec_tlumaczen_jezyk_segment ON pamiec_tlumaczen(jezyk, segment_zrodlowy);
CREATE INDEX idx_pamiec_tlumaczen_panel ON pamiec_tlumaczen(panel_id);
CREATE INDEX idx_pamiec_tlumaczen_projekt ON pamiec_tlumaczen(projekt, jezyk);

-- ── Polityka pamięci okna (`translate.memory.policy.*`) ─────────────────────
-- `prog` jest progiem dopasowania w procentach (kontrakt: `threshold`).
CREATE TABLE polityka_pamieci_okna (
    id                     INTEGER PRIMARY KEY AUTOINCREMENT,
    okno_id                INTEGER NOT NULL UNIQUE REFERENCES okno_tlumaczenia(id) ON DELETE CASCADE,
    zasieg                 TEXT    NOT NULL DEFAULT 'card'
                                   CHECK(zasieg IN ('card','project','team')),
    prog                   INTEGER NOT NULL DEFAULT 75,
    dopasowanie_kontekstu  INTEGER NOT NULL DEFAULT 0 CHECK(dopasowanie_kontekstu IN (0,1)),
    wstepne_tlumaczenie    INTEGER NOT NULL DEFAULT 0 CHECK(wstepne_tlumaczenie IN (0,1)),
    zaktualizowano         INTEGER NOT NULL
);
