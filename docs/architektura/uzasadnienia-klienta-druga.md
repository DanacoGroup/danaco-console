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

## budowa/klient-poprzedni/src/moduly/studio/okno-ingest-ocr-panel.ts

Panel stał na jednym wywołaniu wydobycia tekstu: jeden język, bez silnika,
bez progu pewności, z kolejką prowadzoną w oknie i ginącą z jego
odświeżeniem. Prowadzi teraz rodzinę komend cyfryzacji: kolejkę po stronie
rdzenia, rozpoznanie z pełnym sterowaniem, poprawkę rozpoznanego słowa przed
przyjęciem, przyjęcie wyniku jako dokumentu wraz z pierwszą wersją i wykaz
urządzeń wejściowych.

Poprawianie słów wymaga pracy na obrazie obok tekstu, nie ciasnego paska:
panel stawia obok siebie wykaz słów wraz z ich położeniem na stronie
i pewnością rozpoznania oraz warstwę tekstową pozycji. Obrazu skanu tu nie
ma i panel mówi to wprost — komendy pobierającej bajty zasobu do
przeglądarki kontrakt nie niesie, więc położenie słowa jest podane liczbami,
a nie zaznaczone na obrazku.

Narzędziownia nie jest stałą kolumną: wchodzi przyciskiem jako nakładka
i schodzi. Kolejka z wieloma pozycjami potrzebuje miejsca na wykaz, więc
dostaje nakładkę, a nie pasek, który zabierałby szerokość także wtedy, gdy
nikt nic nie cyfryzuje.

## budowa/klient/src/wejscie/montaz.ts

Okno jest siatką o trzech wierszach: belka, korpus, pas działań. W wierszu
pasa stoją dwie rzeczy — pas przez całą szerokość i nota w kolumnie
tożsamości. Montaż tutaj zostaje wyłącznie przełożeniem stanu na węzły
dokumentu i zdarzeń na wywołania przebiegu, dzięki czemu każda odsłona jest
osiągalna w sprawdzianie bez przeglądarki — sprawdzian prowadzi przebieg, nie
okno.

Okno przygotowania nie ma belki systemowej, bo stoi już wewnątrz ramy
aplikacji, a ta niesie własną; nie ma też odsłon do przełączania.

Wpisanie stałej czasu zwłoki byłoby obietnicą odmierzania, którego nikt nie
mierzył — dlatego licznik dostaje wyłącznie wartość zmierzoną przez przebieg.

Pole, którego dotyczy usterka, niesie atrybut `aria-invalid`, żeby czytnik
ekranu dowiedział się tego samego co oko. Przy kilku usterkach naraz wykaz
niesie same rozpoznania, bez wskazówek co robić — wskazówkę niesie zaznaczone
pole, a jej powtórzenie kilka razy wypchnęłoby formularz poza okno.

Kontrolka niegotowa nie jest wyłączana atrybutem disabled, tylko oznaczona
jako zajęta atrybutem aria-busy, ponieważ przycisk wygaszony nie mówi,
dlaczego nie działa.

Licznik wypisuje czas we wszystkich odsłonach, także ukrytych — inaczej byłby
pusty przez sekundę po pokazaniu odsłony. Ubywa natomiast tylko licznik
widoczny: czas schodzący za plecami doprowadziłby do tego, że użytkownik
zastaje zero, choć odsłonę zobaczył przed chwilą.

## budowa/klient-poprzedni/src/moduly/studio/dziennik-panel.ts

Historia wersji ma swoje okno, a zmiany śledzone swoje miejsce w treści. Ten
panel dokłada to, czego ani jedno, ani drugie nie ma: dziennik czynności
z cofnięciem pojedynczym i nie po kolei wraz z ponowieniem, przy czym
zależności wypisane przy wpisie oznaczają, że odmowa nie jest zaskoczeniem;
zmiany modelu z licznikiem, skakaniem i cofnięciem wszystkiego albo tylko
odhaczonych, z zachowaniem pracy Operatora; różnicę postaci dwóch wersji wraz
z przeniesieniem pojedynczego fragmentu do stanu bieżącego; oraz schowek
dokumentu — odłożenie i wklejenie fragmentu drogą rdzenia, która sprawdza
blokady i odkłada wpis dziennika, więc wklejenie da się cofnąć pojedynczo.

Panel woła rdzeń sam, bo każda z tych czynności oddaje bilans, a bilans jest
treścią dla Operatora, nie wartością pośrednią: co przeszło, co stanęło
i przez którą blokadę. Przepuszczanie go przez okno nadrzędne kosztowałoby
jedno przełożenie na każdej z dziesięciu dróg, a bilans musi trafić na widok
w całości. Skutek czynności — nową treść i postać — panel oddaje oknu
wywołaniem zwrotnym, bo powierzchnia dokumentu jest po jego stronie.

Panel nie zna treści dokumentu ani zaznaczenia: bierze je z kontekstu, który
podaje okno. Dzięki temu ten sam panel obsługuje dokument w zakładce i drugi
w podziale powierzchni.

## budowa/klient-poprzedni/src/uwierzytelnienie/zrodlo-auth.ts

Źródło komend bramki jest jedynym miejscem znającym nazwy komend kontraktu rodziny auth i kształt
ich żądań. Ekran logowania woła cztery z siedmiu komend tej rodziny: wejście, rejestrację,
przedłużenie tokenu — tyle, ile trzeba, żeby wejść — oraz zmianę hasła, bo ta nie wymaga
zalogowania i stoi za odnośnikiem resetu hasła. Dodanie i usunięcie metody należą do czynności
ustawień wykonywanych po zalogowaniu. Potwierdzenie rejestracji nie jest wołane stąd, bo uchwytu
w rdzeniu nie ma; wywołanie wróciłoby kopertą nieznaną i niczym więcej. Odpowiedź niesie typ
koperty, nie tylko wynik, ponieważ rdzeń odpowiada na komendę bez uchwytu kopertą nieznaną
z tym samym kodem, którym wejście odmawia przy bramce nieustawionej — te dwie odmowy znaczą co
innego, więc wysyłka oddaje typ koperty odpowiedzi razem z treścią.

Metoda wejścia widziana przez operatora ogranicza się do hasła i PIN-u, bo tyle rdzeń dziś
naprawdę otwiera; metoda logowania systemu operacyjnego odmawia z powodu pochodzenia dokumentu,
nie z braku kodu, więc nie jest tu metodą — segment wyszarzony byłby bramą, a segment czynny
kłamstwem.

Powitanie połączenia jest podsłuchane z koperty przechodzącej przez kanał, nie wysyłane osobno,
ponieważ powitanie i tak leci przy każdym nawiązaniu połączenia, a subskrypcja staje przy
składaniu źródła, zanim gniazdo się otworzy, więc odpowiedź nie ma jak przepaść. Gdy powitanie
nie wróci w zadanym czasie, wynikiem jest powitanie puste, czyli rdzeń nie wie, i ekran schodzi
na sondę zamiast czekać bez końca; czas oczekiwania zerowy znaczy, że trzeba wziąć to, co już
jest.

Sonda PIN-u idzie tą samą drogą co sonda hasła i z tego samego powodu: kontrakt nie ma komendy
odczytu metod, a segment PIN-u nie ma prawa stanąć na ekranie, jeżeli nie ma czego otworzyć.
Rdzeń szuka metody po parze rodzaj i urządzenie, zanim spojrzy na sekret, więc żądanie bez
sekretu odpowiada wprost: brak wyniku znaczy, że PIN-u tu nie ma, a niepowodzenie walidacji
znaczy, że PIN jest i zabrakło tylko sekretu. Sonda nie podnosi dławika prób, bo zwłokę
zwiększa wyłącznie sekret niezgodny, a tu sekretu nie ma.

Wejście przez bramkę hasłem albo PIN-em niesie do rdzenia także przełącznik przedłużenia sesji,
nie tylko do magazynu przeglądarki, bo bez tego zaznaczenie przedłużałoby wyłącznie życie zapisu
lokalnego, a sesja gasłaby po dobie roboczej.

Rejestracja zakłada jedyne konto właściciela przy pierwszym uruchomieniu i nie zakłada sesji.
Konto powstaje niepotwierdzone, a token dostępu wydaje dopiero potwierdzenie po drodze przysłanej
listem, dlatego odpowiedź rejestracji nie niesie sesji i ekran przechodzi wtedy do kroku
potwierdzenia.

Prośba o drogę odzyskania konta idzie na adres uwierzytelniający, a odpowiedź jest zawsze taka
sama niezależnie od tego, czy adres pasuje do konta — ekran nie ma z czego wywnioskować, czy
konto istnieje, i tak ma być.

Zmiana hasła ze znanym hasłem dotychczasowym należy do ekranu wejścia, choć wygląda na czynność
ustawień: sprawdza wyłącznie zgodność hasła dotychczasowego z zapisem bramki i nie pyta o sesję.
Zmiana jest więc wykonalna przed zalogowaniem i odnośnik resetu hasła ma co wywołać. Odzyskaniem
hasła zapomnianego to nie jest, bo takiej drogi kontrakt nie ma.

Wysyłka jednej komendy odczytuje całą kopertę odpowiedzi. Subskrypcja staje przed wysłaniem,
a dopasowanie idzie po identyfikatorze żądania, dokładnie tak jak koreluje kanał, więc między
subskrypcją a wysyłką nie ma okna, w którym odpowiedź mogłaby przepaść. Transport kolejkuje
ramki do chwili otwarcia połączenia, więc wywołanie przed nawiązaniem łączności czeka, zamiast
przepaść.

Sonda stanu bramki wysyła wejście metodą hasła bez sekretu, bo kontrakt nie ma osobnej komendy
odczytu stanu, a ekran musi wiedzieć, czy pokazać wejście, czy pierwsze ustawienie hasła. Adapter
bramki najpierw szuka kotwicy, a dopiero potem patrzy na sekret, więc odmowa mówi wprost: brak
wyniku znaczy, że hasło bramki nie istnieje i bramki jeszcze nie ustawiono, a niepowodzenie
walidacji znaczy, że kotwica jest i zabrakło tylko sekretu. Sonda jest drogą zapasową, nie
pierwszą: odpowiedź na to samo pytanie niesie powitanie, które nie wyprowadza stanu z odmowy
i nie przechodzi przez dławik prób; sonda zostaje na wypadek milczenia rdzenia, bo pole puste
znaczy, że rdzeń nie wie, a wtedy trzeba zapytać inaczej, nie zgadnąć. Sonda niczego nie zmienia
w rdzeniu i nie niesie żadnego sekretu.

Identyfikator urządzenia idzie tylko przy PIN-ie, bo metoda hasła kontraktu go nie wymaga —
hasło otwiera bramkę z każdej maszyny, a PIN jest właściwy maszynie i bez identyfikatora rdzeń
odmawia wprost. Tożsamość klienta z powitania urządzeniem nie jest, bo nadaje się ją na czas
uruchomienia, więc PIN chodzi po zapisie trwałym urządzenia.

## budowa/klient-poprzedni/src/moduly/studio/czynnosci-warsztatu.ts

Okno buduje pola z tego wykazu i składa żądanie funkcją zloz czynności.
Piętnaście osobnych formularzy dawałoby piętnaście miejsc, w których nazwa
pola kontraktu bywa wpisana z pamięci; w tym pliku każde pole stoi raz.
Nazwy pól są nazwami kontraktu, bo to one jadą do rdzenia. Etykiety są
zdaniem Operatora i z nazwami się nie pokrywają — Operator czyta „zakres
stron", a rdzeń dostaje pole `pages`.

## budowa/klient-poprzedni/src/moduly/studio/przybornik-znakowania.test.ts

Sprawdziany odcinka znakowania, asystenta, schowka i osadzenia mierzą rzeczy, które da się
zmierzyć bez stawiania okna, i tylko te, w których pomyłka jest cicha: rozróżnienie trzech
bytów marginesu, zawężenie wykazu znakowań, kolejność czynności pływaka liczoną z użycia,
zaporę przed pokazaniem tej samej operacji dwa razy w katalogu oraz zapis pochodzenia. Czego
tu nie ma — sprawdzianu, że przycisk wywołuje komendę — jest tak samo istotne: sprawdzian
takiego kształtu mierzyłby atrapę, którą sam stawia, a nie skutek na dokumencie.

## budowa/klient-poprzedni/src/moduly/studio/tabela-panel.ts

Siatka jest wejściem pierwszym, a pola liczbowe stoją obok jako droga druga —
dla tabeli większej niż siatka i dla pracy z klawiatury. Oba wejścia prowadzą
do tej samej komendy wstawienia, więc nie ma dwóch zachowań.

Każda czynność tabeli oddaje bilans i panel go wypisuje zawsze, nie tylko przy
pominięciu. Scalenie komórek w zablokowanym fragmencie wraca odpowiedzią
pomyślną z pominięciem w bilansie — bez wypisania bilansu wyglądałoby to na
scalenie wykonane. Wpis dziennika idzie do tego samego zdania, bo bez niego
Operator nie wie, co ma cofnąć.

Panel nie liczy sam ani szerokości kolumn, ani wyniku sortowania, ani zamiany
tekstu na tabelę. Wszystko to robi rdzeń; panel składa żądanie, czyta tabelę
z odpowiedzi i pokazuje szerokości policzone, a nie założone.

## budowa/klient-poprzedni/src/moduly/studio/wczytanie-dokumentu.ts

Otwarcie dokumentu przyjmuje trzy wzajemnie wykluczające się pola:
identyfikator dokumentu otwartego wcześniej w sesji, identyfikator pliku
repozytorium albo ścieżkę na urządzeniu. Formularz pyta o źródło wprost,
zamiast zgadywać, bo pomyłka kończyłaby się odmową rdzenia o powodzie
trudnym do odczytania. Wybierak źródła nie jest wykazem plików: wykaz
zasobów Library należy do okna Library Explorer, a to pole przyjmuje sam
identyfikator. Wykaz formatów pochodzi z kontraktu, nie z listy zapisanej
w widoku; format wybiera rdzeń przy wczytaniu, okno go tylko pokazuje.

Wniesienie pliku różni się od otwarcia dokumentu tym, że wnosi postać pliku
Operatora — arkusz stylów, sekcje, tabele i obrazy — zamiast otwierać
pozycję, którą rdzeń już prowadzi. Trzy komendy to robią i dlatego stoją
osobno: wniesienie pliku wnosi dokumenty i szablony biurowe wraz
z rozpoznaniem zapisu znaków, wniesienie PDF-u odzyskuje z niego tekst,
akapity, tabele i obrazy na tyle, na ile PDF je niesie, a wniesienie obrazu
wnosi obraz wprost w miejsce kursora.

Odzyskanie z PDF-u jest odtworzeniem, nie odczytem: PDF nie niesie
struktury akapitu ani tabeli wprost. Panel wypisuje więc bilans za każdym
razem — strony z warstwą tekstową i bez niej, tabele rozpoznane
i nierozpoznane, obrazy osadzone i pominięte. PDF ze samych skanów kieruje
na rozpoznanie tekstu, a panel mówi to wprost, zamiast oddać pustą kartkę
jako gotowy dokument.

## budowa/klient-poprzedni/src/moduly/workspace/hub-planowania.ts

Cztery widoki huba planowania dzielą jeden zbiór zadań. Przełącznik widoku nie zmienia danych,
tylko pytanie zadawane rdzeniowi: lista pyta o wykaz zadań, tablica o stan tablicy kanban, oś
czasu o harmonogram, kalendarz o wpisy kalendarza. Zadanie zmienione w jednym widoku jest tym
samym zadaniem w pozostałych trzech, bo drugiego zapisu zadania w module nie ma. Zadanie
wskazane w wykazie staje się przedmiotem czynności paska akcji; wskazanie idzie polem, a nie
stanem ukrytym, żeby było widać, czego dotyczy usunięcie, zanim się je naciśnie.

## budowa/klient-poprzedni/src/punkty-izolacji/obszar-efektywna.ts

Wartość zapisana świadomie i wartość domyślna muszą wyglądać inaczej, ponieważ
stan wyjściowy platformy to pełna swoboda operacyjna — egzekutor odrzuca
naruszenie tylko przy punkcie włączonym, a wartość wyłączona niczego nie
ogranicza. Pokazanie domyślnej tak samo jak zapisu Operatora sugerowałoby
ochronę, której nikt nie włączył.

Pochodzenie czytamy z odpowiedzi rdzenia na isolation.policy.preview: obecne
pole origin znaczy zapis na wskazanym poziomie, brak origin znaczy wartość
z rejestru definicji. Stan „nieznane" wchodzi wtedy, gdy odpowiedź w ogóle nie
niesie przełącznika danego punktu.

Żądanie podglądu musi nieść punkt widzenia — sesję, warstwę, zasięg i okno.
Przy pustym żądaniu rdzeń schodzi kolejką do poziomu globalnego, więc okno
pokazałoby politykę całej platformy pod nazwą „efektywna".

Czwarty stan „brak-odczytu" jest wyłącznie kliencki i nie wolno go zlewać
z wartością domyślną: oba prowadzą do przeciwnych wniosków o bezpieczeństwie.

Jedenaście punktów izolacji odpowiada definicjom rdzenia w module konfiguracji
izolacji. Objaśnienie każdego punktu jest opisem klucza — co znaczy i jaki
jest stan wyjściowy platformy — a nie wynikiem odczytu, dlatego stoi w pliku
klienta: odpowiedź policy.preview niesie wartości, nie objaśnienia.

Miejscem styku granicy asystenta jest komenda zapisu konfiguracji sesji z roli
klawiatury: pisze przedmioty tych samych punktów — katalog roboczy, katalogi
dodatkowe, środowisko, uruchomienie, narzędzia, uprawnienia i konto —
przełączników izolacji jednak nie rusza, więc egzekutor zostaje na miejscu
i granicy nie przekracza.

## budowa/klient-poprzedni/src/moduly/studio/okno-warsztatu-dokumentu.ts

Wstążka jest kontekstowa: wchodzi po naciśnięciu przycisku, a grupa Ochrona
wstążki okna pracy woła przeniesienie ogniska do niej; wchodzi też
samoczynnie, gdy dokument w pracy jest PDF-em. Nie zajmuje miejsca, kiedy
Operator nad PDF-em nie pracuje.

Grupy Strony, Nakładanie, Treść, Bezpieczeństwo i Narzędziownia cyfryzacji
zawężają wykaz czynności do swojego obszaru, a formularz pod nimi jest
jeden: wybór czynności przestawia pola. Piętnaście osobnych formularzy
dałoby wstążkę trudną do przejrzenia i do utrzymania. Pola pochodzą
z katalogu czynności warsztatu, więc nazwa pola kontraktu stoi w drzewie
raz.

Materiał wchodzi z magazynu okna i wynik do niego wraca. Żadna czynność nie
zmienia materiału w miejscu — okno mówi to przy polu materiału, bo Operator
ma wiedzieć, że pomyłka nie kosztuje go dokumentu źródłowego. Odpowiedź
jest opisana skutkiem, nie słowem gotowe: liczbą stron, liczbą części,
identyfikatorem nowego zasobu.

Kolejność grup wstążki jest kolejnością pracy nad plikiem: najpierw strony,
potem to, co się na nich kładzie, potem treść, na końcu bezpieczeństwo
przed wysyłką. Narzędziownia cyfryzacji stoi grupą osobną, bo pracuje przed
dokumentem — na materiale, z którego dokument powstaje.

## budowa/klient-poprzedni/src/moduly/wiedza/okno-wyszukiwania-znaczenia.ts

Wyszukiwanie po słowach ma osobną drogę; to okno prowadzi wyszukiwanie po znaczeniu oraz
przebudowę wskaźnika. Źródło jest treścią wyniku, nie ozdobą przy nim, bo trafienie bez
wskazania źródła jest bezwartościowe: każdy fragment niesie źródło w nagłówku pozycji, nie
w dymku, źródło jest klikalne przyciskiem kopiującym wskazanie do schowka, bo to jedyne, czym
da się po źródło sięgnąć w innym oknie, a fragment bez wskazania dostaje zdanie o braku
wskazania zamiast martwego przycisku. Trafność też jest widoczna jako plakietka w setnych przy
każdej pozycji, bo bez niej wykaz posortowany malejąco wygląda tak samo przy trafieniu bliskim
i przypadkowym. Przebudowa wskaźnika stoi w tym samym oknie, bo pusty wynik szukania ma
dokładnie dwie przyczyny: nie ma czego znaleźć albo wskaźnik nie został zbudowany, a przycisk
przebudowy obok pustego wyniku pozwala rozstrzygnąć to jednym kliknięciem.

Wiersz ze sterem zamiast pola różni się od wiersza biblioteki kontrolek tym, że tamten buduje
etykietę wiążącą się z pierwszym potomkiem dającym się etykietować, czyli z uchwytem menu —
kliknięcie w podpis otwierałoby wtedy menu, którym podpis nie jest, a nazwa dostępna uchwytu
wchodziłaby w spór z opisem, który mechanizm ustawia sam i który niesie bieżącą wartość nastawy.

Zdanie opisujące każdy zakres wiedzy jest tu potrzebne bardziej niż przy większości nastaw:
wszystko i pliki przestrzeni brzmią podobnie, a przeszukują dwa różne zbiory, i to od nich
zależy, czy pusty wynik znaczy, że nie ma czego znaleźć, czy że szuka się nie tam.

Uchwyt do źródła kopiuje wskazanie zamiast skakać do niego, bo kontrakt nie ma komendy
otwierającej źródło fragmentu — samo wskazanie jest jedynym sposobem sięgnięcia po nie w oknie
biblioteki albo w rozmowie z modelem.

Trafność bez miary widocznej sprawiałaby, że wykaz posortowany malejąco wyglądałby tak samo
przy trafieniu bliskim i przy przypadkowym, a to rozstrzygnięcie należy do czytającego, nie do
samego wykazu; rdzeń ma prawo trafności nie podać, a wtedy plakietka mówi o braku miary, zamiast
zmyślać zero.

## budowa/klient-poprzedni/src/moduly/studio/konwersja-dokumentu.ts

Wydania żąda i pasek narzędzi Studio Editora, i Preview Window; czynność jest
jedna i stoi w jednym miejscu, bo dwie kopie rozjechałyby się przy pierwszej
poprawce zdania o wyniku.

Przedmiotem zamiany jest treść zaakceptowana modułu, nie plik na dysku.
Dokument Studia mieszka w rdzeniu, a nie w systemie plików Operatora, więc
komenda zamiany formatu dostaje treść wprost polem. Wynik jest zasobem
magazynu rdzenia — komenda oddaje jego identyfikator i rozmiar, a nie bajty do
pobrania.

Format dokumentu ma cztery wartości, a zamiana formatu przyjmuje nazwy
własnego słownika. Trzy z czterech mają w nim odpowiednik wprost; PDF
odpowiednika użytecznego nie ma, bo treść jedzie napisem, a nie plikiem —
napis nie jest PDF-em, choćby rdzeń tak nazywał dokument, z którego powstał.
Zwraca się wtedy tekst czysty i mówi się o tym Operatorowi osobnym zdaniem.
Bez tego zdania Operator widziałby format dokumentu z rdzenia i zamawiałby
wydanie z formatu, którego rdzeń w tej drodze nie czyta.

Komenda zamiany treści napisem stoi obok czterech komend wydania dokumentu po
stronie rdzenia, nie zamiast nich: pierwsza nie ma czym przenieść postaci
dokumentu prowadzonego jako PDF albo DOCX, bo treścią jest tam sam tekst.
Cztery komendy niżej pracują na dokumencie po stronie rdzenia i przenoszą
arkusz stylów, sekcje, tabele, obrazy i aparat — zapisują pod nową nazwą,
wydają do formatu, wydają wiele dokumentów naraz albo przestawiają format
dokumentu Studia. Starsza droga zostaje dla treści, która nie jest dokumentem
rdzenia.

Format uboższy niż dokument jest normalną sytuacją; przemilczenie straty nie
jest. Wydanie do tekstu czystego, które zgubiło tabelę i przypisy, mówi to
wprost, a wydanie wsadowe wypisuje bilans dokument po dokumencie — odmowa
jednego nie ukrywa się za liczbą wydanych.

## budowa/klient-poprzedni/src/moduly/roundtable/okno-model-panels.ts

Skład uczestników idzie ze stanu debaty, nie z własnej listy okna: stan subskrybuje zmianę
debaty i dokłada uczestników z przyrostu tury oraz z czynności moderatora, a to okno tylko
czyta i rysuje. Odpowiedzi rosną na żywo z fragmentów strumienia wypowiedzi, a rysowanie
składu wraz z odpowiedziami leży w module składu modeli; subskrypcji na zmianę stanu okno nie
zakłada samo, bo jedna na całe złożenie modułu wystarcza.

Dwie tożsamości jednego kanału dają dwa panele. Kluczem panelu jest identyfikator uczestnika,
nigdy identyfikator kanału, bo tak stanowi kontrakt dodania modelu. Liczba dotychczasowych
wystąpień kanału mówi, ile razy kanał już wystąpił — od drugiego wystąpienia wzwyż panel niesie
oznaczenie powtórzenia, żeby dwa identyczne z pozoru panele nie wyglądały na usterkę.

Świeżo otwarte okno nie zna jeszcze składu debaty. Odczyt składu i odczyt całej debaty są
w kontrakcie, ale to okno żadnego z nich jeszcze nie wywołuje, bo obsługi nie zbudowano; do
tego czasu skład narasta z przyrostu zdarzenia i z uczestników zwracanych przez polecenie
moderatora wprost. Stan pusty na starcie jest stanem poprawnym, opisanym wprost, nie udawaną
pustką.

Kontrakt nie niesie górnego limitu instancji tego okna i okno nie narzuca własnego — skład
rośnie tak, jak rośnie w rdzeniu. Ile uczestników mieści debata, mówi nota nad formularzem,
a nie bramka odmowy: uczestnik ponad zwyczajowy skład przechodzi przez rdzeń bez odmowy, więc
przycisk dodania uczestnika nie zna żadnego progu. Nota podaje, ilu uczestników debata ma, czego
nie liczy ani rdzeń, ani kontrakt, i gdzie leży sufit gniazd sceny okien równoległych.

Dwie pozycje paska akcji mają komendę w kontrakcie i nie mają jeszcze obsługi: usunięcie
uczestnika oraz oznaczenie uczestnika jako kluczowego. Wyciszenia w pasku nie ma wcale, bo jest
czynnością moderatora i należy do panelu moderatora, nie do tego okna. Pozostałe czynności
uczestnika bez obsługi siedzą w zestawie akcji warstwy trzeciej: siedem pozycji ma dziś komendę
w kontrakcie, a ósma — liczba wariantów jednej persony — ma pole w uczestniku i nie ma komendy,
która by je zapisała.

## budowa/klient-poprzedni/src/moduly/studio/zapis-formatowany.ts

Okno pracy z dokumentem pokazuje pismo, a nie jego zapis: nagłówek jest
większy, cytat wcięty, lista wypunktowana, tabela ma siatkę. Zapisem pod
tym wyglądem zostaje markdown, bo to jego niesie pole treści dokumentu —
kontrakt nie ma ani pola stylów, ani pola układu, więc drugi zapis, własny
i bogatszy, nie miałby gdzie dojechać do rdzenia i ginąłby przy pierwszym
zapisie dokumentu. Plik nie zna okna ani stanu modułu: wejściem jest napis,
wyjściem bloki albo element, i odwrotnie — dzięki temu przekształcenia
sprawdza się bez stawiania powierzchni.

Trzy drogi i wszystkie trzy muszą się zgadzać: czytanie napisu na bloki
przy czytaniu dokumentu, rysowanie bloku jako element widoczny przy
pokazaniu dokumentu i odczyt elementu z powrotem na napis przy pisaniu
w oknie. Droga trzecia jest odwrotnością drugiej i pierwszej razem: to, co
Operator wpisze w widoku formatowanym, wraca do markdown tą samą składnią,
którą widok odczytał.

Rodzaj bloku przy jego wyrysowaniu idzie do atrybutu danych elementu, bo
z niego czyta go zarówno arkusz stylów, jak i odczyt powierzchni z powrotem
na napis. Wyprowadzanie rodzaju z samej nazwy elementu działałoby tylko
dopóty, dopóki przeglądarka nie wstawi własnego elementu przy naciśnięciu
klawisza Enter.

Podkreślenia markdown nie ma, więc jedzie znacznikiem HTML — markdown
przepuszcza go nietkniętym, a odczyt powierzchni odkłada go z powrotem tak
samo. Odczyt powierzchni edycji jest odwrotnością wyrysu i musi znieść to,
co dokłada przeglądarka samodzielnie: naciśnięcie klawisza Enter zakłada
nowy blok, wklejenie treści dokłada element ze stylami, a pogrubienie
z klawiatury bywa innym elementem niż ten, którym pisze wyrys.

Długość zapisu markdown do wskazanego punktu jest potrzebna do przełożenia
zaznaczenia w widoku formatowanym na zakres znaków w treści dokumentu —
zakres, który jedzie do rdzenia jako granice zaznaczenia. Liczenie go po
samym tekście widocznym dałoby wartość mniejszą od prawdziwej o długość
znaczników składni.

## budowa/klient-poprzedni/src/moduly/studio/linijka-pozioma.ts

Linijka chwyta margines lewy i prawy, wcięcie pierwszego wiersza, wcięcie
lewe i prawe osobnymi znacznikami, tabulatory zakładane naciśnięciem wraz
z rodzajem i znakiem wiodącym oraz szerokości kolumn tabeli. Każdy chwyt
przyjmuje wartość także z klawiatury, w milimetrach albo calach, bo
przestawianie marginesu pisma urzędowego wyłącznie myszą odcięłoby część
Operatorów.

Wcięcie pierwszego wiersza liczy się względem wcięcia lewego, bo tak zachowuje
się linijka pakietu biurowego: przeciągnięcie wcięcia lewego zabiera pierwszy
wiersz ze sobą. Gdyby oba wcięcia były liczone od krawędzi pola, Operator przy
każdej zmianie wcięcia lewego poprawiałby drugie osobno.

Linijka nie zmienia dokumentu i nie woła rdzenia: oddaje nastawę temu, kto ją
zbudował, a nastawy strony i wcięcia trzyma powierzchnia. Inaczej powstałyby
dwa źródła prawdy o marginesie, które rozjechałyby się przy pierwszym profilu
wydania wziętym z rdzenia.

Wybieracz rodzaju tabulatora pokazuje rodzaj na przycisku, żeby Operator
wiedział, co założy, zanim naciśnie podziałkę; jedno naciśnięcie przestawia
rodzaj następnego zakładanego tabulatora.

Krawędź kolumny tabeli nie może przejść przez sąsiednią, bo kolumna
o szerokości ujemnej nie jest kolumną. Granicą jest sąsiad odsunięty
o najmniejsze pole, a nie sam sąsiad — dwie krawędzie w jednym miejscu dałyby
kolumnę zerową.

## budowa/klient/src/wejscie/przebieg.ts

Przebieg nie zna żadnego składnika okna: rozstrzyga wyłącznie, który etap
i która odsłona obowiązuje oraz co wysłać do rdzenia. Dzięki temu każda
odsłona, także odsłona błędu i wstrzymania, jest osiągalna w sprawdzianie bez
przeglądarki. Nazwy komend i kształty ich treści pochodzą wyłącznie
z kontraktu; przebieg nie powtarza ani jednego literału nazwy komendy.

Odsłona `konto-bez-potwierdzenia` obsługuje rejestrację bez konta nadawczego:
w tym wariancie wejście działa hasłem, a okno musi nazwać adres, którego
nikt nie potwierdził, bo adres jest jedyną drogą odzyskania konta. Z kontem
nadawczym list poszedł i okno prowadzi do jego przepisania zamiast do tej
odsłony.

Progu prób logowania nie ma. Kontrakt stanowi wprost, że po stronie rdzenia
nie ma i nie będzie odmowy „za dużo prób" — nieudana próba nakłada na
następną rosnącą zwłokę, nic więcej. Pomiar to potwierdza: osiem kolejnych
nieudanych prób oddaje osiem razy odmowę uwierzytelnienia, a po nich hasło
poprawne wpuszcza. Odsłona z prototypu zostaje i obsługuje zwłokę, nie
zaporę. Ten sam pomiar zwłoki — czas trwania ostatniego wywołania, bo rdzeń
nakłada zwłokę przed odpowiedzią — nazywa też zwłokę żądań odzyskania konta:
progu wysyłek też nie ma, zmierzone siedmioma kolejnymi wysłaniami przyjętymi
bez odmowy.

Miejsce przechowania tokenu bramki nie jest rozstrzygnięte przez ten teren —
należy do powłoki wołającej. Magazyn działający wyłącznie w pamięci procesu
jest więc wartością domyślną: token ginie wraz z procesem, co jest
zachowaniem uczciwym, bo magazyn udający trwałość obiecywałby rozpoznanie
urządzenia, którego nie ma.

Liczba etapów przygotowania środowiska równa się liczbie komend, jakie droga
wejścia potrafi zmierzyć — uwierzytelnienie i przywrócenie kart sesji. Etap
bez komendy w kontrakcie nie jest etapem czekającym, tylko obietnicą, której
nikt nie wykona, a postęp liczony razem z nim nie dobiegłby końca nigdy.

## budowa/klient-poprzedni/src/moduly/terminal/zrodlo-terminala.ts

Każda czynność oddaje typ Wynik, nie samą treść. W terminalu odmowa jest
zjawiskiem zwykłym, nie wyjątkiem: tryb uprawnień okna albo punkt izolacji
zatrzymuje polecenie kodem odmowy uprawnień, a źródło nie zamienia takiej
odmowy w pustą listę.

Wyjście procesu przychodzi zdarzeniem strumienia fragmentów, w którym
identyfikator komunikatu niesie identyfikator procesu z rejestru rdzenia.
Komenda odczytu strumienia zostaje niewpięta, choć rdzeń ją obsługuje: ogon
historii, który ta komenda oddaje, to co do wiersza ta sama treść, która
płynie do wspólnego bufora okna przez zdarzenie strumienia, a w buforze
nie ma klucza, po którym dałoby się odsiać wiersz przyjęty drugi raz.
Komenda odczytu wyjścia jednego procesu jest wpięta, bo zawęża się do
jednego procesu, a jej odpowiedź trafia do własnego widoku pozycji, którego
bufor nie widzi wcale.

Odczyt sięga dalej niż strumień: bufor żyje jedno połączenie, a rejestr
rdzenia jeden bieg rdzenia. Po ponownym podłączeniu gniazda bufor okna jest
pusty, a wykaz procesów oraz odczyt wyjścia wciąż oddają proces wraz z jego
pełnym wyjściem.

Zaplanowanie zadania powłoki jest zapisem automatyki o jednym kroku
rodzaju polecenia wołającym wykonanie w terminalu. Druga rodzina komend po
tej stronie dałaby dwie prawdy o jednym harmonogramie — i to jest jedyny
powód, dla którego okno składa plan z komend cudzego modułu, a nie ze
swoich.

## budowa/klient-poprzedni/src/okna-rownolegle/uklad-okien.ts

Rozsyłanie odczytu łączności idzie z jednego miejsca układu, a nie z każdego gniazda osobno,
bo łącze jest jedno na całego klienta: gniazdo pierwsze i czwarte muszą po jednej zmianie
transportu pokazać dokładnie to samo. Wejście zostaje publiczne, bo scena podglądu układu
pracuje bez rdzenia i podaje odczyty ręcznie.

Kod okna wykonania powstaje dopiero po uzgodnieniu z rdzeniem, już po zmontowaniu układu, bo
panele pomocnicze wołają komendy żądające tego kodu. Dopóki kodu nie ma, panel nazywa wprost,
czego mu brakuje, zamiast pokazywać pustkę.

Subskrypcja strumienia przeżywa usunięcie węzła z drzewa dokumentu, więc zamknięcie panelu musi
odpiąć subskrypcje rdzenia wprost — scena zdjęta bez tego wywołania zostawiałaby żywe nasłuchy.

Pole opcji o wbudowanej rozmowie jest domyślnie wyłączone, bo scena sesji osadza właściwy widok
z zewnątrz; stanowisko podglądu układu włącza je jawnie, bo nie ma z zewnątrz niczego do
osadzenia. Kanał do rdzenia bywa pominięty — wtedy menu paneli jest puste i mówi to wprost,
bo stanowisko podglądu układu pracuje bez rdzenia i właśnie tak ma wyglądać. Transport bywa
pominięty analogicznie — wtedy plakietki łączności milczą, bo scena nie zna łącza.

Moduł bez rozmowy nie zostawia sceny pustej: scena trzyma minimalną liczbę okien, a zdanie
przełącznika mówi wprost, że rozmowy w tym module nie ma, żeby pusta scena nie czytała się jak
awaria. Pytanie, czy dostawienie gniazda cokolwiek zmieni, idzie tą samą rachubą, którą potem
wykona ustawienie liczby, bo sufit składa się z figury modułu i z górnej granicy, a figura
przestawia się z rdzenia — osobne porównanie rozjechałoby się z tamtym przycięciem przy
pierwszym module węższym niż scena. Odpowiedź przecząca znaczy krótszą listę w menu, a nie
wiersz wygaszony.

Znacznik przestawiania ról przez sam układ rozróżnia rolę z zewnątrz, która jest decyzją
i zostaje, od roli nadanej przez zmianę liczby okien, która jest domyślna i ustępuje następnej;
nadawanie ról jest synchroniczne, więc znacznik nie przecieka poza swoją pętlę.

Pozycja otwarcia nowego okna idzie tą samą drogą co przełącznik liczby: podniesienie liczby
wprowadza gniazdo na scenę, a wejście gniazda na scenę zamawia dla niego okno rdzenia. Okno
założone z boku byłoby niewidoczne, więc tu stoją same wywołania — układ nie wie, że po drugiej
stronie jest wiersz menu.

Moduł sceny idzie z gniazda pierwszego, bo to ono niesie okno uzgodnione z rdzeniem i na nie
przychodzi zmiana modułu, także po wejściu do przestrzeni, które przestawia moduł okna zamiast
zakładać drugie. Scena z oknami w różnych modułach bierze figurę z gniazda pierwszego.

Zdanie sceny o liczbie okien modułu składa do trzech zdań, każde tylko wtedy, gdy jest
prawdziwe: czym jest figura modułu, że pozycja wyższa okna nie dokłada, oraz że scena trzyma
okien więcej, niż moduł prowadzi. Ostatnie bierze się stąd, że układ okien nie zdejmuje ich
sam — zdjęcie okna zamyka je również w rdzeniu i jest decyzją operatora.

Żądanie liczby ponad figurę modułu nie jest odmową: pozycja przełącznika zostaje czynna, scena
bierze tyle okien, ile moduł prowadzi, a zdanie pod grupą mówi, dlaczego wyżej się nie da. Cicha
zgoda na okno ponad figurę modułu byłaby sukcesem udawanym.

Przestawienie sceny na moduł nie zdejmuje okien, nawet gdy jest ich więcej, niż moduł prowadzi,
bo zdjęcie okna ze sceny zamyka je również w rdzeniu, a zrobione samoczynnie przy przestawieniu
modułu byłoby zamknięciem okna, o które nikt nie prosił. Scena mówi o nadmiarze wprost
i zostawia zdjęcie operatorowi.

## budowa/klient-poprzedni/src/moduly/studio/okno-petli-wykonawczej.ts

Powierzchnia należy do dokumentu: panele nie zajmują stałych kolumn, wchodzą
na żądanie i schodzą, gdy nie są używane. Okno pętli wchodzi więc znacznikiem
przebiegu i schodzi po zamknięciu, a stała kolumna jest trybem do wyboru,
jawnym, odwracalnym i pamiętanym — nigdy postacią domyślną. Samoczynne
otwarcie zostaje przy zleceniu wielozadaniowym, bo okno wchodzi wtedy, gdy
naprawdę jest używane, i to jest dokładnie ta sama reguła zastosowana wprost.

Kolejka zadań, obsada wykonawców, warsztat łańcucha i tryb wsadowy działają
na tym samym stanie. Krok łańcucha i dokument wsadu są zadaniami tej samej
kolejki, bo kontrakt mówi wprost: przebieg łańcucha i wsadu prowadzi pętla
wykonawcza okna. Druga maszyneria obok byłaby drugą prawdą o tym samym
przebiegu.

Okno nie przestawia nastaw pętli. Nastawy — czy pętla jest czynna, ilu
wykonawców naraz, co przy spięciu — idą osobną rodziną komend. Okno je czyta
i pokazuje wprost, a gdy pętla jest wyłączona, mówi to zdaniem, zamiast
pozwalać naciskać przycisk, który i tak odmówi.

## budowa/klient-poprzedni/src/moduly/studio/zrodlo-postaci-studio.ts

Nastawy strony, marginesy i numeracja stały dotąd wyłącznie w kliencie, prowadzone
przez sesję okna, a jedyną drogą trwałości był profil wydania — osiem pól bez
oprawy, bez sekcji i bez formatu numeru. Rdzeń niesie dziś rodziny komend
obejmujące stronę, sekcję, styl, formatowanie, listę, znak specjalny, widok,
tabulator, tekst, postać dokumentu i pochodzenie fragmentów — i to one są
jedyną drogą, którą nastawa Operatora przeżywa zapis.

Metody biorą treść żądania z kontraktu w całości, a nie rozłożone argumenty.
Rodzina postaci dokumentu ma pola liczone dziesiątkami — same nastawy strony
to czternaście — i połowa z nich jest opcjonalna w znaczeniu „nie ruszaj tej
cechy". Rozkładanie ich na argumenty zmuszałoby wołacza do podawania wartości
tam, gdzie chce ciszy, a „podano zero" i „nie podano" znaczą tu różne rzeczy.

Wywołanie idzie przez osłonę, bo komenda bez uchwytu w rdzeniu wraca kopertą
ogólną bez pola stanu, którego korelacja klienta nie rozstrzyga — okno stałoby
w ładowaniu bez końca. Osłona zamienia to w zwykłą odmowę nazywającą komendę;
rodzina komend jest świeża, więc rdzeń starszy od klienta jest tu przypadkiem
realnym, nie teoretycznym.

## budowa/klient-poprzedni/src/moduly/studio/aparat-panel.ts

Spis treści, spisy ilustracji i tabel, przypisy dolne i końcowe, podpisy,
zakładki, odwołania wzajemne, odsyłacze, powołania, bibliografia, hasła
indeksu i indeks — wszystkie są elementami wyliczanymi z dokumentu. Pola
dokumentu (numer strony, liczba stron, data, właściwość, pole obliczane) mają
tę samą naturę. Wspólne jest to, że po zmianie treści stają się nieświeże,
a odświeżenie liczy je od nowa.

Panel prowadzi jeden wykaz „do odświeżenia" dla obu rodzin i pisze liczbą, ile
elementów rozjechało się z dokumentem. Spis treści pokazany bez znaku
nieświeżości kłamałby o dokumencie, którego nagłówki się zmieniły.

Numer przypisu i powołania nadaje rdzeń przy odświeżeniu, a panel go
wyłącznie pokazuje: okno nie wymyśla żadnego kodu i nie numeruje niczego samo.

## budowa/klient-poprzedni/src/moduly/terminal/katalog-funkcji.ts

Moduł Terminal obiecuje zastąpić emulator terminala, multiplekser sesji,
klienta SSH i SFTP, harmonogram zadań, uruchamianie skryptów, menedżera
sekretów, konsolę szeregową i klienta kontenerów. Część tych zdolności ma
pokrycie w kontrakcie, część nie, a część wymaga programu, którego instalka
Danaco Console nie niesie. Bez tego wykazu Operator poznaje różnicę dopiero
w chwili, w której czegoś potrzebuje — a wtedy milczenie interfejsu czyta się
jak usterka platformy.

Wykaz jest przepisem rozdziału „Katalog funkcji i narzędzi" dokumentacji
modułu, nie zbiorem pomysłów; liczba jego pozycji jest zmierzona z długości
wykazu, nie wpisana osobno.

## budowa/klient-poprzedni/src/moduly/roundtable/katalog-funkcji.ts

Katalog istnieje po to, żeby okno nie wyglądało na kompletne tam, gdzie nie
jest. Opracowanie modułu wymienia siedemdziesiąt cztery narzędzia — bez wykazu
Operator poznawałby rozjazd między nimi a stanem modułu dopiero po naciśnięciu
każdego przycisku z osobna.

Po scaleniu kontraktu rozjazd zmienił naturę. Obszar komend modułu niósł
cztery komendy i niesie ich dziś czterdzieści sześć, więc zdanie „tego nie
wykona żadna komenda obszaru" przestało być prawdziwe co do czterdziestu
dwóch narzędzi naraz. Brakuje ich obsługi, nie kontraktu — i to rozróżnienie
niesie stan pozycji.

Nazwy narzędzi stoją w brzmieniu opracowania modułu, bez numeracji
rozdziałów — numer rozdziału nie jest identyfikatorem niczego w produkcie.
Warstwa pozycji jest warstwą widoczności z opracowania, przypisaną po
elemencie interfejsu, którym narzędzie się otwiera: pierwsza to widoczność
bez interakcji, druga to znacznik kontekstowy, przycisk albo przełącznik,
trzecia to menu zestawu akcji okna, a czwarta to nastawa dostępna poleceniem,
wyszukiwarką funkcji albo trybem administracyjnym — w oknie nie ma jej
w stanie spoczynku.

Stan wykonania narzędzia niesie pięć wartości. „Częściowa" nie jest stanem
pośrednim między porażką a sukcesem: znaczy, że narzędzie działa w części
zakresu, a nie działa w reszcie, i zdanie pozycji mówi, gdzie przebiega
granica. „Bez obsługi" znaczy, że komenda jest w kontrakcie, a moduł jej nie
wywołuje — to najliczniejszy stan modułu i najważniejszy do odróżnienia: nie
brakuje uzgodnienia, brakuje pracy, a zdanie pozycji nazywa komendę, którą
narzędzie wejdzie. „Bez pokrycia" znaczy, że kontrakt nie ma czym narzędzia
wykonać — nie ma komendy albo ma komendę, lecz brakuje pola w jej żądaniu; po
scaleniu kontraktu stan ten zszedł do pojedynczych pozycji i każda nazywa
brakujące pole. „Poza modułem" znaczy, że narzędzie należy do okna wspólnego
platformy.

## budowa/klient-poprzedni/src/moduly/studio/schowek-zrodlo.ts
Schowek przeglądarki żyje tylko w obrębie karty i nie ma historii, ponieważ
`navigator.clipboard` oddaje jedną, ostatnią treść. Wymaganie produktu
obejmuje wykaz wpisów sprzed kilku ruchów oraz wpisy przypięte na stałe,
więc historia leży w rdzeniu poprzez rodzinę komend `clipboard.*`: `push`
dopisuje treść (powtórzenie identycznej treści nie mnoży wpisów, tylko
podnosi zastany wpis na czoło i oddaje `alreadyPresent`), `list` oddaje
historię wraz z przypiętymi, `pin` przypina wpis (przypięty nie wygasa wraz
z retencją), `delete` usuwa wpis albo całą historię nieprzypiętą.

Komendy schowka nie niosą pola autora, więc wpis odłożony przez model i wpis
operatora są w rdzeniu nierozróżnialne. Odkładając fragment za modelem, okno
zapisuje pole `sourceWindowId` — jedyne pole pochodzenia, które wpis niesie.
Rozróżnienie autora wpisu schowka nie wchodzi w zakres kontraktu.

## budowa/klient-poprzedni/src/moduly/workspace/stany-okna.ts
Przedrostek klas CSS „dw-" jest ustalany lokalnie, bo nośnik stanu treści jest współdzielony
między modułami, a każdy moduł nadaje własny przedrostek. Stan błędu przenosi kod i komunikat
z kontraktu, aby okno po niepowodzeniu zapytania nie wyglądało jak okno z pustą, ale poprawną listą.

## budowa/klient-poprzedni/src/moduly/research/dymek-badania.ts
Znak dymku Research zachowuje trzy własności wspólnej fabryki bez zmian: pokazanie na najechaniu albo ustawieniu ogniska bez kliknięcia i bez osobnego zamykania, prowadzenie ogniska przyciskiem z odpowiedzią na każde naciśnięcie oraz treść w atrybucie opisującym, dzięki czemu dymek nie potrzebuje identyfikatora i nie zderza się między oknami. Znak Research różni się od biblioteki wyglądem — jest pierścieniem o średnicy czternaście do szesnastu pikseli ze wskaźnikiem pomocy, a nie kwadratowym przyciskiem biblioteki — dlatego klasa znaku zastępuje klasę biblioteczną, a klasa powłoki dokłada się do klasy dymka.

## budowa/klient-poprzedni/src/sterowanie/zrodlo-kanalow.ts
Zakładanie stoi w tym samym źródle co usuwanie, ponieważ świeża baza ma dokładnie jeden kanał, a okno komunikacji kanału wymaga; panel z samym usuwaniem pozwoliłby nieodwracalnie skasować jedyny kanał produktu bez drogi powrotu. Wykazu kanałów to źródło nie pobiera — rejestr kanałów jest czytającą pamięcią podręczną nad odczytem wykazu, wspólną całemu klientowi, i po udanym zapisie panel odświeża tę pamięć. Każda czynność oddaje wynik z kodem odmowy, nigdy samej treści, aby odmowa rdzenia dojechała do widoku, a nie wyglądała jak wykonanie. Pole rodzaju kanału nie występuje w opisie zmiany, ponieważ kontrakt zmiany rodzaju kanału po założeniu nie przewiduje, a formularz komunikuje to wprost zamiast przyjmować wpis, który rdzeń zignoruje. Sprawdzenie łączności jest narzędziem pomocniczym, nie bramką — jego wynik nie wstrzymuje zapisu, nie wyłącza wiersza ani nie blokuje wysłania tury, bo kanał niedostępny przed chwilą bywa sprawny teraz. Odpowiedź o stanie poświadczenia mówi, czy poświadczenie jest ustawione, jakiego jest rodzaju i gdzie mieszka, lecz nigdy nie niesie samego sekretu.

## budowa/klient-poprzedni/src/uwierzytelnienie/sesja-bramki.ts

Transportem jest jedno gniazdo WebSocket, nie seria żądań HTTP, więc ciasteczka
nie ma: token wraca z `auth.login` i tylko klient może go przechować między
uruchomieniami. Bez zapisu każdy start aplikacji wymagałby hasła, a bramka ma
stać przy wejściu raz, nie przy każdym otwarciu okna.

W zapisie leży sesja w kształcie kontraktu (`AuthSession`): token, czas
wygaśnięcia, metoda. Token jest poświadczeniem na okaziciela — rdzeń trzyma
w bazie wyłącznie jego skrót, więc jedyną kopią tokenu jest właśnie ten zapis.
`localStorage` interfejsu podawanego z `127.0.0.1` i z powłoki natywnej jest
magazynem lokalnym tej samej maszyny, czyli tym samym progiem zaufania, na
którym stoi plik sejfu rdzenia obok bazy.

Magazyny są dwa, bo wybór trwałości logowania jest jawny. Pole „Nie wyloguj
mnie" na ekranie logowania rozstrzyga, gdzie sesja wyląduje: zaznaczone kładzie
ją w `localStorage` i sesja przeżywa zamknięcie aplikacji, niezaznaczone —
w `sessionStorage`, skąd ginie razem z oknem.

Skutek jest dwustronny: ta sama wartość idzie żądaniem `auth.login` jako
`keepSignedIn` i rozstrzyga w rdzeniu o trwaniu sesji — 365 dni zamiast doby
roboczej. Trwanie jest zapisane przy wierszu sesji, więc `auth.token.refresh`
go nie ścina. Zapis w przeglądarce bez tego byłby obietnicą, której rdzeń nie
dotrzymuje.

Odczyt pyta obu magazynów, w kolejności od trwałego. Inaczej zmiana
rozstrzygnięcia między jednym a drugim uruchomieniem zostawiałaby sesję
niewidoczną, a wyglądałoby to jak jej wygaśnięcie. Kasowanie czyści oba
z tego samego powodu.

Żaden błąd pamięci nie zatrzymuje uruchomienia: brak magazynu, zapis nieczytelny
i kształt spoza kontraktu znaczą to samo — sesji zapisanej nie ma i wejście
wymaga hasła. O tym, czy sesja jeszcze żyje, rozstrzyga wyłącznie rdzeń
odpowiedzią na `auth.token.refresh`; ten plik nie porównuje czasów i niczego
nie unieważnia sam.

Zapisana sesja leży wyłącznie w jednym magazynie naraz: pole „Nie wyloguj mnie"
nie pamięta własnego zaznaczenia osobnym zapisem, tylko czyta miejsce, w którym
sesja naprawdę leży. Osobny zapis rozjechałby się z magazynem przy pierwszym
czyszczeniu pamięci i pokazywałby zaznaczone pole nad sesją, która ginie
z oknem.

Czas ważności idzie z odpowiedzi rdzenia (`expiresAt`), nie ze stałej klienta —
długość życia sesji zna wyłącznie rdzeń. Data pełna pojawia się tylko wtedy,
gdy wygaśnięcie wypada innego dnia; w dniu bieżącym wystarcza godzina.

## budowa/klient-poprzedni/src/rozmowa/blok-zwijany.ts
Zwinięcie bloku opiera się na natywnym zachowaniu przeglądarki dla elementów details i summary, dzięki czemu blok działa klawiaturą, ma poprawną semantykę dla technologii wspomagających i nie wymaga obsługiwaczy zdarzeń. Brak treści bloku prowadzi do jego ukrycia, a nie wyszarzenia.

Sterowanie rozwinięciem z zewnątrz jest wołane wyłącznie przy zmianie trybu widoku, nigdy przy odświeżeniu treści: blok przestawiony ręcznie pozostaje w stanie, w jakim go zostawiono, także gdy w trakcie tury przychodzą kolejne fragmenty.

## budowa/klient-poprzedni/src/moduly/studio/stan-okna-studio.ts
Rozróżnienie czterech faz jest treścią: stan sprzed pierwszego wywołania, stan
w trakcie wywołania i stan odmowy rdzenia są osobnymi fazami, ponieważ zlanie
ich zmusza czytającego do zgadywania, czy czekać, czy działać. W Studio różnica
jest ostra, bo rdzeń oddaje tu i treść, i pustkę, i odmowę: otwarcie dokumentu
ze ścieżką oddaje dokument pusty, wykaz repozytorium przed pierwszym zapisem
oddaje wykaz pusty, porównanie różnic bez wskazanej strony oddaje pustkę,
a otwarcie nieznanego dokumentu odmawia kodem błędu. Pusty widok bez powodu
byłby nie do odróżnienia od żadnej z tych czterech sytuacji.

Komunikat nie kasuje treści, tylko ją przesłania; wyjątkiem jest faza
ładowania, w której treść ustępuje miejsca wskaźnikowi odczytu. Stan pusty ma
jedną formę złożoną ze znaku, tytułu i opisu; znak dochodzi z zestawu ikon,
a nie z własnego rysunku, bo forma stanu pustego wariantów stylu nie ma.

Zestaw faz i ich znakowanie pochodzą ze wspólnego modułu faz okna: ta sama
nazwa fazy trafia do atrybutu danych w każdym module, więc wspólny selektor
i wspólny sprawdzian mają się o co oprzeć. W tym module zostaje wyłącznie to,
czym Studio różni się świadomie: wskaźnik odczytu i chowanie treści.

## budowa/klient-poprzedni/src/moduly/research/formularz-zrodla.ts
Formularz wydzielono z okna zarządzania źródłami, ponieważ katalogowanie źródła i ocena jego wiarygodności to osobne funkcje operatora z własnymi polami w kontrakcie (rodzaj, pochodzenie, adres, dokument repozytorium, ocena wiarygodności). Każde pole niesie dymek objaśnienia, bo każde jest elementem konfiguracji zlecenia. Rodzaj źródła i ocena wiarygodności idą przez wspólny wybór nastawy — obsadę bibliotecznego menu drzewa, gdzie uchwyt pokazuje wartość bieżącą, a natywna lista rozwijana pokazywałaby ją dopiero po rozwinięciu; pola tekstowe zostają polami tekstowymi, bo nie ma w nich czego rozwijać.
Rzutowanie wartości sterowania na wyliczenie kontraktu odbywa się wyłącznie w funkcji pomocniczej: ster oddaje napis, którego pozycje pochodzą z wyliczeń kontraktu, więc innej wartości wydać nie może, a pusty wybór wraca wartością zastępczą, nie napisem pustym.

## budowa/klient-poprzedni/src/moduly/wiedza/indeks.ts
Katalog wiedzy stoi osobno od katalogu poczty, bo to dwie różne dziedziny rdzenia: poczta sięga
po skrzynkę Operatora stojącą poza urządzeniem, wiedza po treści własne leżące w rdzeniu. Wspólny
katalog zlepiłby dwa słowniki pojęć w jeden. Rejestr modułów wiąże widok z kodem modułu rdzenia,
a tabela modułów nie niesie osobnego kodu ani dla poczty, ani dla wiedzy — drugi wpis rejestru
czekałby na kod, którego nawigacja nigdy nie poda. Okno wchodzi do sceny przez złożenie modułu
poczty, które ma w rejestrze jeden wpis. Dlatego ten plik nie wystawia stałej modułu.

## budowa/klient-poprzedni/src/moduly/translate/wiersz-terminu.ts
Jeden wiersz przypada na plik, tak samo jak wiersz biblioteki ekspertów w module agentów. Ten sam
wiersz nadaje się zarówno do wykazu terminów, jak i do zestawienia wystąpień, ponieważ nie wywołuje
działań samodzielnie, tylko oddaje naciśnięcia oknu, które wie, co z nimi zrobić.

## budowa/klient-poprzedni/src/moduly/roundtable/zestawienie-udzialu.ts
Moduł mierzy wyłącznie udział: liczbę wypowiedzi i znaków na uczestnika w turze bieżącej. Okno opisuje głosowanie, ocenę rubrykową, ocenę modelami-sędziami i ranking akumulowany między sesjami; odpowiadające im komendy kontraktu (roundtable.vote.start, roundtable.vote.cast, roundtable.vote.get, roundtable.rubric.set, roundtable.judge.run, roundtable.leaderboard.get) istnieją, ale moduł żadnej z nich jeszcze nie wywołuje. Wynik liczony jako udział nie jest oceną jakości wypowiedzi ani progiem quorum — próg zgody niesie pole quorum żądania roundtable.vote.start, którego okno nie wywołuje.
Wypowiedź otwartą ostatnio liczy się z tekstu widocznego w oknie (strumień), nie z pola content: rdzeń rozgłasza wypowiedź najpierw pustą i dopisuje treść strumieniem, więc sam zapis dawałby zero znaków przez cały czas mówienia modelu. Wypowiedzi wcześniejsze są już utrwalone i liczą się z zapisu.

## budowa/klient-poprzedni/src/okna-pomocnicze/indeks.ts
Plik był punktem zbiorczym eksportów całego obszaru okien pomocniczych, więc pierwotny
nagłówek opisywał po kolei skład obszaru, powód wydzielenia go poza katalog modułów
oraz ryzyko rozejścia się kopii przy duplikacji w dwóch modułach. Redakcja scaliła te
trzy wątki w jedno zdanie nazywające zawartość i wspólność dla modułów Developer
i Diagnostics, zachowując wymienione elementy obszaru bez warstwy uzasadniającej wybór
lokalizacji pliku w drzewie źródeł.

## budowa/klient-poprzedni/src/moduly/research/indeks.ts
Powłoka zna stąd jedną rzecz — opis modułu. Kod modułu stoi wewnątrz opisu, bo moduł sam mówi, którym jest modułem, więc rozjazd między nazwą w rejestrze a rzeczywistością staje się niemożliwy. Moduł nie osadza się sam w dokumencie i nie zna powłoki: oddaje element, a warstwa składająca decyduje, gdzie go postawić, dzięki czemu te same okna wchodzą i w obszar roboczy powłoki, i w stanowisko sprawdzianu. Funkcja zamknięcia podpina rozłączenie modułu, bez którego subskrypcja stanu badania i kreator raportu żyłyby dalej po zejściu modułu ze sceny; pole jest w opisie widoku modułu opcjonalne, więc kompilator braku nie zgłosi.

## budowa/klient-poprzedni/src/okna-rownolegle/uchwyt-szerokosci.ts
Obsługa klawiaturą jest obowiązkowa, ponieważ uchwyt osiągalny wyłącznie
wskaźnikiem odcinałby od jedynego sposobu ustawienia szerokości kolumny; stąd
rola separatora, indeks tabulacji, strzałki, Home/End i pełny komplet
atrybutów wartości ARIA. Uchwyt nie zna paneli ani rozmowy — oddaje wołającemu
jedną liczbę w pikselach bez własnych założeń, a przycięcie do minimów
wykonuje osobny moduł odpowiedzialny za szerokości gniazda. Uchwyt stoi po
lewej stronie sterowanej kolumny zgodnie z układem gniazda: rozmowa, uchwyt,
panele. Ruch w lewo poszerza kolumnę, ruch w prawo ją zwęża; strzałki
klawiatury podążają tą samą logiką kierunku, żeby obsługa ręką i klawiaturą
dawały ten sam skutek. Krok strzałki wynosi wielokrotność czterech pikseli
siatki interfejsu — krok jednopikselowy wymagałby wielokrotnego naciskania
strzałki, a stupikselowy uniemożliwiałby dojście do wartości osiąganej
wskaźnikiem bez wysiłku.

## budowa/klient-poprzedni/src/moduly/studio/stan-studio.ts
Pięć okien modułu Studio — edytor, panel narzędzi, panel różnic i wyszukiwania,
podgląd oraz repozytorium sesji — pracuje na tym samym dokumencie czynnym.
Gdyby każde okno prowadziło własną kopię, przywrócenie wersji w repozytorium
nie przestawiłoby treści edytora, a różnica dotyczyłaby innego dokumentu niż
podgląd. Poza własnym działaniem okien jedynym źródłem odświeżenia stanu jest
zdarzenie zmiany dokumentu — zmiana dokonana w innym oknie albo na innym
urządzeniu konta dociera nim, bez odpytywania w pętli.

Zakres wysyłany do rdzenia operacją kontekstową bierze wskaźnik panelu, pasek
zaznaczenia edytora i żądanie operacji z tej samej drogi, żeby istniała jedna
nastawa i jedna prawda o niej. Powodu odmowy ostatniej operacji kontekstowej
pytają okna stojące obok panelu narzędzi, żeby brak propozycji po odmowie nie
wyglądał u nich jak sytuacja sprzed pierwszej operacji.

Parę wersji do porównania wskazuje repozytorium sesji przyciskiem porównania
przy wersji, a czyta panel różnic przy składaniu żądania; nastawa jest jedna
i idzie przez stan współdzielony, ponieważ dotyczy dwóch okien naraz.

Migawkę pól dokumentu czynnego bierze okno pracy przy przełączaniu zakładek:
dwa dokumenty naraz mają niezależne zaznaczenie, bufor i propozycję, a stan
modułu prowadzi jeden dokument czynny. Migawka jest jedyną drogą do tego —
druga kopia stanu modułu rozjechałaby się z nim przy pierwszym zdarzeniu
rdzenia. Nadejście wyniku operacji unieważnia odmowę poprzednią, ponieważ
komunikat o odmowie i komunikat o wyniku nie mogą stać obok siebie prawdziwe.

## budowa/klient-poprzedni/src/rozmowa/etykiety-rozmowy.ts
Napis nigdy nie stoi obok samej barwy: każdy stan niesie słowo, bo stan sygnalizowany wyłącznie kolorem jest niedopuszczalny.

Podpowiedź pod polem wpisywania mówi wprost, że pole jest zarazem filtrem wykazu, i wymienia klawisze, które przy otwartym wykazie znaczą co innego niż zwykle.

Napis wykazu wytworów mówi „narzędzia", a nie „pliki": rdzeń nie nadaje spisu wytworzonych plików, więc wykaz wymienia nazwy wywołanych narzędzi.

Opis stanu pustego trybu Streszczenie tłumaczy, czym streszczenie jest i skąd się bierze, zamiast meldować brak danych.

Odpowiedź na Zatrzymaj bez tury dotyczy przypadku, w którym rdzeń oddaje pole stopped równe fałsz, a kliknięcie i tak musi zostawić ślad na ekranie.

## budowa/klient-poprzedni/src/strona-glowna/czynnosci-sesji.ts
Plik nie zna kanału ani kontraktu — wykonanie czynności podaje wpięcie. Warunki dostępności czynności stoją w tym pliku, a nie w karcie czy w menu, ponieważ karta buduje wiersz, menu buduje przyciski, i żadne z nich nie rozstrzyga, czy na przykład wznowienie sesji czynnej ma sens; jedno miejsce daje jedną odpowiedź i jedno miejsce do poprawienia, gdy rdzeń zmieni stany sesji. Brak czynności zdejmuje przycisk z menu, zamiast go wyszarzać, ponieważ przycisk widoczny, którego naciśnięcie nic nie robi, jest atrapą — wiersz bez czynności pozostaje czysto informacyjny. Zatrzymanie tur ma sens wyłącznie przy turze w biegu, bo zero okien strumieniujących znaczy, że komenda nie miałaby czego zatrzymać. Wznowienie dotyczy sesji, która nie biegnie — wstrzymanej albo zakończonej — ponieważ sesja czynna nie ma czego wznawiać.

## budowa/klient-poprzedni/src/okna-rownolegle/uwaga-podgladu.ts
Komunikat jest blokiem pełnej szerokości rodzica, oznaczonym wyłącznie cienką
kreską po lewej krawędzi, bez tła i obwódki, żeby nie czytał się jako druga
karta. Forma stoi w arkuszu stylów tego widoku pod nazwą własną, ponieważ
biblioteka współdzielona nie niesie gotowej klasy komunikatu blokowego.
Komunikat pozostaje widoczny na scenie, dopóki nie ustanie jego przyczyna —
brak komendy przekazania w kontrakcie — dzięki czemu odpowiada na pytanie
o brak działania także po czasie.

## budowa/klient-poprzedni/src/moduly/research/kody-okien.ts
Kody okien trzymane w jednym miejscu sprawiają, że dopisanie okna po stronie rdzenia rozjeżdża się z klientem widocznie, a nie w siedmiu plikach naraz. Dwa kody, discovery-panel i reading-view, nie mają jeszcze wiersza w katalogu rdzenia — migracje 030 i 031 zakładają dla modułu Research pięć okien, a opracowanie modułu wylicza siedem własnych obok dwóch wspólnych platformie. Kody stoją w brzmieniu, w jakim mają wejść do katalogu, dzięki czemu atrybut okna obu okien jest już dziś tym samym napisem, którym będzie po dobudowie rdzenia, i nawigacja wewnątrzmodułowa nie wymaga zmiany.
Wykaz okien spoza katalogu jest zaporą, nie zgodą: okna są zbudowane i czynne, a ich kody wchodzą do katalogu akcji tak samo jak kody okien znanych. Rdzeń odmówi im na katalogu okien, zanim dojdzie do katalogu akcji, i nazywa okno operatorowi wprost, zamiast pokazywać puste miejsce. Wykaz znika, gdy migracje 030 i 031 dostaną wiersze obu okien.

## budowa/klient-poprzedni/src/moduly/translate/zrodlo-okna-translate.ts
Komendy wykazu okien sesji i odczytu parametrów wykonania wymagają identyfikatora okna, który
bierze się z rdzenia, nigdy z literału po stronie klienta. Idą zwykłą drogą wywołania, nie drogą
odmowy, ponieważ rdzeń je obsługuje i odpowiedź przychodzi kopertą ze statusem, którą korelacja
rozpoznaje.

## budowa/klient-poprzedni/src/strona-glowna/etykiety-sesji.ts
Żadna funkcja tego pliku nie wymyśla danych: każda przekłada wartość oddaną przez rdzeń, a wartość nieobecną oddaje jako `undefined`, dzięki czemu karta pomija cały fragment zamiast pokazać wartość zmyśloną. Nazw środowisk plik nie niesie — daje je osobne źródło zasilane wykazem z rdzenia, żeby na jednym ekranie nie stały dwa źródła tej samej nazwy.

## budowa/klient-poprzedni/src/moduly/research/kreator-raportu.ts
Kreator raportu ma stany zamknięty, otwierający się, otwarty, próbę z brakami, ładowanie i błąd. Wzorzec czterech dróg zamknięcia (kontrolka w nagłówku, klawisz Escape, kliknięcie w nakładkę, akcja w stopce) powtarza wzorzec okna konfiguracji. Przycisk główny jest czynny od otwarcia; próba z brakami nie gasi kontrolki — kreator zostaje otwarty, a brak wraca komunikatem nad stopką, a ponowne naciśnięcie w trakcie wywołania jest bezpieczne po stronie logiki kreatora, nie przez odebranie klikalności.
Reguła wyboru drogi budowy raportu: żądanie z polem sekcji składa raport z samych sekcji podanych i pomija identyfikatory ustaleń; bez tego pola woła kanał modelu po sekcję nadrzędną. Reguła stoi w kreatorze, bo tu Operator rozstrzyga ją każdym wpisanym znakiem. Droga modelu nie gwarantuje odpowiedzi modelu — bez logowania budowa kończy się powodzeniem, raport zostaje zapisany, a treścią sekcji nadrzędnej jest komunikat procesu kanału; kreator uprzedza o tym przed naciśnięciem, bo po nim klient nie ma już czym rozpoznać pochodzenia treści.

## budowa/klient-poprzedni/src/rozmowa/historia-tur.ts
Komenda message.send zwraca wiadomość użytkownika, a identyfikator odpowiedzi znany jest dopiero z pierwszego fragmentu strumienia. Dlatego wpis odpowiedzi powstaje od razu po wysłaniu, jako wpis oczekujący, a pierwszy fragment go przejmuje zamiast zakładać drugi. Bez tego okno stałoby puste do nadejścia pierwszego znaku, a potem pokazało dwa wpisy zamiast jednego.

Wiązanie wypowiedzi z wiadomością rdzenia dotyczy sytuacji, w której wypowiedź trafia do historii natychmiast, przed komendą, więc nie ma jeszcze identyfikatora wiadomości: message.send oddaje go dopiero w odpowiedzi. Rdzeń rozgłasza tę samą wiadomość zdarzeniem message.changed do wszystkich połączeń konta, także do tego, które ją wysłało. Ponieważ okno pokazuje wypowiedzi roli user, to samo zdanie weszłoby do wątku dwa razy: raz jako echo miejscowe, raz ze zdarzenia. Wiązanie idzie po treści, nie po identyfikatorze, bo identyfikatora w tej chwili nie ma po żadnej stronie: dopasowuje pierwszą niezwiązaną wypowiedź o tej samej treści i nadaje jej identyfikator z rdzenia. Metoda zwraca wpis związany albo wartość pustą, gdy wypowiedzi o tej treści w historii nie ma, czyli gdy wiadomość naprawdę przyszła skądinąd.

Czyszczenie historii nie zeruje licznika kluczy. Klucz jest tożsamością pozycji w widoku; gdyby po wyczyszczeniu zaczął się od nowa, pierwszy wpis nowej rozmowy trafiłby w węzeł pozostały po wpisie rozmowy poprzedniej i zamiast nowej pozycji stanęłaby podmieniona stara.

## budowa/klient-poprzedni/src/moduly/translate/odmowa-translate.ts
Dwie komendy modułu wykonuje rdzeń modelem: przekład tekstu źródłowego na język panelu i rozpoznanie
języka tekstu. Gdy rejestr nie ma czynnego kanału, rdzeń odmawia obu z jednym kodem niedostępności
kanału. Sama wiadomość rdzenia jest powodem, ale nie mówi Operatorowi, co ma zrobić, a jej brzmienie
sugeruje usterkę wewnętrzną tam, gdzie stoi zwykły brak włączonego kanału — zdanie o brakującej rzeczy
dokłada klient, idąc za powodem rdzenia, nie zastępując go. Zdania nie dokłada wspólny tłumacz kodów,
bo ten sam kod niesie w module dwie różne odmowy: powyższą oraz odmowę importu glosariusza, gdzie
rdzeń nie czyta plików z dysku Operatora i zdanie o kanale modelu byłoby nieprawdziwe. Zdanie mówi też,
dlaczego Operator nie ma tu pola wyboru: żądania modelowe nie niosą identyfikatora kanału, więc rdzeń
bierze kanał domyślny czynny, a klient nie wymyśla pól kontraktu.

## budowa/klient-poprzedni/src/moduly/translate/pola-zrodla.ts
Skład pola — etykieta, podpowiedź, dymek, wykaz podpowiadanych wartości — to inna odpowiedzialność
niż to, co dzieje się po naciśnięciu przycisku, dlatego pola są wydzielone z okna, które składa
całość i wiąże zdarzenia. Dymek objaśnienia dostaje każde pole zmieniające treść żądania: język
źródłowy oraz ponowna segmentacja.

## budowa/klient-poprzedni/src/strona-glowna/formularz-komponentu.ts
Formularz stoi zwinięty, ponieważ środkiem ciężkości strony głównej są karty środowisk, a formularz stale rozłożony odbierałby im pas ekranu przy czynności wykonywanej rzadko; rozwija go segment „Dodaj nowy” belki strefy trzeciej. Wykaz rodzajów komponentu odmawia rodzaju bez magazynu profili, a powód stoi na ekranie, nie tylko w komentarzu, inaczej zostałaby sama krótsza lista bez wyjaśnienia. Zakładanie i zmiana komponentu idą jedną szufladą, ponieważ są jedną pracą w dwóch krokach, a osobny uchwyt do drugiego kroku kazałby szukać go po założeniu komponentu; drugi uchwyt do jednej szuflady byłby drugą drogą do jednego bytu. Doklejenie kafla z odpowiedzi zakładania pokazałoby stan, którego rdzeń nie potwierdził drugim odczytem, dlatego wykaz dociąga strefa osobnym zapytaniem.

## budowa/klient-poprzedni/src/okna-rownolegle/widok-pelnoekranowy.ts
Scena dzieli szerokość między trzy obszary: kolumnę rozmowy, kolumnę paneli
i widok pełnoekranowy. Ten trzeci jako jedyny nie ma sąsiada — bierze całą
scenę, bo panel w kolumnie trzyma minimum szerokości, a treść, która się
w nim nie mieści, potrzebuje sceny, nie kolumny. Treść jest przenoszona, nie
kopiowana, żeby uniknąć drugiej subskrypcji rdzenia dla tego samego panelu
i subskrypcji osieroconej po zamknięciu jednego z egzemplarzy; funkcja
pokazująca zapamiętuje rodzica i następnik, a funkcja chowająca odkłada treść
dokładnie tam. Klawisz Escape wychodzi z widoku dzięki nasłuchowi osadzonemu
na elemencie widoku, nie na dokumencie, dzięki czemu sprzątanie nasłuchu jest
zbędne i klawisz nie zabiera działania niczemu innemu na scenie. Wyjście nie
jest zamknięciem panelu: treść wraca do stosu, a panel żyje dalej ze swoją
subskrypcją, natomiast właściwe zamknięcie panelu należy do jego obudowy
w stosie. Widok nie buduje treści i nie rozstrzyga zachowania kolumn sceny —
to należy do układu, który go osadził.

## budowa/klient-poprzedni/src/moduly/research/modul-research.ts
Układ wynika z roli okna i z porządku pracy badawczej: wyszukaj, skataloguj, przeczytaj, odnotuj, złóż, wydaj. Research Workspace jest wiodące i jest punktem wejścia, więc stoi w pasie pierwszym na całą szerokość. Dalej idą pary, w których biegnie wiązanie: Discovery Panel obok Sources Manager, bo pozycja wyniku staje się źródłem; Reading View obok Findings Panel, bo wypis z lektury staje się ustaleniem; Report Builder obok Export Panel, bo dokument staje się plikiem. Badanie jest jedno na cały moduł: wybór źródeł, zaznaczenie ustaleń i wskazanie materiału do lektury przestawiają wszystkie siedem okien naraz, bo stan badania jest jeden. Nawigacja wewnątrzmodułowa jest przeniesieniem ogniska, nie zmianą trasy: okna stoją obok siebie, więc odnośnik do innego okna prowadzi wzrok i ognisko, zamiast wymieniać zawartość obszaru.

## budowa/klient-poprzedni/src/strona-glowna/indeks.ts
Powłoka aplikacji sięga po budowę strony głównej wyłącznie przez ten plik i nie zna podziału widoku na strefy, karty ani kafle, dzięki czemu zmiana wewnętrzna katalogu strony głównej nie dotyka niczego poza nim.

## budowa/klient-poprzedni/src/moduly/research/nasluch-odmow.ts
Odmowa nieznanej komendy wraca kopertą bez pola statusu, a klient uznaje za odpowiedź wyłącznie kopertę ze statusem. Korelacja żądania takiej koperty nie rozstrzyga, więc obietnica wywołania zostałaby nierozstrzygnięta na zawsze, a okno wisiałoby w stanie ładowania — odmowę trzeba więc rozpoznać po ładunku, nie po nazwie zdarzenia. Nasłuch wiąże identyfikator żądania z ładunkiem odmowy nieznanej komendy: zwycięża to, co przyjdzie pierwsze — odpowiedź korelacji albo odmowa nieznanej; połączenia nic to nie kosztuje, rdzeń pracuje dalej.
Kod niepowodzenia not_found mówi prawdę: bytu o tej nazwie w rdzeniu nie ma. Nazwa typu zostaje w polu osobnym, żeby okno mogło ją wypisać wprost, zamiast pokazywać pusty wykaz udający brak danych.

## budowa/klient-poprzedni/src/moduly/research/odczyt-badania.ts
Odczyt ma jedną odpowiedzialność: dwa kroki odczytu i przełożenie ich niepowodzeń na fazę pamięci. Wydzielony jest z pamięci, bo pamięć nie zna kontraktu, a odczyt zna wyłącznie kontrakt — rozdzielenie trzyma obie rzeczy w rozmiarze, w którym dają się przeczytać naraz. Kroki są dwa, bo mówią o dwóch różnych rzeczach: wykaz okien mówi, które okno sesji należy do modułu Research, a stan okna mówi, co w nim jest. Klient nie zgaduje identyfikatora okna — bez wskazania rdzenia zostaje uczciwy stan pusty nazywający brak, nie zaszyta wartość.

## budowa/klient-poprzedni/src/uwierzytelnienie/tozsamosc-urzadzenia.ts

PIN jest właściwy maszynie, nie osobie wchodzącej: rdzeń odnajduje metodę po
parze rodzaj i urządzenie, więc PIN założony pod jednym identyfikatorem nie
otworzy bramki podanym pod innym. Ekran logowania musi znać tę wartość, żeby
wiedzieć, czy segment „PIN" ma w ogóle stanąć, i żeby móc PIN wysłać.

Bramka wyłącznie czyta identyfikator; nadaje go zakładanie PIN-u w Ustawieniach,
czyli czynność wykonana na urządzeniu. Ekran logowania staje przy każdym
uruchomieniu, także pierwszym na świeżej przeglądarce — gdyby nadawał
tożsamość, zapisywałby maszynę, która żadnego PIN-u nie ma. Pusty wynik jest
prawdą: ta przeglądarka jeszcze się nie przedstawiła.

Ten sam klucz niesie plik tożsamości urządzenia w module ustawień i tak samo
rozdziela odczyt od nadania. Bramka nie importuje pliku ustawień, bo warstwa
logowania leży niżej niż okno ustawień i zależność szłaby pod prąd warstw;
scalenie obu plików wymaga warstwy wspólnej poniżej nich.

## budowa/klient-poprzedni/src/okna-rownolegle/wiersz-czynnosci.ts
Wiersz czynności stoi osobno od wiersza panelu, bo wiersz panelu jest
przełącznikiem — niesie stan zaznaczenia i stałe miejsce na ptaszek, ponieważ
panel bywa otwarty albo zamknięty, natomiast czynność sesji nie ma stanu
włączenia, więc ptaszek kłamałby o tym, czym pozycja jest. Układ pozostaje ten
sam co w sekcji paneli: ikona, nazwa i przeznaczenie, skrót wyrównany do
prawej, żeby dwie sekcje jednego menu nie miały dwóch rytmów. Skrót
klawiszowy jest tu prawdziwy, bo nasłuch stoi na sekcji i łapie literę, dopóki
menu jest rozwinięte; skrótu globalnego celowo nie rejestruje się, żeby litery
wciśnięte w polu wypowiedzi pisały tekst, a nie wykonywały czynność. Naciśnięcie
wiersza nie zwija menu — tak samo jak przy przełączaniu paneli — a modal
pytania o nazwę oraz modal potwierdzenia usunięcia są natywnymi oknami
z własną nakładką, stającymi nad menu, więc nie ma czego chować. Wyróżnienie
ostrzegawcze jest wyłącznie barwą wiersza, nie blokadą: czynność usuwająca
zapis bez odwrotu idzie czerwienią, ale wiersz pozostaje tak samo klikalny jak
każdy inny.

## budowa/klient-poprzedni/src/rozmowa/historia-z-rdzenia.ts
Komenda message.list daje wiadomości leżące w bazie, więc odświeżenie okna albo powrót do sesji pokazuje wątek, a nie pustkę. Wpisy odtworzone są domknięte: pochodzą z zapisu, nie ze strumienia, więc nie mają tury w biegu ani fragmentów do doliczenia. Wywołujący dostaje dwie drogi, przyjmij na sukces także z pustą tablicą, i zglosNiepowodzenie na odmowę z gotowym zdaniem złożonym z powodu koperty, i ma obsłużyć obie. Odmowa rdzenia i historia naprawdę pusta muszą wyglądać na ekranie inaczej.

Treść wpisu bierze się z pola content; wiadomość przerwana albo błędna wraca z treścią, którą zdążyła zebrać, bo pokazanie jej jest uczciwsze niż ukrycie tury, która się odbyła. Stan wpisu bierze się z pola status wiadomości, nie z samego faktu zapisu: tura przerwana i tura zamknięta błędem mają po odświeżeniu okna wyglądać tak samo jak w chwili, w której się odbyły.

Odczyt wykazu blocks jest tolerancyjny wzorem odczytu fragmentu: brak obszaru, obszar cudzego kształtu albo pozycja bez rodzaju dają mniej bloków, bo historia ma się wyświetlić, a nie zniknąć od nieczytelnej pozycji.

Wypowiedź operatora nie ma stanu tury: jest tym, co napisał, i nie może wrócić z zapisu jako zamknięta bez odpowiedzi, nawet gdyby zapis był pusty.

## budowa/klient-poprzedni/src/moduly/studio/statystyka-roznicy.ts
Panel różnic i wyszukiwania żąda paska statystyki oraz filtrowania różnic po
rodzaju. Kontrakt nie niesie osobnej komendy liczącej statystykę, ponieważ
porównanie stron oddaje fragmenty wraz z treścią przed i po, więc liczby
wynikają z odpowiedzi, którą panel już ma; drugie wywołanie po te same liczby
byłoby pytaniem o coś, co leży na stole. Moduł nie zna struktury dokumentu:
wejściem są fragmenty kontraktu, wyjściem liczby i wykaz przefiltrowany,
dzięki czemu rachunek sprawdza się bez stawiania okna.

## budowa/klient-poprzedni/src/strona-glowna/kafel-komponentu.ts
Kafel dzieli kartę ze środowiskiem, o mniejszej wadze: wariant komponentu odbiera wstęgę i cień sygnału w spoczynku, zmniejsza skalę ikony i zapisuje etykietę krojem bazowym półgrubym, nie nagłówkowym. Różnica krojów niesie znaczenie: krój nagłówkowy oznacza wejście do środowiska, bazowy oznacza zbudowanie komponentu, dlatego kafel nie dostaje klasy akcentu, którą nosi karta środowiska. Wiersz pusty metadanych wyglądałby jak metadane, których nie odczytano, dlatego kafel rodzaju wiersza w ogóle nie dostaje. Czas formatuje widok, bo rdzeń nie zna strefy czasowej Operatora; data pojawia się dopiero, gdy zmiana wypadła innego dnia, a w dniu bieżącym wystarcza sama godzina.

## budowa/klient-poprzedni/src/moduly/workspace/warstwy-instrukcji.ts
Warstwy instrukcji stoją w osobnym pliku od okna, bo to inna odpowiedzialność: okno prowadzi zapis
i odczyt, tu leży wyłącznie przedstawienie warstwy. Wykaz poziomów i ich nazwy pochodzą z modułu
poziomów zasięgu.

## budowa/klient-poprzedni/src/widok-sterowania/indeks.ts

Katalog nie buduje ani jednej kontrolki, tylko oprawę, w której kontrolki
katalogu stają się widoczne. Punkt wejścia klienta montuje widok dla okna
potwierdzonego przez rdzeń, tworząc rejestr kanałów, odświeżając go i wołając
montaż widoku sterowania z kanałem, oknem, rejestrem kanałów, panelem oraz
akcjami paska. Widok powstaje raz na okno; rejestr kanałów raz na klienta,
bo jest katalogiem wyboru, nie ustawieniem okna.

## budowa/klient-poprzedni/src/okna-rownolegle/wiez-koordynacji.ts
Ustalanie więzi wiąże pierwszy widoczny koordynator z pierwszym widocznym
wykonawcą; trzecie okno w roli wykonawcy pozostaje poza tą więzią, ponieważ
pas relacji pokazuje jedno powiązanie naraz, żeby obraz dał się ogarnąć
wzrokiem. Koordynator niesie napis informujący o zleceniu do wykonawcy,
a wykonawca napis o zleceniach otrzymanych od koordynatora; groty obu napisów
wskazują tę samą stronę sceny — kierunek biegu zlecenia — i stoją po stronie
okna partnera. U koordynatora grot wychodzi z napisu ku wykonawcy, u wykonawcy
wchodzi w napis od strony koordynatora, dzięki czemu kierunek pętli czyta się
bez czytania słów; gniazdo spoza więzi nie dostaje napisu.

## budowa/klient-poprzedni/src/moduly/translate/porownanie-paneli.ts
Zestawienie jest czynnością czysto miejscową: niczego nie liczy i o nic nie pyta rdzenia. Panele
są dokładnie dwa; przy innej liczbie wskazań zestawienie mówi wprost, czego oczekuje, zamiast
pokazać cokolwiek.

## budowa/klient-poprzedni/src/moduly/research/okno-discovery-panel.ts
Okno jest zbudowane, choć katalog okien rdzenia jeszcze go nie zna, bo trzon jego pracy ma dziś pokrycie w kontrakcie: dwa z czterech trybów zapytania idą własną komendą, a przeniesienie pozycji do katalogu idzie tą samą drogą co formularz Sources Manager. Ster trybu niesie wszystkie cztery pozycje, także te dwie, których rdzeń jeszcze nie obsługuje — wybór takiego trybu niczego nie gasi, zapytanie wychodzi pod nazwą swojej komendy, a rdzeń odmawia własnymi słowami; ukrycie tych pozycji zataiłoby zakres modułu, a ich wygaszenie łamałoby zasadę zero blokad. Plik składa widok; zachowanie po naciśnięciu leży w module obsługującym czynności odkrywania.
Tryb semantyczny szuka po znaczeniu w bibliotece wiedzy, historii rozmów i przestrzeni roboczej okna badania komendą wyszukiwania wiedzy. Tryb pełnotekstowy szuka po treści w zasobach repozytorium komendą wyszukiwania pliku. Tryb webowy i tryb naukowy dzielą jedną komendę różniącą się wartością pola trybu; rdzeń nie ma dla nich jeszcze uchwytu, więc zapytanie wraca odmową.

## budowa/klient-poprzedni/src/widok-sterowania/montaz-widoku.ts

Kontrolki nie powstają w tym pliku — pochodzą w całości z katalogu sterowania:
środowisko wykonania, host, moduł, model, model zapasowy, nakład rozumowania,
tryb uprawnień, rola okna, katalogi robocze. Ten plik dokłada im miejsce,
nagłówek, podsumowanie wartości i dwa uchwyty rozwijania. Kolumna montuje się
obok sceny okna, a nie w prawym rogu paska górnego — dziewięć kontrolek nie ma
jak się tam zmieścić. Widok powstaje osobno dla każdego okna i domyka się na
jego identyfikatorze, więc dwa okna obok siebie dostają dwie niezależne kolumny
bez ani jednej wspólnej zmiennej. Szuflada startuje rozwinięta, żeby komplet
ustawień był widoczny bez szukania, co nacisnąć.

## budowa/klient-poprzedni/src/strona-glowna/karta-sesji.ts
Każdy napis karty pochodzi z rdzenia: tytuł to tytuł sesji, a gdy rdzeń go nie nadał — identyfikator sesji, nigdy nazwa wymyślona; środowisko, moduł i liczby okien niesie żywy odpis obecności, a przy odpisie nieobecnym cały fragment znika, zamiast pokazywać zero udające odczyt. Przycisk powrotu do sesji pojawia się wyłącznie wtedy, gdy montaż podał czynność powrotu; bez tej czynności wiersz jest czysto informacyjny, bez przycisku wyszarzonego ani martwego. Wykaz środowisk strony jest tym samym bytem, z którego rysuje się karta środowiska w strefie pierwszej — nazwa wzięta z klienckiej stałej rozjeżdżałaby się po cichu z nazwą zapisaną w bazie. Ponowne naciśnięcie przycisku powrotu w trakcie jego zajętości nie dubluje powiązania z sesją.

## budowa/klient-poprzedni/src/widok-sterowania/naglowek-widoku.ts

Tytuł jest jedynym miejscem kroju szeryfowego w tym widoku — krój ten należy
wyłącznie do nagłówków. Rola okna stoi w nagłówku, a nie tylko w podsumowaniu,
ponieważ to ona rozstrzyga, czym okno jest w pętli koordynator–wykonawca: przy
dwóch oknach obok siebie rola musi być widoczna bez rozwijania czegokolwiek.
Plakietka niesie ikonę słowną — nazwę roli — więc stan nie opiera się na samej
barwie.

## budowa/klient-poprzedni/src/moduly/translate/zrodlo-dokumentu-translate.ts
Obszar translate nie ma komendy przyjmującej plik: tekst źródłowy wchodzi napisem. Wydobycie tekstu
z dokumentu i rozpoznanie pisma ze skanu leżą w obszarze document i są jedynym wejściem od strony
pliku, więc moduł woła je wprost, tak jak woła wykaz okien i wykaz kanałów, zamiast je kopiować.
Wywołanie idzie zwykłą drogą protokołu, nie ścieżką odmowy właściwą wyłącznie obszarowi modułu,
bo odmowa obszaru dokumentów przychodzi kopertą ze statusem, którą korelacja rozpoznaje bez pomocy.
Ścieżka pliku jest ścieżką po stronie rdzenia — klient dysku nie czyta ani nie zapisuje, tylko podaje
wskazanie i oddaje Operatorowi odpowiedź rdzenia. Brak pola wymuszenia rozpoznania pisma znaczy, że
rdzeń bierze warstwę tekstową dokumentu, gdy ją ma — to zachowanie domyślne kontraktu, nie wybór okna.

## budowa/klient-poprzedni/src/rozmowa/metadane-konta.ts
Kanał nadaje fragment tego rodzaju także w chwili przełączenia konta w trakcie tury, żeby rotacja była widoczna. Struktura niesie wyłącznie odwołania: kod konta, nazwę profilu, odwołanie do danych dostępowych, a treść poświadczenia nie trafia tu nigdy. Kształt odpowiada strukturze metadanych konta po stronie rdzenia.

## budowa/klient-poprzedni/src/rozmowa/modul-okna.ts
Moduł okna jest czytany z rdzenia, a nie zgadywany po stronie widoku, dwiema drogami z kontraktu. Komenda window.state.get przy złożeniu pozwala oknu zapytać, w czym pracuje, zamiast czekać na pierwszą zmianę; bez tego okno wznowione po odświeżeniu klienta stałoby na module nieustalonym mimo znanego stanu w rdzeniu. Zdarzenie window.changed w toku pracy odzwierciedla sytuację, w której polecenie przestawienia obszaru roboczego przestawia to samo okno na inny moduł zamiast je zamykać, a rdzeń rozgłasza zmianę zdarzeniem. Niepowodzenie zapytania nie jest błędem okna: moduł zostaje nieustalony, a okno pracuje na arsenale wspólnym.

## budowa/klient-poprzedni/src/rozmowa/indeks.ts
Złożenie warstwy rozmowy w punkcie wejścia sprowadza się do wywołania zamontujRozmowe wewnątrz uzgodnienia otwarcia okna, z przekazaniem sceny, kanału, identyfikatora okna oraz persony i roli okna. Nazwy komend i zdarzeń, których warstwa używa — wysłania wiadomości, zatrzymania wiadomości, fragmentu strumienia, zmiany wiadomości, zmiany okna, odczytu stanu okna oraz wykazu akcji — pochodzą wyłącznie ze wspólnego kontraktu. Okno przestawia się na moduł samo: śledzi moduł okna w rdzeniu i przy jego zmianie rekonfiguruje pasek narzędzi promptu, panel akcji i panel kontekstu, zachowując wątek rozmowy. Powłoka może też przestawić okno wprost przez metodę ustawModul zamontowanej rozmowy, gdy zna wynik polecenia wejścia w obszar roboczy.

## budowa/klient-poprzedni/src/moduly/research/okno-export-panel.ts
Export Panel niesie dwie funkcje Operatora: eksport raportu oraz wybór formatu wyjściowego i miejsca docelowego. Plik składa widok; zachowanie po naciśnięciu leży w module obsługującym czynności eksportu. Format wyjściowy idzie sterem nastawy: uchwyt niesie wartość bieżącą, a nie nazwę rodzajową, bo natywna lista pokazuje ją dopiero po rozwinięciu. Historia eksportów bieżącej sesji rośnie wyłącznie z odpowiedzi rdzenia, nie z zamówienia, więc wykaz mówi, co rdzeń naprawdę oddał, razem z brakiem ścieżki tam, gdzie jej nie oddał.

## budowa/klient-poprzedni/src/okna-rownolegle/wpis-sesji-okna.ts
Sekcja czynności menu gniazda rozstrzyga o pozycjach na podstawie wpisu sesji
— stanu, tytułu, projektu i liczby okien strumieniujących — a gniazdo zna
wyłącznie identyfikator sesji, więc to źródło prowadzi od identyfikatora do
pełnego wpisu. Źródło sesji w tle się tu nie nadaje, bo obsługuje sesje spoza
bieżącego połączenia, a nie tę, o którą pyta menu, mimo że korzysta z tej
samej komendy z żywym stanem. Źródło nie zna DOM-u, nie wykonuje żadnej
czynności i nie rozstrzyga, co z sesją wolno zrobić — oddaje wpis albo pustą
wartość, która oznacza brak wiedzy, nie zakaz: dopóki rdzeń nie oddał wykazu,
sekcja czynności pozostaje krótsza, bo czynności bez znanego stanu sesji nie
da się uczciwie nazwać.

## budowa/klient-poprzedni/src/moduly/translate/zrodlo-zrodla.ts
Tekst źródłowy mieszka w module stanu translate, żeby panele tłumaczenia patrzyły na ten sam tekst,
a nie na własną kopię. Wszystkie trzy komendy rdzeń rejestruje i obsługuje wprost; bez zalogowanego
modelu rdzeń oddaje stan zdegradowany z komunikatem o braku zalogowania — to jest odpowiedź rdzenia,
nie brak uchwytu. Ścieżka odmowy zostaje na wypadek starszego rdzenia albo pośrednika.

## budowa/klient-poprzedni/src/moduly/studio/strona-magazyn-widoku.ts
Moduł nastaw operatora widoku od początku przyjmował magazyn podany z zewnątrz
zamiast sięgać po magazyn przeglądarki wprost, żeby zapis dał się przełożyć na
rdzeń bez zmiany wołających go miejsc. Ten moduł jest tym przełożeniem: zapis
idzie do rdzenia komendami odczytu i zapisu widoku, więc nastawy przestają
ginąć przy przesiadce na inne urządzenie.

Magazyn nastaw widoku pozostaje z natury natychmiastowy: odczyt musi oddać
wartość w chwili składania okna, a rdzeń odpowiada później. Zapis dwustopniowy
rozwiązuje to bez kłamstwa — pamięć podręczna trzyma nastawy zapisane
w rdzeniu, a wczytanie ją napełnia; do chwili odpowiedzi obowiązuje odbicie
w magazynie przeglądarki, gdy jest dostępne, albo nastawy domyślne.

Odbicie w przeglądarce zostaje, ponieważ pierwsza klatka okna rysuje się przed
odpowiedzią rdzenia. Odbicie nie jest drugim źródłem prawdy: rdzeń nadpisuje
je przy każdym odczycie, a zapis idzie do obu naraz. Bez odbicia każde wejście
do modułu zaczynałoby się nastawami domyślnymi, choćby operator ustawił co
innego wcześniej.

Nieudany zapis nastawy widoku nie przerywa pracy, ale i nie milczy: idzie
zdaniem zwrotnym do zaplecza, więc korzystający wie, że wybór nie przeżyje
zamknięcia okna.

## budowa/klient-poprzedni/src/strona-glowna/karta-srodowiska.ts
Karta buduje treść w układzie: godło, tytuł krojem nagłówkowym, motto, jednozdaniowy opis trybu pracy krojem bazowym. W spoczynku powierzchnia jest neutralna, bez akcentu; przy najechaniu i w stanie czynnym pojawia się wstęga górna w błękicie sygnałowym, cień sygnału i uniesienie karty, z wariantu akcentu biblioteki komponentów, a arkusz strony wygasza wstęgę w spoczynku, bo biblioteka pokazuje ją stale. Sygnał jest jedyną barwą akcentu bieżącego systemu wizualnego. Nośnikiem karty jest przycisk, nie warstwa z rolą przycisku, ponieważ karta ma być celem nawigacji klawiaturą z pierwszeństwem natywnym, a pierścień fokusu wnosi wariant klikalny biblioteki komponentów.

## budowa/klient-poprzedni/src/moduly/translate/pasek-kontekstu.ts
Trzy pustki kontekstu są trzema różnymi zdaniami: jeszcze nie pytano, pytanie w toku, rdzeń nie zna
ani jednego okna tej sesji. Pierwsze każe czekać na wejście do modułu, drugie na odpowiedź, trzecie
mówi, że komendy panelu źródła nie mają czym zaadresować żądania. Przycisk ponowienia stoi tylko
tutaj, bo kontekst okna jest jedynym odczytem modułu; resztę wyzwala zapis, którego przycisk zostaje
klikalny także po odmowie. Ponowienie wyzwala odczyt i nie zdejmuje komunikatu — ten znika dopiero,
gdy odczyt się powiedzie. Tryb uprawnień, zasięg wykonania, rola okna i katalogi robocze stoją
w panelu sterowania okna na tym samym ekranie.

## budowa/klient-poprzedni/src/widok-sterowania/nazwy-ustawien.ts

Plik nie zna ani jednej nazwy prezentacyjnej z osobna — wszystkie bierze
z istniejących słowników etykiet okna, etykiet sterowania oraz rejestru
kanałów. Jedyne, co dokłada, to złożenie wartości w wiersz czytelny bez
otwierania kontrolki. Host wykonania nie jest osobnym wierszem, lecz dopiskiem
przy środowisku wykonania — uszczegóławia zasięg, a nie stanowi odrębnego
ustawienia. Nazwa hosta zostaje przez to widoczna bez dziewiątej pozycji.

Identyfikator kanału spoza wykazu nie znika i nie zostaje podmieniony — widok
pokazuje prawdę o oknie, a nie pozycję pierwszą z brzegu. Dopóki wykaz jest
pusty, bo odpowiedź rdzenia jest jeszcze w drodze, dopisek o pozycji spoza
wykazu byłby przedwczesny — wiersz pokazuje sam identyfikator ze wskaźnikiem
ładowania. Rejestr nie odróżnia braku odpowiedzi od pustego rejestru.

## budowa/klient-poprzedni/src/powloka/menu-profilu.ts

Menu profilu Operatora w pasku powłoki środowiska udostępnia siedem pozycji ustawień
pogrupowanych w trzy grupy tematyczne — tożsamość i dostęp, widok, obecność globalna —
oraz stopkę prowadzącą do pełnego rejestru ustawień. Pozycje nie są osobnym wykazem:
powstają z rejestru pozycji ustawień, z którego żyje także listwa strony głównej i menu
aplikacji. W tym module leży wyłącznie podział na grupy oraz umieszczenie pozycji
konfiguracji w stopce, bo stopka jest drogą do rejestru, czyli miejscem rzeczy
pełniejszej niż pozycje nad nią. Zmiana nazwy albo ikony w rejestrze zmienia to menu
samoczynnie. Mechanizm budowy menu pochodzi w całości z biblioteki komponentów drzewa
menu, która niesie grupowanie, nastawę dwustanową, stopkę, znacznik wyboru i wędrówkę
klawiaturą. Uchwyt menu niesie inicjały Operatora zamiast napisu rodzajowego —
pojedyncza litera zastępcza nie oznacza inicjału konkretnej osoby, tylko pierwszą literę
jedynej nazwy roli, jaką produkt zna. Motyw ma jedną prawdę: przełącznik w menu i
przycisk w pasku czytają ten sam stan obowiązującego motywu i nasłuchują tego samego
zdarzenia zmiany, więc zmiana w jednym miejscu przerysowuje drugie. Menu nie zawiera
konta, adresu ani polecenia wylogowania, bo platforma nie prowadzi kont ani profili
osobowych — bramka dostępu jest jedna, a rola Operatora pozostaje bezimienna, więc po
zalogowaniu nie ma bram do przekraczania. Zdanie graniczne tłumaczące ten brak stoi na
ekranie w stopce, tym samym brzmieniem co w menu aplikacji.

Pole szukania nad ośmioma pozycjami zabrałoby wiersz i nie skróciłoby ani jednego ruchu:
próg biblioteki, od którego pole szukania pojawia się samoczynnie, i tak by go nie
pokazał przy tej liczbie pozycji — zapis jawny mówi, że to świadomy wybór, a nie
przypadek. Znacznik wyboru pozycji w drzewie menu zostaje zawsze pusty, bo żadna
pozycja nie jest nastawą jednokrotną: wszystkie otwierają okno, więc menu nie ma czego
zaznaczać jako „ostatnio otwarte". Grupa „Widok" stoi między tożsamością a obecnością
i ma dokładnie jedną pozycję — nastawę motywu, a nie wejście do okna — dlatego składa
się osobno od pozostałych grup budowanych z rejestru.

## budowa/klient-poprzedni/src/moduly/research/okno-findings-panel.ts
Zdarzenie zmiany ustalenia jest już w kontrakcie, ale rdzeń go nie rozgłasza, więc wykaz nie dostanie ustalenia zapisanego na innym urządzeniu konta. Okno mówi o tym w stanie pustym zamiast odpytywać rdzeń w pętli — subskrypcja dojdzie razem z rozgłaszaniem. Plik składa widok; zachowanie po naciśnięciu leży w module obsługującym czynności ustaleń.
Trzecia reprezentacja wykazu, kodowanie jakościowe, nie jest tu zbudowana: kody są już w kontrakcie, ale ustalenie oddawane przez rdzeń ich nie niesie i rdzeń nie ma uchwytu książki kodów, więc widok, który nie miałby czego pokazać, tu nie stoi — pozycja jest w panelu akcji, pod nazwą swojej komendy.
Kopia przed sortowaniem w funkcji porządkującej jest konieczna: wykaz pochodzi wprost z pamięci modułu, a sortowanie w miejscu wywróciłoby porządek narastania w pamięci, wspólnej wszystkim oknom.

## budowa/klient-poprzedni/src/moduly/translate/zrodlo-glosariusza.ts
Kontrakt nie ma komendy odczytu glosariusza jako wykazu, choć edycja glosariusza wymaga wczytania.
Źródło nie zmyśla wykazu: okno pokazuje wyłącznie terminy zapisane w bieżącej sesji i mówi wprost,
czego brakuje.

## budowa/klient-poprzedni/src/okna-rownolegle/zrodlo-lacznosci.ts
Doprowadzenie stanu łącza do sceny okien równoległych ma jedną odpowiedzialność:
zamianę subskrypcji transportu na strumień odczytów, które układ rozsyła do
nagłówków gniazd. Nic tu nie jest źródłem prawdy — prawdą jest sam transport,
a ten plik wyłącznie go odpytuje. Licznik kolejki musi być prawdziwy, a nie
zamrożony: transport ogłasza połączenie przed opróżnieniem kolejki, więc odczyt
zrobiony w chwili zmiany stanu zamarzałby na wartości sprzed wysłania, dlatego
po połączeniu odczyt dobija się cyklicznie, aż kolejka spadnie do zera. Tę samą
rachubę prowadzi niezależnie pasek górny aplikacji — nie jako druga prawda,
tylko jako drugi pytający tego samego transportu, ponieważ wspólnego miejsca
dla tej logiki nie ma: pasek zwraca gotowy element, nie strumień odczytów.
Interfejs transportu wystawia sześć podstawowych metod, a numer próby
i zaplanowane ponowienie trzyma prywatnie gniazdo połączenia, więc dojście do
przebiegu ponowienia sprawdza obecność metod jawnie i jednorazowo, zamiast
zakładać ją na sztywno w miejscu wywołania.

## budowa/klient-poprzedni/src/widok-sterowania/obserwator-ustawien.ts

Podsumowanie pokazuje osiem wartości obok siebie także wtedy, gdy szuflada
z kontrolkami jest zwinięta — a kontrolki kompletu sterowania trzymają swój
stan wewnątrz panelu sterowania i nie wystawiają go na zewnątrz. Ten plik jest
odczytem tego samego stanu, a nie drugą jego definicją: buduje stan sterowania,
zmianę okna i zmianę ustawienia z katalogu sterowania i nie dopisuje ani jednej
reguły protokołu.

Obserwator wyłącznie czyta — nie wywołuje czynności zastosowania ani zapisu,
więc nie ma drogi, którą mógłby zmienić okno. Komunikaty o losie zmian pomija,
bo pokazuje je pasek kompletu sterowania; podwójny komunikat byłby szumem. Gdy
panel sterowania wystawi własną migawkę i subskrypcję zmian, ten plik zniknie,
a podsumowanie odczyta stan kompletu wprost.

## budowa/klient-poprzedni/src/rozmowa/lista-wpisow.ts
Wpis rozpoznawany jest po kluczu, nie po pozycji: fragment strumienia odświeża tę samą pozycję, zamiast dokładać kolejną. Dzięki temu tura, która przyniosła prowenancję, dziesiątki fragmentów tekstu, wywołania narzędzi i podsumowanie, pozostaje w historii jednym wpisem. Przewijanie do końca następuje tylko wtedy, gdy lista już stała na końcu, więc czytanie starszej wypowiedzi nie jest przerywane przez nadchodzący strumień. Tryb widoku transkryptu jest stanem listy, nie pojedynczego wpisu: lista trzyma jeden tryb i rozsyła go do widoków, więc wpis założony po przełączeniu rodzi się już w trybie bieżącym. Pamięć wpisów jest od trybu niezależna, przełączenie niczego nie usuwa, a powrót do trybu zwykłego przywraca wątek w całości.

Po zdjęciu wszystkich wpisów stan pusty wraca na wierzch, ale nie zostaje sam: wołający zaraz po wyczyszczeniu dopisuje zdanie o powodzie, żeby w miejscu zniknięcia rozmowy nie stał napis o rozmowie, która się jeszcze nie zaczęła.

Dwie różne pustki dostają dwa różne zdania: to, że wątek się nie zaczął, i to, że tryb Streszczenie nie ma jeszcze czego streścić, to nie ten sam fakt. Drugi wariant tłumaczy, czym streszczenie jest i skąd się bierze.

## budowa/klient-poprzedni/src/moduly/workspace/agenci-czynnosci.ts
Czynności okna zarządcy ekspertów sięgające poza sam wykaz — zmiana uprawnienia eksperta i zdanie
o przypisaniach projektu — stoją w osobnym pliku od okna, bo okno składa kontrolki i prowadzi odczyt,
a tu leży przebieg czynności wraz z tym, co po niej widzi operator, tym samym wzorem co w module
czynności biblioteki. Potwierdzenie mówi to, co zapisał rdzeń, nie to, czego żądało okno: rdzeń
oddaje po zapisie komplet uprawnień eksperta i ten komplet rozstrzyga. Po odpowiedzi okno odrysowuje
wykaz, bo stan ładowania go opróżnia.

## budowa/klient-poprzedni/src/strona-glowna/listwa-ustawien.ts
Segmenty listwy dzieli delikatny separator wewnątrz jednej powierzchni, więc pas czyta się jako jeden byt o kilku wejściach, a nie jako trzecia siatka kart; pozycje listwy nie mają formy kart ani kafli. Segmenty ustawień są trzy: okno konfiguracji, tryb Mobile oraz Always On Display. Czwarty segment nie jest ustawieniem — „Dodaj nowy” otwiera formularz zakładania komponentu własnego i zgłasza to osobnym wywołaniem zwrotnym, nie przez wykaz ustawień, ponieważ pozycja, która ustawieniem nie jest, nie udaje jego kodu; zakładanie komponentu należy do strefy drugiej, więc przeniesienie segmentu jest zgłoszone, a nie wykonane z tego katalogu. Waga wizualna strefy jest najniższa z trzech: segment ma wysokość kontrolki i niesie ikonę oraz nazwę, bez wezwania do działania i bez metadanych. Żadna pozycja nie jest wyszarzona ani pozbawiona klikalności.

## budowa/klient-poprzedni/src/moduly/studio/strona-panel-tresci.ts
Pięć czynności panelu — odczyt treści fragmentu, zmiana brzmienia, odczyt
postaci, zapis postaci i wykaz pochodzenia — dotyczą jednego: treści i postaci
dokumentu jako całości, a nie pojedynczej cechy; rozsypane po osobnych panelach
nastaw byłyby nie do znalezienia. Zapis postaci przyjmuje treść jako pole
opcjonalne kontraktu: brak znaczy brak zmiany treści. Przycisk zapisu odczytuje
postać i zapisuje ją z powrotem, zakładając wersję — utrwala to, co nastawy
strony i style zmieniły, bez dotykania treści, żeby postać przestawała ginąć
nawet wtedy, gdy operator nie pisał.

Brzmienie fragmentu ma pole wielowierszowe, ponieważ poprawa fragmentu bez
przepisywania całości jest osią zamówienia, a fragment bywa akapitem, nie
wyrazem; pole jednowierszowe wymuszałoby wklejanie akapitu w linijkę wysokości
jednego wiersza.

## budowa/klient-poprzedni/src/moduly/research/okno-reading-view.ts
Okno jest zbudowane, choć katalog okien rdzenia jeszcze go nie zna, bo dwie jego czynności mają dziś pokrycie: wczytanie treści dokumentu repozytorium i zamiana zaznaczonego fragmentu w ustalenie z powiązaniem do czytanego źródła. Podświetlenia trwałe, notatki na marginesie i wypisy zbiorcze mają już w kontrakcie własne komendy i własny byt wraz z kotwicą pozycji, ale rdzeń nie ma dla nich uchwytu — okno ich nie udaje, zaznaczenie jest zaznaczeniem przeglądarki, a trwałym staje się dopiero jako ustalenie. Pozostałe operacje na źródle stoją w panelu akcji pod nazwami swoich komend i wracają odmową rdzenia, zamiast znikać z okna. Plik składa widok; zachowanie po naciśnięciu leży w module obsługującym czynności lektury.

## budowa/klient-poprzedni/src/moduly/workspace/okno-rozmowy.ts
Dwie kontrolki modułu wołają przeniesienie kontekstu i obie wymagają identyfikatora okna źródłowego:
udostępnienie zaznaczonych zasobów w bibliotece projektu oraz przejście do budowniczego eksperta
w zarządcy ekspertów. Rdzeń niesie kod modułu w opisie okna, więc okna karty sesji wystarczy zawęzić
do modułu. Moduł montuje się od razu, a nie odracza montażu do chwili znalezienia okna, bo okno
rozmowy jest mu potrzebne do dwóch czynności, nie do istnienia — jego brak nie wygasza pięciu okien
operacyjnych. Kolejność względem wejścia do modułu jest treścią: okno rozmowy jest wtedy już
przestawione na ten moduł, a pytanie zadane wcześniej oddałoby okno modułu poprzedniego. Funkcja nie
zgaduje: gdy rdzeń odmówi albo nie odda okna tego modułu, wraca pusty identyfikator wraz z powodem,
który staje potem w odmowie obu kontrolek. Rdzeń oddaje w wykazie oba stany okna, otwarty i zamknięty.

## budowa/klient-poprzedni/src/okno-komunikacji/dyktowanie/dostarczenie-nagrania.ts
Kontrakt oczekuje w komendzie transkrypcji ścieżki do pliku na maszynie
silnika, a nagranie z mikrofonu istnieje w przeglądarce wyłącznie jako bajty
w pamięci karty; kontrakt nie ma dziś komendy przyjmującej dźwięk, a powłoka
nie ma wtyczki systemu plików, więc pliku nie ma czym zapisać. Ogniwo odmawia
zamiast obchodzić ten brak: nagranie nie jest nigdzie zapisywane, ścieżka nie
jest zmyślana, a dźwięk nie opuszcza maszyny. Drogę otworzy dopiero komenda
kontraktu przyjmująca nagranie zakodowane w base64; do tego czasu dyktowanie
pozostaje wyłączone, a tekst wpisuje się z klawiatury. Pomocnik zamieniający
bajty na base64 jest już gotowy, bo kodowanie bajtów jest jedyną częścią tej
drogi możliwą do zbudowania bez zmiany kontraktu; kodowanie idzie porcjami,
ponieważ przekazanie wszystkich bajtów naraz przekracza limit argumentów
wywołania na nagraniu dłuższym niż kilka sekund.

## budowa/klient-poprzedni/src/rozmowa/montaz-rozmowy.ts
Pominięcie polityki pamięci rozmowy dla kodu modułu oznacza, że każdy moduł pracuje z pamięcią.

Montaż nie buduje paska zlecenia i nie zagląda do jego środka: przepuszcza element do widoku nietknięty. Buduje go warstwa mająca komplet sterowania okna; podgląd rozmowy go nie ma i pomija.

Przestawienie widoku transkryptu nie jest nigdy odrzucane i nigdy nie kosztuje rundy do rdzenia: powłoka podaje kod trybu i dostaje przerysowany zapis, nie musząc znać ani wpisów, ani warstw.

Funkcja osadzająca rozmowę okna w dokumencie odpowiada za jedno: powiązanie warstwy rozmowy z jej widokiem i wstawienie całości w kontener. To jest punkt styku warstwy rozmowy z powłoką, która nie musi znać ani widoku, ani kontraktu. Moduł okna jest śledzony, nie zakładany: okno pyta rdzeń o swój moduł i słucha zmian, więc przestawienie modułu wykonane gdziekolwiek indziej przestawia to okno samo. Powłoka może też przestawić okno wprost przez metodę ustawModul, gdy zna wynik komendy wcześniej. Ognisko ląduje na polu wypowiedzi od razu po złożeniu.

Odwrotna kolejność złożenia śledzenia modułu i rozmowy dałaby okno modułu bez pamięci sesyjnej, które odtwarza wątek z rdzenia, zanim się dowie, że nie miało prawa go odtworzyć.

Wykaz po ukośniku jest jedyną drogą doraźnego dostępu do narzędzia, więc nie może zależeć od tego, którym oknem operator akurat pracuje.

Wybór narzędzia z tego okna melduje się sam z odpowiedzi komendy; poszerzenie zestawu spoza okna przychodzi wyłącznie zdarzeniem dołożenia narzędzia sesji. Powód niepodpięcia nasłuchu nie idzie do wątku przy montażu: gdy zdarzenia nie ma w wygenerowanym kontrakcie, ten sam brak melduje się zdaniem w chwili sięgnięcia po wykaz, więc powtarzanie go przy otwarciu każdego okna byłoby hałasem.

## budowa/klient-poprzedni/src/strona-glowna/macierz-modulow.ts
Dane macierzy niesie odczyt środowisk z rdzenia z dołączonymi modułami: każde środowisko ma pole kodów modułów widocznych w jego bocznej nawigacji, w kolejności wyświetlania, a osobnego zapytania o macierz nie ma. Macierz niczego nie blokuje — zanim rdzeń odpowie, i gdy moduł nie stoi w żadnym środowisku, odczyt środowiska modułu zwraca brak wartości, co znaczy „nie wiem, dokąd", a nie „nie wolno"; co z tym zrobić, rozstrzyga miejsce wyboru. Środowisko początkowe to pierwsza karta według kolejności z rdzenia, a przed pierwszą odpowiedzią rdzenia — pierwszy kod z kontraktu, bez nazwy wpisanej ręcznie, ponieważ to ten sam wykaz źródłowy, z którego żyje strefa pierwsza. Kafel prowadzący gdziekolwiek jest lepszy niż kafel, który nagle przestał wiedzieć, dlatego odmowa rdzenia i cisza nie kasują macierzy już poznanej. Remis w przypisaniu modułu do środowiska bierze się stąd, że moduł bywa widoczny w kilku środowiskach naraz, a kontrakt nie niesie środowiska macierzystego — reguła liczy się z danych, bez kodu wpisanego na sztywno.

## budowa/klient-poprzedni/src/moduly/research/okno-report-builder.ts
Widok ma dwie warstwy, bo niesie dwie rzeczy. Sekcja osadzona w przestrzeni modułu pokazuje raport złożony i pozostaje widoczna zawsze. Modal niesie samą kompozycję — zamknięcie kreatora nie może zdejmować okna z ekranu. Plik składa widok; zachowanie po naciśnięciu leży w module obsługującym czynności raportu. Konspekt stoi obok treści, nie zamiast niej: widok podzielony pokazuje strukturę po lewej i treść sekcji po prawej, a kliknięcie pozycji konspektu prowadzi ognisko do sekcji w podglądzie.
Pozycja konspektu przewiduje przy sekcji znacznik ukończenia. Sekcja raportu niesie identyfikator, tytuł, treść, ustalenia i kolejność — pola stanu nie ma, więc konspekt zamiast wymyślonego znacznika mówi rzecz sprawdzalną: czy sekcja ma już treść.

## budowa/klient-poprzedni/src/widok-sterowania/podsumowanie-ustawien.ts

Podsumowanie ośmiu ustawień okna jest odczytem, nie kontrolką. Widoczne wtedy,
gdy szuflada z kontrolkami jest zwinięta — sterowanie zamyka się, a mimo to
widać, na czym okno pracuje. Wartość każdego ustawienia stoi obok jego nazwy,
więc osiem odpowiedzi widać jednym spojrzeniem, bez rozwijania czegokolwiek.
Podsumowanie nie jest bramą: kliknięcie wiersza otwiera szufladę i prowadzi do
kontrolki, a nie odmawia dostępu. Ikona przy każdym wierszu sprawia, że
rozróżnienie nie opiera się na samej barwie.

Powłoka dymka niesie własny odstęp z prawej strony, bez którego dymek kurczy
się w rzędzie nazwy; znak zastępuje styl biblioteczny, bo jest pierścieniem
ze wskaźnikiem pomocy, a nie kwadratowym przyciskiem ikonowym. Reguły stoją
w arkuszu stylu podsumowania ustawień.

Para znaczników nazwy i wartości mieszka we własnym bloku — dopuszcza to
budowa listy opisowej, a układ zyskuje jedną komórkę siatki na ustawienie
zamiast dwóch niezależnych, które przy zmianie liczby kolumn rozjechałyby się
względem siebie. Nazwa jest przyciskiem prowadzącym do kontrolki — nie
oznakowaniem, którego nie da się nacisnąć. Wartość zostaje tekstem, bo jest
odczytem stanu.

## budowa/klient-poprzedni/src/moduly/studio/strona-pola-postaci.ts
Pole liczbowe jest własne, ponieważ moduł kontrolek formularza niesie pole
tekstowe, listę wyboru i przełącznik, a nastawy postaci są w większości
liczbami z granicami: margines w milimetrach, stopień pisma w punktach, punkt
startu numeracji, krycie znaku wodnego. Pole tekstowe przyjęłoby zapis słowny
i wysłało go wprost do rdzenia.

Pole puste znaczy brak zmiany cechy w całym tym module: rodzina komend strony
i formatu ma pola opcjonalne w tym samym znaczeniu. Margines zerowy jest
nastawą, którą operator może wybrać świadomie, więc podanie zera i brak
podania nie mogą znaczyć tego samego — inaczej każde naciśnięcie przycisku
zerowałoby każdą cechę, której operator nie wpisał.

Bilans czynności masowej jest sednem uczciwości tego odcinka: zamiana
w całym dokumencie, która trafiła w blokadę, wykonuje się poza blokadą i musi
nazwać, którą. Przemilczenie pominięcia jest zakazane, ponieważ liczba
zmienionych miejsc równa zero bez słowa wyglądałaby jak wykonana czynność.

## budowa/klient-poprzedni/src/powloka/most-katalogow.ts

Most do natywnego okna wyboru katalogu jest jedynym konsumentem polecenia powłoki
wybierz_katalog_roboczy. Okno systemu operacyjnego zwraca ścieżkę istniejącą i
rozwiniętą, czego samo pole tekstowe nie zapewnia. Katalog wskazuje się w dwóch
sprawach — katalog roboczy modelu oraz zakres wglądu modelu, czyli punkt dostępu
rodzaju katalogu lokalnego — ale czynność systemu jest w obu ta sama, więc okno
jest jedno, a sprawa zmienia wyłącznie napis w belce okna. Poza powłoką natywną
wskazanie wraca wartością pustą tak samo jak rezygnacja z wyboru, a widok
zostawia drogę wpisania ścieżki ręcznie. Żadna ścieżka wykonania nie rzuca
wyjątkiem i nie odrzuca obietnicy.

## budowa/klient-poprzedni/src/okna-rownolegle/podglad-ukladu.ts
Komentarz przy tworzeniu układu odsyłał czytelnika do typu `OpcjeGniazda` zamiast podać
wprost powód pustych gniazd. Redakcja usunęła odesłanie i przeniosła wyjaśnienie
bezpośrednio nad instrukcję, której dotyczy, jako pełne zdanie: brak rdzenia i osadzonej
rozmowy w tym stanowisku podglądu skutkuje pustymi gniazdami układu przy każdym
uruchomieniu. Komentarz nad złożeniem obszaru skrócono do jednego zdania nazywającego
elementy kompozycji bez wyliczenia wtórnego wobec kodu importów.

## budowa/klient-poprzedni/src/moduly/workspace/poziomy-zasiegu.ts
Poziom aplikacji stoi w wykazie, ale poza wyborem: kontrakt przypisuje mu nastawy samego programu,
nie treści w nim prowadzonej, więc kontrolka oferująca ten poziom proponowałaby zapis, którego rdzeń
nie ma gdzie umieścić. Nazwę zachowuje, bo rdzeń może oddać ten poziom w odpowiedzi, i okno ma go
wtedy nazwać, a nie pokazać surowego kodu. Wykaz jest zapisany jako mapa poziomu na opis, więc poziom
dołożony do kontraktu nie przejdzie kompilacji, dopóki nie zostanie rozstrzygnięte, czy moduł ma go
pokazywać. Pamięć projektu otwiera się na poziomie projektu, bo taki jest jej domyślny zasięg —
zmienia się wyłącznie kolejność prezentacji, nazwy zostają te same, więc oba okna pokazują tę samą
nastawę pod tą samą nazwą. Miejsce puste w kolejności znaczy poza wyborem i wtedy powód pominięcia
mówi, dlaczego.

## budowa/klient-poprzedni/src/strona-glowna/meldunki-sesji.ts
Wspólne miejsce dwóch zdań meldunku trzyma jedną odpowiedź na oba pytania. Odmowa mówi treścią odmowy z rdzenia; zdanie zastępcze wchodzi tylko wtedy, gdy odmowa przyszła bez treści, i nazywa wówczas czynność, której dotyczyła.

## budowa/klient-poprzedni/src/widok-sterowania/szuflada.ts

Zwinięta szuflada pokazuje podsumowanie ośmiu wartości, rozwinięta — komplet
kontrolek katalogu sterowania. Zwijanie nie jest blokadą: żaden element nie
traci klikalności, zmienia się wyłącznie to, która warstwa zajmuje miejsce.
Uchwyt pozostaje czynny w każdym stanie. Stan szuflady jest ogłaszany
magistralą, więc uchwyt paska górnego i uchwyt panelu pokazują tę samą
prawdę, zamiast każdy swoją.

## budowa/klient-poprzedni/src/moduly/research/okno-research-workspace.ts
Zakres badania zapisuje komenda ustawienia przestrzeni roboczej, niosąca zakres i etapy; zapis kończy się odpowiedzią rdzenia albo jego odmową, nigdy ciszą. Nawigacja do pozostałych okien modułu jest przeniesieniem ogniska wewnątrz przestrzeni modułu: okna stoją obok siebie, nie w osobnych trasach, więc nie woła rdzenia.
Wskazanie etapu prowadzi wzrok do okna właściwego temu krokowi pracy. Przypisanie źródła do etapu ma już pole w żądaniu katalogowania, ale samo źródło oddawane przez rdzeń go nie niesie, więc okno nie ma po czym zawężać i mówi to wprost, zamiast udawać filtr.

## budowa/klient-poprzedni/src/powloka/most-rdzenia.ts

Most do wiedzy powłoki natywnej o rdzeniu jest konsumentem poleceń adres_rdzenia
i stan_rdzenia. Powłoka stawia proces rdzenia, zna jego port ze zmiennej środowiska,
jego identyfikator procesu i ścieżkę dziennika. Interfejs nie zna żadnej z tych
rzeczy: adres gniazda wylicza z lokalizacji dokumentu, a powodu ciszy nie zna wcale.
Stąd dwa polecenia: adres_rdzenia zwraca adres HTTP rdzenia lokalnego, potrzebny
gdy okno dostało interfejs z pakietu osadzonego w powłoce, bo strona ma wtedy
pochodzenie lokalne powłoki, a wyliczenie z lokalizacji dałoby adres gniazda, pod
którym nie nasłuchuje nikt; stan_rdzenia zwraca opis rdzenia w tle wraz z
przebiegiem uruchomienia, bez którego wskaźnik łączności umie powiedzieć wyłącznie
że jest rozłączony, bez powodu i bez wskazania dziennika. Poza powłoką natywną
oraz przy niepowodzeniu polecenia odpowiedzią jest wartość pusta. Żadna ścieżka
nie rzuca wyjątkiem i nie odrzuca obietnicy — brak odpowiedzi powłoki niczego nie
wstrzymuje, bo interfejs ma własną drogę ustalenia adresu i własny stan łączności.
Polecenia zmieniającego stan rdzenia tu nie ma: zatrzymanie rdzenia jest czynnością
z zasobnika powłoki, a most nie tworzy drugiej drogi sterowania platformą.

## budowa/klient-poprzedni/src/widok-sterowania/uchwyt-paska.ts

Drugie wejście do tej samej szuflady stoi z paska, obok przełącznika motywu.
Na pasku stoi wyłącznie przycisk ikonowy o wymiarze kontrolki paska; komplet
kontrolek mieszka w kolumnie obok sceny, bo w prawym rogu paska nie miałby
się gdzie zmieścić. Naciśnięty przy zwiniętej szufladzie rozwija ją, przy
rozwiniętej — zwija.

## budowa/klient-poprzedni/src/moduly/workspace/agenci-pozycja.ts
Pozycja wykazu ekspertów stoi w osobnym pliku od okna, bo to inna odpowiedzialność: okno prowadzi
odczyt i przypisanie, pozycja rysuje jednego eksperta. Pozycja nie trzyma własnego stanu — rysuje
komplet uprawnień z pola odpowiedzi, tak jak podał go rdzeń. Uprawnienie ma trzy stany, nie dwa:
grupa, o której rdzeń nie powiedział nic, nie jest ani przyznana, ani odebrana, i tak też jest
wypisana. Grupa bez odpowiedzi rdzenia daje przy naciśnięciu żądanie przyznania; wykaz odrysowuje
okno dopiero z odpowiedzi rdzenia na zmianę uprawnienia.

## budowa/klient-poprzedni/src/okna-rownolegle/role-domyslne.ts
Pierwotny nagłówek wykładał regułę doboru ról względem obsady pętli multitaskingu
wierszami przykładów dla jednego, dwóch, trzech i czterech okien oraz zastrzeżenie
o nadpisywaniu roli ustawionej ręcznie. Redakcja zostawiła w nagłówku samo działanie
funkcji, a wyliczenie obsady i relację z ręcznym nadaniem roli pozostawiono kodowi
i nazwom stałych, które już je wyrażają.

## budowa/klient-poprzedni/src/okna-rownolegle/skroty-paneli.ts
Nagłówek opisywał osobno pochodzenie zestawu skrótów, powód pomijania pozycji
nieotwieralnych, budowę dymka z nazwą i skrótem, powód niekorzystania z komponentu
`.dn-tooltip`, brak nasłuchu klawiszy oraz źródło znacznika nowości. Redakcja
scaliła te wątki do jednego zdania nazywającego złożenie i miejsce wyświetlania
rzędu skrótów, a szczegóły konstrukcyjne pozostały czytelne w samym kodzie modułu.

## budowa/klient-poprzedni/src/okna-rownolegle/pustka-relacji.ts
Nagłówek uzasadniał brak znikania pasa przy pustej relacji względem mylącego
komunikatu dla Operatora oraz zastrzegał jeden wariant stanu pustego w bibliotece.
Redakcja pozostawiła samo działanie modułu, a uzasadnienie wyboru komunikatu
przeniesiono poza komentarz w kodzie.

## budowa/klient-poprzedni/src/okna-rownolegle/opis-gniazda.ts
Nagłówek tłumaczył pośrednią pozycję okna między sesją a wiadomością oraz zasadę
kopiowania zamiast współdzielenia ustawień między oknami. Redakcja zostawiła
zwięzły opis działania funkcji, a rozważania o zasięgu ustawień przeniesiono
poza komentarz w kodzie.

## budowa/klient-poprzedni/src/okna-rownolegle/objasnienie-ukladu.ts
Nagłówek opisywał zasadę obecności objaśnienia przy każdym elemencie konfiguracji,
sposób działania komponentu `.dn-tooltip` przy najechaniu i ognisku oraz podział
odpowiedzialności między biblioterkę chmurki a widok rysujący znak zapytania.
Redakcja scaliła te wątki w jedno zdanie nazywające działanie modułu, a podział
odpowiedzialności między warstwami pozostał czytelny w strukturze kodu.

## budowa/klient-poprzedni/src/moduly/research/okno-sources-manager.ts
Pole zawężania działa po stronie klienta i nie woła rdzenia. Komenda odczytu wykazu wraz z polami zawężającymi jest już w kontrakcie, ale rdzeń nie ma dla niej uchwytu, więc wykaz mieszka w pamięci modułu i tam też się zawęża; po dobudowie zawężanie ma przenieść się do żądania. Zawężenie nie zdejmuje zaznaczenia — źródło niewidoczne w wykazie pozostaje zaznaczone i idzie do ustalenia, a okno mówi o tym liczbą przy polu. Plik składa widok; zachowanie po naciśnięciu leży w module obsługującym czynności źródeł.
Stan okna liczy się z wykazu pełnego, nie z zawężonego: zawężenie bez trafień nie znaczy, że katalog jest pusty, a zaproszenie do skatalogowania pierwszego źródła postawione nad katalogiem pełnym byłoby nieprawdą.

## budowa/klient-poprzedni/vitest.config.ts

Plik jest osobny, a nie sekcją test dopisaną do konfiguracji budowania: budowanie
pakietu i sprawdzanie kodu to dwa różne byty. Vitest czyta ten plik z pierwszeństwem
przed konfiguracją budowania, więc konfiguracja budowania zostaje nietknięta.
Środowisko domyślne node, ponieważ warstwy protokołu i łączności nie dotykają
dokumentu. Widoki dotykają go wprost — tworzą elementy, czytają arkusze,
przełączają motyw — więc ich sprawdziany dostają jsdom przez dopasowanie ścieżek.
Zgoda na katalog nadrzędny obejmuje import wspólnego kontraktu, jedynego źródła
prawdy nazw, tak samo jak w konfiguracji budowania.

## budowa/klient-poprzedni/src/powloka/nawigacja-modulow.ts

Trzeci pas powłoki jest pionową, stałą nawigacją modułów środowiska o jednej
odpowiedzialności: wykaz pozycji bieżącego środowiska i wskazanie pozycji
wybranej. Nagłówek kolumny niesie nazwę środowiska krojem szeryfowym, pod nią
motto i liczbę pozycji; pozycja wybrana dostaje złoty pasek przy lewej krawędzi.
Nawigacja nie zna ani jednej nazwy modułu — pyta o wykaz podłączone źródło
nawigacji korzystające z komend wejścia do środowiska i wykazu modułów. Do
czasu odpowiedzi kolumna pokazuje stan wczytywania, a po odmowie treść odmowy;
kopii katalogu modułów nie ma tu żadnej. Źródło dokłada się po złożeniu:
powłoka powstaje bez połączenia z rdzeniem i dopiero warstwa składająca
aplikację ma czym ją zasilić — stąd osobne podłączenie źródła, a nie parametr
wytwórni. Żadna pozycja nie traci klikalności: wykaz wynika ze środowiska,
a nie z gotowości modułu. Wykaz modułu spoza środowiska bierze się z macierzy
widoczności, a ta rozstrzyga tylko o obecności modułu na liście, nie o prawie
do jego otwarcia — kafel składa wtedy pozycję z katalogu modułów i wskazuje ją
osobną drogą, a kolumna nie zapala żadnego wiersza, bo żaden jej wiersz nie
odpowiada temu modułowi.

## budowa/klient-poprzedni/src/moduly/research/pamiec-badania.ts
Pamięć ma jedną odpowiedzialność: przechowanie i ogłoszenie zmiany. Pamięć nie zna kontraktu i nie woła rdzenia — dzięki temu odczyt badania i wykazy okien patrzą na ten sam zbiór, a nie na jego kopie. Fazy są trzy, nie dwie: stan przed pytaniem, stan pytania w toku i stan, w którym rdzeń nic nie ma, zostają rozróżnialne, bo zlanie ich w jedno kazałoby zgadywać, czy czekać, czy działać.

## budowa/klient-poprzedni/src/moduly/research/panel-akcji.ts
Panel nie zna kontraktu i nie woła rdzenia bezpośrednio: droga wykonania akcji, czy to komenda
`window.action`, komenda dziedzinowa, czy uczciwa odmowa, należy do okna, nie do paska przycisków.
Żaden przycisk nie jest wygaszany ani warunkowany zaznaczeniem — brak warunku merytorycznego
nazywa odpowiedź po kliknięciu, a nie odbiera klikalności. Każdy przycisk ma dymek objaśniający,
czym akcja jest i czym się kończy.

## budowa/klient-poprzedni/src/strona-glowna/menu-sesji.ts
O tym, które czynności są sensowne dla wiersza, rozstrzyga moduł czynności sesji; menu ich nie zna.
Pusty wykaz czynności daje wynik pusty zamiast pustego pojemnika — wiersz czysto informacyjny nie
dostaje ramki po menu, którego nie ma. Na czas wykonania wiersz nie przyjmuje drugiej czynności,
także z innego przycisku niż naciśnięty: archiwizacja w trakcie zmiany nazwy dawałaby dwa żądania
o tę samą sesję o przypadkowej kolejności skutków. Zajętość nie jest bramą i nie gasi przycisków:
atrybut wyłączenia nie jest stosowany nigdzie w produkcie — przycisk zostaje klikalny, a stan
niesie napis: przycisk czynny dopisuje wielokropek, więc zajętość jest widoczna bez samego koloru.
Naciśnięcie w trakcie biegu wraca meldunkiem, nie ciszą. Odmowa wraca meldunkiem, nie wyjątkiem:
wykonanie oddaje zdanie albo wynik pusty, a menu podaje je dalej. Wyjątek nieprzewidziany też
kończy się meldunkiem, żeby wiersz nie został z zablokowanymi przyciskami.

## budowa/klient/src/gniazdo-zastepcze.ts
Gniazdo zastępcze nie udaje rdzenia, lecz sprawdza zachowanie klienta wobec gniazda: kolejkowanie, ponawianie, kolejność stanów oraz przeżycie subskrypcji przy wymianie gniazda. Sama biblioteka gniazda nie podlega tu sprawdzeniu. Rozmowę z rdzeniem rzeczywistym mierzy osobny sprawdzian, który tego rdzenia wymaga.

## budowa/klient-poprzedni/src/moduly/research/panel-postepu.ts
Każdy licznik powstaje z pamięci modułu, a pamięć niesie wyłącznie odpowiedzi rdzenia z tego
połączenia — żadna z liczb nie jest oszacowaniem, „źródeł 14" znaczy czternaście źródeł, które
rdzeń skatalogował i oddał. Panel celowo nie liczy pokrycia pytań badawczych źródłami ani podziału
źródeł na przeczytane i nieprzeczytane, z dwóch różnych powodów: pokrycie pytań czeka na uchwyt
w rdzeniu, kształt odpowiedzi jest już rozstrzygnięty i licznik dojdzie razem z uchwytem; stan
lektury ma pole w żądaniu zapisu źródła, ale odpowiedź, którą rdzeń oddaje, tego pola nie niesie —
wartość da się wysłać, a nie da się jej odczytać z powrotem, i tego uchwyt sam nie naprawi, dopóki
byt źródła nie dostanie pola. Procent policzony z danych, których nie ma, byłby metryką zmyśloną,
więc panel go nie pokazuje i mówi, czego brakuje. Pokrycie ustaleń jest natomiast policzalne
i policzone: sekcje raportu niosą pole `findingIds`, więc liczba zebranych ustaleń, które weszły
do dokumentu, wychodzi z porównania dwóch zbiorów, które rdzeń oddał.

## budowa/klient/src/main.ts
Plik pozostaje wyłącznie kompozycją: gniazdo, kanał, sesja, tożsamość, przebieg i montaż. Etap i odsłonę rozstrzyga przebieg, węzły stawia montaż — plik wejścia nie rozgałęzia drogi ani nie dotyka dokumentu poza wskazaniem korzenia montażowi.
Magazyn tokenu bramki zostaje domyślny, czyli w pamięci procesu. Magazyn trwały należy do powłoki i żadne źródło go dziś nie wskazuje; magazyn udający trwałość obiecywałby rozpoznanie urządzenia, którego nie ma.
Rdzeń wystawia pakiet interfejsu obok gniazda, więc dokument wczytany po HTTP przyszedł z rdzenia i to jego adres jest adresem gniazda. Dokument wczytany inaczej — z pliku albo z protokołu powłoki — pochodzenia nie niesie.

## budowa/klient-poprzedni/src/moduly/translate/kontrolki-translate.ts
Dymek jest budowany lokalnie, a nie brany z biblioteki komponentów, ponieważ nosi klasy arkusza
stylów modułu tłumaczeń, podczas gdy wersja biblioteczna wymaga klas układu, których ten arkusz
nie udostępnia, co dałoby dymek bez pozycjonowania. Dymek pokazuje się na najechaniu albo na
skupieniu klawiaturowym, bez klikania i bez osobnego zamykania; znak jest przyciskiem, więc jego
naciśnięcie przenosi ognisko i zwraca odpowiedź głosową, a treść dymku leży w opisie dostępności,
więc dymek nie potrzebuje identyfikatora i nie koliduje między oknami.

## budowa/klient-poprzedni/src/okno-komunikacji/dyktowanie/dyktowanie.ts
Plik nie rysuje interfejsu — nie tworzy ikony, menu ani paska i nie zna klas CSS; ikona mikrofonu, wykaz urządzeń i przełącznik przytrzymania należą do obszaru paska poleceń, inaczej powstałyby dwa mikrofony w interfejsie. Dostępność jest pytaniem zadawanym przed narysowaniem ikony, nie odpowiedzią po naciśnięciu: mikrofon, którego nie ma czym obsłużyć, nie pojawia się wcale, a wyszarzony albo odmawiający po kliknięciu mikrofon byłby bramą zamiast krótszym paskiem.

## budowa/klient/src/polaczenie/adres-rdzenia.ts
Warstwa połączenia zna wyłącznie miejsce nasłuchu rdzenia. Wyboru adresu nie dokonuje sama — wskazuje go wołający, bo to on wie, czy rdzeń stoi na tej samej maszynie, czy pod adresem podanym przez powłokę.
Nasłuch bez wskazania adresu wiąże się z pętlą zwrotną, więc to jest adres, pod którym rdzeń stoi, dopóki nikt nie wskazał inaczej.
Przekład adresu HTTP rozstrzyga schemat i ścieżkę gniazda; gdy adres nie jest adresem HTTP, wołający rozstrzyga, czy sięgnąć po pętlę zwrotną, czy odmówić. Ścieżka gniazda jest własnością tej warstwy, więc przekład mieszka tutaj, a nie u tego, kto adres HTTP zdobył.

## budowa/klient-poprzedni/src/moduly/research/pasek-etapow.ts
Etapy przychodzą polem `stages` odpowiedzi `research.workspace.set` i wcześniej nigdzie się nie
pokazywały: okno wiodące przyjmowało je w polu tekstowym i odsyłało do rdzenia, a Operator nie
widział, co rdzeń naprawdę zapisał. Pasek celowo nie pokazuje stanu etapu w postaci „ukończony,
bieżący, nierozpoczęty" wraz z checklistą, bo kontrakt niesie etapy jako zwykły wykaz napisów bez
pola na stan — znaczek ukończenia postawiony tutaj byłby danymi zmyślonymi. Wskazanie etapu jest
nastawą widoku i tylko nią: rdzeń nie ma gdzie zapisać, na którym etapie stoi badanie, więc po
ponownym wejściu do modułu wskazanie zaczyna od zera, a okno tego nie ukrywa.

## budowa/klient/src/polaczenie/dziennik-nieznanych.ts
Rdzeń odpowiada na nieznaną komendę zdarzeniem `*.unknown` właściwym dla obszaru nazwy; każdy obszar kontraktu ma własne zdarzenie zapasowe. Dziennik obejmuje je wszystkie, sięgając po komplet z mapy kontraktu zdarzeń nieznanych, nigdy po literał nazwy: dopisanie obszaru w kontrakcie rozszerza dziennik samo, bez zmiany tego pliku.
Osobno przechwytywane są koperty o typie spoza kontraktu — takie, których nie zna ani wykaz komend, ani wykaz zdarzeń. Powstają, gdy rdzeń wyprzedził klienta wersją albo gdy ramka była nieczytelna.

## budowa/klient-poprzedni/src/moduly/research/pustka-okien.ts
Stan pusty tłumaczy, czym okno jest i jak je zapełnić, a nie melduje samą nieobecność pozycji —
takie zdania są tekstem produktu, nie logiką widoku, i mają dać się przeczytać obok siebie oraz
poprawić w jednym miejscu, gdy kontrakt dostanie komendę, której dziś nie ma. Ten sam wzorzec niesie
plik `moduly/library/etykiety-biblioteki.ts`. Brak okna badania to co innego niż pusty wykaz, więc
rozróżnienie stoi w jednym miejscu, nie w siedmiu oknach: komendy zapisu źródła, ustalenia i budowy
raportu mają pole `windowId` obowiązkowe, więc bez okna zaproszenie „skataloguj pierwsze źródło"
prowadzi Operatora wprost pod odmowę. Żadne z tych zdań nie odbiera klikalności ani jednej
kontrolce — nazywa warunek merytoryczny, po czym formularze i przyciski zostają czynne, a odmowa
rdzenia wraca własnymi słowami rdzenia.

## budowa/klient-poprzedni/src/okno-komunikacji/dyktowanie/nagrywanie.ts
Kolejność preferencji rodzaju treści: silnik transkrypcji przyjmuje zapisy `.wav .ogg .m4a .webm`, a `audio/webm` stoi pierwszy, bo w oknie osadzonym opartym o technologię Chromium bywa jedynym wspieranym zapisem; kodek dobiera się tam, gdzie przeglądarka go wymaga do rozstrzygnięcia, a wpis bez kodeka stoi zaraz za nim jako zapasowy. Zwrócony rodzaj treści nie wraca na podstawie samej preferencji — pusty napis, gdy żadna pozycja wykazu nie przechodzi, oznacza wybór własny nagrywarki, odczytany dopiero po fakcie, aby napis rodzaju nie stał się atrapą podpisu pod bajtami, których nie użył. Wspólne zdanie odmowy dla braku zgody i braku sprzętu kazałoby szukać przyczyny po omacku — pierwsza jest decyzją cofalną w ustawieniach przeglądarki, druga stanem sprzętu naprawianym kablem.

## budowa/klient/src/polaczenie/gniazdo.ts
Gniazdo, które błędem kończy samo nawiązywanie, zamyka się bez niczyjej pomocy i ogłasza to zdarzeniem close; wywołanie close na takim gnieździe wywołuje kolejny błąd i wpada w nawrót bez końca — stąd zamknięcie na błędzie dotyczy wyłącznie gniazda już otwartego.
Ponawianie jest bezterminowe, więc bez zaniechania na żądanie wołającego proces klienta nie miałby jak dojść do końca: zaplanowana próba trzymałaby go przy życiu, a każde zamknięcie gniazda planowałoby następną. Zaniechanie dotyczy połączenia, nie treści, która czeka na wysłanie.

## budowa/klient-poprzedni/src/moduly/roundtable/zestawienie-udzialu.ts
Moduł liczy wyłącznie udział mierzalny z bieżących danych: liczbę i objętość wypowiedzi
widocznych przez okno w turze. Głosowanie, oceny rubrykowe i ranking między sesjami mają
osobne komendy w kontrakcie wymiany, których okno jeszcze nie wywołuje, dlatego zestawienie
nie orzeka o quorum ani o wyniku głosowania — dopóki te komendy nie są wywoływane, liczenie
ich skutków byłoby czynnością pozorowaną. Wypowiedź otwartą liczy się z treści widocznej
w oknie, a nie z pola danych źródłowych, ponieważ rdzeń rozgłasza ją najpierw pustą i dopisuje
kolejne fragmenty strumieniem — sam zapis źródłowy dawałby zero znaków przez cały czas trwania
wypowiedzi. Wypowiedzi wcześniejsze, już utrwalone, liczą się wprost z zapisu.

## budowa/klient-poprzedni/src/powloka/nieczynne-w-pasku.ts
Zasada zera blokad zabrania wygaszania kontrolki jako sposobu powiedzenia „to jeszcze nie działa".
Kontrolka zostaje klikalna, a z klikalności wynika obowiązek odwrotny: każde naciśnięcie musi dać
odpowiedź. Odpowiedź mówi, co jest nieczynne i dlaczego, i niczego nie udaje. Każdy powód niżej jest
sprawdzalny w `shared/contract.json`: jeżeli komendy nie ma, tekst mówi wprost, że jej nie ma, i nie
zmyśla nazwy. Zdanie zaprzeczające funkcji, która istnieje, jest tak samo szkodliwe jak przycisk
udający funkcję, której nie ma — prowadzi Operatora od okna, które by mu pomogło.

## budowa/klient-poprzedni/src/moduly/translate/macierz-izolacji.ts
Element trzyma się podziału odczytu i zapisu: platforma ma własne okno punktów izolacji z pełną
obsługą zapisu, a druga kontrolka nad tym samym kluczem dawałaby dwa miejsca zmiany jednej
wartości i dwa różne obrazy stanu przy odmowie jednego z nich. Zakres odczytu to karta sesji i jej
warstwa — moduł pyta o politykę obowiązującą tę pracę, nie o domyślną politykę platformy, a sesji
nieustalonej nie podmienia na globalną, tylko zgłasza brak przedmiotu zapytania. Stan wyjściowy
platformy to zero blokad: żaden zakres nie jest domyślnie odcięty, a macierz nazywa ten skutek przy
każdym kluczu z osobna, ponieważ wykaz zakresów się przewija.

## budowa/klient-poprzedni/src/moduly/research/pustka-okien.ts
Wybór zdania pustki nie obejmuje faz `odczyt` i `blad`: rozstrzygają je czynności okien wcześniej,
odpowiednio wskaźnikiem odczytu i komunikatem odmowy, więc funkcja wybiera między trzema
pozostałymi stanami — brak jeszcze zapytania, zapytanie bez okna w rdzeniu, i zapytanie z miejscem,
w którym treści jeszcze nie ma.

## budowa/klient-poprzedni/src/moduly/roundtable/zrodlo-roundtable.ts
Obszar roundtable niesie czterdzieści sześć komend, z których interfejs klienta wywołuje dziś
tylko cztery — obsługi pozostałych jeszcze nie zbudowano, a wykaz brakujących komend prowadzi
katalog-funkcji.ts. Pole windowId jest wymagane w każdej czynności, bo debata jest bytem okna,
nie sesji: uczestnicy, tury i stanowisko należą do jednego okna debaty, a rdzeń bez tego
identyfikatora nie wie, o którą debatę pytamy — wymóg stoi w kontrakcie wymiany, więc żadne
pole żądania nie ma tu wartości domyślnej. Debata z wieloma modelami odmawia z powodów
zwyczajnych, takich jak nieczynny kanał uczestnika albo już zamknięta tura, dlatego każda
czynność oddaje Wynik zamiast samej treści, żeby okno odróżniło brak wypowiedzi od nieudanego
zapytania.

## budowa/klient-poprzedni/src/okno-komunikacji/dyktowanie/przytrzymanie.ts
Escape jest nasłuchiwany na dokumencie, a nie na elemencie przycisku, ponieważ podczas trzymania myszą ognisko klawiatury bywa gdzie indziej i przycisk wcale nie musi je mieć — nasłuch przy samym elemencie przepuściłby wtedy skrót i nie byłoby czym cofnąć nagrania w połowie wypowiedzi. Zdarzenie `pointercancel`, będące systemowym wyrwaniem gestu przez przewinięcie albo telefon, kończy nagranie zamiast je porzucać, bo użytkownik zdążył już coś powiedzieć, a odrzucenie słów następuje wyłącznie na wyraźne żądanie klawiszem Escape.

## budowa/klient-poprzedni/src/powloka/obszar-roboczy.ts
Jedna odpowiedzialność: kontener na treść wybranej pozycji nawigacji wraz z uczciwym stanem pustym
dla pozycji, która zbudowanego widoku nie ma. Obszar nie zna żadnego widoku i żadnego rdzenia —
przyjmuje gotowy element od warstwy, która go składa. Plansza obszaru niesie dwie warstwy naraz:
widok modułu w wierszu górnym i scenę okien komunikacji pod nim. Okno rozmowy jest oknem wiodącym
każdego modułu, więc obszar pokazujący jedną warstwę naraz odbierałby rozmowę każdemu modułowi,
który doczekał się własnego widoku. Widok raz osadzony nie jest niszczony: `pokaz` dokłada go przy
pierwszym użyciu i dalej wyłącznie przestawia widoczność atrybutem `hidden`. W obszarze stoi scena
z żywym oknem rozmowy — gniazdem WebSocket, strumieniem odpowiedzi modelu, historią wpisów
i obserwatorami układu okien równoległych. Odmontowanie sceny przy każdym przełączeniu modułu
zrywałoby rozmowę w połowie zdania; ukrycie zostawia ją nietkniętą. Stan pusty nie zapowiada modułu,
którego nie ma — nazywa go i pisze wprost, że jego okna operacyjne nie zostały zbudowane, wymieniając
katalog kodów wzięty z rdzenia.

## budowa/klient-poprzedni/src/moduly/workspace/zarzadca-agentow.ts
Wykaz ekspertów przypisanych do projektu pochodzi z zestawienia pulpitu, ponieważ
kontrakt nie udostępnia oddzielnej komendy odczytu przypisań modułu Workspace.
Uprawnienia są konfiguracją możliwości, nie kontrolą dostępu — cztery grupy
zakresu ustawia komenda modułu Agents, a ich stan przychodzi w polu uprawnień
agenta, którego okno samo nie wylicza. O przypisaniu do projektu rozstrzyga
wyłącznie rdzeń: dopóki pulpit nie został odczytany, okno nie zgaduje
przypisań i nie przedstawia nieodczytanego pulpitu jako projektu bez
ekspertów.

Zawężenie wykazu ekspertów pracuje wyłącznie na bibliotece już odczytanej.
Zawężenie do przypisanych nie udaje pustego wykazu przy niewiedzy rdzenia —
nazywa brakującą przesłankę. Plakietka licznika gaśnie razem z wykazem, bo
licznik sprzed odmowy dotyczyłby biblioteki, której okno nie zdołało
odczytać. Po zapisie przypisania okno przyjmuje odpowiedź rdzenia zamiast
nieaktualnego zestawienia pulpitu, żeby wykaz pod potwierdzeniem nie
zaprzeczał samemu potwierdzeniu. Powód wyświetlany przy skoku do modułu
Agents pochodzi z odpowiedzi rdzenia, a znacznik przeniesienia — z okna
otwartego przez rdzeń, nie ze stałej wpisanej w wywołanie. Skasowanie
zestawienia pulpitu wraca stan przypisań do „nie wiadomo”, a nie do zera,
bo o przypisaniach nowego projektu nikt jeszcze nie pytał.

## budowa/klient/src/polaczenie/stan-polaczenia.ts
Ramka wpisana przy rozłączeniu trafia do kolejki wychodzącej i idzie do rdzenia po wznowieniu połączenia.

## budowa/klient-poprzedni/src/moduly/workspace/zrodlo-sasiadow.ts
Okna Workspace wołają komendy sąsiednich modułów wprost, zamiast trzymać
własne odpowiedniki tych samych czynności. Plik stoi osobno od głównego
źródła modułu Workspace, ponieważ niesie inną odpowiedzialność: tamten plik
opisuje obszar własny modułu, ten — jego sąsiedztwo.

## budowa/klient-poprzedni/src/powloka/pasek-gorny.ts
Akcje mieszkają w pliku akcje-paska.ts, wyszukiwanie w pliku wyszukiwanie-globalne.ts (materiał)
i w pliku wykaz-wynikow.ts (obsada menu). Pole poleceń szuka w otwartym środowisku, jego modułach,
kodach okien operacyjnych, kartach sesji i ustawieniach platformy; przede wszystkim nie przeszukuje
treści, bo kontrakt nie ma komendy wyszukiwania globalnego. Wykonania polecenia pole nie udaje: Enter
bez trafienia mówi wprost, że centrum poleceń nie jest zbudowane. Podpowiedź „Ctrl K" nie jest ozdobą
— skrót naprawdę prowadzi ognisko do pola. Powłoka ma drugi pasek o innym składzie, w pliku
aplikacja/pasek-aplikacji.ts (Centrum dowodzenia i Mission Control), który niesie menu aplikacji
i menu Operatora, a nie ma pola wyszukiwania ani dzwonka.

Pominięte źródło materiału wyszukiwania znaczy „ten pasek nie ma czego przeszukać" i pole mówi to
wprost przy zatwierdzeniu, zamiast pokazywać pusty wykaz udający wyniki. Powłoka podaje źródło sama,
więc w produkcie ta gałąź nie zachodzi — zostaje dla podglądu i dla sprawdzianów.

Strzałki w wykazie wyników przesuwają wyróżnienie wirtualne mechanizmu, a nie ognisko; przeniesienie
go na pozycję wyrwałoby Operatorowi klawiaturę spod palców w połowie pisania. Bez źródła materiału
pole zostaje czynne, a Enter mówi, że centrum poleceń nie jest zbudowane — to brak materiału
powiedziany wprost, nie blokada kontrolki.

## budowa/klient-poprzedni/src/okno-komunikacji/dyktowanie/transkrypcja-nagrania.ts
Zatrzymanie na dostarczeniu i odmowa transkrypcji dają ten sam stan nieprzetworzony, bo skutek dla użytkownika jest jeden: tekstu nie ma. Powód jest jednak inny i inna jest naprawa — w pierwszym przypadku brakuje komendy w kontrakcie, w drugim rdzeń miał nagranie i go nie przerobił, dlatego oba zdania odmowy pozostają rozdzielone. Obietnica wyniku nie jest odrzucana: odmowa wraca jako wynik ze stanem nieprzetworzonym, nie jako wyjątek, więc widok nie zakłada przechwytywania błędu, a nieudane dyktowanie nie wywraca okna. Dźwięk nie opuszcza maszyny użytkownika inną drogą niż przez dostarczenie — ogniwo transkrypcji pilnuje tej granicy samodzielnie.

## budowa/klient-poprzedni/src/moduly/workspace/zrodlo-workspace.ts
Obszar workspace liczy w kontrakcie więcej komend, niż plik wystawia: zadania,
tablica, harmonogram, kalendarz, notatki, graf wiedzy, oś czasu i komentarze
mają nazwy, lecz nie mają jeszcze obsługi w rdzeniu ani okien w tym module.
Nieobecność w tym pliku znaczy "niezbudowane", nie "nieistniejące w
kontrakcie" — o pokryciu każdej z nich orzeka źródło braków kontraktu,
pytając rdzeń o jego własny wykaz. Brak wartości domyślnej zamiast wyniku
błędu w całym pliku jest zamierzony: okna modułu mają obowiązkowy stan
błędu i mają odróżnić brak danych od nieudanego zapytania.

Usunięcie wpisu pamięci jest osobną rodziną komend niż zapis i odczyt
instrukcji, choć dotyczy tego samego bytu: kasuje wpis założony zapisem
instrukcji, po czym znika on z odczytu pamięci projektu. Odczyt pamięci w
zasięgu karty sesji nie jest drugim źródłem prawdy dla pamięci projektu —
pyta o to, co widzi karta sesji przy włączonych poziomach, a poziomy
przestawia osobna komenda przełączania pamięci.

## budowa/klient-poprzedni/src/moduly/translate/magazyn-translate.ts
Zbiór jest oddzielony od odczytu: stan-translate.ts wie, jak zapytać rdzeń, ten plik wie wyłącznie,
co moduł już wie, dzięki czemu okna można sprawdzić na samym zbiorze, a odczyt nie miesza się
z pamięcią. Trzy fazy kontekstu są rozróżnialne: faza spoczynku znaczy brak zapytania, faza odczytu
znaczy zapytanie w toku, a faza gotowa z pustym oknem znaczy, że rdzeń nie zna ani jednego okna tej
sesji — zlanie tych faz w jedną kazałoby zgadywać, czy czekać, czy działać.
Ogłoszenie zmiany ma jeden wspólny nasłuch, bo rejestr kanałów modelu leży poza magazynem i jego
danych nie dotyka, ale jego powrót ma przerysować stery, a wszystkie widoki modułu przerysowuje
jedno ogłoszenie stanu — drugi, równoległy nasłuch w oknach trzeba by odpinać w każdej instancji
panelu, a pierwszy przeoczony byłby wyciekiem.

## budowa/klient-poprzedni/src/moduly/roundtable/zrodlo-strumienia-debaty.ts
Rdzeń nadaje wypowiedź uczestnika na żywo, fragment po fragmencie, zanim utrwali ją i rozgłosi
po raz drugi jako zdarzenie zmiany debaty rodzaju „updated" — ten moduł odbiera tylko tamtą
nadawaną treść. Filtr działa po oknie, bo tylko okno jest w zdarzeniu strumienia pewne, więc
oddziela głosy tej debaty od strumieni okna rozmowy, podglądu w tle i każdego innego nadawcy
wspólnej drogi. Pole identyfikatora koperty zostaje surowe, bo jeden strumień niesie pod nim
dwie różne wartości — fragmenty treści dostają identyfikator uczestnika, a fragment domykający
i fragment błędu identyfikator wypowiedzi — a przypisanie do uczestnika należy do modułu
źródła wypowiedzi, który pyta o ten identyfikator stan debaty. Moduł nie woła żadnej komendy
obszaru roundtable i nie gromadzi treści, bo gromadzenie wymaga wiedzy o składzie uczestników,
której subskrypcja strumienia nie niesie. Parametr okna funkcji tworzącej źródło jest funkcją
odczytywaną w chwili nadejścia fragmentu, nie w chwili subskrypcji, żeby przełączenie modułu na
inną debatę nie skutkowało przepuszczaniem fragmentów debaty poprzedniej przez zapamiętany
filtr. Okno puste odrzuca cały ruch fragmentów, bo gospodarz bez okna nadanego przez rdzeń
inaczej otrzymywałby fragmenty okna rozmowy i podglądu w tle podpisane uczestnikami tej debaty.

## budowa/klient/src/polaczenie/warstwa-polaczenia.test.ts
Obietnica warstwy połączenia wobec warstw wyższych jest jedna: ramka wpisana przy rozłączeniu nie ginie, ponawianie nie ustaje, a subskrypcje przeżywają wymianę gniazda. Sprawdziany w tym pliku mierzą dokładnie to.

## budowa/klient-poprzedni/src/moduly/translate/mapowanie-stylow.ts
Ta sama zmienna zapisuje się na różnych platformach różnie, a przeniesienie materiału między nimi
wymaga przełożenia zapisu, nie treści. Przekład działa na tekście źródłowym, ponieważ to on jest
wzorcem dla wszystkich paneli — zamiana zapisu w jednym przekładzie bez zamiany go w źródle i w
pozostałych panelach dałaby materiał niespójny co do zmiennych, czyli dokładnie tę usterkę, której
szuka kontrola jakości. Wykaz zamian towarzyszy wynikowi zawsze, ponieważ przekład zapisu jest
zmianą niewidoczną w treści na pierwszy rzut oka, a od niej zależy podstawienie wartości w gotowym
produkcie.

## budowa/klient-poprzedni/src/rozmowa/nadanie-i-przerwanie.ts
Nadanie wiadomości i przerwanie tury rozmawiają z rdzeniem komendami message.send
i message.stop. Wysłanie w trakcie odpowiedzi samo prosi o przerwanie, więc obie
czynności są od siebie zależne i stoją w jednym pliku. Warstwa nie zna widoku —
dostaje wyłącznie haczyki: co ogłosić, jak zgłosić błąd, jak przestawić stan.
Dzięki temu ta sama logika obsługuje okno komunikacji i podgląd.

Przerwanie jest jawne: message.send skierowany do okna, które odpowiada, odmawia
kodem conflict zamiast skasować odpowiedź w połowie zdania, więc przerwanie
jedzie osobną komendą. Para komend idzie sekwencyjnie, nie równolegle — nadane
naraz dotarłyby w kolejności niegwarantowanej, a wysłanie, które wyprzedziło
zatrzymanie, trafiłoby na okno nadal zajęte i odmówiło. Nieudane zatrzymanie nie
wstrzymuje wysłania: jeżeli tura zdążyła tymczasem dobiec końca sama, wysłanie
przejdzie, a jeżeli nie — odmówi rdzeń.

## budowa/klient/src/polaczenie/zrodlo-zdarzen.ts
Zamiast importu kanału obserwator opisuje dokładnie to, czego potrzebuje: subskrypcję zdarzenia po nazwie z kontraktu i podgląd całego ruchu. Kanał spełnia ten opis kształtem, bez dodatkowej deklaracji — jeden byt odpowiada jednemu modułowi.

## budowa/klient-poprzedni/src/moduly/studio/agenci-obsada.ts
Strona okna pętli pokazuje dwa wykazy o tej samej pracy kilku wykonawców nad jednym dokumentem:
obsadę, czyli który wykonawca zajął który fragment, w jakim jest stanie i kiedy jego zajęcie
wygasa, oraz spięcia, czyli co się stało, gdy dwóch wykonawców sięgnęło po ten sam fragment —
czyja zmiana weszła, czyja została odłożona i wedle jakiej nastawy. Odłożone brzmienie nie
przepada: pole ze zmianą, która nie weszła, jest sednem tego widoku, nie ozdobą, bo praca
wykonawcy odłożona bez pokazania jej operatorowi byłaby pracą wyrzuconą po cichu. Brzmienie stoi
więc wprost do przeczytania, wraz z czynnością „Przyjmij brzmienie", która wnosi je do dokumentu
decyzją operatora. Gdzie rdzeń odłożył brzmienie jako propozycję albo zmianę śledzoną, widok
odsyła do niej po identyfikatorze, bo dwóch kopii tego samego brzmienia nie zakłada.

## budowa/klient-poprzedni/src/okno-komunikacji/dyktowanie/urzadzenia-dzwieku.ts
Pusty wykaz nie tłumaczy się sam: lista bez pozycji oznacza brak mikrofonu albo odmowę przeglądarki, a to dwie różne sytuacje z dwoma różnymi wyjściami dla użytkownika, dlatego pole `powod` towarzyszy wykazowi i bywa niepuste również wtedy, gdy urządzenia istnieją, lecz noszą nazwy zastępcze — pusty `powod` oznacza wykaz kompletny, niewymagający dopowiedzenia. Nazw urządzeń nie zgaduje się: przed pierwszą zgodą etykieta każdego urządzenia jest pusta, a numer porządkowy w miejsce nazwy jest przyznaniem się do niewiedzy, nie nazwą sprzętu.

## budowa/klient-poprzedni/src/moduly/translate/modul-translate.ts
Układ wynika z warstw widoczności opracowania: warstwa pierwsza stoi na ekranie bez interakcji —
Source Panel i Translation Panels obok siebie, bo zapis źródła aktualizuje wszystkie panele naraz
i skutek musi być widoczny w tej samej chwili. Format Studio jest kolumną sąsiadującą warstwy
drugiej, otwieraną przy pracy z dokumentem. Glossary Manager, Translation Memory Panel i QA
Review Center są zarządcami warstwy trzeciej, stoją w kolumnie bocznej i są zwinięte, dopóki ich
nie wywołać. Chat Window i Execution Loop Window są oknami wspólnymi platformy i moduł ich nie
buduje: pierwsze montuje scena sesji, drugie należy do pętli wykonawczej; katalog okien podaje
obok, ile okien modułu rdzeń zna, a ile moduł buduje, licząc tę liczbę z odpowiedzi rdzenia.
Drogi warstwy czwartej są trzy i wszystkie prowadzą do tych samych rzeczy: skrót klawiszowy,
wyszukiwarka funkcji z pełnym katalogiem opracowania oraz rozwinięcia przy oknach, których
dotyczą — zasada jednego kliknięcia jest przez to spełniona także dla pozycji bez własnego
przycisku w oknie.
Sesja ostatniego wejścia jest potrzebna do ponowienia odczytu kontekstu, bo trzeba wiedzieć,
o czyje okna pytać, a moduł nie sięga po sesję sam — dostaje ją z powłoki. Warsztat prowadzi
czynności wsadowe i wymianę z otoczeniem, a nie bieżący przekład, dlatego stoi zwinięty tak jak
trzy okna obok, dopóki operator go nie wywoła. Pasek uczciwości bierze zdanie z bytu wspólnego
katalogu okien, tego samego dla wszystkich modułów, więc zdanie o rozjeździe nie rozjedzie się
między modułami po cichu — byt wypowiada obie strony: okna katalogu rdzenia, których moduł nie
buduje, oraz okna budowane spoza katalogu. Jeden odczyt katalogu okien starcza na cały kanał: gdy
o katalog zapytał już inny moduł albo powłoka, wywołanie nie wysyła drugiego zapytania. Jedyny
odczyt wykonywany bez czynności operatora w oknach warstwy pierwszej dotyczy macierzy izolacji —
idzie obok, bo dotyczy procesu sesji, a jej odmowa nie zatrzymuje modułu, ponieważ powód zapisuje
sama macierz. Bez odpięcia paska uczciwości od wspólnego katalogu przy zejściu modułu wpis kanału
trzymałby przerysowanie elementu zdjętego już z drzewa. Stery kanału zwijają się przed zejściem
modułu, ponieważ rozwinięte trzymają nasłuch na dokumencie, którego zniknięcie elementu nie zdejmuje.

## budowa/klient-poprzedni/src/powloka/podglad.ts
Strona wzorowana jest na podglądzie zetonów motywu. Jedna odpowiedzialność: uruchomienie powłoki
poza aplikacją, żeby cztery pasy dało się obejrzeć i przeklikać przed ich osadzeniem w punkcie
wejścia. Podgląd niczego nie udaje: karty sesji i ich stany są tu materiałem pokazowym, a nie
danymi z rdzenia — dlatego mieszkają na tej stronie, a nie w module powłoki. Uruchomienie: npm run
dev, adres /src/powloka/podglad.html. Motyw wybiera się parametrem ?motyw=jasny albo
?motyw=ciemny; bez parametru rozstrzyga zapisany wybór, a w jego braku preferencja systemu.

## budowa/klient-poprzedni/src/powloka/postac-pasa-kart.ts
Nazywa przyczynę braku karty w pasie oraz treść odpowiedzi rdzenia na czynność Operatora.

## budowa/klient-poprzedni/src/moduly/research/pustka-okien.ts
Stan przed pierwszym odczytem rozróżnia „jeszcze nie pytałem" od „rdzeń nic nie ma": twierdzenie
o zawartości rdzenia postawione bez jego odpowiedzi byłoby zgadywaniem, więc moduł nazywa wprost,
że jeszcze nie zapytał, zamiast domniemywać pustkę.

## budowa/klient-poprzedni/src/moduly/translate/narzedzia-panelu.ts
Narzędzia są wydzielone z panelu, bo panel odpowiada za treść tłumaczenia i jej korektę, a to jest
pasek czynności wykonywanych na tej treści, przy czym każde narzędzie ma własną komendę kontraktu.
Wynik narzędzia nie przesłania panelu: kontrola jakości, tłumaczenie zwrotne i podpowiedź pamięci
pokazują się w wierszu odpowiedzi pod paskiem, żeby dało się je z panelem porównać. Ster kanału
stoi przy pasku, bo dotyczy jednej z jego czynności — pole kanału tłumaczenia zwrotnego jest
jedynym polem kanału w tym pasku, pozostałe pięć komend modelu nie wywołuje albo woła go bez
wskazania, więc podpis steru nazywa czynność wprost, zamiast udawać nastawę całego panelu.
Treść tłumaczenia jest potrzebna tłumaczeniu zwrotnemu, ponieważ rdzeń bywa, że oddaje w kolumnie
zwrotnej dokładnie tę treść, a wtedy przekładu odwrotnego nie było i sprawozdanie ma to powiedzieć
wprost. Po przerysowaniu paska pole tonu panelu niesie tę samą wartość, którą pokazuje nagłówek
panelu. Ta sama zasada rządzi polem tłumaczenia w pliku panel-jezyka.ts: kontrolka, w której
operator właśnie pisze, należy do niego, dopóki jej nie odda, a zdarzenie zmiany tłumaczenia,
przychodzące w środku pisania, nie może podmienić wpisywanego tonu.

## budowa/klient-poprzedni/src/moduly/research/rama-badania.ts
Obudowa okna jest jedna dla całego interfejsu, ale wykaz akcji i ich dymki objaśniające należą do
modułu — panel akcji wnosi generyk akcji, którego rama sama nie potrzebuje; okna Research składają
te dwie rzeczy w jednym miejscu, więc przedrostek klas modułu i sposób osadzenia panelu stoją razem.
Przedrostek `mr` nie wnosi wyglądu, jest uchwytem jednej reguły własnej modułu — obwiedzenia okna,
do którego nawigacja wewnątrzmodułowa przeniosła ognisko. Rama nie buduje elementów treści, dostaje
je gotowe, tak jak przestrzeń modułu dostaje gotowe okna.

## budowa/klient-poprzedni/src/moduly/studio/aparat-panel.test.ts
Testy sprawdzają cztery reguły panelu aparatu dokumentu i pola: spis treści zakłada się
z poziomów nagłówków wskazanych przez operatora, a numer nadaje rdzeń, więc okno nie numeruje
niczego samo; element wymagający odświeżenia jest nazwany liczbą i znacznikiem w danych, więc
spis treści rozjechany z dokumentem nie może wyglądać na zgodny; pole bez policzonej wartości
mówi to wprost zamiast pokazywać pustkę; usunięcie odmówione nie zdejmuje elementu z wykazu.

## budowa/klient-poprzedni/src/moduly/translate/odczyty-zrodla.ts
Oba odczyty biorą się z tego samego napisu i z jednego przebiegu, więc stoją w jednym pliku, i oba
są warstwą pierwszą Source Panel — mają być widoczne bez interakcji, ponieważ orientacja w
rozmiarze materiału i wiedza o tym, czego nie wolno przetłumaczyć, poprzedzają każdą czynność w
tym oknie. Liczby są liczbami okna, nie rdzenia: rdzeń oddaje liczbę pozycji podziału przy zapisie
źródła i to jest jego prawda o segmentach, a słowa i znaki liczy okno z tekstu, który ma przed
sobą, bo kontrakt takiej komendy nie ma. Analizy względem pamięci tłumaczeń ani wyceny nie ma tu
wcale, ponieważ nie ma z czego ich złożyć, a okno tego nie zastępuje szacunkiem.

## budowa/klient-poprzedni/src/strona-glowna/naglowek-strony.ts
Trójka napisów pochodzi z makiety Centrum dowodzenia i ze schematu strony w warstwie projektowej. Znaku marki
nagłówek nie powtarza: godło i nazwa produktu stoją już w pasku górnym, a strona ma nad strefami tytuł czytany,
nie drugi sygnet. Stopień tytułu jest niższy od stopnia tytułów kart środowisk z rozmysłu: środkiem ciężkości
strony pozostaje strefa pierwsza, nie napis nad nią.

## budowa/klient-poprzedni/src/moduly/studio/aparat-zrodlo.ts
Aparat dokumentu niesie trzynaście rodzajów elementów — spis treści, spis ilustracji i tabel,
przypis dolny i końcowy, podpis, zakładkę, odwołanie wzajemne, odsyłacz, powołanie, bibliografię,
hasło indeksu i indeks — a pola dokumentu niosą numer strony, liczbę stron, datę, godzinę,
tytuł i autora dokumentu, właściwość oraz pole obliczane. Obie rodziny dzielą tę samą oś:
element wyliczany z dokumentu, który po zmianie treści staje się nieświeży i wymaga odświeżenia,
dlatego panel prowadzi dla obu jeden wspólny wykaz do odświeżenia zamiast dwóch osobnych.

## budowa/klient/src/protokol/kanal.ts
Kanał nie zna treści dziedzinowej. Nazwy komend i zdarzeń oraz kształty ich treści pochodzą wyłącznie z kontraktu współdzielonego: zmiana nazwy w kontrakcie przerywa kompilację klienta. Stąd jedno wejście wysyłające obsługuje każdą komendę kontraktu, a jedno wejście subskrybujące — każde jego zdarzenie; typ treści wyznacza nazwa. Komunikat nierozpoznany nie jest odrzucany.
Dziennik zakładany jest od razu, bo komunikat nierozpoznany może przyjść przed pierwszą subskrypcją warstwy wyższej. Zapis i wpis do dziennika są jedyną reakcją: ani zdarzenie nieznane, ani koperta o typie spoza kontraktu nie zrywa połączenia i nie blokuje sesji.

## budowa/klient-poprzedni/src/strona-glowna/wykonanie-nazwy.ts
Trzy czynności stoją osobno od pozostałych, bo mają wspólny kształt — pytanie, odmowa, komenda, meldunek —
i wspólny warunek zejścia: zamknięte pytanie kończy rzecz bez żadnego żądania. Zamknięcie pytania zwraca
brak, nie napis pusty. Odmowa operatora nie jest odmową rdzenia i nie ma o niej czego meldować, a napis
pusty bywa odpowiedzią sensowną — przy kopii oddaje nadanie nazwy rdzeniowi — więc te dwa przypadki nie
mogą się zlać.

## budowa/klient-poprzedni/src/strona-glowna/wpiecie-srodowisk.ts
Źródłem prawdy jest wykaz środowisk pobierany z rdzenia, nie stała klienta. Wykaz zastany w kliencie
zostaje jako wartość początkowa: odpowiedź rdzenia przychodzi po pierwszym rysowaniu i wtedy przerysowuje
strefę. Odmowa i wykaz pusty nie gaszą ekranu — na miejscu zostaje wykaz zastany, po którym da się wejść
do pracy.

## budowa/klient-poprzedni/src/moduly/research/sekcje-raportu.ts
Sekcje jadą w polu `sections` komendy `research.report.build`, więc redakcja nie potrzebuje osobnej
komendy; kolejność bierze się z porządku wierszy i wchodzi w pole `order`. Sekcja bez tytułu nie
blokuje zapisu — kreator wysyła to, co ma, a brak nazywa komunikatem przy polu. Wiersz pusty w obu
polach nie jest sekcją i tu wypada. Podanie `sections` przestawia rdzeń na gałąź zapisu wprost: nie
woła modelu i nie zamienia zaznaczonych ustaleń na sekcje. Pusta tablica też jest podanymi sekcjami,
więc warstwa żądania pomija pole `sections` w całości, gdy nie zostało nic — inaczej kreator, który
dokłada pusty wiersz przy otwarciu, składałby raport z jednej sekcji bez tytułu i bez treści zamiast
streszczenia ustaleń. Wiersz wczytany z raportu niesie dalej swoje `findingIds`, więc powtórne
złożenie nie zrywa wiązania ustalenie-sekcja.

Identyfikator sekcji musi być niepowtarzalny poza sesją okna, nie tylko w niej: kolumna
`sekcja_raportu_badania.identyfikator_zewnetrzny` jest unikalna globalnie, a nie w obrębie raportu,
więc kod powtórzony przy nowym raporcie kończy się odmową rdzenia. Ten sam kod przy tym samym
raporcie jest w porządku — rdzeń poprawia sekcję zastaną. Człon losowy pochodzi ze wspólnego
generatora identyfikatorów, żeby nie hodować drugiej reguły nadawania identyfikatorów w kliencie.

## budowa/klient/src/protokol/koperta.ts
Plik nie definiuje własnego typu komunikatu i nie powiela ani jednego literału nazwy — kształt koperty pochodzi wyłącznie ze współdzielonego kontraktu.

## budowa/klient-poprzedni/src/moduly/studio/blokada-fragmentow.test.ts
Sprawdziany mierzą trzy zachowania, których zgubienie przeszłoby niezauważone: zmiana obejmująca
blokadę częściowo oddaje bilans akcji wraz z nazwą blokady zamiast milczeć o niej; odmowa zajęcia
fragmentu nazywa wykonawcę i czas zajęcia zamiast wracać pustą odpowiedzią; zasięg obejmujący
właściciela pisma jest opisany jako ustawienie jawne, bo blokada domyślnie jest skierowana
przeciw modelowi, nie przeciw właścicielowi pisma.

## budowa/klient/src/protokol/korelacja.ts
Rejestr nie wprowadza limitu czasu ani limitu żądań oczekujących: zerwanie połączenia nie kończy pracy rdzenia nad poleceniem.

## budowa/klient-poprzedni/src/moduly/translate/okno-format-studio.ts
Okno stoi na styku dwóch obszarów kontraktu i tylko dlatego ma czym pracować: dokument wchodzi
komendami obszaru document — wydobycie tekstu wraz z rozpoznaniem pisma oraz zamiana formatu —
a wychodzi komendą eksportu panelu z obszaru translate. Obieg jest więc realny na obu końcach, ale
nie jest obiegiem zamkniętym: wydobycie oddaje sam tekst, więc styl, tabela i osadzenie dokumentu
wejściowego nie mają jak przetrwać przekładu, a okno mówi to przy wyniku, zamiast obiecywać
wierność formatu. Podgląd jest warstwą tekstową, nie renderem strony: kontrakt nie ma komendy
rysującej dokument ani porównującej jego układ, więc porównania układów okno nie pokazuje i nie
udaje. Tekst wydobyty z dokumentu wchodzi do pola Source Panel, ale nie zapisuje się sam — zapis
źródła jest czynnością operatora i uruchamia aktualizację wszystkich paneli, a wykonanie go bez
naciśnięcia byłoby przekładem zamówionym przez okno, nie przez człowieka.
Zdanie o wyniku wydobycia tekstu rozróżnia dwie drogi, którymi tekst mógł powstać, bo różnią się
pewnością: warstwa tekstowa dokumentu jest odczytem, rozpoznanie pisma jest odgadnięciem z
pikseli, a odpowiedź mówi, która droga zaszła, i okno tego nie zaciera.
Czynność zdejmowania poprzedniej odmowy jest miejscowa i nie pyta rdzenia, więc nie stawia okna
w ładowaniu, tylko zdejmuje komunikat poprzedniej odmowy, bo od tej chwili okno pokazuje wynik,
a nie powód niewykonania.

## budowa/klient-poprzedni/src/rozmowa/nadawca.ts
Rodzaje nadawcy rozróżniają się trzema nośnikami naraz i żaden z nich nie jest
barwą tła: ikona w medalionie pierwszej kolumny wpisu, barwa medalionu i kreski
krawędzi ustalana wyłącznie w rozdzielczości klasy semantycznej, nie rodzaju,
oraz etykieta słowna wersalikami. Barwa sama nie może być jedynym nośnikiem
znaczenia, więc tło wpisu pozostaje jedno dla wszystkich dziewięciu rodzajów.

Arkusz stylu wpisu zna dokładnie trzy modyfikatory klasy semantycznej i wiąże
z nimi barwę kreski oraz barwę medalionu: człowiek dostaje barwę atramentu
(Operator), inteligencja barwę sygnału (model, agent, koordynator, wykonawca,
walidator), system barwę neutralną (automatyzacja, Always On Display, wynik
narzędzia). Rodzaj wewnątrz klasy różnicuje ikona i etykieta, nigdy kolor.

Rozpoznawanie nadawcy łączy rolę wiadomości z kontraktu, która zna cztery
wartości, z rolą okna rozdzielającą wypowiedź modelu na koordynatora
i wykonawcę pętli. Rola systemowa oznacza wypowiedź warstwy automatycznej
platformy, a nie żadnego z ośmiu pozostałych nadawców.

## budowa/klient/src/protokol/ksztalt-odpowiedzi.ts
Rzutowanie jest obietnicą kompilatora, nie rdzenia — rdzeń starszej wersji albo pośrednik może przysłać treść bez pola obowiązkowego, a wołający dostałby wartość niezdefiniowaną w miejscu, w którym typ obiecuje wartość. Sprawdzian zamienia taką odpowiedź w niepowodzenie wywołania z wpisem do dziennika i błędem walidacji.

## budowa/klient-poprzedni/src/powloka/potwierdzenie-usuniecia.ts
Nazwy komendy ten plik nie zna — wysyłkę dostaje z zewnątrz, więc daje się poddać próbie bez rdzenia.
Operator widzi stratę, zanim kliknie. Kontrakt żąda pola `confirm`, a rdzeń bez niego odmawia
wykonania. Potwierdzenie ma więc treść, a nie samo pytanie „czy jesteś pewien": modal wypisuje tytuły
wskazanych sesji i mówi wprost, że zapis ginie razem z wiadomościami, oknami i artefaktami. Okno stoi
na natywnym elemencie dialogowym, wzorem okna pytania o nazwę ze strony głównej, więc warstwa tła,
pułapka ogniska i Escape należą do przeglądarki. Escape w trakcie wywołania jest wstrzymany: rdzeń już
usuwa sesje, więc zamknięcie okna nie odwołałoby niczego, a Operator zostałby bez odpowiedzi. Trzy
stany obowiązkowe stoją na jednym pasie fazy okna z biblioteki komponentów: stan ładowania w czasie
wywołania, stan pusty gdy rdzeń nie usunął niczego, stan błędu przy odmowie — z kodem i treścią
rdzenia z opisu odmowy. Odmowa nie zamyka modalu: przycisk czynności znika, zostaje samo „Zamknij",
więc nie da się wziąć odmowy za skutek.

Oddaje rozliczenie rdzenia, gdy komenda przeszła — także rozliczenie, w którym nic nie zginęło, bo to
też jest odpowiedź rdzenia i pas kart ma ją powtórzyć. Oddaje `null`, gdy Operator odmówił
potwierdzenia albo gdy rdzeń odmówił wykonania; treść odmowy została wtedy pokazana w modalu.
Obietnica rozstrzyga się z odpowiedzią rdzenia, nie z zamknięciem okna. Modal zostaje otwarty do
przeczytania skutku, ale pas kart ma odpowiedź natychmiast — wiązanie rozstrzygnięcia z zamknięciem
okna kazałoby pasowi milczeć tak długo, jak długo Operator czyta.

Przed powtórzeniem czynności broni strażnik stanu w toku w obsłudze kliknięcia, a przed zamknięciem
okna ten sam strażnik w funkcji zdejmującej okno, więc znacznik zablokowania kontrolki nie dokładałby
ochrony — dokładałby ciszę. Zamiast tego przycisk odpowiada zdaniem, dlaczego w tej chwili nie ma czego
zrobić. Powód idzie dwiema drogami: tytułem pod kursorem i opisem dostępności dla czytnika ekranu.

Zdanie o bezskuteczności kliknięcia idzie do akapitu skutku oznaczonego rolą stanu, więc czytnik
ekranu ogłasza je od razu, a pas stanu dalej trzyma fazę ładowania. Milczenie byłoby tu gorsze od
blokady: Operator wziąłby brak reakcji za zawieszenie.

## budowa/klient-poprzedni/src/strona-glowna/wpiecie-modulow.ts
Zapytanie o wykaz modułów oddaje komplet modułów platformy wraz z polem zawężenia widoczności, którego
pustka znaczy, że moduł jest dostępny wyłącznie ze strony głównej. Gdy rdzeń nie odpowie, strefa zostaje
ukryta zamiast stać z nagłówkiem nad pustką: brak odpowiedzi nie jest wykazem pustym. Przejście idzie
z samym kodem modułu, bez środowiska — moduł w żadnym nie stoi, więc powłoka otwiera go z pominięciem
macierzy widoczności.

## budowa/klient/src/protokol/powitanie.ts
Dopiero z odpowiedzi klient dowiaduje się, jaką wersję protokołu zna rdzeń i które komendy ta wersja rdzenia obsługuje. Rozstrzygnięcie, co zrobić z rozjazdem wersji, należy do warstwy wyższej — protokół oddaje odpowiedź w kształcie kontraktu.
Wersja protokołu w żądaniu pochodzi ze stałej kontraktu, nie z literału: klient przedstawia się tą wersją, z którą został zbudowany.
Token wiąże połączenie z sesją bramki. Jego brak nie jest błędem — rdzeń odpowiada wtedy brakiem uwierzytelnienia. Token dostarcza wołający, bo magazyn sesji bramki nie należy do warstwy protokołu.

## budowa/klient-poprzedni/src/moduly/studio/blokada-fragmentow.ts
Blokada fragmentu i zajęcie fragmentu to dwie różne rzeczy, obie utrzymywane w rdzeniu. Sprawdzenie
blokady stoi na drodze każdej komendy zmieniającej dokument po stronie serwera, przed dotknięciem
treści — blokada pilnowana przez okno byłaby pozorna, bo model woła komendy rdzenia tak samo jak
klient i ominąłby ją bez wysiłku, więc okno tylko pokazuje blokady i ich skutki, nie wykonuje ich.
Blokada jest trwała i skierowana przeciw modelowi: zdejmuje ją wyłącznie operator, a model, który
uzna zmianę za konieczną, zakłada propozycję na marginesie. Zajęcie jest czasowe i skierowane
przeciw drugiemu wykonawcy, żeby dwóch agentów nie pisało po tym samym akapicie — wygasa samo,
a odmowa zajęcia nazywa wykonawcę i czas. Zmiana obejmująca blokadę częściowo wykonuje się poza
blokadą i oddaje bilans: co przeszło, co pominięte i przez którą blokadę, bo przemilczenie pominięcia
jest zakazane, choć odmowa całości zmiany byłaby nieproporcjonalna. Plik nie zna elementów strony
ani wywołań rdzenia.

## budowa/klient-poprzedni/src/moduly/research/skutek-zlozenia.ts
Wydzielone z pliku `czynnosci-raportu.ts`, bo tamten plik prowadzi czynność, a ten jest przekładem
odpowiedzi `research.report.build` na zdanie dla Operatora — dwie rzeczy zmieniane z różnych
powodów, na wzór pliku `assistant/skutek-sterowania.ts`. Znaku „to nie model" w odpowiedzi nie ma:
sekcja niesie komplet pól kontraktu i nic ponad to (`id, title, content, findingIds, order`),
a rdzeń zbiera z kanału same fragmenty tekstu (w `core/adapter_modul_badania.go`, funkcja
`zapytajModel`) — komunikat procesu kanału i zdanie modelu docierają więc tą samą drogą, jako ta
sama treść. Rozpoznanie po napisie jest zakazane, bo napis się zmienia: okno nie orzeka, że
streszczenie jest komunikatem błędu, tylko mówi, czego nie wie, i każe przeczytać treść przed
wydaniem raportu. Znak, który odróżnia naprawdę: sekcja, której identyfikatora okno nie wysłało,
jest sekcją napisaną przez rdzeń — kryterium bierze się z odpowiedzi, a nie z brzmienia tytułu.
Druga rozbieżność widoczna z odpowiedzi: przy podanym polu `sections` rdzeń odkłada `findingIds`
na bok, więc zaznaczenie ustaleń przepada bez śladu, a okno ogłasza to odmową, nie przemilcza.

Tytuł pusty nazywany jest wprost pustym — cudzysłów bez treści wygląda na usterkę wypisywania,
a jest wiernym oddaniem tego, co wróciło, bo pole title puste nie idzie do rdzenia i raport
został bez nazwy; klient jej nie dopowiada za Operatora.

Rozpoznanie sekcji napisanej przez rdzeń bez udziału Operatora ma granicę: zdanie pada wyłącznie
przy złożeniu, w którym sekcja przyszła z rdzenia po raz pierwszy. Kolejne złożenie wysyła ją już
z redakcji kreatora, pod tym samym identyfikatorem, bo wczytanie przejmuje sekcje potwierdzone —
więc kryterium przestaje ją wskazywać i słusznie: Operator miał ją wtedy przed sobą w polach edycji.

## budowa/klient-poprzedni/src/powloka/powloka-natywna.ts
Interfejs Danaco Console działa w dwóch miejscach naraz: w oknie powłoki natywnej i w zwykłej
przeglądarce. Każdy most do powłoki musi najpierw ustalić, czy powłoka w ogóle jest — inaczej
wywołanie IPC rzuca wyjątkiem tam, gdzie żadnego IPC nie ma. Pytanie stoi w osobnym pliku, żeby most
rdzenia nie zależał od mostu katalogów tylko po to, by je zadać; oba pytają tak samo. Brak powłoki nie
jest błędem: sprawdzenie wykonane poza powłoką ma prawo rzucić wyjątkiem, a wtedy odpowiedzią jest
zwyczajne „nie ma powłoki".

## budowa/klient/src/protokol/ramka.ts
Ramka nieczytelna albo o kształcie niezgodnym z kopertą nie blokuje sesji — wraca jako zdarzenie zapasowe z zachowaniem treści surowej w ładunku.

## budowa/klient-poprzedni/src/okno-komunikacji/dyktowanie/wynik-dyktowania.ts
Cisza jest prawidłowym wynikiem pomiaru: nagranie, w którym nie padło słowo, przeszło przez silnik tak samo jak nagranie z pełną wypowiedzią, dlatego stan bez mowy ma własny stan i własne zdanie zamiast być zlewany z odmową, co pokazywałoby awarię tam, gdzie jej nie ma, i rozmywałoby odmowę prawdziwą — brak modelu, brak silnika rozpoznawania albo brak drogi dostarczenia nagrania. Pole powodu wypełnia się tylko przy stanie nieprzetworzonym; przy dwóch pozostałych stanach pomiar się odbył i nie ma czego uzasadniać. Zdanie dla stanu bez mowy nazywa fakt, że nagranie przetworzono, a mowy w nim nie było, bez słowa błąd i bez wezwania do ponownego nagrania, bo cisza nie jest usterką.

## budowa/klient-poprzedni/src/moduly/translate/okno-glossary-manager.ts
Kontrakt nie ma komendy odczytu glosariusza — w obszarze translate nie występuje ani lista, ani
pobranie pojedynczego terminu. Okno pokazuje więc wyłącznie terminy zapisane w tej sesji i mówi
o tym wprost w stanie pustym; wykaz zmyślony albo pusta lista podana jako „glosariusz jest pusty"
byłyby atrapą. Wyszukiwanie działa na tym, co okno ma, czyli na terminach zapisanych w tej sesji,
i nie udaje, że przeszukuje bazę — przesiewa wykaz widoczny obok. Przerysowanie odbudowuje wykaz
terminów, ale nie rusza fazy trwającej ani fazy błędu — te zdejmuje czynność, która je postawiła
(zapis terminu, odczyt wystąpień); inaczej zapis, który napełnia wykaz w środku własnego wywołania,
kasowałby sobie zapowiedź tego wywołania.
Wykaz pusty przy pokazaniu wystąpień jest wynikiem, nie pustką okna: rdzeń odpowiedział i nie
znalazł terminu, a odmowa czyści wykaz i mówi wprost, że wystąpień nie sprawdzono — te dwa stany
nie mogą wyglądać tak samo. Zdanie nazywa termin, którego rdzeń szukał, biorąc go z pola term
odpowiedzi, a nie z żądania: zdanie zbudowane z żądania przypisywałoby rdzeniowi przeszukanie
o zakresie, którego okno nie widziało, a rozbieżność echa jest odmową, bo znaczy, że szukano
czegoś innego.

## budowa/klient-poprzedni/src/moduly/wykaz-komend-rdzenia.ts
Wykaz jest warstwą niższą niż pokrycie komend: nie zna dokumentu ani
kontrolki, zna tylko odpowiedź rdzenia i zdania, które z niej wynikają.
Prawda bierze się wyłącznie z powitania połączenia, nie z wykazu komend
zapisanego w kontrakcie, bo rdzeń oddaje w odpowiedzi powitania komendy,
które naprawdę obsługuje. Obsługiwacz powitania nie czyta żądania i niczego
nie zapisuje, więc powtórne powitanie jest czystym odczytem, nie drugim
uzgodnieniem połączenia. Samych komend świadomie się nie pyta wywołaniem
niszczącym tylko po to, by zobaczyć odmowę.

Dwa pierwsze z pięciu rozstrzygnięć pokrycia rozstrzyga sam kontrakt, bez
pytania rdzenia, więc są widoczne jeszcze przed odpowiedzią powitania.
Jedno powitanie przypada na połączenie, nie na moduł: wykaz leży w pamięci
podręcznej przypisanej do kanału, więc pierwszy odczyt pyta, a pozostałe
czekają na tę samą odpowiedź; odmowa w pamięci nie zostaje, więc kolejne
wywołanie ponawia pytanie. Wykaz zmienia się z wersją rdzenia, nie w toku
sesji, ale transport ponawia połączenie pod tym samym kanałem, więc po
zerwaniu można trafić na rdzeń inny niż odczytany — dlatego nasłuch
powitania obejmuje cały kanał, także powitania cudze, i odświeża wykaz bez
dodatkowego zapytania.

Tożsamość klienta zgłaszana w powitaniu jest stała, a nie świeżo nadawana
przy każdym wywołaniu, ponieważ obsługiwacz powitania nie czyta żądania i
nowy identyfikator rozdzieliłby ognisko od połączenia, które je zgłosiło.

## budowa/klient-poprzedni/src/strona-glowna/wpiecie-sesji.ts
Ten plik jest jedynym w katalogu znającym kanał i czynność powrotu naraz. Powrót wiąże połączenie i po nim
przenosi ognisko: powiązanie odtwarza okna i kieruje ich strumienie na to połączenie, a odmowa przeniesienia
ogniska nie cofa powiązania, bo sesja jest już związana. Oba żądania niosą tożsamość klienta z powitania
połączenia, więc tożsamości nie nadaje się tu po raz drugi. Montaż bez tożsamości klienta dostaje wykaz
czysto informacyjny, bez przycisku powrotu; pozostałe czynności historii — nazwa, kopia, projekt, archiwum
i bieg sesji — idą na sam identyfikator sesji, więc wpinane są przed wyjściem po braku tożsamości: wykaz bez
powrotu nadal daje się porządkować.

## budowa/klient/src/protokol/rozmowa-z-rdzeniem.test.ts
Bez rdzenia nasłuchującego pod adresem lokalnym nie ma czego zmierzyć. Milczenie rdzenia kończy się tu niepowodzeniem nazywającym przeszkodę, nie pominięciem: sprawdzian, który sam siebie odpuszcza przy braku rdzenia, wygląda potem tak samo jak sprawdzian zdany.

## budowa/klient-poprzedni/src/moduly/research/stan-badania.ts
Gdyby każde z siedmiu okien modułu prowadziło swój zbiór, powiązanie źródło-ustalenie-sekcja
raportu dotyczyłoby sześciu różnych bytów zamiast jednego. Wykazy narastają z odpowiedzi: komendy
odczytu wykazu źródeł i ustaleń są już w kontrakcie, ale rdzeń nie ma dla nich uchwytów, więc stan
trzyma to, co rdzeń potwierdził w tym połączeniu, i nie dopowiada reszty — po dobudowie uchwytów
odczyt dołoży się tutaj, obok odświeżenia, bo kształt stanu tego nie wymaga. Plik jest złożeniem
pamięci, odczytu i dwóch źródeł komend; sam nie trzyma ani jednej wartości.

Zaznaczenie jest nastawą wspólną oknom, więc ogłasza się tą samą drogą co treść. Pamięć go nie zna —
nie pochodzi z rdzenia i nie jest treścią badania — więc ogłoszenie ma własny rejestr słuchaczy,
a funkcja obserwująca wpisuje słuchacza do obu. Wołający ma jedną subskrypcję na cały stan i nie
musi wiedzieć, która zmiana skąd pochodzi.

Zdarzenie rdzenia jest jedynym odświeżeniem poza własnym działaniem: raport zmieniony na innym
urządzeniu konta dociera tą samą drogą. Rdzeń rozgłasza zdarzenie zmiany raportu bez wskazania
sesji, więc dochodzi ono na każde połączenie — porównanie pola okna badania w raporcie z oknem
badania odsiewa raporty cudzych okien; bez tego porównania Export Panel wydałby dokument, którego
to okno nie składało.

## budowa/klient-poprzedni/src/rozmowa/odczyt-fragmentu.ts
Pole `data` kontraktu jest typu `unknown` — rdzeń pakuje w nie strukturę
właściwą rodzajowi fragmentu. Odczyt jest w całości tolerancyjny: pole
brakujące, pole innego typu ani ładunek nieznanego kształtu nie przerywają
strumienia. Brak wartości znaczy „nie wiem", nie „błąd".

## budowa/klient/src/protokol/sesja.ts
Kontrakt dopuszcza kopertę bez identyfikatora sesji — pole jest opcjonalne, puste dla powitania połączenia.

## budowa/klient-poprzedni/src/strona-glowna/strefa-zwijana.ts
Strefa jest zbudowana na elemencie szczegółów z podsumowaniem: postać rozwinięcia trzyma przeglądarka,
a dane strefy czytane są dopiero przy rozwinięciu. Tak samo działa archiwum sesji; ten plik uogólnia zabieg
na całe strefy i dokłada pamięć postaci między wejściami, żeby przywołanie nie powtarzało się przy każdym
wejściu na stronę. Strefa zwinięta nie jest strefą ukrytą: zapowiedź — nazwa, zdanie wyjaśnienia i dopisek
z liczbą — stoi na ekranie zawsze i mówi, co jest pod spodem, więc jedno naciśnięcie otwiera treść znaną
z opisu. Rozwinięcia nie dubluje żadna klasa stylu; jedynym jego nośnikiem jest atrybut elementu, bo druga
kopia stanu mogłaby się z pierwszą wyłącznie rozminąć.

## budowa/klient-poprzedni/src/moduly/studio/braki-cyfryzacji.ts
Plik rozróżnia dwa różne braki, których nie wolno zlewać: brak pozycji kontraktu, czyli funkcja
nieobecna w produkcie i okno nie ma czym jej wywołać — to brak do zgłoszenia — oraz brak
składnika pakietu serwera, czyli funkcja zbudowana, a serwer, na którym stoi, jest niekompletny —
to usterka wdrożenia, nie ograniczenie produktu. Ustalono, że rozpoznawanie tekstu, rozpakowywacz
archiwów i silnik mowy idą wraz z aplikacją na serwer, więc arsenał stoi na serwerze, a u operatora
jest tylko cienka instalka; rozpoznanie tekstu i rozpakowanie archiwum nie są więc brakiem
produktu, a odmowa z ich powodu ma nazwać brakujący składnik pakietu serwera. Poprzednia postać
tego pliku wyliczała osiem czynności jako robotę do wykonania w oknie, bo panel stał jeszcze na
starszej komendzie wydobycia tekstu z dokumentu; okno prowadzi dziś całą rodzinę komend cyfryzacji
i zdania w pliku opisują funkcje działające oraz to, na czym stoją, żeby operator wiedział, czym
sterują nastawy, a nie czego brakuje.

## budowa/klient-poprzedni/src/powloka/powloka.ts
Jedna odpowiedzialność: złożenie czterech pasów w jeden układ i związanie ich zdarzeniami. Żaden pas
nie zna pozostałych; wiedzę o ich współpracy trzyma wyłącznie ten plik. Układ czterech pasów: pasek
górny 48 px na atramencie ramy — jedyny pas nieprzełączający się z motywem, poziome karty sesji
o mechanice zakładek, pionowa stała nawigacja modułów środowiska, obszar roboczy wypełniany przez
inne widoki. Wybór modułu przeładowuje kartę sesji: tytuł karty czynnej równa się nazwie otwartego
modułu, kontekst paska pokazuje parę środowisko-moduł, a obszar roboczy zapowiada moduł, którego
okna wejdą w jego miejsce.

Metoda podłączania zaczepów jest osobną metodą, a nie polem opcji powłoki, z tego samego powodu co
podłączenie źródła nawigacji: powłoka powstaje zanim istnieje połączenie z rdzeniem, a czynności bez
kanału podać się nie da. Wywołanie jest jedną linijką dokładaną w widoku trasy i niczego w nim nie
przestawia. Bez wywołania nic się nie psuje: przyciski obecności, pozycje menu profilu i wyniki
wyszukiwania zostają klikalne i mówią, że to pasek nie dostał drogi — nie że okien nie ma.

Powłoka aplikacji żyje tyle, co okno, więc w produkcie zamknięcie nie zachodzi ani razu — ale nasłuch
na dokumencie bez drogi zdjęcia jest wyciekiem, gdy powłoka powstaje wielokrotnie, jak w podglądzie
i w sprawdzianach.

Materiał wyszukiwania: to jedyne miejsce widzące naraz boczną nawigację, z wykazem modułów pobranym
przez komendę wykazu modułów, i pas kart sesji, więc źródło składa się tutaj — tak samo jak tutaj
wiąże się kontekst paska z wyborem modułu. Nowego odczytu z rdzenia nie ma: druga komenda wykazu
modułów obok tej, którą zrobiła już nawigacja, byłaby drugą prawdą o jednym wykazie, a dwie prawdy
rozjeżdżają się przy pierwszej zmianie po stronie rdzenia.

## budowa/klient-poprzedni/src/okno-komunikacji/historia.ts
Klasa neutralna dzielona z systemem oraz ikona pozostają rozróżnione celowo: klasa mówi o randze wypowiedzi, a ikona o jej źródle, więc wynik narzędzia nosi tę samą klasę co komunikat systemowy, lecz nigdy tę samą ikonę.

## budowa/klient-poprzedni/src/moduly/translate/okno-qa-review.ts
Okno robi to, czego pojedynczy panel zrobić nie może: kontrakt ma kontrolę jakości jednego panelu,
więc kontrola zbiorcza jest powtórzeniem tej komendy dla każdego panelu, i okno mówi to wprost,
zamiast sugerować zdolność wsadową, której rdzeń nie ma. Wywołania idą równolegle, bo kontrola
jednego panelu nie zależy od kontroli drugiego, więc szeregowanie ich tylko wydłużałoby czekanie,
a odmowa jednego panelu zostaje przy nim i nie przerywa pozostałym. Przebieg akceptacji tłumaczenie
– korekta – zatwierdzenie pokazuje stany, które kontrakt zna: oczekuje, tłumaczenie w toku, gotowe,
błąd; etapu zatwierdzenia ani autora zmiany w kontrakcie nie ma, więc wskaźnik mówi o stanie
wykonania i nazywa tę różnicę, zamiast malować przebieg, którego rdzeń nie prowadzi.
Bilans sprawdzenia wszystkich paneli jest obowiązkowy: przy wielu wywołaniach część potrafi się
nie udać, a zdanie mówiące wyłącznie o powodzeniu ukryłoby panele, których nie sprawdzono, dlatego
powody odmów idą po nazwie panelu, żeby wiadomo było który. Zastrzeżenia w wykazie nie są wynikiem
ostatniego kliknięcia, tylko odbiciem tego, co rdzeń trzyma przy panelach: kontrola zapisuje je
w panelu, a panel wraca do modułu zdarzeniem zmiany, dzięki czemu wykaz jest prawdziwy także wtedy,
gdy kontrolę uruchomiono z paska narzędzi pojedynczego panelu.

## budowa/klient/src/protokol/warstwa-protokolu.test.ts
Mierzone są trzy rzeczy: koperta i ramka w obie strony, wiązanie odpowiedzi z żądaniem oraz odbiór zdarzeń — wszystkich, jakie zna kontrakt, bo warstwa nie wybiera spośród nich i żadnego nie wyróżnia nazwą wpisaną w kod.

## budowa/klient-poprzedni/src/moduly/research/stan-okna-badania.ts
Rdzeń odpowiada na każde wywołanie, także odmową, a odmowa ma być widoczna tam, gdzie ją wywołano,
nie wyłącznie w konsoli. Forma stanu pustego ma trzy części: ikona, tytuł, opis. Wariantów stylu nie
ma — różnicuje wyłącznie treść, a ikona jest treścią, nie wariantem. Stan początkowy mówi, że okno
jeszcze nie pytało rdzenia; meldunek o pustce rdzenia postawiony przed jego odpowiedzią byłby
zmyśleniem. Nazwa fazy i jej znakowanie pochodzą ze wspólnej biblioteki komponentów: wartość trafia
do atrybutu danych, po którym sięgają arkusze i sprawdziany, więc własny zestaw wartości w module
oznaczałby ten sam stan pod inną nazwą niż w oknach sąsiadów. W tym module zostaje wyłącznie to,
czym Research różni się świadomie: wskaźnik odczytu i chowanie treści na czas ładowania. Stan nie
kasuje treści, tylko ją przesłania — nieudane odświeżenie zostawia to, co już było widoczne,
a powrót do stanu gotowego odsłania treść nietkniętą. Ładowanie niesie wskaźnik obok opisu, nigdy
samodzielnie na pełnym ekranie.

## budowa/klient-poprzedni/src/strona-glowna/strefa-archiwum.ts
Wykaz sesji archiwum idzie z odpytania listy archiwum, a przywrócenie z osobnej komendy; obie wykonuje
wpięcie, ten plik zna wyłącznie ich wynik. Strefa startuje zwinięta i odpytuje rdzeń dopiero przy rozwinięciu,
a potem przy każdym kolejnym — archiwum bywa długie, nie dotyczy pracy bieżącej, a sesja mogła w międzyczasie
wrócić albo dojść. Liczba wszystkich sesji archiwum przychodzi obok strony wyników, więc nagłówek pokazuje
obie liczby, gdy się różnią; sama długość strony nie jest liczbą sesji w archiwum. Odpytywanie, archiwum
puste i odmowa rdzenia mają osobne napisy — wspólny zacierałby różnicę między brakiem danych a brakiem
odpowiedzi.

## budowa/klient/src/protokol/wywolanie.ts
Wywołanie nie jest drugą drogą do rdzenia: każde idzie tym samym wysłaniem kanału, z nazwą komendy wziętą wyłącznie ze stałych kontraktu. Gdy rdzeń nie odpowie w ogóle, bo połączenie padło w trakcie, obietnica pozostaje nierozstrzygnięta: kontrakt nie przewiduje limitu czasu, a rozłączenie klienta nie kończy pracy rdzenia nad poleceniem.

## budowa/klient-poprzedni/src/motyw/motyw.ts
Brak zapisanego wyboru motywu nie ustawia atrybutu koloru na dokumencie:
rozstrzyga wtedy zapytanie o preferencję systemu, a zmiana tej preferencji
działa na żywo. Żaden błąd pamięci trwałej nie zatrzymuje uruchomienia —
wybór degraduje się do preferencji systemu, nigdy do blokady. Oba motywy są
równoprawne, żaden nie jest wartością domyślną produktu.

## budowa/klient-poprzedni/src/powloka/rama-usuniecia.ts
Jedna odpowiedzialność: zbudować modal, który mówi, co zginie. Kolejność czynności, wywołanie
rdzenia i stany należą do warstwy przebiegu potwierdzenia; ten plik nie zna ani komendy, ani kanału.
Ramę niesie biblioteka komponentów: klasy modalu, nagłówka, ciała i stopki pochodzą z arkusza
nakładki, pas stanu z arkusza pustego stanu, a znak zapytania z komponentu dymka objaśnienia.

## budowa/klient-poprzedni/src/moduly/research/wiersz-ustalenia.ts
Wiersz niesie dwie drogi działania: zaznaczenie do raportu oraz wciągnięcie do formularza, czyli
funkcję operatora „edycja", która w kontrakcie jest polem findingId tej samej komendy zapisu.

## budowa/klient/src/sprawdzian.ts
Jego deklaracje typów mieszkają w pakiecie zewnętrznym, którego nowy klient nie zaciąga, więc sprawdzenie typów odmówiłoby każdemu plikowi sprawdzianu. Zbiór pusty jest tu niepowodzeniem, nie wynikiem: sprawdzian, który niczego nie zmierzył, milczałby dokładnie tak samo jak sprawdzian zdany.

## budowa/klient-poprzedni/src/moduly/translate/okno-source-panel.ts
Funkcja operatora z wykazu jest jedna: wprowadzenie albo wklejenie tekstu źródłowego, a zmiana
źródła uruchamia jednoczesną aktualizację wszystkich Translation Panels, więc zapis nie kończy się
na tym oknie — odpowiedź zapisu źródła niesie komplet paneli i wchodzi do stanu modułu. Import
pliku nie ma komendy: panel akcji wymienia „Importuj plik", a obszar translate nie niesie
w kontrakcie komendy przyjmującej plik źródłowy, więc przycisk zostaje klikalny i mówi, czego
brakuje, zamiast być wygaszony bez wyjaśnienia. Drogą wpisania treści bez zapisu jest Format
Studio: wydobyty tekst dokumentu ma trafić tam, gdzie operator go zobaczy i podda zapisowi, bo
zapis źródła uruchamia aktualizację wszystkich paneli i wykonanie go automatycznie byłoby
zleceniem przekładu, którego nikt nie zamówił.
Pole języka źródłowego jest wejściem operatora i miejscem, w które wpisuje się jeszcze
niezapisane rozpoznanie pisma, a osobne zdanie mówi, co rdzeń trzyma u siebie, bo nie wystarczy
przeliczyć go przy ogłoszeniu stanu modułu — rozjazd powstaje także wtedy, gdy operator pisze
w polu i gdy rozpoznanie wpisuje tam wynik, które zdarzeniem stanu nie jest. Odczyty liczą się
z treści pola, a nie za ogłoszeniem stanu modułu, bo tekst przed zapisem jest tym, nad którym
operator pracuje. Przerysowanie nie zdejmuje fazy trwającej ani fazy błędu, bo odświeżenie
przychodzi z każdego ogłoszenia stanu modułu, także w środku zapisu i zaraz po odmowie — gdyby
wtedy stawiało okno na gotowe albo na pustkę, skasowałoby nieprzeczytany komunikat albo zgasiłoby
zapowiedź trwającego wywołania.

## budowa/klient-poprzedni/src/moduly/research/wiersz-wyniku.ts
Przycisk dodania do źródeł jest czynny zawsze, także dla pozycji już przeniesionej: powtórzenie
kończy się odpowiedzią rdzenia, nie odebraniem klikalności. Wiersz nie niesie pola wyboru, bo
zaznaczenie wielokrotne w module dotyczy źródeł i ustaleń jako bytów badania, a pozycja wyniku
badaniem jeszcze nie jest — staje się nim dopiero po skatalogowaniu.

## budowa/klient-poprzedni/src/strona-glowna/pamiec-zwiniecia.ts
Strefy rozwijają się na żądanie, a raz wykonane rozwinięcie ma się utrzymać między wejściami na ekran.
Nastawa siedzi w pamięci przeglądarki, nie w rdzeniu: dotyczy powierzchni na tym urządzeniu i nie ma
swojego bytu w kontrakcie; ten sam wzorzec nosi pamięć motywu i pamięć kolejności okien równoległych.
Pamięć przeglądarki bywa niedostępna w trybie prywatnym albo w osadzeniu w ramce: awaria odczytu albo
zapisu zostaje przy wartości domyślnej i nie jest zgłaszana jako błąd, bo nastawa widoku nie jest powodem,
żeby ekran nie wstał.

## budowa/klient-poprzedni/src/powloka/rozliczenie-usuniecia.ts
Każdy z dwóch wykazów rdzenia idzie osobnym zdaniem, bo sesja bez odpowiednika w historii to nie
sesja skasowana. Odpowiedź udana z pustym wykazem usuniętych znaczy „nic nie zginęło" i tak brzmi
jej zdanie. Sesję nazywamy tytułem karty z pasa; gdy tytułu nie ma, zdanie pokazuje sam identyfikator.
Funkcje są czyste i nie znają DOM.

## budowa/klient-poprzedni/src/okno-komunikacji/okno.ts
Moduł jest właściwością okna, nie wdrożenia: przestawienie modułu zmienia wskaźnik modułu, pasek narzędzi promptu i panel kontekstu, a historii wątku nie dotyka. Panelu akcji w tym złożeniu nie ma, bo pozycje panelu pochodzą z rejestru rdzenia, a złożenie bez kanału nie ma jak ich pobrać — panel bez katalogu byłby atrapą; okno z kanałem składa osobna warstwa rozmowy, w której panel akcji jest pełny. Okno przechowuje warstwę dyktowania i podaje ją elementowi rysującemu mikrofon, ponieważ zna oba końce — kanał i pasek poleceń — a warstwa sama ich nie widzi.

## budowa/klient-poprzedni/src/moduly/studio/czynnosci-cyfryzacji.ts
Wsad wcześniej szedł jednym wywołaniem wydobycia tekstu, bez silnika, bez progu pewności i bez
pojęcia kolejki. Dziś kolejkę prowadzi rdzeń, rozpoznanie ma pełne sterowanie, poprawka słowa
wchodzi przed przyjęciem, a przyjęcie zakłada dokument wraz z pierwszą wersją, nie sam bufor
edytora — poprzednio wynik szedł do bufora i panel musiał tłumaczyć brak pierwszej wersji, czego
wydobycie tekstu nie zakładało. Dokument wchodzi do stanu modułu, żeby okno pracy zobaczyło go
natychmiast. Odmowa jednej pozycji nie zatrzymuje pozostałych: zatrzymanie całej kolejki na
pierwszym nieczytelnym skanie byłoby karą za materiał, a nie obsługą błędu, dlatego bilans na
końcu mówi, ile pozycji przeszło, ile wróciło do ponowienia i ile odpadło.

## budowa/klient-poprzedni/src/motyw/podglad-zetonow.ts
Strona nie dowodzi, że produkt działa: nie ma tu rdzenia, kanału kontraktu,
sesji ani okien — sam katalog wartości żetonów motywu. Poprawny wygląd żetonu
tutaj nie znaczy, że jakikolwiek ekran produktu go używa, a plik nie wchodzi
do pakietu produkcyjnego. Uruchamia się go osobno poleceniem npm run dev pod
adresem lokalnym src/motyw/podglad-zetonow.html.

Plakietka stanu bierze barwę z rodziny żetonów stanu, a ikonę i napis z
katalogu znaczeń stanów. Rodzina neutralna nie ma własnych żetonów — to
plakietka bazowa, a znaczenie niesie w niej wyłącznie znak.

## budowa/klient-poprzedni/src/strona-glowna/sygnal-wyboru.ts
Strona główna nie otwiera środowiska i nie wysyła komendy — skutek wyboru należy do odbiorcy sygnału.
Własny sygnał zamiast zdarzenia niestandardowego na elemencie DOM utrzymuje wybór w pełni typowany aż
do słuchacza: zdarzenie DOM niesie ładunek nietypowany i gubi typ wyboru na granicy ładunku zdarzenia.

## budowa/klient-poprzedni/src/rozmowa/ostatnio-uzyte.ts
Rejestr jest podstawą szeregowania, które trzyma świeżo użyte pozycje wykazu
po ukośniku bliżej wierzchu. Żyje w pamięci okna, przeżywa otwarcie i zamknięcie
wykazu, nie przeżywa odświeżenia strony. Kontrakt takiego pojęcia nie niesie —
`tools.catalog.list` oddaje wykaz, nie historię sięgania po niego — więc
rejestr nie idzie do rdzenia żadną komendą i nie udaje jego stanu.

## budowa/klient-poprzedni/src/rozmowa/podglad-rozmowy.ts
Podgląd pokazuje warstwę rozmowy samą: bez powłoki, bez nawigacji modułów, bez
pasa kart sesji — na gołej scenie i z przełącznikiem motywu. Dzięki temu
usterkę rozmowy widać bez zgadywania, czy winna jest rozmowa, czy warstwa,
która ją osadza. Osobno wykonuje zapewnienie kanału głównego, czyli założenie
wiersza rejestru na świeżej bazie — drogę, którą produkt przechodzi tylko ręką
Operatora w panelu sterowania. Podgląd nie dowodzi, że rozmowa działa
w produkcie: produkt osadza ją w module, w powłoce i w oknie sesji, a tutaj
powłoki w ogóle nie ma.

Uruchomienie: `cd budowa/client && npm run dev`, a następnie adres
`http://localhost:5173/src/rozmowa/podglad-rozmowy.html` (opcjonalnie
z parametrem `?pytanie=…`, żeby zobaczyć pełną turę). Sam plik to wyłącznie
kompozycja: arkusze, odszukanie miejsc w dokumencie, przełącznik motywu,
uruchomienie. Logika mieszka w `polaczenie-podgladu.ts`.

## budowa/klient-poprzedni/src/rozmowa/podsumowanie-tury.ts
Rdzeń dokłada podsumowanie do ostatniego fragmentu strumienia — tego, który
koperta znakuje polem `done`. To ono wybudza koordynatora w pętli koordynatora
i wykonawcy, dlatego jedzie osobnym ładunkiem, a nie tekstem. Kształt odpowiada
`injection.ZakonczenieTury` po stronie rdzenia. Ostatni fragment bywa
domknięciem pustej tury, dlatego ładunek bez ani jednego pola własnego
podsumowania nie jest uznawany za podsumowanie.

## budowa/klient-poprzedni/src/rozmowa/polaczenie-podgladu.ts
Droga uruchomienia podglądu jest ta sama, którą pójdzie powłoka: transport,
kanał kontraktu, zapewnienie kanału głównego w rejestrze, uzgodnienie
(powitanie, sesja, okno), rozmowa okna. Podgląd nie stawia atrapy — rozmawia
z rdzeniem tak samo jak aplikacja. Kanał główny jest zapewniany przed
uzgodnieniem, ponieważ `window.create` wskazuje kanał modelu, a świeża baza
rdzenia nie ma jeszcze ani jednego wiersza rejestru.

## budowa/klient-poprzedni/src/moduly/research/wiersz-zrodla.ts
Pole wyboru niesie zaznaczenie wielokrotne, bo źródła wybiera się grupami do eksportu. Przycisk
„Czytaj" jest akcją podstawową pozycji i otwiera materiał w Reading View — jest czynny przy każdym
źródle, także przy takim, które nie wskazuje dokumentu repozytorium: wtedy odpowiedzią jest zdanie
czytnika o braku drogi do treści, a nie wygaszona kontrolka.

## budowa/klient-poprzedni/src/strona-glowna/strefa-srodowisk.ts
Nazwa strefy jest nazwą ustaloną w warstwie projektowej, nie parafrazą; stoi tak samo w schemacie strony
i w makiecie Centrum dowodzenia, więc nazw sekcji nie tłumaczymy i nie skracamy.

## budowa/klient-poprzedni/src/moduly/translate/okno-translation-memory.ts
Kontrakt nie ma odczytu par, ich edycji, progu dopasowania, zasięgu, wymiany TMX ani operacji
konserwacyjnych — okno robi więc to jedno, co da się zrobić, i przy każdej pozostałej funkcji
mówi, czego brakuje, zamiast stawiać kontrolkę, która nie ma czego wysłać. Konkordancja i
podpowiedź to w tym oknie jedna czynność, bo w kontrakcie są jedną komendą: pole szukania
wypełnia się segmentem wybranym z tekstu źródłowego albo frazą wpisaną ręcznie, a wynikiem jest
wykaz podpowiedzi rdzenia — dwa osobne przyciski nad jedną komendą sugerowałyby dwie różne
zdolności. Wskazanie panelu jest wymagane przez kontrakt, a nie przez okno: pamięć odpowiada
w języku panelu, więc bez panelu nie ma języka, w którym miałaby podpowiadać.
Pozostałe czynności rodziny pamięci — wykaz par, zapis pary, usunięcie, wymiana z plikiem TMX,
utrzymanie, tłumaczenie wstępne, wyrównanie i polityka okna — mają swoje pola w oknie Warsztat
tłumaczenia i nie są tu powtórzone ani nazwane brakiem: jedna czynność w dwóch oknach to dwie
drogi, które rozjadą się przy pierwszej zmianie kontraktu, a napis „brak" przy czynności, która
działa, jest zwykłą nieprawdą.
Wykaz trafień w wyszukiwaniu zostaje na widoku, bo jest wynikiem szukania, a nie odbiciem stanu
modułu, więc zmiana panelu w innym oknie nie ma go kasować; wykaz pusty jest wynikiem, nie pustką
okna — rdzeń odpowiedział i dopasowania nie znalazł, a odmowa czyści wykaz, bo trafienia sprzed
odmowy dotyczyłyby innego zapytania niż to, które właśnie zawiodło.

## budowa/klient-poprzedni/src/okno-komunikacji/opis-okna.ts
Moduł i kanał modelu w opisie początkowym zostają puste, ponieważ pierwsza pozycja wykazu znanych modułów jest kodem środowiska, nie modułu, a pierwsza pozycja wykazu rodzajów kanału jest rodzajem kanału, nie kodem wiersza rejestru — podstawienie którejkolwiek wartości wskazywałoby moduł spoza katalogu i kanał, którego rejestr nie zna, więc pierwsza wypowiedź wracałaby odmową. Widok nie zna rejestru rdzenia w chwili zakładania okna, dlatego moduł i kanał podstawia dopiero pierwszy kanał czynny.

## budowa/klient-poprzedni/src/powloka/sekcje-panelu.ts
Moduł podaje własne sekcje oraz identyfikator panelu; pętlę odczytaj układ, zmień, zapisz, przerysuj
niesie ten plik. Układ adresuje parę okno-panel, więc bez wskazanego okna zmiana zostaje miejscowa
i mówi to wprost. Widok przerysowuje się układem oddanym przez rdzeń w polu odpowiedzi, nie układem
wysłanym — rozjazd obu jest wtedy widoczny od razu. Zapis, który dałby układ tożsamy z bieżącym, nie
idzie.

## budowa/klient/src/wejscie/ekrany/dostep.ts
Dziesięć odsłon leży razem, bo dzielą oprawę — belkę, kolumnę tożsamości i zakładki nad treścią. Zakładki niosą wyłącznie dwie drogi równorzędne: logowanie i rejestrację. Odzyskiwanie dostępu nie jest trzecią drogą — jest wyjściem z logowania, więc w zakładkach zostaje zaznaczone logowanie.
Odsłona konta założonego bez potwierdzenia nie stoi w prototypie. Wymusza ją reguła dziedzinowa: rejestracja bez konta nadawczego kończy się wejściem hasłem, a okno ma wtedy nazwać adres, którego nikt nie potwierdział.
Prototyp stawia w miejscu długości i ważności drogi potwierdzenia wartości przykładowe; okno bierze wartości prawdziwe z pomiaru na rdzeniu.
Cały tor odzyskiwania dostępu stał tak w prototypie od pierwszego kroku.

## budowa/klient-poprzedni/src/moduly/research/wskazanie-lektury.ts
Wydzielone jako osobny byt, bo dotyczy dwóch okien naraz — Sources Manager naciska przycisk
czytania, a Reading View wczytuje wskazany materiał — a zbiorcze zaznaczenie pozycji niesie zbiór
i tu byłoby narzędziem o jeden wymiar za dużym: czytać można jedno źródło. Wskazanie nie warunkuje
klikalności ani jednej kontrolki: Reading View bez wskazanego źródła pokazuje zdanie „nie wskazano
czego czytać", a nie wygaszony przycisk. Ogłoszenie idzie wyłącznie po faktycznej zmianie — odbiorcą
jest przerysowanie okien, a przerysowanie woła ograniczenie wskazania; ogłoszenie bezwarunkowe
zamknęłoby pętlę bez końca. Ten sam wzorzec niesie plik `zaznaczenie-pozycji.ts`.

## budowa/klient-poprzedni/src/okno-komunikacji/panel-akcji.ts
Katalog pusty nie daje panelu pustego: gdy pozycji nie ma, panel pokazuje wyjaśnienie powodu zamiast milczącej pustej przestrzeni.

## budowa/klient-poprzedni/src/strona-glowna/wykaz-srodowisk.ts
Wykaz jest obiektem, a nie funkcją, bo nazwę trzeba odczytać w chwili rysowania wiersza, a nie w chwili
importu modułu — odpowiedź rdzenia przychodzi później niż pierwsza klatka. Wpięcie strony głównej spina
go ze strefą środowisk, więc zasilenie kart bez zasilenia wykazu jest niewykonalne. Do pierwszej odpowiedzi
rdzenia wykaz niesie tę samą stałą, z której rysują się karty — jedna treść zastana, nie dwie. Kod spoza
wykazu wraca dosłownie, bo rdzeń mógł oddać środowisko, o którym klient nie wie, a wiersz sesji ma wtedy
pokazać kod zamiast przemilczeć środowisko; wykaz jest informacyjny, nie jest bramą.

## budowa/klient-poprzedni/src/moduly/studio/czynnosci-edytora.ts
Widok Studio Editora odpowiada za układ, stany i podpięcie zdarzeń, a ten plik za to, co dzieje
się po naciśnięciu: rozmowę z rdzeniem i skutek dla stanu modułu, więc zmiana kształtu żądania
nie wymaga czytania układu okna. Odmowa zostaje w pasie stanu, powodzenie w wierszu odpowiedzi,
bo komunikat błędu ma być trwały w układzie, a potwierdzenie czynności ustępuje następnemu
naciśnięciu. Żądanie otwarcia dokumentu ze wskazaniem ścieżki oddaje dokument bez treści, bo
rdzeń nie ma dostępu do systemu plików operatora; ze wskazaniem pliku repozytorium oddaje
dokument z zapamiętanym odwołaniem, również bez treści, bo doczytanie jej zrobiłoby z adaptera
klienta cudzego modułu — treść dostarcza dopiero pierwszy zapis. Wstawienie wyniku operacji nie
jest przyjęciem propozycji: treść trafia do bufora edytora, a propozycja czeka dalej na decyzję,
bo zrównanie obu czynności odebrałoby możliwość wstawienia fragmentu i odrzucenia reszty.

## budowa/klient/src/wejscie/ekrany/przygotowanie.ts
Kolumna tożsamości niesie tu animację powłok zamiast wykazu zdań: na tym etapie użytkownik nie wybiera już programu, tylko czeka, aż się złoży.
Wykaz niesie wyłącznie etapy, dla których droga wejścia ma komendę, więc postęp dobiega stu procent — pasek stojący w połowie na zawsze czytałoby się jak zawieszenie programu.
Pas działań daje dwie czynności. Ponowienie stoi wyłącznie przy etapie nieudanym: przy przebiegu udanym nie ma czego ponawiać, a kontrolka bez skutku jest gorsza od jej braku. Widoczność nastawia montaż, bo tylko on widzi stan. Czynności pomijającej przywracanie tu nie ma — przywracanie dzieje się w rdzeniu jednym wywołaniem, a kontrakt nie zna komendy, która by je odwołała.

## budowa/klient-poprzedni/src/motyw/znaczenia-stanow.ts
Granica wiedzy pliku: katalog mówi, jak stan ma być pokazany, nie mówi, jaki
stan jest — to rozstrzyga kontrakt współdzielony. Stan „do weryfikacji” nie
ma odpowiednika w kontrakcie i siedzi na tej samej rodzinie barw co stan
„wstrzymany”, dlatego bez ikony i etykiety byłyby na ekranie nierozróżnialne.

Stan „anulowany” kontrakt zna, tyle że pod innym wyliczeniem używanym przez
moduł Assistant, który prowadzi własny wykaz plakietek różniący się od tego
katalogu — dołożenie tu drugiego zapisu bez usunięcia tamtego dałoby trzecią
prawdę zamiast jednej. Stany „do weryfikacji” i „przyjęty” nie mają
odpowiednika, bo pozycja kolejki niesie wyłącznie stan kolejki, a postęp
osobny stan postępu.

Klasa wariantu kropki: praca trwająca bierze tętno, nie barwę rodziny, bo
tętno jest jedynym ruchem ciągłym interfejsu i jest zastrzeżone dla pracy w
tle. Przy preferencji ograniczonego ruchu tętno zamiera, a jego znaczenie
przejmuje pierścień statyczny, również zapisany w arkuszu stylów.

## budowa/klient-poprzedni/src/strona-glowna/wykonanie-czynnosci.ts
Składa komplet czynności i przeprowadza te, które idą do rdzenia bez pytania operatora o cokolwiek.
Czynności pytające o nazwę stoją osobno; dzieli je nie temat, tylko kształt — tam każda ma etap zbierania
danych, tu żadna go nie ma. Komendy zbiorowe oddają wykaz sesji faktycznie przeniesionych, a ten bywa
krótszy od żądania, stąd meldunek składany jest z wyniku. Po każdej udanej zmianie woła się odświeżenie
podane przez wpięcie: rdzeń rozsyła zdarzenie zmiany sesji, ale archiwum jest odpytywane osobno i tego
zdarzenia nie widzi.

## budowa/klient-poprzedni/src/moduly/studio/czynnosci-narzedzi.ts
Zakres operacji kontekstowej bierze się ze stanu modułu, nie z kontrolki wyboru: zakres skuteczny
liczy stan studia i biorą go stamtąd wszyscy trzej odbiorcy — wskaźnik panelu, pasek zaznaczenia
edytora i to żądanie — bo druga kopia nastawy trzymana osobno w liście wyboru panelu rozjeżdżała
się ze stanem. Wybór zakresu zaznaczenia bez samego zaznaczenia nie jest blokowany: schodzi na
cały dokument, a panel pisze o tym we wskaźniku zakresu, natomiast żądanie z zakresem zaznaczenia
i bez jego granic dostałoby odmowę walidacji. Operacja kontekstowa wraca odpowiedzią pomyślną
także wtedy, gdy kanał modelu oddał zamiast wyniku swój własny komunikat — rdzeń tego nie
odróżnia, bo dostał zwykły tekst — dlatego zdanie o wyniku mówi wprost, ile treści przyszło
i że przyjęcie podmieni nią dokument, zamiast zgadywać po samej treści.

## budowa/klient-poprzedni/src/okna-pomocnicze/budzet-rozmowy.ts
Budżet stoi przy oknach pomocniczych, bo pas pomocniczych podaje obie połowy
naraz: ile okien rozmowy moduł prowadzi i gdzie leży jego ciężar. Liczbę
oddaje funkcja, nie samo pole granicy: pole niesie granicę, funkcja — prawo
do otwarcia, więc dla modułu bez rozmowy pole daje jeden, a funkcja zero, co
zapobiega obiecywaniu Operatorowi czatu w module, który go nie prowadzi.
Liczba początkowa pochodzi z tej samej stałej minimalnej, z której korzysta
silnik okien równoległych, żeby obie wartości nie rozjechały się przy jego
zmianie. Granicy zapisanej w profilu nie da się w czasie działania odróżnić
od stałej oznaczającej brak granicy, więc zdanie o budżecie podaje samą
liczbę i nic nie orzeka o jej pochodzeniu.

## budowa/klient/src/wejscie/ekrany/uruchomienie.ts
Odsłony różni stan etapów i to, co stoi pod wykazem. Łączenie i powrót z tokenem nie mają czynności głównej: przechodzą dalej same, gdy rdzeń odpowie. Błąd ją ma, bo tam jest co rozstrzygnąć.

## budowa/klient-poprzedni/src/moduly/research/wybor-nastawy.ts
Ster nastawy niesie wartość bieżącą na uchwycie, na przykład „PDF", a nie ogólną nazwę pola typu
„Format wyjściowy" — bo to ona jest odpowiedzią na pytanie, co jest teraz ustawione; natywna lista
wyboru tego nie robi, pokazuje wartość dopiero po rozwinięciu. Plik nie jest drugim mechanizmem
rozwijania: rozwijanie, znacznik wyboru, opisy pozycji, wędrówka strzałkami, pole szukania po progu
i zdanie o pustym wykazie należą do biblioteki, a ten plik podaje mechanizmowi dane. Ta sama obsada
stoi też w plikach `moduly/apps/wybor-z-menu.ts` (wraz z wymianą pozycji w locie) i
`moduly/browser/ster-wyboru.ts`; miejscem docelowym jest katalog komponentów wspólnych, a do czasu
przeniesienia moduł nie sięga po cudzą obsadę przez granicę modułu, bo import między modułami wiąże
je mocniej niż powtórzenie kilkudziesięciu wierszy. Wykaz jest stały — to jedyna różnica wobec
obsady modułu aplikacji: trzy nastawy Research pochodzą z wyliczeń kontraktu, a nie z odpowiedzi
rdzenia, więc wymiany pozycji w locie tu nie ma. Uchwyt jest klikalny zawsze, także zanim Operator
cokolwiek wybrał.

Podpis wiersza formularza nie jest znacznikiem etykiety, bo uchwyt menu jest przyciskiem,
a przycisk nie jest elementem etykietowalnym: znacznik etykiety owinięty wokół niego nie
przeniosłby ani kliknięcia, ani ogniska, więc udawałby wiązanie, którego nie ma. Nazwę nastawy
niesie etykieta dostępności uchwytu, stawiana przez mechanizm z pola nastawy — podpis jest tu
dla oka, nie dla czytnika ekranu.

## budowa/klient-poprzedni/src/protokol/korelacja.ts
Kontrakt zapowiada dokładnie jedną odpowiedź na komendę, a odpowiedź powtarza identyfikator żądania.
Rejestr nie wprowadza limitu czasu ani limitu żądań oczekujących: zerwanie połączenia nie kończy
pracy rdzenia nad poleceniem.

## budowa/klient-poprzedni/src/protokol/nawigacja-platformy.ts
Plik niesie wyłącznie kształt zależności, nie drogę do rdzenia. Wywołania stoją w plikach
jednokomendowych, a obiekt tego kształtu składa korzeń montażu klienta — jedyne miejsce, które ma
kanał. Powłoka dostaje go gotowego i nie zna nazwy ani jednej komendy. Metody są dwie, bo tyle woła
powłoka: `home.enter` woła wprost widok strony głównej, `workspace.enter` przestrzeń modułu, a
`environment.list` czyta strona główna kanałem.

## budowa/klient-poprzedni/src/protokol/ognisko-sesji.ts
Ognisko jest właściwością klienta, nie konta: to samo konto otwarte na dwóch urządzeniach ma dwa
ogniska, dlatego treść żądania niesie identyfikator klienta. Rozgłoszenie zmiany należy do rdzenia —
po wykonaniu komendy rozsyła on zdarzenie zmiany ogniska, które przechwytuje obserwator ogniska
w warstwie łączności.

## budowa/klient-poprzedni/src/protokol/powiazanie-sesji.ts
Klient albo wiąże się z sesją trwającą na rdzeniu, albo wchodzi na stronę główną. Rozłączenie klienta
nie kończy sesji ani procesów, więc powiązanie zwykle zastaje sesję żywą: pole odtworzenia odróżnia
odtworzenie stanu od założenia go od nowa, a wykaz okien niesie okna, których strumienie od tej chwili
idą na bieżące połączenie. Puste wskazanie okien w żądaniu wiąże wszystkie okna sesji — zawężenie jest
decyzją widoku, nie protokołu.

## budowa/klient-poprzedni/src/protokol/sesja.ts
Sesja pozostaje wspólna dla plików, pamięci, projektu i agentów; okno komunikacji jest wobec niej
bytem podrzędnym. Identyfikator nadaje rdzeń w odpowiedzi na założenie sesji — do tej chwili sesja
jest pusta, a kontrakt taką kopertę dopuszcza.

## budowa/klient-poprzedni/src/protokol/wejscie-do-przestrzeni.ts
Żądanie idzie z identyfikatorem żywego okna rozmowy: okno wskazane przestawia moduł i zachowuje
historię, a dopiero jego brak zakłada okno nowe. Okno zakładane po stronie rdzenia nie dostaje
kanału modelu, więc pierwsze wysłanie wiadomości skończyłoby się odmową. Odpowiedź niesie sesję,
moduł, okno i kody okien operacyjnych modułu — sprawdzian kształtu pyta o wszystkie cztery, bo widok
przestawia planszę na ich podstawie i pusty komplet dałby planszę bez treści.

## budowa/klient-poprzedni/src/protokol/wejscie-do-srodowiska.ts
Jedno wywołanie oddaje komplet bocznej nawigacji: środowisko, jego moduły z katalogiem okien
operacyjnych oraz karty sesji otwarte w tym środowisku. Dlatego wykaz nawigacji nie składa się z dwóch
odczytów — wykaz modułów jest tu wyłącznie drugim podejściem, gdy wejście oddało wykaz pusty mimo
nawigacji modułowej. Warstwa nie buduje ani jednego elementu widoku: oddaje treść odpowiedzi
w kształcie kontraktu i zostawia widokowi rozstrzygnięcie, co z nią zrobić. Nazwa komendy pochodzi
wyłącznie ze stałej kontraktu.

## budowa/klient-poprzedni/src/protokol/wejscie-strony-glownej.ts
Pierwsze ogniwo łańcucha nawigacji: wejście na stronę główną, wykaz środowisk, wejście do środowiska,
wykaz modułów, wejście do przestrzeni. Wejście mówi rdzeniowi, który klient stanął na stronie głównej;
rdzeń wiąże ognisko z klientem, więc bez tego wywołania nie ma komu oddać identyfikatora sesji
ogniskowanej. Ponad wykazem środowisk odpowiedź niesie środowiska wraz z kodami modułów, sesje czynne
konta, sesję ostatnio ogniskowaną oraz stan sesji trwających w tle. Wykaz środowisk bez dołączenia
modułów nie daje żadnej z tych rzeczy.

## budowa/klient-poprzedni/src/protokol/wykaz-modulow.ts
Komenda pełni dwie role: z identyfikatorem środowiska — drugie podejście bocznej nawigacji, gdy
wejście do środowiska oddało wykaz pusty mimo nawigacji modułowej; bez pola — komplet modułów
platformy, niezależny od macierzy środowisko-moduł. Tą drogą otwiera się moduł, którego macierz nie
pokazuje w żadnym środowisku. Odpowiedź niesie licznik obok wykazu, więc sprawdzian kształtu pyta
o oba: wykaz bez licznika znaczy, że rdzeń odpowiedział czymś innym niż wykazem modułów.

## budowa/klient-poprzedni/src/protokol/wynik-czastkowy.ts
Komenda kontraktu oddaje kopertę z polem wykazu, a widok potrzebuje samej zawartości. Pomocnik wycina
pole, przepuszczając odmowę nietkniętą: okna mają obowiązkowy stan błędu i muszą odróżnić brak
zawartości od nieudanego zapytania. Wycięcie pola operatorem pustej wartości w miejscu wywołania
zgubiłoby powód odmowy, a odmowa uprawnienia i brak zasobu są przy sięganiu po cudze okna komunikacji
zjawiskiem zwykłym, nie awarią. Odpowiedź udana, lecz bez treści, jest tu niepowodzeniem bez pola
błędu: rdzeń nie podał powodu, więc pomocnik żadnego nie dopisuje. Pole jest wtedy nieobecne
w obiekcie, nie ustawione na wartość pustą — stąd rozwinięcie warunkowe.

## budowa/klient-poprzedni/src/protokol/wywolanie.ts
Kanał rozdaje wynik przez wywołanie zwrotne, ponieważ tak wygodniej obsłużyć strumień i ruch ciągły.
Widok pyta jednak punktowo — wejdź do środowiska i pokaż, co wróciło — i dla niego obietnica jest
formą naturalną. To nie jest druga droga do rdzenia: każde wywołanie idzie tym samym wysłaniem przez
kanał, z nazwą komendy wziętą wyłącznie ze stałych kontraktu. Obietnica nie jest odrzucana nigdy.
Niepowodzenie wraca jako wynik z polem błędu, więc wywołujący nie musi zakładać obsługi wyjątku,
a brak odpowiedzi nie wywraca widoku. Gdy rdzeń nie odpowie w ogóle, bo połączenie padło w trakcie,
obietnica po prostu pozostaje nierozstrzygnięta: kontrakt nie przewiduje limitu czasu, a rozłączenie
klienta nie kończy pracy rdzenia nad poleceniem.

## budowa/klient-poprzedni/src/moduly/translate/okno-translation-panels.ts
Formularz niesie ster kanału, bo pole kanału żądania dodania celu wskazuje model wykonujący
przekład, a decyzja o nim zapada właśnie tutaj; nota sufitu paneli mówi obok, ile paneli okno zna
(każdy kosztuje jedno wywołanie modelu), nie odbierając dodania kolejnego. Okno nie prowadzi
własnego wykazu paneli — prawdą jest wykaz w stanie modułu, bo napełniają go trzy drogi naraz:
odpowiedź dodania celu, odpowiedź zapisu źródła (komplet paneli po zmianie źródła) i zdarzenie
zmiany tłumaczenia; druga kopia rozjechałaby się przy pierwszej z nich.
Opracowanie mówi o segmencie wymagającym rewizji, ale kontrakt stanu segmentu nie zna — stan
niesie cały panel — więc przejście prowadzi między panelami, których stan nie jest gotowy, a
zdanie o tej różnicy stoi w wykazie skrótów. Panele odrysowują się także wtedy, gdy okno stoi
w ładowaniu, bo panel przeliczony przez rdzeń przychodzi zdarzeniem zmiany tłumaczenia w środku
innego wywołania i ma być widoczny od razu — komunikatu skasować natomiast nie wolno, bo zapowiedź
trwającego dodania języka i powód odmowy zdejmuje czynność, która je postawiła. Wykaz pusty przy
przejściu do uwagi nie jest ciszą — wiersz odpowiedzi mówi, że wszystkie panele są gotowe.
Komenda dodania języka docelowego idzie modelem: rdzeń przekłada tekst źródłowy okna na język
panelu czynnym kanałem modelu, a bez takiego kanału odmawia kodem niedostępności kanału; brzmienie
tej odmowy nie mówi operatorowi, czego brakuje, więc zdanie dokłada funkcja odmowy modelu. Sprawdzenie
rozbieżności odpowiedzi milczy z konstrukcji, a nie przez zaufanie do rdzenia. Kolejność instancji
paneli bierzemy z rdzenia, nie z chwili dodania, inaczej dwa klienty tego samego konta pokazałyby
panele w innym porządku.

## budowa/klient/src/wejscie/ikony.ts
Ten sam sygnet stoi w belce okna i w kolumnie tożsamości, w dwóch różnych wielkościach. Kropka sygnału w godle niesie klasę, nie wpisaną barwę: barwę rozstrzyga arkusz stylu, osobno dla belki i osobno dla kolumny tożsamości. Łańcuchy w tym pliku są rysunkami, nie tekstem — każdy zaczyna się od znacznika rysunku.

## budowa/klient-poprzedni/src/moduly/studio/czynnosci-podgladu.ts
Eksport nie ma osobnej komendy w obszarze studio i mieć jej nie musi: zamiana formatu dokumentu
jest czynnością obszaru dokumentów, ma tam uchwyt w rdzeniu i słownik ośmiu formatów, a przedmiotem
zamiany jest treść zaakceptowana modułu, więc żądanie niesie ją wprost, nie ścieżkę pliku. Wynik
zamiany zostaje zasobem magazynu rdzenia, a komendy wydającej jego bajty wprost do przeglądarki
kontrakt nie niesie, więc okno mówi to wprost przy potwierdzeniu, zamiast pozwalać czytać
„wyeksportowano" jako „pobrano". Przekazanie do modułu Library wykonuje się w całości: rdzeń
zakłada okno modułu docelowego i oddaje potwierdzenie przeniesienia. Odmowa zostaje w pasie
stanu, powodzenie w wierszu odpowiedzi, a wskaźnik odczytu zapala się wyłącznie na czas
rzeczywistego wywołania, żeby warunek sprawdzany przed wysłaniem żądania nie udawał wywołania.

## budowa/klient-poprzedni/src/strona-glowna/strefa-komponentow.ts
Nazwa strefy pochodzi z warstwy projektowej, bez parafrazy. Siatka niesie dwa rodzaje kafli: pierwsze
pochodzą z modułów, które rdzeń oznaczył jako nastawiane na stronie głównej, z czwórką zastaną do pierwszej
odpowiedzi z kontraktu; za nimi stoją kafle personalizowane, po jednym na komponent zbudowany i nazwany
przez operatora. Żadna z dwóch liczb nie jest z góry znana, więc siatka przyjmuje oba wykazy osobno
i przerysowuje się, gdy rdzeń odpowie. Waga wizualna strefy jest niższa niż strefy pierwszej: mniejsza
powierzchnia, mniejszy promień, mniejsza ikona, etykieta krojem bazowym półgrubym, bez wstęgi i bez cienia
sygnału w spoczynku. Strefa bez ani jednego wejścia byłaby regresem widocznym na ekranie, a pusta odpowiedź
bywa też odpowiedzią bazy bez wykonanych migracji; operator, który nic jeszcze nie zbudował, widzi same
kafle modułów. Zwinięcie wykonane ręcznie zostaje zapamiętane pod kluczem strefy, żeby zwinięcie z urzędu
nie ukrywało jej przed operatorem. Druga siatka z własnym podpisem wydłużałaby stronę o dwa pasy i spychała
sesje w tle oraz archiwum poniżej pierwszego ekranu; rodzaj kafla rozpoznaje się po wezwaniu wobec nazwy
własnej. Po odpowiedzi rdzenia na wykaz modułów rozstrzyga rdzeń. Utworzenie komponentu jest czynnością
wykonywaną po obejrzeniu tego, co już stoi.

## budowa/klient-poprzedni/src/moduly/research/wybor-zrodel.ts
Tak Findings Panel wiąże ustalenie ze źródłem, bez potrzeby osobnej komendy zapisu wiązania. Pusty
katalog nie jest błędem ani blokadą: zdanie zastępcze mówi, skąd wziąć źródła, a formularz
ustalenia zostaje w pełni czynny.

## budowa/klient-poprzedni/src/rozmowa/pole-wysylki.ts
Pasek zlecenia pod polem wypowiedzi jest opcjonalny, bo jego stery żądają
kompletu sterowania okna (migawka, subskrypcja, `window.update`, `config.set`),
a ten powstaje dopiero w powłoce — stanowisko podglądu rozmowy go nie ma i ma
dalej działać. Pole przepuszcza element nietknięty: nie zna ani nastaw, ani
kontraktu. Pasek stoi pod polem, w tym samym wierszu co przycisk zatrzymania,
bo Operator ustawia na nim otoczenie zadania, zanim treść pójdzie do modelu.

Wykaz po ukośniku nad polem wypowiedzi też jest opcjonalny z tego samego
powodu: pole nie buduje go i nie wie, skąd biorą się pozycje — zna wyłącznie
cztery ruchy klawiatury, które nim sterują. Wykaz potrzebuje kanału i wątku
rozmowy, a tych pole nie ma; podaje go widok rozmowy.

Odbitka cudzej wypowiedzi pojawia się w polu od razu w całości, na krótko,
i sama znika: wypowiedź jest już wysłana — rdzeń rozgłosił ją zdarzeniem po
przyjęciu komendy — więc literowanie jej znak po znaku pokazywałoby zdanie
niepełne. Odbitka nie zabiera pracy Operatorowi: pole z ogniskiem albo choćby
jednym znakiem treści nie jest ruszane, odbitka nie ma ogniska, nie da się jej
wysłać Enterem i nie wchodzi do treści.

Czas, przez jaki odbitka stoi w polu, ustala kompromis: Operator ma zdążyć
zobaczyć, że tekst się pojawił, i przeczytać początek zdania, zanim zniknie.
Krócej byłoby mrugnięciem nie do złapania wzrokiem, dłużej — zawadą stojącą
w miejscu pracy. Odbitka niczego nie opóźnia, bo wypowiedź poszła do rdzenia,
zanim to połączenie w ogóle się o niej dowiedziało.

Pole wypowiedzi Operatora oraz przyciski wysyłki i zatrzymania nie tracą
klikalności w żadnym stanie. Przycisk zatrzymania jest czynny także wtedy, gdy
nic nie biegnie — rdzeń odpowiada wtedy `stopped` równym fałsz, a nie błędem.
Niegotowość opisuje wskaźnik obok przycisków, nie odebranie możliwości
działania.

Wykaz po ukośniku wyzwala zdarzenie `input`, a nie `keydown`: znak jest już
w polu, więc warunek „ukośnik stoi na początku treści" da się sprawdzić na
treści, a nie zgadywać z klawisza (martwe klawisze, wklejenie, IME). Ukośnik
w środku zdania nie otwiera niczego — ścieżka w rodzaju `/home/ubuntu` wpisana
w wypowiedź jest tekstem, nie komendą.

## budowa/klient-poprzedni/src/okna-pomocnicze/nastawa-retencji.ts
Rdzeń zasadę retencji egzekwuje, nie tylko zapisuje: przemiata pozycje przy
swoim starcie, przy każdym odczycie historii oraz zaraz po zapisie zasady,
więc po zapisie pozycje przestają istnieć naprawdę, a nie tylko formalnie.
Zapowiedź skutku nie jest bramką — nie ma pytania o potwierdzenie, przycisk
zapisu jest czynny zawsze, a zapowiedź daje wiedzę, nie zgodę. Oba progi
puste znaczą brak ograniczenia i tak zdejmuje się zasadę z zakresu, bo
kontrakt nie ma osobnej komendy kasującej — dwie drogi do jednego skutku
byłyby dwiema prawdami o retencji.

Rdzeń przyjmuje zakresy okna, sesji i globalny, przy czym dwa pierwsze
wymagają wskazania bytu. Okno panel zna od gospodarza, sesji nie zna z
własnych opcji — sesja pada wyłącznie we wczytanych pozycjach historii,
więc zakres sesji wchodzi do wyboru dopiero wtedy, gdy panel tę sesję
zobaczył; pozycja obiecująca zakres, którego nie ma czym wypełnić, byłaby
atrapą.

Panel widzi tylko stronę wykazu, nie całą historię okna, a rdzeń liczy próg
liczby pozycji na całości, dlatego zapowiedź jest dolnym oszacowaniem —
pozycji poza zasadą będzie co najmniej tyle, ile panel widzi. Kolejność
pozycji jest kolejnością odczytu historii od najnowszej, więc numer pozycji
w tablicy jest jej numerem od najnowszej i to on wchodzi w próg.

## budowa/klient-poprzedni/src/punkty-izolacji/efektywna-punkt-widzenia.ts
Podgląd polityki bez podanego punktu widzenia schodzi w rdzeniu po kolejce zasięg, okno, sesja,
poziom globalny aż do poziomu globalnego, więc oddałby politykę całej platformy pod nazwą
efektywna. Panel podaje punkt odniesienia jawnie: identyfikator sesji z kanału, warstwę z pasa
narzędzi okna, zasięg i identyfikator zasięgu z selektora zasięgu w panelu lewym, a identyfikator
okna z pola wypełnianego przez Operatora — okno otwiera się z listwy Ustawień i nie jest związane
z żadnym oknem komunikacji, więc kanał niesie wyłącznie sesję. Poziomu ten panel nie wybiera po raz
drugi. Podgląd ma pokazywać wynik dziedziczenia aż do poziomu wskazanego w selektorze zasięgu, jak
opisuje rozdział 6.4 Modelu konfiguracji; własna lista poziomów kazałaby czytać politykę innego
poziomu niż ten, na którym Operator właśnie przestawia macierz.

## budowa/klient-poprzedni/src/punkty-izolacji/stan-zasiegu.ts
Stan jest wspólny dla całego okna z tego samego powodu, dla którego wspólna jest warstwa: poziom
rozstrzyga, który zapis czyta i pisze macierz izolacji, dokąd trafia przypisanie profilu i czego
dotyczy podgląd polityki efektywnej. Osobny selektor w każdym z tych miejsc pokazywałby obok siebie
wartości z trzech różnych poziomów pod jedną nazwą izolacja. Poziom globalny jest warstwą bazową
i jedynym, który nie potrzebuje bytu — każdy węższy wskazuje byt: kod środowiska, identyfikator
projektu, sesji, nazwę roli. Byt pusty przy poziomie węższym nie jest błędem klienta: żądanie idzie
bez pola identyfikatora zasięgu, a odmowę — jeżeli rdzeń bytu wymaga — nazywa on sam. Plik nie woła
rdzenia i nie buduje ani jednego elementu widoku. Selektor stoi w panelu zasięgu, czytelnicy —
w obszarach.

## budowa/klient-poprzedni/src/punkty-izolacji/sterowanie-warstwa.ts
Kontrolka stoi w ramie okna, a nie w obszarze, bo warstwa rozstrzyga, który zapis czyta i pisze
każdy obszar tego okna: kontekst, zakres techniczny, przypisanie profilu i podgląd polityki
efektywnej. Schowana w jednej zakładce byłaby ustawieniem czterech pozostałych, którego z nich nie
widać — dlatego stoi nad paskiem zakładek, tak samo widoczna z każdego obszaru. Warstwa czynna okna
zmienia się dopiero po odpowiedzi rdzenia na komendę ustawienia warstwy — nigdy przed nią i nigdy na
samą wartość wybraną w liście. Gdy rdzeń odmawia, lista wraca do warstwy poprzedniej: pokazywanie
wyboru, którego rdzeń nie przyjął, mówiłoby nieprawdę o stanie maszyny. Lista nigdy nie dostaje
znacznika zablokowania i nie pyta o potwierdzenie — warstwa nie jest kłódką na Operatorze, tylko
wskazaniem, na którym podkładzie pracuje. Odmowa jest meldowana zdaniem trzyczęściowym: co się nie
udało, dlaczego, treścią wprost z rdzenia, i czym Operator to zmieni.

## budowa/klient-poprzedni/src/moduly/studio/czynnosci-pracy.ts
Widok okna pracy z dokumentem odpowiada za układ i stany, a ten plik za skutek naciśnięcia,
tak samo jak w czynnościach Studio Editora, tyle że tutaj dochodzą czynności, których
poprzednie okna nie miały: operacja z poleceniem własnym i nastawami suwaków, decyzja
o zmianach śledzonych, śledzenie, komentarze, szablony i profile wydania — odmowa zostaje
w pasie stanu, powodzenie w wierszu odpowiedzi, tak samo jak w pozostałych czynnościach
modułu. Zakres operacji kontekstowej bierze się ze stanu modułu tą samą drogą, którą liczą
go czynności narzędzi, bo nastawa zakresu jest jedna na moduł; polecenie i nastawy suwaków
jadą polem parametrów, które rdzeń dokłada do treści polecenia dla modelu, a wynik rdzeń
wpisuje do dokumentu jako zmianę śledzoną autorstwa modelu i rozgłasza zmianę dokumentu,
więc okno nie podmienia treści samo. Decyzja o propozycji bez odwołania w rdzeniu zapada
w tym oknie tak samo jak przed scaleniem okien Studio Editora, bo odmowa w tym miejscu
odebrałaby operatorowi decyzję, którą wolno mu podjąć.

## budowa/klient-poprzedni/src/strona-glowna/strefa-modulow.ts
Strefa nie pyta rdzenia i nie otwiera modułu — wykaz podaje jej warstwa wpięcia modułów, a skutek wyboru
należy do warstwy, która stronę zamontowała. Wykaz pusty chowa całą strefę, zamiast zostawiać nagłówek nad
pustym prostokątem: gdy rdzeń przypnie moduły do środowisk, kafle znikną razem ze swoim powodem. Waga
wizualna jest jak w strefie drugiej — ta sama karta, ten sam krój bazowy półgruby, ta sama siatka. Kafel
prowadzi do pracy w module, a nie do zbudowania rzeczy, więc tak jak kafel komponentu nie nosi akcentu
zarezerwowanego dla kart środowisk. Zapowiedź z liczbą modułów stoi zawsze, więc zwinięcie niczego nie
ukrywa przed operatorem, a zapowiedź nad pustką mówiłaby o wykazie, którego nikt jeszcze nie odczytał. Zero
okien jest stanem możliwym i mówi się je wprost, zamiast chować wiersz.

## budowa/klient-poprzedni/src/moduly/translate/panel-jezyka.ts
Okno jest instancją wielokrotną, liczoną jak liczba języków, więc panel jest tu tym, czym wiersz
w liście: jedną rzeczą na jeden plik, a okno zbiorcze nimi zarządza i nie wie, jak są zbudowane.
Każda instancja ma własny pas stanu i własny wiersz odpowiedzi, więc odmowa kontroli jakości dla
jednego języka nie gasi panelu innego języka, dlatego stan okna powstaje osobno dla każdego panelu,
a nie raz na okno zbiorcze. Drogą prowadzenia ogniska jest przejście skrótem między panelami
wymagającymi uwagi: panel wskazany bez ogniska wymagałby jeszcze jednego kliknięcia, żeby zacząć
poprawiać, a skrót ma prowadzić do pracy, nie do widoku. Instancja znika, gdy rdzeń przestaje
oddawać jej panel, a jej ster kanału bywa wtedy rozwinięty i trzyma nasłuch na dokumencie, którego
samo usunięcie elementu nie zdejmuje.
Zapis korekty odpowiedzią niesie cały panel po korekcie razem z jego treścią, więc zgodność da się
sprawdzić bez dodatkowego wywołania. Ton widać w dwóch miejscach i oba biorą z jednego źródła:
nagłówek oraz pole tonu panelu w pasku narzędzi dostają ton z tej samej odpowiedzi rdzenia,
w tym samym przebiegu. Zdanie o pustce panelu mówi, czym ten panel jest i czym się go zapełnia,
zamiast nazywać brak; forma stanu pustego jest w module jedna, tytuł nad opisem, różnicuje ją treść.

## budowa/klient-poprzedni/src/moduly/research/wynik-odkrycia.ts
Odpowiedź wyszukiwania wiedzy oddaje fragmenty treści, a odpowiedź wyszukiwania repozytorium
zasoby repozytorium; panel pokazuje jedną listę wyników i wstawia z niej źródła jedną komendą, więc
obie odpowiedzi sprowadza do jednego bytu tutaj, a nie w widoku — widok nie ma rozstrzygać, z której
komendy pochodzi wiersz, który rysuje. Przeniesienie do badania jest komendą dodania źródła; kształt
zlecenia powstaje z pozycji, więc reguła „fragment wiedzy wchodzi jako notatka, zasób repozytorium
jako dokument" stoi w jednym miejscu.

Adresu pozycje wyszukiwania nie niosą — ani fragment wiedzy, ani zasób repozytorium nie mają pola
z adresem sieciowym — więc pole adresu zostaje puste zamiast być dopowiedziane ze ścieżki
repozytorium, która adresem nie jest.

Identyfikator trafienia jest identyfikatorem tego, z czego fragment pochodzi, ale pochodzić może
z pliku biblioteki, z wiadomości albo z pliku przestrzeni roboczej — zakres wyszukiwania to
rozstrzyga. Wiązanie z dokumentem repozytorium zakładamy wyłącznie dla zakresu biblioteki; poza nim
identyfikator wskazuje byt, którego pole dokumentu repozytorium w komendzie dodania źródła nie
przyjmie.

## budowa/klient-poprzedni/src/punkty-izolacji/okno-punktow-izolacji.test.ts
Pilnowane są cztery rzeczy, które wcześniej były w tym oknie zaszyte na sztywno albo rozdzielone na
zakładki: trzy panele w kolejności selektor, macierz, profil, jak opisuje rozdział 6.1 Modelu
konfiguracji, a nie pasek zakładek; objaśnienie kontekstowe przy każdym przełączniku macierzy, jak
opisuje rozdział 3.3, z treścią, nie samym znakiem; zasięg i warstwa niesione w żądaniach izolacji
zamiast wpisanych w kod; wskazanie poziomu w selektorze przestawiające odczyt macierzy — dowód, że
lewy panel steruje, a nie tylko wygląda na selektor. Rdzeń jest atrapą: sprawdzian pyta o zachowanie
okna, nie o zachowanie serwera. Odpowiedź nieznana atrapie wraca odmową — tak jak wraca z rdzenia,
który komendy nie zna.

## budowa/klient-poprzedni/src/moduly/studio/czynnosci-redakcji.test.ts
Sprawdzian składania żądań redakcji dokumentu pilnuje trzech miejsc, w których łatwo o cichą
pomyłkę: pole puste nie trafia do żądania, bo pusty napis jest dla rdzenia wskazaniem, a nie
jego brakiem; pole trójstanowe „domyślnie" nie wysyła wartości fałszywej, bo w paczce redakcyjnej
taka wartość wyłącza człon, o którego wyłączenie nikt nie prosił; czynność wymagająca dokumentu
odmawia zdaniem, a nie wysłaniem żądania bez identyfikatora dokumentu i czekaniem na odmowę
walidacji rdzenia.

## budowa/klient-poprzedni/src/moduly/studio/czynnosci-warsztatu.test.ts
Sprawdzian składania żądań warsztatu dokumentu pilnuje miejsca, w którym formularz zamienia się
w treść kontraktu: pole puste nie trafia do żądania, a nie trafia jako pusty napis, bo pusty
napis w polu nieobowiązkowym jest dla rdzenia wskazaniem, a nie jego brakiem, i zmienia znaczenie
czynności — w szyfrowaniu decyduje to, czy hasło jest nakładane, czy zdejmowane.

## budowa/klient-poprzedni/src/moduly/research/zaznaczenie-pozycji.ts
Wydzielone jako osobny byt, bo dotyczy dwóch okien naraz — Sources Manager, przy wyborze wielu
źródeł do eksportu, i Findings Panel, przy wyborze ustaleń do raportu. Zaznaczenie nie warunkuje
klikalności: przy pustym zbiorze akcja zostaje czynna i odpowiada zdaniem opisowym, zamiast być
wygaszona. Zbiór ogłasza swoją zmianę i dlatego konstruktor żąda wywołania zwrotnego — zbiór jest
wspólny obu oknom, więc bez ogłoszenia drugie okno przerysowałoby się dopiero przy najbliższej
zmianie treści badania: ta sama nastawa pokazywałaby w dwóch oknach dwie różne wartości, a wiązanie
źródła z ustaleniem brałoby to, czego w panelu ustaleń nie widać. Wywołanie zwrotne jest więc
obowiązkowe, nie domyślne. Ogłoszenie idzie wyłącznie po faktycznej zmianie: odbiorcą jest
przerysowanie okien, a przerysowanie woła ograniczenie zaznaczenia — ogłoszenie bezwarunkowe
zamknęłoby pętlę bez końca.

## budowa/klient-poprzedni/src/strona-glowna/strefa-sesji.ts
Strefa nie zna kanału ani nazwy żadnej komendy — dane przynosi źródło sesji, a wiąże je warstwa wpięcia
sesji. Wykaz powstaje wyłącznie z wpisów rdzenia, w trzech postaciach bez danych miejscowych: oczekiwanie,
gdy rdzeń jeszcze nie odpowiedział, pusto, gdy rdzeń odpowiedział i sesji w tle nie ma, oraz błąd, gdy rdzeń
odmówił — stan pusty nie udaje wtedy braku sesji. Zapowiedź niesie liczbę sesji, więc zwinięcie domyślne
nie ukrywa faktu, że coś trwa. Dopisek niesie liczbę sesji z migawki, nigdy liczbę wymyśloną. Montaż bez
tożsamości klienta zostaje przy wykazie informacyjnym. Inaczej ta sama odpowiedź rdzenia pojawiałaby się
w dwóch postaciach zależnie od tego, kto ją wywołał, gdyby menu wiersza budowało własne miejsce na treść.

## budowa/klient-poprzedni/src/okna-pomocnicze/okno-historii-rozmowy.ts
Panel nie jest drugim widokiem rozmowy: kolumna rozmowy pokazuje turę
bieżącą rosnącą na żywo, panel pokazuje zapis trwały tego samego okna od
pozycji najnowszej wstecz, stronicowany kursorem czasu, i nie zna komendy
nadania wiadomości. Zamiast pytania o potwierdzenie przed czynnością stoi
odpowiedź po niej: rdzeń oddaje liczbę usuniętych pozycji, a panel mówi ją
wprost. Wyjątkiem jest wyczyszczenie całej historii okna — ten sam wyjątek,
który ma kasowanie sesji, obejmujący czynność nieodwracalną i całościową;
usunięcie pozycji wskazanych ręcznie potwierdzenia nie wymaga. Wszystkie
przyciski są czynne zawsze — brak wskazanych pozycji nie wyszarza
przycisku, tylko wraca zdaniem, co się stało. Gospodarz bywa bez okna
nadanego przez rdzeń; odczyt historii wymaga wtedy czegoś, czego nie ma, i
panel mówi to wprost zamiast pokazać pustą listę, podczas gdy nastawa
retencji zostaje czynna, bo zakres globalny okna nie potrzebuje.

Modal czyszczenia woła rdzeń sam, bo to on ma pokazać odpowiedź w chwili
czynności — odmowa i liczba usuniętych padają tam, gdzie Operator patrzy, a
panel powtarza je u siebie i przeładowuje wykaz.

## budowa/klient/src/wejscie/katalog-tresci.test.ts
Łańcuch dopisany poza plikiem treści jest usterką i ma tu upaść. Reguła obejmuje pliki okna — składniki, ekrany, montaż, przebieg, narzędzia i zestaw znaków. Nie obejmuje sprawdzianów: opis sprawdzianu jest zdaniem dla tego, kto czyta wynik uruchomienia, i do okna nie trafia nigdy. Wykaz plików pominiętych jest wypisywany, więc pominięcie nie da się rozrosnąć po cichu.
Reguła jest dwuczłonowa i cała maszynowa: poza katalogiem żaden łańcuch nie niesie polskiego znaku diakrytycznego, a każdy łańcuch ma kształt techniczny — jest nazwą bez odstępu, wykazem klas albo selektorem, wzorem z podstawieniem, samym odstępem rozdzielającym węzły tekstowe, rysunkiem w zestawie znaków albo wpisem diagnostycznym. Zdanie żadnego z tych kształtów nie ma i tu upada.
Odstępstwo jest jedno i wąskie: wpis diagnostyczny — łańcuch oddany do dziennika wywołaniem konsoli albo niesiony wyjątkiem. To zdanie dla wykonawcy sprawdzianu, okno go nie pokazuje. Wpisy są liczone i wypisywane, więc odstępstwo nie rozrośnie się po cichu.
Czego sprawdzian nie wychwyci: pojedynczego słowa bez odstępu i bez polskiego znaku, na przykład nazwy własnej wpisanej wprost w składnik. Granica jest nazwana wprost, bo instrument, który udaje szczelność, jest gorszy od instrumentu o znanym zasięgu.
Bez pomijania komentarzy pomiar mierzyłby polszczyznę komentarzy zamiast łańcuchów, a bez pomijania wyrażeń regularnych ukośnik klasy znaków wyglądałby jak początek komentarza i zjadał resztę pliku.

## budowa/klient-poprzedni/src/moduly/research/zlecenia-badania.ts
Zlecenia nazywają to, co okno zebrało z formularza, zanim warstwa kontraktu przełoży to na treść
żądania. Rozdzielenie jest celowe: okno mówi po polsku o zakresie badania i wiarygodności źródła,
a plik źródła komend mówi nazwami kontraktu — dzięki temu zmiana nazwy pola w kontrakcie dotyka
jednego pliku, nie pięciu okien.

## budowa/klient-poprzedni/src/moduly/translate/pseudolokalizacja.ts
Trzy rzeczy szuka się w interfejsie przed tłumaczeniem, i trzy zabiegi je pokazują: czy pole
zniesie znaki diakrytyczne, więc litery zamieniane są na warianty z diakrytykami, a sylwetka słowa
zostaje czytelna; czy pole zniesie dłuższy tekst, więc treść jest dopełniana do zadanego
wydłużenia, bo przekład bywa dłuższy od źródła; czy tekst nie jest sklejany z kawałków, więc treść
dostaje ramkę, żeby ucięcie i sklejenie było widać na pierwszy rzut oka. Symbole zastępcze
i znaczniki formatu przechodzą nietknięte, bo przekształcenie ich zapisu zamieniłoby test wyglądu
w usterkę podstawienia, a to jest dokładnie ta klasa błędu, której pseudolokalizacja ma nie
wprowadzać. Rdzeń o niczym tu nie wie: kontrakt nie ma komendy pseudolokalizacji, więc wynik jest
wyłącznie do odczytania i przeniesienia ręcznego, a okno mówi to wprost przy wyniku, zamiast
pozwolić sądzić, że coś zapisano.

## budowa/klient-poprzedni/src/okno-komunikacji/pasek-zlecenia.ts
Na pasku ustawia się całe otoczenie zadania, zanim treść pójdzie do modelu: gdzie się wykona, na czym, jakim modelem, z jakim wysiłkiem i w jakim trybie zatwierdzania — to jedna decyzja rozłożona na kilka nastaw, więc stoi w jednym rzędzie, bo dwa rzędy byłyby dwiema kopertami. W rzędzie stoją: urządzenie wchodzące z zewnątrz jako gotowy element, katalog roboczy, model, wysiłek, tryb zatwierdzania i rozszerzenia. Mikrofon jest zbudowany, ale nie staje, dopóki kontrakt nie ma komendy przenoszącej nagranie z pamięci karty na maszynę silnika — pasek jest wtedy krótszy, bo wyszarzona ikona byłaby bramą, a ikona odmawiająca po naciśnięciu atrapą. Nie mają portu w ogóle: projekt, bo rdzeń nie zna komendy dającej wykaz projektów; izolacja i drzewo równoległe, które są poziomami zasięgu, nie sterami; wersjonowanie i zatwierdzanie zmian, należące do koperty budowniczego agenta, nie okna czatu; załączniki, bo komenda wysyłki przyjmuje odwołania, a nie plik wprost; oraz sposób podania, wobec braku komendy dokładającej pozycję do kolejki — pustych uchwytów na te sprawy tu nie ma. Dwa ostatnie stery, rozszerzenia i mikrofon, wchodzą inaczej niż pozostałe: stery koperty czytają migawkę okna i nic więcej nie potrzebują, a rozszerzenia i mikrofon trzymają własne subskrypcje, których pasek pozbywa się przy zejściu okna. Mikrofon staje dopiero po zapytaniu o dostępność silnika, więc wchodzi do rzędu z opóźnieniem jednej odpowiedzi rdzenia, a miejsce ma przygotowane z góry, żeby kolejność rzędu nie zależała od czasu odpowiedzi. Pasek nie buduje menu i nie zna komend: menu buduje osobny komponent, wysyłkę prowadzi osobne źródło zlecenia, a ten plik jedynie ustawia stery w rzędzie i podaje każdemu wycinek migawki, który do niego należy.

## budowa/klient-poprzedni/src/protokol/archiwum-sesji.ts
Jedna odpowiedzialność: trzy komendy jednego przejścia — sesja wychodzi z historii bieżącej i wraca
do niej. Zapis zostaje w całości, więc archiwizacja nie jest usunięciem i widok nie nazywa jej tak.
Wszystkie trzy działają na wielu sesjach: żądania niosą wykaz identyfikatorów, a odpowiedzi oddają
wykaz faktycznie przeniesionych. Rdzeń może przenieść część zbioru, więc zwrócony wykaz idzie do
widoku nietknięty. Wgląd jest stronicowany: wykaz archiwum przyjmuje przesunięcie i limit i oddaje
łączną liczbę wszystkich sesji archiwum, nie długość strony.

## budowa/klient-poprzedni/src/protokol/bieg-sesji.ts
Obie komendy nie są swoimi odwrotnościami. Zatrzymanie przerywa tury biegnące w oknach sesji, a sesja
zostaje otwarta i gotowa na kolejną turę; wznowienie otwiera sesję zamkniętą wraz z jej oknami. Widok,
który postawiłby je jako parę start-stop, obiecałby coś, czego rdzeń nie robi, więc każda z nich
pojawia się osobno i pod własnym warunkiem. Obie oddają zakres skutku. Wykaz zatrzymanych okien mówi,
w ilu oknach turę naprawdę zatrzymano — zero znaczy, że nic nie biegło, i to jest wynik poprawny, nie
odmowa. Wykaz okien po wznowieniu niesie okna otwarte razem z sesją, więc powrót do niej nie musi ich
odpytywać drugi raz.

## budowa/klient-poprzedni/src/protokol/nazwa-sesji.ts
Bieg sesji niesie osobny plik biegu sesji, przynależność plik projektu sesji, a odłożenie plik
archiwum sesji. Kopia jest osobnym bytem: kopiowanie oddaje nową sesję wraz z liczbą przeniesionych
wiadomości — kopiowanie nie jest odnośnikiem do źródła, więc widok mówi wprost, ile zapisu naprawdę
poszło. Pusta nazwa kopii nie jest błędem. Kontrakt pozwala pominąć tytuł i wtedy rdzeń bierze nazwę
źródła z dopiskiem. Klient nie składa tego dopisku sam — nazwę nadaje ta strona, która ją zna.

## budowa/klient-poprzedni/src/protokol/projekt-sesji.ts
Żądanie przypisania niesie albo identyfikator projektu, albo jego nazwę; przy pustym identyfikatorze
rdzeń zakłada projekt o podanej nazwie i oddaje jego identyfikator, więc klient nie zakłada projektu
osobną komendą. Wykaz przeniesionych może być krótszy od żądania — rdzeń oddaje sesje faktycznie
przeniesione, więc widok melduje z wyniku, nie z treści żądania.

## budowa/klient-poprzedni/src/strona-glowna/wpiecie-komponentow.ts
Docelowe środowisko kafla wynika z macierzy widoczności, którą rdzeń podaje w odpowiedzi o środowiskach;
kopia tej macierzy po stronie klienta mogłaby rozjechać się z bazą po cichu, więc jej tutaj nie ma.
Zapytanie o środowiska pada tu osobno, z dołączeniem modułów, bo wpięcie środowisk pyta bez tego pola —
kartom strefy pierwszej moduły nie są potrzebne. Podobnie wykaz modułów wołają dwa wpięcia: to po kafle
strefy drugiej, a wpięcie modułów po kafle modułów poza nawigacją; scalenie wołań wymaga zmiany w warstwie
montażu strony głównej, poza tym plikiem — strona wstaje raz na wejście, nie w pętli, więc dwa wołania są
ceną za niezależność obu wpięć. Kafel jest zawsze klikalny i zawsze prowadzi do pracy: gdy macierz jeszcze
nie przyszła albo moduł nie stoi w żadnej bocznej nawigacji, przejście idzie do środowiska początkowego,
nadal ze wskazaniem modułu, zamiast pokazać odmowę lub kafel wyszarzony. Wykaz personalizowanych przychodzi
bez znacznika włączania wyłączonych, więc komponent wyłączony skraca listę zamiast dawać kafel wyszarzony;
odmowa rdzenia zostawia same kafle rodzajów, bo te wynikają z kontraktu, nie z tej odpowiedzi. Powłoka
otwiera moduł bez pozycji na liście drogą bezpośrednią, z pominięciem wykazu; zgubienie wskazania odbierałoby
kaflowi jedyne wejście do tych okien.

## budowa/klient-poprzedni/src/moduly/research/zrodlo-lektury.ts
Kontrakt ma własną komendę modułu — odczyt treści źródła niezależnie od tego, czy jest ono plikiem
repozytorium — a rdzeń nie ma dla niej jeszcze uchwytu; droga dziś przejezdna jest podgląd zasobu
repozytorium wraz z podziałem na strony i znacznikiem skrócenia. Ta droga wystarcza dla źródeł
związanych z dokumentem repozytorium i tylko dla nich — źródło typu strona internetowa, notatka
albo zbiór danych, którego rdzeń nie trzyma jako pliku, treści nie ma dziś skąd wziąć. Okno nazywa
to Operatorowi zamiast pokazywać pusty czytnik; dobudowa uchwytu własnej komendy tę granicę
zdejmuje. Mechanizm podglądu jest wspólny z podglądem plików modułu biblioteki — Research nie
buduje drugiego czytnika obok tamtego.

## budowa/klient/src/wejscie/skladniki/belka-okna.ts
Znak w belce dorównuje wielkością kontrolkom, więc niesie kropkę w barwie sygnału: przygaszony monochromat czyta się jak brakująca ikona.

## budowa/klient/src/wejscie/skladniki/fraza-nawigacyjna.ts
Fraza stoi w treści, nie w pasie działań — w pasie zostają same czynności, inaczej zdanie, czynność poboczna i główna nie mieszczą się w jednym wierszu.

## budowa/klient/src/wejscie/skladniki/pasek-postepu.ts
Miara jest daną, a wygląd należy do arkusza stylu. Rolę paska postępu niesie tor, bo to on wyraża postęp; nazwa toru opisuje rzecz, a nie powtarza etykietę nad nim.

## budowa/klient/src/wejscie/skladniki/pole-sesji.ts
Samo hasło "pozostań zalogowany" nie niesie ani jednego, ani drugiego bez opisu skutku. Zgoda idzie do rdzenia jako pole trwania sesji: sesja bramki dostaje trwanie długie zamiast doby roboczej. Bramką nie jest — znosi powtarzanie logowania, niczego nie blokuje.

## budowa/klient-poprzedni/src/strona-glowna/pytanie-o-nazwe.ts
Z modalu korzystają czynności zmiany nazwy, kopii i przypisania projektu. Warstwę tła daje pseudoelement
podkładu, a stos okien, pułapkę ogniska i zamknięcie klawiszem Escape zapewnia przeglądarka. Zamknięcie
oznacza odmowę, nie nazwę pustą: Escape i przycisk odmowy dają brak wartości, więc wywołujący odróżnia
odmowę od świadomie pustego napisu. Modal jest doklejany na czas pytania i usuwany po odpowiedzi, żeby
pytania nie nawarstwiały się w drzewie przy każdym wierszu wykazu.

## budowa/klient/src/wejscie/skladniki/naglowek-ekranu.ts
Głowa ekranu ma własny rytm odstępów; wyjęcie toru kroków poza nią rozstraja ten rytm.

## budowa/klient-poprzedni/src/protokol/kanal.ts
Nazwy komend i zdarzeń oraz kształty ich treści pochodzą wyłącznie z kontraktu współdzielonego:
zmiana nazwy w pliku kontraktu przerywa kompilację klienta. Komunikat nierozpoznany nie jest
odrzucany. Zapis i wpis do konsoli są jedyną reakcją na komunikat nierozpoznany: ani zdarzenie
zapasowe, ani koperta o typie spoza kontraktu nie zrywa połączenia i nie blokuje sesji.

## budowa/klient-poprzedni/src/protokol/koperta.ts
Kształt koperty pochodzi w całości z kontraktu współdzielonego — plik nie definiuje własnego typu
komunikatu i nie powiela ani jednego literału nazwy. Odpowiada wyłącznie za nadanie kopercie
wychodzącej identyfikatora oraz czasu nadania i za odczytanie pól odpowiedzi z koperty przychodzącej.
Obecność pola statusu odróżnia odpowiedź od zdarzenia i od fragmentu strumienia. Rzutowanie treści
koperty jest świadome: warstwa transportu nie waliduje ładunku, a błąd kształtu dotyczy wyłącznie
bieżącego komunikatu.

## budowa/klient-poprzedni/src/protokol/ramka.ts
Odczyt jest fail-open: ramka nieczytelna albo o kształcie niezgodnym z kopertą nie zrywa połączenia
i nie blokuje sesji — wraca jako zdarzenie zapasowe z zachowaniem treści surowej w ładunku. Zdarzenie
zapasowe dla ramki nierozpoznanej wskazuje kontrakt — klient nie wybiera go samodzielnie.

## budowa/klient-poprzedni/src/moduly/translate/skroty-translate.ts
Nasłuch wisi na elemencie modułu, nie na dokumencie: moduł znika z drzewa przy zejściu ze sceny
i nasłuch znika razem z nim, a nasłuch dokumentu trzeba by odpinać osobno, więc pierwszy przeoczony
byłby wyciekiem — to samo rozstrzygnięcie, co przy sterze kanału. Skutek uboczny tej decyzji jest
nazwany, nie przemilczany: skrót działa, gdy ognisko stoi wewnątrz modułu, a poza modułem klawisze
należą do powłoki. Wyszukiwarka funkcji zajmuje kombinację klawisza K za opracowaniem modułu,
a ten sam skrót wiąże na dokumencie pole poleceń paska górnego: dopóki ognisko stoi w module,
skrót otwiera wyszukiwarkę modułu i nie idzie dalej, poza modułem prowadzi do pola poleceń bez
zmian, a zbieżność jest realna i pierwszeństwo ma okno, w którym operator pracuje. Dwa skróty
pętli wykonawczej stoją w wykazie, ale moduł ich nie wiąże, bo Execution Loop Window jest oknem
wspólnym platformy i leży poza katalogiem tego modułu — wykaz mówi to wprost, zamiast pomijać
pozycje i sugerować, że skrótów nie ma.

## budowa/klient-poprzedni/src/moduly/research/zrodlo-odkrywania.ts
Cztery tryby zapytania istnieją w zamierzeniu modułu: webowy, naukowy, własny semantyczny
i pełnotekstowy. Kontrakt niesie dziś dwa ostatnie i oba są tu wywołane naprawdę — własny
semantyczny szuka po znaczeniu w bibliotece wiedzy Operatora i w przestrzeni roboczej okna
badania, pełnotekstowy szuka po treści zasobów repozytorium. Tryb webowy i naukowy mają
w kontrakcie własną komendę modułu, gdzie o trybie rozstrzyga pole żądania, ale rdzeń nie ma dla
niej jeszcze uchwytu. Panel ich nie udaje i nie gasi: zapytanie wychodzi pod nazwą tej komendy,
a odmowa rdzenia jest tym, co Operator zobaczy — ta sama droga, jaką opisuje plik akcji okien. Gdy
uchwyt dojdzie, oba tryby przeniosą się tutaj, obok dwóch już działających. Źródło nie ma stanu:
wyniki mieszkają w oknie, bo są zapytaniem, nie treścią badania — do badania wchodzą dopiero przez
komendę dodania źródła.

## budowa/klient-poprzedni/src/strona-glowna/podglad-strony-glownej.ts
Wzorzec jest ten sam co w podglądzie zetonów motywu: widok daje się obejrzeć bez montażu w powłoce
aplikacji, nie wchodzi do pakietu głównego i niczego z niego nie importuje. Bez parametrów w adresie
obowiązuje preferencja systemu i brak środowiska czynnego.

## budowa/klient-poprzedni/src/moduly/studio/dymki-komentarzy.ts
Na marginesie dokumentu stają trzy różne byty i operator ma po samym wyglądzie wiedzieć, na
który patrzy: komentarz mówi o fragmencie i nie niesie brzmienia, więc treści nie zmienia ani
teraz, ani po rozwiązaniu wątku; propozycja zmiany niesie brzmienie fragmentu, które jeszcze
nie weszło w treść, i operator ją przyjmuje, odrzuca albo poprawia; zmiana śledzona jest już
w treści i czeka na decyzję. Zlanie ich w jedną kartę odebrałoby operatorowi rozróżnienie, po
którym poznaje, czy dokument już się zmienił, dlatego każda karta niesie własny rodzaj, własny
nagłówek i własne zdanie o skutku decyzji, a arkusz stylów daje im trzy różne obramowania.
Komentarz stoi jako dymek przy miejscu w treści, nie w osobnym wykazie: kotwica w treści
wskazuje fragment, a karta ustawia się na jej wysokości, natomiast spis do przejścia stoi
osobno w przyborniku znakowania — jedno nie zastępuje drugiego, dymek mówi „tu", spis mówi
„ile jeszcze". Kontrakt nie niesie usunięcia komentarza, tylko dodanie, wykaz i rozwiązanie,
dlatego przycisku usuwającego nie ma, a jego miejsce zajmuje rozwiązanie wątku wraz ze zdaniem
o tej różnicy. Karta propozycji brzmienia różni się od zmiany śledzonej tym, że decyzja idzie
osobną komendą dla propozycji, nie dla zmiany śledzonej — dwie różne komendy dla dwóch różnych
bytów. Karta zmiany śledzonej stoi na marginesie obok oznaczenia w treści, bo oznaczenie mówi
gdzie, a karta mówi, co było przed zmianą, i to drugie jest tym, czego operator potrzebuje do
decyzji.

## budowa/klient/src/wejscie/skladniki/pole-tekstowe.ts
Wygląd obwódki i ogłoszenie czytnika ekranu biorą się z jednego stanu błędu, nie z osobnej klasy — inaczej rozjeżdżają się.

## budowa/klient-poprzedni/src/protokol/ksztalt-odpowiedzi.ts
Kanał nie waliduje ładunku: rzutuje go na typ zapowiedziany przez kontrakt i oddaje wywołującemu.
Rzutowanie jest obietnicą kompilatora, nie rdzenia — rdzeń starszej wersji albo pośrednik może
przysłać treść bez pola obowiązkowego, a widok dostałby wartość pustą w miejscu, w którym typ obiecuje
tablicę. Sprawdzian zamienia taką odpowiedź w zwykłe niepowodzenie wywołania: wpis do dziennika
i wynik z błędem walidacji.

## budowa/klient-poprzedni/src/moduly/research/zrodlo-okna-badania.ts
Moduł pyta o okna, bo każda komenda obszaru research wymaga pola windowId, czyli okna badania.
Klient nie wymyśla tego identyfikatora: bierze go z rdzenia komendą wykazu okien zawężoną do sesji
i do modułu research. Zaszyty identyfikator byłby daną zmyśloną, a okno bez wskazania rdzenia
pokazuje stan pusty, zamiast udawać, że wie. Wszystkie trzy komendy mają uchwyt w rdzeniu; akcja
panelu jest drogą wykonania dla akcji, których kontrakt nie rozróżnia osobną komendą, na przykład
cytowania, grupowania automatycznego czy porównania źródeł.

## budowa/klient-poprzedni/src/strona-glowna/pozycje-ustawien.ts
Strefa ma najniższą wagę wizualną i jest stale dostępna. Wykaz kodów strefy trzeciej liczy trzy pozycje:
okno konfiguracji, mobile i wyświetlanie stałe. Pozostałe pozycje — uwierzytelnianie operatora, punkty
izolacji, dostępy oraz modele i tożsamość — są zakresami okna konfiguracji; własne wejście na listwie
dawałoby dwie drogi do tej samej rzeczy. Pełny wykaz niesie dalej menu aplikacji w pasku, gdzie hierarchii
stref nie ma.

## budowa/klient-poprzedni/src/moduly/studio/dziennik-czynnosci.test.ts
Sprawdziany dziennika czynności mierzą cztery rzeczy, których zgubienie kosztowałoby operatora
pracę albo zaufanie: odmowę cofnięcia rozpoznaną jako odmowa mimo udanego wywołania, zależność
nazwaną wprost zamiast schowaną za ogólnym „nie udało się", zdanie o różnicy drzew obecne przy
każdym cofnięciu oraz liczbę zmian operatora zachowanych przy cofaniu pracy modelu. Sprawdzianu,
że przycisk woła komendę, tu nie ma, bo taki sprawdzian mierzyłby tylko atrapę, którą sam stawia.

## budowa/klient-poprzedni/src/moduly/translate/stan-okna-translate.ts
Faza siedzi w atrybucie danych, żeby arkusz mógł odróżnić błąd od pustki bez mnożenia klas, a
sprawdzian mógł zapytać o stan okna, nie o jego wygląd. Forma jest ta sama co w pozostałych
modułach: ikona, tytuł, opis, treść wyśrodkowana, ze stopniami pisma i szerokością łamania z arkusza
wspólnego, którego moduł u siebie nie nadpisuje. Ikona należy do stanu pustego i tylko do niego,
bo arkusz modułu chowa ją w pozostałych fazach: ładowanie ma własny nośnik, błąd kreskę po lewej.
Komunikat przesłania treść, nie kasuje jej: nieudane odświeżenie zostawia to, co już było widoczne,
więc powrót do treści nie wymaga ponownego odczytu. Ładowanie nie chowa treści i jest to decyzja
modułu — treścią każdego z trzech okien jest formularz, do którego operator właśnie pisze, więc
schowanie go na czas zapisu zabrałoby z oczu to, co wpisane; kontrolka zostaje widoczna i
edytowalna, a wskaźnik odczytu stoi obok niej, nigdy zamiast niej. Bez własnego przycisku ponowienia:
wszystkie trzy okna stoją na zapisach wyzwalanych z formularza, więc ponowieniem jest ten sam
przycisk, który czynność wywołał, zostający klikalny także po odmowie. Podział z bytem wspólnym
niesie zestaw faz i znakowanie powłoki, a tutaj zostaje to, czego wspólny byt nie przesądza: klasy
modułu, treść komunikatu, nośnik ikony i spinnera oraz pozostawienie treści widocznej na czas
ładowania.

## budowa/klient/src/wejscie/skladniki/kolumna-tozsamosci.ts
Kolumna nie jest planszą marki — jest strefą okna, więc powierzchnia różni się od panelu treści odcieniem, nie kontrastem. Nota wydawcy nie należy do tego składnika: stoi w wierszu pasa działań, niżej niż kolumna, i jest jedna dla całego okna. W odsłonie przygotowania zamiast zdań stoi pole animacji powłok: na tym etapie użytkownik już nie wybiera programu, tylko czeka, aż się złoży.

## budowa/klient-poprzedni/src/okna-pomocnicze/okno-podglad-bash.ts
Jedno wywołanie zapisu na zbiorcze wyjście robi dwie rzeczy naraz: zapisuje
okno na wyjście wszystkich otwartych kart terminala i oddaje ogon historii;
nowe wiersze dochodzą potem zdarzeniem strumienia z identyfikatorem tego
okna. Trzy stany dostają trzy różne zdania: odmowa rdzenia pokazuje stan
błędu z treścią odmowy, nigdy pustą listę, bo brak odpowiedzi to co innego
niż brak danych; ogon pusty oznacza udany odczyt bez wierszy, co po
restarcie rdzenia jest prawdą o nim, a nie ukrytą stratą, bo dziennik
wyjścia jest pierścieniem w pamięci żyjącym jeden bieg rdzenia; zapis
niedoszły oznacza, że historia wróciła, ale okno nie jest zapisane na żywo,
bo żądanie poszło bez identyfikatora okna. Widok trzyma ograniczoną liczbę
ostatnich wierszy, bo okno pomocnicze nie jest drugą konsolą, a ucięcie jest
oznaczone, nie zamaskowane.

## budowa/klient/src/wejscie/skladniki/kroki-odzyskiwania.ts
Numer ustępuje ptaszkowi po zrobieniu kroku, więc stan nie stoi na samej barwie — wymaganie kontrastu spełnione jest przez kształt, nie tylko przez kolor. Wstawianie tu osobnego rysunku dublowałoby ten sam znak. Strzałki między krokami są rysunkiem, nie treścią — czytnik ekranu je pomija.

## budowa/klient-poprzedni/src/strona-glowna/pozycje-srodowisk.ts
Kody pochodzą z kontraktu i odpowiadają kolumnie kodu środowiska; opisy trzymane są w rekordzie
kluczowanym tym typem, więc zmiana kodu w kontrakcie przerywa kompilację klienta zamiast dawać cichą lukę.
Środowisko o kodzie spoza kontraktu zostaje na ekranie z godłem zastępczym — wykaz kontraktu jest
informacyjny, nie bramą.

## budowa/klient-poprzedni/src/protokol/uzgodnienie.ts
Kolejność wynika z kontraktu: powitanie połączenia, założenie sesji, otwarcie okna. Okno jest bytem
pośrednim między sesją a wiadomością, więc wiadomość można wysłać dopiero po jego otwarciu.
Niepowodzenie etapu nie blokuje połączenia ani kolejnych prób — dotyczy wyłącznie bieżącego
wywołania. Uzgodnienie wita rdzeń w chwili nawiązania połączenia, czyli zanim podano hasło, więc
niesie wtedy token sesji poprzedniej albo żaden. Po wejściu przez bramkę rdzeń musi dowiedzieć się,
które gniazdo należy teraz do której sesji; inaczej zerowanie hasła rozłączyłoby tego, kto właśnie
zmienił hasło. Powitanie jest odczytem — nie zakłada ani sesji pracy, ani okna — więc powtórzenie
niczego nie dubluje. Ciągu dalszego uzgodnienia to wywołanie nie uruchamia; tamten idzie przy
nawiązaniu połączenia. Pole tożsamości jest tym samym identyfikatorem klienta, którego żądają inne
wejścia platformy. Widok bierze je stąd, zamiast wołać budowę tożsamości klienta po raz drugi: każde
wywołanie nadaje identyfikator nowy, a ognisko jest właściwością klienta — drugi identyfikator
rozdzieliłby ognisko od połączenia, które je zgłosiło. Token dostarcza funkcja, nie wartość, ponieważ
powitanie idzie przy każdym nawiązaniu połączenia — także po zerwaniu i ponownym połączeniu — a sesja
bramki może się między nimi zmienić: wejście, wylogowanie, wygaśnięcie. Wartość zamrożona przy
składaniu warstwy niosłaby stan sprzed uruchomienia i wiązałaby połączenie z sesją, której już nie ma.
Warstwa protokołu nie zna pochodzenia tokenu: magazyn sesji należy do warstwy uwierzytelnienia, stąd
wstrzyknięcie od składającego. Brak tokenu nie jest błędem — rdzeń odpowiada brakiem uwierzytelnienia,
a ekran logowania i tak stoi nad aplikacją.

## budowa/klient-poprzedni/src/okna-pomocnicze/panel-pomocniczy.ts
Tytuł, czynności i stany treści należą do samego panelu, nie do umowy. Plik
nie zna kodów paneli i niczego nie buduje: mapowanie kodu na wytwórnię stoi
w pliku wytwórni paneli, a spis opisowy w rejestrze pomocniczych. Umowa nie
wymaga ramy okna, choć panel wolno w nią ubrać — okno podglądu w tle tak
robi — a umowa i tak oddaje wyłącznie sam element.

## budowa/klient-poprzedni/src/moduly/research/zrodlo-research.ts
Zbiór źródeł, ustaleń i raportu mieszka w stanie badania, żeby pięć okien patrzyło na jeden
komplet, a nie na pięć kopii. Wszystkie pięć komend ma uchwyt w rdzeniu, wpięty przy rejestracji
modułu badań i wypełniony w montażu portów; odmowa merytoryczna wraca zwykłym błędem i okno
pokazuje ją wprost. Pięć komend głównych ma własne metody, bo ich żądanie składa się ze zlecenia
okna, a nie z samych pól; pozostałe komendy obszaru idą metodą ogólną, której żądanie składa plik
wywołań komend, znający zaznaczenie i wskazania okien — druga taka sama metoda per komenda byłaby
wieloma przepisaniami tego samego wywołania. Sprawdzian kształtu warstwy protokołu zna wyłącznie
ogólny wynik i o polu nieznanego typu nie wie; bez tego przełożenia odmowa nieznanej komendy
dochodziłaby do okna jako zwykły błąd i okno nie miałoby czego wypisać.

## budowa/klient/src/wejscie/skladniki/lista-etapow.ts
Znak zostawiony z poprzedniego stanu kłamie — ptaszek przy kroku w toku mówi, że rzecz jest skończona. Jeden składnik obsługuje cztery etapy łączenia i pięć etapów przygotowania: różni je wyłącznie gałąź katalogu, z której biorą się nazwy i miary.

## budowa/klient-poprzedni/src/moduly/translate/stan-translate.ts
Odpowiedź zapisu źródła niesie od razu komplet paneli: gdyby Source Panel prowadził własną kopię
tekstu, a Translation Panels własną kopię paneli, ta jednoczesność musiałaby być odtwarzana
ręcznie w dwóch miejscach i rozjeżdżałaby się przy pierwszej odmowie. Identyfikator okna jest
warunkiem wstępnym: bez niego moduł nie wysyła zapisu źródła ani dodania celu, mówi o tym wprost
zamiast wysyłać żądanie z pustym polem i pokazywać odmowę walidacji jako własną usterkę.
Obszar translate nie ma ani jednej komendy przyjmującej plik, więc bez źródła dokumentu moduł nie
miałby wejścia od strony dokumentu wcale. Jeden rejestr kanałów jest na moduł, nie jeden na ster:
sterów jest jeden plus liczba paneli, a rejestr zakładany osobno przez każdy z nich pytałby rdzeń
wielokrotnie o ten sam wykaz. Warsztat mieści pamięć jako byt operatora, segmentację, terminologię,
korektę, obieg, dokumenty, lokalizację, napisy, silniki i wymianę zewnętrzną.
Parametry wykonania okna są dodatkiem informacyjnym przy niepowodzeniu odczytu stanu (zasada
działania mimo błędu). Rejestr kanałów idzie obok kolejki komend: wykaz kanałów modelu nie zależy
od okna i nie jest warunkiem żadnej czynności, obsadza wyłącznie ster kanału, którego pusty wybór
jest poprawną wartością, więc czekanie na niego opóźniałoby pas kontekstu, a jego odmowa nie
zatrzymuje modułu. Moduł okna zmienia się osobną komendą, więc okno sesji bez modułu Translate
i tak jest tym oknem, w którym operator właśnie pracuje.

## budowa/klient-poprzedni/src/okno-komunikacji/profil-modulu.ts
Narzędzie paska promptu nie jest osobną komendą kontraktu — jest gotowym poleceniem wstawianym do pola wypowiedzi, które model wykonuje wywołaniem narzędzia platformy, dzięki czemu pasek nie zawiera ani jednego przycisku bez działania. Postać rozmowy modułu przyjmuje trzy wartości: zwykłe okno czatu dla większości modułów, pływający awatar z oknem dymkowym i głosem przed tekstem wyłącznie dla modułu Assistant, oraz brak rozmowy wyłącznie dla modułu Library, który jest menedżerem zasobów, nie środowiskiem pracy z modelem. Brak rozmowy nie znaczy pustki: okno trafiające na moduł bez rozmowy pokazuje operatorowi wprost, czym ten moduł jest zamiast czatu, bo pusta lista czyta się jak brak treści i jest błędem. Granica liczby okien równoległych jest ustalona wprost dla wybranych modułów — Apps i MultitaskingAI po cztery, Developer do czterech, Automations i Translate po dwa, Agents i Assistant po jednym — a pozostałe moduły dostają pełny sufit platformy. Ile okien rozmowy można otworzyć rozstrzyga zawsze funkcja, nie samo pole granicy: pole niesie granicę, funkcja niesie prawo do otwarcia.

## budowa/klient-poprzedni/src/moduly/studio/dziennik-czynnosci.ts
Dziennik odwracalny opisuje słowami trzy rzeczy, których nie wolno zgubić. Cofnięcie nakłada
różnicę drzew, nie migawkę: praca naniesiona po cofanej czynności zostaje, a operator, który
tego nie wie, nie odważy się cofnąć czynności ze środka dziennika, więc zdanie o tym stoi przy
każdej czynności, nie w pomocy. Czynność będąca podstawą późniejszej odmawia cofnięcia i nazywa
te czynności — rdzeń oddaje je w zależnościach wpisu oraz w bilansie odpowiedzi, a okno pokazuje
powód, nie samo „nie udało się", inaczej operator zostaje z zamkniętą drogą bez wskazania, co ją
zamknęło. Wykonawca jest nazwany: wpis niesie imię i wersję wykonawcy modelowego, więc dwóch
agentów nie pokazuje się jako jeden model. Plik nie zna elementów strony ani rdzenia — wejściem
są byty kontraktu, wyjściem napisy i rozstrzygnięcia, dzięki czemu sprawdza się bez stawiania
okna. Przycisk cofnięcia czynności zablokowanej pozostaje czynny mimo zapowiedzianej odmowy,
bo prawdę rozstrzyga rdzeń, a wykaz czynności blokujących mógł się zmienić po ostatnim odczycie.
Bilans jest jedyną drogą, którą operator dowiaduje się o pominięciu, bo przemilczenie pominięcia
jest w tym module zakazane wprost. Wywołanie cofnięcia bywa udane, a mimo to nic nie cofa, gdy
rdzeń oddaje wykaz cofniętych pusty i bilans z pominięciem nazywającym zależność — okno musi to
odróżnić od powodzenia, inaczej pokazałoby cofnięcie po czynności, która nie zeszła.

## budowa/klient-poprzedni/src/okna-pomocnicze/panele-otwieralne.ts
Gdy stan zbudowane rejestru rozjeżdża się z istnieniem wytwórni, menu jest
krótsze, bo pozycja bez wytwórni nie wchodzi w ogóle — pas pomocniczych
zgłasza wtedy rozjazd osobno. Taka pozycja nie dostaje wiersza wygaszonego
ani zapowiedzi, bo nieczynny wiersz byłby bramką; jej miejsce jest w wykazie
pozycji nieotwieralnych, wraz z powodem wprost z rejestru. Plik nie buduje
paneli i nie zna DOM, nie rozstrzyga też, ile paneli wolno otworzyć naraz —
panel jest bytem otwieranym i zamykanym pojedynczo, bez limitu. Dla modułu
spoza rejestru obie funkcje oddają wykaz pusty, a zdanie o braku spisu
pokazuje gospodarz pasa pomocniczych.

## budowa/klient/src/wejscie/skladniki/metody-logowania.ts
Pusta lista dróg pobocznych nie mówi, że są jakieś, a lista z wygaszonymi pozycjami mówi, co można włączyć w ustawieniach. Składnik zwraca wykaz węzłów, nie jeden węzeł: etykieta i siatka metod są rodzeństwem w kolumnie panelu. Opakowanie ich w pudełko wprowadziłoby dodatkowy poziom, przez który odstęp kolumny liczyłby się raz zamiast dwa.

## budowa/klient-poprzedni/src/rozmowa/prowenancja.ts
Rdzeń nadaje prowenancję pierwszym fragmentem strumienia, przed jakąkolwiek
treścią odpowiedzi, także wtedy, gdy proces kanału w ogóle nie wystartuje.
Kształt odpowiada strukturze `injection.Prowenancja` po stronie rdzenia; nazwy
pól po stronie sieci są angielskie, po stronie modułu polskie. Prowenancja jest
wyłącznie opisem wywołania — niczego nie dopuszcza ani nie wstrzymuje.

Pola kodu kanału, adaptera, adresu i środowiska wykonania są dodatkowe wobec
pól procesu: serwer emituje je tylko dla kanałów bez własnego procesu (sieć,
echo), a dla kanału głównego uruchamianego lokalnie zostają puste.

## budowa/klient-poprzedni/src/rozmowa/stan-rozmowy.ts
Stan wysyłania jest jawny, ponieważ Operator ma widzieć, że tura biegnie,
zanim przyjdzie pierwszy znak odpowiedzi. Nie jest bramką: przycisk
zatrzymania pozostaje czynny w każdym stanie, a pole wpisywania nigdy nie jest
wyszarzane.

Próg sekund milczenia strumienia jest wysoki, bo model bywa cichy, kiedy
rozumuje — alarm po kilku sekundach byłby fałszywy. Po upływie progu przyczyna
(zerwane łącze, martwy proces kanału, długie rozumowanie) nie ma znaczenia:
Operator dostaje opis ciszy zamiast niezmiennego napisu „Wysyłanie…".

## budowa/klient-poprzedni/src/moduly/roundtable/braki-kontraktu.ts
Napis orzekający na stałe, czego kontrakt nie ma, przestaje być prawdą w dniu powstania komendy
i nikt go wtedy nie zdejmuje. Ten moduł przeszedł dokładnie przez taki dzień: obszar roundtable
niósł cztery komendy, a po scaleniu kontraktu niesie ich czterdzieści sześć — zdania mówiące, że
nie wykonuje danej czynności żadna komenda obszaru, stały się wtedy nieprawdą co do czterdziestu
dwóch czynności naraz. Stąd podział na dwa powody, nie jeden: powód braku obsługi dotyczy komendy,
która JEST w kontrakcie, a okno jej nie wywołuje — nazwa komendy przychodzi jako wartość
wyliczeniowa, nie napisem, więc zmiana nazwy w kontrakcie przerywa kompilację zamiast zostawiać
w dymku napis o komendzie, której już nie ma; powód braku czynności dotyczy sytuacji, w której
czynności nie da się wskazać jedną komendą, bo kontrakt nie ma jej w ogóle albo brakuje pola, nie
komendy. Czego ten plik nie wie: czy rdzeń ma obsługiwacz danej komendy — kontrakt nie daje takiego
odczytu, bo nie ma komendy wyliczającej komendy obsługiwane, jedynie okna modułu. Rdzeń odpowiadający
odmową nieznanej komendy na komendę obecną w wykazie jest z tego miejsca nierozpoznawalny;
rozpoznaje go dopiero odmowa po naciśnięciu, dlatego zdania mówią o obsłudze niezbudowanej,
a nie orzekają, gdzie dokładnie jej brakuje.

## budowa/klient-poprzedni/src/punkty-izolacji/komenda.ts
Kanał przyjmuje wywołanie zwrotne, a obszary tego okna czytają po kilka komend po kolei — splot
wywołań zwrotnych w takim ciągu jest nieczytelny. Wszystkie obszary okna wołają więc stąd, zamiast
powtarzać u siebie ten sam opakowujący zapis. Typowanie zostaje pełne: parametr generyczny wiąże
żądanie z odpowiedzią wprost z kontraktu, więc pomyłka w kształcie żądania przerywa kompilację
klienta zamiast wracać odmową rdzenia w czasie działania. Odmowa nie jest wyjątkiem: obietnica
spełnia się zawsze, także gdy rdzeń odmówił. Odmowa siedzi w polu powodzenia i w polu błędu, żeby
wywołujący rozstrzygnął ją zdaniem dla Operatora, a nie obsługą wyjątku, którą łatwo pominąć.

## budowa/klient-poprzedni/src/punkty-izolacji/indeks.ts
Reszta aplikacji zna stąd jedną czynność otwarcia okna Punktów Izolacji przez podanie kanału rdzenia.
Okno jest jedno na klienta i żyje między otwarciami — powtórne otwarcie wraca do obszaru, na którym
Operator skończył, tak jak Okno Konfiguracji, którego ten plik jest wierną kopią wzorca. Kanał
podajemy przy pierwszym otwarciu. Wywołanie z innym kanałem — po ponownym połączeniu z rdzeniem —
buduje okno na nowo, żeby komendy izolacji nie szły przez transport, którego już nie ma.

## budowa/klient-poprzedni/src/punkty-izolacji/stan-warstwy.ts
Warstwy są dwie, ale nie są dwoma trybami do wyboru. Warstwa domyślna to warstwa bazowa platformy:
obowiązuje przy każdej nowej sesji, karcie, roli, projekcie i oknie. Warstwa sesji nakłada się na
bazową bez jej zmiany i wygasa z zamknięciem karty sesji — to podkład i naklejka na nim, nie dwa
osobne magazyny. Warstwa jest stanem wspólnym, a nie polem jednego obszaru, bo dotyczy każdego
odczytu i każdego zapisu izolacji: kontekstu, zakresu technicznego, przypisania profilu i podglądu
polityki. Gdyby siedziała w jednej zakładce, pozostałe pytałyby rdzeń o coś innego, niż widać na
pasku — dlatego jej kontrolka jest w narzędziach ramy okna, a stan tutaj. Ten plik nie woła rdzenia.
Przełączenie warstwy w rdzeniu należy do pliku sterowania warstwą; tutaj zostaje wyłącznie to, co
okno wie o warstwie po odpowiedzi rdzenia, i powiadomienie obszarów o zmianie.

## budowa/klient/src/wejscie/skladniki/miernik-sily.ts
Warunek niespełniony niesie puste kółko, spełniony ptaszka; ptaszek w każdym stanie czytałby się jako zrobione, a barwa jako jedyna różnica łamałaby wymaganie kontrastu. Miernik wiąże się z polem przez atrybut danych, nie przez sąsiedztwo w drzewie — sąsiedztwo bywa różne w różnych oknach i cicho się rozjeżdża. Regułę oceny niesie przebieg; miernik jej nie powtarza.

## budowa/klient-poprzedni/src/strona-glowna/pozycje-komponentow.ts
Różnica wobec strefy pierwszej: karta środowiska prowadzi do przestrzeni pracy, a kafel komponentu do
zbudowania rzeczy, która w tej przestrzeni potem pracuje; różnicę niesie krój — nagłówkowy dla wejścia do
środowiska, bazowy półgruby dla zbudowania komponentu. Rodzaj modułu i rodzaj komponentu własnego to dwa
różne byty, więc kod spoza kontraktu jest tu błędem kompilacji, nie brakującym kaflem. Czas ostatniej
zmiany przychodzi w milisekundach epoki; sformatowanie należy do widoku, bo rdzeń nie zna strefy czasowej
operatora. O tym, co nastawia się na stronie głównej, rozstrzyga kolumna widoczności wystawiana
w kontrakcie — kolejny moduł tak oznaczony ma dostać kafel bez zmiany w kliencie. Zamkniętego zbioru
rodzajów pilnują nadal treści komponentu i wykaz rodzajów do założenia, bo zakładanie komponentu przyjmuje
wyłącznie rodzaj z wyliczenia; kafel i zakładanie to dwie różne rzeczy. Kafel rodzaju nie ma metadanych
i mieć nie może: rodzaj nie jest bytem w bazie, więc nie ma daty założenia ani stanu czynności — pole
nieobecne znaczy kafel rodzaju, nie że metadanych nie odczytano. Kafel rodzaju, jeden z czterech stałych,
prowadzi do zbudowania komponentu i identyfikatora nie ma; kafel personalizowany wskazuje komponent już
zbudowany, więc niesie jego identyfikator, po którym rozpoznaje się jeden kafel spośród wielu tego samego
rodzaju. Pusty ekran przed pierwszą odpowiedzią byłby gorszy niż cztery kafle, które i tak zostaną
przerysowane po odpowiedzi rdzenia wykazem modułów. Nazwy i wezwania zostają miejscowe wyłącznie dla tej
jednej klatki: kod rodzaju niesie kod, a nie napis na kaflu.

## budowa/klient-poprzedni/src/okno-komunikacji/profile-inzynieria.ts
Diagnostics nie ma okna telemetrii, ponieważ nie jest ustalone, co miałoby ono mierzyć, a rdzeń nie prowadzi tabel połączeń ani procesów sesji potrzebnych do takiego widoku. Brak pamięci sesyjnej modułu Agents oznacza dwa czyszczenia, nie jedno: przy zamknięciu okna i przy zmianie testowanego agenta, o każdym operator jest uprzedzany, bo czat gubiący wątek po cichu wygląda jak awaria.

## budowa/klient/src/wejscie/skladniki/pas-dzialan.ts
Pas urwany na granicy kolumny tożsamości czyta się jak niedokończony rysunek — stąd domyka okno przez całą jego szerokość. Czynności zbierają się przy prawej krawędzi, bo tam ręka szuka ich w każdym oknie tej rodziny, a gdy zostaje jedna, nie zawisa samotnie po lewej.

## budowa/klient-poprzedni/src/okna-pomocnicze/pas-pomocniczych.ts
Pas nie jest drugim złożeniem modułu: nie tworzy stanu modułu, nie zna
komend obszaru inżynierskiego i nie dotyka okien operacyjnych, składając
wyłącznie gotowe okna pomocnicze oraz spis pozostałych pozycji. Zamknięcie
jest obowiązkowe, bo okno podglądu trzyma subskrypcję strumienia wyjścia,
która bez tego wywołania zostaje żywa po zejściu pasa ze sceny.

Nie każdy moduł zna swoje okno przy montażu: część modułów montuje się od
razu z oknem sesji, a inne poznają okno dopiero z późniejszego wykazu okien.
Panele powstają na nowo przy zmianie okna, bo biorą je przy powołaniu — tą
samą drogą co kolumna paneli sceny okien równoległych: stare panele są
zamykane wraz z subskrypcjami rdzenia, nowe budowane z oknem właściwym.
Podanie tego samego okna drugi raz nie robi nic, bo przebudowa panelu już
pytającego o właściwe okno byłaby zerwaniem strumienia bez powodu.

Mapowanie kodu pozycji na wytwórnię stoi poza pasem, więc pas nie zna ani
jednego kodu, a dołożenie panelu nie wymaga jego edycji — tego samego
mapowania używa kolumna paneli sceny okien równoległych. Pozycja nazwana w
rejestrze zbudowaną, dla której nie ma wytwórni, jest rozjazdem między
spisem a kodem, stąd wpis w dzienniku zdarzeń zamiast cichej wartości
pustej udającej, że pozycji w spisie nie ma.

## budowa/klient-poprzedni/src/moduly/studio/filtr-historii.ts
Filtr po autorze jest już wykonalny: wersja niesie pole autora oraz znacznik kluczowości, więc
rozdzielenie zmian operatora od zmian modelu ma po czym przebiegać. Pole autora jest
nieobowiązkowe — wersje założone przed jego wprowadzeniem autora nie niosą, a filtr autora nie
odsiewa ich po cichu, tylko trzyma w wykazie z autorem pustym, bo odsianie ich byłoby ukryciem
historii przed operatorem. Autozapis idzie osobnym szeregiem: zapis samoczynny poznaje się po
tym, że nie jest wersją kluczową i nie ma etykiety własnej, a obie te cechy nadaje wyłącznie
operator — rozróżnienie stoi na osobnej komendzie nadania etykiety. Zapytanie puste o autora
przepuszcza wszystko, a wersja bez zapisanego autora przechodzi wyłącznie przy zapytaniu pustym,
bo przy szukaniu konkretnego autora nie wolno jej ani oddać, ani o niej zapomnieć, dlatego okno
pokazuje liczbę zawężenia obok wykazu.

## budowa/klient-poprzedni/src/okno-komunikacji/profile-tresc.ts
Powierzchnia tekstowa modułu Studio jest jedna: okno pracy z dokumentem niesie treść, podgląd wydania i różnicę jako tryby jednego widoku, zamiast osobnego edytora, kanwy tekstowej, okna podglądu i panelu różnic — pozostałe trzy okna operacyjne modułu powierzchni tekstowej nie mają.

## budowa/klient/src/wejscie/skladniki/pole-hasla.ts
Dwa znaki leżą w przycisku, a widoczność rozstrzyga arkusz stylu po stanie wciśnięcia — przełączanie znaków skryptem rozjeżdżałoby się ze stanem kontrolki. Etykieta po odsłonięciu hasła ma odpowiednik nakładany przez mechanikę okna.

## budowa/klient-poprzedni/src/moduly/translate/ster-kanalu.ts
To jest ster, nie wyświetlacz: etykieta uchwytu niesie wartość bieżącą nastawy, nazwę kanału,
a nie słowo Kanał; kliknięcie rozwija wybór, wybór zmienia nastawę, a nazwa rodzajowa idzie do
opisu dostępności mechanizmu, bo czytnik ekranu musi wiedzieć, czego dotyczy wartość, której nazwa
sama tego nie mówi. Podpis jest blokiem, nie etykietą powiązaną atrybutem for: etykieta bez niego
wiąże się z pierwszym potomkiem, który da się etykietować, a tym jest uchwyt menu — skutek byłby
dwojaki i oba razy zły, więc podpis nie jest sterem. Pusty wybór jest wyborem, nie brakiem: pozycja
kanału czynnego okna stoi w drzewie zawsze i jest domyślna, bo dokładnie to opisuje kontrakt —
dzięki niej jest droga powrotna do zachowania sprzed wskazania. Rejestr, który nie dotarł, nie
odbiera czynności: drzewo ma wtedy samą pozycję domyślną, a powód odmowy stoi przy niej jako opis,
i ster nie wygasza się, nie znika i nie zatrzymuje przycisku obok.
Menu rozwinięte wiesza nasłuch wskazania na dokumencie, żeby zamykać się po kliknięciu poza sobą;
instancja panelu, która znika z wykazu rdzenia, jest z dokumentu usuwana, więc gdyby jej ster był
wtedy rozwinięty, nasłuch zostałby na dokumencie na zawsze i wołałby do elementu, którego już nie
ma — panele przychodzą i znikają z każdą zmianą wykazu, więc to nie jest przypadek teoretyczny.
Kanał wyłączony w innym oknie znika z wykazu kanałów czynnych, a rdzeń odmówiłby przekładu
wskazaniem na niego; cichy powrót byłby jednak zmianą nastawy bez wiedzy operatora, więc
towarzyszy mu zdanie przy pozycji domyślnej.

## budowa/klient-poprzedni/src/strona-glowna/pozycje-modulow.ts
Kafel modułu bez nawigacji jest jedyną drogą do modułu, którego rdzeń nie pokazuje w bocznej nawigacji
żadnego środowiska: wykaz modułów oddaje takie moduły z pustym wykazem środowisk, a bez kafla ich okna
operacyjne byłyby nieosiągalne. Kryterium pochodzi z kontraktu, nie z wykazu nazw: pusty wykaz środowisk
modułu znaczy, że moduł jest dostępny wyłącznie ze strony głównej, więc plik pyta dokładnie o to i nie
wpisuje ani jednego kodu modułu — gdy rdzeń przypnie moduł do środowiska, kafel znika sam. To nie jest
druga macierz widoczności: macierz mówi, gdzie moduł stoi, a ten wykaz wyłącznie to, że nie stoi nigdzie.
Moduł nastawiany na stronie głównej nie wraca w wykazie modułów poza nawigacją: moduł potrafi spełniać
oba warunki naraz, a wtedy rysowałby się na jednym ekranie dwa razy, w dwóch sąsiednich strefach, choć oba
kafle wołałyby ten sam kod — pierwszeństwo ma strefa kafli komponentów. Warunek pyta o dane rdzenia, nie
o kod modułu, więc rozstrzygnięcie przenosi się samo wraz ze zmianą oznaczenia w bazie. Wezwanie kafla
w strefie komponentów jest zdaniem, którym rdzeń sam opisuje moduł; pusty opis wraca do wezwania zastanego
z tej samej klatki, co wykaz komponentów, tak samo jak w kaflu personalizowanym. Wezwanie opisuje moduł,
a nie obiecuje czynności, której kafel nie wykonuje: jedynym jego skutkiem jest przejście do modułu — kafel
zakładający komponent stoi osobno, w przyborniku listwy ustawień.

## budowa/klient-poprzedni/src/punkty-izolacji/katalog-izolacji.ts
Obszar Profili składa i pokazuje ten sam zestaw kluczy, co obszary kontekstu i zakresu technicznego —
profil jest nazwanym zestawem trzech przełączników kontekstu i ośmiu zakresów technicznych. Nazwy
stoją w jednym miejscu, żeby profil nie nazywał klucza inaczej niż zakładka, w której Operator ten
sam klucz przestawia. Plik trzyma wyłącznie nazewnictwo: stany przełączników przychodzą z rdzenia.
Słownik pokrywa cały kontrakt: brak wartości zatrzymuje kompilację.

## budowa/klient-poprzedni/src/moduly/studio/galeria-szablonow.ts
Szablony pochodzą z rdzenia: wykaz szablonów oddaje pismo, umowę, raport, notatkę i ofertę wraz
z polami do wypełnienia, a założenie dokumentu z szablonu wypełnia go wskazanymi wartościami —
obie komendy działały już wcześniej, tylko nie miały czym się pokazać w oknie. Miniatura jest
podglądem wyglądu złożonym z tego, co szablon o sobie mówi: nazwa, przeznaczenie, format i pola.
Obrazka miniatury kontrakt nie niesie i rdzeń obrazów nie generuje, więc miniatura jest wyrysem
układu — kartka z nagłówkiem nazwy i kreskami wierszy w liczbie pól, bo rysunek udający gotowe
pismo pokazywałby wygląd, którego szablon nie obiecuje. Przejście do warsztatu szablonów jest
nieobowiązkowe: gdy warsztat nie jest zamontowany, galeria działa wyłącznie do odczytu i
zakładania dokumentu, a gdy jest zamontowany, miniatura dostaje przejście do niego, bo szablon
zmienia się tam, nie w galerii. Kategoria szablonu jest fabryczna albo własna, bo innego podziału
kontrakt nie niesie, a wymyślanie działów na podstawie nazw byłoby porządkiem zgadniętym.

## budowa/klient-poprzedni/src/okno-komunikacji/profile-wiedza.ts
Roundtable nie podaje własnej liczby okien, ponieważ wpisanie tu wartości dublowałoby sufit sceny albo z nim kolidowało — wiążącym ograniczeniem pozostaje zawsze sufit platformy. Dymek Assistant mówi wprost, że kanału głosowego nie ma, zamiast milcząco zejść na tekst, ponieważ kontrakt i rdzeń nie niosą mowy.

## budowa/klient-poprzedni/src/strona-glowna/pozycje-personalizowane.ts
Kafel rodzaju z czwórki stałej prowadzi do zbudowania nowego komponentu; kafel personalizowany wskazuje
konkretną automatykę, eksperta czy projekt. Kafel nie jest osobnym bytem, tylko widokiem komponentu: wykaz
przychodzi jedną komendą i nie ma osobnej komendy utworzenia kafla. Wykaz komponentów bez włączenia
wyłączonych oddaje wyłącznie komponenty czynne, więc komponent wyłączony daje krótszą listę, a nie kafel
wyszarzony. Kafel bez skutku byłby atrapą, a krótsza lista atrapą nie jest.

## budowa/klient/src/wejscie/skladniki/pole-kodu.ts
Każdy zestaw rządzi się sam — mechanika wiąże pola grupami, więc kursor nie przeskakuje między odsłonami. Składnik zwraca wykaz węzłów: zestaw pól i stopka z odliczaniem są rodzeństwem w kolumnie panelu. Liczba pól jest właściwością, nie stałą tego pliku: rozstrzyga ją długość drogi potwierdzenia wydawanej przez rdzeń, a tę odczytuje się z rdzenia, nie z okna.
Rdzeń wydaje drogę dłuższą niż zestaw pól, więc wklejenie musi unieść ją w całości — inaczej przycięłaby się do liczby pól i rdzeń odmówiłby drogi, która przyszła listem poprawna.

## budowa/klient-poprzedni/src/moduly/roundtable/czytelnosc-glosow.ts
Dwa głosy tego samego modelu to dwa głosy, nie jeden: schemat debaty nie ma więzu jedności między
oknem a kanałem, a kontrakt powtarza to w opisie komendy dodania modelu. Kluczem wypowiedzi jest
więc identyfikator uczestnika, nigdy identyfikator kanału — dwie tożsamości jednego kanału muszą
różnić się na ekranie, bo inaczej dwa głosy czyta się jako jeden. Stąd znacznik mówcy przy każdej
wypowiedzi i zdanie o powtórzonym kanale. Rozróżnienie idzie układem i gęstością, nie samą barwą.
Wstęga barwna zostaje, ale sama nic nie niesie temu, kto barw nie rozróżnia — nośnikiem jest
znacznik mówcy w osobnej kolumnie (te same znaki co w składzie), nagłówek z pełną tożsamością przy
każdej wypowiedzi oraz odstęp: nowy mówca dostaje przerwę, ciąg tego samego mówcy zostaje ciasny.
Moderator nie jest uczestnikiem: rdzeń zapisuje jego interwencję kodem stałym, poza wykazem
uczestników, bo więz obcy do tabeli uczestnika debaty odciąłby ją od zapisu.

Moderator nie ma kanału ani persony, więc rdzeń podpisuje go stałą kodu moderatora; tę samą
wartość niesie interwencja moderatora. Uczestnicy dostają kod z przedrostkiem właściwym
uczestnikom, więc kolizja nie zachodzi.

Odczytu wykazu modeli żadne okno modułu dziś nie wywołuje, więc okno otwarte w trakcie debaty nie
zna jeszcze całego składu.

Miejsce w składzie nie ma skąd wziąć dla uczestnika składowi nieznanego, a zmyślenie go
rozjechałoby obie listy.

Tok rozumowania stoi osobno, bo rdzeń nie liczy go do treści wypowiedzi — sumuje wyłącznie
fragmenty tekstu. Doklejony do zdania dałby na żywo wypowiedź inną niż ta, którą za chwilę utrwali
rdzeń.

## budowa/klient-poprzedni/src/moduly/studio/indeks.ts
Moduł sam mówi, którym jest modułem, więc rozjazd między nazwą w rejestrze a rzeczywistością
jest niemożliwy. Moduł nie osadza się sam w dokumencie i nie zna powłoki: oddaje element,
a warstwa składająca decyduje, gdzie go postawić, dzięki czemu te same okna wchodzą i w obszar
roboczy powłoki, i w stanowisko sprawdzianu. Zamknięcie modułu odpina subskrypcję zmiany
dokumentu; pole jest w umowie rejestru modułów nieobowiązkowe, woła je dziś stanowisko
sprawdzianu, a powłoka zawoła je, gdy dostanie granicę życia widoku.

## budowa/klient-poprzedni/src/punkty-izolacji/obszar-kontekst.ts
Trzy punkty to historia wymiany, pamięć długoterminowa i bieżący stan roboczy. Każdy przyjmuje jedną
z dwóch wartości, odrębna albo współdzielona, z domyślną odrębna: rdzeń startuje z pełną izolacją,
współdzielenie jest zawsze decyzją Operatora, nigdy stanem narzuconym. Treść wyjaśnień odrębna albo
współdzielona przy każdym kluczu jest przepisana z definicji tego samego pliku rdzenia — jedno źródło
prawdy, żeby zdanie widziane przez Operatora nie rozjechało się z tym, co egzekwuje maszyneria. Obszar
woła komendy odczytu i zapisu kontekstu wprost; wysłanie komendy jest opakowane w obietnicę, a odmowa
rozstrzygana opisem odmowy. Poziom zasięgu i warstwa nie są tu zaszyte: oba przychodzą ze stanu
wspólnego okna. Zaszycie ich znaczyłoby, że selektor zasięgu w lewym panelu i wybór warstwy w pasie
narzędzi pokazują jedno, a odczyt idzie po drugie. Wartość obecna ma trzy stany, celowo nierysowane
identycznie: zapisana — odczyt na warstwie czynnej powiódł się, Operator świadomie rozstrzygnął tę
wartość na tym zasięgu i tej warstwie; domyślna — odczyt warstwy czynnej zwrócił brak zasobu, zapisu
Operatora nie ma, więc widok pokazuje wartość domyślną platformy jawnie oznaczoną jako domyślna, a nie
jako świadomy zapis; nieodczytana — odczyt zwrócił inną odmowę niż brak zasobu, to nie jest informacja
o stanie maszyny, tylko brak informacji, i tak też jest nazwana. Żaden z tych stanów nie jest rysowany
jako sukces i nie ma tu zapisu optymistycznego: widok zmienia się dopiero po odpowiedzi rdzenia, nigdy
przed nią. Wybór nowej wartości i przycisk zapisu zostają zawsze czynne; odmowa rdzenia jest meldowana
zdaniem, nie blokadą kontrolki.

## budowa/klient/src/wejscie/skladniki/zakladki-pigulki.ts
Pigułka w pełnym wypełnieniu czyta się jak przycisk działania, a to kontrolka wskazująca położenie, stąd stan wybrany niesie podświetlenie. Odzyskiwanie dostępu nie jest trzecią drogą — jest wyjściem z logowania, więc w zakładkach zostaje wtedy zaznaczone logowanie.

## budowa/klient-poprzedni/src/moduly/roundtable/diagnostyka-kanalow.ts
Diagnostyka mierzy to, co niesie rejestr kanałów rdzenia — nazwa, rodzaj, model, konto poświadczeń,
czynność kanału — zestawione ze składem debaty. Nie sięga po telemetrię kanału, bo telemetrii
kontrakt nie oddaje: ani czasu odpowiedzi kanału, ani liczby błędów, ani stanu połączenia.
Rozstrzygnięcie ważne dla Operatora jest jedno: kanał uczestnika nieznany rejestrowi. Zachodzi ono
naprawdę — uczestnika zakłada się z kanału wybranego w chwili dodania, a rejestr bywa odczytany
później i bywa węższy, gdy kanał zniknął z rejestru rdzenia. Wtedy uczestnik zostaje w debacie,
a jego kanał nie ma nazwy, więc panel podpisuje go identyfikatorem. Diagnostyka mówi to wprost,
zamiast zostawiać identyfikator bez wyjaśnienia.

## budowa/klient-poprzedni/src/moduly/translate/sufit-paneli.ts
Sufit instancji Translation Panels jest opisem stanu, nie bramką: górnej granicy liczby paneli nie
stawia nikt i plik tego nie zmienia, bo wpisanie progu z palca byłoby prawem wymyślonym po stronie
klienta, a rdzeń i tak przyjąłby panel ponad nim. Skąd wiadomo, że granicy nie ma: migracja bazy
nie ma więzu na liczbę wierszy panelu; żądanie dodania celu nie niesie ani jednego licznika, a
granica maksymalnej liczby elementów nie pada w całym kontrakcie; uchwyt zapisu panelu sprawdza
okno i język, po czym woła model i zapisuje panel bez licznika i bez odmowy przy jakiejkolwiek
liczbie paneli okna; siatka paneli układa się jako siatka dopasowująca kolumny, więc panele
dokładają się bez końca. Każdy panel kosztuje przy tym jedno wywołanie modelu, bo zapis panelu
liczy przekład przed zapisem, żeby odmowa modelu nie zostawiła pustego wiersza w bazie, więc
rachunek rośnie liniowo z liczbą paneli. Dwie liczby z tego drzewa nie są odpowiedzią na pytanie
o granicę i nota ma je rozróżnić: sufit gniazd sceny okien równoległych, do którego cały moduł
Translate się mieści, i osobna granica okien komunikacji profilu Translate — okno komunikacji
a panel języka to dwa różne byty. Forma noty idzie za wzorem sufitu uczestników z innego modułu,
a sufit sceny jest importowany, nie przepisywany, więc liczba jedzie z pliku, w którym stoi.

## budowa/klient-poprzedni/src/rozmowa/stopka-wpisu.ts
Rozliczenie tury pod wypowiedzią to stopka i blok błędów, nie transkrypt:
podsumowanie tury, konto kanału, typy zdarzeń i wykaz wytworów mówią, ile tura
kosztowała i co po sobie zostawiła, a nie co model powiedział. Widok wpisu
składa warstwy i steruje nimi; ten plik wie, jak wygląda jedna notka. O tym,
które części pokazać, ten plik nie rozstrzyga — rozkład warstw przychodzi
gotowy z widoku zapisu.

Typy linii przechwyconych w turze niosą przejrzystość kanału, nie diagnostykę
błędu. Tryb „Pełny" jest jedynym miejscem, w którym ta lista ma sens: poza nim
jest szumem nad odpowiedzią. Spisu plików w wykazie wytworów nie ma: rdzeń go
nie nadaje — wpis niesie wywołania narzędzi, a nie listę tego, co po nich
zostało na dysku. Wytworem tury są więc nazwy narzędzi, które tura uruchomiła.

Błąd tury nie podlega trybowi widoku transkryptu. Wpis, w którym tura padła,
mówi o tym zawsze — schowanie błędu za ustawieniem widoku byłoby ciszą
w miejscu, gdzie Operator musi wiedzieć, że kanał odmówił.

## budowa/klient-poprzedni/src/strona-glowna/zalozenie-komponentu.ts
Wywołanie z rodzajem profilu asystenta kończy się odmową rdzenia, bo platforma nie ma magazynu profili
asystenta; przycisk zakładający taki profil byłby przyciskiem pewnej odmowy. Rodzaj wraca do wykazu, gdy
rdzeń dostanie magazyn. Zakładanie z formularza jest drogą drugą, obok kafla rodzaju: naciśnięcie kafla
nadal otwiera moduł, w którym komponent się buduje. Formularz służy operatorowi, który wie, czego chce,
i nie potrzebuje wchodzić do modułu po nazwę.

## budowa/klient/src/wejscie/tresci.ts
Plik przekazuje się tłumaczowi bez dostępu do kodu — nie ma tu ani znacznika, ani rozgałęzienia. Wartości w nawiasach klamrowych to miejsca na dane podstawiane w czasie działania; ich nazw nie tłumaczy się. Poza tym plikiem żaden plik drogi wejścia nie niesie łańcucha widocznego dla użytkownika — ani składnik, ani ekran, ani przebieg. Właściwość jest sprawdzalna: sprawdzian katalogu odrzuca każdy łańcuch spoza tego pliku, który niesie spację albo polski znak diakrytyczny.

## budowa/klient-poprzedni/src/okna-pomocnicze/potwierdzenie-czyszczenia.ts
Czyszczenie całej historii okna dostaje ten sam wyjątek co kasowanie sesji —
i tylko ono; usunięcie pozycji wskazanych ręcznie potwierdzenia nie ma, bo
wskazanie jest zgodą. Potwierdzenie usunięcia sesji nie nadaje się tu wprost,
bo jest związane kształtem z sesją: przyjmuje wykaz identyfikatorów i
tytułów, buduje wykaz jedna-pozycja-na-sesję i rozlicza po identyfikatorach
zdaniami odmieniającymi słowo sesja, a czyszczenie historii nie ma ani
wykazu bytów, ani rozliczenia po identyfikatorach — rdzeń oddaje samą
liczbę usuniętych. Wygląd jest za to powielony co do znaku: ten sam natywny
dialog, ta sama rama biblioteki, ta sama rodzina klas z tego samego
arkusza, ten sam pas stanów i ta sama zasada, że odmowa rdzenia nie zamyka
modalu.

Oddaje odpowiedź rdzenia po udanej komendzie, także odpowiedź "usunięto 0",
bo to również jest odpowiedź, którą panel ma powtórzyć. Gdy odczytu jeszcze
nie było, modal mówi to wprost i podaje liczbę widoczną jako dolną granicę
— tę samą, którą niesie zapowiedź retencji.

## budowa/klient-poprzedni/src/moduly/translate/warstwy-translate.ts
Warstwa nie jest ozdobą opisu: rozstrzyga postać elementu przy wejściu do modułu. Warstwa pierwsza
jest rozwinięta i pozostaje taka; warstwy druga, trzecia i czwarta stoją zwinięte, a zapowiedź nad
nimi mówi, co jest pod spodem — zwinięte nie znaczy ukryte. Postać elementu details trzyma
przeglądarka, więc element działa klawiaturą i ma poprawną semantykę bez ani jednego nasłuchu;
druga kopia stanu w klasie CSS mogłaby się z atrybutem otwarcia wyłącznie rozminąć. Znacznik
wywołania idzie z wykazu ikonografii opracowania i stoi przy nazwie, żeby droga do elementu była
widoczna, zanim się go otworzy.

## budowa/klient-poprzedni/src/rozmowa/ulotnosc.ts
Źródłem ulotności jest profil modułu, nie wykaz kodów prowadzony w tym pliku.
Ulotność to dwa czyszczenia, nie jedno: zamknięcie okna, po którym rozmowa nie
wraca, bo historia nie jest odtwarzana z rdzenia, oraz zmiana kontekstu
roboczego, na przykład zmiana testowanego eksperta. Oba muszą być widoczne:
czat, który po cichu gubi wątek, czyta się jak awaria aplikacji. Dlatego
polityka niesie gotowe zdania, nie samą wartość logiczną, a rozmowa wypisuje
je w oknie.

Zasięg ulotności kończy się na kliencie: rdzeń dopisuje każdą wiadomość
Operatora i każdą odpowiedź modelu do dziennika rozmowy niezależnie od trybu
okna, kontrakt nie ma komendy kasowania wiadomości, a zamknięcie okna zmienia
tylko jego stan, nie usuwa wierszy. Bez pamięci sesyjnej znaczy więc: klient
nie odtwarza i nie pokazuje, a rdzeń zapis trzyma. Okno mówi o tym Operatorowi
wprost.

Kontekst roboczy jest portem ogólnym: warstwa rozmowy nie zna ani modułu
Agents, ani pojęcia eksperta. Wie tylko, że kontekst ma klucz, a zmiana klucza
kończy dotychczasową rozmowę.

Polityka modułu bierze regułę z jego profilu, żeby nie prowadzić drugiej listy
modułów bez pamięci, która rozjechałaby się z profilem przy jego zmianie.
Kontekst roboczy podaje warstwa składająca, bo tylko ona zna moduły. Moduł bez
pamięci sesyjnej i bez podanego kontekstu jest wciąż ulotny — traci rozmowę
przy zamknięciu okna — mówi tylko o jednym czyszczeniu zamiast dwóch.

## budowa/klient-poprzedni/src/moduly/roundtable/glos-biezacy.ts
Rdzeń nadaje wypowiedź dwiema drogami: zdarzeniem zmiany debaty, rodzaju stworzonej z treścią
pustą w chwili otwarcia głosu, rodzaju zaktualizowanej z treścią pełną po domknięciu strumienia,
oraz fragmentami strumienia przez cały czas mówienia. Widok czytający wyłącznie zdarzenia pokazuje
pustą wypowiedź aż do domknięcia strumienia, dlatego obie drogi trzeba złożyć w jedną treść.
Świeższą treść wybiera się po długości, bo fragment strumienia nie niesie znacznika czasu, a
znacznik czasu wypowiedzi zapisanej bywa zerem: rosnąca wygrywa, gdy jest niekrótsza; utrwalona
wchodzi, gdy rosnącej nie ma albo urwała się krótsza — strumień zerwany w połowie, a rdzeń zapisał
całość. Moduł nie tworzy węzłów interfejsu i nie zna klas stylu: trzy widoki rysują inaczej
(kolumna, wykaz, panel), a treść i stan mają mieć to samo. Nazwy stanów są napisami stałymi, bo
trafiają do atrybutu danych, w który arkusze stylu celują wprost; wyliczenie stoi tu, a nie
w widoku, żeby zmiana napisu w jednym miejscu nie rozjechała dwóch arkuszy i trzech okien.

## budowa/klient-poprzedni/src/strona-glowna/zmiana-komponentu.ts
Zakładanie stoi obok, w osobnym pliku; tu jest to, co robi się z komponentem już założonym. Zmiana pól
zmienia pola podane, a pominięte zostawia bez zmian — tak stanowi kontrakt; formularz wysyła więc wyłącznie
to, co operator naprawdę zmienił, bo przesłanie wszystkich pól przy każdej zmianie nazwy nadpisywałoby opis
i stan czynności wartościami z ekranu, który mógł być odczytany dawno. Przypisanie zapamiętuje na wierszu
komponentu, na którym poziomie zasięgu ten komponent obowiązuje, i nic ponad to: nie przenosi bytu
modułowego, nie nadaje uprawnień, nie włącza komponentu do żadnej pętli wykonania. Znaczenie przypisania
nie zostało jeszcze domknięte w rdzeniu; widok mówi to wprost, zamiast obiecywać skutek, którego rdzeń nie
wywołuje. Rdzeń przyjmuje pięć poziomów z dziewięciu: globalny, środowisko, projekt, sesja i okno; poziomy
modułu, pary modułów, roli i aplikacji kończą się odmową sprawdzenia żądania.

## budowa/klient-poprzedni/src/moduly/studio/karty-dokumentow.ts
Zakładki i podział powierzchni są równorzędnymi trybami pracy z dwoma dokumentami, a wybór między
nimi jest zapamiętany. Zakładki dają większe pole pracy nad jednym pismem, bo dwie kartki obok
siebie schodzą do rozmiaru, w którym pisma się nie czyta, podczas gdy podział pokazuje oba naraz.
Przełączenie trybu niczego nie gubi i nie zamyka: stan każdego dokumentu siedzi w jego migawce
niezależnie od trybu, a dokument niewidoczny zostaje otwarty i dostępny modelowi. Zakładka czynna
jest dokumentem czynnym modułu, więc pozostałe panele pracują na tym, co operator widzi, a
zakładka odłożona trzyma swoją migawkę: treść, zaznaczenie, propozycję i parę porównania —
przełączenie jest podmianą migawek, nie drugą kopią stanu modułu. Źródło wstawień jest
nieobowiązkowe, bo zakładki działają też bez niego, zakładając zakładkę bez dokumentu; gdy
źródło jest podane, przycisk nowego dokumentu zakłada pustą stronę gotową do pisania z domyślnym
arkuszem stylów i nastawami strony — to była jedyna droga zakładania dokumentu, której rdzeń
wcześniej nie miał czym obsłużyć.

## budowa/klient-poprzedni/src/moduly/translate/wymiana-glosariusza.ts
Czynności wydzielono z okna, bo formularz terminu opisuje jeden termin, a to są czynności zbiorcze
o jednej odpowiedzialności na plik. Ścieżka pliku jest ścieżką po stronie rdzenia: klient plików
nie czyta i nie zapisuje, więc kontrolką jest pole tekstowe, a nie okno wyboru pliku przeglądarki,
które sugerowałoby przesył nieprzewidziany kontraktem. Trzy czynności stoją poza wytwórnią
elementów: wytwórnia składa pole i pasek przycisków, a każda czynność jest osobną funkcją modułu,
mówi o czym innym i daje się sprawdzić bez klikania w przycisk.
Odpowiedź eksportu niesie wyłącznie liczbę wyeksportowanych terminów: ani ścieżki wyniku, ani
znaku, że plik powstał — odpowiedź wygląda tak samo także przy ścieżce do nieistniejącego
katalogu i przy napisie, który ścieżką nie jest. Liczba mówi o zawartości glosariusza w chwili
zlecenia, bo rdzeń liczy zastane terminy, a nie o zapisie, dlatego wiersz odpowiedzi ma wydźwięk
odmowy.

## budowa/klient-poprzedni/src/okno-komunikacji/przelacznik-srodowiska.ts
Komponent pokazuje maszynę, nie nazwę wartości kontraktu: wyliczenie niesie wartości techniczne, a w oknie stoi opisowa nazwa albo nazwa hosta z ustawienia wykonania. Wybór hosta zdalnego bez wskazanej maszyny zatrzymuje się na miejscu z pełnym zdaniem odmowy — co się nie stało, dlaczego i czym to zmienić — zamiast zapisywać wartość, za którą nie stoi żadna maszyna; wyróżnienie wybranej maszyny bierze się wyłącznie z migawki stanu, więc nieudane przełączenie nie zostawia mylącego podświetlenia. Drzewo menu ma jeden poziom, bo port zna wyłącznie polecenie zastosowania środowiska: nazwa hosta jest osobnym ustawieniem poziomu okna i zmienia się w komplecie sterowania — gałąź z nazwami hostów, z których żadnej nie dałoby się stąd wybrać, byłaby atrapą. Zależności są wąskie i wstrzykiwane: komplet sterowania okna podaje migawkę, subskrypcję i wysyłkę polecenia aktualizacji, a komponent nie zna kanału ani stanu globalnego — gotowy element montuje widok rozmowy tuż nad polem wypowiedzi, bo dopiero powłoka zna naraz okno i jego komplet sterowania.

## budowa/klient-poprzedni/src/rozmowa/wiadomosc-z-rdzenia.ts
Tu mieszka całe przełożenie wiadomości kontraktu na wpis wątku: dwie
równorzędne gałęzie — wypowiedź roli `user` i domknięcie tury modelu — wraz
z obroną przed podwojeniem wpisu. Rozdzielenie od rozmowy głównej biegnie
wzdłuż odpowiedzialności, nie wzdłuż długości pliku. Plik nie wie, kto
wypowiedź napisał, i nie zgaduje: kontrakt nie niesie sprawcy zmiany, więc
jedyne, co da się udowodnić o wypowiedzi nieznanej wątkowi, to że nie powstała
w tym połączeniu.

Autorem wypowiedzi Operatora bywa nie tylko sam Operator: prompt wpisuje też
asystent, innym połączeniem, przez MCP. Wypowiedź wchodzi więc do wątku
zawsze — inaczej pytanie, na które model odpowiada, nie byłoby widoczne ani
w polu wypowiedzi, ani w wątku. Przed podwojeniem broni jej wiązanie po
treści: echo miejscowe dostaje identyfikator z rdzenia zamiast drugiej
pozycji obok.

Wpis składamy wprost, a nie przez odczyt historii dla wiadomości: tamta droga
przejmuje wpis oczekujący, czyli ramkę przygotowaną na odpowiedź modelu,
i wypowiedź wjechałaby wtedy w miejsce odpowiedzi.

Rozstrzygnięcia stanu wpisu domkniętego zdarzeniem zmiany wiadomości są
cztery, nie dwa: tura zamknięta błędem i tura, która nie przyniosła ani
jednego znaku, mają własne stany. Bez nich obie wyglądałyby na ekranie tak
samo jak udana — pustą ramką z napisem „zakończona".

## budowa/klient/src/wejscie/wejscie-z-rdzeniem.test.ts
Sprawdzian przytacza odpowiedzi rdzenia i wymaga bazy świeżej: rejestracja wykonuje się raz, więc rdzeń z założonym już kontem odmawia jej, a wtedy nie ma czego zmierzyć. Wyjątkiem jest bieg wskazujący drogę potwierdzenia: ten kontynuuje rejestrację z biegu poprzedniego, więc konta oczekuje, zamiast go zakładać. Inaczej gałęzi z kontem nadawczym nie dałoby się domknąć — drogę niesie list, a list powstaje dopiero przy rejestracji.
Milczenie rdzenia kończy się tu niepowodzeniem nazywającym przeszkodę, nie pominięciem: sprawdzian, który sam siebie odpuszcza przy braku rdzenia, wygląda potem tak samo jak sprawdzian zdany.
Obie gałęzie mierzy ten sam bieg, uruchomiony dwa razy: raz wobec rdzenia bez konta nadawczego, raz wobec rdzenia z kontem nadawczym. Gałąź rozstrzyga odpowiedź rdzenia, nie nastawa sprawdzianu.
Bez konta nadawczego bramka wpuszcza hasłem od razu. Z kontem nadawczym bramka jest zamknięta do chwili potwierdzenia adresu — i to jest właściwa droga tej gałęzi, nie usterka. Drogę potwierdzenia niesie list, więc sprawdzian bierze ją z zewnątrz; bez niej mierzy samą odmowę.

## budowa/klient-poprzedni/src/okna-pomocnicze/rejestr-pomocniczych.ts
Spis zawiera także pozycje niezbudowane: bez nich moduł wyglądałby na
kompletny, a brak przestałby być widoczny, tak samo jak przycisk bez
komendy w kontrakcie pozostaje widoczny i mówi, czego brakuje. Stan
"nie zbudowano" nie jest tym samym co "nie da się zbudować", a "buduje to
inna praca" nie jest tym samym co "nikt nie rozstrzygnął, czy ma powstać".

Zasoby Designu mają być dostępne w Studio Editorze i w module Apps — panel
obok rozmowy nie jest oknem operacyjnym Designu (tamten ma trzy pasy i
czynność prowadzącą donikąd poza Designem), tylko gęstszym widokiem tych
samych danych czytanym tą samą warstwą wywołań.

Historia rozmowy w spisie Apps nie ma odpowiednika w Roundtable, bo tamten
moduł prowadzi debatę wielu modeli w jednym oknie i pokazuje obok rozmowy
jej przebieg, nie drugi komplet okien operacyjnych — historia dotyczy
tego samego okna i niesie zapis trwały, którego przebieg debaty nie niesie.

## budowa/klient-poprzedni/src/strona-glowna/zrodlo-sesji.ts
Rozłączenie klienta nie kończy sesji ani procesów, więc strona główna musi wiedzieć, że sesja żyje i gdzie
żyje. Odpis żywego stanu niesie środowisko, moduł i liczbę okien — sekcja nie wymyśla ani jednej liczby,
pokazuje wyłącznie to, co oddał rdzeń. Odmowa albo odpowiedź o złym kształcie daje stan błędu z treścią
odmowy, a nie pusty widok udający brak sesji.

## budowa/klient-poprzedni/src/moduly/translate/wystawienie-operacji.ts
Opracowanie nazywa tę funkcję wystawieniem tłumaczenia i kontroli jakości jako operacji
wywoływanych z zewnątrz — jedno z dwojga jest w kontrakcie naprawdę, drugiego nie ma, i element
rozdziela je wprost: operacje modułu są zadeklarowane jako narzędzia modelu z nazwą, komendą
i kompletem pól żądania, co jest interfejsem, którym wywołuje je model, i realnym wystawieniem
operacji; webhooka wywoływanego z zewnątrz kontrakt nie wystawia, bo nie ma komendy zakładającej
odbiornik ani adresu, pod który rdzeń by uderzył. Wykaz deklaracji nie jest przepisany: pochodzi
ze stałej kontraktu, więc dopisanie komendy do obszaru zmienia go samo, a liczba pozycji jest
liczona, nie wpisana.

## budowa/klient-poprzedni/src/moduly/studio/kategorie-operacji.ts
Pozycje wykazu kategorii nie są wierszami rejestru akcji: rejestr dla zasięgu modułu studia
oddaje wyłącznie komendy okna komunikacji, a operacji redakcyjnych w rejestrze nie ma, choć
komenda operacji kontekstowej ma uchwyt i wychodzi do kanału modelu okna — brakuje tylko wierszy
katalogu, które nazwałyby poszczególne operacje. Dlatego panel pokazuje przy każdej pozycji, skąd
ona jest: z rejestru rdzenia albo z tego wykazu. Pozycja wykazu jest identyfikatorem akcji
podawanym rdzeniowi, więc niczego nie udaje — naciśnięcie wychodzi do rdzenia i wraca jego
odpowiedzią albo odmową. Operacje pływaka kontekstowego mają na wierzchu stawać czynności
najczęstsze wedle rzeczywistego użycia, nie wedle domysłu, ale użycia w chwili pierwszego
uruchomienia jeszcze nie ma, a pływak bez żadnej czynności byłby pusty, więc ten wykaz jest
stanem początkowym zastępującym pierwsze użycie operatora, dopóki kolejność licząca użycie nie
przejmie sterowania. Wykaz wszystkich operacji w jednym ciągu jest jeden dla czterech dróg —
pływaka, menu pełnego, wiersza polecenia i narzędzi modelu — bo druga kopia rozjechałaby się
z pierwszą przy pierwszym dołożeniu operacji. Fraza pusta w wyszukiwaniu operacji oddaje wykaz
w całości, a nie pustkę, bo pole szukania niewypełnione nie jest zawężeniem do zera.

## budowa/klient-poprzedni/src/okno-komunikacji/przeplyw-komunikatow.ts
Prawdę o stanie okna zna rdzeń, nie znacznik lokalny: znacznik nieaktualny kosztuje jedno zbędne zatrzymanie, które zawsze odpowiada, nie utratę wiadomości. Klient wysyła zatrzymanie i wysyłkę jako parę sekwencyjną, nie równoległą — komendy nadane naraz dotarłyby w kolejności niegwarantowanej i wysłanie mogłoby wyprzedzić zatrzymanie, trafiając na okno nadal zajęte i kończąc się odmową. Nieudane zatrzymanie nie wstrzymuje wysłania: jeżeli tura zdążyła tymczasem dobiec końca sama, wysłanie przechodzi, a jeżeli nie, odmowa przychodzi z rdzenia.

## budowa/klient-poprzedni/src/moduly/translate/wyszukiwarka-funkcji.ts
Zasada jednego kliknięcia mówi, że każdy element modułu jest osiągalny jednym kliknięciem, jednym
skrótem albo jednym poleceniem, a wyszukiwarka jest tą drogą dla pozycji, które nie mają własnego
przycisku w oknie: nazywa je, mówi, na której warstwie stoją, i podaje przy każdej komendę
kontraktu albo brak. Wyszukiwarka niczego nie uruchamia i nie udaje, że uruchamia: pozycja bez
komendy nie dostaje tu przycisku, który po naciśnięciu przeprosi, tylko zdanie o tym, czego
brakuje, czytelne bez naciskania czegokolwiek. Szukanie idzie środkiem nazwy, opisu i grupy, bo
nazwy pozycji są w większości angielskie i złożone, więc szukanie wyłącznie od początku nazwy nie
trafiłoby przy słowie wpisanym z pamięci ze środka nazwy.
Pozycja z komendami wymienia je co do nazwy, pozycja bez komend mówi, czego brakuje.

## budowa/klient/src/wejscie/skladniki/baner.ts
Barwa stanu obejmuje znak i głowę, a wstęga przy lewej krawędzi niesie stan kształtem — barwa jako jedyna różnica nie wystarcza wymaganiu kontrastu. Głowa z licznikiem rozpada się na trzy części: to, co przed liczbą, sam licznik i to, co po niej — inaczej mechanika musiałaby przepisywać całe zdanie co sekundę, a wtedy czytnik ekranu ogłaszałby je od nowa.
Odsłona zwłoki dostaje wartość licznika z pomiaru dopiero po odmowie rdzenia, a węzeł musi już wtedy stać.

## budowa/klient-poprzedni/src/strona-glowna/panel-komponentu.ts
Kolejność jest kolejnością pracy — najpierw komponent powstaje, potem się go zmienia i wiąże. Panel
wysuwany z kafla wymagałby warstwy nakładek, której strefa druga dziś nie ma, więc czynności stoją
w przyborniku jako rozjazd wobec warstwy projektowej, nie ukryty. Wiązanie mówi, z czym wiąże, zanim
zwiąże: zdanie nad przyciskiem nazywa komponent i byt poziomu po nazwach, nie po identyfikatorach,
i przepisuje się przy każdej zmianie wyboru — przycisk, po którym operator dowiaduje się z odpowiedzi,
co właśnie związał, byłby przyciskiem wiążącym w ciemno. Zmiana wysyła wyłącznie pola dotknięte: pole
zostawione puste nie jedzie wcale, bo kontrakt mówi, że pola pominięte zostają bez zmian — puste pole
nazwy nie jest życzeniem pustej nazwy.

## budowa/klient-poprzedni/src/moduly/roundtable/indeks.ts
Jedno okno rozmawia tu z wieloma kanałami naraz, po tym samym rejestrze kanałów, którym jedzie
okno rozmowy. Moduł nie prowadzi więc własnej listy modeli: uczestnika zakłada się na kanale
wziętym z rejestru, a nazwę kanału składa biblioteka wspólna. Moduł pracuje w oknie, nie w sesji:
każda komenda obszaru, którą moduł wywołuje, wymaga identyfikatora okna, więc złożenie jedzie przez
mechanizm widoku z okna sesji wraz z kodem modułu — przejście pyta rdzeń o okna sesji i odracza
montaż do chwili, gdy okno tego modułu jest znane. Bez kodu przejście wzięłoby okno pierwsze
w wykazie i komendy debaty jechałyby z identyfikatorem okna cudzego modułu. Rejestr kanałów
zakładany jest tutaj: wspólny egzemplarz powstaje w warstwie sceny sesji i schodzi do okna rozmowy
przez wiązanie gniazda, ale umowa rejestru modułów daje modułowi wyłącznie kanał, więc egzemplarza
rejestru nie da się tą drogą podać. Drugi egzemplarz jest dopuszczalny, bo rejestr to czytająca
pamięć podręczna nad wykazem kanałów — katalog wyboru wspólny dla całego klienta, w którym żadne
okno nie zapisuje swojego stanu; od egzemplarza wspólnego różni się jedynie chwilą odświeżenia.
Układ idzie za warstwami widoczności zamierzenia modułu. Warstwa pierwsza stoi w pasie górnym:
Model Panels (skład debaty, osobny panel na każdy model) obok monitora przebiegu (Debate Panel) —
pytanie idzie jednocześnie do wszystkich uczestników i odpowiedzi narastają obok składu. Warstwa
druga to cztery rozszerzenia boczne — Argument Map & Analysis, Voting & Evaluation Center,
Moderator Panel i Consensus Panel — zwinięte do chwili otwarcia przyciskiem z Debate Panelu.
Zwinięcie nie jest blokadą: okno jest zbudowane, zasubskrybowane i o jedno naciśnięcie dalej. Okno
rozmowy modułu nie należy do złożenia: jest bytem sesji i składa je warstwa rozmowy. Strumień
głosów jest jeden na złożenie, nie jeden na okno. Rdzeń nadaje wypowiedź uczestnika fragment po
fragmencie i rozgłasza zdarzenie zmiany debaty rodzaju stworzonej z treścią pustą, a
zaktualizowanej z pełną dopiero po domknięciu strumienia; okno słuchające samego zdarzenia pokazuje
pustą wypowiedź przez cały czas mówienia modelu. Gdyby każde z sześciu okien założyło własną
subskrypcję i własne gromadzenie, ten sam fragment byłby przyjęty sześciokrotnie, a sześć okien
miałoby sześć osobnych obrazów jednego głosu. Subskrypcja stoi zatem tutaj, obok stanu debaty,
i schodzi do okien gotowym gromadzeniem.

Rejestr kanałów przyjmowany z zewnątrz istnieje po to, żeby poszerzenie umowy opisu modułu nie
wymagało przepisywania złożenia.

Nasłuch zmiany tury stoi przed oknami, bo słuchacze stanu są wołani w kolejności zapisania: gdyby
stał po nich, okna zdążyłyby raz narysować głosy tury poprzedniej pod nagłówkiem tury nowej.

Panel debaty nie należy do złożenia okien powyżej i nie jest przez nie stawiany: gospodarzem jest
pas okien pomocniczych modułu albo kolumna paneli sceny okien równoległych, a wpina go wytwórnia
paneli. Reeksport jest dla czytającego, nie dla wytwórni: wytwórnia sięga wprost po plik panelu
debaty, bo import tego pliku wciągnąłby do gospodarza całe złożenie okien operacyjnych
i rejestrację modułu, których gospodarz nie stawia. Wpis tutaj mówi tyle, że panel należy do tego
modułu.

## budowa/klient-poprzedni/src/moduly/studio/kolejka-cyfryzacji.ts
Kolejka wczytywania była wcześniej prowadzona w oknie i ginęła z jego odświeżeniem, bo panel stał
na starszej komendzie jednorazowego wydobycia tekstu, bez pojęcia pozycji. Rodzina komend
cyfryzacji niesie dziś kolejkę po stronie rdzenia — dołożenie materiału i odczyt wykazu pozycji
ze stanem i wynikiem — więc ten plik przestał być kolejką, a stał się jej odbiciem: pamięcią
tego, co rdzeń ostatnio oddał. Rdzeń nie prowadzi wskazania pozycji, na której pracuje operator,
bo to jest nastawa widoku; nie prowadzi też słów rozpoznanych i bloków układu odebranych przy
rozpoznaniu, bo odczyt wykazu ich nie powtarza, więc bez odłożenia tutaj poprawianie słów nie
miałoby na czym pracować; i nie prowadzi znacznika wywołania w toku, żeby dwa naciśnięcia nie
poszły naraz. Stan pozycji nie jest już liczony w oknie: przychodzi z rdzenia jako wartość
pięciostanowa — oczekuje, przetwarzanie, gotowa, ponowienie z powodu pewności poniżej progu,
odmowa — a ponowienie jest stanem osobnym, bo pozycja nie jest ani gotowa, ani odmówiona, tylko
wraca do rozpoznania.

## budowa/klient-poprzedni/src/moduly/studio/kopie-zapasu.test.ts
Sprawdziany autozapisu i kopii mierzą uczciwość zapisu: wskaźnik „zapisano" pokazany po
nieudanym zapisie byłby najgorszym możliwym błędem tego modułu, bo operator zamknąłby okno
i stracił pracę, dlatego sprawdzian mierzy zachowanie na zapisie nieudanym, nie tylko na udanym.

## budowa/klient-poprzedni/src/moduly/roundtable/metryki-odpowiedzi.ts
Cztery metryki mają miejsce przy Model Panels: długość, czas odpowiedzi, koszt tokenów i szacowaną
pewność. Zmierzyć da się dwie pierwsze: długość liczy się z treści, którą okno ma przed sobą,
a czas z dwóch znaczników kontraktu — otwarcia tury i zapisania wypowiedzi. Pewność ma dziś
w kontrakcie pole, lecz rdzeń go jeszcze nie wypełnia; licznika tokenów nie niesie żadne pole
obszaru roundtable. Okno mówi o obu wprost, zamiast wypełniać rubrykę kreską. Czas mierzy odstęp
między otwarciem tury a zapisaniem wypowiedzi, a nie czas pracy modelu: chwili rozpoczęcia
wypowiedzi kontrakt nie niesie — nazwa metryki mówi dokładnie to, co metryka liczy. Rubryka
pokazująca dwie metryki z czterech wygląda na komplet, dopóki nie powie, że kompletem nie jest,
dlatego ostatnie zdanie stoi zawsze.

## budowa/klient-poprzedni/src/okna-pomocnicze/wykaz-historii.ts
Zaznaczenie wskazuje zakres, nie potwierdza czynności: panel ma osobną
czynność na wyczyszczenie całej historii okna i osobną na usunięcie pozycji
wskazanych, a wiersz bez zaznaczenia niczego nie blokuje. Pozycje, których
nastawiona zasada nie utrzyma, są oznaczone, a nie ukryte ani wygaszone —
nastawa retencji ma pokazać skutek zasady przed jej zapisem, więc wiersz
oznaczony nadal da się przeczytać i zaznaczyć. Wykaz nie zna żadnej komendy
i nie wie, skąd pozycje przyszły: dostaje tablicę, oddaje identyfikatory.

## budowa/klient-poprzedni/src/moduly/translate/wywolanie-translate.ts
Rdzeń odpowiada na komendę, której nie obsługuje, kopertą zdarzenia nieznanej komendy. Koperta
niesie identyfikator żądania, ale nie niesie pola statusu, więc korelacja klienta nie rozpoznaje
jej jako odpowiedzi, a zwykłe wywołanie czekałoby na nią bez końca, pozostawiając okno w stanie
ładowania na zawsze. Dlatego wywołanie nasłuchuje zdarzenia nieznanej komendy równolegle
z odpowiedzią i rozstrzyga się na tym, co przyjdzie pierwsze; odmowa wraca nazwana, bo pole
żądanego typu niesie typ z rdzenia, więc okno mówi wprost, której komendy rdzeń nie zna, zamiast
pokazać pustą listę jako wynik. To nie jest druga droga do rdzenia: wysyłka idzie tym samym
kanałem wysyłki, a nazwa komendy pochodzi wyłącznie ze stałych kontraktu, więc obietnica nigdy nie
jest odrzucana. Pole błędu jest wypełnione także przy odmowie, więc widok, który zna tylko zwykły
wynik, zachowuje się poprawnie. Gdy pola identyfikatora żądania zabrakło — starszy rdzeń albo
pośrednik — zostaje dopasowanie po żądanym typie: mylna zbieżność jest mniej szkodliwa niż okno
czekające bez końca. Rzeczą, której nie ma, jest tutaj uchwyt komendy, a zdanie błędu nazywa typ
wprost, więc operator nie musi znać kodu, żeby zrozumieć odmowę.

## budowa/klient-poprzedni/src/okno-komunikacji/ster-katalogow.ts
Katalogi robocze są listą, nie pojedynczą wartością, więc pozycje wykazu są przełącznikami, a nie wyborem jednokrotnym: okno pracuje na wszystkich naraz, a zdjęcie jednego z nich nie jest wyborem innego, i każda zmiana idzie komendą aktualizacji z pełną listą po zmianie. W menu nie ma dodawania, bo nowy katalog wskazuje się wpisaniem ścieżki, a rdzeń nie ma komendy dającej wykaz katalogów do wyboru, więc gałąź „dostępne katalogi” byłaby atrapą — zamiast niej stoi stopka otwierająca kolumnę sterowania, gdzie pole ścieżki działa naprawdę. Na uchwycie stoi nazwa ostatniego odcinka ścieżki, nie cała ścieżka, ponieważ pasek ma zostać jednym rzędem, a wielokropek ucinałby ścieżkę od jedynej części, która ją rozróżnia; cała ścieżka stoi w opisie pozycji.

## budowa/klient/src/wejscie/narzedzia.ts
Nic w tym pliku nie wie o żadnym oknie — to warstwa niżej niż składniki.
Treści nie wstawia się znacznikiem: węzeł powstaje przez tworzenie elementu, a tekst przez tworzenie węzła tekstowego, więc dana z zewnątrz nie ma jak stać się znacznikiem.
Pusty znacznik znaku jest usterką zestawu, a nie sytuacją, którą ma czytać wykonawca sprawdzianu.

## budowa/klient-poprzedni/src/rozmowa/zegar-ciszy.ts
Domknięcie tury przychodzi wyłącznie zdarzeniem z rdzenia, a zerwany strumień
nie ma w kontrakcie własnego zdarzenia. Rdzeń ubity w połowie tury zostawiłby
więc wpis na stanie „wysyłanie" bez końca, dlatego warstwa rozmowy rozpoznaje
zerwanie po jedynym śladzie, jaki zostaje: po tym, że nic nie przychodzi.
Zegar niczego nie przerywa i nie zamyka tury — tura wraca do życia, gdy
nadejdzie fragment. Zegar wyłącznie nazywa to, co widać; różnicę między
długim rozumowaniem a martwym połączeniem rozstrzyga Operator, mając przed
sobą liczbę sekund.

## budowa/klient-poprzedni/src/rozmowa/zlozenie-tury.ts
Rozdzielenie rodzajów fragmentu jest sednem złożenia tury: jedna droga niesie
tekst, prowenancję, metadane konta, tok rozumowania, wywołania narzędzi
i błąd, a odbiorca rozstrzyga po rodzaju, gdzie fragment trafi. Fragment
rodzaju nieznanego nie jest odrzucany — ląduje w treści wpisu z oznaczeniem
rodzaju, żeby żadna treść strumienia nie znikła bez śladu. Wpis jest
zmieniany w miejscu: historia trzyma jedną tożsamość tury, a widok odświeża tę
samą pozycję zamiast dokładać kolejną.

Podsumowanie czyta się z każdego fragmentu, nie tylko z domykającego: ładunek
podsumowania niesie fragment kanału, a strumień domyka osobne zdarzenie
z wersją ostateczną, więc odczyt tylko z domykającego gubiłby podsumowanie
całkiem. Parser odrzuca ładunek bez własnych pól podsumowania.

Odtwarzanie bloku zapisanego w bazie idzie przez ten sam rozdzielacz
rodzajów, co fragment żywego strumienia: historia po restarcie ma składać
wpis tą samą logiką, którą składała go tura na żywo, a drugi czytnik ładunków
byłby drugą prawdą o rozdziale. Blok rodzaju tekstowego przechodzi bez śladu,
bo rdzeń tekstu w blokach nie zapisuje — taki blok mógłby być wyłącznie
powtórką, a doliczenie go zdublowałoby treść wpisu.

Zastąpienie wersją ostateczną, a nie doklejenie, jest całym powodem istnienia
tego rodzaju fragmentu: treść złożona z fragmentów jest przybliżeniem, kanał
może fragment powtórzyć po rotacji konta, a tura zapasowa nadaje własne.
Wersja ostateczna jest tą samą treścią, którą rdzeń zapisuje jako wiadomość —
po podmianie wpis na ekranie i wiersz w bazie mówią to samo.

## budowa/klient-poprzedni/src/okna-pomocnicze/wytwornia-paneli.ts
Wytwórnie wpisuje się strukturalnie, bez przejściówek: panelowi wystarczy,
że niesie element, funkcję odświeżenia i zamknięcia — nadmiarowe pola nie
przeszkadzają, a pole wymagane przechodzi w miejsce opcjonalnego, bo
TypeScript wiąże typy strukturalnie. Gdy okno rozjedzie się z umową panelu,
kompilator zatrzyma się na wpisie w mapie wytwórni. Plik nie czyta rejestru
i nie zna stanów pozycji — posiadanie wytwórni i nazwanie pozycji zbudowaną
w rejestrze to dwie różne prawdy, zestawiane osobno, a rozjazd między nimi
zgłasza pas okien pomocniczych.

## budowa/klient-poprzedni/src/moduly/translate/wyzwalacze-okien.ts
Opracowanie modułu rozstrzyga postać spoczynkową wprost: w stanie spoczynku widoczne są elementy
warstwy pierwszej oraz zwinięte wyzwalacze warstw wyższych, a okna zarządców pozostają zwinięte do
chwili wywołania i po zamknięciu znikają z przestrzeni roboczej. Wyzwalacz stoi na ekranie zawsze,
jest w pełni klikalny, niesie nazwę okna, jego warstwę i znacznik wywołania z ikonografii
opracowania, a stan otwarcia mówi atrybut rozwinięcia dostępności, nie sama barwa. Okno zamknięte
jest ukrywane, nie usuwane: usunięcie zabrałoby wraz z nim treść wpisaną do jego pól i wynik
ostatniego wywołania, a operator zamyka kolumnę, żeby zrobić miejsce, nie żeby stracić pracę.
Sama nazwa okna nie mówiłaby, co się stanie po naciśnięciu uchwytu, a stan otwarcia byłby wtedy
niesiony wyłącznie przez atrybut dostępności, którego wzrokiem się nie czyta.

## budowa/klient-poprzedni/src/moduly/roundtable/obsluga-czynnosci-moderatora.ts
Obsługa wydzielona z okna panelu moderatora na wzór kontekstu budowania z Build Output, zwartą
odpowiedzialnością: wysyłka czynności, zapis tury, warunkowy zapis składu i złożenie zdania
potwierdzenia. Potwierdzenie mówi, co zrobił rdzeń, a nie co wysłało okno: zdanie składa funkcja
biorąca odpowiedź, nie napis wnoszony razem z żądaniem — gotowe zdanie sukcesu nie odróżniłoby
zamknięcia tury od potwierdzenia tury już zamkniętej, a przy zmianie zagadnienia przemilczałoby,
że rdzeń zamknął turę bieżącą i otworzył nową. Tura zapisywana jest zawsze, skład wyłącznie wtedy,
gdy rdzeń go przysłał: pole składu jest w kontrakcie nieobowiązkowe, więc bez tego pola okno nie
ma dowodu, że zmiana składu zaszła. Odmowa nie jest jedynym niepowodzeniem: sprawdzian kształtu
oddaje niepowodzenie także wtedy, gdy rdzeń odpowiedział, ale bez pola obowiązkowego — zdanie
„rdzeń odmówił" byłoby wtedy zarzutem, którego rdzeń nie orzekł. Dlatego mówimy o czynności, która
nie doszła do skutku, a powód niesie treść błędu wprost z kontraktu.

Rdzeń nie odmawia zamknięcia tury już zamkniętej — oddaje ją bez zmiany, więc okno nie ma czym
odróżnić zamknięcia od potwierdzenia stanu zastanego. Stan sprzed wywołania też nie jest świadkiem
pewnym, bo po zamknięciu potrafi jeszcze przyjść rozgłoszenie z turą oznaczoną jako otwarta,
dlatego czasownik „zamknął" pada wyłącznie tam, gdzie stan sprzed wywołania mówił coś innego.

## budowa/klient-poprzedni/src/okna-pomocnicze/zrodlo-historii.ts
Widok panelu nie zna ani jednej nazwy komendy — dostaje cztery funkcje i
tyle. Zdarzenie zmiany historii nie jest czwartą komendą: rdzeń rozgłasza
je po każdym skasowaniu, także po tym z zasady przechowywania nadanej w
innym oknie i po przemiataniu przy starcie rdzenia, więc bez tej subskrypcji
panel pokazywałby pozycje, których w bazie już nie ma. Pole pozycji zdarzenia
jest niewymagane, bo po czyszczeniu zbiorczym żadna pojedyncza pozycja nie
istnieje i panel przeładowuje wtedy wykaz w całości. Źródło nie trzyma
stanu, nie zna okna i nie rozstrzyga, czy zasada ma sens — to należy do
panelu i do rdzenia; zdarzeń po oknie nie filtruje, bo filtr wymaga
znajomości okna gospodarza, która jest w panelu.

## budowa/klient-poprzedni/src/rozmowa/zapewnienie-kanalu.ts
Rejestr kanałów jest sterowany danymi: kanał istnieje wtedy, gdy istnieje
jego wiersz, a nie wtedy, gdy typ dopisano w kodzie. Świeża baza rdzenia nie
ma ani jednego wiersza, więc okno komunikacji nie miałoby czym rozmawiać —
warstwa dokłada wiersz komendą `channel.add`, tak samo jak zrobiłby to
Operator w panelu sterowania okna. Niepowodzenie nie przerywa niczego poza
tym wywołaniem: wynik niesie opis przeszkody, a okno pozostaje czynne.

## budowa/klient-poprzedni/src/punkty-izolacji/obszar-osie.ts
Nazwy i kolejność biorą się ze źródła rdzenia. Oś jest prostopadła do poziomu, nie jego
przedłużeniem — i to jest sedno tego obszaru. Poziom zasięgu rozstrzyga pierwszy: cały łańcuch od
okna do poziomu globalnego, zanim oś w ogóle wejdzie w grę. Dopiero wewnątrz już wybranego poziomu
oś wskazuje adresata wartości — dla kogo ona obowiązuje: konto, model czy tło platformy. Klucz
rozstrzygania jest złożony: klucz, poziom, byt poziomu, oś, byt osi, czyli jedna komórka na
przecięciu dwóch współrzędnych, a nie jeden dłuższy szczebel drabiny. Widok pokazuje obie
współrzędne osobno właśnie po to, żeby nie czytało się ich jako jednej listy do przewinięcia.
Zaślepka jest wyłącznie na działanie, nie na treść. Trzy osie i ich relacja do poziomu to treść
stała — wynika ze źródła rdzenia, nie z odpowiedzi rdzenia, więc pokazuje się zawsze w pełni.
Jedyna czynność tego obszaru, zapis wartości na wskazanej osi, nie ma pokrycia w kontrakcie: żądania
zapisu kontekstu i zapisu zakresu technicznego niosą poziom zasięgu, byt poziomu i warstwę, ale pola
osi w nich nie ma. To jedyne miejsce tego okna, gdzie brak jest po stronie kontraktu, a nie
podłączenia — stoi tu więc jawny, w pełni klikalny stan braku, nie cichy brak i nie martwy przycisk.

## budowa/klient-poprzedni/src/moduly/studio/kopie-zapasu.ts
Wskaźnik „zapisano" pokazany, gdy zapis się nie udał, byłby najgorszym możliwym błędem tego
modułu, bo operator zamknąłby okno i stracił pracę, dlatego stan zapisu niesie cztery wartości,
nie dwie, a nieudany jest jedną z nich, wraz z powodem podanym przez rdzeń i wskazaniem, że
treść leży w kopii zapasowej. Rozróżnienie szeregów wersji nie jest wymyślone przez okno: rdzeń
oddaje je wprost wraz z liczbą wersji w każdym z nich, a okno czyta to pole zamiast odgadywać
autozapis po braku etykiety. Kopia zapasowa jest zakładana przed zapisem, nie po nim, bo kopia
po zapisie nie chroniłaby od niczego — awaria zapisu zostawiłaby dokument uszkodzony bez stanu
sprzed, a kopie idą niezależnie od historii wersji, żeby przetrwały awarię procesu. Plik nie zna
elementów strony ani rdzenia: wejściem są byty kontraktu, wyjściem napisy i rozstrzygnięcia.
Zgłoszenie po nagłym zamknięciu dotyczy wyłącznie kopii najświeższej niosącej zmiany
niezapisane, bo pytanie o wiele kopii naraz byłoby pytaniem, na które operator nie ma jak
odpowiedzieć. Nazwa szeregu wersji jest szeregiem zapytania, nie cechą samego wiersza, bo wersja
pola szeregu nie niesie, więc pojedyncza wersja nie mówi, z którego szeregu pochodzi — wykaz
zawężony do jednego szeregu nazwać się daje, wykaz zbiorczy nie, i okno mówi to wprost zamiast
zgadywać po braku etykiety.

## budowa/klient-poprzedni/src/moduly/translate/wzorce-placeholderow.ts
Obie czynności liczy okno, nie rdzeń, i jest to rozstrzygnięcie wymuszone kontraktem: kontrola
jakości rdzenia zgłasza niezgodność symbolu zastępczego dopiero w gotowym panelu, a komendy
zamieniającej styl zmiennych nie ma w ogóle, więc okno liczy to, co da się policzyć z samego
tekstu, i nie przypisuje wyniku rdzeniowi. Wykaz stylów jest zamknięty i wzięty z opracowania
modułu — printf, indeks w nawiasach, pojedyncza i podwójna nazwa w nawiasach klamrowych oraz
znacznik formatu — a styl spoza wykazu nie jest rozpoznawany i okno tego nie ukrywa: wynik mówi,
ile wystąpień znalazło, a nie że znalazło wszystkie. Przekład stylu zachowuje kolejność wystąpień
i nazwy, gdy styl źródłowy je niesie; gdy styl źródłowy nazwy nie ma, a docelowy jej wymaga,
nazwą zostaje numer kolejny wystąpienia — reguła jest jedna, jawna i wypowiedziana w
sprawozdaniu, bo nazwa zmiennej nie jest czymś, co wolno zgadnąć.
Podwójne nawiasy pasują także do wzorca nazwy w pojedynczych nawiasach, więc bez pilnowania
odcinków zajętych jedno wystąpienie policzyłoby się dwa razy pod dwiema nazwami stylu.

## budowa/klient-poprzedni/src/okna-pomocnicze/zrodlo-podgladu-bash.ts
Żądanie bez identyfikatora okna jest kontraktem dopuszczone i znaczy sam
odczyt ogona — wraca wtedy wartość zapisu na żywo fałszywa, co nie jest
awarią, i okno ma to powiedzieć wprost, zamiast milczeć albo udawać podgląd
na żywo, którego nie ma. Kontrakt nie ma osobnego zdarzenia zbiorczego
wyjścia, więc rdzeń rozsyła wiersze obserwatorowi wspólnym strumieniem
fragmentów z identyfikatorem okna obserwującego i identyfikatorem
wiadomości równym identyfikatorowi procesu, dlatego druga czynność tego
źródła jest subskrypcją, a nie drugą komendą. To nie jest drugie źródło
terminala: źródło terminala niesie komendy kart i procesów, a ten plik
wyłącznie jedną komendę zbiorczego wyjścia i żadnej z tamtych nie powiela.

## budowa/klient-poprzedni/src/strona-glowna/panel-komponentu.test.ts
Sprawdzian pilnuje czterech rzeczy stanowiących o odbiorze: że obie czynności mają drogę z okna, że
zmiana wysyła wyłącznie pola dotknięte, bo pominięte zostają w rdzeniu bez zmian, że okno mówi, z czym
wiąże, zanim zwiąże, i że powtórzone przypisanie nie udaje czynności, której rdzeń nie wykonał.

## budowa/klient-poprzedni/src/okno-komunikacji/ster-mikrofonu.ts
Komponent ma dwa uchwyty w jednym pudełku: przycisk nagrywania i uchwyt drzewa. Nagrywanie jest czynnością, a wybór mikrofonu nastawą; gest przytrzymania zajmuje naciśnięcie przycisku, a uchwyt menu otwiera się kliknięciem — jeden przycisk pełniący obie role zaczynałby nagranie przy każdym otwarciu wykazu. Drzewo ma dwa poziomy: gałąź „Mikrofon” niesie wykaz urządzeń, a przełącznik „Przytrzymaj, aby nagrać” stoi poziom wyżej, bo dotyczy gestu, nie sprzętu. Mikrofonu bez silnika nie stawia się wcale — sprawdzenie dostępności jest pytaniem zadanym zanim cokolwiek powstanie, stąd obietnica w wyniku i wartość pusta zamiast steru; wyszarzona ikona albo ikona odmawiająca po naciśnięciu obiecywałaby zdolność, której nie ma, a krótszy pasek niczego nie obiecuje.

## budowa/klient-poprzedni/src/moduly/roundtable/okno-argument-map.ts
Okno pokazuje strukturę zapisu debaty i nie udaje grafu argumentów. Graf jest dziś w kontrakcie
wykonalny: wypowiedź niesie pole odpowiedzi na inną wypowiedź oraz typ aktu mowy, a węzły
i krawędzie oddaje komenda wykazu argumentów. Niezbudowana jest obsługa — rdzeń pól relacji jeszcze
nie wypełnia, a to okno komendy grafu nie wywołuje. Do tego czasu okno liczy chronologię i zbieżność
leksykalną wypowiedzi, mierzy pokrycie pól relacji i mówi wprost, czego brakuje: obsługi, nie
kontraktu. Nazwy widoków mówią, co liczą, żeby „macierz zbieżności" nie czytała się jako macierz
zgodności stanowisk. Wyliczenia siedzą w pliku struktury argumentów; tutaj zostaje wywołanie,
rysowanie i eksport. Subskrypcji strumienia okno nie zakłada — jedna na całe złożenie stoi
w indeksie modułu, a przyrost przychodzi wywołaniem odświeżenia głosów. Okno nie wywołuje dziś ani
jednej komendy obszaru: wszystko, co pokazuje, pochodzi ze stanu debaty wspólnego oknom modułu.

Pozycje bez obsługi stoją widoczne, bo okno bez nich wyglądałoby na analizę kompletną. Każda
nazywa po naciśnięciu komendę, która ją wykona, gdy powstanie jej obsługa — nie mówi już, że
kontrakt jej nie przewiduje, bo przewiduje.

## budowa/klient-poprzedni/src/moduly/studio/liczniki-dokumentu.ts
Liczniki dokumentu stają w pasku statusu, a rdzeń żadnej komendy liczącej nie niesie i nie musi:
treść stoi w buforze edytora, więc liczenie jej po stronie klienta nie jest obejściem braku,
tylko właściwym miejscem tej czynności — wywołanie rdzenia po liczbę słów byłoby przesyłaniem
dokumentu po odpowiedź, którą klient ma natychmiast. Plik nie zna elementów strony, dzięki czemu
rachunek sprawdza się bez stawiania widoku. Czas czytania liczony jest tempem dwustu słów na
minutę, wartością przyjętą w typografii użytkowej dla tekstu ciągłego, nazwaną stałą, żeby przy
zmianie wymagania nie trzeba było szukać liczby magicznej w kodzie. Kropka rozdzielająca skrót
kończy zdanie tylko wtedy, gdy po niej idzie odstęp i wielka litera albo koniec treści — rachunek
jest przybliżony i takim ma pozostać, bo pełna segmentacja zdań wymaga słownika skrótów, którego
moduł nie ma i którego dla licznika w pasku statusu nie warto zakładać.

## budowa/klient-poprzedni/src/strona-glowna/okno-strona-glowna.ts
Strefy stoją w kolejności malejącej masy wizualnej: karty środowisk jako środek ciężkości z krojem
nagłówkowym, kafle komponentów tą samą kartą krojem bazowym półgrubym, kafle modułów spoza nawigacji tą
samą formą — strefa ukryta, dopóki rdzeń nie poda modułu stojącego poza nawigacją — listwa ustawień
o masie najniższej, i sesje w tle jako wykaz z danych rdzenia zasilany z zewnątrz przez warstwę wpięcia,
bo sama strona danych nie pobiera. Strona zgłasza wybór środowiska, a skutek należy do warstwy, która
stronę zamontowała — przejście przez stronę główną ma być świadome i widoczne. Wykaz środowisk aktualizuje
się, którąkolwiek drogą wykaz przyszedł, tak by zasilenie kart bez zasilenia wykazu nie dało się tu
napisać.

## budowa/klient-poprzedni/src/rozmowa/zrodlo-wykazu-ukosnika.ts
Typ funkcyjny oddający pozycje i powód jest w tym repozytorium wzorcem
zastanym: panel akcji modułu bierze pozycje tak samo i z tego samego powodu —
wykaz pusty bez powodu jest atrapą, wykaz pusty z powodem jest stanem
opisanym. Pozycją jest wprost typ z kontraktu, a nie własny kształt: drugi
model dawałby dwie prawdy o tym, czym jest pozycja wykazu.

Wybór pozycji rodzaju narzędzia oddaje zdanie dla Operatora i to, czy
dołożenie weszło. Zdanie idzie do wątku rozmowy: model dostał narzędzie,
którego nie miał, więc musi to być widoczne — tak samo jak niepowodzenie, bo
cisza po wyborze byłaby nieodróżnialna od powodzenia.

Wykaz idzie w komplecie, bez zawężania po stronie rdzenia, bo filtrowanie ma
działać od pierwszego znaku, a runda do rdzenia na każde uderzenie w klawisz
tego nie daje. Zawężanie, porządek według trafności i wytłuszczenie trafień
niesie mechanizm z biblioteki, któremu wykaz jest podawany w całości.
Identyfikator sesji idzie, gdy sesja stoi — wtedy rdzeń oznacza pozycje już
dołożone i Operator nie dokłada po raz drugi tego, co ma.

Powtórzenie dołożenia nie jest błędem: kontrakt oddaje informację
o wcześniejszym dołożeniu i to pole jest czytane wprost, więc sięgnięcie po
narzędzie już dołożone daje zdanie o tym, a nie drugi wpis. Sprawca dołożenia
przez komendę po ukośniku zapisuje się w danych wprost, choć kontrakt
przyjmuje tę wartość także milczeniem.

Zestaw narzędzi sesji poszerza także asystent działający za Operatora, nie
tylko wybór z tego okna — bez nasłuchu zdarzenia takie dołożenie byłoby
niewidoczne. Nazwy dołożeń zameldowanych przez odpowiedź komendy z tego okna
stoją w małym zbiorze sesyjnym, bo tylko ten plik widzi obie drogi wejścia tej
samej wiadomości: odpowiedź komendy i rozgłoszone zdarzenie. Bez niego
Operator, który dołożył narzędzie sam, przeczytałby o tym w wątku dwa razy.

## budowa/klient-poprzedni/src/punkty-izolacji/obszar-poziomy.ts
Drabina poziomów zapisu, w kolejności rozstrzygania od okna do globalnego — wygrywa pierwszy poziom
mający własny zapis, zanim w grę wchodzi oś. Drabina nie jest wykazem do czytania: wskazanie poziomu
ustawia zasięg czynny okna, a na nim czyta i pisze macierz izolacji, przypisanie profilu i podgląd
polityki efektywnej. Ograniczenie konfiguracji do jednego, z góry narzuconego poziomu byłoby twardą
regułą — a tych okno nie stawia. Drabina woła rdzeń i nie jest drugą kopią słownika. Nazwa, kolejność
i flaga najwęższego poziomu przychodzą z odpowiedzi rdzenia, który czyta je z tabeli poziomów zasięgu
— ten plik ich nie zgaduje ani nie trzyma jako stałej. Pole kolejności w kontrakcie liczy odwrotnie
niż kolejność rozstrzygania: rośnie od poziomu najszerszego, a tabela pokazuje kolejność od
najwęższego — stąd sortowanie malejąco. Czego rdzeń nie niesie: pola opisu, bo adapter pomija je tam,
gdzie baza nie ma treści. Zdania objaśniające przy każdym poziomie zostają więc lokalnym słownikiem
tego pliku — to nie jest druga kopia drabiny, tylko treść, której rdzeń nie ma i nie udaje, że ma.
Obszar nie pokazuje, który poziom wygrywa dla bieżącego Operatora — pokazuje wyłącznie ten, który
Operator wskazał do zapisu. Wykaz poziomów zwraca słownik poziomów, nazwę, kolejność, czy poziom jest
najwęższy w ogóle, a nie zapis Operatora na którymkolwiek z nich. Takie rozstrzygnięcie wymaga innej
rodziny komend oraz pełnej ścieżki bytów — identyfikatorów roli, sesji, projektu i pary modułów —
których klient nie zna. Bliźniacze okno Konfiguracji ma tę samą drabinę i z tego samego powodu jego
odczyt pobiera tylko dwa poziomy: globalny oraz ten wskazany punktem widzenia okna. Braku nie wolno
zamalować milczącą pustką, która wyglądałaby jak nikt nic nie zapisał zamiast uczciwym tego nie
odczytaliśmy.

## budowa/klient-poprzedni/src/punkty-izolacji/obszar-profile.ts
Pięć komend rdzenia: wykaz profili, zapis nowego albo zmiana istniejącego, wczytanie do podglądu,
przypisanie do poziomu i warstwy, usunięcie. Wczytanie to nie przypisanie i na tym polega ten obszar.
Wczytanie zwraca zawartość profilu do podglądu i do formularza; maszyneria po nim pracuje tak samo
jak przedtem. Przypisanie dopiero wiąże profil z poziomem zasięgu i warstwą i zwraca politykę
obowiązującą po przypisaniu. Dwie czynności, dwa skutki, dwa osobne zdania przy przyciskach. Cel
przypisania jest widoczny wcześniej i nie jest wybierany drugi raz tutaj: poziom zasięgu
i identyfikator bytu przychodzą z selektora zasięgu w lewym panelu, a warstwa z pasa narzędzi okna —
oba dotyczą wszystkich paneli naraz. Drugi selektor poziomu w tym obszarze pokazywałby cel przypisania
inny niż zasięg, na którym Operator właśnie przestawia macierz. Przycisk przy profilu mówi
w podpowiedzi, dokąd trafi zapis. Usunięcie idzie od razu, bez potwierdzenia; żadna pozycja nie jest
wygaszona ani zablokowana. Odmowa rdzenia wraca zdaniem trzyczęściowym: co się nie udało, dlaczego,
treścią wprost z rdzenia, i czym Operator to zmieni.

## budowa/klient-poprzedni/src/moduly/studio/linijka-pionowa.ts
Margines górny i dolny są nastawami tej samej wagi co lewy i prawy, a bez linijki pionowej
przestawia się je wyłącznie liczbą w oknie nastaw; do tego linijka pionowa pokazuje, ile miejsca
na kartce zostało, co przy pisaniu pisma na jedną stronę jest informacją, po którą operator
sięga co chwilę. Linijka pionowa nie chwyta wierszy tabeli: wysokość wiersza tabeli bierze się
z jego treści i z nastaw akapitu, a nie z chwytu na linijce, tak samo jak w pakietach biurowych —
udawanie takiego chwytu dawałoby nastawę, której nic nie pilnuje.

## budowa/klient-poprzedni/src/okno-komunikacji/ster-modelu.ts
Model rozdziela wybór na dwie grupy w jednym poziomie menu, nie w gałęzi, bo modele surowe i eksperci mieszczą się razem, a gałąź schowałaby wykaz o jeden ruch dalej. Znaczenie wyboru jest pisane raz — przedrostek eksperta, wartość zaznaczona i przekład wyboru na treść wysyłaną do rdzenia pochodzą z jednego miejsca, bo gdyby pasek składał to zlecenie po swojemu, wybór eksperta w pasku i w kolumnie sterowania znaczyłby dwie różne rzeczy. Etykieta uchwytu jest krótka: w kolumnie sterowania wiersz niesie pełny opis, a w pasku stoi sama nazwa własna, bo pasek ma zostać jednym rzędem — oba napisy składane są jednak z tych samych pól źródłowych. Kanał nieczynny zostaje na wykazie i pozostaje wybieralny, bo wyszarzenie byłoby blokadą, a o stanie kanału mówi opis pozycji.

## budowa/klient-poprzedni/src/rozmowa/wykaz-po-ukosniku.ts
Plik nie rysuje ani jednego wiersza, nie filtruje, nie szereguje i nie
wytłuszcza. Składa dane w płaską tablicę pozycji i podaje ją mechanizmowi
z biblioteki, który to wszystko już niesie. Wykaz po ukośniku jest jego
kolejną obsadą, obok sterów paska zlecenia, z jedną różnicą: pracuje w trybie
bez uchwytu, bo wyzwalaczem jest znak w polu wpisywania, a nie kliknięcie
w przycisk. Wyzwalacz, filtr, kierunek otwarcia, przewijanie, wyróżnienie
pierwszej pozycji i dopasowanie w środku napisu niesie mechanizm; do tego
pliku należy jedno: pozycje ostatnio użyte bliżej wierzchu. Rejestru sięgania
po pozycje nie ma ani w kontrakcie, ani w bibliotece, więc prowadzi go osobny
plik rejestru ostatnio użytych. Sortowanie tego pliku jest stabilne i przy
remisie oceny oddaje kolejność wejścia; świeżość podnosi więc pozycję wśród
równie trafnych i ani o krok wyżej.

Bez podpięcia słuchacza zmiany stanu wykazu pole wypowiedzi nie zauważyłoby
samo, że wykaz zwinął się bez jego udziału — mechanizm zamyka się także po
kliknięciu poza nim i po wyborze pozycji myszą. Bez tego sygnału podpowiedź
zostawałaby pod polem, przy którym nic już nie stoi.

Wybrany przy pozycji oznacza, że jest już dołożona do sesji, a nie że jest
wyróżniona — wyróżnieniem zarządza mechanizm sam. Druga nazwa wchodzi tylko
wtedy, gdy naprawdę jest inna niż nazwa pełna, żeby uniknąć powtórzenia
w nawiasie. Skutek wyboru jest różny dla dwóch rodzajów wpisu w jednym
wykazie: jedno powiększa zestaw narzędzi modelu, drugie wykonuje czynność
aplikacji, i z samej nazwy tego nie widać.

Wykonania nie udajemy: kształt żądania komendy akcji ten plik nie zna,
a wysłanie pustej treści byłoby zgadywaniem — Operator dostaje zdanie o tym,
co się nie stało, zamiast ciszy nieodróżnialnej od powodzenia. Bez
pogodzenia stanu obsady ze stanem mechanizmu obsada sterowałaby dalej
wykazem, którego na ekranie już nie ma.

## budowa/klient-poprzedni/src/ustawienia/sekcja-braku.ts
Dwie sekcje okna ustawień, konta operatora i powiadomień, mówią tę samą rzecz: czego w rdzeniu nie ma
i dlaczego kontrolki tu nie ma. Wspólna wytwórnia trzyma ten kształt w jednym miejscu, zamiast w dwóch
odmianach. Kontrolka wyłączona byłaby tu gorsza niż jej brak: stan wyłączony należy się elementowi, który
w danym miejscu nie ma sensu, a nie zapowiedzi czegoś, czego nie ma. Zaślepka mówi wkrótce — ta sekcja
mówi, co dokładnie sprawdzono i gdzie, z nazwami plików rdzenia i migracji, żeby dało się to zweryfikować;
zdanie fałszywe albo nieaktualne waży tu więcej niż jego brak, bo operator czyta je jako wiedzę o stanie
produktu.

## budowa/klient-poprzedni/src/ustawienia/sekcja-konta-modeli.ts
Rodzina komend kont modeli jest obsłużona w komplecie w warstwie modeli: panel kont, wykaz kont, źródło
kont. Drugi formularz do tych samych komend byłby drugim oknem do tych samych danych rdzenia, a te
rozjechałyby się przy pierwszej zmianie. Sekcja bez pól nie jest brakiem do uzupełnienia, tylko podziałem
odpowiedzialności. Sekcja istnieje, bo projekt okna ją wymienia: operator, który jej tu nie znajdzie,
szukałby jej dalej w tym oknie, a jedno zdanie z drogą jest krótsze niż to szukanie i uczciwsze niż
milczenie.

## budowa/klient-poprzedni/src/ustawienia/sekcja-konto.ts
Konto istnieje: rejestracja zakłada jedyną encję właściciela z loginem, adresem e-mail uwierzytelniającym,
stanem potwierdzenia i datą utworzenia. Sekcja nie ma jednak czym go pokazać, bo brakuje komendy odczytu
profilu — komenda rejestracji konto zapisuje, komenda weryfikacji oddaje sesję, ale żadna komenda nie
zwraca loginu ani adresu. Pole bez komendy odczytu byłoby puste, więc pól tu nie ma do czasu, aż rodzina
odczytu profilu wejdzie do kontraktu. Rodzina kont modeli nie jest tym, czego tu brakuje: opisuje konta
modeli, nie konto operatora, i jest obsłużona w warstwie modeli; sekcja kont modeli niżej w rejestrze
odsyła właśnie tam.

## budowa/klient-poprzedni/src/moduly/translate/zapisy-zrodla.ts
Odczyt i zapis mieszkają osobno, tak jak w sekcji dostępów: widok składa pola i przyciski, a to,
co się dzieje po naciśnięciu, jest tutaj, dzięki czemu sprawdzian może wywołać zapis bez klikania
w element, a plik widoku nie puchnie o obsługę odmów. Zdanie o wyniku mówi tylko to, co niesie
odpowiedź rdzenia: zapis tekstu żadnego przekładu nie uruchamia, oddaje panele z ich dotychczasową
treścią, a przelicza je dopiero dodanie celu albo korekta operatora; przy zapisie bez wskazanego
języka źródłowego rdzeń oddaje pusty język, nazwany, nie przemilczany. Rozpoznanie języka nie
wpisuje się do pola w ciemno: komenda rozpoznania przepuszcza odpowiedź modelu wprost do pola
języka kontraktu, więc w polu potrafi wrócić komunikat kanału zamiast oznaczenia języka, a wartość
niebędąca oznaczeniem języka wraca operatorowi dosłownie, jako odmowa, i nie nadpisuje tego, co
sam wpisał.
Rdzeń sprawdza pusty tekst dokładnie: pustego tekstu odmawia, ale samą spację czy sam znacznik
kolejności bajtów przyjmuje i zapisuje, a okno takiego źródła nie wysyła i nie przypisuje
rdzeniowi odmowy, której by nie było. Wskaźnik ładowania stoi przy tytule pasa, a przesłonięcie
pola skasowałoby operatorowi z oczu to, co właśnie pisze. Odpowiedź zapisu niesie język i podział,
które rdzeń oddał, a nie samą treść źródła. Przy polu pustym żądanie zapisu nie niesie języka,
więc nie ma zamówienia, z którym można by zestawić odpowiedź; rdzeń oddaje wskazany język wprost,
a przy żądaniu bez języka oddaje pole puste, kasując język zapisany wcześniej — obie drogi mają
w oknie własne zdanie. Przerysowanie okna nie zdejmuje fazy trwającej, więc ogłoszenie wykonane
przy zapalonym ładowaniu przeszłoby bez skutku i okno zostałoby z zapowiedzią wywołania, które
się już skończyło.
Kontrakt opisuje pole language jako rozpoznany język, a rdzeń prosi model o samą nazwę albo kod
ISO 639-1 — oznaczenie języka to najwyżej kilka słów złożonych z liter, ewentualnie z kodem
w nawiasie. Komunikat kanału mieści się zwykle w czterdziestu znakach, więc sama długość go nie
odsiewa; odsiewa go zbiór znaków i liczba słów. Klient nie orzeka, że wartość jest błędna: oddaje
ją operatorowi dosłownie i zostawia mu ocenę, nie wpisuje tylko cudzego komunikatu do pola języka
i nie melduje go jako rozpoznania.

## budowa/klient-poprzedni/src/okno-komunikacji/ster-nakladu.ts
W pasku stoi menu, a nie suwak, ponieważ etykieta komponentu ma być bieżącą wartością, a suwak pokazuje położenie i wymaga podpisu obok siebie — suwak w kolumnie sterowania zostaje, bo to ta sama nastawa w dwóch widokach, nie dwa stany. Stopień jest napisem wyliczenia, nie liczbą, a jego wartość idzie ustawieniem poziomu okna, bo komunikat aktualizacji okna nie ma dla niej osobnego pola. Napis pusty jest pełnoprawnym stopniem, oznaczającym rozstrzygnięcie przez kanał modelu, a nie brak ustawienia, więc stoi na wykazie jak każdy inny. Klucz pusty należy w mechanizmie menu do stopki, a stopień „bez wskazania” jest napisem pustym — przedrostek klucza rozdziela jedno od drugiego.

## budowa/klient-poprzedni/src/moduly/roundtable/okno-consensus-panel.ts
Odczyt stanowiska bierze identyfikator tury jako pole opcjonalne: podany zawęża do jednej tury,
pominięty obejmuje całą debatę. Okno nie zgaduje, którą chce Operator — pole „Zakres stanowiska"
rozstrzyga to wprost, a odpowiedź mówi, czego dotyczy. Wybór „cała debata" pomija pole tury
w żądaniu, nie wysyła go pustym napisem, bo to udawałoby żądanie tury. Brak tury bieżącej przy
zakresie „tura" i brak treści stanowiska po udanym odczycie to dwa różne stany puste, oba poprawne.
Odmowa odczytu jest trzecim, osobnym stanem z kodem i treścią rdzenia — puste stanowisko przy
odmowie mówiłoby nieprawdę. Ponowny odczyt idzie na zamknięciu tury, nie przy każdej zmianie: stan
debaty subskrybuje zdarzenie zmiany i budzi nasłuch przy każdym przyroście, więc druga subskrypcja
zdarzenia byłaby powieleniem — okno łapie wyłącznie przejście tury w stan zamknięty. Stanowisko
niesie dziś w kontrakcie znacznie więcej, niż to okno pokazuje: punkty zgody i sporne, zdania
odrębne, poparcie ważone oraz części zapisu decyzji. Okno rysuje treść i metrykę, bo tyle oddaje
mu rdzeń — pozostałe pola pozostają puste do czasu zbudowania ich obsługi, a nie dlatego, że
kontrakt ich nie ma.

## budowa/klient-poprzedni/src/ustawienia/sekcja-wyglad.ts
Motyw to ta sama nastawa, co na pasku górnym. Wiersz nie ma własnego stanu ani własnego zapisu: bierze
most, jedynego właściciela tej nastawy po stronie klienta, i jest jego drugim sterem — pierwszym jest
przełącznik na pasku. Zmiana tutaj przestawia pasek, zmiana na pasku przestawia ten wiersz, a zmiana
w drugim oknie dolatuje do obu zdarzeniem zmiany konfiguracji. Opcje wyboru i ich etykiety przychodzą
z katalogu rdzenia, nie z tego pliku. Język interfejsu jest tu miejscem nazwanym, nie kontrolką: katalog
ustawień rdzenia nie niesie klucza języka interfejsu — jedyny klucz językowy dotyczy rozpoznawania mowy.
Kontrakt nie ma komendy zmiany języka, a klient nie ma warstwy tłumaczeń, więc przełącznik nie miałby
dokąd pójść; zamiast niego stoi zdanie mówiące, czego brakuje.

## budowa/klient-poprzedni/src/ustawienia/stan-ustawien.ts
Rejestr sekcji w warstwie sekcji, nie katalog rdzenia, więc rama nie potrzebuje odczytu ani fazy ładowania
na swoim poziomie. Jedyna rzecz, którą okno pamięta między otwarciami, to sekcja, na której operator
skończył: tak samo jak okno konfiguracji, które nie wraca do pierwszej kategorii przy każdym otwarciu.

## budowa/klient-poprzedni/src/moduly/translate/zgodnosc-odpowiedzi.ts
Zdanie o skutku buduje się z odpowiedzi, a nie z żądania: okno wie, co wysłało, i wie, co rdzeń
oddał, więc porównanie jest darmowe, a zdanie zbudowane z żądania mówiłoby prawdę tylko tak długo,
jak długo rdzeń zapisuje dokładnie to, o co go poproszono. Rozbieżność jest odmową, nie uwagą na
marginesie: gdy rdzeń oddał co innego, niż zamówiono, okno nie potwierdza czynności, mówi wprost,
co wysłało i co wróciło, a przy zgodności zdanie potwierdza normalnie, bo o braku, którego rdzeń
nie pokazał, okno nie orzeka. Znaki niewidoczne wychodzą kodem, bo przycinanie białych znaków
w przeglądarce i w rdzeniu nie są tożsame: jedna strona nie zdejmuje znaku, który zdejmuje druga,
i odwrotnie — gdyby rozbieżność padła na taki znak, dwa cytaty wyglądałyby w oknie identycznie,
a operator zobaczyłby zarzut bez różnicy, dlatego każdy znak sterujący, formatujący i odstęp inny
niż zwykła spacja wychodzi jako kod znaku.

## budowa/klient-poprzedni/src/okna-rownolegle/czynnosci-sesji-menu.ts
Sekcja jest doklejana jako druga sekcja menu nagłówka okna rozmowy; kreska
nad sekcją rysuje się po stronie menu paneli i tylko wtedy, gdy sekcja
niepusta. Kolejność wierszy jest stała: otwarcie w nowym oknie, zmiana
nazwy, widok transkryptu, archiwizacja, usunięcie — pozycja bez pokrycia w
rdzeniu nie powstaje w ogóle, więc nie ma tu wiersza wygaszonego.
Identyfikator sesji czytany jest z kanału przy każdym pytaniu, nie raz przy
montażu, bo gniazdo powstaje przed uzgodnieniem z rdzeniem i sesji wtedy
jeszcze nie ma. Skróty klawiaturowe sekcji łapie nasłuch tej sekcji i
działają tylko, dopóki menu jest rozwinięte — skrótu globalnego nie
rejestrujemy, żeby te same litery wpisane w polu wypowiedzi mogły dalej
pisać litery.

## budowa/klient-poprzedni/src/ustawienia/wiersz-nastawy.ts
Nazwa rodzajowa idzie do etykiety wiersza po lewej i do etykiety dostępności uchwytu, żeby czytnik ekranu
wiedział, czego dotyczy wartość, której nazwa sama tego nie mówi. Wybór nie stoi rozwinięty: sekcja
pokazuje po jednym wierszu na nastawę, a opcje rozwijają się dopiero pod kliknięciem — sześć sekcji
rozwiniętych naraz byłoby sześcioma płachtami, a nie oknem ustawień. Uchwyt, wykaz, opisy pozycji, haczyk
przy wybranej, zwijanie kliknięciem obok i obsługa klawiatury należą do mechanizmu menu drzewa; ten plik
obsadza ten mechanizm danymi i nie odtwarza go u siebie. Klient nie zna ani jednej wartości dopuszczalnej
z góry — gdy katalog dołoży czwarty motyw, wiersz pokaże go bez zmiany tego pliku. Katalog rdzenia ma
opcję o wartości pustej — dla motywu jest nią preferencja systemu i to trzeci pełnoprawny stan nastawy,
nie brak wyboru.

## budowa/klient-poprzedni/src/okna-rownolegle/dziennik-wpisow.ts
Zmiana roli okna przebudowuje widok okna komunikacji, bo rola jest częścią
jego opisu i widnieje w jego własnym nagłówku. Sposób dokładania fragmentu
odpowiada zachowaniu historii okna: fragment dołącza się do ostatniego
wpisu tej samej persony, a przy zmianie persony zaczyna wpis nowy — rozjazd
tych dwóch reguł dałby po przebudowie inną historię niż przed nią.

## budowa/klient-poprzedni/src/moduly/translate/zrodlo-kanalow-translate.ts
Pole kanału żądania dodania celu oraz pole kanału żądania tłumaczenia zwrotnego wskazują kanał
modelu wykonujący przekład; brak pola bierze kanał czynny okna. Rdzeń czyta wskazany kanał, a
kanał, którego nie odnajduje w wykazie, odrzuca odmową nazwaną, nie cichym zejściem na domyślny.
Zawężenie do kanałów włączonych obowiązuje, bo ster ma pokazywać to, czym da się przełożyć: kanał
wyłączony rdzeń i tak odrzuci, więc stawianie go na liście wyboru byłoby zaproszeniem do odmowy.
Rejestr, który nie dotarł, daje wykaz pusty i nie odbiera ani dodania języka, ani tłumaczenia
zwrotnego: żądanie idzie wtedy bez pola kanału, a rdzeń bierze kanał domyślny; powód nieudanego
odczytu nie jest połykany, mówi go ster przy pozycji domyślnej. Faz są cztery, bo trzy pustki są
różne: faza spoczynku znaczy brak zapytania, faza odczytu znaczy zapytanie w toku, faza gotowa
z pustym wykazem znaczy, że rdzeń nie zna ani jednego kanału czynnego, a faza błędu znaczy, że
zapytanie się nie udało — zlanie ich kazałoby zgadywać, czy czekać, czy zakładać kanał.

## budowa/klient-poprzedni/src/rozmowa/rozmowa.ts
Pominięta polityka pamięci znaczy rozmowę z pamięcią, czyli stan
dotychczasowy; podana z wyłączoną pamięcią sesyjną wstrzymuje odtworzenie
historii z rdzenia, bo w module bez pamięci sesyjnej odczyt historii
przywróciłby dokładnie ten wątek, który miał zniknąć razem z oknem.

Wpis automatyzacji stawia w wątku zdanie warstwy automatycznej — zdanie,
którego nikt wprost nie wypowiedział. Tą samą drogą idą powody wyczyszczenia
rozmowy, odmowa odczytu historii i zatrzymanie tury bez tury. Wejście jest
wystawione na zewnątrz, bo dołożenie modelowi narzędzia spoza okna też musi
być widoczne, a jedynym miejscem, w którym Operator patrzy, jest wątek
rozmowy. To nie jest druga droga wypowiedzi: zdanie nie idzie do rdzenia, nie
otwiera tury i nie ma nadawcy.

Wyczyszczenie rozmowy wymaga podanego powodu, bo wyczyszczenie bez słowa jest
dla Operatora nieodróżnialne od awarii, w której aplikacja zgubiła jego
pracę. Wejście w moduł bez pamięci sesyjnej czyści rozmowę i mówi o tym
wprost: wątek modułu poprzedniego nie ma prawa zostać na oczach Operatora
w oknie, które właśnie ogłosiło, że pamięci nie prowadzi.

Subskrypcja wypowiedzi, która weszła do okna spoza tego połączenia, niesie
tylko to, co przyszło i skąd — widok rozstrzyga, czy pokazać to w polu
wypowiedzi i jak długo. Sygnał idzie osobno od wpisu, bo wypowiedź wpisana za
Operatora jest zdarzeniem innego rodzaju niż nowa pozycja w wątku, i tylko
pierwsze ma prawo ruszyć pole Operatora.

Wysyłanie wiadomości Operatora tworzy wpisy od razu, przed jakąkolwiek
komendą — interfejs nie ma blokad. Rdzeń nie przerywa biegnącej tury przy
okazji wysłania, więc przerwanie jedzie osobną komendą warstwy nadania.
Niepowodzenie komendy zamyka turę i nic ponadto: treść błędu idzie sama, bez
doklejonego kodu, bo kod jedzie osobnym polem i dokleja go widok wpisu —
sklejenie obu pokazałoby kod dwa razy.

Wyczyszczenie rozmowy ulotnej stawia najpierw powód, potem stan, bo kolejność
jest tu treścią: widok najpierw zdejmuje wszystkie pozycje, więc gdyby na tym
poprzestać, Operator zostałby z pustą listą nieodróżnialną od okna dopiero co
otwartego. Zmiana kontekstu roboczego, dla Agents zmiana testowanego eksperta,
kończy rozmowę dokładnie tak samo, jak zamknięcie okna. Wejście w moduł
ulotny czyści to, co zostało po module poprzednim, a wyjście z niego czyści
tak samo, żeby wątek testowy nie pojechał za Operatorem do pracy na serio.

Historia okna wraca z rdzenia, nie ze strumienia: wiadomości leżą w bazie
i przeżywają restart rdzenia, więc odświeżenie okna albo powrót do sesji mają
odtworzyć wątek. Odmowa odczytu historii nie może wyglądać jak pusta
rozmowa — inaczej Operator widziałby ten sam pusty stan niezależnie od tego,
czy rozmowa naprawdę jest pusta, czy rdzeń po prostu nie oddał wiadomości.
Niepowodzenie idzie więc do historii jako wpis automatyzacji, i na ekranie
zostaje ślad zamiast ciszy. Rozmowa ulotna historii nie odtwarza. Odpowiedź
spóźniona o zmianę modułu — gdy odczyt historii poszedł jeszcze w module
z pamięcią, a zanim wrócił, okno przestawiło się na moduł bez pamięci
sesyjnej — nie wchodzi do wątku, bo przywróciłaby wątek dopiero co zdjęty na
oczach Operatora, ale pominięcie zostaje nazwane, żeby nie wyglądało jak
zgubienie odpowiedzi rdzenia.

## budowa/klient-poprzedni/src/okno-komunikacji/ster-nastawy.ts
Ster nastawy nie jest drugim mechanizmem menu: rozwijanie, haczyk, opisy, gałęzie, grupy, pole szukania i stopka należą do biblioteki menu-drzewa, a ten plik bierze ten mechanizm gotowy — istnieje, bo mechanizm drzewa nie wykonuje wyboru, tylko oddaje klucz wołającemu. Bez tej obudowy ten sam kawałek — czyszczenie zdania odmowy, sygnalizacja zajętości podczas wysyłki, wyświetlenie błędu i powrót do stanu potwierdzonego — stałby osobno w każdym sterze. Wyróżnienie idzie wyłącznie z migawki stanu: ster nie zapisuje wyboru u siebie, tylko po wysyłce woła odświeżenie u wołającego, a ten czyta stan potwierdzony przez rdzeń, dzięki czemu nieudana zmiana nie zostawia mylącej etykiety na uchwycie. Ster nie traci klikalności ani na czas wysyłki, ani po odmowie — sygnalizacja zajętości mówi o pracy, nie odbiera możliwości działania. Treść pusta zdania pod uchwytem chowa je z układu, więc ster bez zdania nie zostawia pustego pasa w rzędzie.

## budowa/klient-poprzedni/src/ustawienia/zrodlo-bramki.ts
Rozdział względem źródła logowania przebiega po bramce, nie po rodzinie komend: tamto źródło woła komendy
logowania, rejestracji i odświeżenia tokenu, czyli tyle, ile trzeba, żeby wejść. Tych dwóch komend tu nie
ma, bo w oknie ustawień odmawiałyby zawsze: rejestracja jest wykonalna tylko raz, a bramka jest już
założona w chwili, gdy okno da się otworzyć; logowanie jest samą bramką i po jej przejściu nie ma czego
otwierać. Metody rozpoznawania przez system operacyjny rdzeń nie ma zbudowanej i odmawia jej założenia;
przyjmuje hasło i kod PIN. Komendy odczytu wykazu metod kontrakt nie niesie: wykaz dociera zdarzeniem
zmiany bramki, które rdzeń rozsyła do wszystkich gniazd, także po czynności wykonanej w innym oknie.
Odmowę zdjęcia metody orzeka rdzeń; okno jej nie uprzedza. Drogi odzyskania hasła listem na adres e-mail
nie ma — rdzeń poczty nie wysyła. Zdarzenie zmiany bramki niesie powód zmiany i komplet metod po niej;
pole metod jest w kontrakcie opcjonalne, więc gdy go brak, słuchacz dostaje sam powód, bo podstawianie
w to miejsce wykazu poprzedniego byłoby zgadywaniem za rdzeń.

## budowa/klient-poprzedni/src/moduly/translate/zrodlo-paneli.ts
Panele różnią się tylko identyfikatorem, który wchodzi do żądania; gdyby każda instancja
zakładała własne źródło, subskrypcja zmiany tłumaczenia powstałaby wielokrotnie i tyle samo razy
przyszłaby ta sama zmiana. Odmowa jednego panelu zostaje w tym panelu: źródło niczego nie ucisza
i nie przerywa pozostałym instancjom pracy.
Rdzeń przycina wskazanie kanału i pustego traktuje jak brak, ale wysyłanie pustej wartości
nazwałoby wskazaniem coś, czego operator nie wskazał.

## budowa/klient-poprzedni/src/moduly/studio/linijka-podzialka.ts
Chwyt marginesu i zakładanie tabulatora liczą się w milimetrach, a rysują w punktach ekranu przy
skali widoku, która bywa inna niż sto procent — pomyłka o krok podziałki przesuwa margines pisma
urzędowego, więc rachunek ma dać się sprawdzić bez przeglądarki, a linijka ma go wyłącznie
rysować. Jednostka jest wyborem operatora: milimetry albo cale, a podziałka calowa nie jest
podziałką milimetrową z inną etykietą — kreski stoją co jedną ósmą cala, a przyciąganie chodzi
po jednej szesnastej cala, bo tak działa linijka w pakiecie biurowym i tak operator spodziewa się,
że tabulator wskoczy. Kreska, która przy danej skali stanęłaby bliżej niż najmniejszy
rozróżnialny odstęp od poprzedniej, jest pomijana, bo podziałka zlana w szarą wstęgę nie mówi
nic, a udaje, że mówi; napisy niosą pełne centymetry albo cale, bo to one są miarą, w której
operator myśli o marginesie. Chwyt puszczony między kreskami ma wskoczyć na kreskę, bo margines
w rodzaju 19,73 mm nie jest nastawą, którą ktoś świadomie wybiera, a zakres jest przycinany, nie
odrzucany, więc chwyt wleczony za krawędź kartki zostaje na krawędzi. Skala niedodatnia albo
nieskończona bierze wartość jeden zamiast dzielić przez zero, bo chwyt linijki nie jest miejscem,
w którym wolno oddać nieskończoność. Wcięcie pierwszego wiersza akapitu jest liczone względem
wcięcia lewego, nie od krawędzi pola, dzięki czemu przeciągnięcie wcięcia lewego zabiera pierwszy
wiersz ze sobą, co jest zachowaniem, którego operator się spodziewa.

## budowa/klient-poprzedni/src/rozmowa/widok-rozmowy.ts
Wykaz po ukośniku czyta rdzeń i dokłada narzędzia do sesji, więc bez kanału
komend nie ma go z czego złożyć; stanowisko podglądu rozmowy kanału komend
nie prowadzi i ma dalej działać bez wykazu. Stery paska zlecenia potrzebują
kompletu sterowania okna, który powstaje dopiero w powłoce produktu, więc pasek
też jest opcjonalny z tego samego powodu. Widok przepuszcza element paska
nietknięty: nie zna ani nastaw, ani kontraktu.

Widok rozmowy jednego okna wyłącznie składa: powierzchnia modułowa na górze,
lista wpisów pośrodku, pole wysyłki na dole. Nie zna kontraktu, nie buduje
koperty i nie wie, czym jest fragment strumienia. Historia nie znika przy
zmianie modułu: lista wpisów powstaje raz, a przestawienie modułu sięga
wyłącznie do powierzchni modułowej — wygląd, możliwości, narzędzia i kontekst
się zmieniają, a wątek zostaje zachowany.

Wykaz po ukośniku powstaje w tym pliku, bo tu spotykają się jego trzy
potrzeby: źródło pozycji i droga dołożenia z powłoki oraz wątek rozmowy,
w którym ma stanąć zdanie o dołożeniu. Pole wypowiedzi dostaje go gotowego
i steruje nim ruchami klawiatury, nie znając kontraktu.

Pasek zlecenia nie wchodzi do powierzchni modułowej z trzech powodów: jego
miejsce jest przy polu wypowiedzi, nie nad zapisem; powierzchnia ma
ograniczoną wysokość i własne przewijanie, w którym ster mógłby odjechać poza
widok; a powierzchnia jest tym, co okno przestawia przy zmianie modułu,
podczas gdy koperta zlecenia należy do okna, nie do modułu.

Wyczyszczenie rozmowy ulotnej ogłasza zdjęcie pozycji osobno od wpisów, bo
„nie ma już tych wpisów" jest zdarzeniem innego rodzaju niż „jest nowy wpis";
zaraz po nim przychodzą zdania o powodzie, więc lista nie zostaje pusta dłużej,
niż trzeba.

Prompt wpisany spoza tego połączenia idzie przez pole wypowiedzi, żeby było
widać, skąd się wziął, zamiast wejść do wątku jako zastany fakt. Wpis modelu
wchodzi do wątku równocześnie i bez zwłoki, bo odwrócenie kolejności
postawiłoby odpowiedź przed pytaniem, na które odpowiada.

## budowa/klient-poprzedni/src/punkty-izolacji/obszary.ts
Podział na obszary idzie za maszynerią rdzenia, która rozstrzyga punkt izolacji na dwóch
prostopadłych osiach: oś rozstrzygania — konto, model, platforma — dla czego wartość obowiązuje;
poziom zasięgu — okno, rola, sesja, projekt, para modułów, moduł, środowisko, globalny — jak wąsko
wartość obowiązuje. Poziom rozstrzyga pierwszy, oś dopiero w jego ramach. Na te dwie osie nakłada się
przedmiot izolacji — jedenaście kluczy w dwóch grupach: kontekst, trzy klucze, wartość odrębna albo
współdzielona, domyślnie odrębna: historia wymiany, pamięć długoterminowa, bieżący stan roboczy —
pliki, projekt, załączniki, zmienne; zakres techniczny, osiem kluczy, wartość włączony albo
wyłączony, domyślnie wyłączony — stan wyjściowy platformy to pełna swoboda operacyjna: katalog
roboczy sesji, środowisko procesu, katalog danych modelu, dostęp sieciowy, odczyt i zapis plików,
konto i token, model procesu, serwer wykonania. Egzekutor po stronie rdzenia bierze wartość włączoną
i odrzuca wykonanie, które by ją naruszyło; wartość wyłączona niczego nie ogranicza. Polityka
efektywna składa te ustawienia w jeden rozstrzygnięty podgląd na dany kontekst: dla każdego klucza
wynik i poziom, z którego pochodzi. Osobnym pojęciem maszynerii jest profil: nazwany zestaw trzech
przełączników kontekstu i ośmiu zakresów technicznych, niezależny od przypisania — ten sam profil
bywa przypisany kilku sesjom, rolom, projektom i oknom. Stąd sześć obszarów tego okna, rozłożonych na
trzy panele: panel selektora zasięgu, lewy — poziom zasięgu, na którym reguła obowiązuje, oraz oś
rozstrzygania w ramach poziomu; panel macierzy izolacji, środkowy — trzy punkty izolacji kontekstu
oraz osiem punktów izolacji zakresu technicznego; panel profilu i podglądu, prawy — nazwane zestawy
jedenastu punktów z zapisem, wczytaniem, przypisaniem do poziomu i usunięciem, oraz podgląd polityki
efektywnej z rozstrzygnięciem i pochodzeniem każdego z jedenastu punktów dla wybranego kontekstu.
Kontrakt wpięcia: każdy plik obszaru eksportuje dokładnie jedną funkcję wytwórczą i nic więcej
publicznego poza własnymi typami pomocniczymi. Kontrakt obszaru i zależności obszaru są zdefiniowane
tutaj — obszar ich nie definiuje ponownie. Rama okna montuje element każdego obszaru w jego panelu
i woła odświeżenie przy otwarciu okna oraz po każdej zmianie zasięgu, warstwy i po zdarzeniu rdzenia;
zamknięcie przy demontażu okna, jeśli obszar go zwrócił. Wszystkie dwanaście komend izolacji stoją
w kontrakcie współdzielonym, a rdzeń ma je zarejestrowane. Obszary wołają je naprawdę; stan „rdzeń
nie niesie" zostaje wyłącznie tam, gdzie czynności brakuje w kontrakcie — dziś jest to zapis wartości
na wskazanej osi rozstrzygania.

## budowa/klient-poprzedni/src/okna-rownolegle/etykiety-paneli.ts
Plik stoi osobno od etykiet układu, bo tamten plik mówi o scenie — liczbie
okien, rolach i pasie relacji — a panele należą do rozmowy, nie do sceny.
Stoją tu wyłącznie napisy oprawy: nazwy i przeznaczenia paneli przychodzą
ze spisu okien pomocniczych. Konwencja jest wspólna z etykietami układu —
stałe wersalikowe dla napisów stałych, funkcje dla zdań składanych, polska
odmiana rozpisana przypadkami zamiast doklejania końcówek.

## budowa/klient-poprzedni/src/moduly/wiedza/zrodlo-wiedzy.ts
Stała komendy pada wyłącznie w plikach źródłowych modułu: okno zna czynność, poszukaj albo
przebuduj wskaźnik, a nie nazwę komendy, więc zmiana nazwy w kontrakcie przerywa kompilację
w jednym pliku, a nie w każdym widoku z osobna. Obie komendy siedzą razem, choć robią co innego:
szukanie jest odczytem, a przebudowa wskaźnika — przebudową, ale bez wskaźnika odczyt nie ma czego
oddać, więc okno pokazujące wyniki potrzebuje pod ręką drogi do przebudowy, a rozdzielenie ich
dałoby oknu dwie zależności na jedną dziedzinę. Źródło oddaje pełny wynik, nie samą treść, i nie
ma tu ani jednej pustej tablicy jako wartości domyślnej: wyszukiwanie, które po odmowie rdzenia
oddaje pustą tablicę, mówi „nic nie znalazłem" zamiast „nie udało się zapytać", a to są dwa różne
zdania, z których tylko jedno jest prawdziwe. Kontrakt nie niesie zdarzenia zmiany wskaźnika
wiedzy, więc okno nie ma się na czym zawiesić i odświeża się wyłącznie na czynność operatora —
nasłuch stanąłby na zdarzeniu, które nigdy nie przychodzi.

## budowa/klient-poprzedni/src/ustawienia/zrodlo-urzadzen.ts
Rozdział względem źródła bramki przebiega po przedmiocie, nie po oknie: tamto źródło prowadzi metody
wejścia i sesję bramki, to prowadzi urządzenia. Obie rodziny spotykają się w polu identyfikatora
urządzenia, ale odpowiadają za co innego — zdjęcie metody PIN zostawia token urządzenia nietknięty,
a unieważnienie tokenu nie zdejmuje metod. Wykaz po zmianie przychodzi zdarzeniem, nie z odpowiedzi:
unieważnienie oddaje wyłącznie potwierdzenie, a pełny wykaz rdzeń rozgłasza zdarzeniem do wszystkich
połączonych urządzeń — dzięki temu unieważnienie wykonane na jednej maszynie widać natychmiast na
pozostałych ekranach, i dlatego sekcja nie odpytuje rdzenia po każdej czynności. Rdzeń nie broni
unieważnienia własnego tokenu; orzeczenie o skutku należy do niego, a ostrzeżenie do ekranu. Zdarzenie
niesie pełny wykaz po zmianie oraz opcjonalnie urządzenie, którego zmiana dotyczy; pole bywa puste, tak
przychodzi zmiana hasła, która unieważnia tokeny hurtem, i sekcja nie zgaduje za rdzeń, czego dotyczyła.

## budowa/klient-poprzedni/src/ustawienia/tozsamosc-urzadzenia.ts
Czytelnik, który wyłącznie pokazuje, czym maszyna się przedstawia, niczego nadawać nie może: samo
obejrzenie nie jest czynnością operatora na urządzeniu. Jedna funkcja z domyślnym trybem byłaby pułapką:
pomyłka w argumencie nadałaby tożsamość przy odczycie, stąd dwie nazwy i jedno miejsce, w którym stoi
klucz pamięci. Sekcja urządzeń po tę wartość nie sięga: wykaz niesie pole wskazujące maszynę pytającego,
którym rdzeń sam ją oznacza, a tożsamość z pamięci przeglądarki byłaby przy nim drugą, słabszą prawdą.
Wartość mieszka w pamięci przeglądarki, nie w powitaniu połączenia: PIN jest właściwy urządzeniu, a
tożsamość z powitania jest nadawana na czas uruchomienia — PIN założony wczoraj ma zostać PIN-em tej
samej maszyny, więc wartość musi przeżyć zamknięcie okna. Pusty wynik oddania tożsamości znaczy, że
pamięć trwała jest niedostępna — wywołanie idzie wtedy bez wymyślonego identyfikatora i to rdzeń orzeka
odmowę, zamiast klienta zgadującego za niego.

## budowa/klient-poprzedni/src/moduly/roundtable/okno-debate-panel.ts
Okno uruchamia kolejną turę debaty i pokazuje narastające argumenty; debatę inicjuje ono samo albo
Moderator Panel. Okno wysyła jedno pytanie, a rdzeń oddaje skład uczestników, który faktycznie
odpowiada — bywa on węższy od całego składu debaty, gdy uczestnik był wyciszony w poprzedniej
turze, więc panel liczy różnicę i mówi ją wprost. Na żywo znaczy z dwóch dróg, nie z jednej: rdzeń
rozgłasza wypowiedź stworzoną z treścią pustą w chwili otwarcia głosu, nadaje słowa fragmentami
strumienia i dopiero po domknięciu strumienia rozgłasza zaktualizowaną z całością — samo zdarzenie
zmiany debaty pokazywałoby więc pustą wypowiedź przez cały czas mówienia modelu. Okno bierze
strumień wypowiedzi parametrem; subskrypcji nie zakłada, jedna na całe złożenie stoi w indeksie
modułu, bo cztery subskrypcje tego samego zdarzenia byłyby czterokrotnym odbiorem jednej treści.
Ucięcie historii wypowiedzi jest powiedziane wprost: okno bierze wypowiedzi wyłącznie ze zdarzenia,
więc otwarte w trakcie debaty widzi sam ogon, nie transkrypt od pierwszej wypowiedzi. Nie jest to
już brak kontraktu — komenda odczytu debaty oddaje skład, tury i wypowiedzi jednym wywołaniem —
lecz brak obsługi tego odczytu. Nagłówek przebiegu i eksport transkryptu to zapowiadają, żeby pusty
początek listy nie czytał się jako „debata zaczęła się teraz".

Numer i stan tury biorą się z odpowiedzi, nie z żądania okna — okno nie zna numeru tury, dopóki
rdzeń go nie nada. Wykaz pusty jest przypadkiem osobnym: „tura uruchomiona" bez ani jednego
adresata byłoby potwierdzeniem czynności, która nikogo nie dotyczy.

Przyciski rozszerzeń przychodzą gotowe z pasa rozszerzeń modułu, bo to z Debate Panelu zamierzenie
każe otwierać Argument Map & Analysis i Voting & Evaluation Center. Okno nie zna ich wnętrza —
dostaje przycisk i stawia go w pasku.

## budowa/klient-poprzedni/src/okno-komunikacji/ster-rozszerzen.ts
Podział na cztery gałęzie pierwszego poziomu drzewa — serwery MCP, wtyczki, integracje API, skille — pochodzi z wyliczenia kontraktu. Katalog jest czytany przez gotowe źródło komend rozszerzeń, a ster nie wie, jak wygląda koperta; to źródło z kolei nie importuje ani jednego pliku z okna komunikacji i zna wyłącznie kontrakt oraz protokół, więc zależność nie zamyka pętli. Liść jest przełącznikiem, nie wyborem: rozszerzenia nie wykluczają się nawzajem, więc zgaszenie jednego nie jest wybraniem innego. Opis pozycji pochodzi z rdzenia i idzie do menu dosłownie, a rozszerzenie bez opisu zostaje w menu bez zdania. Wykaz jest zawężony do zainstalowanych: menu obsługuje podłączenie, czyli to, czego operator używa teraz, a rejestrację, czyli to, co ma w ogóle, prowadzi katalog rozszerzeń w module Agents — pozycja niezainstalowana nie ma czego włączać.

## budowa/klient-poprzedni/src/okna-rownolegle/plakietka-stanu.ts
Każdy stan ma odrębną ikonę, odrębną etykietę i odrębną kropkę, więc odczyt
stanu nie zależy od rozróżnienia barw.

## budowa/klient-poprzedni/src/okna-rownolegle/tresc-probna.ts
Strona podglądu układu nie łączy się z rdzeniem, więc historia okien byłaby
pusta, a rozkład sceny nieczytelny bez treści przykładowej; treść jest
jawnie oznaczona jako przykładowa, żeby nie została pomylona z prawdziwą.

## budowa/klient-poprzedni/src/ustawienia/wiersz-wymogu-logowania.ts
Wiersz nie odsyła wprost do okna konfiguracji, bo odesłanie byłoby mylące: klucz dopuszcza wyłącznie
poziom zasięgu aplikacji, a warstwa zasięgów nie niesie go w kolejności od najwęższego, więc przecięcie
w warstwie wyboru adresu wychodzi puste i gałąź zapasowa proponuje same poziomy niewłaściwe. Skutek nie
kończy się na niewygodzie: zapis nastawy sprawdza klucz wobec katalogu, ale nie sprawdza poziomu, a poziom
aplikacji jest poziomem najszerszym, więc każdy poziom węższy go przesłania — zapis na poziomie
niedozwolonym przechodzi i potrafi wyłączyć wymóg logowania mimo zapisu przeciwnego na poziomie właściwym.
Dlatego nasłuch łapie każdy poziom, nie tylko właściwy: zapis pod niedozwolonym adresem jest tym, co
wiersz ma nazwać wprost, zamiast pominąć jako nie mój poziom.

## budowa/klient-poprzedni/src/rozmowa/widok-wpisu.ts
Wpis jest siatką dwukolumnową: w pierwszej kolumnie medalion nadawcy,
w drugiej tożsamość i wszystko, co pod nią idzie. Bez medalionu pasek
tożsamości wpadłby w kolumnę awatara i został zgnieciony. Kolejność warstw
odpowiada kolejności strumienia: prowenancja stoi przed treścią, bo rdzeń
nadaje ją przed jakimkolwiek tekstem i przed startem procesu, dalej idzie
podgląd pracy modelu, potem odpowiedź, potem narzędzia, błędy i podsumowanie
tury. Rozróżnienie nadawcy niesie ikona medalionu, klasa semantyczna i
etykieta słowna — nigdy sama barwa tła.

Widok transkryptu jest filtrem nad tym wpisem, nie drugim widokiem. Cztery
tryby sterują wyłącznie tym, które warstwy są rysowane; wpis zostaje tym samym
obiektem na tej samej pozycji listy. Przełączenie nie kasuje niczego, nie
woła rdzenia i nie blokuje — schowana warstwa wraca w całości po powrocie do
trybu, który ją pokazuje.

Komplet klas wpisu łączy budowę z biblioteki, klasę semantyczną nadawcy,
stan pracy i dwie nazwy własne jako uchwyty rozpoznania. Klasa oznaczająca,
że tura jeszcze biegnie, niesie kropkę tętna przy nadawcy.

## budowa/klient-poprzedni/src/moduly/workspace/biblioteka-czynnosci.ts
Osobny plik od okna, bo to inna odpowiedzialność: okno składa kontrolki i prowadzi stany, tu leży
przebieg trzech czynności wraz z ich obsługą niepowodzenia. Wgranie pliku potrafi oddać plik
z identyfikatorem projektu, którego wykaz biblioteki dla tego projektu nie pokaże — sam zapis
pliku nie dowodzi więc wgrania do projektu, dlatego każde zdanie sukcesu bierze wartości
z odpowiedzi i sprawdza skutek w wykazie zasobów.

## budowa/klient-poprzedni/src/okno-komunikacji/ster-uprawnien.ts
Wartości trybu odpowiadają dosłownie przełącznikowi trybu uprawnień kanału głównego i pochodzą z wyliczenia kontraktu — pasek nie prowadzi własnego katalogu trybów, tak samo jak nie prowadzi go kolumna sterowania. Żadna pozycja nie jest ukryta ani wyszarzona, w tym pominięcie kontroli uprawnień: jedyną kontrolą dostępu jest uwierzytelnianie, a interfejs nie stawia blokad — skutek wyboru mówi opis pozycji, bo to jest właściwe miejsce na ostrzeżenie, nie odebranie kliknięcia.

## budowa/klient-poprzedni/src/moduly/studio/modul-studio.ts
Studio Editor, kanwa tekstowa, Preview Window i Diff/Grep Panel pracowały wcześniej nad tą samą
treścią w czterech miejscach, a przy dwóch dokumentach dawało to sześć okien; zeszły się w okno
pracy z dokumentem, gdzie treść z formatowaniem na kartce, podgląd wydania i różnica są trybami
jej widoku, a wynik modelu jest zmianą oznaczoną w miejscu — kanwa tekstowa przestała istnieć,
bo wynik operacji kontekstowej wchodzi wprost do dokumentu. Panel cyfryzacji stoi najwyżej
w układzie, bo cyfryzacja poprzedza redakcję: dopiero jego wynik daje oknu pracy treść, gdy
materiałem wejściowym jest skan. Okno pracy jest wiodące i stoi obok panelu narzędzi niosącego
pełny wykaz operacji; warsztat dokumentu i redakcja pracują na materiale wniesionym do okna,
więc stoją niżej, a repozytorium sesji zamyka układ, bo dotyczy całej sesji, nie bieżącej
czynności. Trzy źródła modułu sięgają poza obszar studia, bo obszar ten nie niesie ani
cyfryzacji, ani zamiany formatu, ani przygotowania obrazu — kontrakt dzieli się po rodzajach
czynności, nie po modułach. Stan studia jest jeden na cały moduł, więc przywrócenie wersji
w repozytorium przestawia treść edytora i podgląd naraz, a odczyty idą równolegle i nie gaszą
się nawzajem.

## budowa/klient-poprzedni/src/okna-rownolegle/figura-modulu.ts
Liczba pochodzi z profilu modułu, role z reguł domyślnych, a sufit z pliku
identyfikatorów — plik nie ma własnego zdania o żadnym module. Liczbę okien
bierze funkcja licząca okna rozmowy, nie samo pole granicy: funkcja oddaje
zero dla modułu bez rozmowy, a pole niesie dla niego jedynkę. Wykaz par
koordynator-wykonawca nie nadaje ról — reguła domyślna i tak daje przy
dwóch oknach koordynatora i wykonawcę; mówi wyłącznie, że dla wymienionego
modułu para jest jego właściwością. Moduł, dla którego role okien nie są
ustalone, w tym wykazie nie stoi.

## budowa/klient-poprzedni/src/rozmowa/widok-zapisu.ts
Klucz techniczny trybu jest oddzielony od etykiety: kod trybu nie zmienia się
z językiem interfejsu, a polska nazwa i zdanie o przeznaczeniu stoją obok
niego w wykazie. Plik nie rysuje niczego i nie zna DOM-u — rozstrzyga
wyłącznie, które warstwy wpisu mają się pokazać; rysowanie zostaje w widoku
wpisu. Tryb nie sięga do rdzenia: wszystkie cztery stoją na polach, które
wpis już niesie. Przełączenie jest filtrem nad pamięcią okna — nie woła
komendy, nie dociąga historii i niczego nie gubi.

Podsumowanie tury i konto kanału są widoczne bez przełącznika trybu, a tryb
ma zapis poszerzać albo zawężać do streszczenia, nie odbierać informacji
dostępnej wszędzie indziej. Błąd kanału nie jest szczegółem diagnostycznym,
który wolno schować za trybem: wpis, w którym tura padła, ma o tym mówić
niezależnie od ustawienia widoku.

Tryb streszczenia filtruje wpis, z którego nie da się nic streścić — bez
domkniętej tury, bez wywołań narzędzi i bez błędu — bo taki wpis zostawiłby na
ekranie samą ramkę z godziną. Trzy pozostałe tryby nie chowają niczego.
Funkcja nie usuwa wpisu z pamięci okna ani z listy; ukrycie jest odwracalne
powrotem do innego trybu.

## budowa/klient-poprzedni/src/moduly/workspace/biblioteka-projektu.ts
Okno jest odpowiednikiem Library Explorer ograniczonym do zakresu projektu, ale okno odrębne, nie
ten sam byt. Trzy komendy należą do modułu Library, nie do Workspace: wgranie, wersje i etykieta
zbiorcza idą wprost do rdzenia, który ma dla nich uchwyty; odmowa, merytoryczna albo awaryjna,
trafia do stanu błędu okna zamiast do przycisku, który milczy. Udostępnienie do modułu
zewnętrznego idzie przenoszeniem kontekstu: zaznaczone pliki jadą wraz z projektem do modułu
docelowego jedną komendą.

## budowa/klient-poprzedni/src/punkty-izolacji/okno-punktow-izolacji.ts
Dwie warstwy, żadnej drugiej ramy. Obudowa modalna jest tym samym idiomem co w oknie konfiguracji,
dostępów i modeli. Wewnątrz ciała modalu stoi panel zbudowany przez bibliotekę ramy okna — komponent
współdzielony z oknami operacyjnymi modułów, a nie drugi, przepisany od nowa. Plik zna cztery rzeczy:
obudowę modalną, wybór warstwy izolacji, podział na trzy panele i miejsce montażu obszarów w każdym
z nich. Nie zna ani jednego klucza izolacji — te dostarczają pliki obszarów przez wspólny kontrakt
obszaru izolacji. Trzy panele, nie zakładki: selektor zasięgu po lewej, macierz izolacji pośrodku,
profil i podgląd polityki efektywnej po prawej. Panele stoją obok siebie i są widoczne naraz, bo mówią
o jednej rzeczy w trzech ujęciach: gdzie reguła obowiązuje, co ustawia i co z tego wynika po
dziedziczeniu. Rozdzielone na zakładki kazałyby Operatorowi pamiętać wybór zasięgu z jednej zakładki,
przestawiając przełącznik w drugiej. Warstwa izolacji oraz zasięg rozstrzygają, który zapis czyta
i pisze każdy obszar. Oba są stanem wspólnym okna: warstwa stoi w pasie narzędzi ramy, zasięg
w panelu lewym, a ich zmiana odświeża wszystkie obszary naraz — inaczej jeden panel pokazywałby
wartości zasięgu, którego w selektorze już nie ma. Okno otwiera się natychmiast, przed jakąkolwiek
odpowiedzią rdzenia; komendy izolacji wołają dopiero obszary, każdy własnym odczytem. Rdzeń ogłasza
dwa zdarzenia zmiany profilu i zmiany polityki, i tutaj są one słuchane raz, dla całego okna, a nie
osobno w każdym obszarze. Jedna subskrypcja odświeża wszystkie sześć obszarów; osobne byłyby
rozjeżdżającymi się odczytami tej samej zmiany. Bez nich punkt przestawiony z drugiego okna albo ręką
asystenta zostawiałby w oknie wartość nieświeżą, bez śladu i bez odświeżenia. Odświeżenie nie jest
jedyną odpowiedzią: sama zmiana wartości pod palcami Operatora byłaby posunięciem niewidzialnym, więc
nad obszarem staje ślad — co się zmieniło i czego to dotyczyło. Ślad nie powie, czyja ręka. Oba
zdarzenia idą bez pól sprawcy, choć rdzeń wypełnia je w innych kopertach, więc ślad mówi „nie
wiadomo, czyja ręka" — napis pewniejszy niż dowód byłby gorszy od jego braku.

## budowa/klient-poprzedni/src/okna-rownolegle/indeks.ts
Montaż na scenie powłoki: tworzy się układ funkcją zamontowania okien
wraz z korzeniem, rdzeniem i opisem, po czym ustawia się liczbę okien,
nadaje role poszczególnym gniazdom komendą nadania roli i pokazuje
przekazanie między dwoma gniazdami komendą pokazania przekazania. Sam
układ powstaje przez osobną funkcję tworzącą i nie wymaga rdzenia — montaż
dokłada wyłącznie podpięcie pierwszego gniazda do łączności.

## budowa/klient-poprzedni/src/ustawienia/sekcje.ts
Okno ma sześć sekcji: konto, uwierzytelnianie, wygląd i język, urządzenia, powiadomienia, konta modeli.
Jedna z nich, konto, stoi jako miejsce nazwane: mówi, czego w rdzeniu nie ma, i nie niesie ani jednego
pola, przycisku czy przełącznika bez pokrycia. Ujawnianie stopniowe obowiązuje wewnątrz sekcji tak samo
jak w oknie: na nastawę przypada jeden wiersz, na wierszu uchwyt z wartością bieżącą, a wybór rozwija się
dopiero pod kliknięciem. Samo okno jest przywoływane z listwy Centrum dowodzenia. Sekcji zbudowanych na
rodzinie komend kont i tożsamości modeli tu nie ma: te same komendy niesie już warstwa modeli — rejestr
kont, katalog kategorii tożsamości, dokumenty per oś, nakładka obowiązująca — a druga rama do tych samych
danych byłaby drugą prawdą; pozycja kont modeli nie niesie ani jednego pola tamtych komend, tylko zdanie
i drogę do tamtego okna. Okno ustawień istnieje, bo pozycja stoi w katalogu okien operacyjnych rdzenia,
a uwierzytelnianie i motyw nie mają dokąd pójść — warstwa modeli niesie tożsamość modelu, nie nastawy
operatora. Stan sześciu sekcji: konto jest miejscem nazwanym, bo bramka nie zna encji konta, a kontrakt
nie ma rodziny profilu; uwierzytelnianie niesie cztery komendy po zalogowaniu z wykazem metod odświeżanym
na żywo zdarzeniem zmiany bramki; wygląd ma jeden ster motywu i jedną prawdę z rdzeniem, a język
interfejsu jest nazwany jako brak, bez kontrolki; urządzenia niosą wykaz powiązany z kontem z
unieważnieniem tokenu, odświeżany na żywo zdarzeniem zmiany urządzeń; powiadomienia niosą macierz nastaw
katalogu — przełącznik główny, siedem klas zdarzeń i kanały dostarczenia — bo model danych mówi wprost, że
zakres klas i kanał dostarczenia są ustawieniami konfiguracyjnymi, więc droga jest ta sama, co do każdej
innej nastawy platformy, a samego doręczania nie ma jeszcze czym wykonać i sekcja mówi to wprost; konta
modeli są pozycją odsyłającą do warstwy modeli. Kolumnę nawigacji osadza rama okna ustawień, gdy rejestr
niesie więcej niż jedną pozycję.

## budowa/klient-poprzedni/src/punkty-izolacji/profile-formularz.ts
Profil jest zestawem, nie skrótem. Zapis profilu przyjmuje nazwę, opis oraz oba zestawy
przełączników, a formularz podaje zawsze komplet jedenastu, nawet gdy Operator ruszył jeden: profil
ma być pełnym zestawem, a nie różnicą wobec czegoś, czego rdzeń w tym żądaniu nie widzi. Zapis nowego
i zmiana istniejącego idą jedną komendą. Puste wskazanie identyfikatora zakłada nowy profil, podane
zmienia istniejący. Formularz mówi wprost, którą z dwóch rzeczy zrobi zapis, i pozwala wrócić do
zakładania nowego jednym naciśnięciem. Formularz nie zastępuje walidacji rdzenia: pole nazwy nie
blokuje zapisu i nie wygasza przycisku, pusta nazwa jedzie do rdzenia i to rdzeń rozstrzyga, czy ją
przyjmie. Ten plik nie zna ani jednej barwy i ani jednego odstępu.

## budowa/klient-poprzedni/src/okna-rownolegle/kolumna-paneli.ts
Panel należy do rozmowy, nie do ekranu: kolumna obsługuje jedno gniazdo,
każde okno sceny ma własny stos, a kolumna nie zna ani identyfikatora
sceny, ani pozostałych kolumn. Pusty stos nie dostaje komunikatu ani
ramki, tylko atrybut ukrycia — miejsce wraca do rozmowy, z którą dzieli
szerokość gniazda. Cykl życia panelu wraz z subskrypcją rdzenia należy do
tego, kto panel powołał — zamykanie paneli przy każdym ustawieniu stosu
ubiłoby panel przeniesiony na pełny ekran. Podział szerokości między
rozmowę a kolumnę należy do sceny i jej arkusza stylów.

## budowa/klient-poprzedni/src/okno-komunikacji/wpis.ts
Przekład roli kontraktu na klasę biblioteki obejmuje cztery wartości — rolę użytkownika na klasę człowieka, rolę modelu na klasę inteligencji, a role systemu i narzędzia na wspólną klasę systemową — i jest wystawiany wraz z medalionem, którego wymaga dwukolumnowa siatka wiersza wpisu.

## budowa/klient-poprzedni/src/moduly/studio/nastawy-wizualne.ts
Dokument w kontrakcie niesie treść, tytuł, format i wersję, ale nie pole stylów i nie ma komendy,
którą styl dojechałby do rdzenia, więc nastawy z tego pliku żyją przez sesję okna i giną z jej
zamknięciem — okno mówi to operatorowi wprost, zamiast udawać zapis, którego nie ma. Trwałe jest
wyłącznie to, co jest składnią treści, takie jak styl nazwany bloku, pogrubienie, kursywa,
podkreślenie, przekreślenie, lista, tabela i podział strony, oraz nastawy strony, bo te mają
w kontrakcie osobne pole i profil wydania. Plik nie zna elementów strony: oddaje wartości
i zdania, a przypięcie ich do elementu należy do powierzchni dokumentu.

## budowa/klient-poprzedni/src/moduly/roundtable/okno-moderator-panel.ts
Wyliczenie akcji moderatora niesie dziś dziewięć wartości; okno wysyła pięć starszych —
bezpośrednia interwencja, zamknięcie tury, następne zagadnienie, wyciszenie i zdjęcie wyciszenia.
Czterech nowszych — ustawienie kolejności głosu, wstrzyknięcie głosu przeciwnego, otwarcie wątku
pobocznego, scalenie wątku pobocznego — jeszcze nie wysyła: obsługi nie zbudowano. Kolejność głosu
jedzie zatem nadal razem z akcją bezpośrednią, samodzielnie albo obok wiadomości, choć wartość
dedykowana kolejności jest już w wyliczeniu i przeniesienie na nią kolejności usunęłoby niepewność
opisaną przy przesunięciu kolejności. Skład uczestników w odpowiedzi jest polem nieobowiązkowym,
więc okno musi umieć jego brak i wtedy nie potwierdza zmiany składu; kiedy skład przyjdzie, okno
wpisuje go do stanu debaty — to jedyna droga, którą Model Panels i Consensus Panel dowiadują się
o składzie. Skład i kanały idą wyłącznie ze stanu, tego samego rejestru, którym jedzie okno
rozmowy — okno nie woła wykazu kanałów na własną rękę.

Pas czynności trzyma treść do następnego zapisu, więc bez rozpoznania tury odmowa dotycząca
jednej tury wisiałaby nad następną; ten sam wzór niesie czyszczenie odpowiedzi przy zmianie
eksperta w oknie zarządzania umiejętnościami. Odmowa czynności nie idzie stanem błędu okna, bo ten
czyści miejsce treści i zabrałby z ekranu turę i cały skład — okno pokazywałoby pustą debatę tam,
gdzie debata jest; stan błędu zostaje zarezerwowany dla nieudanego odczytu, po którym treści
naprawdę nie ma. Powód odmowy stoi obok przerysowanej treści, oznaczony atrybutem, który arkusz
maluje barwą błędu. Przyciski skrajnych wierszy przesunięcia kolejności są klikalne zawsze,
a wygaszenie ich milczałoby o powodzie, gdy ruch jest niewykonalny.

## budowa/klient-poprzedni/src/ustawienia/indeks.ts
Okno jest jedno na klienta i żyje między otwarciami — z tego samego powodu, dla którego router nie
porzuca widoku opuszczonej trasy: powtórne otwarcie ma wrócić do sekcji, na której operator skończył,
a nie zaczynać od początku, i nie gubić stanu niezapisanego formularza sekcji, których operator akurat
nie ogląda. Kanał podajemy przy pierwszym otwarciu; wywołanie z innym kanałem, po ponownym połączeniu
z rdzeniem, buduje okno na nowo, żeby komendy nie szły przez transport, którego już nie ma. Gdyby most
motywu żył wyłącznie wewnątrz okna, zmiana motywu wykonana w drugim oknie albo na drugim urządzeniu
dolatywałaby dopiero po otwarciu ustawień; dlatego podpięcie mostu woła się raz przy wiązaniu gniazda
aplikacji, a sekcja wyglądu bierze ten sam, już podpięty most.

## budowa/klient-poprzedni/src/moduly/workspace/braki-kontraktu.ts
Powód bierze się z powitania: odpowiedź powitania oddaje wykaz komend zarejestrowanych przez
rdzeń, a nie wykaz z kontraktu. Moduł pyta o niego raz przy montażu i z odpowiedzi układa zdanie
każdej nieczynnej kontrolki, więc rozróżnia brak po stronie rdzenia od braku w kontrakcie, a gdy
rdzeń komendę zarejestruje, zdanie zmienia się samo. Stanów jest więcej niż dwa: dopóki rdzeń nie
odpowiedział, kontrolka nie orzeka o braku, tylko mówi, że pytanie jest w drodze, a odmowa
powitania też nie staje się orzeczeniem o braku. Kontrolka bez pokrycia nie znika i nie udaje, że
działa — zostaje widoczna, klikalna i niesie powód wprost.
Parametr etykieta niesie napis na przycisku, czynność jest tym, czego kontrolka miała dokonać,
i wchodzi do zdania powodu, a komendy to komendy, które by tego dokonały; brak nazw znaczy, że
okno nie zna żadnej.

## budowa/klient-poprzedni/src/okna-rownolegle/montaz-ukladu.ts
Pierwsze gniazdo obejmuje okno uzgodnione z rdzeniem — to ono ma
identyfikator nadany komendą tworzenia okna i to na nim wisi przepływ
komunikatów; pozostałe gniazda pracują z własnym opisem, dopóki rdzeń nie
otworzy dla nich osobnych okien. Wskaźnik łączności w pasku górnym nie
sięga okna na pełnym ekranie ani dalszej kolumny sceny, a polecenie
wysłane po zerwanym łączu szłoby w próżnię — łącze jest jedno, więc
odpowiedź w każdym nagłówku ma być ta sama. Okno bywa już otwarte w chwili
montażu, a bywa, że dopiero powstanie, stąd oba wejścia kodu okna, nie
jedno.

## budowa/klient-poprzedni/src/moduly/studio/obiekt-panel.test.ts
Sprawdziany warsztatu obiektów pilnują czterech rzeczy: wskazanie źródła jedzie do pola, które
rdzeń dla tego źródła czyta — węzeł designu, zasób albo plik biblioteki — zamiast wsadzać
wszystko w jedno pole ścieżki, co byłoby odmową walidacji przy każdym źródle poza plikiem;
wykres jest odmową nazwaną, widoczną przed próbą i kierującą do modułu projektowania, bo rdzeń
rachunku wykresu nie ma i okno tego nie udaje; odpowiedź o nieusunięciu jest czytana jako wynik,
nie jako awaria, i nie zdejmuje obiektu z wykazu; brak tekstu zastępczego jest nazwany, bo
wydanie bez obrazów inaczej nie powie, co w dokumencie stało.

## budowa/klient-poprzedni/src/punkty-izolacji/profile-widok.ts
Przy każdym wierszu stoją trzy czynności i dwie z nich łatwo pomylić. Wczytanie pobiera profil do
podglądu i do formularza — żaden poziom zasięgu po nim nie działa inaczej. Przypisanie do poziomu
dopiero wiąże profil z wybranym poziomem i warstwą, i to po niej izolacja faktycznie się zmienia.
Obie stoją obok siebie, więc różnica jest wypisana słowami przy każdym wierszu, a nie domyślana
z kolejności przycisków. Usunięcie usuwa od razu — bez pytania o potwierdzenie, bez wygaszania, bez
uprawnień. Profil izolacji rozstrzyga, co model widzi z sąsiedniego okna, i zmienia się na żądanie;
skutek usunięcia stoi wprost przy przycisku. Ten plik nie woła rdzenia. Buduje węzły i oddaje
naciśnięcia wywołującemu; nie zna ani jednej komendy, ani jednej barwy.

## budowa/klient-poprzedni/src/okno-komunikacji/wpisy-strumienia.ts
Miejsce jest wspólne dla dwóch odbiorców: przepływ komunikatów prowadzi jedno okno sesji uzgodnione z rdzeniem przy starcie powłoki, a panel modułu prowadzi okna zakładane osobno, i obaj biorą stąd obie funkcje, więc rozpoznanie persony i pokazanie fragmentu mają jedną postać. Podpisanie wiadomości systemowej albo wyniku narzędzia kanałem modelu kazałoby plakietce kłamać przy wpisie klasy neutralnej. Rola narzędzia powinna nieść w plakietce jego nazwę, ale kontrakt tej nazwy przy wiadomości nie przenosi — do czasu, aż ją dostanie, plakietka niesie nazwę roli, a nie zmyśloną nazwę narzędzia.

## budowa/klient-poprzedni/src/okna-rownolegle/plakietka-roli.ts
Tło gniazda zostaje neutralne niezależnie od roli — kolor sygnałowy nosi
wyłącznie powierzchnia pigułki plakietki koordynatora, nigdy tło całego
gniazda.

## budowa/klient-poprzedni/src/ustawienia/okno-ustawien.ts
Warstwę tła, stos okien i zamknięcie klawiszem Escape daje przeglądarka, a nie własna nakładka; wygląd
bierze z biblioteki komponentów, plik nie zna ani jednej barwy. Rama nie niesie treści sekcji: zna
wyłącznie rejestr sekcji (kod, nazwa, ikona, fabryka) i wspólny kontrakt sekcji — co sekcja niesie
w środku, jest sprawą pliku sekcji, nie ramy. Sekcje budują się razem, przy otwarciu okna. Rejestr sekcji
jest stały, nie przychodzi katalogiem rdzenia jak kategorie konfiguracji, więc nie ma powodu budować ich
leniwie przy pierwszym kliknięciu — budowa od razu utrzymuje stan każdej sekcji, na przykład wpisany,
jeszcze niezapisany formularz, przy przełączaniu zakładek, zamiast go gubić; kod działa dla jednej sekcji
identycznie jak dla wielu, więc dołożenie drugiej pozycji do rejestru nie wymaga zmiany tego pliku. Okno
otwiera się natychmiast; każda sekcja sama rozstrzyga swój odczyt i swój stan błędu — przycisk odczytu
ponownego w stopce odświeża wyłącznie sekcję czynną, bo sekcje niewidoczne nie ciągną rdzenia w tle bez
powodu. Nawigacja do jednego miejsca jest nawigacją donikąd: kolumna z pojedynczym, zawsze czynnym
przyciskiem sugerowałaby operatorowi wybór, którego nie ma, i byłaby kłamstwem o kształcie okna. Próg
jest samoczynny, a nie ręcznym przełącznikiem — warunek czyta długość rejestru przy każdej budowie okna,
więc po dołożeniu drugiej sekcji kolumna wraca sama, bez zmiany w tym pliku.

## budowa/klient-poprzedni/src/moduly/workspace/czynnosci-planowania.ts
Czynności stoją osobno od źródła workspace z tego samego powodu, dla którego rdzeń ma osobny port
planowania: rodzina komend workspace liczy czterdzieści komend, a jedno źródło byłoby wykazem
wszystkiego, co moduł umie, zamiast wykazem jednego obszaru pracy. Okno huba ma obowiązkowy stan
błędu, więc źródło nie połyka niepowodzenia i nie zwraca pustej listy w jego miejsce — widok musi
odróżnić brak zadań projektu od nieudanego zapytania, stąd brak wartości domyślnej w postaci
pustej listy w całym pliku.
Karta bez klucza idzie na koniec kolumny, lepiej niż na jej początek, bo świeżo dołożona karta nie
ma prawa przeskoczyć kart już ułożonych. Sformułowanie „nic do zrobienia" nie jest tym samym co
„wszystko zrobione", a pasek pełny w pustym projekcie byłby meldunkiem o pracy, której nie było.

## budowa/klient-poprzedni/src/moduly/studio/obiekt-zrodlo.ts
Kształt jest zapleczem modułu projektowania, który zakłada go od razu jako węzły ścieżki, a ikona
jego biblioteką wyszukiwania — źródło obiektów przyjmuje węzeł projektu i nazwę ikony właśnie
dlatego, że studio tylko wskazuje, co osadzić, a rysuje to moduł projektowania, więc drugiego
rachunku kształtu tu nie ma i nie będzie. Rodzaj obiektu niesie też wartość „wykres", ale rdzeń
odmawia jej nazwanym powodem, bo rachunku wykresu po stronie studia nie ma — odmowa stoi w oknie
przed próbą i kieruje do modułu projektowania, bo pokazanie kontrolki kończącej się odmową
rdzenia byłoby obietnicą bez pokrycia.

## budowa/klient-poprzedni/src/moduly/roundtable/okno-voting-evaluation.ts
Okno nie prowadzi głosowania i nie udaje, że je prowadzi — ale powód zmienił się co do istoty.
Kontrakt niesie dziś pełną rodzinę głosowania wraz z metodą agregacji, progiem kworum i wynikiem,
a także rubryki, sędziów, ranking, macierz decyzyjną i kalibrację. Brakuje ich obsługi: żadnej
z tych komend to okno jeszcze nie wywołuje. Suwak, gwiazdki albo przycisk uruchomienia głosowania
postawione bez obsługi zbierałyby wybór Operatora, który nie dociera nigdzie i znika wraz
z odświeżeniem okna — czyli byłyby pozorem czynności. Okno pokazuje zamiast tego jedyną wielkość,
którą samo umie zmierzyć: udział uczestników w turze bieżącej — ile razy i jak obszernie każdy się
odezwał, kto milczy, kto jest wyciszony i jakie ma miejsce w kolejności głosu. Nazwa wielkości mówi,
czym ona jest: udziałem, nie rankingiem i nie kworum. Wyliczenia siedzą w pliku zestawienia
udziału. Subskrypcji strumienia okno nie zakłada — jedna na całe złożenie stoi w indeksie modułu.

## budowa/klient-poprzedni/src/okna-rownolegle/ikony-paneli.ts
Rozdzielenie ikony wykazu menu od ikony rzędu skrótów dałoby tej samej
pozycji dwa różne rysunki w jednym nagłówku, dlatego jedno wiązanie kodu
pozycji z ikoną obsługuje oba miejsca naraz. O tym, które pozycje wolno
otworzyć, rozstrzyga wytwórnia paneli okien pomocniczych — pozycja spoza
wykazu nie jest błędem i nie zostaje bez rysunku, tylko dostaje ikonę
karty okna, bo każda z nich jest oknem obok rozmowy.

## budowa/klient-poprzedni/src/moduly/workspace/czynnosci-wiedzy.ts
Graf tej rodziny rysuje sieć odnośników, kto na kogo wskazuje wprost. Wyszukiwanie po znaczeniu
prowadzi osobna rodzina komend wiedzy i nie ma go tutaj, żeby okno nie obiecywało podobieństwa
tam, gdzie dostaje dopasowanie po słowach.
Okno liczy nazwy odnośników przed zapisem, żeby pokazać operatorowi, do czego strona linkuje,
zanim rdzeń odpowie.

## budowa/klient-poprzedni/src/powloka/srodowiska.ts
Plik nie jest katalogiem modułów: środowiska, moduły i macierz widoczności są sterowane danymi
i przychodzą komendami wejścia do środowiska oraz wykazu modułów. Nie ma tu ani jednej nazwy modułu,
ani jednego wiersza macierzy. Zostały dwie rzeczy, których kontrakt nie niesie: ikony pozycji, bo
moduł kontraktu nie ma pola ikony — rysunek jest zasobem pakietu wizualnego, nie wierszem tabeli, więc
mapa wiąże kod modułu z nazwą ikony zestawu, a kod nieznany dostaje ikonę zastępczą zamiast pustego
miejsca; oraz sekcje panelu orkiestracji, bo dla środowiska o nawigacji orkiestracyjnej kontrakt
zwraca pusty wykaz modułów, ponieważ sekcje panelu modułami nie są. Nie ma ich skąd wziąć z rdzenia,
więc stoją tutaj — a że żadna nie ma modułu, żadna nie ma też zbudowanego widoku operacyjnego. Wykaz
kluczy sekcji stoi osobno i publicznie, bo znać go musi zarówno boczna nawigacja, jak i moduł
MultitaskingAI, żeby zbudować powierzchnię sekcji. Dwa wykazy rozjechałyby się przy pierwszej zmianie
kolejności, a rozjazd objawiłby się pozycją nawigacji bez treści. Kolejność i widoczność podlegają
konfiguracji: układ podsekcji trzyma rdzeń, a brak ustawienia znaczy wartość domyślną, a nie
niedostępność sekcji. Kod spoza mapy ikon nie jest błędem — dostaje ikonę zastępczą, bo nowy moduł to
nowy wiersz w bazie, nie zmiana kodu klienta. Moduł jako pozycja wykazu jest funkcją publiczną, bo
pozycję buduje się także poza wykazem środowiska: moduł bez wiersza widoczności w macierzy nie stoi
w bocznej nawigacji, a mimo to daje się otworzyć — kafel komponentu własnego na stronie głównej sięga
po niego wprost do katalogu wykazu modułów i składa pozycję tą samą funkcją. Drugiego przekładu
modułu na pozycję nie ma.

## budowa/klient-poprzedni/src/okna-rownolegle/identyfikatory.ts
Identyfikator gniazda jest stały przez całe życie układu — rola, model i
katalogi zmieniają się w oknie, ale gniazdo pozostaje tym samym miejscem
na scenie, dzięki czemu komendy nadania roli i pokazania przekazania
przyjmują wprost identyfikator, bez pośrednictwa indeksu tablicy. Górna
granica liczby okien wynosi cztery, bo tyle liczy obsada multitaskingu:
koordynator, dwóch wykonawców i analityk — mniejszy sufit nie mieściłby
pełnej pętli, bo analityk nie miałby gdzie stanąć obok pary, którą ocenia.

## budowa/klient-poprzedni/src/okno-komunikacji/zrodlo-srodowiska.ts
Zasięg wykonania ma w kliencie dwa widoki — listę wyboru w szufladzie ustawień i ten przełącznik nad polem wypowiedzi — a reguła protokołu nie jest pisana po raz drugi: odczyt idzie tym samym stanem sterowania, zapis tą samą zmianą okna, tą samą komendą i tym samym identyfikatorem okna, więc oba widoki przyjmują wyłącznie stan potwierdzony przez rdzeń. Nazwa maszyny zdalnej jest ustawieniem poziomu okna, nie polem samej zmiany okna, dlatego port osobno wczytuje ustawienia zasięgu okna — bez tego przełącznik pokazywałby brak wskazania przy hoście już zapisanym w bazie.

## budowa/klient-poprzedni/src/rozmowa/wpis-rozmowy.ts
Wpis modelu odpowiada jednej turze, nie jednemu fragmentowi. Wszystko, co
turę opisuje — prowenancja wywołania, tok rozumowania, wywołania narzędzi,
konto, błędy, podsumowanie — wisi przy tym samym wpisie, zamiast rozsypywać
się po historii na osobne pozycje. Dzięki temu to, co poszło do modelu, stoi
obok tego, co model odpowiedział.

Pytanie o pustą turę dotyczy wszystkich warstw wpisu, nie samego tekstu: tura,
która oddała sam tok rozumowania albo samo wywołanie narzędzia, coś
przyniosła. Pusta jest dopiero taka, po której na ekranie nie zostaje nic
prócz nagłówka.

## budowa/klient-poprzedni/src/ustawienia/most-motywu.ts
Przy podłączeniu most czyta wartość z katalogu i stosuje wynik; subskrybuje zmianę i stosuje ją, także tę
wykonaną w drugim oknie albo na drugim urządzeniu, bo rdzeń rozgłasza zdarzenie do obu połączeń, również
do tego, które zapisało; nasłuchuje zdarzenia motywu, więc przełącznik paska jest drugim sterem tej samej
nastawy, a jego kliknięcie idzie zapisem do rdzenia; wartość pusta znaczy preferencję systemu i zostaje
trzecim stanem, bo katalog rdzenia go ma — most jej nie spłaszcza do jasnego. Zapis lokalny zostaje, choć
prawdą jest rdzeń: okno rysuje się, zanim rdzeń odpowie na pierwszy odczyt, a bez niego każde uruchomienie
zaczynałoby się mignięciem motywu preferowanego przez system; most zapisuje lokalnie każdą wartość
potwierdzoną przez rdzeń — to odbicie prawdy, nie druga prawda, bo zapis lokalny nigdy nie jedzie z
powrotem do rdzenia jako nastawa. Dwa mosty na tym samym kanale byłyby dwoma właścicielami jednej nastawy
i odbiłyby sobie nawzajem każde zdarzenie, więc podpięcie mostu oddaje most już podpięty, gdy kanał się
zgadza. Klucz nastawy motywu jest jedyną rzeczą z katalogu zaszytą w tym pliku — etykiet, opcji, poziomów
i rodzaju kontrolki most nie zna, czyta je z definicji katalogu; nazwę klucza znać musi, bo most z
definicji dotyczy jednej nastawy.

## budowa/klient-poprzedni/src/moduly/studio/ocena-redaktora.ts
Każda liczba, którą panel pokazuje, jest policzona wprost z napisu, bo produkt nie ma słownika
języka, korpusu ani reguł gramatyki, którymi liczyłby ocenę pisowni tak jak pakiety biurowe.
Zamiast oceny wymyślonej panel pokazuje ocenę policzoną z miar czytelności i nazywa, z czego
ona jest, a pozycje bez pomiaru mówią wprost, że pomiaru nie ma, i wskazują, co byłoby potrzebne,
żeby był. Wskaźnik czytelności liczy się ze średniej długości zdania i średniej długości słowa,
dwóch wielkości liczonych z samego tekstu bez słownika, z rodziny wskaźników mglistości: im
dłuższe zdania i wyrazy, tym wyżej wykształcenia potrzeba, żeby tekst przeczytać bez potykania
się, a wartość nie jest oceną jakości i panel tego nie udaje. Ocena wyrażona punktami czytelności
jest odwrotnością mglistości: mglistość odpowiadająca tekstowi prasowemu daje wynik wysoki,
a mglistość zdań wielokrotnie złożonych z terminami — niski, i panel zawsze podaje przy liczbie
jej podstawę, bo liczba bez podstawy jest kopertą. Pisowni i gramatyki panel nie mierzy, bo jedno
wymaga słownika języka, drugie analizy składniowej, a rdzeń nie ma ani jednego, ani drugiego —
te wiersze zostają więc w panelu z wartością niepodaną i z powodem, żeby usunięcie ich nie
kazało panelowi wyglądać na kompletny; interpunkcję mierzy się częściowo, tylko tam, gdzie
wzorzec jest pewny.

## budowa/klient-poprzedni/src/moduly/workspace/indeks.ts
Układ wynika z ról okien: pulpit projektu jest oknem wiodącym, więc stoi w obszarze głównym jako
punkt wejścia; panel instrukcji jest pomocniczy i stoi w kolumnie obok; pamięć kontekstu,
biblioteka projektu i zarządca agentów zarządzają repozytoriami i tworzą pas pod nimi. Okno
rozmowy modułu nie należy do tego złożenia: jest bytem sesji i składa je warstwa rozmowy.
Kolejność okien odpowiada macierzy modułu: pulpit, biblioteka, pamięć, eksperci, instrukcje.
Rodzina komend workspace liczy czterdzieści komend, a jedno źródło byłoby wykazem wszystkiego, co
moduł umie, zamiast wykazem jednego obszaru pracy.
Kod okna nadaje rama, której każde z pięciu okien podaje swój kod; funkcja oznaczania tylko go
potwierdza, gdy rama go nie ustawiła, i ustawia wartość skoku ogniska wymaganą przez klawiaturę.
Kod modułu siedzi w module, nie w mapie po stronie powłoki: dodanie modułu to jeden wpis, a nie
dwa, więc nie da się dodać modułu i zapomnieć o wytwórni.
Bez pola karty panel poziomów pamięci nie miałby czego wysłać. Tą samą drogą wchodzi okno rozmowy:
moduł pyta o okna karty i bierze to, które rdzeń przypisał jemu, bo bez niego obie kontrolki
przeniesienia kontekstu odmawiają. Odczyt jest tu, a nie przy montażu, bo dopiero wczytanie widoku
niesie kartę sesji, a powłoka woła go po odpowiedzi na wejście do przestrzeni roboczej, kiedy okno
rozmowy jest już przestawione na ten moduł.

## budowa/klient-poprzedni/src/powloka/uklad-sekcji.ts
Rozmowę z rdzeniem prowadzi źródło sekcji paneli, kształt widoku buduje moduł części sekcji. Każda
funkcja tutaj jest czysta — bierze układ, oddaje nowy, niczego nie modyfikuje w miejscu. Kolejność
liczy się od jeden i jest przeliczana po każdej zmianie: układ z dziurami w numeracji rdzeń przyjmie,
ale kolejne przestawienie policzy błędnie. Wpis rdzenia o identyfikatorze, jakiego ten panel nie
buduje, odpada przy scalaniu układu — nie da się go narysować. Sekcje znane rdzeniowi idą w jego
kolejności, sekcje panelu nieobecne w odpowiedzi dochodzą na koniec, bo panel urósł od czasu zapisu,
a wpisy rdzenia bez odpowiednika w panelu odpadają, bo panel się skurczył. Numeracja wychodzi ciągła
niezależnie od tego, który z trzech przypadków zaszedł. Ruch poza wykaz oddaje układ bez zmiany,
a wywołujący porówna go z poprzednim i nie wyśle zapisu bez skutku. Zdjęta sekcja zostaje w układzie
— traci widoczność, nie miejsce. Gdyby wypadała z tablicy, przywrócenie musiałoby zgadywać, gdzie
stała.

## budowa/klient-poprzedni/src/okno-komunikacji/zrodlo-zlecenia.ts
Cztery stery paska — model, nakład rozumowania, tryb zatwierdzania i katalog roboczy — czytają jedną migawkę i piszą jedną drogą; gdyby każdy zakładał własny stan sterowania, pasek miałby cztery kopie tej samej prawdy i cztery komplety subskrypcji na okno, a rozjazd między nimi byłby kwestią czasu, nie możliwości. Drogi zapisu są dwie, bo kontrakt ma dwie: model, tryb uprawnień i katalogi robocze mają pola w treści zmiany okna, a nakład rozumowania takiego pola nie ma i idzie ustawieniem poziomu okna — port oddaje obie drogi osobno, zamiast zlewać je w jedną, bo zlanie ukryłoby przed wołającym fakt, że to dwie różne komendy o dwóch różnych potwierdzeniach. Te same nastawy stoją w kolumnie sterowania, ale reguła protokołu nie jest pisana po raz drugi: odczyt i zapis idą tymi samymi bytami, tą samą komendą i tym samym identyfikatorem okna, więc oba widoki przyjmują wyłącznie stan potwierdzony przez rdzeń, rozgłaszany do wszystkich połączeń konta — zmiana dokonana w pasku dochodzi do kolumny sterowania tą samą drogą, którą dochodzi do drugiego urządzenia konta, i odwrotnie.

## budowa/klient-poprzedni/src/ustawienia/zrodlo-nastaw.ts
To osobne źródło obok źródła wartości okna konfiguracji: tamto czyta poziom w całości, skleja łańcuch
ośmiu zasięgów i rozstrzyga oś, bo tamto okno rysuje każdą pozycję katalogu na każdym poziomie. Sekcja
ustawień pyta o jeden klucz na jednym poziomie i potrzebuje przy tym definicji katalogu — opcje wyboru
przychodzą z katalogu, nie z wykazu zaszytego w kliencie; komendy są te same i idą tym samym wywołaniem
na tym samym kanale. Dwie właściwości rdzenia, na które ten plik jest przygotowany: zdarzenie zmiany
konfiguracji o rodzaju usunięcia niesie w zapisie wartość zdjętą, nie nową, więc subskrypcja zmiany klucza
oddaje przy usunięciu brak, a nie treść wpisu; rdzeń nie sprawdza poziomu zapisu wobec dozwolonych
zasięgów definicji i przyjmuje zapis na poziomie węższym, niż katalog dopuszcza, a taki zapis wygrywa
potem rozstrzyganie osi, więc poziom bierzemy zawsze z dozwolonych zasięgów definicji, nigdy z domysłu
wołającego. Definicja niesie dozwolone zasięgi i to ona rozstrzyga; klient nie wybiera poziomu za katalog,
bo rdzeń takiego zapisu nie odrzuci i pomyłka byłaby cicha. Odmowa nie wywraca odczytu: brak definicji
albo brak wartości wraca jako pole puste wraz ze zdaniem odmowy, a nie jako wyjątek.

## budowa/klient-poprzedni/src/okna-rownolegle/rodzaje-obszaru.ts
Te same liczby czyta rachunek podziału gniazda, uchwyt szerokości i arkusz
toru — trzy kopie tej samej liczby dałyby trzy różne progi czytelności po
pierwszej poprawce, a jedno źródło znaczy, że podniesienie minimum jest
jedną zmianą. Moduł nie zna DOM, nie zna paneli po nazwie i nie rozstrzyga,
ile paneli wolno otworzyć — oddaje liczby, a decyzję podejmuje ten, kto
pyta; gdy miejsca zabraknie, szerokość zmienia się uchwytem, panel nie
chowa się sam. Widok pełnoekranowy jest obszarem biorącym całą scenę, bez
sąsiada, z którym miałby się dzielić.

## budowa/klient-poprzedni/src/moduly/roundtable/panel-debaty.ts
Arkusz roundtable.css daje trzy stany obowiązkowe nośnika stany-okna, arkusz panel-debaty.css
wygląd głosów; arkusza debata.css panel nie wciąga, bo Debate Panel i Consensus Panel u gospodarza
nie stoją. Kształt panelu jest dokładnie ten, którego wymaga panel pomocniczy — element, odświeżenie,
zamknięcie — i ani jedno pole ponad to. Gospodarzem bywa pas okien pomocniczych modułu albo kolumna
paneli sceny okien równoległych; panel żadnego z nich nie zna i niczego o nich nie zakłada.
Zamknięcie zdejmuje trzy rzeczy, nie jedną: panel zakłada subskrypcję fragmentów wypowiedzi,
subskrypcję zmian stanu debaty oraz, przez stan debaty, nasłuch zdarzenia zmiany debaty i rejestru
kanałów. Wszystkie schodzą przy zamknięciu; panel bez tego zostawiłby je żywe po zejściu ze sceny
i rysowałby do elementu, którego nikt już nie ogląda. Własny stan debaty bierze się stąd, że opcje
panelu dają wyłącznie kanał, okno, moduł i przedrostek — egzemplarza stanu debaty tą drogą podać
się nie da, tak samo jak umowa opisu modułu nie przenosi rejestru kanałów. Drugi stan nie jest
drugą prawdą o debacie: czyta te same zdarzenia rdzenia i nie ma ani jednej drogi zapisu — różni
się od stanu złożenia wyłącznie chwilą otwarcia nasłuchu. Panel nie odczytuje przebiegu na żądanie,
choć kontrakt to przewiduje: komenda odczytu debaty oddaje skład, tury i wypowiedzi jednym
wywołaniem, ale obsługi tego odczytu jeszcze nie zbudowano — ani tutaj, ani w oknach złożenia
modułu — więc panel pokazuje wyłącznie to, co usłyszał od swojego otwarcia. Odświeżenie odnawia
zatem wykaz kanałów (nazwy uczestników) i przerysowuje widok, a przycisk nazywający brakującą
obsługę stoi widoczny i klikalny zamiast zniknąć.

Pierwszy rysunek idzie od razu: bez niego panel stałby pusty aż do pierwszego zdarzenia rdzenia,
a stan pusty ma własne zdanie mówiące, dlaczego jest pusty.

## budowa/klient-poprzedni/src/moduly/studio/odmowa-rdzenia.ts
Rdzeń odmawia komendy bez uchwytu kopertą osobnego typu, która nie niesie pola stanu, a korelacja
klienta rozstrzyga wyłącznie koperty ze stanem, więc zwykłe wywołanie na komendzie bez uchwytu
zostałoby obietnicą nierozstrzygniętą i okno stałoby w stanie ładowania bez końca. Wywołanie
łączy więc dwie drogi w jeden wynik: odpowiedź skorelowaną oraz zdarzenie odmowy o tym samym
identyfikatorze żądania — odmowa wraca jako zwykłe niepowodzenie, więc okno obsługuje ją tą samą
ścieżką co każdy inny błąd i nie buduje drugiego mechanizmu. Nazwa zdarzenia pochodzi z kontraktu:
funkcja tnie typ komendy po separatorze obszaru i sięga do mapy zdarzeń odmowy, w której stoją
obszary studia i okna, więc obie drogi odmowy modułu są rozpoznawalne bez literału nazwy. Kod
odmowy jest tu adekwatny, bo brakuje uchwytu, a nie treści żądania, a znacznik niepowtarzalności
powstrzymuje widok przed ponawianiem czegoś, czego rdzeń nie nabędzie przed wdrożeniem nowej
wersji.

## budowa/klient-poprzedni/src/sterowanie/indeks.ts
Punkt wejścia klienta montuje komplet sterowania dla każdego okna z osobna:
tworzy rejestr kanałów, odświeża go, a następnie tworzy panel sterowania
z kanałem, oknem i rejestrem, dołączając jego element do kontenera. Rejestr
kanałów powstaje raz na klienta — jest katalogiem wyboru, nie ustawieniem
okna. Panel powstaje raz na okno i nie ma z innym panelem żadnej wspólnej
zmiennej.

## budowa/klient-poprzedni/src/sterowanie/modul-okna.ts
Wykaz modułów pochodzi z katalogu modułów rdzenia, nie ze stałej listy znanych
identyfikatorów — tamta lista niesie środowiska, a nie moduły. Moduł spoza
katalogu pozostaje widoczny i wybieralny — katalog jest listą informacyjną,
nie bramą.

## budowa/klient-poprzedni/src/sterowanie/komunikat-zmiany.ts
Komunikat potwierdzenia i jego odbiorca istnieją osobno od reszty kompletu
sterowania, bo wiele elementów sterujących zgłasza zmianę tą samą drogą,
a miejsce pokazania komunikatu jest jedno.

## budowa/klient-poprzedni/src/sterowanie/rola-okna.ts
Wskazanie koordynatora, któremu podlega wykonawca, jest osobnym powiązaniem
między oknami — komplet sterowania jednego okna go nie ustanawia.

## budowa/klient-poprzedni/src/powloka/usuniecie-sesji.ts
Usunięcie nie jest zamknięciem ani archiwizacją. Rdzeń rozdziela trzy czynności i tak samo rozdziela
je powłoka: zamknięcie zmienia sam stan sesji, wiadomości, okna i katalog roboczy zostają nietknięte;
archiwizacja wyprowadza sesję z historii bieżącej i pozwala ją przywrócić, zapis zostaje w całości;
usunięcie wyprowadza zapis do kosza rdzenia, po terminie kosza rdzeń kasuje go trwale, a w oknie
terminu zapis wraca przywróceniem. Dlatego czynność nazywa się „Usuń trwale", nigdy „Zamknij" — nazwa
wzięta od sąsiada byłaby kłamstwem o skutku. Potwierdzenie jest warunkiem kontraktu, nie ozdobą
widoku. Pole potwierdzenia przechodzi tędy dokładnie takie, jakie podał wywołujący — powłoka nie
dopisuje go sama i nie uprzedza odmowy rdzenia własnym sprawdzeniem. Żądanie bez potwierdzenia wraca
z rdzenia jako niepowodzenie walidacji i tak ma je zobaczyć Operator. Dwa wykazy znaczą dwie różne
rzeczy. Wykaz usuniętych to sesje, których zapis zniknął; wykaz nieznalezionych to wskazania bez
odpowiednika w historii, które nie są błędem — widok nie ma prawa ani ich przemilczeć, ani wliczyć do
usuniętych, bo milczenie o nich znaczyłoby „usunąłem", gdy nie było czego usuwać. Kształt sprawdzamy
na wykazie usuniętych i na liczniku usuniętych: wykaz nieznalezionych jest w kontrakcie polem
opcjonalnym, więc jego brak jest odpowiedzią poprawną, a nie odpowiedzią o złym kształcie. Dopełnienie
wykazu nieznalezionych do wykazu pustego jest tu wolne od zwykłego zarzutu o operator pustej
wartości w źródłach: tędy przechodzą wyłącznie odpowiedzi udane, niepowodzenie zatrzymuje wcześniej
przeniesienie wraz z powodem, a brak pola opcjonalnego w odpowiedzi udanej znaczy dokładnie „żadne
wskazanie nie było bez odpowiednika". Odmowy ta gałąź nie widzi i nie ma czego przesłonić.

## budowa/klient-poprzedni/src/polaczenie/adres-rdzenia.test.ts
Ustalenie adresu gniazda jest miejscem, w którym klient albo trafia w rdzeń, albo szuka go tam, gdzie nikt nie nasłuchuje — a wtedy ponawianie milczy bez końca i wygląda jak awaria rdzenia, choć jest pomyłką adresu, dlatego sprawdzane są wszystkie drogi w jednym ustalonym pierwszeństwie. Adres z powłoki jest stanem modułu, dlatego każdy przypadek testowy wczytuje moduł na nowo, inaczej mierzyłby ślad po przypadku poprzednim.

## budowa/klient-poprzedni/src/okna-rownolegle/stan-paneli.ts
Przełącznik paneli stoi w nagłówku okna rozmowy, więc każde okno ma własny
zestaw paneli: panel należy do rozmowy, nie do ekranu. Moduł nie ma żadnej
zmiennej na poziomie modułu, bo jeden egzemplarz stanu na moduł znaczyłby
jeden zestaw paneli na całą aplikację — cały stan siedzi w domknięciu
fabryki, a cztery gniazda tworzą cztery niezależne egzemplarze, z których
żaden nie widzi pozostałych. Rozstrzyganie, co da się otworzyć, należy do
spisu okien pomocniczych. Otwarcie panelu już otwartego nie jest błędem, a
wartość spoza zakresu jest przycinana zamiast wstrzymywać wykonanie:
powtórne otwarcie nic nie zmienia i nie budzi subskrybentów.

## budowa/klient-poprzedni/src/moduly/roundtable/pasek-uczciwosci.ts
Pierwsza rzecz: rozjazd katalogu okien rdzenia z oknami, które moduł buduje. Rdzeń przypisuje
modułowi cztery okna operacyjne, a zamierzenie modułu wymienia sześć okien własnych — Argument Map
& Analysis i Voting & Evaluation Center nie mają w katalogu rdzenia ani jednego wiersza. Moduł je
buduje mimo to, bo bez nich okien byłoby cztery zamiast sześciu; zdanie rozjazdu liczy katalog okien
z odczytu wykazu modułów, nie z napisu. Druga rzecz: stosunek liczby narzędzi zamierzenia do stanu
ich wykonania. Bilans mówi, ile moduł wykonuje w całości, ile częściowo, ile ma w kontrakcie
komendę bez zbudowanej obsługi i ile nie ma pokrycia. Rozróżnienie dwóch ostatnich liczb jest po
scaleniu kontraktu najważniejszą treścią pasa: obszar urósł z czterech komend do czterdziestu
sześciu, więc niemal wszystko, co moduł nazywał brakiem, jest dziś pracą do wykonania, a nie
brakiem uzgodnienia. Trzecia rzecz: odczyt niewywoływany. Komenda odczytu debaty oddaje skład,
tury i wypowiedzi, a komenda wykazu modeli sam skład — żadnej z nich moduł jeszcze nie wywołuje,
więc okno otwarte w trakcie debaty zna wyłącznie to, co usłyszało od swojego otwarcia. Dotyczy to
każdego okna modułu naraz, dlatego zdanie stoi na pasie, a nie w jednym z nich. Pas niczego nie
blokuje i niczego nie ocenia — podaje liczby i nazywa granicę.

## budowa/klient-poprzedni/src/powloka/wedrowka-kart.ts
Jedna odpowiedzialność: przełożenie klawisza na czynność pasa. Strzałki przenoszą wybór z zawijaniem,
Home i End skaczą na krańce, Enter i spacja wybierają kartę pod fokusem, Delete ją zamyka. Plik nie
wie, co wybór i zamknięcie znaczą: w pasie związanym z rdzeniem obie czynności są komendami
kontraktu, w pasie samego widoku — zmianą miejscową. Rozstrzyga to moduł kart sesji, który podaje tu
czynności.

## budowa/klient-poprzedni/src/ustawienia/sekcja-uwierzytelnianie.ts
Dwóch komend rodziny tu nie ma: rejestracja odmawia trwale po pierwszym uruchomieniu, a logowanie jest
samą bramką — obie należą do ekranu logowania. Wykaz metod przychodzi z czynności albo ze zdarzenia,
nie z odczytu: kontrakt nie ma komendy odczytu metod, ma za to zdarzenie zmiany bramki z pełnym wykazem
metod po zmianie; sekcja je subskrybuje, więc wykaz nadąża także za czynnością wykonaną w drugim oknie
albo na drugim urządzeniu. Nie obejmuje to pierwszego otwarcia: zanim zajdzie jakakolwiek zmiana, nie ma
czego rozgłosić i wykaz jest pusty — sekcja mówi to wprost, zamiast pokazywać pustkę, którą dałoby się
odczytać jako brak metod. Metoda systemowa jest nieczynna; zdanie na ekranie cytuje odmowę rdzenia co do
słowa, żeby powód na ekranie i powód w odmowie nie mówiły dwóch rzeczy o tej samej niedostępności. Wymóg
logowania jest tu pokazywany, ale zmienia się go w oknie konfiguracji.

## budowa/klient-poprzedni/src/moduly/workspace/pamiec-kontekstu.ts
Pamięć jest odrębna dla projektu, a zakresem współdzielenia jest poziom zasięgu wpisu: zapis
szerszy niż projekt wchodzi do pamięci innych projektów po włączeniu ustaleń wspólnych. Wpis
o pochodzeniu model jest propozycją czekającą na decyzję: przyjęcie zapisuje ten sam wpis
z pochodzeniem operator, a odrzucenie usuwa go komendą usunięcia wpisu pamięci. Zmienia się
kolejność prezentacji zasięgów, nie nastawa: nazwy poziomów są w całym module jedne.
Trzy kontrolki paska narzędzi nie zmieniają pamięci projektu, tylko zawężają wykaz już odczytany.
Pamięć karty sesji prowadzi osobna rodzina komend, oddzielona od pamięci projektu, żeby mieszanie
obu wykazów nie zatarło, który wpis czyją własnością jest. Wpis zapisany w zasięgu szerszym niż
projekt bywa niewidoczny w wykazie poniżej, bo wykaz pokazuje wspólne dopiero po włączeniu
przełącznika. Rdzeń odmawia usunięcia wpisu, którego nie ma, i okno pokazuje tę odmowę wprost.
Zapis scali dołączoną treść w jeden wpis; wpis źródłowy zostaje w pamięci projektu, dopóki nie
usunie się go osobnym przyciskiem, a przycisk bierze powód z rdzenia, więc kopia w innym miejscu
mogłaby się rozjechać. Przy zapytaniu w drodze rysowanie pokazywałoby wykaz sprzed zmiany;
plakietka mówi, ile wpisów rdzeń oddał i ile z nich przeszło przez filtr, bez drugiej liczby
wykaz zawężony wyglądałby jak pamięć uboższa, niż jest.

## budowa/klient-poprzedni/src/moduly/studio/okno-preview-window.ts
Podgląd dokumentu przestał być osobnym oknem: jest trybem widoku okna pracy z dokumentem, bo
pokazywał tę samą treść co edytor, tylko w formacie wyjściowym, a jego trzy czynności poszły tam
wraz z nim bez przepisania — wydanie i przekazanie do Library dalej idą przez czynności podglądu,
a decyzja o wyniku przez decyzję o propozycji. Definicja okna w katalogu rdzenia jest jedna
i miała dwa przypięcia; zostało przypięcie modułu Design, więc rozróżnik wariantu niesie dziś
jedną wartość, ale zostaje jawny, bo drugie przypięcie może wrócić i wtedy ma się rozejść polem,
nie odpisem pliku. Wariant Designu podgląda zasób wizualny: rdzeń nie generuje obrazów, więc
zasób wraca bez adresu, formatu i wymiarów, a podgląd pokazuje nazwę zasobu, wykaz pól wraz z
tymi, których rdzeń nie podał, oraz porównanie wariantów po treści słownej, nie po obrazie.
Przekazanie potwierdza się modułem z pola odpowiedzi, bo samo przekazanie kontekstu nie sprawdza
katalogu modułów, a kod modułu docelowego jest zamówieniem, nie potwierdzeniem. Eksport zasobu
i akceptacja wyniku nie mają odpowiednika w kontrakcie i stoją jako nazwane braki drogi, nie jako
martwe przyciski.

## budowa/klient-poprzedni/src/sterowanie/naglowek-sterowania.ts
Jeden wiersz obsługuje cały komplet sterowania — listę wyboru, suwak nakładu,
pole hosta i wykaz katalogów — bo wszystkie potrzebują tego samego nagłówka.
Etykieta nie obudowuje kontrolki: znak objaśnienia jest przyciskiem,
a wewnątrz obudowy jego naciśnięcie przenosiłoby się na kontrolkę zamiast
pokazać objaśnienie. Wiązanie idzie więc atrybutem wskazującym pole, nie
zagnieżdżeniem.

## budowa/klient-poprzedni/src/sterowanie/pasek-komunikatow.ts
Pasek jest jedyną reakcją interfejsu na niepowodzenie zmiany: informacja. Nie
wyłącza sterowań, nie zamyka okna i nie wymusza potwierdzenia — kolejna próba
idzie zwyczajnie, bo błąd dotyczy wyłącznie bieżącego wywołania.

## budowa/klient-poprzedni/src/sterowanie/wskaznik-odczytu.ts
Wskaźnik stoi obok pola, nigdy zamiast pola: podmiana kontrolki na wskaźnik
byłaby blokadą, a Operator ma móc wybrać wartość także wtedy, gdy katalog
jeszcze jedzie z rdzenia. Wskaźnik mówi też technologiom wspomagającym, co się
dzieje: ogłasza zmianę bez zabierania ogniska, a etykieta dostępności niesie
nazwę katalogu, którego dotyczy odczyt.

## budowa/klient-poprzedni/src/sterowanie/rejestr-modulow.ts
Wykaz idzie z rdzenia, nie ze stałej listy znanych identyfikatorów kontraktu:
tam stoją cztery pozycje, a są to środowiska, nie moduły. Katalog rdzenia
niesie piętnaście modułów zasilanych migracją bazy. Podstawienie wykazu
środowisk pod pole pytające o moduł pokazywałoby każdy prawdziwy moduł jako
spoza wykazu, a Operatorowi proponowałoby środowiska tam, gdzie pyta się
o moduł. Rejestr niczego nie zapisuje — jest wyłącznie odczytem. Wykaz jest
pusty do chwili odpowiedzi rdzenia; osobna flaga rozróżnia „jeszcze nie wiem"
od „katalog jest pusty", bo widok musi te dwa stany rozróżnić.
