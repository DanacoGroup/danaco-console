-- Migracja 369 — tożsamość wykonawcy przy zmianie śledzonej i przy komentarzu.
--
-- ── Dlaczego rodzaj autora nie wystarcza ─────────────────────────────────────
-- `zmiana_sledzona_studio.autor` i `komentarz_studio.autor` mają warunek
-- CHECK(autor IN ('uzytkownik','model')). To rozróżnienie działa i zostaje.
-- Nie wystarcza jednak przy dwóch wykonawcach pracujących naraz nad jednym
-- dokumentem: Operator nie odróżni, który z nich co napisał, a przełącznik
-- podświetlający zmiany wykonawców pokazywałby obu jako jednego.
--
-- ── Dlaczego to NIE jest nowa wartość w warunku CHECK ───────────────────────
-- Słowami Właściciela: w module Agents można założyć setki własnych agentów
-- i później ich powoływać. Wyliczenie ich nigdy nie obejmie — dopisanie
-- wartości „agent-umowy" do warunku CHECK kazałoby zmieniać schemat przy każdym
-- agencie założonym przez Operatora, a warunku CHECK w SQLite nie zmienia się
-- bez przepisania tabeli. Tożsamością jest więc KOD AGENTA w osobnej kolumnie,
-- a rodzaj autora zostaje grubym rozróżnieniem człowiek-wykonawca.
--
-- ── Dlaczego pojęcia agenta nie zakładamy drugiego ──────────────────────────
-- Platforma identyfikuje wykonawcę dwoma bytami, które już są: `Agent`
-- (tabela ekspertów, `id`, `name`, wersja) i `Subagent` (wykonawca uruchomiony
-- przez okno). `autor_agent_kod` odpowiada pierwszemu, `autor_podagent_kod`
-- drugiemu. Trzeciego pojęcia agenta w Studiu nie ma — byłby to ten sam błąd,
-- co dwa wykazy nośników druku.
--
-- Klucza obcego do tabeli ekspertów tu NIE ma i to jest zamysł: zmiana wniesiona
-- przez agenta ma zostać podpisana także wtedy, gdy Operator tego agenta później
-- usunie. Podpis pod pracą przetrwa wykonawcę; inaczej usunięcie eksperta
-- zamieniłoby historię dokumentu w anonim.
--
-- ── Dlaczego wszystkie kolumny są nieobowiązkowe ────────────────────────────
-- Zmiany, komentarze i adnotacje zapisane przed tą dobudową agenta nie mają
-- i mają zostać poprawne. Pusty kod agenta znaczy „wykonawca nienazwany albo
-- Operator" — brak wiedzy, nie twierdzenie.

ALTER TABLE zmiana_sledzona_studio ADD COLUMN autor_agent_kod TEXT;
ALTER TABLE zmiana_sledzona_studio ADD COLUMN autor_agent_nazwa TEXT;
ALTER TABLE zmiana_sledzona_studio ADD COLUMN autor_agent_wersja TEXT;
ALTER TABLE zmiana_sledzona_studio ADD COLUMN autor_podagent_kod TEXT;

-- Zmiana POSTACI, nie treści, ma być widoczna w podświetleniu tak samo jak
-- dopisany akapit — „zmiana kroju bez zmiany liter nie może być niewidzialna".
-- Rodzaj `formatowanie` już stoi w warunku CHECK tej tabeli, ale nie było czym
-- powiedzieć, CO się w postaci zmieniło: `tresc_przed` i `tresc_po` niosą litery,
-- a przy zmianie postaci litery są te same. Stąd wycinek postaci sprzed i po.
ALTER TABLE zmiana_sledzona_studio ADD COLUMN postac_przed_json TEXT;
ALTER TABLE zmiana_sledzona_studio ADD COLUMN postac_po_json TEXT;

-- Czynność dziennika, która zmianę odłożyła. Wiąże dwie drogi cofania: zmianę
-- śledzoną (rozstrzyganą decyzją) i czynność dokumentu (cofaną z dziennika).
-- Bez tego wiązania cofnięcie czynności zostawiłoby zmianę śledzoną wiszącą
-- w powietrzu — Operator widziałby do rozstrzygnięcia zmianę, której już nie ma.
ALTER TABLE zmiana_sledzona_studio ADD COLUMN czynnosc_kod TEXT;

CREATE INDEX idx_zmiana_sledzona_studio_agent
    ON zmiana_sledzona_studio(dokument_id, autor_agent_kod, zakres_od);
-- Podświetlenie zmian wykonawców pyta o wszystkie zmiany autora `model`
-- w całym dokumencie, od pierwszej strony do ostatniej, i liczy je do licznika
-- przy przełączniku. Indeks zastany prowadzi po decyzji, nie po autorze.
CREATE INDEX idx_zmiana_sledzona_studio_autor
    ON zmiana_sledzona_studio(dokument_id, autor, decyzja, zakres_od);

ALTER TABLE komentarz_studio ADD COLUMN autor_agent_kod TEXT;
ALTER TABLE komentarz_studio ADD COLUMN autor_agent_nazwa TEXT;
ALTER TABLE komentarz_studio ADD COLUMN autor_agent_wersja TEXT;
ALTER TABLE komentarz_studio ADD COLUMN autor_podagent_kod TEXT;

CREATE INDEX idx_komentarz_studio_agent
    ON komentarz_studio(dokument_id, autor_agent_kod, rodzaj);
