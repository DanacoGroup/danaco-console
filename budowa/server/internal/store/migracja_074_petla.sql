-- Migracja 074 — pętla wykonawcza automatyki: obsada, wiązanie z agentem, wybudzenie.
--
-- Podstawa (definicja automatyki, kroki, zależności, harmonogram z cronem
-- i strefą czasową, wyzwalacze, przebiegi ze śladem i ponowieniami) już stoi.
-- Brakuje trzech rzeczy:
--
--   1. krok rodzaju `model` wskazuje model, ale nie umie wskazać agenta
--      z portfolio;
--   2. w pętli bierze udział od 1 do 4 modeli z podziałem na role — obsady nie
--      ma gdzie zapisać;
--   3. wybudzenie: bieg zawieszony po wykonaniu akcji i czekający na reakcję
--      ze świata.
--
-- Pauza to nie jest wybudzenie — i na tym rozróżnieniu stoi cała ta migracja.
-- Pauza jest już w schemacie: krok rodzaju `wait` czeka określony czas i idzie
-- dalej, a stan `paused` czeka na Operatora. Wybudzenie jest czym innym:
--
--   pauza        bieg czeka na czas albo na Operatora — wiadomo, kiedy ruszy
--   wybudzenie   bieg czeka na świat — na reakcję, która może przyjść za godzinę,
--                za trzy dni albo nigdy
--
-- Z tego biorą się trzy wymagania, żadne nie jest zwykłym polem:
--
--   · bieg musi przetrwać restart rdzenia. Zawieszenie na dni oznacza, że stan
--     żyje w bazie, nie w pamięci procesu. Stąd stan `oczekuje` w tabeli, a nie
--     mapa w pamięci adaptera.
--   · musi istnieć wiązanie: który bieg czeka na który sygnał.
--     `wyzwalacz_automatyki` uruchamia automatykę od początku; wybudzenie trafia
--     do konkretnego biegu w konkretnym kroku. To inny byt i inna tabela.
--   · musi być rozstrzygnięte, co się dzieje, gdy reakcja nie przyjdzie. Bieg
--     czekający wiecznie to wyciek tej samej klasy co podagent-sierota. Stąd
--     `termin` i `po_terminie` jako pola obowiązkowe, nie ozdoba.
--
-- Zero bramek. Żaden z tych bytów niczego Operatorowi nie zabrania. Obsada
-- opisuje, kto bierze udział w pętli; wybudzenie opisuje, na co bieg czeka.
-- Brak wiersza obsady nie blokuje uruchomienia — automatyka bez obsady biegnie
-- na modelu wskazanym w kroku. Brak wiersza oczekiwania znaczy tyle, że bieg
-- na nic nie czeka.
--
-- Słownik ról nie jest wymyślony. Wartości `koordynator` · `wykonawca` ·
-- `samodzielne` to dosłownie wartości bazy wyliczenia `WindowRole` z kontraktu
-- (pole `baza` przy `coordinator`/`executor`/`standalone`). Czwarte miejsce
-- obsady bierze rolę `samodzielne`, tak jak rozstrzyga to klient
-- w `okna-rownolegle/role-domyslne.ts`. Własne nazwy dałyby drugi słownik ról
-- obok kontraktowego, a dwa słowniki tego samego bytu zawsze się rozjeżdżają.
--
-- Czego tu nie ma: komunikacji między modelami. Wymaga nowej rodziny komend albo
-- rozszerzenia `subagent.*`, czyli nowej powierzchni kontraktu. Obsada zapisuje,
-- kto bierze udział; czym ci uczestnicy się wymieniają, rozstrzyga osobny pakiet.

-- ── 1. Krok automatyki wskazuje agenta z portfolio ──────────────────────────
--
-- Dlaczego to nie jest kolumna na `krok_automatyki`. Kolumna w tym miejscu
-- gubiłaby wiązanie po cichu przy pierwszym zapisie definicji. `ZapiszKroki`
-- (`dane/automations_kroki.go`) podmienia komplet kroków: kasuje wszystkie
-- wiersze automatyki i wstawia je na nowo, wymieniając osiem kolumn, których
-- nazwy zna. Kolumny dziewiątej nie zna i nie przepisze, więc każdy zapis
-- z Workflow Buildera zdejmowałby Operatorowi agentów z kroków — bez komunikatu,
-- bez błędu, bez śladu.
--
-- Schemat ma już na to własną odpowiedź: `zaleznosc_kroku_automatyki` wiąże
-- kroki przez `krok_z`/`krok_do` typu TEXT, czyli przez identyfikator zewnętrzny,
-- a nie przez klucz liczbowy — dokładnie dlatego, że wiersze kroków bywają
-- podmieniane w całości. Wiązanie z agentem jest tym samym rodzajem faktu i idzie
-- tą samą drogą.
--
-- `ON DELETE CASCADE` po stronie agenta jest tu poprawne, a nie surowe: wiersz
-- wiązania bez agenta nie znaczy nic. Krok zostaje nietknięty i wraca do
-- zachowania sprzed wiązania — biegnie na modelu wskazanym w `parametry`.
-- Kaskada zabiera tu wyłącznie wiązanie, nigdy kawałek pętli.
--
-- Krok skasowany z definicji zostawia wiersz, który do niczego nie pasuje.
-- To jest w porządku: wiersz bez kroku nigdy nie zostanie odczytany, a przy
-- skasowaniu automatyki znika kaskadą. Sprzątanie osieroconych wiązań nie należy
-- do schematu — należałoby do zapisu kroków, a ten pliku nie zna.

CREATE TABLE agent_kroku_automatyki (
    automatyka_id   INTEGER NOT NULL REFERENCES automatyka(id) ON DELETE CASCADE,
    -- Identyfikator ZEWNĘTRZNY kroku — ten sam rodzaj klucza, którym posługuje
    -- się `zaleznosc_kroku_automatyki`, i z tego samego powodu.
    krok_zewnetrzny TEXT    NOT NULL,
    agent_id        INTEGER NOT NULL REFERENCES agent(id) ON DELETE CASCADE,
    utworzono       TEXT    NOT NULL DEFAULT (strftime('%Y-%m-%dT%H:%M:%fZ','now')),
    PRIMARY KEY (automatyka_id, krok_zewnetrzny)
);

CREATE INDEX idx_agent_kroku_automatyki_agent
    ON agent_kroku_automatyki(agent_id);

-- ── 2. Obsada automatyki — od 1 do 4 uczestników z rolami ───────────────────
--
-- Granica czterech miejsc jest specyfikacją, nie bramką. Klient liczy wykonawców
-- do dwóch (`LICZBA_WYKONAWCOW`) i szuka analityka wśród okien samodzielnych.
-- Granica siedzi w `CHECK(miejsce BETWEEN 1 AND 4)` razem z `UNIQUE(automatyka_id,
-- miejsce)`: cztery miejsca to najwyżej cztery wiersze, bez licznika w kodzie
-- i bez wyzwalacza. Gdy specyfikacja się zmieni, zmienia się jedna liczba
-- w jednym więzie.
--
-- Uczestnik to agent albo model, nigdy nic. `CHECK(agent_id IS NOT NULL OR
-- model IS NOT NULL)` nie pozwala zapisać miejsca obsady, które nie wskazuje
-- na nikogo. Puste miejsce obsady wyglądałoby w wykazie tak samo jak miejsce
-- wypełnione, a nie byłoby nim.

CREATE TABLE obsada_automatyki (
    id            INTEGER PRIMARY KEY AUTOINCREMENT,
    automatyka_id INTEGER NOT NULL REFERENCES automatyka(id) ON DELETE CASCADE,
    -- Uczestnik: agent z portfolio albo goły model. Agent ma pierwszeństwo —
    -- niesie ze sobą instrukcje systemowe, umiejętności i parametry.
    agent_id      INTEGER          REFERENCES agent(id) ON DELETE SET NULL,
    model         TEXT,
    rola          TEXT    NOT NULL DEFAULT 'wykonawca'
                          CHECK(rola IN ('koordynator','wykonawca','samodzielne')),
    miejsce       INTEGER NOT NULL CHECK(miejsce BETWEEN 1 AND 4),
    utworzono     TEXT    NOT NULL DEFAULT (strftime('%Y-%m-%dT%H:%M:%fZ','now')),
    zaktualizowano TEXT   NOT NULL DEFAULT (strftime('%Y-%m-%dT%H:%M:%fZ','now')),
    UNIQUE (automatyka_id, miejsce),
    CHECK (agent_id IS NOT NULL OR model IS NOT NULL)
);

CREATE INDEX idx_obsada_automatyki_automatyka
    ON obsada_automatyki(automatyka_id, miejsce);

CREATE INDEX idx_obsada_automatyki_agent
    ON obsada_automatyki(agent_id)
    WHERE agent_id IS NOT NULL;

-- ── 3. Przebudowa `przebieg_automatyki` — stan `oczekuje` ───────────────────
--
-- SQLite nie zna polecenia zdejmującego więz CHECK, więc tabela idzie przez
-- przebudowę: kopia z poprawionym więzem → przepisanie wierszy → podmiana nazwy
-- → odtworzenie indeksów.
--
-- Przebudowa obejmuje jedną tabelę: na `przebieg_automatyki` nie wskazuje żaden
-- klucz obcy, więc kaskada nie ma czego zabrać, a samo `DROP TABLE` nie zdejmuje
-- dzieci.
--
-- Stan `oczekuje` dokłada się do sześciu istniejących, żadnego nie zabierając.
-- Wiersze istniejące przepisują się co do kolumny, więc żaden bieg zapisany
-- przed tą migracją nie zmienia stanu ani nie znika.
--
-- Kontrakt tego stanu jeszcze nie niesie. `AutomationExecutionStatus` zna sześć
-- wartości (pending · running · paused · succeeded · failed · stopped). Dopóki
-- kontrakt nie dostanie siódmej, stan `oczekuje` żyje w bazie i widzi go rdzeń,
-- ale nie ma jak dojechać do klienta — schemat idzie pierwszy, bo bez miejsca
-- w bazie bieg czekający nie przetrwa restartu.

CREATE TABLE przebieg_automatyki_nowa (
    id                       INTEGER PRIMARY KEY AUTOINCREMENT,
    identyfikator_zewnetrzny TEXT    NOT NULL UNIQUE,
    automatyka_id            INTEGER NOT NULL REFERENCES automatyka(id) ON DELETE CASCADE,
    kolejka_id               INTEGER REFERENCES kolejka(id) ON DELETE SET NULL,
    stan                     TEXT    NOT NULL DEFAULT 'pending'
                                     CHECK(stan IN ('pending','running','paused','oczekuje',
                                                    'succeeded','failed','stopped')),
    etap_biezacy             INTEGER NOT NULL DEFAULT 0,
    etapow                   INTEGER NOT NULL DEFAULT 0,
    proba                    INTEGER NOT NULL DEFAULT 0,
    komunikat_bledu          TEXT,
    rozpoczeto               TEXT    NOT NULL DEFAULT (strftime('%Y-%m-%dT%H:%M:%fZ','now')),
    zakonczono               TEXT
);

INSERT INTO przebieg_automatyki_nowa
    (id, identyfikator_zewnetrzny, automatyka_id, kolejka_id, stan,
     etap_biezacy, etapow, proba, komunikat_bledu, rozpoczeto, zakonczono)
SELECT
     id, identyfikator_zewnetrzny, automatyka_id, kolejka_id, stan,
     etap_biezacy, etapow, proba, komunikat_bledu, rozpoczeto, zakonczono
  FROM przebieg_automatyki;

DROP TABLE przebieg_automatyki;

ALTER TABLE przebieg_automatyki_nowa RENAME TO przebieg_automatyki;

-- Indeksy migracji 040 odtworzone pod nazwami własnymi — indeksy tabeli
-- skasowanej giną razem z nią, więc nazwy są wolne.
CREATE INDEX idx_przebieg_automatyki_automatyka
    ON przebieg_automatyki(automatyka_id, id DESC);

CREATE UNIQUE INDEX idx_przebieg_automatyki_kolejka
    ON przebieg_automatyki(kolejka_id)
    WHERE kolejka_id IS NOT NULL;

-- ── 4. Oczekiwanie biegu — na co czeka i co, gdy nie doczeka ────────────────
--
-- To nie jest drugi wyzwalacz. `wyzwalacz_automatyki` wiąże się z harmonogramem
-- i uruchamia automatykę od początku. Ten byt wiąże się z przebiegiem i wznawia
-- go w miejscu, w którym stanął. Dwa różne pytania, dwie różne tabele:
--
--   wyzwalacz          „kiedy zacząć nowy bieg"
--   oczekiwanie_biegu  „na co czeka ten bieg i od którego kroku ma ruszyć dalej"
--
-- `po_terminie` jest polem obowiązkowym: bieg czekający wiecznie to wyciek. Trzy
-- wyjścia i każde jest decyzją, nie awarią:
--
--   wznow      termin minął — ruszaj dalej tak, jakby sygnał przyszedł.
--              Domyślne, bo produkt ma pracować dalej, a nie stawać.
--   ponow      wykonaj krok oczekiwania jeszcze raz (`proba` rośnie).
--   przerwij   zakończ bieg stanem `stopped` z jawnym powodem.
--
-- Bez terminu też wolno. `termin IS NULL` znaczy „czekaj bez końca" i jest
-- wyborem Operatora, nie przeoczeniem. Wiersz bez terminu nie trafia do budzika,
-- więc nie kosztuje nic.
--
-- Jedno czynne oczekiwanie na bieg. Wymusza to indeks częściowy UNIQUE po
-- `przebieg_id` z warunkiem `wybudzono IS NULL`. Bieg stoi w jednym miejscu,
-- więc nie może czekać na dwie rzeczy naraz; oczekiwania zamknięte zostają
-- w tabeli jako ślad i nie wchodzą sobie w drogę.

CREATE TABLE oczekiwanie_biegu (
    id               INTEGER PRIMARY KEY AUTOINCREMENT,
    przebieg_id      INTEGER NOT NULL REFERENCES przebieg_automatyki(id) ON DELETE CASCADE,
    -- Krok, w którym bieg stanął — identyfikator ZEWNĘTRZNY, tak samo jak
    -- w `zaleznosc_kroku_automatyki`, bo definicja kroków bywa podmieniana
    -- w całości przy zapisie i klucz liczbowy by tego nie przeżył.
    krok_zewnetrzny  TEXT    NOT NULL,
    -- Na co bieg czeka. Kod sygnału ustala ten, kto sygnał wyśle — rdzeń nie
    -- prowadzi tu słownika, bo świat zewnętrzny nie da się zamknąć w CHECK.
    sygnal           TEXT    NOT NULL,
    -- Do kiedy; NULL = bez terminu.
    termin           TEXT,
    po_terminie      TEXT    NOT NULL DEFAULT 'wznow'
                             CHECK(po_terminie IN ('wznow','ponow','przerwij')),
    zalozono         TEXT    NOT NULL DEFAULT (strftime('%Y-%m-%dT%H:%M:%fZ','now')),
    -- Wypełnione = oczekiwanie zamknięte. NULL = bieg nadal czeka.
    wybudzono        TEXT,
    powod_wybudzenia TEXT    CHECK(powod_wybudzenia IS NULL
                                   OR powod_wybudzenia IN ('sygnal','termin','operator')),
    -- Treść, z którą przyszedł sygnał — wynik reakcji świata, do wykorzystania
    -- przez krok wznawiany. Pusta jest w porządku.
    tresc_sygnalu    TEXT,
    CHECK (wybudzono IS NULL OR powod_wybudzenia IS NOT NULL)
);

-- Jedno czynne oczekiwanie na bieg.
CREATE UNIQUE INDEX idx_oczekiwanie_biegu_czynne
    ON oczekiwanie_biegu(przebieg_id)
    WHERE wybudzono IS NULL;

-- Doręczenie sygnału: „które biegi czekają na TEN sygnał".
CREATE INDEX idx_oczekiwanie_biegu_sygnal
    ON oczekiwanie_biegu(sygnal)
    WHERE wybudzono IS NULL;

-- Budzik terminów: „którym oczekiwaniom minął termin". Wiersze bez terminu
-- zostają poza indeksem, więc czekanie bez zegara nic nie kosztuje.
CREATE INDEX idx_oczekiwanie_biegu_termin
    ON oczekiwanie_biegu(termin)
    WHERE wybudzono IS NULL AND termin IS NOT NULL;
