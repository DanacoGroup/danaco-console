-- Migracja 340 wprowadza do katalogu okien sześć warsztatów modułu Design
-- wraz z ich rolami i kolejnością.

INSERT INTO okno_operacyjne (kod, nazwa, rola, kategoria, kolejnosc)
VALUES
    ('design.photo-workshop',       'Warsztat fotografii',        'pomocnicze', 'narzedzia', 17),
    ('design.vector-workshop',      'Warsztat wektora',           'pomocnicze', 'narzedzia', 18),
    ('design.print-workshop',       'Warsztat druku',             'pomocnicze', 'narzedzia', 19),
    ('design.stock-browser',        'Przeglądarka baz zdjęciowych', 'zarzadca', 'narzedzia', 20),
    ('design.publication-workshop', 'Warsztat publikacji',        'pomocnicze', 'narzedzia', 21),
    ('design.mockup-workshop',      'Warsztat makiety',           'pomocnicze', 'narzedzia', 22)
ON CONFLICT(kod) DO NOTHING;

INSERT INTO okno_operacyjne_modul (okno_operacyjne_id, modul_id, kolejnosc)
SELECT o.id, m.id, o.kolejnosc - 11
  FROM okno_operacyjne o
  JOIN modul m ON m.kod = 'design'
 WHERE o.kod IN ('design.photo-workshop', 'design.vector-workshop',
                 'design.print-workshop', 'design.stock-browser',
                 'design.publication-workshop', 'design.mockup-workshop')
ON CONFLICT(okno_operacyjne_id, modul_id) DO NOTHING;
