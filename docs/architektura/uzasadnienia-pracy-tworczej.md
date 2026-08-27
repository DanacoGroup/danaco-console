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
