# Danaco Console — Uzasadnienia komentarzy klienta poprzedniego

Dokument gromadzi uzasadnienia, które przekraczają dopuszczalną długość
nagłówka komentarza w plikach `budowa/klient-poprzedni/`. Każdy rozdział nosi
nazwę pliku źródłowego, którego uzasadnienie dotyczy.

## budowa/klient-poprzedni/src/moduly/design/czynnosci-warsztatow-designu.ts

Katalog opisuje pięć warsztatów danymi zamiast pięciu odrębnych zestawów
formularzy, wzorem pliku `studio/czynnosci-warsztatu.ts` — okno buduje pola
z tego wykazu i składa żądanie wspólną funkcją `zloz`. Druga, osobna rodzina
typów rozjechałaby się z pierwszą przy każdym nowym rodzaju pola, co
wymuszałoby poprawki w dwóch oknach zamiast w jednym.

Nazwy pól są nazwami kontraktu, ponieważ trafiają wprost do rdzenia. Etykiety
są zdaniem Operatora i z nazwami pól się nie pokrywają — Operator czyta
„proporcje kadru", a rdzeń otrzymuje `aspectRatio`.

Nazwa komendy pochodzi wyłącznie ze stałych kontraktu (`Command.*`), nigdy
z napisu wpisanego ręcznie: napis literowany z pamięci przechodzi sprawdzian
typów, lecz zawodzi dopiero u Operatora, gdy rdzeń takiej komendy nie zna.

Podpowiedź pola wykazu czynności bierze nazwę komendy z tych samych stałych
kontraktu z tego samego powodu — napis wpisany wprost przeżyłby zmianę nazwy
w kontrakcie i podpowiadałby Operatorowi komendę, której rdzeń już nie zna,
a to podpowiedź jest tym, co Operator przepisuje do pola.

## budowa/klient-poprzedni/src/moduly/design/katalog-funkcji-designu.ts

Katalog istnieje po to, żeby stan modułu Design dało się przeczytać, a nie
zgadnąć. Okno pokazuje czynności, które wykonuje; katalog pokazuje komplet
zamierzonych pozycji i przy każdej mówi jedno z trojga: którą komendą
kontraktu jest wykonywana, że wykonuje ją samo okno bez udziału rdzenia, albo
czego brakuje, żeby była. Pozycja bez drogi nie znika z wykazu, ponieważ
zniknięcie byłoby ukryciem braku.

Trzeci stan jest konieczny, bo dwa nie oddają prawdy o tym module. Kanwa,
wyrównanie warstw, siatka pomocnicza, drzewo żetonów i rachunek kontrastu
dzieją się w całości w przeglądarce i komend nie potrzebują; nazwanie ich
brakiem kontraktu byłoby zmyśleniem długu, a nazwanie komendą — zmyśleniem
drogi wykonania, której nie ma.

Nazwy pozycji i podział na grupy nie są tłumaczone ani parafrazowane w żadnym
miejscu, które je wyświetla — pozycję odnajduje się po pełnej nazwie własnej.

Liczba pozycji katalogu nie jest zapisana nigdzie na stałe jako osobna
wartość, tylko liczona z wykazu funkcją length w miejscu użycia, żeby napis
nie mógł rozjechać się z rzeczywistą zawartością wykazu, który stoi obok.

Katalog jest zbiorem danych, nie widokiem — wyszukiwarka funkcji modułu
Design buduje z niego listę do wyświetlenia.

## budowa/klient-poprzedni/src/moduly/apps/narzedzia-apps.ts

Narzędzia są pogrupowane wedle okna modułu Apps, w którym dana czynność się
odbywa: Product Builder prowadzi produkt, etapy, kamienie milowe i oś czasu;
Architecture Designer — walidację układu, wersje, adnotacje i eksport;
warsztaty — podgląd, trasy, motyw, punkty końcowe i schemat; Deployment Panel
— środowiska, zmienne, domenę, skalowanie, kondycję, dzienniki i artefakty;
Publisher Panel — pakowanie, manifest, walidację, podpis i publikację.

Zdanie o skutku jest w każdym narzędziu inne, bo skutek jest inny: liczba
pozycji, nazwa bytu, kod odpowiedzi, odwołanie do pliku, wynik weryfikacji
podpisu. Wspólny komunikat „gotowe" byłby meldunkiem, z którego nic nie
wynika, a to jest dokładnie wzorzec, którego moduł ma nie powtarzać.

Wartości startowe pól są przykładami z domeny produktu, nie wartościami
wymuszonymi: Operator zmienia je przed wykonaniem, a puste pole nieobowiązkowe
oznacza brak zawężenia.

## budowa/klient-poprzedni/src/moduly/apps/zrodlo-apps.ts

Obszar apps.* jest dwukierunkowy: obok trzech komend zapisu (architektura,
plik warsztatu, wdrożenie) stoją trzy komendy odczytu — apps.deployment.list,
apps.architecture.get, apps.workspace.list. Bez nich moduł po odświeżeniu
okna przeglądarki zaczynałby od zera, bo wypełniałyby go tylko zdarzenia
bieżącej sesji gniazda, a historia wdrożeń trzymana w bazie znikałaby z oczu.
Każdy odczyt oddaje ten sam kształt, którym odpowiada jego komenda zapisu.

Zdarzenie apps.workspace.changed jest drogą rozgłoszenia dla
apps.workspace.update: bez niego plik zapisany w jednym oknie nie docierałby
do drugiego, dopóki plik nie zostałby odczytany ręcznie. Wchodzi tą samą
bramą co apps.build.changed — subskrypcją źródła, nie własnym gniazdem.

Źródło nie ma własnego stanu: jest warstwą wywołań i sprawdzianu kształtu
odpowiedzi. Stan produktu mieszka osobno, żeby pięć okien patrzyło na jeden
zbiór, a nie na pięć kopii.

Komendy zapisu mają uchwyt w rdzeniu, więc wywołanie wraca zwykłą
odpowiedzią albo odmową merytoryczną, nie kopertą apps.unknown. Droga przez
warstwę odmowy rdzenia zabezpiecza ścieżkę fail-open: gdyby uchwyt zniknął
albo rdzeń nie rozpoznał którejś komendy, obietnica wywołania ma się czym
rozstrzygnąć zamiast wisieć bez końca, a okno nazywa odmowę zamiast ją ukryć.

Pola tożsamości idą do rdzenia tak, jak je wpisano, bez przycięcia po drodze
po stronie klienta. Przycięcie w przeglądarce i przycięcie w rdzeniu nie są
tym samym działaniem i rozjeżdżają się na dwóch znakach białych, których
jedna strona zdejmuje, a druga zostawia. Ma to znaczenie na ścieżce
warsztatu, bo klucz złożony z okna, warstwy i ścieżki czyni ją tożsamością
pliku: przycięcie po stronie klienta zapisałoby plik pod ścieżką inną niż
wpisana albo odrzuciłoby ścieżkę, którą rdzeń przyjmuje. Rdzeń przycina te
pola po swojemu i to on jest tu jedyną władzą; o pustce pola rozstrzyga
więc pustka dosłowna, a okno zestawia potem wpisane z oddanym i ogłasza
różnicę.

Sprawdzian kształtu odpowiedzi jest przy każdej komendzie osobny i pyta
o pole, którego okno naprawdę używa. Sprawdzanie, czy cokolwiek wróciło,
przepuściłoby odpowiedź o kształcie innym niż kontraktowy, a okno wywróciłoby
się dopiero przy rysowaniu. Pustka nie jest tu uszkodzonym kształtem: wykaz
pusty i pole opcjonalne bez wartości są odpowiedziami prawdziwymi i znaczą,
że w tym oknie danej rzeczy jeszcze nie ma — dlatego przy nich sprawdzian
pyta o tablicę albo przepuszcza wszystko, zamiast żądać obiektu.

## budowa/klient-poprzedni/src/moduly/apps/narzedzia-rozszerzen.ts

Narzędzia są pogrupowane wedle okna strony rozszerzeń: App Catalog prowadzi
wyszukiwarkę, kartę szczegółów, kolekcje i rejestr organizacji; Installed
Apps Manager — aktualizacje, przesyłkę paczki, wersjonowanie, instalację
zestawu, dziennik cyklu życia i tryb administracyjny; Integrations Hub i MCP
and Connector Console — transport, poświadczenie, odkrywanie narzędzi, próbne
wywołanie, log protokołu, piaskownicę, import definicji, webhooki,
odwzorowania, metryki i kondycję; Permissions and Trust Center — uprawnienia,
podpis, skaner i sekrety.

Żadne pole nie niesie treści poświadczenia. Pola sekretów przyjmują wyłącznie
klucz jawny warstwy sekretów — nazwę, po której rdzeń wydaje wartość.
Wprowadzenie hasła, tokenu czy klucza API zostaje po stronie Operatora, poza
tą drogą i poza kontraktem.

Pozycja, na której narzędzie ma pracować, bierze się ze wskazania w oknie
funkcją stan.wybrane(), a nie z pola wpisywanego przy każdym przycisku. Brak
wskazania kończy się odmową z powodem — Operator ma najpierw wybrać pozycję
z listy.

## budowa/klient-poprzedni/src/komponenty/menu-drzewo.ts

Mechanizm stoi w bibliotece komponentów, a nie wewnątrz paska zlecenia, bo
pasek jest tylko jego największym odbiorcą: takie samo menu stawiają
nagłówki okien, panele i ekrany konfiguracji. Mechanizm zamknięty w pasku
byłby dla nich nieosiągalny i zostałby skopiowany, a kopie rozjeżdżają się
w szczegółach — haczyk wyboru pojawiałby się w jednym menu, a w drugim nie,
przy jednakowym wymogu.

To jest mechanizm docelowy dla sterów nastawy: dopisanie obok niego kolejnego
mechanizmu rozwijania odtwarza dokładnie ten rozjazd, który on likwiduje.
Trwałym wyjątkiem zostaje mechanizm menu czynności z osobnym plikiem, bo
obsługuje menu, którego uchwyt jest ikoną, a treść wchodzi jako gotowe
elementy, a nie ster niosący bieżącą wartość — to osobny wzorzec, nie rozjazd
do scalenia.

Mechanizm niesie sam sześć cech, których odbiorcy nie budują u siebie:
zagnieżdżenie gałęzi bez ograniczenia głębokości; znacznik bieżącego wyboru
na liściu i ślad tego wyboru na gałęzi, która ten liść niesie, żeby Operator
widział wybór bez wchodzenia w gałąź; opis jako pole każdej pozycji;
przełącznik dwustanowy wewnątrz gałęzi; grupowanie pozycji po źródle lub
rodzinie; drogę do rejestru na dole menu. Treść wchodzi jako dane, a formę
nadaje mechanizm; przyjmowanie gotowych elementów od odbiorcy zniosłoby
jednolitość tych sześciu cech.

Mechanizm nie wie, co jest w drzewie, i nie wykonuje wyboru — oddaje klucz
pozycji wołającemu; nie zna kontraktu, komendy ani stanu okna. Żadna pozycja
nie dostaje stanu wyłączonego: drzewo puste nie jest błędem, uchwyt otwiera
się i mówi zdaniem, że wykaz jest pusty, zamiast przestać reagować.

Pozycja niesie równocześnie nazwę pełną i nazwę skróconą. W wykazie płaskim
na setki pozycji nazwa pełna musi nieść źródło, bo inaczej dwie komendy
o tej samej nazwie z dwóch wtyczek są nie do rozróżnienia. To samo źródło
czyni jednak wykaz nieczytelnym, gdy sto pozycji zaczyna się tym samym
przedrostkiem — dlatego pozycja niesie obie nazwy, pełną ze źródłem i
skróconą do samej rzeczy. Szukanie obejmuje obie nazwy i opis; trafienie
w nazwę skróconą liczy się wyżej niż w pełną, bo to ona jest tym, czego
Operator naprawdę szukał.

Opis pozycji stoi wprost pod nazwą, a nie za dymkiem wywoływanym najechaniem,
bo w menu o kilkudziesięciu liściach najeżdżanie na każdy liść z osobna nie
jest drogą; dymek zostaje formą dla kontrolek stojących pojedynczo.

Nazwa nastawy (na przykład „Model", „Wysiłek", „Urządzenie") nie trafia na
ekran, bo na uchwycie stoi wartość, nie nazwa nastawy. Idzie do etykiety
dostępności uchwytu i listy, żeby czytnik ekranu wiedział, czego dotyczy
wartość, której nazwa sama tego nie mówi — „Opus 5" nie niesie słowa „model".

Próg liczby liści, od którego drzewo stawia pole szukania, liczy liście, nie
wszystkie pozycje: drzewo o trzech gałęziach i dwóch liściach filtra nie
potrzebuje, drzewo o dwóch gałęziach i stu pozycjach potrzebuje go
natychmiast. Próg, a nie stałe pole szukania, bo pole nad wykazem pięciu
pozycji zabierałoby wiersz i nie skracałoby żadnego ruchu.

Kierunek rozwinięcia wykazu nie jest preferencją wizualną. Ster stojący
w pasku u góry okna ma pod sobą całą wysokość sceny i rozwija się w dół.
Ster przy polu wpisywania stoi u dołu okna, a wykaz rozwinięty w dół nie
miałby dokąd pójść — wyszedłby poza krawędź sceny.

Tryb bez uchwytu obsługuje wykaz komend wywoływany ukośnikiem: wykaz ma się
pojawić natychmiast po wpisaniu znaku, a filtrem ma być to samo pole
wpisywania, nie osobne okienko. W tym trybie uchwyt nie wchodzi do
dokumentu, wewnętrzne pole szukania nie powstaje wcale, a wołający steruje
mechanizmem przez rozwinięcie, ustawienie frazy, przesunięcie wyróżnienia
i wybór pozycji wyróżnionej. Ognisko zostaje w polu wołającego — stąd
wyróżnienie jest wirtualne, a nie ogniskiem.

Opis pokazywany tylko przy pozycji wyróżnionej to ustawienie dla wykazu na
setki pozycji: setka opisów naraz nie jest objaśnieniem, tylko ścianą
tekstu, przez którą nie widać już samych nazw — wtedy opis należy się jednej
pozycji, tej, na którą Operator właśnie patrzy.

Metoda ustawiająca całą zawartość naraz podaje wartość widoczną na uchwycie
i całe drzewo. Drzewo przerysowuje się w całości, bo wykaz bywa czytany
z rdzenia i pozycje przychodzą oraz znikają; rozwinięte gałęzie przeżywają
przerysowanie po kluczu, żeby odświeżenie katalogu nie zwijało menu pod ręką
Operatora.

Fraza filtrująca podana z zewnątrz wchodzi do wewnętrznego pola, gdy pole
stoi na ekranie, żeby Operator widział, czym wykaz jest przycięty. Gdy pole
nie stoi (tryb bez uchwytu albo wykaz poniżej progu szukania), fraza i tak
przycina wykaz — próg rozstrzyga, czy mechanizm stawia pole, a nie czy
w ogóle umie filtrować.

Przesunięcie wyróżnienia nie rusza ogniska: przy obsadzie ukośnikiem
strzałki naciska się w polu wpisywania, więc ognisko musi w nim zostać;
przenoszenie go na pozycję wyrwałoby Operatorowi klawiaturę spod palców
w połowie pisania.

Wybór pozycji wyróżnionej oddaje wartość fałsz, gdy nie ma czego wybrać —
wykaz pusty albo przycięty do zera — i wtedy wołający wie, że klawisz Enter
ma zrobić swoje zwykłe zadanie zamiast niczego. Wyróżniona gałąź nie jest
wyborem — Enter na niej otwiera poziom.

## budowa/klient-poprzedni/src/moduly/browser/sekcje-rodzin.ts

Plik składa wyłącznie wywołania rdzenia w zdania o skutku — elementy stawia
`panel-rodzin.ts`, a stan modułu trzyma `stan-przegladania.ts`. Zdanie
zwracane po czynności niesie zawsze liczbę albo identyfikator wzięty wprost
z odpowiedzi rdzenia, nigdy samo potwierdzenie w rodzaju „gotowe" — zdanie
bez pokrycia w odpowiedzi byłoby meldunkiem bez skutku, którego Operator nie
mógłby sprawdzić.

## budowa/klient-poprzedni/src/moduly/design/zrodlo-designu.ts

Zbiór zasobów mieszka w `stan-designu.ts`, żeby trzy okna patrzyły na jeden
zbiór, a nie na osobne kopie.

Zasób wchodzi do modułu dwiema drogami, które się nie zastępują: `design.asset.generate`
daje zasób z pracy modelu i wymaga kanału obrazowego, a jego brak jest odmową
nazywającą brak, nigdy obrazem zastępczym; `design.asset.upload` daje zasób
z pliku wskazanego przez Operatora i działa bez żadnego kanału modelu, wnosząc
treść już istniejącą.

Przy wgraniu Operator wskazuje plik, klient czyta go i oddaje treść w polu
`contentBase64`, a rdzeń zapisuje ją w swoim magazynie pod sumą kontrolną. Od
tej chwili zasób nie zależy już od pliku na dysku, który Operator może
nadpisać albo skasować.

Treść idzie bajtami, nie ścieżką, mimo że kontrakt zna oba pola. `sourcePath`
każe rdzeniowi otworzyć plik pod wskazaną ścieżką, a rdzeń biegnie na innej
maszynie niż przeglądarka Operatora: ścieżka z jego pulpitu nie znaczy tam nic
albo znaczy coś zupełnie innego. Klient czyta plik sam i wysyła to, co
przeczytał.

Odmowa wraca polem `blad`, nie wyjątkiem: rdzeń odmawia czynności na zasobie,
którego nie zna, a okno pokazuje to zdanie i zostaje czynne.

Wskazanie kanału zlecenia generowania jedzie polem `channelId`: rdzeń sprawdza
po nim rodzaj kanału i odmawia wskazania tekstowego.

Pola opisowe zlecenia wgrania — nazwa, format, wymiary — idą tylko wtedy, gdy
klient je zmierzył; podstawienie wartości domyślnej byłoby wpisaniem rdzeniowi
metadanych, których nikt nie sprawdził.

Skala i jakość zlecenia wydania są w kontrakcie opcjonalne: zero wysłane do
rdzenia byłoby żądaniem obrazu o zerowym boku albo zerowym stopniu kompresji
i wróciłoby odmową, dlatego pole zerowe nie wchodzi do żądania.

## budowa/klient-poprzedni/src/moduly/apps/zrodlo-rozszerzen-dobudowa.ts

Podział na dwa źródła (`zrodlo-rozszerzen-apps.ts` i to źródło) jest podziałem
obowiązków, nie rozmiaru. Pierwsze prowadzi cykl życia pozycji: wykaz,
instalację, konfigurację, przełącznik i odinstalowanie. To źródło niesie
wszystko, co robi się na pozycji już stojącej — wyszukiwanie, kolekcje,
wersjonowanie, rozmowę protokołem, webhooki, uprawnienia, podpis i sekrety.

Żadna z tych komend nie niesie identyfikatora okna: katalog rozszerzeń jest
bytem rdzenia stojącym poziom wyżej niż moduł, Apps jest jego operacyjnym
frontem, nie drugim źródłem prawdy, więc pozycji nie posiada żadne okno.

Sprawdzian kształtu odpowiedzi pyta przy każdej komendzie o pole, którego okno
naprawdę używa. Pustka nie jest tu uszkodzonym kształtem: wykaz pusty i pole
opcjonalne bez wartości są odpowiedziami prawdziwymi i znaczą, że tego jeszcze
nie ma.

## budowa/klient-poprzedni/src/moduly/automations/zrodlo-dobudowy.ts

Rozdział między tym plikiem a `zrodlo-automations.ts` idzie po roli, nie po
objętości. Tamten plik niesie rdzeń modułu — definicję, harmonogram, kolejkę,
zależności i przebiegi, czyli to, czym okna pracują bez otwierania
jakiegokolwiek panelu. Ten plik niesie czynności paneli i szuflad.

Każda czynność oddaje `Wynik` z całą odpowiedzią, nie z wyciętym polem, bo
większość tych odpowiedzi niesie więcej niż jedną rzecz naraz — wykaz i
znacznik przycięcia, automatykę i numer wersji, zlecenie i znacznik
duplikatu, zmienne i zastrzeżenia. Wycięcie jednego pola gubiłoby drugie,
a okno musi pokazać oba.

Sprawdzian kształtu pilnuje pola obowiązkowego odpowiedzi: rdzeń, który oddał
kopertę powodzenia bez treści, jest dla okna odmową, inaczej widok rysowałby
pustkę i twierdził, że to wynik.

## budowa/klient-poprzedni/src/moduly/automations/okno-queue-manager.ts

Okno nie ma własnego silnika kolejek ani własnego cyklu życia zlecenia: jeden
silnik obsługuje pętlę sesyjną i MultitaskingAI, a okno posuwa kolejkę
komendami kontraktu i pokazuje to, co rdzeń o niej oddaje.

Czego kontrakt nie oddaje, tego okno nie zmyśla. Byt `Queue` niesie stan,
licznik obiegów, okna, politykę, liczbę zleceń oczekujących i powiązaną
automatykę, lecz nie niesie samych zleceń — te oddaje osobna komenda wykazu
zleceń. Dopóki rdzeń nie ma jej uchwytu, przegląd pokazuje kolejkę jako
całość, a identyfikator zlecenia wpisuje Operator; okno mówi to wprost,
zamiast pokazywać pustą listę udającą przegląd.

O tym, czy rdzeń ma uchwyt danej komendy, okno samo nie orzeka i nie
przepisuje stanu kontraktu do zdań: powód pozycji paska bierze się z wykazu
komend oddanego przez rdzeń (`pokrycie-komend.ts`), a wykaz działań silnika
z wyliczenia kontraktu. Oba przestają być prawdziwe same, bez edycji tego
pliku.

Zdanie potwierdzenia powstaje z oddanej kolejki, nie z etykiety przycisku:
etykieta mówi o żądaniu, a rdzeń może oddać kolejkę w innym stanie, niż
żądanie zapowiadało.

Odpowiedź zasilenia kolejki krokami niesie stan kolejki i liczbę zleceń
oczekujących, lecz nie mówi, które zlecenia w niej stoją — zdanie po zasileniu
mówi więc to, co odpowiedź naprawdę niesie, i nie orzeka, że kroki weszły.

Priorytet zlecenia jedzie przy działaniu „wstrzymaj”, bo jest polem żądania,
nie osobnym działaniem silnika, a wstrzymanie jest działaniem odwracalnym
jednym naciśnięciem „Wznów”.

## budowa/klient-poprzedni/src/moduly/library/zrodlo-zarzadu.ts

Zarząd repozytorium jest wydzielony ze `zrodlo-biblioteki.ts` wzdłuż
odpowiedzialności: tamten plik prowadzi wykaz i podgląd, ten niesie czynności
zarządcze warstw trzeciej i czwartej modułu biblioteki (opis zasobu, słownik
etykiet, tezaurus, reguły, higienę, cykl życia, utrwalenie, udostępnienia
i sugestie). Kształt pozostaje jeden — `ZrodloBiblioteki` — więc okna nadal
widzą jedno źródło.

## budowa/klient-poprzedni/src/moduly/automations/okno-scheduler.ts

Cykliczność w oknie Schedulera ustala się wzorcem, nie samą składnią: Operator wybiera
opis w rodzaju „co tydzień, poniedziałek, 07:00”, a okno składa z tego zapis cron według
kontraktu wzorców z `wzorce-cyklicznosci.ts`. Pole zapisu cron pozostaje przy tym widoczne
i edytowalne — zapis wpisany wprost jest rozpoznawany z powrotem jako wzorzec, a zapis,
którego żaden wzorzec nie obejmuje, jest wzorcem własnym i idzie do rdzenia w całości, bez
udziału kreatora.

Harmonogram obowiązuje dopiero po powiązaniu z automatyką, dlatego okno wprost nazywa ten
warunek zamiast odmawiać zapisu bez wyjaśnienia przyczyny.

Potwierdzenie po zapisie opiera się na polach `enabled` i `nextRunAt` odpowiedzi rdzenia,
nie na treści żądania: gdyby zdania „wstrzymany” i „wznowiony” zależały od tego, o co okno
prosiło, rdzeń, który zapisu nie przyjął zgodnie z żądaniem Operatora, dostawałby od okna
potwierdzenie czynności, która się nie odbyła.

Pola kreatora cykliczności stoją w rzędzie zawsze, niezależnie od wybranego wzorca, zamiast
chować pola nieużywane przez dany wzorzec: pole schowane przy przełączeniu wzorca
przeskakiwałoby układ pod ręką Operatora, a pole nieużywane jest nieszkodliwe, bo jego
wartość po prostu nie wchodzi do zapisu.

Kalendarz uruchomień czyta przebiegi bez wskazania okna, więc rdzeń nie zakłada na oknie
Schedulera obserwacji telemetrii — kalendarz jest jednorazowym zdjęciem stanu, a stała
obserwacja przebiegów należy do Execution Monitora i ma pozostać jej jedynym odbiorcą
w module Automations.

## budowa/klient-poprzedni/src/moduly/apps/etykiety-apps.ts

Teksty tego pliku stoją osobno od plików budujących elementy okien modułu Apps, na wzór
podziału zastosowanego w `sterowanie/etykiety-sterowania.ts`. Pozycja wykazu braków niesie
tylko to, co okno wie samo: nazwę czynności, czego by do niej trzeba (pole `czego`), nazwę
komendy, której wejście do kontraktu znosi brak (pole `komendaZnoszaca` — pozycja znika
wtedy z wykazu sama) oraz cudze drogi do sprawdzenia w wykazie (pole `komendyCudze`).
Zdanie o stanie kontraktu dokłada `braki-kontraktu.ts`, czytając stałą `KOMENDY` przy
składaniu okna.

Nazwy komend proponowanych w tym pliku nie są nazwami kontraktu i nie idą na drut do
rdzenia: pozycja z takim wskazaniem stoi w wykazie braków, a wpisanie komendy do kontraktu
zdejmuje pozycję z wykazu samo. Nazwy stoją w jednym miejscu, żeby definicja oddana
właścicielowi produktu i zdanie widoczne Operatorowi mówiły o tej samej komendzie. Obszar
nazwy idzie za bytem, którego dotyczy: czynności na katalogu rozszerzeń należą do obszaru
`extension`, bo katalog stoi poziom wyżej niż moduł, a czynności na produkcie budowanym
w module — do obszaru `apps`.

## budowa/klient-poprzedni/src/moduly/apps/okno-integrations-hub.ts

Okno pokazuje dwa rodzaje pozycji katalogu — serwer protokołu MCP i integrację
przez interfejs API — bo tylko te dwa są drogą do usługi zewnętrznej. Wtyczki
i umiejętności nie należą do tego okna: są rozszerzeniami powłoki i sposobami
wykonania zadań, nie drogami do usług zewnętrznych, więc ich miejsce jest
w katalogu aplikacji.

Test połączenia i monitor zdrowia sprawdzają most, przez który dochodzi się do
serwera, a nie usługę stojącą za mostem. Sprawdzenie oddaje stan, czas
i korzenie potwierdzone przez punkt dostępu; czasu odpowiedzi samej integracji
ani liczby jej narzędzi kontrakt nie niesie, więc okno tej wartości nie zmyśla.

Integracja bez punktu dostępu nie jest usterką: pozycja rodzaju API bywa
dostępna wprost adresem, a punkt dostępu opisuje most do maszyny. Okno
odróżnia brak mostu od mostu nieznanego rdzeniowi, ponieważ to dwa różne stany
i drugi z nich jest niespójnością rejestru.

Transport protokołu MCP da się dziś wpisać wyłącznie do nieprzezroczystej
konfiguracji w postaci zapisu JSON. Kontrakt nie nazywa tego pola osobno, więc
rdzeń i okno mogą rozumieć wpisaną wartość odmiennie.

## budowa/klient-poprzedni/src/asystent-plywajacy/stan-dymka.ts

Dymek nie ma własnej warstwy wywołań — bierze `zrodlo-assistant.ts`
i `zrodlo-zaplecza.ts` modułu Assistant, żeby własna warstwa nie rozjechała się
z modułem przy zmianie kształtu odpowiedzi. Żadna ścieżka nie kończy się
milczeniem: ustalenie okna (`window.list`), wydanie polecenia
(`assistant.voice.command`), odczyt dziennika (`assistant.activity.list`)
oraz zamknięcie zlecenia błędem albo anulowaniem zawsze kończy się wypowiedzią
w historii, nigdy cichym `return` — pusta historia czytałaby się jak brak
głosu asystenta, co byłoby nieprawdą. Stan okna rozróżnia sześć wartości,
nie dwie, bo odmowa rdzenia i brak okna asystenta prowadzą do różnych
wniosków Operatora.

Odpowiedź asystenta przychodzi później niż odpowiedź komendy: rdzeń potwierdza
samo przyjęcie zlecenia, a wpis rodzaju `result` dopisuje przy domykaniu
(`adapter_modul_asystent_wykonawca.go` → `domknijZlecenie`), więc zejście
zlecenia z toru, ogłaszane zdarzeniem `assistant.action.changed`, pociąga
odczyt dziennika.

Źródło posunięć jest jedno i wspólne z pasem dolnym (`aplikacja/
pas-posuniec.ts`); druga subskrypcja byłaby drugim rozstrzyganiem sprawcy,
a `utworzRozstrzyganieSprawcy` zużywa odcisk okna przy weryfikacji, więc dwa
egzemplarze wydałyby dwa różne werdykty o jednym zdarzeniu.

Sprawca idzie z koperty zdarzenia: `AssistantActionChangedEvent` niesie
`actor` i `actorClientId` wypełniane przez rdzeń (`core/sprawca.go`).
Rozgłoszenie idzie do całego konta (`transport/rozgloszenie.go`), więc
zlecenie założone poza tym dymkiem wchodzi jako posunięcie (kto, co, w jakim
stanie), ale bez treści dziennika — treść odpowiedzi należy do okna,
w którym padło polecenie.

Odmowa `odnotujOdmowe` trafia do historii oprócz dymka powiadomienia, bo
powiadomienie znika po sekundach, a Operator wraca po przebieg rozmowy do
historii i ma tam znaleźć ślad każdej odmowy.

## budowa/klient-poprzedni/src/dostepy/dodanie-katalogu.ts

Drogą pierwszą i właściwą jest natywne okno wyboru powłoki. Ścieżka wraca z niego
istniejąca i rozwinięta, więc punkt zakładany na niej nie wskazuje miejsca, którego
nie ma.

Drogą drugą jest pole ścieżki wpisywanej z ręki. Nie jest ono atrapą ani zapasem na
gorsze czasy: interfejs bywa otwarty w przeglądarce, bez powłoki, a wtedy natywnego
okna po prostu nie ma. Pole zostaje czynne zawsze, także w powłoce, ponieważ Operator
znający ścieżkę nie musi jej odklikiwać.

Przycisk natywnego wyboru nie jest wygaszany poza powłoką. Naciśnięcie daje odpowiedź:
zdanie o tym, że okno systemowe należy do powłoki, oraz przeniesienie uwagi do pola
ścieżki.

Katalog zakładany tą drogą leży na maszynie bieżącej, czyli na tej, na której stoi
powłoka. Schemat wymaga wskazania urządzenia, a jego identyfikator poda komenda
`device.list`, gdy trafi do kontraktu. Do tego czasu drugi argument funkcji
`zalozKatalogLokalny` zostaje pusty i rdzeń odmawia z powodem, bez atrapy
identyfikatora.

## budowa/klient-poprzedni/src/moduly/agents/podglad-wywolania.ts

Podgląd pyta rdzeń o gotowe wywołanie zamiast składać wiersz polecenia po
swojemu. Gdyby okno składało wywołanie własnymi regułami, pokazywałoby wiersz,
którego tura nigdy nie wykona — rdzeń liczy podgląd tymi samymi funkcjami,
którymi wykonuje wysłanie wiadomości, wraz z nałożeniem eksperta okna.

Wywołanie liczy się dla okna komunikacji, a moduł agentów żadnego okna rozmowy
nie zna: okna modułu i okna rozmowy są odrębnymi bytami. Dlatego identyfikator
okna wpisuje Operator, a bez niego podgląd rdzenia nie pyta, zamiast sięgać po
okno przypadkowe.

Pusty prompt systemowy jest odpowiedzią, nie brakiem odpowiedzi: znaczy, że ani
oś, ani ekspert nie wnoszą treści systemowej. Pole zostaje wtedy widoczne wraz
z takim zdaniem, zamiast się schować i sugerować, że odczyt się nie odbył.

## budowa/klient-poprzedni/src/moduly/developer/okno-warsztatu-kodu.ts

Czynności warsztatu stoją w osobnym oknie zamiast w pasku pływającym nad
zaznaczeniem, ponieważ pole edycji Code Editora jest zwykłym obszarem tekstu
bez modelu dokumentu — zaznaczenie znika przy pierwszym kliknięciu poza nim.
Okno podaje ten sam zestaw czynności w postaci, która nie potrzebuje żywego
zaznaczenia: plik bierze się ze stanu modułu, a zaznaczenie wpisuje się
jawnie w polu formularza.

Refaktoryzacja i operacje kontekstowe modelu wracają jako podgląd zmiany, nie
jako zapis na dysk. Zapis jest osobnym rozstrzygnięciem Operatora, bo to
jedyna chwila, w której da się pracę modelu odrzucić przed jej utrwaleniem.

Nawigacja po symbolach i analiza statyczna niosą w odpowiedzi jawne pole
mówiące, czy serwer języka i program analizy statycznej były osiągalne. Okno
rozróżnia stan, w którym wystąpień nie ma, od stanu, w którym nie było czym
ich szukać — pierwszy naprawia się w kodzie, drugi instalacją programu po
stronie serwera.

## budowa/klient-poprzedni/src/mobile/wywolania-interwencji.ts

Moduł składa cztery drogi interwencji warstwy mobilnej w wywołania komend kontraktu rdzenia i jest jedynym miejscem tej warstwy, które zna nazwy komend zmieniających stan sesji.

Każda droga stoi na określonych komendach kontraktu. Zatwierdzenie kroku wywołuje komendę akcji kolejki z rozstrzygnięciem wznowienia albo powtórzenia. Wstrzymanie stosuje akcję kolejki pauzy, zatrzymanie wiadomości albo zatrzymanie sesji, przy czym wstrzymanie sesji oddaje wykaz okien, w których turę zatrzymano. Nastawienie koordynatora zapisuje konfigurację sesji na zasięgu okna i potwierdza wynik odczytem zwrotnym konfiguracji obowiązującej, który wskazuje źródło nadpisania. Przejęcie sterowania składa trzy wywołania: zatrzymanie wiadomości, zapis konfiguracji sesji z trybem uprawnień ręcznym oraz wysłanie wiadomości z poleceniem Operatora.

Droga przejęcia sterowania zatrzymuje turę, przestawia okno na pytanie o każdy krok i wpuszcza polecenie Operatora. Dedykowanej komendy przejęcia kontrakt nie niesie, dlatego złożenie trzech wywołań zastępuje jedno wywołanie, a ekran pozostaje bez zmian.

Droga zatwierdzenia kroku ma widoczną granicę: zatwierdzenie prowadzi wyłącznie przez kolejkę, więc pozycja bez kolejki nie ma czym zatwierdzić kroku, ponieważ komendy zatwierdzenia pojedynczego kroku kontrakt nie niesie. Funkcja rozstrzygająca dostępność dróg wyraża ten stan zdaniem, a ekran nie rysuje wtedy przycisku zatwierdzenia.

Każde wywołanie oddaje kwit zamiast wartości logicznej czy pustego wyniku: nazwę komendy, rozstrzygnięcie, stan po zmianie wyjęty z odpowiedzi rdzenia albo treść odmowy. Kwit jest jedynym dowodem interwencji, jaki pozostaje po zamknięciu aplikacji mobilnej.

## budowa/klient-poprzedni/src/moduly/apps/odmowa-rdzenia.ts

Powód osobnej ścieżki wywołania jest mechaniczny. Komenda bez uchwytu w rdzeniu wraca kopertą typu
`<obszar>.unknown` z polem `requestedType`. Ta koperta nie niesie pola `status`, a korelacja żądań
rozstrzyga wyłącznie koperty ze statusem, co opisuje `protokol/koperta.ts`. Obietnica zwykłego
wywołania pozostałaby więc nierozstrzygnięta na zawsze, a okno stałoby w stanie ładowania bez końca.

Rozwiązanie nie buduje drugiej drogi do rdzenia. Żądanie idzie tym samym `kanal.wyslij`, a odmowę
odczytuje się z dziennika komunikatów nierozpoznanych, który kanał prowadzi dla wszystkich obszarów
kontraktu w `polaczenie/dziennik-nieznanych.ts`. Moduł nie zna literału `apps.unknown` i nie zakłada
własnej subskrypcji zdarzenia odmowy.

Wynik odmowy jest zwykłym `Wynik` z polem `blad`, dzięki czemu okno pokazuje go tak samo jak każdą
inną odmowę rdzenia i nie potrzebuje osobnej gałęzi widoku.

## budowa/klient-poprzedni/src/moduly/design/modal-kreatora.ts

Cztery drogi zamknięcia są równorzędne i wszystkie obowiązkowe. Dwie z nich są
własnością natywnego elementu `dialog`: klawisz Escape obsługuje przeglądarka,
a kliknięcie w nakładkę rozpoznaje się po tym, że celem zdarzenia jest sam
element dialogu, a nie jego wnętrze — wnętrze przykrywa dialog w całości, więc
kliknięcie w treść nigdy nie dochodzi do elementu nadrzędnego. Ten sam wzorzec
niesie `konfiguracja/okno-konfiguracji.ts`.

Przycisk główny stopki zachowuje klikalność zawsze, także przy brakach w polach
i w czasie ładowania. Powtórzone naciśnięcie rozstrzyga logika akcji, a nie
odebranie klikalności kontrolce: kontrolka nieklikalna nie mówi Operatorowi,
czego brakuje, więc blokada należy do warstwy, która zna powód.

Wywołanie `showModal` daje nakładkę, umieszczenie na stosie okien i obsługę
klawisza Escape. W środowisku, które tej metody nie ma, powłoka otwiera okno
zwykłym ustawieniem stanu otwarcia, ponieważ okno ma się pokazać niezależnie od
dostępności metody.

## budowa/klient-poprzedni/src/moduly/apps/zrodlo-izolacji-apps.ts

Granica pojęciowa jest tu istotna i okno musi ją nazwać wprost. Komenda
`isolation.policy.preview` oddaje politykę platformy rozstrzygniętą po ośmiu
poziomach zasięgu, a nie izolację nadaną pojedynczemu rozszerzeniu. Osiem
zakresów technicznych — katalog roboczy, środowisko procesu, dostęp sieciowy,
odczyt i zapis plików, konto i token, model procesu, serwer wykonania oraz
katalog danych modelu — opisuje warunki, w jakich kod się wykonuje, a w tych
samych warunkach wykona się kod rozszerzenia. Deklaracji uprawnień rozszerzenia
kontrakt nie niesie, więc okno mówi o tym wprost, zamiast podstawiać jedno
za drugie.

Poziom zasięgu źródło bierze najwęższy, jaki moduł zna. Podanie okna modułu każe
rdzeniowi rozstrzygnąć dziedziczenie aż do niego, więc odpowiedź opisuje
politykę obowiązującą pracy prowadzonej w tym oknie. Bez okna źródło pyta
o poziom modułu, a wtedy odpowiedź opisuje warunki wspólne wszystkim jego oknom:
zdanie pozostaje prawdziwe, tyle że szersze.

Warstwy żądanie nie podaje. Warstwa pominięta znaczy warstwę obowiązującą,
a wskazanie którejkolwiek innej byłoby cudzym rozstrzygnięciem przebranym
za odczyt stanu.

## budowa/klient-poprzedni/src/moduly/browser/zapisy-przekazania.ts

Moduł skupia ścieżki zapisu wychodzące poza obszar komend `browser.*` i jest
trzymany osobno od odczytu w `zrodlo-browser.ts`, ponieważ dotyczy innych
obszarów kontraktu i ma innego adresata.

Przekazanie międzymodułowe, uruchamiane poleceniami wysłania do innego modułu
oraz tłumaczenia, idzie komendą `context.transfer`. Jest to jedyna droga
przeniesienia kompletu kontekstu między modułami, jaką niesie kontrakt i jaką
obsługuje rdzeń.

Pytanie o zaznaczenie, uruchamiane poleceniem wyjaśnienia na pasku pływającym,
idzie komendą `message.send`. Oknem, do którego pytanie trafia, jest okno modułu
Browser, to samo, którego identyfikator niosą komendy `browser.*`.

Adnotacja spłaszczona do obrazu rastrowego idzie tą samą drogą co pytanie:
`message.send` kładzie treść w rozmowie wskazanego okna i niesie pole
`attachments`, natomiast `context.transfer` wymaga pola `targetModuleId`
i wysłałby rysunek poza moduł, w którym powstał, nie stawiając go w rozmowie
wcale.

Pole `attachments` żądania `message.send` jest polem kontraktu: rdzeń przepisuje
je do zakładanej wiadomości w adapterze rozmowy i oddaje w odpowiedzi, dzięki
czemu okno ocenia dołączenie po wierszu, który wrócił, a nie po tym, że wysłało.

## budowa/klient-poprzedni/src/komponenty/naglowek-okna.ts

Jeden kształt górnego pasa obowiązuje we wszystkich modułach, ponieważ Operator ma otrzymywać tę samą odpowiedź na pytanie o tożsamość oglądanego okna niezależnie od tego, który moduł je otworzył. Granica wobec ramy okna przebiega tak, że rama obejmuje cały kontener okna wraz z ciałem i stanami, natomiast nagłówek odpowiada wyłącznie za górny pas i nie zna treści umieszczonej pod sobą. Nagłówek nie osadza się samodzielnie w drzewie dokumentu, lecz dostaje go rama okna albo plik okna.

Plakietka roli stoi bezpośrednio przy nazwie, ponieważ nazwy okien powtarzają się między modułami: Process Monitor w module Terminal i Execution Monitor w module Automations to dwa odrębne okna. Podpis roli jest tym elementem, który je na ekranie rozróżnia, dlatego napis w plakietce bywa zarówno samym określeniem roli, jak i zdaniem doprecyzowującym przeznaczenie okna.

Wygląd pasa pochodzi w całości z biblioteki: pas stoi na klasie dn-karta-naglowek, nazwa na klasie dn-karta-tytul, a plakietka roli na klasie dn-plakietka--rola, bez klas własnych i bez barw zapisanych w kodzie. Moduł, który potrzebuje odstępstwa, podaje własną klasę polem klasa i styluje nazwę selektorem potomka, zamiast powielać komponent. Wyrównanie kontrolek do prawej krawędzi należy do modułu wołającego, który opakowuje je klasą dn-pasek-prawa.

## budowa/klient-poprzedni/src/moduly/design/okno-tokens-system-panel.ts

Okno Tokens & System Panel stoi jako rozwinięcie warstwy trzeciej, bo opracowanie wywołuje
je z menu kebab obszaru roboczego, a nie stawia go na widoku spoczynkowym. Wartości żetonów
czyta z motywu obowiązującego w danej chwili (`zetony-systemu.ts`), pary i progi kontrastu
z wykazu progów produktu (`kontrast-wcag.ts`), a powiązania żetonu z komponentami — z arkuszy
wczytanych do dokumentu; nic z tego nie jest przepisane do modułu, więc poprawka w motywie
jest widoczna od razu, a rozjazd nazw wychodzi na wierzch zamiast zniknąć.

Okno nie zapisuje zestawu żetonów, bo kontrakt nie zna takiego bytu — wydanie wychodzi
plikiem do przeglądarki Operatora. Nie nadpisuje motywu wczytanym zestawem, bo motyw jest
własnością powłoki, a import kończy się zestawieniem różnicy. Nie wydaje przewodnika stylu
do modułów Library i Studio, bo kontrakt nie zna przekazania tego bytu między modułami.

Przełącznik motywu podglądu nie jest kopią przełącznika powłoki: woła tę samą czynność
motywu produktu, więc drugiego stanu motywu nie ma. Okno nasłuchuje zmiany motywu i przelicza
wszystkie pomiary, bo oba motywy są równoprawne i każdy wymaga własnego pomiaru.

Panel zestawów trwałych czyta żetony motywu w chwili naciśnięcia, a nie z kopii zrobionej
przy otwarciu okna, ponieważ między otwarciem a naciśnięciem Operator mógł przełączyć motyw,
a zestaw ma opisywać stan obowiązujący w chwili działania.

Porównanie wczytanego zestawu z żetonami motywu idzie po nazwach ról, nie po strukturze
pliku: zestaw obcy bywa zagnieżdżony dowolnie, a jedyne, co da się z nim zrobić uczciwie, to
powiedzieć, które role produktu w nim są, które mają inną wartość i których produkt nie zna.

## budowa/klient-poprzedni/src/dostepy/zapisy-nadan.ts

Odczyt i zapis nadań to dwie różne odpowiedzialności. Stan sekcji pilnuje tego, co widok wie
o punktach dostępu i nadaniach, a ten plik pilnuje tego, co dzieje się po zapisie.

Komplet zastępuje zbiór w całości. Wszystkie trzy komendy zwracają nie tylko zmienione nadanie,
lecz komplet nadań okna po zmianie, ponieważ przestawienie kolejności albo oznaczenia głównego
dotyka pozostałych wierszy. Widok bierze komplet zamiast składać go z domysłów.

Nadanie bez okna rozmowy nie idzie do rdzenia, ponieważ nadanie żyje przy oknie i bez okna nie ma
czego nadać. Odmowa przyjmuje kształt wyniku komendy, żeby widok nie potrzebował drugiej ścieżki
obsługi.

## budowa/klient-poprzedni/src/moduly/apps/okno-backend-workspace.ts

Plik jest wiązaniem, a nie drugim widokiem. Formularz warsztatu stoi raz,
w `okno-warsztatu.ts`, natomiast osobny plik daje jedno miejsce w drzewie, po
którym widać, że okno warstwy usług jest zbudowane, a nie tylko wymienione
w słowniku kodów okien.

Kod okna nie pada tu literałem: ramę woła `utworzRameApps(opis.kodOkna, …)` ze
zmiennej, dzięki czemu `KODY_OKIEN` pozostaje jedynym miejscem, w którym kody
okien modułu Apps są zapisane. Okna Studio rozdziela ta sama zasada —
`okno-studio.ts` jest wspólną ramą, a poszczególne okna mają własne pliki.

## budowa/klient-poprzedni/src/moduly/library/tresc-podgladu.ts

Podgląd graficzny oraz podgląd stron dokumentu PDF przychodzą z rdzenia
wyłącznie jako odnośnik w polu `imageRef`. Klient nie ma komendy pobierającej
bajty spod tego odnośnika, więc ciało podglądu pokazuje sam odnośnik i nazywa
brak treści, zamiast udawać obraz, którego nie otrzymał.

Pole `truncated` odpowiedzi mówi o skróceniu podglądu po stronie rdzenia.
Wypisane pod treścią zdanie odróżnia fragment pliku od jego całości, dzięki
czemu odczyt niepełny nie jest brany za odczyt kompletny.

## budowa/klient-poprzedni/src/konfiguracja/kontrolka.ts

Formularz okna konfiguracji powstaje z katalogu pozycji i z założenia nie zna żadnej kontrolki z osobna. Zna wyłącznie wspólny kształt zadeklarowany w tym pliku, co pozwala złożyć ekran z pozycji katalogu bez wiedzy o rodzaju wartości, jaki za nimi stoi.

Konsekwencją takiego rozdziału jest koszt rozszerzenia. Dodanie kolejnego rodzaju wartości do kontraktu sprowadza się do jednego przypadku w module rozdzielającym wybor-kontrolki.ts oraz do jednej funkcji budującej kontrolkę. Nie powstaje przy tym nowy ekran ani zmiana w samym formularzu.

## budowa/klient-poprzedni/src/moduly/assistant/okno-memory-context-manager.ts

Okno zamyka obszar pamięci modułu Assistant i realizuje zasadę jawności: pamięć
asystenta jest w całości widoczna, edytowalna i usuwalna. Cztery pierwsze
zakładki odpowiadają czterem obszarom opracowania modułu — fakty, pamięć
semantyczna, konteksty i baza wiedzy.

Plik odpowiada wyłącznie za skład okna. Wywołania kontraktu mieszkają
w `zrodlo-pamieci.ts`, a każda zakładka ma własny plik obszaru: ustalenia
w `panel-faktow.ts`, wskaźnik znaczenia w `panel-wiedzy.ts`, poziomy pamięci
karty sesji w `panel-kontekstow.ts`.

Fazy odczytu nie ma na poziomie okna, tylko w obszarach. Cztery zakładki czytają
cztery różne byty i wołają je w różnych chwilach — jedna faza dla całego okna
kazałaby odmowie odczytu pamięci przesłonić zakładkę wskaźnika, która o tej
odmowie nic nie wie.

Piąta zakładka „Zestawy i retencja" zamyka trzy obszary, które kontrakt niesie
w całości: nazwane konteksty pamięci `memory.context.*`, zasady retencji
i wygaszania `memory.retention.*` oraz miernik okna kontekstu
`context.usage.get`. Miernik pokazuje liczbę policzoną tokenizatorem rdzenia
wraz z nazwą słownika, którym policzono, a gdy pomiar jest niewykonalny — sam
powód.

## budowa/klient-poprzedni/src/modele/zakladki-sekcji.ts

Cztery obszary sekcji odpowiadają czterem pytaniom o ten sam model: konta rozstrzygają, czym prowadzone jest połączenie, ustawienia osi rozstrzygają, jak model ma działać, tożsamość rozstrzyga, kim model ma być, a podgląd promptu pokazuje, co ostatecznie trafia do modelu. Wszystkie cztery obszary mówią względem tej samej osi, dlatego trzyma je jedna sekcja, a nie cztery osobne sekcje.

Przełącznik nie porzuca obszaru, z którego Operator wychodzi. Element obszaru zostaje w drzewie dokumentu i jedynie przestaje być widoczny, dzięki czemu tekst wpisany w edytorze tożsamości przeżywa zajrzenie do obszaru kont i powraca w niezmienionej postaci.

## budowa/klient-poprzedni/src/moduly/multitasking/kontrolki.ts

Przycisk, pole, wybór, wykaz i pozycja wykazu pochodzą z pliku
`modele/kontrolki-formularza.ts`, a moduł bierze je stamtąd zamiast trzymać
własne kopie. Wiersz klucz–wartość odpowiednika tam nie ma: niosą go wyłącznie
okna ról, gdzie zastępuje tabelę stanu — więź wykonawcy z koordynatorem,
licznik obiegów oraz powód zatrzymania biegu.

Wygląd w całości pochodzi z arkusza rodziny `dm-`, przez co plik nie zna ani
jednej barwy i ani jednego odstępu. Przedrostek klas stoi jedną stałą, bo
arkusz rodziny należy do modułu, a trzy wykazy okien ról muszą wskazywać ten
sam arkusz.

Wysokość pól redakcyjnych jest wspólna dla trzech okien ról, ponieważ kreator
promptu, polecenie wykonawcy i uzasadnienie oceny stoją obok siebie na scenie.

## budowa/klient-poprzedni/src/aod/sekcja-kontekstu.ts

Pole `executionParams` ma w kontrakcie typ `unknown`, więc sekcja go nie rozbiera i melduje jedynie
obecność parametrów wykonania. Rozbiór wymagałby założenia o kształcie danych, którego kontrakt nie
gwarantuje.

Komplet pusty jest stanem poprawnym, a nie odmową odczytu, ponieważ wszystkie pola `ContextBundle`
są w kontrakcie opcjonalne. Sekcja mówi to osobnym zdaniem zamiast wypisywać wykaz samych oznaczeń
braku.

## budowa/klient-poprzedni/src/mobile/pasek-kwitu.ts

Pasek czyta kwity z pamięci trwałej urządzenia, a nie ze stanu okna.
Potwierdzenie trzymane w stanie okna ginie razem z oknem, natomiast kwit zapisany
w magazynie przeżywa zamknięcie ekranu i ponowne uruchomienie telefonu, dzięki
czemu wykaz wysłanych decyzji pozostaje kompletny po powrocie do aplikacji.

Wiersz paska powstaje ze zdania złożonego przez `zdanieKwitu` z odpowiedzi
rdzenia, nie z zamiaru klienta. Kwit nieudany stoi na tym samym pasku wraz
z treścią odmowy, ponieważ odmowa jest rozstrzygnięciem równie wiążącym jak
przyjęcie i ma być widoczna w tym samym miejscu.

## budowa/klient-poprzedni/src/mission-control/naglowek-pulpitu.ts

Nagłówek ma jedną odpowiedzialność: podaje tytuł ekranu i jawnie nazywa
pochodzenie liczb pokazywanych niżej. Dopóki odczyt z rdzenia nie nadszedł,
a więc dopóki źródło danych ma wartość `ZrodloDanych.Oczekiwanie`, przy tytule
stoi plakietka oczekiwania wraz ze zdaniem wyjaśniającym. Zdanie jest konieczne,
ponieważ bez niego pusty ekran zostałby wzięty za pomiar mówiący, że pracy
nie ma.

Po pierwszym odczycie plakietka mówi o danych pochodzących z rdzenia. Pulpit nie
zna innych źródeł: liczby albo pochodzą z odczytu rdzenia przez kanał kontraktu,
albo nie ma ich wcale, a sekcje pokazują stany puste.

## budowa/klient-poprzedni/src/moduly/library/wykaz-plikow.ts

Struktura katalogów pochodzi z danych, nie z osobnego zapytania: kontrakt nie ma
komendy katalogu folderów, niesie za to pole ścieżki każdego pliku. Wykaz
katalogów składa się więc z pierwszych członów ścieżek zwróconych przez rdzeń,
a zbiór pusty daje jedną pozycję obejmującą cały zbiór.

Formę prezentacji rozstrzyga widok wybrany w oknie i zbudowany w module widoków
wykazu. Wykaz nie zna żadnej z pięciu form: składa katalogi, oddaje zbiór
widoczny i osadza to, co widok zbudował.

## budowa/klient-poprzedni/src/moduly/multitasking/plan-etapow.ts

Etap planu jest kolejką, a nie osobnym bytem klienta: drugiego rejestru etapów po stronie
okna nie ma, więc wykaz odświeża wyłącznie zdarzenie `queue.changed` przychodzące z rdzenia.

Kontrolka zależności etapów pozostaje widoczna i nieczynna. Komenda
`orchestration.dependency.set` wiąże kroki układów automatyk polami `fromStepId` oraz
`toStepId`, a etap tego planu jest kolejką, nie krokiem automatyki. Powód nieczynności bierze
się z wykazu komend oddanego przez rdzeń, nie ze stałej wpisanej w kodzie okna.

Potwierdzenie założenia etapu składa się z pól odpowiedzi rdzenia, nie z pól żądania:
przycięta nazwa albo nieprzyjęte okna wykonawców mają być dla Operatora widoczne.

Wiersz etapu bieżącego niesie plakietkę zamiast wygaszonego przycisku, a kontrolka wyboru
pojawia się wyłącznie w wierszach, w których jest co wybrać.

## budowa/klient-poprzedni/src/moduly/library/tresc-base64.ts

Odczyt jest wspólny dla wgrania pliku komendą `library.file.upload` oraz dla
dołożenia wersji komendą `library.version.add`, ponieważ obie komendy niosą to
samo pole `contentBase64`.

Nieudany odczyt, czyli zniknięcie pliku albo odrzucenie dostępu przez
urządzenie, wraca jako `null`, a nie wyjątkiem wywracającym widok. Wołający ma
wtedy powiedzieć Operatorowi, że nic nie zostało wysłane.

## budowa/klient-poprzedni/src/moduly/apps/rama-okna.ts

Obudowa okna, na którą składają się nagłówek, plakietka roli, pas akcji i ciało, należy do komponentu ramy okna w katalogu komponenty, a znakowanie fazy do komponentu fazy okna w tym samym katalogu. Plik ramy modułu Apps splata oba byty, ponieważ rama biblioteczna nie zna cyklu życia danych, których sama nie pobiera, a pięć okien modułu ma mieć jedną przesłonę stanu zamiast pięciu wariantów tego samego rozwiązania.

Przesłona stanu nie kasuje treści, tylko ją przykrywa. Powrót do fazy gotowej odsłania to, co Operator już widział, dzięki czemu nieudane odświeżenie nie zabiera wyniku poprzedniego odczytu. Wskaźnik odczytu stoi obok komunikatu, a nigdy zamiast kontrolki, więc treść okna zostaje na miejscu i pozostaje klikalna. Wspólna procedura znakowania fazy tego nie przesądza, ponieważ moduły różnią się między sobą pod tym względem.

Wygląd pochodzi w całości z biblioteki komponentów z przedrostkiem dn oraz z żetonów motywu, dlatego plik nie zawiera ani jednej wartości barwy i ani jednej wartości odstępu.

## budowa/klient-poprzedni/src/aplikacja/widok-pulpitu.ts

Pulpit operacyjny stoi obok strony głównej, a nie zamiast niej. Centrum
dowodzenia jest wejściem przy rozpoczynaniu pracy, pulpit — przy powrocie
do pracy już rozpoczętej, a obydwa pozostają dostępne z przełącznika widoków.

Dane pulpitu pochodzą z rdzenia. Źródło pulpitu odpytuje `session.list`,
`window.list` oraz `channel.list`, a następnie subskrybuje zdarzenia
`session.changed`, `window.changed`, `queue.changed` i `progress.changed`.
Do pierwszego odczytu pulpit pokazuje stany puste, ponieważ liczby wymyślonej
nie poda.

Zamiary Operatora idą do rdzenia tym samym kanałem, którym przyszły dane:
`session.open` oraz `queue.action` przechodzą kanałem kontraktu.

## budowa/klient-poprzedni/src/aplikacja/stan-pary-gniazda.ts

Stan pary ma dwa człony. Skład pary mówi, kto jest koordynatorem, a kto
wykonawcą; położenie pary mówi, czy wykonawca pracuje, czy kolejka stoi.
Pierwszy człon niesie zdarzenie `window.changed`, drugi — `queue.changed`.

Właścicielem pola `windowRole` jest rdzeń, nie scena. Po odtworzeniu okna
z bazy, po `window.update` z innego urządzenia albo po `context.transfer` może
wrócić rola inna niż domyślna nadana przy tworzeniu gniazda. Dlatego każde
`window.changed` dotyczące tego okna przestawia nagłówek, a zmianę na tę samą
rolę gniazdo pomija, żeby widok nie przebudowywał się bez powodu.

Okno poza pętlą nie ma pary, więc jego nagłówek zostaje bez plakietki stanu
zamiast pokazywać stan cudzej kolejki.

Każde gniazdo słucha całej magistrali zdarzeń, więc zdarzenie dotyczące innego
okna jest pomijane. Kolejka bez wykazu okien dotyczy całej sesji i wchodzi do
każdego gniazda tej sceny.

## budowa/klient-poprzedni/src/modele/okno-modeli.ts

Okno stoi na natywnym elemencie okna dialogowego: warstwę tła, stos okien
i zamknięcie klawiszem ucieczki daje przeglądarka, a nie własna nakładka. Wygląd
bierze z biblioteki komponentów, z klasy okna modalnego, więc plik nie ustala
barw.

Rama wyłącznie osadza sekcję. Treść i stan mieszkają w module sekcji modeli, aby
tę samą sekcję dało się wstawić także w widok osadzony bez powielania kodu.

## budowa/klient-poprzedni/src/modele/zrodlo-tozsamosci.ts

Kategorie zasad są danymi, nie kodem: konstytucja, profil roli, ekspertyza,
zasady bezpieczeństwa i zasady harnessu przychodzą komendą
`identity.category.list` wraz z warstwą, porządkiem, obowiązkowością i trybem
proponowanym. Dopisanie kategorii jest wtedy nowym wierszem katalogu, a nie
zmianą kodu klienta.

Nakładka obowiązująca jest odczytem rdzenia, nie sklejeniem w kliencie: klient
nie składa promptu z treści kategorii własnym porządkiem, tylko pyta
`identity.effective.get` i pokazuje to, co trafi do modelu. Drugie składanie po
stronie widoku dałoby podgląd rozjeżdżający się z rdzeniem przy pierwszej
zmianie reguł.
