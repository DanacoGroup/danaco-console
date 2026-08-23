-- Migracja 140 — Warsztat dokumentu wchodzi do katalogu okien.
--
-- Moduł Studio dostał okno rodzin `studio.pdf.*` i `studio.security.*`. Bez
-- wiersza w katalogu klient postawiłby okno, o którym rdzeń nie wie: `module.list`
-- zaniżałby zakres modułu, a pas uczciwości okna meldowałby kod „poza katalogiem".
-- Ten sam rozjazd domykała migracja 126 dla siedmiu innych modułów.
--
-- Rola `pomocnicze`, bo okno wspiera pracę Studio Editora, a nie prowadzi jej
-- samo: pracuje na materiale wniesionym do okna, nie na treści redagowanej.
-- Kategoria `narzedzia` wzorem Tools Panelu, z którym dzieli rolę w module.
-- Kolejność w kategorii jest pierwszą wolną po pozycjach migracji 126.

INSERT INTO okno_operacyjne (kod, nazwa, rola, kategoria, kolejnosc)
VALUES
    ('studio.document-workshop', 'Warsztat dokumentu', 'pomocnicze', 'narzedzia', 16)
ON CONFLICT(kod) DO NOTHING;

INSERT INTO okno_operacyjne_modul (okno_operacyjne_id, modul_id, kolejnosc)
SELECT o.id, m.id, 8
  FROM okno_operacyjne o
  JOIN modul m ON m.kod = 'studio'
 WHERE o.kod = 'studio.document-workshop'
ON CONFLICT(okno_operacyjne_id, modul_id) DO NOTHING;
