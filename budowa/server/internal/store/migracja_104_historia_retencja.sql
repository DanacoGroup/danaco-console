-- Migracja 104 — zasada przechowywania historii.
--
-- Rodzina `history.*` dostaje tu jedną nową tabelę: zasadę przechowywania.
-- Drugiej tabeli „historia rozmowy okna" nie ma, bo historia rozmowy okna już
-- istnieje:
--
--   `wiadomosc`        — wypowiedź: rola, treść, czas, kolejność. To jest pozycja
--                      historii; `history.load` czyta ten wykaz.
--   `blok_wiadomosci`  — nietekstowe fragmenty jednej wypowiedzi (rozumowanie,
--                      narzędzia, prowenancja). To nie jest pozycja historii,
--                      tylko jej wnętrze — nie niesie ani roli, ani tekstu, więc
--                      `HistoryEntry` (role, preview) nie miałby się z czego złożyć.
--
-- Założenie trzeciej tabeli na to samo byłoby drugą prawdą o jednej wypowiedzi.
--
-- ── zasada_przechowywania ───────────────────────────────────────────────────
-- Jeden wiersz = jedna zasada dla jednego zakresu. Zakres rozstrzyga się od
-- najwęższego: okno wygrywa z sesją, sesja z globalną (ta sama kolejność, co
-- w warstwowej konfiguracji). Zasady się nie sumują — obowiązuje jedna,
-- najbliższa oknu; sumowanie dawałoby wynik, którego Operator nie umiałby
-- przewidzieć z żadnego pojedynczego ekranu.
--
-- Zakres niesie wartość kontraktu (angielską), nie przekład. Kontrakt dla
-- `RetentionPolicy.scope` żadnego polskiego odwzorowania nie zapisuje, a przekład
-- wolno trzymać wyłącznie w kontrakcie. Własny słownik polski byłby drugim
-- źródłem odwzorowania.
--
-- `zakres_kod` jest identyfikatorem kontraktowym (napisem), nie kluczem obcym.
-- Zasada ma przeżyć okno, którego dotyczy: Operator ustawia retencję na oknie,
-- okno znika wraz z sesją, a zasada zostaje wierszem bez skutku zamiast znikać
-- kaskadą po cichu. Dla zakresu `global` kolumna jest pusta (NULL) — bytu do
-- wskazania nie ma.
--
-- Oba progi są niewymagane i oba mogą być puste naraz. Brak progu znaczy „bez
-- ograniczenia", a nie zero: wiersz z dwoma NULL-ami to zasada wyłączona
-- i egzekucja nie usuwa wtedy niczego. Wartość 0 jest odbita warunkiem CHECK, bo
-- „trzymaj zero pozycji" i „nie ograniczaj" to dwie różne rzeczy, których jedna
-- liczba nie uniesie.
--
-- Tabela wchodzi pusta, więc żadne okno nie ma zasady, a egzekucja nie ma czego
-- zastosować. Retencja zaczyna działać w chwili, w której Operator ją ustawi —
-- nie wcześniej.

CREATE TABLE zasada_przechowywania (
    id               INTEGER PRIMARY KEY AUTOINCREMENT,
    zakres           TEXT    NOT NULL
                             CHECK (zakres IN ('window','session','global')),
    zakres_kod       TEXT,
    dni_trzymania    INTEGER CHECK (dni_trzymania    IS NULL OR dni_trzymania    > 0),
    pozycje_trzymane INTEGER CHECK (pozycje_trzymane IS NULL OR pozycje_trzymane > 0),
    utworzono        TEXT    NOT NULL DEFAULT (strftime('%Y-%m-%dT%H:%M:%fZ','now')),
    zaktualizowano   TEXT    NOT NULL DEFAULT (strftime('%Y-%m-%dT%H:%M:%fZ','now')),

    -- Zakres `global` ma być JEDEN, nie jeden na każdy pusty kod.
    CHECK ((zakres = 'global' AND zakres_kod IS NULL)
        OR (zakres <> 'global' AND zakres_kod IS NOT NULL AND zakres_kod <> ''))
);

-- Jedna zasada na zakres — druga byłaby drugą prawdą o tym samym oknie.
-- Wyrażenie COALESCE zamiast pary kolumn, bo NULL nie równa się NULL i zwykły
-- indeks UNIQUE przepuściłby dowolną liczbę wierszy `global`.
CREATE UNIQUE INDEX idx_zasada_przechowywania_zakres
    ON zasada_przechowywania (zakres, COALESCE(zakres_kod, ''));
