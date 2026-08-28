-- Migracja 080 — Terminal znika z bocznej nawigacji CodeStudio.
--
-- Terminal nie jest samodzielnym modułem: jest oknem pomocniczym wewnątrz
-- modułów Developer, Diagnostics i Apps. Klient trzyma go jako panel pomocniczy
-- (`client/src/okna-pomocnicze/wytwornia-paneli.ts`) i opisuje w trzech
-- modułach, które go goszczą (`client/src/moduly/{developer,diagnostics,apps}/`).
-- Rdzeń nadal podaje Terminal jako pozycję wykazu modułów CodeStudio, bo wiersz
-- `('codestudio','terminal')` w macierzy `srodowisko_modul` ma `widoczny = 1`,
-- a wykaz modułów idzie wyłącznie z rdzenia. Ta migracja przenosi ten fakt do
-- danych, zamiast zostawiać go w warunku wpisanym w kod.
--
-- Bez DELETE: wiersz zostaje w macierzy z `widoczny = 0`. Macierz ma być
-- kompletem par, żeby „moduł niewidoczny" i „moduł, o którym nikt nie
-- zdecydował" dały się w bazie odróżnić; usunięcie wiersza zrobiłoby
-- z odsłonięcia Terminala wstawianie wiersza zamiast zmiany bitu.
--
-- Bez ruchu w innych środowiskach: Terminal ma wiersze także w TalkIn
-- i WorkSpace, oba już `widoczny = 0`. Ten krok ich nie dotyka; zawężenie
-- `srodowisko_id = codestudio` jest tu po to, żeby czytający wiedział, czego
-- migracja nie robi.
--
-- Kolumna `kolejnosc` nie jest pozycją wystawianą kontraktowi:
-- `modulKontraktu` (`core/przeklad_nawigacja.go`) bierze pozycję z porządku
-- zapytania, a `kolejnosc` służy wyłącznie za klucz `ORDER BY sm.kolejnosc,
-- m.kod` (`dane/moduly.go`). Pozostałe siedem par CodeStudio dostaje numery
-- 0–6 bez przerwy, a wiersz Terminala dostaje `kolejnosc = 0` — zero jest tu
-- brakiem treści, nie pierwszym miejscem.
--
-- Co się nie zmienia: moduł `terminal` zostaje w tabeli `modul` z `aktywny`
-- i rodzajem `srodowisko_robocze`. Zostają jego trzy okna operacyjne, tabele
-- `terminal_karta` i `terminal_proces` oraz wszystkie cztery komendy
-- `terminal.*` — żadna z nich nie czyta macierzy widoczności.
--
-- Idempotencja: oba polecenia to `UPDATE` przypisujące wartości bezwzględne
-- (nie przyrosty), zawężone kodami, a nie identyfikatorami wierszy. Powtórne
-- wykonanie daje ten sam stan; migracja nie tworzy wierszy, więc nie potrzebuje
-- `ON CONFLICT`.

-- ── Terminal traci widoczność w CodeStudio ────────────────────────────────────
UPDATE srodowisko_modul
   SET widoczny  = 0,
       kolejnosc = 0
 WHERE srodowisko_id = (SELECT id FROM srodowisko WHERE kod = 'codestudio')
   AND modul_id      = (SELECT id FROM modul      WHERE kod = 'terminal');

-- ── Siedem pozostałych pozycji CodeStudio bez przerwy w numeracji ─────────────
-- Wartości bezwzględne, nie przyrosty — dlatego krok wolno powtórzyć.
WITH nowa_kolejnosc(modul_kod, kolejnosc) AS (
    VALUES
        ('workspace',   0),
        ('roundtable',  1),
        ('design',      2),
        ('developer',   3),
        ('diagnostics', 4),
        ('apps',        5),
        ('agents',      6)
)
UPDATE srodowisko_modul
   SET kolejnosc = (SELECT n.kolejnosc
                      FROM nowa_kolejnosc n
                      JOIN modul m ON m.kod = n.modul_kod
                     WHERE m.id = srodowisko_modul.modul_id)
 WHERE srodowisko_id = (SELECT id FROM srodowisko WHERE kod = 'codestudio')
   AND modul_id IN (SELECT m.id
                      FROM modul m
                      JOIN nowa_kolejnosc n ON n.modul_kod = m.kod);
