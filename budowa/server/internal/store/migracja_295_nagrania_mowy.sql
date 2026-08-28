-- Migracja 295 — magazyn nagrań mowy (`speech.audio.upload`, `speech.audio.fetch`).
--
-- Brak, który ta tabela zamyka, jest platformowy, nie okienny: `speech.transcribe`
-- bierze ŚCIEŻKĘ PLIKU na maszynie silnika, a nagranie z mikrofonu istnieje
-- wyłącznie jako bajty w pamięci karty. Dopóki nie ma dokąd tych bajtów odłożyć,
-- ani dyktowanie w oknie komunikacji, ani mikrofon Voice Console nie mają czym
-- dojechać do silnika.
--
-- Bajty NIE leżą w bazie. Kolumna `sciezka` niesie bezwzględną ścieżkę pliku pod
-- katalogiem danych rdzenia — dokładnie taką, jaką przyjmuje `speech.transcribe`
-- i jaką oddaje `translate.speech.synthesize`. Dzięki temu odnośnik nagrania jest
-- jednym rodzajem odnośnika w całym produkcie, a nie dwoma. Nagranie nie opuszcza
-- maszyny rdzenia: droga sieciowa prowadzi tu i z powrotem, nigdzie indziej.
--
-- Ten sam wiersz obsługuje odsłuch (`speech.audio.fetch`): wykaz odnośników
-- mintowanych przez rdzeń jest tym, co pozwala oddać bajty, nie otwierając przy
-- okazji drogi do czytania dowolnego pliku dysku.
--
-- `wygasa` jest chwilą, po której nagranie przestaje być potrzebne. Puste znaczy
-- nagranie trwałe (Operator włączył zapis nagrań). Wygaszaniem zajmuje się
-- adapter przy kolejnym przyjęciu — osobnego zegara ta tabela nie potrzebuje.

CREATE TABLE nagranie_mowy (
    id                       INTEGER PRIMARY KEY AUTOINCREMENT,
    identyfikator_zewnetrzny TEXT    NOT NULL UNIQUE,
    sciezka                  TEXT    NOT NULL UNIQUE,
    typ_tresci               TEXT    NOT NULL,
    rozmiar_bajtow           INTEGER NOT NULL,
    dlugosc_ms               INTEGER,
    sesja_kod                TEXT,
    okno_kod                 TEXT,
    trwale                   INTEGER NOT NULL DEFAULT 0,
    utworzono                INTEGER NOT NULL,
    wygasa                   INTEGER
);
CREATE INDEX idx_nagranie_mowy_okno ON nagranie_mowy(okno_kod, utworzono DESC);
CREATE INDEX idx_nagranie_mowy_wygasa ON nagranie_mowy(wygasa);
