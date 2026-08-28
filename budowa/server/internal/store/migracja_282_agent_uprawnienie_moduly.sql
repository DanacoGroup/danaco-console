-- Migracja 282 — piąta grupa zakresu uprawnień eksperta: moduły i zasoby.
--
-- Migracja 038 zamknęła kolumnę `grupa` na czterech wartościach. Permissions
-- Center ma cztery GRUPY ZAKRESU (agents.md rozdz. 9.2), ale nie te same
-- cztery: „dostęp do modułów i zasobów" jest jedną z nich i nie mieściła się
-- w słowniku. Bez tej wartości zawężenie do modułów Developer i Terminal —
-- przykład podany wprost w Specyfikacji agentów, Załącznik A — nie miało gdzie
-- się zapisać.
--
-- SQLite nie zmienia warunku CHECK w miejscu, więc tabela idzie przez
-- przebudowę. Na `agent_uprawnienie` nie wskazuje żaden klucz obcy, więc
-- przebudowa obejmuje jedną tabelę i kaskada nie ma czego zabrać.
--
-- Zakres szczegółowy wpisu grupy `modules` niesie KOD MODUŁU z katalogu
-- platformy. Warunku na to nie ma i być nie może: katalog modułów jest tabelą
-- `modul`, a CHECK nie sięga innych tabel. Wskazanie sprawdza rdzeń, który
-- katalog zna.

CREATE TABLE agent_uprawnienie_nowe (
    agent_id  INTEGER NOT NULL REFERENCES agent(id) ON DELETE CASCADE,
    grupa     TEXT    NOT NULL
                      CHECK (grupa IN ('files', 'network', 'processes', 'integrations', 'modules')),
    zakres    TEXT    NOT NULL DEFAULT '',
    przyznane INTEGER NOT NULL DEFAULT 1 CHECK (przyznane IN (0, 1)),
    PRIMARY KEY (agent_id, grupa, zakres)
);

INSERT INTO agent_uprawnienie_nowe (agent_id, grupa, zakres, przyznane)
SELECT agent_id, grupa, zakres, przyznane FROM agent_uprawnienie;

DROP TABLE agent_uprawnienie;
ALTER TABLE agent_uprawnienie_nowe RENAME TO agent_uprawnienie;
