# Danaco Console — Okno instalatora

| | |
|---|---|
| **Produkt** | Danaco Console |
| **Rodzaj** | Platforma AI Workspace OS |
| **Opis** | Platforma jest wielośrodowiskowym systemem operacyjnym dla sztucznej inteligencji, integrującym komunikację, zarządzanie wiedzą, tworzenie treści, projektowanie, automatyzacje procesów oraz rozwój oprogramowania. |
| **Producent** | Danaco Holding Group Sp. z o.o. |
| **Twórca** | Dariusz Naharnowicz |
| **Wersja** | v2.0 |
| **Status** | Deweloperski |
| **Data** | 2026-08-20 |

**Informacje szczegółowe dokumentu:**

| | |
|---|---|
| **Tytuł** | Okno instalatora |
| **Klasa dokumentu** | Specyfikacja docelowa |
| **Odbiorcy** | projektant — buduje ekrany instalatora · deweloper — buduje logikę instalacji i pakietowanie · Operator — instaluje produkt |
| **Przeznaczenie** | Ustala kompletną budowę i przebieg okna instalatora aplikacji — okna przedaplikacyjnego (etap 0 przepływu), prowadzącego przez sześć kroków od warunków licencji do zakończenia i pierwszego uruchomienia, wraz z wariantami technicznego wykonania pakietu na trzech systemach operacyjnych. |
| **Zakres** | charakter okna, wydawca i wiarygodność pakietu, wymagania systemowe, budowa okna i mapa stref, sześć kroków, sześć składników instalacji, walidacje i komunikaty zgodne z zasadą zero blokad, warianty pakietowania na trzech systemach operacyjnych, wycofanie, naprawa i aktualizacja, pełny wykaz etykiet, żetonów i komponentów, stany kontrolek, przebieg nawigacji, punkty łamania, skróty klawiszowe, dostępność, kryteria odbioru |
| **Poza zakresem** | konfiguracja po pierwszym uruchomieniu i stan faktyczny dróg budowy produktu — [Instalacja i konfiguracja](../INSTALACJA-I-KONFIGURACJA.md); rama okna aplikacji (belka, szyna, wstążka, pasek stanu) — [Rama okna](rama-okna.md); okno startowe, rejestracja i logowanie jako kolejne etapy przepływu — [Przepływ okien](przeplyw-okien.md); okno instrukcji użytkowania — [Okno instrukcji](okno-instrukcji.md) |
| **Dokument nadrzędny** | [Elementy okien](elementy-okien.md) |
| **Dokumenty powiązane** | [Elementy okien](elementy-okien.md) · [Instalacja i konfiguracja](../INSTALACJA-I-KONFIGURACJA.md) · [Przepływ okien](przeplyw-okien.md) · [Okno instrukcji](okno-instrukcji.md) · [Rama okna](rama-okna.md) · [System wizualny](system-wizualny.md) · [Katalog komponentów](katalog-komponentow.md) |
| **Prototypy odniesienia** | `design/05-okna/platformowe/instalator.html` |
| **Źródła normatywne** | `design/05-okna/platformowe/instalator.html` (wartości wymagań, sześciu składników, treści kolumny tożsamości, przebieg kroku 3) · `design/zasoby/wejscie.css` (rodzina okien wejściowych, klasy `.we-*`) · `design/zasoby/css/komponenty.css` (przyciski, pole wyboru, pasek postępu) · `design/zasoby/zetony/zetony.css` (punkty łamania `--dn-bp-*` i pozostałe żetony) · [Instalacja i konfiguracja](../INSTALACJA-I-KONFIGURACJA.md) rozdz. 5.1 i 5.6 (stan faktyczny dróg budowy i pakowania) |
| **Zasada nadrzędna** | Instalator jest oknem przedaplikacyjnym — nie niesie ramy aplikacji (szyny nawigacji, wstążki, paska stanu). Korzysta wyłącznie z wyglądu okien wejściowych. Żaden wybór składnika opcjonalnego nie jest ostateczny — wszystkie można zmienić po instalacji. Obowiązuje zasada zero blokad (`komponenty.css`, w. 7): żadna kontrolka instalatora nie traci klikalności — niegotowość komunikuje się opisem albo komunikatem, nigdy odebraniem możliwości kliknięcia. |

---

## Spis treści

1. [Czym jest okno instalatora](#1-czym-jest-okno-instalatora)
   - [1.1 Okno przedaplikacyjne](#11-okno-przedaplikacyjne)
   - [1.2 Wydawca, podpis i wiarygodność pakietu](#12-wydawca-podpis-i-wiarygodność-pakietu)
   - [1.3 Stan faktyczny a stan docelowy](#13-stan-faktyczny-a-stan-docelowy)
   - [1.4 Komendy kontraktu](#14-komendy-kontraktu)
   - [1.5 Słownik terminów pakietowania](#15-słownik-terminów-pakietowania)
2. [Wymagania systemowe](#2-wymagania-systemowe)
3. [Budowa okna i mapa stref](#3-budowa-okna-i-mapa-stref)
   - [3.1 Układ dwukolumnowy](#31-układ-dwukolumnowy)
   - [3.2 Kolumna tożsamości](#32-kolumna-tożsamości)
   - [3.3 Kolumna kroku](#33-kolumna-kroku)
   - [3.4 Punkty łamania](#34-punkty-łamania)
   - [3.5 Cykl życia procesu instalatora](#35-cykl-życia-procesu-instalatora)
4. [Sześć kroków instalacji](#4-sześć-kroków-instalacji)
   - [4.1 Tor kroków i jego stany](#41-tor-kroków-i-jego-stany)
   - [4.2 Krok 1 — warunki licencji](#42-krok-1--warunki-licencji)
   - [4.3 Krok 2 — katalog instalacji](#43-krok-2--katalog-instalacji)
   - [4.4 Krok 3 — składniki](#44-krok-3--składniki)
   - [4.5 Krok 4 — skróty i uruchamianie](#45-krok-4--skróty-i-uruchamianie)
   - [4.6 Krok 5 — instalacja](#46-krok-5--instalacja)
   - [4.7 Krok 6 — zakończenie](#47-krok-6--zakończenie)
5. [Składniki instalacji](#5-składniki-instalacji)
6. [Walidacje i komunikaty](#6-walidacje-i-komunikaty)
7. [Warianty pakietowania na trzech systemach operacyjnych](#7-warianty-pakietowania-na-trzech-systemach-operacyjnych)
   - [7.1 Windows — NSIS i WiX/MSI](#71-windows--nsis-i-wixmsi)
   - [7.2 macOS — DMG i pkg](#72-macos--dmg-i-pkg)
   - [7.3 Linux — DEB, RPM, AppImage](#73-linux--deb-rpm-appimage)
   - [7.4 Podpis kodu](#74-podpis-kodu)
   - [7.5 Instalacja dla użytkownika a dla wszystkich użytkowników](#75-instalacja-dla-użytkownika-a-dla-wszystkich-użytkowników)
   - [7.6 Aktualizacje](#76-aktualizacje)
   - [7.7 Czego nie warto obiecywać](#77-czego-nie-warto-obiecywać)
   - [7.8 Zalecenie](#78-zalecenie)
8. [Wycofanie, naprawa i aktualizacja](#8-wycofanie-naprawa-i-aktualizacja)
9. [Etykiety interfejsu](#9-etykiety-interfejsu)
10. [Żetony i komponenty](#10-żetony-i-komponenty)
11. [Stany kontrolek](#11-stany-kontrolek)
12. [Przebieg nawigacji](#12-przebieg-nawigacji)
13. [Skróty klawiszowe](#13-skróty-klawiszowe)
14. [Dostępność](#14-dostępność)
15. [Scenariusze instalacji](#15-scenariusze-instalacji)
16. [Kryteria odbioru](#16-kryteria-odbioru)

---

## 1. Czym jest okno instalatora

### 1.1 Okno przedaplikacyjne

Instalator jest **oknem przedaplikacyjnym** — etapem 0 przepływu, przed pierwszym uruchomieniem
aplikacji. Prototyp `design/05-okna/platformowe/instalator.html` deklaruje siebie atrybutem
`data-prototyp-okno` jako „Instalator aplikacji — okno przed uruchomieniem, wygląd okien
wejściowych” oraz atrybutem podpisu jako „etap 0 przepływu · przed pierwszym uruchomieniem
aplikacji”.

Z bycia oknem przedaplikacyjnym wynika najważniejsza cecha budowy: **instalator nie niesie ramy
aplikacji.** Nie ma szyny nawigacji, wstążki poziomej ani paska stanu opisanych w
[Ramie okna](rama-okna.md). Nie jest oknem wnętrza aplikacji — jest oknem, które aplikację dopiero
instaluje, więc rama, która zakłada istnienie zalogowanej sesji Operatora, środowisk i modułów, nie
ma tu jeszcze czego opisywać. Instalator korzysta z tej samej rodziny wyglądu co okno startowe i
okno logowania: rodziny okien wejściowych (`design/zasoby/wejscie.css`, klasy `.we-*`), opisanej
w rozdz. 3.

```
   Przepływ uruchomienia — cztery etapy poprzedzające pełną ramę aplikacji
   ┌──────────────┐   ┌──────────────┐   ┌──────────────────┐   ┌─────────────┐
   │  etap 0      │   │  etap 1      │   │  etap 2          │   │  etap 3+    │
   │  INSTALATOR  │──►│  okno        │──►│  rejestracja i   │──►│  strona     │
   │  (ten dok.)  │   │  startowe    │   │  logowanie       │   │  główna     │
   └──────────────┘   └──────────────┘   └──────────────────┘   └─────────────┘
     bez ramy           bez ramy           bez ramy               pełna rama
     .we-* okno          .we-* okno          .we-* okno            .dn-* rama
```

Legenda: instalator jest pierwszym z okien wejściowych; wszystkie trzy pierwsze etapy dzielą tę samą
rodzinę wyglądu `.we-*`; ramę aplikacji (`.dn-*`, opisaną w [Ramie okna](rama-okna.md)) niesie dopiero
strona główna i okna wnętrza aplikacji, otwierane po zalogowaniu.

Instalator otwiera się jako osobne okno powłoki natywnej, bez paska tytułu systemu — pasek tytułu
własny opisany w rozdz. 3.2 pełni tę rolę i jednocześnie działa jako uchwyt do przesuwania okna,
skoro instalator nie ma standardowej belki systemowej, za którą chwyta się okno.

### 1.2 Wydawca, podpis i wiarygodność pakietu

Pakiet instalacyjny jest podpisany cyfrowo certyfikatem wydawcy. Kolumna tożsamości (rozdz. 3.2)
deklaruje wydawcę oraz wskazuje, gdzie sprawdzić autentyczność pakietu — dosłowny tekst stopki
kolumny tożsamości brzmi: „Danaco Core sp. z o.o. — wydawca oprogramowania. Pakiet podpisany
certyfikatem wydawcy; suma kontrolna widoczna we właściwościach pliku.”

| Cecha | Wartość z prototypu |
|---|---|
| Wydawca oprogramowania (tekst prototypu) | „Danaco Core sp. z o.o.” |
| Podpis pakietu | certyfikat wydawcy |
| Suma kontrolna | widoczna we właściwościach pliku |
| Wersja pakietu instalatora | „2.4.0 (kompilacja 2418) · pakiet podpisany cyfrowo” — wersja pliku instalatora, odrębna od wersji produktu v2.0 |

**Rozbieżność zweryfikowana ze źródłem konfiguracyjnym.** [Instalacja i konfiguracja](../INSTALACJA-I-KONFIGURACJA.md)
rozdz. 5.6 opisuje rzeczywistą konfigurację pakowania Tauri (`tauri.conf.json`), w której pole
`publisher` niesie wartość „Danaco Holding Group Sp. z o.o.” — nazwę zgodną z metryką produktową
niniejszego zbioru dokumentacji (wiersz „Producent” powyżej). Tekst prototypu instalatora podaje
„Danaco Core sp. z o.o.”. Dwa źródła normatywne niosą dwie różne nazwy wydawcy dla tego samego
pakietu. **[DO DECYZJI OPERATORA]** — która nazwa wydawcy jest wiążąca dla stopki kolumny
tożsamości instalatora: nazwa z pola `publisher` konfiguracji pakowania, czy nazwa użyta w
prototypie. Do rozstrzygnięcia dokument nie ujednolica żadnej z dwóch wartości domysłem — obie są
przywołane dosłownie, każda ze swojego źródła.

### 1.3 Stan faktyczny a stan docelowy

Niniejszy dokument jest specyfikacją docelową: opisuje okno instalatora, jakie **ma powstać**.
[Instalacja i konfiguracja](../INSTALACJA-I-KONFIGURACJA.md) rozdz. 5.1, dokument klasy Stan
wdrożenia, stwierdza wprost: „Danaco Console w wersji v2.0 nie posiada klasycznego instalatora —
nie ma pakietu MSI ani gotowego pliku `setup.exe` do pobrania i uruchomienia. Instalacja polega na
zbudowaniu produktu ze źródeł i złożeniu artefaktów w katalogu produktu.” Ten sam dokument, rozdz.
5.6, potwierdza że konfiguracja pakowania Windows (NSIS) w `tauri.conf.json` jest **przygotowana**
(`"bundle": { "active": true, "targets": ["nsis"] }`), lecz jej wytworzenie nie należy do drogi
wydania — skrypt wydania buduje powłokę z przełącznikiem `--no-bundle`, świadomie pomijając
pakowanie.

Ten sam dokument odnotowuje też ograniczenie zasadnicze dla zakresu rozdz. 7 niniejszej
specyfikacji: instalator NSIS, gdyby został wytworzony bez dalszych zmian, zainstalowałby
**wyłącznie powłokę** wraz z pakietem interfejsu — nie obejmuje binarki rdzenia
`danaco-console.exe`, której powłoka szuka obok siebie. Pełny instalator produktu (rdzeń + powłoka)
wymaga dołączenia rdzenia jako zasobu zewnętrznego pakietu, co [Instalacja i konfiguracja](../INSTALACJA-I-KONFIGURACJA.md)
rozdz. 23.4, poz. 4, odnotowuje jako kwestię **[DO DECYZJI OPERATORA]**: „Czy wprowadzić pełny
instalator produktu obejmujący rdzeń i czy go podpisywać”. Niniejsza specyfikacja opisuje okno
instalatora zakładające, że ta decyzja zapadnie na „tak” — sześć kroków i sześć składników
niniejszego dokumentu (rozdz. 4–5) obejmuje rdzeń, powłokę i cztery składniki dodatkowe jako jeden
spójny pakiet.

```
   Trzy warstwy prawdy o instalatorze
   ┌─────────────────────────────────────────────────────────────────────┐
   │  STAN FAKTYCZNY (v2.0)                                              │
   │  brak klasycznego instalatora · budowa ze źródeł · --no-bundle      │
   ├─────────────────────────────────────────────────────────────────────┤
   │  PRZYGOTOWANIE CZĘŚCIOWE                                            │
   │  konfiguracja NSIS gotowa, ale pakuje wyłącznie powłokę             │
   ├─────────────────────────────────────────────────────────────────────┤
   │  STAN DOCELOWY (niniejszy dokument)                                 │
   │  sześć kroków · sześć składników · rdzeń + powłoka w jednym pakiecie│
   └─────────────────────────────────────────────────────────────────────┘
```

### 1.4 Komendy kontraktu

**Nie dotyczy.** Instalator działa przed nawiązaniem połączenia klienta z rdzeniem — pierwsza
komenda kontraktu, `connection.hello`, zawiązująca wersję protokołu i zdolności obu stron, jest
wywoływana dopiero z okna startowego (etap 1 przepływu, [Przepływ okien](przeplyw-okien.md)), nie z
okna instalatora (etap 0). Instalator nie ma z czym prowadzić wymiany komend: proces rdzenia nie
jest jeszcze uruchomiony, konto Operatora nie istnieje, żadna sesja ani okno komunikacji z sensu
[Kontraktów komunikacji](../architektura/kontrakty-komunikacji.md) nie zostały jeszcze założone.
Sześć kroków instalatora (rozdz. 4) operuje wyłącznie na stanie lokalnym procesu instalującego —
wyborach Operatora zebranych w pamięci procesu i zapisanych na dysku w kroku 5 — bez udziału
protokołu komunikacji opisanego w `budowa/shared/contract.json`. Wykaz jest tu zapisany z jawną
adnotacją „nie dotyczy”, zgodnie z wymogiem rozdz. 4.2 standardu redakcyjnego, nie pominięty
milczeniem.

### 1.5 Słownik terminów pakietowania

Rozdział 7 posługuje się terminologią właściwą pakietowaniu aplikacji natywnych na trzech systemach
docelowych. Terminy zebrane tu raz, w miejscu pierwszego pełnego użycia, zgodnie z rozdz. 5.3
standardu redakcyjnego zbioru.

| Termin | Definicja w kontekście niniejszego dokumentu |
|---|---|
| NSIS | Nullsoft Scriptable Install System — skryptowalny system budowy instalatorów Windows, jeden z dwóch formatów wytwarzanych przez Tauri (rozdz. 7.1) |
| WiX / MSI | zestaw narzędzi budujących pakiety w formacie Windows Installer (`.msi`) — format deklaratywny, obsługiwany natywnie przez system Windows |
| Bootstrapper | program rozruchowy uruchamiany przed właściwym pakietem MSI; w narzędziu WiX tę rolę pełni komponent Burn — pozwala pokazać własne okno o dowolnym wyglądzie, zanim silnik MSI zainstaluje pakiet po cichu w tle (rozdz. 7.1) |
| DMG | obraz dysku macOS (`.dmg`), zwyczajowy nośnik dystrybucji aplikacji tym systemem, otwierany gestem „przeciągnij do Aplikacji” (rozdz. 7.2) |
| `.pkg` | format instalatora macOS pozwalający na wybór miejsca instalacji i składników, lecz z wyglądem wyznaczonym przez systemowy Instalator (rozdz. 7.2) |
| DEB | format pakietu binarnego dystrybucji Linux pochodnych Debiana (w tym Ubuntu), instalowany menedżerem pakietów systemu (rozdz. 7.3) |
| RPM | format pakietu binarnego dystrybucji Linux pochodnych Red Hat, instalowany menedżerem pakietów systemu (rozdz. 7.3) |
| AppImage | format pakietu Linux niewymagający instalacji — pojedynczy plik wykonywalny niosący aplikację wraz z zależnościami (rozdz. 7.3) |
| Notaryzacja | dodatkowy krok po podpisaniu pakietu macOS, w którym producent systemu potwierdza brak złośliwej zawartości — bez niego Gatekeeper blokuje uruchomienie (rozdz. 7.4) |
| SmartScreen | mechanizm ostrzegawczy Windows uruchamiany przy pliku wykonywalnym bez zaufanego podpisu (rozdz. 7.4) |
| Gatekeeper | mechanizm ostrzegawczy macOS uruchamiany przy aplikacji bez podpisu i notaryzacji (rozdz. 7.4) |
| WebView2 | komponent systemowy Windows osadzający silnik przeglądarki w aplikacji natywnej; wymagany przez powłokę Tauri, doinstalowywany automatycznie, gdy go brak (rozdz. 4.6, 5) |
| Aktualizacja różnicowa | mechanizm pobierający wyłącznie zmianę między wersjami pakietu zamiast całego pakietu na nowo (rozdz. 7.6) |
| Instalacja dla użytkownika / dla wszystkich użytkowników | dwa warianty zasięgu instalacji: katalog profilu Operatora bez podniesienia uprawnień, albo katalog wspólny wymagający uprawnień administratora (rozdz. 7.5) |
| Podpis pakietu (certyfikat wydawcy) | poświadczenie kryptograficzne tożsamości wydawcy dołączone do pliku pakietu, sprawdzane przez system operacyjny przed uruchomieniem (rozdz. 1.2, 7.4) |
| Suma kontrolna | wartość skrótu kryptograficznego pliku instalatora, pozwalająca Operatorowi potwierdzić, że pobrany plik nie został zmieniony po opublikowaniu (rozdz. 1.2) |

---

## 2. Wymagania systemowe

Wartości dosłowne z kolumny tożsamości prototypu (rozdz. 3.2), zapisane tam jako cztery pozycje
wykazu, każda z etykietą pogrubioną i treścią.

| Pozycja | Etykieta prototypu | Wymóg |
|---|---|---|
| System operacyjny | „System operacyjny” | Windows 10 w wersji 1809 lub nowszej · macOS 12 · Ubuntu 22.04 LTS |
| Procesor i pamięć | „Procesor i pamięć” | 4 rdzenie, 8 GB RAM (zalecane 16 GB przy modelach lokalnych) |
| Dysk i sieć | „Dysk i sieć” | 2,4 GB wolnego miejsca · połączenie z siecią przy pierwszym uruchomieniu |
| Wersja instalowana | „Wersja instalowana” | 2.4.0 (kompilacja 2418) · pakiet podpisany cyfrowo |

Instalator sprawdza spełnienie wymagań przed rozpoczęciem kopiowania (krok 5, rozdz. 4.6).
Niespełnione wymaganie twarde (system poniżej minimum, brak miejsca na dysku) zatrzymuje przebieg
z komunikatem opisanym w rozdz. 6; wymaganie zalecane (16 GB RAM) daje ostrzeżenie, nie zatrzymanie —
zgodnie z zasadą zero blokad (rozdz. 6), ostrzeżenie nie odbiera możliwości kontynuowania.

Wymaganie sieciowe („połączenie z siecią przy pierwszym uruchomieniu”) dotyczy etapu 1 przepływu
(okno startowe), nie samego instalatora: instalator kopiuje pliki lokalnie i nie wymaga sieci do
ukończenia sześciu kroków — sieć staje się potrzebna dopiero przy rejestracji i logowaniu (etap 2).

**Zgodność wymogu dyskowego z sumą składników.** Rozdz. 5 podaje sumę zajętości sześciu składników
przy zestawie domyślnym: 1 842 MB, czyli w przybliżeniu 1,8 GB. Wymaganie dyskowe z niniejszego
rozdziału wynosi 2,4 GB — o około 0,6 GB więcej niż suma samych plików docelowych. Różnica jest
spójna z tym, że instalator w kroku 5 (rozdz. 4.6) rozpakowuje archiwa: potrzebuje miejsca
tymczasowego na spakowaną postać składników równolegle z ich rozpakowaną postacią, zanim skasuje
plik tymczasowy. Wartość 2,4 GB jest wartością wymogu wejściowego, sprawdzaną na starcie kroku 5,
niezależnie od tego, ile z sześciu składników Operator faktycznie zaznaczył — wykaz wymagań w
kolumnie tożsamości (rozdz. 3.2) jest stały na wszystkich sześciu krokach i nie przelicza się wraz
z wyborem składników w kroku 3; przelicza się wyłącznie podsumowanie zajętości w samym kroku 3
(`.in-suma`, rozdz. 5).

**Kolejność sprawdzeń w kroku 5.** Cztery wymagania z tabeli powyżej nie są równoważne co do
skutku niespełnienia — dwa są twarde (system operacyjny, dysk), dwa dają wyłącznie ostrzeżenie
(procesor i pamięć poniżej zalecanej, lecz powyżej minimum). Kolejność sprawdzenia w kroku 5 jest:
najpierw system operacyjny (rozdz. 6: „System nie spełnia wymagań minimalnych”, zatrzymanie),
następnie miejsce na dysku (rozdz. 6: „Za mało wolnego miejsca…”, zatrzymanie), na końcu procesor i
pamięć (rozdz. 6: „Zalecane 16 GB RAM…”, ostrzeżenie bez zatrzymania). Sprawdzenie sieci nie
występuje w tym kroku — wymóg sieciowy z tabeli powyżej dotyczy etapu 1, nie etapu 0, więc krok 5
instalatora nie testuje łączności.

---

## 3. Budowa okna i mapa stref

### 3.1 Układ dwukolumnowy

Instalator dziedziczy układ okna wejściowego z `design/zasoby/wejscie.css`: siatka dwukolumnowa
o stałych proporcjach `384px minmax(0, 1fr)`, szerokości docelowej całego okna `min(1040px, 100%)`
i minimalnej wysokości `696px`. Nad tą siatką stoi jeden dodatkowy rząd — własny pasek tytułu
instalatora (rozdz. 3.2), obecny wyłącznie w tym oknie spośród rodziny okien wejściowych, bo
pozostałe okna wejściowe (startowe, logowania) otwierają się już wewnątrz powłoki natywnej z
własną belką aplikacji, a instalator działa przed jej uruchomieniem.

```
   Mapa stref okna instalatora — szerokość ≥ 1000 px
   ┌─────────────────────────────────────────────────────────────────────────┐
   │ [znak]  Instalacja programu Danaco Console         [uchwyt] [ ✕ ]       │  .in-belka — 48 px
   ├─────────────────────────────────┬───────────────────────────────────────┤
   │  KOLUMNA TOŻSAMOŚCI             │  KOLUMNA KROKU                        │
   │  .we-marka — 384 px             │  .we-panel — minmax(0, 1fr)           │
   │  tło: --dn-rama                 │  tło: --dn-powierzchnia               │
   │                                 │                                       │
   │  ● Danaco Console               │  Instalacja · krok 3 z 6              │
   │  AI Operating Environment       │  Składniki instalacji                 │
   │                                 │  {wprowadzenie kroku}                 │
   │  ✓ Wersja instalowana           │                                       │
   │  ✓ System operacyjny            │  ┌──────────────────────────────┐     │
   │  ✓ Procesor i pamięć            │  │  .we-kroki — tor sześciu     │     │
   │  ✓ Dysk i sieć                  │  │  kroków (rozdz. 4.1)         │     │
   │                                 │  └──────────────────────────────┘     │
   │                                 │                                       │
   │                                 │  {treść właściwa kroku}               │
   │                                 │                                       │
   │  Danaco Core sp. z o.o. —       │  ┌──────────────────────────────┐     │
   │  wydawca oprogramowania.        │  │  .we-alarm — komunikat       │     │
   │  Pakiet podpisany certyfikatem  │  │  (gdy dotyczy)               │     │
   │  wydawcy; suma kontrolna        │  └──────────────────────────────┘     │
   │  widoczna we właściwościach     │                                       │
   │  pliku.                         │  ( Przerwij )      ( Wstecz )( Dalej) │  .we-stopka
   └─────────────────────────────────┴───────────────────────────────────────┘
```

Legenda: strefa lewa jest stała na każdym z sześciu kroków (rozdz. 3.2); strefa prawa zmienia treść
właściwą zależnie od kroku bieżącego, ale zachowuje trzy elementy stałe swojej struktury: nagłówek
kroku, tor kroków i pasek działań w stopce.

### 3.2 Kolumna tożsamości

Kolumna tożsamości (`.we-marka`) jest **stała na wszystkich sześciu krokach** — nie zmienia treści
w trakcie przebiegu instalacji. Niesie w kolejności z góry na dół:

1. Godło marki (`.we-marka-godlo`) z kropką tętniącą animacją `we-tetno` co `--dn-czas-tetno` —
   jedyny ruch całej kolumny tożsamości, wyłączany przy `prefers-reduced-motion: reduce`.
2. Nazwę produktu (`.we-marka-nazwa`): „Danaco Console”.
3. Motto (`.we-marka-motto`): „AI Operating Environment”.
4. Wykaz czterech pozycji wymagań systemowych (`.we-marka-lista`, rozdz. 2), każda z ikoną,
   etykietą pogrubioną i treścią w jednym wierszu tekstu.
5. Stopkę kolumny (`.we-marka-stopka`): tekst wydawcy i wiarygodności pakietu (rozdz. 1.2),
   oddzielony górną kreską `--dn-rama-obrys`.

Kolumna tożsamości znika całkowicie poniżej progu 1000 px (rozdz. 3.4) — instalator na wąskim
oknie pokazuje wyłącznie kolumnę kroku, bo w tej szerokości nie ma miejsca na obie kolumny bez
ściśnięcia treści formularza.

**Zachowanie kolumny tożsamości w obu motywach.** Żeton `--dn-rama` i pochodne (`--dn-rama-tekst`,
`--dn-rama-tekst-2`, `--dn-rama-hover`, `--dn-rama-obrys`) są zadeklarowane w `zetony.css` w sekcji
opisanej wprost komentarzem „RAMA KOKPITU — stała w obu motywach”: wartość `--dn-rama` (`var(--dn-szary-925)`,
odcień niemal czarny) nie zmienia się między motywem jasnym a ciemnym. Kolumna kroku
(`.we-panel`), przeciwnie, stoi na żetonie `--dn-powierzchnia`, który **zmienia się** między
motywami — biel (`--dn-szary-0`) w motywie jasnym, ciemny grafit (`--dn-szary-900`) w motywie
ciemnym. Skutek: w motywie jasnym kolumna tożsamości kontrastuje wyraźnie z jasną kolumną kroku; w
motywie ciemnym obie kolumny mają zbliżoną jasność, a jedyną granicą między nimi pozostaje kreska
`--dn-rama-obrys` — komentarz źródłowy `wejscie.css` odnotowuje to wprost jako zamierzone
zachowanie, nie usterkę kontrastu.

| Żeton | Rola | Motyw jasny | Motyw ciemny |
|---|---|---|---|
| `--dn-rama` | tło kolumny tożsamości | `--dn-szary-925` (stałe) | `--dn-szary-925` (stałe) |
| `--dn-powierzchnia` | tło kolumny kroku | `--dn-szary-0` (biel) | `--dn-szary-900` (grafit) |
| `--dn-rama-obrys` | kreska między kolumnami | `rgba(255,255,255,0.10)` (stałe) | `rgba(255,255,255,0.10)` (stałe) |

### 3.3 Kolumna kroku

Kolumna kroku (`.we-panel`) zmienia treść zależnie od kroku bieżącego, lecz zachowuje trzy elementy
stałej struktury w każdym z sześciu kroków:

| Element stały | Klasa | Treść |
|---|---|---|
| Nagłówek kroku | `.we-glowa` z `.we-nadtytul`, `.we-tytul`, `.we-lid` | nadtytuł „Instalacja · krok N z 6”, tytuł kroku, jednozdaniowe wprowadzenie |
| Tor kroków | `.we-kroki` | sześć pozycji stanu (rozdz. 4.1), obecny na każdym ekranie |
| Pasek działań | `.we-stopka` z `.in-dzialania` | trzy przyciski: „Przerwij”, „Wstecz”, „Dalej” (rozdz. 4.1, 11) |

Między torem kroków a paskiem działań stoi treść właściwa kroku — różna dla każdego z sześciu
(rozdz. 4.2–4.7) — oraz opcjonalny komunikat `.we-alarm` (rozdz. 6), gdy krok niesie ostrzeżenie
albo informację.

### 3.4 Punkty łamania

| Próg szerokości | Zachowanie układu | Źródło |
|---|---|---|
| ≥ 1000 px | pełny układ dwukolumnowy: kolumna tożsamości 384 px + kolumna kroku elastyczna | `we-okno { grid-template-columns: 384px minmax(0, 1fr) }` |
| < 1000 px | kolumna tożsamości znika (`display: none`); widoczna wyłącznie kolumna kroku w pełnej szerokości okna | `@media (max-width: 1000px) { .we-marka { display: none } }` |
| ≤ 640 px (`--dn-bp-w1`) | dopełnienie kolumny kroku i sceny zmniejszone z `--dn-od-8`/`--dn-od-6` na `--dn-od-5`/`--dn-od-4`; wiersz składnika (rozdz. 5) traci trzecią kolumnę rozmiaru, rozmiar schodzi pod nazwę | `@media (max-width: 640px)` w `wejscie.css` oraz lokalnie w `instalator.html` |

Próg 1000 px nie odpowiada żadnemu nazwanemu żetonowi `--dn-bp-*` — najbliższy niższy jest
`--dn-bp-w2` (960 px). Wartość 1000 px jest zapisana w arkuszu `wejscie.css` wprost, jako liczba
pikseli, nie jako odwołanie do żetonu. Próg 640 px pokrywa się dokładnie z `--dn-bp-w1`.

**[DO DECYZJI OPERATORA]** — czy próg 1000 px ma otrzymać własny żeton w rodzinie `--dn-bp-*`, czy
pozostać wartością lokalną arkusza `wejscie.css`. Do czasu rozstrzygnięcia wykaz punktów łamania tego
okna niesie jedną wartość bez pokrycia w `design/zasoby/zetony/zetony.css`.

```
   Skutek przejścia progu 1000 px na kolumnę tożsamości
   ┌───────────────┬──────────────────────┐          ┌──────────────────────────────┐
   │  KOLUMNA      │  KOLUMNA KROKU       │   ──►    │      KOLUMNA KROKU           │
   │  TOŻSAMOŚCI   │                      │ < 1000px │    (pełna szerokość)         │
   │  384 px       │  elastyczna          │          │                              │
   └───────────────┴──────────────────────┘          └──────────────────────────────┘
                  ≥ 1000 px                              brak kolumny tożsamości
```

### 3.5 Cykl życia procesu instalatora

Instalator jest procesem odrębnym od procesu aplikacji docelowej — kończy działanie po kroku 6
(rozdz. 4.7), niezależnie od tego, czy Operator zaznaczył uruchomienie aplikacji. Cztery stany
procesu, w kolejności przebiegu:

| Stan procesu | Wejście w stan | Wyjście ze stanu |
|---|---|---|
| Uruchomiony, przed potwierdzeniem | otwarcie pliku pakietu przez Operatora | akceptacja licencji (krok 1) albo zamknięcie okna |
| Zbieranie wyborów (kroki 1–4) | akceptacja licencji | kliknięcie „Dalej” na kroku 4 albo „Przerwij” na dowolnym z kroków 1–4 |
| Zapis na dysku (krok 5) | rozpoczęcie kroku 5 | zakończenie kopiowania (przejście do kroku 6) albo przerwanie z cofnięciem zapisu (rozdz. 6, scenariusz D w rozdz. 15) |
| Zakończony (krok 6) | zakończenie kopiowania | zamknięcie okna instalatora — z przejściem do okna startowego albo bez niego, zależnie od opcji „Uruchom Danaco Console” |

Stany „Zbieranie wyborów” i „Zapis na dysku” różnią się zasadniczo dostępnością przycisku
„Przerwij”: w pierwszym stanie przerwanie nie usuwa niczego z dysku (rozdz. 4.6 — instalator nie
zapisuje nic przed krokiem 5); w drugim przerwanie żąda potwierdzenia i usuwa pliki skopiowane do
tej chwili (rozdz. 6). Proces instalatora nie utrzymuje stanu między własnymi uruchomieniami —
ponowne otwarcie pliku pakietu po przerwaniu zaczyna przebieg od kroku 1, z zestawem domyślnym
opisanym w rozdz. 5, nie od miejsca przerwania.

---

## 4. Sześć kroków instalacji

### 4.1 Tor kroków i jego stany

```
[1] Licencja ─► [2] Katalog ─► [3] Składniki ─► [4] Skróty ─► [5] Instalacja ─► [6] Zakończenie
   akceptacja      ścieżka        wybór            uruchamianie    postęp           gotowe
```

Tor kroków (`.we-kroki`, prototyp: „Instalacja · krok 3 z 6”) towarzyszy każdemu z sześciu ekranów
i niesie stan każdej z sześciu pozycji jednocześnie — Operator widzi cały przebieg, nie tylko krok
bieżący. Stan pozycji jest atrybutem danych `data-stan` na elemencie `.we-krok`, zweryfikowanym w
`wejscie.css` — cztery wartości, nie dwie:

| Wartość `data-stan` | Znaczenie | Wygląd znaku pozycji (`.we-krok-znak`) | Przykład z prototypu |
|---|---|---|---|
| `gotowy` | krok już ukończony | obrys `--dn-sukces-obrys`, tło `--dn-sukces-tlo`, ikona znacznika (fajka) | krok 1 „Warunki licencji”, meta „zaakceptowane” |
| `pracuje` | krok aktualnie widoczny; niesie `aria-current="step"` | obrys `--dn-sygnal-obrys`, tło `--dn-sygnal-tlo`, numer kroku | krok 3 „Składniki”, meta „krok bieżący” |
| `czeka` | krok jeszcze nieosiągnięty | obrys `--dn-obrys-mocny`, bez wypełnienia, numer kroku | kroki 4–6, meta „oczekuje” |
| `blad` *(zdefiniowany w arkuszu, nieużyty w zrzucie kroku 3)* | krok zakończony błędem | obrys `--dn-blad-obrys`, tło `--dn-blad-tlo`, ikona znacznika | — |

Każda pozycja toru niesie też pole meta (`.we-krok-meta`) — krótki tekst po prawej stronie nazwy
kroku. W zrzucie kroku 3 wartości meta są następujące:

| Krok | `data-stan` | Nazwa | Meta |
|---|---|---|---|
| 1 | `gotowy` | Warunki licencji | zaakceptowane |
| 2 | `gotowy` | Katalog instalacji | `C:\Program Files\Danaco Console` |
| 3 | `pracuje` | Składniki | krok bieżący |
| 4 | `czeka` | Skróty i uruchamianie | oczekuje |
| 5 | `czeka` | Instalacja | oczekuje |
| 6 | `czeka` | Zakończenie | oczekuje |

Krok 1 i krok 2 są jedynymi krokami, dla których prototyp zawiera dosłowną wartość zapisaną w
przebiegu (odpowiednio „zaakceptowane” i domyślny katalog Windows) — obie wartości pochodzą z pola
meta toru kroków, widocznego na zrzucie kroku 3, a nie z osobno wyrenderowanego ekranu kroków 1 i 2.
Treść własnych ekranów kroków 1, 2, 4, 5 i 6 (poza torem i paskiem działań, wspólnymi każdemu
krokowi — rozdz. 3.3) nie ma odrębnego zrzutu w prototypie; rozdz. 4.2, 4.3, 4.5, 4.6 i 4.7 opisują
tylko to, co wynika z toru kroków, z wykazu wymagań (rozdz. 2) i z ogólnej zasady zero blokad
(rozdz. 6). Treść przypisana w tych podrozdziałach kontrolkom nieujętym w zrzucie jest oznaczona
wprost jako **[DO DECYZJI OPERATORA]**, gdzie wykracza poza to, co tor kroków i wykaz wymagań
potwierdzają.

### 4.2 Krok 1 — warunki licencji

| Cecha | Wartość | Źródło |
|---|---|---|
| Nazwa kroku (tor) | „Warunki licencji” | tor kroków, zrzut kroku 3 |
| Meta po ukończeniu | „zaakceptowane” | tor kroków, zrzut kroku 3 |
| Treść | pełny tekst warunków licencji, zgodny z [Licencją produktu](../LICENSE.md) | **[DO DECYZJI OPERATORA]** — brak odrębnego zrzutu ekranu kroku 1; zgodność treści z licencją produktu jest wymogiem projektowym, nie zweryfikowaną wartością prototypu |
| Kontrola | pole wyboru „Warunki licencji zaakceptowane” | **[DO DECYZJI OPERATORA]** — brzmienie etykiety pola nie ma odrębnego zrzutu |
| Zachowanie bez akceptacji | przycisk „Dalej” pozostaje klikalny (zasada zero blokad, rozdz. 6); próba przejścia dalej bez akceptacji odpowiada komunikatem walidacji przy polu wyboru | rozdz. 6, zgodnie z `komponenty.css` w. 7 |

### 4.3 Krok 2 — katalog instalacji

| Cecha | Wartość | Źródło |
|---|---|---|
| Nazwa kroku (tor) | „Katalog instalacji” | tor kroków, zrzut kroku 3 |
| Meta po ukończeniu (Windows) | `C:\Program Files\Danaco Console` | tor kroków, zrzut kroku 3 — wartość dosłowna |
| Etykieta pola | „Katalog instalacji” | **[DO DECYZJI OPERATORA]** — nazwa kroku w torze potwierdza temat ekranu, nie dosłowną etykietę pola formularza |
| Zmiana ścieżki | kontrolka wyboru katalogu systemowego | **[DO DECYZJI OPERATORA]** — brak zrzutu ekranu kroku 2 |
| Walidacja | katalog musi być zapisywalny przez proces instalatora (rozdz. 6) | wynika z wymogu instalacji plików w tym katalogu |

Wartość domyślna katalogu na macOS i Linuksie nie ma odpowiednika w torze kroków — prototyp niesie
wyłącznie wartość Windows. **[DO DECYZJI OPERATORA]** — katalog domyślny na macOS (zwyczajowo
`/Applications/Danaco Console.app`) i na Ubuntu (zwyczajowo `/opt/danaco-console` albo katalog
profilu, zależnie od rozstrzygnięcia rozdz. 7.5).

### 4.4 Krok 3 — składniki

Jedyny krok z pełnym zrzutem ekranu w prototypie. Nadtytuł: „Instalacja · krok 3 z 6”. Tytuł:
„Składniki instalacji”. Wprowadzenie kroku, tekst dosłowny z `.we-lid`:

> „Wskaż składniki, które mają zostać zainstalowane na tym stanowisku. Każdy składnik opcjonalny
> można dodać albo usunąć po instalacji, w oknie konserwacji pakietu — wybór podjęty teraz nie jest
> wyborem ostatecznym.”

Treść właściwa kroku to wykaz sześciu składników z polami wyboru, podsumowanie zajętości dysku i
podgląd paska postępu — rozdz. 5 niesie pełny wykaz składników; rozdz. 4.6 wyjaśnia relację
podglądu paska postępu widocznego już na kroku 3 do właściwego kroku 5.

### 4.5 Krok 4 — skróty i uruchamianie

| Opcja | Domyślnie | Źródło |
|---|---|---|
| Skrót na pulpicie | włączony | **[DO DECYZJI OPERATORA]** — brak zrzutu ekranu kroku 4; wartość przyjęta jako zwyczajowa dla instalatorów Windows/macOS/Linux, niezweryfikowana w prototypie |
| Skrót w menu Start (Windows) / Launchpad (macOS) / menu aplikacji (Linux) | włączony | **[DO DECYZJI OPERATORA]**, jak wyżej |
| Uruchamianie przy starcie systemu | wyłączony | **[DO DECYZJI OPERATORA]**, jak wyżej |
| Integracja z powłoką systemu (skojarzenia plików, menu kontekstowe, odnośniki `danaco://`) | zależna od zaznaczenia składnika „Integracja z powłoką systemu” w kroku 3 (rozdz. 5) | zrzut kroku 3 — opis składnika |

Krok 4 nie ma odrębnego zrzutu ekranu w prototypie. Trzy pierwsze wiersze tabeli są oznaczone jako
**[DO DECYZJI OPERATORA]** zgodnie z rozdz. 4.1 — dokument nie przesądza dokładnego brzmienia
etykiet ani domyślnego stanu przełączników tego kroku bez zrzutu potwierdzającego. Wiersz czwarty
jest zweryfikowany pośrednio: opis składnika „Integracja z powłoką systemu” w kroku 3 wprost
wymienia „skojarzenie plików projektu, pozycja „Otwórz w Danaco Console” w menu kontekstowym oraz
obsługa odnośników `danaco://`” jako skutek zaznaczenia tego składnika — krok 4 jest miejscem
logicznym dla ustawień pochodnych tego wyboru, nie miejscem, gdzie wybór ten się odbywa.

### 4.6 Krok 5 — instalacja

Ekran postępu: pasek postępu kopiowania i rozpakowywania składników. Element paska postępu ma
odrębny zrzut — widoczny już na kroku 3, jako podgląd stanu przed rozpoczęciem: pasek `.dn-postep`
z wartością `0%`, wiersz etykiety „Postęp instalacji” i tekst stanu, dosłownie:

> „Kopiowanie plików rozpocznie się w kroku 5. Do tego czasu instalator nie zmienia niczego na
> dysku.”

To zdanie jest zweryfikowanym potwierdzeniem zasady: **instalator nie zapisuje nic na dysku przed
krokiem 5** — kroki 1–4 zbierają wybory Operatora bez skutku ubocznego na systemie plików. W tym
kroku instalator sprawdza wymagania (rozdz. 2) i — w systemie Windows — doinstalowuje brakujący
komponent WebView2, jeśli go nie ma (opis składnika „Powłoka natywna Tauri”, rozdz. 5).

| Element paska postępu | Klasa | Zachowanie |
|---|---|---|
| Tor paska | `.dn-postep-tor`, wysokość `4px`, tło `--dn-powierzchnia-2` | stały |
| Wypełnienie | `.dn-postep-wartosc`, tło `--dn-sygnal-wypelnienie` | szerokość rośnie z `0%` do `100%`, przejście `width` w czasie `--dn-czas-3` |
| Rola dostępności | `role="progressbar"` z `aria-valuenow`, `aria-valuemin="0"`, `aria-valuemax="100"`, `aria-label="Postęp instalacji"` | wartość liczbowa odczytywana przez technologie wspomagające |
| Etykieta liczbowa | wiersz `.in-postep-wiersz` — „Postęp instalacji” i wartość procentowa pogrubiona | aktualizowana wraz z paskiem |

Krok 5 nie udostępnia przycisku „Wstecz” — przerwanie w trakcie kopiowania jest jedyną drogą
odwrotu, opisaną w rozdz. 6.

### 4.7 Krok 6 — zakończenie

| Element | Treść | Źródło |
|---|---|---|
| Potwierdzenie | „Instalacja zakończona” | **[DO DECYZJI OPERATORA]** — brak zrzutu ekranu kroku 6; treść przyjęta analogicznie do wzorca `.we-alarm--sukces` innych okien wejściowych |
| Opcja | „Uruchom Danaco Console” (domyślnie zaznaczona) | **[DO DECYZJI OPERATORA]**, jak wyżej |
| Skutek zamknięcia z zaznaczoną opcją | przejście do okna startowego (etap 1 przepływu, [Przepływ okien](przeplyw-okien.md)) | wynika z rozdz. 1.1 — instalator jest etapem 0, okno startowe etapem 1 |
| Skutek zamknięcia bez zaznaczonej opcji | zamknięcie instalatora bez uruchomienia aplikacji | wynika z natury pola wyboru |

---

## 5. Składniki instalacji

Wartości dosłowne z kroku 3 prototypu — sześć składników, nie trzy: dwa wymagane, dwa zalecane
(domyślnie zaznaczone) i dwa opcjonalne (jeden domyślnie zaznaczony, jeden nie). Każdy składnik
opcjonalny albo zalecany można dodać albo usunąć po instalacji, w oknie konserwacji pakietu
(rozdz. 8).

| Składnik | Status | Zaznaczenie domyślne | Rozmiar | Opis dosłowny |
|---|---|---|---|---|
| Rdzeń platformy | wymagany | zaznaczony, nie do odznaczenia | 412 MB | „Cztery środowiska pracy, magistrala kontekstu, lokalna baza stanu i katalog konektorów. Bez rdzenia aplikacja nie działa — składnika nie da się odznaczyć.” |
| Powłoka natywna Tauri | wymagany | zaznaczony, nie do odznaczenia | 86 MB | „Okno aplikacji, integracja z systemem plików i mechanizm aktualizacji. W systemie Windows korzysta z komponentu WebView2; brakujący komponent zostanie doinstalowany.” |
| Silnik modeli lokalnych | opcjonalny | zaznaczony | 1 240 MB | „Uruchamianie modeli na stanowisku, bez wysyłania treści do usług zewnętrznych. Same wagi modeli nie wchodzą w skład pakietu — pobiera się je później, z poziomu ustawień.” |
| Pakiet czcionek | zalecany | zaznaczony | 38 MB | „Kroje interfejsu i krój o stałej szerokości znaku. Bez tego składnika aplikacja sięga po kroje systemowe, a układ widoków może się nieznacznie różnić od projektu.” |
| Integracja z powłoką systemu | zalecany | zaznaczony | 12 MB | „Skojarzenie plików projektu, pozycja „Otwórz w Danaco Console” w menu kontekstowym oraz obsługa odnośników danaco://. Wymaga uprawnień administratora.” |
| Moduł Mobile — parowanie | opcjonalny | **nie**zaznaczony | 54 MB | „Usługa parowania stanowiska z aplikacją mobilną w sieci lokalnej. Otwiera port nasłuchu tylko na czas parowania; można ją włączyć później w ustawieniach.” |

Suma rozmiarów sześciu składników przy zaznaczeniu domyślnym (wszystkie poza modułem Mobile):
412 + 86 + 1 240 + 38 + 12 = **1 788 MB**. Wartość wyświetlana w prototypie przy tym samym
zestawie zaznaczeń to **1 842 MB** — różnica 54 MB odpowiada dokładnie rozmiarowi składnika „Moduł
Mobile — parowanie”, co wskazuje, że w zrzucie prototypu pole wyboru tego składnika jest w istocie
zaznaczone mimo atrybutu `checked` nieobecnego w znaczniku HTML tego elementu; wartość podsumowania
(„Wybrane składniki zajmą na dysku 1 842 MB”) jest przywołana dosłownie z prototypu i stanowi wartość
wiążącą tej specyfikacji, nad domysłem wynikającym z odczytu atrybutu znacznika.

```
   Składniki instalacji — stan domyślny (krok 3)
   ┌────────────────────────────────────────────────────────────────────┐
   │ [✓] Rdzeń platformy              WYMAGANY          412 MB          │  ← nie do odznaczenia
   │ [✓] Powłoka natywna Tauri        WYMAGANY           86 MB          │  ← nie do odznaczenia
   │ [✓] Silnik modeli lokalnych      OPCJONALNY      1 240 MB          │  ← wybór Operatora
   │ [✓] Pakiet czcionek              ZALECANY           38 MB          │  ← wybór Operatora
   │ [✓] Integracja z powłoką systemu ZALECANY           12 MB          │  ← wybór Operatora
   │ [ ] Moduł Mobile — parowanie     OPCJONALNY         54 MB          │  ← wybór Operatora
   ├────────────────────────────────────────────────────────────────────┤
   │  Wybrane składniki zajmą na dysku 1 842 MB                         │
   │  Wolne miejsce na dysku C: 148,6 GB                                │
   └────────────────────────────────────────────────────────────────────┘
```

Legenda: `[✓]` składnik zaznaczony do instalacji; składniki wymagane są zaznaczone i zablokowane
wizualnie obrysem `--dn-obrys` i tłem `--dn-powierzchnia-2` (klasa `.in-skladnik--wymagany`), lecz
— zgodnie z zasadą zero blokad — pole wyboru nie niesie atrybutu blokady technicznej: reakcja na
próbę odznaczenia jest opisana w rozdz. 6, nie realizowana odjęciem klikalności.

Podsumowanie zajętości (`.in-suma`) niesie dwie wartości w jednym wierszu: sumę wybranych
składników pogrubioną monospace (`.in-suma b`) i wolne miejsce na dysku docelowym
(`.in-suma-wolne`) — w zrzucie prototypu „Wolne miejsce na dysku C: 148,6 GB”.

**Macierz zależności między składnikami.** Opisy dosłowne z tabeli powyżej niosą trzy zależności
funkcjonalne między składnikami, którymi krok 3 nie zarządza automatycznie — zaznaczenie jednego
składnika nie zaznacza ani nie wymusza drugiego, ale skutek jednego zależy od obecności drugiego.

| Składnik zależny | Zależność | Skutek braku zależności |
|---|---|---|
| Powłoka natywna Tauri | wymaga komponentu systemowego WebView2 (Windows) | krok 5 doinstalowuje WebView2 automatycznie, gdy go brak (rozdz. 4.6, 6) — nie jest to blokada instalacji, tylko krok dodatkowy |
| Silnik modeli lokalnych | wymaga osobno pobranych wag modeli, poza zakresem pakietu instalacyjnego | składnik instaluje wyłącznie silnik uruchomieniowy; bez pobrania wag z poziomu ustawień po instalacji funkcja pozostaje bezprzedmiotowa |
| Integracja z powłoką systemu | wymaga uprawnień administratora do wpisów obejmujących całą maszynę | przy instalacji w wariancie „dla użytkownika” (rozdz. 7.5) zakres integracji może być ograniczony do profilu Operatora — **[DO DECYZJI OPERATORA]**, czy krok 3 komunikuje to ograniczenie przy zaznaczeniu składnika w wariancie bez uprawnień administratora |

Żadna z trzech zależności nie jest zależnością blokującą w rozumieniu zasady zero blokad (rozdz. 6)
— każda rozwiązuje się działaniem dodatkowym (doinstalowanie, pobranie później, zawężenie zakresu),
nie odmową wykonania kroku.

---

## 6. Walidacje i komunikaty

Wzorzec instalatorów spoza rodziny Danaco Console rozwiązuje niegotowość formularza wyłączeniem
przycisku „Dalej” — przycisk staje się nieklikalny, dopóki warunek nie jest spełniony. Instalator
Danaco Console tego wzorca nie stosuje. Powodem nie jest estetyka, lecz spójność z resztą produktu:
każde inne okno platformy, opisane w [Katalogu komponentów](katalog-komponentow.md), traktuje
niegotowość jako stan komunikowany, nie jako stan odbierający możliwość działania — Operator zawsze
może spróbować, a system zawsze odpowiada, dlaczego próba się nie powiodła. Instalator, jako
pierwsze okno, z którym Operator się styka, ustanawia ten wzorzec od razu, zanim Operator zobaczy
resztę produktu. Poniższa tabela zbiera sytuacje wymagające walidacji na wszystkich sześciu krokach,
zachowanie zgodne z tą zasadą i dosłowny albo przyjęty komunikat.

Komponenty biblioteki `komponenty.css` obowiązują zasadą zapisaną w nagłówku arkusza (w. 7):
„Zasada zero blokad: żaden wariant nie odbiera klikalności — niegotowość komunikuje się opisem albo
komunikatem.” Instalator stosuje tę zasadę bez wyjątku: żaden przycisk instalatora nie jest
wyłączony atrybutem technicznym; niegotowość do przejścia dalej komunikuje wyłącznie stan błędu
kontrolki (`aria-invalid="true"`, obrys `--dn-blad-obrys`) i towarzyszący komunikat tekstowy.

| Sytuacja | Zachowanie zgodne z zasadą zero blokad | Komunikat | Rodzaj |
|---|---|---|---|
| Warunki licencji niezaakceptowane, próba „Dalej” | przycisk „Dalej” pozostaje klikalny; pole wyboru licencji przechodzi w stan `aria-invalid="true"`, ognisko wraca na pole | „Warunki licencji wymagają akceptacji przed dalszym krokiem” | błąd walidacji |
| Katalog bez prawa zapisu | przycisk „Dalej” pozostaje klikalny; pole katalogu przechodzi w stan błędu | „Brak prawa zapisu do wskazanego katalogu” | błąd walidacji |
| Za mało miejsca na dysku | przycisk „Dalej” na kroku 3 pozostaje klikalny; komunikat `.we-alarm--ostrzezenie` pojawia się nad podsumowaniem zajętości | „Za mało wolnego miejsca — wymagane 2,4 GB, dostępne mniej” | ostrzeżenie |
| System poniżej minimum | krok 5 nie rozpoczyna kopiowania; komunikat `.we-alarm` typu błąd zastępuje pasek postępu | „System nie spełnia wymagań minimalnych” | błąd |
| Mniej niż 16 GB RAM | ostrzeżenie informacyjne, bez wpływu na możliwość kontynuowania | „Zalecane 16 GB RAM przy modelach lokalnych” | ostrzeżenie |
| Brak komponentu WebView2 (Windows) | doinstalowanie w tle na kroku 5, pasek postępu pokazuje etap pośredni | „Instalowanie komponentu WebView2…” | potwierdzenie (informacja o toku pracy) |
| Próba odznaczenia składnika wymaganego | pole wyboru pozostaje klikalne technicznie; zaznaczenie wraca natychmiast do stanu zaznaczonego, obrys `.in-skladnik--wymagany` pozostaje | „Rdzeń platformy jest wymagany do działania aplikacji” (komunikat przy polu, nie odjęcie klikalności) | ostrzeżenie |
| Przerwanie instalacji w trakcie kroku 5 | przycisk „Przerwij” pozostaje dostępny na każdym kroku; potwierdzenie żąda drugiego kliknięcia | „Przerwać instalację? Skopiowane pliki zostaną usunięte.” | potwierdzenie |
| Zakończenie kroku 6 | ekran końcowy zastępuje pasek działań trójprzyciskowy jednym potwierdzeniem | „Instalacja zakończona” | potwierdzenie |

Komunikat instalatora korzysta z trzech wariantów `.we-alarm` zdefiniowanych w `wejscie.css`:
`--informacja` (obrys i tło `--dn-informacja-*`, przykład: informacja o koncie na kroku 3),
`--ostrzezenie` (obrys i tło `--dn-ostrzezenie-*`, przykład: brak miejsca na dysku) i `--sukces`
(obrys i tło `--dn-sukces-*`, przykład: potwierdzenie na kroku 6). Zrzut kroku 3 niesie zweryfikowany
przykład wariantu informacyjnego, dosłowna treść:

> „Instalacja nie obejmuje konta ani stanu pracy. Instalator umieszcza na stanowisku wyłącznie
> pliki programu. Konto Operatora, projekty, karty sesji i pamięć kontekstowa należą do platformy i
> zostaną odtworzone po zalogowaniu, przy pierwszym uruchomieniu aplikacji.”

---

## 7. Warianty pakietowania na trzech systemach operacyjnych

Rozdział zbiera ustalenia techniczne osadzone w prototypie `instalator.html` jako notatka projektowa
(sekcja „Czy instalator może tak wyglądać — możliwości i granice”) oraz w [Instalacji i
konfiguracji](../INSTALACJA-I-KONFIGURACJA.md) rozdz. 5.6, dotyczące tego, jak wygląd okna z rozdz.
3–6 daje się faktycznie zbudować na trzech systemach docelowych. Rozdział odpowiada na pytanie
techniczne wykonawcze — nie zmienia specyfikacji wyglądu okna ustalonej w rozdziałach poprzednich.

### 7.1 Windows — NSIS i WiX/MSI

Powłoka natywna produktu jest budowana narzędziem Tauri, które wytwarza pakiety Windows w dwóch
formatach: NSIS i WiX/MSI.

**NSIS** jest skryptowalny — pozwala podmienić grafikę, teksty, kolejność stron oraz dołożyć
własne strony kreatora; Tauri dopuszcza wskazanie własnego szablonu skryptu instalatora. Wygląd
pozostaje jednak wyglądem kontrolek Win32: własna typografia, własne odstępy i płynne animacje —
takie jak tętno kropki godła (rozdz. 3.2) czy przejście paska postępu w czasie `--dn-czas-3`
(rozdz. 4.6) — są tu trudne, a przy pełnej swobodzie układu dwukolumnowego (rozdz. 3.1) praktycznie
poza zasięgiem.

**WiX/MSI** to format deklaratywny, obsługiwany przez usługę Windows Installer. Zestaw stron
kreatora daje się zmieniać, lecz w ramach silnika dialogów MSI, który jest sztywny i ubogi.
Rozbudowany, „markowy” instalator MSI zwykle buduje się jako oddzielny program rozruchowy
(ang. bootstrapper — w świecie WiX służy do tego Burn), który pokazuje własne okno o dowolnym
wyglądzie, a właściwy pakiet MSI instaluje po cichu w tle. To jest droga techniczna zgodna z
wyglądem ustalonym w rozdziałach 3–6 niniejszego dokumentu.

[Instalacja i konfiguracja](../INSTALACJA-I-KONFIGURACJA.md) rozdz. 5.6 potwierdza stan
przygotowania: konfiguracja `tauri.conf.json` niesie już `"bundle": { "active": true, "targets":
["nsis"] }` wraz z metadanymi wydawcy, prawami autorskimi i kategorią „DeveloperTool”, uruchamiana
poleceniem `cargo tauri build` (bez przełącznika `--no-bundle`) w katalogu
`budowa/desktop/src-tauri`. Instalator NSIS wytworzony tą drogą bez dalszych zmian pakuje jednak
wyłącznie powłokę, nie rdzeń (rozdz. 1.3) — droga zgodna z sześcioma krokami niniejszej specyfikacji
wymaga rozszerzenia konfiguracji o dołączenie rdzenia jako zasobu pakietu, co jest kwestią otwartą
tego samego rozdziału źródłowego.

**Instalator własny.** Najbliżej wyglądu ustalonego w rozdz. 3 jest rozwiązanie, w którym instalator
to mała aplikacja natywna — technicznie ta sama powłoka co produkt — uruchamiana z pobranego pliku.
Wtedy okno instalatora jest po prostu widokiem aplikacji: te same żetony `--dn-*`, te same
komponenty `.dn-*`, ten sam układ dwukolumnowy `.we-*`. Kosztem jest utrzymanie: rozpakowanie
plików, wpisy do rejestru albo do bazy pakietów, dopisanie programu do listy „Aplikacje i funkcje”
oraz poprawny deinstalator (rozdz. 8) trzeba wykonać samodzielnie i samodzielnie przetestować.

### 7.2 macOS — DMG i pkg

Rozpowszechnianie na macOS odbywa się zwykle obrazem DMG z jednym gestem „przeciągnij do
Aplikacji” — kreatora wieloetapowego takiego jak sześć kroków niniejszej specyfikacji w ogóle się
tam nie oczekuje jako zwyczaju platformy. Format `.pkg` pozwala na wybór miejsca instalacji, wybór
składników i własne strony, ale wygląd stron pozostaje wyglądem systemowego Instalatora macOS.
Okno takie jak w rozdziałach 3–6 na macOS oznacza w praktyce: własną aplikację instalacyjną, a nie
przerobiony `.pkg`.

### 7.3 Linux — DEB, RPM, AppImage

Pakiety DEB, RPM i AppImage nie mają własnego interfejsu instalacyjnego — instaluje je menedżer
pakietów dystrybucji albo sklep z oprogramowaniem systemu. Graficzny kreator sześciu kroków byłby
tu obcy przyjętym zwyczajom platformy Linux; jeżeli ma powstać, to jako oddzielna aplikacja obok
pakietu systemowego, nie zamiast niego.

**Zestawienie porównawcze trzech platform.** Tabela zbiera w jednym miejscu rozstrzygnięcia
rozdz. 7.1–7.3 dla szybkiego porównania, bez powtarzania uzasadnień.

| Kryterium | Windows | macOS | Linux |
|---|---|---|---|
| Format natywny narzędzia Tauri | NSIS, WiX/MSI | DMG (domyślnie) | DEB, RPM, AppImage |
| Zgodność formatu natywnego z wyglądem rozdz. 3 | częściowa (NSIS skryptowalny) do żadnej (MSI) | żadna — DMG nie ma kreatora | żadna — pakiety nie mają interfejsu instalacyjnego |
| Droga zgodna z wyglądem rozdz. 3 | bootstrapper przed pakietem MSI, albo aplikacja instalacyjna własna | aplikacja instalacyjna własna | osobna aplikacja obok pakietu systemowego (nietypowe dla platformy) |
| Wymóg podpisu | certyfikat wydawcy (SmartScreen) | certyfikat + notaryzacja (Gatekeeper) | brak wymogu systemowego równoważnego |
| Stan przygotowania w produkcie (rozdz. 1.3) | konfiguracja NSIS obecna w `tauri.conf.json`, pakuje wyłącznie powłokę | brak konfiguracji opisanej w [Instalacji i konfiguracji](../INSTALACJA-I-KONFIGURACJA.md) | brak konfiguracji opisanej w tym samym dokumencie |

### 7.4 Podpis kodu

Bez podpisu instalator zderzy się z ostrzeżeniami systemu operacyjnego: na Windows z ekranem
SmartScreen, na macOS z Gatekeeperem, gdzie osobnym krokiem po podpisaniu jest notaryzacja u
producenta systemu. Podpis kosztuje — certyfikat wydawcy (rozdz. 1.2) i, w przypadku macOS,
uczestnictwo w programie deweloperskim — i wymaga miejsca na klucze w łańcuchu budowania. To
wydatek stały, niezależny od tego, jak instalator wygląda, i [Instalacja i konfiguracja](../INSTALACJA-I-KONFIGURACJA.md)
rozdz. 23.4, poz. 3, odnotowuje jako kwestię otwartą: „Podpisywanie artefaktów wydania i
wyłączenia z monitorowania czasu rzeczywistego”.

### 7.5 Instalacja dla użytkownika a dla wszystkich użytkowników

Instalacja w katalogu profilu użytkownika nie wymaga podniesienia uprawnień, przechodzi bez okna
kontroli konta użytkownika i ułatwia aktualizacje w tle. Instalacja dla wszystkich użytkowników
wymaga uprawnień administratora, za to daje jedną kopię programu na stanowisku. Te dwa warianty
mają różne skutki dla skojarzeń plików i menu kontekstowego (rozdz. 5, składnik „Integracja z
powłoką systemu”) — wpisy dotyczące całej maszyny wymagają uprawnień administratora, co opis tego
składnika w kroku 3 wprost odnotowuje. **[DO DECYZJI OPERATORA]** — czy krok 2 (rozdz. 4.3) niesie
jawny wybór wariantu instalacji, i jak zmiana wariantu po fakcie (zwykle: odinstalowanie i ponowna
instalacja) jest komunikowana Operatorowi.

### 7.6 Aktualizacje

Aktualizacje są osobnym mechanizmem, nie stroną instalatora — nie wchodzą w zakres sześciu kroków
rozdz. 4. Tauri ma wbudowany moduł aktualizacji, sprawdzający podpisany plik z opisem wydania i
podmieniający pakiet aplikacji. Aktualizacje różnicowe — przesyłanie samej różnicy między wydaniami
zamiast całego pakietu — są technicznie możliwe, lecz wymagają własnego zaplecza: budowania łatek,
wersjonowania i procedury awaryjnej, gdy łatka nie da się nałożyć. Przy wielkości pakietu rzędu
kilkuset megabajtów (rozdz. 5: suma domyślna 1 842 MB) rozsądniej zacząć od pełnej podmiany pakietu
i wrócić do różnic dopiero wtedy, gdy koszt pobierania stanie się realnym problemem. Szczegóły
modułu aktualizacji zmieniają się między wydaniami Tauri i wymagają sprawdzenia w bieżącej
dokumentacji przed decyzją wdrożeniową.

### 7.7 Czego nie warto obiecywać

| Obietnica ryzykowna | Dlaczego |
|---|---|
| Jednego, identycznego okna instalatora na trzech systemach | wymaga trzech osobnych ścieżek budowania i trzech osobnych testów (rozdz. 7.1–7.3) |
| Pełnej swobody typograficznej wewnątrz stron NSIS albo dialogów MSI | tam wygląd wyznaczają kontrolki systemowe (rozdz. 7.1) |
| Instalacji „bez żadnego okna systemowego”, gdy program pisze poza profilem użytkownika | podniesienie uprawnień zawsze pokaże okno systemu (rozdz. 7.5) |
| Postępu instalacji odmierzanego w procentach dokładnie (rozdz. 4.6) | silniki pakietów raportują postęp zgrubnie — pasek postępu jest informacją o toku pracy, nie pomiarem |

### 7.8 Zalecenie

Ekran taki jak w rozdziałach 3–6 ma sens jako **własna aplikacja instalacyjna** na Windows i macOS
(rozdz. 7.1–7.2), przy zachowaniu zwykłych pakietów systemowych na Linuksie (rozdz. 7.3) oraz przy
wdrożeniach masowych, gdzie administratorzy oczekują cichej instalacji z wiersza polecenia. Wtedy
jedno okno wejściowe niesie markę i sześć kroków niniejszej specyfikacji, a pakiet MSI albo DEB
pozostaje tym, czym ma być — nośnikiem plików.

---

## 8. Wycofanie, naprawa i aktualizacja

Po instalacji tym samym pakietem obsługuje się utrzymanie — okno konserwacji pakietu, osobne od
instalatora opisanego w rozdz. 3–6, lecz dzielące z nim rodzinę wyglądu okien wejściowych.

| Operacja | Miejsce | Skutek |
|---|---|---|
| Wycofanie instalacji | okno konserwacji pakietu | usunięcie aplikacji i składników opcjonalnych; dane Operatora usuwa osobny wybór, potwierdzony oddzielnie |
| Naprawa | okno konserwacji pakietu | ponowne rozpakowanie rdzenia i powłoki bez utraty konfiguracji |
| Dodanie/usunięcie składnika | okno konserwacji pakietu | zmiana zestawu składników opcjonalnych i zalecanych z rozdz. 5 (dołożenie modułu Mobile — parowanie) |
| Aktualizacja | mechanizm aktualizacji powłoki (rozdz. 7.6) | podmiana pakietu przy zachowaniu danych i konfiguracji |

```
   Po instalacji — okno konserwacji pakietu
   ┌────────────────────────────────────────────────┐
   │  ( Napraw )  ( Zmień składniki )  ( Odinstaluj)│
   └────────────────────────────────────────────────┘
```

Legenda: trzy operacje utrzymaniowe dostępne po instalacji z poziomu tego samego pakietu, jakim
zainstalowano produkt (Panel sterowania Windows, Ustawienia systemowe macOS albo menedżer pakietów
Linux, zależnie od wariantu z rozdz. 7); aktualizacja biegnie osobno, z wnętrza aplikacji, nie z
tego okna.

Zasięg wycofania zależy od wariantu instalacji ustalonego w kroku 2 (rozdz. 4.3, 7.5). Instalacja
w wariancie „dla użytkownika” wycofuje się bez podniesienia uprawnień — okno konserwacji pakietu
usuwa katalog profilu Operatora. Instalacja w wariancie „dla wszystkich użytkowników” żąda przy
wycofaniu tego samego podniesienia uprawnień, jakiego żądała instalacja — system operacyjny pokazuje
własne okno kontroli konta, poza zasięgiem wyglądu ustalonego w rozdz. 3. Dane Operatora (konto,
projekty, karty sesji, pamięć kontekstowa) nie są przedmiotem tej operacji niezależnie od wariantu —
komunikat kroku 3 (rozdz. 6, wariant informacyjny) wyjaśnia to już w trakcie instalacji: dane te
należą do platformy, nie do pakietu plików, i przetrwają wycofanie samego pakietu, jeżeli Operator
nie wybierze osobno ich usunięcia.

---

## 9. Etykiety interfejsu

Pełny wykaz dosłownych brzmień tekstowych okna instalatora, zweryfikowanych w
`design/05-okna/platformowe/instalator.html`.

| Element | Dosłowne brzmienie | Miejsce wystąpienia |
|---|---|---|
| Tytuł belki | „Instalacja programu Danaco Console” | `.in-belka-tytul` |
| Etykieta zamknięcia belki | „Zamknij instalator” | `aria-label` przycisku `.in-belka-zamknij` |
| Nazwa produktu | „Danaco Console” | `.we-marka-nazwa` |
| Motto | „AI Operating Environment” | `.we-marka-motto` |
| Etykieta wymagania 1 | „Wersja instalowana” | `.we-marka-poz b` |
| Treść wymagania 1 | „2.4.0 (kompilacja 2418) · pakiet podpisany cyfrowo” | `.we-marka-poz span` |
| Etykieta wymagania 2 | „System operacyjny” | `.we-marka-poz b` |
| Treść wymagania 2 | „Windows 10 w wersji 1809 lub nowszej, macOS 12, Ubuntu 22.04 LTS” | `.we-marka-poz span` |
| Etykieta wymagania 3 | „Procesor i pamięć” | `.we-marka-poz b` |
| Treść wymagania 3 | „4 rdzenie, 8 GB RAM (zalecane 16 GB przy modelach lokalnych)” | `.we-marka-poz span` |
| Etykieta wymagania 4 | „Dysk i sieć” | `.we-marka-poz b` |
| Treść wymagania 4 | „2,4 GB wolnego miejsca · połączenie z siecią przy pierwszym uruchomieniu” | `.we-marka-poz span` |
| Stopka kolumny tożsamości | „Danaco Core sp. z o.o. — wydawca oprogramowania. Pakiet podpisany certyfikatem wydawcy; suma kontrolna widoczna we właściwościach pliku.” | `.we-marka-stopka` (rozbieżność z konfiguracją opisana w rozdz. 1.2) |
| Nadtytuł kroku 3 | „Instalacja · krok 3 z 6” | `.we-nadtytul` |
| Tytuł kroku 3 | „Składniki instalacji” | `.we-tytul` |
| Wprowadzenie kroku 3 | „Wskaż składniki, które mają zostać zainstalowane na tym stanowisku…” (cytat pełny w rozdz. 4.4) | `.we-lid` |
| Nazwa kroku 1 (tor) | „Warunki licencji” | `.we-krok` |
| Meta kroku 1 | „zaakceptowane” | `.we-krok-meta` |
| Nazwa kroku 2 (tor) | „Katalog instalacji” | `.we-krok` |
| Meta kroku 2 | `C:\Program Files\Danaco Console` | `.we-krok-meta` |
| Nazwa kroku 3 (tor) | „Składniki” | `.we-krok` |
| Meta kroku 3 | „krok bieżący” | `.we-krok-meta` |
| Nazwa kroku 4 (tor) | „Skróty i uruchamianie” | `.we-krok` |
| Nazwa kroku 5 (tor) | „Instalacja” | `.we-krok` |
| Nazwa kroku 6 (tor) | „Zakończenie” | `.we-krok` |
| Meta kroków 4–6 | „oczekuje” | `.we-krok-meta` |
| Nazwa składnika 1 | „Rdzeń platformy” | `.in-skladnik-nazwa` |
| Znacznik składnika 1 | „wymagany” | `.in-znacznik--wymagany` |
| Opis składnika 1 | „Cztery środowiska pracy, magistrala kontekstu, lokalna baza stanu i katalog konektorów. Bez rdzenia aplikacja nie działa — składnika nie da się odznaczyć.” | `.in-skladnik-opis` |
| Nazwa składnika 2 | „Powłoka natywna Tauri” | `.in-skladnik-nazwa` |
| Opis składnika 2 | „Okno aplikacji, integracja z systemem plików i mechanizm aktualizacji. W systemie Windows korzysta z komponentu WebView2; brakujący komponent zostanie doinstalowany.” | `.in-skladnik-opis` |
| Nazwa składnika 3 | „Silnik modeli lokalnych” | `.in-skladnik-nazwa` |
| Znacznik składnika 3 | „opcjonalny” | `.in-znacznik` |
| Opis składnika 3 | „Uruchamianie modeli na stanowisku, bez wysyłania treści do usług zewnętrznych. Same wagi modeli nie wchodzą w skład pakietu — pobiera się je później, z poziomu ustawień.” | `.in-skladnik-opis` |
| Nazwa składnika 4 | „Pakiet czcionek” | `.in-skladnik-nazwa` |
| Znacznik składnika 4 | „zalecany” | `.in-znacznik` |
| Opis składnika 4 | „Kroje interfejsu i krój o stałej szerokości znaku. Bez tego składnika aplikacja sięga po kroje systemowe, a układ widoków może się nieznacznie różnić od projektu.” | `.in-skladnik-opis` |
| Nazwa składnika 5 | „Integracja z powłoką systemu” | `.in-skladnik-nazwa` |
| Znacznik składnika 5 | „zalecany” | `.in-znacznik` |
| Opis składnika 5 | „Skojarzenie plików projektu, pozycja „Otwórz w Danaco Console” w menu kontekstowym oraz obsługa odnośników danaco://. Wymaga uprawnień administratora.” | `.in-skladnik-opis` |
| Nazwa składnika 6 | „Moduł Mobile — parowanie” | `.in-skladnik-nazwa` |
| Znacznik składnika 6 | „opcjonalny” | `.in-znacznik` |
| Opis składnika 6 | „Usługa parowania stanowiska z aplikacją mobilną w sieci lokalnej. Otwiera port nasłuchu tylko na czas parowania; można ją włączyć później w ustawieniach.” | `.in-skladnik-opis` |
| Podsumowanie zajętości | „Wybrane składniki zajmą na dysku 1 842 MB” | `.in-suma` |
| Wolne miejsce | „Wolne miejsce na dysku C: 148,6 GB” | `.in-suma-wolne` |
| Etykieta postępu | „Postęp instalacji” | `.in-postep-wiersz` |
| Stan postępu | „Kopiowanie plików rozpocznie się w kroku 5. Do tego czasu instalator nie zmienia niczego na dysku.” | `.in-postep-stan` |
| Nagłówek komunikatu informacyjnego | „Instalacja nie obejmuje konta ani stanu pracy” | `.we-alarm--informacja b` |
| Treść komunikatu informacyjnego | „Instalator umieszcza na stanowisku wyłącznie pliki programu. Konto Operatora, projekty, karty sesji i pamięć kontekstowa należą do platformy i zostaną odtworzone po zalogowaniu, przy pierwszym uruchomieniu aplikacji.” | `.we-alarm--informacja span` |
| Przycisk 1 stopki | „Przerwij” | `.dn-btn--duch.dn-btn--sm` |
| Przycisk 2 stopki | „Wstecz” | `.dn-btn--zarys.dn-btn--sm` |
| Przycisk 3 stopki | „Dalej” | `.dn-btn--atrament.dn-btn--sm` |

---

## 10. Żetony i komponenty

| Miejsce zastosowania | Żeton `--dn-*` | Czego dotyczy |
|---|---|---|
| Belka tytułowa instalatora | `--dn-wym-belka` (48 px) | wysokość paska tytułu własnego, równa belce ramy aplikacji |
| Tło i tekst belki | `--dn-rama`, `--dn-rama-tekst`, `--dn-rama-tekst-2`, `--dn-rama-obrys` | kolory paska tytułu i kolumny tożsamości |
| Okno wejściowe | `--dn-r-xl`, `--dn-cien-3`, `--dn-obrys` | zaokrąglenie, cień i obrys całego okna `.we-okno` |
| Ikony belki i kolumny | `--dn-wym-ikona-sm` (16 px), `--dn-wym-ikona` (18 px) | rozmiary ikon w belce i w wykazie wymagań |
| Godło marki | rozmiar stały 44 px (poza skalą `--dn-wym-*`), `--dn-sygnal-500` dla kropki tętniącej | godło kolumny tożsamości |
| Tętno kropki | `--dn-czas-tetno`, `--dn-ease` | animacja `we-tetno`, wyłączana przy `prefers-reduced-motion` |
| Typografia nagłówków | `--dn-ff-naglowek`, `--dn-fw-polgruba`, `--dn-fs-2xl` (nazwa), `--dn-fs-xl` (tytuł kroku) | nazwa produktu i tytuł kroku |
| Typografia monospace | `--dn-ff-mono`, `--dn-ls-mono-wersaliki` | motto, nadtytuł kroku, meta toru kroków, rozmiary składników |
| Pole wyboru | `--dn-wym-check` (16 px), `--dn-r-xs`, `--dn-atrament`, `--dn-atrament-tekst` | `.dn-check` — pola składników i licencji |
| Przycisk główny | `--dn-atrament`, `--dn-atrament-hover`, `--dn-atrament-tekst` | „Dalej” (`.dn-btn--atrament`) |
| Przycisk drugorzędny | brak wypełnienia (`background: transparent`) | „Wstecz” (`.dn-btn--zarys`) |
| Przycisk trzeciorzędny | `--dn-tekst-2`, `--dn-hover` | „Przerwij” (`.dn-btn--duch`) |
| Rozmiar przycisków stopki | `--dn-wym-kontrolka` pomniejszony o `--dn-od-1`, `--dn-fs-sm` | `.dn-btn--sm` na wszystkich trzech przyciskach stopki |
| Znak toru kroku „gotowy” | `--dn-sukces-obrys`, `--dn-sukces-tlo`, `--dn-sukces-tekst` | krok ukończony |
| Znak toru kroku „pracuje” | `--dn-sygnal-obrys`, `--dn-sygnal-tlo`, `--dn-sygnal`, `--dn-powierzchnia-2` | krok bieżący |
| Znak toru kroku „blad” | `--dn-blad-obrys`, `--dn-blad-tlo`, `--dn-blad-tekst` | krok zakończony błędem (zdefiniowany, nieużyty w zrzucie) |
| Znacznik składnika wymaganego | `--dn-sygnal-obrys`, `--dn-sygnal-tlo`, `--dn-sygnal` | `.in-znacznik--wymagany` |
| Karta składnika | `--dn-obrys-subtelny`, `--dn-powierzchnia`, `--dn-powierzchnia-2` (wskazanie kursorem) | `.in-skladnik` |
| Pasek postępu | `--dn-powierzchnia-2` (tor), `--dn-sygnal-wypelnienie` (wypełnienie), `--dn-czas-3` (przejście) | `.dn-postep-tor`, `.dn-postep-wartosc` |
| Komunikat informacyjny | `--dn-informacja-obrys`, `--dn-informacja-tlo`, `--dn-informacja-tekst` | `.we-alarm--informacja` |
| Komunikat ostrzegawczy | `--dn-ostrzezenie-obrys`, `--dn-ostrzezenie-tlo`, `--dn-ostrzezenie-tekst` | `.we-alarm--ostrzezenie` |
| Komunikat sukcesu | `--dn-sukces-obrys`, `--dn-sukces-tlo`, `--dn-sukces-tekst` | `.we-alarm--sukces`, krok 6 |
| Pierścień ogniska | `--dn-fokus`, `--dn-wym-fokus` (2 px), `--dn-wym-fokus-odsuniecie` (2 px) | wszystkie kontrolki interaktywne |
| Odstępy sceny i panelu | `--dn-od-8`, `--dn-od-6`, `--dn-od-5`, `--dn-od-4`, `--dn-od-3`, `--dn-od-2`, `--dn-od-1` | dopełnienia na wszystkich poziomach zagnieżdżenia, malejące na progu 640 px |

| Klasa `.dn-*` / `.we-*` / `.in-*` | Rola | Modyfikatory |
|---|---|---|
| `.we-scena` | scena centrująca okno wejściowe | — |
| `.we-okno` | siatka dwukolumnowa okna wejściowego | lokalnie rozszerzona przez `.in-okno` o rząd belki |
| `.we-marka` | kolumna tożsamości | — |
| `.we-marka-godlo`, `.we-marka-nazwa`, `.we-marka-motto`, `.we-marka-lista`, `.we-marka-poz`, `.we-marka-stopka` | elementy kolumny tożsamości | — |
| `.we-panel` | kolumna kroku | — |
| `.we-glowa`, `.we-nadtytul`, `.we-tytul`, `.we-lid` | nagłówek kroku | — |
| `.we-kroki`, `.we-krok`, `.we-krok-znak`, `.we-krok-meta` | tor kroków | `[data-stan]` z wartościami `gotowy` / `pracuje` / `czeka` / `blad` |
| `.we-alarm` | komunikat w kolumnie kroku | `--informacja`, `--ostrzezenie`, `--sukces` |
| `.we-stopka` | pasek działań kroku | — |
| `.in-belka`, `.in-belka-znak`, `.in-belka-tytul`, `.in-belka-uchwyt`, `.in-belka-zamknij` | pasek tytułu własny instalatora | — |
| `.in-skladniki`, `.in-skladnik`, `.in-skladnik-nazwa`, `.in-skladnik-opis`, `.in-skladnik-rozmiar` | wykaz składników | `.in-skladnik--wymagany` |
| `.in-znacznik` | plakietka roli składnika | `.in-znacznik--wymagany` |
| `.in-suma`, `.in-suma-wolne` | podsumowanie zajętości dysku | — |
| `.in-postep`, `.in-postep-wiersz`, `.in-postep-stan` | opakowanie paska postępu | — |
| `.in-dzialania`, `.in-dzialania-odstep` | wiersz trzech przycisków stopki | — |
| `.dn-btn`, `.dn-btn--atrament`, `.dn-btn--zarys`, `.dn-btn--duch`, `.dn-btn--sm` | trzy przyciski stopki | rozdz. 11 |
| `.dn-check` | pole wyboru licencji i składników | — |
| `.dn-postep`, `.dn-postep-tor`, `.dn-postep-wartosc` | pasek postępu krok 5 | — |

---

## 11. Stany kontrolek

| Kontrolka | Spoczynek | Wskazanie kursorem | Wciśnięcie | Ognisko | Stan błędu / komunikat |
|---|---|---|---|---|---|
| Pole „Warunki licencji zaakceptowane” | `.dn-check` puste | obrys mocniejszy | zaznaczenie, ikona znacznika skaluje się z 0 do 1 | pierścień `--dn-fokus` | `aria-invalid="true"`, komunikat rozdz. 6 przy próbie „Dalej” bez akceptacji |
| Przycisk „Dalej” | `.dn-btn--atrament.dn-btn--sm`, tło `--dn-atrament` | tło `--dn-atrament-hover` | przesunięcie `translateY(1px)` | pierścień `--dn-fokus` | zawsze klikalny (zasada zero blokad); warunek niespełniony daje komunikat, nie blokadę |
| Przycisk „Wstecz” | `.dn-btn--zarys.dn-btn--sm`, tło przezroczyste | tło `--dn-hover` | przesunięcie `translateY(1px)` | pierścień `--dn-fokus` | niewidoczny na kroku 1 (brak kroku wcześniejszego) |
| Przycisk „Przerwij” | `.dn-btn--duch.dn-btn--sm`, kolor `--dn-tekst-2` | tło `--dn-hover`, kolor `--dn-tekst` | przesunięcie `translateY(1px)` | pierścień `--dn-fokus` | drugie kliknięcie żąda potwierdzenia (rozdz. 6) |
| Wybór katalogu (krok 2) | przycisk wyboru systemowego | tło `--dn-hover` | otwiera okno systemowe wyboru folderu | pierścień `--dn-fokus` | `aria-invalid="true"` przy katalogu bez prawa zapisu |
| Pole składnika wymaganego | zaznaczone, obrys `.in-skladnik--wymagany` | podpowiedź „wymagany” | zaznaczenie wraca natychmiast (rozdz. 6) | pierścień `--dn-fokus` | komunikat „Rdzeń platformy jest wymagany…” zamiast odjęcia klikalności |
| Pole składnika opcjonalnego/zalecanego | wg wyboru domyślnego (rozdz. 5) | obrys mocniejszy | przełączenie zaznaczenia, suma dysku (rozdz. 5) przelicza się natychmiast | pierścień `--dn-fokus` | — |
| Pasek postępu (krok 5) | wypełnienie `0%` | — | — | — | zastąpiony `.we-alarm` błędu przy systemie poniżej minimum |
| Znak toru kroku | wg `data-stan` (rozdz. 4.1) | — | — | — | wariant `blad` przy kroku zakończonym błędem |

Instalator nie zawiera elementów trwale wyłączonych bez powodu — nieaktywność wizualna (obrys
`.in-skladnik--wymagany`) zawsze towarzyszy zachowaniu klikalnemu i komunikatowi podanemu wprost,
zgodnie z zasadą zero blokad z nagłówka `komponenty.css`.

---

## 12. Przebieg nawigacji

Sekwencja od uruchomienia pliku pakietu do przejścia w okno startowe, z uwzględnieniem ścieżki
przerwania.

```
Operator                            Instalator                      System plików
   │                                     │                                 │
   │  uruchomienie pakietu               │                                 │
   ├────────────────────────────────────►│                                 │
   │                                     │  krok 1 — licencja              │
   │◄────────────────────────────────────┤                                 │
   │  akceptacja                         │                                 │
   ├────────────────────────────────────►│                                 │
   │                                     │  krok 2 — katalog               │
   │◄────────────────────────────────────┤                                 │
   │  potwierdzenie ścieżki              │                                 │
   ├────────────────────────────────────►│                                 │
   │                                     │  krok 3 — składniki (rozdz. 5)  │
   │◄────────────────────────────────────┤                                 │
   │  wybór składników                   │                                 │
   ├────────────────────────────────────►│                                 │
   │                                     │  krok 4 — skróty                │
   │◄────────────────────────────────────┤                                 │
   │  potwierdzenie opcji                │                                 │
   ├────────────────────────────────────►│                                 │
   │                                     │  krok 5 — sprawdzenie wymagań   │
   │                                     ├────────────────────────────────►│
   │                                     │        kopiowanie i rozpakowanie│
   │                                     │◄────────────────────────────────┤
   │  postęp na żywo                     │                                 │
   │◄────────────────────────────────────┤                                 │
   │                                     │  krok 6 — zakończenie           │
   │◄────────────────────────────────────┤                                 │
   │  „Uruchom Danaco Console” zaznaczone                                  │
   ├────────────────────────────────────►│                                 │
   │                                     │  zamknięcie instalatora         │
   │  okno startowe (etap 1)             │                                 │
   │◄────────────────────────────────────┤                                 │
```

Legenda: strzałki poziome — czynność Operatora albo odpowiedź instalatora; para strzałek w kroku 5 —
zapis na dysku, jedyny moment w całym przebiegu, w którym instalator zmienia zawartość systemu
plików (rozdz. 4.6). Ścieżka przerwania (przycisk „Przerwij”, dostępny na każdym kroku) nie jest
narysowana osobno w sekwencji powyżej — poniższy diagram uzupełnia dokładnie tę ścieżkę, zgodnie
ze scenariuszem D z rozdz. 15.

```
Operator                            Instalator                      System plików
   │                                     │                                 │
   │  klik „Przerwij”                    │                                 │
   ├────────────────────────────────────►│  (krok 5, postęp na wartości    │
   │                                     │   pośredniej — rozdz. 4.6)      │
   │  komunikat potwierdzenia            │                                 │
   │◄────────────────────────────────────┤                                 │
   │  drugie kliknięcie                  │                                 │
   │  potwierdzające                     │                                 │
   ├────────────────────────────────────►│                                 │
   │                                     │  usunięcie plików skopiowanych  │
   │                                     │  do tej chwili                  │
   │                                     ├────────────────────────────────►│
   │                                     │◄────────────────────────────────┤
   │  zamknięcie instalatora             │                                 │
   │◄────────────────────────────────────┤                                 │
   │  (bez przejścia do okna startowego — etap 1 nieosiągnięty)            │
```

Legenda: ścieżka przerwania kończy się zawsze zamknięciem instalatora bez przejścia dalej —
w odróżnieniu od ścieżki głównej (diagram powyżej), która kończy się otwarciem okna startowego.

---

## 13. Skróty klawiszowe

| Skrót | Działanie | Zasięg | Kolizje |
|---|---|---|---|
| `Enter` | aktywuje przycisk „Dalej” | cały ekran kroku | brak — instalator nie ma innych kontrolek reagujących na `Enter` |
| `Esc` | otwiera potwierdzenie przerwania instalacji (rozdz. 6), tożsame z kliknięciem „Przerwij” | cały ekran kroku | brak |
| `Tab` / `Shift+Tab` | przejście między kontrolkami ekranu w kolejności: tor kroków (nieinteraktywny) → treść właściwa kroku → „Przerwij” → „Wstecz” → „Dalej” | cały ekran kroku | brak — kolejność stała na wszystkich sześciu krokach, niezależna od liczby pól treści właściwej |
| `Spacja` | przełącza pole wyboru z ogniskiem (licencja, składnik) | pole z ogniskiem | brak — składnik wymagany wraca do stanu zaznaczonego zgodnie z rozdz. 6 |

Instalator nie definiuje skrótów właściwych aplikacji (takich jak skróty ramy opisane w
[Ramie okna](rama-okna.md)) — okno przedaplikacyjne nie ma menu ani wstążki, więc pula skrótów
ramy nie ma tu zastosowania. Kolizji między czterema skrótami instalatora a skrótami systemu
operacyjnego nie odnotowano: żaden z czterech nie nakłada się na skróty zastrzeżone przez Windows,
macOS ani powszechne środowiska graficzne Linuksa (przełączanie okien, zamykanie aplikacji,
zrzuty ekranu), sprawdzone wobec zestawień systemowych tych trzech platform.

---

## 14. Dostępność

| Wymóg | Realizacja |
|---|---|
| Nawigacja klawiaturą | pełna na każdym z sześciu ekranów, kolejność `Tab` opisana w rozdz. 13 |
| Ognisko widoczne | pierścień `--dn-fokus` o grubości `--dn-wym-fokus` (2 px) i odsunięciu `--dn-wym-fokus-odsuniecie` (2 px) na każdej kontrolce |
| Wskaźnik kroków | stan każdego kroku komunikowany tekstem meta (rozdz. 4.1: „zaakceptowane”, „krok bieżący”, „oczekuje”) i atrybutem `aria-current="step"` na kroku bieżącym, nie samym kolorem znaku |
| Postęp | pasek postępu z wartością liczbową czytaną przez `role="progressbar"`, `aria-valuenow`, `aria-valuemin`, `aria-valuemax`, `aria-label="Postęp instalacji"` (rozdz. 4.6) |
| Opis składnika powiązany z polem | `aria-describedby` łączy każde pole wyboru składnika (`.dn-check`) z jego opisem (`.in-skladnik-opis`), tak że technologia wspomagająca czyta opis razem z nazwą |
| Stan błędu | `aria-invalid="true"` na kontrolce, nigdy wyłącznie kolor obrysu (rozdz. 6, 11) |
| Zamknięcie okna | przycisk `.in-belka-zamknij` niesie `aria-label="Zamknij instalator"`, nie polega wyłącznie na ikonie |
| Ruch ograniczony | animacja tętna godła wyłączana przy `prefers-reduced-motion: reduce` (rozdz. 3.2) |
| Kontrast | zgodny z progami `budowa/shared/kontrasty-progi.json`, w obu motywach — okno wejściowe jest zbudowane wyłącznie na żetonach `--dn-*`, reagujących na motyw |
| Rola składnika bez opierania się na kolorze | znaczniki „wymagany”, „zalecany”, „opcjonalny” (rozdz. 5) są tekstem plakietki `.in-znacznik`, czytelnym przez czytnik ekranu niezależnie od koloru obrysu |
| Struktura nagłówków | `.we-tytul` (tytuł kroku) i nagłówek okna `Instalacja programu Danaco Console` tworzą hierarchię odczytywalną narzędziami nawigacji po nagłówkach, bez pomijania poziomów |
| Powiększenie tekstu | typografia w jednostkach względnych żetonów `--dn-fs-*`, bez wartości pikselowych zaszytych w treści — powiększenie ustawień systemowych skaluje tekst okna |

---

## 15. Scenariusze instalacji

Cztery przebiegi ilustrujące sześć kroków (rozdz. 4) w sytuacjach różniących się wyborem Operatora
i stanem stanowiska. Każdy scenariusz odwołuje się wyłącznie do wartości ustalonych w rozdziałach
poprzednich — żaden krok scenariusza nie wprowadza wartości spoza rozdz. 2–8.

**Scenariusz A — instalacja domyślna, system spełnia wymagania.**

1. Operator uruchamia pakiet podpisany wersją „2.4.0 (kompilacja 2418)” (rozdz. 1.2); okno
   instalatora otwiera się na kroku 1, tor kroków pokazuje sześć pozycji w stanie `czeka` poza
   pierwszą, w stanie `pracuje`.
2. Krok 1: Operator zaznacza pole „Warunki licencji zaakceptowane” i klika „Dalej”. Tor kroków
   zmienia pozycję 1 na `gotowy` z metą „zaakceptowane” (rozdz. 4.1–4.2).
3. Krok 2: Operator zostawia katalog domyślny `C:\Program Files\Danaco Console` (Windows) i klika
   „Dalej”. Walidacja zapisu (rozdz. 6) przechodzi bez komunikatu. Tor kroków zmienia pozycję 2 na
   `gotowy` z tą wartością jako metą (rozdz. 4.3).
4. Krok 3: ekran otwiera się z zestawem domyślnym — wszystkie sześć składników poza modułem
   Mobile zaznaczone (rozdz. 5), podsumowanie pokazuje „Wybrane składniki zajmą na dysku 1 842 MB”.
   Operator nie zmienia zestawu i klika „Dalej”.
5. Krok 4: Operator zostawia skróty domyślne (rozdz. 4.5) i klika „Dalej”.
6. Krok 5: instalator sprawdza wymagania (rozdz. 2) — system spełnia je w całości, łącznie z
   zaleceniem 16 GB RAM — i rozpoczyna kopiowanie. Pasek postępu (rozdz. 4.6) przechodzi od `0%` do
   `100%`; w trakcie tego kroku wykrywa brak komponentu WebView2 i doinstalowuje go automatycznie
   (rozdz. 6), bez zatrzymania paska.
7. Krok 6: ekran końcowy pokazuje „Instalacja zakończona” z zaznaczoną opcją „Uruchom Danaco
   Console” (rozdz. 4.7). Operator zamyka okno; aplikacja otwiera okno startowe (etap 1 przepływu).

**Scenariusz B — instalacja minimalna, bez składników opcjonalnych.**

1. Kroki 1–2 przebiegają jak w scenariuszu A.
2. Krok 3: Operator odznacza „Silnik modeli lokalnych”, „Pakiet czcionek”, „Integracja z powłoką
   systemu” i zostawia „Moduł Mobile — parowanie” niezaznaczony (stan domyślny, rozdz. 5). Próba
   odznaczenia „Rdzeń platformy” nie powiodłaby się zgodnie z zasadą zero blokad (rozdz. 6) — pole
   wraca do stanu zaznaczonego z komunikatem, ale Operator w tym scenariuszu takiej próby nie
   podejmuje. Podsumowanie zajętości przelicza się do 412 + 86 = **498 MB**.
3. Kroki 4–6 przebiegają jak w scenariuszu A, z tą różnicą, że krok 4 nie oferuje ustawień
   pochodnych integracji z powłoką systemu (rozdz. 4.5), ponieważ ten składnik nie został wybrany.
4. Wynik: instalacja zajmuje 498 MB zamiast 1 842 MB; Operator może dołożyć dowolny z czterech
   pominiętych składników później, z poziomu okna konserwacji pakietu (rozdz. 8).

**Scenariusz C — system poniżej wymagań minimalnych.**

1. Kroki 1–4 przebiegają bez przeszkód — żadna z wartości zebranych w krokach 1–4 nie zależy od
   parametrów sprzętowych stanowiska.
2. Krok 5: sprawdzenie wymagań (rozdz. 2) wykrywa procesor poniżej 4 rdzeni. Zgodnie z rozdz. 6
   jest to wymaganie twarde: krok 5 nie rozpoczyna kopiowania, pasek postępu zostaje zastąpiony
   komunikatem błędu „System nie spełnia wymagań minimalnych”. Przycisk „Wstecz” pozostaje
   dostępny — Operator może wrócić do kroku 3 i zmienić zestaw składników, choć w tym scenariuszu
   zmiana zestawu nie usunie przyczyny (niedobór rdzeni procesora nie jest funkcją wyboru
   składników).
3. Jedyną drogą naprzód w tym scenariuszu jest zamknięcie instalatora i instalacja na stanowisku
   spełniającym wymagania — dokument nie przewiduje instalacji częściowej poniżej minimum
   twardego.

**Scenariusz D — przerwanie w trakcie kopiowania.**

1. Kroki 1–4 przebiegają jak w scenariuszu A.
2. Krok 5: w trakcie kopiowania, przy pasku postępu na wartości pośredniej (40%), Operator
   klika „Przerwij” — przycisk dostępny na każdym z sześciu kroków (rozdz. 3.3, 13).
3. Instalator pokazuje komunikat potwierdzenia: „Przerwać instalację? Skopiowane pliki zostaną
   usunięte.” (rozdz. 6) — czynność nieodwracalna wymaga drugiego kliknięcia, zgodnie z ogólną
   zasadą potwierdzeń niszczących zawartość.
4. Po potwierdzeniu instalator usuwa pliki skopiowane do tej chwili z katalogu wskazanego w kroku 2
   i zamyka się bez przejścia do okna startowego. Stanowisko wraca do stanu sprzed uruchomienia
   pakietu — rozdz. 4.6 potwierdza, że przed krokiem 5 instalator nie zmienia niczego na dysku, więc
   przerwanie przed rozpoczęciem kopiowania nie zostawia żadnego śladu do usunięcia.
5. Ponowne uruchomienie pakietu zaczyna przebieg od kroku 1 — instalator nie zapamiętuje wyborów
   przerwanej sesji między uruchomieniami procesu instalującego.

Cztery scenariusze powyżej pokrywają łącznie wszystkie sytuacje opisane w rozdz. 6 (walidacje i
komunikaty) z wyjątkiem katalogu bez prawa zapisu, który zależy wyłącznie od uprawnień systemu
plików stanowiska i nie wymaga osobnego scenariusza narracyjnego ponad wiersz tabeli rozdz. 6:
reakcja jest identyczna niezależnie od tego, który krok poprzedza próbę przejścia dalej — pole
katalogu przechodzi w stan błędu, przycisk „Dalej” pozostaje klikalny, Operator poprawia ścieżkę i
ponawia próbę bez utraty wyborów zebranych na wcześniejszych krokach.

---

## 16. Kryteria odbioru

| Warunek | Sposób sprawdzenia |
|---|---|
| Instalator nie niesie ramy aplikacji | brak szyny nawigacji, wstążki i paska stanu na każdym z sześciu ekranów |
| Układ dwukolumnowy zgodny z `wejscie.css` | pomiar szerokości kolumny tożsamości (384 px) i całego okna (`min(1040px, 100%)`) przy szerokości ≥ 1000 px |
| Kolumna tożsamości znika poniżej 1000 px | zmiana szerokości okna przeglądarki/okna aplikacji w prototypie |
| Tor kroków niesie cztery stany zweryfikowane w arkuszu | odczyt `data-stan` na wszystkich sześciu pozycjach w różnych fazach przebiegu |
| Krok 1 nie blokuje technicznie przejścia bez akceptacji licencji, lecz komunikuje błąd | próba kliknięcia „Dalej” bez zaznaczenia zgody — przycisk pozostaje klikalny, pojawia się komunikat |
| Katalog domyślny to `C:\Program Files\Danaco Console` na Windows | otwarcie kroku 2 albo odczyt pola meta kroku 2 w torze |
| Sześć składników z rozdz. 5 obecnych w kroku 3, z rozmiarami zgodnymi z prototypem | porównanie z `design/05-okna/platformowe/instalator.html` |
| Suma zajętości domyślnej wynosi 1 842 MB | zaznaczenie domyślnego zestawu składników i odczyt podsumowania |
| Rdzeń platformy i powłoka natywna są nieodznaczalne na kroku 3, lecz pole pozostaje klikalne technicznie | próba odznaczenia obu składników wymaganych — zaznaczenie wraca, pojawia się komunikat |
| Krok 5 nie zapisuje nic na dysku przed jego rozpoczęciem | obserwacja systemu plików podczas kroków 1–4 |
| Krok 5 doinstalowuje WebView2, gdy go brak (Windows) | instalacja na systemie bez WebView2 |
| Po kroku 6 z zaznaczoną opcją aplikacja uruchamia się do okna startowego | zakończenie i obserwacja przejścia do etapu 1 przepływu |
| Okno konserwacji pakietu oferuje naprawę, zmianę składników i odinstalowanie | uruchomienie pakietu po instalacji |
| Każdy ekran przechodzi się klawiaturą, ognisko widoczne na każdej kontrolce | nawigacja `Tab` przez sześć kroków |
| Żadna kontrolka instalatora nie niesie atrybutu blokady technicznej | przegląd znaczników `.dn-btn`, `.dn-check` w prototypie i w kodzie — brak `disabled` |
| Rozbieżność nazwy wydawcy między prototypem a konfiguracją pakowania jest rozstrzygnięta | odczyt pola `publisher` w `tauri.conf.json` po decyzji Operatora (rozdz. 1.2) |
| Przerwanie instalacji w kroku 5 usuwa pliki skopiowane do tej chwili i nie przechodzi do okna startowego | wykonanie scenariusza D (rozdz. 15) i obserwacja katalogu docelowego po zamknięciu instalatora |
| Kolumna tożsamości zachowuje stałe tło `--dn-rama` niezależnie od motywu, kolumna kroku zmienia tło wraz z motywem | przełączenie motywu jasny/ciemny na dowolnym z sześciu kroków i porównanie z tabelą rozdz. 3.2 |
| Wymóg dyskowy 2,4 GB jest sprawdzany niezależnie od zestawu składników zaznaczonych w kroku 3 | zaznaczenie zestawu minimalnego (scenariusz B, rozdz. 15) i porównanie wartości sprawdzanej w kroku 5 z wykazem wymagań kolumny tożsamości |

---

*Koniec dokumentu. Okno instalatora — Specyfikacja docelowa, wersja 2.0, 2026-08-20.*

---
*Danaco Console — Platforma AI Workspace OS · v2.0 · status Deweloperski*
*© 2026 Danaco Holding Group Sp. z o.o. Wszelkie prawa zastrzeżone — Dariusz Naharnowicz.*
*Warunki korzystania: [Licencja produktu](../LICENSE.md). Kontakt: support@danaco-group.pl*
