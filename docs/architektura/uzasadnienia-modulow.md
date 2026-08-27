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
