-- Tworzy tabelę konsultacja_doradcy, przechowującą pytanie, radę, dane kanału i modelu doradcy z chwili pytania oraz kopię fragmentu prowenancji strumienia, osobno od wiadomości agenta.

CREATE TABLE konsultacja_doradcy (
    id                       INTEGER PRIMARY KEY AUTOINCREMENT,
    -- Tożsamość wpisu widoczna poza bazą; klucz sztuczny zostaje sprawą bazy.
    identyfikator_zewnetrzny TEXT    NOT NULL UNIQUE,
    -- Okno, w którym padła konsultacja. NULL znaczy „spoza okna".
    okno_id                  TEXT,
    -- Kanał agenta pytającego — po nim rozstrzygał się próg siły doradcy.
    pytajacy_kanal           TEXT    NOT NULL,
    -- Kanał i model doradcy zapisane w chwili pytania, niezależnie od późniejszej zmiany kanału.
    doradca_kanal            TEXT    NOT NULL,
    doradca_model            TEXT    NOT NULL,
    -- Pełna treść wysłana doradcy wraz z ramą konsultacji, dowodząca, że doradcy kazano radzić.
    pytanie                  TEXT    NOT NULL,
    -- Treść rady. NULL przy odmowie — nie ma czego zapisać.
    rada                     TEXT,
    -- Skrót z pytania i rady; diagnostyka, nigdy bramka.
    skrot                    TEXT,
    -- Fragment `provenance` strumienia, przepisany co do znaku, nie odtworzony powtórnym złożeniem.
    prowenancja              TEXT,
    stan                     TEXT    NOT NULL CHECK(stan IN ('rada', 'odmowa')),
    -- Nazwany powód odmowy; przy radzie nie ma czego opisywać.
    powod                    TEXT,
    -- Milisekundy epoki podane przez wołającego.
    utworzono                INTEGER NOT NULL,
    CHECK((stan = 'rada'   AND rada  IS NOT NULL AND rada  <> '' AND powod IS NULL)
       OR (stan = 'odmowa' AND powod IS NOT NULL AND powod <> '' AND rada  IS NULL))
);

-- Tworzy indeks konsultacji według okna i czasu utworzenia, obejmujący porządek wyniku, dzięki czemu wykaz konsultacji okna od najnowszej nie wymaga sortowania.
CREATE INDEX idx_konsultacja_doradcy_wykaz
    ON konsultacja_doradcy (okno_id, utworzono, id);

-- Wykaz konsultacji jednego doradcy — odpowiedź na pytanie „ile i o co pytano
-- tego doradcę", które rozstrzyga o koszcie silniejszego modelu.
CREATE INDEX idx_konsultacja_doradcy_doradca
    ON konsultacja_doradcy (doradca_kanal, utworzono);
