-- Migracja 280 — spięcie silnika kolejek Automations z silnikiem kolejek
-- środowiska MultitaskingAI (`orchestration.multitasking.link`).
--
-- Stan wyjściowy to rozłączone, więc brak wiersza znaczy brak spięcia i tabela
-- nie zakłada niczego z góry.
--
-- Spięcie nie jest znacznikiem. Silnik kolejek jest w rdzeniu jeden — drugiego
-- nie ma i mieć nie będzie — a środowisko MultitaskingAI odróżnia się od pętli
-- sesyjnej tym, czyim koordynatorem kolejka jest prowadzona
-- (`kolejka.okno_koordynatora_id`) i jakiego jest rodzaju. Spięcie przestawia
-- właśnie te dwie rzeczy na kolejkach automatyki, a rozłączenie je zdejmuje.
-- Wiersz poniżej niesie zapis decyzji; skutek leży w tabeli `kolejka`.

CREATE TABLE orkiestracja_spiecie_multitasking (
    automatyka_id INTEGER PRIMARY KEY REFERENCES automatyka(id) ON DELETE CASCADE,
    -- Okno roli środowiska MultitaskingAI, której podlega automatyka. NULL
    -- znaczy spięcie bez wskazania roli — kolejki idą wtedy rodzajem
    -- `multitasking`, ale bez koordynatora.
    okno_roli_id  INTEGER REFERENCES okno_komunikacji(id) ON DELETE SET NULL,
    zapisano      TEXT    NOT NULL DEFAULT (strftime('%Y-%m-%dT%H:%M:%fZ','now'))
);
