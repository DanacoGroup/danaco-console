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
zdejmuje pozycję z wykazu samo. Nazwy stoją w jednym miejscu, żeby definicja wpisana
do kontraktu i zdanie widoczne Operatorowi mówiły o tej samej komendzie. Obszar
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

## budowa/klient-poprzedni/src/komponenty/okno-aplikacji.ts

Znakowanie służy oknom, które mają wiersz w katalogu okien operacyjnych, lecz nie
noszą ramy operacyjnej. Okno Konfiguracji i Okno Ustawień są modalami opartymi na
elemencie `dialog`, a ich wygląd niesie klasa `dn-modal`. Strona główna jest widokiem
pełnoekranowym. Okno Rozmowy dostaje ramę od gniazda układu równoległego
`okna-rownolegle/gniazdo-okna.ts`. W każdym z tych przypadków nagłówek, plakietka
roli i pas akcji z ramy operacyjnej byłyby elementem zbędnym albo powtórzonym.

Przypisanie kodu wprost do zbioru danych elementu, rozsiane po tych plikach,
dawałoby kilka miejsc, które mogą się rozejść. Funkcja nadaje przy tej samej
okazji etykietę dostępności z tej samej nazwy, więc kod okna i treść czytana przez
czytnik ekranu pozostają zgodne.

Funkcja zwraca ten sam element i ten sam typ, ponieważ modal potrzebuje typu
`HTMLDialogElement` do wywołania metod `showModal` oraz `close`. Zwracanie typu
`HTMLElement` odbierałoby te metody i zmuszało wywołujących do rzutowania.

## budowa/klient-poprzedni/src/aktualizacja/indeks.ts

Adres kanału wydań, porównanie wersji i most do powłoki są własnością tego
katalogu i na zewnątrz nie wychodzą; punkt wejścia aplikacji zna wyłącznie to,
co jest tu wymienione.

Widok wykazu wydań wychodzi stąd osobno od banera. Baner przywołuje go
przyciskiem, ale widok jest gotowy do osadzenia także w stałym miejscu, więc
jedno użycie nie przesądza o drugim.

## budowa/klient-poprzedni/src/moduly/assistant/indeks.ts

Kod modułu stoi wewnątrz opisu modułu, więc moduł sam mówi, którym jest modułem.
Rozjazd między nazwą zapisaną w `aplikacja/rejestr-modulow.ts` a rzeczywistością
staje się przez to niemożliwy.

Moduł nie zna powłoki: oddaje element, a warstwa składająca rozstrzyga, gdzie go
postawić i czy wpisać go do rejestru modułów.

## budowa/klient-poprzedni/src/aplikacja/arkusze-stylow.ts

Żetony motywu dają barwy, typografię, przestrzeń i ruch, a biblioteka `.dn-*`
stoi na nich. Arkusz `powloka.css` sprowadza zmienne widoku okna do żetonów
motywu, więc idzie po arkuszu tego widoku, a `aplikacja.css` osadza reguły
w obszarze roboczym powłoki i dlatego jest ostatni.

Arkusz kompletu sterowania `sterowanie/sterowanie.css` wciąga sam moduł panelu.
Trafia przez to do pakietu za tym wykazem, jako ostatni styl widoku, i w wykazie
nie figuruje.

## budowa/klient-poprzedni/src/moduly/diagnostics/stany-okna.ts

Rozłączność stanu błędu i stanu pustki waży w tym module szczególnie. Port
diagnostyki jest w rdzeniu odbiorcą odmów wykonania komend, więc okno błędów
pokazujące po nieudanym odczycie pusty wykaz ukrywałoby odmowy dwa razy: cudze
i swoją własną.

Pustka bywa tu zaś stanem poprawnym, ponieważ instalacja bez ani jednego błędu
w zadanym zakresie czasu nie jest usterką.

## budowa/klient-poprzedni/src/konfiguracja/dymek-objasnienia.ts

Sam dymek objaśnienia mieszka w bibliotece komponentów, a stąd jest wyłącznie
re-eksportowany, żeby okno konfiguracji miało jedno wejście.

Zdania objaśnienia są dopisywane, a nie zastępowane. Pozycja bez opisu nadal ma
czym objaśnić swój klucz, a pozycja z opisem zyskuje to, czego opis nie mówi:
jednostkę, poziom domyślny i wymóg ponownego uruchomienia.

## budowa/klient-poprzedni/src/dostepy/dialog-katalogu.ts

Most do powłoki natywnej mieszka w katalogu `powloka/`, ponieważ dotyczy powłoki
natywnej, a nie dostępów. Tutaj zostaje wyłącznie zapis wyniku przyjęty
w widokach dostępów: napis pusty znaczy, że katalogu nie wskazano.

## budowa/klient-poprzedni/src/dostepy/obszar-katalogu-roboczego.ts

Katalog roboczy modelu jest ustawieniem niezależnym od dostępów, ponieważ obie rzeczy odpowiadają na inne pytanie. Model może mieć wgląd w cudzy katalog i niczego w nim nie zapisywać, a własne katalogi sesyjne trzymać w miejscu instalacji aplikacji. Zlanie tych dwóch ustawień w jedno kazałoby otworzyć zapis wszędzie tam, gdzie model ma wyłącznie czytać.

Obszar nie jest kopią okna konfiguracji. Pokazuje dwa klucze, ponieważ bez nich sekcja dostępów byłaby niepełna: po nadaniu modelowi dostępu do katalogów następne pytanie dotyczy miejsca, w którym model będzie pisał. Pełny wybór poziomów zasięgu i osi zostaje w oknie konfiguracji, które prowadzi jedyny rachunek dziedziczenia wartości.

Wczytanie danych nie czeka na rdzeń i nie blokuje osadzenia obszaru. Pola stoją od razu z wartościami domyślnymi, a odpowiedź rdzenia nanosi na nie wartości obowiązujące.

## budowa/klient-poprzedni/src/moduly/design/tor-komendy.ts

Słownik zdań o torze komend stoi osobno od słownika etykiet modułu Design,
ponieważ tamten plik jest słownikiem nazw i braków, a ten niesie wiedzę o drodze
żądania: o kolejności kroków rdzenia, o tym, co zastępuje stan poprzedni,
i o tym, czego kontrakt nie przewiduje wcale.

Zdanie o torze nie orzeka o powodzie odmowy. Powód niesie opis odmowy wyjęty
z odpowiedzi rdzenia, a słownik dokłada wyłącznie to, czego rdzeń o sobie nie
mówi.

## budowa/klient-poprzedni/src/moduly/research/badanie-ksiazka-kodow.ts

Komenda zapisu książki kodów zapisuje ją w całości: pole kodów niesie stan po
zmianie, więc żądanie złożone z samych kodów wpisanych teraz wymazałoby wszystkie
dotychczasowe. Dlatego zapis idzie po odczycie i jest złączeniem, a nie
podmianą. Zdanie odpowiedzi mówi osobno, ile kodów zostało zachowanych i ile
dołożonych, ponieważ to jedyny sposób, żeby Operator rozpoznał wymazanie, gdyby
rdzeń zapisał co innego.

Kod istniejący rozpoznaje się po nazwie, bez względu na wielkość liter, ponieważ
nazwa jest tym, co Operator wpisuje; identyfikator kodu nadaje rdzeń i Operator
go nie zna. Kod dopisany jedzie z pustym identyfikatorem, tak samo jak nowe
pytanie badawcze w komendzie zapisu pytania warsztatu.

Definicja stoi w wierszu za znakiem rozdzielającym. Znak jest wyborem okna
i okno mówi o nim wprost przy chwycie pola. Nazwa kodu bez definicji jest
poprawna, ponieważ kontrakt ma pole opisu nieobowiązkowe.

Zdjęcie kodu z książki jest osobną czynnością i okno mówi Operatorowi, że dziś
jej nie ma. Kod zastany, którego Operator teraz nie wpisał, zostaje: wymazanie
reszty byłoby skutkiem, którego nikt nie zamówił.

## budowa/klient-poprzedni/src/mission-control/sekcja-utworz.ts

Sekcja tworzenia stoi na pulpicie, ponieważ bez niej jedyną drogą do pracy jest kliknięcie sesji już biegnącej w matrycy. Rząd kafli otwiera wejście tam, gdzie nic jeszcze nie biegnie, dlatego każdy kafel jest czynny od razu, żaden nie czeka na spełnienie warunku, a podpis kafla nazywa byt, który powstanie po naciśnięciu.

Etykiety kafli idą krojem bazowym półgrubym, a nie krojem szeryfowym, ponieważ różnica krojów niesie w tym miejscu znaczenie: krój szeryfowy oznacza wejście do środowiska, a krój bazowy oznacza zbudowanie komponentu.

## budowa/klient-poprzedni/src/moduly/apps/zrodlo-okna-modulu.ts

Każde żądanie obszaru niesie pole `windowId`, które kontrakt nazywa oknem
modułu Apps, więc identyfikator musi pochodzić z rdzenia, a nie z wyliczenia po
stronie widoku. Podaje go komenda `window.list` zawężona do sesji i do okien
modułu Apps, przy czym pole `Window.moduleId` niesie to, co rdzeń przypisał
oknu przy wejściu w obszar roboczy komendą `workspace.enter`.

Brak okna jest stanem poprawnym, ponieważ sesja bywa jeszcze nieznana w chwili
wejścia w moduł. Wykaz wraca wtedy pusty, a okna modułu pokazują stan pusty
nazywający brakujący warunek, zamiast wysyłać żądanie bez identyfikatora.

Wybór okna bierze pierwsze okno o module zgodnym z kodem modułu. Gdy rdzeń nie
przypisał żadnego, wynikiem jest pusty łańcuch, a nie okno przypadkowe.

## budowa/klient-poprzedni/src/konfiguracja/widocznosc-pol.ts

Pozycja katalogu może zależeć od innej: pole „host wykonania" ma sens dopiero
przy zasięgu zdalnym, a pole „ścieżka klucza" dopiero przy moście SSH. Warunek
żyje w katalogu, nie w kodzie interfejsu, więc zmiana zależności jest nowym
wierszem katalogu, a nie nową gałęzią w kliencie.

Warunek wskazujący klucz, którego katalog nie zna, nie chowa pola. Ukrycie pola
z powodu braku metadanej byłoby cichym odebraniem Operatorowi ustawienia,
dlatego brak klucza warunkującego zostaje odnotowany w dzienniku klienta, a pole
pozostaje widoczne.

## budowa/klient-poprzedni/src/mission-control/zdarzenia-pulpitu.ts

Pulpit nie wysyła komend samodzielnie: nadaje zamiar, a kopertę składa z niego
powłoka. Dzięki temu nazwa komendy pojawia się w jednym miejscu i pochodzi
z importu ze wspólnego kontraktu, zamiast być przepisana z pamięci do każdego
wywołania.

Pole roli niesie nazwę kolejki widoczną dla Operatora. Kontrakt roli kolejki
nie niesie, więc zamiar jej nie podstawia i nie udaje, że pochodzi z rdzenia.

Podniesienie priorytetu nie ma odpowiednika ani w wyliczeniu `QueueAction`, ani
wśród komend kontraktu. Zamiar niesie wtedy rodzaj priorytetowy z komendą pustą,
zamiast napisu wymyślonego po stronie klienta, którego rdzeń nie rozpozna.

## budowa/klient-poprzedni/src/mission-control/kolumna-zespolu.ts

Liczba podagentów nie ma źródła w kontrakcie. Wiersz mówi to wprost, zamiast pokazać
wartość zastępczą, ponieważ liczba zmyślona wygląda tak samo jak odczytana i wprowadza
w błąd przy ocenie obciążenia zespołu.

## budowa/klient-poprzedni/src/aplikacja/pasek-aplikacji.ts

Pasek ma jedną odpowiedzialność: złożenie elementów paska. Stoi nad widokami,
które własnego paska nie mają — Centrum dowodzenia oraz Mission Control.
Powłoka środowiska ma pasek własny w `powloka/pasek-gorny.ts` i drugiego nie
dostaje; przełącznik tras oraz przełącznik motywu wchodzą wtedy w jej grupę
akcji.

Motyw jest przełączalny z każdego widoku, ponieważ pasek niesie ten sam
przełącznik, którego używa powłoka. Wartości motywu ani żadnej barwy pasek nie
zna — całą pracę wykonuje warstwa `motyw/`.

Poza godłem, przełącznikiem tras i motywem pasek niesie menu aplikacji przy
godle, ikonę ustawień oraz menu Operatora w grupie akcji. Są to te same wejścia,
które ma powłoka środowiska w `powloka/akcje-paska.ts`, żeby produkt zachowywał
się jednakowo na każdej trasie. Żadne z nich nie zakłada nowej czynności ani nie
woła komendy spoza wykazu: pozycje są te same, co w listwie strony głównej,
a skutek jeden i wspólny — wykaz skutków stoi w `akcje-ustawien.ts`.

Trasy stoją w grupie akcji, a nie osobno pośrodku paska: znaki tras i kontrolki
platformy tworzą jeden rząd ikon przy prawej krawędzi, rozdzielony kreską, więc
pasek ma dwie grupy zamiast trzech.

Ikona ustawień jest skrótem, a nie drugą drogą: otwiera dokładnie tę pozycję,
którą listwa strony głównej niesie pod kodem `konfiguracja`, przez ten sam wykaz
skutków. Stoi w pasku, ponieważ ustawienia platformy wykonuje się z każdego
miejsca, a listwa widoczna jest wyłącznie na stronie głównej.

## budowa/klient-poprzedni/src/moduly/browser/zrodlo-izolacji.ts

Panel modułu jest wyłącznie odczytem, dlatego źródło wystawia same komendy
odczytu. Zapis przełączników oraz przypisanie profilu należą do okna
konfiguracji punktów izolacji, ponieważ dwa miejsca zapisujące tę samą politykę
dawałyby dwa różne zdania o tym, co obowiązuje.

Rozdzielenie odpowiada podziałowi punktów izolacji: moduł przeglądania wskazuje
punkty właściwe przeglądaniu, czyli dostęp sieciowy procesu sesji, kontenery
tożsamości oraz zakres pętli, a personalizacja tych punktów odbywa się w oknie
konfiguracji.

## budowa/klient-poprzedni/src/moduly/library/zrodlo-znaczenia.ts

Rodzina komend `knowledge.*` ma drugiego klienta w module Wiedza, z którego moduł Library świadomie
nie korzysta. Powód pierwszy jest własnościowy: moduł nie sięga do plików innego modułu, ponieważ
wiązałby swoje okna z cudzym cyklem życia. Powód drugi jest techniczny i ważniejszy: tamto źródło
woła kanał wprost, a rodzina `knowledge.*` nie ma własnego zdarzenia odmowy, więc rdzeń bez wpiętego
portu Wiedzy odpowiada kopertą `connection.unknown`, pozbawioną pola `status`. Library idzie zatem
przez straż odmów, która wiąże taką odpowiedź z żądaniem po identyfikatorze i zamienia ją w zwykły
wynik z błędem.

Rodzina leży poza obszarem `library`, ale czyta dokładnie ten zbiór. Czytelnikiem zakresu `library`
jest repozytorium modułu Library, a pole `sourceId` trafienia niesie identyfikator pliku biblioteki,
co ustala przejściówka rdzenia `adapter_modul_wiedza_zrodla.go`. Dzięki temu trafienie daje się
odwzorować na wiersz wykazu bez drugiego odczytu.

Rozdział zdolności między rodzinami jest wiążący dla okna. Komenda `library.file.search` dopasowuje
słowa poprzez indeks pełnotekstowy repozytorium, a `knowledge.search` dopasowuje znaczenie poprzez
wskaźnik osadzeń. Tryb hybrydowy okna nie jest komendą kontraktu, lecz złożeniem obu odpowiedzi po
stronie klienta.

Wskaźnik osadzeń nie odświeża się przy wgraniu pliku. Rdzeń buduje go wyłącznie na żądanie
`knowledge.index`, ponieważ osadzanie potrafi trwać minutę i sięga po wagi modelu. Przeliczenie
wskaźnika jest więc osobną czynnością w zakładce higieny, a nie skutkiem ubocznym wgrania pliku.

## budowa/klient-poprzedni/src/aplikacja/przelacznik-motywu.ts

Oba motywy są równoprawne i kontrolka nie wskazuje żadnego z nich jako
domyślnego. Dopóki nie zapadł wybór własny, obowiązuje preferencja systemu,
a kontrolka powrotu stoi w stanie wciśniętym: pokazuje, że rozstrzygnięcie należy
do systemu.

Wartości motywu ani żadnej barwy ten plik nie zna — całą pracę wykonuje warstwa
`motyw/`, a tutaj zostaje wyłącznie obsługa kontrolek. Zmiana preferencji systemu
przy braku wyboru własnego dochodzi tą samą drogą zdarzenia co przełączenie
ręczne, więc kontrolka nadąża bez przeładowania okna.

## budowa/klient-poprzedni/src/konfiguracja/kontrolki-wyboru.ts

Definicja nastawy bez podanych opcji nie daje pustego ekranu bez wyjaśnienia: obie listy
niosą wtedy ostrzeżenie przy polu, a lista jednokrotna zostaje przy samej pozycji pustej.
Operator widzi wtedy, że wartości nie podał katalog nastaw, a nie że kontrolka się nie
zbudowała.

## budowa/klient-poprzedni/src/moduly/library/pasek-kontekstu.ts

Okno komunikacji wybierane jest po oznaczeniu modułu, a nie z brzegu wykazu, ponieważ komenda wykazu okien oddaje okna wszystkich modułów sesji. Sięgnięcie po pierwsze okno z wykazu wysyłałoby kontekst cudzego okna, dlatego pasek szuka okna należącego do modułu Library.

Okno komunikacji jest zarazem oknem źródłowym przeniesienia kontekstu, więc jego brak pasek opisuje wprost, zamiast podstawiać okno innego modułu. Tytuł okna pasek pokazuje razem z liczbą wiadomości, natomiast moduł okna i tryb uprawnień zostają w panelu sterowania okna, a tytuł powtórzony jest już w karcie sesji.

## budowa/klient-poprzedni/src/moduly/design/brak-drogi.ts

Przycisk czynności bez komendy w kontrakcie zachowuje się według wzoru
z `aplikacja/zamiary-pulpitu.ts`: pozostaje aktywny i podpisany, a powód braku
podaje dopiero po naciśnięciu. Kontrolka wyłączona albo pusta nie niesie tej
informacji, więc Operator nie odróżniłby czynności niedostępnej chwilowo od
czynności, której kontrakt w ogóle nie zna.

Znak `aria-description` niesie tę samą wiedzę technologiom wspomagającym, zanim
przycisk zostanie naciśnięty: zapowiedź płynie z etykiety, a powód z dymka
otwieranego naciśnięciem.

## budowa/klient-poprzedni/src/moduly/automations/kody-okien.ts

Rejestr okien operacyjnych rdzenia zna kody `workflow-builder`, `scheduler`,
`queue-manager`, `orchestrator` i `execution-monitor`, zakładane migracją
`migracja_030_rejestr_okien_operacyjnych.sql`; powiązanie tych kodów z modułem
`automations` niesie `migracja_031_okna_modulow.sql`. Kody stoją w jednym
miejscu, ponieważ rdzeń nie sprawdza identyfikatora okna podawanego
w `automation.execution.subscribe`, tylko zapamiętuje go jako klucz obserwacji.
Kod z przedrostkiem modułu nie wywołałby więc odmowy: trafiłby do
`Queue.windowIds` i do telemetrii postępu jako okno nieistniejące.

Znacznik wykazu gotowych pętli stoi poza wykazem kodów okien, bo służy wyłącznie
znacznikowi `data-okno` w układzie modułu, czyli skokom nawigacji i sprawdzianom
widoku. Okno wykazu nie zakłada obserwacji telemetrii i nie podaje `windowIds`
przy zakładaniu kolejki, więc jego znacznik nigdy nie opuszcza klienta.

## budowa/klient-poprzedni/src/modele/zapisy-kont.ts

Odczyt i zapis rejestru kont są rozdzielone celowo. Stan rejestru pilnuje tego, co widok wie o kontach, a moduł zapisów pilnuje tego, co dzieje się po zapisie: które wiersze wolno nanieść wprost z odpowiedzi rdzenia, a kiedy trzeba przeczytać wykaz od nowa. Rozdział pozwala zmienić regułę uzgodnienia bez dotykania ścieżki odczytu.

Komenda wskazująca konto domyślne przestawia dwa wiersze naraz, a odpowiedź rdzenia niesie tylko jeden z nich. Dlatego po jej powodzeniu wykaz zostaje przeczytany ponownie: bez tego konto tracące oznaczenie zostałoby w widoku jako drugie konto domyślne swojego rodzaju.

## budowa/klient-poprzedni/src/aplikacja/menu-aplikacji.ts

Menu nosi pasek Centrum dowodzenia oraz pasek Mission Control. Pozycje pochodzą z wykazu pozycji ustawień, czyli z tego samego wykazu, z którego powstaje listwa ustawień strony głównej, ponieważ drugi wykaz rozjechałby się z pierwszym przy dopisaniu pozycji. Skutek naciśnięcia również jest wspólny i rozstrzyga go wykaz skutków ustawień, więc samo menu nie wie, co stanie się po wyborze pozycji.

Rozwijanie wnosi komponent menu rozwijanego z katalogu okien równoległych wraz z obsługą klawiatury, znacznikiem zapowiadającym menu podręczne, zamykaniem po wskazaniu poza obszarem oraz warstwą przybornika. Menu nie zastępuje listwy strony głównej, lecz daje te same pozycje tam, gdzie listwy nie ma: na Mission Control oraz po zwinięciu strony w dół.

## budowa/klient-poprzedni/src/moduly/assistant/okno-command-tools-hub.ts

Okno zamyka obszary szybkich akcji, wywołań narzędzi oraz proaktywności modułu.
Każda z czterech pierwszych zakładek stoi nad rzeczywistą komendą kontraktu: zakładka
Akcje nad `action.list`, czyli katalogiem akcji zasięgu modułu; zakładka Makra nad
`automation.workflow.save`, jedynym trwałym miejscem dla sekwencji kroków; zakładka
Narzędzia i MCP nad `tools.catalog.list` wraz z rodziną `session.tool.*`; zakładka
Rutyny nad `automation.workflow.list`, `schedule.get` oraz `automation.schedule.set`.

Zakładki „Umiejętności" osobno nie ma z rozstrzygnięcia, nie z przeoczenia: kontrakt
trzyma narzędzia i umiejętności w jednym katalogu i rozróżnia je polem `kind` oraz
przedrostkiem źródła. Druga zakładka nad tą samą komendą udawałaby drugie źródło,
a rozróżnienie robi filtr rodzaju w zakładce narzędzi.

Piąta zakładka stoi nad rodzinami `clipboard.*` (trwała historia schowka), `snippet.*`
(słownik skrótów rozwijanych we wszystkich polach platformy) oraz `launcher.hotkey.*`
(skrót globalny wywoływacza). Schowka maszyny Operatora rdzeń nie czyta i zakładka tego
nie udaje: podział ról jest widoczny na ekranie, w pliku `panel-schowka.ts`.

Plik odpowiada wyłącznie za skład okna. Wywołania mieszkają w `zrodlo-narzedzi.ts`,
a każda zakładka ma własny plik obszaru.

Siatka akcji jest tym samym widokiem katalogu, który stoi w Voice Console: drugim
widokiem, a nie drugim źródłem. Obie sięgają `action.list` zasięgu modułu, a naciśnięcie
kafla wypełnia to samo pole polecenia.

Trzy odczyty składania okna idą równolegle, ponieważ katalog akcji, katalog narzędzi
i wykaz rutyn dotyczą trzech różnych rodzin komend i żaden nie warunkuje pozostałych.
Każdy nazywa swoje niepowodzenie w swoim obszarze.

## budowa/klient-poprzedni/src/aplikacja/akcje-ustawien.ts

Pozycja listwy pozbawiona własnego widoku pozostaje klikalna, a jej naciśnięcie
mówi wprost, że ekran jeszcze nie powstał, zamiast otwierać atrapę. Powierzchnia
Always On Display nie ma odrębnego okna i wchodzi jako rozszerzenie boczne,
dlatego jej pozycja prowadzi do powierzchni interakcji, a nie do okna.

## budowa/klient-poprzedni/src/aod/wyciszenie-aod.ts

Wyciszenie Always On Display ma pięć rodzajów: ten plik trzyma trzy z nich jako
stan (czasowe, kontekstowe, klasy zdarzeń); czwarty — tryb cichy — jest trybem
obecności i mieszka w `tryb-obecnosci.ts`; piąty — wyjątek wagi krytycznej —
nie jest wyciszeniem, tylko regułą przebijającą wszystkie pozostałe, i stoi
w regule ujawniania. Plik nie dotyka kanału ani rdzenia: kontrakt nie niesie
wyciszenia nakładki ani jednym polem (`grep -i wycisz shared/contract.go`
znajduje wyłącznie wyciszenie uczestnika tury w module Roundtable i wyciszone
wyzwolenia reguł alarmowych, nic z rodziny `aod.*`), więc stan wyciszenia
jest dziś stanem okna. Magazyn jest podawany wołaczowi (`MagazynWyciszen`),
nie brany z globalnej przestrzeni na sztywno — wzorem
`moduly/studio/widok-nastawy-operatora.ts` — żeby przełożenie zapisu na
rdzeń nie ruszyło ani jednego wołacza. Braki kontraktu nazywa wprost
`wyciszenie-braki-kontraktu.ts`.

Klasa zdarzeń każdego powodu rozpoznanego przez nakładkę przypisuje „pętlę
wstrzymaną i przerwaną" klasie stanu pętli wykonawczej, a „kolejkę
zatrzymaną, zadanie w stanie błędu i zadanie oczekujące dłużej niż próg" —
klasie stanu kolejki zadań. Pozostałe cztery klasy zdarzeń nie mają dziś
w kontrakcie nośnika sygnału (nie ma zdarzenia kontroli jakości,
harmonogramu ani powtarzalności czynności Operatora), więc ich wyciszenie
zapisze się i zadziała z chwilą, w której sygnał wejdzie; menu mówi to
wprost, zamiast udawać, że wycisza coś, co i tak milczy.

Sugestia opisana wobec wyciszenia niesie moduł i sesję jako pola opcjonalne,
bo nie każda sugestia je zna: telemetria rdzenia niesie okno i sesję, modułu
nie niesie wcale, a `AodSuggestion` kontraktu niesie samo okno. Funkcja
rozstrzygająca, które wyciszenie obejmuje sugestię, nie zna wyjątku wagi
krytycznej — wyjątek nie znosi wyciszenia, tylko przepuszcza sugestię
plakietką mimo niego, a rozstrzyga to reguła ujawniania w
`tryb-obecnosci.ts`, w jednym miejscu dla wszystkich trzech rodzajów.

## budowa/klient-poprzedni/src/moduly/assistant/braki-kontraktu.ts

Okno nie stawia w miejscach bez drogi kontraktu przycisku wygaszonego ani
milczącego. Kontrolka zostaje klikalna i odpowiada zdaniem mówiącym, czego
w kontrakcie nie ma; tak samo postępuje `aplikacja/zamiary-pulpitu.ts`.

Każde zdanie wymienia komendę albo pole, którego brakuje, żeby zgłoszenie braku
dało się napisać wprost z ekranu, bez sięgania do innego opracowania.

## budowa/klient-poprzedni/src/moduly/multitasking/strumien-wykonawcy.ts

Bieg naprawczy prowadzi rdzeń w typie `session.StrumienWykonawcy`. Po stronie
klienta leży wyłącznie odkładanie fragmentów; klient niczego nie wybudza ani nie
rozpoczyna.

Odkładane są wszystkie rodzaje fragmentów, nie sama odpowiedź: tok rozumowania,
wywołania narzędzi, ich wyniki oraz pliki. Bez wywołań narzędzi koordynator nie
odróżnia wykonawcy pracującego od zatrzymanego.

Tekst narasta w obrębie jednej tury. Nowa wiadomość zaczyna wynik od nowa,
inaczej przekazanie wyniku niosłoby sklejkę dwóch odpowiedzi wykonawcy.

## budowa/klient-poprzedni/src/moduly/browser/wypis-zrodel.ts

Wykaz źródeł przychodzi do okna komendą `browser.source.list`, a notacja jest
wyłącznie sposobem jego zapisania, dlatego komendy eksportu bibliografii kontrakt
nie niesie i nie musi.

Rodzajem wpisu jest zasób elektroniczny — moduł zbiera strony internetowe i tylko
o nich może tak zaświadczyć. Datą jest czas dodania źródła do wykazu okna,
ponieważ daty publikacji strony rdzeń nie oddaje.

## budowa/klient-poprzedni/src/aod/sekcja-glosu.ts

Pola `audioRef` i `transcript` kontraktu są opcjonalne, a nakładka nagrań nie
tworzy, więc polecenie wydaje się samą transkrypcją.

Sekcja rozmowy jest osobną drogą. Komenda `aod.chat.send` niesie zdanie do okna
rozmowy, natomiast `aod.voice.command` kieruje polecenie do asystenta i ma własne
pole `speak` na odpowiedź syntezą mowy.

## budowa/klient-poprzedni/src/moduly/library/zrodlo-wersji.ts

Podział wobec pliku `zrodlo-biblioteki.ts` biegnie wzdłuż odpowiedzialności:
historia dokumentu należy tutaj, natomiast plik, wyszukiwanie i kolekcje należą
tam. Kształt pozostaje jeden, `ZrodloBiblioteki`, aby okna widziały jedno źródło,
a nie dwa.

## budowa/klient-poprzedni/src/aplikacja/komunikaty.ts

Żaden przycisk w aplikacji nie jest wyszarzany, więc każde naciśnięcie musi dać
Operatorowi odpowiedź. Gdy zamiar nie ma odpowiednika w komendzie kontraktu albo
rdzeń odmawia wykonania, dymek mówi to wprost, zamiast udawać czynność wykonaną.
Milczenie po naciśnięciu byłoby dla Operatora nieodróżnialne od powodzenia.

## budowa/klient-poprzedni/src/moduly/library/ster-modulu.ts

Wybór modułu docelowego jest wyborem z katalogu, a nie polem tekstowym. Komenda
przeniesienia z nieznanym kodem modułu wraca powodzeniem, a rdzeń zakłada okno
z dokładnie tym kodem, więc literówka kończy się oknem, do którego nikt nie
wejdzie. Katalog pozycji pochodzi z komendy wykazu modułów.

Wybór zmienia nastawę we wspólnym stanie modułu. Nastawa jest jedna dla całego
modułu Library, więc zmiana dokonana w oknie eksploratora jest widoczna w oknie
podglądu pliku i odwrotnie.

Pozycja pusta zostaje w wykazie: bierze ona moduł wytwórcy przeniesionego pliku,
a ster tę drogę nazywa wprost, zamiast kazać jej się domyślać z pustego pola.

Gdy wykaz modułów wróci odmową, ster ma samą pozycję pustą, a pod nim stoi powód
odmowy. Milcząca lista z jedną pozycją wyglądałaby jak platforma z jednym
modułem, czyli mówiłaby Operatorowi nieprawdę o stanie rdzenia.

Moduł Library wypada z własnego wykazu pozycji, ponieważ przeniesienie kompletu
do modułu, w którym Operator już stoi, założyłoby drugie okno tego samego modułu
i nie otworzyłoby niczego nowego. Tak samo zawęża swój katalog moduł Design.

Przy odświeżeniu wartość bierze się ze stanu modułu, a nie z kontrolki: funkcja
ustawiająca pozycje utrzymuje poprzedni wybór, o ile pozycja nadal istnieje, ale
prawdą o nastawie jest stan, także wtedy, gdy zmieniło ją drugie okno. Nastawa
wskazująca moduł, którego katalog już nie zna, zniknęłaby po cichu na pozycję
pustą, czyli przeniesienie poszłoby gdzie indziej, niż mówił ster przed chwilą.

## budowa/klient-poprzedni/src/ikony/zrodla-ikon.ts

Wykaz nazw jest w tym pliku jedynym źródłem prawdy — typ `NazwaIkony` wywodzi się
z rejestru, a nie z osobnej listy. Same wiązania nazw z plikami mieszkają w katalogu
`zrodla/`, w podziale na grupy zastosowań zgodnym z porządkiem `ikony/manifest.json`;
ten plik wyłącznie je scala.

Godło marki stoi poza zestawem 82 ikon, bo jest jedynym znakiem o barwach własnych.
Pozostałe odmiany znaku niesie `ikony/marka.ts`.

Nazwy wycofane przy przebudowie zestawu ikon wskazują na następców z tego
zestawu — nie wnoszą żadnego nowego rysunku ani pliku spoza niego. Znikają wraz z ostatnim wywołaniem.

## budowa/klient-poprzedni/src/moduly/developer/braki-kontraktu.ts

Wykaz `KOMENDY` z `shared/contract.ts` jest wykazem komend kontraktu dostępnym
w czasie działania, wytwarzanym z `contract.json`. Zdanie powodu powstaje
z odczytu tego wykazu przy składaniu okna, więc dopisanie komendy do kontraktu
przepisuje je samo. Zdanie wpisane na stałe przestałoby być prawdziwe w dniu
takiej zmiany i nikt nie miałby po czym tego rozpoznać.

Zdanie nie orzeka, czy złożony rdzeń komendę rejestruje. Na to pytanie
odpowiada `katalog-komend.ts`, który pyta rdzeń o jego własny rejestr. Tutaj
brak jest brakiem po stronie kontraktu i tylko o kontrakcie zdanie mówi.

Bliźniaczy mechanizm stoi w module Diagnostics. Biblioteka `komponenty/` nie ma
dziś wspólnego bytu dla obu, a modułowi nie wolno sięgać do wnętrza sąsiada,
więc każdy z nich składa zdanie powodu u siebie.

## budowa/klient-poprzedni/src/aplikacja/ulotnosc-okna.ts

Informację o pamięci sesyjnej modułu bierze `politykaModulu` wyłącznie z profilu
modułu. Powłoka dokłada do niej jedną rzecz, której warstwa rozmowy znać nie
może: co w danym module jest kontekstem roboczym, którego zmiana kończy rozmowę.

Taki kontekst ma jeden moduł — Agents, gdzie kontekstem jest testowany ekspert.
Powłoka jest jedynym miejscem, w którym obie strony są widoczne naraz: widok
modułu powstaje w `przestrzen-modulu`, rozmowa okna w `wiazanie-gniazda`, a żadna
z tych warstw nie zna drugiej.

Moduł nieobjęty tą funkcją również przechodzi przez `politykaModulu` i dostaje
politykę wprost z własnego profilu, bez kontekstu roboczego. Odjęcie pamięci
sesyjnej kolejnemu modułowi czyni jego okno ulotnym bez zmiany w tym pliku.

## budowa/klient-poprzedni/src/aod/indeks.ts

Reszta aplikacji zna stąd dwie czynności. `zaczepAod(kanal, opis?)` buduje
warstwę do osadzenia w powłoce i to ona wnosi pływający awatar widoczny bez
interakcji. `otworzPowierzchnieAod(kanal, opis?)` otwiera powierzchnię interakcji
jako rozszerzenie boczne i jest wołana z listwy ustawień.

Warstwa jest jedna na klienta i żyje między otwarciami powierzchni: kolejka
decyzji dosypuje się ze zdarzeń `progress.changed` oraz `window.state.changed`
także wtedy, gdy kolumna stoi zwinięta, więc plakietka awatara jest prawdziwa
jeszcze przed otwarciem powierzchni.

Kanał podaje się przy pierwszym zaczepieniu. Wywołanie z innym kanałem, po
ponownym połączeniu z rdzeniem, buduje warstwę na nowo, żeby przyszłe komendy
`aod.*` nie szły przez transport, którego już nie ma.

Powierzchnia nie czeka na rdzeń: odczyty `aod.status.get`, `aod.context.get`
i `aod.suggestion` dojeżdżają do otwartej kolumny odpowiedzią, a czynności
`aod.observe.attach`, `aod.observe.detach`, `aod.chat.send` oraz
`aod.voice.command` jadą z sekcji na żądanie. Warstwa niezaczepiona w powłoce
osadza się w korzeniu dokumentu, ponieważ pozycja „Always On Display" listwy
ustawień ma otwierać powierzchnię, a nie odmawiać z powodu kolejności montażu.

## budowa/klient-poprzedni/src/moduly/diagnostics/braki-kontraktu.ts

Plik `shared/contract.ts` niesie wykaz `KOMENDY`, czyli komendy kontraktu dostępne
w czasie działania, wytwarzane z pliku `contract.json`. Zdanie o powodzie powstaje
z odczytu tego wykazu przy składaniu okna. Zdanie wpisane na stałe przestałoby być
prawdziwe w dniu dopisania komendy do kontraktu i nikt by tego nie zauważył.

Zdanie nie orzeka, czy złożony rdzeń komendę rejestruje; jest to pytanie osobne.
Brak jest tu brakiem po stronie kontraktu i wyłącznie o kontrakcie zdanie mówi.

Treść czynności podaje okno wywołujące, ponieważ to okno zna przeznaczenie danej
pozycji inwentarza kontrolek.

## budowa/klient-poprzedni/src/aod/pola-wykazu.ts

Cztery sekcje okna — stan, podpowiedzi, obecność i kontekst — wypisują pola tym
samym wykazem `dl`. Wspólne miejsce trzyma pętlę `dt`/`dd` w jednym egzemplarzu,
dzięki czemu zmiana układu pola obowiązuje we wszystkich sekcjach naraz.

Plik nie zna barw ani odstępów. Nadaje wyłącznie klasy z przedrostkiem `ao-`,
a ich wygląd pokrywa `aod.css` żetonami `--dn-*`.

## budowa/klient-poprzedni/src/moduly/design/polecenie-z-odmowy.ts

Komenda `design.asset.generate` kończy się błędem, gdy w rejestrze brakuje kanału
obrazowego, kanał jest nieczynny albo wskazany kanał okazuje się tekstowy, gdy
brakuje poświadczenia, gdy odpowiedź nie niesie obrazu oraz gdy bajtów nie da się
pobrać. Każda taka odmowa niesie gotową treść polecenia w polu
`details.polecenie`, ponieważ rdzeń składa prompt w całości, zanim cokolwiek
wyśle. Okno pokazuje ten tekst, aby praca nad promptem nie ginęła wraz z odmową
i dała się użyć poza produktem, a zdanie towarzyszące mówi wprost, że obrazu nie
ma i skąd tekst pochodzi.

Kształt pola `details` jest sprawdzany, a nie zakładany: kontrakt opisuje je jako
`unknown`, więc rzutowanie na własny typ byłoby obietnicą bez pokrycia. Brak pola
albo pole innego kształtu daje pusty wynik, nie wyjątek.

## budowa/klient-poprzedni/src/aplikacja/otwarcie-okna.ts

Pierwsze okno komunikacji otwiera uzgodnienie połączenia, prowadzone w pliku
`protokol/uzgodnienie`. Kolejne okna powstają dopiero wtedy, gdy Operator
wprowadza je na scenę, i przechodzą tą samą drogą kontraktu: mają własny
identyfikator, własną historię oraz własne ustawienia.

Nazwa komendy i kształt żądania pochodzą z pliku `shared/contract`, a przekład
opisu okna na treść żądania wykonuje `zamowienieOkna`. Dzięki temu plik
zamówienia nie zna ani jednej nazwy pola kontraktu i nie rozjeżdża się z nim
przy zmianie kształtu żądania.

Rola przekazana w wywołaniu nadpisuje rolę zapisaną w opisie okna, ponieważ
o roli okna na scenie rozstrzyga jego miejsce w figurze koordynator–wykonawca,
a nie konfiguracja budowania pakietu.

## budowa/klient-poprzedni/src/moduly/multitasking/format-zadan.ts

Zapis liczb i czasów stoi w osobnym pliku, ponieważ reguła zapisu czasu trwania jest potrzebna w trzech miejscach naraz: w wierszu podsumowania przepływu, w kolumnie czasu w tabeli agentów oraz w liczniku pozycji zwiniętych. Trzy kopie tej samej reguły rozjechałyby się przy pierwszej poprawce.

Podagent bez znacznika startu otrzymuje brak, a nie czas zerowy. Zero znaczyłoby, że praca ruszyła i nie trwała ani chwili, co jest twierdzeniem innym niż brak pomiaru.

## budowa/klient-poprzedni/src/moduly/browser/wiersz-zrodla.ts

Każda akcja wiersza jest naciskalna — również usunięcie, którego kontrakt nie
niesie. Przycisk odpowiada wtedy nazwaniem brakującej komendy, zamiast znikać
albo gasnąć, ponieważ kontrolka wygaszona bez wyjaśnienia nie mówi, czego
brakuje.

Wiersz nie usuwa pozycji z widoku na własną rękę. Zniknięcie pozycji bez zapisu
w rdzeniu byłoby udawaniem wykonania czynności, która się nie odbyła.

Wykaz źródeł, formularz dodania oraz przekazanie do modułu Research mieszkają
w oknie modułu, a nie w wierszu. Wiersz odpowiada wyłącznie za własną treść
i za wywołanie czynności podanych w `AkcjeZrodla`.

## budowa/klient-poprzedni/src/dostepy/usuniecie-punktu.ts

Usunięcie punktu dostępu unieważnia wszystkie nadania, które się na ten punkt
powoływały, także w oknach, których Operator w danej chwili nie widzi. Zasięg
skutku jest więc szerszy od tego, co widać na ekranie, i to on rozstrzyga
o dwustopniowym przebiegu czynności: pierwsze naciśnięcie nazywa skutek, drugie
go wywołuje.

Potwierdzenie nie korzysta z okna `confirm` przeglądarki, ponieważ okno takie
blokuje wątek dokumentu i nie da się go ubrać w warstwę wizualną platformy.
Potwierdzenie mieszka więc w samym przycisku.

Zamiar wygasa samoczynnie po upływie czasu trwania. Przycisk zostawiony w stanie
potwierdzenia byłby pułapką dla kolejnego naciśnięcia w to samo miejsce,
a wygaśnięciu towarzyszy komunikat mówiący, że punkt pozostał nietknięty.

W drzewie obecne są definicje typów zarówno przeglądarki, jak i środowiska Node,
dlatego typ zegara bierze się z wyniku samej funkcji `setTimeout`, a nie z typu
liczbowego.

## budowa/klient-poprzedni/src/moduly/apps/kanwa-komponentow.ts

Linia zależności jest nazwana, a nie narysowana. Kontrakt niesie zależność jako
wykaz identyfikatorów komponentów, bez współrzędnych i bez kierunku przepływu,
więc rysunek strzałek wymagałby danych, których w kontrakcie nie ma. Kafel podaje
zależność nazwami komponentów, a identyfikator wskazujący komponent spoza kanwy
zostaje w treści wprost, żeby niespójność układu była widoczna zamiast zniknąć.

## budowa/klient-poprzedni/src/mission-control/stan-zrodla.ts

Plik ustala kształt i wartość początkową wsadu złożenia danych pulpitu. Pole
`odczytano` mówi, czy rdzeń już się odezwał, i tym odróżnia pulpit pusty
z powodu braku odpowiedzi od pulpitu pustego z powodu braku danych.

Nazwa kolumny matrycy środowisk pochodzi z odczytu `environment.list`, a nie
z kopii katalogu utrzymywanej w pulpicie. Dzięki temu kolumna znika i pojawia
się wraz ze środowiskiem po stronie rdzenia, bez osobnej pielęgnacji kopii.

## budowa/klient-poprzedni/src/aod/kolumna-aod.ts

Funkcja Always On Display nie ma własnego okna. Całą jej powierzchnią jest
awatar oraz kolumna otwierana jako rozszerzenie boczne po prawej stronie
obszaru roboczego — na pełną wysokość, o regulowanej szerokości. Otwarcie
zwęża kolumny obszaru roboczego, nie przesłania ich i nie zamyka żadnego
z okien komunikacji operacyjnej.

Stąd trzy różnice wobec modala: brak przyciemnienia tła, brak pułapki
ogniska i brak zabranej pracy pod spodem. Kolumna leży na własnej warstwie
graficznej, wyższej niż warstwa modala.

Ten plik wyłącznie składa powierzchnię: rama z biblioteki, stan treści
z biblioteki, sekcje z osobnych plików, sterowanie obecnością z osobnego
modułu i dwa źródła komend — jedno rodziny funkcji, drugie komend cudzych
rodzin potrzebnych kolejce decyzji.

Kolumna subskrybuje dwa zdarzenia rdzenia i zdejmuje subskrypcje przy
rozłączeniu: zmianę postępu, rozgłaszaną do każdego połączenia i niosącą
proces, etap, procent, stan oraz okno; oraz zmianę stanu okna, jedyny żywy
nośnik stanu pętli, ponieważ zmiana postępu nie wypełnia tego pola. Stan
sprzed otwarcia kolumny nadrabia osobny odczyt stanu nakładki.

Subskrypcje żyją między zbudowaniem warstwy a rozłączeniem, nie między
otwarciem a zamknięciem kolumny: proces zmienia stan niezależnie od tego,
czy Operator patrzy. Dzięki temu plakietka awatara jest prawdziwa także
wtedy, gdy kolumna stoi zamknięta.

## budowa/klient-poprzedni/src/main.ts

Punkt wejścia jest wyłącznie kompozycją: logika, typy i obsługiwacze mieszkają
w modułach. Kolejno pyta powłokę natywną o adres rdzenia (`powloka/most-rdzenia`),
wczytuje warstwy wizualne (`aplikacja/arkusze-stylow`), składa aplikację
(`aplikacja/aplikacja`) i uruchamia ją: najpierw widok, potem łączność.

Adres gniazda rdzenia ustala się z lokalizacji dokumentu (`polaczenie/adres-rdzenia.ts`),
a w oknie powłoki natywnej z pakietem osadzonym wychodzi z tego adres `tauri.localhost`,
pod którym nikt nie nasłuchuje. Powłoka zna adres prawdziwy i podaje go poleceniem
`adres_rdzenia`. Pytanie musi paść przed złożeniem aplikacji, ponieważ złożenie zakłada
transport na adresie; zadane później zastałoby gniazdo już otwarte pod adresem ślepym.
Poza powłoką natywną odpowiedzią jest brak, a pytanie kosztuje jedno sprawdzenie
środowiska.

Kolejność dwóch ostatnich wywołań uruchomienia ma znaczenie. Router pokazuje trasę
zapisaną w adresie dokumentu, zanim ruszy transport, więc Operator widzi Centrum
dowodzenia od pierwszej klatki, a nie puste tło czekające na rdzeń. Uzgodnienia nie
rozpoczyna się tutaj: robi to przepływ komunikatów w chwili, gdy transport zgłosi stan
połączenia. Wywołanie stąd dałoby drugie powitanie i podwojenie historii otwarcia.

Uruchomienie bramki idzie po `polacz()`, ponieważ rozpoznanie bramki wysyła komendy
`auth.*` i ekran potrzebuje kanału, który już rusza; wcześniej pokazywałby formularz,
którego nie ma jak wysłać. Ekran montuje się sam, wywołaniem `document.body.append`
wewnątrz `uruchom()`, dlatego punkt wejścia go nie osadza — inaczej przesłona wisiałaby
w dwóch miejscach naraz.

Wywołanie zwrotne `naWejscie` wiąże świeżą sesję z żywym połączeniem: `connection.hello`
przyjmuje token, a rdzeń zapamiętuje, które gniazdo należy do której sesji bramki. To
powitanie nie dubluje powitania z uzgodnienia, ponieważ tamto pada przy nawiązaniu, zanim
Operator poda hasło, więc niesie token sesji poprzedniej albo żaden. Bez powtórzenia po
wejściu rdzeń znałby gniazdo jako niezwiązane przez całe uruchomienie, a komenda
`auth.password.reset` rozłączałaby tego, kto właśnie zmienił hasło. Powitanie jest czystym
odczytem: nie zakłada sesji pracy ani okna.

Scena wejścia staje w tej samej chwili, w której bramka schodzi. Kolejność jest wiążąca:
scena idzie do dokumentu przed zdjęciem przesłony, ponieważ bramka woła `naWejscie`, zanim
usunie swój element; między jednym ekranem a drugim nie mignie wtedy strona główna z pustym
jeszcze wykazem środowisk. Scena stoi warstwę niżej niż bramka, wedle `ladowanie.css`,
i schodzi, gdy pierwszy odczyt strony `home.enter` się domknie albo gdy minie jej własny
kres czekania, żeby rdzeń milczący nie zamienił jej w zasłonę nad produktem.

Pas aktualizacji pojawia się dopiero wtedy, gdy kanał wydań ma wersję nowszą niż
zainstalowana, a przy pierwszym uruchomieniu pyta o to z opóźnieniem, żeby nie konkurować
z uruchomieniem o łącze.

## budowa/klient-poprzedni/src/moduly/agents/archiwum-kontrolki.ts

Plik trzyma wyłącznie budulec panelu archiwum: wiersz wykazu wraz z kontrolką
pozycji oraz przekład odmowy rdzenia na zdanie czytelne dla Operatora. Składanie
panelu i jego czynności należą do `archiwum-ekspertow.ts`, dzięki czemu kontrolki
pozostają wolne od wiedzy o przebiegu pracy panelu.

## budowa/klient-poprzedni/src/mission-control/kafel-liczby.ts

Liczba w kaflu idzie krojem technicznym `--dn-ff-mono`, ponieważ jest daną
techniczną, a nie nagłówkiem. Krój szeryfowy zostaje zastrzeżony dla tytułów
sekcji, przez co czytelnik odróżnia pomiar od nazwy części widoku samym
kształtem znaków.

Wyróżnienie kafla jest zastrzeżone dla miary krytycznej. Etykieta nazywa miarę
wprost i nie bywa zastępowana samym kolorem, dzięki czemu kafel pozostaje
czytelny bez rozróżniania barw.

## budowa/klient-poprzedni/src/aplikacja/menu-operatora.ts

W menu nie ma nazwy własnej, adresu, awatara ani przełączania kont, ponieważ
platforma kont nie prowadzi: komenda `auth.register` konta nie zakłada i adresu
poczty nie przyjmuje. Menu nie udaje więc profilu osobowego, którego pod spodem
nie ma.

Nie ma również pozycji wylogowania. Rdzeń nie odcina komend po wygaśnięciu
sesji, więc taka pozycja mogłaby jedynie skasować zapis sesji i przeładować
aplikację, obiecując Operatorowi skutek, którego nie osiąga.

Obie pozostałe pozycje pochodzą z tego samego wykazu ustawień, co listwa strony
głównej, i idą tą samą drogą skutku. Dwa wykazy rozjechałyby się przy pierwszej
zmianie nazwy albo trasy pozycji.

## budowa/klient-poprzedni/src/moduly/multitasking/wcielenia-analizy.ts

Wybrane wcielenie utrwala komenda `config.set` na poziomie okna analityka. Rdzeń
nie rejestruje komendy `role.update`, więc profilu wcielenia nie zna; pokrycie
kontraktu mierzy `braki-kontraktu.ts`.

Kontrakt nie ma komendy utrwalającej werdykt oceny jako stan: `monitor.status`
czyta stan procesów i niczego nie zapisuje. Ocena analityka jest poleceniem dla
koordynatora, dlatego zgłoszenie niezgodności idzie zwykłym `message.send`.

## budowa/klient-poprzedni/src/moduly/multitasking/zrodlo-zespolow.ts

Zespół i ekspert stoją w jednym źródle, ponieważ zespół kontraktu jest wyłącznie
wykazem identyfikatorów ekspertów. Widok zespołu odcięty od wykazu ekspertów
pokazywałby Operatorowi ciąg identyfikatorów bez nazw, więc obie rodziny komend
obsługuje to samo źródło.

Zespół nie ma pola ról ani powiązania z kolejkami: komenda zapisu przyjmuje
nazwę, opis i skład, i nic poza tym. Gotowy układ ról, powiązań i kolejek daje
się zatem zapisać wyłącznie jako skład, a reszta nie ma w kontrakcie gdzie
zamieszkać. Sekcja mówi o tym Operatorowi wprost, zamiast wpychać te dane w opis
zespołu, gdzie żaden odbiorca by ich nie odczytał.

## budowa/klient-poprzedni/src/modele/stany-odczytu.ts

Pas stanów stoi pod paskiem osi, a nad zakładkami, czyli tam, gdzie widać go
niezależnie od wybranej zakładki. Zawartość zakładek zostaje na miejscu
i pozostaje czynna także przy niepowodzeniu odczytu; podmiana panelu na komunikat
zablokowałaby pracę w pozostałych obszarach sekcji.

Komunikat błędu nie znika po naciśnięciu. Przycisk ponowienia wyzwala odczyt
i zostawia zdanie na miejscu, dopóki sytuacja nie ustanie, dzięki czemu ponowna
próba nie jest brana za powodzenie.

Stan pusty rejestru niesie sam wykaz kont, ponieważ mówi o zawartości rejestru,
a nie o przebiegu odczytu.

## budowa/klient-poprzedni/src/moduly/browser/wykaz-automatyk.ts

Podział odpowiedzialności jest tu taki sam jak w panelach pomocniczych modułu przeglądarki, czyli w plikach wiersz-zrodla.ts oraz wiersz-notatki.ts. Okno składa formularz i prowadzi rozmowę z rdzeniem, a moduł wykazu zamienia wynik tej rozmowy w wiersze widoczne na ekranie. Dzięki temu postać wykazu można zmienić bez dotykania obsługi komend.

## budowa/klient-poprzedni/src/aod/warstwa-aod.ts

Warstwa Always On Display osadza się w powłoce na wysokiej warstwie graficznej, poza obszarem podmienianym przez moduł, więc przełączenie środowiska ani modułu nie zabiera awatara i nie przeładowuje kolumny. Jest jedynym miejscem, które zna naraz cztery rzeczy i wiąże je jedną regułą: kolejkę decyzji, czyli ile sugestii czeka i o jakiej wadze; stan obecności, czyli tryb pełny, cichy albo ukryty i trwające wyciszenie; progi ujawniania, czyli wagę, limit godzinowy i odstęp między dymkami; oraz układ, czyli szerokość otwartej kolumny, o którą odsuwa się awatar i dymek. Skróty klawiszowe są zaczepione w tym samym module, ponieważ tylko tutaj widać wszystkie trzy elementy naraz.

Kolumnę zwężoną poniżej progu czytelności powierzchnia interakcji zwija do dymka kontekstowego, a dymek skraca treść do jednego zdania z działaniem rozwinięcia.

Menu kebab na powierzchni interakcji awatara daje dodatkową akcję wyciszenia bezpośrednio z pozycji awatara. Menu jest tym samym komponentem, który stoi w nagłówku kolumny, a nie nowym wzorem, i pozostaje osobnym wyzwalaczem: kliknięcie pojedyncze awatara nadal otwiera dymek, podwójne nadal otwiera powierzchnię interakcji, a wyciszenie nie wymaga przejścia do kolumny.

Reguła samoczynnego ujawnienia otwiera dymek wyłącznie przy wadze wysokiej, przy nienaruszonym limicie godzinowym i odstępie, poza trybem cichym i poza wyciszeniem. Sugestia krytyczna wyciszona ujawnia się plakietką, bez dymka, a robi to funkcja odświeżająca wygląd, nie reguła ujawnienia.

Decyzja przekazywana regule wyciszenia niesie klasę zdarzenia wynikającą z powodu rozpoznania, kartę sesji z telemetrii oraz moduł z odczytu bieżącego okna w kontekście wyciszenia. Modułu, którego nakładka nie rozpoznaje, funkcja nie podstawia: sugestia bez modułu nie wpada w wyciszenie modułu.

Plakietka liczy sugestie ujawniane, a nie wszystkie: liczba obejmująca sugestie wstrzymane wyciszeniem wybiórczym obiecywałaby coś, czego wyciszenie nie pokaże. Że coś milczy, mówi stan awatara i menu wyciszania. Wyciszenie czasowe chowa plakietkę, z wyjątkiem sugestii krytycznej, która ujawnia się mimo każdego wyciszenia; wtedy plakietka liczy same sugestie krytyczne z tego samego powodu.

Zwężenie obszaru roboczego o szerokość otwartej kolumny idzie wyściółką rodzica warstwy, w której warstwa siedzi, a nie zmianą w katalogu powłoki: funkcja globalna dokłada się do układu, a nie przepisuje go. Klasa zwężenia i szerokość jadą razem, więc powłoka bez otwartej kolumny nie nosi po niej śladu.

## budowa/klient-poprzedni/src/moduly/multitasking/zwiniete-zakonczone.ts

Przepływów zakończonych bywa kilkadziesiąt i wypisane w całości spychałyby poza
krawędź panelu pracę trwającą. Zwinięcie nie jest ukryciem: licznik podaje ich
liczbę, a rozwinięcie oddaje każdy przepływ po nazwie i stanie wraz
z podsumowaniem czasu odcinka oraz liczby podagentów.

## budowa/klient-poprzedni/src/mobile/indeks.ts

Katalog wystawia reszcie aplikacji jedną czynność — `otworzOknoMobile(kanal)` —
oraz typ okna. Okno jest jedno na klienta i żyje między otwarciami, więc
powtórne otwarcie ponawia odczyt zamiast budować okno od nowa.

Kanał podaje się przy pierwszym otwarciu. Wywołanie z innym kanałem, czyli po
ponownym połączeniu z rdzeniem, rozłącza okno dotychczasowe i buduje je na nowo,
żeby przyszłe komendy `mobile.*` nie szły przez transport, którego już nie ma.

## budowa/klient-poprzedni/src/konfiguracja/adres-ustawienia.ts

Ten sam kształt adresu służy dwóm rolom, dlatego mieszka w jednym module.
Pierwszą rolą jest punkt widzenia, czyli miejsce, z którego oglądana jest
konfiguracja i względem którego liczone jest dziedziczenie. Drugą rolą jest
adres zapisu, czyli miejsce, w którym komenda `config.set` zapisze wartość.

Byt pusty znaczy poziom bez bytu: na osi poziomów jest to poziom globalny, a na
osi rozstrzygania — platforma. Dzięki temu adres nie potrzebuje osobnego
znacznika nieobecności bytu.

## budowa/klient-poprzedni/src/moduly/automations/wykaz-petli-widok.ts

Przycisk uruchomienia otrzymuje każda pozycja wykazu, również pętla wyłączona.
Stan pętli stoi w opisie pozycji, a nie w blokadzie wiersza, ponieważ decyzja
o uruchomieniu należy do okna, a nie do widoku wykazu.

Zawężanie wykazu ma dwie drogi: napis szukania oraz przełącznik obejmujący
wyłącznie pętle czynne, który zawęża żądanie do rdzenia polem `enabledOnly`.
Szukanie przebiega po stronie klienta, ponieważ komenda `automation.workflow.list`
przyjmuje jedynie `enabledOnly` oraz `limit`, a pełny wykaz jest już w oknie.
Dopasowanie bierze kolejno nazwę, identyfikator i opis pętli.

Widok oddaje osobne zdanie dla każdego rodzaju pustki: rdzeń bez ani jednej
zapisanej pętli jest innym brakiem niż wykaz zawężony napisem, który do niczego
nie pasuje.

## budowa/klient-poprzedni/src/aod/sekcja-stanu.ts

Wykaz procesów przypiętych niesie sekcja obecności, a nie sekcja stanu,
ponieważ to w sekcji obecności Operator przypina i odpina proces. Wykaz stoi
przy czynności, która go zmienia, dzięki czemu skutek czynności jest widoczny
w tym samym miejscu, w którym została ona wykonana.

Pola puste zostają puste: widok stawia oznaczenie braku zamiast wartości
zmyślonej po stronie widoku. Świeża nakładka bez okna ogniskowanego jest
stanem poprawnym, więc pusta karta sesji i puste okno nie są sygnałem błędu.

Odmowa jednego odczytu jest faktem o jednym odczycie, a nie o całej nakładce.
Pozostałe sekcje stoją na innych komendach i zostają widoczne, przez co odmowa
komendy `aod.status.get` nie wygasza całego widoku.

## budowa/klient-poprzedni/src/moduly/developer/okno-build-output.ts

Build Output i Run & Debug tworzą jedno okno o dwóch częściach zajmujących tę samą kolumnę,
przełączanych pasem zakładek w jej nagłówku — stąd jedna rama, dwie zakładki i jeden kod
okna `build-output` z katalogu rdzenia; Run & Debug nie jest osobnym oknem i osobnego
wiersza katalogu nie dostaje.

Aktualizacja na żywo idzie ze zdarzenia `developer.build.changed`, które okno subskrybuje
bezpośrednio u źródła, mimo że `stan-developer.ts` subskrybuje to samo zdarzenie i budzi
okna przez `stan.naZmiane(...)`: stan tylko rozgłasza, że coś się zmieniło, po filtrze
`windowId`, nie niesie pola `logLine` ani przebiegu z treści zdarzenia, a kontrakt nie ma
komendy, którą dałoby się dogonić stan inaczej. Ta druga subskrypcja bierze więc co innego
niż subskrypcja stanu i nie jest powieleniem tej samej pracy.

Zgłoszenie przebiegu prowadzi do pliku: zgłoszenia niosą pola `path` i `line`, a każde ze
ścieżką ma przejście „Otwórz w edytorze”, które woła `stan.wskazPlik`, a Code Editor otwiera
plik sam.

Log jest ucięty przy otwarciu okna w trakcie przebiegu i rośnie wyłącznie z pola `logLine`
kolejnych zdarzeń, więc okno otwarte po starcie przebiegu — albo przebiegu uruchomionego
przez inne okno tego konta — widzi tylko ogon logu. Ucięcie jest oznaczone wprost w treści
okna, nie zamaskowane.

Szukanie w logu i zawężanie zgłoszeń wagą są czynnościami wyłącznie klienckimi nad
materiałem już zebranym, bo kontrakt nie ma komendy, którą dałoby się dopytać rdzeń
o wiersze pominięte, i widok mówi to wprost. Waga zgłoszenia pochodzi z pola `severity`
oddanego przez rdzeń, nie z rozpoznawania treści wiersza logu.

Zdanie stanu pustego o tym, czym rozporządza rdzeń, składa `katalog-komend.ts` z rejestru
komend wziętego z odpowiedzi `connection.hello` — z tego, co rdzeń rzeczywiście
zarejestrował, a nie z napisu wpisanego na stałe w kliencie.

Powtórzenie uruchomienia komendą `developer.build.run` znaczy to samo żądanie, a nie to,
co aktualnie stoi w polach kreatora — od tego jest osobny przycisk uruchomienia. Gdy okno
nie wysłało jeszcze niczego samo, bo przebieg zaczęło inne okno tego konta, powtarza samo
zadanie z migawki rdzenia i mówi wprost, że parametrów zadania nie zna, ponieważ
`DeveloperBuild` ich nie niesie.

## budowa/klient-poprzedni/src/moduly/design/zrodlo-warsztatow-designu.ts

Metoda `wykonaj` jest jedną drogą na siedemdziesiąt pięć komend, ponieważ
wszystkie idą tak samo: nazwa komendy ze stałych kontraktu oraz treść żądania
złożona z formularza. Nazwa nigdy nie jest napisem wpisanym wprost — przychodzi
z katalogu czynności `czynnosci-warsztatow-designu.ts`, a ten bierze ją ze
stałych kontraktu.

Kształt żądania sprawdza rdzeń i odsyła odmowę walidacji wraz z nazwą pola.
Klient tego nie zastępuje: kontrakt rozstrzyga po stronie rdzenia, a drugie
sprawdzenie po stronie okna byłoby drugą prawdą o tym, co wolno wysłać,
i rozjechałoby się z pierwszą przy najbliższej zmianie kontraktu.

Materiał pochodzi z jednego magazynu zasobów, czytanego komendą
`design.asset.list` — tego samego, z którego czyta Assets Panel. Drugi wykaz
materiału byłby drugim miejscem, w którym ta sama treść żyje.

## budowa/klient-poprzedni/src/konfiguracja/zrodlo-obszarow-sesji.ts

Komenda `config.get` oddaje pojedyncze klucze katalogu ustawień, czyli płaskie
wpisy `ConfigEntry` ze wszystkich poziomów naraz, więc dziedziczenie jednego
klucza rozstrzyga klient w pliku `rozstrzygniecie.ts`. Komenda
`config.effective.get` oddaje obszary konfiguracji sesji już rozstrzygnięte przez
rdzeń, wraz z pochodzeniem obszaru i z rozejściem katalogu roboczego. Klient tego
rachunku wykonać nie może: nie zna pełnej ścieżki bytów ani rejestrów spoza
rodziny `config.*` — rejestru kont, rejestru kanałów, katalogu tożsamości
i nadań dostępu — z których treść obszaru jest składana.

Wynik obu czynności wraca opakowany, więc okno odróżnia obszar bez zapisu od
nieudanego zapytania.

Kształt odpowiedzi `config.effective.get` powstaje w rdzeniu w pliku
`server/internal/core/sesja_konfiguracja_skladanie.go`. Pola `config`, `origins`
i `workingDirectory` wychodzą zawsze, pole `capabilities` wyłącznie na żądanie,
a `unsupportedFields` tą ścieżką nie wychodzi. Sprawdzian obejmuje trzy pola
obowiązkowe i ani jednego więcej, ponieważ pole nieobowiązkowe w sprawdzianie
zamieniłoby zdrową odpowiedź w fałszywą odmowę. Tak samo w odpowiedzi
`config.session.set` pole `unsupportedFields` nie jest wymagane do uznania zapisu
za udany, choć okno pokazuje je, gdy przyjdzie.

Poziom i byt idą w żądaniu tym samym tłumaczeniem adresu, co komenda
`config.set`.

## budowa/klient-poprzedni/src/konfiguracja/rozstrzygniecie.ts

Rachunek dziedziczenia jest w całości po stronie klienta i opiera się wyłącznie
na tym, co przyszło z rdzenia komendą `config.get`: rdzeń oddaje wpisy wraz
z poziomem, bytem poziomu, osią i bytem osi, więc dziedziczenie da się odtworzyć
bez drugiej komendy.

Klient nie zna pełnej ścieżki bytów od środowiska po okno, więc jej nie zgaduje:
punkt widzenia wskazuje pasek u góry okna. Łańcuch takiego punktu składa się
z wpisów globalnych oraz z wpisów zapisanych dokładnie na wskazanym poziomie
i dla wskazanego bytu; zapis na tym samym poziomie, lecz dla innego bytu, dotyczy
kogoś innego i do rachunku nie wchodzi.

Brak zapisu nie jest błędem ani blokadą: obowiązuje wartość domyślna z katalogu,
a wskaźnik mówi to wprost.

## budowa/klient-poprzedni/src/dostepy/sekcja-dostepow.ts

Sekcja dotyczy dwóch bytów: punktu dostępu wraz z nadaniem, czyli tego, do jakich
maszyn i katalogów model sięga w obrębie jednego okna rozmowy, oraz katalogu
roboczego, w którym model zostawia swoje pliki. Środowiska, czyli profilu
widoczności modułów w bocznej nawigacji, sekcja nie dotyka i nie zapisuje do
niego niczego.

Układ jest dwuczęściowy: po lewej wykaz punktów wraz z dodawaniem katalogu
z „Mój komputer", po prawej zbiór nadań okna. Pod nimi stoi obszar katalogu
roboczego, oddzielony, ponieważ mówi o czym innym.

Sekcja pojawia się natychmiast, przed odpowiedzią rdzenia; puste wykazy niosą
zdanie, nie pusty prostokąt.

## budowa/klient-poprzedni/src/moduly/diagnostics/prowenancja-zrodlo.ts

Źródło jest osobne, zamiast dołożenia metod do `zrodlo-diagnostics.ts`, ponieważ
tamten plik należy do rodziny `diagnostics.*`, a prowenancja jest własną rodziną
komend kontraktu. Rozstrzygnięcie odpowiedzi jest jednak to samo — funkcja
`rozstrzygnij` z tamtego pliku — ponieważ rozróżnienie odmowy rdzenia od
odpowiedzi nieczytelnej i od odpowiedzi bez treści ma tu tę samą wagę: zakładka
poświęcona jawności pracy modeli nie może zamilczeć własnego potknięcia.

Piąta komenda rodziny, `provenance.call.replay`, do tego źródła nie należy:
powtórzenie wywołania wydaje pieniądze Operatora i jest czynnością sprawczą
osobnego odcinka, nie odczytem. Zakładka mówi o niej wprost, zamiast wołać ją po
cichu.

Odcinki są w kontrakcie obowiązkowe i puste dla wywołania bez drzewa, więc ich
brak jest odpowiedzią nieczytelną, a nie wywołaniem prostym. Treść wydania śladu
bywa pusta legalnie, gdy zakres nie obejmuje żadnego wywołania, ale format
i licznik muszą przyjść.

## budowa/klient-poprzedni/src/moduly/research/badanie-zdjecie-adnotacji.ts

Komenda `research.annotation.remove` usuwa podświetlenie albo notatkę, a rdzeń
nie ma komendy, która by je przywróciła: adnotacja niesie cytat, komentarz
i kotwicę pozycji, a po zdjęciu nie ma z czego ich odtworzyć. Dlatego czynność
mówi to przed wykonaniem i wymaga drugiego naciśnięcia tego samego chwytu.

Uzbrojenie jest wiedzą modułu, a nie okna: adnotacja wskazana w widoku czytania
i chwyt naciśnięty w panelu akcji tego samego okna dotyczą jednego badania,
a moduł prowadzi jedno badanie na sesję — stan trzyma `stan-badania.ts`. Osobne
uzbrojenie dla każdego okna pozwoliłoby uzbroić w jednym oknie, a zdjąć
w drugim, czyli bez ostrzeżenia w miejscu naciśnięcia.

Zapowiedź nie jest okienkiem dialogowym i nie odbiera klikalności: pierwsze
naciśnięcie oddaje zdanie odmowy z powodem, drugie wykonuje. Adnotacja wskazana
inna niż uzbrojona uzbraja od nowa, więc zdjęciu zawsze towarzyszy ostrzeżenie
nazywające rzecz zdejmowaną.

## budowa/klient-poprzedni/src/modele/sekcja-modeli.ts

Oś wskazania jest wspólna dla całej sekcji. Ustawienia bytu, tożsamość i podgląd
promptu mówią o tym samym bycie: wskazanym modelu albo wskazanym koncie. Trzy
osobne wybory osi w trzech panelach dałyby trzy rozbieżne wskazania, dlatego
pasek osi jest jeden i stoi nad zakładkami.

Wykaz kont zasila wybór osi: konta znane rdzeniowi stają się podpowiedziami bytu
osi `account`, a identyfikatory modeli składa osobne źródło z rejestru kanałów
oraz z modeli domyślnych kont. Podpowiedź nie zamyka pola — byt spoza wykazu
wolno wpisać wprost. Sekcja otwiera się przed odpowiedzią rdzenia, a każdy
z czterech obszarów ma własny komunikat na wypadek braku danych i żaden nie
blokuje pozostałych.

Identyfikatory modeli trwają między odczytami, ponieważ pochodzą z rejestru
kanałów czytanego osobną komendą. Zmiana rejestru kont odświeża podpowiedzi kont
i zostawia modele bez zmian: podanie pustej listy skasowałoby podpowiedź bez
powodu.

Odczyt sekcji nie rozsyła osi ponownie. Oś zmienia wyłącznie pasek u góry,
a jego zmiana już ją rozesłała; powtórzenie rozesłania przy odczycie kazałoby
rdzeniowi przeczytać te same zapisy i tę samą nakładkę drugi raz pod rząd.

## budowa/klient-poprzedni/src/dostepy/nazwy-dostepow.ts

Plik jest jedynym miejscem, w którym wartość wyliczenia kontraktu zamienia się w zdanie
po polsku. Wartości pochodzą wyłącznie ze stałych `shared/contract`, a ten plik dokłada
do nich warstwę językową i nic więcej. Rozsypanie tych zdań po widokach dałoby dwie nazwy
tego samego stanu.

Wartość spoza wyliczenia nie jest błędem i niczego nie wygasza: wraca jako własny napis,
żeby widoczne było to, co przysłał rdzeń, zamiast pustego miejsca.

## budowa/klient-poprzedni/src/dostepy/klucze-katalogu.ts

Dostęp to nie katalog roboczy. Model może mieć wgląd w katalog wskazany nadaniem
dostępu, a swoje pliki zostawiać w katalogu sesyjnym wewnątrz miejsca
instalacji. Są to dwa niezależne ustawienia i dwa osobne obszary tej sekcji,
mylenie ich prowadziłoby do nadania dostępu tam, gdzie potrzebny jest katalog
roboczy.

Katalog roboczy jest zwykłym ustawieniem: idzie komendami odczytu, zapisu
i przywrócenia wartości domyślnej konfiguracji, przez ten sam rezolwer zasięgów
co reszta ustawień. Osobnej komendy dla niego nie ma.

Katalog roboczy powstaje domyślnie w miejscu instalacji aplikacji głównej.
Klient tej ścieżki nie zna i jej nie zgaduje: mówi Operatorowi, skąd się ona
bierze, zamiast pokazywać zmyśloną ścieżkę jako fakt.

Brak wiersza w katalogu ustawień nie może zabrać możliwości ustawienia katalogu
roboczego, dlatego obszar niesie własne definicje zastępcze. Definicja
z katalogu ma pierwszeństwo, a zastępcza wchodzi wyłącznie wtedy, gdy katalog
milczy, i jest zbudowana z tych samych wartości co odpowiadające jej definicje
rdzenia.

Zasięgi katalogu roboczego są węższe niż komplet zasięgów konfiguracji: katalog
roboczy jest własnością instalacji, środowiska, projektu, sesji oraz okna, a nie
pary modułów.

## budowa/klient-poprzedni/src/moduly/agents/stan-agentow.ts

Pięć okien modułu — Agent Builder, Model Configuration, Skills Manager,
Connectors Manager, Permissions Center — pracuje na jednym zbiorze stanu. Gdyby
każde okno prowadziło własny wykaz i własny wybór, zmiana modelu bazowego nie
odświeżałaby biblioteki, a przypisanie umiejętności dotyczyłoby innego eksperta
niż ten pokazany w edytorze. Poza własnym działaniem stan odświeża wyłącznie
zdarzenie rdzenia; odpytywania w pętli tu nie ma.

Ogłoszenie zmiany stanu ustawia najpierw eksperta testowanego, dopiero potem
woła słuchaczy okien. Czat modułu jest czatem testowym bez pamięci sesyjnej:
rozmowa znika przy zamknięciu okna i przy zmianie testowanego eksperta, a okno
rozmowy stoi obok widoku modułu, więc o przełączeniu dowiaduje się wyłącznie
tym wywołaniem. Kolejność jest treścią: rozmowa zdąży się wyczyścić i nazwać
powód, zanim okna modułu przerysują się na nowego eksperta. Odwrotna kolejność
pokazywałaby przez moment czat poprzedniego eksperta w oknach już opisanych
nazwą następnego.

## budowa/klient-poprzedni/src/moduly/library/panel-akcji.ts

Panel nie jest zaszytym wykazem: pozycje przychodzą z katalogu rdzenia komendą
`action.list` o zasięgu modułu, a każda niesie własny kod, którym wykonuje ją
`window.action`. Zaszycie pozycji po stronie klienta byłoby drugą kopią
katalogu, rozjeżdżającą się z rdzeniem przy każdej zmianie po jego stronie.

Czynności, którym kontrakt nie przypisał komendy — eksport, archiwizacja, kosz,
wykrywanie duplikatów, udostępnienie odnośnikiem oraz porównanie — stoją obok
jako przyciski klikalne. Naciśnięcie takiego przycisku nic nie wysyła i nazywa
brak, zamiast zostawiać Operatora przy kontrolce wygaszonej albo milczącej.

Powodzenie wywołania oznacza wynik akcji, a nie samo przyjęcie zgłoszenia:
akcję, której rdzeń nie wykonuje, rdzeń odrzuca wprost kodem `conflict` wraz
z kodem komendy do wywołania, a panel pokazuje tę odmowę.

## budowa/klient-poprzedni/src/moduly/agents/okno-connectors-manager.ts

Wykaz mostów protokołu MCP pochodzi z tego samego katalogu punktów dostępu,
z którego rdzeń składa plik konfiguracji procesu modelu. Okno nie zakłada
drugiego rejestru mostów i nie zna ani jednego adresu maszyny: podaje kod
punktu dostępu, resztę wie rdzeń.

Katalogiem rozszerzeń zarządza osobne okno modułu, obsługujące komplet pięciu
komend tego katalogu; tutaj stan tych komend bierze się na żywo z bytu
pokrycia, a samo okno wykonuje wyłącznie podłączenie konektora i wtyczki.

Licznik narzędzi stoi także w tym oknie, ponieważ serwer narzędzi czyta
identyfikatory umiejętności i konektorów eksperta jako jeden zbiór kodów —
podłączenie konektora zmienia tę samą liczbę, którą pokazuje Agent Builder.

Wtyczki stoją osobno od konektorów: wykaz wtyczek ma własny panel i własne
komendy podłączenia oraz usunięcia. Konektor jest drogą do usługi zewnętrznej,
wtyczka katalogiem rozszerzeń powłoki — to dwa różne mechanizmy rozszerzania
eksperta.

## budowa/klient-poprzedni/src/moduly/assistant/wiersz-dziennika.ts

Wszystkie trzy czynności wiersza mają drogę w kontrakcie i wszystkie trzy naprawdę
coś robią. „Odtwórz przebieg" oddaje zamiar oknu, „Odsłuchaj nagranie" pobiera bajty
spod odnośnika wpisu komendą `speech.audio.fetch` i odtwarza je w karcie, a „Oznacz
jako ważne" zapisuje wyróżnienie w rdzeniu komendą `assistant.activity.flag`. Znacznik
przeżywa odświeżenie wykazu, ponieważ ma gdzie zamieszkać po stronie rdzenia.

Wpis bez odnośnika nagrania nie jest brakiem produktu, tylko wpisem tekstowym. Przycisk
odsłuchu mówi to wprost, zamiast milczeć albo znikać z wiersza.

## budowa/klient-poprzedni/src/moduly/multitasking/widok-strumienia.ts

Okno koordynatora pokazuje cztery rodzaje treści, a nie samą odpowiedź modelu:
tok rozumowania, wywołania narzędzi, wyniki oraz pliki. Pliki są wyróżnione
znacznikiem osobnym, ponieważ stoją obok pozostałych trzech rodzajów jako
czwarta rzecz widoczna dla koordynatora, a nie jako odmiana wyniku narzędzia.

Rodzaje fragmentów pochodzą z kontraktu przez katalog warstwy rozmowy, więc
moduł nie wprowadza własnego nazewnictwa rodzajów ani własnych wartości.

Wiersz licznika niesie powód zatrzymania biegu wprost z pola `LoopState.stopReason`.
Bez tego wiersza bieg zatrzymany byłby nieodróżnialny od biegu bezczynnego,
a koordynator nie miałby przesłanki do rozstrzygnięcia o wznowieniu.

## budowa/klient-poprzedni/src/moduly/assistant/wiersz-zlecenia.ts

Wszystkie sterowania wiersza zostają klikalne niezależnie od stanu zlecenia. Wstrzymania
zlecenia już wykonanego widok nie blokuje — odpowiada na nie rdzeń, a odpowiedź trafia do
wiersza odpowiedzi okna. Wygaszona kontrolka kazałaby zgadywać, czy przycisk nie działa, czy
tylko nie odpowiada.

Pole priorytetu wysyła sterowanie `none`, choć rdzeń zapisuje priorytet wyłącznie przy
sterowaniu innym niż `none`. Zmiana kolejności obsługi nie jest zmianą stanu zlecenia,
a kontrakt nie zna wartości `AssistantActionControl` znaczącej sam zapis priorytetu — są
wyłącznie `none`, `pause`, `resume`, `cancel` oraz `retry`, a komenda
`assistant.action.status` jest jedyną komendą obszaru niosącą pole `priority`.

Podszycie się pod `pause` albo `retry` po to, żeby przemycić priorytet, przestawiłoby stan
zlecenia, o który nikt nie prosił. Pole zostaje więc przy `none`, a rozbieżność nazywa wiersz
odpowiedzi okna składany przez `skutek-sterowania.ts`.

## budowa/klient-poprzedni/src/moduly/apps/zrodlo-punktow-dostepu.ts

Serwer MCP jest w kontrakcie pozycją katalogu rozszerzeń, a most, którym się do
niego dochodzi — punktem dostępu `AccessPoint`. Oba byty łączy pole
`Extension.accessPointId`, więc okno integracji potrzebuje obu wykazów:
z pierwszego bierze rodzaj i stan włączenia, z drugiego adres mostu, korzenie
oraz wynik ostatniego sprawdzenia.

Komenda `access.point.check` jest jedyną drogą kontraktu do testu połączenia
i do monitora zdrowia: odpowiada, czy punkt odpowiada, jakie korzenie potwierdza
i co poszło nie tak. Nie jest to sprawdzenie kondycji samego rozszerzenia —
sprawdzeniu podlega most, a nie usługa stojąca za nim. Korzenie i szczegół
bywają puste także przy sprawdzeniu udanym, dlatego całą odpowiedź niesie pole
stanu.

Źródło nie ma własnego stanu ani nie zna okna. Wykazy trzyma `stan-rozszerzen.ts`,
żeby cztery okna patrzyły na jeden zbiór danych.

## budowa/klient-poprzedni/src/konfiguracja/zasiegi.ts

Poziom zasięgu mówi, jak wąsko obowiązuje wartość: okno jest najwęższe i wygrywa,
globalny najszerszy i przegrywa z każdym innym.

Oś mówi, dla czego wartość obowiązuje: dla platformy, dla wskazanego modelu albo dla
wskazanego konta. Oś pominięta we wpisie znaczy `platform`.

## budowa/klient-poprzedni/src/moduly/library/wiersz-wersji.ts

Osiągalność treści bierze się z odpowiedzi rdzenia opracowanej w `dostepnosc-tresci.ts`,
a nie z obecności sumy kontrolnej w opisie wersji. Pole `checksum` bywa obecne przy wersji,
której treści rdzeń nie oddaje, i nieobecne przy wersji, którą oddaje w całości, więc jako
przesłanka osiągalności myli.

Odpowiedź rdzenia dotyczy treści, którą dokument niesie jako bieżącą, dlatego werdykt siada
wyłącznie na wierszu bieżącym: kontrakt nie ma komendy pytającej o treść wersji niebieżącej.
Wiersz niebieżący nie jest zatem „nieznany" z braku odpowiedzi — po prostu nie był pytany.

Werdykt `brak` znaczy, że nie ma do czego wracać. Werdykt `odmowa` znaczy, że rdzeń treści
nie oddał, ale też o niej nie orzekł, więc nie odbiera przycisku przywrócenia. Werdykt
`odwolanie` znaczy wskazanie miejsca zamiast bajtów.

## budowa/klient-poprzedni/src/konfiguracja/panel-obszarow-sesji.ts

Zasięg jest widoczny w dwóch miejscach, bo mówi o dwóch rzeczach. Nagłówek
sekcji mówi, dla kogo liczona jest konfiguracja i pod jaki adres pójdzie zapis;
oba biorą się z punktu widzenia ustawionego paskiem u góry okna, więc panel nie
stawia drugiego selektora zasięgu. Plakietka przy obszarze mówi natomiast, skąd
wartość przyszła: z którego rejestru, z którego poziomu i z której osi.

Ładowanie i odmowa idą tym samym pasem stanów odczytu, którym idą stany
katalogu, a pustka idzie stanem pustym panelu kategorii. Żadnego z tych trzech
stanów panel nie buduje sam.

Zapis wybranych obszarów zawsze kończy się zdaniem odpowiedzi. Wybór pusty nie
idzie do rdzenia, ponieważ komenda z pustym wykazem obszarów nie zapisałaby
niczego, a milczący przycisk wyglądałby jak awaria.

## budowa/klient-poprzedni/src/moduly/library/administracja-repozytorium.ts

Siedem rodzin zdolności eksperckich stoi w jednym panelu, bo wszystkie są
sterowaniem repozytorium jako całością, a nie pracą nad wskazanym plikiem.
Rozbicie ich na siedem osobnych paneli dałoby siedem miejsc, w których
Operator szukałby tego samego: gdzie ustawia się zachowanie biblioteki.

Każda czynność panelu mówi wynikiem, co zaszło. Wykazy odczytują się po
wykonaniu czynności, a nie w pętli — repozytorium nie odpytuje się samo.

## budowa/klient-poprzedni/src/moduly/agents/zrodlo-zespolow.ts

Źródło jest osobne od `zrodlo-agentow.ts`, ponieważ zespół nie jest polem
eksperta ani jego wersją: jest własnym bytem kontraktu `Team`, z własnym
zdarzeniem `team.changed` i własnym cyklem życia. Włączenie tych czterech
komend do źródła obszaru `agent.*` dałoby jeden plik o dwóch
odpowiedzialnościach.

Nazwy `Command.*` kończą się w tym pliku. Panel oraz widok składu wołają
czynności interfejsu nazwane po polsku, bytami dziedziny, i kontraktu nie
dotykają. Przemianowanie komendy w `contract.json` przerywa wtedy kompilację
w jednym pliku, a nie w trzech widokach.

Żadna czynność źródła nie rzuca wyjątku. Odmowa wraca polem `blad` wyniku
i pokazuje się w wierszu odpowiedzi tego panelu, w którym czynność wywołano.

Rozróżnienie między założeniem nowego zespołu a zmianą istniejącego opiera się
wyłącznie na obecności identyfikatora, więc pole idzie do żądania tylko wtedy,
gdy naprawdę jest. Wysłanie pustego napisu byłoby dla rdzenia wskazaniem
zespołu o nazwie pustej, a nie brakiem wskazania.

Opis jest w kontrakcie polem opcjonalnym, więc pusty nie trafia do żądania:
przy zmianie zespołu wysłanie pustego napisu skasowałoby opis już zapisany,
choć formularz o skasowanie nie prosił.

Zdanie o zespole podaje liczebność składu, ponieważ jest ona jedyną rzeczą, po
której Operator poznaje na wykazie różnicę między zespołem a jego kopią świeżo
powieloną.

## budowa/klient-poprzedni/src/dostepy/przelacznik-trybu.ts

Ostrzeżenie o zapisie jest częścią przełącznika, ponieważ dotyczy skutku wyboru,
a nie sąsiedztwa na ekranie. Zapis na maszynie chronionej pokazuje pełne zdanie pod
przełącznikiem, a nie w podpowiedzi pod kursorem: skutek zapisu na produkcyjnej
platformie LEX jest nieodwracalny, więc zdanie musi być widoczne bez najeżdżania
wskaźnikiem.

Treść ostrzeżenia liczy się z nazwy maszyny podanej w punkcie dostępu, dlatego
nanoszenie nowego punktu przelicza je od nowa.

## budowa/klient-poprzedni/src/modele/rodzaje-kont.ts

Plik nie wymienia żadnego dostawcy. Dostawca jest w kontrakcie wartością danych,
a nie typem kodu, i wpisuje się go w polu tekstowym formularza konta. Rodzaj
konta pozostaje natomiast wyliczeniem kontraktu, więc jego wartości pochodzą
wyłącznie ze stałych `AccountKind`.

Rodzaj spoza kontraktu nie gaśnie — nazwa gotowa do wydruku pokazuje wtedy
własny kod wartości, żeby było widać, co przyszło z rdzenia. Ta sama zasada
obowiązuje kontrolkę wyboru: napis pusty znaczy brak ograniczenia, a napis
spoza kontraktu również, ponieważ zawężenie wykazu do wartości, której rdzeń nie
zna, dałoby wykaz pusty bez powodu.

## budowa/klient-poprzedni/src/moduly/assistant/stan-probny.ts

Sprawdzianów modułu są dwa — stany okien oraz wykonywanie pracy — a atrapa stanu modułu jest
jedna. Skopiowana do obu plików rozjeżdżałaby się przy każdym nowym polu interfejsu i jeden
ze sprawdzianów badałby wtedy stan, którego moduł już nie ma.

Plik nie należy do produktu: sięgają po niego wyłącznie pliki sprawdzianów, więc punkt wejścia
aplikacji go nie wciąga.

Atrapa nie sięga do rdzenia i nie udaje, że sięga: każda droga wywołania rzuca wyjątkiem,
dopóki sprawdzian jej nie obsadzi. Sprawdzian, który przypadkiem wywoła komendę, dostaje przez
to błąd, a nie ciszę. Dotyczy to również trzech źródeł dobudowanych obok rdzenia modułu —
atrapa oddająca pustą odpowiedź udawałaby wynik, którego nie ma.

## budowa/klient-poprzedni/src/moduly/agents/pola-zalezne.ts

Adapter dostawcy deklaruje dla każdego pola obszaru `model`, czy obsłuży je w całości,
częściowo, czy wcale, i podaje powód. Okno pokazuje tę deklarację wprost, więc widać,
że kanał wiersza poleceń zignoruje `samplingTemperature`, zanim wartość zostanie wpisana.

Wykaz zbudowany z deklaracji, a nie z listy wpisanej w kodzie, nadąża za dostawcą:
nowa wersja jego programu zmienia zestaw pól bez wydania nowego klienta.

## budowa/klient-poprzedni/src/moduly/multitasking/indeks.ts

Scena modułu obsadza cztery okna trzech typów. Okno koordynatora planuje bieg
i steruje nim, lecz nie tworzy produktu końcowego. Dwa okna wykonawców pracują
równolegle i nie zarządzają procesem. Okno analityka zestawia wyniki obu
wykonawców. Wszystkie cztery okna stoją na jednym stanie i jednym źródle biegu,
ponieważ przy osobnych źródłach koordynator widziałby inny przebieg niż
wykonawca, którym steruje.

Wykaz komend rdzenia powstaje raz na scenę i jest wspólny czterem oknom. Cztery
osobne zapytania o ten sam wykaz dałyby cztery odpowiedzi, które mogłyby się
rozejść. Zapytanie wychodzi przed pierwszym odświeżeniem okien, dzięki czemu
kontrolki bez pokrycia w wykazie oznaczają odpowiedź w drodze, a nie brak
komendy w rdzeniu.

Panel obsady stoi przed oknami ról, ponieważ role są przypisaniem okien sesji,
a nie osobnymi bytami, i bez przydziału nie ma czego otworzyć. Panel jest
jedynym miejscem nadania tych ról, a jego meldunki idą do jego własnego stanu
treści, więc odmowa rdzenia nie ginie po drodze. Układ sekcji paneli pochodzi
z powłoki, ponieważ dotyczy okna, a nie dziedziny modułu.

Panel zadań w tle stoi pod oknami ról, ponieważ opisuje pracę już przez nie
zleconą, a podagenci należą do okien wykonawców. Panel jest jeden na kartę
sesji, ponieważ komenda wykazu podagentów przyjmuje identyfikator sesji i dwa
panele pytałyby dwukrotnie o to samo. Oba okna wykonawców korzystają z tego
samego źródła podagentów co panel, więc drugie źródło byłoby kopią.

Nasłuchy odpina się przed zdjęciem sceny, ponieważ zdarzenie przyjęte po
usunięciu elementu odświeżałoby widok już nieistniejący.

## budowa/klient-poprzedni/src/moduly/diagnostics/telemetria-procesow.ts

Zakres czasu nie dotyczy tego odczytu: `monitor.status` i `monitor.subscribe`
nie mają w kontrakcie pól `fromTime` ani `toTime`, więc telemetria jest zawsze
stanem bieżącym rejestru procesów, a zawężanie jej zakresem po stronie okna
byłoby zawężaniem pozorowanym. Obserwacja zakłada się raz, odczyt powtarza
się: pierwsze pytanie idzie `monitor.subscribe` (zapisuje okno na telemetrię),
każde następne `monitor.status` (czyta i niczego nie zapisuje) — pole
`windowId` znaczy w tych dwóch komendach co innego: w pierwszej jest oknem
obserwatora, w drugiej sitem procesów, więc do odczytu nie idzie wcale.

Na żywo idzie `progress.changed`. Zdarzenie niesie stan, etap i stopień
ukończenia, lecz nie niesie czasu zmiany, dlatego wiersz odświeżony
zdarzeniem pokazuje liczby ze zdarzenia, a czas nadal ten z odczytu, wraz ze
zdaniem, skąd się bierze — podstawienie czasu klienta w miejsce czasu rdzenia
byłoby wpisaniem do telemetrii wartości, której rdzeń nie orzekł.

Eksport telemetrii stoi przy danych, a nie w pasku akcji okna: potwierdzenie
eksportu wypisuje się w tej samej zakładce, w której widać eksportowany
materiał, bo potwierdzenie na zakładce zasłoniętej nie byłoby potwierdzeniem.
Wytwórnia sama nie czyta — pierwszy odczyt zleca złożenie modułu wywołaniem
`odswiez()`, tak samo jak w panelu rekomendacji, żeby odczyt w wytwórni obok
odczytu ze złożenia nie dał dwóch żądań na jedno zmontowanie okna.

Po odmowie odczytu wykaz poprzedni też odchodzi: skoro odczyt się nie udał,
wykaz przestaje być bieżący, a eksport oddałby go, jakby wciąż nim był.

## budowa/klient-poprzedni/src/aplikacja/rejestr-modulow.ts

Nawigacja bierze wykaz modułów z rdzenia komendami `module.list`
oraz `environment.enter`, więc powłoka poznaje kody modułów dopiero w czasie
działania i nie może o nich nic zakładać. Bez rejestru każdy moduł musiałby
wpinać się sam, a przestrzeń robocza trzymałaby własną kopię katalogu modułów.
Moduł nieobecny w rejestrze nie jest błędem: rdzeń zna moduły, których klient
jeszcze nie zbudował, a taki moduł dostaje stan pusty z paska uczciwości
przestrzeni roboczej zamiast martwego kliknięcia.

Element widoku i wczytanie są rozdzielone, ponieważ sesja bywa nieznana w chwili
wyboru pozycji nawigacji: przestrzeń robocza odkłada wtedy wejście i ponawia je
po otwarciu okna przez rdzeń. Gdyby wytwórnia oddawała sam element, każdy moduł
obsługiwałby tę samą zwłokę u siebie.

Metoda `zamknij` jest szwem rozbiórki widoku i pozostaje opcjonalna. Powłoka jej
nie woła, ponieważ `przestrzen-modulu.ts` woła `utworzWidok` raz na kod modułu
i oddaje ten sam widok przy każdym powrocie, a `obszar-roboczy.ts` nie zdejmuje
elementu z drzewa, tylko przestawia atrybut `hidden`. Liczba żywych subskrypcji
jest przez to ograniczona z góry liczbą modułów i nie rośnie wraz z liczbą
przełączeń. Gdy w powłoce pojawi się pierwsza granica życia widoku, wołanie
`zamknij` musi towarzyszyć zdejmowaniu elementu w `obszar-roboczy.ts`.

Wpis MultitaskingAI niesie kod `multitaskingai`, który jest kodem środowiska,
a nie modułu. Nawigacja tego środowiska podaje sekcje orkiestracji ze stałej
`SEKCJE_ORKIESTRACJI` w `powloka/srodowiska.ts`, więc wyszukanie opisu modułu
po kodzie środowiska nie daje wyniku.

Moduł Diagnostics montuje się od razu, ponieważ większość jego komend nie
wymaga identyfikatora okna.

## budowa/klient-poprzedni/src/moduly/developer/zrodlo-probne.ts

Jedna atrapa dla całego modułu zastępuje atrapę przepisywaną w każdym sprawdzianie z osobna:
rozjazd z umową `ZrodloDeveloper` przerywa wtedy kompilację, zamiast rozjeżdżać sprawdziany
po cichu. Plik służy wyłącznie sprawdzianom tego modułu i nie jest wciągany przez żadne okno.

## budowa/klient-poprzedni/src/moduly/developer/stan-developer.ts

Cztery okna modułu pracują nad tym samym repozytorium: drzewo projektu wskazuje
plik, edytor kodu go otwiera i zapisuje, panel repozytorium podstawia ścieżkę
wskazaną do pola czynności i sam wskazuje ścieżki ze swojego wyniku, a okno
wyniku budowania wskazuje pliki ze zgłoszeń przebiegu. Gdyby każde okno trzymało
własną ścieżkę i własny korzeń, wskazanie pliku w drzewie nie dotarłoby do
edytora.

Identyfikator okna modułu należy do stanu wspólnego, ponieważ kontrakt wymaga go
w każdej komendzie obszaru, a okna mają go podawać identycznie. Rozbieżność
rozdzieliłaby ich pracę między dwa katalogi robocze rdzenia.

Stan nie wywołuje komend w zastępstwie okien: przechowuje wybór i rozgłasza
zmianę, a odczyt i zapis pozostają w oknach, ponieważ tylko one prowadzą stan
ładowania i obsługę odmowy.

Ostatni plik oddany przez rdzeń nie jest tym samym co ścieżka wskazana:
wskazanie biegnie natychmiast po wskazaniu węzła w drzewie, a plik przychodzi
dopiero odpowiedzią rdzenia. Zgodność obu rozstrzyga porównanie ścieżki pliku ze
ścieżką wskazaną i dlatego plik niesie własną ścieżkę. Edytor kodu czerpie stąd
identyfikator wersji sprzed zapisu, a drzewo projektu odróżnia węzeł wskazany od
węzła wczytanego do edytora.

Treść pliku pochodząca z rdzenia nie jest treścią pola edycji. Pole bywa
zmienione i niezapisane, a rdzeń po zapisie nie zawsze odsyła treść, więc
znacznik czystości pola należy do edytora kodu.

Wskazanie nowej ścieżki nie zeruje pliku: plik niesie własną ścieżkę, więc nie
sposób pomylić go ze wskazaniem, a wyzerowanie usunęłoby jedyną wiedzę o tym, co
leży w polu edytora.

Przyrost budowania przychodzi także z pracy innego okna tego konta, dlatego stan
przyjmuje wyłącznie przyrost dotyczący własnego okna modułu. Budowanie tworzy
pliki, więc po nim drzewo i plik w edytorze bywają nieaktualne i wymagają
ponownego odczytu.

## budowa/klient-poprzedni/src/moduly/browser/material-sesji.ts

Zbiór materiału nie jest odbiciem stanu rdzenia, ponieważ odbijać nie ma czego:
zapis pozycji jako wytworu sesji ma komendę `browser.artifact.add`, natomiast
komendy odczytu wykazu wytworów okna kontrakt nie niesie. Pozycją jest migawka,
którą rdzeń oddał na żądanie, a po przeładowaniu karty wykaz zaczyna się od nowa;
same migawki zostają w rdzeniu pod swoimi identyfikatorami. Okno mówi o tym
wprost, zamiast obiecywać trwałość, której nie ma.

Rodzaj pozycji rozstrzyga zawartość odpowiedzi, a nie zamówienie okna. Rdzeń
pobiera stronę bez uruchamiania przeglądarki (`przegladarka_pobieranie.go`), więc
migawka zamówiona ze zrzutem potrafi przyjść bez niego — pozycja nazywa się
wtedy archiwum albo treścią. Rozpoznanie po odpowiedzi nazwie ją zrzutem samo,
gdy obraz zacznie przychodzić.

Monitor zmian stoi na adresie, a nie na migawce: sprawdzenie polega na ponownym
pobraniu tej samej strony i zestawieniu jej treści z zapamiętaną. Adres jest więc
jedynym rozsądnym kluczem monitora, bo dwa monitory tej samej strony pilnowałyby
dokładnie tego samego. Założenie monitora na adresie już pilnowanym oddaje
wartość fałszywą, żeby okno powiedziało o tym wprost zamiast meldować czynność,
która się nie odbyła.

## budowa/klient-poprzedni/src/aod/zrodlo-decyzji.ts

Źródło stoi osobno od `zrodlo-komend.ts`, ponieważ tamto niesie wyłącznie rodzinę
komend nakładki. Sklejenie obu zatarłoby granicę między komendami własnymi nakładki
a komendami cudzych obszarów, po które nakładka sięga.

Odczyt procesów idzie bez argumentów świadomie: nakładka pyta o wszystkie procesy,
które rdzeń zna, bo proces czekający na decyzję nie musi należeć do sesji tego klienta.

Konfiguracja okna otwierana jest na zasięgu `window`, czyli poziomie najwęższym,
który wygrywa z każdym szerszym. Konfiguracja koordynatora dotyczy tego jednego okna
i nie ma sięgać ustawień sesji ani konta.

## budowa/klient-poprzedni/src/moduly/agents/panel-zakresu-eksperta.ts

Grupy dostępu do rozszerzeń i dostępu do modułów są wartościami wyliczenia
uprawnień i stoją wyżej jako zwykłe wiersze Permissions Center. W tym panelu
mieszkają trzy rzeczy, które wartością logiczną nie są: moduły zastosowania
jako wykaz kodów, gdzie wykaz pusty znaczy brak ograniczenia i jest stanem
wyjściowym; izolacja techniczna jako macierz ośmiu suwaków, identyczna
z macierzą okna konfiguracji punktów izolacji, bo macierz jest jedna;
oraz para ustawień sieci podagentów — czynność i górna liczba podagentów.

Podgląd polityki efektywnej pokazuje wynikowy zestaw ustawień po uwzględnieniu
wszystkich trzech grup. Jest to podgląd, nie bramka — komenda niczego nie
zapisuje i niczego nie rozstrzyga.

Każde zawężenie ustawione w tym panelu jest zawężeniem egzekwowanym: rdzeń
czyta te zapisy przy nakładaniu eksperta na okno i przy powołaniu podagentów,
i na ich podstawie odmawia. Panel mówi o tym Operatorowi wprost, żeby żaden
suwak nie został wzięty za ozdobę.

## budowa/klient-poprzedni/src/moduly/assistant/panel-kontekstow.ts

Kontrakt niesie jeden mechanizm przełączania pamięci w trakcie pracy: komenda
`memory.toggle` ustala, które poziomy zasięgu karta sesji widzi i czy zapis
pamięci jest czynny. To jest kontekst, który rzeczywiście da się przełączyć,
i tak też okno go nazywa.

Nazwanych zestawów pamięci okno nie udaje. Kontrakt nie ma bytu, w którym taki
zestaw miałby zamieszkać: obszar `memory.*` zna wpis oraz poziom zasięgu, a nie
nazwany profil pamięci. Brak nazywa przycisk obok, zamiast pola, którego rdzeń
nie zapisze.

Przestawienie idzie jednym wywołaniem obejmującym oba pola naraz. Rdzeń
odczytuje pole `levels` jako stan po przestawieniu, a pusta lista zostawia
poziomy bez zmian, więc wysyłanie samego jednego przełącznika zamieniałoby
resztę wykazu w brak zmian przy każdym kliknięciu, a Operator odznaczający
ostatni poziom nie miałby jak wyłączyć wszystkich.

Stan pokazany po przestawieniu bierze się z odpowiedzi rdzenia: kontrolki
odbijające samo zamówienie okna kłamałyby przy każdym takim wywołaniu.

## budowa/klient-poprzedni/src/moduly/assistant/skutek-sterowania.ts

Plik odpowiada wyłącznie za przekład odpowiedzi rdzenia na zdanie dla
czytającego. Stoi poza oknem monitora, ponieważ okno składa tabelę, a to jest
ocena tego, co rdzeń oddał.

Reguła jest jedna: porównaj zamówienie ze zleceniem, które wróciło. Odpowiedź
niesie zlecenie po zmianie, a przekład sterowania na stan stoi w rdzeniu
w jednej funkcji, więc sprawdzenie nie kosztuje nic i obowiązuje tak samo dla
priorytetu, jak dla czterech przycisków panelu akcji. Zamówienie i odpowiedź
rozchodzą się z dwóch powodów. Komenda zapisuje priorytet wyłącznie przy
sterowaniu innym niż odczyt, a sama zmiana kolejności idzie torem odczytu
i wraca ze zleceniem bez zmian. Wykonawca zlecenia domyka je natomiast własnym
zapisem, a rdzeń oddaje wiersz odczytany po zapisie, więc zapis wykonawcy
potrafi wejść między zapis sterowania a jego odczyt.

Zdanie odmowy powstaje wyłącznie z rozbieżności między zamówieniem
a odpowiedzią rdzenia. Okno nie orzeka o braku, którego rdzeń nie pokazał, tak
samo jak nie potwierdza skutku, którego rdzeń nie oddał.

Słownik stanów zamówionych jest kluczowany stałymi kontraktu, więc dopisanie
sterowania przerwie kompilację w tym miejscu, zamiast po cichu wpaść w gałąź
braku zamówienia. Zatwierdzenie bramy potwierdzeń puszcza zlecenie w bieg, bo
czekało na rękę operatora, a nie na zasób; odmowa zdejmuje je z kolejki i jest
anulowaniem z powodem, nie osobnym stanem.

## budowa/klient-poprzedni/src/moduly/design/porownanie-wariantow.ts

Porównanie zestawia pola opisowe zasobów, a nie ich obrazy. Pole `uri` zasobu
wskazuje ścieżkę w systemie plików rdzenia, a droga po bajty obrazu, choć opisana
w kontrakcie, nie ma jeszcze uchwytu w rdzeniu, więc przeglądarka nie ma czym
wczytać żadnej ze stron porównania. Suwak przejścia między stroną „przed"
a stroną „po" nabiera znaczenia dopiero po dobudowaniu tej obsługi.

## budowa/klient-poprzedni/src/moduly/browser/ster-wyboru.ts

Trzy nastawy modułu Browser — rodzaj wyodrębnienia w pasku dolnym oraz powiązane
źródło i moduł docelowy w formularzu notatki — korzystają z jednego mechanizmu
rozwijania z komponentu menu-drzewa. Rozwijanie, filtrowanie, haczyk przy pozycji
i wędrówka klawiszami należą do tego mechanizmu, a ster wyłącznie go obsadza
pozycjami i odbiera wybór.

Wartość nastawy przechowuje jedno pole. Napis na uchwycie oraz haczyk przy
pozycji biorą się z tego pola przy każdym przerysowaniu, a odczyt wartości oddaje
to samo pole, więc uchwyt i wykaz nie mogą pokazać dwóch różnych nastaw.

Podpis nad sterem jest blokiem tekstu, a nie etykietą formularza. Etykieta bez
wskazania pola związałaby się z pierwszym potomkiem dającym się etykietować,
czyli z uchwytem menu. Kliknięcie w podpis otwierałoby wtedy wykaz, a nazwa
dostępna uchwytu konkurowałaby z opisem dostępnym, który mechanizm menu składa
sam z nazwy nastawy i wartości bieżącej.

Podpis pojawia się na ekranie tam, gdzie ster sąsiaduje z polami formularza,
ponieważ sama wartość nie mówi, czego dotyczy. Przy pasku dolnym nazwa nastawy
trafia wyłącznie do opisu dostępnego uchwytu.

Ustawienie nastawy z pominięciem Operatora oddaje informację o przyjęciu.
Wartość spoza wykazu nie zostaje na uchwycie: ster wraca do pozycji pierwszej
i oddaje odpowiedź odmowną. Milczące przyjęcie wartości nieobecnej w wykazie
dałoby uchwyt pokazujący nastawę, której nie da się wybrać, więc rozstrzygnięcie
o powiadomieniu Operatora należy do wołającego. Z tego samego powodu wymiana
wykazu utrzymuje wybór tylko wtedy, gdy jego pozycja nadal w wykazie stoi,
a w przeciwnym razie nastawą staje się pierwsza pozycja nowego wykazu.

## budowa/klient-poprzedni/src/aplikacja/aplikacja.ts

Gotowość Centrum dowodzenia jest domknięciem pierwszego odczytu strony. Punkt
wejścia podaje ją scenie wejścia w katalogu `ladowanie/`, żeby ekran między
bramką a produktem schodził w chwili, gdy strona ma czym stanąć. Nieodwiedzona
trasa strony głównej znaczy gotowość natychmiastową: nie ma na co czekać, skoro
widoku nikt nie zbudował.

Motyw rusza przed montażem czegokolwiek, żeby dokument dostał żetony motywu,
zanim pojawi się pierwszy element. Router buduje widok leniwie, przy pierwszym
wejściu na trasę, i już go nie porzuca — karty opuszczonego środowiska trwają
w tle i wracają z pełnym stanem. Środowisko przygotowuje się od razu, ponieważ to
ono podpina przepływ komunikatów.

Uzgodnienia z rdzeniem złożenie aplikacji nie rozpoczyna: robi to przepływ
komunikatów, gdy transport zgłosi stan połączenia, a wywołanie stąd dałoby drugie
powitanie. Tożsamość klienta pochodzi z powitania, ponieważ komendy `session.bind`
i `session.focus` żądają identyfikatora klienta przedstawionego rdzeniowi raz.

Pas kart sesji powstaje w głębi powłoki, która o rdzeniu nie wie, bo buduje ją
także stanowisko podglądu; drogę do rdzenia podaje jej złożenie aplikacji.

## budowa/klient-poprzedni/src/aod/progi-aod.ts

Wszystkie progi funkcji globalnej Always On Display są ustawieniami zasięgu tej
funkcji i podlegają zmianie z okna konfiguracji, z okna Ustawień oraz poleceniem
języka naturalnego. Kontrakt nie niesie ani kategorii ustawień tej funkcji, ani
komendy zapisującej te wartości — dlatego wartości domyślne stoją w pliku jako
stałe, a nie jako odczyt komendy `config.get`. Rozjazd jest zgłoszony, nie
zasypany atrapą.

Chwilę końca wyciszenia do końca dnia liczy zegar maszyny Operatora, ponieważ
funkcja stoi na jego stanowisku, a nie w strefie czasowej rdzenia.

Przekroczenie progu czasu oczekiwania zadania w kolejce oraz przekroczenie progu
wypełnienia kolejki tworzy sugestię klasy „stan kolejki zadań". Przekroczenie
progu powtarzalności czynności ręcznej tworzy sugestię konfiguracji. Sugestia
nieprzyjęta po upływie czasu życia otrzymuje status odrzuconej i znika z listy
oczekujących.

## budowa/klient-poprzedni/src/moduly/design/stan-okna.ts

Komunikat stanu przesłania treść okna, a nie zastępuje jej: nieudane odświeżenie
nie kasuje tego, co Operator już widział, i treść wraca nietknięta, gdy okno
wróci do fazy gotowej. Odmowa rdzenia jest widoczna w oknie, w którym Operator ją
wywołał, a nie wyłącznie w konsoli.

Wygląd bierze się w całości z biblioteki `komponenty/`, ze stanu pustego oraz ze
wskaźnika obrotowego, więc plik nie zna ani jednej barwy. Nazwy faz i znakowanie
powłoki są wspólne dla wszystkich modułów w pliku `komponenty/faza-okna.ts`;
tutaj zostaje wyłącznie to, czym moduł Design różni się świadomie: wskaźnik
odczytu i chowanie treści na czas ładowania.

Kreator Prompt Buildera stoi na natywnym elemencie `dialog` otwartym wywołaniem
`showModal`, a modal otwarty czyni resztę dokumentu bezwładną. Schowanie treści
okna na czas ładowania schowałoby też ten modal — bezwładność by została, a strona
nie miałaby ani jednej drogi wyjścia. Dlatego treści schować nie wolno, gdy w niej
stoi otwarty modal.

## budowa/klient-poprzedni/src/moduly/diagnostics/zrodlo-obserwowalnosci.ts

Rodzina `monitor.*` jest jedynym źródłem obserwowalności, jakie kontrakt niesie
poza obszarem `diagnostics.*`. Metryki wydajności i kontrola stanu samej platformy
wywodzą się z encji `proces_sesji`, a struktura `MonitorStatus` jest dokładnie jej
odwzorowaniem w kontrakcie. Zakładki Metrics & Performance oraz Health & Uptime
jadą więc jednym odczytem: są to dwie perspektywy na ten sam materiał, a nie dwa
niezależne pomiary, które rozejdą się po pierwszej zmianie.

Odmowa komendy `monitor.status` bywa tu zjawiskiem zwykłym: rdzeń bez wpiętego
rejestru telemetrii odmawia głośno w pliku `core/handlers_monitor.go`, a wykaz
pusty ukryłby ten fakt pod zdaniem o braku procesów.

Podział na dwie komendy jest podziałem ról, nie wygodą: `monitor.status` czyta
stan bieżący i niczego nie zapisuje, a `monitor.subscribe` dodatkowo zapisuje okno
na telemetrię. Pole `subscribed` odpowiedzi mówi, czy zapis doszedł do skutku —
żądanie bez identyfikatora okna jest zwykłym odczytem, obsłużonym w pliku
`core/adapter_modul_monitor.go`, a moduł Diagnostics montuje się bez okna.

Zdarzenie `progress.changed` dochodzi do wszystkich połączeń konta niezależnie od
tego, czy komenda `monitor.subscribe` zapisała okno; zapis mówi rdzeniowi
wyłącznie, które okno których procesów pilnuje. Dlatego wykaz zmienia się na żywo
także wtedy, gdy pole zapisu niesie wartość fałszywą.

## budowa/klient-poprzedni/src/mission-control/pas-relacji.ts

Kontrakt nie niesie odczytu powiązań między sesjami. Komplet danych pulpitu podaje
w tym miejscu `null`, a pas mówi o braku źródła. Jest to co innego niż pusty wykaz:
zdanie „brak powiązań" byłoby twierdzeniem, którego nie da się odczytać z rdzenia,
więc pas rozdziela oba stany i nazywa je osobnymi zdaniami.

Znak relacji jest podwójny — strzałka `⇄` albo `→` oraz słowo w podpowiedzi i w treści
czytanej. Kierunek powiązania nie zależy wtedy od samego kształtu znaku, który dla
odczytu ekranowego jest niedostępny.

## budowa/klient-poprzedni/src/moduly/agents/poziomy-pamieci.ts

Konfiguracja pamięci jest czwartym z siedmiu komponentów definicji eksperta.
Poziomy są cztery i wybiera się je niezależnie, ponieważ ekspert korzysta naraz
z kilku: pamięci globalnej Operatora, pamięci projektu, pamięci karty sesji oraz
pamięci środowiska.

Piątego poziomu oznaczającego wyłączenie nie ma i kontrolka go nie dorabia:
wyłączeniem jest zbiór pusty, czyli cztery pola odznaczone. Kontrakt rozstrzyga
to wprost przy wyliczeniu `MemoryLevel` — piąta wartość obok czterech poziomów
pozwalałaby zapisać stan sprzeczny, w którym sesja jest zaznaczona równocześnie
z wyłączeniem.

Dlatego pod grupą stoi zdanie czytające stan bieżący, a nie stała treść:
Operator ma widzieć, że odznaczenie wszystkiego jest wyłączeniem pamięci, zanim
naciśnie zapis. Zdanie zmienia się przy każdym kliknięciu.

## budowa/klient-poprzedni/src/komponenty/faza-okna.ts

Zestaw faz nie zna wartości spoczynku. Zdanie o tym, że rdzeń nie był jeszcze
pytany, opisuje stan źródła danych, a nie stan okna, i ma własny typ
`FazaOdczytu` w `dostepy/stan-dostepow.ts` oraz w stanach pozostałych modułów.
Moduły odwzorowują go na okienne `puste` z osobną treścią komunikatu — tak
robią `moduly/design/okno-assets-panel.ts` oraz
`moduly/library/okno-library-explorer.ts`.

Widoczność wskaźnika odczytu i chowanie treści na czas ładowania zostają
w modułach i różnią się między nimi celowo: Assistant oraz Apps zostawiają treść
widoczną w ładowaniu, Design i Research ją chowają.

Rola `alert` należy się wyłącznie odmowie, ponieważ przerywa czytnikowi ekranu
bieżącą wypowiedź; pozostałe fazy idą jako `status`. Pas komunikatu znika tylko
w fazie `gotowe` — stan przesłania treść, a nie kasuje jej, więc po powrocie do
fazy `gotowe` widać to, co stało na ekranie przed nieudanym odświeżeniem.

## budowa/klient-poprzedni/src/konfiguracja/wybor-adresu.ts

Ta sama kontrolka obsługuje dwie role okna konfiguracji. Pasek u góry ustawia
punkt widzenia, względem którego liczone jest dziedziczenie pól, a panel przy
polu ustawia adres zapisu wartości nadpisującej. Różni je wyłącznie wykaz
dopuszczalnych poziomów i osi: pasek podaje wszystkie, panel pola tylko te,
które dopuszcza katalog ustawienia.

Pola bytu znikają tam, gdzie poziom albo oś bytu nie mają, czyli przy poziomie
globalnym i przy osi platformy.

## budowa/klient-poprzedni/src/moduly/design/karta-zasobu.ts

Karta zasobu jest przyciskiem, a nie prostokątem, ponieważ wybór zasobu
przestawia naraz panel metadanych, kanwę Design Board oraz pole obrazu
referencyjnego Prompt Buildera. Element musi być z tego powodu osiągalny
klawiaturą. Zasób niesie pole `uri` tylko wtedy, gdy rdzeń zna jego położenie;
bez tego pola karta pokazuje pole zastępcze z rodzajem zasobu zamiast miniatury.

Słownik nazw rodzajów pokrywa cały typ kontraktu, więc brak wartości zatrzymuje
kompilację. Wynik pracy modelu wchodzi do wykazu zasobów tą samą drogą, co zasób
wniesiony ręcznie, dlatego słownik obejmuje także rodzaje wytwarzane przez
model — bez nich karta pokazałaby pustkę w miejscu nazwy.

Zawężanie wykazu frazą jest miejscowe i nie jest polem żądania. Komenda
`design.asset.list` zawęża polami identyfikatora okna, rodzaju, etykiet,
znacznika ulubionego oraz ograniczenia liczby wyników, a frazy wśród nich nie
ma. Predykat zawężania stoi w jednej kopii, ponieważ czytają go dwa widoki tego
samego zbioru: okno Assets Panel oraz panel `zasoby-designu` stosu paneli
pomocniczych. Dwie kopie dawałyby na tę samą frazę różne wyniki.

Etykieta zasobu nie niesie stanu, więc wchodzi plakietką bazową. Barwna odmiana
plakietki pozostaje zarezerwowana dla plakietek znaczących stan.

## budowa/klient-poprzedni/src/moduly/automations/stan-automatyki.ts

Pięć okien modułu pracuje nad tą samą automatyką: Workflow Builder buduje jej kroki,
Orchestrator układa między nimi zależności, Scheduler nadaje jej cykliczność, Queue
Manager uruchamia jej kolejkę, a Execution Monitor pokazuje jej przebiegi. Wskazanie
trzymane osobno w każdym oknie rozjeżdżałoby się przy pierwszej zmianie automatyki.

Kolejka bieżąca stoi obok automatyki, ponieważ Queue Manager musi wiedzieć, którą
kolejkę posuwa, a Execution Monitor — po której przyszedł stan przebiegu. Kolejka jest
bytem silnika kolejek, a nie definicji, więc nie należy do struktury `AutomationWorkflow`
i nie da się jej z niej odczytać.

Stan przebiegu przychodzi także z pracy innego okna albo innego urządzenia tego konta.
Gdy dotyczy automatyki bieżącej, definicja trzymana w oknie jest już nieaktualna
i okna odczytują ją ponownie.

## budowa/klient-poprzedni/src/moduly/apps/formularz-komponentu.ts

Identyfikator komponentu składa się z jego nazwy, o czym mówi objaśnienie pola. Kontrakt
wymaga pola `id` w każdym komponencie żądania, a rdzeń nadaje własne oznaczenie dopiero
w odpowiedzi. Identyfikator roboczy jest więc kluczem zależności wewnątrz jednego zapisu,
a nie obietnicą trwałości.

Rodzaj komponentu wybiera się rozwijaniem z biblioteki przez obsadę `wybor-z-menu.ts`,
a nie kontrolką natywną przeglądarki, żeby wykaz rodzajów wyglądał tak samo jak pozostałe
wykazy wyboru w oknie.

## budowa/klient-poprzedni/src/dostepy/indeks.ts

Reszta aplikacji zna z tego katalogu jedną czynność — `otworzOknoDostepow`,
wołaną z kanałem rdzenia i identyfikatorem okna rozmowy. Gospodarz, który ma
własny obszar na ekranie i nie chce okna modalnego, bierze zamiast tego samą
sekcję przez `utworzSekcjeDostepow` i osadza jej element u siebie; czynność
`ustawOkno` wiąże ją wtedy z oknem rozmowy.

Okno jest jedno na klienta i żyje między otwarciami, tak samo jak okno
konfiguracji. Powtórne otwarcie wraca do stanu, na którym się skończyło, i nie
gubi subskrypcji zdarzeń `access.*.changed`. Wywołanie z innym kanałem, po
ponownym połączeniu z rdzeniem, buduje okno od nowa, żeby komendy nie szły przez
transport, którego już nie ma.

Nadanie dostępu żyje w obrębie okna, więc identyfikator okna podaje się przy
otwarciu. Otwarcie bez niego jest poprawne: wykaz punktów i katalog roboczy nie
zależą od okna, a w miejscu nadań sekcja mówi, na co czeka.

## budowa/klient-poprzedni/src/moduly/library/archiwum-pakowanie.ts

Spakowanie archiwum stoi w obszarze archiwum panelu metadanych obok utrwalenia
i paczki migracyjnej, lecz jest od nich oddzielone osobnym zdaniem, ponieważ
pracuje na innym zbiorze. Komendy `library.preservation.run`
i `library.package.export` obejmują zasoby repozytorium biblioteki i oddają wynik
jako zasób biblioteki, natomiast `archive.pack` pakuje katalog albo plik ze
stanowiska operatora i oddaje wynik jako zasób magazynu Designu. Nazwanie jednej
czynności drugą byłoby obietnicą, że spakowana została biblioteka.

Wykaz zasobów nie jedzie tą drogą z tego samego powodu, dla którego nie jedzie
komendami obszaru multimediów: identyfikator pliku biblioteki wraca z magazynu
Designu odmową braku zasobu, a przycisk prowadzący do pewnej odmowy nie pełni
żadnej funkcji.

Postaci archiwum kontrakt nie zamyka wyliczeniem: pole jest napisem, a wartość
pusta oznacza postać domyślną rdzenia. Wykaz trzech dostępnych postaci pochodzi
z opisu komendy w kontrakcie, a nie z domysłu.

Archiwum o zerowej liczbie pozycji jest poprawną odpowiedzią, a nie awarią: tak
wraca spakowany katalog pusty. Zdanie odpowiedzi nazywa to wprost, ponieważ plik
powstał i jego pustka jest informacją należną operatorowi.

## budowa/klient-poprzedni/src/aod/odmowy-aod.ts

Zdanie odmowy nakładki składa się z trzech części. Nazwę nieudanej czynności
i powód podany przez rdzeń składa wspólny komponent odmowy. Trzecia część, czyli
rada dla Operatora, zależy od czynności nakładki i wspólny komponent nie ma jak
jej znać, dlatego wykaz rad stoi w warstwie nakładki.

Część odmów rdzenia w rodzinie komend nakładki jest zachowaniem poprawnym, a nie
usterką. Odpięcie procesu spoza wykazu przypiętych zwraca odmowę „nie znaleziono",
a puste wskazanie procesu przy przypięciu albo odpięciu zwraca odmowę
walidacyjną. Bez zdania trzeciego Operator odczytuje obie odpowiedzi jak awarię.

Nakładka nie stawia bramki przed wywołaniem: nie blokuje pustego pola, nie
wygasza przycisku i nie pyta o potwierdzenie. Rozstrzygnięcie należy do rdzenia,
a nakładka nazywa jego odpowiedź językiem zrozumiałym dla Operatora.

Wykaz rad obejmuje obszary czynności nakładki — przypięcie, odpięcie, rozmowę,
głos, podpowiedzi, stan i kontekst — oraz obszary kolejki decyzji: procesy,
wstrzymanie, konfigurację i przejęcie. Dla kodu odmowy bez osobnego zdania
podawana jest rada domyślna.

## budowa/klient-poprzedni/src/moduly/agents/zasieg-okien.ts

Najwęższym poziomem zasięgu uprawnień jest okno komunikacji i tryb uprawnień
żyje właśnie tam: pole `Window.permissionMode` niesie wartość przełącznika
`--permission-mode` kanału głównego, a zmienia je komenda `window.update`.

Ta część nie powiela ani okna izolacji, ani panelu sterowania: wykaz pokazuje
stan wszystkich okien sesji naraz, czego żadne z nich nie robi.

Lista wyboru pokazuje tryb okna, a nie wybór wskazany przez Operatora.
Przeglądarka przestawia listę natychmiast, a zapis idzie dopiero po niej, więc
po odmowie rdzenia lista stałaby na wartości, której okno nie ma. Dlatego każda
odmowa cofa listę do trybu ostatnio potwierdzonego przez rdzeń.

## budowa/klient-poprzedni/src/konfiguracja/lancuch-zasiegow.ts

Wiersz, z którego pochodzi wartość obowiązująca w punkcie widzenia okna, jest oznaczony
sygnałem. Pozostałe zapisy zostają widoczne, ponieważ Operator ma wiedzieć nie tylko to,
co obowiązuje, lecz także to, co czeka na innym poziomie.

Ostatni wiersz należy do wartości domyślnej katalogu. Brak zapisu nie jest dziurą: jest
wartością domyślną i tak został nazwany.

Każdy zapis wolno usunąć. Usunięcie nie jest niszczeniem ustawienia, tylko zdjęciem
nadpisania, po którym wartość wraca do poziomu szerszego.

## budowa/klient-poprzedni/src/moduly/multitasking/sterowanie-koordynatora.ts

Pas sterowania koordynatora zestawia siedem przycisków, z których pięć — Start, Stop, Pauza, Wznów i Powtórz — dzieli jeden silnik kolejek: idą akcją kolejki na kolejce etapu, tej samej, którą obsługuje pętla sesyjna i moduł Automations. Koordynator nie ma własnego wykonawcy zleceń i nie prowadzi biegu naprawczego: bieg prowadzi rdzeń, a licznik obiegów przychodzi z odczytu stanu okna.

Przycisk „Przekaż” wydaje zlecenie wykonawcy. Adresata i treść rozstrzyga tryb współpracy; nośnikiem jest wysłanie wiadomości, ponieważ okno wykonawcy jest oknem komunikacji. Droga modelu, czyli wywołanie narzędzia platformy, biegnie obok i tej kontrolki nie zastępuje.

Przycisk „Waliduj” sprawdza przebieg, a nie ocenia wynik: kontrakt nie ma komendy oceny, więc przycisk odczytuje stan koordynatora i wykonawców i mówi, czy tura trwa, ile obiegów zliczono i czy bieg stoi wraz z powodem. Ocena treści wyniku należy do modułu Results Analyzer.

Zatrzymanie biegu gasi kolejkę etapu i tury wykonawców naraz. Kolejność ma znaczenie: najpierw gaśnie tura, która właśnie zużywa czas modelu, a dopiero potem kolejka, która mogłaby ją wznowić.

Przekazanie zlecenia to dwie czynności i obie muszą się odbyć. Doręczenie treści oknu wykonawcy pozwala mu ruszyć do pracy. Utrwalenie w bazie rdzenia zapisuje, kto komu co zlecił, i zakłada pozycję kolejki; bez niego więź ginie z zamknięciem przeglądarki, a moduł Mission Control nie ma czego pokazać po ponownym uruchomieniu. Kolejność ma znaczenie: najpierw zapis, potem doręczenie — odwrotna zostawiałaby wykonawcę pracującego nad zleceniem, którego nikt nie odnotował. Nieudany zapis nie wstrzymuje doręczenia i nie liczy się jako niedoręczenie: dwa rodzaje niepowodzenia wracają osobnymi wykazami, bo pierwszy mówi, że wykonawca nie dostał pracy, a drugi, że dostał, ale po ponownym uruchomieniu nikt tego nie odtworzy.

## budowa/klient-poprzedni/src/moduly/apps/karta-rozszerzenia.ts

Pola, których pozycja katalogu nie niesie — dziennik zmian, wykaz narzędzi oraz
wymagane uprawnienia — nie mają na karcie zaślepki. Stoją w wykazie braków okna,
żeby nieobecność danych była widoczna, a nie udawana.

Stan pozycji nigdy nie opiera się na samej barwie: plakietka niesie słowo, a nie
tylko odmianę stylu. Cztery rozróżnialne stany (dostępna, zainstalowana
wyłączona, zainstalowana włączona, dostępna aktualizacja) składa się z dwóch pól
kontraktu, więc czwartego z nich karta nie podaje — rejestr nie niesie wersji
dostępnej do zestawienia z wersją zainstalowaną.

Instalacja i przełączenie mają w rdzeniu osobne komendy, dlatego karta stawia
dwa osobne przyciski. Jeden przycisk o zmiennym napisie kazałby zgadywać, którą
czynność zaraz zleca. Przycisk uprawnień otwiera Permissions & Trust Center na
pozycji karty i pozostaje czynny również dla pozycji niezainstalowanej, ponieważ
zakres dostępu ogląda się przed włączeniem.

## budowa/klient-poprzedni/src/aplikacja/przelacznik-tras.ts

Żaden przycisk trasy nie jest bramą i żaden nie zostaje wyszarzony. Do
środowiska można wejść także wprost, ponieważ środowisko domyślne istnieje od
pierwszej chwili. Trasa bieżąca jest oznaczona atrybutem `aria-current`, a nie
odebraniem klikalności.

Przyciski niosą sam znak, bez napisu, ponieważ trasy stoją w grupie akcji paska
obok ustawień, motywu i menu konta, które również są ikonami, a nazwa widoku
stoi tuż pod paskiem jako tytuł strony. Nazwę niesie atrybut `title` oraz
etykieta dostępności przycisku.

## budowa/klient-poprzedni/src/moduly/design/historia-promptow.ts

Odczyt historii i zapis szablonu mają już komendy w kontrakcie, brakuje im
natomiast uchwytów w rdzeniu. Okno tej drogi jeszcze nie wywołuje, więc
kontrolka zapisu szablonu nazywa stan, zamiast wykonywać zapis.

Porównanie wymienia pola kontraktu, którymi dwa prompty się różnią. Prompt jest
zbiorem pól, więc różnica liczona na samym tekście nie niosłaby informacji
o tym, co się zmieniło.

## budowa/klient-poprzedni/src/mission-control/przyciski-transportu.ts

Żaden przycisk paska nie dostaje `disabled`, `aria-disabled` ani klasy wygaszającej:
blokada nie jest dozwolonym sposobem informowania o stanie. Stan kolejki jest wypisany
słowem przy jej nazwie, a przycisk zatrzymania pozostaje czynny w każdym stanie.

Osiem działań wyliczenia `QueueAction` leżących poniżej sterowania biegiem — wstawienie
i zdjęcie z kolejki, odłożenie w czasie, rozdzielenie, scalenie, skierowanie,
rozgałęzienie i warunek — należy do układania przebiegu, a nie do jego prowadzenia.
Pasek transportu pulpitu ich nie pokazuje, ale mapa ikon musi je znać, ponieważ
wyliczenie kontraktu je niesie, a mapa niepełna nie skompilowałaby się przy pierwszym
ich użyciu.

## budowa/klient-poprzedni/src/dostepy/wybor-korzeni.ts

Reguła pustego zbioru stoi przy kontrolce zdaniem, a nie w podpowiedzi, ponieważ
Operator, który odznaczy wszystko, ma wiedzieć, że właśnie nadał wszystko.

Nadanie nie sięga poza korzenie punktu, więc kontrolka nie przyjmuje ścieżki wpisanej
z ręki: pokazuje korzenie punktu i pozwala je zaznaczyć. Poszerzenie obszaru wymaga
zmiany punktu, nie nadania.

## budowa/klient-poprzedni/src/moduly/library/straz-odmow.ts

Rdzeń odpowiada na komendę pozbawioną uchwytu kopertą zdarzenia `<obszar>.unknown`, która niesie identyfikator żądania, ale nie ma pola `status` (`server/internal/protocol/zadanie.go`). Korelacja po stronie klienta rozstrzyga wyłącznie koperty ze statusem (`protokol/koperta.ts`), więc obietnica zwykłego wywołania po takiej odmowie nigdy się nie rozstrzyga, a okno zostaje w stanie ładowania. Straż wiąże odmowę z żądaniem po identyfikatorze i zamienia ją w zwykły wynik z polem błędu, dzięki czemu okno obsługuje odmowę tą samą drogą co każde inne niepowodzenie i mówi wprost, której komendy rdzeń nie zna.

Sześć obszarów własnych odpowiada temu, co moduł woła: własny obszar `library`, okno komunikacji (`window.action`, `window.state.get`), przenoszenie kontekstu (`context`), katalog akcji (`action`), katalog modułów (`module.list`, który obsadza ster modułu docelowego) oraz komplet kontekstu okna (`aod.context.get`, jedyna komenda kontraktu czytająca to, co przyniosło przekazanie).

Siódma subskrypcja idzie na obszar zapasowy `connection`, ponieważ rodzina `knowledge.*` — wyszukiwanie po znaczeniu oraz wskaźnik znaczenia biblioteki — nie ma własnego zdarzenia odmowy. Wykaz `zdarzeniaNieznanej` (`shared/contract.go`) nie zna klucza `knowledge`, więc rdzeń bez wpiętego portu Wiedzy odpowie kopertą `connection.unknown`. Bez tej subskrypcji obietnica takiego wywołania nigdy by się nie rozstrzygnęła, a okno zostałoby w ładowaniu. Wiązanie idzie po identyfikatorze żądania, więc cudza odmowa obszaru zapasowego niczego tu nie rozstrzyga.

Subskrypcje są wypisane po jednej, a nie złożone pętlą, ponieważ kształt treści zdarzenia bierze się z jego nazwy: pętla po wykazie zgubiłaby typ ładunku i kazałaby go rzutować na ślepo.

## budowa/klient-poprzedni/src/moduly/browser/zebrane-w-sesji.ts

Zbiór jest odbiciem rdzenia, a nie drugą prawdą. Wykaz przychodzi komendami
`browser.source.list` oraz `browser.note.list` i podmienia zawartość w całości, a nie doszywa
się do zastanej: rdzeń wie, co w oknie jest, i to jego odpowiedź rozstrzyga. Dopisanie po
udanym dodaniu zostaje obok — nowa pozycja ma być widoczna od razu, bez czekania na ponowny
odczyt, ale jest tym samym wierszem, który przyjdzie w wykazie.

Przypięcie i klasyfikacja są własnością widoku, nie rdzenia. Kontrakt nie niesie ani pola
przypięcia, ani pola rodzaju notatki, więc jedno i drugie zapamiętane w zbiorze nie udaje
zapisu, a panel notatek mówi o tym w swoim opisie. Dlatego podmiana wykazu notatek zostawia
przypięcia nietknięte, a przypięcie pozycji, której rdzeń już nie oddaje, po prostu niczego
nie porządkuje.

Fragment pusty oraz fragment taki sam jak zastany nie wchodzą do wyodrębnionych, a czynność
dopisania mówi o tym wprost, żeby okno nie meldowało dopisania, którego nie było.

Każda zmiana zbioru jest ogłaszana. Źródło dodane w panelu źródeł ma się pojawić w wyborze
powiązania notatki w tej samej chwili; bez ogłoszenia drugie okno zobaczyłoby je dopiero przy
własnym odświeżeniu, a Operator dostałby wybór bez pozycji, którą właśnie zapisał.

## budowa/klient-poprzedni/src/moduly/design/suwaki-promptu.ts

Parametry generowania podaje suwak, a nie pole liczbowe, ponieważ każdy z trzech
parametrów ma w kontrakcie zakres zamknięty albo naturalny: pole `creativity`
przyjmuje wartości od zera do jedności, pole `variants` jest liczbą wariantów
jednego zlecenia, a pole `seed` liczbą całkowitą powtarzającą wynik. Odczyt
wartości stoi obok suwaka, ponieważ ziarno trzeba umieć przepisać, żeby powtórzyć
wynik generowania.

Wypełnienie toru suwaka niesie żeton `--dn-suwak-pozycja`, czytany przez arkusz
`komponenty/suwak.css`. Żeton jest ustawiany przy każdej zmianie wartości
kontrolki. Arkusz ma dla niego wartość zapasową równą połowie zakresu, więc bez
tego ustawienia tor każdego suwaka pokazywałby połowę niezależnie od wartości
kontrolki.

## budowa/klient-poprzedni/src/moduly/browser/odczyt-migawki.ts

Treść migawki pochodzi z tabeli `migawka_strony`, którą wypełnia `browser.navigate`,
więc przeżywa przeładowanie powłoki i wymianę karty. Odczyt jest jeden i mieści się
w tym pliku, wzorem `wykazy-zebranego`, żeby okna pytały o wynik, a nie o sposób.

Odmowa `not_found` nie jest błędem okna: rdzeń odpowiada nią, dopóki w oknie nie odbyło
się ani jedno przejście, co jest normalnym stanem okna świeżo otwartego. Zostaje wtedy
stan `nietknieta` — podgląd pusty i zdanie o powodzie, bez stanu błędu. Każda inna odmowa
mówi swoim powodem.

Odczyt w toku jest osobnym stanem, a nie odmianą pustki. Zdania „nie ma czego pokazać"
i „pytam" to dwie różne rady: pierwsza każe działać, druga czekać.

Migawki nie kasujemy przy odmowie, ponieważ nieudany odczyt nie unieważnia strony, którą
Operator już czyta. Migawka świeżo wchłonięta unieważnia natomiast zdanie o poprzednim
odczycie: powód „okno nie ma jeszcze migawki" przestaje być wtedy prawdziwy.

## budowa/klient-poprzedni/src/moduly/agents/zrodlo-zaplecza.ts

Komenda `channel.list` podaje rejestr kanałów modelu dla okna Model
Configuration; moduł Agents kanałów nie zakłada, lecz wybiera spośród wierszy
rejestru rdzenia. Komenda `config.capabilities.get` podaje deklarację zdolności
adaptera dostawcy dla pól konfiguracji sesji, dzięki czemu okno pokazuje wprost,
którego parametru wybrany kanał nie obsłuży, zanim Operator go wypełni. Komenda
`access.point.list` podaje katalog mostów MCP dla okna Connectors Manager — ten
sam katalog, z którego rdzeń składa wykaz serwerów procesu modelu. Komenda
`identity.category.list` podaje słownik kategorii tożsamości dla Agent Buildera.
Para komend `window.list` oraz `window.update` obsługuje tryb uprawnień
przypisany jednemu oknu komunikacji w oknie Permissions Center.

Zapytanie o zdolności obejmuje obszar modelu, ponieważ okno Model Configuration
wypełnia wyłącznie pola modelu prowadzącego i parametrów jego wywołania. Pytanie
o komplet obszarów przyniosłoby deklarację, której okno nie pokaże.

## budowa/klient-poprzedni/src/moduly/automations/dziennik-przebiegow.ts

Dziennik nie jest pełnym logiem przebiegu. Zbiera wyłącznie to, co przyszło do
okna otwartego, więc przebieg sprzed otwarcia okna nie ma tu ani jednego wiersza.
Pełny zapis oddaje osobna komenda logu przebiegu — okno nazywa ją pozycją paska
akcji i mówi wprost, czy rdzeń ma dla niej uchwyt.

Poziomu wiersza samo zdarzenie nie niesie. Zamiast zgadywać go z treści, dziennik
zapamiętuje stan przebiegu z chwili wpisu; po nim idzie zawężanie i po nim widać,
czy wiersz powstał w przebiegu, który jeszcze trwał, czy w takim, który już się
załamał.

## budowa/klient-poprzedni/src/moduly/design/pasek-postepu.ts

Komenda `design.asset.generate` postępu nie zgłasza, ponieważ rdzeń nie ma czym
wytworzyć obrazu, więc okno kreatora wycisza pasek zamiast zostawiać go w stanie
oczekiwania.

Pasek bez oczekiwanego procesu milczy: zdarzenie `progress.changed` jedzie
z każdego strumienia rozmowy w rdzeniu, obsługiwanego w pliku
`core/telemetria_strumienia.go`, a nie tylko z pracy tego okna. Przyjmowanie
każdego zdarzenia pokazywałoby tu postęp pracy, której to okno nie zlecało.
Rozstrzygnięcie własnego postępu wygląda tak: ze wskazanym procesem liczy się
zgodność procesu, a dopóki rdzeń procesu nie nazwał, wystarczy zgodność okna.

Miniatura w budowie nie jest podglądem wyniku, bo wyniku jeszcze nie ma; jest
widocznym śladem tego, że rdzeń pracuje nad zasobem tego okna. Wygląd bierze się
z biblioteki komponentów, więc plik nie zna ani jednej barwy.

## budowa/klient-poprzedni/src/moduly/agents/edytor-eksperta.ts

Edytor składa formularz tożsamości — nazwę, imię własne, favikon, opis
i instrukcje — z edytorem warstw promptu oraz z wykazem kategorii tożsamości
pochodzącym z katalogu rdzenia. Każda z tych rzeczy ma własne wywołania i własne
odmowy, więc mieszka we własnym pliku; tutaj zostaje samo złożenie i przekazanie
eksperta czynnego w dół.

Katalog kategorii jest tylko do odczytu: kategorie tożsamości pochodzące
z komendy `identity.category.list` są własnością rdzenia i edytor nie dopisuje ani
jednej własnej. Okno pokazuje je jako kontekst, w który wchodzą warstwy eksperta,
a nie po to, żeby je stąd zmieniać. Sam edytor nie wykonuje ani jednego wywołania
rdzenia.

## budowa/klient-poprzedni/src/moduly/library/formularz-kolekcji.ts

Utworzenie kolekcji oraz przypisanie do niej zasobu to dwie z trzech funkcji
Operatora w oknie kolekcji i znaczników. Kolekcja jest swobodna, a nie regułowa:
komenda `library.collection.create` przyjmuje nazwę i opis, a reguły składającej
kolekcję samoczynnie kontrakt nie niesie. Wariant regułowy stoi więc w panelu
akcji jako czynność nazwana wprost, bez drogi wykonania.

Przypisanie działa na zaznaczeniu wykazu plików modułu: okno nie ma własnego
wykazu, ponieważ wykaz jest jeden na cały moduł.

Zdanie przy nazwie pustej mówi o oknie, a nie o rdzeniu. Rdzeń odmawia wyłącznie
nazwy pustej; nazwę złożoną z samego znacznika kolejności bajtów, z samej spacji
nierozdzielającej albo z samych spacji zwykłych przyjmuje i oddaje co do znaku.
Każdą z nich przycięcie po stronie przeglądarki zamienia w pustkę, więc okno jej
nie wyśle i może powiedzieć tylko tyle, że samo nic nie wysłało.

Nazwa pokazana po założeniu pochodzi z odpowiedzi, a rozbieżność wobec wpisanej
jest wypowiedziana: rdzeń zapisuje nazwę dosłownie, więc różnica może wziąć się
wyłącznie z przycięcia po stronie okna, a Operator ma o niej wiedzieć, zanim
zacznie tej kolekcji szukać po tym, co napisał.

## budowa/klient-poprzedni/src/modele/wybor-osi.ts

Sekcja modeli posługuje się dwiema osiami z trzech, lecz oś platformy musi być
w pasku obecna: tożsamość zapisana dla platformy stanowi tło, na którym leży
zapis modelu oraz konta, a pasek jest jedynym miejscem jej podglądu i zmiany.

Byt osi pozostaje polem otwartym. Podpowiedzi pochodzą z rejestru kanałów
i z rejestru kont, lecz identyfikator modelu dotychczas nieużywanego nie jest
błędem, ponieważ wykaz podpowiedzi nie jest zbiorem zamkniętym. Oś platformy bytu
nie posiada, więc pole bytu znika w całości, zamiast stać wygaszone.

Słownik nazw osi pochodzi z warstwy konfiguracji, ponieważ oś rozstrzygania jest
jedna dla całego systemu i nie może być nazywana odmiennie w poszczególnych
oknach.

## budowa/klient-poprzedni/src/modele/wykaz-kont.ts

Wiersz wykazu podaje cztery rzeczy naraz: nazwę konta, jego rodzaj, dostawcę oraz
stan wyrażony plakietkami. Plakietka poświadczenia jest jedyną informacją
o poświadczeniu, jaką klient ma prawo pokazać, ponieważ kontrakt nie zwraca jego
treści żadną komendą. Wykaz podaje więc wyłącznie to, czy poświadczenie zostało
zapisane.

Wiersz konta nieczynnego nie jest wygaszony ani pozbawiony reakcji na kliknięcie.
Nieczynność jest stanem danych, a nie blokadą interfejsu: w takie konto trzeba
móc wejść, aby je z powrotem uruchomić.

Pusty rejestr nie daje pustego prostokąta. Zdanie opisu dobiera się do fazy
odczytu, ponieważ ten sam pusty wykaz znaczy co innego w trakcie zapytania do
rdzenia, co innego po odmowie rdzenia, a co innego wtedy, gdy rdzeń odpowiedział
i rejestr rzeczywiście nie zawiera konta spełniającego warunek wykazu. Jedno
zdanie na wszystkie trzy przypadki nie rozstrzygałoby, czy czekać, czy działać.
Widok pozostaje przy tym czynny.

## budowa/klient-poprzedni/src/moduly/design/stan-kompozycji.ts

Kompozycja powstaje w kliencie, zanim pojedzie do rdzenia: komenda `design.board.update`
przyjmuje gotowy układ w polu `layers`, więc wywołanie idzie przy zapisie, a nie przy każdym
ruchu warstwy.

Wczytanie kompozycji jest drogą powrotną tej samej komendy — dzięki niej odświeżenie okna
pokazuje kompozycję zapisaną po stronie rdzenia, a nie pustą kanwę.

Położenie i rozmiar warstwy ustawia się jedną zmianą, ponieważ cztery liczby opisują jeden
prostokąt. Cztery osobne wywołania rozgłosiłyby cztery stany pośrednie, w których warstwa ma
już nowe położenie przy jeszcze starym rozmiarze.

## budowa/klient-poprzedni/src/konfiguracja/pasek-punktu-widzenia.ts

Wartość ustawienia nie jest jedna — zależy od tego, dla kogo zadane jest
pytanie. Ten sam klucz może mieć inną wartość globalnie, inną w oknie
komunikacji i jeszcze inną dla wskazanego modelu.

Pasek jest wyłącznie soczewką odczytu, a nie poleceniem. Zmiana punktu widzenia
niczego nie zapisuje; przestawia jedynie miejsce, względem którego liczone jest
dziedziczenie, i każe polom formularza przeliczyć swoje pochodzenie od nowa.

## budowa/klient-poprzedni/src/moduly/design/filtr-zasobow.ts

Filtr ma jedną odpowiedzialność: zebrać warunki zawężenia i oddać je oknu.
Cztery pola filtra idą do rdzenia, piąte zostaje po stronie klienta. Rodzaj
zasobu, etykiety, znacznik ulubionego oraz górna granica liczby zasobów są
polami żądania `design.asset.list` i zawężają odczyt po stronie rdzenia. Fraza
wyszukiwania nie jest polem tego żądania, więc zawęża wyłącznie wynik już
otrzymany; pole mówi o tym wprost, zamiast pozorować wyszukiwanie po stronie
rdzenia.

## budowa/klient-poprzedni/src/moduly/design/warstwy-designu.ts

Warstwa pierwsza jest rozwinięta i pozostaje taka: kanwa, wykaz zasobów i kreator stoją
od wejścia do modułu. Warstwy druga, trzecia i czwarta stoją zwinięte, a zapowiedź nad
nimi mówi, co jest pod spodem — zwinięte nie znaczy ukryte.

Nośnikiem rozwinięcia jest znacznik `details`. Postać trzyma przeglądarka, więc element
działa klawiaturą i ma poprawną semantykę bez ani jednego nasłuchu. Druga kopia stanu
w nazwie klasy arkusza stylów mogłaby się z atrybutem `open` wyłącznie rozminąć.

Znacznik wywołania stoi przy nazwie, żeby droga do elementu była widoczna, zanim się go
otworzy. Sam znacznik jest ozdobą uchwytu, a nie jego nazwą, więc idzie z atrybutem
`aria-hidden`: czytnik ekranu odczyta nazwę elementu, a nie znak graficzny.

## budowa/klient-poprzedni/src/mission-control/zrodlo-pulpitu.ts

Źródło ma jedną odpowiedzialność: zebranie stanu z rdzenia i podawanie odbiorcy świeżego kompletu po każdej zmianie. Odczyty i zdarzenia idą wyłącznie kanałem kontraktu, a przełożenie stanu na komplet należy do `zlozenie-danych.ts`.

Odczyty obejmują `session.list` wraz z obecnością, `window.list`, `channel.list` oraz `environment.list`, z którego biorą się nazwy kolumn matrycy — z rdzenia, a nie z kopii katalogu. Transport kolejkuje ramki do chwili połączenia, więc odczyt wysłany przed otwarciem gniazda dochodzi po nim i własny nasłuch stanu łącza jest zbędny.

Subskrypcje obejmują `session.changed`, `window.changed`, `queue.changed` oraz `progress.changed`. Kolejki nie mają odczytu w kontrakcie, ponieważ brakuje komendy `queue.list`, więc ich stan buduje się wyłącznie ze zdarzeń. Do pierwszego zdarzenia pulpit pokazuje uczciwy stan pusty.

## budowa/klient-poprzedni/src/moduly/developer/akcje-kanoniczne.ts

Ta sama operacja nazwana w dwóch miejscach inaczej rozjeżdża przycisk z komendą,
którą wywołuje, dlatego nazwy kanoniczne mają jedno źródło. Dziewiąta operacja,
konwersja języka, należy do paska narzędzi promptu, a nie do ośmiu operacji
panelu akcji, i tablica zachowuje tę różnicę kolejnością pozycji.

Zmiana nazwy symbolu w całym repozytorium jest czynnością rozstrzygalną i wykonuje
ją serwer języka, który zna graf odwołań; model dałby wynik prawdopodobny zamiast
poprawnego, i to w czynności, której poprawność da się sprawdzić. Rdzeń kieruje
więc operację `renameSymbol` do serwera języka, a panel podaje ją tą samą komendą,
ponieważ na pasku jest to ta sama pozycja.

Próg paska nie buduje menu rozwijanego. Jedyne menu biblioteczne
(`komponenty/menu-drzewo.ts`) jest sterem nastawy: na jego uchwycie stoi wartość
bieżąca, a wybór liścia ją zmienia. Operacje jednorazowe wartości bieżącej nie
mają, więc uchwyt pokazywałby napis, który niczego nie odzwierciedla — wszystkie
dziewięć przycisków stoi wobec tego płasko.

## budowa/klient-poprzedni/src/moduly/design/przyciecie-pol.ts

Samo `trim()` przeglądarki tu nie wystarcza, ponieważ `trim()` JavaScriptu
i `strings.TrimSpace` języka Go nie zdejmują tego samego. Produkcja WhiteSpace
JavaScriptu nie zna znaku NEL (U+0085), a `unicode.IsSpace` w Go go zna;
odwrotnie, `trim()` zdejmuje ZWNBSP (U+FEFF), którego Go nie rusza. Rdzeń
etykiet nie przycina w ogóle — `dane/design_zasoby.go` odrzuca wyłącznie
etykietę pustą — więc etykieta złożona z samego NEL zapisałaby się jako
niewidoczna, a etykieta z NEL na końcu nie dałaby się odnaleźć filtrem wpisanym
bez tego znaku.

Brzegową pustką jest tutaj wszystko, co za biały znak uważa którakolwiek ze
stron. Przycięcie idzie sumą obu zbiorów, więc do rdzenia jedzie to, co człowiek
uznałby za wpisane — tak samo po stronie zapisu, czyli nadania etykiet, jak po
stronie odczytu, czyli filtru wykazu. Środek wartości pozostaje nietknięty.

Powtórzenia odpadają, ponieważ rdzeń i tak nie zapisze etykiety dwa razy: klucz
główny tabeli `etykieta_zasobu_design` obejmuje etykietę, a odpowiedź niosłaby
wtedy co innego niż żądanie.

## budowa/klient-poprzedni/src/moduly/design/tryby-kanwy.ts

Moduł zna siedem trybów narzędzia kanwy: zaznaczanie, pióro, kształt, tekst,
pędzel z maską, retusz oraz ramkę interfejsu. Wykaz podaje wszystkie siedem,
ponieważ ukrycie tych, których ta budowa nie wykonuje, przedstawiłoby moduł
jako mniejszy, niż go zaprojektowano.

Żaden tryb nie jest wygaszony: wybór zawsze się udaje, a pod wykazem staje
zdanie mówiące, co tryb w tej budowie robi albo czego mu brakuje. Tryb
zablokowany niczego nie tłumaczy, tryb wybrany wraz ze zdaniem tłumaczy
wszystko.

Dwa tryby zmieniają zachowanie kanwy. Zaznaczanie jest stanem wyjściowym
i prowadzi wskazywanie oraz przeciąganie warstw. Ramka interfejsu kieruje pracę
do presetów ramek przybornika, ponieważ ramka jest tu warstwą o zadanych
wymiarach, a nie osobnym bytem.

Pięć pozostałych trybów wymaga bytów, których kompozycja nie zna: ścieżki
z węzłami, kształtu innego niż prostokąt, tekstu, maski przezroczystości oraz
pędzla działającego na pikselach zasobu. Kompozycja niesie warstwę o położeniu,
rozmiarze, kolejności, blokadzie i adnotacji, i tyle jedzie do rdzenia.

## budowa/klient-poprzedni/src/asystent-plywajacy/zaczep-asystenta.ts

Warstwa osadza się w `aplikacja/widok-srodowiska.ts`, ponieważ tylko tam znana jest
jednocześnie powłoka i droga do rdzenia. Favikon ląduje w korzeniu powłoki, obok obszaru
roboczego: moduł podmienia zawartość obszaru w `aplikacja/przestrzen-modulu.ts`, więc
favikon leżący poza nim zostaje widoczny przy każdej zmianie modułu.

Dwa inne zaczepy się nie nadają. Plik `aplikacja/scena-sesji.ts` trzyma okna równoległe
liczone do sufitu `LICZBA_MAX` z `okna-rownolegle/identyfikatory.ts`, a asystent nie jest
oknem sceny — profil daje mu `granicaOkien: 1` i postać `dymek-glosowy` — więc zabierałby
gniazdo oknu komunikacji i znikał razem ze sceną. Plik `powloka/powloka.ts` buduje także
stanowisko podglądu `powloka/podglad.ts`, które o rdzeniu nie wie; zaczep tam wymusiłby
albo wersję niemą, albo przeciek drogi do rdzenia w dół.

Dymek jest jeden na widok środowiska. Liczby `granicaOkien` nie przepisano do tego pliku:
`okno-dymkowe.ts` bierze ją z rejestru profilów.

Profil daje `pamiecSesyjna: true`, więc zwinięcie dymka nie czyści historii, a `schowaj()`
wyłącznie chowa element. Rozmowa ginąca po zwinięciu wyglądałaby na awarię.

Favikon i stan spina jedna subskrypcja: `stan.obserwuj` przenosi na favikon pracę asystenta
i licznik nieprzeczytanych posunięć. Otwarcie dymka zeruje licznik, zamknięcie znów go
zbiera — dopóki dymek stoi odsłonięty, posunięcia są czytane na bieżąco.

## budowa/klient-poprzedni/src/moduly/diagnostics/prowenancja.test.ts

Sprawdzian ujawnia pracę modeli w oknie Provenance Explorer i pilnuje czterech rzeczy, o które ta zakładka istnieje: że wykaz wywołań ma drogę z okna, że odczyt śladu przynosi drzewo odcinków i treść, że ocena jest czynnością z własnym chwytem, oraz że wydanie śladu mówi, co plik poniesie. Piąta rzecz jest równie ważna: pustka rejestru ma zdanie, a nie migające puste miejsce.

## budowa/klient-poprzedni/src/moduly/apps/os-etapow.ts

Oś prowadzi wyłącznie do okien, które moduł zbudował. Pozycje przychodzą z pliku
składającego moduł, a nie z wykazu zapisanego w samej osi, dzięki czemu
dopisanie okna nie wymaga poprawki w drugim miejscu.

Kolejność etapów — architektura, warsztaty, wdrożenie — jest kolejnością
procesu, a nie bramą. Naciśnięcie etapu późniejszego jest dozwolone i prowadzi
ognisko do jego okna.

Oś nie udaje stanu etapu. Stan przychodzi zdarzeniem `apps.build.changed`
i rysuje go osobny wykaz okna.

## budowa/klient-poprzedni/src/moduly/library/wiersz-pliku.ts

O braku treści wiersz mówi wyłącznie wtedy, gdy rdzeń sam ją orzekł, czyli przy werdykcie
`brak` z pliku `dostepnosc-tresci.ts`, i powtarza jego powód. Plik, o który nikt jeszcze
nie pytał, nie dostaje żadnego znaku.

Odmowa podglądu nie jest orzeczeniem o treści. Rdzeń odmawia kodem `not_found`, gdy pliku
nie ma, oraz kodem `internal_error`, gdy nośnik nie oddał treści spod odwołania. Oba są
zdaniem o pliku i o nośniku, a nie o zawartości repozytorium, więc odmowa dostaje własne
zdanie: powód rdzenia bez dopisanego zarzutu.

Werdykt `odmowa` powtarza sam powód i mówi, że treść pozostaje nieznana. Werdykt
`odwolanie` mówi, że rdzeń wskazał miejsce treści, a jej samej nie podał. Werdykty
`nieznana` oraz `osiagalna` nie dają zdania: pierwsza nie ma czego powiedzieć, druga
niczego nie zarzuca.

Metryka nie jest świadkiem treści: brak pola `versionId` czy `checksum` niczego o niej nie
orzeka i w drugą stronę tak samo. Świadkiem jest odpowiedź rdzenia.

Fragment, na którym oparło się dopasowanie po znaczeniu, stoi przy wierszu, ponieważ bez
niego trafność jest liczbą bez podstawy do sprawdzenia.

## budowa/klient-poprzedni/src/moduly/library/kafel-pliku.ts

Kafel niesie to samo, co wiersz wykazu: zaznaczenie do czynności zbiorczych oraz
wskazanie pliku czynnego dla trzech pozostałych okien modułu. Różni się układem,
przeznaczonym dla materiału oglądanego, a nie czytanego.

Kafel nie pokazuje miniatury i pokazać jej nie może. Komenda podglądu pliku
oddaje dla obrazu wyłącznie odwołanie do niego, a klient nie ma komendy, którą
pobrałby bajty spod tego odwołania. Miejsce miniatury zajmuje więc znak rodziny
treści złożony z typu MIME: podaje, co to za plik, i nie udaje, że okno widziało
jego zawartość.

Zdanie o dopasowaniu po znaczeniu podaje fragment tekstu przysłany przez rdzeń
jako podstawa trafienia. Pokazany fragment pozwala Operatorowi sprawdzić
dopasowanie zamiast przyjmować je na wiarę, a długi fragment zostaje skrócony.
Trafność wchodzi do zdania tylko wtedy, gdy rdzeń ją podał, ponieważ wartość
dopisana przez okno wyglądałaby na pomiar.

## budowa/klient-poprzedni/src/moduly/assistant/zrodlo-narzedzi.ts

Wywołania dzielą się na trzy rodziny, a każda odpowiada za inne pytanie okna.
Komenda `tools.catalog.list` podaje, co w ogóle da się wywołać po ukośniku:
narzędzia, umiejętności i komendy akcji w jednym wykazie. Rodzina
`session.tool.*` podaje, co z tego jest dołożone do bieżącej karty sesji;
dołożenie żyje w stanie sesji i nie rusza definicji eksperta. Rodzina
`automation.*` prowadzi makro i rutynę asystenta do modułu Automations — rdzeń
nie ma innego magazynu sekwencji kroków, więc jest to jedyne miejsce, w którym
makro przeżywa sesję.

Komenda `schedule.get` stoi po stronie odczytu do pary z
`automation.schedule.set`. Bez niej okno pokazywałoby harmonogram, który samo
wysłało, zamiast tego, który rdzeń trzyma.

## budowa/klient-poprzedni/src/mission-control/kolumna-kolejek.ts

Kolumna ma jedną odpowiedzialność: pokazać kolejki znane ze zdarzeń
`queue.changed`, każdą z własnym paskiem transportu opartym o komendę
`queue.action`. Kontrakt nie niesie roli kolejki, liczby zadań czekających ani
odczytu zbiorczego `queue.list`, więc wykaz kolejek buduje się wyłącznie ze
zdarzeń, a miary bez źródła kolumna wypisuje wprost jako brak źródła.

Bieg naprawczy nie ma limitu, dlatego licznik obiegów stoi przy każdej kolejce:
przejrzystość zastępuje w tym miejscu bramę zatrzymującą bieg. Sterowanie jest
ręczne i żaden przycisk transportu nie jest wyszarzony; stan kolejki poznaje się
po plakietce ze słowem, a nie po dostępności przycisku.

## budowa/klient-poprzedni/src/moduly/katalog-okien.ts

Liczba wpisana w pasek uczciwości na stałe rozjeżdża się z rdzeniem po
cichu: katalog okien zmienia migracja rdzenia, a napis nie jest z nim
niczym połączony. Ten byt liczy rozjazd z odczytu, więc pasek nie orzeka
o czymś, czego rdzeń nie powiedział.

Byt stoi w korzeniu warstwy modułów, a nie w pojedynczym module, bo tę samą
potrzebę ma każdy moduł kontraktu; kopia na moduł dałaby tyle samo
rozjeżdżających się zdań o jednym stanie produktu. Bliźniakiem jest byt
rozstrzygający to samo o pokryciu komend.

Prawda pochodzi z odczytu wykazu modułów: rdzeń oddaje kody okien
operacyjnych otwieranych wraz z modułem. Czego moduł buduje, rdzeń nie wie
i wiedzieć nie może — tę połowę podaje moduł.

Okno rozmowy nie jest oknem operacyjnym modułu: katalog rdzenia niesie kod
okna rozmowy w każdym module, a montuje je scena sesji, nie rejestr
modułów. Wykaz, który tego nie odejmuje, przypisuje każdemu modułowi brak
okna, które w rzeczywistości działa.

Rozstrzygnięcia są cztery, bo każde znaczy co innego: stan nieustalony
znaczy, że rdzeń jeszcze nie odpowiedział albo odmówił, a cisza nie jest
orzeczeniem o braku; stan modułu nieznanego znaczy, że rdzeń nie wymienia
tego modułu wcale, więc katalogu jego okien nie ma czym sprawdzić; stan
katalogu pełnego znaczy, że moduł buduje każde okno operacyjne swojego
katalogu; stan katalogu szerszego znaczy, że katalog niesie kody, których
moduł nie buduje. Niezależnie od nich wychodzi wykaz kodów poza katalogiem:
kod, który moduł buduje, a rdzeń mu go nie przypisuje — to nie brak modułu,
tylko rozjazd z rdzeniem.

Odczyt jest jeden na połączenie, nie jeden na moduł — moduły pytające
o ten sam katalog dałyby tyle samo zbędnych zapytań. Katalog leży w pamięci
podręcznej przypisanej do kanału: pierwszy odczyt pyta, pozostałe czekają
na tę samą odpowiedź. Odmowa nie zostaje w pamięci, więc kolejny odczyt
ponawia pytanie.

Unieważnienie katalogu następuje przy zerwaniu połączenia: transport
ponawia połączenie pod tym samym kanałem, więc po zerwaniu można trafić na
rdzeń o innym katalogu. Dlatego byt nasłuchuje odpowiedzi odczytu wykazu
modułów na całym kanale — każdy odczyt, także cudzy, odświeża katalog bez
ani jednego zapytania stąd. Powłoka może też wymusić zapomnienie katalogu
wprost, osobnym wywołaniem.

## budowa/klient-poprzedni/src/moduly/automations/nastepne-uruchomienia.ts

Rachunek kolejnych terminów prowadzony jest po stronie klienta, ponieważ dotyczy
cykliczności jeszcze niezapisanej. Kontrakt nie ma ani pola z wykazem kolejnych
terminów, ani komendy próbnego wyliczenia, więc podgląd nie ma czego zażądać od
rdzenia. Wartością obowiązującą po zapisie pozostaje pole `AutomationSchedule.nextRunAt`
wyliczane przez rdzeń, a podgląd jest wyłącznie pomocą przy układaniu zapisu.

Rachunek biegnie na składowych czasu UTC, tak samo jak w rdzeniu, aby oba wyniki
dawały się porównać bez przeliczania strefy.

## budowa/klient-poprzedni/src/konfiguracja/nawigacja-kategorii.ts

Kategoria z wypełnionym wskazaniem nadrzędnej staje pod nią jako pozycja wcięta. Kategoria wskazująca rodzica, którego w katalogu nie ma, nie znika, lecz trafia na poziom najwyższy: ukrycie pozycji z powodu niespójności katalogu odebrałoby Operatorowi dostęp do ustawień.

Żadna pozycja nie jest wygaszona. Kategoria z fałszywym znacznikiem czynności nie jest tu w ogóle rysowana, ponieważ katalog pobierany jest bez wierszy nieczynnych.

Przebudowa wykazu nie ogłasza wyboru. Kolumna oznacza kategorię czynną i na tym kończy swoją rolę, a przebudowaniem formularza kieruje warstwa, która wykaz zamontowała. Bez tego rozdziału jedno wczytanie katalogu budowałoby formularz dwa razy.

## budowa/klient-poprzedni/src/moduly/automations/widok-ukladu.ts

Jedyną odpowiedzialnością pliku jest postać układu na ekranie. Żaden fragment nie woła
rdzenia i nie zna stanu okna; jeden bierze wywołanie zwrotne usunięcia, ponieważ przycisk
wiersza musi sięgnąć po zapis układu.

## budowa/klient-poprzedni/src/moduly/design/pasek-kontekstu.ts

Nad kanwą stoją znaczniki projektu, środowiska, modelu, trybu narzędzia
i poziomu powiększenia. Znacznik kontekstowy otwiera się naciśnięciem, więc jest
przyciskiem, a nie napisem; zdanie idzie komunikatem, ponieważ znacznik ma
zostać pigułką, zamiast rozrastać się w akapit.

Selektorów pasek nie dubluje. Wybór silnika stoi w Prompt Builderze, tryb
narzędzia w znaczniku trybu nad kanwą, a powiększenie w przyborniku. Dwie
kontrolki nastawiające jedną wartość rozjeżdżają się przy pierwszej zmianie
i przestaje być wiadomo, która mówi prawdę, więc naciśnięcie znacznika nazywa
miejsce, w którym jego selektor stoi.

Znacznik środowiska jest jedynym, który nie ma czego pokazać. Moduł jest
dostępny w dwóch środowiskach, ale które z nich jest bieżące, wie wyłącznie
powłoka: żadna komenda obszaru środowiska nie niesie, a okno modułu w rdzeniu
wskazuje moduł i sesję. Znacznik mówi to wprost, zamiast wpisywać nazwę wziętą
z niczego. Wartość znacznika nie zmienia się w cyklu życia modułu, więc stoi raz
i nie idzie w odświeżanie.

## budowa/klient-poprzedni/src/modele/stan-kont.ts

Wykaz kont i wybór konta czynnego stoją w jednym stanie, ponieważ każda zmiana
rejestru dotyka obu naraz. Usunięcie konta czynnego musi przestawić wybór,
zamiast zostawiać formularz wskazujący byt, którego rejestr już nie zawiera.

Zdarzenie zmiany konta przychodzi także wtedy, gdy konto założono albo zmieniono
na innym urządzeniu. Stan przyjmuje takie zdarzenie tą samą drogą, którą przyjmuje
własną odpowiedź rdzenia, więc wykaz nie rozjeżdża się między urządzeniami.

Żadna ścieżka nie zatrzymuje widoku. Rdzeń, który nie odda wykazu, zostawia wykaz
pusty wraz z powodem niepowodzenia, a zapis nieudany wraca jako wynik z błędem,
bez wyjątku przerywającego pracę okna.

Faza odczytu jest osobną wartością stanu, ponieważ bez niej pusty rejestr znaczy
trzy rzeczy naraz: odczytu jeszcze nie było, odczyt trwa albo rdzeń nie zna ani
jednego konta. Każdemu z tych przypadków należy się inny obraz: widok bez treści,
wskaźnik odczytu oraz stan pusty ze zdaniem opisu.

## budowa/klient-poprzedni/src/moduly/automations/uruchomienie-petli.ts

Przebiegi są dwa, ponieważ komenda `automation.queue.action` wymaga
identyfikatora kolejki, a wykaz automatyk go nie niesie: struktura
`AutomationWorkflow` pola kolejki nie ma. Kolejka powstaje więc tą samą drogą,
którą zakłada ją okno Queue Manager — komendą `queue.create` z sesją modułu
automatyk — a zaraz po niej idzie działanie rozpoczęcia ze wskazaniem automatyki.
Rdzeń zasila kolejkę krokami wyłącznie wtedy, gdy nie ma ona ani jednej pozycji;
kolejka świeżo założona jest pusta, więc uruchomienie startuje z pełnym wykazem
kroków.

Sesja kolejki jest ta sama, co w oknie Queue Manager, ponieważ pętla uruchamiana
z wykazu nie ma karty sesji: moduł jest komponentem własnym drugiej strefy strony
głównej, a nie przestrzeni roboczej sesji.

Odpowiedź niesie stan kolejki po działaniu, licznik obiegów i liczbę zleceń
oczekujących, lecz nie mówi, które kroki się wykonały. Twierdzenie o wykonanych
krokach byłoby potwierdzeniem czynności, której okno nie zmierzyło.

## budowa/klient-poprzedni/src/komponenty/stan-tresci.ts

Stan treści nie jest tym samym, co `faza-okna.ts`. Tamten byt niesie fazy ramy,
czyli atrybut fazy na powłoce i rolę dostępności pasa; ten niesie miejsce treści
wraz z pasem komunikatu i potwierdzeniem czynności, a stan trzyma w atrybucie
danych na akapicie komunikatu — po tym atrybucie sięgają arkusze modułów.

Stan błędu niesie treść błędu z kontraktu, czyli kod i komunikat, ponieważ okno
pokazujące po odmowie pusty wykaz mówiłoby o braku pozycji zamiast o nieudanym
zapytaniu. Stan pusty ma własne zdanie, ponieważ pustka bywa poprawna: instalacja
bez automatyk czy repozytorium bez zmian do zatwierdzenia nie są usterkami.

Wygląd zostaje w module: klasy noszą przedrostek modułu i pokrywa je arkusz
modułu, dlatego przedrostek jest wartością wejściową, a nie stałą. Z przedrostka
powstają cztery klasy: stanów, stanu, potwierdzenia i treści.

Plik `komponenty/odmowa.ts` składa zdanie odmowy w innym porządku i wymaga nazwy
czynności, której te okna nie podają, dlatego opis błędu powstaje tutaj osobno.

## budowa/klient-poprzedni/src/moduly/diagnostics/stan-diagnostyki.ts

Okno Recommendations Panel wyświetla rekomendacje tej analizy, którą uruchomiło
Diagnostics Center: kontrakt wiąże je polem identyfikatora analizy. Gdyby panel
trzymał własne wskazanie, pokazywałby rekomendacje analizy poprzedniej obok
wyniku nowej, a Operator nie miałby z czego rozpoznać, że patrzy na dwie różne
migawki.

Zakres czasu stoi obok analizy, ponieważ trzy komendy przyjmują czas początkowy
i końcowy niezależnie: dziennik, wykaz błędów i sama analiza. Rozjechany zakres
dałby Errors Panel z błędami jednej doby, dziennik z innej i analizę z trzeciej,
czyli zestawienie wewnętrznie sprzeczne, po którym nie da się orzec przyczyny.
Zakres pusty, z obydwoma końcami nieustawionymi, znaczy brak zawężenia i jest
stanem poprawnym, a nie brakiem.

Analiza zmienia się także pracą innego okna albo innego urządzenia tego konta,
ponieważ rdzeń dopisuje do niej rekomendacje po zakończeniu przebiegu. Dotyczy to
analizy bieżącej, więc migawka w oknach jest nieaktualna i okna mają się odczytać
ponownie. Migawkę podmienia się od razu, ponieważ zdarzenie niesie ją w całości
i drugie pytanie rdzenia byłoby zbędne.

## budowa/klient-poprzedni/src/aplikacja/router.ts

Widok opuszczony nie znika: element zostaje w dokumencie z atrybutem ukrycia,
zamiast być usuwanym i budowanym na nowo. Dzięki temu karty opuszczonego
środowiska trwają w tle i odtwarzają pełny stan po powrocie, a żadna subskrypcja
kanału się nie gubi — powrót nie zakłada drugiego okna komunikacji.

Trasa widnieje w adresie dokumentu, więc odświeżenie strony wraca tam, gdzie
Operator był, a przycisk powrotu przeglądarki działa bez kodu dodatkowego. Adres
nieznany nie zatrzymuje uruchomienia: router otwiera wtedy trasę początkową.

## budowa/klient-poprzedni/src/moduly/agents/wiersz-rozszerzenia.ts

Zestaw akcji wiersza wynika ze stanu pozycji katalogu. Pozycja niezainstalowana ma
jedną drogę — instalację; przełączenia i odinstalowania rdzeń by odmówił. Pozycja
zainstalowana dostaje przełącznik `extension.toggle`, który włącza i wyłącza bez
odinstalowania, oraz odinstalowanie.

Nazwy rodzajów rozszerzeń stoją po polsku, ponieważ kontrakt niesie kody angielskie
wyliczenia `ExtensionKind`, a wiersz czyta człowiek.

## budowa/klient-poprzedni/src/moduly/library/narzedzia-explorera.ts

Przełącznik widoku i zawężenie są w całości klienckie. Tryb wyszukiwania rozstrzyga,
którą komendą pójdzie następne szukanie: `library.file.search`, `knowledge.search` albo
obiema naraz.

Widok mapy zostaje w przełączniku mimo braku źródła współrzędnych. Pozycja usunięta
wyglądałaby na widok, którego nigdy nie przewidziano, a pozycja wybrana mówi wprost,
czego kontrakt nie niesie; rozstrzyga o tym `widoki-wykazu.ts`.

Liczba widocznych pozycji bierze się z wykazu, a nie z długości zbioru wskazanego, ponieważ
raport potrafi wskazać plik, którego świeży odczyt już nie zawiera, a katalog struktury
zawęża wynik dodatkowo.

## budowa/klient-poprzedni/src/moduly/developer/okno-project-tree.ts

Otwarcie pliku w tym oknie idzie wyłącznie przez `stan.wskazPlik`: Code Editor sam nasłuchuje
tego wskazania i sam woła komendę `developer.file.open`, więc wywołanie jej też stąd
czytałoby ten sam plik dwa razy.

Komenda `developer.tree.get` oddaje listę płaską, a hierarchię składa funkcja `zbudujDrzewo`
na podstawie pola `parentPath` każdego węzła. Żaden węzeł nie znika po cichu: samowskazanie
rodzica, cykl wzajemny, zduplikowana ścieżka i rodzic nieobecny w wykazie trafiają do
korzenia z jawną adnotacją, złożoną funkcją `wykryjCykle`. Gdyby po złożeniu drzewa nie
zostało nic, funkcja `rysujDrzewoProjektu` nazywa to stanem błędu, nie stanem pustej treści.

Pole „Katalog” niesie żądanie wypełnione przez Operatora, a pole `root` odpowiedzi rdzenia
niesie miejsce, od którego rdzeń naprawdę czytał drzewo — okno pokazuje jedno i drugie
funkcją `opiszKorzen`, bo bez tego rozróżnienia pusty katalog roboczy i katalog niewłaściwy
wyglądają identycznie.

Menu kontekstowe pod prawym przyciskiem myszy jest w tym oknie jedynym takim przypadkiem
w całej platformie — gdzie indziej „menu kontekstowe” znaczy lewy klik ikony trzech kropek.
Stąd osobny nasłuch zdarzenia `contextmenu` z wywołaniem `preventDefault`. Jedenaście
z trzynastu pozycji tego menu nie ma pokrycia w kontrakcie i stoi jako `przyciskBezKomendy`
z podanym powodem; pozycje „Skopiuj ścieżkę” (schowek) i „Odśwież” (odczyt) są zrobione
naprawdę.

Zawężenie listy węzłów po fragmencie nazwy jest czynnością wyłącznie kliencką nad węzłami
już odczytanymi — nie jedzie do rdzenia i nie jest tym samym co pole „Katalog”, które zmienia
zakres samego odczytu z rdzenia.

## budowa/klient-poprzedni/src/modele/zrodlo-kont.ts

Moduł obsługuje pięć komend obszaru `account.*` oraz zdarzenie
`account.changed`. Poświadczenie idzie w jedną stronę: kontrakt przyjmuje pole
`credential` w treści żądania i nie zwraca go nigdy, a odpowiedź niesie wyłącznie
znacznik `hasCredential`. Moduł nie ma ścieżki odczytu poświadczenia, ponieważ
nie istnieje komenda, którą mógłby o nie spytać.

Odczyt wykazu kont nie zatrzymuje widoku: wykaz, który nie dotarł, wraca pusty
i zostawia wpis w dzienniku, a widok podaje, że rdzeń nie odpowiedział,
i pozostaje czynny. Rejestr pusty po odmowie rdzenia oraz rejestr pusty na
świeżej instalacji są dwoma różnymi stanami, dlatego powód niepowodzenia dociera
do widoku osobnym powiadomieniem.

Zapisy oddają pełny wynik wywołania, ponieważ ich niepowodzenie musi stanąć przy
formularzu, a nie zniknąć w konsoli.

## budowa/klient-poprzedni/src/moduly/diagnostics/prowenancja-ocena.ts

Rodzina prowenancji ma cztery komendy odczytu i jedną zapisującą sąd człowieka
o pracy modelu. Ocena dostaje wobec tego własny, jawny chwyt w oknie, a nie skrót
klawiszowy ani kliknięcie w wiersz: czynność, która zapisuje ocenę po naciśnięciu
czegoś wyglądającego jak odczyt, jest czynnością ukrytą.

Skala trafności pochodzi z wyliczenia `ModelCallQuality` w kontrakcie i jest to
jedyne miejsce, z którego wolno ją wziąć — skala wymyślona w oknie nie miałaby
gdzie się zapisać. Wartość „bez oceny” zostaje w wyborze celowo, ponieważ ocena
nadana omyłkowo musi mieć drogę zdjęcia, a kontrakt tę wartość niesie.

Uzasadnienie jest w kontrakcie nieobowiązkowe (pole `note`) i takie zostaje
w oknie: wymuszenie go byłoby zaporą, której rdzeń nie stawia, a ocena bez słowa
nadal jest oceną. Zdanie po zapisie nazywa stan oddany przez rdzeń, ponieważ samo
potwierdzenie zapisu nie mówiłoby, co w rdzeniu ostatecznie stoi.

## budowa/klient-poprzedni/src/moduly/browser/limity-przebiegu.ts

Utrwalenie granic ma w kontrakcie własne komendy `browser.executor.limits.set` oraz `browser.executor.limits.get`, których panel jeszcze nie wywołuje. Powód stoi pod polami i bierze się z odczytu wykazu komend rdzenia, więc panel nie udaje zapisu w rdzeniu. Robi natomiast to, co zrobić może i co ma znaczenie: sprawdza scenariusz przed wysłaniem.

Stan wyjściowy jest zgodny z zasadą braku blokad domyślnych: pusta granica i pusty wykaz domen znaczą brak ograniczenia. Limit powstaje wtedy, gdy Operator go postawi, a nie wcześniej i nie domyślnie.

## budowa/klient-poprzedni/src/moduly/library/wgranie-pliku.ts

Treść pliku jedzie w zapisie base64, ponieważ takie pole niesie kontrakt, a plik
wskazany z urządzenia nie ma ścieżki widocznej dla rdzenia: przeglądarka podaje
wyłącznie nazwę. Pole ścieżki źródłowej zostaje puste, gdyż wypełnienie go nazwą
pliku byłoby wprowadzeniem danych zmyślonych.

Odczyt pliku wyprzedza wywołanie komendy i może się nie powieść, kiedy plik
zniknął albo odczyt został odrzucony. Niepowodzenie odczytu wraca tą samą drogą
co odmowa rdzenia, czyli jednym zdaniem w wierszu odpowiedzi, bez wyjątku
wywracającego widok.

Powodzenie komendy nie jest powodzeniem wgrania: rdzeń pozbawiony magazynu treści
zakłada wpis o pliku i odmawia jego podglądu. Okno pyta więc rdzeń o dostępność
treści zaraz po wgraniu i podaje to, co rdzeń odpowiedział, kosztem jednej
dodatkowej komendy, ponieważ kontrakt nie ma tańszego świadka.

Odpowiedź ze wskazaniem miejsca treści nie jest odmową: taki podgląd wraca
zarówno dla treści leżącej w magazynie, jak i dla wskazania, które nie sięga
nigdzie.
Zdanie odpowiedzi mówi więc, co przyszło, i nie orzeka ani wgrania, ani jego
braku. Podobnie werdykt odmowy oznacza, że rdzeń treści nie oddał i o niej nie
orzekł, więc okno nie orzeka wtedy ani obecności treści, ani jej braku.

Kontrolka wyboru pliku zgłasza zmianę wyłącznie przy zmianie wartości, dlatego
wyczyszczenie jej po przejęciu uchwytu pozwala wskazać ten sam plik ponownie po
wgraniu nieudanym.

## budowa/klient-poprzedni/src/mission-control/sekcja-operacji.ts

Wysycenia, kolejki, limitu ani kosztu narastającego kontrakt nie niesie, więc plakietka
wypisuje przy nich etykietę braku źródła zamiast liczb wymyślonych po stronie klienta.
Liczba zmyślona wygląda tak samo jak odczytana i wprowadzałaby w błąd przy ocenie
obciążenia kanału.

## budowa/klient-poprzedni/src/moduly/assistant/siatka-akcji.ts

Wykaz pozycji należy do rdzenia, tak samo jak wykaz modułów w nawigacji. Zestawu nie da
się zmienić z okna, ponieważ kontrakt nie ma komendy zapisu biblioteki poleceń szybkich.
Okno nazywa ten brak wprost, zamiast stawiać przycisk dostosowania, który niczego nie
zapisze.

Komenda `action.list` zasięgu modułu Assistant oddaje komendy platformy, więc opis wiersza
jest zdaniem o komendzie, a nie treścią polecenia dla asystenta. Czym wiersz jest, mówi
opis kafla: to katalog komend rdzenia, a nie biblioteka gotowych poleceń.

Wysyłanie komendy wprost z kafla byłoby zgadywaniem ładunku, więc akcja wstawia swoją
nazwę do pola, a polecenie wychodzi dopiero z paska.

## budowa/klient-poprzedni/src/moduly/library/okno-tags-collections.ts

Trzy czynności Operatora z wiersza wykazu mają drogę do rdzenia: nadanie
etykiety, utworzenie kolekcji i przypisanie zasobu do kolekcji. Słownik
etykiet i wykaz kolekcji przychodzą z rdzenia, więc okno pokazuje słownik
całego repozytorium, a nie próbkę zebraną z odczytanej strony wykazu — widać
w nim także etykiety nieużywane, których nie nosi żaden zasób, a które
istnieją i dają się usunąć.

Słownikiem można zarządzać: zmiana nazwy przechodzi po wszystkich zasobach,
łączenie wchłania etykiety duplikujące się, a usunięcie zdejmuje etykietę
z repozytorium i przy etykiecie używanej żąda potwierdzenia. Tezaurus
wychodzi z rdzenia w zapisie z relacjami; relacje ustanawia osobna komenda
i klient nie ma ich skąd wziąć sam. Mapa kolekcji zostaje wywozem okna, bo
składa się z tego, co okno już pokazuje.

Okno nie ma własnego odczytu i mieć go nie może: kontrakt nie zna komendy
katalogu etykiet ani kolekcji, więc jedno i drugie składa się z pól plików
przyniesionych przez Library Explorer. Stan tego okna jest więc stanem
tamtego odczytu. Trzy pustki są tu rozłączne: nikt jeszcze nie pytał,
odczyt trwa i rdzeń odpowiedział bez ani jednej etykiety. Wspólne zdanie
dla wszystkich trzech orzekałoby o odpowiedzi rdzenia, zanim ta przyjdzie.

## budowa/klient-poprzedni/src/moduly/assistant/panel-odpowiedzi.ts

Odsłuch prowadzi przez komendę `speech.audio.fetch`: odnośnik odpowiedzi syntezowanej
z pola `speechRef` idzie do rdzenia, bajty wracają do karty, a odtwarza je przeglądarka.

Odnośnik zostaje widoczny obok przycisku po to, żeby dało się go przekazać dalej,
a nie tylko odsłuchać na miejscu.

## budowa/klient-poprzedni/src/moduly/apps/zdania-wdrozen.ts

Składanie zdań stoi osobno od okna. Okno prowadzi rozmowę z rdzeniem: zbiera
pola, wysyła, przyjmuje odpowiedź i nazywa stan. Składanie zdań z liczb i pól
jest czynnością czystą, bez kanału i bez elementu, więc rozdział pozwala
przeczytać je w całości bez czytania okna.

Zdanie potwierdzające mówi o tym, co rdzeń oddał, a nie o tym, co okno wysłało.
Samo wypisanie pól odpowiedzi jest za małe: zamówiono konkretne środowisko
i konkretną czynność, więc gdy rdzeń odda co innego, rozejście pada wprost,
a nie w drobnym druku. Cofnięcie sprawdza się osobno, ponieważ słowo o cofnięciu
bierze się z zamówienia, a potwierdza je dopiero pole odpowiedzi wskazujące
wdrożenie, z którego cofnięto.

Powód pustego wykazu wdrożeń rozstrzyga najpierw odczyt. Pustka po udanym
odczycie znaczy rzecz ostateczną: rdzeń przejrzał bazę i wdrożeń tego okna nie
ma. To zdanie innej wagi niż brak odczytu i te dwa stany nie są zlewane.
Pozostałe trzy gałęzie dotyczą chwili, gdy nikt jeszcze nie pytał. Zero ramek
znaczy, że nic nie przyszło i nie pytano. Ramki bez pola wdrożenia znaczą, że
przyszły, ale wdrożeń w nich nie było. Ramki z wdrożeniem przy pustym wykazie są
sprzecznością wewnątrz okna, bo wykaz niczego nie zdejmuje, i zdanie nazywa ją
wprost, zamiast zaokrąglać do jednej z pozostałych. Żadna gałąź nie orzeka
o tym, czego rdzeń nie robi.

## budowa/klient-poprzedni/src/moduly/automations/widok-przekazan.ts

Pozycja wykazu mówi, skąd przyszedł scenariusz i ile ma kroków, ponieważ to
jedyne, czym operator może się kierować przed zapisem. Treść samych kroków nie
jest tu pokazywana wcale — należy do okna Workflow Builder i staje się dostępna
po zapisie scenariusza jako automatyki.

Wykaz pusty oddaje wartość `null`, ponieważ sekcja bez ani jednej pozycji nie ma
po co stać na ekranie, a zdanie o pustce niesie stan treści okna.

## budowa/klient-poprzedni/src/dostepy/stan-katalogu-roboczego.ts

Rachunek dziedziczenia oddany jest modułowi `konfiguracja/rozstrzygniecie`, temu samemu, z którego korzysta okno konfiguracji. Druga implementacja tego rachunku byłaby drugą prawdą o tej samej wartości.

Katalog roboczy jest własnością instalacji, więc ta sekcja zapisuje go na poziomie globalnym, na osi platformy. Zapis węższy — dla sesji, okna, modelu albo konta — należy do okna konfiguracji, które ma pasek punktu widzenia i pełny wybór poziomów. Zapis dokonany tam widać tutaj, ponieważ łańcuch rozstrzygnięcia mówi, skąd wartość pochodzi.

## budowa/klient-poprzedni/src/moduly/browser/czynnosci-notatek.ts

Panel składa formularz i wykaz, a w tym pliku mieszka to, co dzieje się po naciśnięciu.

Zdanie końcowe każdej czynności powstaje z odpowiedzi rdzenia, wedle `skutek-zapisu.ts`,
a nie z treści żądania: zapis notatki opisuje jej postać po zapisie, a komenda
`context.transfer` oddaje okno docelowe wraz z jego modułem i znacznikiem `transferred`.

Wykaz dopisuje notatkę także wtedy, gdy postać zapisana rozjeżdża się z wysłaną. Notatka
w rdzeniu wtedy jest, a widok ma pokazać jej prawdziwą postać obok zdania o rozjeździe.

## budowa/klient-poprzedni/src/moduly/developer/okno-code-editor.ts

Dziewięć operacji kontekstowych panelu akcji jedzie komendą
`developer.contextual.op`: rdzeń składa polecenie z treści pliku, zaznaczenia
i wskazania Operatora, a wynik wraca propozycją — okno pokazuje ją w miejscu
treści i niczego nie zapisuje. Nazwy kanoniczne i rodzaje operacji stoją
w `akcje-kanoniczne.ts`, jednym wykazem dla całego modułu.

Wskazanie pliku przychodzi ze stanu: Project Tree woła `stan.wskazPlik(...)`,
to okno nasłuchuje `stan.naZmiane(...)` i na nową ścieżkę otwiera plik przez
`zrodlo.otworzPlik`. Po udanym odczycie i po zapisie woła `stan.ustawPlik`,
żeby Project Tree i Git Panel widziały to samo. Jedyny punkt wejścia do
odczytu to `odswiez()` — wytwórnia okna nie czyta sama, złożenie modułu woła
`odswiez()` raz po montażu, a odczyt w wytwórni dałby podwójne żądanie.
`odswiez()` nie zdejmuje subskrypcji `stan.naZmiane`: ta żyje od konstrukcji
do `zamknij()`, inaczej wskazanie pliku w Project Tree przestałoby cokolwiek
robić po pierwszym odświeżeniu.

Okno nie gubi pracy Operatora cicho. Gdy treść w polu edycji różni się od
treści ostatnio wczytanego pliku, a stan wskazuje inną ścieżkę, okno woła
`tresc.potwierdzenie(...)`, nie `tresc.blad(...)` — ten drugi czyściłby
miejsce treści razem z polem edycji. Pole `.mdev-kod` zostaje widoczne
i edytowalne, a Operator ma dwa jawne wyjścia: „Zapisz plik” albo „Porzuć
zmianę i otwórz wskazany plik”. Potwierdzenie zapisu mówi, co zrobił rdzeń:
zdanie składa `zdania-odpowiedzi.ts` z pól odpowiedzi, nie z przełącznika
wersji — rdzeń robi wersję z treści sprzed zapisu, więc gdy takiej nie było,
`versionId` nie wraca wcale, a potwierdzenie nie ma prawa obiecywać punktu
powrotu.

Znacznik czystości pola nie jest treścią z rdzenia: po zapisie, którego
rdzeń nie potwierdził treścią, znacznikiem czystości zostaje treść pola —
inaczej edytor uznawałby pracę za niezapisaną na zawsze — a zdanie
potwierdzenia mówi wprost, że okno nie ma czym potwierdzić zgodności
z dyskiem. Znacznik ustawia się przed `ustawPlik`, bo ten rozgłasza zmianę
i budzi ponowne wywołanie `otworz`. Reentrantne wywołanie `otworz` kończy się
natychmiast na porównaniu `stan.plik()?.path === stan.sciezka()` — obie
wartości ustawia `ustawPlik` przed powiadomieniem — więc drugie żądanie
odczytu nie leci.

Pasek statusu edytora pokazuje wyłącznie pola, które kontrakt niesie
(język, rozmiar, wersję i czas ostatniej zmiany); pole nieobecne w
odpowiedzi jest nazwane jako niezmierzone, zamiast pokazane wartością
domyślną — „UTF-8” wypisane bez odczytu byłoby zgadywaniem, a „0 błędów”
bez lintera mówiłoby o pliku sprawdzonym, choć nikt go nie sprawdzał.

„Porzuć zmianę” jest jawnym drugim wyjściem z ostrzeżenia o niezapisanej
pracy, obok „Zapisz plik” — bez niego Operator widziałby ostrzeżenie bez
żadnej drogi naprzód poza ręczną edycją treści z powrotem do stanu
wyjściowego. Nazwy operacji nie stoją przy pasku akcji — mają jedno źródło
w `akcje-kanoniczne.ts`, żeby przycisk i komenda dołożona później mówiły
o operacji tym samym słowem.

## budowa/klient-poprzedni/src/moduly/multitasking/wskazanie-analityka.ts

Wyliczenie `WindowRole` nie ma czwartej wartości. Results Analyzer jest oknem
roli `standalone`, a odróżnia go wyłącznie wskazanie zapisane w ustawieniach
sesji. Dopisanie czwartej roli po stronie klienta rozjechałoby wykaz ról z bazą
rdzenia.

Klucz `multitasking.analityk` nie stoi w katalogu ustawień rdzenia, czyli
w tabeli `definicja_ustawienia`, więc komenda `config.set` odpowiada odmową
`validation_failed`, a wskazanie nie przeżywa odświeżenia. Zapis jest z tego
powodu sprawdzany, a odmowa trafia do zdania oddawanego wołającemu; wartość
oddana inna niż wysłana również jest nazwana wprost, ponieważ wskazanie trzyma
wtedy sam widok.

## budowa/klient-poprzedni/src/moduly/multitasking/hierarchia-decyzji.ts

Kontrakt niesie role okien — `WindowRole` z wartościami coordinator, executor
i standalone — oraz podagentów rodziny `subagent.*`, ale nie niesie wykazu poziomów
decyzji ani ich kolejności. Wykaz jest zatem opisem po stronie klienta: tłumaczy
w sekcji Monitor, kto może przerwać kogo, i nie jest bramą. Hierarchię można
konfigurować i pominąć, a interwencja pozostaje możliwa na dowolnym poziomie
w dowolnej chwili.

Tryby nakładki są dwa: obserwator patrzy i doradza, operator dodatkowo zatwierdza
i wstrzymuje kroki. Przełącznik stoi w sekcji Monitor procesu.

Kontrakt nie zna trybu nakładki. Struktura `AodStatus` niesie urządzenie, sesję, okno,
procesy przypięte i licznik procesów w biegu, a komendy `aod.mode.set` nie ma wcale.
Tryb utrwala się więc ustawieniem na poziomie sesji komendą `config.set`, a sekcja
wypisuje tę lukę kontraktu w meldunku braków.

## budowa/klient-poprzedni/src/moduly/automations/przeglad-przebiegow.ts

Zawężenie wykazu i miary niezawodności idą po stronie klienta z tego samego
powodu: komenda `automation.execution.subscribe` przyjmuje wyłącznie automatykę,
pojedynczy przebieg, okno oraz górną granicę, a oddaje przebiegi w całości. Nie
ma w niej pola zapytania, zakresu dat ani zestawienia miar, a wykaz jest już
w oknie, więc pytanie rdzenia o to samo drugi raz nic by nie wniosło.

Miary liczą się z pól, które kontrakt naprawdę niesie: stanu przebiegu oraz
znaczników czasu rozpoczęcia i zakończenia. Miary kosztu — tokeny i wywołania
modelu — składa widok (`widok-przebiegow.ts`), ponieważ ich pola są
nieobowiązkowe, a suma wymaga podania obok siebie liczby przebiegów, które je
wypełniły.

## budowa/klient-poprzedni/src/dostepy/komunikat-czynnosci.ts

Każde naciśnięcie kontrolki w sekcji dostępów kończy się zdaniem przy tej
kontrolce, bo po milczeniu widoku Operator nie wie, czy model ma już dostęp,
czy jeszcze nie. Powodzenie mówi, co się stało, niepowodzenie powtarza treść
odpowiedzi rdzenia.

Odmowa własna widoku ma kształt wyniku komendy, więc odmowa wystawiona bez
wysyłki — na przykład brak okna rozmowy — wygląda dla widoku jak odmowa rdzenia
i tą samą drogą trafia do zdania przy kontrolce. Widok nie potrzebuje przez to
drugiej ścieżki obsługi.

## budowa/klient-poprzedni/src/mission-control/mission-control.ts

Pulpit powołuje sekcje w ustalonej kolejności: aktywność modeli wraz z pasem
decyzji, czyli eskalacją koordynatora; operacje modeli, czyli kanały i koszt;
matrycę sesji wraz z pasem relacji, pokazującą procesy biegnące równolegle;
sekcję zakładania pracy jeszcze nierozpoczętej; oraz trzy kolumny procesów,
kolejek i zespołu.

Pasek działań słucha nadajników pulpitu na własną rękę, żeby każde działanie
miało widoczny skutek niezależnie od tego, czy ktokolwiek podpiął odbiorcę
zewnętrznego.

## budowa/klient-poprzedni/src/moduly/browser/stan-okna.ts

Stany są rozdzielone, ponieważ brak zapytania, zapytanie w toku i odmowa rdzenia to trzy
różne sytuacje. Zlanie ich w jedno kazałoby Operatorowi zgadywać, czy czekać, czy działać;
ten sam podział prowadzi `dostepy/stany-odczytu.ts`.

Nieudane odświeżenie zostawia to, co Operator już widział, a ponowienie odczytu nie
zdejmuje komunikatu błędu. Treść znika wyłącznie na czas odczytu.

Nazwa fazy i jej znakowanie pochodzą z `komponenty/faza-okna`. Wartość trafia do atrybutu
`data-faza`, po którym sięgają arkusze stylów i sprawdziany; własny zestaw nazw w module
znaczyłby, że ten sam stan okna nazywa się gdzie indziej inaczej.

Między złożeniem okna a pierwszym odświeżeniem Operator widzi zdanie, które i tak
zobaczyłby przy pustym wyniku.

## budowa/klient-poprzedni/src/moduly/agents/stan-okna.ts

Rdzeń odpowiada na każde wywołanie, również odmową, a odmowa jest pokazywana
w oknie, w którym została wywołana, a nie wyłącznie w dzienniku przeglądarki.

Stan nie zastępuje treści, tylko ją przesłania. Gdy okno wraca do fazy `gotowe`,
wcześniejsza treść jest nietknięta, dzięki czemu nieudane odświeżenie nie kasuje
tego, co stało już na ekranie.

Zestaw faz i znakowanie powłoki są wspólne dla wszystkich modułów i mieszkają
w `komponenty/faza-okna`. W module Agents zostaje wyłącznie to, czym różni się on
od pozostałych: własne klasy `da-stan*` oraz chowanie treści na czas ładowania.

## budowa/klient-poprzedni/src/moduly/design/plansza-kompozycji.ts

Przeciąganie jest miejscowe: wskaźnik zostaje przechwycony na warstwie przez
`setPointerCapture`, przesunięcie liczy się w jednostkach kompozycji, czyli po
podzieleniu przez powiększenie, a warstwa zablokowana nie rusza się wcale.

Kliknięcie w puste płótno zdejmuje zaznaczenie, ponieważ inaczej nie da się wyjść
z zaznaczenia wielokrotnego bez trafienia we właściwą warstwę.

Plansza nie zna rdzenia. Ruch warstwy zmienia zapis kompozycji, a do rdzenia jedzie
dopiero zapis całości komendą `design.board.update` — z okna, a nie z planszy.

## budowa/klient-poprzedni/src/aod/wyciszenie-braki-kontraktu.ts

Kontrakt niesie wyciszenie nakładki jako byt rdzenia: strukturę wyciszenia wraz
z rodzajem, zakresem i chwilą końca, komendy odczytu i zapisu wyciszenia,
zdarzenie rozgłaszane na pozostałe powłoki, nośnik sygnału klas zdarzeń, moduł
przy stanie i podpowiedzi nakładki oraz kategorię ustawień z pozycją reguł
wyciszania.

Pozostaje jeden brak i leży po stronie nakładki, a nie kontraktu: wołacze tego
okna nadal piszą do magazynu stanowiska, więc wyciszenie założone w tym oknie nie
dojdzie do drugiej powłoki, dopóki magazyn nie zostanie przełożony na rdzeń.
Magazyn wyciszeń jest podawany z zewnątrz właśnie po to, aby przełożenie nie
wymagało zmiany ani jednego wołacza.

Zdanie braku nazywa, czego brakuje i po czyjej stronie brak leży. Nazwy komend
pochodzą z wyliczenia komend kontraktu, a nie z literałów, więc zmiana nazwy
w kontrakcie zostaje wychwycona przy kompilacji. Pozycja braku, której komenda
znosząca weszła już do wykazu komend, wypada z wykazu braków czynnych, a pozycja
bez komendy znoszącej pozostaje, ponieważ takiego braku nie zniesie żadna komenda.

Plik nie obchodzi braku: nie wysyła komendy zastępczej i nie udaje zapisu
w rdzeniu. Wyciszenie pozostaje stanem tego okna i tak jest nazwane w menu.

## budowa/klient-poprzedni/src/dostepy/stany-odczytu.ts

Pusty wykaz punktów dostępu znaczy co innego, gdy rdzeń jeszcze nie
odpowiedział, co innego, gdy odmówił, i co innego, gdy rejestr jest pusty.
Zlanie tych przypadków w jeden kazałoby zgadywać, czy czekać, czy działać,
dlatego pas stanów rozróżnia je osobnymi fazami odczytu.

Przycisk ponowienia wyzwala odczyt i zostawia komunikat na miejscu, dopóki
sytuacja nie ustanie; sekcja pozostaje przez cały ten czas czynna.

Biblioteka `komponenty/` nie ma jeszcze reguł komunikatu blokowego, więc
komunikat stoi na klasie komunikatu tej sekcji.

## budowa/klient-poprzedni/src/aod/zrodlo-komend.ts

Sekcje okna dostają ze źródła czynności, a nie stałe komend, więc żadna sekcja
nie powtarza opakowania kanału. Komenda `aod.voice.command` idzie bez mikrofonu:
pola odwołania do nagrania i transkrypcji są w kontrakcie opcjonalne, więc
polecenie wydane samą treścią jest wywołaniem pełnoprawnym. Granicę nazywa plik
`sekcja-glosu.ts`.

Żadne żądanie nie niesie identyfikatora urządzenia. Rdzeń rozumie przez ten
identyfikator numer wiersza katalogu maszyn — adapter nakładki czyta go jako
liczbę całkowitą, tym samym prawem co rodzina komend punktów dostępu.
Identyfikator sesji dostaje odmowę niepowodzenia sprawdzianu, a kolumna nakładki
przerywa wczytywanie na pierwszym błędzie stanu, więc jedna zła wartość zabiera
całą treść nakładki. Pole puste rdzeń przyjmuje, a żadna komenda kontraktu nie
mówi klientowi, którym numerem katalogu maszyn jest ta maszyna, więc podstawienie
zgadniętej jedynki byłoby wartością zmyśloną.

## budowa/klient-poprzedni/src/modele/wykaz-kategorii.ts

Kategoria spoza warstw kontraktu nie znika: dostaje własną grupę na końcu, żeby
było widać, co przyszło z rdzenia.

Wiersz mówi dwie rzeczy: czy kategoria jest obowiązkowa i czy ma treść zapisaną
na osi czynnej. Brak treści na osi nie znaczy braku treści w ogóle — obowiązuje
wtedy zapis z osi szerszej i wiersz nazywa to wprost.

## budowa/klient-poprzedni/src/moduly/diagnostics/zuzycie-zrodlo.ts

Zestawienie idzie po jednym wymiarze naraz i tak stanowi kontrakt wprost:
zestawienie po dwóch wymiarach naraz jest dwoma żądaniami, nie jednym. Zakładka
nie składa więc dwóch odpowiedzi w jedną tabelę — pokazuje wymiar wybrany i mówi,
który to. Raport bierze wymiary wykazem, ponieważ jest jednym wytworem za okres,
a nie widokiem; wykaz pusty znaczy wszystkie wymiary.

Koszt niepełny jest tu osobnym pojęciem, a nie zaokrągleniem. Kontrakt niesie
udział wywołań objętych cennikiem oraz wywołania pominięte w koszcie, a pole
kosztu opisuje jako puste, gdy cennika kanału nie ma: puste znaczy brak cennika
kanału, nie koszt zerowy. Zakładka przenosi to rozróżnienie na ekran, ponieważ
złożenie kosztu z zer dałoby liczbę wyglądającą na pomiar.

Rodzina zdarzeń nie ma, więc nic tu nie nasłuchuje: zużycie odczytuje się na
żądanie. Pusty wykaz zestawienia jest poprawną odpowiedzią i sprawdzian kształtu
tego nie myli z odpowiedzią bez kształtu — wymagana jest sama tablica oraz
granice okresu, bo bez nich zakładka nie wie, o czym mówi.

## budowa/klient-poprzedni/src/moduly/browser/pasek-kontekstu.ts

Kolejność w pasku jest kolejnością rejestracji rozszerzeń w module.

Wyzwalacze warstwy czwartej nie stoją w pasku, dopóki Operator nie włączy trybu
administracyjnego. Nie jest to blokada: rozszerzenie ze skrótem otwiera się
skrótem niezależnie od trybu, a przycisk nie zostaje wygaszony, tylko zdjęty ze
sceny.

Stan wyzwalacza niesie atrybut wciśnięcia, a nie barwa: pasek czyta się także bez
rozróżniania barw i z czytnika ekranu.

## budowa/klient-poprzedni/src/moduly/browser/pasek-adnotacji.ts

Pasek leży na rysunku i nie zabiera wysokości podglądowi, więc pozycjonowanie należy do warstwy płótna (`adnotacja.css`), a nie do układu okna.

Wybór narzędzia i barwy niesie atrybut `aria-pressed`, a nie klasa arkusza. Stan wypowiedziany atrybutem czyta czytnik ekranu, widzi go sprawdzian i sięga po niego arkusz.

Nazwa próbki barwy pada w atrybucie, ponieważ bez niej wybór byłby sygnalizowany samym kolorem, czego zabrania zasada dostępności żetonów stanu (`motyw/stany.css`).

Pole napisu pozostaje widoczne i edytowalne zawsze, a znaczenie ma wyłącznie przy narzędziu stawiającym tekst; różnicę niesie znacznik odczytywany przez arkusz.

Żaden przycisk nie gaśnie. Przycisk dodania do rozmowy przy pustym płótnie pozostaje naciskalny i odpowiada zdaniem o tym, że nie ma czego wysłać.

## budowa/klient-poprzedni/src/konfiguracja/wskaznik-zasiegu.ts

Plakietka jest naciskalna zawsze, także wtedy, gdy klucz nie ma ani jednego zapisu.
Rozwinięcie mówi wtedy wprost, że obowiązuje wartość domyślna katalogu, i od razu daje
adres, pod którym można to zmienić.

Punkt widzenia nanosi się na wybór celu tylko raz, ponieważ poziom zapisu przestawiony
ręcznie przy jednym polu nie ma wracać do punktu widzenia przy kolejnym odświeżeniu stanu.

## budowa/klient-poprzedni/src/moduly/agents/zrodlo-rozszerzen.ts

Katalog stoi osobno od `zrodlo-zaplecza`, ponieważ obie warstwy odpowiadają za co innego.
Zaplecze niesie rejestry, z których moduł Agents tylko korzysta — kanały modelu, mosty MCP,
okna sesji — a katalogiem rozszerzeń moduł zarządza: instaluje i odinstalowuje pozycje.

Katalog obejmuje rodzaje wymienione w wyliczeniu `ExtensionKind`: serwer MCP, wtyczkę,
integrację API i umiejętność.

O tym, czym jest instalacja pozycji, rozstrzyga rdzeń, a nie ta warstwa. Źródło przekazuje
pola kontraktu i oddaje odpowiedź bez dopowiedzenia, żeby znaczenie czynności zostało
w jednym miejscu.

## budowa/klient-poprzedni/src/dostepy/wykaz-punktow.ts

Maszyna za mostem protokołu i katalog na dysku urządzenia to dwa różne ryzyka
i dwie różne drogi sprawdzenia, więc wykaz rozdziela je rodzajem. Zlanie ich
w jedną listę kazałoby czytać rodzaj z każdego wiersza z osobna.

Karty przeżywają przeliczenie: zmiana stanu nanosi wartości na karty już
zbudowane, a przebudowa następuje wyłącznie po zmianie składu wykazu. Inaczej
wybór trybu i zaznaczenie korzeni ginęłyby przy każdym zdarzeniu z rdzenia. Do
składu wchodzą korzenie, bo z nich powstaje wybór korzeni w karcie; stan punktu
i czas sprawdzenia do składu nie wchodzą, ponieważ zmieniają się przy każdym
sprawdzeniu, a przebudowa zabrałaby wybór trybu wpisany chwilę wcześniej.

Stan pusty należy się wyłącznie odczytowi zakończonemu powodzeniem. W trakcie
pytania rdzenia zdanie o braku punktów byłoby nieprawdą, a po niepowodzeniu
odczytu pomyłką co do przyczyny.

## budowa/klient-poprzedni/src/konfiguracja/wybor-kontrolki.ts

Formularz okna konfiguracji powstaje z katalogu nastaw, więc rozdzielenie
rodzaju wartości na kontrolkę ma dokładnie jeden punkt. Nowy rodzaj wartości
w kontrakcie oznacza nowy przypadek w tym rozdzieleniu i nowego budowniczego
obok, a nie nowy ekran ani zmianę w formularzu i w panelu kategorii.

Rodzaj nieznany nie gasi pola. Rdzeń nowszy od klienta może przysłać rodzaj
wartości, którego ten klient nie zna; zamiast pustego miejsca staje wtedy pole
tekstowe z ostrzeżeniem, wartość pozostaje odczytywalna i zapisywalna jako
napis, a nazwa rodzaju idzie do dziennika przeglądarki.

## budowa/klient-poprzedni/src/moduly/automations/zrodlo-automations.ts

Każda czynność źródła oddaje wynik opakowany, nie samą treść: okna mają
obowiązkowy stan błędu, więc źródło nie połyka odmowy ani nie podstawia
pustego wykazu w jej miejsce — widok musi odróżnić brak danych od
nieudanego zapytania. Po stronie klienta nie stoi drugi silnik kolejek:
Queue Manager zakłada kolejkę jedną komendą i posuwa ją drugą, która po
stronie rdzenia jedzie tym samym adapterem kolejek co pętla sesyjna.

Odczyt zależności układu jest odrębny od zapisu definicji orkiestratora,
choć obie komendy mówią o tym samym grafie: zapis definicji zapisuje układ
w całości i oddaje go przy okazji, a odczyt zależności wyłącznie czyta. Sam
odczyt jest potrzebny, bo okno musi pokazać układ, którego nikt wcześniej
nie zapisywał z tego okna. Zapis pojedynczej zależności nie jest odmawiany:
układ z cyklem zostaje zapisany, a zastrzeżenia wracają w wyniku — ocena
układu nie jest więc oceną samej czynności zapisu, i okno rozdziela te dwa
zdania.

Cztery dopełnienia układu niosą to, czego krawędź nie wyraża. Krawędź mówi,
że kroki się schodzą. Bramka mówi, kiedy tory scalają się w jednym kroku;
grupa mówi, że zbiór kroków biegnie razem; kompensacja mówi, co zrobić, gdy
przebieg pękł w pół; spięcie mówi, że kolejki automatyki prowadzi silnik
środowiska MultitaskingAI.

Odczyt okien komunikacji platformy służy jednej sprawie: przekazanie
scenariusza z innego modułu zakłada okno modułu Automations, i to w jego
konfiguracji sesji mieszka przeniesiony komplet. Bez wykazu okien nie ma
jak dojść do tego, komu przekazano scenariusz.

Zdarzenie zmiany powiązania automatyki przychodzi także wtedy, gdy
scenariusz wpiął tu inny moduł, więc to ono jest sygnałem, że wykaz
automatyk okna jest już nieaktualny.

## budowa/klient-poprzedni/src/moduly/browser/czynnosci-paska-dolnego.ts

Jedyną odpowiedzialnością pliku są wywołania rdzenia paska dolnego. Pasek składa
kontrolki, a rozmowa z rdzeniem stoi osobno. Milczenie po naciśnięciu jest zakazane.

Zdanie końcowe buduje się z treści odpowiedzi, a nie z samego znacznika powodzenia;
składa je `skutek-zapisu.ts`. Rdzeń pobiera dziś stronę bez uruchamiania przeglądarki,
wedle `przegladarka_pobieranie.go`, więc odpowiedź udana potrafi przyjść z pustym polem
`screenshotRef` — okno nie ma wtedy prawa twierdzić o obrazie strony. Zdanie o skutku
czyta pole odpowiedzi, więc pozostanie prawdziwe także wtedy, gdy zrzuty zacznie oddawać
uchwyt `browser.screenshot.capture`.

Zdanie o cytacie ma się opierać na wierszu, który wrócił z rdzenia, a nie na tym, że okno
cytat wysłało; dlatego zamówienie spisane jest raz i to samo idzie w obie strony.

## budowa/klient-poprzedni/src/moduly/apps/warsztat-plikow.ts

Zawartość warsztatu przychodzi dwiema drogami: `apps.workspace.list` daje wykaz z bazy,
a zdarzenie `apps.workspace.changed` — zmianę zaszłą gdziekolwiek indziej. Ten plik jest
jednym miejscem, w którym oba źródła się spotykają.

Tożsamością pliku jest para złożona z warstwy i ścieżki, a nie sama ścieżka. Klucz naturalny
warsztatu to okno, warstwa i ścieżka, przy czym okno rozstrzyga się na poziomie całego modułu,
więc zostają dwa człony. Zbiór trzymany jedną listą po samej ścieżce zlewałby `src/main.ts`
warstwy przedniej z `src/main.ts` warstwy zaplecza w jeden wiersz i pokazywał treść nie tego
pliku, co trzeba.

Ścieżek nie przycinamy przy porównaniu. Rdzeń zapisuje ścieżkę co do znaku, więc dwie ścieżki
różniące się znakiem niewidocznym są dla niego dwoma plikami i tak samo muszą być tutaj.
Zrównanie ich w kliencie kazałoby zbiorowi nadpisać wpis, którego rdzeń nie tknął.

Wchłonięcie wykazu zastępuje warstwę, a nie dokłada do niej. Odpowiedź rdzenia jest pełnym
stanem warstwy w chwili odczytu, więc scalanie zostawiłoby w wykazie plik, który w bazie już
nie istnieje. Zdarzenie o zmianie bez pola pliku jest mijane, zamiast zakładać wiersz z pustki.

## budowa/klient-poprzedni/src/moduly/design/skutek-designu.ts

Zdania o skutku są wydzielone z okien, ponieważ okna składają kontrolki, a to
jest ocena odpowiedzi rdzenia. Zdanie zbudowane z tego, co okno wysłało, bywa
prawdziwe przypadkiem i skłamie, gdy rdzeń zapisze co innego.

Przy przekazaniu komenda `context.transfer` nie sprawdza katalogu modułów:
przepisuje pole `targetModuleId` do okna docelowego bez zmiany i oddaje wartość
prawdziwą w polu `transferred` także dla modułu, którego nie ma. Jedynym polem
mówiącym, gdzie zasób wylądował, jest `window.moduleId` z odpowiedzi, więc zdanie
odpowiada na oba pytania: czy przeniesienie się odbyło i dokąd.

Przy zapisie kompozycji odpowiedź niesie całą kompozycję odczytaną po zapisie,
wraz z nazwą i wykazem warstw; liczba warstw, która wróciła, jest jedyną miarą
tego, ile ich leży w rdzeniu. Przy nadaniu etykiet zestaw zastępuje poprzedni,
więc różnica wobec zestawu zamówionego znaczy, że zapis nie jest tym, o który
proszono, i musi być odmową, a nie milczeniem.

## budowa/klient-poprzedni/src/moduly/multitasking/panel-obsady.ts

Panel obsady odpowiada za scenę; czym rola jest nadawana i skąd bierze się zdanie o skutku, rozstrzyga moduł nadania ról, który woła komendę nadania roli oraz komendę zmiany wcielenia. Okno powstaje na wzór okna istniejącego: założenie okna wymaga kanału modelu, katalogów roboczych, zasięgu wykonania i trybu uprawnień, więc klient bierze te pola z okna, które w sesji już stoi. Sesja bez ani jednego okna nie daje wzorca i panel mówi to wprost, zamiast wysyłać żądanie skazane na odmowę. Panel ma własne miejsce stanu treści i nie oddaje meldunków wywołaniu zwrotnemu sceny, ponieważ inaczej odmowy rdzenia przepadałyby w ciszy.

Wykaz okien obcych wychodzi nieosadzony: jest jedną z sekcji panelu, a o tym, gdzie sekcja stoi i czy jest zwinięta, rozstrzyga układ z rdzenia. Osadzenie go na sztywno byłoby drugą prawdą o kolejności.

## budowa/klient-poprzedni/src/moduly/design/zrodlo-zaplecza.ts

Komendy `window.list` i `window.state.get` dają okno modułu i jego stan; są wspólne
każdemu oknu operacyjnemu.

Komenda `channel.list` daje rejestr kanałów modelu. Jest najbliższym istniejącym
odpowiednikiem wyboru silnika generującego z panelu akcji Prompt Buildera, bo pole
`engine` promptu przyjmuje identyfikator kanału.

Komendy `module.list` i `context.transfer` dają katalog modułów rdzenia oraz jedyną
zbudowaną drogę przekazania międzymodułowego; nią idzie wysłanie zasobu do modułu
docelowego z Assets Panel. Wykaz modułów docelowych bierze się z rdzenia, nigdy z kopii
katalogu po stronie klienta.

Zdarzenie `progress.changed` niesie telemetrię postępu. Generowanie zasobu nim nie jedzie:
komenda `design.asset.generate` wraca w jednej turze, bez pola `processId`, i nie wysyła
po drodze żadnego `progress.changed`.

## budowa/klient-poprzedni/src/moduly/browser/zrodlo-okien.ts

Każda komenda obszaru `browser.*` wymaga pola `windowId`, a moduł dostaje
z powłoki wyłącznie identyfikator sesji. Ustalenie okna stoi więc w jednym
miejscu: pytanie o nie osobno w każdym oknie modułu dałoby kilka różnych okien
dla jednego modułu. Odczyt okien sesji i odczyt stanu okna należą do obszaru
okien, nie do obszaru `browser.*`, dlatego stoją w źródle osobnym wobec
`zrodlo-browser.ts`.

Rdzeń przestawia moduł okna komendą `workspace.enter`, więc oknem przeglądarki
jest to, którego pole `moduleId` równa się kodowi modułu.

## budowa/klient-poprzedni/src/moduly/agents/zrodlo-wersji-eksperta.ts

Historia i archiwum stoją osobno od `zrodlo-agentow.ts`, ponieważ odpowiadają na inne
pytanie: biblioteka mówi, jacy eksperci są, a to źródło mówi, co się z danym ekspertem
działo i gdzie poszedł. Trzymanie obu w jednym pliku łączyłoby dwie odpowiedzialności
i przekraczało próg objętości pliku.

Wersję wskazuje identyfikator, a nie numer, bo tak przyjmuje ją komenda
`agent.version.restore`. Numer jest porządkiem historii, nie tożsamością wersji: po
przywróceniu numery rosną i ten sam numer znaczyłby co innego.

## budowa/klient-poprzedni/src/moduly/browser/formularz-notatki.ts

Błąd pola nie kasuje formularza: alert liniowy staje pod polem, pole pozostaje edytowalne, a wpisana treść zostaje na miejscu. Alert znika razem ze znacznikiem niepoprawności, ponieważ pole zdjęte z alertu, a wciąż oznaczone atrybutem `aria-invalid`, czytałoby się jako błędne bez powodu.

Wypełnienie formularza treścią istniejącej notatki oddaje fałsz, gdy źródła notatki nie ma już w wykazie okna. Ster wraca wtedy do pozycji bez powiązania, a okno ma o czym powiedzieć. Ciche przełknięcie tej różnicy zamieniłoby zapis w notatkę o innym powiązaniu niż pierwowzór.

Wybór źródła i modułu docelowego prowadzą dwa stery zamiast dwóch natywnych list wyboru: nastawa stoi na uchwycie, a mechanizm rozwijania jest jeden dla całego produktu (`komponenty/menu-drzewo.ts`). Podpis pozostaje przy każdym sterze, ponieważ oba stoją w rzędzie pól formularza, gdzie sama wartość nie mówi, czego dotyczy.

## budowa/klient-poprzedni/src/moduly/assistant/selektor-profilu.ts

Osobnej rodziny komend `assistant.profile.*` kontrakt nie ma i mieć nie musi: rodzaj
komponentu nazywa się w wyliczeniu `ComponentKind` wprost profilem asystenta, więc wykaz
profili jest zawężonym odczytem `component.list`.

Do rdzenia jedzie `Component.targetId`, a nie `Component.id`. Rdzeń rozstrzyga `profileId`
po kolumnie `profil_asystenta.identyfikator_zewnetrzny`, a bytem docelowym komponentu jest
właśnie ten kod. Wysłanie identyfikatora kafla byłoby wskazaniem profilu, którego rdzeń nie
zna — zlecenie poszłoby wtedy bez warstwy profilu i nikt by się o tym nie dowiedział.

Komponent bez bytu docelowego zostaje w wykazie i mówi o tym w swojej nazwie. Wybranie go
nie wstawia niczego do pola, bo nie ma czego wstawić; ukrycie takiego wiersza zataiłoby fakt,
że kafel profilu stoi w strefie 2 bez profilu pod spodem.

## budowa/klient-poprzedni/src/moduly/library/stan-okna.ts

Stany są rozdzielne, ponieważ brak zapytania, zapytanie w toku i repozytorium bez ani
jednego pliku to trzy różne rzeczy; zlanie ich w jedno kazałoby Operatorowi zgadywać,
czy czekać, czy działać. Stan pusty niesie dokładnie jeden przycisk pierwszej akcji.

Forma stanu pustego jest jedna: ikona, tytuł i opis, treść wyśrodkowana, bez wariantów
klasy — różnicuje ją sama treść. Stopnie pisma i szerokość łamania niesie wspólny arkusz
`komponenty/drobne.css` klasami `.dn-pusty-stan-tytul` oraz `.dn-pusty-stan-opis`, a nie
arkusz modułu; nadpisanie ich u siebie dałoby dwie formy tego samego stanu.

Ikona należy wyłącznie do stanu pustego. Ładowanie ma własny nośnik `.dn-spinner`, a błąd
kreskę po lewej stronie; trzeci rysunek nad nimi niczego by nie dopowiedział. Widoczność
rozstrzyga arkusz po atrybucie `data-faza`, tak samo jak kreskę błędu.

Nazwy faz i znakowanie powłoki pochodzą ze wspólnego `komponenty/faza-okna`, żeby ten sam
stan nazywał się w całym drzewie tak samo.

Okno stoi w stanie pustym przez chwilę między zbudowaniem układu a pierwszą odpowiedzią
rdzenia i jest to wtedy jedyne zdanie, które Operator czyta.

## budowa/klient-poprzedni/src/moduly/design/wejscie-rozmowy.ts

Plik rozstrzyga, którędy wynik pracy prowadzonej w oknie rozmowy wchodzi na
okno robocze modułu. Kładzeniem warstwy zajmuje się Design Board, a samą
subskrypcją złożenie modułu.

Miarą jest okno zasobu, a nie jego identyfikator. Komenda generowania zasobu
przepisuje okno żądania na pole okna zasobu, więc zasób niesie okno, z którego
padło zlecenie. Zdarzenie zmiany zasobu rozgłaszane jest natomiast dla każdego
zasobu odpowiedzi, zanim odpowiedź wróci do autora żądania, więc po
identyfikatorze rozróżnić się nie da: w chwili zdarzenia własnego generowania
klient nie zna jeszcze identyfikatora własnego zasobu. Miara oparta na oknie
wyścigu nie ma, a jedyną drogą zlecenia modułowi spoza jego okien jest rozmowa.

Puste okno modułu miary nie psuje, ponieważ Prompt Builder bez okna modułu
odmawia zlecenia, więc zasób, który wtedy przyszedł, na pewno nie jest jego.
Miara nie rozstrzyga natomiast przypadku, w którym model wołałby komendę
generowania z identyfikatorem okna tego modułu: wynik rozmowy byłby wtedy nie
do odróżnienia od wyniku Prompt Buildera, na kanwę sam by nie wszedł
i trafiłby do panelu zasobów.

Skutek zdjęcia zasobu istnieje mimo tego, że rdzeń komendy usuwającej dziś nie
ma. Wyliczenie zmian dopuszcza usunięcie, a pominięcie tej gałęzi kazałoby oknu
położyć na kanwie zasób ogłoszony jako usunięty.

Nasłuch wejścia rozmowy jest drugą subskrypcją tego samego zdarzenia, obok tej,
którą prowadzi stan modułu. Rozdzielenie jest celowe: stan wciąga każdy zasób
do wykazu zarządcy, a nasłuch odpowiada wyłącznie na pytanie o pochodzenie
zasobu i wykazu nie zmienia. Zdarzenie zasobu z własnego okna modułu jest
pomijane, bo melduje o nim Prompt Builder wraz ze słowem modelu.

Znacznik obecności na kanwie jest sprawozdaniem okna roboczego, a nie
zapowiedzią: zasób już dołożony do kompozycji nie jest kładziony drugi raz
i zdanie mówi to wprost, zamiast zostawiać licznik, który podskoczył bez powodu.

## budowa/klient-poprzedni/src/asystent-plywajacy/dostepnosc-mowy.ts

Stan drogi głosowej ma cztery ogniwa. Droga polecenia do rdzenia jest zbudowana:
komenda `assistant.voice.command` stoi w kontrakcie i ma uchwyt w rdzeniu, więc
polecenie dojeżdża, zakłada zlecenie i wraca z jego stanem. Rozpoznania mowy nie
ma: żądanie z samym odwołaniem do nagrania, bez transkrypcji, zostawia zlecenie
w stanie oczekiwania, ponieważ rdzeń czeka na tekst i sam go nie wytworzy. Nie ma
też czym wytworzyć odwołania do nagrania — żadna komenda kontraktu nie przyjmuje
nagrania z przeglądarki, więc klient nie ma dokąd wysłać tego, co by nagrał,
a nagrywanie do pamięci i wyrzucanie nagrania byłoby atrapą mikrofonu. Syntezy
odpowiedzi nie ma: odwołanie do mowy w odpowiedzi wraca puste, a pole żądania
o odczytanie odpowiedzi nie ma kolumny w schemacie. Komenda
`translate.speech.synthesize` syntezuje treść panelu tłumaczenia, żąda
identyfikatora tego panelu i odpowiedzi asystenta nie odczyta.

## budowa/klient-poprzedni/src/dostepy/zrodlo-nadan.ts

Trzy komendy zapisujące zwracają nie tylko zmienione nadanie, lecz komplet nadań
okna po zmianie, ponieważ przestawienie kolejności albo oznaczenia głównego
dotyka pozostałych wierszy, a widok musi je zobaczyć w jednej odpowiedzi.

## budowa/klient-poprzedni/src/moduly/automations/maszyna-stanow.ts

Model stanów nie jest wymyślony po stronie okna: stany są kompletem wyliczenia
`AutomationExecutionStatus`, a przejścia wynikają z działań, które kontrakt na
przebiegu dopuszcza — uruchomienia, wstrzymania, wznowienia, zatrzymania oraz
zamknięcia przebiegu powodzeniem albo błędem.

Widok jest wykazem, a nie rysunkiem: przejść jest siedem, a rysunek siedmiu
strzałek nie mówi więcej niż siedem zdań i kosztuje drugą kanwę w module. Pusty
stan bieżący znaczy, że żaden przebieg nie jest wskazany; wykaz stoi wtedy bez
wyróżnienia i jest samym modelem.

## budowa/klient-poprzedni/src/dostepy/lista-nadan.ts

Kolejność nadań rozstrzyga, w jakiej postaci trafią one do konfiguracji mostów,
a nadanie główne wskazuje punkt, od którego model zaczyna.

Lista przebudowuje się przy każdej zmianie zbioru i jest to wybór świadomy:
zmiana kolejności albo oznaczenia głównego przestawia wszystkie wiersze naraz,
więc nanoszenie wartości na wiersze już zbudowane byłoby trudniejsze i mniej
wierne niż zbudowanie ich od nowa z odpowiedzi rdzenia.

## budowa/klient-poprzedni/src/moduly/assistant/filtr-zlecen.ts

Zawężenie wykazu zleceń liczy się w oknie, a nie w rdzeniu. Komenda
`assistant.action.status` nie ma pola stanu w żądaniu — `adapter_modul_asystent_czynnosci.go`
czyta po `windowId` albo po `actionId` — więc wykaz i tak przychodzi w całości,
a odczyt na każdą zmianę pozycji listy byłby wywołaniem bez nowej treści.

Grupy filtra są złożone ze stanów kontraktu, nie z własnych nazw stanów.
Anulowane stoją osobno od nieudanych, bo anulowanie jest decyzją Operatora,
a nie niepowodzeniem zlecenia; zlanie obu w jedną pozycję kazałoby czytać własną
decyzję jako usterkę.

Wykaz otwiera się bez zawężenia, ponieważ zlecenie zakończone chwilę wcześniej
jest dla Operatora wchodzącego do modułu tak samo istotne jak zlecenie w toku,
a wykaz zawężony od razu wyglądałby na pusty rdzeń.

## budowa/klient-poprzedni/src/moduly/design/zestawy-zetonow.ts

Motyw jest własnością powłoki. Zestaw żetonów jest bytem obok niego: panel
odczytuje żetony motywu obowiązującego, a Operator może je odłożyć jako
zestaw, wczytać cudzy system i wydać go do kodu. Zapis zestawu nie zmienia
wyglądu produktu ani o jeden piksel i panel mówi to wprost, zamiast
zostawiać Operatora z domysłem.

Import wnosi role nieznane systemowi projektowemu klienta i je wymienia:
odrzucenie ich byłoby zgubieniem pracy, przemilczenie — obietnicą, że
wszystko pasuje. Rola nieznana wchodzi do zestawu, a jej nazwa wraca
w odpowiedzi rdzenia i panel wypisuje ją w całości.

Postacie wydania — CSS, SCSS, Tailwind, moduł JavaScript, zasoby iOS
i Android — powstają w rdzeniu przez sklejenie napisów, bez jednego
programu spoza instalki. Panel pokazuje treść wydania, a nie samo
potwierdzenie gotowości: plik bez treści jest kopertą udaną i pustą.

## budowa/klient-poprzedni/src/moduly/multitasking/wskaznik-etapu.ts

Struktura `Subagent` nie niesie ani nazwy etapu, ani skali postępu. Niosą je pola
struktury `MonitorStatus`: `stage` z nazwą etapu bieżącego, `stageIndex` z jego numerem,
`stageCount` z liczbą etapów i `completion` ze stopniem ukończenia w procentach. Panel
wiąże telemetrię z przepływem po polu `windowId` — po tym samym, którym `Subagent`
wskazuje okno wykonawcy.

Telemetria milcząca nie jest telemetrią zerową. Gdy monitor nie ma procesu dla okna albo
odmówił odpowiedzi, wiersz nie rysuje pustego wskaźnika, ponieważ pusty wskaźnik znaczyłby
„zero etapów za sobą" — stan, którego nikt nie odczytał. Wiersz mówi wtedy wprost, że etapu
nie oddano, i skąd ta cisza pochodzi.

## budowa/klient-poprzedni/src/moduly/browser/podglad-strony.ts

Przewijanie i zaznaczenie dzieją się w kliencie, na odczytanej migawce, ponieważ to ona
jest treścią widzianą przez model. Przewinięcie strony po stronie rdzenia ma osobną
komendę `browser.scroll`, której podgląd jeszcze nie wywołuje; powód stoi przy przyciskach
nawigacji, w `formularz-nawigacji.ts`.

Treść idzie przez `textContent`, nigdy przez `innerHTML`. Pole `html` migawki jest treścią
obcą, a wstrzyknięcie go do dokumentu powłoki wpuściłoby cudzy znacznik do interfejsu
Operatora.

W trybie czytnika źródło strony znika z widoku, ale nie z migawki, więc wyłączenie trybu
pokazuje je z powrotem bez ponownego odczytu.

Wiersz wskazuje się udziałem w treści, a nie pomiarem wysokości linii: treść jest jednym
węzłem tekstowym, więc pozycji wiersza nie da się odczytać z układu bez rozbicia jej na
elementy, a to zmieniłoby drzewo dokumentu podglądu, które ma pozostać tym samym, co widzi
model.

## budowa/klient-poprzedni/src/moduly/design/tabliczka-rozmowy.ts

Postać rozmowy i okna obowiązkowe tabliczka czyta z profilu modułu prowadzonego
w `okno-komunikacji/rejestr-profilow.ts`, więc zdanie o wiązaniu zmienia się razem
z profilem, zamiast powtarzać jego treść z pamięci.

Liczby okien tabliczka nie podaje, bo profil Designu niesie w tym miejscu wartość
`GRANICA_NIEPODANA`. Pyta wyłącznie o prawo do rozmowy przez `liczbaOkienRozmowy`.

## budowa/klient-poprzedni/src/moduly/automations/okno-execution-monitor.ts

Przycisk zatrzymania nie ma warunku: nie jest wyszarzany ani przy przebiegu zakończonym,
ani przy braku odczytu — odmowa rdzenia jest widoczna w stanie błędu okna. „Na żywo” znaczy
tu ze zdarzenia, nie z odpytywania: okno słucha zdarzenia `automation.execution.status`
i przerysowuje wiersz przebiegu, gdy zdarzenie przyjdzie, także z pracy innego okna albo
innego urządzenia tego samego konta.

Przerwanie przebiegu jest natychmiastowe, a przycisk „Cofnij” stoi przez krótki czas po
akcji. Cofnięciem jest wznowienie tej samej kolejki — innej drogi kontrakt nie ma, bo typ
`QueueAction` nie zna działania odwracającego zatrzymanie. Okno mówi to wprost i nie
obiecuje przywrócenia stanu sprzed zatrzymania: wznawia kolejkę, a co z niej zostało,
rozstrzyga rdzeń.

Eksport zestawienia przebiegów istnieje w dwóch postaciach, bo służą dwóm czynnościom: opis
w Markdown czyta człowiek, a zestawienie rozdzielane średnikiem wchodzi do arkusza. Postaci
przenośnego dokumentu okno nie składa, bo wymagałaby biblioteki składu, której warstwa
kliencka platformy nie ma.

## budowa/klient-poprzedni/src/moduly/browser/panel-wyodrebnien.ts

Panel jest wynikiem pozycji „Wyodrębnij dane” paska dolnego oraz przycisku
„Wyodrębnij” paska zaznaczenia. Wyodrębnienie dzieje się w kliencie, z treści,
którą klient już ma: kontrakt nie ma komendy wyodrębniania danych ze strony,
a migawka jest jedyną treścią, którą moduł dostał od rdzenia. To z niej
wyjmowane są trzy rodzaje danych i dlatego panel nie podstawia żadnego pola,
dopóki migawka nie przyjdzie.

## budowa/klient-poprzedni/src/moduly/browser/zrodlo-automatyk.ts

Obszar `browser.*` nie niesie ani jednej komendy makra. Kontrakt trzyma powtarzalne
przebiegi w jednym miejscu — jako automatyki platformy opisane strukturą `AutomationWorkflow`.
Scenariusz nagrany przy przeglądaniu jest właśnie taką automatyką, więc okno zapisuje go tą
samą drogą, którą czyta go moduł Automations. Drugi, własny magazyn scenariuszy w module
Browser byłby drugą prawdą o tym samym bycie.

Stąd też jedno nazewnictwo dla jednego pojęcia: `AutomationWorkflow` jest w całym produkcie
automatyką, `AutomationStep` krokiem, `AutomationSchedule` harmonogramem, a `AutomationExecution`
przebiegiem. Nazwa okna Automation Studio pochodzi z opracowania modułu i pozostaje nietknięta.

## budowa/klient-poprzedni/src/modele/ustawienia-bytu.ts

Katalog kategorii, katalog definicji, budowa kontrolek, wskaźnik dziedziczenia i zapis
komendą `config.set` pochodzą w całości z modułów okna konfiguracji. Drugi generator pól
oznaczałby dwie prawdy o tej samej wartości i dwa miejsca do zmiany po dopisaniu rodzaju
wartości do kontraktu.

Wskaźnik zasięgu nanosi punkt widzenia na wybór adresu zapisu raz, przy pierwszym
odświeżeniu pola. Formularz zbudowany dla jednego modelu, a pozostawiony po przejściu na
inny, zapisywałby dalej pod adres poprzedniego, więc po zmianie osi panel dostaje pełne
pokazanie, a nie samo odświeżenie.

Pozycja katalogu, która nie dopuszcza wskazanej osi, nie znika z formularza: jej wskaźnik
zasięgu poda wtedy poziomy i osie dopuszczone przez katalog, a rozstrzygnięcie
o dopuszczalności zapisu zostaje przy rdzeniu.

Okno punktów izolacji jest jedno na klienta i pamięta swój stan, więc otwarcie stąd trafia
w ten sam egzemplarz. Pominięcie przejścia zostawiłoby w tym oknie pozycję izolacji, która
nie prowadzi donikąd.

Zapis dokonany gdzie indziej nie przerywa pracy przy polu, ponieważ zmiana stanu nanosi
wartości na pola już zbudowane i składu katalogu nie rusza.

## budowa/klient-poprzedni/src/konfiguracja/zrodlo-katalogu.ts

Klient nie zna żadnego klucza, kategorii ani etykiety pola — bierze je komendami
`settings.category.list` i `settings.definition.list`, których nazwy pochodzą ze stałych
`Command.*`. Dodanie ustawienia jest wtedy nowym wierszem katalogu, a nie zmianą kodu
interfejsu.

Odczyt jest odporny na brak katalogu: odpowiedź nieudana albo o innym kształcie daje wykaz
pusty i wywołanie `NaNiepowodzenie`, a nie odrzucenie obietnicy. Okno rozróżnia dzięki temu
katalog pusty po odmowie rdzenia od katalogu pustego na świeżo otwartym oknie.

## budowa/klient-poprzedni/src/moduly/agents/testowany-agent.ts

Czat modułu jest środowiskiem testowania wybranego eksperta i nie ma pamięci
sesyjnej: rozmowa znika przy zamknięciu okna oraz przy zmianie testowanego
eksperta.

Rozgłos jest jednym bytem po stronie klienta, a nie polem modułu. Widok modułu
staje obok sceny sesji, a nie zamiast niej, bo okno rozmowy jest oknem wiodącym
każdego modułu i zostaje na planszy. Czat modułu jest więc tym samym oknem
komunikacji ze sceny, przestawionym na ten moduł, a nie oknem wewnątrz modułu,
a wybór eksperta żyje po drugiej stronie planszy, w stanie modułu. Byt jest
wspólny, ponieważ nadawanie powstaje przy budowie widoku modułu, a odbiór przy
wiązaniu gniazda sceny z oknem rdzenia; przekazanie z rąk do rąk wymagałoby
przeciągnięcia uchwytu przez rejestr modułów i przez scenę sesji, czyli przez
dwie warstwy, których ta rzecz nie dotyczy.

Plik niczego nie czyści. Ogłasza wyłącznie zmianę testowanego eksperta i podaje
zdanie o niej, a co z tym zrobić, rozstrzyga polityka ulotności rozmowy, która
zna profil modułu i pole pamięci sesyjnej.

Ustawienie tą samą parą identyfikatora i nazwy nie jest zmianą i nie ogłasza
niczego. Rozgłos stanu modułu biegnie przy każdym odświeżeniu wykazu, przy
każdej fazie odczytu i przy każdym wchłonięciu zmiany z rdzenia; gdyby każde
z nich liczyło się jako zmiana testowanego eksperta, rozmowa znikałaby po
zapisaniu instrukcji eksperta właśnie testowanego. Sama zmiana nazwy nie jest
zmianą kontekstu i rozmowy nie kończy.

Pierwszy wybór także nie jest zmianą. Moduł wchodzi z pustym wyborem i sam
wskazuje pierwszego eksperta z biblioteki po odczycie z rdzenia; ogłoszenie
tego jako wyczyszczenia rozmowy kasowałoby zapowiedź polityki wypisaną chwilę
wcześniej i meldowało utratę wątku, którego jeszcze nie było.

Zdanie o czacie testowym stoi w module, a nie tylko w oknie rozmowy, ponieważ
przełączający eksperta w bibliotece ma przeczytać, co się stanie z rozmową,
zanim się to stanie. Zdanie o zmianie nazywa poprzedniego eksperta, bo sama
informacja o wyczyszczeniu rozmowy nie mówi, co ją wyczyściło.

## budowa/klient-poprzedni/src/konfiguracja/wiersze-obszarow-sesji.ts

Trzy pasy poboczne panelu to stan pusty, rozejście katalogu roboczego oraz pola
nieobsłużone.

Wszystkie czynności pliku są czyste: biorą stan i element, nie domykają się na niczym.
Dzięki temu wytwórnia panelu zawiera samo złożenie.

## budowa/klient-poprzedni/src/moduly/assistant/filtr-dziennika.ts

Zawężenie dziennika liczy się w oknie. Komenda `assistant.activity.list`
przyjmuje wyłącznie `windowId`, `actionId` i `limit` — tak czyta ją
`adapter_modul_asystent_czynnosci.go` — więc rdzeń nie ma czym zawęzić wykazu
po treści ani po rodzaju wpisu, a odczyt na każde naciśnięcie klawisza byłby
wywołaniem zwracającym za każdym razem to samo.

Wyszukiwanie jest pełnotekstowe i tylko takie. Wyszukiwania po znaczeniu
w dzienniku asystenta kontrakt nie ma: `knowledge.search` sięga biblioteki,
historii rozmów i plików przestrzeni roboczej, a dziennik asystenta nie jest
żadnym z tych zakresów. Brak nazywa okno wprost, zamiast podstawiać dopasowanie
po literach pod nazwę wyszukiwania semantycznego.

Porównanie idzie po zwinięciu wielkości liter właściwym dla polszczyzny,
przez `toLocaleLowerCase('pl')`, więc „Ł” znajduje „ł”.

## budowa/klient-poprzedni/src/moduly/diagnostics/zuzycie-zakladka.ts

Zakładka Usage & Cost stoi w kontenerze Observability Tools, nie w Diagnostics Center. Zestawia zużycie tokenów, żądań i kosztu wraz z raportem rozliczeniowym, opierając się na komendach odczytu podsumowania zużycia i budowy raportu, dla których rdzeń ma uchwyty.

Zakładka trzyma się trzech rozstrzygnięć. Pierwsze: jeden wymiar naraz, tak jak stanowi kontrakt — wykaz mówi, po którym wymiarze jest zebrany, ponieważ tabela z dwoma wymiarami pod jednym nagłówkiem byłaby wykazem, którego rdzeń nigdy nie oddał. Drugie: puste zestawienie nie jest brakiem — okres bez ani jednego wywołania jest poprawną odpowiedzią i zakładka nazywa to zdaniem, osobno od odmowy rdzenia i osobno od stanu, w którym o zestawienie jeszcze nie zapytano. Trzecie: koszt niepełny mówi o sobie — wywołanie kanału bez cennika nie wchodzi do kosztu, a kontrakt niesie na to osobne pole pokrycia cennikiem i pole kosztu poza cennikiem; suma podana bez tego zastrzeżenia wyglądałaby na pełny rachunek, będąc tylko jego częścią.

Raport wytwarza plik: komenda budowy raportu oddaje treść, nie zasób w magazynie, więc wytworem jest plik pobrany na dysk Operatora, nazwany okresem i postacią, żeby dwa raporty nie nadpisały się wzajemnie.

## budowa/klient-poprzedni/src/moduly/design/slowo-modelu.ts

Rdzeń odkłada bajty obrazu w magazynie pod sumą sha256 i zakłada wiersz zasobu z polami
`uri`, `format` oraz wymiarami odczytanymi z nagłówka pliku. Pole `name` jest opcjonalne
i bywa odpowiedzią modelu tekstowego, dlatego okno traktuje je jako słowo modelu, a nie
tytuł zasobu. Wynik pusty znaczy, że rdzeń nie podał nazwy w ogóle; to inny stan niż nazwa
pusta w znakach i okno musi je rozróżniać.

Generowanie zakłada wiersz dopiero po utrwaleniu bajtów w magazynie, a wgranie — po zapisaniu
pliku, więc puste pole `uri` znaczy zasób założony drogą, która tego pola nie wypełniła.

Pole `uri` jest ścieżką w systemie plików rdzenia: magazyn treści oddaje ścieżkę złożoną
z katalogu i sumy kontrolnej, i ta ścieżka wchodzi do wiersza jako odwołanie. Przeglądarka
ścieżki dyskowej nie otworzy, więc wstawiona w `<img src>` daje zdarzenie `error`, którego
przyczyną nie jest sama przeglądarka.

Komenda oddająca treść zasobu jest już w kontrakcie i obejmuje cały magazyn, a nie sam obszar
Designu — zasób Design, plik Library i dokument Studia leżą w jednym repozytorium. Brakuje jej
uchwytu w rdzeniu, więc okno nadal nie ma skąd wziąć bajtów; jest to brak obsługi, a nie brak
drogi.

## budowa/klient-poprzedni/src/moduly/multitasking/wiez-automations.ts

Sekcje Kolejki, Orkiestracja, Harmonogram i Monitor pracują na oknach modułu Automations —
kolejno Queue Manager, Orchestrator, Scheduler i Execution Monitor. Panel orkiestracji buduje
samo sterowanie i nazywa okno, na którym ono pracuje.

Powiązanie nie jest domyślne, więc sekcja bez wskazanej automatyki opisuje jego brak i sposób
założenia. Brak powiązania nie wyłącza kontrolek pracujących na samej sesji, czyli kolejek
i monitora; zawęża wyłącznie to, co adresuje pole `workflowId` — zależności i harmonogram —
ponieważ bez niego żądanie nie ma adresu.

## budowa/klient-poprzedni/src/moduly/browser/pasek-zaznaczenia.ts

Pasek niesie pięć czynności na zaznaczonym fragmencie strony: wyjaśnienie,
wyodrębnienie, notatkę, tłumaczenie i pytanie własne. Wyjaśnienie i pytanie idą
tą samą komendą i różnią się treścią pytania: pierwsze pyta zdaniem stałym,
drugie tym, co wpisano obok. Tłumaczenie wywołuje czynność strony wspólną
z paskiem dolnym, więc przekazanie fragmentu do modułu tłumaczącego jest jednym
wywołaniem, nie dwoma podobnymi. Notatka należy do panelu notatek i pasek
przenosi do niego wyłącznie fragment.

Pasek nie znika i nie gaśnie: przy pustym zaznaczeniu przyciski zostają
naciskalne i nazywają, czego brakuje, a znacznik zaznaczenia niesie tę różnicę
do arkusza stylów i do sprawdzianu.

Wyjaśnienie idzie komendą wysłania wiadomości skierowaną do okna modułu, tego
samego, którego identyfikator niosą pozostałe komendy obszaru. Osobna komenda
wyjaśniania nie miałaby w rdzeniu odbiorcy.

Zdanie po wyodrębnieniu bierze się ze skutku dopisania, a nie z samego
naciśnięcia: fragment już wyodrębniony nie wchodzi do wykazu po raz drugi,
a zdanie o dopisaniu byłoby wtedy potwierdzeniem czynności, która się nie
odbyła.

## budowa/klient-poprzedni/src/moduly/agents/katalog-narzedzi.ts

Katalog jest jedynym miejscem w kliencie, które wie, z czego składa się wyposażenie
eksperta. Źródłem jest stała `NARZEDZIA_MODELU` z `shared/contract.ts`.

Grupą jest obszar nazwy komendy, czyli człon przed pierwszym separatorem obszaru: komenda
`agent.list` należy do grupy `agent`. Własny wykaz wiążący narzędzie z grupą byłby drugą
prawdą, rozjeżdżającą się przy pierwszym narzędziu dopisanym do kontraktu. Tę samą regułę
stosuje `grupaKomendy` w `server/internal/narzedzia/grupa.go`, a porządek alfabetyczny
grup ustala tam `Grupy`.

Kod z pól `Agent.skillIds` oraz `Agent.connectorIds` czyta się dwiema drogami: jako nazwę
narzędzia, wskazującą jedną pozycję, albo jako nazwę grupy, wskazującą wszystkie pozycje
obszaru. Rozstrzyga o tym `server/internal/narzedzia/ekspert_wykaz.go`.

Narzędzia wskazane kodami idą w kolejności katalogu, a nie kodów: dwa kody wskazujące tę
samą pozycję, czyli nazwa narzędzia i jego grupa naraz, dają ją raz. Tak samo przesiewa
`przesiej` w `ekspert_wykaz.go`.

## budowa/klient-poprzedni/src/mobile/port-kolejki-decyzji.ts

Implementacja domyślna nie udaje kolejki pustej. Pusty wykaz znaczy „nic nie czeka"
i jest zdaniem uspokajającym; bez źródła byłoby ono nieprawdziwe dokładnie w tej sytuacji,
dla której ekran powstał. Dlatego `BRAK_ZRODLA_KOLEJKI` oddaje `dostepna: false` wraz
z powodem, a ekran pisze ten powód wprost — tak samo jak pulpit w `mission-control/pas-decyzji.ts`.

Wpięcie źródła to jedno wywołanie `zainstalujKolejkeDecyzji` przy montażu warstwy mobilnej.
Ekran się przez nie nie zmienia: pozycje portu mają ten sam kształt co pozycje własne,
opisany typem `PozycjaDecyzji`, więc wchodzą do tego samego wykazu i tego samego arkusza dróg.

## budowa/klient-poprzedni/src/moduly/agents/zrodlo-zakresu-eksperta.test.ts

Sprawdzian główny pilnuje jednej rzeczy: czy każda komenda dołożona rdzeniowi
ma drogę z okna. Wykaz nie jest przepisany z pamięci — powstaje z wywołań
źródła, a porównywany jest ze stałymi kontraktu, więc komenda, którą ktoś kiedyś
z okna wyjmie, zostanie tu nazwana.

Pozostałe sprawdziany dotyczą rozstrzygnięć, które warstwa kliencka podejmuje
sama i które łatwo cofnąć nieuważną poprawką. Wszystkie sprowadzają się do
jednej zasady: zbiór pusty bywa żądaniem, a pole pominięte znaczy brak zmiany.
Pomylenie tych dwóch rzeczy w jedną stronę odbiera ekspertowi wszystko,
a w drugą nie zdejmuje niczego.

## budowa/klient-poprzedni/src/moduly/assistant/etykiety-assistant.ts

Pliki budujące elementy modułu nie noszą napisów, bo napis zmienia się
z innego powodu niż układ — zmieniany w dziesięciu miejscach rozjeżdżałby
się między oknami. Wartości słowników są kluczowane stałymi kontraktu, więc
dopisanie stanu w kontrakcie przerywa kompilację tutaj, zamiast wypuścić na
ekran pusty napis.

Nazwa klasy plakietki stanu jest wykazem jawnym, nie sklejką składaną
w czasie działania: nazwa złożona doklejeniem stanu do przedrostka byłaby
poza zasięgiem kontroli klas arkusza stylów.

Poziomy zasięgu pamięci przestawiane w zakładce kontekstów są węższym
wykazem niż pełny zasięg konfiguracji platformy, bo dotyczą wyłącznie
poziomów widoczności pamięci karty sesji: globalny, projekt, karta sesji
i okno komunikacji — te, na których pamięć asystenta ma treść.

Każde zdanie stanu pustego ma dwie części: czym okno jest i czym się je
zapełnia. Zdanie mówiące wyłącznie, czego nie ma, zostawiałoby Operatora
przed oknem, o którym nie wie ani po co ono stoi, ani co miałby zrobić,
żeby coś w nim zobaczyć. Każde okno ma swoje zdanie spoczynku osobne od
zdania po odpowiedzi rdzenia, bo to dwa różne fakty: przed pytaniem
i po odpowiedzi.

Zdania wykazu zleceń spoza okna modułu nie są zdaniami stanu pustego,
ponieważ sekcja jest ukryta, dopóki nie ma czego pokazać — pustki nie ma
więc czym nazywać. Zdanie zasięgu mówi wprost, czego w wykazie nie będzie:
zleceń zamkniętych przed otwarciem modułu. Wykaz podany bez tego
zastrzeżenia wyglądałby na komplet pracy asystenta w sesji, a nim nie jest,
ponieważ kontrakt nie ma odczytu zleceń zamkniętych przed wejściem do
modułu.

## budowa/klient-poprzedni/src/aplikacja/pas-posuniec.ts

Lewa strona pasa mówi, że asystent pracuje i nad czym, czyli o zleceniach modułu
Assistant. Prawa wymienia posunięcia: nawigację, okna i polecenia. Praca modelu
docelowego, czyli odpowiedź i strumień, na pas nie wchodzi — dzieje się w oknie i tam jest
jej miejsce. Zlanie obu warstw odebrałoby rozeznanie, kto co zrobił.

Pas nie jest bramką: nie ma w nim przycisku zgody, odmowy ani wstrzymania. Jedyny
przycisk przesuwa widok wtedy, gdy podążanie zostało wstrzymane, ponieważ Operator pisał.
Posunięcie asystenta dokonało się przed pojawieniem się przycisku i nic na niego nie
czeka; pas zamienia wstrzymane przejście w jedno naciśnięcie, zamiast je zgubić.

Nazwa pasa idzie atrybutem `aria-label`, a nie nagłówkiem, ponieważ pas jest obszarem
uzupełniającym, a nie sekcją pracy. Atrybut `aria-live` w wierszach niżej niesie zmiany
bez zabierania ogniska.

Zdanie o źródle ruchu mówi najwyżej tyle, ile niesie kontrakt: przy pewności niepełnej
posunięcie pochodzi spoza tego połączenia i tak też jest nazwane. Rozstrzyga o tym
nagłówek `zrodlo-posuniec.ts`.

## budowa/klient-poprzedni/src/moduly/assistant/eksport-dziennika.ts

Eksport dziennika nie potrzebuje komendy kontraktu: wpisy są już w oknie, a plik
powstaje z tego, co Operator widzi na ekranie. Tę samą drogę ma eksport pamięci
projektu w `moduly/workspace/pamiec-pozycja.ts` oraz pomocnik `pobierzPlik`
z `modele/kontrolki-formularza-braki.ts`.

Log obejmuje wpisy po zawężeniu, a nie cały zapis rdzenia. Plik ma odpowiadać
temu, co widać: eksport szerszy niż widok kazałby zgadywać, skąd wzięły się
wiersze, których na ekranie nie było. Nagłówek pliku nazywa dlatego zarówno
chwilę pobrania, jak i liczbę wpisów, które do niego weszły.

Nagrania log nie niesie. Wpis ma wyłącznie odnośnik `audioRef`, a kontrakt nie
ma komendy pobrania dźwięku, więc odnośnik idzie do pliku wprost, bo jest tym,
co rdzeń naprawdę oddał.

## budowa/klient-poprzedni/src/moduly/automations/eksport-grafu.ts

Opracowanie wymienia cztery postacie zapisu mapy; trzy z nich powstają po stronie okna,
ponieważ treść jest już tutaj: zapis wektorowy bierze się wprost z kanwy budowanej
w `graf-krokow.ts`, a zapisy DOT i Mermaid składa ten plik.

Czwarta postać, mapa rastrowa, wymagałaby przerysowania kanwy na płótno i zostaje poza
oknem. Postać wektorowa niesie to samo, a otwiera się w każdej przeglądarce.

Żadna z tych postaci nie potrzebuje komendy kontraktu i mieć jej nie musi.

Kroki ścieżki krytycznej dostają w zapisie DOT grupę własną, a nie barwę wpisaną wprost:
barwa należy do narzędzia rysującego, a nie do treści mapy.

## budowa/klient-poprzedni/src/moduly/developer/warsztat-drogi.test.ts

Sprawdziany pilnują, że przycisk, który niczego nie woła, nie wygląda tak samo jak przycisk działający — usterka tego rodzaju widać dopiero u Operatora. Sprawdziany naciskają kontrolki okien i pytają, czy z każdej rodziny komend naprawdę wyszło wywołanie, nie o to, czy okno się narysowało.

Sprawdziany nie badają treści rdzenia — od tego są sprawdziany skutku po stronie serwera, które schodzą do bazy i na dysk. Tutaj mierzy się jedno: czy droga istnieje i czy okno mówi prawdę, gdy rdzeń odpowiada brakiem narzędzia.

## budowa/klient-poprzedni/src/aod/rozpoznanie-decyzji.ts

Pojęcia „zdarzenie wymagające decyzji" nie ma w kontrakcie ani w rdzeniu: nie
ma rodziny `decision.*`, nie ma pola `awaitingDecision`, a `ProgressStatus`
ma sześć wartości (pending, running, paused, stopped, done, failed) i żadna
z nich nie znaczy „czekam na Ciebie". Regułę wywodzi więc nakładka i tak też
jest nazwana Operatorowi w oknie (`sekcja-decyzji.ts`).

Sedno reguły: budzi to, przy czym proces stoi — `running` i `done` nie budzą
nigdy. Powody dzielą się na pewne (stan rdzenia mówi wprost, że proces stoi)
i sporne (to ocena nakładki); sporne nie są ukrywane, tylko oznaczone własną
wagą, żeby Operator wiedział, czyja to ocena. `bieg-stanal` (pewna) to
`loop.stopped = true`: koordynator sam ogłosił, że bieg naprawczy stanął,
i podał powód — sprawdzany pierwszy, bo jest najbardziej szczegółowy, niesie
okno koordynatora i przyczynę. `wstrzymany` i `zatrzymany` (pewne) to `status
= paused` i `status = stopped`. `usterka` (sporna) to `status = failed` —
nakładka budzi, bo proces stoi, a nikt inny go nie ruszy. `bez-ruchu`
(sporna) to `status = pending` dłużej niż próg — pozycja, która nigdy nie
ruszyła, jest nieodróżnialna od pozycji zapomnianej.

Źródło sygnału jest częścią odpowiedzi: odczyt nadrabiający i zdarzenie na
żywo mają różną świeżość. Pole `loop` w `progress.changed` jest martwe —
kontrakt je obiecuje, `telemetria_proces.go` nigdy go nie wypełnia. Stan
biegu naprawczego dojeżdża wyłącznie zdarzeniem `window.state.changed`,
którego treść jest odpisem `WindowStateGetResponse`, i stamtąd ta reguła
bierze `loop`.

Próg braku ruchu dla stanu `pending` wynosi 15 minut i pochodzi z progów AOD:
przekroczenie tworzy sugestię klasy „stan kolejki zadań". Pozycja `pending`
to właśnie zadanie oczekujące w kolejce, więc reguła nakładki bierze próg
stamtąd, zamiast stanowić własny. Kontrakt progu nie niesie; `LoopState` ma
własny `threshold`, ale liczy obiegi, nie czas, i dotyczy wyłącznie biegu
naprawczego — próg jest argumentem reguły, więc jego zmiana dotyka jednej
stałej, a nie kodu.

Rodzaj sugestii przypisany każdemu powodowi: cztery powody mówiące, że
proces stoi, są wskazaniem problemu („zadanie w stanie błędu, kolejka
zatrzymana, pętla przerwana"); powód `bez-ruchu` jest kolejnym krokiem
(klasa „stan kolejki zadań": zadanie oczekujące dłużej niż próg czasu) —
nic się nie zepsuło, czeka na ruszenie.

Waga ujawnienia każdego powodu: wysoka oznacza „punkt decyzyjny
wstrzymujący proces, kolejka zatrzymana", średnia oznacza „przekroczony
czas oczekiwania zadania" — i to jest dokładnie `bez-ruchu`. Wyjątek wagi
krytycznej: taka sugestia ujawnia się mimo wyciszenia — plakietką, bez
dymka. Powody będące punktem decyzyjnym wstrzymującym proces (bieg
naprawczy, który stanął, i proces wstrzymany) zatrzymują pracę i czekają
wprost na rozstrzygnięcie człowieka; zatrzymanie, usterka i brak ruchu
przez wyciszenie przeczekają.

## budowa/klient-poprzedni/src/moduly/automations/okno-workflow-builder.ts

Cofnij i ponów działają miejscowo, na pracy sprzed zapisu. Po zapisie
śladem jest numer wersji definicji, a powrót do wersji wcześniejszej
należy do pozycji „Historia wersji" — ta nazywa swoją komendę i mówi, czy
rdzeń ma dla niej uchwyt. Zdanie potwierdzenia powstaje z odpowiedzi
rdzenia — z identyfikatora, numeru wersji i liczby kroków; przy
duplikowaniu okno porównuje identyfikator sprzed czynności z tym, który
wrócił, bo dopiero różnica dowodzi, że powstała nowa definicja.

Walidacja definicji idzie przy każdej zmianie kroków, nie dopiero przy
zapisie: zastrzeżenie postawione w chwili wpisywania jest poprawką, a to
samo zastrzeżenie postawione po zapisie jest już tylko wiadomością
o definicji wadliwej, która zdążyła trafić do magazynu automatyk.
Sygnalizacja jest podwójna, bo służy dwóm czynnościom: znacznik przy kroku
pokazuje, który krok poprawić, a wykaz w treści okna mówi, co dokładnie —
funkcja oddaje zastrzeżenia, żeby zapis mógł dopowiedzieć o nich
w potwierdzeniu.

Zapis idzie mimo zastrzeżeń — walidacja nie jest bramą, tak samo jak
walidacja układu po stronie rdzenia. Walidacja na żądanie sprawdza
zastrzeżenia dotyczące kroków po stronie okna; cykle i ścieżkę krytyczną
całego układu rozstrzyga rdzeń w Orchestratorze, więc zdanie potwierdzenia
mówi wprost, gdzie szukać drugiej połowy oceny.

## budowa/klient-poprzedni/src/aod/tryb-obecnosci.ts

Plik trzyma tryb obecności funkcji Always On Display, sięga po wykaz wyciszeń i rozstrzyga z obu trzy rzeczy, o które pyta warstwa widoku: czy awatar jest w polu widzenia, bo tryb ukryty go zabiera; czy plakietka liczbowa się pokazuje, bo wyciszenie czasowe ją chowa; oraz czy dana sugestia otwiera dymek w tej chwili, biorąc razem wagę, tryb, trzy rodzaje wyciszenia, limit godzinowy i odstęp między dymkami.

Wyjątek wagi krytycznej jest zaszyty w regule, a nie zostawiony wołającemu, i przechodzi przez wszystkie rodzaje wyciszenia: czasowe, kontekstowe, klasy zdarzeń oraz tryb cichy. Punkt decyzyjny pętli wykonawczej wstrzymujący proces ujawnia się mimo wyciszenia plakietką, bez dymka, a tryb cichy zachowuje ten wyjątek bez syntezy mowy. Osłabienie tego wyjątku byłoby jedyną ciszą, po której Operator nie dowiaduje się, że praca stoi, więc reguła sprawdza go pierwsza, przed wszystkim innym.

Stan obecności i wyciszeń stoi na stanowisku Operatora, a magazyn zapisu jest podawany parametrem: kontrakt nie ma jeszcze kategorii ustawień ani komendy zapisującej tryb obecności czy wyciszenie nakładki jako byt rdzenia. Gdy taka komenda powstanie, magazyn zmieni się w jednym miejscu, bez zmiany ani jednego wołacza.

## budowa/klient-poprzedni/src/moduly/browser/zrodlo-browser.ts

Zebrane źródła i notatki mieszkają w `stan-przegladania.ts`, żeby trzy okna modułu Browser
patrzyły na jeden wspólny zbiór, a nie na trzy osobne kopie. Komendy `source.list`
i `note.list` czytają dokładnie to, co moduł zapisał komendami `browser.source.add`
i `browser.note.add`, i nic ponad to.

Odmowa wykazu ma dwa znaczenia, które okno rozróżnia: kod `not_found` znaczy, że rdzeń nie
zna tego okna przeglądania, a nie że nic nie zebrano — wykaz pusty przychodzi ze statusem
udanym i pustą tablicą. Żadne wywołanie tego źródła nie rzuca wyjątkiem i nie odrzuca
obietnicy: niepowodzenie wraca w polu `blad` wyniku, a nazwy komend biorą się wyłącznie ze
stałych kontraktu.

Rodziny komend dołożone ponad migawkę, źródło i notatkę — otwieranie i zarządzanie kartami,
przestrzenie robocze, monitory, kanały, listy do przeczytania, zakładki i narzędzia
deweloperskie przeglądarki — są każda jedną komendą kontraktu bez własnego stanu warstwy:
sprawdzają wyłącznie kształt odpowiedzi rdzenia.

## budowa/klient-poprzedni/src/moduly/automations/panel-zlecen.ts

Okno Queue Managera pokazywało dotąd kolejkę jako całość i mówiło wprost, że wykazu zleceń rdzeń nie oddaje. Ten panel jest odpowiedzią na to zdanie: odczyt wykazu zleceń oddaje zlecenia wraz z ładunkiem, próbami i terminem, odczyt głębokości kolejki oddaje obciążenie w czasie, a jedenaście pozostałych czynności posuwa pojedyncze zlecenie.

Kolejka bierze się ze stanu modułu, a pole jej wskazania zostaje: Operator bywa w Queue Managerze przy kolejce innej niż bieżąca automatyka, na przykład kierując zlecenie do kolejki przeglądu ręcznego.

## budowa/klient-poprzedni/src/moduly/library/higiena-repozytorium.ts

Duplikat dokładny poznaje się po sumie kontrolnej — tej samej, którą rdzeń
wyliczył przy wgraniu — więc raport nie jest domysłem okna, tylko odczytem
jego odpowiedzi. Plik osierocony to plik bez etykiety i bez kolekcji: oba
pola niesie `LibraryFile`, więc reguła audytu ma na czym stanąć. Cztery
czynności, których okno zrobić nie może samo, prowadzą do rdzenia: rozpoznanie
duplikatów wraz z niemal-duplikatami po treści i po obrazie
(`library.duplicate.scan`, porównanie całego zbioru z całym), weryfikacja
integralności (`library.fixity.check`, przeliczenie sumy kontrolnej z bajtów
leżących pod odwołaniem), normalizacja nazw (`library.name.normalize`) wraz
z przebiegiem próbnym, oraz dziennik audytu (`library.audit.list`) i pulpit
stanu (`library.stats.get`) liczony po całym zbiorze, nie po odczytanej
stronie wykazu.

Raporty liczone z wykazu zostają obok tych z rdzenia, bo odpowiadają
natychmiast i bez ruchu do rdzenia — ale to pulpit z rdzenia jest miarą
repozytorium, a wykaz w pamięci okna jest próbką. Przeliczenie wskaźnika
znaczenia idzie komendą `knowledge.index` w zakresie biblioteki — jedyną
czynnością konserwacyjną repozytorium, którą kontrakt niesie. Wskaźnik nie
odświeża się przy wgraniu pliku, więc bez tej czynności wyszukiwanie po
znaczeniu opisuje stan sprzed ostatniego napływu.

Wskaźnik pomija pliki bez treści i pliki, których treści nie da się odczytać
jako tekstu; zero wprowadzonych pozycji przy niepustym wykazie zdanie nazywa
wprost, zamiast milczeć o różnicy.

Karta raportu: zawężenie Library Explorera jest odwracalne jednym przyciskiem
w pasku narzędzi Explorera i znika przy każdym nowym odczycie wykazu. Karta
zgodna (bez pozycji do zgłoszenia) zachowuje przycisk i po naciśnięciu mówi,
że zawężać nie ma do czego — wygaszenie przycisku robiłoby bramę tam, gdzie
jest wynik pomiaru.

## budowa/klient-poprzedni/src/moduly/diagnostics/okno-diagnostics-center.ts

Diagnostics Center jest oknem wiodącym modułu Diagnostics: przegląd zagregowanego stanu systemu i uruchomienie analizy, punkt wejścia agregujący Logs Viewer oraz Errors Panel. Uruchomienie analizy jest jedynym przyciskiem sprawczym widoku. Po udanym biegu okno zapisuje identyfikator analizy do stanu, a panel rekomendacji pyta o rekomendacje tej analizy po jej identyfikatorze. Zakres czasu idzie do stanu modułu, bo Logs Viewer i Errors Panel jadą tym samym zakresem. Druga subskrypcja zdarzenia zmiany analizy tu nie stoi: stan diagnostyki sam nasłuchuje zdarzenia i filtruje po analizie bieżącej, oknu wystarcza jego własny nasłuch zmiany stanu.

Skróty zakresu czasu obejmują ostatnią godzinę, dzień, tydzień oraz zakres własny. Pozycja bez okna znosi zawężenie: zakres pusty jest stanem poprawnym modułu, nie brakiem. Skrót nie liczy niczego, czego nie widać: wyliczoną chwilę wpisuje w pola początku i końca, więc Operator czyta z okna dokładnie te liczby, które idą do rdzenia.

Rdzeń nie honoruje dziś dwóch pól żądania analizy: uruchomienie analizy dobiera błędy wyłącznie zakresem czasu, a pole porównania migawki nie ustawia się nigdy, więc identyfikator analizy porównywanej nie pojawi się w odpowiedzi nawet wtedy, gdy porównywana analiza istnieje. Potwierdzenie zawsze udane robiłoby z tych dwóch pól bez skutku pola pozornie działające. Zdanie potwierdzenia mówi więc liczbami z odpowiedzi, a każde pominięte pole żądania wychodzi na wierzch tonem nieudanym. Gdy rdzeń zacznie te pola honorować, zastrzeżenia znikną same, bo nic o zachowaniu rdzenia nie jest tu wpisane na sztywno.

## budowa/klient-poprzedni/src/moduly/developer/okno-git-panel.ts

Czynności typu `GitActionKind` idą przez komendę `developer.git.action`; świeżo otwarty
panel nie wie o repozytorium nic i mówi to wprost stanem pustym, zamiast udawać czyste
repozytorium. Okno stoi we wspólnym stanie modułu: nasłuchuje `stan.naZmiane` i podstawia
wskazaną ścieżkę do pola „Ścieżki”, dopóki Operator nie wpisze własnej, a ścieżki z własnego
wyniku oddaje z powrotem przez `stan.wskazPlik`, czym otwiera je w Code Editorze.

Zdanie stanu pustego, o tym czego kontrakt nie niesie, składa `katalog-komend.ts` z rejestru
komend wziętego z odpowiedzi `connection.hello`, a nie stała zapisana w tym pliku: napis na
stałe zostałby na scenie jako nieprawda w dniu dołożenia brakującej komendy.

Etykieta przycisku wysłania wymuszonego nie nazywa polecenia gita wprost: pole kontraktu to
`force`, a czym rdzeń je wykonuje, mówi pole `output` wyniku — napis o poleceniu rozjechałby
się przy zmianie rdzenia, a wyjście polecenia jest zawsze świeże.

Potwierdzenie zatwierdzenia bez opisu używa `potwierdzenie(zdanie, false)`, nie
`blad(...)`, bo to nie jest odmowa rdzenia ani nieudany odczyt, tylko lokalna walidacja przed
wysłaniem — treść okna z wynikiem ostatniej czynności ma zostać nietknięta, a `potwierdzenie`
właśnie tego nie kasuje.

## budowa/klient-poprzedni/src/moduly/developer/zrodlo-warsztatu.ts

Interfejs `ZrodloDeveloper` obsługuje okna wiodące modułu Developer — edytor, drzewo,
repozytorium i budowanie — a ten warsztat obsługuje rodziny komend należące do zakładek
Dev Tools, panelu Run & Debug i paska operacji kontekstowych. Rozdział między obu źródłami
idzie po odbiorcy komend, nie po wielkości pliku: Project Tree nie ma nic wspólnego
z kontenerami, a jedna umowa na wszystko kazałaby każdemu oknu przyjmować zależność od
czterdziestu metod, z których używa czterech.

Każda czynność tego źródła oddaje `Wynik`, nie samą treść odpowiedzi: okna mają obowiązkowy
stan błędu, więc źródło nie połyka odmowy i nie zwraca w jej miejsce pustego wykazu. Pusty
wykaz kontenerów i odmowa odczytu to dwa różne zdania, a Operator ma prawo wiedzieć, które
z nich obowiązuje — stąd brak zastępowania odmowy pustą tablicą w całym pliku.

Każda odpowiedź przechodzi przez `sprawdzKsztalt`: rdzeń rozminięty z kontraktem nie ma
prawa dojść do okna jako wartość niezdefiniowana w środku rysowania widoku, bo to jest ta
klasa usterki, która ujawnia się dopiero u Operatora.

Komendy tego źródła grupują się tematycznie w warstwę językową Code Editora, historię pliku
edytora, odczyt okna Build Output, Run & Debug, API Client, Data Console, Containers,
zależności i bezpieczeństwo oraz operacje kontekstowe i sondę warsztatu.

## budowa/klient-poprzedni/src/moduly/browser/okno-notes-panel.ts

Notes Panel niesie jedną odpowiedzialność: złożenie okna i wykaz notatek. Formularz stoi w osobnym pliku formularza notatki, rozmowa z rdzeniem w osobnym pliku czynności notatek, jeden wiersz wykazu w osobnym pliku wiersza notatki.

Edycja nie udaje zmiany: komenda aktualizacji notatki stoi w kontrakcie, ale to okno jeszcze jej nie wywołuje. Pozycja „Edytuj” wypełnia formularz treścią notatki i mówi wprost, powodem wziętym z odczytu wykazu komend rdzenia, że zapis utworzy notatkę nową, nie zmieni zastanej.

Powód, dla którego zapis zakłada notatkę nową zamiast zmieniać zastaną, bierze się z odczytu wykazu komend rdzenia, nie z napisu w module: komenda aktualizacji notatki stoi w kontrakcie, więc od dnia, w którym rdzeń dostanie jej uchwyt, powód ma brzmieć inaczej i zabrzmi bez wchodzenia w ten plik.

Trzy stany obowiązkowe wykazu notatek stosują ten sam zestaw reguł, co w panelu źródeł, bo oba wykazy przychodzą jednym zaciągnięciem. Kolejność pytań: czekanie, odmowa, pustka. Odmowa wykazu jest błędem panelu tylko przy pustym wykazie — z notatkami na ekranie wpisy zostają, a powód idzie zdaniem przy wykazie.

## budowa/klient-poprzedni/src/moduly/automations/okno-orchestrator.ts

Walidacja nie jest bramą — rdzeń zapisuje układ także wtedy, gdy ma cykl,
i oddaje zastrzeżenia, które okno pokazuje wprost. Zastrzeżenie do układu
nie jest niepowodzeniem czynności: werdykt o czynności należy do odpowiedzi
na wywołanie, werdykt o układzie stoi osobno — w ocenie układu i w wykazie
zastrzeżeń. Pole `criticalPathStepIds` liczy rdzeń, który zna cały układ;
drugi rachunek w oknie byłby drugą prawdą o tym samym grafie. Układ
pokazują trzy widoki tej samej treści: kanwa grafu (`graf-krokow.ts`),
wykaz zależności i model stanów przebiegu (`maszyna-stanow.ts`) — kanwa
mówi, jak długi jest łańcuch i gdzie tory się rozchodzą, wykaz co z czym
jest związane i którą zależność usunąć, model co przebiegowi wolno dalej.

Panel zależności stoi w narzędziach kontekstowych, bo zawęża i bada ten
sam układ, którym gospodaruje panel akcji okna. Panel dopełnień układu
(bramka dołączenia, grupa kroków, krok wycofujący i spięcie kolejek
z silnikiem środowiska MultitaskingAI) stoi obok panelu zależności, bo
gospodaruje tym samym układem — innym jego wymiarem.

Usunięcie zależności komendą `orchestration.dependency.remove` oddaje pole
`removed` oraz `dependencies` po usunięciu; pełny obraz układu — ocenę
i ścieżkę krytyczną — odczytuje się po niej ponownym
`automation.orchestrator.define` bez zmiany, bo komenda usunięcia oddaje
same zależności. Nazwa pliku eksportu bierze identyfikator automatyki;
automatyka niewskazana daje nazwę rodzajową, bo plik ma się zapisać mimo
wszystko.

## budowa/klient-poprzedni/src/aod/wyciszenie-menu.ts

Dwie kopie tego samego menu (powierzchnia interakcji awatara — wyciszenie od
ręki bez otwierania kolumny — i nagłówek powierzchni interakcji) czytają
ten sam stan i rysują się jego powiadomieniem, więc nigdy nie mówią dwóch
różnych rzeczy. Żadna pozycja nie pyta o potwierdzenie i żadna nie jest
wyszarzana — reguła przyjęta w `cztery-stery.ts`. Pozycja, której nakładka
dziś nie ma czym wykonać (na przykład wyciszenie modułu, którego rdzeń nie
wskazał), zostaje klikalna i mówi, czego brakuje oraz po czyjej stronie.

Menu jest osobnym wyzwalaczem wobec awatara: kliknięcie pojedyncze i podwójne
awatara zostają nietknięte, a znak „⋮” jest zwykłym przyciskiem — osiągalnym
klawiszem tabulacji, otwieranym Enter i Spacją, zamykanym Esc, z ruchem
strzałkami po pozycjach. Menu jest dostępne z klawiatury bez skrótu
własnego — znak stoi w kolejności tabulacji obok awatara, więc żaden własny
skrót nie jest potrzebny.

Kontekst nakładki (bieżący moduł i bieżąca karta sesji) bywa pominięty —
wtedy pozycje kontekstowe zostają widoczne i mówią wprost, że menu nie zna
bieżącego bytu, zamiast zniknąć.

## budowa/klient-poprzedni/src/moduly/research/akcje-okien.ts

Droga wykonania akcji stoi w tym katalogu, a nie w oknie, bo rozstrzyga o odbiorze: akcja
bez wykonawcy idzie generycznym `window.action`, zamiast znikać z paska albo udawać
wykonanie. Nazwa komendy nigdy nie jest tu napisem wpisanym wprost: kod akcji, która ma
w kontrakcie własną komendę, bierze się ze stałej `Command.*` — inaczej zmiana nazwy
w `contract.json` zostawiłaby w panelu martwe wywołanie, którego kompilator nie wychwyci.
Zaporę na to niesie plik `src/kontrakt.test.ts`.

Wykaz tego katalogu nie jest kopią rejestru akcji rdzenia. Rejestr, czytany komendą
`action.list`, opisuje akcje modułu zapisane w bazie, a zaczyn katalogu akcji (migracja
`migracja_009_zaczyn_akcji.sql`) zakłada dla modułu Research wyłącznie trzy pozycje paska
promptu — pozycje wiadomości — bo powstają złączeniem z tabelą `modul`. Pozostałe kody
z tego pliku wiersza w katalogu rdzenia nie mają, więc panel zbudowany wyłącznie z rejestru
byłby pusty; wykaz w tym pliku znika dopiero wtedy, gdy rejestr odda te same pozycje.

Zdolność drogi `komenda` niesie kod akcji wzięty ze stałej `Command.*`, a wywołanie składa
`wywolania-komend.ts`. Zdolność drogi `akcja` jest nastawą pola komendy istniejącej albo
należy do okna konfiguracji; rdzeń tę komendę zna, więc odmowa jest merytoryczna — kod
`not_found` z katalogu akcji albo `conflict` braku wykonawcy z adaptera okna akcji rdzenia.

Dla Discovery Panel i Reading View odmowa dotycząca okna spoza katalogu przychodzi o krok
wcześniej niż przy akcji bez wiersza w katalogu akcji: Discovery Panel i Reading View mają
wiersz w opracowaniu modułu, a nie mają go jeszcze w migracjach katalogu okien rdzenia, więc
akcja ma to zapowiadać, zamiast obiecywać drogę, której dziś nie ma nawet do połowy.

## budowa/klient-poprzedni/src/moduly/diagnostics/okno-errors-panel.ts

Stan pusty i stan odmowy odczytu są w oknie rozróżnione celowo: Errors Panel
jest jedynym miejscem w produkcie, w którym widać odrzucenie komendy przez
rdzeń, więc zrównanie tego stanu z pustym wykazem ukryłoby jedyny ślad
odmowy dostępny dla Operatora.

Kod odmowy nie leży w polu message błędu, tylko osobno. Pole errorCode
niesie go w kontrakcie docelowym, a dzisiejszy rdzeń wkłada go jeszcze do
context.errorCode; po tym kodzie rozpoznaje się rodzinę odmów zakończoną na
unknown. Odczyt sięga do obu miejsc w tej kolejności, żeby wypełnienie pola
kontraktowego przez rdzeń nie wymagało zmiany w tym pliku.

Zakres czasu wykazu bierze się wyłącznie ze wspólnego stanu modułu przez
subskrypcję zmiany tego stanu; okno nie prowadzi drugiego, własnego zakresu,
bo dwa niezależne zakresy w jednym module rozjeżdżałyby się przy każdej
zmianie jednego z nich.

Subskrypcja zakresu jest wyłączana na czas zapisu migawki analizy: zapis
migawki budzi ten sam nasłuch, który go wywołał, więc bez wyłączenia okno
odczytywałoby wykaz dwa razy i migało stanem ładowania obok świeżo
pokazanego potwierdzenia.

Odczyt wykazu błędów należy wyłącznie do złożenia modułu, nie do wytwórni
okna: wytwórnia woła odswiez() sama, tak jak pozostałe okna Diagnostics,
a odczyt już w wytwórni obok odczytu ze złożenia wysyłałby dwa żądania
error.list na jedno zmontowanie modułu.

Pasek akcji niesie trzy pozycje bez pokrycia w kontrakcie: diagnostics.error.list
jest wyłącznie odczytem, a DiagnosticError niesie pola status, priority i note
bez komendy zapisu. Panel może po tych polach filtrować, nigdy ich nadawać.
Powód każdej nieczynnej pozycji składa się z pól kontraktu, a nie z gotowego
napisu, żeby zdanie o braku komendy zmieniło się samo w dniu, w którym
komenda się pojawi.

Zdanie potwierdzenia przekazania do analizy mówi, co objęła analiza, a nie
co wysłało okno: rdzeń dobiera błędy uruchamianej analizy wyłącznie zakresem
czasu i pola errorIds żądania nie używa, więc liczba wzięta z żądania
przeczyłaby temu, co w tej samej chwili pokazuje Diagnostics Center. Rozjazd
między wykazem wysłanym a objętym jest w zdaniu powiedziany wprost i tonem
nieudanym.

Wykaz błędów niesie liczbę wszystkich pozycji spełniających warunki osobno
od samych pozycji: rdzeń liczy total osobnym zapytaniem, biorącym stan,
priorytet i zakres czasu, a nie biorącym granicy, więc wykaz bywa ucięty
granicą z pola okna albo granicą domyślną rdzenia. Rozjazd między total
a liczbą oddanych pozycji jest widoczny wprost nad wykazem — ucięcie bez
ostrzeżenia byłoby tej samej rodziny co pusty wykaz przy odmowie odczytu.

Kod odmowy w polu errorCode ma pierwszeństwo przed kontekstem zapisu, bo
odmowa komendy jest faktem o samym błędzie, a nie o okolicznościach jego
zapisu.

## budowa/klient-poprzedni/src/moduly/automations/panel-zaleznosci.ts
Panel należy do rodziny komend `orchestration.*`: okno Orchestratora zapisuje cały układ komendą
`automation.orchestrator.define`, natomiast ta rodzina działa na tym samym grafie w sposób bardziej
szczegółowy — odczyt bez zapisu (`dependency.list`), zmiana pojedynczej krawędzi (`dependency.set`)
i sprawdzenie układu (`validate`). Zależności są panelem, a nie osobnym oknem, ponieważ rejestr okien
operacyjnych rdzenia nie przewiduje dla nich wpisu; panel wchodzi w pas narzędzi kontekstowych
Orchestratora. Rdzeń przyjmuje również krawędź tworzącą cykl i zwraca zastrzeżenia w wyniku, więc
niepowodzenie czynności zapisu jest odrębne od oceny poprawności układu — stąd rozdzielenie komunikatu
o powodzeniu czynności od zdania opisującego zastrzeżenia układu.

## budowa/klient-poprzedni/src/moduly/agents/licznik-narzedzi.ts
Nadmiar narzędzi degraduje wywołanie po cichu: docierają do modelu bez opisów, więc
model po nie nie sięga, a nic tego nie zgłasza. Liczba pochodzi z jednego wyliczenia
(policz) obsadzonego w trzech oknach modułu — Agent Builder, Skills Manager, Connectors
Manager — więc przypisanie umiejętności w jednym oknie przestawia liczbę w pozostałych.
Wtyczki nie wchodzą do liczby: server/internal/narzedzia/ekspert_definicja.go bierze
do doboru narzędzi wyłącznie Agent.SkillIds i Agent.ConnectorIds, a wtyczka jedzie
katalogiem rozszerzeń powłoki (--plugin-dir), nie wykazem tools/list. Serwer narzędzi
przesiewa nie wykaz kontraktu, lecz wykaz okna — kontrakt powiększony o pozycje
dokładane przez rolę okna (WykazZasiegu w zasieg_roli.go). Licznik zna sam kontrakt,
bo klient nie wie, w jakiej roli okno eksperta zostanie otwarte: kod wskazujący
narzędzie roli zostanie policzony jako nierozpoznany, choć rdzeń go rozpozna. Liczba
jest dolnym oszacowaniem dla okien roli asystenta i dokładna dla okien roboczych.
Progu znaczeniowego nie ma i licznik go nie udaje — nie maluje pasma zielony/bursztyn/
czerwony po zmyślonych wartościach. Bursztyn zapala się wyłącznie przy stanie
wyprowadzonym ze złożenia wykazu eksperta: gdy zawężenie nie weszło i model dostanie
wykaz w całości.

Funkcja policz rozstrzyga trzy przypadki, każdy osobnym zdaniem, bo znaczą co innego:
kodów nie ma — zawężenia nie ma czym wykonać, idzie wykaz w całości; kody są, ale żaden
nie nazywa narzędzia ani grupy — jak wyżej, tyle że z winy kodów, więc zdanie wymienia
je z nazwy; rozpoznano co najmniej jeden kod — wykaz zawężony, liczba jest doborem.

Widok licznika: stan nigdy nie jest samym kolorem, bursztyn niesie zdanie i znacznik
stanu, bo żeton barwy nie zwalnia komponentu z etykiety tekstowej. Atrybut aria-live
ogłasza zmianę liczby, bo liczba zmienia się wskutek czynności wykonanej w innym oknie
modułu.

Zdanie uwagi łączy trzy rzeczy naraz, bez wchodzenia w drugie okno: powód braku
zawężenia, gdy jest, kody nierozpoznane, gdy są, i wtyczki, gdy są — te ostatnie
z zaznaczeniem, że jadą inną drogą niż wykaz narzędzi. Na końcu zawsze zastrzeżenie,
że licznik nie orzeka, czy liczba jest bezpieczna.

## budowa/klient-poprzedni/src/moduly/automations/przyjecie-przekazania.ts
Rodzina komend `automation.*` jest jedynym magazynem scenariuszy w platformie — moduły Browser,
Assistant i Terminal nie mają własnego magazynu i oddają scenariusze tutaj. Do magazynu prowadzą
dwie drogi: zapis wprost komendą `automation.workflow.save`, którym idzie rutyna Assistanta oraz
scenariusz zapisany z okna Browsera, trafiający do wykazu automatyk od razu; oraz przeniesienie
kompletu komendą `context.transfer`, którym Browser oddaje scenariusz do modułu docelowego. Rdzeń
zakłada wtedy okno modułu Automations, a przeniesiony komplet trafia do konfiguracji sesji tego
okna, do obszaru kontekstu rozmowy, pod polem `transferredContext`. Ten plik obsługuje drugą drogę
od strony odczytu: wykaz okien zawężony do modułu, a dla każdego okna konfiguracja obowiązująca
zawężona do obszaru kontekstu rozmowy. Plik niczego nie zapisuje i niczego nie przyjmuje sam —
przeniesiony komplet jest propozycją, a zapis do magazynu automatyk pozostaje osobną, jawną
czynnością operatora w oknie wykazu. Pole `executionParams` kontrakt opisuje jako `json`, więc
kształt sprawdzany jest jawnie zamiast rzutowany: przekazanie z modułu, który ułoży komplet inaczej,
zostaje pominięte zamiast trafić do okna jako scenariusz bez kroków. Okna zamknięte są pomijane przy
odczycie, ponieważ przekazanie do okna zamkniętego jest przekazaniem odbytym i zakończonym.

## budowa/klient-poprzedni/src/aktualizacja/baner-aktualizacji.ts
Pas nie zadaje pytania potwierdzającego przed założeniem wydania: kliknięcie jest
zgodą, a jedynym miejscem, gdzie produkt o cokolwiek pyta, zostaje logowanie. Pas
pojawia się wyłącznie wtedy, gdy jest co zakładać — brak sieci, brak wydań
i wydanie nie nowsze od zainstalowanego znaczą to samo: baneru nie ma. Sam
z siebie nic nie aktualizuje, dopóki nikt nie kliknie.

Odmowa zostaje na widoku: gdy powłoka nie zdoła założyć wydania, baner nie znika,
tylko zamienia się w zdanie mówiące, co poszło nie tak i co z tym zrobić. Zdanie
zachęty rozstrzyga powłoka, a nie samo rozpoznanie środowiska, ponieważ na kopii
z pakietu instalacyjnego powłoka odmawia podmiany zaraz po kliknięciu — pas
obiecuje wtedy tylko to, co wiadomo na pewno.

Tor postępu rysuje się wyłącznie przy znanej całości pobrania. Gdy powłoka poda
same bajty, jest sam licznik megabajtów bez toru; gdy nie poda nic, pas pokazuje
upływ czasu od kliknięcia jako jedyną liczbę, którą zna, zamiast rysować pasek
udający procenty.

Obudowa trzyma pas i przycisk wykazu jako rodzeństwo, nie zagnieżdżenie:
zagnieżdżenie przycisku wykazu w banerze byłoby niepoprawnym znacznikiem
i pułapką dla czytnika ekranu, bo jedno kliknięcie trafiałoby w dwa sterowniki
naraz.

Odmowa trwała wyłącza pas na stałe, ponieważ powtórzone kliknięcie dałoby tę samą
odmowę co do słowa. Odmowa przemijająca (brak łączności, przerwane pobieranie)
pozostawia pas klikalnym, bo ponowienie jest jedyną sensowną czynnością. Droga
niedostępna rozpoznana przed pobraniem dostaje inny znak niż odmowa, ponieważ
nic się jeszcze nie zaczęło; powłoka z góry mówi, że tej kopii nie podmieni, więc
żaden bajt nie idzie po nic.

Obieg pytania o wydania pomija czynność trwającą albo już zakończoną, żeby nie
nadpisać widoku pracy albo wyniku zachętą. Drugie sprawdzenie stanu po czekaniu
na odpowiedź o drodze chroni przed tym samym: między pytaniem o wykaz a pytaniem
o drogę mija czas sieciowy, w trakcie którego zakładanie mogło już ruszyć.

Zapora wejścia w zakładaniu wydania dopuszcza jedną aktualizację naraz, bo druga
pisałaby w ten sam plik roboczy obok aplikacji; powłoka ma tę samą zaporę
u siebie, ale rozstrzygnięcie po stronie pasa oszczędza migającego stanu pracy,
po którym natychmiast przychodzi odmowa. Nasłuch postępu działa wyłącznie na
czas jednego przebiegu pobierania, a zdarzenie spóźnione albo dotyczące już
zdjętego banera nie ma prawa go wskrzesić.

Zdanie powodzenia po założeniu wydania układa powłoka, bo to ona wie, co się
stało z plikiem i z rdzeniem w tle; powłoka odkłada też sam restart, więc jest
chwila na pokazanie tego zdania. Chronologia wydań rozwija się dopiero na
żądanie, żeby nie pytać kanału bez potrzeby, a zwinięcie jej nie kasuje stanu.

Poza powłoką natywną baner nie ma czynności do wykonania po kliknięciu, bo
podmiana pliku i restart są własnością powłoki — pytanie kanału wydań byłoby
wtedy ruchem w sieć po nic.

## budowa/klient-poprzedni/src/moduly/agents/modul-agents.ts
Układ wynika z roli okna. Agent Builder jest kreatorem i punktem wejścia modułu, więc
stoi w pasie pierwszym na całą szerokość. Model Configuration jest oknem pomocniczym,
a cztery pozostałe są zarządcami, więc stoją w pasie drugim obok siebie. Trzy z nich
odnoszą się do eksperta wybranego w kreatorze; Katalog rozszerzeń stoi na końcu pasa,
bo jako jedyny mówi o platformie, nie o ekspercie, i dlatego nie gaśnie, gdy żaden
ekspert nie jest wybrany. Jeden ekspert jest czynny na cały moduł: stan agentów jest
jeden, więc wybór w bibliotece przestawia okna eksperta naraz. Dlatego zakładka paska
edytora nie otwiera drugiego formularza — przenosi ognisko do okna, które daną rzeczą
zarządza.

Kod zakładki nie jest kodem okna: zakładka „Tożsamość" prowadzi do okna agent-builder.
Pętla czyszcząca porównuje się z kodem okna, nie z kodem zakładki — inaczej skasowałaby
znacznik postawiony dla okna docelowego.

Wykaz braków ma na końcu modułu, poza pasami okien, jeszcze jedno uzasadnienie: pozycje
mają też własne kontrolki tam, gdzie operator ich szuka, a tutaj stoi ich liczba.

Odczyty przy wczytaniu modułu idą równolegle, ponieważ odmowa jednego zostaje w jego
oknie i nie gasi pozostałych.

Drzewo wyboru narzędzi zakłada nasłuchy zamknięcia kliknięciem obok i klawiszem Escape
na dokumencie, więc bez jawnego zamknięcia przeżyłoby własne okno i reagowało na
klawiaturę w module, którego już nie ma.

## budowa/klient-poprzedni/src/moduly/agents/okno-agent-builder.ts
Okno składa się z trzech części: widoku „Biblioteka ekspertów", widoku „Edytor" z paskiem
zakładek oraz panelu „Historia wersji". Zakładka „Tożsamość" pokazuje edytor tego okna;
pozostałe zakładki przenoszą ognisko do właściwego okna modułu — model bazowy mieszka
w Model Configuration, umiejętności w Skills Manager i tak dalej. Pasek nie powiela więc
żadnego formularza. Panel „Zespoły ekspertów" składa nazwany skład z tego samego wykazu,
który okno pokazuje po lewej, dlatego dostaje bibliotekę z metody odświeżenia i nie
odpytuje rdzenia po raz drugi. Licznik narzędzi stoi pod tożsamością, bo ekspert to
tożsamość plus dobór narzędzi. Liczba nie jest tu liczona po raz drugi: licznik narzędzi
jest jednym bytem obsadzonym także w Skills Managerze i Connectors Managerze, więc
przypisanie kodu w tamtych oknach przestawia tę liczbę w tej samej klatce.

Wywołanie agent.assignment.list bez wskazania eksperta oddaje przypisania wszystkich —
dokładnie po to, żeby karta każdego miała licznik po jednym odczycie, a nie po jednym
na kartę. Odmowa nie gasi biblioteki: karty zostają bez plakietki, bo zero wpisane
z ciszy byłoby nieprawdą.

Panel historii bierze kanał, a nie stan modułu, bo wykaz wersji nie da się wyprowadzić
z wykazu biblioteki. Po przywróceniu ekspert wraca tą samą drogą co po każdym innym
zapisie — wchłonięciem, żeby okna eksperta przerysowały się w tej klatce.

Bez oddzwonienia po archiwizacji ekspert odłożony wisiałby na liście do najbliższego
odczytu, czyli wyglądałby na niezarchiwizowanego.

Atrybut aria-selected ustawione raz przy budowie byłoby fałszywym stanem: pasek
meldowałby wybór tożsamości nawet po kliknięciu innej zakładki, a czytnik ekranu
dostawałby zapewnienie o wyborze, którego operator nie dokonał. Ognisko przenosi moduł,
ale to pasek wie, którą zakładkę przycisnięto, więc znakowanie należy do niego.

## budowa/klient-poprzedni/src/moduly/agents/okno-katalog-rozszerzen.ts
Skills Manager i Connectors Manager pracują na przypisaniach do jednego eksperta;
katalog jest bytem szerszym — wykazem pozycji platformy, z których dopiero się wybiera.
Dlatego stoi w module agents, w pasie zarządców. Kontrakt nie rozstrzyga, co znaczy
instalacja pozycji, więc okno nazywa tę granicę wprost pod formularzem, zamiast
obiecywać pobranie paczki. Filtr rodzaju idzie do rdzenia jako pole kind komendy
extension.list, a nie ukrywa wierszy w przeglądarce — inaczej licznik pozycji mówiłby
o czymś innym niż wykaz pod nim.

Wykaz przerysowany z odpowiedzi pojedynczej pozycji potrafiłby pokazać stan, którego
rdzeń nie ma, dlatego wywołanie zmieniające katalog zawsze kończy się ponownym
odczytem wykazu.

## budowa/klient-poprzedni/src/moduly/automations/walidacja-definicji.ts
Rdzeń ocenia układ zależności komendą `orchestration.validate`, lecz ocenia go dopiero po zapisie
i wyłącznie na krawędziach grafu. Zastrzeżenia rozstrzygalne bez pytania rdzenia — krok bez
identyfikatora, identyfikator powtórzony, krok odwołujący się do poprzednika, którego w definicji
nie ma, krok bez treści właściwej jego rodzajowi — okno wypowiada od razu, przy wpisywaniu, ponieważ
czekanie z nimi na odpowiedź rdzenia oznaczałoby zapis definicji wadliwej. Zastrzeżenie nie jest
bramą: zapis pozostaje możliwy, a wynik jest ostrzeżeniem sygnalizowanym przy kroku, którego dotyczy;
tak samo postępuje rdzeń z układem zależności zawierającym cykl. Plik jest czysty: nie dotyka
dokumentu i nie woła rdzenia. Krok bez połączenia jest zgłaszany dopiero wtedy, gdy definicja
w ogóle używa zależności — definicja bez ani jednej zależności wykonuje kroki w kolejności zapisu
i jest poprawna, więc każdy jej krok byłby wtedy zgłoszony bez powodu. Rachunek cyklu jest zwykłym
przeglądem w głąb ze znacznikiem odwiedzin: krok napotkany powtórnie na tej samej ścieżce zamyka
cykl, a kroki cyklu wracają kompletem, bo sygnalizacja stoi przy każdym z nich, a nie przy jednym
wybranym.

## budowa/klient-poprzedni/src/aktualizacja/most-aktualizacji.ts
Przeglądarka nie podmieni pliku wykonywalnego i nie uruchomi aplikacji ponownie;
potrafi to wyłącznie powłoka natywna, bo to ona stawia proces i zna swoje
miejsce na dysku. Interfejs rozpoznaje, że jest co zakładać, i przekazuje
żądanie. Suma kontrolna idzie w żądaniu, bo powłoka ma odmówić założenia pliku,
którego suma się nie zgadza — interfejs nigdy nie woła aktualizacji bez sumy
z wykazu wydań. Poza powłoką natywną i przy powłoce, która polecenia nie zna,
wynikiem jest nazwana odmowa, nie wyjątek i nie cisza.

Wykaz poleceń powłoki tego polecenia jeszcze nie zawiera, więc zapytanie o drogę
oddaje `null` — powłoka nie mówi — a baner nie obiecuje wtedy restartu.
Dołożenie polecenia po stronie powłoki niczego tu nie łamie: odpowiedź zaczyna
przychodzić, zdanie banera robi się dokładniejsze. Podobnie zdarzenie postępu:
powłoka go jeszcze nie rozgłasza, pętla pobrania liczy bajty, ale ich nie
wysyła; nasłuch stoi założony, żeby dołożenie rozgłoszenia po stronie powłoki
było jedyną potrzebną zmianą.

Odpowiedź o przebiegu dociera przed restartem: powłoka odkłada ponowne
uruchomienie, żeby baner zdążył powiedzieć „udało się”. Bez tej zwłoki okno
ginęłoby przed odebraniem odpowiedzi i wyglądałoby to jak awaria.

Odmowa niesie osobno kod i zdanie, bo to dwie różne rzeczy dla dwóch różnych
odbiorców: zdanie czyta Operator, kod czyta baner. Bez kodu każda odmowa
wyglądałaby na trwałą, także brak łączności, który minie za minutę. Kody odmów
przemijających opisują brak, który może wrócić; reszta kodów opisuje brak
trwały, więc baner zostaje przy przeczytanym powodzie aż do następnego obiegu
pytania.

Sama obecność powłoki nie wystarcza za odpowiedź o drodze aktualizacji: na
kopii z pakietu instalacyjnego powłoka jest, ale podmiany nie wykona. `null`
znaczy „powłoka nie mówi”, a nie „nie da się” — wywołujący ma wtedy milczeć
o restarcie, a nie zgadywać w którąkolwiek stronę.

Pole całości pobrania bywa puste, i to nie jest brak danych do załatania:
serwer nie musi podać nagłówka długości treści, wtedy znana jest wyłącznie
liczba bajtów już pobranych, a widok pokazuje licznik megabajtów zamiast
wymyślonego udziału procentowego.

Obietnica założenia wydania rozstrzyga się wyłącznie odmową: przy powodzeniu
powłoka zamyka proces i nikt na wynik nie czeka, bo nie ma już czego czekać.
Wywołujący kod musi zakładać, że może się nie doczekać odpowiedzi.

Wyjątek bez kształtu odmowy oznacza powłokę, która polecenia nie zna, albo
przerwany kanał komunikacji z powłoką; dostaje osobny kod zamiast zgadywanego
kodu powłoki. Zdanie odmowy układa powłoka, nie interfejs: tylko ona wie, czy
zabrakło łączności, prawa zapisu, czy zgodności sumy kontrolnej.

## budowa/klient-poprzedni/src/moduly/agents/okno-model-configuration.ts
Lista kanałów pochodzi z komendy channel.list — tego samego rejestru, z którego biorą
kanał okna rozmowy. Okno nie ma własnej listy dostawców i nie zna nazw programów CLI;
kanały „Code CLI", „Agent SDK" i „API" są wierszami tego rejestru wraz z transportem,
nie gałęziami w kodzie. Zmiana kanału przełącza zestaw pól zależnych. Zestaw bierze się
z deklaracji zdolności adaptera (obszar model konfiguracji sesji), więc przełącza się
sam, gdy zmieni się kanał albo transport. Komendy okna to zapis modelu bazowego oraz
odczyt rejestru kanałów i deklaracji zdolności. Komenda nadania kanału modelu oknu
komunikacji tu nie występuje, bo Model Configuration dotyczy jednego eksperta i żadnego
okna komunikacji nie prowadzi — nie ma czym wskazać identyfikatora okna. Wszystkie
cztery pola są odczytywalne z bytu eksperta, a bazą trzyma je migracja agentów. Formularz
pokazuje więc stan rdzenia, a nie własną pamięć — przy zmianie eksperta pola przyjmują
wartości nowego, a nie zostają po poprzedniku. Puste pole znaczy „ekspert tego nie ma",
i tak je opisuje.

Parametry są w kontrakcie wartością JSON dowolnego kształtu, więc do pola idą tekstem
sformatowanym, nie surowym rzutowaniem na napis, bo to dałoby treść nieczytelną zamiast
tej, którą operator ma poprawić. Brak wartości daje pole puste, a nie napis „undefined".

## budowa/klient-poprzedni/src/moduly/diagnostics/zrodlo-kondycji.ts
Sonda bez ani jednego przebiegu nie ma dostępności, a rdzeń oddaje wtedy pustą
wartość procentową; okno pokazuje tę pustkę zamiast stu procent, ponieważ brak
pomiaru nie jest dowodem sprawności. Powtórzenie wywołania modelu nie jest
odczytem, tylko nowym wywołaniem kanału z własnym kosztem i własnym wierszem
śladu, dlatego okno woła je wyłącznie na jawne żądanie użytkownika, nigdy przy
odświeżeniu wykazu.

## budowa/klient-poprzedni/src/moduly/agents/panel-zakresow-narzedzi.ts
Panel stoi w oknie Permissions Center mimo że dotyczy profilu asystenta, a nie eksperta, ponieważ to jedyne miejsce, w którym Operator ustala zakres działania wykonawcy w jego imieniu. Zakres jest nastawą, nie bramką wbudowaną: pozycja bez wiersza pozostaje dostępna bez granicy wywołań, a platforma niczego nie zawęża z góry — zawężenie zapisane w panelu rdzeń odczytuje przed każdym wywołaniem narzędzia i na jego podstawie odmawia, a kolumna zużycia pokazuje, ile z granicy zostało.

## budowa/klient-poprzedni/src/moduly/diagnostics/okno-observability-tools.ts
Kontrakt niesie z całej diagnostyki jedną rodzinę komend wraz ze zdarzeniem
postępu, więc każda zakładka bez pokrycia nie znika ze sceny: znika dopiero
funkcja, o której nikt nie wie, że jej nie ma, a powód każdego braku składa
się z wykazu komend kontraktu przy składaniu okna. Widoczność zakładek wiąże
się z ustawieniami obszaru diagnostyki, których rdzeń dziś nie zna ani
jednego, więc okno pokazuje wszystkie pięć zakładek zamiast ukrywać jedną na
podstawie wartości zmyślonej. Prowenancja wchodzi osobnym źródłem, bo rodzina
jej komend, wraz z oceną wywołania i wydaniem śladu, nie należy do portu
diagnostyki, a bez tego źródła zakładka nazywa własny brak po stronie
złożenia modułu. Reguła alertu zbudowana wyłącznie w oknie żyłaby do
zamknięcia karty i nie zadziałałaby, gdy przeglądający nie patrzy — czyli
dokładnie wtedy, gdy alert ma sens.

## budowa/klient-poprzedni/src/moduly/diagnostics/prowenancja-slowa.ts
Wartości kontraktu są angielskie i techniczne, a w oknie stoi pełna nazwa
polska, ponieważ zakaz numeracji i kodów w produkcie znaczy również zakaz
pokazywania wartości pola zamiast jej nazwy. Przekład jest jednostronny: do
rdzenia jedzie wyłącznie wartość kontraktu, nigdy napis z tego pliku. Wykazy
są pełne wobec kontraktu i kompilator tego pilnuje, ponieważ typ wykazu po
typie wartości nie skompiluje się, gdy kontrakt dołoży stan albo format, więc
nowa wartość nie przemknie do okna jako pusty napis. Format wydania śladu
obejmuje wszystkie cztery formaty telemetrii, które kontrakt przyjmuje:
OpenTelemetry Protocol i JSON dla śladu, JSON Lines dla dziennika oraz
wartości rozdzielone przecinkiem dla kosztu i błędów; odjęcie formatu, który
rdzeń przyjmuje, byłoby brakiem funkcji zrobionym w oknie. Pole nieoddane
przez rdzeń nie staje się zerem: koszt zerowy i kanał bez cennika to dwa różne
zdania o instalacji, a wywołanie w biegu nie ma jeszcze opóźnienia ani liczby
tokenów, co nie jest usterką.

## budowa/klient-poprzedni/src/moduly/diagnostics/prowenancja-wiersz.ts
Odczyt śladu, ocena odpowiedzi i wskazanie wywołania do wydania mieszkają
w jednym wierszu, ponieważ wszystkie trzy dotyczą tego samego wywołania i
nigdzie indziej nie miałyby czego dotyczyć. Rozwinięcie wiersza nie pyta
jednocześnie rdzeń o treść promptu, bo byłoby to odczytem materiału
wrażliwego zrobionym przez pomyłkę w celowaniu, dlatego odczyt śladu ma
własny przycisk. Ocena zjawia się dopiero po naciśnięciu przycisku oceny:
sąd nad każdym wierszem wykazu z góry byłby formularzem zawsze, a przy
pięćdziesięciu wywołaniach ścianą pól, w której nie widać samych wywołań.
Drzewo odcinków śladu wychodzi z kontraktu w kolejności od korzenia, więc
głębokość liczy się z pola wskazującego odcinek nadrzędny, a nie z kolejności
w wykazie, ponieważ wywołanie z podagentem ma odcinki zagnieżdżone i płaska
lista zgubiłaby to, kto kogo wywołał. Wspólny zapis braku zlałby wyłączony
zapis treści z pustą odpowiedzią modelu, dlatego okno rozróżnia te dwa stany
osobnym zdaniem.

## budowa/klient-poprzedni/src/moduly/diagnostics/prowenancja-wydanie.ts
Treść promptu i odpowiedzi jest w kontrakcie domyślnie wyłączona, więc
przełącznik startuje wyłączony i jego zdanie mówi, co jego włączenie wynosi
z instalacji; wartość początkowa odwrotna wynosiłaby treść rozmów poza rdzeń
przez samo naciśnięcie przycisku. Plik oddaje przeglądarka, a nie rdzeń, bo
treść wydania przychodzi w odpowiedzi i nie ma po co wracać do rdzenia po
drugie zapisanie tego samego. Rdzeń ma prawo wydać w innym formacie niż
poproszono i ma prawo zredagować treść; jedno i drugie musi być widoczne
w zdaniu o pliku, bo od tego zależy, czym plik wolno się posłużyć.

## budowa/klient-poprzedni/src/moduly/diagnostics/zakladka-prowenancji.ts
Zakładka wykonuje cztery czynności z własnym, widocznym chwytem każda: wykaz
wywołań zawężony filtrami paska i zakresem czasu wspólnym modułowi, odczyt
jednego śladu przy wierszu wraz z drzewem odcinków oraz treścią promptu
i odpowiedzi, ocenę odpowiedzi jako osobny formularz przy wierszu oraz
wydanie śladu do pliku na urządzenie wraz ze zdaniem o tym, co ten plik
niesie. Powtórzenia wywołania zakładka nie ma: komenda jest w kontrakcie, ale
brak chwytu jest rozstrzygnięciem po stronie klienta, ponieważ powtórzenie
wysyła prompt do modelu ponownie, więc wydaje pieniądze instalacji i nie
stoi obok przycisków odczytu. Korelacja z dziennikiem została polem, a nie
widokiem: pole identyfikujące proces jest wspólne wpisowi dziennika,
telemetrii procesu i wywołaniu modelu, więc wpis z przeglądarki dzienników
prowadzi tutaj przez filtr procesu, a drugie przeszukanie dziennika w tej
zakładce byłoby powtórzeniem tamtej pracy nad materiałem, który przychodzi
wprost z rejestru wywołań. Odczyt jest leniwy: zakładka pyta rdzeń dopiero
przy otwarciu i sama pamięta, że jeszcze nie pytała.

## budowa/klient-poprzedni/src/moduly/automations/wzorce-cyklicznosci.ts
Operator ustala cykliczność wzorcem („co tydzień, poniedziałek, 07:00”), a kontrakt niesie ją jednym polem `AutomationSchedule.cron`. Plik jest przekładem w obie strony: składa zapis z nastaw i rozpoznaje wzorzec w zapisie odczytanym z rdzenia, żeby okno otwarte na harmonogramie zastanym pokazywało wzorzec, a nie samą składnię cron. Przekład jest zawężony do wzorców, które da się zapisać pięcioma polami notacji: zapis odczytany z rdzenia i niepasujący do żadnego wzorca zostaje wzorcem własnym, a jego treść idzie do pola zapisu bez zmiany. Rachunek kolejnych terminów uruchomienia należy do pliku `nastepne-uruchomienia.ts`.

## budowa/klient-poprzedni/src/moduly/diagnostics/zakladki-narzedzi.ts
Zakładka niewidoczna nie jest zakładką porzuconą: obszar zostaje w drzewie
i traci wyłącznie widoczność, więc wpisany filtr i odczytany wykaz przeżywają
zajrzenie do sąsiedniej zakładki. Wędrówka strzałkami należy do wzorca
zakładek: pas ma jeden przystanek tabulatora na zakładce czynnej, a strzałki
przenoszą wybór między zakładkami, więc pas pięciu przycisków nie staje się
pięcioma przystankami przed treścią. Bliźniaczy mechanizm stoi w module
zakładek sekcji dla okna modeli, lecz jest przywiązany do własnej rodziny
klas i do arkusza tamtego okna; wspólnego komponentu zakładek biblioteka
komponentów dziś nie ma. Obudowa ciała zakładki stoi w jednym miejscu, a nie
w każdej zakładce z osobna, ponieważ pięć zakładek składających własne
pudełko rozjechałoby się przy pierwszej zmianie odstępu, a różnią się
treścią, nie kształtem. Akapit objaśnienia jest częścią wyposażenia
kontenera, a nie ozdobą jednej zakładki, ponieważ zakładki warstwy
eksperckiej mówią o własnej granicy zdaniem, nie milczeniem.

## budowa/klient-poprzedni/src/moduly/browser/czynnosci-automatyk.ts
Plik ma jedną odpowiedzialność: wywołania rdzenia w imieniu okna. Okno składa kontrolki i wykaz; tutaj mieszka to, co dzieje się po naciśnięciu, wraz z oceną odpowiedzi w pliku `skutek-zapisu.ts`. Każda czynność odpowiada zapowiedzią przed wywołaniem, a po odpowiedzi rdzenia albo wynikiem, albo powodem odmowy. Wynik pusty znaczy, że rdzeń nie oddał treści i powód już podano — okno nie ma wtedy czym zastąpić wykazu i zostawia ten, który operator widzi.

## budowa/klient-poprzedni/src/moduly/diagnostics/zrodlo-diagnostics.ts
Port diagnostyki jest w rdzeniu odbiorcą odmów wykonania komend, a przeglądarka
błędów pokazuje te odmowy; źródło, które połknęłoby odmowę odczytu i oddało
pusty wykaz, kłamałoby o własnym niepowodzeniu w oknie poświęconym
niepowodzeniom cudzym. Asymetria kontraktu jest zamierzona: identyfikator
okna występuje wyłącznie przy uruchomieniu analizy i jest tam nieobowiązkowy,
ponieważ analiza zapisuje, z którego okna ją uruchomiono, a pozostałe trzy
komendy czytają dziennik, błędy i rekomendacje całej instalacji, więc pojęcia
okna nie mają. Wynik warstwy protokołu nie odróżnia sam z siebie odmowy
rdzenia od odpowiedzi, której klient nie zrozumiał, ani od odpowiedzi bez
treści, a to trzy różne zdania, z których tylko jedno znaczy odmowę rdzenia;
rozstrzygnięcie zapada w tym miejscu, bo tylko tutaj widać surową odpowiedź
przed sprawdzianem kształtu. Kolejność sprawdzeń w rozstrzygnięciu jest
istotna: dopóki nie wiadomo, czy odpowiedź surowa była odmową, nie wolno
sprawdzać jej kształtu i wziąć jego kodu za kod rdzenia. Rozstrzygnięcie ma
też odbiorcę poza tym plikiem, w źródle obserwowalności, gdzie rozróżnienie
odmowy od odpowiedzi nieczytelnej ma tę samą wagę co tutaj; druga kopia tego
rozstrzygnięcia dałaby dwa zdania o jednej ciszy rdzenia.

## budowa/klient-poprzedni/src/moduly/diagnostics/zuzycie-zakladka.test.ts
Sprawdzian zakładki zużycia pilnuje czterech rzeczy stanowiących o odbiorze:
że zestawienie ma drogę z okna i idzie jednym wymiarem, że pusty okres ma
zdanie zamiast pustego miejsca, że koszt niepełny mówi o sobie, i że raport
rozliczeniowy wytwarza plik nazwany, a nie ciszę po naciśnięciu. Pobranie
pliku podstawia w sprawdzianie wyłącznie te dwa punkty styku z API
przeglądarki, którego środowisko sprawdzianu nie ma w całości.

## budowa/klient-poprzedni/src/moduly/agents/wykaz-wtyczek.ts
Wykaz wtyczek jest bytem odrębnym od konektora: konektor to droga do usługi, serwer MCP wskazany punktem dostępu podawanym powłoce przełącznikiem --mcp-config, a wtyczka to katalog rozszerzeń powłoki — zbiór poleceń, zaczepów i umiejętności wgrywany przez --plugin-dir. Kontrakt rozdziela oba byty osobnymi polami i komendami, więc rozdziela je i okno. Pole eksperta niesie same identyfikatory wtyczek, dlatego nazwę, wersję i źródło wiersz bierze z osobnego odczytu definicji; identyfikator zostaje w podpowiedzi wiersza, bo nim posługuje się czynność odłączenia. Gdy odczyt definicji odmówi albo jeszcze nie wrócił, wiersze pokazują same identyfikatory — uboższa treść nie jest awarią i nie odbiera przycisku odłączenia. Odczyt definicji odrzuca odpowiedź przedawnioną: między wysłaniem żądania a powrotem odpowiedzi operator może wybrać innego eksperta, a wtedy wykaz pokazywałby cudze wtyczki pod właściwymi kodami.

## budowa/klient-poprzedni/src/moduly/library/archiwum-repozytorium.ts
Rozdział zdolności wynika z kontraktu, nie z wygody: migawka repozytorium
i paczka migracyjna powstają w części opisowej, a nie jako archiwum, więc
nazwanie manifestu paczką migracyjną byłoby obietnicą archiwum, którego
w pliku nie ma. Utrwalenie archiwalne idzie do rdzenia własną komendą modułu:
zapis w formacie PDF/A, pakiet BagIt i profil PREMIS/METS pracują na zasobie
biblioteki, a nie na magazynie zasobów projektowych, więc identyfikator się
zgadza i wynik wraca wraz z zapisem walidacji. Paczka migracyjna i migawka
repozytorium niosą bajty: archiwum powstaje w rdzeniu i ląduje w repozytorium
jako nowy zasób, a wywóz opisowy zostaje obok jako dwie różne rzeczy, których
okno nie myli. Wskaźnik retencji czyta politykę z rdzenia wraz z raportem
zasobów zbliżających się do końca okresu.

## budowa/klient-poprzedni/src/moduly/browser/czynnosci-materialu.ts
Plik ma jedną odpowiedzialność: wywołania rdzenia w imieniu panelu. Panel składa kontrolki i wykaz, pamięć materiału stoi w pliku `material-sesji.ts`, a ocena odpowiedzi w pliku `skutek-zapisu.ts`. Sprawdzenie monitora przechodzi pod jego adres, ponieważ polecenie `browser.snapshot.get` oddaje treść strony bieżącej okna, a nie dowolnej strony. Sprawdzenie przestawia więc wspólny podgląd — operator ma o tym wiedzieć przed naciśnięciem, nie po nim, więc mówi mu to zarówno dymek przy pozycji, jak i zdanie odpowiedzi.

## budowa/klient-poprzedni/src/moduly/library/czynnosci-zbiorcze.ts
Panel Explorera wymienia osiem czynności zbiorczych: etykietę, przypisanie do
kolekcji, otwarcie w module źródłowym, przeniesienie w strukturze,
archiwizację i przywrócenie, usunięcie trwałe oraz udostępnienie
odnośnikiem. Usunięcie trwałe pyta osobno, ponieważ jest jedyną czynnością
nieodwracalną w tym module, więc pierwsze naciśnięcie zapowiada skutek,
a dopiero drugie go wykonuje; potwierdzenie w tym samym kliknięciu byłoby
zgodą pozorną. Plik buduje wyłącznie kontrolki i pokazuje odpowiedź,
a rozmowę z rdzeniem prowadzi moduł zapisów zbiorczych, zachowując jedną
odpowiedzialność na plik.

## budowa/klient-poprzedni/src/moduly/library/dolozenie-wersji.ts
Dołożenie wersji jest jedynym wejściem, którym rośnie historia dokumentu,
ponieważ wgranie pliku zawsze zakłada nowy dokument, więc drugie wgranie
dałoby drugi dokument, a nie drugą wersję pierwszego. Kamień milowy stawia
rdzeń na treści bieżącej pliku, przepisując ją do wersji, więc powrót do
znacznika niczego nie wymazuje. Wersja zapisana nie jest jeszcze wersją
z treścią: zapis z treścią wraca powodzeniem i zakłada wersję z własną sumą
kontrolną także wtedy, gdy podgląd dokumentu odmawia, więc do takiej wersji
przywrócenie nie ma po co wracać, dopóki okno nie zapyta rdzeń o treść po
dołożeniu. Zdanie powodzenia bierze osobno wskazanie wersji przez dokument
i osobno zgodność sum kontrolnych, każde z odpowiedzi rdzenia, nie z żądania
okna. Rozbieżność pola wskazania wersji znaczy, że wiersz historii powstał,
ale dokument został przy innej wersji, więc okno ogłasza to odmową.
Rozbieżność sum kontrolnych znaczy, że dokument niesie wersję, ale nie jej
bajty; wtedy zdanie mówi o wersji tyle, ile rdzeń podał, i ani słowa więcej.

## budowa/klient-poprzedni/src/moduly/agents/panel-zespolow.ts
Panel zespołów ekspertów stoi przy bibliotece ekspertów, bo zespół jest nazwanym składem tej samej biblioteki; osobne okno musiałoby wykaz ekspertów powielić albo pokazywać skład samymi identyfikatorami. Zapis bez wskazanego zespołu zakłada nowy — to drugie znaczenie tego samego przycisku, wypisane przy nim wprost. Zespół założony na innym urządzeniu konta dochodzi zdarzeniem zmiany zespołu, więc panel nie odpytuje rdzenia w pętli.

## budowa/klient-poprzedni/src/moduly/browser/czynnosci-zrodel.ts
Plik ma jedną odpowiedzialność: rozmowa z rdzeniem w imieniu panelu źródeł. Panel składa kontrolki i wykaz, a tutaj mieszka to, co dzieje się po naciśnięciu. Usunięcie źródła nie usuwa: kontrakt niesie polecenie `browser.source.remove`, ale ta czynność jeszcze go nie wywołuje. Powód bierze się z odczytu wykazu komend rdzenia, a pozycja zostaje w wykazie, bo zniknięcie wiersza bez zapisu w rdzeniu byłoby udawaniem wykonania. Adres w zdaniu potwierdzenia bierze się z odpowiedzi, nie z pola formularza, opisanego w pliku `skutek-zapisu.ts` — do rdzenia idzie adres przycięty, więc zdanie o skutku mówi o tym, co wróciło, a nie o surowej treści pola.

## budowa/klient-poprzedni/src/moduly/apps/okno-app-catalog.ts
Katalog rozszerzeń stoi poziom wyżej niż okno modułu, bo żadna komenda obszaru extension nie niesie identyfikatora okna — działa więc także wtedy, gdy rdzeń nie wskazał jeszcze okna modułu, i nie powtarza zdania o brakującym oknie, które go nie dotyczy. Instalacja niesie pochodzenie, bo pochodzenie rozstrzyga stan wyjściowy pozycji: pozycja danaco staje włączona, pozycja personal wyłączona; to jedyne miejsce, w którym źródło zmienia zachowanie rdzenia, poza nim jest faktem do pokazania operatorowi. Okno mówi o tym skutku przed instalacją, zamiast zostawiać operatora ze zdziwieniem, że pozycja nie działa po zainstalowaniu. Zawężenie idzie po stronie klienta nad tym, co rdzeń oddał, bo odczyt wykazu rozszerzeń zawęża wyłącznie rodzajem i stanem zainstalowania; pusty wykaz przy czynnym filtrze niesie zdanie o zawężeniu, inaczej czytałby się jak pusty rejestr. Przełącznik widoku siatka i lista jest czynnością wyłącznie okienną: nic nie jedzie do rdzenia, więc nośnikiem stanu jest atrybut dostępności kontrolki, nie komenda.

## budowa/klient-poprzedni/src/moduly/library/dostepnosc-tresci.ts
Wgranie pliku odpowiada powodzeniem także wtedy, gdy bajtów nie ma gdzie
zapisać, więc samo powodzenie zapisu nie jest świadkiem wgrania treści.
Kontrakt nie ma pola osiągalności treści, a metryka pliku o niej nie mówi:
plik może mieć sumę kontrolną i identyfikator wersji, a nie mieć treści,
więc jedynym świadkiem jest odpowiedź rdzenia na podgląd, nigdy własny
domysł modułu. Stan nieznany jest osobny od stanu braku: nikt jeszcze nie
pytał to co innego niż rdzeń odmówił treści. Stan odmowy jest osobny, bo
rdzeń odmawia podglądu trzema kodami i tylko jeden z nich orzeka o treści;
pozostałe dwa mówią o pliku albo o nośniku, nie o samej treści, a sprawdzian
kształtu odpowiedzi po stronie klienta sam wytwarza ten sam kod, gdy
odpowiedź nie ma zapowiedzianego pola, więc i taka odpowiedź trafia na stan
braku. Stan odwołania oddziela wskazanie miejsca treści od samej treści:
podgląd pliku, którego odwołanie wskazuje ścieżkę nieistniejącą, wraca
powodzeniem tak samo jak podgląd pliku leżącego w magazynie rdzenia, więc
świadkiem treści jest sama treść odpowiedzi. Werdykt jest wydzielony do
osobnej funkcji, ponieważ powstaje w dwóch miejscach: przy pytaniu zadanym
wprost o treść i przy zwykłym podglądzie okna podglądu pliku, a oba mają
nazywać ten stan tak samo.

## budowa/klient-poprzedni/src/moduly/assistant/panel-faktow.ts
Zakładka realizuje zasadę jawności modułu pamięci: cztery czynności kontraktu wystarczają, odczyt zakłada, zmienia i kasuje wpis, w tym przypięcie i pochodzenie. Zasięgiem odczytu jest karta sesji, nie projekt, bo moduł asystenta pracuje w karcie sesji środowiska rozmów i nie ma pojęcia projektu; zapis idzie tą samą drogą, z zasięgiem wskazanym jawnie przez operatora. Potwierdzenie mówi to, co zapisał rdzeń, a nie to, co wysłało okno: wpis zapisany na poziomie szerszym niż karta sesji bywa niewidoczny w wykazie poniżej, więc samo odświeżenie listy niczego by nie potwierdzało. Reguł retencji, wygaszania i znaczników wrażliwości okno nie udaje — wpis pamięci nie niesie czasu życia ani wrażliwości, więc brak nazywa przycisk, zamiast stawiać formularz, którego rdzeń nie zapisze.

## budowa/klient-poprzedni/src/moduly/browser/edytor-scenariusza.ts
Kroki mają dwie postacie: karty w oknie i tekst tutaj. Tekst jest potrzebny tam, gdzie kart nie starcza — przy warunkach, zależnościach i treści żądania kroku, których nie da się wyklikać, bo są dowolnym obiektem JSON. Zapis przyjmuje wyłącznie JSON. Widok YAML jest odczytem: klient nie niesie czytnika YAML, a napisanie własnego na potrzeby jednego pola byłoby budowaniem parsera języka, nie okna. Notacja jest więc przełącznikiem prezentacji, a nie dwiema drogami zapisu — okno mówi to wprost, zamiast przyjąć YAML i po cichu go zgubić. Sprawdzian treści jest surowy i wskazuje pozycję: definicja z krokiem bez rodzaju przeszłaby do rdzenia i wróciła jako scenariusz, który nic nie robi.

## budowa/klient-poprzedni/src/moduly/browser/edytor-scenariusza.ts — zapis YAML
Wartości napisowe idą w cudzysłowie zapisanym regułą JSON, bo napis JSON jest zarazem poprawnym skalarem YAML w cudzysłowie podwójnym. Rozbieranie dowolnego obiektu żądania kroku na osobne wiersze YAML niczego by tu nie wyjaśniło.

## budowa/klient-poprzedni/src/moduly/automations/panel-dopelnien-ukladu.ts
Krawędź mówi tylko, że kroki się schodzą; panel odpowiada osobno na cztery pytania, których krawędź nie wyraża: kiedy tory scalają się w kroku wspólnym (bramka dołączenia), czy zbiór kroków biegnie razem czy jeden po drugim (grupa), co wycofuje skutki kroku, gdy przebieg pękł w pół (kompensacja), oraz czyim silnikiem jadą kolejki tej automatyki (spięcie z silnikiem kolejek). Panel stoi obok panelu zależności w narzędziach kontekstowych okna Orchestrator, gospodarując tym samym układem, tylko innym jego wymiarem. Ocena bramki nie blokuje zapisu, tak samo jak przy krawędzi: bramka na kroku, którego jeszcze nie ma, zapisuje się, a zastrzeżenie wraca w odpowiedzi, bo Workflow Builder buduje układ krok po kroku i odmowa kazałaby operatorowi układać go w jedynej dopuszczonej kolejności.

## budowa/klient-poprzedni/src/moduly/library/eksporty-biblioteki.ts
Żadna z tych czynności nie ma komendy kontraktu i mieć jej nie musi:
materiałem jest odpowiedź rdzenia leżąca w oknie, a nie bajty, po które
trzeba by wrócić do repozytorium. Granica jest ostra i przebiega przy treści
plików: wywóz opisu składa się tutaj, wywóz zawartości nie, bo klient nie ma
komendy pobierającej bajty pliku biblioteki. Każdy wytwór nazywa w nagłówku,
czego nie zawiera, ponieważ manifest bez bajtów podpisany jako paczka
migracyjna byłby obietnicą archiwum, którym nie jest. Kolekcja jest
w kontrakcie samym identyfikatorem: pole kolekcji niesie kody, a założenie
kolekcji oddaje kod i nazwę wyłącznie w chwili założenia, więc mapa wypisuje
kody i przypisane im pliki i mówi to wprost zamiast podstawiać kod w miejsce
nazwy. SKOS i RDF to standardy branżowe, a nie oznaczenia wymyślone na
potrzeby modułu; brak relacji jest wypisany w komentarzu wytworu, bo tezaurus
bez relacji wygląda jak tezaurus, w którym relacji nie ustalono. Manifest
odpowiada części opisowej migawki repozytorium i paczki migracyjnej, a części
z bajtami nie ma.

## budowa/klient-poprzedni/src/moduly/apps/strona-dystrybucji.test.ts
Sprawdziany pilnują pięciu zachowań, które łatwo zepsuć po cichu: zawężenie katalogu odsiewa po polach niesionych przez pozycję i nazywa zakres szukania zamiast pozwalać czytać pustkę jako brak pozycji; odczyt katalogu zastępuje zbiór, a nie dokłada do niego, inaczej w wykazie zostawałaby pozycja, której rdzeń już nie zna; zdarzenie zmiany rozszerzenia zmienia zbiór niezależnie od tego, gdzie zaszła zmiana, a usunięcie zdejmuje też wskazanie panelu bocznego; żadna kontrolka okien nie ma atrybutu wyłączenia zgodnie z zasadą zero blokad platformy, więc kontrolka bez pokrycia ma być klikalna i nazywać brak; instalacja niesie pochodzenie pozycji, bo to jedyne pole rozstrzygające stan wyjściowy rejestracji.

## budowa/klient-poprzedni/src/moduly/library/etykiety-biblioteki.ts
Teksty stoją w osobnym pliku, bo to inna odpowiedzialność niż budowa widoku:
zdania mówiące, czego kontrakt nie niesie, dają się przeczytać w jednym
miejscu i tam poprawić, gdy komenda wejdzie do kontraktu. Wykaz czynności
panelu metadanych i archiwum jest długi, bo warstwa czwarta modułu jest
obszerna, a kontrakt niesie dla niej dokładnie jedną komendę, przeliczenie
wskaźnika znaczenia w zakładce higieny; inaczej wykaz wyglądałby na listę
rzeczy zapomnianych, a jest pomiarem kontraktu.

## budowa/klient-poprzedni/src/moduly/browser/etykiety-browser.ts
Plik ma jedną odpowiedzialność: słowo mówione do operatora. Pliki budujące elementy nie trzymają ani jednego zdania, bo wtedy zmiana brzmienia wymagałaby wejścia w widok, a te same zdania powtarzałyby się w kilku oknach naraz, na wzór pliku `sterowanie/etykiety-sterowania.ts`. Zdania o brakach stoją tu razem z resztą: brak drogi w kontrakcie jest treścią widoku tak samo jak nazwa przycisku — okno ma go wypowiedzieć, nie przemilczeć.

## budowa/klient-poprzedni/src/moduly/browser/etykiety-browser.ts — kody okien operacyjnych
Kody okien operacyjnych modułu są bezmodułowe: migracja rejestru okien operacyjnych nadaje definicjom okien kody bez przedrostka modułu, więc kody przychodzące z pola `Module.operationalWindowCodes` zestawiają się wprost z tymi wartościami. Rejestr rdzenia niesie dziś wiersze trzech pierwszych okien. Automation Studio i Capture & Monitor Panel opracowanie modułu wymienia na równi z nimi, więc moduł je buduje i podaje ich kody katalogowi — pasek uczciwości wypowiada wtedy rozjazd („kody budowane, których rdzeń modułowi nie przypisuje") zamiast go przemilczeć.

## budowa/klient-poprzedni/src/moduly/browser/etykiety-browser.ts — klasy dymka
Klasy własne dymka stoją w jednym miejscu, bo sięga po nie każdy dymek modułu: rozpisane literałem rozjechałyby się przy pierwszej zmianie nazwy klasy w arkuszu stylu. Fabryka jest wspólna, wygląd znaku pozostaje modułu — pierścień o wymiarze pola wyboru zamiast bibliotecznego kwadratu ikony.

## budowa/klient-poprzedni/src/moduly/browser/etykiety-browser.ts — barwy adnotacji
Żeton i klasa barwy adnotacji stoją w jednym wierszu, bo opisują tę samą barwę dwiema drogami: arkusz maluje próbkę na przycisku regułą klasy barwy, a płótno rysunkowe nie zna zmiennej stylu i musi dostać żeton po nazwie, żeby rozwiązać go w chwili rysowania. Rozdzielone na dwa wykazy rozjechałyby się przy zmianie palety. Żetony są prymitywne, nie semantyczne: barwa stanu przełącza się wraz z motywem, a tusz adnotacji ma być ten sam w obrazie wysłanym z motywu jasnego i z ciemnego, bo załącznik ogląda się poza motywem.

## budowa/klient-poprzedni/src/moduly/browser/etykiety-browser.ts — klasyfikacje notatki
Kod klasyfikacji notatki jest kluczem pamięci widoku, nazwa — napisem dla operatora. Kontrakt pola klasyfikacji nie niesie, więc oznaczenie żyje w karcie sesji i panel mówi o tym wprost.

## budowa/klient-poprzedni/src/moduly/browser/etykiety-browser.ts — pozycje nieobsłużone
Zdania o powodzie braku obsługi nie ma tutaj ani jednego. Powód rozstrzyga się przy oknie, z odczytu wykazu komend rdzenia: kontrakt komendę niesie, a rdzeń może mieć albo nie mieć jej uchwytu — i to się zmienia wraz z rdzeniem, nie wraz z tym plikiem. Zdanie wpisane tu na sztywno przestałoby być prawdziwe w dniu dobudowy obsługi i nikt by go nie zdjął. Nazwa komendy pochodzi wyłącznie ze stałych kontraktu, żeby zmiana jej nazwy w kontrakcie nie zostawiła w oknie zdania o komendzie, której już nie ma.

## budowa/klient-poprzedni/src/moduly/agents/panel-doradcy.ts
Treść rady stoi pod etykietą rady i nie trafia do pola instrukcji eksperta sama z siebie — przeniesienie jest osobnym kliknięciem i niesie nagłówek prowenancji. Wybierany jest kanał, nie nazwa modelu: rejestr kanałów jest jedynym miejscem, w którym rdzeń wie, czym się połączyć i czyim poświadczeniem; wykaz obejmuje kanały czynne poza kanałem bazowym eksperta. Powód nie odbiera przycisku i nie może tego robić, bo platforma nie stawia bram, a niegotowość sygnalizuje się po naciśnięciu komunikatem albo opisem obok kontrolki — przycisk wygaszony zabierałby Operatorowi jedyną drogę dowiedzenia się, czego brakuje, skoro tytuł bywa niedostępny z klawiatury i milczy na urządzeniu dotykowym.

## budowa/klient-poprzedni/src/moduly/browser/etykiety-browser.ts — ścieżka rysunku modelu
Rysunek idzie jako URI danych, rdzeń materializuje go do pliku na własnym nośniku i podaje modelowi ścieżkę w treści zapytania, wraz z prośbą o sięgnięcie po nią narzędziem odczytu. Obrazu wklejonego w wiadomość model nie dostaje — kanał wywołania nie ma na to pola — więc zobaczy rysunek dopiero po otwarciu pliku.

## budowa/klient-poprzedni/src/moduly/library/magazyn-biblioteki.ts
Ster modułu docelowego stoi w dwóch oknach: w przeglądarce biblioteki jako
czynność zbiorcza otwarcia w module źródłowym, i w podglądzie pliku jako ta
sama czynność dla pojedynczego pliku. Wspólne pole trzyma oba okna zgodne,
bo osobne pola pokazywałyby po zmianie dwie różne wartości. Werdykty rdzenia
o treści plików stoją osobno od zbioru, ponieważ nie są polem kontraktu:
plik nie niesie żadnej wartości mówiącej, czy repozytorium ma jego treść.
Werdykt o treści znika wraz ze zmianą pliku, ponieważ dołożenie wersji,
przywrócenie wcześniejszej i zdarzenie zmiany pliku z innego modułu
zmieniają dokument, więc poprzednia odpowiedź rdzenia przestaje go opisywać,
a stan nieznana jest tu bezpieczniejszy niż nieaktualna pewność. Z tego
samego powodu znika trafienie wskaźnika znaczenia, bo fragment pochodzi
z treści sprzed zmiany. Zawężenie wskazujące skasowany plik przestałoby
zgadzać się z wykazem: licznik zbioru mówiłby o pozycji, której rdzeń już
nie zna.

## budowa/klient-poprzedni/src/moduly/browser/etykiety-browser.ts — wpis adnotacji w wytworach sesji
Czynność „Dodaj do rozmowy" daje dwie rzeczy naraz: załącznik rozmowy oraz wpis w wytworach sesji. Pierwszą niesie komenda wysyłania wiadomości polem załączników, drugą komenda dodania wytworu — obie komendy kontrakt ma, więc okno pyta o ich pokrycie w rdzeniu, zamiast orzekać o braku.

## budowa/klient-poprzedni/src/moduly/agents/zrodlo-agentow.ts
Źródło nie ma własnego stanu i niczego nie pamięta — jest wyłącznie warstwą wywołań i sprawdzianu kształtu odpowiedzi; stan biblioteki ekspertów mieszka osobno, żeby pięć okien modułu patrzyło na jeden zbiór, a nie na pięć osobnych kopii. Żadne wywołanie nie rzuca wyjątkiem ani nie odrzuca obietnicy: niepowodzenie wraca polem błędu wyniku, a okno pokazuje je w swoim stanie błędu; dotyczy to także treści JSON wpisanej przez Operatora, gdzie niepoprawny zapis jest odmową wywołania, nie wyjątkiem wywracającym widok. Trybu podania zapis warstwy nie niesie, bo kontrakt go nie ma: warstwa eksperta dopisuje się do promptu systemowego jako zakres użytkownika i nigdy go nie zastępuje.

## budowa/klient-poprzedni/src/moduly/browser/etykiety-browser.ts — stany puste okien
Stan pusty opisuje sytuację oczekiwaną, czyli pierwsze użycie, a nie awarię odczytu; napis „Brak danych" nie mówi ani czym okno jest, ani co operator ma zrobić. Dlatego opisu żąda już fabryka budująca stan okna — okna bez zdania o sobie nie da się zbudować.

## budowa/klient-poprzedni/src/moduly/browser/etykiety-browser.ts — pochodzenie wykazów
Panel pokazujący pozycje bez słowa o pochodzeniu wygląda tak samo, gdy czyta rdzeń, i wtedy, gdy pokazuje własną pamięć — a to dwie różne obietnice wobec operatora.

## budowa/klient-poprzedni/src/moduly/browser/etykiety-browser.ts — warstwy wyzwalaczy
Znaczniki źródeł i notatek należą do warstwy drugiej interfejsu, menu operacji do warstwy trzeciej, a narzędzia warstwy czwartej nie mają w pasku nazwy, dopóki operator nie włączy trybu administracyjnego.

## budowa/klient-poprzedni/src/moduly/library/material-archiwum.test.ts
Sprawdzian pilnuje tego, co stanowi o odbiorze: że każda z trzech komend
arsenału ma drogę z okna, że żądanie idzie ścieżką z dysku, bo identyfikator
zasobu biblioteki wróciłby odmową, że czynność trwająca mówi o sobie przed
końcem, i że wynik jest nazwany, a nie potwierdzony ciszą.

## budowa/klient-poprzedni/src/moduly/agents/zrodlo-zakresu-eksperta.ts
Źródło jest osobne od źródła agentów: tamto opisuje bibliotekę — założenie, wykaz, zmianę tożsamości — a to opisuje zakres, czyli co ekspertowi wolno zrobić w systemie; rozdział jest ten sam, który przebiega w rdzeniu między portem agentów a portem zakresu eksperta. Zakres narzędzi profilu asystenta mieszka tutaj, choć dotyczy profilu, a nie eksperta, bo Permissions Center jest jedynym oknem, w którym operator ustala, jak szeroko działa wykonawca w jego imieniu, więc rozstrzygnięcie o wywoływaniu narzędzi stoi tam, gdzie operator go szuka, a nie w oknie, którego moduł Agents nie ma. Żadne wywołanie nie rzuca wyjątkiem: niepowodzenie wraca polem błędu wyniku, a okno pokazuje je w swoim stanie odmowy.

## budowa/klient-poprzedni/src/moduly/agents/warstwy-promptu.ts
Nakładkę składa rdzeń, a kontrakt niesie warstwę tożsamości i cztery komendy obszaru layer. Warstwa nie ma własnego trybu, bo dopisanie do globalnego promptu systemowego albo jego zastąpienie dotyczy instrukcji eksperta jako całości: przełącznik stoi raz, przy tożsamości; dlatego zlecenie zapisu warstwy pola trybu nie ma, a sam ekspert ma. W miejscu przełącznika stoi przy każdej warstwie zdanie czytające tryb eksperta czynnego. Edytor nie dotyka pola instrukcji systemowej eksperta — to pole płaskie zapisuje formularz tożsamości osobną komendą aktualizacji.

## budowa/klient-poprzedni/src/moduly/library/material-narzedzia.ts
Wszystkie trzy komendy arsenału rozwiązują identyfikator zasobu przez
repozytorium zasobów modułu projektowego, a plik biblioteki leży w innym
rejestrze, więc jego identyfikator wraca stamtąd odmową. Kontrakt niesie
jednak drugą drogę źródła, ścieżkę z dysku, i mówi o niej wprost: treść jest
wciągana do magazynu, nie dowiązywana, więc tą drogą czynność wykonuje się
naprawdę, a przycisk wysyłający identyfikator pliku biblioteki zawodziłby
zawsze i byłby przyciskiem pewnej odmowy. Bajty wyniku idą do tego samego
magazynu zasobów projektowych, którym jedzie wgranie zasobu. Pola okna
źródło nie podaje i to jest rozstrzygnięcie, nie przeoczenie: kontrakt każe
podać okno modułu projektowego, a moduł biblioteki zna wyłącznie okno
komunikacji sesji, więc okno zmyślone nie zapisałoby się w ogóle, a okno
cudze pokazałoby wynik w wykazie, do którego nie należy; bajty i tak
trafiają do magazynu pod sumą kontrolną, ale zasób nie pojawi się w wykazie
okna, i widok mówi to zdaniem. Rodzina komend materiału nie ma w kontrakcie
ani jednego zdarzenia, więc nic tu nie nasłuchuje: czynność kończy się swoją
odpowiedzią i niczym więcej.

## budowa/klient-poprzedni/src/moduly/agents/podsumowanie-definicji.ts
Definicja eksperta powstaje w pięciu oknach naraz i żadne z nich nie widzi całości: okno konfiguracji modelu nie wie, ile ekspert ma umiejętności, a menedżer umiejętności nie wie, jakim kanałem ekspert mówi. Panel jest jedynym miejscem, w którym widać komplet, więc stoi przy edytorze niezależnie od tego, które okno ma ognisko. Panel niczego nie zapisuje i nie woła ani jednej komendy — czyta eksperta czynnego ze stanu modułu. Wiersz jest za to przenośnikiem: kliknięcie przenosi ognisko do okna, które daną rzeczą zarządza, więc wartość pusta jest drogą do okna konfiguracji, a nie samym stwierdzeniem braku. Wartości nie są tu liczone po raz drugi: każdy wiersz czyta pole bytu eksperta, a przy pamięci i zakresie możliwości regułę stanu wyjściowego zapisaną w kontrakcie — cztery poziomy pamięci i pełny dostęp operacyjny; zdanie o braku ustawienia jako wartości domyślnej obowiązuje tu tak samo jak w Permissions Center. Stanu zarchiwizowanego w plakietce nie ma i nie jest to przeoczenie: wykaz biblioteki ekspertów zarchiwizowanych nie oddaje, więc ekspert czynny w edytorze nigdy nim nie jest — archiwum ma własny panel i własny wykaz.

## budowa/klient-poprzedni/src/moduly/browser/formularz-nawigacji.ts
Plik ma jedną odpowiedzialność: formularz przejścia i dwie komendy, które z niego wychodzą. Podgląd, pasek zaznaczenia i pasek dolny są osobno. Przewijanie przewija migawkę, nie stronę w rdzeniu: przewinięcie po stronie rdzenia ma własną komendę, której formularz jeszcze nie wywołuje. Zdanie pod przyciskami bierze powód z odczytu wykazu komend rdzenia, więc zmieni się samo w dniu dobudowy obsługi. Adres w zdaniu końcowym pochodzi z migawki, nie z pola formularza: funkcja przycinania białych znaków w przeglądarce nie jest tą samą funkcją co odpowiednik rdzenia, więc zdanie zbudowane z pola byłoby prawdziwe przypadkiem. Polecenie nawigacji oddaje migawkę z adresem, pod którym strona została faktycznie pobrana.

## budowa/klient-poprzedni/src/moduly/library/material-panel.ts
Powierzchnia stoi w obszarze archiwum panelu metadanych, przy pozostałych
czynnościach wykonywanych nad treścią, a nie nad opisem, i miejsce dla tej
rodziny komend w dokumentacji modułów pozostaje do rozstrzygnięcia; stanęła
tutaj, bo tutaj użytkownik pracuje nad zasobem i tutaj czynności arsenału
mają sąsiadów o tej samej naturze. Rozpoznanie jest osobnym krokiem, nie
ozdobą przetworzenia, dlatego jego odpowiedź zostaje na widoku. Czynność,
która trwa, mówi to zanim skończy, i mówi to samo pole, które potem poniesie
wynik, bo dwa miejsca na jedną wiadomość dałyby wybór, w które patrzeć.
Zdanie o czasie trwania zerowym nazywa to wprost, ponieważ sam zapis
zerowego czasu wyglądałby na pomiar nieudany, a dotyczy też kontenera bez
nagłówka czasu.

## budowa/klient-poprzedni/src/moduly/library/metadane-pliku.ts
Blok techniczny czyta pola pliku oraz, na żądanie, metadane osadzone
w bajtach pliku: wymiary, liczba stron, czas trwania nagrania i podobne pola
techniczne, które rdzeń oddaje po włączeniu odczytu technicznego. Odczyt
osadzonych idzie osobnym przyciskiem, bo otwiera bajty zasobu, co jest
kosztem, którego przegląd wykazu nie potrzebuje. Formularz Dublin Core
scala domyślnie: pole zostawione puste zostaje bez zmiany, a pole
wyczyszczone jawnie kasuje wartość, dokładnie tak, jak mówi kontrakt.

## budowa/klient-poprzedni/src/moduly/apps/braki-kontraktu.ts
Panele akcji okien Apps wymieniają więcej czynności, niż obszar apps niesie komend; czynność bez komendy zostaje widoczna i klikalna, a naciśnięcie mówi, czego brakuje. Powód składa się z wykazu komend kontraktu przy składaniu okna, więc dopisanie komendy do kontraktu przepisuje zdanie samo. Wobec bliźniaczego pliku modułu Diagnostics dochodzą tu dwie rzeczy: pozycja znika sama, gdy wskazana komenda wejdzie do kontraktu, a wskazanie cudzej drogi jest sprawdzane — zdanie o komendzie z innego obszaru pada wyłącznie wtedy, gdy ta komenda stoi w wykazie, w przeciwnym razie zdanie mówi o jej zniknięciu. Zdanie nie orzeka, czy złożony rdzeń komendę rejestruje — brak jest po stronie kontraktu i tylko o kontrakcie zdanie mówi. Pełny powód dla jednej pozycji jest wyeksportowany, bo ta sama treść idzie do dymka, do tytułu i do opisu dostępności przycisku, a sprawdzian sięga po nią bez budowania dokumentu. Które pozycje są jeszcze brakiem jest wyeksportowane, żeby dało się to sprawdzić bez DOM.

## budowa/klient-poprzedni/src/moduly/browser/macierz-izolacji.ts
Punkty izolacji właściwe przeglądaniu obejmują dostęp sieciowy procesu sesji, kontenery tożsamości i zakres pętli wykonawczej. Panel pokazuje, co z nich obowiązuje dla okna przeglądarki, i tylko to. Zapisu tu nie ma: przełączniki i profile zapisuje okno konfiguracji punktów izolacji, a dwa miejsca zapisujące tę samą politykę dawałyby dwa różne zdania o tym, co obowiązuje. Nazwy punktów izolacji przychodzą ze słownika okna punktów izolacji — własny wykaz nazw w module znaczyłby, że ten sam przełącznik nazywa się w dwóch oknach inaczej. Stan przełącznika nigdy nie jest samą barwą: przy każdej pozycji stoi słowo „odcięty" albo „wspólny", bo panel czyta się także bez rozróżniania barw.

## budowa/klient-poprzedni/src/moduly/library/modul-library.ts
Układ wynika z roli okna: Library Explorer jest oknem wiodącym i jedynym
z własnym wejściem, a File Preview, Versioning Panel, Tags & Collections
i Metadata & Archive Panel odnoszą się do wykazu i do pliku wskazanego wyżej,
więc wskazanie w przeglądarce przestawia pozostałe cztery okna naraz, bo
stan jest jeden. Przekazanie z innego modułu ma odbiorcę w oknie wiodącym:
komenda przekazania rozgłasza zmianę okna, a nie zmianę pliku, i nie zakłada
pliku w repozytorium, więc nasłuch wyłącznie na zdarzeniu zmiany pliku
przekazania by nie zobaczył. Dziesięć komend obszaru biblioteki ma uchwyt
w rdzeniu, a cztery z nich — wgranie, dołożenie i przywrócenie wersji oraz
ustawienie etykiet — rozgłaszają po udanym wykonaniu zdarzenie zmiany pliku,
które jest jedyną drogą odświeżenia okien poza ich własnym działaniem.
Rodziny arsenału moduł bierze wyłącznie drogą ścieżki na dysku, co jest
rozstrzygnięciem mierzonym, nie ostrożnością, ponieważ wszystkie trzy
rodziny rozwiązują identyfikator zasobu przez repozytorium modułu
projektowego, a plik biblioteki leży w innym rejestrze i wróciłby stamtąd
odmową; przycisk wysyłający identyfikator zasobu biblioteki zawodziłby
zawsze, więc go nie ma. Kontrakt niesie jednak drugą drogę źródła, treść
wciąganą do magazynu pod sumą kontrolną, i tą drogą trzy czynności stoją
w obszarze archiwum: rozpoznanie, przetworzenie i spakowanie materiału. Ich
wynik jest zasobem magazynu projektowego, nie zasobem biblioteki, i widok
mówi to zdaniem, żeby nikt nie wziął jednego za drugie. Straż odmów i odmowa
nieznanej komendy zostają na miejscu, bo mierzą stan rdzenia w chwili
wywołania: gdyby uchwyt wypadł, okno powie to z pomiaru, a nie z komentarza.

## budowa/klient-poprzedni/src/moduly/apps/okno-mcp-connector-console.ts
Okno stoi w całości na czynnościach, których kontrakt nie prowadzi, i jest zbudowane właśnie po to, żeby ten brak był widoczny i policzony: kontrolki są obecne, klikalne i nazywają brakującą komendę, żadna nie jest wygaszona i żadna nie udaje wykonania. Powód każdej kontrolki bierze się z bytu pokrycia komend, a nie z napisu wpisanego tutaj — byt pyta rdzeń o wykaz jego komend i sam przerysowuje zdanie, więc w dniu, w którym rdzeń dostanie uchwyt dla którejś z tych komend, kontrolka powie to sama, bez tknięcia tego pliku; to ważne, bo brak jest stanem przejściowym, komendy są zaprojektowane do natychmiastowej implementacji. Formularz argumentów jest jednym polem tekstowym, nie polami wyprowadzonymi ze schematu, bo schemat przychodzi z odkrycia narzędzi, którego nie ma — pola zmyślone ze zgadniętego schematu wyglądałyby na wiedzę o serwerze, której okno nie ma. Okno otwiera się na pozycji wskazanej w Integrations Hubie; inspektor bez wskazanej integracji nie ma czego inspekcjonować i mówi to wprost.

## budowa/klient-poprzedni/src/moduly/library/odbior-przekazania.ts
Przekazanie zakłada w rdzeniu okno modułu i rozgłasza wyłącznie zmianę okna;
pliku w repozytorium nie zakłada, więc zmiana pliku po przekazaniu nie
przychodzi, a bez tej subskrypcji przybycie kompletu byłoby dla modułu nieme.
Zdarzenie zmiany okna niesie okno, nie powód jego zmiany, a biblioteka ma
także okno, w imieniu którego moduł sam działa. Przekazanie rozpoznają dwa
warunki: okno jest inne niż okno modułu, a okno modułu jest już znane —
dopóki pasek kontekstu nie odczytał własnego okna, moduł milczy; oraz rdzeń
ma dla tego okna komplet z dokumentami, bo okno, którego nikt nie
przekazywał, oddaje komplet bez dokumentów. Identyfikatory z kompletu są
identyfikatorami dokumentów modułu nadawcy i biblioteka może ich nie znać,
więc przycisk odświeża wykaz i mówi wprost, czy rdzeń ma pod tym
identyfikatorem plik: gdy ma, wskazuje go jako plik czynny, a podgląd,
wersje i etykiety przestawiają się razem; gdy nie ma, nazywa to brakiem
zamiast pokazać pustą pozycję. Powtórne naciśnięcie przycisku w trakcie
odczytu mówi, co się właśnie dzieje, zamiast wysyłać drugie żądanie o to
samo.

## budowa/klient-poprzedni/src/moduly/browser/modul-browser.ts
Układ modułu wynika z roli okna i z warstwy widoczności. Browser Window jest oknem wiodącym i punktem wejścia modułu, więc stoi w pasie pierwszym na całą szerokość i nie chowa się nigdy. Sources Panel, Notes Panel, Automation Studio, Capture & Monitor Panel oraz macierz izolacji są rozszerzeniami bocznymi: w spoczynku zwinięte, otwierane wyzwalaczem paska kontekstu albo skrótem klawiszowym. Stan jest jeden na cały moduł: zaznaczenie zrobione w podglądzie strony trafia do formularza notatki, źródło dodane w panelu źródeł pojawia się w wyborze powiązania notatki, a strona przechwycona w Capture & Monitor Panel przestawia wspólny podgląd, bo stan przeglądania jest jeden.

## budowa/klient-poprzedni/src/moduly/apps/okno-permissions-trust-center.ts
Okno otwiera się na pozycji wskazanej gdzie indziej — znacznikiem uprawnień karty App Catalogu albo przyciskiem wiersza Integrations Hubu — a wskazanie mieszka w stanie modułu. Granica tego okna jest ostra i okno musi ją wypowiedzieć, bo bez tego czytałoby się jak deklaracja bezpieczeństwa, której nikt nie złożył: pochodzenie pozycji jest znane i pokazane, lecz podpis cyfrowy i suma kontrolna pakietu nie są polami pozycji katalogu, więc oznaczenie zweryfikowanego wydawcy nie ma tu pokrycia i nie pada; wymagane uprawnienia leżą w manifeście rozszerzenia, którego kontrakt nie przenosi, więc okno pokazuje nieprzezroczystą konfigurację pozycji i mówi, że to nie jest to samo; osiem zakresów izolacji technicznej opisuje warunki, w jakich wykonuje się kod na tej platformie, i w tych samych warunkach wykona się kod rozszerzenia — nie jest to izolacja nadana pojedynczej pozycji. Zasada zero blokad obowiązuje tu wprost: żadne ostrzeżenie nie wstrzymuje instalacji ani włączenia, kontrola zostaje po stronie Operatora przez świadome włączenie i przez zakres uprawnień. Trzy oznaczenia pochodzenia wymagają weryfikacji podpisu, lecz kontrakt niesie samo pochodzenie, więc okno nazywa pochodzenie i mówi wprost, że o podpisie nie wie nic, zamiast wyprowadzać z pochodzenia orzeczenia o zaufaniu.

## budowa/klient-poprzedni/src/moduly/browser/modul-browser.ts — pasek uczciwości
Byt wspólny wypowiada okna, które katalog rdzenia modułowi przypisuje, a moduł ich nie buduje, oraz te, które moduł buduje, a rdzeń mu ich nie przypisał.

## budowa/klient-poprzedni/src/moduly/browser/modul-browser.ts — przerysowanie z katalogu
Element przerysowuje się także po odczycie cudzym, na przykład powłoki budującej nawigację. Treść wpisana na sztywno byłaby dokładnie tym, co pasek uczciwości ma tropić: twierdzeniem o stanie rdzenia wypowiedzianym bez zapytania rdzenia.

## budowa/klient-poprzedni/src/moduly/apps/okno-publisher-panel.ts
Jedna część okna stoi na kontrakcie i działa: dziennik wydań produktu składa się z wdrożeń oddanych przez odczyt wdrożeń — wersja i notatki wydania są polami wdrożenia, więc chronologia wydań produktu istnieje naprawdę; nie jest to jednak dziennik wydań rozszerzenia, i okno tę różnicę nazywa. Reszta okna czeka na komendy: formularz manifestu stoi i trzyma wpisane wartości, nie jest zaślepką, bo pola manifestu są znane z opracowania — tożsamość, wersja semantyczna, deklaracja narzędzi, wymagane uprawnienia, zależności — ale nie ma dokąd ich wysłać, i przycisk zapisu mówi to po naciśnięciu, zamiast milczeć albo udawać zapis. Powody kontrolek biorą się z bytu pokrycia, więc przerysują się same, gdy rdzeń dostanie uchwyty dla tych komend.

## budowa/klient-poprzedni/src/moduly/library/odczyty-biblioteki.ts
Wykaz plików zawężają pola takie jak etykieta, kolekcja i projekt, dopasowanie
słów pracuje w indeksie pełnotekstowym repozytorium, a dopasowanie znaczenia
w zakresie biblioteki pracuje we wskaźniku osadzeń; to trzy różne komendy
kontraktu i okno ma dla każdej osobne wejście. Tryb hybrydowy nie jest
czwartą komendą: okno wysyła obie i składa odpowiedzi, słowa przed
znaczeniem, bo dopasowanie dosłowne jest sprawdzalne, a semantyczne jest
przybliżeniem. Nieudany odczyt zostawia powód, nie pustkę, ponieważ pusty
wykaz po odmowie i pusty wykaz na świeżej instalacji to dwa różne stany.
Pusty wykaz w trybie semantycznym może znaczyć, że wskaźnika znaczenia nie
zbudowano, a nie że nie ma takich plików, dlatego zdanie wyniku nazywa
drogę, którą wynik powstał. Trafienie bez pliku w wykazie znaczy, że
wskaźnik zna dokument, którego wykaz nie oddał, bo wypadł poza granicę
odczytu albo zniknął z repozytorium po zbudowaniu wskaźnika. Kolejność dróg
w trybie hybrydowym jest rozstrzygnięciem, nie wygodą: dopasowanie słów da
się sprawdzić w treści pliku, dopasowanie znaczenia jest przybliżeniem
wskaźnika, a kontrakt nie niesie skali wspólnej dla obu indeksów, więc
liczba złożona z dwóch niewspółmiernych wyglądałaby na pomiar.

## budowa/klient-poprzedni/src/moduly/browser/modul-browser.ts — jeden odczyt na kanał
Tak samo działa wykaz komend rdzenia, z którego pozycje modułu biorą powód swojego bezruchu — odczyt idzie raz na kanał, nie raz na moduł.

## budowa/klient-poprzedni/src/moduly/browser/modul-browser.ts — odpięcie od katalogu
Warstwa adnotacji odpina obserwatora rozmiaru płótna z tego samego powodu, co pasek odpina się od wspólnego katalogu okien.

## budowa/klient-poprzedni/src/moduly/library/okno-file-preview.ts
Okno nie ma własnego wejścia i bez wskazania stoi w stanie pustym, nie
pytając rdzenia o nic. Przycisk zamknięcia nie zamyka okna komunikacji
sesji, tylko zdejmuje wskazanie pliku, a dymek okna mówi dlaczego. Odpowiedź
podglądu jest jedynym świadkiem tego, czy repozytorium ma treść pliku, dzięki
czemu wiersz w przeglądarce mówi o braku treści bez drugiego wywołania, a
werdykt bierze sam podgląd, nie powodzenie komendy, ponieważ rdzeń odpowiada
powodzeniem także wtedy, gdy zamiast treści oddaje odnośnik do niej. Bez
rozróżnienia odmowy od pustki okno przechodziłoby w stan gotowy i pokazywało
pusty prostokąt, po którym nie da się rozpoznać, czy plik jest pusty, czy
podgląd się nie udał. Nastawę modułu docelowego zmienia także drugie okno
z czynnościami zbiorczymi Explorera, a poniższe gałęzie kończą odświeżanie
wcześnie przy braku wskazania i przy pliku niezmienionym; postawiony niżej,
ster pokazywałby wartość sprzed tamtej zmiany i przeniesienie poszłoby gdzie
indziej, niż mówi.

## budowa/klient-poprzedni/src/moduly/browser/narzedzia-inspekcyjne.ts
Warstwa niesie pięć narzędzi: inspektor DOM i stylów, monitor sieci z eksportem HAR, konsolę strony, emulację urządzeń oraz podgląd źródła i różnic. Dwa ostatnie panel wykonuje sam z migawki: źródło strony przychodzi z polem zawierającym HTML, a porównanie zestawia dwie migawki tego samego okna. Trzy pierwsze wraz z emulacją mają własne komendy prowadzone protokołem narzędzi deweloperskich i stoją tu jako pozycje pytające rdzeń o pokrycie — ich powód bierze się z odczytu wykazu komend, nie z napisu w module. Porównanie zestawia wiersze treści, a nie znaki: pełny algorytm różnicowy należy do rdzenia, a okno ma powiedzieć, czy i o ile strona się zmieniła, oraz w którym wierszu zaczyna się różnica.

## budowa/klient-poprzedni/src/moduly/apps/przybornik-apps.test.ts
Pilnowane jest jedno, ale najważniejsze: czy z okna prowadzi droga do każdej komendy obszaru. Moduł ma czterdzieści jeden komend; sześć z nich prowadzą kontrolki formularzy okien architektury, warsztatu, wdrożenia i trzech odczytów, pozostałe trzydzieści pięć — narzędzia przyborników. Wykaz komend liczy się z kontraktu w czasie działania, nie z listy wpisanej w sprawdzianie: nowa komenda obszaru apps dołożona do kontraktu ma ten sprawdzian złamać, bo znaczy, że okno o niej nie wie. Drugi sprawdzian pilnuje zasady zero blokad: ani jedna kontrolka przybornika nie ma atrybutu wyłączenia — narzędzie bez pokrycia, brak okna modułu, puste pole wymagane, ma być klikalne i nazwać brak, a nie milczeć pod wyszarzonym przyciskiem.

## budowa/klient-poprzedni/src/moduly/browser/narzedzia-inspekcyjne.ts — czynności bez zastępczych przycisków
Drzewo DOM, rejestr sieciowy, konsola i emulacja urządzenia mają w rdzeniu uchwyty prowadzone protokołem narzędzi deweloperskich, a okno je wywołuje. Zastępczych przycisków bez obsługi już tu nie ma.

## budowa/klient-poprzedni/src/moduly/library/okno-library-explorer.ts
Cztery funkcje użytkownika z wiersza wykazu mają każda własną drogę do
rdzenia: nawigacja po strukturze, wyszukiwanie w trzech trybach, otwarcie
zasobu, gdzie wskazanie pliku przestawia cztery pozostałe okna, a otwarcie
w module źródłowym idzie osobną komendą przekazania, oraz wgranie pliku.
Prezentacja wykazu ma pięć postaci, wszystkie liczą się z tej samej
odpowiedzi rdzenia i nie wysyłają ani jednej komendy. Komplet przybyły
z innego modułu trafia do okna wiodącego, gdzie użytkownik może z nim
cokolwiek zrobić. Przeglądarka przycina pola przed wysłaniem, a rdzeń
biblioteki dopasowuje etykietę dosłownie, nie przycinając ani zapisu, ani
filtru; sam mechanizm przycinania przeglądarki zdejmuje przy tym znaki,
które rdzeń zostawia, więc etykieta zapisana z takim znakiem na brzegu jest
z tego pola nieosiągalna, a pusty wykaz wyglądałby jak zdanie o
repozytorium. Pusty wykaz w trybie semantycznym może znaczyć brak
wskaźnika znaczenia, a nie brak takich plików, dlatego o powodzeniu orzeka
faza wykazu, nie treść zdania. Przy frazie albo etykiecie w polu rdzeń
odpowiedział o tym, o co go pytano, a nie o całym zbiorze. Katalog modułów
obsadza ster modułu docelowego w tym oknie i w podglądzie pliku: jedna
nastawa, więc jeden odczyt.

## budowa/klient-poprzedni/src/moduly/apps/przybornik-apps.ts
Każde narzędzie jest przyciskiem, a przy narzędziach, które czegoś od Operatora potrzebują, obok przycisku stoją pola wejściowe; naciśnięcie woła rdzeń i pokazuje pod przyciskiem zdanie o skutku — nazwę bytu, liczbę pozycji, kod odpowiedzi, odwołanie do pliku — składane przez samo narzędzie, bo tylko ono wie, co w jego odpowiedzi jest skutkiem. Przybornik jest czymś innym niż wykaz braków: tamten wymienia czynności, których kontrakt nie niesie, a naciśnięcie mówi, czego brakuje; ten wymienia czynności, które kontrakt niesie i rdzeń obsługuje, więc naciśnięcie robi robotę — dwa wykazy obok siebie mówią Operatorowi wprost, co w tym oknie działa, a co jest jeszcze zapowiedzią. Zasada zero blokad: ani jedna kontrolka nie dostaje atrybutu wyłączenia; narzędzie, któremu brakuje okna modułu albo wypełnionego pola, pozostaje klikalne i po naciśnięciu nazywa brak, zamiast milczeć pod wyszarzonym przyciskiem. Odmowa rdzenia nie jest wyjątkiem widoku: wraca zwykłym wynikiem z polem błędu, a przybornik pokazuje ją tym samym zdaniem, którym opisuje odmowy reszta platformy. Wartość startowa pola wejściowego jest przykładem z domeny produktu, nie wartością wymuszoną.

## budowa/klient-poprzedni/src/moduly/library/okno-metadata-archive.ts
Zakładki obiecywałyby trzy równorzędne widoki jednego bytu, a tu Metadane
mówią o pliku wskazanym, a Archiwum i Higiena o całym repozytorium, dlatego
przełącza je selektor obszaru. Zależność wejściowa różni się obszarem:
Metadane bez wskazania pliku nie mają o czym mówić i okno stoi wtedy
w stanie pustym, a Archiwum i Higiena pracują na odczytanym wykazie
i wskazania nie potrzebują. Okno nie ma własnego odczytu poza przeliczeniem
wskaźnika znaczenia, bo kontrakt nie ma komendy metadanych zasobu, więc
wszystko, co panel pokazuje, pochodzi z wykazu przyniesionego przez
przeglądarkę plików i z odpowiedzi o treści odłożonej przez podgląd pliku;
jedno zdanie dla wszystkich trzech obszarów orzekałoby o wskazaniu pliku
także tam, gdzie wskazanie nie jest do niczego potrzebne. Cztery puste stany
przycisku pierwszej akcji prowadzą do czterech różnych czynności: odśwież
wykaz, dodaj plik, wskaż plik, przełącz obszar, a jeden przycisk o stałym
napisie kierowałby w trzech z nich w złą stronę, więc każdy stan mówi to
własnym zdaniem.

## budowa/klient-poprzedni/src/moduly/browser/okno-automation-studio.ts
Rozmowa z rdzeniem stoi w pliku `czynnosci-automatyk.ts`, widok tekstowy kroków w pliku `edytor-scenariusza.ts`, granice przebiegu w pliku `limity-przebiegu.ts`. Kroki scenariusza mają jedno miejsce: treść widoku tekstowego. Wykaz kart pod formularzem jest jego odczytem, a nie drugim zbiorem — dwa zbiory kroków rozjechałyby się przy pierwszej ręcznej poprawce, a operator nie wiedziałby, który z nich pojechał do rdzenia. Nagrywarka makra ma w kontrakcie własną komendę, której to okno jeszcze nie wywołuje — pozycja pyta więc rdzeń o jej pokrycie. Krok powstaje tymczasem z bieżącej migawki: przejście pod adres, który operator właśnie otworzył.

## budowa/klient-poprzedni/src/moduly/browser/okno-automation-studio.ts — kroki przyjęte
Karty kroków rysują się z ostatnio przyjętej treści, a nie z każdej litery wpisywanej w pole: treść w połowie poprawiona jest niepoprawnym zapisem JSON i wyczyściłaby wykaz kroków, których operator nie usuwał.

## budowa/klient-poprzedni/src/moduly/browser/okno-automation-studio.ts — nazwa komendy kroku
Krok z nazwą przepisaną ręcznie przeżyłby zmianę nazwy w kontrakcie i wracałby z rdzenia zdarzeniem obszaru nieznanego.

## budowa/klient-poprzedni/src/moduly/apps/stan-produktu.ts
Dwa równoległe stany dałyby dwie prawdy o tym samym produkcie: Architecture Designer definiuje komponenty, oba warsztaty przypisują do nich pliki, Product Builder rysuje z nich oś etapów, a Deployment Panel wdraża, wszystkie czytając z jednego stanu. Trzy komendy odczytu obszaru wypełniają stan tym, co rdzeń trzyma w bazie, zamiast tym, co przeleciało gniazdem w bieżącej sesji — bez nich odświeżenie okna przeglądarki zerowałoby moduł, choć historia wdrożeń, architektura i pliki warsztatu leżą w rdzeniu. Etapy budowy zostają wyłącznie przy zdarzeniu, bo komendy ich odczytu kontrakt nie niesie, więc pusty wykaz etapów na starcie jest stanem prawdziwym, a nie brakiem odczytu, i okno mówi o nim inaczej niż o wdrożeniach. Odczyt nie gasi stanu, który już jest: każdy z trzech odczytów jest niezależny, odmowa jednego zostawia dwa pozostałe nietknięte i nie kasuje tego, co moduł już wie; powód odmowy trafia do osobnego pola dla każdego odczytu, bo okna czytają je w różnych miejscach ekranu. Który z trzech odczytów obszaru rozstrzyga klucz powodu odmowy: nazwy są własne modułu, nie nazwami komend, bo warstwa widoku nazw komend nie zna.

## budowa/klient-poprzedni/src/moduly/library/okno-versioning-panel.ts
Dołożenie wersji jest tu potrzebne, bo wgranie pliku zakłada nowy dokument,
więc bez niego historia nie rośnie ponad jeden wpis, a przywrócenie nie ma
dokąd wracać. Filtr historii nie przyjmuje pola zawężającego w kontrakcie,
a raport składa się z odpowiedzi, którą okno już ma; wywóz historii wraz
z treścią każdej wersji nie powstaje, bo treści wersji niebieżącej nie
oddaje żadna komenda kontraktu. To nie jest repozytorium sesji ze Studia:
tam wersje żyją w toku sesji, tu w repozytorium biblioteki, i narastają przy
zmianie dokumentu w dowolnym module, także przy zmianie wykonanej przez
model, dlatego panel odświeża się także zdarzeniem zmiany pliku, a nie
wyłącznie własnym działaniem.

## budowa/klient-poprzedni/src/moduly/browser/okno-browser-window.ts — zdanie stanu pustego
Obie części zdania są konieczne. Sam opis okna nie mówi, czy przeszkodą jest brak sesji, brak okna przeglądarki czy tylko brak przejścia; sam powód z rdzenia nie mówi, czym okno jest ani jak je zapełnić. Bez okna przeglądarki nie ma o co pytać o migawkę, więc zdanie o migawce opisywałoby wtedy skutek zamiast przyczyny — operator ma dostać powód braku okna, bo to on rozstrzyga, co da się zrobić dalej.

## budowa/klient-poprzedni/src/moduly/browser/okno-browser-window.ts — wskaźnik obecności
Wskaźnik mówi o migawce i o chwili jej pobrania, a nie o samej „obecności" bez pokrycia. Bez migawki mówi wprost, że model nie ma na czym pracować.

## budowa/klient-poprzedni/src/moduly/browser/okno-browser-window.ts
Okno składa gotowe części: formularz nawigacji, podgląd strony, pasek zaznaczenia, pasek dolny i panel wyodrębnień. Trzy czynności operatora: nawigacja do strony przez formularz, przewijanie i zaznaczenie fragmentu przez podgląd. Nawigacja idzie do rdzenia; przewijanie ma komendę, której okno jeszcze nie wywołuje, a zaznaczenie dzieje się wyłącznie w kliencie i komendy nie potrzebuje. Odmowa rdzenia jest treścią okna, nie jego awarią. Każda komenda obszaru przeglądania może odmówić — okno pokazuje wtedy powód odmowy i mówi wprost, co odczytało naprawdę, zamiast pustego prostokąta.

## budowa/klient-poprzedni/src/moduly/browser/okno-browser-window.ts — kolejność warstwy adnotacji
Przełącznik nie może wskazywać czegoś, czego jeszcze nie ma. Zwrotne wywołanie zmiany trybu domyka pętlę: tryb zamknięty pływającym paskiem gasi przycisk paska dolnego.

## budowa/klient-poprzedni/src/moduly/assistant/panel-narzedzi.ts
Wykaz pochodzi z jednego katalogu, w którym rdzeń trzyma narzędzia, umiejętności i komendy akcji naraz; kontrakt rozróżnia je polem rodzaju i przedrostkiem źródła, nie osobną rodziną komend, więc osobna zakładka złożona z tej samej komendy byłaby drugim widokiem jednego wykazu, udającym drugie źródło. Zawężenie liczy rdzeń, nie okno: katalog liczy setki pozycji, a żądanie przyjmuje tekst, rodzaj, grupę i granicę wykazu, bo filtrowanie po stronie klienta wymagałoby ściągnięcia całości przy każdym naciśnięciu klawisza. Dołożenie idzie do karty sesji i żyje w jej stanie, definicji eksperta nie rusza; pole odpowiedzi katalogu mówi, co jest dołożone już teraz, więc przycisk wiersza nazywa czynność zgodnie ze stanem, który rdzeń oddał, a nie ze stanem zapamiętanym po ostatnim kliknięciu. Zakresu uprawnień i limitu wywołań pozycji okno nie udaje — nazywa brak.

## budowa/klient-poprzedni/src/moduly/library/stan-biblioteki.ts
Gdyby każde z pięciu okien prowadziło własny wykaz, etykieta trafiłaby na
inny plik niż ten pokazany w podglądzie, a panel wersji pokazywałby historię
trzeciego, dlatego zależności wykazu mieszkają w stanie, nie w oknach. Forma
prezentacji wykazu i tryb wyszukiwania też mieszkają tu: wybór widoku
zmienia to, co użytkownik widzi w przeglądarce, a tryb rozstrzyga, którą
komendą pójdzie następne szukanie, więc obie wartości muszą przeżyć
przerysowanie okna. Poza własnym działaniem stan odświeża wyłącznie
zdarzenie zmiany pliku: artefakt wytworzony w innym module dociera przez to
zdarzenie, a stan wciąga go tak samo jak własne wgranie, bez odpytywania
w pętli.

## budowa/klient-poprzedni/src/moduly/browser/okno-browser-window.ts — treść wskaźnika obecności
Stan wskaźnika obecności niesie napis, nie samą barwę, bo model widzi dokładnie treść migawki i dokładnie z tej chwili jej pobrania.

## budowa/klient-poprzedni/src/moduly/browser/okno-browser-window.ts — płótno adnotacji jako rodzeństwo
Płótno adnotacji jest rodzeństwem podglądu strony, nie jego dzieckiem, więc oznaczanie strony nie wstrzykuje w jej dokument ani jednego węzła.

## budowa/klient-poprzedni/src/moduly/assistant/edytor-makra.ts
Makro asystenta nie ma własnego magazynu w kontrakcie i opracowanie modułu prowadzi je tam, gdzie magazyn jest — przekazanie powtarzalnego makra lub rutyny do modułu Automations przyciskiem wysyłki. Kroki wpisuje się w JSON, bo dokładnie taki kształt niesie kontrakt; sprawdzenie jest tu, a nie w rdzeniu, z jednego powodu — błąd składni ma się nazwać przy polu, w którym powstał, zanim cokolwiek pojedzie do rdzenia. Sprawdzian pilnuje wyłącznie tego, co kontrakt uznaje za wymagane, identyfikatora kroku i jego rodzaju, reszty pól nie zgaduje; formatu YAML edytor nie przyjmuje i nie udaje, że przyjmuje. Pól opcjonalnych rozbiór kroków nie egzekwuje: rozstrzyga o nich rdzeń, a klient, który by je narzucił, odmawiałby definicji poprawnych.

## budowa/klient-poprzedni/src/moduly/browser/okno-browser-window.ts — panel rodzin
Rodziny prowadzone przez rdzeń obejmują karty, grupy kart, przestrzenie robocze, zakładki, przewinięcie, zrzut i narzędzia inspekcyjne.

## budowa/klient-poprzedni/src/moduly/browser/okno-browser-window.ts — ponowienie po odmowie
Ponowienie odczytujące tylko stan okna zostawiałoby operatora z komunikatem, którego nie da się zdjąć czynnością, którą mu podano — więc ponowienie sięga po stan okna i po treść strony naraz.

## budowa/klient-poprzedni/src/moduly/browser/okno-browser-window.ts — kolejność trzech stanów
Odwrotna kolejność kazałaby operatorowi czytać „jest pusto" w chwili, w której odczyt jeszcze trwa.

## budowa/klient-poprzedni/src/moduly/assistant/panel-rutyn.ts
Rutyna jest w kontrakcie automatyką z harmonogramem: odczyt wykazu automatyk mówi, jakie automatyki są, odczyt harmonogramu mówi, kiedy biegną, a zapis harmonogramu nadaje im cykliczność; osobnego bytu rutyny asystenta kontrakt nie ma i okno go nie zakłada. Wyzwalaczy zdarzeniowych okno nie ustawia: automatyka niesie cztery ich rodzaje, ale każdy wymaga wyrażenia właściwego dla swojego rodzaju — adres wywołania zdalnego, ścieżka pliku, warunek na wyniku modelu — a ich redakcja należy do modułu Automations, którego Workflow Builder jest miejscem budowy procesów. Assistant, zgodnie z granicą tematyczną modułu, inicjuje i nadzoruje pojedyncze zlecenia, a nie projektuje pełnych procesów.

## budowa/klient-poprzedni/src/moduly/library/widoki-wykazu.ts
Widok, który nie ma czego pokazać mimo niepustego wykazu, oddaje zdanie
w polu braku zamiast pustego prostokąta, bo galeria bez obrazów i
repozytorium bez plików to dwa różne stany. Wciągnięcie pliku bez pola
rodzaju treści do galerii na wszelki wypadek stawiałoby w niej dokumenty.
Znacznik czasu dodania w kontrakcie jest liczbą milisekund epoki,
a przeglądarka zna strefę czasową przeglądającego, więc doba jest jednostką
najmniejszą, którą da się nazwać bez ustawienia strefy. Widok mapy opiera
się na geolokalizacji z metadanych EXIF, której kontrakt dziś nie niesie
w żadnej postaci, a pozycja usunięta z przełącznika wyglądałaby na widok,
którego nigdy nie było.

## budowa/klient-poprzedni/src/moduly/browser/okno-browser-window.ts — pustka a odmowa
„Rdzeń nie ma jeszcze migawki tego okna" jest stanem pustym; nieudany odczyt treści strony jest odmową i ma się nią przedstawić.

## budowa/klient-poprzedni/src/moduly/assistant/panel-schowka.ts
Zakładka nie czyta schowka maszyny Operatora i nie udaje, że umie: schowek należy do tamtej maszyny, a rdzeń stoi na serwerze. Podział jest jawny i widoczny na ekranie — Operator wkleja skopiowaną treść w pole treści do zapamiętania, okno oddaje ją rdzeniowi, a rdzeń daje jej trwałość; powrót idzie tą samą drogą, kliknięcie wpisu kopiuje go do schowka karty przeglądarką, a nie rdzeniem. Nastawę skrótu globalnego trzyma rdzeń, klawisze przechwytuje powłoka programu okiennego; odpowiedź mówi wprost, czy rejestracji ma kto dokonać, a okno powtarza to zdanie, zamiast obiecywać skrót, który nikogo nie obudzi.

## budowa/klient-poprzedni/src/moduly/library/zapisy-zbiorcze.ts
Funkcje zapisu nie dotykają dokumentu i nie znają kontrolek — widok decyduje,
gdzie odpowiedź pokazać. Komplet etykiet zbudowany z nieświeżej kopii wykazu
kasowałby etykiety, których okno nie zdążyło zobaczyć, a rdzeń oddaje w tej
samej odpowiedzi plik po zapisie wraz z jego etykietami. Przycięcie po
stronie okna jest wypowiedziane, bo rdzeń nie przycina wcale, a wykaz
dopasowuje etykietę dosłownie, więc etykieta przycięta po cichu byłaby inna
niż wpisana i wykaz przestałby ją znajdować. Kod okna operacyjnego katalogu
i okno komunikacji to dwa różne byty, a brak okna jest powiedziany wprost,
a nie zamieniony w ciche nic. Rdzeń nie sprawdza przy przenoszeniu, czy
moduł docelowy istnieje: kod nieznanego modułu wraca powodzeniem wraz
z nowo założonym oknem, a że okno nie ma komendy sprawdzającej istnienie
modułu, nie orzeka o tym nic i podaje kod oraz numer okna oddane przez
rdzeń.

## budowa/klient-poprzedni/src/moduly/assistant/panel-wiedzy.ts
Odczyt wyszukiwania odnajduje fragmenty po znaczeniu i oddaje je wraz ze źródłem, żeby dało się je zacytować zamiast streszczać; budowa wskaźnika buduje wskaźnik, z którego to wyszukiwanie korzysta — rozdzielenie ich na dwa pliki dałoby dwa miejsca mówiące o jednym wskaźniku. Zakładka rozstrzyga wyłącznie o tym, co jest w panelu widoczne: pamięć semantyczna pokazuje szukanie, baza wiedzy budowę wskaźnika; wspólny jest zakres, bo wskaźnik jest jeden i szuka się w tym, co się zindeksowało. Wgrywania dokumentów tu nie ma i nie powinno być: pliki wchodzą do platformy przez moduł Library, a budowa wskaźnika obejmuje nim to, co w bibliotece już leży — drugie wejście dla plików znaczyłoby dwa repozytoria.

## budowa/klient-poprzedni/src/moduly/browser/okno-capture-monitor.ts
Rozmowa z rdzeniem stoi w pliku `czynnosci-materialu.ts`, pamięć zgromadzonego materiału w pliku `material-sesji.ts`. Panel mówi o trwałości prawdę: migawki zostają w rdzeniu, ale komendy odczytu wykazu wytworów okna kontrakt nie niesie, więc po przeładowaniu karty lista zaczyna się od nowa. Pobrania, kanały RSS, kolejka czytania i cykliczne sprawdzanie monitora mają w kontrakcie własne komendy, których panel jeszcze nie wywołuje — pytają rdzeń o ich pokrycie i mówią jego odpowiedź, zamiast orzekać o braku z napisu w module.

## budowa/klient-poprzedni/src/moduly/library/zrodlo-biblioteki.ts
Zbiór plików mieszka w stanie biblioteki, żeby cztery okna patrzyły na jeden
wykaz, a nie na cztery kopie. Odmowa rdzenia jest tu drogą równoprawną:
komenda, dla której rdzeń nie ma uchwytu, wraca zdarzeniem nieznanej
komendy, a straż zamienia je w wynik z błędem nazywającym żądany typ.
Wartość nazwy nieobjęta sprawdzianem kształtu potrafi dojść jako brak
wartości i wpisać się w zdanie potwierdzające jako nazwa kolekcji, dlatego
nazwa wchodzi do samego sprawdzianu.

## budowa/klient-poprzedni/src/moduly/assistant/panel-wybudzenia.ts
Rdzeń stoi na serwerze i mikrofonu tej maszyny nie widzi: nasłuch jest umową między oknem a rdzeniem — okno nagrywa u siebie, wysyła odcinki wraz z identyfikatorem okna, a rdzeń ogłasza, co w nich usłyszał, zdarzeniem częściowego rozpoznania oraz zdarzeniem wykrycia frazy wybudzającej; panel mówi to wprost, zamiast rysować mikrofon sugerujący, że rdzeń słucha sam. Odczyt stanu wykonalności oddaje niedostępność wraz z powodem, gdy wybudzenia nie da się wykonać, a panel powtarza powód i nie stawia przycisku obiecującego czynność, której nie ma czym wykonać.

## budowa/klient-poprzedni/src/moduly/browser/okno-capture-monitor.ts — panel rodzin
Rodziny prowadzone przez rdzeń obejmują monitory, kanały, kolejkę czytania, pobrania i wytwory sesji.

## budowa/klient-poprzedni/src/moduly/assistant/panel-zestawow.ts
Kontekst jest zestawem wskazań: usunięcie kontekstu kasuje wskazanie, a nie wpisy pamięci, i okno mówi to wprost przy przycisku, żeby Operator nie bał się sprzątać zestawów roboczych, a zarazem nie sądził, że kasuje ustalenia. Zasada retencji obejmuje zapisy kolejne: zapis zasady nie rusza wstecz wpisów zastanych, odpowiedź niesie ich policzoną liczbę i okno ją pokazuje, żeby Operator wiedział, ilu ustaleń zasada dotknie przy najbliższym wygaszaniu, zanim to nastąpi. Miernik okna kontekstu liczy żetony tokenizatorem rdzenia, a odpowiedź niesie nazwę słownika, którym policzono; pomiar bywa niewykonalny, wtedy okno pokazuje powód zamiast paska wobec granicy, której nikt nie ustalił.

## budowa/klient-poprzedni/src/moduly/library/zrodlo-otoczenia.ts
Okno komunikacji sesji bierze się z rejestru, bo przenoszenie kontekstu żąda
okna źródłowego, a moduł zna wyłącznie kody okien operacyjnych katalogu
rdzenia — to dwa różne byty. Kontekst okna jest wymagany przez wiersz
każdego z czterech okien wykazu, więc moduł pyta raz i dzieli odpowiedź.
Panel akcji nie jest zaszytym wykazem po stronie klienta: pozycje
przychodzą z rdzenia albo panel zostaje pusty i mówi to wprost, a wykonanie
pozycji bez uchwytu w rdzeniu wraca odmową nieznanej komendy. Katalog
modułów obsadza ster modułu docelowego: rdzeń przyjmuje każdy kod i zakłada
okno z dokładnie tym kodem, więc kodu wpisanego z ręki nie miałby kto
sprawdzić. Odczyt kompletu przeniesionego do okna jest jedyną komendą
kontraktu, która czyta magazyn zapisany przy przeniesieniu kontekstu; bez
niej przybycie przekazania jest dla modułu nieme. Zdarzenie przybycia nie
rozgłasza zmiany pliku i nie zakłada pliku w repozytorium — rozgłasza
wyłącznie zmianę okna z oknem docelowym.

## budowa/klient-poprzedni/src/moduly/assistant/pasek-polecenia.ts
Mikrofon nie jest wygaszony, choć kontrakt nie ma przesyłu dźwięku: naciśnięcie odpowiada zdaniem mówiącym, czego brakuje, i prowadzi ognisko do pola transkrypcji; rozpoznanie mowy jest warstwą wejścia, nie drugą drogą rozmowy — po transkrypcji treść wchodzi tam, gdzie weszłaby wpisana ręcznie. Wybudzenie stoi osobno od mikrofonu, bo to dwie różne czynności: mikrofon nagrywa jedno polecenie i wysyła je do rozpoznania, a wybudzenie prowadzi nasłuch ciągły i frazę, na którą asystent reaguje — sklejone w jedną kontrolkę dałyby jeden przycisk o dwóch znaczeniach.

## budowa/klient-poprzedni/src/moduly/multitasking/braki-kontraktu.ts
Powitanie rdzenia oddaje wykaz komend zarejestrowanych po montażu,
obsługiwanych naprawdę, a nie tylko wypisanych w kontrakcie. Moduł pyta o to
raz przy montażu i układa z odpowiedzi zdanie każdej nieczynnej kontrolki,
więc gdy rdzeń domknie kolejną komendę, zdanie zmienia się samo. Stany są
trzy, nie dwa: dopóki rdzeń nie odpowiedział, kontrolka nie orzeka o braku,
a odmowa powitania też nie jest orzeczeniem braku. Kontrolka bez pokrycia
nie znika i nie udaje, że działa — zostaje widoczna, nieczynna i niesie
powód wprost, żeby brak pozostał widoczny w oknie.

## budowa/klient-poprzedni/src/moduly/assistant/wysylka-polecenia.ts
Plik stoi poza oknem, bo okno składa kontrolki, a to jest rozmowa z rdzeniem. Polecenie idzie jednym wywołaniem z polem transkrypcji, tą samą drogą co polecenie wpisane ręcznie — nie ma tu drugiej drogi do rdzenia ani własnego modelu. Przerwanie dotyczy zlecenia, nie nagrania: kontrakt nie niesie strumienia dźwięku, ale niesie sterowanie anulowania; bez zlecenia w toku przycisk mówi wprost, czego po stronie audio brakuje, zamiast milczeć. Fazę nazywa wysyłka, nie okno: niesie wartość ze wspólnego słownika fazy okna, tego samego, którym mówią Actions Monitor i Activity Feed — nie jest to drugi mechanizm stanu, tylko ta sama faza, wpuszczana wprost w ten sam stan okna, nazwana w miejscu, które jako jedyne wie, co się właśnie stało; przełącznik trwa/nie trwa nie wystarcza, bo nie umie powiedzieć o odmowie rdzenia.

## budowa/klient-poprzedni/src/moduly/multitasking/formularz-powolania.ts
Formularz nie zgaduje za rdzeń: górną granicę powołania niesie kontrakt, więc
pole liczbowe ją pokazuje, ale zdanie potwierdzenia liczy podagentów
z odpowiedzi rdzenia, nie z tego, o ilu prosił formularz. Odmowa jest
nazwana: wiersz odpowiedzi niesie treść i kod wprost z rdzenia oraz zdanie
o tym, co można z tym zrobić. Formularz stawia wyłącznie pola, które
kontrakt zna.

## budowa/klient-poprzedni/src/moduly/assistant/zrodla-przekrojowe.test.ts
Sprawdzian pyta o jedno: czy każda komenda rodzin mowy, pamięci kontekstowej, retencji, zużycia kontekstu, schowka, skrótów, wywoływacza i asystenta ma drogę z okna do rdzenia. Wykaz oczekiwany nie jest tu przepisany — bierze się ze stałych kontraktu, a wykaz rzeczywisty z komend, które źródła naprawdę wysłały; komenda dołożona do kontraktu i pominięta w oknie wypadnie tu jako brak, bez dopisywania czegokolwiek w tym pliku. Sprawdzian mierzy warstwę kliencką, nie rdzeń: kanał jest próbny i tylko zapamiętuje nazwy — to wystarcza, bo pytanie brzmi, czy okno ma czym zawołać, a nie czy rdzeń odpowie, na to drugie odpowiadają sprawdziany skutku po stronie rdzenia.

## budowa/klient-poprzedni/src/moduly/browser/okno-sources-panel.ts
Panel ma jedną odpowiedzialność: formularz źródła i wykaz źródeł. Rozmowa z rdzeniem jest w pliku `czynnosci-zrodel.ts`, jeden wiersz w pliku `wiersz-zrodla.ts`. Wykaz pochodzi z rdzenia: polecenie odczytu źródeł oddaje źródła okna od najnowszego, więc panel pokazuje komplet zebrany w oknie, nie tylko pozycje dodane w tej karcie. Gdy rdzeń odmówi albo nie zna jeszcze okna, panel wypowiada powód zamiast pokazywać pusty wykaz bez wyjaśnienia.

## budowa/klient-poprzedni/src/moduly/browser/okno-sources-panel.ts — samoczynne dopisanie źródła
Bez progu identyfikatora migawki każde ogłoszenie stanu, także zmiana zaznaczenia, próbowałoby dopisać to samo źródło jeszcze raz. Ponownie odwiedzona strona nie zakłada drugiego wpisu, bo jej adres stoi już w wykazie okna — to jest reguła scalania duplikatów panelu.

## budowa/klient-poprzedni/src/moduly/multitasking/karta-przeplywu.ts
Kolumny żetonów i wywołań narzędzi stoją puste, ponieważ struktura podagenta
w kontrakcie nie niesie ani licznika żetonów, ani licznika wywołań
narzędzi. Komórki dostają znak braku wraz z powodem zamiast zera, które
czytałoby się jako wykonany pomiar, i z tego samego powodu wiersz
podsumowania nie sumuje żetonów.

## budowa/klient-poprzedni/src/moduly/browser/okno-sources-panel.ts — wypis bibliografii
Bibliografia wklejana ze schowka i bibliografia wczytywana do menedżera to dwa różne zastosowania tej samej treści, a rozdzielenie ich na dwie czynności kazałoby operatorowi wybierać, zanim zobaczy wynik.

## budowa/klient-poprzedni/src/moduly/assistant/zrodlo-assistant.ts
Plik odpowiada wyłącznie za warstwę wywołań kontraktu wraz ze sprawdzianem kształtu odpowiedzi; źródło nie ma własnego stanu i nie buduje ani jednego elementu — stan zleceń mieszka osobno, żeby trzy okna modułu patrzyły na jeden zbiór, a nie na trzy kopie. Żadne wywołanie nie rzuca wyjątkiem: niepowodzenie wraca polem błędu wyniku, a okno pokazuje je w swoim stanie błędu; tą samą drogą wraca odmowa merytoryczna rdzenia i koperta drogi bez uchwytu, więc okno nazywa je wprost zamiast udawać wykonanie. Wyróżnienie ma gdzie zamieszkać po stronie rdzenia, więc okno go nie udaje — powód bez znacznika byłby notatką do wpisu, którego nikt nie wyróżnił. Słuchacz zdarzenia dostaje kopertę, bo koperta niesie sesję okna zlecenia; bez niej zlecenie założone w innym oknie tej samej sesji byłoby nie do odróżnienia od zlecenia cudzej sesji.

## budowa/klient-poprzedni/src/moduly/browser/okno-sources-panel.ts — uwaga o pochodzeniu wykazu
Przy zawężonym wykazie wykaz krótszy od zebranego, bez zdania o liczbie pozycji ukrytych przez filtr, wyglądałby jak wykaz niepełny zamiast przefiltrowanego.

## budowa/klient-poprzedni/src/moduly/browser/okno-sources-panel.ts — kolejność trzech stanów wykazu
Odwrotna kolejność pokazywałaby pustkę w trakcie odczytu i po odmowie odczytu. Jedno nieudane odświeżenie nie unieważnia pozycji już wyświetlonych na ekranie.

## budowa/klient-poprzedni/src/moduly/multitasking/nadanie-rol.ts
Panel obsady odpowiada za scenę: zakłada okna, pokazuje skład i blokuje
drugi egzemplarz roli. Rolę nadaje komenda przypisania roli, nie zmiana
okna: ta druga komenda przy roli spoza kontraktu odpowiada powodzeniem,
a oddaje okno o roli samodzielnej bez przypięcia do koordynatora. Wcielenie
niesie wyłącznie komenda zmiany roli, a pasek jest jej jedynym wołaczem;
zmiana okna zostaje przy tytule i katalogach. Założenie okna z polem
przypięcia do koordynatora kończy się powodzeniem, ale nie utrwala więzi
w bazie, a wiersz okna zakładany przy tym zostaje w tym miejscu pusty, więc
obsada zostaje na niepełnym składzie, a okna wykonawców nie mają adresata
polecenia. Funkcja powtarza więc żądanie nadania roli zaraz po założeniu
okna i mówi w zdaniu wprost, że więź poszła drugim żądaniem; wołający
sprawdza obsadę odczytaną już po założeniu okna, więc gdy rdzeń więź
utrwali, drugie żądanie nie idzie wcale.

## budowa/klient-poprzedni/src/moduly/assistant/zrodlo-kontekstow.ts
Kontekst jest zestawem wskazań, nie właścicielem treści: usunięcie kontekstu kasuje wskazanie, a wpisy pamięci zostają, i okno mówi to wprost przy kasowaniu, żeby Operator nie bał się posprzątać zestawów roboczych. Zasada retencji obejmuje zapisy kolejne i nie rusza wstecz wpisów zastanych; odpowiedź niesie policzoną liczbę wpisów, których zasada dotknie przy najbliższym wygaszaniu — liczbę wierszy, nie oszacowanie. Pomiar zajętości okna kontekstu bywa niewykonalny, gdy kanał nie zadeklarował wielkości okna; to jest odpowiedź, nie awaria — pole dostępności niesie fałsz wraz z powodem, a okno pokazuje powód zamiast paska wobec granicy, której nikt nie ustalił.

## budowa/klient-poprzedni/src/moduly/browser/panel-rodzin.ts
Panel niesie komendy kart i przestrzeni roboczych, monitorów, kanałów, kolejki czytania, zakładek, pobrań, wytworów, zrzutów, narzędzi inspekcyjnych, nagrywarki makr i granic Wykonawcy. Panel nie pamięta niczego między naciśnięciami — stan modułu stoi w pliku `stan-przegladania.ts`, a wykaz wyniku jest odczytem, nie kopią. Powód jednego panelu zamiast kontrolki przy każdej pozycji: rodzin jest kilkanaście, a każda ma dwa–trzy pola. Rozsypane po oknach dałyby kilkadziesiąt kontrolek w miejscach, w których operator ich nie szuka; zebrane w sekcje przy właściwym oknie zostają w zasięgu jednego kliknięcia i nie zasłaniają pracy podstawowej. Panel mówi prawdę o odmowie: kod i treść odmowy rdzenia idą wprost do wiersza odpowiedzi. Cisza po naciśnięciu, albo zdanie „gotowe" bez pokrycia, byłaby tą samą szkodą, przed którą stoją sprawdziany skutku po stronie rdzenia.

## budowa/klient-poprzedni/src/moduly/assistant/zrodlo-mowy.ts
Do niedawna okno nie miało czym wysłać nagrania z mikrofonu: transkrypcja przyjmuje ścieżkę pliku na maszynie silnika, a nagranie z mikrofonu karty istnieje wyłącznie jako obiekt binarny w pamięci — przesył bajtów jest brakującym ogniwem, przyjmuje bajty i oddaje odnośnik, którym posługują się transkrypcja oraz polecenie asystenta; ta sama komenda otwiera dyktowanie w oknie komunikacji. Odsłuch jest drogą powrotną: bajty wracają do karty, która je wysłała, i wyłącznie dla odnośników, które rdzeń sam wystawił. Nasłuch ciągły nie jest otwarciem cudzego mikrofonu i nikt tu tego nie udaje: okno nagrywa u siebie, wysyła odcinki wraz z identyfikatorem okna, a rdzeń ogłasza, co w nich usłyszał, zdarzeniem częściowego rozpoznania oraz zdarzeniem wykrycia frazy wybudzającej. Żadne wywołanie nie rzuca wyjątkiem — niepowodzenie wraca polem błędu.

## budowa/klient-poprzedni/src/moduly/multitasking/okno-coordinator-chat.ts
Okno nie tworzy produktu końcowego, więc nie ma ani pola redakcyjnego, ani
zapisu wyniku. Koniec tury wykonawcy wybudza koordynatora po stronie
rdzenia; tu widać wyłącznie licznik obiegów, próg braku postępu i powód
zatrzymania odczytane ze stanu okna, a klient nie rozpoczyna obiegu i nikogo
nie wybudza. Powód, którym stan relacji da się odwrócić i skąd rdzeń to wie,
potrzebuje pełnego wiersza, bo plakietka nagłówka mieści tylko dwa słowa.
Wartość potwierdzająca subskrypcję bierze się z odpowiedzi, nie z faktu
wysłania żądania, bo rdzeń mógłby obserwacji nie założyć. Ocena treści
wyniku walidacji należy do przeglądarki wyników, która czyta stan procesów
osobną komendą monitorowania. Odczyt stanu wchodzi wywołaniem zwrotnym: sam
raport jest złożeniem wierszy, a sięgnięcie po rdzeń zostaje po stronie
wywołującego. Odczyt biegu koordynatora jest wyjęty poza źródło okien
i stan wspólny, bo nie dotyka niczego z wnętrza okna.

## budowa/klient-poprzedni/src/moduly/assistant/zrodlo-schowka.ts
Schowek należy do maszyny Operatora i rdzeń go nie czyta: Operator kopiuje u siebie, okno oddaje skopiowaną treść rdzeniowi, a rdzeń daje jej trwałość — historia przestaje ginąć razem z kartą i jest ta sama na każdej maszynie tego samego Operatora; wklejenie jest ruchem powrotnym i wykonuje je okno, u siebie. Nastawę skrótu globalnego trzyma rdzeń, przechwycenie klawiszy należy do powłoki programu okiennego; odpowiedź mówi wprost, czy rejestracji ma kto dokonać, a okno powtarza to zdanie zamiast obiecywać skrót, który nikogo nie obudzi.

## budowa/klient-poprzedni/src/moduly/automations/automations.test.ts
Sprawdziany obejmują wyłącznie byty rozstrzygalne bez rdzenia: przekład wzorca cykliczności na zapis cron i z powrotem, walidację definicji, miary i zawężenie wykazu przebiegów, układ warstwowy grafu, zapisy eksportu, kalendarz oraz wyjęcie scenariusza z przeniesionego kompletu. Dane przykładowe pochodzą ze świata produktu — kroki i przebiegi automatyki raportowej — i są oznaczone jako przykładowe nazwą. Zapora na powtórzenie usterki: wykazy przepisane ręcznie przestały raz odpowiadać kontraktowi po jego rozszerzeniu, a okno pokazywało wtedy mniej rodzajów i działań, niż silnik naprawdę zna, nie mówiąc o tym ani słowem; kompletność wymusza już typ mapy zupełnej po wyliczeniu, sprawdzian pilnuje drugiej połowy — że wykaz na ekranie powstaje z tej mapy, a nie obok niej.

## budowa/klient-poprzedni/src/moduly/browser/pasek-dolny.ts
Rozmowa z rdzeniem jest w pliku `czynnosci-paska-dolnego.ts`, podgląd i pasek zaznaczenia osobno. Przyciski są trzech rodzajów. Podział ekranu, tryb czytnika i tryb adnotacji dzieją się w kliencie — adnotacja rysuje na płótnie nad sceną, a do rdzenia idzie dopiero jej wynik przez komendę wysyłania wiadomości. Zrzut ekranu i tłumaczenie mają uchwyt w rdzeniu. Zakładka, makro i pobrania mają w kontrakcie własne komendy, których to okno jeszcze nie wywołuje — pytają więc rdzeń o ich pokrycie i mówią jego odpowiedź.

## budowa/klient-poprzedni/src/moduly/browser/pasek-dolny.ts — przełącznik trybu adnotacji
Tryb adnotacji zamyka się także przyciskiem zamknięcia trybu na pływającym pasku, a wtedy przycisk paska dolnego musi przestać twierdzić, że tryb trwa.

## budowa/klient-poprzedni/src/moduly/multitasking/okno-executor-chat.ts
Okno nie ma sterowania kolejką ani planu etapów. Identyfikator koordynatora
okna wykonawczego jest jedynym oznaczeniem więzi, więc okno wypisuje je
wprost, inaczej wykonawca własnej pary jest nie do odróżnienia od cudzego.
Tryb współpracy jest wspólny obu wykonawcom i rozstrzyga, kto dostaje
zlecenie przy przekazaniu, więc zapisuje się na oknie koordynatora, a nie
osobno w każdym wykonawcy, bo dwa zapisy tej samej rzeczy rozjeżdżałyby się
przy pierwszej zmianie. Numer wykonawcy spoza zakresu dawałby kod pusty,
czyli okno zbudowane bez wiersza katalogu, o czym nikt by nie zameldował.
Zapis wiadomości przez rdzeń zwraca zapisaną wiadomość użytkownika i nic nie
orzeka o turze modelu. Wywołanie zmiany stanu woła się także na fragment
strumienia, więc odczyt bezwarunkowy zasypałby rdzeń wywołaniami wykazu
podagentów w tempie strumienia.

## budowa/klient-poprzedni/src/moduly/browser/pasek-dolny.ts — nazwy rodzajów wyodrębnienia
„Treść renderowana" i „źródło strony" brzmią podobnie, a dają dwie różne rzeczy — dlatego opis przy pozycji rozstrzyga, co wyjdzie z wyodrębnienia.

## budowa/klient-poprzedni/src/moduly/automations/droga-komend.test.ts
Wykaz oczekiwany bierze się z kontraktu, nie z tego pliku: komenda dołożona do którejkolwiek z trzech rodzin i pominięta w źródle wypadnie tu jako brak, bez dopisywania czegokolwiek w sprawdzianie — sprawdzian nie mierzy więc tego, co ktoś pamiętał, tylko to, czego rodzina naprawdę wymaga. Przelot woła każdą czynność źródła raz i nie sprawdza jej wyniku, od tego są sprawdziany skutku po stronie rdzenia, tylko to, że okno ma czym daną komendę wysłać; moduł, który wygląda na kompletny, a nie umie wysłać jednej komendy z czterdziestu dziewięciu, jest modułem niekompletnym w miejscu, którego nie widać.

## budowa/klient-poprzedni/src/moduly/browser/pasek-dolny.ts — nazwa nastawy sterownika
Nazwa nastawy idzie do etykiety dostępności sterownika, bo w pasku narzędzi nazwy rodzajowe zabierają miejsce i nic nie mówią same z siebie.

## budowa/klient-poprzedni/src/moduly/browser/pasek-dolny.ts — bez dublowania czynności
Zakładka, nagrywarka makr i menedżer pobrań mają w rdzeniu uchwyty, a okno je wywołuje: czynności stoją w panelu rodzin przy oknie, do którego są przypisane.

## budowa/klient-poprzedni/src/moduly/automations/edytor-krokow.ts
Kolejność zmieniają przyciski w górę i w dół, a nie przeciąganie myszą: kontrakt niesie ją liczbą pola porządkowego kroku, a przyciski są dostępne z klawiatury. Mapa nazw rodzajów jest zupełna po wyliczeniu, nie wykazem przepisanym ręcznie: gdy rodzaj kroku urośnie w kontrakcie, kompilacja zatrzyma się tutaj i nowy rodzaj dostanie nazwę, zamiast zniknąć z pola wyboru bez śladu.

## budowa/klient-poprzedni/src/moduly/multitasking/okno-results-analyzer.ts
Przycisk stanu monitora melduje stan procesów telemetrii oddany przez
rdzeń. Zgłoszenie niezgodności idzie do okna koordynatora zwykłym zapisem
wiadomości. Wcielenie roli analityka niesie własny wykaz kryteriów oceny
i własny nagłówek zgłoszenia; wartość utrwala zapis ustawienia na poziomie
okna analityka.

## budowa/klient-poprzedni/src/moduly/browser/plotno-adnotacji.ts
Pasek narzędzi jest osobno, a złożenie jednego z drugim w warstwę nad podglądem stoi w pliku `warstwa-adnotacji.ts`. Płótno nie dopisuje ani jednego węzła do elementu podglądu, nie zmienia mu klas i niczego z niego nie czyta — leży nad nim. Tło jest przezroczyste, bo zrzutu strony nie ma: rdzeń pobiera stronę biblioteką HTTP i zostawia pole zrzutu puste. Rysunek powstaje więc po tym, co pokazuje scena, i nie utrwala tła. Kontekst rysunkowy oddaje wartość pustą w środowisku bez rasteryzacji — wtedy nie ma czego wykreślić, ale wykaz śladów, wybór narzędzia i czyszczenie działają nietknięte.

## budowa/klient-poprzedni/src/moduly/browser/plotno-adnotacji.ts — rozwiązanie barwy żetonu
Płótno rysunkowe nie zna zmiennej stylu wprost, więc żeton trzeba rozwiązać z wyliczonego stylu elementu. Gdy arkusz nie jest wczytany, w sprawdzianie bez stylów, wraca się do koloru tekstu elementu, a nie do barwy zapisanej wprost w kodzie — żadna z nich nie ma prawa tam paść.

## budowa/klient-poprzedni/src/moduly/browser/plotno-adnotacji.ts — bufor po zmianie rozmiaru
Zmiana rozmiaru zeruje bufor płótna, więc obraz odtwarza się z wykazu śladów. Obserwator rozmiaru bywa nieobecny w środowisku sprawdzianu, a jego brak nie jest powodem, by tryb adnotacji przestał działać.

## budowa/klient-poprzedni/src/moduly/browser/plotno-adnotacji.ts — spłaszczenie płótna
Gdy przeglądarka nie oddała obrazu, okno ma o tym powiedzieć, a nie wysłać pustkę zamiast rysunku adnotacji.

## budowa/klient-poprzedni/src/moduly/multitasking/panel-orkiestracji.ts
Środowisko nie udostępnia modułów w bocznej nawigacji, więc okien do modułu
nie przypina. Sekcja Role montuje istniejącą scenę czterech okien
roboczych, wraz z panelem podagentów w oknie wykonawcy; dwa wystąpienia
wykonawcy to jeden typ okna w dwóch wystąpieniach, nie dwa osobne okna.
Pozostałe cztery sekcje pracują na oknach modułu automatyzacji, bo własnych
okien nie mają, i każda niesie podsekcję powiązania mówiącą, czy powiązanie
jest skonfigurowane i jak je założyć, zamiast rysować pustkę albo udawać
własne okno. Sekcja nieznana temu plikowi nie gaśnie i nie znika: przestrzeń
robocza pokazuje wtedy stan pusty z nazwą sekcji, a brak ustawienia układu
podsekcji znaczy układ domyślny, nigdy niedostępność. Adresem układu jest
okno, a panel orkiestracji należy do środowiska; sesja bez ani jednego okna
zostawia układ miejscowy, o czym powłoka mówi przy pierwszej zmianie.

## budowa/klient-poprzedni/src/moduly/automations/indeks.ts
Moduł nie ma okna modułowego w żadnym środowisku: jest komponentem własnym strefy drugiej strony głównej, więc jego okna otwierają się stamtąd, a nie z przestrzeni roboczej karty sesji, dlatego złożenie przyjmuje automatykę i kolejkę, nie środowisko ani kartę sesji. Układ wynika z ról okien: Workflow Builder jest kreatorem budującym strukturę i stoi w obszarze głównym jako punkt wejścia, obok niego drugi kreator, Orchestrator, który tę strukturę układa; pod nimi pas wykonania — dwaj zarządcy, Scheduler i Queue Manager, i monitor, Execution Monitor — bo dopiero po zbudowaniu automatyki jest co uruchamiać i co obserwować. Okno rozmowy modułu nie należy do tego złożenia: jest bytem sesji i składa je warstwa rozmowy. Nad układem stoi pas przywołania: menu skraca drogę do okna do jednego wskazania i niczego nie chowa — siatka zostaje taka sama, żadne okno nie znika i żadne nie dochodzi z urzędu. Wskazanie w nawigacji przestawia się także przy skoku z okna, bo uchwyt ma nieść okno bieżące, a nie ostatnie wybrane w menu.

## budowa/klient-poprzedni/src/moduly/browser/skutek-zapisu.ts
Plik stoi osobno od plików czynności, bo tamte prowadzą rozmowę z rdzeniem, a ten ocenia jego odpowiedź. Zdanie zbudowane z flagi żądania twierdziłoby o rzeczy, której w odpowiedzi może nie być. Rdzeń pobiera dziś stronę biblioteką sieciową i zostawia pole zrzutu puste; zdanie czyta to pole, więc mówi prawdę i dziś, i po dobudowie zrzutów. Tak samo z wartością pola tekstowego — funkcja przycinania białych znaków w przeglądarce i odpowiednik rdzenia to dwie różne funkcje o różnym zestawie usuwanych znaków. Adres, kod i nazwa jednoznaczna jadą do rdzenia jako identyfikator, a o tym, co się zapisało, rozstrzyga wyłącznie odpowiedź.

## budowa/klient-poprzedni/src/moduly/browser/skutek-zapisu.ts — pozycja materiału sesji
Przechwycenie zamawia i zrzut ekranu, i źródło strony, a wraca z tym, co rdzeń ma. Bez zdania o brakującym zrzucie operator zobaczyłby w wykazie „archiwum" i nie wiedziałby, czemu nie „zrzut".

## budowa/klient-poprzedni/src/moduly/automations/panel-dozoru.ts
Trzynaście czynności szuflad i paneli popover okna przebiegów obejmuje przeglądarkę logów, drążenie do poziomu kroku, punkty wznowienia, wznowienie od punktu, podgląd i odtworzenie ładunku, reguły alarmowania, budżety czasu, skarbiec poświadczeń i dziennik audytu. Execution Monitor pokazuje historię wielu uruchomień naraz, a te czynności dotyczą jednego z nich — podstawienie przebiegu bieżącego kazałoby zgadywać, którego. Skarbiec i audyt stoją tu, bo obie rodziny przenikają okno przebiegów: poświadczenie jest tym, czego krok potrzebował, żeby przebieg się powiódł, a audyt odpowiada na pytanie, kto zmienił definicję między dwoma przebiegami; osobnego okna dla nich dokument projektowy nie zna, obie rodziny wprost nie mają odrębnych okien.

## budowa/klient-poprzedni/src/moduly/browser/skutek-zapisu.ts — liczba kroków automatyki
Liczba kroków ma znaczenie osobne od samego zapisu: samo słowo „zapisano" nie pokazuje, że scenariusz przyjęty z krokami odrzuconymi nic nie zrobi.

## budowa/klient-poprzedni/src/moduly/automations/panel-nadzoru.ts
Cztery z pięciu czynności żądają harmonogramu, nie automatyki — kontrakt wskazuje go osobnym polem; panel pyta więc o niego wprost i nie podstawia w jego miejsce automatyki bieżącej, bo harmonogram ma własny identyfikator i to on jest przedmiotem tych czynności. Piąta, adres webhooka, żąda automatyki i tę bierze ze stanu modułu. Klucz podpisu nie wraca i nie może wrócić: odpowiedź niesie referencję, a nie wartość — panel pokazuje dokładnie to, co dostał, i niczego nie dopowiada.

## budowa/klient-poprzedni/src/moduly/browser/skutek-zapisu.ts — pytanie w oknie rozmowy
Polecenie wysyłania wiadomości oddaje wiersz wiadomości wraz z oknem; gdyby wiadomość trafiła do cudzego okna, zdanie o „przyjęciu przez okno rozmowy" byłoby prawdą o czymś innym niż to okno.

## budowa/klient-poprzedni/src/moduly/multitasking/panel-podagentow.ts
Rdzeń rozgłasza zmianę podagenta przy powołaniu, wejściu w bieg, zakończeniu
i zatrzymaniu, a przycisk odświeżenia wykazu zostaje, bo zdarzenie mówi
o zmianie, nie o stanie zastanym, i po ponownym nawiązaniu łącza wykaz
trzeba odczytać raz od nowa. Panel wykonawcy nie ma odświeżać się na cudzym
podagencie. Sekcja niesie wyżej definicje podagentów i kontrolki wytwórni,
zanim żywy wykaz się otworzy. Rdzeń oddaje osobno podagentów, których
zatrzymał, i tych, którzy w chwili wywołania nie pracowali; ten drugi
przypadek nie jest błędem i panel go tak nie nazywa.

## budowa/klient-poprzedni/src/moduly/automations/panel-wersji.ts
Przebieg próbny sprawdza definicję automatyki przed wpięciem, wersje ją cofają i porównują, zmienne i mapowania opisują przepływ danych między jej krokami, notatka i położenie należą do jej węzłów, a szablon jest jej odbiciem zapisanym do wielokrotnego użycia. Automatyka bierze się ze stanu modułu, nie z pola panelu: pięć okien pracuje nad jedną automatyką naraz, panel z własnym polem automatyki pozwalałby zapisać wersję jednej, patrząc przy tym na kanwę drugiej. Odpowiedź rdzenia niesie automatykę po zmianie i to ona jest prawdą — panel pokazuje ją w całości pod przyciskiem.

## budowa/klient-poprzedni/src/moduly/browser/skutek-zapisu.ts — adnotacja dołączona do rozmowy
Polecenie wysyłania wiadomości przepisuje załączniki do zakładanej wiadomości i oddaje ją w odpowiedzi, więc rozbieżność jest widoczna wprost: wiadomość przyjęta bez obrazu znaczy, że w rozmowie stoi sam opis adnotacji, a zdanie „adnotacja dołączona" mówiłoby o rysunku, którego rdzeń nie zapisał.

## budowa/klient-poprzedni/src/moduly/browser/skutek-zapisu.ts — nazwa czynności przeniesienia
Czynność nazywa, co przeniesiono, na przykład „Przekazano 3 źródła", bo tego rdzeń nie oddaje: odpowiedź przeniesienia kontekstu niesie okno docelowe i znacznik przeniesienia, nie zawartość kompletu. Zdanie mówi więc osobno, co wysłano, i osobno, co potwierdził rdzeń.

## budowa/klient-poprzedni/src/moduly/multitasking/panel-subagent-network.ts
Definicje podagentów dostępnych w sesji stoją w obszarze narzędzi
konfiguracji sesji, odczytywanym na poziomie sesji, a panel zestawia ten
wykaz z górnym limitem podagentów na wykonawcę. Powołanie, wykaz biegnących
i zebranie wyników wymagają okna wykonawcy, a to panel dostaje wyłącznie
w polu opcjonalnym. Gdy wołający je podał, panel stawia formularz powołania
oraz żywy wykaz z odświeżeniem i zbieraniem; gdy nie podał, nie ma
identyfikatora okna, więc nie ma czego wołać, i w to samo miejsce wchodzą
kontrolki nieczynne z powodem wziętym z wykazu komend oddanego przez
rdzeń. Definicje jako pole nieokreślonego kształtu należą do dostawcy
kanału, nie do kontraktu, dlatego widok bierze zapis surowy zamiast udawać,
że definicji nie ma.

## budowa/klient-poprzedni/src/moduly/browser/slady-adnotacji.ts
Plik zawiera model rysunku i wykreślenie go pisakiem; elementy dokumentu i zdarzenia wskaźnika należą do pliku `plotno-adnotacji.ts`. Dzięki temu rysunek da się sprawdzić bez przeglądarki, w której nie ma kontekstu rysunkowego. Piksele płótna giną przy każdej zmianie jego bufora, na przykład przy zmianie rozmiaru okna; z wykazu śladów odtwarza się obraz po każdym przerysowaniu, a polecenie „Wyczyść adnotacje" opróżnia wykaz, nie zamalowuje płótno. Punkt zapamiętany bezwzględnie wyskoczyłby poza rysunek, gdy okno zwęzi się między narysowaniem a przerysowaniem; ułamek trzyma oznaczenie tam, gdzie operator je postawił względem oglądanej ramki.

## budowa/klient-poprzedni/src/moduly/automations/okno-wykaz-petli.ts
Przycisk uruchomienia nie pyta o potwierdzenie ani o wskazanie kolejki — zakłada kolejkę i posuwa ją działaniem start; pętla wyłączona ma ten sam przycisk co czynna, jej stan stoi w opisie pozycji. Zawężanie zastępuje blokowanie wierszy: napis szukania skraca wykaz po stronie klienta, bo kontrakt nie ma pola zapytania, a przełącznik tylko czynne zawęża samo żądanie polem stanu włączenia. Uruchomienie przestawia wspólny stan modułu na tę pętlę i jej kolejkę, więc Execution Monitor pokazuje jej przebiegi, a Queue Manager posuwa jej kolejkę bez przepisywania identyfikatorów. Nad wykazem stoi wykaz przekazań: scenariusze przeniesione tu z innych modułów komendą przeniesienia kontekstu; okno jest wejściem do modułu, więc to tutaj widać, że ktoś coś oddał — przekazanie zostaje propozycją do chwili, w której Operator naciśnie zapis, nic nie wchodzi do magazynu automatyk samo. Zdanie potwierdzenia po odczycie wykazu podaje liczbę oddanych pozycji, żeby wykaz krótki dał się odróżnić od przyciętego granicą; odczyt przekazań mówi, ile okien sprawdzono i ile z nich niesie scenariusz — pustka po sprawdzeniu jest odpowiedzią, nie usterką.

## budowa/klient-poprzedni/src/moduly/multitasking/panel-zadan-w-tle.ts
Każdy przepływ dostaje kartę z nazwą, odznaką stanu, wierszem podsumowania,
etapem z telemetrii i tabelą agentów. Panel nie stawia pauzy ani usuwania
podagenta, bo kontrakt nie niesie komendy, która by je wykonała: rodzina
komend podagenta obejmuje tylko powołanie, wykaz i zebranie wyników,
a komenda czynności kolejki wymaga identyfikatora kolejki, którego podagent
nie niesie.

## budowa/klient-poprzedni/src/moduly/browser/slady-adnotacji.ts — oprawa wykreślenia
Płótno rysunkowe nie zna zmiennej stylu ani skrótu czcionki z arkusza i przyjmuje wyłącznie wartości gotowe. Rozwiązanie żetonu barwy należy więc do warstwy, która ma element w dokumencie, a nie do samego rysowania.

## budowa/klient-poprzedni/src/moduly/multitasking/sekcja-harmonogram.ts
Cały mechanizm harmonogramu w kontrakcie to para komend: jedna zapisuje
regułę, druga ją oddaje, a okno robocze jest cudze, więc sekcja ustawia
regułę i wskazuje, gdzie ona potem pracuje. Przełącznik obowiązywania jest
stanem danych, nie bramką: harmonogram wyłączony zostaje w rdzeniu i włącza
się jednym naciśnięciem. Definicja wyzwalacza niesie wyłącznie rodzaj
i wyrażenie, a żadna komenda nie oddaje wykazu wyzwalaczy osobno.

## budowa/klient-poprzedni/src/moduly/automations/panel-dobudowy.ts
Cztery panele — wersje, nadzór harmonogramu, zlecenia kolejki, dozór przebiegu — mają tę samą budowę, a pod nimi jeden nośnik stanu treści, zamiast czterech razy tego samego rusztowania. Panel nie jest oknem operacyjnym i nie udaje nim być: okien modułu jest pięć, a panel osadza się wewnątrz okna, do którego należy jego praca — wersje w Workflow Builderze, zlecenia w Queue Managerze; szósty kafel na siatce byłby szóstym oknem, którego dokument projektowy nie zna. Odpowiedź rdzenia pokazuje się w całości, zapisem strukturalnym: panel jest powierzchnią roboczą Operatora nad czynnościami, których kontrakt oddaje bardzo różne kształty — wykaz wersji, różnicę pól, ładunek kroku, punkty wznowienia; rysunek zmyślony osobno dla każdej z nich pokazywałby mniej, niż rdzeń oddał, a to jest gorsze niż surowy zapis, Operator ma widzieć odpowiedź, a nie jej streszczenie napisane przez okno. Jedno miejsce wykonania czynności wystarcza, bo wszystkie czterdzieści trzy czynności kończą się tak samo i różnią się wyłącznie zdaniem. Pole nieobecne w żądaniu znaczy co innego niż pole o wartości pustej, a zapis nieczytelny sprawia, że wołający odmawia wysłania — żądanie z uszkodzonym ładunkiem odbiłoby się od rdzenia komunikatem o kopercie, a Operator ma zobaczyć, że to on pomylił nawias.

## budowa/klient-poprzedni/src/moduly/multitasking/sekcja-kolejki.ts
Komendy tworzenia, sterowania i wykazu kolejki obsługują pętlę sesyjną,
moduł Automations i to środowisko; sekcja nimi steruje, zamiast budować
własnego wykonawcę zleceń, a po skonfigurowaniu powiązania kolejkę prowadzi
Queue Manager modułu Automations. Przycisk łączący dwie akcje w jedną
robiłby co innego, niż mówi jego napis, dlatego brakujące czynności sekcja
tylko wypisuje. Kolejka niesie sesję, okna, stan i licznik obiegów, ale pola
zasięgu nie ma, więc zasięg jest wykazem opisowym wraz z kształtem
brakującego pola, a nie kontrolką bez skutku.

## budowa/klient-poprzedni/src/moduly/automations/kalendarz-uruchomien.ts
Wykaz pięciu najbliższych terminów mówi kiedy najbliżej; kalendarz mówi jak gęsto i pokazuje to razem z historią, dzięki czemu widać dzień, w którym uruchomienie wypadało, a przebiegu nie było. Terminy zaplanowane liczy okno z wpisanej cykliczności, bo dotyczą reguły jeszcze niezapisanej; terminem obowiązującym pozostaje najbliższe uruchomienie liczone przez rdzeń — kalendarz jest podglądem i mówi to wprost w swoim podpisie. Dni układają się według czasu miejscowego przeglądarki, bo tak Operator czyta kalendarz; rachunek terminu idzie w UTC, tak samo jak w rdzeniu.

## budowa/klient-poprzedni/src/moduly/browser/stan-przegladania.ts
Migawkę odświeżają dwie drogi: własne wywołanie odczytu treści strony i zdarzenie zmiany strony — odpytywania w pętli nie ma. Wykazy źródeł i notatek oraz migawka idą zaraz po ustaleniu okna, bo ich treść mieszka w rdzeniu, nie w pamięci karty. Stan warstw widoczności stoi tutaj razem z resztą, a nie osobno przy pasku kontekstu: odsłonięcie rozszerzenia jest zmianą, na którą okna reagują tak samo jak na nową migawkę, więc idzie tym samym ogłoszeniem.

## budowa/klient-poprzedni/src/moduly/multitasking/sekcja-monitor.ts
Przełącznik trybu nakładki rozstrzyga, co nakładka może zrobić z procesem:
obserwator wyłącznie patrzy, podpowiada i ostrzega, a operator dodatkowo
zatwierdza i wstrzymuje kroki z dowolnego miejsca. Stan nakładki podaje
urządzenie, sesję, okno, procesy przypięte i licznik procesów w biegu, ale
komendy ustawiającej tryb nie ma wcale, dlatego wartość operatora przy
rdzeniu, który nic nie zapisał, mówiłaby nieprawdę o tym, kto może
wstrzymać krok. Statusy idą na żywo: subskrypcja monitora zapisuje okno na
telemetrię i oddaje stan bieżący jednym ruchem, a zdarzenie postępu donosi
zmiany, choć nie niesie etykiety ani sesji, więc widok złożony z samych
zdarzeń gubiłby nazwy przebiegów po pierwszej zmianie. Okna analizy wyników
i monitora wykonania należą do tej sekcji, ale są osobnymi oknami produktu,
więc panel nie buduje trzeciego monitora.

## budowa/klient-poprzedni/src/moduly/browser/stan-przegladania.ts — zdarzenie zmiany strony
Jedna sesja bywa oglądana w kilku oknach, a migawka nie swojego okna przestawiłaby podgląd na stronę, której operator tu nie otwierał.

## budowa/klient-poprzedni/src/moduly/browser/stan-przegladania.ts — migawka po ustaleniu okna
Rdzeń zna treść strony tego okna również wtedy, gdy przejście odbyło się w innej karcie albo przed przeładowaniem powłoki.

## budowa/klient-poprzedni/src/moduly/browser/stan-przegladania.ts — odpięcie pozycji pokrycia
Wpis kanału trzymałby inaczej przerysowanie kontrolek zdjętych już z drzewa, tak samo jak nieodpięty pasek uczciwości trzymałby przerysowanie po zamknięciu katalogu okien.

## budowa/klient-poprzedni/src/moduly/multitasking/sekcja-orkiestracja.ts
Ocenę układu wydaje rdzeń: zapis zależności oddaje pole poprawności wraz
z zastrzeżeniami, a sprawdzenie układu dokłada ścieżkę krytyczną. Widok nie
wylicza cyklu samodzielnie, tylko pokazuje odpowiedź rdzenia łącznie
z odmową.

## budowa/klient-poprzedni/src/moduly/automations/nawigacja-okien.ts
Nawigacja skraca drogę już istniejącą, nie otwiera nowej: wszystkie sześć okien modułu stoi w jednej przewijanej siatce, bez niej Operator dochodzi do Execution Monitora przewijaniem, z nią jednym wskazaniem. Mechanizm jest z biblioteki: rozwijanie, wędrówkę strzałkami, ślad wyboru na gałęzi, pole szukania po progu i zdanie o pustym wykazie niesie wspólny komponent menu drzewa, więc tutaj nie ma własnej listy rozwijanej ani nakładki wyboru. Ujawnianie jest stopniowe: okno wiodące stoi na pierwszym poziomie, bo jest wejściem do modułu, a dwie rodziny ról, kreatory i pas wykonania, odsłaniają się dopiero po wskazaniu gałęzi; siatka pod menu nie zmienia się, żadne okno nie znika i żadne nie dochodzi, więc przywołanie jest skokiem, nie otwarciem. Na uchwycie stoi bieżąca wartość nastawy, nie jej nazwa rodzajowa: przed pierwszym przywołaniem żadna wartość nie jest prawdziwa, Operator niczego nie przywołał, a wszystkie okna są na ekranie, więc uchwyt mówi to wprost i żaden liść nie niesie wtedy znacznika wyboru. Kod spoza wykazu okien modułu jest pomijany przy przestawianiu wskazania: menu nie zaczyna twierdzić, że przywołało coś, czego w module nie ma.

## budowa/klient-poprzedni/src/moduly/browser/stan-przegladania.ts — odmowa wykazu bez błędu okna
Okno jest ustalone i działa, nieudany jest jeden odczyt wykazu — dlatego faza okna zostaje nietknięta. Panel sam rozstrzyga, co z tym zrobić: wykaz z pozycjami zostaje na widoku, a wykaz pusty po odmowie nie przedstawia się jako pusty.

## budowa/klient-poprzedni/src/moduly/multitasking/sekcja-zespoly.ts
Szablony stoją tu jako gotowe nazwy do zapisania: naciśnięcie szablonu
zakłada zespół zapisem tego samego kształtu, a nadanie ról oknom robi
osobna sekcja Role. Skład pokazuje się nazwami ekspertów, bo zespół niesie
same identyfikatory, a nazwy dokłada osobny odczyt wykazu ekspertów; ekspert
nieznany wykazowi zostaje w składzie z identyfikatorem i zdaniem, że rdzeń
go nie oddał.

## budowa/klient-poprzedni/src/moduly/multitasking/stan-multitaskingu.ts
Własna obsada w każdym oknie sprawiłaby, że przepięcie wykonawcy pod innego
koordynatora zmieniłoby jedno okno i zostawiło trzy z obrazem nieaktualnym.
Okna żyją ze zdarzeń: fragment strumienia niesie strumień wykonawcy do
koordynatora, zmiana wiadomości domknięty wynik do analityka, a zmiana
kolejki stan etapu; odczyty okien, stanu okna i wiadomości służą pierwszemu
wypełnieniu i wznowieniu po rozłączeniu. Wewnątrz stoją trzy rejestry —
obsada ról, tury wykonawców i kolejki etapów — każdy z własnymi polami, a na
zewnątrz okna widzą je wyłącznie przez wspólny stan, bo koordynator musi
widzieć ten sam przebieg co wykonawca, którym steruje. Cudza rozmowa nie
pokazuje się jako praca własnego wykonawcy, bo strumień odsiewa fragmenty
spoza obsady.

## budowa/klient-poprzedni/src/moduly/browser/stan-przegladania.ts — pole pokrycia komend
Pokrycie niesie pozycje, których okno jeszcze nie wykonuje, wspólne dla wszystkich takich pozycji modułu.

## budowa/klient-poprzedni/src/moduly/browser/stan-przegladania.ts — zdanie o zaciągnięciu wykazów
Zdanie jest osobne od powodu ustalenia okna, bo dotyczy czego innego niż samo ustalenie okna.

## budowa/klient-poprzedni/src/moduly/browser/stan-przegladania.ts — stan zaciągnięcia wykazów
Samo zdanie o odczycie tego nie niosło: powód bywa niepusty również wtedy, gdy rdzeń po prostu nie zna jeszcze okna.

## budowa/klient-poprzedni/src/moduly/automations/graf-krokow.ts
Opracowanie modułu opisuje graf zależności jako widok własny, oddzielny od wykazu: wykaz mówi, co z czym jest związane, rysunek mówi, jak długi jest łańcuch i gdzie tory się rozchodzą — a tego z wykazu wierszy nie widać. Rysunek nie jest edytorem: zapis układu idzie komendami kontraktu z panelu akcji i z wykazu, kanwa pokazuje stan po zapisie. Układ węzłów liczy się z samych zależności, więc jest powtarzalny, lecz nie jest układem, który Operator mógłby ułożyć myszą — przenoszenie węzłów wymaga zapisania ich położenia osobną komendą, a tej rdzeń jeszcze nie obsługuje, do tego czasu kanwa układa graf sama i nie obiecuje, że zapamięta cudze ułożenie. Barwy i grubości nie stoją tutaj: węzły i krawędzie noszą klasy rodziny modułu, a wygląd niesie osobny arkusz stylów na żetonach koloru. Warstwa mówi, ile zależności trzeba przejść, zanim krok może ruszyć, więc układ warstwowy pokazuje kolejność wykonania wprost; węzeł stojący w cyklu nie ma najdłuższej drogi, rachunek zatrzymuje się wtedy na węźle już odwiedzonym i węzeł zostaje na warstwie, do której doszedł — cykl jest zastrzeżeniem układu, nie powodem, żeby nie narysować niczego.

## budowa/klient-poprzedni/src/moduly/browser/warstwa-adnotacji.ts — treść wiadomości adnotacji
Rdzeń materializuje URI danych do pliku i wplata jego ścieżkę w zapytanie, więc model zobaczy rysunek dopiero, gdy sięgnie po plik narzędziem odczytu. Adres strony i liczba śladów docierają do niego bez otwierania pliku i dlatego stoją w treści wiadomości, a nie tylko w podpisie załącznika.

## budowa/klient-poprzedni/src/moduly/multitasking/stany-relacji.ts
Znaczniki nazywają wyłącznie to, co rdzeń zlicza i oddaje: kolejka
wstrzymana odpowiada stanowi wstrzymania kolejki etapu bieżącego,
a koordynator wybudzony odpowiada dodatniemu licznikowi obiegów przy biegu
niezatrzymanym i przy żadnej turze wykonawcy w biegu, przy czym licznik
obiegów prowadzi pętla sesji w rdzeniu, nie klient. Wybudzenia nie zgadujemy
z czasu: znacznik czasu ostatniej zmiany biegu mówi, kiedy licznik ruszył
ostatnio, ale próg liczby sekund bez obiegu, który miałby znaczyć, że
koordynator myśli, byłby liczbą wymyśloną w kliencie. Rama okna ma jeden
znacznik, a przesłanki potrafią zajść naraz, więc kolejność sprawdzeń idzie
od stanu, który zatrzymuje pracę, do stanu, który ją opisuje, a pełny obraz
zostaje w polu powodu.

## budowa/klient-poprzedni/src/moduly/browser/warstwa-adnotacji.ts
Rysowanie mieszka w pliku `plotno-adnotacji.ts`, kontrolki w pliku `pasek-adnotacji.ts`, a ocena odpowiedzi rdzenia w pliku `skutek-zapisu.ts`. Rama okna stawia warstwę jako rodzeństwo podglądu w tym samym kontenerze pozycjonującym. Wyłączony tryb chowa warstwę atrybutem ukrycia — schowana nie łapie wskaźnika, więc zaznaczanie tekstu w podglądzie pozostaje możliwe. Rdzeń zostawia pole zrzutu puste, więc pod rysunkiem nie ma zrzutu strony. Zdanie o braku tła stoi w warstwie na stałe, aby operator wiedział przed wysłaniem, że w załączniku pójdzie sam rysunek. Wysyłka idzie poleceniem wysłania wiadomości, a nie przeniesieniem kontekstu: adresatem jest rozmowa tego okna, a przeniesienie międzymodułowe zaniosłoby rysunek gdzie indziej.

## budowa/klient-poprzedni/src/moduly/browser/warstwa-adnotacji.ts — barwa początkowa
Dwa zapisy tej samej domyślnej barwy rozjechałyby się przy zmianie palety, dlatego barwa początkowa nie stoi osobnym literałem.

## budowa/klient-poprzedni/src/moduly/browser/warstwa-adnotacji.ts — zdanie o braku tła i wytworze
Najpierw stoi to, czym model dysponuje, bo od tego zależy, czy operator dopisze zdanie wyjaśniające; dopiero potem brak tła i brak wpisu w wytworach. Zdanie o wpisie w wytworach przerysowuje się po każdej zmianie wykazu komend rdzenia, więc powód zmieni się sam w dniu, w którym rdzeń dostanie jej uchwyt.

## budowa/klient-poprzedni/src/moduly/assistant/zapis-modulu.ts
Plik odpowiada wyłącznie za to, co moduł wie i jak ta wiedza zmienia się po odpowiedzi rdzenia; powiadamianie widoków i subskrypcja zdarzeń mieszkają osobno, rozdzielenie pozwala sprawdzić przejścia bez budowania widoku. Wykaz zleceń obcych stoi osobno, bo odczyt zleceń czyta po oknie i podmienia nim cały wykaz — po takim odczycie zlecenie z innego okna zniknęłoby bez śladu, choć wciąż biegnie; zbiory są rozłączne z definicji, a rozdział pozwala jednemu z nich być podmienianym, a drugiemu trwać. Wpisy zleceń obcych biorą się wyłącznie ze zdarzenia rdzenia, które rozgłasza je do wszystkich połączeń konta; najczęstszą drogą jest polecenie głosowe, gdzie nakładka dobiera okno sama i nie sprawdza modułu, więc zlecenie zwykle siada na oknie rozmowy, nie na oknie asystenta. Wykaz nie obejmuje zleceń spoza okna założonych, zanim moduł został otwarty: kontrakt nie ma odczytu zleceń po sesji, więc pierwszym momentem, w którym klient dowiaduje się o takim zleceniu, jest jego pierwsza zmiana stanu. Znacznik pytania o zlecenia zostaje nietknięty przy wciągnięciu zlecenia obcego: zdarzenie o cudzym oknie nie jest odpowiedzią rdzenia w sprawie zleceń okna modułu, więc po nim wciąż nie wiadomo, czy okno modułu prowadzi coś swojego — postawienie znacznika zamieniłoby stan jeszcze-nie-pytałem w rdzeń-nie-prowadzi-nic bez odczytu, który by za tym zdaniem stał. Zlecenie okna modułu przysłane drogą zdarzenia, zanim ustalenie okna je poznało, nie zostaje w wykazie obcym: odczyt zleceń czyści z niego wszystko, co należy już do okna modułu.

## budowa/klient-poprzedni/src/moduly/browser/warstwa-adnotacji.ts — bufor płótna i spłaszczenie
Element schowany mierzy zero i dałby płótno o boku jednego piksela, dlatego bufor dostaje rozmiar dopiero, gdy warstwa jest widoczna. Płótno bez rasteryzacji oddaje pusty napis przy spłaszczeniu, więc sprawdzenie stoi przed wysyłką załącznika.

## budowa/klient-poprzedni/src/moduly/multitasking/tryby-wspolpracy.ts
Tryb mówi wyłącznie, który wykonawca dostaje zlecenie i co do niego jedzie;
automatyczne podawanie wyniku dalej po każdej turze byłoby drugim silnikiem
pętli obok pętli sesji w rdzeniu, więc moduł go nie buduje. Rdzeń nie
rejestruje komendy aktualizacji roli, więc profil roli po jego stronie się
nie zmienia. Selektor trybu stoi w oknie koordynatora i w obu oknach
wykonawców, a zapis jest jeden i ten sam, więc mieszka w tym pliku. Zapis
klucza trybu współpracy kończy się odmową walidacji, bo klucz leży poza
katalogiem ustawień, i odmowa cofa widok do wartości poprzedniej, a zdanie
mówi wprost, że zapisu nie ma. Rdzeń oddaje zapisany wpis, więc to jego
wartość mówi, co stoi w konfiguracji okna koordynatora.

## budowa/klient-poprzedni/src/moduly/assistant/okno-activity-feed.ts
Zapis jest chronologiczny i wspólny dla poleceń głosowych i tekstowych. Wpisy stoją w grupach dziennych: nagłówek doby rozdziela zapis tak, jak rozdziela go pamięć Operatora, a grupę składa okno, bo rdzeń oddaje wykaz płaski bez pola grupowania. Zawężenie, pole szukania i rodzaj wpisu, liczy się również w oknie, bo rdzeń nie ma czym zawęzić tego wykazu. Pustki są rozróżnione na pięć stanów: jeszcze nie pytałem, pytam, rdzeń nie zna ani jednego wpisu, pustka po zawężeniu — wpisy są, tylko żaden nie pasuje — i odmowa rdzenia wraz z jej powodem; komunikat zostaje w układzie, a ponowienie go nie usuwa. Powrót do wyniku nie potrzebuje osobnej komendy: wpis rodzaju wynik niesie treść wyniku, a zawężenie dziennika do jednego zlecenia idzie polem identyfikatora zlecenia w odczycie dziennika. Odsłuch pobiera bajty i odtwarza je w karcie, bo rdzeń oddaje wyłącznie nagrania, które sam wystawił; wyróżnienie idzie do rdzenia, więc przeżywa odświeżenie wykazu — po zapisie czytamy dziennik ponownie, żeby wykaz pokazywał stan zapisany, a nie przewidywany.

## budowa/klient-poprzedni/src/moduly/design/adnotacje-kompozycji.ts
Pole notatki warstwy niesie jedno zdanie bez autora i bez wątku, a przy każdym zapisie układu jedzie razem z całym układem i wraca przepisane od nowa. Uwaga zostawiona przez jedną osobę znikała więc przy pierwszym przesunięciu warstwy przez drugą — adnotacja dostała więc własny wiersz i własny czas. Pola autora tu nie ma i nie będzie: rdzeń bierze go z kontekstu wywołania, bo pole, w które da się wpisać cudze nazwisko, odbierałoby oznaczeniom osób w wątku całe ich znaczenie. Położenie kursora sprzed godziny nie jest wiedzą o niczym: panel zgłasza obecność na żądanie operatora, a nie w pętli, bo odpytywanie dziesięć razy na sekundę byłoby ruchem, którego nikt nie zamawiał. Kursory pozostałych przychodzą zdarzeniem obecności kompozycji.

## budowa/klient-poprzedni/src/moduly/design/adnotacje-kompozycji.ts — obecność zdarzeniem
Kursory pozostałych przychodzą zdarzeniem, nie odpytywaniem, i niosą komplet obecnych, więc panel je podmienia zamiast doliczać stan z ciągu przyrostów, którego początku nie widział.

## budowa/klient-poprzedni/src/moduly/multitasking/wykaz-nadan-rol.ts
Panel obsady składa scenę z wykazu okien: bierze okna i odczytuje z nich
pole roli, co wystarcza do pokazania sceny, ale nie jest rejestrem nadań;
rejestrem jest rdzeń, a jego odpowiedź niesie wcielenie, pole, którego okno
nie ma w ogóle. Przesiewanie odpowiedzi w kliencie dałoby ten sam obraz, ale
kłamałoby przy wykazie uciętym po stronie rdzenia i byłoby drugą regułą
przynależności do pętli obok tej z kontraktu. Zdejmij rolę wykonuje się bez
pytania, czy na pewno, bez wygaszania i bez uprawnienia per okno: okno
zostaje, traci wyłącznie rolę, a nadać ją z powrotem można paskiem obsady
stojącym wyżej w tym samym panelu. Zgoda nie jest dowodem: odpowiedź udana
z polem usunięcia fałszywym znaczy, że nie było czego zdjąć, a zdanie
o skutku powstaje z tego pola, nie z faktu, że wywołanie nie zwróciło
błędu. Zawężenie do wykonawców koordynatora wykonuje rdzeń, bo pole
identyfikatora koordynatora z pustym napisem byłoby zawężeniem do okna,
którego nie ma, i oddałoby wykaz pusty bez powodu. Odczyt wywołany po
zdjęciu roli zawraca też do wykazu nadań, więc drugiego wywołania stąd nie
ma: byłoby tym samym pytaniem zadanym dwa razy pod rząd.

## budowa/klient-poprzedni/src/moduly/design/adnotacje-kompozycji.ts — treść przy zamknięciu wątku
Kontrakt wymaga pola treści przy każdym zapisie, a podstawienie treści z pola okna przepisałoby cudzą uwagę przy okazji zamykania wątku.

## budowa/klient-poprzedni/src/moduly/design/adnotacje-kompozycji.ts — położenie z przesunięcia kanwy
Kursor myszy śledzony w pętli byłby ruchem na łączu, którego nikt nie zamawiał, a widok przesunięcia kanwy mówi, na co operator patrzy.

## budowa/klient-poprzedni/src/moduly/multitasking/zadania-w-tle.ts
Kontrakt nie zna bytu przepływu: zna podagenta z polami okna i zadania.
Powołanie podagentów przyjmuje jedno zadanie i liczbę podagentów do
piętnastu, więc podagenci jednego uruchomienia mają wspólne okno i wspólną
treść zadania, a przepływ jest dokładnie tą parą; nie zgaduje się go z nazw
podagentów, bo nazwa jest nieobowiązkowa, a rozbiór nazwy po dwukropku
byłby umową, której kontrakt nie zawiera. Podagent nie ma liczby żetonów,
liczby wywołań narzędzi ani nazwy etapu, więc model nie stawia dla nich pól
z zerem; nazwę etapu i skalę postępu bierze osobno telemetria monitora,
a żetony i narzędzia zostają brakiem zgłoszonym w widoku. Chwila odniesienia
dla czasu trwania przychodzi z zewnątrz funkcji składającej przepływy, żeby
wszystkie wiersze jednego przerysowania mierzyły się do tej samej sekundy.
Dopóki choć jeden podagent pracuje, odcinek czasu biegnie do chwili
odniesienia, bo przepływ trwa nadal i domknięcie go w zapisie byłoby
nieprawdą.

## budowa/klient-poprzedni/src/moduly/design/stan-designu.ts
Trzy okna korzystają z tego samego zbioru zasobów. Prompt Builder oddaje wynik generowania do Assets Panel, Assets Panel oddaje zasób na kanwę Design Board. Gdyby każde okno prowadziło własny wykaz, zasób wygenerowany w kreatorze nie pojawiłby się w panelu, a kanwa układałaby warstwy z zasobów, których panel już nie ma. Zdarzenie zmiany zasobu jest drugim źródłem odświeżenia: wciąga zasób powstały gdziekolwiek, także po stronie rdzenia, dokładnie tak samo jak własny odczyt. Odpytywania w pętli tu nie ma.

## budowa/klient-poprzedni/src/moduly/design/stan-designu.ts — zdjęcie zasobu po usunięciu
Zdjęcie idzie tą samą funkcją co gałąź zdarzenia usunięcia. Rdzeń rozgłasza usunięcie i zdarzenie i tak przyjdzie, ale okno, które właśnie kazało zasób usunąć, nie ma prawa pokazywać go dalej ani przez chwilę.

## budowa/klient-poprzedni/src/moduly/multitasking/zrodlo-biegu.ts
Koordynator widzi pełny strumień wykonawcy niesiony wspólnym zdarzeniem
fragmentu, w którym pole okna mówi, czyja to tura; okno koordynatora nie
zakłada drugiego kanału podglądu, tylko subskrybuje ten sam strumień
i odsiewa okna swoich wykonawców. Doręczenie polecenia niesie treść do okna
wykonawcy, więc wykonawca rusza do pracy, ale sama ta droga zostawia więź
wyłącznie w pamięci przeglądarki i gubi ją z jej zamknięciem, dlatego zapis
w bazie rdzenia utrwala, kto komu co zlecił, i zakłada pozycję kolejki.
Kontrakt nie ma osobnej komendy zatrzymania podagenta, a sterowanie kolejką
żąda jej identyfikatora, którego podagent nie niesie, więc ogniwem jest
nazwa kolejki, pod którą rdzeń zakłada kolejkę podagenta jego własnym
identyfikatorem. Żądanie z oknem przy subskrypcji monitora zakłada
obserwację.

## budowa/klient-poprzedni/src/moduly/multitasking/zrodlo-nadzoru.ts
Trzy sekcje mają jedno źródło, bo pracują na tym samym bycie: zależność
zapisuje się na układzie, harmonogram ustawia się na układzie, a monitor
pokazuje jego przebieg, więc osobne źródła pytałyby trzy razy o ten sam
wykaz automatyk. Komendy nadzoru nakładki leżą tutaj, bo przełącznik trybu
nakładki należy do sekcji Monitor procesu; stan nakładki niesie urządzenie,
sesję, okno i licznik procesów, ale nie niesie wartości obserwator albo
operator. Wiązanie kolejki z ekspertem, projektem albo automatyką jest
czynnością panelu, nie źródła biegu, bo bieg pętli koordynator–wykonawca
porusza rolami i nie zna projektów.

## budowa/klient-poprzedni/src/moduly/multitasking/zrodlo-okien.ts
Rola okna pochodzi z kontraktu, nie z oznaczenia własnego: pole roli
i identyfikator koordynatora niosą całą przynależność okna do pętli
koordynator–wykonawca, a moduł nie zakłada drugiego rejestru ról poza
komendą rdzenia. Nadanie roli idzie komendą przeznaczoną do nadawania roli,
a wcielenie osobną komendą aktualizacji roli, bo zmiana okna pola wcielenia
nie niesie; zmiana okna zostaje przy tym, co jest jej: tytuł, katalogi,
kanał modelu.

## budowa/klient-poprzedni/src/moduly/pokrycie-komend.ts
Zdanie wpisane w moduł na sztywno przestaje być prawdą w dniu, w którym
zmieni się rdzeń albo kontrakt, i nikt go nie zdejmuje, bo nic go z rdzeniem
nie łączy; taki napis myli też stronę braku, bo komenda bywa w kontrakcie,
a nie ma uchwytu w złożonym rdzeniu, więc pokrycie bierze rozstrzygnięcie
z odczytu i nazywa brak tam, gdzie jest. Byt stoi w korzeniu katalogu
modułów, a nie w jednym module, bo tę samą potrzebę ma każdy moduł, a
przepisywanie dałoby tyle samo rozjeżdżających się zdań o jednym stanie
produktu; sam fakt i zdania o nim leżą warstwę niżej, w źródle wykazu
komend rdzenia. Gdy żadna komenda nie jest nawet pomyślana dla czynności
wyłącznie okiennej, pozycja zostaje przy przycisku bez komendy z modeli
kontrolek formularza, bo bez nazwy komendy nie ma czego sprawdzać u
rdzenia, a byt nie zgaduje. Wywołanie przyjmuje etykietę pozycji w panelu
akcji, komendę, która tę pozycję by wykonała — nazwę z kontraktu albo
z wykazu okien operacyjnych, także taką, której kontrakt nie ma — oraz opis
czynności, który wchodzi w zdanie powodu. Powitanie idzie raz na połączenie,
nie raz na moduł, więc wszystkie wywołania odczytu czekają na jedną
odpowiedź.

## budowa/klient-poprzedni/src/moduly/design/czynnosci-zasobu.ts
Usunięcie wykonuje się bez pytania „czy na pewno": bramka potwierdzająca zabiera Operatorowi
jedno kliknięcie i uczy odruchowego potwierdzania, zamiast podnosić bezpieczeństwo. Okno wykonuje
czynność i mówi, co się stało.

Odpowiedź niesie treść, a nie samo powodzenie: `design.asset.remove` oddaje pole `removed`, gdzie
`false` znaczy „takiego zasobu nie było" i jest odpowiedzią udaną, a zarazem inną wiadomością niż
„usunięto". Zlanie ich w jedno zdanie potwierdzałoby czynność, która się nie odbyła.

Przyciski są czynne zawsze: brak wskazanego zasobu nie gasi kontrolki, a naciśnięcie mówi wtedy,
czego brakuje.

Stan ulubionego docelowy liczy się z zasobu, nie z napisu na przycisku. Napis bywa o ułamek sekundy
starszy od zbioru, ponieważ zdarzenie mogło właśnie przestawić ulubionego z drugiego okna, a wtedy
przełącznik wysłałby wartość, którą zasób już ma.

Rdzeń oddaje zasób odczytany z bazy po zapisie, więc porównanie zamówienia z odpowiedzią jest
tanie: gdy rdzeń zapisał co innego, okno mówi to wprost zamiast potwierdzać własne zamówienie.

Cisza kanału nie jest odmową usunięcia: rdzeń mógł zasób skasować, a odpowiedź zginąć z gniazdem.
Zdanie mówi o braku rozstrzygnięcia, a zasób zostaje w wykazie do czasu odczytu potwierdzającego
stan (`czuwanie-rdzenia.ts`).

Zasób był w wykazie, a rdzeń go nie zna — wykaz był nieaktualny. Zdejmujemy go z wykazu, ponieważ
pokazywanie dalej byłoby pokazywaniem bytu, o którym rdzeń właśnie powiedział, że go nie ma.

Napis mówi, co przycisk zrobi, a nie w jakim stanie zasób jest — przełącznik opisany stanem
bieżącym czyta się dokładnie odwrotnie do tego, co wykonuje.

## budowa/klient-poprzedni/src/dostepy/zrodlo-punktow.ts
Punkt dostępu określa, do czego model ma wgląd: do maszyny przez most MCP albo do katalogu lokalnego. Nie jest środowiskiem — środowisko pozostaje profilem widoczności modułów w bocznej nawigacji i tej sekcji nie dotyczy. Nie jest też katalogiem roboczym: ten mieszka w ustawieniu `katalog.roboczy.podstawa` i ma w tej sekcji własny obszar.

## budowa/klient-poprzedni/src/ikony/ikony.ts
Element `<svg>` bez atrybutu `width` bierze całą szerokość rodzica, a `preserveAspectRatio` wyśrodkowuje znak w pustym polu — logotyp wygląda wtedy na wielokrotnie mniejszy, niż wynika z podanej wysokości; stąd `szerokoscZnaku` wylicza szerokość jawnie.

Funkcja `elementGodla` daje wyłącznie sygnet i dobiera przy tym odmianę uproszczoną poniżej 16 pikseli. Odmiany złożone, dostępne w `marka.ts`, mają inne proporcje niż kwadrat, więc `elementZnaku` idzie osobną drogą: `zbudujElement` nadaje bok kwadratowy z parametru `rozmiar` przy sygnecie, a przy odmianach złożonych `rozmiar` rozstrzyga wyłącznie wysokość, a szerokość zostaje wyliczona osobno.

## budowa/klient-poprzedni/src/komponenty/rama-okna.ts
Rola okna stoi w nagłówku, bo rozstrzyga układ: okna dzielą się na wiodące, pomocnicze, monitory, kreatory i zarządców, moduł ustawia je w pasach według roli, a Operator ma widzieć, dlaczego okno stoi tam, gdzie stoi.

Fazy okna (puste, ładowanie, błąd, gotowe) nie należą do ramy. Rama daje `cialo`; przesłonę stanu buduje osobny byt modułu (`stany-okna`) i moduł osadza ją w ciele. Inaczej obudowa znałaby cykl życia danych, których nie pobiera.

Klasy modułu w `OpisRamyOkna.przedrostek` nie wnoszą wyglądu — wygląd jest w bibliotece — ale bywają uchwytem reguł własnych arkusza: `.dt-okno` niesie rozmiar pisma terminala i selektor `[data-zawijanie]`. Moduł bez takiej reguły nie podaje nic i dostaje wygląd biblioteczny.

## budowa/klient-poprzedni/src/konfiguracja/kontrolki-tekstowe.ts
Zatwierdzenie wartości następuje na zdarzeniu `change`, czyli po opuszczeniu pola, nie po każdym znaku. Zapis co znak zasypałby rdzeń komendami `config.set` i odbierał możliwość poprawienia wartości przed wysyłką.

Wartość rodzaju `secret` nie wraca z rdzenia: pole pozostaje puste i mówi to wprost. Kontrolka przyjmuje wartość nową, nie pokazuje wartości zapisanej.

## budowa/klient-poprzedni/src/konfiguracja/obszary-sesji.ts
Plik nie buduje ani jednego elementu widoku i nie woła rdzenia — tak samo jak `zasiegi.ts`, którego jest odpowiednikiem dla drugiego kształtu wartości. Nazwy poziomu i osi bierze z `adres-ustawienia.ts`: pochodzenie obszaru składa się na `AdresUstawienia` i idzie przez `opisAdresu`, więc poziom nazywa się w oknie tak samo, jak nazywa się w łańcuchu zapisów. Wartości wyliczeń pochodzą wyłącznie z kontraktu; literału nazwy obszaru w kodzie nie ma.

Rozejście katalogu roboczego z `checkedAt` równym zeru składane jest z samej konfiguracji obowiązującej, bez dotykania dysku, zgodnie z `core/sesja_konfiguracja_skladanie.go`.

## budowa/klient-poprzedni/src/konfiguracja/okno-konfiguracji.ts
Okno stoi na natywnym `<dialog>`, więc warstwę tła, stos okien i zamknięcie klawiszem Esc daje przeglądarka, a nie własna nakładka. Wygląd bierze z biblioteki `komponenty/` (`dn-modal`) — plik nie zna ani jednej barwy. Kategorie i pola przychodzą z katalogu rdzenia; plik zna wyłącznie trzy obszary układu i sposób ich związania. Przycisk odświeżenia pozwala spytać rdzeń ponownie, gdy katalog nie dotarł.

Obszary sesji rozstrzyga rdzeń (`config.effective.get`, `config.session.set`), bo tylko on zszywa je z rejestrami spoza rodziny `config.*` — konta, kanały, tożsamości, dostępy. Łańcuch pojedynczych kluczy zostaje po stronie klienta (`rozstrzygniecie.ts`), bo podgląd dziedziczenia potrzebuje wszystkich zapisów, nie samego zwycięzcy.

Kanały modelu (`channel.add`, `channel.update`, `channel.remove`) mają panel w oknie konfiguracji, a nie w komplecie sterowania, bo rejestr kanałów jest bytem globalnym — katalogiem wyboru, nie ustawieniem okna (`sterowanie/panel-sterowania.ts`) — a komplet sterowania jest per okno: wstawiony tam panel powstawałby raz na każde otwarte okno. Egzemplarz rejestru zakładany w oknie konfiguracji jest czytającą pamięcią podręczną nad `channel.list`, bez ani jednej drogi zapisu, więc drugi egzemplarz to drugi odczyt tej samej prawdy, nie druga prawda. Po każdym udanym zapisie panel woła `rejestr.odswiez()`, co ogłasza zmianę wszystkim czytelnikom naraz — oknu rozmowy, obu sterowaniom modelu, panelowi modeli i Roundtable.

Pozycja „Izolacja" w nawigacji zakresów prowadzi do okna punktów izolacji, jednego na klienta, które pamięta swój stan — otwarcie z panelu kategorii trafia w ten sam egzemplarz, co otwarcie z listwy Ustawień, nie zakłada drugiego.

Drugi selektor zasięgu byłby powieleniem, a rozjazd między nimi pokazywałby wartości z dwóch różnych zasięgów obok siebie.

## budowa/klient-poprzedni/src/konfiguracja/panel-kategorii.ts
Zaszycie w tym panelu listy pól oznaczałoby, że dodanie ustawienia wymaga zmiany kodu w dwóch miejscach — czego to okno ma właśnie nie robić.

Zakres „Izolacja" jest jedyną kategorią z przejściem dalej: rozwinięty jest w osobnym, trzypanelowym oknie punktów izolacji, bo ma własną złożoność — dwa rodzaje izolacji, siedem poziomów zasięgu i profile. Jedenaście wierszy katalogu pokazywanych w tym panelu to te same klucze, ale bez selektora zasięgu, macierzy i podglądu polityki efektywnej — przejście stoi nad formularzem, żeby ta pozycja nawigacji prowadziła tam, gdzie zakres jest konfigurowany w całości. To, że pozycja „Izolacja" otwiera okno punktów izolacji, jest wymaganiem struktury okna, nie metadaną ustawienia; katalog niesie wyłącznie wiersze.

Różnica między jedenastoma wierszami katalogu w tym panelu a oknem punktów izolacji z selektorem zasięgu i podglądem polityki efektywnej jest powodem, dla którego przejście do tego okna tu stoi. Przycisk przejścia niczego nie wygasza: obie drogi pozostają czynne.

O błędzie odczytu katalogu mówi komunikat blokowy nad stopką okna, nie stan pusty formularza.

## budowa/klient-poprzedni/src/konfiguracja/panel-zaczepow.ts
Panel zaczepów stoi osobno od panelu obszarów, bo tamten utrwala wartości obowiązujące i nie redaguje pól, a zaczep trzeba złożyć z punktu cyklu życia, polecenia i zawężenia. Zaczep bez zdarzenia albo bez polecenia odrzuca rdzeń (`adapter_rozmowa_powierzchnia.go`). Punkt cyklu życia jest polem wpisu z podpowiedzią, nie listą zamkniętą — to wartość danych, a nie typ kodu.

## budowa/klient-poprzedni/src/konfiguracja/pole-ustawienia.ts
Wiersz nie zna ani jednego klucza z osobna — wszystko, co rysuje, pochodzi z pozycji katalogu przekazanej w zależnościach. Rodzaj kontrolki rozstrzyga jeden moduł rozdzielający, więc nowy rodzaj wartości nie dotyka tego pliku.

Adres zapisu wskazuje domyślnie punkt widzenia okna, po zmianie poziomu — poziom wybrany przez Operatora. Wynik zapisu widnieje przy polu, żeby naciśnięcie zawsze dało odpowiedź: powodzenie mówi, gdzie zapisano, niepowodzenie mówi, co odpowiedział rdzeń.

Przemilczenie adnotacji o stanie klucza w objaśnieniu kazałoby Operatorowi wierzyć, że zapisana wartość steruje wykonaniem, mimo że żadna ścieżka rdzenia jej nie czyta.

## budowa/klient-poprzedni/src/moduly/assistant/stan-assistant.ts
Trzy okna obserwują ten sam zapis zleceń i dziennika: Voice Console zakłada zlecenie, Actions
Monitor nim steruje, Activity Feed pokazuje jego przebieg. Gdyby każde okno prowadziło własny
wykaz, wstrzymanie zlecenia w monitorze nie zmieniłoby tego, co Voice Console uważa za polecenie
w toku.

Poza własnym działaniem jedynym źródłem odświeżenia jest zdarzenie `assistant.action.changed`: ono
wciąga zmianę dokonaną gdzie indziej tak samo jak zmianę własną, więc nie ma tu odpytywania w pętli.

Zejście zlecenia z toru pociąga odczyt dziennika. Rdzeń dopisuje wpis rodzaju `result` dopiero przy
domykaniu zlecenia, czyli po odczycie, który Voice Console robi zaraz po wysłaniu polecenia.
Zdarzenie o stanie końcowym jest jedyną chwilą, w której klient wie, że dziennik urósł, więc odczyt
jedzie właśnie tu.

Zdarzenie przychodzi z całego konta, nie z jednego okna: rdzeń rozgłasza je do wszystkich połączeń
konta, a kanał klienta nie zawęża zdarzeń do sesji. Zlecenie cudzej sesji wciągnięte do tego zapisu
byłoby cudzą historią pokazaną jako własna, więc granicą jest sesja z koperty, nie okno z treści.
Rdzeń wypełnia pole sesji koperty sesją okna zlecenia, więc porównanie z sesją kanału opiera się na
danych, nie na domyśle. Zlecenia tej samej sesji spoza okna modułu, na przykład z nakładki Always
On Display, która dobiera okno sama, idą do wykazu osobnego, bo są pracą asystenta widoczną na
ekranie, a nie cudzą historią.

Trzy źródła dobudowane obok rdzenia modułu, każde nad własną rodziną komend: mowa, konteksty
pamięci wraz z zajętością okna oraz historia schowka wraz ze słownikiem skrótów i skrótem globalnym.
Osobno, bo osobno znikają: maszyna bez silnika mowy ma sprawną historię schowka, a maszyna bez
powłoki okiennej — sprawne konteksty pamięci.

Dziennik czytamy dopiero po zejściu zlecenia z toru, bo wcześniej rdzeń nie dopisał do niego ani
jednego wpisu. Zawężenie `actionId` czyta dziennik zlecenia bez oglądania się na okno, więc przebieg
zlecenia z nakładki Always On Display jest widoczny w Activity Feed po jego wskazaniu w monitorze.

## budowa/klient-poprzedni/src/konfiguracja/stan-konfiguracji.ts
Katalog (kategorie i definicje) oraz wpisy konfiguracji trzymane są razem, ponieważ pole formularza potrzebuje obu naraz: definicja mówi, jaką ma być kontrolką, wpisy mówią, skąd bierze się jej wartość. Dwa równoległe stany dałyby dwie prawdy o tej samej wartości. Rdzeń, który nie odda katalogu, zostawia wykazy puste; okno pokazuje wtedy komunikat i pozostaje otwarte.

Bez fazy odczytu pusty katalog znaczy trzy rzeczy naraz: „jeszcze nie pytałem", „pytam" i „rdzeń nie zna ani jednej kategorii". Każdej należy się inny stan okna: nic, wskaźnik odczytu, stan pusty.

Wykaz bez zapisu obsługuje jeden przepis na dwie drogi: przywrócenie własne (przycisk „Przywróć" tego okna) i przywrócenie cudze (zdarzenie `config.changed` z innego okna albo urządzenia) zdejmują zapis dokładnie tak samo. Dwa przepisy rozeszłyby się na wskaźniku pochodzenia: jeden egzemplarz stanu pokazywałby wartość przywróconą, drugi wartość zapisaną i widmowy wiersz łańcucha dziedziczenia.

Rodzaj zmiany pominięty w odpowiedzi rdzenia znaczy zapis — tak wchodzi odpowiedź na własną komendę `config.set`.

Zmiana punktu widzenia dociąga poziom, który do tej pory nie był czytany; bez tego wartość spod okna, sesji czy projektu nie ma pokrycia w stanie, a pochodzenie wskazuje poziom globalny. Funkcja odczytu wpisów ogłasza od razu, względem wpisów już znanych, i ponownie, gdy poziom dojedzie.

## budowa/klient-poprzedni/src/konfiguracja/stan-obszarow-sesji.ts
Punkt widzenia jest wspólny z resztą okna: pasek u góry okna ustala, dla kogo liczymy wartości obowiązujące, a ten stan bierze od niego ten sam adres, którym jedzie łańcuch zapisów kluczy. Dwa punkty widzenia w jednym oknie znaczyłyby, że plakietka klucza i plakietka obszaru mówią o dwóch różnych miejscach przestrzeni konfiguracji.

Zapis utrwala to, co obowiązuje: panel nie redaguje pól obszaru, tylko zapisuje obszar odziedziczony z poziomu szerszego albo z innego rejestru jako zapis własny na poziomie wskazanym punktem widzenia. Treść zapisu bierze się z konfiguracji obowiązującej, którą oddał rdzeń. Obszar spoza wykazu `areas` rdzeń zostawia nietknięty, a obszar w wykazie bez treści usuwa, wracając do dziedziczenia.

Obszar `hooks` jest jedynym obszarem redagowanym w oknie konfiguracji: zaczep trzeba móc złożyć, a nie tylko utrwalić odziedziczony.

## budowa/klient-poprzedni/src/moduly/design/czuwanie-rdzenia.ts
Czuwanie nad czynnością, której rdzeń nie rozstrzygnął, mówi prawdę o wywołaniu bez odpowiedzi
i odróżnia rdzeń pracujący od kanału milczącego. Gdy gniazdo padnie w trakcie oczekiwania, okno
samo się nie odnajdzie: obietnica wywołania nigdy nie jest odrzucana, a odbiorca odpowiedzi żyje
w rejestrze korelacji przypisanym do gniazda, które padło — odpowiedź nie przyjdzie już nigdy,
także po ponownym połączeniu.

Zwykły limit czasu nie odróżnia rdzenia, który pracuje długo, od kanału, który zamilkł, a komendy
modułu bywają wolne i długie oczekiwanie jest tu stanem poprawnym. Czuwanie pyta więc kanał
o życie: wysyła jedno tanie żądanie kontraktu tą samą drogą. Gdy próba odpowiedziała, kanał żyje
i czynność nadal trwa; gdy próba zamilkła, kanał zamilkł i skutek pozostaje nieznany.

Ta sama próba jest zarazem czujnikiem powrotu i nie kosztuje dodatkowej ramki: żądanie wysłane
przy rozłączeniu czeka w kolejce wychodzącej, a rdzeń odpowiada na nie w chwili ponownego
połączenia. Jedna obietnica próby mówi więc najpierw, że kanał milczy, a potem, że łączność
wróciła; pętli odpytującej tu nie ma.

Czuwanie nie ogłasza niepowodzenia czynności i nie zgaduje jej skutku — rdzeń mógł żądanie odebrać
i wykonać, zanim gniazdo padło. Jedyną drogą do prawdy jest odczyt po powrocie łączności i tak
brzmi zdanie dla Operatora.

## budowa/klient-poprzedni/src/konfiguracja/zrodlo-wartosci.ts
Pola osi są w kontrakcie nieobowiązkowe — puste znaczy `platform` — więc adres platformy wychodzi bez nich. `config.get` bez poziomu oddaje samą wartość obowiązującą, bez pozostałych zapisów, a okno pokazuje również, z którego poziomu wartość pochodzi i jakie zapisy stoją obok. Łańcuch punktu widzenia liczy klient (`rozstrzygniecie.ts`) z wpisów globalnych oraz wpisów zapisanych dokładnie na wskazanym poziomie i bycie, więc odczyt pobiera surowe wpisy tych poziomów osobnymi zapytaniami z podanym `scope`.

Surowy odczyt zawęża się do granulacji poziomu (zasięg i byt poziomu), bez osi: łańcuch bierze wszystkie zapisy poziomu niezależnie od osi, a którą oś przyjąć rozstrzyga klient względem punktu widzenia.

Wynik `config.get` zasila rozstrzyganie pochodzenia po stronie klienta. Pominięty punkt znaczy widok globalny.

Bez rodzaju zmiany w zdarzeniu `config.changed` nie da się odróżnić zapisu od usunięcia, bo rdzeń rozgłasza przywrócenie wartości domyślnej wpisem niosącym starą wartość (`core/handlers_config.go`, `core/adapter_ustawienia.go` funkcja `Przywroc`); odczyt samego wpisu wstawiłby skasowany zapis z powrotem do wykazu.

## budowa/klient-poprzedni/src/moduly/design/eksport-tokenow.ts
Wszystko tutaj składa przeglądarka z wartości odczytanych z motywu obowiązującego — rdzeń nie
bierze w tym udziału i nie musi. To jest ta część grupy wydań tokenów projektowych, która komendy
nie potrzebuje, więc nazwanie jej brakiem kontraktu byłoby zmyśleniem długu.

Opracowanie wymienia przy eksporcie także wydanie dla systemów mobilnych oraz wydanie przewodnika
do modułów Library i Studio. Pierwszego nie ma, bo wymagałoby przekładu ról systemu wizualnego na
pojęcia dwóch obcych platform — a taki przekład jest rozstrzygnięciem projektowym, nie zapisem
pliku. Drugiego nie ma, bo wydanie czegokolwiek do innego modułu wymaga komendy, której kontrakt
nie zna; okno nazywa ten brak zamiast wysyłać plik w próżnię.

Postacie wydania nie są tu ozdobą: zmienne CSS to postać, w której ten produkt żetony trzyma;
pozostałe trzy są postaciami, w których przyjmuje je kod korzystający z systemu.

## budowa/klient-poprzedni/src/moduly/design/etykiety-designu.ts
Wykaz braków jest danymi, a nie zdaniami rozsianymi po widokach, żeby każde okno nazywało ten sam
brak tak samo. Wpis znika stąd z chwilą, w której kontrakt dostaje komendę dla danej czynności.

Kod modułu z kolumny katalogu rdzenia jest jedynym miejscem w kliencie wiążącym ten katalog z kodem
modułu Design; nie jest kopią katalogu modułów — moduł mówi wyłącznie, którym modułem sam jest,
wykazu pozostałych czternastu tu nie ma.

Okno Tokens & System Panel jest piątym oknem operacyjnym modułu: rozjazd nie jest przemilczany,
moduł zgłasza katalogowi okien komplet pięciu kodów, a byt wspólny wypowiada obie strony różnicy —
okna rejestru, których moduł nie buduje, oraz okna budowane spoza rejestru. Dopisanie wiersza do
rejestru rdzenia zdejmie tę drugą połowę bez zmiany ani jednej linii tutaj.

Wykaz czynności bez drogi powstał, gdy kontrakt nie niósł dla nich ani jednej komendy. Po scaleniu
rodziny komend design większość z nich komendę już ma — brakiem jest już uchwyt w rdzeniu i droga
z okna, a to jest inne zdanie. Rozstrzygnięcie o pokryciu nie należy jednak do tego pliku: napis nie
jest z rdzeniem połączony i zestarzeje się znowu.

## budowa/klient-poprzedni/src/kontrakt.test.ts
Kontrakt jest jedynym źródłem prawdy nazw, a `contract.ts` jego wytworem. Trzy rzeczy mogą się tu rozejść po cichu i żadnej nie wychwyci ani kompilator, ani przegląd: generat starszy od źródła (`contract.json` zmienione, generator niepuszczony — kompilacja przechodzi, bo stała nadal istnieje), literał nazwy powielony w kodzie klienta zamiast wzięty z generatu (zmiana nazwy w kontrakcie zostawia wtedy w interfejsie martwe wywołanie, które rdzeń odbije jako `*.unknown`), oraz port rdzenia zaszyty w kliencie rozjechany z portem domyślnym rdzenia — klient szuka gniazda tam, gdzie nikt nie nasłuchuje.

Nierozpoznana komenda wraca pod nazwą swojego obszaru, a nie pod nazwą połączenia — odmowa poczty przedstawia się jako sprawa poczty. Sprawdzian obszarów wypada niepomyślnie zarówno wtedy, gdy taki obszar się pojawi, jak i wtedy, gdy wiersz zostanie tu po obszarze już domkniętym.

Zmiana nazwy w `contract.json` przechodzi przez kompilację obu stron i zostawia w interfejsie martwe wywołanie, które rdzeń odbije zdarzeniem `*.unknown` — dług nie rośnie i nie znika po cichu.

## budowa/klient-poprzedni/src/modele/stan-tozsamosci.ts
Trzy elementy stanu — katalog kategorii, zapisy treści i nakładka — występują razem, ponieważ edytor korzysta z wszystkich naraz: katalog wskazuje dostępne kategorie i proponowany tryb, zapisy niosą treść, a nakładka określa, co z tej treści trafia do modelu. Oś wyznacza zakres odczytu: zapisy pobierane są dla osi wskazanej, nie dla wszystkich naraz, więc kategoria bez zapisu na danej osi ma treść pustą, a obowiązuje dla niej treść z osi szerszej. Oś wymagająca bytu, dla której bytu nie ma, jest traktowana jak platforma — adresowanie takiej osi nie ma znaczenia w kontrakcie.

## budowa/klient-poprzedni/src/moduly/design/wyszukiwarka-funkcji-designu.ts
Wyszukiwarka niczego nie uruchamia i nie udaje, że uruchamia: pozycja bez drogi
nie dostaje przycisku, który po naciśnięciu przeprosi, tylko zdanie o tym, czego
brakuje. Szukanie idzie środkiem nazwy, opisu, grupy i okna, bo nazwy pozycji są
w części angielskie i złożone — szukanie wyłącznie od początku nazwy nie
znalazłoby pozycji po słowie wpisanym z pamięci. Liczba w zdaniu o zasięgu
katalogu liczy się z wykazu przy każdym odświeżeniu, bo zapisanie jej wprost
rozjeżdżałoby się z listą pozycji przy zmianie wykazu.

## budowa/klient-poprzedni/src/moduly/design/inspektor-warstwy.ts
Inspektor mówi też, czego warstwa NIE niesie: wypełnienie, obrys, efekty, więzy responsywne
i auto-layout — warstwa kompozycji nie ma pola na żadne z nich, więc wykaz stawia je jako brak
nazwany, zamiast pokazywać puste pola sugerujące, że wartość istnieje, tylko jest niewypełniona.

Zmiana liczby idzie do zapisu kompozycji, nie do rdzenia. Do rdzenia jedzie dopiero cały układ,
zapisem kompozycji — tak samo jak przy przeciąganiu warstwy po kanwie.

Przy zaznaczeniu wielokrotnym inspektor opisuje warstwę pierwszą z wykazu i mówi o tym wprost:
cztery liczby opisują jeden prostokąt, a nie zbiór.

Pole zostawione z wpisem, którego nie da się odczytać jako liczby, zostawia wymiar bez zmiany,
zamiast zsuwać warstwę do lewego górnego rogu.

## budowa/klient-poprzedni/src/moduly/design/zapis-designu.ts
Powiadamianiem widoku zajmuje się osobny zapis stanu; odczyt z rdzenia i
rozgłoszenie zmiany to dwie różne czynności rozdzielone między te dwa pliki.
Nadawcą zdarzenia usunięcia zasobu jest procedura po stronie rdzenia, która
wysyła je wyłącznie wtedy, gdy wiersz naprawdę zniknął — bez tej gałęzi zdarzenie
trafiałoby do gałęzi wciągającej i okno pokazywałoby jako obecny zasób, o którym
rdzeń właśnie powiedział, że go nie ma. Przyjęcie odpowiedzi stoi osobno od
wywołania, bo odpowiedź bywa, że nie przyjdzie: odczyt zlecony przed zerwaniem
połączenia nie dostaje odpowiedzi nigdy, więc czuwanie nad połączeniem musi
rozstrzygnąć stan okna samo. Do odmowy okno dokleja zdanie o torze komendy, bo
powód dotyczy wtedy żądania; zerwanego połączenia żadne pole żądania nie
tłumaczy, więc przy nim dopisek byłby wyłącznie hałasem. Granicę wyjściową
zmienia się w filtrze panelu zasobów.

## budowa/klient-poprzedni/src/modele/formularz-konta.ts
Pole poświadczenia jest wyłącznie wejściem: kontrakt nie zwraca zapisanej wartości żadną komendą, a puste pole przy zmianie zostawia poświadczenie dotychczasowe, zamiast je kasować. Komenda zmiany konta nie przyjmuje rodzaju, więc przy zmianie konta rodzaj jest plakietką informacyjną, a nie kontrolką. Katalog konfiguracji dotyczy wyłącznie kont programu code CLI; dla pozostałych rodzajów pole znika w całości.

## budowa/klient-poprzedni/src/moduly/design/kolekcje-designu.ts
Kolekcja jest bytem osobnym od etykiety, choć obie grupują zasoby. Etykieta jest słowem: nie ma
nazwy własnej ani porządku, a przemianowanie jej wymaga przepisania każdego zasobu z osobna.
Kolekcja ma nazwę, opis i kolejność, i przeżywa odświeżenie okna, bo leży w bazie stanowiska.

Panel nie wysyła nigdy stanu kolekcji, tylko zmianę jednego zasobu: kolekcja bywa duża,
a przepisywanie jej przy każdej zmianie jest drogą do zgubienia zawartości, gdy dwa okna wyślą
swój stan naraz.

Liczba zasobów w podpisie pochodzi z odpowiedzi rdzenia, a nie z długości wykazu, który okno akurat
trzyma. Prawdą o kolekcji jest to, co w niej leży w bazie.

Zawężenie do zasobu ukryłoby te kolekcje, do których Operator akurat chce zasób dołożyć.

Zasób, który w kolekcji już był albo go w niej nie było, nie zmienił niczego, a potwierdzenie
zmiany byłoby potwierdzeniem czynności, która się nie odbyła.

## budowa/klient-poprzedni/src/modele/kontrolki-formularza-braki.ts
Okno operacyjne modułu nie stawia formularza z etykietą nad polem: wkłada samą kontrolkę do panelu akcji albo do paska narzędzi, a nazwę niesie atrybut opisujący. Kontrolki hubowe zwracają wiersz złożony z elementu i kontrolki, którego w pasek akcji włożyć się nie da — tu stoją ich odpowiedniki bez wiersza. Wejście zostaje jedno: hub reeksportuje ten plik w całości. Wygląd pochodzi z biblioteki komponentów współdzielonych, więc plik nie zna ani jednej barwy i ani jednego odstępu. Klasa rodziny modułu przychodzi parametrem i należy do modułu wywołującego.

## budowa/klient-poprzedni/src/modele/kontrolki-formularza-braki.ts (przycisk bez komendy)
Pozycja przycisku nie znika ze sceny, bo okno bez niej wyglądałoby na kompletne, a brak przestałby być widoczny. Wygaszenie jest tu tak samo niedopuszczalne jak milczenie: element nie traci obsługi zdarzeń ani nie zmienia kursora. Powód idzie równolegle trzema drogami: tytułem, opisem dostępności i znacznikiem danych dla bram i sprawdzianów.

## budowa/klient-poprzedni/src/mission-control/model-danych.ts
Model danych pulpitu ma jedną odpowiedzialność: kształt danych, które widok pulpitu umie wyrysować; wartości buduje `zlozenie-danych.ts` wyłącznie z odczytów i zdarzeń rdzenia. Pole `null` znaczy brak źródła: miara, której kontrakt nie niesie, ma w modelu typ `X | null`, a widok wypisuje przy niej etykietę „brak źródła danych w kontrakcie" zamiast liczby. Nazwy stanów pochodzą z kontraktu (`shared/contract`) zamiast z powtarzanych literałów.

Typ kodu środowiska jest napisem, a nie unią wywiedzioną z `KnownModuleIds`, bo wykaz środowisk należy do rdzenia jako dane i `environment.list` niesie go w całości — unia zamykałaby matrycę na kody znane klientowi, a środowisko spoza niej trafiałoby do wykazu „poza środowiskami" jako sesja bez wskazania środowiska. `KnownModuleIds` pilnuje kodów tam, gdzie klient sam je wymienia (`strona-glowna/pozycje-srodowisk.ts`); kolumna matrycy przepisuje to, co przyszło z rdzenia.

Widok wypisuje etykietę braku źródła przy miarach wysycenia, kolejki, limitu i kosztu kanału.

## budowa/klient-poprzedni/src/moduly/design/kontrast-wcag.ts
Pary i progi kontrastu nie są wymyślone w tym module — pochodzą z wykazu progów kontrastu produktu,
tego samego, którym mierzy przyrząd pomiaru produktu. Ułożenie tu własnej listy par dałoby drugą
miarę jednego stanu, a dwie miary zawsze się rozjeżdżają.

Wartości barw czytane są z motywu obowiązującego, więc tabela mówi o produkcie w tej chwili,
a nie o zapisie sprzed przełączenia motywu. To jest zarazem jedyna droga, żeby ocenić motyw jasny
i ciemny osobno — a system wizualny traktuje je jako równoprawne.

Wykaz progów nie podaje roli pary. WCAG dopuszcza próg niższy dla tekstu dużego oraz dla obrysów
i wskaźników skupienia, więc para obrysu wychodzi tu poniżej progu, choć wobec właściwej reguły
może być zgodna. Wynik oznaczamy więc jako pomiar, a nie jako werdykt.

## budowa/klient-poprzedni/src/modele/kontrolki-formularza.ts
Formularz konta i edytor tożsamości opisują byty o polach stałych, wynikających wprost z kontraktu, a nie z katalogu wierszy, więc nie mogą korzystać z generatora pól okna konfiguracji, który buduje kontrolkę z definicji ustawienia — nie ma czego mu podać. Zamiast dwóch równoległych sposobów budowania pola w dwóch plikach sekcja ma jeden ten. Wygląd pochodzi w całości z biblioteki komponentów współdzielonych, więc plik nie zna ani jednej barwy i ani jednego odstępu.

## budowa/klient-poprzedni/src/modele/kontrolki-formularza.ts (pole poświadczenia)
Kontrakt przyjmuje poświadczenie w żądaniu i nie zwraca go żadną komendą. Pole jest zatem wyłącznie wejściem: puste znaczy „nie zmieniaj”, wypełnione znaczy „zapisz nowe”. Nigdy nie pokazuje wartości zapisanej, ponieważ klient jej nie ma.

## budowa/klient-poprzedni/src/moduly/design/modul-design.ts
Jedna odpowiedzialność: złożenie okien modułu i rozdanie im jednego stanu. Układ idzie warstwami
widoczności opracowania, nie kolejnością plików. W pasie pierwszym i drugim stoją okna warstwy
pierwszej — te, które są widoczne bez interakcji. W pasie trzecim stoją rozwinięcia warstw
wyższych: Tokens & System Panel (warstwa trzecia, wywoływany menu kebab), wyszukiwarka funkcji
i wykaz skrótów (warstwa czwarta). Zwinięte nie znaczy ukryte — zapowiedź nad każdym mówi, co jest
pod spodem.

Układ wynika z roli okna. Design Board jest wiodące i stoi w pasie pierwszym na całą szerokość —
na nim odbywa się praca koncepcyjna. W pasie drugim stoją trzy pozostałe w kolejności katalogu
rdzenia: Assets Panel (zarządca) wskazuje zasób, Preview Window (pomocnicze) pokazuje zasób
wskazany, Prompt Builder (kreator) zleca nowy. Podgląd stoi między nimi, bo patrzy i na to, co
zarządca wskazał, i na to, co kreator dopiero przyniósł.

Jeden zbiór zasobów na cały moduł: wynik generowania z kreatora wchodzi do wykazu zarządcy, stamtąd
na kanwę wiodącego, a podgląd czyta ten sam wybór, bo stan modułu jest jeden. Zdarzenie zmiany
zasobu wciąga zasób tą samą drogą także wtedy, gdy zlecenie przyszło z obcego połączenia.

Druga droga na kanwę prowadzi z rozmowy. Zasób zlecony spoza okien modułu wchodzi zdarzeniem
zmiany zasobu: stan wciąga go do wykazu zarządcy, a to złożenie kładzie go warstwą na kanwie
wiodącego. Wiązanie mieszka tutaj, bo wiąże dwa okna i nie jest sprawą żadnego z nich z osobna.

Moduł nie osadza się sam — oddaje element; gdzie stanie, rozstrzyga warstwa składająca.

Okno modułu musi być znane przed odczytem zasobów: odczyt zasobów przyjmuje identyfikator okna,
a bez niego dotyczyłby czegoś innego niż to okno. Zaplecze idzie równolegle — jest niezależne.

## budowa/klient-poprzedni/src/modele/podglad-promptu.ts
Treść nakładki bierze się z komendy odczytu tożsamości skutecznej, a nie ze sklejenia warstw w kliencie: własny porządek składania rozjechałby się z rdzeniem przy pierwszej zmianie reguł. Kategoria wymagana bez treści nie wstrzymuje uruchomienia, więc podgląd wylicza ją imiennie.

## budowa/klient-poprzedni/src/mission-control/sekcja-matrycy.ts
Pas relacji dokłada `mission-control.ts` pod matrycą — to osobny plik, bo niesie inną treść: powiązania, nie sesje. Jednym spojrzeniem widać, że procesy w różnych środowiskach biegną razem: karty strony głównej mówią, gdzie wejść, matryca — co już biegnie.

Kontrakt wskazuje środowisko sesji wyłącznie w `presence.environmentCode` sesji trwających — sesję bez tego odpisu matryca wypisuje osobno pod kolumnami zamiast zgadywać przypisanie; wiersz sesji poza środowiskami nie jest kontrolką wejścia, bo wejście wymaga środowiska.

Wskaźnik pracy w tle to kropka `dn-kropka--tetno` przy tytule; postęp etapów mieszka osobno, w kolumnie „Procesy w tle".

## budowa/klient-poprzedni/src/moduly/design/nadanie-etykiet.ts
Zbiera pełny zestaw etykiet jednego zasobu i oddaje go rdzeniowi komendą ustawienia etykiet.
Kontrolka stoi przy filtrze wykazu, bo pole zasila to samo zawężenie co pole etykiet żądania
odczytu wykazu zasobów.

Pole pokazuje stan bieżący zasobu przy każdym przewybraniu, a wyczyszczenie pola zdejmuje
wszystkie etykiety; pusty zestaw jest drogą udaną, nie odmową.

Etykietę przycina wspólny moduł przycinania pól, nie metoda standardowa: rdzeń zapisuje etykietę
co do znaku, a przycinanie przeglądarki zostawia znak NEL (U+0085), więc zestaw złożony z samego
NEL dałby etykietę niewidoczną o długości jednego znaku. Przycięcie jest wspólne z filtrem wykazu,
żeby etykieta nadana i szukana były jednym ciągiem znaków.

Rdzeń oddaje etykiety odczytane z bazy po zapisie, nie echo żądania, więc zdanie skutku stoi na
odpowiedzi, a różnica wobec zestawu zamówionego jest odmową.

Cisza kanału nie jest odmową nadania — rdzeń mógł etykiety zapisać, mimo że gniazdo padło przed
odpowiedzią.

Wciągnięcie zasobu do zbioru modułu gubi wskazanie promptu źródłowego: odpowiedź nadania etykiet
niesie te same klucze co odczyt wykazu, bez klucza promptu, bo repozytorium nie ma przekładu klucza
wiersza promptu na kod kontraktu. Zasób świeżo wygenerowany traci więc po otagowaniu wskazanie
promptu w Preview Window; okno mówi o tym wprost w wierszu „Prompt źródłowy" zamiast dorabiać
wartość z poprzedniej odpowiedzi.

## budowa/klient-poprzedni/src/moduly/design/narzedzia-planszy.ts
Jedna odpowiedzialność: kontrolki zmieniające kompozycję i widok kanwy. Grupy dzielą się tym, co
narzędzie zmienia: widok kanwy (powiększenie, przesunięcie, siatka), układ warstw zaznaczonych
(wyrównanie, rozmieszczenie), skład kompozycji (szablony, elementy pomocnicze, zaznaczenie
wszystkiego), czynności bez drogi w kontrakcie (wersjonowanie, eksport, kursor współpracy).

Grupa bez drogi w kontrakcie stoi w przyborniku, a nie poza nim, żeby komplet narzędzi był widoczny
razem z tym, które z nich nie mają wykonania — nazwane, klikalne i mówiące dlaczego.

Liczby w presetach ramek nie są wzięte z cudzych urządzeń: to punkty łamania kierunku projektowego
platformy, te same, na których stoi cały jej układ. Dzięki temu ramka makiety ma dokładnie tę
szerokość, przy której produkt zmienia postać, a nie szerokość telefonu, który akurat był
w sprzedaży. Wysokości preset nie ustawia — punkty łamania jej nie rozstrzygają.

## budowa/klient-poprzedni/src/mission-control/zlozenie-danych.ts
Złożenie jest czystym przełożeniem stanu zebranego przez `zrodlo-pulpitu.ts` na kształt, który widok umie wyrysować. Każda wartość pochodzi z odczytu kontraktu, a miara bez źródła zostaje `null` i widok pokazuje ją jako stan pusty.

Kolumny matrycy pochodzą z rdzenia w całości — kod, nazwa, motto i kolejność z odczytu `environment.list`. Przed pierwszą odpowiedzią kolumn nie ma: nagłówek pulpitu mówi „oczekiwanie na rdzeń", a matryca pokazuje ten sam stan zamiast zastępczego kompletu kolumn zbudowanego z kodów modułów.

Treść karty środowiska strony głównej jest osobna od motta kolumny matrycy i nie jest tu powielana. Kolejność kolumn `Environment.order` odpowiada kolumnie `srodowisko.kolejnosc`; klient nie sortuje po nazwie ani po kodzie, bo porządek kart jest zapisany w bazie.

## budowa/klient-poprzedni/src/moduly/agents/archiwum-ekspertow.ts
Historia wersji ma własny panel i tu jej nie ma: archiwum odpowiada na pytanie, gdzie ekspert poszedł, a historia — co się z jego tożsamością działo. Dwa wykazy wersji w jednym oknie byłyby dwiema prawdami o tym samym. Kontrolka pyta rdzeń, a nie stałą: wywołanie dostaje wyłącznie czynność, którą rdzeń melduje przy powitaniu połączenia; pozostałe zostają kontrolką nazywającą brak. Dzięki temu panel mówi prawdę także przed rdzeniem starszym niż on sam — przy wdrożeniach on-premise to stan normalny. Archiwizacja nie jest wyłączeniem: wyłączony ekspert zostaje w bibliotece i da się go edytować, a wyłączenie znaczy „nie obsługuje okien”, nie „zeszedł z drogi”. Archiwum ma własne komendy — przywrócenie dotyczy pozycji, nie panelu, i stoi przy każdym wierszu wykazu, bez kontrolki zbiorczej, która musiałaby pytać, którego eksperta dotyczy.

## budowa/klient-poprzedni/src/moduly/agents/archiwum-ekspertow.ts (wiersz czynności)
Czynność, którą rdzeń melduje, dostaje kontrolkę wywołującą. Czynność bez drogi zostaje kontrolką klikalną, która naciśnięta nazywa brak dymkiem, zamiast milczeć albo być wygaszona. Czynność pozycji nie ma kontrolki zbiorczej — jej miejsce jest przy wierszu wykazu.

## budowa/klient-poprzedni/src/moduly/design/okna-warsztatow-designu.ts
Bez tych okien siedemdziesiąt pięć komend byłoby funkcjami, których Operator nie ma. Nastawy stoją
danymi, nie pięcioma plikami po jednym oknie: różnią się kodem katalogu rdzenia, tytułem, rolą,
grupą czynności i zdaniem objaśnienia, a poza tym są tym samym oknem.

Kody są kodami katalogu rdzenia. Bez wiersza w katalogu klient postawiłby okno, o którym rdzeń nie
wie: wykaz modułów zaniżałby zakres modułu, a pas uczciwości meldowałby, że okno stoi poza
katalogiem.

Pięć źródeł byłoby pięcioma połączeniami do tego samego kanału. Kanał wchodzi wprost, a nie ze
stanu modułu: stan oddaje źródła obszaru Design i zaplecza, a warsztaty idą dowolną komendą
z katalogu czynności. Dołożenie kanału do stanu tylko dla nich otwierałoby wszystkim oknom drogę
obok źródeł, które stan dla nich trzyma.

## budowa/klient-poprzedni/src/mobile/arkusz-drog.ts
Układ arkusza jest podporządkowany kciukowi: kontekst decyzji zdaniami u góry (czyta się go raz), drogi u dołu (dotyka się ich w biegu). Arkusz wjeżdża od dołu, bo tam sięga kciuk trzymający telefon jedną ręką.

Droga nieprzejezdna dostaje zdanie mówiące, czego brakuje, zamiast przycisku wyszarzonego, który obiecywałby przyszłe działanie. Przycisk przejęcia z pustym polem oddaje kwit ze zdaniem o niepodanym poleceniu, zamiast być martwą kontrolką.

## budowa/klient-poprzedni/src/aktualizacja/wykaz-wydan-widok.ts
Widok pokazuje tę samą chronologię co strona „Pobierz” (`budowa/witryna/tresc/pobierz.mjs`),
w tych samych siedmiu kolumnach: Wersja, Data, System, Plik, Rozmiar, Suma SHA-256,
Co się zmieniło. Oba widoki czytają jeden plik `wydania.json`; inny wykrój tych samych
danych po jednej ze stron dałby drugą prawdę o wydaniach.

Zdanie o wydaniu bez sumy brzmi tak samo po obu stronach, bo mówi rzecz prawdziwą u obu:
most do powłoki odmawia wywołania bez sumy (`most-aktualizacji.ts`, kod `wydanie-bez-sumy`),
a puste pole w tabeli wyglądałoby na brak danych, nie na przeszkodę.

Każdy stan kanału ma własne zdanie. Widok pokazuje to, co oddał `pobierzWykazWydan()`,
i ani słowa więcej: odczyt rozróżnia brak łączności, brak pliku pod adresem i pusty wykaz,
więc „nie ma jeszcze żadnego wydania” pada tylko wtedy, gdy kanał tak odpowiedział.

Widok nie zawiera odnośników, bo odnośnik wyprowadziłby okno aplikacji pod obcy adres.
Adres jest wypisany jako tekst do skopiowania, a pobieranie zostaje pod przyciskiem
banera albo na witrynie. Widok nie montuje się sam: oddaje element, a osadza go ten,
kto go przywołał — przycisk „Wykaz wydań” przy banerze.

## budowa/klient-poprzedni/src/moduly/design/okno-assets-panel.ts
Zasób wchodzi do wykazu dwiema drogami. Wgranie zasobu przyjmuje plik wskazany w oknie — klient
czyta jego bajty i oddaje je rdzeniowi; kontrolka stoi w oknie pierwsza. Generowanie oddaje bajty
z kanału obrazowego, a powstały zasób trafia do wykazu zdarzeniem zmiany, bez czynności w tym
oknie.

Zawartość wykazu przestawiają jeszcze trzy komendy: nadanie etykiet zasila filtr etykiet, ustawienie
ulubionego zasila przełącznik „Tylko ulubione", a usunięcie jest nadawcą rodzaju zmiany usunięcia;
obie ostatnie stoją w module czynności zasobu. Usunięcie idzie bez pytania „czy na pewno".

Odmowa odczytu nie gasi okna: odczyt zakończony odmową wchodzi w stan błędu, a okno zostaje czynne.
Zasób przysłany zdarzeniem zmiany wejdzie do wykazu mimo to, bo wciąga go stan modułu, nie ta
odpowiedź.

Czynności na zasobie wskazanym (ulubiony, usunięcie) idą zaraz pod wykazem, obok nadania etykiet —
wszystkie trzy dotyczą jednego zasobu wskazanego kartą.

Zerwanego gniazda nie tłumaczy żadne pole żądania, więc przy braku rozstrzygnięcia dopisek
zostaje pominięty.

## budowa/klient-poprzedni/src/mobile/ekran-interwencji.ts
Od otwarcia ekranu do wykonanej decyzji prowadzą dwa dotknięcia: karta, potem droga. Nagłówek melduje każdą odmowę odczytu osobno. Wykaz pusty to nie to samo co brak źródła — „Nic nie czeka” pada wyłącznie wtedy, gdy odczyty doszły; przy odmowach ekran mówi, że nie wie.

Zdarzenie `mobile.process.changed` nie ma po stronie rdzenia producenta (`handlers_mobile.go`), więc nasłuch przerysowania wykazu stoi na trzech zdarzeniach, które go mają — telefon nie ma być odpytywany palcem.

## budowa/klient-poprzedni/src/moduly/design/zapis-kompozycji.ts
Rozgłaszaniem zmian zajmuje się osobny zapis stanu kompozycji: zmiana układu jest
rachunkiem na liczbach, rozgłoszenie obsługą obserwatorów, rozdzielone czytają
się i sprawdzają osobno. Wstawienie kompozycji z rdzenia zastępuje układ, nie
scala go: kompozycja z rdzenia jest pełnym stanem planszy, a nie jej dokładką,
więc scalenie dałoby układ, którego nie ma ani na ekranie, ani w bazie. Okno
mówi wprost, że odczyt nadpisuje to, co ma na kanwie. Identyfikator wchodzi
razem z układem, bo bez niego kolejny zapis założyłby kompozycję nową zamiast
zaktualizować odczytaną i plansza rozmnażałaby się w bazie po jednej sztuce na
każde wejście w moduł. Licznik warstw przesuwa się ponad wczytane, żeby kod
warstwy dokładanej po odczycie nie zderzył się z kodem warstwy, która przyszła
z rdzenia. Kolumna zewnętrznego identyfikatora warstwy jest unikalna w całej
tabeli, nie w obrębie pojedynczej kompozycji, a wstawienie warstw idzie bez
scalania po konflikcie — kod z samego licznika okna wracałby po każdym
przeładowaniu do tej samej wartości i zderzałby się z warstwą zapisaną
wcześniej, co rdzeń odrzuca naruszeniem unikalności. Człon losowy identyfikatora
zapewnia unikalność, a numer porządkowy zostaje z przodu, żeby kod warstwy dało
się przeczytać w panelu warstw. Formaty społecznościowe stoją w jednym rzędzie,
a nie w siatce, żeby oglądać obok siebie kadr tego samego materiału w trzech
proporcjach. Szerokość ramki obszaru roboczego pochodzi z punktu łamania
produktu, a wysokość zostaje przy boku wyjściowym, bo punkty łamania opisują
wyłącznie szerokość; wysokość ustawia się osobno w inspektorze właściwości.

## budowa/klient-poprzedni/src/mobile/ekran-procesow.ts
Ekran niesie dokładnie te dwa elementy, które nazywa funkcja globalna Mobile: listę procesów z filtrem stanu oraz zestaw czynności przy pozycji — uruchom ponownie, zatrzymaj, wstrzymaj, wznów, zatwierdź, modyfikuj. Dwa dotknięcia dla czynności nieodwracalnej to nie utrudnienie: telefon nosi się w kieszeni, a zatrzymanie pracy, która biegnie bez Operatora, jest jedyną czynnością tego ekranu, której nie da się cofnąć niczym.

Filtr stanu zawęża po stronie rdzenia, polem `status` żądania, a nie po stronie widoku: wykaz przefiltrowany w oknie kłamałby o liczbie procesów, których rdzeń nie przysłał. Urządzenie mobilne nazywa się kartą sesji kanału — tak samo jak w kafelku stanu platformy, żeby rdzeń widział jedno urządzenie, a nie dwa.

Jedno sterowanie bywa widoczne w kilku wierszach naraz (kolejka, tura); odpowiedź opisuje jeden proces, o pozostałych rozstrzyga rdzeń. Napis przycisku zmienia się na czas uzbrojenia, żeby nie dało się go pomylić z pierwszym dotknięciem.

## budowa/klient-poprzedni/src/moduly/design/okno-design-board.ts
Panel akcji modułu wymienia kanwę swobodną, panel warstw, wyrównanie, siatkę, szablony układu,
adnotacje, wersjonowanie, eksport, zaznaczenie wielokrotne, kursor współpracy i bibliotekę
elementów pomocniczych, a kontrakt niesie jedną komendę zbiorczą aktualizacji planszy. Okno dzieli
to tak: zestawianie układu dzieje się na kanwie i jedzie do rdzenia jednym zapisem w polu warstw;
to, czego pole warstw nie unosi — wersje, eksport, obecność drugiego Operatora — jest nazwane
brakiem, nie pozorowane.

Potwierdzenie zapisu opisuje odpowiedź rdzenia, nie wysłane żądanie: odpowiedź niesie całą
kompozycję odczytaną po zapisie, z nazwą i wykazem warstw. Liczba warstw, która wróciła, jest
jedyną miarą tego, ile ich leży w rdzeniu; rozbieżność wobec wysłanego układu jest odmową, bo
kanwa pokazuje wtedy co innego niż baza.

Okno robocze modułu stoi w parze z oknem rozmowy. Profil modułu wymienia Design Board wśród okien
obowiązkowych, a zasób zlecony w Chat Window przychodzi zdarzeniem zmiany i kładzie się warstwą na
kanwie bez czynności w tym oknie. Tabliczka nad kanwą opisuje tę drogę także wtedy, gdy nic
jeszcze nią nie weszło.

Kolejny zapis po odświeżeniu bez wcześniejszego odczytu zakładałby kompozycję nową, bo
identyfikator kompozycji przepada razem z ekranem.

Bez zaznaczenia okno mówi, czego brakuje, zamiast milczeć: skrót naciśnięty bez skutku i skrót
niedziałający wyglądają dla Operatora tak samo.

Okno prowadzi jedną kanwę, więc wyboru między kilkoma planszami nie ma czym wyrazić; data
aktualizacji jest jedyną miarą świeżości, którą niesie kontrakt.

Faza wraca z ładowania, żeby odczyt w toku nie został przykryty stanem pustym.

Zapis pod czuwaniem nie ogłasza niepowodzenia zapisu, który mógł się w rdzeniu odbyć mimo
zerwanego gniazda.

## budowa/klient-poprzedni/src/mobile/okno-mobile.ts
Kafelek stanu platformy woła `mobile.status.get` i pokazuje odpowiedź rdzenia bez interpretacji; odmowa dotyczy wyłącznie kafelka i nie gasi niczego poza nim. Ekran interwencji stoi na komendach, które rdzeń obsługuje — `monitor.status`, `queue.list`, `window.list`, `window.state.get`, `role.list` — i niesie cztery drogi interwencji w dwóch dotknięciach. Przegląd zadań i procesów woła `mobile.process.list` i `mobile.process.control` — wykaz procesów wraz ze sterowaniem nimi, czyli jedyną czynność sprawczą okna nad procesem; czynność nieodwracalna mówi to przed wykonaniem.

Warstwy okna nie są trzema źródłami prawdy: wszystkie czytają ten sam rejestr telemetrii procesów, z którego czyta `monitor.status`. Różnią się drogą i tym, co potrafią — obraz interwencji rozstrzyga, wykaz procesów steruje.

Trzecia warstwa okna czyta inną rodzinę komend niż obraz interwencji i niesie to, czego obraz nie ma — sterowanie procesem. Odmowa jednej warstwy nie gasi pozostałych dwóch. Kafelek stanu platformy, obraz interwencji i wykaz procesów mówią o tym samym rdzeniu trzema różnymi rodzinami komend.

## budowa/klient-poprzedni/src/moduly/design/okno-prompt-builder.ts
Pola strukturalne mieszkają w modalu, bo klasa okna kreatora żąda nośnika modalnego w stanie
zamkniętym; sekcja okna zostaje widoczna zawsze, bo jest jednym z okien operacyjnych modułu.
Przycisk generowania stoi w obu miejscach i prowadzi do tej samej czynności.

Generowanie kończy się jedną z dwóch dróg i okno obsługuje obie. Rdzeń wysyła polecenie kanałem
obrazowym, odbiera fragment obrazu, odkłada bajty w magazynie pod sumą kontrolną, mierzy format
i wymiary z nagłówka utrwalonego pliku i zakłada wiersz zasobu — odpowiedź udana niesie zasoby
z treścią. Odmowa przychodzi przy braku kanału obrazowego w rejestrze, kanale nieczynnym, kanale
tekstowym, braku poświadczenia albo odpowiedzi bez obrazu; niesie wtedy gotową treść polecenia,
więc okno pokazuje odmowę wraz z oddanym tekstem — złożenie promptu odbyło się także wtedy.

Odpowiedź udanego generowania ma w kontrakcie pole identyfikatora procesu, którego rdzeń nie
wypełnia, a zdarzenia postępu nie przychodzą. Bez tego pola pasek postępu jest więc wyciszany,
zamiast stać na „zlecenie przyjęte" po pracy już zakończonej.

Panel oddaje prompt z powrotem w pola kreatora, więc szablon da się użyć, a nie tylko obejrzeć.

Cisza kanału nie jest odmową generowania: rdzeń mógł zlecenie odebrać i wykonać, więc pasek
postępu milknie, a zdanie mówi o skutku nieznanym, nie o niepowodzeniu.

Odmowa wskazania (bez okna, bez tematu) różni się od odmowy braku silnika tym, że tam składać
nie było czego — złożenie promptu jeszcze się nie odbyło.

## budowa/klient-poprzedni/src/mobile/pozycje-decyzji.ts
Pozycja powstaje wyłącznie z tego, co rdzeń umie udowodnić odpowiedzią komendy — cztery dowody, cztery rodzaje pozycji: `queue.list` z kolejką w stanie `paused` (praca stoi na kolejce), `window.state.get` z `loop.stopped = true` (pętla stoi, z powodem), `monitor.status` z procesem `failed` (krok padł) i `monitor.status` z procesem `paused` (krok wstrzymany).

Kolejek eskalacji — przepływów wstrzymanych z pytaniem do człowieka — rdzeń nie wystawia: nie ma na nie ani komendy odczytu, ani zdarzenia. Pulpit mówi o tym wprost, a warstwa mobilna nie zamalowuje braku pozycjami zmyślonymi. Gdy rdzeń wykaz eskalacji wystawi, wejdzie on portem `port-kolejki-decyzji.ts`.

Kontekst decyzji jest częścią pozycji: Operator otwiera telefon na minutę i musi wiedzieć, na czym praca stoi, zanim cokolwiek naciśnie. Każde zdanie kontekstu niesie nazwę komendy, z której przyszło, więc da się sprawdzić jego źródło.

## budowa/klient-poprzedni/src/aktualizacja/wykaz-wydan.ts
Odczyt nie jest warunkiem pracy. Brak sieci, adres nieosiągalny, odpowiedź
nieczytelna i wykaz pusty znaczą dla banera to samo: nie wiadomo o żadnym
nowszym wydaniu, więc baner się nie pokazuje. Żadna ścieżka nie rzuca
wyjątkiem — kanał wydań nie może popsuć uruchomienia produktu.

Dla wykazu pokazywanego w oknie aplikacji te stany znaczą co innego i muszą
być rozróżnione, inaczej okno mówi „nie ma wydań” wtedy, gdy prawdą jest
„nie ma sieci”. Dlatego są tu dwie drogi odczytu: `pobierzWykazWydan()`
oddaje nazwany stan świata (typ `OdczytWykazu`), a `nowszeWydanie()` jest
nakładką na tamtą, oddającą `Wydanie|null` dla banera, któremu wystarczy
„jest co zakładać”.



Człony porównuje się liczbowo, nie napisami — inaczej `1.10` byłoby starsze
niż `1.9`. Człon nieliczbowy (np. `1.2.0-rc1`) schodzi do zera: wydanie próbne
nie ma prawa udawać nowszego niż wydanie właściwe.

Dwa miejsca, w których `parseInt` sam z siebie nie wystarcza: `Number.parseInt('4-rc1')`
oddaje 4, nie NaN, więc `1.2.4-rc1` wyszłoby równe `1.2.4` — dlatego człon musi być
liczbą w całości (`/^\d+$/`), inaczej jest zerem. Przedrostek `v` (`v1.9.0`) zbija
pierwszy człon do zera, więc jest zdejmowany z napisu, bo `v1.9.0` i `1.9.0` to
zapis tej samej wersji.



Rozpoznanie idzie po `navigator`, bo interfejs działa w oknie przeglądarkowym
także wtedy, gdy siedzi w powłoce natywnej — i to jest jedyna rzecz o systemie,
jaką strona wie bez pytania powłoki. Gdy nie wiadomo, oddawany jest pusty napis:
niewiedza nie jest podstawą do odrzucenia wydania.



Pole `system` deklaruje plik `wydania.json` i rozróżnia nim pliki strona „Pobierz”.
Bez tego sprawdzenia baner na Linuksie podałby plik `.exe` dla Windowsa, a powłoka
podstawiłaby go w miejsce AppImage. Wydanie bez pola `system` przechodzi: wykaz
jednosystemowy jest zgodny z takim kształtem pliku, a brak deklaracji to brak
wiedzy, nie deklaracja obcego systemu.
Rozróżnienie istnieje obok `nowszeWydanie()`, bo okno wykazu musi umieć powiedzieć
co innego przy zerwanym łączu, a co innego przy kanale, który wprost deklaruje brak
wydań. Zwinięte do jednego `null` obie sytuacje wyglądałyby dla okna tak samo
i przy braku sieci pisałoby ono „nie ma jeszcze żadnego wydania”.

Nazwy stanów (`brak-lacznosci`, `odpowiedz-serwera`) są wspólne z kodami powłoki
w `budowa/desktop/src-tauri/src/aktualizacja/pobranie.rs` — ta sama rzecz nazywa
się tak samo po obu stronach mostu.

`brak-wydania-pod-adresem` stoi osobno od `odpowiedz-serwera`: „serwer odpowiedział,
że tego pliku nie ma” to informacja o kanale, a nie o łączu, i nie wolno wtedy
kazać sprawdzać połączenia.



Pierwsza pozycja wykazu jest najnowsza (tak składa go witryna), ale nie ufamy
kolejności — porównanie idzie po wszystkich pozycjach. Plik wydania musi być
wskazany: wydanie bez pliku jest wpisem historycznym, a nie czymś, co da się
zainstalować.

Funkcja jest nakładką na `pobierzWykazWydan()` gubiącą rozróżnienie stanów: baner
pyta o jedno — „czy jest co zakładać”. Brak sieci, 404 i pusty wykaz odpowiadają
na to tak samo, a baner ma się wtedy nie pokazać. Kto potrzebuje zdania o świecie,
woła `pobierzWykazWydan()` wprost.
## budowa/klient-poprzedni/src/moduly/design/okno-warsztatu-designu.ts
Pięć warsztatów (fotografia, wektor, druk, bazy zdjęciowe, publikacja) dzieli ten sam budowniczy,
bo dzielą to samo zadanie: wybór czynności przestawia pola, pola pochodzą z katalogu, żądanie
składa katalog, a odpowiedź jest opisana skutkiem. Pięć osobnych budowniczych byłoby pięcioma
miejscami, w których pomyłka w składaniu żądania mieszka osobno.

Kontrolki bierze wspólny budowniczy pól z warsztatu dokumentu Studio — ten sam przełącznik
rodzajów pól, bo katalogi opisują pola tymi samymi typami. Druga kopia tego przełącznika
rozjechałaby się z pierwszą przy pierwszym nowym rodzaju pola.

Materiał wchodzi z magazynu okna i wynik do niego wraca. Żadna czynność nie zmienia materiału
w miejscu — okno mówi to przy polu materiału, bo Operator ma wiedzieć, że pomyłka nie kosztuje go
zdjęcia źródłowego.

Odpowiedź jest opisana liczbą, nie słowem „gotowe": nowy zasób, liczba stron, udział punktów
przezroczystych, liczba kafli, dostawcy, którzy nie odpowiedzieli. Meldunek bez liczby nie
odróżnia czynności wykonanej od czynności przyjętej — a to jest wzorzec szkody, który ten moduł
ma w historii.

Bilans zamiast ciszy także w oknie: pole, które rdzeń wypełnił powodem niepowodzenia, ma być
widoczne w meldunku. Meldunek „gotowe" nad odpowiedzią z trzema nieudanymi rozmiarami byłby tym
samym kłamstwem, przed którym broni się rdzeń.

## budowa/klient-poprzedni/src/aod/awatar-aod.ts
Pływający awatar jest CAŁĄ powierzchnią funkcji Always On Display w stanie spoczynku:
pojedyncze koło przy prawej krawędzi obszaru roboczego, bez etykiety, bez ramki
kontenera, ponad całą powłoką aplikacji. Nie znika przy przełączeniu środowiska
ani modułu, bo warstwa, w której siedzi, leży poza obszarem podmienianym przez
moduł.

Żaden ze stanów nie jest niesiony samą barwą: przy każdym stoi etykieta dostępności
i tytuł, a stany „sugestia oczekująca” i „waga wysoka” niosą dodatkowo plakietkę
liczbową — stan nigdy samym kolorem.

Awatar nie zna kanału, kolejki ani dymka. Przyjmuje dwa wywołania zwrotne
(kliknięcie, kliknięcie podwójne) i cztery czynności nastawcze; co za nimi
stoi, rozstrzyga `warstwa-aod.ts`.

Zero jako liczba oczekujących sugestii chowa plakietkę samo z siebie, stan „ukryta
(zero)”. Plakietka ukryta ustawieniem wyciszenia idzie osobno przez
`ustawPlakietkeWidoczna`.

Przeglądarka wysyła `click` przed `dblclick`, więc bez tego odstępu każde kliknięcie
podwójne otwierałoby najpierw dymek. Odstęp jest krótszy niż czas reakcji na otwarty
dymek, więc pojedyncze kliknięcie nadal działa od ręki.
## budowa/klient-poprzedni/src/moduly/design/panel-metadanych.ts
Przekazanie do modułu docelowego idzie komendą przekazania kontekstu — jedyną zbudowaną w kontrakcie
drogą przeniesienia kompletu kontekstu. Wykaz modułów docelowych przychodzi z rdzenia, nie z kopii
katalogu po stronie klienta.

Paczka kontekstu nie ma pola zasobu wizualnego, więc zasób jedzie polem identyfikatorów dokumentów
wraz z poleceniem nazywającym go po imieniu — to przybliżenie kontraktu, nie jego pełne pokrycie.

Potwierdzenie opisuje okno, które wróciło, a nie moduł zamówiony w polu wyboru: przekazanie nie
sprawdza katalogu modułów — przepisuje identyfikator modułu docelowego do okna i oddaje odpowiedź
udaną także dla kodu, którego katalog nie zna. Jedynym polem mówiącym, gdzie zasób wylądował, jest
pole modułu w odpowiedzi okna; rozbieżność wobec zamówienia jest odmową.

Etykietowanie, wgranie i usunięcie zasobu mają własne kontrolki wykonujące.

## budowa/klient-poprzedni/src/moduly/design/zetony-systemu.ts
Panel żetonów nie kopiuje ani jednej wartości: wykaz niesie wyłącznie nazwy, role
które system wizualny nazywa, a wartość każdej czytana jest z dokumentu przy
każdym odświeżeniu. Dzięki temu panel pokazuje stan produktu, a nie jego opis —
przełączenie motywu zmienia tu wszystko bez linii kodu, a barwa poprawiona
w arkuszu motywu jest widoczna natychmiast. Konsekwencja jest zamierzona: żeton
wymieniony w wykazie, którego motyw nie definiuje, wychodzi na wierzch jako brak
definicji zamiast zniknąć, bo rozjazd nazw między modułem a motywem ma być
widoczny jako usterka produktu. Wykaz obejmuje żetony semantyczne, te po które
wolno sięgać komponentom; prymitywów skali szarości tu nie ma, bo sięganie po nie
wprost jest w tym produkcie zabronione. Moduł nie zapisuje żetonów i nie
nadpisuje motywu: motyw jest własnością powłoki, więc trwałego zapisu nie ma
gdzie odłożyć, i panel mówi to wprost zamiast udawać edytor. Przegląd selektorów
arkuszy idzie po arkuszach wczytanych do dokumentu; górna liczba selektorów
w wyniku kończy przegląd wcześniej, bo wykaz na kilkaset pozycji nie jest
odpowiedzią, tylko zrzutem. Zagnieżdżenie reguł ma znaczenie, bo definicje
motywu ciemnego stoją w regule warunkowej, więc przegląd płaski przeoczyłby
połowę produktu.

## budowa/klient-poprzedni/src/aod/cztery-stery.ts
Cztery stery decyzji stoją w jednym pasie przy wpisie: zatwierdzenie kroku,
wstrzymanie, konfiguracja Koordynatora, przejęcie bezpośredniego sterowania.

Ster ma przycisk tylko wtedy, gdy stoi za nim komenda kontraktu; inaczej wyświetla
zdanie nazywające granicę. Zatwierdzenia kroku kontrakt nie zna, a przestawienie
ogniska wymaga identyfikatora klienta z powitania połączenia, którego okno
nakładki nie otrzymuje. Żaden przycisk nie pyta o potwierdzenie i żaden nie jest
wyszarzany — rozstrzyga rdzeń, nakładka pokazuje odpowiedź.
## budowa/klient-poprzedni/src/mobile/procesy-mobilne.ts
Powierzchnia stoi na ekranie „Przegląd zadań i procesów", z listą procesów wraz z filtrem stanu i z zestawem czynności przy pozycji. Rodzina jest w rdzeniu wpięta i to zostało zmierzone, nie założone: `montaz_porty.go` wnosi do portu nawigacji ogniwo `ZWarstwaMobilna`, więc asercja w `handlers_mobile.go` przechodzi i trzy komendy mają obsługiwaczy. Zdanie z `zrodlo-interwencji.ts` o niewpiętej warstwie opisywało stan wcześniejszy i przestało być prawdziwe.

Wykaz procesów nie zastępuje obrazu interwencji. Obraz stoi na pięciu komendach czytających stan pracy z różnych stron (telemetria, kolejki, okna, role, bieg naprawczy); `mobile.process.list` jest jednym skrótem do procesów i tak jest tu użyty — jako druga, węższa droga do tego samego stanu, przynosząca to, czego obraz nie ma: sterowanie procesem.

Sterowanie nie ma w kontrakcie zdarzenia własnego, więc po każdym udanym poleceniu wykaz czyta się na nowo. Proces oddany w odpowiedzi opisuje jeden wiersz; o pozostałych rozstrzyga rdzeń, a nie widok, który je wcześniej widział.

Zatrzymanie i ponowne uruchomienie przerywają pracę, która biegnie: tego, co proces zdążył zrobić, żadne z nich nie odda. Zatwierdzenie i modyfikacja idą naprzód, nie w tył.

## budowa/klient-poprzedni/src/moduly/design/panel-zasobow.ts
Ma kształt panelu pomocniczego, więc staje w pasie paneli dowolnego modułu i w kolumnie paneli
sceny okien równoległych. Poza nim zasoby Designu są osiągalne wyłącznie wewnątrz złożenia modułu
Design.

Panel jest osobnym, gęstszym widokiem tych samych danych, a nie opakowaniem okna Assets Panel.
Tamto okno nie ma czynności zamknięcia — subskrypcję zmiany zasobu trzyma stan modułu, a zdejmuje
ją rozłączenie wołane przez moduł — i żąda całego stanu modułu wraz z komendami odczytu okien,
stanu okna, kanałów, modułów oraz kanwy Design Board, której poza modułem Design nie ma. Tu mieści
się jeden wiersz nagłówka, nie trzy pasy kontrolek okna operacyjnego.

Własnych reguł o zasobie panel nie pisze — pożycza komplet od modułu: źródło danych jako jedyna
warstwa wywołań odczytu i subskrypcji zmiany, zapis jako całe wciąganie zmian wraz z gałęzią
usunięcia, kartę zasobu i predykat frazy jako jedną kopię na moduł, stan okna jako trzy stany
obowiązkowe, tor komendy jako zdanie doklejane do odmowy. Oba widoki zbiegają się na tym samym
zdarzeniu rdzenia.

Panel czyta zasoby wszystkich okien Designu i nazywa to w stanie pustym. Okno gospodarza niesie
okno panelu pomocniczego, a nie okno modułu Design; pole okna żądania jest opcjonalne, a warunki
filtra dokładają warunek okna tylko wtedy, gdy pole przyszło. Podstawienie okna gospodarza
zawęziłoby wykaz do zasobów obcego okna, czyli najczęściej do pustki, więc panel pola nie
podstawia.

Panel nie oddaje zasobu gospodarzowi: komendy wstawiającej zasób Designu w rozmowę obcego modułu
kontrakt nie ma, a przekazanie kontekstu biegnie przeciwnie — z okna źródłowego do modułu
docelowego, otwierając tam okno. Nie generuje, nie nadaje etykiet i nie zapisuje kompozycji; te
czynności zostają w oknach modułu Design.

Rozdział po rodzaju zmiany idzie z modułu wraz z gałęzią usunięcia, żeby panel nie pokazał jako
obecnego zasobu, o którym rdzeń właśnie powiedział, że go nie ma.

Czytanie przy pierwszym odczycie zleciłoby odczyt zasobów dwa razy pod rząd.

## budowa/klient-poprzedni/src/moduly/agents/archiwum-ekspertow.ts
Historia wersji ma własny panel i tu jej nie ma: archiwum odpowiada na pytanie, gdzie ekspert poszedł, a historia — co się z jego tożsamością działo. Dwa wykazy wersji w jednym oknie byłyby dwiema prawdami o tym samym. Kontrolka pyta rdzeń, a nie stałą: wywołanie dostaje wyłącznie czynność, którą rdzeń melduje przy powitaniu połączenia; pozostałe zostają kontrolką nazywającą brak. Dzięki temu panel mówi prawdę także przed rdzeniem starszym niż on sam — przy wdrożeniach on-premise to stan normalny. Archiwizacja nie jest wyłączeniem: wyłączony ekspert zostaje w bibliotece i da się go edytować, a wyłączenie znaczy „nie obsługuje okien”, nie „zeszedł z drogi”. Archiwum ma własne komendy — przywrócenie dotyczy pozycji, nie panelu, i stoi przy każdym wierszu wykazu, bez kontrolki zbiorczej, która musiałaby pytać, którego eksperta dotyczy.

## budowa/klient-poprzedni/src/mobile/zrodlo-interwencji.ts
Obraz interwencji stoi na pięciu komendach czytających stan pracy z różnych stron, bo z nich składa się rozstrzygnięcie: gdzie stoi pętla, kto jest koordynatorem i co czeka w kolejce. Rodzina `mobile.*` żadnego z tych pytań nie zastępuje.

Rodzina `mobile.*` jest już wpięta i to zostało zmierzone, nie założone: `montaz_porty.go` wnosi do portu nawigacji ogniwo `ZWarstwaMobilna`, więc asercja w `handlers_mobile.go` przechodzi i trzy komendy mają obsługiwaczy. Rodzinę woła osobne źródło (`procesy-mobilne.ts`), bo niesie ona to, czego obraz nie ma — sterowanie procesem — a nie drugi odczyt tego samego.

Ekran, któremu odmówiono ról, dalej pokazuje kolejki — odmowa jednej komendy nie unieważnia pozostałych.

Bieg naprawczy pytany jest wyłącznie o okna koordynatorów, bo `loop` jest puste dla okna samodzielnego i wykonawczego — pozostałe pytania byłyby ruchem bez odbiorcy.

## budowa/klient-poprzedni/src/moduly/design/zrodlo-designu.test.ts
Sprawdzian pilnuje, czy każda komenda rodziny design ma drogę z okna do rdzenia,
przy czym wykaz oczekiwany nie jest tu przepisany — powstaje z wywołań źródła,
a porównywany jest ze stałymi kontraktu, więc komenda dołożona do kontraktu
i pominięta w oknie zostanie tu nazwana. Dróg z okna do rdzenia są dwie i
sprawdzian przechodzi obie: pierwsza to źródło obszaru, metoda na komendę
typowana kontraktem, tą drogą jadą okna warstwy pierwszej; druga to warsztaty,
jedna droga na komendę wskazaną katalogiem czynności, tą jadą okna warsztatowe,
bo osobna metoda na każdą z wielu komend byłaby wieloma miejscami na tę samą
pomyłkę. Pozostałe sprawdziany dotyczą rozstrzygnięć, które warstwa kliencka
podejmuje sama i które łatwo cofnąć nieuważną poprawką: pola opcjonalne idą do
rdzenia wyłącznie wskazane, bo pole wysłane „na wszelki wypadek" jest zdaniem
o woli operatora, którego operator nie wypowiedział. Żądanie w przelocie drugą
drogą jest puste, bo sprawdzian pyta o drogę, nie o kształt żądania — kształt
sprawdza rdzeń osobno, a odmowa walidacji z nazwą pola jest odpowiedzią, którą
kanał próbny i tak zwraca jako powodzenie.

## budowa/klient-poprzedni/src/moduly/design/pasek-uczciwosci.ts
Jedna odpowiedzialność: powiedzenie wprost tego, czego okna nie mówią same.

Pas nie liczy już niczego napisem. Liczba komend obszaru, liczba okien katalogu i to, czy rdzeń ma
uchwyt danej komendy, są tu mierzone odczytem z rdzenia — katalogiem okien i pokryciem komend.
Powód jest w tym module policzalny: zdanie o liczbie komend obszaru stało tu, gdy komend było już
osiem, i nic go z kontraktem nie łączyło. Napis o stanie produktu, który nie jest z produktem
połączony, staje się fałszem w dniu, w którym produkt się zmienia — i nikt go nie zdejmuje, bo
nikt nie wie, że skłamał.

Zostają tu wyłącznie zdania, których żaden odczyt nie zastąpi, bo nie dotyczą liczby, tylko
mechanizmu: bajty wygenerowanego obrazu nie mają drogi do przeglądarki — rdzeń generowanie
wykonuje i zasób zakłada z prawdziwą treścią, ale pole odsyłacza jest ścieżką w jego systemie
plików; pytanie o Preview Window pozostaje otwarte — katalog rdzenia niesie jedną definicję okna
i dwa przypięcia do modułów, a to, czy Studio i Design mają kiedyś zejść się w jedno okno
konfigurowalne, jeszcze nikt nie rozstrzygnął.

Pozycja druga brakiem nie jest: pierwsza mówi, czego moduł nie umie, druga — czego nikt jeszcze
nie rozstrzygnął.

Wykaz idzie do pomiaru pokrycia, więc mówi nie o tym, co kontrakt niesie, ale o tym, czy rdzeń ma
dla każdej z nich uchwyt.

Zdanie przy każdej pozycji zmieniło się samo, bez dotykania tego wykazu: pomiar czyta wykaz komend
rdzenia i mówi teraz o braku uchwytu, a nie o braku nazwy. To są dwa różne stany i pas ich nie
zlewa — kontrakt komendę ma, rdzeń jej jeszcze nie obsługuje.

## budowa/klient-poprzedni/src/moduly/agents/biblioteka-ekspertow.ts
Kontrakt nie ma komendy kopiującej eksperta, więc duplikat powstaje z komendy zakładającej treścią oryginału i z komendy przypisującej skill dla każdej jego umiejętności. Konektory nie idą do kopii: ich definicję odczytuje dziś komenda, dla której rdzeń nie ma jeszcze uchwytu, a przepisanie samych identyfikatorów dałoby wpisy bez treści. Widok mówi o tym wprost po każdym duplikowaniu — i o tym, że jest to stan przejściowy, nie granica projektu. Zawężanie jest podzielone między rdzeń a przeglądarkę, bo kontrakt dzieli je tak samo: frazę wyszukiwania przyjmuje komenda odczytu biblioteki, więc jedzie do rdzenia i wraca węższym wykazem. Zasięg widoczności i stan czynności są polami bytu, który już przyszedł — zawężenie po nich w przeglądarce nie pyta rdzenia po raz drugi o to, co klient trzyma w ręku.

## budowa/klient-poprzedni/src/moduly/agents/biblioteka-ekspertow.ts (duplikowanie eksperta)
Mapa pusta i mapa nieustawiona to dwie różne rzeczy: dopóki odczyt nie wrócił, karty nie pokazują plakietki wcale — zero wpisane z ciszy byłoby orzeczeniem, którego nikt nie wydał. Imię własne i favikon idą do kopii razem z resztą tożsamości: bez nich duplikat wracałby w wykazie modeli bez znaku i bez imienia, choć powielany ekspert oba miał. Odstępstwo kopii, które wróciłoby do dopisywania, pracowałoby na innym prompcie systemowym niż powielany oryginał — a widać to dopiero po treści odpowiedzi modelu. Zasięg i pamięć kopii, które wróciłyby do stanu wyjściowego, dałyby eksperta widzianego szerzej niż powielany i czytającego pamięć, której tamten nie czyta. Kopia już jest w rdzeniu, więc wykaz trzeba odświeżyć tak samo jak po duplikowaniu udanym; odmowa przypisania umiejętności idzie po odświeżeniu, żeby jej nie przykryło zdanie o powodzeniu. Pustka po zawężeniu znaczy co innego niż pusta biblioteka: Operator ma wiedzieć, że eksperci są, tylko filtr ich nie przepuścił — inaczej sięgnąłby po założenie nowego eksperta zamiast po filtr.

## budowa/klient-poprzedni/src/modele/edytor-tozsamosci.ts
Tryb bierze się kolejno z zapisu i z trybu proponowanego przez kategorię, a gdy katalog nie podaje żadnego — z wartości domyślnej klucza `tozsamosc.tryb_domyslny`. Kategoria bez zapisu na osi czynnej bierze treść z osi szerszej i edytor podaje to wprost, zamiast pokazywać puste pole bez wyjaśnienia.

Zmiana kategorii albo osi znaczy inną treść w edytorze; samo przeliczenie stanu nie znaczy zmiany treści.

Jedyną drogą zapisu treści demonstracyjnej pozostaje komenda `identity.document.set`, wywoływana przyciskiem „Zapisz treść kategorii".

## budowa/klient-poprzedni/src/aod/dymek-sugestii.ts
Dymek otwiera się po lewej stronie awatara, nie przesuwa kolumn obszaru roboczego,
nie przyciemnia tła i nie zabiera ogniska klawiatury siłą. Zamknięcie klawiszem
`Esc` ani odejście ogniska nie zmieniają statusu sugestii.

Treść dymka mówi cztery rzeczy: co się stało (zdanie decyzji), jakiego jest
rodzaju, jak głośno wchodzi (waga) i czyja to ocena — rdzenia czy nakładki
(`WagaDecyzji`). Pod treścią stoją działania: te, za którymi stoi komenda
kontraktu, przyciskiem (`cztery-stery.ts`), a cykl życia sugestii — „Odłóż"
i „Odrzuć" — osobno, bo należy do nakładki, nie do rdzenia.

Przy wąskiej kolumnie obszaru roboczego dymek skraca się do jednego zdania
i działania „Rozwiń" otwierającego powierzchnię interakcji.
## budowa/klient-poprzedni/src/moduly/developer/indeks.ts
Moduł pracuje w oknie, nie w sesji: wszystkie komendy obszaru wymagają kodu
okna, a umowa ogólna widoku modułu niesie identyfikator sesji, więc złożenie
odracza montaż do chwili, gdy okno tego modułu jest znane rdzeniowi — bez tego
kodu przejście brałoby pierwsze okno w wykazie, także cudze. Układ wynika
z ról: w pasie górnym wiodące okno edytora kodu wraz z pomocniczym drzewem
projektu, które wskazuje plik odczytywany przez edytor; w pasie środkowym
zarządca repozytorium i monitor budowania; w pasie dolnym narzędzia
deweloperskie jako kolumna czterech integracji. Okno rozmowy modułu nie należy
do tego złożenia, bo jest bytem sesji i składa je warstwa rozmowy. Monitor jest
jednym oknem o dwóch częściach przełączanych zakładkami w nagłówku kolumny,
nie dwoma oknami, dlatego ma jeden kod katalogu rdzenia. Ostatni pas niesie
okna pomocnicze zbudowane oraz spis pozycji jeszcze nieistniejących wraz
z powodem każdej; terminal ma w tym spisie miejsce i nie jest budowany drugi
raz, nie powiela go też zakładka integracji kontenerów. Zdarzenie zmiany
budowania ma dwie subskrypcje, bo każda bierze co innego: stan modułu
unieważnia po nim drzewo i edytor, bo budowanie generuje pliki, a monitor
budowania bierze przyrost logu, którego stan nie przenosi, więc log narasta
wyłącznie ze zdarzenia. Warsztat jest drugim źródłem, bo ma innych odbiorców:
zakładki narzędzi deweloperskich, panel uruchamiania i debugowania, historię
przebiegów i pasek operacji edytora kodu — jedna umowa na wszystko kazałaby
drzewu projektu przyjmować zależność od wielu metod, z których używa niewielu.
Kolumna zakładek integracji dev tools jest rozszerzeniem bocznym obszaru
roboczego, nie sąsiadem zarządcy repozytorium, i potrzebuje pełnej szerokości
na wykaz zależności zewnętrznych. Pas okien pomocniczych zamyka się pierwszy,
bo trzyma subskrypcję strumienia, która żyje niezależnie od stanu modułu
i po zejściu ze sceny nikt by jej nie zdjął.

## budowa/klient-poprzedni/src/moduly/design/plyta-podgladu.ts
Pokazuje jeden zasób takim, jakim rdzeń go zna. Powierzchnia ma trzy stany: rdzeń podał adres —
płyta wstawia obraz i mówi, skąd go bierze; rdzeń podał adres, a treść spod niego się nie
wczytała — płyta nazywa przyczynę, zamiast zrzucać ją na przeglądarkę; rdzeń adresu nie podał —
płyta mówi, czego nie ma, i nie dorabia obrazu.

Drugi stan jest stanem każdego zasobu z generowania: rdzeń wypełnia pole adresu ścieżką w swoim
systemie plików, a nie adresem do pobrania. Obraz powstaje poprawnie i mimo to się nie wyświetli,
dopóki droga po treść zasobu — opisana już w kontrakcie i wspólna całemu magazynowi — nie dostanie
uchwytu w rdzeniu.

Wykaz pól mówi także o tym, czego nie ma: wartość pusta to nie pustka w wierszu, tylko zdanie
„rdzeń tego nie podał". Format nieznany i format nieistniejący to dwa różne stany.

Dla pola prompt źródłowy zdanie mówi więcej niż „nie podał": klucza nie niesie żadna odpowiedź
rdzenia, bo nie ma on przekładu klucza wiersza promptu na kod kontraktu i nie zgaduje go. Samo
„rdzeń tego nie podał" kazałoby sądzić, że zasób prompt zgubił.

## budowa/klient-poprzedni/src/aod/kolejka-decyzji.ts
Kolejka decyzji czekających jest magazynem bez DOM i bez kanału. Trzyma, co czeka,
od kiedy, czego dotyczy i skąd o tym wiadomo; drogi wyjścia zostają poza magazynem,
bo zależą od stanu rdzenia w chwili czynności (`cztery-stery.ts`).

Dosypują dwa źródła: odczyt nadrabiający `monitor.status` wnosi to, co stanęło przed
otwarciem okna, a sygnały `progress.changed` i `window.state.changed` — to, co staje
przy otwartym oknie.

Powtórzenie tego samego powodu nie zeruje `czekaOd`; zegar rusza od nowa dopiero
przy zmianie powodu. Kluczem jest identyfikator procesu, a gdy sygnał go nie niesie
(`window.state.changed`) — okno z przedrostkiem `okno:`. Sygnał o biegnącym procesie
zdejmuje wpis natychmiast. Wykaz idzie po `czekaOd` rosnąco, a przy równych chwilach
rozstrzyga klucz, żeby nie migotał.

Porównywane są pola rysowane, a nie całe struktury: `monitor.status` przychodzi
co odczyt i przy identycznej treści nie ma powodu przerysowywać wykazu —
przerysowanie gubiłoby ognisko klawiatury na sterach.
## budowa/klient-poprzedni/src/moduly/developer/katalog-komend.ts
Zdania stanu pustego orzekają o tym, czym rdzeń rozporządza — napis wpisany na
stałe byłby stanem wiedzy z chwili wpisania i nie zdjąłby się sam w dniu, w
którym rdzeń komendę dostanie, dlatego zdania budują się z rejestru rdzenia.
Źródłem jest powitanie połączenia, nie kontrakt: kontrakt mówi, co obiecano,
a powitanie rdzenia mówi, co rdzeń naprawdę zarejestrował. Różnica między nimi
to klasa usterki „martwy port": komenda obecna w kontrakcie i uchwyt obecny
w kodzie, ale w złożonym rdzeniu żadnej rejestracji, więc wywołanie wraca
odmową nieznanej komendy. Odczyt idzie raz, przy montażu okna, bo rejestr
komend zmienia się wraz z wersją rdzenia, a nie w toku sesji. Wykaz komend
wołanych przez okna modułu jest pisany ręcznie z rozmysłem: to deklaracja, co
okna modułu naprawdę wołają, a porównanie jej z rejestrem rdzenia wykrywa dwie
różne usterki — komendę wołaną, a niezarejestrowaną (martwy port), oraz komendę
zarejestrowaną, a niewołaną (funkcję rdzenia bez drogi z okna). Wykaz
wywiedziony wprost z kontraktu nie ujawniłby żadnej z tych dwóch usterek.

## budowa/klient-poprzedni/src/moduly/agents/braki-modulu.ts
Wykaz stoi w module, a nie rozsypany po oknach, bo odpowiada na pytanie zadawane raz i o całość: czego ten moduł nie potrafi i po czyjej stronie brak leży. Nazwa komendy pochodzi ze stałej kontraktu, nie z napisu: napis przetrwałby zmianę nazwy w kontrakcie i zostawiłby zdanie o pozycji, której już nie ma pod tą nazwą, podczas gdy stała przerywa kompilację i każe poprawkę wykonać. Zdanie o każdej pozycji bierze się z odpowiedzi rdzenia, nie ze stałej wpisanej w moduł. Byt pokrycia rozstrzyga pięć stanów i odróżnia brak w kontrakcie od braku uchwytu w rdzeniu. Pozycja nie znika i nie jest wygaszona — kontrolka zostaje klikalna i po naciśnięciu nazywa stan pozycji.

## budowa/klient-poprzedni/src/moduly/agents/braki-modulu.ts (dwie pozycje wykazu)
Wykaz liczył czternaście pozycji, dopóki rodzina komend agenta nie miała uchwytów. Dwanaście z nich zeszło stąd nie dlatego, że przestały być potrzebne, lecz dlatego, że mają już drogę z okna do rdzenia: umiejętności i konektory w swoich zarządcach, podgląd wersji w panelu historii, licznik przypisań na karcie eksperta, a cztery grupy zakresu w panelu uprawnień. Zostają dwie i obie należą do rodziny komend rejestru kanałów modelu, nie do modułu Agents — moduł ich potrzebuje, ale nie jest ich właścicielem, a dobudowanie ich stąd byłoby wejściem w cudzy obszar.

## budowa/klient-poprzedni/src/moduly/agents/braki-modulu.ts (przerysowanie wykazu)
Liczby nie ma przed odpowiedzią rdzenia i nie jest to niedopatrzenie: policzenie braków z ciszy byłoby orzeczeniem, którego nikt nie wydał. Dopiero powitanie mówi, ile z tych komend rdzeń faktycznie rejestruje.

## budowa/klient-poprzedni/src/moduly/developer/panel-historii-budowan.ts
Build Output prowadzi przebieg bieżący: log narasta zdarzeniem, a okno pokazuje
go na żywo; historia mówi o przebiegach zakończonych i o tym, co po nich
zostało, wyniku testów i pokryciu kodu — to dwa różne pytania i dwa różne
czasy, jedno o tym, co się teraz dzieje, drugie o tym, co wyszło wtedy.
W dzienniku przebiegu zostaje ogon logu, bo budowanie dużego projektu ma
dziesiątki tysięcy wierszy; odpowiedź niesie znacznik przycięcia i okno go
pokazuje, żeby odróżnić „ostatnie wiersze" od „tyle ich było", a operator
szukający wiersza z początku budowania ma wiedzieć, że go tu nie ma. Przebieg,
w którym nikt nie uruchamiał testów, nie ma wyników, i okno pisze to wprost,
zamiast pokazać zero na zero — zero przy zerze wygląda jak powodzenie,
a znaczy brak pomiaru.

## budowa/klient-poprzedni/src/aod/rodzaje-sugestii.ts
Rodzaje sugestii są cztery: `doradztwo`, `konfiguracja`, `problem`, `kolejny_krok`.

Ten katalog nie dokłada rodzaju do kontraktu i nie zgaduje go z treści zdania:
struktura `AodSuggestion` kontraktu (`shared/contract.ts`) niesie `id`, `text`,
`commandType`, `windowId` i `createdAt`, pola rodzaju nie ma. Sugestia przychodząca
z rdzenia ma rodzaj nieznany i tak jest opisana. Rodzaj mają wyłącznie te sugestie,
których autorem jest nakładka — decyzje rozpoznane regułą (`rozpoznanie-decyzji.ts`),
bo tam nakładka wie, co rozpoznała.

Waga sugestii jest osobną osią wobec wagi rozpoznania (`WagaDecyzji`: pewna/sporna
z `rozpoznanie-decyzji.ts`). Waga ujawnienia mówi, jak głośno sugestia ma się
ujawnić; `WagaDecyzji` mówi, czyja to ocena. Obie żyją obok siebie, żadna nie
zastępuje drugiej.

Katalog działań każdego rodzaju jest opisem, nie wykonaniem: mówi, jakie działania
rodzaj niesie i co każde robi. Które z nich mają dziś za sobą komendę kontraktu,
rozstrzyga `cztery-stery.ts` przy konkretnej decyzji — i tam, gdzie komendy nie ma,
stoi zdanie nazywające granicę zamiast przycisku-atrapy.

Wartość „waga sugestii ujawnianej samoczynnie — wysoka” jest ustawieniem
konfiguracyjnym; do czasu, aż kontrakt poniesie ustawienia zasięgu Always On
Display, obowiązuje wartość domyślna.
## budowa/klient-poprzedni/src/moduly/design/pola-promptu.ts
Jedna odpowiedzialność: zebranie treści promptu z pól i oddanie jej jako prompt kontraktu. Siedem
pól odwzorowuje pola promptu — temat, styl, kompozycję, oświetlenie, paletę, proporcje, wykluczenia
— bez pola spoza kontraktu i bez pominięcia któregoś z nich. Każde niesie dymek objaśnienia, bo
jest elementem konfiguracji.

Biblioteka stylów jest podpowiedzią, nie wykazem zamkniętym: kontrakt nie ma komendy katalogu
stylów, więc podpowiedź składa się ze stylów użytych w tej sesji, a pole pozostaje otwarte.

Stoi osobno, bo etykieta niesie nastawę bieżącą, a nie samą nazwę pola. Nazwa i nastawa pochodzą
z jednego miejsca, inaczej rozjadą się przy pierwszej zmianie.

Pole kompozycji stoi poza wykazem rodzajów, bo kompozycja powstaje z układu warstw Design Board,
a nie z promptu. Rodzaje audio, wideo i archiwum kontrakt zna, ale kanał obrazowy oddaje wyłącznie
fragment obrazu — wskazanie ich opisałoby bajty obrazu nazwą innego rodzaju.

Ta sama wartość silnika idzie także w pole silnika promptu, ale w roli opisowej: rdzeń wpisuje
prompt do wiersza w całości, więc pole zostaje śladem w zapisie promptu — po nim poznać, którym
kanałem zasób powstał.

Bez przycięcia temat złożony z samego znaku NEL (U+0085) przechodzi przez przycięcie standardowe
cało, mija sprawdzian „temat jest wymagany" i zakłada w rdzeniu zasób z tematem niewidocznym.

Ukrycie kanału tekstowego kazałoby Operatorowi szukać kanału, który założył, a wybieralność
prowadziłaby prosto w odmowę rdzenia.

Operator czyta etykietę, a nie identyfikator kanału; wpisanie identyfikatora dałoby etykietę
mówiącą co innego niż rozwinięty wykaz.

Rejestr bez ani jednego kanału obrazowego jest stanem, po którym generowanie odmówi — Operator
widzi to przed naciśnięciem przycisku generowania.

## budowa/klient-poprzedni/src/aod/sekcja-decyzji.ts
Pusta kolejka jest stanem poprawnym i tak jest opisana, zamiast ostrzeżeniem.

Każdy wpis mówi cztery rzeczy: co czeka (proces, stan, etap), od kiedy (z policzonym
odstępem), czego dotyczy (okno, sesja, kolejka) i czym to wykryto — odczyt
nadrabiający i zdarzenie na żywo mają różną świeżość. Wpis rozpoznany regułą sporną
(`usterka`, `bez-ruchu`) jest oznaczony jako ocena nakładki: rdzeń pojęcia decyzji
nie ma.

Kolejka jest magazynem w pamięci okna, zasilanym odpowiedzią rdzenia i zdarzeniami;
zamknięcie okna nic nie utrwala.

`MonitorStatus` niesie proces, okno i sesję, ale nie nazywa kolejki. Dojście
prowadzi więc przez `queue.list`: najpierw kolejka obsługująca to okno, potem
kolejka tej sesji. Przy niejednoznaczności zwracane jest `undefined`, a ster
wstrzymania powie, że kolejki nie dopasowano.
## budowa/klient-poprzedni/src/moduly/developer/panel-run-debug.ts
Monitor jest jednym oknem o dwóch częściach: Build Output prowadzi budowanie
i testy, Run & Debug prowadzi konfiguracje uruchomień i debugger krokowy,
a obie części przełącza pas zakładek w nagłówku kolumny. Punkt przerwania
stawia się w pliku i wierszu, zanim cokolwiek ruszy — należy do okna, a nie do
sesji, i przeżywa kolejne biegi; rdzeń trzyma go w bazie i podaje adapterowi
przy starcie, a panel pokazuje komplet punktów okna, bo margines edytora
rysuje wszystkie naraz. Zatrzymanie na punkcie zgłasza debugowany proces, a nie
operator, więc panel po każdym kroku pyta rdzeń o stos wywołań, zamiast
rysować stan z samego naciśniętego przycisku, co pokazywałoby stan życzeniowy.
Panel nie pokazuje przycisków kroku, dopóki sesji nie ma, bo wykonanie
wyrażenia wymaga ramki, a ramka istnieje tylko w programie zatrzymanym.
Kontrakt nie niesie zdarzenia zatrzymania, więc panel odczytuje stan na
żądanie, zamiast udawać widok na żywo.

## budowa/klient-poprzedni/src/aod/sekcja-obecnosci.ts
Rejestr obecności dotyczy procesów przypiętych do nakładki komendami
`aod.observe.attach` i `aod.observe.detach`. Obie komendy zwracają
`attachedProcessIds` po zmianie, więc wykaz rysuje odpowiedź rdzenia, a nie
przewidywanie klienta: po odmowie wykaz zostaje niezmieniony i nie trzeba
powtarzać odczytu stanu.

Odpięcie idzie bez pytania o potwierdzenie, a puste pole identyfikatora procesu
wolno wysłać — pustą wartość ocenia rdzeń (`validation_failed`). `detach`
procesu nieprzypiętego kończy się `not_found`; to stan poprawny, nie awaria,
i zdanie odmowy z `odmowy-aod.ts` mówi to wprost.
## budowa/klient-poprzedni/src/aod/sekcja-podpowiedzi.ts
Rdzeń układa podpowiedzi od bytu najwęższego do platformy — najpierw okno
ogniskowane, potem sesja, na końcu platforma. Sekcja tej kolejności nie
przestawia; numer porządkowy listy pokazuje ją wprost.

Przy pozycji stoi `commandType`, czyli nazwa proponowanej komendy. Przycisku
„wykonaj” nie ma: kontrakt daje nazwę komendy, ale nie daje jej żądania.

Sugestie dzielą się na cztery rodzaje — `doradztwo`, `konfiguracja`, `problem`,
`kolejny_krok` — i każdemu przypisany jest własny komplet działań. Struktura
`AodSuggestion` kontraktu pola rodzaju nie niesie, więc sugestia przychodząca
z rdzenia ma rodzaj nieznany. Sekcja mówi to wprost i wypisuje katalog obok
wykazu, zamiast zgadywać rodzaj z treści zdania — zgadnięty rodzaj podstawiłby
cudzy komplet działań pod cudzą sugestię. Rodzaj mają wyłącznie decyzje
rozpoznane przez samą nakładkę, w sekcji decyzji czekających.

Wykaz pusty jest stanem poprawnym.

Katalog rodzajów sugestii stoi w powierzchni interakcji na stałe, a nie przy
pozycji, właśnie dlatego, że pozycji nie da się do rodzaju przypisać: kontrakt
nie niesie tego pola. Operator widzi więc, jakie rodzaje funkcja zna i jakie
działania każdy niesie, i widzi zarazem, dlaczego przy pozycji rodzaju nie ma.
## budowa/klient-poprzedni/src/moduly/agents/drzewo-narzedzi.ts
Komponent jest sterem, nie wyświetlaczem: etykieta uchwytu niesie wskazany kod, kliknięcie rozwija wybór, wybór zmienia nastawę. Nastawa jest jedna i dzieli ją z polem otwartym Skills Managera — wpisanie kodu ręcznie przestawia uchwyt, wskazanie w drzewie wypełnia pole. Pole otwarte zostaje obok drzewa, ponieważ lista umiejętności w kontrakcie jest listą napisów bez narzuconego słownika, a serwer narzędzi rozpoznaje też kody, których katalog kontraktu nie zna. Liście są dwojakie: pierwszy liść każdej gałęzi wskazuje cały obszar jednym kodem, pozostałe — pojedyncze narzędzia. Obszar jest jednostką doboru, bo pozycji jest ponad dwieście. Gałąź zwinięta nie jest brakiem: widoczne są nazwy obszarów wraz z liczbą narzędzi, a pozycje odsłania dopiero wejście w obszar.

## budowa/klient-poprzedni/src/moduly/agents/drzewo-narzedzi.ts (budowa drzewa)
Znacznik przypisania znakuje pozycje, które ekspert już ma, i idzie zdaniem w opisie, a nie samym wskaźnikiem wyboru: wskaźnik wyboru niesie nastawę tej kontrolki, a przypisanie jest czymś innym — stanem eksperta. Zlanie obu w jeden znacznik kazałoby uchwytowi pokazywać naraz kilkanaście wartości, czyli przestać być sterem nastawy.

## budowa/klient-poprzedni/src/moduly/developer/widok-przebiegu.ts
Fragment nie zna ani źródła, ani stanu modułu: dostaje przebieg, zebrany log
i zawężenie, oddaje element. Zawężenie jest czynnością wyłącznie kliencką nad
materiałem, który już przyszedł — log narasta ze zdarzenia i nie ma komendy,
którą dałoby się dopytać rdzeń o wiersze pominięte, więc szukanie odbywa się
w tym, co okno usłyszało, a widok mówi to wprost, zamiast pozorować
przeszukanie całości. Zgłoszenie ze ścieżką jest przejściem do pliku: nowej
komendy to nie wymaga, bo wystarczy wskazać plik we wspólnym stanie modułu,
a edytor kodu otworzy go sam. Numer wiersza zostaje w napisie i niczego nie
otwiera, bo polecenie otwarcia pliku nie ma pola wiersza, a pole edycji nie ma
numeracji — przejście otwiera plik, nie miejsce w pliku, i przycisk mówi
dokładnie tyle.

## budowa/klient-poprzedni/src/moduly/design/skroty-designu.ts
Skąd biorą się kombinacje: dokumentacja modułu wymienia skrót klawiszowy jako drogę równorzędną
kliknięciu — przy adnotacjach na kanwie, przy powiększeniu podglądu i przy całej warstwie czwartej
— ale konkretnych kombinacji nie podaje. Kombinacje są więc rozstrzygnięciem projektowym tej
budowy, podjętym wedle jednej reguły: bierzemy to, co ten produkt już związał w module Translate,
żeby Operator przechodzący między modułami nie uczył się dwóch układów klawiatury. Czynności, dla
których dokumentacja skrótu nie przewiduje, skrótu tu nie dostają — wymyślanie ich byłoby
dokładaniem funkcji.

Nasłuch wisi na elemencie modułu, nie na dokumencie. Moduł znika z drzewa przy zejściu ze sceny
i nasłuch znika razem z nim; nasłuch dokumentu trzeba by odpinać osobno, a pierwszy przeoczony
byłby wyciekiem.

Skutek uboczny tej decyzji jest nazwany, nie przemilczany: skrót działa, gdy ognisko stoi wewnątrz
modułu. Poza modułem klawisze należą do powłoki.

Dwa skróty stoją w wykazie, a moduł ich nie wiąże. Powiększenie podglądu należy do Preview Window,
którego wytwórnia leży poza katalogiem tego modułu; sterowanie pętlą należy do Execution Loop
Window, okna wspólnego platformy. Wykaz mówi to wprost, zamiast pomijać pozycje i sugerować, że
skrótów nie ma.

Ctrl i Cmd są tu równoważne, tak jak w zapisie skrótów dokumentacji: na komputerach Apple
modyfikatorem polecenia jest klawisz Meta.
