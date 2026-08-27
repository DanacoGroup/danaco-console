-- Migracja 077 ujednolica obsadę biegu: jeden komplet bytów obsługuje zarówno
-- automatyzacje, jak i przepływ wieloagentowy, zamiast dwóch osobnych
-- słowników ról.

-- Bieg orkiestracji to nazwany bieg prowadzony w sesji, grupujący agentów pod
-- orkiestratorem; korzysta z tego samego silnika kolejki co automatyzacje,
-- zamiast osobnego mechanizmu.
CREATE TABLE bieg_orkiestracji (
    id                       INTEGER PRIMARY KEY AUTOINCREMENT,
    identyfikator_zewnetrzny TEXT    NOT NULL UNIQUE,
    sesja_id                 INTEGER NOT NULL REFERENCES sesja(id) ON DELETE CASCADE,
    -- Okno prowadzące bieg; wartość pusta oznacza bieg bez koordynatora,
    -- dopuszczalny przy jednym modelu.
    okno_koordynatora_id     INTEGER REFERENCES okno_komunikacji(id) ON DELETE SET NULL,
    -- Kolejka niosąca pracę; wiązanie luźne, bo skasowanie kolejki nie może
    -- skasować śladu biegu.
    kolejka_id               INTEGER REFERENCES kolejka(id) ON DELETE SET NULL,
    nazwa                    TEXT    NOT NULL,
    stan                     TEXT    NOT NULL DEFAULT 'pending'
                                     CHECK(stan IN ('pending','running','paused','oczekuje',
                                                    'succeeded','failed','stopped')),
    -- Nazwa etapu widoczna w panelu; rdzeń zna etap, nie zna nazwy procesu,
    -- dlatego jest to zwykłe pole.
    etap                     TEXT,
    etap_biezacy             INTEGER NOT NULL DEFAULT 0,
    etapow                   INTEGER NOT NULL DEFAULT 0,
    rozpoczeto               TEXT    NOT NULL DEFAULT (strftime('%Y-%m-%dT%H:%M:%fZ','now')),
    zakonczono               TEXT
);

CREATE INDEX idx_bieg_orkiestracji_sesja
    ON bieg_orkiestracji(sesja_id, id DESC);

CREATE UNIQUE INDEX idx_bieg_orkiestracji_kolejka
    ON bieg_orkiestracji(kolejka_id)
    WHERE kolejka_id IS NOT NULL;

-- Obsada łączy dwa nośniki w jednym bycie: wiersz wisi albo na automatyce,
-- albo na biegu orkiestracji, nigdy na obu ani na żadnym, a jednoznaczność
-- miejsca pilnują dwa indeksy częściowe.
CREATE TABLE obsada_biegu (
    id                 INTEGER PRIMARY KEY AUTOINCREMENT,
    automatyka_id      INTEGER REFERENCES automatyka(id) ON DELETE CASCADE,
    bieg_id            INTEGER REFERENCES bieg_orkiestracji(id) ON DELETE CASCADE,

    -- Uczestnik: agent z portfolio albo goły model; agent ma pierwszeństwo,
    -- bo niesie własne wyposażenie.
    agent_id           INTEGER REFERENCES agent(id) ON DELETE SET NULL,
    model              TEXT,
    rola               TEXT    NOT NULL DEFAULT 'wykonawca'
                               CHECK(rola IN ('koordynator','wykonawca','samodzielne')),
    miejsce            INTEGER NOT NULL CHECK(miejsce BETWEEN 1 AND 4),

    -- Wyposażenie stanowiska. Puste znaczy „bierz z agenta albo
    -- z ustawień biegu" — nie jest brakiem.
    prompt             TEXT,
    narzedzia          TEXT,     -- wykaz nazw narzędzi zapisem JSON
    profil_izolacji_id INTEGER REFERENCES profil_izolacji(id) ON DELETE SET NULL,

    utworzono          TEXT    NOT NULL DEFAULT (strftime('%Y-%m-%dT%H:%M:%fZ','now')),
    zaktualizowano     TEXT    NOT NULL DEFAULT (strftime('%Y-%m-%dT%H:%M:%fZ','now')),

    CHECK ((automatyka_id IS NOT NULL) + (bieg_id IS NOT NULL) = 1),
    CHECK (agent_id IS NOT NULL OR model IS NOT NULL)
);

-- Przepisanie obsady zastanej. Wszystkie wiersze migracji 074 wiszą na
-- automatyce, więc `bieg_id` zostaje pusty, a wyposażenie — nieustawione.
INSERT INTO obsada_biegu
    (id, automatyka_id, agent_id, model, rola, miejsce, utworzono, zaktualizowano)
SELECT
     id, automatyka_id, agent_id, model, rola, miejsce, utworzono, zaktualizowano
  FROM obsada_automatyki;

DROP TABLE obsada_automatyki;

CREATE UNIQUE INDEX idx_obsada_biegu_miejsce_automatyki
    ON obsada_biegu(automatyka_id, miejsce)
    WHERE automatyka_id IS NOT NULL;

CREATE UNIQUE INDEX idx_obsada_biegu_miejsce_orkiestracji
    ON obsada_biegu(bieg_id, miejsce)
    WHERE bieg_id IS NOT NULL;

CREATE INDEX idx_obsada_biegu_agent
    ON obsada_biegu(agent_id)
    WHERE agent_id IS NOT NULL;

CREATE INDEX idx_obsada_biegu_profil
    ON obsada_biegu(profil_izolacji_id)
    WHERE profil_izolacji_id IS NOT NULL;

-- Podagent to zadanie w tle wykonywane przez silnik kolejki; trwałą
-- tożsamością jest pozycja kolejki, a wiersz niesie to, czego sama pozycja
-- nie wie — obsadę, bieg i koszt.
CREATE TABLE podagent (
    id                       INTEGER PRIMARY KEY AUTOINCREMENT,
    identyfikator_zewnetrzny TEXT    NOT NULL UNIQUE,
    -- Okno wykonawcy, które podagenta powołało. To ono jest jego orkiestratorem.
    okno_wykonawcy_id        INTEGER NOT NULL REFERENCES okno_komunikacji(id) ON DELETE CASCADE,
    -- Bieg, w którym podagent pracuje; puste znaczy podagenta powołanego poza
    -- biegiem orkiestracji.
    bieg_id                  INTEGER REFERENCES bieg_orkiestracji(id) ON DELETE SET NULL,
    -- Pozycja kolejki niosąca pracę; pusta wyłącznie między powołaniem
    -- a zasileniem kolejki.
    pozycja_kolejki_id       INTEGER REFERENCES pozycja_kolejki(id) ON DELETE SET NULL,
    nazwa                    TEXT,
    zadanie                  TEXT    NOT NULL,
    stan                     TEXT    NOT NULL DEFAULT 'pending'
                                     CHECK(stan IN ('pending','running','done','failed','stopped')),
    wynik                    TEXT,
    -- Rachunek widoczny w panelu zadań w tle.
    zetonow                  INTEGER NOT NULL DEFAULT 0,
    narzedzi                 INTEGER NOT NULL DEFAULT 0,
    rozpoczeto               TEXT,
    zakonczono               TEXT,
    utworzono                TEXT    NOT NULL DEFAULT (strftime('%Y-%m-%dT%H:%M:%fZ','now'))
);

-- Wykaz podagentów okna wykonawcy służy odpytywaniu o listę podagentów,
-- jakie to konkretne okno powołało.
CREATE INDEX idx_podagent_okno
    ON podagent(okno_wykonawcy_id, id);

-- Wykaz podagentów przepływu tworzy wiersze, które pokazuje operatorowi
-- panel zadań pracujących w tle.
CREATE INDEX idx_podagent_bieg
    ON podagent(bieg_id, id)
    WHERE bieg_id IS NOT NULL;

-- Droga powrotna od pozycji kolejki do podagenta: silnik kolejki wie o pozycji,
-- panel pyta o podagenta.
CREATE UNIQUE INDEX idx_podagent_pozycja
    ON podagent(pozycja_kolejki_id)
    WHERE pozycja_kolejki_id IS NOT NULL;

-- Zbieranie wyników korzysta z tego indeksu przy odpytywaniu o podagentów,
-- których praca się nie zakończyła.
CREATE INDEX idx_podagent_stan
    ON podagent(stan, id);
