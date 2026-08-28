-- Migracja 103 — zespoły ekspertów: nazwany skład biblioteki modułu Agents
-- (rodzina `team.*`).
--
-- Kontrakt niesie cztery komendy `team.save`, `team.load`, `team.list`
-- i `team.duplicate` oraz zdarzenie `team.changed`, a w schemacie nie było
-- niczego, co mogłoby je zasilić. Zespół nie jest ani projektem przestrzeni
-- roboczej (migracja 035), ani składem debaty Roundtable (migracja 044): tamten
-- skład żyje w jednej debacie i ginie razem z nią, a zespół jest bytem trwałym
-- Operatora, wybieranym wielokrotnie i kopiowanym.
--
-- Skład wskazuje kod, nie numer wiersza. Kolumna `agent_kod` niesie
-- identyfikator zewnętrzny eksperta — ten sam, którym posługuje się kontrakt
-- (`Agent.id` = `agent.kod`). Numer wiersza `agent.id` żyje wyłącznie wewnątrz
-- bazy i po przeniesieniu bazy przestałby się zgadzać, a zespół ma przeżyć
-- eksport i import biblioteki.
--
-- Bez klucza obcego do `agent` — celowo. Klucz obcy z kaskadą kasowałby wiersz
-- składu w chwili usunięcia eksperta, a zespół po cichu zmieniałby zawartość;
-- klucz obcy bez kaskady odmawiałby usunięcia eksperta należącego do zespołu,
-- czyli wprowadzałby bramę, której nikt nie ustanowił. Zamiast tego
-- ślad po ekspercie zostaje, a odczyt składu rozstrzyga, czy ekspert jest dziś
-- dostępny — łącząc `agent_kod` z tabelą `agent` z pominięciem archiwum
-- (kolumna `zarchiwizowano_o`, migracja 102). Ekspert wrócony z archiwum wraca
-- do składu sam, bo jego wiersz nigdy nie zniknął.
--
-- Kolejność nadaje Operator. Kontrakt mówi wprost: „eksperci wchodzący w skład,
-- w kolejności nadanej przez Operatora". Kolumna `kolejnosc` przenosi tę
-- kolejność; porządek alfabetyczny byłby zmyśleniem innej prawdy niż zapisana.
--
-- Czas w milisekundach epoki, nie w ISO 8601 — bo taki kształt niosą pola
-- `Team.createdAt` i `Team.updatedAt` kontraktu. Przekład formatu w kodzie
-- rdzenia byłby drugą prawdą o tej samej chwili.

CREATE TABLE zespol (
    id             INTEGER PRIMARY KEY AUTOINCREMENT,
    -- Identyfikator zewnętrzny; wraca polem `Team.id`.
    kod            TEXT    NOT NULL UNIQUE,
    nazwa          TEXT    NOT NULL,
    opis           TEXT    NOT NULL DEFAULT '',
    utworzono      INTEGER NOT NULL,
    -- NULL znaczy „zespół nie był zmieniany od założenia"; pole `updatedAt`
    -- kontraktu jest opcjonalne dokładnie z tego powodu.
    zaktualizowano INTEGER
);

-- Wykaz zespołów biegnie zawsze po nazwie (kontrakt: „zespoły w kolejności
-- nazwy"), a zawężenie idzie po tej samej kolumnie.
CREATE INDEX idx_zespol_nazwa ON zespol (nazwa, kod);

CREATE TABLE zespol_sklad (
    zespol_id INTEGER NOT NULL REFERENCES zespol(id) ON DELETE CASCADE,
    -- Kod eksperta z biblioteki modułu Agents (`agent.kod`). Bez klucza obcego —
    -- uzasadnienie w nagłówku.
    agent_kod TEXT    NOT NULL,
    kolejnosc INTEGER NOT NULL,
    PRIMARY KEY (zespol_id, agent_kod)
);

-- Skład czytany jest zawsze dla wskazanych zespołów, w kolejności Operatora.
CREATE INDEX idx_zespol_sklad_kolejnosc ON zespol_sklad (zespol_id, kolejnosc);
