-- Migracja 364 — odwracalny dziennik czynności dokumentu.
--
-- ── Czego rdzeń NIE miał ─────────────────────────────────────────────────────
-- Zmiana śledzona cofa się pojedynczo już dziś (`studio.tracking.decide`
-- przyjmuje ChangeIds). Zmiana JUŻ WNIESIONA do treści nie cofała się wcale —
-- nie było czego cofnąć, bo nie było zapisu, że coś się stało. Ten dziennik jest
-- tym zapisem: każda czynność (wpis, zmiana postaci, wstawienie tabeli,
-- przyjęcie propozycji modelu) odkłada wiersz, który da się wycofać.
--
-- ── Dlaczego stan sprzed i po, a nie sam opis ───────────────────────────────
-- Cofnięcie ma obejmować POSTAĆ, nie tylko treść: pomyłkowa zmiana kroju
-- w całym dokumencie musi się cofać tak samo jak skasowany akapit. Opis
-- słowny na to nie wystarcza — trzeba mieć stan, do którego się wraca. Dlatego
-- wiersz niesie `stan_przed_json` i `stan_po_json`: wycinek postaci objęty
-- czynnością, nie cały dokument.
--
-- Wycinek, nie całość, bo cofnięcie ma działać NIE PO KOLEI. Gdyby wiersz
-- trzymał migawkę całego dokumentu, cofnięcie czynności ze środka dziennika
-- zabrałoby ze sobą wszystko, co po niej weszło — czyli byłoby przywróceniem
-- wersji, a nie cofnięciem czynności.
--
-- ── Dlaczego zależności są osobną tabelą ────────────────────────────────────
-- Cofnięcie czynności ze środka dziennika, która jest podstawą późniejszej, ma
-- ODMÓWIĆ i nazwać zależność, a nie zostawić dokument w stanie niespójnym.
-- Zależność jest relacją wiele-do-wielu (wstawienie tabeli jest podstawą
-- scalenia komórki I policzenia szerokości), a wykaz kodów w kolumnie nie
-- pozwoliłby zapytać „co stoi na tej czynności" bez przeszukania wszystkich
-- wierszy. Kierunek: `czynnosc_id` stoi na `podstawa_id`.
--
-- ── Dlaczego kolejność jest osobną kolumną, a nie kluczem wiersza ───────────
-- Ponowienie czynności cofniętej nie zakłada nowego wpisu — przestawia stan
-- zastanego. Kolejność musi więc zostać ta sama, co przy pierwszym wykonaniu,
-- żeby dziennik dalej opisywał porządek pracy, a nie porządek zapisu.

CREATE TABLE czynnosc_dokumentu_studio (
    id                       INTEGER PRIMARY KEY AUTOINCREMENT,
    identyfikator_zewnetrzny TEXT    NOT NULL UNIQUE,
    dokument_id              INTEGER NOT NULL REFERENCES dokument_studio(id) ON DELETE CASCADE,
    kolejnosc                INTEGER NOT NULL,
    rodzaj                   TEXT    NOT NULL
                                     CHECK(rodzaj IN ('textEdit','formatChange','styleChange','pageChange',
                                           'listChange','tableChange','objectChange','apparatusChange',
                                           'markupChange','proposalAccept','importChange')),
    -- Rodzaj autora zostaje grubym rozróżnieniem człowiek/wykonawca; tożsamością
    -- jest kod agenta z modułu Agents. Agentów jest dowolnie wielu i są
    -- zakładani przez Operatora, więc wyliczenie nigdy by ich nie objęło.
    -- Kolumny agenta są nieobowiązkowe: czynności zapisane przed tą dobudową
    -- agenta nie mają i mają zostać poprawne.
    autor_rodzaj             TEXT    NOT NULL DEFAULT 'uzytkownik'
                                     CHECK(autor_rodzaj IN ('uzytkownik','model')),
    autor_agent_kod          TEXT,
    autor_agent_nazwa        TEXT,
    autor_agent_wersja       TEXT,
    autor_podagent_kod       TEXT,
    opis                     TEXT    NOT NULL,
    zakres_od                INTEGER,
    zakres_do                INTEGER,
    stan_przed_json          TEXT,
    stan_po_json             TEXT,
    zmiana_sledzona_kod      TEXT,
    zadanie_kod              TEXT,
    stan                     TEXT    NOT NULL DEFAULT 'active'
                                     CHECK(stan IN ('active','reverted')),
    utworzono                TEXT    NOT NULL DEFAULT (strftime('%Y-%m-%dT%H:%M:%fZ','now')),
    zaktualizowano           TEXT    NOT NULL DEFAULT (strftime('%Y-%m-%dT%H:%M:%fZ','now')),
    UNIQUE(dokument_id, kolejnosc)
);
CREATE INDEX idx_czynnosc_dokumentu_studio_dokument
    ON czynnosc_dokumentu_studio(dokument_id, kolejnosc DESC);
CREATE INDEX idx_czynnosc_dokumentu_studio_autor
    ON czynnosc_dokumentu_studio(dokument_id, autor_rodzaj, autor_agent_kod);
CREATE INDEX idx_czynnosc_dokumentu_studio_stan
    ON czynnosc_dokumentu_studio(dokument_id, stan, kolejnosc);

CREATE TABLE zaleznosc_czynnosci_studio (
    id                       INTEGER PRIMARY KEY AUTOINCREMENT,
    czynnosc_id              INTEGER NOT NULL
                                     REFERENCES czynnosc_dokumentu_studio(id) ON DELETE CASCADE,
    podstawa_id              INTEGER NOT NULL
                                     REFERENCES czynnosc_dokumentu_studio(id) ON DELETE CASCADE,
    powod                    TEXT,
    UNIQUE(czynnosc_id, podstawa_id),
    -- Czynność nie stoi na sobie samej; taki wiersz uczyniłby cofnięcie
    -- niemożliwym bez podania powodu, którego nikt nie mógłby naprawić.
    CHECK(czynnosc_id <> podstawa_id)
);
CREATE INDEX idx_zaleznosc_czynnosci_studio_podstawa
    ON zaleznosc_czynnosci_studio(podstawa_id, czynnosc_id);
