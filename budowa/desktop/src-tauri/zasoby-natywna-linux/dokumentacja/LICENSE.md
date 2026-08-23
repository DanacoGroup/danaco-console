# Licencja produktu Danaco Console

**Umowa licencyjna oprogramowania Danaco Console — Platforma AI Workspace OS**

| | |
|---|---|
| **Produkt** | Danaco Console |
| **Rodzaj produktu** | Platforma AI Workspace OS — system operacyjny pracy z modelami sztucznej inteligencji |
| **Producent (Licencjodawca)** | Danaco Holding Group Sp. z o.o., ul. Gen. J. Hallera 81, 43-400 Cieszyn |
| **Twórca** | Dariusz Naharnowicz |
| **Kontakt** | support@danaco-group.pl |
| **Wersja produktu** | v2.0 |
| **Wersja dokumentu licencji** | 2.0 |
| **Data dokumentu** | 2026-08-18 |
| **Zakres obowiązywania** | wszystkie egzemplarze Danaco Console v2.0 |

---

## Spis treści

1. [Postanowienia wstępne](#1-postanowienia-wstępne)
2. [Definicje](#2-definicje)
3. [Przedmiot licencji](#3-przedmiot-licencji)
4. [Udzielenie licencji i zakres uprawnień](#4-udzielenie-licencji-i-zakres-uprawnień)
5. [Zasady korzystania i ograniczenia](#5-zasady-korzystania-i-ograniczenia)
6. [Własność intelektualna](#6-własność-intelektualna)
7. [Komponenty osób trzecich i rozszerzenia](#7-komponenty-osób-trzecich-i-rozszerzenia)
8. [Aktualizacje i wersjonowanie](#8-aktualizacje-i-wersjonowanie)
9. [Dane i treści Operatora](#9-dane-i-treści-operatora)
10. [Rękojmia i odpowiedzialność](#10-rękojmia-i-odpowiedzialność)
11. [Okres obowiązywania i zakończenie licencji](#11-okres-obowiązywania-i-zakończenie-licencji)
12. [Zgodność licencyjna](#12-zgodność-licencyjna)
13. [Poufność](#13-poufność)
14. [Postanowienia końcowe](#14-postanowienia-końcowe)
15. [Załącznik A — komponenty osób trzecich](#15-załącznik-a--komponenty-osób-trzecich)

---

## 1. Postanowienia wstępne

### 1.1. Charakter dokumentu

Niniejszy dokument („**Licencja**") określa warunki, na jakich Danaco Holding
Group Sp. z o.o. udostępnia oprogramowanie Danaco Console — Platformę AI
Workspace OS („**Oprogramowanie**") do korzystania. Licencja reguluje zakres
przyznanych uprawnień, granice dozwolonego użycia, zasady dotyczące komponentów
osób trzecich wykorzystanych w Oprogramowaniu, postępowanie z danymi
wytworzonymi w toku korzystania oraz zasady odpowiedzialności Licencjodawcy.

Oprogramowanie jest **produktem własnościowym**. Nie jest oprogramowaniem
otwartoźródłowym ani darmowym w rozumieniu swobody redystrybucji i nie podlega
żadnej z licencji wolnego oprogramowania. Wykorzystanie w Oprogramowaniu
komponentów osób trzecich rozpowszechnianych na licencjach otwartych (rozdz. 7)
nie zmienia własnościowego charakteru Oprogramowania ani nie rozciąga warunków
tych licencji na kod autorski Licencjodawcy — z zastrzeżeniem obowiązków
wynikających wprost z licencji poszczególnych komponentów.

### 1.2. Strony stosunku licencyjnego

Stronami stosunku licencyjnego są:

1. **Licencjodawca** — Danaco Holding Group Sp. z o.o. z siedzibą w Cieszynie,
   ul. Gen. J. Hallera 81, 43-400 Cieszyn, podmiot, któremu przysługują
   autorskie prawa majątkowe do Oprogramowania, występujący w dokumentacji
   produktu również jako **Producent**;
2. **Licencjobiorca** — podmiot, osoba fizyczna albo osoba prawna, któremu
   Licencjodawca udostępnił Egzemplarz Oprogramowania na warunkach niniejszej
   Licencji.

Dokumentacja produktu posługuje się pojęciem **Operator** na oznaczenie osoby,
która faktycznie obsługuje Oprogramowanie: prowadzi konto właściciela,
konfiguruje kanały modeli, komponenty własne, rozszerzenia i punkty izolacji,
prowadzi sesje i oba kanały komunikacji operacyjnej oraz odpowiada za wynik
pracy wykonanej z udziałem modelu. Licencja jest udzielana Licencjobiorcy;
Operator wykonuje uprawnienia Licencjobiorcy w jego imieniu, a Licencjobiorca
odpowiada za jego działania i zaniechania jak za własne.

### 1.3. Hierarchia dokumentów

W razie rozbieżności między dokumentami dotyczącymi Oprogramowania stosuje się
następującą kolejność pierwszeństwa:

1. **odrębna umowa** zawarta między stronami na piśmie, w zakresie, w jakim
   reguluje daną kwestię odmiennie;
2. **niniejsza Licencja**;
3. **dokumentacja produktu** przekazana Licencjobiorcy wraz z Egzemplarzem —
   pliki `README.md`, `INSTALACJA-I-KONFIGURACJA.md`,
   `INSTRUKCJA-UZYTKOWANIA.md` oraz `CHANGELOG.md`;
4. dokumentacja projektowa i wewnętrzna Licencjodawcy — **nieprzekazywana
   Licencjobiorcy i nieskładająca się na treść zobowiązania**.

Wykaz z pkt 3 jest zamknięty: dokumentem produktu jest wyłącznie plik wprost
w nim wymieniony. Dokument nienależący do tego wykazu — choćby dotyczył
Oprogramowania i pochodził od Licencjodawcy — nie tworzy zobowiązania co do
zakresu funkcjonalnego (rozdz. 14.5). Materiały prezentacyjne, poglądowe
i koncepcyjne nie stanowią części Licencji.

---

## 2. Definicje

### 2.1. Podmioty

| Pojęcie | Znaczenie |
|---|---|
| **Licencjodawca**, **Producent** | Danaco Holding Group Sp. z o.o. |
| **Twórca** | Dariusz Naharnowicz — podmiot autorskich praw osobistych do Oprogramowania |
| **Licencjobiorca** | Podmiot, któremu udostępniono Egzemplarz na warunkach Licencji |
| **Operator** | Osoba fizyczna obsługująca Oprogramowanie w imieniu Licencjobiorcy; właściciel konta w rozumieniu dokumentacji produktu |
| **Dostawca modelu** | Podmiot udostępniający model sztucznej inteligencji albo narzędzie dostępu do niego |

### 2.2. Oprogramowanie i jego warstwy

| Pojęcie | Znaczenie |
|---|---|
| **Oprogramowanie** | Danaco Console w postaci wykonywalnej: Aplikacja serwerowa, Klient, Powłoka, Kontrakt, schemat trwałości wraz z migracjami oraz zasoby wbudowane |
| **Aplikacja serwerowa** (Rdzeń) | Serwerowa część Oprogramowania w języku Go: prowadzi trwałość, obsługuje Kontrakt, zarządza sesjami, procesami i kolejkami, uruchamia wywołania modeli i przekazuje strumień odpowiedzi |
| **Klient** | Część Oprogramowania w języku TypeScript odpowiadająca za prezentację i przyjmowanie poleceń Operatora |
| **Powłoka** | Natywna warstwa okna na urządzeniu Operatora, osadzająca Klienta w komponencie prezentacji systemu operacyjnego |
| **Pakiet kliencki** | Instalowalny pakiet Powłoki wraz z Klientem — cienki klient urządzenia |
| **Kontrakt** | Zbiór definicji nazw poleceń, zdarzeń, strumieni i ładunków komunikacji Klienta z Aplikacją serwerową |
| **Egzemplarz** | Kopia Oprogramowania udostępniona Licencjobiorcy |
| **Wydanie** | Oznaczona przez Licencjodawcę postać Oprogramowania udostępniana do korzystania |
| **Instalacja** | Zainstalowany i skonfigurowany Egzemplarz zdatny do pracy |

### 2.3. Pojęcia funkcjonalne produktu

| Pojęcie | Znaczenie |
|---|---|
| **Środowisko** | Najwyższy poziom organizacji pracy: TalkIn, WorkSpace, CodeStudio, MultitaskingAI |
| **Moduł** | Wyspecjalizowany obszar roboczy wewnątrz Środowiska |
| **Okno operacyjne** | Przestrzeń wykonywania zadania wewnątrz Modułu |
| **Karta sesji** | Jednostka pracy równoległej z własnym układem, historią i kontekstem |
| **Kanały komunikacji operacyjnej** | Chat Window (Użytkownik ↔ Wykonawca) oraz Execution Loop Window (Koordynator ↔ Wykonawca) |
| **Wykonawca** | Model sztucznej inteligencji, agent albo system wykonawczy realizujący zadania |
| **Koordynator** | Komponent orkiestrujący pracę Wykonawców |
| **Komponent własny** | Nazwany wytwór Operatora: automatyka, agent, projekt albo profil asystenta |
| **Kanał modelu** | Skonfigurowany sposób osiągnięcia modelu: API, CLI, SSH albo HTTP |
| **Punkt izolacji** | Ustawienie rozstrzygające zakres współdzielenia albo odrębności kontekstu i zasobów technicznych |
| **Rozszerzenie** | Wtyczka, umiejętność, konektor albo serwer MCP podłączony do Oprogramowania |

### 2.4. Pojęcia dotyczące danych

| Pojęcie | Znaczenie |
|---|---|
| **Baza** | Trwałe miejsce zapisu stanu Oprogramowania prowadzone przez Aplikację serwerową |
| **Dane Operatora** | Całość danych zapisanych przez Oprogramowanie w związku z pracą Operatora |
| **Treści Operatora** | Treści wprowadzone przez Operatora oraz artefakty powstałe w toku pracy |
| **Poświadczenia** | Klucze dostępu, tokeny, hasła i klucze prywatne wskazywane Oprogramowaniu |
| **Serwer platformy** | Serwer, na którym działa Aplikacja serwerowa |

### 2.5. Reguły wykładni

1. Pojęcia zdefiniowane pisane są wielką literą i mają w całym dokumencie
   jednakowe znaczenie.
2. Odesłania do rozdziałów oznaczają rozdziały niniejszej Licencji, o ile nie
   wskazano inaczej.
3. Wyliczenia poprzedzone zwrotem „w szczególności" nie są wyczerpujące.
4. Zwroty opisowe użyte w dokumentacji produktu („opis funkcjonalny",
   „przewodnik Operatora") nie są nazwami odrębnych dokumentów.

---

## 3. Przedmiot licencji

### 3.1. Zakres przedmiotowy

Przedmiotem Licencji jest korzystanie z Oprogramowania Danaco Console
w postaci wykonywalnej, obejmującego:

1. **Aplikację serwerową** wraz z procesami sesji, silnikiem kolejek, warstwą
   orkiestracji, warstwą dostawcy modelu i warstwą rozszerzeń;
2. **Klienta** wraz z warstwą wizualną (żetony motywu, motyw jasny i ciemny,
   zestaw ikon, kroje pisma);
3. **Powłokę** — natywne okno aplikacji na urządzeniu Operatora;
4. **Kontrakt** wraz z wygenerowanymi z niego wiązaniami warstw;
5. **schemat trwałości** wraz z kompletem migracji;
6. **zasoby wbudowane** — ikony, logotypy i sygnety marki, pliki krojów pisma;
7. **dokumentację produktu** przekazaną Licencjobiorcy wraz z Egzemplarzem
   (rozdz. 3.5).

Przedmiot Licencji obejmuje wyłącznie **postać wykonywalną** Oprogramowania
wraz z dokumentacją produktu. Kod źródłowy, repozytorium projektu, dokumentacja
projektowa i wewnętrzna Licencjodawcy, materiały koncepcyjne, pakiety
projektowe warstwy wizualnej oraz narzędzia budowania **nie stanowią przedmiotu
Licencji** i nie są udostępniane Licencjobiorcy, chyba że co innego wynika
z odrębnej umowy zawartej na piśmie.

### 3.2. Model wdrożenia i postać dystrybucyjna

Oprogramowanie pracuje w modelu hybrydowym: **Aplikacja serwerowa działa na
Serwerze platformy, a Operator korzysta z Oprogramowania przez Pakiet kliencki
zainstalowany na urządzeniu.** Praca obliczeniowa, dostęp do plików, procesy
sesji i wywołania modeli wykonują się po stronie serwera; Pakiet kliencki
odpowiada za prezentację i przyjmowanie poleceń. Stan trwały Oprogramowania
istnieje wyłącznie po stronie serwera.

| Cecha | Aplikacja serwerowa | Pakiet kliencki |
|---|---|---|
| Umiejscowienie | Serwer platformy | Urządzenie Operatora |
| Logika i praca obliczeniowa | Tak | Nie |
| Stan trwały | Tak — jedyne źródło prawdy | Nie |
| Instalacja | Jednorazowa po stronie serwera | Jednorazowa na każde urządzenie |
| Aktualizacja | Centralna, po stronie serwera | Przy zmianie samego okna aplikacji |

Komunikacja Klienta z Aplikacją serwerową odbywa się jednym stałym połączeniem
WebSocket w formacie JSON. Jedno konto właściciela obsługuje wiele urządzeń
Operatora; każde urządzenie uzyskuje własny token dostępu, a odwołanie dostępu
jednemu urządzeniu nie wpływa na pozostałe.

Pakiet kliencki dystrybuowany jest w postaci instalatora dla systemu Windows
oraz pakietów dla systemów Linux; dostęp z urządzeń przenośnych zapewnia
funkcja globalna Mobile po sparowaniu urządzenia.

**Kto prowadzi Serwer platformy**, rozstrzyga odrębna umowa albo zamówienie.
Jeżeli Serwer platformy prowadzi Licencjodawca albo podmiot z nim powiązany,
zastosowanie znajduje rozdz. 9.6 pkt 4.

### 3.3. Elementy wyłączone z przedmiotu licencji

Przedmiotem Licencji **nie są**:

1. **modele sztucznej inteligencji** jakiegokolwiek dostawcy — Oprogramowanie
   nie zawiera modelu, nie dostarcza jego wag ani nie udziela dostępu do usługi
   modelu;
2. **zewnętrzne narzędzia wiersza poleceń** uruchamiane przez Kanał modelu
   rodzaju CLI oraz **usługi sieciowe** wywoływane przez Kanał modelu rodzaju
   API albo HTTP; adres, poświadczenie i warunki korzystania pochodzą od ich
   dostawcy;
3. **hosty zdalne** wskazywane Kanałem modelu rodzaju SSH, **serwery MCP**
   i maszyny wskazane w konfiguracji, wraz z oprogramowaniem na nich
   zainstalowanym;
4. **Komponenty osób trzecich** w zakresie, w jakim ich własne licencje
   przyznają uprawnienia szersze albo nakładają obowiązki odrębne (rozdz. 7);
5. **kod źródłowy** Oprogramowania, repozytorium projektu oraz dokumentacja
   projektowa i wewnętrzna Licencjodawcy;
6. **znaki towarowe, logotypy i oznaczenia** Licencjodawcy poza zakresem
   niezbędnym do zwykłego korzystania z Oprogramowania (rozdz. 6.4).

### 3.4. Zależności zewnętrzne warunkujące działanie

Licencjobiorca przyjmuje do wiadomości, że pełne wykorzystanie Oprogramowania
wymaga zależności pozostających poza zakresem Licencji:

1. **konto u Dostawcy modelu** dla co najmniej jednego Kanału modelu; uzyskanie
   i utrzymanie konta, pokrycie kosztów oraz przestrzeganie limitów obciąża
   Licencjobiorcę;
2. **narzędzia i usługi zewnętrzne** właściwe wybranemu rodzajowi Kanału modelu
   (narzędzie wiersza poleceń, punkt końcowy usługi, host zdalny, interfejs
   webowy dostawcy);
3. **Serwer platformy** wraz z jego utrzymaniem, dostępnością i kopiami
   zapasowymi — w zakresie wynikającym z odrębnej umowy albo zamówienia;
4. **dostęp sieciowy** urządzeń Operatora do Serwera platformy;
5. **środowisko systemowe urządzenia** — system operacyjny z komponentem
   prezentacji wykorzystywanym przez Powłokę. Pakiet kliencki nie wymaga
   instalowania na urządzeniu Operatora żadnych dodatkowych programów ani
   lokalnego układu GPU.

### 3.5. Dokumentacja objęta licencją

Licencją objęta jest dokumentacja produktu przekazana Licencjobiorcy wraz
z Egzemplarzem, w zakresie niezbędnym do korzystania z Oprogramowania. Jest to
zamknięty zestaw pięciu plików umieszczanych w katalogu produktu:

| Plik | Rola |
|---|---|
| `LICENSE.md` | niniejszy dokument — warunki licencyjne |
| `README.md` | karta produktowa: koncepcja, środowiska, moduły, architektura, wymagania |
| `INSTALACJA-I-KONFIGURACJA.md` | instalacja Pakietu klienckiego, pierwsze uruchomienie, konfiguracja platformy, diagnostyka |
| `INSTRUKCJA-UZYTKOWANIA.md` | praca z platformą: nawigacja, oba kanały komunikacji operacyjnej, moduły, scenariusze |
| `CHANGELOG.md` | historia zmian Wydań |

Licencjobiorca może sporządzać kopie dokumentacji na własne potrzeby
i udostępniać ją swoim Operatorom. Publikowanie dokumentacji, jej fragmentów
ani opracowań na jej podstawie poza organizacją Licencjobiorcy wymaga
uprzedniej zgody Licencjodawcy udzielonej na piśmie.

---

## 4. Udzielenie licencji i zakres uprawnień

### 4.1. Udzielenie licencji

Licencjodawca udziela Licencjobiorcy licencji **niewyłącznej**,
**nieprzenoszalnej**, **bez prawa udzielania sublicencji**, na korzystanie
z Oprogramowania w postaci wykonywalnej, wyłącznie na warunkach i w granicach
określonych niniejszym dokumentem.

Licencja obejmuje wyłącznie pola eksploatacji wskazane w rozdz. 4.2.
Uprawnienia niewymienione nie są udzielane; w szczególności nie jest udzielane
prawo zwielokrotniania Oprogramowania w celu wprowadzenia do obrotu, prawo
najmu ani prawo użyczenia Egzemplarza. Licencja nie przenosi autorskich praw
majątkowych ani praw zależnych — prawa te pozostają przy Licencjodawcy.

### 4.2. Uprawnienia przyznane

W ramach Licencji Licencjobiorca jest uprawniony do:

1. **uruchamiania Oprogramowania** na Serwerze platformy oraz **instalowania
   i uruchamiania Pakietu klienckiego** na urządzeniach Operatorów objętych
   zakresem podmiotowym Licencji (rozdz. 4.3), w celach wewnętrznej
   działalności Licencjobiorcy;
2. **korzystania z pełnej funkcjonalności** udostępnionej w danym Wydaniu,
   w granicach opisanych w Licencji i dokumentacji produktu;
3. **konfigurowania** Oprogramowania w zakresie przewidzianym mechanizmami
   produktu: kanałów modeli i kont, komponentów własnych, rozszerzeń,
   punktów izolacji, profili konfiguracji, ustawień obu kanałów komunikacji
   operacyjnej oraz warstw widoczności funkcji;
4. **sporządzania kopii zapasowych** Bazy i danych Oprogramowania — bez
   ograniczeń liczbowych, przy czym Licencjobiorca odpowiada za zabezpieczenie
   tych kopii stosownie do wrażliwości zawartych w nich danych;
5. **sporządzenia kopii zapasowej Pakietu klienckiego** w liczbie niezbędnej
   do zabezpieczenia przed utratą, z zachowaniem oznaczeń praw autorskich;
6. **korzystania z wytworów pracy** — Treści Operatora powstałe przy użyciu
   Oprogramowania należą do Licencjobiorcy w granicach rozdz. 6.5;
7. **wskazywania Operatorów** uprawnionych do obsługi Oprogramowania w imieniu
   Licencjobiorcy, w liczbie wynikającej z rozdz. 4.3.

### 4.3. Zakres podmiotowy

1. Licencja jest liczona **na Operatora**. Jeden Operator prowadzi jedno konto
   właściciela i może korzystać z Oprogramowania na dowolnej liczbie własnych
   urządzeń powiązanych z tym kontem.
2. Liczbę Operatorów objętych Licencją określa odrębna umowa albo zamówienie.
   W braku odmiennego ustalenia Licencja obejmuje **jednego Operatora**.
3. Licencjobiorca prowadzi ewidencję Operatorów i Instalacji na potrzeby
   wykazania zgodności korzystania z zakresem Licencji.
4. Udostępnienie Oprogramowania podmiotowi spoza organizacji Licencjobiorcy —
   w tym świadczenie usług na rzecz osób trzecich przy użyciu Oprogramowania —
   wymaga uprzedniej zgody Licencjodawcy udzielonej na piśmie.
5. Korzystanie przez podmioty powiązane z Licencjobiorcą wymaga objęcia ich
   odrębną umową albo wyraźnego rozszerzenia zakresu Licencji.

### 4.4. Zakres terytorialny i czasowy

Licencja jest udzielana bez ograniczeń terytorialnych, z zastrzeżeniem, że
Licencjobiorca zapewnia zgodność korzystania z przepisami obowiązującymi
w miejscu korzystania, w szczególności z przepisami o ochronie danych
osobowych, o kontroli eksportu oraz z regulacjami dotyczącymi systemów
sztucznej inteligencji. Okres obowiązywania reguluje rozdz. 11.

### 4.5. Warunki finansowe

Wysokość wynagrodzenia, model rozliczenia oraz zasady fakturowania określa
odrębna umowa albo zamówienie. Licencja nie ustala wynagrodzenia i nie może być
odczytywana ani jako zobowiązanie do zapłaty, ani jako zrzeczenie się prawa do
wynagrodzenia. Jeżeli Egzemplarz udostępniono nieodpłatnie, zastosowanie
znajduje rozdz. 10.4 zdanie ostatnie.

### 4.6. Brak przeniesienia praw

Licencja nie stanowi sprzedaży Oprogramowania. Licencjobiorca nabywa wyłącznie
prawo korzystania w zakresie rozdz. 4.2; wszelkie prawa nieudzielone wprost
pozostają zastrzeżone. Licencjobiorca nie może przenieść praw ani obowiązków
wynikających z Licencji na osobę trzecią — w drodze czynności prawnej,
wniesienia aportem, połączenia, podziału ani przejęcia — bez uprzedniej zgody
Licencjodawcy udzielonej na piśmie pod rygorem nieważności.

---

## 5. Zasady korzystania i ograniczenia

### 5.1. Zasada ogólna

Licencjobiorca korzysta z Oprogramowania zgodnie z jego przeznaczeniem,
dokumentacją produktu i niniejszą Licencją, w sposób nienaruszający praw
Licencjodawcy, praw osób trzecich ani przepisów prawa.

### 5.2. Ograniczenia dotyczące Oprogramowania

Licencjobiorca zobowiązuje się **nie**:

1. zwielokrotniać Oprogramowania w zakresie wykraczającym poza czynności
   niezbędne do jego zainstalowania, uruchomienia i sporządzenia kopii
   zapasowej;
2. rozpowszechniać Oprogramowania ani jego części — odpłatnie ani
   nieodpłatnie — w szczególności przez udostępnianie plików wykonywalnych,
   obrazów instalacyjnych, kontenerów ani obrazów maszyn wirtualnych;
3. najmować, użyczać, oddawać w leasing ani udostępniać Oprogramowania osobom
   trzecim w modelu współdzielenia dostępu, hostingu albo świadczenia usług
   przy jego użyciu na rzecz osób trzecich;
4. dokonywać dekompilacji, dezasemblacji ani innej postaci odtwarzania kodu
   źródłowego, z zastrzeżeniem rozdz. 5.6;
5. modyfikować Oprogramowania, tworzyć jego opracowań, tłumaczeń ani adaptacji,
   w tym modyfikować plików wykonywalnych, zasobów wbudowanych, migracji
   schematu ani warstwy wizualnej;
6. usuwać, zasłaniać ani zmieniać oznaczeń praw autorskich, znaków towarowych,
   not licencyjnych ani oznaczeń wersji — w interfejsie Oprogramowania i w
   plikach dystrybucji;
7. obchodzić ani wyłączać mechanizmów uwierzytelniania, autoryzacji, dziennika
   audytu oraz zapisu historii;
8. odwzorowywać chronionych elementów Oprogramowania wskazanych w rozdz. 6.1 —
   w szczególności Kontraktu, modelu danych i warstwy wizualnej — w innym
   produkcie; ograniczenie to nie zakazuje tworzenia produktów tej samej
   kategorii, o ile nie odwzorowują one tych elementów;
9. badać bezpieczeństwa Instalacji innych niż własna Instalacja Licencjobiorcy
   ani przeprowadzać testów obciążeniowych obciążających zasoby współdzielone
   bez uprzedniej zgody Licencjodawcy udzielonej na piśmie. Badanie
   bezpieczeństwa własnej Instalacji jest dozwolone; ustalenia zgłasza się
   trybem rozdz. 13.4.

### 5.3. Zastosowania wyłączone

Licencjobiorca zobowiązuje się **nie stosować** Oprogramowania:

1. w procesach, w których błąd, przerwanie pracy albo utrata danych mogą
   spowodować zagrożenie życia, zdrowia albo bezpieczeństwa osób;
2. w systemach sterowania infrastrukturą krytyczną, urządzeniami medycznymi
   ani środkami transportu;
3. jako jedynego miejsca przechowywania danych o znaczeniu istotnym dla
   działalności Licencjobiorcy — bez niezależnej kopii zapasowej;
4. jako narzędzia podejmującego automatycznie decyzje wywołujące skutki prawne
   wobec osób fizycznych albo w podobny sposób istotnie na nie wpływające, bez
   udziału człowieka.

### 5.4. Korzystanie z modeli i odpowiedzialność za polecenia

1. Oprogramowanie jest narzędziem pośredniczącym: przekazuje modelowi treść
   wskazaną przez Operatora wraz z konfiguracją sesji i roli, po czym zwraca
   strumień odpowiedzi. **Licencjodawca nie tworzy, nie weryfikuje ani nie
   zatwierdza treści wytwarzanych przez model.**
2. Operator odpowiada za treść poleceń wydawanych modelowi oraz za
   wykorzystanie wyników jego pracy, w tym za weryfikację ich poprawności
   przed użyciem.
3. Operator odpowiada za zakres uprawnień nadanych modelowi i agentom: za
   uprawnienia ustalone w Centrum uprawnień, zakres autonomii Wykonawcy,
   punkty obowiązkowego zatwierdzenia, katalogi i zasoby udostępnione
   procesowi sesji oraz konfigurację rozszerzeń i mostów. Nadanie uprawnienia
   zapisu oznacza możliwość dokonywania przez model zmian w objętych nim
   zasobach.
4. Operator zapewnia, że zasoby udostępniane modelowi przysługują
   Licencjobiorcy albo zostały udostępnione za zgodą uprawnionego.

### 5.5. Obowiązki w zakresie bezpieczeństwa

1. Licencjobiorca zabezpiecza Serwer platformy oraz urządzenia Operatorów
   w sposób adekwatny do wrażliwości Danych Operatora — w zakresie, w jakim
   prowadzenie serwera należy do niego.
2. Licencjobiorca odpowiada za ochronę Poświadczeń wskazywanych Oprogramowaniu
   przez odwołania (klucze dostępu, klucze prywatne, katalogi konfiguracji
   kont). Oprogramowanie przechowuje wyłącznie odwołania do Poświadczeń, nie
   ich treść (rozdz. 9.3).
3. Licencjobiorca chroni dane dostępowe do konta właściciela, korzysta
   z dostępnych metod dodatkowych uwierzytelniania i niezwłocznie odwołuje
   dostęp urządzeń utraconych albo wycofanych z użycia.
4. Licencjobiorca nie udostępnia punktu nasłuchu Aplikacji serwerowej w sieci
   publicznej bez zabezpieczenia właściwego dla usług dostępnych publicznie.

### 5.6. Uprawnienia wynikające z przepisów bezwzględnie obowiązujących

Ograniczenia z rozdz. 5.2 pkt 4 i 5 nie naruszają uprawnień przysługujących
Licencjobiorcy na podstawie bezwzględnie obowiązujących przepisów prawa
właściwego, w szczególności uprawnień do zwielokrotniania kodu i tłumaczenia
jego formy w zakresie niezbędnym do uzyskania współdziałania z innymi
programami, na zasadach określonych w art. 75 ustawy z dnia 4 lutego 1994 r.
o prawie autorskim i prawach pokrewnych. Licencjobiorca zamierzający
skorzystać z tych uprawnień zwraca się uprzednio do Licencjodawcy
o udostępnienie informacji niezbędnych do współdziałania.

### 5.7. Zgodność z warunkami Dostawców modeli

Korzystanie z Oprogramowania prowadzące do wywołania modelu podlega równolegle
warunkom umownym Dostawcy modelu. Licencjobiorca zobowiązuje się:

1. przestrzegać regulaminów i limitów Dostawcy modelu, w tym zakazów
   dotyczących współdzielenia kont i obchodzenia limitów;
2. nie wykorzystywać mechanizmu wielu kont i ich rotacji w sposób sprzeczny
   z warunkami Dostawcy modelu; mechanizm ten służy uporządkowanej pracy przy
   wielu uprawnionych kontach;
3. zapewnić, że treści przekazywane modelowi mogą być mu przekazane zgodnie
   z prawem i ze zobowiązaniami Licencjobiorcy wobec osób trzecich.

Licencjodawca nie jest stroną stosunku prawnego między Licencjobiorcą
a Dostawcą modelu i nie odpowiada za jego treść, wykonanie ani skutki
zakończenia.

### 5.8. Skutki naruszenia ograniczeń

Naruszenie ograniczeń z rozdz. 5.2, 5.3 albo 5.7 stanowi istotne naruszenie
Licencji i uprawnia Licencjodawcę do jej wypowiedzenia w trybie rozdz. 11.3,
niezależnie od dalej idących roszczeń. Odpowiedzialność za naruszenie
kształtują zasady ogólne; Licencja nie zastrzega kar umownych.

---

## 6. Własność intelektualna

### 6.1. Prawa do Oprogramowania

Oprogramowanie stanowi utwór w rozumieniu prawa autorskiego. Autorskie prawa
majątkowe przysługują Licencjodawcy — Danaco Holding Group Sp. z o.o.
Autorskie prawa osobiste przysługują Twórcy — Dariuszowi Naharnowiczowi —
i nie podlegają zrzeczeniu ani przeniesieniu.

Ochronie podlega całość Oprogramowania oraz jego elementy dające się
wyodrębnić, w szczególności:

1. kod źródłowy i kod wynikowy wszystkich warstw;
2. **Kontrakt** — dobór, nazewnictwo i struktura poleceń, zdarzeń, strumieni
   i ładunków wraz z konwencją generowania wiązań;
3. **model danych** — struktura schematu trwałości, dobór i nazewnictwo tabel
   i kolumn oraz łańcuch więzów sesji, kart, okien i wiadomości;
4. **warstwa wizualna** — żetony motywu, motyw jasny i ciemny, siatka i zasady
   zestawu ikon, kompozycja widoków;
5. **struktura funkcjonalna produktu** — podział na Środowiska, Moduły i Okna
   operacyjne, model dwóch kanałów komunikacji operacyjnej, model warstw
   widoczności oraz nazwy własne tych elementów;
6. **oznaczenia marki** — logotypy, sygnety i ich warianty;
7. dokumentacja produktu.

### 6.2. Zastrzeżenie praw

Wszelkie prawa nieudzielone Licencjobiorcy wprost w rozdz. 4.2 pozostają
zastrzeżone. Milczenie Licencji co do określonego pola eksploatacji oznacza
brak zgody, nie zgodę dorozumianą.

### 6.3. Zakaz działań naruszających

Licencjobiorca zobowiązuje się nie podejmować działań zmierzających do
podważenia praw Licencjodawcy, w szczególności nie rejestrować na swoją rzecz
oznaczeń zbieżnych z oznaczeniami Oprogramowania, nazw domen zawierających
oznaczenie „Danaco Console" ani nie zgłaszać do ochrony rozwiązań
odwzorowujących chronione elementy Oprogramowania.

### 6.4. Znaki towarowe i oznaczenia

Oznaczenia „Danaco", „Danaco Console", „Danaco Holding Group", nazwy Środowisk
(TalkIn, WorkSpace, CodeStudio, MultitaskingAI) oraz logotypy i sygnety
produktu stanowią oznaczenia Licencjodawcy, chronione prawem autorskim oraz
przepisami o zwalczaniu nieuczciwej konkurencji, niezależnie od rejestracji.

Licencja **nie przyznaje** prawa używania tych oznaczeń poza zakresem
niezbędnym do zwykłego korzystania z Oprogramowania i wewnętrznego odwoływania
się do niego. Użycie oznaczeń w materiałach marketingowych, w komunikacji
publicznej, w nazwach produktów Licencjobiorcy albo w sposób sugerujący
powiązanie gospodarcze z Licencjodawcą wymaga uprzedniej zgody udzielonej na
piśmie.

### 6.5. Prawa do Treści Operatora

1. Licencjodawca nie nabywa praw do Treści Operatora ani do Danych Operatora.
2. Licencjodawca nie uzyskuje dostępu do Danych Operatora w związku z samym
   korzystaniem przez Licencjobiorcę z Oprogramowania.
3. Licencjodawca nie składa oświadczenia co do tego, czy i w jakim zakresie
   wytwór wygenerowany przez model podlega ochronie prawnoautorskiej ani komu
   przysługują do niego prawa; ocena ta zależy od prawa właściwego, sposobu
   powstania wytworu i wkładu twórczego człowieka.
4. Licencjobiorca odpowiada za to, aby korzystanie z wytworów pracy modelu nie
   naruszało praw osób trzecich.

### 6.6. Informacje zwrotne

Uwagi, zgłoszenia błędów i propozycje funkcji przekazane Licencjodawcy przez
Licencjobiorcę mogą być przez Licencjodawcę wykorzystane w rozwoju
Oprogramowania nieodpłatnie i bez ograniczeń czasowych. Przekazanie zgłoszenia
nie rodzi po stronie Licencjodawcy obowiązku jego wdrożenia ani obowiązku
zapłaty. Jeżeli zgłoszenie zawiera rozwiązanie o charakterze twórczym,
Licencjobiorca udziela Licencjodawcy nieodpłatnej, niewyłącznej licencji na
jego wykorzystanie w Oprogramowaniu.

---

## 7. Komponenty osób trzecich i rozszerzenia

### 7.1. Zasada nadrzędna

Oprogramowanie zawiera oraz wykorzystuje komponenty osób trzecich
(„**Komponenty**") rozpowszechniane na własnych licencjach. Do każdego
Komponentu stosuje się **jego własną licencję**; niniejsza Licencja nie
ogranicza uprawnień przyznanych Licencjobiorcy przez te licencje ani nie
rozszerza ich na kod autorski Licencjodawcy. W razie sprzeczności — w zakresie
dotyczącym danego Komponentu — pierwszeństwo ma licencja Komponentu.

Zestawienie Komponentów sporządzono na podstawie deklaracji zależności
i plików zamknięcia zależności bieżącego Wydania. Zestawienie zbiorcze zawiera
Załącznik A.

### 7.2. Komponenty warstwy serwerowej

Aplikacja serwerowa jest budowana jako moduł Go z 23 zależnościami
bezpośrednimi i 86 zależnościami pośrednimi deklarowanymi w pliku modułu.
Zależności bezpośrednie wraz z rodzajem licencji ustalonym z treści plików
licencyjnych tych komponentów:

| Komponent | Wersja | Licencja | Rola w produkcie |
|---|---|---|---|
| `github.com/coder/websocket` | v1.8.15 | ISC | transport połączenia klient–serwer |
| `modernc.org/sqlite` | v1.56.0 | BSD 3-Clause | silnik trwałości bez zależności natywnych |
| `golang.org/x/crypto` | v0.54.0 | BSD 3-Clause | funkcje kryptograficzne, tor połączeń SSH |
| `golang.org/x/sys` | v0.47.0 | BSD 3-Clause | dostęp do funkcji systemowych |
| `golang.org/x/text` | v0.40.0 | BSD 3-Clause | przetwarzanie tekstu i kodowań |
| `golang.org/x/image` | v0.44.0 | BSD 3-Clause | przetwarzanie obrazu |
| `github.com/jackc/pgx/v5` | v5.10.0 | MIT | obsługa baz PostgreSQL |
| `github.com/go-sql-driver/mysql` | v1.10.0 | **MPL-2.0** | obsługa baz MySQL |
| `github.com/docker/docker` | v28.5.2 | Apache-2.0 | obsługa silnika kontenerów |
| `github.com/docker/go-connections` | v0.8.1 | Apache-2.0 | połączenia silnika kontenerów |
| `github.com/go-git/go-git/v5` | v5.19.2 | Apache-2.0 | operacje kontroli wersji |
| `github.com/hexops/gotextdiff` | v1.0.3 | BSD 3-Clause | porównywanie treści |
| `github.com/pdfcpu/pdfcpu` | v0.15.0 | Apache-2.0 | przetwarzanie dokumentów PDF |
| `github.com/tdewolff/canvas` | (wersja Wydania) | MIT | skład i rysowanie dokumentów |
| `github.com/disintegration/imaging` | v1.6.2 | MIT | przekształcenia obrazu |
| `github.com/HugoSmits86/nativewebp` | v1.3.0 | MIT | obsługa formatu WebP |
| `github.com/rwcarlsen/goexif` | (wersja Wydania) | BSD 2-Clause | odczyt metadanych obrazu |
| `github.com/lucasb-eyer/go-colorful` | v1.4.0 | MIT | operacje na barwach |
| `github.com/emersion/go-imap/v2` | v2.0.0-beta.8 | MIT | obsługa poczty przychodzącej |
| `github.com/emersion/go-message` | v0.18.2 | MIT | przetwarzanie wiadomości poczty |
| `github.com/getkin/kin-openapi` | v0.146.0 | MIT | obsługa opisów interfejsów zewnętrznych |
| `github.com/tiktoken-go/tokenizer` | v0.8.1 | MIT | pomiar długości treści przekazywanej modelowi |
| `go.yaml.in/yaml/v3` | v3.0.5 | MIT | przetwarzanie plików konfiguracyjnych |

Wymienione licencje — poza wskazaną wprost MPL-2.0 — są licencjami
zezwalającymi, dopuszczającymi rozpowszechnianie w postaci wynikowej pod
warunkiem zachowania not o prawach autorskich i tekstu licencji.
Komponent na licencji MPL-2.0 jest wykorzystywany jako biblioteka
w postaci niezmodyfikowanej; obowiązek udostępnienia kodu źródłowego dotyczy
wyłącznie tego komponentu i jest spełniany przez wskazanie jego publicznego
źródła.

### 7.3. Komponenty warstwy Klienta

Zależnością **produkcyjną** Klienta — to jest taką, której kod trafia do
pakietu dostarczanego Licencjobiorcy — jest wyłącznie:

| Komponent | Wersja | Licencja | Rola w produkcie |
|---|---|---|---|
| `@tauri-apps/api` | 2.5.0 | Apache-2.0 **lub** MIT (do wyboru) | wywołanie poleceń Powłoki z poziomu interfejsu |

Pozostałe zależności warstwy Klienta są narzędziami budowania i testowania
i nie wchodzą do pakietu dostarczanego Licencjobiorcy: `typescript` (5.8.3,
Apache-2.0), `vite` (6.3.5, MIT), `vitest` (3.2.7, MIT), `jsdom` (MIT),
`happy-dom` (MIT), `playwright` (Apache-2.0).

Zamknięcie zależności warstwy Klienta obejmuje 150 pozycji. Rozkład rodzajów
licencji odczytany z pliku zamknięcia: MIT — 133, Apache-2.0 — 5, ISC — 3,
MIT-0 — 2, BSD-2-Clause — 2, BSD-3-Clause — 2, „Apache-2.0 lub MIT" — 1,
BlueOak-1.0.0 — 1, CC0-1.0 — 1. Zestawienie nie ujawnia licencji o charakterze
wzajemnym w tej warstwie.

### 7.4. Komponenty warstwy Powłoki

Zależności bezpośrednie Powłoki:

| Komponent | Wersja | Licencja | Rola w produkcie |
|---|---|---|---|
| `tauri` | 2.11.5 | Apache-2.0 **lub** MIT | powłoka okna, ikona zasobnika, obsługa obrazów ikon |
| `tauri-plugin-dialog` | 2.7.2 | Apache-2.0 **lub** MIT | okna dialogowe systemu |
| `serde` | 1.0.229 | MIT **lub** Apache-2.0 | serializacja danych |
| `serde_json` | 1.0.151 | MIT **lub** Apache-2.0 | obsługa formatu JSON |
| `ureq` | 3.4.0 | MIT **lub** Apache-2.0 | żądania HTTP warstwy okna |
| `sha2` | 0.10.9 | MIT **lub** Apache-2.0 | sumy kontrolne |
| `tauri-build` | 2.6.3 | Apache-2.0 **lub** MIT | narzędzie etapu budowania (nieobecne w pliku wynikowym) |

Zamknięcie zależności tej warstwy obejmuje 407 pozycji, w tym pozycje na
licencjach wzajemnych MPL-2.0 wykorzystywane jako biblioteki w postaci
niezmodyfikowanej. Pełne zestawienie wraz z tekstami licencji dołączane jest
do Wydania (rozdz. 7.8).

### 7.5. Kroje pisma

Warstwa wizualna opiera się na trzech krojach pisma, których pliki
w formacie WOFF2 są wbudowane w Oprogramowanie i rozpowszechniane wraz z nim:

| Krój | Rola w produkcie | Licencja | Uprawniony |
|---|---|---|---|
| **Space Grotesk** | nagłówki, tytuły środowisk, logotyp | SIL Open Font License 1.1 | Copyright 2020 The Space Grotesk Project Authors |
| **IBM Plex Sans** | interfejs i treść | SIL Open Font License 1.1 | Copyright 2019 IBM Corp. |
| **IBM Plex Mono** | dane techniczne, identyfikatory, terminal | SIL Open Font License 1.1 | Copyright 2019 IBM Corp. |

Kroje osadzono w podzbiorach `latin` oraz `latin-ext`, obejmujących pełny
zestaw polskich znaków diakrytycznych. Pełne teksty licencji dołączono do
zasobów Oprogramowania. Z licencji wynikają obowiązki, które Licencjobiorca
przyjmuje do wiadomości:

1. pliki krojów mogą być rozpowszechniane wyłącznie **wraz z tekstem licencji
   i notą o prawach autorskich**;
2. **zakazana jest sprzedaż samych plików krojów** jako odrębnego towaru;
3. zmodyfikowana wersja kroju nie może być rozpowszechniana pod nazwą
   zastrzeżoną przez autora oryginału;
4. licencja dotyczy plików krojów, nie dokumentów ani grafik złożonych przy ich
   użyciu — wytwory Operatora nie stają się przez to objęte tą licencją.

### 7.6. Zestaw ikon i zasoby marki

Zestaw ikon interfejsu obejmuje **82 pozycje** na jednolitej siatce 24×24
z barwą dziedziczoną z kontekstu:

1. **78 pozycji** wywiedzionych z biblioteki **Lucide**, rozpowszechnianej na
   licencji **ISC**, z obrysem ujednoliconym do wartości przyjętej w zestawie;
2. **4 pozycje własne** — emblematy Środowisk (TalkIn, WorkSpace, CodeStudio,
   MultitaskingAI) wykonane w siatce i kresce zestawu; stanowią utwór
   Licencjodawcy.

Warunkiem licencji ISC jest zachowanie noty o prawach autorskich i treści
zezwolenia we wszystkich kopiach; tekst licencji dołączono do zasobów
Oprogramowania.

**Zasoby marki** — logotypy, sygnety i ich warianty dla motywu jasnego
i ciemnego oraz znak marki — stanowią utwór Licencjodawcy, nie podlegają
licencjom Komponentów, a ich użycie reguluje rozdz. 6.4.

### 7.7. Rozszerzenia

Oprogramowanie przyjmuje rozszerzenia dwóch źródeł: dostarczane z platformą
oraz instalowane przez Operatora. Licencjodawca odpowiada wyłącznie za
rozszerzenia dostarczone z platformą. Do rozszerzeń instalowanych przez
Operatora stosuje się zasady następujące:

1. Licencjobiorca odpowiada za tytuł prawny do rozszerzenia i za zgodność jego
   użycia z licencją jego dostawcy;
2. rozszerzenie działa w zakresie uprawnień nadanych mu przy podłączeniu;
   nadanie uprawnień jest czynnością Operatora i jego odpowiedzialnością;
3. Licencjodawca nie odpowiada za działanie rozszerzeń niepochodzących od
   niego ani za skutki ich podłączenia.

### 7.8. Obowiązki dokumentacyjne przy dystrybucji

Licencjodawca dołącza do Wydania zestawienie Komponentów wraz z pełnymi
tekstami ich licencji i notami o prawach autorskich. Licencjobiorca
zobowiązuje się nie usuwać tych plików z Instalacji ani nie modyfikować ich
treści.

---

## 8. Aktualizacje i wersjonowanie

### 8.1. Wersjonowanie

Wydania oznaczane są przez Licencjodawcę. Numer wersji produktu identyfikuje
zakres funkcjonalny opisany dokumentacją produktu dostarczoną wraz z Wydaniem;
identyfikacja Wydania następuje przez oznaczenie nadane przez Licencjodawcę
oraz zapis w pliku `CHANGELOG.md`.

### 8.2. Charakter Aktualizacji

1. Aktualizacja może obejmować poprawki błędów, zmiany zakresu funkcjonalnego,
   zmiany interfejsu, zmiany Kontraktu oraz zmiany schematu trwałości.
2. Aktualizacja może usunąć funkcję dostępną w Wydaniu poprzednim, jeżeli
   funkcja ta została zastąpiona rozwiązaniem równoważnym albo jej utrzymanie
   byłoby niezgodne z architekturą produktu; informację o usunięciu podaje
   `CHANGELOG.md`.
3. Aktualizacja Aplikacji serwerowej wykonywana jest po stronie Serwera
   platformy i nie wymaga czynności Operatora. Aktualizacja Pakietu
   klienckiego wykonywana jest przez Licencjobiorcę wtedy, gdy zmienia się
   samo okno aplikacji.

### 8.3. Migracje schematu trwałości

1. Schemat Bazy tworzą i aktualizują migracje stosowane samoczynnie przy
   starcie Aplikacji serwerowej, w jednej transakcji, z kontrolą sumy
   kontrolnej kroku.
2. **Migracje działają wyłącznie w kierunku naprzód.** Oprogramowanie nie
   udostępnia mechanizmu cofnięcia migracji; po uruchomieniu nowszego Wydania
   na istniejącej Bazie powrót do Wydania wcześniejszego może być niemożliwy.
3. Przed instalacją Aktualizacji wykonuje się kopię zapasową Bazy. Zaniechanie
   tej czynności obciąża podmiot prowadzący Serwer platformy.

### 8.4. Uprawnienie do Aktualizacji i wsparcie

1. W okresie obowiązywania Licencji Licencjobiorcy przysługuje prawo
   korzystania z Aktualizacji udostępnianych przez Licencjodawcę dla posiadanej
   wersji produktu.
2. Licencjodawca nie zobowiązuje się do wydawania Aktualizacji w oznaczonych
   terminach ani do wdrożenia funkcji nieopisanych dokumentacją produktu.
3. Zakres, kanał i czasy reakcji wsparcia technicznego określa odrębna umowa
   albo zamówienie.

### 8.5. Zmiany treści Licencji

Licencjodawca może wydać nową wersję dokumentu Licencji wraz z nowym Wydaniem
Oprogramowania. Nowa wersja obowiązuje wobec Wydań udostępnionych po jej
wydaniu; do Wydania już posiadanego stosuje się wersję Licencji dostarczoną
wraz z tym Wydaniem. Zmiana Licencji nie może pozbawić Licencjobiorcy
uprawnień nabytych w stosunku do Wydania już posiadanego.

---

## 9. Dane i treści Operatora

### 9.1. Miejsce przechowywania danych

Stan trwały Oprogramowania — sesje, karty sesji, okna komunikacji, historia
wiadomości, pamięć kontekstowa, artefakty, ustawienia, definicje kanałów
modeli, kont, komponentów własnych i rozszerzeń — mieści się w Bazie
prowadzonej przez Aplikację serwerową na Serwerze platformy. Pakiet kliencki
nie przechowuje trwale Danych Operatora poza danymi niezbędnymi do
nawiązania połączenia i uwierzytelnienia urządzenia.

**Licencjodawca nie otrzymuje kopii Danych Operatora w związku z korzystaniem
z Oprogramowania**, o ile nie prowadzi Serwera platformy na podstawie odrębnej
umowy (rozdz. 9.6 pkt 4).

### 9.2. Zakres danych zapisywanych przez Oprogramowanie

W Bazie zapisywane są w szczególności:

1. konto właściciela, wykaz urządzeń i stan ich poświadczeń;
2. sesje, karty sesji i okna komunikacji wraz z ich ustawieniami;
3. **treść wiadomości Operatora oraz treść odpowiedzi modelu**;
4. pamięć kontekstowa, artefakty i ich wersje;
5. wartości ustawień konfiguracji na wszystkich warstwach i poziomach zasięgu;
6. definicje kanałów modeli, kont (bez treści Poświadczeń), komponentów
   własnych, rozszerzeń i punktów izolacji;
7. zapisy dziennika audytu dotyczące decyzji podejmowanych w kanałach
   komunikacji operacyjnej.

### 9.3. Postępowanie z Poświadczeniami

Oprogramowanie **nie przechowuje treści Poświadczeń w Bazie.** Baza przechowuje
wyłącznie odwołania — nazwę wpisu w magazynie danych dostępowych, ścieżkę
katalogu konfiguracji konta albo ścieżkę klucza. Kontrakt przyjmuje treść
poświadczenia wyłącznie w żądaniach dodania i aktualizacji; żadna odpowiedź ani
żadne zdarzenie nie zwraca tej treści. Odpowiedzialność za zabezpieczenie
miejsc, w których Poświadczenia się znajdują, spoczywa na Licencjobiorcy.

### 9.4. Przekazywanie danych na zewnątrz

Oprogramowanie przekazuje dane poza Serwer platformy wyłącznie w następstwie
czynności Operatora i wyłącznie:

1. **przy wywołaniu modelu** — treścią zapytania, zgodnie z rodzajem
   skonfigurowanego Kanału modelu (API, CLI, SSH albo HTTP);
2. **przy pracy rozszerzeń i mostów** — w zakresie wynikającym z konfiguracji
   punktów dostępu i nadanych uprawnień;
3. **przy powiadomieniach** — w zakresie i kanałem wskazanym ustawieniami
   powiadomień;
4. **do urządzeń Operatora** powiązanych z kontem — w ramach synchronizacji
   stanu platformy.

Poza powyższymi przypadkami Oprogramowanie nie nawiązuje połączeń
wychodzących. **Oprogramowanie nie zawiera telemetrii przekazywanej
Licencjodawcy, nie raportuje zdarzeń użycia i nie wysyła samoczynnie zgłoszeń
o błędach.**

Operator przyjmuje do wiadomości, że treść przekazana modelowi opuszcza
kontrolę Oprogramowania i podlega dalej zasadom przetwarzania stosowanym przez
Dostawcę modelu. Ocena dopuszczalności przekazania określonej kategorii danych
należy do Licencjobiorcy i jest dokonywana **przed** przekazaniem.

### 9.5. Kopie zapasowe i usuwanie danych

1. Kopie zapasowe Bazy wykonuje podmiot prowadzący Serwer platformy.
   Oprogramowanie nie wykonuje kopii zapasowych samoczynnie i nie zawiera
   mechanizmu przywracania danych z kopii.
2. Usunięcie Bazy powoduje nieodwracalną utratę całości stanu platformy.
3. Odinstalowanie Pakietu klienckiego z urządzenia nie usuwa Danych Operatora
   — stan trwały pozostaje na Serwerze platformy.
4. Dostęp do danych poza Oprogramowaniem możliwy jest przez odczyt Bazy
   narzędziami zgodnymi z jej formatem.

### 9.6. Dane osobowe

1. Jeżeli Dane Operatora obejmują dane osobowe, ich administratorem jest
   Licencjobiorca.
2. Licencjobiorca zobowiązuje się przetwarzać dane osobowe przy użyciu
   Oprogramowania zgodnie z przepisami, w szczególności z rozporządzeniem (UE)
   2016/679.
3. Przed przekazaniem danych osobowych modelowi Licencjobiorca ocenia podstawę
   prawną przekazania oraz warunki przetwarzania stosowane przez Dostawcę
   modelu, w tym miejsce przetwarzania i ewentualne przekazanie poza
   Europejski Obszar Gospodarczy.
4. **Jeżeli Serwer platformy prowadzi Licencjodawca albo podmiot z nim
   powiązany**, Licencjodawca przetwarza dane osobowe zawarte w Bazie jako
   podmiot przetwarzający, na polecenie Licencjobiorcy i wyłącznie na
   podstawie umowy powierzenia przetwarzania danych osobowych zawartej
   z Licencjobiorcą. Umowa ta określa lokalizację przetwarzania, okres
   retencji, zasady usuwania danych, środki bezpieczeństwa oraz zasady
   korzystania z dalszych podmiotów przetwarzających.

### 9.7. Dostęp Licencjodawcy do danych przy czynnościach wsparcia

Przekazanie Licencjodawcy plików diagnostycznych, wycinków Bazy albo zrzutów
interfejsu następuje wyłącznie z inicjatywy Licencjobiorcy i na jego
odpowiedzialność. Licencjobiorca usuwa z przekazywanych materiałów dane,
których przekazanie byłoby niedopuszczalne. Licencjodawca wykorzystuje
przekazane materiały wyłącznie w celu udzielenia wsparcia i usuwa je po jego
zakończeniu.

---

## 10. Rękojmia i odpowiedzialność

### 10.1. Charakter świadczenia

Oprogramowanie jest udostępniane w stanie odpowiadającym opisowi zawartemu
w dokumentacji produktu dostarczonej wraz z Wydaniem. Licencjodawca oświadcza,
że przysługują mu prawa pozwalające na udzielenie Licencji w zakresie
opisanym w rozdz. 4.

### 10.2. Zakres rękojmi

1. Licencjodawca nie udziela gwarancji, że Oprogramowanie będzie działać bez
   przerw i bez błędów, ani że będzie odpowiadać potrzebom Licencjobiorcy
   wykraczającym poza opis zawarty w dokumentacji produktu.
2. Licencjodawca nie odpowiada za treści wytworzone przez model ani za ich
   poprawność, aktualność i kompletność (rozdz. 5.4).
3. Licencjodawca nie odpowiada za dostępność, warunki ani zmiany usług
   Dostawców modeli i innych podmiotów zewnętrznych.
4. Uprawnienia Licencjobiorcy wynikające z przepisów bezwzględnie
   obowiązujących pozostają nienaruszone.

### 10.3. Zgłaszanie nieprawidłowości

Nieprawidłowość działania Oprogramowania Licencjobiorca zgłasza Licencjodawcy
niezwłocznie po jej wykryciu, ze wskazaniem okoliczności umożliwiających jej
odtworzenie. Tryb i terminy usuwania nieprawidłowości określa odrębna umowa
albo zamówienie (rozdz. 8.4 pkt 3).

### 10.4. Ograniczenie odpowiedzialności

W najszerszym zakresie dopuszczalnym przez prawo właściwe Licencjodawca nie
ponosi odpowiedzialności za:

1. utracone korzyści, utratę przychodu, utratę spodziewanych oszczędności,
   utratę renomy ani szkody pośrednie;
2. utratę, uszkodzenie albo ujawnienie Danych Operatora oraz Treści Operatora,
   w zakresie, w jakim nie wynikają one z niewykonania obowiązków
   Licencjodawcy jako podmiotu prowadzącego Serwer platformy;
3. skutki decyzji podjętych na podstawie wyników pracy modelu;
4. działania modelu i agentów w zakresie zasobów udostępnionych im przez
   Operatora, w tym za zmiany dokonane w plikach objętych nadanym
   uprawnieniem zapisu;
5. przerwy w działaniu albo zmianę warunków świadczenia usług przez Dostawcę
   modelu, w tym wyczerpanie limitów i zawieszenie konta;
6. niezgodność korzystania z Oprogramowania z warunkami Dostawcy modelu albo
   z przepisami obowiązującymi Licencjobiorcę;
7. szkody wynikłe z korzystania z Oprogramowania w sposób sprzeczny
   z rozdz. 5.3.

Całkowita odpowiedzialność Licencjodawcy z tytułu Licencji ograniczona jest do
wysokości wynagrodzenia zapłaconego przez Licencjobiorcę w okresie dwunastu
miesięcy poprzedzających zdarzenie szkodzące. Jeżeli Egzemplarz udostępniono
nieodpłatnie, odpowiedzialność Licencjodawcy ograniczona jest do przypadków
winy umyślnej. Ograniczenia nie dotyczą szkody wyrządzonej umyślnie ani
odpowiedzialności, której wyłączyć nie można na podstawie przepisu
bezwzględnie obowiązującego.

### 10.5. Siła wyższa

Żadna ze stron nie odpowiada za niewykonanie zobowiązań spowodowane
okolicznościami pozostającymi poza jej rozsądną kontrolą, w szczególności
działaniem siły wyższej, awarią infrastruktury teleinformatycznej niezależną
od strony, decyzją organu władzy publicznej ani zaprzestaniem świadczenia
usług przez Dostawcę modelu.

---

## 11. Okres obowiązywania i zakończenie licencji

### 11.1. Wejście w życie

Licencja wiąże z chwilą pierwszego z następujących zdarzeń: zainstalowania
Pakietu klienckiego, pierwszego uruchomienia Oprogramowania albo
zaakceptowania warunków Licencji w toku instalacji. Przystąpienie do
korzystania z Oprogramowania oznacza przyjęcie warunków Licencji w całości.
Podmiot, który warunków nie akceptuje, odstępuje od instalacji, a Egzemplarz
usuwa.

### 11.2. Okres obowiązywania

Licencja udzielana jest na **czas nieoznaczony**, o ile odrębna umowa albo
zamówienie nie wskazuje okresu oznaczonego. Każda ze stron może wypowiedzieć
Licencję udzieloną na czas nieoznaczony z zachowaniem **trzydziestodniowego**
okresu wypowiedzenia, ze skutkiem na koniec miesiąca kalendarzowego.

### 11.3. Wypowiedzenie z przyczyny naruszenia

1. Licencjodawca może wypowiedzieć Licencję ze skutkiem natychmiastowym,
   jeżeli Licencjobiorca istotnie narusza jej postanowienia, w szczególności
   ograniczenia z rozdz. 5.2, 5.3 albo 5.7, i nie zaprzestaje naruszenia
   w terminie **trzydziestu dni** od wezwania skierowanego na piśmie albo
   drogą elektroniczną.
2. Wezwania nie wymaga wypowiedzenie z powodu naruszenia nieusuwalnego —
   w szczególności rozpowszechnienia Oprogramowania, usunięcia oznaczeń praw
   autorskich albo obejścia mechanizmów uwierzytelniania i dziennika audytu.
3. Licencjobiorca może w każdym czasie zakończyć Licencję przez zaprzestanie
   korzystania z Oprogramowania i usunięcie Egzemplarzy.
4. Strony mogą zakończyć Licencję w każdym czasie za porozumieniem wyrażonym
   na piśmie.

### 11.4. Skutki zakończenia

Z chwilą zakończenia Licencji:

1. wygasają uprawnienia Licencjobiorcy opisane w rozdz. 4.2;
2. Licencjobiorca zaprzestaje korzystania z Oprogramowania, odinstalowuje
   Pakiet kliencki ze wszystkich urządzeń i usuwa Egzemplarze oraz ich kopie
   zapasowe;
3. Licencjobiorca zachowuje prawo do Danych Operatora i Treści Operatora oraz
   prawo do zachowania kopii Bazy — zakończenie Licencji **nie zobowiązuje do
   usunięcia własnych danych**. Jeżeli Serwer platformy prowadzi Licencjodawca,
   wydanie kopii danych i termin ich usunięcia określa umowa powierzenia
   przetwarzania (rozdz. 9.6 pkt 4);
4. Licencjobiorca zachowuje prawo korzystania z wytworów pracy powstałych przy
   użyciu Oprogramowania w okresie obowiązywania Licencji;
5. na żądanie Licencjodawcy Licencjobiorca potwierdza wykonanie obowiązków
   z pkt 2 oświadczeniem złożonym na piśmie albo drogą elektroniczną.

### 11.5. Postanowienia trwające

Zakończenie Licencji nie wpływa na moc obowiązującą postanowień, które
z natury obowiązują dłużej, w szczególności rozdz. 6 (własność
intelektualna), rozdz. 10 (odpowiedzialność), rozdz. 13 (poufność) oraz
rozdz. 14 (postanowienia końcowe).

---

## 12. Zgodność licencyjna

### 12.1. Obowiązek zgodności

Licencjobiorca korzysta z Oprogramowania w granicach zakresu podmiotowego
Licencji (rozdz. 4.3) i utrzymuje ewidencję pozwalającą wykazać zgodność
korzystania z tym zakresem.

### 12.2. Weryfikacja

Licencjodawcy przysługuje prawo żądania oświadczenia o liczbie Operatorów
i Instalacji, składanego nie częściej niż raz w roku, w terminie czternastu dni
od żądania. **Oprogramowanie nie zawiera mechanizmu technicznej kontroli
licencji**: nie weryfikuje klucza licencyjnego, nie łączy się z serwerem
aktywacji i nie raportuje faktu uruchomienia. Weryfikacja zgodności odbywa się
wyłącznie środkami organizacyjnymi.

### 12.3. Skutki stwierdzenia niezgodności

W razie stwierdzenia korzystania wykraczającego poza zakres Licencji strony
ustalają zakres niezgodności i sposób jej usunięcia. Nieusunięcie niezgodności
w uzgodnionym terminie stanowi istotne naruszenie Licencji w rozumieniu
rozdz. 11.3 pkt 1.

---

## 13. Poufność

### 13.1. Informacje poufne

Za informacje poufne uważa się nieujawnione publicznie informacje dotyczące
Oprogramowania, przekazane Licencjobiorcy przez Licencjodawcę albo powzięte
w związku z korzystaniem z Oprogramowania, w szczególności dokumentację
techniczną nieprzeznaczoną do publikacji, opisy architektury wewnętrznej,
strukturę Kontraktu i modelu danych, informacje o podatnościach oraz
informacje o planach rozwoju produktu.

Za informacje poufne **nie uważa się** informacji, które: są publicznie
dostępne bez naruszenia zobowiązania do poufności; były znane Licencjobiorcy
przed ich przekazaniem; zostały uzyskane zgodnie z prawem od osoby trzeciej
nieobjętej zobowiązaniem do poufności; albo zostały opracowane samodzielnie
bez wykorzystania informacji poufnych.

### 13.2. Obowiązki

Licencjobiorca zachowuje informacje poufne w tajemnicy, nie ujawnia ich osobom
trzecim bez zgody Licencjodawcy udzielonej na piśmie oraz udostępnia je
własnym Operatorom i współpracownikom wyłącznie w zakresie niezbędnym do
korzystania z Oprogramowania, po zobowiązaniu ich do poufności.

Obowiązek zachowania poufności nie stoi na przeszkodzie ujawnieniu informacji
na żądanie uprawnionego organu, w zakresie i trybie wynikającym z przepisów;
o takim żądaniu Licencjobiorca niezwłocznie zawiadamia Licencjodawcę, o ile
nie jest to zabronione.

### 13.3. Okres obowiązywania

Zobowiązanie do poufności obowiązuje w okresie trwania Licencji oraz przez
**trzy lata** po jej zakończeniu. W odniesieniu do informacji stanowiących
tajemnicę przedsiębiorstwa zobowiązanie obowiązuje przez cały okres, w którym
informacje zachowują taki charakter, bez ograniczenia terminem.

### 13.4. Ujawnienie podatności

Licencjobiorca, który poweźmie informację o podatności bezpieczeństwa
Oprogramowania, zawiadamia o niej wyłącznie Licencjodawcę — na adres
support@danaco-group.pl (rozdz. 14.7) — i powstrzymuje się od jej publikacji
przez okres uzgodniony z Licencjodawcą, nie dłuższy jednak niż konieczny do
wydania poprawki.

---

## 14. Postanowienia końcowe

### 14.1. Prawo właściwe

Licencja podlega prawu polskiemu, z wyłączeniem norm kolizyjnych. Wybór prawa
nie pozbawia Licencjobiorcy ochrony wynikającej z bezwzględnie obowiązujących
przepisów prawa państwa jego zwykłego pobytu, jeżeli Licencjobiorcą jest
konsument.

### 14.2. Rozstrzyganie sporów

Strony podejmą w dobrej wierze próbę polubownego rozwiązania sporu przed
skierowaniem sprawy na drogę postępowania sądowego. Spory nierozwiązane
polubownie rozstrzyga sąd powszechny właściwy miejscowo dla siedziby
Licencjodawcy; postanowienie to nie ma zastosowania, jeżeli Licencjobiorcą
jest konsument, a właściwość wyłączna byłaby wobec niego bezskuteczna.

### 14.3. Forma zmian

Zmiany Licencji w stosunku dwustronnym wymagają formy pisemnej albo
dokumentowej pod rygorem nieważności, z zastrzeżeniem rozdz. 8.5.

### 14.4. Klauzula salwatoryjna

Nieważność albo bezskuteczność któregokolwiek z postanowień Licencji nie
wpływa na ważność pozostałych. W miejsce postanowienia nieważnego strony
stosują postanowienie skuteczne, najbliższe celowi gospodarczemu
postanowienia zastępowanego. Dotyczy to w szczególności rozdz. 5.2, 10.2
i 10.4, których zakres podlega ograniczeniu do granic dopuszczalnych przez
prawo właściwe.

### 14.5. Całość porozumienia

Licencja wraz z Załącznikiem A stanowi całość porozumienia stron w zakresie
korzystania z Oprogramowania i zastępuje wcześniejsze ustalenia, oświadczenia
i zapewnienia dotyczące tego przedmiotu, z zastrzeżeniem rozdz. 1.3 pkt 1.
Materiały prezentacyjne i koncepcyjne nie stanowią części Licencji i nie
tworzą zobowiązania co do zakresu funkcjonalnego.

### 14.6. Język dokumentu

Językiem Licencji jest język polski. W razie sporządzenia tłumaczenia
rozstrzygające znaczenie ma wersja polska.

### 14.7. Doręczenia

Oświadczenia i zawiadomienia związane z Licencją składa się na następujące
adresy Licencjodawcy:

| Kanał | Adres |
|---|---|
| Adres do doręczeń pisemnych | Danaco Holding Group Sp. z o.o., ul. Gen. J. Hallera 81, 43-400 Cieszyn |
| Adres poczty elektronicznej | support@danaco-group.pl |

Tym samym adresem poczty elektronicznej Licencjobiorca posługuje się przy
zgłoszeniach nieprawidłowości (rozdz. 10.3) oraz przy zgłoszeniach podatności
bezpieczeństwa (rozdz. 13.4). Adres do doręczeń Licencjobiorcy wskazuje odrębna
umowa albo zamówienie. Zmianę adresu strona zgłasza drugiej stronie
niezwłocznie; do czasu zgłoszenia doręczenie na adres dotychczasowy jest
skuteczne.

### 14.8. Załącznik

Załącznik A stanowi integralną część Licencji. W razie rozbieżności między
treścią rozdziału a treścią Załącznika rozstrzyga treść rozdziału, z wyjątkiem
zestawienia Komponentów, gdzie rozstrzyga Załącznik jako zestawienie
szczegółowe.

---

## 15. Załącznik A — komponenty osób trzecich

Zestawienie obejmuje komponenty osób trzecich wykorzystane w bieżącym Wydaniu,
w podziale na warstwy Oprogramowania. Pełne teksty licencji oraz noty o prawach
autorskich, w tym dla zależności pośrednich, dołączane są do Wydania
(rozdz. 7.8).

### A.1. Warstwa serwerowa

23 zależności bezpośrednie wymienione w rozdz. 7.2 oraz 86 zależności
pośrednich deklarowanych w pliku modułu. Rodzaje licencji zależności
bezpośrednich: BSD 3-Clause, MIT, Apache-2.0, ISC, BSD 2-Clause oraz jedna
licencja wzajemna MPL-2.0 (`github.com/go-sql-driver/mysql`).

### A.2. Warstwa Klienta

Jedna zależność produkcyjna: `@tauri-apps/api` 2.5.0 (Apache-2.0 lub MIT).
Zamknięcie zależności warstwy: 150 pozycji, rozkład licencji wskazany
w rozdz. 7.3. Narzędzia budowania i testowania nie wchodzą do pakietu
dostarczanego Licencjobiorcy.

### A.3. Warstwa Powłoki

Siedem zależności bezpośrednich wymienionych w rozdz. 7.4 (Apache-2.0 lub MIT,
MIT lub Apache-2.0). Zamknięcie zależności warstwy: 407 pozycji, w tym pozycje
na licencjach wzajemnych MPL-2.0 wykorzystywane jako biblioteki w postaci
niezmodyfikowanej.

### A.4. Kroje pisma

Space Grotesk, IBM Plex Sans, IBM Plex Mono — SIL Open Font License 1.1;
podzbiory `latin` i `latin-ext`. Obowiązki licencyjne wskazuje rozdz. 7.5.

### A.5. Zestaw ikon

82 pozycje zestawu: 78 wywiedzionych z biblioteki Lucide (ISC), 4 emblematy
Środowisk stanowiące utwór Licencjodawcy. Zasoby marki nie należą do zestawu
ikon i podlegają rozdz. 6.4.

### A.6. Zależności zewnętrzne nieobjęte dystrybucją

Modele sztucznej inteligencji, narzędzia wiersza poleceń dostawców, usługi
sieciowe dostawców, hosty zdalne, serwery MCP oraz oprogramowanie Serwera
platformy niebędące częścią Oprogramowania — do każdego z nich stosuje się
warunki jego dostawcy (rozdz. 3.3, 3.4).

---

*Danaco Console — Platforma AI Workspace OS · v2.0 · Licencja w wersji 2.0*
*© 2026 Danaco Holding Group Sp. z o.o. Wszelkie prawa zastrzeżone — Dariusz Naharnowicz*
