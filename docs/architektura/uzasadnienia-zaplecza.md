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

## budowa/server/internal/store/migracja_077_obsada_biegu.sql

Obsada od jednego do czterech uczestników z rolami służy modułowi
automatyzacji, a przepływ wieloagentowy potrzebuje dokładnie tego samego
bytu: udział bierze od jednego do czterech modeli z podziałem na role. Drugi
komplet tabel dałby dwa słowniki ról, dwa sufity liczebności i dwa silniki
wybudzeń — dwie prawdy o jednym mechanizmie. Przeszkodą był warunek
obowiązkowego wskazania automatyzacji w dawnej tabeli obsady: bieg
orkiestracji sesyjnej zakłada się w oknie i nie ma wiersza w tabeli
automatyzacji. Ta migracja rozluźnia nośnik: obsada wisi albo na
automatyzacji, albo na biegu orkiestracji, nigdy na obu i nigdy na żadnym —
tym samym wzorcem, którym silnik kolejek jest wspólny dla pętli sesyjnej
i przepływu wieloagentowego.

Górna granica czterech miejsc siedzi w warunku sprawdzającym: koordynator,
dwaj wykonawcy i analityk. Nie jest to sprzeczne z górną granicą piętnastu
podagentów: obsada liczy stanowiska w scenie, a podagenci pracują pod jednym
wykonawcą, w tle, bez własnego gniazda — dwa różne sufity dwóch różnych
bytów.

Ponad obsadę odziedziczoną z modułu automatyzacji każde stanowisko niesie
własny prompt, własne narzędzia i własny profil izolacji, inaczej udział
kilku modeli z podziałem na role znaczyłby tylko kilka nazw zamiast kilku
różnie wyposażonych stanowisk. Obsada niczego nie zabrania: bieg bez obsady
rusza na modelu wskazanym w kroku albo w oknie, a pusty wykaz obsady to
krótsza lista, nie odmowa.

Bieg orkiestracji nie korzysta z tabeli przebiegów automatyzacji, bo tamta
wisi na obowiązkowym wskazaniu zapisanej definicji; bieg orkiestracji
zakłada się w oknie i bywa jednorazowy, a zmuszanie operatora do zapisania
automatyzacji, żeby móc zestawić dwa modele, byłoby zbędną przeszkodą. Stan
oczekujący jest tu wartością pierwszej klasy: bieg zawieszony na sygnał ze
świata to stan tego bytu, nie dopisek.

SQLite nie zdejmuje warunku obowiązkowości ani warunku sprawdzającego, więc
przebudowa dawnej tabeli obsady idzie przez założenie nowej tabeli
i przepisanie wierszy; na dawną tabelę nie wskazywał żaden klucz obcy, więc
przebudowa obejmuje jedną tabelę i kaskada nie ma czego zabrać. Warunek
dokładnie jednego nośnika jest tu sednem poprawności: wiersz bez nośnika to
obsada niczyja, wiersz z dwoma — obsada dwóch biegów naraz; SQLite liczy
wyrażenia logiczne jako zero i jeden, więc suma dwóch testów równa jeden
wyraża dokładnie jeden nośnik bez wyzwalacza. Jednoznaczność miejsca idzie
dwoma indeksami częściowymi, nie jednym złożonym warunkiem unikalności, bo
taki warunek przepuściłby dwa wiersze o tym samym miejscu dla różnych
biegów orkiestracji, ponieważ wartość pusta nie równa się samej sobie; dwa
indeksy częściowe pilnują każdego nośnika osobno i dokładnie.

Trwałą tożsamością podagenta jest pozycja kolejki, a proces modelu jest
wyłącznie sposobem jej wykonania. Podagent nie jest oknem, bo agentów bywa
więcej niż gniazd sceny; nie jest samym procesem, bo wykaz zakończonych
przeżywa restart rdzenia, a proces nie; jego pracę wykonuje silnik kolejek,
który już istnieje. Ponad pozycję kolejki ten wiersz niesie to, czego
pozycja nie wie: kto ją wykonuje w rozumieniu obsady, pod jakim biegiem
biegnie i ile kosztowała — żetony, narzędzia i czas są polami, bo pokazuje
je panel zadań w tle, a struktura kontraktowa podagenta ich nie niesie.

## budowa/desktop/src-tauri/src/wskazanie.rs

Rdzeń stoi na serwerze wdrożenia, a instalator adresu tego serwera nie zna
i znać go nie może: w chwili rozpakowania plików nikt jeszcze nie wie, pod
jaką nazwą stoi rdzeń danego Operatora. Wie to sam Operator i podaje adres
przy pierwszym uruchomieniu okna; drugą drogą wskazania jest zmienna
środowiska, właściwa wykonawcy i jednostce usługi, nie Operatorowi. Produkt
występuje wyłącznie jako hybryda, więc wybór „rdzeń na tym urządzeniu" nie
istnieje: rdzenia na urządzeniu Operatora nie ma, a powłoka nie ma czym go
postawić.

Pole `host` wraca do pola okna po zapisie, żeby Operator poprawiał to, co
napisał, zamiast pisać wskazanie od nowa. Pole `adres` odpowiada postaci
oczekiwanej przez warstwę połączenia klienta (`klient/src/polaczenie/adres-
rdzenia.ts`), która samego adresu nie składa. Warstwa `brak` oznacza pierwsze
uruchomienie i jest dla okna sygnałem, że ma zapytać Operatora o adres;
wskazania pochodzącego ze zmiennej środowiska okno nie nadpisuje zapisem —
ma o tym poinformować, zamiast przyjmować zapis bez skutku.

Kolejność działań w `wskaz` jest wiążąca: najpierw próba połączenia, potem
zapis. Zapis przed próbą utrwaliłby wskazanie, które nie działa, a Operator
zobaczyłby skutek dopiero przy następnym starcie, gdy okno zamieniłoby się
w ekran milczący. `adres` przychodzi z okna w postaci, w jakiej Operator go
napisał: nazwa serwera albo `serwer:port`, z przedrostkiem `http://` albo bez.

Rozbiór wskazania przyjmuje `serwer`, `serwer:17870`, `http://serwer:17870`
oraz adres IPv6 w nawiasach (`[::1]:17870`), bo Operator wpisuje to, co ma
zapisane, a nie to, co wygodne dla rozbioru. Brak podanego portu znaczy port
obowiązujący — ten sam, na którym rdzeń nasłuchuje domyślnie. Więcej niż
jeden dwukropek poza nawiasami jest adresem IPv6 podanym bez nich: portu
w takim zapisie nie ma, cała treść jest hostem.

## budowa/server/internal/store/migracja_076_rodzaj_modulu.sql

Moduły nie są jednym gatunkiem: część to środowiska robocze — okna czatu z
narzędziami i oknami pomocniczymi, część to kompozytory wytwarzające
pozycje używane w innych modułach, jeden jest repozytorium plików, a
czwartym rodzajem jest sekcja konfiguracyjna. Bez kolumny rodzaj tabela
modul nie odróżniała ich niczym, więc podziału nie dało się wyprowadzić z
rdzenia i strona główna musiała trzymać własny wykaz kafli strefy 2.
Kolumna jest opcjonalna: moduł dołożony później, którego rodzaju jeszcze
nie ustalono, ma prawo pozostać z wartością pustą, odróżnialną od każdej z
czterech wartości właściwych; wymuszenie wartości obowiązkowej skutkowałoby
rodzajem wpisanym byle jak przy zakładaniu wiersza. Rodzaj sekcja
konfiguracyjna nie dostaje w tej migracji ani jednego modułu.

Zbiór czterech wartości pilnuje treść tej migracji i warstwa danych, nie
warunek CHECK w schemacie, ponieważ SQLite nie umie dołożyć takiego
warunku przez ALTER TABLE bez przepisania całej tabeli, a przepisanie
pociągnęłoby za sobą klucze obce wskazujące tabelę modul.

Ikony piętnastu modułów przenoszą się z mapy klienta do bazy, ponieważ
kolumna ikony w tabeli modul stała pusta, więc boczna nawigacja musiała
trzymać kopię faktu należącego do bazy; każda z przenoszonych nazw
istnieje już w zestawie ikon klienta.

Moduł Automations traci widoczność w bocznej nawigacji, bo działa ze
strefy 2 strony głównej i nie otwiera własnego okna modułowego. Droga
wejścia zostaje: kafel strefy 2 niesie kod modułu zawsze, a klient otwiera
moduł spoza wykazu nawigacji przez osobną ścieżkę odczytu kompletu modułów
platformy. Wiersze pary moduł-środowisko zostają w macierzy z widocznością
wyłączoną, a nie znikają: macierz ma być kompletem par, a odsłonięcie
modułu jest zmianą bitu, nie wstawieniem wiersza.

ALTER TABLE ADD COLUMN w SQLite przepisuje sam nagłówek schematu, nie
tabelę, więc kolumna dopuszczająca wartość pustą i bez wartości domyślnej
nie dotyka żadnego wiersza istniejących danych. Aktualizacje wartości idą
w tej samej transakcji co zmiana schematu, więc schemat nie zostanie
zastosowany bez treści ani treść bez schematu. Dopasowanie wiersza idzie
po kolumnie kodu, jedynej z warunkiem jednoznaczności, więc każda
aktualizacja trafia w dokładnie jeden wiersz albo w żaden; brak wiersza o
danym kodzie nie jest tu błędem.

## budowa/server/internal/store/migracja_128_studio_dobudowa.sql

Moduł Studio miał w kontrakcie sześć komend i trzy tabele (dokument, wersja,
propozycja zmiany), a opracowanie modułu opisuje osiem okien i dwanaście rodzin
funkcji. Migracja dobudowuje byty, na których te funkcje stoją: komentarze
i adnotacje, zmiany śledzone, gałęzie dokumentu, kolejkę cyfryzacji, operacje
własne, łańcuchy, profile wydania, szablony i odwołania do wersji.

Komentarz redakcyjny i adnotacja różnicy dzielą jedną tabelę, nie dwie: poza
polem zakresu (znaki dla komentarza, numer fragmentu porównania dla adnotacji)
mają te same sześć kolumn — dokument, wątek nadrzędny, autor, treść, stan
rozwiązania, czas. Dwie tabele o sześciu wspólnych kolumnach rozjechałyby się
przy pierwszej poprawce, a zapytanie o wszystko, co ktoś napisał przy
dokumencie, wymagałoby sumy dwóch zapytań.

Autor jest kolumną wersji, nie osobną tabelą: rozdział 3.6 opracowania żąda
przy każdej wersji rozróżnienia zmiany operatora od zmiany modelu, a autor jest
cechą wersji, nie bytem samodzielnym — nie ma stanu, nie ma historii i nie
istnieje bez wersji. Kolumna dopuszcza NULL, więc wersje założone przed tą
migracją pozostają poprawne i czytają się jako autor nieznany, nie jako
operator, którym mogły nie być.

Krok łańcucha i pole szablonu są tekstem (JSON), nie osobną tabelą podrzędną:
nie są wyszukiwane, nie mają własnego cyklu życia i nie wiąże się z nimi nic
z zewnątrz — istnieją wyłącznie jako zawartość swojego rodzica. Tabela
podrzędna dałaby złączenie przy każdym odczycie i nic w zamian.

Warstwa słów rozpoznanych i bloki układu w `pozycja_wczytywania_studio` stoją
tekstem w formacie JSON z tego samego powodu: są odczytem z materiału, nie
bytem samodzielnym, i giną razem z pozycją.

Ingest/OCR Panel jest ósmym oknem modułu z opracowania. Katalog rdzenia znał
dla Studia pięć okien, więc okno zbudowane w kliencie nie miało wiersza,
do którego mogłoby się odwołać; kategoria „narzedzia" kończyła się na pozycji
15, więc szesnasta dokleja się na końcu bez przestawiania wiersza zastanego.

## budowa/server/internal/store/migracja_365_studio_znakowanie.sql

Trzy rzeczy nie wolno pomieszać. Komentarz nie niesie brzmienia — zmiana
śledzona jest już w treści. Propozycja na marginesie niesie brzmienie i nie
jest w treści — operator ją przyjmuje, odrzuca albo poprawia. To trzy różne
czynności i okno ma je odróżniać wyraźnie.

Komentarz ma już swoją tabelę (`komentarz_studio`) i zmiana śledzona swoją
(`zmiana_sledzona_studio`). Tutaj leży to, czego nie miały: wyróżnienie barwą,
znacznik własny operatora i propozycja z brzmieniem. Trzy rodzaje w jednej
tabeli, bo pracuje się nimi tak samo — wszystkie są przypięte do fragmentu,
wszystkie mają autora i stan, i wszystkie wchodzą do jednego wykazu znakowań,
po którym operator skacze i który odhacza.

Wyróżnienie tła jest cechą postaci znaku (`highlightColor` w drzewie postaci)
— i tam musi być, bo inaczej nie wyszłoby przy wydaniu do docx ani do PDF.
Wiersz w tabeli znakowań jest czymś innym: jest pozycją wykazu, po której się
przechodzi, którą się filtruje wedle autora i którą się zdejmuje jednym
poleceniem. Bez tego wiersza wyróżnienie modelu w długim dokumencie ginie:
postać wie, że tło jest żółte, ale nie wie, kto je nadał ani po co.

Znacznik własny („do sprawdzenia", „wymaga źródła", „gotowe") ma nazwę, barwę
i wykaz. Trzymanie barwy przy każdym użyciu znaczyłoby, że zmiana barwy rodzaju
wymaga przejścia wszystkich znakowań — a operator, który zmienia barwę „do
sprawdzenia", zmienia ją dla wszystkich, nie dla jednego miejsca. Rodzaje są
zasięgu operatora, nie dokumentu: znacznik „wymaga źródła" obowiązuje we
wszystkich pismach, a nie zakłada się go od nowa w każdym.

## budowa/server/internal/store/migracja_115_wskaznik_znaczenia.sql

Pełnotekstowy indeks treści biblioteki dopasowuje słowa, ten wskaźnik —
znaczenia: trzyma wektor, czyli ciąg liczb, w którym bliskość odpowiada
bliskości sensu, więc łączy zdania niemające wspólnego wyrazu. Zakresy obu
struktur są różne: indeks pełnotekstowy widzi wyłącznie pliki biblioteki,
wskaźnik obejmuje ponadto historię rozmów i pliki przestrzeni roboczej okna.

Fragment i wektor są bytem wtórnym, odtwarzalnym przebiegiem budowy
wskaźnika; bajty treści leżą w magazynie biblioteki pod sumą kontrolną,
a odwołanie do treści pozostaje jedyną drogą do nich.

Wektor jest polem binarnym, nie tabelą współrzędnych: czyta się go zawsze
w całości i zawsze po to, żeby policzyć jeden iloczyn skalarny, a wiersz na
współrzędną dałby przy dziesięciu tysiącach fragmentów blisko osiem milionów
wierszy, o które nikt nie pyta pojedynczo. Kolumna wymiar pozwala rozpoznać
wiersz uszkodzony, zanim trafi do porównania.

Model stoi w wierszu i w warunku jednoznaczności, bo wektory dwóch modeli
leżą w różnych przestrzeniach i ich iloczyn skalarny nie jest trafnością. Po
zmianie ustawienia modelu w tabeli stoją dwa komplety, a odczyt zawęża się
po modelu; wiersze poprzedniego modelu zostają i czyści je przebieg
przebudowy wskaźnika.

Kluczy obcych do źródła nie ma, bo jeden z trzech zakresów — plik przestrzeni
roboczej — nie ma wiersza w żadnej tabeli, a trzy wzajemnie wykluczające się
kolumny byłyby kształtem gorszym od jednego napisu. Fragmenty po skasowanym
źródle sprząta przebieg budowy wskaźnika, który zna dziś istniejące źródła,
a nie kaskada bazy, która ich nie zna.

Kolumna utworzono niesie milisekundy epoki podane przez wołającego, a nie
wyrażenie czasu bazy danych — bez wskazania źródła fragmentu model cytowałby
bez możliwości sprawdzenia, czego kontrakt tej rodziny komend zakazuje.

## budowa/server/internal/wiedza/pomocnik_przesiewu.py

Osadzenia liczy `pomocnik_osadzen.py`, a ten pomocnik jest przebiegiem
drugim. Różnica nie jest w modelu, tylko w tym, co model dostaje na wejście:
osadzarka widzi pytanie i fragment osobno i sprowadza każde z nich do
wektora, więc spotykają się dopiero jako dwie liczby, a krzyżowy koder czyta
pytanie razem z fragmentem w jednym przebiegu i oddaje jedną ocenę ich
dopasowania. Stąd bierze się kolejność inna niż z kosinusa, i stąd bierze się
koszt: ocen jest tyle, ile kandydatów, a nie jedna na tekst raz na zawsze.

Droga rozmowy jest ta sama co u osadzarki i to jest warunek, nie zbieg
okoliczności: zlecenie przychodzi ścieżką pliku JSON w argumencie, odpowiedź
wraca jednym obiektem JSON na standardowym wyjściu, a diagnostyka biblioteki
idzie na strumień diagnostyczny. Rdzeń woła procesy wyłącznie przez
`zewnetrzne/wolanie.go`, a ta droga nie podaje procesowi standardowego
wejścia.

Katalog modeli bywa dwiema różnymi rzeczami i pomocnik je rozróżnia tak samo
jak osadzarka: pusty jest miejscem, do którego biblioteka dopiero pobierze
wagi, a katalog, w którym wagi już leżą, jest samym modelem — wtedy
biblioteka dostaje go wprost, pobieranie jest wyłączone zmienną
`HF_HUB_OFFLINE`, a druga kopia tego, co stoi na dysku, nie powstaje.
Rozstrzyga obecność pliku `model.safetensors`, bo tylko on jest tu wagami.

Kształtu wag pomocnik nie odczytuje z osobnych deklaracji, inaczej niż
osadzarka, ponieważ krzyżowy koder jest klasyfikatorem pary: `config.json`
niesie komplet, ustrój transformera wraz z głową oceniającą, a warstwy
łączącej tokeny w wektor tu po prostu nie ma, więc nie ma czego zgadywać ani
skąd doczytywać.

Brak jest odpowiedzią, a nie wywróceniem: gdy biblioteki nie ma albo wag nie
da się wczytać, pomocnik oddaje `{"ok": false, "brak": …}` z nazwą braku
i wagą modelu do dociągnięcia, a kod wyjścia zostaje zerowy.

Katalog pobrania jest podawany jawnie z tego samego powodu, dla którego robi
to pomocnik osadzeń: droga wołania procesu nie niesie `HOME` ani `HF_HOME`,
więc bez tego wagi lądowałyby w katalogu pamięci podręcznej zależnym od tego,
gdzie akurat stoi rdzeń, a nie w miejscu wskazanym ustawieniem.

Wykaz tekstów pusty znaczy pytanie, czy koder stoi, a nie polecenie
przesiania niczego: rdzeń pyta o to, zanim przeczyta wskaźnik, żeby odmówić
wcześnie i nazwać brak. Ta sama umowa obowiązuje w pomocniku osadzeń.

## budowa/server/internal/store/migracja_207_rozszerzenia_cykl_zycia.sql

Migracja 070 dała katalogowi jeden wiersz na pozycję i nic poza nim. Wszystko,
co rodzina `extension.*` robi z pozycją w czasie — kolekcjonuje ją, odnotowuje
zmiany, przypina wersję, cofa do wcześniejszej, przyjmuje przesłaną paczkę —
nie miało dotąd gdzie usiąść. Cztery tabele tej migracji są tymi miejscami.

Kolekcja jest nazwanym zestawem, nie etykietą pozycji. `extension.collection.save`
nadsyła `extensionIds` w komplecie przy każdym zapisie, a
`extension.collection.apply` włącza albo wyłącza cały zestaw jednym wywołaniem
— więc związek ma tabelę złącznikową wymienianą „usuń, wstaw od nowa", a nie
kolumnę listy w wierszu kolekcji.

Dziennik cyklu życia jest dziennikiem, nie stanem. `ExtensionHistoryEntry`
niesie czynność, wersję przed i po oraz czas — wiersz na zdarzenie, nigdy
nadpisywany. Wartości kolumny `czynnosc` są wartościami kontraktu
(ExtensionLifecycleAction).

Wersja pozycji ma wiersz, bo inaczej cofnięcie nie ma dokąd wrócić.
`extension.version.rollback` przyjmuje `targetVersion` i ma przywrócić stan
tamtej wersji; pozycja z jedną kolumną `wersja` pamięta wyłącznie tę bieżącą.

Przypięcie jest kolumną pozycji, nie wierszem wersji. `extension.version.pin`
przypina jedną wersję pozycji, a wersja przypięta w dwóch wierszach naraz
byłaby sprzecznością, której nikt by nie wykrył.

Paczka przesłana leży na dysku. `extension.package.upload` przyjmuje bajty
i oddaje `uploadRef`, którym woła się potem instalację; wiersz bez pliku byłby
meldunkiem o przesyłce, której nie ma.

## budowa/server/internal/store/migracja_150_badania_dobudowa.sql

Moduł Research miał w schemacie pięć tabel (źródło, ustalenie, wiązanie, raport,
sekcja, eksport, przestrzeń) obsługujących pięć komend, a kontrakt niesie ich
siedemdziesiąt sześć wobec dziewięciu okien i ośmiu rodzin funkcji opracowania
modułu. Migracja dobudowuje brakujące byty: katalogowanie źródeł, załączniki,
treść do lektury, adnotacje i wypisy, książkę kodów, sprzeczności, weryfikacje,
wątki, ślad prowenancji, pytania badawcze, odkrycia i ich przesiew, monitory,
style cytowania, szablony, wersje raportu, komentarze recenzji, bloki wstawek
oraz szablony i udostępnienia eksportu.

Treść źródła jest kolumną źródła, nie osobną tabelą: tekst wydobyty ze źródła
nie ma tożsamości własnej, nie da się go wskazać bez wskazania źródła, nie ma
historii i ginie razem ze źródłem. Tabela podrzędna dałaby złączenie przy
każdym otwarciu lektury i nic w zamian. Liczba stron i warstwa tekstowa idą
obok tekstu jako jego cechy: „skan bez warstwy tekstowej" jest stanem treści,
nie stanem źródła.

Wynik odkrycia jest bytem trwałym, nie wartością liczoną na żądanie:
`research.discovery.reject` odrzuca wynik po kluczu, a `research.prisma.get`
liczy przesiew — bez wiersza wyniku nie ma czego odrzucić ani czego policzyć,
liczniki PRISMA byłyby liczbami wymyślonymi w chwili odpytania, a odrzucenie
nie przeżyłoby odświeżenia panelu. Wynik zapisuje się w chwili wyszukania,
a odrzucenie jest zmianą jego stanu.

Wersja raportu niesie migawkę sekcji, nie różnicę: `research.report.diff`
porównuje dwie wskazane wersje, także nieprzyległe. Łańcuch różnic wymagałby
odtworzenia stanu przez złożenie wszystkich kroków pośrednich i rozsypałby się
przy pierwszym kroku zgubionym — migawka jest samowystarczalna, dwie wersje
wystarczą do porównania.

Tabela `eksport_raportu_badania` powstaje na nowo (kopia, przepisanie wierszy,
podmiana nazwy), bo więz CHECK na kolumnie `format` wymieniał pięć formatów,
a kontrakt niesie osiem (doszły pptx, xlsx, latex), a SQLite nie zna zmiany
warunku CHECK w miejscu; przeniesienie wierszy zachowuje ślad eksportów już
wykonanych.

## budowa/server/internal/store/migracja_117_profil_asystenta.sql

Wywołanie polecenia głosowego asystenta niosło identyfikator profilu od
pierwszej wersji kontraktu, lecz krok wcześniejszy, wprowadzający zlecenie
i dziennik, nie znał ani profilu, ani kolumny, w której dałoby się go
zapamiętać: adapter przyjmował wskazanie profilu i porzucał je bez śladu,
bo pole żądania było etykietą bez bytu po drugiej stronie.

Profil nie jest ozdobą wykazu okien. Niesie warstwę promptu, która mówi
modelowi, że działa jako narzędzie operatora, nie jako autor samodzielny —
a to ona rozstrzyga, czy model sięgnie po narzędzia platformy, czy
odpowie tekstem we własnym oknie. Bez miejsca na jej zapis asystent
wracałby do bycia zwykłym czatem przy pierwszym uruchomieniu, bo warstwy
nie miałby gdzie postawić; stąd byt trwały, nie etykieta przejściowa.

Profil nie jest zestawem uprawnień blokujących: żadna kolumna tabeli nie
dopuszcza ani nie odmawia czynności, wszystkie są nastawą zasięgu pracy, a
zasięg domyślny jest pełny. Kolumna trybu uprawnień startuje więc z
wartości najszerszej, a nie najwęższej, bo profil świeżo założony ma
pracować, nie prosić o zgodę na każdym kroku.

Pola tabeli pochodzą wyłącznie z tego, co niesie kontrakt i widok okna
asystenta: nazwa, warstwa promptu, kanał modelu, głos syntezy, środowisko
wykonania jako zasięg maszyny, tryb uprawnień jako zasięg pracy i znacznik
profilu domyślnego. Osobne rozmowy przypisane do profilu wymagałyby
kolumny w tabeli okna komunikacji albo w wykazie rozmów, więc ten krok
migracji ich nie dotyka. Nie ma też osobnego przełącznika włączenia
syntezy: pusty głos syntezy już znaczy brak syntezy, a osobna flaga
byłaby drugą prawdą o tym samym fakcie. Czas jest liczbą — milisekundami
epoki, tak jak w krokach migracji wcześniejszych tego obszaru.

## budowa/server/internal/store/migracja_102_wersje_i_archiwum_eksperta.sql

Historia wersji dostaje trwały nośnik, bo dostają go też komendy, które ją
czytają i zapisują (`agent.version.list`, `agent.version.restore`,
`agent.archive`, `agent.restore`, `agent.archive.list`). Bez tej tabeli
historia żyła tylko w zdarzeniach `agent.changed` widzianych w toku sesji,
a `agent.delete` kasował eksperta wraz z całą przeszłością.

Wiersz `agent_wersja` niesie pełną treść tożsamości eksperta w danej wersji,
a nie zapis zmiany — jest migawką, nie różnicą. Różnicę („co się zmieniło")
wylicza warstwa `dane` przy odczycie, porównując sąsiednie migawki. Zapis
różnicowy wymagałby odtwarzania stanu przez złożenie całej historii: jedna
luka w łańcuchu i przywrócenie oddaje eksperta, którego nigdy nie było.

Migawkę zakłada wyzwalacz bazy, nie kod aplikacji, ponieważ tożsamość
eksperta zapisują trzy różne drogi (`Dodaj` i `Aktualizuj`
w `dane/agenci_zapis.go` oraz `ZapiszTozsamosc` i `UstawTrybNakladki`
repozytorium warstw). Historia oparta o jedną z tych dróg byłaby dziurawa,
a każda nowa droga zapisu cicho by ją omijała; wyzwalacz widzi każdą z nich
i każdą przyszłą.

Tabela trzyma jeden wiersz na numer wersji: klucz na parze (agent_id, numer)
wraz z `ON CONFLICT DO UPDATE` sprawia, że zapisy niepodnoszące licznika
(imię własne, favikon, tryb nakładki) uzupełniają migawkę wersji bieżącej
zamiast mnożyć wiersze, więc historia ma tyle pozycji, ile wersji widział
Operator.

Kolumna `autor` niesie napis, nie klucz obcy do konta: wyzwalacz nie ma
dostępu do sesji ani do konta wywołującego, a produkt jest jednoosobowy.
Wartość domyślna to `operator`; przywrócenie wpisuje w to miejsce `restore`,
bo wtedy powód powstania wersji jest znany warstwie `dane`.

Archiwum jest znacznikiem, podobnie jak kosz sesji. Osobna kolumna
`zarchiwizowano_o`, a nie wartość `aktywny`, istnieje dlatego, że ekspert ma
wrócić z archiwum dokładnie w tym stanie czynności, w którym go
archiwizowano, a `aktywny` niesie już inne znaczenie (`enabledOnly`
kontraktu). Definicja i historia zostają nietknięte — archiwizacja nie usuwa
ani jednego wiersza.

## budowa/server/internal/store/migracja_361_studio_postac_dokumentu.sql

Do tej pory rdzeń znał z dokumentu `documentId`, `content` i `title`. Model był
wobec dokumentu ślepy na jego postać: widział tekst, nie widział kroju,
wcięcia, tabeli ani obrazu — więc nie mógł ich ani przeczytać, ani zmienić.
Postać przestaje być stanem klienckim i staje się wierszem w bazie; zapis
dokumentu przenosi ją razem z treścią, zamiast ją gubić.

Drzewo dokumentu (sekcje → bloki → fragmenty o jednolitej postaci znaku,
z tabelami w środku) jest strukturą zagnieżdżoną, którą czyta się i zapisuje
całą: każda czynność na postaci przelicza sąsiedztwo — zmiana wcięcia rusza
łamanie, scalenie komórki rusza szerokości. Rozłożenie tego na wiersze
kazałoby przy każdej czynności składać drzewo z kilkuset wierszy i pilnować
ich kolejności, a żadne zapytanie po pojedynczym fragmencie nie jest do
niczego potrzebne. Dlatego `postac_json` niesie drzewo jednym zapisem.

Styl nazwany i sekcja to co innego i dlatego mają wiersze. Po stylu się pyta:
„ile miejsc go używa", „które style dziedziczą po tym", „usuń styl i przenieś
jego miejsca użycia". Sprawdzian odbioru mierzy wprost, że zmiana stylu
przestawiła wszystkie miejsca użycia — a to jest pytanie do bazy, nie do
drzewa. Sekcja niesie własne nastawy strony i własne nagłówki, po których pyta
podgląd wydruku i wydanie do PDF.

Nazwa stylu jest jego jedynym identyfikatorem w obrębie dokumentu — tak samo
jak w pakiecie biurowym — dlatego warunek UNIQUE stoi na parze (dokument,
nazwa), a nie na osobnym kluczu. `styl_nadrzedny` trzyma nazwą, nie kluczem
wiersza: dziedziczenie ma przetrwać przejęcie arkusza stylów z szablonu, gdzie
klucze wierszy są inne, a nazwy te same.

Nagłówki i stopki sekcji idą jednym zapisem JSON, bo są wykazem najwyżej
trzech pozycji (strony zwykłe, pierwsza strona, strony parzyste) i nikt nie
pyta o nie osobno — pyta o nie sekcja, w całości.

Kolumna `wersja_postaci` rośnie z każdym zapisem. Służy dwóm rzeczom: oknu,
żeby wiedziało, czy trzyma stan świeży, i dziennikowi czynności, żeby cofnięcie
wiedziało, na jakim stanie postaci czynność stała.

## budowa/server/internal/store/migracja_049_badania.sql

Źródło badania nie jest odciskiem odwiedzonej strony przeglądarki: tamten
cykl życia zaczyna się od nawigacji i kończy z kartą, a źródło badania
zaczyna się od decyzji operatora, że dane źródło wchodzi do badania, i niesie
własną, ręczną ocenę wiarygodności, której odcisk strony nie ma i mieć nie
może — to inna prawda o innym momencie. Stąd tabela źródeł badania jest
osobna, bez więzu do tabeli źródeł przeglądarki.

Identyfikator pliku biblioteki jest tekstem bez więzu obcego. Plik biblioteki
żyje w osobnym module, budowanym równolegle — więz obcy do niego wiązałby
kolejność migracji, której ta migracja nie kontroluje, i wymagałby, żeby
wiersz pliku biblioteki istniał już w chwili katalogowania źródła. Kontrakt
dopuszcza źródło bez pliku repozytorium — źródło może być samym adresem albo
notatką — więc kolumna tekstowa dopuszczająca wartość pustą bez sztywnego
więzu jest właściwym wyborem, nie ustępstwem.

Treść ustalenia i sekcji raportu bywa krótką notatką albo długim akapitem
przeniesionym z dokumentu — para kolumn obsługuje oba przypadki bez osobnej
ścieżki dla długiej treści. Inne moduły trwałości używają samego odwołania
do treści, bo tam treść zawsze jest plikiem; tu treść bywa krótkim zdaniem
wpisanym wprost przez operatora, więc krótka ścieżka musi zostać dostępna.

Eksport raportu jest bytem trwałym, nie czynnością bez śladu: identyfikator
pliku biblioteki, ścieżka i rozmiar w bajtach są trzema faktami o wyniku,
które muszą przeżyć samo wywołanie komendy, inaczej operator traci
możliwość odpowiedzieć na pytanie, czy i dokąd dany raport już
wyeksportował, bez ponownego eksportu. Panel raportu pokazuje historię
eksportów obok wersji raportu, więc ślad ma własną tabelę, osobną od tabeli
raportów — dwa eksporty tego samego raportu do różnych formatów to dwa
wiersze, nie nadpisanie jednego.

Przestrzeń badania jest bytem własnym bez okna, tak jak definicja
automatyzacji: żądanie ustawienia przestrzeni nie niesie identyfikatora
okna ani identyfikatora rozbudowywanego bytu, bo zakres i etapy badania są
jedną, bieżącą definicją całego modułu, nie stanem karty — tak jak zapis
definicji automatyzacji nadpisuje jej pola, nie mnoży wierszy.

Powiązanie ustalenia ze źródłem i powiązanie sekcji z ustaleniem mają własne
tabele złącznikowe, po wzorze innych relacji wiele do wielu w tej samej
warstwie trwałości.

## budowa/server/internal/store/migracja_075_mowa.sql

Schemat daje silnikowi mowy dwie rzeczy, których w bazie nie było: sterowanie
i ślad. Silnik jest lokalny i liczy na procesorze — dźwięk nie wychodzi
z maszyny operatora, rozpoznanie robi proces Pythona uruchamiany obok rdzenia
na tej samej maszynie, bez wywołania sieciowego. Dlatego migracja nie niesie
ani kolumny na klucz API, ani na adres usługi, ani na konto dostawcy — nie ma
dostawcy — ani kolumny na wybór urządzenia liczącego, bo silnik chodzi zawsze
na procesorze. Nagranie jest ścieżką, nie bajtami: rdzeń bajtów nie kopiuje,
otwiera plik w miejscu, tak samo jak biblioteka trzyma treść pliku na dysku;
kolumna BLOB urosłaby o rząd wielkości ponad wszystko inne w tym pliku, a kopia
dźwięku obok oryginału byłaby drugą prawdą o tym samym nagraniu. Czas jest
liczbą, a zegar jeden: kolumna `utworzono` niesie milisekundy epoki podane
przez wołającego, nie wyrażeniem `strftime` bazy, bo baza z własnym „teraz"
byłaby drugim zegarem obok zegara rdzenia.

Sterowanie silnikiem idzie czterema wpisami katalogu ustawień, nie stałymi
kodu: ścieżka interpretera, rozmiar modelu, język i katalog pobrania stają się
sterowalne od zaraz przez komendy `config.set`/`config.get` już wpięte
w rdzeń, zamiast czekać na wydanie binarium przy każdej zmianie. Kategoria
`mowa` jest nowa, nie dopisana do `modele`, bo `modele` opisuje kanał modelu
odpowiadającego operatorowi, a silnik mowy niczego nie odpowiada i nie dzieli
z kanałem modelu ani dostawcy, ani nakładu, ani konta. Oś jest wyłącznie
`platform`: rozmiar modelu i ścieżka Pythona są własnością maszyny operatora,
nie zależą od modelu odpowiadającego ani od konta. `mowa_program`
i `mowa_katalog_modeli` opisują jedną instalację na maszynie i stoją na
poziomie `globalny`; `mowa_model` i `mowa_jezyk` opisują pojedyncze zlecenie
i są ustawialne na wszystkich ośmiu poziomach. Wartość pusta ma znaczenie
własne (interpreter na ścieżce systemowej, katalog domyślny biblioteki), nie
jest brakiem danych, dlatego katalog nie niesie tu wpisów `wymagane`.

Dziennik transkrypcji zapisuje odmowy razem z powodzeniami: transkrypcja
nieudana jest zdarzeniem, o które operator zapyta jako pierwsze, a odmowa bez
śladu nie daje się zdiagnozować. Kolumna `stan` niesie trzy wartości —
`gotowa`, `bez_mowy`, `odmowa` — bo nagranie ciszy albo szumu nie jest ani
sukcesem ze znakow=0, ani odmową, skoro niczego nie odmówiono; stan nazywa tę
sytuację wprost. Więz CHECK wiąże stan z kolumną `powod` w obie strony, żeby
w bazie nie dało się zapisać odmowy bez wyjaśnienia. `trwanie_ms` to długość
nagrania, nie czas przetwarzania, który migracja świadomie pomija — byłby
miarą maszyny i chwili, nie faktem o nagraniu. `okno_id` dopuszcza NULL i nie
jest kluczem obcym, bo transkrypcję wołają też ścieżki spoza okna komunikacji
(kolejka wykonawcy asystenta), a odmowa zapisu śladu z powodu nieznanego okna
kasowałaby dowód zdarzenia, które się wydarzyło.

Indeks `idx_transkrypcja_wykaz` biegnie kolumnami dokładnie w porządku
zapytania wykazu z `mowa/dziennik.go` (`okno_id` zawęża, `utworzono, id`
porządkują), więc przy wskazanym oknie SQLite czyta indeks wstecz zamiast
sortować wynik; kolumny `stan` w indeksie nie ma celowo, bo stanęłaby między
zawężeniem a porządkiem i zepsułaby porządek.

## budowa/server/internal/store/migracja_200_apps_produkt.sql

Produkt jest jeden na okno, nie wiele. Kontrakt daje `apps.product.get`
z żądaniem niosącym samo `windowId` i wynikiem o jednym polu `product`,
a `apps.product.save` nie ma pola `productId` — nie ma czym wskazać drugiego
produktu tego samego okna. Bez warunku UNIQUE dwa zapisy z rzędu zakładałyby
dwa wiersze, a odczyt musiałby zgadywać, który z nich jest produktem okna.

Platformy docelowe leżą w jednej kolumnie tekstowej rozdzielonej znakiem nowej
linii, tym samym wzorcem co `architektura_apps.zastrzezenia_walidacji`
(migracja 051). Nikt nie filtruje ani nie sortuje po pojedynczej platformie:
`AppProduct.platforms` wychodzi zawsze w komplecie razem z produktem.

Etap ma własną tabelę, nie kolumnę produktu. `apps.stage.save` zmienia etap po
jego identyfikatorze, `apps.stage.list` zawęża po stanie, a
`AppMilestone.stageIds` wiąże kamień milowy z etapami — każda z tych trzech
dróg wymaga wiersza na etap. Etap należy do okna, nie do produktu: żądanie
`apps.stage.list` niesie `windowId`, a nie identyfikator produktu, więc okno
bez zapisanego produktu ma prawo mieć etapy.

Kolejność etapu jest kolumną, nie porządkiem wstawiania. Tracker etapów rysuje
oś w kolejności zadanej przez operatora (`AppStage.order`), a ta zmienia się
bez zakładania wierszy na nowo.

Związek kamienia milowego z etapami ma tabelę złącznikową, nie kolumnę listy.
`apps.milestone.save` nadsyła `stageIds` w komplecie przy każdym zapisie
(kontrakt nie ma trybu częściowej zmiany), więc zapis wymienia wiersze związku
„usuń, wstaw od nowa" — jedna prawda o krawędzi, nie kopia w dwóch miejscach.
Kolumna `etap_kod` trzyma kod zewnętrzny etapu, a nie więz obcy: kamień milowy
ma prawo wskazywać etap usunięty po zapisie, tak samo jak
`plik_warsztatu_apps.komponent_id` wskazuje komponent zdjęty z architektury.

## budowa/server/internal/store/migracja_015_tozsamosc_modelu.sql

Tożsamość modelu jest zamieniana, nie dołączana do ustawień fabrycznych.
Silnik nakładki niesie trzy warstwy — konstytucja, profil, ekspertyza — oraz
dwa tryby podania promptu systemowego, zastąpienie i dołączenie. Te dwie
tabele niosą sterowanie silnikiem, skąd bierze się treść warstw i który tryb
obowiązuje, i ani jednego zdania promptu: treść wchodzi wyłącznie
z konfiguracji.

Kolumna warstwa przyjmuje dosłownie wartości wyliczenia kontraktu, które
kontrakt sam wiąże ze stałymi silnika nakładki. Kontrakt nie daje dla tego
wyliczenia osobnego słownika przekładu bazy, więc kolumna trzyma wartość
kontraktu wprost, bez drugiego, równoległego nazewnictwa.

Tryb domyślny kategorii i tryb zapisu przyjmują wartości wyliczenia trybu
podania promptu, z wartością domyślną odpowiadającą wartości domyślnej
klucza ustawień tożsamości. Baza nie powtarza reguły składania trybu
nakładki: rozstrzyga ją składacz promptu w warstwie rdzenia, tak samo jak
kolejność poziomów zasięgu rozstrzyga warstwa konfiguracji, a nie schemat.

Oś przyjmuje wartości wyliczenia osi konfiguracji i mówi, dla czego treść
obowiązuje; jest prostopadła do poziomu zasięgu. Byt osi wskazuje konkretny
byt tej osi: identyfikator modelu dla osi modelu, identyfikator konta dla
osi konta. Dla osi platformy byt osi jest pusty — pilnuje tego warunek
sprawdzający, bo pusty i niepusty byt to dwa różne klucze rozstrzygania,
a pomyłka dawałaby dwa wiersze o tym samym znaczeniu.

Byt osi nie jest kluczem obcym: jedna kolumna niosłaby odwołania do dwóch
różnych rodziców — identyfikatora modelu, który w ogóle nie ma własnej
tabeli, bo model jest napisem w kanale modelu oraz w polu okna, i
identyfikatora konta z tabeli kont. Klucz obcy do dwóch tabel naraz nie
istnieje, a wybór jednej z nich zamykałby drogę drugiej osi. Spójność osi
konta sprawdza warstwa danych przy zapisie; brak bytu osi nie wywraca
odczytu, tylko znaczy, że dany wiersz nie obowiązuje.

Kolumna odcisk treści jest skrótem kryptograficznym treści liczonym przez
warstwę wyższą. Służy diagnostyce prowenancji i rozpoznaniu, czy nakładka
się zmieniła — niczego nie dopuszcza i niczego nie blokuje; kolumna pusta
jest dopuszczalna.

Zaczyn katalogu kończy się klauzulą braku działania przy konflikcie kodu,
więc migracja przechodzi także na bazie, w której część wierszy katalogu już
jest. Obowiązkowa jest wyłącznie konstytucja — to jej odcisk niesie panel
prowenancji jako osobny blok. Obowiązkowość jest oznaczeniem dla okna
konfiguracji, nie bramą: brak treści nie wstrzymuje uruchomienia, wraca
wykazem brakujących kategorii obowiązkowych.

## budowa/server/internal/store/migracja_281_agent_wersja_tozsamosc.sql

`AgentVersion` niesie etykietę, opis zmiany, sprawcę i sumę kontrolną, ale nie
niesie treści. Podglądu wersji nie ma z czego złożyć, a przywrócenie było
dotąd jedynym sposobem zobaczenia, co w wersji stało — czyli obejrzenie
historii wymagało jej zmiany. Te dwie kolumny domykają migawkę do kształtu,
którego wymaga kontrakt: widoczność i poziomy pamięci są w `Agent` polami
wymaganymi, więc migawka bez nich nie da się złożyć bez zgadywania.

Czego tu celowo nie ma: warstw promptu. Warstwy zmienia `agent.layer.set`,
która nie podnosi licznika `agent.wersja` i nie zakłada migawki. Kolumna
z warstwami wypełniana wyzwalaczem niosłaby więc warstwy bieżące wpisane do
wersji starej — czyli odpowiedź nieprawdziwą o tym, co w tej wersji stało.
Pole `layers` snapshotu zostaje puste dopóty, dopóki warstwy nie wersjonują
się razem z tożsamością.

Poziomy pamięci wchodzą napisem rozdzielonym przecinkiem, a nie tabelą
podrzędną migawki: migawkę zakłada wyzwalacz, a wyzwalacz nie umie wstawić
wielu wierszy z podzapytania w jednym kroku bez pętli, której SQLite nie ma.
Napis jest tu zapisem migawkowym — nikt go nie zapytuje po wartości, tylko
odczytuje w całości razem z wersją.

## budowa/scripts/pakiet-serwera.sh

Zgodnie z modelem wdrożenia całość platformy stoi na serwerze, a u operatora
zostaje cienkie okno. Ten pakiet jest więc sercem platformy, nie oknem
klienta, i dlatego nie idzie bundlerem Tauri, tylko dpkg-deb po własnym
drzewie katalogów.

Pakiet niesie: rdzeń zbudowany natywnie na Linuksa; serwer narzędzi modelu,
którego rdzeń szuka obok siebie na dysku, stąd oba pliki w jednym katalogu
i stąd katalog roboczy jednostki systemd; pakiet interfejsu, który rdzeń
serwuje klientom; cztery dokumenty produktu — README, instalacja
i konfiguracja, instrukcja, licencja; jednostkę systemd z kontem usługi,
katalogiem danych, restartem i portem nasłuchu; pomocników pythonowych
wołanych przez rdzeń, którego prowizjonowanie bierze stamtąd plik wymagań
środowiska rozpoznawania mowy; oraz skrypt prowizjonowania arsenału na
serwerze.

Skrypt nie stawia arsenału i nie przepisuje jego wykazu. Wykaz zależności
rdzenia stoi w jednym rejestrze deklaracji wraz z deklaracjami w adapterach,
a rdzeń wypisuje go sam poleceniem wykazu zależności. Pole zależności
w pliku kontrolnym pakietu niesie tę część rejestru, którą ma dystrybucja;
resztą zajmuje się skrypt arsenału, wskazany operatorowi przez skrypt
poinstalacyjny.

Wywołanie bez przełącznika buduje wszystko i składa pakiet; przełącznik
pominięcia budowy składa pakiet z tego, co już zbudowane. Zmienna wersji
pakietu domyślnie pochodzi z pliku kontrolnego, a zmienna katalogu wydania
wskazuje, gdzie odłożyć wynik.

## budowa/server/internal/wiedza/pomocnik_osadzen.py

Rdzeń jest w Go, a modele osadzeń są wydawane jako wagi ONNX obsługiwane
bibliotekami Pythona; przepisanie inferencji transformera do Go byłoby drugą
implementacją tej samej rzeczy, rozjeżdżającą się z wagami przy każdym wydaniu
modelu, więc liczenie wektorów jest procesem obok rdzenia — tak samo jak
rozpoznawanie mowy (`mowa/pomocnik.go`). Zlecenie przychodzi ścieżką pliku JSON
w argumencie, odpowiedź wraca jednym obiektem JSON na standardowym wyjściu;
rdzeń woła procesy wyłącznie przez `zewnetrzne/wolanie.go`, a ta droga nie
podaje procesowi standardowego wejścia i nie dziedziczy środowiska, dlatego
diagnostyka biblioteki idzie na strumień diagnostyczny, żeby nie zaśmiecić
odpowiedzi JSON.

Bez dziedziczenia środowiska nie ma `HOME` ani `HF_HOME`, więc biblioteka nie
zna swojego katalogu pamięci podręcznej — rdzeń podaje katalog modeli wprost
w zleceniu, pod katalogiem danych rdzenia, tam gdzie baza i magazyn biblioteki,
więc model pobrany raz zostaje na dysku. Katalog modeli bywa dwiema różnymi
rzeczami: pusty jest miejscem, do którego biblioteka dopiero pobierze wagi,
a katalog, w którym model już leży, jest samym modelem — biblioteka dostaje go
wprost i nie pobiera nic drugi raz. Rozstrzyga obecność pliku ONNX, jedynych
tu wag.

Modelu stojącego biblioteka nie umie opisać sama: jej wykaz obejmuje wyłącznie
wydania, które sama publikuje, a dla każdego innego nie wie, jak złożyć tokeny
w wektor ani czy wynik normalizować. Pomocnik te trzy rzeczy odczytuje
z deklaracji leżących przy wagach (`modules.json`, `config.json` warstwy
łączącej) zamiast zgadywać je — zgadnięcie dałoby wektory bez znaczenia,
a rozpoznać to dałoby się dopiero po jakości wyszukiwania, czyli za późno.

Brak jest odpowiedzią, a nie wywróceniem programu: gdy biblioteki nie ma albo
wag nie da się pobrać, pomocnik oddaje `{"ok": false, "brak": …}` z nazwą
braku i wagą modelu do dociągnięcia, a kod wyjścia zostaje zerowy — rdzeń
zamienia to na odmowę nazywającą brak, bo ślad stosu Pythona nie powiedziałby,
ile waży to, czego nie ma.

Rozmieszczenie plików modelu stojącego w `UKLADY_WAG` odpowiada temu, które
zapisuje `sentence-transformers` przy eksporcie i które ma repozytorium,
z którego biblioteka pobiera własne wagi: model ONNX pod `onnx/model.onnx`,
opis obok niego w korzeniu; eksport pojedynczy zostawia sam plik ONNX
w korzeniu.

## budowa/desktop/src-tauri/src/ustawienia.rs

Wskazanie hosta rdzenia zna dwie warstwy, od słabszej: nastawy zapisane
trwale (`nastawy.rs`), potem zmienna środowiska. Trzeciej warstwy — wartości
domyślnej hosta — nie ma i nie może być: rdzeń stoi na serwerze wdrożenia,
a jego nazwy nie zna ani powłoka, ani instalator. Brak obu warstw znaczy więc
„wskazania nie złożono", a nie „rdzeń stoi tu obok". Nazwy `DANACO_PORT`
i `DANACO_KATALOG_DANYCH` są własnością rdzenia
(`server/internal/konfiguracja/srodowisko.go`); powłoka je wyłącznie czyta.

Zmienna stoi nad plikiem nastaw, ponieważ plik niesie wskazanie Operatora
złożone w oknie i ma przetrwać zamknięcie okna, a zmienna niesie wskazanie
tego, kto stawia proces — wykonawcy przy budowie, jednostki usługi na
serwerze — i musi brać górę, bo inaczej plik z jednej maszyny sterowałby
uruchomieniem na drugiej po skopiowaniu profilu.

Nazwa warstwy wchodzi do odpowiedzi polecenia `wskazanie_rdzenia`, żeby okno
mogło powiedzieć Operatorowi, dlaczego pola nie da się zmienić: wskazanie ze
zmiennej środowiska jest silniejsze od zapisu w oknie.

Warstwa środowiska jest ustalona raz, przy starcie procesu. Warstwa nastaw
żyje dalej — Operator składa wskazanie w oknie już po starcie — więc siedzi
za zamkiem i jest wspólna dla wszystkich kopii ustawień (`Arc`). Bez tego
polecenie `adres_rdzenia` odpowiadałoby starym adresem do końca pracy
procesu, a interfejs łączyłby się nie tam, gdzie Operator wskazał.

Pole `wskazanie` zwraca `None` przy pierwszym uruchomieniu po instalacji:
powłoka nie zgaduje wtedy żadnego adresu, bo każdy zgadnięty byłby adresem
cudzym albo pustym. Zapis wskazania sprawdza kolejność wiążącą — najpierw
próba połączenia, potem zapis — i zwraca zdanie o niepowodzeniu, gdy zapis
się nie udał, bo wskazanie nieutrwalone nie zostaje przyjęte: zniknęłoby
przy następnym starcie i Operator dowiedziałby się o tym dopiero wtedy.
Wskazanie ze zmiennej środowiska nie znika przez ten zapis: `wskazanie()`
pyta zmienną pierwszą, a okno dostaje warstwę w odpowiedzi
(`warstwa_wskazania`) i wie, że zapis nie rozstrzyga.

Odczyt nastaw zapisanych po zamku zatrutym panika daje nastawy puste, nie
panikę samą, ponieważ odczyt nastawy nie jest wart przerwania pracy okna.

## budowa/pomocniki/transkrypcja/transkrypcja.py

Ograniczenia wpisane w konstrukcję pomocnika: nagranie nie opuszcza maszyny,
bo nie ma tu wywołania sieciowego ani importu biblioteki obsługującej sieć;
praca odbywa się wyłącznie na procesorze, z ustawieniami wpisanymi na stałe
w warstwie silnika — wolniej, ale jednakowo na każdej maszynie; tryb
sprawdzenia wersji niczego nie pobiera, tylko patrzy na dysk, żeby sprawdzenie
gotowości nie mogło zająć maszyny na kilka minut ani ściągnąć wag bez
pytania; brak silnika i brak modelu są dwiema osobnymi odpowiedziami
z dwiema osobnymi naprawami, nie jednym wspólnym niepowodzeniem.

Na standardowe wyjście idzie wyłącznie dokument JSON, bo rdzeń parsuje ten
strumień; ostrzeżenia, ślady wyjątków i komunikaty bibliotek idą na
standardowe wyjście błędów, ponieważ pojedynczy dodatkowy wiersz na
standardowym wyjściu psuje odczyt po stronie rdzenia.

Ładowanie modelu przy sprawdzeniu gotowości zostaje jako droga zapasowa dla
nietypowych układów katalogów, gdy pliku modelu nie widać wprost na dysku;
wyjątek z tej próby trafia do pola powodu odpowiedzi, nie na wyjście jako
ślad stosu.

## budowa/server/internal/store/migracja_030_rejestr_okien_operacyjnych.sql

Okno jest bytem katalogu, a obecność okna w module jest relacją między
dwoma bytami. Dzięki rozdzieleniu definicji od przypisania okno wspólne dla
wielu modułów ma jeden wiersz definicji i wiele przypięć, a jego pozycja
może być inna w każdym module — czego jedna kolumna kolejności na definicji
nie potrafi wyrazić.

Poprzednia kolumna wskazująca moduł odchodzi, więc tabela powstaje na nowo
i przejmuje nazwę starej. Katalog jest słownikiem wnoszonym migracją — nie
ma w nim danych operatora, dlatego wiersze zastępuje się kompletem
inwentarza. Kolumna kolejności na definicji porządkuje okno w obrębie jego
kategorii; kolejność w module należy do osobnej macierzy.

## budowa/server/internal/store/migracja_073_agent_warstwy.sql

Ekspert miał `kod`, `nazwa`, `opis` i jedno pole na prompt —
`instrukcje_systemowe`. Portfolio ekspertów potrzebuje trzech rzeczy, których
w schemacie nie było: imienia własnego, bo `nazwa` jest napisem technicznym,
po którym idzie sortowanie wykazu (`idx_agent_nazwa`,
`dane/agenci.go:listaAgentow`); znaku graficznego, bo kolumny na favikon nie
było nigdzie w bazie; oraz warstw promptu, bo `injection/nakladka.go` składa
prompt z warstw (konstytucja, profil, ekspertyza), a ekspert nie miał gdzie
żadnej z nich zapisać, tylko jeden worek tekstu. Kontrakt
(`shared/contract.go`) niósł już struktury `AgentLayer` i `AgentPlugin`, pola
`Agent.displayName`, `Agent.favicon`, `Agent.layers`, `Agent.pluginIds` oraz
cztery komendy `agent.layer.*` i `agent.plugin.*`; bez tej migracji byłyby to
komendy bez miejsca zapisu.

Migracja świadomie nie zakłada katalogu dostępnych wtyczek: kontrakt zna
`agent.plugin.add` z nazwą, źródłem i wersją podanymi wprost, a komendy
przeglądania katalogu rozszerzeń w rodzinie `agent.*` nie ma, więc tabela
katalogu byłaby zapisem bez czytelnika. Nie zakłada też archiwum treści
warstw, bo `agent.wersja` jest licznikiem, nie archiwum, i ta migracja tego
nie zmienia. Nie zakłada wreszcie kolumny na obraz favikonu: `favikon` niesie
odwołanie — napis, który klient umie pokazać (nazwa znaku, ścieżka,
identyfikator zasobu) — a bajtów obrazu baza nie trzyma, tak samo jak
biblioteka trzyma na dysku treść pliku, a w kolumnie wyłącznie odwołanie.

Migracja wyłącznie dodaje: dwie kolumny z wartością domyślną i dwie tabele.
Nie ma tu ani jednego `UPDATE`, ani jednego `DROP`, ani jednej zmiany
istniejącej kolumny — `nazwa`, `instrukcje_systemowe`, `agent_umiejetnosc`,
`agent_konektor` i `agent_uprawnienie` zostają takie, jakie były.

`nazwa` zostaje tym, czym była: napisem technicznym, po którym biegnie
porządek wykazu i wyszukiwanie frazą (`dane/agenci.go`). `imie_wlasne` jest
tym, co widzi Operator, i tylko tym; gdyby imię wpisać w `nazwa`, zmiana
imienia przestawiałaby wykaz i rozjeżdżała wyszukiwanie. Pusty napis znaczy
„nie nadano", nie „błąd": obie kolumny są `NOT NULL DEFAULT ''`, więc każdy
ekspert założony wcześniej dostaje je puste i pracuje dalej bez żadnej
zmiany, a klient pokazuje wtedy `nazwa` i znak zastępczy, bo pola kontraktu
`displayName` i `favicon` są opcjonalne właśnie po to.

Warstw promptu jest dziś trzy, ale liczba trzy nie jest nigdzie przesądzona:
`injection/nakladka.go` dopuszcza warstwę o nazwie spoza katalogu —
`kolejnoscWarstw` zna trzy nazwy, a `pozycjaWarstwy` zwraca dla każdej innej
`len(kolejnoscWarstw) + 1`, czyli ustawia ją na końcu zamiast odrzucić. Trzy
stałe kolumny zamknęłyby ekspertowi drogę, którą silnik nakładki ma otwartą,
a czwarta warstwa wymagałaby zmiany schematu. Klucz główny na parze
(agent_id, warstwa) sprawia, że powtórzone `agent.layer.set` nadpisuje treść
zamiast dokładać drugi wiersz tej samej warstwy — zapis warstwy jest więc
ustaleniem stanu, a nie dopisaniem zdarzenia. `instrukcje_systemowe` zostaje
nietknięte, bo kolumna jest w kontrakcie polem `Agent.systemPrompt` i czyta
ją klient: kasując ją, migracja zabrałaby oknu modułu Agents treść, którą ono
dziś pokazuje, a warstwy są bytem nowym, obok istniejącego pola. Kolumny
`warstwa` i `tryb` przyjmują dosłownie wartości kontraktu
(`shared.IdentityLayer`, `shared.IdentityMode`), bo kontrakt nie daje dla
tych dwóch wyliczeń słownika przekładu bazy — tak samo jak `agent.transport`
i `agent_konektor.rodzaj`; warunki CHECK są jedynym miejscem, w którym ten
katalog stoi po stronie bazy, a warstwa `dane` go nie powtarza.

Konektor jest drogą do usługi: wiersz `agent_konektor` wskazuje most
z katalogu punktów dostępu, z którego rdzeń składa wpisy `mcpServers`
podawane procesowi modelu przełącznikiem `--mcp-config`
(`core/most_okna.go`). Wtyczka jest katalogiem rozszerzeń powłoki — nie ma
adresu, nie ma poświadczenia, nie ma punktu dostępu; ma nazwę, źródło
i wersję, a program dostaje ją przełącznikiem `--plugin-dir`. To dwa różne
przełączniki i dwa różne byty: skille, konektory i pluginy są trzema bytami,
nie dwoma. Wartość `'plugin'` w `agent_konektor.rodzaj` zostaje, bo taki był
wcześniej jedyny sposób zapisania czegokolwiek o wtyczce — tej wartości
migracja nie kasuje i nie przepisuje wierszy, bo kasowanie rodzaju
wywróciłoby `core/adapter_modul_agents_zasoby.go` i zabrałoby treść
konektorom już zapisanym. Kod jest identyfikatorem trwałym tak samo jak
w `agent_konektor`: `id` służy powiązaniom w bazie, a `kod` wychodzi na
zewnątrz jako `AgentPlugin.id` kontraktu, dzięki czemu `agent.plugin.remove`
wskazuje wtyczkę kodem, a nie numerem wiersza, którego kontrakt nie zna.
Tabela nie ma klucza na parze (agent_id, nazwa), bo ekspert może mieć dwie
wtyczki tej samej nazwy w różnych wersjach albo z różnych źródeł, a kontrakt
kasuje wtyczkę po `pluginId`, nie po nazwie — klucz na nazwie odbierałby tę
możliwość bez powodu.

## design/03-marka/emblematy/generator-godel.py

Warianty barwne emblematów z wypaloną barwą (dla rastrów i osadzeń, które nie
potrafią dziedziczyć currentColor) trzymają zasadę bezwzględną: cały emblemat
niesie jedną barwę, kropka nigdy nie odrywa się barwą od obrysu. Tę zasadę
ustala księga znaku razem z arkuszem komponenty.css, a generator jej pilnuje
programowo zamiast zdawać się na ręczne przestrzeganie przy każdym nowym
osadzeniu.

## budowa/desktop/src-tauri/src/aktualizacja/droga.rs

Trzy sytuacje kończą się inaczej przy założeniu wydania: na Linuksie jako
AppImage plik wskazany zmienną `APPIMAGE` zostaje podmieniony i powłoka wstaje
ponownie sama; na Linuksie z pakietu `.deb` podmiana jest niemożliwa, bo
binarka leży w `/usr/bin`, gdzie proces operatora nie ma prawa zapisu, a
powłoka nie podnosi sobie tych praw — ta gałąź zwraca odmowę nazywającą brak
i mówiącą operatorowi, co zrobić samemu (dpkg, strona pobrania); na Windows
wydanie przychodzi jako instalka NSIS, którą powłoka uruchamia jako osobny
program i schodzi jej z drogi.

Podmiana idzie przez `rename`, nie przez zapis w miejsce: na Linuksie nie
wolno nadpisać pliku wykonywalnego, który właśnie biegnie (jądro zwraca
`ETXTBSY`), a `rename` w obrębie tego samego katalogu tego zakazu nie łamie —
podstawia nowy i-węzeł pod starą nazwę, proces już uruchomiony dopracowuje na
starym, który znika dopiero po jego zakończeniu. Plik roboczy pobiera się więc
obok celu, nie do katalogu tymczasowego, bo `rename` działa wyłącznie w obrębie
jednego systemu plików.

Brak zmiennej `APPIMAGE` na Linuksie znaczy, że aplikacja pochodzi z pakietu
systemowego (`.deb` → `/usr/bin`) albo z budowy deweloperskiej
(`target/release`) — w obu wypadkach podmiana z wnętrza aplikacji jest
niewłaściwa, bo pakietem zarządza `dpkg`, a budowę deweloperską nadpisuje
`cargo`.

Sprawdzenie zawartości pliku roboczego zadaje dwa pytania o brak, nigdy
o zgodę: czy jest w nim cokolwiek i czy wygląda na plik wykonywalny Linuksa.
Zgodna suma SHA-256 nie odpowiada na pytanie „czy to w ogóle aplikacja" —
wykaz wydań może wskazywać wydanie na inny system albo plik pusty, a suma
będzie się zgadzać co do znaku. Przy zakładaniu: ostatnie spojrzenie na
zawartość dzieje się tuż przed podmianą, bo `zaloz` podmienia plik aplikacji
nieodwracalnie i nie wierzy, że ktoś wcześniej sprawdził właściwą rzecz. Bit
wykonywalny nadaje się dopiero po sprawdzeniu sumy, żeby plik pobrany,
a jeszcze niesprawdzony, nie miał prawa być uruchamialny nawet przez pomyłkę.
Instalka NSIS budowana przez Tauri przyjmuje `/S` (przebieg cichy); powłoka
musi zejść z drogi, bo instalator nie podmieni pliku trzymanego przez biegnący
proces, dlatego `restartuje_powloka` jest tu fałszem — aplikację z powrotem
stawia instalator, nie powłoka.

Moduł testów jest cały pod `#[cfg(not(windows))]`, bo sprawdza wyłącznie
wariant `AppImage` i funkcję `sprawdz_zawartosc`; gałąź windowsowa nie ma tu
testów i mieć ich nie będzie, bo sprawdzenie jej wymaga uruchomienia
instalatora NSIS na Windowsie, a test to udający byłby atrapą. `rozpoznaj()`
czyta zmienną środowiska, wspólną dla całego procesu testowego, więc dostęp do
niej zamyka zamek, bez którego testy mrugałyby przy równoległym biegu.

## budowa/desktop/src-tauri/src/aktualizacja/pobranie.rs

To jest jedyne miejsce w powłoce, w którym dane z sieci stają się plikiem
uruchamianym na maszynie użytkownika. Suma kontrolna jest tu warunkiem, nie
diagnostyką: plik o niezgodnej sumie zostaje skasowany, zanim funkcja wróci.
Kolejność kroków jest wiążąca: pobierz do pliku roboczego obok celu (nie pod
nazwą celu), policz SHA-256 z tego, co leży na dysku, porównaj z sumą żądaną,
dopiero potem wolno zakładać plik (robi to `droga.rs`). Suma liczona jest
w locie z tego samego strumienia, który trafia na dysk — liczona z osobnego
bufora w pamięci sprawdzałaby co innego niż zapisany plik.

Łańcuch zaufania ma trzy ogniwa: kanał pobrań wkompilowany w powłokę
(`ADRES_KANALU`) ogranicza, skąd plik w ogóle może przyjść; HTTPS chroni
wykaz wydań, z którego pochodzi suma; a suma chroni pobrany plik. Adres inny
niż `https://` jest odmawiany, bo po zwykłym HTTP pośrednik podmienia plik
i sumę naraz; adres spoza kanału jest odmawiany, bo adres i sumę podaje
powłoce ta sama strona `budowa/witryna/wydania.json`, więc bez kanału nic nie
wiąże ich z wydawcą — strona podstawiona wskazałaby własny plik wraz z jego
poprawną sumą. Kanał jest częścią powłoki z tego samego powodu, z którego
suma jest częścią wykazu: ogniwo zaufania nie może pochodzić od strony, którą
wiąże. Porównanie przedrostka jest dosłowne i kończy się ukośnikiem, żeby host
`pobierz.danaco-group.pl.obcy-serwer` nie przeszedł jako pasujący przedrostek.
Wartość `ADRES_KANALU` musi pozostać kopią `kanal.adres` z wykazu wydań —
pilnuje tego sprawdzian `kanal_powloki_zgadza_sie_z_wykazem_wydan`, który
czyta wykaz przy kompilacji, żeby przeniesienie kanału w jednym miejscu bez
drugiego nie rozjechało obu wartości po cichu.

Pułap wielkości pliku wydania (512 MiB) zabezpiecza przed zapełnieniem dysku
przez odpowiedź bez końca z serwera zepsutego albo podstawionego; dobrany
z zapasem względem wielkości pakietu powłoki, która idzie w dziesiątki
megabajtów.

`sprawdz_zadanie` sprawdza samo żądanie, zanim cokolwiek poleci przez sieć
i zanim powstanie plik roboczy — dlatego jest wydzielone z
`pobierz_i_sprawdz`: to jedyne odmowy padające przed pierwszym bajtem
z gniazda, więc test może je wywołać bez dostępu do sieci. Zwraca sumę
sprowadzoną do małych liter, bo `format!("{:x}")` daje małe litery, a wykaz
wydań pisany ręcznie potrafi mieć wielkie.

`ocen_pobrane` sprawdza kolejność pytań jako warunek poprawności: najpierw
dolna granica wielkości, dopiero potem suma. Plik pusty ma poprawną,
64-znakową sumę SHA-256, więc wykaz podający właśnie ją — albo serwer
oddający 200 z pustym ciałem — przeszedłby przez sam warunek sumy,
a `droga.rs` podstawiłoby zero bajtów w miejsce aplikacji i zaplanowało
restart, po którym nie ma z czego wrócić. Kasowanie pliku roboczego zostaje po
stronie wołającego (`posprzataj`), bo to on wie, jaki plik utworzył; każde
wyjście błędem po utworzeniu pliku przechodzi przez `posprzataj`, żeby plik
pobrany, ale niesprawdzony albo wprost niezgodny, nie został na dysku obok
właściwej aplikacji, wyglądając na jej część.

Testy modułu biegną bez sieci, bo `sprawdz_zadanie` pada przed pierwszym
bajtem z gniazda; osią jest przywiązanie do kanału pobrań — adres spoza
`ADRES_KANALU` ma zostać odrzucony, zanim cokolwiek poleci przez sieć.
Wykaz wydań nie jest w testach przepisany, tylko wczytany przy kompilacji
przez `include_str!`, żeby nie stał się drugą kopią tych samych wartości,
rozjeżdżającą się cicho z wykazem tak samo jak sam `ADRES_KANALU` mógłby się
rozjechać bez sprawdzianu wiążącego.

## budowa/server/internal/store/migracja_377_katalog_powiadomien.sql

Sekcja Powiadomień nie wymaga nowej rodziny kontraktu, ponieważ jest macierzą
nastaw i idzie tą samą drogą co każde inne ustawienie platformy: `config.get`,
`config.set`, `settings.definition.list` oraz zdarzenie `config.changed`.
Kody klas są kodami z modelu danych (`powiadomienie.klasa`: zakonczenie,
decyzja, blad, wzmianka, termin, automatyka, system) — jeden zapis na całą
platformę, bo drugi zestaw nazw rozjechałby ustawienie z wierszem, którego
dotyczy.

Kanał centrum nie jest tu wyborem i nie ma dla niego kolumny: centrum
powiadomień jest kanałem podstawowym każdej klasy, a każde zdarzenie objęte
ustawieniem trafia do rejestru centrum niezależnie od pozostałych kanałów.
Nastawa wybiera więc kanały dodatkowe; centrum stoi zawsze, dopóki klasa jest
czynna.

Zasięgi zapisu obejmują warstwę globalną, środowisko i kartę sesji, spełniając
jednocześnie każdy z zapisów źródłowych, które nazywają część tych warstw dla
tej sekcji — zawężenie do mniejszego zbioru wybierałoby, który z zapisów jest
ważniejszy, a to nie jest rozstrzygnięcie tej migracji. Żadna pozycja nie
wymaga restartu: nastawy czyta się drogą `config.get` przy otwarciu sekcji,
a zmiany dolatują zdarzeniem `config.changed`.

Wyłączenie przełącznika głównego wygasza wszystkie klasy naraz, zachowując
ich ustawienia: jest osobnym kluczem, a nie zapisem wartości fałszywej do
siedmiu kluczy klas, bo to skasowałoby wybór Operatora, a ponowne włączenie
przywróciłoby stan domyślny zamiast poprzedniego.

## budowa/server/internal/store/migracja_003_kolejki.sql
Nagłówek tabeli dziennika akcji kolejki rozszerzono o cel przechowywanych
wpisów, ponieważ krótkie zdanie nie mieściło się w wymaganym przedziale
długości nagłówka, mimo że treściowo było wystarczające.

## budowa/przygotuj-drzewo.sh
Nagłówek skryptu zwięźle łączy trzy fakty osobnego akapitu źródłowego —
przyczynę dowiązania, sposób działania i bezpieczeństwo wielokrotnego
uruchamiania — w jedno zdanie mieszczące się w wymaganej długości nagłówka.

## budowa/scripts/pokaz.sh
Nagłówek skryptu połączono w jedno zdanie mieszczące cel i sposób działania;
przykłady wywołania usunięto z nagłówka, ponieważ powtarzają treść widoczną
w komunikacie końcowym skryptu i w opisie flagi zatrzymania.
## budowa/server/internal/store/migracja_099_wyzwalacz_poczty.sql
Odbiór listu uruchamia automatyki dwiema drogami jednocześnie: startuje bieg każdej czynnej automatyki z wyzwalaczem rodzaju mail od początku oraz doręcza sygnał mail z kodem skrzynki silnikowi wybudzeń, wznawiając bieg zawieszony krokiem oczekiwania. Wyrażenie wyzwalacza mail wskazuje skrzynkę: kod skrzynki pocztowej albo gwiazdkę oznaczającą każdą skrzynkę. Ponieważ zapis harmonogramu pomija wyzwalacze z pustym wyrażeniem, kontrakt „każda skrzynka" zapisuje się jako gwiazdka, nie jako pole puste. Na tabelę wyzwalacz_automatyki nie wskazuje żaden klucz obcy, więc przebudowa obejmuje wyłącznie tę jedną tabelę.

## budowa/server/internal/zewnetrzne/.dowod-wpiecia/uruchom.sh
Nagłówek połączono w jedno zdanie opisujące cel dowodu i sposób weryfikacji
nietkniętych oryginałów; wskazanie katalogu uruchomienia usunięto z nagłówka,
ponieważ jest wywołaniem, nie opisem stanu.

## budowa/server/internal/zewnetrzne/.dowod-wpiecia/zastosuj.py
Dodano nagłówek w postaci komentarza, ponieważ narzędzie redakcji nie
rozpoznaje ciągu dokumentacyjnego modułu jako nagłówka; ciąg dokumentacyjny
pozostał nietknięty jako treść kodu.
## budowa/server/internal/store/migracja_100_zywotnosc_podagentow.sql
Tabela podagent nosiła pełny cykl stanów, lecz wiersz nie wiedział, które uruchomienie rdzenia prowadzi jego pracę. Rdzeń przerwany w trakcie pracy podagenta zostawiał wiersz w stanie running na zawsze, ponieważ proces modelu ginął razem z rdzeniem, a panel zadań w tle po restarcie pokazywał pracę, której nikt nie wykonuje. Podagent należy do uruchomienia rdzenia, nie do pliku bazy: trwała jest jego tożsamość i wynik, nietrwała jest jego praca, która żyje w procesie modelu i ginie z rdzeniem. Kolumna uruchomienie_rdzenia zapisuje tę granicę wprost — wiersz niezakończony ze znacznikiem innym niż znacznik rdzenia, który właśnie wstał, jest podagentem osieroconym i rdzeń zamyka go przy starcie. Wybrano znacznik uruchomienia zamiast progu czasowego, ponieważ sam czas daje tylko domysł o awarii, a znacznik rozstrzyga bez domysłu; oznaka życia pozostaje jako zapis ostatniego dotknięcia wiersza, nie jako kryterium osierocenia. Krok użył ALTER TABLE zamiast przebudowy tabeli, ponieważ podagent jest celem indeksu częściowego i nośnikiem wyniku, żaden klucz obcy na niego nie wskazuje i nic w nim nie ubywa, więc przebudowa nie miałaby czego naprawić. Krok jest dokładający — cofnięcie to usunięcie nowego indeksu i przebudowa tabeli bez trzech dodanych kolumn, ponieważ SQLite zdejmuje kolumnę wyłącznie przez przepisanie tabeli. Słownik powodu zakończenia jest polski, ponieważ to nazwa nadana w warstwie bazy, a nie wzięta z kontraktu, który powodu nie niesie.

## design/zasoby/marka/otwarcie/zrodlo-blender.py
Jedyny nagłówek rozpoznawany w pliku — poprzedzający funkcję łagodzenia
ruchu — rozwinięto do pełnego zdania o rozkładzie narastania animacji, bo
skrócony zapis nie mieścił się w wymaganej długości nagłówka.

## design/zbuduj-licencje.py
Za znacznikiem interpretera dopisano zdanie streszczające zadanie skryptu,
ponieważ sam znacznik stanowi osobny, zbyt krótki nagłówek. Komentarz o polu
przewijanym zestawień rozwinięto o przyczynę i miejsce zastosowania.
## budowa/server/internal/store/migracja_101_doradca.sql
Doradca to kanał modelu, którego agent woła po radę, a nie po wykonanie pracy. Pojęcie stoi na trzech własnościach, które ta migracja utrwala. Rada jest jawna: kolumny pytania, kanału i modelu doradcy oraz rady stoją obok siebie w jednym wierszu, więc nie da się odczytać rady bez odczytania, kto ją dał i o co był pytany; gdyby treść rady lądowała w tabeli wiadomości agenta, nie dałoby się po czasie odróżnić rady cudzej od odpowiedzi własnej, dlatego tabela istnieje osobno. Konsultacja ląduje w prowenancji: kolumna prowenancja niesie ten sam napis JSON, który poszedł strumieniem jako fragment provenance, kopiowany, a nie odtwarzany powtórnym złożeniem, aby ślad w bazie nie rozjechał się ze śladem w oknie. Konsultacja to nie delegacja: w tabeli nie ma kolumny stanu pracy, kolumny wyniku ani wiązania z pozycją kolejki, ponieważ rada niczego nie zleca i odpowiedzialność za wynik zostaje przy agencie pytającym. To, kto może być doradcą, rozstrzygają dane rejestru kanałów, nie ta migracja — nowy doradca to nowy wiersz istniejącej tabeli kanałów, nie nowa kolumna. Odmowy zapisywane są razem z radami, ponieważ odpowiedź na pytanie, dlaczego nie zapytano doradcy, jest możliwa tylko wtedy, gdy nieodbyta konsultacja zostawia wiersz; ograniczenie CHECK wiąże stan z treścią w obie strony, więc wiersz rady bez rady jest niemożliwy w bazie. Kolumna utworzono niesie milisekundy epoki podane przez wołającego, tak samo jak w migracjach dziennika transkrypcji i katalogu rozszerzeń, ponieważ drugi zegar rozjeżdżałby chwilę rady z chwilą zapisu. Kolumna okno_id jest napisem, nie kluczem obcym, ponieważ konsultację prowadzi też praca spoza okna rozmowy, a odmowa zapisu śladu z powodu nieznanego okna kasowałaby dowód zdarzenia, które się wydarzyło.
## budowa/server/internal/store/migracja_103_zespoly.sql
Kontrakt niesie cztery komendy oraz zdarzenie dla zespołów, a schemat nie miał niczego, co mogłoby je zasilić. Zespół nie jest ani projektem przestrzeni roboczej, ani składem debaty Roundtable: tamten skład żyje w jednej debacie i ginie razem z nią, a zespół jest bytem trwałym Operatora, wybieranym wielokrotnie i kopiowanym. Skład wskazuje kod eksperta, nie numer wiersza, ponieważ numer wiersza żyje wyłącznie wewnątrz bazy i po przeniesieniu bazy przestałby się zgadzać, a zespół ma przeżyć eksport i import biblioteki. Kolumna agent_kod celowo nie ma klucza obcego do tabeli agent: klucz obcy z kaskadą kasowałby wiersz składu przy usunięciu eksperta, a bez kaskady odmawiałby usunięcia eksperta należącego do zespołu; zamiast tego ślad po ekspercie zostaje, a odczyt składu łączy kod z tabelą agenta z pominięciem archiwum, więc ekspert wrócony z archiwum wraca do składu sam. Kolejność w składzie nadaje Operator, zgodnie z kontraktem, dlatego porządek alfabetyczny byłby zmyśleniem innej prawdy niż zapisana. Czas jest liczony w milisekundach epoki, nie w formacie ISO 8601, ponieważ taki kształt niosą pola czasu kontraktu, a przekład formatu w kodzie rdzenia byłby drugą prawdą o tej samej chwili.

## budowa/server/internal/store/migracja_104_historia_retencja.sql
Rodzina komend historii dostaje tu jedną nową tabelę, zasadę przechowywania; drugiej tabeli na historię rozmowy okna nie ma, ponieważ historia rozmowy okna już istnieje w tabeli wiadomości i w tabeli nietekstowych bloków wypowiedzi, a założenie trzeciej tabeli na to samo byłoby drugą prawdą o jednej wypowiedzi. Jeden wiersz zasady odpowiada jednemu zakresowi, a zakres rozstrzyga się od najwęższego: okno wygrywa z sesją, sesja z zasięgiem globalnym, tą samą kolejnością co w warstwowej konfiguracji, i zasady się nie sumują, ponieważ sumowanie dawałoby wynik, którego Operator nie umiałby przewidzieć z żadnego pojedynczego ekranu. Kolumna zakresu niesie wartość kontraktu w oryginalnym języku, nie przekład, ponieważ własny słownik polski byłby drugim źródłem odwzorowania. Kolumna zakres_kod jest identyfikatorem kontraktowym, nie kluczem obcym, ponieważ zasada ma przeżyć okno, którego dotyczy — okno znika wraz z sesją, a zasada zostaje wierszem bez skutku zamiast znikać kaskadą po cichu. Oba progi retencji są niewymagane i mogą być puste naraz: brak progu znaczy brak ograniczenia, a nie zero, więc wartość zero jest odrzucana warunkiem CHECK, bo ograniczenie do zera pozycji i brak ograniczenia to dwie różne rzeczy. Tabela wchodzi pusta, więc żadne okno nie ma zasady na starcie, a retencja zaczyna działać dopiero w chwili, w której Operator ją ustawi.

## budowa/server/internal/store/migracja_105_sekcje_paneli.sql
Panel okna operacyjnego pokazuje kilka sekcji jedna pod drugą; Operator przestawia je kolejnością, zwija te, których w danej pracy nie czyta, a niepotrzebne zdejmuje z widoku, i ta tabela jest nośnikiem tego układu, bez którego układ ginąłby z zamknięciem okna. Kluczem układu jest para okna i panelu, nie sam panel, ponieważ ten sam panel stoi w wielu oknach naraz i w każdym z nich Operator układa go pod inną pracę. Kolumna okno_id jest napisem, nie kluczem obcym, z tego samego powodu co w dzienniku transkrypcji i w konsultacji doradcy: wiersz okna komunikacji powstaje leniwie, przy pierwszej wiadomości, a układ panelu Operator przestawia natychmiast po otwarciu okna, więc klucz obcy odmawiałby zapisu układu dokładnie wtedy, gdy Operator go ustawia. Kolejność sekcji liczy się od jeden i jest ciągła, co wymusza ograniczenie CHECK na dolnej granicy, ponieważ zero i liczby ujemne mają być w bazie niemożliwe, nie tylko niezalecane w kodzie. Zwinięcie i zdjęcie z widoku to dwa osobne stany, bo są to dwie różne rzeczy: zwinięcie oznacza sekcję obecną na widoku, lecz zawiniętą do nagłówka, a zdjęcie oznacza brak sekcji na widoku w całości; jedna kolumna trójstanowa zlewałaby oba stany w jeden i odbierała Operatorowi możliwość zdjęcia sekcji rozwiniętej. Tabela celowo nie niesie nazwy ani treści sekcji, ponieważ katalog sekcji panelu należy do warstwy widoku, a wpisanie tu nazwy założyłoby drugi katalog obok tego, którym posługuje się klient.

## budowa/server/internal/store/migracja_106_pulap_kosztu.sql
Rdzeń dotąd umiał o koszcie tylko rzeczy wsteczne: czytał koszt tury, która już się odbyła, i wykrywał granicę, o którą tura już uderzyła, lecz nie umiał powiedzieć programowi wywołującemu model, ile wolno wydać, mimo że program przyjmuje przełącznik ograniczający budżet wywołania. Wykrycie po fakcie mówi Operatorowi, że pieniądze już wydano, a pułap sprawia, że nie zostaną wydane — pierwsze jest meldunkiem, drugie nastawą zapobiegawczą. Pułap jest nastawą okna, nie stałą produktu: wchodzi katalogiem ustawień i rozstrzyga się tymi samymi poziomami zasięgu co nakład rozumowania i kanał zapasowy, więc Operator ustawia go globalnie, na projekcie, na karcie sesji albo na jednym oknie, a okno wygrywa. Wartość domyślna zero oznacza brak pułapu, dzięki czemu bazy zastane po tej migracji zachowują się jak przed nią, a pułap zaczyna działać dopiero w chwili, w której Operator wpisze liczbę dodatnią, wyłącznie na tym poziomie zasięgu, na którym ją wpisał. Nastawa jest liczbą zmiennoprzecinkową w dolarach, nie liczbą całkowitą w groszach, ponieważ przełącznik programu i pole kosztu zdarzenia wyniku niosą kwotę w dolarach, a przeliczanie w dwie strony na granicy z programem byłoby miejscem, w którym jednostka może się zgubić. Krok jest dokładający: cofnięcie polega na usunięciu wiersza katalogu, po którym przypięcia do poziomów i osi znikają kaskadą, a każde wstawienie kończy się klauzulą pomijającą konflikt, więc krok przechodzi także na bazie, w której wiersz już istnieje.

## budowa/server/internal/store/migracja_224_workspace_wersje_instrukcji.sql
Instrukcja obowiązująca leży w tabeli ustawienie jako ustawienie ośmiu poziomów zasięgu, bez drugiego porządku poziomów. Tabela wersja_instrukcji_projektu przechowuje historię: każdy zapis odkłada wersję, a panel wersji czyta wyłącznie ten wykaz. Przywrócenie zakłada nową wersję o treści wskazanej wersji i zapisuje w kolumnie przywrocono_z, skąd treść wzięto. Historia nie jest nadpisywana: wersja raz zapisana pozostaje w wykazie na stałe, ponieważ przywrócenie jest czynnością odnotowaną, a nie cofnięciem czasu.

## budowa/server/internal/store/migracja_225_workspace_kalendarz.sql
Kalendarz projektu łączy dwa źródła: zadania z terminem, przechowywane w tabeli zadanie_projektu, oraz wydarzenia wciągnięte plikiem iCal bez zakładania zadań. Tabela pozycja_kalendarza_projektu trzyma to drugie źródło; bez niej wciągnięcie z opcją pomijającą zakładanie zadań nie miałoby gdzie osiąść, a kalendarz pozostawałby pusty mimo udanego wczytania pliku. Kolumna uid_zewnetrzny niesie identyfikator wydarzenia nadany w źródle, a warunek jednoznaczności na parze projekt–identyfikator sprawia, że powtórne wciągnięcie tego samego pliku aktualizuje istniejące wiersze zamiast tworzyć ich kopie.

## budowa/server/internal/store/migracja_226_workspace_stan_wstrzymany.sql
Zestaw stanów projektu w kontrakcie obejmuje trzy wartości: active, paused i archived, podczas gdy poprzedni schemat dopuszczał tylko dwie. W praktyce oznaczało to, że polecenie ustawiające stan projektu na wstrzymany kończyło się odmową wynikającą z warunku kolumny, a nie z zamierzonego działania. SQLite nie pozwala zmienić warunku CHECK istniejącej kolumny, dlatego migracja odtwarza tabelę projekt od nowa z rozszerzonym warunkiem i przenosi do niej istniejące wiersze; więzy obce wskazujące na projekt(id) pozostają spójne, ponieważ klucze główne są przepisywane bez zmiany wartości.

## budowa/server/internal/store/migracja_230_design_kolekcje.sql
Kolekcja jest bytem odrębnym od etykiety, choć obie grupują zasoby. Etykieta zasobu jest wolnym słowem bez własnej tożsamości: nie ma nazwy własnej, opisu ani kolejności, a dwa zasoby niosące to samo słowo należą do siebie wyłącznie przez nie. Kolekcja ma nazwę, opis i kolejność, więc ma własny wiersz i własny identyfikator zewnętrzny; inaczej zmiana nazwy kolekcji wymagałaby przepisania etykiety w każdym zasobie z osobna. Przypisanie zasobu do kolekcji jest dokładką albo odjęciem, nigdy zastąpieniem, dlatego para kolekcja–zasób stanowi klucz główny bez osobnego surogatu: powtórzone dołożenie tego samego zasobu nie zakłada drugiego wiersza. Zasób jest wskazywany identyfikatorem zewnętrznym tekstowym, a nie więzem obcym, dzięki czemu kolekcja przeżywa usunięcie zasobu, który się w niej znalazł.

## budowa/server/internal/store/migracja_231_design_szablony_promptu.sql
Szablon promptu ma własną tabelę, a nie flagę w tabeli prompt_design. Wydany prompt jest zdarzeniem historii: powstał raz, przeszedł kanałem, wygenerował zasoby i pozostaje zapisem tego, co się stało. Szablon jest bytem żywym, który operator nadpisuje, przemianowuje i wykorzystuje wielokrotnie. Wspólna tabela wymagałaby odróżniania jednego od drugiego kolumną filtrującą każdy odczyt historii, a nadpisanie szablonu przepisywałoby wiersz, na który wskazują powstałe z niego zasoby, czyli zmieniałoby prompt źródłowy. Pola promptu są przepisane z prompt_design, ponieważ opisują ten sam kształt kontraktu; odwołanie do wiersza promptu zamiast kopii pól związałoby szablon z jednym wydaniem, a usunięcie tamtego wydania zabrałoby operatorowi szablon.

## budowa/server/internal/store/migracja_246_terminal_ksiazka_hostow.sql
Do tej migracji książka adresowa hostów istniała wyłącznie w widoku klienta: wpisy znikały przy odświeżeniu strony, a jedyną drogą ich zachowania poza jedno posiedzenie był wywóz do pliku. Wpis hosta jest jednak nastawą operatora, nie stanem widoku, więc ma przetrwać zamknięcie okna, zamknięcie przeglądarki i restart rdzenia. Wpis nie zawiera hasła, frazy klucza ani żadnego materiału tajnego: kolumna klucz_kod wskazuje pozycję wykazu kluczy, a ten niesie wyłącznie ścieżkę klucza prywatnego na maszynie rdzenia. Kolumna host_posredni_kod wskazuje inny wiersz tej samej tabeli, czyli host pośredniczący w połączeniu metodą ProxyJump programu OpenSSH. Klucz obcy wskazuje sam na siebie i przy usunięciu ustawia się na wartość pustą: usunięcie hosta pośredniego nie unieważnia wpisów, które przez niego prowadziły połączenie, lecz przywraca je do połączenia bezpośredniego.

## budowa/server/internal/store/migracja_278_narzedzia_zakresy_profilu.sql
Zakres jest nastawą zasięgu, a nie bramką wbudowaną: stanem wyjściowym jest pełny dostęp bez limitu, a wiersz powstaje dopiero wtedy, gdy operator coś przestawił. Brak wiersza oznacza więc pozycję dostępną bez żadnej granicy. Zapis, którego nikt nie czyta przy wykonaniu, byłby gorszy niż jego brak, ponieważ operator widziałby ograniczenie, którego nikt nie egzekwuje; dlatego obok zakresu stoi tabela zużycia, bez której pole limitu wywołań byłoby samą etykietą. Zużycie liczy się wierszem na wywołanie, a nie licznikiem narastającym w kolumnie zakresu, ponieważ limit obowiązuje w oknie czasu i licznik narastający musiałby być zerowany przez proces wiedzący, kiedy okno się przesunęło. Wiersze niosące znacznik czasu odpowiadają na pytanie o liczbę wywołań w ostatnich sekundach jednym zapytaniem i nie wymagają żadnego procesu w tle.

## budowa/server/internal/store/migracja_282_agent_uprawnienie_moduly.sql
Poprzedni schemat zamykał kolumnę grupa na czterech wartościach, podczas gdy zestaw grup zakresu uprawnień eksperta liczy ich pięć: dostęp do modułów i zasobów jest jedną z nich i nie mieścił się w słowniku kolumny. Bez tej wartości zawężenie uprawnień do wybranych modułów nie miało gdzie się zapisać. SQLite nie zmienia warunku CHECK istniejącej kolumny w miejscu, dlatego tabela przechodzi przebudowę: żaden klucz obcy nie wskazuje na agent_uprawnienie, więc przebudowa obejmuje wyłącznie tę tabelę. Zakres szczegółowy wpisu grupy modułów niesie kod modułu z katalogu platformy; warunku sprawdzającego istnienie tego kodu nie ma i być nie może, ponieważ katalog modułów jest osobną tabelą, a warunek CHECK nie sięga innych tabel — wskazanie sprawdza rdzeń, który ten katalog zna.

## budowa/server/internal/store/migracja_373_wylaczenia_pamieci.sql
Do tej migracji schemat nie znał żadnej tabeli wyłączeń, więc żądanie wyłączenia ustalenia wspólnego w projekcie albo module kończyło się odmową, mimo że powinno być możliwe wyłączenie pamięci częściowo albo całkowicie z poziomu konfiguracji operatora. Wyłączenie nie jest usunięciem: wiersz tabeli wylaczenie_pamieci nie zmienia wpisu pamięci projektu, treść pozostaje nietknięta, a zniesienie wyłączenia przywraca wpis do kontekstu w całości. Wyłączenie nie jest też odpięciem: odpięcie zawęża zasięg samego wpisu, podczas gdy wyłączenie zostawia zasięg wpisu nietknięty i jedynie wstrzymuje jego obowiązywanie w zasięgu wskazanym. Wyłączeniu podlega jeden z dwóch bytów: albo pojedynczy wpis, albo cały poziom pamięci; wiersz z obydwoma polami pustymi nie wyłącza niczego i schemat go odrzuca warunkiem sprawdzającym. Zasięg wyłączenia korzysta z tej samej tabeli poziomów zasięgu, którą posługuje się każdy wpis pamięci i każde ustawienie platformy, obejmującej zasięg globalny, środowisko, projekt, moduł, parę modułów i kartę sesji. Poziom pamięci zapisuje się wartością zgodną z kontraktem, ponieważ poziomy pamięci eksperta korzystają z tego samego wyliczenia, a drugie odwzorowanie tych samych czterech wartości byłoby drugim źródłem prawdy.

## budowa/server/internal/store/migracja_374_wyciszenie_nakladki.sql
Do tej migracji wyciszenie nakładki było stanem jednego okna, przechowywanym w zapisie miejscowym przeglądarki, więc wyciszenie sugestii w jednej powłoce nie obejmowało pozostałych. Wyciszenie nakładki ma sięgać wszystkich powłok jednocześnie, dlatego wymaga wiersza w bazie rdzenia. Tabela rozróżnia trzy rodzaje wyciszenia: czasowe, kontekstowe oraz wyciszenie klasy zdarzeń; tryb cichy obecności i wyjątek wagi krytycznej nie są wyciszeniami w tym sensie i nie mają tu wiersza, ponieważ mieszałyby dwa różne znaczenia jednej tabeli. Kolumna czasu końca niesie wyłącznie wyciszenie czasowe; pozostałe dwa trwają do zniesienia ręcznego, co pusta wartość wyraża wprost. Nazwa bytu zakresu stoi obok jego identyfikatora, aby wykaz wyciszeń dało się odczytać bez odpytywania wykazu modułów i wykazu kart sesji. Identyfikator urządzenia niesie informację, skąd wyciszenie przyszło, a nie gdzie obowiązuje, ponieważ obowiązuje ono wszędzie.

## budowa/server/internal/store/migracja_375_sygnal_nakladki.sql
Zestaw sześciu klas zdarzeń wyzwalających nakładkę obejmuje więcej sygnałów, niż rdzeń widział własną telemetrią: stan pętli wykonawczej i stan kolejki zadań były widoczne, natomiast wynik kontroli jakości, harmonogram przebiegów automatyk i powtarzalność czynności operatora nie miały żadnego nośnika, więc ich wyciszenie nie miało czego wyciszać, a sugestia tych klas nie miała z czego powstać. Tabela sygnal_nakladki jest nośnikiem, a nie drugą telemetrią: sygnał odkłada tu ten składnik rdzenia, który go zaobserwował — okno pętli wykonawczej, okno modułu po kontroli jakości, przebieg automatyki po uruchomieniu z harmonogramu albo dziennik czynności operatora po przekroczeniu progu powtarzalności. Sygnał wyciszony odkłada się nadal, ponieważ wyciszenie wstrzymuje wyłącznie ujawnienie sugestii operatorowi, a nie zapis zdarzenia; skasowanie sygnału w chwili wyciszenia zabrałoby operatorowi możliwość zobaczenia go po zniesieniu wyciszenia. Liczba wystąpień stoi w wierszu, ponieważ próg powtarzalności czynności ręcznej liczy się z niej. Sześć klas zdarzeń jest zapisanych warunkiem sprawdzającym, a nie osobną tabelą słownikową, ponieważ stanowią wyliczenie ustalone w kontrakcie, a nie katalog, do którego operator dopisuje własne wiersze.
## budowa/server/internal/store/migracja_111_indeks_tresci_biblioteki.sql
Bajty każdego wgrania leżą na dysku pod swoją sumą kontrolną; bez tego indeksu każde żądanie wyszukiwania musiałoby otworzyć i przeczytać wszystkie bloby repozytorium, więc koszt wyszukiwania rósłby z rozmiarem całej biblioteki, a nie z liczbą trafień. Wybrano FTS5, nie dopasowanie wzorca po wyciągu tekstu, ponieważ budowa niesie SQLite z obsługą FTS5, więc pełnotekstowy indeks nie dokłada zależności ani procesu, a dopasowanie wzorca byłoby skanem każdego wiersza przy każdym żądaniu i nie umiałoby dopasować słów niezależnie od wielkości liter i znaków diakrytycznych. Baza nie przechowuje samego pliku: w indeksie leży wyłącznie wyciąg tekstowy treści bieżącej, obcięty granicą, jako byt wtórny odtwarzalny z blobów magazynu, po którego utracie biblioteka traci tylko trafność wyszukiwania. Rowid indeksu równa się identyfikatorowi pliku biblioteki, ponieważ tabela FTS5 nie ma kluczy obcych ani kaskady, więc wiązanie idzie przez rowid nadawany wprost przy zapisie; wersje historyczne pliku nie są indeksowane, bo dawałyby trafienia w pliki, w których szukanego słowa już nie ma.

## budowa/server/internal/store/migracja_112_trwanie_sesji_bramki.sql
Wcześniejsza migracja dała sesji bramki kolumnę chwili końca i nic więcej, a przełącznik pozostawania zalogowanym rozstrzygał wyłącznie o pierwszym wyliczeniu tej chwili, po czym ginął — odnowienie tokenu nie miało skąd wiedzieć, na jak długo sesję przedłużyć, więc brało stałą rdzenia. Ponieważ pierwsze uruchomienie klienta woła odnowienie zawsze, sesja roczna kurczyła się do doby roboczej przy pierwszym starcie, dokładnie tak jak przełącznik miał to uniemożliwić. Migracja zapisuje trwanie, nie flagę logiczną, ponieważ rdzeń potrzebuje odpowiedzi na pytanie, o ile przesunąć wygaśnięcie, a nie na pytanie, czy Operator chciał pozostać zalogowany; trwanie przeżywa zmianę stałych rdzenia i nie wymusza, żeby baza znała nazwę przełącznika z kontraktu. Wartość domyślna zero oznacza sesję sprzed tej migracji — jawny ślad wiersza założonego, gdy kolumny jeszcze nie było — i dla takich wierszy rdzeń bierze trwanie podstawowe, zachowując się dokładnie tak, jak zachowywał się dotąd.

## budowa/server/internal/store/migracja_113_rodzaje_zasobow_arsenalu.sql
Kolumna rodzaju zasobu dopuszczała warunkiem CHECK tylko trzy wartości, a rdzeń wytwarza też film, dźwięk, dokument i archiwum, których kontrakt niesie cztery brakujące wartości, rozbijające się dotąd o warunek schematu. Migracja przebudowuje tabelę, ponieważ SQLite nie zmienia warunku CHECK w miejscu — jest on częścią tekstu tworzenia tabeli — a kolumny, typy, wartości domyślne i pozostałe warunki są przy tym przepisane co do znaku, bez innej zmiany niż poszerzenie jednej listy wartości. Tabela etykiet jest przebudowywana razem z zasobami, ponieważ wskazuje na zasób kaskadą kasowania, a rdzeń trzyma klucze obce włączone przez cały czas migracji w transakcji, więc porzucenie starej tabeli zasobów przy włączonych kluczach skasowałoby kaskadą wszystkie etykiety; stąd kolejność jest odwrócona — najpierw kopiowane są etykiety do tabeli tymczasowej, potem porzucany oryginał, potem tabela zasobów, na końcu zmieniane nazwy, a SQLite przy zmianie nazwy tabeli sam przepisuje odwołania klucza obcego w tabeli potomnej. Wiersze zastane zostają nietknięte celowo: poprawianie starych błędnych zapisów zgadywaniem na podstawie formatu pliku dałoby wiersz w panelu nieodróżnialny od prawdziwego, więc stare wiersze mówią to, co mówiły, a prawdę niosą wiersze zakładane od tej migracji.

## budowa/server/internal/store/migracja_116_skrzynka_operatora.sql
Aplikacja nie ma własnego serwera poczty i używa skrzynki, którą Operator ma już skonfigurowaną na urządzeniu albo w chmurze, więc migracja nie zakłada żadnej nowej tabeli opisującej pocztę, tylko dopowiada, jak rdzeń łączy się z cudzą skrzynką. Druga tabela na skrzynkę Operatora byłaby drugą prawdą o tym, do czego rdzeń się loguje, ponieważ rodzina poczty widziałaby jedne skrzynki, a wyzwalacz automatyk rodzaju poczty drugie, dlatego migracja dokłada sześć kolumn do bytu, który już istnieje. Protokół staje się polem, ponieważ odbiór może iść przez IMAP, JMAP albo POP3, więc przestał być stałą; domyślną wartością jest IMAP, zgodnie z tym, czego oczekuje kontrakt przy braku wskazania. Źródło nastaw zapisuje fakt o pochodzeniu wartości — czy wpisał je Operator, czy odczytano je z klienta poczty na urządzeniu, czy pochodzą z chmury — ponieważ bez tej kolumny podpowiedź odczytana automatycznie wyglądałaby po zapisie identycznie jak nastawy wpisane ręcznie, a to są dwie różne rzeczy co do aktualności. Szyfrowanie jest dwiema kolumnami, bo odbiór i wysyłka jadą osobnymi gniazdami do osobnych serwerów z osobnymi ustawieniami szyfrowania, więc jedna wspólna kolumna zmuszałaby do zgadywania trybu z numeru portu. Domyślność rozstrzyga za Operatora, gdy komenda nie wskazała konkretnej skrzynki wprost, a indeks częściowy pilnuje, że domyślna skrzynka jest co najwyżej jedna.

## narzedzia/nowy-teren.sh
Wieloakapitowy nagłówek złożono w jedno zdanie o przyczynie i skutku otwarcia
terenu; przykład wywołania usunięto z nagłówka, bo powtarza treść komunikatu
błędu widocznego przy niepoprawnym wywołaniu skryptu.

## narzedzia/przygotuj-drzewo.sh
Nagłówek pliku scalono w jedno zdanie łączące cel i wywołującego; osobny
komentarz przed pętlą dowiązań rozwinięto o wzmiankę, że dowiązanie dotyczy
każdego pakietu z osobna, bo krótszy zapis nie mieścił się w wymaganej długości.
## budowa/server/internal/store/migracja_118_zalaczniki_wiadomosci.sql
Wysyłanie wiadomości przyjmuje pole załączników, niesie je w odpowiedzi i podaje modelowi, więc bez tej kolumny tabela wiadomości nie miała gdzie ich trzymać, a odczyt listy wiadomości oddawał w tym miejscu wartość pustą — model wracający do rozmowy po restarcie albo po przewinięciu historii nie widział, że w rozmowie były pliki. Kolumna, nie osobna tabela, ponieważ kontrakt opisuje załączniki jako zwykły wykaz odwołań, nie byt z własnym życiem: nie ma pola, którego załącznik nie dzieliłby z wiadomością, więc osobna tabela dałaby złączenie i drugie repozytorium po to, żeby przechować listę napisów. Kolumna niesie tablicę JSON napisów, dokładnie tak, jak brzmi pole kontraktu, a wartość pusta oznacza wiadomość bez załączników i tak też wraca do kontraktu jako pole nieobecne, nie pusta tablica, ponieważ wszystkie wiersze sprzed tej migracji mają wartość pustą z braku wiedzy o ich załącznikach, a brak wiedzy nie ma prawa wyglądać jak pewność, że załączników nie było.

## budowa/server/internal/store/migracja_119_zasieg_aplikacja.sql
Osiem wcześniejszych poziomów zasięgu opisuje treść prowadzoną w platformie, lecz żaden nie opisuje samego programu, a nastawy takie istnieją i są używane, jak wymóg logowania czy adres nasłuchu; bez tego poziomu wymóg logowania dawał się ustawić wyłącznie przy starcie rdzenia, bo nie było poziomu, na którym wolno by go zapisać, ani drogi, którą zapis dotarłby do warstwy nasłuchu bez restartu. Nowej tabeli nastaw migracja nie zakłada, ponieważ osobne miejsce na ustawienia serwera byłoby drugim źródłem prawdy o nastawach — zamiast tego dokłada wiersz do istniejącego już słownika poziomów, dzięki czemu cała rodzina komend konfiguracji obsługuje nowy poziom bez jednej nowej komendy. Pierwszeństwo zero stawia poziom aplikacji poniżej wszystkich ośmiu dotychczasowych, więc zapis na którymkolwiek z nich nadal wygrywa i rozstrzyganie sprzed tej migracji nie zmienia się ani o jeden klucz; poziom nie ma własnego bytu do wskazania, tak samo jak poziom globalny. Migracja poprawia warunek CHECK schematu zamiast przebudowywać tabelę, ponieważ zwykła przebudowa jest wykluczona przy włączonych kluczach obcych, które podczas migracji przepisałyby odwołania wszystkich tabel wskazujących przenoszoną tabelę na nazwę odstawioną, a klucza obcego nie da się wyłączyć wewnątrz trwającej transakcji migracji; zostaje więc poprawka tekstu schematu, po której numer wersji schematu rośnie o jeden, aby połączenia trzymające schemat w pamięci przeczytały go na nowo.

## budowa/server/internal/store/migracja_120_katalog_wymogu_logowania.sql
Poprzednia migracja dołożyła poziom zasięgu aplikacji jako miejsce, w którym nastawa może zamieszkać, lecz sam poziom nie wystarcza, ponieważ zapis ustawienia sprawdza klucz wobec katalogu, a katalog żyje w osobnej tabeli definicji ustawień; bez wiersza katalogu zapis klucza wymogu logowania byłby odrzucany, a odmowa byłaby prawdziwa, bo cichego zapisu bez skutku nie ma, więc ten krok dokłada brakujący wiersz. Kategoria bezpieczeństwa jest wybrana, ponieważ tam już mieszka uwierzytelnianie, a nowa kategoria byłaby drugim miejscem na to samo pojęcie w oknie ustawień. Wartość domyślna jest pusta, nie fałszywa, ponieważ wymóg logowania ma trzy stany, nie dwa: brak wskazania, po którym rozstrzyga adres nasłuchu, wskazanie twierdzące i wskazanie przeczące — wartość domyślna fałszywa skasowałaby stan pierwszy i po cichu zniosłaby wymóg na nasłuchu wystawionym poza pętlę zwrotną. Zmiana nie wymaga restartu, ponieważ warstwa nasłuchu sprawdza wymóg przy każdym nowym połączeniu, a rdzeń czyta nastawę przy powitaniu połączenia, więc nowa wartość obowiązuje od następnego połączenia. Dozwolony jest wyłącznie zasięg aplikacji, bo nasłuch jest jeden dla całego programu i nie ma czym go zawęzić do okna czy sesji.

## budowa/server/internal/store/migracja_232_design_prompt_okno.sql
Odczyt historii promptów zwraca prompty wydane w jednym oknie, a tabela prompt_design okna nie znała: prompt powstawał jako parametr wywołania, a nie jako zapis zdarzenia w oknie. Bez tej kolumny historia mieszałaby prompty wszystkich okien albo musiałaby wnioskować okno z zasobów powstałych z promptu, przez co prompt bez ani jednego zasobu wypadałby z historii, choć jest właśnie tym, czego operator w niej szuka. Kolumna kanału jest zapisywana obok, ponieważ kontrakt o nią pyta, a rdzeń zna ją wyłącznie w chwili wydania promptu. Pusta wartość domyślna jest zamierzona: wiersze zastane powstały przed tą kolumną i ich okna nie da się odtworzyć, więc takie prompty nie wejdą do żadnej historii, co odzwierciedla stan rzeczywisty zamiast go ukrywać.

## budowa/server/internal/store/migracja_233_design_wersje_kompozycji.sql
Zapis układu bieżącego zastępuje poprzedni układ w całości: warstwy są usuwane i wstawiane od nowa w jednej transakcji, więc ciągu kolejnych stanów nie ma skąd wziąć — wersja jest migawką, którą ten zapis pomija. Warstwy wersji leżą w osobnej tabeli, a nie w jednym polu zapisu strukturalnego, ponieważ przywrócenie wersji przepisuje je z powrotem do tabeli warstw kompozycji kolumna w kolumnę; ten sam kształt po obu stronach oszczędza przekład, który mógłby zgubić pole. Odczyt wykazu wersji nie sięga po warstwy, dlatego liczba warstw jest utrwalona osobno w wierszu wersji, żeby wykaz nie ważył tyle, co cały zapis. Identyfikator warstwy wersji jest kopią identyfikatora warstwy żywej, a nie odwołaniem do niej, ponieważ warstwa, z której migawkę zdjęto, bywa już usunięta z kompozycji, a wersja ma ją odtworzyć pod tym samym identyfikatorem; z tego powodu identyfikator nie jest jednoznaczny w całej tabeli, tylko powtarza się w każdej wersji, w której dana warstwa została zapisana.

## budowa/server/internal/store/migracja_234_design_adnotacje.sql
Pole adnotacji w warstwie kompozycji niosło jedno zdanie bez autora, bez czasu i bez wątku, a przy każdym zapisie układu wracało przepisane od nowa razem z całym zestawem warstw. Uwaga zostawiona przez jedną osobę znikała więc przy pierwszym przesunięciu warstwy przez drugą osobę. Adnotacja ma odtąd własny wiersz, własny czas i własnego autora, dzięki czemu przeżywa każdy kolejny zapis układu. Wątek powstaje przez pole wskazujące inną adnotację tej samej kompozycji jako nadrzędną; wskazanie jest identyfikatorem zewnętrznym, a nie więzem obcym, ponieważ odpowiedź w wątku zakłada się zaraz po adnotacji nadrzędnej, oba wiersze bywają wstawiane w jednym przebiegu okna, a więz obcy narzuciłby porządek wstawiania, którego kontrakt nie zna. Adnotacja przypięta do całej kompozycji, a nie do pojedynczej warstwy, ma pole warstwy puste, co kontrakt dopuszcza wprost.

## budowa/server/internal/store/migracja_235_design_zestawy_zetonow.sql
Motyw produktu jest własnością powłoki i rdzeń go nie nadpisuje. Zestaw żetonów jest bytem obok motywu: operator zakłada go, wczytuje z zapisu zewnętrznego, wydaje do kodu i porównuje z motywem obowiązującym; bez własnej tabeli panel czytałby żetony z motywu i nie miałby ich gdzie odłożyć, przez co każda praca nad systemem projektowym kończyłaby się z zamknięciem karty. Żeton jest wierszem, a nie polem zapisu strukturalnego, ponieważ wydanie do arkuszy stylów i języków programowania idzie żeton po żetonie, a odsyłacz roli rozstrzyga się po nazwie roli w obrębie zestawu; para zestawu i nazwy roli jest kluczem głównym bez osobnego surogatu, ponieważ dwie wartości tej samej roli w jednym zestawie tworzyłyby sprzeczny system. Motyw jest kolumną zestawu, a nie żetonu, ponieważ oba motywy są równoprawne i niosą własne wartości tych samych ról — motyw jasny i ciemny stanowią więc dwa osobne zestawy, a nie jeden z podwójnymi wierszami.

## budowa/desktop/src-tauri/build.rs
Opis funkcji ogłaszającej zależność budowy złożono w jedno zdanie łączące
przyczynę i skutek; komentarz o katalogu zasobów skrócono poniżej stu znaków,
zachowując istotę: bez zejścia w assets/ zmiana zasobów nie wznawia budowy.

## budowa/server/internal/store/migracja_247_terminal_biblioteka_skryptow.sql
Zapis skryptu zakłada kolejną wersję, a numer nadaje rdzeń; pozycja z jedną kolumną treści traciłaby poprzednie brzmienie przy każdym zapisie, a skrypt uruchamiany na maszynach operatora jest dokładnie tym rodzajem treści, do której trzeba wrócić po nieudanej poprawce. Tabela pozycji niesie więc treść bieżącą, żeby wykaz dało się odczytać jednym zapytaniem, a tabela wersji niesie pełny ślad zmian. Alias jest skrótem przywołującym snippet z palety poleceń, więc musi być jednoznaczny w obrębie biblioteki; warunek jednoznaczności stoi na wyrażeniu, a nie wprost na kolumnie, ponieważ aliasu nie ma większość pozycji, a wartość pusta w SQLite nie zderza się sama ze sobą, więc wiele pozycji bez aliasu współistnieje bez przeszkody.

## budowa/server/internal/store/migracja_249_terminal_tunele.sql
Do tej migracji tunel dało się założyć wyłącznie poleceniem wydanym w karcie: proces przekierowania portu biegł jako zwykłe polecenie, a jego stan i to, czy port po stronie rdzenia w ogóle nasłuchuje, pozostawały niewidoczne. Tunel dostaje własny wiersz, ponieważ jest bytem długożyjącym o własnym stanie, w odróżnieniu od polecenia, które ma początek i koniec. Wiersz istnieje także wtedy, gdy po restarcie rdzenia tunel już nie biegnie, ponieważ operator ma zobaczyć, że tunel był założony i dlaczego przestał działać; przy montażu rdzeń przestawia tunele zastane w stanie czynnym na nieczynny, bo wykazywanie ich jako działających byłoby nieprawdą. Kolumn liczników przesłanych bajtów tabela nie ma, ponieważ przepustowość jest wielkością chwili, a nie dziennika: zapisana w bazie starzeje się między dwoma odczytami, a rdzeń liczy ją wyłącznie przy żywym tunelu i nie wpisuje zera tam, gdzie jej nie zna, ponieważ zero oznaczałoby brak transferu, a nie brak wiedzy o nim.

## budowa/server/internal/store/migracja_250_terminal_obserwacje.sql
Obserwacja uruchamia polecenie karty przy zmianie plików pasujących do wzorca. Wyzwalacz plikowy istniał dotąd wyłącznie w automatykach, poza powłoką terminala, więc nie dało się powiązać zmiany w katalogu z poleceniem uruchamianym w konkretnej karcie, jej katalogu i jej środowisku. Licznik wyzwoleń i chwila ostatniego wyzwolenia są dziennikiem, a nie stanem żywym: po restarcie rdzenia obserwacja nie działa i przy montażu przechodzi w stan zatrzymany, ale liczba dotychczasowych wyzwoleń pozostaje, dzięki czemu operator odróżnia obserwację założoną i niedziałającą od takiej, która po prostu nic jeszcze nie wykryła. Klucza obcego do tabeli kart terminala nie ma z tego samego powodu, dla którego nie ma go karta do okna: karta bywa bytem pamięci rdzenia bez własnego wiersza, a wpis obserwacji ma powstać niezależnie od tego.

## budowa/server/internal/store/migracja_251_terminal_karta_cel_zdalny.sql
Adres powłoki zdalnej wchodził dotąd zmienną środowiskową karty, a zmienne środowiska karty z zamysłu nie mają własnej kolumny, ponieważ bywają nośnikiem poświadczeń, a baza nie jest sejfem. Skutkiem było to, że karta zdalna odtworzona po restarcie rdzenia traciła adres celu, a pierwsze polecenie wydane w niej kończyło się odmową. Adres celu poświadczeniem nie jest — jest tą samą wartością, która widnieje w książce hostów i w nazwie karty — więc jego kolumna niczego z sejfu do bazy nie przenosi; hasło i klucz pozostają poza tabelą tak samo jak dotychczas. Kolumny dokładają się osobnymi poleceniami zmiany tabeli, bez jej przebudowy, ponieważ warunek sprawdzający na kolumnie powłoki pozostaje nietknięty, a wiersze zastane dostają wartości puste oznaczające kartę lokalną.

## budowa/server/internal/store/migracja_252_terminal_karta_powloki_urzadzen.sql
Kontrakt zna dziesięć rodzajów powłoki karty, a poprzedni warunek sprawdzający kolumny wymieniał sześć. Rdzeń nauczył się obsługiwać pozostałe cztery rodzaje — kontener, pod, konsolę szeregową i sesję Telnet — więc karta takiego rodzaju powstawała w pamięci, ale jej zapis odbijał się od warunku kolumny; skutek pozostawał niezauważony, ponieważ karta działała do restartu rdzenia, a po nim znikała bez żadnego sygnału ostrzegawczego. Warunku sprawdzającego nie da się w używanej bazie zmienić samym poleceniem zmiany tabeli, dlatego tabela jest odtwarzana od nowa: powstaje nowa tabela, treść zostaje przepisana, nazwa zamieniona, a indeks odtworzony. Kolumny są wymienione wprost, a nie przez odczyt wszystkich kolumn na raz, aby przepisanie zależało od schematu zapisanego w tej migracji, a nie od kolejności kolumn zastanej w bazie.
## budowa/server/internal/store/migracja_122_narzedzia_sesji.sql
Raz załadowane narzędzie jest obecne agentowi do końca pracy w sesji: nie przechodzi na stałe do zestawu agenta ani nie jest na jednorazowe użycie, lecz trwa do zakończenia sesji albo usunięcia rozmowy. Definicja eksperta jako miejsce zapisu odpada, ponieważ obowiązuje wszystkie sesje eksperta, nie tylko tę jedną, a pamięć tury odpada, bo dołożenie nie znika po jednym użyciu, więc nie może żyć w bycie kończącym się razem z turą. Zostaje więc stan sesji: klucz obcy do sesji z kasowaniem kaskadowym daje dokładnie potrzebny czas życia — dołożenia przeżywają rozłączenie klienta, bo wiersz leży w bazie, a giną razem z rozmową przy jej usunięciu, więc czasu życia nie pilnuje ani zadanie sprzątające, ani warstwa wyższa, tylko klucz obcy. Wiersz zapisuje pozycję wraz z nazwami i opisem, nie sam klucz obcy do wykazu, ponieważ wykaz po ukośniku nie jest tabelą, tylko składa się na bieżąco z komend kontraktu i katalogów, więc odinstalowanie rozszerzenia zamieniłoby dołożenie oparte na samej nazwie w wiersz bez opisu. Rodzaj wylicza dwie wartości, bo po ukośniku idą powołania narzędzia i komendy akcji, choć dołożone bywają wyłącznie pozycje rodzaju narzędzia, a tego pilnuje warstwa wyższa, nie baza, aby warunek bazy nie stał się drugą regułą tego samego wyboru. Kolumna źródła zapisuje, czyja ręka dołożyła narzędzie — komenda wpisana przez Operatora czy dołożenie zrobione przez asystenta — tym samym uzasadnieniem, co pole sprawcy przy zdarzeniu zmiany sesji. Czas jest liczbą milisekund epoki, tak jak w kolumnie aktualizacji rozszerzenia, ponieważ baza nie wstawia własnego znacznika czasu, by uniknąć dwóch prawd o jednej chwili.

## budowa/server/internal/store/migracja_123_wstrzymanie_kroku.sql
Krokiem w tym produkcie jest wiersz pozycji kolejki, nazywany krokiem wprost przez silnik kolejki; drugiego bytu na krok zlecenia nie ma i ta migracja go nie zakłada. Stan pozycji kolejki nie zna stanu wstrzymany i nie ma gdzie zapisać decyzji Operatora, a wstrzymanie całej kolejki nie dotyka pojedynczej pozycji, więc ta tabela zapisuje wstrzymanie pojedynczego kroku wraz z decyzją Operatora, żeby wznowienie wracało do stanu sprzed wstrzymania, a nie traktowało pracy przerwanej jak skończonej. Schemat pilnuje trzech rzeczy: wstrzymać da się wyłącznie krok żywy, ponieważ warunek CHECK dopuszcza cztery stany robocze i żaden stan końcowy; jedno czynne wstrzymanie na krok, bo indeks częściowy nie dopuszcza drugiego wstrzymania tego samego kroku przed zastosowaniem pierwszej decyzji; oraz że decyzja nie ginie, ponieważ trzy osobne znaczniki czasu rozdzielają fakt podjęcia decyzji, jej zastosowania i jej doręczenia do wykonawcy. Tabela rozróżnia stan sterowania od stanu pracy: zapisuje wyłącznie wartości wstrzymany, zatwierdzony i odrzucony, a stany czekania i biegu czyta się ze stanu pozycji kolejki, bo zapisanie ich tutaj byłoby drugą prawdą o tym samym kroku. Wiersz wstrzymania nie kasuje się po zastosowaniu decyzji, ponieważ historia wstrzymań kroku jest częścią przejrzystości pętli wykonania, tak jak dziennik akcji kolejki, do którego te zdarzenia również trafiają.

## budowa/server/internal/store/migracja_125_konto_wlasciciela.sql
Konto właściciela jest jedno, a schemat pilnuje tego warunkiem na kluczu głównym, nie umową w kodzie, ponieważ platforma prowadzi jednego właściciela na dowolnej liczbie urządzeń, więc drugi wiersz nie jest stanem, który wolno osiągnąć. Login i adres e-mail są tożsamością konta, a hasła w tej tabeli nie ma: skrót hasła leży poza bazą, w sejfie poświadczeń, a wiersz metody uwierzytelnienia niesie wyłącznie odwołanie do niego, tak samo jak przed dołożeniem konta. Kolumna potwierdzenia rozdziela dwa stany konta wymagane przy rejestracji — konto powstaje niepotwierdzone i pozostaje takie do chwili potwierdzenia adresu, a dopiero potwierdzenie wydaje urządzeniu token dostępu. Droga potwierdzenia tożsamości korzysta z jednej tabeli na dwa cele, ponieważ mechanizm jest ten sam: jednorazowy materiał wysłany listem, ważny przez czas ograniczony. W bazie leży skrót drogi, nigdy sama droga, więc wyciek kopii bazy nie daje możliwości potwierdzenia cudzej tożsamości, ponieważ ze skrótu nie da się odtworzyć materiału, który poszedł listem; kolumna użycia zamyka drogę po pierwszym użyciu, aby ta sama droga nie otwierała konta wielokrotnie.

## budowa/server/internal/store/migracja_260_automations_wersje.sql
Kolumna wersji na wierszu automatyki mówi wyłącznie, ile razy definicję zapisano, a nie jak wyglądała za którymś razem, więc porównanie dwóch wersji i przywrócenie wcześniejszej nie miały z czego powstać. Tabela wersja_automatyki daje im źródło: migawkę kroków zapisywaną przy każdym zapisie definicji. Kroki leżą w tej migawce jako zapis strukturalny, a nie jako osobne wiersze, ponieważ wersja jest migawką martwą — nikt jej nie edytuje po fakcie ani nie pyta o pojedynczy krok wybranej wersji; rozbicie migawki na tabelę kroków dołożyłoby drugą prawdę o kroku obok tabeli kroków automatyki i wymagałoby jej utrzymywania przy każdej zmianie schematu kroku.

## budowa/server/internal/store/migracja_261_automations_etykiety.sql
Etykieta automatyki jest wierszem, a nie polem z wartościami rozdzielonymi przecinkiem, ponieważ wyszukiwanie po etykiecie ma korzystać z indeksu, a nie z dopasowania podciągu, które trafiałoby krótszą nazwę wewnątrz dłuższej pokrewnej.

## budowa/server/internal/store/migracja_262_automations_zmienne.sql
Wartość domyślna zmiennej przepływu leży jako zapis strukturalny, ponieważ kontrakt niesie ją typem ogólnym: zmienna bywa liczbą, tekstem albo zapisem złożonym, a kolumna o jednym typie skalarnym zmuszałaby rdzeń do zgadywania przy odczycie. Odwołanie do sekretu nigdy nie niesie wartości poświadczenia, wyłącznie referencję do skarbca; sama wartość poświadczenia pozostaje poza bazą danych.

## budowa/server/internal/store/migracja_263_automations_mapowania.sql
Krok jest wskazywany w mapowaniu własnym kodem, a nie kluczem wiersza tabeli kroków automatyki, ponieważ zapis mapowania następuje równocześnie z zapisem kroków, a zapis kroków podmienia wiersze w całości; klucz obcy do wiersza kroku kasowałby więc mapowanie przy każdym zapisie definicji. Mapowanie wskazujące krok, który już nie istnieje, jest zastrzeżeniem oddawanym w wyniku odczytu, a nie odmową na poziomie schematu.

## budowa/server/internal/store/migracja_248_terminal_klucze_ssh.sql
Wykaz kluczy SSH nie jest sejfem: w wierszu stoi ścieżka klucza prywatnego na maszynie rdzenia, treść klucza publicznego i odcisk, nigdy materiał tajny. Klucz prywatny nie opuszcza dysku maszyny rdzenia ani przy wytworzeniu, ani przy wciągnięciu do wykazu przez wskazanie ścieżką. Kolumna hasła niesie jedną wartość logiczną informującą, że klucz jest chroniony hasłem, a samo hasło leży poza tą tabelą, w osobnym sejfie poświadczeń. Warunek jedności na kolumnie ścieżki zapobiega dwóm wpisom wskazującym ten sam plik, co byłoby nierozróżnialne w chwili usuwania pliku z dysku.

## budowa/desktop/src-tauri/src/aktualizacja/probne.rs
Testy sprawdzianu aktualizacji pracują na prawdziwych plikach: tworzą je, podmieniają przez zmianę nazwy i sprawdzają stan dysku, więc katalog próbny musi być prawdziwym katalogiem systemowym, nie atrapą w pamięci. Sprzątanie biegnie w implementacji `Drop`, więc uruchamia się również przy odwijaniu paniki, gdy test padnie w trakcie sprawdzenia. Katalogi dwóch testów nie mogą na siebie wpadać: testy biegną równolegle w jednym procesie i te same testy potrafią biec jednocześnie w dwóch drzewach roboczych, dlatego ścieżka niesie zarówno nazwę podaną przez test, jak i numer procesu.

## budowa/server/internal/store/migracja_264_automations_notatka_polozenie.sql
Notatka przy kroku i położenie węzła na kanwie leżą kluczowane kodem kroku, a nie kluczem wiersza tabeli kroków automatyki, ponieważ zapis definicji podmienia komplet wierszy kroków przy każdym zapisie przepływu; kolumna w tabeli kroków kasowałaby się więc przy każdej zmianie nazwy kroku. Adnotacja kroku, którego już nie ma, zostaje w tabeli i nie przeszkadza, ponieważ odczyt idzie od kroków do adnotacji, nigdy odwrotnie — krok przywrócony z wersji wcześniejszej zastaje swoje położenie na miejscu.
Odwołania do skarbca idą kolumną kroku, a nie adnotacją, celowo przeciwnie do rozdziału powyżej: odwołanie do poświadczenia jest częścią samej definicji kroku, więc krok, który przestał wywoływać usługę zewnętrzną, przestaje potrzebować klucza, a kolumna znika razem z krokiem. Bez tej kolumny usuwanie poświadczenia nie miałoby jak nazwać kroków, które straciły pokrycie: odwołania przychodziłyby żądaniem zapisu definicji i ginęłyby w locie, a skarbiec zdejmowałby klucz bez możliwości powiadomienia o skutku. Wykaz leży jako zapis strukturalny, ponieważ kontrakt niesie go wykazem tekstów ustalanym w całości wraz z krokiem.

## budowa/desktop/src-tauri/src/awaria_startu.rs
Powłoka Tauri panikuje samodzielnie, gdy funkcja składania okna i zasobnika zwróci błąd podczas startu, zanim jakiekolwiek okno zdąży powstać, ponieważ okna nie są deklarowane statycznie w konfiguracji. Profil wydania ma ustawioną strategię paniki przerywającą proces natychmiast, a podsystem okienkowy poza kompilacją debugową nie daje konsoli — bez własnego haka panika kończy proces bez konsoli, bez okna i bez wpisu w dzienniku, znikając bez śladu dla Operatora. Hak paniki uruchamia się zawsze przed przerwaniem procesu niezależnie od strategii paniki, dlatego jego instalacja jest pierwszą instrukcją funkcji main.
## budowa/server/internal/store/migracja_131_szablony_studia.sql
Migracja 131 — szablony fabryczne modułu Studio.

Opracowanie modułu wymienia pięć układów predefiniowanych: pismo, umowa,
raport, notatka, oferta. Bez nich `studio.template.list` oddaje pustkę,
a `studio.template.apply` nie ma czego zastosować — komenda działałaby
poprawnie i bezużytecznie.

── Dlaczego treść stoi w migracji, a nie w kodzie ──────────────────────────
Szablon jest DANĄ, nie zachowaniem: Operator go czyta, kopiuje i zmienia,
a nowy szablon ma być wierszem tabeli, nie wydaniem produktu. Zaszycie go
w kodzie kazałoby przebudować rdzeń, żeby dołożyć układ pisma.

── Dlaczego pola są osobne od treści ───────────────────────────────────────
Treść niesie znaczniki `{{nazwa}}`, a wykaz pól mówi, które z nich Operator
ma wypełnić i jak się nazywają po ludzku. Bez wykazu okno musiałoby zgadywać
pola, parsując treść — i pomyliłoby się przy pierwszym znaczniku w cytacie.

Treści są układami pustymi, nie przykładami: żadnych zmyślonych stron umowy,
kwot ani nazwisk. Miejsce na dane wskazuje znacznik pola.
## budowa/server/internal/store/migracja_140_okno_warsztatu_dokumentu.sql
Migracja 140 — Warsztat dokumentu wchodzi do katalogu okien.

Moduł Studio dostał okno rodzin `studio.pdf.*` i `studio.security.*`. Bez
wiersza w katalogu klient postawiłby okno, o którym rdzeń nie wie: `module.list`
zaniżałby zakres modułu, a pas uczciwości okna meldowałby kod „poza katalogiem".
Ten sam rozjazd domykała migracja 126 dla siedmiu innych modułów.

Rola `pomocnicze`, bo okno wspiera pracę Studio Editora, a nie prowadzi jej
samo: pracuje na materiale wniesionym do okna, nie na treści redagowanej.
Kategoria `narzedzia` wzorem Tools Panelu, z którym dzieli rolę w module.
Kolejność w kategorii jest pierwszą wolną po pozycjach migracji 126.
## budowa/server/internal/store/migracja_141_developer_punkty_przerwania.sql
Migracja 141 — punkty przerwania okna Run & Debug modułu Developer.

Punkt przerwania przeżywa sesję debugowania i musi ją przeżyć. Operator
stawia go w marginesie Code Editora zanim cokolwiek uruchomi, a potem
uruchamia debugowanie po raz drugi i trzeci — punkt postawiony w pamięci
sesji zniknąłby razem z nią i trzeba by go stawiać od nowa przy każdym biegu.
Dlatego punkt należy do okna i pliku, nie do sesji debugowania.

`zweryfikowany` mówi, czy adapter DAP potwierdził, że pod tym wierszem da się
zatrzymać. Punkt niezweryfikowany nie jest usterką: plik bywa jeszcze
nieskompilowany, a Operator ma prawo postawić punkt zanim program powstanie.

Warunek, warunek trafień i treść wpisu stoją osobnymi kolumnami, bo są
osobnymi rodzajami punktu (BreakpointKind kontraktu) i pytanie o nie zadaje
się osobno — sklejenie ich w jedno pole kazałoby rdzeniowi zgadywać, które
z trzech znaczeń niesie zapisany tekst.
## budowa/server/internal/store/migracja_142_developer_kolekcje_api.sql
Migracja 142 — kolekcje zapytań okna API Client modułu Developer.

Kolekcja jest zestawem zapytań HTTP wraz ze środowiskami, w których je się
uruchamia. Zapytanie wpisane raz i zgubione po zamknięciu okna nie jest
klientem API, tylko polem tekstowym — dlatego kolekcja ma miejsce w bazie.

`zapytania` i `srodowiska` trzymamy jako tekst JSON, a nie tabelami
podrzędnymi. Kontrakt niesie je jako surowy JSON (`json.RawMessage`) i nikt
nie pyta o pojedyncze zapytanie kolekcji osobnym żądaniem: kolekcja przychodzi
i odchodzi w całości. Rozbicie jej na wiersze byłoby rozkładaniem i składaniem
tej samej struktury po obu stronach bez jednego odbiorcy tej pracy.

Import OpenAPI zapisuje się tą samą tabelą — kolekcja z importu nie różni się
niczym od kolekcji ułożonej ręcznie poza tym, skąd wzięła zapytania.
## budowa/server/internal/store/migracja_143_developer_polaczenia_danych.sql
Migracja 143 — połączenia bazodanowe okna Data Console modułu Developer.

Wiersz opisuje, DO CZEGO się łączyć, i nie niesie hasła. Hasło stoi w sejfie
pod odwołaniem zapisanym w `poswiadczenie` — ta sama zasada, którą kanały
modeli stosują do kluczy dostawców. Trzymanie hasła w tej tabeli oznaczałoby
kopię sekretu w drugim miejscu, którego nikt nie rotuje.

`tylko_odczyt` jest nastawą połączenia, nie podpowiedzią interfejsu.
Połączenie oznaczone jako tylko do odczytu odrzuca zapytanie zmieniające dane
w rdzeniu, zanim cokolwiek wyjdzie do silnika — konsola SQL nad produkcją bez
takiej nastawy jest jednym nieuważnym poleceniem od szkody.
## budowa/server/internal/store/migracja_144_developer_skany.sql
Migracja 144 — przebiegi skanowania i ich znaleziska (Developer, zakładka
Bezpieczeństwo w Dev Tools).

Skan i znalezisko rozdzielono na dwie tabele, bo pytanie o nie zadaje się
osobno: `developer.scan.run` zakłada przebieg i oddaje jego nagłówek,
a `developer.scan.result.list` czyta znaleziska z filtrem po rodzaju i wadze,
często dla kilku przebiegów naraz. Jedna tabela kazałaby powtarzać nagłówek
przy każdym znalezisku.

Rodzaje skanu przebiegu zapisujemy jako tekst rozdzielony przecinkiem.
Rodzajów jest cztery i są zamkniętym słownikiem kontraktu (ScanKind); tabela
pośrednia na cztery wartości byłaby złożonością bez odbiorcy.

Znalezisko nie ma stanu „przyjęte/odrzucone”. Skan jest pomiarem stanu
repozytorium w danej chwili, a nie listą zadań: kolejny przebieg zakłada nowe
znaleziska, a poprzednie zostają śladem tamtego pomiaru.
## budowa/server/internal/store/migracja_145_developer_testy_pokrycie.sql
Migracja 145 — wyniki testów i pokrycie kodu przebiegu budowania (Developer,
okno Build Output).

Wynik testu powstaje z rozbioru logu przebiegu w chwili jego domknięcia,
a nie z ponownego odpytywania logu przy każdym żądaniu. Powód jest praktyczny:
w dzienniku przebiegu zostaje wyłącznie OGON logu (rdzeń przycina go do
kilkuset wierszy), więc wynik testu odczytany godzinę później nie miałby
z czego powstać. Rozbiór idzie raz, kiedy pełne wyjście jeszcze płynie.

Pokrycie ma własną tabelę, bo jest pomiarem pliku, a nie testu: jeden przebieg
daje setki wyników testów i dziesiątki wierszy pokrycia, i nic ich nie łączy
poza przebiegiem.

`wiersze_bez_pokrycia` trzymamy jako tekst z numerami rozdzielonymi
przecinkiem. Odbiorcą jest nakładka pokrycia w edytorze, która bierze ten
zbiór w całości dla jednego pliku; tabela wiersz-na-wiersz rosłaby o rząd
wielkości bez jednego pytania, na które odpowiadałaby lepiej.
## budowa/server/internal/store/migracja_164_dokument_tlumaczenia.sql
Migracja 164 — dokument wniesiony do tłumaczenia (`translate.document.load`,
`translate.document.render`, `translate.document.layout.compare`) razem
z jego segmentami.

Dokument jest bytem trwałym, bo kontrakt adresuje go identyfikatorem
w dwóch kolejnych komendach: `document.render` i `document.layout.compare`
przyjmują `documentId` wydany przy wczytaniu. Bez wiersza identyfikator
byłby napisem, którego rdzeń przy następnym wywołaniu nie rozpozna.

Segment dokumentu to nie to samo co segment okna (migracja 161). Segment
okna jest kawałkiem tekstu źródłowego; segment dokumentu niesie dodatkowo
miejsce w strukturze pliku: ścieżkę węzła (akapit, komórka, kształt), numer
strony i nazwę stylu. To one pozwalają złożyć dokument z powrotem
z zachowaniem układu, i to ich brak sprawiłby, że `document.render` oddałby
goły tekst zamiast dokumentu.

`uzyto_ocr` jest własnością wczytania, nie dokumentu na dysku: ten sam plik
wczytany dwa razy — raz z warstwy tekstowej, raz z rozpoznania pisma — daje
dwa różne materiały i model ma prawo wiedzieć, na którym pracuje.
## budowa/server/internal/store/migracja_171_zakladki_przegladania.sql
Migracja 171 — zakładki okna przeglądarki (`browser.bookmark.*`).

Zakładka nie jest źródłem przeglądania. Źródło (`zrodlo_przegladania`,
migracja 047) narasta samo w miarę odwiedzania stron i opisuje przebieg
sesji; zakładkę Operator zakłada świadomie i ma ona trwać dłużej niż sesja,
razem z folderem, etykietami i własną notatką. Wspólna tabela kazałaby
odsiewać jedno od drugiego kolumną rodzaju przy każdym odczycie obu wykazów.

Etykiety idą kolumną JSON, bo kontrakt niesie je jako `tags: string[]`
wewnątrz `BrowserBookmark` — nie jest to byt samodzielny, który ktokolwiek
wyszukuje niezależnie od zakładki. Osobna tabela wiązań dawałaby złączenie
przy każdym odczycie i ani jednego nowego pytania, na które umiałaby
odpowiedzieć.
## budowa/server/internal/store/migracja_172_monitory_przegladania.sql
Migracja 172 — monitory zmian strony (`browser.monitor.*`).

Monitor pilnuje strony między sprawdzeniami, więc musi trzymać odniesienie:
treść z chwili założenia albo z chwili ostatniego przyjętego sprawdzenia.
Bez odniesienia „zmieniło się" nie ma względem czego być prawdą, a komenda
`browser.monitor.check` oddawałaby zawsze `changed:false` albo zawsze
`changed:true` — i jedno, i drugie jest meldunkiem bez pomiaru.

Odniesienie idzie odwołaniem do pliku, nie treścią w kolumnie: strona bywa
setkami kilobajtów tekstu, a wzorem jest `migawka_strony.tekst_odwolanie`
z migracji 047. Obok odwołania stoi długość odniesienia — dzięki niej wykaz
monitorów mówi o rozmiarze pilnowanej treści bez sięgania po plik.

Stan sprawdzenia (`pending`, `unchanged`, `changed`, `failed`) jest kolumną,
a nie wyliczeniem z dat: sprawdzenie nieudane (`failed`) też jest
sprawdzeniem, ma swój czas i nie wolno go pomylić z „bez zmian".
## budowa/server/internal/store/migracja_173_kanaly_przegladania.sql
Migracja 173 — kanały RSS/Atom/JSON Feed i ich wpisy (`browser.feed.*`).

Wpis kanału jest bytem osobnym, nie polem kanału: ma własny adres, własny
czas publikacji i własne oznaczenie przeczytania, a wykaz `browser.feed.list`
pyta o kanały z wpisami albo o same kanały (`includeEntries`). Wpisy zapisane
kolumną JSON w wierszu kanału nie dałyby się oznaczyć pojedynczo bez
przepisywania całej kolumny przy każdym przeczytanym wpisie.

Wpis ma więz obcy do kanału z kasowaniem kaskadowym, bo kontrakt mówi wprost:
`browser.feed.remove` zdejmuje subskrypcję „wraz z jej wpisami". Kaskada
w schemacie jest tu jedyną gwarancją, że zdjęcie kanału nie zostawia wpisów
bez rodzica — sprzątanie w kodzie pominęłoby je przy pierwszym błędzie.

Ten sam adres w tym samym oknie nie zakłada drugiej subskrypcji: warunek
UNIQUE(okno, url) czyni z ponownego wywołania `browser.feed.subscribe`
odświeżenie zastanego kanału, a nie jego duplikat.
## budowa/server/internal/store/migracja_174_kolejka_czytania.sql
Migracja 174 — kolejka czytania (`browser.readlist.*`).

Kolejka czytania jest odłożeniem strony na później wraz z przypomnieniem
(opracowanie modułu, rozdz. 2.2). Pozycja przeczytana nie znika z tabeli:
`browser.readlist.remove` z polem `markRead` oddaje pozycję, a nie sam fakt
usunięcia — kolejka ma pamiętać, co już przeczytano, żeby ta sama strona nie
wracała jako nowa.

Przypomnienie jest chwilą, nie flagą: kontrakt niesie `remindAt` jako czas
w milisekundach epoki, więc kolumna trzyma znacznik ISO 8601 tej chwili,
a jego brak znaczy „bez przypomnienia".
## budowa/server/internal/store/migracja_175_wytwory_przegladania.sql
Migracja 175 — wytwory sesji przeglądania i zrzuty stron
(`browser.artifact.add`, `browser.screenshot.capture`,
`browser.snapshot.screenshot.get`).

Wytwór jest wynikiem czynności Operatora na stronie: zrzutem, archiwum,
wyodrębnionymi danymi, adnotacją albo rejestrem sieciowym. Bajty leżą
w magazynie treści pod sumą kontrolną, a tabela trzyma odwołanie — tak samo
jak treść biblioteki i zasoby modułu Design. Wiersz bez odwołania nie ma
prawa powstać: wytwór, za którym nie ma ani jednego bajtu, jest dokładnie tą
szkodą, przed którą stoją sprawdziany skutku tego produktu.

Zrzut ma własną tabelę obok wytworu, bo niesie pomiary, których wytwór nie
zna: tryb (widok, cała strona, obszar, element), format pliku oraz wymiary
w pikselach. Wciśnięcie ich w kolumnę JSON wytworu odebrałoby możliwość
odczytu zrzutu po odwołaniu (`browser.snapshot.screenshot.get` pyta
`screenshotRef` albo `snapshotId`, a nie identyfikatorem wytworu).
## budowa/server/internal/store/migracja_176_pobrania_przegladania.sql
Migracja 176 — menedżer pobrań modułu Browser (`browser.download.*`).

Pobranie jest bytem o własnym cyklu życia: czeka w kolejce, biegnie, bywa
wstrzymane, kończy się powodzeniem, błędem albo przerwaniem. Postęp
(`odebrano_bajtow` wobec `razem_bajtow`) jest liczbą mierzoną w trakcie, nie
opisem — wykaz pobrań ma pokazywać, ile naprawdę leży na dysku.

Komunikat błędu stoi w kolumnie obok stanu, bo „nie udało się" bez powodu
każe Operatorowi zgadywać, czy ponowienie ma sens. Ta sama zasada rządzi
odmowami rdzenia.
## budowa/server/internal/store/migracja_177_makra_przegladania.sql
Migracja 177 — nagrywarka makr przeglądania (`browser.macro.record`).

Makro jest ciągiem kroków w kolejności wykonania, a krokiem kontraktu jest
`AutomationStep` — ten sam kształt, którym jedzie moduł Automations. Kroki
idą kolumną JSON, bo kolejność i kształt kroku należą do makra: krok
wyjęty z makra nie jest niczym samodzielnym, nikt go nie wyszukuje osobno,
a tabela kroków kazałaby przy każdym doklejeniu kroku przenumerować resztę.

Kolumna `nagrywanie` jest stanem nagrywarki, nie własnością zapisu: makro
z ustawionym `nagrywanie = 1` przyjmuje kolejne kroki komendą
`browser.macro.record` z czynnością `step`. Zakończenie zapisu (`stop`)
zdejmuje ten stan i makro staje się zamkniętym scenariuszem.

Kolumna `automatyka_zewnetrzna_id` niesie odwołanie do przebiegu modułu
Automations, gdy scenariusz zostanie tam przekazany (opracowanie, rozdz.
2.13). Nie jest więzem obcym: przekazanie jest konfigurowalne i makro ma
prawo istnieć bez niego.
## budowa/server/internal/store/migracja_178_granice_wykonawcy.sql
Migracja 178 — granice działania Wykonawcy na stronie
(`browser.executor.limits.set`, `browser.executor.limits.get`).

Granice są konfiguracją zasięgu, nie bytem okna: kontrakt niesie je z polem
`scope` (`ConfigScope`) i `scopeId`, a żądanie wskazuje okno ALBO kartę
sesji. Dlatego kluczem jest para (zasięg, wskazanie zasięgu), a nie samo
okno — dwa wiersze o tym samym zasięgu byłyby dwiema odpowiedziami na jedno
pytanie „ile kroków wolno Wykonawcy tutaj".

Domeny dozwolone i zablokowane idą kolumnami JSON, bo są wykazem wewnątrz
jednego ustawienia, a nie bytem wyszukiwanym osobno.

Wiersza domyślnego migracja nie zakłada. Brak wiersza znaczy „granice
domyślne rdzenia" i tak też odpowiada `browser.executor.limits.get` — wartość
domyślna należy do kodu, który ją stosuje, a nie do schematu; wiersz
zaszczepiony w migracji rozjechałby się z nią przy pierwszej zmianie.
## budowa/server/internal/store/migracja_179_zestawy_zrodel_i_watki.sql
Migracja 179 — zestawy tematyczne źródeł, wątki notatek oraz brakujące
własności źródła i notatki (`browser.source.group.*`, `browser.note.thread.*`,
`browser.note.update`).

Kontrakt niesie w `BrowserSource` pole `groupId`, a w `BrowserNote` —
`classification`, `threadId` i `pinned`. Migracja 047 zakładała te tabele,
zanim rodzina komend obejmowała grupowanie i klasyfikację, więc kolumn tych
w nich nie ma. Dopisanie ich tutaj jest jedynym sposobem, żeby okno mogło
oddać notatkę tak oznaczoną, jak ją Operator oznaczył — oznaczenie żyjące
wyłącznie w kliencie ginie przy przeładowaniu karty.

Przynależność stoi po stronie źródła i notatki, nie w kolumnie zestawu:
źródło należy do jednego zestawu, a notatka do jednego wątku, więc skład
zestawu jest zapytaniem po kolumnie, a nie drugą listą do utrzymania.
Kontrakt oddaje `sourceIds` i `noteIds` — rdzeń wylicza je z tej samej
kolumny, którą zapisuje.
## budowa/server/internal/store/migracja_182_biblioteka_reguly.sql
Migracja 182 — moduł Library: reguły repozytorium.

Jedna tabela na trzy rodzaje reguły (kolekcja inteligentna, reguła napływu,
folder obserwowany), bo wszystkie trzy mają ten sam kształt: warunek, cel
i przełącznik czynności. Rozdzielenie ich na trzy tabele powieliłoby warunek
i wykaz — a `library.rule.list` pyta o wszystkie naraz, zawężając rodzajem.

Warunek stoi w zapisie JSON, nie w kolumnach. Kontrakt niesie go jako `json`
(`LibraryRule.condition`: moduł źródłowy, etykiety, rodzaj treści, zakres
dat), a kolumna na każdy człon warunku zamieniłaby dodanie członu w migrację.
Rozbiór warunku należy do rdzenia, który go stosuje.

Przypisanie do kolekcji dostaje pochodzenie. Bez niego usunięcie reguły
kolekcji inteligentnej nie miałoby jak odróżnić zasobu wciągniętego regułą od
zasobu przypisanego ręką Operatora — a `library.rule.remove` z żądaniem
`detachFiles` ma zdjąć wyłącznie te pierwsze. Przypisanie ręczne jest
wartością domyślną, więc wiersze zastane opisują się same.
## budowa/server/internal/store/migracja_185_biblioteka_udostepnienia_webhooki.sql
Migracja 185 — moduł Library: udostępnienia odnośnikiem i nasłuchy zewnętrzne.

Token udostępnienia leży w kolumnie jawnie, nie jako skrót. To jest wybór
zgodny z zasadą jawności kluczy platformy: Operator ma móc odczytać wystawiony
odnośnik i przekazać go powtórnie, a nie wystawiać nowy, bo pierwszy da się
wyłącznie sprawdzić. Zawężeniem dostępu jest termin i odwołanie, nie
nieodczytywalność.

Odwołanie udostępnienia jest znacznikiem czasu, nie skasowaniem wiersza:
„odnośnik przestaje działać, wpis zostaje w dzienniku audytu" — wiersz
odwołany świadczy, że odnośnik istniał.

Nasłuch i jego zdarzenia stoją w dwóch tabelach, bo zdarzenie jest wartością
z wykazu, a nie tekstem: warunek CHECK po stronie wiersza zdarzenia wychwytuje
literówkę przy zapisie, czego lista sklejona w jednej kolumnie zrobić nie może.
## budowa/server/internal/store/migracja_186_biblioteka_sugestie.sql
Migracja 186 — moduł Library: sugestie porządkujące.

Domyślnym zachowaniem modułu jest sugestia z akceptacją Operatora, nie zapis
bez pytania: klasyfikacja wsadowa (`library.classify.run`) wytwarza wiersze
tej tabeli, a dopiero `library.suggestion.apply` zamienia je w zmianę zasobu.
Sugestia musi więc przeżyć między jednym żądaniem a drugim — stąd tabela,
a nie wynik oddany i zapomniany.

Uzasadnienie jest kolumną obowiązkową. Sugestia bez powodu jest poleceniem
podanym bez podstawy, a Operator ma decydować, nie zgadywać, skąd wzięła się
propozycja.

Decyzja zostaje przy wierszu (`stan`), zamiast kasować go przy odrzuceniu:
sugestia odrzucona ma nie wracać przy kolejnym przebiegu klasyfikacji, więc
rdzeń musi wiedzieć, że raz już padła i została odsunięta.
## budowa/server/internal/store/migracja_193_roundtable_analiza.sql
Migracja 193 — ustalenia analizy i rejestr dowodów.

Ustalenie analizy zostaje w zapisie, bo `roundtable.analysis.run` kosztuje
wywołanie kanału modelu. Odczyt wyniku bez powtórnego wywołania jest tu
warunkiem użyteczności: transkrypt z ustaleniami (`includeAnalysis`) wydaje
się długo po tym, jak analiza przebiegła.

Rejestr dowodów jest osobny od ustaleń, bo ma inne pytanie: nie „co model
rozpoznał", tylko „czym poparto twierdzenie". Kolumna `poparte` niesie
odpowiedź, a brak źródła przy `poparte = 0` jest właśnie tym, co
`unsupportedOnly` wyciąga na wierzch.
## budowa/server/internal/store/migracja_194_roundtable_glosowanie.sql
Migracja 194 — głosowanie nad stanowiskami: warianty, głosy i wynik agregacji.

Wynik agregacji NIE jest kolumną. Liczy się go z głosów przy każdym odczycie
(`roundtable.vote.get`), bo metoda agregacji jest własnością głosowania,
a głos może dojść po pierwszym odczycie. Kolumna z wynikiem byłaby drugą
prawdą, rozjeżdżającą się z pierwszą przy każdym kolejnym głosie.

Głos jest jeden na wyborcę. Powtórne oddanie zastępuje poprzedni (ON CONFLICT
w zapisie), bo zmiana zdania w trakcie otwartego głosowania jest czynnością
dozwoloną, a dwa głosy tej samej osoby liczone dwukrotnie nie są.

Kształt głosu zależy od metody, więc kolumny są trzy i wszystkie mogą być
puste: aprobata wypełnia `aprobaty`, metody rankingowe `ranking`, skala
punktowa i metoda kwadratowa `punkty_json`.
## budowa/server/internal/store/migracja_197_roundtable_macierz.sql
Migracja 197 — macierz decyzyjna wariantów i kryteriów z wagami.

Wynik wariantu (`total`) liczy się przy odczycie z ocen i wag, a nie stoi
w kolumnie: zmiana wagi jednego kryterium przestawia wynik każdego wariantu
naraz, więc kolumna wymagałaby przeliczenia całej macierzy przy każdym
zapisie i rozjeżdżałaby się z ocenami przy pierwszym pominięciu.

Oceny wariantu w kryteriach idą jednym polem JSON, bo są mapą „kryterium →
ocena" o kształcie zadanym przez kryteria tej macierzy. Tabela wiążąca
dawałaby ten sam kształt kosztem trzeciego złączenia przy każdym odczycie.
## budowa/server/internal/store/migracja_199_roundtable_wydanie.sql
Migracja 199 — szablony moderacji i artefakty wydane z debaty.

Szablon moderacji zapisuje format, liczbę tur, granicę czasu i kolejność
głosu (opracowanie, 2.8.5). Kolejność jest listą kodów uczestników
rozdzieloną znakiem nowego wiersza, a nie tabelą wiążącą: uczestnik należy do
okna, a szablon ma przeżyć okno, więc więz obcy do składu zabiłby szablon
razem z debatą, z której go zdjęto.

Artefakt debaty niesie odwołanie do bajtów w magazynie treści, nie same bajty.
Transkrypt, graf i nagranie idą w megabajtach, a baza rdzenia trzyma stan,
nie treść — tak samo jak w bibliotece i w module Design.
## budowa/server/internal/store/migracja_204_apps_dziennik.sql
Migracja 204 — moduł Apps: wiersze dziennika usług i wdrożeń.

Jedna tabela na oba dzienniki, bo jest to jeden byt. `apps.service.log.read`
i `apps.deployment.log.read` mają identyczny kształt wyniku (`lines`,
`total`, `streaming`) i różnią się wyłącznie zawężeniem: pierwsza po
komponencie, druga po wdrożeniu. Dwie tabele o tych samych kolumnach byłyby
dwiema prawdami o wierszu dziennika.

Wiersz powstaje z pracy, nie z żądania odczytu. Zapisuje go silnik wykonania
wdrożenia (`adapter_modul_aplikacje_wdrozenie_bieg.go`) przy każdym kroku
przebiegu oraz serwer podglądu przy podniesieniu i zatrzymaniu. Odczyt
niczego nie dopisuje — dziennik pokazujący własne odczyty byłby dziennikiem
o sobie.

`chwila` niesie milisekundy epoki, bo `since` w obu żądaniach jest liczbą
milisekund; porównanie z kolumną tekstową wymagałoby przekładu przy każdym
odczycie.

`wdrozenie_kod` i `komponent_kod` są kodami zewnętrznymi, nie więzami obcymi:
wiersz dziennika przeżywa byt, którego dotyczy — po to się dziennik prowadzi.
## budowa/server/internal/store/migracja_205_apps_artefakty.sql
Migracja 205 — moduł Apps: artefakty budowania.

Artefakt jest plikiem na dysku, nie wpisem o pliku. Wiersz powstaje wtedy,
gdy silnik wykonania wdrożenia spakuje przestrzeń roboczą okna do archiwum
w magazynie treści rdzenia — kolumna `sciezka` trzyma odwołanie względne
magazynu (ten sam wzorzec co `zasob_designu.sciezka`), a `rozmiar`
i `suma_kontrolna` opisują bajty, które tam naprawdę leżą. Wiersz bez pliku
byłby dokładnie tym wzorcem szkody, który w tym produkcie już wystąpił:
wykazem zasobów, za którymi nie ma ani jednego bajtu.

`wdrozenie_kod` jest kodem zewnętrznym, nie więzem obcym — artefakt przeżywa
przebieg, z którego powstał, bo to on idzie potem do pakowania
(`apps.package.build` bierze „artefakt ostatniego wdrożenia udanego").
## budowa/server/internal/store/migracja_206_apps_pakiety.sql
Migracja 206 — moduł Apps, Publisher Panel: pakiety rozszerzenia zbudowane
z produktu.

Pakiet jest archiwum na dysku i wierszem obok niego. `apps.package.build`
składa archiwum z artefaktu budowania i manifestu, kładzie je w magazynie
treści rdzenia i zapisuje tu odwołanie, rozmiar i format. Kolejne komendy
rodziny pracują na tym samym wierszu: `apps.package.manifest.save` wymienia
manifest, `apps.package.validate` czyta go do raportu zastrzeżeń,
`apps.package.sign` dopisuje podpis, `apps.package.publish` — kod pozycji
katalogu, która z pakietu powstała.

Manifest i podpis leżą jako surowy JSON kontraktu (`AppPackageManifest`,
`ExtensionSignature`). Rozłożenie manifestu na kolumny znaczyłoby drugą
definicję kształtu, którego jedynym źródłem jest kontrakt, a narzędzia
i uprawnienia pakietu wychodzą zawsze w komplecie razem z pakietem — nie ma
po czym filtrować.

`rozszerzenie_kod` wskazuje pozycję katalogu kodem, nie więzem obcym: pozycja
żyje w tabeli `rozszerzenie` (migracja 070) własnym cyklem życia i jej
odinstalowanie nie ma prawa skasować pakietu, z którego powstała.
## budowa/server/internal/store/migracja_216_workspace_zadania.sql
Migracja 216 — zadania projektu modułu Workspace.

Zadanie jest jednym bytem trzech widoków huba planowania: pozycją listy,
kartą tablicy i słupkiem osi czasu. Trzech tabel nie ma, bo trzy tabele
rozjechałyby się przy pierwszej zmianie stanu wykonanej z innego widoku.

Etykiety i lista kontrolna leżą w kolumnach tego samego wiersza: etykiety
rozdzielone znakiem nowego wiersza, lista kontrolna zapisem JSON. Osobne
tabele wiążące dawałyby tu wyłącznie koszt złączeń — żadne okno nie pyta
o etykietę bez zadania ani o krok bez zadania.

Klucz porządkowy karty jest NAPISEM, nie liczbą: wstawienie karty między dwie
sąsiednie ma dopisać klucz pośredni, a nie przepisać całą kolumnę.
## budowa/server/internal/store/migracja_217_workspace_zaleznosci.sql
Migracja 217 — zależności między zadaniami projektu.

Zależność wiąże dwa zadania tego samego projektu: poprzednik i następnik.
Warunek UNIQUE na parze pilnuje, żeby ta sama krawędź nie powstała dwa razy —
graf zależności ma jedną krawędź między dwoma zadaniami, a nie tyle, ile razy
Operator kliknął.

Cyklu baza nie wykryje: to jest sprawdzenie rdzenia przed zapisem, bo cykl
rozpoznaje się przejściem grafu, a nie warunkiem kolumny.
## budowa/server/internal/store/migracja_218_workspace_tablica.sql
Migracja 218 — kolumny tablicy kanban projektu.

Kolumna jest nastawą Operatora i odwzorowuje się na stan zadania z zestawu
wyjściowego platformy. Kolumna własna nie zakłada stanu nowego — dlatego
kolumna `stan` ma ten sam warunek co zadanie, a nie własne słownictwo.

Granica prac w toku jest sygnalizowana, nie egzekwowana: zero znaczy brak
granicy, a przekroczenie wraca ostrzeżeniem odpowiedzi, nie odmową.
## budowa/server/internal/store/migracja_219_workspace_notatki.sql
Migracja 219 — notatki i strony wiki projektu wraz z odnośnikami treści.

Strona jest notatką: hierarchię daje wskazanie strony nadrzędnej, a nie druga
tabela. Nagłówki treści leżą w kolumnie obok treści, bo spis treści notatki
czyta się przy każdym otwarciu strony, a parsowanie Markdowna przy każdym
odczycie byłoby liczeniem tego samego po raz drugi.

Odnośnik treści jest osobnym wierszem, bo panel „co linkuje tutaj" pyta
ODWROTNIE niż zapisuje edytor: szuka stron wskazujących tę stronę. Bez
osobnego wiersza trzeba by przeszukiwać treść wszystkich stron projektu.

Odnośnik do strony jeszcze niezałożonej ma pustą kolumnę `notatka_docelowa`.
To jest stan poprawny wiki, nie usterka: nazwa czeka na stronę.
## budowa/server/internal/store/migracja_220_workspace_tablica_wizualna.sql
Migracja 220 — tablica wizualna projektu (płótno, mapa myśli).

Scena leży jednym zapisem JSON i rdzeń jej nie rozbiera. Kształt sceny —
kartki, strzałki, grupy, osadzenia — należy do widoku, który ją rysuje;
rozbiór na wiersze związałby schemat bazy z rysunkiem interfejsu i każda
zmiana kształtu kartki byłaby migracją.
## budowa/server/internal/store/migracja_222_workspace_os_czasu.sql
Migracja 222 — oś czasu aktywności projektu.

Zdarzenie zapisuje się w chwili czynności, a nie wylicza z bytów przy
odczycie. Wyliczanie z bytów pokazałoby wyłącznie to, co jeszcze istnieje:
usunięte zadanie znikałoby również z historii, a wtedy oś czasu przestaje być
zapisem zdarzeń i staje się drugim widokiem stanu bieżącego.
## budowa/server/internal/store/migracja_223_workspace_komentarze.sql
Migracja 223 — komentarze przy bytach projektu.

Jedna tabela na komentarze wszystkich bytów: cel opisuje para rodzaj + byt,
a nie osobna tabela per byt. Wątek składa wskazanie komentarza nadrzędnego —
tak samo jak strona wiki składa hierarchię wskazaniem strony nadrzędnej.

Przywołania znakiem małpy leżą kolumną obok treści: rozpoznaje się je raz,
przy zapisie, bo panel powiadomień pyta o nie częściej, niż komentarz się
zmienia.

## budowa/server/internal/store/migracja_266_automations_parametry_szablonu.sql
Parametr szablonu przepływu jest wierszem, a nie polem zapisu strukturalnego szablonu, ponieważ formularz wpięcia buduje się z niego pole po polu, a zastosowanie szablonu musi umieć nazwać parametr, dla którego nie podano wartości.

## budowa/server/internal/store/migracja_267_automations_publikacja.sql
Publikacja, udostępnienie i budżety czasu automatyki leżą jako cztery kolumny na tabeli automatyki, a nie cztery osobne tabele, ponieważ każda jest polem pojedynczym o krotności jeden do jednego z automatyką i nie ma własnego cyklu życia; osobna tabela na jedną liczbę byłaby złączeniem bez powodu. Pusta wersja opublikowana oznacza, że automatyki nigdy nie opublikowano, a nie że opublikowano wersję zerową — wykonywana produkcyjnie jest wtedy wersja bieżąca. Budżet równy zero oznacza brak granicy czasu, zgodnie z kontraktem pól budżetu.

## budowa/desktop/src-tauri/src/dialog_katalogu.rs
Katalogi robocze są listą na oknie komunikacji, a ich wskazanie musi być czynnością systemu operacyjnego, nie polem tekstowym. Interfejs wskazuje katalog w dwóch sprawach o odmiennym znaczeniu: wskazanie katalogu roboczego, gdzie moduł zostawia swoje pliki, oraz dodanie katalogu jako punktu dostępu, do którego model sięga po treść. Czynność systemu operacyjnego jest w obu przypadkach ta sama, więc okno jest jedno, a zastosowania różni wyłącznie napis w belce podany przez wywołującego. Konsumentem zdarzenia po stronie interfejsu jest most katalogów w kliencie.

## budowa/server/internal/store/migracja_268_kolejka_zlecenia.sql
Tabela pozycji kolejki wprowadzona wcześniej opisuje etap pętli koordynator-wykonawca: tytuł, treść zlecenia, werdykt weryfikacji, licznik obiegów. Zlecenie kolejki jest bytem odrębnym, niosącym ładunek strukturalny, priorytet, termin wykonania, klucz idempotencji i warunek przetworzenia — wtłoczenie jednego w drugie kazałoby kolumnie tytułu nieść ładunek, a kolumnie werdyktu weryfikacji stan zlecenia o innym słowniku. Klucz idempotencji jest unikatowy w obrębie jednej kolejki, nie globalnie, ponieważ ten sam klucz w dwóch kolejkach opisuje dwa różne zlecenia dwóch różnych torów. Zlecenie martwe zachowuje kolejkę źródłową, ponieważ odczyt zleceń martwych bez wskazania kolejki oddaje zadania martwe wszystkich kolejek naraz, więc rozdzielenie ich na osobną tabelę odebrałoby im pochodzenie.
## budowa/server/internal/store/migracja_126_okna_modulow_dobudowane.sql
Migracja 126 — okna dobudowane modułom wchodzą do katalogu rdzenia.

Siedem modułów buduje okna, których katalog rdzenia nie zna, więc `module.list`
zaniża zakres modułu: klient stawia okno, a rdzeń o nim nie wie i pas uczciwości
modułu melduje kod „poza katalogiem". Migracja domyka ten rozjazd — definicje
idą do `okno_operacyjne`, przypięcia do `okno_operacyjne_modul`, tak jak
ustaliły to migracje 030 i 031.

Rola i kategoria każdego okna są dobrane wzorem pozycji już obecnych w katalogu:
okno prowadzące pracę operatora jest `wiodace`, okno pokazujące przebieg jest
`monitor`, okno rządzące zbiorem jest `zarzadca`, okno wspierające inne jest
`pomocnicze`, a okno prowadzące przez tworzenie bytu jest `kreator`.
Kolejność definicji jest pierwszą wolną w obrębie kategorii.

Kolejność w module jest dopisaniem na koniec, nie przestawieniem. Okna zastane
zostają na swoich miejscach: zmiana porządku wyświetlania należy do projektu
interfejsu, nie do migracji domykającej katalog.

## budowa/server/internal/store/migracja_269_kolejka_polityka.sql
Polityka kolejki i dodanie rodzaju kolejki uruchamianej z zegara idą w jednym kroku, ponieważ obie zmiany dotyczą tabeli kolejki, a przenoszenie danych nie zmienia warunku kolumny w miejscu — rozdzielenie ich na dwa kroki oznaczałoby dwa przepisania tej samej tabeli. Polityka idzie kolumnami tabeli kolejki, a nie tabelą obok, ponieważ jest polem bytu kolejki o krotności jeden do jednego, bez własnego cyklu życia; tabela obok kazałaby każdemu odczytowi kolejki wykonywać złączenie, by dowiedzieć się rzeczy o samej kolejce. Rodzaj kolejki dostaje trzecią wartość dla kolejki uruchamianej z zegara, ponieważ wcześniejszy warunek dopuszczał tylko kolejkę sesyjną i kolejkę wieloczynnościową, a bez tej wartości automatyka z harmonogramem nigdy się nie uruchamiała — przebieg odbijał się od warunku kolumny, mimo poprawnie złożonego żądania. Wartości domyślne polityki są wartościami domyślnymi modelu konfiguracji: zasięg lokalny, bez ograniczenia współbieżności i tempa, trzy próby z wycofaniem wykładniczym i rozproszeniem czasowym, kolejka zadań martwych włączona; kolejka założona przed tym krokiem dostaje je tak samo jak kolejka założona po nim.

## budowa/desktop/src-tauri/src/dziennik.rs
Powłoka pracuje bez konsoli w podsystemie okienkowym, więc brak możliwości otwarcia pliku dziennika nie wstrzymuje startu — zapis jest wtedy pomijany, bo dziennik nie jest bramą. Katalog danych jest własnością powłoki, nie rdzenia: zmienna wskazująca katalog danych rdzenia dotyczy maszyny serwera wdrożenia, a powłoka pisze dziennik u siebie, obok pliku nastaw, bo oba pliki są jej własne.
## budowa/server/internal/store/migracja_127_okna_wspolne_platformy.sql
Migracja 127 — okna wspólne platformy w katalogu rdzenia.

Dwa okna towarzyszą każdemu modułowi: Chat Window (kanał Użytkownik ↔
Wykonawca) i Execution Loop Window (kanał Koordynator ↔ Wykonawca). Wszystkie
opracowania modułów wymieniają je razem, na czele wykazu okien, jako kolumny
stałe układu.

Katalog rdzenia znał do tej pory jedno z nich, i to nie dla wszystkich:

  * `execution-loop-window` nie miał wiersza definicji w ogóle, więc żaden
    moduł nie mógł go wskazać, a pas uczciwości modułów meldował ten kod jako
    stojący poza katalogiem;
  * `chat-window` przypina wszystkim modułom migracja 031, ale moduł Library
    wszedł do rejestru po niej, więc przypięcia nie dostał. Jako jedyny
    z piętnastu.

Kolejność zero należy do obu okien wspólnych, nie do jednego. Odczyt sortuje
`ORDER BY om.kolejnosc, o.kod` (`dane/okna_operacyjne.go`), więc remis przy
zerze rozstrzyga kod: `chat-window` stoi przed `execution-loop-window`, czyli
dokładnie tak, jak układa je każde opracowanie. Okna własne modułu zaczynają
się od jedynki i nie są przestawiane.

## budowa/desktop/src-tauri/src/main.rs
Plik pełni wyłącznie kompozycję: brak w nim logiki, typów i obsługi zdarzeń, bo każda odpowiedzialność mieszka w osobnym module. Powłoka niesie okno wraz z wkompilowanym interfejsem i nie niesie rdzenia — rdzeń stoi na serwerze wdrożenia, więc przy starcie nie ma czego stawiać ani na co czekać, a powłoka jedynie czyta wskazanie, gdzie tego rdzenia szukać, i otwiera okno. Lista poleceń wywoływalnych z interfejsu wchodzi na listę zamkniętą: natywne okno wyboru katalogu roboczego, wskazanie serwera rdzenia oraz podmiana pliku aplikacji przy aktualizacji.
## budowa/server/internal/store/migracja_160_pamiec_tlumaczen_pelna.sql
Migracja 160 — pamięć tłumaczeń modułu Translate w kształcie, którego żąda
kontrakt rodziny `translate.memory.*`, oraz polityka pamięci okna.

Tabela `pamiec_tlumaczen` (migracja 054) powstała pod jedną komendę
(`memory.suggest`) i pod jedną drogę zapisu: parę segmentów zdjętą
z zatwierdzonego panelu. Stąd `panel_id NOT NULL`. Rodzina `memory.*` żąda
czego innego: `memory.set` zakłada parę BEZ panelu (Operator wpisuje ją
wprost), `memory.import` wnosi parę z pliku wymiany, a kontrakt
(`TranslationMemoryEntry`) niesie projekt, klienta, autora i kontekst sąsiedni.
Kolumny da się dołożyć poleceniem ALTER; zdjęcia warunku NOT NULL z kolumny
`panel_id` już nie — SQLite tego nie umie. Dlatego tabela powstaje od nowa,
a wiersze zastane przechodzą do niej przepisaniem: pary zebrane dotąd
z paneli zostają parami z panelem, reszta pól zostaje pusta, bo nikt jej
nigdy nie podał.

Zasięg pary (`zasieg`) niesie wartości `TranslationMemoryScope` kontraktu
wprost. Domyślną jest `card` — pamięć własna karcie sesji; pary szersze
(projekt, zespół) powstają, gdy Operator wskaże zasięg sam.

Polityka pamięci jest osobną tabelą, nie kolumnami okna: `memory.policy.get`
odpowiada wtedy „polityki nie ustawiono” brakiem wiersza, zamiast czterema
kolumnami NULL, których nie da się odróżnić od polityki wyzerowanej.

## budowa/server/internal/store/migracja_270_automations_dziennik_przebiegu.sql
Dziennik przebiegu automatyki musi przeżyć restart rdzenia, ponieważ przeglądarka logów pokazuje pełny zapis zdarzeń pojedynczego uruchomienia, także sprzed otwarcia okna, a strumień na żywo niesie wyłącznie zdarzenia zaistniałe przy otwartym oknie. Dziennik trzymany wyłącznie w pamięci procesu byłby dziennikiem, którego po restarcie nie ma, choć przebieg dalej widnieje w historii.
## budowa/server/internal/store/migracja_161_segmentacja_tlumaczenia.sql
Migracja 161 — segmentacja modułu Translate: zestawy reguł podziału
(`translate.segmentation.rules.*`) i trwałe segmenty okna
(`translate.segment.merge`, `translate.segment.split`).

Migracja 053 stwierdzała, że segment nie jest bytem trwałym, bo jedyna
komenda, która go dotykała (`source.segment`), dzieliła tekst w locie.
Scalanie i podział segmentów zmieniają to wprost: po scaleniu dwóch zdań
podział wynikający z tekstu źródłowego przestaje być prawdą o segmentach
okna, a kolejne wywołanie musi zobaczyć wynik poprzedniego. Bez tabeli
`segment_okna_tlumaczenia` obie komendy oddawałyby wykaz, który znika razem
z odpowiedzią — czyli meldunek zamiast skutku.

Wiersze zakładane są leniwie: dopóki nikt nie scalał ani nie dzielił,
okno nie ma ani jednego wiersza i podział liczy się z tekstu źródłowego.
Pierwsze scalenie albo pierwszy podział utrwala cały bieżący wykaz, a potem
zmienia w nim jedną rzecz — inaczej numer segmentu w żądaniu wskazywałby na
inny segment niż ten, który Operator widział.

Zestaw reguł segmentacji jest odrębny od okna: to nastawa wielokrotnego
użytku (kontrakt: `SegmentationRuleset` z własnym identyfikatorem i językiem).
`srx` niesie treść pliku SRX podaną przez Operatora — standard branżowy,
którego rdzeń nie rozkłada na własne reguły; przechowywany w całości, żeby
eksport oddał to samo, co przyszło.
## budowa/server/internal/store/migracja_162_jakosc_i_zatwierdzenia_tlumaczenia.sql
Migracja 162 — trzy rzeczy jednej fasety kontroli w module Translate:
rozszerzenie terminu słownika o stan i dziedzinę (`translate.glossary.list`),
profile kontroli jakości (`translate.qa.profile.*`) oraz obieg zatwierdzeń
panelu (`translate.approval.*`).

Termin dostaje dwie kolumny, nie tabelę: `GlossaryTerm.status`
i `GlossaryTerm.domain` są polami terminu, po których `glossary.list` zawęża
wykaz. Kolumna `stan` ma wartości `GlossaryTermStatus` kontraktu wprost.
Wiersze zastane zostają bez stanu (NULL) — Operator nigdy ich nie oznaczył,
a wpisanie im `approved` byłoby nadaniem zgody, której nikt nie wydał.

Profil kontroli jakości ma dziecko, bo `QaProfileCheck` to trójka
(rodzaj, waga, czy włączony) na każdy z sześciu rodzajów niezgodności —
kolumna z JSON-em byłaby wykazem, którego baza nie umie zawęzić ani sprawdzić.

Zatwierdzenie jest jednocześnie zapisem historii i stanem bieżącym panelu.
Dlatego rośnie tabela `zatwierdzenie_panelu` (kontrakt: `ApprovalRecord`
z własnym identyfikatorem, autorem i chwilą), a panel dostaje trzy kolumny
migawki (`TranslationPanel.approvalStage/approvedBy/approvedAt`) — inaczej
każdy odczyt panelu musiałby dociągać ostatni wiersz obiegu.

## budowa/server/internal/store/migracja_271_automations_kroki_przebiegu.sql
Stan kroku przebiegu wraz z ładunkami wejściowymi i wyjściowymi jest zapisem trwałym, odrębnym od stanu przebiegu wyprowadzanego z pozycji kolejki: kolejka po opróżnieniu traci pozycje, a drążenie do poziomu kroku i podgląd ładunku mają dalej mieć co pokazać, inaczej historia przebiegu kończyłaby się w chwili sprzątnięcia kolejki. Ładunki są redagowane przy odczycie, a nie przy zapisie, ponieważ reguła redakcji sekretów bywa zmieniana, a zapis już zredagowany nie da się odredagować.

## budowa/desktop/src-tauri/src/menu_zasobnika.rs
Menu zasobnika celowo nie ma pozycji zatrzymania rdzenia: rdzeń stoi na serwerze wdrożenia, powłoka go nie postawiła i nie ma czym go wygasić, więc taka pozycja obiecywałaby władzę, której powłoka nie ma.
## budowa/server/internal/store/migracja_163_korekta_tlumaczenia.sql
Migracja 163 — ustalenia korekty językowej modułu Translate
(`translate.proofread.run`, `translate.proofread.apply`).

Ustalenie korekty musi być trwałe, inaczej `proofread.apply` nie ma czego
zastosować: kontrakt każe wskazać ustalenia identyfikatorami
(`findingIds`) w osobnym wywołaniu, a identyfikator wydany w odpowiedzi
i zapomniany po niej byłby identyfikatorem donikąd. To odróżnia korektę od
kontroli jakości (`quality.check`), której zastrzeżenia są migawką
wymienianą w całości i nikt ich nie adresuje pojedynczo.

Ustalenie ma trzy stany życia i wszystkie trzy są tu widoczne: nowe
(`zastosowano` i `odrzucono` puste), zastosowane (Operator przyjął
poprawkę) i odrzucone (Operator ustalenie oddalił). Kasowanie odrzuconych
byłoby zgubieniem odpowiedzi „nie, tak ma być" — kolejny przebieg korekty
zgłosiłby to samo po raz drugi.

Wynik czytelności (`ReadabilityScore`) nie ma tu tabeli. Jest funkcją treści
panelu w chwili pomiaru: liczba zdań, długość słowa, długość zdania. Wiersz
trzeba by unieważniać przy każdej korekcie panelu, a `proofread.run` i tak
liczy go od nowa; nikt też nie adresuje wyniku czytelności identyfikatorem.

## budowa/server/internal/store/migracja_272_automations_punkty_wznowienia.sql
Punkt wznowienia przebiegu musi być trwały, ponieważ służy wznowieniu po awarii, a awaria zabiera pamięć procesu. Wykaz kroków ukończonych leży jako zapis strukturalny, ponieważ czyta się go w całości i tylko w całości — wznowienie pyta, czego już nie powtarzać, a nie czy dany krok jest ukończony.

## budowa/server/internal/store/migracja_273_automations_alarmy.sql
Kanały reguły alarmowania leżą jako zapis strukturalny, ponieważ kontrakt niesie je wykazem tekstów ustalanym w całości: reguła ma kanały takie, jakie zapisano ostatnio. Rozbicie na tabelę wierszy dałoby możliwość stanu, którego kontrakt nie zna — kanału dopisanego bez przepisania reguły.

## budowa/server/internal/store/migracja_274_automations_skarbiec_audyt.sql
Tabela poświadczeń automatyki nie ma i nie będzie miała kolumny na wartość: wartość leży w sejfie plikowym katalogu danych, tym samym, którym idą sekrety kont i punktów dostępu. Baza zna wyłącznie nazwę, zasięg i odwołanie, dzięki czemu odczyt bazy nie może wynieść sekretu — kolumny na niego po prostu nie ma, a nie dlatego, że filtr o tym pamięta.
Dziennik audytu jest zapisem niezmiennym: wiersz raz dopisany nie jest zmieniany ani kasowany przez żadną komendę modułu. Automatyka usunięta zabiera swoje wpisy, lecz wpis o automatyce nieznanej jest dopuszczony, ponieważ audytowana bywa czynność, która żadnej automatyki nie dotyczy.
## budowa/server/internal/store/migracja_165_lokalizacja_i_napisy.sql
Migracja 165 — lokalizacja oprogramowania (`translate.resource.*`) i napisy
(`translate.subtitle.*`, `translate.dubbing.script.build`).

Zasób lokalizacyjny to plik kluczy wniesiony do okna: JSON, YAML, properties,
Android XML, iOS strings i stringsdict, RESX, gettext PO. Klucz jest bytem
adresowanym po nazwie w obrębie zasobu (`resource.key.context.set` przyjmuje
`key`, nie identyfikator), więc para (zasób, klucz) jest kluczem naturalnym.

`formy_mnogie` trzyma JSON, i to jest wyjątek świadomy: liczba form mnogich
zależy od języka (angielski ma dwie, polski trzy, arabski sześć), a nazwy
form są nazwami CLDR (`one`, `few`, `many`, `other`). Tabela dziecka miałaby
tyle wierszy, ile form, i ani jednego zapytania, które by po nich zawężało —
formy czyta się zawsze kompletem razem z kluczem.

`zrzut_zasob_id` wskazuje zasób modułu Design (`zasob_design`) ze zrzutem
ekranu, na którym klucz widać. Odwołanie jest miękkie (kod zewnętrzny, bez
klucza obcego), bo zrzut należy do innego modułu i jego usunięcie nie ma
prawa skasować kontekstu klucza.

Kwestia napisów mieszka przy panelu, nie przy oknie: napisy są tłumaczone,
więc każdy język ma własne taktowanie i własny podział linii. Import napisów
wnosi kwestie źródłowe do okna, dlatego `panel_id` bywa pusty — wtedy kwestia
jest kwestią materiału źródłowego, nie przekładu.
## budowa/server/internal/store/migracja_170_karty_przegladania.sql
Migracja 170 — karty przeglądania, grupy kart i przestrzenie robocze modułu
Browser.

Opracowanie modułu (rozdz. 2.1) stawia kartę, grupę kart i przestrzeń roboczą
obok siebie jako trzy różne byty jednego okna przeglądarki: karta niesie
adres, grupa niesie nazwę i barwę zwijanego zestawu, a przestrzeń robocza jest
zapisanym kompletem kart przełączanym bez utraty stanu. Trzy byty, trzy
tabele — jedna tabela z kolumną „rodzaj" kazałaby każdemu odczytowi odsiewać
dwie trzecie wierszy.

Karta zamknięta zostaje w tabeli. `browser.tab.close` oddaje `closed:true`,
a nie „karty nigdy nie było": karta zamknięta jest częścią historii okna
i przestrzeni roboczej, która ją zapamiętała. Kolumna `zamknieta` odsiewa ją
z wykazu `browser.tab.list`, zamiast kasować wiersz, do którego odwołuje się
zapisany zestaw.

Grupa i przestrzeń nie niosą składu w kolumnie tekstowej. Skład jest po
stronie karty (`grupa`, `przestrzen`), bo karta należy do jednej grupy
i jednej przestrzeni naraz. Wykaz identyfikatorów przepisany do kolumny JSON
byłby drugą prawdą o tej samej przynależności — i rozjechałby się przy
pierwszym zamknięciu karty.

Okno jest kolumną tekstową, nie więzem obcym — tak samo jak w migracji 047:
`windowId` modułu Browser jest oknem operacyjnym, nie oknem komunikacji.

## budowa/server/internal/store/migracja_275_harmonogram_okna_wyzwolenia.sql
Okno wykonania, nadzór obecności uruchomień, wyzwalacz webhook oraz historia wyzwoleń mają przeżyć restart rdzenia: trzymanie ich w pamięci procesu przestałoby ograniczać budzik po pierwszym ponownym złożeniu rdzenia, a harmonogram pokazywałby ograniczenie, które już niczego nie ogranicza. Okno wykonania jest wierszem, ponieważ jest ich wiele na harmonogram — godziny robocze poniedziałku bywają inne niż piątku. Dni tygodnia leżą jako zapis strukturalny, ponieważ kontrakt niesie je wykazem liczb ustalanym w całości.
Nadzór obecności uruchomień i adres wejściowy webhooka są polami pojedynczymi harmonogramu, po jednym na harmonogram, bez własnego cyklu życia. Kolumna tolerancji równa zero wyłącza nadzór, zgodnie z kontraktem. Kolumna webhooka niesie odwołanie do sejfu, nigdy wartość klucza podpisu.

## budowa/desktop/src-tauri/src/nastawy.rs
Adresu serwera wdrożenia nie zna instalator i znać go nie może: w chwili rozpakowania plików nikt jeszcze nie wie, pod jaką nazwą stoi rdzeń danego Operatora. Wiedza pojawia się przy pierwszym uruchomieniu okna i musi przetrwać jego zamknięcie, inaczej przy każdym starcie powłoka pytałaby o to samo. Zmienna środowiskowa wskazująca host rdzenia stoi wyżej niż ten plik jako narzędzie wykonawcy i środowiska serwerowego, gdzie nastawę wnosi jednostka usługi, a nie okno; kolejność warstw rozstrzyga moduł ustawień. Plik leży w katalogu danych powłoki, obok jej dziennika, bo oba pliki należą do powłoki i drugie miejsce zapisu byłoby drugim stanem do pogodzenia przy przenoszeniu profilu.

Odczyt nastaw nie jest bramą i nie wstrzymuje startu okna: brak pliku, plik nieczytelny i treść niezgodna z umową dają nastawy puste, a powód nieczytelności trafia do dziennika powłoki, żeby Operator nie zobaczył ekranu pierwszego uruchomienia bez wyjaśnienia, dlaczego jego poprzednie wskazanie zniknęło. Zapis idzie przez plik przejściowy i przemianowanie, żeby przerwanie w trakcie pisania nie zostawiło pliku obciętego — wskazanie odczytane w połowie byłoby gorsze niż wskazanie nieodczytane, bo okno łączyłoby się z adresem złożonym z połowy nazwy hosta. W przeciwieństwie do odczytu, niepowodzenie zapisu jest bramą: wywołujący ma odmówić Operatorowi, a nie przyjąć wskazanie, które zniknie przy następnym starcie.
## budowa/server/internal/store/migracja_180_biblioteka_zasob_opis.sql
Migracja 180 — moduł Library: cykl życia zasobu, jego miejsce w strukturze
repozytorium oraz opis w schemacie Dublin Core wraz z polami niestandardowymi.

Trzy braki naraz, bo wszystkie trzy dotyczą jednego bytu — zasobu:

  1. `stan` rozdziela wykaz czynny od archiwum. Kosz repozytorium
     (`library.file.archive` / `library.file.restore`) jest przeniesieniem
     między stanami, nie usunięciem wiersza: zasób zarchiwizowany zachowuje
     wersje, etykiety i kolekcje, więc przywrócenie oddaje go w całości.
     Usunięcie trwałe (`library.file.delete`) zdejmuje wiersz i wtedy dopiero
     kaskada zabiera wersje.
  2. `sciezka_repozytorium` jest drogą WEWNĄTRZ biblioteki
     (`LibraryFile.path`), rozłączną z kolumną `sciezka`, która niesie
     ścieżkę źródłową z maszyny Operatora i z rdzenia nie wychodzi
     (`dane/library.go`). Bez osobnej kolumny `library.file.move` nie miałby
     dokąd przenieść zasobu, a wypełnienie pola kontraktu ścieżką źródłową
     wyniosłoby do klienta układ cudzego dysku.
  3. Opis Dublin Core mieszka w tabeli towarzyszącej, nie w kolumnach
     `plik_biblioteki`: piętnaście pól opisowych obciążałoby każdy odczyt
     wykazu, a wykaz opisu nie pokazuje. Jeden wiersz opisu na jeden zasób —
     klucz główny jest kluczem obcym.

Pola niestandardowe stoją dwutorowo, bo są dwiema różnymi rzeczami:
DEFINICJA pola należy do repozytorium (`pole_schematu_biblioteki`,
`library.schema.set`), a WARTOŚĆ pola do zasobu (kolumna
`pola_niestandardowe` opisu, mapa kod→wartość w zapisie JSON). Rozdział ten
ma skutek wprost w kontrakcie: zdjęcie definicji nie kasuje wartości
zapisanych przy zasobach i wartości wracają, gdy pole zostanie założone
ponownie.
## budowa/server/internal/store/migracja_181_biblioteka_slownik_kolekcje.sql
Migracja 181 — moduł Library: słownik etykiet, tezaurus i hierarchia kolekcji.

Etykieta była dotąd wolnym tekstem w tabeli złącznikowej
(`etykieta_pliku_biblioteki`, migracja 045) i tym pozostaje przy zasobie.
Słownik jest bytem NAD tym tekstem: niesie barwę, czas założenia i sam fakt
istnienia etykiety, której dziś nie nosi żaden zasób. Bez słownika
`library.tag.list` mógłby pokazać wyłącznie etykiety użyte, więc etykieta
nieużywana — ta, którą komenda ma umieć wskazać i usunąć — nie istniałaby
dla rdzenia wcale.

Wpis słownika nie jest warunkiem noszenia etykiety. Zasób otagowany
`library.tag.set` etykietą nową dostaje ją natychmiast, a wpis słownika
powstaje przy okazji; klucz obcy w drugą stronę zamieniłby tagowanie
w dwuetapowy obrządek i wywrócił zapis przy wyścigu dwóch wgrań.

Tezaurus łączy dwie etykiety relacją modelu SKOS. Relacja jest bytem
symetrycznym w zapisie (para nazw + rodzaj), a odwrotność wyprowadza odczyt:
`nadrzedna` czytana od drugiej strony jest `podrzedna`, więc zapisywanie obu
kierunków dałoby dwa wiersze mówiące to samo i rozjazd, gdy zniknie jeden.

Kolekcje dostają rodzica i regułę. Rodzic daje hierarchię o dowolnej
głębokości (`LibraryCollection.parentId`), reguła — kolekcję inteligentną
(`LibraryCollection.ruleId`), której zawartość wynika z warunku, a nie
z ręcznego przypisania. Kolumna reguły nie ma klucza obcego, bo tabela reguł
powstaje krok dalej (182), a kolejność kroków jest jednokierunkowa.

## budowa/desktop/src-tauri/src/okno.rs
Interfejs pochodzi wyłącznie z pakietu wkompilowanego w powłokę i z żadnego innego miejsca — to jest cały produkt na urządzeniu Operatora, okno wraz z interfejsem, bez rdzenia. Adresu strony nie ma czego rozstrzygać w czasie pracy: nastawa budowania rozstrzyga go raz, a powłoka nie niesie żadnego adresu zapasowego, w tym adresu serwera rozwojowego.
## budowa/server/internal/store/migracja_183_biblioteka_audyt.sql
Migracja 183 — moduł Library: dziennik audytu repozytorium.

Dziennik jest przyrostowy: wiersz raz dopisany nie jest zmieniany ani
kasowany żadną komendą rodziny `library.*`. Stąd brak kolumny zmiany i brak
kolumny stanu — wpis opisuje zdarzenie, które już zaszło, a zdarzenie zajść
nie przestaje.

Wskazanie zasobu jest luźne z zamysłu: `plik_kod` niesie identyfikator
zewnętrzny zasobu, a nie klucz obcy do `plik_biblioteki(id)`. Klucz obcy
z kaskadą zabrałby wpisy razem z zasobem usuniętym trwale — czyli zdjąłby
ślad dokładnie tej czynności, dla której dziennik istnieje. Wpis o usunięciu
ma przeżyć usunięcie.

Czynności obejmujące całość repozytorium (wywóz paczki, skanowanie
duplikatów) wchodzą bez wskazania zasobu — kolumna jest pusta i to jest
stan opisany kontraktem (`LibraryAuditEntry.fileId` opcjonalne).
## budowa/server/internal/store/migracja_184_biblioteka_retencja_utrwalenie.sql
Migracja 184 — moduł Library: polityki przechowywania i utrwalenie archiwalne.

Polityka retencji nie usuwa niczego sama. Przechowuje trzy rzeczy: poziom
zasięgu, na którym obowiązuje, liczbę dni przechowywania liczoną od ostatniej
zmiany zasobu oraz czynność po upływie okresu. Czynność `usuniecie` znaczy
„zgłoś do usunięcia", nie „usuń": trwałe usunięcie ma własną komendę
i własne potwierdzenie (`library.file.delete`). Bez tego rozdziału polityka
zapisana pomyłkowo zabierałaby zasoby bez śladu decyzji człowieka.

Poziom zasięgu zapisuje się kodem z katalogu `poziom_zasiegu` (kontrakt:
`ConfigScope`), więc warunku CHECK tu nie ma — wykaz poziomów należy do
tabeli katalogu, a jego powielenie w warunku byłoby drugą prawdą o zasięgu.

Zadanie utrwalenia jest zapisem SKUTKU, nie zleceniem do wykonania. Wiersz
powstaje po pracy: niesie wynik walidacji i jej zapis, żeby Operator mógł
wrócić do pytania „czy ten dokument naprawdę przeszedł normalizację" bez
powtarzania utrwalenia. Zasób wytworzony wskazywany jest kodem, bo bywa go
brak — profil PREMIS/METS opisuje zasób, nie wytwarza nowego pliku.
## budowa/server/internal/store/migracja_190_roundtable_sklad.sql
Migracja 190 — skład debaty rozszerzony do kształtu kontraktu: tożsamość
uczestnika, rola w naradzie, waga kompetencji, oznaczenie kluczowego oraz
zespoły i biblioteka ról.

Migracja 044 zakładała skład pod cztery komendy obszaru i wprost odnotowała,
że kolumny „kluczowy" nie ma, bo żadna komenda nie potrafiła jej ustawić.
Kontrakt ma dziś `roundtable.model.update` z polami `key`, `weight`, `role`,
`avatar`, `roleDescription` i `agentId` — kolumny przestały być miejscem,
do którego nic nie pisze, więc powstają.

Zespół (`roundtable.team.save`) jest kopią składu, nie odwołaniem do niego.
Skład okna zmienia się po zapisaniu zespołu, a zespół ma zostać taki, jaki
był w chwili zapisu — inaczej „wnieś zespół" wnosiłoby stan bieżący cudzego
okna zamiast zapamiętanego układu.

## budowa/desktop/src-tauri/src/polecenia.rs
Powłoka udostępnia interfejsowi wyłącznie polecenia, których przeglądarka wykonać nie może: reszta pracy idzie kanałem WebSocket kontraktu, więc powłoka nie tworzy drugiej drogi sterowania platformą. Zatrzymania ani postawienia rdzenia tu nie ma i nie będzie: rdzeń stoi na serwerze wdrożenia, powłoka go nie niesie i nie ma nad nim władzy.

Polecenie wyboru katalogu obsługuje jednym poleceniem oba zastosowania interfejsu — wskazanie katalogu roboczego i dodanie katalogu jako punktu dostępu — bo czynność systemu operacyjnego jest ta sama, a różni je wyłącznie napis w belce podany przez wywołującego. Interfejs nie ma jak wyliczyć adresu HTTP rdzenia z lokalizacji dokumentu, bo strona pochodzi z pakietu wkompilowanego w powłokę — polecenie zwracające adres jest jedyną drogą, którą warstwa połączenia interfejsu poznaje ten adres. Okno pyta o stan wskazania rdzenia przed złożeniem aplikacji: odpowiedź bez wskazania znaczy pierwsze uruchomienie po instalacji, wtedy staje ekran wskazania, bo instalator adresu serwera wdrożenia nie zna i znać go nie może. Wskazanie serwera jest jedynym poleceniem powłoki zmieniającym jej stan trwały, bo dotyczy pliku nastaw na dysku, którego przeglądarka dotknąć nie może.

Polecenie aktualizacji pobiera adres i sumę kontrolną z wykazu wydań czytanego przez interfejs po HTTPS; powłoka wykazu nie czyta i wersji nie porównuje, bo do tego wystarczy przeglądarka, a powłoka robi dwie rzeczy niemożliwe ze strony: pisze po dysku pod plikiem aplikacji i stawia proces na nowo. Praca idzie wątkiem roboczym, bo pobranie wydania to dziesiątki megabajtów — na wątku głównym okno stałoby zamrożone przez cały czas ściągania, a baner aktualizacji nie zdążyłby się nawet przerysować.

## budowa/desktop/src-tauri/src/rdzen/mod.rs
Powłoka rdzenia nie stawia i nie wygasza: rdzeń nie stoi ani w instalce, ani na urządzeniu Operatora, tylko na serwerze wdrożenia.
## budowa/server/internal/store/migracja_191_roundtable_tura.sql
Migracja 191 — tura i wypowiedź doprowadzone do kształtu kontraktu.

Trzy braki naraz, wszystkie w warunku CHECK albo w brakującej kolumnie:

 1. Format. Migracja 044 dopuszczała cztery formaty, a kontrakt ma sześć:
    doszły `delphi` (rundy anonimowe) i `expertPanel` (panel ekspercki).
    Warunku CHECK nie da się poszerzyć poleceniem ALTER — stąd przebudowa.
 2. Tura nadrzędna. `roundtable.debate.branch` zakłada wariant tury,
    a `roundtable.debate.followup` wątek boczny. Obie potrzebują wskazania
    tury, przy której stoją; bez niego wariant byłby zwykłą kolejną turą.
 3. Granice tury. Opracowanie modułu ma zegar tury i granicę długości
    wypowiedzi (2.8.2), a kontrakt pola `timeLimitMs`, `maxStatementChars`
    i `anonymous`.

Wypowiedź dostaje redakcję (`roundtable.statement.regenerate` zastępuje
treść, a numer redakcji odróżnia zastąpienie od pierwszego głosu), akt mowy
(klasyfikacja z `roundtable.analysis.run`), pewność deklarowaną (wejście do
kalibracji) i odwołanie do wypowiedzi, na którą odpowiada.

Przebudowa idzie parami, bo `debata_wypowiedz` wiąże się kluczem obcym
z `debata_tura`: przemianowanie tabeli wskazywanej przeciąga za sobą
deklarację klucza w tabeli wskazującej, więc obie muszą powstać na nowo
w jednym kroku.

## budowa/desktop/src-tauri/src/rdzen/nasluch.rs
Rdzeń stoi na serwerze wdrożenia, a nie na urządzeniu Operatora, więc pytanie o łączność idzie zawsze pod wskazany serwer, nigdy pod pętlę zwrotną — pętla zwrotna dawałaby odpowiedź fałszywą niezależnie od rzeczywistego stanu rdzenia.
## budowa/server/internal/store/migracja_195_roundtable_ocena.sql
Migracja 195 — ocena Operatora, rubryki oceny i werdykty modeli-sędziów.

Ocena Operatora i werdykt sędziego są osobnymi bytami, choć obie „oceniają".
Różnią się autorem i skutkiem: ocena Operatora wchodzi do rankingu jako
pojedynek rozstrzygnięty ręcznie, werdykt sędziego niesie punkty w kryteriach
rubryki i uzasadnienie wypowiedziane przez model.

Rubryka bywa wspólna dla platformy albo związana z jednym oknem. Okno puste
znaczy rubrykę wspólną — `roundtable.rubric.list` bez wskazania okna oddaje
wtedy same wspólne, a ze wskazaniem wspólne i te jednego okna.
