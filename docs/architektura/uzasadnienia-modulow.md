# Uzasadnienia modułów

Dokument gromadzi uzasadnienia decyzji projektowych, które nie mieszczą się
w nagłówku pliku źródłowego. Rozdział nosi nazwę pliku, którego dotyczy.

## budowa/server/internal/core/adapter_modul_extension_integracje.go

Rozmowa z serwerem MCP przy odkrywaniu narzędzi, próbnym wywołaniu i sprawdzeniu
kondycji jest rzeczywista: prowadzi ją klient MCP wkompilowany w rdzeń, po
transporcie ustawionym dla pozycji. Każda ramka rozmowy trafia do dziennika
protokołu, a każde wywołanie narzędzia — do dziennika użycia; metryki i audyt
liczą się z tych wierszy, nie z licznika w pamięci procesu.

Import definicji API czyta opis zamiast zgadywać operacje. Opis OpenAPI 3
rozkłada biblioteka `getkin/kin-openapi`, wkompilowana w binarium; schemat
GraphQL rozkłada czytnik własny, wyodrębniający pola typów `Query`
i `Mutation`. Żadna z dwóch dróg nie woła programu zewnętrznego.

Piaskownica uruchamia wyłącznie to, co integracja ma podłączone: pozycja
z transportem procesu lokalnego dostaje wejście na standardowe wejście procesu
i oddaje jego wyjście; pozycja rozmawiająca protokołem MCP dostaje wywołanie
`tools/call`. Pozycja bez jednego i drugiego wraca odmową z powodem, a nie
wynikiem udającym przebieg.

Odwołanie do poświadczenia integracji wchodzi także do rejestru referencji
sekretów rozszerzeń, ponieważ centrum uprawnień pyta o nie osobno — referencja
znana wyłącznie jednej integracji byłaby niewidoczna dla rejestru rotacji
sekretów.

Adres zgody OAuth2 powstaje wyłącznie wtedy, gdy integracja ma ustawiony adres,
do którego można pokierować przeglądarkę — adres wymyślony prowadziłby
donikąd.

Polecenie uruchamiane w piaskownicy pochodzi wprost od integracji podłączonej
rdzeniowi jako rozszerzenie; rdzeń nie składa go z danych zewnętrznych, stąd
wyłączenie ostrzeżenia gosec G204.

Log przebiegu piaskownicy ląduje w magazynie treści rdzenia, ponieważ kontrakt
oddaje do niego wyłącznie odwołanie, za którym musi leżeć zapisany plik.

Pole `output` odpowiedzi piaskownicy przyjmuje wyjście wyłącznie wtedy, gdy
jest poprawnym JSON-em, zgodnie z typem `json` zadeklarowanym w kontrakcie.

Operacje zaimportowane z opisu API lądują w tym samym wykazie narzędzi
pozycji, który wypełnia odkrywanie MCP — konektor zbudowany z opisu API jest
z punktu widzenia eksperta tym samym co narzędzie odkryte protokołem.

Adres nasłuchu webhooka przychodzącego składa rdzeń, nie integracja
zewnętrzna, ponieważ tylko rdzeń zna ścieżkę, pod którą odbiera zdarzenia.

Okno domyślne metryk użycia obejmuje ostatnią dobę, aby uniknąć sumowania
wywołań od początku istnienia pozycji.

Sprawdzenie kondycji bez wskazanego kodu pozycji obejmuje wyłącznie pozycje
włączone, aby uniknąć ruchu sieciowego do integracji wyłączonych.

Stan `degraded` sprawdzenia kondycji oznacza serwer, który odpowiedział na
powitanie, lecz nie oddał wykazu narzędzi — stan pośredni między dostępnością
a awarią.

Uprawnienia nadane pozycji rozszerzenia audyt czyta raz na pozycję, nie raz na
wpis, ponieważ audyt jednej pozycji obejmuje zwykle setki wywołań.

## budowa/server/internal/core/adapter_modul_developer_jezyk.go

Wiedzę o kodzie repozytorium niesie program serwera, nie logika rdzenia: gopls
zna typy repozytorium i graf odwołań, gofmt i prettier znają styl,
golangci-lint zna reguły analizy. Wołanie tych programów jest jedynym
miejscem tego modułu, które sięga po proces serwera — odtworzenie tej wiedzy
w kodzie rdzenia oznaczałoby napisanie drugiego kompilatora. Wszystkie stoją
na serwerze i wchodzą do sondy zależności.

Brak programu na serwerze nie zamienia odpowiedzi w odmowę. Kontrakt niesie
pola `serverAvailable` i `linterAvailable` właśnie po to: pusty wykaz symboli
przy `serverAvailable: false` znaczy brak narzędzia do sprawdzenia, a przy
`true` — sprawdzenie bez wyniku. Pierwszy stan naprawia instalacja programu na
serwerze, drugi — poprawka w kodzie.

Serwer języka jest wołany w trybie jednorazowym na każde pytanie, a nie jako
utrzymywana sesja: komenda kontraktu jest pytaniem i odpowiedzią bez stanu
między nimi, a proces utrzymywany na czas okna wymagałby własnego rejestru
uchwytów, sprzątania i przerywania — trzeciego takiego mechanizmu w module
obok budowania i debugowania.

Serwer języka TypeScriptu nie ma trybu jednorazowego wywołania, jaki ma gopls
— rozmawia wyłącznie sesją protokołu LSP na strumieniu wejścia: uzgodnienie,
otwarcie dokumentu, pytanie, zamknięcie. Jedyna droga rdzenia do procesu
zewnętrznego z zamysłu nie pisze na jego wejście, ponieważ jednoczesne pisanie
i czytanie tego samego procesu jest klasą zakleszczenia, której ten pakiet ma
nie mieć. Dlatego pliki TypeScriptu dostają odpowiedź nazywającą brak serwera
zamiast wyniku z serwera Go, który tego pliku nie rozumie; program pozostaje
zadeklarowany w wykazie zależności, więc sonda startowa zgłasza jego brak.

Formater dobiera się po rozszerzeniu pliku, bo styl należy do języka. Treść
z żądania wygrywa z treścią na dysku, ponieważ edytor formatuje bufor,
którego jeszcze nie zapisano, i to on jest przedmiotem czynności. Dla plików
Go pierwszeństwo ma goimports: robi to samo co gofmt, a dodatkowo porządkuje
wykaz importów, którego ręczne pilnowanie jest najczęstszym powodem
niekompilującego się pliku po refaktoryzacji; gdy goimports nie ma na
serwerze, formatuje gofmt.

Plik roboczy formatowania powstaje obok pliku źródłowego, nie w katalogu
tymczasowym systemu, ponieważ reguły formatowania zależą od położenia
w repozytorium — plik przeniesiony gdzie indziej sformatowałby się wedle
innych reguł niż ten, którego dotyczy czynność.

Zawężenie wyniku do zakresu bierze z formatowanej treści wyłącznie wskazane
wiersze, ponieważ formatery pracują na całym pliku, a kontrakt dopuszcza
zakres — reszta pliku zostaje nietknięta, bo zaznaczenie dotyczyło zakresu,
nie całego pliku.

Repozytorium bywa wielojęzyczne, a każdy analizator zna jeden język:
golangci-lint czyta Go, Ruff — Pythona, Stylelint — arkusze CSS, typos szuka
literówek niezależnie od języka. Zgłoszenia wszystkich programów idą do
jednego wykazu, ponieważ kontrakt niesie w każdym wpisie pole `source`.
Pole `linterAvailable` pochodzi z czasu, gdy program analizy był jeden,
i znaczy: choć jeden program odpowiedział. Fałsz tego pola jest więc stanem,
w którym pusty wykaz zgłoszeń kłamałby, bo nie sprawdzono niczym; które
programy milczały, ta odpowiedź nie rozróżnia bez zmiany kontraktu.

Dobór programów analizy idzie po rozszerzeniu plików repozytorium, nie po
tym, co stoi na maszynie, ponieważ program uruchomiony tam, gdzie nie ma ani
jednego pliku jego języka, nie milczy, tylko odmawia — golangci-lint
w repozytorium bez plików Go kończy się błędem o braku plików do analizy.
Odmowa jednego programu nie ma prawa zabrać analizy pozostałym: repozytorium
samego Pythona albo samego frontendu zostałoby wtedy bez analizy w ogóle.
Żądanie bez wskazania plików obejmuje całe drzewo i każdy program dostaje
swój katalog.

Nonzero exit code programów analizy jest tu wynikiem, nie usterką: Ruff,
Stylelint i typos kończą się niezerowo dokładnie wtedy, gdy mają co zgłosić.
Miarą przeprowadzenia analizy jest odczytana odpowiedź, a nie kod wyjścia ani
strumień, na który program ją napisał — program, który przewrócił się bez
zwrócenia czytelnego wyniku, analizy nie przeprowadził, i jego milczenie nie
ma prawa wyglądać jak brak zastrzeżeń.

Waga zgłoszenia bez wagi nadanej przez linter jest ostrzeżeniem, nie błędem:
literówka w identyfikatorze bywa nazwą celowo skróconą, a program nie ma jak
tego rozstrzygnąć — podniesienie wagi kazałoby poprawiać rzeczy, których
nikt tak nie oznaczył.

Podgląd jest domyślnym trybem refaktoryzacji. Refaktoryzacja semantyczna
dotyka wielu plików naraz, więc zmiana zapisuje się na dysk wyłącznie na
wyraźne żądanie z wyłączonym podglądem — w przeciwnym razie wraca różnica do
zaakceptowania.

Różnica zunifikowana jest jedynym kształtem, w którym serwer języka mówi
o zmianie wielu plików naraz; kontrakt niesie ją jako wykaz zmian tekstu,
więc rozbiór wyniku idzie po nagłówkach plików i fragmentach zmiany.

## budowa/server/internal/core/adapter_modul_developer_zaleznosci.go

Wykaz zależności powstaje z bezpośredniego odczytu manifestów repozytorium
(go.mod, package.json, requirements.txt, Cargo.toml), nie z wywołania
narzędzia języka: odczyt pliku tekstowego o ustalonym kształcie jest
kilkudziesięcioma wierszami kodu i działa bez narzędzia zainstalowanego na
serwerze oraz bez pobranych zależności — przeglądarka zależności otwiera się
także w repozytorium, którego jeszcze nikt nie zbudował.

Skan sekretów rozpoznaje kształty kluczy regułami wyrażeń regularnych, tym
samym środkiem co wyszukiwanie w repozytorium i tą samą metodą co program
gitleaks. Skan kodu szuka wzorców znanych źródeł podatności w Go
i TypeScripcie — zakres węższy niż semgrep, nazwany wprost: reguła nazywa
się w znalezisku, więc widać, co sprawdzono.

Podatności zależności rozpoznaje się po zakresie wersji, nie po bazie CVE:
serwer bywa odcięty od sieci, a skan milczący bez sieci byłby skanem
twierdzącym „nic nie znaleziono" tam, gdzie nie szukał. Skan zależności
zgłasza więc wyłącznie to, co da się orzec z samego manifestu: zależność
bez przypiętej wersji i zależność wskazującą gałąź zamiast wydania.

Zawężenie „tylko przestarzałe" oddaje pozycje bez przypiętej wersji, ponieważ
rdzeń nie pyta rejestrów pakietów o najnowsze wydania i nie udaje, że zna
wersję nowszą.

Kolejność wykrywania manifestu jest ustalona i celowa: go.mod przed
package.json, bo repozytorium Go z narzędziami frontendowymi ma oba pliki,
a jego zależnościami są moduły Go. Wskazanie jawne w żądaniu zawsze wygrywa
z tym rozpoznaniem.

Zależności deweloperskie i towarzyszące manifestu Node stoją pod własnym
rodzicem: są zależnościami repozytorium, lecz nie wchodzą do wydania,
a wykaz płaski zacierałby tę różnicę.

Reguły sekretów nazywają się tak, jak nazywa się to, co znajdują: znalezisko
niosące „klucz prywatny" mówi więcej niż numer reguły. Wzorce są celowo
wąskie — reguła łapiąca każde słowo `password` dałaby wykaz, w którym
prawdziwe znaleziska toną wśród nazw pól formularzy.

Skan biegnie w całości przed odpowiedzią, a nie w tle: przechodzi po plikach
repozytorium wyrażeniami regularnymi i po manifestach, więc kończy się
w sekundach. Bieg w tle wymagałby własnego rejestru przebiegów i własnego
zdarzenia przyrostu — trzeciego takiego mechanizmu w module, bez odbiorcy,
który by na to czekał.

Przejście po plikach przy skanie treści idzie tą samą drogą, co wyszukiwanie
w repozytorium: z poszanowaniem `.gitignore` i z pominięciem plików
binarnych. Skan zgłaszający sekret w katalogu `node_modules` byłby wykazem,
którego nikt nie czyta.

Treść dopasowania sekretu nie wchodzi do znaleziska: wykaz znalezisk
z wypisanymi sekretami byłby drugim miejscem, w którym te sekrety leżą —
tym razem w bazie produktu.

Reguły wyrażeń regularnych są drogą podstawową skanu kodu: pracują zawsze,
bez niczego spoza rdzenia, i to one rozstrzygają, że skan kodu w ogóle coś
zmierzył. Programy zewnętrzne tę drogę poszerzają: semgrep orzeka po
składni tam, gdzie wyrażenie regularne widzi tylko wiersz, a jscpd i dupl
znajdują powtórzony fragment, którego żadna reguła wierszowa nie zobaczy.
Program nieobecny na maszynie nie psuje skanu i nie zmienia jego stanu —
skan oddaje wtedy to, co zmierzyła droga podstawowa; każde znalezisko niesie
nazwę reguły, więc widać, co je wystawiło.

Próg powtórzenia kodu rdzeń nie narzuca: każdy z programów jscpd i dupl ma
własny próg i to on obowiązuje. Próg mówi, od ilu żetonów zbieżność jest
powtórzeniem, a nie przypadkiem — jest rozstrzygnięciem o tym, co w danym
języku jest powieleniem kodu; wartość wpisana w rdzeniu byłaby wzięta
znikąd i cichym nadpisaniem tego, co program o swoim języku wie.

Skan semantyczny Semgrepa używa wyłącznie zestawu reguł repozytorium, nie
zestawu z rejestru, ponieważ zestaw rejestru pobiera się przez sieć przy
każdym przebiegu — ten sam powód, dla którego skan zależności nie pyta bazy
CVE. Reguły są też rozstrzygnięciem o tym, co w danym repozytorium jest
podatnością, a tego rdzeń nie orzeka samodzielnie. Repozytorium bez
własnego zestawu reguł nie jest więc skanowane semantycznie, i żadne
znalezisko nie twierdzi inaczej.

Raport jscpd idzie do katalogu tymczasowego maszyny rdzenia, bo program
pisze go do pliku, nie na wyjście; katalog znika po odczycie, ponieważ
repozytorium nie ma prawa dostać pliku, o który nikt nie prosił. Zawężenie
tego skanu do TypeScriptu i JavaScriptu jest podziałem pracy, nie
oszczędnością: powtórzenia w plikach Go liczy program dupl osobno, a oba
programy puszczone na ten sam plik zgłosiłyby to samo powtórzenie dwa razy.
Brak raportu — czy to dlatego, że program go nie stworzył, czy dlatego, że
nie ma go na maszynie — pozostawia wynik drogi podstawowej bez zmian.

Para powtórzeń z programu dupl wraca dwukrotnie, raz z każdej strony, więc
zapisywana jest jedna strona pary: wykaz z tym samym fragmentem dwa razy
mówiłby o dwóch spostrzeżeniach tam, gdzie jest jedno.

## budowa/server/internal/core/adapter_modul_extension_katalog.go

Wyszukiwarka szuka wyłącznie w tym, co pozycja katalogu niesie: nazwa, opis,
kategoria, udostępniane narzędzia i znaczniki. Rdzeń przeszukuje więc kod,
nazwę i opis wiersza oraz nazwy narzędzi odkrytych u integracji — wszystko,
co naprawdę leży w bazie. Podpowiedzi składają się z nazw pozycji i narzędzi,
które trafienie zawierają, a nie z listy wpisanej w kodzie.

Prywatny rejestr organizacji nie jest drugim katalogiem: wykaz rejestru
oddaje te pozycje katalogu, które powstały z publikacji pakietu — publikacja
zostawia w konfiguracji pozycji odwołanie do pakietu i jego archiwum. Osobna
tabela rejestru byłaby drugą prawdą o tym samym zbiorze pozycji. Adres
rejestru jest lokalny, bo rejestr jest lokalny: pozycje leżą w bazie tego
rdzenia, a archiwa w jego magazynie treści; adres wskazujący cudzy serwer
byłby obietnicą, za którą nic nie stoi.

Sprawdzenie aktualizacji liczy się z wersji, które istnieją: zestawia wersję
zainstalowaną z najwyższą wersją zapisaną w tabeli wersji rozszerzenia, a te
wiersze powstają przy publikacji pakietu i przy przesłaniu paczki. Rdzeń nie
pyta o aktualizacje żadnego serwera w sieci — nie ma dokąd pytać, a udawanie
odpowiedzi byłoby meldunkiem bez pokrycia.

Żadna z komend tej rodziny niczego nie blokuje: operacje zbiorcze, cofnięcie
wersji i instalacja zestawu wykonują się od razu, a pozycje, których wykonać
się nie dało, wracają w polu odrzuceń z powodem, nie jako przerwany przebieg.

Odsunięcie i granica strony wyszukiwarki liczą się po zawężeniu wyników do
zapytania: pole całkowitej liczby opisuje zbiór spełniający warunki, nie
długość oddanej strony.

Dziennik zmian karty szczegółów składa się z wpisów wersji niosących opis —
jedyne miejsce, w którym rdzeń go trzyma; wersja bez wpisu nie dokłada
pustego wiersza do dziennika.

Zależności i znaczniki karty szczegółów niesie konfiguracja pozycji, nie
osobne pola żądania, ponieważ manifest pakietu odkłada je właśnie tam.

Suma kontrolna podana przy przesyłce paczki jest sprawdzana, nie przyjmowana
bezkrytycznie: przesyłka, która dojechała uszkodzona, ma się o tym dowiedzieć
przy przesłaniu, a nie dopiero przy instalacji.

Przypięcie wersji dopuszcza wyłącznie wersję, która istnieje w dzienniku
pozycji: przypięcie do numeru wymyślonego byłoby obietnicą, której nikt nie
spełni.

## budowa/server/internal/core/adapter_modul_badania_zrodla.go

Skutkiem każdej czynności tej rodziny jest byt, nie odpowiedź: wiersz w
bazie, wiązanie zdjęte albo bajty w magazynie — mierzalny stan po jej
zakończeniu, niezależnie od treści odpowiedzi. Odpowiedź powodzenia z pustym
wynikiem jest wzorcem szkody, dlatego przechwycenie strony zapisuje treść,
nie sam wiersz o niej.

Import bibliografii czyta formaty BibTeX, RIS, CSL-JSON, EndNote XML i CSV
parserami napisanymi w tym module, bibliotecznymi środkami języka, bez
wołania programu zewnętrznego — import ma działać niezależnie od tego, co
jest doinstalowane.

Całkowita liczba pozycji wykazu źródeł liczy się przed wycięciem strony:
opisuje zbiór spełniający zawężenie żądania, nie długość oddanej strony,
inaczej licznik zmieniałby się przy przewijaniu.

Podstawa dopasowania duplikatów źródeł jest nazywana wprost w odpowiedzi:
identyczny identyfikator jest pewnością, zbieżny adres prawie pewnością,
podobny tytuł wyłącznie podpowiedzią. Scalenie zostaje decyzją zewnętrzną —
rdzeń nie scala pozycji samodzielnie.

Liczba odcinków transkrypcji liczy się z akapitów transkryptu, ponieważ
kontrakt pyta o liczbę odcinków, nie o liczbę znaków, i wartość musi być
policzona, nie wpisana na sztywno.

## budowa/server/internal/core/adapter_modul_aplikacje_pakiety.go

Pakiet jest archiwum na dysku: budowanie pakietu bierze artefakt budowania
(wskazany albo ostatni z udanego wdrożenia), rozpakowuje go i składa nowe
archiwum wraz z manifestem, po czym kładzie je w magazynie treści rdzenia.
Wiersz pakietu wskazuje ten plik, jego rozmiar i sumę kontrolną — pakiet bez
zapisanych bajtów byłby wzorcem szkody, którego pilnują sprawdziany skutku.

Archiwa zip i tar.gz składają biblioteki wkompilowane w binarium (pakiety
`archive/zip`, `archive/tar`, `compress/gzip` biblioteki standardowej), bez
wołania programu zewnętrznego — pakowanie działa na instalacji niosącej sam
rdzeń.

Podpis pakietu jest prawdziwym podpisem Ed25519: pakiet `crypto/ed25519`
biblioteki standardowej podpisuje sumę kontrolną archiwum, a podpis jest od
razu weryfikowany kluczem publicznym, więc pole potwierdzenia mówi
o sprawdzeniu, które naprawdę przeszło, nie o samym zamiarze podpisania.

Klucz wydawcy nie przechodzi przez kontrakt: żądanie niesie wyłącznie
odwołanie do klucza w warstwie sekretów. Materiał klucza leży w sejfie
poświadczeń rdzenia; odwołanie użyte po raz pierwszy zakłada tam nowy klucz
wydawcy, ponieważ inaczej nie dałoby się podpisać pierwszego pakietu bez
żądania treści klucza wprost, co wniosłoby sekret do kontraktu i do
dziennika.

Publikacja nie zakłada drugiego rejestru: prywatny rejestr organizacji to ta
sama tabela, którą prowadzi rodzina komend rozszerzeń — pozycja opublikowana
trafia do tego samego katalogu obok pozycji zainstalowanych wprost. Osobny
rejestr obok tamtego byłby drugą prawdą o tym samym katalogu.

Manifest pakietu przy budowaniu zaczyna od tożsamości, którą rdzeń naprawdę
zna — kodu pakietu i nazwy produktu okna; resztę pól wypełnia dopiero zapis
manifestu osobną komendą, ponieważ wymyślenie ich przy budowaniu byłoby
wpisaniem treści, której nikt nie podał.

Weryfikacja podpisu następuje natychmiast po jego złożeniu, aby pole
potwierdzenia mówiło o sprawdzeniu, które przeszło, a nie o samym fakcie, że
podpis powstał. Sam podpis idzie do wiersza wraz z jego postacią bajtową —
bez niej ponowna weryfikacja nie miałaby czego sprawdzić przy pozycji, która
z tego pakietu powstanie.

Konfiguracja pozycji katalogu powstałej z publikacji niesie wyłącznie to, co
odróżnia ją od pozycji instalowanej ręcznie: pakiet, z którego powstała,
jego archiwum i widoczność nadaną przy publikacji — bez pola wymyślonego.
Podpis pakietu przechodzi na pozycję katalogu wraz z materiałem do jego
ponownego sprawdzenia, ponieważ przepisanie samego werdyktu nie byłoby
weryfikacją.

Pozycja katalogu założona z publikacji wchodzi wyłączona, tak samo jak
każda pozycja spoza zestawu wbudowanego: włącza ją świadoma decyzja po
przejrzeniu uprawnień, nie sam fakt publikacji w organizacji.

Wersja manifestu wchodzi do rejestru wersji pozycji zawsze, także przy
pakiecie niepodpisanym — bez tego wpisu cofnięcie wersji nie miałoby dokąd
wrócić po drugim wydaniu, a sprawdzenie aktualizacji nie miałoby czego
z czym zestawić. Pakiet niepodpisany po prostu nie zostawia podpisu na
pozycji, i to jest prawda, którą walidator manifestu potem pokazuje.

Kolejność wpisów archiwum jest ustalona alfabetycznie, aby archiwum tej
samej treści miało tę samą sumę kontrolną — mapa Go przechodzi się
w kolejności losowej.

Uprawnienie zadeklarowane w manifeście bez wskazania bytu jest uprawnieniem
na wszystko: uprawnienie sieciowe bez domeny albo zapis plików bez korzenia
katalogu. Walidator manifestu ma o tym ostrzec jako sygnał, nie jako bramę
wstrzymującą publikację.

## budowa/server/internal/core/adapter_modul_library_archiwum.go

Trzy postacie utrwalenia archiwalnego są trzema odrębnymi standardami
branżowymi, nie wariantami jednego zapisu, więc każda ma własną drogę
przez rdzeń. PDF/A przechodzi przez bibliotekę `pdfcpu`: struktura dokumentu
zostaje uporządkowana, a wynik zwalidowany tą samą biblioteką; zapis
walidacji mówi wprost, co sprawdzono — rdzeń nie niesie walidatora profilu
PDF/A-2b i zgodności z profilem nie orzeka, bo orzeczenie bez sprawdzenia
byłoby pieczątką, nie utrwaleniem. BagIt składa pakiet w postaci Biblioteki
Kongresu (`bagit.txt`, `bag-info.txt`, manifest sum kontrolnych i katalog
`data/`), a walidacja polega na przeliczeniu manifestu z bajtów spakowanych —
pakiet, którego manifest się nie zgadza, jest wynikiem niepoprawnym i tak
wchodzi do odpowiedzi. PREMIS/METS składa dokument XML opisujący obiekt,
jego sumę kontrolną i zdarzenie utrwalenia; wynikiem jest opis, nie kopia
zasobu, więc pole zasobu wytworzonego wskazuje osobny plik opisu.

Pakowanie archiwum powstaje w całości kodem wkompilowanym w binarium —
pakiety standardowe archiwizacji, kompresji i sum kontrolnych biblioteki Go
wraz z biblioteką `pdfcpu`. Format 7z jest jedynym wyjątkiem: sięga po
arsenał serwerowy programem `7z`, ponieważ biblioteka standardowa tego
formatu nie zapisuje, i tylko wtedy, gdy żądanie wprost o ten format poprosi;
brak arsenału jest wtedy odmową nazwaną wprost, z podaniem formatów, które
rdzeń składa bez programu zewnętrznego.

Utrwalenie w miejsce oryginału jest dołożeniem nowej wersji zasobu, nie
nadpisaniem: poprzednia treść zostaje w historii wersji, więc jest dokąd
wrócić, gdy normalizacja PDF/A coś w dokumencie przestawi.

Walidacja pakietu BagIt idzie po bajtach już spakowanych, nie po tych,
z których pakiet powstał — pakiet, którego nie da się odczytać, jest
pakietem nieważnym, choćby materiał źródłowy był w porządku.

Liczba zasobów objętych zapisaną polityką retencji jest liczbą rzeczywistą,
przeliczoną z repozytorium, nie zapowiedzią: polityka zasięgu projektowego
obejmuje zasoby tego projektu, polityka zasięgu globalnego — wszystkie
zasoby czynne. Termin retencji każdego zasobu liczy się w chwili pytania,
bo wynika z jego ostatniej zmiany i okresu polityki, a obie wartości
zmieniają się w czasie.

Migawka wywozu paczki jest zapisem stanu całego repozytorium — zawężenia
żądania do zbioru plików albo kolekcji przestają wtedy obowiązywać, bo
migawka części repozytorium nie byłaby migawką. Zasób bez odczytanej treści
nie wywraca wywozu paczki: wchodzi do indeksu jako pozycja bez bajtów, żeby
paczka mówiła prawdę o stanie repozytorium, zamiast pomijać milczeniem to,
czego nie udało się odczytać.

## budowa/server/internal/core/adapter_modul_tlumaczenie_korekta_silniki.go

Trzy silniki zewnętrzne korekty uzupełniają reguły wbudowane komendy, które
są mechaniczne i przez to sprawdzalne: to samo wywołanie na tej samej treści
daje te same ustalenia o tych samych identyfikatorach, więc zastosowanie
poprawki ma co adresować. Wywołanie modelu tej własności nie ma i dlatego
w tym miejscu nie wchodzi. LanguageTool prowadzi całość — gramatykę,
ortografię, interpunkcję, typografię i styl, każdą regułą nazwaną i z
propozycją poprawki — i jest drogą pierwszą, bo jako jedyny widzi zdanie,
nie samo słowo. Hunspell prowadzi samą ortografię i wchodzi wyłącznie wtedy,
gdy LanguageToola nie ma: uruchomione razem podwoiłyby każdą literówkę,
ponieważ ortografię LanguageToola liczy ta sama rodzina słowników
morfologicznych. Jest to ten sam układ pierwszeństwa, co program `piper`
przed `espeak-ng` przy syntezie mowy tego modułu. Vale prowadzi styl prozy —
powtórzenia i terminy — obok LanguageToola, nie zamiast niego: mierzy tekst
jako całość, a nie zdanie, i nie ma z tamtym wspólnych reguł.

Wybór padł na hunspella, nie na `enchant-2`: `enchant-2` jest pośrednikiem,
nie słownikiem — na maszynie budowy wypisuje trzech dostawców i dla polskiej
treści oddaje wynik znak w znak taki sam jak hunspell wołany wprost, bo woła
właśnie jego. Pośrednik dokłada warstwę, która nie mówi, który dostawca
odpowiedział, oraz zestaw modułów wtyczkowych do spakowania obok programu.
Hunspell jest jednym plikiem wykonywalnym z plikami słownika obok, czyli
układem, który pakowanie produktu niesie wprost, a jego słowniki są tymi
samymi plikami, które ma już LibreOffice stojący w wykazie zależności.

LanguageTool umie chodzić jako usługa długożyjąca, ale rdzeń tej drogi nie
ma: jedyna brama do procesu prowadzi uruchomienie od startu do końca,
z obowiązkową granicą czasu i ubiciem całego drzewa procesów, czyli robi
dokładnie to, czego usłudze długożyjącej robić nie wolno. Wiersz poleceń
działa dziś i idzie tą samą bramą, co każdy inny program rdzenia.

Ustalenie zewnętrznego silnika niesie w polu propozycji całą treść panelu
po poprawce, tak samo jak przy regułach wbudowanych, ponieważ zastosowanie
poprawki wstawia ją w miejsce całej treści, nie w miejsce samego słowa.

Silnik nieobecny na maszynie nie jest odmową komendy: jego brak jest stanem
maszyny, o którym sonda startowa powiedziała przy uruchomieniu rdzenia,
a komenda oddaje to, co da się zmierzyć pozostałymi silnikami i regułami
wbudowanymi. Silnik obecny, który zawiódł w trakcie czynności, jest czymś
innym i wraca odmową, bo cisza w jego miejscu byłaby brakiem pomiaru
podanym jako brak zastrzeżeń. Oba silniki pisowni pracują tylko wtedy, gdy
język panelu da się sprowadzić do oznaczenia rozpoznawanego przez słownik —
Vale pracuje niezależnie od tego rozstrzygnięcia, bo mierzy powtórzenia
i terminy, które widać bez słownika.

Położenie zgłoszenia, którego nie da się potwierdzić w treści panelu, jest
brakiem pomiaru, nie ustaleniem: wstawienie poprawki pod zgadniętym numerem
znaku popsułoby panel w miejscu wybranym przypadkiem.

Przekład kategorii reguły LanguageToola na rodzaj kontroli kontraktu idzie
po kategorii, polu wyjścia programu, a nie po treści komunikatu, który jest
zdaniem w języku naturalnym i zmieni się przy pierwszym poprawionym
tłumaczeniu reguł programu. Kategoria nierozpoznana wraca jako rodzaj
gramatyczny — najszerszy z rodzajów kontraktu — zamiast być pominięciem:
ustalenie i tak dochodzi do odbiorcy, bo zgubione byłoby stratą, a źle
nazwane jest wciąż prawdziwe.

LanguageTool wagi ustalenia nie podaje, a kontrakt jej wymaga, więc rdzeń
rozstrzyga sam jedną zasadą: błędem jest to, co jest faktem o słowie — słowa
nie ma w słowniku — a ostrzeżeniem to, co jest orzeczeniem reguły o zdaniu.
Ta sama zasada dzieli reguły wbudowane komendy, gdzie odstęp przed
przecinkiem jest błędem, a cudzysłów prosty jedynie wskazówką.

Rozpoznanie języka panelu jest trzystopniowe. Napis o kształcie oznaczenia
języka idzie do programu bez zmiany, ponieważ rdzeń nie ma i nie powinien
mieć własnej tabeli języków LanguageToola: program zna ich kilkadziesiąt,
wykaz zmienia się z wydaniami, a oznaczenie nieznane program odrzuca sam,
wymieniając wszystkie, które zna — rdzeń sprawdza więc kształt napisu, nie
przynależność do wykazu. Nazwa własna języka sprowadza się do oznaczenia
podstawowego, bez odmiany krajowej, bo na przykład „angielski" nie mówi, czy
chodzi o pisownię brytyjską czy amerykańską, a te różnią się tysiącami słów —
dopisanie odmiany byłoby rozstrzygnięciem podjętym za użytkownika zamiast
przez niego. Napis, który nie jest ani oznaczeniem, ani znaną nazwą, nie
idzie do programu wcale: wysłany na chybił trafił dałby korektę treści
regułami niewłaściwego języka, a wynik wyglądałby jak korekta, będąc
zmyśleniem.

Kod wyjścia wywołania `hunspell -D` nie rozstrzyga niczego i nie jest
czytany: ten tryb diagnostyczny kończy się kodem błędu zawsze, także wtedy,
gdy wypisał komplet słowników. Odpowiedzią na pytanie, co program ma, jest
sam wykaz, nie kod wyjścia; wykaz pusty jest jedynym stanem, który znaczy,
że nie ma czym mierzyć. Dobór słownika hunspella rozstrzyga pytaniem
o rzeczywisty stan maszyny, nie tabelą wpisaną w kod: `hunspell -D` wypisuje
słowniki, które na tej maszynie stoją, i to jest jedyna prawda o tym, co
program otworzy — tabela języków w rdzeniu byłaby drugą, konkurencyjną
prawdą. Oznaczenie bez kraju dopasowuje się do słownika po samym języku
wyłącznie wtedy, gdy kandydat jest dokładnie jeden: dwóch kandydatów rdzeń
nie ma prawa rozstrzygnąć sam, bo `en_US` i `en_GB` różnią się pisownią
tysięcy słów.

Podmiana proponowana przez hunspella wchodzi wyłącznie przy jednym
wystąpieniu słowa w treści panelu: hunspell podaje przesunięcie względem
wiersza, nie całej treści, więc przy drugim wystąpieniu rdzeń nie wie,
o które chodzi, i wtedy ustalenie idzie bez propozycji.

Konfiguracja Vale wyłącza kontrolę pisowni tego programu rozstrzygnięciem,
nie przeoczeniem: jej słownik jest wyłącznie angielski i na treści w innym
języku zgłaszałby każde słowo jako błąd, podczas gdy ortografię w tej
komendzie prowadzą LanguageTool i hunspell, oba ze słownikiem języka panelu.
Wywołanie Vale idzie z opcją, która zdejmuje z programu kod wyjścia
niezerowy przy znalezionych zastrzeżeniach — bez niej każdy udany pomiar
wracałby przez arsenał serwerowy jako niepowodzenie programu, czyli
znaleziona usterka tekstu wyglądałaby jak usterka rdzenia. Kolejność
wyjścia programu Vale jest kolejnością mapy plików w pamięci, więc porządek
ustaleń nie byłby powtarzalny między wywołaniami bez sortowania jawnego.

Katalog roboczy wywołania programu korekty jest podany wprost, nie
zostawiony samej bramie wywołania procesu: Vale szuka swojej konfiguracji
w katalogu, z którego ruszył, więc uruchomienie gdzie indziej kończyłoby
się błędem braku konfiguracji, choć plik konfiguracji leży przygotowany.

Wycięty fragment treści potwierdza, że wskazanie programu zewnętrznego
mieści się w granicach treści panelu — potwierdzenie jest tu warunkiem, nie
ostrożnością, bo numer znaku podany przez program zewnętrzny liczy się
w jego własnej jednostce, a wstawienie poprawki pod numerem niesprawdzonym
popsułoby treść panelu w miejscu wybranym przypadkiem.

## budowa/server/internal/core/adapter_modul_badania_analiza.go

Sprzeczność wykryta między ustaleniami jest bytem zapisanym w bazie, nie
widokiem wyliczonym w locie: bez wiersza zapisanego w bazie rozstrzygnięcie
sprzeczności nie miałoby czego adresować, a samo rozstrzygnięcie nie
przeżyłoby odświeżenia panelu — sprzeczność wróciłaby do rozstrzygania
w kółko przy każdym kolejnym wykryciu.

Detekcja sprzeczności liczbowych heurystyką poprzedza wywołanie modelu:
rozbieżność dwóch liczb w zdaniach o tym samym przedmiocie da się wskazać
bez modelu, więc para ustaleń z tym samym rdzeniem słownym i różnymi
liczbami wchodzi jako sprzeczność z podaną różnicą już samą heurystyką.
Model dokłada sprzeczności, których sama liczba nie widzi, ale nie jest
warunkiem, żeby detektor w ogóle zadziałał.

Sprzeczność raz rozstrzygnięta nie wraca jako nowa przy kolejnym wykryciu:
wykrycie ma pokazywać rozbieżności wciąż otwarte, nie budzić te już zamknięte.

## budowa/server/internal/core/adapter_modul_developer_kontenery.go

Rozmowa z silnikiem kontenerów idzie biblioteką klienta Dockera wkompilowaną
w rdzeń, przez gniazdo silnika — ta sama zasada, co przy repozytorium, gdzie
rdzeń używa biblioteki `go-git` zamiast programu `git`: rdzeń nie startuje
procesu potomnego i nie zależy od tego, czy ktoś doinstalował klienta wiersza
poleceń. Podman wystawia to samo API OCI pod własnym gniazdem, więc obsługuje
się go tą samą drogą, wskazanym zmienną środowiskową gniazda silnika.

Biblioteka rozmawia z silnikiem, a silnika nie da się wkompilować: gdy na
serwerze nie ma ani Dockera, ani Podmana, odpowiedź wykazu kontenerów
przychodzi z polem dostępności silnika ustawionym na fałsz i pustym wykazem —
mówi jawnie o braku, zamiast udawać, że kontenerów po prostu nie ma.
Instalacja silnika po stronie serwera jest zmianą ciężką, poza zakresem tego
modułu.

Plik `docker-compose.yml` czyta i wykonuje rdzeń sam: rozbiera opis usług,
zakłada sieć stosu i startuje kontenery przez to samo API klienta. Obsługiwany
jest zakres używany w oknie — obraz, polecenie, zmienne, porty, wolumeny,
zależności — a nie każda konstrukcja, jaką format Compose zna; konstrukcja
nieznana wraca zdaniem nazywającym ją po nazwie, nie cichym pominięciem.

Silnik przeplata wyjście i diagnostykę kontenera jednym strumieniem,
znakując każdą porcję ośmiobajtowym nagłówkiem multipleksowania: bajtem
strumienia, trzema zerowymi i czterobajtową długością. Bez zdjęcia tych
nagłówków log w oknie miałby co kilkadziesiąt znaków wtrącone znaki
sterujące.

Silnik przyjmuje kontekst budowania obrazu wyłącznie jako strumień archiwum
tar — nie ma drogi wskazania mu katalogu, bo gniazdo silnika bywa po drugiej
stronie sieci. Pakowanie idzie strumieniem przez potok, nie do pliku
tymczasowego: kontekst dużego repozytorium ma setki megabajtów, a plik
tymczasowy tej wielkości zostawałby na dysku serwera po każdym nieudanym
budowaniu. Reguły `.dockerignore` są brane pod uwagę, bo bez nich do obrazu
wchodziłby katalog `.git` i wszystko, co repozytorium ma z założenia poza
obrazem. Dowiązania i pliki urządzeń zostają poza archiwum kontekstu, bo
kontekst budowania jest zbiorem plików zwykłych i katalogów, a dowiązanie
wskazujące poza katalog wyprowadziłoby budowanie z jego obszaru.

## budowa/server/internal/core/adapter_modul_badania_lektura.go

Źródło niesie tekst w swoim wierszu, jeżeli został już wydobyty —
przechwyceniem strony, transkrypcją albo wcześniejszym rozpoznaniem pisma.
Gdy tekstu nie ma, treść powstaje z załącznika pełnotekstowego przez port
arsenału dokumentowego, ten sam Pandoc i poppler, którymi czyta moduł
Library, nie drugą ich odmianę. Tekst raz wydobyty zostaje przy źródle:
druga lektura tej samej pozycji nie ma po co uruchamiać arsenału powtórnie.

Rozmowa oparta na korpusie odpowiada wyłącznie z treści wskazanych źródeł
i mówi wprost, czy odpowiedź jest zakotwiczona w nich. Żądanie z wymogiem
zakotwiczenia przy pustym korpusie kończy się odmową, nie odpowiedzią
modelu z pamięci — odpowiedź badawcza bez źródła jest w tym module gorsza
niż brak odpowiedzi.

Wykrywanie tabel w tekście źródła idzie po siatce znaków: wiersz tabeli ma
ten sam rozkład separatorów co jego sąsiedzi. Jest to heurystyka i tak się
nazywa — tabela wykryta wchodzi do odpowiedzi, a rozstrzygnięcie, czy ją
utrwalić, zostaje po stronie odbiorcy.

Dobór fragmentów korpusu najbliższych pytaniu idzie po pokryciu słów
pytania: bez magazynu wektorów jest to miara uboższa, ale prawdziwa — mierzy
treść, która rzeczywiście leży w korpusie, a nie podobieństwo obiecane przez
usługę, której instalacja nie niesie.

## budowa/server/internal/core/adapter_modul_roundtable_analiza.go

Cztery rodzaje analizy debaty są pracą modelu, bo wymagają rozumienia
treści: wydobycie argumentów, klasyfikacja aktów mowy, wykrycie błędów
logicznych, kontrola steelman, weryfikacja faktyczności i sygnalizacja tonu.
Jeden rodzaj modelu nie wymaga: scalenie powtórzeń liczy się podobieństwem
treści węzłów, a rdzeń robi to sam — wołanie modelu, żeby porównał dwa
zdania, kosztowałoby wywołanie kanału za robotę, którą wykonuje arytmetyka,
i dawałoby wynik niepowtarzalny między przebiegami.

Każdy przebieg analizy zostawia ślad w bazie, nie tylko w odpowiedzi
komendy: wydobycie argumentów zastępuje graf, klasyfikacja aktów mowy
znakuje wypowiedzi, wykrycie błędów stawia oznaczenia na węzłach, weryfikacja
faktyczności wypełnia rejestr dowodów. Odpowiedź komendy jest odczytem tego,
co zostało zapisane, nie jedynym miejscem, w którym wynik istnieje.

Próg scalenia węzłów jest wysoki rozmyślnie: scalenie dwóch argumentów,
które tylko brzmią podobnie, zabiera z grafu jeden głos i zawyża poparcie
drugiego.

Model dostaje zapis debaty z kodami wypowiedzi i ma je powtórzyć przy
każdej wydobytej jednostce argumentacyjnej. Dzięki temu węzeł grafu wraca do
swojej wypowiedzi i do mówcy — bez tego graf byłby zbiorem zdań bez autora,
a okno analizy nie miałoby czego pokazać w kolumnie uczestnika.

Krawędzie grafu argumentów powstają też z relacji, którą rdzeń widzi bez
modelu: węzły z tej samej wypowiedzi wspierają się nawzajem po kolei —
pierwszy jest tezą, każdy następny wsparciem poprzedniego. Relacji między
wypowiedziami różnych uczestników rdzeń nie zgaduje; od tego jest kontrola
steelman i wykrywanie sporu.

Znakowanie aktu mowy idzie do zapisu wypowiedzi, nie tylko do osobnego
ustalenia, ponieważ okno debaty pokazuje etykietę przy samej wypowiedzi,
nie w osobnym wykazie.

Zakres wykrywania błędów logicznych bierze się z katalogu okna: błąd
wyłączony w katalogu nie jedzie w poleceniu do modelu i nie zostaje
oznaczony, choćby model go i tak nazwał — wyłączenie, które nie wyłącza,
byłoby ustawieniem bez skutku.

Brak rozpoznanych błędów logicznych jest wynikiem analizy, nie usterką:
debata bez chwytów erystycznych jest debatą poprawną. Ustalenie zbiorcze
mówi to wprost, zamiast oddawać pustkę nie do odróżnienia od analizy, która
nie ruszyła.

Kod wypowiedzi wiodący wiersz w nawiasie kwadratowym jest odrzucany, gdy
nie należy do wykazu znanych wypowiedzi, a reszta wiersza zostaje w całości:
model, który nawiasu użył do czegoś innego, nie ma prawa przypiąć ustalenia
do wypowiedzi, której nie ma. Wykaz pusty przepuszcza każdy kod, bo woła się
tak wtedy, gdy wołający nie ma wykazu, z którym by porównywał.

## budowa/server/internal/core/adapter_modul_aplikacje_produkt.go

Oś czasu projektu nie ma własnej tabeli zdarzeń i mieć jej nie może:
`AppTimelineKind` niesie dokładnie pięć wartości — architektura, warsztat,
etap, wdrożenie, pakiet — i każda z nich ma już w bazie swój wiersz ze
znacznikiem czasu. Osobna tabela zdarzeń byłaby szóstą kopią tych samych
faktów, rozjeżdżającą się przy pierwszym zapisie, który zapomni ją dopisać.
Oś czasu składa się więc z odczytów tych pięciu źródeł i sortuje je po
czasie, malejąco, a przy równym czasie po identyfikatorze zdarzenia, bo pięć
źródeł zapisuje znaczniki z rozdzielczością milisekundy i dwa zdarzenia tej
samej milisekundy bez drugiego klucza zamieniałyby się miejscami między
wywołaniami.

Powiązania modułów są mierzone, nie deklarowane: kontrakt nie daje komendy,
którą dałoby się włączyć powiązanie ręcznie, więc gdyby rdzeń oddawał tu
stałą listę ze stanem włączenia niezależnym od zawartości okna, meldowałby
wykaz, za którym nic nie stoi. Każde powiązanie liczy więc w bazie danych to,
czym naprawdę żyje w tym oknie, a pole `detail` mówi, co policzono. Powiązanie
czynne znaczy „w tym oknie to powiązanie ma już treść", nie że ktoś zaznaczył
przełącznik.

Pole `ownerAgentId` przy zapisie etapu rozróżnia trzy stany wykonawcy: brak
wskazania zostawia zastane przypisanie, pusty łańcuch zdejmuje przypisanie
wykonawcy, a wartość przypisuje nowego. Zapis etapu rozgłasza też zdarzenie
`apps.build.changed`, bo panel Product Buildera rysuje oś etapów z tego
samego zdarzenia, którym dostaje przejścia wdrożenia.

Usunięcie kamienia milowego, którego nie ma, kończy się odmową `not_found`,
nie polem `deleted: false`: klient odróżnia „usunięto" od „nie było czego
usunąć" po odmowie, bo pole logiczne oddane jako powodzenie kazałoby mu
zgadywać, czy operacja się odbyła.
