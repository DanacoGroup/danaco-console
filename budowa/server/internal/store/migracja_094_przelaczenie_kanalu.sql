-- Migracja zakłada ślad przełączeń kanału jako jawny zapis, z którego kanału na który
-- tura pojechała i dlaczego.

CREATE TABLE przelaczenie_kanalu (
    id            INTEGER PRIMARY KEY,
    chwila        INTEGER NOT NULL,          -- epoka w milisekundach
    okno_kod      TEXT    NOT NULL,
    wiadomosc_kod TEXT    NOT NULL,
    z_kanalu      TEXT    NOT NULL,          -- kod kanału, który odmówił
    na_kanal      TEXT    NOT NULL,          -- kod kanału, którym tura pojechała dalej
    powod         TEXT    NOT NULL DEFAULT ''
);

CREATE INDEX idx_przelaczenie_okno ON przelaczenie_kanalu (okno_kod, chwila);
