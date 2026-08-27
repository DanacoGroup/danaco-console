-- Migracja 373 zakłada tabelę wyłączeń pamięci w zasięgu, pozwalającą
-- wstrzymać obowiązywanie ustalenia bez usunięcia ani zawężenia jego treści.

CREATE TABLE IF NOT EXISTS wylaczenie_pamieci (
    id                     INTEGER PRIMARY KEY AUTOINCREMENT,
    identyfikator_zewnetrzny TEXT  NOT NULL UNIQUE,
    wpis_id                INTEGER          REFERENCES wpis_pamieci_projektu(id) ON DELETE CASCADE,
    poziom_pamieci         TEXT    NOT NULL DEFAULT ''
                                   CHECK(poziom_pamieci IN ('', 'global', 'environment',
                                                            'project', 'session')),
    poziom_zasiegu_id      INTEGER NOT NULL REFERENCES poziom_zasiegu(id) ON DELETE CASCADE,
    klucz_zasiegu          TEXT    NOT NULL DEFAULT '',
    utworzono              TEXT    NOT NULL DEFAULT (strftime('%Y-%m-%dT%H:%M:%fZ','now')),
    zaktualizowano         TEXT    NOT NULL DEFAULT (strftime('%Y-%m-%dT%H:%M:%fZ','now')),
    -- Wyłączenie bez wskazanego bytu nie wyłącza niczego; wskazanie obu bytów naraz jest niedopuszczalne.
    CHECK((wpis_id IS NOT NULL AND poziom_pamieci = '')
          OR (wpis_id IS NULL AND poziom_pamieci <> ''))
);

-- Indeks jednoznaczny ogranicza wyłączenie tego samego bytu i zasięgu do jednego wiersza,
-- tak aby powtórzone żądanie nie zakładało kolejnego.
CREATE UNIQUE INDEX IF NOT EXISTS idx_wylaczenie_pamieci_byt
    ON wylaczenie_pamieci(IFNULL(wpis_id, 0), poziom_pamieci,
                          poziom_zasiegu_id, klucz_zasiegu);

CREATE INDEX IF NOT EXISTS idx_wylaczenie_pamieci_zasieg
    ON wylaczenie_pamieci(poziom_zasiegu_id, klucz_zasiegu);
