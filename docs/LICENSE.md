# Licencja produktu Danaco Console

**Umowa licencyjna oprogramowania Danaco Console — Platforma AI Workspace OS**

| | |
|---|---|
| **Produkt** | Danaco Console |
| **Rodzaj** | Platforma AI Workspace OS |
| **Producent** | Danaco Holding Group Sp. z o.o. |
| **Twórca** | Dariusz Naharnowicz |
| **Wersja** | v2.0 |
| **Status** | Deweloperski |
| **Data** | 2026-08-20 |

**Informacje szczegółowe dokumentu:**

| | |
|---|---|
| **Tytuł** | Umowa licencyjna oprogramowania Danaco Console — Platforma AI Workspace OS |
| **Klasa dokumentu** | Akt prawny |
| **Odbiorcy** | Licencjobiorca · Operator · Właściciel |
| **Przeznaczenie** | Ustala warunki, na jakich Danaco Holding Group Sp. z o.o. udostępnia produkt Danaco Console — zakres uprawnień, ograniczenia, odpowiedzialność oraz ujawnienie rzeczywistego stanu wykonania. Na jego podstawie rozstrzyga się, co Licencjobiorcy wolno, a czego nie. |
| **Zakres** | postanowienia wstępne, definicje, przedmiot licencji, udzielenie licencji i zakres uprawnień, zasady korzystania i ograniczenia, własność intelektualna, komponenty i rozszerzenia, aktualizacje i wersjonowanie, dane i treści Operatora, ograniczenie rękojmi i odpowiedzialności, okres obowiązywania, zgodność licencyjna, poufność, postanowienia końcowe oraz załączniki A, B i C |
| **Poza zakresem** | techniczny opis produktu — [Opis produktu](README.md); obsługa produktu przez Operatora — [Instrukcja użytkowania](INSTRUKCJA-UZYTKOWANIA.md); instalacja i konfiguracja — [Instalacja i konfiguracja](INSTALACJA-I-KONFIGURACJA.md) |
| **Dokument nadrzędny** | [Spis opracowań](SPIS-OPRACOWAN.md) |
| **Dokumenty powiązane** | [Opis produktu](README.md) · [Instrukcja użytkowania](INSTRUKCJA-UZYTKOWANIA.md) · [Instalacja i konfiguracja](INSTALACJA-I-KONFIGURACJA.md) · [Standard redakcyjny i językowy](STANDARD-REDAKCYJNY-I-JEZYKOWY.md) |
| **Prototypy odniesienia** | nie dotyczy — akt prawny nie opisuje okna produktu ani jego układu |
| **Źródła normatywne** | `budowa/shared/contract.json` · `budowa/server/internal/store/` · `budowa/client/src/` · audyt repozytorium opisany w Załączniku C |
| **Zasada nadrzędna** | Licencja nie przedstawia jako działającego niczego, czego nie potwierdza stan faktyczny kodu — rozdz. 1.3 i rozdz. 10.3. |
| **Podstawa faktograficzna** | audyt repozytorium przeprowadzony w stanie sprzed prac naprawczych, zaktualizowany po zamknięciu prac naprawczych do stanu bieżącego repozytorium, z której zbudowano i uruchomiono doręczany Egzemplarz — patrz Załącznik C.1 i C.3 pkt 4 |
| **Data ustalenia stanu faktycznego** | 2026-08-12 — dzień pomiaru repozytorium; pole **Data** metryki produktowej niesie datę wydania dokumentu, nie datę pomiaru |
| **Zakres obowiązywania** | wszystkie egzemplarze Danaco Console v2.0 o statusie Deweloperskim |
| **Stan dokumentu** | **projekt wymagający przeglądu prawnego przed publikacją** — patrz rozdz. 1.5 |

---

## Spis treści

1. [Postanowienia wstępne](#1-postanowienia-wstępne)
   - [1.1 Charakter dokumentu](#11-charakter-dokumentu)
   - [1.2 Strony stosunku licencyjnego](#12-strony-stosunku-licencyjnego)
   - [1.3 Zasada zgodności ze stanem faktycznym produktu](#13-zasada-zgodności-ze-stanem-faktycznym-produktu)
   - [1.4 Hierarchia dokumentów projektu](#14-hierarchia-dokumentów-projektu)
   - [1.5 Zastrzeżenie o przeglądzie prawnym](#15-zastrzeżenie-o-przeglądzie-prawnym)
   - [1.6 Sposób oznaczania kwestii nierozstrzygniętych](#16-sposób-oznaczania-kwestii-nierozstrzygniętych)
2. [Definicje](#2-definicje)
   - [2.1 Podmioty](#21-podmioty)
   - [2.2 Oprogramowanie i jego warstwy](#22-oprogramowanie-i-jego-warstwy)
   - [2.3 Pojęcia funkcjonalne produktu](#23-pojęcia-funkcjonalne-produktu)
   - [2.4 Pojęcia dotyczące danych](#24-pojęcia-dotyczące-danych)
   - [2.5 Pojęcia dotyczące dystrybucji i wersji](#25-pojęcia-dotyczące-dystrybucji-i-wersji)
   - [2.6 Reguły wykładni](#26-reguły-wykładni)
3. [Przedmiot licencji](#3-przedmiot-licencji)
   - [3.1 Zakres przedmiotowy](#31-zakres-przedmiotowy)
   - [3.2 Postać dystrybucyjna i topologia](#32-postać-dystrybucyjna-i-topologia)
   - [3.3 Elementy wyłączone z przedmiotu licencji](#33-elementy-wyłączone-z-przedmiotu-licencji)
   - [3.4 Zależności zewnętrzne warunkujące działanie](#34-zależności-zewnętrzne-warunkujące-działanie)
   - [3.5 Dokumentacja objęta licencją](#35-dokumentacja-objęta-licencją)
4. [Udzielenie licencji i zakres uprawnień](#4-udzielenie-licencji-i-zakres-uprawnień)
   - [4.1 Udzielenie licencji](#41-udzielenie-licencji)
   - [4.2 Uprawnienia przyznane](#42-uprawnienia-przyznane)
   - [4.3 Zakres podmiotowy — liczba stanowisk i użytkowników](#43-zakres-podmiotowy--liczba-stanowisk-i-użytkowników)
   - [4.4 Zakres terytorialny i czasowy](#44-zakres-terytorialny-i-czasowy)
   - [4.5 Odpłatność](#45-odpłatność)
   - [4.6 Brak przeniesienia praw](#46-brak-przeniesienia-praw)
   - [4.7 Zakres wykonania modelu a zakres licencji](#47-zakres-wykonania-modelu-a-zakres-licencji)
5. [Zasady korzystania i ograniczenia](#5-zasady-korzystania-i-ograniczenia)
   - [5.1 Zasada ogólna](#51-zasada-ogólna)
   - [5.2 Ograniczenia dotyczące samego Oprogramowania](#52-ograniczenia-dotyczące-samego-oprogramowania)
   - [5.3 Ograniczenia wynikające ze statusu Deweloperskiego](#53-ograniczenia-wynikające-ze-statusu-deweloperskiego)
   - [5.4 Zasady korzystania z modelu i odpowiedzialność za polecenia](#54-zasady-korzystania-z-modelu-i-odpowiedzialność-za-polecenia)
   - [5.5 Obowiązki Licencjobiorcy w zakresie bezpieczeństwa](#55-obowiązki-licencjobiorcy-w-zakresie-bezpieczeństwa)
   - [5.6 Zgodność z warunkami Dostawców modeli](#56-zgodność-z-warunkami-dostawców-modeli)
   - [5.7 Uprawnienia wynikające z przepisów bezwzględnie obowiązujących](#57-uprawnienia-wynikające-z-przepisów-bezwzględnie-obowiązujących)
   - [5.8 Skutki naruszenia ograniczeń](#58-skutki-naruszenia-ograniczeń)
6. [Własność intelektualna](#6-własność-intelektualna)
   - [6.1 Prawa do Oprogramowania](#61-prawa-do-oprogramowania)
   - [6.2 Zastrzeżenie praw](#62-zastrzeżenie-praw)
   - [6.3 Zakaz działań naruszających](#63-zakaz-działań-naruszających)
   - [6.4 Znaki towarowe i oznaczenia](#64-znaki-towarowe-i-oznaczenia)
   - [6.5 Prawa do Treści Operatora](#65-prawa-do-treści-operatora)
   - [6.6 Informacje zwrotne](#66-informacje-zwrotne)
7. [Komponenty i rozszerzenia](#7-komponenty-i-rozszerzenia)
   - [7.1 Zasada nadrzędna](#71-zasada-nadrzędna)
   - [7.2 Komponenty warstwy Rdzenia (Go)](#72-komponenty-warstwy-rdzenia-go)
   - [7.3 Komponenty warstwy Klienta (TypeScript)](#73-komponenty-warstwy-klienta-typescript)
   - [7.4 Komponenty warstwy Powłoki (Rust / Tauri)](#74-komponenty-warstwy-powłoki-rust--tauri)
   - [7.5 Kroje pisma](#75-kroje-pisma)
   - [7.6 Zestaw ikon i zasoby marki](#76-zestaw-ikon-i-zasoby-marki)
   - [7.7 Program zewnętrzny kanału głównego](#77-program-zewnętrzny-kanału-głównego)
   - [7.8 Rozszerzenia](#78-rozszerzenia)
   - [7.9 Obowiązki dokumentacyjne przy dystrybucji](#79-obowiązki-dokumentacyjne-przy-dystrybucji)
8. [Aktualizacje i wersjonowanie](#8-aktualizacje-i-wersjonowanie)
   - [8.1 Zasada wersjonowania](#81-zasada-wersjonowania)
   - [8.2 Charakter Aktualizacji](#82-charakter-aktualizacji)
   - [8.3 Migracje schematu trwałości](#83-migracje-schematu-trwałości)
   - [8.4 Uprawnienie do Aktualizacji i sposób dostarczania](#84-uprawnienie-do-aktualizacji-i-sposób-dostarczania)
   - [8.5 Brak zobowiązania do rozwoju](#85-brak-zobowiązania-do-rozwoju)
   - [8.6 Zmiany treści Licencji](#86-zmiany-treści-licencji)
9. [Dane i treści Operatora](#9-dane-i-treści-operatora)
   - [9.1 Miejsce przechowywania danych](#91-miejsce-przechowywania-danych)
   - [9.2 Zakres danych zapisywanych przez Oprogramowanie](#92-zakres-danych-zapisywanych-przez-oprogramowanie)
   - [9.3 Postępowanie z Poświadczeniami](#93-postępowanie-z-poświadczeniami)
   - [9.4 Przekazywanie danych na zewnątrz](#94-przekazywanie-danych-na-zewnątrz)
   - [9.5 Kopie zapasowe i usuwanie danych](#95-kopie-zapasowe-i-usuwanie-danych)
   - [9.6 Dane osobowe](#96-dane-osobowe)
   - [9.7 Dostęp Licencjodawcy do danych przy czynnościach wsparcia](#97-dostęp-licencjodawcy-do-danych-przy-czynnościach-wsparcia)
10. [Ograniczenie rękojmi i odpowiedzialności](#10-ograniczenie-rękojmi-i-odpowiedzialności)
   - [10.1 Charakter świadczenia](#101-charakter-świadczenia)
   - [10.2 Wyłączenie rękojmi](#102-wyłączenie-rękojmi)
   - [10.3 Ujawnienie rzeczywistego stanu wykonania](#103-ujawnienie-rzeczywistego-stanu-wykonania)
   - [10.4 Ograniczenie odpowiedzialności](#104-ograniczenie-odpowiedzialności)
   - [10.5 Rozkład ryzyka](#105-rozkład-ryzyka)
   - [10.6 Siła wyższa](#106-siła-wyższa)
11. [Okres obowiązywania i zakończenie licencji](#11-okres-obowiązywania-i-zakończenie-licencji)
   - [11.1 Wejście w życie](#111-wejście-w-życie)
   - [11.2 Okres obowiązywania](#112-okres-obowiązywania)
   - [11.3 Wypowiedzenie i rozwiązanie](#113-wypowiedzenie-i-rozwiązanie)
   - [11.4 Skutki zakończenia](#114-skutki-zakończenia)
   - [11.5 Postanowienia trwające](#115-postanowienia-trwające)
12. [Zgodność licencyjna i kontrola korzystania](#12-zgodność-licencyjna-i-kontrola-korzystania)
   - [12.1 Obowiązek zgodności](#121-obowiązek-zgodności)
   - [12.2 Weryfikacja](#122-weryfikacja)
   - [12.3 Skutki stwierdzenia niezgodności](#123-skutki-stwierdzenia-niezgodności)
13. [Poufność](#13-poufność)
   - [13.1 Informacje poufne](#131-informacje-poufne)
   - [13.2 Obowiązki](#132-obowiązki)
   - [13.3 Okres obowiązywania](#133-okres-obowiązywania)
   - [13.4 Ujawnienie podatności](#134-ujawnienie-podatności)
14. [Postanowienia końcowe](#14-postanowienia-końcowe)
   - [14.1 Prawo właściwe](#141-prawo-właściwe)
   - [14.2 Rozstrzyganie sporów](#142-rozstrzyganie-sporów)
   - [14.3 Forma zmian](#143-forma-zmian)
   - [14.4 Klauzula salwatoryjna](#144-klauzula-salwatoryjna)
   - [14.5 Całość porozumienia](#145-całość-porozumienia)
   - [14.6 Język dokumentu](#146-język-dokumentu)
   - [14.7 Dane kontaktowe i doręczenia](#147-dane-kontaktowe-i-doręczenia)
   - [14.8 Nagłówki i załączniki](#148-nagłówki-i-załączniki)
15. [Załącznik A — zestawienie komponentów osób trzecich](#załącznik-a--zestawienie-komponentów-osób-trzecich)
   - [A.1. Warstwa Rdzenia — moduł Go `danacoconsole` (wymagane `go 1.26`)](#a1-warstwa-rdzenia--moduł-go-danacoconsole-wymagane-go-126)
   - [A.2. Warstwa Klienta — pakiet npm `danaco-console-client`](#a2-warstwa-klienta--pakiet-npm-danaco-console-client)
   - [A.3. Warstwa Powłoki — pakiet Cargo `danaco-console-powloka`](#a3-warstwa-powłoki--pakiet-cargo-danaco-console-powloka)
   - [A.4. Kroje pisma](#a4-kroje-pisma)
   - [A.5. Zestaw ikon i zasoby marki](#a5-zestaw-ikon-i-zasoby-marki)
   - [A.6. Zależności zewnętrzne nieobjęte dystrybucją](#a6-zależności-zewnętrzne-nieobjęte-dystrybucją)
16. [Załącznik B — wykaz kwestii pozostawionych do decyzji Operatora](#załącznik-b--wykaz-kwestii-pozostawionych-do-decyzji-operatora)
17. [Załącznik C — metryka weryfikacji faktograficznej](#załącznik-c--metryka-weryfikacji-faktograficznej)
   - [C.1. Podstawa faktograficzna](#c1-podstawa-faktograficzna)
   - [C.2. Elementy zweryfikowane bezpośrednio w kodzie](#c2-elementy-zweryfikowane-bezpośrednio-w-kodzie)
   - [C.3. Ograniczenia weryfikacji](#c3-ograniczenia-weryfikacji)
   - [C.4. Zakres wymagający ponownego przeglądu przy każdym Wydaniu](#c4-zakres-wymagający-ponownego-przeglądu-przy-każdym-wydaniu)

---

## 1. Postanowienia wstępne

### 1.1 Charakter dokumentu

Niniejszy dokument („**Licencja**”) określa warunki, na jakich Danaco Holding Group
Sp. z o.o. udostępnia oprogramowanie Danaco Console — Platforma AI Workspace OS
(„**Oprogramowanie**”) do korzystania. Licencja reguluje zakres uprawnień
przyznanych korzystającemu, granice dozwolonego użycia, zasady dotyczące
komponentów osób trzecich wbudowanych w Oprogramowanie, sposób postępowania
z danymi wytworzonymi w toku korzystania, a także zasady odpowiedzialności
Licencjodawcy odpowiadające rzeczywistemu, potwierdzonemu stanowi wykonania
produktu.

Oprogramowanie jest **produktem własnościowym**. Nie jest oprogramowaniem
otwartoźródłowym, nie jest oprogramowaniem darmowym w rozumieniu swobody
redystrybucji i nie podlega żadnej z powszechnie stosowanych licencji wolnego
oprogramowania. Fakt, że Oprogramowanie zawiera komponenty osób trzecich
rozpowszechniane na licencjach otwartych (rozdz. 7), nie zmienia własnościowego
charakteru samego Oprogramowania ani nie rozciąga warunków tych licencji na kod
autorski Licencjodawcy — z zastrzeżeniem obowiązków wynikających wprost
z licencji poszczególnych komponentów, opisanych w rozdz. 7.

Licencja odnosi się do wersji **v2.0 o statusie Deweloperskim**. Status ten nie
jest formułą marketingową ani zastrzeżeniem ostrożnościowym: wynika z przyjętej
decyzji architektonicznej („v2.0 do pierwszej publikacji; zakaz wersjonowania w
trakcie budowy; status Deweloperski”) i odpowiada faktycznemu stanowi
wykonania produktu, ustalonemu w drodze pełnego audytu repozytorium
przeprowadzonego pierwotnie w stanie sprzed prac naprawczych, a następnie zweryfikowanemu
ponownie i uzupełnionemu po zamknięciu prac naprawczych do stanu
bieżącego stanu repozytorium tej samej gałęzi — rewizji, z której zbudowano i uruchomiono
doręczany Egzemplarz (Załącznik C.1, C.3 pkt 4). Konsekwencje statusu Deweloperskiego dla
rękojmi i odpowiedzialności opisuje rozdz. 10; konsekwencje dla zakresu
funkcjonalnego, do którego Operator ma prawo mieć zaufanie, opisuje rozdz. 3.1
oraz rozdz. 10.3.

### 1.2 Strony stosunku licencyjnego

Stronami stosunku licencyjnego są:

1. **Licencjodawca** — Danaco Holding Group Sp. z o.o., podmiot, któremu
   przysługują autorskie prawa majątkowe do Oprogramowania, występujący
   w dokumentacji projektu również jako **Producent**;
2. **Licencjobiorca** — podmiot, któremu Licencjodawca udostępnił egzemplarz
   Oprogramowania na warunkach niniejszej Licencji.

Dokumentacja techniczna produktu posługuje się konsekwentnie pojęciem
**Operator** na oznaczenie osoby, która faktycznie obsługuje Oprogramowanie:
konfiguruje kanały modelu, konta, tożsamość modelu, punkty dostępu i nadania,
prowadzi sesje i okna komunikacji oraz odpowiada za wynik pracy wykonanej
z udziałem modelu. Pojęcia **Licencjobiorca** i **Operator** nie są tożsame:
Licencjobiorcą może być osoba prawna, podczas gdy Operatorem jest zawsze osoba
fizyczna zasiadająca przy Oprogramowaniu. Relacja między tymi rolami — a w
szczególności to, czy Licencjobiorcą jest osoba fizyczna korzystająca
z Oprogramowania, czy podmiot, na którego rzecz Operatorzy pracują —
pozostaje kwestią modelu licencjonowania.

> **[DO DECYZJI OPERATORA] — definicja Licencjobiorcy.**
> Dokumentacja projektu (README, zapis decyzji architektonicznych, brief
> projektowy, specyfikacja Fali 2) nie zawiera rozstrzygnięcia, kto jest
> Licencjobiorcą w rozumieniu prawnym. Warianty do wyboru:
> **(a)** licencja imienna na osobę fizyczną (Operatora), niezbywalna, jedno
> stanowisko na osobę;
> **(b)** licencja na podmiot gospodarczy z prawem wskazania określonej liczby
> Operatorów;
> **(c)** licencja na urządzenie (instalację), niezależna od liczby osób, które
> z niej korzystają;
> **(d)** licencja wewnętrzna grupy kapitałowej — korzystanie wyłącznie przez
> Danaco Holding Group Sp. z o.o. i podmioty powiązane, bez dystrybucji na
> zewnątrz.
> Wybór wariantu przesądza treść rozdz. 4.3, 11.4 i 12 i musi zostać dokonany
> przed pierwszym udostępnieniem Oprogramowania podmiotowi zewnętrznemu.

### 1.3 Zasada zgodności ze stanem faktycznym produktu

Licencja opisuje produkt rzeczywisty, nie zamierzony. Wszędzie tam, gdzie
konieczne jest odniesienie do funkcji Oprogramowania — w szczególności
w rozdz. 3 (przedmiot licencji), rozdz. 9 (dane Operatora) i rozdz. 10
(odpowiedzialność) — przywołany jest stan wykonania potwierdzony audytem
repozytorium, a nie stan zaprojektowany w dokumentacji koncepcyjnej.

Reguła ta ma znaczenie prawne, nie redakcyjne: zakres, w jakim Licencjodawca
mógłby ponosić odpowiedzialność za niezgodność Oprogramowania z opisem, wyznacza
opis zawarty w niniejszej Licencji i w dokumentacji produktu przekazanej
Licencjobiorcy — to jest w zamkniętym zestawie dokumentów wymienionym
w rozdz. 1.4 pkt 3 i rozdz. 3.5 — a nie materiały koncepcyjne, prezentacyjne
ani rejestry decyzji
architektonicznych opisujące zamierzenia. Funkcje przewidziane koncepcją, lecz
w wersji v2.0 niezaimplementowane albo niepodłączone do interfejsu, oznaczono
w treści Licencji jednoznacznie — nie stanowią one przedmiotu świadczenia.

### 1.4 Hierarchia dokumentów projektu

W razie rozbieżności między dokumentami dotyczącymi Oprogramowania stosuje się
następującą kolejność pierwszeństwa:

1. niniejsza Licencja — w zakresie uprawnień, ograniczeń, odpowiedzialności
   i zasad dotyczących komponentów osób trzecich;
2. odrębna umowa zawarta na piśmie między Licencjodawcą a Licencjobiorcą, jeżeli
   została zawarta i jeżeli wprost modyfikuje postanowienia Licencji;
3. **dokumentacja produktu** przekazana Licencjobiorcy wraz z Egzemplarzem,
   obejmująca — poza niniejszym plikiem [LICENSE.md](LICENSE.md) — trzy dokumenty
   oznaczone nazwami plików:
   - **[README.md](README.md)** — karta techniczno-produktowa Danaco Console: koncepcja,
     architektura trzech warstw, stos technologiczny, model katalogów oraz
     zakres funkcjonalny wraz z oznaczeniem stanu wykonania poszczególnych
     funkcji;
   - **[INSTRUKCJA-UZYTKOWANIA.md](INSTRUKCJA-UZYTKOWANIA.md)** — opis pracy z gotową Instalacją: strona
     główna, Środowiska, Okno komunikacji, okna równoległe, okno konfiguracji,
     punkty dostępu, konta i tożsamość modelu;
   - **[INSTALACJA-I-KONFIGURACJA.md](INSTALACJA-I-KONFIGURACJA.md)** — procedura instalacji, budowania
     warstw i doprowadzenia świeżej Instalacji do stanu zdatnego do pracy,
     wykaz zmiennych środowiska i argumentów wywołania oraz postępowanie
     w razie niepowodzeń;
4. dokumentacja projektowa i wewnętrzna Licencjodawcy — **nieprzekazywana
   Licencjobiorcy i nieskładająca się na treść zobowiązania** — w szczególności
   zapis decyzji architektonicznych, wykaz luk, brief projektowy,
   specyfikacje fal budowy, pakiety warstwy wizualnej oraz opracowania
   klasy Specyfikacja docelowa katalogów `architektura/`, `moduly/`,
   `interfejs-uzytkownika/`, `srodowiska/`, `specyfikacje/` i
   `funkcje-globalne/` — dokumentacja techniczna towarzysząca kodowi
   repozytorium budowy, odrębna od czterech dokumentów produktu wymienionych
   w pkt 3.

Wykaz z pkt 3 jest zamknięty: dokumentem produktu jest wyłącznie plik wprost
w nim wymieniony. Dokument nienależący do tego wykazu — choćby dotyczył
Oprogramowania i pochodził od Licencjodawcy — nie tworzy zobowiązania co do
zakresu funkcjonalnego (rozdz. 14.5).

### 1.5 Zastrzeżenie o przeglądzie prawnym

**Niniejszy dokument jest projektem przygotowanym na podstawie dokumentacji
technicznej i stanu faktycznego repozytorium. Przed publikacją, przed
dołączeniem do instalatora oraz przed udostępnieniem Oprogramowania
jakiemukolwiek podmiotowi zewnętrznemu dokument wymaga przeglądu przez
radcę prawnego lub adwokata.**

Przegląd prawny powinien objąć w szczególności:

- skuteczność wyłączenia rękojmi i ograniczenia odpowiedzialności wobec
  konsumentów oraz wobec przedsiębiorców, w świetle prawa właściwego wybranego
  zgodnie z rozdz. 14.1;
- zgodność klauzul dotyczących danych osobowych z rozporządzeniem (UE) 2016/679
  (RODO) oraz ustalenie ról administratora i podmiotu przetwarzającego
  (rozdz. 9.6);
- kompletność i poprawność zestawienia komponentów osób trzecich wraz
  z wymaganymi notami licencyjnymi (rozdz. 7 i Załącznik A), w tym pełne
  zbadanie zamknięcia zależności warstwy Rust, którego wykaz w Załączniku A
  jest zweryfikowany częściowo;
- dopuszczalność ograniczeń korzystania w świetle bezwzględnie obowiązujących
  przepisów o dozwolonym użytku programów komputerowych;
- zgodność sposobu korzystania z zewnętrznych usług modeli językowych
  z warunkami umownymi ich dostawców (rozdz. 3.4 i 5.6);
- ocenę dopuszczalności i skuteczności klauzuli zakazu budowy produktu
  konkurencyjnego oraz zakazu badań bezpieczeństwa (rozdz. 5.2 pkt 8 i 9),
  a także oznaczenia terminu naprawczego (rozdz. 11.3 pkt 2) i okresu
  poufności następczej (rozdz. 13.3);
- sprawdzenie zgodności niniejszego dokumentu z pozostałymi dokumentami
  produktu wymienionymi w rozdz. 1.4 pkt 3 — w szczególności co do opisu stanu
  wykonania funkcji i co do danych liczbowych warstwy wizualnej;
- rozstrzygnięcie wszystkich pozycji oznaczonych jako
  **[DO DECYZJI OPERATORA]** i zebranych w Załączniku B.

Do czasu zakończenia przeglądu prawnego dokument należy traktować jako materiał
roboczy Licencjodawcy.

### 1.6 Sposób oznaczania kwestii nierozstrzygniętych

Wszędzie tam, gdzie dokumentacja projektu nie rozstrzyga kwestii prawnej lub
biznesowej niezbędnej dla kompletności Licencji, w treści umieszczono oznaczenie
**[DO DECYZJI OPERATORA]** wraz z krótkim opisem możliwych wariantów i wskazaniem
skutków wyboru. Oznaczenia te są **celowe**: dokument nie przesądza za
Licencjodawcę kwestii, które nie zostały przez niego rozstrzygnięte, i nie
zastępuje decyzji biznesowej domysłem redakcyjnym.

Komplet oznaczeń zebrano w Załączniku B. Dokument nie może zostać opublikowany
ani dołączony do dystrybucji przed usunięciem wszystkich takich oznaczeń przez
wpisanie rozstrzygnięć.

---

## 2. Definicje

Pojęcia pisane w Licencji wielką literą mają znaczenie nadane im w niniejszym
rozdziale. Terminologia funkcjonalna odpowiada terminologii obowiązującej
w dokumentacji technicznej Danaco Console; rozbieżność między znaczeniem
potocznym a znaczeniem zdefiniowanym rozstrzyga się na rzecz znaczenia
zdefiniowanego.

### 2.1 Podmioty

**Licencjodawca** — Danaco Holding Group Sp. z o.o., podmiot uprawniony
z autorskich praw majątkowych do Oprogramowania; w dokumentacji technicznej
występujący jako Producent.

**Twórca** — Dariusz Naharnowicz, autor koncepcji produktu i twórca
Oprogramowania w rozumieniu prawa autorskiego, wskazany jako twórca w metryce
każdego dokumentu projektu.

**Licencjobiorca** — podmiot, któremu udzielono licencji na korzystanie
z Oprogramowania na warunkach niniejszego dokumentu.

**Operator** — osoba fizyczna faktycznie obsługująca Oprogramowanie:
konfigurująca je, prowadząca sesje pracy z modelem oraz odpowiadająca za
polecenia wydawane modelowi i za wykorzystanie wyników jego pracy. Operator
działa w imieniu i na rzecz Licencjobiorcy.

**Dostawca modelu** — podmiot trzeci świadczący usługę udostępniania modelu
sztucznej inteligencji albo dostarczający program kliencki umożliwiający
korzystanie z takiego modelu, w szczególności dostawca programu
`claude` (Claude Code CLI) wykorzystywanego przez kanał główny Oprogramowania.

**Osoba trzecia** — podmiot inny niż Licencjodawca i Licencjobiorca,
w szczególności autor lub uprawniony z praw do Komponentu osoby trzeciej.

### 2.2 Oprogramowanie i jego warstwy

**Oprogramowanie** (Danaco Console) — całość programu komputerowego objętego
niniejszą Licencją, obejmująca Rdzeń, Klienta, Powłokę, Kontrakt, schemat
trwałości wraz z migracjami, warstwę wizualną oraz zasoby wbudowane, w postaci
wykonywalnej dostarczonej Licencjobiorcy.

**Rdzeń** — serwerowa część Oprogramowania napisana w języku Go (moduł
`danacoconsole`). Rdzeń prowadzi trwałość, obsługuje komendy Kontraktu,
zarządza sesjami i oknami komunikacji, rozstrzyga konfigurację, uruchamia proces
Kanału modelu i przekazuje strumień jego odpowiedzi.

**Klient** (interfejs) — część Oprogramowania napisana w języku TypeScript i
renderowana przez komponent webview systemu operacyjnego. Klient odpowiada za
prezentację, przyjmowanie poleceń Operatora i komunikację z Rdzeniem po
Kontrakcie.

**Powłoka** — natywna warstwa okna zbudowana w technologii Tauri (język Rust),
uruchamiająca Rdzeń w tle, otwierająca okno interfejsu i osadzająca ikonę w
zasobniku systemowym.

**Kontrakt** — zbiór definicji nazw komend, zdarzeń, strumieni, narzędzi
i kształtu koperty komunikatu, stanowiący jedyne źródło prawdy nazewnictwa
w komunikacji Klienta z Rdzeniem. Kontrakt v2.0 obejmuje **53 komendy, 31
zdarzeń i 39 deklaracji narzędzi modelu**, z których generowane są wiązania dla
obu warstw.

Koperta komunikatu jest jedna dla żądania, odpowiedzi i fragmentu strumienia;
przesyłana jest w formacie JSON po połączeniu WebSocket. Jej pola
**wspólne** to: `type` (wymagane — typ komunikatu w notacji
`obszar.zasob.akcja`), `id` (wymagane — identyfikator zadania powtarzany
w odpowiedzi i we fragmentach strumienia), `timestamp` (wymagane — czas nadania
w milisekundach epoki), a ponadto `sessionId` (opcjonalne — sesja, której
komunikat dotyczy; puste dla powitania połączenia) oraz `payload` (opcjonalne —
treść właściwa, której kształt wyznacza typ komunikatu). Wyłącznie w odpowiedzi
wypełniane są pola `status` i — przy statusie błędu — `error`; wyłącznie we
fragmentach strumienia pola `seq` (numer fragmentu, liczony od jedności)
i `done` (prawda w ostatnim fragmencie).

Skrót `{type, sessionId, payload}`, którym posługują się
[Instrukcja użytkowania](INSTRUKCJA-UZYTKOWANIA.md) rozdz. 1.4 oraz
[Instalacja i konfiguracja](INSTALACJA-I-KONFIGURACJA.md), nazywa wyłącznie trzy
z powyższych pól — te przywoływane w tamtych dokumentach najczęściej — i nie
zastępuje wykazu pełnego podanego tutaj.

Brzmienie powyższych nazw jest wiążące zgodnie z regułą wykładni z rozdz. 2.6
pkt 5; źródłem prawdy pozostaje plik Kontraktu, a przytoczony wykaz odpowiada
jego stanowi w bieżącym stanie repozytorium (liczba komend, zdarzeń i deklaracji narzędzi
nie uległa zmianie względem stanu sprzed prac naprawczych, na której przeprowadzono audyt
pierwotny — zgodność potwierdzona odczytem `budowa/shared/contract.json`, źródła
normatywnego Kontraktu).

**Katalog danych** — katalog systemu plików, w którym Rdzeń przechowuje
trwałość Oprogramowania. Domyślnym katalogiem danych jest
`%LOCALAPPDATA%\DanacoConsole`; katalog może zostać wskazany zmienną środowiska
`DANACO_KATALOG_DANYCH` albo argumentem wywołania `--dane`.

**Baza** — jeden plik bazy danych SQLite o nazwie `danaco-console.db`
umieszczony w Katalogu danych, stanowiący całość trwałości Rdzenia.
Schemat Bazy powstaje i jest aktualizowany przez migracje wkompilowane
w Rdzeń, stosowane samoczynnie przy starcie.

### 2.3 Pojęcia funkcjonalne produktu

**Środowisko** — profil widoczności modułów w nawigacji Oprogramowania.
Zdefiniowane są cztery Środowiska: **TalkIn** (wiedza), **WorkSpace** (praca),
**CodeStudio** (technologia) i **MultitaskingAI** (orkiestracja inteligencji).
Środowisko nie jest pojemnikiem na maszyny ani katalogi.

**Moduł** — jednostka funkcjonalna widoczna w nawigacji Środowiska.

**Okno operacyjne** — wyspecjalizowany widok roboczy Modułu.

**Sesja** — nadrzędna jednostka pracy Operatora, wspólna dla plików, pamięci,
projektu i agentów; prezentowana jako **karta sesji**.

**Okno komunikacji** — byt pośredni między Sesją a wiadomością, niosący moduł,
kanał modelu, listę katalogów roboczych, środowisko wykonania, tryb uprawnień i
rolę. Okno komunikacji jest najwęższym, ósmym poziomem zasięgu konfiguracji.

**Kanał modelu** — wiersz rejestru sterowanego danymi opisujący sposób
wywołania modelu. Kolumna `rodzaj_kanalu` schematu trwałości dopuszcza cztery
wartości: `cli` (kanał główny, uruchamiający zewnętrzny program wiersza
poleceń), `api` (kanał generyczny, wywołujący usługę sieciową opisaną wierszem
rejestru), `sdk` oraz `lokalny`. Adapter kanału **nie jest** tożsamy
z wartością kolumny `rodzaj_kanalu` — dobierany jest kluczem `adapter`
z parametrów wiersza, wartością danych, nie nazwą typu. Kanał testowy **echo**
(bez sieci i bez procesu, odsyłający treść zapytania) zakłada
się jako wiersz rodzaju `lokalny` z parametrem `{"adapter":"echo"}` — „echo”
nie jest samodzielną wartością kolumny `rodzaj_kanalu`.

**Konto** — wpis opisujący tożsamość używaną przez Kanał modelu. Dla kanału
rodzaju `cli` Konto wskazuje katalog konfiguracji przekazywany procesowi modelu
zmienną `CLAUDE_CONFIG_DIR`. Konto nie przechowuje treści poświadczenia —
przechowuje wyłącznie odwołanie do niego.

**Pula kont** — uporządkowany zbiór Kont jednego rodzaju, w obrębie którego
Rdzeń dokonuje rotacji przy rozpoznaniu wyczerpania limitu danego Konta.

**Nakładka tożsamości** — trójwarstwowy mechanizm kształtowania zachowania
modelu w kolejności krytyczności: **konstytucja → profil → ekspertyza**
przekazywany procesowi modelu w trybie `ZASTAP` albo `DOLACZ`.

**Punkt dostępu** — wpis opisujący maszynę lub katalog, do którego model ma
mieć wgląd. **Nadanie** — przypisanie Punktu dostępu do Okna komunikacji wraz z
trybem (`read` albo `write`) i podzbiorem korzeni.

**Prowenancja** — zbiór informacji o tym, co faktycznie zostało przekazane
modelowi w danym wywołaniu (argumenty wywołania, treść nakładki tożsamości,
ustawienia, sumy kontrolne). Prowenancja pełni funkcję przejrzystości, nie
kontroli dostępu. **Prowenancja nie jest utrwalana**: powstaje przed
uruchomieniem procesu Tury, żyje wyłącznie w pamięci tego procesu i jest
przekazywana do interfejsu na czas jej trwania. Schemat trwałości nie zawiera
tabeli ani kolumny prowenancji, w szczególności nie zawiera jej tabela
wiadomości. Po zamknięciu Tury odtworzenie z Bazy, co faktycznie zostało
przekazane modelowi, nie jest możliwe (rozdz. 10.3 wykaz ograniczeń, pkt 20).

**Tura** — pojedyncze wywołanie modelu w Oknie komunikacji, od przyjęcia
wiadomości Operatora do zamknięcia strumienia odpowiedzi.

### 2.4 Pojęcia dotyczące danych

**Dane Operatora** — wszelkie dane wprowadzone do Oprogramowania przez Operatora
albo powstałe w wyniku korzystania z Oprogramowania po stronie Licencjobiorcy,
w szczególności: treść wiadomości i odpowiedzi modelu zapisana w Bazie, nazwy
i struktura sesji, okien i kart, wartości ustawień, definicje kanałów modelu
i kont, definicje punktów dostępu i nadań, dokumenty tożsamości modelu, a także
pliki wytworzone przez model w katalogach roboczych.

**Treści Operatora** — wytwór intelektualny powstały przy użyciu Oprogramowania,
w szczególności teksty, dokumenty, kod źródłowy i inne materiały wygenerowane
albo przetworzone w toku pracy z modelem.

**Poświadczenie** — wartość umożliwiająca uwierzytelnienie wobec Dostawcy modelu
albo wobec maszyny wskazanej Punktem dostępu (klucz, token, hasło, klucz
prywatny). Poświadczenie nie jest przechowywane w Bazie — Baza przechowuje
wyłącznie odwołanie do niego.

**Katalog roboczy** — katalog systemu plików wskazany Oknu komunikacji, w którym
model zostawia i odczytuje pliki pracy.

### 2.5 Pojęcia dotyczące dystrybucji i wersji

**Egzemplarz** — pojedyncza kopia Oprogramowania udostępniona Licencjobiorcy
w postaci wykonywalnej.

**Instalacja** — czynność doprowadzenia Egzemplarza do stanu zdatnego do
uruchomienia na urządzeniu, wraz z utworzeniem Katalogu danych i Bazy.

**Wydanie** — oznaczony zestaw plików wykonywalnych i zasobów Oprogramowania
udostępniony przez Licencjodawcę jako całość.

**Aktualizacja** — Wydanie następujące po Wydaniu posiadanym przez
Licencjobiorcę, zastępujące je w całości albo w części.

**Status Deweloperski** — status wersji v2.0 wynikający z przyjętej zasady
wersjonowania, oznaczający, że produkt znajduje się w fazie budowy, a jego
zakres funkcjonalny nie jest domknięty. Konsekwencje statusu opisuje rozdz. 10.

**Rdzeń lokalny** i **rdzeń serwerowy** — dwa umiejscowienia Rdzenia
przewidziane topologią docelową: w fazie budowy Rdzeń działa lokalnie na
urządzeniu
Operatora, docelowo zaś przewidziane jest jego przeniesienie na serwer
`danaco-system` wraz z dystrybucją instalatorów zawierających Klienta i agenta
lokalnego na urządzenia użytkowników.

### 2.6 Reguły wykładni

1. Tytuły rozdziałów i podrozdziałów służą wyłącznie orientacji i nie wpływają
   na wykładnię postanowień.
2. Wyliczenia poprzedzone zwrotem „w szczególności” mają charakter przykładowy
   i nie wyczerpują zakresu pojęcia.
3. Odwołanie do rozdziału obejmuje wszystkie jego podrozdziały.
4. Liczba pojedyncza obejmuje liczbę mnogą i odwrotnie, o ile z kontekstu nie
   wynika inaczej.
5. Terminy techniczne zapisane czcionką o stałej szerokości (`w ten sposób`)
   oznaczają dosłowne nazwy elementów Oprogramowania: zmiennych środowiska,
   argumentów wywołania, plików, tabel, komend Kontraktu i kluczy konfiguracji.
   Nazwy te są wiążące co do brzmienia.

---

## 3. Przedmiot licencji

### 3.1 Zakres przedmiotowy

Przedmiotem Licencji jest korzystanie z Oprogramowania Danaco Console v2.0
w postaci wykonywalnej, obejmującego:

1. **Rdzeń** — serwer w języku Go realizujący obsługę Kontraktu, trwałość
   w pliku SQLite, rozstrzyganie konfiguracji, prowadzenie sesji i okien
   komunikacji, uruchamianie procesu kanału modelu oraz przekazywanie strumienia
   odpowiedzi;
2. **Klienta** — interfejs w języku TypeScript wraz z warstwą wizualną
   (żetony motywu, dwa motywy: jasny i ciemny, zestaw ikon, kroje pisma);
3. **Powłokę** — natywne okno aplikacji zbudowane w technologii Tauri wraz
   z osadzeniem ikony w zasobniku systemowym;
4. **Kontrakt** wraz z wygenerowanymi z niego wiązaniami obu warstw;
5. **schemat trwałości** wraz z kompletem migracji wkompilowanych w Rdzeń;
6. **zasoby wbudowane** — ikony, logotypy i sygnety marki, pliki krojów pisma;
7. **dokumentację produktu** przekazaną Licencjobiorcy wraz z Egzemplarzem —
   pliki [README.md](README.md), [INSTRUKCJA-UZYTKOWANIA.md](INSTRUKCJA-UZYTKOWANIA.md)
   i [INSTALACJA-I-KONFIGURACJA.md](INSTALACJA-I-KONFIGURACJA.md) wraz z niniejszym plikiem [LICENSE.md](LICENSE.md)
   (rozdz. 1.4 pkt 3 i rozdz. 3.5).

Przedmiot Licencji obejmuje wyłącznie **postać wykonywalną** Oprogramowania
wraz z dokumentacją produktu wymienioną w pkt 7. Kod źródłowy, repozytorium
projektu, dokumentacja projektowa i wewnętrzna Licencjodawcy (w tym zapis
decyzji architektonicznych, wykaz znanych luk oraz dokumentacja techniczna
repozytorium budowy), materiały koncepcyjne, pakiety projektowe warstwy
wizualnej oraz narzędzia budowania **nie stanowią przedmiotu Licencji** i nie są
udostępniane Licencjobiorcy, chyba że co innego wynika z odrębnej umowy
zawartej na piśmie.

Rozgraniczenie to biegnie po wykazie z rozdz. 1.4 pkt 3, nie po miejscu
powstania dokumentu: plik [README.md](README.md) dostarczany w katalogu produktu jest
**dokumentacją produktu objętą Licencją**, mimo że pełni jednocześnie funkcję
karty technicznej projektu. Nie jest natomiast objęta Licencją dokumentacja
projektowa pozostająca w repozytorium budowy i nieprzekazywana
Licencjobiorcy.

### 3.2 Postać dystrybucyjna i topologia

Oprogramowanie jest przeznaczone do pracy jako aplikacja desktopowa systemu
Windows z Rdzeniem działającym lokalnie. Powłoka uruchamia Rdzeń jako proces
w tle i otwiera okno interfejsu; komunikacja Klienta z Rdzeniem odbywa się po
połączeniu WebSocket, domyślnie na porcie **17870**. Konfiguracja Rdzenia jest
przyjmowana warstwowo — wartość domyślna, następnie zmienna środowiska,
następnie argument wywołania — a zmiennymi rozpoznawanymi przez Rdzeń są
wyłącznie: `DANACO_ROLA`, `DANACO_PORT`, `DANACO_KATALOG_DANYCH`,
`DANACO_KATALOG_KLIENTA` oraz `DANACO_KATALOG_PROFILI`.

Konfiguracja pakietu instalacyjnego Powłoki wskazuje nazwę produktu
„Danaco Console”, identyfikator `pl.danaco.console`, wydawcę
„Danaco Holding Group Sp. z o.o.”, notę „© 2026 Danaco Holding Group Sp. z o.o.”
oraz format instalatora NSIS dla systemu Windows.

Topologia docelowa przewiduje przeniesienie Rdzenia na serwer `danaco-system`
oraz dystrybucję instalatorów zawierających Klienta i agenta
lokalnego na urządzenia użytkowników. **W wersji v2.0 topologia docelowa nie
jest wdrożona**: przewidziana jest koncepcją, a Oprogramowanie w bieżącej wersji
pracuje z Rdzeniem lokalnym. Zmiana umiejscowienia Rdzenia nie zmienia zakresu
niniejszej Licencji; zasady licencjonowania korzystania z Rdzenia serwerowego
przez wielu Operatorów jednocześnie wymagają odrębnego rozstrzygnięcia.

> **[DO DECYZJI OPERATORA] — model dystrybucji końcowej.**
> Dokumentacja projektu ustala kierunek techniczny, lecz nie rozstrzyga modelu
> dystrybucji w rozumieniu handlowym. Warianty:
> **(a)** dystrybucja wyłącznie wewnętrzna w ramach grupy Danaco, bez
> udostępniania podmiotom zewnętrznym;
> **(b)** dystrybucja instalatora Windows bezpośrednio wskazanym
> Licencjobiorcom, z licencją per instalacja;
> **(c)** model usługowy — dostęp do Rdzenia serwerowego `danaco-system`
> z instalacją Klienta i agenta lokalnego na urządzeniu użytkownika, licencja
> per Operator;
> **(d)** model mieszany — instalacja lokalna dla pracy indywidualnej,
> Rdzeń serwerowy dla pracy zespołowej.
> Wybór przesądza treść rozdz. 4.3, 4.4, 8.4 i 12.

### 3.3 Elementy wyłączone z przedmiotu licencji

Przedmiotem Licencji **nie są**:

1. **modele sztucznej inteligencji** jakiegokolwiek dostawcy — Oprogramowanie
   nie zawiera modelu, nie dostarcza wag modelu ani nie udziela dostępu do
   usługi modelu;
2. **program `claude` (Claude Code CLI)** ani żaden inny zewnętrzny program
   uruchamiany przez Kanał modelu rodzaju `cli`; program taki Operator zapewnia
   we własnym zakresie (rozdz. 3.4);
3. **usługi sieciowe** wywoływane przez Kanał modelu rodzaju `api` — adres,
   poświadczenie i warunki korzystania z takiej usługi pochodzą od jej dostawcy;
4. **serwery MCP** i maszyny wskazane Punktami dostępu, wraz z oprogramowaniem
   na nich zainstalowanym;
5. **Komponenty osób trzecich** w zakresie, w jakim ich własne licencje
   przyznają uprawnienia szersze niż niniejsza Licencja albo nakładają
   obowiązki odrębne — do komponentów tych stosuje się ich własne licencje
   (rozdz. 7);
6. **kod źródłowy** Oprogramowania, repozytorium projektu oraz dokumentacja
   projektowa i wewnętrzna Licencjodawcy wymieniona w rozdz. 1.4 pkt 4;
   wyłączenie to **nie obejmuje** dokumentacji produktu z rozdz. 1.4 pkt 3,
   która jest Licencją objęta (rozdz. 3.1 pkt 7 i rozdz. 3.5);
7. **znaki towarowe, logotypy i oznaczenia** Licencjodawcy, poza zakresem
   niezbędnym do zwykłego korzystania z Oprogramowania (rozdz. 6.4).

### 3.4 Zależności zewnętrzne warunkujące działanie

Operator przyjmuje do wiadomości, że pełne wykorzystanie Oprogramowania wymaga
zależności pozostających poza zakresem Licencji:

1. **Program modelu dla kanału głównego.** Kanał modelu rodzaju `cli` uruchamia
   zewnętrzny program `claude` (Claude Code CLI), którego ścieżka pochodzi
   z wiersza rejestru kanału albo — w braku wskazania — z domyślnej nazwy
   `claude` odnajdywanej w zmiennej `PATH` systemu operacyjnego. Program ten
   nie jest dostarczany z Oprogramowaniem, a korzystanie z niego podlega
   warunkom jego dostawcy, w tym warunkom korzystania z konta u tego dostawcy.
   Uwaga: klucz konfiguracji `harness.program_claude` istnieje w katalogu
   ustawień Oprogramowania, lecz **w wersji v2.0 nie jest odczytywany** —
   wskazanie ścieżki programu innej niż domyślna następuje wyłącznie przez
   wiersz rejestru kanału.
2. **Konto Dostawcy modelu.** Kanał główny wymaga co najmniej jednego Konta
   rodzaju `cli` wskazującego katalog konfiguracji przekazywany procesowi
   modelu. Uzyskanie i utrzymanie takiego konta, w tym pokrycie kosztów
   i przestrzeganie limitów, obciąża Licencjobiorcę.
3. **Doprowadzenie świeżej instalacji do pracy.** W wersji v2.0 świeża
   Instalacja **nie jest zdolna do wykonania wywołania modelu bez czynności
   przygotowawczych**: rejestr kanałów modelu nie posiada zaczynu w migracjach,
   a pula kont pozostaje pusta. Doprowadzenie Instalacji do pracy wymaga
   utworzenia wiersza kanału głównego (rodzaj `cli`) oraz co najmniej jednego
   Konta rodzaju `cli`. Zakres dostępności tych dwóch czynności jest różny
   i wymaga rozróżnienia:
   - **kanał modelu** — **interfejs zakładania kanału nie istnieje w wersji
     v2.0**. Widoki aplikacji wywołują wyłącznie `channel.list`; utworzenie
     wiersza rejestru możliwe jest jedynie komendą Kontraktu `channel.add`
     wydaną poza interfejsem aplikacji. Komendą tą posługuje się samodzielna
     strona podglądu warstwy rozmowy (`podglad-rozmowy`), która przy braku
     kanału głównego dokłada jego wiersz; strona ta jest narzędziem podglądu
     warstwy, nie oknem aplikacji, i nie stanowi interfejsu zakładania kanałów
     w rozumieniu opisu przedmiotu świadczenia;
   - **konto** — pełny cykl życia Konta (`account.list`, `account.add`,
     `account.update`, `account.remove`) jest dostępny **z poziomu
     interfejsu**, w oknie „Modele, konta i tożsamość” (rozdz. 10.3 wykaz
     funkcji działających, pkt 7). Świeża Instalacja wymaga zatem założenia
     Konta czynnością Operatora w interfejsie, nie zaś komendą wydaną poza nim.

   Okoliczność ta jest istotnym elementem opisu przedmiotu świadczenia
   i została ujawniona wprost, aby wyłączyć wątpliwość co do zgodności
   Oprogramowania z opisem.
4. **Środowisko systemowe.** Oprogramowanie wymaga systemu operacyjnego
   Windows z dostępnym komponentem webview wykorzystywanym przez Powłokę oraz
   uprawnieniami pozwalającymi na zapis w Katalogu danych i nasłuch na porcie
   lokalnym.

### 3.5 Dokumentacja objęta licencją

Licencją objęta jest dokumentacja produktu przekazana Licencjobiorcy wraz
z Egzemplarzem, w zakresie niezbędnym do korzystania z Oprogramowania. Jest to
zamknięty zestaw czterech plików umieszczanych w katalogu produktu obok plików
wykonywalnych:

| Plik | Rola |
|---|---|
| [LICENSE.md](LICENSE.md) | niniejszy dokument — warunki licencyjne |
| [README.md](README.md) | karta techniczno-produktowa: koncepcja, architektura, stos, zakres funkcjonalny wraz ze stanem wykonania |
| [INSTRUKCJA-UZYTKOWANIA.md](INSTRUKCJA-UZYTKOWANIA.md) | praca z gotową Instalacją — widoki, Okno komunikacji, konfiguracja, konta i tożsamość modelu |
| [INSTALACJA-I-KONFIGURACJA.md](INSTALACJA-I-KONFIGURACJA.md) | instalacja, budowanie warstw, doprowadzenie świeżej Instalacji do pracy, zmienne środowiska i argumenty wywołania |

Wszędzie, gdzie Licencja posługuje się zwrotem „dokumentacja przekazana
Licencjobiorcy” albo „dokumentacja produktu”, rozumie się przez to wyłącznie
powyższy zestaw. Zwroty opisowe („dokumentacja wdrożeniowa”, „opis
funkcjonalny”, „instrukcja Operatora”) nie są nazwami odrębnych dokumentów;
odpowiadające im treści mieszczą się odpowiednio w plikach
[INSTALACJA-I-KONFIGURACJA.md](INSTALACJA-I-KONFIGURACJA.md), [README.md](README.md) i [INSTRUKCJA-UZYTKOWANIA.md](INSTRUKCJA-UZYTKOWANIA.md).

Licencjobiorca może sporządzać kopie dokumentacji na własne potrzeby
i udostępniać ją swoim Operatorom. Publikowanie dokumentacji, jej fragmentów
ani opracowań na jej podstawie poza organizacją Licencjobiorcy wymaga uprzedniej
zgody Licencjodawcy udzielonej na piśmie.

---

## 4. Udzielenie licencji i zakres uprawnień

### 4.1 Udzielenie licencji

Licencjodawca udziela Licencjobiorcy licencji **niewyłącznej**,
**nieprzenoszalnej**, **bez prawa udzielania sublicencji**, na korzystanie
z Oprogramowania w postaci wykonywalnej, wyłącznie na warunkach i w granicach
określonych niniejszym dokumentem.

Licencja obejmuje wyłącznie pola eksploatacji wskazane wprost w rozdz. 4.2.
Uprawnienia niewymienione nie są udzielane; w szczególności nie jest udzielane
prawo zwielokrotniania Oprogramowania w celu wprowadzenia do obrotu, prawo
najmu ani prawo użyczenia Egzemplarza.

Licencja nie przenosi na Licencjobiorcę żadnych autorskich praw majątkowych
ani praw zależnych do Oprogramowania. Prawa te pozostają przy Licencjodawcy.

### 4.2 Uprawnienia przyznane

W ramach Licencji Licencjobiorca jest uprawniony do:

1. **instalowania i uruchamiania** Oprogramowania na urządzeniach objętych
   zakresem podmiotowym Licencji (rozdz. 4.3), w celach wewnętrznej działalności
   Licencjobiorcy;
2. **korzystania z pełnej funkcjonalności** udostępnionej w danym Wydaniu,
   w granicach opisanych w niniejszej Licencji i dokumentacji;
3. **konfigurowania** Oprogramowania w zakresie przewidzianym mechanizmami
   samego produktu: definiowania kanałów modelu, kont i puli kont, punktów
   dostępu i nadań, dokumentów tożsamości modelu, katalogów roboczych, wartości
   ustawień na przewidzianych poziomach zasięgu i osiach rozstrzygania —
   przy czym uprawnienie to obejmuje **zapis i rozstrzyganie** wartości
   ustawień, nie zaś zapewnienie, że każda zapisana wartość zostanie
   zastosowana przy wywołaniu modelu; zakres wartości faktycznie stosowanych
   oraz ograniczenia zakładania kanałów i punktów dostępu ujawnia rozdz. 10.3
   (wykaz ograniczeń, pkt 21 i 22) oraz rozdz. 3.4 pkt 3;
4. **sporządzenia kopii zapasowej** Egzemplarza w liczbie niezbędnej do
   zabezpieczenia przed utratą, z zachowaniem wszystkich oznaczeń praw
   autorskich; kopia zapasowa nie może być używana równolegle z Egzemplarzem
   podstawowym;
5. **sporządzania kopii zapasowych Katalogu danych i Bazy** — bez ograniczeń
   liczbowych, przy czym Licencjobiorca odpowiada za zabezpieczenie tych kopii
   stosownie do wrażliwości zawartych w nich Danych Operatora;
6. **korzystania z wytworów pracy** — Treści Operatora powstałe przy użyciu
   Oprogramowania należą do Licencjobiorcy w granicach opisanych w rozdz. 6.5;
7. **wskazywania Operatorów** uprawnionych do obsługi Oprogramowania w imieniu
   Licencjobiorcy, w liczbie wynikającej z rozdz. 4.3, przy czym Licencjobiorca
   odpowiada za ich działania i zaniechania jak za własne.

### 4.3 Zakres podmiotowy — liczba stanowisk i użytkowników

> **[DO DECYZJI OPERATORA] — liczba stanowisk i model liczenia licencji.**
> Dokumentacja projektu nie zawiera rozstrzygnięcia o liczbie dozwolonych
> instalacji ani o jednostce, według której licencja jest liczona.
> Warianty:
> **(a)** licencja jednostanowiskowa — jedna Instalacja na jednym urządzeniu,
> obsługiwana przez jednego Operatora;
> **(b)** licencja imienna — jeden Operator, dowolna liczba urządzeń należących
> do Licencjobiorcy, przy zakazie równoczesnego korzystania;
> **(c)** licencja na urządzenie — jedna Instalacja, dowolna liczba Operatorów;
> **(d)** licencja zbiorowa z limitem — określona liczba równocześnie
> pracujących Operatorów, weryfikowana organizacyjnie;
> **(e)** licencja nielimitowana w obrębie organizacji Licencjobiorcy.
> Wariant (c) i (e) wymagają odrębnego rozstrzygnięcia dla topologii
> z Rdzeniem serwerowym, w której jedna Instalacja Rdzenia obsługuje wielu
> Operatorów. Do czasu rozstrzygnięcia obowiązuje zasada ostrożności: brak
> zgody Licencjodawcy na korzystanie wykraczające poza zakres uzgodniony
> indywidualnie.

Niezależnie od wybranego wariantu obowiązują zasady następujące:

1. Licencjobiorca prowadzi ewidencję Instalacji na potrzeby wykazania zgodności
   korzystania z zakresem Licencji.
2. Udostępnienie Oprogramowania podmiotowi spoza organizacji Licencjobiorcy —
   w tym w modelu świadczenia usług na rzecz osób trzecich przy użyciu
   Oprogramowania — wymaga uprzedniej zgody Licencjodawcy udzielonej na piśmie.
3. Korzystanie przez podmioty powiązane z Licencjobiorcą nie jest objęte
   Licencją, chyba że wynika to wprost z odrębnej umowy.

### 4.4 Zakres terytorialny i czasowy

Licencja jest udzielana bez ograniczeń terytorialnych, z zastrzeżeniem, że
Licencjobiorca zapewnia zgodność korzystania z przepisami obowiązującymi
w miejscu korzystania, w szczególności z przepisami o ochronie danych
osobowych, o kontroli eksportu oraz z regulacjami dotyczącymi systemów
sztucznej inteligencji.

Okres obowiązywania Licencji reguluje rozdz. 11.

### 4.5 Odpłatność

> **[DO DECYZJI OPERATORA] — odpłatność licencji.**
> Dokumentacja projektu nie rozstrzyga, czy Licencja jest odpłatna, a jeżeli
> tak — w jakim modelu. Warianty:
> **(a)** licencja nieodpłatna — korzystanie wewnętrzne w grupie Danaco;
> **(b)** opłata jednorazowa za Egzemplarz, z odrębnie płatnymi Aktualizacjami;
> **(c)** opłata abonamentowa (okresowa) obejmująca prawo korzystania
> i Aktualizacje w okresie abonamentu;
> **(d)** model mieszany — opłata podstawowa za Instalację plus opłata
> okresowa za wsparcie;
> **(e)** licencja nieodpłatna dla wersji o statusie Deweloperskim, odpłatna
> od pierwszego Wydania produkcyjnego.
> Wybór przesądza treść rozdz. 8.4 (uprawnienia do Aktualizacji), rozdz. 11.2
> (okres obowiązywania) i rozdz. 11.3 (przyczyny wypowiedzenia), a także
> konieczność uzupełnienia dokumentu o zasady fakturowania, waloryzacji
> i skutków opóźnienia w zapłacie.

Do czasu rozstrzygnięcia powyższej kwestii żadne postanowienie niniejszej
Licencji nie może być odczytywane jako zobowiązanie do zapłaty ani jako
zrzeczenie się prawa do wynagrodzenia.

### 4.6 Brak przeniesienia praw

Licencja nie stanowi sprzedaży Oprogramowania. Licencjobiorca nabywa wyłącznie
prawo do korzystania w zakresie opisanym w rozdz. 4.2. Wszelkie prawa
nieudzielone wprost pozostają zastrzeżone na rzecz Licencjodawcy.

Licencjobiorca nie może przenieść praw ani obowiązków wynikających z Licencji na
osobę trzecią — w drodze czynności prawnej, wniesienia aportem, połączenia,
podziału ani przejęcia — bez uprzedniej zgody Licencjodawcy udzielonej na
piśmie pod rygorem nieważności.

### 4.7 Zakres wykonania modelu a zakres licencji

Zgodnie z przyjętą architekturą zakres wykonania modelu (`local`, `core`,
`remote`) jest parametrem Okna komunikacji i pozostaje niezależny od
umiejscowienia Rdzenia.
Licencja nie ogranicza Operatora co do wyboru maszyn, na których model pracuje,
z zastrzeżeniem, że:

1. Licencjobiorca musi posiadać tytuł prawny do korzystania z maszyn i zasobów
   udostępnianych modelowi przez Punkty dostępu i Katalogi robocze;
2. udostępnienie modelowi zasobów osoby trzeciej wymaga jej zgody;
3. w wersji v2.0 wybór zakresu wykonania **jest wyłącznie etykietą w prowenancji
   i nie zmienia miejsca uruchomienia procesu** — proces modelu startuje zawsze
   na maszynie Rdzenia, a kanał połączenia zdalnego (SSH) jako kanał modelu nie
   istnieje. Operator nie może zatem opierać się na wartości zakresu wykonania
   jako na mechanizmie ograniczającym; okoliczność ta ma bezpośredni wpływ na
   ocenę ryzyka udostępniania zasobów i została ujawniona wprost.

---

## 5. Zasady korzystania i ograniczenia

### 5.1 Zasada ogólna

Licencjobiorca korzysta z Oprogramowania zgodnie z jego przeznaczeniem,
dokumentacją i niniejszą Licencją, w sposób nienaruszający praw Licencjodawcy,
praw osób trzecich oraz przepisów prawa. Wątpliwość co do zakresu dozwolonego
korzystania rozstrzyga się przez zapytanie skierowane do Licencjodawcy, nie
przez samodzielne rozszerzającą wykładnię uprawnień.

### 5.2 Ograniczenia dotyczące samego Oprogramowania

Licencjobiorca zobowiązuje się **nie**:

1. zwielokrotniać Oprogramowania w zakresie wykraczającym poza czynności
   niezbędne do jego zainstalowania, uruchomienia i sporządzenia kopii
   zapasowej;
2. rozpowszechniać Oprogramowania ani jego części — odpłatnie ani nieodpłatnie —
   w szczególności przez udostępnienie plików wykonywalnych, obrazów
   instalacyjnych, kontenerów ani obrazów maszyn wirtualnych zawierających
   Oprogramowanie;
3. najmować, użyczać, oddawać w leasing ani udostępniać Oprogramowania osobom
   trzecim w modelu współdzielenia dostępu, hostingu ani świadczenia usług
   przy jego użyciu na rzecz osób trzecich;
4. dokonywać dekompilacji, dezasemblacji ani innej postaci odtwarzania kodu
   źródłowego, z zastrzeżeniem rozdz. 5.7;
5. modyfikować Oprogramowania, tworzyć jego opracowań, tłumaczeń ani adaptacji,
   w tym modyfikować plików wykonywalnych, zasobów wbudowanych, migracji
   schematu ani warstwy wizualnej;
6. usuwać, zasłaniać ani zmieniać oznaczeń praw autorskich, znaków towarowych,
   not licencyjnych, oznaczeń wersji ani statusu wersji, zarówno w interfejsie
   Oprogramowania, jak i w plikach dystrybucji;
7. obchodzić ani wyłączać mechanizmów Oprogramowania służących identyfikacji
   wersji, prowenancji wywołania modelu ani zapisowi historii;
8. **[DO DECYZJI OPERATORA]** wykorzystywać Oprogramowania, jego architektury,
   Kontraktu, modelu danych ani warstwy wizualnej do zbudowania produktu
   konkurencyjnego wobec Danaco Console;
9. **[DO DECYZJI OPERATORA]** przeprowadzać testów penetracyjnych, testów
   obciążeniowych ani badań bezpieczeństwa Oprogramowania bez uprzedniej zgody
   Licencjodawcy udzielonej na piśmie; ujawnienie ustalonych podatności
   następuje wyłącznie Licencjodawcy.

> **[DO DECYZJI OPERATORA] — zakaz budowy produktu konkurencyjnego
> (pkt 8 powyżej).**
> Dokumentacja projektu — w tym zapis decyzji architektonicznych — nie zawiera
> rozstrzygnięcia o objęciu Licencjobiorcy klauzulą niekonkurencji.
> Warianty: **(a)** brak klauzuli — ochrona wyłącznie środkami prawa
> autorskiego, prawa znaków towarowych i przepisów o zwalczaniu nieuczciwej
> konkurencji; **(b)** zakaz wąski — wyłącznie zakaz odwzorowania Kontraktu,
> modelu danych i warstwy wizualnej, bez zakazu tworzenia produktów tej samej
> kategorii; **(c)** zakaz szeroki w brzmieniu pkt 8, ograniczony czasowo
> i przedmiotowo; **(d)** zakaz szeroki bez ograniczenia czasowego.
> Warianty (c) i (d) wymagają oceny dopuszczalności w świetle prawa
> właściwego (rozdz. 14.1) — klauzula niekonkurencji o nieoznaczonym zakresie
> bywa uznawana za nieskuteczną. Do czasu rozstrzygnięcia pkt 8 należy
> odczytywać w brzmieniu wariantu (b), to jest jako zakaz odwzorowania
> chronionych elementów Oprogramowania wskazanych w rozdz. 6.1.

> **[DO DECYZJI OPERATORA] — badania bezpieczeństwa i testy obciążeniowe
> (pkt 9 powyżej).**
> Dokumentacja projektu nie rozstrzyga, czy Licencjobiorca może samodzielnie
> badać bezpieczeństwo Oprogramowania. Kwestia ma wagę praktyczną: rozdz. 5.5
> pkt 3 ujawnia, że uwierzytelnianie dostępu do Rdzenia jest w fazie budowy
> wyłączone, a rozdz. 13.4 nakłada obowiązek zgłaszania podatności — obowiązek
> ten trudno pogodzić z zakazem ich wyszukiwania. Warianty: **(a)** brak
> zakazu — Licencjobiorca bada Oprogramowanie we własnym środowisku i zgłasza
> ustalenia trybem rozdz. 13.4; **(b)** dopuszczenie badań we własnym
> środowisku Licencjobiorcy z zakazem badania instalacji cudzych i z obowiązkiem
> zgłoszenia ustaleń; **(c)** zakaz w brzmieniu pkt 9 — badania wyłącznie za
> zgodą na piśmie; **(d)** zakaz wyłącznie testów obciążeniowych obciążających
> zasoby wspólne (topologia z Rdzeniem serwerowym), przy swobodzie badań
> bezpieczeństwa. Wybór przesądza spójność rozdz. 5.2 pkt 9 z rozdz. 13.4 i z
> polityką zgłaszania podatności.

### 5.3 Ograniczenia wynikające ze statusu Deweloperskiego

Ze względu na status Deweloperski wersji v2.0 Licencjobiorca zobowiązuje się
**nie stosować** Oprogramowania:

1. w procesach, w których błąd, przerwanie pracy albo utrata danych mogą
   spowodować zagrożenie życia, zdrowia albo bezpieczeństwa osób;
2. w systemach sterowania infrastrukturą krytyczną, urządzeniami medycznymi,
   środkami transportu ani w innych zastosowaniach o podwyższonym ryzyku;
3. jako jedynego miejsca przechowywania danych o znaczeniu istotnym dla
   działalności Licencjobiorcy — bez niezależnej kopii zapasowej prowadzonej
   poza Katalogiem danych;
4. jako narzędzia podejmującego automatycznie decyzje wywołujące skutki prawne
   wobec osób fizycznych albo w podobny sposób istotnie na nie wpływające, bez
   udziału człowieka.

Ograniczenia te wynikają z rzeczywistego stanu wykonania produktu opisanego
w rozdz. 10.3 i pozostają w mocy do czasu wydania przez Licencjodawcę wersji
o statusie innym niż Deweloperski.

### 5.4 Zasady korzystania z modelu i odpowiedzialność za polecenia

1. Oprogramowanie jest narzędziem pośredniczącym: przekazuje modelowi treść
   wskazaną przez Operatora wraz z Nakładką tożsamości, ustawieniami Okna
   komunikacji i konfiguracją mostów, po czym zwraca strumień odpowiedzi.
   **Licencjodawca nie tworzy, nie weryfikuje ani nie zatwierdza treści
   wytwarzanych przez model.**
2. Operator odpowiada za treść poleceń wydawanych modelowi oraz za
   wykorzystanie wyników jego pracy, w tym za weryfikację ich poprawności
   przed użyciem.
3. Operator odpowiada za zakres uprawnień nadanych modelowi: za wybór trybu
   uprawnień, listę Katalogów roboczych, Punkty dostępu i Nadania (`read`
   albo `write`) oraz za konfigurację mostów MCP. Operator przyjmuje do
   wiadomości, że nadanie trybu `write` oznacza możliwość dokonywania przez
   model zmian w plikach na maszynach objętych nadaniem.
4. Operator przyjmuje do wiadomości, że **w wersji v2.0 mechanizmy izolacji
   opisane dokumentacją (jedenaście ustawień izolacji rozstrzyganych na ośmiu
   poziomach zasięgu) są zdefiniowane i rozstrzygane, lecz nie mają
   egzekutorów — nie wpływają na sposób uruchomienia procesu modelu.**
   Ustawienia izolacji nie mogą być traktowane jako zabezpieczenie techniczne.
5. Operator przyjmuje do wiadomości, że **w wersji v2.0 model nie zachowuje
   ciągłości rozmowy między turami**: każda tura uruchamia nowy proces bez
   historii poprzednich wiadomości Okna komunikacji. Okoliczność ta ma wpływ
   na sposób formułowania poleceń i na wynik pracy.

### 5.5 Obowiązki Licencjobiorcy w zakresie bezpieczeństwa

1. Licencjobiorca zabezpiecza urządzenie, na którym działa Rdzeń, w sposób
   adekwatny do wrażliwości Danych Operatora zapisanych w Bazie.
2. Licencjobiorca odpowiada za ochronę Poświadczeń wskazywanych Oprogramowaniu
   przez odwołania (klucze API, klucze prywatne SSH, katalogi konfiguracji
   kont). Oprogramowanie nie przechowuje treści Poświadczeń w Bazie — ich
   bezpieczeństwo pozostaje po stronie Licencjobiorcy.
3. Licencjobiorca przyjmuje do wiadomości, że mechanizm uwierzytelniania jest
   **w fazie budowy wyłączony**, a rozstrzyganie dostępu do Rdzenia opiera się
   wyłącznie na zabezpieczeniu środowiska, w którym Rdzeń działa.
   W konsekwencji **każdy podmiot mający dostęp sieciowy do portu nasłuchu
   Rdzenia ma dostęp do pełnej funkcjonalności Oprogramowania**.
   Ograniczenie dostępu do tego portu należy do Licencjobiorcy.
4. Licencjobiorca nie udostępnia portu nasłuchu Rdzenia w sieci publicznej.

### 5.6 Zgodność z warunkami Dostawców modeli

Korzystanie z Oprogramowania w sposób prowadzący do wywołania modelu podlega
równolegle warunkom umownym Dostawcy modelu. Licencjobiorca zobowiązuje się:

1. przestrzegać regulaminów i limitów Dostawcy modelu, w tym zakazów dotyczących
   współdzielenia kont, automatyzacji i obchodzenia limitów;
2. nie wykorzystywać mechanizmu Puli kont i rotacji Kont w sposób sprzeczny
   z warunkami Dostawcy modelu; mechanizm ten służy uporządkowanej pracy przy
   wielu uprawnionych kontach, nie zaś obchodzeniu ograniczeń;
3. zapewnić, że treści przekazywane modelowi mogą być temu modelowi przekazane
   zgodnie z prawem i z zobowiązaniami Licencjobiorcy wobec osób trzecich.

Licencjodawca nie jest stroną stosunku prawnego między Licencjobiorcą
a Dostawcą modelu i nie odpowiada za jego treść, wykonanie ani skutki
zakończenia.

### 5.7 Uprawnienia wynikające z przepisów bezwzględnie obowiązujących

Ograniczenia opisane w rozdz. 5.2 pkt 4 i 5 nie naruszają uprawnień
przysługujących Licencjobiorcy na podstawie bezwzględnie obowiązujących
przepisów prawa właściwego, w szczególności — jeżeli prawem właściwym będzie
prawo polskie — uprawnień do zwielokrotniania kodu i tłumaczenia jego formy
w zakresie niezbędnym do uzyskania współdziałania z innymi programami, na
zasadach i przy zachowaniu warunków określonych w art. 75 ustawy z dnia
4 lutego 1994 r. o prawie autorskim i prawach pokrewnych. Licencjobiorca
zamierzający skorzystać z tych uprawnień zwraca się uprzednio do Licencjodawcy
o udostępnienie informacji niezbędnych do współdziałania; Licencjodawca
udostępnia je w rozsądnym terminie, o ile pozostaje to w jego możliwościach.

### 5.8 Skutki naruszenia ograniczeń

Naruszenie ograniczeń opisanych w rozdz. 5.2, 5.3 albo 5.6 stanowi istotne
naruszenie Licencji i uprawnia Licencjodawcę do jej wypowiedzenia w trybie
rozdz. 11.3, niezależnie od dalej idących roszczeń.

> **[DO DECYZJI OPERATORA] — sankcje umowne za naruszenie.**
> Dokumentacja projektu nie rozstrzyga, czy Licencja ma przewidywać kary
> umowne za naruszenie ograniczeń (w szczególności za rozpowszechnianie
> Oprogramowania), a jeżeli tak — w jakiej wysokości i za jakie kategorie
> naruszeń. Warianty: **(a)** brak kar umownych, odpowiedzialność wyłącznie
> na zasadach ogólnych; **(b)** kara umowna zryczałtowana za każde naruszenie;
> **(c)** kara umowna powiązana z wartością licencji, z zastrzeżeniem prawa
> dochodzenia odszkodowania przenoszącego jej wysokość.

---

## 6. Własność intelektualna

### 6.1 Prawa do Oprogramowania

Oprogramowanie stanowi utwór w rozumieniu prawa autorskiego. Autorskie prawa
majątkowe do Oprogramowania przysługują Licencjodawcy — Danaco Holding Group
Sp. z o.o. Autorskie prawa osobiste przysługują Twórcy — Dariuszowi
Naharnowiczowi — i nie podlegają zrzeczeniu ani przeniesieniu.

Ochronie podlega całość Oprogramowania oraz jego elementy dające się
wyodrębnić, w szczególności:

1. kod źródłowy i kod wynikowy wszystkich trzech warstw (Rdzeń, Klient,
   Powłoka);
2. **Kontrakt** — dobór, nazewnictwo i struktura komend, zdarzeń, strumieni
   i deklaracji narzędzi, wraz z konwencją generowania wiązań;
3. **model danych** — struktura schematu trwałości, dobór i nazewnictwo tabel,
   kolumn oraz łańcuch więzów sesja → karta → okno → wiadomość;
4. **warstwa wizualna** — zestaw żetonów motywu, dwa motywy (jasny i ciemny),
   siatka i zasady zestawu ikon, kompozycja widoków;
5. **oznaczenia marki** — logotypy, sygnety i ich warianty;
6. dokumentacja produktu.

Struktura katalogów, nazewnictwo modułów oraz konwencje architektoniczne
Oprogramowania stanowią wyraz autorskiego doboru i również podlegają ochronie
w zakresie, w jakim mają charakter twórczy.

### 6.2 Zastrzeżenie praw

Wszelkie prawa nieudzielone Licencjobiorcy wprost w rozdz. 4.2 pozostają
zastrzeżone na rzecz Licencjodawcy. Milczenie Licencji co do określonego pola
eksploatacji oznacza brak zgody, nie zaś zgodę dorozumianą.

### 6.3 Zakaz działań naruszających

Licencjobiorca zobowiązuje się nie podejmować działań zmierzających do
podważenia praw Licencjodawcy, w szczególności nie rejestrować na swoją rzecz
oznaczeń zbieżnych z oznaczeniami Oprogramowania, nazw domen zawierających
oznaczenie „Danaco Console” ani nie zgłaszać do ochrony rozwiązań
odwzorowujących architekturę albo model danych Oprogramowania.

### 6.4 Znaki towarowe i oznaczenia

Oznaczenia „Danaco”, „Danaco Console”, „Danaco Holding Group”, nazwy Środowisk
(TalkIn, WorkSpace, CodeStudio, MultitaskingAI) oraz logotypy i sygnety produktu
stanowią oznaczenia Licencjodawcy. Licencja **nie przyznaje** prawa do
używania tych oznaczeń poza zakresem niezbędnym do zwykłego korzystania
z Oprogramowania i wewnętrznego odwoływania się do niego.

Użycie oznaczeń w materiałach marketingowych, w komunikacji publicznej,
w nazwach produktów Licencjobiorcy albo w sposób sugerujący powiązanie
gospodarcze z Licencjodawcą wymaga uprzedniej zgody udzielonej na piśmie.

> **[DO DECYZJI OPERATORA] — status ochronny znaków.**
> Dokumentacja projektu nie zawiera informacji o zgłoszeniu albo rejestracji
> znaków towarowych „Danaco Console” ani oznaczeń graficznych produktu.
> Warianty: **(a)** powołanie się wyłącznie na ochronę oznaczeń
> nierejestrowanych i prawo autorskie do logotypów; **(b)** zgłoszenie znaków
> do rejestracji przed publikacją i uzupełnienie dokumentu o numery praw
> ochronnych; **(c)** rezygnacja z rozdziału o znakach towarowych.

### 6.5 Prawa do Treści Operatora

1. Licencjodawca nie nabywa praw do Treści Operatora ani do Danych Operatora.
   Prawa te — w zakresie, w jakim powstają — przysługują Licencjobiorcy albo
   podmiotom uprawnionym zgodnie z prawem właściwym.
2. Licencjodawca nie uzyskuje dostępu do Danych Operatora w związku z samym
   korzystaniem przez Licencjobiorcę z Oprogramowania. Dane te pozostają
   w Katalogu danych na urządzeniu Licencjobiorcy (rozdz. 9).
3. Licencjodawca nie składa żadnego oświadczenia co do tego, czy i w jakim
   zakresie wytwór wygenerowany przez model podlega ochronie prawnoautorskiej,
   ani co do tego, komu przysługują do niego prawa. Ocena ta zależy od prawa
   właściwego, sposobu powstania wytworu i wkładu twórczego człowieka
   i pozostaje po stronie Licencjobiorcy.
4. Licencjobiorca odpowiada za to, aby korzystanie z wytworów pracy modelu nie
   naruszało praw osób trzecich.

### 6.6 Informacje zwrotne

> **[DO DECYZJI OPERATORA] — status informacji zwrotnych.**
> Dokumentacja projektu nie rozstrzyga, na jakich zasadach Licencjodawca może
> korzystać z uwag, zgłoszeń błędów i propozycji funkcji przekazywanych przez
> Licencjobiorcę. Warianty: **(a)** swobodne, nieodpłatne i nieograniczone
> korzystanie przez Licencjodawcę, bez roszczeń zgłaszającego; **(b)**
> korzystanie za zgodą zgłaszającego, w zakresie każdorazowo uzgodnionym;
> **(c)** brak regulacji — zastosowanie zasad ogólnych.
> Wybór wariantu (a) wymaga wyraźnego oświadczenia Licencjobiorcy o udzieleniu
> nieodpłatnej licencji na korzystanie ze zgłoszenia; wariant ten bywa
> kwestionowany, jeżeli zgłoszenie zawiera rozwiązanie o charakterze twórczym.

---

## 7. Komponenty i rozszerzenia

### 7.1 Zasada nadrzędna

Oprogramowanie zawiera oraz wykorzystuje komponenty osób trzecich
(„**Komponenty**”), rozpowszechniane na własnych licencjach. Do każdego
Komponentu stosuje się **jego własną licencję**; niniejsza Licencja nie
ogranicza uprawnień przyznanych Licencjobiorcy przez te licencje ani nie
rozszerza ich na kod autorski Licencjodawcy.

W razie sprzeczności między postanowieniem niniejszej Licencji a warunkiem
licencji Komponentu — w zakresie dotyczącym tego Komponentu — pierwszeństwo ma
licencja Komponentu.

Zestawienie Komponentów w niniejszym rozdziale sporządzono na podstawie
faktycznej zawartości repozytorium w stanie bieżącym: plików deklaracji
zależności (`budowa/go.mod`, `budowa/client/package.json`,
`budowa/desktop/src-tauri/Cargo.toml`) oraz plików zamknięcia zależności
(`budowa/client/package-lock.json`, `budowa/desktop/src-tauri/Cargo.lock`).
Zestawienie warstwy Rdzenia jest niezmienione względem rewizji pierwotnej
stanu sprzed prac naprawczych; zestawienie warstwy Klienta uwzględnia dodanie zależności
narzędziowej `jsdom` (środowisko DOM podkładane pod sprawdziany widoków,
prace nad sprawdzianami widoków) — patrz rozdz. 7.3, pozostałe pozycje warstwy
Klienta są niezmienione; zestawienie warstwy Powłoki uwzględnia usunięcie
nieużywanej zależności bezpośredniej `serde_json` z `Cargo.toml`, dokonane
w pracach nad skryptami wydania (prace nad skryptami wydania) — patrz rozdz. 7.4. Zestawienie
zbiorcze zawiera Załącznik A.

### 7.2 Komponenty warstwy Rdzenia (Go)

Rdzeń jest budowany jako moduł Go `danacoconsole` (wymagana wersja języka:
`go 1.26`). Zależności bezpośrednie i pośrednie deklarowane w `budowa/go.mod`
oraz rodzaje ich licencji ustalone z treści plików licencyjnych
w lokalnej pamięci podręcznej modułów:

| Komponent | Wersja | Rodzaj licencji | Rola w produkcie |
|---|---|---|---|
| `github.com/coder/websocket` | v1.8.15 | ISC | transport WebSocket Rdzenia |
| `golang.org/x/sys` | v0.47.0 | BSD 3-Clause | dostęp do funkcji systemowych |
| `modernc.org/sqlite` | v1.56.0 | BSD 3-Clause | silnik SQLite bez zależności natywnych |
| `github.com/dustin/go-humanize` | v1.0.1 | MIT | zależność pośrednia silnika SQLite |
| `github.com/google/uuid` | v1.6.0 | BSD 3-Clause | zależność pośrednia silnika SQLite |
| `github.com/mattn/go-isatty` | v0.0.24 | MIT | zależność pośrednia silnika SQLite |
| `github.com/ncruces/go-strftime` | v1.0.0 | MIT | zależność pośrednia silnika SQLite |
| `github.com/remyoudompheng/bigfft` | v0.0.0-20230129092748-24d4a6f8daec | BSD 3-Clause | zależność pośrednia silnika SQLite |
| `modernc.org/libc` | v1.74.4 | BSD 3-Clause | zależność pośrednia silnika SQLite |
| `modernc.org/mathutil` | v1.7.1 | BSD 3-Clause | zależność pośrednia silnika SQLite |
| `modernc.org/memory` | v1.11.0 | BSD 3-Clause | zależność pośrednia silnika SQLite |

Wszystkie wymienione licencje są licencjami zezwalającymi, dopuszczającymi
rozpowszechnianie w postaci wynikowej pod warunkiem zachowania not
o prawach autorskich i tekstu licencji. Licencjodawca spełnia ten warunek
przez dołączenie zestawienia z Załącznika A wraz z pełnymi tekstami licencji do
Wydania (rozdz. 7.9).

### 7.3 Komponenty warstwy Klienta (TypeScript)

Zależnością **produkcyjną** Klienta — to znaczy taką, której kod trafia do
pakietu dostarczanego Licencjobiorcy — jest wyłącznie:

| Komponent | Wersja | Rodzaj licencji | Rola w produkcie |
|---|---|---|---|
| `@tauri-apps/api` | 2.5.0 | Apache-2.0 **lub** MIT (do wyboru) | wywołanie poleceń Powłoki z poziomu interfejsu |

Pozostałe zależności Klienta mają charakter **narzędzi budowania i testowania**
i nie wchodzą do pakietu dostarczanego Licencjobiorcy:

| Komponent | Wersja | Rodzaj licencji | Rola |
|---|---|---|---|
| `typescript` | 5.8.3 | Apache-2.0 | kontrola typów (`tsc --noEmit`) |
| `vite` | 6.3.5 | MIT | budowanie pakietu interfejsu |
| `vitest` | 3.2.7 | MIT | sprawdziany jednostkowe interfejsu |
| `jsdom` | ^29.1.1 | MIT | środowisko DOM dla sprawdzianów widoków (prace nad sprawdzianami widoków) |

Zamknięcie zależności zapisane w `package-lock.json` obejmuje **141 pozycji**.
Rozkład rodzajów licencji w tym zamknięciu, odczytany z pól `license` pliku
zamknięcia: MIT — 126 pozycji, Apache-2.0 — 3, ISC — 3, MIT-0 — 2,
BSD-2-Clause — 2, BSD-3-Clause — 2, „Apache-2.0 OR MIT” — 1, BlueOak-1.0.0 — 1,
CC0-1.0 — 1. Zestawienie nie ujawnia licencji o charakterze wzajemnym
(copyleft) w tej warstwie.

### 7.4 Komponenty warstwy Powłoki (Rust / Tauri)

Zależności zadeklarowane w `budowa/desktop/src-tauri/Cargo.toml`, z rozróżnieniem
sekcji `[dependencies]` (zależność produkcyjna, wchodząca do pliku wynikowego)
i `[build-dependencies]` (narzędzie etapu budowania, nieobecne w pliku
wynikowym) — rozróżnienie prowadzone konsekwentnie z rozdz. 7.3 dla warstwy
Klienta:

| Komponent | Sekcja `Cargo.toml` | Wersja wymagana | Rodzaj licencji | Rola w produkcie |
|---|---|---|---|---|
| `tauri` | `[dependencies]` | 2 (rozwiązana: 2.11.5) | Apache-2.0 **lub** MIT | powłoka okna, ikona zasobnika, obsługa obrazów ikon |
| `tauri-plugin-dialog` | `[dependencies]` | 2 (rozwiązana: 2.7.2) | Apache-2.0 **lub** MIT | okno wyboru katalogu roboczego |
| `serde` | `[dependencies]` | 1 (rozwiązana: 1.0.229) | MIT **lub** Apache-2.0 | serializacja danych poleceń powłoki |
| `tauri-build` | **`[build-dependencies]`** — zależność budowania | 2 (rozwiązana: 2.6.3) | Apache-2.0 **lub** MIT | krok budowania powłoki; nie wchodzi do pliku wynikowego |

**Zmiana względem stanu sprzed prac naprawczych.** Zależność bezpośrednia
`serde_json` została usunięta z `Cargo.toml` w pracach nad ścieżką wydania i bramami mierzącymi wywołania jako nieużywana.
`serde_json` w wersji 1.0.151 nadal
występuje w zamknięciu `Cargo.lock` jako zależność **pośrednia**, wciągana
przez `tauri`; nota licencyjna MIT **lub** Apache-2.0 pozostaje należna z tego
tytułu, lecz komponent nie jest już zależnością bezpośrednią warstwy Powłoki
i nie figuruje w tabeli powyżej.

Zamknięcie zależności zapisane w `Cargo.lock` obejmuje **447 pakietów** —
jest to konsekwencja architektury Tauri, która wciąga zależności warstwy
systemowej dla wielu platform. Rozkład rodzajów licencji ustalony przez odczyt
pola `license` z manifestów pakietów obecnych w lokalnej pamięci podręcznej
rejestru (**274 z 447 pakietów zweryfikowane lokalnie; 173 pakiety
niezweryfikowane** — w większości pakiety właściwe dla platform innych niż
Windows oraz pakiet własny Powłoki):

| Zadeklarowana licencja | Liczba pakietów |
|---|---|
| MIT OR Apache-2.0 (w różnych zapisach) | 144 |
| MIT | 47 |
| Apache-2.0 OR MIT | 31 |
| Unicode-3.0 | 18 |
| Unlicense OR MIT (w różnych zapisach) | 11 |
| **MPL-2.0** | **5** |
| BSD-3-Clause i warianty z BSD-3-Clause | 6 |
| Zlib i warianty z Zlib | 5 |
| Apache-2.0 | 2 |
| pozostałe warianty złożone | 5 |

**Uwaga o licencjach wzajemnych.** W zamknięciu zależności Powłoki występuje
pięć pakietów na licencji **MPL-2.0** (Mozilla Public License 2.0):
`cssparser` 0.36.0, `cssparser-macros` 0.6.1, `dtoa-short` 0.3.5,
`option-ext` 0.2.0 oraz `selectors` 0.36.1. MPL-2.0 jest licencją wzajemną
o zasięgu plikowym: obowiązek udostępnienia kodu źródłowego dotyczy plików
objętych tą licencją i ich modyfikacji, nie zaś całości produktu, z którym
zostały połączone. Licencjodawca nie modyfikuje tych pakietów, w związku z czym
obowiązek sprowadza się do wskazania miejsca uzyskania ich kodu źródłowego
i zachowania not licencyjnych.

**Zastrzeżenie o kompletności.** Wykaz warstwy Rust jest zweryfikowany
częściowo. Przed publikacją Wydania konieczne jest sporządzenie pełnego
zestawienia licencji zamknięcia zależności właściwego dla budowy docelowej
(narzędziem inwentaryzującym uruchomionym na maszynie budującej) i dołączenie
go do Wydania.

> **[DO DECYZJI OPERATORA] — sposób prowadzenia inwentarza licencji.**
> Warianty: **(a)** inwentarz generowany automatycznie przy każdym Wydaniu
> i dołączany jako plik `NOTICE`; **(b)** inwentarz utrzymywany ręcznie
> i weryfikowany przy zmianie zależności; **(c)** inwentarz publikowany
> w dokumentacji online. Wariant (a) jest jedynym odpornym na rozjazd przy
> aktualizacji zależności.

### 7.5 Kroje pisma

Warstwa wizualna Oprogramowania opiera się na trzech krojach pisma, których
pliki w formacie WOFF2 są **wbudowane w Oprogramowanie
i rozpowszechniane wraz z nim**:

| Krój | Rola w produkcie | Licencja | Uprawniony |
|---|---|---|---|
| **Space Grotesk** | nagłówki, tytuły środowisk, logotyp | SIL Open Font License 1.1 | Copyright 2020 The Space Grotesk Project Authors |
| **IBM Plex Sans** | interfejs i treść | SIL Open Font License 1.1 | Copyright 2019 IBM Corp. |
| **IBM Plex Mono** | dane techniczne, identyfikatory, terminal | SIL Open Font License 1.1 | Copyright 2019 IBM Corp. |

Kroje osadzono w podzbiorach `latin` oraz `latin-ext`, obejmujących pełny
zestaw polskich znaków diakrytycznych. Pełne teksty licencji obu rodzin
znajdują się w plikach dołączonych do zasobów Oprogramowania
(`LICENCJA-space-grotesk.txt`, `LICENCJA-ibm-plex.txt`).

Z licencji SIL OFL 1.1 wynikają obowiązki, które Licencjobiorca przyjmuje do
wiadomości i zobowiązuje się respektować:

1. pliki krojów mogą być rozpowszechniane wyłącznie razem z Oprogramowaniem
   albo samodzielnie, **zawsze wraz z tekstem licencji i notą o prawach
   autorskich**;
2. **zakazana jest sprzedaż samych plików krojów** jako odrębnego towaru;
3. zmodyfikowana wersja kroju nie może być rozpowszechniana pod nazwą
   zastrzeżoną (Reserved Font Name) przez autora oryginału;
4. licencja OFL dotyczy plików krojów, nie zaś dokumentów ani grafik
   złożonych przy ich użyciu — wytwory Operatora nie stają się przez to
   objęte OFL.

Nieobowiązujące — wyłączone decyzją projektową i niewystępujące
w Oprogramowaniu — są kroje Cormorant Garamond, Inter oraz JetBrains Mono.

### 7.6 Zestaw ikon i zasoby marki

Zestaw ikon interfejsu opisany jest manifestem warstwy wizualnej w wersji v2.0
(`ikony/manifest.json`) i obejmuje **82 pozycje** na jednolitej siatce 24×24
z obrysem 1,75 i barwą dziedziczoną z kontekstu (`currentColor`). Liczba ta
wynika wprost z manifestu (pole `liczba-ikon` o wartości 82 oraz 82 wpisy
wykazu) i odpowiada wiązaniom zestawu w kodzie Klienta (13 + 16 + 11 + 14 + 24 +
4 = 82 pozycje w sześciu grupach zastosowań). Źródła pozycji, ustalone z pola
`zrodlo` każdego wpisu manifestu:

1. **78 pozycji** wywiedzionych z biblioteki **Lucide**, rozpowszechnianej na
   licencji **ISC**, z obrysem ujednoliconym do wartości przyjętej w zestawie;
2. **4 pozycje własne** — emblematy Środowisk (`srodowisko-talkin`,
   `srodowisko-workspace`, `srodowisko-codestudio`, `srodowisko-multitaskingai`)
   wykonane w siatce i kresce zestawu; stanowią utwór Licencjodawcy.

Licencja ISC jest licencją zezwalającą; jej warunkiem jest zachowanie noty
o prawach autorskich i treści zezwolenia we wszystkich kopiach.

**Uzgodnienie z liczbą 83 pojawiającą się w katalogu plików zestawu.**
Katalog plików SVG zestawu (`client/src/ikony/svg/`) zawiera **83 pliki**, lecz
jednym z nich jest `logo-danaco.svg` — znak marki wczytywany w kodzie odrębnie
i **nienależący do zestawu ikon** ani do wykazu manifestu. Dla celów
licencyjnych rozstrzyga podział rzeczowy, nie liczba plików katalogu: **82
pozycje zestawu ikon** (78 na licencji ISC, 4 własne) oraz **znak marki**
podlegający rozdz. 6.4 wraz z pozostałymi zasobami marki.

Redakcja niniejszej Licencji zweryfikowała bezpośrednio treść [README.md](README.md)
i [INSTRUKCJA-UZYTKOWANIA.md](INSTRUKCJA-UZYTKOWANIA.md) dołączonych do Egzemplarza: oba dokumenty
posługują się liczbą **82 pozycji zestawu ikon** i wprost objaśniają
czytelnikowi ten sam podział rzeczowy — 82 pozycje wykazu manifestu wobec 83
plików katalogu z uwzględnieniem znaku marki. Wcześniejsza redakcja niniejszego
rozdziału odnotowywała rozbieżność liczby 83 przypisywaną obu tym dokumentom;
rozbieżność ta nie istnieje w treści dokumentów na dzień sporządzenia
niniejszej wersji Licencji i została usunięta z treści rozdziału jako
nieaktualna. Kwestia jednolitej liczby pozycji zestawu ikon w dokumentacji
produktu nie wymaga już rozstrzygnięcia Operatora — wszystkie trzy dokumenty
(niniejsza Licencja, [README.md](README.md), [INSTRUKCJA-UZYTKOWANIA.md](INSTRUKCJA-UZYTKOWANIA.md)) zgodnie
posługują się liczbą wynikającą z manifestu jako źródła prawdy warstwy
wizualnej, a odpowiadająca temu pozycja została usunięta z Załącznika B.

**Zasoby marki** — logotypy, sygnety i ich warianty dla motywu jasnego
i ciemnego (dwanaście plików katalogu zasobów marki) oraz znak
`logo-danaco.svg` — stanowią utwór Licencjodawcy i nie podlegają licencjom
komponentów. Ich użycie reguluje rozdz. 6.4.

### 7.7 Program zewnętrzny kanału głównego

Kanał modelu rodzaju `cli` uruchamia zewnętrzny program `claude`
(Claude Code CLI). Program ten **nie jest częścią Oprogramowania i nie jest
rozpowszechniany przez Licencjodawcę**. Oprogramowanie:

1. buduje argumenty wywołania warunkowo, zależnie od konfiguracji Okna
   komunikacji i Nakładki tożsamości (`--system-prompt` albo
   `--append-system-prompt`, `--permission-mode`, `--add-dir`, `--mcp-config`);
2. ustawia dla procesu zmienną `CLAUDE_CONFIG_DIR` wskazującą katalog
   konfiguracji Konta;
3. odczytuje strumień odpowiedzi w formacie `stream-json` i przekazuje go
   Klientowi w postaci fragmentów.

Korzystanie z programu `claude` oraz z usług Dostawcy modelu podlega wyłącznie
warunkom umownym tego Dostawcy. Licencjobiorca zapewnia legalność korzystania
z programu i konta.

### 7.8 Rozszerzenia

Przyjęta decyzja architektoniczna przewiduje katalog rozszerzeń (przeglądanie,
instalacja, zarządzanie, kategorie) oparty na tabeli `rozszerzenie` i rodzinie
komend
`extension.*`, przy czym marketplace jako model wydawniczy (publikowanie
i dystrybucja rozszerzeń do innych operatorów) został odłożony do odrębnego
opracowania.

**W wersji v2.0 mechanizm rozszerzeń nie istnieje**: rodzina komend
`extension.*` nie występuje w Kontrakcie wykonawczym, a Oprogramowanie nie
udostępnia drogi instalowania ani uruchamiania rozszerzeń. Funkcjonalność ta
jest przewidziana koncepcją i w bieżącej wersji niezaimplementowana.

W konsekwencji Licencja **nie reguluje** zasad tworzenia, dystrybucji ani
korzystania z rozszerzeń.

> **[DO DECYZJI OPERATORA] — zasady licencyjne rozszerzeń.**
> Przed udostępnieniem mechanizmu rozszerzeń konieczne jest rozstrzygnięcie:
> **(a)** czy rozszerzenia będą mogły pochodzić wyłącznie od Licencjodawcy,
> czy także od osób trzecich; **(b)** na jakiej licencji będzie udostępniane
> ewentualne API rozszerzeń; **(c)** kto odpowiada za treść i działanie
> rozszerzenia osoby trzeciej; **(d)** czy Licencjodawca będzie prowadził
> weryfikację rozszerzeń przed udostępnieniem. Do czasu rozstrzygnięcia
> instalowanie w Oprogramowaniu komponentów rozszerzających jest niedozwolone
> na podstawie rozdz. 5.2 pkt 5.

### 7.9 Obowiązki dokumentacyjne przy dystrybucji

Licencjodawca zobowiązuje się dołączać do każdego Wydania Oprogramowania
zestawienie Komponentów wraz z pełnymi tekstami ich licencji oraz notami
o prawach autorskich, w postaci pliku not licencyjnych dostępnego
Licencjobiorcy wraz z Egzemplarzem.

Licencjobiorca, rozpowszechniając Oprogramowanie w przypadkach, w których
rozpowszechnianie jest dopuszczalne na podstawie odrębnej zgody, zobowiązany
jest przekazać ten plik w postaci niezmienionej.

---

## 8. Aktualizacje i wersjonowanie

### 8.1 Zasada wersjonowania

Oprogramowanie podlega zasadzie wersjonowania ustalonej decyzją
architektoniczną: **oznaczenie v2.0 obowiązuje do pierwszej publikacji, w
trakcie budowy
zakazane jest podnoszenie numeru wersji, a status wersji brzmi
„Deweloperski”**. W konsekwencji numer wersji Oprogramowania **nie jest
miernikiem zakresu funkcjonalnego** ani stopnia dojrzałości: dwa Wydania
oznaczone tym samym numerem v2.0 mogą różnić się zakresem działających funkcji.

Licencjobiorca przyjmuje do wiadomości, że identyfikacja Wydania następuje
przez oznaczenie Wydania nadane przez Licencjodawcę, nie przez numer wersji
produktu.

Numer `1.0.0` występujący w metadanych pakietu Powłoki, w manifeście pakietu
Klienta oraz w manifeście pakietu Rust jest odwzorowaniem wersji produktu w
rozumieniu powyższej zasady i podlega tej samej uwadze.

### 8.2 Charakter Aktualizacji

1. Aktualizacja może obejmować poprawki błędów, zmiany zakresu funkcjonalnego,
   zmiany interfejsu, zmiany Kontraktu oraz zmiany schematu trwałości.
2. Ze względu na status Deweloperski Licencjodawca **zastrzega prawo do zmian
   niezachowujących zgodności wstecznej**, w szczególności do zmiany nazw
   komend i zdarzeń Kontraktu, zmiany kluczy konfiguracji oraz zmiany struktury
   tabel.
3. Aktualizacja może usunąć funkcję dostępną w Wydaniu poprzednim, jeżeli
   funkcja ta okazała się niezgodna z przyjętą architekturą albo nie została
   doprowadzona do stanu użytecznego.

### 8.3 Migracje schematu trwałości

1. Schemat Bazy jest tworzony i aktualizowany przez migracje wkompilowane
   w Rdzeń i stosowane samoczynnie przy jego starcie, w jednej transakcji,
   z kontrolą sumy kontrolnej kroku.
2. **Migracje działają wyłącznie w kierunku naprzód.** Oprogramowanie nie
   udostępnia mechanizmu cofnięcia migracji ani przywrócenia schematu do stanu
   poprzedniego. Po uruchomieniu nowszego Wydania na istniejącym Katalogu
   danych powrót do Wydania wcześniejszego może być niemożliwy.
3. Przed instalacją Aktualizacji Licencjobiorca zobowiązany jest wykonać kopię
   zapasową Katalogu danych. Zaniechanie tej czynności obciąża
   Licencjobiorcę; Licencjodawca nie odpowiada za skutki niemożności powrotu do
   Wydania wcześniejszego.

### 8.4 Uprawnienie do Aktualizacji i sposób dostarczania

Oprogramowanie **nie zawiera mechanizmu samoczynnej aktualizacji**: Powłoka nie
korzysta z komponentu aktualizatora, a Rdzeń nie nawiązuje połączeń z
serwerami Licencjodawcy w celu sprawdzenia dostępności nowego Wydania.
Aktualizacja następuje wyłącznie przez świadome działanie Licencjobiorcy —
pobranie i zainstalowanie Wydania udostępnionego przez Licencjodawcę.

> **[DO DECYZJI OPERATORA] — uprawnienie do Aktualizacji i zasady wsparcia.**
> Dokumentacja projektu nie rozstrzyga: **(a)** czy Licencjobiorca ma prawo do
> Aktualizacji w ramach Licencji, czy stanowią one odrębne świadczenie;
> **(b)** czy Licencjodawca zobowiązuje się do wydawania Aktualizacji, a jeżeli
> tak — z jaką częstotliwością i przez jaki okres; **(c)** czy świadczone jest
> wsparcie techniczne, w jakich godzinach, jakim kanałem i z jakimi czasami
> reakcji; **(d)** czy przewidziane jest utrzymanie Wydań wcześniejszych;
> **(e)** czy istnieje zobowiązanie do usuwania podatności bezpieczeństwa
> w określonym terminie.
> Do czasu rozstrzygnięcia obowiązuje rozdz. 8.5.

### 8.5 Brak zobowiązania do rozwoju

Do czasu odmiennego rozstrzygnięcia Licencjodawca **nie zaciąga zobowiązania**
do wydawania Aktualizacji, rozwijania Oprogramowania, doprowadzania do stanu
działającego funkcji przewidzianych koncepcją ani do świadczenia wsparcia
technicznego. Postanowienie to jest bezpośrednią konsekwencją statusu
Deweloperskiego i nieodpłatnego charakteru udostępnienia w tej fazie; ulegnie
zmianie wraz z rozstrzygnięciem kwestii z rozdz. 4.5 i 8.4.

### 8.6 Zmiany treści Licencji

Licencjodawca może wydać nową wersję dokumentu Licencji wraz z nowym Wydaniem
Oprogramowania. Nowa wersja Licencji obowiązuje wobec Wydań udostępnionych po
jej wydaniu. Do Wydania już posiadanego przez Licencjobiorcę stosuje się
wersję Licencji dostarczoną wraz z tym Wydaniem, chyba że strony postanowią
inaczej.

Zmiana Licencji nie może pozbawić Licencjobiorcy uprawnień nabytych w stosunku
do Wydania już posiadanego.

---

## 9. Dane i treści Operatora

### 9.1 Miejsce przechowywania danych

Całość trwałości Oprogramowania mieści się w **jednym pliku bazy SQLite**
o nazwie `danaco-console.db`, umieszczonym w Katalogu danych. Katalogiem danych
jest domyślnie `%LOCALAPPDATA%\DanacoConsole`; wskazanie innego katalogu
następuje zmienną środowiska `DANACO_KATALOG_DANYCH` albo argumentem
wywołania `--dane`. Zmiana Katalogu danych przenosi całość trwałości — nie
istnieje drugie, niezależne miejsce zapisu.

Dane Operatora pozostają na urządzeniu, na którym działa Rdzeń.
**Licencjodawca nie otrzymuje kopii Danych Operatora w związku z korzystaniem
z Oprogramowania.**

### 9.2 Zakres danych zapisywanych przez Oprogramowanie

W Bazie zapisywane są w szczególności:

1. sesje, karty sesji i okna komunikacji wraz z ich ustawieniami (moduł, kanał
   modelu, katalogi robocze, środowisko wykonania, tryb uprawnień, rola).
   **Moment powstania wiersza wymaga zastrzeżenia:** wiersz sesji i wiersz okna
   komunikacji **nie powstają w chwili utworzenia sesji ani okna**
   (`session.create`, `window.create`), lecz dopiero przy zapisie pierwszej
   wiadomości w danym oknie; ponadto zmiana ustawień okna (`window.update`)
   **nie jest utrwalana** — obowiązuje wyłącznie w pamięci Rdzenia do czasu
   jego zatrzymania (rozdz. 10.3 wykaz ograniczeń, pkt 17). Sesja i okno, które
   nie doszły do pierwszej wiadomości, nie zostaną odtworzone po restarcie
   Rdzenia i nie znajdą się w kopii zapasowej Katalogu danych (rozdz. 9.5,
   rozdz. 11.4 pkt 3);
2. **treść wiadomości Operatora oraz treść odpowiedzi modelu** — w tabeli
   wiadomości powiązanej łańcuchem więzów sesja → karta → okno → wiadomość;
3. wartości ustawień konfiguracji na przewidzianych poziomach zasięgu i osiach
   rozstrzygania — z zastrzeżeniem, że zapis wartości nie jest równoznaczny z jej
   stosowaniem przy wywołaniu modelu (rozdz. 10.3 wykaz ograniczeń, pkt 21);
4. definicje kanałów modelu, kont (bez treści Poświadczeń), punktów dostępu
   i nadań, dokumentów tożsamości modelu.

**Czego Baza nie zapisuje.** Poza treścią Poświadczeń (rozdz. 9.3) Baza nie
zapisuje **danych prowenancji wywołań modelu**: struktura prowenancji powstaje
przed uruchomieniem procesu Tury i jest przekazywana do interfejsu, lecz nie ma
dla niej ani tabeli, ani kolumny w schemacie trwałości. Po zamknięciu Tury nie
da się z Bazy odtworzyć, co faktycznie zostało przekazane modelowi — dostępna
pozostaje wyłącznie treść wiadomości i odpowiedzi (rozdz. 2.3, rozdz. 10.3
wykaz ograniczeń, pkt 20). Okoliczność ta ma znaczenie dla rozliczalności pracy
wykonanej z udziałem modelu i dla odtwarzania przebiegu zdarzeń po fakcie.

**Uwaga o treściach obszernych.** Przyjęta decyzja architektoniczna przewiduje
przechowywanie treści obszernych jako plików poza Bazą, z zapisem odwołania w
Bazie.
Mechanizm ten jest przewidziany koncepcją i **w wersji v2.0
niezaimplementowany** — pole odwołania do treści nie jest wypełniane, a treść
zapisywana jest w całości w Bazie. Okoliczność ta ma znaczenie dla planowania
rozmiaru Katalogu danych i dla procedur usuwania danych.

### 9.3 Postępowanie z Poświadczeniami

Zgodnie z przyjętą architekturą Oprogramowanie **nie przechowuje treści
Poświadczeń w Bazie**. Baza przechowuje wyłącznie odwołania — nazwę wpisu
w magazynie sekretów, ścieżkę katalogu konfiguracji Konta albo ścieżkę klucza.
Kontrakt przyjmuje treść poświadczenia wyłącznie w żądaniach dodania
i aktualizacji Konta oraz Punktu dostępu; żadna odpowiedź ani żadne zdarzenie
nie zwraca tej treści — zwracana jest wyłącznie informacja o obecności
poświadczenia albo jego odwołanie.

Odpowiedzialność za zabezpieczenie miejsc, w których faktyczne Poświadczenia
się znajdują (katalogi konfiguracji kont, pliki kluczy, magazyn sekretów
systemu operacyjnego), spoczywa na Licencjobiorcy.

### 9.4 Przekazywanie danych na zewnątrz

Oprogramowanie przekazuje dane poza urządzenie wyłącznie w następstwie
czynności Operatora i wyłącznie w następujących sytuacjach:

1. **wywołanie modelu kanałem rodzaju `cli`** — treść wiadomości, Nakładka
   tożsamości i parametry wywołania są przekazywane uruchamianemu lokalnie
   programowi zewnętrznemu, który we własnym zakresie komunikuje się z usługą
   Dostawcy modelu;
2. **wywołanie modelu kanałem rodzaju `api`** — treść zapytania jest wysyłana
   na adres wskazany wierszem rejestru kanału, przy użyciu poświadczenia
   wskazanego konfiguracją. Droga ta jest w wersji v2.0 **praktycznie
   nieosiągalna**: Oprogramowanie nie udostępnia interfejsu zakładania kanałów
   (rozdz. 3.4 pkt 3), a Konta rodzaju `api` nie mają drogi wykonawczej
   (rozdz. 10.3 wykaz ograniczeń, pkt 19). Jej uruchomienie wymagałoby ręcznego
   zasiania tabeli `kanal_modelu` poza Oprogramowaniem albo wydania komendy
   `channel.add` poza interfejsem aplikacji. Ujawnienie tej drogi ma charakter
   pełnego opisu możliwych wyjść danych, nie zaś opisu funkcji działającej;
3. **most MCP** — konfiguracja mostu przekazana procesowi modelu może
   powodować nawiązanie połączenia z maszyną wskazaną Punktem dostępu, w tym
   połączenia z użyciem programu `ssh` obecnego w systemie operacyjnym.

Poza powyższymi przypadkami Oprogramowanie nie nawiązuje połączeń wychodzących.
**Oprogramowanie nie zawiera telemetrii przekazywanej Licencjodawcy, nie
raportuje zdarzeń użycia, nie sprawdza dostępności aktualizacji i nie wysyła
zgłoszeń o błędach.**

Operator przyjmuje do wiadomości, że treść przekazana modelowi opuszcza
kontrolę Oprogramowania i podlega dalej zasadom przetwarzania stosowanym przez
Dostawcę modelu. Ocena dopuszczalności przekazania określonej kategorii danych
(w tym danych osobowych, tajemnicy przedsiębiorstwa, informacji poufnych osób
trzecich) należy do Licencjobiorcy i jest dokonywana **przed** przekazaniem.

### 9.5 Kopie zapasowe i usuwanie danych

1. Licencjobiorca odpowiada za wykonywanie kopii zapasowych Katalogu danych.
   Oprogramowanie **nie wykonuje kopii zapasowych samoczynnie** i nie zawiera
   mechanizmu przywracania danych z kopii.
   **Zakres kopii wyznacza zakres utrwalenia:** kopia Katalogu danych obejmuje
   wyłącznie to, co zostało zapisane w Bazie. Nie obejmuje zatem sesji i okien,
   które nie doszły do pierwszej wiadomości, zmian ustawień okna dokonanych
   komendą `window.update`, ani danych prowenancji wywołań modelu — żaden
   z tych elementów nie jest utrwalany (rozdz. 9.2 pkt 1, rozdz. 10.3 wykaz
   ograniczeń, pkt 17 i 20).
2. Usunięcie Katalogu danych powoduje nieodwracalną utratę całości trwałości —
   sesji, okien, historii wiadomości, ustawień, definicji kanałów, kont,
   punktów dostępu i nadań.
3. Odinstalowanie Oprogramowania nie usuwa Katalogu danych; jego usunięcie
   wymaga odrębnej czynności Licencjobiorcy.
4. Oprogramowanie nie udostępnia w wersji v2.0 funkcji eksportu całości Danych
   Operatora do formatu przenośnego. Dostęp do danych poza Oprogramowaniem
   możliwy jest przez odczyt pliku Bazy narzędziami zgodnymi z SQLite.

### 9.6 Dane osobowe

Jeżeli Dane Operatora obejmują dane osobowe, administratorem tych danych jest
Licencjobiorca. Licencjodawca nie przetwarza tych danych w związku
z korzystaniem przez Licencjobiorcę z Oprogramowania, ponieważ nie uzyskuje do
nich dostępu (rozdz. 9.1 i 9.4).

Licencjobiorca zobowiązuje się:

1. przetwarzać dane osobowe przy użyciu Oprogramowania zgodnie z przepisami,
   w szczególności z rozporządzeniem (UE) 2016/679;
2. przed przekazaniem danych osobowych modelowi ocenić podstawę prawną tego
   przekazania oraz warunki przetwarzania stosowane przez Dostawcę modelu,
   w tym miejsce przetwarzania i ewentualne przekazanie poza Europejski Obszar
   Gospodarczy;
3. wdrożyć adekwatne środki techniczne i organizacyjne zabezpieczające
   urządzenie, na którym działa Rdzeń, z uwzględnieniem faktu, że
   uwierzytelnianie dostępu do Rdzenia jest w fazie budowy wyłączone
   (rozdz. 5.5 pkt 3).

> **[DO DECYZJI OPERATORA] — role w rozumieniu RODO przy topologii serwerowej.**
> Przy przeniesieniu Rdzenia na serwer `danaco-system` prowadzony przez
> Licencjodawcę albo podmiot z nim powiązany role stron ulegają zmianie:
> Licencjodawca może stać się podmiotem przetwarzającym. Wymaga to
> **(a)** rozstrzygnięcia, kto prowadzi serwer i na czyją rzecz;
> **(b)** zawarcia umowy powierzenia przetwarzania danych osobowych;
> **(c)** określenia lokalizacji przetwarzania, okresu retencji i zasad
> usuwania danych; **(d)** uzupełnienia dokumentu o klauzule odpowiadające
> tej roli. Do czasu rozstrzygnięcia postanowienia rozdz. 9.6 dotyczą wyłącznie
> instalacji z Rdzeniem lokalnym.

### 9.7 Dostęp Licencjodawcy do danych przy czynnościach wsparcia

Jeżeli Licencjobiorca zwróci się o wsparcie techniczne, przekazanie
Licencjodawcy plików diagnostycznych, wycinków Bazy albo zrzutów interfejsu
następuje wyłącznie z inicjatywy Licencjobiorcy i na jego odpowiedzialność.
Licencjobiorca zobowiązany jest usunąć z przekazywanych materiałów dane, których
przekazanie byłoby niedopuszczalne. Licencjodawca wykorzystuje przekazane
materiały wyłącznie w celu udzielenia wsparcia i usuwa je po jego zakończeniu.

---

## 10. Ograniczenie rękojmi i odpowiedzialności

### 10.1 Charakter świadczenia

Oprogramowanie jest udostępniane **w stanie, w jakim się znajduje**
(„as is”), z zastrzeżeniem statusu Deweloperskiego wersji v2.0. Licencjodawca
udostępnia produkt w toku budowy, o niedomkniętym zakresie funkcjonalnym,
którego rzeczywisty stan wykonania został ustalony audytem i ujawniony
w niniejszym dokumencie.

### 10.2 Wyłączenie rękojmi

W najszerszym zakresie dopuszczalnym przez prawo właściwe Licencjodawca
wyłącza rękojmię za wady Oprogramowania oraz nie udziela żadnych gwarancji,
w szczególności nie zapewnia, że:

1. Oprogramowanie będzie działać nieprzerwanie, bez błędów i bez przestojów;
2. Oprogramowanie będzie przydatne do określonego celu zamierzonego przez
   Licencjobiorcę;
3. wyniki pracy modelu uzyskane przy użyciu Oprogramowania będą poprawne,
   kompletne, aktualne albo nadające się do wykorzystania bez weryfikacji;
4. dane zapisane w Bazie nie ulegną utracie ani uszkodzeniu;
5. funkcje przewidziane koncepcją produktu zostaną doprowadzone do stanu
   działającego.

Wyłączenie rękojmi nie dotyczy przypadków, w których przepis bezwzględnie
obowiązujący nie dopuszcza takiego wyłączenia, ani szkody wyrządzonej umyślnie.

### 10.3 Ujawnienie rzeczywistego stanu wykonania

Licencjodawca ujawnia poniżej stan wykonania wersji v2.0 w zakresie istotnym
dla oceny przydatności Oprogramowania. Ujawnienie ma charakter oświadczenia
o stanie przedmiotu świadczenia i wyłącza możliwość powołania się przez
Licencjobiorcę na nieświadomość tych ograniczeń.

**Funkcje potwierdzone jako działające:**

1. transport komunikatów po WebSocket wraz z kopertą Kontraktu, odpornością na
   nieznaną komendę i przetrwaniem pracy Rdzenia mimo rozłączenia Klienta;
2. obsługa wszystkich komend Kontraktu po stronie Rdzenia;
3. trwałość w pliku SQLite wraz z samoczynnym stosowaniem migracji
   i odtwarzaniem sesji oraz okien po restarcie Rdzenia;
4. kanał główny modelu (uruchomienie programu zewnętrznego, katalog
   konfiguracji per Konto, odczyt strumienia `stream-json`, rozpoznanie
   wyczerpania limitu, rotacja Kont w obrębie puli, nadanie prowenancji przed
   startem procesu);
5. rejestr kanałów modelu sterowany danymi wraz z odświeżeniem po zmianie;
6. okno konfiguracji ze sterowanym danymi katalogiem kategorii, definicji
   i opcji oraz zapisem wartości na poziomach zasięgu i osiach — **z istotnym
   zastrzeżeniem**: działa katalog, rozstrzyganie i zapis wartości, natomiast
   **większość zapisanych kluczy nie jest odczytywana przy wywołaniu modelu**
   (wykaz ograniczeń, pkt 21). Zapis wartości nie jest równoznaczny z jej
   stosowaniem;
7. zarządzanie kontami, punktami dostępu, nadaniami i tożsamością modelu
   z poziomu interfejsu — **z zastrzeżeniem**, że nie da się zapisać Punktu
   dostępu rodzaju `localDirectory` (wykaz ograniczeń, pkt 22); funkcja działa
   w zakresie kont, tożsamości modelu, punktów dostępu rodzaju `mcpBridge`
   i nadań;
8. Nakładka tożsamości w trybach `ZASTAP` i `DOLACZ` przekazywana procesowi
   modelu;
9. zapis historii rozmowy do Bazy;
10. warstwa wizualna: dwa motywy, zwarta gęstość, zestaw ikon;
11. zatrzymanie bieżącej tury.

**Ograniczenia istotne — funkcje przewidziane koncepcją, w wersji v2.0
niedziałające albo działające częściowo:**

1. **Świeża instalacja nie wykona wywołania modelu bez czynności
   przygotowawczych** — brak zaczynu rejestru kanałów, pusta pula kont odmawia
   wywołania, a Klient przekazuje rodzaj kanału zamiast identyfikatora wiersza
   rejestru (rozdz. 3.4 pkt 3).
2. **Brak ciągłości rozmowy** — każda tura uruchamia nowy proces bez historii;
   model nie widzi poprzednich wiadomości Okna komunikacji.
3. **Historia rozmowy nie jest odtwarzana w interfejsie** — mimo zapisu do Bazy
   Klient nie pobiera zapisanych wiadomości; po odświeżeniu widoku historia
   znika z ekranu.
4. **Okna operacyjne modułów nie istnieją** — z siedemdziesięciu dziewięciu wierszy
   katalogu okien operacyjnych w Bazie, będących liczbą obowiązującą w całej
   dokumentacji produktu ([INSTRUKCJA-UZYTKOWANIA.md](INSTRUKCJA-UZYTKOWANIA.md)
   rozdz. 7), zbudowany jest **jeden widok** — okno komunikacji, któremu odpowiada
   piętnaście wierszy katalogu, po jednym na moduł; wybór modułu
   w nawigacji nie przeładowuje przestrzeni roboczej. Wszystkie piętnaście
   modułów (Studio, Workspace, Browser, Research, Library, Translate,
   Roundtable, Design, Assistant, Terminal, Developer, Diagnostics, Apps,
   Agents, Automations) występuje wyłącznie jako pozycje nawigacji.
5. **Środowiska TalkIn, WorkSpace, CodeStudio i MultitaskingAI** mają
   nawigację i opis, lecz nie mają realizacji funkcjonalnej; środowisko
   MultitaskingAI nie posiada panelu orkiestracji, ról, kolejek ról ani
   monitora.
6. **Pulpit dowodzenia (Mission Control)** — dane przykładowe wbudowane
   w interfejs (plik `dane-przykladowe.ts`) zostały usunięte; pulpit rysuje się
   obecnie na podstawie stanu Rdzenia, odczytywanego komendami `session.list`,
   `window.list` i `channel.list` oraz podtrzymywanego subskrypcją zdarzeń
   `session.changed`, `window.changed`, `queue.changed` i `progress.changed`,
   z uczciwymi stanami pustymi tam, gdzie Rdzeń nie zwrócił jeszcze danych.
   Zmiana ta nie obejmuje zamiarów: zamiary wyrażone na pulpicie (kolejka,
   utworzenie zasobu, decyzja) nadal **nie wykonują komend Kontraktu** —
   zapisują się wyłącznie do lokalnego paska działań pulpitu, bez wywołania
   kanału.
7. **Kolejki** — komendy kolejek zmieniają wyłącznie stan wiersza; pozycje
   kolejki nie powstają, silnik wykonania nie istnieje.
8. **Przekazanie zlecenia koordynator → wykonawca w interfejsie** jest
   wyłącznie animacją lokalną, bez wywołania komendy Kontraktu. Sama pętla
   koordynator–wykonawca istnieje w Rdzeniu i uruchamia się wraz z turą.
9. **Wybór modelu, modelu zapasowego i nakładu rozumowania na poziomie Okna
   komunikacji** — wartości zapisują się pod kluczami nierozpoznawanymi przez
   Rdzeń i nie wpływają na wywołanie.
10. **Zakres wykonania i host wykonania** — wyłącznie etykieta w prowenancji;
    proces modelu startuje zawsze na maszynie Rdzenia (rozdz. 4.7 pkt 3).
11. **Izolacja** — jedenaście ustawień rozstrzyganych na ośmiu poziomach
    zasięgu, bez egzekutorów (rozdz. 5.4 pkt 4).
12. **Sterowanie platformą przez narzędzia modelu** — trzydzieści dziewięć
    deklaracji narzędzi jest generowanych z Kontraktu, lecz nie ma konsumenta;
    model nie steruje platformą.
13. **Nawigacja platformy przez Kontrakt** — komendy strony głównej, środowisk,
    modułów i przestrzeni roboczej są zaimplementowane po obu stronach, lecz
    w większości nie są wywoływane przez żaden widok. **Wyjątek**: `session.list`,
    `session.bind` i `session.focus` są wywoływane przez stronę główną w strefie
    „Sesje w tle” — odczyt sesji trwających na Rdzeniu, powrót do sesji przez
    powiązanie połączenia (`session.bind`) i przeniesienie ogniska karty czynnej
    (`session.focus`), dostępne wyłącznie po ustaleniu tożsamości klienta
    z powitania połączenia.
14. **Pamięć wielopoziomowa** — istnieje warstwa danych, brak komend Kontraktu
    i konsumentów.
15. **Rozszerzenia i katalog rozszerzeń** — nie występują w Kontrakcie
    wykonawczym (rozdz. 7.8).
16. **Telemetria postępu w interfejsie** — zdarzenie `progress.changed` **ma już
    konsumenta**: źródło danych pulpitu dowodzenia subskrybuje je wraz
    z `session.changed`, `window.changed` i `queue.changed` i wprowadza jego
    treść do składania kompletu danych pulpitu. Zapytanie o stan okna po
    stronie Rdzenia nadal zwraca zawsze stan oczekujący, dopóki dla tego okna
    nie istnieje proces Tury — ograniczenie to dotyczy odczytu stanu okna
    niezależnie od tego, czy zdarzenie postępu ma subskrybenta po stronie
    Klienta.
17. **Wiersz sesji i wiersz okna powstają dopiero przy pierwszej wiadomości**,
    nie w chwili `session.create` ani `window.create`; **zmiana ustawień okna**
    (`window.update`) nie jest utrwalana w Bazie, lecz obowiązuje wyłącznie
    w pamięci Rdzenia do jego zatrzymania; sesje trafiają zawsze do pierwszego
    środowiska. Skutki dla zakresu utrwalenia opisuje rozdz. 9.2 pkt 1, dla
    kopii zapasowych — rozdz. 9.5, dla stanu po zakończeniu Licencji —
    rozdz. 11.4 pkt 3.
18. **Instalator NSIS Powłoki nie powstaje.** Na rewizji audytu pierwotnego
    stanu sprzed prac naprawczych konfiguracja Powłoki (`tauri.conf.json`) nie zawierała kroku
    poprzedzającego budowanie (`beforeBuildCommand`) ani adresu trybu
    rozwojowego (`devUrl`), a repozytorium nie zawierało skryptu składającego
    katalog produktu — złożenie katalogu `C:\DanacoConsole_App` było wówczas
    czynnością ręczną. **Stan ten uległ zmianie w pracach nad ścieżką wydania**: `tauri.conf.json` zawiera obecnie `beforeBuildCommand`
    i `devUrl`; repozytorium zawiera `budowa/scripts/wydanie.sh`
    i `budowa/scripts/pakowanie.sh`, którymi jednym poleceniem składa się
    produkt (rdzeń `danaco-console.exe`, powłoka `Danaco Console.exe`,
    `client/dist`) do katalogu `C:\DanacoConsole_App` — z tej właśnie ścieżki
    zbudowano i uruchomiono Egzemplarz, do którego dołączona jest niniejsza
    Licencja. Punkt niniejszy, dawniej ujawniający brak jakiejkolwiek
    automatyzacji budowania, pozostaje aktualny w zakresie węższym:
    **wygenerowanie instalatora Windows w formacie NSIS nadal nie następuje** —
    krok `cargo tauri build` wytwarzający instalator wymaga obecności narzędzia
    NSIS na maszynie budującej; ścieżka wydania pomija budowę Powłoki
    automatycznie przy braku `cargo` lub `tauri-cli`, a sam instalator pozostaje
    krokiem odrębnym, udokumentowanym w skrypcie wydania jako wymagający NSIS,
    którego na maszynie audytowanej brak. Punkt niniejszy podlega ponownej
    weryfikacji przy każdym Wydaniu (Załącznik C.4).
19. **Zmiany definicji Kont wymagają ponownego uruchomienia Rdzenia** — pula
    rotacji budowana jest jednorazowo przy starcie; stan wyczerpania limitu nie
    jest utrwalany. Konta rodzaju `api` i `sdk` nie mają drogi wykonawczej
    (skutki dla kanału rodzaju `api` opisuje rozdz. 9.4 pkt 2).
20. **Prowenancja nie jest utrwalana.** Struktura prowenancji powstaje przed
    uruchomieniem procesu Tury i jest przekazywana do interfejsu, lecz w całym
    schemacie trwałości nie ma dla niej tabeli ani kolumny — w szczególności
    nie ma jej tabela wiadomości, której polami treściowymi są: rola, persona,
    okno źródłowe, rodzaj treści, stan, treść, odwołanie do treści, liczniki
    tokenów wejścia i wyjścia, kolejność i czas utworzenia — obok klucza
    głównego `id` i klucza obcego `okno_komunikacji_id` wiążącego wiadomość
    z Oknem komunikacji (`budowa/server/internal/store/migracja_002_okna.sql`).
    Wyliczenie to obejmuje pola treściowe tabeli, nie stanowi pełnego wykazu
    jej schematu. **Po zamknięciu Tury nie da
    się odtworzyć, co faktycznie poszło do modelu** — jakie argumenty wywołania,
    jaka treść Nakładki tożsamości i jakie ustawienia zostały zastosowane.
    Ograniczenie to jest istotne wszędzie tam, gdzie Licencjobiorca opiera się
    na rozliczalności pracy wykonanej z udziałem modelu; wyklucza ono również
    posługiwanie się prowenancją jako dowodem po fakcie (rozdz. 2.3,
    rozdz. 9.2).
21. **Większość kluczy katalogu ustawień zapisuje się, lecz nie jest
    odczytywana przy wywołaniu modelu.** W wykonaniu Tury stosowane są wyłącznie:
    Nakładka tożsamości (tryb `ZASTAP` albo `DOLACZ`), tryb uprawnień, katalogi
    robocze, katalog roboczy sesji oraz mosty MCP. Pozostałe klucze katalogu
    ustawień — w tym **cała rodzina `harness.*`** (`harness.program_claude`,
    `harness.plik_ustawien`, `harness.konfiguracja_mcp`) — zapisują się w Bazie
    **bez wpływu na wywołanie modelu**. Ujawnione wcześniej przypadki
    szczegółowe (rozdz. 3.4 pkt 1 — program kanału głównego; pkt 9 powyżej —
    model, model zapasowy i nakład rozumowania; pkt 10 — zakres wykonania;
    pkt 11 — izolacja) są przykładami tej reguły ogólnej, nie wyjątkami od niej.
    W konsekwencji uprawnienie do konfigurowania Oprogramowania przyznane
    w rozdz. 4.2 pkt 3 obejmuje **zapis i rozstrzyganie wartości**, nie zaś
    zapewnienie, że każda zapisana wartość wpłynie na zachowanie modelu.
22. **Punktu dostępu rodzaju `localDirectory` nie da się zapisać.** Schemat
    trwałości wymaga dla tego rodzaju powiązania z wierszem urządzenia (kolumna
    `urzadzenie_id` wraz z więzem sprawdzającym `rodzaj <> 'localDirectory' OR
    urzadzenie_id IS NOT NULL`), a wiersz w tabeli `urzadzenie` nie powstaje
    w żaden sposób produkcyjny: nie istnieje repozytorium urządzeń w warstwie
    danych Rdzenia, migracje nie niosą zaczynu tej tabeli, a jedyne jej
    zasilenia występują w plikach sprawdzianów. W wersji v2.0 zarządzanie
    punktami dostępu z poziomu interfejsu obejmuje zatem wyłącznie punkty
    rodzaju `mcpBridge`.

Wykaz powyższy odzwierciedla stan w bieżącym stanie repozytorium — audyt pierwotny
przeprowadzono w stanie sprzed prac naprawczych, a punkty 6, 16 i 18 zaktualizowano po
weryfikacji zmian wprowadzonych pracami naprawczymi, poprzedzającymi
bieżący stan repozytorium (Załącznik C.1, C.3 pkt 4). Punkty nieoznaczone jako
zaktualizowane pozostają prawdziwe na obu rewizjach. Wykaz nie stanowi
zobowiązania do usunięcia wymienionych ograniczeń.

### 10.4 Ograniczenie odpowiedzialności

W najszerszym zakresie dopuszczalnym przez prawo właściwe Licencjodawca nie
ponosi odpowiedzialności za:

1. utracone korzyści, utratę przychodu, utratę spodziewanych oszczędności,
   utratę renomy ani szkody pośrednie;
2. utratę, uszkodzenie albo ujawnienie Danych Operatora oraz Treści Operatora;
3. skutki decyzji podjętych na podstawie wyników pracy modelu;
4. działania modelu w zakresie zasobów udostępnionych mu przez Operatora,
   w tym za zmiany dokonane w plikach na maszynach objętych nadaniem trybu
   `write`;
5. przerwy w działaniu albo zmianę warunków świadczenia usług przez Dostawcę
   modelu, w tym za wyczerpanie limitów, zawieszenie konta albo zaprzestanie
   udostępniania programu zewnętrznego;
6. niezgodność korzystania z Oprogramowania z warunkami Dostawcy modelu albo
   z przepisami obowiązującymi Licencjobiorcę;
7. szkody wynikłe z korzystania z Oprogramowania w sposób sprzeczny
   z rozdz. 5.3.

Ograniczenia nie dotyczą szkody wyrządzonej umyślnie ani odpowiedzialności,
której wyłączyć nie można na podstawie przepisu bezwzględnie obowiązującego.

> **[DO DECYZJI OPERATORA] — górna granica odpowiedzialności.**
> Dokumentacja projektu nie rozstrzyga, czy odpowiedzialność Licencjodawcy ma
> być ograniczona kwotowo. Warianty: **(a)** brak limitu kwotowego, wyłącznie
> wyłączenia rodzajowe z rozdz. 10.4; **(b)** limit równy wynagrodzeniu
> zapłaconemu przez Licencjobiorcę w okresie dwunastu miesięcy poprzedzających
> zdarzenie; **(c)** limit kwotowy oznaczony wprost; **(d)** przy licencji
> nieodpłatnej — ograniczenie do przypadków winy umyślnej.
> Wybór jest ściśle powiązany z rozstrzygnięciem kwestii odpłatności
> (rozdz. 4.5) i wymaga oceny skuteczności w świetle prawa właściwego.

### 10.5 Rozkład ryzyka

Strony przyjmują, że opisany rozkład odpowiedzialności odpowiada charakterowi
świadczenia: Oprogramowanie o statusie Deweloperskim udostępniane jest
z ujawnieniem rzeczywistych ograniczeń (rozdz. 10.3), a Licencjobiorca
podejmuje decyzję o korzystaniu z pełną wiedzą o tych ograniczeniach.

### 10.6 Siła wyższa

Żadna ze stron nie odpowiada za niewykonanie zobowiązań spowodowane
okolicznościami pozostającymi poza jej rozsądną kontrolą, w szczególności
działaniem siły wyższej, awarią infrastruktury teleinformatycznej niezależną od
strony, decyzją organu władzy publicznej ani zaprzestaniem świadczenia usług
przez Dostawcę modelu.

---

## 11. Okres obowiązywania i zakończenie licencji

### 11.1 Wejście w życie

Licencja wiąże z chwilą pierwszego z następujących zdarzeń: zainstalowania
Oprogramowania, pierwszego uruchomienia Oprogramowania albo zaakceptowania
warunków Licencji w toku instalacji. Przystąpienie do korzystania
z Oprogramowania oznacza przyjęcie warunków Licencji w całości.

Podmiot, który nie akceptuje warunków Licencji, zobowiązany jest odstąpić od
instalacji, a Egzemplarz — usunąć.

### 11.2 Okres obowiązywania

> **[DO DECYZJI OPERATORA] — okres obowiązywania Licencji.**
> Dokumentacja projektu nie rozstrzyga okresu, na jaki Licencja jest
> udzielana. Warianty: **(a)** czas nieoznaczony, z prawem wypowiedzenia przez
> każdą ze stron z zachowaniem terminu; **(b)** czas oznaczony odpowiadający
> okresowi abonamentu, z automatycznym przedłużeniem albo bez niego;
> **(c)** czas oznaczony powiązany ze statusem wersji — do dnia wydania
> pierwszego Wydania o statusie innym niż Deweloperski; **(d)** licencja
> bezterminowa na Wydanie posiadane, z odrębnym uprawnieniem do Aktualizacji.
> Wybór jest powiązany z rozstrzygnięciem odpłatności (rozdz. 4.5)
> i uprawnienia do Aktualizacji (rozdz. 8.4).

Niezależnie od wybranego wariantu Licencja wygasa najpóźniej z chwilą
zaprzestania korzystania z Oprogramowania przez Licencjobiorcę połączonego
z usunięciem wszystkich Egzemplarzy.

### 11.3 Wypowiedzenie i rozwiązanie

1. **Wypowiedzenie przez Licencjobiorcę.** Licencjobiorca może w każdym czasie
   zakończyć Licencję przez zaprzestanie korzystania z Oprogramowania
   i usunięcie wszystkich Egzemplarzy oraz kopii zapasowych Oprogramowania.
   Zakończenie nie rodzi po stronie Licencjodawcy obowiązku zwrotu
   świadczeń, chyba że co innego wynika z rozstrzygnięcia kwestii odpłatności.
2. **Wypowiedzenie przez Licencjodawcę z przyczyny naruszenia.** Licencjodawca
   może wypowiedzieć Licencję ze skutkiem natychmiastowym, jeżeli
   Licencjobiorca istotnie narusza jej postanowienia, w szczególności
   ograniczenia z rozdz. 5.2, 5.3 albo 5.6, i nie zaprzestaje naruszenia
   w terminie wyznaczonym w wezwaniu skierowanym na piśmie albo drogą
   elektroniczną — **[DO DECYZJI OPERATORA]** co do długości tego terminu.

   > **[DO DECYZJI OPERATORA] — termin usunięcia naruszenia.**
   > Dokumentacja projektu nie rozstrzyga, jaki termin należy wyznaczyć
   > Licencjobiorcy na zaprzestanie naruszenia przed wypowiedzeniem.
   > Warianty: **(a)** termin czternastu dni — jednolity dla wszystkich
   > naruszeń; **(b)** termin trzydziestu dni — dający czas na usunięcie
   > naruszeń organizacyjnych (przekroczenie liczby stanowisk, brak ewidencji);
   > **(c)** termin różnicowany: naruszenia usuwalne — trzydzieści dni,
   > naruszenia nieusuwalne (rozpowszechnienie Oprogramowania, obejście
   > oznaczeń) — wypowiedzenie bez wyznaczania terminu; **(d)** termin
   > „odpowiedni” wyznaczany każdorazowo w wezwaniu, nie krótszy niż oznaczony
   > minimalny. Wybór jest powiązany z rozstrzygnięciem kwestii kar umownych
   > (rozdz. 5.8) i zakresu podmiotowego (rozdz. 4.3).
3. **Wypowiedzenie bez przyczyny.** Uprawnienie Licencjodawcy do wypowiedzenia
   Licencji bez wskazania przyczyny oraz termin takiego wypowiedzenia zależą od
   rozstrzygnięcia kwestii okresu obowiązywania (rozdz. 11.2) i odpłatności
   (rozdz. 4.5).
4. **Rozwiązanie za porozumieniem.** Strony mogą zakończyć Licencję
   w każdym czasie za obopólnym porozumieniem wyrażonym na piśmie.

### 11.4 Skutki zakończenia

Z chwilą zakończenia Licencji:

1. wygasają wszystkie uprawnienia Licencjobiorcy opisane w rozdz. 4.2;
2. Licencjobiorca zobowiązany jest zaprzestać korzystania z Oprogramowania,
   odinstalować je ze wszystkich urządzeń i usunąć wszystkie Egzemplarze oraz
   kopie zapasowe Oprogramowania;
3. Licencjobiorca zachowuje prawo do Danych Operatora i Treści Operatora oraz
   prawo do zachowania kopii Katalogu danych — zakończenie Licencji **nie
   zobowiązuje do usunięcia własnych danych**. Zakres danych możliwych do
   zachowania wyznacza jednak zakres ich faktycznego utrwalenia w Bazie
   (rozdz. 9.2 pkt 1 oraz rozdz. 9.5 pkt 1); Licencjobiorca zamierzający
   zabezpieczyć stan pracy przed zakończeniem Licencji powinien wziąć to pod
   uwagę, tym bardziej że Oprogramowanie nie udostępnia funkcji eksportu
   całości Danych Operatora do formatu przenośnego (rozdz. 9.5 pkt 4);
4. Licencjobiorca zachowuje prawo korzystania z wytworów pracy powstałych
   przy użyciu Oprogramowania w okresie obowiązywania Licencji;
5. na żądanie Licencjodawcy Licencjobiorca potwierdza wykonanie obowiązków
   z pkt 2 oświadczeniem złożonym na piśmie albo drogą elektroniczną.

### 11.5 Postanowienia trwające

Zakończenie Licencji nie wpływa na moc obowiązującą postanowień, które
z natury mają obowiązywać dłużej, w szczególności rozdz. 6 (własność
intelektualna), rozdz. 10 (odpowiedzialność), rozdz. 13 (poufność) oraz
rozdz. 14 (postanowienia końcowe).

---

## 12. Zgodność licencyjna i kontrola korzystania

### 12.1 Obowiązek zgodności

Licencjobiorca korzysta z Oprogramowania w granicach zakresu podmiotowego
Licencji (rozdz. 4.3) i utrzymuje ewidencję pozwalającą wykazać zgodność
korzystania z tym zakresem. **Obowiązek niniejszy dzieli los rozstrzygnięcia
oznaczonego [DO DECYZJI OPERATORA] w rozdz. 4.3 (Załącznik B poz. 3)**: dopóki
nie zostanie wybrany wariant liczby stanowisk i jednostki liczenia licencji,
ewidencja, o której mowa wyżej, nie ma mierzalnego przedmiotu — nie da się
stwierdzić, jaki stan faktyczny ewidencja ma dokumentować jako zgodny.
Obowiązek ewidencyjny z rozdz. 4.3 pkt 1 i obowiązek niniejszy są jednym
i tym samym obowiązkiem opisanym w dwóch miejscach dokumentu ze względu na
układ rozdziałów; rozstrzygnięcie rozdz. 4.3 czyni zadość obu.

### 12.2 Weryfikacja

> **[DO DECYZJI OPERATORA] — prawo kontroli zgodności licencyjnej.**
> Dokumentacja projektu nie rozstrzyga, czy Licencjodawcy przysługuje
> uprawnienie do weryfikacji sposobu korzystania z Oprogramowania.
> Warianty: **(a)** brak uprawnienia kontrolnego — wyłącznie oświadczenie
> Licencjobiorcy na żądanie; **(b)** prawo żądania oświadczenia o liczbie
> Instalacji i Operatorów, składanego nie częściej niż raz w roku;
> **(c)** prawo audytu w siedzibie Licencjobiorcy, po uprzedzeniu,
> w godzinach pracy, na koszt Licencjodawcy z zastrzeżeniem obciążenia
> Licencjobiorcy kosztami przy stwierdzeniu naruszenia przekraczającego
> ustalony próg.
> Wariant (c) wymaga uregulowania zasad zachowania poufności danych
> Licencjobiorcy w toku kontroli oraz wyłączenia dostępu do treści
> zapisanych w Bazie.

**Oprogramowanie nie zawiera mechanizmu technicznej kontroli licencji**: nie
weryfikuje klucza licencyjnego, nie łączy się z serwerem aktywacji ani nie
raportuje faktu uruchomienia. Weryfikacja zgodności może zatem odbywać się
wyłącznie środkami organizacyjnymi.

### 12.3 Skutki stwierdzenia niezgodności

W razie stwierdzenia korzystania wykraczającego poza zakres Licencji strony
w pierwszej kolejności ustalają zakres niezgodności i sposób jej usunięcia.
Nieusunięcie niezgodności w uzgodnionym terminie stanowi istotne naruszenie
Licencji w rozumieniu rozdz. 11.3 pkt 2.

---

## 13. Poufność

### 13.1 Informacje poufne

Za informacje poufne uważa się nieujawnione publicznie informacje dotyczące
Oprogramowania, przekazane Licencjobiorcy przez Licencjodawcę albo powzięte
w związku z korzystaniem z Oprogramowania, w szczególności: dokumentację
techniczną nieprzeznaczoną do publikacji, opisy architektury wewnętrznej,
struktury Kontraktu i modelu danych, informacje o niedziałających funkcjach
i podatnościach, a także informacje o planach rozwoju produktu.

Za informacje poufne **nie uważa się** informacji, które: są publicznie
dostępne bez naruszenia zobowiązania do poufności; były znane Licencjobiorcy
przed ich przekazaniem; zostały uzyskane zgodnie z prawem od osoby trzeciej
nieobjętej zobowiązaniem do poufności; albo zostały opracowane samodzielnie
bez wykorzystania informacji poufnych.

### 13.2 Obowiązki

Licencjobiorca zobowiązuje się zachować informacje poufne w tajemnicy, nie
ujawniać ich osobom trzecim bez zgody Licencjodawcy udzielonej na piśmie oraz
udostępniać je własnym Operatorom i współpracownikom wyłącznie w zakresie
niezbędnym do korzystania z Oprogramowania, po zobowiązaniu ich do poufności.

Obowiązek zachowania poufności nie stoi na przeszkodzie ujawnieniu informacji na
żądanie uprawnionego organu, w zakresie i w trybie wynikającym z przepisów;
o takim żądaniu Licencjobiorca niezwłocznie zawiadamia Licencjodawcę, o ile
nie jest to zabronione.

### 13.3 Okres obowiązywania

Zobowiązanie do poufności obowiązuje w okresie trwania Licencji oraz przez
oznaczony okres po jej zakończeniu — **[DO DECYZJI OPERATORA]** co do długości
tego okresu — a w odniesieniu do informacji stanowiących tajemnicę
przedsiębiorstwa: przez cały okres, w którym zachowują one taki charakter.
Ta druga część zobowiązania obowiązuje niezależnie od wybranego wariantu
i nie jest ograniczona terminem.

> **[DO DECYZJI OPERATORA] — okres poufności po zakończeniu Licencji.**
> Dokumentacja projektu nie rozstrzyga, jak długo po zakończeniu Licencji trwa
> zobowiązanie do poufności co do informacji niebędących tajemnicą
> przedsiębiorstwa. Warianty: **(a)** brak okresu następczego — zobowiązanie
> wygasa wraz z Licencją, z wyjątkiem tajemnicy przedsiębiorstwa;
> **(b)** trzy lata; **(c)** pięć lat; **(d)** okres bezterminowy dla całości
> informacji poufnych. Wybór ma znaczenie praktyczne ze względu na zakres
> informacji objętych rozdz. 13.1 — należą do niego informacje o funkcjach
> niedziałających i o podatnościach, których wrażliwość maleje wraz z wydaniem
> kolejnych Wydań; wariant (d) bywa kwestionowany jako nadmierny.

### 13.4 Ujawnienie podatności

Licencjobiorca, który poweźmie informację o podatności bezpieczeństwa
Oprogramowania, zawiadamia o niej wyłącznie Licencjodawcę i powstrzymuje się
od jej publikacji przez okres uzgodniony z Licencjodawcą, nie dłuższy jednak niż
konieczny do wydania poprawki.

> **[DO DECYZJI OPERATORA] — kanał i polityka zgłaszania podatności.**
> Dokumentacja projektu nie wskazuje adresu ani trybu zgłaszania podatności.
> Warianty: **(a)** wskazanie adresu poczty elektronicznej do zgłoszeń
> bezpieczeństwa; **(b)** przyjęcie polityki skoordynowanego ujawniania
> z określonym terminem; **(c)** brak odrębnego kanału — zgłoszenia zwykłą
> drogą wsparcia. Rozstrzygnięcie wymaga wpisania danych kontaktowych do
> rozdz. 14.7.

---

## 14. Postanowienia końcowe

### 14.1 Prawo właściwe

> **[DO DECYZJI OPERATORA] — prawo właściwe.**
> Dokumentacja projektu nie rozstrzyga prawa właściwego dla Licencji.
> Warianty: **(a)** prawo polskie — naturalne ze względu na siedzibę
> Licencjodawcy, z wyłączeniem norm kolizyjnych; **(b)** prawo państwa
> siedziby Licencjobiorcy — przy dystrybucji międzynarodowej;
> **(c)** prawo polskie z zastrzeżeniem bezwzględnie obowiązujących przepisów
> ochrony konsumenta państwa zwykłego pobytu Licencjobiorcy, jeżeli
> Licencjobiorcą miałby być konsument.
> Wybór przesądza skuteczność rozdz. 5.7, 10.2 i 10.4 oraz treść rozdz. 14.2.

Do czasu rozstrzygnięcia przyjmuje się, że wykładni Licencji dokonuje się
z uwzględnieniem prawa polskiego jako prawa siedziby Licencjodawcy, bez
przesądzania właściwości.

### 14.2 Rozstrzyganie sporów

> **[DO DECYZJI OPERATORA] — właściwość sądu i tryb rozstrzygania sporów.**
> Warianty: **(a)** sąd powszechny właściwy miejscowo dla siedziby
> Licencjodawcy; **(b)** sąd właściwy według przepisów ogólnych;
> **(c)** sąd polubowny (arbitraż) ze wskazaniem instytucji i regulaminu;
> **(d)** obowiązkowa próba mediacji przed skierowaniem sprawy na drogę sądową.
> Przy Licencjobiorcy będącym konsumentem klauzula wyłącznej właściwości
> miejscowej może być nieskuteczna — wymaga oceny prawnej.

Niezależnie od wybranego wariantu strony podejmą w dobrej wierze próbę
polubownego rozwiązania sporu przed skierowaniem sprawy na drogę
postępowania sądowego.

### 14.3 Forma zmian

Zmiany Licencji w stosunku dwustronnym wymagają formy pisemnej albo
dokumentowej pod rygorem nieważności, z zastrzeżeniem rozdz. 8.6, który dotyczy
wydawania nowych wersji dokumentu Licencji wraz z nowymi Wydaniami.

### 14.4 Klauzula salwatoryjna

Nieważność albo bezskuteczność któregokolwiek z postanowień Licencji nie
wpływa na ważność pozostałych postanowień. W miejsce postanowienia nieważnego
albo bezskutecznego strony stosują postanowienie skuteczne, najbliższe celowi
gospodarczemu postanowienia zastępowanego. Dotyczy to w szczególności
postanowień rozdz. 5.2, 10.2 i 10.4, których zakres podlega ograniczeniu do
granic dopuszczalnych przez prawo właściwe.

### 14.5 Całość porozumienia

Licencja wraz z załącznikami stanowi całość porozumienia stron w zakresie
korzystania z Oprogramowania i zastępuje wcześniejsze ustalenia, oświadczenia
i zapewnienia dotyczące tego przedmiotu, z zastrzeżeniem rozdz. 1.4 pkt 2.
Materiały prezentacyjne, koncepcyjne i poglądowe dotyczące produktu nie
stanowią części Licencji i nie tworzą zobowiązania co do zakresu
funkcjonalnego.

### 14.6 Język dokumentu

Językiem Licencji jest język polski. W razie sporządzenia tłumaczenia
rozstrzygające znaczenie ma wersja polska.

### 14.7 Dane kontaktowe i doręczenia

Oświadczenia i zawiadomienia związane z Licencją składa się na adres siedziby
Licencjodawcy albo na wskazany przez niego adres poczty elektronicznej.

> **[DO DECYZJI OPERATORA] — dane identyfikacyjne i kontaktowe Licencjodawcy.**
> Dokumentacja projektu wskazuje wyłącznie firmę Licencjodawcy
> („Danaco Holding Group Sp. z o.o.”) oraz Twórcę. Do publikacji dokumentu
> konieczne jest uzupełnienie: adresu siedziby, numeru KRS, numeru NIP,
> adresu poczty elektronicznej do korespondencji oraz — jeżeli zostanie
> przyjęty odrębny kanał — adresu do zgłaszania podatności (rozdz. 13.4).

### 14.8 Nagłówki i załączniki

Załączniki A, B i C stanowią integralną część Licencji. W razie rozbieżności
między treścią rozdziału a treścią załącznika rozstrzyga treść rozdziału,
z wyjątkiem Załącznika A w zakresie zestawienia Komponentów, gdzie rozstrzyga
treść załącznika jako zestawienia szczegółowego.

---

## Załącznik A — zestawienie komponentów osób trzecich

Zestawienie sporządzono na podstawie bieżącego stanu repozytorium
(audyt pierwotny w stanie sprzed prac naprawczych, zaktualizowany po zamknięciu prac naprawczych —
Załącznik C.3 pkt 4). Zestawienie **nie zastępuje** pełnych tekstów licencji, które Licencjodawca
dołącza do Wydania zgodnie z rozdz. 7.9.

### A.1. Warstwa Rdzenia — moduł Go `danacoconsole` (wymagane `go 1.26`)

| Lp. | Komponent | Wersja | Licencja | Zależność |
|---|---|---|---|---|
| 1 | `github.com/coder/websocket` | v1.8.15 | ISC | bezpośrednia |
| 2 | `golang.org/x/sys` | v0.47.0 | BSD 3-Clause | bezpośrednia |
| 3 | `modernc.org/sqlite` | v1.56.0 | BSD 3-Clause | bezpośrednia |
| 4 | `github.com/dustin/go-humanize` | v1.0.1 | MIT | pośrednia |
| 5 | `github.com/google/uuid` | v1.6.0 | BSD 3-Clause | pośrednia |
| 6 | `github.com/mattn/go-isatty` | v0.0.24 | MIT | pośrednia |
| 7 | `github.com/ncruces/go-strftime` | v1.0.0 | MIT | pośrednia |
| 8 | `github.com/remyoudompheng/bigfft` | v0.0.0-20230129092748-24d4a6f8daec | BSD 3-Clause | pośrednia |
| 9 | `modernc.org/libc` | v1.74.4 | BSD 3-Clause | pośrednia |
| 10 | `modernc.org/mathutil` | v1.7.1 | BSD 3-Clause | pośrednia |
| 11 | `modernc.org/memory` | v1.11.0 | BSD 3-Clause | pośrednia |

### A.2. Warstwa Klienta — pakiet npm `danaco-console-client`

| Lp. | Komponent | Wersja | Licencja | Charakter |
|---|---|---|---|---|
| 1 | `@tauri-apps/api` | 2.5.0 | Apache-2.0 OR MIT | produkcyjna |
| 2 | `typescript` | 5.8.3 | Apache-2.0 | narzędzie budowania |
| 3 | `vite` | 6.3.5 | MIT | narzędzie budowania |
| 4 | `vitest` | 3.2.7 | MIT | narzędzie testowe |
| 5 | `jsdom` | ^29.1.1 | MIT | narzędzie testowe (środowisko DOM, prace nad sprawdzianami widoków) |

Zamknięcie zależności: 141 pozycji; rozkład licencji — MIT 126, Apache-2.0 3,
ISC 3, MIT-0 2, BSD-2-Clause 2, BSD-3-Clause 2, „Apache-2.0 OR MIT” 1,
BlueOak-1.0.0 1, CC0-1.0 1. Brak licencji wzajemnych.

### A.3. Warstwa Powłoki — pakiet Cargo `danaco-console-powloka`

Zależności bezpośrednie (bieżący stan repozytorium; kolumna „Sekcja” wskazuje, czy
zależność pochodzi z `[dependencies]` — wchodzi do pliku wynikowego — czy
z `[build-dependencies]` — narzędzie etapu budowania):

| Lp. | Komponent | Sekcja | Wersja rozwiązana | Licencja |
|---|---|---|---|---|
| 1 | `tauri` | `[dependencies]` | 2.11.5 | Apache-2.0 OR MIT |
| 2 | `tauri-plugin-dialog` | `[dependencies]` | 2.7.2 | Apache-2.0 OR MIT |
| 3 | `serde` | `[dependencies]` | 1.0.229 | MIT OR Apache-2.0 |
| 4 | `tauri-build` | `[build-dependencies]` | 2.6.3 | Apache-2.0 OR MIT |

`serde_json` (rozwiązana 1.0.151, MIT OR Apache-2.0) figurowała jako zależność
bezpośrednia w stanie sprzed prac naprawczych; usunięta z `Cargo.toml` jako
nieużywana w pracach nad skryptami wydania, pozostaje w zamknięciu `Cargo.lock` wyłącznie
jako zależność pośrednia wciągana przez `tauri` — patrz rozdz. 7.4.

Zamknięcie zależności: 447 pakietów, z czego 274 zweryfikowane co do
zadeklarowanej licencji. Komponenty na licencji wzajemnej MPL-2.0:
`cssparser` 0.36.0, `cssparser-macros` 0.6.1, `dtoa-short` 0.3.5,
`option-ext` 0.2.0, `selectors` 0.36.1.

**Zestawienie warstwy Powłoki wymaga uzupełnienia przed publikacją** —
patrz rozdz. 7.4 (zastrzeżenie o kompletności).

### A.4. Kroje pisma

| Lp. | Krój | Warianty | Licencja | Plik licencji |
|---|---|---|---|---|
| 1 | Space Grotesk | 500, 600, 700 (latin, latin-ext) | SIL OFL 1.1 | `LICENCJA-space-grotesk.txt` |
| 2 | IBM Plex Sans | 400, 500, 600, 700 (latin, latin-ext) | SIL OFL 1.1 | `LICENCJA-ibm-plex.txt` |
| 3 | IBM Plex Mono | 400, 500, 600 (latin, latin-ext) | SIL OFL 1.1 | `LICENCJA-ibm-plex.txt` |

### A.5. Zestaw ikon i zasoby marki

| Lp. | Zasób | Liczba | Licencja / status |
|---|---|---|---|
| 1 | Ikony wywiedzione z biblioteki Lucide | 78 | ISC |
| 2 | Emblematy Środowisk (własne) | 4 | utwór Licencjodawcy |
| 3 | **Zestaw ikon łącznie — wykaz manifestu `ikony/manifest.json`** | **82** | 78 ISC + 4 własne |
| 4 | Znak marki `logo-danaco.svg` (poza zestawem ikon, w katalogu plików zestawu) | 1 plik | utwór Licencjodawcy |
| 5 | Logotypy, sygnety i ich warianty (katalog zasobów marki) | 12 plików | utwór Licencjodawcy |

Katalog plików SVG zestawu zawiera 83 pliki: 82 pozycje zestawu ikon oraz znak
marki z poz. 4. [README.md](README.md) i [INSTRUKCJA-UZYTKOWANIA.md](INSTRUKCJA-UZYTKOWANIA.md) posługują się zgodnie
liczbą 82 pozycji zestawu ikon i objaśniają ten sam podział rzeczowy wobec
83 plików katalogu — rozdz. 7.6 opisuje weryfikację tej zgodności; kwestia nie
figuruje już w Załączniku B jako nierozstrzygnięta.

### A.6. Zależności zewnętrzne nieobjęte dystrybucją

| Lp. | Zależność | Status |
|---|---|---|
| 1 | Program `claude` (Claude Code CLI) | nie jest dystrybuowany; zapewnia Operator; warunki Dostawcy modelu |
| 2 | Usługa modelu wywoływana kanałem rodzaju `api` | nie jest dystrybuowana; warunki dostawcy usługi |
| 3 | Program `ssh` systemu operacyjnego | wykorzystywany przez most MCP; element systemu operacyjnego |
| 4 | Komponent webview systemu Windows | element systemu operacyjnego |

---

## Załącznik B — wykaz kwestii pozostawionych do decyzji Operatora

Poniższe kwestie nie zostały rozstrzygnięte w dokumentacji projektu. Dokument
nie może zostać opublikowany ani dołączony do dystrybucji przed ich
rozstrzygnięciem i usunięciem oznaczeń **[DO DECYZJI OPERATORA]** z treści.

| Lp. | Rozdział | Kwestia | Skutek nierozstrzygnięcia |
|---|---|---|---|
| 1 | 1.2 | Definicja Licencjobiorcy — osoba fizyczna, podmiot, urządzenie czy grupa kapitałowa | brak podmiotowego zakresu Licencji |
| 2 | 3.2 | Model dystrybucji końcowej — wewnętrzna, instalator, usługa, mieszany | brak podstawy do rozliczenia i kontroli zgodności |
| 3 | 4.3 | Liczba stanowisk i jednostka liczenia licencji | niemożność ustalenia naruszenia zakresu; obowiązek ewidencyjny rozdz. 12.1, oparty na tym samym zakresie podmiotowym, pozostaje bez mierzalnego przedmiotu do czasu rozstrzygnięcia |
| 4 | 4.5 | Odpłatność licencji i model rozliczeń | brak podstawy roszczenia o wynagrodzenie |
| 5 | 5.2 pkt 8 | Zakaz budowy produktu konkurencyjnego — brak, zakaz wąski, zakaz szeroki ograniczony czasowo, zakaz szeroki bezterminowy | klauzula o nieoznaczonym zakresie zagrożona nieskutecznością; do rozstrzygnięcia obowiązuje wykładnia zawężająca |
| 6 | 5.2 pkt 9 | Dopuszczalność badań bezpieczeństwa i testów obciążeniowych przez Licencjobiorcę | sprzeczność zakazu badań z obowiązkiem zgłaszania podatności (rozdz. 13.4) |
| 7 | 5.8 | Kary umowne za naruszenie ograniczeń | odpowiedzialność wyłącznie na zasadach ogólnych |
| 8 | 6.4 | Status ochronny znaków towarowych | ograniczona ochrona oznaczeń |
| 9 | 6.6 | Status informacji zwrotnych od Licencjobiorcy | ryzyko sporu o prawa do zgłoszonych rozwiązań |
| 10 | 7.4 | Sposób prowadzenia inwentarza licencji Komponentów | ryzyko rozjazdu not licencyjnych przy aktualizacji zależności |
| 11 | 7.8 | Zasady licencyjne rozszerzeń i odpowiedzialność za rozszerzenia osób trzecich | brak podstawy do udostępnienia mechanizmu rozszerzeń |
| 12 | 8.4 | Uprawnienie do Aktualizacji, zakres i warunki wsparcia technicznego | brak zobowiązania serwisowego; ryzyko oczekiwań niepokrytych umową |
| 13 | 9.6 | Role w rozumieniu RODO przy topologii z Rdzeniem serwerowym | brak umowy powierzenia przetwarzania przy przeniesieniu Rdzenia |
| 14 | 10.4 | Górna granica odpowiedzialności Licencjodawcy | odpowiedzialność bez limitu kwotowego |
| 15 | 11.2 | Okres obowiązywania Licencji | niepewność co do trwałości uprawnień |
| 16 | 11.3 pkt 2 | Termin wyznaczany na usunięcie naruszenia przed wypowiedzeniem | brak oznaczonego terminu naprawczego; ryzyko sporu o skuteczność wypowiedzenia |
| 17 | 12.2 | Prawo kontroli zgodności licencyjnej | brak narzędzia weryfikacji zakresu korzystania |
| 18 | 13.3 | Okres poufności po zakończeniu Licencji | brak oznaczonego okresu następczego dla informacji niebędących tajemnicą przedsiębiorstwa |
| 19 | 13.4 | Kanał i polityka zgłaszania podatności | brak trybu obsługi zgłoszeń bezpieczeństwa |
| 20 | 14.1 | Prawo właściwe | niepewność co do skuteczności wyłączeń odpowiedzialności |
| 21 | 14.2 | Właściwość sądu i tryb rozstrzygania sporów | brak przewidywalnego forum sporu |
| 22 | 14.7 | Dane identyfikacyjne i kontaktowe Licencjodawcy | dokument niekompletny formalnie |
| 23 | 1.5 | **Przegląd prawny całości dokumentu przed publikacją** | ryzyko nieskuteczności kluczowych klauzul |

Poprzednia redakcja niniejszego załącznika zawierała dodatkową pozycję 11
(rozdz. 7.6 — jednolita liczba pozycji zestawu ikon w dokumentach produktu).
Pozycja została usunięta po weryfikacji bezpośredniej treści [README.md](README.md)
i [INSTRUKCJA-UZYTKOWANIA.md](INSTRUKCJA-UZYTKOWANIA.md), która wykazała, że oba dokumenty posługują się
zgodnie liczbą 82 pozycji zestawu ikon — rozbieżność, dla której pozycja
istniała, nie występuje już w treści dokumentów produktu (rozdz. 7.6). Numeracja
pozostałych pozycji została odpowiednio przesunięta; poniższa uwaga posługuje
się numeracją bieżącą.

**Uwaga o czterech pozycjach nazwanych wprost: zakazie budowy produktu
konkurencyjnego (poz. 5, rozdz. 5.2 pkt 8), dopuszczalności badań bezpieczeństwa
(poz. 6, rozdz. 5.2 pkt 9), terminie na usunięcie naruszenia (poz. 16,
rozdz. 11.3 pkt 2) oraz okresie poufności po zakończeniu Licencji (poz. 18,
rozdz. 13.3).** Pozycje te dotyczą rozstrzygnięć
handlowo-prawnych, które we wcześniejszej redakcji dokumentu zostały przyjęte
w treści bez oznaczenia i bez oparcia w zapisie decyzji architektonicznych ani
w pozostałej dokumentacji projektu (odpowiednio: zakaz konkurencji, zakaz badań
bezpieczeństwa, termin czternastu dni na usunięcie naruszenia, pięcioletni
okres poufności). Ponieważ dokument w kwestiach tej samej wagi — kary umowne
(rozdz. 5.8), okres obowiązywania (rozdz. 11.2), górna granica
odpowiedzialności (rozdz. 10.4) — konsekwentnie stosuje oznaczenie
**[DO DECYZJI OPERATORA]**, przyjęte wcześniej wartości oznaczono w ten sam
sposób i opatrzono wariantami. Nie stanowi to zmiany stanowiska
Licencjodawcy, lecz usunięcie niekonsekwencji redakcyjnej: rozstrzygnięcie
należy do Licencjodawcy i musi zostać wpisane przed publikacją dokumentu.

---

## Załącznik C — metryka weryfikacji faktograficznej

### C.1. Podstawa faktograficzna

Twierdzenia o funkcjach, ograniczeniach, zależnościach i sposobie działania
Oprogramowania zawarte w niniejszym dokumencie oparto na:

1. skonsolidowanym wyniku audytu repozytorium przeprowadzonego pierwotnie na
   stanu sprzed prac naprawczych, uzupełnionym
   ponowną weryfikacją punktową w bieżącym stanie repozytorium tej samej gałęzi po fali
   naprawczych — stanu, z którego zbudowano i uruchomiono
   doręczany Egzemplarz (C.3 pkt 4);
2. dokumencie nadrzędnym projektu (README) oraz zapisie decyzji
   architektonicznych;
3. dokumentacji technicznej towarzyszącej kodowi (kontrakt komunikacyjny —
   [Kontrakty komunikacji](architektura/kontrakty-komunikacji.md); model danych —
   [Model danych](architektura/model-danych.md); mechanizm izolacji —
   [Izolacja i zależności](architektura/izolacja-i-zaleznosci.md); katalog modułów
   i środowisk — `moduly/`, `srodowiska/`; katalog okien operacyjnych w bazie
   i wykaz luk — rozdz. 18 i 20 [README.md](README.md));
4. bezpośrednim odczycie plików repozytorium: deklaracji zależności, plików
   zamknięcia zależności, konfiguracji pakietu Powłoki, wzorca konfiguracji
   środowiska, migracji schematu, manifestu warstwy wizualnej oraz plików
   licencyjnych krojów pisma.

### C.2. Elementy zweryfikowane bezpośrednio w kodzie

| Twierdzenie dokumentu | Źródło weryfikacji |
|---|---|
| Domyślny port nasłuchu 17870 | wzorzec konfiguracji środowiska, konfiguracja Rdzenia |
| Wykaz zmiennych środowiska (pięć pozycji) | wzorzec konfiguracji środowiska |
| Katalog danych `%LOCALAPPDATA%\DanacoConsole` | moduł katalogu danych Rdzenia |
| Nazwa pliku Bazy `danaco-console.db` | moduł pliku bazy Rdzenia |
| Nazwa produktu, identyfikator, wydawca, format instalatora | konfiguracja pakietu Powłoki |
| Brak komponentu samoczynnej aktualizacji | konfiguracja i manifest Powłoki |
| Zależności i wersje trzech warstw | pliki deklaracji i zamknięcia zależności |
| Rodzaje licencji Komponentów warstwy Go | pliki licencyjne w lokalnej pamięci podręcznej modułów |
| Licencje krojów pisma (OFL 1.1) | pliki licencyjne dołączone do zasobów |
| Liczba i źródło pozycji zestawu ikon | manifest warstwy wizualnej |
| Przechowywanie odwołań zamiast treści poświadczeń | schemat tabel kont i punktów dostępu |
| Brak połączeń wychodzących poza kanał modelu i most | przegląd użyć warstwy sieciowej w Rdzeniu |
| Pola koperty Kontraktu: wymagane `type`, `id`, `timestamp`; opcjonalne `sessionId`, `payload`; pola odpowiedzi `status`, `error`; pola strumienia `seq`, `done` | plik Kontraktu, sekcja koperty |
| Wersja rozwiązana `serde` 1.0.229 (zależność bezpośrednia) i `serde_json` 1.0.151 (zależność pośrednia po usunięciu z `Cargo.toml` w pracach nad skryptami wydania) | plik zamknięcia zależności warstwy Rust w bieżącym stanie repozytorium |
| `tauri-build` jest zależnością sekcji `[build-dependencies]`, nie `[dependencies]` | `budowa/desktop/src-tauri/Cargo.toml` |
| Zestaw ikon: 82 pozycje wykazu manifestu (78 wywiedzionych z Lucide, 4 własne) przy 83 plikach katalogu zestawu, z których jeden jest znakiem marki; zgodność [README.md](README.md) i [INSTRUKCJA-UZYTKOWANIA.md](INSTRUKCJA-UZYTKOWANIA.md) z liczbą 82 | manifest warstwy wizualnej, wiązania zestawu w kodzie Klienta, wykaz plików katalogu, odczyt bezpośredni obu dokumentów |
| Zdarzenie `progress.changed` ma konsumenta; dane przykładowe pulpitu dowodzenia usunięte, pulpit odczytuje `session.list`/`window.list`/`channel.list` | `budowa/client/src/mission-control/zrodlo-pulpitu.ts` |
| `session.list`, `session.bind`, `session.focus` wywoływane przez stronę główną (strefa „Sesje w tle”) | `budowa/client/src/strona-glowna/wpiecie-sesji.ts`, `zrodlo-sesji.ts` |
| Kolumny CHECK `rodzaj_kanalu IN ('cli','api','sdk','lokalny')`; kanał `echo` zakłada się jako rodzaj `lokalny` z parametrem `{"adapter":"echo"}` | `budowa/server/internal/store/migracja_001_fundament.sql`; `budowa/server/internal/models/kanal.go` |
| Istnienie `budowa/scripts/wydanie.sh` i `budowa/scripts/pakowanie.sh`; katalog `C:\DanacoConsole_App` zbudowany tą ścieżką i uruchomiony | zawartość skryptów; zawartość katalogu produktu |
| Brak utrwalenia prowenancji — brak tabeli i kolumny w schemacie; wykaz pól tabeli wiadomości | migracje schematu trwałości; moduł prowenancji Rdzenia |
| Niemożność zapisu Punktu dostępu rodzaju `localDirectory` — więz wymagający wiersza urządzenia i brak produkcyjnego źródła tego wiersza | migracja punktów dostępu; przegląd warstwy danych Rdzenia i migracji pod kątem tabeli urządzeń |
| Zakres ustawień stosowanych przy wywołaniu modelu; brak odczytu rodziny `harness.*` | katalog ustawień w migracji; przegląd składania wywołania w Rdzeniu |
| Dostępność cyklu życia Konta z poziomu interfejsu przy braku interfejsu zakładania kanału | warstwa kont i formularz konta w Kliencie; wywołania `channel.list` w widokach |
| Nazwy plików dokumentacji produktu | zawartość katalogu produktu; wykaz dokumentów w [README.md](README.md) |

### C.3. Ograniczenia weryfikacji

1. Zestawienie licencji zamknięcia zależności warstwy Rust jest niepełne —
   zweryfikowano 274 z 447 pakietów (rozdz. 7.4).
2. Weryfikacja miała charakter statyczny: nie uruchamiano procesu budowania
   ani samego Oprogramowania w celu potwierdzenia zachowania w czasie
   działania. Twierdzenia o działaniu funkcji pochodzą z wyniku audytu
   repozytorium.
3. Dokument opisuje stan na dwóch oznaczonych rewizjach tej samej gałęzi:
   rewizji audytu pierwotnego stanu sprzed prac naprawczych oraz bieżącego stanu repozytorium, na której
   punktowo zweryfikowano i zaktualizowano twierdzenia unieważnione falą
   naprawczą (niniejszy punkt 4). Twierdzenie nieoznaczone żadną z tych
   rewizji z osobna obowiązuje na obu. Każda kolejna zmiana zależności,
   warstwy wizualnej albo zakresu funkcjonalnego wymaga ponownej weryfikacji
   rozdz. 7, 10.3 oraz Załącznika A na rewizji faktycznie wydawanej.
4. **Rewizja stanu sprzed prac naprawczych a bieżący stan gałęzi roboczej.** Między rewizją audytu
   pierwotnego stanu sprzed prac naprawczych a HEAD gałęzi roboczej (bieżącego stanu repozytorium) leżą trzy rewizje
   **zatwierdzone**, nie niezatwierdzone zmiany robocze: prac nad skryptami wydania (ścieżka
   wydania — `tauri.conf.json` z `beforeBuildCommand` i `devUrl`,
   `budowa/scripts/wydanie.sh` i `pakowanie.sh`, brama mierząca wywołania,
   usunięcie `serde_json` z zależności bezpośrednich Powłoki), prac nad zasilaniem widoków z rdzenia
   (domknięcie warstwy wizualnej v2.0, pulpit dowodzenia zasilany z Rdzenia
   zamiast danych przykładowych, testy widoków) oraz porządek dokumentacyjny. Egzemplarz zbudowano z zatwierdzonej bieżącego stanu repozytorium;
   niniejszy opis dotyczy tej rewizji, a nie późniejszego stanu katalogu
   roboczego repozytorium, w którym w chwili sporządzania tego zestawienia
   równolegle trwa kolejna, odrębna fala prac nieobjęta niniejszym
   dokumentem. Różnica wobec rewizji
   stanu sprzed prac naprawczych nie ogranicza się do konfiguracji budowy Powłoki: obejmuje również
   pulpit dowodzenia i konsumenta zdarzenia `progress.changed` (rozdz. 10.3
   wykaz ograniczeń, pkt 6 i 16), usunięcie `serde_json` jako zależności
   bezpośredniej (rozdz. 7.4, A.3) oraz stronę główną wywołującą
   `session.list`/`session.bind`/`session.focus` (rozdz. 10.3 pkt 13).
   Twierdzenia niniejszego dokumentu, poza punktami zaktualizowanymi wprost,
   odnoszą się do stanu potwierdzonego w bieżącym stanie repozytorium; różnicę wobec
   stanu sprzed prac naprawczych oznaczono w każdym z punktów, których dotyczy (rozdz. 10.3 wykaz
   ograniczeń, pkt 6, 13, 16, 18). Przed każdym kolejnym Wydaniem należy
   powtórzyć weryfikację na rewizji wówczas faktycznie wydawanej.
5. Twierdzenia o zgodności niniejszego dokumentu z pozostałymi dokumentami
   produktu zweryfikowano przez bezpośredni odczyt tych dokumentów w katalogu
   produktu. Rozbieżność co do liczby pozycji zestawu ikon, opisywana we
   wcześniejszej redakcji niniejszej Licencji, nie występuje w treści
   [README.md](README.md) ani [INSTRUKCJA-UZYTKOWANIA.md](INSTRUKCJA-UZYTKOWANIA.md) na dzień tej weryfikacji — oba
   dokumenty posługują się liczbą 82 zgodnie z manifestem warstwy wizualnej;
   odpowiadająca temu pozycja Załącznika B została usunięta (rozdz. 7.6).

### C.4. Zakres wymagający ponownego przeglądu przy każdym Wydaniu

1. rozdz. 7 i Załącznik A — zestawienie Komponentów i ich licencji, w tym
   wersje rozwiązane zależności bezpośrednich warstwy Powłoki oraz liczba
   i źródła pozycji zestawu ikon;
2. rozdz. 10.3 — ujawnienie rzeczywistego stanu wykonania, ze szczególnym
   uwzględnieniem pkt 6 (źródło danych pulpitu dowodzenia), pkt 13 (wyjątek
   nawigacji platformy — sesje w tle), pkt 16 (konsument `progress.changed`),
   pkt 18 (instalator NSIS Powłoki), pkt 21 (zakres ustawień stosowanych przy
   wywołaniu) i pkt 22 (punkty dostępu rodzaju `localDirectory`);
3. rozdz. 3.4 — zależności zewnętrzne warunkujące działanie, w tym podział na
   czynności dostępne z interfejsu i czynności wymagające komendy Kontraktu;
4. rozdz. 9.2 i 9.4 — zakres danych zapisywanych i przekazywanych na zewnątrz,
   w tym moment powstania wiersza sesji i okna oraz brak utrwalenia
   prowenancji;
5. rozdz. 2.2 — zestaw pól koperty Kontraktu, wiążący co do brzmienia na
   podstawie rozdz. 2.6 pkt 5;
6. rozdz. 1.4 pkt 3 i rozdz. 3.5 — wykaz plików dokumentacji produktu objętej
   Licencją, wraz z ich zgodnością co do stanu opisywanego produktu.

**Niniejszy dokument jest projektem wymagającym przeglądu prawnego przed
publikacją (rozdz. 1.5) oraz rozstrzygnięcia kwestii wskazanych w Załączniku B.**

---

*Koniec dokumentu. Licencja produktu — Akt prawny, wersja 2.0, 2026-08-20.*

---
*Danaco Console — Platforma AI Workspace OS · v2.0 · status Deweloperski*
*© 2026 Danaco Holding Group Sp. z o.o. Wszelkie prawa zastrzeżone — Dariusz Naharnowicz.*
*Warunki korzystania: [Licencja produktu](LICENSE.md). Kontakt: support@danaco-group.pl*
