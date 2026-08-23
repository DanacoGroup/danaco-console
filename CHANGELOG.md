# Danaco Console — Historia zmian

| | |
|---|---|
| **Produkt** | Danaco Console |
| **Rodzaj** | Platforma AI Workspace OS |
| **Producent** | Danaco Holding Group Sp. z o.o., ul. Gen. J. Hallera 81, 43-400 Cieszyn |
| **Wsparcie** | support@danaco-group.pl |
| **Twórca** | Dariusz Naharnowicz |
| **Wersja bieżąca** | v2.0 |
| **Data dokumentu** | 2026-08-18 |

Dokument opisuje zmiany widoczne dla Operatora: nowe możliwości, zmiany
zachowania i usunięcia funkcji. Wpis powstaje po udostępnieniu Wydania.
Szczegóły pracy z opisanymi funkcjami zawiera `INSTRUKCJA-UZYTKOWANIA.md`,
a zasady konfiguracji — `INSTALACJA-I-KONFIGURACJA.md`.

---

## v2.0 — 2026-08-18

Wydanie porządkujące platformę wokół czterech środowisk, piętnastu modułów
i dwóch kanałów komunikacji operacyjnej.

### Struktura pracy

- **Cztery środowiska** — TalkIn, WorkSpace, CodeStudio i MultitaskingAI —
  dostępne ze strony głównej pełniącej rolę centrum dowodzenia z trzema
  strefami: wybór środowiska, komponenty własne, ustawienia.
- **Piętnaście modułów** z własnymi oknami operacyjnymi; dostępność modułu
  w środowisku wyznacza macierz dostępności, a nie sztywne przypisanie.
- **Karty sesji** jako jednostka pracy równoległej — własny układ, historia
  i kontekst każdej karty, z konfigurowalnym zakresem współdzielenia.
- **Magistrala kontekstu** — przekazywanie artefaktów między modułami przez
  odwołanie, bez tworzenia kopii i nowych wersji.

### Sterowanie pracą

- **Chat Window** — kanał Użytkownik ↔ Wykonawca: polecenia w języku
  naturalnym, strumień odpowiedzi na żywo, zatwierdzanie i przerywanie działań,
  wyjaśnianie wyniku.
- **Execution Loop Window** — kanał Koordynator ↔ Wykonawca: dekompozycja
  zlecenia, kolejka zadań, kontrola jakości, decyzje o ponowieniu, wstrzymanie,
  wznowienie i korekta zlecenia.
- **MultitaskingAI** — zespół modeli w rolach (dwóch wykonawców, koordynator,
  rola kontrolna), Subagent Network do piętnastu podagentów, silnik kolejek
  z jedenastoma akcjami, warstwa orkiestracji zależności oraz panel orkiestracji
  z sześcioma sekcjami. Po spięciu z modułem Automations — pętla pracy ciągłej.

### Konfiguracja

- **Okno konfiguracji** z siedemnastoma zakresami ustawień i czterema warstwami
  konfiguracji: globalna → środowisko → projekt → sesja, z dziedziczeniem
  i zasadą „brak ustawienia = wartość odziedziczona".
- **Punkty izolacji** — izolacja kontekstu (historia, pamięć, kontekst) oraz
  osiem pozycji izolacji technicznej, ustawianych na siedmiu poziomach zasięgu,
  z profilami izolacji i podglądem polityki efektywnej.
- **Szybka zmiana konfiguracji** bieżącej sesji z menu kontekstowego okna
  operacyjnego, bez przerywania pracy.
- **Objaśnienia kontekstowe** przy każdej pozycji ustawienia — co robi i jaki
  ma wpływ na działanie platformy.

### Modele, agenci, rozszerzenia

- **Cztery kanały integracji modeli** — API, CLI, SSH i HTTP — z jednolitym
  strumieniem odpowiedzi niezależnym od kanału.
- **Rejestr kont modeli** z jawnym przełączaniem konta aktywnego; dane
  dostępowe przechowywane poza bazą danych.
- **Moduł Agents** — budowa agentów: tożsamość, instrukcje systemowe,
  umiejętności, konektory, pamięć i uprawnienia nadawane w Centrum uprawnień.
- **Rozszerzenia** czterech rodzajów (wtyczki, umiejętności, konektory, serwery
  MCP) w jednolitym kontrakcie, z dwoma źródłami: dostarczane z platformą oraz
  instalowane przez Operatora.

### Możliwości wykonawcze modułów

- **Treść i dokumenty** — konwersja dokumentów, odczyt treści z rozpoznawaniem
  pisma, porównywanie wersji, operacje kontekstowe na tekście.
- **Obraz i media** — przegląd, przekształcenia i konwersje obrazu,
  powiększanie, usuwanie tła, przetwarzanie materiałów wideo i audio.
- **Wiedza** — indeksowanie i wyszukiwanie po znaczeniu, nie po samych słowach.
- **Poczta** — praca ze skrzynką Operatora: odnajdywanie, odczyt z załącznikami,
  szkice i wysyłka; platforma jest klientem skrzynki, nigdy serwerem poczty.
- **Praca inżynierska** — kontrola wersji, konsole i procesy, silnik kontenerów,
  diagnostyka na podstawie logów i błędów.
- **Archiwa** — pakowanie i rozpakowywanie z zabezpieczeniem przed wyjściem
  poza katalog i przed nadmiernym rozpakowaniem.

### Dostęp i bezpieczeństwo

- **Jedno konto właściciela, wiele urządzeń** — rejestracja przy pierwszym
  uruchomieniu, logowanie kolejnych urządzeń, parowanie urządzeń przenośnych,
  własny token dostępu na urządzenie i zdalne odwołanie dostępu.
- **Metody uwierzytelniania** — hasło i adres e-mail uwierzytelniający aktywne
  domyślnie, PIN i Windows Hello do włączenia przez Operatora.
- **Wymóg logowania** jako ustawienie Operatora, z domyślnym stanem aktywnym
  w środowisku produkcyjnym.
- **Zabezpieczenie połączenia** — połączenie szyfrowane, kontrola pochodzenia
  żądania, token wiązany z sesją połączenia.
- **Dziennik audytu** decyzji podejmowanych w kanałach komunikacji operacyjnej;
  działania nieodwracalne wymagają zatwierdzenia Operatora.

### Funkcje globalne

- **Mobile** — pełny nadzór nad procesami z telefonu i tabletu: oba kanały
  komunikacji operacyjnej, przegląd i sterowanie zadaniami, decyzje procesu,
  powiadomienia wypychane, praca przy braku połączenia.
- **Always On Display** — agent towarzyszący obecny we wszystkich częściach
  platformy: proaktywne sugestie, pomoc kontekstowa, tor głosowy, nadzór nad
  procesem w środowisku MultitaskingAI.

### Warstwa wizualna i interfejs

- **System wizualny v2.0** — granat i złoto, motyw jasny i ciemny, zestaw 82
  ikon na jednolitej siatce, trzy kroje pisma z pełnym zestawem polskich znaków
  diakrytycznych.
- **Cztery warstwy widoczności** — funkcja niepotrzebna do bieżącego zadania
  pozostaje niewidoczna, a każda ukryta jest osiągalna jednym kliknięciem,
  skrótem klawiszowym albo poleceniem w języku naturalnym.
- **Układ wyłącznie kolumnowy** obszaru roboczego; regulacji podlega szerokość
  kolumn.

### Model wdrożenia

- **Serwer platformy i cienki klient** — cała praca obliczeniowa i stan trwały
  po stronie serwera, na urządzeniu Operatora wyłącznie okno aplikacji.
  Stanowisko nie wymaga układu GPU ani instalowania dodatkowych programów.
- **Pakiety klienckie** dla systemu Windows (instalator) oraz Linux (pakiet
  systemowy i wariant przenośny).
- **Aktualizacja platformy** wykonywana centralnie po stronie serwera;
  aktualizacja okna aplikacji z poziomu samej aplikacji, z kontrolą sumy
  kontrolnej pobranego pliku przed założeniem.
- **Synchronizacja na żywo** między urządzeniami powiązanymi z kontem.

---

## Wydania wcześniejsze

| Data | Zakres |
|---|---|
| 2026-08-15 | Możliwości wykonawcze modułów: obraz, media, dokumenty z rozpoznawaniem pisma, archiwa, indeks znaczeniowy wiedzy, poczta. Zabezpieczenie brzegu połączenia i tor połączeń z maszynami Operatora. Asystent głosowy prowadzący pracę w oknie docelowym. Pakiety instalacyjne produktu i aktualizacja z poziomu aplikacji |
| 2026-08-13 | Rozmowa dostępna w każdym module. Uporządkowanie warstwy sesji i okien komunikacji. Poprawki zgłoszone po uruchomieniach produktu |
| 2026-08-12 | Kontrakt komunikacji obejmujący piętnaście modułów, trwała historia sesji, okno konfiguracji sterowane danymi, rejestr kont i tożsamości modeli, punkty dostępu i nadania, warstwa wizualna v2.0 |
| 2026-08-11 | Fundament platformy: kontrakt komunikacji, schemat trwałości, protokół połączenia klient–serwer, warstwa rdzenia z sesjami i oknami komunikacji, powłoka natywna okna |

---

*Danaco Console — Platforma AI Workspace OS · v2.0*
*© 2026 Danaco Holding Group Sp. z o.o. Wszelkie prawa zastrzeżone — Dariusz Naharnowicz*
