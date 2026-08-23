-- Migracja 070 — katalog rozszerzeń (rodzina `extension.*`: list, install,
-- configure, toggle, uninstall).
--
-- Pozycja katalogu nie jest ani mostem MCP, ani konektorem eksperta:
--  1. Most mieszka w `punkt_dostepu`, skąd `core/most_okna.go` składa wpisy
--     `mcpServers` procesu modelu, i opisuje wyłącznie maszynę udostępnioną
--     mostem albo katalog lokalny. `ExtensionKind` niesie cztery wartości —
--     `mcp`, `plugin`, `api`, `skill` — a wtyczki, integracji API ani skilla
--     w punkcie dostępu zapisać się nie da.
--  2. `agent_konektor` ma kolumnę `agent_id NOT NULL`, więc konektor należy do
--     jednego eksperta. Katalog stoi poziom wyżej: pole `extension.list.agentId`
--     jest ekspertem, „dla którego liczona jest przypisywalność pozycji", więc
--     pozycja istnieje niezależnie od eksperta. `AgentConnectorKind` ma przy tym
--     trzy wartości, a `ExtensionKind` cztery — skilla konektorem nie zapiszesz
--     (umiejętności eksperta trzyma `agent_umiejetnosc`).
--  3. Pozycję katalogu ekspert dopiero bierze: rodzaje `mcp`/`plugin`/`api` —
--     zakładając sobie `agent_konektor`, rodzaj `skill` — zakładając
--     `agent_umiejetnosc`. Ta tabela jest warstwą nad tamtymi dwiema, nie kopią.
--
-- Rozszerzenie rodzaju `mcp` wskazuje most, a nie powiela go: kolumna
-- `punkt_dostepu_id` wskazuje wiersz `punkt_dostepu` tak samo jak
-- `agent_konektor.punkt_dostepu_id`. Drugiego rejestru serwerów MCP platforma
-- nie ma. Więz jest ON DELETE SET NULL: skasowanie punktu odłącza pozycję
-- katalogu, lecz jej nie kasuje.
--
-- Zawartości kolumny `zrodlo_deklarowane` rdzeń nie pobiera i nie uruchamia —
-- kontrakt zostawia znaczenie instalacji otwarte („pobranie paczki,
-- zarejestrowanie adresu czy zapis punktu dostępu"), a pobieranie i wykonywanie
-- cudzego kodu zmieniłoby klasę bezpieczeństwa produktu. Kolumna trzyma napis
-- podany w polu `source`, żeby go nie zgubić. Kształt `Extension` pola `source`
-- nie ma, więc kolumna nie wychodzi w odpowiedzi.
--
-- Identyfikatory są dwa, bo kontrakt ma dwa. `Extension.id` jest tożsamością
-- wiersza (kolumna `identyfikator_zewnetrzny`, którą niosą `configure`, `toggle`
-- i `uninstall`), a `Extension.code` — kodem pozycji „stałym między wydaniami"
-- (kolumna `kod`, którą niesie `install`). `kod` jest UNIQUE: dwie pozycje
-- katalogu o jednym kodzie znaczyłyby, że kod nie jest stały.
--
-- Odinstalowanie nie kasuje wiersza — inaczej `extension.list` nie miałby czego
-- zawężać polem `installedOnly`. Zdejmuje `zainstalowane` i `wlaczone`,
-- zostawiając pozycję w katalogu razem z konfiguracją; `extension.install` tym
-- samym kodem ją przywraca.
--
-- `Extension.updatedAt` niesie milisekundy epoki, więc kolumna jest typu INTEGER
-- i przenosi wartość kontraktu bez przekładu. Kolumny `utworzono` nie ma, bo
-- kształt `Extension` nie niesie `createdAt`.

-- ── Pozycja katalogu rozszerzeń ──────────────────────────────────────────────
CREATE TABLE rozszerzenie (
    id                       INTEGER PRIMARY KEY AUTOINCREMENT,
    -- `Extension.id` — tożsamość wiersza; nią wołają configure, toggle, uninstall.
    identyfikator_zewnetrzny TEXT    NOT NULL UNIQUE,
    -- `Extension.code` — kod pozycji stały między wydaniami; nim woła install.
    kod                      TEXT    NOT NULL UNIQUE,
    -- Wartości kontraktu (ExtensionKind) wprost, bez tłumaczenia — tak samo jak
    -- `agent_konektor.rodzaj`. Kontrakt nie daje dla tego wyliczenia słownika
    -- przekładu bazy, więc go tu nie ma.
    rodzaj                   TEXT    NOT NULL
                                     CHECK(rodzaj IN ('mcp', 'plugin', 'api', 'skill')),
    nazwa                    TEXT    NOT NULL,
    opis                     TEXT,
    wersja                   TEXT,
    zainstalowane            INTEGER NOT NULL DEFAULT 0 CHECK(zainstalowane IN (0, 1)),
    wlaczone                 INTEGER NOT NULL DEFAULT 0 CHECK(wlaczone IN (0, 1)),
    -- Wskazanie mostu; patrz nagłówek.
    punkt_dostepu_id         INTEGER          REFERENCES punkt_dostepu(id) ON DELETE SET NULL,
    -- Napis podany w polu `source`, nie adres do pobrania; patrz nagłówek.
    zrodlo_deklarowane       TEXT    NOT NULL DEFAULT '',
    -- Konfiguracja rozszerzenia. Poprawność JSON-a sprawdza warstwa `dane`, tak
    -- samo jak przy `komponent.konfiguracja`.
    konfiguracja             TEXT    NOT NULL DEFAULT '{}',
    -- Milisekundy epoki — patrz nagłówek.
    zaktualizowano           INTEGER NOT NULL DEFAULT 0
);

-- Kontrakt każe oddać „pozycje katalogu w kolejności wyświetlania". Indeks
-- biegnie porządkiem zapytania wykazu (`rodzaj, nazwa, id` — patrz
-- `dane/extension.go`), więc kolejność jest stała między wywołaniami. Kolumny
-- `zainstalowane` w indeksie nie ma: stałaby między kolumną zawężenia
-- a kolumnami porządku i zepsułaby ten drugi. Przełącznik `installedOnly` jest
-- warunkiem zapytania, nie kolumną porządku.
CREATE INDEX idx_rozszerzenie_wykaz ON rozszerzenie(rodzaj, nazwa, id);

-- Odpowiedź na pytanie „które pozycje katalogu wiszą na tym moście". Bez indeksu
-- odczyt szedłby przeglądem całej tabeli.
CREATE INDEX idx_rozszerzenie_punkt ON rozszerzenie(punkt_dostepu_id);
