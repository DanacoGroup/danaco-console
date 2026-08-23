-- Zdjęcie okna rozmowy (`chat-window`) z modułu Library.
--
-- Library nie prowadzi rozmowy z modelem: przyjmuje pliki, porządkuje je
-- w katalogi i dzieli na projekty. Okno rozmowy przypina wszystkim modułom
-- reguła zakładająca macierz `okno_operacyjne_modul`; ten krok zdejmuje z niej
-- jedno wystąpienie. Definicja okna w `okno_operacyjne`, przypięcia pozostałych
-- modułów, pozostałe okna Library oraz rodzina komend `library.*` zostają
-- nietknięte.
--
-- Oba wskazania idą podzapytaniem po kodzie, a nie po kluczu wpisanym z ręki:
-- klucze sztuczne nadaje AUTOINCREMENT i nie są one tym samym w każdej bazie.
-- Gdy któregokolwiek z wierszy nie ma, podzapytanie oddaje NULL, porównanie nie
-- jest prawdziwe dla żadnego wiersza i krok kasuje zero wierszy zamiast pomylić
-- się w drugą stronę.
DELETE FROM okno_operacyjne_modul
 WHERE okno_operacyjne_id = (SELECT id FROM okno_operacyjne WHERE kod = 'chat-window')
   AND modul_id           = (SELECT id FROM modul WHERE kod = 'library');
