-- Migracja 370 — pętla wykonawcza Studia i praca dwóch agentów naraz.
--
-- Oba narzędzia są DOMYŚLNIE WYŁĄCZONE i włączane jawnym, odwracalnym
-- ustawieniem Operatora. Nastaw nie ma w tym schemacie z zamysłu: idą zasięgami
-- rodziny `config.*`, która ma własny magazyn. Drugiego magazynu ustawień Studio
-- nie zakłada — inaczej Operator wyłączyłby pętlę w jednym miejscu, a ona
-- chodziłaby dalej wedle drugiego.
--
-- ── Dlaczego rozkład zlecenia jest bytem trwałym ────────────────────────────
-- Zlecenie dokumentowe rozłożone na zadania musi przetrwać przeładowanie rdzenia
-- i zamknięcie okna: pętla wykonawcza chodzi obiegami, a obieg drugi ma wiedzieć,
-- co zrobił obieg pierwszy. Rozkład trzymany w pamięci procesu znaczyłby, że
-- każde przeładowanie zaczyna pracę od nowa — na dokumencie, który już jest w pół
-- przerobiony.
--
-- ── Dlaczego zadanie niesie wykonawcę, a nie tylko stan ─────────────────────
-- Przy dwóch agentach naraz „zadanie w biegu" bez wskazania wykonawcy nie mówi
-- niczego: nie wiadomo, czy stoi, czy ktoś nad nim pracuje, ani kogo zapytać
-- o wynik. Tożsamością jest kod agenta z modułu Agents, tak samo jak przy zmianie
-- śledzonej.
--
-- ── Dlaczego zajęcie fragmentu wygasa ────────────────────────────────────────
-- Wykonawca ubity w pół pracy nie może trzymać akapitu na zawsze — drugi agent
-- stałby bezczynnie, a Operator nie wiedziałby, dlaczego. Stąd `wygasa`:
-- zajęcie bez odnowienia przestaje obowiązywać samo. To nie jest rozjemca; to
-- zapora przed zakleszczeniem.
--
-- ── Dlaczego spięcie ma wiersz, choć rdzeń nie rozstrzyga sporu ─────────────
-- Dwóch wykonawców zmieniających ten sam akapit nie może dać dokumentu, w którym
-- jeden nadpisał drugiego W CISZY. Rdzeń nie orzeka, kto ma rację — odkłada
-- prawdę o tym, co się stało: czyja zmiana weszła, czyja została odłożona
-- i dlaczego. Odłożone brzmienie zostaje w wierszu, więc nie przepada: Operator
-- może je wnieść sam albo znaleźć jako propozycję na marginesie.

CREATE TABLE rozklad_zlecenia_studio (
    id                       INTEGER PRIMARY KEY AUTOINCREMENT,
    identyfikator_zewnetrzny TEXT    NOT NULL UNIQUE,
    dokument_id              INTEGER NOT NULL REFERENCES dokument_studio(id) ON DELETE CASCADE,
    zlecenie                 TEXT    NOT NULL,
    stan                     TEXT    NOT NULL DEFAULT 'draft'
                                     CHECK(stan IN ('draft','running','paused','done','stopped')),
    zalozyl_rodzaj           TEXT    NOT NULL DEFAULT 'uzytkownik'
                                     CHECK(zalozyl_rodzaj IN ('uzytkownik','model')),
    zalozyl_agent_kod        TEXT,
    zalozyl_agent_nazwa      TEXT,
    obiegi                   INTEGER NOT NULL DEFAULT 0,
    obiegi_bez_postepu       INTEGER NOT NULL DEFAULT 0,
    powod_zatrzymania        TEXT,
    utworzono                TEXT    NOT NULL DEFAULT (strftime('%Y-%m-%dT%H:%M:%fZ','now')),
    zaktualizowano           TEXT    NOT NULL DEFAULT (strftime('%Y-%m-%dT%H:%M:%fZ','now'))
);
CREATE INDEX idx_rozklad_zlecenia_studio_dokument
    ON rozklad_zlecenia_studio(dokument_id, utworzono DESC, id DESC);

CREATE TABLE zadanie_rozkladu_studio (
    id                       INTEGER PRIMARY KEY AUTOINCREMENT,
    identyfikator_zewnetrzny TEXT    NOT NULL UNIQUE,
    rozklad_id               INTEGER NOT NULL
                                     REFERENCES rozklad_zlecenia_studio(id) ON DELETE CASCADE,
    kolejnosc                INTEGER NOT NULL DEFAULT 0,
    rodzaj                   TEXT    NOT NULL DEFAULT 'custom'
                                     CHECK(rodzaj IN ('research','draft','format','apparatus','review',
                                           'proofread','export','custom')),
    tytul                    TEXT    NOT NULL,
    polecenie                TEXT,
    zakres_od                INTEGER,
    zakres_do                INTEGER,
    stan                     TEXT    NOT NULL DEFAULT 'pending'
                                     CHECK(stan IN ('pending','running','done','failed','blocked','skipped')),
    wykonawca_rodzaj         TEXT
                                     CHECK(wykonawca_rodzaj IS NULL
                                           OR wykonawca_rodzaj IN ('uzytkownik','model')),
    wykonawca_agent_kod      TEXT,
    wykonawca_agent_nazwa    TEXT,
    wykonawca_podagent_kod   TEXT,
    wynik                    TEXT,
    powod_niepowodzenia      TEXT,
    rozpoczeto               TEXT,
    zakonczono               TEXT,
    utworzono                TEXT    NOT NULL DEFAULT (strftime('%Y-%m-%dT%H:%M:%fZ','now')),
    zaktualizowano           TEXT    NOT NULL DEFAULT (strftime('%Y-%m-%dT%H:%M:%fZ','now'))
);
CREATE INDEX idx_zadanie_rozkladu_studio_rozklad
    ON zadanie_rozkladu_studio(rozklad_id, kolejnosc, id);
CREATE INDEX idx_zadanie_rozkladu_studio_stan
    ON zadanie_rozkladu_studio(rozklad_id, stan);

-- Zależność między zadaniami. Osobna tabela z tego samego powodu, co przy
-- dzienniku czynności: pytanie „co stoi na tym zadaniu" pada przy każdym obiegu
-- pętli i nie może wymagać przeszukania wszystkich wierszy.
CREATE TABLE zaleznosc_zadania_studio (
    id                       INTEGER PRIMARY KEY AUTOINCREMENT,
    zadanie_id               INTEGER NOT NULL
                                     REFERENCES zadanie_rozkladu_studio(id) ON DELETE CASCADE,
    podstawa_id              INTEGER NOT NULL
                                     REFERENCES zadanie_rozkladu_studio(id) ON DELETE CASCADE,
    UNIQUE(zadanie_id, podstawa_id),
    CHECK(zadanie_id <> podstawa_id)
);
CREATE INDEX idx_zaleznosc_zadania_studio_podstawa
    ON zaleznosc_zadania_studio(podstawa_id, zadanie_id);

CREATE TABLE zajecie_fragmentu_studio (
    id                       INTEGER PRIMARY KEY AUTOINCREMENT,
    identyfikator_zewnetrzny TEXT    NOT NULL UNIQUE,
    dokument_id              INTEGER NOT NULL REFERENCES dokument_studio(id) ON DELETE CASCADE,
    wykonawca_rodzaj         TEXT    NOT NULL DEFAULT 'model'
                                     CHECK(wykonawca_rodzaj IN ('uzytkownik','model')),
    wykonawca_agent_kod      TEXT,
    wykonawca_agent_nazwa    TEXT,
    wykonawca_podagent_kod   TEXT,
    stan                     TEXT    NOT NULL DEFAULT 'working'
                                     CHECK(stan IN ('idle','working','waiting','conflicted','stopped')),
    zakres_od                INTEGER NOT NULL,
    zakres_do                INTEGER NOT NULL,
    zadanie_kod              TEXT,
    zajeto                   TEXT    NOT NULL DEFAULT (strftime('%Y-%m-%dT%H:%M:%fZ','now')),
    wygasa                   TEXT,
    CHECK(zakres_do >= zakres_od)
);
-- Najczęstsze pytanie tej tabeli: „czy zakres tego żądania dotyka fragmentu
-- zajętego przez kogo innego" — pada przed każdą zmianą dokumentu przy czynnej
-- pracy wielu agentów, więc musi być tanie.
CREATE INDEX idx_zajecie_fragmentu_studio_zakres
    ON zajecie_fragmentu_studio(dokument_id, zakres_od, zakres_do);
CREATE INDEX idx_zajecie_fragmentu_studio_wykonawca
    ON zajecie_fragmentu_studio(dokument_id, wykonawca_agent_kod, stan);

CREATE TABLE spiecie_wykonawcow_studio (
    id                       INTEGER PRIMARY KEY AUTOINCREMENT,
    identyfikator_zewnetrzny TEXT    NOT NULL UNIQUE,
    dokument_id              INTEGER NOT NULL REFERENCES dokument_studio(id) ON DELETE CASCADE,
    zakres_od                INTEGER NOT NULL,
    zakres_do                INTEGER NOT NULL,
    weszla_rodzaj            TEXT    NOT NULL DEFAULT 'model'
                                     CHECK(weszla_rodzaj IN ('uzytkownik','model')),
    weszla_agent_kod         TEXT,
    weszla_agent_nazwa      TEXT,
    odlozona_rodzaj          TEXT    NOT NULL DEFAULT 'model'
                                     CHECK(odlozona_rodzaj IN ('uzytkownik','model')),
    odlozona_agent_kod       TEXT,
    odlozona_agent_nazwa     TEXT,
    nastawa                  TEXT    NOT NULL DEFAULT 'queue'
                                     CHECK(nastawa IN ('refuse','queue','fragmentLock')),
    powod                    TEXT    NOT NULL,
    -- Odłożone brzmienie NIE PRZEPADA. Bez tej kolumny bilans mówiłby „zmiana
    -- drugiego wykonawcy została odłożona" i nie byłoby czego wnieść.
    brzmienie_odlozone       TEXT,
    zmiana_sledzona_kod      TEXT,
    znakowanie_kod           TEXT,
    domkniete                INTEGER NOT NULL DEFAULT 0,
    utworzono                TEXT    NOT NULL DEFAULT (strftime('%Y-%m-%dT%H:%M:%fZ','now')),
    CHECK(zakres_do >= zakres_od)
);
CREATE INDEX idx_spiecie_wykonawcow_studio_dokument
    ON spiecie_wykonawcow_studio(dokument_id, domkniete, utworzono DESC, id DESC);
