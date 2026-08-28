-- Migracja zdejmuje okno rozmowy z modułu biblioteki, bo moduł biblioteki przyjmuje pliki
-- i porządkuje je, nie prowadzi rozmowy z modelem.
DELETE FROM okno_operacyjne_modul
 WHERE okno_operacyjne_id = (SELECT id FROM okno_operacyjne WHERE kod = 'chat-window')
   AND modul_id           = (SELECT id FROM modul WHERE kod = 'library');
