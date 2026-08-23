-- Migracja 059 — rejestr komponentów własnych Strefy 2 Strony głównej
-- (rodzina `component.*`: list, create, update, delete, assign oraz zdarzenie
-- `component.changed`).
--
-- Komponent własny jest bytem odrębnym — kaflem Strefy 2 — który wskazuje byt
-- magazynu modułowego kolumną `byt_docelowy` (kontraktowe `Component.targetId`).
-- Tabela nie powiela bytu modułowego: nie ma tu kroków automatyki, umiejętności
-- eksperta ani pamięci projektu. Prawdą o bycie modułowym pozostaje jego własna
-- tabela i jego własne komendy (`agent.*`, `automation.workflow.*`,
-- `workspace.*`).
--
-- Rodzina `component.*` nie składa się jako czysta fasada nad tabelami
-- modułowymi — stąd własny magazyn:
--  1. `ComponentKind` niesie wartość `assistant` („Profil asystenta”), a w bazie
--     stoją `zlecenie_asystenta` (wykonywane zlecenie wieloetapowe) i
--     `wpis_dziennika_asystenta` (dziennik czynności). Żadna z nich nie jest
--     profilem, a kontrakt nie ma komendy `assistant.profile.*`.
--  2. `Component.config` (json „przekazywany adapterowi modułu”) nie ma kolumny
--     w `projekt`, `automatyka` ani `agent`. Kolumna `agent.parametry_json` znaczy
--     co innego — parametry wywołania kanału modelu; wpisanie w nią konfiguracji
--     kafla byłoby dwiema prawdami o jednej kolumnie.
--  3. `component.assign` niesie parę (poziom zasięgu, identyfikator bytu poziomu),
--     a żadna z tabel modułowych nie ma kolumny zasięgu.
--  4. `Component.id` trzeba by składać z pary rodzaj+byt docelowy, a `component.list`
--     — sumować cztery niejednorodne zapytania bez wspólnego porządku wyświetlania,
--     którego kontrakt wymaga wprost („komponenty w kolejności wyświetlania”).
--
-- `byt_docelowy` jest tekstem bez więzu obcego, bo wskazuje na cztery różne
-- tabele zależnie od kolumny `rodzaj` (`projekt.kod`, `agent.kod`,
-- `automatyka.identyfikator_zewnetrzny`), a SQLite nie zna więzu warunkowego.
-- Kolumna niesie identyfikator zewnętrzny wprost i jest opcjonalna: kontrakt
-- oznacza `targetId` jako niewymagane.
--
-- `Component.createdAt/updatedAt` niosą milisekundy epoki, więc kolumny czasu są
-- typu INTEGER i przenoszą wartość kontraktu bez przekładu w obie strony.
--
-- Zasięg przypisania idzie zastanym słownikiem `poziom_zasiegu`, a nie własnym
-- wyliczeniem — poziomy zasięgu są jedne na całą platformę. Para kolumn
-- `poziom_zasiegu_id` + `klucz_zasiegu` powtarza wzorzec tabeli `ustawienie`.
-- Obie zostają NULL/puste, dopóki nie padnie `component.assign`: komponent
-- nieprzypisany nie ma poziomu i nie udaje, że stoi na globalnym.

-- ── Komponent własny Strefy 2 Strony głównej ─────────────────────────────────
CREATE TABLE komponent (
    id                       INTEGER PRIMARY KEY AUTOINCREMENT,
    identyfikator_zewnetrzny TEXT    NOT NULL UNIQUE,
    -- Wartości kontraktu (ComponentKind) wprost, bez tłumaczenia — kontrakt nie
    -- daje dla tego wyliczenia słownika przekładu bazy, więc go tu nie ma.
    rodzaj                   TEXT    NOT NULL
                                     CHECK(rodzaj IN ('automations','agents','workspace','assistant')),
    nazwa                    TEXT    NOT NULL,
    opis                     TEXT,
    -- Identyfikator zewnętrzny bytu magazynu modułowego; patrz nagłówek.
    byt_docelowy             TEXT,
    czynny                   INTEGER NOT NULL DEFAULT 1 CHECK(czynny IN (0,1)),
    -- Konfiguracja przekazywana adapterowi modułu. Poprawność JSON-a sprawdza
    -- warstwa `dane`, tak samo jak przy `agent.parametry_json`.
    konfiguracja             TEXT    NOT NULL DEFAULT '{}',
    -- Przypisanie do poziomu zasięgu (`component.assign`). NULL znaczy
    -- „nieprzypisany”, nie „globalny”.
    poziom_zasiegu_id        INTEGER          REFERENCES poziom_zasiegu(id) ON DELETE SET NULL,
    klucz_zasiegu            TEXT    NOT NULL DEFAULT '',
    -- Milisekundy epoki — patrz nagłówek.
    utworzono                INTEGER NOT NULL DEFAULT 0,
    zaktualizowano           INTEGER NOT NULL DEFAULT 0
);

-- Kontrakt każe oddać „komponenty w kolejności wyświetlania”. Indeks biegnie
-- porządkiem zapytania wykazu (`rodzaj, nazwa, id` — patrz
-- `dane/komponenty.go`), więc kolejność jest stała między wywołaniami i nie
-- wymaga sortowania po odczycie. Kolumny `czynny` w indeksie nie ma celowo:
-- stałaby między kolumną zawężenia a kolumnami porządku i zepsułaby ten drugi.
-- Przełącznik `includeDisabled` jest odsiewem na odczytanych wierszach.
CREATE INDEX idx_komponent_wykaz ON komponent(rodzaj, nazwa, id);

-- Wskazanie bytu modułowego służy odpowiedzi na pytanie „czy ten projekt ma już
-- swój kafel”. Bez indeksu odczyt szedłby przeglądem całej tabeli.
CREATE INDEX idx_komponent_byt_docelowy ON komponent(rodzaj, byt_docelowy);
