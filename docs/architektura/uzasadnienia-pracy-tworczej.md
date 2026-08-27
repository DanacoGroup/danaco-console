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

## adapter_modul_studio_wydanie_formatu.go

Format uboższy niż dokument jest sytuacją normalną, przemilczenie straty nie
jest — każde wydanie oddaje wykaz cech pominiętych pozycję po pozycji, z nazwą
cechy i jej rozmiarem policzonym z postaci dokumentu, a nie zdaniem ogólnym.
Rachunek stoi bez programu zewnętrznego: txt, md, html i rtf idą rachunkiem
własnym tego pliku, docx i odt przez złożenie OOXML i ODF, a pdf biblioteką
`pdfcpu` przez profil wydania. Bajty wyniku idą do magazynu zasobów rdzenia;
ścieżka na maszynie serwera, gdy operator ją wskazał, jest drogą dodatkową,
nie zamiast — instalka cienka sięga po wynik zasobem, nie po ścieżce serwera.

Wyróżnienie tła (znakowanie sesji pracy modelu i operatora) schodzi z postaci
przed złożeniem pliku, jednakowo dla każdego formatu: pismo wysłane na
zewnątrz nie ma wyjść w plamach roboczych, tak jak nie ma wynieść komentarzy
z marginesu. Zdjęcie idzie na kopii postaci wydania, nie w bazie, więc dokument
zastany zostaje ze znakowaniem nietkniętym, a zdjęcie jest nazwane w bilansie.
Styl nazwany z wyróżnieniem zostaje nietknięty, bo to postanowienie o wyglądzie
dokumentu, a nie znakowanie sesji — zdejmowanie go byłoby zmianą, o którą nikt
nie prosił. Ten sam podział — kopia na czas wydania, baza nietknięta — obowiązuje
przy narzuceniu postaci szablonu: szablon podmienia arkusz stylów, nastawy
strony, nagłówek i stopkę na czas wydania, a nie treść ani zapis w bazie.

Przypis dolny wychodzi do pliku tekstowego czy markdown wykazem na końcu, bo
w tych formatach nie ma gdzie postawić odsyłacza w miejscu — ale przypis, który
stał pod stroną przy swoim zdaniu, i przypis na końcu pliku to nie ta sama
cecha: znika powiązanie z miejscem i numeracja odświeżalna, więc strata jest
nazwana z liczbą po rodzaju elementu aparatu (przypisy, spisy, odsyłacze,
podpisy), a nie jednym workiem.

RTF wchodzi do wykazu formatów wykonalnych rachunkiem własnym: format tekstowy
o prostej gramaturze, w którym rachunek niesie postać znaku, akapit, wyrównanie
i tabelę; obrazy osadzone i aparat odświeżalny wychodzą w wykazie pominiętych.

Profil wydania PDF niesie nastawy paginacji i stopki; profil niewskazany nie
jest odmową — dokument wychodzi z paginacją i stopką domyślną, a bilans mówi
wprost, że profilu nie było.

## adapter_modul_studio_wejscie_czynnosci.go

Każda czynność tego pliku wnosi treść wprost do edytora, gotową do pracy,
z zachowaną postacią — nie jako załącznik, nie jako pozycja kolejki do
przyjęcia. Dlatego każda czynność kończy się zapisaną postacią dokumentu,
zapisem przez ten sam obszar postaci (`wejscieUtrwalPostac`), a dziennik
czynności dostaje wpis rodzaju `importChange`, którym operator cofa
„wniesienie pliku", odróżnione od zwykłej zmiany treści. Jedyny wyjątek jest
PDF ze samych skanów: nie ma czego wnieść do edytora, więc idzie na
rozpoznanie pisma, a odpowiedź to mówi wprost, zamiast oddać pustą kartkę.
Każde wniesienie oddaje bilans tego, co odzyskano, a czego nie — plik wnoszony
niesie rzeczy, których rdzeń nie odczytuje, i przemilczenie tego zamieniłoby
dokument okaleczony w dokument wczytany bez uwag.

Format pliku wnoszonego rozpoznaje się po zawartości, nie po rozszerzeniu,
a zapis znaków ustala się przed rozbiorem treści: pismo w kodowaniu innym niż
utf-8, wczytane jakby nim było, daje ciąg znaków bez sensu, gorszy niż odmowa,
bo model przeczyta go jako słowa i zacznie na nich pracować.

Wniesienie pliku w miejsce kursora do dokumentu istniejącego wnosi treść wraz
z jej postacią znaku fragmentów, nie całą postać wnoszonego pliku — wstawienie
całej postaci w środek dokumentu zastanego podmieniłoby jego arkusz stylów
i nastawy strony, o co nikt nie prosił; strata jest nazwana w bilansie.

Kopia dokumentu i dokument założony z szablonu dostają nowe identyfikatory dla
wszystkich bytów postaci (`wejsciePrzepiszIdentyfikatory`): identyfikator
wspólny z oryginałem byłby jednym wierszem bazy widzianym z dwóch dokumentów,
a zapis pod takim kluczem przeniósłby wiersz oryginału do kopii zamiast założyć
wiersz nowy — kopia zabrałaby wtedy oryginałowi jego własną postać. Blokady
fragmentów w tym przepisaniu nie uczestniczą: przenosi je droga osobna, razem
z nowym kodem blokady, a przy kopii dokumentu przechodzą domyślnie (fragment
wzorcowy pisma ma zostać wzorcowy też w kopii), podczas gdy znakowanie
i komentarze — domyślnie nie, bo są rozmową o dokumencie, nie jego treścią.

Wniesienie z modułu Design przyjmuje kod węzła do zapisu pochodzenia, a bajty
bierze z zasobu magazynu wskazanego obok, bo węzeł wektorowy nie jest obrazem
rastrowym, dopóki `design.vector.export` go nie wyrysuje do zasobu.

Wniesienie ze strony sieci daje pierwszeństwo fragmentowi wskazanemu w oknie
przeglądarki nad pobraniem całej strony — to fragment, który operator widział
dokładnie. Bez niego rdzeń pobiera stronę przez `net/http`, bez silnika
przeglądarki, więc strona zbudowana wyłącznie skryptem oddaje mało treści;
to jest cena znana i nazwana w bilansie, nie przeoczenie rachunku.

## adapter_modul_studio_wydanie.go

Archiwum eksportu składa biblioteka standardowa Go (`archive/zip`,
`archive/tar`), wkompilowana w binarium rdzenia. Wywołanie zewnętrznego
programu pakującego byłoby zależnością spoza instalki: u Operatora
„Eksportuj historię" kończyłoby się odmową, choć rdzeń meldowałby komendę
jako obsłużoną.

Manifest jest częścią archiwum, nie dodatkiem. Archiwum bez manifestu to
katalog plików, o których nie wiadomo, w jakiej kolejności powstały, kto je
zapisał ani do czego należą, dlatego każde wydanie niesie `manifest.json`,
nawet gdy wersja jest jedna.

Paczka redakcyjna wyłącza niezależnie historię, raport zmian i adnotacje, bo
nie każde przekazanie jest przekazaniem redakcyjnym; dokument finalny
wyłączyć się nie da, ponieważ bez niego paczka nie jest paczką.

Wersja wskazana do wydania repozytorium, ale nie należąca do danego
dokumentu, kończy się odmową, a nie cichym pominięciem: Operator dostałby
archiwum krótsze, niż prosił, i nie dowiedziałby się dlaczego.

Format raportu różnicy bierze się ze słownika zamiany formatu dokumentu, więc
ten sam raport da się oddać do wglądu (PDF), do dalszej redakcji (DOCX,
Markdown) albo do odczytu maszynowego (TXT). Format nieznany kończy się
odmową, nie zejściem na zapis tekstowy: Operator, który poprosił o DOCX,
a dostałby plik tekstowy z rozszerzeniem `.docx`, dowiedziałby się o tym
dopiero przy otwieraniu dokumentu u odbiorcy.

DOCX jest archiwum ZIP z trzema częściami obowiązkowymi, więc składa się go
`archive/zip` bez biblioteki biurowej i bez klucza komercyjnego, których
instalka nie niesie. Dokument nie ma stylów ani tabel i mieć nie udaje: niesie
akapity tekstu, otwiera się w każdym edytorze i daje się dalej redagować,
a to jest dokładnie to, po co paczka redakcyjna powstaje.

## adapter_modul_studio_symbole.go

Tablica znaków, jej punkty kodowe i grupy stanowią wiedzę rdzenia, tak samo jak
arkusz stylów fabryczny i wykaz nośników druku. Do bazy schodzi wyłącznie to, co
Operator zmienił — jego zasady autozamiany — albo czym się posłużył — jego znaki
ostatnio użyte. Wpisanie tablicy do migracji dałoby dwa wykazy, które rozjadą się
przy pierwszym uzupełnieniu.

Wstawienie znaku idzie drogą zmiany treści, ponieważ znak specjalny jest znakiem,
nie obiektem: ma się liczyć do długości akapitu, znaleźć w wyszukiwaniu i
przenieść przy zmianie formatu tak samo jak litera. Wstawienie idzie przez
funkcję zamiany treści, tę samą drogę, którą idzie pisanie, dzięki czemu zmiana
śledzona autora i wpis dziennika odkładają się bez osobnego mechanizmu.

Autozamiana jest nastawą, a nie czynnością na dokumencie: zamiana zachodzi w
chwili pisania, w oknie, na naciśnięcie klawisza. Rdzeń trzyma wykaz zasad,
wystawia go oknu i pozwala go zmienić. Gdyby rdzeń przepuszczał treść przez
zasady przy zapisie, znak wpisany świadomie zostałby zamieniony wbrew woli
piszącego bez możliwości cofnięcia tego pojedynczego przypadku.

Wykaz zasad, także fabrycznych, stoi w bazie, nie w tym pliku. Migracja 368
założyła tabelę zasad autozamiany wraz z kolumną oznaczającą zasadę fabryczną i
wpisała zasady fabryczne wierszami, aby dało się je wyłączyć; migracja 371
dołożyła do tego wykazu znaki prawnicze i ułamki. Powtórzenie wykazu fabrycznego
w tym pliku byłoby drugą prawdą o tym, co wchodzi w miejsce skrótu, a pierwsza
poprawka rozjechałaby obie wersje. Jest to odwrotne podejście niż przy tablicy
znaków: tablicy znaków Operator nie zmienia, a zasadę autozamiany zmienia i
wyłącza.
