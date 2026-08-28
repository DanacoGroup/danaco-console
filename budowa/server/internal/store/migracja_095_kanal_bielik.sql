-- Migracja dokłada kanał drugiego silnika modelu, Bielik przez lokalną usługę, stojący
-- nieaktywny do czasu instalacji silnika.

INSERT INTO kanal_modelu (kod, nazwa, dostawca, identyfikator_modelu, rodzaj_kanalu,
                          parametry_json, multimodalny, aktywny, kolejnosc)
VALUES ('bielik-ollama', 'Bielik (Ollama, lokalnie)', 'speakleash', 'bielik', 'api',
        '{"base_url":"http://127.0.0.1:11434/v1/chat/completions"}', 0, 0, 2)
ON CONFLICT(kod) DO NOTHING;
