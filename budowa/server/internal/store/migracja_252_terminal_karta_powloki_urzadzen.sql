-- Migracja 252 — cztery powłoki urządzeniowe w warunku kolumny `powloka`.
--
-- Kontrakt zna dziesięć rodzajów powłoki karty; warunek CHECK z migracji 041
-- wymieniał sześć. Rdzeń nauczył się pozostałych czterech — kontenera, poda,
-- konsoli szeregowej i sesji Telnet — więc karta takiego rodzaju powstaje
-- w pamięci, ale jej ZAPIS odbijał się od warunku. Skutek byłby cichy w najgorszy
-- możliwy sposób: karta działa do restartu rdzenia, a po nim znika, i nic tego
-- nie zapowiada (zapis karty z zamysłu nie wywraca czynności — patrz
-- `zapiszKarte` w adapterze).
--
-- Warunku CHECK nie da się w SQLite zmienić poleceniem ALTER: tabelę trzeba
-- przebudować. Kroki idą w kolejności, która nie gubi ani jednego wiersza:
-- nowa tabela, przepisanie treści, zamiana nazwy, odtworzenie indeksu. Kolumny
-- są wymienione WPROST, a nie przez `SELECT *`, żeby przepisanie zależało od
-- schematu zapisanego tutaj, a nie od kolejności kolumn zastanej w bazie.
--
-- Migracja stoi PO 251, więc tabela ma już kolumny `cel_zdalny`, `port_zdalny`
-- i `host_kod` — one też wchodzą do nowego kształtu.

CREATE TABLE terminal_karta_nowa (
    id              INTEGER PRIMARY KEY AUTOINCREMENT,
    kod             TEXT    NOT NULL UNIQUE,
    okno_kod        TEXT    NOT NULL,
    powloka         TEXT    NOT NULL
                            CHECK(powloka IN ('powershell','cmd','bash','node','python','ssh',
                                              'container','pod','serial','telnet')),
    tytul           TEXT,
    katalog_roboczy TEXT,
    stan            TEXT    NOT NULL DEFAULT 'running' CHECK(stan IN ('running','exited')),
    utworzono       TEXT    NOT NULL DEFAULT (strftime('%Y-%m-%dT%H:%M:%fZ','now')),
    cel_zdalny      TEXT    NOT NULL DEFAULT '',
    port_zdalny     INTEGER,
    host_kod        TEXT
);

INSERT INTO terminal_karta_nowa
    (id, kod, okno_kod, powloka, tytul, katalog_roboczy, stan, utworzono,
     cel_zdalny, port_zdalny, host_kod)
SELECT id, kod, okno_kod, powloka, tytul, katalog_roboczy, stan, utworzono,
       cel_zdalny, port_zdalny, host_kod
FROM terminal_karta;

DROP TABLE terminal_karta;
ALTER TABLE terminal_karta_nowa RENAME TO terminal_karta;

CREATE INDEX idx_terminal_karta_okno ON terminal_karta(okno_kod, utworzono DESC);
