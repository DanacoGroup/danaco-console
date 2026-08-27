-- Migracja 127 dopisuje do katalogu okien execution-loop-window wraz z brakującym przypięciem chat-window do modułu Library, domykając wykaz okien wspólnych platformy.

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
