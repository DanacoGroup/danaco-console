-- Migracja zakłada zaczyn kanału lokalnego, żeby świeża instalacja mogła rozmawiać
-- bez ręcznej konfiguracji od pierwszego okna.

INSERT INTO kanal_modelu (kod, nazwa, dostawca, identyfikator_modelu, rodzaj_kanalu,
                          parametry_json, multimodalny, aktywny, kolejnosc)
VALUES ('lokalny-claude', 'Kanał lokalny (claude)', 'anthropic', '', 'cli',
        '{}', 0, 1, 0)
ON CONFLICT(kod) DO NOTHING;
