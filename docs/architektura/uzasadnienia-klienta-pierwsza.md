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

## budowa/klient-poprzedni/src/moduly/design/symulacja-widzenia.ts

Zakres symulacji jest tu ważniejszy niż sam rachunek: obejmuje ona próbki
żetonów, a nie obrazy modułu. Obrazu nie ma czym pobrać, ponieważ pole odsyłacza
zasobu niesie ścieżkę w systemie plików rdzenia, więc nałożenie symulacji na
zasób wizualny nie ma dziś drogi wykonania i panel tego nie udaje.

Rachunek jest przybliżeniem i tak jest nazwany. Macierze odwzorowują trzy
dichromazje wprost w przestrzeni sRGB, bez przejścia przez przestrzeń
długofalową, co jest uproszczeniem przyjętym w narzędziach projektowych.
Wystarcza ono, żeby zobaczyć, które dwie barwy systemu zlewają się w jedną,
i nie wystarcza do orzeczenia medycznego. Panel mówi to Operatorowi wprost,
ponieważ różnica między podglądem a badaniem jest tutaj istotna.

Granica zlania jest odległością w składowych sRGB — miarą zgrubną, dobraną tak,
żeby wskazywała pary wymagające obejrzenia przez człowieka, a nie żeby
rozstrzygała za niego.

## budowa/klient-poprzedni/src/moduly/apps/dymek-objasnienia.ts

Wiersz `mp-pole-z-dymkiem` jest wstawką układu właściwą modułowi: pole rośnie na
całą dostępną szerokość, a znak zapytania stoi przy jego prawej krawędzi.

Klasy własne idą do wspólnej fabryki parametrem, ponieważ znak modułu nie jest
bibliotecznym przyciskiem ikonowym. Klasa `mp-dymek__znak` opisuje obwódkę
o wymiarze `--dn-wym-ikona-sm` i narożniku `--dn-r-pill` ze wskaźnikiem myszy
`help`, natomiast biblioteczna klasa `dn-btn-ikona` opisuje kwadrat o boku
trzydziestu dwóch pikseli. Przekazanie klas parametrem pozwala zachować obie
postacie bez rozgałęziania fabryki.

## budowa/klient-poprzedni/src/dostepy/elementy-karty.ts

Elementy kart sekcji dostępów powstają w jednym miejscu, żeby karta punktu,
wiersz nadania i obszar dodawania katalogu miały spójny wygląd. Plik nie zna
barw, odstępów ani reguł widoku: całość wyglądu niosą klasy biblioteki
`komponenty/` oraz klasa modyfikująca podana przez wywołującego. Plakietka
znaku przyjmuje klasę biblioteki i klasę miejsca osobno, ponieważ pierwsza
opisuje wygląd, a druga położenie w układzie karty.

## budowa/klient-poprzedni/src/moduly/assistant/stan-okna.ts

Plik odpowiada wyłącznie za nośnik komunikatu stanu wraz z miejscem na treść okna. Nazwę fazy oraz znakowanie drzewa dokumentu wnosi komponent fazy okna z katalogu komponenty, więc tutaj zostaje jedynie to, czym moduł Assistant różni się od reszty drzewa: wskaźnik odczytu stoi obok treści, a treść pozostaje widoczna również w fazie ładowania, dzięki czemu kontrolki są klikalne przez cały czas trwania wywołania.

Brak zapytania do rdzenia nie jest fazą okna, lecz fazą źródła danych, którą niesie zapis modułu w polach opisujących zadane pytania. Okno pokazuje ten przypadek jako fazę pustą z własnym zdaniem wyjaśniającym. Stan nie zastępuje treści, tylko ją przesłania, więc po powrocie do fazy gotowej wcześniejsza treść pozostaje nietknięta, a nieudane odświeżenie nie kasuje tego, co było już na ekranie.

## budowa/klient-poprzedni/src/dostepy/ostrzezenie-zapisu.ts

Specyfikacja mostu mcp-danaco-pulpit-console zabrania nadawania trybu zapisu na maszynie danaco-data bez wyraźnej potrzeby. Maszyna ta jest hostem produkcyjnej platformy LEX, a zapis modelu sięga tam zbiorów, z których korzysta cała kancelaria. Zakazu nie egzekwuje jednak kod klienta: nadanie zapisu pozostaje możliwe, ponieważ rozstrzyga o nim Operator. Egzekwuje go jawne, widoczne zdanie postawione przy przełączniku trybu, a nie podpowiedź ukryta pod kursorem.

Wykaz maszyn chronionych stoi w kodzie, a nie w katalogu konfiguracji, ponieważ jest to ostrzeżenie bezpieczeństwa i nie może dać się wyłączyć zapisem w bazie. Rozszerzenie wykazu sprowadza się do jednego wiersza w stałej MASZYNY_CHRONIONE.

## budowa/klient-poprzedni/src/moduly/design/kanaly-obrazowe.ts

Żądanie `design.asset.generate` niesie pole `channelId`, a rdzeń sprawdza rodzaj
wskazanego kanału zanim cokolwiek wyśle — kanał tekstowy odmawia kodem
`validation_failed` w `core/adapter_modul_design_kanal.go`. Wykaz silników
pokazuje więc kanały obrazowe jako wybieralne, a pozostałe wymienia z nazwy
i z powodem, zamiast je ukrywać.

Kolejność rozpoznania jest przepisana z funkcji `KluczAdaptera()`
w `models/definicja.go`: parametr `adapter` ma pierwszeństwo, a rodzaj wiersza
jest wartością zapasową. Do kontraktu obie strony jadą tym samym wierszem
(`core/adapter_kanaly.go`, pola `Kind: k.RodzajKanalu` oraz
`Config: ParametryJSON`), więc `Channel.kind` odpowiada polu `Rodzaj`,
a `Channel.config.adapter` polu `Parametr`. Kluczem adaptera obrazowego jest
napis `obrazy` — stała `AdapterObrazy` w `models/kanal.go`.

Rozpoznanie może rozejść się z rdzeniem: klient czyta `config` jako treść
nieokreśloną (`config?: unknown` kontraktu), więc kanał o parametrach zapisanych
inaczej niż obiektem trafi między nieobrazowe. Ostateczną kontrolę wykonuje
rdzeń — dlatego wskazanie kanału spoza wykazu obrazowego nie jest w oknie
blokowane, tylko opisane.

## budowa/klient-poprzedni/src/moduly/browser/indeks.ts

Moduł nie osadza się sam w dokumencie: oddaje element, a warstwa składająca
rozstrzyga, gdzie go postawić. Dzięki temu te same okna wchodzą i w obszar
roboczy powłoki, i w podgląd sprawdzianu.

Wytwórnia oddaje moduł wprost, ponieważ `ModulBrowser` niesie pola `element`
i `wczytaj` żądane przez powłokę oraz czynność `rozlacz`, po którą sięga
sprawdzian modułu.

## budowa/klient-poprzedni/src/konfiguracja/kontrolki-liczbowe.ts

Granice i skok pola liczbowego pochodzą z katalogu ustawień, a nie z kodu
klienta. Katalog, który granic nie podaje, daje pole bez granic: brak metadanej
nie jest błędem i niczego nie blokuje, więc kontrolka nie wymyśla ograniczenia,
którego rdzeń nie postawił.

Pole puste znaczy brak wartości, a nie zero. Odczyt oddaje wtedy `null`, przez
co zapis czyści wartość zamiast wpisywać liczbę, której Operator nie podał.
Rozróżnienie jest konieczne, ponieważ zero bywa wartością obowiązującą
i podstawienie go pod pustkę zmieniałoby ustawienie bez wiedzy Operatora.

## budowa/klient-poprzedni/src/mission-control/pasek-dzialan.ts

Pasek ostatniego działania ma jedną odpowiedzialność: wypisuje, co pulpit właśnie
nadał powłoce. Żaden przycisk pulpitu nie jest wyszarzony, a przycisk czynny bez
widocznego skutku nie mówi Operatorowi, czy zamiar poszedł dalej. Pasek pokazuje
zatem fakt nadania zamiaru, a nie wynik jego wykonania.

Stąd bierze się dobór słów w treści paska. Zdanie mówi „zamiar nadany", nigdy
„wykonano", ponieważ w chwili wypisania powłoka mogła jeszcze nie wysłać komendy,
a rdzeń mógł jej nie potwierdzić.

## budowa/klient-poprzedni/src/ikony/marka.ts

Znak czyta się jako podwójny grot z kropką sygnału ustawioną na linii bazowej. Godło uproszczone
niesie jeden grot i obowiązuje od szesnastu pikseli boku w dół, ponieważ poniżej tej wartości drugi
grot oraz prześwit między grotami przestają być czytelne.

Odmianę znaku dobiera się do podłoża, a nie do motywu interfejsu, ponieważ barwy znaku są wpisane
w plik źródłowy. Pasek kokpitu pozostaje atramentowy w obu motywach, więc leży na nim odmiana
przeznaczona na podłoże ciemne, niezależnie od tego, który motyw jest czynny.

## budowa/klient-poprzedni/src/moduly/apps/warsztaty-apps.ts

Opisy warsztatów stoją osobno od ramy warsztatu, ponieważ rama jest czynnością
budującą okno, a opisy są wykazem danych. Wykaz odwzorowuje warstwy wyliczenia
`AppWorkspaceLayer` z kontraktu, więc trzeci warsztat może powstać dopiero po
rozszerzeniu kontraktu, a nie po zmianie samego widoku.

## budowa/klient-poprzedni/src/aplikacja/polaczenie-z-rdzeniem.ts

Pole nawigacji platformy niesie komendy wejścia do środowiska oraz wykazu
modułów osadzone na wspólnym kanale i służy wyłącznie powłoce: powłoka powstaje
bez połączenia i nie ma jak sięgnąć po kanał sama, więc dostaje te dwie drogi
gotowe. Widoki mające kanał wołają opakowania wprost, widok strony głównej
komendę wejścia na stronę główną, a przestrzeń modułu komendę wejścia do
przestrzeni pracy, bez obiektu pośredniego.

Plik wyłącznie składa warstwy: nie zna ramki, nie buduje koperty i nie zna nazwy
żadnej komendy. Nazwy pochodzą z pakietu wspólnego i żyją w warstwie protokołu.

Złożenie nie otwiera połączenia. Rozpoczęcie łączności należy do cyklu życia,
żeby moment jej nawiązania był jednym miejscem, a nie skutkiem ubocznym budowy
obiektów.

Token sesji bramki jest czytany domknięciem przy każdym powitaniu, a nie raz
przy składaniu, ponieważ po ponownym nawiązaniu połączenia obowiązuje sesja
bieżąca, nie ta sprzed zerwania.

Nawigacja powłoki składa się w tym pliku, a nie w uzgodnieniu, ponieważ
uzgodnienie odpowiada wyłącznie za powitanie, sesję i okno, i o nawigacji nic
nie wie. Korzeń montażu klienta jest jedynym punktem znającym jednocześnie kanał
i odbiorcę tych dróg.

## budowa/klient-poprzedni/src/moduly/apps/okno-architecture-designer.ts

Okno kreatora modułu Apps jest punktem wyjścia procesu tworzenia architektury. Definiowanie komponentów rozwiązania i ustalenie zależności stoją na jednej komendzie kontraktu, która przyjmuje szablon i komplet komponentów wraz z ich zależnościami. Przycisk walidacji uruchamia tę samą komendę, nie osobną drogę: kontrakt nie ma osobnej komendy walidacji, zwraca za to zastrzeżenia w architekturze po zapisie. Przycisk pozostaje klikalny przy pustej kanwie, ponieważ odpowiedź jest opisowa, nie blokadą. Zdanie potwierdzające bierze treść z odpowiedzi rdzenia, nie z zamówienia: wymienia identyfikator, wersję, szablon i komponenty tak, jak oddał je rdzeń, a każdą rozbieżność z zamówieniem i każde zastrzeżenie walidacji ogłasza odmową.

Odczyt architektury zapisanej w rdzeniu nadpisuje kanwę i mówi o tym wprost: architektura oddana przez rdzeń wchodzi wraz ze swoim kompletem komponentów, więc komponenty zestawione na kanwie, a jeszcze niezapisane, po tym kroku znikają. Zdanie wymienia liczbę komponentów, która przyszła. Brak architektury jest odpowiedzią, nie odmową — rdzeń oddaje wtedy wynik udany bez pola, a kanwa zostaje nietknięta.

Funkcja porównująca oddany układ z zamówionym zestawia zbiory identyfikatorów, nie listy. Rdzeń trzyma komponenty i krawędzie zależności w osobnych tabelach i oddaje je w porządku własnym, a powtórzoną krawędź zdejmuje kluczem pierwotnym. Porządek i powtórzenie nie są więc rozbieżnością, tylko odmiennym sposobem przechowywania danych; rozbieżnością jest komponent, który zniknął, albo taki, którego nikt nie zamawiał. Zależność jest wskazaniem identyfikatora innego komponentu, więc jej zgubienie zmienia układ tak samo jak zgubienie komponentu, a kanwa zaciągnięta z odpowiedzi sama z siebie nie pokazuje, że zamówiono więcej.

Zdanie o pustym polu zastrzeżeń walidacji orzeka o odebranej ramce, nie o zachowaniu rdzenia. Rdzeń tego pola nie wylicza: adapter czyta zastrzeżenia z bazy, zapisuje je z powrotem tym samym zapisem scalającym i oddaje w odpowiedzi, więc jedynym pisarzem kolumny jest ten, kto ją przed chwilą przeczytał. Zdanie o braku zastrzeżeń byłoby zapewnieniem o sprawdzeniu, którego nikt nie wykonał. Własny walidator w oknie byłby drugą prawdą o tym samym układzie i zacierał granicę odpowiedzialności: zastrzeżenia widoczne w oknie sugerowałyby, że zna je rdzeń. Rozróżnienie pola nieobecnego od pola obecnego i pustego sprawia, że w chwili, w której rdzeń zacznie zastrzeżenia liczyć, zdanie zmieni treść bez tknięcia tego pliku.

## budowa/klient-poprzedni/src/aplikacja/korzen-dokumentu.ts

Całym dokumentem zajmuje się `gospodarz-dokumentu.ts`, a kształt
`KorzenAplikacji` jest wejściem funkcji `zamontujUkladOkien` z katalogu
`okna-rownolegle/`. Korzeń pozostaje przez to warstwą montażową bez wiedzy
o rdzeniu i bez wiedzy o układzie okien.

Miejsce akcji nie powstaje w korzeniu — przychodzi z paska górnego powłoki
środowiska, żeby uchwyt szuflady sterowania stał w pasku wspólnym z resztą
akcji, a nie w osobnym pasku założonym dla jednego przycisku.

Klasa `dn-sesja` należy do karty sesji w pasie powłoki
(`powloka/karta-sesji.ts`), a jej arkusz zawęża szerokość do pasa nawigacji bez
zawężenia selektora rodzicem. Korzeń sceny nosi własną nazwę `dn-scena-sesji`,
ponieważ w przeciwnym razie wpadłby pod tamtą regułę i zszedł do szerokości
słupka nawigacji.

## budowa/klient-poprzedni/src/moduly/library/indeks.ts

Moduł nie osadza się sam w dokumencie i nie zna powłoki: oddaje element,
a warstwa składająca rozstrzyga, gdzie go postawić.

Kod `library` odpowiada kolumnie `modul.kod` w rdzeniu, a wpis do rejestru
modułów wiąże widok z modułem właśnie po tym kodzie.

## budowa/klient-poprzedni/src/moduly/agents/zlecenia-agentow.ts

Okno pracuje polami formularza: napisami z kontrolek i wartościami list wyboru.
Kontrakt pracuje polami opcjonalnymi, których pusta wartość ma znaczyć „nie
podano", a nie „podano pustkę". Struktury zleceń są granicą między jednym
a drugim: pole puste zostaje w oknie, a do rdzenia idzie żądanie bez niego.

Poziomy pamięci jadą zawsze, także jako zbiór pusty. Kontrakt rozstrzyga to
wprost: pominięcie pola znaczy „bez zmiany", a lista pusta znaczy „pamięć
wyłączona w całości" — i jest jedynym zapisem wyłączenia, ponieważ `MemoryLevel`
piątej wartości nie ma. Gdyby okno pomijało pole przy wszystkich poziomach
odznaczonych, Operator odznaczyłby cztery pola i nie wyłączyłby niczego.

Odstępstwo od globalnego promptu systemowego jest wartością logiczną, a nie
wartością `IdentityMode`, ponieważ dwie wartości pola `Agent.mode` nie są
równorzędne. Prompt systemowy ustawiany globalnie w oknie konfiguracji
obowiązuje domyślnie, a moduł Agents daje albo instrukcję dopisywaną do niego
(stan domyślny, wartość fałszywa), albo jawne odstępstwo (wartość prawdziwa), po
którym instrukcja eksperta staje się promptem systemowym. Przekład na pole
`mode` wykonuje `zrodlo-agentow.ts`, na granicy formularza i kontraktu.

## budowa/klient-poprzedni/src/moduly/assistant/wiersz-pamieci.ts

Wiersz pamięci nie prowadzi żadnego wywołania samodzielnie: zamiar oddaje oknu, które prowadzi wywołania rodziny memory. Ta sama granica obowiązuje wiersz zlecenia w pliku wiersz-zlecenia.ts, dzięki czemu oba wiersze pozostają wymiennymi elementami wykazu, a odpowiedzialność za wywołania skupia się w oknie.

Wpis o pochodzeniu model jest propozycją czekającą na decyzję Operatora. Przyjęcie zapisuje tę samą treść z pochodzeniem operator, a odrzucenie usuwa wpis. Rozróżnienie pochodzenia należy do kontraktu przez typ MemoryEntryOrigin, a nie do okna: bez niego pamięć zapisana przez model byłaby nie do odróżnienia od ustalenia wprowadzonego przez Operatora.

## budowa/klient-poprzedni/src/moduly/agents/prowenancja-rady.ts

Rada doradcy jest jawna i nie wolno pokazywać jej jako własnej odpowiedzi eksperta. Gdyby zdanie o pochodzeniu rady składał widok, pierwszy widok, który by tego zaniechał, pokazałby radę bez źródła. Dlatego zdanie o pochodzeniu powstaje jeden raz w zapisie prowenancji, a widok pobiera treść rady razem z tym zdaniem z jednego miejsca.

Kontrakt nie zawiera pola prowenancji ani komendy zapisującej ją w bazie, więc prowenancja żyje dokładnie tyle, co widok, a wykaz konsultacji mówi o tym wprost. Przeniesienie rady do instrukcji eksperta zabiera ze sobą etykietę oraz oba zdania prowenancji, dzięki czemu po wklejeniu nadal widać, skąd akapit pochodzi.

## budowa/klient-poprzedni/src/konfiguracja/indeks.ts

Okno konfiguracji jest jedno na klienta i trwa między otwarciami z tego samego
powodu, dla którego router nie porzuca widoku opuszczonej trasy: powtórne
otwarcie wraca do kategorii, na której poprzednie się zakończyło, zamiast
zaczynać od początku. Trwałość okna utrzymuje przy okazji subskrypcję zdarzenia
`config.changed`, dzięki czemu zapis dokonany w innym miejscu produktu dociera
również wtedy, gdy okno pozostaje zamknięte.

Kanał podaje się przy pierwszym otwarciu. Wywołanie z innym kanałem, właściwe
dla ponownego połączenia z rdzeniem, buduje okno od nowa, aby komendy nie szły
przez transport, którego już nie ma.

## budowa/klient-poprzedni/src/mobile/kwit-decyzji.ts

Telefon jest kanałem interwencji, a nie miejscem pracy: Operator podejmuje jedną decyzję
i odchodzi od urządzenia. Potwierdzenie wypisane w oknie ginie razem z oknem, a pytanie
o wynik wraca później, bez dostępu do pulpitu. Kwit leży w pamięci trwałej przeglądarki,
więc pierwsze, co telefon pokazuje po ponownym otwarciu, to zdanie o wykonanej interwencji
wraz z godziną, nazwą pozycji i stanem potwierdzonym przez rdzeń.

Kwit niesie zdanie rdzenia, nie widoku. Każdy krok zapisuje nazwę komendy, rozstrzygnięcie
wywołania i stan po zmianie wyjęty z odpowiedzi rdzenia — stan kolejki, zatrzymanie, źródło
konfiguracji. Kwit nieudany zapisuje treść odmowy, bo niepowodzenie także jest wiadomością,
na którą Operator czeka.

Magazyn jest wstrzykiwany, wzorem `uwierzytelnienie/sesja-bramki.ts`: domyślnie pamięć trwała
okna, a w sprawdzianie atrapa. Błąd pamięci nie zatrzymuje drogi interwencji — brak magazynu
znaczy, że kwitu nie da się zapisać, a nie że decyzja się nie wykonała. Stąd kolejność:
najpierw wywołanie rdzenia, potem kwit.

Urządzenie bez pamięci trwałej pokazuje wyłącznie potwierdzenie bieżące, a po zamknięciu
telefonu go nie ma. Taki jest stan faktyczny takiego urządzenia i okno go nie ukrywa.

## budowa/klient-poprzedni/src/moduly/design/wczytanie-pliku.ts

Treść pliku czyta klient, a nie rdzeń. Żądanie `design.asset.upload` zna
wprawdzie pole `sourcePath`, lecz ścieżka każe otworzyć plik rdzeniowi, a ten
stoi na innej maszynie niż przeglądarka, w której Operator plik wskazał.

Wymiary mierzy się wyłącznie dla rastra. Wektor oraz obraz, którego przeglądarka
nie zdekoduje, otrzymują zera znaczące brak pomiaru; oba pola wymiaru są
w kontrakcie opcjonalne, więc zasób bez nich jest poprawny i wgranie ma się
odbyć mimo nieudanego pomiaru.

Bajty odczytuje `readAsDataURL`, ponieważ oddaje zapis base64 bez ręcznego
przepisywania bajtów przez `btoa`, które na treści binarnej wymaga przejścia
przez ciąg znaków jednobajtowych i zawodzi na pierwszym bajcie powyżej 0xFF.

## budowa/klient-poprzedni/src/aplikacja/akcje-tras.ts

Powłoka środowiska ma własny pasek górny wraz z przełącznikiem motywu,
powiadomieniami i awatarem. Grupa akcji tras nie zakłada więc drugiego paska,
lecz dokłada się do grupy akcji paska istniejącego, przy jego prawej krawędzi,
przed akcjami samej powłoki. Grupa niesie tylko to, czego pasek powłoki nie ma:
oznaczenie trasy bieżącej i stan łączności z rdzeniem.

## budowa/klient-poprzedni/src/modele/zrodlo-modeli.ts

Kontrakt nie ma komendy `model.list`, ponieważ model nie jest bytem
rejestrowanym, tylko wartością danych, którą niesie kanał modelu
(`Channel.model`) albo konto (`Account.defaultModel`). Wykaz identyfikatorów
powstaje ze złożenia tych dwóch źródeł.

Wykaz służy wyłącznie jako podpowiedź do pola tekstowego. Model jeszcze
nieużywany wpisuje się identyfikatorem wprost, ponieważ pole nie jest listą
zamkniętą.

Rejestr kanałów, który nie dotarł, daje podpowiedź pustą, a nie pusty formularz.
Adresowanie osi modelu działa wtedy dalej, a nieudany odczyt zostaje odnotowany
w dzienniku przeglądarki zamiast zatrzymywać pracę.

## budowa/klient-poprzedni/src/moduly/apps/tabela-wdrozen.ts

Pusty stan wchodzi w miejsce ciała tabeli, a nie zamiast całej tabeli. Nagłówek
kolumn zostaje, ponieważ Operator ma widzieć, czego wykaz dotyczy, zanim
cokolwiek się w nim znajdzie.

Błąd wdrożenia znakuje plakietka w kolumnie stanu, bez zmiany struktury wiersza.
Stan nigdy nie opiera się na samej barwie: plakietka niesie słowo.

Gdy silnik wdrożeń wpisze do pola `logRef` zdanie o tym, dlaczego przebieg się nie
powiódł, tabela pokazuje jego treść tak, jak przyszła, i nazywa pole, z którego
pochodzi. Samo słowo `failed` zostawiałoby Operatora bez wyjaśnienia, które
przyszło w tej samej ramce.

## budowa/klient-poprzedni/src/moduly/agents/indeks.ts

Powłoka zna z tego katalogu jedną czynność: wytwarza moduł wytwórnią
`utworzModulAgents`, przyjmując kanał rdzenia, stawia oddany element w obszarze
roboczym i wywołuje czynność wczytania z identyfikatorem sesji. Kolejność jest
wiążąca, ponieważ wczytanie żąda danych sesji, a element musi już stać
w dokumencie, gdy odpowiedź wraca.

Moduł nie osadza się sam w dokumencie i nie zna powłoki, dzięki czemu te same
okna wchodzą zarówno w obszar roboczy powłoki, jak i w podgląd sprawdzianu.

Kontekst testowanego eksperta jest wystawiony w interfejsie katalogu, ponieważ
sięga po niego powłoka, a nie moduł. Zmianę testowanego agenta rozgłasza plik
`aplikacja/ulotnosc-okna.ts`, na którym stoi okno rozmowy sesji.

## budowa/klient-poprzedni/src/aplikacja/zamiary-pulpitu.ts

Pulpit sam niczego nie wysyła: nadaje zamiar wyrażony nazwami kontraktu, a kopertę
komendy składa dopiero ta warstwa. Dzięki temu pulpit zna wyłącznie słownik zamiarów,
a wiedza o komendach i o ich wymaganych polach zostaje w jednym miejscu.

Dwa zamiary nie mają dziś drogi przez kontrakt i żaden z nich nie jest udawany.
Podniesienie priorytetu nie ma odpowiednika ani w `QueueAction`, ani wśród komend.
Przekazanie kontekstu wymaga kompletu `ContextBundle`, którego zamiar kolejki nie
niesie, a którego nie wolno złożyć po stronie klienta z wartości domyślnych.

W obu przypadkach warstwa wystawia komunikat mówiący wprost, czego brakuje, zamiast
wysyłać komendę niekompletną albo milczeć.

## budowa/klient-poprzedni/src/mission-control/etykiety-pulpitu.ts

Plik ma jedną odpowiedzialność: przekłada wartości kontraktu na słowo widoczne
dla Operatora. Widok nie zna literałów kontraktu i o nazwę pyta wyłącznie tutaj,
dzięki czemu zmiana słownictwa pulpitu nie sięga do samych widoków.

Każdy słownik jest rekordem zupełnym po dziedzinie swojego typu. Dopisanie
wartości do kontraktu przerywa więc kompilację tego pliku, i jest to skutek
zamierzony: nowy stan ma dostać polską nazwę, zamiast wypaść z widoku po cichu
albo pokazać się Operatorowi surowym literałem kontraktu.

Etykieta braku źródła zastępuje liczbę wszędzie tam, gdzie kontrakt miary dziś
nie niesie. Pulpit nie pokazuje wartości wymyślonej — brak źródła danych jest
nazywany wprost.

## budowa/klient-poprzedni/src/konfiguracja/stany-odczytu.ts

Komunikat blokowy nie znika po naciśnięciu przycisku ponowienia. Naciśnięcie wyzwala kolejny odczyt
i zostawia zdanie na miejscu, dopóki sytuacja nie ustanie; zdejmuje je dopiero odczyt zakończony
powodzeniem.

Stan pustego katalogu ustawień należy do panelu kategorii, a nie do tego pasa, ponieważ mówi
o zawartości katalogu, a nie o przebiegu odczytu.

Zawężenie zależności pasa do fazy odczytu i powodu niepowodzenia czyni z niego jeden byt dla obu
odczytów okna konfiguracji: katalogu ustawień oraz konfiguracji obowiązującej. Stan konfiguracji
spełnia ten kształt bez żadnej zmiany.

## budowa/klient-poprzedni/src/asystent-plywajacy/indeks.ts

Warstwa wystawia na zewnątrz jedno wywołanie `zaczepAsystenta`, przyjmujące
kanał oraz opcjonalne posunięcia i identyfikator klienta. Pozostałe części,
czyli favikon, dymek, stan rozmowy i dostępność głosu, są wnętrzem warstwy
i nie mają powodu być widoczne z aplikacji.

Pierwszym wyjątkiem są `czyGlosDziala` oraz `brakujaceOgniwa`. Stan kanału
głosowego mówi o produkcie, a nie o widoku, więc jest potrzebny każdej warstwie
budującej obsługę głosu.

Drugim wyjątkiem jest `sprawca-zdarzenia.ts`. Odczyt pól `actor`
i `actorClientId` z koperty przydaje się poza tą warstwą, ponieważ
`aplikacja/zrodlo-posuniec.ts` rozpoznaje sprawcę odciskiem okna, mając prawdę
w kopercie. Wystawienie tego odczytu w punkcie zbiorczym zapobiega powstaniu
drugiej, odmiennej realizacji tej samej reguły.

## budowa/klient-poprzedni/src/mission-control/sekcja-aktywnosci.ts

Nazwy miar w kaflach są dosłowne, ponieważ mają nie dopuścić do pomylenia wysycenia kanałów modelu z obciążeniem maszyny. Miara bez źródła w kontrakcie pokazuje kreskę i etykietę braku źródła zamiast wartości zastępczej, aby pulpit nie sugerował danych, których rdzeń nie dostarcza.

Pas decyzji nie należy do tej sekcji. Sekcja udostępnia jedynie miejsce montażu pod rzędem kafli, a sam pas dokłada moduł mission-control.ts, który zna kolejność elementów pulpitu.

## budowa/klient-poprzedni/src/moduly/apps/wykaz-plikow-warsztatu.ts

Wykaz pokazuje wszystkie pliki warsztatu okna, dzięki czemu widać, że wpisana ścieżka należy już do istniejącego pliku, zanim zapis go nadpisze. Wiersz jest przyciskiem, ponieważ wiersz nieklikalny kazałby przepisywać ścieżkę ręcznie, a przepisana ścieżka bywa w rzeczywistości innym plikiem: klucz warsztatu tworzą okno, warstwa i ścieżka rozumiane co do znaku.

Wykaz niczego nie streszcza. Rozmiar i wersja idą tak, jak oddał je rdzeń, a brak pola zostaje nazwany brakiem, zamiast zostać zastąpiony wartością zerową, która wyglądałaby jak wynik pomiaru.

## budowa/klient-poprzedni/src/komponenty/dymek.ts

Znak zapytania jest przyciskiem, nie ozdobą: naciśnięcie prowadzi do niego
ognisko, a ognisko pokazuje objaśnienie. Treść czytają technologie wspomagające
z atrybutu `aria-label` znaku, więc dymek nie potrzebuje identyfikatora i nie
zderza się z dymkami innych okien.

Klasy wywołującego godzą jedną fabrykę z wyglądem właściwym dla miejsca
wywołania. Klasa `.dn-btn-ikona` to kwadrat 32 px, na dotyku 40 px, o narożniku
`--dn-r-sm`, z przezroczystym obrysem i wskaźnikiem `pointer`. Znaki modułów to
obwódki 14–16 px o narożniku `--dn-r-pill`, z widocznym obrysem i wskaźnikiem
`help`. Powłoka niesie z kolei zachowanie w rzędzie — `flex: none`, odstęp
i `vertical-align` — którego `.dn-tooltip` nie zna.

## budowa/klient-poprzedni/src/moduly/library/filtr-wersji.ts

Filtr działa w całości po stronie okna, ponieważ komenda wykazu wersji
biblioteki nie przyjmuje pola zawężającego, a jej odpowiedź niesie komplet pól,
po których dokumentacja każe zawężać, czyli sprawcę zmiany i etykietę wersji.
Wysłanie w tym celu drugiego odczytu byłoby pytaniem o to, co okno już ma.

Sprawca zmiany jest napisem, a nie wyliczeniem kontraktu: rdzeń zapisuje w polu
sprawcy to, co poda wołający. Pozycje zawężające do zmian Operatora oraz do
zmian modelu dopasowują więc po zawartości napisu, a wersja o sprawcy
nienazwanym wchodzi wyłącznie do pozycji obejmującej wszystkie wersje.
Zgadywanie po stronie okna kazałoby takiej wersji trafić do jednej z dwóch grup
bez podstawy.

## budowa/klient-poprzedni/src/moduly/library/sterowanie-widoku.ts

Stopień obrotu i krotność powiększenia trafiają do zbioru danych elementu, a cały
wygląd wynikający z tych wartości bierze się z arkusza stylów modułu. Plik widoku
nie zna zatem żadnej wartości wizualnej: nie podaje wymiarów, odstępów ani barw,
a zmiana wyglądu obrotu i powiększenia odbywa się wyłącznie w arkuszu.

## budowa/klient-poprzedni/src/modele/warstwy-tozsamosci.ts

Warstwa mówi, jak krytyczna jest treść nakładki: konstytucja stoi wyżej niż
profil roli, a profil wyżej niż ekspertyza. Ta kolejność rządzi zarówno układem
wykazu kategorii, jak i porządkiem złożonego promptu. Wartości pochodzą wyłącznie
ze stałych kontraktu.

Tryb rozstrzyga los ustawień fabrycznych: `ZASTAP` podmienia prompt fabryczny
w całości, `DOLACZ` dokłada treść do niego. Wybór decyduje o tym, którym
argumentem uruchomienia pojedzie nakładka, więc każdy przełącznik trybu w tej
sekcji podaje skutek wprost.

## budowa/klient-poprzedni/src/komponenty/odmowa.ts

Kod błędu zostaje przy wiadomości z zamysłem. Kody `validation_failed`
i `not_found` znaczą dla Operatora co innego, a treść wiadomości bywa dla obu
tym samym zdaniem, więc zdanie bez kodu odbierałoby Operatorowi rozróżnienie.
Brak wiadomości nie zostawia pustki: idzie wtedy sam kod, a gdy nie ma i kodu,
zdanie nazywa milczenie rdzenia wprost, żeby Operator nie patrzył na pusty
prostokąt.

Kształt powodu odmowy jest luźniejszy od kształtu `ErrorInfo` z kontraktu,
ponieważ moduły dostają powód z pola `blad` wyniku komendy, które bywa puste,
a pojedyncze okna trzymają własne, częściowe zapisy odmowy.

Nazwa nieznanej komendy obsługuje przypadek osobny, otwarty dla każdego modułu:
rdzeń odpowiedział zdarzeniem mówiącym, że komendy nie zna. Nie jest to awaria
wykonania, a Operator musi widzieć różnicę między usterką a czynnością, której
rdzeń jeszcze nie umie, dlatego nazwa typu idzie na początek zdania.

Zdanie składane z powodu wziętego w całości istnieje po to, żeby wywołanie
po odmowie komendy nie powtarzało w każdym miejscu pary pól `code` i `message`.
Para rozjeżdża się przy przepisywaniu, ponieważ łatwo podać kod z jednego
wyniku, a wiadomość z drugiego.

## budowa/klient-poprzedni/src/moduly/browser/przyciski-browser.ts

Panele pomocnicze i wiersze wykazu modułu Browser potrzebują tego samego:
zbudować przycisk i podpiąć nasłuch naciśnięcia. Czynność ta stoi w jednym
pliku, żeby cztery pliki modułu nie rozjechały się przy pierwszej poprawce.
Nazwa funkcji budującej zawiera wyraz „przycisk”, ponieważ po nim rozpoznają
kontrolkę narzędzia zestawiające etykiety czynności okien; opakowanie nazwane
inaczej chowa etykietę przed takim zestawieniem. Nazwy odmian zastępują napis
klasy powtarzany przy każdym wywołaniu: `glowny` stoi tam, gdzie panel ma jedną
czynność wiodącą, na przykład zapis albo dodanie, `zarys` przy pozostałych
czynnościach panelu, `duch` przy pozycjach wykazu, gdzie przycisków w jednym
wierszu jest wiele.

## budowa/klient-poprzedni/src/modele/panel-kont.ts

Panel pamięta wyłącznie, które konto jest już wypełnione w formularzu. Bez tej pamięci każde zdarzenie zmiany konta kasowałoby treść właśnie wpisywaną przez Operatora, dlatego wypełnienie formularza następuje przy zmianie konta czynnego, a nie przy każdym przeliczeniu stanu.

Ograniczenie rodzaju kont jest częścią żądania wykazu kont kierowanego do rdzenia, a nie filtrem zakładanym w kliencie. Rejestr kont bywa długi, a kolejność rotacji zna rdzeń, więc wybór rodzaju musi trafić do żądania, żeby wykaz pozostał zgodny z rejestrem.

## budowa/klient-poprzedni/src/moduly/diagnostics/okno-logs-viewer.ts

Port Diagnostyka jest w rdzeniu odbiorcą odmów wykonania komend. Każda odmowa
dostaje wpis dziennika, w którym pole źródła niesie nazwę odrzuconej komendy —
także wtedy, gdy zapis wiersza błędu się nie powiedzie i wpis dziennika jest
jedynym śladem tej odmowy. Dlatego filtr źródła i wzorzec wyszukiwania stoją
na pierwszym miejscu paska narzędzi okna.

Przy włączonym przełączniku wyrażenia regularnego okno próbuje zbudować wzorzec
z wpisanego pola, zanim żądanie pójdzie do rdzenia: błąd składni zatrzymuje
żądanie i mówi to wprost, zamiast wysyłać wzorzec, który rdzeń odrzuciłby po
cichszej stronie swojej.

Zakres czasu jest wspólny całemu modułowi. Okno nie prowadzi własnych pól
zakresu: czyta zakres ze stanu modułu i nasłuchuje jego zmiany, żeby
przestawienie zakresu w Diagnostics Center przeliczyło wykaz bez osobnej
czynności okna.

Pola informujące o obcięciu wyniku i o łącznej liczbie spełniających warunki
trafiają do zdania potwierdzenia, ponieważ wykaz ucięty granicą bufora bez
wiedzy czytelnika jest kłamstwem tej samej rodziny co pusty wykaz pokazany
przy odmowie.

Odświeżanie na żądanie znaczy tu odczyt wywołany czynnością, nie strumień na
żywo. Kontrakt nie niesie zdarzenia nowego wpisu dziennika, więc okno nie
odpytuje rdzenia w pętli — odczyt biegnie na żądanie Operatora i na zmianę
zakresu wspólnego modułu.

Rozwijana lista poziomu wpisu korzysta z biblioteki komponentu drzewa menu,
osadzonego przez moduł wyboru z menu w warstwie okien aplikacji, zamiast
z natywnego elementu wyboru.

## budowa/klient-poprzedni/src/moduly/browser/wiersz-notatki.ts

Powiązanie ze źródłem jest widoczne, a nie domyślne. Notatka bez identyfikatora źródła mówi o tym
wprost, ponieważ powiązanie zakłada się z wykazu okna, a jego brak jest informacją, a nie pustką do
przemilczenia.

Klasyfikacja stoi przy wierszu jako lista wyboru, a nie jako ikona. Rodzaj notatki ma być czytelny
bez najeżdżania na znak i bez rozróżniania barw, a zmiana klasyfikacji ma być jednym gestem.

## budowa/klient-poprzedni/src/aplikacja/gospodarz-dokumentu.ts

Gospodarz ma jedną odpowiedzialność: przygotowuje dokument oraz jeden element,
do którego router wstawia widoki tras. Nie zna rdzenia ani żadnego widoku i nie
wie nawet, ile tras ma aplikacja.

Brak kontenera przygotowanego przez `index.html` nie zatrzymuje uruchomienia.
Gospodarz trafia wtedy do ciała dokumentu, a aplikacja rusza normalnie.

Tytuł dokumentu składa nazwę projektu wziętą z opisu okna z nazwą produktu,
ponieważ nazwa projektu pochodzi z konfiguracji budowania i różni się między
osadzeniami, natomiast nazwa produktu jest stała.

## budowa/klient-poprzedni/src/moduly/agents/sklad-zespolu.ts

Kolejność zaznaczania jest treścią, nie porządkiem wyświetlania. Pole `agentIds`
niesie ekspertów w kolejności nadanej przy wyborze; składanie tablicy przejściem
po bibliotece brałoby kolejność z wykazu ekspertów, więc zespół wracałby po zapisie
w innym porządku niż złożony. Wybór trzyma zatem własną listę i dopisuje na jej koniec.

Zespół wczytany może wskazywać eksperta spoza bieżącego wykazu: zarchiwizowanego albo
odciętego frazą zawężającą. Taki wpis zostaje w składzie i pokazuje się wierszem
nazywającym brak, żeby ponowny zapis nie okroił zespołu po cichu.

## budowa/klient-poprzedni/src/moduly/assistant/zrodlo-zaplecza.ts

Odczyt `window.list` podaje okna komunikacji sesji. Komenda `assistant.voice.command`
wymaga pola `windowId`, a moduł okien sam nie zakłada: bierze to, które rdzeń
przypisał modułowi.

Odczyt `action.list` podaje katalog akcji zasięgu modułu. Siatka szybkich akcji Voice
Console jest widokiem tego katalogu, a nie listą zaszytą w kliencie: nowa akcja jest
nowym wierszem rdzenia, nie zmianą kodu.

Odczyt `component.list` podaje komponenty własne rodzaju `assistant`, czyli profile
asystenta ze strefy drugiej strony głównej. Jest to jedyny odczyt kontraktu, którym da
się wypełnić selektor profilu Voice Console. Komponenty wyłączone wchodzą do wykazu:
profil wyłączony w strefie drugiej wciąż istnieje w rdzeniu i wskazanie go należy do
Operatora, a nie jest pomyłką okna.

Odczyt `speech.availability.get` mówi, czy silnik mowy stoi na tej maszynie. Brak
silnika jest odpowiedzią, a nie awarią, więc okno mówi o nim wprost.

Wywołanie `speech.transcribe` rozpoznaje nagranie wskazane ścieżką na maszynie silnika.
Pusta transkrypcja jest wynikiem prawidłowym, gdy pole `processed` niesie prawdę, ponieważ
nagranie mogło zawierać ciszę albo szum; sprawdzian kształtu pyta zatem o obecność pól,
a nie o niepustą treść.

## budowa/klient-poprzedni/src/moduly/design/panel-warstw.ts

Podział ról między kanwą a panelem jest następujący: kanwa pokazuje położenie
warstwy, panel pokazuje porządek, blokadę i adnotację. Adnotacja jest polem
notatki warstwy w kontrakcie, więc zostaje w kompozycji na trwałe.

Blokada dotyczy wyłącznie położenia warstwy na kanwie. Wiersz warstwy
zablokowanej pozostaje klikalny i edytowalny, a stan blokady niesie atrybut
danych wiersza.

## budowa/klient-poprzedni/src/moduly/apps/wykaz-komponentow-warstwy.ts

Okno Frontend Workspace i okno Backend Workspace korzystają z tego samego zbioru komponentów pochodzącego z Architecture Designera i różnią się wyłącznie tym, które rodzaje komponentów do nich należą. Dlatego oba okna składa jeden moduł wykazu, sparametryzowany warstwą warsztatu.

Przypisanie rodzaju komponentu do warstwy jest wyborem klienta, a nie kontraktu. Ani typ AppWorkspaceLayer, ani typ AppComponentKind nie łączą tych dwóch zbiorów, więc podział został wypisany wprost w stałej RODZAJE_WARSTWY, aby stał w jednym miejscu i nie rozchodził się po oknach.

## budowa/klient-poprzedni/src/moduly/browser/warstwy-widocznosci.ts

Ujawnianie nie jest bramą. Element zwinięty nie jest zablokowany — jest schowany, a każda
jego czynność zostaje osiągalna wyzwalaczem albo skrótem. Dlatego skrót warstwy czwartej
działa również wtedy, gdy tryb administracyjny jest wyłączony: tryb decyduje o obecności
wyzwalacza w pasku, nie o dostępie do funkcji.

Skrót podpina się do elementu modułu, a nie do dokumentu. W powłoce stoi obok siebie kilka
modułów, a nasłuch założony na dokumencie odpowiadałby także na klawisze wciśnięte w cudzym
oknie. Skrót działa więc wtedy, gdy fokus stoi wewnątrz modułu. Powłoka gospodarza może
przechwycić kombinację przed stroną — wtedy pozostaje wyzwalacz w pasku.

Licznik pozycji rozszerzenia jest polem pominiętym, gdy rozszerzenie niczego nie zlicza.
Wartość zero jest nieodróżnialna od braku licznika, więc pominięcie niesie własne znaczenie.

## budowa/klient-poprzedni/src/aplikacja/godlo-aplikacji.ts

Znak marki ma barwy własne, więc nie dziedziczy barwy tekstu i nie zmienia się
wraz z motywem; odmianę znaku dobiera podłoże, a pasek aplikacji jest atramentowy
w obu motywach, stąd podłoże ciemne. Droga przez `elementGodla` niesie ponadto
próg odmiany uproszczonej: przy rozmiarze szesnastu pikseli i mniejszym znak
przechodzi na jeden grot.

Nazwa produktu jest stała, a nazwa projektu przychodzi z opisu okna, czyli
z konfiguracji budowania utrwalonej w `okno-komunikacji/opis-okna.ts`. Blok
godła nie jest kontrolką i nią nie udaje: nie prowadzi nigdzie i nie zmienia
kursora na wskazujący.

## budowa/klient-poprzedni/src/ladowanie/indeks.ts

Poza punktem wejścia katalog nie wystawia niczego. Arkusz stylu sceny i postać
jej bryły pozostają sprawą wewnętrzną, więc zmiana układu ekranu wejścia nie
sięga warstw, które scenę stawiają.

## budowa/klient-poprzedni/src/modele/znak-wykazu.ts

Klasa wyglądu pochodzi z biblioteki — jest to `dn-plakietka` wraz z jej
odmianami — a klasa miejsca należy do wykazu, w którym plakietka stoi. Dlatego
obie stoją w osobnych parametrach: wygląd znaku i jego położenie w wierszu są
dwiema niezależnymi decyzjami wywołującego.

Plik nie wymienia żadnej barwy. Barwa plakietki wynika wyłącznie z odmiany klasy
bibliotecznej, więc zmiana palety nie sięga tego pliku.

## budowa/klient-poprzedni/src/modele/pola-konta.ts

Zestaw pól jest oddzielony od formularza, ponieważ formularz odpowiada za rozgałęzienie
zapisu na dwie komendy kontraktu i za to, co wolno nadpisać w polu już wypełnionym.
Trzymanie obu odpowiedzialności w jednym miejscu wiązałoby kształt danych z przebiegiem
zapisu.

Opis poświadczenia należy do pola, a nie do widoku: kontrakt nie zwraca poświadczenia
żadną komendą, więc pole konta, które poświadczenie ma, i tak wygląda na puste. Opis
jest jedynym miejscem, w którym ten stan zostaje nazwany.

## budowa/klient-poprzedni/src/moduly/browser/pozycje-pokrycia.ts

Powód nieczynnej pozycji rozstrzyga byt wspólny wszystkich modułów, stojący
w pliku `moduly/pokrycie-komend.ts`. Byt pyta rdzeń o wykaz jego komend
powitaniem `connection.hello` i rozróżnia trzy stany: kontrakt danej komendy
nie zna, kontrakt ją zna, lecz rdzeń nie ma uchwytu, albo rdzeń ma uchwyt,
a okno jeszcze go nie wywołuje. Zdanie wpisane w moduł na sztywno mówiłoby
nieprawdę od chwili, w której rdzeń dostanie obsługę, i nic nie wymuszałoby
jego zdjęcia, ponieważ nic nie łączyłoby go z rdzeniem.

Przycisk nieczynnej pozycji nie jest wygaszony i nie milczy: naciśnięcie daje
dymek z powodem oraz wiersz odpowiedzi okna, a żądanie do rdzenia nie wychodzi
i nic nie udaje wykonania. Znaczniki `data-brak-komendy` oraz `data-pokrycie`
niesie kontrolka biblioteczna i po nich sięgają sprawdziany produktu.

Powód jest czytany z atrybutu `title` dopiero przy naciśnięciu, tak samo jak
robi to dymek kontrolki. Byt przerysowuje ten atrybut po każdej zmianie wykazu,
więc domknięcie na treści z chwili budowy podawałoby stan odczytu w toku długo
po odpowiedzi rdzenia.

## budowa/klient-poprzedni/src/moduly/multitasking/obsada-rol.ts

Rola okna pochodzi z pola roli w kontrakcie, nie z domysłu. Koordynatorem jest
okno roli koordynatora, a wykonawcą okno roli wykonawcy, którego wskazanie okna
zarządzającego trafia w tego koordynatora — to jest cała definicja pętli.
Znakowanie ról po tytule albo po kolejności rozjechałoby rdzeń z widokiem przy
pierwszym przepięciu wykonawcy pod innego koordynatora.

Analityk stoi poza pętlą. Okno analizy wyników nie jest wykonawcą: nie kończy
tury, która miałaby wybudzić koordynatora, więc jego rolą kontraktową jest rola
samodzielna. Odróżnia go wskazanie zapisane na poziomie sesji, ponieważ kontrakt
nie ma czwartej wartości roli okna, a wymyślenie jej po stronie klienta
rozjechałoby wykaz ról z bazą.

Kody okien są bezmodułowe: rejestr okien operacyjnych niesie kod rozmowy
koordynatora, dwa kody rozmowy wykonawcy oraz kod analizy wyników, każdy bez
przedrostka modułu. Przynależność do modułu niesie osobna macierz, w której tych
czterech okien nie ma — liczą się tam jako okna pozamodułowe. Kod zapisany
z przedrostkiem nie trafiłby w żaden wiersz rejestru, a nawigacja po kodzie okna
nie miałaby dokąd skoczyć.

## budowa/klient-poprzedni/src/aplikacja/trasy.ts

Aplikacja ma trzy widoki najwyższego rzędu. Trasa `strona-glowna` prowadzi do
Centrum dowodzenia i stanowi wejście do produktu. Trasa `srodowisko` otwiera
powłokę środowiska wraz z kartami sesji i modułami. Trasa `pulpit` otwiera
Mission Control jako widok stojący obok strony głównej, a nie zagnieżdżony w niej.

Nazwa trasy pełni dwie role naraz: jest kluczem katalogu i wartością zapisywaną
w adresie dokumentu. Dzięki temu odświeżenie strony wraca do tego samego widoku,
a nie na początek przepływu.

## budowa/klient-poprzedni/src/moduly/apps/indeks.ts

Moduł nie osadza się sam w dokumencie: oddaje element, a warstwa składająca
rozstrzyga, gdzie go postawić. Dzięki temu te same okna wchodzą i w obszar
roboczy powłoki, i w stanowisko sprawdzianu.

Pole `zamknij` podpina `rozlacz()` złożenia modułu. Bez niego subskrypcje stanu
`apps.build.changed` oraz `progress.changed`, a wraz z nimi subskrypcje okien
pomocniczych, żyją dalej po zejściu modułu ze sceny. Pole `zamknij` jest
w typie `WidokModulu` opcjonalne, więc kompilator jego braku nie zgłosi.

## budowa/klient-poprzedni/src/moduly/multitasking/zrodlo-podagentow.ts

Źródło jest jednym z trzech źródeł modułu, obok `zrodlo-biegu.ts` niosącego bieg
pracy i `zrodlo-okien.ts` niosącego obsadę ról. Stoi osobno od biegu pracy,
ponieważ bieg porusza oknami ról — tura, przekazanie, kolejka etapu, telemetria
— a podagent jest bytem pod oknem wykonawcy, o którym kontrakt nie rozstrzyga,
czy jest oknem, procesem, czy pozycją kolejki.

Telemetria zostaje w biegu pracy: panel bierze `monitor.status` ze źródła biegu,
które tę komendę już niesie. Drugie wywołanie tej samej komendy w drugim źródle
byłoby kopią bez powodu.

Powołanie oddaje podagentów tak, jak założył je rdzeń, a nie tak, jak prosił
formularz: liczba powołanych pochodzi z odpowiedzi, ponieważ rdzeń ma prawo
powołać ich mniej, a górna granica piętnastu należy do rdzenia, nie do widoku.
Zawężenie wykazu po stanie również wykonuje rdzeń — żądanie niesie pole
`status`, więc sito po stronie klienta byłoby drugim, rozjeżdżającym się sitem.

Zebranie wyników oddaje komplet podagentów wraz z polem `complete`, którego
panel nie zgaduje ze stanów pojedynczych podagentów, tylko bierze z odpowiedzi.
Zdarzenie zmiany stanu rdzeń rozgłasza przy powołaniu, wejściu w bieg,
zakończeniu i zatrzymaniu, więc wykaz nadąża bez ręcznego odpytywania.

Zatrzymanie idzie komendą `subagent.stop` wprost, a nie obejściem przez
`queue.action` na kolejce o nazwie podagenta: podagent powołany bez wpiętej
kolejki nie miałby wtedy czego zatrzymać. Odpowiedź rozróżnia `stopped` od
`notRunning`, ponieważ podagent już zakończony nie jest błędem, tylko innym
stanem, o którym panel mówi osobno.

## budowa/klient-poprzedni/src/aplikacja/ustawienia-okna-sledzone.ts

Migawka ustawień wraz z wykazem różnic mówi, co się zmieniło. Bez nich pas
meldunkowy podawałby samo zdanie o zmianie okna przy każdym dotknięciu, zamiast
nazwać dobrany model, rolę albo zasięg wykonania.

Odcisk okna mówi, czy zmiana jest własna: zdarzenie przynoszące dokładnie ten
stan okna, co świeża odpowiedź na własną komendę, jest skutkiem tej komendy.
Do odcisku nie wchodzi czas zmiany zapisany w polu `updatedAt`, ponieważ rdzeń
nadaje go osobno w zdarzeniu i osobno w odpowiedzi; doklejenie czasu sprawiłoby,
że własna czynność nigdy nie zrównałaby się sama ze sobą.

Plik nie zna ani widoku, ani kanału. Niesie sam przekład pól okna na wartości
porównywalne oraz na zdania przeznaczone dla czytającego.

## budowa/klient-poprzedni/src/moduly/diagnostics/niepowodzenie-odczytu.ts

Wynik oznaczony jako nieudany powstaje w warstwie protokołu na więcej niż jeden sposób, a tylko jeden z nich jest odmową rdzenia. Pierwsza droga to odmowa rdzenia, czyli koperta ze statusem błędu i opisem błędu z kontraktu. Druga droga to odpowiedź nieczytelna: rdzeń odpowiedział powodzeniem, lecz treść nie ma kształtu przewidzianego kontraktem, więc sprawdzanie kształtu odpowiedzi zamienia ją na niepowodzenie z kodem nieudanej walidacji.

Nazwanie odpowiedzi nieczytelnej odmową byłoby oskarżeniem rdzenia o czyn, którego nie popełnił, ponieważ rdzeń odpowiedział i odpowiedział powodzeniem, a nie zrozumiał go klient. Okno nie ma prawa twierdzić o braku, którego rdzeń nie orzekł, tak samo jak nie ma prawa zamienić odmowy w pustkę.

Stan odpowiedzi udanej bez treści pełni rolę straży, a nie drogi osiągalnej w obecnym układzie. Wynik cząstkowy oddaje niepowodzenie bez pola błędu, gdy odpowiedź jest udana, lecz pusta, jednak do przeniesienia taka odpowiedź nie dochodzi, bo sprawdzanie kształtu przechwytuje ją wcześniej i okno mówi wówczas zdanie o odpowiedzi nieczytelnej. Stan zostaje w słowniku na wypadek zmiany warstwy protokołu.

Zdania stoją w słowniku, a nie w oknach, z tego samego powodu, dla którego stan treści stoi w osobnym komponencie: cztery okna mówiące o tej samej ciszy czterema różnymi zdaniami rozjeżdżają się przy pierwszej poprawce.

## budowa/klient-poprzedni/src/moduly/apps/panel-akcji.ts

Panel akcji modułu powstaje wyłącznie z rejestru rdzenia. Nowa akcja modułu jest nowym wierszem rejestru, a nie zmianą kodu klienta, dlatego w tym pliku nie ma ani jednej pozycji katalogu zapisanej na stałe.

Wiersz rejestru niesie nazwę komendy, którą akcja wywołuje, jednak generycznej drogi wywołania window.action nie woła żaden widok klienta. Naciśnięcie przycisku pozycji nazywa więc komendę wiersza i stwierdza, że droga generyczna czeka na konsumenta. Takie rozwiązanie przyjęto zamiast milczącego przycisku oraz zamiast pozorowania wykonania akcji.

## budowa/klient-poprzedni/src/mission-control/indeks.ts

Katalog pozwala powłoce zamontować pulpit i wpiąć źródło danych rdzenia jednym
ruchem:

    import {
      pustyKomplet,
      utworzMissionControl,
      utworzZrodloPulpitu,
    } from './mission-control/indeks';

    const pulpit = utworzMissionControl(pustyKomplet());
    const zrodlo = utworzZrodloPulpitu(kanal, (dane) => pulpit.odswiez(dane));
    zrodlo.uruchom();

Komplet danych pulpitu pochodzi wyłącznie z odczytów `session.list`,
`window.list` i `channel.list` oraz ze zdarzeń `session.changed`,
`window.changed`, `queue.changed` i `progress.changed`. Do pierwszego odczytu
pulpit pokazuje stany puste zamiast wartości zastępczych. Sam katalog nie
zawiera logiki: jest wykazem tego, co warstwa wystawia na zewnątrz.

## budowa/klient-poprzedni/src/moduly/browser/wykazy-zebranego.ts

Zaciągnięcie składa się z dwóch odczytów rdzenia: `browser.source.list`
i `browser.note.list`. Wynik obu przekłada się na zbiór zebranego w sesji.

Odmowa „okna nie ma" nie jest błędem okna. Rdzeń odpowiada kodem `not_found`,
gdy moduł Browser nie widział jeszcze wskazanego okna, a to normalny stan okna
świeżo otwartego, w którym nikt niczego nie zebrał. Rozpoznanie kończy się wtedy
stanem `nietkniete`: wykaz zostaje pusty, powód jest podany, ale okno nie wchodzi
w stan błędu. Każda inna odmowa jest błędem i mówi swoim powodem. Okno nieznane
zgłasza się przy obu wykazach naraz, więc jest to jedno rozpoznanie, a nie dwa
osobne niepowodzenia.

Stanu `odczyt` nie wpisuje tu nikt: wynik oddawany jest dopiero po obu komendach,
więc stan oczekiwania ustawia wołający w `stan-przegladania.ts` przed czekaniem.
Wartość stoi mimo to w jednym typie z pozostałymi, ponieważ panele czytają jeden
stan wykazu, a nie dwa niezależne.

## budowa/klient-poprzedni/src/moduly/multitasking/stany-okna.ts

Stan błędu jest rozłączny ze stanem pustki. Okno roli, które po odmowie odczytu
pokazałoby wyzerowany licznik zamiast błędu, sugerowałoby zatrzymany bieg tam,
gdzie bieg trwa, a odmowa odczytu i brak pozycji w wykazie są dla Operatora
dwiema różnymi wiadomościami.

## budowa/klient-poprzedni/src/moduly/apps/okno-frontend-workspace.ts

Plik jest wiązaniem, nie drugim widokiem: formularz warsztatu stoi raz,
w `okno-warsztatu.ts`, a stąd dostaje wyłącznie stan produktu i opis warsztatu
frontendu.

Kod okna nie pada tu wprost napisem. Ramę modułu Apps woła wspólne okno
warsztatu, biorąc kod ze zmiennej opisu, dzięki czemu wykaz `KODY_OKIEN`
pozostaje jedynym miejscem, w którym kody okien tego modułu są wypisane.

## budowa/klient-poprzedni/src/ikony/zrodla/srodowiska.ts

Nazwy i kolejność pozycji tego pliku odpowiadają grupie środowisk w wykazie
`ikony/manifest.json`, dzięki czemu wykaz ikon i wykaz źródeł dają się porównać
pozycja po pozycji.

Emblematy środowisk są znakami własnymi budowy i nie pochodzą z zestawu Lucide,
więc nie da się ich podmienić aktualizacją biblioteki ikon.

## budowa/klient-poprzedni/src/moduly/developer/stany-okna.ts

Stan błędu jest rozłączny ze stanem pustki. Okno, które po odmowie odczytu
pokazuje puste drzewo, mówi Operatorowi, że katalog jest pusty, zamiast
powiedzieć, że drzewa nie udało się odczytać.

Pustka bywa w tym module stanem poprawnym: repozytorium bez zmian do
zatwierdzenia nie jest usterką, więc pusty wykaz nie może być pokazywany jako
niepowodzenie.

## budowa/klient-poprzedni/src/aod/sekcja-rozmowy.ts

Ognisko okien zna rdzeń i mogło się ono zmienić między odczytem stanu a wysyłką,
więc sekcja nie podstawia w pole `windowId` okna odczytanego wcześniej. Puste pole
oznacza wprost „okno ogniskowane w chwili przyjęcia wiadomości". Odpowiedź niesie
`messageId` oraz `windowId` okna, które wiadomość przyjęło, i meldunek pokazuje oba.

Przycisk wysyłki pozostaje czynny także przy pustym polu treści. Pustą treść ocenia
rdzeń i to on zwraca odmowę `validation_failed`; blokada po stronie nakładki
rozdzielałaby ocenę wiadomości na dwa miejsca.

## budowa/klient-poprzedni/src/dostepy/okno-dostepow.ts

Sekcja dostępów sama nie zakłada, gdzie zostanie osadzona: oddaje element. Rama
okna jest jednak potrzebna od razu, ponieważ pierwszym gospodarzem sekcji jest
listwa ustawień Centrum dowodzenia, a ta otwiera widoki jako okna modalne, tak
samo jak okno konfiguracji.

Okno stoi na natywnym elemencie okna dialogowego, więc warstwę tła, stos okien
i zamknięcie klawiszem ucieczki daje przeglądarka, a nie własna nakładka. Wygląd
bierze z biblioteki komponentów, z klasy okna modalnego, więc plik nie zna ani
jednej barwy.

Okno otwiera się natychmiast, przed odpowiedzią rdzenia, a wykazy dojeżdżają do
niego odpowiedzią.

## budowa/klient-poprzedni/src/mission-control/kolumna-procesow.ts

Wykaz procesów w tle buduje się wyłącznie ze zdarzeń telemetrii, ponieważ kontrakt nie ma komendy
odczytu procesów bieżących. Do nadejścia pierwszego zdarzenia kolumna pokazuje stan pusty, który
nazywa tę okoliczność wprost, zamiast udawać brak procesów.

Stan procesu jest wypisany słowem obok paska postępu. Barwa paska wspiera odczyt, ale go nie
zastępuje, dzięki czemu wykaz pozostaje czytelny bez rozróżniania barw.

## budowa/klient-poprzedni/src/moduly/multitasking/powierzchnia-sekcji.ts

Sekcje panelu — Zespoły, Kolejki, Orkiestracja, Harmonogram i Monitor — mają ten sam
szkielet. Kolejność i widoczność podsekcji prowadzi wzorzec powłoki
`powloka/sekcje-panelu.ts` przez komendy `panel.sections.*`; ten plik wzorzec woła
i sam niczego nie przestawia.

Komendy `panel.sections.*` adresują parę okna i panelu, a panel orkiestracji należy do
środowiska, nie do okna. Adresem zostaje więc okno, przy którym sekcja pracuje: okno
koordynatora, a gdy obsada go nie ma — pierwsze okno sesji; bez żadnego okna układ
zostaje miejscowy. Brak zapisanego układu znaczy układ domyślny, więc odmowa
`panel.sections.get` zostawia sekcję kompletną.

Identyfikatorem panelu jest klucz sekcji, ten sam, którym boczna nawigacja wskazuje
sekcję. Napis składany z tytułu rozjechałby się przy pierwszej zmianie tytułu, a układ
zapisany w rdzeniu przestałby mieć odpowiednik na ekranie.

## budowa/klient-poprzedni/src/moduly/design/indeks.ts

Moduł nie osadza się sam w dokumencie i nie zna powłoki: oddaje element,
a warstwa składająca rozstrzyga, gdzie go postawić.

Rozłączenie zostaje wewnątrz modułu, ponieważ typ `WidokModulu` powłoki nie ma
czynności odpięcia. Subskrypcje kanału odpina `rozlacz()` widoku modułu.

## budowa/klient-poprzedni/src/modele/indeks.ts

Okno modeli jest jedno na klienta i trwa między otwarciami z tego samego powodu,
dla którego router nie porzuca widoku opuszczonej trasy: powtórne otwarcie wraca
do zakładki i kategorii, na której poprzednie się zakończyło. Trwałość okna
utrzymuje subskrypcję zdarzeń `account.changed` oraz `identity.changed`, dzięki
czemu zmiana dokonana na innym urządzeniu konta dociera również wtedy, gdy okno
pozostaje zamknięte.

Wywołanie z innym kanałem, właściwe dla ponownego połączenia z rdzeniem, buduje
okno od nowa, aby komendy nie szły przez transport, którego już nie ma. Sekcję
można też osadzić bez ramy okna: `utworzSekcjeModeli` oddaje ten sam byt bez
elementu `dialog` wokół niego, czyli jedna implementacja w dwóch oprawach.

## budowa/klient-poprzedni/src/moduly/automations/stany-okna.ts

Nośnik stanu treści okien modułu Automations różni się od nośnika wspólnego
zachowaniem metod `blad()` i `pusto()`: obie zdejmują wcześniejsze potwierdzenie,
zanim postawią własny komunikat. Nośnik wspólny (`komponenty/stan-tresci.ts`)
przestawia wyłącznie pas komunikatu, więc bez tego zdjęcia odmowa rdzenia stoi
obok potwierdzenia czynności poprzedniej, a zielone zdanie nad odmową potwierdza
czynność, która się nie odbyła.

Różnica siedzi w module, a nie w nośniku wspólnym, ponieważ ten sam nośnik
obsługuje okna pozostałych modułów. Przedrostek klas `da` należy do arkusza
modułu, dlatego wiązanie nośnika z przedrostkiem stoi po stronie modułu.

## budowa/klient-poprzedni/src/moduly/assistant/zrodlo-pamieci.ts

Wzorem jest `zrodlo-assistant.ts`: żadne wywołanie nie rzuca wyjątkiem, odmowa
wraca polem `blad` wyniku, a okno pokazuje ją w swoim stanie błędu.

Komenda `memory.detach` nie jest tu wystawiona. Kontrakt oznacza znaczenie
odpięcia jako nierozstrzygnięte, więc okno nie nadaje mu własnego sensu — tak
samo postępuje moduł Workspace w pliku `moduly/workspace/pamiec-pozycja.ts`.

Rodziny są dwie, ponieważ mówią o dwóch różnych bytach. Rodzina `memory.*`
prowadzi ustalenia, czyli zdania, które Operator kazał zapamiętać. Rodzina
`knowledge.*` prowadzi wskaźnik znaczenia zbudowany z treści już istniejących:
biblioteki, historii rozmów i plików przestrzeni roboczej. Zlanie ich w jedno
źródło zatarłoby, co jest ustaleniem, a co odnalezionym fragmentem.

## budowa/klient-poprzedni/src/asystent-plywajacy/sprawca-zdarzenia.ts

Napis nie bywa pewniejszy niż dowód, na którym stoi, dlatego brak pola `actor`
w kopercie zdarzenia daje rozstrzygnięcie „nieznany", a nie domysł.

Plik jest osobny od `aplikacja/rozstrzyganie-sprawcy.ts`, który rozwiązuje inne
zadanie: odpowiada na pytanie, czy czynność wykonało bieżące połączenie, i robi
to dla zdarzeń, które pola `actor` w ogóle nie niosą.

## budowa/klient-poprzedni/src/mission-control/pas-decyzji.ts

Pas jest jedynym miejscem pulpitu, w którym praca stoi, dlatego jako jedyny nosi
akcent: wstęgę trzech pikseli u góry, tło żetonu `--dn-akcent-tlo` oraz wezwanie
z cieniem akcentu. Akcent obejmuje wstęgę i przycisk, nie całą powierzchnię.

Kontrakt nie niesie odczytu przepływów wstrzymanych do decyzji: komplet danych
podaje wtedy wartość pustą, pas mówi o braku źródła i chowa wezwanie, ponieważ
przycisk wzywający do rozstrzygnięcia nieodczytanego wykazu byłby atrapą. Gdy
źródło jest, a nic nie czeka, pas zmienia treść na zdanie o braku oczekujących
przepływów, a wezwanie pozostaje czynne.

## budowa/klient-poprzedni/src/aod/sterowanie-obecnoscia.ts

Menu kebab nie jest tu budowane po swojemu: to ten sam komponent
`wyciszenie-menu.ts`, który stoi przy awatarze. Dwie kopie jednego menu czytają
jeden stan, więc wyciszenie założone przy awatarze widać w nagłówku bez żadnego
zszywania.

Stan nie jest tu przechowywany: obie kontrolki czytają i przestawiają
`StanObecnosci` z pliku `tryb-obecnosci.ts`, a rysują się z powrotem jego
powiadomieniem. Dzięki temu skrót klawiszowy i menu nigdy się nie rozjeżdżają.

## budowa/klient-poprzedni/src/moduly/apps/filtr-katalogu.ts

Filtrowanie idzie po stronie klienta, a nie w żądaniu, ponieważ komenda `extension.list`
zawęża wyłącznie rodzajem i stanem zainstalowania. Pozostałych faset — źródła, stanu
włączenia i tekstu — kontrakt w żądaniu nie ma, a cztery odczyty po jednym na fasetę dałyby
cztery migawki z czterech różnych chwil.

Zakres szukania jest ograniczony i okno musi to powiedzieć: przeszukiwane są nazwa, kod
i opis, bo tyle niesie pozycja katalogu. Szukanie po udostępnianych narzędziach
i znacznikach stoi w wykazie braków kontraktu, zamiast udawać, że pusty wynik znaczy brak
takiego rozszerzenia.

Porównanie tekstu idzie bez rozróżnienia wielkości liter, ale bez normalizacji znaków
diakrytycznych. Napisy różniące się wyłącznie znakami diakrytycznymi zostają dwoma różnymi
napisami, bo zrównanie ich w oknie kazałoby mu twierdzić coś o dopasowaniu, czego rdzeń
przy własnym szukaniu nie potwierdzi.

Zdanie podsumowujące wykaz jest konieczne, ponieważ pusty wykaz przy czynnym filtrze czyta
się jak pusty rejestr. Zdanie wymienia zakres szukania, bo tekst nieznaleziony w opisie bywa
nazwą narzędzia, której pozycja katalogu nie niesie.

## budowa/klient-poprzedni/src/mission-control/nadajnik.ts

Pulpit operacyjny nie sięga po magistralę połączenia, ponieważ jego zdarzenia
są zdarzeniami widoku, a nie kontraktu; to powłoka rozstrzyga, którą komendę
kontraktu z takiego zdarzenia zbuduje. Rozesłanie ładunku biegnie po kopii
zbioru odbiorców i w bloku ochronnym, więc odbiorca zgłaszający usterkę nie
odbiera zdarzenia pozostałym, a odpięcie w trakcie rozgłoszenia nie narusza
przebiegu pętli.
