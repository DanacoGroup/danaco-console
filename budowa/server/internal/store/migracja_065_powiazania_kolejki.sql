-- Migracja 065 — powiązania kolejki (komenda `queue.link`).
--
-- Kontrakt niesie komendę `queue.link`: „Wiaze kolejke z oknami, ekspertem,
-- projektem albo automatyka". Sam schemat kolejki zna jedno wiązanie poza sesją
-- — kolumnę `kolejka.okno_koordynatora_id` wskazującą jedno okno komunikacji.
-- Kontrakt mówi o `windowIds` w liczbie mnogiej i o trzech dalszych bytach, dla
-- których w tabeli `kolejka` nie ma kolumny.
--
-- Wiązania niesie jedna tabela, a nie cztery kolumny w `kolejka`: kolumny
-- obsłużyłyby po jednym bycie każdego rodzaju, a kontrakt niesie wykaz okien.
-- Tabela wiążąca obsługuje wszystkie cztery rodzaje jednym kształtem i nie
-- wymaga zmiany schematu przy piątym.
--
-- `byt` jest tekstem, a nie kluczem obcym, bo kontrakt niesie identyfikatory
-- zewnętrzne: `windowIds` to identyfikatory okien kontraktu (nie
-- `okno_komunikacji.id`), `agentId` to kod eksperta, `projectId` — kod projektu,
-- `workflowId` — identyfikator automatyki. Cztery klucze obce wymagałyby
-- czterech rozwiązań identyfikatora w chwili zapisu i odmawiałyby powiązania
-- z bytem, który nie ma jeszcze wiersza. Kasowanie kolejki zabiera powiązania
-- kaskadą, bo to one wiszą na kolejce, nie odwrotnie.
--
-- Kolumna `kolejka.okno_koordynatora_id` zostaje nietknięta: opisuje okno
-- koordynatora pętli, a nie okna obsługiwane przez kolejkę.

CREATE TABLE powiazanie_kolejki (
    id            INTEGER PRIMARY KEY AUTOINCREMENT,
    kolejka_id    INTEGER NOT NULL REFERENCES kolejka(id) ON DELETE CASCADE,
    -- Cztery rodzaje bytów wymienione wprost przez kontrakt `queue.link`.
    rodzaj        TEXT    NOT NULL
                          CHECK(rodzaj IN ('okno','ekspert','projekt','automatyka')),
    -- Identyfikator bytu w kształcie, w jakim niesie go kontrakt.
    byt           TEXT    NOT NULL,
    utworzono     TEXT    NOT NULL DEFAULT (strftime('%Y-%m-%dT%H:%M:%fZ','now')),
    -- Powiązanie powtórzone nie jest drugim faktem — jest tym samym faktem.
    UNIQUE(kolejka_id, rodzaj, byt)
);

CREATE INDEX idx_powiazanie_kolejki_kolejka ON powiazanie_kolejki(kolejka_id, rodzaj);
CREATE INDEX idx_powiazanie_kolejki_byt ON powiazanie_kolejki(rodzaj, byt);
