-- Migracja 375 — nośnik sygnału klas zdarzeń wyzwalających nakładkę.
--
-- Powód. Rozdz. 3.2 opracowania `funkcje-globalne/always-on-display.md` wymienia
-- SZEŚĆ klas zdarzeń wyzwalających. Rdzeń widział własną telemetrią dwie: stan
-- pętli wykonawczej i stan kolejki zadań (`aod/rozpoznanie-decyzji.ts` daje pięć
-- powodów i wszystkie wpadają w te dwie). Trzech pozostałych — wyniku kontroli
-- jakości, harmonogramu przebiegów automatyk i powtarzalności czynności
-- Operatora — nie widział NICZYM, więc ich wyciszenie zapisywało się i nie miało
-- czego wyciszać, a sugestia tych klas nie miała z czego powstać.
--
-- Tabela jest nośnikiem, nie drugą telemetrią. Sygnał odkłada tu ten, kto go
-- widzi: okno pętli wykonawczej, okno modułu po kontroli jakości, przebieg
-- automatyki po uruchomieniu z harmonogramu, dziennik czynności Operatora po
-- przekroczeniu progu powtarzalności (rozdz. 3.4 — trzy wystąpienia w karcie
-- sesji). Sygnał nie jest sugestią: sugestia powstaje dopiero wtedy, gdy sygnał
-- przekroczy próg i nie wpadnie w wyciszenie (rozdz. 3.1).
--
-- Sygnał wyciszony ODKŁADA SIĘ NADAL. Wyciszenie wstrzymuje UJAWNIENIE, nie
-- zapis: rozdz. 3.5 mówi wprost, że sugestie gromadzą się w liście oczekujących,
-- a rozdz. 3.1 — że zdarzenie niespełniające warunków „zostaje odnotowane
-- w kontekście funkcji i nie ujawnia się w interfejsie". Skasowanie sygnału
-- w chwili wyciszenia zabrałoby Operatorowi to, co miał zobaczyć po zniesieniu.
--
-- Liczba wystąpień (`liczba_wystapien`) stoi w wierszu, bo próg powtarzalności
-- czynności ręcznej z rozdz. 3.4 liczy się z niej — bez niej klasa „kontekst
-- pracy Operatora" nie miałaby po czym przekroczyć progu.
--
-- Sześć klas z CHECK, nie z tabeli słownikowej: klasy są wyliczeniem kontraktu
-- (`AodEventClass`), a nie katalogiem, do którego Operator dopisuje wiersze.
--
-- Wymaga restartu: nie.

CREATE TABLE IF NOT EXISTS sygnal_nakladki (
    id                     INTEGER PRIMARY KEY AUTOINCREMENT,
    identyfikator_zewnetrzny TEXT  NOT NULL UNIQUE,
    klasa_zdarzen          TEXT    NOT NULL
                                   CHECK(klasa_zdarzen IN ('executionLoopState', 'taskQueueState',
                                                           'qualityControlResult', 'moduleEvent',
                                                           'schedule', 'operatorWorkContext')),
    tresc                  TEXT    NOT NULL,
    modul_kod              TEXT    NOT NULL DEFAULT '',
    sesja_kod              TEXT    NOT NULL DEFAULT '',
    liczba_wystapien       INTEGER NOT NULL DEFAULT 1 CHECK(liczba_wystapien >= 1),
    zdarzylo_sie           TEXT    NOT NULL DEFAULT (strftime('%Y-%m-%dT%H:%M:%fZ','now'))
);

CREATE INDEX IF NOT EXISTS idx_sygnal_nakladki_klasa
    ON sygnal_nakladki(klasa_zdarzen, zdarzylo_sie DESC);
