-- Migracja 279 — trzy dopełnienia układu zależności: bramka dołączenia, grupa
-- kroków i krok wycofujący (`orchestration.gate.set`,
-- `orchestration.group.set`, `orchestration.compensation.set`).
--
-- Wszystkie trzy wiszą na `automatyka(id)` i posługują się identyfikatorem
-- ZEWNĘTRZNYM kroku — tym samym napisem, którym posługuje się
-- `zaleznosc_kroku_automatyki` (migracja 039). Klucz wiersza kroku nie wchodzi
-- tu z rozmysłem: układ zależności przepisuje się w całości
-- (`automation.orchestrator.define`), więc kroki bywają zakładane od nowa,
-- a wtedy klucz wiersza jest inny, choć krok w oczach Operatora ten sam.
--
-- Bramka dołączenia jest jedna na krok: krok scalający tory ma jedną regułę
-- scalenia i drugiej mieć nie może. Stąd klucz pierwotny na parze
-- (automatyka, krok), a nie osobny identyfikator.
--
-- Grupa ma identyfikator własny, bo kontrakt `orchestration.group.set` przyjmuje
-- `groupId` przy zmianie i pozwala go pominąć przy założeniu. Kroki grupy leżą
-- w tabeli podrzędnej z zachowaną kolejnością — grupa sekwencyjna to grupa,
-- w której kolejność jest treścią, a nie ozdobą.
--
-- Kompensacja jest jedna na krok główny: krok wycofujący skutki kroku N jest
-- jeden, bo dwa wycofania tego samego kroku dałyby stan, którego nikt nie
-- umiałby rozstrzygnąć. Zdjęcie kompensacji kasuje wiersz — to jest znaczenie
-- pustego `compensationStepId` z kontraktu.

CREATE TABLE orkiestracja_bramka (
    automatyka_id INTEGER NOT NULL REFERENCES automatyka(id) ON DELETE CASCADE,
    krok          TEXT    NOT NULL,
    regula        TEXT    NOT NULL DEFAULT 'all'
                          CHECK (regula IN ('all', 'any', 'count')),
    -- Liczba torów wymagana przy regule licznikowej; przy pozostałych zero.
    licznik       INTEGER NOT NULL DEFAULT 0 CHECK (licznik >= 0),
    zapisano      TEXT    NOT NULL DEFAULT (strftime('%Y-%m-%dT%H:%M:%fZ','now')),
    PRIMARY KEY (automatyka_id, krok),
    -- Reguła licznikowa bez liczby torów byłaby bramką, która nigdy się nie
    -- otworzy albo otworzy zawsze — zależnie od tego, kto ją później czyta.
    CHECK (regula <> 'count' OR licznik > 0)
);

CREATE TABLE orkiestracja_grupa (
    id            INTEGER PRIMARY KEY AUTOINCREMENT,
    kod           TEXT    NOT NULL UNIQUE,
    automatyka_id INTEGER NOT NULL REFERENCES automatyka(id) ON DELETE CASCADE,
    nazwa         TEXT    NOT NULL,
    rodzaj        TEXT    NOT NULL DEFAULT 'parallel'
                          CHECK (rodzaj IN ('sequential', 'parallel', 'conditional')),
    zapisano      TEXT    NOT NULL DEFAULT (strftime('%Y-%m-%dT%H:%M:%fZ','now'))
);

CREATE INDEX idx_orkiestracja_grupa_automatyka
    ON orkiestracja_grupa(automatyka_id, nazwa);

CREATE TABLE orkiestracja_grupa_krok (
    grupa_id  INTEGER NOT NULL REFERENCES orkiestracja_grupa(id) ON DELETE CASCADE,
    krok      TEXT    NOT NULL,
    kolejnosc INTEGER NOT NULL DEFAULT 0,
    PRIMARY KEY (grupa_id, krok)
);

CREATE TABLE orkiestracja_kompensacja (
    automatyka_id INTEGER NOT NULL REFERENCES automatyka(id) ON DELETE CASCADE,
    krok          TEXT    NOT NULL,
    krok_wycofu   TEXT    NOT NULL,
    zapisano      TEXT    NOT NULL DEFAULT (strftime('%Y-%m-%dT%H:%M:%fZ','now')),
    PRIMARY KEY (automatyka_id, krok),
    -- Krok wycofujący sam siebie nie wycofuje niczego.
    CHECK (krok <> krok_wycofu)
);
