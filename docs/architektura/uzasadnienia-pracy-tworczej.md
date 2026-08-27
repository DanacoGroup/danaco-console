# Uzasadnienia pracy twórczej

Dokument gromadzi uzasadnienia decyzji projektowych, które w kodzie źródłowym
nie mieszczą się w nagłówku komentarza. Każdy rozdział nosi nazwę pliku,
którego uzasadnienia dotyczą.

## adapter_modul_studio_ooxml.go

Rachunek nie korzysta z żadnej biblioteki obcej ponad standardową bibliotekę
Go. `.docx` jest archiwum ZIP z dokumentem XML w środku, a biblioteka wzorcowa
Go niesie oba potrzebne składniki — `archive/zip` i `encoding/xml`. Biblioteka
`unioffice` odpada na licencji handlowej, a LibreOffice i `pandoc` byłyby
procesem potomnym tam, gdzie proces potomny jest regresem: jedno binarium
serwera zamienia się w dwa, dochodzi koszt uruchomienia i rozjazd wersji
między maszynami.

Rozbiór dokumentu idzie przez drzewo węzłów, nie przez strumień tokenów
z ręcznym stanem. Postać akapitu OOXML leży w węźle `w:pPr`, który stoi PRZED
fragmentami tekstu, ale odnosi się do całego akapitu, a postać znaku leży
w `w:rPr` wewnątrz każdego fragmentu. Automat stanów musiałby pamiętać oba
konteksty naraz na każdej głębokości; drzewo pozwala je odczytać tam, gdzie
są. Ten sam rozbiór obsługuje ODF, bo i tam postać leży w węzłach obok treści.

OOXML mierzy odległości w twipach (1/1440 cala), stopień pisma w półpunktach,
a odstępy akapitowe w twipach, podczas gdy kontrakt platformy mówi
w milimetrach i punktach. Przeliczenie stoi w jednym miejscu rachunku, bo
przeliczenie rozsypane po wielu miejscach rozjeżdża się przy pierwszej
poprawce.

Bilans wniesienia dokumentu jest obowiązkowy, nie ozdobny: plik wejściowy
niesie rzeczy, których rachunek nie odczytuje (pola obliczane, wykresy,
kształty rysowane, osadzone obiekty innych programów), a przemilczenie ich
zamieniłoby dokument okaleczony w dokument wczytany bez uwag.

Zasięgi nagłówka i stopki — domyślny, pierwszej strony i stron parzystych —
idą osobno, bo nagłówek urzędowy bywa dla tych trzech przypadków różny.
Zasięg pierwszej strony dokłada do sekcji przełącznik `w:titlePg`, bez
którego dokument nagłówka pierwszej strony nie pokaże.

## adapter_modul_studio_odf.go

Rachunek OpenDocument stoi na tym samym drzewie węzłów co rachunek OOXML,
bo `.odt` jest archiwum ZIP z `content.xml` i `styles.xml` w środku — ten
sam układ, w którym postać leży w węzłach obok treści. Drugiego rozbioru
XML rachunek nie zakłada.

Trzy różnice rozstrzygają kształt tego pliku wobec OOXML. Po pierwsze,
postać bezpośrednia nie istnieje: w ODF każde odstępstwo od stylu nazwanego
musi być stylem automatycznym o własnej nazwie, wymienionym w
`office:automatic-styles` — odczyt trzyma mapę stylów automatycznych,
a zapis je wytwarza; nazwy tych stylów są nazwami technicznymi formatu
pliku, takie same nazywa LibreOffice. Po drugie, miary idą jednostką
zapisaną w napisie (`fo:page-width="21cm"`, `fo:font-size="12pt"`), stąd
przeliczenie stoi w jednym miejscu rachunku. Po trzecie, nagłówek i stopka
wiszą na stronie wzorcowej (`style:master-page`), nie na sekcji — sekcje ODF
nie mają własnych nastaw strony w sensie, w którym mają je sekcje OOXML,
więc dokument wniesiony z `.odt` dostaje jedną sekcję, a nie sekcje udawane.

Nagłówek strony pierwszej i stron lewych wychodzą osobnymi węzłami, bo
w pismach urzędowych te trzy przypadki bywają różne.

Rodzaj archiwum (`mimetype`) musi być pierwszym składnikiem i bez kompresji,
bo tak stanowi norma OpenDocument i tak to sprawdzają czytniki.

OpenDocument zapisuje scalenie poziome komórek tabeli dwa razy: raz jako
`table:number-columns-spanned` komórki scalającej, raz jako
`table:covered-table-cell` na każdej z przykrytych kolumn — jest ich
dokładnie o jedną mniej niż rozpiętość. Liczenie obu naraz dawało tabelę
trzykolumnową jako czterokolumnową, z czwartą kolumną bez szerokości,
niewidoczną w dokumencie. Rozbiór dlatego prowadzi licznik przykrytych
komórek do pominięcia, a nie porównanie z granicą kolumny: po przejściu
komórki scalającej numer kolumny stoi już za scaleniem, więc porównanie
z granicą nie odróżnia przykrycia poziomego od pionowego.

## adapter_modul_studio_uchwyty.go

Metody portu Studio leżą w wielu plikach adaptera (nastawy, postać, aparat,
katalogi, gałęzie, wydanie i inne) — ten plik jest jedynym miejscem, które
wpina je razem do rejestru komend.

Zdarzenie `studio.document.changed` rozgłasza się wyłącznie tam, gdzie treść
widoczna w oknie pracy z dokumentem naprawdę się zmienia, a odpowiedź komendy
tej zmienionej treści nie niesie już sama. Stąd trzy grupy zachowań w pliku:

- Operacja kontekstowa, zapis, przywrócenie wersji, założenie gałęzi (gdy
  zostaje otwarta jako treść bieżąca), scalenie zakończone powodzeniem oraz
  osadzenie zasobu zmieniają treść i jej odpowiedź nie niesie — rozgłaszają
  zmianę, doczytując dokument przez `document.open` tam, gdzie trzeba.
- Czynności postaci dokumentu (formatowanie, strona, listy, tabele, obiekty,
  aparat, pola) oddają postać po zmianie wprost w odpowiedzi, więc zdarzenia
  nie rozgłaszają — okno ma już to, czym miałoby się odświeżyć. Wyjątek jest
  jeden: zapis postaci oddaje także dokument wraz z wersją, więc rozgłasza
  zmianę tak jak zwykły zapis, bo po nim odświeżają się także wykazy i historia.
- Odczyty, wykazy i nastawy nie ruszają treści i zdarzenia nie mają. Kolejka
  wczytywania i jej odczyt, rozpoznanie oraz korekta rozgłaszają odrębne
  zdarzenie stanu pozycji, nie zdarzenie dokumentu — dopiero przyjęcie pozycji
  zakłada dokument i rozgłasza zmianę dokumentu. Schowek zdarzenia nie
  rozgłasza, mimo że wycięcie i wklejenie zmieniają treść: kontrakt tych
  odpowiedzi dokumentu nie niesie, a doczytanie go drugą drogą tylko po to,
  żeby rozgłosić zmianę, kosztowałoby dwa odczyty na każde wklejenie — okno
  odświeża się postacią, którą odpowiedź już niesie.
- Kontrola pracy (dziennik, zmiany modelu, znakowanie, autozapis, kopie
  zapasowe, zajęcia wykonawców) rozgłasza zmianę dokumentu tylko tam, gdzie
  czynność podmienia treść widoczną w oknie i oddaje dokument w odpowiedzi:
  cofnięcie i ponowienie czynności, cofnięcie zmian modelu, rozstrzygnięcie
  znakowania i propozycji, przeniesienie fragmentu różnicy, przywrócenie
  kopii i powrót do wersji założycielskiej.

Podpis wykonawcy wpina się w rejestrze, a nie w każdym obsługiwaczu z osobna,
z tego samego powodu, dla którego zapora blokad stoi w drzwiach, a nie przy
każdym stoliku: tożsamość wykonawcy niesie żądanie, ale droga wyjścia
czynności postaci tego żądania nie widzi. Rozłożenie podpisu po trzydziestu
ośmiu sygnaturach czynności postaci znaczyłoby trzydzieści osiem miejsc do
pominięcia przez pomyłkę — a pominięcie nie byłoby widoczne: zmiana
zapisałaby się dalej, tylko podpisana błędnym autorem. Owinięcie rejestru
obejmuje jednym warunkiem wszystkie te komendy, także te jeszcze nienapisane.

Czynności postaci dokumentu rozstrzygają autora z pola żądania, bo tak
stanowi kontrakt platformy. Wykonawca, który tego pola nie poda, zostałby
zapisany jako autor domyślny, a jego zmiana nie odłożyłaby się jako zmiana
śledzona i nie dałoby się jej podświetlić przełącznikiem pracy modelu —
to jest usterka do naprawy, nie ograniczenie do zgłoszenia. Naprawa stoi
w rejestrze, nie w trzydziestu ośmiu czynnościach postaci: gdy fakt gniazda
mówi, że woła wykonawca, wpięcie dopisuje autora modelu do ładunku, zanim
ładunek zobaczy obsługiwacz. Stempel idzie wyłącznie w jedną stronę — podnosi
autora domyślnego do wykonawcy, nigdy odwrotnie — bo żądanie, które samo
podaje się za wykonawcę, jest twierdzeniem modelu o sobie, nie faktem
gniazda serwera narzędzi.
