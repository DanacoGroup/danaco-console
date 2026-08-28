-- Migracja 346 — łańcuch edycji zasobu obrazowego modułu Design (rodzina
-- `design.photo.*`).
--
-- ── Dlaczego łańcuch NIE potrzebuje nowego bytu na zasób ────────────────────
-- `zasob_design.wariant_zasobu_id` już istnieje (migracja 048), więc każda
-- obróbka zakłada WARIANT źródła i oryginał zostaje nietknięty. Łańcuch edycji
-- da się więc przejść samymi zasobami. Czego zasoby nie niosą, to CZYNNOŚĆ,
-- która wariant wytworzyła, i jej NASTAWY — a bez nich `design.photo.history.get`
-- pokazywałby wykaz obrazków bez słowa o tym, co je od siebie różni, a
-- `design.photo.preset.*` nie miałoby skąd wziąć powtarzalnego zestawu.
--
-- ── Dlaczego nastawy jadą jednym zapisem JSON ───────────────────────────────
-- Nastawy różnią się między czynnościami: kadr niesie prostokąt, korekcja barwy
-- dziewięć suwaków, filtr nazwę i siłę. Kolumna na każde pole każdej czynności
-- dałaby tabelę o czterdziestu kolumnach, z których przy każdym wierszu pusta
-- jest trzydzieści siedem. Żadne zapytanie nie pyta o pojedynczą nastawę —
-- czyta się je kompletem, tym samym, którym przyszły w żądaniu.
--
-- ── Droga rachunku jest ZAPISANA, nie odtwarzana ────────────────────────────
-- Cztery czynności (`upscale`, `background.remove`, `inpaint`, `expand`) mają
-- wariant neuronowy i wariant rachunkowy, a ich wyniki różnią się jakością.
-- Kolumna `policzone_przez` trzyma to, co odpowiedź powiedziała polem
-- `computedBy`. Bez niej po tygodniu nie dałoby się powiedzieć, czy dany wariant
-- wyszedł z kanału modelu, czy z rachunku wkompilowanego — a to jest pierwsze
-- pytanie, gdy wynik zawodzi.
--
-- Więz obcy z kasowaniem kaskadowym: czynność bez zasobu, który z niej powstał,
-- nie opisuje niczego.

CREATE TABLE czynnosc_fotografii_design (
    id              INTEGER PRIMARY KEY AUTOINCREMENT,
    zasob_id        INTEGER NOT NULL REFERENCES zasob_design(id) ON DELETE CASCADE,
    zasob_zrodla_id INTEGER          REFERENCES zasob_design(id) ON DELETE SET NULL,
    komenda         TEXT    NOT NULL,
    nastawy_json    TEXT,
    policzone_przez TEXT,
    utworzono       TEXT    NOT NULL DEFAULT (strftime('%Y-%m-%dT%H:%M:%fZ','now'))
);

-- Odczyt idzie po zasobie wynikowym (łańcuch wstecz) i po źródle (co z tego
-- zasobu powstało).
CREATE INDEX idx_czynnosc_fotografii_design_zasob
    ON czynnosc_fotografii_design(zasob_id, id);
CREATE INDEX idx_czynnosc_fotografii_design_zrodlo
    ON czynnosc_fotografii_design(zasob_zrodla_id, id);
