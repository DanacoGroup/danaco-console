-- Migracja 077 — obsada biegu: jeden komplet bytów dla obu kombajnów wykonawczych.
--
-- Obsada 1-4 uczestników z rolami służy modułowi Automations (migracja 074),
-- a MultitaskingAI potrzebuje dokładnie tego samego bytu: udział bierze od 1 do 4
-- modeli z podziałem na role. Drugi komplet tabel dałby dwa słowniki ról, dwa
-- sufity i dwa silniki wybudzeń — dwie prawdy o jednym mechanizmie. Przeszkodą
-- było `obsada_automatyki.automatyka_id NOT NULL`: bieg orkiestracji sesyjnej
-- zakłada się w oknie i nie ma wiersza w `automatyka`. Ta migracja rozluźnia
-- nośnik: obsada wisi albo na automatyce, albo na biegu orkiestracji, nigdy na
-- obu i nigdy na żadnym — tym samym wzorcem, którym silnik kolejek jest wspólny
-- dla pętli sesyjnej i MultitaskingAI (`adapter_modul_automations_kolejka.go`).
--
-- Górna granica czterech miejsc siedzi w więzie (miejsce BETWEEN 1 AND 4):
-- koordynator, dwaj wykonawcy i analityk. Nie jest to sprzeczne z górną granicą
-- piętnastu podagentów (`subagent.spawn`): obsada liczy stanowiska w scenie,
-- a podagenci pracują pod jednym wykonawcą, w tle, bez własnego gniazda — dwa
-- różne sufity dwóch różnych bytów.
--
-- Ponad obsadę modułu Automations każde stanowisko niesie własny prompt, własne
-- narzędzia i własny profil izolacji (migracja 060), inaczej „cztery modele
-- z podziałem na role" znaczyłoby tylko cztery nazwy zamiast czterech różnie
-- wyposażonych stanowisk. Obsada niczego nie zabrania: bieg bez obsady rusza na
-- modelu wskazanym w kroku albo w oknie, a pusty wykaz obsady to krótsza lista,
-- nie odmowa.

-- ── 1. Bieg orkiestracji — przepływ MultitaskingAI ──────────────────────────
--
-- Bieg orkiestracji to nazwany bieg prowadzony w sesji, grupujący agentów
-- pracujących pod orkiestratorem. Ma stan, etap i kolejkę, którą jedzie — bo
-- silnik jest jeden. Nie używa `przebieg_automatyki`, bo tamten wisi na
-- `automatyka_id NOT NULL`, czyli na zapisanej definicji; bieg orkiestracji
-- zakłada się w oknie i bywa jednorazowy, a zmuszanie Operatora do zapisania
-- automatyki, żeby móc zestawić dwa modele, byłoby bramką. Stan `oczekuje` jest
-- tu wartością pierwszej klasy: bieg zawieszony na sygnał ze świata to stan tego
-- bytu, nie dopisek.

CREATE TABLE bieg_orkiestracji (
    id                       INTEGER PRIMARY KEY AUTOINCREMENT,
    identyfikator_zewnetrzny TEXT    NOT NULL UNIQUE,
    sesja_id                 INTEGER NOT NULL REFERENCES sesja(id) ON DELETE CASCADE,
    -- Okno prowadzące bieg. Puste znaczy bieg bez koordynatora — dopuszczalne,
    -- bo jeden model w obsadzie nie ma kogo koordynować.
    okno_koordynatora_id     INTEGER REFERENCES okno_komunikacji(id) ON DELETE SET NULL,
    -- Kolejka niosąca pracę. Wiązanie luźne: bieg powstaje przed kolejką,
    -- a skasowanie kolejki nie ma prawa skasować śladu biegu.
    kolejka_id               INTEGER REFERENCES kolejka(id) ON DELETE SET NULL,
    nazwa                    TEXT    NOT NULL,
    stan                     TEXT    NOT NULL DEFAULT 'pending'
                                     CHECK(stan IN ('pending','running','paused','oczekuje',
                                                    'succeeded','failed','stopped')),
    -- Nazwa etapu widoczna w panelu („Wdrożenie"). Rdzeń zna etap, nie zna
    -- nazwy procesu — dlatego etap jest polem, a nie wyliczeniem.
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

-- ── 2. Przebudowa obsady — dwa nośniki, jeden byt ───────────────────────────
--
-- SQLite nie zdejmuje NOT NULL ani CHECK, więc tabela idzie przez przebudowę.
-- Na `obsada_automatyki` nie wskazuje żaden klucz obcy, więc przebudowa obejmuje
-- jedną tabelę i kaskada nie ma czego zabrać.
--
-- Więz „dokładnie jeden nośnik" jest tu sednem poprawności: wiersz bez nośnika
-- to obsada niczyja, wiersz z dwoma — obsada dwóch biegów naraz. SQLite liczy
-- wyrażenia logiczne jako 0/1, więc suma dwóch testów równa 1 wyraża „dokładnie
-- jeden" bez wyzwalacza.
--
-- Jednoznaczność miejsca idzie dwoma indeksami częściowymi, nie jednym UNIQUE:
-- `UNIQUE(automatyka_id, miejsce)` przepuściłoby dwa wiersze o tym samym miejscu
-- dla różnych biegów orkiestracji, bo NULL nie równa się NULL. Dwa indeksy
-- częściowe pilnują każdego nośnika osobno i dokładnie.

CREATE TABLE obsada_biegu (
    id                 INTEGER PRIMARY KEY AUTOINCREMENT,
    automatyka_id      INTEGER REFERENCES automatyka(id) ON DELETE CASCADE,
    bieg_id            INTEGER REFERENCES bieg_orkiestracji(id) ON DELETE CASCADE,

    -- Uczestnik: agent z portfolio albo goły model. Agent ma pierwszeństwo —
    -- niesie instrukcje systemowe, umiejętności i parametry.
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

-- ── 3. Podagent — zadanie w tle pod przepływem ──────────────────────────────
--
-- Trwałą tożsamością podagenta jest pozycja kolejki, a proces modelu jest
-- wyłącznie sposobem jej wykonania — kontrakt podpowiada to polem
-- `Subagent.queueItemId` („Pozycja kolejki niosąca pracę podagenta"). Podagent
-- nie jest oknem, bo agentów bywa więcej niż gniazd sceny; nie jest samym
-- procesem, bo wykaz zakończonych przeżywa restart rdzenia, a proces nie; jego
-- pracę wykonuje silnik kolejek, który już jest.
--
-- Ponad pozycję kolejki ten wiersz niesie to, czego pozycja nie wie: kto ją
-- wykonuje w rozumieniu obsady, pod jakim biegiem biegnie i ile kosztowała.
-- Żetony, narzędzia i czas są polami, bo pokazuje je panel zadań w tle;
-- kontraktowy `Subagent` ich nie niesie, więc zapisujemy je w bazie.

CREATE TABLE podagent (
    id                       INTEGER PRIMARY KEY AUTOINCREMENT,
    identyfikator_zewnetrzny TEXT    NOT NULL UNIQUE,
    -- Okno wykonawcy, które podagenta powołało. To ono jest jego orkiestratorem.
    okno_wykonawcy_id        INTEGER NOT NULL REFERENCES okno_komunikacji(id) ON DELETE CASCADE,
    -- Bieg, w którym podagent pracuje. Puste znaczy podagenta powołanego poza
    -- biegiem orkiestracji — wykonawca może go uruchomić w zwykłej rozmowie.
    bieg_id                  INTEGER REFERENCES bieg_orkiestracji(id) ON DELETE SET NULL,
    -- Pozycja kolejki niosąca pracę. Puste wyłącznie między powołaniem
    -- a zasileniem kolejki; potem wskazuje zawsze.
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

-- Wykaz podagentów okna wykonawcy — `subagent.list` pyta dokładnie o to.
CREATE INDEX idx_podagent_okno
    ON podagent(okno_wykonawcy_id, id);

-- Wykaz podagentów przepływu — wiersze panelu zadań w tle.
CREATE INDEX idx_podagent_bieg
    ON podagent(bieg_id, id)
    WHERE bieg_id IS NOT NULL;

-- Droga powrotna od pozycji kolejki do podagenta: silnik kolejki wie o pozycji,
-- panel pyta o podagenta.
CREATE UNIQUE INDEX idx_podagent_pozycja
    ON podagent(pozycja_kolejki_id)
    WHERE pozycja_kolejki_id IS NOT NULL;

-- Zbieranie wyników (`subagent.result.collect`) pyta o niezakończonych.
CREATE INDEX idx_podagent_stan
    ON podagent(stan, id);
