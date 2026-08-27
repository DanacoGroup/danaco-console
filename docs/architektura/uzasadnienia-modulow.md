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

## budowa/server/internal/core/adapter_modul_aod_wyciszenie.go

Wyciszenie nakładki Always On Display jest bytem rdzenia, nie wyłącznie
ustawieniem okna: zapis miejscowy w przeglądarce (`client/src/aod/wyciszenie-aod.ts`)
pozwalał wyciszyć sugestie w jednej powłoce, podczas gdy druga powłoka nadal je
ujawniała. Rdzeń trzyma jeden stan wyciszenia współdzielony między powłokami.

Rodzaje wyciszenia są trzy: czasowe, kontekstowe (moduł albo karta sesji)
i klasy zdarzeń. Tryb cichy nakładki nie jest wyciszeniem, lecz ustawieniem
trybu obecności, a wyjątek wagi krytycznej przebija każde z trzech wyciszeń
niezależnie od ich zakresu.

Założenie i zniesienie wyciszenia idą jedną komendą `aod.mute.set`,
rozróżnianą polem `muted`. Osobna komenda znoszenia zmuszałaby okno do
pamiętania dwóch dróg do jednej czynności zamiast jednej.

Sygnał odkłada się w rdzeniu nawet wtedy, gdy wpada w wyciszenie czynne:
wyciszenie wstrzymuje wyłącznie ujawnienie sygnału operatorowi, nie jego
zapis. Odpowiedź `aod.signal.report` i odpowiedź `aod.signal.list` nazywają
wprost, które wyciszenie wstrzymało dany sygnał, żeby brak sygnału w oknie
nie wyglądał na jego brak w rzeczywistości.

## budowa/server/internal/core/adapter_modul_badania_uchwyty.go

Plik wiąże jedną funkcją komplet komend portu Badania z rejestrem rdzenia,
niezależnie od tego, który z dziesięciu plików adaptera modułu Research je
wypełnia — ten sam wzór stosuje wpięcie komend modułu Library.

Rozgłaszane są wyłącznie działania, które naprawdę opisuje jedno z czterech
zdarzeń zmiany kontraktu: raportu, źródła, ustalenia i monitora. Działanie
bez odpowiadającego mu zdarzenia w kontrakcie nie wprowadza zastępczego
zdarzenia własnego.

Trzy drogi wnoszące źródło — dodanie, przechwycenie strony i import —
rozgłaszają to samo zdarzenie zmiany źródła z rodzajem „utworzono”, żeby okno
katalogu zobaczyło nową pozycję bez odpytywania; usunięcie rozgłasza to
zdarzenie z samym wskazaniem zdjętej pozycji, bo pełny byt już nie istnieje.

## budowa/server/internal/core/adapter_modul_roundtable_ocena.go

Wskazanie wypowiedzi bardziej przekonującej jest pojedynkiem dwóch tożsamości
i tak jest liczone: punktacja obu przesuwa się natychmiast po zapisaniu
oceny, nie dopiero po zamknięciu debaty. Ranking, który aktualizuje się
dopiero na koniec, nie pokazywałby niczego w trakcie debaty, a Operator
ocenia właśnie w trakcie.

Elo przesuwa obie punktacje o wartość zależną od różnicy między nimi.
Glicko i TrueSkill dokładają do tego niepewność oszacowania: im mniej
pojedynków ma tożsamość, tym większy krok przesunięcia, bo tym mniej
wiadomo o jej sile. Różnica między tymi dwoma algorytmami leży w tempie
zawężania tej niepewności; niepewność maleje z każdym pojedynkiem, ale nie
schodzi do zera, bo tożsamość, o której „wiadomo wszystko", przestałaby
reagować na kolejne wyniki.

Suma wag kryteriów rubryki ma wynosić jedność, tak jak mówi kontrakt.
Rubryka o sumie innej dawałaby wynik werdyktu, którego nie da się porównać
z wynikiem z innej rubryki, więc odmowa zapisu jest tu jedyną uczciwą
odpowiedzią.

Sędzia oceniający wypowiedź jest uczestnikiem składu debaty, więc ocenia
własnym kanałem modelu i własną tożsamością. Punkty odczytuje się z jego
odpowiedzi tekstowej; odpowiedź bez rozpoznanych liczb zostawia punkty
zerowe dla danego kryterium, a uzasadnieniem werdyktu jest to, co sędzia
naprawdę powiedział — rdzeń nie wystawia oceny za niego.

Kryterium, którego sędzia nie ocenił w odpowiedzi, dostaje zero i tak też
wchodzi do wyniku ważonego. Pominięcie takiego kryterium w mianowniku
dawałoby wynik wyższy za odpowiedź uboższą niż za odpowiedź pełną.

Rachunek dziesięciu do potęgi w krzywej oczekiwania Elo liczy się szeregiem
wykładniczym 10^x = e^(x·ln10), zbieżnym do wystarczającej dokładności po
kilkunastu wyrazach, bo argument jest tu zawsze mały — różnica punktacji
dzielona przez czterysta.

## budowa/server/internal/core/adapter_modul_developer_debug.go

Punkt przerwania jest trwały i należy do okna: Operator stawia go w marginesie
edytora, zanim cokolwiek uruchomi, i oczekuje go zastać przy drugim i trzecim
biegu. Dlatego leży w bazie danych i przeżywa sesję.

Sesja debugowania jest żywa: ma uchwyt do procesu adaptera i do procesu
debugowanego, i kończy się razem z nimi. Nie ma jej w bazie, ponieważ wiersz
opisujący sesję, do której rdzeń stracił uchwyt, byłby wpisem o czymś, czego
już nie ma. Ta sama zasada rządzi przebiegami budowania.

Adapter dowiaduje się o punktach przerwania raz, między komendą inicjalizacji
a domknięciem konfiguracji, ponieważ tak wymaga protokół. Sesja startująca
bierze więc komplet punktów okna z bazy i podaje je adapterowi plik po pliku;
punkt postawiony później dosyła się do sesji czynnej tą samą drogą.

Nasłuch zakłada rdzeń, a adapter dzwoni do niego adresem klienta. Odwrotna
kolejność wymagałaby odczytania portu z wyjścia adaptera, czyli rozbioru
tekstu spoza protokołu. Port wybiera jądro systemu, więc dwie sesje naraz nie
zderzą się o ten sam numer.

## budowa/server/internal/core/adapter_modul_aplikacje_architektura.go

Walidacja architektury liczy zastrzeżenia, nie tylko je przechowuje: migracja
opisana kontraktem zostawiła w architekturze kolumnę zastrzeżeń walidacji,
a definicja układu przepisuje ją bez zmiany, bo nie ma z czego liczyć nowych
wartości. Komenda walidacji przechodzi komponenty i graf zależności i wykrywa
cztery rzeczy: komponent bez ani jednej krawędzi, krawędź wskazującą
komponent spoza układu, cykl w grafie zależności oraz układ bez ani jednego
komponentu. Policzone zastrzeżenia wracają do kolumny, więc kolejny odczyt
architektury pokazuje je bez powtarzania rachunku, a definicja układu, która
je przepisuje bez zmiany, nie kasuje pracy walidatora.

Żadne zastrzeżenie niczego nie blokuje: kontrakt nazywa je ostrzeżeniem, nie
bramą, więc walidacja jest komendą odczytu z zapisem wyniku, nie warunkiem
zapisu układu.

Eksport wytwarza plik, nie opis pliku. Cztery formaty kontraktu powstają
bibliotekami wkompilowanymi w binarium — svg, mermaid i markdown są tekstem
składanym w tym module, png rysuje obraz rastrowy biblioteką standardową
wraz z czcionką rastrową wkompilowaną w binarium. Żaden format nie woła
programu zewnętrznego, więc eksport działa na instalacji, która niesie sam
rdzeń.

## budowa/server/internal/core/adapter_modul_aod.go

Proces przypinany komendą `aod.observe.attach` jest procesem telemetrii
postępu, tym samym, którego identyfikator oddaje `monitor.status` i który
jedzie w `progress.changed`. Nakładka przypina więc proces, nie okno ani
kartę sesji: okno i karta wchodzą do stanu nakładki osobnymi polami.

Przypięcia procesów nie mają wiersza w bazie, w przeciwieństwie do wyciszeń.
Rejestr procesów telemetrii żyje w pamięci rdzenia i proces ginie razem
z rdzeniem, więc wiersz przypięcia, który przeżyłby restart, wskazywałby
proces nieistniejący. Przypięcia mieszkają zatem w pamięci, tak samo jak
rejestr procesów telemetrii, którego dotyczą.

Ognisko karty sesji i okna jest własnością klienta, a żądania rodziny aod.*
niosą wyłącznie urządzenie: wiązania urządzenie-klient nie ma ani w
schemacie, ani w kontrakcie. Rdzeń wskazuje więc fakt, który sam zna: ostatni
punkt pracy odnotowany przez telemetrię — to samo źródło, którym strona
główna rozstrzyga, dokąd prowadzi powrót do sesji.

Niewpięty rejestr telemetrii odmawia kodem błędu wewnętrznego w komendzie
`aod.status.get`, bo zero procesów w biegu byłoby wtedy ciszą udającą pomiar,
nie faktem.

Odpięcie procesu nie wymaga jego istnienia w rejestrze telemetrii: po
opróżnieniu rejestru przypięcie ma dać się zdjąć, choć wtedy jest to
najbardziej potrzebne. Bytem, którego brak rozstrzyga odmowę, jest samo
przypięcie, nie proces.

Identyfikatorem urządzenia w kontrakcie jest numer wiersza katalogu maszyn;
tak samo czyta go rodzina komend accessPoint.

## budowa/server/internal/core/adapter_modul_przegladarka_inspekcja.go

Wszystkie pięć komend tego modułu wymagają strony wykonanej, nie jej źródła:
drzewo elementów po zbudowaniu przez skrypty, rejestr żądań, komunikaty
konsoli, metryki emulowanego urządzenia i stan po przewinięciu istnieją
dopiero wtedy, gdy stronę ktoś naprawdę uruchomił. Dlatego idą silnikiem
przeglądarki, nie samym pobraniem HTTP, i dlatego odmawiają wprost, gdy
silnika na maszynie nie ma, zamiast oddawać puste wykazy udające inspekcję.

Adres strony bierze się z ostatniej migawki okna: kontrakt nie niesie w tych
żądaniach pola adresu, tylko pyta o bieżącą stronę okna, a bieżącą stroną
jest ta, dokąd okno ostatnio przeszło. Okno bez migawki dostaje odmowę
`not_found`.

Zapis wytworu rejestru sieciowego ląduje w magazynie modułu, a pole `harRef`
odpowiedzi wskazuje plik, który naprawdę leży na dysku, nie sam rejestr
osadzony w treści odpowiedzi.

Przewinięcie strony jest czynnością na stronie żywej, więc idzie silnikiem,
a jego skutek widać dopiero w migawce pobranej po przewinięciu: tekst po
przewinięciu bywa inny niż przed nim, bo treść bywa dogrywana przy
przewijaniu, a położenie przewinięcia jest mierzone, nie zakładane.

Nastawa emulacji nazwana daje wymiary urządzenia typowe dla tej nazwy,
nastawa `custom` bierze wymiary z pól żądania, a nastawa `reset` wraca do
wymiarów biurka — zdjęcie emulacji też jest emulacją o znanych metrykach,
nie brakiem odpowiedzi.

Nazwa poziomu komunikatu konsoli nierozpoznana przez rdzeń spada na poziom
`log`: komunikat ma się pokazać Operatorowi, a nie zniknąć przez nieznane
słowo w polu poziomu zdarzenia protokołu.

## budowa/server/internal/core/adapter_modul_library_sugestie.go

Domyślnym zachowaniem modułu jest sugestia z akceptacją operatora, nie zapis
bez pytania: klasyfikacja wytwarza wiersze sugestii, a dopiero zatwierdzenie
zamienia je w zmianę zasobu. Żądanie z polem zatwierdzenia od razu wykonuje
obie czynności naraz, co pozostaje jawnym wyborem wołającego, nie domyślnym
zachowaniem rdzenia.

Sugestie mają dwa źródła i oba są rzeczywiste. Model proponuje etykiety, gdy
żądanie wskazuje kanał, a rdzeń ma wpięty rejestr kanałów; jego odpowiedź jest
propozycją, więc trafia do tabeli sugestii, nie wprost na zasób. Pomiar rdzenia
działa zawsze: zasób bez etykiety i bez kolekcji jest osierocony, zasób
dzielący sumę kontrolną z innym jest duplikatem, a zasób niosący numer PESEL,
NIP albo numer rachunku jest treścią wrażliwą do przeglądu. Te trzy reguły są
policzone, nie zgadnięte, i działają także wtedy, gdy modelu nie ma wcale.

Uzasadnienie towarzyszy każdej sugestii obowiązkowo: przyjęcie propozycji bez
podanego powodu byłoby zaufaniem bez podstawy.

## budowa/server/internal/core/adapter_modul_developer_api.go

Cała rodzina komend stoi na bibliotekach wkompilowanych w rdzeń: `net/http`
wykonuje zapytanie, `getkin/kin-openapi` czyta kontrakt. Nie startuje tu ani
jeden proces potomny i nie ma tu żadnej zależności od programu spoza
instalki — wołanie zewnętrznego narzędzia wiersza poleceń dałoby to samo,
tyle że zależne od tego, czy to narzędzie stoi na maszynie.

Adres, nagłówki i treść zapytania przechodzą przez podstawienie `{{nazwa}}`
wartościami wskazanego środowiska. Bez tego kroku każde zapytanie kolekcji
miałoby wpisany na stałe adres jednego serwera, a przeniesienie kolekcji
między środowiskiem przejściowym a produkcyjnym byłoby przepisywaniem jej
w całości.

Rodzina nie wykonuje zapytań do adresów spoza sieci, do których serwer i tak
nie ma dostępu, i nie zna poświadczeń Operatora — nagłówek uwierzytelniający
podaje wołający albo środowisko kolekcji. Sekret wpisany w środowisko leży
w bazie jawnie i tak też jest opisany; miejscem na sekret jest sejf,
a odwołanie do niego wchodzi jako wartość nagłówka.

Rodzaj treści zapytania uzupełnia nagłówek `Content-Type` tylko wtedy, gdy
wołający sam go nie podał: nagłówek jawny jest zawsze mocniejszy od domysłu
opartego na rodzaju treści.

Nagłówek odpowiedzi powtórzony kilka razy zostaje sklejony przecinkiem przy
spłaszczeniu do mapy, bo kontrakt niesie mapę nazwa-wartość, a zgubienie
drugiej wartości byłoby gorsze od jej sklejenia z pierwszą.

Kontrakt OpenAPI przychodzi z pliku repozytorium albo z adresu. Plik jest
drogą podstawową: kontrakt leżący w repozytorium jest wersjonowany razem
z kodem, więc kolekcja z niego wytworzona opisuje ten sam stan usługi, co
gałąź, w której pracuje Operator.

Kolejność ścieżek i metod przy imporcie z kontraktu jest ustalona sortowaniem,
bo mapa Go oddaje wpisy w kolejności losowej — dwa importy tego samego pliku
dałyby bez sortowania dwie różne kolekcje, a ich porównanie nie mówiłoby
niczego o rzeczywistej zmianie kontraktu.

## budowa/server/internal/core/adapter_modul_developer_dane.go

Rodzina stoi na jednolitym dostępie do baz danych oraz na trzech sterownikach
wkompilowanych w rdzeń: PostgreSQL, MySQL i SQLite, ten sam silnik, na którym
stoi baza produktu. Konsola SQL otwiera się więc wszędzie tam, gdzie stoi
rdzeń, i nie zależy od zewnętrznych klientów wiersza poleceń.

Wiersz połączenia opisuje, do czego się łączyć: silnik, host, port, bazę
i użytkownika. Hasło stoi w sejfie pod odwołaniem z opisu połączenia, tak samo
jak klucze dostawców modeli. Kopia sekretu w drugim miejscu, którego nikt nie
rotuje, jest usterką bezpieczeństwa, a nie wygodą.

Połączenie oznaczone jako tylko do odczytu odrzuca polecenie zmieniające dane,
zanim cokolwiek wyjdzie do silnika. Poleganie na uprawnieniach po stronie bazy
byłoby poleganiem na nastawie, której z tego okna nie widać; konsola SQL nad
produkcją bez tej bramy jest jedną nieuwagą od szkody.

## budowa/server/internal/core/adapter_modul_badania_raport_dobudowa.go

Wersja zapisywana wyłącznie na żądanie oznaczałaby, że osoba, która o nią nie
poprosiła, nie ma do czego wrócić, a to jest dokładnie ta chwila, w której
wersje są potrzebne. Dlatego migawkę zakłada każda operacja zmieniająca treść
raportu, nie tylko jawne żądanie zapisu.

Porównanie wersji liczy się z migawek odczytanych z bazy, nie ze stanu
w pamięci procesu: porównanie po restarcie rdzenia daje ten sam wynik co przed
restartem.

## budowa/server/internal/core/adapter_modul_mobile.go

Proces tej rodziny nie jest tym samym bytem, co proces terminala: pole stanu
procesu mobilnego niesie słownik telemetrii postępu, podczas gdy stan procesu
terminala niesie własny, odrębny słownik. Proces mobilny jest więc wpisem
rejestru telemetrii postępu, tym samym, który pokazuje Process Monitor
warstwy wspólnej — plik jest fasadą nad tym rejestrem, nie osobnym magazynem
ani osobną tabelą.

Adapter mobilny osadza adapter monitora zamiast zakładać port równoległy nad
tym samym rejestrem i dokłada dwa porty czynności, których warstwa wspólna
sama nie ma: rozmowę, czyli zatrzymanie tury okna, oraz kolejki, czyli
działania na kolejce.

Sterowanie procesem nie ma własnego silnika: proces kolejki idzie tą samą
drogą, co działanie kolejki, a tura okna tą samą, co zatrzymanie wiadomości,
razem z ich telemetrią i rozgłoszeniami. Gdyby warstwa mobilna zatrzymywała
turę własnym wywołaniem, stan procesu zmieniłby się bez zdarzenia postępu
i pulpit pokazywałby proces w biegu, którego już nie ma.

Czego rdzeń nie umie, tego ta rodzina nie udaje: tura okna nie ma w rdzeniu
ani wstrzymania, ani wznowienia, ani ponownego uruchomienia, więc takie
sterowanie odmawia głośno zamiast oddać proces rzekomo już po sterowaniu.
Modyfikacja zlecenia zmienia jego treść, a żądanie sterowania niesie
wyłącznie proces, sterowanie i urządzenie, więc korekta treści zlecenia
odmawia i wskazuje drogę, którą treść naprawdę się poprawia: polecenie
w oknie rozmowy.

Trzy przełożenia sterowania kolejki są jednoznaczne, bo oba słowniki nazywają
tę samą czynność. Ponowne uruchomienie idzie na start pozycji, nie na
ponowienie, bo ponowienie dotyczy jednej pozycji kolejki, a sterowanie
mobilne dotyczy procesu jako całości. Zatwierdzenie idzie na wznowienie i nie
jest to przełożenie na skróty: krok naprzód znaczy przyjęcie wyniku, tak jak
nazywa to protokół silnika kolejki, a pozycja czekająca na weryfikację
przechodzi tym krokiem do stanu ukończonej z werdyktem przyjęcia — czyli
zatwierdzenie ma skutek zapisany, nie samą etykietę odpowiedzi.

Zatwierdzenia tura okna nie ma i nie jest to przeoczenie: punktu decyzyjnego
kontrakt dla niej nie zna, bo zatwierdzać można wyłącznie tam, gdzie stan
czekający naprawdę stoi zapisany, czyli w pozycji kolejki. Tura, która nie
biegła, jest sporem stanu, nie powodzeniem, więc odpowiedź mobilna, która nie
ma pola mówiącego wprost, że nic się nie zatrzymało, mówi prawdę odmową.

Pole nazwy procesu mobilnego jest wymagane, więc musi nieść coś prawdziwego:
rejestr telemetrii zna nazwę etapu, nie nazwę procesu, a etap w polu nazwy
byłby podaniem jednej rzeczy za drugą, więc nazwą procesu jest nazwa bytu,
którego dotyczy — tytuł okna albo nazwa kolejki, a bez nazwy własnej opis
wskazujący byt jednoznacznie. Pola chwili uruchomienia nie ma i nie jest to
przeoczenie: rejestr telemetrii nie znakuje procesów chwilą startu, zna
wyłącznie stan i etap, a pole jest opcjonalne kontraktem.

## budowa/server/internal/core/adapter_modul_isolation.go

Jedenaście punktów izolacji to jedenaście kluczy tabeli `ustawienie` — tych
samych, które rozstrzyga rejestr definicji i egzekwuje warstwa sesji.
Rodzina jest drugim wejściem do tego samego magazynu, nie drugim magazynem:
wartość zapisana tutaj wraca też przez odczyt konfiguracji ogólnej,
i odwrotnie. Wykaz punktów stoi w rejestrze definicji izolacji i tylko tam;
ten plik odwzorowuje wyliczenia kontraktu na klucze i nie zna nazwy
pojedynczego punktu.

Wartość zapisu musi nieść dokładnie tę wartość, którą czyta egzekutor przy
rozstrzyganiu polityki; wyprowadzenie jej z wartości domyślnej rejestru
odwracałoby znaczenie zapisu przy zmianie stanu wyjściowego platformy.

Pusty wykaz poziomów zasięgu znaczy bazę bez słownika poziomów — awarię
podłoża, nie brak poziomów jako stan normalny — więc odpowiedzią jest
odmowa, nie pusty wynik. Pola żądania wskazujące sesję i okno nie mają
odpowiednika w wyniku, bo struktura poziomu zasięgu nie niesie pola na byt
poziomu, więc rdzeń niczym ich nie zawęża.

Wykaz przełączników kontekstu zwracany przy odczycie jest zawsze pełny:
macierz okna ma trzy wiersze niezależnie od tego, ile z nich naprawdę
zapisano, a punkt bez zapisu na tym poziomie niesie wartość domyślną
z rejestru definicji. Żądanie zapisu bez ani jednego przełącznika jest
odmawiane, bo pusta tablica zamieniłaby zapis w potwierdzenie bez zmiany.

Adresem zapisu ustawienia jest czwórka: poziom zasięgu, byt poziomu, oś
i byt osi — warstwy w niej nie ma, ani w schemacie, ani w rozstrzyganiu, ani
w egzekutorze. Karta sesji jest natomiast jednym z poziomów zasięgu, więc
warstwa domyślna kieruje zapis pod wskazany poziom, a warstwa sesyjna
kieruje zapis na poziom karty sesji; warstwa sesyjna wskazana razem z innym
poziomem jest sprzecznością i wraca odmową, zamiast po cichu wybrać jedno ze
wskazań. Warstwa sesyjna nie ma osobnej przestrzeni kluczy ani kolumny:
wartość zapisana poza adresem rozstrzygania nie doszłaby do egzekutora,
a zapis zostałby potwierdzony mimo braku skutku.

Byt poziomu jest wymagany poza poziomem globalnym: wiersz zapisany z pustym
bytem na poziomie węższym nie należy do żadnej karty, roli ani okna, więc
rozstrzygacz nigdy po niego nie sięgnie — zapis wyglądałby na udany, nie
robiąc nic.

Reguła odcięcia punktu jest przepisana z egzekutora co do znaku: wymiar
kontekstu zostaje odrębny, dopóki nie zapisano wprost współdzielenia,
a zakres techniczny jest włączony wyłącznie wtedy, gdy zapisano wprost
włączenie.

## budowa/server/internal/core/adapter_modul_badania.go

Typ adaptera i konstruktor stoją w tym pliku, a nie w plikach raportu
i eksportu tego samego modułu, bo wszystkie metody wiszą na tym samym typie
adaptera, tak jak warstwa danych rozkłada jedno repozytorium modułu na kilka
plików wedle odpowiedzialności. Deklaracja typu w dwóch miejscach byłaby
dwiema prawdami o jednym bycie.

Rdzeń nie ocenia wiarygodności źródła sam: komenda dodania źródła zapisuje
ocenę dokładnie taką, jaką podało żądanie, nie wylicza jej z adresu ani
z treści i nie sięga do sieci po źródło. Brak oceny w żądaniu zostaje pusty,
a warstwa danych nadaje mu wtedy wartość domyślną oznaczającą źródło
niezweryfikowane — brak oceny nie udaje oceny wyliczonej.

Odwołania do bajtów, które wychodzą z portów modułu, są ścieżkami
względnymi wobec korzenia magazynów rdzenia — to granica kontraktu: rdzeń
nie wypuszcza ścieżek bezwzględnych z dysku operatora. Odczyt pliku modułu
sprowadza więc wskazanie względne do tego korzenia przed otwarciem.

Żądania modułu Research nie niosą identyfikatora kanału modelu, więc operacje
słowne — streszczanie źródła, porównanie źródeł, redakcja sekcji raportu —
biorą domyślnie pierwszy czynny kanał rejestru, ten sam, który komenda
wykazu kanałów pokazuje jako gotowy. Brak wpiętego rejestru albo brak
czynnego kanału kończy się odmową wprost, nigdy streszczeniem złożonym bez
udziału modelu.

Kod źródła nieistniejącego, wskazany przy zapisie ustalenia, jest pomijany
po stronie warstwy danych, nie wywraca całego zapisu; odpowiedź komendy
niesie komplet źródeł naprawdę powiązanych, nie powtórzenie żądania.

Ślad prowenancji ustalenia zaczyna się przy jego powstaniu, nie przy
pierwszej zmianie, żeby historia ustalenia nie zaczynała się od drugiego
zdarzenia.

## budowa/server/internal/core/adapter_modul_tlumaczenie_uchwyty.go

Port `Tlumaczenie` jest jeden, choć metody adaptera leżą w kilku plikach na
wspólnym typie: trzy porty nad jednym repozytorium byłyby trzema prawdami
o jednej powierzchni kontraktu.

Rdzeń tłumaczy modelem, związany słownikiem operatora. Dodanie panelu
przekłada tekst źródłowy okna na język panelu czynnym kanałem modelu; do
polecenia wchodzą terminy słownika, zakazy tłumaczenia, ton panelu i zasady
jakości, a wynik przechodzi mechaniczną podmianę terminów i migawkę kontroli
jakości. Bez czynnego kanału operacja odmawia wprost. Zapis źródła sam
wyłącznie zapisuje tekst źródłowy i zakłada okno. Tłumaczenie zwrotne też
woła model: przekłada treść panelu z powrotem na język źródłowy okna, bez
słownika, żeby kontrola wierności miała co kontrolować; odmawia, gdy panel
jest pusty albo okno nie ma języka źródłowego.

Rdzeń rozpoznaje język modelem: rozpoznanie źródła pyta model o język
tekstu, a bez czynnego kanału odmawia wprost, nie zgaduje. Kanał wskazuje
wołający polem nieobowiązkowym; jego brak bierze kanał domyślny czynny,
a kanał wskazany, lecz nieznany albo nieczynny, kończy się nazwaną odmową
zamiast cichego zejścia na domyślny. Rozpoznanie źródła tego pola nie ma
i jedzie kanałem domyślnym.

Rdzeń syntezuje mowę lokalnym syntezatorem, tym samym portem i przez tę samą
bramę izolacji, co silnik rozpoznawania mowy, po czym oddaje ścieżkę
powstałego nagrania. Brak programu, brak głosu dla języka panelu albo brak
treści panelu kończy się odmową nazywającą brak, nie pustą ścieżką udającą
nagranie.

Pamięć tłumaczeń ma pisarza: pary segmentów, z których budowane są
podpowiedzi, zapisują dodanie panelu po przekładzie modelu i ustawienie
tłumaczenia po korekcie operatora. Sparowanie zdanie-do-zdania idzie
wyłącznie przy równej liczbie segmentów po obu stronach; przy nierównej
zapisywana jest jedna para całościowa zamiast zmyślonego dopasowania.

Rdzeń nie ma magazynu blobów: eksport słownika i eksport panelu zapisują
ślad, ale nie wytwarzają pliku — ten sam brak co w module Library i module
Research.

## budowa/server/internal/core/adapter_modul_library_higiena.go

Wszystkie pięć czynności higieny repozytorium są czynnościami rdzenia, nie
okna, z jednego powodu: sięgają po to, czego okno nie ma. Rozpoznanie
przybliżone żąda porównania każdego zasobu z każdym, weryfikacja
integralności — przeliczenia sumy kontrolnej z bajtów leżących na dysku,
pulpit stanu — przejścia po całym zbiorze. Okno liczące te rzeczy z odczytanej
strony wykazu orzekałoby tylko o próbce, nie o repozytorium.

Rozpoznanie dokładne idzie po sumie kontrolnej i jest pewne. Rozpoznania
przybliżone niosą trafność, bo to wnioski: podobieństwo treści liczy się
odciskiem słów, podobieństwo obrazu — odciskiem percepcyjnym zbudowanym
z obrazu zeskalowanego do siatki. Oba liczy wkompilowany kod Go; żaden nie
woła programu z zewnątrz.

Zasoby o identycznej sumie kontrolnej są pomijane przy rozpoznaniu
przybliżonym: należą do rozpoznania dokładnego, a wystawione tu po raz drugi
kazałyby Operatorowi rozstrzygać tę samą parę dwa razy.

Zasoby wchłaniane przy scaleniu duplikatów trafiają do archiwum, nie
znikają: scalenie bywa pomyłką, a kosz repozytorium jest odwracalny —
usunięcie trwałe ma własną komendę i własne potwierdzenie. Kod wersji jest
unikalny w całym repozytorium, więc wersja wchłonięta przy zachowaniu
historii wchodzi jako nowy wpis pod własnym identyfikatorem, nie
przeniesiona wprost — przeniesienie wprost byłoby zderzeniem kluczy.

Brak treści zasobu pod odwołaniem przy weryfikacji integralności jest
trzecim wynikiem, nie odmianą niezgodności: sumy nie ma czego porównać,
a zasób i tak jest uszkodzony.

Przebieg próbny normalizacji nazw pokazuje wynik bez zapisu: nazwa jest tym,
po czym Operator odnajduje zasób, więc masowa zmiana bez podglądu byłaby
zmianą w ciemno.

Transliteracja znaków diakrytycznych rozkłada literę kanonicznie, zdejmuje
znak łączący i składa z powrotem: „ą" staje się „a", a „ł" zostaje, bo nie
jest literą ze znakiem łączącym — zamianę liter osobnych jak „ł" na „l"
rdzeń robi osobno, wprost.

## budowa/server/internal/core/adapter_modul_extension.go

Rozszerzenie jest pozycją katalogu platformy, którą ekspert dopiero bierze,
nie mostem MCP i nie konektorem eksperta. Most mieszka w punkcie dostępu,
z którego rdzeń składa listę serwerów MCP procesu modelu; pozycja rodzaju MCP
wskazuje most polem odwołania i tyle, drugiego rejestru serwerów MCP platforma
nie ma. Konektor eksperta należy do jednego eksperta, katalog stoi poziom
wyżej i do eksperta nie należy. Katalog jest więc warstwą nad tymi tabelami,
a nie ich kopią: ekspert bierze pozycję zakładając sobie odpowiedni konektor
albo umiejętność, a ten adapter tych tabel nie dotyka.

Instalacja zakłada pozycję w katalogu platformy i znaczy ją jako
zainstalowaną, nic więcej: nie sięga do sieci, nie pobiera paczki, nie
rozpakowuje archiwum, nie uruchamia procesu i nie zakłada punktu dostępu.
Szersze znaczenie zmieniłoby klasę bezpieczeństwa całego produktu, z programu
wykonującego wyłącznie własny kod na program wykonujący kod przyniesiony
z zewnątrz. Napis podany w polu źródła nie jest wyrzucany, ląduje w kolumnie
źródła zadeklarowanego jako zapis faktu, że taki adres podano; kształt
kontraktu pola źródła nie ma, więc nie wychodzi w odpowiedzi.

Odinstalowanie zdejmuje znaczniki, nie wiersz: gdyby kasowało pozycję, wykaz
katalogu nie miałby czego zawężać przełącznikiem instalacji, bo katalog
zawierałby wtedy wyłącznie pozycje zainstalowane. Odinstalowana pozycja
zostaje w katalogu ze swoją konfiguracją, a ponowna instalacja tym samym
kodem ją przywraca. Katalog jest rejestrem i nikt w rdzeniu go jeszcze nie
czyta: włączenie pozycji rodzaju MCP nie dokłada mostu do procesu modelu, bo
lista serwerów MCP składa się z punktów dostępu, a drugi pisarz tej listy
byłby drugą prawdą.

## budowa/server/internal/core/adapter_modul_badania_siec.go

Wywołanie sieciowe nie jest programem zewnętrznym: idzie biblioteką
standardową wkompilowaną w binarium, bez wywołania procesu i bez pomocnika
na dysku, więc zasada zabraniająca opierania funkcji o program spoza
instalki nie dotyczy sieci. Odmowa widoczna wtedy mówi o niedostępności
usługi, nie o braku na maszynie — to prawda o stanie świata, nie o brakach
instalacji.

Crossref, OpenAlex, arXiv, OpenLibrary i PubMed odpowiadają bez klucza, więc
wyszukiwanie naukowe działa od pierwszego uruchomienia, bez konfiguracji.
Dostawca wymagający klucza może dojść obok, ale nie jest warunkiem, żeby
moduł w ogóle wyszukiwał.

Plik nie rozstrzyga, co zrobić z wynikiem: nie zapisuje źródeł, nie liczy
przesiewu i nie zna kontraktu poza pozycją wyniku odkrycia. Oddaje pozycje
i błąd; decyzje należą do rodzin, które go wołają.

Limit czasu wywołania dostawcy wynosi dwadzieścia pięć sekund, bo dostawca
milczący dłużej jest dostawcą niedostępnym — czekanie bez granicy zawiesiłoby
okno na czas nieokreślony. Limit rozmiaru pobieranej treści chroni pamięć
przed plikiem, którego rozmiaru nikt nie zapowiedział. Nagłówek zgłaszający
moduł wchodzi w każde żądanie, bo Crossref prowadzi pulę grzecznościową dla
wywołań, które się przedstawiają, a wywołanie anonimowe bywa dławione.
Klient HTTP jest jeden na moduł, nie jeden na wywołanie, żeby korzystać
z jednej puli połączeń zamiast otwierać nowe gniazdo do tego samego
dostawcy przy każdej pozycji listy.

Adres przekierowania wyniku webowego rozwija się do adresu docelowego, bo
inaczej źródło zapisałoby adres pośrednika, który przestałby działać, gdy
pośrednik zmieni postać odnośnika.

## budowa/server/internal/core/adapter_modul_developer_wezly.go

Czynności plikowe idą przez bibliotekę standardową, a wyszukiwanie przez
wyrażenia regularne wkompilowane w rdzeń, nie przez program zewnętrzny:
wyszukiwanie w repozytorium jest czynnością, bez której Code Editor
przestaje być edytorem kodu, więc nie może zależeć od programu, którego
instalka nie niesie. Każda ścieżka przechodzi przez sprowadzenie do wnętrza
katalogu roboczego okna i odmawia, gdy z niego wychodzi — także wtedy, gdy
wyjściem jest dowiązanie symboliczne, nie zapis `..` w napisie.

Odpowiedź usunięcia węzłów nie niesie znacznika cofnięcia, i to jest stan
świadomy: cofnięcie wymagałoby odłożenia treści usuniętych węzłów
w magazynie rdzenia, a takiego magazynu moduł nie ma. Znacznik wypełniony
bez pokrycia obiecywałby Operatorowi powrót, którego nikt nie wykona.

Katalog niepusty schodzi przy usunięciu wyłącznie przy jawnym wskazaniu
rekurencji: usunięcie katalogu razem z zawartością na skutek pomyłki
w zaznaczeniu byłoby stratą, po której nie ma powrotu.

Kontrakt nie niesie osobnego pola „szukaj po składni", więc rozstrzyga sam
wzorzec, jego własną, udokumentowaną cechą: metazmienna `$NAZWA` jest
zapisem należącym do programu składniowego i nie znaczy nic w wyszukiwaniu
po napisie — napis `$ARG` jako napis szukany jest zapytaniem, którego nikt
nie zadaje. Wyrażenie regularne wyłącza tę drogę bezwarunkowo, bo w nim `$`
jest kotwicą końca wiersza, więc wzorzec zakończony tym znakiem byłby wzięty
za składniowy wbrew temu, co wołający napisał wprost.

Wiersze i kolumny programu składniowego liczą się od zera, a kontrakt od
jedynki — przeliczenie stoi w miejscu odczytu, żeby nie rozjechało się
między wyszukaniem a zamianą.

Zamiana masowa w repozytorium dotyka wielu plików naraz i nie ma po niej
cofnięcia, więc odpowiedź mówi wprost, czy zmiany zapisano, a wykaz zmian
wraca także przy podglądzie bez zapisu — Operator ma zobaczyć, co się
stanie, zanim to się stanie.

Przejścia zamiany po składni są dwa i jest to rozstrzygnięcie celowe:
pierwsze, bez zapisu, daje wykaz zmian, który wraca w odpowiedzi także przy
zapisie. Drugie zapisuje. Wyprowadzenie wykazu z samego zapisu nie da się
zrobić, bo przy zapisie masowym program oddaje liczbę zmian, nie ich treść.

Brak programu składniowego na maszynie kończy się odmową, nie cichym
powrotem do drogi napisu: wzorzec składniowy szukany jako zwykły napis nie
znajdzie niczego, a zero trafień wyglądałoby jak odpowiedź prawidłowa, nie
jak brak programu — wynik wyglądający dobrze jest groźniejszy od odmowy,
bo nie wzywa do sprawdzenia.

Wykaz pusty przy przejściu po składni jest wynikiem prawidłowym; wykaz
nieczytelny przy niepowodzeniu programu nie jest, bo wtedy nie zmierzono
niczego i odpowiedź ma to powiedzieć odmową, nie pustym wykazem.

## budowa/server/internal/core/adapter_modul_tlumaczenie_pamiec_wykaz.go

Para wniesiona ręcznie i para z pliku wymiany nie mają panelu, a usunięcie
panelu nie zabiera ze sobą par, które z niego kiedyś zdjęto: klucz obcy do
panelu jest opcjonalny i wraca do wartości pustej, nie kasuje wiersza pamięci.

Plik CSV o kolumnach segmentu źródłowego, segmentu docelowego i języka jest
drugą drogą wymiany, bo tyle właśnie eksportuje większość arkuszy, w których
prowadzi się terminologię poza rdzeniem.

Pełny standard TMX niesie nadto nagłówek z metrykami narzędzia i notami;
rdzeń przy zapisie wypisuje nagłówek minimalny, a przy odczycie nie wymaga
niczego ponad jednostkę tłumaczeniową, bo tyle właśnie niesie para pamięci.

## budowa/server/internal/core/adapter_modul_badania_odkrycia.go

Każda pozycja zwrócona przez dostawcę ląduje w bazie pod swoim kluczem: bez
tego odrzucenie pozycji nie miałoby czego odrzucić, a liczniki diagramu
PRISMA byłyby liczbami wymyślonymi w chwili odpytania. Zapamiętanie znaczy
też, że pozycja już dodana do źródeł jest przy powtórnym wyszukaniu
oznaczona identyfikatorem źródła, zamiast wchodzić drugi raz jako nowa.

Wyszukiwanie chodzi po kilku dostawcach naraz. Milczenie jednego z nich
wchodzi do wykazu dostawców, które zawiodły, a pozostali oddają swoje
trafienia — odmowa całości przez jeden zerwany strumień byłaby karą za
cudzą awarię.

Wybór dostawców wyszukiwania: wskazanie wprost ma pierwszeństwo — Operator,
który wskazał dostawcę, dostaje jego wyniki, nie zestaw domyślny obok nich.

Partia importu wsadowego dostaje własny identyfikator, który wchodzi w pole
pochodzenia każdego źródła z niej pozyskanego. Dzięki temu identyfikator
partii w odpowiedzi jest wskazaniem, po którym da się odnaleźć skutek
partii w bazie, nie numerem zadania, które nigdzie nie stoi.

## budowa/server/internal/core/adapter_modul_tlumaczenie.go

Rozpoznanie źródła rozpoznaje język modelem: rdzeń nie ma własnego silnika
rozpoznania, tylko most do rejestru kanałów, który pyta domyślny czynny
kanał, jaki to język. Bez wpiętego rejestru albo bez czynnego kanału odmawia
wprost, zamiast zgadywać po znakach diakrytycznych.

Zapis źródła nie wytwarza treści tłumaczenia — zapisuje sam tekst źródłowy
i liczbę segmentów. Treść panelu powstaje dopiero przy dodaniu panelu, gdzie
znany jest język docelowy; bez czynnego kanału dodanie panelu odmawia
zamiast zakładać panel pusty.

Segmentacja tekstu idzie podziałem własnym, bez biblioteki zewnętrznej: tekst
dzieli się na zdania po znaku końca zdania, po którym stoi biały znak albo
koniec tekstu. Skróty w rodzaju „np." albo „ul." rozłamią zdanie tam, gdzie
językoznawczo się nie kończy — podział semantyczny wymagałby słownika
skrótów albo modelu, których rdzeń nie ma. Komenda podziału na segmenty nie
zapisuje nic do bazy: segmenty nie mają własnej tabeli, więc powtórne
wywołanie na tym samym tekście daje ten sam wynik bez efektu ubocznego.

Pole pewności rozpoznania języka zostaje puste: model oddaje nazwę języka,
nie miarę pewności, a rdzeń nie dorabia liczby, której nikt uczciwie nie
wypełni.

Niezgodności panelu nie są polem repozytorium panelu — wypełnia je odczyt
obszaru jakości, tak samo jak ustalenia badania nie są polem źródła badania.
Kod okna w bycie kontraktu panelu jest kodem zewnętrznym, doczytanym
złączeniem w warstwie danych, nie numerem wiersza wewnętrznego, bo numer
wiersza nie zaadresowałby po stronie klienta żadnego okna.

## budowa/server/internal/core/adapter_modul_isolation_profile.go

Profil jest szablonem: wartość obowiązująca leży w tabeli ustawień
i wyłącznie stamtąd czyta ją rozstrzygacz, więc zapisanie profilu nie zmienia
izolacji — skutek daje dopiero przypisanie profilu, które przepisuje jego
przełączniki pod wskazany adres. Punkt nieujęty w profilu zostaje wtedy
nietknięty: obowiązuje to, co na poziomie stoi, a w ostateczności wartość
domyślna. Profil bez ani jednego przełącznika jest odmawiany, bo przypisanie
nie zmieniłoby żadnej wartości, a odpowiedź niosłaby politykę wyliczoną tak
samo jak przed wywołaniem — potwierdzenie czynności, której nie było.

Usunięcie profilu nie cofa wartości już przepisanych: przełączniki, które
przypisanie wniosło do poziomu, są odtąd wartościami tego poziomu i cofa je
osobne ustawienie kontekstowe albo techniczne. Profil, którego nie ma, wraca
odmową nieznalezienia, a nie odpowiedzią mówiącą, że nic nie usunięto.

Wybór warstwy zmienia wyłącznie sam wybór, utrwalony osobno od wartości
izolacji, więc panel otwarty ponownie pokazuje warstwę zapisaną poprzednio,
a żadna wartość izolacji się nie zmienia. Wybór warstwy należy do bytu, dla
którego zapadł, więc żądanie bez karty sesji ani okna komunikacji wraca
odmową: wiersz bez bytu nie zostałby odczytany.

Podgląd polityki niczego nie zapisuje: liczy politykę po ośmiu poziomach
zasięgu i oddaje ją wraz ze wskazaniem, skąd wartości pochodzą. Poziom
wskazany wprost obowiązuje bez zmian; pominięty znaczy poziom najwęższy
z bytów podanych w żądaniu, czyli okno komunikacji przed kartą sesji,
a przy braku obu — poziom globalny. Warstwa wskazana wprost obowiązuje bez
zmian; pominięta znaczy warstwę zapisaną wcześniej dla tego bytu, a przy
braku zapisu — warstwę domyślną platformy.

Liczenie polityki obowiązującej rozstrzyga pakiet konfiguracji, nie ten
plik: tu składany jest wyłącznie kontekst, czyli byt na właściwym poziomie,
a wynik przekładany na kształt kontraktu. Błąd odczytu kończy podgląd
odmową, mimo że pakiet konfiguracji oddaje politykę mimo błędu źródła dla
samego wykonania — podgląd ma odpowiadać na pytanie, co obowiązuje, więc
pokazanie wartości domyślnych jako obowiązujących, gdy zapisów nie udało się
odczytać, byłoby odpowiedzią nieprawdziwą.

Kontrakt ma na pochodzenie polityki jedno pole, a każdy z jedenastu punktów
bywa zapisany na innym poziomie, więc poziom wraca tylko wtedy, gdy
wszystkie zapisane punkty pochodzą z tego samego poziomu — przy polityce
złożonej z kilku poziomów pole zostaje puste, bo jedna nazwa byłaby
wskazaniem nieprawdziwym dla pozostałych punktów.

## budowa/server/internal/core/adapter_modul_agents_zakres.go

Zawężenie zapisane jest zawężeniem egzekwowanym: zapis, którego rdzeń nie
czyta przy wykonaniu, byłby gorszy niż jego brak, bo dawałby wrażenie
ograniczenia, którego nikt nie pilnuje. Wyliczenie polityki efektywnej
niczego nie rozstrzyga i niczego nie zapisuje, jest wyłącznie
przezroczystością stanu. Dziedziczenie z poziomów zasięgu bierze ten sam
rozstrzygacz, którym jedzie okno konfiguracji punktów izolacji technicznej —
drugiej macierzy izolacji nie ma.

Zero w kolumnie granicy podagentów znaczy wyłączony i jest jedynym zapisem
wyłączenia: włączenie bez podanej granicy przywraca piętnaście, maksimum
techniczne platformy tego samego rzędu wielkości, jakie zna uruchomienie
podagenta.

## adapter_modul_aplikacje_srodowiska.go

Kontrakt nie ma komendy zakładającej środowisko wdrożeniowe — `AppEnvironment`
wychodzi wyłącznie z `apps.environment.list` — a selektor środowiska Deployment
Panelu musi mieć co pokazać w oknie świeżo otwartym. Wykaz zakłada więc
brakujące wiersze dla trzech wartości `AppDeployEnvironment` i oddaje je wraz
z tym, co Operator zdążył im nadać. Zakładanie jest operacją UPSERT
nietykającą domeny, więc kolejne wywołanie wykazu nie kasuje adresu nadanego
między jednym wywołaniem a drugim.

`AppsDeploymentDomainSetResponse` niesie pole `verified`; jedyną uczciwą jego
treścią jest wynik rozwiązania nazwy przez system. Nazwa nierozwiązywalna nie
jest odmową — domena bywa nadawana, zanim wpisy DNS się rozejdą — ale
`verified: false` mówi Operatorowi wprost, że jeszcze nie działa.

`apps.deployment.health.get` sprawdza produkt naprawdę: gdy środowisko ma
domenę albo w oknie stoi podgląd, idzie tam zapytanie HTTP; gdy nie ma dokąd
pójść, dostępność wynika ze stanu ostatniego przebiegu wdrożenia. Każde
sprawdzenie dopisuje wiersz do dziennika kondycji, a udział dostępności liczy
się z tych wierszy, zamiast być liczbą wziętą znikąd.

## budowa/server/internal/core/adapter_modul_developer_wynik.go

Wynik testów rozbiera się w chwili biegu przebiegu, nie w chwili pytania
o niego, ponieważ dziennik przebiegu przechowuje wyłącznie ogon logu: rdzeń
przycina go do kilkuset wierszy, bo budowanie dużego projektu daje ich
dziesiątki tysięcy. Gdyby wynik testu powstawał z odczytu dziennika, przebieg
z tysiącem testów oddałby ich tylko kilkadziesiąt, a odpowiedź wyglądałaby
przy tym poprawnie. Dlatego wiersze niosące wynik testu zbiera się na
bieżąco, podczas gdy log jeszcze płynie, a złożony wynik zapisuje się do
bazy przy domknięciu przebiegu.

Pokrycie kodu pochodzi z dwóch źródeł w ustalonej kolejności. Gdy zadanie
budowania niosło parametr profilu pokrycia, rdzeń czyta ten plik i uzyskuje
pomiar co do instrukcji i wiersza, który zasila nakładkę pokrycia
w edytorze. Gdy profilu nie ma, zostaje wiersz podsumowania testu
z procentem pokrycia pakietu, który daje wynik uboższy: zna procent, a nie
wiersze niepokryte. Drugie źródło jest uboższe, lecz prawdziwe: przyjęcie,
że bez profilu wiadomo, które wiersze są niepokryte, byłoby zmyśleniem.

Zapis wyniku przebiegu bez testów nie zakłada ani jednego wiersza: zapisanie
pustego pomiaru pokazywałoby w oknie Build Output podsumowanie zerowych
testów tam, gdzie testów nikt nie uruchamiał, co nie jest tym samym, co
informacja, że testy przeszły.

## adapter_modul_library_slownik.go

Etykieta ma tu tożsamość, której `library.tag.set` jej nie daje: barwę, czas
założenia i istnienie niezależne od tego, czy nosi ją jakikolwiek zasób. Dzięki
temu okno Tags & Collections pokazuje słownik, a nie próbkę zebraną z odczytanej
strony wykazu.

Wywóz tezaurusa idzie w SKOS/RDF w trzech serializacjach i powstaje w rdzeniu,
bo relacje leżą w bazie i klient nie ma ich skąd wziąć. Zapis składa się
z tekstu, bez programu zewnętrznego: RDF w postaci Turtle, RDF/XML i JSON-LD
to formaty tekstowe o znanym kształcie.

Łączenie etykiet jest zmianą nazwy wykonaną wielokrotnie: zasób noszący
etykietę źródłową dostaje docelową, a źródłowa znika ze słownika. Etykieta
wskazana jako źródłowa i docelowa naraz jest pomijana — wchłonięcie siebie
samej zdjęłoby etykietę z zasobów bez powodu.

Relacja tezaurusa zapisuje się w jednym kierunku, tym wskazanym przez
Operatora. Odwrotność wyprowadza odczyt — pojęcie nadrzędne czytane od drugiej
strony jest podrzędnym — a zapis obu kierunków dałby dwa wiersze mówiące
to samo i rozjazd przy zdjęciu jednego z nich.
## budowa/server/internal/core/adapter_modul_auth.go

Bramka jest jedna, a Operator bezimienny: żądania rodziny `auth` nie niosą ani
nazwy, ani adresu e-mail, a `auth.register` nie zakłada konta w sensie
katalogowym — ustawia sekret bramki przy pierwszym uruchomieniu i od razu
wpuszcza. Konto z `migracja_014_katalog_kont.sql` jest poświadczeniem do
kanału modelu i z bramką nie ma nic wspólnego.

Dwanaście godzin trwania sesji bramki to jedna doba robocza: Operator, który
rano wszedł, nie loguje się w połowie dnia, a maszyna zostawiona na noc bramkę
zamyka. Rok dla sesji z zaznaczonym „nie wyloguj mnie” to długość, po której
zapomniane urządzenie przestaje wchodzić samo, a Operator pracujący codziennie
nie zobaczy okna logowania ani razu, bo wygasanie jest przesuwne i każde
wejście je odnawia.

`ZalozBramke` poczty nie wymaga: na świeżej instalacji konta nadawczego
platformy nie ma jeszcze czym wskazać, a ustawia się je w oknie Konfiguracji —
za bramką. Rejestracja idzie więc dwiema drogami: z pocztą nadaje list z drogą
potwierdzenia i zostawia konto niepotwierdzone, bez poczty zakłada konto
i zapamiętuje w sejfie, że adresu nikt nie potwierdził. W obu razach hasło
otwiera bramkę od razu; sesji rejestracja nie zakłada — konto powstaje
niepotwierdzone i pozostaje w tym stanie do `auth.verify`, bo wpuszczanie od
razu czyniłoby weryfikację adresu ozdobą. Czynność jest wykonalna tylko raz:
istniejąca kotwica daje `conflict`, nie ciche `registered: false`, a powtórzone
żądanie jest próbą podmiany hasła bez znajomości starego, do czego służy
odzyskanie konta. List wysyła się przed oddaniem odpowiedzi i jego
niepowodzenie schodzi na drogę bez poczty, zamiast cofać rejestrację — cofanie
było tu wcześniej ratunkiem przed platformą nie do otwarcia, a stało się
zbędne i ryzykowne, odkąd konto bez potwierdzonego adresu wchodzi hasłem:
literówka w nastawach nadajnika zamykałaby pierwsze uruchomienie równie
szczelnie jak brak poczty w ogóle. `cofnijRejestracje` nie kasuje drogi
potwierdzenia: leży jako sam skrót, wygasa po godzinie, a bez konta nie ma
czego otworzyć.

`WejdzPrzezBramke` odmawia kontu niepotwierdzonemu, bo adres jest jedyną drogą
odzyskania dostępu i musi być sprawdzony, zanim się nią stanie — inaczej
literówka w adresie wyszłaby na jaw dopiero w dniu, w którym trzeba nim
odzyskać konto. Sekret niezgodny daje `not_authenticated`, nie
`validation_failed`: ten drugi kod zostaje przy brakach kształtu żądania, żeby
sonda stanu bramki po stronie klienta, wysyłająca żądanie bez sekretu, dalej
czytała go jako „bramka ustawiona”. Urządzenie sesji bierze się z żądania,
a wiersz metody uzupełnia je tylko wtedy, gdy żądanie milczy: kotwica hasła nie
jest materiałem jednej maszyny, więc branie urządzenia wyłącznie z wiersza
metody porzucało `deviceId` przy każdym wejściu hasłem, maszyna nie pojawiała
się w wykazie `device.list`, a `device.revoke` z jej identyfikatorem wracał
`revoked: false`, zostawiając token czynny.

`zamekZmiany` szereguje czynności, które sprawdzają stan bramki i zaraz potem
go zmieniają. Bez niego sprawdzenie i zapis są dwiema czynnościami, a między
nie wchodzi drugie żądanie z tego samego gniazda — dwie równoległe zmiany
hasła odpowiadają wtedy obie `changed: true`, a bramkę otwiera tylko jedno
z dwóch nowych haseł, bo drugi zapis do sejfu nadpisuje pierwszy. Sejf jest
plikiem z wpisami i transakcji nie zna, więc niepodzielność musi stanąć tutaj.

## budowa/server/internal/core/adapter_modul_wiedza.go

Wyszukiwanie po słowach działa w module Library indeksem pełnotekstowym. Ta
rodzina komend wnosi wyszukiwanie po znaczeniu, a różnica leży w tym, czego
tamto nie umie: pytanie o postępowanie przy awarii maszyny nie ma z dokumentem
opisującym awarię węzła ani jednego wspólnego słowa, więc indeks liter go nie
znajdzie.

Pomocnik osadzeń startuje tym samym uruchamiaczem i przez tę samą bramę
izolacji co każdy inny proces drzewa, więc potrzebuje trójki okno, zasady i
obszar zasięgu platformy. Żądania rodziny `knowledge.*` niosą co najwyżej
identyfikator okna, i to po to, żeby wskazać przestrzeń do przeszukania, nie
po to, żeby w tym oknie liczyć. Wskaźnik jest jeden na maszynę i wspólny dla
wszystkich okien, więc trójka składa się w zasięgu platformy, tą samą drogą co
w silniku mowy: pusty kontekst zasięgu jest poprawnym adresem najszerszego
z poziomów, a nie podstawieniem pustych struktur po cichu.

Silnik osadzeń powstaje na każde wywołanie, tak samo jak w module mowy: metoda
ustawień mutuje byt, więc jedna instancja współdzielona przez równoległe
żądania oznaczałaby wyścig o nastawy. Nastawy są przy tym świeże — zmiana
modelu wiedzy wchodzi w następnym przebiegu, bez restartu rdzenia.

Składanie trwałości wskaźnika nad bazą montażu stoi w osobnej funkcji, a nie
w wyrażeniu w miejscu wpięcia: montaż bez bazy jest stanem, który zdarza się
przy rdzeniu składanym do sprawdzenia transportu, a wartość pusta przechodzi
tędy bez warunku po stronie wołającego.

Wnoszenie dokumentów do wskaźnika idzie dokument po dokumencie, a nie
wszystko naraz w jednej transakcji: biblioteka Operatora bywa gigabajtem
tekstu, a jedna transakcja na całość znaczyłaby komplet wektorów w pamięci
rdzenia naraz. Kasowanie poprzednich fragmentów źródła przed zapisem nowych
zapobiega temu, żeby fragmenty treści skróconej od poprzedniego przebiegu
zostały i wracały jako cytat z dokumentu, w którym ich już nie ma.

Wyszukiwanie oddaje wynik pusty jako odpowiedź, nie jako odmowę: wskaźnik
pusty albo wiedza bez związku z pytaniem znaczą brak trafień i model ma to
usłyszeć wprost. Odmiennie brak silnika, który jest odmową, ponieważ wtedy
rdzeń nie wie, czy wskaźnik coś ma, czy nie ma.

Pole przesiewu dokłada drugi przebieg wyszukiwania. Pierwszy przebieg zostaje
niezmieniony i wykonuje się zawsze: przesiew nie zastępuje podobieństwa
kosinusowego, tylko układa na nowo tych kandydatów, których podobieństwo
kosinusowe wybrało, bo krzyżowym koderem nie da się przejrzeć całego
wskaźnika. Kandydatów jest przy tym więcej niż fragmentów oddawanych
w odpowiedzi: przesiew może wynieść na czoło fragment, który po samym
podobieństwie wektorów był daleko za progiem odpowiedzi, a gdyby kandydatów
było dokładnie tyle, ile fragmentów wraca, przesiew przestawiałby wyłącznie
kolejność wewnątrz zbioru już wybranego. Odmowa przesiewu jest odmową całego
żądania, a nie cichym zejściem na wynik pierwszego przebiegu, ponieważ
wołający prosił o kolejność ułożoną na nowo.

Znakowanie odmów pakietu wiedzy trzyma trzy gałęzie: brak silnika osadzeń daje
kod niedostępności kanału, w praktyce nieponawialny mimo ponawialności samego
kodu, bo zaplecze liczenia jest niedostępne i to samo żądanie powiedzie się
bez zmiany dopiero po naprawie opisanej treścią odmowy; naruszenie izolacji
daje kod odmowy uprawnień, tak samo jak znakuje je moduł Terminal, żeby dwie
reguły dla jednej bramy nie rozjechały się; każda pozostała odmowa dostaje kod
błędu wewnętrznego jako najostrzejszy z zamysłem, bo nieznana odmowa jest
przypadkiem, którego rdzeń nie przewidział.

## budowa/server/internal/core/adapter_modul_tlumaczenie_napisy.go

Cztery formaty napisów kontraktu czyta i pisze ten rdzeń sam: SRT i WebVTT są
tekstowe, TTML jest zapisem XML, a EBU STL jest zapisem dwójkowym o stałej
ramce, złożonym z bloku nagłówkowego GSI liczącego tysiąc dwadzieścia
cztery bajty oraz bloków tekstowych TTI po sto dwadzieścia osiem bajtów.
Wszystkie cztery formaty rozbiera się i składa bibliotekami wkompilowanymi
w rdzeń, bez wywołania programu zewnętrznego.

Kwestie napisów są trwałe, zapisane w osobnej tabeli, ponieważ sprawdzenie
taktowania napisów i złożenie scenariusza dubbingu nie miałyby czego mierzyć
bez zapisanego przebiegu: taktowanie jest własnością materiału, a nie tekstu
panelu, więc musi przeżyć poza pojedynczym wywołaniem eksportu.
## budowa/server/internal/core/adapter_modul_tlumaczenie_wymiana_zewnetrzna.go

Plik obsługuje `translate.handoff.build`, `.receive`,
`translate.bridge.source.receive`, `translate.bridge.result.send`,
`translate.artifact.publish` i `translate.step.list`. Pakiet przekazania
niesie dokładnie te zawartości, o które prosi żądanie: XLIFF z jednostkami,
TMX z pamięcią, TBX z terminologią i plik instrukcji.

`artifact.publish` odkłada wytwór do magazynu treści rdzenia i zakłada wiersz
pliku biblioteki. Jest to jedyna droga, którą wytwór Translate staje się
widoczny dla reszty platformy — bez wiersza plik leżałby na dysku bez jednego
bytu, który by o nim wiedział.

Wykaz kroków `translate.step.list` stoi w kodzie, a nie w bazie, bo opisuje
zdolności rdzenia, nie dane Operatora: krok istnieje dokładnie wtedy, gdy
istnieje obsługująca go komenda, i znika razem z nią. Nazwy komend biorą się
ze stałych kontraktu, więc wykaz nie ma jak rozjechać się z rejestrem.

## adapter_modul_role.go

Rodzina `role.*` jest fasadą, nie drugą prawdą o roli. Rola okna ma w rdzeniu
jednego właściciela: pakiet `session` (rejestr okien i `rola_okna.go`, który
zna słownik ról i normalizuje więź), a jej ślad trwały — kolumny
`okno_komunikacji.rola_okna` i `okno_komunikacji.okno_koordynatora_id`
z `migracja_002_okna.sql`. Rodzina `role.*` jedzie dokładnie na nich, tak samo
jak `window.update`. Gdyby założyła własny zapis roli, produkt miałby dwie
odpowiedzi na pytanie, jaką rolę ma to okno — jedną z `window.list`, drugą
z `role.*`.

Wobec `window.update` rodzina dokłada trzy rzeczy. Pierwsza to wcielenie
(`persona`), którego `window.update` nie zna wcale. Druga to odmowa zamiast
cichego pominięcia: `window.update` ze wskazaniem koordynatora dla okna, które
wykonawcą nie jest, cicho zdejmuje wskazanie, a rodzina `role.*` odmawia
z powodem, bo jej jedynym tematem jest właśnie rola i jej więź. Trzecia to
ślad trwały — `window.update` zmienia wyłącznie rejestr pamięciowy, a `role.*`
zapisuje rolę także do wiersza okna, gdy wiersz istnieje, bez czego
`window.state.get` po restarcie rdzenia oddawałby rolę sprzed nadania.

Wcielenie mieszka tam, gdzie już mieszka. Klient utrwala wcielenie okna
komendą `config.set` na poziomie zasięgu okna pod kluczem
`multitasking.wcielenie`. Rdzeń pisze i czyta ten sam adres, więc wcielenie
nadane komendą `role.update` widzi selektor analityka i odwrotnie. Własny
klucz albo własna tabela byłyby drugim wcieleniem tego samego okna.

Drogi zapisu roli są dwie, bo okno ma dwa życia. Okno otwarte w tym
uruchomieniu rdzenia stoi w rejestrze pamięciowym i to on jest jego prawdą
bieżącą; okno sprzed restartu jest wyłącznie wierszem. Rodzina `role.*`
obsługuje oba, zamiast odmawiać oknu, o którym Operator wie z
`window.state.get`, że istnieje. Wskazania niewypełnione — rola i koordynator
oba puste — niczego nie zmieniają: obie drogi sprowadzają się wtedy do odczytu
stanu obowiązującego.

Pomocnik `wskaznikPolaRoli` jest własny, a nie wspólny, bo pakiet niesie dwa
pomocniki o nazwie zbliżonej i o różnym znaczeniu pustego napisu — jeden
oddaje brak pola, drugi wskaźnik na pustkę. Rodzina `role.*` trzyma się
zasady: koordynator pusty i wcielenie puste są brakiem pola.

Więź z oknem, którego nie ma, byłaby potwierdzeniem relacji, która nie
powstała. Rejestr pamięciowy sprawdza to sam, a na drodze wiersza sprawdza to
`nadajRoleWWierszu`, bo tam nikt inny tego nie robi.

## budowa/server/internal/core/adapter_modul_badania_eksport.go

Eksport wytwarza plik, a nie tylko wiersz o nim: zapisuje bajty pod sumą
kontrolną we wspólnym magazynie rdzenia i oddaje rzeczywistą ścieżkę pliku,
żeby odpowiedź o powodzeniu nigdy nie wracała z pustą ścieżką.

Formaty tekstowe składa rdzeń sam. Markdown, HTML, tekst i LaTeX powstają
w Go, bez ani jednego procesu zewnętrznego, więc działają zawsze. Pliki PDF,
DOCX i PPTX idą przez port arsenału dokumentowego, ponieważ składu tych
formatów nie da się napisać od nowa uczciwiej niż dojrzałym programem, który
stoi na serwerze razem z rdzeniem. Plik XLSX powstaje w tym module wprost,
biblioteką archiwizującą wkompilowaną w rdzeń: arkusz jest spakowanym
dokumentem XML, a nie składem, więc osobny program go nie wymaga.

## budowa/server/internal/core/adapter_modul_library.go

Wyszukiwanie `library.file.search` dopasowuje frazę do nazwy pliku albo do
jego treści, tę drugą przez indeks pełnotekstowy zasilany przy każdym zapisie
treści. Bajty leżą poza bazą, indeks jest ich odtwarzalnym wyciągiem
tekstowym; dopasowanie jest trafieniem w słowo, nie w znaczenie. Wykaz
`library.file.list` zawęża wynik po nazwie i nie szuka w dokumentach.

Katalog danych magazynu treści jest przy konstrukcji adaptera domyślny;
katalog obowiązujący na uruchomieniu zna wyłącznie montaż i podaje go osobnym
wywołaniem przy składaniu portu. Wartość domyślna stoi w konstruktorze dla
wywołania bez montażu, żeby konstruktor nigdy nie oddał adaptera bez
magazynu. Sprzątanie magazynu treści wypada przy tym samym wywołaniu, w
jedynej chwili startu rdzenia, w której moduł zna już swój prawdziwy katalog
danych i jeszcze nie obsługuje żadnej komendy; sprzątanie idzie synchronicznie,
bo obchód katalogu jest tani wobec odtworzenia stanu, a rdzeń przyjmujący
wgrania w trakcie przemiatania widziałby wykaz żywych odwołań sprzed nich.

Obie drogi wgrania pliku prowadzą do jednego magazynu: treść przysłana base64
i treść wciągnięta spod ścieżki źródłowej lądują w magazynie treści rdzenia,
a odwołaniem jest ścieżka bloba. Ścieżka źródłowa zostaje przy pliku w
osobnej kolumnie jako prowenancja — mówi, skąd plik przyszedł, i nie jest
wskaźnikiem na treść żywą, którą ktoś z zewnątrz mógłby nadpisać po wgraniu.

Kolekcje pliku czyta się z bazy, nie z żądania, żeby każda odpowiedź niosąca
plik mówiła tę samą przynależność, niezależnie od drogi, którą przyszła.
Ścieżka WEWNĄTRZ repozytorium, którą nadaje przenoszenie pliku, jest polem
osobnym od ścieżki źródłowej z maszyny Operatora: ta druga z rdzenia nie
wychodzi, bo wyniosłaby na zewnątrz układ cudzego dysku wraz z nazwami
katalogów.

## budowa/server/internal/core/adapter_modul_tlumaczenie_lokalizacja.go

Formaty zasobów — JSON, YAML, properties, Android XML, iOS strings
i stringsdict, RESX, gettext PO — czyta i pisze ten rdzeń sam, bibliotekami
wkompilowanymi, bez wywołania programu zewnętrznego.

Formy mnogie idą regułami CLDR. Reguła mnogości jest własnością języka, nie
tłumaczenia: polski ma trzy formy, angielski dwie, czeski trzy, rosyjski
trzy, arabski sześć. Zastosowanie form zakłada klucze wariantów, których
język docelowy wymaga, i zdejmuje te, których nie zna, bo inaczej plik
wyniku miałby formę, której język nie ma, i program lokalizowany nigdy by
jej nie użył.

## budowa/server/internal/core/adapter_modul_aplikacje_warsztaty.go

Trasy, punkty końcowe i schemat są odczytane z pracy, nie zapisane obok
niej. Kontrakt nie daje komend zapisu tych bytów — trzy komendy odczytu
stoją w rodzinie same. Nie jest to przeoczenie kontraktu, tylko jego
rozstrzygnięcie: trasa, punkt końcowy i tabela produktu są w tym, co
Operator napisał, a osobna tabela byłaby drugą prawdą, rozjeżdżającą się
z kodem przy pierwszej edycji, która zapomni ją odświeżyć.

Trasy czyta się z plików warstwy frontendu, z deklaracji ścieżki oraz
z atrybutu ścieżki znacznika trasy. Punkty końcowe pochodzą z kontraktów
API komponentów architektury, wypełnianych w panelu kontraktu API. Schemat
pochodzi z poleceń tworzenia tabeli w plikach warstwy backendu. Każdy z tych
trzech odczytów daje wynik pusty, gdy Operator jeszcze niczego nie napisał,
i to jest odpowiedź prawdziwa, nie brak.

Zapytanie próbne idzie po sieci naprawdę: składa żądanie pod adres
środowiska, a gdy domeny nie ma, pod adres stojącego podglądu, i mierzy
czas oraz kod odpowiedzi zegarem, nie zgadywaniem. Brak adresu i brak
podglądu razem znaczy odmowę z powodem, a nie wynik powodzenia wzięty
znikąd.

## budowa/server/internal/core/adapter_modul_terminal_analiza.go

Mapa programów analizy dla poszczególnych powłok jest nieunikniona, ponieważ
rozpoznanie składni Basha, PowerShella i Pythona to trzy odrębne, dojrzałe
programy, których żaden rozsądny nakład pracy nie przepisze do Go. Mapa jest
POMIERZONA, nie zgadnięta: bash ma ShellCheck do analizy i shfmt do
formatowania, oba zmierzone i zainstalowane na maszynie rdzenia; powershell
ma PSScriptAnalyzer wywoływany przez `Invoke-ScriptAnalyzer`
i `Invoke-Formatter`, również zmierzony i zainstalowany; node ma wyłącznie
`node --check`, czyli sam interpreter, którego karta node i tak wymaga, więc
formatowania dla tej powłoki nie ma i odpowiedź go nie obiecuje; python ma
Ruff (`ruff check`, `ruff format`), a gdy Ruffa na maszynie nie ma, zostaje
`python -m py_compile` — sam interpreter, orzekający wtedy wyłącznie
o składni, nie o regułach, i odpowiedź nazywa w takim wypadku interpreter,
nie Ruffa. Cmd i ssh nie mają pozycji: dla wsadu cmd nie istnieje powszechnie
przyjęty analizator, a ssh nie jest językiem, tylko transportem. Wykaz
zawiera to, co zmierzono, i nic ponad to; powłoka spoza wykazu dostaje
odpowiedź mówiącą wprost, że analizatora dla niej nie ma, ponieważ pusty
wykaz uwag znaczyłby fałszywie, że treść jest bez zastrzeżeń.

Treść skryptu trafia do programu analizy przez plik tymczasowy, nie przez
strumień wejścia procesu, z dwóch niezależnych powodów. Każdy z programów
analizy czyta plik; ShellCheck umie czytać też strumień, ale wtedy gubi
nazwę pliku w uwagach. Jedyna droga uruchomienia rdzenia do procesu
zewnętrznego (`zewnetrzne.Wolaj`) nie pisze na wejście procesu z zamysłu,
ponieważ jednoczesne pisanie na wejście i czytanie wyjścia jest klasą
zakleszczeń, której ten pakiet ma nie mieć. Ta sama droga sprawia, że Ruff
formatuje plik, a nie strumień, mimo że `ruff format -` dałby ten sam
wynik: droga rdzenia do procesu z zamysłu nie pisze na jego wejście, więc
czytanie idzie z pliku, który po analizie znika z katalogu tymczasowego
maszyny rdzenia.

Python ma dwa programy, a pierwszeństwo należy do Ruffa: `python -m
py_compile` orzeka wyłącznie o składni i jest wobec Pythona tym, czym
ShellCheck jest wobec basha i PSScriptAnalyzer wobec PowerShella dla
pozostałych powłok — orzeczeniem o składni oraz o regułach. Ruff obejmuje
jedno i drugie: błąd składni wraca z niego jako uwaga `invalid-syntax`, więc
pierwszeństwo Ruffa niczego nie odbiera. Gdy Ruffa na maszynie nie ma,
zostaje interpreter, i wtedy odpowiedź nazywa interpreter, bo to on
naprawdę sprawdzał treść — nazwa programu w odpowiedzi ma zgadzać się z tym,
co ją wystawiło.

Żądanie analizy nie niesie okna, a każdy proces rdzenia ma mieć obszar
i zasady izolacji, więc SprawdzSkrypt bierze obszar pusty i zasady puste —
rozstrzygnięcie, nie przeoczenie: program analizy czyta wyłącznie plik
tymczasowy założony przez rdzeń, nie sięga do obszaru żadnego okna i nie ma
z niego czego wynieść. Kod wyjścia różny od zera jest tu wynikiem analizy,
nie usterką: ShellCheck i `node --check` kończą pracę niezerowo dokładnie
wtedy, gdy mają co zgłosić; odpowiedź odmawia wyłącznie wtedy, gdy programu
nie ma czym uruchomić. Formatowanie nieudane nie unieważnia analizy z tego
samego powodu, dla którego formatowanie jest osobnym krokiem: uwagi są tym,
po co komenda powstała, a treść sformatowana jest dodatkiem — gdy się nie
uda, pole treści sformatowanej zostaje nieobecne, zgodnie z kontraktem.

## budowa/server/internal/core/adapter_modul_przegladarka_uchwyty.go

Dodanie źródła (`browser.source.add`) i dodanie notatki (`browser.note.add`)
nie rozgłaszają zdarzenia domenowego, w odróżnieniu od nawigacji. Kontrakt nie
przewiduje dla nich zdarzenia zmiany (nie ma odpowiednika `browser.source.changed`
ani podobnego dla notatek), więc port odkłada wynik wyłącznie do odpowiedzi
komendy, zamiast wprowadzać zdarzenie, którego kontrakt nie niesie.

Rozgłoszenie migawki strony po nawigacji jedzie tym samym emiterem rdzenia,
co pozostałe zmiany obszarów aplikacji, zachowując jeden wspólny wzorzec
rozgłaszania zdarzeń niezależnie od modułu źródłowego.

Otwarcie i zmiana stanu karty rozgłaszają zdarzenie zmiany karty z rodzajem
zmiany zapisanym w polu kontraktu `ChangeKind`, a nie w dowolnym napisie —
klient odróżnia po tym polu założenie karty od zmiany jej stanu.

Zdarzenie zmiany monitora idzie wyłącznie wtedy, gdy sprawdzenie rzeczywiście
wykryło zmianę pilnowanej strony; rozgłoszenie przy każdym sprawdzeniu
byłoby sygnałem bez treści.

Emulacja urządzenia i przewinięcie strony rozgłaszają zmianę strony z powodem
interakcji, w odróżnieniu od powodu nawigacji — to rozróżnienie niesie pole
`reason` zdarzenia zmiany strony, ponieważ opisuje inne źródło zmiany niż
przejście pod nowy adres. Odczyt drzewa dokumentu, audyt dostępności, odczyt
konsoli i rejestru sieciowego niczego w stronie nie zmieniają, więc zdarzenia
nie rozgłaszają.

Zdarzenie zmiany strony nie niesie wskazania sesji, ponieważ okno przeglądarki
nie jest bytem karty sesji — migawka jest przypisana do okna, nie do sesji.

## budowa/server/internal/core/adapter_modul_developer_api_obciazenie.go

Pojedyncze żądanie developer.api.request oddaje status, czas i rozmiar jednej
odpowiedzi — to wystarcza, by sprawdzić, czy punkt końcowy odpowiada, ale nie
mówi nic o usłudze pod obciążeniem: jeden pomiar nie ma percentyla ani
przepustowości, a to ogon rozkładu, nie średnia, rozstrzyga, czy usługa jest do
użycia. Rzetelny przebieg obciążeniowy wymaga utrzymania zadanej liczby
połączeń równolegle, zbierania histogramu czasów bez wpływu na pomiar
i liczenia percentyli z pełnego rozkładu, nie z próbki — pętla wywołań po
net/http pisana od nowa mierzyłaby w dużej mierze samo siebie.

Wybór padł na autocannon, nie na k6, mimo że oba programy stoją na maszynie
i oba liczą percentyle. K6 opisuje przebieg skryptem w JavaScripcie, a nie
parametrami: wpięcie go tutaj znaczyłoby albo kontrakt niosący program do
wykonania, czyli powierzchnię znacznie szerszą niż adres z parametrami, albo
skrypt składany przez rdzeń, czyli generowanie cudzego języka. Jedyna droga
rdzenia do procesu zewnętrznego zbiera wyłącznie jego wyjście: autocannon
oddaje cały wynik na wyjście we własnym trybie maszynowym, k6 pisze
podsumowanie do pliku, więc wymagałby pisania i odczytu plików pośrednich,
których ta droga nie obsługuje. Kształt, o który pyta kontrakt — percentyle
czasu, żądania na sekundę, bajty na sekundę, rozbicie po kodach stanu —
autocannon oddaje wprost. Siłą k6 są przebiegi narastające, progi
i scenariusze, nieosiągalne bez skryptu, a skryptu nikt tu nie zamawiał.

Program kończy się powodzeniem także wtedy, gdy ani jedno żądanie nie doszło
do skutku: punkt końcowy milczy, a wynik niesie same zera obok licznika
błędów. Podanie takiego wyniku jako pomiaru byłoby brakiem pomiaru
w przebraniu — zero żądań na sekundę czyta się jak usługa skrajnie wolna, a nie
jak usługa, której nie ma — dlatego przebieg bez ani jednej odpowiedzi wraca
odmową nazywającą liczbę błędów. Odpowiedzi spoza klasy 2xx to co innego: punkt
końcowy odpowiedział, pomiar się odbył, więc wynik wychodzi wraz z licznikiem
non2xx i rozbiciem po kodach, żeby wołający zobaczył, że mierzył ścieżkę
błędu, a nie zgadywał.

## adapter_modul_library_technika.go

Rodzaj treści pochodzi z zawartości bajtów, nie z rozszerzenia pliku
(`http.DetectContentType`). Wymiary obrazu pochodzą z nagłówka formatu
(`image.DecodeConfig`), bez dekodowania całego obrazu do pamięci. Liczba stron
dokumentu pochodzi z biblioteki `pdfcpu`, wkompilowanej w binarium. EXIF i GPS
czyta własny czytnik w tym pliku, ponieważ struktura TIFF mieści się
w kilkudziesięciu wierszach kodu i nie uzasadnia zależności zewnętrznej. Czas
trwania nagrania czyta program `ffprobe`, stojący na serwerze razem z rdzeniem.
Pola IPTC, XMP i ID3 czyta program `exiftool`: w odróżnieniu od EXIF nie są
jedną strukturą — IPTC jest zapisem rekordowym w segmencie APP13, XMP drzewem
RDF/XML osadzanym inaczej w każdym formacie kontenera, ID3 dwiema niezgodnymi
rodzinami wersji — więc napisanie własnego czytnika byłoby powtórzeniem pracy
wieloletniej, nieuzasadnionym przy dostępności gotowego programu.

Odczyt metadanych osadzonych nie odmawia z żadnego pojedynczego powodu: brak
narzędzia, plik bez danego rodzaju metadanych albo odpowiedź niemożliwa do
rozebrania zostawiają odpowiednie pole puste, a opis zasobu jest wtedy węższy,
nie błędny.

Współrzędne GPS są jedynym źródłem widoku mapy w module Library — bez nich
widok nie ma czego nanieść, więc czytnik EXIF wydobywa je zawsze, gdy są
obecne w pliku.

## budowa/server/internal/core/adapter_modul_tlumaczenie_korekta.go

Korekta w tym pliku jest mechaniczna: ustalenia powstają z reguł wskazywalnych
w tekście (podwójna spacja, spacja przed znakiem interpunkcyjnym, cudzysłów
prosty zamiast drukarskiego, wielokropek złożony z trzech kropek, zdanie
dłuższe od stu dwudziestu znaków), a każde niesie propozycję poprawki, którą
`proofread.apply` wstawia w treść panelu.

Wywołanie modelu językowego jest w tym miejscu świadomie pominięte. Ustalenie
modelu bywa trafne, ale nie da się go zastosować mechanicznie ani powtórzyć:
ten sam panel dałby przy drugim przebiegu inne ustalenia o innych
identyfikatorach, a `proofread.apply` wskazywałby na ustalenia, których już
nie ma. Pole `channelId` żądania zostaje w kontrakcie nietknięte, ponieważ
komenda nie korzysta z niego.

Reguła napisowa nie rozstrzyga o odmianie słowa ani o jego istnieniu w
słowniku — to zadanie trzech programów zewnętrznych zebranych w
`adapter_modul_tlumaczenie_korekta_silniki.go`. Wchodzą tą samą drogą co
reguły wbudowane: ustalenie z propozycją, którą `proofread.apply` wstawia w
treść. Warunek powtarzalności obowiązuje je tak samo — słownik odpowiada dwa
razy tak samo, model nie.

`consistency.check` porównuje panele okna między sobą oraz pary pamięci
tłumaczeń: to samo zdanie źródłowe przełożone dwoma różnymi zdaniami i ten sam
termin oddany dwoma różnymi słowami w jednym języku. Obie niezgodności są
faktem o danych, nie opinią o stylu.

## adapter_modul_aplikacje_uchwyty.go

Rejestr komend rdzenia potrzebuje jednego miejsca wiazacego nazwe komendy z metoda
portu; ten sam uklad stosuje adapter modulu Library. Zdarzenie apps.build.changed
rozglasza sie z dwoch niezaleznych zrodel niosacych dwa rozne byty: zmiane etapu
budowy produktu oraz zmiane stanu przebiegu wdrozenia. Silnik wykonania wdrozenia
przesuwa przebieg przez kolejne stany po odeslaniu odpowiedzi komendy, dlatego
rozgloszenie tych przejsc nie moze wychodzic wylacznie z obslugiwacza zadania —
adapter przyjmuje droge do emitera przez metode portu, tak jak modul Developer
przyjmuje ja dla przyrostu budowania. Zdarzenie apps.workspace.changed powstaje
inaczej: rozgloszenie idzie z obslugiwacza komendy, ale droga podpieta tym samym
sposobem, poniewaz emiter nalezy do rdzenia, a nie do adaptera; bez tej drogi drugie
okno tej samej przestrzeni nie dowiaduje sie o zmianie pliku warsztatu.

Sesja komunikatu apps.workspace.changed zostaje pusta, poniewaz warsztat nalezy do
okna, a nie do sesji rdzenia: zdarzenie idzie do wszystkich polaczen konta i niesie
identyfikator okna, po ktorym klient je przypisuje. Sesja komunikatu
apps.build.changed dla wdrozenia zostaje pusta z innego powodu: rdzen rozglasza
koniec wdrozenia takze wtedy, gdy okna nie ma juz w rejestrze, bo przebieg przezywa
zamkniecie okna, a klient ma prawo zobaczyc jego wynik. W zdarzeniu tym pole Stage
niesie wylacznie identyfikator okna, ktorego zmiana dotyczy, poniewaz nie ma tu
etapu Product Buildera do pokazania.

## budowa/server/internal/core/adapter_modul_workspace_szukanie.go

Wyszukiwanie w projekcie nie zastępuje przeszukania biblioteki centralnej:
library.file.search przeszukuje bibliotekę i wyłącznie pliki, a ta komenda
przeszukuje jeden projekt i więcej niż pliki — zadania, notatki, wpisy pamięci
i instrukcje. Wyszukiwanie po znaczeniu prowadzi osobna rodzina komend
knowledge; tutaj idzie dopasowanie po słowach.

Treść pliku wchodzi do wyszukiwania przez wyciąg zbudowany komendą
workspace.library.text.extract; plik bez wyciągu jest dopasowywany po samej
nazwie. To jest różnica widoczna dla operatora projektu, więc trafienie
z wyciągu niesie fragment treści, a trafienie po samej nazwie — nie.

## budowa/server/internal/core/adapter_modul_workspace_pamiec_komendy.go

Typ adapterPamieciPrzestrzeni osadza adapter modułu Workspace, więc niesie
komplet jego metod i jest tym samym bytem, którym pracuje okno Context Memory.
Deklaracja repozytoriumWpisowPamieci stoi w tym pliku, a nie przy repozytorium
przestrzeni roboczej, ponieważ jej wymaga wyłącznie rodzina komend memory.*.

Odpięcie wpisu pamięci (memory.detach) zwęża zasięg samego wpisu do poziomu
projektu, który go niesie. Poziom szerszy niż projekt jest jedynym wiązaniem
wpisu poza własnym projektem — to on wprowadza wpis do pamięci innych
projektów na żądanie includeShared. Odpięcia od projektu, który wpisu nie
niesie, nie da się wykonać tą komendą: wstrzymanie ustalenia wspólnego
w cudzym projekcie, module, parze modułów albo karcie sesji robi osobna
komenda memory.disable.set, która nie rusza ani zasięgu wpisu, ani jego
treści. Żądanie odpięcia wskazujące inny projekt kończy się odmową conflict
kierującą do tej komendy. Zdjęcie samego przypięcia (pinned) należy do
memory.set, nie do memory.detach.

Poziomy zasięgu szersze niż projekt są w kodzie wymienione wprost, ponieważ
pierwszeństwo poziomów prowadzi baza danych (kolumna poziom_zasiegu.pierwszenstwo);
porównanie liczbowe w kodzie Go byłoby drugą, rozjeżdżającą się kopią tego
samego porządku.

Wpis wyłączony w memory.list nie wchodzi do wykazu wpisów czynnych, ale też
nie znika bez śladu: idzie osobnym wykazem disabledEntries wraz z zasięgiem,
który go wyłączył, i tożsamością wyłączenia, którym operator znosi je jednym
ruchem. Wyłączenie odsiewa się przed zastosowaniem granicy limitu i osobno od
zawężenia zasięgu, ponieważ cisza bez podania powodu byłaby gorsza niż samo
wyłączenie. Przy aktywnym sicie zasięgu granica limitu stosuje się dopiero na
końcu, a nie przy odczycie z bazy — obcięcie przed sitem oddałoby mniej
wpisów, niż prosi żądanie, i wyglądałoby błędnie na koniec wykazu.

Komenda memory.set niesie pole scopeId, którego workspace.context.set nie ma:
wskazany byt zasięgu jest brany wprost z żądania, a bez wskazania obowiązuje
reguła okna Context Memory — bytem poziomu projektu jest sam projekt, a byt
poziomu szerszego pozostaje pusty.

## adapter_modul_przegladarka_karty.go

`browser.tab.open` wywołane z polem `url` woła to samo pobranie strony, którym
idzie `browser.navigate`, i odkłada migawkę. Inaczej karta byłaby wierszem
w bazie z adresem, którego nikt nie odwiedził, a odpowiedź niosłaby pole
`snapshot` wzięte znikąd.

Przestrzeń robocza zapisuje skład kart, nie ich kopię: `browser.workspace.save`
przypisuje wskazane karty do przestrzeni, a `browser.workspace.open` odtwarza
z nich rząd kart okna. Kopia kart dałaby dwa byty o tym samym adresie i dwie
sprzeczne odpowiedzi na pytanie, która karta jest tą otwartą.

Przywrócenie przestrzeni budzi karty stopniowo: `browser.tab.update` z polem
`active` odsyła kartę po stronę dopiero wtedy, gdy Operator na nią przechodzi,
zamiast odpalać wszystkie przejścia naraz przy samym otwarciu przestrzeni.
Karty zamknięte w chwili zapisu też wracają, ponieważ przestrzeń pamięta
zestaw z chwili zapisu, nie to, co akurat było otwarte — przywrócenie tylko
kart żywych byłoby przywróceniem zestawu okrojonego bez odnotowania ubytku.

## adapter_modul_workspace_pamiec_wylaczenia.go

Trzy czynnosci na wpisie pamieci nie sa zamiennikami jedna drugiej. Komenda
memory.delete usuwa tresc: wpisu po niej nie ma. Komenda memory.detach zweza
zasieg samego wpisu do jego projektu — wpis zostaje, ale obowiazuje wezej, bo
zmienia sie wiersz pamieci. Komenda memory.disable.set wstrzymuje wpis albo
caly poziom pamieci we wskazanym zasiegu: wiersz pamieci zostaje nietkniety,
wpis nie wchodzi do kontekstu i wraca w calosci po zniesieniu wylaczenia.

Wylaczony wpis jest nazwany w odpowiedzi, a nie przemilczany, poniewaz cisza,
po ktorej Operator nie wie, ze wpis jest wylaczony, jest gorsza od braku
wyciszenia. Komenda memory.list oddaje dlatego wykaz wpisow wstrzymanych obok
wykazu czynnych i przy kazdym podaje zasieg, ktory go wylaczyl, razem
z tozsamoscia wylaczenia, ktorym Operator znosi je jednym ruchem. Zniesienie
wylaczenia idzie ta sama komenda memory.disable.set z polem disabled rownym
falszowi: odwracalnosc jednym ruchem jest wymogiem produktu, nie wygoda okna.
Odpowiedz komendy niesie zawsze wykaz po zmianie oraz pole changed, ktore
mowi wprost, czy wykaz naprawde sie ruszyl; odmowa nazywa brak zamiast
milczec — zadanie bez wskazania wpisu ani poziomu nie ma czego wylaczyc,
a zniesienie wylaczenia, ktorego nie ma, konczy sie odmowa not_found.

## budowa/server/internal/core/adapter_modul_aplikacje.go

Zależności adaptera opisane skrótowo przy polach struktury niosą dodatkowe
uzasadnienie poniżej.

`przyrostWdrozenia` rozgłasza `apps.build.changed`. Silnik wykonania wdrożenia
przesuwa stan przebiegu poza wykonaniem komendy, więc rozgłoszenie nie może
iść wyłącznie z obsługiwacza żądania — ten sam wzorzec co `przyrost` w module
Developer.

`okna` jest zależnością opcjonalną: bez niej moduł pracuje, ale montaż ją
wpina, ponieważ wdrożenie bierze z okna przestrzeń roboczą, którą wysyła.

`magazyn` to ten sam magazyn treści, którym jadą zasoby modułu Design i pliki
biblioteki: wiersz w bazie wskazuje plik na dysku, a nie udaje, że go ma.

`katalogDanych` chroni przed odwołaniem, które wypuszczone z rdzenia
wynosiłoby układ katalogów maszyny.

`sejf` dostaje w żądaniu `apps.package.sign` wyłącznie odwołanie
(`signingKeyRef`), nigdy treść klucza — materiał nie przechodzi przez bazę
modułu.

`katalogRozszerzen` zakłada w rejestrze pozycję z manifestu pakietu przy
`apps.package.publish` — prywatny rejestr organizacji nie jest drugim
rejestrem obok `extension.*`, tylko tym samym.

`podglady` żyje wyłącznie w pamięci: serwer podglądu nie przeżywa restartu
rdzenia, więc wiersz w bazie mówiłby po restarcie o nasłuchu, którego nie ma.

`uruchamiacz` jest zależnością opcjonalną: moduł sięga po port startu procesu
w jednym miejscu, audycie wydajności strony w przeglądarce
(`adapter_modul_aplikacje_wydajnosc.go`); bez niej audyt odmawia zdaniem
nazywającym brak, a reszta modułu pracuje dalej.

`PodepnijPrzyrostWarsztatu` i `PodepnijPrzyrostWdrozenia` niosą osobne drogi
rozgłoszenia, bo zdarzenie warsztatu i zdarzenie wdrożenia niosą różne byty —
jedna funkcja o dwóch znaczeniach byłaby dwiema prawdami.

## budowa/server/internal/core/adapter_modul_roundtable_sklad.go

Odczyt stanu debaty istnieje po to, żeby okno otwarte w trakcie debaty
wiedziało więcej niż to, co usłyszało zdarzeniem od swojego otwarcia — bez
niego operator, który odświeżył okno, widziałby debatę pustą aż do następnej
wypowiedzi.

Granica i przesunięcie komendy StanDebaty tną wypowiedzi, nie tury: tura bez
wypowiedzi jest nadal turą i musi zostać pokazana, a wypowiedzi w długiej
debacie idą w tysiącach, więc tylko one wymagają stronicowania.

Stanowisko wchodzi do migawki stanu debaty tylko wtedy, gdy naprawdę leży
w bazie. Składanie stanowiska przy okazji każdego odczytu stanu podbijałoby
licznik wersji przy każdym otwarciu okna, a wersja liczy redakcje, nie odczyty.

Format zapisywanego zespołu bierze się z ostatniej tury debaty, ponieważ to
w niej dany skład właśnie debatował; debata bez ani jednej tury nie ma
formatu do zapisania i zapisany zespół zostaje bez niego.

Przy wniesieniu zapisanego zespołu do okna rola, waga i opis uczestnika idą
osobnym zapisem po samym dopisaniu uczestnika, ponieważ dopisanie opisuje go
wyłącznie parą kanał modelu i nazwa tożsamości, a reszta profilu uczestnika
należy już do osobnej zmiany.

## budowa/server/internal/core/adapter_modul_automations_dziennik.go

Telemetria WebSocket niesie wyłącznie zdarzenia, które padły przy otwartym
oknie; kontrakt żąda pełnego zapisu zdarzeń pojedynczego uruchomienia, także
sprzed otwarcia okna, a tego nie da się oddać z bufora, który znika razem
z procesem — stąd wszystkie czynności tego pliku czytają bazę.

Redakcja sekretów idzie przy odczycie, nie przy zapisie. Reguła redakcji bywa
zmieniana, a zapis raz zredagowany nie da się odredagować — zapis surowy
z redakcją dopiero przy wydaniu zachowuje obie możliwości.

## adapter_modul_terminal_tunele.go

Tunel dało się dotąd założyć wyłącznie poleceniem wydanym w karcie: biegł
wtedy jako zwykły proces polecenia, bez wykazu tuneli i bez ich stanu,
a zamknięcie sprowadzało się do odszukania właściwego wiersza w Process
Monitorze. Prowadzenie tunelu jako bytu długożyjącego zastępuje ten stan.

Wybór programu `ssh` zamiast własnego przekierowania w Go nie jest zasadą
biblioteki zamiast programu, stosowaną gdzie indziej dla git, PDF
i wyszukiwania: tu nie chodzi o czynność biblioteczną, tylko o transport, dla
którego torem jest SSH — program `ssh` niesie uwierzytelnienie, szyfrowanie
i sprawdzenie klucza hosta (`known_hosts`), których własne przekierowanie
musiałoby dorobić od zera. Przekierowanie napisane w Go byłoby drugim,
słabszym torem obok istniejącego, nie usunięciem zależności.

Przełącznik `-o ExitOnForwardFailure=yes` sprawia, że `ssh` kończy się, gdy
przekierowania nie udało się założyć, zamiast biec z otwartym połączeniem
i zamkniętym portem. Rdzeń czeka chwilę na taki koniec przed odpowiedzią,
a potem dogląda procesu do końca jego życia: stan `active` znaczy, że proces
biegnie z założonym przekierowaniem, nie że polecenie zostało tylko wysłane.

## adapter_modul_library_magazyn.go

Magazyn lezy w katalogu danych rdzenia, w podkatalogu biblioteka/tresc, obok
bazy i sejfu poswiadczen, a nie w katalogu roboczym sesji, poniewaz przezywa
restart rdzenia tak samo jak wiersz w bazie, ktory go wskazuje. Obie drogi
zapisu tresci koncza sie blobem w magazynie, bo tresc, ktora ktos z zewnatrz
moze nadpisac, nie jest tresci wersji; dzieki temu czytelnik podgladu nie musi
wiedziec, ktora droga plik przyszedl, bo zawsze czyta sciezke z dysku.

Nazwa pliku jest suma kontrolna tresci. Adapter i tak liczy sume kazdej
przyslanej tresci, a nazwa zbudowana z tej sumy niesie trzy wlasnosci naraz:
dwa wgrania tej samej tresci dziela jeden blob zamiast dwoch kopii, ponowny
zapis tej samej tresci jest bezczynnoscia zamiast nadpisaniem, a nazwa nie
zalezy od nazwy pliku z zadania, wiec nie da sie nia wyjsc z katalogu magazynu.
Pierwsze dwa znaki sumy tworza podkatalog, zeby jeden katalog nie urosl do
dziesiatek tysiecy wpisow.

Zapis jest niepodzielny: plik tymczasowy powstaje na tym samym nosniku co plik
docelowy, wiec przemianowanie nie jest kopiowaniem, i dopiero przemianowanie
czyni tresc widoczna pod odwolaniem. Awaria w polowie zapisu zostawia plik
tymczasowy, nigdy plik obciety, ktory podglad pokazalby jako pelna tresc.
Niepowodzenie zapisu jest odmowa komendy u wolajacego, nie pustym odwolaniem
podanym jako powodzenie.

W komendzie ZapiszZePliku kopia powstaje mimo tego, ze plik zrodlowy gdzies juz
lezy na dysku: odwolanie do cudzego pliku wskazywaloby tresc zywa, a historia
wersji wymaga tresci zamrozonej. Bez kopii plik nadpisany na dysku po wgraniu
zmienialby tresc swojej utrwalonej wersji, a przywrocenie wersji nie mialoby
do czego wrocic. Suma kontrolna liczy sie w locie, w trakcie przepisywania,
bez drugiego przebiegu po pliku i bez wciagania calej tresci do pamieci —
kopiowanie idzie strumieniem, wiec plik o dowolnym rozmiarze przechodzi tak
samo. Nazwa bloba jest suma, znana dopiero po przeczytaniu calosci, wiec plik
tymczasowy powstaje w korzeniu magazynu i dopiero stamtad wedruje pod swoja
sume.

Postacia odwolania wypuszczanego z rdzenia jest sciezka wzgledna magazynu,
liczona od jego korzenia w katalogu danych. Sciezka bezwzgledna wynosilaby do
klienta uklad katalogow maszyny, na ktorej dziala rdzen, a klient stojacy na
innej maszynie i tak pod nia nie siegnie. Klucz nieprzezroczysty odrzucono,
poniewaz wymagalby osobnej komendy rozwiazujacej go do bajtow, a takiej
kontrakt nie przewiduje. Sciezka wzgledna nie niesie zadnego czlonu ukladu
maszyny, jest ta sama na kazdej maszynie i po przeniesieniu katalogu danych,
a rdzen rozwiazuje ja z powrotem do bajtow jednym zlozeniem sciezki z
katalogiem danych. Sciezka spoza magazynu oddaje pustke, poniewaz nie da sie
jej wyrazic wzgledem magazynu; suma kontrolna tresci i tak jedzie w odpowiedzi
jawnie osobnym polem, wiec odwolanie nie wynosi informacji, ktorej odbiorca
by nie mial.

## budowa/server/internal/core/adapter_modul_komponenty.go

Czystej fasady nad samymi magazynami modułowymi zbudować się nie da: pole
config nie ma kolumny w żadnej z tabel modułowych, para poziomu zasięgu
i klucza z component.assign nie ma gdzie usiąść, a component.list nie miałby
wspólnego porządku wyświetlania — stąd własny rejestr komponentów.

Komenda component.create zakłada komponent w magazynie właściwym jego
rodzajowi, a żądanie nie niesie targetId, więc byt docelowy powstaje przy
tworzeniu: projekt przez ZapewnijProjekt, ekspert przez Dodaj, automatyka
przez ZapiszAutomatyke. Komenda component.update zmienia kafel, nie byt
magazynu — nazwa eksperta ma swoją komendę agent.update, definicja automatyki
swoją automation.workflow.save, więc pisanie tych kolumn stąd dałoby dwóch
pisarzy jednej kolumny. Komenda component.delete zdejmuje kafel, a byt
modułowy zostaje pod swoim identyfikatorem; DELETE istnieje wyłącznie dla
tabeli agent, a projekt i automatyka nie mają go w warstwie danych.

Przedrostek kaflowego identyfikatora nie jest przedrostkiem właściwym
magazynowi bytu docelowego: ekspert dostaje ag-, automatyka automat-,
projekt prj-. Stałe modułów są używane wprost, żeby klucz założony przez
component.create wyglądał dokładnie tak jak klucz założony komendą własną
modułu.

Rodzaj assistant odmawia założenia bytu docelowego, ponieważ warstwa danych
oddaje profil asystenta wyłącznie do odczytu i nie ma drogi zapisu, którą
dałoby się tu założyć byt. Kafel bez targetId byłby sukcesem bez skutku,
więc rodzaj odmawia z powodem zamiast zakładać kafel wskazujący na nic.

Przypisanie komponentu (component.assign) nie przenosi bytu modułowego, nie
nadaje uprawnień i nie włącza komponentu do żadnej pętli wykonania — zapisuje
wyłącznie parę poziomu zasięgu i identyfikatora bytu poziomu na wierszu
komponentu. Odpowiedź assigned: false znaczy powtórzenie tego samego
przypisania, przy którym nic nie doszło do skutku.

Sprawdzenie wpięcia rejestru istnieje, ponieważ bez niego component.list
oddałby pusty wykaz, a component.delete — deleted: false; jedno i drugie
wyglądałoby jak stan platformy, a nie jak niezłożony port. Kod internal_error
mówi, że żądanie jest poprawne, a wina leży po stronie montażu adaptera.

## budowa/server/internal/core/adapter_modul_terminal_obserwacje.go

Wyzwalacz plikowy przed tym plikiem istniał w kontrakcie wyłącznie dla
automatyk, czyli poza powłoką terminala — nie dało się powiedzieć „po każdej
zmianie w tym katalogu zbuduj projekt w tej karcie, w jej katalogu, jej
powłoką i jej środowiskiem".

Obserwacja rozpoznaje zmianę przeglądem po czasach modyfikacji, nie zdarzeniami
jądra systemu plików (inotify, ReadDirectoryChangesW), bo te wymagałyby nowej
zależności modułowej i osobnej implementacji na każdy system operacyjny.
Przegląd oparty na bibliotece standardowej zachowuje się jednakowo wszędzie
i jest dokładnie tak dokładny, jak trzeba: obserwacja i tak tłumi powtórzenia
parametrem debounceMs, więc rozdzielczość poniżej progu tłumienia i tak
zostałaby wyrzucona przez samo tłumienie. Cena takiego podejścia — przegląd
katalogu co ustalony odstęp — jest znikoma wobec polecenia, które ten przegląd
wyzwala.

Skutkiem wyzwolenia obserwacji jest proces, nie zapis: wyzwolenie idzie tą samą
drogą co uruchomienie polecenia z poziomu powłoki — brama trybu uprawnień
okna, egzekutor izolacji, port uruchamiacza procesów, wpis w rejestrze
procesów i strumień do konsoli wyjścia. Obserwacja nie jest drugą drogą
uruchamiania procesów, jest wyzwalaczem tej jedynej.

## budowa/server/internal/core/adapter_modul_tlumaczenie_dokument_formaty.go

DOCX, PPTX, XLSX i ODT są archiwami ZIP z dokumentami XML w środku —
`archive/zip` i `encoding/xml` biblioteki Go czytają je bez pomocy z zewnątrz.
Markdown, HTML i tekst czyta się wprost. PDF idzie `pdfcpu`, tą samą
biblioteką, którą warsztat dokumentu modułu Studio: PDF w tym produkcie robią
biblioteki Go, nigdy program obcy.

Wyjątkiem jest rozpoznanie pisma. Dokument bez warstwy tekstowej (skan) nie ma
czego oddać żadnej bibliotece, więc `document.load` z `ocr` woła Tesseracta —
program arsenału serwerowego, który stoi razem z rdzeniem. Jest to droga
wskazana zasadą produktu, a nie obejście: arsenał wolno wołać, programów
spoza arsenału nie wolno.

`zapiszDokumentWyniku` odmawia zapisu formatów, których rdzeń nie umie złożyć
bez utraty układu (PDF, PPTX, XLSX, ODT). Odmowa nazwana jest tu uczciwsza niż
plik z rozszerzeniem, którego treść nie odpowiada rozszerzeniu; DOCX powstaje
jako poprawne archiwum OOXML złożone z dwóch plików XML — tyle wystarcza, żeby
otworzył go edytor tekstu.

## adapter_modul_tlumaczenie_jakosc_zrodlo.go

Porownanie ze zrodlem jest wykonalne, poniewaz panel widzi swoje okno: komenda
quality.check niesie identyfikator panelu, a kod zewnetrzny okna doczytany
zlaczeniem daje to, czego wymaga odczyt okna z repozytorium tlumaczen.

Kontrola swiadomie nie sprawdza dwoch rodzajow wad. Poprawnie przelozona data
zmienia zapis miedzy jezykami, wiec mechaniczne porownanie napisow uznaloby
poprawny przeklad za wade; falszywy alarm w kontroli jakosci jest gorszy od
braku kontroli, poniewaz uczy odbiorce ignorowac wynik, a uczciwe sprawdzenie
wymagaloby rozpoznania daty w obu jezykach, czego rdzen nie ma — zasada
dotyczaca dat idzie za to do polecenia dla modelu, gdzie zadac mozna wiecej,
niz da sie zweryfikowac. Pominiecie zdania jest rzecza znaczeniowa, nie
napisowa: stwierdzenie, czy segment zrodlowy ma swoj odpowiednik w przekladzie,
wymaga rozumienia tresci, a zgrubny zastepnik liczacy zdania mylnie oskarzalby
kazdy przeklad, ktory laczy albo dzieli zdania, co jest normalna, dobra robota
tlumacza; sprawdzenie przyblizone dlugoscia zostaje pod osobnym rodzajem wady,
gdzie jest uczciwe.

Tekst zrodlowy zyjacy w pliku zrodlowym nie jest czytany przez ta warstwe,
poniewaz nie siega ona po pliki spoza repozytorium — porownanie dziala
wylacznie na tresci zapisanej przy oknie. Wykaz walut kontrolowanych jest
zamkniety i krotki, zeby nie brac za kod waluty kazdego skrotowca zlozonego
z trzech wielkich liter. Rozstrzygniecie, czy przecinek w zapisie liczby
oddziela czesc dziesietna, czy tysiace, jest niejednoznaczne z samego napisu;
przyjeta regula uznaje separator z dokladnie trzema cyframi po nim za separator
tysiecy, a kazdy inny za dziesietny — regula myli sie na zapisie w rodzaju
tysiaca pieciuset z przecinkiem dziesietnym, ale trafia w zdecydowanej
wiekszosci zapisow. Kierunek sprawdzenia liczb i walut jest jednostronny:
wartosc dolozona w przekladzie nie jest zglaszana, bo bywa dorobiona uczciwie,
a wartosc zgubiona jest zawsze wada.

## budowa/server/internal/core/adapter_modul_developer_kontekst.go

Wartością komendy developer.contextual.op nie jest samo wywołanie modelu —
to potrafi okno rozmowy. Wartością jest kontekst, którego okno rozmowy nie ma:
treść pliku, na którym Operator stoi, zaznaczenie, na które wskazał, pliki,
które sam dołączył, i stan repozytorium.

Rodzaj renameSymbol jest w kontrakcie razem z operacjami modelu, lecz nie jest
operacją modelu: zmiana nazwy symbolu w całym repozytorium jest czynnością
rozstrzygalną i robi ją serwer języka, który zna graf odwołań. Skierowanie jej
do modelu dałoby wynik prawdopodobny zamiast poprawnego, w czynności, której
poprawność da się sprawdzić w całości.

Granica wielkości pliku wciąganego do kontekstu istnieje, ponieważ kontekst
modelu ma własną granicę, a plik wciągnięty w całości wypchnąłby z niego
zaznaczenie, o które chodziło; plik większy od granicy wchodzi początkiem.

Operacje wyjaśniające oddają samą treść wyniku bez zmiany tekstu, ponieważ
wyjaśnienie nie jest zmianą pliku i wstawianie go do kodu byłoby szkodą; tylko
operacje przepisujące kod oddają dodatkowo zmianę do przyjęcia w edytorze.

Kolejność trzech części polecenia wysyłanego do kanału modelu — czego się
oczekuje, na czym się pracuje, co jest kontekstem — jest stała, ponieważ
model czyta polecenie od początku, a zadanie postawione po tysiącu wierszy
kodu bywa przeczytane jako komentarz do tego kodu.

Plik kontekstu, którego nie da się odczytać, nie zatrzymuje operacji: Operator
dołączył go jako pomoc, a nie jako przedmiot zadania. Cisza byłaby jednak
nieuczciwa, więc rdzeń dopisuje do polecenia informację o pominięciu.

Zakres zmiany zależy od tego, na czym Operator pracował: przy zaznaczeniu jest
nim samo zaznaczenie, przy całym pliku — cały plik od pierwszego wiersza.
Treść wyniku wraca w polu result niezależnie od zakresu zmiany.

Model odpowiada kodem w ogrodzeniu znaczników nawet wtedy, gdy poproszono
o samą treść, a ogrodzenie wstawione wprost do pliku źródłowego jest błędem
składni w każdym języku, więc zdejmijOgrodzenieKodu usuwa je przed zapisem.

Bez sprawdzenia osobneSlowo dopasowanie symbolu Plik trafiałoby też
w PlikRoboczy, a zmiana nazwy ruszyłaby symbol, którego nikt nie wskazał.

## budowa/server/internal/core/adapter_modul_badania_przestrzen.go

Pokrycie pytania badawczego liczy się z wiązań w katalogu źródeł, nie
z deklaracji: pytanie jest pokryte wtedy, gdy istnieje źródło przypisane do
niego wprost. Licznik pokrycia wyliczony z samej liczby źródeł mówiłby, że
badanie posuwa się do przodu, nawet gdy operator dorzuca materiał niezwiązany
z żadnym pytaniem — a to jest dokładnie ten stan, przed którym panel postępu
ma ostrzegać.

Świeżość badania jest faktem, a sugestia odświeżenia progiem: data najnowszego
źródła bierze się wprost z bazy, a sugestia odświeżenia jest porównaniem tej
daty z progiem świeżości. Próg wchodzi do odpowiedzi razem z sugestią, żeby
operator wiedział, wobec czego rdzeń mierzy, a nie dostawał samego ostrzeżenia
bez podstawy.

## adapter_modul_asystent_wykonawca.go

Typ `adapterAsystenta` i przyjęcie polecenia deklaruje moduł asystenta
podstawowy; stan i dziennik prowadzi moduł czynności asystenta. Ten plik
dokłada metody na tym samym typie oraz zależności wykonawcy wpinane
w montażu: rejestr kanałów, nadzorca sesji, nadajnik zdarzeń i kontekst życia
rdzenia. Rdzeń nie ma drugiego silnika modelu — parametry tury bierze
z okna asystenta, najwęższego poziomu zasięgu.

Kontrakt komendy głosowej niesie samą transkrypcję, wolny tekst, bez pola
mówiącego, że to komenda platformy. Zgadywanie po treści, czy padło polecenie
dla modelu czy dla rdzenia, byłoby wykonaniem czynności, której nikt
jednoznacznie nie zlecił — dlatego wykonawca prowadzi wyłącznie drogę modelu
i nie udaje akcji platformy rozpoznaniem wolnego tekstu.

Zlecenie w stanie innym niż stan wejściowy — podjęte inną drogą, wstrzymane,
anulowane — wykonawca zostawia bez ruchu: to jedyne ciche wyjście, bo nikt
niczego nie zlecił drugi raz. Każdy pozostały brak kończy zlecenie głośno
stanem `failed` z nazwanym powodem, zamiast zostawiać wiersz stojący
w `queued` bez wyjaśnienia — Operator ma wtedy z czego zdecydować: dosłać
tekst, wskazać kanał, ponowić zlecenie.

## adapter_modul_tlumaczenie_silniki.go

Komenda engine.compare wola model naprawde, kanal po kanale, i oddaje warianty
obok siebie. Bez wpietego rejestru kanalow albo przy kanale nieczynnym odmawia
wprost, tak samo jak zapis przekladu w panelu — wariant pusty udawalby, ze
model odpowiedzial.

Komenda batch.run wykonuje prace, a nie zapowiada ja. Zlecenie zaklada wiersze
pozycji zlozonych z panelu i operacji, i przechodzi je po kolei: operacja
przekladu woła model, kontrola jakosci zapisuje niezgodnosci, korekta zaklada
ustalenia, a wydanie zapisuje slad eksportu panelu. Pozycja, ktora sie nie
powiodla, zostaje w stanie bledu wraz ze szczegolem, a nie znika z kolejki —
niepowodzenie jednej pozycji nie przewraca calego przebiegu, poniewaz pakiet
ma dojsc do konca i pokazac, co sie udalo, a co nie.

## budowa/server/internal/core/adapter_modul_automations_wersje.go

Porównanie wersji jest strukturalne, nie tekstowe. Kontrakt oddaje
identyfikator kroku wraz z rodzajem zmiany i nazwą pola, więc różnica liczy
się po krokach — krok dodany, krok usunięty, krok o zmienionym polu. Różnica
tekstowa dwóch zapisów strukturalnych mówiłaby o wierszach zapisu, a nie
o krokach procesu.

Symulacja nie wywołuje niczego: przebieg próbny czyta definicję, sprawdza
spójność każdego kroku i oddaje wynik kroku po kroku, nie ruszając ani
kolejki, ani kanału modelu, ani żadnego interfejsu zewnętrznego, zgodnie
z zapewnieniem kontraktu, że efekty uboczne kroków są przy symulacji
wstrzymane. Rdzeń, który przy przebiegu próbnym wysłałby raport pocztą, byłby
rdzeniem, któremu nie wolno ufać.

## budowa/server/internal/core/adapter_modul_tlumaczenie_mowa_silnik.go

Metody tego pliku stoją na wspólnym `*adapterTlumaczenia` (typ i przedrostki
deklaruje `adapter_modul_tlumaczenie.go`); `SyntezujMowe` mieszka
w `adapter_modul_tlumaczenie_mowa.go`, a wybór głosu —
w `adapter_modul_tlumaczenie_mowa_glos.go`.

Dźwięk nie opuszcza maszyny Operatora: nie ma tu klienta HTTP ani adresu,
a oba syntezatory są programami lokalnymi. Silniki są dwa, w tej kolejności
pierwszeństwa: `piper` — synteza neuronowa, wchodzi pierwszy, gdy stoi
binarium oraz jest głos dla języka panelu; głos leży na maszynie jak każde
inne binarium arsenału, więc jest zależnością środowiska, a nie stanem
produktu. `espeak-ng` — synteza formantowa, program jednym plikiem bez stanu
i bez modeli, droga zapasowa, gdy pipera nie ma albo nie ma dla tego języka
głosu.

Operator ma wiedzieć, którym silnikiem słucha: głos zapasowy brzmi inaczej
niż dobry i nie ma być mylony z usterką nagrania. Kontrakt
(`TranslateSpeechSynthesizeResponse`) niesie same `panelId` i `path`, bez pola
na nazwę silnika, więc nazwa silnika idzie w nazwę pliku:
`pan-…-piper-….wav` albo `pan-…-espeak-ng-….wav`. Ta sama nazwa ląduje
w kolumnie `nagranie_odnosnik`, więc ślad również mówi, kto czytał.

Ścieżka programu i głosu bierze się z dwóch źródeł, w tej kolejności: zmienna
środowiska (`DANACO_PIPER`, `DANACO_PIPER_GLOSY`, `DANACO_ESPEAK`) ma
pierwszeństwo, bo arsenał jest instalowany poza produktem i bywa na każdej
maszynie w innym miejscu; następnie wykrycie w miejscach typowych — nazwa
goła w PATH, a przy jej braku katalog arsenału. Katalog ustawień produktu
źródłem nie jest: wpisywałby położenie cudzego binarium do stanu produktu,
a to jest fakt maszyny, nie nastawa Operatora.

Brak głosu to inna odmowa niż brak binarium: program stoi, więc `Stoi` mówi
„jest", a czynność i tak nie wyjdzie. Rozróżnienie widać w treści odmowy, bo
naprawy są różne. Gdy zawiodą oba silniki, odmowa wymienia obie przyczyny
osobno.

## budowa/server/internal/core/adapter_modul_aod_rozmowa.go

Nakładka niczego nie zakłada sama. Wiadomość idzie tym samym portem rozmowy,
którym jedzie message.send, więc wchodzi do tej samej historii okna, uruchamia
tę samą turę modelu i tę samą telemetrię. Polecenie głosowe idzie tym samym
modułem Assistant, którym jedzie assistant.voice.command, więc zlecenie
z nakładki widać w Actions Monitor. Druga droga do wiadomości albo do
zlecenia byłaby drugą prawdą o tym samym bycie.

Podpowiedzi pochodzą z katalogu akcji (tabela akcja) — jedynego zbioru w tym
rdzeniu, który wiąże byt zasięgu z komendą kontraktu do wywołania, a taki
właśnie kształt ma AodSuggestion z polem commandType. Podpowiedzi nie są
układane w kodzie: dopisanie wiersza katalogu daje nową podpowiedź bez zmiany
rdzenia. Komenda aod.suggestion jest odczytem, bo ma żądanie i wynik; zdarzeń
aod.suggestion.* kontrakt nie zna, więc rdzeń podpowiedzi nie wypycha.

Komenda WyslijZNakladki oddaje obok odpowiedzi kontraktu całą przyjętą
wiadomość, ponieważ rozgłoszenie zdarzenia zmiany wiadomości należy do
uchwytu i bez treści wiadomości nie miałoby czego rozgłosić. Okno oddane
w odpowiedzi jest tym, które wiadomość naprawdę przyjął port rozmowy z
założonego wiersza; wskazanie z żądania jest tu tylko punktem wyjścia.

Komenda KontekstNakladki oddaje ten sam komplet kontekstu okna, który przenosi
context.transfer i który zasila Context Panel. Magazyn kompletu jest jeden na
cały rdzeń; nakładka czyta z niego, a nie z drugiego, własnego magazynu.
Metoda KompletOkna mieszka przy adapterze przenoszenia, a nie przy nakładce,
ponieważ przenoszenie kontekstu składa komplet po drodze do okna docelowego
i nie miało dotąd czytelnika pytającego o komplet okna wprost.

W metodzie ogonOdwolan historia rośnie w przód, więc obcięcie do granicy idzie
od początku wykazu — nakładce potrzebne są wiadomości najświeższe, nie
najstarsze.

Żądanie oknoPodpowiedzi wskazujące byt nieistniejący jest odmawiane; żądanie
niewskazujące niczego odmowy nie dostaje — zostają wtedy podpowiedzi całej
platformy, bo są prawdziwe niezależnie od okna. Proces nieznany rejestrowi
telemetrii jest bytem nieistniejącym, więc odpowiedź złożona z podpowiedzi
globalnych byłaby wtedy ciszą udającą wynik i komenda odmawia zamiast jej
oddać.

Kolejność zasięgów w pozycjeZasiegow jest kolejnością ważności podpowiedzi:
to, co dotyczy okna, stoi przed tym, co dotyczy platformy — od bytu
najwęższego do najszerszego, bez powtórzeń.

Kolumna akcja.warunek_dostepnosci nazywa byt, którego wymaga akcja (okno,
sesja, kanał, środowisko, kolejka; puste znaczy: żadnego). Nakładka
rozstrzyga dwa z nich — okno i jego kartę sesji — bo tyle wynika z żądania
aod.suggestion. Podpowiedź wymagająca bytu, którego nakładka nie zna, byłaby
przyciskiem bez celu, więc do wykazu nie wchodzi.

Opisu i ikony akcji AodSuggestion nie niesie, więc zostają w katalogu akcji.
Pole moduleId idzie wprost z okna, a nie okrężnie przez window.list po
stronie nakładki, ponieważ bez tego pola wyciszenie bieżącego modułu nie
miałoby po czym rozpoznać swojej sugestii, a cisza bez podstawy jest gorsza
niż ujawnienie. Klasy zdarzenia podpowiedź z katalogu akcji nie niesie i nie
ma nieść: pozycja katalogu jest czynnością osiągalną, a nie skutkiem
zdarzenia wyzwalającego — klasę niosą sygnały aod.signal.*, tam gdzie
zdarzenie naprawdę zaszło; zgadnięta klasa wpuszczałaby podpowiedź
w wyciszenie, którego Operator na nią nie założył.

Komenda PolecenieGlosoweNakladki nie rozpoznaje mowy sama: nagranie
przepisuje moduł Assistant przez wpięty port Mowa. Drugi silnik po tej
stronie byłby drugą prawdą o tym samym bycie. Nagranie przetworzone bez mowy
odmawia po stronie Assistanta, a nakładka oddaje tę odmowę bez zmiany, więc
wymagane pole transcript nie wraca puste. Pole speak nie jest spełniane:
syntezy mowy rdzeń nie ma, więc speechRef zostaje pusty — pole jest
opcjonalne, a brak jest odpowiedzią zgodną z kontraktem, tak samo jak
w assistant.voice.command i speech.synthesize.

## adapter_modul_orchestration.go

Rodzina orchestration.* jedzie na tej samej maszynerii co okno Orchestrator,
zapisanej w tabelach kroku automatyki i zaleznosci kroku automatyki, a nie na
wlasnej. Roznica wobec okna jest ziarno: okno Orchestrator przyjmuje komplet
lukow, bo Workflow Builder wysyla caly uklad po zmianie, a rodzina
orchestration.* pracuje pojedynczym lukiem — dependency.set doklada jeden,
dependency.remove zdejmuje jeden. Zapis idzie dlatego osobnymi metodami zapisu
i usuniecia pojedynczej zaleznosci, a nie podmiana kompletu, bo inaczej
dolozenie jednego luku przepisywaloby wszystkie pozostale. Metody pliku stoja
na adapterze modulu Automations, wiec uklad zaleznosci ma w rdzeniu jednego
wlasciciela; sprawdzenie ukladu, wykrycie cyklu i sciezka krytyczna pochodza
z osobnego pliku tego samego adaptera.

Ocena ukladu nie blokuje zapisu. Kontrakt komendy dependency.set mowi to
wprost: luk domykajacy cykl albo prowadzacy do kroku, ktorego jeszcze nie ma,
zapisuje sie, a zastrzezenia wracaja w odpowiedzi osobnym polem oznaczajacym
uklad niepoprawny. Odmawiane sa wylacznie zadania, ktorych schemat nie zna:
luk bez wskazania kroku, petla wlasna oraz rodzaj zaleznosci spoza trzech
wartosci kontraktu — sa to bledy zadania, nie stan ukladu.

## adapter_modul_roundtable_zgoda.go

Wszystkie sześć komend liczy rdzeń bez wywołania modelu, bo panel odczytuje
wskaźniki przy każdym otwarciu okna i po każdej turze — muszą być
natychmiastowe, powtarzalne i darmowe. Miara jest słabsza od zanurzeń
semantycznych i nie rozpoznaje synonimów, ale liczba, którą oddaje, jest ta
sama przy każdym odczycie i nie zależy od tego, który model akurat
odpowiedział. Stanowiskiem uczestnika w turze jest złożenie wszystkiego, co
w niej powiedział; uczestnik, który w turze milczał, nie ma w niej stanowiska
i nie wchodzi do żadnego z tych rachunków — cisza nie jest zgodą ani sporem.

Między progiem zgody a progiem sporu leży pas, w którym para stanowisk nie
jest ani zgodna, ani sporna: wymuszenie rozstrzygnięcia w tym pasie dawałoby
punkt zgody tam, gdzie uczestnicy powiedzieli po prostu co innego o czym
innym.

Trafność w kalibracji liczy się z tego, co naprawdę rozstrzygnięto: wypowiedź
jest trafiona, gdy Operator wskazał ją jako bardziej przekonującą, dał jej
wysoką ocenę, albo gdy wygrała głosowanie. Uczestnik bez ani jednej ocenionej
wypowiedzi nie ma trafności do zmierzenia i nie wchodzi do wyniku — liczba
wzięta z zera pomiarów byłaby wymysłem.

## budowa/server/internal/core/adapter_modul_przegladarka_kanaly.go

`browser.feed.subscribe` pobiera wskazany adres, rozpoznaje postać (RSS, Atom,
JSON Feed) i odkłada wpisy przy subskrypcji — wiersz z samym adresem, bez ani
jednego wpisu, byłby subskrypcją, o której nie wiadomo nawet, czy pod tym
adresem stoi kanał.

Rozbiór idzie biblioteką standardową (`encoding/xml`, `encoding/json`), bo RSS
i Atom są dokumentami XML o ustalonym kształcie, a JSON Feed dokumentem JSON.
Zewnętrzna biblioteka kanałów nie dołożyłaby tu niczego poza kolejną
zależnością.

## adapter_modul_roundtable_wydanie.go

Cztery formaty transkryptu powstaja trzema drogami. Markdown i JSON skada
rdzen wprost, bez niczego z zewnatrz. PDF powstaje biblioteka wkompilowana
w binarium — dokument i kryptografia sa w tym produkcie wyjatkiem
bezwzglednym od wolania programow serwerowych. DOCX powstaje Pandociem,
poniewaz formatu biurowego nie da sie zlozyc bibliotecznie w rdzeniu,
a Pandoc jest programem serwerowym zadeklarowanym w sondzie zaleznosci
zewnetrznych i uzywanym juz przez modul Studio, modul Translate i modul
Biblioteki.

Pandoc zamienia transkrypt, a nie przetwarza materialu filmowego;
przekroczenie granicy czasu na zamiane formatu znaczy plik uszkodzony albo
proces, ktory utknal. Zamiana na DOCX przekazuje material plikiem, a nie
strumieniem, poniewaz Pandoc rozpoznaje format wyjsciowy po rozszerzeniu
pliku docelowego, a zapis do strumienia wymagalby wskazania formatu osobno
i tak samo tworzylby plik posredni. Lamanie wierszy dokumentu PDF idzie po
slowach, poniewaz lamanie w srodku slowa dawaloby zapis, ktorego nie da sie
przeczytac ani przeszukac. Zadania modulu Roundtable niosa okno debaty,
a nie okno sesji terminalowej, wiec zasady izolacji dla wywolan arsenalu
biora sie z zasiegu platformy, tak samo jak w rodzinie narzedzi mediow —
punkt izolacji wlaczony globalnie dziala tu tak samo jak dla modulu
Terminal.

## budowa/server/internal/core/adapter_modul_terminal_klucze.go

Para kluczy powstaje biblioteką standardową Go — crypto/ed25519, crypto/rsa
i crypto/ecdsa wytwarzają materiał, a golang.org/x/crypto/ssh zapisuje go
w postaci OpenSSH i liczy odcisk. Program ssh-keygen byłby tu narzędziem spoza
instalki wołanym po to, żeby zrobić rzecz, którą biblioteka standardowa robi
w kilku wierszach — a wytworzenie klucza jest czynnością, bez której książka
hostów przestaje mieć czym się łączyć. Ta sama decyzja zapadła przy obsłudze
git, PDF i wyszukiwania w innych modułach rdzenia.

Klucz prywatny nie opuszcza dysku maszyny rdzenia w żadną stronę: wytworzenie
oddaje sam odcisk i klucz publiczny, wciągnięcie do wykazu bierze ścieżkę, nie
treść, a wykaz nie ma pola, w którym materiał tajny mógłby się znaleźć.

Kontrakt każe podawać hasło klucza odwołaniem do sejfu, nigdy treścią. Rdzeń
nie ma dziś czytnika sejfu, więc żądanie z odwołaniem kończy się odmową
nazywającą ten brak, zamiast po cichu wytworzyć klucz bez hasła w odpowiedzi
na prośbę o klucz z hasłem — takie zachowanie byłoby cichym obniżeniem
ochrony, gorszym niż odmowa.

Klucz chroniony hasłem przy wciąganiu do wykazu czyta się tylko z hasłem,
którego rdzeń nie ma, i to nie jest powód odmowy: wykaz ma nieść taki klucz,
a program ssh odczyta go sam przy połączeniu. Odcisk bierze się wtedy z klucza
publicznego — z pliku obok albo z części publicznej, którą niesie sam błąd
odczytu.

## budowa/server/internal/core/adapter_modul_aplikacje_wydajnosc.go

Moduł umiał dotąd powiedzieć o wdrożonym produkcie jedno: czy odpowiada
(apps.deployment.health.get — dostępność, czas nieprzerwanego działania,
wynik ostatniego sprawdzenia kondycji). To jest odpowiedź na pytanie, czy
produkt stoi, nie na pytanie, jak szybko się otwiera; produkt odpowiadający
w cztery sekundy jest dostępny w stu procentach i nie do użycia.

Core Web Vitals nie są czasem odpowiedzi serwera: największe wymalowanie
treści i przesunięcia układu powstają w przeglądarce, po wykonaniu skryptów,
a ich wartość zależy od emulacji urządzenia i dławienia sieci. Rdzeń, który
mierzyłby to własnym klientem HTTP, oddałby czas pobrania dokumentu i nazwał
go wydajnością strony — liczbę prawdziwą, odpowiadającą na inne pytanie.

Program pomiarowy mówi o nieodbytym pomiarze wprost: przebieg, w którym
strona się nie wczytała, niesie w odpowiedzi pole runtimeError wraz z kodem
powodu i nie niesie ocen. Rdzeń czyta to pole przed czymkolwiek innym — bez
tego odczytu odpowiedź o produkcie, którego pod adresem nie ma, składałaby
się z samych zer i wyglądałaby jak strona wolna, a nie jak strona
niezmierzona.

Audyt trwa kilkanaście sekund przy stronie zdrowej i nie kończy się nigdy
przy stronie, która nie przestaje się wczytywać. Granica idzie do programu
przez opcję --max-wait-for-load i osobno do arsenału, z zapasem — pierwszy
mija program, więc przekroczenie nazywa ten, kto wie, na co czekał.

Deklaracja narzedzieLighthouse stoi przy miejscu użycia; wykaz zależności
odwołuje się do niej, zamiast powtarzać nazwę programu po raz drugi.

Kolejność miar w miaryWydajnosciStrony jest ustalona, żeby dwa kolejne audyty
tej samej strony dawały wykaz w tym samym porządku — wynik ma się różnić
wtedy, gdy zmieniła się strona, a nie wtedy, gdy inaczej ułożyła się mapa
odpowiedzi programu. Wykaz obejmuje Core Web Vitals wraz z miarami, z których
te się liczą; miara dopisana tu bez pokrycia w odpowiedzi programu wyszłaby
z audytu jako zero, dlatego brak miary w odpowiedzi pomija się, zamiast
wypełniać wartością zastępczą.

Postać biurkowa ma w programie własną nastawę zbiorczą, ponieważ sama opcja
--form-factor zmienia sposób liczenia oceny, lecz zostawia emulację
i dławienie telefonu, więc bez tej nastawy wynik byłby oceną biurka
policzoną na warunkach telefonu.

Program kończy się kodem niezerowym także wtedy, gdy pomiar się nie odbył,
a powód opisał w odpowiedzi. Odpowiedź czytamy więc przed rozpatrzeniem
odmowy arsenału, inaczej sytuacja, w której strony nie ma pod danym adresem,
wyszłaby jako zwykłe zakończenie programu niepowodzeniem.

Miara, której program nie policzył, nie wchodzi z wartością zero, ponieważ
zero jest w tych miarach wynikiem najlepszym z możliwych.

Zaokrąglenie oceny programu w dół dałoby 99 dla strony ocenionej idealnie,
więc wSkaliStu zaokrągla do najbliższej liczby całkowitej.

Brak programu jest brakiem, który operator serwera usuwa jedną instalacją,
a przekroczenie granicy czasu jest przekroczeniem, nie awarią. Obie sytuacje
bez tego rozróżnienia wychodziłyby jako internal_error, mówiące czytającemu
coś nieprawdziwego o tym, co się stało.

## adapter_modul_roundtable_konsensus.go

Do chwili redakcji rdzen skladal tresc stanowiska z zapisu tur przy kazdym
odczycie. Od chwili, w ktorej Operator nada stanowisku wlasna tresc, zlozenie
z tur go nie dotyka — pilnuje tego osobna kolumna oznaczajaca redakcje oraz
warunek w zapytaniu zapisujacym; bez tego pierwsze otwarcie panelu po redakcji
kasowaloby prace Operatora. Kazda redakcja odklada osobny wiersz wersji,
a nie nadpisuje jeden licznik: licznik w kolumnie wersji mowi jedynie, ile
redakcji bylo, a porownac dwie redakcje da sie dopiero wtedy, gdy kazda
zostala zapisana osobnym wpisem.

Poparcie wazone stanowiska liczy sie z wag wszystkich uczestnikow debaty
w mianowniku, a w liczniku wylacznie z wag tych, ktorzy nie podpisali zdania
odrebnego. Uczestnik, ktory zglosil zdanie odrebne, nie poparl stanowiska,
lecz jego waga i tak wchodzi do mianownika rachunku, bo poparcie mierzy udzial
w calej debacie, nie tylko wsrod uczestnikow bez zastrzezen.

## adapter_modul_extension_zaufanie.go

Nic w tym pliku nie stanowi bramy instalacji: żadne ostrzeżenie skanera ani
brak podpisu nie blokuje instalacji ani włączenia rozszerzenia. Kontrola idzie
przez stan wyjściowy i zakres uprawnień, nie przez odmowę na wejściu. Skaner
manifestu wystawia spostrzeżenia, weryfikacja podpisu wystawia werdykt,
nadanie uprawnień zapisuje decyzję Operatora — żadna z tych trzech dróg
niczego nie wstrzymuje.

Uprawnienie nadmiarowe liczy się, nie jest zgadywane: `extension.permission.
list` oddaje zakresy nadane, których manifest wcale nie deklaruje. To różnica
dwóch zbiorów leżących w bazie, a nie heurystyka.

Weryfikacja podpisu liczy podpis od nowa: Ed25519 ze standardowej biblioteki
sprawdza bajty podpisu kluczem publicznym odłożonym przy pozycji. Przepisanie
zapamiętanego wyniku poprzedniej weryfikacji nie byłoby weryfikacją, tylko
powtórzeniem cudzego zdania.

## budowa/server/internal/core/adapter_modul_extension_protokol.go

Klient protokołu Model Context Protocol jest wkompilowany, nie pożyczony:
protokół to JSON-RPC 2.0 nad jednym z trzech transportów — procesem lokalnym
przez stdio, strumieniem zdarzeń SSE albo zwykłym HTTP. Wszystkie trzy
obsługuje biblioteka standardowa Go, więc rdzeń rozmawia z serwerem sam, bez
ani jednego programu obok instalki.

Program serwera należy do operatora, nie do platformy: transport stdio
uruchamia polecenie, które operator sam wpisał przy podłączaniu rozszerzenia.
To nie jest zależność rdzenia od cudzego programu — to jest cudzy program,
który operator świadomie podłączył, i którego brak wraca nazwaną odmową,
a nie awarią platformy.

Każda ramka rozmowy idzie do dziennika: żądanie i odpowiedź zapisują się wraz
z korelacją, więc wykaz dziennika protokołu pokazuje rozmowę, która naprawdę
się odbyła, a diagnoza błędu integracji ma z czego wyjść.

## budowa/server/internal/core/adapter_modul_asystent_czynnosci.go

Typ `adapterAsystenta` i przyjęcie polecenia głosowego deklaruje
`adapter_modul_asystent.go`; ten plik dokłada metody na tym samym typie, bez
drugiej deklaracji. Zlecenie asystenta ma automat stanu
`queued/running/paused/done/failed/cancelled` z kontraktu — inny byt niż
katalog akcji (`core/akcje.go`, statyczny spis dostępnych czynności) i inny
niż stan pętli sesyjnej kolejki. Nieznane zlecenie wraca odmową, nie cichą
zgodą: sterowanie zleceniem, którego nie ma, nie może wyglądać jak sukces.

Stan sprzed zapisu w `StanCzynnosci` czytany jest zawsze, bo od niego zależą
dwie rzeczy: `resume` ma wznowić zlecenie wstrzymane, a nie dokładać drugiej
tury zleceniu, które właśnie biegnie (oba stoją po zapisie w stanie
`running`, więc po nim już ich nie odróżnić); a zlecenie zamknięte
(`done`/`failed`/`cancelled`) nie ma czego wznawiać ani wstrzymywać. Bez tego
odczytu `resume` na zleceniu zamkniętym przestawiłby je na `running`, nikt by
go nie podjął — wykonawca wchodzi tylko ze stanu `queued`/`paused` — i wiersz
stałby w toku bez końca, bez wykonawcy i bez zdarzenia.

Anulowanie i wstrzymanie sięgają do biegu, nie tylko do wiersza: zlecenie
biegnące ma turę modelu w locie, a samo przestawienie stanu zostawiłoby ją
pracującą dalej na cudzy rachunek. Przerwanie idzie przed zapisem, żeby
wykonawca zastał już decyzję Operatora, gdy tura wróci z błędem.

Priorytet idzie osobnym zapisem, bo kontrakt pozwala przy sterowaniu zmienić
kolejność obsługi w Actions Monitor — przyjęcie priorytetu bez zapisania go
byłoby potwierdzeniem czynności, która się nie odbyła.

`retry` przestawia zlecenie na `queued`, które znaczy „do wykonania", nie
„czeka na kolejny ręczny ruch": wykonawca podejmuje je tą samą drogą, którą
podejmuje zlecenie świeżo złożone. `resume` oddaje zlecenie temu samemu
wykonawcy, tylko ze stanem wejścia `running`, bo wykonawca podejmuje
wyłącznie `queued`. Bez wpiętego wykonawcy żadna z dwóch dróg nie robi nic.

## adapter_modul_przegladarka_dostepnosc.go

Moduł czyta stronę żywą trzema sondami: drzewem elementów, rejestrem żądań
i konsolą. Audyt dostępności jest czwartą sondą i pyta o to samo — o stronę po
zbudowaniu przez skrypty, nie o jej źródło. Reguła dostępności orzeka
o etykiecie kontrolki wstawionej skryptem tak samo jak o etykiecie wpisanej
w źródle; audyt czytający sam HTML odpowiedziałby o dokumencie, którego nikt
nie ogląda. Adres bierze się z ostatniej migawki okna, tak samo jak w trzech
sondach starszych; kontrakt nie niesie w żądaniu pola adresu, bo pyta
o bieżącą stronę okna.

Reguły WCAG są cudzą wiedzą i rdzeń jej nie przepisuje: między normą a jej
sprawdzeniem stoją setki reguł, które ktoś utrzymuje wraz z kolejnymi
wydaniami normy. Dlatego audyt idzie zewnętrznym programem, a nie własnym
obchodem drzewa. Program dostaje tę samą przeglądarkę, którą rdzeń już
deklaruje dla sond starszych, nie własną kopię pobieraną z sieci przy
pierwszym uruchomieniu.

Program audytujący nie mówi, czy strona się wczytała: dokument błędu 404 bywa
poprawny wobec normy i wychodzi z audytu jako pusty wykaz naruszeń. Odpowiedź
zero naruszeń dla strony, której pod danym adresem nie ma, byłaby brakiem
pomiaru podanym jako pomiar — i to najgorszą jego postacią, bo wygląda dobrze
i nie wzywa nikogo do sprawdzenia. Dlatego przed audytem rdzeń sięga po stronę
własną drogą modułu, która orzeka o stanie odpowiedzi i o tym, czy zasób jest
w ogóle stroną; odmowa stąd nazywa, czego nie zmierzono, zamiast podawać zero.
Drugie potwierdzenie jest po stronie odpowiedzi programu: pusty wykaz jest
wynikiem tylko wtedy, gdy program oddał tablicę JSON.

Program audytu dostaje przeglądarkę wskazaną, nie szukaną: jego własna
warstwa sterowania przeglądarką pobiera wydanie Chrome do katalogu pamięci
podręcznej użytkownika, a rdzeń takiego pobrania nie robi. Wskazanie idzie
plikiem nastaw, bo jedyna droga do procesu nie przekazuje zmiennych
środowiska — i ma nie przekazywać, bo binarium arsenału nie ma powodu widzieć
zmiennych rdzenia.

## budowa/server/internal/core/adapter_modul_mowa.go

Silnik mowy wymaga trójki okno, zasady i obszar, ponieważ pomocnik
transkrypcji startuje tym samym uruchamiaczem i przez tę samą bramę izolacji,
co każdy inny proces drzewa. Terminal i Developer biorą tę trójkę z okna
żądania: rejestr okien daje okno, rozstrzygacz nad zasięgiem daje zasady dla
kontekstu okna, a ustalacz katalogu roboczego daje obszar. Żądania rodziny
speech okna nie niosą: rodzina jest zdolnością platformy, nie okna, więc trójka
składa się dla pustego kontekstu zasięgu, tą samą drogą co dla okna. Pusty
kontekst zasięgu jest poprawnym adresem najszerszego poziomu: rozstrzygacz
oddaje wtedy politykę platformy, a ustalacz katalog roboczy platformy. Wpisane
z ręki zasady i obszar puste znaczyłyby izolację wyłączoną niezależnie od
ustawień. Okno jest jedyną wartością, którą adapter wypełnia sam, ze
środowiskiem wykonania rdzenia, bo ścieżka nagrania wskazuje maszynę silnika
i pomocnik musi ruszyć na hoście rdzenia. Identyfikatora okna nie ma skąd
wziąć, więc wpis dziennika transkrypcji nie dostaje odnośnika okna.

Silnik mowy powstaje na każde wywołanie, ponieważ ustawienie nastaw mutuje
byt: jedna instancja współdzielona przez równoległe żądania oznaczałaby wyścig
o nastawy. Przy okazji nastawy są świeże — zmiana modelu mowy komendą
konfiguracji obowiązuje od następnej transkrypcji, bez restartu rdzenia.

Trzy pola adaptera obsługujące komendy dobudowane obok transkrypcji —
przyjęcie i oddanie bajtów nagrania oraz nastawę wybudzania — są zależnościami
opcjonalnymi: bez nich te komendy odmawiają, nazywając brak, a rozpoznawanie
mowy pracuje bez zmian.

Odsłuch dla availability liczony jest osobno od dyktowania, bo jedzie innym
łańcuchem: syntezator głosu, nie rozpoznawanie mowy. Nawet gdy dyktowanie
odmawia, odsłuch bywa gotowy, dlatego wynik idzie do obu gałęzi odpowiedzi,
a nie tylko do udanej. Trzy dobudowane zdolności rodziny meldują się osobno,
bo osobno znikają: przyjęcie nagrania zależy wyłącznie od magazynu rdzenia
i działa nawet bez silnika mowy, bo bajty da się odłożyć i odsłuchać bez
rozpoznawania czegokolwiek; wybudzenie i nasłuch ciągły rozpoznają każdy
odcinek, więc znikają razem z silnikiem. Pole modelu odpowiedzi niesie model
zastany na dysku, a przy jego braku model ustawiony w konfiguracji: kontrakt
pyta o ten drugi, a oddanie pustki, gdy wag jeszcze nie pobrano, gubiłoby
nastawę widoczną w oknie konfiguracji.

Pomiar gotowości odsłuchu używa tej samej drogi doboru syntezatora, którą idzie
faktyczny odsłuch, więc pomiar nie może rozejść się z wykonaniem: jeżeli dobór
silnika kończy się odmową, odsłuch odmówi tak samo, a jego powód jest tym, co
Operator zobaczy.

Przekład odmów silnika mowy na kody kontraktu rozróżnia cztery przypadki na
trzy kody. Brak nagrania daje kod walidacji, bo jest jedyną odmową wywołaną
daną przysłaną przez klienta, a kod jest nieponawialny. Brak pomocnika daje kod
niedostępności kanału, bo żądanie było poprawne, a produkt mówi wprost, czego
dołożyć; kod jest ponawialny. Brak interpretera i brak silnika dają ten sam kod
niedostępności kanału, bo katalog kodów kontraktu nie ma pozycji odróżniającej
brak interpretera od braku biblioteki w nim; rozróżnienie, którego wymaga
naprawa, niesie treść odmowy, bo trzy ogniwa łańcucha naprawia się trzema
różnymi czynnościami: dołożeniem katalogu pomocników, instalacją interpretera
i instalacją biblioteki rozpoznawania. Naruszenie izolacji daje kod odmowy
uprawnień, tak samo jak znakuje je Terminal. Odmowa nierozpoznana schodzi na
kod błędu wewnętrznego, bo jest przypadkiem, którego rdzeń nie przewidział.

## budowa/server/internal/core/adapter_modul_monitor.go

Rejestr procesów rdzenia jest w pamięci, tabeli procesów nie ma. Proces okna
jest procesem systemowym objętym uchwytem rdzenia i ginie razem z rdzeniem,
więc wiersz, który przeżyłby restart ze stanem `running`, opisywałby proces
nieistniejący. Monitor czyta ten sam rejestr, który rozgłasza
`progress.changed` w `telemetria.go` — jedno źródło faktu, nie dwa.

Gdy port nie niesie telemetrii, monitor nie ma czego czytać i odmawia kodem
`internal_error`, nazywając brakujący byt; pusty wykaz procesów byłby wtedy
odpowiedzią nieprawdziwą. Rejestr wpięty i pusty to co innego — wtedy pusty
wykaz jest prawdą.

Monitor procesów jest oknem warstwy wspólnej, nie modułu, tak jak
`zdarzenia.go`. Warstwę wspólną — stronę główną, środowiska, moduły i stan
okna operacyjnego — obsługuje port Nawigacja, więc monitor osadza jego
adapter zamiast zakładać port równoległy.

Metoda `ZTelemetriaProcesow` wywołuje się po `ZObecnoscia`: chwilę ostatniego
zgłoszenia telemetrii zna pamięć czynności rejestru obecności, a licznik
obiegów zna rejestr biegów.

Zegar odpisów procesów pamięta chwilę, w której rdzeń stwierdził stan
procesu, którego chwili zmiany nie zapisał nikt inny; tę rolę pełni metoda
`chwilaZmiany`.

Proces wskazany wprost w `monitor.status`, którego rejestr nie zna, jest
bytem nieistniejącym: zamiast pustego wykazu idzie odmowa `not_found` z
nazwą procesu.

Dwa pola żądania `monitor.subscribe` mówią o dwóch różnych oknach:
`windowId` jest oknem odbierającym telemetrię, obserwatorem, a nie sitem
procesów — sitem są `processIds` (pusta lista znaczy komplet) oraz
`sessionId`. Tak samo rozstrzyga to bliźniacza komenda
`automation.execution.subscribe`. Żądanie bez `windowId` jest zwykłym
odczytem — nie ma czego zapisać, więc pole `subscribed` niesie `false`.
Zdarzenie `progress.changed` i tak dociera do wszystkich połączeń konta;
zapis mówi rdzeniowi, które okno których procesów pilnuje.

Obserwacja procesu, którego rejestr nie zna, jest obserwacją niczego.
Rejestr telemetrii procesów nie wykreśla, więc identyfikator nieznany
znaczy identyfikator nieistniejący, a nie proces właśnie domknięty.

Sesja procesu bywa uzupełniana z okna i dzieje się to przed sitem, nie po
nim. Telemetria bywa uboższa od stanu rdzenia: proces otwarty na
niepowodzeniu przyjęcia wiadomości zna okno, lecz nie zna jeszcze sesji
(`opisTury(z.WindowId, "")` w `telemetria_tury.go`). Sito po samym zapisie
telemetrii oddałoby pustkę przy procesie, który do wskazanej sesji należy.
Sesję okna rozstrzyga rejestr obecności — tą samą metodą, którą rozstrzyga
ją nasłuch telemetrii.

Metoda `Odpisy` oddaje odpis, a nie wskaźniki: proces bywa zmieniany w tej
samej chwili przez wątek tury albo kolejki, a czytelnik ma dostać stan
spójny, nie stan w połowie zmiany. Metoda mieszka w pliku monitora, bo to
monitor jej potrzebuje — producent telemetrii nie ma czytelników poza nim.

Chwilę ostatniej zmiany procesu rozstrzyga w pierwszej kolejności pamięć
czynności rejestru obecności: zapisuje ona znacznik przy każdym zgłoszeniu
telemetrii okna w `stan_sesji_czynnosc.go`, więc dla procesu okna jest to
dokładnie chwila ostatniej zmiany. Proces bez okna chwili zmiany nie ma
gdzie zapisanej: telemetria trzyma etap i stan, lecz nie czas, a pamięć
czynności jest kluczowana oknem i proces bez okna, na przykład kolejka
założona przed pierwszą wiadomością, do niej nie trafia. Dla takiego procesu
monitor podaje chwilę, w której rdzeń stwierdził ten stan — wartość będącą
górnym ograniczeniem chwili zmiany, nie nią samą.

Zegar odpisów procesów bez okna nie jest drugim rejestrem procesów: nie
trzyma ani etapu, ani stanu jako faktu, tylko znacznik przypięty do odcisku
odpisu. Stan niezmieniony znacznika nie przesuwa, więc kolejne odczyty
monitora nie odmładzają procesu, który stoi. Rejestr telemetrii procesów nie
wykreśla, więc bez granicy pojemności zegar rósłby razem z nim przez cały
czas życia rdzenia; po przekroczeniu granicy zegar zaczyna od nowa, zamiast
puchnąć.

Pamięć okien obserwujących procesy wiąże okno z procesami, których
telemetrię obserwuje. Pusty wykaz procesów znaczy obserwację kompletu — tak
mówi kontrakt komendy wprost. Zdarzenie `progress.changed` dociera do
wszystkich połączeń konta i rejestr niczego w tym nie zmienia. Rejestr
istnieje po to, by pole `subscribed` niosło prawdę, a rdzeń wiedział, które
okno których procesów pilnuje. Ten sam wzorzec prowadzi
`pamiecObserwatorowPrzebiegow` dla Execution Monitora.

## budowa/server/internal/core/adapter_modul_workspace_notatki.go

Odnośniki liczy zapis strony, nie jej odczyt. Przy każdym zapisie strony
rdzeń wyjmuje z treści odnośniki zapisane wzorem podwójnego nawiasu
kwadratowego i zapisuje je wierszami. Panel „co linkuje tutaj" pyta wtedy
jednym zapytaniem, zamiast przeszukiwać treść wszystkich stron projektu przy
każdym otwarciu. Odnośnik do strony jeszcze niezałożonej jest stanem
poprawnym wiki: nazwa czeka na stronę i domyka się sama w chwili jej
założenia. Dlatego zapis oddaje wykaz brakujących nazw — Operator ma widzieć,
które nazwy jeszcze nie mają strony, a nie odkrywać to po kliknięciu w martwy
odnośnik.

Usunięcie strony bez znacznika usunięcia stron podrzędnych przenosi je pod
stronę nadrzędną usuwanej — wiki nie gubi wtedy gałęzi razem z jej korzeniem.

## adapter_modul_design.go

Generowanie oddaje bajty obrazu, nie samą kopertę: komenda albo zapisuje
zasób o prawdziwej treści, albo odmawia zdaniem nazywającym brak rzeczywisty,
składany przy każdym wywołaniu osobno. Obrazu zastępczego, zasobu bez bajtów
ani powodzenia bez treści tu nie ma — każda droga bez bajtów kończy się
błędem.

Sejf poświadczeń jest magazynem sekretów rdzenia. Klucz dostawcy wchodzi do
sejfu drogą kont i punktów dostępu, a moduł Design ma go tylko odczytać. Brak
sejfu nie odbiera modułowi funkcji — dostawcy bez klucza pracują dalej,
a dostawcy z kluczem wracają jako pominięci wśród nieudanych.

Uruchamiacz procesów, rozstrzygacz konfiguracji i katalog roboczy służą
jednej czynności: czytnika liter w czystym Go nie ma, więc rozpoznanie pisma
idzie zewnętrznym programem, jedyną drogą w drzewie, która sprawdza obecność
programu i nakłada bramę izolacji. Uruchomienie stoi w jednym pliku obszaru
Design objętym nazwanym wyjątkiem zapory; nazwa wywołania programu nie ma
prawa być nawet w komentarzu, bo zapora czyta treść pliku, nie składnię. Brak
tych zależności nie odbiera modułowi układu makiety: rozpoznanie odmawia
wtedy zdaniem nazywającym brak, a obszary wychodzą jak dotąd.

Konstruktor adaptera wiąże (część) portu z repozytorium modułu wzorem
konstruktora modułu biblioteki. Sejf poświadczeń i magazyn biblioteki jadą
tą samą zmienną katalogu danych, żeby o katalogu była jedna prawda.

Pole PromptId zasobu kontraktu niesie identyfikator zewnętrzny promptu.
Repozytorium czyta go podzapytaniem obok wiersza zasobu, więc pole jest
prawdziwe albo puste — nigdy odgadnięte. Puste znaczy, że zasób nie powstał
z promptu: zasób wniesiony przez operatora promptu nie ma i mieć nie może.
Drugą stronę tego powiązania oddaje polecenie historii promptów.

Pole Uri wychodzi jako ścieżka względna magazynu, nie jako ścieżka na dysku.
Wiersz trzyma ścieżkę bezwzględną bloba, bo tą ścieżką rdzeń otwiera plik, ale
wypuszczenie jej odsłania układ katalogów maszyny. Przełożenie stoi
w jedynym miejscu składania zasobu kontraktu z wiersza, więc obejmuje
wszystkich wołających naraz i nie da się go ominąć nową drogą. Puste
odwołanie znaczy, że nie ma czego podać: wiersz, którego uri nie leży
w magazynie zasobów, nie dostaje pola zastępczego, bo kontrakt czyni je
niewymaganym właśnie po to.

Sprawdzenie rodzaju zasobu odrzuca wartość ósmą przed wejściem do kolumny.
Bez tego sprawdzenia wartość spoza wykazu odbijała się od warunku bazy
i wracała jako awaria rdzenia, którą wolno ponawiać, a żądanie z rodzajem
spoza wykazu nie uda się przy żadnym ponowieniu — klient z pętlą ponowień
powtarzałby je bez końca. To jest pomyłka wołającego i ma wracać jako pomyłka
wołającego; rdzeń stosuje to samo rozróżnienie w module Library. Odmowa
wymienia dopuszczalne rodzaje, zamiast cytować warunek schematu bazy wraz
z nazwą kolumny.

## budowa/server/internal/core/adapter_modul_library_uchwyty.go

Plik wpina dziesięć komend modułu Library na porcie `Biblioteka`, przez
który rejestr komend rdzenia dociera do adaptera złożonego z dwóch plików
tego samego typu `adapterBiblioteki`: `adapter_modul_library.go` niesie
plik i wersje, `adapter_modul_library_kolekcje.go` niesie kolekcje i
etykiety. Port wymienia wszystkie dziesięć komend niezależnie od tego, który
plik adaptera je implementuje: rejestr rdzenia potrzebuje jednego miejsca
wiążącego nazwę komendy z metodą portu — tak samo jak
`adapter_modul_automations_uchwyty.go` rejestruje w jednej funkcji komendy
swojego modułu.

Pole `LibraryTagSetRequest.collectionIds` ustawia komplet kolekcji pliku, ze
zdejmowaniem włącznie, metodą `UstawKolekcjePliku` z
`dane/biblioteka_kolekcje_pliku.go`. Pola `LibraryFile.collectionIds` i
`LibraryFile.versionId` wychodzą wypełnione metodą `a.zloz`, a
`library.file.search` sięga treści przez indeks pełnotekstowy FTS5.
Nieosiągalne pozostaje dopasowanie semantyczne, o którym mówi kontrakt
wyszukiwania: rdzeń dopasowuje słowa, nie znaczenia.

Cztery komendy rozgłaszają `library.file.changed`. Kontrakt niesie to
zdarzenie jako `shared.EventLibraryFileChanged`, a klient je subskrybuje,
więc każda komenda zmieniająca stan pliku repozytorium — wgranie jako
`created`, dołożenie i przywrócenie wersji oraz ustawienie etykiet jako
`updated` — rozgłasza plik po zmianie zaraz po udanym wykonaniu.
Rozgłoszenie jedzie tym samym emiterem rdzenia, co pozostałe zmiany
obszarów. Nieudana komenda nic nie rozgłasza.

Cztery komendy cyklu życia zasobu, niosące zasób po zmianie, rozgłaszają
`library.file.changed` dla każdego zasobu osobno, ponieważ Library Explorer
odświeża pozycje, a nie cały wykaz.

Metoda `plikBiblioteki` rozgłasza `library.file.changed`, czyli plik
repozytorium po zmianie. Jest to metoda emitera per moduł, wzorowana na
`przebiegAutomatyki`, zadeklarowana w tym pliku, ponieważ to obszar Library
nazywa własne zdarzenie. Plik nie jest bytem karty sesji, więc zdarzenie
idzie bez jej wskazania, a Library Explorer odświeża się ze strony głównej.
## budowa/server/internal/core/adapter_modul_badania_cytowania.go

Style są wkompilowane, nie doczytywane. Repozytorium CSL liczy około dwóch
i pół tysiąca plików XML i procesor, który je czyta, jest biblioteką
JavaScriptu. Oparcie cytowań o taki procesor znaczyłoby, że u Operatora,
który nie doinstalował środowiska JS, cytowanie nie działa wcale — a
cytowanie jest w module badawczym czynnością codzienną, nie ozdobą. Dlatego
pięć stylów najczęściej wymaganych (APA, MLA, Chicago, IEEE, Vancouver) jest
złożonych tutaj, w Go, i działa od pierwszego uruchomienia. Styl własny
Operatora dokłada się obok jako nazwany wariant.

Kontrola kompletności mówi, czego brak. `research.citation.check` nie mówi
„metadane niekompletne". Mówi, które pole brakuje przy którym źródle i skąd
da się je uzupełnić — bo to jest jedyna postać tej informacji, z którą
Operator może cokolwiek zrobić.
## budowa/server/internal/core/adapter_modul_orkiestracja_uklad.go

Cztery komendy pliku dopełniają układ zależności: `orchestration.gate.set`,
`orchestration.group.set`, `orchestration.compensation.set` i
`orchestration.multitasking.link`. Jadą tą samą maszynerią co
`adapter_modul_orchestration.go`: układ zależności ma w rdzeniu jednego
właściciela, więc metody stoją na adapterze modułu Automations, a nie na
własnym adapterze obok.

Łuk nie wyraża wszystkiego, co układ musi umieć powiedzieć. Bramka mówi, KIEDY
tory scalają się w jednym kroku — łuk mówi tylko, że się schodzą. Grupa mówi,
że zbiór kroków biegnie razem — łuk wiąże parami. Kompensacja mówi, co zrobić,
gdy przebieg pękł w pół — łuk o błędzie nie mówi nic.

Spięcie z MultitaskingAI nie jest znacznikiem. Silnik kolejek jest w rdzeniu
jeden i drugiego nie ma; środowisko MultitaskingAI odróżnia się od pętli
sesyjnej tym, czyim koordynatorem kolejka jest prowadzona i jakiego jest
rodzaju. Spięcie przestawia właśnie to na kolejkach automatyki, rozłączenie
zdejmuje. Zapis bez tego skutku byłby polem, które Operator przestawia, a
system ignoruje.
## budowa/server/internal/core/adapter_modul_przegladarka_pobrania.go

Pobranie naprawdę ściąga plik. Ponowienie (retry) idzie po treść spod adresu
pobrania i odkłada ją w magazynie modułu, a postęp w wierszu jest liczbą
bajtów, które na dysku leżą — nie deklaracją. Wstrzymanie i wznowienie
zmieniają stan kolejki, przerwanie ją kończy, zdjęcie usuwa wpis.

Makro zapisuje kroki w kształcie kontraktu (AutomationStep), tym samym,
którym jedzie moduł Automations. Dzięki temu przekazanie scenariusza do
Automations jest przełożeniem wiersza, a nie tłumaczeniem jednego kształtu na
drugi.

Granice Wykonawcy mają wartość domyślną w kodzie, nie w schemacie. Brak
wiersza znaczy granice domyślne rdzenia i tak też odpowiada odczyt — zamiast
odmawiać, że nikt jeszcze niczego nie ustawił.

Funkcja odlozPobranie ściąga zasób, którego nie da się pokazać jako strony,
i zakłada dla niego wiersz w menedżerze pobrań. Oddaje odmowę komendy
browser.navigate — bo migawki strony z tego nie ma — ale odmowa nazywa skutek,
który naprawdę zaszedł: pobranie o podanym identyfikatorze, z bajtami leżącymi
w magazynie. To jest jedyna droga, którą pobrania powstają, i jest to droga
naturalna: w przeglądarce plik pobiera się przez wejście pod jego adres,
a nie osobnym poleceniem dodaj pobranie — takiego kontrakt zresztą nie niesie.

### UstawBramke (adapter_modul_orkiestracja_uklad.go)

Ocena nie blokuje zapisu — tak samo jak przy `orchestration.dependency.set`.
Bramka na kroku, którego jeszcze nie ma, zapisuje się, a zastrzeżenie wraca
w odpowiedzi: Workflow Builder buduje układ krok po kroku i odmowa kazałaby
Operatorowi układać go w jedynej dopuszczonej kolejności.

### oknoRoliSpiecia (adapter_modul_orkiestracja_uklad.go)

Rola środowiska MultitaskingAI jest rolą nadaną oknu (`role.assign`), więc
wskazanie roli jest wskazaniem okna. Rola nierozpoznana wraca odmową: spięcie
z rolą, której nie ma, nie spięłoby niczego, a odpowiedź brzmiałaby udanie.
## budowa/server/internal/core/adapter_modul_roundtable_graf_wydanie.go

Wszystkie szesc formatow (DOT, GraphML, Argdown, AIF, SVG, PNG) sklada rdzen
sam, bez ani jednego programu z zewnatrz. Cztery pierwsze sa formatami
tekstowymi i pisze sie je wprost. SVG jest dokumentem XML, wiec tez. PNG
powstaje rysowaniem po mapie bitowej biblioteka standardowa - rasteryzator
zewnetrzny bylby zaleznoscia, ktorej instalka nie niesie, po to, zeby narysowac
prostokaty i podpisy. Uklad jest kolumnowy i wynika z tresci: wezly stoja
w kolumnach wedlug aktu mowy (teza, argument, kontrargument), wiec czytelnik
widzi strukture sporu, zanim przeczyta chocby jedno zdanie.

Format AIF: wezel tresci ma typ I (information), a relacja typ zalezny od jej
rodzaju - wsparcie RA (rule application), podwazenie CA (conflict application),
przeformulowanie MA (preference/restatement). Podzial pochodzi z samego
standardu AIF, nie jest oznaczeniem wprowadzonym w tym kodzie.

## budowa/server/internal/core/adapter_modul_auth_metody.go

Hasło jest kotwicą bramki i ta zasada przechodzi przez cały plik.
`auth.method.add` hasła nie zakłada (kotwica powstaje przy `auth.register`),
`auth.method.remove` hasła nie zdejmuje, a `auth.password.reset` zmienia je
wyłącznie ze znajomością hasła bieżącego — drogi odzyskania listem nie ma,
bo rdzeń poczty nie wysyła.

### ZalozMetodeWejscia

Rodzaj `password` odmawia, bo kotwicę zakłada `auth.register`. Rodzaj `hello`
odmawia, bo rdzeń go nie obsługuje — jawną odmową, nie cichym pominięciem. PIN
na urządzeniu, które PIN już ma, odmawia `conflict` zamiast dokładać drugi:
baza trzyma parę (urządzenie, rodzaj) jako jednoznaczną, a zmiana PIN-u to
zdjęcie starego i założenie nowego, czyli dwie jawne decyzje Operatora.
## budowa/server/internal/core/adapter_modul_workspace_zadania.go

Zmiana zadania jest łatą, nie podmianą: `workspace.task.update` zmienia
wyłącznie pola podane w żądaniu. Pole pominięte zostaje bez zmiany, bo okno
wysyła jedno pole na jedną czynność Operatora (zmiana stanu, przypisanie
wykonawcy, przesunięcie terminu), a podmiana całego zadania kasowałaby przy
każdej z nich to, czego akurat nie było na ekranie.

Przesunięcie następników po zmianie terminu zadania idzie wszerz z licznikiem
odwiedzin, żeby zależność zapętlona (gdyby powstała inną drogą niż komenda)
nie zawiesiła zapisu.

Klient ma z czego zdjąć pozycję z wykazu pobrań, zamiast zgadywać, która
zniknęła, bo odpowiedź komendy download.control niesie to, co zostało zdjęte.

Odczyt granic Wykonawcy oddaje wraz z brakiem wiersza wskazanie zasięgu, do
którego granice domyślne rdzenia się odnoszą.

Żądanie bez wskazania okna ani sesji dotyczy całej aplikacji — najszerszego
z dziewięciu poziomów zasięgu, tego, który ustępuje każdemu węższemu.

Rok pozyskania nie jest rokiem wydania i tak jest traktowany: wchodzi
wyłącznie jako data dostępu do zasobu sieciowego, a przy pozycji
recenzowanej zostaje brakiem, który zgłosi kontrola.

Źródła niecytowane w raporcie: wykaz literatury, której nikt nie użył, jest
drugą połową kontroli kompletności — pierwszą są braki metadanych.

### Zamek zmiany metody (ZalozMetodeWejscia, ZdejmijMetodeWejscia)

Sprawdzenie „wolno” i samo założenie idą pod jednym zamkiem: inaczej dwa
równoległe żądania na to samo urządzenie przechodzą oba sprawdzenie, a drugie
rozbija się dopiero o indeks bazy.

Metoda szybkiego wejścia bez kotwicy byłaby jedynym wejściem do platformy
i dałaby się zdjąć razem z urządzeniem — bramka zostałaby wtedy bez hasła.

Sekret bez właściciela byłby śmieciem w sejfie; niepowodzenie sprzątania
nie cofa zdjęcia metody, bo wiersza już nie ma.

### ZmienHasloBramki

Żądanie tej komendy tokenu nie niesie, więc sesję wołającego wskazuje więź
gniazda z sesją zawiązana w `wiez_polaczenia.go`. Połączenie niezwiązane
z żadną sesją traci wszystkie: skrót pusty znaczy „nie wiadomo, kto woła”,
a wtedy oszczędzenie którejkolwiek sesji byłoby zgadywaniem.

Sprawdzenie hasła bieżącego i podmiana sekretu idą pod jednym zamkiem.
Rozdzielone przepuszczają dwie równoległe zmiany, obie potwierdzone
`changed: true`, po których bramkę otwiera tylko jedno z dwóch nowych haseł:
drugi zapis nadpisuje pierwszy w sejfie, a Operator dostaje potwierdzenie
hasła, którym nie wejdzie.

### podmienSekret

Sejf nadpisuje wpis bytu i oddaje to samo odwołanie, więc zapis do bazy jest
tu asekuracją na wypadek, gdyby magazyn kiedyś zmienił postać odwołania —
nie drugą prawdą.

### PrzedluzSesjeBramki

Sesja bierze się wtedy z więzi tego gniazda (`wiez_polaczenia.go`), nigdy
z domysłu „jedyna czynna” — inaczej przedłużałoby się cudze wejście.
Połączenie niezwiązane — bieg wewnętrzny albo gniazdo, które nie przedstawiło
tokenu — dostaje odmowę opisującą właśnie ten brak.

Rdzeń trzyma wyłącznie skrót tokenu, więc odtworzyć tokenu nie może, a wpisanie
tam skrótu byłoby oddaniem klientowi napisu, którym nie da się wejść. Wołający,
który tokenu nie podał, ma go u siebie; z odpowiedzi bierze nowy czas
wygaśnięcia.

### trwanieSesji

Wynik przełącznika „nie wyloguj mnie” zostaje zapisany przy wierszu sesji
i odczytany z powrotem przy odnowieniu, zamiast brać stałą. Bez tego każde
`auth.token.refresh` — a klient woła je przy każdym uruchomieniu
(`client/src/uwierzytelnienie/ekran-logowania.ts`) — ścinałoby sesję roczną
do dwunastu godzin. Wiersz bez zapisanego trwania niesie zero i dostaje
trwanie podstawowe.

## budowa/server/internal/core/adapter_modul_studio.go

Jedna droga zapisu: document.save jest jedynym miejscem, które zapisuje treść
dokumentu — document.open tylko czyta albo zakłada wiersz pusty. Dwie ścieżki
zapisu tej samej treści byłyby dwiema prawdami o tym samym bycie.

Bez rejestru kanałów operacja kontekstowa — sedno modułu Studio — odmawia
kodem channel_unavailable, bo nie ma czym wywołać modelu. Rejestr okien jest
potrzebny, bo kanał modelu należy do okna: Studio pracuje w imieniu okna
komunikacji i ma sięgnąć po kanał, który Operator ustawił temu oknu.

Trzy zależności ZNarzedziami nie sprowadzają się do jednej: uruchamiacz
startuje proces, rozstrzygacz mówi, jakie zasady obowiązują okno, a katalog
wskazuje obszar, w którym proces wolno puścić. Ten sam komplet bierze adapter
narzędzi dokumentu pod nazwą ZIzolacja. Zależność jest opcjonalna: bez niej
czynności sięgające po arsenał odmawiają zdaniem nazywającym brak.

Katalog danych ZZasobami jest ten sam, nad którym stoi magazyn zasobów
Designu — wynik wydania Studia i wynik warsztatu PDF mają leżeć w jednym
miejscu, bo dwa magazyny znaczyłyby dwa katalogi, z których jeden prędzej czy
później zostałby przy kopii. Zależność jest opcjonalna: bez niej czynności
wydania, wyrysu i osadzenia odmawiają zdaniem nazywającym brak.

OtworzDokument nie doczytuje treści z Library przy wskazaniu pliku
repozytorium — zrobiłoby to z adaptera Studio klienta modułu Library, którego
konstruktor nie zna. Ścieżka urządzenia jest lokalna dla klienta; rdzeń nie ma
dostępu do systemu plików Operatora.

dane.ZapiszWersje zostawia wersja_biezaca_id nietknięty (patrz komentarz
w studio_wersje.go) i oddaje wywołującemu decyzję, kiedy nowa wersja staje się
bieżącą. document.save z createVersion=true zakłada wersję po to, żeby
Operator dalej edytował od niej, więc od razu staje się bieżącą. Inaczej niż
w PrzywrocWersje (studio_wersje.go), gdzie bieżącą staje się wersja wskazana
przez Operatora, a nie najświeższa zapisana.

Assets Panel i Preview Window widzą wynik po wierszu repozytorium Designu; bez
tej pary wydanie archiwum, wyrys strony i osadzenie grafiki nie mają gdzie
odłożyć tego, co zrobiły, a odczyt zasobu nie ma skąd wziąć materiału.
## budowa/server/internal/core/adapter_modul_terminal_wyjscie_dziennik.go

Byt osobny od adapter_modul_terminal_strumien.go: tamten rozsyla wyjscie
jednego procesu do okna, w ktorym stoi karta (WindowId: proces.oknoKod),
a tutaj zbiorcze wyjscie wszystkich kart odbiera okno wskazane w zadaniu,
wraz z ogonem historii. Okno Output Console nie musi byc oknem karty.

Dziennik nie czyta potokow procesow po raz drugi: owija nadajnik, ktorym
istniejaca pompa juz wysyla fragmenty (nadajnikZDziennikiem). Fragment idzie
swoja dotychczasowa droga nietkniety, a do dziennika i do okna obserwatora
wchodzi jego kopia. Dwie pompy na jeden potok czytalyby sobie nawzajem bajty.

Fragment nie jest wierszem: pompa czyta bloki po 32 KiB, a kontrakt niesie
TerminalOutputLine, czyli wiersz. Blok bywa urwany w polowie wiersza, wiec
resztka czeka na ciag dalszy, osobno dla wyjscia zwyklego i diagnostycznego
(dwie gorutyny, dwa niezalezne strumienie). Resztka bez znaku konca linii
dluzsza niz maksDlugoscWierszaWyjscia idzie do dziennika jako wiersz mimo
wszystko, bo pasek postepu nie konczy sie nigdy, a zuzycie pamieci ma
pozostac ograniczone.

Wierszy wyjscia nie ma gdzie zapisac: schemat bazy ma tabele terminal_karta
i terminal_proces, zadnej tabeli wyjscia. Dziennik jest wiec pierscieniem
w pamieci, tak samo jak pamieciowa warstwa historii procesow
(pojemnoscHistorii w adapter_modul_terminal_rejestr.go). Po restarcie ogon
jest pusty.

Kolejnosc w nadajnikZDziennikiem jest istotna: najpierw droga do okna karty,
potem dziennik, zeby skladanie wierszy i rozsylka do obserwatorow nie
opozniaja wyjscia karty.

## budowa/server/internal/core/adapter_modul_tlumaczenie_dokument.go

Wczytanie dokumentu robi trzy rzeczy naraz i wszystkie trzy są trwałe: zakłada
wiersz dokumentu, zapisuje jego segmenty i wstawia treść dokumentu jako tekst
źródłowy okna. Trzecia jest tą, dla której Operator w ogóle wczytuje dokument:
bez niej `target.add` nie miałby czego przetłumaczyć, a moduł meldowałby
wczytanie dokumentu, po którym okno zostaje puste.

### zapasDlugosciUkladu

Powyżej piętnastu procent tekst realnie wychodzi poza ramkę, w której stał
oryginał.

### WczytajDokument — wskazanie zasobu

Kontrakt dopuszcza wskazanie zasobu zamiast ścieżki. Rdzeń modułu Translate
nie ma dostępu do magazynu zasobów Designu, więc nazywa to wprost, zamiast
oddać pusty dokument z identyfikatorem donikąd.

### WczytajDokument — rozpoznanie pisma

Dokument bez warstwy tekstowej: rozpoznanie pisma idzie Tesseraktem
z arsenału serwerowego — jedyna droga do treści skanu, i droga wskazana
zasadą produktu (nagłówek `*_dokument_formaty.go`).

### WczytajDokument — podział okna

Podział okna idzie po segmentach dokumentu, nie po zdaniach: akapit dokumentu
jest jednostką, którą Operator widzi w pliku źródłowym.

### PorownajUklad — przesunięcie strony

Akapity nadmiarowe po stronie przekładu przesuwają treść dalej niż
w oryginale — to jest przesunięcie strony, nie brak.
## budowa/server/internal/core/adapter_modul_przegladarka_zaplecze.go

Magazyn bajtów sesji przeglądania jest ten sam co w Library i Design co do
mechaniki (blob pod sumą sha256, zapis niepodzielny), a inny co do miejsca:
bajty modułu Browser leżą w `<dane>/przegladarka/tresc`. Wspólny katalog
z biblioteką mieszałby materiał trwały (dokument wniesiony do repozytorium
wiedzy) z materiałem sesji (zrzut strony, archiwum, rejestr sieciowy) — a te
dwa mają różny cykl życia.

Silnik przeglądarki wchodzi tu jednym polem, nie jednym na rodzinę: Chromium
startuje ten sam dla zrzutu, dla drzewa DOM i dla konsoli, więc drugie pole
byłoby drugą prawdą o tym, czym rdzeń renderuje stronę.

ZMagazynem jest wołane przy montażu adaptera — bez niego rodziny wytwarzające
materiał (zrzut, archiwum, rejestr) nie mają gdzie odłożyć bajtów i mówią to
wprost, zamiast meldować powodzenie bez treści.

adresOstatniejStrony: komendy inspekcyjne kontrakt opisuje bez pola adresu —
pytają o „bieżącą stronę okna", a bieżącą stroną okna jest ostatnia jego
migawka.

Kontrakt nie ma kodu „brakuje programu"; `validation_failed` kłamałby o winie
żądania, a `internal_error` o usterce rdzenia. Treść odmowy `bladSilnikaPrzegladarki`
nazywa różnicę wprost: przy braku niesie nazwę programu i podpowiedź
instalacyjną z samego `zewnetrzne.BrakNarzedzia`.

Wstawienie wartości do wyrażeń (jakoLiteral, jakoLiczba, jakoLogiczna) chroni
przed wstrzyknięciem: selektor przychodzi z żądania, a selektor z apostrofem
albo z domknięciem nawiasu przerwałby wyrażenie i wykonał na stronie coś
innego, niż rdzeń napisał.
## budowa/server/internal/core/adapter_modul_aplikacje_podglad.go

Podgląd jest serwerem, który naprawdę stoi. Odpowiedź niesie adres, pod który
Operator ma wejść: adres wymyślony byłby wzorcem szkody, którego pilnują
sprawdziany skutku — odpowiedzią udaną, za którą nie ma niczego. Dlatego
uruchomienie podglądu podnosi nasłuch na pętli zwrotnej, oddaje pod nim treść
plików warsztatu wskazanej warstwy i zwraca adres wydany przez system
operacyjny.

Port wydaje system, nie konwencja. Nasłuch idzie na porcie zerowym pętli
zwrotnej, więc dwa okna podglądane naraz nie walczą o ten sam numer, a rdzeń
nie musi zgadywać, co na maszynie jest wolne. Adres jest zawsze pętlą
zwrotną: podgląd służy Operatorowi tej maszyny, a wystawienie warsztatu na
świat byłoby udostępnieniem kodu produktu bez czyjejkolwiek zgody.

Serwer nie przeżywa restartu rdzenia i wiersz w bazie tego nie ukrywa:
rejestr nasłuchów żyje w pamięci, a uruchomienie po restarcie podnosi nowy
serwer pod nowym adresem. Odtwarzanie nasłuchów przy starcie stawiałoby
podgląd, którego nikt w tej sesji nie zamówił.

Zatrzymanie zamyka nasłuch. Zatrzymanie podglądu woła zamknięcie serwera i
dopiero po jego powrocie melduje stan zatrzymany — inaczej pole potwierdzenia
znaczyłoby jedynie zgłoszenie zamiaru, a port zostawałby zajęty.

Treść obsługiwacza żądań serwera podglądu czyta się z bazy przy każdym
żądaniu, nie z migawki z chwili podniesienia: warsztat zmienia się w trakcie
pracy, a podgląd ma pokazywać stan bieżący.

Uruchamiacz, rozstrzygacz i katalog dają dostęp do programów zewnętrznych
rozpoznania pisma, zamiany formatów i pakowania. Brak pola form znaczy bez
zmiany postaci, a nie postać wyzerowaną — zwykły zapis treści nie ma prawa
zetrzeć arkusza stylów ani tabel. Wersje dokumentu założone przed dobudową
pola autora go nie mają, a podstawienie tam Operatora zamieniłoby brak wiedzy
w twierdzenie fałszywe dla każdej wersji, którą naprawdę zapisał model.
## budowa/server/internal/core/adapter_modul_workspace_pliki.go

Wydobycie nie zaklada drugiego warsztatu dokumentow: idzie ta sama droga co
document.text.extract - najpierw warstwa tekstowa dokumentu, a dopiero po
jej braku rozpoznanie pisma z pikseli. Dzieki temu pole method mowi prawde
o tym, skad wziely sie znaki, a nie o tym, czego rdzen probowal.

Duplikaty rozpoznaje tresc, nie nazwa. Grupe sklada suma kontrolna SHA-256
liczona z bajtow pliku (biblioteka wkompilowana, zadnego programu
z zewnatrz). Dwa pliki o roznych nazwach i tej samej tresci sa duplikatami;
dwa pliki o tej samej nazwie i roznej tresci nie sa nimi wcale.

Wykaz niczego nie scala: duplicate.list wskazuje plik proponowany do
zachowania - najstarszy w grupie - i na tym konczy. Scalenie jest osobna
komenda, bo usuniecie pliku z dysku jest czynnoscia nieodwracalna i ma byc
decyzja Operatora.

Podgląd pustej warstwy jest odmową, nie serwerem oddającym pustkę: Operator
ma się dowiedzieć, że nie ma czego pokazać, a nie oglądać białą stronę i
zgadywać, czy to wina warsztatu, czy nasłuchu.
## budowa/server/internal/core/adapter_modul_auth_konto.go

Droga potwierdzenia idzie listem, bo adres e-mail jest jedynym elementem
tożsamości niezależnym od urządzenia: PIN i klucz Windows Hello zostają na
maszynie, która mogła zginąć. W bazie leży wyłącznie skrót materiału wysłanego
listem, tą samą drogą co tokeny sesji, więc kopia bazy nie pozwala potwierdzić
cudzej tożsamości. Odpowiedź komendy auth.recover jest zawsze taka sama dla
adresu właściciela i dla adresu obcego, inaczej komenda byłaby wyrocznią
zdradzającą pytającemu adres Operatora, a jest osiągalna przed zalogowaniem.

Brak konta nadawczego nazywa się przed zapisaniem drogi w bazie, żeby baza nie
zostawała ze skrótem drogi, która nigdy nie wyszła listem — odmowa i tak mówi
o poczcie, tylko o jeden wiersz później.

Konto założone bez poczty (znacznik w sejfie) wchodzi hasłem, a nie drogą
listowną: pierwszy klient zostawałby przed platformą na zawsze, bo listu nie
ma skąd nadać, a drugiej rejestracji nie ma, bo wykonuje się raz. Adres
pozostaje wtedy niepotwierdzony i nikt tego nie udaje — wiersz konta stoi na
niepotwierdzone, a potwierdzenie czeka na pocztę.

## budowa/server/internal/core/adapter_modul_automations_alarmy.go

Zapis skarbca oddaje samą referencję, wykaz oddaje same referencje, a usunięcie
oddaje sam skutek. Wartość idzie do sejfu plikowego katalogu danych — tego
samego, którym jadą sekrety kont i punktów dostępu — i baza jej nie widzi, bo
kolumny na nią nie ma.

Usunięcie poświadczenia NIE jest wstrzymywane tym, że kroki je przywołują.
Kontrakt mówi to wprost: pole `referencingStepIds` nazywa kroki, które
straciły pokrycie, a usunięcie i tak następuje. Skarbiec, który odmawiałby
zdjęcia wykradzionego klucza, byłby skarbcem działającym przeciw właścicielowi.

Sejf kluczuje bytem, a referencja niesie przedrostek `sejf:` — rozbiera go
wołający, bo sejf zapisał sam byt.

## budowa/server/internal/core/adapter_modul_workspace_graf.go

Rodzina `workspace.knowledge` prowadzi wskaźnik ZNACZENIA i odpowiada na
pytanie „co jest podobne”. Ta rysuje sieć ODNOŚNIKÓW: kto na kogo wskazuje
wprost — strona na stronę, komentarz na zadanie, podzadanie na zadanie
nadrzędne. Krawędź istnieje wtedy, gdy ktoś ją napisał, a nie wtedy, gdy dwa
byty są sobie bliskie.

Graf przycięty granicą wielkości oddaje `truncated: true`. Bez tego pola
obraz częściowy wyglądałby na kompletny obraz projektu, a Operator
wnioskowałby o brakach powiązań z braku miejsca w odpowiedzi.

## budowa/server/internal/core/adapter_modul_orkiestracja.go

Podagent jest zadaniem w tle, którego tożsamością jest pozycja kolejki.
Wynikają z tego trzy rzeczy: powołanie nie uruchamia drugiego silnika — praca
idzie przez `adapterKolejek.Wykonaj`, tym samym silnikiem co pętla sesyjna
i Automations; powołanie nie czeka na wynik, stąd goroutine; stan przeżywa
restart rdzenia, bo stanem podagenta jest wiersz, nie pole struktury w
pamięci. Każdy podagent dostaje własną kolejkę, nie jedną wspólną: silnik
posuwa pierwszą czynną pozycję kolejki, więc piętnastu podagentów w jednej
kolejce jechałoby gęsiego. Granica piętnastu przycina, nie odmawia.

### rodzajKolejkiPodagenta

Więz CHECK kolumny `kolejka.rodzaj` (`store/migracja_003_kolejki.sql`
i `store/migracja_016_stan_kolejki_wyczerpana.sql`) dopuszcza wyłącznie
'sesyjna' i 'multitasking'; wartość spoza tego słownika kończyłaby każde
powołanie odmową bazy. Od kolejki etapu odróżnia ją nazwa równa
identyfikatorowi podagenta oraz okno koordynatora — po nazwie odnajduje ją
panel przy zatrzymaniu (`queue.list` → `queue.action` stop).

### Powolaj

Drogą przewidzianą jest wywołanie narzędzia `danaco_subagent_spawn` w trakcie
tury: serwer narzędzi (`server/internal/narzedzia`) uzupełnia wtedy
`windowId` oknem rozmowy, z którego wywołanie przyszło. Wpis dziennika mówi
raz na proces, czy ta droga jest w kontrakcie wpięta
(`podagenci/narzedzia_modelu.go`).

### Powolaj — zakres eksperta

Ekspert z wyłączonym Subagent Network nie powołuje ani jednego podagenta,
a jego granica przycina żądanie mocniej niż granica platformy — zapis,
którego nikt by tu nie przeczytał, byłby suwakiem bez skutku
(`straz_eksperta.go`).

### Powolaj — oznaczenie prowadzenia

Podagent nieoznaczony, a już pracujący, zostałby zamknięty jako sierota przy
najbliższym starcie. Błąd oznaczenia nie przerywa powołania, zostawia jedynie
ślad w dzienniku.

### Powolaj — rozgłoszenie

Rozgłoszenie idzie przed puszczeniem pracy w tle, żeby panel zobaczył
podagenta `pending`, zanim praca przestawi go na `running` — inaczej dwa
zdarzenia mogłyby dojść w kolejności odwrotnej do faktów.

### puscWTle — nadzór

Nadzór mówi o procesie orkiestratora, czyli okna, które powołało. Proces
samej pozycji podagenta nie ma wpisu w rejestrze procesów
(`podagenci/zywotnosc.go`).

### puscWTle — kontekst

Kontekst własny na podagenta daje `subagent.stop` uchwyt do jednej pracy.
Wisi na życiu rdzenia, więc zatrzymanie rdzenia nadal zabiera wszystkich,
a odwołanie pojedyncze zabiera wyłącznie tego jednego
(`adapter_modul_orkiestracja_zatrzymanie.go`). Odwołanie zwalnia się zawsze:
inaczej praca zakończona zostawiałaby po sobie kontekst bez odbiorcy, a wykaz
prac rósłby z każdym powołaniem.

### wykonaj

Gdyby rdzeń padł pomiędzy zapisem pozycji a wejściem w stan `running`,
zostaje podagent `pending` z gotową pozycją, czyli praca do podjęcia;
odwrotna kolejność zostawiałaby podagenta „w biegu" bez pozycji, która ten
bieg niesie.

### przepiszStanPozycji

Dlatego stan bierze się stąd, gdzie naprawdę powstał.
## budowa/server/internal/core/adapter_modul_agents_warstwy.go

Adapter agentów składa się w montaz_porty.go bez ogniwa .ZWarstwami(...), więc
pole warstwy bywa nil; każda z czterech komend odmawia wtedy, nazywając komendę,
powód i miejsce wpięcia, zamiast panikować albo udawać zapis. Komendy
agent.list i agent.create nie przechodzą tędy — jadą dalej ekspertKontraktu
i przy warstwy == nil oddają eksperta bez warstw i bez wtyczek. Nazwy warstw
sprawdzane są wobec wartości kontraktu (shared.IdentityLayer, shared.IdentityMode)
przed zapisem, żeby powodem odmowy było zdanie po polsku, a nie naruszony
warunek CHECK.

Tożsamość osi (identity.document.set) tryb wybiera i domyślnie zastępuje, bo
jest konfiguracją platformy, a nie warstwą nałożoną na pojedyncze wywołanie.
Wartość musi pochodzić z kontraktu, nie z silnika nakładki: kolumna
agent_warstwa.tryb ma warunek CHECK na ZASTAP albo DOLACZ, a stała silnika
injection.TrybDopisz jest napisem "dopisz"; przekład między jednym a drugim
robi trybSilnika.

Ekspert bez wtyczek oddaje tablicę pustą — brak wtyczek jest poprawnym stanem,
a nie awarią odczytu; odmowa idzie wyłącznie wtedy, gdy eksperta o wskazanym
kodzie nie ma w katalogu. Odczyt tożsamości pełnej nie odmawia z powodu
niewpiętego katalogu warstw: czytać nie ma czego, a to inny przypadek niż
zapis, który nigdzie nie trafia.

Osobna czynność repozytorium zapisu tożsamości dotyka wyłącznie kolumn
imie_wlasne i favikon, zamiast zmieniać drogę, którą jadą wszystkie pozostałe
pola. Pominięte pole nie jest polem pustym: nil zostawia wartość zastaną,
pusty napis czyści ją i zostaje zapisany. Zlanie obu w jedno kasowałoby imię
przy każdej zmianie samego opisu. Brak wpiętego repozytorium nie jest tu
odmową — ekspert bez imienia własnego jest ekspertem, bo tożsamość niesie
pole nazwa.

Zbieranie tożsamości wykazu idzie dwoma zapytaniami na wywołanie, nie dwa na
eksperta: odczyt po jednym dałby przy stu ekspertach dwieście zapytań na jedno
otwarcie biblioteki. Wzorzec ten sam co dolaczPowiazania w warstwie danych.
Błąd odczytu nie wywraca wykazu — ekspert, którego warstw nie udało się
odczytać, wchodzi do wykazu bez warstw, tak samo jak ekspert, który ich nie
ma; wykaz ekspertów ma się pokazać także wtedy, gdy tożsamość jest chwilowo
nieczytelna.

Pominięte pole trybu nałożenia zostawia wartość zastaną — nil nie znaczy
powrotu do domyślnego trybu; raz oznaczone odstępstwo nie ma prawa zniknąć
przy zmianie samego opisu eksperta. Wartość spoza katalogu odrzuca warunek
CHECK na kolumnie; drugiej listy dopuszczonych trybów tu nie ma. Brak
wpiętego repozytorium nie jest tu odmową: ekspert bez zapisanego trybu
dopisuje się do promptu globalnego.
## budowa/server/internal/core/adapter_modul_tlumaczenie_model.go

Metody stoja na wspolnym *adapterTlumaczenia, ktorego typ, konstruktor
i przedrostki deklaruje adapter_modul_tlumaczenie.go. Droga modelu (przeklad
panelu, rozpoznanie jezyka) to jedna odpowiedzialnosc wolana z trzech komend
(target.add, backtranslation.run, source.detect), wiec lezy w jednym pliku
i ma jeden opis granicy: brak kanalu jest odmowa wprost, nigdy pustym
napisem udajacym przeklad.

Kanal wskazuje Operator polem channelId, a rozstrzyga to kanalZadania: brak
wskazania to kanal domyslny czynny; wskazanie kanalu czynnego i gotowego
idzie tym kanalem; wskazanie kanalu nieznanego albo nieczynnego to odmowa
nazwana. Ciche zejscie na kanal domyslny wypelniloby panel przekladem
modelu, ktorego Operator nie wybral, i nic by o tym nie powiedzialo.

ZWyjsciem: wpiecie robi zarejestrujTlumaczenie (ten sam emiter, ktory
dostaja pozostale moduly w kompozycja.go), wiec montaz portow nie musi znac
tej zaleznosci. Nadajnik niepodlaczony nie jest bledem: rdzen tlumaczy takze
wtedy, gdy nikt nie sluchaz zdarzen - emiter.wyslij sam odsiewa pusty
nadajnik.

kanalZadania: odmowy sa trzy i kazda mowi o czym innym - kanal nieznany,
kanal znany lecz wylaczony, kanal wlaczony lecz bez zbudowanego adaptera.
Naprawia sie je trzema roznymi ruchami (poprawic identyfikator, wlaczyc
kanal, poprawic konfiguracje), wiec kazda ma wlasne zdanie. Kanal gotowy
znaczy to samo co w Rejestr.Kontrakt(true): wiersz czynny z zbudowanym
adapterem. Wiersz aktywny=1 bez adaptera odmowi przy pierwszej turze, wiec
lepiej powiedziec to teraz niz w polowie przekladu.

domyslnyKanalModelu: pierwszy wiersz rejestru czynny i gotowy do pracy,
wzor adapter_kolejki.go: rozwiazKanalPozycji. Brak rejestru albo brak
czynnego kanalu znaczy nie ma czym wolac modelu.

przetlumaczModelem: slownik Operatora wchodzi dwa razy - raz do tresci
polecenia (poleceniePrzekladu dostaje wiazania), raz po odpowiedzi modelu
jako mechaniczna podmiana terminow zostawionych w brzmieniu zrodlowym
(zastosujTerminySlownika, ta sama funkcja co w glossary.apply). Uzasadnienie
obu drog niesie naglowek adapter_modul_tlumaczenie_polecenia.go. Nieudany
odczyt slownika odmawia calego przekladu. Ton panelu tez idzie do polecenia,
inaczej kolumna ton (panel.tone.set) bylaby zapisem bez skutku.

Siatka bezpieczenstwa slownika: podmiana terminow nietykalnych byloby
zlamaniem zakazu Operatora, a nie jego pilnowaniem. Liczba podmian nie idzie
nigdzie dalej: kontrakt target.add nie ma pola na taka liczbe, a
ChangedCount nalezy do glossary.apply, nie do przekladu.

przetlumaczZwrotnieModelem: slownik tu nie wchodzi (powod przy
polecenieTlumaczeniaZwrotnego), bo kontrola wiernosci, ktora sama naprawia
terminologie, niczego nie kontroluje.

rozglosZmianePanelu: nadajnik niepodlaczony jest odsiewany w emiter.wyslij,
a nil-emiter w jego odbiorniku nil, wiec wywolanie jest bezpieczne bez
wpietej szyny.

## budowa/server/internal/core/adapter_modul_agents_przeklad.go

Identyfikatorem kontraktu jest kod wiersza: pole Agent.id niesie agent.kod,
nie numer wiersza. Numer żyje wyłącznie wewnątrz bazy i do klienta nie
wychodzi, bo przy przeniesieniu bazy przestałby się zgadzać.

Piąta wartość poziomu pamięci pozwoliłaby przysłać zestaw sprzeczny
("session","disabled"); zbiór pusty tego wyrazić nie umie, dlatego
wyłączenie pamięci jest zbiorem pustym, a nie piątą wartością.

Widocznosc nierozpoznana czyta się jako global, bo ekspert, którego
widoczności nikt nie rozpoznaje, ma się pokazać, a nie zniknąć z biblioteki.
Katalogu wartości tu nie ma — pilnuje go warunek CHECK w bazie.

Tryb pusty warstwy promptu zostaje pominięty — kontrakt opisuje tryb
domyślny brakiem wartości, nie napisem. Wartość trybu nałożenia spoza
katalogu czyta się jako dołączenie, bo baza pilnuje warunku CHECK, a funkcja
trybKontraktu jest drugą siatką, nie drugim katalogiem wartości.

Punkt dostępu konektora wychodzi kodem trwałym, tym samym, którym posługują
się komendy access.point.* i nadania okna rozmowy.

Ekspert niesie wskazania wtyczek, a nie ich treść — pełna struktura wtyczki
wraca wynikiem agent.plugin.add, dokładnie jak konektor przy
agent.connector.add.

Tryb nałożenia oddawany zawsze, także gdy jest domyślny: okno ma pokazać, że
ekspert dopisuje się do promptu globalnego, a pominięcie pola zostawiłoby
domysł zamiast odpowiedzi. Ta sama zasada dotyczy poziomów pamięci i modułów
zastosowania — pominięcie pola zostawiłoby domysł zamiast odpowiedzi.
## budowa/server/internal/core/adapter_modul_library_metadane.go

Opis stoi w tabeli towarzyszącej zasobowi, a nie w kolumnach wykazu:
piętnaście pól Dublin Core obciążałoby każdy odczyt Library Explorera, który
opisu nie pokazuje.

Zapis scala domyślnie, podmienia na żądanie. Różnica jest widoczna wprost:
przy scalaniu pole pominięte w żądaniu zostaje takie, jakie było, a pole
przysłane puste jest kasowane — bo inaczej Operator nie miałby jak wyczyścić
raz wpisanej wartości. Przy podmianie opis staje się dokładnie tym, co
przyszło.

Pola niestandardowe są mapą kod-wartość i przechodzą przez surowy zapis JSON
kontraktu. Rdzeń ich nie tłumaczy na kolumny: definicja pola należy do
Operatora, więc kolumna na pole znaczyłaby migrację przy każdym polu.

Metadane osadzone w pliku (EXIF, IPTC, XMP, ID3) czyta się z bajtów, więc
wchodzą wyłącznie na wyraźne żądanie — tak mówi kontrakt i tak działa ten
odczyt: bez żądania technicznego bajty zasobu nie są w ogóle otwierane.

Liczba zasobów z wartością pola liczy się przed zdjęciem definicji: po
usunięciu wartości zostają przy zasobach, więc liczba mówiłaby to samo, ale
kolejność ma znaczenie przy zakładaniu — Operator ma zobaczyć, ile zasobów
pole zastanie już wypełnione.

Pole przysłane puste kasuje wartość, pole pominięte przy scalaniu ją
zostawia — dlatego przekład idzie po wskaźnikach, a nie po wartościach:
brak wartości znaczy nieodnoszenie się do pola, pusty łańcuch znaczy chęć
jego wyczyszczenia.

Scalanie pól niestandardowych idzie po kluczach: klucz przysłany z wartością
pustą znika, klucz pominięty zostaje. Przy podmianie mapa staje się
dokładnie tą przysłaną.

Odczyt opisu jest dostępem do zasobu — dziennik audytu ma odpowiadać na
pytanie kto to oglądał, więc zaglądanie też zostawia ślad.

Zapis zastany nieczytelny nie może zablokować zapisu nowego: wartość
uszkodzona ustępuje wartości przysłanej.

## budowa/server/internal/core/adapter_modul_przegladarka_monitory.go

Założenie monitora POBIERA stronę i odkłada jej treść jako odniesienie;
sprawdzenie pobiera ją ponownie i zestawia obie treści wiersz po wierszu.
Monitor bez odniesienia oddawałby zawsze „bez zmian” albo zawsze „zmiana” —
jedno i drugie jest meldunkiem bez pomiaru, a to wzorzec szkody, którego ten
produkt już raz doświadczył.

`BrowserContentDiff` niesie liczbę wierszy dodanych, usuniętych, numer
pierwszego wiersza różnicy i przyrost znaków. Wszystkie cztery liczone są
z dwóch treści, nie z niczego.

Próg zmiany odsiewa drgania. Strona z zegarem albo licznikiem odwiedzin różni
się przy każdym pobraniu; próg podany w znakach mówi, od jakiej różnicy zmiana
jest zmianą. Bez progu każdy taki monitor alarmowałby co godzinę.

## budowa/server/internal/core/adapter_modul_workspace_harmonogram.go

Zależność domykająca cykl nie zostaje zapisana. Powód nie jest formalny:
cyklu nie da się ułożyć w czasie, więc zapisany cykl unieruchomiłby
przesuwanie terminów przy każdej późniejszej zmianie zadania.

Słupek wymaga początku i końca, więc zadanie bez obu granic do wykresu nie
wchodzi. Wchodzi jednak do pola unscheduledTaskIds — bez tego pola zadanie
znikałoby z widoku bez śladu, a Operator nie miałby skąd wiedzieć, że
harmonogram pokazuje mniej niż projekt.

Rachunek ścieżki krytycznej chroni przed zapętleniem licznikiem odwiedzin,
gdyby cykl powstał inną drogą niż komenda ZalozZaleznosc.
## budowa/server/internal/core/adapter_modul_library_porownanie.go

Kontrakt rozstrzyga porownanie dwoch zasobow albo dwoch wersji jednego
zasobu brakiem prawej strony - zadanie bez rightFileId porownuje wersje
zasobu lewego. Dokument binarny porownuje sie po tekscie z niego wydobytym
i odpowiedz mowi o tym wprost (comparedAsText). Operator ma wiedziec, ze
nie porownano bajtow: dwa dokumenty PDF o identycznej tresci roznia sie
bajtami w kazdej linii, wiec porownanie bajtowe oddaloby wszystko inne
i byloby bezuzyteczne.

Roznice liczy najdluzszy wspolny podciag wierszy (algorytm LCS) - ten sam
sposob, ktorym idzie porownanie w module Studio. Rdzen nie wola programu
diff: to byloby zaleznoscia od cudzego programu w czynnosci, ktora kod robi
sam.

Wydobycie tekstu z dokumentu PDF idzie programem pdftotext, ktory stoi na
serwerze razem z rdzeniem i jest wolany jedyna dozwolona droga
(zewnetrzne.Wolaj). To jest ODCZYT, nie przetwarzanie dokumentu: sam PDF
modul rusza wylacznie biblioteka wkompilowana.

Fragment obecny po obu stronach w innym brzmieniu wychodzi jako JEDNA
roznica rodzaju changed, a nie jako para usuniecie-dodanie: para kazalaby
czytelnikowi samodzielnie skojarzyc, ze to ten sam fragment.
## budowa/server/internal/core/adapter_modul_workspace_tablica.go

Przeciągnięcie karty między dwie sąsiednie dopisuje klucz leżący pomiędzy ich
kluczami — jeden zapis. Numer pozycji wymagałby przepisania całej kolumny przy
każdym przeciągnięciu, a tablica projektu bywa przeciągana kilkanaście razy
pod rząd. Przekroczenie granicy WIP wraca polem wipExceeded, a karta i tak
staje w kolumnie — platforma nie stawia twardych blokad w interfejsie, granica
jest sygnałem dla Operatora, nie bramką.

Klucz porządkowy karty składa się z liter a-z; wynik jest zawsze większy od
lewego i mniejszy od prawego klucza sąsiada, a przy sąsiadujących literach
schodzi o znak niżej zamiast oddawać klucz równy któremuś z sąsiadów.
## budowa/server/internal/core/adapter_modul_library_reguly.go

Reguła jest warunkiem wraz z tym, co ma się stać z zasobem, który go
spełnia. Trzy rodzaje różnią się chwilą zastosowania, nie kształtem:
kolekcja inteligentna przelicza zawartość na żądanie, reguła napływu
stosuje się przy wejściu zasobu do repozytorium, folder obserwowany wciąga
pliki spod ścieżki.

Przeliczenie jest tutaj, a nie w warstwie danych, bo to rdzeń zna znaczenie
członów warunku. Warunek niesie moduł źródłowy, etykiety, typ zawartości,
zakres dat oraz frazę wyszukiwania, w postaci zbliżonej do zapisu JSON.
Każdy człon jest opcjonalny, a człony łączą się koniunkcyjnie — tak samo
jak filtry fasetowe Library Explorera, bo to ten sam sposób zawężania
widziany z drugiej strony.

Przeliczenie nie opróżnia kolekcji: dokłada zasoby spełniające warunek ze
znacznikiem pochodzenia reguły. Zdjęcie tych zasobów ma własne żądanie —
regułą wolno dodać, a zabrać wyłącznie na wyraźne polecenie.

Folder obserwowany nie jest przeliczany przy żądaniu: wciąganie plików
spod ścieżki jest napływem z zewnątrz, a nie przeglądem repozytorium.
Reguła obserwacji zapisuje się i czeka na napływ, więc jej przeliczenie
oddaje zero i tak jest uczciwie — zamiast udawać, że coś zrobiono.

Ponowny odczyt niesie czas przeliczenia zapisany przy regule — bez niego
odpowiedź mówiłaby o regule sprzed własnego przebiegu.

Czas przeliczenia zapisuje się przy regule, żeby wykaz reguł mówił, kiedy
reguła ostatnio pracowała — inaczej stan czynny znaczyłby tylko włączenie.

Zapis daty jest porównywalny leksykograficznie, więc data podana samym
dniem działa jako granica bez przekładu na czas. Granica górna obejmuje
cały wskazany dzień: znacznik zasobu przycina się do długości granicy,
więc data końcowa nie odcina zasobu powstałego tego dnia.

## budowa/server/internal/core/adapter_modul_model.go

Rodzina `model.*` to jedna komenda kontraktu, `model.channel.set`. Kanały
zakłada i zmienia rodzina `channel.*`; ta komenda niczego nie zakłada,
wyłącznie wskazuje jeden z wierszy rejestru. Metody stoją na `adapterOkien` —
tym samym bycie, którym jedzie `window.update`. Kanał okna ma w rdzeniu
jednego właściciela: gdyby rodzina `model.*` dostała własny adapter nad
własnym dojściem do rejestru okien, ten sam kanał miałby dwa miejsca zmiany
i dwie prawdy o tym, który wygrywa.

Kanał okna mieszka w dwóch miejscach naraz i tylko jedno z nich ma wpływ na
wywołanie modelu: `okno_komunikacji.kanal_modelu_id` — kolumna wiersza okna
oraz jej odpowiednik pamięciowy `session.Okno.KanalModelu`, która jedzie do
wywołania tury i którą odbudowuje `odtworzenie_stanu.go` po restarcie; oraz
`ustawienie.klucz = 'kanal_modelu'` — pozycja katalogu ustawień, ustawialna
na wszystkich ośmiu poziomach zasięgu, w tym `karta_sesji` i `okno`, której
rozstrzygacz konfiguracji nie czyta — żaden kod nie sięga po ten klucz,
`konfig/definicje_wykonania.go` tylko go deklaruje. Komenda zapisuje więc
miejsce pierwsze — to, które działa. Zapis do pozycji konfiguracji byłby
ciszą udającą skutek: Operator przestawiłby kanał, panel pokazałby zmianę,
a tura poszłaby starym kanałem.

### UstawKanalModelu

Rozgłoszeniem `window.changed` zajmuje się wpięcie (handlers_model.go), tak
samo jak przy `window.update`.

### kanalRejestru

Kolumna `okno_komunikacji.kanal_modelu_id` jest kluczem obcym z ON DELETE
RESTRICT, więc zapis kodu nieistniejącego i tak by nie przeszedł, a okno
w pamięci wskazywałoby kanał widmo. Odmowa niesie kod `not_found` i nazwę
bytu. Wykaz idzie po wszystkie wiersze, także nieczynne: kanał wyłączony
istnieje i Operator ma prawo go wskazać, a odpowiedź mówi o tym wprost polem
`enabled`. Odmowa wyboru kanału wyłączonego byłaby zasadą, której kontrakt
nie stawia.

### kanalOkna

Pominięcie eksperta dawałoby Operatorowi wybierającemu eksperta z menu
„Modele" potwierdzenie bez skutku: kanał by się zmieniał, ekspert nie.
Wskaźnik pusty zostawia wybór bez zmiany, wskaźnik na pusty napis zdejmuje
eksperta — znaczenie jest jedno i pochodzi z `session.Zmiana`, żeby ta sama
wartość nie znaczyła tu czegoś innego niż w `window.update`.

### kanalKartySesji

Dwa inne znaczenia są tu świadomie nieprzyjęte: zapis pozycji `kanal_modelu`
na poziomie zasięgu `karta_sesji` — pozycja istnieje w katalogu ustawień,
lecz nikt jej nie czyta, więc Operator dostałby potwierdzenie bez skutku;
oraz „kanał domyślny dla okien zakładanych później" — okno zakładane bez
wskazania bierze kanał z rejestru (`adapterOkien.kanalDomyslny`), a nie
z karty; druga reguła domyślności zrobiłaby z jednego pytania dwa źródła
prawdy. Okna zamknięte zostają nietknięte: zamknięte okno nie prowadzi
pracy, a jego kanał jest zapisem tego, czym pracowało. Karta bez ani
jednego otwartego okna nie ma czemu nadać kanału — odpowiedź „kanał
ustawiony" byłaby tu ciszą udającą skutek.

### kanalOdpowiedzi

`channel.list` podaje czas utworzenia wypełniony (rejestr kanałów warstwy
modeli), więc zero w odpowiedzi tej komendy przeczyłoby tej samej wartości
na tym samym ekranie. Znacznik nieczytelny zostawia zero — dana pomocnicza
nie wywraca odpowiedzi.

### dopiszZmianeKanalu

Nadanie kanału karcie znaczy „nadaj go każdemu oknu, które ta karta
prowadzi"; wyłączenie eksperta z tej samej reguły zrobiłoby z jednego
żądania dwa zasięgi i Operator nie miałby jak zgadnąć, który obowiązuje.

### utrwalKanalOkna

Brak wiersza nie jest błędem. Wiersz okna powstaje leniwie, przy pierwszej
wiadomości (`dane/rozmowa_lancuch.go`), i bierze kanał wprost z opisu okna
w rejestrze rdzenia — czyli z wartości ustawionej właśnie teraz. Zapis nie
ma więc czego dogonić. Błąd zapisu jest błędem komendy: wiersz istnieje,
a nie przyjął zmiany — wybór nie przeżyje restartu i milczenie o tym byłoby
obietnicą bez pokrycia.
## budowa/server/internal/core/adapter_modul_queue.go

To rozszerzenie adaptera, nie drugi adapter. Metody wisza na
adapterKolejek z adapter_kolejki.go, wiec jada tym samym repozytorium, tym
samym silnikiem wykonania i tym samym przekladem kolejki na kontrakt, co
queue.create i queue.action. Drugiego silnika kolejek nie ma nigdzie.

Okna kolejki maja dwa zrodla: okna podane przy queue.create leza w pamieci
powiazan rdzenia (pamiec_sesji_kolejek.go), a okna podane przy queue.link
w tabeli powiazanie_kolejki. Kolejka kontraktu oddaje sume obu - inaczej po
ponownym uruchomieniu rdzenia okna z queue.create znikalyby bez sladu,
a okna z queue.link nie bylyby widoczne w ogole.

Interfejs repozytoriumWiazanKolejek stoi po stronie czytelnika. Deklaracja
mieszka w tym pliku, a nie w dane.RepozytoriumKolejek, bo wymagaja jej
wylacznie dwie komendy tej rodziny; port kolejek tych czynnosci nie zna
i nie musi.

Wykaz: sesja kolejki bywa znana wylacznie z pamieci powiazan (kolejka
zalozona przed pierwsza utrwalona wiadomoscia sesji nie ma czym wypelnic
kolejka.sesja_id). Zadanie bez zadnego warunku jest zgodne z kontraktem -
wszystkie trzy pola sa opcjonalne, zwraca sie wtedy wszystkie kolejki, nie
odmowe. Sesja bez kolejek to prawdziwa odpowiedz zero kolejek, nie nie ma
takiej sesji - sesja zyje w pamieci nadzorcy i bywa bez wiersza w bazie.

Zwiaz: cztery pola wiazania sa opcjonalne z osobna, ale wszystkie naraz
puste znacza zadanie bez tresci - odpowiedz kolejka po powiazaniu byloby
wtedy kolejka po niczym. Kod odmowy to validation_failed. Kontrakt zna
wylacznie wiazanie - komendy rozwiazujacej nie ma, wiec queue.link niczego
nie zdejmuje. Powiazanie powtorzone nie jest drugim faktem i nie jest
bledem; zapis je pomija. Istnienia bytu wiazanego rdzen tu nie sprawdza:
adapter kolejek nie ma dostepu do repozytoriow ekspertow, projektow ani
automatyk, a kolejka nie moze zalezec od tego, czy inny modul zdazyl sie
utrwalic. Zapisany zostaje identyfikator w ksztalcie, w jakim przyszedl.
## budowa/server/internal/core/adapter_modul_roundtable_glosowanie.go

Wariant wskazuje wypowiedź, gdy Operator podał jej kod. Głosowanie nad
stanowiskami debaty ma prowadzić od wyniku z powrotem do słów, które ten
wynik wywołały — wariant będący samą etykietą urywa tę drogę, więc kod
wypowiedzi rozpoznaje się i zapisuje, kiedy tylko padnie.

Głosowanie zamknięte odmawia: głos oddany po rozstrzygnięciu przestawiłby
wynik, który już ogłoszono. Wyborca spoza wykazu uprawnionych odmawia
osobno — uprawnienie zawężone i nieegzekwowane byłoby ustawieniem bez
skutku.

Wynik zerowy przy zerowej liczbie głosów wyglądałby jak rozstrzygnięcie,
którego nikt nie podjął.

Głosowanie nad jednym wariantem ma wynik znany przed oddaniem pierwszego
głosu.

Wariant podany kodem wypowiedzi wskazuje ją wprost; etykietą zostaje
wtedy treść tej wypowiedzi, żeby panel nie pokazywał samego kodu.

Remis znakuje się w stanie głosowania, bo Voting & Evaluation Center
pokazuje stan przy nagłówku, zanim Operator otworzy wynik.
## budowa/server/internal/core/adapter_modul_roundtable_agregacja.go

Remis jest wynikiem, nie usterka. Kazda z pieciu metod potrafi nie wylonic
zwyciezcy. Kontrakt przewiduje to wprost: winnerOptionId jest polem
niewymaganym, a stan glosowania ma wartosc tied. Zwyciezca dopisany bo
trzeba - pierwszy z brzegu przy rownej liczbie glosow - bylby
rozstrzygnieciem wymyslonym przez rdzen.

## budowa/server/internal/core/adapter_modul_agents.go

Model bazowy leży w `adapter_modul_agents_model.go`, umiejętności, konektory
i uprawnienia w `adapter_modul_agents_zasoby.go`. Kanał bazowy wskazuje kod
wiersza `kanal_modelu` — tego samego rejestru, z którego korzysta okno
rozmowy. Adapter sprawdza wskazanie w rejestrze kanałów rdzenia i nie
przechowuje własnych definicji dostawcy; pole `kanaly` służy tylko temu
sprawdzeniu i podpowiedzi modelu domyślnego, adapter nie zakłada kanałów.

Zmiana tożsamości jest częściowa: `agent.update` niesie same pola zmieniane,
a pole pominięte zostaje takie, jakie było. Bez tego okno musiałoby odesłać
komplet tożsamości przy każdej poprawce nazwy i skasowałoby instrukcje
systemowe pierwszym niepełnym żądaniem.

Pole `punkty` jest katalogiem punktów dostępu; konektor rodzaju `mcp`
wskazuje most z tego katalogu, ten sam, z którego rdzeń składa `mcpServers`.

Pole `warstwy` jest repozytorium tożsamości własnej eksperta — warstw jego
promptu i jego wtyczek; `nil` znaczy „nie wpięto”, a wtedy cztery komendy
warstw odmawiają, a reszta modułu pracuje bez zmiany.

Pominięte pole pamięci przy założeniu eksperta znaczy komplet, nie pustkę:
`nil` to „Operator o pamięci nie mówił”, a wtedy obowiązuje stan wyjściowy
platformy — pełny. Lista pusta `[]` jest czym innym: świadomym żądaniem
wyłączenia pamięci.

Imię własne i favikon podane przy zakładaniu eksperta idą tą samą drogą co
przy zmianie tożsamości. Ekspert założony bez nich nie jest ekspertem
niepełnym: `nazwa` niesie tożsamość, a imię własne jest tym, jak Operator
go woła.

Przy zmianie tożsamości imię własne i favikon idą osobną drogą, bo zapytanie
aktualizujące eksperta nie zna kolumn `imie_wlasne` i `favikon` —
repozytorium warstw ma na to własną czynność. Bez tego wywołania formularz
tożsamości przyjmowałby imię i favikon, rdzeń odpowiadałby powodzeniem,
a wartość przepadałaby.

Poziomy pamięci przy zmianie idą osobną drogą, bo siedzą w tabeli podrzędnej.
`nil` znaczy „pole pominięte” i zostawia zastane poziomy; lista pusta wyłącza
pamięć — dlatego wyłączenie nie potrzebuje własnej wartości wyliczenia.

Wskazany projekt zawęża wykaz do ekspertów w nim widocznych. Bez tego pole
`visibility` byłoby etykietą: dałoby się je zapisać i odczytać, ale nic by
z niego nie wynikało.

Warstwy i wtyczki dochodzą przy wykazie dwoma zapytaniami na cały wykaz, nie
dwoma na eksperta. Bez nich edytor warstw pokazywałby puste pola przy
ekspercie, który warstwy ma; pytanie o nie po jednym dałoby przy stu
ekspertach dwieście zapytań na jedno otwarcie biblioteki.
## budowa/server/internal/core/adapter_modul_tlumaczenie_polecenia.go

Własnego silnika tłumaczeń tu nie ma; przekład idzie kanałem modelu, więc
cała wiedza modułu o tym, jak ma wyglądać dobry przekład, mieści się w
treści polecenia. Plik adapter_modul_tlumaczenie_model.go odpowiada za
drogę do modelu, ten plik za treść polecenia.

Słownik Operatora wiąże przekład dwiema drogami. Przed wywołaniem terminy
wchodzą do polecenia jako wykaz obowiązkowych odpowiedników i wykaz nazw
nietykalnych: model, który zna słownik, użyje właściwego słowa w
odmienionej formie i we właściwym miejscu zdania, czego podmiana napisu
nie potrafi. Po wywołaniu wynik przechodzi tę samą mechaniczną podmianę,
co zastosowanie terminów słownika: wyłapuje termin zostawiony w brzmieniu
źródłowym, lecz nie zastępuje pierwszej drogi — trafia wyłącznie w formę
podstawową.

Zasady jakości także wchodzą do polecenia. Kontrakt zna sześć rodzajów
niezgodności: liczbę, datę, walutę, symbol zastępczy, długość i
pominięcie. Moduł szuka ich potem w wyniku. Rodzaj pominięcia jest tu
obecny, choć kontrola jakości go nie sprawdza, bo wymagałby rozumienia
treści — polecenie może żądać rzeczy, której kontrola potem nie
zweryfikuje.

Rozdział słownika na dwa wykazy jest istotny, nie kosmetyczny. Odpowiednik
dotyczy jednego języka docelowego. Nietykalność dotyczy terminu, nie
języka: nazwa własna produktu ma zostać nietknięta w każdym przekładzie,
więc wykaz nietykalnych zbiera terminy wszystkich języków słownika.
Zawężenie ich do języka panelu byłoby cichym rozluźnieniem zakazu
Operatora.

Kolejność treści polecenia nie jest przypadkowa — polecenia idą przed
tekstem, żeby model czytał je jako instrukcję, a nie jako część materiału
do przełożenia; sam tekst źródłowy zamyka polecenie pod wyraźnym
nagłówkiem z tego samego powodu. Żądanie oddania wyłącznie przekładu
zostaje, bo treść panelu ma być tłumaczeniem, nie rozmową o tłumaczeniu.

Gdyby słownik wszedł do tłumaczenia zwrotnego, przekład zwrotny
naprawiałby terminologię z powrotem na brzmienie źródłowe i Operator
zobaczyłby zgodność tam, gdzie jej nie ma — kontrola przestałaby
cokolwiek kontrolować. Z tego samego powodu polecenie żąda przekładu
dosłownego, nie gładkiego.

Termin bez odpowiednika nie niesie żadnego polecenia dla modelu: Operator
zaznaczył słowo, ale nie powiedział, czym je zastąpić.
## budowa/server/internal/core/adapter_modul_orkiestracja_zlozenie.go

Podzial wobec adapter_modul_orkiestracja.go idzie wzdluz odpowiedzialnosci:
tam mieszka powolanie podagentow (subagent.spawn) wraz z praca w tle, tutaj
wylacznie budowanie portu i jego wiazania. Wszystkie dokladki znosza sie
same przy pustej zaleznosci.

drogaNarzedzia: powolan bywa kilkanascie na ture. prace: bez tego wykazu
subagent.stop nie mialby czego zatrzymac i przepisywalby wylacznie wiersz
(adapter_modul_orkiestracja_zatrzymanie.go). rozgloszenie: tresc rozgloszen
niesie adapter_modul_orkiestracja_rozgloszenie.go.

zlozPodagentow stoi w tym pliku, a nie w montaz_porty.go, z dwoch powodow:
montaz przeklada byty na porty jednym wierszem na port, a wiazan podagentow
jest piec; dzieki temu montaz nie musi tez znac pakietu podagenci - wiedze
o ocenie zywotnosci trzyma adapter, ktory jako jedyny jej uzywa. Bez
silnika kolejek podagent zostaje pending, bez biegow powstaje poza biegiem,
bez nadzoru idzie bez zdania o procesie orkiestratora.

zWykonaniem: bez wpiecia podagent zostaje w stanie pending i to jest stan
prawdziwy, nie udawane wykonanie. subagent.result.collect oddaje tresc,
a nie puste pole, dzieki wpieciu adaptera jako ujscia wyniku silnika.
Wpiecie idzie stad, bo tylko ten adapter wie, ze pozycja miewa podagenta -
silnik zna wylacznie pozycje.

zZywotnoscia: wywolanie pozniejsze zamknieloby prace powolana przez ten
sam rdzen. Sierota to wiersz w stanie pending/running prowadzony przez
uruchomienie inne niz biezace - jego proces zginal razem z rdzeniem, wiec
stan running po restarcie klamie.

przyjmijWynikPozycji: brak trafienia to pozycja spoza podagentow (jeden
silnik); blad zapisu idzie do dziennika, bo cichej utraty wyniku nikt by
nie zobaczyl.

## budowa/server/internal/core/adapter_modul_terminal_plik.go

Bez osobnej komendy klient czytał manifest projektu POLECENIEM POWŁOKI, więc
wykrycie zadań zależało od tego, czy na maszynie stoi program wypisujący
plik, i od składni każdej z powłok osobno (`cat`, `Get-Content`, `type`).

Karta lokalna czyta plik przez `os` biblioteki standardowej. Programu powłoki
nie ma tu wcale: `cat` na karcie lokalnej byłby uruchomieniem procesu po to,
żeby zrobić to, co pakiet `os` robi jednym wywołaniem, i wprowadzałby
zależność od zawartości maszyny w czynność, która jej nie potrzebuje.

Karta ZDALNA czyta plik na swojej maszynie — tak mówi kontrakt — więc tam
program jest nieunikniony: plik leży po drugiej stronie łącza. Idzie tą samą
drogą co każdy proces rdzenia (`zewnetrzne.Wolaj` → `session.Uruchamiacz`),
z granicą czasu i objęciem drzewa potomstwa.

Ścieżka bezwzględna jest w kontrakcie dopuszczona, ale punkt izolacji „pliki”
obowiązuje ją tak samo jak względną: przy włączonym punkcie odczyt spoza
obszaru okna kończy się odmową `permission_denied`, a nie treścią.
## budowa/server/internal/core/adapter_modul_asystent.go

Stan zleceń (assistant.action.status) i dziennik (assistant.activity.list)
leżą w adapter_modul_asystent_czynnosci.go, a port Asystent i
zarejestrujAsystenta — w adapter_modul_asystent_uchwyty.go; wszystkie te
pliki piszą metody na jednym typie adapterAsystenta. Pakiet server/internal/mowa
wnosi silnik rozpoznawania mowy, a port Mowa (handlers_mowa.go) wystawia go
w kształcie kontraktu. Komenda assistant.voice.command ma pole Transcript
osobno od AudioRef: tekst poprawiony ręcznie ma pierwszeństwo przed
rozpoznaniem nagrania. Bez wpiętego silnika zlecenie powstaje ze śladem
nagrania, a transkrypcja zostaje pusta — odmowy nie ma, moduł pracuje
w zakresie, w którym może. Syntezy mowy to nie dotyczy: Speak i SpeechRef idą
w drugą stronę i rdzeń ich nie spełnia.

Przedrostki identyfikatorów bytów modułu są zadeklarowane w całości w tym
pliku, łącznie z przedrostkiem wpisu dziennika, którego plik nie nadaje sam,
żeby nie deklarować ich po raz drugi gdzie indziej.

Zależności wykonawcy zleceń są opcjonalne i wpinane w montażu metodami
budującymi, jak w module Roundtable: rejestr kanałów modelu, nadzorca sesji
(źródło kanału i katalogów okna) oraz nadajnik zdarzeń. Bez nich moduł
przyjmuje polecenia i prowadzi ich stan ręcznie, ale nie podejmuje zleceń
z kolejki — zlecenie zostaje queued, zamiast udawać wykonanie, tak jak debata
bez rejestru kanałów odmawia uruchomienia tury.

Pole mowa jest zależnością opcjonalną widzianą portem kontraktu, nie typem
pakietu mowa: moduł potrzebuje tu wyłącznie przekładu nagranie-tekst, a przez
port dostaje tę samą instancję silnika, którą obsługuje rodzina speech.*.
Pole mosty składa konfigurację MCP okna wraz z wpisem serwera narzędzi
modelu (most_narzedzi.go); bez tej zależności tura zlecenia idzie samym
tekstem, bez ani jednego narzędzia (adapter_modul_asystent_sterowanie.go).
Pole zycie jest kontekstem rdzenia, nie połączenia. Pole biegi wiąże kod
zlecenia z przerwaniem jego tury: bez tego wykazu sterowanie Operatora
(cancel, pause) zmieniałoby sam wiersz, a tura pracowałaby dalej i domykała
zlecenie mimo odwołania; zamki strzegą wykazu, bo tury biegną gorutynami,
a sterowanie przychodzi z pętli komend. Pole odwolane niesie decyzje
Operatora, które zapadły, zanim wykonawca zdążył przestawić zlecenie na
running — bez tego znacznika cancel zapisuje cancelled, a wykonawca, który
stan odczytał chwilę wcześniej, nadpisuje go z powrotem na running i dowozi
zlecenie do done; odwołanie ma być mocniejsze od tury, nie odwrotnie.

Trzy wyniki rozpoznania nagrania odpowiadają trzem zachowaniom, tym samym,
które rozdziela pole processed kontraktu speech.transcribe: rozpoznano tekst
— tekst wraca i staje się treścią polecenia; przetworzono, mowy brak —
odmowa nazywająca ten fakt, bo assistant.voice.command ma wytworzyć
polecenie, a puste polecenie nie jest poleceniem, zlecenie z pustą treścią
zostawiłoby w dzienniku modułu wpis, po którym nic się nie dzieje; nie
przetworzono — odmowa silnika wraca nietknięta, z jej własnym kodem
kontraktu i pełnym komunikatem (adapter_modul_mowa.go).

PolecenieGlosowe zakłada zlecenie i pierwszy wpis dziennika w jednej
transakcji — woła PrzyjmijPolecenie warstwy danych, nie dwa osobne zapisy,
żeby nigdy nie powstał wpis bez zlecenia ani zlecenie bez śladu w rozmowie.
Profil ma byt trwały: ProfileId z żądania odkłada się w kolumnie
zlecenie_asystenta.profil_kod, bo niesie warstwę promptu rozstrzygającą, czy
model sięgnie po narzędzia platformy, czy odpisze samym tekstem; wejście tej
warstwy do tury opisuje adapter_modul_asystent_profil.go. Kod profilu
nieznany kończy się nazwaną odmową jeszcze przed zapisem, bo klucz obcy
kolumny zamieniłby go w usterkę bez powodu. Bieg wykonawcy jest asynchroniczny
— komenda potwierdza przyjęcie od razu, a tura modelu trwa dłużej niż
wykonanie komendy; z wpiętym silnikiem mowy treść zlecenia jest złożona,
zanim zlecenie powstanie, a bez silnika zlecenie z samym nagraniem zostaje
w kolejce i czeka na tekst przysłany wprost.

## budowa/server/internal/core/adapter_modul_automations.go

Automations nie ma okna modułowego — jest komponentem własnym strony głównej
i nie pojawia się jako moduł w żadnym środowisku. Dlatego adapter nie zna ani
środowiska, ani karty sesji. Ten adapter buduje definicję i zapisuje
przebieg; wykonania nie prowadzi.

### adapterAutomatyk — pola

uklad niesie trzy dopełnienia układu zależności i spięcie kolejek
(`adapter_modul_orkiestracja_uklad.go`). Wpina je `ZUkladem`; nil znaczy
„nie wpięto", a wtedy cztery komendy odmawiają, a reszta modułu pracuje.

sejf jest magazynem WARTOŚCI poświadczeń, leżącym poza bazą. Wpina go
`ZSejfem`; nil znaczy „nie wpięto", a wtedy `automation.secret.set` i wymiana
klucza podpisu webhooka odmawiają wprost, zamiast zapisywać referencję
wskazującą na nic.

okna są rejestrem okien komunikacji. Służą wyłącznie przełożeniu roli
środowiska MultitaskingAI na okno przy spięciu kolejek.

### ZKolejkami, ZSejfem

`automation.queue.action` odmawia wprost bez wpiętego adaptera kolejek —
brak wykonawcy nie może udawać wykonania. Sejf: jedna instancja pod jednym
zamkiem, ten sam sejf plikowy co sekrety kont i punktów dostępu — dwa sejfy
nad tym samym plikiem ścigałyby się o zapis.

### Zapisz

Kroki podmieniają się w całości, gdy pole `steps` przyszło: Workflow Builder
oddaje po zmianie całą definicję, więc pole obecne znaczy „tak ma wyglądać
automatyka". Pole nieobecne zostawia kroki nietknięte, żeby zapis samej
nazwy albo samego przełącznika „czynna" nie skasował pracy Operatora.
Migawka wersji idzie PO złożeniu automatyki, bo zapisuje to, co naprawdę
stoi w bazie po zapisie — a nie to, co przyszło żądaniem. Kroki zastane
(żądanie bez pola `steps`) trafiają wtedy do wersji tak samo jak
podmienione, więc historia nie ma dziur po zapisie samej nazwy.

### odlozWersje

Nieudany zapis migawki nie wywraca zapisu definicji: definicja już stoi
w bazie, a odmowa komendy mówiłaby Operatorowi, że jego praca przepadła,
choć nie przepadła.

### Automatyka — etykiety i wersja

Etykiety idą razem z definicją, bo wykaz Workflow Buildera filtruje po nich
bez drugiej komendy. Nieudany odczyt etykiet daje wykaz pusty zamiast
wywracać odczyt automatyki — brak etykiet jest stanem poprawnym.

Wersja oddawana kontraktem jest wersją WYKONYWANĄ: opublikowana, gdy
Operator rozdzielił roboczą od opublikowanej, a bieżąca, gdy rozdziału nie
wprowadził. Inaczej okno pokazywałoby numer wersji roboczej przy
automatyce, która produkcyjnie wykonuje wersję wcześniejszą.

### bladAutomatyki

Błąd, któremu kod już nadano — odmowa wskazania, brak bytu, brak wykonawcy —
przechodzi tędy bez zmiany kodu; dopiero usterka bez kodu staje się usterką
wewnętrzną rdzenia.

## budowa/server/internal/core/adapter_modul_tlumaczenie_slownik.go

Obszar dotyczy typu `*adapterTlumaczenia` zadeklarowanego w
`core/adapter_modul_tlumaczenie.go` — ten plik nie deklaruje ani typu
adaptera, ani konstruktora, ani przedrostków identyfikatorów. Flaga
`NieTlumaczyc` jest zdefiniowana w `migracja_054_slownik_tlumaczenia.sql`.

Słownik nie żyje wyłącznie tutaj. `glossary.apply` jest narzędziem
naprawczym: ujednolica terminologię treści, która już powstała. Właściwym
miejscem słownika jest sam przekład — terminy i zakazy Operatora wchodzą do
polecenia dla modelu przy `target.add` (`adapter_modul_tlumaczenie_polecenia.go`),
a `zastosujTerminySlownika` z tego pliku przechodzi jeszcze po wyniku modelu
jako siatka bezpieczeństwa. Ta sama funkcja w dwóch zastosowaniach, nie dwie
kopie zasady.

Zakres obu komend bez wskazania panelu bierze się z `WszystkiePanele`
repozytorium: `glossary.apply` z pustym `PanelId` przechodzi po komplecie
paneli, a `glossary.occurrences` — które w kontrakcie niesie samo `Term` —
przeszukuje ten sam komplet.

Zapis `UstawTerminy` idzie przez `ZapiszTerminy` (jedna transakcja, choć
niesiona jest tu zawsze lista jednoelementowa) — warstwa danych i tak nie ma
osobnej drogi INSERT/UPDATE (`dane/slownik.go`), więc adapter nie dubluje
tego rozróżnienia.

Panel bez treści w bazie w `zastosujWPanelu` daje zero — treść żyjąca
wyłącznie poza bazą (`TrescOdwolanie`) nie jest czytana, bo adapter nie sięga
po pliki spoza repozytorium.

`zastosujTerminySlownika` oddaje zmienioną treść i liczbę faktycznie
podmienionych wystąpień, żeby `ChangedCount` mówił prawdę o skutku, nie
o liczbie terminów w słowniku.

Zakresem `WystapieniaTerminow` są wszystkie panele: kontrakt nie niesie
wskazania okna ani panelu, a słownik jest jeden na instalację, więc jedynym
uczciwym odczytaniem jest przeszukanie kompletu paneli. Otoczenie wystąpienia
to wycinek treści wokół trafienia — kontrakt chce „wystąpień wraz
z otoczeniem”, a nie samych pozycji.
## budowa/server/internal/core/adapter_modul_terminal.go

Powloki lezy w adapter_modul_terminal_powloki.go, tryb uprawnien w
adapter_modul_terminal_uprawnienia.go, ewidencja procesow w
adapter_modul_terminal_rejestr.go, strumien wyjscia w
adapter_modul_terminal_strumien.go, a uruchomienie procesu
w adapter_modul_terminal_wykonanie.go.

Przed kazdym procesem stoja trzy sprawdzenia, w tej kolejnosci: tryb
uprawnien okna - czy wolno w ogole uruchomic proces i czyim poleceniem
(PermissionMode, adapter_modul_terminal_uprawnienia.go); egzekutor izolacji
- czy polecenie miesci sie w obszarze okna (session.SprawdzPolecenie nad
zasadami z izolacja.go); uruchamiacz warstwy kanalu - jedyna droga startu
procesu w drzewie (port session.Uruchamiacz). Rdzen nie buduje wlasnego
exec.Cmd ani drugiego egzekutora, korzysta wylacznie z tych dwoch warstw.

przedrostekProcesuTerminala znakuje identyfikator procesu terminala,
odrebny od proc- telemetrii: tamten opisuje bieg tury modelu, ten proces
urzadzenia, i mylenie ich w Process Monitorze byloby kosztowne.

uruchamiacz jest portem warstwy kanalu - jedyna droga startu procesu.
rozstrzygacz i katalog skladaja zasady izolacji obowiazujace w oknie.
wyjscie rozsyla fragmenty strumienia do okna Output Console. zmiana
rozglasza terminal.process.changed, podpina ja obslugiwacz.

OtworzKarte: karta bez okna nie ma ani trybu uprawnien, ani obszaru
izolacji, a wiec nie da sie jej pozniej wykonac.
## budowa/server/internal/core/adapter_modul_poczta.go

Aplikacja nie ma własnego serwera poczty — używa skrzynki, którą Operator już
ma skonfigurowaną na urządzeniu albo w chmurze. Rdzeń jest więc klientem
cudzej skrzynki i nikim więcej: nie stawia serwera, nie zakłada kont
pocztowych, nie pośredniczy przez żadną infrastrukturę Danaco i nie trzyma
cudzej poczty u siebie. Podpina skrzynkę, którą Operator już ma, i rozmawia
z nią jej protokołem.

Hasło albo token przychodzi polem `secret` żądania `mail.account.add`
i natychmiast ląduje w sejfie poświadczeń — tym samym, którym jadą konta
i punkty dostępu. Baza dostaje odwołanie ("sejf:poczta:<kod>"), nigdy sekret;
kolumny na sekret nie ma w schemacie w ogóle. Do dziennika sekret nie trafia,
bo ten moduł nie loguje niczego. Do odpowiedzi komendy nie trafia, bo
`MailAccount` w kontrakcie nie ma pola, w które dałoby się go włożyć.

Brak skrzynki, brak poświadczenia i brak łączności to trzy różne rzeczy,
z których każdą Operator naprawia inaczej: pierwszą podpięciem skrzynki,
drugą podaniem hasła, trzecią zajrzeniem do sieci albo do dostawcy. Jedno
wspólne "poczta niedostępna" kazałoby mu zgadywać, którą.

Magazyn załączników znajduje repozytorium Designu, to samo, do którego pisze
`design.asset.upload`. Drugiego magazynu zasobów rdzeń nie ma i mieć nie
będzie: załącznik listu wciągnięty osobną drogą byłby zasobem, którego
narzędzia obrazu i dokumentów nie widzą.

Kod `channel_unavailable` przy niepowodzeniu połączenia, nie `internal_error`:
skrzynka jest po drugiej stronie sieci i jej niedostępność bywa chwilowa,
więc odmowa jest ponawialna (`shared.KodyPonawialne`). Rdzeń nie zawinił
i nic tu nie naprawi — a klient, który ponowi za minutę, ma szansę trafić.

Sekret dobrany dla nastaw połączenia żyje tylko do końca komendy. Brak
sekretu nie jest tu błędem: `poczta.Polacz` nazwie go drugą z trzech odmów —
brakiem poświadczenia, odróżnionym od braku skrzynki i braku łączności.
## budowa/server/internal/core/adapter_modul_asystent_profil.go

Profil niesie warstwę promptu, która mówi modelowi, że jest klawiaturą
Operatora, a nie autorem. To ona rozstrzyga, czy model sięgnie po
narzędzia platformy, czy odpisze tekstem we własnym oknie. Bez niej
asystent zachowuje się jak zwykły czat, dlatego warstwa ma byt trwały, a
nie żyje w polu żądania, które ginie razem z odpowiedzią na komendę.

Warstwa profilu ląduje w nakładce zapytania, dokładnie tam, gdzie rozmowa
kładzie warstwy osi tożsamości i warstwy eksperta. Druga droga do promptu
znaczyłaby dwie prawdy o tym, co model naprawdę dostał.

Profil dopisuje, nigdy nie zastępuje — tak samo jak ekspert. Prompt
wbudowany programu kanału zostaje w mocy, a profil dokłada do niego
paczkę treści. Profil nie ma żadnej drogi, którą mógłby prompt platformy
zdjąć.

Syntezy mowy tu nie ma. Kolumna głosu syntezy niesie nastawę głosu
odczytu, ale rdzeń syntezy nie wykonuje — nastawa jest odkładana i
czytana, a nie udawana wywołaniem, którego nikt nie spełnia.

Kolejność rozstrzygania profilu zlecenia jest ustaleniem, nie wygodą:
profil wskazany przez żądanie ma pierwszeństwo, bo Operator powiedział
wprost, którym profilem pracuje; profil domyślny jest kolejny, bo okno
asystenta ma ruszać bez wybierania profilu przy każdym poleceniu; brak
obu nie jest usterką.

Brak profilu nie wstrzymuje pracy. Instalacja świeża nie ma ani jednego
wiersza profilu, bo komenda zakładająca profil nie ma na to kontraktu.
Gdyby brak warstwy odmawiał tury, asystent nie ruszyłby na takiej
instalacji w ogóle. Zlecenie idzie więc bez warstwy, a fakt zostaje
nazwany wpisem dziennika, zamiast zniknąć.

Wskazanie na profil, którego nie ma, kończy się tak samo jak brak
wskazania: bez warstwy i z nazwanym powodem. Zerwanie zlecenia byłoby tu
karą za nastawę okna, na którą Operator w tej komendzie nie ma wpływu.

Warstwa promptu profilu idzie do warstwy profilu nakładki, a nie do
konstytucji ani do ekspertyzy. Zdanie opisujące rolę klawiatury Operatora
opisuje rolę, w której model ma wystąpić; konstytucja jest warstwą
platformy, nie profilu jednego okna, a ekspertyza jest wiedzą zadaniową.
Włożenie tego zdania w konstytucję postawiłoby nastawę jednego okna wyżej
niż zasady całej platformy.

Dokładanie treści warstwy, a nie podstawianie jej, znaczy, że dołożenie
osi tożsamości w przyszłości nie skasuje warstwy profilu po cichu.

Nastawy poza promptem — kanał, zasięg urządzenia i zasięg pracy — są
nastawą zasięgu, nie bramką: profil mówi, dokąd sięga praca, a nie czego
zabrania. Nakładane są tylko wtedy, gdy profil naprawdę coś w danej
sprawie mówi: pole puste znaczy brak zdania profilu i wtedy zostaje
nastawa okna. Bez tego strażnika profil bez wskazanego kanału zabrałby
turze jedyny kanał, jaki miała.

Sprawdzenie profilu wskazanego przez żądanie jest konieczne: kolumna ma
klucz obcy do tabeli profili, a więzy obce są na połączeniu włączone —
kod nieznany wywróciłby cały zapis polecenia usterką więzów, czyli
błędem wewnętrznym bez powodu czytelnego dla Operatora. Sprawdzenie
wcześniejsze pozwala nazwać powód.

To odmowa, a nie ciche pominięcie: brak profilu w ogóle pracy nie
wstrzymuje, ale wskazanie profilu, którego nie ma, jest czym innym —
Operator powiedział wprost, którą warstwą promptu ma pracować model.
Ciche pominięcie puściłoby turę z warstwą profilu domyślnego albo bez
żadnej, meldując wykonanie polecenia, którego nikt nie wydał.

Kontrakt nie ma pola na powód braku warstwy profilu: pole wyniku akcji
jest miejscem na odpowiedź modelu — dopisanie tam zdania rdzenia
zmieszałoby dwa głosy w jednym polu i zafałszowało wynik tury. Dziennik
aktywności jest tą samą rozmową, w której stoi polecenie i wynik, więc
Operator czyta powód dokładnie tam, gdzie patrzy — i zostaje mu ślad po
zleceniu, a nie sam komunikat, który znika.

## budowa/server/internal/core/adapter_modul_workspace_projekt.go

Oś czasu ma pokazać, co się w projekcie działo, więc zapis zdarzenia idzie
w tej samej czynności, która zmienia byt — a nie z wyliczenia stanu przy
odczycie. Wyliczenie pokazywałoby wyłącznie to, co jeszcze istnieje.

Historia po przywróceniu instrukcji zostaje nietknięta: „przywróć” jest
czynnością odnotowaną, a nie cofnięciem czasu.

Niepowodzenie zapisu zdarzenia osi czasu nie przerywa czynności, która je
wywołała: zadanie już powstało, a brak wiersza w dzienniku jest ubytkiem
zapisu, nie unieważnieniem skutku. Dlatego funkcja `odnotujZdarzenieWorkspace`
nie oddaje błędu.
## budowa/server/internal/core/adapter_modul_poczta_wiadomosci.go

Załącznik wchodzi do magazynu tą samą drogą, co wniesienie zasobu w
module Design: bajty lądują w magazynie rdzenia pod sumą sha256, a obok
powstaje wiersz zasobu w tabeli zasobów Designu — taki sam, jaki zakłada
wniesienie pliku do Assets Panelu. Dzięki temu model dostaje w odpowiedzi
identyfikator zasobu, którym woła narzędzia treści. Rdzeń nie ma drugiego
magazynu bajtów, więc załącznik odłożony osobno byłby plikiem, którego
żadne narzędzie nie widzi.

Kolumna okna zasobu Designu niesie tu odwołanie do skrzynki, nie do okna.
Jest wymagana, bo zasób Designu należy do okna modułu, a załącznik listu
do żadnego okna nie należy — należy do skrzynki. Wpisywana jest więc
wartość złożona z przedrostka poczty i kodu skrzynki: to prawda o
pochodzeniu zasobu, daje się odfiltrować i nie miesza się z zasobami
okien. Zmyślenie identyfikatora istniejącego okna byłoby wstawieniem
cudzych plików do cudzego panelu.

Wciągnięcie załączników jest warunkowe, bo kosztuje: list z wielomegabajtowym
skanem odczytany po to, żeby sprawdzić datę, nie ma powodu zostawiać po sobie
tych megabajtów w katalogu danych.

Wciąganie załączników odkłada bajty w magazynie i zakłada wiersze zasobów w
kolejności zamierzonej: najpierw bajty, potem wiersz, tak jak przy wniesieniu
zasobu w module Design. Wiersz wskazujący blob, którego nie ma, byłby
zasobem, po który model sięgnie i niczego nie znajdzie.

Bład odczytu listu odróżnia brak wiadomości od braku skrzynki i od braku
łączności. List przeniesiony albo skasowany w kliencie poczty Operatora
między wykazem a odczytem jest normalnym stanem cudzej skrzynki, więc kod
jest brakiem znalezienia, a nie awarią rdzenia.

## budowa/server/internal/core/adapter_modul_terminal_powloki_urzadzen.go

Powłoki z `adapter_modul_terminal_powloki.go` są jednym wierszem wykazu:
program i argumenty poprzedzające treść. Odmowa idzie przy SKŁADANIU
polecenia, a nie przy otwarciu karty. Karta jest profilem powłoki i wolno ją
otworzyć z pustym wskazaniem, tak samo jak kartę zdalną bez adresu; dopiero
polecenie musi wiedzieć, gdzie się wykonać.

Powłoki te nie zarządzają kontenerem, podem ani urządzeniem: nie zakładają,
nie usuwają i nie zmieniają ich stanu. Wykonują polecenie w bycie, który już
istnieje. Kontenerami zarządza moduł Developer własną drogą (biblioteka
Dockera), a nie karta terminala.

W `polecenieUrzadzenia`: `session.Polecenie` niesie ścieżkę, a program
dołożony do pakietu produktu leży poza ścieżką wyszukiwania systemu i po
samej nazwie by nie wystartował. Brak programu kończy się odmową nazywającą
go wraz z pakietem — tą samą, którą oddaje cały arsenał
(`zewnetrzne.BrakNarzedzia`).
## budowa/server/internal/core/adapter_modul_badania_raport.go

Sekcje redagowane wprost (z.Sections) przechodzą bez zmiany — to redakcja
Operatora, bez modelu. FindingIds bez sekcji własnych uruchamia redakcję
modelem: most zapytajModel streszcza zebrane ustalenia w jedną sekcję
nadrzędną (Streszczenie ustaleń), po której idą sekcje szczegółowe, po jednej
na ustalenie, z treścią ustalenia skopiowaną dosłownie dla śladu — tytuł to
kod ustalenia, bo kontrakt nie daje ustaleniu własnego tytułu. Streszczenie
bierze domyślny czynny kanał rejestru (domyslnyKanalBadania), a brak czynnego
kanału to odmowa wprost.

WyeksportujRaport zapisuje ślad eksportu w bycie trwałym
eksport_raportu_badania, bo LibraryFileId/Path/SizeBytes muszą przeżyć
wywołanie. Pliku na dysk rdzeń nie zapisuje — nie ma magazynu blobów, ten sam
brak co w module Library. Odpowiedź zostawia LibraryFileId, Path i SizeBytes
puste, zamiast udawać wytworzony plik.

UstawPrzestrzen operuje na przestrzeni jednowierszowej — kontrakt
research.workspace.set nie niesie identyfikatora okna ani sesji (sprawdzone
w dane/badania_raport.go), więc jedna przestrzeń na instalację jest jedynym
uczciwym odczytaniem.

sekcjaStreszczenia wymaga modelu: składa z treści ustaleń polecenie redakcji
i zwraca sekcję z odpowiedzią modelu; błąd kanału to odmowa wprost, nie
streszczenie zmyślone bez modelu.

## budowa/server/internal/core/adapter_modul_aplikacje_wdrozenie_bieg.go

Silnik pracuje poza żądaniem. Komenda `apps.deployment.run` kończy się, gdy
przebieg ruszy, a nie gdy się skończy — rozłączenie klienta w połowie nie
przerywa wdrożenia. Dlatego bieg idzie własną gorutyną i własnym kontekstem
(`context.Background`), a nie kontekstem komendy.

Stan końcowy wynika z wykonanej pracy. Krokiem wdrożenia jest sprawdzenie,
czy jest co wdrożyć: wdrożenie w przód udaje się, gdy przestrzeń robocza
okna niesie choć jeden plik; pusta przestrzeń kończy się `failed`;
cofnięcie udaje się, gdy wdrożenie docelowe kiedykolwiek weszło w
`succeeded`; cofnięcie do przebiegu, który nigdy się nie powiódł, kończy
się `failed`.

Rdzeń nie hostuje produktu, więc wdrożenie w przód zostawia pole `Url`
puste — postawienie serwera produktu wymagałoby uruchamiacza procesu
wpiętego w adapter. Cofnięcie dziedziczy `Url` wprost z wdrożenia
docelowego.

Przebieg zostawia po sobie artefakt i dziennik. Udane wdrożenie w przód
pakuje przestrzeń roboczą okna w archiwum `zip` i kładzie je w magazynie
treści rdzenia, a wiersz `artefakt_apps` wskazuje ten plik wraz z
rozmiarem i sumą kontrolną — to on jest wejściem `apps.package.build`
i pozycją `apps.artifact.list`. Każdy krok przebiegu dopisuje wiersz do
`wiersz_dziennika_apps`, skąd czyta go `apps.deployment.log.read`. Bez tych
dwóch rzeczy trzy komendy rodziny meldowałyby pustkę przy przebiegu, który
naprawdę się odbył.

### wykonajWdrozenie

Kontekst jest własny, nie komendy — wdrożenie przeżywa rozłączenie
klienta, tak jak opisano dla całego silnika.

### zapiszIRozglosWdrozenie — powód

Żądanie cofnięcia do przebiegu bez werdyktu odrzuca już obsługiwacz komendy
(`adapter_modul_aplikacje_wdrozenie.go`); tutaj zostaje przypadek, w którym
cel zmienił stan między odczytem komendy a odczytem silnika.
## budowa/server/internal/core/adapter_modul_tlumaczenie_profile_qa.go

Zatwierdzenie zostawia dwa slady naraz i to jest zamierzone: wiersz obiegu
(kto, kiedy, na jakim etapie, z jaka uwaga) oraz migawke na panelu, ktora
widzi kazdy odczyt panelu bez dociagania historii. Oba zapisy ida jedna
transakcja (dane/tlumaczenie_kontrola.go).

Autor zatwierdzenia bierze sie z sesji wywolujacego, nie z zadania. Gdyby
przychodzil zadaniem, obieg zatwierdzen bylby polem tekstowym, w ktore
mozna wpisac dowolne nazwisko - a wtedy nie jest obiegiem, tylko notatka.

## budowa/server/internal/core/adapter_modul_library_podglad.go

Metoda `Podglad` (`adapter_modul_library.go`) rozstrzyga rodzaj podglądu
i składa odpowiedź kontraktu; tutaj leży sam odczyt treści spod odwołania
pliku, wydzielony, bo to osobna odpowiedzialność — dostęp do nośnika.

Odwołanie jest ścieżką na dysku i tylko tyle ten plik o nim zakłada. Wgranie
przez `sourcePath` czyni ścieżkę odwołaniem do treści; wgranie przez
`contentBase64` odkłada bajty w magazynie treści rdzenia i odwołaniem czyni
ścieżkę bloba (`Wgraj`, `adapter_modul_library_magazyn.go`). Obie drogi
kończą się ścieżką, więc czytelnik jest jeden i czyta ją z dysku tak, jak
moduł Developer czyta pliki repozytorium do Code Editor
(`adapter_modul_developer_plik.go`). Plik bez odwołania — wiersz sprzed
wpięcia magazynu albo wersja będąca samym znacznikiem — jest odmawiany
wcześniej (`bladBrakuTresciBiblioteki`), więc tu trafia wyłącznie plik
z realnym odwołaniem.

Granica podglądu na znaki ma pokazać początek dokumentu w oknie modułu,
a nie przesłać cały plik przez gniazdo zdarzeń. Próg rozpoznania tekstu to
głowa dokumentu — dość, by trafić na bajt zerowy albo na niepoprawny UTF-8
formatu binarnego, i mało, by nie czytać całości.

`trescPodgladuBiblioteki`: odwołanie to ścieżka pliku na dysku (`Wgraj`,
`sourcePath`), więc rdzeń czyta ją tak, jak moduł Developer czyta pliki
repozytorium. Czyta o jeden bajt więcej niż górna granica w bajtach — nadmiar
mówi, że treść jest dłuższa niż podgląd, więc podgląd jest skrócony.
Odwołanie, którego nie da się odczytać, jest odmową wprost, a nie pustą
treścią udającą podgląd.

`odwolanieObrazuPodgladu`: pierwsze zawężenie — tylko obraz. Kontrakt nazywa
to pole odwołaniem do obrazu, więc wypełnia je wyłącznie podgląd obrazowy.
Przy `binary` czy `pdf` klient dostałby wskazanie, którego nie umie
otworzyć, a przy okazji wyszedłby na zewnątrz układ katalogu danych
Operatora. Drugie zawężenie — tylko postać względna. Wychodzi ścieżka
względna magazynu, ta sama postać, co `asset.uri` modułu Design — ścieżka
bezwzględna nazywałaby katalog danych rdzenia, a odbiorca bywa na innej
maszynie i tak by jej nie otworzył. Pustka znaczy „nie mam czego podać”,
a nie „obraz bez treści”.

`rodzajPodgladuTresci`: plik tekstowy bez `mimeType` uznany z góry za
`binary` nie pokazałby ani jednego znaku treści, choć rdzeń ma i bajty,
i czytnik tekstu. Bajt zerowy albo ciąg niebędący poprawnym UTF-8 znaczy
treść nietekstową i wtedy `binary` jest prawdą. Sekwencja ma najwyżej cztery
bajty, więc urwany ogon to najwyżej trzy: dalsze skracanie zjadałoby treść
binarną aż do pustej, a pusta przechodzi jako poprawny UTF-8 i cały plik
wyszedłby tekstem.

`rodzajPodgladu` stoi tutaj, przy odczycie treści podglądu, a nie przy
składaniu odpowiedzi (`adapter_modul_library.go`), bo to jedna decyzja tego
samego obszaru.
