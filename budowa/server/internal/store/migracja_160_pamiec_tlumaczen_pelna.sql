-- Migracja 160 przebudowuje tabelę pamiec_tlumaczen do kształtu kontraktu rodziny translate.memory, luzując wymóg panelu, i zakłada tabelę polityki pamięci okna.

-- Zakłada tabelę pamiec_tlumaczen_nowa niosącą parę segmentów pamięci tłumaczeń wraz z kontekstem sąsiednim i zasięgiem widoczności pary.
CREATE TABLE pamiec_tlumaczen_nowa (
    id                       INTEGER PRIMARY KEY AUTOINCREMENT,
    identyfikator_zewnetrzny TEXT    NOT NULL UNIQUE,
    -- Panel, z którego zdjęto parę; wartość pusta oznacza parę wniesioną wprost albo z pliku wymiany.
    panel_id                 INTEGER          REFERENCES panel_tlumaczenia(id) ON DELETE SET NULL,
    jezyk                    TEXT    NOT NULL,
    segment_zrodlowy         TEXT    NOT NULL,
    segment_docelowy         TEXT    NOT NULL,
    projekt                  TEXT,
    klient                   TEXT,
    autor                    TEXT,
    -- Zdanie przed i po segmentem; dopasowanie kontekstowe podnosi wynik pary o zgodnym sąsiedztwie.
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
