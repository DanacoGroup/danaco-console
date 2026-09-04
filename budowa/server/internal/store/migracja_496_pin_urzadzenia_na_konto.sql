-- PIN i klucz Hello należą do urządzenia, a wskaźnik jednoznaczności pary
-- (urządzenie, rodzaj) nie znał konta. Na maszynie dzielonej przez dwóch
-- Operatorów drugi nie mógł założyć własnego PIN-u: para była już zajęta przez
-- pierwszego, a rdzeń czytał to jako zmianę PIN-u zastanego.
--
-- Kotwica hasła zna konto od kroku 485 (`idx_metoda_kotwica_konta`); ten krok
-- domyka drugą metodę wejścia tym samym wzorem. Wskaźnik jest zwykły, nie
-- wbudowany w tabelę, więc wystarczy go założyć na nowo — przebudowy nie ma.

DROP INDEX IF EXISTS idx_metoda_uwierzytelnienia_urzadzenie;

CREATE UNIQUE INDEX idx_metoda_uwierzytelnienia_urzadzenie
    ON metoda_uwierzytelnienia (urzadzenie_kod, rodzaj, COALESCE(konto_id, 0))
    WHERE urzadzenie_kod IS NOT NULL;
