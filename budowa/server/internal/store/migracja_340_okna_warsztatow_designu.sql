-- Migracja 340 — pięć warsztatów modułu Design wchodzi do katalogu okien.
--
-- Komenda bez okna jest funkcją, której Operator nie ma. Moduł Design dostał
-- siedemdziesiąt pięć nowych komend (wektor, makieta, barwa, ikony, marketing,
-- druk, fotografia, publikacje) i bez wierszy w katalogu klient postawiłby okna,
-- o których rdzeń nie wie: `module.list` zaniżałby zakres modułu, a pas
-- uczciwości meldowałby kod „poza katalogiem". Ten sam rozjazd domykała
-- migracja 140 dla warsztatu dokumentu w Studio.
--
-- Warsztatów jest sześć, nie pięć: makieta dostała własne okno, bo ramki, więzy,
-- komponenty i prototyp są innym rodzajem pracy niż rysowanie i inny materiał
-- biorą. Wciśnięcie ich w warsztat wektora dałoby okno o dwudziestu pięciu
-- czynnościach, którego Operator nie przejrzy.
--
-- Role: warsztat fotografii i warsztat wektora PROWADZĄ pracę na materiale, więc
-- `wiodace` im nie przysługuje (wiodącym modułu jest Design Board — tam stoi
-- praca koncepcyjna), ale nie są też tylko pomocnicze — pracują na treści.
-- Rola `narzedziowe` byłaby zmyśleniem, bo katalog jej nie zna, więc idą jako
-- `pomocnicze`: wspierają Design Board materiałem, który mu podają. Przeglądarka
-- baz zdjęciowych jest `zarzadca` — jej praca to wybór z wykazu, jak w Assets
-- Panelu. Warsztat druku i warsztat publikacji są `pomocnicze`: wydają to, co już
-- powstało.
--
-- Kolejność w module idzie po pozycjach zajętych (Design Board 1, Preview
-- Window 2, Assets Panel 3, Prompt Builder 4) i po Tokens & System Panelu,
-- który wszedł osobno.

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
