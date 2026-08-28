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

## adapter_modul_studio_postac.go

Pozostałe obszary postaci dokumentu stoją w plikach sąsiednich (formatowanie,
styl, strona, sekcje, listy, tabele, wstawienia, aparat, widok) i wołają
rachunek z tego pliku, zamiast liczyć zakresy drugi raz. Kontrakt mówi
o zakresie „w znakach", a pismo polskie ma znaki dwubajtowe — cały plik liczy
w runach, nie w bajtach, bo zakres liczony bajtami rozjechałby się na
pierwszym „ż" i dałby wytłuszczoną połowę litery.

Treść dokumentu jest prawdą złożoną w `dokument_studio.tresc`: tam ją zapisuje
`document.save` i tam jej szuka podgląd. Drzewo postaci niesie tę samą treść
rozłożoną na bloki i fragmenty, więc wczytanie postaci uzgadnia drzewo
z treścią — gdy treść zmieniła się drogą, której postać nie zna, bloki idą za
treścią, a postać akapitów i fragmentów zostaje przypisana po kolejności.
Zapis postaci przepisuje treść z drzewa, więc od tej chwili obie strony mówią
to samo. Z tego samego powodu zamiana treści idzie po elementach, nie po
napisie: zamiana napisu w napisie gubi kroje, wcięcia i granice akapitów,
więc treść rozkłada się na elementy — znak wraz ze swoją postacią albo granica
akapitu wraz z postacią akapitu, który za nią się zaczyna — i zamiana jest
wycięciem odcinka elementów oraz wstawieniem odcinka nowego.

Bloki nietekstowe (tabela, obiekt, podział) nie zajmują ani jednego znaku
treści — ich miejsce niesie kolejność bloków, dlatego wstawienie tabeli nie
przesuwa ani jednego zakresu zaznaczenia. Rozcięcie fragmentów na granicach
zakresu jest warunkiem wymagania „zastosowanie postaci do fragmentu zmienia
wyłącznie ten fragment": bez rozcięcia postać nakładałaby się na cały
fragment, który zaznaczenie przecina, nie na samo zaznaczenie.

Blokada odcinka jest skierowana przeciw modelowi, nie przeciw właścicielowi
dokumentu. Zmiana obejmująca blokadę częściowo wykonuje się poza blokadą
i oddaje bilans pominięcia — odmowa całości byłaby nieproporcjonalna, a
przemilczenie pominięcia jest zakazane.

Czynność modelu odkłada się jako zmiana śledzona zawsze, bo na tym stoi
przełącznik „pokaż wszystko, co zrobił model" i wymaganie, żeby zmiana
postaci bez zmiany liter nie była niewidzialna. Czynność operatora odkłada
się wyłącznie wtedy, gdy śledzenie zmian w dokumencie jest włączone — inaczej
operator dostawałby wykaz zmian do przyjęcia po każdym własnym kliknięciu
pogrubienia.

Rodzaj zmiany śledzonej i rodzaj czynności dziennika są dwoma osobnymi
słownikami z osobnymi polami. Zmiana śledzona zna trzy wartości (wstawienie,
usunięcie, formatowanie), bo tyle rozróżnia adiustacja; dziennik zna
jedenaście wartości `StudioActionKind`, bo operator cofa „zmianę stylu", nie
„formatowanie". Podanie jednego słownika w miejsce drugiego przechodzi
kompilację i wywraca się dopiero na ograniczeniu tabeli, dlatego pola są
rozdzielone.

Podpis wykonawcy czynności wchodzi z kontekstu jednym czytaniem: wkłada go
tam wpięcie rejestru `podpisWykonawcyStudia`, które czyta z ładunku pola
`agentId`, `agentName` i `subagentId` — te same, które kontrakt niesie przy
każdej komendzie Studia z polem `author`. Bez tego wpięcia rozbicie zmian po
wykonawcy pokazywałoby wszystkie czynności postaci jako czynności
nienazwanego: rodzaj autora zapisywałby się, a kod agenta nie, choć żądanie
go niosło. Rodzaj autora rozstrzyga się zasadą „szerszy wygrywa" — tą samą,
którą stosuje `kontrolaRozpoznajWykonawce`.

Dziennik czynności jest wspólny dla całego modułu: wykaz i cofanie wystawia
inny odcinek, a każda czynność zmieniająca dokument tam odkłada swój wpis —
inaczej cofnięcie pojedynczej zmiany postaci nie miałoby czego cofnąć. Stan
sprzed czynności idzie całym drzewem postaci, nie samą treścią, bo pomyłkowa
zmiana kroju w całym dokumencie nie rusza ani jednej litery, więc treść
sprzed jej nie odtworzy. Kod agenta w tym wpisie rozróżnia dwóch wykonawców
pracujących naraz nad jednym pismem, czego samo pole rodzaju autora (gruby
podział człowiek-wykonawca) nie rozdziela.

Stempel wykonawcy na zmianie śledzonej idzie po odłożeniu czynności, bo niesie
także jej kod dla wiązania — tak samo jak w obszarze schowka i różnicy wersji,
gdzie drugiej drogi zakładania wiersza zmiany Studio nie ma. Bez tego wiązania
cofnięcie czynności zostawiłoby zmianę śledzoną wiszącą w powietrzu, a operator
widziałby do rozstrzygnięcia zmianę, której już nie ma. Brak tabel obszaru
kontroli pracy nie unieważnia samej czynności — postać jest już zapisana,
i odmowa w tym miejscu byłaby nieprawdą o tym, co się stało.

Przesunięcie zakotwiczeń po zmianie długości treści obejmuje obiekty, pola,
aparat, sekcje i blokady, bo bez niego przypis wstawiony na stronie pierwszej
wskazywałby po dopisaniu akapitu na zdanie ze strony drugiej — cichą szkodę,
po której nikt nie wie, kiedy dokument się rozjechał.

Zapis drzewa sprzed czynności przy wczytaniu postaci jest jedynym miejscem,
w którym da się go jeszcze zapisać: po czynności drzewo jest już zmienione,
więc dziennik nie miałby skąd go wziąć jako stan „przed". Odczyt płaci za to
jednym zapisem drzewa do napisu — tą samą cenę dziennik płaci już za stan po
czynności.

## budowa/server/internal/core/adapter_modul_design_fotografia_maski.go

Cztery czynności tego pliku mają wariant neuronowy lepszy od rachunku: powiększenie, odcięcie tła, domalowanie i rozszerzenie kadru. Gdy kanał modelu obrazowego działa, liczy kanał; gdy nie działa, liczy rachunek zawarty w tym pliku, a odpowiedź podaje, którą drogą policzyła, w polu `computedBy`. Obie drogi oddają piksele, różni je jedynie jakość wyniku.

Odcięcie tła rachunkiem opiera się na barwie tła odczytanej z obwodu obrazu i na rozroście obszaru od brzegów, dzięki czemu jasny przedmiot na jasnym tle nie znika w środku kadru, choć tło ma tam tę samą barwę. Metoda działa poprawnie na zdjęciach produktowych i na grafice na jednolitym tle; na portrecie w tłumie wynik jest gorszy niż wynik kanału modelu, dlatego pole `transparentShare` jest pomiarem — pozwala rozpoznać, jaki udział obrazu zniknął.

## adapter_modul_studio_agenci.go

Zajęcie fragmentu nie jest blokadą. Blokada Operatora jest trwała, skierowana
przeciw wykonawcom i zdejmuje ją wyłącznie Operator; zajęcie fragmentu jest
chwilowe, wygasa samo i chroni przed drugim wykonawcą. Zlanie ich dałoby albo
blokadę, która wygasa, czyli żadną, albo zajęcie, którego nikt nie zdejmie po
agencie ubitym w pół pracy. Zajęcie ma dlatego własną tabelę i własny czas
wygaśnięcia, liczony zegarem bazy — zegar rdzenia rozjechałby się z kolumną
wygaśnięcia przy pierwszej różnicy strefy.

Odmowa zajęcia nazywa wykonawcę, a nie mówi „zajęte". Wykonawca, który dostał
„zajęte", nie wie, czy czekać, czy odstąpić; odmowa nazywa więc, kto trzyma
fragment i do kiedy, i to samo zdanie widzi Operator w wykazie zajęć.

Nastawy pętli wykonawczej i pracy wielu agentów są domyślnie wyłączone,
włącza je jawne, odwracalne ustawienie Operatora. Wartości jadą tą samą
tabelą ustawień, którą jedzie cała platforma, tym samym rozstrzyganiem
dziewięciu poziomów — osobny magazyn nastaw Studia byłby drugim miejscem,
w którym trzeba by szukać wartości obowiązującej. Pierwszeństwo poziomów
należy do pakietu konfiguracji, rdzeń o nie nie pyta; gdy rozstrzygacz nie
został podany przy montażu, nastawy czyta się wprost z zasięgu okna, a potem
globalnego, węższy przed szerszym, tą samą zasadą.

Domyślna ważność zajęcia wynosi dziewięćdziesiąt sekund, bo tyle mieści się
jedna czynność modelu na fragmencie, a wykonawca ubity w pół pracy nie trzyma
akapitu dłużej niż półtorej minuty.

Brak zapisu nastawy liczbowej zostaje brakiem, nie zerem: zero znaczyłoby
„ani jeden wykonawca", a brak znaczy „bez granicy ustawionej przez Operatora".

## adapter_modul_studio_tabele_pomocniki.go

Tabela siedzi w drzewie postaci dokumentu, a jej miejsce w kolejności czytania
niesie blok rodzaju tabela ze wskazaniem tabeli. Osobnego wiersza w bazie tabela
nie ma i mieć nie powinna: nie pyta się, które tabele są nieświeże, tak jak pyta
się o przypisy — tabelę czyta się razem z dokumentem, więc jedzie jego drzewem.

Siatka komórek jest pełna, czyli obejmuje wiersze razy kolumny, zamiast trzymać
rzadko tylko komórki wypełnione. Komórki trzymane rzadko wymagałyby przy każdym
odczycie zgadywania, czy komórki nie ma, bo jest pusta, czy bo została wchłonięta
scaleniem. Komórka wchłonięta scaleniem niesie oznaczenie i zostaje na swoim
miejscu, więc wstawienie wiersza w środek scalenia ma co przesunąć, a nie zgaduje.

Szerokości kolumn nigdy nie schodzą do zera: tabela po scaleniu komórek ma mieć
szerokości policzone. Dlatego każda czynność zmieniająca budowę tabeli przelicza
szerokości od nowa, zamiast zostawiać tablicę krótszą niż liczba kolumn.

Sortowanie tabeli używa kolacji pisma polskiego z biblioteki `golang.org/x/text`
zamiast własnej tablicy znaków. Porównanie napisów bajt po bajcie stawia „ł" za
„z", bo tak leżą one w Unikodzie — wykaz nazwisk wychodziłby wtedy z
Łukasiewiczem na końcu, za Zawadzkim. Własna tablica znaków polskich
rozwiązałaby to prowizorycznie i rozjechałaby się na pierwszym wyrazie z
ligaturą, z apostrofem albo z cyfrą. Siła porównania sortowania jest
drugorzędna: różnica wielkości liter nie decyduje o kolejności, a różnica znaku
diakrytycznego decyduje, bo „laska" i „łaska" są dwoma różnymi wyrazami i mają
stanąć osobno.

## budowa/server/internal/core/adapter_modul_design_makiety_zrzut.go

Obie komendy tego pliku oddają ramkę i warstwy, nie plik graficzny, bo makieta ma być czymś, w czym Operator pracuje dalej, a obrazu ekranu nie da się przesunąć piórem. Rozkład na sekcje liczy rdzeń własną regułą; kanał modelu wchodzi wyłącznie tam, gdzie dodaje wiedzę, przy rozpisaniu opisu na nazwy sekcji. Bez wskazanego kanału rdzeń rozkłada opis regułą, która rozdziela go na części i rozpoznaje nazwy sekcji typowe dla makiet ekranów; wynik jest zawsze, na każdej maszynie i bez konta u zewnętrznego dostawcy. Ze wskazanym kanałem rdzeń pyta model o rozpisanie tego samego opisu, a odpowiedź nieczytelna albo błąd kanału nie kończy komendy, bo makieta ma powstać niezależnie.

Komenda design.mockup.import rozkłada zrzut na obszary rachunkiem własnym na pikselach: barwa tła z obwodu obrazu, progowanie odstępstwa od tła, spójne bloki, prostokąty otaczające. Treść napisów czyta program tesseract, składnik pakietu serwera wołany tą samą drogą, co w warsztacie Studia i w module Translate; czytnika liter w czystym Go nie ma, więc plik ma nazwany wyjątek zapory izolacji na własne wywołanie programu zewnętrznego.

Pole recognizeText domyślnie znaczy tak, a układ makiety nie ma prawa zależeć od odczytu liter. Gdy pole podano wprost jako prawdziwe, brak programu jest odmową nazwaną, nie pustym odczytem udającym, że napisów nie było. Gdy pole jest pominięte, obowiązuje domyślne włączenie, ale odczyt jest dodatkiem do układu: brak programu daje wtedy makietę z układem, a każda linia tekstu niesie w adnotacji zdanie o tym, że treści nie odczytano i dlaczego. Obszary, których rdzeń nie rozłożył na elementy, wracają w polu unrecognizedRegions jako bilans zamiast ciszy.
## adapter_modul_design_makiety.go

Wykaz nastaw urządzeń stoi w jednym miejscu i czytają go dwie czynności:
zakładanie ramki, które wypełnia wymiary z nazwy nastawy, oraz wykaz ramek,
który oddaje ten sam wykaz w polu `devicePresets`. Osobny wykaz w kliencie
rozjechałby się przy pierwszej poprawce, a Operator dostałby ramkę telefonu
o wymiarach tabletu.

Ramka wyznacza obszar wydania, a kompozycja pozostaje samym płótnem: to
rozmiar ramki zmienia się przy sprawdzaniu układu na innym urządzeniu.
Warstwa należy najwyżej do jednej ramki, bo warstwa w dwóch ramkach naraz
musiałaby przy `design.frame.resize.apply` przyjąć dwa różne położenia.

Więz responsywny przelicza się rachunkiem, nie zapowiedzią: `design.frame.
resize.apply` liczy nowe położenia warstw z kotwic i oddaje warstwy po
przeliczeniu, odczytane z bazy. Warstwa bez więzu zostaje tam, gdzie była —
brak więzu jest rozstrzygnięciem Operatora, nie luką do wypełnienia domysłem.

Prototyp niesie bilans zamiast ciszy: `design.prototype.get` niesie
`unreachableFrameIds`, wykaz ramek, do których nie prowadzi żadne połączenie,
bo prototyp z ramką osieroconą wygląda w oknie jak prototyp kompletny, a
dopiero przy przejściu okazuje się, że tam nie da się dojść. `design.
constraint.set` niesie `changed`, liczbę więzi rzeczywiście zmienionych, nie
liczbę nadesłanych.

Znaczenie kotwic więzu responsywnego przy zmianie rozmiaru ramki: `start`
trzyma odległość od krawędzi początkowej, lewej albo górnej; `end` trzyma
odległość od krawędzi końcowej; `center` trzyma środek warstwy w środku
ramki; `stretch` trzyma obie odległości naraz, więc warstwa rośnie razem
z ramką; `scale` skaluje położenie i rozmiar warstwy proporcjonalnie do
zmiany rozmiaru ramki.

Usunięcie ramki nie usuwa jej warstw: warstwa jest bytem kompozycji, a ramka
tylko ją grupowała, więc warstwy zwolnione wracają w `releasedLayerIds`
zamiast zniknąć bez śladu.

Wynik układu automatycznego jest zapisany, nie policzony na boku: układ
automatyczny jest zmianą kompozycji, a nie podglądem układu.

Instancja komponentu wskazuje komponent kluczem obcym i czyta z niego
warianty, więc zmiana komponentu dochodzi do wszystkich jego instancji bez
osobnego zapisu; liczba instancji w odpowiedzi na zapis komponentu pochodzi
z odczytu po zapisie, nie z żądania, bo żądanie nie zna tej liczby.

Przy wskazanej ramce początkowej rachunek nieosiągalności jest przejściem po
grafie od niej: nieosiągalna jest ramka, do której nie prowadzi żadna droga,
a nie tylko ta bez połączenia wchodzącego wprost. Bez wskazania początku
nieosiągalna jest ramka bez ani jednego połączenia wchodzącego, bo nie
wiadomo, skąd Operator zaczyna przegląd prototypu.

## adapter_modul_studio_aparat.go

Element aparatu ma swój wiersz w tabeli `element_aparatu_studio`, bo się o niego
pyta wprost: „które spisy są nieświeże", „ile jest przypisów", „czy indeks
zgadza się z treścią". Drzewo postaci niesie go tylko jako odczyt złożony
przez `postacZlozAparat`; zapis idzie zawsze wierszem, dlatego `postacZapisz`
wycina aparat z drzewa — dwa zapisy jednego bytu rozjechałyby się przy
pierwszej poprawce.

Spis treści, spisy ilustracji i tabel oraz indeks muszą być odświeżalne
i mieć znacznik nieświeżości — kolumnę `nieswiezy`, a nie domysł okna: zmiana
nagłówka zapala znacznik natychmiast, żeby operator wiedział, że spis pokazuje
stan sprzed zmiany.

Numeracja przypisów, podpisów i powołań liczy się od miejsca w treści, nie od
kolejności zapisu: przypis wstawiony przed innym ma przenumerować oba, dlatego
numer nie jest nadawany przy zapisie, a przy odświeżeniu, z porządku
zakotwiczeń. Numer nadany przy zapisie byłby numerem, który po pierwszym
wstawieniu w środek dokumentu przestaje być prawdziwy.

Rachunek stron akapitów bierze geometrię z nastaw każdej sekcji z osobna, nie
ze stałego rozmiaru strony: akapit należy do sekcji, której zakres go
obejmuje, a sekcja rozpoczynająca się od nowej strony przerywa rachunek
wierszy i zaczyna kartkę od nowa. Sekcja `continuous` łamania nie przerywa,
bo jej sensem jest ciągłość — nowe nastawy obowiązują wtedy od miejsca cięcia,
a nie od nowej kartki. Rachunek idzie tym samym silnikiem, którym jedzie
podgląd wydruku (`krojWyrysuStudia`, `zlamWierszStudia`), bez wyrysu obrazów,
bo do numeru strony obraz nie jest potrzebny — drugi rachunek łamania dałby
spis treści wskazujący inne strony niż podgląd, co jest gorsze niż brak
numerów.

`newPage`, `evenPage` i `oddPage` przerywają stronę na początku sekcji;
parzystości i nieparzystości rachunek na razie nie dopełnia pustą kartką, bo
wyrys jej nie rysuje i numer wskazywałby stronę, której w podglądzie nie ma.
`continuous` i `newColumn` nie przerywają: pierwsze z zamysłu, drugie dlatego,
że wyrys składa jedną kolumnę.

## adapter_modul_studio_strona_pomocniki.go

Wykaz formatów nośnika stoi raz w rdzeniu, w wykazie nośników druku wspólnym
dla obszaru Design i obszaru Studio. Studio ten wykaz czyta, nie kopiuje:
zmienna leży w tym samym pakiecie, więc odczyt nie wymaga wejścia w plik innego
modułu ani drugiego wykazu wymiarów.

Nagłówek dokumentu mieszka w sekcji, także dla dokumentu bez wyraźnego podziału
na sekcje. Kontrakt trzyma nagłówki i stopki w sekcji, wedle zasięgu (strony
zwykłe, pierwsza strona, strony parzyste), bo tak stawia je zlecenie i tak są
w pakiecie biurowym. Dokument bez ani jednej sekcji nie miałby więc gdzie
trzymać własnego nagłówka. Dlatego czynność na nagłówku, numeracji i znaku
wodnym bez wskazania sekcji sięga do sekcji pierwszej, a gdy dokument nie ma
żadnej, zakłada jedną, obejmującą całą treść — dokument jest wtedy jedną
sekcją, dokładnie jak w pakiecie biurowym po założeniu nowego pisma. Drugie
miejsce na nagłówek całego dokumentu byłoby drugą prawdą, którą pierwsza
zmiana sekcji by rozjechała.

Bilans przeliczenia układu po zmianie nośnika przelicza treść zamiast ją
obcinać, ponieważ sama cisza byłaby najgorszą odpowiedzią: dokument
zobaczyłby obcięty załącznik dopiero na wydruku. Samo zgłoszenie bez
przeliczenia jest niewiele lepsze — zostawiałoby dokument w stanie, którego
nie da się wydrukować, i kazałoby poprawiać każdą tabelę osobno.

## adapter_modul_studio_blokady_zapora.go

Sprawdzenie blokad stoi w rejestrze komend, nie w poszczególnych
obsługiwaczach. Wymaganie mówi: sprawdzenie stoi na drodze każdej komendy
zmieniającej dokument, po stronie serwera, przed dotknięciem treści. Komend
zmieniających dokument Studio jest ponad setka i pisze je czterech
wykonawców naraz — wywołanie sprawdzenia w każdym obsługiwaczu z osobna
znaczyłoby sto miejsc do pominięcia przez pomyłkę, a każde pominięcie to
cicha dziura w blokadzie. Zapora wpięta w rejestr obejmuje wszystkie te
komendy jednym warunkiem, w tym te, których jeszcze nikt nie napisał, bo
działa po nazwie rodziny komend, nie po wykazie obsługiwaczy. Rejestr jest
jedynym miejscem, w którym rdzeń rozstrzyga, co wykonać, więc jest też
jedynym miejscem, przez które przechodzi każde wywołanie — także wywołanie
modelu, bo model woła komendy tą samą drogą co klient.

Sprawdzenie idzie dwiema drogami wedle kształtu komendy. Komenda na
fragmencie niesie zakres znany przed wykonaniem: zakres stykający się
z blokadą wiążącą kończy się odmową nazwaną i obsługiwacz nie rusza. Komenda
na całym dokumencie zakresu przed wykonaniem nie niesie — powstaje dopiero
z rachunku obsługiwacza. Odmowa całości byłaby tu nieproporcjonalna, więc
zapora zakłada kopię zapasową, puszcza obsługiwacza, a potem uzgadnia wynik
z blokadami: fragmenty zablokowane wracają do brzmienia zastanego,
a odpowiedź dostaje bilans pominięć. Blokada nie zostaje przy tym naruszona
na zewnątrz, bo uzgodnienie zamyka się w tym samym wywołaniu.

Bilans pominięć dopisuje się do odpowiedzi tej samej komendy, którą wywołanie
zawołało, inaczej pominięcie zostałoby przemilczane. Dopisuje pole balance do
gotowego ładunku JSON, zamiast żądać od każdego z czterech wykonawców, żeby
dołożył u siebie to samo pole i pamiętał o nim w każdej nowej komendzie.

Wykaz komend bez zmiany treści jest wykazem wyjątków, nie wykazem objętych,
celowo: gdyby zapora obejmowała wykaz komend zmieniających, komenda dopisana
i niewpisana do wykazu byłaby cichą dziurą w blokadzie. Wykaz wyjątków myli
się w drugą stronę — komenda nowa jest domyślnie sprawdzana, a najgorsze, co
może wyjść, to sprawdzenie zbędne przy odczycie. Czynności markup.add,
markup.remove, clipboard.copy, clipboard.paste i diff.hunk.apply sprawdzają
blokadę same, dokładniej niż zrobiłaby to zapora: markup.add jest jedyną
drogą wykonawcy do fragmentu pod blokadą i musi przechodzić, a wyróżnienie
barwą sprawdza blokadę osobno; markup.remove zdejmuje barwę z zakresu
własnego znakowania, którego zapora z żądania nie zna; clipboard.paste jedzie
polem offset, nie zakresem; diff.hunk.apply niesie w żądaniu zakres wersji
źródłowej, nie treści bieżącej, więc uzgodnienie liczy się osobno, na treści
bieżącej.

Kopia zapasowa przed czynnością nieodwracalną obejmuje trzy przypadki:
przyjęcie wszystkich zmian modelu, zamianę w całym dokumencie i zmianę
formatu nośnika. Zapora zakłada kopię przed każdą czynnością na całym
dokumencie, bo wszystkie trzy wchodzą tą drogą, a dołożenie do wykazu
czwartej nie powinno wymagać osobnego pamiętania o kopii.

Siatka śladu autora dopisuje ślad tam, gdzie obsługiwacz go nie odłożył,
zamiast polegać na zaufaniu do obsługiwaczy. Powodem jest zmierzony
przypadek: wywołanie zapisu dokumentu przez wykonawcę zmieniało treść
i nie odkładało ani zmiany śledzonej, ani wpisu dziennika, więc praca modelu
wchodziła do pisma niewidzialna dla przełącznika pokazującego działania
modelu. Naprawa nie mogła stanąć w obsługiwaczu tej jednej komendy, bo komend
zmieniających treść jest w Studiu ponad setka i pisze je czterech
wykonawców, a każda nowa mogłaby przeoczyć ślad tak samo. Siatka mierzy więc
skutek — czy treść się zmieniła — i dopisuje ślad tylko wtedy, gdy
obsługiwacz go nie odłożył sam. Siatka nie zastępuje śladu odkładanego przez
obsługiwacza: tamten zna zakres i rodzaj zmiany dokładnie, a siatka zna
tylko to, że treść jest inna, dlatego jej zakres to zakres różnicy treści,
nie zgadywane miejsce.

## budowa/server/internal/core/adapter_modul_studio_ooxml_test.go

Materiał próbny wpisuje się tu wprost jako archiwum ZIP ze składnikami XML w postaci, w jakiej wychodzą z pakietu biurowego, zamiast składać go własnym składaczem. Gdyby plik próbny powstawał składaczem rdzenia, sprawdzian mierzyłby zgodność składacza z własnym rozbiorem, czyli że rdzeń czyta to, co sam napisał — taki pomiar przechodzi także wtedy, gdy oba końce mylą się w ten sam sposób, a plik jest dla docelowego programu biurowego nieczytelny. Rozbiór ma zdać egzamin z cudzego pliku, nie ze swojego.

Porównanie idzie po odczycie, nie po bajtach, bo ten sam dokument da się zapisać na wiele poprawnych sposobów: inna kolejność węzłów, inne nazwy stylów automatycznych, inne zaokrąglenie jednostek. Bajt w bajt nie zgodzi się nigdy i nie ma się zgodzić; miarą jest to, czy po wczytaniu wyniku postać jest ta sama — nazwa stylu akapitu, orientacja i marginesy sekcji, wymiary tabeli, jej wiersz nagłówkowy, scalenie komórek i szerokości kolumn.

Plik wyklucza pięć rodzajów szkody: wczytanie dokumentu, po którym w postaci stoi sam tekst, a styl, sekcja i tabela przepadły; wydanie dokumentu, które zapisuje treść i gubi postać dokładnie w miejscu, które sprawdzian bada; wydanie do formatu uboższego, które o stracie milczy; wydanie wielostronicowe oddające jedną stronę, choć odpowiedź mówi inaczej; oraz plik wyjściowy niosący warstwę znakowania sesji, czyli komentarze i wyróżnienia w piśmie wysłanym na zewnątrz.

## budowa/server/internal/core/adapter_modul_design_fotografia.go

Każda czynność tego pliku zapisuje nowy zasób i wskazuje źródło polem variantOfAssetId; oryginał nie jest nigdy nadpisywany, dzięki czemu łańcuch edycji da się przejść wstecz do zdjęcia wniesionego pierwotnie, a cofnięcie nie wymaga historii operacji w pamięci okna. Obok wariantu powstaje wiersz czynności z nastawami, którymi poszła, bez którego odczyt historii pokazywałby wykaz obrazków bez słowa o tym, co je różni.

Wynik wychodzi jako PNG, bezstratnie i z kanałem krycia: format stratny traciłby jakość przy każdym ogniwie łańcucha edycji i nie uniósłby przezroczystości po odcięciu tła. Format wyjścia zmienia się dopiero przy wydaniu zasobu, gdzie strata jest jednorazowa i świadoma.

Cztery czynności — powiększenie, odcięcie tła, domalowanie i rozszerzenie kadru — mają dwie drogi rachunku i pole computedBy mówi, którą poszły. Gdy kanał obrazowy jest wskazany, liczy kanał i odpowiedź niesie computedBy: kanalModelu; gdy pole jest pominięte, liczy rachunek wkompilowany i odpowiedź niesie computedBy: rachunekRdzenia. Zejście na pierwszy dostępny kanał przy polu pominiętym wysyłałoby materiał do zewnętrznego dostawcy, o którego nikt nie prosił, dlatego kontrakt tych czterech pól rozstrzyga wyłącznie wskazanie, nigdy dostępność. Wymiar wyniku jest zawsze tym, który czynność obiecała, także na drodze kanału: kanał oddaje obraz o rozmiarze ze swojej nastawy, a sprowadzenie do zamówionego wymiaru wykonuje rachunek rdzenia, natomiast treść pozostaje kanału. Niepowodzenie kanału wskazanego wprost jest odmową, nie cichym zejściem na rachunek, bo wynik policzony inną drogą byłby odpowiedzią na inne żądanie.

## adapter_modul_studio_style.go

Zmiana stylu nazwanego przestawia wszystkie miejsca, które go używają — to nie
jest wygoda, to jest cała jego treść: gdyby stosowanie stylu kopiowało jego
postać na fragment, dokument o dwustu nagłówkach wymagałby dwustu poprawek.
Dlatego fragment i akapit trzymają nazwę stylu, a nie jego postać, a postać
skuteczna liczy się przy odczycie: styl, jego styl nadrzędny i dalej po
łańcuchu dziedziczenia, a na końcu postać własna fragmentu, która ma
pierwszeństwo. Bilans zmiany stylu oddaje liczbę przestawionych miejsc, żeby
zmiana nie była czynnością, po której nie wiadomo, czy coś się stało.

Postać znaku fragmentu liczy się z trzech warstw w kolejności pierwszeństwa:
postać znaku stylu akapitu, postać znaku stylu znaku nałożonego na fragment i
postać własna fragmentu. Styl akapitu, na przykład nagłówek poziomu pierwszego,
niesie też postać znaku — stopień pisma i pogrubienie — tak jak w każdym
pakiecie biurowym. Gdyby postać znaku płynęła wyłącznie ze stylu znaku,
nagłówek dostawałby stopień tekstu zasadniczego, a zmiana stylu nagłówkowego
nie ruszałaby ani jednej litery. Warstwa bliższa fragmentowi ma pierwszeństwo:
pogrubienie zdjęte ręcznie z jednego słowa nagłówka zostaje zdjęte, choć styl
nagłówka pogrubia.

## adapter_modul_studio_wejscie_skutek_test.go

Ten plik sprawdza trzy szkody, które moduł wejścia i zapisu ma wykluczyć.
Pierwsza: konwersja PDF meldująca powodzenie i oddająca dokument okaleczony
bez ani jednego słowa o tym, czego nie odzyskała, albo PDF ze samych skanów
przepuszczony jako skonwertowany, czyli dokument pusty podany jako gotowy do
pracy. Druga: kopia dokumentu będąca drugim odwołaniem do tego samego bytu,
po której poprawka w kopii zmienia oryginał. Trzecia: zapis dokumentu gubiący
postać dokumentu, przez co praca modelu nad postacią stawała się niewidoczna
dla rdzenia po zapisie.

Materiał próbny PDF powstaje biblioteką pdfcpu, tą samą, którą rdzeń go
czyta. To jest świadome: mierzona jest uczciwość bilansu konwersji, a nie
zgodność dwóch bibliotek PDF między sobą. Tekst materiału jest zapisany
wprost jako operator pokazania tekstu, więc warstwa tekstowa jest w nim
prawdziwa, a nie pozorna.

Bilans konwersji PDF musi nieść policzone strony, w tym rozbicie na strony
z warstwą tekstową i bez niej, oraz zdanie o stanie wyniku, bo PDF nie niesie
struktury akapitu ani tabeli wprost — odzyskanie jest odtworzeniem, nie
odczytem. Liczba stron w bilansie jest porównywana z liczbą policzoną w pliku
niezależną biblioteką, a nie brana na słowo. PDF bez warstwy tekstowej,
złożony wyłącznie ze skanów, nie może zostać udawany jako skonwertowany:
dokument pusty oddany jako gotowy do pracy byłby najgorszą możliwą
odpowiedzią, bo dalsza praca zaczęłaby się na pliku bez treści. Bilans
zaznacza wtedy potrzebę rozpoznania pisma, a odpowiedź wskazuje pozycję
kolejki rozpoznania, nazwaną drogą dalszą zamiast samej odmowy.

Miara osobności kopii sprawdza obie warstwy dokumentu, treść i postać, bo
kopia dzieląca postać z oryginałem jest tak samo zepsuta jak kopia dzieląca
treść, tylko trudniej to zauważyć. Przeniesienie historii wersji przy
kopiowaniu jest wyborem jawnym, nie zachowaniem zaszytym na stałe, a liczba
przeniesionych wersji wychodzi kontraktem, więc sprawdzian mierzy tę liczbę,
nie samo powodzenie komendy.

Zapis dokumentu ma przenosić postać, nie samą treść: komenda zapisu
przyjmowała pola treści dokumentu bez pola postaci, więc postać przy każdym
zapisie ginęła, a praca nad nią stawała się niewidoczna po stronie rdzenia.
Miara sprawdza postać podaną polem form, odczytaną ponownie osobnym
wywołaniem — arkusz stylów, nastawy strony, tabelę i styl akapitu — a zapis
bez pola form nie ma prawa zetrzeć postaci zastanej, bo brak pola znaczy
brak zmiany postaci, nie postać wyzerowaną.

## adapter_modul_studio_format.go

Praca na fragmencie jest osią tego pliku: postać nałożona na zaznaczenie ma
zmienić wyłącznie to zaznaczenie. Rachunek granic stoi w rachunku postaci
(`postacRozetnij`, `postacFragmentyZakresu`) i ten plik się go tylko woła —
dwa liczenia granic zaznaczenia rozjechałyby się przy pierwszej poprawce.

Zabrana postać malarza leży wierszem tabeli `postac_malarza_studio`, nie
w pamięci procesu, bo między zabraniem a położeniem stoją dwie osobne
komendy — postać musi przeżyć czas między nimi, także przeładowanie rdzenia.
Wiersz wygasa po godzinie: malarz jest narzędziem jednej czynności, a postać
zabrana wczoraj i położona dziś byłaby zaskoczeniem, nie pomocą. Godzina
obejmuje ciąg czynności modelu, który zabiera postać, robi kilka innych
rzeczy i kładzie ją na końcu — dłuższy czas czyniłby z malarza magazyn,
krótszy odbierałby postać w połowie roboty. Brak tabeli malarza jest brakiem
montażu rdzenia i mówi to wprost, zamiast cicho wracać do pamięci procesu:
cichy powrót dawałby malarza, który u jednego operatora przeżywa restart,
a u drugiego nie.

Zamiana wielkości liter jest zmianą treści, nie postaci: „Kowalski" zamienione
na „KOWALSKI" ma inne litery, a nie inny krój, dlatego idzie drogą zamiany
treści i odkłada zmianę śledzoną rodzaju wstawienie, a nie formatowanie —
inaczej cofnięcie nie miałoby czego przywrócić.

Odmowa przy braku wiersza malarza nazywa powód prawdziwy: wpis wygasły
i wpis nieistniejący są dla warstwy danych tym samym brakiem wiersza, więc
odmowa mówi o obu i podaje drogę wyjścia — zabranie postaci ponownie.

Wyszukiwanie w `ZamienZPostacia` działa też samą postacią, bez tekstu:
„wszystkie fragmenty czerwone zamień na czarne" jest zwykłym poleceniem
redakcyjnym. Zaznaczanie wedle podobnego formatowania jest czynnością
obowiązkową zlecenia, bo bez niej nie ma pracy na postaci fragmentami
w długim dokumencie. Podobieństwo postaci znaku sprawdza cechy, które
operator widzi gołym okiem — krój, stopień i grubość; barwa wyróżnienia
i język nie liczą się do podobieństwa. Wzorzec przy wyszukiwaniu samą
postacią sprawdza wyłącznie cechy, które sam niesie: wzorzec „barwa
czerwona" ma trafiać we wszystko czerwone, niezależnie od kroju.

## budowa/server/internal/core/adapter_modul_studio_cyfryzacja.go

Ingest/OCR Panel różni się od rodziny komend document.text.extract tym, że tamta jest narzędziem modelu — jedno wywołanie, jeden napis na wyjściu, brak stanu — podczas gdy panel jest stanowiskiem, na którym materiał czeka w kolejce, wraca do niej po poprawce obrazu, dostaje korektę słowa i dopiero na końcu staje się dokumentem. Dlatego kolejka ma tabelę, a wynik rozpoznania niesie słowa wraz z pewnością, bez których korekta rozpoznania nie ma czego poprawiać.

Słowa i pewność oddaje sam program Tesseract wyjściem w postaci tabelarycznej: jeden wiersz na słowo, z ramką i pewnością w skali od zera do stu. Rdzeń niczego tu nie szacuje ani nie dopowiada — pewność pozycji jest średnią pewności jej słów, a gdy Tesseract nie oddał ani jednego słowa, pewności nie ma wcale i pole zostaje puste. Pusta pewność i pewność zerowa to dwie różne rzeczy.

Cztery nastawy kontraktu opisują obróbkę wstępną skanu — prostowanie skosu, odszumianie, progowanie i przycinanie marginesów — i prowadzi je program unpaper, napisany dokładnie do tego zadania; rdzeń nie liczy tego sam, bo skos wykrywa się przemiataniem obrazu pod kątem, a nie jedną pętlą po pikselach. Gdy żadna z czterech nastaw nie jest włączona, unpaper w ogóle nie rusza, bo przebieg bez niczego i tak przepisałby obraz przez konwersję do postaci PNM i z powrotem, płacąc uruchomieniem procesu za wynik, o który nikt nie prosił. unpaper ma wszystkie swoje filtry włączone domyślnie, a kontrakt ma je domyślnie wyłączone, więc rdzeń wyrównuje te dwie domyślności: filtr, o który nikt nie prosił, jest wyłączany jawnie, inaczej jedna włączona nastawa włączyłaby po cichu także pozostałe. Filtr czarnej ramki, którego kontrakt nie ma wcale, jest wyłączony zawsze.

## adapter_modul_studio_wejscie_dane.go

Droga zapisu postaci stoi w adapter_modul_studio_postac.go, a ten plik ją
woła, dokładając jedno: wiersze, których zapis obszaru postaci świadomie
z drzewa wycina — arkusz stylów, sekcje, obiekty i pola. Wycina je, bo dla
dokumentu już istniejącego one w bazie stoją i drzewo nie ma być ich drugą
prawdą. Dokument wnoszony z pliku albo kopiowany jest przypadkiem odwrotnym:
wierszy jeszcze nie ma, a postać przyszła z zewnątrz. Gdyby ten odcinek
zapisał samo drzewo, dokument po ponownym wczytaniu tracił arkusz stylów
wniesiony z docx, czyli dokładnie tę cichą stratę, którą utrwalenie postaci
ma wykluczyć. Styl fabryczny dokumentu zapisuje się jako fabryczny wtedy
i tylko wtedy, gdy tak przyszedł: arkusz wniesiony z docx jest arkuszem
dokumentu, nie arkuszem platformy, i ma dać się zmienić bez odmowy o stylu
fabrycznym.

Kontrakt warstwy danych tego odcinka jest węższy niż cały interfejs
repozytorium Studia celowo: repozytorium deklaruje się w całości w jednym
pliku, a dopisanie tam metod byłoby wejściem w plik cudzego odcinka. Rdzeń
bierze więc dokładnie te metody, których używa, rzutowaniem dwuwartościowym
— brak nazywa się wprost, zamiast wywracać montaż rdzenia. Ten sam wzór
trzyma odcinek kontroli pracy.

Granica wielkości pliku wejściowego chroni rdzeń przed plikiem, którego
rozbiór wypełniłby pamięć procesu: sto dwadzieścia osiem megabajtów mieści
dokument z setkami obrazów. Kolejność dróg do bajtów materiału jest
rozstrzygnięciem, nie przypadkiem: bajty podane wprost są najpewniejsze, bo
zostały wysłane właśnie teraz, potem zasób rdzenia, potem plik Biblioteki,
a ścieżka na końcu, bo jest ścieżką na maszynie serwera, nie na maszynie
wołającego.

Nazwa pliku ma znaczenie poza samym rozpoznaniem formatu: plik
umowa najmu.docx ma zostać dokumentem o nazwie „umowa najmu", nie
dokumentem bez nazwy, dlatego nazwa dokumentu wyjmuje się z nazwy pliku bez
rozszerzenia.

## budowa/server/internal/core/adapter_modul_design_fotografia_wsad.go

Wsad puszcza ten sam zestaw czynności na każdym zasobie osobno, przez te same uchwyty, którymi jadą czynności pojedyncze; nie ma tu drugiej drogi rachunku, bo gdyby wsad liczył po swojemu, wynik wsadowy różniłby się od pojedynczego. Zasób, którego nie udało się przetworzyć, wraca w bilansie polem failedAssetIds razem z powodem — wsad na stu zdjęciach, z których trzy padły, wyglądałby bez tego pola jak wsad kompletny, a brak wyszedłby na jaw dopiero przy przeglądaniu wyników.

Metadane zasobu są pomiarem z pliku, nie echem wiersza bazy: odczyt czyta nagłówek pliku leżącego w magazynie, biorąc wymiary, format, model barwny, obecność kanału krycia, rozdzielczość i pola EXIF. Wiersz bazy niesie tylko to, co zmierzono przy wniesieniu, a plik może być jedyną prawdą o tym, co użytkownik naprawdę ma.
## adapter_modul_studio_wejscie_tekst.go

Rozpoznanie zapisu znaków jest tu osobną pracą, bo pliki Operatora bywają
starsze niż UTF-8. Pismo urzędowe pisane w Windows-1250, wczytane jako
UTF-8, daje albo błąd, albo tekst z krzaczkami w miejscu liter „ą", „ę",
„ł" — a krzaczki są gorsze niż odmowa, bo model przeczyta je jako słowa
i zacznie na nich pracować. Dlatego zapis znaków rozpoznaje się przed
rozbiorem treści, a rozpoznanie wychodzi kontraktem w bilansie, żeby
Operator wiedział, jak rdzeń odczytał jego plik.

Rozpoznanie idzie bibliotekami `golang.org/x/net/html/charset`
i `golang.org/x/text/encoding`, obiema z rodziny wzorcowej Go, już
obecnymi w drzewie zależności. Kolejność rozstrzygania: znacznik
kolejności bajtów, potem deklaracja w treści — nagłówek XML albo
`meta charset` HTML — potem sprawdzenie poprawności UTF-8, a na końcu
miara rozkładu bajtów rozstrzygająca między stronami kodowymi używanymi
w polskich dokumentach.

Wskazanie zapisu znaków przez Operatora ma pierwszeństwo nad rozpoznaniem: kto
wie, w czym jest jego plik, wie to lepiej niż miara rozkładu bajtów; wskazanie
nazwy nieznanej jest jednak odmową, nie cichym zejściem na UTF-8, inaczej
Operator dostałby krzaczki i myślał, że jego wskazanie zadziałało.

Rozpoznanie po rozkładzie bajtów, między stronami kodowymi jednobajtowymi,
nie jest odczytem deklaracji, tylko zgadywaniem — dlatego wychodzi w bilansie,
żeby Operator wiedział, że rdzeń zgadywał, a nie czytał wprost.

Miara polskości jest miarą, nie dowodem, dlatego nazwa zapisu rozpoznanego
tą miarą idzie do bilansu wniesienia, żeby Operator wiedział, że rdzeń
zgadywał, a nie odczytał deklaracji wprost.

Rozbiór markdown jest własny, wierszowy, a nie przez bibliotekę `goldmark`:
potrzebny jest tu przekład na style nazwane i bloki dokumentu, nie na HTML,
który `goldmark` oddaje, a przekład HTML na postać dokumentu drugą drogą
byłby dłuższy i gubiłby to samo.

Styl „kod" nie stoi w arkuszu domyślnym dokumentu, więc krój stałej
szerokości bloku kodu idzie postacią znaku wprost — to strata nazwana
w bilansie, nie przemilczana.

Szerokości kolumn tabeli wniesionej z markdown liczą się z obszaru pisania,
nie zostają zerowe: tabela markdown szerokości kolumn nie niesie, a zero po
zapisie do docx dałoby kolumny niewidoczne.

Postać znaku płynie w dół drzewa HTML: `<b><i>tekst</i></b>` daje fragment
pogrubiony i pochylony naraz, bo każdy poziom dokłada swoją cechę do postaci
odziedziczonej po przodku; postać liczona osobno na każdym poziomie
zgubiłaby cechę odziedziczoną z zewnątrz.

Odstęp między znacznikami blokowymi HTML nie jest treścią i znika, ale
odstęp między znacznikami tekstowymi jest treścią i zostaje — inaczej
`<b>a</b> <i>b</i>` dałoby „ab" zamiast „a b".

Rozbiór HTML nie pobiera bajtów obrazu wskazanego znacznikiem: pobranie
z sieci należy do czynności wniesienia ze strony (`studio.insert.from.web`),
która ma na to kontekst żądania i granicę czasu — tutaj powstaje tylko
obiekt obrazu wraz z jego adresem, czyli prawda o tym, co rdzeń w tym
miejscu wie.

Postać dokumentu z RTF czyta akapity zapisane rozkazem `\par`, pogrubienie,
kursywę, podkreślenie, stopień pisma zapisany rozkazem `\fsN` w półpunktach,
wyrównanie zapisane rozkazami `\qc`, `\qr`, `\qj`, oraz znaki spoza zakresu
jednobajtowego zapisane jako `\'hh` i `\uN`.

Tabele RTF, zapisane rozkazem `\trowd`, i obrazy RTF, zapisane rozkazem
`\pict`, nie są odzyskiwane: tabela RTF jest ciągiem akapitów z granicami
komórek zapisanymi w rozkazach składu, więc jej odzyskanie byłoby
odtworzeniem układu, nie odczytem struktury.

Znak RTF zapisany szesnastkowo czytamy zawsze stroną windows-1250, typową
dla RTF pisma polskiego, bo nagłówek `\ansicpgN` bywa nieprawdziwy częściej
niż strona kodowa, którą naprawdę użyto do zapisu pliku.

## budowa/server/internal/core/adapter_modul_studio_listy.go

Definicja listy — jej rodzaj, poziomy, znaki wypunktowania, formaty numeracji
i wcięcia — stoi w drzewie postaci dokumentu. Akapit należący do listy trzyma
sam kod listy i numer poziomu, a nie kopię jej nastaw: zmiana znaku
wypunktowania poziomu ma przestawić od razu wszystkie punkty tego poziomu, bez
przechodzenia po akapitach. Przy kopii w akapicie lista czterdziestu punktów
wymagałaby czterdziestu poprawek.

Poziom listy niesie wcięcie wzorcowe, lecz wcięcie skuteczne akapitu musi być
widoczne w jego własnej postaci, aby narzędzie do przestawiania wcięć miało
czym operować. Zastosowanie listy ustawia wcięcie akapitu według poziomu,
a późniejsza zmiana wcięcia narzędziem jest zmianą samego akapitu i poziomu
listy nie rusza — tak samo działa pakiet biurowy, z którym format jest zgodny.

Kontrakt niesie punkt startu na liście i na poziomie, ale nie w akapicie, bo
akapit nie jest miejscem na nastawę listy. Wznowienie numeracji od wskazanego
miejsca jest rozdzieleniem listy: punkty od tego miejsca w dół przechodzą do
listy nowej o tych samych poziomach i własnym punkcie startu. Skutek —
numeracja zaczynająca się od nowa od wskazanego miejsca — jest sprawdzalny
w bazie, a nie udawany polem, którego kontrakt nie ma.

Rodzaj none w czynności zastosowania listy zdejmuje listę i jest czynnością
prawdziwą, nie brakiem: zdjęcie listy z fragmentu, który do żadnej listy nie
należy, wraca odmową nazwaną, aby cisza nie sugerowała zdjęcia listy, której
tam nigdy nie było.

Lista wskazana kodem musi istnieć w dokumencie; wskazanie kodu nieznanego jest
odmową nazwaną, ponieważ cicha zamiana na listę nową dałaby dwie listy tam,
gdzie oczekiwano jednej.

Format prawniczy wielopoziomowy bierze się z pola wzoru numeru, a nie
z domysłu, aby pismo z podstawami prawnymi dało się ponumerować dokładnie
według wzoru urzędowego.

Znak wypunktowania bierze się ze znaku gotowego, dowolnego symbolu, ikony albo
obrazu własnego — wszystkie cztery źródła stoją w kontrakcie. Ikona i obraz
idą zasobem, a nie wklejonym rysunkiem: ikony wystawia moduł Design, a zasoby
przechowuje magazyn rdzenia, więc Studio nie zakłada drugiego rachunku ikon.

Wcięcie i odstęp poziomu idą za poziomem listy: przestawienie wcięcia
w pakiecie biurowym zmienia poziom punktu, a nie samo wcięcie akapitu.
Zmniejszenie poniżej pierwszego poziomu nie zdejmuje listy samo z siebie —
zdjęcie listy jest osobną czynnością zastosowania listy z rodzajem none.

## budowa/server/internal/core/adapter_modul_studio_roznica_wersji.go

Różnica postaci działa jako komenda osobna od różnicy treści, ponieważ rachunek
treści liczy zmianę wierszami i nie dostrzega zmiany kroju ani wcięcia — taka
zmiana nie rusza żadnej litery, więc pozostaje niewidoczna, mimo że w
dokumencie jest zmianą widoczną. Dlatego postać porównuje się cechą po cesze,
obszar po obszarze, a wykaz niesie brzmienie stanu przed zmianą i po niej.

Przeniesienie fragmentu numeruje fragmenty tym samym rachunkiem, którego
używa okno różnicy, aby numer fragmentu w żądaniu wskazywał dokładnie ten
fragment, który widać w oknie. Osobny rachunek numeracji ponumerowałby
fragmenty inaczej, więc wskazanie fragmentu określonym numerem znaczyłoby co
innego dla okna i dla rdzenia.

Blokada fragmentu działa przed dotknięciem treści: fragment objęty blokadą
zostaje w brzmieniu zastanym, a odpowiedź niesie bilans pominięć, tak samo
jak przy zaporze rejestru.

Postać przenoszonego fragmentu jedzie razem z jego treścią, gdy żądanie tego
nie wyłączy, ponieważ samo brzmienie zostawiłoby akapit w kroju bieżącym, a
przeniesienie ograniczone do treści byłoby wtedy niepełne.

Postać runów przenosi się wyłącznie wtedy, gdy brzmienie bloku jest to samo:
run niesie zarówno tekst, jak i postać, a przy różnym brzmieniu podstawienie
runów źródłowych podmieniłoby treść pod pozorem postaci.

Tożsamością bloku przy przenoszeniu postaci jest jego identyfikator: blok,
którego wersja źródłowa nie zna, zostaje w postaci bieżącej. Zgadywanie
odpowiedniości bloków po innej cesze byłoby gorsze niż nieprzeniesienie
postaci, ponieważ nałożyłoby krój obcego akapitu na istniejącą treść.

## adapter_modul_design_wektor.go

Kształt zakładany przez `design.vector.shape.add` wchodzi do bazy od razu jako
komplet węzłów ścieżki z uchwytami, a nie jako prostokąt czy elipsa osobnym
bytem do późniejszej zamiany w ścieżkę: byt pośredni wymagałby komendy „zamień
w ścieżkę”, której kontrakt nie ma, i dawałby na kompozycji dwa rodzaje
kształtu różniące się tym, czego z nimi wolno zrobić.

Operacja logiczna (`ZlozSciezkiLogicznie`) liczy sumę, różnicę, część wspólną
i wykluczenie przez bibliotekę `tdewolff/canvas`, a wynik wraca do bazy znowu
jako węzły, nie jako gotowy zapis SVG — ścieżka po operacji ma dać się ciągnąć
piórem dalej, inaczej pierwsza suma dwóch kół kończyłaby edycję kształtu.
Ścieżki źródłowe usunięte po operacji (`keepSources` bez wskazania) wracają w
`removedPathIds`: pole puste tam, gdzie ścieżki zniknęły, byłoby ciszą, a okno
pokazywałoby kształty, których w bazie już nie ma.

Zamiana tekstu w kontury (`TekstNaSciezce`) jest nieodwracalna dla wyniku:
kontury nie wiedzą, że były literami, więc literówki w nich już nie da się
poprawić. Ścieżka nośna wskazana przez wywołującego rozstrzyga też o początku
tekstu — bierze się z jej pierwszego węzła.

`UsunSciezke` odmawia wyłącznie wtedy, gdy zapis się nie powiedzie: ścieżki,
której nie ma, nie traktuje jak usterki, bo kontrakt pyta polem `removed`, czy
wiersz istniał, nie czy polecenie bazy się udało.

`OczyscSciezki`: ubytek bajtów w `savedBytes` jest pomiarem, nie oszacowaniem —
rdzeń mierzy długość zapisu przed przycięciem precyzji i po nim, więc wartość
ujemna (zapis dłuższy, bo zastana precyzja była wyższa) jest tu prawdą
rachunku, nie usterką.

`WydajWektor`: wydanie wektorowe zostaje wektorem — rasteryzacja odebrałaby mu
jedyną własność, dla której jest wektorem. SVG składa dokument tekstowy w tym
samym pliku, PDF i EPS biblioteka `tdewolff/canvas` wkompilowana w binarium.
PDF i EPS mierzą stronę w punktach typograficznych i liczą oś Y od dołu, a
kompozycja liczy Y od góry — stąd odbicie kształtów przy wydaniu, bez którego
dokument wyszedłby lustrzanym odbiciem tego, co widać na kompozycji.

`stylWydaniaWektoraDesignu`: ścieżka bez wypełnienia i bez obrysu dostaje obrys
włoskowy, bo kształt bez żadnego z dwóch byłby w wydaniu niewidoczny, a plik,
w którym nie widać niczego, wygląda identycznie jak plik uszkodzony.

`uporzadkujBilansDesignu` oddaje nil dla wykazu pustego, żeby pole
niewymagane kontraktu nie weszło do odpowiedzi wcale, a dla wykazu niepustego
ustala stałą kolejność, żeby dwa wywołania tej samej komendy nie różniły się
porządkiem zastrzeżeń.

## adapter_modul_studio_wejscie_pdf.go

PDF nie niesie struktury akapitu ani tabeli, tylko rozkaz postawienia napisu w
danym miejscu strony. Akapit, wiersz tabeli i kolumna są wnioskiem z układu
napisów, a nie zapisem w pliku, dlatego wynik nazywa się odzyskaniem i idzie
z bilansem: ile stron miało warstwę tekstową, ile akapitów odtworzono, ile
układów tabelarycznych rozpoznano, a ile pozostało nierozpoznanych. Dokument
złożony z samych skanów nie przechodzi przez konwersję udawaną — wchodzi do
kolejki rozpoznania pisma, a odpowiedź wprost podaje pozycję tej kolejki.

Rozbiór PDF idzie biblioteką `pdfcpu`, wkompilowaną w binarium rdzenia.
Programów zewnętrznych zdjętych z rdzenia nie wywołuje się — wykaz zdjętych
trzyma zapora rdzenia w pliku `zapora_warsztatu_pdf_test.go`. Biblioteka jest
wkompilowana, bo proces potomny byłby tu regresem: jedno binarium serwera
zamieniłoby się w dwa, doszedłby koszt uruchomienia procesu i ryzyko rozjazdu
wersji. Rozpoznanie pisma ze skanu idzie osobną drogą, kolejką wczytywania
modułu, ponieważ wymaga biblioteki Tesseract, składnika pakietu serwera —
i tak też jest zgłoszone, jako brak w wykazie zależności pakietu.

Poziom nagłówka w odtworzonym dokumencie rozpoznaje się po układzie wiersza:
krótkości, braku kropki na końcu, zapisie wersalikami albo numeracji własnej
pisma źródłowego (na przykład „1.2.3"), a nie po stylu, którego PDF nie niesie.
Ten poziom idzie do bloku akapitu, bo bez niego spis treści złożony po
wniesieniu dokumentu nie miałby czego zebrać, choć nagłówki w tekście stoją.

Tabela powstaje z bloku wierszy tabelarycznych wtedy i tylko wtedy, gdy liczba
kolumn jest zgodna w całym bloku albo różni się o jedną. Blok o rozjeżdżającej
się liczbie kolumn nie jest tabelą, tylko tekstem ułożonym w kolumny, i wtedy
wchodzi do dokumentu jako akapity z zachowanym rozkładem odstępów, a bilans
liczy go jako układ nierozpoznany.

## adapter_modul_design_druk.go

Kompozycja Design Board nie niesie jednostki w kontrakcie — warstwa niesie
same liczby. Część drukarska czyta je jako milimetry i stosuje to
rozstrzygnięcie spójnie w kontroli przeddrukowej, wydaniu i podziale na
kafle, dzięki czemu kompozycja 210×297 jest arkuszem A4, a nie prostokątem
bez nazwanego rozmiaru. Wyrys ekranowy liczy te same liczby jako piksele,
bo tam nie ma nośnika, więc nie ma czego mierzyć w milimetrach.

Komenda design.print.export odmawia wydania przy wadzie o wadze błędu: plik
nie do druku wydany jako gotowy do druku idzie do drukarni, wraca po dniu
i kosztuje nakład. Pominięcie kontroli jest świadomym wyborem operatora
i wraca w odpowiedzi polem preflightSkipped, więc pominięcie nie przechodzi
w ciszy.

Zastrzeżenia kontroli przeddrukowej wychodzą uporządkowane po wadze:
najpierw to, co blokuje druk, potem to, co grozi jakością, w kolejności
stabilnej między dwiema kontrolami tego samego materiału.

Rdzeń nie rozkłada profili ICC i nie przelicza barw przez nie, więc wydanie
w CMYK idzie bez osadzonego profilu, a kontrola przeddrukowa zwraca to jako
zastrzeżenie zamiast milczeć.

Profil wydania domyślny, bez wskazania wprost i bez identyfikatora, niesie
CMYK w rozdzielczości drukarskiej i bez spadu, ponieważ spad dołożony po
cichu zmieniałby wymiar strony.

Publikacja wielostronicowa przepuszczona ze stroną nie do druku wraca
z drukarni tak samo jak pojedynczy arkusz, tylko drożej, więc kontrola
przeddrukowa obejmuje każdą stronę osobno.

Oprawa zeszytowa wymaga liczby stron podzielnej przez cztery, bo arkusz
zgięty na pół daje cztery strony; publikacja niespełniająca tego warunku
dostaje odmowę, a nie ciche dołożenie wakatów, ponieważ strona pusta
w środku książki jest rozstrzygnięciem, nie zaokrągleniem.

## budowa/server/internal/core/adapter_modul_studio_wstawienia.go

Wstawienie obrazu bez zapisanego pochodzenia jest brakiem, nie skrótem: każdy
obiekt niosący bajty ma zapisane, skąd jest — zasób magazynu rdzenia, węzeł
modułu Design, plik Biblioteki, baza zdjęciowa albo adres w sieci. Obraz bez
pochodzenia jest za czas obrazem, o którym nikt nie wie, czy wolno go było użyć.

Obiekt nie nosi bajtów w swoim wierszu: nosi wskazanie zasobu, a bajty leżą
w magazynie pod sumą kontrolną. Droga odłożenia bajtów w magazynie zasobów
rdzenia jest jedna dla całego modułu i nie zakłada się jej drugi raz.

Ikona wchodzi z katalogu modułu Design, tego samego, z którego korzysta
wyszukiwanie ikon w tym module — drugiego katalogu ikon Studio nie zakłada
i mieć nie będzie. Kształt rysowany na miejscu opisuje się rodzajem, rozmiarem
i wyglądem, a nie własnym rachunkiem ścieżek; kształt wymagający ścieżek
edytowalnych wskazuje się węzłem modułu Design, który powstaje osobną
czynnością tego modułu.

Odmowa wstawienia wykresu jest nazwana, nie cicha: rdzeń nie ma rachunku
wykresu i nie udaje, że ma. Wykres składa się w module Design i wstawia się
do dokumentu jako obraz albo węzeł Designu.

Zdjęcie z bazy zdjęciowej leży w magazynie zasobów jak każdy inny zasób,
lecz pochodzenie „baza zdjęciowa” zostaje zapisane osobno, bo mówi o prawach
do użycia obrazu.

Adres źródła w sieci zapisuje się jako pochodzenie, ale bajtów rdzeń stąd nie
pobiera: pobranie treści z sieci ma w rdzeniu własną, osobną drogę — druga
droga pobrania byłaby drugą prawdą o tym, co i skąd weszło do dokumentu.

## budowa/server/internal/core/adapter_modul_studio_blokady_dane.go

Kontrakt danych odcinka jest własny i węższy niż `dane.RepozytoriumStudia`,
ponieważ `dane/studio.go` deklaruje to repozytorium w całości, a dopisanie tam
metod byłoby wejściem w plik cudzego odcinka. Rdzeń bierze więc dokładnie te
metody, których używa, rzutowaniem dwuwartościowym: rdzeń złożony z
repozytorium bez tych tabel pracuje dalej w pozostałych czynnościach Studia, a
brak nazywa się wprost, zamiast wywracać proces przy montażu. Ten sam wzór
trzyma odcinek postaci dokumentu i wyposażenie Terminala.

Rozstrzygnięcie, czyja ręka wykonuje czynność, stoi na jednym pytaniu: czy
komendę zawołał operator, czy wykonawca. Kontrakt daje dwie drogi odpowiedzi:
fakt gniazda (`transport.Tozsamosc.Narzedzia()` — serwer narzędzi modelu
przedstawia się przy nawiązaniu gniazda, czego okno operatora nigdy nie robi
i czego model nie układa sobie sam w treści żądania) oraz pole żądania
(`author`, `agentId`, `agentName`, `subagentId` — jedyna droga, która mówi,
który z wielu wykonawców pracuje, skoro operator zakłada w module Agents
dowolnie wielu agentów, a gniazdo ich nie rozróżnia).

Zasada łącząca obie drogi: wykonawcą jest ten, kogo wskazuje szerszy z dwóch
sygnałów. Czynność jest czynnością operatora wtedy i tylko wtedy, gdy milczą
oba sygnały — gniazdo nie jest serwerem narzędzi i żądanie nie podpisuje się
wykonawcą; wystarczy jeden sygnał, żeby czynność była czynnością wykonawcy.
Samo pole żądania osobno byłoby zaporą, którą model omija pominięciem pola —
to usterka do naprawy, nie zamierzone ograniczenie. Samo gniazdo osobno nie
rozdziela dwóch agentów pracujących naraz, więc podświetlenie zmian
wykonawcy pokazywałoby obu jako jednego.

Tożsamość agenta (kod, nazwa, wersja, podagent) bierze się wyłącznie z pola
żądania, bo gniazdo jej nie niesie; wchodzi także wtedy, gdy rodzaj wyszedł
z gniazda, ponieważ podpis mówi wtedy, który z wielu wykonawców to był.
Zmiana bez wskazanego agenta zostaje poprawna, ponieważ dokumenty i wiersze
sprzed dobudowy tego pola nie mają go, a odmowa przy nich karałaby za wiek
wiersza.

Granica zasady: `Rodzaj=narzedzia` jedzie parametrem nawiązania gniazda, a
tożsamość połączenia jest opisem, nie zaporą. Opis wzięty tu pod uwagę może
czynność wyłącznie zawęzić w prawach (uczynić ją czynnością wykonawcy), nigdy
jej praw nie rozszerza. Zatajenie tego parametru nie otwiera niczego, dopóki
żądanie podpisuje się wykonawcą, ani wtedy, gdy oba sygnały milczą — wtedy
rdzeń nie ma po czym rozpoznać wykonawcy, co jest brakiem wiedzy rdzenia, nie
luką tej zasady. Wiązanie tożsamości serwera narzędzi przy jego starcie leży
poza tym plikiem.

Kontekst wywołania, nie parametr sygnatury, niesie rozpoznanego wykonawcę do
czynności postaci dokumentu, ponieważ takich czynności jest trzydzieści osiem
i wszystkie kończą jedną drogą. Przełożenie podpisu przez trzydzieści osiem
sygnatur znaczyłoby trzydzieści osiem miejsc do pominięcia przez pomyłkę, a
każde pominięcie zmieniałoby zapis na podpisany „nienazwanym" — tym samym
rachunkiem, którym zapora blokad stanęła w rejestrze, a nie w
obsługiwaczach. Kontekst obejmuje też komendy, których jeszcze nikt nie
napisał.

Zakres żądania liczy się w znakach, nie w bajtach, ponieważ kontrakt nazywa
początek fragmentu w znakach, a cięcie po bajtach rozcinałoby polskie litery
dwubajtowe: blokada założona w tekście z „ą" chroniłaby wtedy pół znaku, a
bilans mówiłby o zakresie, którego w treści nie ma. Zakres pusty jest
punktem wstawienia i styka się z blokadą, w której środku leży, ponieważ
wpis w środek zablokowanego cytatu jest zmianą tego cytatu, choć nie usuwa
ani jednej litery; punkt na samej krawędzi blokady nie styka się, ponieważ
dopisanie za zablokowanym fragmentem go nie rusza.

Odmowa blokady nazywa fragment i blokadę zamiast milczeć, ponieważ cicha
bezczynność byłaby najgorszą możliwą odpowiedzią — operator myślałby, że
model wykonał polecenie.

## adapter_modul_studio_pola.go

Pole ma swój wiersz w tabeli `pole_dokumentu_studio` z kolumną `nieswieze`,
odpowiadającą na pytanie, które pola wymagają odświeżenia. Wartość zapisana
w wierszu jest wartością z chwili odświeżenia, a nie z chwili odczytu — inaczej
dokument wydany do PDF-u i dokument otwarty w edytorze pokazywałyby dwie różne
wartości tej samej daty.

Numer strony liczy się tym samym silnikiem łamania, co podgląd wydruku,
funkcją `aparatStronyAkapitow`, i tymi samymi nastawami sekcji. Osobny rachunek
dałby numer strony inny niż widziany w podglądzie, co jest gorsze niż brak
numeru.

Pole obliczane liczy rachunek wkompilowany w tym pliku: liczby, cztery
działania, nawiasy oraz działania na kolumnie tabeli (SUMA, ŚREDNIA, MIN,
MAKS, LICZBA). Wyrażenia, którego rachunek nie zna, nie liczy w przybliżeniu —
odmowa nazywa miejsce i znak, na którym rachunek się zatrzymał.

Rdzeń nie trzyma imienia i nazwiska autora dokumentu, tylko rodzaj autora
ostatniej wersji: Operator albo model. Pole autora oddaje więc to, co rdzeń
wie na pewno, i nazywa wprost to, czego nie wie, zamiast wpisywać nazwę
Operatora, której rdzeń nie ma.

## budowa/server/internal/core/adapter_modul_studio_schowek.go

Schowek jest jeden, wspólny z rodziną `clipboard.*`, a nie osobny dla modułu
Studio: wpis odłożony w jednym miejscu daje się wkleić w drugim — model
odkłada fragment, operator go wkleja, i odwrotnie. Osobny schowek Studia
rozdzieliłby jedną historię na dwie i wklejenie wpisu sprzed kilku ruchów
przestałoby działać między oknem pracy z dokumentem a resztą platformy.
Wiersze idą tą samą tabelą i tym samym rachunkiem odcisku, więc powtórzenie
identycznej treści nie mnoży wpisów.

Pole `formatOnly` w żądaniu kopiowania znaczy „kopiuj postać, nie treść” —
dokładnie to, co robi malarz formatów z rodziny `studio.format.painter.*`.
Kopiowanie schowka woła tę drogę i oddaje jej uchwyt, zamiast zakładać drugi
magazyn zabranej postaci; dwa magazyny rozjechałyby się przy pierwszym
naniesieniu.

Wklejenie zapisuje pochodzenie, ponieważ fragment wklejony ze schowka bywa
jedynym śladem, na czym pismo się opiera. Zapis pochodzenia rodzaju
`clipboard` jest podstawą pod wykaz podobieństw w panelu redaktora i pod
bibliografię — bez niego nie dałoby się odtworzyć, skąd wziął się akapit.

Postać przenoszona przy odłożeniu do schowka jest postacią początku zakresu,
nie postacią każdego znaku z osobna: zakres odłożony bywa niejednolity, a
zapisanie postaci znak po znaku dałoby wpis schowka rosnący razem z
długością treści i wklejenie wymagające dopasowania znak w znak przy każdej
różnicy długości. Kontrakt mówi o zachowaniu postaci źródła, nie o
przeniesieniu całego drzewa postaci, więc postacią zachowywaną jest postać
początku zakresu — dokładnie tak, jak działa malarz formatów. Różnica postaci
wewnątrz zakresu jest pominięciem świadomym, a bilans wklejenia o niej mówi.

Schowek nie jest drugim malarzem formatów, mimo tego samego rachunku postaci
skutecznej: malarz przenosi postać bez treści i ma na to własną tabelę z
wygasaniem, schowek przenosi treść i postać razem, bo tego wymaga kontrakt
sposobu `keepFormat`.

## budowa/server/internal/core/adapter_modul_design_bazy_zdjeciowe.go

Dostawcy dzielą się na dwie grupy, a kolejność zapytywania jest zamierzona.
Grupa bez klucza — Openverse, Wikimedia Commons, Met Museum, NASA — pracuje
u każdego użytkownika od pierwszej minuty, bez konta i bez konfiguracji,
dlatego wyszukanie bez wskazanego dostawcy pyta najpierw ją. Grupa z kluczem
darmowym — Unsplash, Pexels, Pixabay, Smithsonian — czyta klucz z sejfu
poświadczeń rdzenia; brak klucza nie jest odmową całej komendy, dostawca
wraca w polu odpowiedzi niosącym dostawców, którzy nie odpowiedzieli,
a pozostali oddają swoje wyniki. Smithsonian trafił do tej drugiej grupy,
choć jego klucz jest darmowy i wydawany natychmiast, ponieważ interfejs
api.si.edu/openaccess mimo to klucza wymaga.

Drogi wszystkich ośmiu dostawców są zmierzone próbą na serwerze próbnym: adres,
klucz, rozbiór odpowiedzi i złożenie wspólnej postaci zasobu. Kształty
odpowiedzi pochodzą w większości z dokumentacji dostawców, z jednym wyjątkiem
zmierzonym na żywym interfejsie: odpowiedź Smithsoniana została odczytana
wprost z api.si.edu i wykazała, że rdzeń pierwotnie szukał mediów pod kluczem
nazwanym content.descriptiveNonRepeating.media, którego odpowiedź nie niesie
— pole z mediami leży zagnieżdżone głębiej, pod kluczem online_media.

Cisza jest zakazana: odpowiedź niesie wykaz dostawców zapytanych i wykaz
dostawców, którzy nie odpowiedzieli. Wykaz zasobów krótszy, niż oczekiwano,
ma mieć obok siebie powód, inaczej czytelnik odpowiedzi uzna, że fraza nie
ma zdjęć, gdy w istocie część dostawców nie odpowiedziała.

Zasób bez zapisanej licencji jest usterką, nie zasobem: wciągnięcie zasobu
zapisuje licencję razem z nim i odmawia, jeśli zapis licencji się nie udał —
materiał, o którym nikt później nie powie, czy wolno go było użyć, jest
gorszy niż brak materiału.

## adapter_modul_studio_petla.go

Pętla wykonawcza i praca kilku wykonawców naraz są domyślnie wyłączone;
włącza je Operator jawnym, odwracalnym ustawieniem. Uruchomienie pętli przy
wyłączonej nastawie wraca odmową nazywającą brak nastawy — started: false
wraz z refusalReason mówiącym, czego brakuje i czym się to włącza — zamiast
ciszy albo udanym przebiegiem bez pracy: plan „uruchomiony”, który nie
wykonał ani jednego zadania, jest wzorcem szkody, którego produkt nie
powtarza. Rozkład zlecenia nastawy nie wymaga, bo ułożenie zadań niczego
nie uruchamia i nie tyka dokumentu.

Nastawy pętli mieszkają w rodzinie studio.agents.settings.* i idą zasięgami
rodziny config.*. Pętla ich nie zapisuje i nie trzyma własnej kopii — czyta
je wywołaniem obsługiwacza z rejestru komend, a nie własnym odczytem
konfiguracji, bo nazwy kluczy należą do tego, kto nastawy wystawia; brak
obsługiwacza nastaw znaczy dla pętli to samo, co nastawa wyłączona.

Rodzina studio.plan.* okna w żądaniu nie niesie: dokument Studia zna swoje
okno, więc pętla bierze je stąd, tą samą drogą, którą wzięłaby je operacja
kontekstowa.

Zadanie jedzie rejestrem, nie wywołaniem adaptera wprost. Wywołanie metody
adaptera wprost byłoby krótsze o dwie warstwy i ominęłoby oba owinięcia,
które stoją w rejestrze: podpis wykonawcy i zaporę blokad fragmentu. Pętla
omijająca zaporę byłaby cichą dziurą w blokadzie — Operator oznaczyłby
podstawę prawną jako nietykalną, a pętla przepisałaby ją, bo „to nie klient
wołał”. Model woła komendy tą samą drogą co klient, więc pętla też.

Tożsamość wykonawcy jedzie kontekstem, bo żądanie operacji kontekstowej pól
podpisu nie ma. Bez tego zapora wzięłaby pracę pętli za pracę Operatora,
a blokada jest skierowana przeciw wykonawcy, nie przeciw właścicielowi
dokumentu — i przepuściłaby ją.

Podpis wykonawcy w ładunku jest konieczny, nie ozdobny: zapora blokad
fragmentu rozpoznaje wykonawcę z dwóch źródeł, z faktu gniazda i z podpisu
w ładunku żądania. Pętlę puszcza Operator przyciskiem w oknie, więc gniazdo
mówi „Operator”, a pracę wykonuje model. Podpis jedzie mapą, nie strukturą
żądania, bo studio.contextual.op pól podpisu w kontrakcie nie ma — ładunek
jest wywołaniem wewnętrznym rdzenia, składanym po to, żeby przejść tą samą
drogą, którą idzie klient, wraz z jej owinięciami; obsługiwacz komendy pól
nadmiarowych nie widzi.

## adapter_modul_studio_wejscie_postac.go

Komenda `studio.document.save` przyjmowała wcześniej tylko `documentId`,
`content` i `title`, więc postać dokumentu ginęła przy każdym zapisie: model
widział tekst, nie widział kroju ani tabeli. Kontrakt dostał pole `form`, drogę,
którą postać dojeżdża do rdzenia i wraca z niego nieuszkodzona.

Wcześniejsza wersja tego odcinka nie miała jeszcze warstwy danych postaci
i odkładała postać ładunkiem JSON w tabeli katalogowej modułu, pod własnym
zasięgiem. Po wniesieniu warstwy danych postaci obie metody czytania i zapisu
zostały przełożone na `postacWczytaj` i `postacZapisz` obszaru postaci, żeby
dwa magazyny jednej postaci nie niosły dwóch prawd o tym samym dokumencie:
odcinek kontroli pracy, wołający `wejsciePostacDokumentu` przy zakładaniu
kopii zapasowej, czytając dawny magazyn dostawał postać pustą dla dokumentu,
który postać ma.

Gdy żądanie zapisu niesie i treść, i postać: treść z pola `content` jest
prawdą o literach, to, co użytkownik ma w edytorze; postać z pola `form` jest
prawdą o strukturze — arkusz stylów, nastawy strony, sekcje, tabele, obiekty
i aparat. Utrwalenie postaci składa treść z bloków i wpisuje ją do wiersza,
a wołający nadpisuje ją potem treścią z żądania, inaczej zapis cofałby litery
dopisane w oknie do stanu, który zna drzewo postaci.

Nazwa nośnika podana przez Operatora ma pierwszeństwo przy nastawach strony
i nie jest sprawdzana wykazem, ponieważ wykaz nośników stoi w module Design,
do którego ten odcinek nie sięga.

## adapter_modul_design_ikony_kroj.go

Krój ikonowy trzeba złożyć, nie tylko zapisać: trzeba zbudować tabele
TrueType. Biblioteka Go, która pisze krój od zera, w tym drzewie nie ma —
`tdewolff/font` krój czyta i przepisuje, ale złożenie go z konturów wymagałoby
zbudowania wszystkich tych samych tabel. Zapis stoi więc w rdzeniu, w jednym
pliku, i jest wkompilowany.

Ikony katalogu są rysowane kreską (`stroke`), a glif kroju jest obszarem
wypełnianym — kreska idzie więc przez `Path.Stroke` biblioteki
`tdewolff/canvas`, którego obrys staje się konturem glifu. Bez tego kroku
wszystkie glify wyszłyby jako cienkie zamknięte pętle albo jako plamy,
zależnie od kształtu.

Obrys kreski wychodzi z biblioteki jako łamana, bo `Path.Stroke` spłaszcza
krzywe do odcinków. Łamana zapisana wprost do tabeli `glyf` dawała plik kilka
razy większy niż potrzeba, dlatego jest tu dopasowywana z powrotem do krzywych
kwadratowych, tych samych, którymi TrueType opisuje glif
(`dopasujKrzyweKonturuDesignu`). Naroża są rozpoznawane i nigdy nie wchodzą
w środek krzywej: krzywa przeciągnięta przez naroże zaokrągliłaby kształt,
który ma zostać ostry.

`dopasujKrzyweKonturuDesignu` biegnie zachłannie: od każdego wierzchołka bieg
przedłuża się tak długo, jak cały jego ciąg mieści się w tolerancji — najpierw
jako odcinek, a gdy odcinek nie mieści się, jako krzywa. Bieg zachłanny nie
daje zapisu najkrótszego z możliwych, ale daje zapis kilka razy krótszy od
łamanej, w jednym przejściu po konturze. Naroża wyznacza się z góry i bieg
zawsze się na nich urywa — krzywa przeciągnięta przez naroże zaokrągliłaby
kształt po obu jego stronach mimo przejścia przez sam wierzchołek. Kontur
obraca się tak, żeby zaczynał się od naroża, inaczej naroże wypadające na
styku końca i początku konturu zostałoby zaokrąglone.

Ikony nie są literami, więc nie mają własnych punktów kodowych Unicode. Krój
przypisuje im kolejne punkty obszaru prywatnego od U+E000 — tak robią
wszystkie kroje ikonowe i tego oczekuje arkusz stylów, który je odczyta.

Ikona, z której nie da się wyciągnąć ani jednej ścieżki, jest odmową całego
pakietu składania kroju: krój z pustym glifem w środku wygląda jak krój
gotowy, a w miejscu tej ikony pokazuje nic.

## adapter_modul_studio_tabele.go

Rachunek siatki tabeli stoi w pliku pomocniczym, a droga zapisu, zmiany
śledzonej i dziennika w obszarze postaci — ten plik nie liczy żadnej z tych
rzeczy drugi raz.

Tabela wchodzi do drzewa dokumentu blokiem nietekstowym, więc jej wstawienie
nie przesuwa ani jednego zakresu zaznaczenia, przypisu ani blokady. Gdyby
tabela zajmowała znaki treści, każde jej wstawienie rozjeżdżałoby wszystko,
co wisi na miejscu w treści.

Ustawianie postaci tabeli i komórek rozstrzyga zasięg wskazania tak samo, jak
w pakiecie biurowym: brak wskazania wiersza i kolumny znaczy całą tabelę,
wskazanie samego wiersza — cały wiersz, samej kolumny — całą kolumnę.

## budowa/server/internal/core/adapter_modul_design_wektor_sciezki.go

Czynności kontraktu rodziny wektorowej stoją w
`adapter_modul_design_wektor.go`; ten plik niesie sam rachunek na ścieżkach,
żeby miał jedno miejsce i jedną prawdę.

Rachunek jest wkompilowany, nie wołany: operacje logiczne na ścieżkach,
wydanie PDF i wydanie EPS idą przez bibliotekę Go wkompilowaną w binarium
serwera, bez uruchamiania procesu zewnętrznego. Programy do obrysowywania
konturów, rysowania wektorowego i rasteryzacji języka strony leżą poza
instalacją operatora, więc funkcja od nich zależna byłaby u niego odmową,
nie funkcją.

Węzeł jest bytem produktu, ścieżka biblioteki jest tylko rachunkiem: kontrakt
niesie węzły z uchwytami i to one leżą w bazie. Ścieżka biblioteki powstaje na
czas rachunku i ginie po nim; wynik wraca znowu jako węzły, żeby operator
mógł go dalej ciągnąć piórem. Zapisanie wyniku jako gotowego napisu SVG
odebrałoby mu edycję — kształt przestałby mieć węzły, a zostałby obrazkiem.

Uchwyt jest odsunięciem od węzła, nie punktem bezwzględnym — tak opisuje pole
kontraktu. Punkt sterujący krzywej jest więc węzłem plus odsunięciem; odczyt
odsunięcia jako współrzędnej bezwzględnej przesuwałby krzywe ku początkowi
układu przy każdym przejściu przez bazę.

Odcinek między dwoma węzłami ścieżki biblioteki jest krzywą sześcienną, gdy
którykolwiek węzeł niesie uchwyt po tej stronie odcinka; brak uchwytu z
jednej strony bierze punkt sterujący na samym węźle, tak jak działa pióro w
programie wektorowym — węzeł narożny z jednej strony i wygładzony z drugiej
daje krzywą wchodzącą z jednej strony prosto.

Przy rozkładzie ścieżki biblioteki z powrotem na węzły łuki wchodzą jako
krzywe, bo kontrakt nie ma węzła łukowego, a łuk zamilczany zgubiłby kawałek
kształtu. Wielościeżkowy wynik operacji logicznej daje jeden wykaz węzłów,
ponieważ kontrakt niesie jedną ścieżkę wynikową, więc rozdzielone kawałki
idą po sobie, a nie giną.

Ostatni węzeł zbieżny z pierwszym przy ścieżce zamkniętej jest zbędny, bo
domknięcie samo prowadzi z ostatniego do pierwszego; uchwyt wchodzący
zostaje jednak przeniesiony na pierwszy węzeł, ponieważ opisuje krzywiznę
domknięcia.

Czyszczenie precyzji węzłów nie usuwa węzłów: węzeł, który operator
postawił, jest jego rozstrzygnięciem o kształcie, a ubytek bajtów bierze się
z krótszego zapisu liczb, nie z gubienia jego pracy.

Kształt podstawowy powstaje od razu jako węzły ścieżki, a nie jako osobny
byt do późniejszej zamiany: prostokąt dorysowany na kanwie ma dać się
natychmiast ciągnąć piórem za narożnik, bez komendy zamiany na ścieżkę,
której kontrakt nie ma i mieć nie będzie.

## budowa/server/internal/core/adapter_modul_studio_dziennik.go

Wpis dziennika niesie całe drzewo postaci, nie samą treść, ponieważ zmiana
formatowania obejmująca cały dokument, na przykład kroju, nie rusza ani jednej
litery treści — treść sprzed czynności więc jej nie odtworzy. Wpis dziennika
niesie stan przed i stan po jako pełne drzewo postaci, a cofnięcie przywraca
drzewo, nie napis.

Cofnięcie czynności nie jest przywróceniem wersji. Przywrócenie wersji cofa
wszystko, co po niej weszło, także pracę wykonaną później. Cofnięcie czynności
ze środka dziennika ma zdjąć wyłącznie to, co zrobiła ta jedna czynność:
dlatego rdzeń nie wgrywa stanu przed w miejsce dokumentu, tylko liczy różnicę
między stanem po i stanem przed i nakłada ją odwrotnie na stan bieżący, byt
po bycie — blok, styl, sekcja, obiekt, pole, nastawy strony. Byt, którego
czynność nie tknęła, zostaje nietknięty.

Czynność, na której stoi zależność, odmawia cofnięcia, nie cofa się po cichu.
Tabela zależności czynności istnieje wyłącznie po to: cofnięcie czynności, na
której stoi późniejsza, zostawiłoby dokument w stanie niespójnym, na przykład
zdjęcie tabeli, do której później wstawiono wiersze. Odmowa nazywa zależność,
zamiast zostawiać pytanie, dlaczego nic się nie stało.

Zależność liczy się z dwóch źródeł, nie z jednego. Zależność zapisana mówi
o powiązaniach, których z zakresu nie widać — wstawienie tabeli jest podstawą
scalenia komórki, choć zakresy mogą się nie stykać; tę wiedzę ma czynność,
która wpis odkładała, i ona ją zapisuje. Zależność wywiedziona z zakresu jest
siatką pod tym: czynność późniejsza, która ruszyła ten sam fragment, stoi na
tej wcześniejszej z samej natury rzeczy — cofnięcie wcześniejszej wgrałoby
w to miejsce stan sprzed niej i zabrałoby ze sobą pracę późniejszą. Oba źródła
liczą się naraz, ponieważ wykaz zależności zapisanych jest dziś niepełny,
a niepełny wykaz sam w sobie znaczyłoby ciche cofnięcie psujące dokument;
zależność wywiedziona z zakresu myli się w drugą stronę, odmawiając czasem
cofnięcia bezpiecznego, ale mówi wtedy wprost, która czynność stoi na drodze
— odmowa nazwana jest tańsza niż dokument niespójny.

Rodzaj czynności dziennika i rodzaj zmiany śledzonej są dwoma osobnymi
słownikami kontraktu, których nie wolno pomieszać. Stan wpisu dziennika jest
wyłącznie aktywny albo cofnięty — inna wartość nie wchodzi do tabeli dziennika.

Aparat dokumentu wraca przy przywracaniu bytów tak samo jak sekcja i obiekt:
przekład bytu kontraktu na wiersz istniał w obszarze aparatu, lecz w dzienniku
pozostawał nieużyty, więc cofnięcie czynności odtwarzało treść i postać, a spis
treści, przypis albo bibliografia zostawały w brzmieniu bieżącym. Kotwica
i znacznik nieświeżości elementu aparatu idą z samego elementu, nie
z przeliczenia od nowa: cofnięcie ma przywrócić stan sprzed czynności, a nie
policzyć aparat na nowo z treści bieżącej, co wstawiłoby numery stron, których
w stanie sprzed czynności nie było.

## adapter_modul_studio_strona.go

Każda czynność tego pliku idzie jedną drogą obszaru postaci: wczytanie
postaci, zmiana postaci, zakończenie postaci. Dzięki temu zmiana śledzona
autora model i wpis odwracalnego dziennika odkładają się same, tą samą
drogą co zmiana kroju czy wcięcia — model, który przestawił nośnik
załącznika na poziomą A3, jest widoczny w podświetleniu zmian modelu tak
samo jak model, który dopisał akapit. Rodzaj czynności dziennika jest dla
całego pliku jeden, bo Operator cofa „zmianę strony”, a nie „formatowanie”;
rodzaj zmiany śledzonej idzie osobnym polem, bo wstawienie podziału jest
wstawieniem, a wszystko inne formatowaniem.

Sekcja bez własnych nastaw dziedziczy nastawy dokumentu — jej pole nastaw
zostaje puste. Skopiowanie nastaw dokumentu do każdej sekcji przy
zakładaniu dałoby dokument, w którym zmiana nośnika „dla całości” nie
rusza ani jednej strony, bo każda sekcja trzymałaby własną kopię sprzed
zmiany. Dlatego nastawa skuteczna liczy się przy odczycie, nakładając
nastawy sekcji na nastawy dokumentu.

Koperta bez położenia adresata i nadawcy byłaby nadrukiem w punkcie zero,
czyli na samej krawędzi, więc położenie niepodane dostaje nastawę z normy:
adresat w prawej dolnej ćwiartce, nadawca w lewym górnym narożniku.

Podział jest blokiem nietekstowym: nie zajmuje ani jednego znaku treści,
więc jego wstawienie nie przesuwa żadnego zaznaczenia ani zakotwiczenia.
Podział wypadający w środku akapitu rozdziela ten akapit na dwa, a dopiero
to rozdzielenie dokłada jeden znak podziału wiersza, o który przesuwa się
wszystko, co za nim wisi; ten jeden znak trzeba potem uwzględnić
w zakotwiczeniach, inaczej przypis za podziałem wskazywałby o literę
wcześniej. Podział sekcji zakłada nową sekcję od miejsca podziału do końca
dokumentu i skraca sekcję poprzednią — bez tego „podział sekcji” byłby
kreską w treści, za którą nic się nie zmienia, a sensem sekcji są jej
własne nastawy.

Postać akapitu przy rozdzieleniu jedzie kopią: dwa akapity na jednej
strukturze znaczyłyby, że wcięcie jednego zmienia drugi.

Dokument bez ani jednej sekcji zapisanej oddaje sekcję jedną, obejmującą
całą treść, bo taki dokument jest jedną sekcją, a pusty wykaz kazałby
oknu twierdzić, że dokument nie ma strony; ta sekcja nie jest przy tym
zakładana w bazie, bo odczyt nie ma prawa zapisywać.

## adapter_modul_studio_autozapis.go

Wersje nazwane i kluczowe zakłada Operator, i to jest historia jego decyzji.
Gdyby zapisy samoczynne wchodziły do tego samego wykazu, po godzinie pracy
historia przestałaby być historią decyzji i stałaby się dziennikiem naciśnięć
klawisza, w którym Operator nie znalazłby własnej wersji nazwanej. Rozróżnienie
idzie kolumną szeregu, a pojęcie wersji kluczowej nie zakłada się drugi raz:
Studio ma je w `studio.version.label.set`.

Kopia zapasowa jest osobna od wersji, bo ma przetrwać awarię procesu i awarię
zapisu. Wersja leży w repozytorium i zakłada się ją tym samym zapisem, który
właśnie się nie udał, więc wersja nie ochroni pracy przed nieudanym zapisem —
kopia zakłada się niezależnie, przed próbą zapisu właściwego.

Wskaźnik „zapisano” pokazany, gdy zapis się nie udał, jest najgorszym możliwym
błędem tego modułu: Operator zamknie okno i straci pracę. Dlatego nieudany
zapis samoczynny nie wraca odmową, która przepadłaby w dzienniku zdarzeń —
wraca odpowiedzią niosącą `saved: false`, nazwany powód i kopię, w której
praca została. Wiersz nastaw zapamiętuje ten skutek, więc kolejne
`autosave.get` powie prawdę o ostatnim zapisie także wtedy, gdy okno tymczasem
się przeładowało.

Gdy zapis się udaje, kopia zapasowa przestaje nieść pracę, której w dokumencie
już nie ma. Wiersz kopii zostaje w wykazie, ale zakłada się nowy wiersz
z wyzerowanym wskaźnikiem niezapisanych zmian, zamiast nadpisać istniejący,
żeby nie zatrzeć śladu, że kopia była założona przed zapisem.

Przywrócenie kopii do nowego dokumentu jest osobną drogą, bo przywracanie samo
nie ma kasować tego, co jest: Operator, który po awarii nie jest pewien, która
wersja jest lepsza, ma dostać obie.

## budowa/server/internal/core/adapter_modul_design_kolor.go

Czynności czytające piksele leżą w `adapter_modul_design_kolor_obraz.go`; sam
rachunek barwy leży w `adapter_modul_design_barwy.go`. Ten plik niesie
czynności kontraktu, które liczą bez dotykania obrazu.

Harmonia jest rachunkiem na kole barw, nie tablicą gotowych palet: dopełnienie,
triada, tetrada i analogia to obroty odcienia o ustalony kąt. Rachunek idzie w
przestrzeni HCL biblioteki `go-colorful`, nie w HSL, ponieważ HCL jest
percepcyjnie równomierna — dwie barwy o tym samym odcieniu i jasności
wyglądają na równie jasne, a nie tylko mają równe liczby. Paleta zbudowana w
HSL dałaby żółć jaskrawą i błękit przygaszony przy identycznych nastawach.

Zapis barwy, którego rdzeń nie zrozumiał, jest odmową, nigdy podstawieniem
czerni domyślnej: paleta zbudowana wokół czerni podstawionej za nieczytelne
wejście wyglądałaby w oknie jak paleta poprawna.

Badanie dostępności zestawu żetonów mierzy pary, nie pojedyncze żetony: bierze
wszystkie pary barwnych żetonów zestawu i mierzy kontrast każdej. Odpowiedź
niesie liczbę par sprawdzonych, liczbę spełniających próg i wykaz łamiących
go — bilans, nie samo stwierdzenie, że są problemy.

## adapter_modul_studio_znakowanie.go

Trzy byty tego pliku nie wolno zlać, bo rdzeń je odróżnia w danych tak, jak
mają być odróżnione w oknie: wyróżnienie (`highlight`) jest cechą postaci
dokumentu — barwa tła fragmentu idzie tą samą drogą co reszta formatowania,
czyli w runy drzewa postaci, a wiersz znakowania jest tylko wykazem do
przejścia, nie drugim miejscem, w którym trzymana jest barwa. Znacznik
(`mark`) jest nazwą własną Operatora i treści dokumentu nie rusza. Propozycja
(`suggestion`) jest brzmieniem proponowanym, stojącym na marginesie: nie jest
zmianą śledzoną, bo tamta jest już w treści, i nie jest komentarzem, bo tamten
nie niesie brzmienia — wchodzi do treści dopiero decyzją Operatora, a wtedy
odkłada się jako zmiana śledzona jej autora.

Blokadę fragmentu zdejmuje wyłącznie Operator. Wykonawca, który uzna, że
fragment wymaga zmiany, zakłada propozycję zamiast wnosić zmianę wprost —
dlatego założenie propozycji na fragmencie zablokowanym jest dozwolone,
a wniesienie jej do treści już nie: decyzja należy do Operatora, a blokada
o zasięgu modelu nie wiąże jego decyzji.

Przy rozstrzyganiu propozycji zakresy liczone są w treści sprzed decyzji,
więc rozstrzyganie idzie od końca dokumentu — przyjęcie od początku
przesuwałoby zakresy propozycji jeszcze nierozpatrzonych.

## adapter_modul_studio_szablony_pism.go

Szablon niesie naraz cztery rzeczy — arkusz stylów, nastawy strony, nagłówek
i stopkę oraz pola do wypełnienia — nie samą treść. Zapis szablonu z dokumentu
bierze całą postać dokumentu, nie tylko jego treść; wniesienie z pliku
przejmuje wszystkie te rzeczy, bo czyta plik tym samym rachunkiem, którym
czyta się dokument; dokument zakładany z szablonu dostaje tę postać
podstawioną, a nie domyślną platformy; blokady wzorcowe idą razem, bo fragment
wzorcowy pisma ma zostać wzorcowy.

Szablon fabryczny nie jest do usunięcia: odmowa nazywa powód wprost, wzorem
`studio.operation.delete`, który robi to samo dla operacji fabrycznych. Cicha
bezczynność albo `deleted: false` bez zdania byłyby tu odpowiedzią, po której
nie dałoby się poznać, że szablon zszedł.

Pola szablonu leżą w `szablon_studio.pola_json`, jednym miejscem — tak stanowi
migracja 367 i tak je czyta `studio.template.apply`. Własna tabela pól dałaby
dwa miejsca prawdy: `apply` czytałby jedno, warsztat drugie, a rozjazd
wyszedłby przy pierwszym szablonie założonym starą drogą.

Wypełnianie szablonu (`WypelnijSzablon`) jest czynnością dostępną także
modelowi AI — dlatego autor czynności wchodzi do dziennika, a czynność modelu
odkłada się jako jego, tym samym rachunkiem autorstwa co czynność człowieka.

## adapter_modul_design_wykresy.go

Wykres bierze serie danych, nie obraz: kontrakt żąda serii liczb, żeby wykres
dawał się przerysować po zmianie liczb. Gdyby komenda przyjmowała obraz,
zmiana jednej wartości wymagałaby narysowania wykresu od nowa poza produktem,
a rdzeń byłby tylko miejscem, w którym ten obraz leży.

Wykres i schemat powstają raz, jako płótno biblioteki `tdewolff/canvas`,
a potem wychodzą albo jako SVG (wydawca `renderers/svg`), albo jako PNG
(rasteryzator `renderers/rasterizer`). Dwie osobne drogi rysowania dałyby dwa
wykresy różniące się szczegółami zależnie od formatu, mimo że mają być tym
samym wykresem.

Podpisy są konturami, nie elementem tekstowym: napisy jadą przez
`sciezkaTekstuDesignu`, jako kontury glifów kroju wkompilowanego. Element
`<text>` w SVG pokazałby podpisy wyłącznie tam, gdzie ten krój jest
zainstalowany, a wykres ma być plikiem przenośnym.

`design.diagram.render` niesie `unplacedNodeIds` — węzły, których układ nie
umieścił. Schemat z połową węzłów wygląda bez tego pola jak schemat kompletny,
dlatego bilans nazywa brak wprost zamiast milczeć o nim.

## adapter_modul_design_fotografia_rachunek.go

Podstawa każdej czynności liczy się w procesie, czystym Go:
`disintegration/imaging` (Lanczos, kadr, obrót, rozmycie, wyostrzenie,
korekcje barwne), `golang.org/x/image/draw` (przekształcenia afiniczne,
kompozycja warstw), `golang.org/x/image` (WEBP, TIFF) oraz rachunek własny na
maski, progowanie, filtr bilateralny odszumiania i obrysowanie konturów. Nie
ma tu ani jednego uruchomienia zewnętrznego procesu i mieć nie będzie: silniki
obrazu i narzędzia obrysowywania konturów, po które sięga się przy takiej
pracy, leżą poza instalką docelową, więc funkcja od nich zależna byłaby
odmową, nie funkcją. Zapora `zapora_fotografii_test.go` pilnuje tego
maszynowo.

Każda funkcja tego pliku oddaje obraz. Żadna nie oddaje „powodzenia” bez
obrazu: droga bez pikseli jest błędem nazywającym brak, a nie odpowiedzią
statusu poprawnego z pustym wynikiem.

Odszumianie bilateralne (`odszumBilateralnieDesignu`) różni się od rozmycia
tym, że waży sąsiada nie tylko odległością, ale i różnicą jasności: rozmycie
waży wyłącznie odległością, więc na granicy dwóch obszarów miesza jeden
z drugim i krawędź traci ostrość, a filtr bilateralny dokłada wagę zakresową
— sąsiad po drugiej stronie krawędzi różni się jasnością i wchodzi do
średniej z wagą znikomą, więc szum wewnątrz obszaru znika, a granica zostaje.
Biblioteka wkompilowana filtru bilateralnego nie ma, więc rachunek jest
własny, złożony z dwóch przejść — poziomego i pionowego — zamiast jednego
jądra kwadratowego: jądro bilateralne nie jest rozdzielne dokładnie, to
przybliżenie, ale krawędzie zachowuje w obu kierunkach przy koszcie rzędu
rozmycia rozdzielnego, dwa razy 2r+1 odczytów na punkt zamiast (2r+1)².
Rachunek pełny na obrazie stumilionowym byłby stu sześćdziesięcioma
dziewięcioma odczytami na punkt, a odszumienie przestałoby być czynnością
wykonywalną w rozsądnym czasie. Kanał krycia przechodzi nietknięty:
odszumienie zmienia barwę, nie przezroczystość, a uśrednienie krycia rozmyłoby
wycinek odcięty od tła.

## adapter_modul_design_ikony.go

Ikona własna i wzór rdzenia to dwa byty, nie jeden: wzory rdzenia są
wkompilowane i niezmienne, ikony własne leżą w bazie i mają okno.
Identyfikator rozstrzyga, o który byt chodzi — wzór niesie przedrostek
zestawu (`rdzen-24/dom`), ikona własna identyfikator zewnętrzny wiersza.
Zmiana wzoru rdzenia jest niemożliwa i tak ma być: gdyby dała się zapisać,
dwie instalacje produktu miałyby dwa różne katalogi pod tą samą nazwą.

`design.icon.set` niesie `gridWarnings` — miejsca, w których ikona nie trzyma
siatki (współrzędne poza polem, grubość obrysu inna niż zadeklarowana, brak
pola widoku). `design.icon.generate` niesie `failedConcepts` — pojęcia, dla
których ikona nie powstała. Zestaw, w którym połowa ikon nie weszła, wygląda
bez tych pól jak zestaw kompletny.

`design.font.preview` mówi polem `available`, czy rdzeń krój ma. Podgląd
złożony krojem zastępczym wygląda identycznie jak prawdziwy, a wybrana
wtedy typografia nie jest tą, którą zobaczy się u siebie.

## adapter_modul_design_uchwyty.go

Komenda `design.asset.tag.set` domyka lukę odczytu: `design.asset.list` zawęża
wykaz zasobów polem `tags`, a bez tej komendy nie byłoby czym etykiet nadać.

Zdarzenie `design.asset.created` ma dwóch nadawców: generowanie i wniesienie
zasobu. Obie drogi odkładają bajty w magazynie rdzenia pod sumą kontrolną,
zanim powstanie wiersz zasobu, a odmowa nie rozgłasza niczego, ponieważ
opakowanie milczy przy błędzie — droga bez bajtów jest więc też drogą bez
zdarzenia.

`UsunZasob` oddaje trzy wartości jako jedyna metoda portu Design: odpowiedź
kontraktu niesie samo pole `removed`, a zdarzenie `design.asset.deleted` musi
nieść cały usunięty zasób, bo po usunięciu wiersza z bazy nie ma go już skąd
odczytać — adapter podaje go obok odpowiedzi. Rozgłoszenie z wnętrza adaptera
byłoby drugą drogą do szyny zdarzeń, skoro emiter należy do rejestru, a nie do
modułu.

Metody odczytu kompozycji dla szyny zdarzeń oddają brak wiersza jako wartość
`false`, a nie jako błąd: zmiana kompozycji już się udała w chwili odczytu,
więc niepowodzenie odczytu ma zamknąć usta szynie zdarzeń, a nie unieważnić
wykonaną komendę.

## budowa/server/internal/mowa/dziennik.go

Plik powiela kształt repozytoriów pakietu `dane` (blok stałych ze składanym SQL,
typ wiersza, filtr, interfejs repozytorium, asercja `var _`, odczyt wiersza przez
interfejs skanera), ale stoi w pakiecie `mowa` i bierze `*sql.DB` wprost. Stąd
własny interfejs skanera: odpowiednik `dane.skaner` jest nieeksportowany i nie
da się go użyć spoza pakietu `dane`.

Dziennik nie ma własnego zegara — czas przychodzi z góry w polu `Utworzono`
wpisu, a baza go nie wstawia. Jeden zegar na jedno zdarzenie: drugi
rozjeżdżałby ślad z chwilą zapisu.

Dziennik zapisuje także odmowy: `Zapisz` przyjmuje wpis o stanie `odmowa`
i jest to jego zwykłe użycie, bo powód odmowy jest pierwszą informacją
potrzebną, gdy nic się nie przepisało. Więz spójności stanu z powodem pilnuje
CHECK tabeli `transkrypcja`, a sprawdzenia w pliku odrzucają niespójny wpis
wcześniej, żeby wołający dostał zdanie po polsku zamiast komunikatu sterownika
bazy danych.

Zapytanie listaTranskrypcji obsługuje oba warianty żądania: puste zawężenie
okna wyłącza pierwszy warunek. Ze wskazanym oknem porządek biegnie indeksem
idx_transkrypcja_wykaz (okno_id, utworzono, id) czytanym wstecz, więc
„najnowsze najpierw” nie kosztuje sortowania wyniku (EXPLAIN QUERY PLAN:
SEARCH ... USING COVERING INDEX, bez kroku ORDER BY). Bez wskazania okna
sortowanie zostaje — to zapytanie diagnostyczne przez cały dziennik, a nie
widok otwierany przy każdym zleceniu. Kolumna `id` w porządku rozstrzyga
wpisy z tej samej milisekundy, żeby kolejność była stała między wywołaniami.

Funkcja sprawdzWpis odrzuca wpisy, których baza i tak by nie przyjęła — po to,
żeby wołający dostał zdanie po polsku zamiast komunikatu o naruszonym CHECK-u.
Zdublowaniem więzu to nie jest: baza pozostaje ostatecznym strażnikiem, bo
pisać do niej może też przyszła ścieżka, która tej funkcji nie wywoła.

Funkcja napisDoKolumny przekłada napis pusty na NULL. Kolumny `okno_id`
i `powod` dopuszczają pustkę, a pustka ta coś znaczy — zlecenie spoza okna
oraz brak odmowy. Napis pusty zapisany wprost udawałby wartość podaną,
a przy `powod` naruszałby CHECK tabeli `transkrypcja`.

## adapter_modul_design_generowanie.go

Wiersz zasobu powstaje wyłącznie po utrwaleniu bajtów w magazynie: ciąg biegnie
od złożenia polecenia, przez wywołanie kanału obrazowego, odłożenie bajtów pod
sumą sha256, pomiar formatu i wymiarów z nagłówka utrwalonego pliku, aż po
założenie wiersza z uri, format, width, height. Kafelek w panelu zasobów, za
którym nic nie leży, byłby kłamstwem koperty.

Każdy brak kończy komendę odmową nazywającą brak, nigdy obrazem zastępczym:
brak kanału obrazowego, kanał nieczynny, kanał tekstowy zamiast obrazowego,
brak poświadczenia, odpowiedź bez obrazu, bajty nie do pobrania. Nie ma drogi,
którą wracałby status ok bez treści. Wykaz braków w odmowie niesie tylko braki
danego wywołania, obok gotowej treści polecenia.

Adres zamiast bajtów też ląduje w magazynie. Dostawcy zgodni z OpenAI Images
oddają b64_json albo url, a url bywa domyślny; odsyłacz dostawcy wygasa, więc
zapisanie go wprost jako uri zasobu dałoby zasób, który po godzinie przestaje
mieć treść. Bajty spod adresu wciąga rdzeń od razu i dopiero one idą do
magazynu; niepowodzenie pobrania jest odmową, nie zasobem bez treści.

Wariant nieudany przerywa całość generowania. Gdyby drugi wariant padł, a
pierwszy został, odpowiedź niosłaby mniej zasobów, niż zamówiono, bez słowa o
tym, czemu — kontrakt nie ma pola na wynik częściowy. Zasoby, które zdążyły
powstać, zostają w bazie i w magazynie, bo ich bajty są prawdziwe i kasowanie
ich byłoby niszczeniem cudzej treści z powodu, który jej nie dotyczy.

Prompt utrwala się raz na całe wywołanie, przed pierwszym wariantem, żeby
wszystkie warianty miały jeden wspólny wiersz promptu w historii — kontrakt
history.list czyta prompty okna i kanał, którym poszły. Niepowodzenie zapisu
promptu nie przerywa generowania, bo prowenancja jest wiedzą o zasobie, a nie
samym zasobem.

## budowa/server/internal/core/adapter_modul_studio_podglad.go

Plik obsługuje komendy `studio.preview.render` oraz `studio.diff.visual`.
Silnik wyrysu stoi w `adapter_studio_wyrys.go`; ten plik odpowiada wyłącznie
za to, skąd wziąć treść, jakimi nastawami ją wyrysować i gdzie odłożyć wynik.

Podgląd oddaje strony, nie tekst. Kontrakt obu komend mówi o zasobach:
`pageAssetIds` i `overlayAssetIds`. Podgląd ma pokazać układ — typografię,
paginację, nagłówek i stopkę — a nie ten sam tekst, który operator widzi
w edytorze. Strona wychodzi więc obrazem, bo obraz jest jedyną postacią,
w której układ da się zobaczyć bez drugiego silnika składu po stronie okna.

Format docelowy komendy WyrenderujPodglad rozstrzyga o postaci wyniku
dwustopniowo: strony powstają zawsze jako obrazy, a przy formacie `pdf`
z tych samych stron składa się dodatkowo dokument. Dwa różne silniki —
jeden dla podglądu, drugi dla wydania — dawałyby dwa różne układy, a wtedy
podgląd przestaje być podglądem.

Komenda PorownajWizualnie idzie po wyrysie obu wersji, nie po ich tekście:
sedno tej czynności jest w tym, żeby zobaczyć zmianę, której różnica
tekstowa nie widzi — przesunięcie akapitu na następną stronę, zmianę
łamania, przestawienie nagłówka. Obie strony rysuje ten sam silnik tymi
samymi nastawami, więc różnica pikseli jest różnicą treści, a nie różnicą
sposobu rysowania.

Funkcja geometriaDokumentuStudia oddaje geometrię wyrysu wyliczoną z nastaw
strony dokumentu, a nie ze stałej A4. Postać dokumentu nieczytelna albo
niezapisana daje geometrię domyślną, a nie odmowę: podgląd dokumentu, który
nastaw strony jeszcze nie ma, jest normalną drogą, a nie usterką. Nastawy
sekcji ta droga bierze przez nastawy dokumentu, na które sekcja pierwsza
się nakłada. Dokument o sekcjach różnych nośników wychodzi w wyrysie jednym
rozmiarem — wyrys składa jeden ciąg kartek i drugiego rozmiaru w tym samym
ciągu nie umie. Rachunek stron aparatu liczy za to sekcja po sekcji, bo
numer strony musi być prawdziwy nawet wtedy, gdy obraz kartki jest
przybliżeniem.

## adapter_modul_studio_galezie.go

Rozgałęzienie zakłada wiersz gałęzi i jedną wersję startową na niej, a nie
drugi dokument. Drugi dokument oderwałby wariant od historii, z której wyrósł:
nie dałoby się powiedzieć, od czego wariant odszedł, a scalanie nie miałoby
wspólnego przodka, na którym stoi cała rodzina operacji gałęzi.

Scalanie trójstronne rozstrzyga samo tam, gdzie zmieniła jedna strona. Tam,
gdzie zmieniły obie i zmieniły inaczej, rdzeń nie wybiera — wybór strony
docelowej byłby cichym skasowaniem cudzej redakcji, dlatego konflikt wraca
kontraktem i czeka na rozstrzygnięcie w interfejsie.

Wspólnym przodkiem scalania jest wersja startowa gałęzi scalanej: od niej
wariant odszedł, więc to ona mówi, co w każdej z dwóch treści jest zmianą,
a co stanem zastanym. Bez przodka porównanie dwóch czół dałoby konflikt na
każdym fragmencie, który zmieniła tylko jedna strona.

Odwołanie do wersji jest adresem w obrębie platformy, nie odnośnikiem
sieciowym: rdzeń nie wystawia treści dokumentu na zewnątrz i nie ma jak
zapewnić, że adres wyprowadzony na świat byłby czytelny wyłącznie dla
uprawnionych.

Porównanie bloków zmiany opiera się na najdłuższym wspólnym podciągu wierszy,
a nie na przycinaniu wspólnego przedrostka i sufiksu, ponieważ przycinanie
dałoby jeden wielki blok na całą treść — wtedy każde scalenie dwóch redakcji
tego samego dokumentu kończyłoby się konfliktem. Przycinanie zostaje jedynie
drogą zapasową dla treści zbyt długich na tablicę podobieństwa, której rozmiar
rośnie iloczynem długości porównywanych wierszy.

## budowa/server/internal/core/adapter_modul_studio_adnotacje.go

Decyzja zmienia treść dokumentu, a nie tylko znacznik. Zmiana śledzona
i propozycja zmiany są bytami tymczasowymi: istnieją po to, żeby ktoś je
przyjął albo odrzucił. Decyzja, która przestawia sam stan wiersza i zostawia
treść dokumentu nietkniętą, byłaby decyzją bez skutku — Operator zobaczyłby
„przyjęto", a w dokumencie dalej stałby tekst sprzed zmiany. Dlatego obie
drogi kończą się zapisem treści i założeniem wersji.

Zmiany stosuje się od końca. Każda zmiana niesie zakres znaków liczony
w treści sprzed decyzji. Zastosowana od początku przesuwałaby zakresy zmian
jeszcze nierozpatrzonych o różnicę długości, więc druga zmiana trafiłaby
w niewłaściwe miejsce. Od końca — przesunięcie dotyczy wyłącznie tego, co już
rozpatrzono.

## budowa/server/internal/core/adapter_modul_design_zetony.go

Wydanie zestawu do kodu leży w `adapter_modul_design_zetony_wydanie.go`. Metody
pliku stoją na `*adapterDesignu` (`adapter_modul_design.go`).

Rdzeń nie nadpisuje motywu produktu: motyw jest własnością powłoki, a zestaw
żetonów jest bytem obok niego. Operator go zakłada, wczytuje z zapisu
zewnętrznego i wydaje do kodu, a produkt dalej wygląda tak, jak wygląda. Import
oddaje przy okazji bilans ról, których system produktu nie zna — bilans zamiast
ciszy.

Rola nieznana wchodzi do zestawu i jest wymieniona: odrzucenie takiej roli
byłoby zgubieniem pracy Operatora, który wczytuje system projektowy klienta,
a ten ma własne nazwy. Przemilczenie jej byłoby obietnicą, że wszystko pasuje.
Rola więc wchodzi, a jej nazwa wraca w `unknownNames`.

Autor komentarza bierze się z żądania i z faktu gniazda uruchomieniowego, nie
z zaszytej wartości. Model zakłada komentarz podpisany jako `model` tą drogą
i tylko tą; autor zaszyty jako Operator kazałby przełącznikowi pokazującemu
wszystko, co zrobił model, przemilczeć każdy komentarz modelu.
## server/internal/core/adapter_modul_studio_wersje.go

PrzywrocWersje: w przeciwienstwie do ZapiszDokument z createVersion=true,
adapter nie dopisuje tu drugiego zapisu, bo wersja docelowa jest podana
wprost.

Wykaz w `roleSystemuWizualnegoDesignu` jest kopią wykazu z panelu żetonów po
stronie klienta i jest to cena świadoma: kontrakt nie niesie bytu „rola
systemu wizualnego", więc wspólnego źródła dla obu stron nie ma. Skutek
rozjazdu jest ograniczony z zamysłu: wykaz służy wyłącznie do wypełnienia
`unknownNames` przy imporcie, a rola spoza niego i tak wchodzi do zestawu.
Rozjazd daje więc bilans zbyt obszerny, nigdy zgubiony żeton.

Rodzaj żetonu spoza kontraktu w `ZapiszZestawZetonow` jest odmową wołającemu,
nie odbiciem od schematu: tabela `zeton_design` świadomie nie ma warunku CHECK
(nagłówek migracji 235), więc sprawdzenie stoi w kodzie i wymienia rodzaje
dopuszczalne.

Postać zapisu w `WczytajZestawZetonow` rozpoznaje rdzeń, gdy wołający jej nie
wskazał. Rozpoznanie idzie po treści, nie po nazwie: zapis zaczynający się
nawiasem klamrowym jest JSON-em niezależnie od tego, jak nazywał się plik
u Operatora.

Rdzeń w `rozlozZapisZetonowDesignu` czyta zapis JSON wzorem W3C Design Tokens.
Postać spoza JSON-a i zmiennych CSS jest odmową wymieniającą postacie czytane
— zgadywanie dałoby zestaw złożony z przypadkowych napisów.

W `rodzajZetonuZWartosciDesignu` rodzaj wychodzący jako miara przy niepewności
jest najczęstszym rodzajem w systemach projektowych i jedynym, który nie
obiecuje niczego szczególnego.

Stan śledzenia po przełączeniu czyta się z bazy, a nie oddaje wprost tego,
o co poproszono w żądaniu: odpowiedź ma mówić, jak jest po zapisie,
a nie powtarzać żądanie.

Decyzja o zmianie śledzonej: przyjęcie usunięcia zdejmuje tekst, odrzucenie
go przywraca. Wersję zakłada się domyślnie, bo decyzja redakcyjna jest
punktem, do którego wypada wrócić.

Autor wersji zakładanej po decyzji bierze się z gniazda, które decyzję
podjęło. Autor zaszyty jako Operator kazałby historii twierdzić, że wersję
założył Operator także wtedy, gdy zmiany rozstrzygnął wykonawca drogą
narzędzi.

Odrzucenie propozycji zostawia dokument bez zmiany i wersji nie zakłada —
nie ma czego utrwalać.

Fragmenty liczy `policzFragmentyRoznicy` — ten sam rachunek, który panel
różnic pokazuje Operatorowi, więc numer fragmentu w żądaniu znaczy dokładnie
ten fragment, który Operator widział. Drugi rachunek fragmentów, wykonany na
potrzeby samej decyzji, mógłby ponumerować je inaczej.

Podział i złożenie idą wierszami, bo wierszami liczą się fragmenty: fragment
pominięty oddaje wiersze strony bazowej, fragment wskazany — wiersze strony
docelowej, fragment dodany i pominięty nie oddaje ani jednego wiersza. Wykaz
numerów spoza rachunku jest odmową, nie ciszą: przyjęcie fragmentu siódmego
w różnicy o trzech fragmentach byłoby przyjęciem czegoś, czego nie ma,
a odpowiedź pomyślna kazałaby czytać to jako wykonane.

## budowa/server/internal/core/adapter_modul_design_zetony_wydanie.go

Zapis i odczyt zestawów leży w `adapter_modul_design_zetony.go`.

Wszystkie postacie składa Go napisami: CSS, SCSS, konfiguracja Tailwind, moduł
JavaScript, zasoby Swift i Kotlin powstają w rdzeniu przez sklejenie napisów.
Nie ma tu ani jednego uruchomienia programu z zewnątrz i mieć nie będzie —
postać, której nie da się złożyć bez cudzego programu, byłaby u Operatora
odmową, a nie funkcją.

Wydanie do modułu jest zasobem, nie obietnicą. Wskazanie modułu docelowego nie
kończy się polem `delivered: true` postawionym z góry. Treść ląduje w
magazynie rdzenia pod sumą swojej zawartości i dostaje wiersz zasobu — ten sam
magazyn obsługuje rodziny `design.*`, `document.*`, `media.*` i `archive.*`,
więc moduł docelowy sięga po nią identyfikatorem zasobu
(`design.asset.content.get`). `delivered` mówi wtedy prawdę: bajty leżą i mają
adres. Niepowodzenie zapisu jest odmową całej komendy, nie polem
`delivered: false` postawionym obok treści oddanej wołającemu — bo Operator
zamawiał wydanie do modułu.

Rozstrzygnięcie zmiany śledzonej jest symetryczne i dlatego da się je zapisać
w czterech wierszach: przyjęcie wstawienia i odrzucenie usunięcia zostawiają
treść „po", odrzucenie wstawienia i przyjęcie usunięcia zostawiają treść
„przed". Zakres liczy się w znakach, tak jak nazywa go kontrakt (poczatek
zmiany w znakach). Cięcie po bajtach rozcinałoby polskie litery dwubajtowe
i zmiana przyjęta w tekście z znakiem „ą" wstawiałaby treść w środek znaku.

Zakres zmiany śledzonej spoza treści znaczy, że dokument zmienił się od czasu
zarejestrowania zmiany. Nowego miejsca nie zgaduje się — decyzja zapisuje się
w wierszu, a treść zostaje nietknięta.
## server/internal/zewnetrzne/brak_poswiadczenia.go

Plik stoi osobno od uwierzytelnienie.go, bo to odmowa innej klasy: tam
dostawca odpowiedzial i trzeba te odpowiedz przetlumaczyc, tutaj nikt nie
zostal zapytany, bo pod odwolaniem wskazanym w wierszu kanalu nie ma sekretu.

BrakPoswiadczenia: warstwa wyzej ma odroznic 'nie ma czym wyslac' od 'wyslano
i odbilo sie' — pierwsze naprawia sie w produkcie i ponawianie nic nie da,
drugie bywa chwilowe. Bez typu obie klasy wygladaja jak ten sam napis.

Zakładanie wersji jest wydzielone, bo trzy rodziny czynności — decyzja
o zmianach, decyzja o propozycji, przyjęcie cyfryzacji — robią dokładnie to
samo w trzech krokach. Trzy kopie tej sekwencji rozjechałyby się przy
pierwszej poprawce.
## server/internal/mowa/nagranie.go

AudioRef jest sciezka pliku na dysku Operatora, nie trescia zakodowana i nie
identyfikatorem w skladnicy. Nagranie zostaje tam, gdzie je nagrano, a
pomocnik otwiera je w miejscu. Przyjmowanie tresci przepuszczaloby kazde
nagranie przez pamiec procesu i przez gniazdo, kilkadziesiat megabajtow na
jedno zdanie, a rdzen musialby je gdzies odlozyc, czyli prowadzic druga
skladnice plikow obok tej, ktora Operator juz ma. Plik nie otwiera
deskryptora, nie czyta ani jednego bajtu tresci i nie sprawdza, czy plik jest
naprawde dzwiekiem: sprawdzenie rozszerzenia jest bramka na oczywiste
pomylki, nie rozpoznaniem formatu, bo rozpoznaje go silnik.

OtworzNagranie: kolejnosc sprawdzen idzie od najtanszego do najdrozszego i od
najbardziej ogolnego do najbardziej szczegolowego: pusty odnosnik,
rozwiniecie sciezki, istnienie, rodzaj wpisu, rozmiar, rozszerzenie.
Rozszerzenie sprawdzane jest ostatnie: gdy pliku nie ma, zdanie o
nieprzyjmowanym formacie byloby odpowiedzia na niezadane pytanie.

rozwinSciezke: sciezka domowa rozwijana z tyldy, bo powloka rozwija ja sama,
a rdzen polecenia od klienta nie dostaje przez powloke; sciezka wzgledna
liczy sie od katalogu biezacego rdzenia jako droga zapasowa dla wywolan
recznych, klient ma przesylac sciezke bezwzgledna.

Para wersji zawęża wykaz adnotacji, bo adnotacja opisuje różnicę, a nie
dokument: ta sama treść porównana z inną wersją daje inne fragmenty
i adnotacja przypięta do fragmentu trzeciego znaczyłaby wtedy co innego.
## server/internal/mowa/dostepnosc.go

Sprawdzenie poprzedza mikrofon: klient pyta o gotowosc, zanim narysuje
przycisk nagrywania, zamiast tlumaczyc jego milczenie po nieudanej probie.
Bledem jest dopiero brak odpowiedzi pomocnika: nie dalo sie go uruchomic
albo odpowiedzial czyms, co nie jest jego odpowiedzia. Pierwsze naprawia sie
instalacja, drugie zgloszeniem usterki. Dwa braki sa rozrozniane osobno, bo
maja dwie rozne naprawy: nie ma czym uruchomic skryptu (interpreter) kontra
skrypt sie uruchomil, lecz nie zastal silnika rozpoznawania.

Klucze pol Dostepnosc sa polskie, bo pomocnik jest czescia tego produktu,
a produkt jest polskojezyczny — to nie jest nazewnictwo kontraktu, ktorego
stale zostaja angielskie.

Dostepnosc (metoda): brak interpretera to ten sam rodzaj wiadomosci dla
Operatora, co brak biblioteki — jedno pytanie daje jedna odpowiedz,
niezaleznie od tego, na ktorym ogniwie lancuch sie urwal. Bledem zostaje
wylacznie odpowiedz nieczytelna: pomocnik odezwal sie czyms, co nie jest
jego odpowiedzia.

odczytajDostepnosc: wyjscie puste znaczy, ze pomocnik nie doszedl do
wypisania odpowiedzi — wtedy jedyna wiadomoscia jest diagnostyka i to ona
idzie w bledzie zamiast zdania o niepoprawnym JSON-ie, ktore niczego nie
tlumaczy. Odmowa bez powodu zostawilaby Operatora z "nie da sie" bez
zdania, co z tym zrobic.

powodZUruchomienia: samo "nie mozna uruchomic" nie wskazuje naprawy.

## budowa/server/internal/core/adapter_modul_studio_wsad.go

Wsad wykonuje operację, a nie kolejkuje ją na później: kontrakt oddaje liczbę
dokumentów przyjętych i wykaz odrzuconych. Gdyby wsad tylko wpisywał pozycje do
kolejki, obie liczby mówiłyby o zapisie do tabeli, a nie o pracy — Operator
dostałby „przyjęto 10" i nie dowiedziałby się nigdy, że siedem z nich odmówiło.
Wsad wykonuje więc operację dokument po dokumencie i dopiero wynik każdego
z nich rozstrzyga o liczbie — odmowa jednego nie przerywa pozostałych.

Osadzenie zasobu w `OsadzZasob` wstawia odwołanie do zasobu, nie jego bajty:
dokument Studia jest tekstem, a wklejona w niego grafika byłaby drugą kopią
czegoś, co już leży w magazynie pod swoją sumą kontrolną. Odwołanie idzie
zapisem Markdown, bo to jedyna postać, którą rozumie i edytor, i wyrys
podglądu, i wydanie.

Zestawienie ze źródłem w `PorownajZeZrodlem` bez żadnego wskazania materiału
wejściowego bierze plik, z którego dokument otwarto. Materiał nieodczytany
oddaje `sourceResolved: false` wraz z pustym wykazem — kontrakt pyta o to
wprost, więc odpowiedź „nie udało się" jest odpowiedzią, a nie milczeniem.

Wyszukiwanie znaczeniowe idzie dwiema drogami: kanałem modelu okna, do którego
dokument należy — model widzi znaczenie, którego miara na słowach nie zobaczy
— albo miarą arytmetyczną w rdzeniu, liczącą zbieżność słów znaczących z wagą
rzadkości. Cisza zamiast wyniku jest niedopuszczalna: dokument bez kanału
modelu ma dostać odpowiedź gorszą, ale prawdziwą, a nie żadnej. Nazwy dróg,
którymi liczy się bliskość znaczeniowa, wychodzą kontraktem w polu `mode`, bo
Operator ma wiedzieć, czy pytał model, czy rdzeń policzył sam — te dwie
odpowiedzi znaczą co innego i mają inną wiarygodność.

Funkcja `ocenyModeluStudia` zwraca `false` wszędzie tam, gdzie odpowiedzi nie
dało się wziąć za prawdę: brak rejestru kanałów, okno bez kanału, kanał
milczący, odpowiedź, której nie da się odczytać. Każdy z tych przypadków
schodzi na miarę arytmetyczną i mówi o tym wprost polem `mode`, zamiast
oddawać wynik modelu, którego nie było.

Wykaz w `slowaNieznaczaceStudia` jest krótki z zamysłu: każde słowo wykreślone
z miary jest słowem, którego Operator nie może użyć w zapytaniu, więc lista
długa szkodziłaby bardziej, niż pomaga.

W `ocenyMiaryStudia` słowo występujące w każdym akapicie nie odróżnia
akapitów, więc waży mało; słowo rzadkie waży dużo. Wynik dzieli się przez wagę
całego zapytania, więc mieści się w przedziale od zera do jedynki i da się
porównywać między dokumentami — inaczej próg `minScore` znaczyłby co innego
w każdym z nich.


## budowa/server/internal/core/adapter_modul_studio_roznice.go

Plik dopisuje się na `adapterStudia` zadeklarowanym w `adapter_modul_studio.go`.
Fragmenty różnicy liczy ten plik w locie — `DiffHunk` nie ma tabeli,
a porównanie wierszowe jest własne, bez biblioteki zewnętrznej.

Wynik operacji kontekstowej wchodzi do dokumentu, a nie stoi obok: Operator
i model pracują nad tą samą treścią w tym samym miejscu. Wynik odłożony
wyłącznie jako propozycja wymagałby drugiej powierzchni tekstowej, a takiej
nie ma. Wynik wchodzi więc w miejsce zakresu operacji, a jego przyjęcie albo
odrzucenie idzie drogą, która już istnieje: `studio.tracking.list` pokazuje
zmiany oczekujące, `studio.tracking.decide` rozstrzyga je pojedynczo albo
grupą. Propozycja zostaje zapisana dalej, bo po jej identyfikatorze
porównuje się strony w `studio.diff.compare`.

Silnik operacji kontekstowej jedzie tym samym rejestrem kanałów, co okno
rozmowy i moduł Roundtable. Odmowa pada tylko tam, gdzie czegoś naprawdę
brakuje: rdzeń złożony bez rejestru (silnik jest dodatkiem, nie warunkiem
startu), okno bez kanału, kanał spoza rejestru, silnik bez ani jednego
fragmentu treści. Każda z tych odmów nazywa brak wprost.

Propozycja zmiany powstaje po wykonaniu operacji kontekstowej, nie przed nim.
Zapis przed wywołaniem zostawiałby w panelu narzędzi propozycje puste po
każdej nieudanej próbie.

Zakres zmiany śledzonej wskazuje miejsce wyniku w treści nowej: przyjęcie
zostawia wtedy treść bez ruchu, a odrzucenie wstawia w to miejsce treść
sprzed operacji. Zakres liczony w treści starej wskazywałby po zapisie nie
ten fragment, o który szło.

Wersję zakłada `zalozWersjeDokumentu` z autorem `model` i odwołaniem do
propozycji — ta sama droga, którą idzie decyzja o propozycji. Drugiej drogi
zakładania wersji Studio nie ma.

Zakres poza treścią i zakres odwrócony w operacji kontekstowej schodzą na
cały dokument, a nie na odmowę: odmowa zapisu wyrzuciłaby pracę już
wykonaną. Cały dokument jest przy tym zakresem prawdziwym — tyle właśnie
model dostał w treści polecenia.

Parametry operacji jadą do modelu, a nie tylko do bazy. Kontrakt niesie je
jako parametry operacji wymagane przez pozycję rejestru, a wiersz polecenia
okna pracy wkłada tam słowa Operatora i nastawy suwaków koncepcyjnych.
Pominięte tutaj byłyby nastawą, której model nigdy nie przeczyta — suwak
przestawiałby wtedy pole bez skutku.

Zakres operacji niesie kontrakt polami selectionStart/selectionEnd, nie
osobnym napisem zaznaczenia. Wycinek bierze się po runach, nie po bajtach:
dokument polski ma znaki dwubajtowe i cięcie po bajtach rozcinałoby litery.

Porównanie idzie po dwóch treściach: wersja wobec wersji albo wersja wobec
propozycji, bo pole `proposalId` jest zamienne z `targetVersionId`.

Żądanie porównania, z którego nie da się policzyć ani fragmentów różnicy,
ani trafień wzorca, jest żądaniem bez odpowiedzi. Powodzenie z kopertą pustą
mówiłoby oknu, że porównano i nie ma różnic, a rdzeń niczego nie porównał:
fragmenty potrzebują dwóch stron, a wzorzec potrzebuje strony, po której ma
szukać — sama treść dokumentu stroną porównania nie jest.

Strona porównania pominięta w `trescStrony` oddaje pusty identyfikator:
`Porownaj` odróżnia po nim brak wskazania strony od strony o treści pustej.
## server/internal/mowa/pomocnik.go

Wzorzec odnajdywania jest ten sam, co w narzedzia/wpiecie.go: najpierw obok
binarium rdzenia, potem droga zapasowa, a gdy zawiodą obie, typowany blad
zamiast sciezki. Sciezka zmyslona jest gorsza od odmowy, bo odmowe da sie
powiedziec Operatorowi, a zmyslona sciezka wraca dopiero jako niezrozumialy
blad uruchomienia procesu. Plik skladla jedynie nazwe programu i wykaz
argumentow; start procesu nalezy wylacznie do portu session.Uruchamiacz.
Jedyny wyjatek to exec.LookPath — ono nie uruchamia procesu, tylko przeglada
sciezke wyszukiwania systemu, dokladnie tak jak robi to narzedzia/wpiecie.go.

interpreterPreferowany: trojka jawnie, bo gola nazwa python na wielu
systemach wciaz wskazuje wydanie drugie, w ktorym pomocnik sie nie uruchomi.

interpreterZapasowy: typowo Windows i czesc obrazow kontenerowych.

OdnajdzPomocnika: interpreter — wskazanie Operatora wchodzi wprost, bez
sprawdzania na dysku, bo moze byc nazwa do rozwiniecia przez system albo
dowiazaniem srodowiska wirtualnego, a odmowa na podstawie wlasnego
sprawdzenia uniewazniałaby to ustawienie. Gdy wskazania nie ma, szukany jest
python3, a gdy i tego nie ma, zostaje python. Ta ostatnia wartosc jest
zgadywana i moze nie istniec; odmowa przyjdzie wtedy z uruchomienia procesu,
bo tylko ono zna prawde o wykonywalnosci.

Argumenty: pomocnik oddaje tekst transkrypcji na standardowe wyjscie,
a Python bez przelacznika -X utf8 koduje je wedlug ustawien regionalnych
systemu. Na polskim Windowsie znaczy to strone kodowa 1250, w ktorej
transkrypcja rozpada sie na krzaki, zanim rdzen zdazy ja odczytac.
Wymuszenie UTF-8 w jednym miejscu jest jedyna obrona, ktora nie zalezy od
tego, kto pomocnika wola.

Rdzeń nie ma mechanizmu odczytu pliku odwołania na warstwie `core`, więc
`trescBytuPorownania` odmawia wprost zamiast oddać porównanie połowy treści.

Druga wartość zwracana przez `trescBytuPorownania` to kod strony, nie
znacznik jej istnienia. Idzie wprost do `StudioTextMatch.VersionId`, które
kontrakt opisuje jako wersję, w której wystąpiło trafienie — po tym kodzie
poznaje się, gdzie wzorzec się znalazł. Stały napis w tym miejscu sprawiłby,
że każde trafienie wyglądałoby tak samo.

Odmowa `brakStronPorownania` idzie do Operatora, więc mówi, co dopisać:
sam wzorzec bez wskazanej strony jest innym brakiem niż jedna strona bez
drugiej, choć obydwa kończą się tą samą pustą odpowiedzią.

Porównanie w `policzFragmentyRoznicy` jest własne i proste, nie algorytm
klasy diff (Myers, LCS wielowierszowy) — dla dokumentów Studio Editor (tekst
krótki, nie repozytorium kodu) rozstrzygnięcie, co się zmieniło pomiędzy
niezmienionym początkiem a niezmienionym końcem, wystarcza i jest o rząd
prościej opisać i sprawdzić niż pełny LCS.

Rodzaj fragmentu w `fragmentZmiany`: obie strony dają "changed", sama
docelowa "added", sama bazowa "removed" — panel różnic pokazuje trzy różne
oznaczenia, jedno na fragment.

Wiersz pusty do przeszukania w `szukajWzorca` (strona porównania pominięta)
oddaje wykaz pusty, nie błąd — szukanie bez treści jest pytaniem poprawnym,
tylko bez odpowiedzi.

Zmiana śledzona ma swoje wiersze w `dane/studio_adnotacje.go`.

## budowa/server/internal/core/adapter_modul_design_marketing.go

Komplet kampanii jest wydaniem, nie zapowiedzią: `design.campaign.set.build`
wyrysowuje kompozycję w każdym zamówionym rozmiarze i zakłada z każdego zasób
w magazynie. Rozmiar, którego nie udało się wydać, wraca w `failedSizes` —
komplet kampanii z połową rozmiarów wygląda bez tego pola jak komplet
gotowy, a brak wyjdzie na jaw dopiero u zamawiającego reklamę.

Makieta produktowa idzie homografią, nie prostym nałożeniem. Projekt nakłada
się na zdjęcie produktu z obrotem i kryciem. Przekształcenie liczy
`golang.org/x/image/draw` macierzą afiniczną — obrót o kąt inny niż
wielokrotność dziewięćdziesięciu stopni bez niej wymagałby własnego
próbkowania i dawał krawędzie w schodkach.

Skala kompletu kampanii bierze mniejszy ze współczynników, więc materiał
wchodzi w kadr cały. Rozciągnięcie do proporcji rozmiaru zniekształciłoby
projekt, a przycięcie ucięłoby jego część bez słowa o tym.

Płótno docelowe kompletu kampanii ma dokładny rozmiar zamówiony, a wyrys
ląduje w jego środku: materiał 1200×628 musi mieć 1200×628 pikseli, bo taki
rozmiar przyjmuje system reklamowy.

Brak klucza dostawcy w wyszukiwaniu baz zdjęciowych nie jest odmową całej
komendy — dostawca wraca w bilansie razem z powodem, a pozostali oddają
swoje wyniki.

Odmowa wyszukiwania baz zdjęciowych należy się wyłącznie sytuacji, w której
nikt nie odpowiedział: wtedy odpowiedź "zero zasobów" byłaby nieprawdą
o frazie, a prawdą o sieci.

Licencja przy wciąganiu zasobu z bazy zdjęciowej zapisuje się razem
z zasobem i jej brak jest odmową: materiał z katalogu zewnętrznego,
o którym nikt później nie powie, czy wolno go było użyć, jest gorszy niż
brak materiału.

Zasób bez zapisanej licencji jest usterką, nie zasobem — komenda importu
odmawia w całości. Bajty zostają w magazynie, są prawdziwe, ale wiersz
zasobu znika, żeby panel zasobów nie pokazał materiału bez prowenancji.

Przeskalowanie tła w makiecie produktowej unieważniłoby wszystkie cztery
liczby żądania.

Macierz przekształcenia makiety produktowej: skalowanie projektu do
zamówionego obszaru, obrót wokół jego środka, przesunięcie na wskazane
miejsce. Kolejność ma znaczenie — obrót po skalowaniu obraca prostokąt
docelowy, a nie źródłowy.

`draw.Transformer` przyjmuje macierz przejścia z układu źródła do układu
płótna, zapisaną wierszami (a b c / d e f).

Krycie częściowe w makiecie produktowej idzie przez maskę jednolitą, bo
`draw.Transformer` przyjmuje maskę w nastawach, a mnożenie składowych
obrazu źródłowego zmieniłoby jego barwy zamiast jego przezroczystości.

## budowa/server/internal/core/adapter_modul_studio_postac_pomocniki.go

Postać scala się, a nie nadpisuje: czynność na postaci podaje tylko te pola,
które zmienia. Pogrubienie ustawione nie ma prawa zdjąć kursywy, a wcięcie
lewe nie ma prawa zdjąć wyrównania. Dlatego każde scalenie bierze pole ze
żądania, gdy jest, a zastane, gdy nie ma, i nigdzie nie podmienia całej
struktury. Podmiana całości byłaby tą samą szkodą, którą w oknie widać jako
poprawę stopnia, po której znika pogrubienie.
## server/internal/core/adapter_modul_studio_arsenal.go

Rodzina cyfryzacji wola Tesseracta, wsad wola 7-Zipa. Obie potrzebuja tego
samego: uruchamiacza, zasad izolacji obowiazujacych okno i obszaru, w ktorym
proces ma pracowac. Gdyby kazda rodzina skladala ten komplet u siebie, jedna
z nich predzej czy pozniej pominelaby zasady izolacji — a to jest dokladnie
ta pomylka, ktorej punkt izolacji ma nie dopuscic. Warsztat dokumentu tedy
NIE idzie i isc nie bedzie: studio.pdf.* oraz studio.security.* pracuja
biblioteka wkompilowana w rdzen, bez ani jednego procesu potomnego. Pilnuje
tego zapora zapora_warsztatu_pdf_test.go. Rdzen zlozony bez warstwy kanalu
ma powiedziec, czego mu brakuje, i pracowac dalej w pozostalych
czynnosciach.

katalogWsadu: nazwa bierze sie z sumy nazwy archiwum, wiec rozpakowanie
tego samego archiwum dwa razy trafia w to samo miejsce, zamiast mnozyc
katalogi.

sciezkaZasobu: rdzen podaje droge, ktora dziala bez niego — wskazanie
materialu sciezka widziana przez rdzen. Milczace zejscie na pusta sciezke
dałoby odczyt pliku, ktorego nie ma, i wynik wygladajacy na prawdziwy.

Indeks górny i dolny w postaci znaku wykluczają się wzajemnie — litera nie
stoi jednocześnie nad i pod wierszem.

Repozytorium modułu Studio ogłasza tabele postaci osobnym kontraktem
(`dane.RepozytoriumPostaciStudia`), więc obszar postaci sięga po nie przez
ten kontrakt, a nie po całe repozytorium. Brak repozytorium jest brakiem
montażu rdzenia i mówi to wprost — nie udaje pustego dokumentu.

Nastawy strony domyślne jadą z jednego miejsca, bo A4 z marginesami 25 mm
jest nastawą pisma urzędowego dla całego modułu. Dwa wykazy domyślnych
rozjechałyby się przy pierwszej poprawce i dokument wczytany wyglądałby
inaczej niż założony.

Cztery składacze postaci z wierszy warstwy danych trzymają jedną zasadę:
kolumna jest prawdą, a pole JSON niesie tylko to, na co kolumny nie ma.
Dlatego najpierw odczytywany jest zapis JSON, a potem nadpisywane są pola
kolumnowe — inaczej nieświeży zapis w JSON-ie przebiłby to, co warstwa
danych wie na pewno.
## server/internal/core/adapter_modul_design_adnotacje.go

Pole note warstwy niesie JEDNO zdanie bez autora i bez watku, a przy kazdym
design.board.update jedzie z calym ukladem i wraca przepisane od nowa —
uwaga jednej osoby znikala wiec przy pierwszym przesunieciu warstwy przez
druga. Autor podany przez wolajacego bylby polem, w ktore da sie wpisac
cudze nazwisko, a oznaczenia osob w watku maja znaczyc to, co znacza.

UstawAdnotacje: watek rozpiety miedzy dwiema tablicami nie jest watkiem —
druga tablica pokazywalaby odpowiedz na uwage, ktorej u siebie nie ma.

Adnotacje (funkcja): watek zamyka sie adnotacja po adnotacji, a odpowiedz
otwarta pod zamknieta uwaga jest sprawa wciaz otwarta i ma zostac widoczna.

adnotacjaZWykazuDesignu: dwa odczyty po kodzie robilyby te sama prace drugi
raz i mogly trafic na stan zmieniony w miedzyczasie.

autorAdnotacjiDesignu: "nieznany" wpisany w kolumne wygladalby przy uwadze
jak czyjs podpis.

Autor zostaje ten, kto adnotacje zalozyl — zmiana tresci przez druga osobe
nie czyni jej autorka cudzej uwagi.

## budowa/server/internal/zewnetrzne/uwierzytelnienie.go

Surowe ciało odpowiedzi dostawcy nie jest komunikatem: nie mówi, czy to odmowa
uwierzytelnienia, wyczerpany limit, zły model czy awaria dostawcy (wszystkie
wracają tą samą drogą), nie mówi, którego poświadczenia kanał użył — sejf ma
wiele wpisów, a odwołanie bywa też nazwą zmiennej środowiskowej — i nie mówi,
co z tym zrobić. Operator z dwoma kontami u jednego dostawcy nie ma z takiego
zdania jak zgadnąć, który klucz wygasł.

Plik rozpoznaje i nazywa: bierze odwołanie, nie sekret, kod stanu i to, co
powiedział dostawca, a oddaje zdanie dla człowieka. Odwołania nie rozwiązuje
i sejfu nie dotyka — robi to `models.poswiadczenieZOdwolania`, bo kanały
powstają fabryką z wiersza rejestru i sejf mają uchwytem pakietowym.

Sekret przez ten plik nie przechodzi: wchodzi odwołanie — nazwa wpisu sejfu
albo nazwa zmiennej — nigdy wartość. `BezSekretu` jest siatką bezpieczeństwa
na drugą stronę: dostawcy wklejają fragment klucza we własny komunikat, a ten
komunikat idzie Operatorowi i do dziennika.

Wpis sejfu konta zakłada i nadpisuje `account.add` albo `account.update` polem
`credential`; wskazanie kanałowi innego odwołania robi `channel.update`.

Kontrakt nie ma pierwszorzędnego pola na odwołanie: żądanie aktualizacji
kanału niesie nazwę, model, stan włączenia i konfigurację, a odwołanie jedzie
parametrem `credentialRef` wewnątrz `config`. Samo wskazanie komendy
`channel.update` wysłałoby Operatora szukać pola, którego w żądaniu nie ma.

`OdmowaKanaluZewnetrznego` jest osobnym typem, a nie wynikiem zwykłego błędu
tekstowego, żeby warstwa wyżej mogła zapytać, czy zawiodło poświadczenie —
i wtedy nie ponawiać z tym samym kluczem — czy dostawca ma chwilową usterkę.
Napis sklejony w miejscu wywołania tej wiedzy nie niesie.

`KomunikatZCiala` wyjmuje z ciała odpowiedzi zdanie, które dostawca powiedział
o sobie sam, gdy ścieżka z wiersza rejestru (`sciezka_bledu`) bywa
nieustawiona, a wtedy bez rozpoznania Operatorowi szłoby całe surowe ciało
JSON. Rozpoznanie ogólne łapie kształty, w których błąd zwracają dostawcy
zgodni z OpenAI, Anthropic i większość bram API. Gdy wiersz ścieżkę ma, i tak
wygrywa (`ZKomunikatem`) — to jest wartość zapasowa, nie druga prawda.

`BezSekretu` wycina wartość klucza z tekstu idącego do Operatora i do
dziennika, choć sekretu tu formalnie nie ma — przychodzi z drugiej strony:
dostawcy wklejają klucz, bywa, że w całości, we własny komunikat błędu, a ten
komunikat jedzie dalej jako treść odmowy. Sekret w dzienniku jest sekretem
ujawnionym, a nikt nie zauważy tego przy zwykłej pracy. Wycina także sam
klucz z wartości nagłówka. Wołający ma pod ręką zwykle całe „Bearer sk-…",
a dostawca cytuje samo „sk-…" — porównanie wprost chybiłoby. Klucze nie mają
spacji, więc człon po ostatniej spacji jest kluczem.

`SekretZOdpowiedzi` odzyskuje wartość klucza z żądania, które tę odpowiedź
wywołało — po to i tylko po to, żeby ją z komunikatu wyciąć. Kanał rozwiązuje
odwołanie przy budowie żądania i sekretu nigdzie nie odkłada — i dobrze, bo
każde miejsce przechowania jest miejscem wycieku. `http.Response.Request`
niesie wysłane żądanie, więc wartość jest pod ręką dokładnie tam, gdzie
potrzebna, i nie przechodzi przez żadne pole ani zmienną po drodze.

Pole `Byt` struktury `Poswiadczenie`: dla drogi środowiskowej równa się
nazwie zmiennej.

Blokady fragmentów widziane przez postać dokumentu składa `blokadaZlozKontrakt`
— tabela blokad należy do obszaru kontroli pracy i to on wie, jak wiersz
przełożyć na kontrakt. Drugie przełożenie tej samej tabeli byłoby drugą
prawdą o jednym wierszu.
## server/internal/mowa/silnik.go

Odbiorcy dostaja silnik wstrzyknieciem i nie buduja wlasnego. Czego silnik
nie robi: nie syntezuje mowy — translate.speech.synthesize idzie w druga
strone (tekst na dzwiek) i potrzebuje innego silnika, ktorego ten pakiet nie
wnosi. Nie nagrywa: nagranie powstaje po stronie Operatora, a tutaj
przychodzi juz jako sciezka pliku. Nie zna kontraktu: mowi wlasnymi typami.
Dzwiek nie opuszcza maszyny Operatora: nie ma tu klienta HTTP ani adresu,
pod ktory cokolwiek by poszlo. Jedyne wyjscie na zewnatrz procesu to
uruchomienie pomocnika lokalnego.

Silnik: wolajacy sklada ustawienia z konfiguracji zasiegu (config.get)
i podaje gotowe; silnik nie siega do bazy po nastawy, bo dwie drogi do tej
samej wartosci bylyby dwiema prawdami.

uruchamiacz: bez niego silnik nie ma czym wywolac pomocnika lokalnego.

ustawienia: zlozony z konfiguracji zasiegu przy zakladaniu.

dziennik: silnik i tak rozpoznaje mowe tak samo jak z nim.

teraz: czas wpisu podaje wolajacy.

NowySilnik: ustawienia startuja wartosciami domyslnymi — brak wskazania
Operatora znaczy wartosc domyslna, nie odmowe pracy.

katalogPracy: gdy punkt izolacji jest wylaczony, proces rusza w katalogu
biezacym rdzenia; inaczej brama izolacji uzupelnia katalog wlasny okna
(session.SprawdzPolecenie).

odnotujOdmowe: blad wraca nietkniety — jego typ niesie rozroznienie, po
ktorym wolajacy nada wlasciwy kod kontraktu, a owiniecie go tutaj to
rozroznienie by zatarlo.

zapiszWpis: dziennik opisuje, co zrobiono z nagraniem, a wiersz o nagraniu,
ktorego nie wskazano, nie odpowiada na zadne pytanie — i tak zostalby
odrzucony przez sprawdzenie wpisu. Odmowa idzie wtedy do wolajacego sama
droga bledu, bez wiersza.

nowyIdentyfikator: licznik wymagalby wspolnego stanu, ktorego ten pakiet
nie ma. Blad zrodla losowosci nie przerywa transkrypcji — identyfikator
opada wtedy na znacznik czasu, bo wpis bez tozsamosci nie zapisze sie
wcale, a to gorsza strata niz tozsamosc mniej odporna na zbieg.

## budowa/server/internal/core/adapter_modul_design_szablony_materialu.go

Metody stoją na `*adapterDesignu` (`adapter_modul_design.go`).

Szablon jest układem, nie poleceniem: szablon materiału niesie rozmiar (baner
1200×628, wizytówka 90×50) i komplet warstw. Zastosowanie zakłada z niego
kompozycję gotową do pracy — po to istnieje. Wykaz szablonów, z którego nie da
się szablonu użyć, byłby spisem cudzej pracy.

Podstawienie treści idzie po nazwie warstwy: warstwa niesie adnotację (`note`)
i to ona jest jej nazwą w szablonie — „logo", „nagłówek", „zdjęcie produktu".
Podstawienie wskazuje zasób, który ma w tej warstwie stanąć. Nazwa, której
szablon nie ma, wraca w `unmatchedNames` — bilans zamiast ciszy. Cicha zgoda
oznaczałaby komplet kampanii złożony z szablonu, w którym połowa podstawień
nie weszła, a Operator dowiedziałby się o tym dopiero z wydruku.

Szablon wskazany a nieznany w `ZapiszSzablonMaterialu` jest odmową, nie cichym
założeniem nowego — wzorem szablonu promptu: Operator, który nadpisuje,
oczekuje, że nadpisał ten jeden, a nie że dostał drugi obok.

Numery stron w `stronySzablonuDoZapisuDesignu` są sprawdzane przed zapisem:
numer niedodatni nie jest numerem strony, a numer powtórzony odbiłby się od
unikatu schematu i wrócił jako awaria rdzenia oznaczona jako ponawialna —
a to jest pomyłka wołającego.

Wartość podstawienia w `podstawieniaSzablonu` musi być napisem: podstawienie
wskazuje zasób, który ma stanąć w warstwie. Liczba ani obiekt nie jest
identyfikatorem zasobu, a przepuszczone po cichu dałyby warstwę pustą przy
powodzeniu komendy.

Strony czynią z szablonu publikację. Zapis idzie po zapisie samego szablonu,
bo strona wskazuje jego klucz. Szablon bez stron zostaje jednostronicowy
i jego warstwy leżą tam, gdzie leżały — baner nie ma stron.

Wszystkie strony naraz na jednym płótnie leżałyby jedna na drugiej, bo
kompozycja jest arkuszem, nie plikiem.

Warstwy szablonu dostają własne identyfikatory w kompozycji: kompozycja jest
odtąd bytem osobnym i jej zmiana nie ma prawa ruszyć szablonu, z którego
powstała.

Strony wracają w odpowiedzi, żeby okno publikacji wiedziało, ile stron
zostało — bez tego Operator dostawałby jedną kompozycję i nie miałby po czym
poznać, że publikacja ma jeszcze dwadzieścia trzy.
## server/internal/mowa/uruchomienie.go

Pomocnik startuje tym samym uruchamiaczem, co okno rozmowy, terminal i git:
w calym drzewie jest dokladnie jedno exec.Command (injection/rozruch.go).
Dzieki temu pomocnik przechodzi przez te sama brame izolacji okna i przez
to samo obejmowanie potomstwa, co pozostale procesy — a proces Pythona,
ktory rozgalezia wlasne watki dekodera, bez objecia drzewem zostawialby
sieroty. Sekwencja jest ta sama, co w core/adapter_modul_terminal_bieg.go
i core/adapter_modul_developer_git_wykonanie.go: straznik nil, sprawdzenie
izolacji, UruchomProces, PrzejmijDrzewo(pid), defer Zwolnij, pompy obu
strumieni, select na Czekaj/timer/ctx.Done, Ubij calego drzewa przy
przekroczeniu. Odstepstwo od niej konczy sie wyciekiem procesu albo
uchwytu. Plik nie rozstrzyga, czy transkrypcja sie udala — oddaje wyjscie
i diagnostyke takie, jakie przyszly.

Wynik: sklejenie strumieni wstawilyby ostrzezenie biblioteki w srodek
przepisanego zdania, a rozpoznanie braku silnika (BrakSilnika) stracilyby
jedyne miejsce, w ktorym go widac.

Uruchom: zasady i obszar sa w sygnaturze, bo session.SprawdzPolecenie
innego ksztaltu nie ma, a bez bramy izolacji pomocnik startowalby jako
jedyny proces w drzewie poza katalogiem roboczym okna. Limit niedodatni to
odmowa, nie "bez granicy": domyslenie granicy za wolajacego ukryloby brak
wskazania az do pierwszego zawieszonego procesu na maszynie Operatora.

Punkt izolacji srodowiska: gdy okno ma go wlaczony, brama nizej to
odrzuci — i tak ma byc: rozstrzyga ustawienie Operatora, nie ten plik.

zbierz: pompy ruszaja przed czekaniem, nie po nim. Bufor potoku ma
kilkadziesiat kilobajtow; pomocnik piszacy dluzsza transkrypcje
zablokowalby sie na zapisie, a rdzen czekalby na koniec procesu, ktory
czeka na rdzen — zakleszczenie, ktore konczy sie dopiero granica czasu
i wyglada jak zawieszony silnik.

Odbior z obu pomp: gorutyny zostalyby inaczej zawieszone na zapisie do
kanalu, a to wyciek na kazde przekroczenie granicy.

Diagnostyka jedzie razem z odmowa, bo zwykle stoi w niej jedyne zdanie
mowiace, na czym pomocnik utknal.

czytajCalosc: przyrostowe zdarzenia bylyby tu kosztem bez odbiorcy. Przy
bledzie odczytu zwracane jest to, co zdazylo przyjsc — niepelna
transkrypcja niesie wiecej niz pusty wynik.

## budowa/server/internal/core/adapter_modul_design_kolor_obraz.go

Paleta z obrazu jest pomiarem, nie zgadywaniem. Barwy dominujące liczy
skupianie metodą k-średnich w przestrzeni CIE Lab, z zasiewem rozłożonym po
histogramie. Lab, nie sRGB: w sRGB odległość między dwiema barwami nie
odpowiada temu, jak różne wydają się oku, więc skupianie łączyłoby zieleń
z żółcią i rozdzielało dwa odcienie granatu. Udział barwy (`share`) jest
ułamkiem punktów przypisanych do jej skupienia — liczbą zmierzoną, nie oceną.

Symulacja wady widzenia idzie przez macierze LMS. Protanopia, deuteranopia
i tritanopia to brak jednego z trzech rodzajów czopków. Rachunek przechodzi
sRGB do LMS (macierz Hunt-Pointer-Estevez), tam zeruje brakujący kanał
zastępując go kombinacją pozostałych (macierze Brettela-Vienota-Mollona),
i wraca do sRGB. Achromatopsja jest luminancją wedle wag WCAG — tych samych,
którymi liczy się kontrast.

Wynik symulacji jest zasobem, nie base64 w odpowiedzi. Kontrakt oddaje
`DesignAsset`, więc bajty idą do magazynu pod sumą kontrolną, a wiersz
powstaje po nich — tą samą drogą, co przy wniesieniu i przy generowaniu.
Zasób wskazuje źródło polem `variantOfAssetId`: symulacja jest wariantem
obrazu, nie osobnym obrazem znikąd.

Granica punktów pomiaru palety: obraz w 24 megapikselach nie potrzebuje
wszystkich punktów, żeby oddać barwy dominujące, próbkowanie równomierne
daje ten sam wynik za setną część rachunku.

Liczba przebiegów skupiania palety: skupienia przestają się przesuwać po
kilkunastu, dwadzieścia jest zapasem, a nie nadzieją.

## budowa/server/internal/core/adapter_modul_studio_wczytanie.go

Strona sieciowa idzie przez klienta HTTP standardowej biblioteki, nie
przeglądarką: pobranie strony to żądanie HTTP i odczyt odpowiedzi. Silnik
przeglądarki wykonałby jeszcze skrypty strony — i byłby zależnością spoza
instalki, a więc odmową u Operatora. Strona zbudowana wyłącznie skryptem
zwróci tu mało treści; to jest cena znana i wybrana, a nie przeoczenie.
Migawkę strony wykonanej skryptem oddaje `browser.snapshot.get`, co kontrakt
mówi wprost w opisie tej komendy.

Skaner idzie warstwą urządzeń systemu: skanowanie prowadzi warstwa urządzeń
rozdzielona po systemie (`urzadzenia_skaner.go`) — SANE na Linuksie, WIA przez
PowerShell na Windowsie. Ta komenda nie ma własnej drogi do urządzenia i nie
ma własnej odmowy — obie rzeczy należą do warstwy, bo inaczej rozjechałyby się
z wykazem (`studio.ingest.device.list`), który pyta tę samą warstwę.

Wersja natywna Windows jest tu rzeczą rozstrzygającą: instalka natywna nie
niesie SANE i nigdy nie poniesie, więc odmowa „brak scanimage" byłaby na
Windowsie odmową na zawsze. Dlatego rozstrzygnięcie po systemie stoi w
warstwie, a nie w tej komendzie.

Rdzeń nie oddaje tu pustej kolejki pozycji, gdy czegoś brakuje: pusta kolejka
jest twierdzeniem „skanowałem i nic nie przyszło", którego rdzeń bez
odpowiedzi urządzenia nie ma prawa postawić. Każdy brak — warstwy, programu,
urządzenia, sterownika — jest odmową nazywającą, czego brakuje.

Stan „oczekuje" w pozycji wczytywania kazałby Operatorowi puścić rozpoznanie
pisma na tekście, który tekstem już jest — dlatego pozycja z adresu wchodzi
od razu jako gotowa.

Pobrane strony wchodzą do kolejki jako pozycje w stanie „oczekuje" — tak samo,
jak materiał dołożony ścieżką. Skan jest obrazem, więc tekstu jeszcze nie ma;
wpisanie stanu „gotowa" kazałoby Operatorowi przyjąć pustą treść jako wynik.

Warstwa oddająca powodzenie bez ani jednego pliku to nie jest pusta kolejka
do przekazania dalej, a usterka warstwy — i jako usterka ma zostać nazwana,
zamiast wyglądać na „skaner nic nie podał".

Czyszczenie treści strony w `tekstZeStronyStudia` jest własne i proste, bez
biblioteki czytelności: te biblioteki rozstrzygają, która część strony jest
treścią główną, a rozstrzygnięcie błędne wycina Operatorowi połowę artykułu
bez ostrzeżenia. Tutaj nic nie znika poza nawigacją i reklamą, które i tak
nie mają tekstu własnego — Operator dostaje więcej, niż prosił, a nie mniej.

Wykaz obrazów w `wciagnijObrazyStronyStudia` wraca zapisem osadzenia. Dzięki
temu obraz wciągnięty ze strony wstawia się do dokumentu tak samo jak grafika
z Design, a nie drugim sposobem zapisu.

Obraz niepobrany nie unieważnia strony: Operator prosił o tekst, a obrazy są
dodatkiem. Odmowa całości z powodu jednego martwego odnośnika byłaby odmową
wczytania artykułu.

Barwy palety wychodzą w kolejności udziału: kontrakt tak opisuje pole,
a Operator patrzy najpierw na to, czego w obrazie jest najwięcej.

Punkty całkowicie przezroczyste nie wchodzą do próby palety: barwa piksela
o zerowym kryciu nie jest barwą obrazu, a w plikach PNG z przezroczystością
bywa czernią, która przeważyłaby całą paletę.

Krok próbkowania punktów obrazu liczony jest z powierzchni: obraz mniejszy
niż granica wchodzi w całości, większy — co n-ty punkt w obu osiach.

Zasiew skupiania barw jest rozłożony po posortowanej próbie, nie losowy: dwa
wywołania na tym samym obrazie mają dać tę samą paletę, inaczej odpowiedź
zmieniałaby się przy każdym wywołaniu bez sposobu poznania, która jest
prawdziwa.

Porządek zasiewu skupiania po jasności percepcyjnej obejmuje zakres od
najciemniejszych do najjaśniejszych barw obrazu.

Nowe środki skupień liczone jako średnia w Lab — średnia w sRGB dałaby barwę
jaśniejszą od wszystkich składowych, bo sRGB niesie gamma.

Skupienie puste nie wchodzi do palety: barwa, do której nie należy ani jeden
punkt obrazu, nie jest barwą tego obrazu.
## server/internal/mowa/bledy.go

Odmowa jest osobnym typem, a nie napisem, bo obslugiwacz komendy kontraktu
rozroznia te trzy przypadki przez errors.As: kazdy dostaje inny kod bledu
i inna podpowiedz naprawy. Brak Pythona to niedokonczona instalacja
produktu, brak silnika to niedoinstalowana zaleznosc Pythona, a brak
nagrania to zla dana przyslana przez klienta. ZADNA ODMOWA NIE NAZYWA
BRAKU, KTOREGO NIE ZMIERZONO. Lancuch transkrypcji ma trzy ogniwa
dokladane osobno i naprawiane osobno: interpreter Pythona, biblioteka
faster-whisper w tym interpreterze, wagi modelu na dysku. Odmowa, ktora
przypisuje brak niewlasciwemu ogniwu, prowadzi Operatora do naprawy
bezskutecznej. Dlatego rozpoznanie ma tu wartosc "nie wiem" i przy niej
odmowa oddaje sam zmierzony powod, zamiast zgadywac rozpoznanie i naprawe.

BrakInterpretera: typ osobny od BrakSilnika, bo naprawa jest inna
i pomylenie ich prowadzi Operatora donikad: instalowanie biblioteki do
interpretera, ktorego nie ma, konczy sie drugim komunikatem o tym samym
braku. Osobny tez od BrakPomocnika, bo tam brakuje CZESCI PRODUKTU
(katalogu ze skryptem), a tu brakuje programu, ktory produkt zastaje na
maszynie. Rozpoznanie idzie po bledzie zrodlowym uruchomienia
(exec.ErrNotFound), nie po zgadywaniu z tresci wyjscia.

odmowaUruchomienia: bez tego rozstrzygniecia kazde niepowodzenie
uruchomienia szlo jako brak silnika — a wiec odmowa twierdzila
"interpreter odnaleziony" takze wtedy, gdy w tym samym zdaniu, w nawiasie,
stalo executable file not found. Operator czytal zdanie glowne i instalowal
biblioteke do interpretera, ktorego nie ma. Dopasowanie po tekscie zostaje
jako druga droga, bo blad bywa owiniety przez warstwe uruchamiania i wtedy
nie niesie juz sygnalu typowanego.

Error (BrakSilnika): zdanie NIE twierdzi, ze interpreter zostal
odnaleziony — twierdzilo tak wczesniej i bylo to twierdzenie niezmierzone:
pomocnik bywa nieuruchomiony z wielu powodow, a odmowa niosla wtedy
w nawiasie prawde przeciwna do zdania glownego. Brak interpretera
rozpoznany wprost ma wlasna odmowe (BrakInterpretera); tutaj zostaje
przypadek, w ktorym uruchomienie nie powiodlo sie z powodu
nierozstrzygnietego.

Unwrap (BrakSilnika): metoda istnieje, zeby wszystkie trzy odmowy pakietu
dały sie obsluzyc jednakowo, i zwraca nil zgodnie z umowa errors.Unwrap —
nil znaczy "lancuch konczy sie tutaj", a nie usterke.

BrakNagrania: kod bledu kontraktu jest wtedy "zly argument", nie "usterka
rdzenia".

Wariant symulacji wady widzenia wskazuje źródło: symulacja jest tym samym
obrazem widzianym inaczej, więc powiązanie jest tu prawdziwe i pozwala oknu
pokazać parę przed i po.

Kanał brakujący w macierzach wad widzenia zastępuje kombinacja dwóch
pozostałych — dlatego wiersz odpowiadający brakującemu czopkowi nie jest
zerowy, a wypełniony.

Rachunek symulacji wady widzenia idzie punkt po punkcie: obraz
kilkumegapikselowy przechodzi w czasie niezauważalnym dla Operatora,
a próbkowanie oszczędzające rachunek dałoby obraz o niższej rozdzielczości
niż źródło — czyli mniej, niż Operator wniósł.

Achromatopsja w symulacji wady widzenia jest całkowitym brakiem widzenia
barw. Wynikiem jest luminancja względna wedle wag WCAG — tych samych, którymi
rdzeń liczy kontrast, żeby dwie czynności nie miały dwóch prawd o tym, co
jest jasne.

Składowe symulacji mnoży się przez krycie, bo `image/draw` liczy w formacie
z kryciem wmnożonym.

Zasób bez treści w magazynie jest odmową przy odczycie obrazu, nie pustym
obrazem: wiersz bez bajtów jest kafelkiem, za którym nic nie leży.
## server/internal/poczta/imap_zapis.go

Szkic zapisuje sie na serwerze, nie u siebie. Szkic trzymany w bazie rdzenia
bylby szkicem, ktorego Operator nie zobaczy w swoim kliencie poczty — a caly
sens tego kroku polega na tym, zeby zobaczyl go tam, gdzie zawsze, zanim
cokolwiek wyjdzie w swiat. IMAP nie zna poprawiania w miejscu. Zmiana
szkicu to APPEND nowej wersji i skasowanie starej — tak robi kazdy klient
poczty, bo protokol nie daje nic innego. Kolejnosc jest tu zamierzona:
NAJPIERW dokladany jest nowy, POTEM kasowany stary. Odwrotnie awaria
w polowie zostawilaby Operatora bez obu wersji odpowiedzi, ktora model dla
niego napisal.

Niepowodzenie kasowania POPRZEDNIEJ wersji szkicu: Operator zobaczy wtedy
dwie wersje — co jest stanem gorszym niz jedna, ale nieporownanie lepszym
niz odmowa po udanym zapisie.

Dwie osobne komendy STORE: zadanie ma prawo zrobic naraz jedno i drugie
(np. "przeczytana i wyrozniona" z listu nieprzeczytanego i niewyroznionego).

znacznikiZmiany: zadanie unread: true ZDEJMUJE znacznik, a unread: false go
doklada; pomylenie tych dwoch kierunkow oznaczaloby, ze rdzen oznacza listy
dokladnie na odwrot, niz prosi Operator.

dolozDoFolderu: serwer bez UIDPLUS konczy APPEND powodzeniem bez podania
UID-u. Zwracane jest wtedy zero, a wolajacy odnajduje wiadomosc osobnym
szukaniem, zamiast odmawiac czynnosci, ktora sie udala.

odnajdzFolder: folder szkicow nazywa sie "Drafts" na jednym serwerze,
"INBOX.Drafts" na drugim i "[Gmail]/Wersje robocze" na trzecim. Zaszycie
ktorejkolwiek nazwy oznaczaloby, ze szkic laduje w NOWYM folderze o nazwie
"Drafts", ktorego Operator nigdy nie otwiera — czyli ze znika.

Folder domyslny: CREATE na folderze zastanym konczy sie bledem, ktory jest
tu pomijany swiadomie — "juz jest" to dokladnie ten stan, o ktory chodzilo.

## budowa/server/internal/core/adapter_modul_design_wgranie.go

Metoda stoi na typie `*adapterDesignu` zadeklarowanym w
`adapter_modul_design.go`; plik osobny wedle odpowiedzialności, jak
`adapter_modul_library_tresc.go` w module Library.

Drogi zasobu do modułu są dwie: wniesienie przez Operatora i generowanie
kanałem. Obie odkładają bajty tym samym magazynem
(`adapter_modul_design_generowanie.go`) — jeden magazyn, jedna miara formatu
i wymiarów, żadnej drugiej prawdy.

Treść zawsze ląduje w magazynie rdzenia; `sourcePath` jest źródłem bajtów, nie
miejscem składowania. Odwołanie do cudzego pliku jest wskaźnikiem na treść
żywą, a zasób Assets Panelu — jak wersja w bibliotece — ma być treścią
zamrożoną: warstwa kompozycji ułożona na obrazie nie ma prawa zacząć leżeć na
innym obrazie dlatego, że Operator posprzątał katalog pobrań. Ścieżka
źródłowa nie jest tu nigdzie zapisywana: tabela `zasob_design` nie ma kolumny
`sciezka`, a `uri` wskazuje blob.

Magazyn jest ten sam co w bibliotece, tylko nad innym katalogiem. Typ
`magazynTresciBiblioteki` (`adapter_modul_library_magazyn.go`) niesie
wszystko, czego ten moduł potrzebuje — blob pod sumą sha256, zapis atomowy
(temp w katalogu docelowym, `Sync`, `Rename`), dedup po nazwie — a drugi taki
magazyn byłby drugą prawdą o tym, jak rdzeń trzyma bajty poza bazą. Tak samo
pożycza go magazyn załączników rozmowy (`adapter_rozmowa_zalaczniki.go`).
Cena jest jedna i widoczna: teksty odmów tego typu mówią „magazyn treści
biblioteki", więc wołający ubiera je we własne zdanie modułu Design.

Wymiary pochodzą z nagłówka pliku albo nie pochodzą wcale: `image.DecodeConfig`
czyta sam nagłówek (nie rozpakowuje obrazu), więc PNG, JPEG i GIF oddają
prawdziwe piksele za cenę kilkudziesięciu bajtów odczytu. Format spoza tej
trójki (SVG, WEBP, plik wideo) zostawia `width` i `height` puste, bo pusto
znaczy „nie wiem", a zgadnięta liczba wygląda w panelu identycznie jak
zmierzona i nie da się jej odróżnić.

Import pusty w rejestracji dekoderów (`image/gif`, `image/jpeg`, `image/png`)
działa wyłącznie skutkiem ubocznym rejestracji — tak przewiduje pakiet
`image`; bez tych trzech importów każdy zasób zostałby bez wymiarów.

`przedrostekZasobuDesign` znakuje wiersze obu dróg — wniesionej i wygenerowanej
— bo jedne i drugie są tym samym bytem: zasobem z bajtami w magazynie.

`podkatalogDesignu` i `podkatalogZasobowDesignu`: osobne podkatalogi, bo
moduły nie dzielą stanu i skasowanie zasobów Designu nie ma prawa ruszyć
treści biblioteki.

`odwolanieZasobuDesignu` oddaje postać `uri` tę samą, co dla treści biblioteki
— powód i cena stoją w nagłówku `odwolanieMagazynu`
(`adapter_modul_library_magazyn.go`), bo to jedna decyzja dla całego produktu,
a nie dwie zbieżne. Baza zostaje przy ścieżce bezwzględnej. Kolumna `uri`
wiersza `zasob_design` jest bookkeepingiem rdzenia: czytają ją narzędzia
modelu (`adapter_narzedzia_obraz.go`, `_media.go`, `_archiwum.go`,
`_dokument.go`) i wysyłka poczty (`adapter_modul_poczta_wysylka.go`),
otwierając plik wprost. Przełożenie jest więc granicą kontraktu, nie zmianą
schematu — zmiana schematu ruszyłaby pięciu czytelników i wymagała migracji,
a wyciek dotyczy wyłącznie tego, co opuszcza rdzeń.

Wiersz zasobu w `WniesZasob` powstaje dopiero po utrwaleniu treści — inaczej
Assets Panel pokazywałby kafelek, za którym nie ma nic, a brak ujawniłby się
dopiero przy próbie obejrzenia zasobu. Ten sam porządek co `Wgraj` w module
Library. Etykiety idą po zapisie wiersza, bo `etykieta_zasobu_design.zasob_id`
wskazuje klucz wiersza, którego przed zapisem nie ma. Ich niepowodzenie jest
odmową całej komendy: zasób wniesiony z etykietami, które nie weszły,
wypadłby z własnego filtra w panelu i wyglądałby na zgubiony.

W `trescZasobu` ścieżka ma pierwszeństwo przed treścią, gdy przyszły obie:
bajty spod ścieżki są tym, co Operator naprawdę wskazał, a base64 bywa wtedy
podglądem złożonym przez okno (ten sam rozstrzyg co w Library). Żądanie bez
jednego i drugiego jest odmową — inaczej niż w Library, i to świadomie: tam
wgranie bez treści zakłada wersję-znacznik, a tu zasób wizualny bez treści
wizualnej nie ma czego pokazać (nagłówek `adapter_modul_design.go`).

Rozbiór base64 w `trescZasobu` idzie tam, a nie pomocnikiem Library
(`trescZadania`), wyłącznie dla odmowy: tamten mówi „moduł Library" w zdaniu,
które czytałby Operator Assets Panelu, a moduł, którego nie wołał, nie ma
prawa być autorem jego odmowy. Rachunek jest ten sam — sha256 z bajtów, ta
sama suma, którą magazyn robi nazwą bloba.

`rozpoznajObrazZasobu`: nieudany odczyt też kończy się pustkami — wniesienie
już się powiodło, a brak wymiarów nie jest powodem, żeby odebrać Operatorowi
zasób, którego bajty leżą utrwalone.

`bladZapisuZasobuDesignu` nazywa niepowodzenie utrwalenia bajtów: nośnik
pełny, brak praw do katalogu danych, ścieżka źródłowa nieczytelna, magazyn
niewpięty; komenda odmawia w całości, bo zasób, którego bajtów nie ma
nigdzie, nie ma prawa trafić do panelu jako wniesiony.
## server/internal/wiedza/bledy_test.go

Roznica jest cala po stronie czytelnika: usterka wewnetrzna mowi "cos sie
zepsulo" i zacheca do ponowienia, ktore da to samo, a odmowa mowi, czego
nie ma, ile to wazy i co zrobic, zeby bylo. Uruchomienie idzie droga
produkcyjna — Silnik, zewnetrzne.Wolaj, injection.UruchamiaczOkien — a nie
wlasnym startem procesu. Sprawdzian omijajacy te droge mierzylby skrypt,
a nie zachowanie rdzenia, i przespalby kazda zmiane w bramie izolacji albo
w rozpoznawaniu braku narzedzia.

zasiegSprawdzianu: mierzy odmowe silnika, a nie brame izolacji, ktora ma
wlasne sprawdziany.

TestBrakInterpreteraNazywaBrakINaprawe: nazwa programu jest tu celowo taka,
jakiej nikt nie zainstaluje — sprawdzian ma mierzyc brak, a nie to, co
akurat stoi na maszynie.

TestWagiBezOpisuKsztaltuDajaBrakNazwany: ta droga jest ta, ktora wskaznik
chodzi na wdrozeniu — nastawa wiedza_katalog_modeli wskazuje katalog
z gotowym modelem, a pomocnik czyta przy nim sposob skladania tokenow,
normalizacje i wymiar. Katalog bez tych deklaracji jest brakiem, bo
zgadniecie ich daloby wektory, ktore sa liczbami i nie znacza nic. Maszyna
bez biblioteki osadzen odpowie na to samo zlecenie brakiem biblioteki —
wczesniejszym w kolejnosci i rownie nazwanym. Sprawdzian przyjmuje oba, bo
mierzy KLASE odpowiedzi, nie stan maszyny.

Plik pusty przeszedlby sprawdzenie obecnosci, nie bedac wagami, wiec niesie
kilka bajtow.

TestOdmowaNazywaBrakWagePrzyczyneINaprawe: sprawdzian jest tabela, a nie
czterema funkcjami, bo mierzy jedna regule w czterech miejscach: rodzaj
dopisany bez wlasnego zdania trafi w galaz domyslna i wyjdzie z niej
zdaniem "silnik nie odpowiedzial zrozumiale", ktore nie mowi ani czego
brak, ani co zrobic.

TestOdpowiedzPomocnikaZamieniaSieWNazwanyBrak: rozjazd o jedna litere nie
psuje niczego widocznie — brak przechodzi przez odczytaj jako brak
nierozpoznany i wychodzi z niego zdaniem "silnik osadzen nie odpowiedzial
zrozumiale", czyli odmowa gorsza od tej, ktora pomocnik juz napisal.
Sprawdzian czyta napisy ze SKRYPTU, a nie powtarza ich za stalymi, bo
powtorzenie mierzyloby samo siebie.

## budowa/server/internal/core/adapter_modul_design_barwy.go

Rachunek barwy jest wkompilowany, nie wołany. Konwersje przestrzeni
percepcyjnych i harmonie idą przez `lucasb-eyer/go-colorful`, bibliotekę Go
wkompilowaną w binarium. Programu zewnętrznego do przekształceń obrazu w tej
drodze nie ma i mieć nie będzie: funkcja zależna od programu spoza instalki
jest u Operatora odmową, nie funkcją.

Luminancja liczy się wzorem WCAG, nie jasnością z biblioteki. Współczynnik
kontrastu WCAG 2.1 stoi na luminancji względnej liczonej ze składowych sRGB
po zdjęciu gamma (`LinearRgb`), z wagami 0.2126 / 0.7152 / 0.0722. Jasność
Lab (`L*`) jest inną wielkością i dałaby inne liczby — para czerni na bieli
ma dawać dokładnie 21, bo taki jest kres tej skali, i po tej liczbie
sprawdzian poznaje, że rachunek jest ten, o który chodzi.

CMYK jest przeliczeniem wprost. Przeliczenie RGB do CMYK bez profilu ICC
jest przeliczeniem naiwnym: oddaje wartości, których drukarnia użyje jako
punktu wyjścia, a nie barwę rozdzieloną pod konkretną maszynę. Rdzeń nie
udaje, że zna profil, którego nie dostał — pole `cmyk` niesie przeliczenie
wprost i tyle.

## budowa/server/internal/core/adapter_modul_design_druk_wspolne.go

Nośnik jest wiedzą rdzenia, nie zgadywanką okna: wykaz rozmiarów stoi
w miejscu wspólnym całego rdzenia (`nosniki_druku_wspolne.go`), a ten plik go
przekłada na pola kontraktu Designu. Czytają go trzy czynności: wykaz
nośników, wykaz profili (pole `paperSizes`) i wydanie do druku, gdy profil
wskazuje nośnik nazwą. Drugi wykaz — czy to w tym pliku, czy w kliencie —
rozjechałby się przy pierwszej poprawce i Operator dostałby A3 o wymiarach A4.

Profil ICC podaje, co serwer ma, a nie co warto by mieć: `design.print.profile.list`
niesie `iccProfiles` — wykaz profili leżących na tej maszynie. Gdy nie ma ani
jednego, wykaz jest pusty i to jest odpowiedź: wydanie w CMYK powie wtedy
wprost, że idzie bez osadzonego profilu, zamiast zamilczeć brak. Rdzeń nie
rozkłada profili ICC i nie przelicza barw przez nie — do tego trzeba
biblioteki, której instalka nie niesie; niesie natomiast prawdę o tym, czy
profil w ogóle jest.

Zasób z bajtów wytworzonych przez rdzeń: kafle, wykresy, schematy, makiety
produktowe i wydania do druku wracają jako zasoby magazynu, nie jako base64
w odpowiedzi. Droga jest ta sama, co przy wniesieniu: najpierw bajty
w magazynie pod sumą kontrolną, potem wiersz. Wiersz bez bajtów byłby
kafelkiem, za którym nie ma nic.

Wykaz nośników w `nosnikiDruku` jest funkcją, nie zmienną: wykaz wspólny jest
źródłem, a to jest jego przekład na pola kontraktu. Zmienna pakietowa
składana raz przy starcie byłaby trzecią kopią tych samych liczb.

Wykaz nośników modułu Design był kiedyś osobny od wykazu wspólnego rdzenia
i rozjeżdżał się z nim naprawdę — koperta DL miała w Designie wymiar wkładki,
nie koperty, a wizytówka stała położona wbrew zdaniu nad wykazem. Wykaz stoi
odtąd w miejscu wspólnym, bo dwa wykazy przy pierwszej poprawce dają
Operatorowi jeden rozmiar o wymiarach drugiego. Przestawienie wniosło
Designowi pięć pozycji, których jego wykaz nie miał: B4, B5 oraz koperty C4,
C5 i C6.

Wiersz zasobu w `zalozZasobZBajtowDesignu` przed bajtami byłby kafelkiem, za
którym nie ma nic — dokładnie ta szkoda, którą ten moduł ma w swojej
historii.

Wykaz nazwanych barw jest krótki z zamysłu: wciągnięcie pełnej listy X11
dałoby nazwy w rodzaju "papayawhip", których nikt w module nie wpisze,
a każda z nich musiałaby być tu utrzymywana.

Zapis barwy nierozpoznany jest odmową, nie barwą domyślną: podstawienie
czerni za tekst, którego rdzeń nie zrozumiał, dałoby paletę zbudowaną wokół
barwy, której nikt nie wskazał.

Zapis szesnastkowy barwy bez krzyżyka jest częstym skrótem Operatora i tak
samo jednoznacznym, więc rozpoznanie idzie po nim, zamiast odmawiać za znak.
## server/internal/wiedza/skladnica.go

Druga baza obok pierwszej to drugi plik do przeniesienia, drugi do kopii
zapasowej i drugi, ktory da sie zgubic osobno. Wektor jest bytem wtornym,
odtwarzalnym z tresci jednym przebiegiem knowledge.index, wiec jego utrata
kosztuje czas procesora, a nie wiedze Operatora. Tak samo wtorny jest
indeks FTS5 tresci biblioteki i lezy w tej samej bazie. Baza nadal nie
przechowuje pliku: bajty tresci leza w magazynie biblioteki pod suma
sha256, a tu lezy wektor i fragment tekstu, z ktorego go policzono.
Fragment jest w tabeli, bo kontrakt kaze oddac w trafieniu text, a
odtwarzanie go przy kazdym zapytaniu znaczyloby otwieranie plikow
zrodlowych i ponowne dzielenie ich na fragmenty. Skladnica mowi wlasnymi
typami, nie typami kontraktu: pakiet wiedza nie zna kontraktu, przeklad na
shared.KnowledgeHit nalezy do adaptera rdzenia.

Zapisz: fragmenty jednego dokumentu wniesione polowicznie dalyby wskaznik,
ktory o tym dokumencie wie, ale zna go do polowy — a Operator nie ma jak
tego zobaczyc: zapytanie o druga polowe wroci puste tak samo, jak wraca
dla dokumentu nieindeksowanego. Dlatego przerwanie w srodku cofa calosc.
Powtorne indeksowanie nadpisuje, a nie doklada — warunek jednoznacznosci
(zakres, kod zrodla, kolejnosc, model) czyni z zapisu upsert, wiec
dokument zmieniony i zaindeksowany ponownie ma tyle fragmentow, ile ma
tresci, a nie sume wszystkich swoich wersji.

Pozycje: Operator, ktory zmienil ustawienie modelu, ma w tabeli wektory
z dwoch przestrzeni; porownanie pytania z wektorem cudzego modelu daje
liczbe, ktora wyglada jak trafnosc i nia nie jest. Stare wiersze zostaja
w tabeli swiadomie — wracaja do uzytku, gdy Operator wroci do poprzedniego
modelu, a rebuild je czysci.

Rozdzielnik liczb wnętrza nawiasu w zapisie barwy bywa przecinkiem albo
spacją (składnia CSS Color 4), a odsetek zamienia się na ułamek od razu —
inaczej `hsl(210, 50%, 40%)` i `hsl(210 0.5 0.4)` znaczyłyby co innego,
choć to ten sam zapis.

Składowe rgb() bywają podane jako 0-255 albo jako odsetek zamieniony wyżej
na ułamek. Rozróżnia je kres: wartość powyżej jedynki nie może być ułamkiem
sRGB.

L* w zapisie CSS idzie 0-100, a biblioteka liczy go 0-1. Ułamek podany
wprost też jest tu poprawny i wchodzi bez skalowania.

Nazwa barwy wchodzi wyłącznie przy trafieniu dokładnym. Nazwa "najbliższa"
mówiłaby o barwie, której w żądaniu nie było — a pole jest niewymagane
właśnie po to, żeby brak nazwy dało się powiedzieć wprost.

Wartości niepomnożone przez krycie dawałyby przy kryciu częściowym barwę
jaśniejszą, niż wskazano, i wyrys nie zgadzałby się z podglądem w oknie.

Tekst duży wedle WCAG: od 18 punktów (24 px) albo od 14 punktów (18.66 px)
przy pogrubieniu.

`passesLargeAA` w wyniku kontrastu liczy się zawsze, a nie tylko wtedy, gdy
wołający podał rozmiar pisma: pole odpowiada na pytanie, czy ta para nadaje
się na nagłówek, co ma sens także bez wskazania rozmiaru. Rozmiar
i pogrubienie rozstrzygają natomiast o `passesAA` — para na 4.0 jest zgodna
dla tekstu dużego i niezgodna dla zwykłego, więc jedna odpowiedź na oba
przypadki byłaby nieprawdziwa dla jednego z nich.

## budowa/server/internal/wiedza/bledy.go

Odpowiedź trafień po literach ma ten sam kształt co odpowiedź trafień po
znaczeniu, więc pytający nie rozpoznałby podmiany i uznałby, że szukanej
treści w wiedzy nie ma. Treść odmowy niesie trzy człony tak samo jak odmowy
w `zewnetrzne/wolanie.go`.

`brakWagStojacych` jest osobnym rodzajem od `brakModelu`, bo naprawa jest
odwrotna: tam trzeba wagi ściągnąć, tu są już na dysku i odesłanie po nie
kierowałoby Operatora po to, co już ma.

`BrakSilnika`: brak usuwa się jedną instalacją, więc adapter rozpoznaje ten
typ i znakuje go jako niedostępność zaplecza (`adapter_modul_wiedza.go`),
a nie jako usterkę wewnętrzną. Różnica jest cała w treści, którą czytelnik
dostaje: usterka wewnętrzna nie mówi ani czego brak, ani ile to waży, ani co
zainstalować.

Pole `Silnik` rozstrzyga o treści odmowy, nie o jej kodzie: brakuje trzech
różnych bibliotek, trzech różnych kompletów wag i wskazuje się trzy różne
ustawienia, a odmowa odsyłająca po `fastembed` w miejscu, w którym brakuje
krzyżowego kodera, kierowałaby Operatora po rzecz, którą już ma.

Pole `WagaMb` — gdy waga jest nieznana, człon o wadze po prostu nie wchodzi
do zdania odmowy.

`opisNaprawyDolozonej` prowadzi trzema drogami, bo trzy różne przyczyny: brak
wag na dysku usuwa pobranie, wagi leżące a nieczytelne — wskazanie innego
katalogu albo innego wydania modelu, a brak bibliotek — jedna instalacja
w interpreterze, który już wskazano ustawieniem osadzarki (interpreter jest
w pakiecie jeden, więc i ustawienie jest jedno).

## budowa/server/internal/core/adapter_modul_design_obrazy.go

Biblioteka wkompilowana, nigdy program zewnętrzny: nie ma tu ani jednego
uruchomienia procesu i mieć nie będzie. Skalowanie robi
`golang.org/x/image/draw` filtrem `CatmullRom`, zapis PNG i JPEG —
biblioteka standardowa, dokument — `pdfcpu`, ikonę — koder w tym pliku.
Funkcja zależna od programu, którego instalka nie niesie, jest u Operatora
odmową, nie funkcją, a sprawdzian na maszynie deweloperskiej świeciłby przy
niej zielono.

Formaty, których biblioteka nie umie, są odmawiane. WEBP wchodzi (dekoder
`golang.org/x/image/webp`), ale nie wychodzi: enkodera WEBP w Go bez
zależności zewnętrznej nie ma. AVIF nie wchodzi i nie wychodzi. Wydanie
takiego formatu kończy się odmową wymieniającą formaty obsługiwane — cichy
PNG pod nazwą .avif byłby plikiem, który Operator wyśle dalej jako AVIF
i który odbiorcy nie otworzy się tam, gdzie miał się otworzyć.

SVG przechodzi bez rasteryzacji. Zasób wektorowy wydany jako SVG to te same
bajty, które leżą w magazynie: rasteryzacja odebrałaby mu jedyną własność,
dla której jest wektorem. Skala nie ma wtedy zastosowania i nie jest po
cichu stosowana.

Import pusty dekoderów wejścia niesie wyłącznie skutek rejestracji, tak
przewiduje pakiet `image`.

Wykaz formatów wydania wchodzi wprost w treść każdej odmowy formatu, żeby
Operator dostał drogę wyjścia, a nie samo zaprzeczenie.

Granica skali wydania nie jest polityką jakości, tylko granicą rachunku:
skala tysiąckrotna z obrazu 4000x4000 to bilion pikseli.

Skala pusta albo równa jedności w przeskalowaniu obrazu oddaje obraz bez
dotknięcia — przepuszczenie go przez filtr dla porządku kosztowałoby jakość
bez żadnego zysku. Miniatura ikony wydana filtrem najbliższego sąsiada
wygląda na uszkodzoną, a to jest najczęstsze wydanie tego modułu.
## server/internal/poczta/imap_odczyt.go

Zapis (szkic, oznaczenie) lezy w imap_zapis.go, wysylka w smtp.go — plik
wedle odpowiedzialnosci. Zawezanie robi serwer, nie rdzen: mail.message.list
niesie nadawce, fraze, date i "tylko nieprzeczytane"; wszystkie cztery jada
do IMAP SEARCH, wiec serwer oddaje same pasujace UID-y. Odsianie tego po
stronie rdzenia oznaczaloby sciagniecie calej skrzynki po to, zeby wyrzucic
z niej 99% — przy skrzynce Operatora z dziesiecioma tysiacami listow to nie
jest szczegol wykonania, tylko roznica miedzy odpowiedzia a zawieszeniem.
Zapowiedz idzie czesciowym odczytem (Partial), nie cala trescia. IMAP
pozwala poprosic o pierwsze N bajtow wskazanej czesci, wiec wykaz
dwudziestu listow kosztuje dwadziescia razy pol kilobajta zamiast
dwudziestu razy "ile wazyl zalacznik".

sekcjaZapowiedzi: czesc 1 jest natomiast pierwsza czescia listu, czyli ta,
ktora klienty poczty pokazuja jako tresc; dla listu jednoczesciowego RFC
3501 kaze rozumiec ja jako cale cialo, wiec jedna droga obsluguje oba
przypadki. Peek: podglad wykazu nie ma prawa oznaczyc listu jako
przeczytanego — bez tego samo wyszukanie listu zmienialoby stan skrzynki
Operatora, czynnosc uboczna, ktorej nikt nie zlecil.

Wykaz: liczba calkowita jest inna liczba niz dlugosc wykazu i to jest
zamierzone — kontrakt niesie total obok messages przycietych granica, zeby
okno moglo powiedziec "pokazuje 20 z 137" bez drugiego pytania. Mylenie
tych dwoch liczb jest usterka, nie szczegolem (wzor z design.asset.list).

WEBP i AVIF dostają zdanie odmowy własne, bo ich brak ma powód, którego nie
widać z wykazu formatów: WEBP rdzeń czyta, więc "nie znam" byłoby nieprawdą.

Skala niedodatnia i skala nad granicą w sprawdzeniu skali wydania są
pomyłką wołającego, nie awarią rdzenia.

Zawartość wpisu ikony ICO jest PNG, nie mapa bitowa DIB: postać PNG
w ikonie przyjmują wszystkie systemy od Windows Vista i wszystkie
przeglądarki, a DIB wymagałby własnej maski przezroczystości zapisanej
odwróconymi wierszami, czyli drugiego kodera obrazu w tym pliku.

Bok ponad 256 pikseli w ikonie ICO jest odmawiany, a nie przycinany po
cichu: pole szerokości w katalogu ikony ma jeden bajt, zero znaczy w nim
256, więc obraz większy zapisałby się jako ikona o boku wziętym z reszty
z dzielenia.

Osadzenie obrazu w dokumencie PDF idzie tą samą biblioteką `pdfcpu`, którą
pracuje warsztat dokumentu modułu Studio. Droga wiedzie przez PNG
w pamięci, nie przez plik pośredni: `pdfcpu` przyjmuje strumień, a plik
pośredni byłby trzecim miejscem, w którym ta sama treść żyje.

## budowa/server/internal/poczta/rozpoznanie.go

Haseł rozpoznanie skrzynek stąd nie pobiera. Thunderbird trzyma je
w `logins.json` zaszyfrowane kluczem z `key4.db`, Evolution w pęku kluczy
GNOME, mutt bywa, że jawnie w `.muttrc`. Rdzeń nie sięga do żadnego z tych
miejsc i nie próbuje ich odszyfrować: odczytanie pęku kluczy jest czynnością
innego rodzaju niż odczytanie pliku nastaw. Poświadczenie podaje Operator
komendą `mail.account.add` i idzie ono wyłącznie do sejfu rdzenia.

Outlook (Windows) trzyma nastawy w rejestrze pod
HKCU\Software\Microsoft\Office\*\Outlook\Profiles, w postaci binarnej
i innej w każdym wydaniu pakietu. Apple Mail w ~/Library/Mail/V*/MailData,
w plistach binarnych. Obu tu nie ma: rdzeń nie czyta formatów, których nie
umie przeczytać wiarygodnie, i nie zgaduje ich zawartości. Mówi wtedy, że
nic nie znalazł, a nie że nic nie ma.
## server/internal/core/adapter_modul_design_ikony_katalog.go

design.icon.library.search szuka wsrod tych wzorow, design.icon.set uzywa
ich jako punktu wyjscia, a design.icon.generate bierze, gdy pojecie trafia
w ktorys z nich. Ikony WLASNE Operatora leza w bazie (dane/design_ikony.go);
tutaj lezy to, co produkt ma bez ani jednego zapisu w bazie i bez ani
jednego pobrania z sieci. Instalka Operatora jest cienka, a arsenal stoi na
serwerze wkompilowany w binarium. Katalog ikon w bazie musialby byc zasiany
migracja i dalby sie skasowac; katalog w plikach obok binarium wymagalby,
zeby wdrozenie kopiowalo katalog danych. Katalog w kodzie jest zawsze i
jest jeden. Wszystkie wzory sa na siatce 24 i rysowane KONTUREM (obrys, nie
wypelnienie), o grubosci dwoch jednostek. Zestaw mieszajacy obrys
z wypelnieniem wyglada w oknie jak zestaw zebrany z dwoch zrodel — a to
jest dokladnie to, czemu zestaw ikon ma zapobiegac. Grubosc obrysu jest
cecha ikony (strokeWidth), wiec Operator zmienia ja i dostaje te sama ikone
cienszą, a nie inna ikone. Wzor niesie sama tresc d sciezek. Dokument SVG
sklada sie z niej na wyjsciu (svgIkonyKataloguDesignu) razem z gruboscia
obrysu, ktora Operator wskazal — gdyby wzory byly gotowymi dokumentami,
zmiana grubosci wymagalaby przepisania czterdziestu napisow.

katalogIkonDesignu: wykaz obejmuje pojecia, ktore w narzedziu pracy
naprawde wystepuja: nawigacja, praca z plikami, stany, komunikacja, media,
dane. Wciagniecie kilku tysiecy ikon dałoby nazwy, ktorych nikt nie wpisze,
a kazda musialaby byc tu utrzymywana.

wzorDlaPojeciaDesignu: Operator pisze "strzalka w prawo" i "prawo", majac
na mysli to samo.

svgIkonyKataloguDesignu: skalowanie idzie polem viewBox, nie przeliczaniem
wspolrzednych — wzor ma jedna prawde o ksztalcie, a rozmiar dokumentu jest
jego oprawa. Dzieki temu ikona na siatce 16 i na siatce 48 to ten sam
rysunek, a nie dwa zaokraglone inaczej.

sciezkiZDokumentuSvgDesignu: rdzen czyta SCIEZKI, bo tylko z nich sklada
sie kontur ikony w pakiecie i w kroju. Dokument z prostokatami i okregami
zamiast sciezek oddaje wykaz pusty, a wolajacy nazywa to wprost — cichy
pakiet z pustymi glifami bylby plikiem, w ktorym nie widac nic.

Konto Thunderbirda zapisuje się kluczami `mail.account.accountN.identities`
i `.server`, wskazującymi serwer (`mail.server.serverM.hostname`)
i tożsamość (`mail.identity.idK.useremail`). Serwer bez tożsamości też
wchodzi do wyniku rozpoznania — host i port są tym, czego Operator
najbardziej nie chce przepisywać.

## budowa/server/internal/wiedza/ustawienia.go

Klucze osadzarki są przepisane z migracji 115, a nie wymyślone w tym pliku.
Prawdą o nazwie ustawienia jest wiersz `definicja_ustawienia.klucz` ze
`store/migracja_115_wskaznik_znaczenia.sql`; cztery stałe osadzarki są jego
kopią co do znaku. Prawdą o wartości domyślnej jest ta sama kolumna po
nadpisaniu migracją 401.

Klucze przesiewu i osi obrazu idą tym samym wzorem nazw i tą samą drogą
odczytu — rozstrzyganie nastawy czyta zapis niezależnie od katalogu definicji
(`konfig/rozstrzyganie.go`), a wiersz katalogu, który wystawia je oknu
konfiguracji, zakłada migracja nastaw.

Katalog nie odmawia nieznanego klucza, więc literówka przechodzi bez błędu
i ustawienie po prostu nic nie robi (ten sam warunek pilnuje nagłówek
`mowa/ustawienia.go`).

`KluczKatalogModeli`: katalog wag rozpoznaje się w `Silnik.katalogWag`.

Wartości domyślne dziś pochodzą z `migracja_401_nastawy_wag_stojacych.sql`,
która nadpisuje wartości założone migracją 115. Rozjazd znaczyłby dwie prawdy
o tym, co zobaczy Operator, który niczego nie ustawił, i dlatego pilnuje go
sprawdzian `TestNastawyWiedzyWskazujaWagiStojace`
(`store/nastawy_wiedzy_test.go`). Wartość z bazy ma pierwszeństwo przed
stałą: rozstrzygacz zasięgu oddaje `definicja_ustawienia.wartosc_domyslna`
jako rozstrzygnięcie o pochodzeniu „domyślna", więc stałe wchodzą wyłącznie
tam, gdzie rozstrzygacza nie ma wcale (`core/adapter_modul_wiedza.go`).

`KatalogModeliDomyslny`: pusta wartość znaczy „pobierz wagi od nowa do
katalogu danych" (rozpoznanie w `katalogWag`, `pomocnik.go`), a wagi modelu
domyślnego już stoją: 4,3 GB w `/opt/danaco-modele/embedder`, z czego 2,2 GB
to wydanie ONNX, którym liczy pomocnik. Wdrożenie ma wystartować bez czynności
po instalacji, a nie pobrać drugą kopię tego, co leży na dysku. Katalog,
w którym wag nie ma, pomocnik traktuje jak pamięć podręczną i zachowuje się
dokładnie tak jak przy wartości pustej.

`ModelDomyslny`: nazwa jest nazwą wag, które leżą w `KatalogModeliDomyslny` —
wydanie ONNX modelu `BAAI/bge-m3`, transformer XLM-R o wymiarze wektora 1024,
składanie tokenów po pierwszym z nich i normalizacja wyniku — wszystko
odczytane z deklaracji leżących przy wagach, a nie przyjęte z góry
(`pomocnik_osadzen.py`). Wykaz własny biblioteki fastembed tej nazwy nie zna
i znać nie musi: model stojący opisuje się deklaracjami, a nie wykazem
wydawcy. Wektory dwóch modeli leżą w dwóch nieporównywalnych przestrzeniach
(`migracja_115_wskaznik_znaczenia.sql`); po zmianie modelu wskaźnik trzeba
przebudować.

`WagaModeluMb`: wartość jest rozmiarem wydania ONNX modelu domyślnego, czyli
tego, po co pomocnik sięga, gdy wag na dysku nie zastanie.

`dlugoscFragmentuDomyslna`: uzasadnienie liczby stoi w `fragmenty.go`.

`ModelPrzesiewuDomyslny`: wielojęzyczność, z tego samego powodu, dla którego
wielojęzyczna jest osadzarka — wiedza Operatora jest po polsku. Koder
jednojęzyczny oceniałby polskie fragmenty przez podobieństwo do angielskiego
pytania, czyli układałby kolejność gorzej niż pierwszy przebieg.

`oknoPrzesiewuTokenow`: fragment wskaźnika ma rząd 700 znaków, więc ucięcie
zdarza się wyłącznie przy fragmentach z ustawienia podniesionego do granicy.

`ModelObrazuDomyslny`: wydanie duże, nie podstawowe, bo oś obrazu wchodzi na
żądanie i liczy się raz na zapytanie, więc rozstrzyga trafność, a nie czas.

`Ustawienia` to struktura, a nie mapa, żeby literówka w kluczu rozstrzygała
się przy kompilacji, a nie przy uruchomieniu.

`Nanies`: klucz nieznany jest pomijany bez błędu, bo konfiguracja poziomu,
z którego pary przychodzą, niesie ustawienia zupełnie innych warstw produktu.
Wartość niepoprawna wraca na domyślną z tego samego powodu, dla którego robi
to silnik mowy — przerwanie budowy wskaźnika z powodu literówki w liczbie
byłoby odmową gorszą od sprowadzenia do wartości domyślnej.

`liczbaLubDomyslna` czyta liczbę własnym odczytem zamiast `strconv`, bo
granice sprawdza się i tak: wartość spoza przedziału sensownego dla okna
modelu jest tym samym co wartość nieczytelna — obie znaczą „Operator nie
podał liczby, którą da się użyć".

Profil Thunderbirda z wieloma tożsamościami jest rzadki, a wybór między
nimi należy do Operatora, nie do rdzenia: podpowiedź ma mu skrócić pisanie,
nie podjąć za niego decyzji.

Hasła w .muttrc bywają jawne (set imap_pass=) i odczyt je pomija: sekret
leżący w pliku możliwym do odczytania nadal nie jest czymś, co rdzeń
zabiera do siebie.
## server/internal/core/adapter_modul_studio_widok.go

Zasada Wlasciciela: gdzie da sie zrobic dwojako i obie drogi maja sens,
wybor nalezy do Operatora i jest JAWNYM, ODWRACALNYM ustawieniem — nie
rozstrzygnieciem wykonawcy zapisanym w kodzie. Ustawienie zyjace wylacznie
w kliencie przepada przy zamknieciu okna, wiec "przelaczenie trybu niczego
nie gubi" nie byloby prawda po ponownym otwarciu Studia. Skala widoku jest
pamietana PRZY DOKUMENCIE (Wlasciciel wymienia to wprost), a tryb
powierzchni — przy oknie. Wiersz bez dokumentu jest nastawa okna, wiersz
z dokumentem nastawa tego dokumentu; dwie tabele znaczylyby dwa odczyty
przy kazdym otwarciu okna i pytanie, ktora wygrywa. Ten sam wiersz niesie
nastawy autozapisu (adapter_modul_studio_autozapis.go) — inne kolumny,
inne polecenie zapisu. Jedno polecenie na obie grupy kazaloby widokowi
przepisywac nastawy autozapisu, ktorych nie zmienial.

UstawWidok: wartosci spoza wyliczen kontraktu wracaja odmowa nazywajaca,
co wolno — tabela i tak ich nie przyjmie, a odmowa nazwana mowi to przed
zapisem, nie po nim.

widokNastawa: wiersz zaklada NastawaPracy — jedna droga zakladania dla
obu grup kolumn. Druga rozjechalaby sie przy pierwszej zmianie wartosci
domyslnej.

## budowa/server/internal/poczta/poczta.go

Wybór bibliotek pakietu poczta:

Odbiór idzie biblioteką `github.com/emersion/go-imap/v2` (z `go-sasl` jako
zależnością). IMAP-a nie ma w bibliotece standardowej Go. Własny parser
protokołu musiałby unieść literały (`{123}` z kontynuacją serwera), zbiory
sekwencji i UID-ów, rekurencyjny BODYSTRUCTURE, odpowiedzi niezamówione
przychodzące w środku cudzej komendy, kodowanie modified-UTF-7 w nazwach
folderów oraz rozbieżności IMAP4rev1 wobec rev2. Wybrana biblioteka ma
klienta i serwer w jednym drzewie, więc jest testowana z obu stron, i mówi
obydwoma wydaniami protokołu. Cena wyboru: gałąź v2 jest rozwojowa
(v2.0.0-beta.8). Gałąź v1 jest zamrożona, a jej API oddaje wyniki kanałami,
co przy wzorcu jedno wywołanie, jedna odpowiedź wymusza gorszą obsługę
błędów. POP3, którym jedzie starsza poczta rdzenia (`dane/poczta.go`), nie
zna folderów, szkiców ani oznaczeń i nie umie APPEND, więc
`mail.folder.list`, `mail.draft.save` i `mail.message.flag` nie miałyby
czym zadziałać.

Rozbiór i składanie MIME idzie biblioteką `github.com/emersion/go-message`
(z `.../mail`). `net/mail` z biblioteki standardowej czyta nagłówki i na
tym kończy: nie schodzi w zagnieżdżony `multipart/*`, nie odkodowuje
`Content-Transfer-Encoding`, nie rozpoznaje `Content-Disposition`. Ręczny
rozbiór przez `mime/multipart` i `mime/quotedprintable` byłby
przepisaniem tej samej rekursji, a `go-message` i tak wchodzi jako
zależność `go-imap/v2`.

Wysyłka idzie biblioteką standardową `net/smtp`, bez nowej zależności.
Klient SMTP to EHLO, STARTTLS, AUTH oraz MAIL/RCPT/DATA i stdlib robi
komplet; `github.com/emersion/go-smtp` wnosi ponad to serwer SMTP, którego
ten pakiet nie buduje.

Pakiet poczta nie zna bazy, sejfu ani kontraktu i nie loguje niczego.
Sekret dostaje argumentem, trzyma go w polu i oddaje wyłącznie serwerowi,
do którego się loguje. Przekład na typy kontraktu robi adapter rdzenia.

## budowa/server/internal/wiedza/przesiew.go

Dwa przebiegi, nie jeden lepszy: osadzarka liczy wektor pytania i wektor
fragmentu osobno, więc wektory fragmentów wolno policzyć raz i trzymać
w bazie — cała biblioteka Operatora jest przejrzana jednym iloczynem
skalarnym na wiersz. Krzyżowy koder tego nie umie: czyta pytanie razem
z fragmentem, więc jego oceny nie da się policzyć zawczasu, bo pytania jeszcze
nie ma. Przejrzenie nim całego wskaźnika znaczyłoby przebieg transformera na
każdym fragmencie biblioteki, czyli minuty na zapytanie. Podział pracy jest
stąd: tani przebieg zbiera kilkudziesięciu kandydatów, drogi układa ich
kolejność.

Silnik nie ma dostępu do składnicy i nie zna kontraktu — dostaje pytanie
i teksty, oddaje po jednej ocenie na tekst. Wybór kandydatów, przycięcie do
`limit` i przekład na trafienia należą do wołającego, tak samo jak przy
osadzarce.

Proces startuje wyłącznie przez `zewnetrzne.Wolaj` i wczytuje wagi na każde
wołanie — powód stoi w nagłówku `silnik.go` i jest tu ten sam. Granica czasu
jest jedna, bo przesiew ma jedną chwilę użycia: zapytanie. Jest szersza od
granicy osadzenia pytania, bo pierwsze uruchomienie może pobierać wagi.

`KandydaciDomyslni`: pięćdziesiąt to tyle, ile koder ocenia w kilku
sekundach, a jednocześnie na tyle więcej od dziesięciu oddawanych fragmentów,
żeby przesiew miał co przestawić.

`Gotowy`: wołający sięga po nie przed zebraniem kandydatów, żeby nie
przeglądać wskaźnika po to, by zaraz odmówić.

`Przesiej`: dwa fragmenty bywają identyczne co do znaku, a pochodzić z dwóch
różnych dokumentów. Liczba ocen różna od liczby tekstów jest usterką
pomocnika, nie wynikiem gorszym, więc kończy się odmową zamiast rankingiem
przesuniętym o jedną pozycję. Wykaz pusty oddaje wykaz pusty bez uruchamiania
procesu: „przesiej nic" nie jest pytaniem o gotowość (od tego jest `Gotowy`).

`-X utf8` idzie zawsze — powód ten sam co w silniku osadzeń: pomocnik oddaje
polski tekst, a Python bez tego przełącznika koduje wyjście według ustawień
regionalnych systemu.

`GranicaKandydatowZadania`: wołający, który prosi o dziesięć fragmentów
i pięciu kandydatów, dostaje pięć — przesiew nie ma prawa dołożyć do
odpowiedzi fragmentu, którego pierwszy przebieg nie wybrał.

`PoPrzesiewie`: kosinus wektorów i prawdopodobieństwo dopasowania pary nie
leżą w jednej skali, a przesiew istnieje właśnie po to, żeby zdanie kodera
było ostateczne. Sortowanie stabilne — dwa kandydaty o równej ocenie
zachowują kolejność pierwszego przebiegu, więc to samo pytanie zadane dwa
razy daje tę samą odpowiedź.

Pakiet poczta nie nasłuchuje na żadnym porcie, nie zakłada kont, nie
kolejkuje poczty i nie pośredniczy przez żadną infrastrukturę Danaco.

Kontrakt poczty zna trzy wartości protokołu (imap, jmap, pop3). Stałej na
JMAP w pakiecie poczta nie ma, bo nie miałaby ani jednego wołacza
i wyglądałaby w kodzie jak zdolność wpięta. Skrzynkę opisaną JMAP-em
`Polacz` odrzuca zdaniem nazywającym brak wprost.

Pole SzyfrujOdbior/SzyfrujWysylke rozstrzyga je raz, przy podpięciu, adapter
rdzenia (`core/adapter_modul_poczta.go`), żeby prawda o tym, jak rdzeń
rozmawia z tą skrzynką, leżała w jednym miejscu i przeżywała restart.

WeryfikujTLS zapisuje się w kolumnie `tls_weryfikacja`, nie jest tylną
furtką.
## server/internal/core/adapter_modul_design_kroje.go

Czytaja to design.font.preview, design.font.glyphs.get,
design.font.pair.suggest, design.vector.text.path i wydania drukarskie.
Kroje sa dwojakie i roznia sie tym, czego Operator moze byc pewien: kroje
WKOMPILOWANE w binarium (rodzina Go z golang.org/x/image/font/gofont) sa
zawsze, na kazdej maszynie, bo leza w pliku wykonywalnym serwera; kroje
SERWERA — pliki .ttf/.otf lezace w katalogach krojow tej maszyny —
odczytywane, nie doinstalowywane. Pole available kontraktu mowi prawde
o tym rozroznieniu. Kroju, ktorego nie ma, rdzen nie podstawia innym:
podglad zlozony krojem zastepczym wyglada identycznie jak podglad
prawdziwy i Operator wybralby typografie, ktorej u siebie nie zobaczy.
Kontury glifow czyta golang.org/x/image/font/sfnt — biblioteka Go
wkompilowana w binarium, rozkladajaca zarowno kontury TrueType, jak
i PostScript (CFF). Programy do przegladania i rozkladania krojow tu nie
wchodza: instalka Operatora ich nie niesie.

przegladajKrojeSerweraDesignu: otwarcie stu plikow tylko po nazwe
kosztowaloby tyle, co zlozenie stu podgladow. Nazwa z tablicy wchodzi
dopiero przy wczytaniu jednego kroju. Katalog nieistniejacy nie jest
awaria: na maszynie bez pakietu krojow go po prostu nie ma, a kroje
wkompilowane zostaja.

kluczKrojuDesignu: Operator wpisuje "Go Mono", "go-mono" i "GoMono" majac
na mysli ten sam kroj, a odmowa za znak rozdzielajacy byłaby odmowa za
zapis poprawny.

sciezkaTekstuDesignu: os Y rosnie w dol — tak samo jak w kompozycji Design
Board — wiec kontur nie jest tu nigdzie odwracany i tekst nie wychodzi do
gory nogami.

glifyKrojuDesignu: kontur wchodzi do odpowiedzi, bo bez niego wykaz glifow
bylby wykazem liczb. Znak, ktorego kroj nie ma, nie wchodzi do wykazu
w ogole: kontrakt pyta o glify KROJU, a nie o zakres, ktorego Operator sie
spodziewal.

szerokoscTekstuDesignu: pomiar ma sie udac nawet dla tekstu z jednym
znakiem spoza kroju, bo sluzy do ulozenia podgladu, nie do wydania.

Identyfikator listu jest parą folder:UID, bo kontrakt daje
`mail.message.get` i `mail.message.flag` samo messageId, bez folderu, a UID
w IMAP-ie jest unikalny wyłącznie w obrębie folderu: bez folderu
w identyfikatorze nie dałoby się odnaleźć listu, którego wykaz przed chwilą
oddał, a szukanie po wszystkich folderach byłoby zgadywaniem, który z kilku
listów o tym samym UID jest ten właściwy.

Pole NazwyZalacznikow mówi, co jest w liście, a bajty ciągnie dopiero
`Pobierz`.

Folder szkiców i wysłanych bywa nazwany przez serwer inaczej (INBOX.Drafts,
[Gmail]/Wersje robocze).

Klient nie jest pulą połączeń: skrzynka należy do Operatora, więc trzymanie
w niej stałego uchwytu i zajmowanie limitu połączeń jego dostawcy byłoby
zabraniem sobie czegoś, czego nie dano.

Odmowa Polacz nazywa swój brak osobnym zdaniem: brak hosta, brak
poświadczenia, protokół nieobsługiwany, serwer nieosiągalny i odrzucone
poświadczenie to pięć różnych braków; jedno wspólne zdanie zostawiłoby
Operatora bez wskazówki, co ma naprawić.

Adres w Polacz idzie przez net.JoinHostPort tym samym powodem, co przy
wysyłce (`smtp.go`): adres IPv6 niesie własne dwukropki.

Weryfikacja łańcucha TLS wyłącza się wyłącznie wtedy, gdy Operator tak
zapisał przy skrzynce — nigdy sama, nigdy z powodu localhost.

Port 143 ze STARTTLS jest u dostawców regułą, a `DialInsecure` puściłoby
hasło Operatora łączem jawnym także tam, gdzie serwer oferuje szyfrowanie.
Biblioteka nie wystawia wariantu podnieś-jeśli-ogłoszone — ma
`DialStartTLS` (podnosi zawsze) albo `DialInsecure` (nie podnosi nigdy).
Odpowiedź BAD/NO na komendę STARTTLS znaczy, że serwer jej nie zna, i wtedy
gniazdo jawne jest jedyną drogą. Zerwany uścisk TLS albo niezaufany łańcuch
znaczą co innego — szyfrowanie jest, tylko nie jest zaufane — a ciche
zejście na jawne obniżyłoby ochronę dokładnie w chwili, w której pojawił
się powód do czujności.

Gniazdo jawne w polaczJawnymGniazdem jest jedyną drogą tak samo dla
skrzynki na tej samej maszynie, jak dla przekaźnika w sieci Operatora.

## budowa/server/internal/wiedza/wskaznik.go

Ten plik zbiera drogę od dokumentu do trafienia w jeden byt.

Wnoszenie idzie dokument po dokumencie: biblioteka Operatora bywa gigabajtem
tekstu. Jedna transakcja na całość znaczyłaby komplet wektorów w pamięci
rdzenia, a przerwanie w połowie — cofnięcie wszystkiego, co już policzono
kwadransem procesora. Dlatego granica niepodzielności biegnie po dokumencie:
każdy, który wszedł, wszedł w całości (`Skladnica.Zapisz`), przerwany
przebieg zostawia wskaźnik niepełny, ale spójny, a powtórzenie dokańcza
resztę.

Kasowanie źródła poprzedza zapis: dokument skrócony od poprzedniego
przebiegu zostawiłby inaczej fragmenty treści, której już w nim nie ma —
a wracałyby jako cytat z dokumentu, w którym ich nie ma. Sam zapis jest
wprawdzie nadpisujący (warunek jednoznaczności: zakres, kod źródła,
kolejność, model), ale nadpisuje tylko te numery fragmentów, które przyszły;
nadmiarowe zostają. Kasowanie jest w tej samej chwili co zapis i dla tego
samego źródła, więc okno, w którym dokument jest niewidoczny, trwa tyle, ile
jedna transakcja.

Czytanie źródeł dokumentu (repozytorium biblioteki, historii rozmów, plików
przestrzeni roboczej) należy do warstwy, która te repozytoria zna
(`core/adapter_modul_wiedza_zrodla.go`). Druga droga do wierszy biblioteki
założona w tym pliku byłaby drugą prawdą o tym, co Operator w niej ma.

Pole `Zakres` struktury `Dokument`: napis przychodzi z adaptera, bo to on zna
wyliczenie kontraktu.

Pole `ZrodloKod`: pusty jest stanem poprawnym dla źródeł, które kodu nie mają.

Pole `skladnica` struktury `Wskaznik`: trwałość przy bazie rdzenia
(`store/migracja_115_wskaznik_znaczenia.sql`).

Pole `dlugoscFragmentu`: wartość spoza przedziału sensownego sprowadza
`Podziel` do domyślnej — jeden strażnik tej liczby, nie dwóch.

`Wnies`: dokument bez treści dającej się podzielić jest pomijany, a nie
odmawiany; plik pusty w bibliotece nie ma prawa odebrać Operatorowi wskaźnika
pozostałych. Liczba wniesionych do tej chwili wraca razem z odmową: przebieg
przerwany w połowie zostawia wskaźnik niepełny i wołający ma wiedzieć, ile
weszło, zamiast zgadywać, czy weszło cokolwiek.

`Szukaj`: wektor pytania z modelu innego niż wektory wskaźnika daje iloczyn
skalarny, który jest liczbą i nie znaczy nic. Zawężenie odczytu po nazwie
modelu jest jedyną obroną przed tym po zmianie ustawienia — a nazwa jest
jedna, bo pyta się o nią tego, kto liczy. Wynik pusty jest odpowiedzią, nie
odmową: wskaźnik pusty albo wiedza bez związku z pytaniem znaczą „nie mam na
to nic". Inaczej niż brak silnika, który jest odmową, bo wtedy rdzeń nie wie,
czy ma coś, czy nie ma.
## server/internal/core/adapter_modul_design_wyrys.go

Dotad kompozycja jezdzila do rdzenia i z powrotem jako uklad warstw i nie
miala drogi wyjscia poza rdzen: Operator widzial tablice na ekranie i nie
mogl jej nikomu wyslac. Warstwy schodza na jedno plotno przez image/draw
(PNG), przez pdfcpu (PDF) i przez sklejenie XML-a (SVG). Przegladarki
bezglowej tu nie ma i miec nie bedzie — instalka Operatora jej nie niesie,
a funkcja zalezna od programu spoza instalki jest u niego odmowa, nie
funkcja. Kompozycja bez ani jednej warstwy z bajtami nie ma czego
wyrysowac. Rdzen odmawia zamiast oddac przezroczysty prostokat: plik,
ktory po otwarciu jest pusty, wyglada identycznie jak plik uszkodzony
i Operator nie ma z czego poznac, ze to jego tablica byla pusta. Warstwa
wskazujaca zasob usuniety z Assets Panel (kolumna zasob_id jest TEXT, nie
wiezem obcym) zostaje pominieta. Gdy przez to nie zostaje NIC, odmowa
nazywa liczbe warstw pominietych, zeby Operator wiedzial, ze tablica nie
byla pusta, tylko rozsypana.

kafleWyrysuDesignu: warstwa bez zasobu i warstwa wskazujaca zasob bez
bajtow sa pomijane, nie odmawiane — kompozycja bywa robocza i jedna
zgubiona warstwa nie ma prawa odebrac Operatorowi wyrysu pozostalych.

obszarWyrysuDesignu: Operator zaznaczyl ramka fragment kanwy i marginesu,
ktorego zazadal, nie odbiera mu sie tego po cichu.

zlozWyrysRastrowyDesignu: kolejnosc wykazu jest juz kolejnoscia warstw
(ORDER BY kolejnosc, id), wiec warstwa pozniejsza klada sie na
wczesniejszej, tak jak na kanwie. Skalowanie kazdej warstwy idzie filtrem
CatmullRom, tym samym co przy wydaniu zasobu: dwa wyrysy tej samej grafiki
nie maja prawa roznic sie ostroscia zaleznie od tego, ktora komenda
powstaly.

zlozWyrysSvgDesignu: tresc jest osadzona, nie dowiazana. Odsylacz do pliku
w magazynie rdzenia nie otworzy sie nigdzie poza ta maszyna, a wyrys ma
byc plikiem, ktory Operator wysyla dalej — dokument wskazujacy cudze
sciezki bylby pusta ramka u kazdego odbiorcy.

## budowa/server/internal/poczta/list.go

Rozdzielenie rozbioru i złożenia MIME na dwa pliki dawałoby dwie prawdy
o tym, co rdzeń uważa za treść listu.

List składa się raz, wysyła dwiema drogami. Szkic zapisany w folderze
Drafts i list nadany SMTP-em to dokładnie ten sam dokument — dlatego
`zlozWychodzacy` ma dwóch wołających (`imap_zapis.go` i `smtp.go`), a nie
dwie kopie. Gdyby szkic składał się inaczej niż wysyłka, Operator
oglądałby przed wysłaniem list inny niż ten, który wyjdzie w świat.

Treść bierze się z `text/plain`, a `text/html` dopiero przy jego braku.
Model ma analizować treść, a nie znaczniki; list wyłącznie HTML-owy
oddaje się takim, jaki jest — obcięcie go do treści wymagałoby
renderowania HTML, czego rdzeń nie robi i o czym nie kłamie.

Rejestracja dekoderów stron kodowych działa skutkiem ubocznym importu
pustego: pakiet podstawia `message.CharsetReader`. Bez niego polski list
w iso-8859-2 albo windows-1250 — a takich w skrzynce Operatora jest pełno —
rozbiera się błędem "nieznana strona kodowa" zamiast oddać treść.

Wątek wiąże się identyfikatorem listu, na który ten odpowiada, a przy jego
braku — własnym. Dzięki temu list otwierający wątek i wszystkie odpowiedzi
w nim mają tę samą wartość, bez pytania serwera o THREAD, którego część
serwerów nie zna.

Odkodowanie transportowe pozostaje nietknięte w zapowiedzi: część pobrana
częściowo bywa urwana w połowie czwórki base64 albo w połowie sekwencji
=XX, więc dekoder i tak nie miałby czego domknąć. Zapowiedź listu
w quoted-printable wygląda przez to nieco surowo (=C5=BC zamiast ż) —
i tak jest uczciwiej niż zgadywać brakujące bajty.

Załącznik bez nazwy istnieje i ma bajty — nazwa zastępcza jest tu etykietą,
a nie zmyśloną własnością listu: bez niej magazyn rdzenia nie miałby czym
opisać zasobu.

Klient poczty Operatora układa rozmowę po nagłówku References — bez niego
odpowiedź wyląduje w jego skrzynce obok wątku, a nie w nim.
## server/internal/core/adapter_modul_studio_zadania.go

Petla, ktora te zadania prowadzi, stoi w adapter_modul_studio_petla.go.
Podzial jest taki: tutaj mieszka to, CZYM petla pracuje, tam to, JAK
pracuje. Operator mowi "napisz pismo w tej sprawie na podstawie tych
materialow, sprawdz terminologie i przygotuj wersje do druku" i ma dostac
cztery zadania, nie jedno. Rozklad idzie dwiema drogami po kolei: najpierw
modelem (kontrakt mowi wprost "brak znaczy rozklad ulozony przez model"),
a gdy modelu nie ma albo jego odpowiedz nie sklada sie w zadania —
rozkladem wlasnym, po czynnosci nazwanych w zleceniu. Druga droga nie jest
atrapa: rozpoznaje czynnosci z wykazu StudioTaskKind po slowach, ktorymi
Operator o nich mowi, i zachowuje ich kolejnosc ze zdania. Zlecenie,
w ktorym nie da sie rozpoznac ani jednej czynnosci, daje jedno zadanie
rodzaju custom niosace cale zlecenie — bo "nie rozpoznalem" nie moze
znaczyc "rozklad pusty". Rozklady leza w magazynie pamieciowym tego
procesu. Warstwa danych nie ma dzis tabeli rozkladu ani zadania,
a zakladanie jej nalezy do odcinka fundamentu, nie tutaj. Skutek jest
nazwany, nie przemilczany: rozklad nie przezywa ponownego uruchomienia
rdzenia. Tabela plan_studio i zadanie_studio sa wypisane w sprawozdaniu
jako brak do domkniecia.

petlaZnacznikCzynnosci: dopisanie wariantu jest dopisaniem wiersza, nie
zmiana rozkladu; kolejnosc zadan bierze sie z miejsca, w ktorym slowo
stoi w zleceniu, bo Operator wymienia czynnosci w tej kolejnosci, w jakiej
maja sie wykonac.

petlaRozlozZlecenie: "sprawdz terminologie i napisz pismo" ma dac korekte
po napisaniu tylko wtedy, gdy Operator tak powiedzial — dwa zadania tej
samej czynnosci z jednego zdania byłyby zdublowana robota. Rozklad
zlecenia dokumentowego jest z natury szeregowy — nie da sie poprawic
jezyka pisma, ktorego jeszcze nie ma. Rozklad podany wprost przez
Operatora albo przez model swoje zaleznosci niesie wlasne i nie sa
nadpisywane.

## budowa/server/internal/core/adapter_modul_design_eksport.go

Eksport różni się od `image.convert` celem: konwersja zakłada nowy zasób
w magazynie i tam się kończy, eksport oddaje bajty gotowe do zapisania poza
produktem. Wydanie nie zostawia więc po sobie ani wiersza, ani bloba —
magazyn Designu nie ma sprzątania, a każde obejrzenie ikony w trzech
skalach zakładałoby trzy zasoby, których nikt nie zamawiał.

Partia eksportu zlicza bilans, nie milczy. Odmowa jednego zasobu nie
wstrzymuje pozostałych (kontrakt), a odrzucony trafia do wykazu wraz
z powodem — panel ma pokazać, czego nie wydał, a nie oddać krótszą listę
bez słowa.

Braki wspólne całej partii — format, którego rdzeń nie zapisuje, skala poza
granicą, pusty wykaz zasobów — są odmową całej komendy, nie odrzuceniem
każdego zasobu z osobna: wynik z pięcioma odrzuceniami o tym samym
powodzie ukrywałby, że pomyłka leży w żądaniu, a nie w zasobach.

## budowa/server/internal/core/adapter_modul_design_szablony.go

Metody stoją na `*adapterDesignu` (`adapter_modul_design.go`).

Wcześniej historia promptów i szablony żyły w oknie do zamknięcia karty
przeglądarki i tyle o nich było wiadomo. Prompt Builder wypracowywał
polecenie, a praca znikała po zamknięciu karty.

Historia domyka też prowenancję. Zasób niesie pole `promptId`, którego rdzeń
dotąd nie wypełniał, bo nie było przekładu klucza wiersza promptu na kod
kontraktu. Odczyt zasobu bierze teraz kod promptu podzapytaniem
(`dane/design_zasoby.go`), a komenda historii oddaje drugą stronę tego samego
powiązania: prompt wraz z kodami zasobów, które z niego powstały.

Nadpisanie w `ZapiszSzablonPromptu` idzie po `templateId` z żądania; brak
zakłada szablon nowy. Operator, który nadpisuje szablon, oczekuje, że
nadpisał ten jeden, a nie że dostał drugi obok — a wykaz z dwoma szablonami
tej samej nazwy wygląda dokładnie tak, jak wygląda zgubiona praca.

W `HistoriaPromptow` kanał bywa odmawiał i wtedy prompt jest zapisem próby —
czyli dokładnie tego, po co Operator do historii sięga, żeby poprawić
polecenie i wydać je jeszcze raz.

Drugi warunek w `sprawdzSzablonPromptuDesignu` nie jest formalnością:
nadpisanie cudzego szablonu przez podanie jego kodu byłoby zmianą stanu,
którego wołający nawet nie widzi.

Wynik kontraktu partii eksportu niesie wykaz wydań (files []DesignExportFile),
z których każde ma własną nazwę, typ treści i skalę, a archiwum byłoby
jednym plikiem bez tych pól. Spakowanie zbioru plików umie moduł Library
(archive.*) i tam jest jego miejsce — udawanie tutaj, że pliki są
spakowane, dałoby wykaz wydań, których treść nie jest tym, co zapowiada
mediaType.

Skala poza granicą w trzeciej pozycji unieważnia całą partię eksportu,
a wydanie dwóch pierwszych zasobów przed odmową zostawiłoby Operatora
z połową zamówienia.

Partia skal jednego zasobu jest jednym zamówieniem (zestaw gęstości),
a wydanie @1x bez @2x wygląda na komplet i nim nie jest.

Wariant gęstości nazwy wydania dokleja się z krotności skali (@2x) — tak
nazywa się zestawy gęstości i tak je czytają narzędzia po drugiej stronie.

Nazwa wydania jedzie do klienta jako propozycja, ale klient zapisuje pod
nią plik, więc separator ścieżki i kropki wiodące wychodzą w jednym
miejscu składania nazwy, w `oczyscNazwePlikuDesignu`.

Kolumna `format` niesie to, co zmierzył nagłówek albo zadeklarował
Operator, a wydanie svg przepuszcza bajty i to one muszą być wektorem —
stąd sprawdzenie `czyZapisWektorowyDesignu` idzie po treści.

## budowa/server/internal/wiedza/fragmenty.go

Akapit to myśl wydzielona przez autora tekstu, natomiast cięcie co N znaków
rozłupuje zdanie w połowie słowa i daje cytat bezużyteczny. Tekst rozpada się
najpierw na akapity (pusty wiersz); akapit krótszy niż docelowa długość
doklejany jest do sąsiada, dopóki mieści się w granicy, bo akapit
jednozdaniowy osadzony osobno niesie za mało kontekstu, żeby dać się
odróżnić od innego jednozdaniowego. Akapit dłuższy od granicy rozpada się na
zdania (kropka, wykrzyknik, pytajnik, koniec wiersza), a zdanie dłuższe od
granicy cięte jest po ostatniej spacji przed granicą, nigdy w środku słowa.

Każdy fragment poza pierwszym zaczyna się od ostatniego zdania fragmentu
poprzedniego. Bez tej zakładki zdanie stojące na styku dwóch fragmentów traci
połowę kontekstu po każdej stronie granicy; koszt zakładki to około jednej
piątej więcej wektorów.

Granica 1100 znaków wynika z okna modelu, przy którym te liczby powstały:
`paraphrase-multilingual-mpnet-base-v2` obcinał wejście na 384 tokenach
podziału XLM-R, a polszczyzna kosztuje w nim około trzech znaków na token.
Fragment dłuższy niż okno zostaje obcięty po cichu, a cytat byłby wtedy
dłuższy niż to, co model przeczytał. Model domyślny jest dziś inny — `bge-m3`
przyjmuje 8192 tokeny (`max_position_embeddings` w opisie jego wag) — więc
okno przestało być ciasne i granica przestała być jego odwzorowaniem. Zostaje
jednak nietknięta, bo jest zarazem granicą cytatu: fragment jest tym, co
wraca Operatorowi i modelowi jako przytoczenie, a przytoczenie na kilka
tysięcy znaków przestaje być przytoczeniem. Docelowe 700 znaków zostawia pod
tą granicą zapas na zakładkę i na słowa łamane na kilka tokenów.

Pole `Kolejnosc` struktury `Fragment` wchodzi do tożsamości wiersza
wskaźnika, żeby powtórne indeksowanie tego samego dokumentu nadpisywało
fragmenty, a nie dokładało ich drugi komplet.

Akapit dłuższy od granicy w `naCzesci` rozkłada się na zdania, a zdanie
dłuższe od granicy — na kawałki cięte po ostatniej spacji.

`zdania`: skróty pisane z kropką („ul.", „art.") rozdzielą tu zdanie
w miejscu, które zdaniem nie jest; skutek jest kosmetyczny — fragment o jedno
zdanie krótszy — a obroną byłby dopiero słownik skrótów polszczyzny
utrzymywany w rdzeniu.

`ostatnieZdanie`: cięcie pada na pierwszej spacji za granicą liczoną od
końca, a gdy spacji tam nie ma — na najbliższym początku znaku UTF-8, więc
zakładka nie zaczyna się od rozłupanej litery.

## budowa/server/internal/wiedza/silnik.go

Osadzenia liczy się w dwóch chwilach — przy budowaniu wskaźnika
(knowledge.index) i przy zapytaniu (knowledge.search) — i musi je liczyć
ten sam model tym samym sposobem. Wektor dokumentu policzony jednym
modelem, a wektor pytania drugim, dają iloczyn skalarny, który jest liczbą
i nawet wygląda sensownie, a nie znaczy nic.

Proces silnika startuje wyłącznie przez `zewnetrzne.Wolaj`: w całym
produkcie stoi dokładnie jedno exec.Command, a każde uruchomienie idzie tą
samą bramą izolacji okna i tym samym obejmowaniem potomstwa. Proces
Pythona liczący na wielu wątkach bez objęcia drzewem zostawiałby sieroty
po każdym przekroczeniu czasu.

Każde wołanie silnika wczytuje model od nowa: prawie cały czas zlecenia to
start procesu i wczytanie wag (rząd gigabajta) do pamięci, a nie samo
porównanie wektorów. Proces rezydentny skróciłby zapytanie, ale wymaga
dwukierunkowej rozmowy z procesem żyjącym między żądaniami, a
`zewnetrzne.Wolaj` prowadzi rozmowę jednorazową i jest jedyną dozwoloną
drogą startu procesu.

Granica czasu silnika jest dwojaka. Pierwsze uruchomienie pobiera wagi
modelu, więc granica budowania wskaźnika jest liczona w minutach; zapytanie
ma wagi już na dysku i granica jest liczona w sekundach. Jedna wspólna
granica byłaby albo za krótka na pobranie, albo tak długa, że zawieszone
zapytanie wyglądałoby na pracujące.

Limit budowania wskaźnika obejmuje pobranie wag przy pierwszym uruchomieniu.

Wielkość partii silnika: uruchomienie procesu kosztuje kilka sekund
(wczytanie wag), więc partia ma być duża; wykaz w pliku JSON o kilkuset
fragmentach to kilkaset kilobajtów, czyli nic.

## budowa/server/internal/core/adapter_modul_design_kompozycje.go

Typ adaptera `adapterDesignu` i jego konstruktor deklaruje
`adapter_modul_design.go`; ten plik dokłada wyłącznie metody obszaru
kompozycji.

Zapis jest zawsze pełny, jedną ścieżką. `design.board.update` nadsyła całą
listę warstw na nowo — kontrakt (`DesignBoardUpdateRequest.Layers`) nie zna
trybu częściowej zmiany. Adapter nie dogaduje różnicy względem stanu
zastanego; warstwa danych (`ZapiszKompozycje`) usuwa i wstawia komplet od
nowa w jednej transakcji. Brak `BoardId` zakłada kompozycję nową — adapter
nadaje wtedy nowy identyfikator zewnętrzny przed wywołaniem repozytorium, bo
repozytorium samo zna wyłącznie „załóż albo nadpisz" po tym identyfikatorze.

Zasób warstwy wskazuje identyfikatorem zewnętrznym, nie kluczem obcym.
Warstwa może wskazywać zasób spoza Assets Panelu w chwili zapisu — kolumna
`warstwa_kompozycji_design.zasob_id` jest w `migracja_048_design.sql` typu
TEXT — więc adapter nie sprawdza istnienia zasobu przed zapisem.

Odczyt w `Kompozycje` jest drugą stroną zapisu: kod kompozycji nadaje adapter
przy pierwszym `design.board.update` (`plansza-…`), więc bez tej komendy
Operator po odświeżeniu okna nie miałby ani planszy, ani czym o nią zapytać,
a kolejny zapis zakładałby kompozycję nową obok zastanej.

Repozytorium nie przycina wykazu kompozycji — gdyby `Total` i długość wykazu
miały prawo się różnić, byłaby to strona, a nie komplet. Kompozycja bez
warstw zostaje w wykazie, bo pominięcie takiej kompozycji ukryłoby przed
Operatorem planszę, którą sam założył.

Warstwy idą osobnym odczytem na kompozycję, tak samo jak składa je
`ZapiszKompozycje`, bo schemat `migracja_048_design.sql` trzyma je w osobnej
tabeli. Jedno okno niesie jednostki plansz, więc pętla odczytów jest tu
tańsza niż złączenie rozklejane potem w pamięci.

Kolejność w wykazie kontraktu jest kolejnością renderowania, więc brak
`Order` w `przelozWarstwyDoZapisu` odziedzicza pozycję w liście.

Sprawdzenie Gotowy idzie przed budowaniem wskaźnika i przed zapytaniem, bo
odmowa "nie ma czym" jest dla Operatora czymś innym niż "liczyło i się
wywróciło". Pomocnik przygotowuje model i wraca, więc przy pierwszym razie
pobierze też wagi.

Osadzanie partiami, nie wszystko naraz: wykaz kilkudziesięciu tysięcy
fragmentów w jednym pliku zlecenia zająłby pomocnikowi pamięć
proporcjonalną do całej biblioteki. Kolejność wektorów odpowiada
kolejności tekstów i to jest warunek — wołający wiąże je pozycją, nie
treścią. Wykaz pusty do osadzenia nie jest pytaniem o gotowość silnika,
tylko pracą, której nie ma (od pytania jest `Gotowy`).
## server/internal/core/adapter_modul_studio_zmiany_modelu.go

Wlasciciel oznaczyl to wymaganie jako WAZNE i nazwal je swoim glownym
narzedziem kontroli nad praca modelu w dokumencie. Model, ktory przestawil
kroj albo wciecie, ma byc widoczny tak samo jak ten, ktory dopisal akapit.
Zmiana postaci bez zmiany liter NIE MOZE byc niewidzialna — dlatego
rachunek bierze zmiany sledzone rodzaju formatowanie na rowni
z wstawienie i usuniecie, a wykaz oddaje tez czynnosci dziennika autora
model, bo tam stoi cale drzewo postaci. Agentow Operator zaklada w module
Agents dowolnie wielu i dwoch moze pracowac nad jednym dokumentem naraz.
Przelacznik pokazujacy ich jako jednego bylby bezuzyteczny wlasnie wtedy,
kiedy jest najbardziej potrzebny. Dlatego odpowiedz niesie byAgent —
rozbicie wedle kodu agenta i podagenta. Wymaganie mowi wprost: dokument
ma wrocic do stanu sprzed pracy modelu Z ZACHOWANIEM zmian Operatora
naniesionych w tym czasie. Przywrocenie wersji skasowaloby prace
Operatora. Dlatego cofa sie POJEDYNCZE zmiany sledzone autora model —
od konca dokumentu, bo zakresy liczone sa w tresci sprzed decyzji —
i pojedyncze czynnosci dziennika autora model.

PrzeskocDoZmianyModelu: "nastepna zmiana modelu" znaczy nastepna, do
ktorej okno ma przewinac, a nie nastepna zapisana.

Wykaz pusty: licznik zero mowi to wprost, a change pozostaje pusty
z zamyslem kontraktu ("brak znaczy koniec wykazu").

zmianyModeluWybierz: trzy kopie tego przesiewu rozjechalyby sie przy
pierwszej poprawce i licznik przy przelaczniku przestalby zgadzac sie
z tym, po czym Operator skacze.

## budowa/server/internal/wiedza/podobienstwo.go

Produkt wozi sterownik `modernc.org/sqlite`, czyli SQLite przepisany na Go,
bez zależności od C. Rozszerzenie wektorowe `sqlite-vec` wydawane jest jako
biblioteka natywna (`.so`, `.dylib`, `.dll`) ładowana przez
`sqlite3_load_extension`, a sterownik nie wystawia drogi włączenia ładowania
rozszerzeń — `load_extension('vec0')` kończy się `SQL logic error: not
authorized (1)`. Nie miałby też jak jej wystawić: kod przepisany na Go nie
wciągnie do siebie natywnej biblioteki C. Drogą alternatywną byłaby zamiana
sterownika na wariant z CGO w całym produkcie.

Zostaje przegląd zupełny: wektory leżą w tabeli, a podobieństwo liczy się
w Go dla każdego wiersza. Przy 768 wymiarach i jednym rdzeniu przegląd stu
tysięcy wektorów zajmuje rząd sześćdziesięciu milisekund, a biblioteka rzędu
tysiąca dokumentów daje kilkanaście tysięcy fragmentów, czyli kilkanaście
milisekund na zapytanie — wobec sekund, które zajmuje wczytanie wag przez
pomocnika (nagłówek `silnik.go`). Indeks przybliżony zaczyna się opłacać
dopiero powyżej setek tysięcy pozycji.

Wektory wchodzą do bazy znormalizowane. Podobieństwo kosinusowe to iloczyn
skalarny podzielony przez iloczyn długości, a długość wektora dokumentu nie
zmienia się między zapytaniami; wektor podzielony przez własną długość już
przy zapisie sprawia, że porównanie jest samym iloczynem skalarnym, a wynik
wprost kosinusem.

`rozmiarLiczby`: cztery bajty, nie osiem — osiem bajtów zapisywałoby zera po
przecinku, których w wektorze nie ma — dwukrotność miejsca i dwukrotność
odczytu z dysku za zero dokładności.

`Znormalizuj`: wektor zerowy wraca bez zmiany, a nie przez dzielenie przez
zero, bo jego iloczyn skalarny z czymkolwiek wynosi zero, więc nigdy nie
wygra rankingu, podczas gdy wektor z NaN psułby porównania.

`NaBajty`: porządek bajtów ustalony jawnie, a nie odziedziczony po maszynie,
bo baza bywa przenoszona między maszynami razem z katalogiem danych, a
wektor odczytany w odwrotnym porządku bajtów to nie gorsze trafienia, tylko
liczby bez żadnego związku z tekstem.

`blisko`: różna długość wektorów znaczy różne modele, więc te dwa wektory
nie leżą w jednej przestrzeni i porównanie ich części byłoby liczbą udającą
wynik.

`Najblizsze`: kosinus niedodatni znaczy, że fragment nie ma z pytaniem nic
wspólnego, a oddanie go jako trafienia dałoby cytat do zbudowania
odpowiedzi z niczego.

`posortujMalejaco`: sortowanie stabilne, żeby dwa fragmenty o równej ocenie
wracały zawsze w tej samej kolejności — inaczej to samo pytanie zadane dwa
razy dawałoby dwie różne odpowiedzi bez żadnej zmiany w wiedzy. Wspólne dla
obu przebiegów: pierwszy układa po kosinusie, drugi po ocenie kodera
(`przesiew.go`), a wymóg powtarzalności jest ten sam.

`WSetnych`: zaokrąglenie w górę od połowy; `score` większe od stu byłoby
liczbą spoza zakresu obiecanego przez kontrakt.
## server/internal/core/adapter_modul_studio_katalogi.go

Operator zapisuje wlasny prompt raz i siega po niego w kazdym dokumencie,
a zespol pracujacy nad projektem ma widziec lancuchy tego projektu, nie
cudze. Zasieg jest wiec czescia tozsamosci wpisu, nie ozdoba: ten sam
identyfikator w dwoch zasiegach to dwa rozne wpisy. chain.run sprawdza
lancuch, dokument i zakres, po czym oddaje identyfikator przebiegu wraz
z liczba krokow. Samych krokow nie puszcza synchronicznie: kazdy jest
operacja kontekstowa wychodzaca do kanalu modelu, a odpowiedz kontraktu
nie niesie ich wynikow — niesie przebieg, po ktorym poznaje sie zadania.
Wykonanie kroku po kroku prowadzi petla wykonawcza okna, a postep idzie
zdarzeniem. Udawanie tu wykonania oddawaloby steps policzone z definicji
i przebieg, za ktorym nic nie stoi.

ZastosujSzablon: pole bez wartosci zostaje w tresci WIDOCZNE jako
znacznik, a nie znika — dokument z pustym miejscem po polu wyglada na
kompletny, a nie jest.

wypelnijSzablon: wartosci przychodza surowym JSON-em, bo kontrakt opisuje
je typem json — szablon nie ma z gory znanego zbioru pol i miec nie moze.
Tresc nieczytelna zostawia szablon nietkniety: dokument ze znacznikami
jest odroznialny od dokumentu wypelnionego, a dokument wypelniony
bzdura — nie.

UstawFormatDokumentu: przestawienie formatu bez zamiany zostawia ten sam
tekst pod nowa nazwa formatu, co bywa wlasciwe (markdown oznaczony jako
tekst czysty). Zamiana tresci miedzy formatami nalezy do obszaru
dokumentow i idzie osobna komenda — tutaj bylaby druga droga do tej samej
czynnosci.

## budowa/server/internal/zewnetrzne/wolanie.go

Produkt daje modelowi arsenał czynności — obraz, dźwięk, dokument, archiwum
— a większość z nich to opakowanie dojrzałego programu, nie pisanie go od
nowa w Go. Bez pakietu zewnetrzne każda rodzina narzędzi zbudowałaby własne
uruchamianie procesu, własny limit czasu i własne sprzątanie potomstwa —
cztery prawdy o jednej rzeczy, a przy czwartej ktoś zapomniałby ubić drzewo.

Sekwencja pakietu jest przepisana z `mowa/uruchomienie.go`. W całym drzewie
stoi dokładnie jedno exec.Command (`injection/rozruch.go`); wszystko inne
idzie portem `session.Uruchamiacz`, przez tę samą bramę izolacji okna i to
samo obejmowanie potomstwa. Każde odstępstwo od tej sekwencji jest
wyciekiem procesu albo uchwytu — a binarium przetwarzające film rozgałęzia
wątki tak samo jak dekoder mowy.

Pakiet zewnetrzne nie rozpoznaje formatów, nie czyta obrazów i nie
tłumaczy komunikatów ImageMagicka na polski. Rozpoznanie należy do
rodziny narzędzi, która zna kształt odpowiedzi swojego programu.

Pole Pakiet wchodzi do treści odmowy, żeby Operator nie musiał szukać, czym
dociągnąć brakujące narzędzie.

Wyjście i diagnostyka Wynik są dwoma strumieniami, tak samo jak w silniku
mowy: na wyjściu stoi wynik pracy (bywa nim binarna treść obrazu), na
diagnostyce ostrzeżenia programu. Sklejenie ich wstawiłoby ostrzeżenie
w środek pliku PNG.

Rodzina narzędzi ma odróżnić BrakNarzedzia, żeby powiedzieć Operatorowi,
co dociągnąć.

Limit czasu w Wolaj jest obowiązkowy — tak samo jak w silniku mowy.
Transkodowanie filmu bez granicy potrafi zająć maszynę na godziny,
a program, który utknął na uszkodzonym pliku, nie kończy się nigdy. Limit
niedodatni jest odmową, nie brakiem granicy: wołający, który granicy nie
podał, o niej zapomniał.

Silnik mowy dziedziczy środowisko, bo interpreter Pythona musi odnaleźć
własne biblioteki; narzędzia arsenału nie mają powodu widzieć zmiennych
rdzenia. Gdy okno ma włączony punkt izolacji środowiska, brama niżej i tak
rozstrzyga — ale domyślna wartość ma być węższa, nie szersza.

## budowa/server/internal/core/adapter_modul_design_kolekcje.go

Metody stoją na `*adapterDesignu` (`adapter_modul_design.go`).

Kolekcja jest bytem osobnym od etykiety. Etykieta jest słowem, kolekcja ma
nazwę, opis i porządek — powód rozdziału stoi w nagłówku migracji 230.

Przypisanie jest dokładką albo odjęciem, nigdy zastąpieniem. Zasoby wchodzące
do kolekcji sprawdza się przed zapisem: kolekcja pełna kodów, za którymi nie
stoi żaden zasób, byłaby dokładnie tym rodzajem koperty, który ten moduł już
raz oddał — wykazem bez treści. Zasób nieznany jest więc odmową całej
komendy, nie cichym pominięciem, bo `changed` nie ma pola na „ten nie wszedł".

`ZalozKolekcje`: kontrakt mówi wprost, że panel dostaje identyfikator,
znacznik czasu i licznik zasobów takie, jakie naprawdę zostały zapisane.

`PrzypiszDoKolekcji`: przy zdejmowaniu zasób z kolekcji bywa już usunięty
z Assets Panelu, a odmowa zdjęcia go zamknęłaby Operatorowi jedyną drogę do
posprzątania kolekcji.

`kolekcjaKontraktuDesignu`: `AssetCount` i długość wykazu są tu zwykle równe,
ale prawdą o kolekcji jest ta z bazy, a nie ta z tego, co akurat udało się
wczytać.

Program dołożony do pakietu leży poza ścieżką wyszukiwania systemu i po
samej nazwie nie wystartowałby, więc do uruchomienia idzie ścieżka
odnaleziona.

Proces bez uchwytu drzewa, zostawiony tak, znaczyłby sierotę poza
rejestrem rdzenia.

Pompy w zbierz ruszają przed czekaniem: bufor potoku ma kilkadziesiąt
kilobajtów, program piszący obraz na wyjście zablokowałby się na zapisie,
a rdzeń czekałby na koniec procesu, który czeka na rdzeń. Zakleszczenie
kończy się dopiero granicą czasu i wygląda jak zawieszone narzędzie.

Odbiór z obu pomp po zakończeniu i bezwarunkowo, także przy przerwaniu:
gorutyny zostałyby inaczej zawieszone na zapisie do kanału.

Plik obrazu urwany w połowie jest gorszy niż jego brak, dlatego wynik
przerwany nie wraca jako udany.

Bez opisu diagnostyki Operator czyta tylko, że narzędzie zakończyło się
niepowodzeniem, i nie wie nic więcej. Pierwsze zdanie programu prawie
zawsze niesie powód, a reszta jest śladem stosu.
## server/internal/core/adapter_modul_studio_blokady.go

Wymaganie rozstrzygajace Wlasciciela. Model wola komendy rdzenia tak samo
jak klient, wiec blokada pilnowana wylacznie przez okno jest pozorna —
ominalby ja bez wysilku. Sprawdzenie stoi wiec na drodze KAZDEJ komendy
zmieniajacej dokument (adapter_modul_studio_blokady_zapora.go), a nie
w poszczegolnych obslugiwaczach: obslugiwaczy zmieniajacych dokument jest
w Studiu ponad setka i pisanych przez czterech wykonawcow, a zapora wpieta
w rejestr obejmuje tez te, ktorych jeszcze nie ma. Zderzenie zmiany
z blokada nie ma odpowiedzi "wolno / nie wolno", ma trzy: CALOSC
W BLOKADZIE daje odmowe NAZWANA (ktory fragment i jaka blokada — cicha
bezczynnosc bylaby najgorsza mozliwa odpowiedzia, bo Operator myslalby, ze
model wykonal polecenie); CZESC W BLOKADZIE — zmiana wchodzi POZA blokada,
a odpowiedz niesie BILANS (co zmienione, co pominiete i przez ktora
blokade — odmowa calosci bylaby tu nieproporcjonalna, a przemilczenie
pominiecia zakazane); POZA BLOKADA — zmiana wchodzi bez slowa. Blokada
jest skierowana przeciw wykonawcom, nie przeciw wlascicielowi dokumentu:
Operator zmienia fragment zablokowany BEZ przeszkod. Blokada dzialajaca
takze na Operatora jest osobnym, jawnym ustawieniem (scope = everyone),
a nie zachowaniem domyslnym. Zdejmuje blokade WYLACZNIE Operator.
Wykonawca, ktory uzna, ze fragment wymaga zmiany, zaklada propozycje na
marginesie (studio.markup.add o rodzaju suggestion) — i tyle. Dlatego
studio.lock.remove nie stoi tez w wykazie narzedzi modelu. Blokada wisi
przy DOKUMENCIE, nie przy wersji, wiec przywrocenie wczesniejszej wersji
jej nie gubi — i to jest caly mechanizm, bez ani jednej linii kodu
w drodze przywracania. Zakresy nadazaja za trescia przez
PrzesunBlokadyFragmentow.

blokadaNaZakresie: wykaz trzech blokad w tresci bledu nie pomaga bardziej
niz jedna, a wykaz pelny Operator ma w studio.lock.list.

blokadaUzgodnijTresc: uzgodnienie musi rozdzielic "ten fragment zmiany
wolno" od "tego nie wolno", a do tego trzeba miejsc, w ktorych zmiane da
sie PRZERWAC. Wiersz jest takim miejscem od poczatku: policzFragmentyRoznicy
liczy roznice wierszami, a studio.diff.hunk.apply wierszami przenosi
fragmenty. Znak bylby miejscem gestszym, ale uzgodnienie znakami
wymagaloby dopasowania klasy LCS, ktorego rdzen nie ma. Gdy strony maja
TYLE SAMO wierszy w srodku, wiersze daja sie sparowac jeden do jednego
i uzgodnienie jest dokladne: wiersz trafiony blokada zostaje zastany,
pozostale wchodza. Gdy liczba wierszy sie rozni, sparowac ich nie sposob
bez zgadywania, KTORY wiersz odpowiada ktoremu. Wtedy srodek jest jedna
caloscia: styka sie z blokada — nie wchodzi w ogole i wraca pominieciem;
nie styka — wchodzi w calosci. Zgadywanie parowania byloby tu gorsze niz
odmowa, bo wstawiloby tresc w srodek zablokowanego cytatu i nazwalo to
bilansem.

## budowa/server/internal/zewnetrzne/kod_odmowy.go

Zdanie dla człowieka (`uwierzytelnienie.go`, `brak_poswiadczenia.go`) to
dopiero połowa tego, co Operator czyta. Drugą połowę — kod kontraktu i
wynikającą z niego ponawialność — klient dokleja do wpisu rozmowy sam. Rdzeń
domyka turę jednym kodem zapasowym dla każdej odmowy kanału
(`channel_unavailable`), a katalog kontraktu (`shared.KodyPonawialne`) uznaje
ten kod za ponawialny. Bez kodu własnego odwołany albo wygasły klucz
przychodziłby więc podpisany „ponowienie ma sens", choć zdanie obok mówi
„wymień klucz".

`protocol.BladZeZrodla` najpierw pyta `errors.As`, czy błąd niesie już swój
kod, i dopiero gdy nie niesie — nakłada kod zapasowy wołającego. Obie odmowy
tego pakietu niosą kod przez `Unwrap`, więc kod zapasowy nie wygrywa i żaden
plik poza tym pakietem nie musi się zmieniać.

Wykaz kodów w `KodOdmowy`, każdy opisany zdaniem z `shared/contract.go`:
401 → `not_authenticated`, dostawca nie uznał klucza, nieponawialny — ten sam
klucz odbije się tak samo; 403 → `permission_denied`, klucz zna, czynności
nie da; 402 → `permission_denied`, konto jest uwierzytelnione, a mimo to nie
ma prawa do czynności, bo powodem jest rozliczenie, nie zasięg klucza; 429 →
`rate_limited`, ponawialny, i słusznie — to jedyna z tych odmów, która mija
sama; 5xx → `channel_unavailable`, awaria po stronie dostawcy jest właśnie
chwilowa.

Odrzucone żądanie (400, zły model, złe ciało) nie jest wprawdzie chwilowe,
ale jego naprawa leży w parametrach wiersza kanału, a `validation_failed`
mówiłby wołającym, że to ich żądanie było niezgodne z kontraktem rdzenia —
dlatego zostaje przy `channel_unavailable`.

`KodBraku`: gdyby brak poświadczenia szedł jako `not_authenticated`, Operator
dostałby maszynowe potwierdzenie, że zawinił jego klucz — a klucza nikt tam
jeszcze nie czytał, bo nie było czym. Pozostałe powody to rzeczywiście brak
uwierzytelnienia i żadne ponowienie ich nie usunie.

`Unwrap`: na tym opiera się przeniesienie kodu do Operatora. Bez tego ogniwa
kod musiałby wybierać `core/strumien_odpowiedzi.go` i każde inne miejsce,
które odmowę kanału przenosi do protokołu. `Error()` pozostaje zdaniem dla
człowieka, a `errors.As(err, &OdmowaKanaluZewnetrznego{})` działa dalej —
ogniwo dokłada się za typem, nie zamiast niego.

## budowa/server/internal/zewnetrzne/odnajdywanie_test.go

Instalka wnosi programy towarzyszące do katalogu `pomocniki` obok binarium
rdzenia. Dopóki szukanie szło wyłącznie ścieżką wyszukiwania systemu, taki
program był niewidoczny: Operator dostawał odmowę „nie ma na tej maszynie",
stojąc nad plikiem, który przyszedł razem z produktem.

Sprawdziany mierzą rozstrzygnięcie, nie stan maszyny: czy pakiet idzie przed
ścieżką i czy katalog o nazwie programu nie udaje programu.

`TestPakietMaPierwszenstwoPrzedSciezkaSystemu`: wersja dołożona do pakietu
jest tą, którą sprawdzono przed wydaniem; wersja zastana na maszynie bywa
starsza albo okrojona i nie jest niczyją obietnicą. `sh` jest właściwym
materiałem na spór o pierwszeństwo, bo stoi na ścieżce systemu każdej
maszyny, na której się uruchomi.

## budowa/server/internal/core/adapter_modul_design_kanal.go

Metoda wyboru kanału stoi na *adapterDesignu (adapter_modul_design.go).

Kanał obrazowy to nie dowolny kanał modelu. Rejestr rdzenia prowadzi
w większości do modeli tekstowych oddających fragmenty text; kanał
obrazowy poznaje się po kluczu adaptera "obrazy" (models.AdapterObrazy,
wiersz rodzaju api z parametrem adapter). Rozróżnienie pada przed
wywołaniem — inaczej moduł wysłałby prompt kanałowi tekstowemu i mówił
Operatorowi, że kanał nie oddał obrazu, zamiast wprost, że wskazał zły
kanał.

Odmowy wyboru kanału są rozdzielone, bo Operator naprawia je różnymi
ruchami: dopisać kanał obrazowy, poprawić identyfikator, włączyć kanał,
wskazać kanał obrazowy zamiast tekstowego, naprawić parametry wiersza.
Żadna z nich nie schodzi po cichu na kanał inny niż wskazany — obraz
wygenerowany innym silnikiem byłby obrazem cudzym, a Operator nie miałby
jak się o tym dowiedzieć.

Braku poświadczenia plik wyboru kanału nie rozstrzyga. Wiersz bez
credentialRef buduje się celowo (models/adapter_obrazy.go), a odmowa pada
w chwili wywołania, zdaniem samego kanału: rozróżnia ono brak odwołania od
odwołania bez wartości w sejfie, czego ten moduł nie widzi. Powielenie
tego sprawdzenia tutaj dałoby dwie prawdy o poświadczeniach i uboższą
treść odmowy.

kanalObrazowyZadania oddaje wiersz kanału, nie sam identyfikator, bo
wołający wypisuje identyfikator kanału w treściach odmów.

Pierwszy kanał obrazowy znaczy: pierwszy w kolejności wykazu rejestru,
czynny i mający zbudowany adapter; wiersz włączony bez adaptera pomija się
milcząco przy szukaniu domyślnego, bo Operator o niego nie prosił
i odmówiłby przy pierwszej turze.

Wykaz kanałów niesie także wiersze nieczynne, więc da się odróżnić "nie ma
takiego kanału" od "jest, ale wyłączony". Gdyby szukać wyłącznie wśród
czynnych, obie sytuacje zlałyby się w jedną odmowę i Operator nie
wiedziałby, czy pomylił identyfikator, czy zapomniał włączyć kanał.

Kanał tekstowy włączony jest tą samą pomyłką, co kanał tekstowy wyłączony
— włączanie go niczego nie naprawi, więc zdanie o czynności wysłałoby
Operatora w złą stronę.

Brak domyślnego kanału obrazowego nie jest brakiem produktu: adapter
obrazowy w rdzeniu jest, magazyn jest, pole channelId w kontrakcie jest.
Dlatego zdanie odmowy mówi, jak taki kanał założyć, zamiast opisywać
granicę rdzenia.

Wykaz kanałów niesie też wiersze nieczynne, więc da się odróżnić brak
kanału od kanału wyłączonego.

## budowa/server/internal/core/adapter_modul_studio_magazyn.go

Po bajty magazynu sięga pięć rodzin naraz: wydanie archiwum, paczka
redakcyjna, raport różnicy, wyrys stron i osadzenie zasobu w treści.
Gdyby każda składała sobie sumę kontrolną i wiersz zasobu u siebie,
wystarczyłaby jedna pomyłka w kolejności najpierw treść, potem wiersz,
żeby w panelu zasobów pojawił się kafelek bez zawartości.

Kolejność jest ta sama, co w całym rdzeniu (design.asset.upload): bajty
trafiają do magazynu pod swoją sumą sha256, dopiero potem powstaje wiersz,
który je wskazuje.

Okno puste w odlozTrescStudia znaczy "wynik do żadnego okna nie należy" —
odlozWynikArsenalu oddaje wtedy zasób złożony z ręki, z odwołaniem, którym
plik da się odczytać. Kontrakt trzech czynności wydania ma windowId jako
pole nieobowiązkowe, a Operator, który go nie podał, o kafelek nie prosił.

Treść obszerna wersji dokumentu, odłożona poza bazą przez TrescOdwolanie,
doczytuje się z magazynu, a nie odmawia: wydanie połowy dokumentu byłoby
wydaniem dokumentu innego niż ten, który Operator widzi w edytorze.

## budowa/server/internal/zewnetrzne/odnajdywanie.go

Program dołożony do pakietu (pomocniki/<program>) ma być znaleziony
niezależnie od tego, co Operator ma w systemie, i ma mieć pierwszeństwo,
bo to jego wersję sprawdzono przed wydaniem. Wersja zastana na maszynie
bywa starsza, nowsza albo okrojona; zgodność z nią nie jest niczyją
obietnicą. Gdy w pakiecie nie ma nic, zostaje ścieżka systemu — tak
pracuje maszyna deweloperska i tak działa produkt na maszynie, gdzie
Operator ma własną instalację programu.

Pakiet zna dwa układy — płaski pomocniki/<program> oraz własny katalog
pomocniki/<program>/<program> dla programów niosących własne środowisko.
Powód rozdzielenia stoi przy miejscaPakietu, bo tam składa się ścieżki.

Rdzeń i katalog pomocniki wychodzą z jednego pakowania i stoją obok
siebie, więc pierwszym miejscem jest katalog binarium, które właśnie
pracuje. Drugim — ta sama ścieżka względna liczona od katalogu bieżącego:
obsługuje uruchomienie z go run, gdzie binarium leży w katalogu
tymczasowym, a obok niego nie ma nic z produktu. Sekwencja jest ta sama,
którą od dawna stosuje pomocnik mowy (mowa/pomocnik.go).

Ścieżka programu wskazana wprost jest wskazaniem Operatora albo montażu
i nie podlega szukaniu — sprawdza się wyłącznie, czy plik nadaje się do
uruchomienia.

Dwa układy pakietu, nie jeden. Układ płaski (pomocniki/<program>) jest
podstawowy i wystarcza programom, które są jednym plikiem albo niosą obok
siebie kilka własnych bibliotek. Nie wystarcza programom, które przynoszą
całe własne środowisko: PowerShell to kilkaset zestawów .NET wraz
z podkatalogiem Modules, LibreOffice ma własne drzewo program, Chromium
swoje pliki wydania, a rembg cały osadzony Python. Windows szuka bibliotek
w katalogu pliku wykonywalnego, więc zsypanie ich wszystkich do jednego
pomocniki znaczy kolizję nazw bibliotek: kilka pakietów wnosi plik o tej
samej nazwie i innej zawartości, a wygrywa ten, który skopiowano później.
Taka wygrana jest przypadkiem, nie rozstrzygnięciem, i objawia się dopiero
u Operatora.

Drugi układ, pomocniki/<program>/<program>, to własny katalog o nazwie
programu, z nietkniętym układem wydania wewnątrz. Program widzi swoje
biblioteki obok siebie, a żaden inny pakiet mu w to nie wchodzi. Płaski
idzie pierwszy, żeby dotychczasowe pakiety zachowały się dokładnie tak jak
przedtem, a katalog własny był drogą dla tych, których płasko położyć się
nie da.

Katalog własny nie jest plikiem wykonywalnym i nie nosi rozszerzenia .exe,
choć plik w jego środku nosi.

Na Windowsie plik wykonywalny nosi rozszerzenie .exe, a nazwa programu
w deklaracji narzędzia go nie niesie — pakowanie nie zmienia nazw, więc
rozszerzenie dokłada się w miejscaPakietu.

## budowa/server/internal/core/adapter_modul_design_tresc.go

Komenda oddania treści zasobu domyka lukę, przez którą moduł miał
precedens szkody. Pole uri zasobu jest ścieżką w systemie plików rdzenia —
przeglądarka nie wczyta spod niego niczego, także wtedy, gdy zasób powstał
bez zarzutu. Bez tej drogi panel zasobów pokazywał kafelek z odwołaniem,
którego nie da się otworzyć, i wyglądało to identycznie jak zasób bez
bajtów.

Zasób nieznany to odmowa not_found. Zasób bez odwołania — odmowa nazywająca
brak treści. Odwołanie prowadzące donikąd — odmowa nazywająca odwołanie.
Treść większa niż granica wołającego — odmowa podająca zmierzoną wielkość,
nigdy treść ucięta: klient, który dostałby połowę pliku ze stanem ok,
zapisałby ją jako plik cały.

Suma kontrolna liczy się z bajtów odczytanych, nie z nazwy bloba. Nazwą
bloba jest wprawdzie suma jego zawartości, więc obie wartości powinny być
równe — i właśnie dlatego liczy się je osobno: rozjazd znaczy, że plik pod
odwołaniem przestał być tym, za który się podaje, a przemilczenie tego
byłoby oddaniem cudzej treści pod nazwą zasobu.

Granica wołającego sprawdza się przed złożeniem odpowiedzi, ale po
zmierzeniu treści: odmowa ma podać wielkość rzeczywistą, a nie samą
wiadomość o przekroczeniu.

Kontrakt czyni contentBase64 niewymaganym właśnie po to, żeby postać
odsyłania mogła oddać miarę i sumę bez bajtów.

Format zapisany w wierszu jest zmierzony przy wniesieniu z nagłówka pliku.
Zgadywania typu treści po nazwie pliku nie ma: nazwa jest wolnym tekstem
Operatora i mówi o pliku dokładnie tyle, ile Operator w nią wpisał.

DetectContentType dokleja parametr zestawu znaków do typów tekstowych.

## budowa/server/internal/core/adapter_modul_design_zasoby.go

Wniesienie zasobu leży w adapter_modul_design_wgranie.go, etykiety
w adapter_modul_design_etykiety.go — podział wedle odpowiedzialności.

Blob leży pod sumą swojej zawartości, więc dwa zasoby o identycznej treści
(to samo tło wniesione w dwóch oknach) dzielą jeden plik; skasowanie go
przy usunięciu jednego zasobu odebrałoby treść drugiemu, który o niczym
nie wie. Zliczania odwołań rdzeń nie prowadzi i ten moduł go nie zakłada,
więc bajty usuniętego zasobu zostają w katalogu danych do czasu, aż
magazyn dostanie sprzątanie.

Zasób odczytuje się przed zapisem oznaczenia ulubionego, bo komenda niesie
identyfikator kontraktu, a oznaczenie wisi na kluczu wiersza; przy okazji
ten sam odczyt odróżnia zasób nieznany (odmowa not_found, panel ma po czym
poznać, że jego wykaz jest nieaktualny) od usterki bazy.

Odpowiedź niesie stan z bazy, nie echo żądania — wzorem
UstawEtykietyZasobu. Podstawienie z.Favorite do zasobu odczytanego przed
zapisem opisywałoby stan, którego w bazie może nie być; stąd drugi odczyt.

Ustawienie stanu ulubionego, który już obowiązuje, jest drogą udaną.
Kontrakt nie pyta, czy się zmieniło, tylko żąda stanu docelowego — odmowa
za powtórzenie zabrałaby oknu prawo do wysłania tego, co Operator widzi na
przełączniku.

Zasób czyta się przed usunięciem, żeby było co rozgłosić: design.asset.changed
niesie cały zasób obok rodzaju zmiany (DesignAssetChangedEvent.Asset), więc
po usunięciu wiersza nie ma już z czego złożyć zdarzenia — a zdarzenie
deleted z pustym zasobem powiedziałoby panelowi "coś zniknęło" bez
wskazania czego. Odczytany zasób wraca wołającemu osobnym wyjściem, bo
odpowiedź kontraktu niesie samo removed (patrz zarejestrujDesign).

Zasób nieznany przy usunięciu nie jest odmową. Kontrakt pyta wprost, czy
zasób istniał i został usunięty — removed: false odpowiada na to pytanie
prawdziwie, a odmowa kazałaby oknu obsługiwać błąd tam, gdzie stan
docelowy (zasobu nie ma) już obowiązuje. Powtórzone usunięcie tego samego
kodu też oddaje false.

Rozpoznanie zasobu nieznanego przy usunięciu idzie po błędzie braku,
a nie odesłaniem go dalej.

Etykiety czyta się jeszcze przed usunięciem zasobu: etykieta_zasobu_design
znika kaskadą razem z wierszem, więc po UsunZasob zdarzenie niosłoby zasób
bez etykiet, których w chwili usunięcia miał pełny zestaw.

## budowa/server/internal/wiedza/osadzarka.go

Silnik osadzeń (silnik.go) startuje proces Pythona przez zewnetrzne.Wolaj
i potrzebuje do tego trójki session.Okno + session.Zasady + session.Obszar.
Trójka jest sprawą uruchamiania procesu, a nie sprawą wskaźnika, który
pyta wyłącznie o zamianę tekstów na wektory. Port zdejmuje ze wskaźnika tę
jedną zależność; podział na fragmenty, normalizacja, zapis do bazy,
odczyt, iloczyn skalarny, ranking i przycięcie do limit zostają tym samym
kodem, który pracuje na maszynie docelowej.

Port ma w drzewie jedno wypełnienie produkcyjne — Silnik, co potwierdza
var _ Osadzarka niżej; adapter rdzenia wiąże silnik i podaje go
wskaźnikowi. Drugie wypełnienie żyje wyłącznie w pliku _test.go, więc nie
da się go wkompilować w binarium.

Model() należy do portu, bo nazwa modelu wchodzi do wiersza wskaźnika przy
zapisie i zawęża odczyt przy szukaniu. Gdyby wskaźnik brał ją z ustawień,
a wektory liczył kto inny, rozjazd zamieniłby wskaźnik w zbiór wektorów,
o których nie wiadomo, z czym wolno je porównywać.

Wskaźnik jest jeden na maszynę i wszystkie jego zlecenia jadą tym samym
zasięgiem platformy; podawanie trójki przy każdym wołaniu sugerowałoby,
że wolno ją zmieniać między fragmentami jednego dokumentu, a wtedy część
wskaźnika powstawałaby pod inną izolacją niż reszta.

## budowa/server/internal/core/adapter_modul_design_etykiety.go

Nadawanie etykiet jest jednym wypełnieniem portu, plik osobny wedle
odpowiedzialności, tak jak adapter_modul_design_kompozycje.go.

Komenda domyka drogę zapisu etykiet: design.asset.list zawęża wykaz polem
tags, a warstwa danych niesie UstawEtykietyZasobu, więc brakowało wyłącznie
przejścia komenda, uchwyt, baza. Własnej migracji ta komenda nie potrzebuje.

Odmowa nadania etykiet należy się wyłącznie zasobowi, którego rdzeń nie
zna — pomylenie tych przypadków odebrałoby Operatorowi jedyny sposób
odetykietowania zasobu, więc pustka nie jest tu sprawdzana ani odrzucana.

Zasób nie zostaje przejściowo bez etykiet, gdy zapis padnie w połowie, bo
podmiana etykiet idzie jedną transakcją warstwy danych.

Warstwa danych pomija wartości puste i zwraca zestaw uporządkowany, więc
odesłanie z.Tags wprost mówiłoby o stanie, którego w bazie nie ma. Drugi
odczyt etykiet jest ceną prawdy o skutku.

PromptId w podmianie etykiet zostaje pusty z tego samego powodu, co przy
Zasoby: repozytorium nie ma przekładu klucza wiersza promptu na kod
kontraktu — pole nie jest zgadywane.

czyBrakZasobuDesignu woła się z design.asset.remove, dla którego brak
zasobu jest odpowiedzią (removed: false), a nie odmową.

Semantyka podmiany etykiet, nie dokładania, jest zapisana w kontrakcie.
Odpowiedź niesie etykiety odczytane z bazy, nie echo żądania.

## budowa/server/internal/poczta/smtp.go

Pakiet niczego nie pyta o zgodę: wysyła, gdy dostanie polecenie, i zostawia
ślad w dwóch miejscach, bo każde odpowiada na inne pytanie: kopia
w folderze Sent skrzynki Operatora odpowiada na pytanie, co wyszło z jego
skrzynki (widzi ją w swoim kliencie poczty, obok listów wysłanych
własnoręcznie); wiersz list_wyslany w bazie rdzenia odpowiada na pytanie,
co wysłała platforma i czy się udało (zapisuje go adapter rdzenia
adapter_modul_poczta_wysylka.go, także dla wysyłki nieudanej).

List leżący w wysłanych, który nigdy nie wyszedł, byłby dowodem czynności,
której nie było. Odwrotnie — nieudane odłożenie kopii po udanym nadaniu
nie unieważnia nadania i nie jest odmową komendy: listu i tak nie da się
cofnąć, więc odpowiedź musi mówić, że wyszedł.

SMTP jedzie biblioteką standardową, tym samym powodem, co odbiór IMAP-em.

Identyfikator wysyłki jest jedyną wartością, po której da się ten list
rozpoznać u odbiorcy i w folderze wysłanych.

Kopia w wysłanych nosi znacznik Seen, bo listu, który Operator sam
wysłał, nie czyta się jako nowego.

Uwierzytelnienie SMTP jest warunkowe. Serwer wysyłkowy dostawcy zawsze go
żąda, ale serwer stojący na tej samej maszynie (albo przekaźnik w sieci
Operatora) często nie ogłasza AUTH wcale. Wpychanie mu wtedy poświadczenia
kończy się odmową serwera przy komendzie, która bez AUTH przeszłaby.

Adres net.JoinHostPort, a nie sklejenie z dwukropkiem: adres IPv6 sam
niesie dwukropki, więc "::1:587" byłoby adresem, którego nie da się
rozebrać.

STARTTLS wymuszany na serwerze, który go nie ma, zerwałby rozmowę zamiast
zabezpieczyć ją mocniej.

Po domknięciu strumienia serwer bierze list na siebie — od tej chwili nie
da się go cofnąć.

PLAIN posyła hasło wprost i biblioteka standardowa dopuszcza go wyłącznie
po TLS-ie (albo na pętli zwrotnej). Dzięki kolejności CRAM-MD5 przed
PLAIN hasło Operatora nie wyjdzie nieszyfrowanym łączem.

## budowa/server/internal/mowa/ustawienia.go

Prawdą o nazwie ustawienia jest wiersz definicja_ustawienia.klucz
z server/internal/store/migracja_075_mowa.sql; stałe pliku są jego kopią
co do znaku. Katalog nie odmawia nieznanego klucza — zapis pod zmyśloną
nazwą przechodzi bez błędu, a kontrolka nie ma wtedy żadnego skutku, więc
jedynym zabezpieczeniem jest to, że napisy stoją po tej stronie w jednym
miejscu i pochodzą z migracji.

Odczyt konfiguracji robi warstwa wyżej (rezolwer ustawień rdzenia, komendy
config.get); tutaj przychodzą już gotowe pary klucz-wartość i jedynym
zadaniem jest złożyć z nich komplet. Dzięki temu silnik mowy nie zna ani
bazy, ani zasięgów, ani osi.

Wartość spoza wykazu nie jest odmową. Klucz nieznany jest pomijany bez
błędu, bo konfiguracja poziomu, z którego przychodzi, niesie także
ustawienia zupełnie innych warstw produktu. Rozmiar modelu spoza wykazu
wraca na wartość domyślną, bo silnik musi dostać nazwę modelu, którą
biblioteka zna.

Trzy pierwsze wartości domyślne odpowiadają migracja_075_mowa.sql, katalog
wag odpowiada migracja_403_nastawa_wag_mowy.sql, która nadpisuje wartość
założoną migracją 075. Rozjazd między tymi stałymi a migracją oznaczałby
dwie różne prawdy o tym, co zobaczy Operator, który niczego nie ustawił.

Katalog modeli domyślny jest w układzie pamięci podręcznej Huba
(models--Systran--faster-whisper-*), bo pomocnik przekazuje tę ścieżkę
bibliotece jako download_root i takiego układu w niej szuka
(pomocniki/transkrypcja/silnik.py). Wagi stoją: 464 MB
w /opt/danaco-modele/mowa, obok wag pozostałych zdolności liczących
lokalnie.

Katalog, w którym wag nie ma, pomocnik traktuje jak miejsce pobrania,
czyli zachowuje się dokładnie tak jak przy wartości pustej — wskazanie
ścieżki nieistniejącej nie jest odmową.

Ustawienia jest strukturą, nie mapą: pola nazwane rozstrzygają się przy
kompilacji, mapa napisów rozstrzygałaby się dopiero przy uruchomieniu —
a literówka w kluczu jest tu dokładnie tym błędem, przed którym plik ma
chronić.

Rozmiary modelu mowy stoją w kolejności rosnącej wierności — tej samej,
w której stoją wiersze opcja_ustawienia w migracja_075_mowa.sql (kolumna
kolejnosc). Inaczej niż przy naklad_rozumowania, gdzie pusty napis znaczy
rozstrzyga kanał modelu, tu wartości pustej nie ma.

Konfiguracja zasięgu, z którego przychodzą wpisy Nanies, niesie ustawienia
całego produktu — kanał modelu, izolację, motyw. Odmowa na każdy nieswój
klucz zamieniłaby zwykły odczyt konfiguracji w awarię silnika mowy.

Komplet wraca przez wartość w Nanies, nie przez wskaźnik: wołający składa
go pętlą po wpisach i nie ma powodu, by dwie takie pętle biegnące obok
siebie pisały po tej samej strukturze.

Sprowadzenie rozmiaru modelu do wartości domyślnej, a nie odmowa: rozmiar
nieznany bibliotece kończyłby się błędem procesu Pythona — komunikatem
cudzego narzędzia, z którego Operator nie wyczyta, że winna jest jedna
literówka w ustawieniu. Sprowadzenie do wartości domyślnej daje
transkrypcję zamiast odmowy, a wykaz dopuszczalnych wartości Operator
i tak widzi w oknie konfiguracji.

Pustka w polu Jezyk znaczyłaby co innego niż brak wpisu (rozpoznanie
automatyczne zamiast polskiego), więc każdy kształt wartości konfiguracji
niesie w tekstUstawienia własną wartość zastępczą.

## budowa/server/internal/wiedza/przesiew_test.go

Samo liczenie ocen sprawdza się tam, gdzie stoją wagi: internal/core,
sprawdziany skutku rodziny knowledge.*. Tutaj mierzone jest to, co ma
działać na każdej maszynie — bo kolejność ułożona źle jest usterką
niezależną od tego, czy model odpowiedział dobrze.

Ocena zerowa trzyma pozycję w wykazie po to, żeby wiązanie po pozycji się
nie przesunęło — a nie po to, żeby stanąć w wyniku.

## budowa/server/internal/poczta/list.go

Rozdzielenie rozbioru i złożenia MIME na dwa pliki dawałoby dwie prawdy
o tym, co rdzeń uważa za treść listu.

List składa się raz, wysyła dwiema drogami. Szkic zapisany w folderze
Drafts i list nadany SMTP-em to dokładnie ten sam dokument — dlatego
zlozWychodzacy ma dwóch wołających (imap_zapis.go i smtp.go), a nie dwie
kopie. Gdyby szkic składał się inaczej niż wysyłka, Operator oglądałby
przed wysłaniem list inny niż ten, który wyjdzie w świat.

Treść bierze się z text/plain, a text/html dopiero przy jego braku. Model
ma analizować treść, a nie znaczniki; list wyłącznie HTML-owy oddaje się
takim, jaki jest — obcięcie go do treści wymagałoby renderowania HTML,
czego rdzeń nie robi i o czym nie kłamie.

Rejestracja dekoderów stron kodowych działa skutkiem ubocznym importu
pustego: pakiet podstawia message.CharsetReader. Bez niego polski list
w iso-8859-2 albo windows-1250 — a takich w skrzynce Operatora jest pełno
— rozbiera się błędem "nieznana strona kodowa" zamiast oddać treść.

Wątek wiąże się identyfikatorem listu, na który ten odpowiada, a przy jego
braku — własnym. Dzięki temu list otwierający wątek i wszystkie odpowiedzi
w nim mają tę samą wartość, bez pytania serwera o THREAD, którego część
serwerów nie zna.

Odkodowanie transportowe zostaje nietknięte w zapowiedzi świadomie: część
pobrana częściowo bywa urwana w połowie czwórki base64 albo w połowie
sekwencji =XX, więc dekoder i tak nie miałby czego domknąć. Zapowiedź
listu w quoted-printable wygląda przez to nieco surowo (=C5=BC zamiast ż)
— i tak jest uczciwiej niż zgadywać brakujące bajty.

Załącznik bez nazwy istnieje i ma bajty — nazwa zastępcza jest tu etykietą,
a nie zmyśloną własnością listu: bez niej magazyn rdzenia nie miałby czym
opisać zasobu.

Klient poczty Operatora układa rozmowę po nagłówku References — bez niego
odpowiedź wyląduje w jego skrzynce obok wątku, a nie w nim.

application/octet-stream znaczy nie wiem, co to jest, i tak jest uczciwie:
zgadnięty application/pdf przy pliku, który PDF-em nie jest, wprowadzałby
w błąd odbiorcę listu.

Skutek podwójnych nawiasów w identyfikatorze wątku widać dopiero
u odbiorcy: klient poczty nie dopasowuje takiej odpowiedzi do wątku, więc
odpowiedź Operatora ląduje obok rozmowy zamiast w niej, a rdzeń melduje
wysyłkę udaną. Obie postaci są w obiegu naraz: koperta IMAP oddaje
identyfikator goły (Naglowek.Watek), a nagłówek listu i człowiek piszą go
w nawiasach — bezNawiasowKatowych przyjmuje jedno i drugie, zamiast
wymagać właściwej postaci od wołającego.

## budowa/server/internal/wiedza/silnik.go

Osadzenia liczy się w dwóch chwilach — przy budowaniu wskaźnika
(knowledge.index) i przy zapytaniu (knowledge.search) — i musi je liczyć
ten sam model tym samym sposobem. Wektor dokumentu policzony jednym
modelem, a wektor pytania drugim, dają iloczyn skalarny, który jest liczbą
i nawet wygląda sensownie, a nie znaczy nic.

Proces silnika startuje wyłącznie przez zewnetrzne.Wolaj: w całym
produkcie stoi dokładnie jedno exec.Command, a każde uruchomienie idzie tą
samą bramą izolacji okna i tym samym obejmowaniem potomstwa. Proces
Pythona liczący na wielu wątkach bez objęcia drzewem zostawiałby sieroty
po każdym przekroczeniu czasu.

Każde wołanie silnika wczytuje model od nowa: prawie cały czas zlecenia to
start procesu i wczytanie wag (rząd gigabajta) do pamięci, a nie samo
porównanie wektorów. Proces rezydentny skróciłby zapytanie, ale wymaga
dwukierunkowej rozmowy z procesem żyjącym między żądaniami, a
zewnetrzne.Wolaj prowadzi rozmowę jednorazową i jest jedyną dozwoloną
drogą startu procesu.

Granica czasu silnika jest dwojaka. Pierwsze uruchomienie pobiera wagi
modelu, więc granica budowania wskaźnika jest liczona w minutach;
zapytanie ma wagi już na dysku i granica jest liczona w sekundach. Jedna
wspólna granica byłaby albo za krótka na pobranie, albo tak długa, że
zawieszone zapytanie wyglądałoby na pracujące.

Dwie drogi do tej samej wartości ustawień silnika byłyby dwiema prawdami,
dlatego ustawienia trzymane są w silniku, a nie odczytywane przy każdym
zleceniu.

Sprawdzenie Gotowy idzie przed budowaniem wskaźnika i przed zapytaniem, bo
odmowa "nie ma czym" jest dla Operatora czymś innym niż "liczyło i się
wywróciło". Pomocnik przygotowuje model i wraca, więc przy pierwszym razie
pobierze też wagi.

Osadzanie partiami, nie wszystko naraz: wykaz kilkudziesięciu tysięcy
fragmentów w jednym pliku zlecenia zająłby pomocnikowi pamięć
proporcjonalną do całej biblioteki. Kolejność wektorów odpowiada
kolejności tekstów i to jest warunek — wołający wiąże je pozycją, nie
treścią. Wykaz pusty do osadzenia nie jest pytaniem o gotowość silnika,
tylko pracą, której nie ma (od pytania jest Gotowy).

`-X utf8` w wołaniu silnika idzie zawsze: pomocnik oddaje polski tekst
w odpowiedzi, a Python bez tego przełącznika koduje wyjście według
ustawień regionalnych systemu — na polskim Windowsie stroną 1250, w której
JSON rozpada się na krzaki.

## budowa/server/internal/wiedza/pomocnik.go

Silnik mowy szuka swojego pomocnika w katalogu pomocniki/ obok binarium
i płaci za to odmową na każdym wdrożeniu, w którym ktoś przeniósł samo
binarium. Skrypt pomocnika osadzeń jedzie w binarium przez go:embed — tak
samo jak migracje schematu (store/zrodlo_migracji.go) — żeby wdrożenie nie
zależało od obecności plików obok programu.

Pomocnik wykłada się do katalogu danych, nie do katalogu tymczasowego. Ten
sam katalog niesie bazę, sejf poświadczeń i magazyn treści biblioteki —
pomocnik przeżywa restart rdzenia tak jak one, a Operator ma go gdzie
obejrzeć, zanim pozwoli mu ruszyć. Zapis jest atomowy (plik tymczasowy
w katalogu docelowym, potem przemianowanie), bo skrypt obcięty w połowie
wystartowałby i wywrócił się komunikatem o składni, którego nikt nie
powiąże z przerwanym zapisem.

Start procesu należy wyłącznie do zewnetrzne.Wolaj. exec.LookPath nie
uruchamia procesu, a jedynie przegląda ścieżkę wyszukiwania systemu.

Ten plik niczego nie uruchamia poza exec.LookPath, które nie uruchamia
procesu, a jedynie przegląda ścieżkę wyszukiwania systemu.

Każdy pomocnik leży pod własną nazwą w tym samym katalogu: jeden plik
o zmiennej treści nie dałby się obejrzeć przed uruchomieniem, a właśnie po
to katalog danych jest miejscem wyłożenia.

Skrypt pomocnika jest bytem wtórnym wobec binarium, więc po podmianie
rdzenia na nowsze wydanie na dysku ma leżeć wersja z tego binarium. Koszt
to zapis kilku kilobajtów raz na żądanie indeksowania.

Druga droga wykładania pomocnika byłaby drugą prawdą o tym, gdzie Operator
ma szukać kodu, który rdzeń uruchamia na jego maszynie.

Droga wołania procesu w zewnetrzne/wolanie.go nie podaje mu standardowego
wejścia, dlatego zapiszZlecenie odkłada treść w pliku i oddaje funkcję
sprzątającą — wołający kasuje plik defer-em, także na ścieżce błędu, żeby
katalog zleceń nie rósł o jeden plik na każde żądanie.

Wskazanie interpretera Operatora idzie wprost i bez sprawdzania na dysku,
bo może być nazwą do rozwinięcia przez system albo dowiązaniem środowiska
wirtualnego, a odmowa na podstawie własnego sprawdzenia unieważniałaby
ustawienie. Gdy wskazania nie ma, szuka się python3, a gdy i tego nie ma —
zostaje python; ta ostatnia wartość jest zgadywana i odmowa przyjdzie
dopiero z uruchomienia, bo tylko ono zna prawdę o wykonywalności.

Brak wskazania katalogu wag znaczy podkatalog katalogu danych rdzenia,
a nie katalog domyślny biblioteki. Różnica jest istotna: biblioteka
domyślnie pisze do katalogu pamięci podręcznej użytkownika, a droga
wołania procesu nie dziedziczy środowiska, więc pomocnik nie zna nawet
HOME i wagi wylądowałyby w miejscu zależnym od tego, gdzie akurat stoi
rdzeń.

Biblioteka wykłada wagi modelu stojącego wprost w katalogu wskazanym, więc
dwa modele w jednym katalogu byłyby dwoma plikami model.safetensors
w tym samym miejscu — czyli jednym z nich nadpisanym przez drugi.

## budowa/server/internal/mowa/zamiar.go

Rozstrzygnięcie nie jest zgadywaniem po treści. Rdzeń nie ma prawa domyślać się
z wolnego tekstu, że Operator chciał wykonać komendę platformy. Plik tego nie
łamie, bo rozstrzyga wyłącznie na deklaracjach, nigdy na podobieństwie:

1. Wskazanie Operatora — pole żądania kontraktu (`intent`). Deklaracja wprost;
   nie ma czego zgadywać.
2. Leksykon akcji — zamknięty wykaz fraz, z których każda jest przypisana
   jednej komendzie kontraktu. Dopasowanie idzie po całej wypowiedzi, znak
   w znak po normalizacji, a nie po fragmencie, tak jak robi to funkcja
   dopasujAkcje.
3. Domyślnie — model. Wolny tekst bez deklaracji jedzie do modelu, więc to
   rozstrzygnięcie niczego nie zmienia za plecami Operatora.

Gdzie deklaracji brak, a domyślnej drogi wziąć nie wolno (Operator zadeklarował
akcję platformy, lecz nie nazwał żadnej), wynikiem jest zamiar nierozstrzygnięty
z nazwanym powodem. Nierozstrzygnięcie jest tu wynikiem pełnoprawnym: lepsze od
wybrania komendy za Operatora.

Plik nie woła komend, nie zna rdzenia i — tak jak reszta pakietu — nie zna
kontraktu: mówi własnymi typami, żeby dało się go wpiąć niezależnie od tego,
kiedy pole `intent` wejdzie do `shared/contract.json`.

### Zamiar, Trzy wartości pola wskazania, AkcjaPlatformy

Wartość zerowa typu Zamiar jest nierozstrzygnięciem i to jest wybór rozmyślny:
struktura wynikowa, która powstała przez pomyłkę zamiast przez rozstrzygnięcie,
ma znaczyć „nikt niczego nie zlecił", a nie „jedź do modelu" ani tym bardziej
„wykonaj komendę". Najbezpieczniejszy wynik ma być najtańszy do przypadkowego
wytworzenia.

Dzięki dosłownym, niewyliczanym nastawom leksykon akcji platformy da się
przeczytać jak umowę i wskazać palcem, co która wypowiedź robi.

### Rozstrzygnij

Kolejność reguł nie jest dowolna:

    tekst pusty                      -> nierozstrzygniety (nie ma czego rozstrzygac)
    wskazanie nieznane               -> nierozstrzygniety (napis spoza slownika)
    wskazanie = model                -> MODEL, chocby fraza pasowala do akcji
    wskazanie = platform + fraza     -> PLATFORMA (komenda z leksykonu)
    wskazanie = platform, brak frazy -> nierozstrzygniety (ktorej komendy?)
    brak wskazania + fraza           -> PLATFORMA
    brak wskazania, brak frazy       -> MODEL (droga dotychczasowa)

Odwrotne pierwszeństwo deklaracji Operatora i leksykonu znaczyłoby, że pole
kontraktu można przegłosować leksykonem, czyli że deklaracja nie jest
deklaracją.

Leksykon pusty (`akcje == nil`) jest poprawnym wejściem, nie brakiem: rdzeń bez
wpiętego leksykonu rozstrzyga wtedy wszystko na model, dokładnie tak, jak
zachowuje się bez leksykonu w ogóle. Podstawienie tu leksykonu domyślnego
z własnej inicjatywy dałoby zachowanie, którego wołający nie zamówił; kto chce
domyślnego leksykonu, woła funkcję AkcjeDomyslne u siebie i widzi to w swoim
kodzie.

### BrakDanychAkcji, NieAkcjaPlatformy, ZlozWywolanie, dopasujAkcje, znormalizuj

BrakDanychAkcji jest osobnym typem z tego samego powodu, dla którego pakiet ma
trzy odmowy w bledy.go: komunikat jest trójczęściowy jak każdy w tym pakiecie.

Puste `Wywolanie{}` bez błędu byłoby atrapą wyglądającą na wynik — stąd
NieAkcjaPlatformy jako osobna odmowa.

ZlozWywolanie niesie trzy odmowy, każda inna: rozstrzygnięcie, które nie jest
akcją platformy, kończy się odmową, bo wołający, który mimo tego złożyłby
wywołanie, wykonałby czynność niezleconą; brak pola z wykazu Wymaga kończy się
odmową BrakDanychAkcji z nazwą pola; pole puste liczy się jak brak, bo puste
„windowId" pojechałoby do rdzenia jako żądanie bez okna i wróciło odmową gorzej
opisaną niż odmowa tutaj. Gdyby kontekst mógł podmienić `control`, akcja
„wstrzymaj" dałaby się w locie zamienić w „anuluj" i leksykon przestałby być
umową.

dopasujAkcje szuka wypowiedzi równej frazie, nie zawierającej frazy: dopasowanie
po fragmencie zamieniłoby wzmiankę w wykonanie — „przypomnij mi, co robi anuluj
zlecenie" anulowałoby zlecenie, choć Operator prosił o wyjaśnienie. Wypowiedź
równa frazie jest natomiast deklaracją samą w sobie — nie da się jej powiedzieć
przypadkiem. Rozstrzyganie sprzeczności między dopasowaniami należy do funkcji
Rozstrzygnij, a nie do wyszukiwania.

znormalizuj zdejmuje wyłącznie to, czego mowa nie niesie — wielkość liter
i interpunkcję, które dokłada pomocnik transkrypcji — i nie skraca treści.

### AkcjeDomyslne

Wszystkie pięć fraz prowadzi do jednej komendy assistant.action.status i nie
jest to niedoróbka, tylko granica postawiona z dwóch własności kontraktu:

1. Komenda ta ma wszystkie pola żądania nieobowiązkowe poza jednym: wystarcza
   identyfikator okna, jedyny identyfikator, który komenda głosowa naprawdę
   niesie. Akcje w rodzaju window.create (wymaga identyfikatora sesji, modułu,
   kanału modelu, katalogów roboczych, środowiska wykonania, trybu uprawnień,
   roli okna) albo session.stop (wymaga identyfikatora sesji) musiałyby te
   pola skądś wziąć, a „skądś" znaczyłoby „z głowy rdzenia".
2. Sterowanie własnym zleceniem jest tą czynnością, której droga przez model
   jest wprost szkodliwa: „anuluj zlecenie" powiedziane do modelu ląduje jako
   treść tury tego właśnie zlecenia, więc zamiast je przerwać — przedłuża.

Wykaz jest startowy i wołający ma prawo podać własny. To, które frazy platforma
rozumie, nie należy do tego pliku; wnosi on mechanizm i pięć fraz, przy których
mechanizm daje się sprawdzić na żywym rdzeniu. Wspólna kopia wyniku
pozwoliłaby wołającemu zmienić leksykon wszystkim naraz, nie chcąc tego.

### Wywolanie

Pola typu Wywolanie są napisami, bo napisami przychodzą i z leksykonu akcji,
i z żądania głosowego; przełożenie ich na typ kontraktu należy do wołającego,
który kontrakt zna.

## budowa/server/internal/wiedza/obraz.go

Osadzarka tekstu tu nie wystarczy i nie chodzi o jakość, tylko o przestrzeń.
Wektor zdania z modelu tekstowego leży w przestrzeni, w której obrazu nie ma;
model osi obrazu ma dwie wieże, jedną dla pikseli, drugą dla słów, uczone tak,
żeby kończyły w jednej przestrzeni. Dopiero tam iloczyn skalarny znaczy „to
zdanie opisuje ten obraz".

Wektorów obrazów nie ma we wskaźniku i to jest rozstrzygnięcie, nie brak.
Tabela fragment_wiedzy trzyma przy każdym wektorze fragment tekstu, z którego
go policzono (migracja store/migracja_115_wskaznik_znaczenia.sql), a obraz
takiego fragmentu nie ma; wpisanie tam nazwy pliku dałoby kolumnę, która dla
jednych wierszy jest cytatem, a dla innych etykietą. Kolumna zakres ma ponadto
warunek dopuszczający trzy wartości kontraktu i czwartej nie przyjmie bez
migracji, a migracje nastaw prowadzi inny teren. Dlatego oś obrazu liczy
wektory na każde zapytanie: koszt to jedno wczytanie wag i jeden przebieg
wieży obrazu na plik, sekundy przy bibliotece rzędu setek obrazów. Trwałość
tych wektorów jest pracą do zrobienia, nie założeniem tego pliku.

Silnik nie zna kontraktu ani biblioteki: dostaje zdanie i ścieżki plików,
oddaje po jednej ocenie na ścieżkę. Wybór obrazów i przekład na trafienia
należą do adaptera rdzenia, tak samo jak przy osadzarce.

### PominietyObraz

Bez wykazu obrazów pominiętych Operator patrzący na wynik bez swojego zdjęcia
nie ma jak odróżnić przypadku, w którym model obrazu go nie wybrał, od
przypadku, w którym plik jest uszkodzony.

### Dopasuj

Obraz pominięty ma ocenę zerową i zajmuje swoje miejsce w wykazie, żeby
pozycje wyników się nie przesunęły.

### NajblizszeObrazy

Sortowanie jest stabilne: powtarzalność odpowiedzi jest warunkiem, nie
wygodą.

### Sciezka (pole Obraz)

Ścieżka do bajtów obrazu nie wychodzi z rdzenia poza jego granicę: odpowiedź
komendy niesie tylko nazwę i identyfikator obrazu, nigdy układ dysku
Operatora.
