-- Migracja 037 — ekspert modułu Agents jako komponent własny.
--
-- Powód. Kontrakt niesie osiem komend obszaru `agent.*` i strukturę Agent, ale
-- schemat nie miał ani jednej tabeli, w której ekspert mógłby zamieszkać.
-- Rejestr okien operacyjnych zna już pięć okien modułu
-- Agents; ta migracja daje im byt, na którym pracują.
--
-- Kod jest identyfikatorem trwałym. Wzorcem jest `kanal_modelu` i `punkt_dostepu`:
-- `id` służy powiązaniom wewnątrz bazy, a `kod` wychodzi na zewnątrz jako pole
-- `Agent.id` kontraktu. Zmiana nazwy eksperta nie zmienia jego kodu, więc
-- przypisanie w projekcie (`workspace.agent.assign`) nie rozjeżdża się z nazwą.
--
-- Kanał, nie drugi rejestr modeli. `kanal_kod` wskazuje wiersz
-- `kanal_modelu.kod` — ten sam rejestr, z którego korzysta okno rozmowy. Kolumna
-- nie jest kluczem obcym z zamysłem: ekspert bywa zakładany zanim kanał
-- powstanie, a kanał bywa kasowany bez kasowania ekspertów. Spójność sprawdza
-- warstwa `dane` przy zapisie, a kanał nieistniejący znaczy „ekspert bez modelu
-- bazowego”, nie „ekspert niezdatny do pracy”.
--
-- Wersja jest licznikiem, nie archiwum. `wersja` rośnie przy każdej zmianie
-- tożsamości (agent.update, agent.model.set) i odpowiada polu `Agent.version`.
-- Tabeli archiwum treści tu nie ma, ponieważ kontrakt nie ma komendy, która
-- mogłaby ją odczytać — założenie takiej tabeli byłoby zapisem bez czytelnika.
-- Panel „Historia wersji” składa wykaz ze zdarzeń `agent.changed`
-- widzianych w toku sesji, a przywrócenie wykonuje przez `agent.update`.
--
-- Umiejętność jest napisem, nie wierszem katalogu. Kontrakt przenosi skille
-- polem `Agent.skillIds` i przyjmuje je komendą `agent.skill.add` jako `skillId`.
-- Katalog dostępnych rozszerzeń szedłby komendami `extension.*`, których kontrakt
-- na dziś nie ma — tabeli katalogu więc nie zakładamy, bo nie miałaby czytelnika.

-- ── Ekspert ───────────────────────────────────────────────────────────────────
CREATE TABLE agent (
    id                   INTEGER PRIMARY KEY AUTOINCREMENT,
    kod                  TEXT    NOT NULL UNIQUE,
    nazwa                TEXT    NOT NULL,
    opis                 TEXT    NOT NULL DEFAULT '',
    instrukcje_systemowe TEXT    NOT NULL DEFAULT '',
    kanal_kod            TEXT,
    model                TEXT,
    -- Droga wywołania kanału; wartości shared.ProviderTransport wprost, bo
    -- kontrakt nie daje dla tego wyliczenia słownika przekładu bazy.
    transport            TEXT,
    -- Parametry wywołania zależne od kanału; sprawdzane jako poprawny JSON
    -- przez warstwę `dane`, tak samo jak `kanal_modelu.parametry_json`.
    parametry_json       TEXT    NOT NULL DEFAULT '{}',
    wersja               INTEGER NOT NULL DEFAULT 1,
    aktywny              INTEGER NOT NULL DEFAULT 1 CHECK (aktywny IN (0, 1)),
    utworzono            TEXT    NOT NULL DEFAULT (strftime('%Y-%m-%dT%H:%M:%fZ', 'now')),
    zaktualizowano       TEXT    NOT NULL DEFAULT (strftime('%Y-%m-%dT%H:%M:%fZ', 'now'))
);

-- Wykaz biblioteki ekspertów idzie po nazwie (kontrakt: „eksperci w kolejności
-- nazw”), a filtr `enabledOnly` po kolumnie `aktywny`.
CREATE INDEX idx_agent_nazwa ON agent(aktywny, nazwa);

-- ── Umiejętności przypisane ekspertowi ────────────────────────────────────────
-- Klucz główny na parze nie pozwala przypisać tej samej umiejętności dwa razy,
-- więc powtórzone `agent.skill.add` jest bezskutkowe zamiast dublować wiersz.
CREATE TABLE agent_umiejetnosc (
    agent_id  INTEGER NOT NULL REFERENCES agent(id) ON DELETE CASCADE,
    kod       TEXT    NOT NULL,
    dodano    TEXT    NOT NULL DEFAULT (strftime('%Y-%m-%dT%H:%M:%fZ', 'now')),
    PRIMARY KEY (agent_id, kod)
);
