-- Migracja 121 — trzy braki bytów: poziomy pamięci eksperta, widoczność
-- eksperta, źródło rozszerzenia.
--
-- Wszystkie trzy są brakiem bytu, nie brakiem komendy: komendy `agent.create`,
-- `agent.update` i `extension.install` istnieją i działają, ale nie miały gdzie
-- zapisać trzech faktów, których wymaga od nich specyfikacja. Idą w jednej
-- migracji, bo dotyczą jednego zestawu komend; rozbicie na trzy numery nie
-- dołożyłoby ani jednej informacji.
--
-- ── 1. Poziomy pamięci: dlaczego tabela, a nie kolumna ──────────────────────
-- Definicja eksperta wskazuje, z których poziomów pamięci korzysta domyślnie
-- (`globalna | projekt | sesja | środowisko`), albo pamięć jest wyłączona
-- w całości. „Wyłączona" nie jest piątym poziomem i nie wolno jej tak zapisać:
-- piąta wartość obok czterech pozwala ułożyć wiersz sprzeczny — `sesja`
-- i `wyłączona` naraz — który rdzeń musiałby rozstrzygać zgadywaniem. Wyłączenie
-- jest pustym zbiorem poziomów, a jedynym kształtem, który zbiór pusty niesie bez
-- udawania, jest tabela podrzędna: brak wierszy znaczy brak poziomów i nic więcej.
-- Kolumna napisowa z listą po przecinku dawałaby to samo, ale bez więzu CHECK na
-- każdej wartości i bez klucza pilnującego, że poziom się nie powtórzy.
--
-- Wartość wyjściowa to komplet czterech poziomów, nie zbiór pusty: stanem
-- wyjściowym platformy jest pełny dostęp operacyjny, więc ekspert świeżo założony
-- ma pracować, nie prosić o włączenie pamięci. Wyłączenie jest świadomą decyzją
-- Operatora, więc to ono wymaga czynności. Wartość ta jest wyjściowa, nie
-- ostateczna: przypisanie eksperta do projektu i do roli nadpisują ją, bo
-- pierwszeństwo ma zasięg najbardziej szczegółowy. Ta tabela trzyma wyłącznie
-- wartość z definicji eksperta; nadpisania mają własne poziomy zasięgu.
--
-- ── 2. Widoczność: dlaczego bez wskazania projektu ──────────────────────────
-- Widoczność ma dwie wartości (`globalny | projektowy`). Dołożenie obok nich
-- kolumny `projekt_id` byłoby drugą prawdą o przynależności eksperta do projektu
-- obok tabeli `przypisanie_agenta_projektu` (migracja 035), którą wypełnia
-- `workspace.agent.assign`; dwie prawdy rozjechałyby się przy pierwszym
-- przypisaniu zrobionym drugą drogą. Widoczność mówi więc tylko to, czego tamta
-- tabela nie mówi: czy ekspert pokazuje się wszędzie, czy wyłącznie tam, gdzie
-- został przypisany.
--
-- ── 3. Źródło rozszerzenia: jeden wyjątek, reszta wspólna ───────────────────
-- Rdzeń serwera nie rozróżnia źródła rozszerzenia — obowiązuje wspólny kontrakt
-- integracji. Ładowanie, wywołanie i użycie operacyjne są dla `danaco`
-- i `personal` te same, i ta migracja nie zakłada bytu, który by je rozdzielał.
-- Jedyny wyjątek to stan wyjściowy przy rejestracji: Danaco Plugin staje włączone,
-- bo zestaw wbudowany jest częścią funkcjonalności bazowej; Personal staje
-- wyłączone, bo Operator włącza je świadomie, po przejrzeniu konfiguracji
-- i zakresu. Rozstrzyga to warstwa instalacji, nie ta tabela — kolumna niesie sam
-- fakt pochodzenia. Wiersze zastane dostają `personal`, bo wszystkie powstały
-- wywołaniem `extension.install` przez Operatora; wpisanie im `danaco` byłoby
-- ogłoszeniem, że coś jest częścią pakietu serwera, choć nikt tego nie dostarczył.

-- ── 1. Poziomy pamięci eksperta ─────────────────────────────────────────────
-- Kolumna `poziom` trzyma wartość kontraktu wprost (shared.MemoryLevel), bez
-- przekładu — tak jak `tryb_nakladki` z migracji 081. CHECK wylicza cztery
-- poziomy i ani jednego więcej; wartości „wyłączona" tu nie ma i nie może być,
-- bo wyłączenie to brak wiersza.
CREATE TABLE agent_pamiec_poziom (
    agent_id INTEGER NOT NULL REFERENCES agent(id) ON DELETE CASCADE,
    poziom   TEXT    NOT NULL
             CHECK (poziom IN ('global', 'project', 'session', 'environment')),

    PRIMARY KEY (agent_id, poziom)
);

-- Eksperci zastani zachowują stan wyjściowy platformy: pełny komplet poziomów.
-- Bez tego wypełnienia biblioteka zastana pokazałaby się jako „pamięć
-- wyłączona" — czyli migracja zmieniłaby zachowanie ekspertów, których nikt
-- nie tknął.
INSERT INTO agent_pamiec_poziom (agent_id, poziom)
SELECT a.id, p.poziom
  FROM agent a
  JOIN (SELECT 'global' AS poziom UNION ALL SELECT 'project'
        UNION ALL SELECT 'session' UNION ALL SELECT 'environment') p;

-- ── 2. Widoczność eksperta ──────────────────────────────────────────────────
-- Wartość domyślna `global` odpowiada stanowi sprzed tej migracji: ekspert był
-- widoczny wszędzie, bo nie było czym go zawęzić. Eksperci zastani nie zmieniają
-- więc widoczności ani o jotę.
ALTER TABLE agent ADD COLUMN widocznosc TEXT NOT NULL DEFAULT 'global'
    CHECK (widocznosc IN ('global', 'project'));

-- ── 3. Źródło rozszerzenia ──────────────────────────────────────────────────
-- Kolumna niesie wartość kontraktu wprost (shared.ExtensionOrigin). Nie jest
-- przesłanką ładowania — patrz nagłówek; rozstrzyga wyłącznie stan wyjściowy
-- przy rejestracji, a to dzieje się piętro wyżej, w adapterze instalacji.
ALTER TABLE rozszerzenie ADD COLUMN zrodlo_pochodzenia TEXT NOT NULL DEFAULT 'personal'
    CHECK (zrodlo_pochodzenia IN ('danaco', 'personal'));
