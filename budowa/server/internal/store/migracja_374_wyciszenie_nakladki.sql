-- Migracja 374 — wyciszenie nakładki Always On Display jako byt rdzenia.
--
-- Powód jest mierzalny: do tej migracji wyciszenie było stanem JEDNEGO OKNA
-- (`client/src/aod/wyciszenie-aod.ts`, zapis miejscowy przeglądarki), więc
-- Operator wyciszał sugestie w jednej powłoce, a w drugiej wchodziły dalej.
-- Rozstrzygnięcie Właściciela z 17.08.2026 żąda wyciszenia od ręki z pozycji
-- awatara; żeby cisza sięgnęła wszystkich powłok, musi mieć wiersz tutaj.
--
-- Trzy rodzaje, nie pięć. Tabela rozdz. 3.5 opracowania
-- `funkcje-globalne/always-on-display.md` wymienia pięć wierszy i tylko trzy
-- pierwsze są wyciszeniami: czasowe, kontekstowe i klasy zdarzeń. Tryb cichy
-- jest TRYBEM OBECNOŚCI (rozdz. 9.2) i jedzie ustawieniem konfiguracji, a wyjątek
-- wagi krytycznej nie jest wyciszeniem, lecz regułą przebijającą każde z trzech.
-- Wpisanie ich tutaj dałoby dwa znaczenia jednej tabeli.
--
-- Czas końca (`konczy_sie`) niesie WYŁĄCZNIE wyciszenie czasowe. Pozostałe dwa
-- trwają do zniesienia ręką Operatora i pusta wartość mówi to wprost — zamiast
-- udawać koniec datą odległą, po której nikt nie umiałby powiedzieć, do kiedy
-- właściwie milczy.
--
-- Nazwa bytu zakresu (`nazwa_zakresu`) stoi obok jego identyfikatora, bo zasada
-- zlecenia jest bezwzględna: cisza, po której Operator nie wie, CO milczy, jest
-- gorsza od braku wyciszenia. Wykaz wyciszeń ma się dać wypisać bez odpytywania
-- wykazu modułów i wykazu kart sesji.
--
-- Urządzenie (`urzadzenie_id`) mówi, SKĄD wyciszenie przyszło, a nie gdzie
-- obowiązuje. Obowiązuje wszędzie — po to jest ta tabela.
--
-- Wymaga restartu: nie.

CREATE TABLE IF NOT EXISTS wyciszenie_nakladki (
    id                     INTEGER PRIMARY KEY AUTOINCREMENT,
    identyfikator_zewnetrzny TEXT  NOT NULL UNIQUE,
    rodzaj                 TEXT    NOT NULL
                                   CHECK(rodzaj IN ('timed', 'contextual', 'eventClass')),
    zakres                 TEXT    NOT NULL DEFAULT ''
                                   CHECK(zakres IN ('', 'module', 'session', 'eventClass')),
    klucz_zakresu          TEXT    NOT NULL DEFAULT '',
    nazwa_zakresu          TEXT    NOT NULL DEFAULT '',
    klasa_zdarzen          TEXT    NOT NULL DEFAULT ''
                                   CHECK(klasa_zdarzen IN ('', 'executionLoopState', 'taskQueueState',
                                                           'qualityControlResult', 'moduleEvent',
                                                           'schedule', 'operatorWorkContext')),
    urzadzenie_id          TEXT    NOT NULL DEFAULT '',
    konczy_sie             TEXT    NOT NULL DEFAULT '',
    utworzono              TEXT    NOT NULL DEFAULT (strftime('%Y-%m-%dT%H:%M:%fZ','now')),
    -- Wyciszenie czasowe bez chwili końca byłoby ciszą bez terminu; wyciszenie
    -- kontekstowe bez bytu i wyciszenie klasy bez klasy nie mają czego wyciszyć.
    CHECK((rodzaj = 'timed' AND zakres = '' AND klasa_zdarzen = '' AND konczy_sie <> '')
          OR (rodzaj = 'contextual' AND zakres IN ('module', 'session') AND klucz_zakresu <> '')
          OR (rodzaj = 'eventClass' AND zakres = 'eventClass' AND klasa_zdarzen <> ''))
);

-- Jedno wyciszenie na byt. Wyciszenie czasowe jest jedno w całej platformie:
-- dwa czasy naraz nie dałyby Operatorowi jednej odpowiedzi na pytanie „do kiedy".
CREATE UNIQUE INDEX IF NOT EXISTS idx_wyciszenie_nakladki_byt
    ON wyciszenie_nakladki(rodzaj, zakres, klucz_zakresu, klasa_zdarzen);
