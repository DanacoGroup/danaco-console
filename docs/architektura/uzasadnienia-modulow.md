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
