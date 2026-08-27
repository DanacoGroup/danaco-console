*Uzasadnienia komentarzy i decyzji zaplecza, przeniesione z komentarzy kodu przy redakcji; rozdziały nazwane ścieżkami plików.*

# Uzasadnienia — zaplecze

## budowa/server/internal/store/migracja_012_katalog_ustawien.sql

Oś zasięgu `ConfigAxis` (`platform` · `model` · `account`) jest prostopadła do
ośmiu poziomów zasięgu: poziom mówi jak wąsko (okno → … → globalny), oś mówi
dla czego (platforma, model, konto). Klucz rozstrzygania jest złożony: klucz +
poziom + byt poziomu + oś + byt osi, dlatego `ustawienie` niesie kolumny `os`
i `klucz_osi`, a więz jednoznaczności obejmuje obie.

Kolumna `os_zasiegu.pierwszenstwo` porządkuje osie od najszerszej (platforma)
do najwęższej (konto): konto jest bytem konkretnym, model klasą, platforma
tłem. Rozstrzyganie idzie najpierw po poziomie, a dopiero w ramach poziomu po
osi: konto → model → platforma. Ustawienie per konto zapisane globalnie nie
bije więc ustawienia zapisanego na oknie — oś opisuje adresata wartości, nie
jej wagę.

Tabela `ustawienie` nie jest wskazywana kluczem obcym z żadnej innej tabeli,
więc przebudowa (nowa tabela → przepisanie wierszy → podmiana nazwy) nie rusza
niczyich odwołań. Indeks `idx_ustawienie_klucz` wraca pod tą samą nazwą.

Pozycje katalogu ustawień odpowiadają bytom czytanym przez rdzeń:

1. `server/internal/konfig/definicje_wykonania.go` — osiem parametrów
   wykonania okna komunikacji wraz z wartościami domyślnymi pochodzącymi ze
   stałych kontraktu (PermissionMode.manual, WindowRole.standalone,
   ExecutionEnv.local).
2. `server/internal/konfig/definicje_izolacji.go` — jedenaście punktów
   izolacji: trzy wymiary kontekstu (`odrebna`) i osiem zakresów technicznych
   (`wylaczony`).
3. Klucze katalogu roboczego i tożsamości: `katalog.roboczy.podstawa`,
   `katalog.roboczy.wzorzec_sesji`, `tozsamosc.tryb_domyslny`.
4. `server/internal/injection/ustawienia.go` — pola struktury Ustawienia:
   Program (`harness.program_claude`), PlikUstawien (`harness.plik_ustawien`,
   przełącznik --settings), KonfiguracjaMCP (`harness.konfiguracja_mcp`,
   przełącznik --mcp-config).
5. `server/internal/session/obieg.go` — ProgBrakuPostepuDomyslny
   (`petla.prog_braku_postepu`, pole LoopState.threshold kontraktu).
6. Egzekwowanie uwierzytelniania włączane przez Operatora
   (`bezpieczenstwo.egzekwowanie_uwierzytelniania`).
7. `client/src/motyw/motyw.ts` — wybór motywu `light` albo `dark`; brak
   wyboru znaczy preferencję systemu (`personalizacja.motyw`).

Poza katalogiem pozostają parametry startu procesu — port, katalog danych,
katalog klienta, katalog profili, rola procesu. Czyta je
`server/internal/konfiguracja` w chwili startu z warstw: wartość domyślna →
zmienna środowiska → argument wywołania (`.env.example`), a plik bazy jest
dopiero skutkiem tych parametrów. Wiersz katalogu byłby dla nich drugim
źródłem prawdy, którego nikt nie czyta.

Każde wstawienie kończy się klauzulą ON CONFLICT DO NOTHING, więc migracja
przechodzi także na bazie, w której część wierszy już jest. Klauzula
`WHERE true` przed ON CONFLICT jest wymogiem składni SQLite dla
INSERT ... SELECT z upsertem.

`wartosc_domyslna` trzymana jest tak samo jak `ustawienie.wartosc`: napisem
w postaci właściwej dla rodzaju. Brak wiersza w tabeli `ustawienie` znaczy
właśnie tę wartość, a brak wiersza w katalogu nie jest awarią — rezolwer
schodzi wtedy na rejestr wbudowany rdzenia.

## budowa/server/internal/store/migracja_016_stan_kolejki_wyczerpana.sql

Kontrakt zna pięć stanów kolejki: QueueStatus.idle, running, paused, stopped,
done. Mapa przekładu `WartosciBazyQueueStatus` wiąże je kolejno z kolumnami
'bezczynna', 'pracuje', 'wstrzymana', 'zatrzymana', 'wyczerpana'. Więz CHECK
wprowadzony w `migracja_003_kolejki.sql` dopuszczał tylko cztery pierwsze, więc
zapis stanu 'wyczerpana' kończył się odmową schematu, mimo że kontrakt ten
stan zna. Kolejka, w której nie została ani jedna pozycja czynna, musi mieć
jak się nazwać, inaczej silnik wykonania nie ma stanu końcowego dla kolejki
i po ostatniej pozycji błędnie zgłasza dalszą pracę.

SQLite nie udostępnia polecenia zdejmującego więz CHECK z istniejącej kolumny,
więc tabela `kolejka` przechodzi przebudowę: nowa tabela z poprawionym więzem,
przepisanie wierszy, podmiana nazwy — tak samo jak `ustawienie` w
`migracja_012_katalog_ustawien.sql`. Odmienność wynika z tego, że `kolejka`
jest wskazywana kluczem obcym przez `pozycja_kolejki` i przez
`log_akcji_kolejki`, w obu przypadkach z ON DELETE CASCADE. Przy włączonym
więzie kluczy obcych (`store/baza.go`, pragma `foreign_keys(1)` w DSN
każdego połączenia) polecenie DROP TABLE na tabeli rodzica wykonuje niejawne
DELETE wszystkich jej wierszy, co uruchamia kaskadę i kasuje wszystkie
pozycje kolejek oraz cały dziennik akcji. Pragmy nie da się wyłączyć na czas
kroku, ponieważ migracja biegnie w transakcji (`store/migracje.go`), a
`PRAGMA foreign_keys` w transakcji jest bez skutku. Z tego powodu przebudowa
obejmuje wszystkie trzy tabele obszaru w porządku bezpiecznym dla danych:
kopie z poprawionym więzem i przepisaniem wierszy co do kolumny; skasowanie
tabel starych od dziecka do rodzica, tak że w chwili kasowania `kolejka`
żaden wiersz jej już nie wskazuje i kaskada nie ma czego zabrać; podmiana
nazw poleceniem ALTER TABLE ... RENAME, które przepisuje odwołania kluczy
obcych w tabelach pozostałych; odtworzenie indeksów pod nazwami z
`migracja_003_kolejki.sql`, wolnymi po skasowaniu tabeli, do której należały.

Słownik stanów pozycji kolejki pozostaje bez zmiany: siedem wartości z
`migracja_003_kolejki.sql` wystarcza silnikowi wykonania na pełny cykl życia
zlecenia. Migracja naprawia więz stanu kolejki, nie rozszerza model danych.

## budowa/server/internal/store/migracja_192_roundtable_graf.sql

Graf argumentów jest trwały, a nie liczony przy każdym odczycie, ponieważ węzeł
da się oznaczyć jako kluczowy (`roundtable.argument.pin`). Oznaczenie postawione
na węźle wyliczanym w locie znikałoby przy następnym odczycie, bo węzeł
dostawałby wtedy nowy identyfikator.

Katalog błędów logicznych jest dwuwarstwowy. Definicje wnosi migracja i są
wspólne dla całej platformy — nie należą do okna. Zakres wykrywania
(`roundtable.fallacy.catalog.set`) należy do okna, więc leży w osobnej tabeli
wiążącej kod błędu z oknem; bez tego rozdziału wyłączenie błędu w jednym oknie
wyłączałoby go wszystkim.

## budowa/server/internal/store/migracja_051_aplikacje.sql

Warsztat modułu Apps trzyma treść pliku warstwy frontendu i backendu produktu,
a nie stan projektu. Kontrakt `apps.workspace.update` niesie `Layer`
(frontend/backend), `Path`, `Content` i opcjonalny `ComponentId` — kształt
zgodny z zapisem pliku edytora w module Developer (`developer_wersja_pliku`),
odmienny od tabeli `projekt`, która trzyma nazwę, opis i stan, bez treści
pliku. Współdzielenie jednej tabeli `projekt` obciążałoby każdy jej wiersz
kolumnami treści pliku, których większość wierszy nigdy nie użyje, dlatego
warsztat Apps ma własną tabelę budowaną na wzór `developer_wersja_pliku`.
Commity kodu produktu idą wspólną komendą `developer.git.action`, poza
zasięgiem tej migracji.

Plik warsztatu trzyma jeden wiersz na parę (okno, warstwa, ścieżka), nie
historię wersji: kontrakt `AppsWorkspaceUpdateRequest` nie niesie odpowiednika
`createVersion` z modułu Developer, więc `apps.workspace.update` jest zwykłym
nadpisaniem stanu bieżącego, a tabela odzwierciedla to wprost przez UPSERT po
kluczu (okno, warstwa, ścieżka), bez osobnej tabeli-dziennika.

Architektura trzyma komponenty jako tabelę własną, nie jako zapis w kolumnie
JSON, ponieważ `AppComponent` niesie pola, po których trzeba filtrować
i wiązać przy odczycie: `Kind` wchodzi do walidacji układu, a `DependsOn`
tworzy graf zależności; parsowanie JSON przy każdym odczycie byłoby kosztem
walidacji układu (`AppArchitecture.ValidationIssues`) i zapytań o zależność.
`apps.architecture.define` nadsyła całą listę komponentów na nowo (kontrakt:
`Components []AppComponent`, bez trybu częściowej zmiany), więc zapis jest
zawsze usunięciem komponentów architektury i wstawieniem przysłanych od nowa.
Zależność między komponentami ma z tego samego powodu własną tabelę
złącznikową: `AppComponent.DependsOn []string` jest listą identyfikatorów
zewnętrznych innych komponentów tej samej architektury.

Dziennik wdrożeń powtarza wzór `developer_budowanie`, ponieważ
`apps.deployment.run` i `developer.build.run` mają identyczny kształt zadania:
zlecenie z zewnętrznym kodem, oknem, czasem startu i końca, stanem i śladem
tekstowym. Rdzeń niczego nie wdraża naprawdę — tabela niesie wyłącznie ślad
zlecenia (środowisko, strategia, wersja, notatki, adres po wdrożeniu,
odnośnik do logu) i jego stan, a sam przebieg prowadzi rdzeń w pamięci.
`RollbackToDeploymentId` i `RolledBackFromId` są parą pól tego samego łuku,
więc jedna kolumna samoodwołania wystarcza, bo odczyt idzie zawsze od strony
cofnięcia. Log wdrożenia jest odwołaniem, nie treścią w bazie: kontrakt niesie
`LogRef *string`, więc kolumna `log_odwolanie` przechowuje ten odnośnik wprost.

## budowa/server/internal/store/migracja_013_punkty_dostepu.sql

Punkt dostępu nie jest środowiskiem: kolumna środowisko pozostaje profilem
widoczności modułów w bocznej nawigacji, więc dopisanie maszyny do tamtej
tabeli rozbiłoby nawigację platformy, a niniejsza migracja tamtej tabeli nie
dotyka. Punkt dostępu nie jest też katalogiem roboczym modelu: katalog
roboczy to ustawienie kluczy katalog.roboczy.podstawa i
katalog.roboczy.wzorzec_sesji, mówiące, gdzie model zostawia własne pliki,
podczas gdy punkt dostępu mówi, do czego model ma wgląd — model może czytać
jeden katalog, a pliki zostawiać w zupełnie innym miejscu instalacji,
ponieważ są to dwa niezależne ustawienia.

Odwzorowany most MCP nosi nazwę mcp-danaco-pulpit-console. Skrypt
uruchamiający bierze tryb z argumentu albo ze zmiennej środowiskowej
DANACO_MOST_TRYB; brak obu oznacza tryb odczytu. Proces mostu ogranicza
działanie do korzeni podanych zmienną środowiskową rozdzielonych
dwukropkiem; poza te korzenie most nie wychodzi, a w trybie odczytu
narzędzia zapisu nie są ogłaszane w wykazie narzędzi.

Kolumny rodzaj, tryb_domyslny, tryb i stan niosą wartości kontraktu wprost:
mcpBridge, localDirectory, read, write, unknown, reachable, unreachable.
Kontrakt nie deklaruje przy tych wyliczeniach pola baza, w odróżnieniu od
wyliczenia rodzaju konta, które to pole ma. Własny przekład polsko-angielski
w warstwie trwałości byłby drugim źródłem przekładu obok pliku kontraktu.

Słowa odczyt i zapis są słowami argumentu uruchamiającego dany most, nie
wartościami wyliczenia platformy. Są własnością konkretnego mostu, więc leżą
w danych tabeli argument_trybu_mostu, nie w kodzie: inny most może żądać
innych słów, a to oznacza nowy wiersz danych, nie nową gałąź kodu.

Kolumna identyfikator_zewnetrzny nadania dostępu niesie identyfikator
tekstowy rdzenia; wartość pusta oznacza nadanie założone wprost w bazie
danych, bez odpowiednika w pamięci rdzenia.

## budowa/server/internal/store/migracja_002_okna.sql

Wartości wyliczeniowe kolumn tego obszaru pochodzą z kontraktu —
`shared/contract.json` jest jedynym źródłem prawdy. Kontrakt zapisuje nazwy
po angielsku, model danych po polsku; odwzorowanie jest jeden do jednego i
leży wyłącznie w kontrakcie, w polu `baza` przy każdej wartości wyliczenia.
Generator wytwarza z kontraktu słowniki `WartosciBazy*` i `WartosciKontraktu*`
w `contract.go` i `contract.ts`, a warstwa trwałości sięga po te słowniki
zamiast wpisywać przekład u siebie. Odwzorowanie kolumn na typy kontraktu:
`okno_komunikacji.srodowisko_wykonania` na ExecutionEnv,
`okno_komunikacji.tryb_uprawnien` na PermissionMode,
`okno_komunikacji.rola_okna` na WindowRole,
`okno_komunikacji.stan` na WindowStatus,
`proces_sesji.stan` na ProgressStatus,
`wiadomosc.rola` na MessageRole,
`wiadomosc.stan` na MessageStatus,
`wiadomosc.rodzaj_tresci` na ChunkKind.

Rola i persona wiadomości niosą rozróżnienie odrębne: rola mówi, czym jest
nadawca w rozumieniu kontraktu (MessageRole, cztery wartości), persona mówi,
w jakiej funkcji nadawca wystąpił w danej wypowiedzi. Persona jest warstwą
prezentacji i nie mnoży wartości pola rola; pozostaje pusta dla zwykłej
wypowiedzi użytkownika i modelu. Atrybucję wypowiedzi w pętli
koordynator-wykonawca niesie kolumna okno_zrodlowe_id: okno źródłowe zna
własną rola_okna, więc nie jest potrzebny drugi, równoległy słownik ról.

## budowa/server/internal/store/migracja_402_nastawy_przesiewu_i_obrazu.sql

Klucze `wiedza_model_przesiewu`, `wiedza_katalog_przesiewu`, `wiedza_model_obrazu`
i `wiedza_katalog_obrazu` stoją w kodzie od czasu dołożenia przesiewu i osi obrazu
(`wiedza/ustawienia.go`), ale wiersza w `definicja_ustawienia` nie miały. Zdolność
mimo to działała: rozstrzyganie nastawy czyta zapis niezależnie od katalogu
definicji (`konfig/rozstrzyganie.go`), więc wartość zapisana wprost w tabeli
`ustawienie` dochodziła do silnika. Czego bez wiersza katalogu nie było, to drogi
dla operatora okna konfiguracji: `config.set` odmawia klucza spoza katalogu,
a okno wystawia wyłącznie pozycje katalogu. Cztery nastawy były więc ustawialne
ręcznym zapisem do bazy i tylko nim.

Migracji 115 się nie zmienia — jej suma kontrolna stoi w rejestrze `migracja`
u każdego, kto rdzeń postawił, a niezgodność sumy wywraca start rdzenia
(`store/migracje.go`). Cztery wiersze dokłada więc osobny krok, wzorem tego,
jak migracja 401 zmieniła wartości domyślne dwóch nastaw z migracji 115.

Kategoria, zasięg i oś są te same, co u czterech nastaw migracji 115, i to
z tych samych powodów. Kategoria `wiedza`, bo to ten sam silnik. Zasięg
wyłącznie globalny, bo wskaźnik znaczenia jest jeden na maszynę: przesiew
układający kolejność dwoma różnymi koderami w dwóch oknach dawałby dwie
nieporównywalne kolejności tego samego wyniku. Oś wyłącznie `platform`, bo
katalog wag i nazwa modelu liczącego lokalnie są własnością maszyny, a nie
konta ani kanału modelu.

Wartości domyślne są kopią stałych `wiedza/ustawienia.go` co do znaku —
`ModelPrzesiewuDomyslny`, `ModelObrazuDomyslny` i `katalogNiewskazany`. Rozjazd
znaczyłby dwie prawdy o tym, czym rdzeń liczy, zależne od drogi wywołania:
rozstrzygacz zasięgu oddaje wartość z tej kolumny, a stała pakietu wchodzi tam,
gdzie rozstrzygacza nie ma (`core/adapter_modul_wiedza.go`).

Oba katalogi wag zostają puste, choć wagi obu modeli leżą na maszynie
(`/opt/danaco-modele/reranker`, `/opt/danaco-modele/clip`). Wartość niepusta
byłaby tutaj drugą prawdą wobec stałej `katalogNiewskazany`, a `internal/wiedza/`
leży poza terenem tej zmiany. Wskazanie wag stojących jest osobnym krokiem,
obejmującym zarazem stałą i ten wiersz — dokładnie tak, jak migracja 401
zrobiła to dla osadzarki.

Klauzula `ON CONFLICT DO NOTHING` czyni krok idempotentnym i nieszkodliwym na
bazie, gdzie te wiersze z jakiegoś powodu już stoją.

## budowa/server/internal/store/migracja_007_zaczyn_slownikow.sql

Łańcuch `srodowisko → modul → karta_sesji → sesja → okno_komunikacji →
wiadomosc` stoi na więzach klucza obcego: bez wierszy w `srodowisko` i `modul`
nie da się zapisać żadnego okna komunikacji ani żadnej wiadomości, dlatego ta
migracja zakłada je jako pierwsza.

Klauzula `WHERE true` poprzedzająca `ON CONFLICT` w zapytaniach `INSERT ...
SELECT` jest wymogiem składni SQLite dla upsertu tej postaci, nie ozdobnikiem.

Macierz widoczności modułów w środowiskach odpowiada kolejności pozycji
w nawigacji bocznej interfejsu. Automations nie ma okna modułowego w żadnym
środowisku, a MultitaskingAI nie udostępnia modułów — oba są w niej celowo
nieobecne.

## budowa/server/internal/wiedza/pomocnik_obrazu.py

Model osi obrazu jest dwuwieżowy: ma osobną wieżę dla obrazu i osobną dla
tekstu, a obie kończą w jednej przestrzeni, w której iloczyn skalarny znaczy,
że dane zdanie opisuje dany obraz. Osadzarki tekstu nie da się tu użyć, bo
jej wektor leży w przestrzeni, w której obrazu nie ma i nigdy nie było, a
porównanie dałoby liczbę bez związku z czymkolwiek.

Obrazy przychodzą do pomocnika ścieżkami plików, nie bajtami. Bajty
biblioteki leżą w magazynie treści rdzenia i tamtą ścieżką czyta je podgląd
modułu Library; przepisanie ich do zlecenia oznaczałoby drugi komplet
obrazów w pliku JSON, rosnący w megabajtach na każde zapytanie. Droga
rozmowy z rdzeniem jest ta sama co u pomocnika osadzeń i u pomocnika
przesiewu: zlecenie ścieżką pliku JSON w argumencie, odpowiedź jednym
obiektem JSON na standardowym wyjściu, diagnostyka biblioteki na strumieniu
diagnostycznym.

Katalog modeli bywa dwiema rzeczami i pomocnik je rozróżnia tak samo jak
pomocnik osadzeń: pusty katalog jest miejscem na wagi, a katalog z wagami
jest samym modelem, i wtedy pobieranie z sieci jest wyłączone. Rozstrzyga
o tym obecność pliku wag modelu stojącego w katalogu. Katalog pobrania jest
podawany bibliotece jawnie, bo bez zmiennej środowiskowej HOME biblioteka
nie zna swojego domyślnego katalogu pamięci podręcznej, a katalog docelowy
wag jest ustawieniem, które kontroluje Operator.

Obraz nieczytelny nie przerywa całego zapytania: jeden plik uszkodzony albo
w postaci, której biblioteka obrazu nie otwiera, nie ma prawa odebrać
odpowiedzi o pozostałych obrazach. Taka pozycja dostaje ocenę zerową i
wraca w wykazie pominiętych, żeby rdzeń wiedział, że nie porównał
wszystkiego, o co prosił.

## budowa/server/internal/store/migracja_045_biblioteka.sql

Treść pliku trzyma dysk, nie baza: kolumna `tresc_odwolanie` niesie odwołanie do
pliku na dysku, baza nie dostaje kolumny BLOB. Ten sam wzorzec powtarza się na
dwóch poziomach — plik biblioteki i każda jego wersja — bo
`library.version.restore` musi umieć przywrócić treść sprzed zmiany, więc treść
poprzednich wersji musi przeżyć nadpisanie bieżącej.

Wersja jest własnym bytem, nie polem licznika. `library.version.list` zwraca
listę `LibraryVersion` z własnym `id`, autorem i sumą kontrolną — to nie jest
rosnący numer przy pliku, tylko osobny wiersz historii, bo każda wersja niesie
własną treść i własnego autora (kontrakt: `LibraryVersion.author`).

Etykiety i kolekcje mają rozłączne tabele — to dwie różne prawdy o pliku:
etykieta jest wolnym tekstem bez własnej tożsamości, kolekcja jest bytem z nazwą
i opisem tworzonym osobną komendą `library.collection.create`. Stąd etykieta
żyje jako wiersz w tabeli złącznikowej z gołym tekstem, a przypisanie do
kolekcji odwołuje się do wiersza `kolekcja_biblioteki`.

Podgląd (`LibraryPreview`) nie ma własnej tabeli. To widok obliczany w locie
z bieżącej wersji pliku (rodzaj podglądu wynika z `mime_type`, treść
z `tresc_odwolanie` wersji) — trwały byt tu jest jeden: wersja pliku.

## budowa/scripts/instalka-hybryda-win-x64.sh

Produkt ma jedną postać: hybrydę. U operatora staje samo okno, a rdzeń,
serwer narzędzi i arsenał stoją na serwerze wdrożenia. Instalka rdzenia nie
niesie i nieść nie ma — to jest sens tego produktu, nie oszczędność na
rozmiarze: rdzeń jest jeden, utrzymuje go administrator w jednym miejscu,
a okno u operatora nie ma czego aktualizować poza sobą.

Powłoka zna trzy stany wskazania rdzenia: wskazanie niezłożone, rdzeń na tym
urządzeniu, rdzeń na serwerze. Ta instalka jest dla trzeciego stanu. Po
instalacji obowiązuje jeszcze stan pierwszy — powłoka nie stawia niczego
i czeka, aż operator wskaże host rdzenia w oknie albo aż wskaże go zmienna
środowiskowa DANACO_HOST_RDZENIA.

Mimo braku rdzenia instalka niesie pakiet interfejsu — nie jako osobny plik,
lecz w środku pliku wykonywalnego powłoki. Ustawienie frontendDist w
konfiguracji Tauri jest ustawieniem budowy, nie zasobem instalatora: pakiet
klienta zostaje wkompilowany w binarkę powłoki i nie da się go z instalki
wyjąć, nie odbierając powłoce wyjścia awaryjnego.

Skrypt woła złożenie pakietu, nie pełną budowę: pełna budowa uruchomiłaby
przebudowę klienta, a ta do budowy powłoki nie należy; złożenie pakietu
niczego nie kompiluje i bierze gotową binarkę z katalogu docelowego. Nakładki
konfiguracyjnej nie ma i nie jest potrzebna, bo konfiguracja opisuje wprost
produkt hybrydowy, ponieważ innego produktu nie ma. Skrypt nie podpisuje
instalatora — podpis Authenticode wymaga certyfikatu i hosta Windows — i nie
sprawdza, czy instalator się uruchamia, bo tego nie da się sprawdzić bez
maszyny z Windows.

Cecha budowy tauri/custom-protocol nie jest ozdobna: bez niej gotowy plik
zachowuje się jak budowa deweloperska, a rozstrzygnięcie źródła interfejsu
zatrzymuje się na drugim warunku, nigdy nie sprawdzając warunku „rdzeń
nasłuchujący" — jedynego, na którym ten produkt stoi. Cecha jest podana
z wiersza poleceń, nie dopisana do pliku manifestu budowy, bo ten manifest
jest wspólny dla obu celów budowy i nie należy wyłącznie do tego skryptu.

Zapora rozstrzyga o osadzonych zasobach, a nie o obecności adresu serwera
rozwojowego w binarce: konfiguracja produktu wchodzi do binarki zawsze,
niezależnie od tego, która droga do interfejsu jest czynna, więc taka zapora
milczałaby przy pliku zepsutym. Sondą jest nazwa pliku interfejsu z sumą
treści w nazwie — taki napis nie ma jak trafić do binarki inaczej niż przez
osadzenie pakietu klienta, i zmienia się z każdą przebudową klienta, dlatego
skrypt czyta go z katalogu, a nie wpisuje na sztywno. Sondy oparte na
rozszerzeniu czcionki albo na pliku indeksu są mylące, bo obie zapalają się
w binarce także bez osadzonych zasobów.

Cel złożenia zawiera 7z, ponieważ tylko wykaz zawartości gotowego pliku
dowodzi, że rdzeń nie wrócił do produktu hybrydowego przy zmianie
konfiguracji. Ostatnia zapora bada binarkę wypakowaną z instalatora, a nie
tę z katalogu docelowego, bo tylko wypakowana binarka jest tym, co dostanie
operator.

## budowa/server/internal/store/migracja_074_petla.sql

Migracja dokłada do pętli wykonawczej automatyki trzy brakujące elementy: krok
rodzaju `model` nie umiał wskazać agenta z portfolio, obsada od jednego do
czterech uczestników z podziałem na role nie miała gdzie się zapisać, a bieg
zawieszony i czekający na reakcję ze świata nie miał odrębnego stanu od zwykłej
pauzy. Pauza jest w schemacie od dawna: krok rodzaju `wait` czeka określony czas
i idzie dalej, a stan `paused` czeka na operatora — w obu wypadkach wiadomo,
kiedy bieg ruszy. Wybudzenie czeka na reakcję świata, która może przyjść za
godzinę, za trzy dni albo nigdy, i stąd trzy wymagania: bieg musi przetrwać
restart rdzenia (stan żyje w tabeli `oczekiwanie_biegu`, nie w pamięci procesu),
musi istnieć wiązanie między biegiem a sygnałem, na który czeka (inny byt niż
`wyzwalacz_automatyki`, który uruchamia automatykę od początku), i musi być
rozstrzygnięte, co się dzieje, gdy reakcja nie przyjdzie — stąd `termin`
i `po_terminie` jako pola obowiązkowe. Żaden z tych bytów niczego operatorowi
nie zabrania: brak wiersza obsady oznacza bieg na modelu wskazanym w kroku, brak
wiersza oczekiwania oznacza, że bieg na nic nie czeka.

Wartości roli obsady (`koordynator`, `wykonawca`, `samodzielne`) są przepisane
z bazy wyliczenia `WindowRole` kontraktu (pole `baza` przy
`coordinator`/`executor`/`standalone`), tak jak rozstrzyga to klient
w `okna-rownolegle/role-domyslne.ts` — własny słownik ról obok kontraktowego
rozjeżdżałby się z nim przy pierwszej zmianie. Komunikacja między uczestnikami
obsady wykracza poza tę migrację: wymaga nowej rodziny komend kontraktu, którą
rozstrzyga osobny pakiet.

Tabela `agent_kroku_automatyki` nie dokłada kolumny do `krok_automatyki`, bo
`ZapiszKroki` (`dane/automations_kroki.go`) przy każdym zapisie z Workflow
Buildera kasuje i wstawia na nowo komplet kroków, wymieniając tylko osiem
kolumn, które zna — dziewiątej by nie przepisał i cicho zdejmował agentów
z kroków. Wiązanie idzie więc tą samą drogą co `zaleznosc_kroku_automatyki`:
przez identyfikator kroku zewnętrzny (TEXT), nie przez klucz liczbowy, bo
wiersze kroków bywają podmieniane w całości. `ON DELETE CASCADE` po stronie
agenta zabiera wyłącznie wiązanie, nigdy krok — krok bez wiązania wraca do
zachowania sprzed niego, czyli biegnie na modelu z `parametry`. Wiersz
wiązania, który zostaje po skasowaniu kroku z definicji, nigdy nie jest
odczytany i znika kaskadą przy skasowaniu całej automatyki; sprzątanie takich
osieroconych wiązań należałoby do zapisu kroków, nie do schematu.

Granica czterech miejsc obsady jest specyfikacją, nie techniczną barierą: klient
dziś liczy wykonawców do dwóch i szuka analityka wśród okien samodzielnych.
Więz `CHECK(miejsce BETWEEN 1 AND 4)` razem z `UNIQUE(automatyka_id, miejsce)`
trzyma tę granicę jedną liczbą w jednym więzie, bez licznika w kodzie i bez
wyzwalacza. Więz `CHECK(agent_id IS NOT NULL OR model IS NOT NULL)` nie
pozwala zapisać puste miejsce, które w wykazie wyglądałoby tak samo jak miejsce
wypełnione.

Tabela `przebieg_automatyki` idzie przez przebudowę (kopia z poprawionym
więzem, przepisanie wierszy, podmiana nazwy, odtworzenie indeksów), bo SQLite
nie zna polecenia zdejmującego więz CHECK. Przebudowa obejmuje wyłącznie tę
tabelę: żaden klucz obcy na nią nie wskazuje, więc `DROP TABLE` nie zdejmuje
żadnych dzieci. Stan `oczekuje` dokłada się do sześciu istniejących bez ruszania
pozostałych, więc żaden zapisany wcześniej bieg nie zmienia stanu ani nie
znika. Kontrakt tego stanu jeszcze nie niesie — `AutomationExecutionStatus`
zna sześć wartości — ale schemat idzie pierwszy, bo bez miejsca w bazie bieg
czekający nie przetrwałby restartu.

Tabela `oczekiwanie_biegu` odpowiada na inne pytanie niż wyzwalacz: wyzwalacz
rozstrzyga, kiedy zacząć nowy bieg, oczekiwanie — na co czeka bieg już
rozpoczęty i od którego kroku ma ruszyć dalej. Pole `po_terminie` ma trzy
wartości, z których każda jest świadomą decyzją: `wznow` (domyślne, termin
minął — ruszaj dalej, jakby sygnał przyszedł, bo produkt ma pracować dalej,
a nie stawać), `ponow` (wykonaj krok oczekiwania jeszcze raz, `proba` rośnie)
i `przerwij` (zakończ bieg stanem `stopped` z jawnym powodem). Termin pusty
(`termin IS NULL`) jest wyborem operatora — czekaj bez końca — i taki wiersz
nie trafia do budzika terminów, więc nic nie kosztuje. Indeks częściowy UNIQUE
po `przebieg_id` z warunkiem `wybudzono IS NULL` wymusza jedno czynne
oczekiwanie na bieg naraz: oczekiwania zamknięte zostają w tabeli jako ślad.

## budowa/server/internal/store/migracja_370_studio_petla_i_agenci.sql

Oba narzędzia — pętla wykonawcza i praca dwóch agentów naraz — są domyślnie
wyłączone i włączane jawnym, odwracalnym ustawieniem Operatora. Nastaw nie ma
w tym schemacie z zamysłu: idą zasięgami rodziny `config.*`, która ma własny
magazyn. Drugiego magazynu ustawień Studio migracja nie zakłada, ponieważ
inaczej Operator wyłączyłby pętlę w jednym miejscu, a ona chodziłaby dalej
wedle drugiego.

Zlecenie dokumentowe rozłożone na zadania musi przetrwać przeładowanie rdzenia
i zamknięcie okna: pętla wykonawcza chodzi obiegami, a obieg drugi ma wiedzieć,
co zrobił obieg pierwszy. Rozkład trzymany w pamięci procesu znaczyłby, że
każde przeładowanie zaczyna pracę od nowa na dokumencie, który już jest w pół
przerobiony, dlatego rozkład zlecenia jest bytem trwałym.

Przy dwóch agentach naraz zadanie w biegu bez wskazania wykonawcy nie mówi
niczego: nie wiadomo, czy stoi, czy ktoś nad nim pracuje, ani kogo zapytać
o wynik. Tożsamością wykonawcy jest kod agenta z modułu Agents, tak samo jak
przy zmianie śledzonej, dlatego zadanie niesie wykonawcę, a nie tylko stan.

Wykonawca ubity w pół pracy nie może trzymać zajętego fragmentu dokumentu na
zawsze — drugi agent stałby bezczynnie, a Operator nie wiedziałby, dlaczego.
Stąd kolumna `wygasa`: zajęcie bez odnowienia przestaje obowiązywać samo, jako
zapora przed zakleszczeniem, a nie jako rozjemca sporu.

Dwóch wykonawców zmieniających ten sam fragment dokumentu nie może dać
dokumentu, w którym jeden nadpisał drugiego bez śladu. Rdzeń nie orzeka, kto
ma rację — zapisuje wiersz spięcia z prawdą o tym, co się stało: czyja zmiana
weszła, czyja została odłożona i dlaczego. Odłożone brzmienie zostaje w
wierszu, więc nie przepada: Operator może je wnieść sam albo znaleźć jako
propozycję na marginesie.

## budowa/server/internal/store/migracja_129_prowenancja_wywolan.sql

Podsystem prowenancji: rdzeń zapisuje, co poszło do modelu i co wróciło. Bez
tego zapisu Provenance Explorer modułu Diagnostics nie ma czego pokazać,
a rozliczenie zużycia nie ma po czym liczyć — dlatego obie rodziny komend stoją
na tej jednej tabeli, a nie na dwóch osobnych. Dwie tabele o tych samych liczbach
rozjechałyby się przy pierwszej korekcie cennika.

Prompt i odpowiedź to kolumny osobne, ale wiersza tego samego. Wyłączenie zapisu
treści ustawieniem zostawia wiersz z liczbami i przebiegiem, a kolumny treści
puste — ślad wywołania zostaje mierzalny nawet wtedy, gdy jego treść świadomie
nie jest przechowywana. Kolumna `tresc_zapisana` odróżnia treść pustą od treści
niezapisanej: bez niej wiersz bez promptu znaczyłby dwie różne rzeczy naraz.

Odcinek wywołania (wejście modelu, użycie narzędzia, tura podagenta) niesie
wskazanie rodzica. Dzięki temu zapytanie o jeden odcinek czyta jeden wiersz,
zamiast rozbierać cały dokument JSON, żeby dojść do gałęzi drzewa.

## budowa/server/internal/core/adapter_narzedzia_obraz_pomocnik_twarzy.py

Rdzeń jest napisany w Go, a GFPGAN jest wydany jako wagi PyTorcha
(`GFPGANv1.4.pth`). Przepisanie tej sieci do Go byłoby drugą implementacją
cudzej architektury i rozjeżdżałoby się z wagami przy każdym kolejnym
wydaniu modelu, więc przebieg twarzowy jest procesem obok rdzenia — tak samo
jak liczenie wektorów znaczenia (`internal/wiedza/pomocnik_osadzen.py`) i
rozpoznawanie mowy (`internal/mowa/pomocnik.go`). Zlecenie przychodzi
czterema argumentami wiersza poleceń, a odpowiedź wraca jednym obiektem
JSON na standardowym wyjściu. Diagnostyka bibliotek idzie na strumień
diagnostyczny, bo ostrzeżenie wstawione w środek JSON-a uczyniłoby
odpowiedź nieczytelną. Standardowego wejścia pomocnik nie dostaje, bo
funkcja wywołania procesów zewnętrznych go nie podaje.

Przebieg jest osobny, a nie wpięty w powiększanie: kontrakt `image.upscale`
nazywa to polem `faces`, opisanym jako osobny przebieg, i tak też jest
liczone. Najpierw Real-ESRGAN powiększa cały obraz, potem ten pomocnik
odnajduje w wyniku twarze, odtwarza każdą z osobna w rozdzielczości
512×512 i wkleja ją z powrotem. Wymiary wyniku pochodzą więc wyłącznie z
powiększenia — sieć twarzowa ich nie rusza, i odpowiedź kontraktu niesie
te same wartości szerokości i wysokości, co przebieg bez poprawki twarzy.

Wybór sieci GFPGAN zamiast CodeFormer jest świadomy, nie brakiem: obok
wag GFPGAN leży też `codeformer.pth`, ale rdzeń go nie woła, ponieważ
CodeFormer stoi na własnej architekturze (VQGAN wraz z transformerem
przewidującym kod słownika), której wydanie nie niesie w wagach — trzeba
by wnieść drugi zestaw cudzego kodu obok tego, którym GFPGAN już liczy.

Gdy biblioteki nie ma albo wagi są nie do wczytania, pomocnik oddaje
odpowiedź z polem ok ustawionym na false i polem powod, po czym kończy
pracę kodem zerowym; rdzeń zamienia to na odmowę nazywającą brak. Ślad
stosu Pythona sam z siebie nie powiedziałby Operatorowi, czego brakuje.

## budowa/scripts/instalka-hybryda-win-arm.sh

Produkt ma jedną postać: hybrydę. Ta instalka wiezie samą powłokę; rdzeń,
serwer narzędzi i arsenał stoją na serwerze wdrożenia, a operator dostaje
okno, nie drugą kopię serca platformy. Instalka jest cienka i nie sprawdza
obecności rdzenia, bo jego brak nie jest tu usterką, jest założeniem.

Celem budowy jest architektura aarch64-pc-windows-msvc, składana przez
cargo xwin, które pobiera i trzyma zestaw nagłówków oraz bibliotek importu
Microsoftu. Kompilatorem krzyżowym jest clang, a bibliotekarzem llvm-lib
z wydania llvm-mingw, dlatego katalog binarny tego wydania musi być na
ścieżce tej budowy. Cel x64 ma wymaganie odwrotne: idzie przez systemowe
mingw-w64, a llvm-mingw na ścieżce mu przeszkadza — oba skrypty ustawiają
ścieżkę same, żeby jedna budowa nie psuła drugiej.

Skrypt woła złożenie pakietu, nie pełną budowę: pełna budowa uruchomiłaby
przebudowę klienta, a ta do budowy powłoki nie należy; złożenie pakietu
niczego nie kompiluje i bierze gotową binarkę z katalogu docelowego. Nakładki
konfiguracyjnej nie ma i nie jest potrzebna, bo konfiguracja opisuje wprost
produkt hybrydowy, ponieważ innego produktu nie ma. Skrypt nie podpisuje
instalatora — podpis Authenticode wymaga certyfikatu i hosta Windows — i nie
sprawdza, czy instalator się uruchamia, bo tego nie da się sprawdzić bez
maszyny z Windows na architekturze ARM64.

Cecha budowy tauri/custom-protocol nie jest ozdobna: bez niej produkt jest
zepsuty w sposób niewidoczny, bo rozstrzygnięcie źródła interfejsu
zatrzymuje się na sprawdzeniu budowy deweloperskiej, nigdy nie sprawdzając
warunku „rdzeń nasłuchujący" — jedynej drogi do interfejsu dla hybrydy;
okno pokazałoby wtedy pustą stronę.

Sprawdzenie architektury binarki dopuszcza zarówno pisownię „ARM64", jak
i „Aarch64", bo nazwa architektury zależy od wersji narzędzia odczytu typu
pliku, a nie od produktu; sprawdzenie stoi przed złożeniem instalatora, bo
instalka z binarką dla innej architektury w środku byłaby bezużyteczna.

Zapora rozstrzyga o osadzonych zasobach, a nie o obecności adresu serwera
rozwojowego w binarce: konfiguracja produktu wchodzi do binarki zawsze,
niezależnie od tego, która droga do interfejsu jest czynna, więc taka zapora
milczałaby przy pliku zepsutym. Sondą jest nazwa pliku interfejsu z sumą
treści w nazwie — taki napis nie ma jak trafić do binarki inaczej niż przez
osadzenie pakietu klienta, i zmienia się z każdą przebudową klienta, dlatego
skrypt czyta go z katalogu, a nie wpisuje na sztywno. Sondy oparte na
rozszerzeniu czcionki albo na pliku indeksu są mylące, bo obie zapalają się
w binarce także bez osadzonych zasobów.

Cel złożenia zawiera 7z, ponieważ tylko wykaz zawartości gotowego pliku
dowodzi, że rdzeń nie wrócił do produktu hybrydowego przy zmianie
konfiguracji. Ostatnia zapora bada binarkę wypakowaną z instalatora, a nie
tę z katalogu docelowego, bo tylko wypakowana binarka jest tym, co dostanie
operator.

## budowa/server/internal/store/migracja_362_studio_obiekty_i_aparat.sql

Obraz, kształt, ikona i pole tekstowe są bytami, po których się pyta
niezależnie od tego, gdzie w treści wiszą: które obiekty w dokumencie
przyszły z modułu Design, co jest osadzone z bazy zdjęciowej, które obiekty
są bez tekstu zastępczego. Drzewo postaci trzyma tylko zakotwiczenie
obiektu, a jego opis leży w tabeli obiektu — dlatego wstawienie obiektu nie
przepisuje całego drzewa, a wykaz obiektów nie wymaga jego rozbierania.
Bajty obiektu leżą w magazynie zasobów pod sumą kontrolną, tym samym, którym
posługuje się warsztat PDF i moduł Design; drugiego magazynu Studio nie
zakłada, kolumna trzyma odwołanie, nie zawartość. Fragment i obraz
wciągnięty ze strony albo z Biblioteki niesie zapis pochodzenia przy samym
obiekcie, bo obiekt ma jedno źródło i nie dzieli się na fragmenty; dla
treści dzielącej się na fragmenty służy do tego osobna tabela z migracji 368.

Spis treści, spis ilustracji, przypis dolny i końcowy, podpis, zakładka,
odwołanie wzajemne, odsyłacz, powołanie, bibliografia, hasło i indeks
różnią się tym, co niosą, ale nie tym, jak się nimi pracuje: każdy jest
przypięty do miejsca w treści, każdy ma numer nadawany przy odświeżeniu i
każdy może być nieświeży. Trzynaście tabel o tym samym kształcie
znaczyłoby trzynaście zapytań przy każdym odświeżeniu aparatu, więc rodzaj
elementu rozstrzyga kolumna rodzaj, a to, co swoiste dla rodzaju, leży w
kolumnie dane_json.

Pole jest przeliczane, nie zbierane: numer strony, liczba stron, data i
pole obliczane liczą się z układu i z właściwości dokumentu, a nie z
nagłówków. Odświeżenie pól i odświeżenie aparatu są dwiema czynnościami o
różnym koszcie i różnej porze — pola odświeżają się przy każdym podglądzie,
a elementy aparatu dopiero na żądanie.

## budowa/server/internal/store/migracja_001_fundament.sql

Migracja przyjmuje spójne konwencje dla całego modelu danych: nazwy tabel
i kolumn po polsku w zapisie snake_case, klucz główny jako
`INTEGER PRIMARY KEY AUTOINCREMENT`, daty jako `TEXT` w formacie ISO 8601,
wartości logiczne jako `INTEGER` z ograniczeniem `CHECK(... IN (0,1))`
i wyliczenia jako `TEXT` z ograniczeniem `CHECK`.

## budowa/server/internal/store/migracja_048_design.sql

Prompt dostaje własną tabelę, nie kolumny powielone w wierszu zasobu.
Kontrakt niesie identyfikator promptu jako pole opcjonalne w DesignPrompt, a
DesignAsset odwołuje się do promptu osobnym polem — dwa sygnały, że prompt
bywa bytem trwałym, nie tylko parametrem jednego wywołania generowania
zasobu. Warianty i generowanie obraz-do-obrazu zakładają wprost, że wiele
zasobów powstaje z tego samego promptu; gdyby prompt żył jako kolumny
powielone w każdym wierszu zasobu, każdy wariant niósłby własną kopię tych
samych pól. Tabela prompt_design jest jedną prawdą o promptcie, a kolumna
prompt_id w tabeli zasobu jest jedynym miejscem odwołania do niej.

Warstwa kompozycji dostaje własną tabelę, nie zapis strukturalny w
kolumnie. DesignBoardLayer niesie własną pozycję, wymiary, kolejność i
znacznik zablokowania — pola, po których trzeba by filtrować i sortować
przy odczycie, gdyby leżały w jednym polu JSON. Zapis warstw kompozycji
nadsyła zawsze całą listę na nowo, bez trybu częściowej zmiany, więc zapis
jest zawsze usunięciem warstw istniejących i wstawieniem przysłanych od
nowa.

Treść zasobu trzyma dysk lub usługa zewnętrzna, nie baza: pole adresu
zasobu w kontrakcie już jest odnośnikiem, nie surową treścią, więc kolumna
uri przechowuje ten odnośnik wprost, bez pośredniej kolumny na dane
binarne. Etykiety zasobu mają własną tabelę złącznikową, bo etykieta jest
wolnym tekstem bez własnej tożsamości, więc para złożona z zasobu i
etykiety jest kluczem bez sztucznego identyfikatora.

## budowa/server/internal/store/migracja_376_ustawienia_pamieci_i_wyciszenia.sql

Rodzina `config.*` niesie odczyt i zapis ustawienia zasięgiem ogólnym
(`config.get`, `config.set`, `config.reset`), a `settings.category.list`
wraz z `settings.definition.list` budują z katalogu formularz okna
Konfiguracji. Drugiej drogi komend zakres nie wymaga: brakowało wyłącznie
wierszy katalogu, więc okno nie miało czym sterować, a migracja wnosi te
wiersze bez dodania żadnej komendy.

Zakres modelu konfiguracji nazywa cztery pozycje pamięci: poziom pamięci,
stan włączenia pamięci na poziomie, odłączenie pamięci w sesji i zawartość
zasobu pamięci. Zawartość zasobu nie jest ustawieniem katalogu — jest treścią
wpisu, bytem tabeli `wpis_pamieci_projektu`, i idzie rodziną `memory.*`; do
katalogu wchodzą więc trzy pierwsze pozycje, a czwarta pozostaje w pamięci.
Wpisanie treści wpisu jako ustawienia dałoby dwa magazyny jednego bytu.

Reguły wyciszania nakładki Always On Display są jedną pozycją warstwy
globalnej — zakresy i czasy wyciszenia dostępne w menu nakładki. Wyciszenia
czynne nie są ustawieniem: mają własną tabelę, ponieważ powstają i giną
w toku pracy, a nie przy nastawianiu platformy.

Model konfiguracji daje pamięci wszystkie cztery warstwy ogólne (globalna,
środowisko, projekt, sesja), a odłączeniu pamięci w sesji wyłącznie warstwę
sesji. Reguły wyciszania obejmują wyłącznie warstwę globalną. Zasięgi zapisu
w migracji idą dokładnie za tym podziałem: poziomu szerszego niż wskazany nie
dokłada się żadnej pozycji, bo zapis na poziomie, którego zakres nie obejmuje,
byłby nastawą bez wskazanego pochodzenia.

Objaśnienia kontekstowe w kolumnie `opis` odpowiadają na pytania, co
ustawienie robi i jaki ma wpływ; przy wyłączeniu pamięci odpowiadają też na
pytanie o los treści, bo bez tego nie da się odróżnić wyłączenia od
usunięcia — treść zostaje. Żadna z pozycji nie wymaga restartu: pamięć czyta
rozstrzygacz przy każdym złożeniu kontekstu, a reguły wyciszania — nakładka
przy każdym otwarciu menu.

## budowa/server/internal/store/migracja_009_zaczyn_akcji.sql

Wiersze pochodzą z `shared/contract.json`, z sekcji `komendy` (nazwa komendy,
opis, pola obowiązkowe żądania) oraz z sekcji `narzedzia`. Zaczyn obejmuje te
komendy, które kontrakt wskazuje jako sterowanie platformą; poza wykazem
zostają `connection.hello` i `session.bind`, ponieważ są czynnościami warstwy
połączenia klienta, nie akcjami panelu. Osobno wchodzi pasek narzędzi promptu
każdego modułu: każdy moduł niesie okno rozmowy, a jego pasek promptu niesie
wysłanie polecenia, zatrzymanie odpowiedzi i historię poleceń. Akcja
wskazująca komendę spoza kontraktu byłaby pozycją, której nie da się wywołać,
więc do zaczynu nie wchodzi — takie akcje dochodzą wierszami, bez zmiany kodu,
gdy ich komendy wejdą do kontraktu.

Poziom zasięgu akcji jest wyprowadzony z bytu, na którym komenda działa:
`environment.enter` działa na środowisku, komendy sesji i kolejek na karcie
sesji, komendy okna na oknie komunikacji, komendy wiadomości na oknie czatu
modułu, a pozostałe na całej platformie.

Każde wstawienie kończy się `ON CONFLICT(kod) DO NOTHING`, dzięki czemu
migracja przechodzi także na bazie, w której część wierszy już istnieje.
Klauzula `WHERE true` przed `ON CONFLICT` jest wymogiem składni SQLite dla
zapisu `INSERT ... SELECT` z upsertem.

## budowa/server/internal/store/migracja_124_urzadzenia_powiadomien.sql

Nośnik doręczenia już istnieje w warstwie transportu: rozgłoszenie serwera
wysyła kopertę do otwartych gniazd i zwraca liczbę urządzeń, które ją
przyjęły, a rejestr połączeń zna tożsamość każdego gniazda. Ta migracja nie
zakłada własnego kanału transmisji — opisuje wyłącznie to, czego nośnikowi
brakuje: kogo wołać, czym i czy doszło.

Nie ma tu drugiej tabeli urządzeń. Urządzenie opisuje osobna tabela
urzadzenie; druga tabela maszyn byłaby drugą prawdą o tym, czym operator
dysponuje. Tabela urzadzenie_powiadomien nie opisuje maszyny — opisuje zgodę
tej maszyny na wołanie i drogę, którą wołanie idzie. Jedna maszyna może mieć
kilka takich dróg — pulpit i przeglądarka to dwa osobne gniazda o dwóch
tożsamościach — więc relacja jest jeden do wielu, a nie kolumną doklejoną do
tabeli urządzeń.

Słownik kanałów jest zamknięty na to, co rdzeń dziś potrafi doręczyć: kolumna
kanał dopuszcza jedną wartość, oznaczającą połączenie. Warunek dopuszczający
kanał usługi zewnętrznej bez kodu, który go obsłuży, przepuściłby rejestrację,
której żaden takt nie doręczy — wiersz stanąłby w kolejce na zawsze. Kanały
zewnętrzne wejdą osobną migracją, razem z kodem, który je obsłuży.

Klucz kanału dla kanału połączenia jest tym samym napisem, który transport
już dziś niesie jako tożsamość klienta gniazda; migracja nie zakłada nowego
identyfikatora urządzenia, bo byłby drugą tożsamością tego samego gniazda.

Kolejka zastępuje wysyłkę wprost, bo bez tabeli powiadomienie zgłoszone przy
zamkniętej aplikacji znika. Kolejka sprawia, że brak odbiorcy jest stanem,
a nie ciszą: wiersz oczekujący z licznikiem prób mówi wprost, że wołanie się
odbyło i nie było komu odpowiedzieć.

Termin ważności jest obowiązkowy z rozmysłu: powiadomienie bez terminu
wisiałoby wiecznie, a po dłuższym postoju rdzenia operator dostałby lawinę
budzików o sprawach dawno nieaktualnych. Termin wymuszony schematem znaczy,
że każdy wołający musi odpowiedzieć na pytanie, do kiedy dane powiadomienie
ma sens.

Rodzaj i identyfikator bytu tworzą kotwicę miękką, bez klucza obcego.
Powiadomienie dotyczy czegoś — dziś przede wszystkim kroku wstrzymanego,
czekającego na słowo operatora. Klucza obcego do jednej tabeli tu nie ma, bo
powiadomienie ma z założenia dotyczyć różnych bytów — kroku, zlecenia, biegu
automatyki — a klucz obcy zamknąłby je na jeden byt i wymusił kolumnę na
każdy następny. Rodzaj i klucz bytu chodzą parą, a więz sprawdzający tę parę
stoi dalej w definicji tabeli, razem z pozostałymi więzami, ponieważ SQLite
nie pozwala wrócić do definicji kolumn po pierwszym więzie tabeli. Ceną
takiego rozwiązania jest brak kaskady: powiadomienie o bycie usuniętym
zostaje w kolejce i wygasa własnym terminem.

Doręczenie jest osobną tabelą, bo samo dotarcie bez wskazania, do którego
urządzenia, jest odpowiedzią nie do sprawdzenia — operator ma pulpit i telefon
naraz. Kolumna dostarczono w tabeli powiadomień mówi, że dotarło gdziekolwiek,
i to wystarcza kolejce do zamknięcia sprawy; tabela doręczeń mówi gdzie
i kiedy, i to jest odpowiedź, którą można pokazać operatorowi bez zmyślania.

## budowa/server/internal/store/migracja_078_komunikacja_biegu.sql

Model kończy turę, a następny dostaje jego wynik jako wejście. Zakres tego, co
widzi następny, wybiera się osobno dla każdego wiązania; domyślnie `artefakt`
niesie wyłącznie wytwór poprzednika — plik, poprawkę, odpowiedź — najmniejszy
i najtańszy w żetonach, bo walidator ocenia wynik, nie tok myślenia. Zakres
`streszczenie` niesie skrót wypowiedzi i wymaga osobnej tury modelu, więc ma
własne miejsce, w którym może zawieść. Zakres `wypowiedz` niesie całą wypowiedź
poprzednika — najbogatszy i najdroższy; przy czterech stanowiskach i długiej
pętli zalewa okno kontekstu. Jeden sztywny zakres byłby albo zbyt ubogi dla
koordynatora przekazującego zlecenie, albo zbyt drogi dla walidatora
oceniającego wynik.

Ślad narzędzi poprzednika jest polem osobnym (`ze_sladem_narzedzi`), domyślnie
wyłączonym: bywa większy od samej wypowiedzi i najczęściej jest szumem, ale
walidator sprawdzający, czy wykonawca uruchomił testy, bez niego nie ma czego
sprawdzić. Pole jest prostopadłe do `zakres` — wolno chcieć samego artefaktu ze
śladem i całej wypowiedzi bez śladu.

Pola wyłączającego podgląd tej rozmowy nie ma: ruch idzie tą samą drogą co
każda tura — fragmentami strumienia do okien — a wiersz w `wiadomosc_biegu`
zostaje jako ślad.

Ta migracja nie buduje drugiego mechanizmu widoczności kontekstu. Rodzina
`isolation.*` rozstrzyga warstwowo, co okno widzi z historii, pamięci
i kontekstu sąsiada. Podział jest ostry: ta migracja niesie to, co poprzednik
jawnie przekazuje następnemu, a profil izolacji niesie to, co następny może
zobaczyć z kontekstu poprzednika poza tym, co mu przekazano. Stanowisko obsady
wskazuje swój profil izolacji kolumną `profil_izolacji_id`.

Wiązanie tabeli `przekazanie_biegu` jest projektem, nie ruchem: „gdy skończy
stanowisko 2, jego artefakt ze śladem narzędzi idzie do stanowiska 3". Operator
układa te wiązania, projektując bieg, a silnik je czyta, gdy stanowisko kończy
pracę. Wiązanie wskazuje `od_miejsca` i `do_miejsca` — numery 1-4 z obsady, nie
identyfikator stanowiska — tak samo jak `zaleznosc_kroku_automatyki` wiąże
kroki po identyfikatorze zewnętrznym, a nie po kluczu wiersza. Powód ten sam:
obsadę wolno podmienić w całości, a plan rozmowy ma to przeżyć.

Tabela `wiadomosc_biegu` ma dwa pola treści, nie jedno. `tresc` niesie to, co
objął `zakres`; `slad_narzedzi` niesie ślad, jeśli wiązanie go żądało. Sklejenie
ich w jedno uniemożliwiłoby późniejsze pokazanie samej treści bez szumu.
Wiadomość wie też, która pozycja kolejki ją wytworzyła — dzięki temu panel
zadań w tle umie pokazać, że wynik podagenta wszedł do rozmowy, a nie zginął.

## budowa/pomocniki/transkrypcja/silnik.py

Moduł nie ma interfejsu wiersza poleceń i nie drukuje na standardowe
wyjście. Rozdzielenie jest celowe: standardowe wyjście pomocnika jest
kanałem danych dla rdzenia i musi zawierać wyłącznie jeden dokument JSON.
Gdyby wykrywanie silnika albo ładowanie modelu miało prawo cokolwiek
dopisać do standardowego wyjścia, rdzeń dostałby dane nieczytelne zamiast
odpowiedzi. Dlatego moduł zawiera wyłącznie funkcje zwracające wartości,
a o tym, co i gdzie wypisać, rozstrzyga moduł transkrypcja.py.

Wykaz przyjmowanych formatów nie obejmuje wszystkiego, co potrafi
przeczytać ffmpeg, bo pomocnik nie ma jak zagwarantować, że ffmpeg jest w
systemie w wersji, która dany kontener otworzy. Odmowa z nazwanym wykazem
mówi więcej niż wyjątek z wnętrza biblioteki.

Sprawdzenie obecności silnika samym zajrzeniem na dysk pozwala uniknąć
ładowania setek megabajtów do pamięci; w trybie sprawdzenia wersji to
różnica między odpowiedzią natychmiastową a kilkusekundowym mieleniem.

Odmowy braku silnika i braku modelu mają zdanie rozpoznawcze dosłowne i
niezmienne, ponieważ rdzeń rozróżnia po treści powodu, który z dwóch
rodzajów niegotowości zaszedł, a Operator po tym samym zdaniu trafia do
właściwego akapitu dokumentacji instalacji; przeredagowanie zepsułoby oba
mechanizmy naraz.

## budowa/server/internal/store/migracja_366_studio_kopie_i_nastawy_pracy.sql

Kopia zapasowa jest osobna od historii wersji, bo ma przetrwać awarię
procesu i awarię zapisu jednocześnie. Wersja leży w repozytorium sesji i
zakłada się ją zapisem, który właśnie mógł się nie udać, więc wersja sama
nie ochroni pracy przed nieudanym zapisem; kopia jest zakładana niezależnie
i dlatego ma własny wiersz. Kolumna udalo_sie jest obowiązkowa, bo wskaźnik
zapisania pokazany mimo nieudanego zapisu jest najgorszym możliwym błędem
tego modułu — operator zamknąłby wtedy okno i stracił pracę. Nieudany
zapis musi być widoczny i nazwany, stąd kolumny udalo_sie i
powod_niepowodzenia przy każdej kopii, nie tylko przy kopiach udanych.
Kolumna zmiany_niezapisane znaczy kopię niosącą pracę, której w dokumencie
jeszcze nie ma; po niej Studio samo zgłasza istnienie niezapisanego
dokumentu z ofertą przywrócenia, zamiast czekać, aż operator się domyśli.

Autozapis odkłada wersje w osobnym szeregu, bo wersje nazwane i kluczowe
zakłada operator ręcznie, a zapisy samoczynne mają być odróżnialne w
wykazie i mieć własną zasadę wygasania — inaczej po godzinie pracy
historia wersji przestaje być historią decyzji i staje się dziennikiem
naciśnięć klawisza. Rozróżnienie niesie kolumna szereg; pojęcie wersji
kluczowej Studio ma już przez odrębne wywołanie oznaczenia wersji i
drugiego nie zakłada.

Nastawy pracy są jedną tabelą na parę okno-dokument, bo skala widoku jest
pamiętana przy dokumencie, a tryb powierzchni przy oknie. Jedna tabela z
nieobowiązkowym dokumentem obsługuje oba przypadki: wiersz bez dokumentu
jest nastawą okna, wiersz z dokumentem nastawą tego dokumentu. Dwie tabele
znaczyłyby dwa odczyty przy każdym otwarciu okna i pytanie, która z nich
wygrywa.

## budowa/desktop/src-tauri/src/aktualizacja/mod.rs

Przedmiotem aktualizacji jest wyłącznie plik powłoki na urządzeniu operatora.
Rdzeń stoi na serwerze wdrożenia i jest utrzymywany tam, więc ten przebieg go
nie dotyka: niczego mu nie podmienia i niczego nie wygasza.

Droga nie idzie przez wtyczkę aktualizacji Tauri: wtyczka nie występuje
w konfiguracji powłoki. Wymaga własnego podpisu minisign, czyli pary kluczy
wydawcy, i narzuca własny kształt pliku wykazu wydań, podczas gdy wykazem
wydań jest istniejący plik serwisu, ten sam, z którego strona pobierania
bierze chronologię. Zamiast tego droga używa HTTPS po plik i sumy SHA-256
z wykazu: suma wiąże plik z wykazem tak samo jak podpis, a wykaz przychodzi
po HTTPS z domeny wydawcy. Przejście na podpis dotknie jednego pliku
odpowiedzialnego za pobranie; reszta przebiegu zostaje bez zmiany.

Zwłoka przed restartem nie jest blokadą: nic nie pyta i niczego nie
wstrzymuje, daje tylko odpowiedzi polecenia czas dolecieć do banera, zanim
okno zniknie — restart natychmiastowy wyglądałby wtedy jak awaria.

Zapora wyłączności chroni przed sytuacją, w której dwa równoległe wywołania
otwierają ten sam plik roboczy obok aplikacji, piszą w niego przeplotem
i każde liczy sumę kontrolną z własnego strumienia, a nie z tego, co
ostatecznie leży na dysku — suma zgadzałaby się wtedy dla pliku, którego
w tej postaci nie ma, po czym oba przebiegi próbowałyby podmienić plik
aplikacji. Straż zwalniająca zaporę działa również przy błędzie i panice,
inaczej jedno niepowodzenie zamykałoby aktualizacje do końca życia procesu.
Funkcja zajmująca wyłączność jest wydzielona osobno, żeby dało się ją
sprawdzić bez stawiania okna, a zwracana straż jest oznaczona jako
wymagająca użycia, ponieważ jej natychmiastowe upuszczenie zwalnia zaporę
i przywraca usterkę, przed którą funkcja stoi.

Restart obrazu przenośnego sięga po zmienną środowiskową wskazującą bieżący
obraz, a nie po ścieżkę bieżącego pliku wykonywalnego, która wskazywałaby
chwilowo podmontowany obraz starego wydania; dzięki temu po podmianie wstaje
wydanie nowe.

## budowa/server/internal/store/migracja_044_roundtable.sql

Debata jest bytem okna, nie sesji. Kontrakt kieruje wszystkie cztery komendy
obszaru `roundtable.*` przez `windowId`, więc identyfikator okna jest tu
kluczem grupującym. Kolumna nie ma więzu obcego do `okno`, bo okno debaty bywa
oknem operacyjnym rejestru 030 (`roundtable.debate-panel`), a nie oknem
komunikacji z migracji 002 — dwa różne byty pod jedną nazwą.

Ten sam kanał może wystąpić dwukrotnie, więc więzu jednoznaczności na parze
(okno, kanał) nie ma: dwaj uczestnicy stoją na tym samym kanale modelu i różnią
się wyłącznie tożsamością — nazwą persony i promptem systemowym. Jednoznaczny
jest identyfikator uczestnika, nic więcej.

Wypowiedź wskazuje uczestnika kodem, nie kluczem. Moderator debaty nie jest
uczestnikiem (nie ma kanału ani persony), a jego interwencja jest wypowiedzią
tury — więz obcy do `debata_uczestnik` odciąłby ją od zapisu. Kod moderatora
jest wartością danych, tak jak kod uczestnika.

Kolumny „kluczowy” w tabeli `debata_uczestnik` nie ma, choć kontrakt ma pole
`RoundtableParticipant.key`. Żadna z czterech komend obszaru nie potrafi go
ustawić (`model.add` pola nie przyjmuje, `ModeratorAction` nie ma wartości
oznaczającej), więc kolumna byłaby miejscem, do którego nic nie pisze.

Pusty kod tury w tabeli `debata_stanowisko` znaczy „stanowisko całej debaty”.
Wartość pusta zamiast NULL jest tu wyborem świadomym: w SQLite dwa NULL-e są
różne, więc więz UNIQUE(okno, tura) na kolumnie dopuszczającej NULL nie
powstrzymałby powielenia stanowiska całej debaty.

## budowa/server/internal/store/migracja_166_silniki_wymiana_tlumaczenia.sql

Profil silnika wskazuje kanały modelu wykazem, więc kanały mają tabelę
dziecka: wykaz w jednej kolumnie nie dałby się złączyć z tabelą kanałów
modelu ani zawęzić zapytaniem o profile używające danego kanału. Kolejność
w wykazie jest znacząca, bo pierwszy kanał jest kanałem pierwszego wyboru,
stąd osobna kolumna kolejności.

Polityka pivota jest jedna na zasięg konfiguracji, a pary języków są jej
dzieckiem. Para mówi, że przekład z języka źródłowego na język celu idzie
przez trzeci, pośredniczący język; język domyślny jest pivotem dla par,
których nikt jawnie nie wymienił.

Pakiet przekazania ma stan, bo jego cykl życia jest realny: złożony,
przekazany, zwrot przyjęty. Przyjęcie zwrotu przestawia ten sam wiersz,
zamiast zakładać drugi — inaczej nie dałoby się jednoznacznie odpowiedzieć,
czy materiał wrócił.

Przebieg pakietowy ma własną kolejkę, nie kolejkę modułu Automations:
uruchomienie przebiegu oddaje identyfikator kolejki, po którym okno
Translate pyta o postęp, a pozycje tej kolejki niosą operacje własne
modułu — przekład, kontrolę jakości, korektę, wydanie. Wpięcie w cudzą
kolejkę związałoby ten moduł z modułem, który daje się wyłączyć
niezależnie.

Most do dokumentu jest wiązaniem dwustronnym: okno Translate wie, z
którego dokumentu wzięło materiał, a wysyłka wyniku wie, dokąd go
odesłać. Bez takiego wiersza wysyłka wyniku musiałaby dostać wskazanie
dokumentu drugi raz, czego kontrakt nie przewiduje.

## budowa/server/internal/store/migracja_121_pamiec_widocznosc_zrodlo.sql

Wszystkie trzy dodatki tej migracji są brakiem bytu, nie brakiem komendy:
komendy `agent.create`, `agent.update` i `extension.install` istnieją
i działają, ale nie miały gdzie zapisać trzech faktów, których wymaga od nich
specyfikacja. Idą w jednej migracji, bo dotyczą jednego zestawu komend;
rozbicie na trzy numery nie dołożyłoby ani jednej informacji.

Definicja eksperta wskazuje, z których poziomów pamięci korzysta domyślnie
(`globalna | projekt | sesja | środowisko`), albo pamięć jest wyłączona
w całości. „Wyłączona" nie jest piątym poziomem i nie wolno jej tak zapisać:
piąta wartość obok czterech pozwala ułożyć wiersz sprzeczny — `sesja`
i `wyłączona` naraz — który rdzeń musiałby rozstrzygać zgadywaniem. Wyłączenie
jest pustym zbiorem poziomów, a jedynym kształtem, który zbiór pusty niesie bez
udawania, jest tabela podrzędna: brak wierszy znaczy brak poziomów i nic
więcej. Kolumna napisowa z listą po przecinku dawałaby to samo, ale bez więzu
CHECK na każdej wartości i bez klucza pilnującego, że poziom się nie powtórzy.

Wartość wyjściowa to komplet czterech poziomów, nie zbiór pusty: stanem
wyjściowym platformy jest pełny dostęp operacyjny, więc ekspert świeżo założony
ma pracować, nie prosić o włączenie pamięci. Wyłączenie jest świadomą decyzją
operatora, więc to ono wymaga czynności. Wartość ta jest wyjściowa, nie
ostateczna: przypisanie eksperta do projektu i do roli nadpisują ją, bo
pierwszeństwo ma zasięg najbardziej szczegółowy. Ta tabela trzyma wyłącznie
wartość z definicji eksperta; nadpisania mają własne poziomy zasięgu.

Widoczność ma dwie wartości (`globalny | projektowy`). Dołożenie obok nich
kolumny `projekt_id` byłoby drugą prawdą o przynależności eksperta do projektu
obok tabeli `przypisanie_agenta_projektu` (migracja 035), którą wypełnia
`workspace.agent.assign`; dwie prawdy rozjechałyby się przy pierwszym
przypisaniu zrobionym drugą drogą. Widoczność mówi więc tylko to, czego tamta
tabela nie mówi.

Rdzeń serwera nie rozróżnia źródła rozszerzenia — obowiązuje wspólny kontrakt
integracji. Ładowanie, wywołanie i użycie operacyjne są dla `danaco`
i `personal` te same, i ta migracja nie zakłada bytu, który by je rozdzielał.
Jedyny wyjątek to stan wyjściowy przy rejestracji: rozszerzenie danaco staje
włączone, bo zestaw wbudowany jest częścią funkcjonalności bazowej; rozszerzenie
personal staje wyłączone, bo operator włącza je świadomie, po przejrzeniu
konfiguracji i zakresu. Rozstrzyga to warstwa instalacji, nie ta tabela —
kolumna niesie sam fakt pochodzenia. Wiersze zastane dostają `personal`, bo
wszystkie powstały wywołaniem `extension.install` przez operatora; wpisanie im
`danaco` byłoby ogłoszeniem, że coś jest częścią pakietu serwera, choć nikt
tego nie dostarczył.
