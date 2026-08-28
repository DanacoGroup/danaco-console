-- Migracja 073 — tożsamość własna eksperta, warstwy jego promptu i wtyczki.
--
-- Czego brakowało. Ekspert miał `kod`, `nazwa`, `opis` i jedno pole na prompt —
-- `instrukcje_systemowe`. Portfolio ekspertów potrzebuje trzech rzeczy, których
-- w schemacie nie ma:
--   · imienia własnego — `nazwa` jest napisem technicznym, po którym idzie
--     sortowanie wykazu (`idx_agent_nazwa`, `dane/agenci.go:listaAgentow`);
--   · znaku graficznego — kolumny na favikon nie ma nigdzie w bazie;
--   · warstw promptu — `injection/nakladka.go` składa prompt z warstw
--     (konstytucja · profil · ekspertyza), ale ekspert nie ma gdzie żadnej
--     z nich zapisać; miał jeden worek tekstu.
-- Kontrakt (`shared/contract.go`) niesie już struktury `AgentLayer`
-- i `AgentPlugin`, pola `Agent.displayName`, `Agent.favicon`, `Agent.layers`,
-- `Agent.pluginIds` oraz cztery komendy `agent.layer.*` / `agent.plugin.*`.
-- Bez tej migracji byłyby to komendy bez miejsca zapisu.
--
-- Czego świadomie nie ma.
--   · Nie ma katalogu dostępnych wtyczek. Kontrakt zna `agent.plugin.add`
--     z nazwą, źródłem i wersją podanymi wprost — komendy przeglądania katalogu
--     rozszerzeń w rodzinie `agent.*` nie ma, więc tabela katalogu byłaby
--     zapisem bez czytelnika.
--   · Nie ma archiwum treści warstw. `agent.wersja` jest licznikiem, nie
--     archiwum, i ta migracja tego nie zmienia.
--   · Nie ma kolumny na obraz favikonu. `favikon` niesie odwołanie — napis,
--     który klient umie pokazać (nazwa znaku, ścieżka, identyfikator zasobu).
--     Bajtów obrazu baza nie trzyma, tak samo jak biblioteka trzyma na dysku
--     treść pliku, a w kolumnie wyłącznie odwołanie.
--
-- Co się nie zmienia. Migracja wyłącznie dodaje: dwie kolumny z wartością
-- domyślną i dwie tabele. Nie ma tu ani jednego `UPDATE`, ani jednego `DROP`,
-- ani jednej zmiany istniejącej kolumny. `nazwa`, `instrukcje_systemowe`,
-- `agent_umiejetnosc`, `agent_konektor` i `agent_uprawnienie` zostają takie,
-- jakie były.

-- ── (a) Tożsamość własna eksperta ─────────────────────────────────────────────
--
-- Imię własne nie jest drugą nazwą. `nazwa` zostaje tym, czym była: napisem
-- technicznym, po którym biegnie porządek wykazu i wyszukiwanie frazą
-- (`dane/agenci.go`). `imie_wlasne` jest tym, co widzi Operator — i tylko tym.
-- Dwie kolumny, dwa zastosowania; gdyby imię wpisać w `nazwa`, zmiana imienia
-- przestawiałaby wykaz i rozjeżdżała wyszukiwanie.
--
-- Pusty napis znaczy „nie nadano", nie „błąd". Obie kolumny są
-- NOT NULL DEFAULT '', więc każdy ekspert założony wcześniej dostaje je puste
-- i pracuje dalej bez żadnej zmiany. Klient pokazuje wtedy `nazwa` i znak
-- zastępczy — pola kontraktu `displayName` i `favicon` są opcjonalne właśnie
-- po to.
ALTER TABLE agent ADD COLUMN imie_wlasne TEXT NOT NULL DEFAULT '';
ALTER TABLE agent ADD COLUMN favikon     TEXT NOT NULL DEFAULT '';

-- ── (b) Warstwy promptu eksperta ──────────────────────────────────────────────
--
-- Dlaczego tabela, a nie trzy kolumny. Warstw jest dziś trzy, ale liczba trzy
-- nie jest nigdzie przesądzona. `injection/nakladka.go` dopuszcza warstwę
-- o nazwie spoza katalogu: `kolejnoscWarstw` zna trzy nazwy, a `pozycjaWarstwy`
-- zwraca dla każdej innej `len(kolejnoscWarstw) + 1`, czyli ustawia ją na końcu
-- zamiast odrzucić. Trzy kolumny `tresc_konstytucja`, `tresc_profil`,
-- `tresc_ekspertyza` zamknęłyby ekspertowi drogę, którą silnik nakładki ma
-- otwartą, a czwarta warstwa wymagałaby zmiany schematu.
--
-- Klucz na parze — jedna treść na warstwę. Klucz główny na parze sprawia, że
-- powtórzone `agent.layer.set` nadpisuje treść zamiast dokładać drugi wiersz tej
-- samej warstwy. Zapis warstwy jest więc ustaleniem stanu, a nie dopisaniem
-- zdarzenia.
--
-- `instrukcje_systemowe` zostaje nietknięte. Kolumna jest w kontrakcie polem
-- `Agent.systemPrompt` i czyta ją klient — kasując ją, migracja zabrałaby oknu
-- modułu Agents treść, którą ono dziś pokazuje. Ta migracja niczego nie kasuje
-- i niczego nie przenosi: warstwy są bytem nowym, obok istniejącego pola.
--
-- Wartości wprost z kontraktu. `warstwa` przyjmuje dosłownie wartości
-- `shared.IdentityLayer` ('constitution','profile','expertise'), a `tryb` —
-- `shared.IdentityMode` ('ZASTAP','DOLACZ'). Kontrakt nie daje dla tych dwóch
-- wyliczeń słownika przekładu bazy, więc kolumny trzymają wartość kontraktu
-- wprost — tak samo jak `agent.transport` i `agent_konektor.rodzaj`. Warunki
-- CHECK są jedynym miejscem, w którym ten katalog stoi po stronie bazy;
-- warstwa `dane` go nie powtarza.
CREATE TABLE agent_warstwa (
    agent_id       INTEGER NOT NULL REFERENCES agent(id) ON DELETE CASCADE,
    warstwa        TEXT    NOT NULL CHECK (warstwa IN ('constitution', 'profile', 'expertise')),
    tresc          TEXT    NOT NULL DEFAULT '',
    tryb           TEXT    NOT NULL DEFAULT 'DOLACZ' CHECK (tryb IN ('ZASTAP', 'DOLACZ')),
    -- Warstwa wyłączona zostaje zapisana i widoczna w oknie; nie wchodzi
    -- wyłącznie do złożonego promptu. Wyłączenie nie jest kasowaniem.
    aktywna        INTEGER NOT NULL DEFAULT 1 CHECK (aktywna IN (0, 1)),
    zaktualizowano TEXT    NOT NULL DEFAULT (strftime('%Y-%m-%dT%H:%M:%fZ', 'now')),
    PRIMARY KEY (agent_id, warstwa)
);

-- ── (c) Wtyczki eksperta ──────────────────────────────────────────────────────
--
-- Dlaczego osobno od `agent_konektor`. Konektor jest drogą do usługi: wiersz
-- `agent_konektor` wskazuje most z katalogu punktów dostępu, z którego rdzeń
-- składa wpisy `mcpServers` podawane procesowi modelu przełącznikiem
-- `--mcp-config` (`core/most_okna.go`). Wtyczka jest katalogiem rozszerzeń
-- powłoki — nie ma adresu, nie ma poświadczenia, nie ma punktu dostępu; ma
-- nazwę, źródło i wersję, a program dostaje ją przełącznikiem `--plugin-dir`.
-- To dwa różne przełączniki i dwa różne byty: skille, konektory i pluginy są
-- trzema bytami, nie dwoma.
--
-- Wartość 'plugin' w `agent_konektor.rodzaj` zostaje. Taki był wcześniej jedyny
-- sposób zapisania czegokolwiek o wtyczce. Tej wartości nie kasujemy i nie
-- przepisujemy wierszy — kasowanie rodzaju wywróciłoby
-- `core/adapter_modul_agents_zasoby.go` i zabrałoby treść konektorom już
-- zapisanym.
--
-- Kod jest identyfikatorem trwałym. Tak samo jak w `agent_konektor`: `id` służy
-- powiązaniom w bazie, a `kod` wychodzi na zewnątrz jako `AgentPlugin.id`
-- kontraktu. Dzięki temu `agent.plugin.remove` wskazuje wtyczkę kodem, a nie
-- numerem wiersza, którego kontrakt nie zna.
--
-- Bez klucza na parze (agent_id, nazwa). Ekspert może mieć dwie wtyczki tej
-- samej nazwy w różnych wersjach albo z różnych źródeł, a kontrakt kasuje
-- wtyczkę po `pluginId`, nie po nazwie — klucz na nazwie odbierałby tę
-- możliwość bez powodu.
CREATE TABLE agent_wtyczka (
    id        INTEGER PRIMARY KEY AUTOINCREMENT,
    kod       TEXT    NOT NULL UNIQUE,
    agent_id  INTEGER NOT NULL REFERENCES agent(id) ON DELETE CASCADE,
    nazwa     TEXT    NOT NULL,
    -- Skąd wtyczka pochodzi (katalog na dysku, adres repozytorium, nazwa
    -- rejestru). Puste znaczy „nie podano", nie „brak".
    zrodlo    TEXT,
    wersja    TEXT,
    aktywna   INTEGER NOT NULL DEFAULT 1 CHECK (aktywna IN (0, 1)),
    utworzono TEXT    NOT NULL DEFAULT (strftime('%Y-%m-%dT%H:%M:%fZ', 'now'))
);

-- Wykaz wtyczek eksperta idzie po nazwie, tak samo jak wykaz konektorów
-- (`idx_agent_konektor_agent`).
CREATE INDEX idx_agent_wtyczka_agent ON agent_wtyczka(agent_id, nazwa);
