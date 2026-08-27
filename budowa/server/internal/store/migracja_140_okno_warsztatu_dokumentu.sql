-- Migracja 140 dopisuje do katalogu okien wiersz Warsztatu dokumentu z rolą pomocniczą i kategorią narzędzia, kontynuując numerację pozycji migracji 126.

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
