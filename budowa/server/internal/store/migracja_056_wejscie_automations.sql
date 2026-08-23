-- Migracja 056 — wejście do modułu Automations.
--
-- Boczna nawigacja bierze wykaz pozycji wyłącznie z macierzy
-- `srodowisko_modul` (`environment.enter` → `module.list`), więc drogę wejścia
-- do modułu dokłada się wierszem w macierzy, nie literałem w kliencie.
--
-- WorkSpace — środowisko procesów i automatyzacji; CodeStudio — harmonogramy
-- i kolejki uruchomień są narzędziem pracy nad kodem tak samo jak terminal.
-- TalkIn zostaje bez Automations: to środowisko pracy z treścią. Kolejność
-- dokłada moduł na końcu wykazu obu środowisk.
--
-- ON CONFLICT DO NOTHING — migracja przechodzi także na bazie, w której wiersz
-- już założono.

WITH macierz(srodowisko_kod, modul_kod, kolejnosc) AS (
    VALUES
        ('workspace',  'automations', 9),
        ('codestudio', 'automations', 8)
)
INSERT INTO srodowisko_modul (srodowisko_id, modul_id, kolejnosc, widoczny)
SELECT s.id, m.id, macierz.kolejnosc, 1
  FROM macierz
  JOIN srodowisko s ON s.kod = macierz.srodowisko_kod
  JOIN modul m      ON m.kod = macierz.modul_kod
 WHERE true
ON CONFLICT(srodowisko_id, modul_id) DO NOTHING;
