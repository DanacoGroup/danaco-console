-- Migracja 050 — trwałość modułu Assistant: zlecenie wieloetapowe wydane
-- poleceniem głosowym albo tekstowym (Actions Monitor) i chronologiczny
-- dziennik działań asystenta.
--
-- Zlecenie asystenta nie jest pozycją silnika kolejek. Kolejka obsługuje pętlę
-- sesyjną i MultitaskingAI — jej wiersz wisi na `sesja_id` i `okno_komunikacji`,
-- a stan biegnie po ścieżce pracy Wykonawcy nad kartą repozytorium (`oczekuje →
-- przydzielona → wykonywana → do_weryfikacji → ukonczona/bledna/anulowana`,
-- z licznikiem obiegów). Zlecenie asystenta wisi na oknie Assistant
-- (`WindowId`), nie na sesji, i ma własną ścieżkę stanu kontraktu
-- (`queued/running/paused/done/failed/cancelled`) oraz pola bez odpowiednika
-- w kolejce — `currentStep`/`totalSteps` jako para liczb kroku bieżącego wobec
-- łącznego i `priority` sterowalny z Actions Monitor. Stąd osobny byt.
--
-- Zlecenie asystenta nie jest też pozycją katalogu akcji (`core/akcje.go`,
-- komenda `action.list`). Katalog akcji to spis dostępnych czynności — wiersz
-- opisuje przycisk panelu albo narzędzie kanału modelu (`Id`, `Name`, `Command`,
-- `Scope`), czynny niezależnie od tego, czy ktoś go użył. Zlecenie asystenta to
-- wykonywana instancja jednej rozmowy z Operatorem, z historią etapów i wynikiem.
-- Katalog odpowiada na pytanie „co można zrobić”, zlecenie — „co się właśnie
-- dzieje”. Osobne pytania, osobne tabele.
--
-- Czas jest liczbą. `AssistantAction.createdAt/updatedAt` i
-- `AssistantActivityEntry.createdAt` niosą milisekundy epoki — kolumna INTEGER
-- przenosi wartość kontraktu bez przekładu w obie strony.
--
-- Dziennik działań jest osobnym bytem od zlecenia, nie widokiem na nie. Wpis
-- dziennika (`assistant.activity.list`) obejmuje też wpisy własne asystenta
-- (`kind = note`) bez zlecenia w tle — kolumna `zlecenie_kod` jest więc
-- opcjonalna, a nie węzłem, przez który dziennik zawsze musiałby przechodzić.
--
-- Rozpoznawanie mowy nie mieszka w rdzeniu. `assistant.voice.command` przyjmuje
-- `audioRef` (odnośnik do nagranego pliku) albo `transcript` (tekst poprawiony
-- przez Operatora) — rdzeń nie transkrybuje i nie syntezuje mowy sam, więc
-- schemat nie zakłada kolumn na krok, którego nikt nie dostarcza. Treść nagrania
-- zostaje plikiem zewnętrznym; w wierszu dziennika ląduje wyłącznie jego
-- odnośnik, tak samo jak w kontrakcie.

-- ── Zlecenie wieloetapowe asystenta (Actions Monitor) ────────────────────────
-- Wartości kolumn `stan` i `droga` są wartościami kontraktu
-- (AssistantActionStatus, AssistantOrigin) wprost, bez tłumaczenia.
CREATE TABLE zlecenie_asystenta (
    id                       INTEGER PRIMARY KEY AUTOINCREMENT,
    identyfikator_zewnetrzny TEXT    NOT NULL UNIQUE,
    okno_kod                 TEXT    NOT NULL,
    tytul                    TEXT,
    stan                     TEXT    NOT NULL DEFAULT 'queued'
                                     CHECK(stan IN ('queued','running','paused','done','failed','cancelled')),
    droga                    TEXT    NOT NULL DEFAULT 'text'
                                     CHECK(droga IN ('voice','text')),
    etap_biezacy             INTEGER,
    liczba_etapow            INTEGER,
    priorytet                INTEGER,
    wynik                    TEXT,
    utworzono                INTEGER NOT NULL,
    zaktualizowano           INTEGER NOT NULL
);
CREATE INDEX idx_zlecenie_asystenta_okno ON zlecenie_asystenta(okno_kod, zaktualizowano DESC);
CREATE INDEX idx_zlecenie_asystenta_stan ON zlecenie_asystenta(stan, zaktualizowano DESC);

-- ── Dziennik działań asystenta (Actions Monitor / rozmowa) ───────────────────
-- Klucz obcy jest miękki (`ON DELETE SET NULL`): usunięcie zlecenia nie ma
-- prawa wymazać śladu rozmowy, która realnie się odbyła.
CREATE TABLE wpis_dziennika_asystenta (
    id                INTEGER PRIMARY KEY AUTOINCREMENT,
    kod               TEXT    NOT NULL UNIQUE,
    okno_kod          TEXT    NOT NULL,
    zlecenie_kod      TEXT    REFERENCES zlecenie_asystenta(identyfikator_zewnetrzny) ON DELETE SET NULL,
    rodzaj            TEXT    NOT NULL DEFAULT 'note'
                               CHECK(rodzaj IN ('command','result','note')),
    tresc             TEXT    NOT NULL,
    nagranie_odnosnik TEXT,
    utworzono         INTEGER NOT NULL
);
CREATE INDEX idx_wpis_dziennika_asystenta_okno ON wpis_dziennika_asystenta(okno_kod, utworzono DESC);
CREATE INDEX idx_wpis_dziennika_asystenta_zlecenie ON wpis_dziennika_asystenta(zlecenie_kod, utworzono DESC);
