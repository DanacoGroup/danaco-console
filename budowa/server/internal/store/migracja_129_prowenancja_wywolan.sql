-- Migracja 129 — ślad wywołania kanału modelu wraz z drzewem odcinków.
--
-- Podsystem prowenancji: rdzeń zapisuje, co poszło do modelu i co wróciło.
-- Bez tego zapisu Provenance Explorer modułu Diagnostics nie ma czego pokazać,
-- a rozliczenie zużycia nie ma po czym liczyć — dlatego obie rodziny komend
-- stoją na tej jednej tabeli, a nie na dwóch osobnych. Dwie tabele o tych samych
-- liczbach rozjechałyby się przy pierwszej korekcie cennika.
--
-- ── Dlaczego treść leży w tej samej tabeli, a nie obok ────────────────────────
-- Prompt i odpowiedź to kolumny osobne, ale wiersza tego samego. Wyłączenie
-- zapisu treści ustawieniem zostawia wiersz z liczbami i przebiegiem, a kolumny
-- treści puste — ślad wywołania zostaje mierzalny nawet wtedy, gdy jego treść
-- świadomie nie jest przechowywana. Kolumna `tresc_zapisana` odróżnia treść
-- pustą od treści niezapisanej: bez niej wiersz bez promptu znaczyłby dwie różne
-- rzeczy naraz.
--
-- ── Dlaczego drzewo jest listą krawędzi ──────────────────────────────────────
-- Odcinek wywołania (wejście modelu, użycie narzędzia, tura podagenta) niesie
-- wskazanie rodzica. Zapytanie o jeden odcinek czyta wtedy jeden wiersz, a nie
-- rozbiera cały dokument JSON, żeby dojść do gałęzi.

CREATE TABLE prowenancja_wywolanie (
    id                       INTEGER PRIMARY KEY AUTOINCREMENT,
    kod                      TEXT    NOT NULL UNIQUE,
    -- slad_kod spina wywołania jednego zlecenia Operatora; rodzic_kod wiąże
    -- wywołanie podagenta z wywołaniem, które je zleciło.
    slad_kod                 TEXT,
    rodzic_kod               TEXT,
    proces_kod               TEXT,
    sesja_kod                TEXT,
    okno_kod                 TEXT,
    wiadomosc_kod            TEXT,
    kanal_kod                TEXT,
    dostawca                 TEXT,
    model                    TEXT,
    konto_kod                TEXT,
    -- projekt_kod i srodowisko istnieją po to, żeby rozliczenie zużycia miało
    -- po czym grupować. Bez nich wymiary „projekt" i „środowisko" byłyby
    -- wymiarami bez danych.
    projekt_kod              TEXT,
    srodowisko               TEXT,
    -- Wartości kolumn biorą się z odwzorowania w kontrakcie (pole `baza` przy
    -- każdej wartości wyliczenia), a nie z wyobrażenia migracji. Kontrakt jest
    -- źródłem prawdy nazw i to on rozstrzyga, co kolumna przyjmuje.
    stan                     TEXT    NOT NULL DEFAULT 'biegnie'
                                     CHECK(stan IN ('biegnie','zakonczone','bledne',
                                                    'limit_czasu','przerwane')),
    kod_bledu                TEXT,
    poczatek                 INTEGER NOT NULL,
    koniec                   INTEGER,
    opoznienie_ms            INTEGER,
    tokeny_promptu           INTEGER,
    tokeny_odpowiedzi        INTEGER,
    tokeny_cache             INTEGER,
    tokeny_razem             INTEGER,
    koszt                    REAL,
    waluta                   TEXT,
    liczba_narzedzi          INTEGER NOT NULL DEFAULT 0,
    liczba_podagentow        INTEGER NOT NULL DEFAULT 0,
    -- ocena jest zdaniem Operatora o wywołaniu, nie pomiarem rdzenia.
    ocena                    TEXT    NOT NULL DEFAULT 'bez_oceny'
                                     CHECK(ocena IN ('bez_oceny','trafna','czesciowa','nietrafna')),
    ocena_notatka            TEXT,
    tresc_zapisana           INTEGER NOT NULL DEFAULT 0 CHECK(tresc_zapisana IN (0,1)),
    zredagowane              INTEGER NOT NULL DEFAULT 0 CHECK(zredagowane IN (0,1)),
    prompt                   TEXT,
    odpowiedz                TEXT
);

CREATE INDEX idx_prowenancja_wywolanie_poczatek ON prowenancja_wywolanie(poczatek DESC);
CREATE INDEX idx_prowenancja_wywolanie_slad     ON prowenancja_wywolanie(slad_kod);
CREATE INDEX idx_prowenancja_wywolanie_proces   ON prowenancja_wywolanie(proces_kod, poczatek DESC);
CREATE INDEX idx_prowenancja_wywolanie_kanal    ON prowenancja_wywolanie(kanal_kod, poczatek DESC);
CREATE INDEX idx_prowenancja_wywolanie_sesja    ON prowenancja_wywolanie(sesja_kod, poczatek DESC);

CREATE TABLE prowenancja_odcinek (
    id             INTEGER PRIMARY KEY AUTOINCREMENT,
    kod            TEXT    NOT NULL UNIQUE,
    wywolanie_kod  TEXT    NOT NULL,
    rodzic_kod     TEXT,
    nazwa          TEXT    NOT NULL,
    rodzaj         TEXT    NOT NULL
                           CHECK(rodzaj IN ('prompt','odpowiedz','narzedzie',
                                            'podagent','kontekst','cache')),
    poczatek       INTEGER NOT NULL,
    koniec         INTEGER,
    czas_ms        INTEGER,
    tokeny         INTEGER,
    koszt          REAL,
    stan           TEXT    NOT NULL DEFAULT 'biegnie'
                           CHECK(stan IN ('biegnie','zakonczone','bledne',
                                          'limit_czasu','przerwane')),
    kod_bledu      TEXT,
    trafienie_cache INTEGER NOT NULL DEFAULT 0 CHECK(trafienie_cache IN (0,1)),
    atrybuty       TEXT
);

CREATE INDEX idx_prowenancja_odcinek_wywolanie ON prowenancja_odcinek(wywolanie_kod, poczatek);
CREATE INDEX idx_prowenancja_odcinek_rodzic    ON prowenancja_odcinek(rodzic_kod);

-- ── Korelacja dziennika ze śladem ────────────────────────────────────────────
-- Wpis dziennika zdarzeń zyskuje wskazanie śladu, więc Provenance Explorer
-- pokazuje obok siebie wywołanie i to, co rdzeń o nim zapisał.
ALTER TABLE diagnostyka_wpis ADD COLUMN slad_kod TEXT;

-- Kolumna procesu stoi w dzienniku od migracji zakładającej diagnostykę, ale bez
-- indeksu: zawężenie po procesie skanowało całą tabelę.
CREATE INDEX idx_diagnostyka_wpis_proces  ON diagnostyka_wpis(proces_kod, chwila DESC);
CREATE INDEX idx_diagnostyka_wpis_slad    ON diagnostyka_wpis(slad_kod);
CREATE INDEX idx_diagnostyka_blad_zrodlo  ON diagnostyka_blad(zrodlo, ostatnie DESC);
