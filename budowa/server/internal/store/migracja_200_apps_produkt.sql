-- Migracja 200 — moduł Apps, strona budowy produktu: produkt okna (Product
-- Builder), etapy budowy, kamienie milowe i ich związek z etapami.
--
-- Produkt jest jeden na okno, nie wiele. Kontrakt daje `apps.product.get`
-- z żądaniem niosącym samo `windowId` i wynikiem o jednym polu `product`,
-- a `apps.product.save` nie ma pola `productId` — nie ma czym wskazać drugiego
-- produktu tego samego okna. Dlatego kolumna `okno` jest UNIQUE, a zapis jest
-- UPSERT-em po niej. Bez tego warunku dwa zapisy z rzędu zakładałyby dwa
-- wiersze, a odczyt musiałby zgadywać, który z nich jest produktem okna.
--
-- Platformy docelowe leżą w jednej kolumnie tekstowej rozdzielonej znakiem
-- nowej linii, tym samym wzorcem co `architektura_apps.zastrzezenia_walidacji`
-- (migracja 051). Nikt nie filtruje ani nie sortuje po pojedynczej platformie:
-- `AppProduct.platforms` wychodzi zawsze w komplecie razem z produktem.
--
-- Etap ma własną tabelę, nie kolumnę produktu. `apps.stage.save` zmienia etap
-- po jego identyfikatorze, `apps.stage.list` zawęża po stanie, a
-- `AppMilestone.stageIds` wiąże kamień milowy z etapami — każda z tych trzech
-- dróg wymaga wiersza na etap. Etap należy do okna, nie do produktu: żądanie
-- `apps.stage.list` niesie `windowId`, a nie identyfikator produktu, więc okno
-- bez zapisanego produktu ma prawo mieć etapy.
--
-- Kolejność etapu jest kolumną, nie porządkiem wstawiania. Tracker etapów
-- Product Buildera rysuje oś w kolejności zadanej przez Operatora
-- (`AppStage.order`), a ta zmienia się bez zakładania wierszy na nowo.
--
-- Związek kamienia milowego z etapami ma tabelę złącznikową, nie kolumnę listy.
-- `apps.milestone.save` nadsyła `stageIds` w komplecie przy każdym zapisie
-- (kontrakt nie ma trybu częściowej zmiany), więc zapis wymienia wiersze
-- związku „usuń, wstaw od nowa" — jedna prawda o krawędzi, nie kopia w dwóch
-- miejscach. Kolumna `etap_kod` trzyma kod zewnętrzny etapu, a nie więz obcy:
-- kamień milowy ma prawo wskazywać etap usunięty po zapisie, tak samo jak
-- `plik_warsztatu_apps.komponent_id` wskazuje komponent zdjęty z architektury.

-- ── Produkt okna — Product Builder ───────────────────────────────────────────
CREATE TABLE produkt_apps (
    id                       INTEGER PRIMARY KEY AUTOINCREMENT,
    identyfikator_zewnetrzny TEXT    NOT NULL UNIQUE,
    -- Jeden produkt na okno — patrz rozstrzygnięcie na czole pliku.
    okno                     TEXT    NOT NULL UNIQUE,
    nazwa                    TEXT    NOT NULL,
    opis                     TEXT,
    -- Wartości kontraktu (AppProductPlatform) rozdzielone znakiem nowej linii.
    platformy                TEXT,
    repozytorium             TEXT,
    utworzono                TEXT    NOT NULL DEFAULT (strftime('%Y-%m-%dT%H:%M:%fZ','now')),
    zaktualizowano           TEXT    NOT NULL DEFAULT (strftime('%Y-%m-%dT%H:%M:%fZ','now'))
);

-- ── Etap budowy produktu — tracker Product Buildera ──────────────────────────
-- Wartości kolumny `stan` są wartościami kontraktu (AppStageStatus).
CREATE TABLE etap_apps (
    id                       INTEGER PRIMARY KEY AUTOINCREMENT,
    identyfikator_zewnetrzny TEXT    NOT NULL UNIQUE,
    okno                     TEXT    NOT NULL,
    nazwa                    TEXT    NOT NULL,
    kolejnosc                INTEGER NOT NULL DEFAULT 0,
    stan                     TEXT    NOT NULL DEFAULT 'pending'
                                     CHECK(stan IN ('pending','active','done','blocked')),
    -- Wykonawca etapu; pusty łańcuch żądania zdejmuje przypisanie, więc kolumna
    -- dopuszcza brak wartości.
    wykonawca                TEXT,
    utworzono                TEXT    NOT NULL DEFAULT (strftime('%Y-%m-%dT%H:%M:%fZ','now')),
    zaktualizowano           TEXT    NOT NULL DEFAULT (strftime('%Y-%m-%dT%H:%M:%fZ','now'))
);
CREATE INDEX idx_etap_apps_okno ON etap_apps(okno, kolejnosc, id);

-- ── Kamień milowy produktu ───────────────────────────────────────────────────
-- Wartości kolumny `stan` są wartościami kontraktu (AppMilestoneStatus).
-- `termin` niesie milisekundy epoki wprost z kontraktu (`AppMilestone.dueAt`),
-- więc kolumna jest typu INTEGER i nie potrzebuje przekładu.
CREATE TABLE kamien_milowy_apps (
    id                       INTEGER PRIMARY KEY AUTOINCREMENT,
    identyfikator_zewnetrzny TEXT    NOT NULL UNIQUE,
    okno                     TEXT    NOT NULL,
    nazwa                    TEXT    NOT NULL,
    termin                   INTEGER,
    stan                     TEXT    NOT NULL DEFAULT 'planned'
                                     CHECK(stan IN ('planned','active','reached','missed')),
    utworzono                TEXT    NOT NULL DEFAULT (strftime('%Y-%m-%dT%H:%M:%fZ','now')),
    zaktualizowano           TEXT    NOT NULL DEFAULT (strftime('%Y-%m-%dT%H:%M:%fZ','now'))
);
CREATE INDEX idx_kamien_milowy_apps_okno ON kamien_milowy_apps(okno, termin, id);

-- ── Etapy składające się na kamień milowy ────────────────────────────────────
CREATE TABLE kamien_milowy_etap_apps (
    kamien_id INTEGER NOT NULL REFERENCES kamien_milowy_apps(id) ON DELETE CASCADE,
    -- Kod zewnętrzny etapu, nie więz obcy — patrz czoło pliku.
    etap_kod  TEXT    NOT NULL,
    PRIMARY KEY (kamien_id, etap_kod)
);
