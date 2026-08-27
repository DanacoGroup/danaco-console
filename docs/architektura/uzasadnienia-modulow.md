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
