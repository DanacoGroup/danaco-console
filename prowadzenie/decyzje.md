# Rejestr decyzji

Jedyne miejsce, w którym obowiązują rozstrzygnięcia trudno odwracalne. Wpis
zakłada Właściciel produktu. Rozstrzygnięcie nieobecne w tym rejestrze nie
obowiązuje, niezależnie od tego, gdzie zostało wypowiedziane.

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

## 8. Model wdrożenia: hybryda z dwoma wariantami instalacji

**Data:** 2026-08-23 · **Stan:** obowiązuje

**Kontekst.** Prototyp instalatora opisywał produkt jako program do pracy na
komputerze Operatora. To jest nieprawda i przekłamuje istotę wdrożenia.

**Decyzja.** Produkt jest **hybrydowy** i występuje w dwóch wariantach
instalacji:

| Wariant | Co staje u Operatora | Serwer Danaco |
|---|---|---|
| **cienki (hybryda)** | okno aplikacji wraz z interfejsem | **wymagany** — rdzeń i całe zaplecze stoją tam |
| **pełny natywny** | **wszystko**: okno, rdzeń, zaplecze, wagi modeli | **niepotrzebny** — produkt pracuje sam |

Rozróżnienie idzie po **zależności od serwera Danaco**, nie po miejscu obliczeń.
Natywna nie może niczego dobierać z serwera Danaco po instalacji, więc instalka
niesie komplet tego, po co rdzeń sięga.

**Zasady korzystania są jedne i uniwersalne.** Nie rozdziela się ich na dwie
wersje wedle wariantu instalacji i nie opisuje się produktu tak, jakby istniał
tylko jeden z nich. Treść, która zakłada pracę wyłącznie lokalną, jest błędna;
treść zakładająca wyłącznie serwerową — również.

**Zakres składników do zainstalowania wynika z wybranego wariantu.** Wykaz
przedstawiany w instalatorze nie jest stały: wariant cienki nie niesie rdzenia
ani jego zaplecza, więc nie może ich proponować. Krok wyboru składników zależy
od wariantu wybranego wcześniej.

**Konsekwencje.** Instalator ma dwie ścieżki, nie jedną. Opis produktu
w zasadach korzystania obejmuje oba warianty bez wskazywania któregokolwiek jako
domyślnego sposobu pracy. Skrypt składający wariant cienki
(`instalka-hybryda-win-x64.sh` wraz z odmianami dla ARM i Linuksa) **nie leży
w repozytorium budowy** — stoi w materiale zabezpieczonym poza gitem,
`~/robocze/material/repo-2.0/budowa/scripts/`. Jest zgodny z decyzją: wkompilowuje
`client/dist` w powłokę i sprawdza wykazem zawartości, że rdzenia w instalce nie
ma. Wchodzi do budowy wraz z terenem instalatora.

---

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

**Skutek dla okna.** Przepływ wejścia ma dwie gałęzie, nie jedną.

| Nadajnik | Odpowiedź rejestracji | Co okno mówi |
|---|---|---|
| ustawiony | `PendingVerification: true` | list wysłany, przepisz drogę potwierdzenia |
| brak | `PendingVerification: false` | konto założone, adres niepotwierdzony — wejście działa hasłem, **odzyskanie konta listem nie zadziała do chwili potwierdzenia** |

Ostrzeżenie z drugiego wiersza jest wymagane, nie zalecane. Adres jest jedyną
drogą odzyskania konta; Operator, który nie wie, że jego adres nie został
sprawdzony, dowie się o tym w dniu, w którym będzie go potrzebował.


---

## Pozycje otwarte

Pozycja otwarta czeka na rozstrzygnięcie Właściciela i blokuje wskazany etap.

| Pozycja | Blokuje | Opis |
|---|---|---|
| Nazewnictwo metaforyczne w zachowanych plikach wykonawczych | etap 2 | Pliki w `narzedzia/` posługują się określeniami metaforycznymi zamiast terminologii zawodowej. Terminologia docelowa powstaje w etapie 2 i wtedy pliki wymagają uzgodnienia. |
| Skład dokumentacji przekroju pionowego | etap 2 | Zakres opracowań powstających wraz z przekrojem — wyłącznie to, co przekrój obsługuje, nie pełny zbiór pięćdziesięciu jeden. |
| Odłączenie karty do osobnego okna | etap 1 | Propozycja Właściciela: karta otwierana w drugim oknie aplikacji, poza oknem głównym — użyteczne na szerokim monitorze. Wymaga własnej nazwy, bo „okno robocze" jest już zajęte przez sesję. |
| Grupowanie kart w paśmie | etap 1 | Pasmo ma dwa ułożenia: poziome oraz karty zgrupowane w jedną, rozwijalną. Reguły grupowania — co wolno zgrupować i kto o tym rozstrzyga — wymagają domknięcia. |
| Okno boczne: samouczek i karta naraz | etap 1 | Okno boczne niesie samouczek, a także przyjmuje wskazaną kartę. Do rozstrzygnięcia, czy mieszczą się tam obie rzeczy jednocześnie, czy karta zastępuje samouczek na czas pobytu. |
| Które moduły tworzą sesję | etap 2 | Sesja powstaje przy interakcji z modelem. Moduły wytwórcze — agenci, automatyzacje — jej nie tworzą. Podział wszystkich modułów na te dwa rodzaje wymaga rozstrzygnięcia, bo wyznacza zawartość wykazu sesji. |
| Izolacja okna bez sesji | etap 2 | Okno robocze modułu wytwórczego nie ma sesji, a mimo to powstają w nim pliki. Czy podlega izolacji, a jeśli tak — czego dotyczy. |
| Zamknięcie ostatniej karty czatu | etap 2 | Rozstrzygnięte jest, że zamknięcie okna nie kończy sesji. Pozostaje pytanie węższe: czy zamknięcie w oknie ostatniej karty z interakcją zmienia cokolwiek w stanie sesji. |
| Zakres trybów `manual`, `auto`, `plan` | etap 2 | Czy `auto` obejmuje tworzenie i usuwanie plików obok edycji; czy `plan` dopuszcza odczyt, bez którego nie ma z czego planować. Zakres `bypass permissions` jest rozstrzygnięty: zniesienie zupełne. |
| Warunek ukończenia zadania | etap 2 | Nakaz pracy do wyniku wymaga sprawdzalnej odpowiedzi na pytanie, kiedy zadanie jest ukończone. Bez niej praca ciągła albo kończy się przedwcześnie, albo nie kończy wcale. |
| Sygnał ukończenia w kontrakcie | etap 2 | Ukończenie wyzwala czynność następną, więc „ukończone z wynikiem" i „przerwane" muszą być rozróżnialne maszynowo, nie z treści odpowiedzi. To pozycja kontraktu, nie stan interfejsu. |
| Trwałość pętli | etap 2 | Praca licząca się w miesiącach przetrwa restart maszyny, awarię i zamknięcie aplikacji albo nie przetrwa wcale. Stan pętli musi być odtwarzalny z zapisu, nie z pamięci procesu. |
| Zachowanie przy wyczerpaniu budżetu | etap 2 | Praca ciągła konsumuje zasób bez nadzoru. Do rozstrzygnięcia, co dzieje się z zadaniem, gdy zasób kończy się w połowie. Rozstrzygnięcie nie może przywracać pytania do Operatora — wyczerpanie zasobu jest zatrzymaniem z braku środków, nie prośbą o zgodę. |
| Przerwanie pracy — natychmiastowość | etap 2 | Komponent przerwania w oknie czatu ma zatrzymywać pracę natychmiast, nie po zakończeniu bieżącej czynności. To wymaganie wobec architektury wykonania; komponent sam w sobie należy do okna czatu i jest niezależny od uprawnień. |
| Dziennik pracy bez nadzoru | etap 2 | Praca prowadzona bez człowieka wymaga zapisu przebiegu, bo inaczej nie da się odtworzyć, dlaczego zadanie potoczyło się tak, a nie inaczej. |
| Nazwy poziomów w bibliotece | etap 2 | Uzgodnienie `.dn-karty-sesji` i `.dn-karta-widoku` z modelem czterech poziomów, razem z przebudową mechanizmu kart. |
| Wskaźnik izolacji — poziom przypisania | etap 2 | `izolacja.css` stanowi, że izolacja jest cechą okna roboczego; wskaźnik stoi w ramie aplikacji. |
| Wartość domyślna `createVersion` | etap 2 | Kontrakt nie ustala, co znaczy brak pola w `studio.document.save`, choć w trzech komendach siostrzanych mówi wprost „brak znaczy tak". Rdzeń stosuje oba idiomy niejednolicie. Do rozstrzygnięcia wraz z tym, czy reguła obejmuje `document.save.as` i `document.form.save`. |
| Potwierdzenie adresu po instalacji bez poczty | etap 2 | Po rejestracji bez nadajnika adres zostaje niepotwierdzony **na zawsze**: `auth.verify` przyjmuje wyłącznie token, `auth.register` odmawia kodem `conflict`, a droga z `auth.recover` jest wydana do innej czynności. Żadna z dziewięciu komend `auth.*` nie wydaje drogi potwierdzenia po raz drugi — choć opis `auth.verify` w kontrakcie stanowi, że wygasła droga „pozwala poprosić o nową". Skutek: adres, który pozycja 11 nazywa jedyną drogą odzyskania konta, nie daje się sprawdzić, a znacznik `bramka-bez-poczty` nie ma osiągalnej drogi zdjęcia. Naprawa wymaga nowej komendy, czyli zmiany kontraktu. |
| Sprzeczne opisy odzyskiwania w kontrakcie | etap 2 | `auth.password.reset` stanowi, że „drogi odzyskania listem na adres e-mail dziś NIE MA — rdzeń poczty nie wysyła", a `auth.recover` w tym samym kontrakcie mówi, że serwer wysyła drogę potwierdzenia na wskazany adres. Dwa opisy przeczą sobie; jeden z nich kieruje wykonawcę przeciwko zachowaniu rdzenia. |
| Metryki w planszach i indeksie designu | etap 1 | Podają wartości rozjechane ze stanem plików, a metody liczenia „interakcji" nie da się odtworzyć. Albo dostają definicję, albo znikają. |
