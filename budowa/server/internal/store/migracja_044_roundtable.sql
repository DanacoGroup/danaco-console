-- Migracja 044 — trwałość modułu Roundtable: uczestnicy debaty, tury,
-- wypowiedzi i stanowisko końcowe.
--
-- Debata jest bytem okna, nie sesji. Kontrakt kieruje wszystkie cztery komendy
-- obszaru `roundtable.*` przez `windowId`, więc identyfikator okna jest tu
-- kluczem grupującym. Kolumna nie ma więzu obcego do `okno`, bo okno debaty
-- bywa oknem operacyjnym rejestru 030 (`roundtable.debate-panel`), a nie
-- oknem komunikacji z migracji 002 — dwa różne byty pod jedną nazwą.
--
-- Ten sam kanał może wystąpić dwukrotnie. Dlatego
-- więzu jednoznaczności na parze (okno, kanał) NIE MA: dwaj uczestnicy stoją na
-- tym samym kanale modelu i różnią się wyłącznie tożsamością — nazwą persony
-- i promptem systemowym. Jednoznaczny jest identyfikator uczestnika, nic więcej.
--
-- Wypowiedź wskazuje uczestnika kodem, nie kluczem. Moderator debaty nie jest
-- uczestnikiem (nie ma kanału ani persony), a jego interwencja jest wypowiedzią
-- tury — więz obcy do `debata_uczestnik` odciąłby ją od zapisu. Kod moderatora
-- jest wartością danych, tak jak kod uczestnika.

-- ── Uczestnik debaty — Model Panels ───────────────────────────────────────────
CREATE TABLE debata_uczestnik (
    id                       INTEGER PRIMARY KEY AUTOINCREMENT,
    identyfikator_zewnetrzny TEXT    NOT NULL UNIQUE,
    okno                     TEXT    NOT NULL,
    kanal_modelu             TEXT    NOT NULL,
    nazwa_tozsamosci         TEXT,
    prompt_systemowy         TEXT,
    -- Wyciszenie obowiązuje w turze: uczestnik zostaje w panelu, lecz pytanie
    -- do niego nie idzie („wycisz w turze”, panel akcji Model Panels).
    wyciszony                INTEGER NOT NULL DEFAULT 0 CHECK(wyciszony IN (0,1)),
    -- Kolumny „kluczowy” tu nie ma, choć kontrakt ma pole
    -- `RoundtableParticipant.key`. Żadna z czterech komend obszaru nie potrafi
    -- go ustawić (`model.add` pola nie przyjmuje, `ModeratorAction` nie ma
    -- wartości oznaczającej), więc kolumna byłaby miejscem, do którego nic nie
    -- pisze.
    -- Kolejność głosu ustawia Moderator Panel; zero znaczy „bez wskazania”,
    -- a wtedy porządkiem jest kolejność dołączenia.
    kolejnosc                INTEGER NOT NULL DEFAULT 0,
    utworzono                TEXT    NOT NULL DEFAULT (strftime('%Y-%m-%dT%H:%M:%fZ','now')),
    zaktualizowano           TEXT    NOT NULL DEFAULT (strftime('%Y-%m-%dT%H:%M:%fZ','now'))
);
CREATE INDEX idx_debata_uczestnik_okno ON debata_uczestnik(okno, kolejnosc, id);

-- ── Tura debaty — Debate Panel ────────────────────────────────────────────────
-- Wartości `format` i `stan` są wartościami kontraktu (RoundtableFormat,
-- RoundtableTurnStatus), nie ich tłumaczeniem — przekład wiersza na turę
-- kontraktu nie potrzebuje wtedy słownika pośredniego.
CREATE TABLE debata_tura (
    id                       INTEGER PRIMARY KEY AUTOINCREMENT,
    identyfikator_zewnetrzny TEXT    NOT NULL UNIQUE,
    okno                     TEXT    NOT NULL,
    numer                    INTEGER NOT NULL,
    zagadnienie              TEXT,
    -- Pytanie tury zostaje w zapisie, bo transkrypt bez pytania jest zbiorem
    -- odpowiedzi na nieznane („Eksportuj transkrypt”, panel akcji Debate Panel).
    pytanie                  TEXT    NOT NULL DEFAULT '',
    format                   TEXT    NOT NULL DEFAULT 'free'
                                     CHECK(format IN ('free','structured','oxford','roundRobin')),
    stan                     TEXT    NOT NULL DEFAULT 'open'
                                     CHECK(stan IN ('open','closed')),
    -- Granica liczby tur ustawiona przy uruchomieniu; zero znaczy „bez granicy”.
    granica_tur              INTEGER NOT NULL DEFAULT 0,
    rozpoczeto               TEXT    NOT NULL DEFAULT (strftime('%Y-%m-%dT%H:%M:%fZ','now')),
    zamknieto                TEXT,
    UNIQUE (okno, numer)
);
CREATE INDEX idx_debata_tura_okno ON debata_tura(okno, numer DESC);

-- ── Wypowiedź uczestnika w turze ──────────────────────────────────────────────
CREATE TABLE debata_wypowiedz (
    id                       INTEGER PRIMARY KEY AUTOINCREMENT,
    identyfikator_zewnetrzny TEXT    NOT NULL UNIQUE,
    tura_id                  INTEGER NOT NULL REFERENCES debata_tura(id) ON DELETE CASCADE,
    uczestnik                TEXT    NOT NULL,
    tresc                    TEXT    NOT NULL DEFAULT '',
    utworzono                TEXT    NOT NULL DEFAULT (strftime('%Y-%m-%dT%H:%M:%fZ','now'))
);
CREATE INDEX idx_debata_wypowiedz_tura ON debata_wypowiedz(tura_id, id);

-- ── Stanowisko końcowe — Consensus Panel ──────────────────────────────────────
-- Pusty kod tury znaczy „stanowisko całej debaty”. Wartość pusta zamiast NULL
-- jest tu wyborem świadomym: w SQLite dwa NULL-e są różne, więc więz
-- UNIQUE(okno, tura) na kolumnie dopuszczającej NULL nie powstrzymałby
-- powielenia stanowiska całej debaty.
CREATE TABLE debata_stanowisko (
    id                       INTEGER PRIMARY KEY AUTOINCREMENT,
    identyfikator_zewnetrzny TEXT    NOT NULL UNIQUE,
    okno                     TEXT    NOT NULL,
    tura                     TEXT    NOT NULL DEFAULT '',
    tresc                    TEXT,
    -- Wersja rośnie przy każdym złożeniu stanowiska. Panel akcji Consensus
    -- Panel ma „wersjonowanie stanowiska”, a bez licznika nie ma co pokazać.
    wersja                   INTEGER NOT NULL DEFAULT 1,
    zaktualizowano           TEXT    NOT NULL DEFAULT (strftime('%Y-%m-%dT%H:%M:%fZ','now')),
    UNIQUE (okno, tura)
);
