# Danaco Console — Moduł WorkSpace — dokumentacja projektowa

| | |
|---|---|
| **Produkt** | Danaco Console |
| **Rodzaj** | Platforma AI Workspace OS |
| **Opis** | Platforma jest wielośrodowiskowym systemem operacyjnym dla sztucznej inteligencji, integrującym komunikację, zarządzanie wiedzą, tworzenie treści, projektowanie, automatyzacje procesów oraz rozwój oprogramowania. |
| **Producent** | Danaco Holding Group Sp. z o.o. |
| **Twórca** | Dariusz Naharnowicz |
| **Wersja** | v2.0 |
| **Status** | Deweloperski |
| **Data** | 2026-08-06 |

**Informacje szczegółowe dokumentu:**

| | |
|---|---|
| **Moduł** | WorkSpace (rozdz. 11.2 Koncepcji platformy) |
| **Zakres** | Dokumentacja projektowa modułu: rola, okna operacyjne, elementy interfejsu, przepływy pracy, katalog funkcji i narzędzi, punkty sterowania z okna konfiguracji |
| **Środowiska dostępności** | TalkIn, WorkSpace, CodeStudio |
| **Forma udostępnienia** | Komponent własny (Projekt), konfigurowany w strefie 2 strony głównej, oraz okno modułowe w bocznej nawigacji trzech środowisk |
| **Odbiorcy dokumentu** | Designer (co, gdzie, w jakiej formie, do czego), Deweloper (co zbudować) |
| **Zasady wiążące** | Zero blokad w UI · klucze jawne · pełna konfigurowalność z okna konfiguracji · jawność zależności |
| **Źródło faktów** | architektura/koncepcja-platformy.md (rozdz. 2.2, 2.3, 3.2, 4, 5.2, 8, 9.2), specyfikacje/specyfikacja-modulow.md (4.2), specyfikacje/specyfikacja-okien-operacyjnych.md (6.2), interfejs-uzytkownika/strona-glowna-i-nawigacja.md (3.3), architektura/izolacja-i-zaleznosci.md (rozdz. 8, 13.3), architektura/model-konfiguracji.md (5.6, 5.9–5.13, 6), architektura/architektura.md (rdzeń Go, SQLite, WebSocket/JSON, klient Tauri), interfejs-uzytkownika/system-wizualny.md |

---

## Spis treści

1. Przeznaczenie i kontekst
2. Komplet okien operacyjnych modułu — przegląd
   - 2.1. Warstwy widoczności w module
   - 2.2. Kanoniczne rozmieszczenie okien modułu
3. Specyfikacja okien operacyjnych — pełny arsenał narzędzi
4. Katalog elementów interfejsu
5. Przepływy pracy w module
6. Stany, dane i powiązania z innymi modułami
7. Scenariusze użycia
8. Katalog funkcji i narzędzi
9. Punkty sterowania z okna konfiguracji

---

## 1. Przeznaczenie i kontekst

### 1.1. Rola modułu

WorkSpace jest **jednostką organizacji pracy projektowej** platformy: izolowaną przestrzenią, w której jeden projekt gromadzi w jednym miejscu swoje zadania, pliki, notatki, wiedzę, instrukcje, pamięć kontekstową i przypisanych wykonawców AI. Moduł stanowi kompletny system operacyjny pojedynczego przedsięwzięcia — użytkownik prowadzący projekt nie sięga po zewnętrzny menedżer zadań, notatnik, dysk w chmurze, wiki ani system zarządzania dokumentami.

Moduł jest jednym z elementów warstwy modułów (Workspace Layer), nie jej odpowiednikiem — nazwa modułu odnosi się wyłącznie do tej jednej, konkretnej, izolowanej przestrzeni projektowej.

### 1.2. Dla kogo

WorkSpace jest przestrzenią dla użytkowników prowadzących **równolegle wiele odrębnych projektów** — konsultantów, zespołów projektowych, kancelarii, deweloperów pracujących nad kilkoma zleceniami jednocześnie. Moduł adresuje sytuację, w której jedna, wspólna historia rozmowy z Wykonawcą i jedna, wspólna pamięć przestają wystarczać, bo ustalenia jednego przedsięwzięcia zaczynają przenikać do innego.

### 1.3. Po co

WorkSpace pozwala **skonfigurować rozdzielenie kontekstu, instrukcji i pamięci** między projektami — tak, aby ustalenia projektu A nie przenikały do projektu B, gdy taka izolacja jest pożądana, a jednocześnie nie stały na przeszkodzie, gdy w danej chwili pożądane jest ich świadome połączenie.

### 1.4. Obszar pokrycia modułu

| Obszar | Zakres |
|---|---|
| Przestrzeń robocza | Zakładanie, konfiguracja, duplikowanie, archiwizacja i szablonowanie projektów; wielość równoległych projektów z izolacją domyślną |
| Zarządzanie plikami | Repozytorium plików projektu: wgrywanie, foldery, tagi, kolekcje, podgląd, wersjonowanie, deduplikacja, wyszukiwanie pełnotekstowe, OCR |
| Zadania i planowanie | Lista, tablica kanban, oś czasu (Gantt), kalendarz, zależności, sprinty, przypomnienia, cele i kamienie milowe |
| Notatki i baza wiedzy | Notatki Markdown, dokumenty projektowe, wiki z odnośnikami wstecznymi, tablica wizualna, mapy myśli |
| Wiedza i pamięć dla AI | Pamięć kontekstowa projektu, wyszukiwanie semantyczne, wpisy ręczne i sugestie Wykonawcy, karty ustaleń |
| Instrukcje i szablony | Instrukcje systemowe projektu, biblioteka szablonów, snippety, zmienne dynamiczne, wersjonowanie |
| Wykonawcy AI | Przypisywanie agentów do projektu, zakres uprawnień, dziennik działań, domyślny wykonawca |
| Organizacja materiałów | Wyszukiwanie globalne, paleta poleceń, katalogowanie, powiązania między elementami projektu |

### 1.5. Granica tematyczna modułu

| Poza granicą | Właściwy moduł / miejsce | Uzasadnienie |
|---|---|---|
| Repozytorium centralne, wieloprojektowe | Library | Project Library jest odpowiednikiem ograniczonym do jednego projektu; magazyn wspólny należy do modułu Library |
| Zaawansowana edycja i przekształcanie treści dokumentów (korekta, styl, streszczenia) | Studio | WorkSpace przechowuje i wersjonuje dokument; operacje redakcyjne AI to natywny mechanizm Studio |
| Budowa i harmonogramowanie procesów automatycznych | Automations | WorkSpace *wiąże* projekt z gotową automatyką; silnik harmonogramów należy do Automations |
| Tworzenie i globalna konfiguracja agentów | Agents | WorkSpace *przypisuje* istniejących agentów; ich budowa i uprawnienia platformowe to moduł Agents |
| Kod źródłowy, kontrola wersji repozytoriów, budowanie | Developer / Terminal | Wersjonowanie w WorkSpace dotyczy dokumentów i notatek, nie kodu |
| Grafika, generowanie i edycja obrazów | Design | WorkSpace odbiera i przechowuje artefakty wizualne, nie generuje ich |
| Tłumaczenia równoległe, badania wieloźródłowe | Translate / Research | WorkSpace przechowuje wyniki; sam proces należy do modułów tematycznych |

Granica jest **kompozycyjna, nie zamknięta**: każdy z powyższych modułów wiąże się z WorkSpace jawnym powiązaniem konfigurowalnym (rozdz. 6.3). WorkSpace pełni rolę domu projektu, z którego materiał trafia do modułu wyspecjalizowanego i wraca do niego jako artefakt.

### 1.6. Charakter pracy

| Cecha | Opis |
|---|---|
| Jednostka organizacji | Projekt — jednostka grupująca instrukcje, pamięć kontekstową, bibliotekę i przypisanych agentów właściwe jednemu przedsięwzięciu |
| Punkt wyjścia | Zamiar prowadzenia wyodrębnionej pracy, a nie pojedyncza rozmowa |
| Wielość | Dowolna liczba projektów jednocześnie, każdy z własną kartą sesji lub wieloma kartami |
| Domyślne zachowanie | Każdy nowy projekt otrzymuje odrębny układ, historię, pamięć i kontekst; współdzielenie ustanawia się świadomą decyzją, nie regułą wymuszoną |

### 1.7. Miejsce modułu w architekturze platformy

```
STRONA GŁÓWNA — Centrum dowodzenia
├─ Strefa 1 · środowiska ──────────────► TalkIn · WorkSpace · CodeStudio · MultitaskingAI
├─ Strefa 2 · komponenty własne
│    kafel „WorkSpace" → „Załóż projekt" ──► okno konfiguracji → PROJEKT (komponent własny)
└─ Strefa 3 · ustawienia

                                    │
                                    ▼
        PROJEKT trafia do pamięci aplikacji jako zasób użytkownika
                                    │
                                    ▼
     wybieralny w bocznej nawigacji modułu WorkSpace w środowiskach:
              TalkIn · WorkSpace · CodeStudio
                                    │
                                    ▼
                    OKNA OPERACYJNE MODUŁU WORKSPACE
     Chat Window · Execution Loop Window · Project Dashboard ·
     Instructions Panel · Context Memory · Project Library · Agent Manager
```

WorkSpace jest jednym z dwóch modułów (obok Agents) dostępnych jednocześnie jako komponent własny **i** jako okno modułowe pełnoprawne we wszystkich trzech środowiskach modułowych — stąd podwójna droga wejścia widoczna na schemacie powyżej.

---

## 2. Komplet okien operacyjnych modułu — przegląd

| Okno | Rola w module | Forma wiodąca | Punkt wejścia? | Warstwa | Sposób wywołania |
|---|---|---|---|---|---|
| Chat Window | Główne okno komunikacji Użytkownik ↔ Wykonawca — centralny punkt pracy i podstawowy mechanizm sterowania procesami projektu | Lewa kolumna, stała, pełna wysokość obszaru roboczego | Nie — stale obecne | 1 | Widoczne bez interakcji |
| Execution Loop Window | Okno pętli wykonawczej Koordynator ↔ Wykonawca — dekompozycja zleceń projektu na zadania, nadzór i kontrola realizacji | Kolumna sąsiadująca z Chat Window, otwierana | Nie — stale dostępne | 2 | Kliknięcie zwiniętego wyzwalacza `Execution Loop ▸` w lewej kolumnie albo zlecenie wydane w Chat Window; okno otwiera się samoczynnie z chwilą przyjęcia zlecenia przez Koordynatora |
| Project Dashboard | Ogólny widok stanu i postępu projektu, hub planowania | Widok wiodący w obszarze roboczym, punkt wejścia do projektu | **Tak** | 1 | Widoczny bez interakcji — otwiera się po wybraniu projektu |
| Instructions Panel | Definiowanie instrukcji systemowych projektu | Okno edycyjne w obszarze roboczym | Nie | 2 | Kafel nawigacyjny `Instructions` na Project Dashboard, pozycja bocznej nawigacji albo paleta poleceń `Ctrl/Cmd + K` |
| Context Memory | Pamięć kontekstowa dedykowana projektowi | Okno listy wpisów w obszarze roboczym | Nie | 2 | Kafel nawigacyjny `Context Memory`, pozycja bocznej nawigacji albo paleta poleceń `Ctrl/Cmd + K` |
| Project Library | Repozytorium plików, notatek i wiki ograniczone do projektu | Okno eksploratora w obszarze roboczym | Nie | 2 | Kafel nawigacyjny `Project Library`, pozycja bocznej nawigacji albo paleta poleceń `Ctrl/Cmd + K` |
| Agent Manager | Przypisywanie agentów do projektu | Okno listy z akcjami w obszarze roboczym | Nie | 2 | Kafel nawigacyjny `Agent Manager`, pozycja bocznej nawigacji albo paleta poleceń `Ctrl/Cmd + K` |

Siedem okien łącznie. Chat Window i Execution Loop Window są oknami wspólnymi wszystkim modułom platformy (rozdz. 5 Specyfikacji okien operacyjnych) — poniżej opisane w zakresie rekonfiguracji właściwej WorkSpace; pozostałe pięć to okna właściwe wyłącznie temu modułowi.

### 2.1. Warstwy widoczności w module

Moduł WorkSpace stosuje zasadę nadrzędną interfejsu platformy: **jeżeli funkcja nie jest potrzebna do realizacji aktualnego zadania, nie jest widoczna**. Projekt gromadzi zadania, pliki, notatki, wiki, pamięć kontekstową, instrukcje i przypisanych wykonawców — bogactwo to istnieje w architekturze modułu i pozostaje niewidoczne w interfejsie do chwili wystąpienia potrzeby użycia danej funkcji. Liczba projektów, agentów, plików i wpisów pamięci nie wpływa na postrzeganą prostotę okna.

Każdy element interfejsu modułu należy do dokładnie jednej z czterech warstw widoczności:

| Warstwa | Zakres w module WorkSpace | Sposób dostępu |
|---|---|---|
| 1 — zawsze widoczna | Chat Window w lewej kolumnie; aktywne okno wiodące obszaru roboczego (domyślnie Project Dashboard: nagłówek projektu, kafle nawigacyjne okien modułu, lista zadań, pasek postępu); boczna nawigacja modułów i lista projektów; pasek kart sesji projektu; wskaźniki stanu wykonania — postęp zlecenia, przebieg pętli wykonawczej, zajętość pamięci i pojemności. Zajmuje ponad 80% powierzchni interfejsu modułu | Widoczna bez interakcji |
| 2 — widoczna na żądanie | Wybór wykonawcy zadania spośród agentów przypisanych do projektu, wybór warstwy zapisu instrukcji (projekt · sesja), profil izolacji projektu, zasięg wpisu pamięci i zasięg przypisania agenta, folder docelowy artefaktów, zakładki huba planowania (Metryki · Cele · Oś czasu), zakładki Project Library (Notatki i wiki · Tablica), pola wyszukiwania, panele podglądu pliku i uprawnień agenta, sterowanie przebiegiem pętli wykonawczej | Znacznik kontekstowy, ikona, przycisk lub przełącznik w pasku kontekstu okna; po użyciu element zwija się samoczynnie, a panel wysuwany po zamknięciu znika całkowicie z przestrzeni roboczej |
| 3 — rozwinięcia kontekstowe | Akcje wiadomości, akcje pliku, akcje wpisu pamięci i akcje agenta; operacje na projekcie (duplikuj · archiwizuj · eksportuj · usuń); akcje zbiorcze na zaznaczeniu; szablony, snippety, zmienne dynamiczne i wersje instrukcji; historia wersji i porównanie dokumentu; komunikaty sterujące, wyniki kontroli jakości i ponowienia na karcie zadania; graf wiedzy i odnośniki wsteczne; bloki osadzane notatki; eksporty i dzienniki | Menu kebab (`⋮`), menu hamburger (`☰`), menu kontekstowe, panel popover, lista rozwijana (`Szablony ▾`, `Wersje ▾`, `Filtr ▾`, `Operacje ▾`, `Sterowanie ▾`), znak `/` lub `@` w polu edycji |
| 4 — funkcje eksperckie | Nadpisania instrukcji per agent, powiązanie Project Library z modułem Library i wybiórcza synchronizacja repozytoriów, deduplikacja zbioru plików, ekstrakcja tekstu i OCR na zbiorze materiałów, przegląd i modyfikacja uprawnień agenta w Permissions Center, konfiguracja punktów izolacji projektu | Polecenie języka naturalnego w Chat Window, skrót klawiszowy, wyszukiwarka funkcji palety poleceń `Ctrl/Cmd + K`, tryb administracyjny albo konfiguracja roli. Użytkownik podstawowy nie widzi tych elementów |

Mechanizmy ukrywania funkcjonalności stosowane w module: menu progresywne (`Wykonawca ▾`, `Szablony ▾`), panele wysuwane (podgląd pliku, uprawnienia agenta — otwierane jako rozszerzenie boczne po prawej stronie obszaru roboczego), grupowanie logiczne akcji w element zbiorczy (`Operacje ▾`, `Sterowanie ▾`) oraz znaczniki kontekstowe w pasku kontekstu projektu, na przykład `[Projekt „Klient X"] [Zamknięty] [Agent „Analityk"] [warstwa: projekt]` — kliknięcie znacznika otwiera odpowiedni selektor.

**Zasada jednego kliknięcia.** Każda ukryta funkcja modułu jest osiągalna jednym kliknięciem, jednym skrótem klawiszowym albo jednym poleceniem języka naturalnego wydanym w Chat Window. Zagnieżdżanie funkcji głęboko w hierarchii menu jest zabronione — ukrycie zmniejsza chaos wizualny, nie utrudnia dostępu.

### 2.2. Kanoniczne rozmieszczenie okien modułu

Makiety w rozdz. 3 przedstawiają interfejs w stanie spoczynku: widoczne są wyłącznie elementy warstwy 1 oraz zwinięte wyzwalacze warstw 2–3 (`▼`, `⋮`, `☰`) i znaczniki kontekstowe.

```
 ═══════════════════════════════════════════════════════════════════════════
  Boczna    │ Chat Window          │ Obszar roboczy modułu    │ Panel
  nawigacja │ Użytkownik ↔         │ WorkSpace                │ pomocniczy
  modułów   │ Wykonawca            │ (Project Dashboard ·     │ (rozszerzenie
  i lista   │                      │  Instructions Panel ·    │  boczne:
  projektów │ ─────────────────    │  Context Memory ·        │  podgląd pliku,
            │ Execution Loop       │  Project Library ·       │  wersje,
            │ Koordynator ↔        │  Agent Manager)          │  uprawnienia)
            │ Wykonawca            │                          │
 ═══════════════════════════════════════════════════════════════════════════
```

---

## 3. Specyfikacja okien operacyjnych — pełny arsenał narzędzi

### 3.0. Konwencja opisu

Każde okno opisano: celem, zawartością źródłową, **pełnym arsenałem narzędzi** (funkcje maksymalnie oprzyrządowane — wymienione hojnie, nie oszczędnie), zachowaniem i stanami, oraz makietą tekstową układu. Wszystkie makiety przedstawiają układ pionowy (podział lewa–prawa); regulacji podlega wyłącznie szerokość kolumn. Stany wspólne wszystkim oknom platformy (Specyfikacja okien operacyjnych, rozdz. 3): trwałość stanu procesu sesji po rozłączeniu klienta, aktualizacja na żywo kanałem WebSocket dla okien monitorujących, objaśnienie kontekstowe `[?]` przy każdym elemencie konfiguracji, domyślna izolacja kontekstu per karta sesji ze współdzieleniem konfigurowalnym z okna konfiguracji punktów izolacji.

### 3.1. Chat Window (kanał Użytkownik ↔ Wykonawca)

| Pole | Treść |
|---|---|
| Cel | Główne okno komunikacji między Użytkownikiem a Wykonawcą (AI / agent / system wykonawczy). Stanowi centralny punkt pracy użytkownika w module i podstawowy mechanizm sterowania wszystkimi procesami projektu — poleceniami wydawanymi wobec instrukcji, pamięci, biblioteki, zadań i agentów tego projektu |
| Waga i miejsce | Lewa kolumna, stała, pełna wysokość obszaru roboczego; obecna w każdym oknie modułu w tym samym miejscu układu |
| Zawartość | Historia rozmowy projektu; pole wprowadzania poleceń w języku naturalnym; strumień odpowiedzi i wyników na żywo; sterowanie zatwierdzaniem i przerywaniem działań; pasek narzędzi kontekstowych |

**Pełny arsenał narzędzi:**

| Narzędzie | Funkcja | Warstwa | Sposób wywołania |
|---|---|---|---|
| Kompozytor wiadomości | Pole tekstowe wieloliniowe, formatowanie Markdown na żywo (pogrubienie, listy, blok kodu), skróty klawiszowe | 1 | Widoczny bez interakcji — pole poleceń w dolnej części lewej kolumny |
| Menedżer snippetów i szablonów poleceń | Zapisane, wielokrotnego użytku polecenia i fragmenty tekstu wstawiane skrótem | 3 | Lista rozwijana `Snippety ▾` w pasku kompozytora |
| Załączniki | Dołączanie plików wprost z Project Library lub z urządzenia; podgląd miniatur przed wysłaniem | 2 | Ikona `📎` w pasku kompozytora; po wstawieniu załącznika pasek wyboru zwija się |
| Wzmianki kontekstowe (`@`) | Odwołanie do pliku Project Library, wpisu Context Memory, zadania lub przypisanego agenta wprost w treści polecenia | 2 | Znak `@` wpisany w polu poleceń otwiera listę podpowiedzi |
| Wybór wykonawcy | Przełącznik modelu lub przypisanego agenta jako adresata polecenia (lista z Agent Managera) | 2 | Znacznik kontekstowy `[Wykonawca ▾]` w pasku kontekstu okna |
| Zlecenie realizacji zadania | Przekazanie polecenia do Koordynatora — zlecenie pojawia się w Execution Loop Window jako pozycja dekomponowana na zadania | 1 | Wysłanie polecenia z pola poleceń |
| Strumień odpowiedzi | Odbiór na żywo kanałem WebSocket, z przerwaniem generowania w toku | 1 | Widoczny bez interakcji — strumień w historii rozmowy |
| Zatwierdzanie i przerywanie działań | Akceptacja wyniku, żądanie korekty, przerwanie działania Wykonawcy z poziomu okna komunikacji | 1 | Widoczne bez interakcji w czasie działania Wykonawcy |
| Akcje wiadomości | Kopiuj, cytuj, edytuj i wyślij ponownie, rozgałęź wątek, przypnij do Context Memory, oznacz jako ustalenie projektu | 3 | Menu kebab (`⋮`) przy wiadomości |
| Blok kodu w odpowiedzi | Podświetlanie składni, kopiowanie, przekazanie do modułu docelowego (np. Developer) | 1 | Widoczny bez interakcji w treści odpowiedzi; akcje bloku w menu kebab (`⋮`) |
| Wyszukiwanie w historii | Przeszukiwanie treści dotychczasowej rozmowy projektu | 2 | Ikona lupy w nagłówku okna |
| Eksport rozmowy | Zapis wątku jako dokumentu (Markdown/PDF) do Project Library | 3 | Menu kebab (`⋮`) nagłówka okna |
| Wskaźnik zasięgu pamięci | Widoczna adnotacja, czy odpowiedź korzysta z pamięci projektu, czy z pamięci szerszego zasięgu (zgodnie z konfiguracją izolacji) | 1 | Widoczny bez interakcji — znacznik kontekstowy przy odpowiedzi |
| Wyjaśnienie wyniku i kontekstu | Rozwinięcie informacji o źródłach użytych przez Wykonawcę: wpisy pamięci, pliki, instrukcje | 3 | Rozwinięcie `Źródła ▾` pod odpowiedzią |

**Zachowanie i stany:** domyślnie odrębna historia i pamięć czatu per karta sesji projektu; ciągła, współdzielona historia między kartami tego samego projektu jest ustawieniem konfiguracyjnym okna konfiguracji punktów izolacji. Stan ładowania — animowany wskaźnik pisania podczas generowania odpowiedzi. Stan błędu — komunikat inline z ponowieniem.

```
 Makieta — Chat Window w module WorkSpace
 ═════════════════╤══════════════════════════════╤════════════════════════
  Boczna          │ Chat Window                  │ Obszar roboczy
  nawigacja       │ Użytkownik ↔ Wykonawca       │ modułu WorkSpace
  modułów         │                              │
  i projektów     │ Historia rozmowy projektu    │ (Project Dashboard
                  │ (przewijana, znaczniki czasu)│  lub inne okno
                  │                          [⋮] │  modułu)
                  │ [Projekt „Klient X"]         │
                  │ [Zamknięty] [pamięć projektu]│
                  │ [Agent „Redaktor" ▼]         │
                  │ ┌──────────────────────────┐ │
                  │ │ pole poleceń — Markdown  │ │
                  │ └──────────────────────────┘ │
                  │                       [ ⏎ ]  │
                  │ ───────────────────────────  │
                  │ Execution Loop ▸             │
 ═════════════════╧══════════════════════════════╧════════════════════════
```

### 3.2. Execution Loop Window (kanał Koordynator ↔ Wykonawca)

| Pole | Treść |
|---|---|
| Cel | Okno pętli wykonawczej prezentujące komunikację między Koordynatorem a Wykonawcą. Odpowiada za prowadzenie pętli wykonawczej, koordynację zadań projektu, nadzór nad realizacją, orkiestrację działań oraz kontrolę realizacji procesów |
| Waga i miejsce | Kolumna sąsiadująca z Chat Window, otwierana; pełna wysokość obszaru roboczego |
| Zawartość | Bieżące zlecenie projektu i jego dekompozycja na zadania, kolejka i stan zadań, wymiana komunikatów sterujących, wyniki kontroli jakości i decyzje o ponowieniu, wskaźniki przebiegu pętli, sterowanie przebiegiem |

**Pełny arsenał narzędzi:**

| Narzędzie | Funkcja | Warstwa | Sposób wywołania |
|---|---|---|---|
| Karta bieżącego zlecenia | Treść zlecenia przyjętego z Chat Window, projekt źródłowy, wskazany wykonawca, znacznik czasu przyjęcia | 1 | Widoczna bez interakcji na szczycie okna |
| Dekompozycja zlecenia na zadania | Drzewo zadań i podzadań wyprowadzonych przez Koordynatora ze zlecenia; każde zadanie z wykonawcą, terminem i kryterium ukończenia | 1 | Widoczna bez interakcji — drzewo zadań pod kartą zlecenia |
| Kolejka i stan zadań | Lista zadań w stanach: oczekujące, w realizacji, w kontroli, ukończone, ponawiane, przerwane | 1 | Widoczna bez interakcji |
| Wymiana komunikatów sterujących | Strumień komunikatów Koordynator → Wykonawca i Wykonawca → Koordynator: przydział zadania, raport postępu, zgłoszenie przeszkody, zwrot wyniku | 3 | Kliknięcie karty zadania rozwija strumień komunikatów |
| Wyniki kontroli jakości | Ocena rezultatu zadania wobec kryterium ukończenia, wraz z uzasadnieniem i decyzją o przyjęciu lub ponowieniu | 3 | Kliknięcie karty zadania rozwija wynik kontroli |
| Decyzje o ponowieniu | Rejestr ponowień zadania z licznikiem prób i zmienionymi parametrami zlecenia | 3 | Rozwinięcie `Ponowienia ▾` na karcie zadania |
| Wskaźniki przebiegu pętli | Liczba zadań w każdym stanie, czas trwania przebiegu, tempo realizacji, obłożenie przypisanych agentów | 1 | Widoczne bez interakcji — pasek wskaźników przebiegu |
| Sterowanie przebiegiem | Wstrzymanie, wznowienie, przerwanie pętli oraz korekta zlecenia bez utraty dotychczasowych wyników | 2 | Element zbiorczy `Sterowanie ▾` (wstrzymaj · wznów · przerwij · koryguj) |
| Powiązanie z zadaniami projektu | Zadania pętli wykonawczej zapisywane są jako pozycje listy zadań projektu (Project Dashboard) i widoczne w osi czasu aktywności | 1 | Widoczne bez interakcji — zadania pętli w liście zadań Project Dashboard |
| Odbiór artefaktów | Wyniki zadań kierowane do wskazanego folderu Project Library oraz ustalenia do Context Memory | 2 | Znacznik kontekstowy folderu docelowego; kliknięcie otwiera selektor |
| Dziennik przebiegu | Pełny, chronologiczny zapis pętli do podglądu i eksportu jako dokument projektu | 3 | Menu kebab (`⋮`) nagłówka okna |

**Zachowanie i stany:** okno otwierane jest jako kolumna sąsiadująca z głównym oknem komunikacji i aktualizuje się na żywo kanałem WebSocket. Stan przebiegu utrzymuje się po rozłączeniu klienta — po powrocie okno prezentuje bieżący stan pętli. Stan pusty — brak aktywnego zlecenia, widoczna lista zakończonych przebiegów projektu. Stan przeszkody — zadanie zgłoszone przez Wykonawcę jako zablokowane wyróżnione wizualnie, z akcją korekty zlecenia.

```
 Makieta — Execution Loop Window
 ═════════════════╤══════════════════════════════╤════════════════════════
  Boczna          │ Chat Window                  │ Obszar roboczy
  nawigacja       │ Użytkownik ↔ Wykonawca       │ modułu WorkSpace
  modułów         │ (zwinięte ▸)                 │
  i projektów     │ ───────────────────────────  │ Project Dashboard
                  │ Execution Loop               │ — zadania projektu
                  │ Koordynator ↔ Wykonawca      │   aktualizowane
                  │                              │   na żywo
                  │ Zlecenie: „Przygotuj raport  │
                  │ tygodniowy projektu"      [⋮]│
                  │ ├ 1 Zebranie danych   ✔    ▼ │
                  │ ├ 2 Analiza           ▶    ▼ │
                  │ ├ 3 Redakcja          ⏳   ▼ │
                  │ └ 4 Kontrola jakości  ⏳   ▼ │
                  │                              │
                  │ 2/4 zadań · czas 04:12       │
                  │ [Sterowanie ▼] [folder: /RAP]│
 ═════════════════╧══════════════════════════════╧════════════════════════
```

### 3.3. Project Dashboard

| Pole | Treść |
|---|---|
| Cel | Ogólny widok stanu i postępu projektu wraz z hubem planowania — punkt wejścia otwierany po wybraniu projektu |
| Zawartość | Zestawienie zadań, statusu prac, metryk i powiązanych zasobów projektu |

**Pełny arsenał narzędzi:**

| Narzędzie | Funkcja | Warstwa | Sposób wywołania |
|---|---|---|---|
| Nagłówek projektu | Nazwa, opis, data założenia, właściciel, plakietka statusu (aktywny / wstrzymany / zarchiwizowany) | 1 | Widoczny bez interakcji na szczycie obszaru roboczego |
| Kafle nawigacyjne | Skróty do pozostałych okien modułu, z licznikiem pozycji (np. „12 wpisów pamięci", „4 agentów") | 1 | Widoczne bez interakcji — siatka kafli okien modułu |
| Hub planowania | Zakładki Lista / Kanban / Gantt / Kalendarz nad wspólnym zbiorem zadań projektu | 1 | Widoczny bez interakcji — pasek zakładek nad zbiorem zadań |
| Lista/tablica zadań | Widok listy oraz widok tablicy kanban, grupowanie wg statusu, przypisanie wykonawcy (agent lub użytkownik) | 1 | Widoczna bez interakcji — widok wiodący huba planowania |
| Pasek postępu | Wizualizacja procentowa ukończenia projektu na podstawie zamkniętych zadań i kamieni milowych | 1 | Widoczny bez interakcji |
| Dashboard metryk | Kafle: zadania wg statusu, obłożenie wykonawców, tempo prac, zbliżające się terminy, wykres wypalania iteracji | 2 | Zakładka `Metryki` w pasku zakładek dolnej części widoku |
| Cele i kamienie milowe | Punkty kontrolne i cele mierzalne powiązane z zadaniami, z procentem realizacji | 2 | Zakładka `Cele` |
| Oś czasu aktywności | Chronologiczny zapis zdarzeń projektu: zmiany instrukcji, nowe wpisy pamięci, przebiegi pętli wykonawczej, przypisania agentów, aktualizacje biblioteki | 2 | Zakładka `Oś czasu` |
| Zasoby powiązane | Kafle skrótów do modułów powiązanych konfiguracyjnie (Library, Automations) wraz ze stanem powiązania | 3 | Rozwinięcie `Powiązania ▾` |
| Profil izolacji projektu | Widoczny znacznik przypisanego profilu izolacji, z odnośnikiem do okna konfiguracji punktów izolacji | 2 | Znacznik kontekstowy `[profil izolacji]` w nagłówku; kliknięcie otwiera selektor |
| Karty sesji projektu | Lista otwartych kart sesji pracujących w kontekście tego projektu, z przełączeniem | 1 | Widoczne bez interakcji — pasek kart sesji nad obszarem roboczym |
| Operacje na projekcie | Duplikuj projekt, zarchiwizuj, eksportuj podsumowanie, usuń — akcja wykonuje się od razu, dla usunięcia dostępny mechanizm „Cofnij" po wykonaniu; modal potwierdzenia to ustawienie konfiguracyjne Operatora | 3 | Menu kebab (`⋮`) przy nagłówku projektu |
| Notatka właściciela | Wolne pole opisowe — cel projektu, uwagi, przypomnienia | 3 | Rozwinięcie `Notatka ▾` w nagłówku projektu |

**Zachowanie i stany:** otwiera się jako punkt wejścia po wybraniu projektu z listy komponentów własnych lub bocznej nawigacji. Stan pusty — projekt świeżo założony pokazuje zachętę do wypełnienia Instructions Panel i Context Memory. Stan ładowania — szkielet kafli podczas pobierania danych.

```
 Makieta — Project Dashboard
 ═════════════════╤════════════════════╤═══════════════════════════════════
  Boczna          │ Chat Window        │ ◈ Nazwa projektu    [● aktywny][⋮]
  nawigacja       │ Użytkownik ↔       │   założony 2026-06-01
  modułów         │ Wykonawca          │   [Zamknięty]
  i projektów     │                    │ ─────────────────────────────────
                  │ historia rozmowy   │ ┌────────┬────────┬────────┬─────┐
                  │ projektu           │ │Instruc-│Context │Project │Agent│
                  │                    │ │tions   │Memory  │Library │Mgr  │
                  │ [pole poleceń]     │ │1 zestaw│12 wpis.│34 pliki│3 ag.│
                  │                    │ └────────┴────────┴────────┴─────┘
                  │ ─────────────────  │ Zadania [Lista ▼]
                  │ Execution Loop ▸   │ Postęp ▓▓▓▓▓▓▓░░░ 68%
                  │                    │ ☐ Zadanie 1 — Agent „Analityk"
                  │                    │ ☑ Zadanie 2 — użytkownik
                  │                    │ ─────────────────────────────────
                  │                    │ ☰ Metryki · Cele · Oś czasu
 ═════════════════╧════════════════════╧═══════════════════════════════════
```

### 3.4. Instructions Panel

| Pole | Treść |
|---|---|
| Cel | Definiowanie instrukcji systemowych właściwych projektowi |
| Zawartość | Treść instrukcji obowiązujących Wykonawcę w ramach danego projektu |

**Pełny arsenał narzędzi:**

| Narzędzie | Funkcja | Warstwa | Sposób wywołania |
|---|---|---|---|
| Edytor tekstu instrukcji | Edycja pełnoekranowa, Markdown, numerowanie linii, zawijanie | 1 | Widoczny bez interakcji — pole edycji zajmuje główną powierzchnię okna |
| Biblioteka szablonów instrukcji | Gotowe wzorce instrukcji do zastosowania i dalszej edycji (np. „redaktor techniczny", „analityk danych") | 3 | Lista rozwijana `Szablony ▾` |
| Zmienne i placeholdery | Wstawianie odwołań dynamicznych (np. `{{nazwa_projektu}}`, `{{data}}`) rozwijanych w czasie wykonania | 3 | Sekwencja `{{` w edytorze otwiera listę zmiennych |
| Wersjonowanie instrukcji | Lista kolejnych wersji z datą i autorem zmiany, porównanie (diff) dwóch wersji, przywrócenie wcześniejszej | 3 | Lista rozwijana `Wersje ▾` |
| Podgląd instrukcji efektywnej | Wynikowa treść instrukcji po uwzględnieniu dziedziczenia z warstwy szerszej (środowisko, globalna) — zgodnie z regułą „brak ustawienia = wartość domyślna" | 2 | Przycisk `Podgląd efektywny` w pasku okna |
| Test instrukcji | Wysłanie zapytania testowego do Chat Window z bieżącą wersją instrukcji, bez zapisu jako historia właściwa | 2 | Przycisk `Test w Chat` w pasku okna |
| Licznik długości | Liczba znaków i przybliżona liczba tokenów zajmowanych przez instrukcję | 1 | Widoczny bez interakcji — stopka okna edycji |
| Menedżer snippetów | Fragmenty instrukcji wielokrotnego użytku wstawiane skrótem | 3 | Lista rozwijana `Snippety ▾` |
| Import / eksport | Wczytanie instrukcji z pliku Project Library, zapis bieżącej treści jako plik | 3 | Menu kebab (`⋮`) nagłówka okna |
| Nadpisania per agent | Wskazanie, czy przypisany agent (Agent Manager) korzysta z instrukcji projektu wprost, czy z własnych instrukcji nadrzędnych (Agent Builder) | 4 | Polecenie języka naturalnego w Chat Window, wyszukiwarka funkcji albo konfiguracja roli — element niewidoczny dla użytkownika podstawowego |
| Warstwa konfiguracji | Przełącznik: instrukcja obowiązuje projekt (warstwa domyślna) czy wyłącznie bieżącą kartę sesji (warstwa sesji) | 2 | Przełącznik `Warstwa ▾` (projekt · sesja) w pasku kontekstu okna |

**Zachowanie i stany:** instrukcje obowiązują domyślnie w zakresie projektu; współdzielenie z innym projektem jest ustawieniem konfiguracyjnym okna konfiguracji punktów izolacji. Stan „niezapisane zmiany" sygnalizowany plakietką przy tytule okna. Stan pusty — projekt bez zdefiniowanych instrukcji dziedziczy wartość domyślną szerszej warstwy.

```
 Makieta — Instructions Panel
 ═════════════════╤════════════════════╤═══════════════════════════════════
  Boczna          │ Chat Window        │ Instrukcje systemowe — projekt
  nawigacja       │ Użytkownik ↔       │ „Raport tygodniowy" [niezapisane●]
  modułów         │ Wykonawca          │ [warstwa: projekt ▼]          [⋮]
  i projektów     │                    │ ─────────────────────────────────
                  │ historia rozmowy   │ ┌───────────────────────────────┐
                  │ projektu           │ │1 Jesteś asystentem redakcyj…  │
                  │                    │ │2 Zawsze stosuj ton formalny.  │
                  │ [pole poleceń]     │ │3 Odwołuj się do {{nazwa_pro…  │
                  │                    │ └───────────────────────────────┘
                  │ ─────────────────  │ 1 284 znaki · ~320 tokenów
                  │ Execution Loop ▸   │
                  │                    │
 ═════════════════╧════════════════════╧═══════════════════════════════════
```

### 3.5. Context Memory

| Pole | Treść |
|---|---|
| Cel | Pamięć kontekstowa dedykowana projektowi |
| Zawartość | Zapisane ustalenia, fakty i decyzje właściwe projektowi |

**Pełny arsenał narzędzi:**

| Narzędzie | Funkcja | Warstwa | Sposób wywołania |
|---|---|---|---|
| Lista wpisów pamięci | Karty pozycji z treścią, znacznikiem czasu i źródłem (sesja, użytkownik lub agent, który zapisał wpis) | 1 | Widoczna bez interakcji — treść wiodąca okna |
| Dodanie wpisu ręcznie | Formularz nowego ustalenia z polem treści, kategorią i zasięgiem | 1 | Przycisk `+ Dodaj wpis` — jedyny stale widoczny przycisk wiodący okna |
| Sugestie Wykonawcy | Wykonawca przedstawia zapis nowego ustalenia z toczącej się rozmowy; użytkownik akceptuje, edytuje lub odrzuca | 1 | Widoczne bez interakcji — wpis sugerowany pojawia się na liście wyróżniony |
| Wyszukiwanie semantyczne | Odnajdywanie wpisów znaczeniowo bliskich zapytaniu, nie tylko po słowach kluczowych | 2 | Ikona lupy w pasku okna |
| Kategoryzacja i tagi | Przypisanie własnych etykiet do wpisów (np. „decyzja", „dane liczbowe", „kontakt") | 3 | Menu kebab (`⋮`) karty wpisu |
| Wyszukiwanie i filtrowanie | Przeszukiwanie treści wpisów, filtr po tagu, dacie lub źródle | 2 | Element zbiorczy `Filtr ▾` (tag · data · źródło) |
| Przypinanie | Oznaczenie wpisów priorytetowych, wyświetlanych zawsze na początku listy i zawsze obecnych w kontekście Wykonawcy | 3 | Menu kebab (`⋮`) karty wpisu |
| Edycja i usuwanie wpisu | Modyfikacja treści istniejącego wpisu; usunięcie pojedynczego wpisu | 3 | Menu kebab (`⋮`) karty wpisu |
| Scalanie wpisów | Połączenie kilku powiązanych wpisów w jeden, spójny zapis; usunięcie sprzeczności | 3 | Element zbiorczy `Operacje ▾` po zaznaczeniu wielu wpisów |
| Poziom zasięgu wpisu | Wskazanie, czy wpis obowiązuje na poziomie projektu, czy — decyzją użytkownika — został podniesiony do zasięgu szerszego (środowisko, globalny) | 2 | Znacznik kontekstowy zasięgu na karcie wpisu; kliknięcie otwiera selektor |
| Wizualizacja zajętości pamięci | Wskaźnik wykorzystania dostępnej pojemności pamięci projektu | 1 | Widoczna bez interakcji — wskaźnik w rogu panelu |
| Eksport pamięci | Zapis całości lub zaznaczonych wpisów jako dokument do Project Library | 3 | Menu kebab (`⋮`) nagłówka okna |

**Zachowanie i stany:** domyślnie odrębna per projekt; zakres współdzielenia z innym projektem lub środowiskiem konfigurowalny z okna konfiguracji punktów izolacji. Stan „sugestia oczekująca" — wpis przedstawiony przez Wykonawcę wyróżniony wizualnie do chwili decyzji użytkownika.

```
 Makieta — Context Memory
 ═════════════════╤════════════════════╤═══════════════════════════════════
  Boczna          │ Chat Window        │ Pamięć projektu (12 wpisów)
  nawigacja       │ Użytkownik ↔       │ [+ Dodaj wpis] [🔍][Filtr ▼]  [⋮]
  modułów         │ Wykonawca          │ Zajętość: ▓▓▓░░ 34%
  i projektów     │                    │ ─────────────────────────────────
                  │ historia rozmowy   │ 📌 „Budżet ustalony na 120 000 zł."
                  │ projektu           │    #decyzja · 2026-07-02        ⋮
                  │                    │ ⬡ „Klient preferuje kontakt
                  │ [pole poleceń]     │    mailowy." — sugestia Wykonawcy
                  │                    │    [Akceptuj|Edytuj|Odrzuć]    ⋮
                  │ ─────────────────  │ ⬡ „Termin oddania: 30.09."     ⋮
                  │ Execution Loop ▸   │
                  │                    │
 ═════════════════╧════════════════════╧═══════════════════════════════════
```

### 3.6. Project Library

| Pole | Treść |
|---|---|
| Cel | Repozytorium plików, notatek i wiki ograniczone do projektu — odpowiednik modułu Library w zakresie jednego projektu, obejmujący pełny system zarządzania dokumentami i bazy wiedzy |
| Zawartość | Dokumenty, artefakty, notatki, strony wiki i tablice właściwe jednemu projektowi |

**Zakładki okna:** Pliki · Notatki i wiki · Tablica.

**Pełny arsenał narzędzi:**

| Narzędzie | Funkcja | Warstwa | Sposób wywołania |
|---|---|---|---|
| Eksplorator plików | Widok siatki oraz widok listy, przełącznik, struktura folderów zagnieżdżonych, sortowanie, kolumny metadanych | 1 | Widoczny bez interakcji — treść wiodąca okna |
| Przeciągnij i upuść | Wgrywanie plików z urządzenia bezpośrednio na obszar eksploratora, wiele naraz, pasek postępu, wznawianie | 1 | Widoczne bez interakcji — obszar eksploratora przyjmuje upuszczone pliki |
| Tagi i kolekcje | Nadawanie etykiet, grupowanie w kolekcje tematyczne niezależne od folderów | 3 | Menu kebab (`⋮`) pozycji albo element zbiorczy `Tagi ▾` |
| Podgląd pliku | Podgląd zawartości bez opuszczania modułu — dokumenty, obrazy, PDF, arkusze, audio, wideo | 2 | Kliknięcie pozycji otwiera panel podglądu jako rozszerzenie boczne; zamknięcie usuwa panel z przestrzeni roboczej |
| Ekstrakcja tekstu i OCR | Wydobycie treści z PDF, DOCX i TXT oraz rozpoznanie tekstu z obrazów i skanów; treść trafia do indeksu wyszukiwania | 4 | Polecenie języka naturalnego w Chat Window albo wyszukiwarka funkcji |
| Wyszukiwanie pełnotekstowe | Przeszukiwanie nazw i treści zgromadzonych materiałów z podświetleniem trafień | 2 | Ikona lupy w pasku okna |
| Panel wersjonowania i porównanie | Historia kolejnych wersji dokumentu, w tym wersji wygenerowanych przez Wykonawcę w innych oknach projektu; porównanie linia po linii, przywrócenie | 3 | Menu kebab (`⋮`) pozycji → `Wersje` |
| Deduplikacja | Wykrywanie identycznych plików po skrócie treści i scalanie pozycji | 4 | Polecenie języka naturalnego w Chat Window albo wyszukiwarka funkcji |
| Akcje zbiorcze | Zaznaczenie wielu plików — pobieranie, przenoszenie, usuwanie, tagowanie grupowe, spakowanie do archiwum | 3 | Element zbiorczy `Operacje ▾` po zaznaczeniu wielu pozycji |
| Odbiór artefaktów | Wyniki pracy Wykonawcy z Chat Window, Execution Loop Window i powiązanych modułów trafiają automatycznie do wskazanego folderu | 1 | Widoczny bez interakcji — artefakt pojawia się w folderze docelowym |
| Udostępnienie do modułu zewnętrznego | Wysłanie pliku jako załącznika do sesji innego modułu (np. Design, Studio, Research) | 3 | Menu kebab (`⋮`) pozycji → `Wyślij do modułu` |
| Edytor notatek Markdown | Pełny edytor: nagłówki, listy, tabele, bloki kodu z podświetlaniem składni, cytaty, listy kontrolne | 2 | Zakładka `Notatki i wiki` |
| Strony i drzewo dokumentów | Hierarchia stron zagnieżdżonych z drzewem nawigacji i automatycznym spisem treści notatki | 2 | Zakładka `Notatki i wiki` — drzewo stron w kolumnie zakładki |
| Odnośniki wsteczne i graf wiedzy | Wiązanie stron zapisem `[[nazwa]]`, panel „co linkuje tutaj", wizualizacja sieci powiązań notatek, zadań i plików | 3 | Panel `Co linkuje tutaj` oraz element zbiorczy `Graf ▾` |
| Bloki osadzane | Wstawienie do notatki żywego widoku zadania, pliku, wpisu pamięci lub grafiki | 3 | Znak `/` w edytorze notatki otwiera listę bloków |
| Szablony notatek i adnotacje | Wzorce (protokół spotkania, notatka decyzyjna, brief), zaznaczenia fragmentów, komentarze marginesowe, kolory wyróżnień | 3 | Lista rozwijana `Szablony ▾` w edytorze notatki |
| Tabele osadzane | Lekka tabela danych z sortowaniem i filtrem wewnątrz notatki, z eksportem do arkusza | 3 | Znak `/` w edytorze notatki → `Tabela` |
| Zakładki i wycinki webowe | Zapis odnośnika z metadanymi, miniaturą i notatką; odbiór wycinków z modułu Browser | 3 | Element zbiorczy `Operacje ▾` zakładki `Notatki i wiki` |
| Dziennik projektu | Nota dzienna z agregacją zdarzeń projektu | 2 | Zakładka `Notatki i wiki` → `Dziennik projektu` |
| Tablica wizualna i mapy myśli | Nieskończone płótno z kartkami, strzałkami i grupami, osadzaniem zadań i plików; rozgałęziona struktura węzłów z konwersją do listy zadań | 2 | Zakładka `Tablica` |
| Wskaźnik pojemności | Wykorzystanie limitu pamięci masowej przypisanego projektowi; ostrzeżenia miękkie, bez blokady | 1 | Widoczny bez interakcji — wskaźnik w rogu panelu |
| Powiązanie z Library | Odnośnik konfiguracyjny do modułu Library — wybiórcza synchronizacja materiałów z repozytorium centralnym | 4 | Okno konfiguracji punktów izolacji albo polecenie języka naturalnego w Chat Window |

**Zachowanie i stany:** stanowi zakres ograniczony do jednego projektu w odróżnieniu od modułu Library (repozytorium centralne, wieloprojektowe). Stan pusty — zachęta do wgrania pierwszego pliku lub utworzenia pierwszej notatki. Stan synchronizacji — wskaźnik trwającego wgrywania i indeksowania.

```
 Makieta — Project Library
 ═════════════════╤════════════════════╤═══════════════════════════════════
  Boczna          │ Chat Window        │ Biblioteka projektu (34 pliki)
  nawigacja       │ Użytkownik ↔       │ [Pliki ▼] [+ Wgraj]
  modułów         │ Wykonawca          │ [🔍] [Operacje ▼]             [⋮]
  i projektów     │                    │ ─────────────────────────────────
                  │ historia rozmowy   │ 🗎 umowa_v3.docx                ⋮
                  │ projektu           │ 🖼 baner.png                    ⋮
                  │                    │ 🗎 raport_final.pdf             ⋮
                  │ [pole poleceń]     │ 🗒 notatka-decyzja.md           ⋮
                  │                    │ ─────────────────────────────────
                  │ ─────────────────  │ Pojemność ▓▓░░ 22%
                  │ Execution Loop ▸   │
                  │                    │
 ═════════════════╧════════════════════╧═══════════════════════════════════
```

### 3.7. Agent Manager

| Pole | Treść |
|---|---|
| Cel | Przypisywanie agentów do projektu |
| Zawartość | Lista agentów skonfigurowanych w module Agents, dostępnych do przypisania |

**Pełny arsenał narzędzi:**

| Narzędzie | Funkcja | Warstwa | Sposób wywołania |
|---|---|---|---|
| Lista dostępnych agentów | Wszyscy agenci utworzeni w module Agents, z awatarem, nazwą i statusem | 1 | Widoczna bez interakcji — treść wiodąca okna |
| Przypisanie do projektu | Dodanie agenta jako Wykonawcy dostępnego w zakresie projektu, przez przeciągnięcie lub przycisk | 1 | Przycisk `+ Przypisz agenta` — jedyny stale widoczny przycisk wiodący okna |
| Zakres uprawnień w projekcie | Wskazanie, do których okien projektu przypisany agent ma dostęp (Chat Window, Execution Loop Window, Project Library, Context Memory); klucze jawne | 2 | Kliknięcie pozycji agenta otwiera panel uprawnień jako rozszerzenie boczne |
| Rola domyślnego wykonawcy | Oznaczenie jednego z przypisanych agentów jako Wykonawcy sugerowanego domyślnie w Chat Window i w pętli wykonawczej | 3 | Menu kebab (`⋮`) pozycji agenta → `Ustaw jako domyślny` |
| Podgląd aktywności agenta | Dziennik działań agenta wykonanych w kontekście tego projektu | 3 | Menu kebab (`⋮`) pozycji agenta → `Dziennik działań` |
| Skrót do Agent Builder | Utworzenie nowego agenta bez opuszczania kontekstu projektu | 2 | Przycisk `+ Nowy agent` w stopce listy |
| Skrót do Permissions Center | Przegląd i modyfikacja uprawnień agenta na poziomie platformowym | 4 | Tryb administracyjny albo polecenie języka naturalnego w Chat Window — element niewidoczny dla użytkownika podstawowego |
| Odłączenie agenta | Usunięcie przypisania agenta do projektu (nie usuwa agenta jako komponentu własnego) | 3 | Menu kebab (`⋮`) pozycji agenta |
| Filtr i wyszukiwanie | Odnalezienie agenta po nazwie, umiejętności lub podłączonym konektorze | 2 | Ikona lupy w pasku okna |
| Zasięg przypisania | Wskaźnik, czy przypisanie obowiązuje wyłącznie ten projekt, czy — decyzją użytkownika — także inne, równolegle prowadzone projekty | 2 | Znacznik kontekstowy zasięgu przy pozycji agenta; kliknięcie otwiera selektor |
| Obłożenie wykonawców | Rozkład zadań i czasu pracy przypadający na poszczególnych agentów projektu | 1 | Widoczne bez interakcji — wskaźnik w stopce okna |

**Zachowanie i stany:** agenci przypisani domyślnie działają w zakresie tego projektu; wpływ na inne, równolegle prowadzone projekty użytkownik określa w oknie konfiguracji. Stan „brak przypisanych agentów" — zachęta do przejścia do modułu Agents.

```
 Makieta — Agent Manager
 ═════════════════╤════════════════════╤═══════════════════════════════════
  Boczna          │ Chat Window        │ Agenci przypisani (3)
  nawigacja       │ Użytkownik ↔       │ [+ Przypisz agenta] [🔍]      [⋮]
  modułów         │ Wykonawca          │ ─────────────────────────────────
  i projektów     │                    │ ⬡ „Analityk" ●domyślny         ⋮
                  │ historia rozmowy   │ ⬡ „Redaktor"                   ⋮
                  │ projektu           │ ⬡ „Weryfikator"                ⋮
                  │                    │ ─────────────────────────────────
                  │ [pole poleceń]     │ [+ Nowy agent]
                  │ ─────────────────  │ Obłożenie ▓▓▓░ 61%
                  │ Execution Loop ▸   │
                  │                    │
 ═════════════════╧════════════════════╧═══════════════════════════════════
```

**Powierzchnie przekrojowe** (dostępne z każdego okna modułu): wyszukiwanie globalne projektu, paleta poleceń `Ctrl/Cmd + K`, objaśnienia kontekstowe `[?]`, strumień powiadomień projektu.

---

## 4. Katalog elementów interfejsu

Zestawienie kluczowych elementów interfejsu modułu, element po elemencie, wg schematu: co to jest · do czego służy · forma i waga · stany · zachowanie po interakcji · gdzie występuje. Klasy komponentów odwołują się do biblioteki `components.css` systemu wizualnego marki (prefiks `.dn-`).

| Element | Co to jest | Do czego służy | Forma i waga | Stany | Zachowanie po interakcji | Gdzie występuje | Warstwa | Sposób wywołania |
|---|---|---|---|---|---|---|---|---|
| Kolumna Chat Window | Stała kolumna komunikacji `.dn-kolumna--komunikacja` | Główne okno komunikacji Użytkownik ↔ Wykonawca; sterowanie procesami modułu | Lewa kolumna, stała, pełna wysokość obszaru roboczego; regulowana szerokość | domyślny, zwężony, zwinięty do paska ikon | Regulacja szerokości uchwytem; zwinięcie nie przerywa strumienia | Wszystkie okna modułu | 1 | Widoczna bez interakcji |
| Kolumna Execution Loop Window | Kolumna pętli wykonawczej `.dn-kolumna--petla` | Prezentacja komunikacji Koordynator ↔ Wykonawca i sterowanie przebiegiem | Kolumna sąsiadująca z Chat Window, otwierana, pełna wysokość | zamknięta, otwarta, aktywny przebieg (obrys akcentowy), przeszkoda (obrys ostrzegawczy) | Otwarcie odsłania kolejkę zadań; sterowanie przebiegiem działa natychmiast | Wszystkie okna modułu | 2 | Kliknięcie zwiniętego wyzwalacza `Execution Loop ▸`; kolumna otwiera się samoczynnie z chwilą przyjęcia zlecenia |
| Karta zadania pętli | `.dn-karta--zadanie` z ikoną stanu | Prezentacja jednego zadania dekomponowanego ze zlecenia | Wiersz drzewa zadań | oczekujące, w realizacji, w kontroli, ukończone, ponawiane, przerwane | Kliknięcie rozwija komunikaty sterujące i wynik kontroli jakości | Execution Loop Window | 1 | Widoczna bez interakcji w drzewie zadań; komunikaty sterujące i wynik kontroli w warstwie 3 — kliknięcie karty |
| Kafel nawigacyjny okna | Karta `.dn-karta--interaktywna` z licznikiem | Skrót do pozostałych okien modułu | Średni panel, siatka 4 kolumn w obszarze roboczym | domyślny, hover (tło), aktywny (obrys złoty), ładowanie (szkielet) | Otwiera docelowe okno w tej samej karcie sesji | Project Dashboard | 1 | Widoczny bez interakcji |
| Plakietka statusu projektu | `.dn-plakietka--stan` z kropką koloru | Sygnalizacja stanu projektu (aktywny/wstrzymany/zarchiwizowany) | Mała pigułka, obok nagłówka | sukces (aktywny), ostrzeżenie (wstrzymany), neutralny (zarchiwizowany) | Kliknięcie otwiera menu zmiany statusu | Project Dashboard | 1 | Widoczna bez interakcji; menu zmiany statusu w warstwie 3 — kliknięcie plakietki |
| Przycisk „+ Dodaj wpis" | `.dn-btn--zloty`, rozmiar `sm` | Jedyny, najważniejszy CTA okna Context Memory | Mały przycisk wypełniony, cień złoty | domyślny, hover (uniesienie), aktywny (wciśnięcie) | Zawsze klikalny — otwiera formularz nowego wpisu inline; próba zapisania pustej treści sygnalizowana ostrzeżeniem, nie blokadą | Context Memory | 1 | Widoczny bez interakcji — jedyny stale widoczny przycisk wiodący okna |
| Karta wpisu pamięci | `.dn-karta--pozycja` | Prezentacja pojedynczego ustalenia | Wiersz listy, rozdzielony kreską | domyślny, przypięty (ikona pinezki), sugestia oczekująca (obrys akcentowy) | Kliknięcie rozwija treść i akcje (edytuj/usuń/przypnij) | Context Memory | 1 | Widoczna bez interakcji; akcje wpisu w warstwie 3 — menu kebab (`⋮`) |
| Przełącznik warstwy konfiguracji | `.dn-zakladki--pigulki`, dwie pozycje | Wybór, czy zmiana dotyczy projektu, czy tylko bieżącej sesji | Mały segmentowany przełącznik | projekt (domyślny), sesja | Zmienia zasięg zapisu kolejnej edycji | Instructions Panel | 2 | Znacznik `[warstwa: projekt ▼]` w pasku kontekstu okna; po wyborze zwija się samoczynnie |
| Ikona `[?]` objaśnienia kontekstowego | `.dn-tooltip`, wyzwalacz ikonowy 14 px | Wyjaśnienie działania ustawienia i jego wpływu na aplikację | Bardzo mała ikonka | domyślny, najechanie/fokus (dymek widoczny) | Wyświetla dymek z treścią objaśnienia | Wszystkie okna z ustawieniami | 2 | Ikona wyzwalacza; dymek pojawia się po najechaniu lub uzyskaniu fokusu i znika po odsunięciu |
| Przycisk „Test w Chat" | `.dn-btn--zarys`, rozmiar `sm` | Akcja drugorzędna — wysłanie zapytania testowego z bieżącą instrukcją | Mały przycisk z obrysem | domyślny, hover | Zawsze klikalny — kieruje zapytanie testowe do Chat Window, nie zapisuje jako historia właściwa; przy pustej treści instrukcji dołączone ostrzeżenie o braku treści zamiast blokady | Instructions Panel | 2 | Przycisk `Test w Chat` w pasku okna |
| Awatar agenta | `.dn-awatar--kwadrat` | Identyfikacja wizualna agenta (nie-osobowa) | Mały kwadrat z inicjałem lub ikoną | domyślny, ze wskaźnikiem statusu online (`.dn-awatar-stan`) | Kliknięcie otwiera podgląd szczegółów agenta | Agent Manager, Chat Window (wzmianki), Execution Loop Window | 1 | Widoczny bez interakcji przy pozycji agenta; szczegóły agenta w warstwie 2 — kliknięcie otwiera panel boczny |
| Plakietka „domyślny wykonawca" | `.dn-plakietka--zloto` | Oznaczenie agenta sugerowanego domyślnie | Bardzo mała pigułka | aktywna (widoczna) / nieaktywna (ukryta) | Kliknięcie „Ustaw jako domyślny" w menu agenta przełącza oznaczenie | Agent Manager | 1 | Widoczna bez interakcji przy agencie domyślnym; przełączenie w warstwie 3 — menu kebab (`⋮`) |
| Pasek postępu projektu | Pasek wypełnienia poziomego | Wizualizacja procentowego ukończenia zadań projektu | Cienki pasek pełnej szerokości panelu | 0–100%, kolor złoty wypełnienia | Aktualizuje się automatycznie po zmianie statusu zadania | Project Dashboard | 1 | Widoczny bez interakcji |
| Pole wgrywania plików | Strefa przeciągnij-i-upuść z obrysem przerywanym | Wgranie nowego pliku do Project Library | Duży obszar w pustym stanie eksploratora | pusty (zachęta), najechanie pliku (podświetlenie), wgrywanie (pasek postępu) | Upuszczenie pliku rozpoczyna wgrywanie i dodaje pozycję do listy | Project Library | 1 | Widoczne bez interakcji — obszar eksploratora przyjmuje upuszczone pliki |
| Menu kontekstowe pliku (`⋯`) | Przycisk ikonowy `.dn-btn-ikona` 36×36 px | Dostęp do akcji na pojedynczym pliku (podgląd, wersje, tagi, usuń, wyślij do modułu) | Mała ikona | domyślny, hover, otwarte (rozwinięta lista) | Otwiera listę akcji kontekstowych | Project Library | 3 | Menu kebab (`⋮`) przy pozycji pliku |
| Modal potwierdzenia | `.dn-modal` z nakładką `.dn-nakladka` | Dodatkowe sygnalizowanie wagi akcji nieodwracalnej (np. usunięcie projektu); ustawienie konfiguracyjne Operatora, nie bramka domyślna | Duży panel wyśrodkowany | otwarty / zamknięty | Domyślnie akcja wykonuje się od razu, z dostępnym „Cofnij" po wykonaniu; gdy Operator włączy to ustawienie, modal poprzedza wykonanie | Project Dashboard (usuń, zarchiwizuj) | 4 | Ustawienie konfiguracyjne Operatora w oknie konfiguracji — element niewidoczny dla użytkownika podstawowego |
| Pole wyszukiwania | `.dn-input` z ikoną lupy | Filtrowanie listy (plików, wpisów pamięci, agentów, zadań) | Pole tekstowe, pełna szerokość kolumny | pusty, wpisywanie, z wynikami, brak wyników | Filtruje listę na żywo podczas wpisywania | Context Memory, Project Library, Agent Manager, Project Dashboard | 2 | Ikona lupy w pasku okna; pole zwija się po wyczyszczeniu zapytania |
| Paleta poleceń | `.dn-paleta`, wywołanie `Ctrl/Cmd + K` | Przejście do dowolnego okna, elementu lub akcji z klawiatury | Nakładka wyśrodkowana z polem i listą trafień | zamknięta, otwarta, wyniki rozmyte, brak wyników | Wybór pozycji wykonuje akcję lub otwiera element | Wszystkie okna modułu | 4 | Skrót klawiszowy `Ctrl/Cmd + K` — wyszukiwarka funkcji obejmująca również operacje eksperckie |
| Wskaźnik zajętości (pamięć/pojemność) | Mały pasek lub pierścień procentowy | Orientacyjna informacja o wykorzystaniu limitu | Bardzo mały element, róg panelu | poniżej progu / blisko limitu (kolor ostrzegawczy) | Najechanie pokazuje wartość liczbową w dymku | Context Memory, Project Library | 1 | Widoczny bez interakcji — wskaźnik w rogu panelu; wartość liczbowa w dymku po najechaniu |
| Karta sesji projektu | Element paska kart środowiska | Reprezentacja jednej równoległej przestrzeni roboczej projektu | Mały prostokąt w pasku kart nad obszarem roboczym | aktywna, w tle (wskaźnik pracy), zamykana | Kliknięcie przełącza widok obszaru roboczego | Pasek kart sesji | 1 | Widoczna bez interakcji — pasek kart sesji nad obszarem roboczym |

---

## 5. Przepływy pracy w module

### 5.1. Założenie nowego projektu i pierwsza konfiguracja

```
Strona główna, strefa 2
        │  kafel „WorkSpace” → „Załóż projekt”
        ▼
Okno konfiguracji komponentu własnego (szablon: nazwa, przeznaczenie, widoczność)
        │  zapis
        ▼
PROJEKT trafia do pamięci aplikacji jako zasób użytkownika
        │
        ▼
Otwarcie modułu WorkSpace w wybranym środowisku → Project Dashboard (stan pusty)
        │
        ├──► Instructions Panel   — zdefiniowanie instrukcji systemowych
        ├──► Context Memory       — pierwsze ustalenia
        ├──► Project Library      — wgranie materiałów wejściowych
        └──► Agent Manager        — przypisanie wykonawców
                │
                ▼
        Praca właściwa przez Chat Window, w kontekście
        w pełni skonfigurowanego projektu
```

### 5.2. Zlecenie realizowane w pętli wykonawczej

```
Chat Window (Użytkownik → Wykonawca)
        │  polecenie w języku naturalnym: „przygotuj raport tygodniowy”
        ▼
Koordynator przyjmuje zlecenie i dekomponuje je na zadania
        │
        ▼
Execution Loop Window (Koordynator ↔ Wykonawca)
   ├─ przydział zadań przypisanym agentom projektu
   ├─ raporty postępu i zgłoszenia przeszkód
   ├─ kontrola jakości wobec kryterium ukończenia
   └─ decyzja: przyjęcie wyniku albo ponowienie zadania
        │
        ▼
Zadania widoczne w Project Dashboard · artefakty w Project Library ·
ustalenia w Context Memory · zapis przebiegu w osi czasu aktywności
        │
        ▼
Chat Window prezentuje wynik; Użytkownik zatwierdza albo koryguje zlecenie
```

### 5.3. Codzienna praca równoległa nad wieloma projektami

```
[ Karta sesji A · WorkSpace · Projekt „Klient X” ]
[ Karta sesji B · WorkSpace · Projekt „Klient Y” ]   [ + nowa karta ]
        │                          │
   własna historia,           własna historia,
   pamięć, biblioteka         pamięć, biblioteka
   (odrębne domyślnie — rozdz. 6 Koncepcji platformy)
        │                          │
        ▼                          ▼
  Instructions Panel A       Instructions Panel B
  Context Memory A           Context Memory B
  pętla wykonawcza A         pętla wykonawcza B
        (izolacja domyślna; współdzielenie
         konfigurowalne z okna konfiguracji
         punktów izolacji)
```

### 5.4. Wykorzystanie agenta jako Wykonawcy zadania w projekcie

```
Moduł Agents (strefa 2 lub okno modułowe)
        │  utworzenie i konfiguracja agenta
        ▼
Agent zapisany jako komponent własny (zasób platformowy)
        │
        ▼
WorkSpace → Agent Manager → „+ Przypisz agenta”
        │
        ▼
Agent widoczny w Chat Window projektu jako wybieralny Wykonawca
i w Execution Loop Window jako adresat przydzielanych zadań
        │
        ▼
Polecenie wydane w Chat Window → Koordynator przydziela zadanie agentowi,
który realizuje je w zakresie uprawnień z Permissions Center i w kontekście
instrukcji, pamięci i biblioteki tego projektu
```

### 5.5. Świadome połączenie kontekstu dwóch projektów

```
Stan domyślny:  Projekt A ─╳─ Projekt B     (kontekst odrębny)
                              │
              okno konfiguracji punktów izolacji
              poziom zasięgu: „Para modułów” lub „Projekt”
                              │
                              ▼
Decyzja świadoma: „współdziel pamięć projektu A z projektem B”
                              │
                              ▼
Stan po zmianie:  Projekt A ─────► Projekt B   (pamięć współdzielona,
                                                 historia nadal odrębna,
                                                 zgodnie z wybranym zakresem)
```

---

## 6. Stany, dane i powiązania z innymi modułami

### 6.1. Dane wykorzystywane przez Wykonawcę w module

| Źródło danych | Okno pochodzenia | Uwaga |
|---|---|---|
| Instrukcje systemowe projektu | Instructions Panel | Obowiązują domyślnie w zakresie projektu |
| Pamięć kontekstowa projektu | Context Memory | Domyślnie odrębna per projekt |
| Materiały projektu | Project Library | Repozytorium ograniczone do projektu |
| Lista i konfiguracja przypisanych agentów | Agent Manager | Determinuje dostępnych Wykonawców w Chat Window i w pętli wykonawczej |
| Zlecenia, zadania i wyniki kontroli jakości | Execution Loop Window | Stan pętli wykonawczej projektu, źródło zapisów osi czasu |
| Stan i postęp projektu | Project Dashboard | Kontekst orientacyjny, nie treść instrukcji |

### 6.2. Stany projektu jako komponentu własnego

| Stan | Znaczenie | Gdzie widoczny |
|---|---|---|
| Aktywny | Projekt w bieżącym użyciu, dostępny do pracy | Plakietka w Project Dashboard, na kaflu strefy 2 |
| Wstrzymany | Projekt tymczasowo nieaktywny, dane zachowane | Plakietka ostrzegawcza |
| Zarchiwizowany | Projekt zamknięty, dostępny do podglądu, nieaktywny operacyjnie | Plakietka neutralna, ograniczony zestaw akcji |
| Współdzielony (częściowo) | Co najmniej jeden zakres kontekstu połączony z innym projektem lub środowiskiem | Znacznik profilu izolacji w Project Dashboard |

### 6.3. Powiązania konfigurowalne z innymi modułami

```
                     ┌───────────────────────────────┐
   Agents  ────────► │                                │
   (Agent Manager)   │           WORKSPACE             │
                      │        (moduł, ta karta)        │
   Library  ◄──────── │  Project Library — odpowiednik   │
   (repozytorium      │  ograniczony do jednego projektu │
    centralne)        │                                │
                      │                                │ ────────► Automations
                      └───────────────────────────────┘   (powiązanie projektu
                                                             z procesem automatycznym)
```

| Moduł/funkcja docelowa | Charakter powiązania | Typ | Miejsce ustanowienia |
|---|---|---|---|
| Agents | Agenci skonfigurowani w module Agents udostępniani — decyzją użytkownika — jako Wykonawcy zadań projektu | Konfiguracyjne | Agent Manager |
| Library | Project Library jako odpowiednik modułu Library ograniczony do zakresu projektu; wybiórcza synchronizacja | Konfiguracyjne | Project Library, okno konfiguracji |
| Automations | Projekty prowadzone w WorkSpace wiąże się z procesami automatycznymi z modułu Automations | Konfiguracyjne | Strona główna, strefa 2 (Automations); Project Dashboard |
| Studio, Design, Research, Translate | Materiał projektu przekazywany do modułu wyspecjalizowanego i odbierany z powrotem jako artefakt | Konfiguracyjne | Project Library (udostępnienie do modułu) |
| Developer / Terminal | Bloki kodu z notatek i odpowiedzi Wykonawcy przekazywane do modułu właściwego kodowi | Konfiguracyjne | Chat Window, Project Library |
| Browser | Odbiór wycinków webowych i zakładek do notatek projektu | Konfiguracyjne | Project Library (zakładka Notatki i wiki) |
| Okno konfiguracji punktów izolacji | Przypisanie profilu izolacji do projektu; poziom zasięgu „Projekt” w selektorze zasięgu | Globalny mechanizm platformy | Atrybut projektu „przypisany profil izolacji” |

Wszystkie powyższe powiązania są jawne — stan wyjściowy nowego projektu to pełne rozdzielenie kontekstu od innych projektów, przy jednoczesnym braku aktywnej izolacji technicznej procesu (żaden z ośmiu zakresów technicznych nie jest domyślnie włączony). Każde z powiązań włącza się świadomą decyzją z okna konfiguracji lub z poziomu okna modułu.

---

## 7. Scenariusze użycia

### 7.1. Kancelaria prowadząca równolegle kilka spraw klienckich

1. Użytkownik zakłada osobny projekt WorkSpace dla każdej sprawy klienckiej (strefa 2 strony głównej).
2. Dla każdego projektu definiuje odrębne instrukcje systemowe w Instructions Panel (np. rejestr terminologii właściwy danej sprawie).
3. Ustalenia i fakty specyficzne dla sprawy trafiają do Context Memory — ręcznie albo jako sugestie Wykonawcy akceptowane w toku rozmowy.
4. Dokumenty sprawy gromadzone są w Project Library, z zachowaniem wersji, rozpoznaniem tekstu ze skanów i wyszukiwaniem pełnotekstowym.
5. Dedykowany agent „Redaktor pism" przypisany przez Agent Manager wspiera przygotowanie pism w kontekście wyłącznie tej jednej sprawy — kontekst innej sprawy pozostaje niewidoczny, zgodnie z domyślną izolacją.
6. Zlecenie „przygotuj projekt pisma" wydane w Chat Window trafia do Koordynatora; Execution Loop Window prezentuje dekompozycję na zadania, kontrolę jakości i decyzje o ponowieniu.
7. Karty sesji poszczególnych spraw pozostają otwarte równolegle; przełączanie między nimi nie miesza ich pamięci ani historii.

### 7.2. Zespół produktowy łączący świadomie dwa powiązane projekty

1. Dwa projekty — „Produkt — badania" i „Produkt — wdrożenie" — prowadzone są odrębnie.
2. Gdy wnioski badawcze mają zasilić prace wdrożeniowe, użytkownik otwiera okno konfiguracji punktów izolacji, wybiera poziom zasięgu „Para modułów"/„Projekt" i włącza współdzielenie pamięci między oboma projektami.
3. Context Memory obu projektów prezentuje wspólny zbiór ustaleń oznaczonych źródłem pochodzenia.
4. Po zakończeniu prac współdzielenie zostaje wyłączone — Context Memory wraca do pełnej odrębności, bez utraty już zapisanych wpisów.

### 7.3. Automatyzacja cyklicznego raportu projektowego

1. W projekcie WorkSpace „Raport tygodniowy" gromadzone są dane źródłowe w Project Library.
2. W module Automations (strefa 2 strony głównej) budowana jest automatyka pobierająca dane z tego projektu i generująca raport wg ustalonego harmonogramu.
3. Powiązanie projektu z automatyką ustanawiane jest jako konfiguracyjne, widoczne w Project Dashboard jako aktywny zasób powiązany.
4. Koordynator prowadzi realizację raportu w pętli wykonawczej; przebieg widoczny jest w Execution Loop Window wraz z kontrolą jakości wyniku.
5. Wygenerowany raport trafia do Project Library jako nowa wersja dokumentu, widoczna w panelu wersjonowania tego pliku.

### 7.4. Migracja pracy z narzędzi zewnętrznych do jednego projektu

1. Użytkownik wczytuje zbiór notatek Markdown wraz ze strukturą folderów oraz eksport tablicy zadań w formacie JSON.
2. Materiał trafia do Project Library jako strony wiki z odnośnikami wstecznymi oraz do listy zadań projektu.
3. Kalendarz `.ics` wczytany do projektu zasila zadania terminami i pozycjami kalendarza.
4. Wyszukiwanie pełnotekstowe i semantyczne obejmuje cały wczytany materiał; paleta poleceń daje dostęp do dowolnego elementu z klawiatury.

---

## 8. Katalog funkcji i narzędzi

Konwencja pozycji: **nazwa** · co robi · zależności (biblioteki Go rdzenia / formaty / integracje wewnątrzplatformowe). Wszystkie funkcje działają bez twardych blokad — brak konfiguracji oznacza wartość domyślną, nie przerwanie pracy. Brak skonfigurowanej biblioteki oznacza wyłączenie danej funkcji z komunikatem informacyjnym, nie przerwanie pracy modułu. Zależności wymagające CGo (`gosseract`, `go-fitz`) oraz o licencji komercyjnej (`unioffice`) przyjmuje się świadomą decyzją — dla każdej wskazano wariant czysto-Go jako alternatywę, zgodnie z zasadą jawności zależności.

### Grupa A — Projekty i przestrzeń robocza

| # | Nazwa | Co robi | Zależności |
|---|---|---|---|
| A1 | Kreator projektu | Zakłada projekt z nazwą, opisem, ikoną, kolorem, celem i profilem izolacji; projekt jest operacyjny od razu, przed uzupełnieniem konfiguracji | SQLite (encja `project`), `google/uuid`; komunikat `project.create` (MODEL-KONFIGURACJI 5.13) |
| A2 | Szablony projektów | Zakłada projekt z gotowego wzorca (zestaw zadań, folderów, instrukcji, tagów) — np. „sprawa kliencka", „raport cykliczny", „badanie" | JSON szablonu w SQLite; `text/template` do rozwinięcia zmiennych |
| A3 | Duplikacja projektu | Klonuje projekt w całości lub wybiórczo (struktura bez treści / z treścią) | Transakcja SQLite; kopiowanie plików treści |
| A4 | Archiwizacja i wznowienie | Przenosi projekt w stan „zarchiwizowany" (podgląd, bez pracy) i z powrotem, bez utraty danych; „Cofnij" po akcji | Pole stanu w SQLite; brak twardego usuwania |
| A5 | Miękkie usuwanie z „Cofnij" | Usuwa projekt z natychmiastowym mechanizmem odtworzenia; modal potwierdzenia to ustawienie konfiguracyjne Operatora, nie bramka | Kolejka „soft-delete" w SQLite z TTL; `robfig/cron` do czyszczenia |
| A6 | Karty sesji projektu | Prowadzi wiele równoległych kart sesji jednego projektu; przełączanie bez mieszania historii i pamięci | Izolacja per karta (MODEL-KONFIGURACJI 5.12); stan trwały (ARCHITEKTURA 3.x) |
| A7 | Profil izolacji projektu | Przypisuje projektowi profil rozdzielenia kontekstu, pamięci i historii (od „zamknięty" po „współdzielony") | Okno konfiguracji punktów izolacji (MODEL-KONFIGURACJI 6) |
| A8 | Eksport podsumowania projektu | Generuje jednoplikowe podsumowanie stanu (zadania, pliki, ustalenia, oś czasu) do Markdown/PDF | `yuin/goldmark` (MD→HTML), `go-fitz`/maroto (PDF); trafia do Project Library |

### Grupa B — Zadania, planowanie i harmonogram

| # | Nazwa | Co robi | Zależności |
|---|---|---|---|
| B1 | Widok listy zadań | Lista zadań z polami: tytuł, status, wykonawca, priorytet, termin, etykiety, podzadania | SQLite (encja `task`); WebSocket na żywo |
| B2 | Tablica kanban | Kolumny wg statusu, przeciąganie kart, limity WIP, swimlane wg wykonawcy | Model kolumn w JSON; kolejność przez pole rank |
| B3 | Oś czasu / Gantt | Wykres słupkowy zadań w czasie, przeciąganie ram, zależności „koniec-początek" | Obliczenia ścieżki na rdzeniu Go; render po stronie klienta |
| B4 | Kalendarz projektu | Zadania i kamienie milowe w siatce miesiąc/tydzień/dzień; eksport i import iCal | `arran4/golang-ical` (format `.ics`) |
| B5 | Zależności zadań | Definiuje relacje poprzednik/następnik, blokady, automatyczne przesunięcia terminów | Graf zależności w SQLite; walidacja cykli w Go |
| B6 | Podzadania i listy kontrolne | Rozbicie zadania na kroki z paskiem postępu; szablony list kontrolnych | Encja `subtask`; agregacja postępu |
| B7 | Priorytety i macierz ważności | Etykiety pilności i ważności, widok macierzy 2×2 | Pola priorytetu; widok filtrowany |
| B8 | Iteracje | Grupuje zadania w okresy z celem iteracji, prędkością i przeglądem | Encja `sprint`; liczniki agregujące |
| B9 | Kamienie milowe i cele | Definiuje punkty kontrolne i cele mierzalne powiązane z zadaniami | Encja `milestone`; procent ukończenia |
| B10 | Przypomnienia i terminy | Powiadomienia o zbliżających się i przekroczonych terminach; kanałem WebSocket i push | `robfig/cron`; kolejka powiadomień; komunikat `notify` |
| B11 | Zadania cykliczne | Definiuje powtarzalność (dzienna, tygodniowa, reguła RRULE) | `arran4/golang-ical` (RRULE); scheduler Go |
| B12 | Przypisanie wykonawcy (agent/użytkownik) | Wskazuje agenta z Agent Managera lub użytkownika jako odpowiedzialnego; zlecenie realizacji wprost z karty zadania | Powiązanie z Agent Manager; Chat Window jako kanał zlecenia, Execution Loop Window jako kanał nadzoru |
| B13 | Śledzenie czasu | Rejestruje czas pracy nad zadaniem (ręcznie lub stoperem); raport sumaryczny | Encja `time_entry`; agregacja w Go |
| B14 | Automatyzacje tablicy | Reguły „gdy status = X, wykonaj Y" (przenieś, przypisz, oznacz), lokalnie lub przez powiązanie z Automations | Silnik reguł Go; powiązanie z modułem Automations |

### Grupa C — Notatki, wiki i baza wiedzy

| # | Nazwa | Co robi | Zależności |
|---|---|---|---|
| C1 | Edytor notatek Markdown | Pełny edytor WYSIWYG/Markdown: nagłówki, listy, tabele, bloki kodu, cytaty, listy kontrolne | `yuin/goldmark` (+ rozszerzenia GFM); zapis jako plik treści |
| C2 | Dokumenty projektowe (strony) | Hierarchia stron zagnieżdżonych z drzewem nawigacji; strona jako notatka długa | Drzewo w SQLite; pliki treści `.md` |
| C3 | Odnośniki wsteczne (backlinks) | Automatyczne wiązanie stron przez `[[nazwa]]`; panel „co linkuje tutaj" | Parser wikilinków w Go; graf w SQLite |
| C4 | Graf wiedzy | Wizualizuje sieć powiązań notatek, zadań i plików projektu | Graf z backlinków; render po stronie klienta |
| C5 | Tablica wizualna | Nieskończone płótno z kartkami, strzałkami, grupami; osadzanie zadań i plików | Model sceny JSON w SQLite; render na kliencie |
| C6 | Mapy myśli | Rozgałęziona struktura węzłów wokół tematu; konwersja do listy zadań | Struktura drzewa JSON; eksport do zadań B1 |
| C7 | Bloki osadzane | Wstawia do notatki żywy widok zadania, pliku, wpisu pamięci lub grafiki | System referencji `@`; render żywy przez WebSocket |
| C8 | Szablony notatek | Wielokrotne wzorce (protokół spotkania, notatka decyzyjna, brief) z polami | `text/template`; biblioteka szablonów w SQLite |
| C9 | Adnotacje i wyróżnienia | Zaznaczanie fragmentów, komentarze marginesowe, kolory wyróżnień | Zakresy tekstu w SQLite; warstwa nakładki na kliencie |
| C10 | Bloki kodu z podświetlaniem | Notatki techniczne z kolorowaniem składni i kopiowaniem; przekazanie do Developer | `alecthomas/chroma`; powiązanie z Developer |
| C11 | Osadzanie tabel i arkuszy | Lekka tabela danych z sortowaniem i filtrem wewnątrz notatki | Model tabeli JSON; eksport CSV `encoding/csv` |
| C12 | Zakładki i wycinki webowe | Zapis odnośnika z metadanymi, miniaturą i notatką; odbiór wycinków z modułu Browser | `mimetype`; powiązanie z Browser (Notes/Sources Panel) |
| C13 | Dziennik / notatki dzienne | Nota per dzień z agregacją zdarzeń projektu | Generacja z osi czasu; `robfig/cron` |
| C14 | Spis treści i nawigacja notatki | Automatyczny spis nagłówków z przewijaniem; zwijanie sekcji | Parser nagłówków `goldmark` (AST) |

### Grupa D — Pliki, repozytorium i zarządzanie dokumentami

| # | Nazwa | Co robi | Zależności |
|---|---|---|---|
| D1 | Eksplorator plików (siatka/lista) | Przegląda pliki projektu w folderach, dwa widoki, sortowanie, kolumny metadanych | SQLite (metadane) + pliki treści; `spf13/afero` (abstrakcja FS) |
| D2 | Przeciągnij i upuść / wgrywanie | Wgrywa pliki z urządzenia, wiele naraz, pasek postępu, wznawianie | Chunkowany transfer WebSocket; zapis pliku treści |
| D3 | Foldery i struktura | Tworzy zagnieżdżone foldery, przenosi, zmienia nazwy | Drzewo folderów w SQLite |
| D4 | Tagi i kolekcje | Nadaje etykiety wielokrotne i grupuje pliki w kolekcje tematyczne niezależne od folderów | Relacja tag↔plik w SQLite |
| D5 | Podgląd plików w oknie | Renderuje bez pobierania: obrazy, PDF, DOCX, TXT/MD, CSV, audio, wideo | `go-fitz` (PDF→obraz), `mimetype`, `disintegration/imaging` (obrazy) |
| D6 | Ekstrakcja tekstu (PDF/DOCX/TXT) | Wydobywa treść tekstową dokumentu do indeksu i podglądu | PDF: `ledongthuc/pdf` (czysty Go) albo `gen2brain/go-fitz` (MuPDF, CGo); DOCX: `nguyenthenguyen/docx` albo `unidoc/unioffice` |
| D7 | OCR obrazów i skanów | Rozpoznaje tekst z obrazów i skanów PDF, dokłada do indeksu wyszukiwania | `otiai10/gosseract` (Tesseract), języki `pol`+`eng` |
| D8 | Wyszukiwanie pełnotekstowe | Przeszukuje nazwy i treść wszystkich materiałów projektu, z podświetleniem trafień | SQLite FTS5 jako mechanizm domyślny (spójny z jednym plikiem danych); `blevesearch/bleve` przy bogatszych zapytaniach i fasetach; indeks inkrementalny |
| D9 | Wersjonowanie dokumentów | Historia wersji pliku (ręcznych i generowanych przez Wykonawcę), porównanie, przywrócenie | Snapshoty treści; `sergi/go-diff` |
| D10 | Porównanie wersji (diff) | Pokazuje różnice między dwiema wersjami dokumentu tekstowego linia po linii | `sergi/go-diff` (diffmatchpatch) albo `hexops/gotextdiff` (format unified) |
| D11 | Deduplikacja i wykrywanie duplikatów | Wykrywa identyczne pliki po skrócie treści, przedstawia scalenie | `crypto/sha256`; indeks skrótów w SQLite |
| D12 | Akcje zbiorcze | Zaznacza wiele plików: pobierz, przenieś, usuń, taguj, spakuj do ZIP | `archive/zip` (stdlib); transakcja SQLite |
| D13 | Odbiór artefaktów | Wyniki pracy Wykonawcy z Chat Window, Execution Loop Window i powiązanych modułów trafiają automatycznie do wskazanego folderu | Reguła kierowania w konfiguracji; komunikat `artifact.created` |
| D14 | Udostępnienie do modułu zewnętrznego | Wysyła plik jako załącznik do sesji Studio, Design, Research itd. | Powiązania konfigurowalne (SPECYFIKACJA-MODULOW 4.2) |
| D15 | Synchronizacja z Library | Wybiórcza synchronizacja plików projektu z repozytorium centralnym Library | Powiązanie Project Library ↔ Library; kolejka sync |
| D16 | Wskaźnik pojemności i limity | Pokazuje wykorzystanie limitu pamięci masowej projektu; miękkie ostrzeżenia, nie blokada | Agregacja rozmiarów w SQLite; próg z konfiguracji |

### Grupa E — Pamięć kontekstowa i wiedza dla Wykonawcy

| # | Nazwa | Co robi | Zależności |
|---|---|---|---|
| E1 | Lista wpisów pamięci | Karty ustaleń z treścią, znacznikiem czasu, źródłem (sesja/użytkownik/agent) | Encja `memory_entry` w SQLite |
| E2 | Ręczne dodawanie ustaleń | Formularz nowego faktu lub decyzji z kategorią i zasięgiem | SQLite; pola kategorii |
| E3 | Sugestie Wykonawcy z rozmowy | Wykonawca przedstawia zapis nowego ustalenia z toczącej się rozmowy; akceptacja, edycja albo odrzucenie | Analiza wątku w rdzeniu; stan „sugestia oczekująca" |
| E4 | Wyszukiwanie semantyczne | Znajduje wpisy znaczeniowo bliskie zapytaniu, nie tylko po słowach kluczowych | Adapter osadzeń (INTEGRACJA-MODELI) + `asg017/sqlite-vec` (rozszerzenie SQLite, spójne z jednym plikiem danych) albo `philippgille/chromem-go` (baza wektorowa w czystym Go) |
| E5 | Kategoryzacja i tagi wpisów | Etykiety własne („decyzja", „liczba", „kontakt"), filtr po tagu, dacie i źródle | Relacja tag↔wpis w SQLite |
| E6 | Przypinanie priorytetowych | Utrzymuje kluczowe ustalenia na początku listy i zawsze w kontekście Wykonawcy | Pole „pinned"; kolejność |
| E7 | Scalanie i porządkowanie wpisów | Łączy powiązane wpisy w jeden spójny zapis; usuwa sprzeczności | Operacja merge; historia zmian wpisu |
| E8 | Poziom zasięgu wpisu | Podnosi wpis z poziomu projektu do środowiska lub poziomu globalnego (świadoma decyzja) | Warstwowość pamięci (MODEL-KONFIGURACJI 5.10, 6.2) |
| E9 | Wskaźnik zajętości pamięci i eksport | Pokazuje wykorzystanie pojemności; eksportuje pamięć do dokumentu w Project Library | Agregacja; `goldmark` do eksportu MD |

### Grupa F — Instrukcje, polecenia i szablony

| # | Nazwa | Co robi | Zależności |
|---|---|---|---|
| F1 | Edytor instrukcji systemowych | Pełnoekranowa edycja instrukcji projektu: Markdown, numeracja linii, zawijanie | `goldmark`; zapis na warstwie projektu (MODEL-KONFIGURACJI 5.6) |
| F2 | Biblioteka szablonów instrukcji | Gotowe wzorce ról („redaktor techniczny", „analityk") do zastosowania i edycji | Szablony w SQLite |
| F3 | Zmienne i placeholdery | Odwołania dynamiczne `{{nazwa_projektu}}`, `{{data}}` rozwijane w czasie wykonania | `text/template` |
| F4 | Wersjonowanie instrukcji | Lista wersji z autorem i datą, porównanie (diff), przywrócenie | Snapshoty + `sergi/go-diff` |
| F5 | Podgląd instrukcji efektywnej | Pokazuje wynikową instrukcję po dziedziczeniu z warstw szerszych (środowisko, globalna) | Rozwinięcie warstw (MODEL-KONFIGURACJI 4) |
| F6 | Test instrukcji w Chat | Kieruje zapytanie próbne z bieżącą wersją bez zapisu do historii właściwej | Kanał testowy Chat Window; brak zapisu |
| F7 | Licznik długości i tokenów | Liczba znaków i przybliżona liczba tokenów instrukcji | `pkoukk/tiktoken-go` |
| F8 | Menedżer snippetów i poleceń | Zapisane, wielokrotne polecenia i fragmenty tekstu wstawiane skrótem | SQLite; wywołanie z kompozytora Chat Window |
| F9 | Nadpisania per agent | Wskazuje, czy agent używa instrukcji projektu, czy własnych (z Agent Builder) | Powiązanie Agent Manager ↔ Instructions Panel |

### Grupa G — Wykonawcy AI (Agent Manager)

| # | Nazwa | Co robi | Zależności |
|---|---|---|---|
| G1 | Lista i przypisanie agentów | Pokazuje agentów z modułu Agents; przypisuje ich do projektu (przeciągnięcie lub przycisk) | Powiązanie z modułem Agents |
| G2 | Zakres uprawnień w projekcie | Ustala, do których okien (Chat Window, Execution Loop Window, Project Library, Context Memory) agent ma dostęp; klucze jawne | Uprawnienia jawne, bez ukrytych bram |
| G3 | Domyślny wykonawca | Oznacza jednego agenta jako sugerowanego w Chat Window i w pętli wykonawczej | Pole „default executor" |
| G4 | Dziennik działań agenta | Chronologiczny zapis czynności agenta w kontekście projektu | Encja `agent_activity` |
| G5 | Skrót do Agent Builder | Tworzy nowego agenta bez opuszczania projektu | Głębokie powiązanie z modułem Agents |
| G6 | Skrót do Permissions Center | Przegląd i zmiana uprawnień platformowych agenta | Powiązanie z Permissions Center |
| G7 | Zasięg przypisania | Wskazuje, czy przypisanie obowiązuje tylko ten projekt, czy także inne (decyzja użytkownika) | Warstwowość zasięgu (MODEL-KONFIGURACJI 6.2) |

### Grupa H — Wyszukiwanie, nawigacja i paleta poleceń

| # | Nazwa | Co robi | Zależności |
|---|---|---|---|
| H1 | Wyszukiwanie globalne projektu | Jedno pole przeszukujące zadania, notatki, pliki, wpisy pamięci i instrukcje | `bleve` / FTS5 nad wszystkimi encjami |
| H2 | Paleta poleceń (`Ctrl/Cmd + K`) | Szybkie przejście do dowolnego okna, elementu lub akcji z klawiatury | Indeks akcji; dopasowanie rozmyte `sahilm/fuzzy` |
| H3 | Wyszukiwanie rozmyte | Toleruje literówki i częściowe dopasowania nazw | `sahilm/fuzzy` |
| H4 | Filtry zapisane i widoki własne | Zapisuje kombinacje filtrów jako nazwany widok (np. „moje pilne zadania") | Definicje widoków w SQLite |
| H5 | Ostatnio używane i ulubione | Lista ostatnio otwartych i przypiętych elementów projektu | Rejestr dostępu; pole „favorite" |
| H6 | Nawigacja skrótami klawiszowymi | Kompletny zestaw skrótów (przełączanie okien, tworzenie zadań i notatek, wyszukiwanie) | Mapa skrótów konfigurowalna |
| H7 | Powiązania krzyżowe elementów | Łączy zadanie z notatką, plikiem lub wpisem pamięci referencją `@` | Tabela referencji uniwersalnych |

### Grupa I — Współpraca, komentarze i aktywność

| # | Nazwa | Co robi | Zależności |
|---|---|---|---|
| I1 | Oś czasu aktywności | Chronologiczny zapis zdarzeń projektu (zmiany instrukcji, wpisy pamięci, pliki, zadania, przebiegi pętli wykonawczej) | Log zdarzeń; WebSocket na żywo |
| I2 | Komentarze do elementów | Wątki komentarzy przy zadaniu, notatce lub pliku; wzmianki `@` | Encja `comment`; powiązania `@` |
| I3 | Wzmianki kontekstowe `@` | Odwołania do pliku, wpisu pamięci, zadania lub agenta wprost w treści polecenia i komentarza | Parser `@`; render plakietek |
| I4 | Notatka właściciela projektu | Wolne pole opisowe: cel, uwagi, przypomnienia widoczne na dashboardzie | Pole tekstowe projektu |
| I5 | Eksport rozmowy do biblioteki | Zapisuje wątek Chat Window jako dokument (MD/PDF) do Project Library | `goldmark` + PDF; zapis pliku |
| I6 | Oznaczanie ustaleń projektu | Oznacza wiadomość czatu jako ustalenie i przypina do Context Memory | Powiązanie Chat Window ↔ Context Memory |
| I7 | Powiadomienia projektu | Zbiorczy strumień powiadomień (terminy, wzmianki, zakończone zlecenia pętli wykonawczej) | Kolejka `notify`; WebSocket + push |

### Grupa J — Import, eksport, migracja i integracje

| # | Nazwa | Co robi | Zależności |
|---|---|---|---|
| J1 | Import z Markdown/folderu | Wciąga zbiór plików `.md` z zachowaniem struktury folderów jako notatki | Rekurencyjny odczyt; `goldmark` |
| J2 | Import eksportu MD/CSV | Wczytuje eksport w formacie Markdown wraz z bazami CSV do stron i tabel | Parser MD + `encoding/csv` |
| J3 | Import tablic z JSON | Odtwarza tablice, listy i karty z eksportu JSON | Parser JSON; mapowanie na zadania |
| J4 | Import kalendarza (iCal) | Wciąga wydarzenia `.ics` jako zadania z terminem | `arran4/golang-ical` |
| J5 | Eksport projektu (paczka) | Pakuje cały projekt (pliki + JSON struktury) do archiwum przenośnego | `archive/zip`; manifest JSON |
| J6 | Eksport do PDF/DOCX | Renderuje notatkę lub podsumowanie do PDF albo DOCX | maroto / `go-fitz` (PDF), `unioffice` (DOCX) |
| J7 | Eksport tabel do CSV/XLSX | Zapisuje tabele zadań i danych do arkusza | `encoding/csv`, `xuri/excelize` (XLSX) |
| J8 | Wiązanie projektu z Automations | Spina projekt z gotową automatyką (np. cykliczny raport do Project Library) | Powiązanie z modułem Automations (SPECYFIKACJA-MODULOW 4.2) |
| J9 | Udostępnienie przez WebDAV | Udostępnia repozytorium projektu przez WebDAV klientom zewnętrznym; sterowane ustawieniem konfiguracyjnym | `golang.org/x/net/webdav` |
| J10 | Obserwacja zmian plików | Wykrywa zmiany w powiązanym folderze lokalnym i aktualizuje bibliotekę | `fsnotify/fsnotify` |

### Grupa K — Automatyzacja projektu i analityka

| # | Nazwa | Co robi | Zależności |
|---|---|---|---|
| K1 | Pasek postępu projektu | Procent ukończenia na podstawie zamkniętych zadań i kamieni milowych | Agregacja w Go; render klienta |
| K2 | Dashboard metryk projektu | Kafle: zadania wg statusu, obłożenie wykonawców, tempo, zbliżające się terminy | Zapytania agregujące SQLite |
| K3 | Wykres wypalania | Wizualizuje postęp iteracji w czasie | Szereg czasowy z logu zadań |
| K4 | Raport obciążenia wykonawców | Rozkład zadań i czasu na agentów i użytkowników | Agregacja `time_entry` + przypisań |
| K5 | Reguły automatyczne projektu | „Gdy zdarzenie X → akcja Y" lokalnie, bez modułu Automations | Silnik reguł Go; log wykonania |
| K6 | Cykliczne migawki stanu | Zapisuje okresową migawkę projektu do audytu i porównań | `robfig/cron`; snapshot do Project Library |
| K7 | Cele i wskaźniki | Śledzi mierzalne cele projektu z procentem realizacji | Encja `goal`; agregacja |
| K8 | Objaśnienia kontekstowe `[?]` | Przy każdym ustawieniu dymek wyjaśniający działanie i wpływ na aplikację | Mechanizm `.dn-tooltip`; treść jako część definicji ustawienia |

**Łączna liczba pozycji katalogu: 92** (A: 8, B: 14, C: 14, D: 16, E: 9, F: 9, G: 7, H: 7, I: 7, J: 10, K: 8).

---

## 9. Punkty sterowania z okna konfiguracji

Wszystko poniżej ustawia Operator z okna konfiguracji (strefa 3 strony głównej) lub z uproszczonego menu kontekstowego okna; każde ustawienie ma objaśnienie `[?]`, klucz jawny, wartość domyślną i warstwę (globalna → środowisko → projekt → sesja). Zero twardych blokad.

| Zakres (wg MODEL-KONFIGURACJI) | Co Operator personalizuje | Domyślna |
|---|---|---|
| Izolacja (5.11, 6) | Profil izolacji projektu; poziom zasięgu (para modułów / projekt / karta sesji); współdzielenie pamięci, instrukcji i historii między projektami | brak aktywnej izolacji technicznej |
| Pamięć (5.10) | Poziom pamięci projektu; włączenie wyszukiwania semantycznego; wybór adaptera osadzeń; próg akceptacji sugestii Wykonawcy; pojemność | poziom sesja / semantyka wył. |
| Prompty systemowe (5.6) | Instrukcje projektu; domyślny szablon roli; źródło instrukcji agenta (projekt albo własne) | brak / dziedziczenie |
| Historia (5.9) | Retencja historii rozmów, przebiegów pętli wykonawczej i migawek projektu | bez limitu |
| Komponenty własne (5.13) | Widoczność projektu (globalny/projektowy); szablon pól projektu; reguły duplikacji | globalny |
| Karty sesji (5.12) | Przypisanie projektu do karty; izolacja per karta albo współdzielenie kart | brak przypisania / izolacja per karta |
| Akcje (5.3) | Zasięg kolejki zleceń projektu; domyślny wykonawca (agent); liczba ponowień zadania w pętli wykonawczej | lokalna |
| Pętla wykonawcza | Szerokość i tryb otwarcia Execution Loop Window; poziom szczegółowości komunikatów sterujących; kryteria kontroli jakości; automatyczne ponowienie zadania | kolumna otwarta / szczegółowość średnia |
| Integracje (5.8) | Powiązania: Agents, Library (synchronizacja wybiórcza), Automations, WebDAV, obserwacja plików | brak powiązania |
| Aplikacja / UX | Zachowanie „Cofnij" wobec modalu potwierdzenia usuwania; skróty klawiszowe; szerokości kolumn obszaru roboczego; domyślny widok hubu planowania (Lista/Kanban/Gantt/Kalendarz); kolumny kanban i stany zadań; folder odbioru artefaktów | „Cofnij" / Lista / stany domyślne |
| Pliki i dokumenty | Włączenie OCR i języki; głębokość wersjonowania; automatyczna deduplikacja; limity i progi ostrzeżeń; mechanizm indeksu pełnotekstowego (FTS5 albo Bleve) | OCR wł. (pol+eng) / wersje bez limitu |
| Powiadomienia | Kanały (WebSocket/push); wyprzedzenie przypomnień terminów; wyciszenia | wł. / 1 dzień |

---

*Koniec dokumentu. Moduł WorkSpace — dokumentacja projektowa, wersja 2.0, 2026-08-06.*

---
*Danaco Console — AI Workspace OS · v2.0*
*© 2026 Danaco Holding Group Sp. z o.o. Wszelkie prawa zastrzeżone — [LICENSE](LICENSE). Kontakt: support@danaco-group.pl*
