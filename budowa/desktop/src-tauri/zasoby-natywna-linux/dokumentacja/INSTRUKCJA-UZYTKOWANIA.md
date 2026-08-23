# Danaco Console — Instrukcja użytkowania

**Praca na platformie: od uruchomienia okna do prowadzenia autonomicznej pętli pracy ciągłej.**

| | |
|---|---|
| **Produkt** | Danaco Console |
| **Rodzaj** | Platforma AI Workspace OS |
| **Wersja** | v2.0 |
| **Adresat dokumentu** | Operator — właściciel konta pracujący na platformie |
| **Producent** | Danaco Holding Group Sp. z o.o., ul. Gen. J. Hallera 81, 43-400 Cieszyn |
| **Wsparcie** | support@danaco-group.pl |

Instalację pakietu klienckiego, zakładanie konta i pełną konfigurację opisuje
`INSTALACJA-I-KONFIGURACJA.md`. Niniejszy dokument opisuje pracę.

---

## Spis treści

1. [Jak czytać tę instrukcję](#1-jak-czytać-tę-instrukcję)
2. [Uruchomienie i wejście do platformy](#2-uruchomienie-i-wejście-do-platformy)
3. [Strona główna — Centrum dowodzenia](#3-strona-główna--centrum-dowodzenia)
4. [Nawigacja: środowisko → moduł → okno](#4-nawigacja-środowisko--moduł--okno)
5. [Karty sesji — praca równoległa](#5-karty-sesji--praca-równoległa)
6. [Chat Window — sterowanie platformą](#6-chat-window--sterowanie-platformą)
7. [Execution Loop Window — pętla wykonawcza](#7-execution-loop-window--pętla-wykonawcza)
8. [Cztery środowiska](#8-cztery-środowiska)
9. [Moduły — przewodnik po piętnastu obszarach pracy](#9-moduły--przewodnik-po-piętnastu-obszarach-pracy)
10. [MultitaskingAI — prowadzenie zespołu modeli](#10-multitaskingai--prowadzenie-zespołu-modeli)
11. [Komponenty własne](#11-komponenty-własne)
12. [Magistrala kontekstu — przekazywanie artefaktów](#12-magistrala-kontekstu--przekazywanie-artefaktów)
13. [Funkcje globalne: Mobile i Always On Display](#13-funkcje-globalne-mobile-i-always-on-display)
14. [Szybka zmiana konfiguracji w trakcie pracy](#14-szybka-zmiana-konfiguracji-w-trakcie-pracy)
15. [Jak dotrzeć do funkcji, której nie widać](#15-jak-dotrzeć-do-funkcji-której-nie-widać)
16. [Historia, pamięć i artefakty](#16-historia-pamięć-i-artefakty)
17. [Zakończenie i wznowienie pracy](#17-zakończenie-i-wznowienie-pracy)
18. [Typowe scenariusze pracy](#18-typowe-scenariusze-pracy)
19. [Słowniczek](#19-słowniczek)

---

## 1. Jak czytać tę instrukcję

Praca na platformie przebiega zawsze tą samą drogą: **strona główna → środowisko
→ moduł → okno operacyjne**, a sterowanie idzie dwoma kanałami komunikacji —
Chat Window i Execution Loop Window. Rozdziały 2–7 opisują tę drogę i oba kanały;
kto je przeczyta, umie pracować w każdym module. Rozdziały 8–11 opisują, co
konkretnie da się zrobić w poszczególnych środowiskach i modułach. Rozdziały 12–17
opisują mechanizmy przechodzące przez całą platformę.

Zasada, którą warto znać od początku: **jeżeli funkcja nie jest potrzebna do
bieżącego zadania, nie jest widoczna** — ale każda ukryta funkcja jest osiągalna
jednym kliknięciem, jednym skrótem klawiszowym albo jednym poleceniem w języku
naturalnym. Rozdział 15 pokazuje trzy drogi dotarcia do niej.

---

## 2. Uruchomienie i wejście do platformy

```
Uruchomienie Danaco Console na urządzeniu
        ▼
OKNO STARTOWE — godło marki, „łączenie z serwerem"
        │
        ├── urządzenie ma ważny token ──► STRONA GŁÓWNA
        │                                 albo ostatnio aktywna karta sesji
        │
        └── brak tokenu ──► OKNO LOGOWANIA ──► STRONA GŁÓWNA
```

Okno startowe jest stanem przejściowym — trwa od ułamka sekundy do kilku sekund
i nie wymaga żadnej decyzji. Przy braku połączenia pojawia się w nim komunikat
o nieosiągalnym serwerze i kontrolka „Spróbuj ponownie".

W oknie logowania Operator podaje login i hasło albo korzysta z aktywnej metody
dodatkowej — PIN, Windows Hello, adres e-mail uwierzytelniający. Po pozytywnej
weryfikacji urządzenie otrzymuje token dostępu i przy kolejnych uruchomieniach
wchodzi bez logowania.

---

## 3. Strona główna — Centrum dowodzenia

Główne okno aplikacji nie otwiera od razu żadnej przestrzeni roboczej. Jest
przedpokojem przed strefą pracy — miejscem, z którego Operator wybiera środowisko,
buduje własne komponenty albo przechodzi do ustawień.

```
STRONA GŁÓWNA — Centrum dowodzenia
├─ Strefa 1 · wybór środowiska    → cztery duże karty       (waga główna)
│     TalkIn · WorkSpace · CodeStudio · MultitaskingAI
├─ Strefa 2 · komponenty własne   → cztery kafle            (waga pośrednia)
│     Automations · Agents · Workspace · Assistant
└─ Strefa 3 · ustawienia          → zwarta listwa           (waga najniższa)
      Okno konfiguracji · Mobile · Always On Display
```

| Strefa | Co robi kliknięcie |
|---|---|
| **Strefa 1** | Otwiera przestrzeń roboczą wybranego środowiska. Jest to **jedyna droga wejścia** do środowiska |
| **Strefa 2** | Otwiera okno konfiguracji komponentu własnego — automatyki, agenta, projektu albo profilu asystenta (rozdz. 11) |
| **Strefa 3** | Otwiera okno konfiguracji platformy albo aktywuje funkcję globalną (rozdz. 13) |

Chat Window obecne jest także na stronie głównej — polecenie można wydać, jeszcze
nie wchodząc do żadnego środowiska. Górny pasek marki niesie wyszukiwanie globalne
i wskaźnik konta; kliknięcie godła wraca na stronę główną z dowolnego miejsca
platformy.

---

## 4. Nawigacja: środowisko → moduł → okno

```
Strona główna
    │  klik karty środowiska (strefa 1)
    ▼
Okno środowiska  ──(boczna nawigacja modułów / panel orkiestracji)──►  Moduł
                                                                          │
                                                                          ▼
                                                              Okno operacyjne — praca właściwa
```

| Krok | Co się dzieje |
|---|---|
| 1. Wybór środowiska | Strona główna zostaje zastąpiona przestrzenią roboczą; otwiera się nowa karta sesji albo przywracana jest ostatnio aktywna |
| 2. Wybór modułu | Boczna nawigacja pokazuje moduły dostępne w tym środowisku; w MultitaskingAI zamiast listy modułów jest panel orkiestracji |
| 3. Okno operacyjne | Moduł otwiera właściwe mu okna — tu zaczyna się praca |
| 4. Karty sesji | Praca równoległa przez wiele kart otwartych obok siebie (rozdz. 5) |

**Układ okna środowiska** jest stały i wyłącznie kolumnowy:

```
 ═══════════════════════════════════════════════════════════════════════════
  Boczna    │ Chat Window          │ Obszar roboczy modułu    │ Panel
  nawigacja │ Użytkownik ↔         │                          │ pomocniczy
  modułów   │ Wykonawca            │  (okna edycyjne,         │ (rozszerzenie
            │ ─────────────────    │   podglądu, monitory)    │  boczne)
            │ Execution Loop       │                          │
            │ Koordynator ↔        │                          │
            │ Wykonawca            │                          │
 ═══════════════════════════════════════════════════════════════════════════
```

Boczna nawigacja pozostaje widoczna niezależnie od otwartego modułu. Execution
Loop Window otwiera się w kolumnie sąsiadującej z Chat Window w chwili przyjęcia
pierwszego zlecenia. Okna pomocnicze dochodzą jako kolejne kolumny po prawej
stronie. Regulacji podlega wyłącznie szerokość kolumn — żaden element nie jest
rozmieszczany w podziale poziomym.

**Zmiana modułu w tej samej karcie przeładowuje przestrzeń:** znikają okna
poprzedniego modułu, pojawia się zestaw dopasowany do nowego kontekstu. Zachowane
zostają: sama karta, boczna nawigacja, pozostałe karty i funkcje globalne.

---

## 5. Karty sesji — praca równoległa

Karta sesji jest jednostką pracy równoległej — mechanika odpowiada kartom
przeglądarki. W obrębie jednej karty widoczny jest interfejs jednego modułu; praca
nad kilkoma rzeczami naraz odbywa się przez otwarcie kilku kart.

| Czynność | Skutek |
|---|---|
| Otwarcie nowej karty | Nowa, pusta przestrzeń robocza oczekująca na wybór modułu |
| Przełączenie karty | Przejście do innej sesji z jej własnym układem, historią i kontekstem |
| Zamknięcie karty | Znika widok tej karty; pozostałe pracują dalej |

Każda karta ma własny układ, historię i kontekst. **Zakres, w jakim karty
współdzielą kontekst, historię i pamięć, ustala Operator** — stanem wyjściowym
jest uporządkowany podział pracy, a współdzielenie jest świadomą decyzją
podejmowaną w oknie konfiguracji (rozdz. 14) albo jednorazowo na poziomie karty.

W środowisku MultitaskingAI karta sesji odpowiada roli — Executor 1, Executor 2,
Coordinator, Executor 3 / Validator.

---

## 6. Chat Window — sterowanie platformą

Chat Window jest głównym oknem komunikacji z Wykonawcą — modelem AI, agentem albo
systemem wykonawczym — i **podstawowym mechanizmem sterowania wszystkimi procesami
platformy**. Zajmuje lewą kolumnę obszaru roboczego, o pełnej wysokości
i szerokości regulowanej przez Operatora, i jest obecne w każdym module.

| Funkcja | Jak z niej korzystać |
|---|---|
| **Polecenie** | Sformułować zadanie w języku naturalnym i wysłać do Wykonawcy |
| **Strumień odpowiedzi** | Odpowiedź, wyniki cząstkowe i artefakty pojawiają się na żywo, w miarę generowania |
| **Zatwierdzenie** | Zaakceptować proponowane działanie, zanim zostanie wykonane |
| **Przerwanie** | Zatrzymać działanie w toku |
| **Wyjaśnienie** | Zażądać wyjaśnienia wyniku, przyjętych założeń i kontekstu, w jakim zadanie zrealizowano |

Chat Window rekonfiguruje się przy każdej zmianie środowiska albo modułu: wygląd,
możliwości, historia, dostępne narzędzia i kontekst dopasowują się do bieżącej
pracy. Czat w module Developer zna kod i repozytorium, czat w module Research —
zebrane źródła.

W pasku kontekstu nad rozmową widoczne są lekkie znaczniki bieżącego układu pracy —
środowisko, projekt, model, wykonawca, tryb. Kliknięcie znacznika otwiera właściwy
selektor; zmiana obowiązuje bieżącą sesję.

**Chat Window przyjmuje też polecenia sterujące samą platformą** — w języku
naturalnym. To najkrótsza droga do funkcji z wyższych warstw widoczności
(rozdz. 15).

---

## 7. Execution Loop Window — pętla wykonawcza

Drugi kanał komunikacji pokazuje pracę między Koordynatorem a Wykonawcą.
Koordynator dekomponuje zlecenie, przydziela zadania, nadzoruje przebieg
i kontroluje realizację; Operator ten przebieg widzi i steruje nim.

| Zawartość okna | Co daje Operatorowi |
|---|---|
| Zlecenie i dekompozycja | Widok bieżącego zlecenia rozłożonego na zadania składowe |
| Kolejka i stan zadań | Co czeka, co się wykonuje, co skończone, co wymaga uwagi |
| Komunikaty sterujące | Wymiana między Koordynatorem a Wykonawcą |
| Kontrola jakości | Wyniki kontroli i decyzje o ponowieniu zadania |
| Wskaźniki przebiegu | Postęp pętli jako całości |
| Sterowanie przebiegiem | Wstrzymanie, wznowienie, przerwanie i korekta zlecenia |

Okno otwiera się w chwili przyjęcia pierwszego zlecenia i pozostaje otwarte przez
cały czas trwania pętli. Wyzwalacz otwarcia i wskaźnik stanu pętli są widoczne
zawsze; szczegółowa zawartość pojawia się na żądanie.

**Granice autonomii Wykonawcy ustala Operator** — w konfiguracji okna pętli:
zakres autonomii, punkty obowiązkowego zatwierdzenia, progi kontroli jakości
i politykę ponowień. Działania nieodwracalne wymagają zatwierdzenia niezależnie od
ustawionego zakresu autonomii. Przebieg pętli i podjęte decyzje zapisywane są
w dzienniku audytu.

---

## 8. Cztery środowiska

| Środowisko | Kiedy je wybrać |
|---|---|
| **TalkIn** | Przedmiotem pracy jest treść: rozmowa z AI, dokumenty, analizy, tłumaczenia, badania, raporty, współpraca wielu modeli |
| **WorkSpace** | Przedmiotem pracy jest zorganizowane działanie: projekty, procesy, automatyzacja, aplikacje biznesowe, projektowanie |
| **CodeStudio** | Przedmiotem pracy jest kod i środowisko wykonawcze: tworzenie oprogramowania, debugowanie, terminale, procesy developerskie |
| **MultitaskingAI** | Przedmiotem pracy jest zespół modeli i agentów realizujący wspólny proces — także w trybie ciągłym (rozdz. 10) |

**Który moduł jest gdzie dostępny:**

| Moduł | TalkIn | WorkSpace | CodeStudio |
|---|:---:|:---:|:---:|
| Studio | ● | ● | |
| Workspace | ● | ● | ● |
| Browser | ● | ● | |
| Research | ● | ● | |
| Library | ● | ● | |
| Translate | ● | | |
| Roundtable | ● | ● | ● |
| Design | | ● | ● |
| Assistant | ● | | |
| Terminal | | | ● |
| Developer | | | ● |
| Diagnostics | | | ● |
| Apps | | ● | ● |
| Agents | ● | ● | ● |

Moduł **Automations** nie ma okna w bocznej nawigacji — konfiguruje się go ze
strony głównej, a gotowe automatyki wpina się do sesji jako komponenty własne.
Środowisko MultitaskingAI nie ma listy modułów; jego boczną nawigacją jest panel
orkiestracji.

---

## 9. Moduły — przewodnik po piętnastu obszarach pracy

Każdy moduł udostępnia oba kanały komunikacji operacyjnej oraz własne okna
robocze. Poniżej: po co moduł istnieje, jakie ma okna i jak wygląda w nim typowy
przebieg pracy.

### 9.1. Studio — praca z tekstem i dokumentami

**Po co:** tworzenie, redagowanie i przekształcanie materiałów pisanych, od notatek
po obszerne dokumenty — w jednym miejscu, bez rozdzielania edytora i okna rozmowy.

**Okna:** Studio Editor · Tools Panel · Diff/Grep Panel · Session Repository ·
Preview Window

**Zakres:** edycja dokumentów, obsługa PDF, DOCX, TXT i Markdown, tłumaczenia,
korekta, analiza treści, przepisywanie, zmiana stylu, streszczenia, rozwijanie
treści, Diff, Grep, operacje kontekstowe AI.

**Typowy przebieg:**
1. Wczytać dokument do Studio Editor.
2. Zlecić z Chat Window korektę albo zmianę stylu wybranego fragmentu.
3. Porównać wersję przed i po zmianie w Diff/Grep Panel.
4. Zaakceptować wynik; wcześniejsze wersje pozostają dostępne w Session Repository.

### 9.2. Workspace — izolowane środowisko projektowe

**Po co:** prowadzenie równolegle wielu odrębnych projektów z rozdzielonym
kontekstem, instrukcjami i pamięcią, tak aby ustalenia jednego nie przenikały do
drugiego.

**Okna:** Project Dashboard · Instructions Panel · Context Memory ·
Project Library · Agent Manager

**Typowy przebieg:**
1. Założyć projekt.
2. Zdefiniować odrębny zestaw instrukcji systemowych w Instructions Panel.
3. Zbudować dedykowaną pamięć kontekstową w Context Memory.
4. Przypisać do projektu konkretnych agentów przez Agent Manager.

Project Library pełni rolę biblioteki ograniczonej do jednego projektu.

### 9.3. Automations — procesy automatyczne

**Po co:** zamienić powtarzalne czynności — cykliczne raporty, regularne
przetwarzanie danych, zaplanowane wywołania modeli — w procesy działające bez
stałego nadzoru.

**Okna:** Workflow Builder · Scheduler · Queue Manager · Orchestrator ·
Execution Monitor

**Typowy przebieg:**
1. Zaprojektować proces w Workflow Builder.
2. Ustalić jego cykliczność w Scheduler.
3. Obserwować kolejne uruchomienia w Execution Monitor — status każdego przebiegu
   i ewentualne błędy wymagające interwencji.

Automatyka zapisana jako komponent własny wpina się później do sesji dowolnego
modułu. Silnik kolejek modułu można spiąć z silnikiem kolejek MultitaskingAI
(rozdz. 10).

### 9.4. Browser — wspólne przeglądanie internetu

**Po co:** sytuacje, w których Operator i AI muszą pracować nad tą samą treścią
internetową w tym samym czasie, zamiast przekazywać sobie odsyłacze i fragmenty.

**Okna:** Browser Window · Sources Panel · Notes Panel

**Typowy przebieg:**
1. Otworzyć stronę w Browser Window.
2. AI — mając ten sam podgląd — odpowiada na pytania o treść i wyszukuje
   informacje powiązane.
3. Odnotować istotne fragmenty w Notes Panel.
4. Zbudować listę wykorzystanych źródeł w Sources Panel.

### 9.5. Research — badania, analizy, opracowania

**Po co:** praca wymagająca systematycznego zbierania, porządkowania i syntezy
informacji z wielu źródeł — od analiz rynkowych po opracowania konkurencyjne.

**Okna:** Research Workspace · Sources Manager · Findings Panel ·
Report Builder · Export Panel

**Typowy przebieg:**
1. Zgromadzić źródła w Sources Manager.
2. Odnotowywać ustalenia cząstkowe w Findings Panel w miarę postępu analizy.
3. Skompletować raport końcowy w Report Builder.
4. Wyeksportować go przez Export Panel do formatu wymaganego przez odbiorcę.

### 9.6. Library — repozytorium wiedzy i plików

**Po co:** jedno uporządkowane miejsce na materiały wykorzystywane i wytwarzane
w pracy — trwała pamięć zewnętrzna wspólna dla środowisk i modułów.

**Okna:** Library Explorer · Tags & Collections · File Preview ·
Versioning Panel

**Typowy przebieg:**
1. Porządkować materiały przez Tags & Collections.
2. Przeglądać zawartość bez opuszczania modułu w File Preview.
3. Śledzić kolejne wersje dokumentu w Versioning Panel — w tym wersje wytworzone
   przez AI w innych modułach.

### 9.7. Translate — tłumaczenia wielojęzyczne

**Po co:** praca równoległa z treścią w wielu językach, której nie obsłuży
sekwencyjne tłumaczenie w oknie czatu.

**Okna:** Source Panel · Translation Panels · Glossary Manager

**Typowy przebieg:**
1. Umieścić tekst źródłowy w Source Panel.
2. Tłumaczyć równocześnie na wiele języków widocznych w osobnych panelach.
3. Utrzymać spójność terminologii przez Glossary Manager, definiując preferowane
   odpowiedniki kluczowych pojęć.

### 9.8. Roundtable — współpraca wielu modeli

**Po co:** zadania, w których odpowiedź jednego modelu nie wystarcza, a wartość
wynika z konfrontacji perspektyw i wypracowania wspólnego stanowiska.

**Okna:** Model Panels · Debate Panel · Moderator Panel · Consensus Panel

**Typowy przebieg:**
1. Kilka modeli widocznych równolegle w Model Panels odpowiada na to samo
   zagadnienie.
2. Debate Panel rejestruje wymianę argumentów.
3. Moderator Panel pozwala ukierunkować dyskusję.
4. Consensus Panel gromadzi uzgodnione stanowisko.

### 9.9. Design — zasoby wizualne

**Po co:** generowanie i edycja materiałów graficznych — od ilustracji
i elementów brandingowych po makiety interfejsu — ze wsparciem AI na każdym etapie.

**Okna:** Design Board · Assets Panel · Prompt Builder · Preview Window

**Typowy przebieg:**
1. Sformułować precyzyjne polecenia generujące grafikę w Prompt Builder.
2. Gromadzić wygenerowane zasoby w Assets Panel.
3. Zestawiać je i prowadzić dalszą pracę koncepcyjną w Design Board.

### 9.10. Assistant — komunikacja głosowa

**Po co:** sytuacje, w których mowa jest wygodniejsza od pisania — praca w ruchu,
sterowanie zadaniami bez klawiatury.

**Okna:** Voice Console · Actions Monitor · Activity Feed

**Typowy przebieg:**
1. Wydawać polecenia głosowe przez Voice Console.
2. Obserwować status ich realizacji w Actions Monitor.
3. Odtworzyć przebieg wieloetapowego zlecenia z chronologicznego zapisu
   w Activity Feed.

Profil asystenta konfiguruje się na stronie głównej, a wykorzystuje operacyjnie
w środowisku TalkIn.

### 9.11. Terminal — konsole i środowiska wykonawcze

**Po co:** bezpośredni dostęp do powłok systemowych i narzędzi wiersza polecenia
z poziomu platformy, bez przełączania się do zewnętrznej aplikacji.

**Okna:** Terminal Tabs · Output Console · Process Monitor

**Zakres:** PowerShell, CMD, Bash, Node.js, Python, narzędzia wiersza polecenia.

**Typowy przebieg:**
1. Prowadzić równolegle wiele sesji w różnych powłokach w Terminal Tabs.
2. Gromadzić ich wynik w Output Console.
3. Obserwować i kontrolować uruchomione procesy w Process Monitor — w tym te
   zainicjowane poleceniem AI.

### 9.12. Developer — tworzenie i rozwój kodu

**Po co:** praca programistyczna z pełnym zestawem narzędzi inżynierskich
i wsparciem AI działającym w kontekście konkretnego repozytorium.

**Okna:** Code Editor · Project Tree · Git Panel · Build Output

**Typowy przebieg:**
1. Nawigować po strukturze projektu w Project Tree.
2. Edytować kod w Code Editor przy wsparciu AI z Chat Window.
3. Zarządzać zmianami przez Git Panel.
4. Obserwować wynik kompilacji albo budowania w Build Output.

### 9.13. Diagnostics — analiza i usuwanie problemów

**Po co:** sytuacje, w których punktem wyjścia jest nie nowa funkcja, lecz
zrozumienie przyczyny błędu albo spadku wydajności.

**Okna:** Diagnostics Center · Logs Viewer · Errors Panel ·
Recommendations Panel

**Typowy przebieg:**
1. Zebrać materiał źródłowy w Logs Viewer i Errors Panel.
2. Zagregować go w spójny obraz stanu w Diagnostics Center.
3. Przejrzeć kroki naprawcze wypracowane przez AI w Recommendations Panel.

### 9.14. Apps — budowa produktów cyfrowych

**Po co:** realizacja pełnych projektów aplikacyjnych — od architektury, przez
frontend i backend, po wdrożenie — w jednym spójnym procesie.

**Okna:** Product Builder · Architecture Designer · Frontend Workspace ·
Backend Workspace · Deployment Panel

**Typowy przebieg:**
1. Zaprojektować architekturę rozwiązania w Architecture Designer.
2. Pracować równolegle w Frontend Workspace i Backend Workspace.
3. Poprowadzić udostępnienie gotowego produktu przez Deployment Panel.

### 9.15. Agents — budowa własnych agentów

**Po co:** trwałe, wyspecjalizowane jednostki AI z określoną tożsamością,
umiejętnościami i uprawnieniami, wykorzystywane później jako wykonawcy zadań
w całej platformie.

**Okna:** Agent Builder · Model Configuration · Skills Manager ·
Connectors Manager · Permissions Center

**Typowy przebieg:**
1. Utworzyć agenta w Agent Builder.
2. Wybrać model bazowy w Model Configuration.
3. Dobrać umiejętności w Skills Manager.
4. Podłączyć integracje zewnętrzne przez Connectors Manager.
5. Ustalić zakres uprawnień w Permissions Center — **przed** udostępnieniem agenta
   do pracy operacyjnej.

Agent działa ponad wszystkimi środowiskami; po przypisaniu do roli pracuje także
w MultitaskingAI.

---

## 10. MultitaskingAI — prowadzenie zespołu modeli

W tym środowisku przedmiotem konfiguracji przestaje być pojedyncza interakcja
z modelem — staje się nim cały zespół realizujący wspólny proces.

### 10.1. Panel orkiestracji

Zamiast listy modułów boczna nawigacja zawiera sześć sekcji:

| Sekcja | Co w niej robisz |
|---|---|
| **Zespoły** | Zapisujesz, wczytujesz i duplikujesz konfiguracje zespołu — presety ról, powiązań i kolejek |
| **Role** | Prowadzisz cztery okna robocze ról i przypisujesz do nich własnych agentów; włączasz Subagent Network |
| **Kolejki** | Definiujesz kolejki globalne, lokalne, modeli, agentów i projektów oraz ich akcje |
| **Orkiestracja** | Ustalasz zależności między modelami, agentami, zadaniami, kolejkami, automatykami i projektami |
| **Harmonogram i automatyki** | Ustawiasz harmonogram pracy ciągłej i wpinasz automatyki z modułu Automations |
| **Monitor procesu** | Obserwujesz przebieg pętli, statusy i hierarchię decyzji; tu nadzoruje Always On Display |

Kolejność i widoczność sekcji podlega konfiguracji — powyższy zestaw jest zestawem
domyślnym.

### 10.2. Role i podział pracy

| Rola | Co robi |
|---|---|
| **Executor 1** | Wykonuje pracę właściwą: dokumenty, analizy, kod, projekty, aplikacje. Nie zarządza procesem |
| **Subagent Network** | Do 15 wyspecjalizowanych podagentów uruchamianych przez wykonawcę; wyniki agreguje wykonawca |
| **Coordinator** | Planuje etapy, zależności i harmonogramy, dzieli pracę, buduje prompty, steruje procesem i kolejką |
| **Executor 2** | Drugi niezależny wykonawca — praca równoległa: praca niezależna, przekazywanie wyników, praca naprzemienna albo iteracyjna |
| **Executor 3 / Validator** | Rola kontrolna albo doradcza: Validator, Reviewer, Security Auditor, Architect, Product Owner, QA Lead, Arbitrator |

Rozdzielenie Koordynatora od Wykonawcy jest celowe — żaden model nie musi
jednocześnie planować i realizować.

### 10.3. Kolejki i orkiestracja

Kolejki definiuje się na poziomie globalnym, lokalnym, modelu, agenta albo
projektu, z akcjami: `enqueue`, `dequeue`, `delay`, `retry`, `pause`, `resume`,
`split`, `merge`, `route`, `branch`, `condition`.

Warstwa orkiestracji pilnuje, by ustalone zależności były respektowane w toku
wykonania: zadanie drugiego wykonawcy nie rozpocznie się przed zakończeniem etapu
pierwszego, jeśli tak ustalono, a wynik trafi do właściwej kolejki albo roli.

### 10.4. Praca ciągła 24/7/365

Po spięciu środowiska z modułem Automations powstaje pełna pętla autonomiczna:

| Element | Rola w pętli |
|---|---|
| Coordinator | Planuje proces |
| Automations | Prowadzi harmonogramy i kolejki |
| Executor 1, Executor 2 | Realizują zadania |
| Executor 3 / Validator | Kontroluje jakość |
| Always On Display | Nadzoruje całość z perspektywy Operatora |

Spięcie nie jest domyślne — wymaga decyzji Operatora w oknie konfiguracji
i w ustawieniach modułu Automations. Nadzór nad taką pętlą działa również
z urządzenia przenośnego (rozdz. 13).

---

## 11. Komponenty własne

Komponent własny to nazwany wytwór Operatora, budowany na stronie głównej
i wpinany później tam, gdzie ma zastosowanie.

| Komponent | Gdzie powstaje | Gdzie się go używa |
|---|---|---|
| **Automatyka** | Strefa 2 → Automations | Wpinana w sesji dowolnego modułu; zaplecze harmonogramów MultitaskingAI |
| **Agent** | Strefa 2 → Agents | Wykonawca zadań w TalkIn, WorkSpace, CodeStudio; po przypisaniu do roli — w MultitaskingAI |
| **Projekt** | Strefa 2 → Workspace | Izolowana przestrzeń projektowa z własnym kontekstem, pamięcią i biblioteką |
| **Profil asystenta** | Strefa 2 → Assistant | Praca głosowa w środowisku TalkIn |

Po utworzeniu komponent trafia do zasobów Operatora i pozostaje wybieralny
w trakcie sesji. Moduły Agents i Workspace mają dodatkowo własne okna wewnątrz
środowisk — komponent można więc budować przed pracą i doprecyzowywać w jej toku.

---

## 12. Magistrala kontekstu — przekazywanie artefaktów

Magistrala przenosi artefakty między modułami. Występuje pod tą samą nazwą
w każdym środowisku i module.

| Krok | Jak go wykonać |
|---|---|
| **Odłożenie** | Pozycja „Do magistrali" w menu kontekstowym artefaktu albo przeciągnięcie na krawędź obszaru roboczego |
| **Przegląd** | Kliknięcie licznika magistrali w pasku kontekstu — kolumna pokazuje odłożone pozycje z modułem pochodzenia, typem i znacznikiem czasu |
| **Pobranie** | „Otwórz w…", „Dołącz do rozmowy" albo przeciągnięcie pozycji do modułu docelowego |
| **Skierowanie wprost** | „Przekaż do modułu" — jednym krokiem, z pominięciem przeglądu |
| **Zwolnienie** | „Usuń z magistrali" — pozycja znika, artefakt pozostaje w obu modułach |

**Magistrala przechowuje odwołania, nie kopie.** Przekazanie nie tworzy nowego
artefaktu ani nowej wersji — moduł docelowy pracuje na tym samym artefakcie,
a każda utrwalona zmiana powstaje jako jego kolejna wersja. Licznik pozycji
magistrali widoczny jest w pasku kontekstu bez żadnej interakcji.

---

## 13. Funkcje globalne: Mobile i Always On Display

Funkcje globalne nie tworzą przestrzeni roboczej i nie podlegają przełączaniu —
ich aktywacja nie zmienia niczego w bieżącej pracy.

### 13.1. Mobile — nadzór z telefonu

| Co można z urządzenia przenośnego | Zakres |
|---|---|
| Chat Window | Pełny: polecenie, strumień odpowiedzi, zatwierdzanie, przerywanie, wyjaśnianie |
| Execution Loop Window | Pełny nadzór: dekompozycja, kolejka, komunikaty sterujące, kontrola jakości, wstrzymanie, wznowienie, przerwanie, korekta |
| Przegląd procesów | Projekty, automatyzacje, agenci, aplikacje, sesje, przebiegi MultitaskingAI, kolejki |
| Zarządzanie zadaniami | Uruchamianie, zatrzymywanie, zatwierdzanie, wstrzymywanie, wznawianie, modyfikowanie |
| Decyzje procesu | Zatwierdzanie punktów decyzyjnych, odrzucenie wyniku, decyzja o ponowieniu |
| Podgląd artefaktu | Treść, wersja i metadane — na tyle, by podjąć decyzję o zatwierdzeniu |
| Praca bez połączenia | Przewidziana; po odzyskaniu łączności stan uzgadnia się z serwerem |

Przy stanowisku stacjonarnym pozostaje praca konstrukcyjna: zakładanie kart sesji
i przełączanie środowisk, okna operacyjne modułów, okno konfiguracji i okno
punktów izolacji, pełny panel orkiestracji, rejestracja konta i inicjowanie
parowania urządzeń, definiowanie agentów, przepływów i rozszerzeń.

### 13.2. Always On Display — agent towarzyszący

Pływający avatar obecny we wszystkich częściach platformy jednocześnie. Ma dostęp
do środowisk, projektów, sesji, historii rozmów i agentów Operatora, dzięki czemu:

- proponuje działania, sugeruje konfiguracje, wskazuje problemy i rekomenduje
  kolejne kroki;
- rozumie kontekst — wie, w jakim module, środowisku i projekcie trwa praca;
- rozmawia głosem i tekstem, przyjmuje szybkie polecenia;
- w MultitaskingAI pełni rolę obserwatora albo operatora nadzorującego proces.

Reguły wyzwalania sugestii, ich progi i częstotliwość, reguły wyciszania oraz tor
głosowy są ustawieniami — funkcję da się prowadzić od pełnej proaktywności do
całkowitej ciszy.

---

## 14. Szybka zmiana konfiguracji w trakcie pracy

Nie każda zmiana wymaga wyjścia do okna konfiguracji. **Menu kontekstowe okna
operacyjnego** daje dostęp do ustawień warstwy sesji bez przerywania pracy.

| Droga | Zasięg zmiany | Kiedy używać |
|---|---|---|
| Menu kontekstowe okna operacyjnego | Wyłącznie bieżąca karta sesji | Jednorazowa korekta zachowania modelu, izolacji albo pamięci w trakcie pracy |
| Pełne okno konfiguracji (listwa ustawień strony głównej) | Warstwa globalna, środowiska, projektu i sesji | Trwała polityka pracy |

Zmiana na warstwie sesji **nakłada się** na wartość szerszą, nie zmieniając jej —
ustawienie wprowadzone dla jednej karty nie dotyczy pozostałych. Przy każdej
pozycji ustawienia widoczny jest poziom warstwy, z której pochodzi wartość
bieżąca, oraz objaśnienie `[?]` mówiące, co to ustawienie robi.

---

## 15. Jak dotrzeć do funkcji, której nie widać

Interfejs ujawnia możliwości stopniowo. To, że funkcji nie widać, nie znaczy, że
jest niedostępna — znaczy, że nie jest potrzebna do bieżącej czynności.

| Chcesz | Zrób |
|---|---|
| Zmienić model, wykonawcę, tryb pracy | Kliknij znacznik w pasku kontekstu albo zwinięty selektor `▼` |
| Znaleźć zestaw akcji dla elementu | Otwórz menu `⋮` przy elemencie albo menu `☰` środowiska |
| Wywołać funkcję zaawansowaną | Wpisz polecenie w języku naturalnym w Chat Window |
| Znaleźć funkcję, której nazwy nie pamiętasz | Użyj wyszukiwarki funkcji |
| Pracować szybciej | Użyj skrótu klawiszowego; paleta poleceń otwiera się skrótem `Ctrl/Cmd + K` |
| Odsłonić wszystkie funkcje eksperckie | Włącz tryb administracyjny |

**Zagnieżdżanie funkcji głęboko w menu jest w tym produkcie zabronione** — każda
ukryta funkcja jest osiągalna jednym kliknięciem, jednym skrótem albo jednym
poleceniem. Zakres widocznych funkcji zależy od profilu warstw przypisanego roli
użytkownika; profile te są ustawieniem konfiguracji.

---

## 16. Historia, pamięć i artefakty

| Mechanizm | Co robi | Czym steruje Operator |
|---|---|---|
| **Historia** | Zapis rozmów i sesji | Zakresem przechowywania, dostępnością i retencją — globalnie albo dla jednej sesji |
| **Pamięć kontekstowa** | To, co AI pamięta poza treścią bieżącej rozmowy | Zasobami pamięci na czterech poziomach oraz tym, co jest współdzielone, a co odrębne |
| **Artefakty** | Wytwory pracy: dokumenty, raporty, grafiki, pliki kodu | Miejscem odbioru (Library, Project Library), wersjonowaniem, przekazywaniem magistralą |
| **Punkty izolacji** | Rozstrzygnięcie, co jest wspólne, a co osobne między przestrzeniami pracy | Macierzą izolacji na siedmiu poziomach zasięgu i profilami izolacji |

**Ciągła, wspólna historia jest ustawieniem, nie stanem narzuconym.** Stanem
wyjściowym jest uporządkowany podział pracy — współdzielenie historii, pamięci
i kontekstu między kartami, modułami, projektami i środowiskami włącza Operator
świadomie, tam gdzie jest pożądane. Library pełni rolę trwałego repozytorium
wspólnego dla środowisk i modułów; wersje dokumentu — także te wytworzone przez AI
w innym module — widać w Versioning Panel.

---

## 17. Zakończenie i wznowienie pracy

**Sesja jest procesem po stronie serwera, nie stanem okna na urządzeniu.**

```
Rozłączenie klienta ─▶ proces sesji pracuje dalej ─▶ ponowne połączenie
                                                     odtwarza pełny stan:
                                                     układ, historię, kontekst
```

| Zdarzenie | Co pozostaje |
|---|---|
| Powrót na stronę główną i wejście do innego środowiska | Karty sesji poprzedniego środowiska trwają w tle i odtwarzają pełny stan przy powrocie |
| Zamknięcie okna aplikacji | Procesy sesji pracują dalej na serwerze; po ponownym uruchomieniu klient wraca do ostatnio aktywnej karty |
| Przejście na inne urządzenie | Ten sam stan platformy — synchronizacja idzie przez serwer na żywo |
| Aktywacja funkcji globalnej | Nic się nie zmienia — funkcje globalne nie tworzą przestrzeni roboczej |

Ta trwałość jest tym, co pozwala nadzorować z telefonu procesy uruchomione przy
biurku — a pętli pracy ciągłej pracować niezależnie od tego, czy okno aplikacji
jest otwarte.

---

## 18. Typowe scenariusze pracy

### 18.1. Dokument od materiału źródłowego do wersji końcowej

TalkIn → Studio. Wczytać dokument do Studio Editor, zlecić z Chat Window korektę
i zmianę stylu, porównać wersje w Diff/Grep Panel, odłożyć wynik magistralą do
Library. Jeżeli potrzebna jest wersja obcojęzyczna — przekazać artefakt do modułu
Translate i tłumaczyć równolegle na kilka języków z jednym glosariuszem.

### 18.2. Opracowanie wieloźródłowe

TalkIn → Browser → Research. Przeglądać strony wspólnie z AI w Browser Window,
odnotowywać fragmenty w Notes Panel, przekazać zebrane źródła do Research, tam
zgromadzić je w Sources Manager, prowadzić ustalenia w Findings Panel, złożyć
raport w Report Builder i wyeksportować przez Export Panel.

### 18.3. Decyzja wymagająca kilku perspektyw

TalkIn albo WorkSpace → Roundtable. Postawić to samo zagadnienie kilku modelom
w Model Panels, prześledzić wymianę argumentów w Debate Panel, ukierunkować
dyskusję z Moderator Panel, zebrać stanowisko w Consensus Panel.

### 18.4. Zmiana w kodzie z weryfikacją

CodeStudio → Developer → Terminal → Diagnostics. Nawigować po Project Tree,
edytować w Code Editor przy wsparciu Chat Window, zbudować i sprawdzić wynik
w Build Output, uruchomić polecenia w Terminal Tabs, a przy błędzie przejść do
Diagnostics — Logs Viewer, Errors Panel, Recommendations Panel — i wrócić
z poprawką do modułu Developer.

### 18.5. Cykliczny raport bez nadzoru

Strona główna → Automations. Zaprojektować proces w Workflow Builder, ustalić
cykliczność w Scheduler, wskazać moduł wykonujący pracę, obserwować kolejne
przebiegi w Execution Monitor. Powiadomienia doprowadzą do interwencji wtedy, gdy
przebieg będzie jej wymagał.

### 18.6. Projekt prowadzony zespołem modeli w trybie ciągłym

Strona główna → Agents (zbudować agentów) → MultitaskingAI. W sekcji Role
przypisać agentów do Executora 1, Executora 2, Coordinatora i Validatora,
w Kolejkach zdefiniować przepływ zadań, w Orkiestracji ustalić zależności,
w Harmonogramie wpiąć automatykę zapewniającą ciągłość. Ustalić w konfiguracji
zakres autonomii i punkty obowiązkowego zatwierdzenia, zapisać układ w sekcji
Zespoły. Przebieg obserwować w Monitorze procesu — przy biurku albo z telefonu.

---

## 19. Słowniczek

| Pojęcie | Znaczenie |
|---|---|
| **Operator** | Właściciel konta prowadzący pracę na platformie |
| **Środowisko** | Najwyższy poziom organizacji pracy — „w jakim trybie pracuję?" |
| **Moduł** | Wyspecjalizowany obszar roboczy wewnątrz środowiska — „jakie zadanie wykonuję?" |
| **Okno operacyjne** | Właściwa przestrzeń wykonania zadania wewnątrz modułu |
| **Karta sesji** | Jednostka pracy równoległej z własnym układem, historią i kontekstem |
| **Boczna nawigacja** | Stała lista modułów dostępnych w środowisku |
| **Panel orkiestracji** | Boczna nawigacja środowiska MultitaskingAI — sześć sekcji sterowania zespołem |
| **Chat Window** | Okno kanału Użytkownik ↔ Wykonawca — centralny punkt sterowania platformą |
| **Execution Loop Window** | Okno kanału Koordynator ↔ Wykonawca — pętla wykonawcza i nadzór |
| **Wykonawca** | Model AI, agent albo system wykonawczy realizujący zadania |
| **Koordynator** | Komponent orkiestrujący: dekompozycja zlecenia, przydział zadań, nadzór pętli |
| **Subagent Network** | Do 15 podagentów uruchamianych przez wykonawcę |
| **Komponent własny** | Nazwany wytwór Operatora: automatyka, agent, projekt, profil asystenta |
| **Magistrala kontekstu** | Mechanizm przekazywania artefaktów między modułami |
| **Artefakt** | Wytwór pracy: dokument, raport, grafika, plik kodu — wersjonowany |
| **Punkt izolacji** | Ustawienie rozstrzygające, co jest współdzielone, a co odrębne |
| **Warstwa widoczności** | Poziom, na którym element interfejsu jest ujawniany (1–4) |
| **Pasek kontekstu** | Listwa znaczników bieżącego układu pracy: środowisko, projekt, model, wykonawca |

---

*Danaco Console — Platforma AI Workspace OS · v2.0*
*© 2026 Danaco Holding Group Sp. z o.o. Wszelkie prawa zastrzeżone — Dariusz Naharnowicz*
