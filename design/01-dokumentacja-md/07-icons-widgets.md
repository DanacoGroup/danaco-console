# Danaco Console — Ikonografia i widżety

| | |
|---|---|
| **Produkt** | **Danaco Console** — AI Operating Environment (warstwa wizualna v2.0) |
| **Rodzaj** | Opracowanie merytoryczno-techniczne — warstwa ikonografii i wzorców złożonych (widżetów) |
| **Producent** | Danaco Holding Group Sp. z o.o. |
| **Twórca** | Dariusz Naharnowicz |
| **Wersja** | v2.0 |
| **Status** | Deweloperski |
| **Data** | 2026-08-14 |
| **Odbiorcy** | DESIGNER (co rysować, jak nazywać, kiedy użyć) · DEWELOPER (co zbudować, z jakich klas `.dn-*`, jakie stany) · REDAKTOR DOKUMENTACJI (jak opisywać ikony i widżety) |
| **Zakres** | Część I — pełny zestaw **82 ikon** (`zasoby/ikony/`): zasady konstrukcji, pochodzenie, nazewnictwo, katalog pozycja po pozycji, siedem grup semantycznych, cztery emblematy środowisk, zasady użycia, dostępność, procedura dodawania ikony. Część II — **13 widżetów**, czyli złożonych, samodzielnych zestawów komponentów prezentujących stan systemu, wraz z anatomią, danymi, stanami, zachowaniem i składem z klas `.dn-*` |
| **Czego NIE zawiera** | Definicji żetonów (patrz `04-tokens.md`) · pełnej biblioteki komponentów prostych `.dn-*` (patrz `06-components.md`) · księgi znaku i logotypu (patrz `09-brand-system.md`) · makiet okien operacyjnych (katalog `05-okna/`) · plików źródłowych SVG (są w `zasoby/ikony/svg/`) · licencji Lucide w brzmieniu pełnym · ikon spoza zestawu — zestaw jest zamknięty i rozszerzany wyłącznie procedurą z rozdz. 9 |

---

## Spis treści

**Część I — Ikonografia**

1. [Zasady zestawu](#1-zasady-zestawu)
2. [Pochodzenie zestawu](#2-pochodzenie-zestawu)
3. [Nazewnictwo](#3-nazewnictwo)
4. [Pełny katalog ikon](#4-pełny-katalog-ikon)
5. [Grupy semantyczne](#5-grupy-semantyczne)
6. [Cztery emblematy środowisk](#6-cztery-emblematy-środowisk)
7. [Zasady użycia](#7-zasady-użycia)
8. [Dostępność ikon](#8-dostępność-ikon)
9. [Jak dodać ikonę](#9-jak-dodać-ikonę)

**Część II — Widżety**

10. [Definicja widżetu w Danaco Console](#10-definicja-widżetu-w-danaco-console)
11. [Katalog widżetów](#11-katalog-widżetów)
12. [Karty widżetów](#12-karty-widżetów)
13. [Decyzje projektowe](#13-decyzje-projektowe)

---

# CZĘŚĆ I — IKONOGRAFIA

## 1. Zasady zestawu

Ikona w Danaco Console jest **przyrządem odczytu**, nie ozdobą. Kontrakt kierunku (kierunek systemu projektowego) wiąże ją z językiem „instrumentu pomiarowego": jedna siatka, jedna grubość kreski, jedna barwa dziedziczona z kontekstu.

### 1.1. Sześć niezmienników

| # | Niezmiennik | Wartość | Uzasadnienie |
|---|---|---|---|
| 1 | **Siatka** | `24 × 24` (`viewBox="0 0 24 24"`) | Jedna siatka źródłowa dla wszystkich 82 pozycji — kreska nie łamie się przy skalowaniu; każde renderowanie to prosta zmiana `width`/`height` |
| 2 | **Obrys** | `stroke-width="1.75"` | Ujednolicony wobec domyślnych 2 px Lucide. 1,75 daje kreskę wyraźnie cieńszą — „precyzja instrumentu" — a jednocześnie nie znika przy 14 px, jak dzieje się przy 1,5 |
| 3 | **Wypełnienie** | `fill="none"` | Zestaw jest wyłącznie konturowy. Wyjątek: **kropka sygnału** w czterech emblematach środowisk oraz plamki w `paleta` — mają jawne `fill="currentColor"` i `stroke="none"` |
| 4 | **Barwa** | `stroke="currentColor"` | Ikona **nigdy** nie deklaruje własnej barwy. Dziedziczy `color` rodzica — dzięki temu ta sama ikona jest poprawna na przycisku duchu, w plakietce błędu i na atramentowym pasku górnym |
| 5 | **Zakończenia i łączenia** | `stroke-linecap="round"` · `stroke-linejoin="round"` | Zaokrąglenia łagodzą kreskę 1,75 i utrzymują czytelność narożników w małych stopniach |
| 6 | **Renderowanie** | `14 · 16 · 20 · 24 px` | Cztery stopnie z żetonów `--dn-wym-ikona-sm / -ikona / -ikona-lg / -ikona-xl`. Poza tymi czterema wartościami ikony **nie renderujemy** |

### 1.2. Cztery stopnie renderowania — kiedy który

| Stopień | Żeton | Gdzie stosowany |
|---|---|---|
| **14 px** | `--dn-wym-ikona-sm` | Medalion wpisu komunikacji (`.dn-wpis-medalion > svg`), znak kroku kolejki (12 px — patrz uwaga), drobne metadane |
| **16 px** | `--dn-wym-ikona` | **Wartość domyślna.** Przyciski (`.dn-btn`), przycisk ikonowy (`.dn-btn-ikona`), pozycje bocznej nawigacji, pole wyszukiwania, toast |
| **20 px** | `--dn-wym-ikona-lg` | Nagłówki paneli, kafel komponentu własnego, wyróżnione akcje w przyborniku |
| **24 px** | `--dn-wym-ikona-xl` | Pusty stan, karty ról, nagłówki dużych sekcji dokumentacyjnych |

> **Uwaga wykonawcza.** Dwa miejsca w `komponenty.css` schodzą poniżej skali: `.dn-plakietka > svg` (12 px) i `.dn-krok-znak > svg` (12 px). Są to ikony wewnątrz elementów wielkości 16–20 px, gdzie 14 px rozsadziłoby pojemnik. Traktujemy je jako **udokumentowany wyjątek dwóch komponentów**, a nie piąty stopień skali. Emblemat karty środowiska (`.dn-karta-srodowiska-godlo`, 40 px) jest emblematem, nie ikoną interfejsu — patrz rozdz. 6.

### 1.3. Anatomia pliku źródłowego

Każdy plik w `zasoby/ikony/svg/` ma tę samą głowę. Przykład: `dom.svg`.

```
<svg aria-hidden="true"
  xmlns="http://www.w3.org/2000/svg"
  width="24" height="24" viewBox="0 0 24 24"
  fill="none" stroke="currentColor"
  stroke-width="1.75" stroke-linecap="round" stroke-linejoin="round">
  <path d="M15 21v-8a1 1 0 0 0-1-1h-4a1 1 0 0 0-1 1v8" />
  <path d="M3 10a2 2 0 0 1 .709-1.528l7-6a2 2 0 0 1 2.582 0l7 6A2 2 0 0 1 21 10v9a2 2 0 0 1-2 2H5a2 2 0 0 1-2-2z" />
</svg>
```

`aria-hidden="true"` jest **wartością źródłową**: ikona domyślnie nie mówi. Odsłonięcie jej dla technologii wspomagających jest decyzją miejsca użycia, nie pliku (rozdz. 8).

### 1.4. Manifest jako rejestr

`zasoby/ikony/manifest.json` jest jedynym rejestrem zestawu. Struktura:

```
{
  "produkt": "Danaco Console",
  "warstwa": "zestaw ikon",
  "wersja": "v2.0",
  "data": "2026-08-11",
  "zasady": { siatka · obrys · wypelnienie · barwa · zakonczenia · laczenia · renderowanie },
  "zrodlo": { biblioteka: "Lucide (ISC), obrys ujednolicony do 1.75",
              wlasne:     "emblematy środowisk — siatka i kreska zestawu,
                           jedna wypełniona kropka sygnału" },
  "ikony": [ { nazwa, zrodlo, zastosowanie } × 82 ],
  "liczba-ikon": 82
}
```

Trzy pola każdej pozycji odpowiadają na trzy pytania: **jak się nazywa** (`nazwa`, ona sama jest nazwą pliku bez rozszerzenia), **skąd pochodzi** (`zrodlo`), **do czego służy** (`zastosowanie`). Pole `liczba-ikon` jest sumą kontrolną — po każdej zmianie musi zgadzać się z długością tablicy `ikony` i z liczbą plików w `svg/`.

---

## 2. Pochodzenie zestawu

### 2.1. Dwa źródła, jedna kreska

| Źródło | Liczba pozycji | Licencja | Traktowanie |
|---|---|---|---|
| **Lucide** | **78** | ISC | Przejęte bez zmian geometrii; zmieniony wyłącznie `stroke-width` (2 → 1,75) |
| **Własne emblematy domenowe** | **4** | Danaco Holding Group Sp. z o.o. | Narysowane na siatce 24×24 tą samą kreską 1,75; każdy zawiera dokładnie jedną wypełnioną kropkę sygnału |

Cztery pozycje własne to wyłącznie emblematy środowisk: `srodowisko-talkin`, `srodowisko-workspace`, `srodowisko-codestudio`, `srodowisko-multitaskingai`. W manifeście mają `"zrodlo": "wlasna"`.

### 2.2. Zasada „nie rysuj ręcznie, jeśli istnieje w bibliotece"

> **Ikon nie rysuje się ręcznie, jeżeli istnieją w bibliotece. Dorysowuje się wyłącznie brakujące pojęcia domenowe.** (kierunek systemu projektowego kontrakt systemu projektowego)

Reguła ma trzy skutki praktyczne, widoczne w samym zestawie:

| Pojęcie domenowe | Rozstrzygnięcie | Ikona |
|---|---|---|
| Rola **Coordinator** — plan i delegacja | Znaleziono odpowiednik w bibliotece: `lucide:waypoints` (węzły połączone trasą) | `wezly` |
| Rola **Executor** — realizacja zlecenia | Znaleziono odpowiednik: `lucide:circle-play` | `wykonawca` |
| Rola **Validator** — weryfikacja kroku | Znaleziono odpowiednik: `lucide:badge-check` | `walidator` |

Trzy role zespołu MultitaskingAI są pojęciami **wyłącznie domenowymi**, a mimo to nie zostały narysowane od zera — bo biblioteka miała kształt niosący to znaczenie. Rysowanie zaczęło się dopiero tam, gdzie żaden kształt Lucide nie mógł zastąpić tożsamości środowiska.

### 2.3. Co robimy z geometrią Lucide

| Wolno | Nie wolno |
|---|---|
| Zmienić `stroke-width` na 1,75 | Zmienić `viewBox` |
| Ustawić `stroke="currentColor"` | Dorysować element do kształtu bibliotecznego |
| Dodać `aria-hidden="true"` | Zmienić proporcje przez `transform: scale()` |
| Zapisać plik pod nazwą polską | Wypełnić kontur (`fill`) — poza emblematami |
| Renderować w 14/16/20/24 px | Obracać ikonę, by uzyskać wariant kierunkowy |

Ostatni punkt jest istotny: zestaw zawiera `grot-dol`, `grot-gora`, `grot-prawo` jako **trzy odrębne pliki**, a nie jeden obracany. Obrót przez CSS łamie `stroke-linejoin` na narożnikach i wprowadza niejednoznaczność w kodzie.

---

## 3. Nazewnictwo

### 3.1. Konwencja

**Nazwy ikon są polskie i opisowe.** Nazwa opisuje **rzecz albo czynność**, nie bibliotekę pochodzenia i nie moduł, w którym ikona akurat wystąpiła.

| Zasada | Przykład zgodny | Przykład niezgodny |
|---|---|---|
| Polski, małe litery, bez znaków diakrytycznych | `ostrzezenie`, `klodka`, `wezly` | `ostrzeżenie`, `warning`, `alertTriangle` |
| Rozdzielenie członów łącznikiem | `grot-dol`, `link-zewnetrzny`, `tabela-danych` | `grotDol`, `link_zewnetrzny` |
| Rzeczownik dla obiektu, rzeczownik odczasownikowy lub tryb rozkazujący dla czynności | `dokument`, `folder` · `uruchom`, `wyslij`, `odswiez` | `plikDokumentu`, `wysylanie` |
| Opis **kształtu i pojęcia**, nie miejsca użycia | `dzwonek` (nie `powiadomienia`) | `ikona-paska-gornego` |
| Prefiks `srodowisko-` wyłącznie dla czterech emblematów | `srodowisko-talkin` | `talkin`, `emblemat-1` |

### 3.2. Nazwa pliku = nazwa w manifeście

Nazwa w polu `nazwa` manifestu jest **dokładnie** nazwą pliku bez rozszerzenia:

```
manifest.json  →  "nazwa": "grot-prawo"
katalog svg/   →  grot-prawo.svg
odwołanie      →  zasoby/ikony/svg/grot-prawo.svg
```

Nie ma mapowania pośredniego, tablicy aliasów ani nazw kodowych. Jeżeli nazwa się zmienia — zmienia się plik, wpis w manifeście i wszystkie odwołania jednocześnie.

### 3.3. Trzy pary nazw wymagające czujności

| Para | Różnica |
|---|---|
| `agent` vs `agenci` | `agent` (`lucide:bot`) = **jedna** jednostka inteligencji. `agenci` (`lucide:users`) = **zespół**, moduł Agents, role |
| `blad` vs `ostrzezenie` | `blad` (`lucide:circle-x`) = **niepowodzenie operacji**. `ostrzezenie` (`lucide:triangle-alert`) = **skutki wymagające uwagi**, operacja mogła się powieść |
| `ptaszek` vs `ptaszek-kolo` | `ptaszek` (`lucide:check`) = **zaznaczenie**, potwierdzenie wykonania w polu wyboru. `ptaszek-kolo` (`lucide:circle-check`) = **stan** „proces zakończony poprawnie" |

### 3.4. Ikona nieoczywista — `karta-okna`

`karta-okna` (`lucide:app-window`) niesie trzy znaczenia naraz: okno platformy, kartę przeglądarki i moduł **Browser**. Nazwa opisuje kształt (rama okna z paskiem tytułu), a nie żadne z trzech zastosowań — to celowe, bo zastosowania są trzy, a plik jeden.

---

## 4. Pełny katalog ikon

82 pozycje w kolejności manifestu. Kolumna **Grupa** odsyła do rozdz. 5. Kolumna **Moduł / okno** wskazuje udokumentowane miejsce wystąpienia; gdy ikona jest ogólnosystemowa, wpisano zakres, w jakim występuje.

| # | Nazwa | Źródło | Grupa | Zastosowanie (manifest) | Moduł / okno |
|---:|---|---|---|---|---|
| 1 | `dom` | `lucide:house` | nawigacja | Centrum dowodzenia — powrót do strony głównej | Pasek górny (rama, wszystkie 4 powłoki) · Strona główna |
| 2 | `menu` | `lucide:menu` | nawigacja | Rozwinięcie nawigacji bocznej, menu podręczne | Powłoki TalkIn · WorkSpace · CodeStudio · MultitaskingAI |
| 3 | `strzalka-lewo` | `lucide:arrow-left` | nawigacja | Cofnięcie, powrót | Okno rejestracji i logowania (Odzyskiwanie konta krok 1/2) · Browser Window |
| 4 | `strzalka-prawo` | `lucide:arrow-right` | nawigacja | Przejście dalej, wejście w pozycję | Browser Window · Panel orkiestracji, sekcja Kolejki |
| 5 | `grot-dol` | `lucide:chevron-down` | nawigacja | Rozwinięcie sekcji, lista rozwijana | Okno Konfiguracji (13 zakresów) · każde pole `.dn-wybor` |
| 6 | `grot-gora` | `lucide:chevron-up` | nawigacja | Zwinięcie sekcji | Okno Konfiguracji · Subagent Network (zwijanie listy) |
| 7 | `grot-prawo` | `lucide:chevron-right` | nawigacja | Wejście w głąb, okruszki nawigacji | Library Explorer · Project Tree · Sources Manager |
| 8 | `wiecej` | `lucide:ellipsis` | nawigacja | Menu dodatkowych czynności pozycji | Nagłówek każdego okna operacyjnego (menu „⋯") |
| 9 | `link-zewnetrzny` | `lucide:external-link` | nawigacja | Odesłanie poza bieżący widok | Panel orkiestracji, sekcja Orkiestracja · Sources Panel |
| 10 | `szukaj` | `lucide:search` | nawigacja | Wyszukiwanie — pasek górny, Command Center | Pasek górny (`.dn-pasek-szukaj`) · Command Center |
| 11 | `filtr` | `lucide:funnel` | nawigacja | Zawężenie wykazu, panel warunków | Panel orkiestracji, sekcja Kolejki · Queue Manager · Logs Viewer |
| 12 | `zamknij` | `lucide:x` | nawigacja | Zamknięcie karty, modala, powiadomienia | Pas kart sesji (`✕`) · `.dn-modal` · `.dn-toast` |
| 13 | `karta-okna` | `lucide:app-window` | nawigacja | Okno, karta przeglądarki, moduł Browser | Moduł **Browser** — Browser Window, Sources Panel |
| 14 | `plus` | `lucide:plus` | działanie | Utworzenie bytu: sesja, okno, kanał | Pas kart sesji (`+`) · Agent Builder · Tags & Collections |
| 15 | `olowek` | `lucide:pencil` | działanie | Zmiana nazwy, edycja | Sources Manager · Library Explorer · Assets Panel |
| 16 | `kosz` | `lucide:trash-2` | działanie | Usunięcie pozycji (z potwierdzeniem) | Sources Manager · Library Explorer · Assets Panel · Agent Manager |
| 17 | `pobierz` | `lucide:download` | działanie | Pobranie pliku, eksport wyniku | Export Panel (Research) · Report Builder · Build Output |
| 18 | `wgraj` | `lucide:upload` | działanie | Wgranie pliku, import | Library Explorer · Assets Panel · Sources Manager |
| 19 | `wyslij` | `lucide:send-horizontal` | działanie | Wysłanie komunikatu w oknie komunikacji | **Chat Window** — pasek promptu, wszystkie 15 modułów |
| 20 | `odswiez` | `lucide:refresh-cw` | działanie | Ponowienie próby, odświeżenie danych | Okno startowe (stan Błąd połączenia) · Queue Manager (retry) |
| 21 | `odpowiedz` | `lucide:reply` | działanie | Odpowiedź, cytowanie wpisu | **Chat Window** — kontrolka przy wpisie historii |
| 22 | `uruchom` | `lucide:play` | działanie | Uruchomienie procesu, zadania, automatyzacji | Execution Monitor · Workflow Builder · Process Monitor |
| 23 | `zatrzymaj` | `lucide:square` | działanie | Zatrzymanie procesu | Execution Monitor · Process Monitor · Coordinator Chat (stop) |
| 24 | `wstrzymaj` | `lucide:pause` | działanie | Wstrzymanie kolejki, pauza | Queue Manager · Coordinator Chat (pauza) |
| 25 | `spinacz` | `lucide:paperclip` | działanie | Załącznik komunikatu | **Chat Window** — przybornik paska promptu |
| 26 | `ustawienia` | `lucide:settings` | działanie | Konfiguracja — okna, kanału, środowiska | Pasek górny (szybka konfiguracja sesji) · Strefa 3 · Okno Konfiguracji |
| 27 | `oko` | `lucide:eye` | działanie | Podgląd, odsłonięcie wartości ukrytej | Panel orkiestracji, sekcja Monitor procesu · Preview Window · pole hasła |
| 28 | `gwiazdka` | `lucide:star` | działanie | Wyróżnienie pozycji, presety zespołów | Panel orkiestracji, sekcja Zespoły · Library Explorer |
| 29 | `kopiuj` | `lucide:copy` | działanie | Kopiowanie treści, identyfikatora | Output Console · Logs Viewer · Chat Window |
| 30 | `ptaszek` | `lucide:check` | stan | Potwierdzenie wykonania, zaznaczenie | `.dn-check` · Skills Manager · Permissions Center |
| 31 | `ptaszek-kolo` | `lucide:circle-check` | stan | Powodzenie — proces zakończony poprawnie | `.dn-krok--poprawny` · `.dn-toast--sukces` · Monitor procesu |
| 32 | `blad` | `lucide:circle-x` | stan | Błąd — niepowodzenie operacji | `.dn-toast--blad` · Errors Panel · Diagnostics Center |
| 33 | `ostrzezenie` | `lucide:triangle-alert` | stan | Ostrzeżenie — skutki wymagające uwagi | Okno startowe (błąd połączenia) · `.dn-krok--bledy` · `.dn-toast--ostrzezenie` |
| 34 | `info` | `lucide:info` | stan | Wyjaśnienie, podpowiedź kontekstowa | Objaśnienie `[?]` przy każdej pozycji Okna Konfiguracji |
| 35 | `klodka` | `lucide:lock` | stan | Treść chroniona, dane dostępowe | Okno rejestracji i logowania · Okno Ustawień (uwierzytelnianie) |
| 36 | `tarcza` | `lucide:shield` | stan | Uprawnienia, zakres dostępu, izolacja | Permissions Center · Okno punktów izolacji · Panel prowenancji |
| 37 | `zegar` | `lucide:clock` | stan | Oczekiwanie, znacznik czasu | Panel orkiestracji, sekcja Harmonogram i automatyki · Scheduler |
| 38 | `historia` | `lucide:history` | stan | Historia sesji, wersje | Session Repository · Versioning Panel · Okno Konfiguracji (Historia) |
| 39 | `dzwonek` | `lucide:bell` | stan | Powiadomienia | Pasek górny · Strefa 3 (Always On Display) · Okno Ustawień |
| 40 | `aktywnosc` | `lucide:activity` | stan | Praca w tle, telemetria procesu | Monitor procesu · Actions Monitor · Execution Monitor |
| 41 | `dokument` | `lucide:file-text` | obiekt | Dokument, opracowanie — moduł Studio | Moduł **Studio** — Studio Editor, Session Repository |
| 42 | `plik` | `lucide:file` | obiekt | Plik o nieokreślonym rodzaju | File Preview · Project Tree · Library Explorer |
| 43 | `folder` | `lucide:folder` | obiekt | Katalog roboczy, grupa zasobów | Project Tree · Library Explorer · Tags & Collections |
| 44 | `archiwum` | `lucide:archive` | obiekt | Zasoby odłożone, Session Repository | **Session Repository** (Studio) · Library Explorer |
| 45 | `obraz` | `lucide:image` | obiekt | Załącznik graficzny, moduł Design (zasoby) | Moduł **Design** — Assets Panel · Chat Window (załącznik) |
| 46 | `kod` | `lucide:code-xml` | obiekt | Fragment kodu, moduł Developer | Moduł **Developer** — Code Editor · `.dn-kod` w Chat Window |
| 47 | `koperta` | `lucide:mail` | obiekt | Wiadomość, korespondencja | Okno rejestracji i logowania (e-mail) · Odzyskiwanie konta |
| 48 | `kalendarz` | `lucide:calendar` | obiekt | Harmonogram, automatyzacja cykliczna | **Scheduler** (Automations) · Panel orkiestracji, sekcja Harmonogram |
| 49 | `tabela-danych` | `lucide:table` | obiekt | Zestawienie, dane tabelaryczne | Queue Manager · Orchestrator · Monitor procesu |
| 50 | `biblioteka` | `lucide:library-big` | moduł | Moduł Library — repozytorium wiedzy | Moduł **Library** — Library Explorer, Project Library |
| 51 | `baza` | `lucide:database` | obiekt | Trwałość, magazyn danych | Context Memory · Okno Konfiguracji (Pamięć) |
| 52 | `wykres` | `lucide:chart-line` | moduł | Metryki, moduł Research (wyniki) | Moduł **Research** — Findings Panel, Report Builder |
| 53 | `terminal` | `lucide:square-terminal` | moduł | Moduł Terminal — powłoki wykonawcze | Moduł **Terminal** — Terminal Tabs, Output Console |
| 54 | `galaz` | `lucide:git-branch` | obiekt | Repozytorium, moduł Developer | **Git Panel** (Developer) · Project Tree |
| 55 | `agent` | `lucide:bot` | rola | Agent — jednostka inteligencji | Agent Builder · `.dn-wpis--inteligencja` · karty ról |
| 56 | `agenci` | `lucide:users` | moduł | Zespół, moduł Agents / role | Moduł **Agents** · Panel orkiestracji, sekcja Role · Agent Manager |
| 57 | `uzytkownik` | `lucide:user` | rola | Operator, konto, tożsamość | Pasek górny (profil) · `.dn-wpis--czlowiek` · Okno Ustawień |
| 58 | `siec` | `lucide:network` | moduł | Orkiestracja, topologia — MultitaskingAI | Powłoka **MultitaskingAI** · Orchestrator · Subagent Network |
| 59 | `wezly` | `lucide:waypoints` | rola | Koordynator — plan i delegacja | **Coordinator** · Coordinator Chat |
| 60 | `wykonawca` | `lucide:circle-play` | rola | Wykonawca — realizacja zlecenia | **Executor 1** · **Executor 2** · Executor Chat |
| 61 | `walidator` | `lucide:badge-check` | rola | Walidator — weryfikacja kroku | **Executor 3 / Validator** · Results Analyzer |
| 62 | `warstwy` | `lucide:layers` | moduł | Moduł Workspace — przestrzeń projektowa | Moduł **Workspace** — Project Dashboard, Context Memory |
| 63 | `automatyzacja` | `lucide:workflow` | moduł | Moduł Automations — procesy i kolejki | Moduł **Automations** — Workflow Builder, Queue Manager |
| 64 | `debata` | `lucide:messages-square` | moduł | Moduł Roundtable — debata modeli | Moduł **Roundtable** — Model Panels, Debate Panel |
| 65 | `tlumacz` | `lucide:languages` | moduł | Moduł Translate | Moduł **Translate** — Translation Panels, Glossary Manager |
| 66 | `badanie` | `lucide:telescope` | moduł | Moduł Research — analizy i raporty | Moduł **Research** — Research Workspace, Sources Manager |
| 67 | `aplikacje` | `lucide:layout-grid` | moduł | Moduł Apps — budowa produktów | Moduł **Apps** — Product Builder, Architecture Designer |
| 68 | `paleta` | `lucide:palette` | moduł | Moduł Design — materiały wizualne | Moduł **Design** — Design Board, Prompt Builder |
| 69 | `rozmowa` | `lucide:message-square` | moduł | Okno komunikacji — pas komunikacji | **Chat Window** — wspólne dla wszystkich 15 modułów |
| 70 | `mikrofon` | `lucide:mic` | moduł | Moduł Assistant — interfejs głosowy | Moduł **Assistant** — Voice Console, Activity Feed |
| 71 | `diagnostyka` | `lucide:stethoscope` | moduł | Moduł Diagnostics — analiza problemów | Moduł **Diagnostics** — Diagnostics Center, Recommendations Panel |
| 72 | `monitor` | `lucide:monitor` | obiekt | Stanowisko, środowisko wykonania | Okno Ustawień (urządzenia i parowanie) · Process Monitor |
| 73 | `telefon` | `lucide:smartphone` | obiekt | Widok mobilny — pozycja „Mobile" | Strefa 3 (**Mobile**) · pasek górny · Okno Ustawień |
| 74 | `cpu` | `lucide:cpu` | obiekt | Zasoby wykonawcze, serwer rdzenia | Process Monitor · Diagnostics Center · Deployment Panel |
| 75 | `polecenie` | `lucide:command` | nawigacja | Command Center — paleta poleceń | **Command Center** (warstwa 1300) |
| 76 | `slonce` | `lucide:sun` | działanie | Przełącznik motywu — jasny | Pasek górny (`data-przelacz-motyw`) · Okno Ustawień (wygląd) |
| 77 | `ksiezyc` | `lucide:moon` | działanie | Przełącznik motywu — ciemny | Pasek górny (`data-przelacz-motyw`) · Okno Ustawień (wygląd) |
| 78 | `globus` | `lucide:globe` | obiekt | Sieć, źródło zewnętrzne, danaco-web | Browser Window · Sources Manager · Connectors Manager |
| 79 | `srodowisko-talkin` | **własna** | środowisko | Środowisko TalkIn — Myśl. Analizuj. Rozumiej. | Karta środowiska (Strefa 1) · nagłówek powłoki **TalkIn** |
| 80 | `srodowisko-workspace` | **własna** | środowisko | Środowisko WorkSpace — Planuj. Organizuj. Realizuj. | Karta środowiska (Strefa 1) · nagłówek powłoki **WorkSpace** |
| 81 | `srodowisko-codestudio` | **własna** | środowisko | Środowisko CodeStudio — Projektuj. Buduj. Rozwijaj. | Karta środowiska (Strefa 1) · nagłówek powłoki **CodeStudio** |
| 82 | `srodowisko-multitaskingai` | **własna** | środowisko | Środowisko MultitaskingAI — Deleguj. Koordynuj. Nadzoruj. | Karta środowiska (Strefa 1) · nagłówek powłoki **MultitaskingAI** |

**Suma kontrolna:** 82 pozycje manifestu = 82 pliki w `zasoby/ikony/svg/` = 82 wiersze powyżej.

---

## 5. Grupy semantyczne

Grupy **nie występują w manifeście** — są warstwą porządkującą wprowadzoną w tym opracowaniu (patrz rozdz. 13, decyzje projektowe opracowania). Służą do dwóch rzeczy: filtrowania galerii i rozstrzygania, czy nowa ikona ma prawo powstać.

### 5.1. Siedem grup

```
                      ZESTAW 82 IKON
                             │
   ┌───────────┬─────────────┼─────────────┬───────────┐
   ▼           ▼             ▼             ▼           ▼
NAWIGACJA  DZIAŁANIE       STAN         OBIEKT      MODUŁ
   14         18            11            15          15
                             │
                   ┌─────────┴─────────┐
                   ▼                   ▼
              ŚRODOWISKO             ROLA
                   4                   5
```

| Grupa | Liczba | Odpowiada na pytanie | Barwa w użyciu | Reguła |
|---|---:|---|---|---|
| **nawigacja** | 14 | „dokąd się przemieszczam?" | `--dn-tekst-2`, na hover `--dn-tekst` | Nigdy nie niesie barwy stanu |
| **działanie** | 18 | „co zrobię?" | dziedziczy z przycisku | Zawsze wewnątrz elementu klikalnego |
| **stan** | 11 | „co się dzieje / co się stało?" | `--dn-sukces-tekst`, `--dn-ostrzezenie-tekst`, `--dn-blad-tekst`, `--dn-informacja-tekst` | **Obowiązkowo z etykietą** — stan nigdy samym kolorem |
| **obiekt** | 15 | „czym to jest?" | `--dn-tekst-2` lub `--dn-tekst-3` | Neutralna; nie sygnalizuje stanu obiektu |
| **moduł** | 15 | „w którym module jestem?" | `--dn-tekst-2`, aktywna pozycja `--dn-tekst` | Jedna ikona = jeden moduł; przypisanie stałe |
| **środowisko** | 4 | „w jakim trybie pracuję?" | `--dn-tekst`; kropka odziedziczona z `currentColor` | Wyłącznie emblemat; nie stosowana jako ikona przycisku |
| **rola** | 5 | „kto to robi?" | `--dn-sygnal` w klasie „inteligencja", `--dn-atrament` przy człowieku | Zawsze z plakietką roli albo etykietą nadawcy |

### 5.2. Skład grup

| Grupa | Ikony |
|---|---|
| **nawigacja** | `dom` · `menu` · `strzalka-lewo` · `strzalka-prawo` · `grot-dol` · `grot-gora` · `grot-prawo` · `wiecej` · `link-zewnetrzny` · `szukaj` · `filtr` · `zamknij` · `karta-okna` · `polecenie` |
| **działanie** | `plus` · `olowek` · `kosz` · `pobierz` · `wgraj` · `wyslij` · `odswiez` · `odpowiedz` · `uruchom` · `zatrzymaj` · `wstrzymaj` · `spinacz` · `ustawienia` · `oko` · `gwiazdka` · `kopiuj` · `slonce` · `ksiezyc` |
| **stan** | `ptaszek` · `ptaszek-kolo` · `blad` · `ostrzezenie` · `info` · `klodka` · `tarcza` · `zegar` · `historia` · `dzwonek` · `aktywnosc` |
| **obiekt** | `dokument` · `plik` · `folder` · `archiwum` · `obraz` · `kod` · `koperta` · `kalendarz` · `tabela-danych` · `baza` · `galaz` · `monitor` · `telefon` · `cpu` · `globus` |
| **moduł** | `biblioteka` · `wykres` · `terminal` · `warstwy` · `automatyzacja` · `debata` · `tlumacz` · `badanie` · `aplikacje` · `paleta` · `rozmowa` · `mikrofon` · `diagnostyka` · `agenci` · `siec` |
| **środowisko** | `srodowisko-talkin` · `srodowisko-workspace` · `srodowisko-codestudio` · `srodowisko-multitaskingai` |
| **rola** | `agent` · `uzytkownik` · `wezly` · `wykonawca` · `walidator` |

### 5.3. Grupa „moduł" wobec piętnastu modułów platformy

Platforma ma piętnaście modułów; grupa „moduł" liczy piętnaście ikon — ale **zbiory nie pokrywają się jeden do jednego**:

| Moduł | Ikona modułowa | Uwaga |
|---|---|---|
| Studio | `dokument` | Ikona z grupy **obiekt** — moduł nazwany po rodzaju treści, którą obrabia |
| Research | `badanie` | Dodatkowo `wykres` dla wyników (Findings Panel, Report Builder) |
| Library | `biblioteka` | — |
| Translate | `tlumacz` | — |
| Browser | `karta-okna` | Ikona z grupy **nawigacja** |
| Assistant | `mikrofon` | — |
| Roundtable | `debata` | — |
| Workspace | `warstwy` | — |
| Automations | `automatyzacja` | Jedyny moduł bez okna w bocznej nawigacji |
| Design | `paleta` | — |
| Apps | `aplikacje` | — |
| Terminal | `terminal` | — |
| Developer | `kod` | Ikona z grupy **obiekt**; dodatkowo `galaz` dla Git Panel |
| Diagnostics | `diagnostyka` | — |
| Agents | `agenci` | — |

Trzy moduły (**Studio**, **Browser**, **Developer**) mają ikonę spoza grupy „moduł". Grupa „moduł" mieści za to trzy pozycje niebędące modułami: `rozmowa` (Chat Window — okno wspólne), `siec` (orkiestracja MultitaskingAI), `wykres` (wyniki Research). To nie jest niespójność — grupa opisuje **rejestr znaczeniowy ikony**, a nie przypisanie do modułu.

---

## 6. Cztery emblematy środowisk

### 6.1. Zasada wspólna

Cztery emblematy są jedynymi ikonami rysowanymi od zera. Wiąże je pięć reguł:

| Reguła | Treść |
|---|---|
| **Siatka i kreska** | Ta sama co reszta zestawu: 24×24, 1,75, `round`, `currentColor` |
| **Jedna wypełniona kropka** | **Każdy emblemat zawiera dokładnie jedną wypełnioną kropkę sygnału** — `fill="currentColor" stroke="none"`. Nie zero, nie dwie |
| **Kropka dziedziczy barwę** | W pliku SVG kropka ma `currentColor`, nie wpisany błękit. Barwa sygnałowa nadawana jest w miejscu użycia (`color: var(--dn-kropka)` na elemencie zawierającym) |
| **Metafora, nie logo** | Emblemat opowiada **czynność środowiska**, nie jest znakiem marki. Znakiem marki jest sygnet „Delegacja" (`zasoby/marka/logo/sygnet.svg`) |
| **Renderowanie** | 24 px jako ikona nawigacji, **40 px** jako godło karty środowiska (`.dn-karta-srodowiska-godlo`) |

### 6.2. Konstrukcja pozycja po pozycji

| Emblemat | Środowisko i motto | Konstrukcja | Gdzie leży kropka |
|---|---|---|---|
| `srodowisko-talkin` | **TalkIn** — *Myśl. Analizuj. Rozumiej.* | Dymek rozmowy z ogonkiem w lewym dolnym rogu (`M20.5 5.5a2.5 2.5 0 0 0-2.5-2.5H6…`), wewnątrz dwie linie tekstu: pełna (`7.5 7.5 → 16.5`) i skrócona (`7.5 10.8 → 12.5`) | `cx 16,2 · cy 10,8 · r 1,5` — **domyka drugą, krótszą linię**: myśl jest w toku, nie dopowiedziana |
| `srodowisko-workspace` | **WorkSpace** — *Planuj. Organizuj. Realizuj.* | Cztery zaokrąglone kwadraty 7,2×7,2 (`rx 2`) w układzie 2×2 — plan podzielony na zadania | `cx 16,9 · cy 16,9 · r 1,6` — **w prawym dolnym kwadracie**: ostatnie zadanie planu jest w realizacji |
| `srodowisko-codestudio` | **CodeStudio** — *Projektuj. Buduj. Rozwijaj.* | Rama okna 18×16 (`rx 2,5`), wewnątrz znak zachęty powłoki `❯` (`M7 9.2l3.1 2.8L7 14.8`) i linia polecenia (`12.6 15.4 → 16.6`) | `cx 18,4 · cy 9,9 · r 1,5` — **na wysokości znaku zachęty**, po prawej: polecenie przyjęte, kompilacja biegnie |
| `srodowisko-multitaskingai` | **MultitaskingAI** — *Deleguj. Koordynuj. Nadzoruj.* | Trzy okręgi konturowe (`r 1,9`) na obwodzie — góra, lewy dół, prawy dół — połączone trzema odcinkami z węzłem środkowym | `cx 12 · cy 12 · r 2,6` — **węzeł centralny, największa kropka w zestawie**: Coordinator, przez którego przechodzi cała praca |

### 6.3. Diagram porównawczy

```
  TalkIn              WorkSpace           CodeStudio          MultitaskingAI
 ┌──────────┐        ┌────┐ ┌────┐       ┌────────────┐        ○   ─  ─  ─
 │ ──────── │        │    │ │    │       │ ❯          │         \        /
 │ ───── ●  │        └────┘ └────┘       │      ●     │          ●  ← węzeł
 └─┐  ┌─────┘        ┌────┐ ┌────┐       │   ────     │         /        \
   └──┘              │    │ │ ●  │       └────────────┘        ○   ─  ─  ○
                     └────┘ └────┘
  kropka domyka      kropka w ostatnim    kropka na linii     kropka = środek
  niedokończoną      kwadracie planu      znaku zachęty       delegacji
  myśl
```

### 6.4. Emblemat a znak marki — rozgraniczenie

| | Sygnet „Delegacja" | Emblemat środowiska |
|---|---|---|
| **Co znaczy** | Marka Danaco Console | Tryb pracy Operatora |
| **Kształt** | Podwójny grot `»»` + kropka na linii bazowej | Metafora czynności środowiska |
| **Plik** | `zasoby/marka/logo/sygnet.svg` (siatka 96×96) | `zasoby/ikony/svg/srodowisko-*.svg` (siatka 24×24) |
| **Barwa kropki** | Wpisana wprost — błękit sygnałowy | `currentColor`, nadawana w miejscu użycia |
| **Gdzie** | Pasek górny, favicon, ikona aplikacji | Karta środowiska (Strefa 1), nagłówek powłoki |
| **Wariant uproszczony** | Tak — jeden grot poniżej 24 px | Nie — emblemat nie schodzi poniżej 24 px |

---

## 7. Zasady użycia

### 7.1. Rozmiar ikony względem tekstu

| Stopień tekstu | Żeton tekstu | Ikona towarzysząca | Uzasadnienie |
|---|---|---|---|
| 11 px (`xs`) | `--dn-fs-xs` | 12 px (wyjątek plakietki/kroku) albo 14 px | Ikona nigdy mniejsza niż 12 px — poniżej kreska 1,75 zlewa się |
| 12 px (`sm`) | `--dn-fs-sm` | **14 px** | Metadane, wpis komunikacji |
| 13 px (`base`) | `--dn-fs-base` | **16 px** | Wartość domyślna interfejsu |
| 14 px (`md`) | `--dn-fs-md` | **16 px** | Treść wpisów; ikona nie rośnie razem z tekstem |
| 16 px (`lg`) | `--dn-fs-lg` | **20 px** | Nagłówki paneli |
| ≥ 20 px (`xl`+) | `--dn-fs-xl`… | **24 px** | Nagłówki okien i modali, pusty stan |

**Zasada kciuka:** ikona jest o jeden do trzech pikseli **większa** od stopnia tekstu, przy którym stoi. Nigdy mniejsza — optycznie zapada się poniżej linii pisma.

### 7.2. Odstęp od etykiety

| Kontekst | Odstęp | Żeton |
|---|---|---|
| Ikona w przycisku (`.dn-btn`) | **8 px** | `--dn-od-2` (`gap`) |
| Ikona w plakietce (`.dn-plakietka`) | **4 px** | `--dn-od-1` |
| Ikona w pozycji bocznej nawigacji | **12 px** | `--dn-od-3` |
| Ikona w toaście (`.dn-toast`) | **12 px** | `--dn-od-3` |
| Ikona w pustym stanie (nad tytułem) | **8 px** | `--dn-od-2` (kierunek pionowy) |
| Ikona w polu wyszukiwania | **32 px** wcięcia pola | `--dn-od-8` (`padding-left`) |

Odstęp jest **zawsze z żetonu** i zawsze realizowany przez `gap`, nie przez `margin` na samej ikonie. Ikona nie nosi własnych marginesów — nosi je pojemnik.

### 7.3. Ikona sama vs. ikona z etykietą

```
  IKONA + ETYKIETA                     IKONA SAMA
  ────────────────                     ──────────
  ✔ akcja rzadka lub nieodwracalna     ✔ akcja powszechnie rozpoznawalna
  ✔ pozycja nawigacji                    (zamknij ✕ · nowa karta + · wyślij ➤)
  ✔ plakietka stanu                    ✔ pasek narzędzi o wysokiej gęstości
  ✔ przycisk główny widoku             ✔ powtarzalny wiersz tabeli
  ✔ pozycja Strefy 3                   ✔ przycisk ikonowy paska górnego
                                       ── wymóg bezwzględny: aria-label
                                          + dymek .dn-tooltip
```

| Sytuacja | Rozstrzygnięcie |
|---|---|
| Akcja występuje raz w widoku | **Z etykietą** |
| Akcja powtarza się w każdym wierszu listy | Sama, z `aria-label` i dymkiem |
| Akcja jest nieodwracalna (`kosz`) | **Zawsze z etykietą** albo z potwierdzeniem w modalu |
| Ikona stoi na pasku górnym | Sama — pasek jest gęsty, a zestaw ikon paska jest stały i wyuczalny |
| Ikona oznacza stan | **Zawsze z etykietą** — patrz 7.4 |

### 7.4. Zakaz: ikona jako jedyny nośnik znaczenia stanu

> **Stan nigdy samym kolorem — zawsze ikona albo etykieta.** (kontrakt systemu projektowego)

Zasadę czytamy w obie strony. Kolor sam nie wystarczy — ale **ikona sama też nie wystarczy**, gdy komunikuje stan. Powód: użytkownik nie ma słownika ikon stanu przed oczami, a różnica między `ptaszek-kolo` a `walidator` przy 14 px jest subtelna.

| Niepoprawnie | Poprawnie |
|---|---|
| `<span class="dn-kropka dn-kropka--sukces"></span>` samodzielnie w komórce tabeli | `<span class="dn-plakietka dn-plakietka--sukces"><svg…ptaszek-kolo/> sukces</span>` |
| Zielona ikona `ptaszek-kolo` jako całe oznaczenie kroku | `.dn-krok--poprawny` z `.dn-krok-znak` (ikona) **oraz** treścią kroku i `.dn-krok-meta` |
| Sama ikona `ostrzezenie` w wierszu kolejki | Ikona + etykieta „błędy" + licznik obiegów w `.dn-krok-meta` |

Wyjątek jest jeden i wynikowo pozorny: **kropka sygnału na karcie sesji**. Kropka nie komunikuje tam wyniku (sukces/błąd), lecz fakt „coś biegnie" — a to samo niesie już animacja tętna i tytuł karty. Kropka jest wzmocnieniem, nie jedynym nośnikiem.

### 7.5. Czego nie robimy z ikoną

| Zakaz | Zamiast tego |
|---|---|
| **Emoji jako ikona** | Wyłącznie SVG z `zasoby/ikony/svg/` |
| Ikona spoza zestawu wklejona ad hoc | Procedura z rozdz. 9 |
| Wpisanie barwy w atrybut `stroke` | `currentColor` + `color` na rodzicu |
| Skalowanie przez `transform: scale()` | `width`/`height` z żetonu wymiaru |
| Obrót w celu uzyskania wariantu kierunkowego | Osobny plik (`grot-dol`, `grot-gora`, `grot-prawo`) |
| Ikona jako tło (`background-image`) | Element `<svg>` inline — inaczej `currentColor` przestaje działać |
| Dwie ikony obok siebie bez tekstu | Jedna ikona + etykieta |
| Ikona wewnątrz ikony | Jeden symbol na jedno pojęcie |

---

## 8. Dostępność ikon

### 8.1. Dwa tryby, jedna decyzja

Wybór trybu rozstrzyga jedno pytanie: **czy usunięcie ikony odbiera informację?**

| Odpowiedź | Tryb | Zapis |
|---|---|---|
| Nie — obok stoi tekst niosący to samo | **Dekoracyjna** | `aria-hidden="true"` (wartość źródłowa plików) |
| Tak — ikona jest jedynym nośnikiem | **Znacząca** | `role="img"` + `aria-label="…"` |

```
 ┌─────────────────────────────────────────────────────────┐
 │  Czy obok ikony jest tekst niosący to samo znaczenie?     │
 └───────────────┬───────────────────────┬─────────────────┘
            TAK  ▼                       ▼  NIE
   ┌────────────────────────┐  ┌──────────────────────────────┐
   │ <svg aria-hidden="true">│  │ Czy ikona jest w przycisku?   │
   │  ikona dekoracyjna      │  └────┬──────────────────┬──────┘
   └────────────────────────┘   TAK  ▼                  ▼  NIE
                    ┌────────────────────────┐  ┌────────────────────┐
                    │ aria-label NA PRZYCISKU │  │ role="img"          │
                    │ svg zostaje aria-hidden │  │ + aria-label na svg │
                    └────────────────────────┘  └────────────────────┘
```

### 8.2. Wzorce zapisu

**Ikona dekoracyjna — obok jest etykieta**

```html
<button class="dn-btn dn-btn--atrament">
  <svg aria-hidden="true" width="16" height="16" viewBox="0 0 24 24" fill="none"
       stroke="currentColor" stroke-width="1.75" stroke-linecap="round"
       stroke-linejoin="round"><path d="M5 12h14"/><path d="M12 5v14"/></svg>
  Nowa karta sesji
</button>
```

**Przycisk ikonowy — etykieta na przycisku, ikona milczy**

```html
<button class="dn-btn-ikona" aria-label="Zamknij kartę sesji Studio">
  <svg aria-hidden="true" width="16" height="16" viewBox="0 0 24 24" fill="none"
       stroke="currentColor" stroke-width="1.75" stroke-linecap="round"
       stroke-linejoin="round"><path d="M18 6 6 18"/><path d="m6 6 12 12"/></svg>
</button>
```

**Ikona znacząca poza kontrolką — np. emblemat w nagłówku powłoki**

```html
<svg role="img" aria-label="Środowisko CodeStudio" width="24" height="24"
     viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="1.75"
     stroke-linecap="round" stroke-linejoin="round">…</svg>
```

### 8.3. Przycisk ikonowy — komplet wymagań

| Wymaganie | Realizacja |
|---|---|
| Nazwa dostępna | `aria-label` na `<button>`, nie na `<svg>` |
| Cel dotykowy | 32 px (`--dn-wym-ikonowy`); przy `pointer: coarse` żeton podnosi do **40 px** — bez wyjątku w komponencie |
| Fokus widoczny | `:focus-visible` → `outline: 2px solid var(--dn-fokus)` + odsunięcie 2 px |
| Stan wciśnięty | `aria-pressed="true"` → tło `--dn-sygnal-tlo`, barwa `--dn-sygnal` |
| Stan zajęty | `aria-busy="true"` → spinner; **przycisk pozostaje klikalny** (zasada zero blokad) |
| Objaśnienie | `.dn-tooltip` przy każdej ikonie bez etykiety |
| Zakaz | **`disabled` nie występuje** — zamiast blokady komunikat po naciśnięciu |

### 8.4. Kontrast ikony

Ikona jest kreską, nie tekstem — obowiązuje ją próg **3:1** (WCAG 2.1, kryterium 1.4.11 „Non-text Contrast").

| Barwa ikony | Na tle | Ocena |
|---|---|---|
| `--dn-tekst` | `--dn-tlo`, `--dn-powierzchnia` | Z dużym zapasem (16,1:1 jasny / 15,0:1 ciemny) |
| `--dn-tekst-2` | `--dn-tlo`, `--dn-powierzchnia` | Z zapasem (6,4:1 jasny / 6,8:1 ciemny) |
| `--dn-tekst-3` | `--dn-tlo` | **Wyłącznie ikony dekoracyjne**, nigdy nośnik jedynej informacji |
| `--dn-rama-tekst-2` | `--dn-rama` | Pasek górny — ikona w spoczynku; na hover przechodzi na `--dn-rama-tekst` |
| barwy stanu | odpowiednie `--dn-*-tlo` | Pary zmierzone w `zasoby/zetony/kontrasty.json` |

### 8.5. Nawigacja klawiaturą

| Klawisz | Działanie |
|---|---|
| `Tab` | Przejście do kolejnego przycisku ikonowego; kolejność DOM = kolejność wzrokowa |
| `Enter` / `Spacja` | Uruchomienie akcji przycisku ikonowego |
| `Esc` | Zamknięcie dymka, menu, modala otwartego ikoną |
| — | **Sama ikona (`<svg>`) nigdy nie jest w kolejności fokusu.** Fokusowalny jest przycisk, nie symbol |

---

## 9. Jak dodać ikonę

### 9.1. Procedura — siedem kroków

```
 1. UZASADNIENIE      ─►  Czy pojęcie ma pokrycie w dokumentacji produktu?
        │                 NIE → koniec. Zestaw jest zamknięty.
        ▼
 2. GRUPA             ─►  Do której z 7 grup semantycznych należy?
        │                 Brak grupy → pojęcie nie jest ikonograficzne.
        ▼
 3. KOLIZJA           ─►  Czy któraś z 82 pozycji już to znaczy?
        │                 TAK → użyj istniejącej.
        ▼
 4. LUCIDE            ─►  Czy Lucide ma ten kształt?
        │                 TAK → pobierz, ustaw stroke-width 1.75.
        │                 NIE → rysuj na siatce 24×24 tą samą kreską.
        ▼
 5. NAZWA             ─►  Polska, opisowa, małe litery, łącznik.
        ▼
 6. PLIK + MANIFEST   ─►  svg/<nazwa>.svg  +  wpis {nazwa, zrodlo, zastosowanie}
        │                 + inkrementacja "liczba-ikon".
        ▼
 7. WALIDACJA         ─►  Lista sprawdzeń 9.2 — wszystkie pozycje na TAK.
```

### 9.2. Lista sprawdzeń nowej ikony

| # | Sprawdzenie | Kryterium |
|---:|---|---|
| 1 | `viewBox` | dokładnie `0 0 24 24` |
| 2 | `width` / `height` | `24` / `24` w pliku źródłowym |
| 3 | `fill` | `none` na `<svg>`; `currentColor` dopuszczalne wyłącznie na kropce/plamce z `stroke="none"` |
| 4 | `stroke` | `currentColor` — **żadnej wartości szesnastkowej w pliku** |
| 5 | `stroke-width` | `1.75` |
| 6 | `stroke-linecap` / `stroke-linejoin` | `round` / `round` |
| 7 | `aria-hidden` | `true` w pliku źródłowym |
| 8 | Nazwa pliku | polska, małe litery, bez diakrytyków, łącznik jako separator |
| 9 | Wpis w manifeście | `nazwa` = nazwa pliku bez rozszerzenia; `zrodlo` = `lucide:<slug>` albo `wlasna`; `zastosowanie` = zdanie po polsku |
| 10 | `liczba-ikon` | równa długości tablicy `ikony` i liczbie plików w `svg/` |
| 11 | Czytelność | rozpoznawalna przy **14 px** w obu motywach |
| 12 | Kontrast | ≥ 3:1 wobec tła, na którym będzie stała |
| 13 | Brak kolizji | żadna z pozostałych pozycji nie znaczy tego samego |
| 14 | Emblemat (jeśli własny) | dokładnie **jedna** wypełniona kropka sygnału |

### 9.3. Aktualizacja manifestu — wzór wpisu

```json
{
  "nazwa": "<polska-nazwa-opisowa>",
  "zrodlo": "lucide:<slug-biblioteczny>",
  "zastosowanie": "<zdanie po polsku: do czego służy, w którym module/oknie>"
}
```

Pozycję wstawiamy **w miejscu wynikającym z grupy semantycznej**, nie na koniec tablicy. Manifest jest uporządkowany tematycznie (nawigacja → działanie → stan → obiekt → moduł → środowisko), a nie alfabetycznie — kolejność jest częścią jego czytelności.

### 9.4. Usunięcie ikony z zestawu

| Krok | Treść |
|---|---|
| 1 | Sprawdź, czy nazwa nie występuje w żadnym prototypie `05-okna/` ani w opracowaniach `02-dokumentacja-html/` |
| 2 | Usuń plik `svg/<nazwa>.svg` |
| 3 | Usuń wpis z tablicy `ikony` |
| 4 | Zmniejsz `liczba-ikon` |
| 5 | Zaktualizuj listę nazw w kontrakcie systemu projektowego i katalog w rozdz. 4 niniejszego opracowania |

---

# CZĘŚĆ II — WIDŻETY

## 10. Definicja widżetu w Danaco Console

### 10.1. Definicja

> **Widżet** to złożony, samodzielny zestaw komponentów `.dn-*`, który prezentuje **stan systemu** — nie samą treść — i jest gotowy do osadzenia w oknie jako całość, bez dodatkowej kompozycji.

Trzy człony definicji są rozłączne i wszystkie muszą być spełnione:

| Człon | Znaczenie | Kontrprzykład |
|---|---|---|
| **złożony** | Powstaje z co najmniej dwóch komponentów albo z komponentu + żywych danych | `.dn-btn` sam nie jest widżetem |
| **samodzielny** | Można go wstawić w dowolne okno bez wiedzy o otoczeniu | Nagłówek okna operacyjnego nie jest widżetem — zależy od okna |
| **prezentuje stan systemu** | Odczytujemy z niego, co robi platforma | `.dn-tabela` z listą źródeł prezentuje dane, nie stan → nie jest widżetem |

### 10.2. Miejsce widżetu w warstwach systemu

```
  1 · PRYMITYWY          --dn-szary-500, --dn-sygnal-500
        ▼
  2 · ŻETONY SEMANTYCZNE --dn-tekst, --dn-kropka, --dn-czas-tetno
        ▼
  3 · KOMPONENTY .dn-*   .dn-btn · .dn-kropka · .dn-plakietka · .dn-postep
        ▼
  4 · WIDŻETY            ◄── TA WARSTWA (13 pozycji, rozdz. 11–12)
        ▼
  5 · OKNO OPERACYJNE    Studio Editor · Workflow Builder · Chat Window
        ▼
  6 · MODUŁ              Studio · Automations · Agents …
        ▼
  7 · ŚRODOWISKO         TalkIn · WorkSpace · CodeStudio · MultitaskingAI
```

Widżet sięga **wyłącznie po warstwę 3** (komponenty) i po żetony semantyczne. Nigdy po prymitywy, nigdy po inne widżety.

### 10.3. Widżet a komponent — pięć różnic

| Kryterium | Komponent `.dn-*` | Widżet |
|---|---|---|
| Definicja | Klasa CSS w `komponenty.css` | Zestaw klas + reguła kompozycji |
| Zawartość | Struktura, brak danych | Struktura **plus** dane operacyjne |
| Stan | Stany interakcji (hover, fokus, wybrany) | Stany **systemu** (pracuje, poprawny, błędy, wstrzymany) |
| Ruch | Brak, poza mikroprzejściem | Może nieść **jeden ruch znaczący** (tętno) |
| Nazwa | Angielska klasa `.dn-*` | Polska nazwa opisowa |

### 10.4. Kryterium wejścia do katalogu

Katalog obejmuje **wyłącznie widżety mające pokrycie w dokumentacji źródłowej**: w `komponenty.css` (istniejąca klasa złożona), w kontrakcie systemu projektowego (opisany element powłoki) albo w inwentarzach `brief/`. Widżet bez pokrycia nie wchodzi do katalogu — zgodnie z zasadą nadrzędną projektu.

---

## 11. Katalog widżetów

| # | Widżet | Klasa wiodąca | Grupa | Ruch | Gdzie występuje |
|---:|---|---|---|---|---|
| **1** | **Karta środowiska** | `.dn-karta-srodowiska` | wejście | wstęga na hover | Strona główna — Strefa 1 (4 wystąpienia) |
| **2** | **Kafel komponentu własnego** | `.dn-kafel` | wejście | brak | Strona główna — Strefa 2 (4 wystąpienia) |
| **3** | **Karta sesji z kropką tętna** | `.dn-karta-sesji` + `.dn-kropka--tetno` | stan | **tętno 2,4 s** | Pas kart sesji — 4 powłoki środowisk |
| **4** | **Wskaźnik postępu kolejki** | `.dn-postep` | stan | przejście szerokości 0,22 s | Execution Monitor · Queue Manager · Monitor procesu |
| **5** | **Krok kolejki** | `.dn-krok` | stan | brak | Queue Manager · Workflow Builder · Coordinator Chat |
| **6** | **Monitor procesu** | `.dn-tabela` + `.dn-plakietka` | stan | zmiana plakietek na żywo | Panel orkiestracji, sekcja Monitor procesu |
| **7** | **Rdzeń Always On Display** | `.dn-aod` | tożsamość | **tętno pierścienia 2,4 s** | Warstwa AOD (z-index 1200) — ponad całą powłoką |
| **8** | **Plakietka roli** | `.dn-plakietka--rola` | tożsamość | brak | Okna ról MultitaskingAI · Chat Window (klasa `tool`) |
| **9** | **Wpis komunikacji** | `.dn-wpis` | komunikacja | tętno kropki w `--pracuje` | **Chat Window** — wszystkie 15 modułów |
| **10** | **Pasek promptu** | `.dn-prompt` | komunikacja | pierścień fokusu 0,16 s | **Chat Window** — wszystkie 15 modułów |
| **11** | **Pusty stan** | `.dn-pusty-stan` | stan | brak | Nowa karta sesji · panele przed konfiguracją |
| **12** | **Toast** | `.dn-toast` | powiadomienie | wejście 0,22 s | Warstwa powiadomień (z-index 1000) — globalnie |
| **13** | **Przybornik okna** | `.dn-przybornik` | działanie | mikroprzejścia przycisków | Pasek promptu · nagłówki okien operacyjnych |

### 11.1. Rozkład widżetów na okna

| Okno / powłoka | Widżety obecne |
|---|---|
| **Strona główna — Centrum dowodzenia** | Karta środowiska (×4) · Kafel komponentu własnego (×4) · Rdzeń Always On Display · Toast |
| **Powłoka TalkIn / WorkSpace / CodeStudio** | Karta sesji · Wpis komunikacji · Pasek promptu · Pusty stan · Toast · Przybornik okna · Rdzeń Always On Display |
| **Powłoka MultitaskingAI** | Karta sesji · Wskaźnik postępu kolejki · Krok kolejki · Monitor procesu · Plakietka roli · Wpis komunikacji · Pasek promptu · Toast · Przybornik okna · Rdzeń Always On Display |
| **Chat Window** (wspólne, 15 modułów) | Wpis komunikacji · Pasek promptu · Przybornik okna · Pusty stan (historia pusta) |
| **Execution Monitor / Queue Manager** (Automations) | Wskaźnik postępu kolejki · Krok kolejki · Monitor procesu · Pusty stan · Toast |
| **Okno Konfiguracji** | Toast · Przybornik okna |
| **Always On Display** | Rdzeń Always On Display · Monitor procesu w trybie nadzoru |

---

## 12. Karty widżetów

### Karta środowiska

**Anatomia**

```
 ┌────────────────────────────────────────────────┐  ← ::before  wstęga 2 px
 │                                                 │    (hover / aria-current)
 │   ◈  emblemat 40 px                             │  ← -godlo
 │                                                 │
 │   TalkIn                                        │  ← -tytul   Space Grotesk 30 px
 │   Myśl. Analizuj. Rozumiej.                     │  ← -motto   Plex Mono 12 px
 │                                                 │
 │   Środowisko wiedzy, komunikacji i pracy        │  ← -opis    13 px, max 44ch
 │   z treścią. Dziewięć modułów.                  │
 └────────────────────────────────────────────────┘
   promień 14 px (--dn-r-xl) · padding 24 px · translateY(-2px) na hover
```

| Aspekt | Treść |
|---|---|
| **Dane** | Emblemat środowiska (Monitor procesu zestawu ikon) · nazwa własna środowiska · motto · zdanie opisowe · liczba modułów |
| **Stany** | Spoczynek (bez wstęgi) · najechanie (wstęga 2 px `--dn-kropka`, `--dn-cien-2`, `translateY(-2px)`) · `aria-current="true"` — środowisko z otwartymi kartami sesji (wstęga trwała) · fokus (pierścień 2 px) |
| **Zachowanie** | Klik / `Enter` zamyka stronę główną i otwiera powłokę środowiska; przy powrocie odtwarza wszystkie karty sesji tego środowiska |
| **Gdzie** | Strona główna — **Strefa 1**, dokładnie cztery wystąpienia |
| **Skład** | `.dn-karta-srodowiska` · `-godlo` · `-tytul` · `-motto` · `-opis` · `<svg>` emblematu 40 px |

**Przykładowe treści (z domeny produktu)**

| Środowisko | Motto | Moduły |
|---|---|---|
| TalkIn | Myśl. Analizuj. Rozumiej. | 9 |
| WorkSpace | Planuj. Organizuj. Realizuj. | 9 |
| CodeStudio | Projektuj. Buduj. Rozwijaj. | 8 |
| MultitaskingAI | Deleguj. Koordynuj. Nadzoruj. | panel orkiestracji (6 sekcji) |

---

### Kafel komponentu własnego

**Anatomia**

```
 ┌──────────────────────────────────────────┐
 │  ┌────┐   Automations                     │  ← -etykieta  13 px, półgruba
 │  │ ◈  │   zbuduj automatykę               │  ← -opis      12 px, --dn-tekst-3
 │  └────┘                                   │
 └──────────────────────────────────────────┘
   -ikona 36×36 (tu: automatyzacja) · promień 10 px · tło --dn-panel → --dn-powierzchnia na hover
```

| Aspekt | Treść |
|---|---|
| **Dane** | Ikona komponentu · nazwa komponentu własnego · etykieta działania („zbuduj" vs „wejdź") |
| **Stany** | Spoczynek · najechanie (obrys `--dn-obrys-mocny`, tło `--dn-powierzchnia`) · fokus (pierścień 2 px) |
| **Zachowanie** | Klik otwiera **okno budowy komponentu** albo modal kreatora — nie przestrzeń roboczą środowiska. Etykieta działania rozróżnia „zbuduj" (nowy wytwór) od „wejdź" (istniejący) |
| **Gdzie** | Strona główna — **Strefa 2**, dokładnie cztery wystąpienia |
| **Skład** | `.dn-kafel` · `-ikona` · `-etykieta` · `-opis` |
| **Waga wobec Karta środowiska** | Krój **bazowy**, nie nagłówkowy; brak wstęgi; mniejszy promień. Różnica wagi niesie hierarchię stref |

**Cztery kafle (dokumentacja)**

| Kafel | Ikona | Etykieta działania | Otwiera |
|---|---|---|---|
| **Automations** | `automatyzacja` | zbuduj automatykę | Workflow Builder |
| **Agents** | `agenci` | zbuduj agenta | Agent Builder |
| **Workspace** | `warstwy` | otwórz projekt | Project Dashboard |
| **Assistant** | `mikrofon` | skonfiguruj profil | Voice Console |

---

### Karta sesji z kropką tętna

**Anatomia**

```
  ▔▔▔▔▔▔▔▔▔▔▔▔▔▔▔▔  ← ::before wstęga 2 px (tylko aria-selected)
 ┌────────────────┐┌───────────────┐┌───┐
 │ ● Studio     ✕ ││   Research  ✕ ││ + │
 └────────────────┘└───────────────┘└───┘
   ▲ aktywna         ▲ w tle           ▲ nowa
   ● = .dn-kropka--tetno (praca w tle, 2,4 s)
   ✕ = .dn-btn-ikona, ikona zamknij, widoczna na hover
```

| Aspekt | Treść |
|---|---|
| **Dane** | Tytuł karty (nazwa modułu w środowiskach modułowych; **nazwa procesu orkiestracji** w MultitaskingAI) · wskaźnik pracy w tle · kontrolka zamknięcia |
| **Stany** | Aktywna (`aria-selected="true"` — tło `--dn-powierzchnia`, obrys, wstęga sygnału) · w tle (tło przezroczyste, tekst `--dn-tekst-2`) · najechanie · z pracą w tle (kropka tętni) · bez pracy (brak kropki) · fokus |
| **Zachowanie** | Klik w dowolnym miejscu poza `✕` przełącza obszar roboczy na układ, historię i kontekst tej karty. `✕` zamyka widok karty — proces sesji trwa po stronie serwera. `+` tworzy nową kartę w stanie pustym (Pusty stan) |
| **Ruch** | **Tętno kropki 2,4 s** — jedyny ruch ciągły interfejsu. Przy `prefers-reduced-motion` żeton skraca animację, a `.dn-kropka--tetno` otrzymuje **pierścień statyczny** 2 px |
| **Gdzie** | Pas kart sesji (`--dn-wym-pas-kart` = 36 px) we wszystkich czterech powłokach |
| **Skład** | `.dn-karty-sesji` › `.dn-karta-sesji` + `.dn-kropka--tetno` + `.dn-btn-ikona` (`zamknij`) + `.dn-btn-ikona` (`plus`) |

**Przykładowe karty (z domeny produktu)**

| Powłoka | Karty |
|---|---|
| TalkIn | `● Studio` · `Research` · `+` |
| CodeStudio | `● Developer` · `Terminal` · `Diagnostics` · `+` |
| MultitaskingAI | `● Zespół — analiza repozytorium` · `Zespół — przegląd źródeł` · `+` |

---

### Wskaźnik postępu kolejki

**Anatomia**

```
  ┌──────────────────────────────────────────┬──────────────┐
  │ ████████████████████░░░░░░░░░░░░░░░░░░░░ │  7 / 11      │
  └──────────────────────────────────────────┴──────────────┘
     -tor 4 px, promień pill, tło --dn-powierzchnia-2
     -wartosc: --dn-sygnal-wypelnienie, transition width 0,22 s
     -etykieta: Plex Mono 11 px, tabular-nums, --dn-tekst-2
```

| Aspekt | Treść |
|---|---|
| **Dane** | Liczba kroków zakończonych · liczba kroków ogółem · wartość procentowa jako szerokość wypełnienia |
| **Stany** | Zero (tor pusty) · w toku (wypełnienie rośnie) · pełny (100%) · wstrzymany (wypełnienie zamiera, etykieta niesie słowo „wstrzymana") |
| **Zachowanie** | Szerokość `.dn-postep-wartosc` zmienia się przejściem `--dn-czas-3` (0,22 s) przy każdej aktualizacji z kanału WebSocket. **Etykieta jest obowiązkowa** — pasek sam nie komunikuje stanu (zasada 7.4) |
| **Dostępność** | `role="progressbar"` + `aria-valuenow` / `aria-valuemin` / `aria-valuemax` na `.dn-postep-tor`; etykieta powiązana `aria-labelledby` |
| **Gdzie** | Execution Monitor · Queue Manager · Panel orkiestracji, sekcja **Monitor procesu** · Build Output |
| **Skład** | `.dn-postep` › `.dn-postep-tor` › `.dn-postep-wartosc` + `.dn-postep-etykieta` |

**Przykładowe etykiety** — `7 / 11 kroków` · `#129 · 64%` · `kolejka wstrzymana · 3 / 11`

---

### Krok kolejki

**Anatomia**

```
 │◤ ┌──┐                                                        │
 │  │ ✓│  Analiza repozytorium — pobranie drzewa plików   00:12  │  --poprawny
 │  └──┘                                                        │
 │◤ ┌──┐                                                        │
 │  │ ◐│  Budowa promptu wykonawczego dla Executor 1      00:04  │  --pracuje
 │  └──┘                                                        │
 │◤ ┌──┐                                                        │
 │  │ ! │  Walidacja wyniku — 2 niezgodności · obieg 2/3  00:31  │  --bledy
 │  └──┘                                                        │
 │◤ ┌──┐                                                        │
 │  │ ✕│  Przekazanie do Executor 2 — kolejka wstrzymana   —     │  --wstrzymany
 │  └──┘                                                        │
  ▲ border-left 2 px = klasa semantyczna kroku
     -znak 20×20 (ikona 12 px)      -meta: Plex Mono 11 px
```

| Aspekt | Treść |
|---|---|
| **Dane** | Znak kroku (ikona stanu albo numer porządkowy) · treść kroku · metadane (czas, licznik obiegów, rola docelowa) |
| **Stany** | Cztery klasy semantyczne + stan neutralny (oczekuje) — patrz tabela poniżej |
| **Zachowanie** | Krok nie jest klikalny sam z siebie; klikalny jest wiersz w Queue Manager, gdzie krok bywa osadzony. Zmiana stanu następuje z kanału WebSocket i **nie animuje się** — jedynie przełącza klasę |
| **Gdzie** | Queue Manager · Workflow Builder · Coordinator Chat (stan kolejki) · Execution Monitor |
| **Skład** | `.dn-kolejka` › `.dn-krok` + `.dn-krok-znak` (`<svg>` 12 px) + `.dn-krok-meta` |

**Trzy wyjścia weryfikacji + dwa stany przebiegu**

| Klasa | Ikona znaku | Barwa krawędzi | Znaczenie |
|---|---|---|---|
| *(bez modyfikatora)* | numer porządkowy | `--dn-obrys` | Oczekuje w kolejce |
| `--pracuje` | `aktywnosc` | `--dn-sygnal-wypelnienie` | Krok w toku |
| `--poprawny` | `ptaszek-kolo` | `--dn-sukces-obrys` | Weryfikacja pozytywna |
| `--bledy` | `ostrzezenie` | `--dn-ostrzezenie-obrys` | Błędy + **licznik obiegów** — kolejka biegnie dalej |
| `--wstrzymany` | `blad` | `--dn-blad-obrys` | Zdarzenie nieoczekiwane — **kolejka wstrzymana** |

---

### Monitor procesu

**Anatomia**

```
 ┌──────────┬────────────────────┬───────────────────────┬─────────┐
 │ PRZEBIEG │ ROLA               │ STATUS                │ CZAS    │  ← th: Plex Mono
 ├──────────┼────────────────────┼───────────────────────┼─────────┤    11 px, wersaliki
 │ #128     │ Executor 1         │ ⬤ ✓ sukces            │ 12:04   │
 │ #129     │ Coordinator        │ ⬤ ◐ w toku            │ 12:05   │
 │ #130     │ Executor 2         │ ⬤ ◷ oczekuje          │ —       │
 └──────────┴────────────────────┴───────────────────────┴─────────┘
   wiersz 36 px (--dn-wym-wiersz) · status = .dn-plakietka + ikona + słowo
```

| Aspekt | Treść |
|---|---|
| **Dane** | Identyfikator przebiegu (mono, `tabular-nums`) · nazwa roli · plakietka stanu · znacznik czasu |
| **Stany** | Wiersz w spoczynku · najechanie (`--dn-hover`) · wybrany (`aria-selected="true"` → `--dn-sygnal-tlo`) · tabela pusta (Pusty stan w miejsce ciała) |
| **Zachowanie** | Aktualizacja na żywo kanałem WebSocket — zmienia się plakietka i czas, wiersz nie przeskakuje. Klik wiersza otwiera okno robocze roli. Nagłówek `th` jest przyklejony (`position: sticky`, `--dn-z-przybornik`) |
| **Ruch** | Brak własnego. Tętno należy do Karta sesji i Rdzeń Always On Display — monitor jest tabelą odczytu, nie źródłem ruchu |
| **Gdzie** | Panel orkiestracji, sekcja **Monitor procesu** (MultitaskingAI) · Execution Monitor · Process Monitor (Terminal) · Actions Monitor (Assistant) · tryb nadzoru **Always On Display** |
| **Skład** | `.dn-tabela` + `.dn-plakietka--sukces / --informacja / --ostrzezenie / --blad` + `.dn-kropka` + `.dn-dane` |

**Przykładowe wiersze (z domeny produktu)**

| Przebieg | Rola | Status | Czas |
|---|---|---|---|
| `#128` | Executor 1 | sukces | 12:04 |
| `#129` | Coordinator | w toku | 12:05 |
| `#130` | Executor 2 | oczekuje | — |
| `#131` | Executor 3 / Validator | błędy · obieg 2/3 | 12:07 |

---

### Rdzeń Always On Display

**Anatomia**

```
                              ┌─ ::after: pierścień 1 px --dn-sygnal-obrys,
                              │          animacja dn-tetno 2,4 s
                              ▼
 ┌───────────────────────────────────────────────────────────┐
 │  ╭─────╮                                                   │
 │  │ ◈   │   Coordinator kończy plan etapów —                │  ← -tresc 12 px,
 │  │     │   przebieg #129, 3 z 11 kroków                    │    max 36ch
 │  ╰─────╯                                                   │
 └───────────────────────────────────────────────────────────┘
    -rdzen 36×36, --dn-grad-sygnal, promień pill
    całość: position fixed, prawy dolny róg, z-index 1200
```

| Aspekt | Treść |
|---|---|
| **Dane** | Rdzeń awatara (gradient sygnałowy — jedno z trzech dozwolonych zastosowań gradientu) · jedno zdanie o stanie systemu · opcjonalnie identyfikator przebiegu |
| **Stany** | Spoczynek (rdzeń + tętno pierścienia, treść zwięzła) · aktywny (rozwinięta powierzchnia interakcji) · tryb nadzoru MultitaskingAI (plakietka roli przy nazwie) · przy `prefers-reduced-motion` — pierścień statyczny bez animacji |
| **Zachowanie** | Pływa **ponad całą powłoką** (`position: fixed`, `--dn-z-aod` = 1200) — nad modalem, pod Command Center. Nie tworzy przestrzeni roboczej i nie przeładowuje widoku. Aktywacja z listwy Strefy 3 albo z paska górnego |
| **Ruch** | **Tętno pierścienia 2,4 s** — drugie i ostatnie miejsce ruchu ciągłego w systemie (po kropce karty sesji) |
| **Gdzie** | Warstwa AOD — obecna jednocześnie we wszystkich środowiskach, modułach i oknach |
| **Skład** | `.dn-aod` › `.dn-aod-rdzen` (`::after` = pierścień tętna) + `.dn-aod-tresc` (`<strong>` dla wyróżnienia) |

**Przykładowe komunikaty (z domeny produktu)**

- „**Coordinator** kończy plan etapów — przebieg #129, 3 z 11 kroków"
- „**Execution Monitor**: automatyka zakończona, 11 z 11 kroków"
- „**Executor 3 / Validator** zgłosił 2 niezgodności — obieg 2 z 3"

---

### Plakietka roli

**Anatomia**

```
  ┌──────────────┐  ┌────────────┐  ┌──────────────────────┐
  │ COORDINATOR  │  │ EXECUTOR 1 │  │ EXECUTOR 3/VALIDATOR │
  └──────────────┘  └────────────┘  └──────────────────────┘
   Plex Mono 11 px · letter-spacing 0,14em · wersaliki
   promień 3 px (--dn-r-xs) — kanciasta, w odróżnieniu od plakietki stanu (pill)
```

| Aspekt | Treść |
|---|---|
| **Dane** | Nazwa roli albo nazwa narzędzia — **dokładnie jak w dokumentacji**, bez tłumaczenia |
| **Stany** | Bez stanów interakcji — plakietka nie jest klikalna. Wariant semantyczny nadaje kontekst: `--sygnal` dla ról inteligencji, wariant neutralny dla klasy `tool` |
| **Zachowanie** | Statyczna. Towarzyszy nazwie nadawcy w `.dn-wpis-tozsamosc` albo tytułowi karty roli. **Nie zastępuje ikony roli** — stoi obok niej |
| **Kształt jako znaczenie** | Promień `--dn-r-xs` (3 px) zamiast `--dn-r-pill`: plakietka roli jest **maszynowa**, plakietka stanu — **organiczna**. Różnica kształtu pozwala je rozróżnić bez czytania |
| **Gdzie** | Coordinator Chat · Executor Chat (1, 2) · Results Analyzer · karty ról w sekcji Role · Chat Window przy wpisach klasy `tool` |
| **Skład** | `.dn-plakietka .dn-plakietka--rola` (+ opcjonalnie `--sygnal`) |

**Zbiór wartości udokumentowanych**

| Kontekst | Wartości |
|---|---|
| Role zespołu MultitaskingAI | `COORDINATOR` · `EXECUTOR 1` · `EXECUTOR 2` · `EXECUTOR 3 / VALIDATOR` |
| Wcielenia Results Analyzer | `VALIDATOR` · `REVIEWER` · `SECURITY AUDITOR` · `ARCHITECT` · `PRODUCT OWNER` · `QA LEAD` · `ARBITRATOR` |
| Sieć podagentów | `SUBAGENT NETWORK` (do 15 pozycji zagnieżdżonych) |

---

### Wpis komunikacji

**Anatomia**

```
 ┃ ┌───┐  OPERATOR                                        12:04 │
 ┃ │ ◈ │  ────────────────────────────────────────────────────  │
 ┃ └───┘  Przygotuj zestawienie źródeł dla przebiegu #129.       │
   ▲ border-left 2 px = klasa semantyczna
   -medalion 24×24 (ikona 14 px) · -nadawca 11 px wersaliki
   -godzina Plex Mono 11 px, margin-left auto · -tresc 14 px, interlinia 1,6
```

| Aspekt | Treść |
|---|---|
| **Dane** | Ikona nadawcy w medalionie · etykieta nadawcy (wersaliki) · plakietka roli (Plakietka roli, opcjonalnie) · godzina (mono) · treść |
| **Trzy klasy semantyczne** | `--czlowiek` (atrament) · `--inteligencja` (sygnał) · `--system` (neutralna, tło wycofane). **Nie dziewięć barw dla dziewięciu nadawców** — rozróżnienie niesie komplet: ikona + etykieta + klasa krawędzi |
| **Stany** | Spoczynek · `--pracuje` (kropka tętna przy etykiecie nadawcy, `::after` na `.dn-wpis-nadawca`) · najechanie (ujawnia kontrolkę `odpowiedz`) · strumień na żywo (treść narasta) |
| **Zachowanie** | Treść dopisuje się w czasie rzeczywistym kanałem WebSocket. Kontrolka `odpowiedz` wstawia cytat do paska promptu (Pasek promptu). Zaznaczenie fragmentu czyni go przedmiotem operacji modułu |
| **Gdzie** | **Chat Window** — okno wspólne wszystkich piętnastu modułów, we wszystkich czterech środowiskach |
| **Skład** | `.dn-wpis` (+ `--czlowiek` / `--inteligencja` / `--system` / `--pracuje`) › `-medalion` › `<svg>` · `-tozsamosc` › `-nadawca` + `.dn-plakietka--rola` + `-godzina` · `-tresc` |

**Mapowanie ról kontraktu na klasy**

| Rola kontraktu | Klasa | Ikona medalionu | Plakietka roli |
|---|---|---|---|
| `user` | `--czlowiek` | `uzytkownik` | — |
| `assistant` | `--inteligencja` | `agent` | nazwa modelu albo roli |
| `system` | `--system` | `info` | — |
| `tool` | `--system` | ikona narzędzia | **nazwa narzędzia** |

**Przykładowa wymiana (z domeny produktu)**

| Nadawca | Klasa | Treść |
|---|---|---|
| OPERATOR | `--czlowiek` | „Przygotuj zestawienie źródeł dla przebiegu #129." |
| COORDINATOR | `--inteligencja --pracuje` | „Dzielę zadanie na 11 kroków. Executor 1 przejmuje kroki 1–6." |
| EXECUTION MONITOR | `--system` | „Przebieg #129 uruchomiony · 3 z 11 kroków" |

---

### Pasek promptu

**Anatomia**

```
 ┌──┬────────────────────────────────────────────────────┬──────┐
 │❯ │  napisz polecenie…                                  │  ➤   │
 │  │                                                     │      │
 └──┴────────────────────────────────────────────────────┴──────┘
  ▲ -grot: „❯" Plex Mono półgruby, --dn-kropka — sygnatura wejścia
    -obszar: textarea 40–160 px, resize none, 14 px
    :focus-within → obrys --dn-fokus + --dn-cien-sygnal
 ┌──────────────────────────────────────────────────────────────┐
 │  [spinacz] [biblioteka] [agent] [ustawienia] model           │  ← Przybornik okna
 └──────────────────────────────────────────────────────────────┘
```

| Aspekt | Treść |
|---|---|
| **Dane** | Grot `❯` (stały) · treść polecenia · przybornik (Przybornik okna) · przycisk wysyłki |
| **Stany** | Puste (podpowiedź `--dn-tekst-3`) · wypełnione · fokus (`:focus-within` → pierścień sygnału) · generowanie odpowiedzi — **pole pozostaje edytowalne** |
| **Zachowanie** | `Enter` wysyła, `Shift+Enter` łamie wiersz. Pole rośnie od 40 do 160 px, potem przewija. **Przycisk wysyłki jest zawsze klikalny**: przy pustym polu klik pokazuje krótkie ostrzeżenie zamiast blokady (zasada zero blokad) |
| **Sygnatura** | Grot `❯` w barwie `--dn-kropka` jest wizualnym echem podwójnego grotu `»»` z sygnetu marki — znak „tu wchodzi zlecenie" |
| **Gdzie** | **Chat Window** — dolna strefa pasa komunikacji (`--dn-wym-pas-komunikacji` = 320 px), wszystkie 15 modułów |
| **Skład** | `.dn-prompt` › `-grot` + `-obszar` (`<textarea>`) + `.dn-btn-ikona` (`wyslij`) · pod spodem `.dn-przybornik` |

---

### Pusty stan

**Anatomia**

```
              ┌──────────────────────┐
              │                       │
              │          ◇            │  ← <svg> 28 px, --dn-tekst-3
              │                       │
              │   Nowa karta sesji    │  ← -tytul  Space Grotesk 16 px
              │                       │
              │  Wybierz moduł z      │  ← -opis   12 px, max 40ch
              │  bocznej nawigacji,   │
              │  aby rozpocząć pracę  │
              │                       │
              │  [ Otwórz Studio ]    │  ← opcjonalny .dn-btn
              └──────────────────────┘
                padding 40 / 20 px, wyśrodkowanie w obu osiach
```

| Aspekt | Treść |
|---|---|
| **Dane** | Ikona kontekstowa · tytuł stanu · zdanie wyjaśniające · opcjonalna akcja wyjścia ze stanu |
| **Stany** | Pusty stan **sam jest stanem** innego elementu — nie ma stanów własnych poza stanami osadzonego przycisku |
| **Zachowanie** | Zastępuje treść panelu, tabeli albo obszaru roboczego. Znika, gdy pojawi się pierwsza pozycja. **Komunikuje stan oczekiwany, nie błąd ładowania** — błąd ma własną formę |
| **Dobór ikony** | Ikona wskazuje **czego brakuje**, nie „pustkę": `dokument` w Session Repository · `badanie` w Findings Panel · `automatyzacja` w Queue Manager · `karta-okna` w nowej karcie sesji |
| **Gdzie** | Nowa karta sesji (obszar roboczy) · Sources Manager · Findings Panel · Queue Manager · Library Explorer · Assets Panel · Errors Panel · historia Chat Window przed pierwszym poleceniem |
| **Skład** | `.dn-pusty-stan` › `<svg>` 28 px + `-tytul` + `-opis` (+ `.dn-btn--zarys`) |

**Przykładowe treści (z domeny produktu)**

| Miejsce | Tytuł | Opis |
|---|---|---|
| Nowa karta sesji | Nowa karta sesji | Wybierz moduł z bocznej nawigacji, aby rozpocząć pracę |
| Queue Manager | Brak zdefiniowanych kolejek | Utwórz kolejkę albo wczytaj zapisany zespół z sekcji Zespoły |
| Findings Panel | Brak ustaleń | Ustalenia cząstkowe pojawią się po pierwszym przebiegu Research Workspace |

---

### Toast

**Anatomia**

```
 ┃ ┌────────────────────────────────────────────────────┬────┐
 ┃ │ ✓  Profil izolacji zapisany                         │ ✕  │  ← -tytul 13 px
 ┃ │    Zakres: karta sesji · warstwa: sesji             │    │  ← -tresc 12 px
 ┃ └────────────────────────────────────────────────────┴────┘
  ▲ border-left 2 px w barwie semantycznej + ikona w tej samej barwie
    280–420 px · cień --dn-cien-3 · wejście dn-wejscie 0,22 s
    stos: .dn-toasty, prawy dolny róg, z-index 1000
```

| Aspekt | Treść |
|---|---|
| **Dane** | Ikona stanu · tytuł (co się stało) · treść uzupełniająca (kontekst) · kontrolka zamknięcia |
| **Warianty** | `--sukces` (`ptaszek-kolo`) · `--ostrzezenie` (`ostrzezenie`) · `--blad` (`blad`) · `--informacja` (`info`) |
| **Stany** | Pojawienie się (`dn-wejscie` 0,22 s: `opacity` + `translateY` + `scale 0.98`) · widoczny · znikanie (automatyczne po czasie ekspozycji albo ręczne przez `✕`) |
| **Zachowanie** | Nie przerywa pracy i nie wymaga reakcji. Stos pionowy w `.dn-toasty`; nowe toasty dokładają się od dołu. Potwierdza czynność, **która już się dokonała** — nie pyta o zgodę |
| **Dostępność** | Pojemnik `.dn-toasty` z `role="status"` i `aria-live="polite"`; wariant `--blad` z `aria-live="assertive"` |
| **Gdzie** | Warstwa powiadomień (`--dn-z-powiadomienie` = 1000) — globalnie, ponad każdym oknem |
| **Skład** | `.dn-toasty` › `.dn-toast` (+ wariant) › `<svg>` 16 px + `-tytul` + `-tresc` + `.dn-btn-ikona` (`zamknij`) |

**Przykładowe powiadomienia (z domeny produktu)**

| Wariant | Tytuł | Treść |
|---|---|---|
| `--sukces` | Profil izolacji zapisany | Zakres: karta sesji · warstwa: sesji |
| `--informacja` | Przebieg #129 uruchomiony | Coordinator przejął plan etapów |
| `--ostrzezenie` | Agent zapisany niekompletnie | Brak przypisanego modelu bazowego — uzupełnij w Model Configuration |
| `--blad` | Kolejka wstrzymana | Krok 7 zgłosił zdarzenie nieoczekiwane — Results Analyzer oczekuje decyzji |

---

### Przybornik okna

**Anatomia**

```
 ┌────────────────────────────────────────────────────────────────────────────────┐
 │ [spinacz] [biblioteka] [agent] [ustawienia] · [ model: kanał domyślny ▾ ] · [⋯]│
 └────────────────────────────────────────────────────────────────────────────────┘
   .dn-btn-ikona 32×32 (ikona 16 px) · gap 4 px (--dn-od-1) · flex-wrap
```

| Aspekt | Treść |
|---|---|
| **Dane** | Zestaw akcji kontekstowych okna albo paska promptu — zmienny per moduł |
| **Stany** | Każdy element przejmuje stany własnego komponentu: `.dn-btn-ikona` (spoczynek · najechanie · `aria-pressed` · fokus · `aria-busy`), `.dn-wybor` (rozwinięcie) |
| **Zachowanie** | Zawija się (`flex-wrap`) przy zwężeniu okna, nie przewija poziomo. Kolejność od najczęstszej akcji po najrzadszą. Zestaw ikon jest **stały w obrębie modułu** — Operator uczy się go raz |
| **Zasada doboru** | Wyłącznie ikony z grup **działanie** i **nawigacja**. Ikona stanu w przyborniku byłaby myląca — przybornik oferuje czynności, nie raportuje |
| **Gdzie** | Pod paskiem promptu w Chat Window · nagłówki okien operacyjnych · listwa Strefy 3 (wariant lekki) |
| **Skład** | `.dn-przybornik` › `.dn-btn-ikona` × n (+ `.dn-separator--pionowy` + `.dn-wybor` + `.dn-tooltip`) |

**Przykładowe zestawy per moduł (z domeny produktu)**

| Moduł | Ikony przybornika |
|---|---|
| **Studio** | `spinacz` · `dokument` · `historia` · `ustawienia` · `wiecej` |
| **Terminal** | `spinacz` · `uruchom` · `zatrzymaj` · `kopiuj` · `wiecej` |
| **Automations** | `spinacz` · `uruchom` · `wstrzymaj` · `odswiez` · `kalendarz` · `wiecej` |
| **Research** | `spinacz` · `globus` · `badanie` · `pobierz` · `wiecej` |

---

## 13. Decyzje projektowe

Rozstrzygnięcia podjęte **ponad źródła** — tam, gdzie dokumentacja nie rozstrzygała formy. Każde odnotowane jawnie, zgodnie z zasadą 5 rozdz. 10 kontrakt systemu projektowego.

| # | Zagadnienie | Stan źródeł | Rozstrzygnięcie | Uzasadnienie |
|---|---|---|---|---|
| **1** | **Grupy semantyczne ikon** | `manifest.json` nie zawiera pola grupy; kontrakt systemu projektowego podaje płaską listę 82 nazw | Wprowadzono **siedem grup**: nawigacja · działanie · stan · obiekt · moduł · środowisko · rola, z pełnym przypisaniem wszystkich 82 pozycji (rozdz. 5.2) | Zadanie wprost wymaga grup. Podział wyprowadzono z pola `zastosowanie` manifestu — żadna ikona nie została przypisana wbrew opisowi. Grupa jest warstwą dokumentacyjną, **nie trafia do manifestu** |
| **2** | **Liczba pozycji Lucide vs własnych** | Manifest podaje regułę („Lucide + emblematy własne"), nie podaje liczb | **78 Lucide / 4 własne** — policzone po polu `zrodlo` (`"wlasna"` występuje dokładnie 4 razy) | Liczba jest wyprowadzalna z manifestu jednoznacznie; podanie jej ułatwia audyt licencyjny |
| **3** | **Dwa wyjątki 12 px** | `komponenty.css` renderuje `.dn-plakietka > svg` i `.dn-krok-znak > svg` w 12 px, poza skalą 14/16/20/24 z manifestu | Uznano za **udokumentowany wyjątek dwóch komponentów**, nie za piąty stopień skali | Skala manifestu pozostaje wiążąca dla ikon samodzielnych. Wyjątek dotyczy ikon wewnątrz elementów 16–20 px, gdzie 14 px rozsadza pojemnik. Odnotowanie zamyka lukę zamiast ją ukrywać |
| **4** | **Zakaz „ikona jako jedyny nośnik stanu"** | kontrakt systemu projektowego formułuje zasadę jako „stan nigdy samym kolorem — zawsze ikona albo etykieta" | Zaostrzono: przy komunikacie **stanu** ikona sama także nie wystarcza — wymagana etykieta słowna (rozdz. 7.4) | Spójnik „albo" w źródle dotyczy uzupełnienia koloru. Dla stanu wymagamy pary ikona **+** etykieta, bo różnica między `ptaszek-kolo` a `walidator` przy 14 px jest nieczytelna. Wyjątek (kropka na karcie sesji) opisany i uzasadniony |
| **5** | **Przypisanie ikon do piętnastu modułów** | Manifest wiąże z modułem 12 ikon; trzy moduły (**Studio**, **Browser**, **Developer**) nie mają ikony z grupy „moduł" | Przyjęto ikony wskazane w polu `zastosowanie`: Studio → `dokument`, Browser → `karta-okna`, Developer → `kod`; rozbieżność opisano jawnie (rozdz. 5.3) | Manifest wprost przypisuje te ikony do tych modułów w opisie. Nie tworzymy nowych ikon dla trzech modułów — zestaw jest zamknięty |
| **6** | **Definicja widżetu** | Źródła używają pojęcia „wzorzec złożony" (inwentarz komponentów, 6 pozycji); nie definiują „widżetu" | Przyjęto definicję trójczłonową: **złożony + samodzielny + prezentuje stan systemu** (rozdz. 10.1), umieszczoną jako **warstwa 4** między komponentami a oknem operacyjnym | Zadanie wymaga definicji. Trójczłon pozwala jednoznacznie odrzucić kandydatów (nagłówek okna, tabela źródeł) i zgadza się z warstwowaniem systemu z katalogu komponentów |
| **7** | **Zamknięcie katalogu na trzynastu pozycjach** | Zadanie wymienia trzynaście widżetów | Każdą z trzynastu pozycji zweryfikowano wobec `komponenty.css` — **wszystkie mają klasę wiodącą albo jawny opis w inwentarzu**; nie dodano żadnej pozycji ponad listę | Katalog bez pokrycia w kodzie byłby projektowaniem, nie dokumentowaniem |
| **8** | **Kształt plakietki roli jako nośnik znaczenia** | `komponenty.css` nadaje `.dn-plakietka--rola` promień `--dn-r-xs` zamiast `--dn-r-pill`; nie wyjaśnia dlaczego | Odczytano jako **celowe rozróżnienie**: plakietka roli maszynowa (kanciasta), plakietka stanu organiczna (pill) — i opisano jako regułę (Plakietka roli) | Różnica jest w kodzie; nadanie jej znaczenia czyni ją powtarzalną, zamiast pozostawiać przypadkową |
| **9** | **Dwa miejsca ruchu ciągłego** | kontrakt systemu projektowego wymienia pięć miejsc kropki sygnału, `komponenty.css` animuje `dn-tetno` w trzech selektorach (`.dn-kropka--tetno`, `.dn-wpis--pracuje`, `.dn-aod-rdzen::after`) | Skatalogowano jako **trzy nośniki tętna** w widżetach: Karta sesji (karta sesji), Wpis komunikacji (wpis pracujący), Rdzeń Always On Display (rdzeń AOD); pozostałe widżety mają ruch **zerowy** albo mikroprzejście | `INTENSYWNOSC_RUCHU = 3/10` wymaga policzalności ruchu. Lista zamknięta zapobiega rozlaniu animacji na kolejne widżety |
| **10** | **Dostępność wskaźnika postępu** | `komponenty.css` nie deklaruje ról ARIA dla `.dn-postep` | Wskazano `role="progressbar"` + `aria-valuenow/min/max` na `.dn-postep-tor` oraz obowiązkową etykietę (Wskaźnik postępu kolejki) | WCAG 2.1 AA jest warunkiem wejściowym; pasek bez roli i bez etykiety nie przekazuje wartości technologii wspomagającej |
| **11** | **Dostępność stosu powiadomień** | `komponenty.css` nie deklaruje roli dla `.dn-toasty` | Wskazano `role="status"` + `aria-live="polite"`, a dla `--blad` — `aria-live="assertive"` (Toast) | Toast pojawia się bez akcji użytkownika; bez obszaru żywego jest niesłyszalny |
| **12** | **Zestaw ikon przybornika per moduł** | Źródła opisują przybornik jako „zależny od modułu — konkretna forma poza zakresem" | Podano **cztery przykładowe zestawy** (Studio, Terminal, Automations, Research) złożone wyłącznie z ikon istniejących w zestawie i oznaczone jako przykładowe (Przybornik okna) | Zestawy ilustrują regułę doboru (tylko grupy działanie i nawigacja), nie ustanawiają wiążącego kontraktu — źródło jawnie zostawia formę fazie wykonawczej |

### 13.1. Czego świadomie nie rozstrzygnięto

| Zagadnienie | Powód pozostawienia otwartym |
|---|---|
| Forma rozwiniętej powierzchni Always On Display (dymek / panel boczny / pełne okno) | Źródło (inwentarz komponentów, rozdz. 12.2) jawnie odsyła do fazy wykonawczej |
| Zachowanie pasa kart sesji przy przepełnieniu (przewijanie vs zwężanie) | Źródło jawnie odsyła do fazy wykonawczej |
| Czas ekspozycji toastu przed samoczynnym zniknięciem | Brak wartości w żetonach i w dokumentacji; wpisanie liczby byłoby wymyśleniem |
| Ikony dla wcieleń Results Analyzer poza `walidator` | Zestaw jest zamknięty; sześć pozostałych wcieleń różnicuje **plakietka roli**, nie ikona |

---

## Podsumowanie liczbowe

| Wielkość | Wartość |
|---|---|
| Ikon w zestawie | **82** |
| — z Lucide (ISC) | 78 |
| — emblematów własnych | 4 |
| Grup semantycznych | 7 |
| Stopni renderowania | 4 (14 · 16 · 20 · 24 px) |
| Udokumentowanych wyjątków rozmiaru | 2 (plakietka, znak kroku — 12 px) |
| Widżetów w katalogu | **13** |
| — z ruchem ciągłym (tętno) | 3 (Karta sesji, Rdzeń Always On Display, Wpis komunikacji) |
| — bez ruchu własnego | 6 (Kafel komponentu własnego, Krok kolejki, Monitor procesu, Plakietka roli, Pusty stan, Przybornik okna) |
| Klas `.dn-*` wykorzystanych w widżetach | 34 |
| Decyzji projektowych ponad źródła | 12 |

---

*Danaco Console — AI Operating Environment · v2.0*
*© 2026 Danaco Holding Group Sp. z o.o. — Dariusz Naharnowicz*
