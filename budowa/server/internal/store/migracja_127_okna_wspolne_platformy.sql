-- Migracja 127 — okna wspólne platformy w katalogu rdzenia.
--
-- Dwa okna towarzyszą każdemu modułowi: Chat Window (kanał Użytkownik ↔
-- Wykonawca) i Execution Loop Window (kanał Koordynator ↔ Wykonawca). Wszystkie
-- opracowania modułów wymieniają je razem, na czele wykazu okien, jako kolumny
-- stałe układu.
--
-- Katalog rdzenia znał do tej pory jedno z nich, i to nie dla wszystkich:
--
--   * `execution-loop-window` nie miał wiersza definicji w ogóle, więc żaden
--     moduł nie mógł go wskazać, a pas uczciwości modułów meldował ten kod jako
--     stojący poza katalogiem;
--   * `chat-window` przypina wszystkim modułom migracja 031, ale moduł Library
--     wszedł do rejestru po niej, więc przypięcia nie dostał. Jako jedyny
--     z piętnastu.
--
-- Kolejność zero należy do obu okien wspólnych, nie do jednego. Odczyt sortuje
-- `ORDER BY om.kolejnosc, o.kod` (`dane/okna_operacyjne.go`), więc remis przy
-- zerze rozstrzyga kod: `chat-window` stoi przed `execution-loop-window`, czyli
-- dokładnie tak, jak układa je każde opracowanie. Okna własne modułu zaczynają
-- się od jedynki i nie są przestawiane.

INSERT INTO okno_operacyjne (kod, nazwa, rola, kategoria, kolejnosc)
VALUES
    ('execution-loop-window', 'Execution Loop Window', 'zarzadca', 'komunikacja', 6)
ON CONFLICT(kod) DO NOTHING;

-- Oba okna wspólne dostają przypięcie do KAŻDEGO modułu — także tego, który
-- wejdzie do rejestru później. Zapytanie liczy się z rejestru modułów, a nie
-- z wykazu wypisanego wprost, więc nie da się go rozjechać przez przeoczenie.
INSERT INTO okno_operacyjne_modul (okno_operacyjne_id, modul_id, kolejnosc)
SELECT o.id, m.id, 0
  FROM okno_operacyjne o
 CROSS JOIN modul m
 WHERE o.kod IN ('chat-window', 'execution-loop-window')
ON CONFLICT(okno_operacyjne_id, modul_id) DO NOTHING;
