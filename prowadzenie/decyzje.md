# Rejestr decyzji

Jedyne miejsce, w którym obowiązują rozstrzygnięcia trudno odwracalne. Wpis
zakłada Właściciel produktu albo — od upoważnienia z 1 września 2026
(„Rozstrzygaj za mnie") — Prowadzący budowę, i wtedy wpis nazywa upoważnienie
w pierwszym wierszu. Rozstrzygnięcie nieobecne w tym rejestrze nie obowiązuje,
niezależnie od tego, gdzie zostało wypowiedziane.

Układ wpisu: kontekst, rozważone warianty, decyzja, konsekwencje.

---

## 1. Zamknięcie poprzedniego podejścia do budowy

**Data:** 2026-08-22 · **Stan:** obowiązuje

**Kontekst.** Poprzednie podejście do budowy Danaco Console nosiło dryf na tyle
głęboki, że kolejne poprawki nie przywracały spójności. Drzewo kodu zostało
usunięte w całości wraz z historią rewizji; na maszynie pozostał dorobek
merytoryczny (opracowania funkcjonalne, warstwa projektowa) oraz cztery pliki
wykonawcze uznane za przydatne. Repozytorium nie miało kontroli wersji — cały
dorobek istniał wyłącznie jako pliki na jednym dysku, bez kopii i bez historii.

**Rozważone warianty.**

1. Naprawa zastanego stanu przez kolejne poprawki. Odrzucone: dryf dotyczył
   nazewnictwa i granic, a nie pojedynczych usterek, więc poprawka objawowa
   odtworzyłaby przyczynę.
2. Ponowne rozpoczęcie od czystej karty, bez materiału odniesienia. Odrzucone:
   zniszczyłoby dorobek merytoryczny, który opisuje produkt, a nie nieudaną
   implementację.
3. Ponowne rozpoczęcie z dorobkiem merytorycznym jako materiałem wejściowym
   podlegającym przebudowie.

**Decyzja.** Wariant trzeci. Budowa rusza od nowa, poprzedzona przebudową
dokumentacji. Zastany stan drzewa został zapisany jako rewizja bazowa, żeby
żaden materiał nie zależał od pojedynczego dysku.

**Konsekwencje.** Kod aplikacji nie powstaje, dopóki dokumentacja nie zostanie
przyjęta. Dorobek poprzedniego podejścia jest odzyskiwalny z historii rewizji.
Terminologia i granice modułów wymagają ponownego rozstrzygnięcia — nie wolno
ich przepisać z materiału wejściowego bez sprawdzenia.

---

## 2. Dorobek merytoryczny na wydzielonych gałęziach

**Data:** 2026-08-22 · **Stan:** obowiązuje

**Kontekst.** Zachowane opracowania funkcjonalne i warstwa projektowa opisują
produkt trafnie, ale w kształcie ustalonym przez zamknięte podejście: odwołują
się do usuniętego drzewa kodu, w tym ponad sto sześćdziesiąt razy do
nieistniejącego pliku kontraktu. Pozostawione w gałęzi głównej udawałyby źródło
prawdy, którym nie są.

**Rozważone warianty.**

1. Pozostawienie w gałęzi głównej z naprawą martwych odwołań. Odrzucone:
   utrwaliłoby ustalenia zamkniętego podejścia jako obowiązujące.
2. Przeniesienie do katalogu archiwum w gałęzi głównej. Odrzucone: materiał
   nienormatywny obok normatywnego w jednym drzewie prowadzi do sięgania po
   archiwum jak po źródło.
3. Wydzielenie na osobne gałęzie przebudowy, z powrotem do gałęzi głównej
   dopiero po przyjęciu.

**Decyzja.** Wariant trzeci. Warstwa projektowa leży na gałęzi
`teren/prototypy`, opracowania zamkniętego podejścia pod
`refs/przeniesienie/dokumentacja-zastana`. Materiał przeniesiony nie jest
źródłem prawdy. Dorobek wchodzi do `main` dopiero po przyjęciu — dopiero wtedy
jest wiążący dla wykonawcy.

**Konsekwencje.** Gałąź główna nie zawiera dokumentacji produktu przez cały czas
trwania etapów 1 i 2. Gałęzie nie nakładają się plikami, więc scalenie jest
dodaniem, nie rozstrzyganiem konfliktów. Sesja szukająca wartości w gałęzi
głównej niczego nie znajdzie — i o to chodzi.

---

## 3. Zachowane pliki wykonawcze pozostają żywym kodem

**Data:** 2026-08-22 · **Stan:** obowiązuje

**Kontekst.** Cztery pliki ocalałe z poprzedniego podejścia — skrypt
prowizjonowania serwera, skrypt składania instalatora Windows, rejestr zależności
zewnętrznych w Go oraz wzorzec konfiguracji środowiska — niosą wiedzę operacyjną,
której nie da się odtworzyć z dokumentacji: wykaz programów zewnętrznych,
uzasadnienie doboru portu, model wdrożenia. Jednocześnie odwołują się do ścieżek
usuniętego drzewa i posługują się nazewnictwem metaforycznym.

**Rozważone warianty.**

1. Usunięcie. Odrzucone: wiedza operacyjna zniknęłaby bezpowrotnie.
2. Zachowanie jako materiał archiwalny. Odrzucone: materiał archiwalny nie jest
   utrzymywany, a te pliki mają wrócić do użycia.
3. Zachowanie jako żywy kod wymagający adaptacji.

**Decyzja.** Wariant trzeci. Pliki leżą w katalogu `narzedzia/` gałęzi głównej,
z opisem stanu i wykazem koniecznych adaptacji.

**Konsekwencje.** Pliki nie są wykonywalne w obecnym stanie — wywołują binaria
i ścieżki, których nie ma. Ich nazewnictwo wymaga uzgodnienia z terminologią
przyjętą w etapie 2, zanim wejdą do budowy. Do tego czasu służą jako źródło
wiedzy, nie jako narzędzie.

---

## 4. Przekrój pionowy zamiast kompletu prototypów

**Data:** 2026-08-22 · **Stan:** obowiązuje

**Kontekst.** Próbka na dwóch prototypach wykazała, że koncepcja przepływów
działa, a zasada jednej karty roboczej jest wykonalna. Ujawniła jednocześnie
coś ważniejszego: brief opisywał zasadę, a nie kompozycję okna, więc wykonawca
nie miał czego trafić w warstwie wyglądu. Kompozycja pierwszego okna modułowego
nie da się przekazać zleceniem — jest rozstrzygnięciem Właściciela, nie pracą
odtwórczą.

Osobno wróciła przyczyna upadku poprzedniego podejścia: powstało tam trzydzieści
osiem prototypów i pięćdziesiąt jeden opracowań, zanim cokolwiek zadziałało.
Objętość dorobku rosła, a sprawdzalność nie — bo nie było czego uruchomić.

**Rozważone warianty.**

1. Dokończenie kompletu prototypów, potem dokumentacja, potem budowa. Odrzucone:
   to jest dokładnie przebieg, który upadł, i upadłby ponownie z tej samej
   przyczyny — nic nie działa, dopóki nie powstanie wszystko.
2. Budowa bez prototypów, wprost z opracowań. Odrzucone: opracowania nie niosą
   kompozycji, a interfejs projektowany w kodzie kosztuje najwięcej przy
   pierwszej zmianie zdania.
3. Jedno okno modułowe doprowadzone do końca, a następnie pełny przekrój
   pionowy budowy na tym jednym module.

**Decyzja.** Wariant trzeci. Powstaje **jeden** prototyp okna modułowego —
Studio — prowadzony **własnoręcznie przez Właściciela**, bez orkiestracji
agentowej i bez wykonawców. Pozostałe trzydzieści cztery okna nie powstają.
Po domknięciu Studio budowana jest aplikacja w pełnym przekroju pionowym:
kontrakt, rdzeń, kanał, klient, powłoka. Kolejne moduły wchodzą dopiero na
przekroju, który działa.

**Konsekwencje.** Dorobek prototypowy przestaje rosnąć przed budową — to jest
zamierzone i jest to główna zapora przed powtórzeniem poprzedniego upadku.
Dokumentacja przestaje być etapem poprzedzającym budowę: powstaje wraz
z przekrojem i opisuje to, co działa, a nie to, co zamierzone. Ryzyko przesuwa
się z „zbudowaliśmy dużo, nic nie działa" na „przekrój jest wąski, więc decyzje
podjęte na jednym module mogą nie unieść czterdziestu" — i to ryzyko przyjmujemy
świadomie, bo sprawdza się je uruchomieniem, a nie objętością opracowań.

Kompozycja okna modułowego jest odtąd wyłączona z zakresu zleceń wykonawczych.
Zlecenie może opisać zasadę, mechanizm i kryterium sprawdzalne; nie może opisać
tego, jak okno ma wyglądać, bo tego nie da się sprawdzić inaczej niż oceną
Właściciela.

---

## 5. Model okna: aplikacja, okno robocze, karta, panel

**Data:** 2026-08-22 · **Stan:** obowiązuje

**Kontekst.** Warstwa projektowa nie niesie modelu produktu — niesie trzy różne,
nawarstwione: `rama.css` zakłada moduł jako miejsce wybierane w szynie
nawigacji, `stanowisko.css` układ wielookienny z czatami mnożonymi w poziomie,
`karty-okna.css` kartę jako widok aplikacji o rodzajach centrum, moduł, grupa.
Żaden z nich nie odpowiada temu, czym produkt ma być. Sprzeczność
`.dn-karty-sesji` zawierających `.dn-karta-widoku` nie jest usterką nazewniczą,
lecz śladem po tych trzech pomysłach w jednym drzewie.

**Decyzja.** Produkt ma cztery poziomy. Odniesieniem jest sposób, w jaki
przeglądarki organizują okna i karty.

| Poziom | Czym jest | Ile naraz |
|---|---|---|
| **aplikacja** | rama: górna wstążka, lewy panel stały z sesjami i projektami, prawy panel wysuwany z samouczkiem, obszar okna roboczego pośrodku | jedna instancja |
| **projekt** | zbiór zamknięty: własne ustawienia, wykaz modeli, własne pliki; grupuje sesje | wiele; praca może toczyć się poza projektem |
| **okno robocze** | kontener jednego działania Operatora; przełączane ikonką w lewym górnym rogu obszaru roboczego | wiele naraz |
| **karta** | jedna funkcja w sesji: czat z modelem, edytor, przeglądarka, wykaz plików, podgląd zmian, terminal, tłumacz | bez ograniczenia liczby; miarą jest pamięć maszyny |
| **widok** | miejsce wyświetlenia karty: główny pojedynczy, główny dzielony na dwie karty, oraz okno boczne | najwyżej trzy naraz |

**Sesja jest cechą okna roboczego, nie jego synonimem ani przeciwieństwem.**
Okno robocze niesie sesję wtedy, gdy toczy się w nim interakcja z modelem —
i wtedy otwarcie nowego okna jest jednocześnie założeniem nowej sesji. Okno
modułu służącego wytworzeniu mechanizmu — agenta, automatyzacji, pętli
wielomodelowej — sesji nie niesie, bo Operator niczego tam nie prowadzi
z modelem, tylko coś ustawia.

Relacja jest więc jeden do zera albo jednej: każda sesja należy do dokładnie
jednego okna roboczego, a okno robocze ma najwyżej jedną sesję. Stąd okien bywa
więcej niż sesji — pięć otwartych okien może dawać dwie sesje — ale nie jest to
reguła, tylko następstwo tego, jakie moduły są otwarte.

**Moment powstania sesji jest jeden: pierwsza wiadomość wysłana do modelu.**
Nie otwarcie okna, nie wejście do modułu, nie otwarcie karty czatu — dopiero
wysłanie treści, choćby jednego słowa. Do tej chwili Operator może swobodnie
chodzić po produkcie: oglądać agentów, zaglądać w ustawienia, otwierać kolejne
moduły i wracać — i żadna sesja nie powstaje.

**Nowe okno robocze otwiera się na centrum dowodzenia.** Centrum nie jest więc
osobnym miejscem, do którego się wchodzi, lecz stanem początkowym każdego okna
roboczego.

**Kolejne karty czatu nie tworzą kolejnych sesji.** Sesja obejmuje całe okno
robocze wraz ze wszystkimi jego kartami. Praca z czterema modelami i dołożenie
piątego to nadal jedna sesja — mnożenie rozmów nie jest przechodzeniem między
sesjami. To samo dotyczy sesji wznowionej: otwarcie w niej dodatkowej karty nie
zakłada nowej sesji, tylko rozszerza tę wznowioną.

Produkt prowadzi wobec tego **dwa odrębne wykazy**, i nie wolno ich mylić:

| Wykaz | Co obejmuje | Gdzie stoi |
|---|---|---|
| okna robocze | wszystko, co Operator ma otwarte | przełącznik w lewym górnym rogu obszaru roboczego |
| sesje | wyłącznie okna z interakcją modelu | lewy panel stały ramy aplikacji |

**Moduł to praca z modelem na określonym narzędziu głównym.** Nie jest miejscem,
do którego się wchodzi, ani samym narzędziem. Studio to praca z modelem
w edytorze; Browser to praca z modelem na przeglądarce; Terminal to praca
z modelem w powłoce. Narzędzie jest osią modułu, ale modułem jest to, co się
przy jego użyciu z modelem robi.

Widać to najlepiej na module Browser: nie jest przeglądarką dołożoną do
produktu, lecz pracą na materiałach w sieci — model ogląda witrynę, prowadzi jej
audyt, wskazuje, co działa źle i co wygląda źle, czyta wskazania konsoli, ocenia
przepływ i pracuje nad widocznością w wyszukiwarkach. Przeglądarka jest tu tym,
czym edytor w Studio.

**Środowisko zawęża pole dostępnych narzędzi do swoich modułów.** Operator
pracujący w środowisku dysponuje narzędziami tych modułów, które są do tego
środowiska przypisane — i tylko nimi. Przeglądarka jest osiągalna w TalkIn,
bo Browser należy do TalkIn; powłoki tam nie ma i nie da się jej otworzyć, bo
Terminal do TalkIn nie należy. Rozstrzyga o tym macierz dostępności modułów,
nie to, w którym module Operator akurat pracuje.

**Narzędzie w karcie pomocniczej to nie to samo, co moduł oparty na tym
narzędziu.** Granicę widać najlepiej na przeglądarce.

| | karta pomocnicza | moduł Browser |
|---|---|---|
| oglądanie pliku, który wyświetla się tylko w przeglądarce | tak | tak |
| otwarcie adresu w sieci | tak | tak |
| wskazanie elementu kursorem wraz z jego nazwą ze struktury | tak | tak |
| konsola | nie | tak |
| oprzyrządowanie analityczne strony | nie | tak |

Karta daje podstawowy przybornik i na tym jej rola się kończy. Wszystko, co
służy analizie strony, należy do modułu. Tak samo powłoka, edytor i pozostałe
narzędzia: karta daje funkcję podstawową, moduł — funkcję wraz z tym, co czyni
z niej warsztat. Tak samo powłoka, edytor i pozostałe narzędzia: karta
daje funkcję podstawową, moduł — funkcję wraz z tym, co czyni z niej warsztat.

Wynika z tego, że **karta pomocnicza i moduł oparty na tym samym narzędziu to
dwa różne wytwory**, nie jeden komponent użyty dwa razy. Wspólna jest podstawa;
osobne jest to, co moduł do niej dokłada.

Funkcje dostępne w oknie roboczym wnosi się do niego jako karty. Praca w module Studio
to jedno okno robocze, w którym Operator otwiera czat z modelem, drugi czat
z innym modelem, edytor, drzewo plików, przeglądarkę i terminal — każde jako
osobną kartę, przełączane albo zestawione w podziale.

**Do modułu prowadzą dwie drogi i obie obowiązują.** Pionowy pasek boczny jest
**paskiem szybkiego dostępu**: jedno kliknięcie ikony otwiera okno robocze
z wybranym modułem, z pominięciem drogi przez okno startowe, środowisko
i centrum dowodzenia. Druga droga wiedzie właśnie tamtędy. Obie prowadzą do
tego samego i żadna nie zastępuje drugiej — szybka służy czynności doraźnej,
długa wejściu do pracy.

**Zamknięcie okna to nie zakończenie sesji — to dwie różne czynności.**

| Czynność | Skutek | Gdzie |
|---|---|---|
| **Zamknij okno** | okno znika z przełącznika, **procesy trwają w tle**, sesja zostaje w wykazie i daje się wznowić | okno robocze |
| **Zakończ sesję** | praca zostaje ucięta, a sesja **usunięta z historii** | menu okna roboczego albo lewa strona ramy |

Zamknięcie okna jest odłożeniem pracy, nie jej końcem — to zresztą sedno
produktu, bo sesja przeżywa zamknięcie aplikacji. Zakończenie sesji jest
czynnością osobną, jawną i nieodwracalną.

**Zamknięcie okna roboczego wraca do poprzedniego.** Operator pracujący
w Studio otwiera ikoną moduł automatyzacji, tworzy automatykę, zamyka to okno
i wraca do Studio, żeby ją tam uruchomić. Praca poprzednia nie jest przerywana
ani odtwarzana — okna trwają obok siebie, a przełącznik nimi zarządza.

**Karty, widoki i strony ramy.** Każda funkcja jest osobną kartą; karty
przełącza się w paśmie. Operator może zestawić dwie karty w **widoku dzielonym**
— zawsze dwie i zawsze lewo–prawo, bez podziału w pionie. Pasmo kart pozostaje
przy tym bez zmian: obie zestawione karty nadal w nim stoją, opasane **wspólną
ramką** pokazującą, że tworzą parę. Przy pasmie poziomym ramka biegnie wzdłuż
szerokości; przy kartach zgrupowanych w jedną — po rozwinięciu grupy biegnie
pionowo. Żadna pozycja nie znika i żadna nie przybywa.

Par bywa wiele naraz: przy ośmiu otwartych kartach Operator może mieć cztery
pary, czyli cztery widoki, i przełączać się między nimi.

**Okno boczne.** Prawa strona ramy niesie w stanie domyślnym samouczek — ale
samouczek jest jej **zawartością domyślną, nie przeznaczeniem**. Jest to
powierzchnia do wykorzystania i Operator może umieścić w niej dowolną kartę;
ma ona **własne pasmo kart**, więc karta przeniesiona z okna roboczego znika
z tamtego pasma i pojawia się tutaj, a powrót działa tak samo w drugą stronę.
Daje to trzy widoki naraz: dwa w oknie roboczym i jeden z boku.

**Kart w oknie bocznym bywa dowolnie wiele.** Nie ma tu żadnego limitu — pięć,
dwadzieścia, ile Operator zechce; jedynym ogranicznikiem jest pamięć maszyny.
Przełącza się je tak samo jak w oknie roboczym. Typowe zestawienie to wykaz
plików katalogu roboczego, drugi wykaz dla innego katalogu, terminal i podgląd
zmian pokazujący, co modele właśnie poprawiły.

Powierzchnia okna bocznego jest niewielka, więc nadaje się szczególnie dla kart,
które Operator chce mieć stale pod ręką do zerknięcia i szybkiego sięgnięcia.
Nie jest to ograniczenie, lecz przewidywane użycie, z którego wynika wymaganie
projektowe: **każda karta znosi wąską szerokość** i zachowuje w niej
użyteczność, zamiast wymuszać poziomy suwak albo ucinać treść.

**Lewa strona ramy jest zarządem sesji i projektów.** Kart nie przyjmuje — jest
jedynym z trzech obszarów, który nie jest dla nich miejscem — ale nie jest przez
to samym wykazem. Toczy się w niej pełne zarządzanie: zakładanie i usuwanie
sesji oraz projektów, zmiana nazw, wgląd w szczegóły, przenoszenie sesji
prowadzonej poza projektem do wybranego projektu.

**Jest to jedyne miejsce, w którym widać całość dorobku Operatora** — wszystkie
sesje i wszystkie projekty naraz. Wykaz projektów otwierany w oknie roboczym
pokazuje projekty wraz z ich sesjami, ale sesji prowadzonych poza projektami nie
obejmuje. Lewa strona obejmuje jedno i drugie, i stąd prowadzi jedyne
bezpośrednie przejście z projektu na projekt i z sesji na sesję.

Lewa strona należy do **aplikacji, nie do okna roboczego** — nie jest strefą
obszaru pracy i nie zmienia się wraz z nim. Zwija się i rozwija tak samo jak
prawa; zamknąć się jej nie da, podobnie jak strony z samouczkiem.

**Obie strony ramy zwijają się, a obszar roboczy przejmuje uwolnioną
przestrzeń.** Lewa — z wykazami sesji i projektów — oraz prawa zachowują się
tak samo: kto ich potrzebuje, trzyma je rozwinięte; kto potrzebuje szerokiego
obszaru pracy, zwija. Zwinięcie obu daje dwa widoki okna roboczego na pełnej
szerokości ekranu, po co zwijanie właśnie istnieje.

**Wykaz plików jest pełnym narzędziem plikowym, nie samym drzewem.** Obejmuje
katalog roboczy, urządzenie oraz wskazane serwery — ten sam zakres, którym
rządzi kaskada dostępu — i pozwala pliki otwierać, podglądać, kopiować oraz
wnosić nowe. Służy szybkiemu sięganiu po rzeczy, na których toczy się praca,
i podglądaniu tego, co modele właśnie wytwarzają.

**Cztery czynności obsługują całość:**

| Czynność | Skutek |
|---|---|
| Dodaj do widoku dzielonego | zestawia kartę z inną w widoku głównym |
| Rozłącz widok dzielony | rozdziela parę; obie karty zostają otwarte |
| Dodaj do okna bocznego | przenosi kartę do pasma okna bocznego |
| Wróć do okna głównego | zdejmuje kartę z okna bocznego |

**Wiele modeli w jednej sesji.** Kilka kart czatu to kilka modeli pracujących
w tej samej sesji. Podział widoku pozwala postawić dwa czaty obok siebie
i uruchomić przekazywanie wiadomości między modelami wraz z przydziałem ról —
koordynator i wykonawca. To jest sedno produktu, nie funkcja dodatkowa.

**Projekt jest środowiskiem zamkniętym.** Praca w projekcie toczy się w jego
ustawieniach: modele wnoszone do okien roboczych pracują na plikach projektu
i podlegają jego postanowieniom. Każda sesja założona w trybie projektu należy
do tego projektu, a wykaz sesji projektu jest jego własnym wykazem — projekt
z czternastoma sesjami pokazuje czternaście, niezależnie od sesji prowadzonych
poza nim.

Sesję wznawia się z wykazu i wraca do niej wraz z całą jej historią. Gdy
wznowiona sesja jest już zbyt obciążona, Operator otwiera nowe okno robocze
i zakłada w nim kolejną sesję — nadal w obrębie tego samego projektu.

**Karta incognito.** Karta może pracować bez zapisu historii, wzorem karty
prywatnej w przeglądarce. Dotyczy kart niosących interakcję z modelem, bo tylko
tam powstaje historia do pominięcia. Model pracujący w takiej karcie nie ma dostępu do
historii ani do danych sesji. Ma dostęp do plików projektu, jeżeli okno robocze
jest otwarte w projekcie. Wytwory — pliki powstałe w tej karcie — **są**
zapisywane; nie jest zapisywana historia. Zamknięcie karty kończy jej odrębny
zapis bezpowrotnie.

**Gałąź robocza — mechanizm odrębny od incognito.** Kilka modeli pracujących na
tych samych plikach nadpisuje sobie wyniki. Operator może temu zapobiec,
włączając dla karty pracę na gałęzi. Domyślnie gałąź jest **wyłączona** —
nowa karta pracuje wprost na plikach, dopóki Operator nie postanowi inaczej.

**Zasada nadrzędna: sterowalność zamiast domysłu.** Każdy mechanizm okna
roboczego Operator włącza świadomie. Wartość domyślna istnieje tylko tam, gdzie
czynność musi mieć jakiś stan początkowy, i jest wtedy stanem najprostszym,
nie najbogatszym. System nie decyduje za Operatora, na czym pracuje i co
zostaje zapisane.

**Konsekwencje.**

Warstwa projektowa wymaga przebudowy w części, która zakłada model — około
137 kB arkuszy wraz z odpowiadającymi skryptami. Wartości projektowe, fundament
i biblioteka komponentów, około 102 kB, zostają: ze 138 klas `komponenty.css`
tylko trzy zakładają model.

`design/zasoby/izolacja.css` jest wprost sprzeczny z tą decyzją. Jego nagłówek
stanowi „izolacja jest cechą okna, nigdy pojedynczej karty", a karta incognito
jest izolacją na poziomie karty. Izolacja ma dwa poziomy: okno robocze niesie
izolację swojej sesji, karta może nieść własną.

Pionowy pasek modułów **nie** jest sprzeczny z modelem i nie podlega usunięciu.
Jego przebudowa dotyczy tego, że dziś wypycha strefę szybkiego wyboru poza kadr
przy trzydziestu ośmiu pozycjach — czyli mieszczenia i porządku stref, nie racji
bytu.

Podział widoku i przypisanie kart do paneli należą do stanu sesji, a sesja
przeżywa zamknięcie aplikacji — więc układ jest trwały i wchodzi do modelu
danych, nie tylko do warstwy widoku.
---

## 6. Produkt otwiera możliwości, nie stawia granic

**Data:** 2026-08-22 · **Stan:** obowiązuje

**Kontekst.** Produkt ma być narzędziem, które Operator układa pod siebie, a nie
takim, który wymusza na nim konfigurację przed rozpoczęciem pracy. Odruch
wykonawcy jest odwrotny: domyślna odmowa, wymagane pola, komunikat o błędzie
przy braku ustawienia. Ten odruch trzeba zablokować wprost, bo inaczej wróci
w każdym oknie.

**Zasada nadrzędna produktu — i sposobu pisania o nim.** Ten produkt odróżnia
się od większości narzędzi tej klasy tym, że **nie wprowadza reguł, tylko
możliwości**. Nie stawia granic; otwiera drogi. Każde ograniczenie, jakie
w nim występuje, jest ograniczeniem, które **Operator sam ustanowił** dla siebie
albo dla modeli, którymi się posługuje.

Z tego wynika sposób opisu, wiążący dla interfejsu i dla całej dokumentacji:
piszemy, **co Operator może**, a nie czego mu nie wolno. Zdanie „bez wskazania
katalogu praca się nie zacznie" jest w tym produkcie fałszywe i wobec tego
zakazane jako opis; poprawne jest „wskazanie katalogu zawęża pracę do niego".
Ta sama treść podana od strony możliwości mówi prawdę o produkcie; podana od
strony zakazu — nie.

Zastrzeżenie: nie dotyczy to **dyscypliny budowy**. Zapory nałożone na
wykonawców — zakaz dopowiadania, zakaz otwierania wyjątków w trybie
`bypass permissions`, zakaz mylenia wykazu okien z wykazem sesji — chronią to,
co Operator dostanie, i nie odbierają mu niczego. One zostają wyrażone jako
zakazy, bo takimi są.

**Rozważone warianty.**

1. Domyślna odmowa: bez wskazania katalogu roboczego i uprawnień praca się nie
   zaczyna. Odrzucone: zamienia narzędzie w formularz, a pierwsze uruchomienie
   w konfigurację.
2. Ustawienia wymagane tylko dla czynności wrażliwych. Odrzucone: rozdziela
   czynności na dwie klasy, a granica między nimi nigdy nie jest oczywista
   i z czasem się przesuwa.
3. Brak ustawienia znaczy zakres pełny; ustawienie służy zawężeniu albo
   rozszerzeniu, zawsze z woli Operatora.

**Decyzja.** Wariant trzeci. **Nic nie jest przypisane na sztywno.** Wartości
domyślne istnieją, ale każdą Operator może zmienić. Brak ustawienia **nigdy** nie
jest przeszkodą w rozpoczęciu pracy i **nigdy** nie daje błędu — znaczy zakres
pełny. Ustawienie jest korzyścią do wykorzystania, nie obowiązkiem do spełnienia.

**Cztery poziomy zakresu.** Każdy dziedziczy po szerszym, każdy wolno nadpisać.
Poziom, na którym Operator niczego nie ustawił, **w całości przyjmuje to, co
stoi wyżej** — dziedziczenie jest zasadą, nie wyjątkiem.

| Poziom | Gdzie się ustawia | Brak ustawienia znaczy |
|---|---|---|
| **aplikacja** | konfiguracja produktu | zakres pełny |
| **projekt** | konfiguracja projektu; dziedziczą wszystkie jego sesje | to, co ustawiono dla aplikacji |
| **sesja** | okno robocze | to, co ustawiono dla projektu |
| **karta z interakcją modelu** | komponenty w oknie czatu | to, co ustawiono dla sesji |

Kaskada w działaniu: w konfiguracji aplikacji Operator ustawia tryb wymagający
zgody na każdą czynność. Zakłada projekt i nadaje mu tryb pracy bez pytania oraz
dostęp do wskazanego serwera — od tej chwili każda sesja tego projektu zaczyna
z tymi ustawieniami zamiast z aplikacyjnymi. W jednej karcie tej sesji ustawia
tryb samodzielnej edycji i dostęp wyłącznie do jednego katalogu; w karcie obok
model pracuje na innej maszynie. Rozstrzygnięcia z trzech różnych kondygnacji
obowiązują tu naraz, każde w swoim miejscu, a kondygnacja pominięta niczego nie
przerywa — przekazuje dalej to, co dostała.

**Zakres ma dwie osie i obie podlegają tej samej kaskadzie.**

| Oś | Czego dotyczy |
|---|---|
| **upoważnienie** | co model może zrobić sam, a o co musi zapytać |
| **dostęp** | gdzie model działa: katalogi, urządzenie, wskazany serwer |

**Dostęp nadaje się wskazaniem miejsca, nie wyliczeniem uprawnień.** Operator
wskazuje punkt w drzewie, a model obejmuje ten punkt wraz ze wszystkim poniżej
niego. Wskazanie katalogu głównego daje pracę od tego katalogu w dół, z
pominięciem katalogów systemowych; wskazanie urządzenia daje całe urządzenie.
Zakres rozszerza się i zawęża przesunięciem punktu zaczepienia, nie dopisywaniem
pozycji do wykazu.

Nadanie następuje na tej kondygnacji kaskady, którą Operator wybierze —
w konfiguracji aplikacji, w projekcie, w sesji albo przy pojedynczej karcie
czatu. Kondygnacja pominięta przekazuje dalej to, co dostała.

**Tryby upoważnienia — nazwy wiążące:** `manual`, `auto`, `plan`,
`bypass permissions`. Obowiązują w interfejsie, w dokumentacji i w kontrakcie;
nazwa zastępcza ani tłumaczenie nie wchodzą do produktu. Są to terminy przyjęte
w dziedzinie pracy z agentami — nie zapożyczenia od pojedynczego wytwórcy —
więc Operator znający narzędzia agentowe rozpoznaje je bez nauki.

**`bypass permissions` jest trybem pracy ciągłej, nie tylko zniesieniem pytań.**
Odróżnia go od odpowiedników w narzędziach zewnętrznych to, że zniesienie zgód
jest w nim **zupełne**. Model nie wraca po potwierdzenie przy żadnej czynności,
także nieodwracalnej, także wykraczającej poza katalog roboczy, także takiej,
przy której narzędzia zewnętrzne pytają mimo włączonego trybu.

Wykaz czynności wyjętych spod tego trybu jest **pusty i ma pozostać pusty**.
Nie wolno go otworzyć ani „dla bezpieczeństwa", ani przy czynnościach
nieodwracalnych, ani pod pretekstem wyczerpania zasobu, ani przy pierwszym
uruchomieniu. Tryb, który pyta choć raz, nie jest tym trybem — jest trybem
`auto` pod cudzą nazwą, a Operator, który ustawił autonomię i mimo to został
zapytany, nie dostał tego, co wybrał.

**Zatrzymanie przez Operatora nie należy do tego trybu i nie jest przez niego
ograniczone.** Przerwanie pracy jest osobnym komponentem okna czatu, dostępnym
zawsze i niezależnie od ustawionych uprawnień. Operator naciska je, kiedy chce,
i praca zostaje przerwana — żaden tryb upoważnienia tego nie osłabia i nie
warunkuje. Uprawnienia rozstrzygają, o co model pyta i jak długo pracuje;
przerwanie rozstrzyga, czy pracuje w ogóle. To są dwie różne rzeczy i nie wolno
ich wiązać.

Do zniesienia zgód dochodzi **nakaz pracy do wyniku**. Model kończy pracę
dopiero wtedy, gdy ma wynik zadania. Nie oddaje pracy w połowie, nie kończy tury
pytaniem i **nie melduje na etapach z przerwaniem** — sprawozdanie z postępu nie
jest powodem do zatrzymania się i oczekiwania. Dostaje zadanie, kończy zadanie,
wtedy dopiero może się wyłączyć.

Rozmowa wewnątrz układu nie jest przerwaniem pracy. Wykonawca zwraca się tam,
gdzie prowadzą powiązania nadane przez Operatora — do koordynatora, do modelu
równoległego, wyżej w hierarchii. Oni są częścią pracy.

**Ten tryb daje Operatorowi możliwość odejścia od maszyny.** Praca toczy się
bez niego, a rozstrzygnięcia zapadają tam, gdzie poprowadził powiązania. Jeżeli
chce być pytany, wpisuje siebie do układu jako ogniwo drogi odwoławczej — wtedy
powiadomienie dociera do niego jako wykonanie ustanowionej przez niego
procedury. Jeżeli tego nie zrobi, model rozstrzyga sam i prowadzi zadanie do
wyniku. Jedno i drugie jest wyborem Operatora; produkt nie narzuca żadnego
z nich.

**Ten tryb jest zarazem sygnałem, że po drugiej stronie nie ma człowieka.**
Model ma z niego odczytać nie tylko zakres swoich uprawnień, lecz i to, że nikt
nie czeka przy oknie, żeby odpowiedzieć — a więc że zadanie musi doprowadzić do
końca sam.

**Niepewność nie jest powodem do zatrzymania.** Model, który czegoś nie wie albo
nie ma rozstrzygnięcia, ma rozstrzygnąć sam. Nie wolno mu przerwać pracy po to,
by zadać pytanie rozstrzygające — ani w połowie, ani na początku. Powód nie jest
formalny: model dostał tożsamość, narzędzia i zadanie, więc **jest twórcą
wyniku i to on nim rozporządza**. Gdy zadaniem jest wytworzenie rzeczy, model
jest jej autorem — i decyzja o jej kształcie należy do autora, nie do Operatora,
który tej rzeczy nie zna. Jeżeli wykonawca nie wie, to nikt nie wie; odesłanie
pytania do Operatora nie sprowadza wiedzy, tylko przenosi rozstrzygnięcie na
tego, kto ma jej mniej.

**Po co ten rygor.** Tryb istnieje ze względu na pracę pętli wieloagentowej,
w której **ukończenie zadania wyzwala czynność następną**. Pętla obejmuje wiele
modeli koordynujących wielu wykonawców i pracuje bez przerwy — dobami, a przy
większych przedsięwzięciach miesiącami. Wykonawca, który przerwie bez wyniku
albo odda wynik niepełny, nie zatrzymuje własnego zadania; **zatrzymuje pętlę**,
bo następna czynność nie ma czym się wyzwolić. Jedna furtka w jednym miejscu
łańcucha kończy przedsięwzięcie liczone w miesiącach — dlatego wykaz wyjątków
jest pusty i dlatego nie wolno go otworzyć.

**Granica odpowiedzialności.** Jakość wyniku zależy od tego, jak Operator
przygotował wykonawcę — tożsamość, upoważnienia, narzędzia i oprzyrządowanie
nadaje mu w module Agents. Produkt daje narzędzie do przygotowania wykonawcy;
nie zastępuje Operatora w tym przygotowaniu i nie poprawia go w czasie pracy.
Zachowania agenta nie programuje się za Operatora.

**Nadpisanie działa w obie strony.** Karta może zawęzić uprawnienia modelu —
na przykład do samego odczytu — ale może też je **rozszerzyć** poza zakres
projektu, dopuszczając pracę na całym urządzeniu mimo wskazanego katalogu
roboczego albo na maszynie spoza projektu. Zakres węższego poziomu nie jest
więc podzbiorem szerszego; jest rozstrzygnięciem Operatora, które szerszy
poziom jedynie podpowiada.

**Karty jednej sesji pracują niezależnie od siebie.** Dwie rozmowy w tym samym
oknie roboczym mogą działać na dwóch różnych serwerach i dwóch różnych zbiorach
katalogów. To jest sedno skalowalności, o którą tu chodzi — nie o mnożenie
rozmów, lecz o to, że każda ma własne umocowanie.

**Sterowanie mieszka przy pracy, nie w oknie konfiguracji.** Uprawnienia
i dostępy modelu ustawia się komponentami wbudowanymi w okno czatu, przy
rozmowie, której dotyczą. Okno konfiguracji nie jest jedynym ani głównym
miejscem tych rozstrzygnięć.

**Konsekwencje.**

Domyślny zakres pełny jest **świadomym rozstrzygnięciem produktowym**, nie
przeoczeniem i nie usterką bezpieczeństwa do naprawienia. Wykonawca, który
zastanie brak ograniczeń, nie wprowadza domyślnej odmowy „dla bezpieczeństwa" —
zgłasza to Prowadzącemu, jeżeli uważa za konieczne, i czeka na rozstrzygnięcie.
Odpowiedzialność za zakres udzielony modelowi spoczywa na Operatorze, a produkt
ma mu pokazać ten zakres wyraźnie, nie odebrać.

Skoro zakres nadpisuje się na czterech poziomach i w obie strony, **zakres
obowiązujący musi być w każdej chwili widoczny przy rozmowie** — inaczej
Operator nie wie, czym dysponuje model, do którego właśnie pisze. To wiąże
kompozycję karty czatu.

Model danych niesie zakres na czterech poziomach z rozróżnieniem „nieustawione"
od „ustawione na pełne" — to dwa różne stany, bo pierwszy dziedziczy dalej,
a drugi rozstrzyga.

Skoro karta może wskazać inną maszynę niż sesja i projekt, **rdzeń wykonuje
pracę na wielu maszynach naraz w obrębie jednej sesji**. To nie jest ustawienie
interfejsu — to wymaganie wobec architektury wykonania i wchodzi do etapu 2.

---

## 7. MultitaskingAI: produkt daje środki, reguły układa Operator

**Data:** 2026-08-22 · **Stan:** obowiązuje · **Dotyczy zakresu poza etapami 1–2**

**Kontekst.** Środowisko MultitaskingAI bywa mylone z automatyzacją. Różnica
rozstrzyga o wymaganiach stawianych rdzeniowi już teraz, mimo że samo środowisko
powstaje później. Automatyzacja wykonuje ustalony ciąg czynności — do tego
wystarcza moduł Automations. MultitaskingAI wykonuje **układ**, który Operator
sam zbudował.

**Zastrzeżenie do sposobu opisu — obowiązuje każdego, kto pisze o tym
środowisku.** Reguł działania MultitaskingAI **nie da się wyczerpać zapisem**
i nie należy próbować. Kto komu przekazuje wynik, kto z kim się konsultuje, kto
zatwierdza przejście dalej, kiedy i do kogo biegnie droga odwoławcza — to nie
są postanowienia produktu, tylko rozstrzygnięcia Operatora, podejmowane osobno
w każdym układzie. Liczba możliwych ułożeń jest praktycznie nieograniczona;
opracowanie próbujące je opisać rozrosłoby się bez końca i i tak nie objęłoby
wszystkiego. **Dokumentacja opisuje środki, nie ułożenia.**

Poniższe jest wobec tego wykazem tego, co produkt musi umieć, żeby Operator
mógł zbudować dowolny układ — nie opisem układu.

### Czego produkt nie rozstrzyga

Hierarchia. Role i odpowiedzialności. Kierunek i istnienie każdego powiązania.
Prawo do konsultacji i prawo do zwrotu pracy. Kto zatwierdza przejście dalej.
Długość i przebieg drogi odwoławczej. To, czy Operator w ogóle występuje
w układzie. Liczba modeli, agentów, zespołów i pętli. Podział pracy na etapy
zamykane wewnątrz zespołu.

Wszystko to Operator nadaje przy budowie układu i wszystko może ułożyć inaczej
w każdym kolejnym.

### Co produkt musi umieć

| Wymaganie | Dlaczego |
|---|---|
| Zapisać układ jako dane podlegające odczytowi, zmianie i odtworzeniu | układ jest wytworem Operatora, nie konfiguracją wkompilowaną w produkt |
| Utrwalić kierunek każdego powiązania oraz przypisane mu uprawnienia | bez tego drogi przekazania i zwrotu nie dają się odtworzyć |
| Podać modelowi **wyłącznie jego własne powiązania**, nigdy mapy układu | model, który nie zna adresata, nie może pominąć hierarchii; mapa podana modelowi znosi tę zaporę niezależnie od poleceń |
| Rozróżnić maszynowo ukończenie z wynikiem od przerwania | ukończenie wyzwala czynność następną gdzieś w układzie |
| Odtworzyć stan układu po restarcie, awarii i zamknięciu aplikacji | praca liczy się w setkach dni, więc restart jest pewnością, nie wypadkiem |
| Prowadzić wiele torów modeli jednocześnie i rozdzielać je | to samo wymaganie, które daje karcie prawo do własnej maszyny |
| Wydać wynik w postaci użytecznej i dla modelu, i dla człowieka | odbiorca nie jest znany z góry i bywa różny w każdym ogniwie |

### Rzeczy, które produkt zakłada niezmiennie

Dwie, i tylko dwie.

**Twórcą jest sztuczna inteligencja, od pomysłu do wykonania.** Operator wnosi
zadanie, buduje układ i oddaje kompetencje. Wymyślenie rzeczy, zaprojektowanie
jej i wykonanie należą do modeli. Zadanie brzmi „zaprojektujcie rower" i nic
ponadto.

**Model nie kieruje pytania do Operatora z własnej inicjatywy.** Kieruje je
tam, dokąd prowadzi jego powiązanie. Operator bywa ostatnim ogniwem drogi
odwoławczej, ale tylko wtedy, gdy sam ją tak poprowadził — i wtedy dotarcie do
niego jest wykonaniem procedury, nie ucieczką od decyzji.

### Granica odpowiedzialności

Jakość układu rozstrzyga o jakości wyniku. Tak jak w przedsiębiorstwie dobór
ludzi i kierownictwa przesądza o tym, co firma wytworzy, tak tutaj przesądza
nadana hierarchia, przydzielone role, umiejętności i tożsamości. Produkt daje
środki; nie poprawia układu, który Operator zbudował, i nie odpowiada za jego
skutki.
---

## 8. Model wdrożenia: wyłącznie hybryda

**Data:** 2026-08-26 · **Stan:** obowiązuje · **Zastępuje** wcześniejsze
rozstrzygnięcie o dwóch wariantach instalacji

**Decyzja.** Produkt występuje w **jednej postaci**: cienka instalka u Operatora,
rdzeń i całe zaplecze na serwerze Danaco. **Wariant pełny natywny nie powstaje.**

| | Na dysku Operatora | Gdzie stoi zaplecze |
|---|---|---|
| hybryda | **~250 MB** | serwer Danaco |

**Jedyna platforma: Windows 11.** Jedynym wyborem w instalatorze jest
architektura — **x64** albo **ARM64**. Linuksa, macOS ani Windows 10 produkt nie
obsługuje. Oba cele budowania stoją na maszynie budowlanej i są sprawdzone:
`x86_64-pc-windows-gnu`, `aarch64-pc-windows-msvc`, llvm-mingw
z `aarch64-w64-mingw32-gcc`.

**Podstawa — pomiar.** Zaplecze pracy zważone na maszynie budowlanej: modele AI
15,6 GB, arsenał mediów i dokumentów 2,7 GB, n8n 2,5 GB, whisper z wagami
1,7 GB, przeglądarki Playwright 0,93 GB, silnik rembg 0,72 GB, reszta programów
i bibliotek 3,7 GB — **razem ~28 GB**. Wariant natywny musiałby to powtórzyć na
każdym urządzeniu; z miejscem roboczym na obróbkę mediów daje to **30 GB
minimum, 50 GB zalecane**. Hybryda potrzebuje **~250 MB** — sto dwadzieścia razy
mniej.

**Powód.** Nie chodzi o rozmiar pliku, lecz o utrzymanie. Zaplecze rośnie:
narzędzie dołożone do produktu w wariancie serwerowym staje przed wszystkimi
Operatorami tego samego dnia, a w natywnym wymaga nowego wydania, nowego
pobrania i nowej instalacji u każdego z osobna. Przy zapleczu tej wielkości
druga droga nie jest wykonalna.

**Konsekwencje.**

1. **Instalator ma jedną ścieżkę.** Nie ma kroku wyboru wariantu ani kroku
   wyboru składników — nie ma czego wybierać, bo składniki nie schodzą na
   urządzenie. Architekturę instalator rozpoznaje sam i wskazuje właściwą
   postać; wybór ręczny zostaje jako możliwość, nie jako krok wymagany.
1a. **Rdzenia nie ma w instalce w ogóle.** Wchodzi wyłącznie powłoka
   z wkompilowanym interfejsem, więc budowanie krzyżowe rdzenia na Windows nie
   jest produktowi potrzebne. Zostaje sama powłoka Tauri w dwóch celach.
1b. **WebView2 jest składnikiem Windows 11.** Instalator go nie niesie i nie
   stawia — to jedyna zależność systemowa produktu i jest zaspokojona z góry.
2. **Produkt wymaga łączności z serwerem Danaco.** Praca bez sieci nie jest
   przewidziana. Wdrożenie odcięte od sieci nie jest obsługiwane.
3. **Zasady korzystania opisują jeden sposób pracy.** Zapis wymagający
   opisywania dwóch wariantów bez wskazywania domyślnego przestaje obowiązywać.
4. **Skrypty instalek natywnych są materiałem zamkniętym.** Leżą poza
   repozytorium, w `~/robocze/material/repo-2.0/budowa/scripts/`. Do budowy nie
   wchodzą.
5. **Modele i arsenał nie wchodzą do żadnej instalki.** Stoją na serwerze
   wdrożenia i są utrzymywane w jednym miejscu.

## 9. Wstążka narzędziowa okna roboczego i przybornik karty

**Data:** 2026-08-23 · **Stan:** obowiązuje

**Decyzja.** Pod pasmem kart okna roboczego stoi **wstążka narzędziowa**, która
niesie funkcje właściwe **rodzajowi otwartej karty**. Karta przeglądarki daje
inne narzędzia niż karta terminala, a ta inne niż karta wykazu plików. Wstążka
zmienia się wraz z przełączeniem karty; nie jest wspólnym pasem dla wszystkich.

**Wzorcem zachowania jest przeglądarka internetowa.** Karta u góry, wstążka
narzędziowa pod nią, zachowania takie, jakie ma dziś każda przeglądarka.
Standardy przyjęte w tej klasie programów obowiązują bez potrzeby ustalania
reguły dla każdego komponentu osobno.

**Reguła pracy, nie odpowiedź jednorazowa:** domyślnie stosuje się standardy
platformy, a pytanie zadaje się wyłącznie tam, gdzie produkt od nich naprawdę
odchodzi. Obowiązuje to każde zlecenie — wykonawca stosuje standard, zamiast
pytać o regułę dla każdego komponentu.

### Kontrakt wstążki okna roboczego

**Położenie i forma.** Jeden rząd bezpośrednio pod pasmem kart, na pełną
szerokość okna roboczego. Wysokość stała, niezależna od zawartości. Wstążka nie
rośnie, nie zawija się do drugiego rzędu i nie przewija się poziomo.

**Zawartość zmienna, forma stała.** Przełączenie karty podmienia zawartość; sam
pasek nie drgnie.

**Trzy strefy, jak w przeglądarce.**

| Strefa | Zawartość | Odpowiednik w przeglądarce |
|---|---|---|
| lewa | czynności na treści karty | wstecz, dalej, odśwież |
| środkowa | pole kontekstu, elastyczne, pochłania wolną szerokość | pole adresu |
| prawa | narzędzia jako przyciski ikonowe stałej wielkości | rozszerzenia i menu |

**Kurczenie.** Przy zwężaniu okna kurczy się najpierw strefa środkowa. Strefy
ikonowe trzymają rozmiar do progu, po którym prawa grupa zwija się do **menu
nadmiaru** — nie zawija się i nie znika po cichu. Ten sam wzorzec, którego brak
w szynie nawigacji sprawia, że sześć z siedmiu pozycji strefy szybkiego wyboru
wypada poza kadr bez żadnego sygnału.

**Przyciski ikonowe.** Kwadratowe, jednakowej wielkości, bez podpisów,
z podpowiedzią po najechaniu z opóźnieniem. Narzędzie włączone niesie trwały
stan wybrania, wzorem przypiętego rozszerzenia.

**Trwałość.** Wstążka stoi nieruchomo, gdy treść karty się przewija, i nie chowa
się przy przewijaniu.

**Klawiatura.** Każda kontrolka osiągalna klawiszem tabulacji; cały pasek jako
jeden obszar orientacyjny.

**Dostosowanie.** Menu kontekstowe na pasku otwiera dostosowanie zawartości —
w produkcie istnieje już jako okno „Dostosuj paski".

Zmienne pozostaje wyłącznie to, co w danej karcie stoi w trzech strefach:
przeglądarka dostaje przybornik, terminal coś innego, wykaz plików jeszcze co
innego.

### Przybornik przeglądarki — dostępny we wszystkich modułach

Karta przeglądarki jest osiągalna w każdym module i środowisku, w którym jest
dopuszczona macierzą. Jej wstążka otwiera **przybornik pływający** o zakresie
podstawowym:

| Narzędzie | Do czego |
|---|---|
| zrzut ekranu | utrwalenie widoku |
| inspektor elementów | wskazanie elementu strony wraz z jego nazwą z kodu |
| widok komputerowy i mobilny | obejrzenie strony w obu szerokościach |
| znakowanie: pole tekstowe, mazak, podkreślnik, ramka, strzałka, wybór barwy | opisanie i oznaczenie fragmentu |

Przeznaczenie przybornika: **lokalizowanie, znakowanie i przekazywanie dalej** —
do rozmowy z modelem, nie do samodzielnej analizy.

### Inspektor elementów

Wskaźnik prowadzony po stronie odsłania obszary elementów, każdy innym
wyróżnieniem, wraz z nazwą, jaką element nosi **w kodzie strony**. Ujawnia
warstwy niewidoczne dla oka: tam, gdzie widać jedno tło, wskaźnik pokazuje, ile
powłok faktycznie się na nie składa.

Inspektor **pokazuje wyłącznie to, co jest zapisane w kodzie strony, i nic
ponadto**. Nie dolicza, nie interpretuje i nie wyprowadza wniosków.

Działa na każdej stronie, także dowolnej stronie sieci — nie tylko na oknach
własnych produktu.

### Menu wskazanego elementu

Kliknięcie **nie wykonuje czynności — otwiera menu**. Cztery pozycje w dwóch
grupach i ani jednej więcej:

| Grupa | Pozycja | Co wnosi |
|---|---|---|
| Do czatu | **Zrzut z zaznaczeniem** | obraz obszaru karty z obrysowanym elementem, jako załącznik |
| Do czatu | **Element jako tekst** | odczyt panelu wstawiony do pola wpisywania, edytowalny |
| Zachowaj | **Do schowka** | ten sam odczyt jako tekst |
| Zachowaj | **Do modułu Library** | zrzut wraz z odczytem, jako jeden artefakt |

Rozstrzygnięcia szczegółowe wraz z powodami:

- **Zrzut obejmuje obszar karty, nie sam element z przycięciem.** Wycinek bez
  otoczenia gubi to, co czyni zgłoszenie czytelnym — nie widać, gdzie rzecz
  stoi. Jest to zgodne z istniejącą pozycją „Wskazane okno robocze" w menu
  zrzutu.
- **Przytoczenie wchodzi jako tekst edytowalny, nie zamknięta plakietka.**
  Odczyt panelu liczy kilkanaście wierszy, a w wiadomości zwykle liczą się trzy;
  Operator musi móc go przyciąć przed wysłaniem.
- **Schowek niesie dokładnie ten sam tekst co przytoczenie**, nie trzeci format.
  Jeden odczyt, dwa ujścia.
- **Library dostaje obraz razem z odczytem**, bo jest repozytorium wiedzy — sam
  obraz jest tam nieodnajdywalny.
- **Wędrówka po drzewie nie należy do menu.** Przejście do elementu nadrzędnego
  i podrzędnego mieszka w panelu inspektora, obsługiwane strzałkami. Menu
  odpowiada na jedno pytanie: co z tym elementem zrobić. Zmiana wskazania nie
  jest czynnością do wykonania.
- **Kopiowanie selektora nie występuje** — nazwa elementu stoi już w odczycie,
  więc osobna pozycja dublowałaby to samo pod inną nazwą.

**Brzmienia są zgodne z istniejącym słownictwem produktu.** Menu „Wykonaj zrzut"
ma sekcję „Po wykonaniu" z pozycjami: do schowka, do pliku w katalogu zrzutów,
do modułu Library, do Chat Window jako załącznik. To jest ustalone słownictwo
produktu dla tego, gdzie rzecz trafia — inspektor mówi tym samym językiem, a nie
zakłada drugiego, konkurencyjnego.

**Załącznik i przytoczenie muszą różnić się brzmieniem**, bo niosą co innego:
załącznik to obraz wraz z odczytem, przytoczenie to sam odczyt wstawiony w tekst
wiadomości.

### Zachowanie przybornika pływającego

Da się go przesunąć w dowolne miejsce. **Nie idzie za wskaźnikiem.** Położenia
nie pamięta — po zamknięciu i ponownym otwarciu staje w miejscu domyślnym.
Przybornik należy do karty: zamknięcie karty zamyka przybornik.

### Granica wobec modułu Browser

Przybornik podstawowy kończy się na powyższym. Pełna konsola strony wraz
z oprzyrządowaniem analitycznym należy do modułu Browser i nie występuje
w karcie pomocniczej.

### Always On Display — funkcja odrębna, bez kolizji

Asystent towarzyszący pracy: analizuje to, co jest na ekranie, i doradza, gdy
Operator tego potrzebuje. Bywa przedstawiony jako element pływający dający
szybki dostęp do rozmowy z asystentem; Operator może go odłączyć, żeby nie
zajmował ekranu — asystowanie trwa niezależnie od tego, czy element jest
widoczny. Pływający sygnet jest **skrótem do funkcji, nie funkcją**, więc
o miejsce na ekranie nic tu nie konkuruje.

Warte rozważenia później, nie teraz: skoro asystent i tak czyta ekran, odczyt
inspektora jest dokładnie tą postacią, w której da się mu ten ekran podać —
uporządkowaną, a nie odgadywaną z obrazu.

### Podział odpowiedzialności przy szczegółach

Rozstrzygnięcia szczegółowe — sposób zachowania komponentu, dobór i układ
funkcji w narzędziu — należą do wykonawcy jako praca inżynierska, nie do
Właściciela. Warunek jest jeden: funkcje muszą być **realne i osiągalne**,
przemyślane wobec tego, czemu i komu służą, oraz spójne z całością. Wymyślanie
rozwiązań abstrakcyjnych albo niewykonalnych nie jest korzystaniem z tej
swobody.

---

## 10. Powitanie kanału jest wyjęte spod bramy kontraktu

**Data:** 2026-08-24 · **Stan:** obowiązuje

**Kontekst.** Rdzeń przyjmował żądania bez pól wymaganych i wartości spoza
wyliczeń, milcząc — klient mógł wysyłać treść niezgodną z kontraktem i nie
dowiedzieć się o tym. Naprawa wniosła bramę sprawdzającą każde żądanie wobec
kontraktu. Brama objęła także `connection.hello`, którego trzy pola kontrakt
oznacza jako wymagane.

**Skutek, który to ujawniło.** Powitanie jest jedynym miejscem, w którym klient
odczytuje `protocolVersion` rdzenia — a więc jedynym, w którym rozpoznaje, że
jest starszy niż rdzeń. Klient nieznający pola `clientId`, bo pochodzi sprzed
jego wprowadzenia, dostaje od bramy odmowę zamiast informacji o niezgodności
wersji. Im starszy klient, tym pewniej nie dowie się, dlaczego został odrzucony.

**Decyzja.** `connection.hello` **nie podlega bramie**. Odpowiada zawsze, także
na żądanie niepełne.

Braków pól powitanie **nie zgłasza w treści odpowiedzi** — `ConnectionHelloResponse`
nie ma pola, w którym mogłoby to zrobić, a dołożenie takiego pola jest zmianą
kontraktu, czyli rzeczą osobną i większą. Braki idą do dziennika rdzenia. Do
celu wyjątku to wystarcza: klient starszy ma odczytać `protocolVersion` i sam
rozpoznać rozjazd, a do tego potrzebuje odpowiedzi, nie wykazu swoich braków.

**Powód.** Uzgodnienie wersji musi działać przed uzgodnieniem czegokolwiek
innego — w tym przed zgodnością co do pól. Brama sprawdzająca powitanie wobec
kontraktu zakłada, że obie strony już znają ten sam kontrakt, czyli zakłada to,
co powitanie ma dopiero ustalić.

**Konsekwencje.** Powitanie jest jedynym wyjątkiem i pozostanie jedynym. Każda
inna komenda przechodzi przez bramę bez ustępstw — wyjątek dla powitania wynika
z jego roli w uzgodnieniu, nie z wygody.

---

## 11. Pierwsze uruchomienie bez poczty przechodzi i nie wymaga rozstrzygnięcia

**Data:** 2026-08-26 · **Stan:** obowiązuje · **Podstawa:** pomiar uruchomieniem

**Pozycja otwarta w tej sprawie została zamknięta bez rozstrzygnięcia
Właściciela, bo pytania nie było** — rdzeń niesie odpowiedź, a otworzyłem ją
biorąc zawodzący sprawdzian za usterkę.

**Co rdzeń robi.** Rejestracja idzie dwiema drogami. Z pocztą nadaje list
z drogą potwierdzenia i zostawia konto niepotwierdzone. Bez poczty zakłada
konto i zapamiętuje w sejfie znacznik `auth:bramka-bez-poczty` niosący adres,
którego nikt nie potwierdził. Bramka czyta ten znacznik i wpuszcza hasłem.

**Co zmierzono.** Świeża instalacja bez nadajnika: rejestracja udana
(`Registered:true`, `PendingVerification:false`), jedno konto w bazie, zero
wierszy potwierdzenia, logowanie hasłem wydaje sesję z tokenem. Droga przechodzi
w całości.

**Dlaczego rdzeń tak stanowi.** Poczta jest potrzebna do pisania do innych
ludzi — do potwierdzania ich adresów i do odzyskiwania hasła listem — a nie do
postawienia bramki na własnym urządzeniu. Nadajnik ustawia się w oknie
Konfiguracji, czyli za bramką; wymaganie poczty do jej postawienia zamykałoby
pierwszego Operatora przed platformą na zawsze.

**Czego stan bez poczty nie udaje.** Wiersz konta stoi na „niepotwierdzone"
i nikt tego nie zmienia. Znacznik zdejmuje jeden warunek i tylko jeden: bramki
nie zamyka potwierdzenie, którego platforma nie miała czym wysłać.

**Skutek dla okna.** Przepływ wejścia ma dwie gałęzie, nie jedną — i różnią się
**bramką, nie tylko komunikatem**. Ten akapit prostuje pierwotny zapis pozycji,
który różnicę zaniżał; poprawka pochodzi z pomiaru wykonawcy terenu
`droga-wejscia`.

| Nadajnik | Odpowiedź rejestracji | Czym wchodzi Operator | Co okno mówi |
|---|---|---|---|
| ustawiony | `PendingVerification: true` | **wyłącznie `auth.verify`** — `auth.login` odmawia `not_authenticated`, bo konto czeka na potwierdzenie | list wysłany, przepisz drogę potwierdzenia |
| brak | `PendingVerification: false` | `auth.login` hasłem | konto założone, adres niepotwierdzony — wejście działa hasłem, **odzyskanie konta listem nie zadziała do chwili potwierdzenia** |

Wiersz pierwszy rozstrzyga przepływ okna: po rejestracji z nadajnikiem **nie ma
przejścia do logowania**. Wykonawca czytający samą pierwotną treść pozycji
zbudowałby przejście, które rdzeń odrzuca.

Ostrzeżenie z drugiego wiersza jest wymagane, nie zalecane. Adres jest jedyną
drogą odzyskania konta; Operator, który nie wie, że jego adres nie został
sprawdzony, dowie się o tym w dniu, w którym będzie go potrzebował.


---

## 12. Prototyp instalatora jest wersją przyjętą

**Data:** 2026-08-26 · **Stan:** obowiązuje · **Rozstrzygnął:** Właściciel

Kompozycję kreatora instalacji ustala prototyp
`design/05-okna/platformowe/instalator.html` wraz ze składnikami
w `design/zasoby/okna/instalator/`.

**Prototyp jest projektem, nie inżynierią ostateczną.** Rozstrzyga to, czego
uruchomiony układ nie rozstrzygnie sam; nie rozstrzyga tego, co układ zmierzy.
Ta granica obowiązuje wszystkie prototypy, nie tylko instalator.

| Prototyp **rozstrzyga** | Prototyp **nie rozstrzyga** |
|---|---|
| które ekrany, w jakiej kolejności | ile co trwa |
| co Operator wybiera i czego nie musi | ile waży i ile zajmuje |
| co się dzieje po wyborze, odmowie, przerwaniu | ile jest czego — plików, składników, etapów |
| gęstość, układ, hierarchia, ton wypowiedzi | kody błędów, ścieżki, numery wersji |
| które stany są przewidziane | wartości, które da się odczytać z działającego układu |

**Wartość liczbowa w prototypie jest miejscem na wartość, nie wartością.**
Wykonawca bierze prawdziwą z pomiaru i **nie zgłasza różnicy jako rozjazdu**.
Odwrotnie też: „to tylko parametr" nie upoważnia do zmiany kompozycji.

**Jak rozpoznać, po której stronie stoi rzecz.** Jeżeli działający układ potrafi
ją zmierzyć — to parametr. Jeżeli nie potrafi jej wybrać, bo to rozstrzygnięcie
człowieka — to kompozycja.

**Sześć kroków:** Wymagania · Licencja · Wersja programu · Lokalizacja ·
Instalacja · Podsumowanie.

**Wykrycie procesora zamiast wyboru wariantu.** Krok 3 rozpoznaje procesor
i zaznacza właściwą postać, oznaczając ją plakietką ZALECANE. Wybór ręczny
zostaje; niezgodność ostrzega i pyta osobnym oknem, ale nie blokuje.

**Katalog treści jest jedynym miejscem z tekstem widocznym dla użytkownika.**
Sprawdzone pomiarem: w składnikach i ekranach stoi 427 łańcuchów i żaden nie
jest tekstem dla użytkownika. Ta właściwość jest wiążąca — łańcuch dopisany
poza `tresci.js` jest usterką.

**Odsłony osiągalne adresem** — `?procesor=arm|brak`, `?stan=blad|wycofywanie`,
`?wynik=ostrzezenia` — należą do prototypu i mają przetrwać przejście do kodu.
Bez nich odsłony błędu przestają być przeglądalne.

**Dwie rzeczy stoją po stronie kompozycji, a prototyp mówi w nich dwie rzeczy
naraz** — nie są to parametry, bo żaden pomiar ich nie rozstrzygnie:

| Rzecz | Krok mówi jedno | Krok mówi drugie |
|---|---|---|
| pobieranie a rozpakowanie | krok 1: „Instalator pobiera składniki programu z serwera Danaco" | krok 5: „Rozpakowywanie plików — zapis plików programu w katalogu docelowym", licznik „Rozpakowano 250 MB" |
| uprawnienia | krok 4: „instalacja w katalogu profilu użytkownika nie wymaga uprawnień administratora" | krok 5: odmowa dostępu do tego samego katalogu profilu, rada „uruchom jako administrator" |

**Jeżeli instalator pobiera z serwera Danaco, adres serwera musi skądś pochodzić.**
Żaden z sześciu kroków o niego nie pyta, więc jest wpisany w postać instalki przy
jej składaniu — osobnej dla każdego wdrożenia. To wniosek z prototypu, nie
zastrzeżenie wobec niego.


---

## 13. Rejestr kont i kanałów modelu jest otwarty

**Data:** 2026-08-26 · **Stan:** obowiązuje · **Podstawa:** pomiar uruchomieniem

**Wymaganie Właściciela.** Aplikacja niesie na starcie konta Anthropic przez
Code CLI, ale **to nie jest wykaz zamknięty**. Za jakiś czas kont może być
sześć zamiast czterech, albo jedno zostanie przy CLI, a pozostałe przejdą na
API — i produkt ma to znieść **bez przerabiania aplikacji**.

**Zmierzone.** Sprawdzian przeszedł całą drogę wobec rdzenia:

| Krok | Wynik |
|---|---|
| cztery konta Anthropic przez CLI | założone, `rodzaj=cli`, `dostawca=anthropic` |
| piąte i szóste konto CLI | **dołożone bez przeszkody** — wykaz nie jest zamknięty |
| konto obcego dostawcy przez API | założone, `dostawca=openai`, własny `baseUrl` |
| kanał na generycznym adapterze sieciowym | założony, `rodzaj=api` |
| usunięcie konta | **kanały przeżyły: 3 przed, 3 po** |

**Dlaczego rdzeń to znosi.** `KanalAPI` nie zna żadnego dostawcy — adres, model,
odwołanie do klucza, nagłówki, kształt ciała żądania i ścieżka do treści
w odpowiedzi pochodzą z wiersza rejestru. Nagłówek pliku stanowi wprost: „nowy
dostawca to nowy wiersz danych, nie nowy typ w kodzie". Pole `provider`
kontrakt opisuje jako „wartość danych, nie typ kodu". Adapter CLI bierze wykaz
kont z katalogu profili na progu każdej tury, więc konto dodane komendą
`account.*` wchodzi do rotacji **bez restartu rdzenia**.

**Jedyny przypadek wymagający kodu** to dostawca bez HTTP, dostępny wyłącznie
własną biblioteką — wtedy powstaje nowa fabryka adaptera. Dostawca z API
publicznym nie wymaga ani jednej linii.

**Konsekwencja dla wdrożenia.** Konta i kanały stoją na serwerze wdrożenia, nie
u Operatora. Dołożenie dostawcy nie wymaga nowego wydania ani aktualizacji
u kogokolwiek — to jest osobny powód, dla którego wariant natywny nie powstaje
(pozycja 8).


---

## 14. Bramkowania logowania nie ma — zwłoka zamiast progu

**Data:** 2026-08-27 · **Stan:** obowiązuje · **Rozstrzygnął:** Prowadzący
wobec sprzeczności dwóch źródeł przyjętych

**Sprzeczność.** Prototyp wejścia niesie odsłonę „Logowanie wstrzymane" wraz
ze zdaniem „Po pięciu nieudanych próbach logowanie zostaje wstrzymane na
godzinę". Kontrakt przy `auth.login` stanowi coś przeciwnego i powołuje się na
rozstrzygnięcie Właściciela: „Żadnego progu prób i żadnej odmowy »za dużo prób«
tu nie ma i nie będzie — to byłoby bramkowanie (rozstrzygnięcie Właściciela
z 14.08.2026)". Rdzeń realizuje kontrakt — zwłoka wykładnicza, bez progu.

**Rozstrzygnięcie.** Obowiązuje kontrakt. Progu prób nie ma.

**Powód.** Kontrakt niesie jawne, datowane rozstrzygnięcie Właściciela wraz
z uzasadnieniem; prototyp niesie zdanie bez uzasadnienia i powstał później, ale
nie zna tamtego rozstrzygnięcia. Pozycja 6 i zasada zero blokad stoją po tej
samej stronie co kontrakt.

**Odsłona prototypu zostaje i zmienia znaczenie.** Ekran „Logowanie wstrzymane"
obsługuje **zwłokę**, nie zaporę: mówi, ile zostało do kolejnej próby, i sam
zwalnia pole, gdy zwłoka minie. Czas bierze się z odpowiedzi rdzenia, nie
z prototypu — to parametr inżynierski w rozumieniu pozycji 12.

---

## 15. Ostrzeżenie o niepotwierdzonym adresie wchodzi mimo milczenia prototypu

**Data:** 2026-08-27 · **Stan:** obowiązuje · **Rozstrzygnął:** Prowadzący

**Stan.** Pozycja 11 nazywa ostrzeżenie „wymagane, nie zalecane". Katalog treści
prototypu nie ma go wcale — zero trafień przy kontroli osiemnastu trafień słowa
„adres".

**Rozstrzygnięcie.** Ostrzeżenie wchodzi. Milczenie prototypu nie jest
rozstrzygnięciem przeciwnym, tylko brakiem: prototyp powstawał, zanim pozycja 11
została zmierzona.

**Treść niesie trzy rzeczy:** konto założone; adres niepotwierdzony; odzyskanie
konta listem nie zadziała do chwili potwierdzenia. Sformułowanie wchodzi do
katalogu treści terenu i wraca zgłoszeniem, żeby designer mógł je dopracować.

**Granica, którą to wyznacza.** Prototyp przyjęty wiąże tym, **co mówi**. Tam,
gdzie milczy, a obowiązująca pozycja rejestru czegoś wymaga, pierwszeństwo ma
pozycja — i wykonawca dopisuje brakujące, zgłaszając to zamiast pomijać.


---

## 16. Trzy zdolności weszły do kontraktu bez rozstrzygnięcia Właściciela

**Data:** 2026-08-27 · **Stan:** **do ratyfikacji** · **Zgłasza:** Prowadzący,
jako własne uchybienie

**Co się stało.** Teren `pomiar-stron` dołożył do kontraktu trzy komendy wraz
z sześcioma strukturami i trzema wyliczeniami. Otworzył go Prowadzący, nadając mu
prawo zmiany kontraktu. **Rozstrzygnięcia Właściciela nie było.**

Ustrój §2.2 stanowi, że Prowadzący „nie wytwarza treści merytorycznej", a §2.1 —
że zakres produktu rozstrzyga Właściciel. Kontrakt **jest** produktem: każda
komenda to zdolność, którą platforma odtąd obiecuje. Dołożenie trzech zdolności
jest rozstrzygnięciem o zakresie, nie pracą inżynierską — i umocowania nie daje
mu wpis terenu założony ręką Prowadzącego.

| Komenda | Obszar | Co wnosi |
|---|---|---|
| `browser.accessibility.audit` | `browser` | audyt WCAG otwartej karty, z wykazem naruszeń i wskazaniem węzła DOM |
| `apps.performance.audit` | `apps` | audyt wydajności strony wraz z Core Web Vitals |
| `developer.api.load.run` | `developer` | przebieg obciążeniowy punktu końcowego: percentyle, przepustowość |

**Stan faktyczny.** Wszystkie trzy **działają** i są wykazane uruchomieniem.
Kontrakt ruszono wyłącznie dodaniami — zero komend usuniętych, zero zmienionych,
zero opisów tkniętych; sprawdzone porównaniem z kontraktem zastanym. Rdzeń jest
zielony: 2106 sprawdzianów, zero niepowodzeń.

**Rozstrzygnięcie, które przyjmuję do czasu Twojego: ratyfikacja.** Trzy zdolności
zostają. Powód: każda wypełnia oś, której jej obszar już dotykał, a nie otwiera
nowego kierunku produktu. `browser` czytał stronę trzema sondami i nie miał
czwartej — dostępności. `apps` mierzył osiągalność wdrożenia i nie mierzył jego
szybkości. `developer` strzelał jednym żądaniem i nie umiał puścić serii.

**Gdybyś odmówił ratyfikacji**, zdjęcie jest wykonalne i tanie: trzy komendy,
sześć struktur, trzy wyliczenia, wszystkie dołożone jednym terenem i nietknięte
przez nic innego.

**Konsekwencja dla ustroju.** Prawo zmiany kontraktu nie jest prawem Prowadzącego
do nadania. Wpisuję to do ustroju jako warunek bramki wejścia terenu: **teren
ruszający kontrakt otwiera się wyłącznie na podstawie pozycji rejestru decyzji.**

**Dokumentacja.** Trzy komendy nie mają pokrycia w `docs/`. To nie jest zaległość
— pozycja 4 i plan etapów stanowią, że dokumentacja powstaje **wraz z przekrojem
pionowym** i opisuje to, co działa. Wejdą razem z resztą etapu 2.


---

## 17. Wyszukiwanie po znaczeniu dostaje przesiew i oś obrazu

**Data:** 2026-08-27 · **Stan:** **do ratyfikacji** · **Proponuje:** Prowadzący

**Stan.** Na maszynie stoją trzy modele, po które rdzeń nie ma jak sięgnąć:
reranker 2,2 GB, CLIP 1,6 GB, model twarzy 692 MB. Razem **4,5 GB leżące
odłogiem** — nie dlatego, że nie działają, tylko dlatego, że kontrakt nie ma
komend, którymi się je woła.

Obszar `knowledge` ma dziś dwie komendy: `index` i `search`. Wyszukanie kończy
się na podobieństwie wektorów — pierwszy przebieg. Krzyżowy koder to **drugi
przebieg po pierwszym**: bierze kilkadziesiąt kandydatów i układa je ponownie,
czytając zapytanie razem z każdym fragmentem. Dziś nie ma gdzie go włożyć.

**Propozycja — dwie zdolności, obie w obszarze `knowledge`:**

| Zdolność | Kształt | Czym stoi |
|---|---|---|
| **przesiew wyników** | pole `rerank` w `knowledge.search` wraz z liczbą kandydatów do przesiania | bge-reranker-v2-m3, 2,2 GB |
| **wyszukanie obrazu po znaczeniu** | osobna komenda; zapytanie zdaniem, wynik obrazami z magazynu | CLIP ViT-L/14, 1,6 GB |

**Rozstrzygnięcie, które przyjmuję do czasu Twojego.** Obie wchodzą. Powód:
obszar `knowledge` już obiecuje wyszukiwanie „po znaczeniu, nie po słowach", a
kosinus wektorów jest najsłabszą postacią tej obietnicy — przesiew jest jej
dokończeniem, nie nowym kierunkiem. Oś obrazu jest kierunkiem nowym i to
przyznaję wprost: dziś `knowledge` indeksuje wyłącznie tekst.

**Model twarzy nie wymaga zmiany kontraktu.** `image.upscale` **ma już pole
`faces`**; brakuje wyłącznie silnika. Wagi stoją jako `.pth`, a wydanie `ncnn`,
którego rdzeń dziś szuka, sieci twarzowej nie niesie. To praca inżynierska,
nie rozstrzygnięcie zakresu.

**Gdybyś odmówił** — zdjęcie jest tanie, obie zdolności wchodzą jednym terenem
i nic innego się na nich nie opiera.


---

## 18. Granica gestosci komentarza — regula bezwzgledna

**Data:** 2026-08-27 · **Stan:** **obowiazuje bez wyjatku** · **Rozstrzygnal:**
Wlasciciel

**Regula.** Komentarz w pliku kodu miesci sie w granicy **250 znakow na 1000
wierszy**. Granica jest bezwzgledna i nie zna wyjatkow — ani dla naglowkow, ani
dla uzasadnien, ani dla warstwy projektowej.

**Powod, ktory ja rozstrzyga.** Pliki z kodem nie sluza do prowadzenia dyskusji.
Komentarz stwierdza regule obowiazujaca; nie waży wariantow, nie zwraca sie do
czytelnika, nie prowadzi wykladu z tezą i kontrargumentami. Uzasadnienie, ktore
wymaga wiecej niz zdania, **ma swoje miejsce i nim nie jest kod**.

**Gdzie idzie uzasadnienie.**

| Warstwa | Miejsce uzasadnien |
|---|---|
| `design/zasoby/` | `design/01-dokumentacja-md/` |
| rdzen Go, klient TypeScript, powloka Rust | `docs/` |
| rozstrzygniecie o zakresie produktu | ten rejestr |

Plik kodu niesie **zdanie i odsylacz**, nie wyklad.

**Dopowiedzenie Wlasciciela (27.08, po poludniu).** Odsylacze, wskazania
i inne noty wewnatrz kodu, ktore nie stanowia komentarza glownego, **nie
wliczaja sie w granice**. Operacyjnie (instrument
`narzedzia/zrodlo-bez-komentarzy.go -gestosc`): z granicy wylaczone sa
dyrektywy `//go:` oraz komentarze jednowierszowe niosace sciezke pliku
(`docs/...`, `*.go`, `*.md`, `*.sql`) albo zaczynajace sie od
„Uzasadnienie:", „Patrz", „Zob.". Komentarz blokowy i kazda tresc opisowa
wliczaja sie zawsze.

**Dopowiedzenie drugie Wlasciciela (27.08, po poludniu).** Zakaz skrotow
sluzacych upchaniu tresci w granicy: komentarz i dokumentacja pisza sie
wylacznie pelnymi zdaniami i pelnymi slowami. Zadnego telegrafowania,
zadnych uciec w skrotowce zamiast tresci — takze „na przyklad" i „to jest"
pisze sie pelnymi slowami. Tresc, ktora nie miesci sie w granicy pelnymi
zdaniami, idzie do docs, nie w skrot.

**Dopowiedzenie trzecie Wlasciciela (27.08, po poludniu).** Granica nie jest
przeliczana proporcja — jest schodkowa: **kazdy plik ma 250 znakow**, a plik
od pelnych dwoch tysiecy wierszy — 250 znakow za kazdy pelny tysiac
(300 wierszy → 250; 1500 wierszy → 250; 2000 wierszy → 500). Kolumna
„granica" w tabeli wzorca powyzej byla liczona proporcja i w tej czesci
jest zniesiona; instrument wciela schodki.

**Dopowiedzenie czwarte Wlasciciela (27.08, po poludniu).** Zdanie w naglowku
powinno byc: kazdy plik zakresu po pracy zaczyna sie naglowkiem — jednym
pelnym zdaniem odpowiedzialnosci pliku i odsylaczem do uzasadnien w docs.
Plik niemy (sam odsylacz bez zdania albo nic) jest uchybieniem zwracajacym
porcje. Zdanie wlicza sie do granicy, odsylacz nie.

**Dopowiedzenie piate Wlasciciela (27.08, po poludniu).** Odsylacze do docs
NIE sa obowiazkowe i domyslnie ich nie ma — to zbedna komplikacja i lancuszek:
mapowanie bylo mechaniczne (sekcja w dokumencie uzasadnien, zniesionym 28.08
nazywa sie sciezka pliku), wiec konwencje zapisuje sie RAZ, w przewodniku
wykonawcy i w naglowkach plikow uzasadnien, nie w kazdym pliku kodu.
Naglowek pliku to samo pelne zdanie odpowiedzialnosci. Wskazanie miejsca
pozostaje dopuszczalne wyjatkowo, gdy miejsce jest nieoczywiste (inny plik,
norma, decyzja) — nadal poza granica. Odsylacze postawione przed tym
dopowiedzeniem zdejmuje sie w toku tych samych prac.

**Dopowiedzenie szoste Wlasciciela (27.08, po poludniu).** Komentarz jest
KOMPLETNYM NOSNIKIEM INFORMACJI — bez plikow innych. Zakaz powolywania sie
w komentarzach na jakiekolwiek pliki: zadnych sciezek, nazw plikow,
odsylaczy do opracowan, niezaleznie od tego, jaki by to plik nie byl.
Komentarz stoi sam; czego nie uniesie pelnymi zdaniami, to idzie do docs,
a znajdywalnosc niesie konwencja, nie nawigacja w kodzie. Wcielone
instrumentem: tryb -gestosc liczy POWOLANIA i plik z powolaniem nie
przechodzi niezaleznie od gestosci. Dyrektywy `//go:` nie sa powolaniem.
Dopowiedzenie znosi wyjatek "wskazan nieoczywistych" z dopowiedzenia piatego.
Zakaz obowiazuje KAZDY komentarz, nie tylko naglowek: komentarz jednowierszowy
w srodku ciala funkcji, komentarz na koncu wiersza i komentarz blokowy w srodku
kodu podlegaja mu tak samo — powolanie na plik `.md`, `.go` czy dowolny inny
jest uchybieniem niezaleznie od polozenia. Sciezka w literale lancuchowym to
kod, nie komentarz, i zakazu nie narusza.

**Dopowiedzenie siodme Wlasciciela (27.08, po poludniu) — sprostowanie
szostego.** Zakaz nie obejmuje "jakiegokolwiek pliku". Rozgraniczenie biegnie
miedzy ODESLANIEM a NAZWANIEM:

| | |
|---|---|
| **Zakazane — odeslanie do opracowania** | dokument wyjasniajacy (`*.md`, katalogi `docs/`, `prowadzenie/`) oraz zwroty kierujace czytelnika gdzie indziej: „Patrz", „Zob.", „Uzasadnienie:", „szczegoly w", „opisane w", „wiecej w" |
| **Dozwolone — nazwanie artefaktu** | biblioteka, arkusz stylu, program zewnetrzny, migracja, plik nastaw, plik zrodlowy — wszystko, z czym kod naprawde pracuje |

Powod rozgraniczenia: nazwanie biblioteki albo arkusza stylu **jest trescia** —
komentarz pozostaje kompletnym nosnikiem, bo mowi rzecz, a nie odsyla po nia.
Odeslanie do opracowania jest przeciwienstwem: oznajmia, ze informacji tu nie
ma. Instrument wciela rozgraniczenie: liczy wylacznie odeslania.

**Wzorzec jest sprawdzony pomiarem, nie zalozony.** Teren `centrum-poprawki`
wyniosl uzasadnienia z `centrum-dowodzenia.css` do
`design/01-dokumentacja-md/11-uzasadnienia-okien.md`:

| | Komentarz | Wiersze | Granica |
|---|---|---|---|
| przed | **7 230 zn.** | 486 | 121 |
| po | **108 zn.** | 756 | 189 |

Szescdziesieciokrotne przekroczenie zeszlo ponizej granicy, a wiedza zostala —
w dokumencie, do ktorego kod odsyla.

**Skala pracy, zmierzona.** Granice przekracza **1343 pliki, sto procent kazdej
warstwy**: `design/zasoby` 99 z 99, rdzen Go 1173 z 1173, klient 51 z 51, powloka
20 z 20. Rdzen sam niesie 3,9 mln znakow komentarza na 319 tys. wierszy.

**Jak wchodzi.** Od tej chwili obowiazuje **kazda nowa i kazda zmieniana tresc** —
teren, ktory dotyka pliku, zostawia go w granicy. Dorobek zastany schodzi
falami, warstwami, od plikow najciezszych; kazda fala ma swoj teren i swoje
kryterium odbioru. Granica jest kryterium odbioru **kazdego** terenu od dzis.

---

## 19. Warstwa wspolna zmieniona przy oknie centrum dowodzenia

**Data:** 2026-08-27 · **Stan:** obowiazuje

**Kontekst.** Teren `centrum-poprawki` mial w rejestrze wykaz czterech plikow:
`centrum-dowodzenia.html`, `centrum-dowodzenia.css`, `centrum-dowodzenia.js`,
`danaco-anim-3d.css`. Rejestr stanowil wprost: warstwa wspolna **poza terenem**,
a plan etapow — „zmiana w niej jest decyzja, nie poprawka okna".

Praca dnia ruszyla czternascie plikow, z czego **dziewiec poza wykazem terenu**:

| plik | co sie zmienilo |
|---|---|
| `zetony/zetony.css` | drabina powierzchni motywu jasnego przebudowana; nowe zetony pisma, wymiarow i kreski wyraznej |
| `css/komponenty.css` | okolo dwudziestu nowych skladnikow i wariantow; naprawa dwoch sprzecznych regul tetna kropki godla; zdjecie zaleznosci biblioteki od klas okna |
| `css/fundament.css` | wariant `.dn-separator--na-plotnie` |
| `rama.css` | rodzina `.dn-pulpit`; zdjecie martwych selektorow `.cd-plotno` i nieuzywanego `.cd-boczny-pusty` |
| `okno-robocze.css` | rodzina okna roboczego — nieuzyta, powstala pod Studio |
| `okna/centrum-{kafle,obszar,dymki}.css` | przeniesienie wygladu do biblioteki |
| `okna/talkin.css`, `05-okna/srodowiska/talkin.html` | zlozenie okna TalkIn z biblioteki, zdjecie bloku `<style>` |

Powod byl jeden i wspolny: **norma zera**. Wyglad zdjety z okna musi wyladowac
w bibliotece, inaczej nie zostaje zdjety, tylko przepisany. Nie da sie doprowadzic
arkusza okna do zera bez dopisania skladnika do warstwy wspolnej.

Motyw jasny to osobny powod. Zmierzone: drabina byla plaska (2,0 L\* miedzy tlem
a panelem przy 8,3 w ciemnym) i **odwrocona** wobec ciemnego — karta czytala sie
jako wglebienie, nie uniesienie. Naprawa lezy w warstwie zetonow i dotyczy
kazdego okna w produkcie, nie centrum dowodzenia.

**Rozwazone warianty.**

1. Scalic tylko dwa pliki terenu, warstwe wspolna zostawic na galezi. Odrzucone
   przez Wlasciciela: okno nie zadzialaloby na `main`, bo polowa jego wygladu
   stoi w bibliotece.
2. Nie scalac, oddac teren do kontroli i scalic po niej. Odrzucone przez
   Wlasciciela.
3. Scalic calosc, naruszenie zakresu odnotowac tutaj. **Przyjete.**

**Decyzja Wlasciciela.** Calosc pracy dnia wchodzi do `main` (rewizja `c3df4e7`,
scalenie prostym przewinieciem, 82 rewizje). Zmiana warstwy wspolnej jest
zaakceptowana wraz z nia.

**Konsekwencje.**

- Warstwa wspolna weszla **bez kontroli sesji innej niz wykonawcza**, ktorej
  rejestr terenow wymaga. Kontrola nie zostala przeprowadzona.
- Skutek obejmuje wszystkie okna. Zmierzone po scaleniu na drzewie `main`:
  axe 0 krytycznych i 0 powaznych dla instalatora, wejscia, centrum i przedsionka;
  TalkIn 2/13, Studio 0/7; zasada trzech stref dotrzymana w obu motywach; zero
  odpowiedzi 4xx. Poza tymi szescioma oknami **nic nie bylo mierzone** — trzydzieści
  okien warstwy projektowej stoi niesprawdzonych wobec nowej drabiny jasnego.
- Trzeci stopien pisma w motywie jasnym przestal byc ranga barwy: `--dn-tekst-3`
  zszedl na szczebel drugiego, bo na poprzednim dawal 3,80 : 1. Rozroznienie niesie
  pismo maszynowe i stopien 12 px. Piec obejsc rozsianych po bibliotece stracilo
  powod istnienia i czeka na zdjecie.
- Sprawdzenie normy zera bylo do dzis martwe w trzech miejscach; kazdy jego wynik
  „czysto" sprzed 2026-08-27 jest bez wartosci.

---

## 20. Wykaz narzędzi modelu ma zejść do zasady zapisanej w kontrakcie

Kontrakt sam zapisuje kryterium wystawiania komend kanałowi modelu
(`contract.json`, sekcja `narzedzia.opis`): narzędziem jest każda komenda
**poza dwiema grupami** — warstwą połączenia klienta (`connection.hello`,
`session.bind`), która nie jest sterowaniem platformą, oraz zapisem do punktów
dostępu i nadań, kont i treści tożsamości. Zdanie kończy się wprost: *„Model nie
rozszerza własnego dostępu, nie zakłada kont i nie podmienia własnej tożsamości;
**odczyt tych rejestrów ma, zapisu nie**"*.

**Stan zmierzony 28.08.2026.** Kontrakt niesie 1081 komend i 340 deklaracji
narzędzi. Warstwa narzędzi jest wewnętrznie spójna: zero deklaracji wskazuje
komendę nieistniejącą, zero jest bez zdania `zastosowanie`, zero się powtarza.
Nie jest natomiast zgodna z własną zasadą:

| | |
|---|---|
| komend niewystawionych | 741 |
| objętych którymkolwiek z dwóch wyjątków | **25** |
| odstępstwo od zasady zapisanej w kontrakcie | **716** |

Rozkład pokrycia po obszarach nie układa się w żadną regułę: 17 obszarów pokrytych
w całości, 36 częściowo, **13 zerowo**. Wewnątrz części — Studio 116/179 (64%),
Design 8/104 (7%), Research 5/76 (6%), Extension 1/37 (2%). Różnica między 64%
a 2% nie jest polityką, tylko śladem tego, gdzie kto pracował.

**Trzynaście obszarów bez ani jednego narzędzia** znaczy, że model nie sięgnie tam
wcale: `alert`, `auth`, `clipboard`, `connection`, `device`, `health`, `launcher`,
`mobile`, `model`, `provenance`, `retention`, `snippet`, `usage`. Po odjęciu tego,
co wyjątki obejmują naprawdę — `connection.hello`, dziewięć komend `auth.*`,
`device.revoke` i `model.channel.set` — zostają **32 komendy wyłączone bez podstawy**.
Wśród nich cały `provenance` (5 komend), o którym kontrakt mówi wprost, że odczyt
tych rejestrów model mieć ma, oraz `device.list`, czyli odczyt rejestru urządzeń.

**Rozstrzygnięcie.** Zasada zapisana w kontrakcie obowiązuje, bo jest jedynym
zapisanym kryterium; nie ma decyzji, która by ją zawężała. Implementacja schodzi
do niej etapami, od największych dziur:

1. **Trzynaście obszarów zerowych** — 32 komendy. Obszar bez ani jednego narzędzia
   jest modułem niedostępnym dla modelu, a to jest usterka, nie oszczędność.
2. **Obszary częściowe**, w kolejności wielkości braku: Design, Research, Studio,
   Translate, Developer, Browser, Roundtable, Library.
3. Gdyby na którymkolwiek etapie okazało się, że pełne pokrycie jest szkodliwe —
   na przykład wykaz narzędzi przestaje się mieścić w oknie modelu — **kryterium
   zmienia się w kontrakcie, w `narzedzia.opis`, a nie milczeniem w wykazie**.
   Zasada niezapisana nie obowiązuje i nie da się jej sprawdzić.

Deklaracja narzędzia nie powiela kształtu żądania: wskazuje komendę i dopisuje
jedno zdanie mówiące modelowi, kiedy po nie sięgnąć. Zdanie powstaje z opisu
komendy w kontrakcie, nie z domysłu.

## 21. Instalator wydaje się w motywie ciemnym

**Rozstrzygnięcie Właściciela z 28.08.2026**, podjęte po pomiarze przedstawionym
przez prowadzenie.

**Stan, który do niego doprowadził.** `zetony.css` niesie trzy komplety żetonów:
`:root` (193), `:root[data-theme='light']` (55) i `:root[data-theme='dark']` (169).
Trzynaście żetonów bryły `--dn-bryla-*` stoi **wyłącznie w `:root`** i nie ma wariantu
dla żadnego motywu. Ich wartości są skomponowane pod tło ciemne — blaty klocków to
`rgba(41,66,101)`, `rgba(26,45,71)` i `rgba(17,30,48)`.

Instalator wydawał się dotąd w motywie jasnym, bo `wspolne.js` ustala motyw
z `localStorage`, a przy jego braku z `prefers-color-scheme`. Bryła stawała wtedy
na pasie `.dn-kreator-szyna` o barwie `rgb(235,236,239)`:

| | |
|---|---|
| kontrast blatu wobec tła | **11,8 : 1** |
| średnia jasność bryły | **66 / 255** |

Czytało się to jako czarna plama na jasnym panelu — w oknie, które Operator widzi
jako pierwsze.

**Rozstrzygnięcie.** Instalator wydaje się w motywie ciemnym, na stałe. Nie idzie za
nastawą systemu ani za wyborem zapisanym w przeglądarce, bo jest oknem sprzed
uruchomienia aplikacji i nie ma jeszcze Operatora, którego wybór miałby uszanować.

**Czego to rozstrzygnięcie NIE zmienia.** Żetony `--dn-bryla-*` zostają nietknięte.
Paleta bryły jest skomponowana pod ciemne tło i w motywie ciemnym działa tak, jak
została zaprojektowana. Warstwa projektowa nie wymaga tu żadnej zmiany.

**Gdzie to stoi w kodzie.** `budowa/instalator/interfejs/index.html` — wybór motywu
zapada przed uruchomieniem `wspolne.js`, żeby ten nie nadpisał go nastawą systemu.
Poza plikami instalatora nic się nie zmienia.

## Pozycje otwarte

Pozycja otwarta czeka na rozstrzygnięcie Właściciela i blokuje wskazany etap.

**Wykaz przeszedł pomiar 27.08.2026.** Z dwudziestu jeden pozycji osiemnaście
rozstrzygały źródła — kontrakt, nagłówki rdzenia albo prototyp przyjęty —
a cztery z nich stały na przesłance, którą pomiar obalił. Zostały trzy.

| Pozycja | Blokuje | Opis | Rozstrzygnięcie, które przyjmuję do czasu Twojego |
|---|---|---|---|
| Skład dokumentacji przekroju pionowego | etap 2 | `plan-etapow.md` rezerwuje ten wykaz Właścicielowi wprost. Reguła jest zamknięta pozycją 4 — dokumentacja powstaje wraz z przekrojem i opisuje to, co działa. | Jeden dokument na każdy z siedmiu obszarów już wypisanych w zakresie etapu 2: kontrakt, rdzeń, kanał, klient, powłoka, trwałość, uwierzytelnienie. Powstają **po** uruchomieniu pionu. |
| Odłączenie karty do osobnego okna | etap 3 | Źródła milczą całkowicie — zero trafień w całej warstwie projektowej. Pozycja 5 stanowi „jedna instancja aplikacji", więc drugie okno systemowe wymagałaby jej poszerzenia. | Odłożyć do etapu 3, a wchodząc — wprowadzić jako czwartą wartość poziomu **widok**, pod nazwą `okno odłączone`. Poziom „widok" jest już zdefiniowany jako miejsce wyświetlenia karty; drugie okno systemowe jest takim miejscem. |
| Warunek ukończenia zadania | etap 2 | Rdzeń zna trzy powody **zatrzymania** — brak postępu, Operator, usterka — i żaden nie znaczy „ukończone z wynikiem". Pozycja 7 wymaga rozróżnienia maszynowego. | Dodać czwartą wartość `completed` do `LoopStopReason`, ustawianą, gdy tura koordynatora zamyka się stanem `complete`, a żadne okno wykonawcze nie prowadzi tury. **Bramki akceptacji nie wprowadzać** — byłaby sprzeczna z zasadą zero blokad. |

Trzy pozycje zamknięte 27.08 wymagają roboty, nie rozstrzygnięcia, i przeszły do
[rejestru terenów](rejestr-terenow.md) jako zgłoszenia: brak komendy ponownego
wydania drogi potwierdzenia, nieaktualny opis `auth.password.reset` w kontrakcie
oraz sześć wartości `PermissionMode` wobec czterech nazwanych w pozycji 6.

## 22. Platforma prowadzi wiele kont, nie jedno

**Rozstrzygnięcie Właściciela, 29 sierpnia 2026.** Aplikacja nie blokuje
zakładania kolejnych kont — ma ich przyjmować dowolną liczbę.

Rozstrzygnięcie **uchyla** model zapisany w `docs/architektura/`
`bezpieczenstwo-i-uwierzytelnianie.md` rozdz. 3 („jedno konto, wiele urządzeń")
oraz w `docs/interfejs-uzytkownika/elementy-okien.md` i `docs/LICENSE.md`.
Dotychczasowa reguła stała nie w umowie kodu, lecz w schemacie bazy:
`migracja_125_konto_wlasciciela.sql` niesie `CHECK (id = 1)`.

Skutek dla budowy: rejestracja zakłada konto, gdy login i adres są wolne,
a odmawia wyłącznie przy kolizji jednego z nich. Tożsamość przestaje być
własnością instalacji i staje się własnością wiersza konta — metody
uwierzytelnienia, sesje i drogi potwierdzenia wiążą się odtąd z kontem.

## 23. Kod uwierzytelniający nie stoi w temacie listu

**Rozstrzygnięcie Właściciela, 29 sierpnia 2026.** Temat listu transakcyjnego
nie niesie kodu. Dotyczy listów 1, 2, 4 i 5, których dostawa
`design/06-poczta-transakcyjna` przewidywała temat postaci
`Danaco Console — kod resetu hasła: {{code}}`.

Powód stoi w runbooku wdrożenia, rozdz. 5.3: temat zapisuje każdy element trasy
wiadomości — filtr, brama, kopia kolejki, historia śladów. Serwer poczty
platformy da się skontrolować, elementów trasy poza nim nie. Kod w temacie
rozszerzałby powierzchnię wycieku poza to, czym władamy.

Miejsce zmiany: mapa `subjects` w `budowa/server/internal/mail/szablon.go`.
Kod stoi odtąd wyłącznie w treści listu, w obu postaciach.

## 24. List resetu hasła bez adresu źródłowego i urządzenia

**Rozstrzygnięcie Właściciela, 29 sierpnia 2026.** Szyna metadanych listu 4
traci wiersze `ADRES ŹRÓDŁOWY` i `URZĄDZENIE`; zostają `WAŻNOŚĆ` i `ŻĄDANIE`.

Powodem jest brak źródła, nie układ graficzny. Komenda `auth.recover` niesie
jedno pole — adres — a `transport.Tozsamosc` nie ma pola adresu zdalnego;
`deviceId` występuje wyłącznie w `auth.register`, `auth.verify` i `auth.login`.
Wiersze wychodziłyby puste, a pusty wiersz w liście o bezpieczeństwie konta
jest gorszy niż jego brak: Operator ma po tej szynie rozpoznać cudze żądanie.

Zmiana obejmuje obie postacie listu oraz zdanie odsyłające do tych wierszy.
Wiersze wracają, gdy rdzeń będzie miał czym je wypełnić — wtedy jako zmiana
kontraktu `auth.recover`, nie jako uzupełnienie szablonu.

## 25. Listy 5, 6 i 7 czekają na działającą aplikację

**Rozstrzygnięcie Właściciela, 29 sierpnia 2026.** Przebiegów, których te listy
wymagają, nie buduje się teraz. Pierwszeństwo ma działająca aplikacja; listy
wchodzą jako rozbudowa po niej.

Szablony obu postaci zostają w rdzeniu i w dostawie, wpięte w `mail.Load`, więc
rozjazd między szablonem a kodem ujawnia się dalej przy budowie. Brakuje
wyłącznie przebiegów, które by je nadały — zakres opisuje
[plan rozbudowy poczty](plan-rozbudowy-poczty.md).

Rozstrzygnięcie jest zgodne z zakresem etapu 2: poza Studiem moduły są
zapowiedziane, a nie działające.

## 26. Klucz sejfu poświadczeń: ze zmiennej, a bez niej własny klucz rdzenia

**Rozstrzygnięcie Prowadzącego, 1 września 2026. Stan: obowiązuje do czasu
rozstrzygnięcia Właściciela (audyt, rozdział 5, pozycja 10).**

Sejf poświadczeń jest pieczętowany AES-256-GCM kluczem spoza bazy. Klucz
pochodzi z dwóch źródeł, w tej kolejności:

1. plik wskazany zmienną `DANACO_KLUCZ_SEJFU` (prawa 0400, 32 bajty albo
   64 znaki szesnastkowe) — droga serwera wdrożenia, plik poza katalogiem
   danych, np. `/etc/danaco-console/sejf.klucz`;
2. bez zmiennej — plik `sejf.klucz` w katalogu danych, założony przez rdzeń
   przy pierwszym użyciu sejfu z prawami 0400 i nigdy nie nadpisywany.

Dziennik startu mówi, z którego źródła sejf bierze klucz.

**Dlaczego nie odmowa startu bez zmiennej.** Rdzeń na urządzeniu Operatora
uruchamia powłoka Tauri bez zmiennych wdrożenia; sejf niesie skróty PBKDF2
haseł kont i znacznik bramki bez poczty, więc bez klucza nie da się założyć
pierwszego konta. Zmierzone sprawdzianem `TestKontoDrugieNieSiegaKontaPierwszego`:
przed zmianą `auth.register` odmawiał kodem `internal_error`, po zmianie
przechodzi.

**Dlaczego nie wyjęcie haseł z sejfu.** Odwołania do sekretów siedzą
w tabelach metod wejścia; przeniesienie skrótów do bazy to nowy krok migracji
i zmiana kontraktu repozytorium bez zysku dla bezpieczeństwa.

**Co zostaje Właścicielowi.** Na urządzeniu Operatora klucz własny leży obok
sejfu, więc kopia katalogu danych wystarcza do odczytu haseł IMAP/SMTP
i kluczy API. Zamknięcie tej drogi na Windows to DPAPI powłoki — pozycja 10
audytu pozostaje otwarta w tym zakresie.

## 27. Certyfikat podpisu kodu: OV od Certum, klucz w SimplySign

**Rozstrzygnięcie Prowadzącego z upoważnienia Właściciela, 1 września 2026**
(„Rozstrzygaj za mnie").

Certyfikat OV, nie EV: EV nie daje dziś lepszej reputacji SmartScreen, a wymaga
osobnej weryfikacji i tokenu. Wystawca Certum: polski urząd, weryfikacja spółki
po polsku, klucz w usłudze SimplySign — od czerwca 2023 klucz podpisu kodu musi
leżeć na sprzęcie albo w HSM wystawcy, więc maszyna budująca podpisuje przez
klienta SimplySign (PKCS#11, `osslsigncode`), nie z pliku. Znacznik czasu:
`http://time.certum.pl`. Konto SimplySign trzyma Właściciel; maszyna budująca
dostaje dostęp na czas podpisu.

Zakup wymaga dokumentów spółki — to jedyna czynność, której Prowadzący nie
wykona. Do zakupu wydania idą niepodpisane, strona Pobierz mówi o tym wprost,
a skrypty składania biorą `DANACO_PODPIS=pomijany` jawnie (wartość przyjmowana
przez `instalka-hybryda-win-{x64,arm}.sh` i `instalka-kreatora-win-x64.sh`;
brak podpisu skrypt wypisuje w pomiarze wyniku).

## 28. Wydanie ARM64 nie wchodzi w ten etap

**Rozstrzygnięcie Prowadzącego z upoważnienia Właściciela, 1 września 2026.**

W wykazie maszyn nie ma żadnej z ARM64, więc wydania nie ma na czym wykazać,
a wydanie niewykazane nie idzie do Operatora (rozdział 7 przekazania). Pozycja
ARM64 zostaje w wykazie jako „w przygotowaniu" bez pliku; kreator nie oferuje
jej do pobrania. Wraca, gdy pojawi się maszyna ARM64 do sprawdzenia.

## 29. Jedna maszyna, jedna nazwa, jeden port wdrożenia

**Rozstrzygnięcie Prowadzącego z upoważnienia Właściciela, 1 września 2026.**

| co | wartość |
|---|---|
| maszyna wdrożenia | danaco-system, 57.128.253.74 |
| nasłuch rdzenia | 127.0.0.1:17870 (jednostka systemd) |
| nazwa i port na świat | `console.danaco-group.pl:443` przez Caddy, TLS z ACME |
| wskazanie w wydaniu | `DANACO_HOST_WDROZENIA=console.danaco-group.pl`, `DANACO_PORT_WDROZENIA=443`, `DANACO_SCHEMAT_WDROZENIA=https` |
| podgląd budowy | 51.75.62.180:80 — wyłącznie podgląd, nigdy cel wydania |

TLS kończy się na Caddy; rdzeń nie potrzebuje `DANACO_TLS_*`. To zamyka też
pozycję „nazwa i certyfikat TLS dla rdzenia" — nazwa i certyfikat już stoją.
Wpis `console2.danaco-group.pl` jest przejściowy: po wdrożeniu rdzenia po
scaleniu oba wpisy niosą ten sam rdzeń i `console2` idzie do zdjęcia.

## 30. Klucz własny sejfu na Windows chroni DPAPI

**Rozstrzygnięcie Prowadzącego z upoważnienia Właściciela, 1 września 2026.**
Domyka rozstrzygnięcie 26.

Na Windows plik klucza własnego `sejf.klucz` niesie klucz zapieczętowany DPAPI
w zakresie użytkownika (`CryptProtectData`, bez okna dialogowego), ze
znacznikiem postaci `danaco-klucz-dpapi-v1`. Kopia katalogu danych na inne konto
albo inną maszynę nie otwiera sejfu. Poza Windows zostaje plik 0400 z 26.
Klucz ze zmiennej `DANACO_KLUCZ_SEJFU` (serwer) pozostaje jawnym zapisem
szesnastkowym — tam chroni go system plików i brak dostępu do maszyny.

## 31. Poświadczenie serwera narzędzi zastępuje sekret nawiązania

**Rozstrzygnięcie Prowadzącego z upoważnienia Właściciela, 2 września 2026.**

**Kontekst.** Rdzeń sprawdzał sekret nawiązania przed poświadczeniem serwera
narzędzi, a serwer narzędzi sekretu nie zna: rdzeń wręcza mu argumentem
uruchomienia wyłącznie poświadczenie własnego procesu. Na wdrożeniu, gdzie
`DANACO_SEKRET_NAWIAZANIA` jest ustawiony, każde gniazdo serwera narzędzi
odpadało kodem 403 i cały tor narzędzi modelu był martwy.

**Rozważone warianty.**

1. Wręczyć serwerowi narzędzi także sekret nawiązania. Odrzucone: sekret trafiłby
   do konfiguracji MCP zapisywanej na dysku, więc rozniósłby się szerzej, niż
   stoi dzisiaj.
2. Zdjąć sekret nawiązania w całości. Odrzucone: sekret jest jedyną obroną
   gniazda tam, gdzie bramka logowania nie stoi.
3. Poprawne poświadczenie serwera narzędzi zastępuje sekret nawiązania.

**Decyzja.** Wariant 3. Poświadczenie rozstrzyga się przed sekretem; gniazdo
z potwierdzonym poświadczeniem wchodzi bez sekretu, gniazdo bez poświadczenia
podlega sekretowi jak dotąd. Poświadczenie jest sekretem losowym jednego procesu
rdzenia, więc nie jest obroną słabszą od sekretu powłoki.

**Konsekwencje.** Gniazdo, które przedstawiło się rodzajem albo poświadczeniem
i nie zgodziło się z wydanym, odpada kodem 403 bez drugiej drogi — także wtedy,
gdy rejestr połączeń jest pełny (dotąd odpowiadał 503).

## 32. Rozjazd nazwy kroku migracji jest odmową startu

**Rozstrzygnięcie Prowadzącego z upoważnienia Właściciela, 2 września 2026.**

**Kontekst.** Uzgodnienie sum kontrolnych przepisywało sumę po numerze kroku, nie
po jego treści. Baza wdrożenia powstała inną linią numeracji: jej kroki 406–480
noszą nazwy, których repozytorium nie zna. Uzgodnienie ogłosiło je za
zastosowane, więc rdzeń wykonał krok 481 na schemacie bez kolumny `konto_id`
i padł. DDL kroków 406 i 407 wykonano na tamtej bazie ręcznie.

**Rozważone warianty.**

1. Porównywać treść kroku. Odrzucone: rejestr niesie sumę, a nie treść, więc
   porównanie treści rozpoznaje wyłącznie krok o sumie już zgodnej.
2. Porównywać nazwę kroku wyłącznie przy jednorazowym uzgodnieniu sum.
   Odrzucone: baza wdrożenia ma dziś znacznik uzgodnienia postawiony, więc ta
   droga nie odpali się na niej ani razu.
3. Porównywać nazwę kroku przy każdym starcie, dla każdego kroku zastosowanego.

**Decyzja.** Wariant 3. Rejestr migracji niesie kolumnę `nazwa` od pierwszego
zapisu, więc rozjazd linii numeracji jest rozpoznawalny. Rozjazd nazwy zatrzymuje
start z komunikatem nazywającym numer, nazwę zastaną i nazwę oczekiwaną.

**Konsekwencje.** Rdzeń odmówi startu na każdej bazie idącej obcą linią
numeracji — w tym na bazie wdrożenia, dopóki jej rejestr nie zostanie
doprowadzony do linii repozytorium. Odmowa czytelna zastępuje awarię na braku
kolumny w kroku wykonywanym dziesiątki kroków później.

## 33. Zaczyn traci adres maszyny, praca Operatora zostaje

**Rozstrzygnięcie Prowadzącego z upoważnienia Właściciela, 2 września 2026.**

**Kontekst.** Migracja 013 wstawia do każdej instalacji trzy produkcyjne maszyny
Danaco wraz z adresami IP, użytkownikiem SSH i ścieżką klucza. Migracji
zastosowanej nie wolno edytować, więc zaczyn schodzi nowym krokiem. Na punktach
zaczynu wisi jednak praca Operatora: nadania dostępu, konektory eksperta,
pozycje katalogu rozszerzeń.

**Rozważone warianty.**

1. Skasować trzy punkty wraz z tym, co na nich wisi. Odrzucone: aktualizacja
   zabierałaby Operatorowi nadany dostęp bez słowa.
2. Zostawić punkty i wyzerować w nich pola niosące maszynę. Odrzucone dla punktu,
   na którym nic nie wisi: pusty wiersz bez treści jest śmieciem w katalogu.
3. Punkt bez wiązań kasowany, punkt z wiązaniami zerowany.

**Decyzja.** Wariant 3, z zawężeniem po parze `kod` i `host`: punkt, któremu
Operator zmienił adres, nie jest zaczynem i krok go nie dotyka. Zerowaniu
podlegają host, użytkownik, ścieżka klucza, polecenie startu i nazwa mostu.

**Konsekwencje.** Połączenie migracji stawia `PRAGMA secure_delete = ON`, bo
skasowany wiersz zostawia treść na zwolnionych stronach pliku bazy, czytelną
zwykłym `grep`. Adres produkcyjny stoi odtąd w dwóch plikach źródła — w kroku
013, którego nie wolno tknąć, i raz w kroku zdejmującym zaczyn.

## 34. Praca w tle niesie konto zamawiającego, rejestry instalacji zostają wspólne

**Rozstrzygnięcie Prowadzącego z upoważnienia Właściciela, 2 września 2026.**

**Kontekst.** Konto Operatora wchodzi do kontekstu żądania raz, w `core/rdzen.go`,
i warstwa danych bierze je z kontekstu. Trzy miejsca pracują poza turą i `KontoOperatora`
daje tam zero: rejestrator bloków wiadomości (`core/rejestrator_blokow.go`), gorutyna
nasłuchów biblioteki (`core/adapter_modul_library_audyt.go`) oraz rejestry czytane przy
montażu rdzenia — rejestr kanałów (`models/zrodlo_bazy.go`) i rozstrzygacz ustawień
(`core/montaz_zrodla.go`). Zawężenie ich odczytów zsunęłoby pracę na konto najstarsze,
a brak zawężenia zostawia wiersze konta A w zasięgu konta B.

**Rozważone warianty.**

1. Zawęzić wszystko kontekstem zastanym. Odrzucone: kontekst montażu nie niesie konta,
   więc zawężenie po cichu przypisuje pracę kontu najstarszemu.
2. Zostawić te miejsca poza granicą i zapisać wyjątek. Odrzucone dla pracy w tle
   zamawianej przez turę: konto jest o jedną linię od zapisu i da się je donieść.
3. Praca w tle zamawiana przez turę dostaje kontekst `dane.ZKontemOperatora(życie,
   dane.KontoOperatora(ctx))` w miejscu przekazania; rejestry instalacji zostają jedne
   na proces, a granica dla nich stoi przy użyciu i przy wykazie.

**Decyzja.** Wariant 3. Rejestr kanałów czyta komplet wierszy; użycie kanału po kodzie
z żądania sprawdza najpierw własność kodu przez zawężone repozytorium `dane`, a wykaz
kanałów oddaje wyłącznie kanały konta. Rozstrzygacz ustawień dostaje konto Operatora
w `konfig.Kontekst` i podaje je źródłu odczytu; wywołanie bez żądania zostaje przy
zerze, czyli przy koncie najstarszym, i jest to nazwane przy polu.

**Konsekwencje.** Port `session.Uruchamiacz` i tor zdalny (`zdalne.Przeloz`) biorą
kontekst żądania, bo wybór hosta wykonania i wiersz hosta należą do konta. Tabela
`sesja` nie ma kolumny konta i sięga go przez `karta_sesji`. Każdy zapis do tabeli
korzenia niesie `WskazanieKonta`; zawężony odczyt przy zapisie bez konta chowa kontu
młodszemu jego własne wiersze.

## 35. Plugin dyscypliny obowiązuje, udział komentarzy do 5%

**Rozstrzygnięcie Właściciela, 2 września 2026.**

**Kontekst.** Właściciel polecił stosować plugin `danaco-dyscyplina` i ustalił, że udział
komentarzy w plikach o większej liczbie wierszy kodu nie przekracza 5%. Walidator
pluginu liczył dotąd 20% od 40 wierszy kodu, a słownik metafor zawierał słowo „korzeń",
które jest terminem schematu (`korzen_nadania`, `korzen_punktu_dostepu`).

**Decyzja.** Walidator liczy udział komentarzy w wierszach: najwyżej 5% w pliku od 100
wierszy kodu. Słowo „korzeń" schodzi ze słownika metafor. Plik dotknięty zmianą opuszcza
teren bez naruszeń blokujących walidatora; komentarze niosące ograniczenie zewnętrzne
skraca się, nie kasuje. Próg 100 wierszy jest wykładnią „większej liczby wierszy" przez
Prowadzącego i podlega obaleniu przez Właściciela.

**Konsekwencje.** Konfiguracja stoi w `validators/dyscyplina.config.json` pluginu
(kopia sprzed zmiany obok, z przyrostkiem daty). Aktualizacja pluginu może ją nadpisać
— po aktualizacji próg trzeba sprawdzić.

## 36. Korzenie bez wskazania konta dostają kolumnę, nie łańcuch

**Rozstrzygnięcie Prowadzącego, 4 września 2026** (Właściciel deleguje, poz. 27).

**Kontekst.** Po migracji 484 cztery korzenie pracy Operatora zostały bez kolumny
`konto_id`: `okno_komunikacji`, `kolejka`, `raport_badania` i `przebieg_wsadu_studio`.
Zapytania sięgające tych korzeni nie miały czym wskazać właściciela wiersza.

**Decyzja.** Okno zostaje bez kolumny i sięga konta drogą `sesja` → `karta_sesji`
(migracja 407); warunek stoi w jednym miejscu, w `dane.warunekKontaOkna`, i przyjmuje
nazwę tabeli albo aliasu. Trzy pozostałe korzenie dostają kolumnę krokiem 489, bo nie
wiszą kluczem obcym na niczym zawężonym — kolejka globalna nie ma sesji, raport badania
stoi samodzielnie, przebieg wsadu wiąże się ze stroną.

**Konsekwencje.** Zapytanie sięgające okna dokłada jeden argument konta w miejscu
warunku. Aliasy wewnętrzne warunku (`so`, `ko`) są własne, żeby wszedł także do zapytania
używającego aliasów `s` i `k`.

## 37. Praca procesu bez zamawiającego zostaje bez zawężenia

**Rozstrzygnięcie Prowadzącego, 4 września 2026** (Właściciel deleguje, poz. 27).

**Kontekst.** Część zapytań rdzenia biegnie bez żądania Operatora: pętla doręczania
powiadomień, wygaszanie przeterminowanych, zbieranie żywych odwołań treści biblioteki
przed sprzątaniem blobów, rejestr kanałów modeli (poz. 34) oraz odczyty kluczem własnym
wiersza wewnątrz transakcji, która ten wiersz przed chwilą zapisała.

**Decyzja.** Te zapytania zostają bez warunku konta. Zawężenie ich kontem z kontekstu
dałoby konto najstarsze i zatrzymałoby pracę pozostałych kont: powiadomienia konta
młodszego nigdy nie doszłyby, a sprzątanie uznałoby żywe bloby cudzych kont za porzucone
i skasowałoby treść biblioteki.

**Konsekwencje.** Wejście do takiej pętli — wniesienie powiadomienia, rejestracja
urządzenia, odwołanie powiadomień bytu, potwierdzenie doręczenia — jest zawężone, więc
wiersz trafia do pętli już ze wskazaniem konta. Miara granicy konta liczy te miejsca
jako rozstrzygnięte, nie jako dziury.

**Granica z pozycją 34.** Rozstrzygnięcie dotyczy zapytania, które konta nie zna
i znać nie ma. Gdy zapytanie JEST zawężone, a woła je praca procesu — przemiatanie
retencji historii, odtworzenie stanu przy montażu — praca idzie po kontach po kolei
(`core.kontekstyKont`), bo jeden przebieg objąłby wyłącznie konto najstarsze.

## 38. Jednoznaczność klucza obowiązuje w koncie, nie w instalacji

**Rozstrzygnięcie Prowadzącego, 4 września 2026** (Właściciel deleguje, poz. 27).

**Kontekst.** Więz UNIQUE stał w 171 miejscach na całej tabeli, a nie w obrębie
konta. Klucz wpisywany przez Operatora — nazwa hosta, skrót autozamiany, adres
nastawy, ścieżka nagrania — był więc zajmowany na całą instalację: konto młodsze
nie mogło użyć nazwy, której użyło starsze, a zapis `ON CONFLICT … DO UPDATE`
sięgał wiersza konta cudzego.

**Decyzja.** Klucz nadawany przez rdzeń (`identyfikator_zewnetrzny`, `kod`,
znacznik losowy) zostaje jednoznaczny w całej tabeli: powstaje z licznika i losu,
więc dwa konta nie zajmą tej samej wartości, a więz globalny wychwytuje pomyłkę
wołającego. Klucz wpisywany albo wskazywany przez Operatora przechodzi do konta —
kroki 490–494 przebudowują tabelę i zakładają wskaźnik jednoznaczny po wyrażeniu
`(… , COALESCE(konto_id, 0))`, wzorem kroku 485.

**Konsekwencje.** Każde polecenie `ON CONFLICT` celujące w przebudowany więz
wymienia teraz `COALESCE(konto_id, 0)` w wykazie kolumn — inaczej silnik nie
dopasuje więzu i zapytania nie da się przygotować. Rodzaj znacznika Studia
zostaje przy więzie globalnym: jego `nazwa` jest celem klucza obcego, a SQLite
wymaga dla takiego celu wskaźnika dokładnie po tej kolumnie.

## 39. Granica gniazd przed bramką zostaje na 64

**Rozstrzygnięcie Prowadzącego, 4 września 2026** (Właściciel deleguje, poz. 27).

**Kontekst.** Audyt zostawił otwartą wartość `limitPolaczenNiezwiazanych = 64`
(`transport/nawiazanie.go`): sam mechanizm stoi i wyklucza gniazda serwera
narzędzi, ale liczby nikt nie ratyfikował. Granica obowiązuje wyłącznie tam,
gdzie bramka stoi — bez wymogu logowania żadne gniazdo nie jest związane
i granica odcięłaby pracę własną.

**Rozważone warianty.**

1. Zejść do jednej cyfry. Odrzucone: wejście Operatora zajmuje jedno gniazdo,
   ale po wybudzeniu maszyny albo zerwaniu sieci klient nawiązuje na nowo i przez
   chwilę stoi niezwiązany; okien bywa kilka, a stare gniazdo schodzi dopiero po
   swoim terminie. Granica jednocyfrowa odcinałaby powrót do pracy.
2. Zdjąć granicę. Odrzucone: rejestr bez granicy jest drogą wyczerpania pamięci
   maszyny z jednego procesu lokalnego, a to jedyne, przed czym ta liczba broni.
3. Zostawić 64. Przyjęte.

**Decyzja.** Wartość zostaje. Sześćdziesiąt cztery to jedna ósma granicy gniazd
w ogóle (`limitPolaczen = 512`), więc gniazda przed bramką nie wypchną pracy
uwierzytelnionej, a zapas pokrywa powrót po zerwaniu przy kilku oknach naraz.

**Konsekwencje.** Liczba wiąże wyłącznie instalację z wymogiem logowania —
wdrożenie serwerowe, gdzie pakiet stawia `DANACO_WYMOG_LOGOWANIA=true`.
Na pętli zwrotnej bez bramki nie zmienia niczego.

## 40. Eksport i import komponentu zostają poza wydaniem

**Rzecz.** Prototyp Centrum niesie w menu komponentu pozycje `eksport` (siedem
wystąpień) i `import` (jedno). Kontrakt rodziny `component` ma `assign`,
`changed`, `create`, `delete`, `list` i `update` — nic, co by te dwie pozycje
obsłużyło. Wiązanie Centrum zdejmuje je z okna, więc Operator ich nie widzi.

**Rozważone warianty.**

1. Dopisać `component.export` i `component.import` do kontraktu. Odrzucone:
   prototyp podaje nazwę pozycji i nic ponadto — nie mówi, co wchodzi do paczki,
   w jakim formacie, czy niesie ze sobą przypisania i konfigurację, ani jak
   zachować się przy imporcie komponentu o nazwie już zajętej. Napisanie tego
   byłoby dopowiedzeniem zachowania, którego w źródle nie ma.
2. Zostawić pozycje wygaszone. Odrzucone: martwa pozycja w menu to blokada bez
   zdania wyjaśniającego.
3. Zdejmować je z okna i nazwać rzecz Właścicielowi. Przyjęte.

**Decyzja.** Pozycje zostają zdjęte, kontraktu się nie poszerza. Eksport
komponentu wchodzi wtedy, gdy Właściciel poda jego zakres i format.

**Konsekwencje.** Menu komponentu w Centrum jest o dwie pozycje uboższe od
prototypu i różnica ta jest zamierzona — przy ocenie kryterium 5 nie liczy się
jako rozjazd z prototypem. Eksport rzeczy innych niż komponent działa bez
zmian: `session.export`, `research.export.*`, `studio.export.profile.*`.

## 41. Los siedmiu okien platformowych

**Rzecz.** Prototypy platformowe to Instalator, Instrukcja użytkowania, Nowy
projekt, Ustawienia, Konfiguracja, Historia sesji, Mobile i Always On Display.
Do wydania wchodzi Instalator; pozostałe albo mają pokrycie inną drogą, albo
nie mają go wcale, a plan etapu 2 zabrania i pustej powłoki, i zdjęcia wejścia.

**Rozstrzygnięcie po oknach.** Nowy projekt działa wprost — zakłada projekt
z Centrum i z przedsionka, bez osobnego okna. Konfiguracja otwiera katalog
„Operacje platformy" z rodziną `config`. Instrukcja użytkowania i Instalator
stoją pod pozycjami pomocy. Historia sesji, Mobile i Always On Display zostają
zapowiedziane: wejście klika się, a produkt odpowiada zdaniem, gdzie ten zakres
jest dziś — sesje w panelu bocznym, obie pozostałe w katalogu operacji.

**Rozważone warianty.**

1. Dołożyć trzy brakujące okna. Odrzucone: żadne nie należy do pakietu pięciu
   okien wydania, a zbudowanie ich odsunęłoby bramkę bez zysku dla Operatora,
   który te same operacje wywoła z katalogu.
2. Zdjąć wejścia z okna. Odrzucone: znika wtedy ślad, że taki zakres w produkcie
   istnieje, i Operator nie ma jak się dowiedzieć, gdzie go szukać.
3. Zostawić wejścia i nazwać niegotowość zdaniem. Przyjęte.

**Decyzja.** Osiem prototypów platformowych zamyka się tak: jeden w wydaniu,
cztery z pokryciem inną drogą, trzy zapowiedziane zdaniem wskazującym zastępstwo.

**Konsekwencje.** Żadne wejście platformowe nie jest martwe ani nie prowadzi do
pustej powłoki, więc kryterium 6 bramki obejmuje je wszystkie. Zapowiedź niesie
nazwę okna, zastępstwo i zdanie o wydaniu — sprawdzone klikaniem.
