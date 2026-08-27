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

## budowa/klient-poprzedni/src/moduly/terminal/okno-process-monitor.ts

Okno Process Monitor pokazuje na żywo rejestr procesów rdzenia terminala, w tym proces
zainicjowany poleceniem sztucznej inteligencji, zgodnie z rejestrem procesów rdzenia serwera.
Aktualizacja na żywo znaczy ze zdarzeń: okno odczytuje wykaz przy wejściu i na wyraźne żądanie,
a każdą późniejszą zmianę przynosi zdarzenie rdzenia. Odpytywanie w pętli dałoby ten sam obraz
drożej i z opóźnieniem, a przy stu procesach zalałoby gniazdo. Filtr stanu i inicjatora jedzie do
rdzenia parametrem komendy, bo tak stanowi kontrakt i tak wynik jest spójny z dziennikiem rdzenia;
grupowanie jest wyłącznie porządkiem wyświetlania i zostaje w oknie. Wyjście na żywo idzie do
wspólnego bufora konsoli wyjścia, ale bufor żyje jedno połączenie: po rozłączeniu i ponownym
podłączeniu ma zero fragmentów, choć rdzeń wciąż oddaje pełną treść. Osobny przycisk pozycji pyta
o wyjście procesu wprost i pokazuje je przy pozycji, a nie w buforze, więc żaden wiersz nie
wchodzi do konsoli dwa razy.

Okno nie ma dziś ani jednej pozycji bez pokrycia w rdzeniu: wstrzymanie procesu było ostatnią
i już stoi przy wierszu wykazu.

Wstrzymanie albo wznowienie procesu w rdzeniu niesie pole odpowiedzi mówiące, czy system maszyny
rdzenia w ogóle zna wstrzymanie obcego drzewa procesów. Fałsz tego pola znaczy, że proces został
nietknięty, ponieważ platforma tego nie umie — co jest czymś innym niż niepowodzenie czynności.
Okno mówi to wprost, zamiast pokazywać powodzenie przy procesie, który dalej zajmuje procesor.

Czas oczekiwania rdzenia na domknięcie procesu przed odczytem wyjścia wynosi trzy sekundy.
Kontrakt dopuszcza sześćdziesiąt tysięcy milisekund, ale czekanie zajmuje zadanie gniazda,
a minuta bez odpowiedzi wygląda jak zawieszone okno. Trzy sekundy wystarczają, aby polecenie
krótkie zdążyło się domknąć i oddało komplet wraz z kodem wyjścia; polecenie długie i tak odda
wyjście dotychczasowe ze stanem biegnącym, bo czekanie nie jest warunkiem odpowiedzi, a kto chce
zobaczyć resztę, odczytuje ponownie.

Zwinięcie podglądu wyjścia nie pyta rdzenia: drugie kliknięcie zdejmuje treść z widoku i mówi to
wprost, inaczej nie dałoby się odróżnić zwinięcia od odczytu, który wrócił pusty.

Zdanie potwierdzenia odczytu wyjścia mierzy treść, nie sam fakt odpowiedzi. Wyjście puste jest
przebiegiem udanym, bo polecenie mogło nic nie wypisać, ale różni się od wyjścia niepustego —
inaczej potwierdzenie znaczyłoby to samo w obu przypadkach. Rozjazd stanu też idzie wprost: wykaz
w oknie jest kopią wcześniejszego odczytu, a odpowiedź na odczyt przychodzi z tej chwili, więc gdy
się różnią, świeższa jest odpowiedź.

Odmowa rdzenia przy braku wyjścia procesu idzie dosłownie, ponieważ sama nazywa dokładny powód:
czy proces nigdy nie ruszył, wypadł z historii, czy rdzeń był uruchomiony ponownie. Parafraza
zgubiłaby wszystkie trzy powody.

Złożenie powierzchni rejestru procesów jest konstrukcją czystą: nie domyka się na stanie okna ani
na rdzeniu, więc dała się wyjąć bez przenoszenia zależności.

## budowa/klient-poprzedni/src/moduly/roundtable/zrodlo-arsenalu.ts

Rozdział komend modułu Roundtable na dwa pliki idzie po roli, nie po
wielkości. Cztery komendy źródła debaty prowadzą samą debatę: dodają
uczestnika, otwierają turę, moderują ją i czytają stanowisko. Komendy tego
pliku pracują nad zapisem debaty — czytają go, analizują, oceniają, wydają.
Okno rozmowy może stać wyłącznie na komendach źródła debaty; rozszerzenia
boczne stoją na komendach tego pliku.

Każda czynność oddaje typ Wynik, nie samą treść, i każda sprawdza kształt
odpowiedzi. Powód jest ten sam co w źródle debaty: obszar odmawia z powodów
zwyczajnych — kanał uczestnika bywa nieczynny, głosowanie bywa zamknięte —
a okno musi odróżnić brak treści od nieudanego zapytania, więc żaden
odczyt nie zastępuje braku odpowiedzi pustą wartością domyślną.

## budowa/klient-poprzedni/src/moduly/roundtable/czynnosci-arsenalu.ts

Plik istnieje po to, żeby żadna komenda obszaru nie została bez drogi z okna.
Wcześniej pozycje bez obsługi stały jako przyciski nazywające brak; obecnie
każda z nich naprawdę woła swoją komendę, a wynik melduje w oknie, z którego
padła.

Dane żądania pochodzą ze stanu debaty wspólnego oknom modułu: okno, tura
bieżąca, skład, wypowiedzi. Czynność, dla której stan nie ma jeszcze wskazania
— nie ma tury, nie ma uczestników, nie ma wypowiedzi — nie idzie do rdzenia po
to, żeby dostać odmowę: melduje brak wskazania od razu i nazywa, czego
brakuje. Rdzeń odmówiłby tak samo, tylko po podróży tam i z powrotem.

Plik nie podstawia wartości domyślnych za Operatora. Tam, gdzie komenda
potrzebuje treści — stanowisko, zdanie odrębne, warianty głosowania — treść
przychodzi z pola okna, a przycisk bez wypełnionego pola melduje, czego
brakuje.

## budowa/klient-poprzedni/src/moduly/studio/styl-panel-arkusza.ts

Panel prowadzi trzynaście czynności na fragmencie dokumentu: style nazwane
(wykaz, założenie, zmiana, stosowanie, usunięcie), styl znaku, styl akapitu,
czyszczenie formatowania, wielkość liter, malarz formatów, zaznaczenie wedle
podobnego formatowania, zamianę wraz z postacią oraz tabulator zakładany
liczbą. Praca na fragmentach jest osią całego zamówienia.

Zapis stylu nazwanego przestawia wszystkie miejsca dokumentu, które danego
stylu używają. Panel mówi to wprost przy przycisku i oddaje liczbę
przestawionych miejsc z bilansu odpowiedzi, żeby Operator widział skutek
zapisu, a nie samo potwierdzenie wykonania.

Malarz formatów prowadzi dwie czynności rdzenia: pobranie postaci oddaje
uchwyt, naniesienie stosuje go gdzie indziej. Uchwyt pamięta warstwa wyżej;
panel pokazuje wyłącznie, czy coś jest pobrane.

Panel nie woła rdzenia i nie zna dokumentu ani zaznaczenia. Składa treść
żądania z pól bez identyfikatora dokumentu i bez zakresu — zakres dokłada
warstwa wyżej z bieżącego zaznaczenia, bo ona jedna wie, co Operator
zaznaczył.

## budowa/klient-poprzedni/src/moduly/studio/wstazka-pracy.ts

Czynności wstążki stoją zebrane wedle rodzaju pracy, wzorem pakietu biurowego,
zamiast jednego paska ze stoma czynnościami naraz: taki pasek byłby wykazem,
po którym Operator szuka wzrokiem, a zakładka mówi wprost, gdzie czego szukać,
i nazwana grupa — dlaczego te czynności stoją razem.

„Asystent" jest osobną zakładką wstążki, nie pozycją w menu: niesie suwaki
koncepcyjne, wykaz operacji z pliku kategorii operacji i decyzję o wyniku.
Wykaz operacji nie powstaje w tym module od nowa — pochodzi z tego samego
pliku, co wykaz Tools Panelu, żeby obie powierzchnie pokazywały ten sam zestaw.

Wstążka nie woła rdzenia i nie zna stanu modułu. Buduje kontrolki i zgłasza
naciśnięcia oknu przez czynności wstążki, a złożone kawałki — formularz
wczytania, panel znajdź/zamień, pola różnicy — przyjmuje gotowe jako gniazda.
Dzięki temu jedno miejsce trzyma układ wstążki, a inne prowadzi rozmowę
z rdzeniem.

## budowa/klient-poprzedni/src/moduly/studio/kopie-panel.ts

Panel łączy cztery wymagania trwałości dokumentu w jednym miejscu. Autozapis:
odstęp i zapis przy zdarzeniach okna jako jawne, odwracalne ustawienie;
obejmuje treść i postać dokumentu i idzie osobnym szeregiem wersji, żeby nie
zaśmiecał historii Operatora. Kopia zapasowa: zakładana przed zapisem
i niezależnie od historii wersji, żeby przetrwała awarię procesu i awarię
zapisu. Przywrócenie po nagłym zamknięciu: Studio zgłasza je samo, zdaniem
o niezapisanym dokumencie z podaniem godziny, zamiast czekać, aż Operator się
domyśli; przywrócenie do nowego dokumentu stoi obok przywrócenia na miejsce,
żeby przywracanie nie kasowało tego, co już jest. Powrót do stanu pierwotnego:
jedno polecenie, bez szukania wersji założycielskiej w wykazie; wersje nowsze
zostają, więc powrót jest odwracalny.

Wskaźnik stanu zapisu ma cztery wartości i „nieudany" jest jedną z nich.
Nieudany zapis samoczynny wraca odpowiedzią udaną z polem `saved: false` —
panel pokazuje go jako niepowodzenie wraz z powodem, bo pokazanie „zapisano"
po nieudanym zapisie kosztowałoby Operatora pracę.

Panel woła rdzeń sam; treść i postać do odłożenia bierze z kontekstu, bo
powierzchnia dokumentu jest po stronie okna.

## budowa/klient-poprzedni/src/moduly/studio/styl-panel-list.ts

Wzór numeru numeracji jest polem, nie wyborem z wykazu, ponieważ numeracja
prawnicza wielopoziomowa typu 1.1.2 jest wzorem, nie pozycją wykazu gotowych.
Format `legal` mówi rdzeniowi rodzaj numeracji, a wzór mówi jej postać —
pismo urzędowe wymaga obu wartości naraz i nie da się go zamknąć w wykazie
gotowych wzorów.

Znaki tablicy są znakami ostatnio użytymi, nie ulubionymi, ponieważ kontrakt
oddaje pole `recentlyUsed` przy każdym znaku i to ono rozstrzyga, co stoi pod
ręką. Panel nie zakłada drugiego pojęcia obok tego pola: wykaz jest jeden
i pochodzi z rdzenia, więc znaki ostatnio użyte przeżywają zamknięcie karty.

Panel nie woła rdzenia i nie zna dokumentu — składa treść żądania i oddaje ją
warstwie wyżej, która dokłada dokument i zakres przed wywołaniem komendy.

## budowa/klient-poprzedni/src/moduly/studio/zrodlo-wstawien-studio.ts

Siedemnaście komend studio wnoszą treść do dokumentu i wydają go dalej, obejmując wejście,
wydanie, wniesienie ze źródła i cyfryzację. Wszystkie prowadzą tę samą drogę: coś spoza
dokumentu wchodzi do dokumentu albo dokument wychodzi na zewnątrz. Założenie dokumentu zakłada
pustą stronę, wniesienie pliku i PDF wnosi materiał wprost do edytora, wniesienie obrazu stawia
go w miejscu kursora, wniesienie z biblioteki i ze strony sieci niesie fragment wraz z zapisem
pochodzenia, a cyfryzacja prowadzi kolejkę wczytywania po stronie rdzenia. Wydanie do formatu,
wydanie wsadowe, zapis pod nazwą i kopia zamykają drogę w drugą stronę.

Bilans nie jest ozdobą odpowiedzi. Rdzeń oddaje przy wniesieniu bilans odzyskania, przy
czynności bilans zmiany, a przy wydaniu wykaz cech pominiętych przez format docelowy. Zdania
składające te trzy bilanse stoją w jednym miejscu, bo czytają je wszystkie panele wstawień
naraz. Milczące zgubienie tabeli przy wydaniu do tekstu czystego jest ciszą niedopuszczalną,
więc okno nie ma drogi, którą mogłoby bilans pominąć.

Wywołanie uczciwe osłania każdą komendę, ponieważ koperta nieznanej komendy studia nie niesie
pola stanu i bez tej osłony okno stałoby w ładowaniu bez końca.

Zdanie o bilansie czynności traktuje zero zmian jako wynik prawdziwy, nie jako powodzenie:
czynność, która nie tknęła ani jednego miejsca, jest odpowiedzią „nic się nie stało" i tak ma
być przeczytana. Pominięcie z powodu blokady nazywa tę blokadę, aby było wiadomo, która zapora
zatrzymała czynność, zamiast pozostawiać wrażenie częściowego powodzenia.

Zdanie o bilansie wniesienia odzwierciedla to, że odzyskanie z formatu PDF jest odtworzeniem,
nie odczytem: liczby stron z warstwą tekstową i bez niej, tabel rozpoznanych i nierozpoznanych,
obrazów osadzonych i pominiętych stoją obok siebie. Skierowanie na rozpoznanie pisma jest
zdaniem osobnym, ponieważ to inna decyzja niż samo dalsze wniesienie materiału.

Zdanie o wyniku wydania nazywa wprost stratę cech, gdy format docelowy jest uboższy niż
dokument: format uboższy jest sytuacją normalną, ale przemilczenie straty nie jest, więc
wydanie, które zgubiło tabelę albo przypisy, mówi to wprost zamiast pokazywać samo „zapisano".

## budowa/klient-poprzedni/src/moduly/studio/obiekt-panel.ts

Panel nie rysuje ani jednego kształtu i nie prowadzi biblioteki ikon. Komenda
studio.object.insert przyjmuje rodzaj kształtu i nazwę ikony, a rachunek stoi
w module Design, który zakłada kształt od razu jako węzły ścieżki albo
wyszukuje ikonę. Wskazanie węzła Designu idzie osobnym polem — to jedyna droga
osadzenia kształtu złożonego w Designie.

Rodzaj „wykres" stoi w kontrakcie, ale rdzeń go odmawia, ponieważ rachunku
wykresu po stronie Studia nie ma. Odmowa jest widoczna zanim Operator naciśnie
kontrolkę i mówi, co zrobić zamiast tego: złożyć wykres w module Design
i osadzić go w dokumencie jako obiekt wskazany węzłem. Kontrolka kończąca się
odmową rdzenia byłaby obietnicą bez pokrycia, a jej ukrycie zabrałoby
Operatorowi wiedzę, że taka droga istnieje.

## budowa/klient-poprzedni/src/moduly/terminal/okno-terminal-tabs.ts

Do rdzenia idą dwie komendy: otwarcie karty, obejmujące nową kartę i duplikat,
oraz wysłanie polecenia. Podział widoku, rozmiar pisma, schemat barw,
przypięcie, nazwa karty i eksport transkryptu są czynnościami widoku — kontrakt
nie ma dla nich komend, bo nie zmieniają niczego po stronie rdzenia. Zamknięcie
karty jest wyjątkiem opisanym niżej.

Zdanie potwierdzenia po otwarciu karty albo uruchomieniu polecenia składa się
wyłącznie z pól odpowiedzi rdzenia — stan karty, katalog roboczy, stan procesu,
polecenie — i rozstrzyga o powodzeniu tym, co rdzeń oddał, a nie tym, że
żądanie poszło. Czynności czysto widokowe mówią wprost, że rdzeń o nich nie
wie.

Zamknięcie karty w rdzeniu ma dziś komendę kontraktu, ale okno jej jeszcze nie
wywołuje. Pozycja mówi o tym zdaniem liczonym z odczytu wykazu komend rdzenia,
a nie stałym napisem: napis orzekałby o stanie rdzenia z chwili pisania kodu
i nie przestałby go orzekać w dniu, w którym rdzeń dostanie odpowiednią
obsługę.

Zamknięcie karty zachodzi najpierw w rdzeniu, potem w widoku. Kolejność jest
rozstrzygnięciem: karta zdjęta z ekranu przed odpowiedzią rdzenia znikałaby
Operatorowi także wtedy, gdy rdzeń zamknięcia odmówił, a wtedy powłoka
biegłaby dalej bez żadnego widoku na siebie. Przypięcie jest bramą, nie
ozdobą: pasek kart zapowiada, że karty przypiętej nie zamyka się jednym
kliknięciem, więc panel akcji odmawia tak samo jak krzyżyk na pasku. Procesy
karty zostają biegnące — zamknięcie zakładki nie ma prawa przerwać
budowania, które trwa trzecią minutę; zakończenie procesów jest osobną
czynnością.

Rozmiar pisma i schemat barw stawiane są na węźle modułu, a nie na węźle
okna: pismo konsoli ma być jedno w całym module, bo podgląd wyjścia
w monitorze procesów niesie ten sam strumień co ogon karty. Węzeł modułu
bierze się z drzewa, bo okno nie zna swojego gospodarza w chwili budowy.

## budowa/klient-poprzedni/src/moduly/translate/zrodlo-warsztatu-translate.ts

Warsztat modułu translate jest jednym źródłem dla wszystkich rodzin komend, które nie mieszczą
się w czterech oknach pierwotnych: pamięci jako bytu widocznego dla operatora, segmentacji,
terminologii, korekty i spójności, profili kontroli jakości, obiegu zatwierdzeń, dokumentów,
lokalizacji oprogramowania, napisów i dubbingu, silników, polityki pivota, przebiegu
pakietowego oraz wymiany zewnętrznej. Jedno źródło zamiast jednego pliku na okno wynika
z tego, że wszystkie te komendy idą tą samą drogą wywołania i różnią się wyłącznie nazwą
i ładunkiem; rozbicie ich na dwanaście plików dałoby dwanaście kopii tej samej obudowy, z których
każda musiałaby osobno pamiętać o odmowie komendy nieznanej.

Sprawdzian kształtu towarzyszy każdej komendzie i nie jest formalnością: odróżnia rdzeń, który
odpowiedział wynikiem, od rdzenia, który odpowiedział pustą kopertą. Okno przyjmujące pustą
kopertę jako wynik pokazałoby pustkę jako skutek, czyli dokładnie tę szkodę, przed którą stoi
cały moduł.

## budowa/klient/src/wejscie/wejscie.test.ts

Sprawdzian obejmuje przebieg bez rdzenia uruchomionego: każda odsłona, w tym
odsłona błędu połączenia i odsłona zwłoki nałożonej przez rdzeń, jest tu
osiągalna naprawdę, a nie tylko opisana. Rozmowę z rdzeniem naprawdę
uruchomionym mierzy osobny sprawdzian, który rdzenia wymaga.

## budowa/klient-poprzedni/src/moduly/studio/blokada-panel.ts

Panel łączy blokady fragmentów, zajęcia wykonawców i nastawy pracy kilku agentów naraz, bo
wszystkie trzy odpowiadają na jedno pytanie: czego modelowi nie wolno tknąć i kto teraz pisze
po którym akapicie. Blokada jest trwała i skierowana przeciw modelowi; zdejmuje ją wyłącznie
operator konta, a zasięg obejmujący także operatora jest osobnym, jawnym ustawieniem, nie
zachowaniem domyślnym. Zajęcie fragmentu jest czasowe i skierowane przeciw drugiemu wykonawcy;
odmowa nazywa wykonawcę i czas, bo cicha odmowa kazałaby zgadywać, dlaczego fragment nie
drgnął. Nastawy rozstrzygają, czy pętla wykonawcza i praca wielu agentów w ogóle stoją, ilu
wykonawców pracuje naraz i co się dzieje przy spięciu; oba narzędzia są domyślnie wyłączone
i włącza je operator konta, nie okno. Panel woła rdzeń sam, bo skutkiem każdej z tych czynności
jest wykaz albo odmowa nazwana — jedno i drugie jest treścią widoczną, nie wartością pośrednią
przekazywaną wyżej.

## budowa/klient-poprzedni/src/moduly/terminal/okno-script-library.ts

Do rdzenia idą dwie komendy uruchomienia: jedna wykonuje treść skryptu,
a druga oddaje wynik sprawdzenia składni. Rdzeń podaje treść programowi
powłoki jako pojedynczy argument, więc skrypt wieloliniowy wykonuje się bez
zapisywania go do pliku — dlatego biblioteka nie potrzebuje ani komendy
zapisu pliku, ani ścieżki na dysku serwera.

Każdy zapis pozycji zakłada kolejną wersję, więc poprzednie brzmienie treści
zostaje w dzienniku wersji; usunięcie zabiera pozycję wraz ze wszystkimi jej
wersjami naraz.

Analizę treści prowadzi program leżący na maszynie rdzenia, osobny od
kontroli wstępnej okna. Odpowiedź niesie pole mówiące, czy program analizy
był dostępny, i okno je pokazuje, bo pusty wykaz uwag przy braku programu
znaczyłby fałszywie „treść bez zastrzeżeń". Kontrola wstępna okna zostaje
obok analizy, nie zamiast niej: sprawdza to, co da się rozstrzygnąć bez
żadnego programu, i robi to w chwili pisania.

## budowa/klient-poprzedni/src/moduly/studio/nastawy-strony.ts

Kształt nastaw strony nie jest wymyślony: to `StudioPageSetup` z kontraktu,
ten sam, który jedzie w profilu wydania i w renderze podglądu. Okno pracy
pokazuje więc kartkę o tych samych nastawach, którymi rdzeń wyda dokument.
Rozmiar nośnika w milimetrach stoi w tym pliku, bo kontrakt niesie samo
oznaczenie, a nie jego wymiary; wykaz jest zamknięty i nazwany, a nośnik
spoza wykazu bierze wymiary A4 i okno mówi o tym wprost, zamiast rysować
kartkę o zgadniętym rozmiarze. Plik nie zna DOM: wejściem są nastawy,
wyjściem liczby i zdania.

Źródłem prawdy o wymiarach nośników jest rdzeń: komenda właściwa temu
obszarowi oddaje szereg A, szereg B, Letter, Legal, Tabloid oraz koperty
DL, C4, C5 i C6 wraz z wymiarami, czytając je z jednego wspólnego wykazu
rdzenia. Okno pracy woła tę komendę przy wczytaniu dokumentu i wchłania
odpowiedź, więc kartka rysuje się wymiarami rdzenia, a nie wymiarami
wbudowanymi. Wykaz wbudowany stoi w tym pliku wyłącznie dlatego, że kartkę
trzeba narysować, zanim odpowiedź rdzenia dojedzie — rysowanie w rozmiarze
zgadniętym byłoby nieprawdą o nośniku. Poprawki wymiarów wnosi się do
wykazu rdzenia; wpis wbudowany ma się do niego równać, a nie odwrotnie.
Sprawdzian pokrywający ten plik pilnuje, żeby wchłonięcie naprawdę
nadpisywało wymiary wbudowane.

Rodzaj nośnika przy wchłanianiu wykazu rdzenia bierze się z pola `kind`
odpowiedzi, gdy ta go niesie, a nie z wpisu wbudowanego — inaczej koperta
dołożona kiedyś do wykazu rdzenia wchodziłaby do okna jako arkusz i nadruk
koperty nie miałby się na czym wykonać.

Margines na oprawę nie ma pola w kontrakcie — `StudioPageSetup` niesie
cztery marginesy i nic poza nimi — więc oprawa żyje przez sesję okna
i jest zgłoszona jako brak pozycji kontraktu; rysowana jest prawdziwie,
bo pas oprawy zabiera pole pisania. Marginesy odbicia włączone znaczą, że
margines lewy jest marginesem wewnętrznym: na stronie nieparzystej stoi po
lewej, na parzystej po prawej — bez tego oprawa wypadałaby raz w rowku, raz
na krawędzi.

Wysokość bloku przy rozdzielaniu stron podaje wołający, bo zmierzyć ją
potrafi wyłącznie przeglądarka, a rachunek podziału ma dać się sprawdzić
bez niej. Blok wyższy od całej strony zostaje na stronie własnej, ponieważ
dzielenie go w środku wymagałoby łamania wiersza wewnątrz akapitu, czego ta
warstwa nie umie i czego nie udaje.

## budowa/klient-poprzedni/src/moduly/studio/czynnosci-redakcji.ts

Katalog jest osobny od katalogu czynności warsztatu, choć typy dzieli. Powód
leży w materiale, na którym te dwa zbiory pracują: warsztat bierze dokument
PDF z magazynu okna, redakcja bierze dokument Studia wczytany w edytorze.
Jeden wspólny wykaz kazałby Operatorowi wybierać czynność, która nie ma na
czym pracować, i dowiadywać się o tym dopiero z odmowy.

Nazwy pól są nazwami kontraktu i jadą do rdzenia bez zmiany; etykiety są
zdaniem Operatora i z nazwami się nie pokrywają.

## budowa/klient-poprzedni/src/moduly/studio/zrodlo-kontroli-studio.ts

Źródło pracy niesie zmiany śledzone, profile wydania, podgląd, szablony
i komentarze — to czynności redakcji. To źródło niesie czynności kontroli nad
tym, co się z dokumentem stało i czego modelowi nie wolno: odwracalny dziennik
czynności z cofnięciem pojedynczym i ponowieniem; wszystko, co zrobił model,
wraz z licznikiem, skakaniem i cofnięciem z zachowaniem pracy Operatora;
różnicę postaci dwóch wersji i przeniesienie pojedynczego fragmentu do stanu
bieżącego; zapis samoczynny osobnym szeregiem, kopie zapasowe i powrót do
wersji założycielskiej; znakowanie fragmentu w rdzeniu wraz z rodzajami
znaczników własnych Operatora; blokadę fragmentu obowiązującą w rdzeniu;
zajęcie fragmentu przez wykonawcę wraz z nastawami pracy kilku wykonawców
naraz; oraz schowek dokumentu, który sprawdza blokady i odkłada wpis
dziennika, więc wklejenie da się cofnąć pojedynczo.

Każde wywołanie idzie drogą odporną na odmowę: komenda bez uchwytu w rdzeniu
wraca kopertą bez pola stanu, której korelacja nie rozstrzyga, więc zwykłe
wywołanie zostawiłoby okno w ładowaniu bez końca. Droga odporna zamienia taki
brak na zwykłe niepowodzenie nazywające komendę, żeby okno mogło pokazać
Operatorowi, po czyjej stronie jest brak.

Źródło nie ma stanu i nie buduje elementu — sprawdza kształt odpowiedzi
i oddaje ją oknu. Nazwy metod noszą wspólny przedrostek, więc suma
z pozostałymi źródłami modułu nie ma kolizji nazw.

## budowa/klient-poprzedni/src/punkty-izolacji/obszar-techniczny.ts

Obszar zakresu technicznego niesie osiem punktów izolacji technicznej: katalog roboczy sesji,
środowisko procesu, katalog danych modelu, dostęp sieciowy, odczyt i zapis plików, konto i token,
model procesu oraz serwer wykonania. Wartości kluczy są zawsze domyślnie wyłączone. Różnica
wobec obszaru kontekstu jest zasadnicza: tam wartość domyślna opisuje zakres widoczności i obie
wartości coś znaczą, tu wyłączony nie jest słabszym wariantem ochrony. Egzekutor rdzenia bierze
wyłącznie klucze włączone i odrzuca wykonanie, które by je naruszyło; klucz wyłączony nie jest
sprawdzany w ogóle, więc nikt nie pilnuje tego zakresu, dopóki nie zostanie świadomie włączony.
Stan wyjściowy platformy to pełna swoboda operacyjna, a widok nazywa ten skutek przy każdym
z ośmiu kluczy osobno, nie tylko raz na górze, żeby przewijanie wykazu nie zgubiło ostrzeżenia.

Odczyt woła komendę odczytu zakresu technicznego, a przełączenie klucza woła komendę zapisu
z pełnym zestawem ośmiu przełączników naraz, bo kontrakt niesie cały zestaw, nie pojedynczy
klucz. Wartość na ekranie zmienia się wyłącznie po potwierdzeniu z rdzenia: stan modułu to
zawsze ostatnia odpowiedź odczytu albo zapisu, nigdy przewidywanie kliknięcia. Odmowa zapisu nie
idzie do paska ogólnego obszaru, tylko do komórki zmiany wiersza klucza, którego dotyczy, aby
było wiadomo, który z ośmiu kluczy się nie zapisał; wartość wyświetlona zostaje ta sprzed próby,
bo rdzeń jej nie potwierdził.

Plakietka wartości ma trzy stany i trzy wyglądy. Stan nieznany, gdy rdzeń nie zwrócił klucza,
i stan wyłączony, gdy rdzeń zwrócił, że nikt nie pilnuje, prowadzą do przeciwnych wniosków
o bezpieczeństwie maszyny, więc nie mogą wyglądać tak samo. Warianty plakietki są zdefiniowane
w bibliotece komponentów: sukces, ostrzeżenie, błąd, informacja, sygnał, rola, atrament — poza
tym zbiorem plakietka nie dostaje żadnego wyróżnienia.

## budowa/klient-poprzedni/src/moduly/studio/okno-session-repository.ts

Historia wersji stała dotąd wąską kolumną na dole modułu i zabierała miejsce
także wtedy, gdy nikt do niej nie zaglądał. Element wyzwalacza jest teraz
wąskim paskiem ze znacznikiem liczby wersji i przyciskiem; cała tabela
otwiera się nakładką nad treścią i schodzi naciśnięciem albo klawiszem
Escape — powierzchnia dokumentu wraca do dokumentu, gdy historia nie jest
używana.

Repozytorium narasta i niczego nie usuwa samo: przywrócenie wersji
przestawia treść bez usuwania wersji nowszych, więc samo cofnięcie jest
odwracalne. „Usuń z wykazu" jest czynnością miejscową i jawną — schowaniem
pozycji w tym oknie, nie usunięciem wersji w rdzeniu; okno mówi to wprost,
żeby Operator nie sądził, że stracił wersję.

Zapisy samoczynne są domyślnie ukryte, z przełącznikiem pokazującym też je.
Gdyby wchodziły do jednego wykazu z wersjami nazwanymi, historia zasypałaby
się w kilka minut. Rozróżnienie bierze się z pola oznaczającego kamień
milowy i z etykiety wersji — drugiego pojęcia okno nie zakłada.

Czynności rozgałęzienia debaty to okno dziś nie pokazuje: rodzina komend
rozgałęzienia jest zbudowana i pracuje dalej, po prostu nie dostaje w tej
turze własnego miejsca w interfejsie.

Zaplecze czynności obejmuje etykietowanie wersji, eksport historii i paczkę
przekazania. Bez zaplecza trzy czynności stoją jako brak nazwany wraz
z komendą, która czeka gotowa w rdzeniu, zamiast jako przycisk milczący po
naciśnięciu.

## budowa/klient-poprzedni/src/moduly/studio/strona-postac-dokumentu.test.ts

Miarą odbioru tego odcinka jest numeracja stron konfigurowana przez Operatora
— styl, umiejscowienie i numer początkowy — oraz marginesy zmienialne.
Sprawdziany mierzą więc treść żądania, które wychodzi z okna do rdzenia, a nie
samo to, że przycisk dał się nacisnąć: nastawa, która nie dojedzie do rdzenia,
nie przeżyje zapisu i jest usterką, nie uproszczeniem.

Druga rzecz mierzona wprost: pole puste znaczy „nie ruszaj tej cechy", a nie
zero. Bez tego każde naciśnięcie przycisku zerowałoby nastawy, których Operator
nie tknął.
