# Danaco Console — Uzasadnienia komentarzy klienta poprzedniego

Dokument gromadzi uzasadnienia, które przekraczają dopuszczalną długość
nagłówka komentarza w plikach `budowa/klient-poprzedni/` i `budowa/klient/`.
Każdy rozdział nosi nazwę pliku źródłowego, którego uzasadnienie dotyczy.

## budowa/klient-poprzedni/src/moduly/studio/okno-pracy-z-dokumentem.ts

Okno łączy w jednej powierzchni to, co dawniej działało jako cztery osobne
okna nad tą samą treścią: edycję z formatowaniem, podgląd wydania i różnicę
wersji jako tryby jednego widoku, przy czym różnica jest warstwą nakładaną na
treść, a nie drugą kolumną obok niej. Wynik operacji wchodzi do dokumentu jako
zmiana śledzona autora `model`, a jej poprawianie przed decyzją odbywa się tam,
gdzie zmiana stoi — w treści, nie w osobnym buforze.

Tryb wykazu operacji (narzędzia ukryte albo stały panel) jest jedną nastawą na
cały moduł, więc rozstrzyga ją i zapisuje okno pracy, a nie pływak i stały
panel osobno — inaczej te dwa miejsca rozjechałyby się w tym, co pokazują.

Siedem źródeł postaci dokumentu, kontroli pracy i wstawień (rodziny poleceń
`studio.page.*`, `style.*`, `format.*`, `table.*`, `object.*`, `apparatus.*`,
`field.*`, `journal.*`, `markup.*`, `lock.*`, `backup.*`, `template.*`) nie
należą do portu pracy, dlatego wchodzą osobnym bytem `ZrodlaPostaciStudia`,
a nie dopiskiem do `ZrodloPracyStudio`. Byt jest nieobowiązkowy: gdy któregoś
źródła nie podano, okno działa dalej, tylko odpowiednia rodzina poleceń jest
wyłączona — brak leży wtedy po stronie złożenia modułu, nie rdzenia.

Górna granica wielkości jednej kartki (sześć megabajtów) wynika z tego, że
kartka A4 w 96 punktach na cal jako obraz PNG mieści się w rzędzie setek
kilobajtów, więc sześć megabajtów daje zapas także przy nośniku
wielkoformatowym; rdzeń dostaje tę granicę, bo kontrakt każe mu wtedy odmówić
z podaniem zmierzonej wielkości, zamiast oddać treść uciętą, co wyglądałoby na
kartkę zepsutą przez rdzeń. Górna liczba kartek pobieranych jednym renderem
(dwadzieścia cztery) wynika z tego, że pismo dwustustronicowe dałoby dwieście
żądań i dwieście obrazów w pamięci karty na raz, a podgląd wydania służy
sprawdzeniu składu, nie czytaniu całego pisma obrazkami.

Postać dokumentu musi powstać przed powierzchnią, bo powierzchnia bierze od
niej magazyn nastaw widoku i zgłasza jej chwyty linijki; odwrotna kolejność
zostawiłaby chwyty bez odbiorcy. Kartki wyrysowane przez rdzeń tracą
aktualność w chwili, gdy Operator zmieni treść, dlatego są zdejmowane od razu:
pokazywanie starego wyrysu jako podglądu treści bieżącej byłoby nieprawdą
o dokumencie, nie oszczędnością jednego wywołania.

Mikrofon wiersza polecenia dostaje z powłoki okno komunikacji, nie kartę
sesji, bo `speech.audio.upload` przyjmuje puste `sessionId` jako nagranie bez
przypisania do karty — zgadywanie karty z okna byłoby wskazaniem nieprawdziwym.
Licznik użycia operacji jedzie do rdzenia po każdym uruchomieniu, bo z niego
bierze się kolejność czynności na wierzchu pływaka: z użycia, nie z domysłu.

Zakładka założona przyciskiem przybornika jest elementem aparatu dokumentu
i zapisuje się w rdzeniu, więc przeżywa zapis; wykaz w oknie zostaje obok, bo
powrót ma działać od ręki, bez odczytu, ale prawdą o dokumencie jest wiersz
w rdzeniu, nie ten wykaz. Trzeci argument przybornika znakowania jest rdzeniem
znakowania: bez niego przybornik pokazuje znakowania samej sesji, z nim
znakowanie jest trwałe i przechodzi przez zapis dokumentu.

Dziennik, kopie i blokady stoją nakładkami przy powierzchni, bo dotyczą tej
samej treści: dziennik cofa jej czynności, kopie ją ratują, blokady jej
pilnują — każdy przyjmuje skutek tą samą drogą co reszta okna, treścią roboczą
stanu. Postać po czynności czyta się z rdzenia, a nie przyjmuje nieopisanym
ładunkiem z odpowiedzi: postać jest drzewem o kształcie kontraktu, a wartość
niesprawdzonego typu wstawiona w stan okna byłaby drugą, niepewną prawdą
o dokumencie.

Wskaźnik zakresu operacji rozróżnia trzeci przypadek — cały dokument z
zaznaczeniem pominiętym wyborem ręcznym z Tools Panelu — osobno od pozostałych
dwóch, bo bez niego wybór ręczny wyglądałby jak usterka: zaznaczenie widać,
a operacja idzie na całość.

Panel osadzenia źródeł stoi w powłoce treści, nie na całej powierzchni okna,
żeby wstążka i pasek statusu zostały widoczne i Operator nie stracił drogi
powrotu. Pas otwarć trzech paneli (schowek, źródła, przybornik) stoi osobno od
wstążki, bo wstążka niesie formatowanie i widok, a te trzy są narzędziami
bocznymi treści; siedem powierzchni postaci i kontroli pracy w dalszej części
okna idzie tym samym wzorem — przycisk w pasie, treść nakładką, panel
zamknięty nie zajmuje miejsca — jako reguła stopniowego ujawniania, nie
oszczędność powierzchni.

Przerysowanie marginesu dzieje się jednym wywołaniem na trzy byty naraz
(komentarz, propozycja, zmiana śledzona), bo stoją w jednej kolumnie i muszą
się ustawić wobec siebie — dwa osobne przerysowania kładłyby je jedne na
drugich.

Zasoby stron w odpowiedzi renderu stoją na początku wykazu, a dokument
w formacie docelowym dochodzi na jego końcu, więc kartek jest dokładnie tyle,
ile podaje pole `pages`, i do tej liczby wykaz się przycina; rozpoznawanie
kartki po typie treści wymagałoby pobrania także dokumentu, czyli megabajtów
pobranych po to, żeby je odrzucić.

Wstawianie treści w miejsce kursora jest jedną drogą wspólną dla schowka,
dyktowania i wniesienia ze źródła: zaznaczenie zostaje zastąpione tak jak
w każdym edytorze, a brak zaznaczenia znaczy „na końcu treści", bo kursor bez
zaznaczenia stoi tam, gdzie stan modułu ostatnio go widział.

Malarz formatów kopiuje postać akapitu, nie treść, więc nie jedzie schowkiem
rdzenia: rodzina poleceń `clipboard.*` niesie treść i rodzaj wpisu, a nie
arkusz nastaw akapitu, dlatego postać czyta się z nastaw akapitu okna.

## budowa/klient-poprzedni/src/moduly/studio/powierzchnia-dokumentu.ts

Powierzchnia niesie zaznaczenie do rdzenia jako zakres znaków treści w zapisie
markdown, bo tak liczy je komenda kontekstowa operacji, a kursor po
przerysowaniu wraca na miejsce w tekście widocznym bloku, bo tam Operator go
widział; dwie miary służą dwóm różnym odbiorcom i zlanie ich przesuwałoby
kursor o długość znaczników markdown. Podział na kartki jest prawdziwy —
liczą go wysokości bloków zmierzone przez przeglądarkę, nie licznik znaków —
bo blok, który nie mieści się na kartce, ma zejść na następną tak, jak zejdzie
w wydaniu. Rdzeń liczy własną paginację osobno, silnikiem swojego składu,
i oddaje ją obrazami stron; podgląd wydania sięga po te obrazy dopiero, gdy
rdzeń je odda, a do tego czasu pokazuje kartki liczone tutaj i mówi to wprost
paskiem stanu, żeby wyrys nieaktualny nie uchodził za wynik ostateczny.

Kartka rdzenia wchodzi obrazem w zapisie `data:`, bo pole `uri` zasobu
magazynu jest ścieżką w systemie plików rdzenia, niedostępną przeglądarce;
bajty przychodzą odpowiedzią komendy pobrania treści zasobu i stąd zapis
`data:`, bez jednego żądania do sieci. Kartki rdzenia tracą aktualność w
chwili zmiany treści, więc wykaz pusty zdejmuje wyrys rdzenia i wraca do
kartek liczonych w oknie — pokazywanie starego wyrysu jako podglądu treści
bieżącej byłoby nieprawdą o dokumencie.

Magazyn nastaw widoku rozróżnia argument pominięty od podanego jako `null`:
pominięty daje magazyn tego urządzenia, `null` znaczy „bez zapisu" — nastawy
działają przez sesję i nikt nie obiecuje, że przeżyją zamknięcie. Nastawy
strony przestawione chwytem linijki albo paskiem widoku idą do osobnego
wykazu nadpisań, nie tylko do nastaw bieżących, bo okno pcha własną kopię
całych nastaw strony przy każdej swojej zmianie; bez wykazu nadpisań ta kopia
zabierałaby Operatorowi to, co właśnie ustawił chwytem. Gdy okno pcha nowe
nastawy, wygrywa wyłącznie to pole, które okno naprawdę zmieniło — nadpisanie
zostaje, dopóki wartość okna się od niego nie różni.

Nastawa gotowa skali liczy się z pola widoku zmierzonego w chwili wywołania,
bo tylko wymiar zmierzony teraz jest prawdziwy; nastawa nieznana albo pole
jeszcze nieosadzone nie zmienia nic, bo skala zgadnięta byłaby gorsza od
nastawy niewykonanej. Wydruk bierze te kartki, które Operator widzi w
podglądzie, nie drugi wyrys ani surowy tekst źródłowy, więc przed drukiem
powierzchnia schodzi do trybu wydania, oznacza strony poza zakresem druku i
zdejmuje adiustację, gdy Operator drukuje pismo do wysłania — wszystkie trzy
kroki wraca po zamknięciu okna drukarki.

Decyzja o zmianie śledzonej stoi w miejscu zmiany, w treści, a nie w osobnym
oknie, bo Operator rozstrzyga tam, gdzie patrzy; pasek zatwierdzenia pod
blokiem używa tej samej drogi decyzji co znaczniki w treści i wykaz w
okienku recenzowania, nie jest trzecim miejscem decyzji. Kotwica komentarza
jest znakiem w treści, nie wpisem w wykazie, bo komentarz przypięty do
fragmentu ma pokazywać, do którego — komentarz bez zakresu kotwicy nie
dostaje, bo dotyczy dokumentu, a nie miejsca w nim. Panele okna stoją nad
powierzchnią jako nakładka, bo powierzchnia należy do dokumentu — stałe
kolumny zostają trybem do wyboru, a znacznik układu idzie na wspólnego
przodka powierzchni i paneli, ponieważ to on rozkłada kolumny.

Krawędzie kolumn tabeli mierzy powierzchnia z wyrysu, nie zgaduje z liczby
kolumn, bo tabela ustawia szerokości wedle treści komórek, więc krawędź na
linijce ma stać tam, gdzie naprawdę stoi krawędź komórki; krawędź ostatniej
komórki jest krawędzią całej tabeli, nie granicą między kolumnami, więc nie
podlega chwytowi przesunięcia. Położenie punktu w zapisie treści liczy się
sumą zapisów bloków poprzedzających i zapisem bloku bieżącego do miejsca
kliknięcia — tak samo, jak treść dokumentu jedzie do rdzenia — a punkt poza
treścią oddaje `null`, nie zero, bo zero jest położeniem prawdziwym i pomyłka
tutaj wysłałaby operację na początek dokumentu.

## budowa/klient-poprzedni/src/moduly/studio/strona-postac-dokumentu.ts

Warstwa postaci spaja panele okna pracy z komendami rdzenia: panele składają
treść żądania z pól i nic więcej wiedzieć nie muszą, a warstwa dokłada do
każdego żądania identyfikator dokumentu, zakres z bieżącego zaznaczenia oraz
autora czynności, woła rdzeń, czyta jego bilans i mówi Operatorowi, co się
stało — łącznie z tym, czego czynność nie zrobiła i przez którą blokadę.
Zakres dokłada się w warstwie, nie w panelu, bo zaznaczenie jest jedno na
moduł i mieszka w stanie — gdyby panel je czytał sam, powstałoby drugie
miejsce wiedzące, nad czym Operator pracuje, i te dwa rozjechałyby się przy
pierwszym przełączeniu zakładki dokumentu.

Autor czynności jedzie jawnie: czynność Operatora odkłada się jako zmiana
autora `uzytkownik`, czynność modelu jako zmiana autora `model`, inaczej
przełącznik pokazujący pracę modelu nie miałby czego podświetlić. Warstwa
zawsze podaje autorem Operatora, bo to jego panel — droga modelu do tych
samych komend idzie osobno, narzędziami modelu, i pola autora tam nie
zapomina. Brak wczytanego dokumentu jest odmową nazwaną, nie ciszą: każda
z komend przyjmuje identyfikator dokumentu jako pole obowiązkowe, bo żądanie
bez dokumentu nie ma czego dotyczyć, a wysłanie go z polem pustym wróciłoby
błędem walidacji, którego Operator by nie zrozumiał — warstwa mówi więc
wprost, czego brakuje i po czyjej stronie leży brak.

## budowa/klient-poprzedni/src/moduly/research/wywolania-komend.ts

Zdanie odpowiedzi komendy mówi o mierzonym skutku czynności — ile pozycji,
jaki plik, ile luk — a nie o tym, że wywołanie się powiodło: „rdzeń oddał
wynik" jest zdaniem, po którym Operator nadal nie wie, czy coś się stało.
Akcje panelu niosą w polu kodu dokładnie nazwę komendy rdzenia, więc
rozdzielnik komend tego pliku jest odwzorowaniem jeden do jednego i nie ma
w nim ani jednej nazwy pisanej z ręki. Nie ma w nim natomiast żądań, których
nie da się złożyć bez tekstu od Operatora — zapytania wyszukiwania, treści
notatki, uzasadnienia odrzucenia — bo takie pozycje mają mówić wprost, czego
brakuje, zamiast wysyłać żądanie z polem pustym i wracać odmową walidacji,
z której nic nie wynika.

## budowa/klient-poprzedni/src/moduly/terminal/okno-task-schedule.ts

Plan zadania powłoki jest drugą powierzchnią jednej rodziny komend automatyki,
nie drugą rodziną komend: cykliczność wiąże się z automatyką, zapisuje ją
`automation.schedule.set`, czyta `schedule.get`, a rodziny `terminal.schedule.*`
nie ma i nie będzie, bo dwie rodziny na jeden harmonogram byłyby dwiema
prawdami o jednym bycie. Zaplanowanie zadania idzie dwoma krokami tej rodziny:
`automation.workflow.save` zapisuje automatykę o jednym kroku rodzaju
`command` wołającym `terminal.command.exec`, a `automation.schedule.set`
dokłada jej wyrażenie cron; edytora automatyk to okno nie stawia, bo pełna
praca na krokach, wersjach i zmiennych zostaje w module Automations.

Plan nie jest budzikiem: rdzeń wylicza chwilę najbliższego uruchomienia, ale
nie ma czym odpalić automatyki samodzielnie, bo kontrakt nie zna komendy ani
zdarzenia, którym harmonogram zgłaszałby wyzwolenie — zapisany plan jest
zapisem obowiązującym wraz z terminem, a uruchomienie prowadzi Operator
z okien modułu Automations, i potwierdzenie zapisu mówi to wprost. Kolejki
okno nie zakłada, bo `queue.create` żąda identyfikatora karty sesji, którego
rdzeń jeszcze nie wypełnia. Potok jest sekwencjonowaniem po stronie klienta,
nie silnikiem w rdzeniu: okno wysyła krok, czeka na jego domknięcie odczytem
wyjścia i dopiero wtedy decyduje o kroku następnym; zamknięcie okna w trakcie
przerywa sekwencjonowanie, ale nie krok już uruchomiony — ten kończy się
w rdzeniu i widać go w Process Monitorze.

Granica czekania na domknięcie kroku jest dłuższa niż w podglądzie wyjścia
Process Monitora, bo tu czekanie jest warunkiem decyzji: bez kodu wyjścia
kroku nie ma jak rozstrzygnąć, czy wolno ruszyć z krokiem następnym. Krok
automatyki wołający polecenie powłoki niesie nazwę komendy kontraktu i treść
jej żądania, nie polecenie powłoki wprost, dzięki czemu plan przechodzi tę
samą bramę uprawnień i ten sam egzekutor izolacji, co polecenie wydane ręcznie
z karty; krok omijający komendę byłby drugą drogą do powłoki, bez żadnego
z tych sprawdzeń.

## budowa/klient-poprzedni/src/moduly/translate/okno-warsztat-translate.ts

Rodziny komend spoza czterech okien pierwotnych dzielą trzy wskazania — okno
tłumaczenia, panel języka i ścieżkę pliku — więc stoją w jednym warsztacie,
a nie w dwunastu osobnych oknach: rozbicie kazałoby Operatorowi wpisywać te
same trzy wskazania po kilkanaście razy, a każda rodzina prowadziłaby własną
kopię wykazu paneli. Tutaj wskazania stoją raz, na górze, i wchodzą do
każdego żądania. Każda sekcja kończy się zdaniem o mierzonym skutku, nie
słowem „gotowe" — liczbą par wniesionych do pamięci, ścieżką pliku, który
powstał, liczbą pozycji przebiegu — bo komunikat „udało się" bez liczby jest
meldunkiem zamiast skutku.

## budowa/klient-poprzedni/src/moduly/studio/przybornik-znakowania.ts

Znaczenie fragmentu jest czynnością powtarzaną kilkadziesiąt razy na
dokument, a droga przez menu kosztuje przy każdym razie dwa ruchy więcej,
dlatego wszystko, czym się dokument znaczy, stoi w bocznym przyborniku przy
krawędzi treści, nie w menu ani osobnym oknie. Przybornik jest narzędziem do
pracy na fragmentach, nie na dokumencie w całości: gdy zaznaczenia nie ma,
pozycje nie milkną, tylko mówią, że znakowanie obejmie cały dokument, i to
jest odpowiedź, nie odmowa.

Wszystko, czym Operator znaczy dokument, model musi umieć założyć i
przeczytać: komentarz i adnotacja idą komendami rdzenia niosącymi pole
autora, więc znakowanie modelu jest podpisane jako `model`, a Operatora jako
`uzytkownik`. Znacznik własny pola autora w kontrakcie nie ma, więc
przybornik trzyma go u siebie wraz z autorem — to brak nazwany, nie
zatajony. Czynności bez zaplecza w rdzeniu stoją w osobnej części, z
nazwaniem, co dokładnie brakuje — przypis, odsyłacz, odwołanie wzajemne,
wstawienie tabeli i pola należą do aparatu i postaci dokumentu, których
kontrakt jeszcze nie niesie — bo przycisk wychodzący do rdzenia po komendę,
której nie ma, byłby uprzejmą odmową udającą funkcję.

Zaplecze rdzenia znakowania jest nieobowiązkowe, bo przybornik stoi w oknie
pracy z dokumentem, którego montaż powstaje etapami: zaplecze podane znaczy,
że droga do rdzenia jest wpięta i znakowanie jest wtedy trwałe oraz widoczne
modelowi, a zaplecze `null` znaczy, że tej drogi jeszcze nie wpięto, i
przybornik pisze to Operatorowi wprost, zamiast pokazywać przyciski, które
nic nie robią. Przycisk dyktafonu, wzorem pakietu biurowego, jest przyciskiem
dyktowania do treści dokumentu; przybornik go tylko osadza, a nagrywanie
i przepisanie liczy zaplecze mowy, jednym rachunkiem wspólnym z mikrofonem
wiersza polecenia.

## budowa/klient-poprzedni/src/moduly/studio/styl-panel-arkusza.test.ts

Sprawdziany panelu arkusza stylów oraz panelu list i znaków mierzą treść zgłoszonego żądania,
a nie samo naciśnięcie przycisku. Trzy zachowania są sprawdzane osobno jako miejsca dawnej szkody
w tym produkcie. Cecha logiczna niesie trzy stany, nie dwa: naniesienie pogrubienia nie zdejmuje
przy okazji kursywy, której operator nie wskazał. Czynność bez ani jednego wypełnionego pola jest
odmawiana, a nie wysyłana jako żądanie, które niczego nie zmieni. Nastawa poziomu listy wymaga
wskazania listy zastanej — bez niej okno nazywa brakujący element zamiast wysyłać żądanie z pustym
identyfikatorem.

## budowa/klient-poprzedni/src/moduly/studio/szablon-panel.ts

Zapis szablonu z bieżącego dokumentu przenosi domyślnie blokady fragmentów
wzorcowych — pole `includeLocks` bez wartości znaczy „tak", ponieważ fragmenty
wzorcowe pisma mają pozostać wzorcowe także w dokumentach założonych z tego
szablonu.

Szablonu fabrycznego nie da się usunąć: rdzeń odmawia i zwraca powód odmowy
w polu `deleted`. Panel zostawia pozycję w wykazie, ponieważ szablon nadal
istnieje — zdjęcie jej z widoku przy jednoczesnym powrocie po kolejnym odczycie
byłoby pokazaniem skutku, który się nie wydarzył.

## budowa/klient-poprzedni/src/moduly/studio/osadzenie-panel.ts

Panel osadzony pokazuje dokument przez cały czas pracy: Operator nie opuszcza
edytora, a wybór położenia panelu — obok treści albo na całej powierzchni —
należy do Operatora. Wniesienie fragmentu z Biblioteki albo ze strony
sieciowej trafia wprost w miejsce kursora w dokumencie, nie do kolejki ani do
zasobów; komenda `studio.ingest.url` oddaje pole tekstowe, które panel
prowadzi do dokumentu, a gdy tekst jeszcze nie istnieje, bo pozycja czeka na
rozpoznanie pisma, panel nazywa ten stan zamiast wnosić pustkę.

Komendy `studio.insert.from.library` i `studio.insert.from.web` wnoszą
fragment do dokumentu i oddają zapis pochodzenia — adres albo plik, wersję,
czas sięgnięcia i autora. Panel woła te komendy, gdy dostał źródło wstawień,
i pokazuje pochodzenie oddane przez rdzeń wraz z bilansem czynności. Wiersz
pochodzenia wnoszony do treści dokumentu pozostaje drogą zapasową, ponieważ
przeżywa wydanie dokumentu do formatu, który zapisu pochodzenia nie niesie.

Panel nie rysuje strony sieciowej. Migawka `browser.snapshot.get` oddaje
adres, tytuł, treść renderowaną i źródło strony jako tekst, nie jako obraz;
zrzut ekranu jest odnośnikiem zasobu, a komendy pobierającej jego bajty do
przeglądarki kontrakt nie niesie. Panel pokazuje więc treść, którą rdzeń
naprawdę oddaje, i nazywa wprost to, czego nie oddaje. Ramka z cudzą stroną
wewnątrz okna byłaby drugą przeglądarką, a moduł Browser jest jeden.

## budowa/klient-poprzedni/src/moduly/studio/widok-pasek-widoku.ts

Panele nastaw wchodzą na żądanie i schodzą, gdy nie są używane; stała kolumna
zabierałaby kartce szerokość na stałe, dlatego pasek jest wąskim rzędem przy
krawędzi powierzchni, a nastawy rzadsze — układ kartek, przewijanie, jednostka
linijki, tryb dwóch dokumentów — stoją w nakładce rozwijanej przyciskiem
„Nastawy widoku". Nakładka stoi nad treścią i schodzi naciśnięciem, klawiszem
Escape albo naciśnięciem poza nią, nie zabierając kartce ani milimetra
szerokości.

Nastawy widoku stoją przy powierzchni, a nie wyłącznie na wstążce, ponieważ
dotyczą tego, na co Operator patrzy, i sięga po nie stale — skala oraz skok
o stronę są czynnościami ciągłymi, nie wyprawą na osobną zakładkę. Wstążka
niesie te same nastawy jako gniazdo tego samego paska, więc nie powstają dwa
miejsca, które mogłyby się rozjechać.

Szybkie drukowanie stoi w rzędzie stałym, nie w nakładce: jedno naciśnięcie
uruchamia druk ostatnimi nastawami, bez okna nastaw, tak jak pasek szybkiego
dostępu pakietu biurowego. Bez sterownika druku przycisk pozostaje widoczny
i nazywa powód niedostępności zamiast milczeć.

Tryb źródłowy ze znacznikami i podgląd wydruku wykluczają się: znaczniki nie
mają paginacji, a podgląd jej wymaga. Wyłączony przełącznik mówi to wprost,
zamiast oddawać widok, którego Operator nie zamawiał.

## budowa/klient-poprzedni/src/moduly/studio/strona-panel-nastaw.ts

Numeracja stron i marginesy stoją w tym panelu jako nastawy zapisywane osobnymi
komendami rdzenia, nie jako pojedynczy przełącznik profilu wydania — numeracja
niesie styl, umiejscowienie i numer startowy, a marginesy przeżywają zapis
zamiast żyć wyłącznie przez czas trwania sesji.

Panel jest nakładką otwieraną przyciskiem i domyślnie schowaną, ponieważ
powierzchnia należy do dokumentu — stała kolumna nastaw zajmowałaby ją na
stałe, choć Operator sięga po te nastawy rzadko.

Każda nastawa strony da się ustawić osobno dla sekcji, dlatego wybór sekcji
stoi raz, u góry panelu, i dotyczy wszystkich pól poniżej: pismo z załącznikiem
w orientacji poziomej pozostaje wtedy jednym dokumentem, nie dwoma, a Operator
nie wskazuje sekcji osobno przy każdej nastawie.

Panel nie woła rdzenia i nie zna dokumentu, na którym pracuje Operator — składa
treść żądania z wypełnionych pól i oddaje ją wywołującemu. Identyfikator
dokumentu, autor czynności i odczyt odpowiedzi rdzenia należą do warstwy
wyżej, żeby panel nie stał się drugim miejscem, które śledzi, nad czym
Operator pracuje.

## budowa/klient-poprzedni/src/moduly/translate/katalog-funkcji-translate.ts

Wykaz istnieje po to, żeby stan modułu dało się przeczytać, a nie zgadnąć. Okno
pokazuje czynności, które wykonuje; katalog pokazuje komplet zamierzony
w opracowaniu modułu i przy każdej pozycji mówi jedno z dwojga: którą komendą
kontraktu jest wykonywana albo czego brakuje, żeby była. Pozycja bez komendy
nie znika z wykazu, bo zniknięcie byłoby ukryciem braku.

Nazwy pozycji i podział na grupy pochodzą z opracowania modułu i nie są
tłumaczone ani parafrazowane. Numeracji opracowania katalog nie przenosi:
pozycję odnajduje się po pełnej nazwie, nie po oznaczeniu, które poza
dokumentem nic nie znaczy.

Katalog jest zbiorem danych, nie widokiem — wyszukiwarka funkcji buduje z niego
listę, a okna sięgają po pozycje swojej grupy.

## budowa/klient-poprzedni/src/moduly/terminal/okno-session-manager.ts

Do rdzenia idzie jedna komenda otwarcia karty, z powłoką ssh i adresem celu
podanym zmienną środowiska SSH_TARGET — tak i tylko tak rdzeń przekazuje adres
programowi ssh. Karta powstaje w tym samym stanie modułu, w którym stoją karty
otwarte w oknie Terminal Tabs, więc połączenie z tego okna jest od razu
widoczne w oknie wiodącym i w konsoli wyjścia poleceń.

Książka hostów, klucze SSH, tunele portowe i wykaz kart są bytami rdzenia,
a okno jest ich widokiem. Wpis hosta zapisuje i czyta rodzina komend hosta,
klucze prowadzi rodzina komend klucza, tunele — rodzina komend tunelu, a karty
pokazuje wykaz sesji rdzenia, więc wykaz sięga dalej niż pamięć tego
połączenia: rdzeń odtwarza karty przy starcie i okno je widzi po ponownym
podłączeniu gniazda.

Czytanie i pisanie pliku konfiguracyjnego OpenSSH zostaje, bo służy czemu
innemu niż trwałość: wnosi wpisy z maszyny operatora i wynosi je z powrotem.

Pozycja panelu akcji bez odpowiadającej jej komendy kontraktu — polityka
znanych hostów — stoi jawnie nieczynna wraz z powodem liczonym z odczytu
wykazu komend rdzenia, tak samo jak każda inna pozycja zależna od pokrycia
komend.

## budowa/klient-poprzedni/src/uwierzytelnienie/ekran-logowania.ts

Ekran jest przesłoną nad aplikacją, nie osobną trasą: kładzie się nad
gospodarzem dokumentu, a aplikacja pod nim składa się i łączy z rdzeniem
w tym samym czasie, więc po wejściu przesłona znika i widoczne jest gotowe
Centrum dowodzenia zamiast drugiego ładowania. Trasy aplikacji zostają
nietknięte, ponieważ bramka nie jest miejscem pracy.

Ekran nie zakłada budzika, nie liczy czasu sesji i nie przerywa pracy
pytaniem o tożsamość. Rdzeń nie odcina komend po wygaśnięciu sesji — bramka
jest progiem wejścia, nie strażnikiem każdego żądania — więc wygaśnięcie
ujawnia się wyłącznie przy następnym uruchomieniu, zdaniem nad formularzem.

Który z dwóch formularzy pokazać — wejście hasłem czy pierwsze ustawienie
hasła — rozstrzyga rdzeń, nie domysł klienta. Pole nie jest blokowane,
przycisk nie jest wyszarzany; jedyny sprawdzian po stronie formularza to
zgodność hasła z powtórzeniem przy zakładaniu, ponieważ kontrakt opisuje
powtórzenie jako sprawę formularza klienta.

Wymóg logowania jest nastawą `gateway.requireLogin` — gdy jest wyłączony,
przesłona nie staje.

Stan pola „nie wyloguj mnie” bierze się z miejsca, w którym leży sesja
zapisana poprzednio, a nie ze stałej: raz odznaczone pole ma zostać
odznaczone. Zaznaczone pole kładzie sesję w pamięci trwałej, odznaczone
w pamięci okna, a ta sama wartość idzie do rdzenia i rozstrzyga o trwaniu
sesji.

Odnośnik resetu hasła prowadzi do odzyskania konta, nie do zmiany hasła ze
znanym hasłem bieżącym: naciska go ten, kto hasła nie pamięta, a wtedy
zmiana ze znanym hasłem jest drogą donikąd. Zmiana hasła ze znanym hasłem
jest czynnością Ustawień.

Odnośnik potwierdzenia z listu wprowadza krok drugi rejestracji i pozwala
z niego wyjść. Bez tego odnośnika krok potwierdzenia prowadziłby wyłącznie
z udanej rejestracji w tym samym oknie; odświeżenie strony między listem
a przepisaniem drogi zostawiałoby Operatora przed formularzem wejścia bez
żadnej drogi dalej, ponieważ rejestracja odmawia wtedy kodem `conflict`,
a logowanie oczekiwaniem na potwierdzenie adresu. To samo dotyczy listu
odczytanego na innej maszynie.
