-- Migracja 035 — trwałość modułu Workspace: projekt, pamięć projektu (Context
-- Memory) i przypisanie eksperta do projektu (Agent Manager).
--
-- Projekt jest bytem, nie napisem. Kolumna `sesja.projekt` niesie sam napis,
-- bez nazwy, stanu i czasu ostatniej zmiany — a Project Dashboard ma pokazać
-- dokładnie te trzy rzeczy. Tabela `projekt` daje im miejsce, a kolumna `kod`
-- wiąże wiersz z napisem, którym posługują się sesje i kontrakt
-- (`WorkspaceProject.id`). Kolumny `sesja.projekt` nie ruszamy: wiąże się po
-- kodzie, więc sesje wiążące się samym napisem zostają czytelne.
--
-- Instrukcji projektu tu nie ma i nie jest to przeoczenie. Instrukcje systemowe
-- podlegają dziedziczeniu warstwowemu ośmiu poziomów zasięgu, które
-- rozstrzyga pakiet `internal/konfig` nad tabelą `ustawienie`. Druga tabela
-- instrukcji znaczyłaby drugi rozstrzygacz, więc `workspace.instructions.set`
-- zapisuje ustawienie o kluczu `workspace.instrukcje` pod adresem wskazanego
-- poziomu — dokładnie tam, gdzie mieszka reszta konfiguracji warstwowej.
--
-- Biblioteki projektu tu nie ma — również z zamysłem. Pliki projektu leżą
-- w katalogu roboczym projektu, a `workspace.library.list`
-- czyta ten katalog. Tabela plików bez komendy zapisu w module Workspace byłaby
-- pusta i bez pisarza; repozytorium wiedzy jest bytem modułu Library.

-- ── Projekt przestrzeni roboczej ──────────────────────────────────────────────
CREATE TABLE projekt (
    id             INTEGER PRIMARY KEY AUTOINCREMENT,
    kod            TEXT    NOT NULL UNIQUE,
    nazwa          TEXT    NOT NULL,
    opis           TEXT,
    stan           TEXT    NOT NULL DEFAULT 'active' CHECK(stan IN ('active','archived')),
    utworzono      TEXT    NOT NULL DEFAULT (strftime('%Y-%m-%dT%H:%M:%fZ','now')),
    zaktualizowano TEXT    NOT NULL DEFAULT (strftime('%Y-%m-%dT%H:%M:%fZ','now'))
);

-- ── Pamięć projektu — Context Memory ──────────────────────────────────────────
-- Wpis niesie poziom zasięgu współdzielenia: domyślnie projekt, lecz Operator
-- może wynieść ustalenie wyżej albo zawęzić je do karty sesji czy okna.
-- Pochodzenie odróżnia ustalenie Operatora od propozycji modelu, bo propozycja
-- czeka na przyjęcie i nie ma prawa wejść do pracy sama.
CREATE TABLE wpis_pamieci_projektu (
    id                       INTEGER PRIMARY KEY AUTOINCREMENT,
    projekt_id               INTEGER NOT NULL REFERENCES projekt(id) ON DELETE CASCADE,
    identyfikator_zewnetrzny TEXT    NOT NULL UNIQUE,
    tresc                    TEXT    NOT NULL,
    przypiety                INTEGER NOT NULL DEFAULT 0 CHECK(przypiety IN (0,1)),
    pochodzenie              TEXT    NOT NULL DEFAULT 'operator'
                                     CHECK(pochodzenie IN ('operator','model')),
    poziom_zasiegu_id        INTEGER NOT NULL REFERENCES poziom_zasiegu(id),
    klucz_zasiegu            TEXT    NOT NULL DEFAULT '',
    utworzono                TEXT    NOT NULL DEFAULT (strftime('%Y-%m-%dT%H:%M:%fZ','now')),
    zaktualizowano           TEXT    NOT NULL DEFAULT (strftime('%Y-%m-%dT%H:%M:%fZ','now'))
);
CREATE INDEX idx_wpis_pamieci_projektu_projekt
    ON wpis_pamieci_projektu(projekt_id, przypiety DESC, zaktualizowano DESC);

-- ── Przypisanie eksperta do projektu — Agent Manager ──────────────────────────
-- Ekspert jest komponentem własnym modułu Agents i mieszka poza tą
-- tabelą; tu zapisujemy wyłącznie fakt przypisania go do projektu wraz z rolą
-- i wskazaniem domyślnego wykonawcy. Kolumna `agent_kod` niesie identyfikator
-- kontraktu, więc przypisanie nie czeka na tabelę modułu Agents.
CREATE TABLE przypisanie_agenta_projektu (
    projekt_id         INTEGER NOT NULL REFERENCES projekt(id) ON DELETE CASCADE,
    agent_kod          TEXT    NOT NULL,
    rola               TEXT,
    domyslny_wykonawca INTEGER NOT NULL DEFAULT 0 CHECK(domyslny_wykonawca IN (0,1)),
    przypisano         TEXT    NOT NULL DEFAULT (strftime('%Y-%m-%dT%H:%M:%fZ','now')),
    PRIMARY KEY (projekt_id, agent_kod)
);
CREATE INDEX idx_przypisanie_agenta_projektu_domyslny
    ON przypisanie_agenta_projektu(projekt_id, domyslny_wykonawca DESC, agent_kod);
