-- Migracja 095 — kanał drugiego silnika: Bielik przez Ollama.
--
-- Silnik: Ollama serwuje API zgodne z OpenAI pod http://127.0.0.1:11434/v1/…
-- lokalnie, CPU-only dozwolone. Kanał jedzie adapterem `api`
-- (models/adapter_api.go) — generycznym kanałem HTTP bez zaszytego dostawcy;
-- domyślne ścieżki odczytu (`choices.0.message.content`
-- i `choices.0.delta.content`) są dokładnie ścieżkami odpowiedzi OpenAI,
-- więc żaden nowy adapter nie powstaje: nowy silnik to wiersz rejestru, nie
-- gałąź w kodzie.
--
-- Kanał stoi z `aktywny = 0`, bo silnika nie ma w standardowej instalacji.
-- Kanał czynny bez silnika byłby przyciskiem-atrapą: okno dałoby się na niego
-- przełączyć, a każda tura kończyłaby się odmową. Wiersz stoi w rejestrze (droga
-- jest zbudowana), a włącza go Operator komendą `channel.update` (enabled=true)
-- po zainstalowaniu silnika: `ollama pull <model>` i start usługi. Włączony kanał
-- bez silnika odmawia jawnie błędem kanału — nie udaje odpowiedzi.
--
-- Identyfikator modelu to znacznik po stronie Ollamy; wpisany jest bazowy
-- `bielik` — Operator wskazuje dokładny znacznik (np. wariant kwantyzacji)
-- komendą `channel.update` polem model, bo znacznik zależy od tego, co
-- faktycznie pobrał (konfiguracja, nie kod).
--
-- Parametr `kanal_zapasowy` celowo nieobecny: kolejność zapasowa jest nastawą
-- Operatora (channel.update, config.kanal_zapasowy), nie rozstrzygnięciem
-- schematu — migracja 094 niesie mechanizm i ślad, nie politykę.

INSERT INTO kanal_modelu (kod, nazwa, dostawca, identyfikator_modelu, rodzaj_kanalu,
                          parametry_json, multimodalny, aktywny, kolejnosc)
VALUES ('bielik-ollama', 'Bielik (Ollama, lokalnie)', 'speakleash', 'bielik', 'api',
        '{"base_url":"http://127.0.0.1:11434/v1/chat/completions"}', 0, 0, 2)
ON CONFLICT(kod) DO NOTHING;
