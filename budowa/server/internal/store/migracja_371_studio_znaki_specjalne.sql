-- Migracja 371 — znaki ostatnio użyte i uzupełnienie wykazu autozamiany.
--
-- ── Czego ta migracja NIE zakłada i dlaczego ────────────────────────────────
-- Nie zakłada tabeli `autozamiana_znaku_studio`. Ta stoi już w migracji 368
-- wraz z kolumną `fabryczna` i zasadami fabrycznymi wpisanymi wierszami. Druga
-- tabela o tej nazwie — albo wykaz zasad fabrycznych powtórzony w kodzie
-- rdzenia — dałaby dwie prawdy o tym, co wchodzi w miejsce skrótu „(c)".
-- Rdzeń czyta więc wykaz z tamtej tabeli, a tutaj dokłada wyłącznie to, czego
-- tamten wykaz nie niósł.
--
-- ── Dlaczego znaki ostatnio użyte mają wiersz na znak, nie na użycie ────────
-- Pytanie brzmi „czego Operator używa", nie „ile razy dziś nacisnął". Licznik
-- i czas ostatniego użycia wystarczają, żeby wykaz ułożyć od najbliższego ręce,
-- a jeden wiersz na użycie zasypałby tabelę w godzinę pisania.
--
-- ── Dlaczego wykaz nie ma dokumentu ────────────────────────────────────────
-- `studio.symbol.list` i `studio.symbol.autoreplace.*` nie przyjmują w kontrakcie
-- dokumentu, i słusznie: znak, którym Operator posłużył się wczoraj w umowie, ma
-- być pod ręką i dziś w notatce. Wykaz przypięty do dokumentu byłby pusty przy
-- każdym nowym pismie.

CREATE TABLE znak_ostatnio_uzyty_studio (
    id        INTEGER PRIMARY KEY AUTOINCREMENT,
    kod       TEXT    NOT NULL UNIQUE,
    znak      TEXT    NOT NULL,
    ile_uzyc  INTEGER NOT NULL DEFAULT 1,
    uzyto     TEXT    NOT NULL DEFAULT (strftime('%Y-%m-%dT%H:%M:%fZ','now'))
);
CREATE INDEX idx_znak_ostatnio_uzyty_studio_kolejnosc
    ON znak_ostatnio_uzyty_studio(uzyto DESC, ile_uzyc DESC);

-- Uzupełnienie wykazu zasad fabrycznych z migracji 368. Tamten wykaz niesie
-- znaki interpunkcyjne i matematyczne, ale nie niesie znaków prawniczych ani
-- ułamków, a Właściciel wymienia znaki prawnicze wprost. `OR IGNORE` sprawia, że
-- krok jest bezpieczny wobec bazy, w której któryś ze skrótów już stoi —
-- kolumna `skrot` ma warunek UNIQUE i bez tego migracja wywróciłaby się na
-- pierwszej powtórce.
INSERT OR IGNORE INTO autozamiana_znaku_studio (skrot, zamiennik, fabryczna) VALUES
    ('(par)', '§',  1),
    ('(nr)',  '№',  1),
    ('(st)',  '°',  1),
    ('<->',   '↔',  1),
    ('=>',    '⇒',  1),
    ('1/2',   '½',  1),
    ('1/4',   '¼',  1),
    ('3/4',   '¾',  1);
