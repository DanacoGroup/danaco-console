-- Tworzy tabele zespol i zespol_sklad, przechowujące nazwane, trwałe zestawy ekspertów kopiowane wielokrotnie przez Operatora, niezależnie od projektów i debat Roundtable.

CREATE TABLE zespol (
    id             INTEGER PRIMARY KEY AUTOINCREMENT,
    -- Identyfikator zewnętrzny; wraca polem `Team.id`.
    kod            TEXT    NOT NULL UNIQUE,
    nazwa          TEXT    NOT NULL,
    opis           TEXT    NOT NULL DEFAULT '',
    utworzono      INTEGER NOT NULL,
    -- NULL znaczy, że zespołu nie zmieniano od założenia, stąd pole updatedAt jest opcjonalne.
    zaktualizowano INTEGER
);

-- Wykaz zespołów biegnie zawsze po nazwie (kontrakt: „zespoły w kolejności
-- nazwy"), a zawężenie idzie po tej samej kolumnie.
CREATE INDEX idx_zespol_nazwa ON zespol (nazwa, kod);

CREATE TABLE zespol_sklad (
    zespol_id INTEGER NOT NULL REFERENCES zespol(id) ON DELETE CASCADE,
    -- Kod eksperta z biblioteki Agents, bez klucza obcego, aby usunięcie eksperta nie kasowało składu.
    agent_kod TEXT    NOT NULL,
    kolejnosc INTEGER NOT NULL,
    PRIMARY KEY (zespol_id, agent_kod)
);

-- Tworzy indeks składu zespołu po identyfikatorze zespołu i kolejności, wspierający odczyt zawsze w porządku nadanym przez Operatora.
CREATE INDEX idx_zespol_sklad_kolejnosc ON zespol_sklad (zespol_id, kolejnosc);
