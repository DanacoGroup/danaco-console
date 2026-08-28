-- Migracja tworzy okno komunikacji jako byt pośredniczący między sesją a wiadomością;
-- niesie moduł, kanał modelu, katalogi robocze, środowisko wykonania, tryb uprawnień
-- i rolę w pętli koordynator-wykonawca. Wartości wyliczeniowe pochodzą z kontraktu.

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
    -- Pole odpowiada coordinatorWindowId kontraktu: okno wykonawcy wskazuje własnego koordynatora.
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

-- Katalog roboczy okna jest listą ścieżek, nie pojedynczym polem, ponieważ jedno
-- okno komunikacji może pracować na wielu katalogach jednocześnie.
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

-- Wiadomość należy do okna, nie wprost do sesji; rola ma cztery wartości kontraktu,
-- a okno_zrodlowe_id niesie atrybucję w pętli koordynator-wykonawca. Treść obszerna
-- trafia do pliku, baza trzyma odwołanie.
CREATE TABLE wiadomosc (
    id                     INTEGER PRIMARY KEY AUTOINCREMENT,
    okno_komunikacji_id    INTEGER NOT NULL REFERENCES okno_komunikacji(id) ON DELETE CASCADE,
    rola                   TEXT    NOT NULL
                                   CHECK(rola IN ('uzytkownik','model','system','narzedzie')),
    -- Rola nazywa nadawcę, persona nazywa pełnioną funkcję; jest pusta dla zwykłej wypowiedzi.
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
