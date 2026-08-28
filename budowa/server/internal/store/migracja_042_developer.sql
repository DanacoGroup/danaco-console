-- Migracja 042 — trwałość modułu Developer: wersja pliku zakładana przy zapisie
-- (okno Code Editor) i dziennik przebiegów budowania (okno Build Output).
--
-- Dlaczego wersja pliku ma własną tabelę. Kontrakt `developer.file.save` niesie
-- pole `createVersion` („czy założyć wersję w repozytorium"), a `DeveloperFile`
-- niesie `versionId`. Bez miejsca, w którym wersja naprawdę powstaje, oba pola
-- byłyby atrapą: rdzeń oddawałby identyfikator wersji, której nikt nigdzie nie
-- zapisał. Wersja jest migawką treści sprzed zapisu — to ona pozwala
-- Operatorowi wrócić do stanu, który sam nadpisał, i to ona odróżnia zapis
-- wersjonowany od zwykłego nadpisania pliku na dysku.
--
-- Wersja nie jest commitem i go nie zastępuje. Historia repozytorium należy do
-- Gita i jedzie komendą `developer.git.action`. Wersja jest zapisem roboczym
-- edytora — powstaje przed zatwierdzeniem, często wielokrotnie w obrębie jednej
-- zmiany, i ginie z projektem, a nie z gałęzią.
--
-- Dziennik budowań jest dziennikiem, nie stanem żywym. Przebieg czynny prowadzi
-- rdzeń w pamięci — tylko on ma uchwyt do drzewa procesu i tylko on potrafi je
-- przerwać. Tabela niesie ślad po przebiegu: zadanie, kod wyjścia, czasy i ogon
-- logu. Czyta ją montaż rdzenia (osierocenie przebiegów zostawionych przez
-- poprzedni bieg) oraz `developer.build.run` z `stop=true` dla okna, w którym
-- nic już nie biegnie — zamiast odmowy Operator dostaje ostatni znany przebieg.

-- ── Wersja pliku edytora (Code Editor) ───────────────────────────────────────
-- Kolumna `okno_kod` niesie identyfikator okna komunikacji w postaci kontraktu.
-- Klucza obcego do `okno_komunikacji` nie ma, bo okno bywa bytem pamięci
-- nadzorcy, który nie ma jeszcze wiersza, a wersja nie ma prawa nie powstać
-- z tego powodu.
--
-- `tresc` trzyma migawkę w całości. Wersja zakładana jest wyłącznie na wyraźne
-- żądanie (`createVersion`), więc nie rośnie przy każdym naciśnięciu klawisza;
-- wersjonowanie automatyczne byłoby tu kopiowaniem repozytorium do bazy.
CREATE TABLE developer_wersja_pliku (
    id         INTEGER PRIMARY KEY AUTOINCREMENT,
    kod        TEXT    NOT NULL UNIQUE,
    okno_kod   TEXT    NOT NULL,
    sciezka    TEXT    NOT NULL,
    tresc      TEXT    NOT NULL,
    rozmiar    INTEGER NOT NULL DEFAULT 0,
    utworzono  TEXT    NOT NULL DEFAULT (strftime('%Y-%m-%dT%H:%M:%fZ','now'))
);
CREATE INDEX idx_developer_wersja_sciezka ON developer_wersja_pliku(okno_kod, sciezka, utworzono DESC);

-- ── Dziennik przebiegów budowania (Build Output) ─────────────────────────────
-- `argumenty` przechowujemy jako jeden wiersz rozdzielony znakiem nowej linii,
-- a nie tabelą podrzędną: parametr zadania nie jest bytem, o który ktokolwiek
-- pyta osobno, a tabela na dwa pola byłaby złożonością bez odbiorcy.
--
-- `log` trzyma ogon przebiegu, nie całość. Log na żywo jedzie zdarzeniem
-- `developer.build.changed` i to ono jest nośnikiem podstawowym; kolumna służy
-- wyłącznie temu, żeby przebieg zamknięty przed połączeniem klienta dało się
-- jeszcze pokazać. Przycinanie ogona robi rdzeń — baza nie zna rozmiaru.
CREATE TABLE developer_budowanie (
    id          INTEGER PRIMARY KEY AUTOINCREMENT,
    kod         TEXT    NOT NULL UNIQUE,
    okno_kod    TEXT    NOT NULL,
    zadanie     TEXT    NOT NULL,
    argumenty   TEXT,
    stan        TEXT    NOT NULL DEFAULT 'running'
                        CHECK(stan IN ('running','succeeded','failed','stopped')),
    kod_wyjscia INTEGER,
    log         TEXT,
    uruchomiono TEXT    NOT NULL DEFAULT (strftime('%Y-%m-%dT%H:%M:%fZ','now')),
    zakonczono  TEXT
);
CREATE INDEX idx_developer_budowanie_okno ON developer_budowanie(okno_kod, uruchomiono DESC);
CREATE INDEX idx_developer_budowanie_stan ON developer_budowanie(stan, uruchomiono DESC);
