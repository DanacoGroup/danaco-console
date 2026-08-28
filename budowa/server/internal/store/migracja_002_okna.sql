-- Migracja 002 — okno komunikacji jako byt pośredni między sesją a wiadomością.
-- Okno niesie moduł, kanał modelu, listę katalogów roboczych, środowisko wykonania,
-- tryb uprawnień oraz rolę w pętli koordynator–wykonawca.
--
-- Wartości wyliczeniowe pochodzą z kontraktu — `shared/contract.json` jest
-- jedynym źródłem prawdy. Kontrakt zapisuje nazwy po angielsku, model danych
-- po polsku; odwzorowanie jest jeden do jednego i leży wyłącznie w kontrakcie
-- (pole `baza` przy każdej wartości wyliczenia). Generator wytwarza z niego słowniki
-- `WartosciBazy*` / `WartosciKontraktu*` w `contract.go` i `contract.ts` — warstwa
-- trwałości sięga po nie zamiast wpisywać przekład u siebie.
--
--   okno_komunikacji.srodowisko_wykonania  ← ExecutionEnv
--   okno_komunikacji.tryb_uprawnien        ← PermissionMode
--   okno_komunikacji.rola_okna             ← WindowRole
--   okno_komunikacji.stan                  ← WindowStatus
--   proces_sesji.stan                      ← ProgressStatus
--   wiadomosc.rola                         ← MessageRole
--   wiadomosc.stan                         ← MessageStatus
--   wiadomosc.rodzaj_tresci                ← ChunkKind

CREATE TABLE okno_komunikacji (
    id                     INTEGER PRIMARY KEY AUTOINCREMENT,
    sesja_id               INTEGER NOT NULL REFERENCES sesja(id) ON DELETE CASCADE,
    modul_id               INTEGER NOT NULL REFERENCES modul(id) ON DELETE RESTRICT,
    kanal_modelu_id        INTEGER NOT NULL REFERENCES kanal_modelu(id) ON DELETE RESTRICT,
    tytul                  TEXT,
    srodowisko_wykonania   TEXT    NOT NULL DEFAULT 'lokalne'
                                   CHECK(srodowisko_wykonania IN ('lokalne','rdzen','zdalny')),
    tryb_uprawnien         TEXT    NOT NULL DEFAULT 'reczny'
                                   CHECK(tryb_uprawnien IN ('reczny','akceptuj_zmiany','plan','automatyczny','bez_zapytan',
                                                            'bez_pytania')),
    rola_okna              TEXT    NOT NULL DEFAULT 'samodzielne'
                                   CHECK(rola_okna IN ('samodzielne','wykonawca','koordynator')),
    -- Odpowiednik pola `coordinatorWindowId` kontraktu: okno wykonawcy wskazuje
    -- swojego koordynatora, dzięki czemu pętla koordynator–wykonawca ma kierunek.
    okno_koordynatora_id   INTEGER REFERENCES okno_komunikacji(id) ON DELETE SET NULL,
    tryb_komunikacji       TEXT    NOT NULL DEFAULT 'tekst',
    stan                   TEXT    NOT NULL DEFAULT 'otwarte'
                                   CHECK(stan IN ('otwarte','zamkniete')),
    kolejnosc              INTEGER NOT NULL DEFAULT 0,
    utworzono              TEXT    NOT NULL DEFAULT (strftime('%Y-%m-%dT%H:%M:%fZ','now')),
    zaktualizowano         TEXT    NOT NULL DEFAULT (strftime('%Y-%m-%dT%H:%M:%fZ','now'))
);
CREATE INDEX idx_okno_sesja ON okno_komunikacji(sesja_id, kolejnosc);
CREATE INDEX idx_okno_modul ON okno_komunikacji(modul_id);
CREATE INDEX idx_okno_kanal ON okno_komunikacji(kanal_modelu_id);
CREATE INDEX idx_okno_koordynator ON okno_komunikacji(okno_koordynatora_id);

-- ── Katalog roboczy okna jest listą, nie pojedynczym polem ───────────
CREATE TABLE katalog_okna (
    id                     INTEGER PRIMARY KEY AUTOINCREMENT,
    okno_komunikacji_id    INTEGER NOT NULL REFERENCES okno_komunikacji(id) ON DELETE CASCADE,
    sciezka                TEXT    NOT NULL,
    kolejnosc              INTEGER NOT NULL DEFAULT 0,
    UNIQUE(okno_komunikacji_id, sciezka)
);
CREATE INDEX idx_katalog_okna_okno ON katalog_okna(okno_komunikacji_id, kolejnosc);

-- ── Proces sesji: 1:N wobec sesji, 1:1 wobec okna; telemetria postępu ─
-- Stan procesu jest słownikiem ProgressStatus zdarzenia progress.changed.
CREATE TABLE proces_sesji (
    id                     INTEGER PRIMARY KEY AUTOINCREMENT,
    okno_komunikacji_id    INTEGER NOT NULL UNIQUE
                                   REFERENCES okno_komunikacji(id) ON DELETE CASCADE,
    identyfikator_procesu  TEXT    NOT NULL,
    identyfikator_systemowy INTEGER,
    stan                   TEXT    NOT NULL DEFAULT 'oczekuje'
                                   CHECK(stan IN ('oczekuje','pracuje','wstrzymany',
                                                  'zatrzymany','zakonczony','bledny')),
    etap_biezacy           INTEGER NOT NULL DEFAULT 0,
    liczba_etapow          INTEGER NOT NULL DEFAULT 0,
    stopien_ukonczenia     INTEGER NOT NULL DEFAULT 0
                                   CHECK(stopien_ukonczenia BETWEEN 0 AND 100),
    opis_etapu             TEXT,
    uruchomiono            TEXT    NOT NULL DEFAULT (strftime('%Y-%m-%dT%H:%M:%fZ','now')),
    zaktualizowano         TEXT    NOT NULL DEFAULT (strftime('%Y-%m-%dT%H:%M:%fZ','now')),
    zakonczono             TEXT
);

-- ── Wiadomość należy do okna, nie wprost do sesji ────────────────────
-- Rola ma cztery wartości kontraktu (MessageRole). Atrybucję w pętli
-- koordynator–wykonawca niesie okno_zrodlowe_id: okno źródłowe zna własną
-- rola_okna, więc nie potrzeba drugiego, równoległego słownika ról.
-- Treści obszerne trafiają do pliku — baza trzyma odwołanie.
CREATE TABLE wiadomosc (
    id                     INTEGER PRIMARY KEY AUTOINCREMENT,
    okno_komunikacji_id    INTEGER NOT NULL REFERENCES okno_komunikacji(id) ON DELETE CASCADE,
    rola                   TEXT    NOT NULL
                                   CHECK(rola IN ('uzytkownik','model','system','narzedzie')),
    -- Atrybucja wypowiedzi. Rola mowi, czym jest nadawca; persona mowi, w jakiej
    -- roli wystapil. Persony sa warstwa prezentacji i nie mnoza wartosci pola
    -- `rola`. Puste dla zwyklej wypowiedzi uzytkownika i modelu.
    persona                TEXT
                                   CHECK(persona IS NULL OR persona IN (
                                       'coordinator','executor','validator',
                                       'automation','aod','tool')),
    okno_zrodlowe_id       INTEGER REFERENCES okno_komunikacji(id) ON DELETE SET NULL,
    rodzaj_tresci          TEXT    NOT NULL DEFAULT 'tekst'
                                   CHECK(rodzaj_tresci IN ('tekst','rozumowanie','wywolanie_narzedzia',
                                                           'wynik_narzedzia','obraz','dzwiek','blad')),
    stan                   TEXT    NOT NULL DEFAULT 'oczekuje'
                                   CHECK(stan IN ('oczekuje','strumien','zakonczona',
                                                  'zatrzymana','bledna')),
    tresc                  TEXT,
    tresc_odwolanie        TEXT,
    tokeny_wejscia         INTEGER NOT NULL DEFAULT 0,
    tokeny_wyjscia         INTEGER NOT NULL DEFAULT 0,
    kolejnosc              INTEGER NOT NULL DEFAULT 0,
    utworzono              TEXT    NOT NULL DEFAULT (strftime('%Y-%m-%dT%H:%M:%fZ','now'))
);
CREATE INDEX idx_wiadomosc_okno ON wiadomosc(okno_komunikacji_id, kolejnosc, id);
CREATE INDEX idx_wiadomosc_okno_zrodlowe ON wiadomosc(okno_zrodlowe_id);
