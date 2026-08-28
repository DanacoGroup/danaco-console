-- Migracja 072 — dopełnienie macierzy widoczności modułów do pełnych 45 wierszy.
--
-- Macierz `srodowisko_modul` niosła wyłącznie pary widoczne, a para niewidoczna
-- istniała jako brak wiersza — „moduł niewidoczny w tym środowisku" i „moduł,
-- o którym nie zdecydowano" były w bazie nieodróżnialne. Komplet wierszy czyni
-- z macierzy konfigurację: żeby pokazać moduł w środowisku, wystarczy zmienić
-- jeden bit.
--
-- Wiersz `widoczny = 0` nie jest zakazem i niczego nie blokuje. Boczna
-- nawigacja bierze pozycje zapytaniem `WHERE sm.widoczny = 1`
-- (`dane/moduly.go`), więc niewidoczna para daje krótszą listę, a nie pozycję
-- wyszarzoną ani odmowę.
--
-- Migracja wyłącznie dopisuje brakujące pary — nie ma tu ani jednego `UPDATE`,
-- a `ON CONFLICT DO NOTHING` czyni ją idempotentną także na bazie, w której
-- część par już założono.
--
-- MultitaskingAI zostaje poza macierzą. `core/przeklad_nawigacja.go` rozstrzyga
-- rodzaj nawigacji środowiska wyłącznie pustością jego wykazu modułów: brak
-- modułów daje `NavigationKindOrchestration`, czyli panel orkiestracji zamiast
-- listy. Wiersz w macierzy — nawet niewidoczny — odebrałby tę własność
-- i wymusiłby rozstrzyganie po kodzie środowiska wpisanym w warunek. Pełną
-- macierz tworzy zatem 15 modułów × 3 środowiska kolumnowe = 45 par.
--
-- Kolumna `kolejnosc` niesie pozycję w bocznej nawigacji, a pozycji
-- niewidocznej w nawigacji nie ma. Zero jest tu brakiem treści, nie pierwszym
-- miejscem; pozycję nadaje się, gdy para stanie się widoczna.

-- ── Siedemnaście par niewidocznych ──────────────────────────────────────────
WITH macierz(srodowisko_kod, modul_kod) AS (
    VALUES
        -- TalkIn — środowisko pracy z treścią; poza wykazem zostają narzędzia
        -- wytwórcze i programistyczne.
        ('talkin',     'automations'),
        ('talkin',     'design'),
        ('talkin',     'terminal'),
        ('talkin',     'developer'),
        ('talkin',     'diagnostics'),
        ('talkin',     'apps'),
        -- WorkSpace — środowisko prowadzenia projektu; poza wykazem zostaje
        -- warsztat kodu i dwa moduły operacyjne TalkIn.
        ('workspace',  'translate'),
        ('workspace',  'assistant'),
        ('workspace',  'terminal'),
        ('workspace',  'developer'),
        ('workspace',  'diagnostics'),
        -- CodeStudio — środowisko programowania; poza wykazem zostaje praca
        -- z dokumentem i źródłami.
        ('codestudio', 'studio'),
        ('codestudio', 'browser'),
        ('codestudio', 'research'),
        ('codestudio', 'library'),
        ('codestudio', 'translate'),
        ('codestudio', 'assistant')
)
INSERT INTO srodowisko_modul (srodowisko_id, modul_id, kolejnosc, widoczny)
SELECT s.id, m.id, 0, 0
  FROM macierz
  JOIN srodowisko s ON s.kod = macierz.srodowisko_kod
  JOIN modul m      ON m.kod = macierz.modul_kod
 WHERE true
ON CONFLICT(srodowisko_id, modul_id) DO NOTHING;
