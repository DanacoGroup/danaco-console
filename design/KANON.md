# KANON — jedyne źródło prawdy dla wszystkich zespołów projektowych

> **Ten plik jest wiążący. Każdy agent MUSI go przeczytać przed rozpoczęciem pracy.**
> Zasada nadrzędna: **nic, co nie ma pokrycia w dokumentacji źródłowej, nie może pojawić się w opracowaniu.**
> Zakaz wymyślania modułów, okien, funkcji, nazw, metryk, osób.

---

## 0. Metryka produktu

| | |
|---|---|
| **Produkt** | **Danaco Console** — AI Operating Environment (warstwa wizualna v2.0) |
| **Produkt (warstwa funkcjonalna)** | **Danaco Pilot** — Platforma AI Workspace OS (dokumentacja projektowa v1.0) |
| **Producent** | Danaco Holding Group Sp. z o.o. |
| **Twórca** | Dariusz Naharnowicz |
| **Kontakt** | support@danaco-group.pl |
| **Status** | Deweloperski |
| **Data pakietu wizualnego** | 2026-08-11 (v2.0) |
| **Data dokumentacji UI** | 2026-08-06 (v1.0) |

**Nazewnictwo w opracowaniach:** produkt nazywamy **Danaco Console**. Nazwa **Danaco Pilot** występuje wyłącznie tam, gdzie cytujemy dokumentację funkcjonalną v1.0 (można ją podać jako nazwę warstwy funkcjonalnej). Nie mieszać w jednym zdaniu bez wyjaśnienia.

---

## 1. Kontrakt kierunku projektowego (KIERUNEK.md — nienaruszalny)

> Platforma operacyjna (**kokpit dowodzenia**) dla zawodowego **Operatora** zarządzającego cyfrową organizacją, w języku **monochromatycznej precyzji** — odcienie bieli w motywie jasnym, odcienie czerni w motywie ciemnym, **jeden chłodny sygnał** — z ciążeniem ku własnemu systemowi „**instrumentu pomiarowego**": zwarta gęstość, dane krojem mono, zero dekoracji bez funkcji.

### Trzy pokrętła
| Pokrętło | Wartość |
|---|---|
| `WARIANCJA_PROJEKTOWA` | **4/10** — kokpit wymaga przewidywalności; asymetria tylko gdy niesie hierarchię |
| `INTENSYWNOSC_RUCHU` | **3/10** — mikroprzejścia 100–220 ms + jeden ruch znaczący (tętno kropki) |
| `GESTOSC_WIZUALNA` | **8/10** — gęstość **zwarta**; kontrolki 32 px, wiersze 36 px, odstępy 12–16 px |

### Anty-domyślne — katalog zablokowanych odruchów (BEZWZGLĘDNY ZAKAZ)
| Odruch | Zamiast tego |
|---|---|
| Fioletowy gradient „AI" | monochrom + jeden błękit sygnałowy |
| Inter + slate-900 jako baza | Plex Sans + czysta neutralna skala |
| Trzy równe karty funkcji | strefy o malejącej masie; asymetria koordynator–wykonawca |
| Wyszarzone przyciski jako bramy | **zero blokad** — komunikat po naciśnięciu albo opis obok |
| Stan samym kolorem | **zawsze** ikona lub etykieta |
| Emoji jako ikony | wyłącznie SVG z zestawu |
| Dziewięć kolorów tła dla dziewięciu nadawców | trzy klasy semantyczne + ikona + etykieta + plakietka roli |
| Glassmorfizm wszędzie | powierzchnie kryjące; rozmycie **tylko** w nakładce modala |
| `#000000` / `#FFFFFF` jako tło strony | `#0F0F0F` / `#F4F4F4` |
| Cień „domyślny szary wszędzie" | dwa komplety cieni tonowanych per motyw; `cien-sygnal` zarezerwowany |
| **Lorem ipsum, „Jan Kowalski", zmyślone metryki** | **treści operacyjne z domeny produktu, oznaczone jako przykładowe** |

---

## 2. Żetony — wartości wiążące (`zasoby/zetony/zetony.css`)

**ZASADA: nigdy nie wpisuj wartości szesnastkowej wprost. Zawsze `var(--dn-*)`.**
**ZASADA: komponent sięga wyłącznie po żetony SEMANTYCZNE (`--dn-tlo`, `--dn-tekst`…), nigdy po prymitywy (`--dn-szary-500`).**

### Skala szarości (prymitywy — zakaz użycia wprost)
`--dn-szary-0` #FFFFFF · `-25` #FAFAFA · `-50` #F4F4F4 · `-100` #ECECEC · `-150` #E3E3E3 · `-200` #D7D7D7 · `-300` #C0C0C0 · `-400` #9E9E9E · `-500` #7C7C7C · `-600` #616161 · `-700` #4A4A4A · `-750` #3A3A3A · `-800` #2A2A2A · `-850` #212121 · `-900` #181818 · `-925` #131313 · `-950` #0F0F0F · `-1000` #0A0A0A
**`#000000` jest ZAKAZANE.**

### Sygnał — jedyna barwa akcentu
`--dn-sygnal-100` #EDF3FE · `-200` #C9DAFB · `-300` #8FB2F5 · `-400` #5C8CEC · `-500` **#3B6FE0** (bazowy — kropka sygnału, fokus) · `-600` #2457C9 · `-700` #1D49AF · `-800` #173A8C

### Stany
zieleń `#1E7A46 / #4FBD85 / #BCE3CD / #E7F5EE` · bursztyn `#8A5B0C / #E0A63C / #EED9A7 / #FBF2DE` · czerwień `#C0362F / #EE7168 / #F3C8C4 / #FCEDEB` · **informacja = rodzina sygnału**

### Rama kokpitu — stała w OBU motywach
`--dn-rama: var(--dn-szary-925)` · `--dn-rama-tekst: szary-100` · `--dn-rama-tekst-2: szary-400` · `--dn-rama-hover: rgba(255,255,255,.08)` · `--dn-rama-obrys: rgba(255,255,255,.10)`
**Pasek górny jest zawsze atramentowy — jedyny element nieprzełączający się z motywem.**

### Typografia — trzy role, trzy kroje
| Rola | Krój | Żeton |
|---|---|---|
| Nagłówki, tytuły środowisk, logotyp | **Space Grotesk** 500–700 | `--dn-ff-naglowek` |
| Interfejs i treść | **IBM Plex Sans** 400–700 | `--dn-ff-bazowa` |
| Dane, identyfikatory, terminal | **IBM Plex Mono** 400–600 | `--dn-ff-mono` |

Stopnie: `xs` 11 · `sm` 12 · `base` **13** · `md` 14 · `lg` 16 · `xl` 20 · `2xl` 24 · `3xl` 30 · `display` 40 px
Wagi: `400 / 500 / 600 / 700`. Interlinia: ciasny 1.25 · bazowy 1.45 · luźny 1.6.
Odstępy liter: `--dn-ls-naglowek` −0.01em · `--dn-ls-wersaliki` 0.08em · `--dn-ls-mono-wersaliki` 0.14em.
Liczby: **`tabular-nums`** w Plex Mono.
Test diakrytyków: `ŁÓDŹ ŻÓŁĆ GĘŚLĄ JAŹŃ`.

### Przestrzeń — jednostka 4 px [NIENEGOCJOWALNE]
`--dn-od-0/1/2/3/4/5/6/8/10/12/16` = 0/4/8/12/16/20/24/32/40/48/64 px
Promienie: `xs` 3 · `sm` 6 · `md` 8 · `lg` 10 · `xl` 14 · `pill` 999 px

### Ruch
`--dn-ease: cubic-bezier(0.2, 0, 0, 1)` · `--dn-czas-1` 0.1s (mikroreakcje) · `--dn-czas-2` 0.16s (barwy, motyw) · `--dn-czas-3` 0.22s (modal, panel, toast) · `--dn-czas-tetno` **2.4s** (tętno kropki — jedyny ruch ciągły)

### Wymiary (gęstość zwarta)
kontrolka **32** · ikonowy **32** · pasek górny **48** · pas kart **36** · wiersz **36** · boczna nawigacja **224** · pas komunikacji **320** · modal max **560** · awatar 24/28/36 · ikona 14/16/20/24 · przełącznik 36×20 · check 16 · **kropka 6** · wstęga 2 · spinner 14 · fokus 2 px + odsunięcie 2 px
Dotyk (`pointer: coarse`): kontrolka/ikonowy **40**, wiersz **44**, check 20, przełącznik 44×24.

### Punkty łamania
`w1` 640 (telefon poziomo / Mobile) · `w2` 960 (tablet — boczna zwija się do ikon) · `w3` 1280 (biurko — pełny kokpit) · `w4` 1600 (szerokie biurko — dwa okna komunikacji)
Siatka 12 kolumn, przerwa 24 px, treść max 1200 px.

### Warstwy (z-index)
podłoga 0 · przybornik 10 · pasek 100 · boczna 200 · pas komunikacji 300 · nakładka 800 · modal 900 · powiadomienie 1000 · tooltip 1100 · **AOD 1200** · **Centrum poleceń 1300**

### Gradienty — wyłącznie ilustracyjne
`--dn-grad-atrament` · `--dn-grad-sygnal`. Zakres: awatary bez zdjęcia, rdzeń AOD, grafiki brandowe.
**Gradient NIGDY nie jest tłem przycisku, karty ani sekcji.**

---

## 3. Znak marki i element sygnaturowy

### Sygnet „Delegacja"
**Podwójny grot promptu „»" + kropka sygnału na linii bazowej** — znak czyta się „**»».**" i opowiada pętlę produktu: zlecenie przechodzi od koordynatora do wykonawcy, kropka to praca, która właśnie ruszyła.
- jednobarwny (atrament / biel) + **kropka w błękicie sygnałowym**
- wersje mono w komplecie
- **wariant uproszczony (jeden grot) obowiązuje poniżej 24 px**

### Logotyp
`DANACO` (Space Grotesk 700) + `CONSOLE` (IBM Plex Mono 500, rozstrzelone wersaliki).
Para krojów opowiada produkt: **marka + maszyna**.

### Kropka sygnału — element sygnaturowy
Wypełniony punkt błękitu sygnałowego oznaczający „**tu biegnie praca**":
- w **godle** — kropka zamykająca podwójny grot,
- w **emblematach czterech środowisk** — każdy zawiera dokładnie jedną wypełnioną kropkę,
- na **kartach sesji** — pulsująca kropka = proces w tle,
- w **oknie komunikacji** — kropka przy nadawcy aktywnie piszącym,
- w **Always On Display** — rdzeń awatara.

Tętno (2,4 s) to **jedyny ruch ciągły** w interfejsie; przy `prefers-reduced-motion` zastępuje je **pierścień statyczny**.

### Pliki znaku (w `zasoby/marka/`)
`logo/logotyp.svg` · `logo/logotyp-na-ciemnym.svg` · `logo/logo-poziomy.svg` · `logo/logo-poziomy-na-ciemnym.svg` · `logo/logo-poziomy-mono-bialy.svg` · `logo/logo-poziomy-mono-czarny.svg` · `logo/logo-pionowy.svg` · `logo/logo-pionowy-na-ciemnym.svg` · `logo/sygnet.svg` · `logo/sygnet-na-ciemnym.svg` · `logo/sygnet-mono-bialy.svg` · `logo/sygnet-mono-czarny.svg` · `logo/sygnet-uproszczony.svg` · `logo/sygnet-uproszczony-na-ciemnym.svg` · `logo/png/*` · `favicon/*` · `ikona-aplikacji/*` · `srodowiska/srodowisko-{talkin,workspace,codestudio,multitaskingai}.svg`

---

## 4. Ikony

Zestaw z **Lucide** (ISC), siatka **24×24**, obrys ujednolicony do **1,75**, wypełnienie `none`, barwa `currentColor`, zakończenia i łączenia `round`, renderowanie w 14/16/20/24 px.
**Ikon nie rysuje się ręcznie, jeżeli istnieją w bibliotece.** Dorysowuje się wyłącznie brakujące pojęcia domenowe (emblematy środowisk, rola koordynatora) na tej samej siatce.
Nazewnictwo **polskie, opisowe**. Manifest: `zasoby/ikony/manifest.json`. Pliki: `zasoby/ikony/svg/*.svg`.

Dostępne nazwy (kompletne): agenci, agent, aktywnosc, aplikacje, archiwum, automatyzacja, badanie, baza, biblioteka, blad, cpu, debata, diagnostyka, dokument, dom, dzwonek, filtr, folder, galaz, globus, grot-dol, grot-gora, grot-prawo, gwiazdka, historia, info, kalendarz, karta-okna, klodka, kod, koperta, kopiuj, kosz, ksiezyc, link-zewnetrzny, menu, mikrofon, monitor, obraz, odpowiedz, odswiez, oko, olowek, ostrzezenie, paleta, plik, plus, pobierz, polecenie, ptaszek-kolo, ptaszek, rozmowa, siec, slonce, spinacz, srodowisko-codestudio, srodowisko-multitaskingai, srodowisko-talkin, srodowisko-workspace, strzalka-lewo, strzalka-prawo, szukaj, tabela-danych, tarcza, telefon, terminal, tlumacz, uruchom, ustawienia, uzytkownik, walidator, warstwy, wezly, wgraj, wiecej, wstrzymaj, wykonawca, wykres, wyslij, zamknij, zatrzymaj, zegar

---

## 5. Biblioteka komponentów — klasy `.dn-*` (`zasoby/css/komponenty.css`)

Istniejące klasy (WIĄŻĄCE — nie zmieniać nazw, wolno rozszerzać o nowe modyfikatory):

**Przyciski:** `.dn-btn` + `--atrament --sygnal --zarys --duch --niebezpieczny --wybrany --sm --lg` · `.dn-btn-ikona` + `--na-ramie`
**Pola:** `.dn-pole` (`-etykieta -kontrolka -opis -blad`) · `.dn-wybor` · `.dn-check` · `.dn-radio` · `.dn-przelacznik` · `.dn-suwak` · `.dn-szukaj`
**Nawigacja:** `.dn-pasek` (`-godlo -logotyp -szukaj -prawa`) · `.dn-boczna` (`-naglowek -pozycja`) · `.dn-karty-sesji` · `.dn-karta-sesji` · `.dn-zakladki` · `.dn-zakladka` · `.dn-listwa` (`-pozycja`) · `.dn-przybornik`
**Dane:** `.dn-karta` (`--klikalna --wybrana`, `-naglowek -tytul -cialo`) · `.dn-karta-srodowiska` (`-godlo -tytul -motto -opis`) · `.dn-kafel` (`-ikona -etykieta -opis`) · `.dn-tabela` · `.dn-dane` · `.dn-plakietka` + `--sukces --ostrzezenie --blad --informacja --sygnal --rola` · `.dn-kropka` + `--tetno --sukces --ostrzezenie --blad --neutralna` · `.dn-pusty-stan` (`-tytul -opis`) · `.dn-postep` (`-etykieta -tor -wartosc`) · `.dn-kolejka` · `.dn-krok` + `--pracuje --poprawny --bledy --wstrzymany` (`-znak -meta`)
**Nakładki:** `.dn-modal` (`-naglowek -tytul -cialo -stopka`) · `.dn-toasty` · `.dn-toast` + `--sukces --ostrzezenie --blad --informacja` (`-tytul -tresc`) · `.dn-tooltip` (`-tresc`) · `.dn-spinner`
**Tożsamość:** `.dn-awatar` + `--sm --lg --kwadrat --inteligencja` (`-stan`) · `.dn-aod` (`-rdzen -tresc`)
**Komunikacja:** `.dn-wpis` + `--czlowiek --inteligencja --system --pracuje` (`-medalion -nadawca -tozsamosc -godzina -tresc`) · `.dn-prompt` (`-grot -obszar`)

**Mapowanie ról kontraktu → klasy:** `user → --czlowiek` · `assistant → --inteligencja` · `system → --system` · `tool → --system` (z plakietką roli = nazwa narzędzia).

---

## 6. Architektura platformy — model pojęciowy

```
URUCHOMIENIE ─► OKNO STARTOWE ─► REJESTRACJA/LOGOWANIE ─► STRONA GŁÓWNA (Centrum dowodzenia)
                                                                   │
        ┌──────────────────────────────┬───────────────────────────┴──────────┐
        ▼ STREFA 1                       ▼ STREFA 2                    ▼ STREFA 3
   ŚRODOWISKO                      KOMPONENT WŁASNY                 USTAWIENIA
   TalkIn · WorkSpace ·            Automations · Agents ·           Okno konfiguracji
   CodeStudio · MultitaskingAI     Workspace · Assistant            Mobile · Always On Display
        │                                │                                │
        ▼                                ▼ (wpięcie w sesji modułu)        ▼
   MODUŁ (boczna nawigacja)                                         funkcja globalna
   albo ROLA (panel orkiestracji)
        │
        ▼
   OKNO OPERACYJNE = Chat Window (wspólne) + okna właściwe modułowi
```

**Trójstopniowa hierarchia:** Środowisko („w jakim trybie pracuję?") → Moduł („jakie zadanie wykonuję?") → Okno operacyjne („jakim narzędziem realizuję?").

### Cztery środowiska
| Środowisko | Emblemat | Moduły w bocznej nawigacji |
|---|---|---|
| **TalkIn** | `srodowisko-talkin.svg` | 9 |
| **WorkSpace** | `srodowisko-workspace.svg` | 9 |
| **CodeStudio** | `srodowisko-codestudio.svg` | 8 |
| **MultitaskingAI** | `srodowisko-multitaskingai.svg` | panel orkiestracji (6 sekcji) zamiast modułów |

### Piętnaście modułów × środowisko
| Moduł | TalkIn | WorkSpace | CodeStudio | MultitaskingAI |
|---|:-:|:-:|:-:|:-:|
| Studio | ● | ● | | |
| Research | ● | ● | | |
| Library | ● | ● | | |
| Translate | ● | | | |
| Browser | ● | ● | | |
| Assistant | ● | | | |
| Roundtable | ● | ● | ● | |
| Workspace | ● | ● | ● | |
| Automations | ○ | ○ | ○ | integracja 24/7 |
| Design | | ● | ● | |
| Apps | | ● | ● | |
| Terminal | | | ● | |
| Developer | | | ● | |
| Diagnostics | | | ● | |
| Agents | ● | ● | ● | rola/ekspert |

`●` = okno w bocznej nawigacji · `○` = wyłącznie komponent własny strefy 2 (bez okna modułowego)

**Podwójna rola:** Assistant (Profil asystenta), Workspace (Projekt), Automations (Automatyka) są jednocześnie modułami i komponentami własnymi strefy 2. **Automations to jedyny moduł bez okna w bocznej nawigacji.**

---

## 7. INWENTARZ OKIEN — pełny zakres opracowania

### 7.1. Okna przepływu głównego
| # | Okno | Stany / warianty |
|---|---|---|
| 1 | **Okno startowe (ładowania)** | Łączenie (domyślny) · Powrót z ważnym tokenem · Błąd połączenia |
| 2 | **Okno rejestracji i logowania** | Rejestracja · Konto niepotwierdzone · Logowanie · Logowanie z błędem · Odzyskiwanie konta krok 1 · Odzyskiwanie konta krok 2 |
| 3 | **Strona główna — Centrum dowodzenia** | Strefa 1 (4 karty środowisk) · Strefa 2 (4 kafle) · Strefa 3 (listwa 3 pozycji) |

### 7.2. Powłoki środowisk
| # | Powłoka | Zawartość |
|---|---|---|
| 4 | **Powłoka TalkIn** | pasek górny · pas kart sesji · boczna nawigacja (9) · obszar roboczy |
| 5 | **Powłoka WorkSpace** | jw. (9 pozycji) |
| 6 | **Powłoka CodeStudio** | jw. (8 pozycji) |
| 7 | **Powłoka MultitaskingAI** | pasek górny · pas kart sesji · **panel orkiestracji (6 sekcji)** · obszar roboczy |

### 7.3. Okna platformowe
| # | Okno | Zakres |
|---|---|---|
| 8 | **Okno Konfiguracji** | 13 zakresów + „Nakładka i wywołanie modelu" (14. obszar) + **Panel prowenancji** + modal „Pula kont Code CLI" + okno punktów izolacji |
| 9 | **Okno Ustawień** | konto Operatora · uwierzytelnianie · wygląd (motyw, język) · urządzenia i parowanie Mobile · powiadomienia · konta modeli |
| 10 | **Always On Display** | warstwa centralna, rdzeń awatara, monitor procesu |
| 11 | **Mobile** | funkcja globalna, widok mobilny |

### 7.4. Okno wspólne
| # | Okno | Uwaga |
|---|---|---|
| 12 | **Chat Window / okno komunikacji** | pas komunikacji wspólny WSZYSTKIM modułom; rekonfigurowany w kontekście każdego modułu |

### 7.5. Okna operacyjne modułów (poza Chat Window)
| Moduł | Okna (wiodące **pogrubione**) |
|---|---|
| **Studio** | **Studio Editor**, Tools Panel, Diff/Grep Panel, Session Repository, Preview Window |
| **Research** | **Research Workspace**, Sources Manager, Findings Panel, Report Builder, Export Panel |
| **Library** | **Library Explorer**, Tags & Collections, File Preview, Versioning Panel |
| **Translate** | Source Panel, **Translation Panels**, Glossary Manager |
| **Browser** | **Browser Window**, Sources Panel, Notes Panel |
| **Assistant** | **Voice Console**, Actions Monitor, Activity Feed |
| **Roundtable** | **Model Panels**, Debate Panel, Moderator Panel, Consensus Panel |
| **Workspace** | **Project Dashboard**, Instructions Panel, Context Memory, Project Library, Agent Manager |
| **Automations** | **Workflow Builder**, Scheduler, Queue Manager, Orchestrator, Execution Monitor |
| **Design** | **Design Board**, Assets Panel, Prompt Builder, Preview Window |
| **Apps** | **Product Builder**, Architecture Designer, Frontend Workspace, Backend Workspace, Deployment Panel |
| **Terminal** | **Terminal Tabs**, Output Console, Process Monitor |
| **Developer** | **Code Editor**, Project Tree, Git Panel, Build Output |
| **Diagnostics** | **Diagnostics Center**, Logs Viewer, Errors Panel, Recommendations Panel |
| **Agents** | **Agent Builder**, Model Configuration, Skills Manager, Connectors Manager, Permissions Center |

### 7.6. Okna robocze ról MultitaskingAI
**Executor 1** · **Executor 2** · **Coordinator** · **Executor 3 / Validator**
Makiety szczegółowe: **Executor Chat**, **Coordinator Chat**, **Results Analyzer**, okno konfiguracji punktów izolacji na poziomie „Rola".
Piąta pozycja zespołu: **Subagent Network**.
Panel orkiestracji — 6 sekcji. Silnik kolejek — **11 akcji**.

---

## 8. Zasada nadrzędna platformy — ZERO BLOKAD (ADL-017)

> **Danaco Console nie narzuca twardych blokad, bram bezpieczeństwa ani wymuszonych zgód. Domyślne zachowanie systemu to wykonanie polecenia.**

Izolacja techniczna, izolacja kontekstu, uprawnienia, profile izolacji i Permissions Center są **wyłącznie opcjami konfigurowalnymi** przez Operatora, ze stanem wyjściowym „wyłączone / pełny dostęp".

**Konsekwencje projektowe (bezwzględne):**
- **Nie stosuj atrybutu `disabled`.** Przycisk jest zawsze klikalny.
- Zamiast blokady: **komunikat po naciśnięciu** albo **opis obok**.
- Przycisk „Pomiń" w oknie logowania jest obecny i klikalny **w obu fazach wdrożenia**.
- Odliczanie przy „Wyślij ponownie" ma charakter **wyłącznie informacyjny** — nie blokuje kliknięcia.
- Metoda uwierzytelniania wyłączona w konfiguracji **nie jest w ogóle renderowana** („brak metody = mniej segmentów, nie zablokowany segment").

---

## 9. Dostępność — warunek wejściowy

- **WCAG 2.1 AA** jako warunek wejściowy, nie cel.
- Każda para tekst/tło **zmierzona** — komplet w `zasoby/zetony/kontrasty.json` (33 pomiary).
- `--dn-tekst-3` (#7C7C7C) — **wyłącznie metadane i tekst ≥ 18,66 px półgruby**.
- **Stan nigdy samym kolorem** — zawsze ikona albo etykieta.
- Pierścień fokusu: 2 px + odsunięcie 2 px, barwa `--dn-fokus`.
- `prefers-reduced-motion` obsługiwane **globalnie w żetonach** — nie per komponent.
- Cele dotykowe rosną żetonem (`pointer: coarse`), nie wyjątkiem.
- Nawigacja klawiaturą we wszystkich prototypach; `aria-*` i role semantyczne.
- Język dokumentu: `lang="pl"`.

---

## 10. Zasady redakcyjne opracowań

1. **Język:** polski. Komentarze i opisy po polsku; **identyfikatory kodu, klasy CSS i nazwy własne okien po angielsku** (zgodnie z dokumentacją: `Studio Editor`, `Workflow Builder`, `Chat Window`…).
2. **Nazwy własne modułów, okien, ról, sekcji — dokładnie jak w dokumentacji.** Zakaz tłumaczenia i parafrazowania.
3. **Treści przykładowe pochodzą z domeny produktu** i są oznaczone jako przykładowe. Zakaz: Lorem ipsum, „Jan Kowalski", zmyślonych metryk, fikcyjnych firm, wymyślonych nazwisk.
4. Przykładowe dane operacyjne budujemy ze świata produktu: nazwy sesji, identyfikatory zadań w kolejce, ścieżki repozytoriów, nazwy modeli, nazwy automatyk — spójne z modułem.
5. **Zakaz projektowania środowisk, komponentów i elementów, które nie mają bezpośredniego pokrycia w dokumentacji.** Jeśli dokumentacja nie rozstrzyga formy — rozstrzygnij ją zgodnie z katalogiem komponentów i odnotuj w sekcji „decyzje projektowe".
6. Każdy dokument MD: nagłówek metrykowy (produkt, wersja, status, data, odbiorcy, zakres) + stopka `*Danaco Console — AI Operating Environment · v2.0*` + `*© 2026 Danaco Holding Group Sp. z o.o.*`.

---

## 11. Standard techniczny prototypów HTML

**Każdy plik HTML prototypu MUSI:**

1. Być **samodzielny w działaniu** po otwarciu z dysku (`file://`) — bez serwera, bez budowania.
2. Zaczynać się od:
```html
<!doctype html>
<html lang="pl" data-theme="dark">
<head>
<meta charset="utf-8">
<meta name="viewport" content="width=device-width, initial-scale=1">
<title>Danaco Console — <nazwa okna></title>
<link rel="icon" href="<ścieżka>/zasoby/marka/favicon/favicon.svg" type="image/svg+xml">
<link rel="stylesheet" href="<ścieżka>/zasoby/zetony/fonty.css">
<link rel="stylesheet" href="<ścieżka>/zasoby/zetony/zetony.css">
<link rel="stylesheet" href="<ścieżka>/zasoby/css/fundament.css">
<link rel="stylesheet" href="<ścieżka>/zasoby/css/komponenty.css">
```
   i kończyć `<script src="<ścieżka>/zasoby/wspolne.js"></script>` (przełącznik motywu + zapas Invoker Commands).
3. Mieć **przełącznik motywu** (`data-przelacz-motyw`) na pasku górnym — oba motywy równoprawne, obowiązkowo działające.
4. **Być KLIKALNY i DYNAMICZNY.** Wymagane minimum interakcji:
   - przełączanie zakładek / paneli / sekcji bez przeładowania,
   - stany `:hover`, `:active`, `:focus-visible` widoczne i zgodne z żetonami,
   - otwieranie i zamykanie modali (`<dialog>`), dymków, menu kontekstowych,
   - przejścia stanów (ładowanie → wynik, praca → gotowe), **animowane** czasami `--dn-czas-*`,
   - **tętno kropki sygnału** tam, gdzie w tle biegnie praca,
   - przełączanie kart sesji, wybór pozycji bocznej nawigacji,
   - w oknach z listami — wybór wiersza, filtrowanie, sortowanie (jeśli okno je przewiduje).
5. **Cały JavaScript inline** w `<script>` na końcu pliku (poza `wspolne.js`). Bez zewnętrznych bibliotek, bez CDN, bez frameworków.
6. **Zakaz `localStorage`** poza tym, co robi `wspolne.js`.
7. Respektować `prefers-reduced-motion` (obsłużone żetonami — nie dubluj).
8. Zawierać **panel informacyjny „O tym opracowaniu"** (rozwijalny `<details>`, na dole): jakie okno przedstawia, źródło w dokumentacji, jakie interakcje są odwzorowane.
9. Używać **wyłącznie ikon z `zasoby/ikony/svg/`** — wklejanych inline jako `<svg>` (dla `file://` `<img>` też działa; preferuj inline dla `currentColor`).
10. **Zakaz `disabled`** — zgodnie z zasadą zero blokad.

### Wzorzec animacji (INTENSYWNOSC_RUCHU = 3/10)
- mikroprzejścia 100–220 ms, `--dn-ease`,
- brak animacji dekoracyjnych, brak parallaxu, brak scrollytellingu w oknach roboczych,
- **jeden ruch znaczący na widok** — tętno pracy w tle,
- przejścia widoków: `opacity` + `translateY(4px)`, maks. 220 ms,
- ruch zawsze niesie informację o stanie systemu.

---

## 12. Struktura repozytorium wynikowego

```
WYNIK/
├── KANON.md                       ← ten plik
├── INDEKS.html                    ← nawigacja po całości
├── zasoby/                        ← żetony, css, ikony, fonty, marka, wspolne.js
├── 01-dokumentacja-md/            ← opracowania merytoryczno-techniczne (.md)
├── 02-dokumentacja-html/          ← opracowania merytoryczno-graficzne (.html)
├── 03-marka/                      ← portfolio logotypu + księga znaku
├── 04-portfolio/                  ← portfolio design identity (tematyczne)
└── 05-okna/                       ← klikalne prototypy wszystkich okien
    ├── przeplyw/                  ← ładowanie, logowanie, centrum dowodzenia
    ├── srodowiska/                ← powłoki 4 środowisk
    ├── platformowe/               ← konfiguracja, ustawienia, AOD, Mobile
    └── moduly/                    ← okna operacyjne 15 modułów
```

---

## 13. Dokumentacja źródłowa — gdzie czytać

| Zakres | Ścieżka |
|---|---|
| Indeks projektu UI | `/home/claude/danaco/dok/projekt-ui/README.md` |
| Przepływ okien i stany | `/home/claude/danaco/dok/projekt-ui/przeplyw/przeplyw-okien.md` |
| Anatomia elementów powłoki | `/home/claude/danaco/dok/projekt-ui/przeplyw/elementy-okien.md` |
| Katalog komponentów | `/home/claude/danaco/dok/projekt-ui/komponenty/katalog-komponentow.md` |
| Moduły (14 plików) | `/home/claude/danaco/dok/projekt-ui/moduly/*.md` |
| Okno Konfiguracji | `/home/claude/danaco/dok/projekt-ui/okna/konfiguracja.md` |
| Okno Ustawień | `/home/claude/danaco/dok/projekt-ui/okna/ustawienia.md` |
| Moduł Agents | `/home/claude/danaco/dok/projekt-ui/okna/agents.md` |
| Środowisko MultitaskingAI | `/home/claude/danaco/dok/projekt-ui/srodowiska/multitaskingai.md` |
| Propozycje rozbudowy | `/home/claude/danaco/dok/propozycje-rozbudowy/` |
| Kierunek projektowy | `/home/claude/danaco/design/opracowania/design/01-kierunek/KIERUNEK.md` |
| Handoff wdrożeniowy | `/home/claude/danaco/design/opracowania/design/HANDOFF.md` |
| Istniejące makiety (wzorzec) | `/home/claude/danaco/design/opracowania/design/06-okna/*.html` |
| Istniejąca księga marki | `/home/claude/danaco/design/opracowania/design/07-ksiega/ksiega-marki.html` |
| Galeria ikon | `/home/claude/danaco/design/opracowania/design/05-ikony/ikony.html` |
| Pakiet systemu wizualnego | `/home/claude/danaco/design/opracowania/system-wizualny/pakiet/*.md` |
| Inwentarz okien (opracowany) | `/home/claude/danaco/brief/INWENTARZ-OKIEN.md` |
| Inwentarz komponentów (opracowany) | `/home/claude/danaco/brief/INWENTARZ-KOMPONENTOW.md` |

---

## 14. Kontrola jakości — lista sprawdzeń przed oddaniem pliku

- [ ] Wszystkie nazwy własne zgodne z dokumentacją (zero parafraz)
- [ ] Zero elementów bez pokrycia w dokumentacji
- [ ] Zero wartości szesnastkowych wpisanych wprost — tylko `var(--dn-*)`
- [ ] Zero `#000000`
- [ ] Zero `disabled`
- [ ] Zero emoji jako ikon
- [ ] Zero Lorem ipsum / zmyślonych nazwisk / zmyślonych metryk
- [ ] Oba motywy działają i wyglądają poprawnie
- [ ] Pasek górny atramentowy w obu motywach
- [ ] Stan nigdy samym kolorem (ikona/etykieta towarzyszy)
- [ ] Fokus widoczny na każdej kontrolce
- [ ] Interakcje działają po otwarciu z `file://`
- [ ] Panel „O tym opracowaniu" obecny
- [ ] Polskie diakrytyki poprawne w całym pliku

---

*Danaco Console — AI Operating Environment · v2.0*
*© 2026 Danaco Holding Group Sp. z o.o. — Dariusz Naharnowicz*

---

## 15. ANEKS v2.1 — WIĄŻĄCY (re-baseline pod Danaco Pilot, plansza 1–21)

> Zatwierdzony przez Właściciela 2026-08-14. **Zastępuje** wcześniejszy układ okien operacyjnych.
> Zasada: **zakres i komponenty funkcyjne + okna narzędziowe** wg planszy 1–21; **wygląd nasz** (monochrom + kropka, żetony `--dn-*`).

**Nowa warstwa układu (obowiązkowa dla wszystkich okien roboczych):**
- `zasoby/stanowisko.css` + `zasoby/stanowisko.js` — układ KOLUMNOWY (czaty w pionie + panele robocze), obszar polecenia, tryby, menu, dokowane panele. **Zastępuje** stackowanie góra/dół.
- Wzorzec odniesienia: `05-okna/WZORZEC-STANOWISKA.html`. Źródło zakresu: `KATALOG-FUNKCJI-DANACO-PILOT.md`.

**Obowiązkowe w każdym oknie roboczym:**
1. Powłoka `.sta-powloka` (pełny viewport, zero max-width) · pasek z hamburgerem (1), trybami **Talk/Work/Code** (0), okruszkami, szukajką, kontem `support@danaco-core.pl` (2).
2. Szyna `.szyna` z filtrem Stan/Sortowanie/Grupowanie (13), grupami projektów (16), „Więcej” (14).
3. `.sta-obszar`: strefa czatów w pionie (1/2/4) + strefa dokowanych paneli. Panele wg modułu, sterowane menu „⋮ → Panele” (3).
4. **Okno komunikacji** z pełnym **obszarem polecenia** (23): pasek kontekstu (Ten komputer · projekt · repo · worktree · diff · Zatwierdź zmiany) + wiersz sterowania: tryb uprawnień (5, domyślnie „Z pominięciem zgód”), „+” z konektorami/wtyczkami (9/15), mikrofon (8), model Fable 5/Opus 5/Sonnet 5/Haiku 4.5 (6), suwak wysiłku Ultra (7), nagrywanie, **↑ Do kolejki** (Work/Code) / **Wyślij** (Talk).
5. Dokowane panele wg katalogu (22): Plan, Zadania w tle, Artefakty, Pliki, Terminal, Kolejka, Subagenci, Przeglądarka (z trybem urządzenia 18, Serwery 19, ⋮ 17, ołówek 21).
6. Konfiguracja per okno (koło zębate w belce) pozostaje.

**Model nazewnictwa:** produkt = **Danaco Pilot**; modele: Fable 5, Opus 5, Sonnet 5, Haiku 4.5; e-mail: support@danaco-core.pl.

**Zakazy bez zmian:** zero `disabled`, zero `#000000`, zero hex w `<style>`, zero Lorem ipsum; „Z pominięciem zgód” jako stan domyślny (ADL-017).
