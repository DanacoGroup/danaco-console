-- Migracja 051 — trwałość modułu Apps: architektura produktu (Architecture
-- Designer), plik warsztatu frontendu/backendu (Workspace okna Apps) oraz
-- dziennik przebiegów wdrożenia (Deployment).
--
-- Warsztat Apps to nie projekt modułu Workspace. `apps.workspace.update` niesie
-- `Layer` (frontend/backend), `Path`, `Content` i opcjonalny `ComponentId` —
-- kształt identyczny z zapisem pliku edytora w module Developer
-- (`developer_wersja_pliku`), nie z `projekt`. `projekt` to byt bez treści:
-- nazwa, opis, stan — kontener, do którego przypięta jest sesja i pamięć
-- kontekstu. Warsztat Apps to sama treść pliku warstwy frontendu/backendu
-- produktu, bez nazwy, bez stanu, bez przypisania eksperta. Współdzielenie jednej
-- tabeli `projekt` znaczyłoby, że wiersz projektu dostaje kolumny treści pliku,
-- których większość wierszy nigdy nie użyje. Dlatego warsztat Apps ma własną
-- tabelę, budowaną na wzór `developer_wersja_pliku`. Commity kodu produktu jadą
-- wspólną komendą `developer.git.action`, poza zasięgiem tej migracji.
--
-- Plik warsztatu trzyma jeden wiersz na (okno, warstwa, ścieżka), nie historię
-- wersji. Kontrakt `AppsWorkspaceUpdateRequest` nie niesie odpowiednika
-- `createVersion` z modułu Developer — nie ma pola, którym Operator zakłada
-- migawkę. `apps.workspace.update` jest więc zwykłym nadpisaniem stanu bieżącego;
-- tabela odzwierciedla to wprost przez UPSERT po kluczu (okno, warstwa, ścieżka),
-- bez tabeli-dziennika obok.
--
-- Architektura trzyma komponenty jako tabelę własną, nie zapis w kolumnie.
-- `AppComponent` niesie pola, po których trzeba filtrować i wiązać przy odczycie:
-- `Kind` (rodzaj komponentu wchodzi do walidacji układu) i `DependsOn` (graf
-- zależności) — gdyby leżały w jednym polu JSON rodzica, walidacja układu
-- (`AppArchitecture.ValidationIssues`) i zapytanie „które komponenty zależą od X"
-- musiałyby parsować JSON przy każdym odczycie. `apps.architecture.define`
-- nadsyła całą listę komponentów na nowo (kontrakt: `Components []AppComponent`,
-- bez trybu częściowej zmiany), więc zapis jest zawsze „usuń komponenty
-- architektury, wstaw przysłane od nowa".
--
-- Zależność między komponentami ma własną tabelę złącznikową, nie kolumnę JSON.
-- `AppComponent.DependsOn []string` jest listą identyfikatorów zewnętrznych
-- innych komponentów tej samej architektury.
--
-- Dziennik wdrożeń powtarza wzór `developer_budowanie`. `apps.deployment.run`
-- i `developer.build.run` mają identyczny kształt zadania: zlecenie z zewnętrznym
-- kodem, oknem, czasem startu i końca, stanem i śladem tekstowym. Rdzeń niczego
-- naprawdę nie wdraża — tabela niesie wyłącznie ślad zlecenia (środowisko,
-- strategia, wersja, notatki, adres po wdrożeniu, odnośnik do logu) i jego stan;
-- sam przebieg prowadzi rdzeń w pamięci. `RollbackToDeploymentId`
-- i `RolledBackFromId` są parą pól tego samego łuku — jedna kolumna samoodwołania
-- wystarcza, bo odczyt idzie zawsze od strony cofnięcia.
--
-- Log wdrożenia jest odwołaniem, nie treścią w bazie. Kontrakt niesie
-- `LogRef *string` — nazwa pola już mówi „odnośnik", nie „treść" — więc kolumna
-- `log_odwolanie` przechowuje ten odnośnik wprost.

-- ── Architektura produktu — Architecture Designer ────────────────────────────
CREATE TABLE architektura_apps (
    id                       INTEGER PRIMARY KEY AUTOINCREMENT,
    identyfikator_zewnetrzny TEXT    NOT NULL UNIQUE,
    okno                     TEXT    NOT NULL,
    nazwa                    TEXT,
    szablon                  TEXT    NOT NULL DEFAULT 'monolith'
                                     CHECK(szablon IN ('monolith','microservices','serverless')),
    -- Wersja rośnie przy każdym zapisie definicji: panel Architecture Designer
    -- musi umieć powiedzieć, którą wersję układu Operator ogląda.
    wersja                   INTEGER NOT NULL DEFAULT 1,
    -- Zastrzeżenia walidacji układu (ValidationIssues []string) są wynikiem
    -- obliczonym przy zapisie, nie wejściem Operatora — jeden wiersz tekstu
    -- rozdzielony znakiem nowej linii wystarcza, bo nikt nie filtruje ani nie
    -- sortuje po pojedynczym zastrzeżeniu.
    zastrzezenia_walidacji   TEXT,
    utworzono                TEXT    NOT NULL DEFAULT (strftime('%Y-%m-%dT%H:%M:%fZ','now')),
    zaktualizowano           TEXT    NOT NULL DEFAULT (strftime('%Y-%m-%dT%H:%M:%fZ','now'))
);
CREATE INDEX idx_architektura_apps_okno ON architektura_apps(okno, zaktualizowano DESC);

-- ── Komponent architektury — Architecture Designer ───────────────────────────
-- Wartości kolumny `rodzaj` są wartościami kontraktu (AppComponentKind), nie
-- ich tłumaczeniem.
CREATE TABLE komponent_architektury_apps (
    id                       INTEGER PRIMARY KEY AUTOINCREMENT,
    architektura_id          INTEGER NOT NULL REFERENCES architektura_apps(id) ON DELETE CASCADE,
    identyfikator_zewnetrzny TEXT    NOT NULL,
    nazwa                    TEXT    NOT NULL,
    rodzaj                   TEXT    NOT NULL
                                     CHECK(rodzaj IN ('frontend','backend','service','database','queue','external')),
    stos                     TEXT,
    opis                     TEXT,
    kontrakt_api             TEXT,
    utworzono                TEXT    NOT NULL DEFAULT (strftime('%Y-%m-%dT%H:%M:%fZ','now')),
    UNIQUE (architektura_id, identyfikator_zewnetrzny)
);
CREATE INDEX idx_komponent_architektury_apps_architektura
    ON komponent_architektury_apps(architektura_id, id);

-- ── Zależność między komponentami — Architecture Designer ───────────────────
-- Więz pierwotny obejmuje parę komponentów, pętla własna odrzucona wprost, cykl
-- dłuższy rozstrzyga walidacja układu w rdzeniu.
CREATE TABLE zaleznosc_komponentu_apps (
    architektura_id INTEGER NOT NULL REFERENCES architektura_apps(id) ON DELETE CASCADE,
    komponent_z     TEXT    NOT NULL,
    komponent_do    TEXT    NOT NULL,
    utworzono       TEXT    NOT NULL DEFAULT (strftime('%Y-%m-%dT%H:%M:%fZ','now')),
    PRIMARY KEY (architektura_id, komponent_z, komponent_do),
    CHECK (komponent_z <> komponent_do)
);
CREATE INDEX idx_zaleznosc_komponentu_apps_do
    ON zaleznosc_komponentu_apps(architektura_id, komponent_do);

-- ── Plik warsztatu — Workspace (okno Apps, nie moduł Workspace) ─────────────
-- Klucz (okno, warstwa, ścieżka) niesie tożsamość pliku; `apps.workspace.update`
-- działa jako UPSERT po tym kluczu, bez tabeli-dziennika wersji — patrz
-- rozstrzygnięcie na czole pliku.
CREATE TABLE plik_warsztatu_apps (
    id             INTEGER PRIMARY KEY AUTOINCREMENT,
    okno           TEXT    NOT NULL,
    warstwa        TEXT    NOT NULL CHECK(warstwa IN ('frontend','backend')),
    sciezka        TEXT    NOT NULL,
    tresc          TEXT    NOT NULL,
    rozmiar        INTEGER NOT NULL DEFAULT 0,
    -- Komponent architektury, którego dotyczy zmiana — wartość danych (kod
    -- zewnętrzny), nie więz obcy: plik może wskazywać komponent usunięty
    -- z architektury po zapisie.
    komponent_id   TEXT,
    utworzono      TEXT    NOT NULL DEFAULT (strftime('%Y-%m-%dT%H:%M:%fZ','now')),
    zaktualizowano TEXT    NOT NULL DEFAULT (strftime('%Y-%m-%dT%H:%M:%fZ','now')),
    UNIQUE (okno, warstwa, sciezka)
);
CREATE INDEX idx_plik_warsztatu_apps_okno ON plik_warsztatu_apps(okno, warstwa, sciezka);

-- ── Dziennik przebiegów wdrożenia — Deployment ───────────────────────────────
-- Kształt zadania jak w `developer_budowanie` — patrz nagłówek pliku. Wartości
-- kolumny `stan` są wartościami kontraktu (AppDeployStatus).
CREATE TABLE wdrozenie_apps (
    id                    INTEGER PRIMARY KEY AUTOINCREMENT,
    kod                   TEXT    NOT NULL UNIQUE,
    okno                  TEXT    NOT NULL,
    srodowisko            TEXT    NOT NULL CHECK(srodowisko IN ('dev','staging','production')),
    strategia             TEXT    NOT NULL DEFAULT 'immediate'
                                  CHECK(strategia IN ('immediate','staged','blueGreen')),
    stan                  TEXT    NOT NULL DEFAULT 'pending'
                                  CHECK(stan IN ('pending','running','succeeded','failed','rolledBack')),
    wersja                TEXT,
    notatki_wydania       TEXT,
    adres                 TEXT,
    log_odwolanie         TEXT,
    -- Wdrożenie źródłowe cofnięcia — patrz rozstrzygnięcie o jednej kolumnie
    -- samoodwołania na czole pliku.
    cofniete_do_kodu      TEXT,
    rozpoczeto             TEXT    NOT NULL DEFAULT (strftime('%Y-%m-%dT%H:%M:%fZ','now')),
    zakonczono             TEXT
);
CREATE INDEX idx_wdrozenie_apps_okno ON wdrozenie_apps(okno, rozpoczeto DESC);
CREATE INDEX idx_wdrozenie_apps_stan ON wdrozenie_apps(stan, rozpoczeto DESC);
