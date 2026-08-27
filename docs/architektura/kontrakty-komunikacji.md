*Dokument specyfikuje kontrakty komunikacji Danaco Console: kanały, komunikaty, format WebSocket/JSON oraz reguły wymiany danych.*

# Danaco Console — Kontrakty komunikacji

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
| **Źródło** | Koncepcja platformy i architektura |

Dokument definiuje pełne kontrakty komunikacji między klientem (oknem aplikacji na urządzeniu Operatora) a rdzeniem (serwerem wykonawczym) platformy Danaco Console. Kontrakty realizują w warstwie komunikacyjnej strukturę koncepcji platformy oraz zasady techniczne ustalone w dokumencie Architektury: model wdrożenia hybrydowego, pojedynczy kanał WebSocket, serwer jako jedyne źródło prawdy i warstwową konfigurację. Zakres obejmuje stronę główną i jej trzy strefy, wybór środowiska, karty sesji, komponenty własne, piętnaście modułów według macierzy dostępności, oba kanały komunikacji operacyjnej — Chat Window (Użytkownik ↔ Wykonawca) i Execution Loop Window (Koordynator ↔ Wykonawca) — warstwy widoczności interfejsu, środowisko MultitaskingAI wraz z panelem orkiestracji, silnik kolejek, orkiestrację zależności, agentów, pamięć, konfigurację, okno konfiguracji punktów izolacji oraz rozszerzenia. Opracowanie ujmuje kontrakt w tabelach zestawczych, tabelach schematów ładunku, diagramach sekwencji oraz diagramach przepływu.

---

## Spis treści

- [Wprowadzenie](#wprowadzenie)
1. [Zasady ogólne komunikacji](#1-zasady-ogólne-komunikacji)
2. [Cykl życia połączenia](#2-cykl-życia-połączenia)
3. [Struktura komunikatu](#3-struktura-komunikatu)
4. [Kategorie komunikatów](#4-kategorie-komunikatów)
5. [Uwierzytelnianie i tożsamość urządzenia](#5-uwierzytelnianie-i-tożsamość-urządzenia)
6. [Strona główna — kontrakty trzech stref](#6-strona-główna--kontrakty-trzech-stref)
7. [Nawigacja i warstwy widoczności](#7-nawigacja-i-warstwy-widoczności)
   - 7.1. [Środowiska, moduły, okna kontekstowe](#71-środowiska-moduły-okna-kontekstowe)
   - 7.2. [Warstwy widoczności](#72-warstwy-widoczności)
8. [Karty sesji](#8-karty-sesji)
9. [Kanały komunikacji operacyjnej](#9-kanały-komunikacji-operacyjnej)
   - 9.1. [Chat Window — kanał Użytkownik ↔ Wykonawca](#91-chat-window--kanał-użytkownik--wykonawca)
   - 9.2. [Execution Loop Window — kanał Koordynator ↔ Wykonawca](#92-execution-loop-window--kanał-koordynator--wykonawca)
10. [Moduły — kontrakty](#10-moduły--kontrakty)
11. [MultitaskingAI — panel orkiestracji](#11-multitaskingai--panel-orkiestracji)
12. [Pamięć](#12-pamięć)
13. [Historia i retencja](#13-historia-i-retencja)
14. [Konfiguracja](#14-konfiguracja)
15. [Izolacja — okno konfiguracji punktów izolacji](#15-izolacja--okno-konfiguracji-punktów-izolacji)
16. [Funkcje globalne](#16-funkcje-globalne)
17. [Rozszerzenia](#17-rozszerzenia)
18. [Synchronizacja wielourządzeniowa](#18-synchronizacja-wielourządzeniowa)
19. [Obsługa błędów](#19-obsługa-błędów)
20. [Wersjonowanie kontraktu](#20-wersjonowanie-kontraktu)
- [Załącznik A. Zbiorcza tabela poleceń](#załącznik-a-zbiorcza-tabela-poleceń)
- [Załącznik B. Zbiorcza tabela zdarzeń](#załącznik-b-zbiorcza-tabela-zdarzeń)
- [Załącznik C. Przykładowe wymiany komunikatów](#załącznik-c-przykładowe-wymiany-komunikatów)
- [Załącznik D. Kody błędów](#załącznik-d-kody-błędów)
- [Załącznik E. Słowniczek pojęć kontraktu](#załącznik-e-słowniczek-pojęć-kontraktu)
- [Załącznik F. Schematy ładunku wybranych komunikatów](#załącznik-f-schematy-ładunku-wybranych-komunikatów)

---

## Wprowadzenie

Koncepcja platformy opisuje Danaco Console jako system operacyjny do pracy z modelami sztucznej inteligencji, zbudowany na trójstopniowej hierarchii środowisko → moduł → okno operacyjne, ponad którą działają funkcje globalne i agenci. Dokument Architektury ustala, że cała ta struktura komunikuje się jednym, dwukierunkowym kanałem WebSocket, a serwer pozostaje jedynym źródłem prawdy o stanie systemu. Niniejszy dokument spina oba źródła w jeden kontrakt: dla każdego obszaru funkcjonalnego koncepcji — strony głównej, wyboru środowiska, kart sesji, komponentów własnych, piętnastu modułów, środowiska MultitaskingAI, silnika kolejek, orkiestracji, agentów, pamięci, konfiguracji, izolacji i rozszerzeń — definiuje odpowiadające mu polecenia, zdarzenia i strumienie.

Kontrakt jest zgodny z zasadą nadrzędną pełnej konfigurowalności (koncepcja, rozdz. 14, zasada 2): żadne polecenie opisane w tym dokumencie nie wymusza zachowania, które nie byłoby jednocześnie dostępne do zmiany z poziomu okna konfiguracji. Brak ustawienia oznacza wartość domyślną, a nie ograniczenie możliwości — zasada ta obowiązuje w warstwie komunikacji tak samo, jak w warstwie funkcjonalnej i w warstwie danych.

Nazwy typów komunikatów zapisane są w notacji `obszar.zasób.akcja`, w języku angielskim, zgodnie z regulaminem projektu; opisy, nazwy pól o znaczeniu biznesowym oraz wszystkie objaśnienia pozostają w języku polskim. Nazwy własne środowisk, modułów, okien operacyjnych, ról i funkcji globalnych powtarzają dosłownie nazewnictwo koncepcji platformy.

---

## 1. Zasady ogólne komunikacji

Poniższe pięć zasad, ustalonych w dokumencie Architektury (rozdz. 11) i rozwiniętych tu do poziomu kontraktu, obowiązuje w całym dokumencie i nie jest powtarzane przy każdym rozdziale szczegółowym.

| Nr | Zasada | Reguła kontraktu |
|---|---|---|
| 1 | Jeden kanał | Cała komunikacja między klientem a serwerem odbywa się jednym połączeniem WebSocket, nawiązywanym raz na sesję połączenia urządzenia i utrzymywanym jako stałe, dwukierunkowe. |
| 2 | Format JSON | Każdy komunikat jest obiektem JSON zgodnym ze strukturą koperty opisaną w rozdziale 3. |
| 3 | Serwer jako jedyne źródło prawdy | Klient nigdy nie modyfikuje stanu lokalnie w oderwaniu od serwera — wysyła polecenia i odbiera zdarzenia oraz strumienie. Widok klienta jest zawsze odbiciem stanu serwera. |
| 4 | Rozgłaszanie zmian | Każda zmiana stanu wykonana na skutek polecenia jednego urządzenia jest rozgłaszana przez serwer w postaci zdarzenia do wszystkich urządzeń powiązanych z tym samym kontem (rozdz. 18). |
| 5 | Transmisja na żywo | Odpowiedzi modeli oraz zmiany synchronizowane między urządzeniami przekazywane są strumieniowo, fragment po fragmencie, bez oczekiwania na komplet danych. |

---

## 2. Cykl życia połączenia

Cykl życia połączenia obejmuje pięć etapów, zgodnych z rozdziałem 2 dokumentu Architektury, doprecyzowanych tu do poziomu wymienianych komunikatów.

| Etap | Inicjator | Komunikaty | Opis |
|---|---|---|---|
| 1. Nawiązanie połączenia | Klient | otwarcie kanału WebSocket | Klient otwiera kanał z adresem serwera skonfigurowanym w pakiecie klienta. |
| 2. Uzgodnienie i uwierzytelnienie | Klient → Serwer → Klient | `connection.hello`, `connection.hello.ack` | Klient przedstawia wersję protokołu, identyfikator urządzenia i token dostępu (o ile posiada); serwer potwierdza wersję, fazę wdrożenia i wynik uwierzytelnienia (rozdz. 5). |
| 3. Powiązanie | Klient → Serwer | `session.bind` albo `home.enter` | Klient zgłasza, czy wraca do aktywnej karty sesji, czy wchodzi na stronę główną; serwer wiąże połączenie z odpowiednim stanem. |
| 4. Praca | Klient ↔ Serwer | polecenia, zdarzenia, strumienie (rozdz. 4–17) | Właściwa wymiana komunikatów w toku pracy z platformą. |
| 5. Rozłączenie | Klient albo sieć | zamknięcie kanału WebSocket | Zamknięcie kanału nie kończy sesji ani procesów po stronie serwera (Architektura, rozdz. 7); ponowne połączenie odtwarza etapy 1–3 i podłącza klienta do stanu zastanego na serwerze. |

Poniższy diagram sekwencji przedstawia wymianę komunikatów w pięciu etapach cyklu życia połączenia.

```
Klient                                        Serwer
  │   otwarcie kanału WebSocket                 │   Etap 1
  │ ──────────────────────────────────────────►│   nawiązanie połączenia
  │                                             │
  │   connection.hello                          │   Etap 2
  │   (protocolVersion, deviceId, token)        │   uzgodnienie
  │ ──────────────────────────────────────────►│   i uwierzytelnienie
  │   response (serverVersion,                  │
  │   deploymentPhase, authenticated)           │
  │ ◄──────────────────────────────────────────│
  │                                             │
  │   session.bind albo home.enter              │   Etap 3
  │ ──────────────────────────────────────────►│   powiązanie
  │                                             │
  │   polecenia · zdarzenia · strumienie        │   Etap 4
  │ ◄─────────────────────────────────────────►│   praca
  │                                             │
  │   zamknięcie kanału WebSocket               │   Etap 5
  │ ──────────────────────────────────────────►│   rozłączenie (sesja
  │                                             │   i procesy trwają dalej)
```

Etap 3 rozstrzyga różnicę wynikającą z rozdziału 7 koncepcji: aplikacja nie otwiera od razu żadnego środowiska, lecz stronę główną — centrum dowodzenia — chyba że klient jednoznacznie zgłasza powrót do wcześniej otwartej karty sesji.

```json
// Klient → Serwer, etap 2

{ "type": "connection.hello", "id": "h-1", "payload": {
    "protocolVersion": "1.0.0", "deviceId": "dev-8841", "token": "•••" },
  "timestamp": "2026-08-06T09:00:00Z" }

// Serwer → Klient, etap 2

{ "type": "response", "id": "h-1", "status": "ok", "payload": {
    "protocolVersion": "1.0.0", "serverVersion": "0.4.2",
    "deploymentPhase": "build", "authenticated": true },
  "timestamp": "2026-08-06T09:00:00Z" }
```

---

## 3. Struktura komunikatu

Każdy komunikat wymieniany kanałem WebSocket ma wspólną kopertę. Pola obecne zależą od kategorii komunikatu (rozdz. 4).

| Pole | Obecność | Opis |
|---|---|---|
| `type` | zawsze | Typ komunikatu w notacji `obszar.zasób.akcja` (dla poleceń i zdarzeń) albo `response` / `stream.chunk` (dla odpowiedzi i strumieni). |
| `id` | zawsze | Identyfikator komunikatu. Dla polecenia — unikatowy identyfikator nadany przez klienta. Dla odpowiedzi i fragmentów strumienia — powiela `id` polecenia, które je zainicjowało (korelacja). Dla zdarzenia — własny identyfikator nadany przez serwer. |
| `payload` | zawsze | Ładunek właściwy komunikatu; obiekt JSON, może być pusty (`{}`). |
| `timestamp` | zawsze | Znacznik czasu w formacie ISO 8601, nadawany przez nadawcę komunikatu. |
| `status` | tylko w odpowiedziach | `"ok"` albo `"error"`. |
| `error` | tylko gdy `status = "error"` | Obiekt błędu — struktura w rozdziale 19. |
| `seq` | tylko w fragmentach strumienia | Numer porządkowy fragmentu, rosnący od 1. |
| `done` | tylko w fragmentach strumienia | `true` przy ostatnim fragmencie strumienia, w przeciwnym razie `false`. |

Odpowiedź na polecenie ma zawsze `type: "response"` — klient odróżnia, którego polecenia dotyczy, po zgodności pola `id`. Ten wzorzec obowiązuje jednolicie we wszystkich obszarach opisanych w rozdziałach 6–17 i nie jest powtarzany przy każdym poleceniu z osobna.

---

## 4. Kategorie komunikatów

Platforma rozróżnia trzy kategorie komunikatów wymienianych kanałem WebSocket: polecenia, zdarzenia i strumienie. Poniższa tabela zestawia ich właściwości; szczegółowy opis każdej kategorii zawierają podrozdziały 4.1–4.3.

| Właściwość | Polecenie | Zdarzenie | Strumień |
|---|---|---|---|
| Kierunek | klient → serwer | serwer → klient | serwer → klient |
| Inicjator | klient | serwer (samodzielnie) | serwer, w odpowiedzi na polecenie inicjujące |
| Wartość pola `type` | `obszar.zasób.akcja` | `obszar.zasób.akcja` | `stream.chunk` |
| Korelacja przez `id` | `id` nadane przez klienta | własny `id` nadany przez serwer | `id` polecenia inicjującego |
| Odpowiedź serwera | dokładnie jeden komunikat `response` | brak (komunikat jednokierunkowy) | strumień otwierany odpowiedzią `response` z polem `streamId` |
| Pola swoiste koperty | — | — | `seq`, `done` |
| Zakończenie | komunikat `response` | — | fragment z `done: true` |
| Zasięg dostarczenia | serwer, z wynikowym zdarzeniem (rozdz. 4.2) | wszystkie urządzenia powiązane z kontem | urządzenie odbierające dany strumień |
| Przykłady | `message.send`, `home.enter` | `session.changed`, `config.changed` | strumień odpowiedzi modelu, strumień synchronizacji stanu |

Poniższy diagram przedstawia ogólny wzorzec interakcji na przykładzie polecenia inicjującego pracę modelu: skorelowaną odpowiedź otwierającą strumień, fragmenty strumienia oraz zdarzenie rozgłoszone przez serwer do innego urządzenia tego samego konta.

```
Urządzenie A (klient)          Serwer            Urządzenie B (to samo konto)
      │                          │                          │
      │  polecenie (id = c-1)    │                          │
      │ ────────────────────────►│                          │
      │  response (id = c-1)     │                          │
      │ ◄────────────────────────│                          │
      │  stream.chunk (id = c-1, │                          │
      │  seq = 1 … done = true)  │                          │
      │ ◄────────────────────────│                          │
      │                          │   zdarzenie (*.changed)  │
      │                          │ ────────────────────────►│
```

### 4.1. Polecenia (klient → serwer)

Żądanie wykonania operacji. Serwer zawsze odpowiada dokładnie jednym komunikatem `response`, powiązanym przez `id`. Polecenie może dodatkowo skutkować jednym lub wieloma zdarzeniami rozgłoszonymi do wszystkich urządzeń (rozdz. 4.2), a w przypadku poleceń inicjujących pracę modelu — otwarciem strumienia (rozdz. 4.3).

### 4.2. Zdarzenia (serwer → klient)

Powiadomienie o zmianie stanu. Serwer inicjuje je samodzielnie, niezależnie od tego, czy zmianę wywołało polecenie tego samego urządzenia, polecenie innego urządzenia, proces działający w tle (na przykład harmonogram Automations) czy działanie agenta. Zdarzenie niesie stan wystarczający do odświeżenia widoku klienta bez dodatkowego zapytania (Architektura, rozdz. 11).

### 4.3. Strumienie (serwer → klient)

Ciąg fragmentów przekazywanych na żywo. Strumień otwiera odpowiedź `status: "ok"` na polecenie inicjujące (na przykład `message.send`), z ładunkiem zawierającym `streamId` równy `id` polecenia; kolejne fragmenty mają `type: "stream.chunk"` i to samo `id`. Ostatni fragment niesie `done: true`. Dwa rodzaje strumieni występują w platformie: strumień odpowiedzi Wykonawcy w Chat Window (rozdz. 9.1) oraz strumień synchronizacji stanu między urządzeniami (rozdz. 18).

---

## 5. Uwierzytelnianie i tożsamość urządzenia

Rozdział realizuje kontrakt komunikacyjny dla mechanizmu opisanego w rozdziale 13 dokumentu Architektury: rejestracja, logowanie, metody dodatkowe oraz fazowanie wdrożenia.

| Typ komunikatu | Kierunek | Opis |
|---|---|---|
| `auth.register` | polecenie | Rejestracja przy pierwszym uruchomieniu: login, adres e-mail uwierzytelniający, hasło; skutkuje wysłaniem potwierdzenia przez e-mail. |
| `auth.confirm` | polecenie | Potwierdzenie rejestracji kodem otrzymanym e-mailem. |
| `auth.login` | polecenie | Logowanie loginem i hasłem; odpowiedź niesie token dostępu. |
| `auth.method.add` | polecenie | Dodanie metody dodatkowej z okna konfiguracji: PIN, Windows Hello albo e-mail uwierzytelniający. |
| `auth.method.remove` | polecenie | Usunięcie metody dodatkowej. |
| `auth.password.reset` | polecenie | Odzyskanie konta i ustawienie nowego hasła przez adres e-mail uwierzytelniający. |
| `auth.token.refresh` | polecenie | Odświeżenie tokenu dostępu bez ponownego logowania. |
| `auth.changed` | zdarzenie | Zmiana metod uwierzytelniania konta, rozgłaszana do pozostałych urządzeń. |

Token uzyskany po uwierzytelnieniu jest przedstawiany w polu `token` komunikatu `connection.hello` (rozdz. 2) przy każdym kolejnym nawiązaniu połączenia. Zgodnie z fazowaniem wdrożenia ustalonym w architekturze, w fazie budowy pole `deploymentPhase` odpowiedzi na `connection.hello` przyjmuje wartość `"build"`, a serwer przepuszcza połączenie bez wymogu tokenu; w fazie produkcyjnej (`"production"`) token jest wymagany. Realizacja opiera się na jednym przełączniku trybu po stronie serwera — kontrakt komunikatów pozostaje identyczny w obu fazach.

| Faza wdrożenia | Wartość pola `deploymentPhase` | Token wymagany | Zachowanie serwera |
|---|---|---|---|
| Budowa | `build` | Nie | Serwer przepuszcza połączenie bez wymogu tokenu. |
| Produkcyjna | `production` | Tak | Serwer wymaga tokenu do nawiązania połączenia. |

---

## 6. Strona główna — kontrakty trzech stref

Strona główna jest centrum dowodzenia, do którego klient trafia po etapie powiązania (rozdz. 2), o ile nie zgłasza powrotu do aktywnej karty sesji. Polecenie `home.enter` pobiera stan wszystkich trzech stref jednym wywołaniem — zgodnie z podziałem koncepcji (rozdz. 7) na strefę wyboru środowiska, strefę komponentów własnych i strefę ustawień.

```json
// Klient → Serwer

{ "type": "home.enter", "id": "c-1", "payload": {}, "timestamp": "2026-08-06T09:00:05Z" }

// Serwer → Klient

{ "type": "response", "id": "c-1", "status": "ok", "payload": {
    "zone1Environments": [ /* rozdz. 6.1 */ ],
    "zone2Components": [ /* rozdz. 6.2 */ ],
    "zone3Settings": { /* rozdz. 6.3 */ } },
  "timestamp": "2026-08-06T09:00:05Z" }
```

Ładunek odpowiedzi na `home.enter` niesie stan trzech stref strony głównej w jednym wywołaniu.

| Pole ładunku | Strefa | Zawartość |
|---|---|---|
| `zone1Environments` | Strefa 1 — wybór środowiska | Karty czterech środowisk: TalkIn, WorkSpace, CodeStudio, MultitaskingAI (rozdz. 6.1). |
| `zone2Components` | Strefa 2 — komponenty własne | Komponenty własne użytkownika: Automations, Agents, Workspace, Assistant (rozdz. 6.2). |
| `zone3Settings` | Strefa 3 — ustawienia | Wejście do okna konfiguracji oraz funkcje globalne Mobile i Always On Display (rozdz. 6.3). |

### 6.1. Strefa 1 — wybór środowiska

Realizuje karty środowisk opisane w rozdziale 7.4 koncepcji: cztery duże karty wejścia do środowisk TalkIn, WorkSpace, CodeStudio i MultitaskingAI.

| Typ komunikatu | Kierunek | Opis |
|---|---|---|
| `environment.list` | polecenie | Wykaz czterech środowisk platformy wraz z metadanymi karty (nazwa, jednozdaniowy opis trybu pracy, godło). |
| `environment.enter` | polecenie | Wybór karty środowiska w strefie 1; otwiera przestrzeń roboczą środowiska opisaną w rozdziale 7. |

### 6.2. Strefa 2 — komponenty własne

Realizuje kafle komponentów własnych opisane w rozdziale 7.4 koncepcji: siatkę czterech kafli — Automations, Agents, Workspace, Assistant — z etykietą zorientowaną na tworzenie. Zgodnie z rozdziałem 2.3 koncepcji komponent własny jest nazwanym wytworem konkretnego użytkownika, odrębnym od modułu.

| Typ komunikatu | Kierunek | Opis |
|---|---|---|
| `component.list` | polecenie | Wykaz komponentów własnych użytkownika, z filtrem po rodzaju (`automation`, `agent`, `project`, `assistantProfile`). |
| `component.create` | polecenie | Utworzenie nowego komponentu własnego wskazanego rodzaju w oknie konfiguracji właściwym temu rodzajowi. |
| `component.update` | polecenie | Zmiana konfiguracji istniejącego komponentu własnego. |
| `component.delete` | polecenie | Usunięcie komponentu własnego. |
| `component.assign` | polecenie | Przypisanie utworzonego komponentu do modułu lub sesji, w której ma zastosowanie (na przykład wpięcie automatyki w sesji modułu Developer, zgodnie z rozdz. 7.2 koncepcji). |
| `component.changed` | zdarzenie | Zmiana komponentu własnego, rozgłaszana do wszystkich urządzeń. |

Rodzaj komponentu determinuje kształt pola `payload` w `component.create` i `component.update`: dla `automation` odpowiada zasobom opisanym w rozdziale 10.3, dla `agent` — w rozdziale 10.15, dla `project` — w rozdziale 10.2, dla `assistantProfile` — w rozdziale 10.10.

### 6.3. Strefa 3 — ustawienia

Realizuje listwę ustawień opisaną w rozdziale 7.4 koncepcji: wejście do okna konfiguracji oraz dwie funkcje globalne, Mobile i Always On Display.

| Typ komunikatu | Kierunek | Opis |
|---|---|---|
| `config.window.open` | polecenie | Otwarcie okna konfiguracji (rozdz. 14) z poziomu strefy 3. |
| `mobile.status.get` | polecenie | Stan funkcji globalnej Mobile (rozdz. 16.1). |
| `aod.status.get` | polecenie | Stan funkcji globalnej Always On Display (rozdz. 16.2). |

---

## 7. Nawigacja i warstwy widoczności

Rozdział realizuje trzy pierwsze kroki nawigacji ustalone w rozdziale 8 koncepcji: strona główna → środowisko, boczna nawigacja → moduł, moduł → okno kontekstowe. Krok czwarty, karty sesji, opisano w rozdziale 8.

Obszar roboczy każdego modułu układany jest wyłącznie pionowo, w kolumnach sąsiadujących poziomo: lewa kolumna stała mieści Chat Window i otwierane obok Execution Loop Window (rozdz. 9), prawa kolumna dominująca mieści okna operacyjne modułu, a okna pomocnicze otwierają się jako kolejne kolumny boczne po prawej stronie obszaru roboczego. Kontrakt nie przewiduje komunikatów sterujących innym rozmieszczeniem — polecenia układu operują wyłącznie szerokością kolumn.

### 7.1. Środowiska, moduły, okna kontekstowe

| Typ komunikatu | Kierunek | Opis |
|---|---|---|
| `module.list` | polecenie | Wykaz modułów dostępnych w bocznej nawigacji wskazanego środowiska, zgodnie z macierzą dostępności modułów (koncepcja, rozdz. 10). Dla środowiska MultitaskingAI zwraca zamiast tego strukturę panelu orkiestracji (rozdz. 11.7 niniejszego dokumentu, koncepcja rozdz. 13.8). |
| `workspace.enter` | polecenie | Wejście do wybranego modułu; tworzy nową kartę sesji (rozdz. 8) i buduje zestaw okien operacyjnych właściwy temu modułowi (Załącznik A koncepcji). |
| `window.state.get` | polecenie | Odczyt bieżącego stanu wskazanego okna operacyjnego otwartej karty sesji. |
| `window.action` | polecenie | Wzorzec ogólny wywołania akcji w oknie operacyjnym, gdy akcja nie ma dedykowanego typu komunikatu w rozdziale 10; ładunek niesie `windowId`, `action` i `params`. |
| `window.state.changed` | zdarzenie | Zmiana stanu okna operacyjnego, rozgłaszana do urządzeń powiązanych z daną kartą sesji. |

`module.list` przyjmuje parametr `environmentId`; zwrócony wykaz jest podzbiorem piętnastu modułów platformy wyznaczonym przez macierz dostępności — moduł Automations nie występuje w żadnym wykazie modułowym, ponieważ zgodnie z rozdziałem 10 koncepcji nie ma własnego okna w bocznej nawigacji środowisk i pozostaje dostępny wyłącznie przez stronę główną (rozdz. 6.2) oraz jako komponent własny wpięty w sesję.

```json
// Klient → Serwer

{ "type": "module.list", "id": "c-2", "payload": { "environmentId": "codestudio" },
  "timestamp": "2026-08-06T09:01:00Z" }

// Serwer → Klient

{ "type": "response", "id": "c-2", "status": "ok", "payload": {
    "modules": ["workspace", "roundtable", "design", "terminal", "developer",
                "diagnostics", "apps", "agents"] },
  "timestamp": "2026-08-06T09:01:00Z" }
```


### 7.2. Warstwy widoczności

Każdy element interfejsu należy do dokładnie jednej z czterech warstw widoczności; przynależność elementu do warstwy jest stanem serwera i podlega tym samym zasadom kontraktu, co pozostałe obszary — klient odczytuje ją poleceniem, a zmiana jest rozgłaszana zdarzeniem. Profil warstw wyznaczany jest rolą użytkownika: elementy warstwy 4 nie są zwracane klientowi działającemu w roli podstawowej.

| Warstwa | Nazwa | Sposób dostępu w kontrakcie |
|---|---|---|
| 1 | Zawsze widoczna | Elementy zwracane w stanie okna bez dodatkowego wywołania. |
| 2 | Widoczna na żądanie | Element ujawniany poleceniem `layer.element.get` po użyciu ikony, przycisku, przełącznika albo znacznika kontekstowego. |
| 3 | Rozwinięcia kontekstowe | Element ujawniany poleceniem `layer.element.get` przy rozwinięciu menu kebab, menu hamburger, menu kontekstowego, panelu popover albo listy rozwijanej. |
| 4 | Funkcje eksperckie | Element wywoływany poleceniem `layer.function.invoke` — poleceniem języka naturalnego w oknie komunikacji, skrótem klawiszowym albo wyborem z wyszukiwarki funkcji. |

| Typ komunikatu | Kierunek | Opis |
|---|---|---|
| `layer.element.get` | polecenie | Pobranie przypisania wskazanego elementu interfejsu do warstwy widoczności; ładunek niesie `elementId`, odpowiedź — `layer`, `trigger` i `visible` dla bieżącej roli. |
| `layer.element.set` | polecenie | Zapis przypisania elementu interfejsu do warstwy widoczności; ładunek niesie `elementId`, `layer` oraz `trigger`. |
| `layer.profile.get` | polecenie | Odczyt profilu warstw obowiązującego dla roli użytkownika — wykazu elementów widocznych bez interakcji oraz elementów ujawnianych wyzwalaczami warstw 2–3. |
| `layer.search.query` | polecenie | Zapytanie wyszukiwarki funkcji; ładunek niesie frazę `query` i zakres `scope` (środowisko, moduł, karta sesji), odpowiedź — dopasowane funkcje wraz z ich warstwą i sposobem wywołania. |
| `layer.function.invoke` | polecenie | Wywołanie funkcji warstwy 4 poleceniem języka naturalnego, skrótem klawiszowym albo pozycją wyszukiwarki funkcji; ładunek niesie `functionId` albo `naturalCommand` oraz `invocation` (`naturalLanguage`, `shortcut`, `search`). |
| `layer.profile.changed` | zdarzenie | Zmiana przypisania elementu do warstwy albo profilu warstw roli, rozgłaszana do wszystkich urządzeń konta. |

Zasada jednego kliknięcia obowiązuje w kontrakcie wprost: każda funkcja warstw 2–4 osiągalna jest pojedynczym poleceniem `layer.element.get` albo `layer.function.invoke`, bez łańcucha wywołań pośrednich.

```json
// Klient → Serwer — odczyt profilu warstw dla roli użytkownika

{ "type": "layer.profile.get", "id": "c-60", "payload": {
    "sessionId": "s-501", "roleId": "operator" },
  "timestamp": "2026-08-06T09:20:00Z" }

// Serwer → Klient

{ "type": "response", "id": "c-60", "status": "ok", "payload": {
    "roleId": "operator",
    "layer1": ["chatWindow", "executionLoopWindow", "contextBar", "moduleNav"],
    "layer2": ["modelSelector", "executorSelector", "effortLevel"],
    "layer3": ["operationsMenu", "quickSettings"],
    "layer4": [] },
  "timestamp": "2026-08-06T09:20:00Z" }

// Klient → Serwer — wywołanie funkcji warstwy 4 poleceniem języka naturalnego

{ "type": "layer.function.invoke", "id": "c-61", "payload": {
    "sessionId": "s-501", "naturalCommand": "otwórz diagnostykę niskiego poziomu",
    "invocation": "naturalLanguage" },
  "timestamp": "2026-08-06T09:20:30Z" }
```

Warunki błędu rozdziału: `not_found` — wskazany `elementId` albo `functionId` nie istnieje; `permission_denied` — rola użytkownika nie obejmuje funkcji warstwy 4 wskazanej w `layer.function.invoke`; `validation_failed` — wartość `layer` spoza zakresu 1–4 albo brak zarówno `functionId`, jak i `naturalCommand`; `command_not_understood` — treść `naturalCommand` nie została dopasowana do żadnej funkcji.

---

## 8. Karty sesji

Rozdział realizuje mechanizm kart sesji opisany w rozdziale 8 koncepcji: aktywne sesje trzymane są w kartach poziomych w oknie środowiska, a przełączanie kart oznacza przełączanie między sesjami — mechanika analogiczna do kart w przeglądarce internetowej.

| Typ komunikatu | Kierunek | Opis |
|---|---|---|
| `session.create` | polecenie | Utworzenie nowej karty sesji dla wskazanego środowiska i modułu; alternatywna droga do tego samego skutku co `workspace.enter` z jawnym utworzeniem karty. |
| `session.list` | polecenie | Wykaz kart sesji aktywnych na koncie, z metadanymi (środowisko, moduł, projekt, tytuł, znacznik ostatniej aktywności). |
| `session.open` | polecenie | Otwarcie zapisanej, wcześniej zamkniętej karty sesji, wraz z zapisaną migawką układu okien (Model danych, rozdz. 4.1). |
| `session.focus` | polecenie | Oznaczenie karty sesji jako aktywnej na danym urządzeniu; nie wpływa na stan pozostałych urządzeń. |
| `session.update` | polecenie | Zmiana atrybutów karty sesji (tytuł, przypisany projekt, przypisany profil izolacji). |
| `session.bind` | polecenie | Powiązanie połączenia z aktywną kartą sesji przy etapie 3 cyklu życia połączenia (rozdz. 2), przy powrocie klienta do wcześniej otwartej karty. |
| `session.close` | polecenie | Zamknięcie karty sesji; zapis do powrotu zależny od ustawienia konfiguracyjnego (Model danych, rozdz. 4.1) — proces sesji po stronie serwera trwa dalej niezależnie od zamknięcia karty po stronie klienta (Architektura, rozdz. 7). |
| `session.delete` | polecenie | Trwałe usunięcie sesji wraz z jej historią, zgodnie z zasadą pełnego dostępu do operacji na danych (Architektura, rozdz. 10). |
| `session.changed` | zdarzenie | Zmiana karty sesji (utworzenie, aktualizacja, zamknięcie), rozgłaszana do wszystkich urządzeń. |
| `session.focus.changed` | zdarzenie | Zmiana karty aktywnej na innym urządzeniu tego samego konta — informacyjne, nie wymusza przełączenia u odbiorcy. |

Zakres, w jakim karty sesji współdzielą kontekst, historię i pamięć, jest konfigurowalny na poziomie zasięgu „karta sesji” opisanym w rozdziale 15 — `session.create` i `session.update` przyjmują pole `isolationProfileId`, którego pominięcie oznacza dziedziczenie profilu z zasięgu szerszego (rozdz. 15.3).

---

## 9. Kanały komunikacji operacyjnej

Platforma prowadzi pracę operacyjną dwoma kanałami komunikacji, obecnymi w każdym środowisku i każdym module jako elementy pierwszoplanowe architektury. Kanał pierwszy — Chat Window — łączy Użytkownika z Wykonawcą i zajmuje lewą kolumnę obszaru roboczego, stałą i o pełnej wysokości. Kanał drugi — Execution Loop Window — prezentuje komunikację Koordynatora z Wykonawcą i otwiera się jako kolumna sąsiadująca z oknem pierwszym. Oba kanały korzystają z tej samej koperty komunikatu (rozdz. 3) i tych samych trzech kategorii komunikatów (rozdz. 4); różnią się rolami stron i zakresem ładunku.

| Kanał | Okno | Strony | Zakres kontraktu |
|---|---|---|---|
| Kanał pierwszy | Chat Window | Użytkownik ↔ Wykonawca | Przesłanie polecenia, strumień odpowiedzi, zatwierdzenie i przerwanie działania, żądanie wyjaśnienia wyniku, zmiana kontekstu pracy (rozdz. 9.1). |
| Kanał drugi | Execution Loop Window | Koordynator ↔ Wykonawca | Przyjęcie zlecenia, dekompozycja na zadania, przydział zadania, zgłoszenie postępu, zwrot wyniku, wynik kontroli jakości, decyzja o ponowieniu, zamknięcie zlecenia, sterowanie przebiegiem (rozdz. 9.2). |

Role kontraktu: **Użytkownik** zleca i zatwierdza; **Koordynator** jest komponentem orkiestrującym platformy, dekomponuje zlecenie, przydziela i nadzoruje zadania; **Wykonawca** — AI, agent albo system wykonawczy — realizuje zadania. Kanał drugi jest uruchamiany zleceniem powstałym z polecenia przesłanego kanałem pierwszym; oba kanały prowadzone są jednym połączeniem WebSocket i odnoszą się do tej samej karty sesji.

### 9.1. Chat Window — kanał Użytkownik ↔ Wykonawca

Chat Window jest głównym oknem komunikacji i podstawowym mechanizmem sterowania wszystkimi procesami platformy. Zgodnie z rozdziałem 2.4 koncepcji stanowi jedną instancję kontraktu, rekonfigurowaną kontekstem bieżącej karty sesji, a nie osobny kontrakt dla każdego modułu.

| Typ komunikatu | Kierunek | Opis |
|---|---|---|
| `message.send` | polecenie | Przesłanie polecenia Użytkownika do Wykonawcy w kontekście aktywnej karty sesji i okna operacyjnego; ładunek niesie `sessionId`, `windowId` i `content`. Odpowiedź `status: "ok"` niesie `streamId` otwieranego strumienia odpowiedzi. |
| `stream.chunk` | strumień | Pojedynczy fragment strumienia odpowiedzi Wykonawcy; `payload` niesie `delta` (przyrost treści) oraz `toolCall`, gdy fragment reprezentuje wywołanie narzędzia zamiast treści tekstowej. |
| `message.approval.request` | zdarzenie | Żądanie zatwierdzenia działania, kierowane przez Wykonawcę do Użytkownika przed wykonaniem operacji wymagającej zgody; ładunek niesie `approvalId`, `action` i `impact`. |
| `message.approve` | polecenie | Zatwierdzenie albo odrzucenie działania przez Użytkownika; ładunek niesie `approvalId` oraz `decision` (`approve`, `reject`). |
| `message.stop` | polecenie | Przerwanie działania Wykonawcy i trwającego strumienia odpowiedzi; ładunek niesie `sessionId` i `streamId`. |
| `message.explain` | polecenie | Żądanie wyjaśnienia wyniku: uzasadnienia odpowiedzi, użytych źródeł, wywołanych narzędzi i przyjętych założeń; ładunek niesie `messageId` oraz `aspect` (`reasoning`, `sources`, `tools`, `assumptions`). Odpowiedź otwiera strumień wyjaśnienia. |
| `chat.context.set` | polecenie | Zmiana kontekstu pracy kanału: środowisko, moduł, projekt, okno operacyjne, model i wykonawca; ładunek niesie `sessionId` oraz zmieniane pola kontekstu. |
| `chat.context.get` | polecenie | Odczyt bieżącego kontekstu pracy kanału, prezentowanego jako znaczniki kontekstowe w interfejsie. |
| `message.changed` | zdarzenie | Zapis kompletnej wiadomości do historii po zakończeniu strumienia (`done: true`). |
| `chat.context.changed` | zdarzenie | Zmiana kontekstu pracy kanału, rozgłaszana do wszystkich urządzeń konta. |
| `message.list` | polecenie | Odczyt wiadomości bieżącej sesji (skrót do `history.load` ograniczony do rozmowy — rozdz. 13). |

```json
// Klient → Serwer — przesłanie polecenia

{ "type": "message.send", "id": "c-9", "payload": {
    "sessionId": "s-501", "windowId": "chat", "content": "Podsumuj ten dokument." },
  "timestamp": "2026-08-06T09:05:00Z" }

// Serwer → Klient — otwarcie strumienia odpowiedzi

{ "type": "response", "id": "c-9", "status": "ok",
  "payload": { "streamId": "c-9" }, "timestamp": "2026-08-06T09:05:00Z" }

// Serwer → Klient — fragmenty strumienia

{ "type": "stream.chunk", "id": "c-9", "seq": 1, "done": false,
  "payload": { "delta": "Dokument opisuje" }, "timestamp": "2026-08-06T09:05:01Z" }
{ "type": "stream.chunk", "id": "c-9", "seq": 2, "done": true,
  "payload": { "delta": " trzy główne wnioski." }, "timestamp": "2026-08-06T09:05:03Z" }

// Serwer → Klient — żądanie zatwierdzenia działania

{ "type": "message.approval.request", "id": "e-9110", "payload": {
    "sessionId": "s-501", "approvalId": "ap-4",
    "action": "developer.file.save", "impact": "Nadpisanie 3 plików w repozytorium." },
  "timestamp": "2026-08-06T09:05:20Z" }

// Klient → Serwer — zatwierdzenie działania

{ "type": "message.approve", "id": "c-10", "payload": {
    "approvalId": "ap-4", "decision": "approve" },
  "timestamp": "2026-08-06T09:05:25Z" }

// Klient → Serwer — żądanie wyjaśnienia wyniku

{ "type": "message.explain", "id": "c-11", "payload": {
    "messageId": "m-882", "aspect": "sources" },
  "timestamp": "2026-08-06T09:06:00Z" }

// Klient → Serwer — zmiana kontekstu pracy

{ "type": "chat.context.set", "id": "c-12", "payload": {
    "sessionId": "s-501", "moduleId": "diagnostics", "modelId": "fable-5",
    "executorId": "agt-14" },
  "timestamp": "2026-08-06T09:07:00Z" }
```

Poniższy diagram sekwencji przedstawia pełny cykl kanału pierwszego: przesłanie polecenia, strumień odpowiedzi, zatwierdzenie działania, przerwanie oraz zapis wiadomości do historii.

```
Użytkownik (Chat Window)                      Wykonawca (serwer)
  │   message.send (id = c-9)                   │
  │ ──────────────────────────────────────────►│
  │   response (streamId = c-9)                 │
  │ ◄──────────────────────────────────────────│
  │   stream.chunk (seq = 1, done = false)      │
  │ ◄──────────────────────────────────────────│
  │   message.approval.request (approvalId)     │
  │ ◄──────────────────────────────────────────│
  │   message.approve (decision = approve)      │
  │ ──────────────────────────────────────────►│
  │   stream.chunk (seq = 2, done = false)      │
  │ ◄──────────────────────────────────────────│
  │        (przerwanie przez użytkownika)       │
  │   message.stop ────────────────────────────►│
  │   stream.chunk (seq = n, done = true)       │
  │ ◄──────────────────────────────────────────│
  │   message.changed (zapis do historii)       │
  │ ◄──────────────────────────────────────────│
  │   message.explain (aspect = sources)        │
  │ ──────────────────────────────────────────►│
  │   stream.chunk (wyjaśnienie wyniku)         │
  │ ◄──────────────────────────────────────────│
```

Ładunek fragmentu strumienia (`stream.chunk`) ma następującą strukturę.

| Pole | Obecność | Opis |
|---|---|---|
| `delta` | zawsze | Przyrost treści odpowiedzi Wykonawcy względem fragmentów poprzednich. |
| `toolCall` | warunkowo | Wywołanie narzędzia, gdy fragment reprezentuje działanie zamiast treści tekstowej. |
| `seq` (koperta) | zawsze | Numer porządkowy fragmentu, rosnący od 1 (rozdz. 3). |
| `done` (koperta) | zawsze | `true` przy ostatnim fragmencie strumienia, w przeciwnym razie `false`. |

Warunki błędu kanału pierwszego: `not_found` — wskazany `sessionId`, `messageId` albo `approvalId` nie istnieje; `conflict` — `message.stop` dla strumienia już zamkniętego albo `message.approve` dla żądania już rozstrzygniętego; `permission_denied` — zatwierdzane działanie przekracza uprawnienia konta lub roli; `validation_failed` — pusty `content` albo wartość `decision` bądź `aspect` spoza dopuszczalnego zbioru; `channel_unavailable` — kanał integracji modelu jest chwilowo niedostępny, polecenie podlega ponowieniu.

Wybór modelu i kanału integracji (API, CLI, SSH, HTTP — Architektura, rozdz. 9) dla danej sesji lub roli ustala polecenie `model.channel.set` (rozdz. 14).

### 9.2. Execution Loop Window — kanał Koordynator ↔ Wykonawca

Execution Loop Window prezentuje pętlę wykonawczą: przyjęte zlecenie i jego dekompozycję na zadania, kolejkę i stan zadań, wymianę komunikatów sterujących, wyniki kontroli jakości i decyzje o ponowieniu, wskaźniki przebiegu pętli oraz sterowanie przebiegiem. Okno otwiera się jako kolumna sąsiadująca z Chat Window. Zlecenie powstaje z polecenia Użytkownika przesłanego kanałem pierwszym albo z harmonogramu pracy ciągłej (rozdz. 11.5); Koordynator przyjmuje je, dekomponuje i przydziela zadania Wykonawcom.

| Typ komunikatu | Kierunek | Opis |
|---|---|---|
| `loop.order.accept` | zdarzenie | Przyjęcie zlecenia przez Koordynatora; ładunek niesie `orderId`, `sessionId`, `source` (`chat`, `schedule`, `automation`) oraz `objective`. |
| `loop.order.decompose` | zdarzenie | Dekompozycja zlecenia na zadania; ładunek niesie `orderId` oraz wykaz `tasks` z `taskId`, `title`, `dependsOn` i `acceptanceCriteria`. |
| `loop.task.assign` | zdarzenie | Przydział zadania Wykonawcy; ładunek niesie `taskId`, `executorId`, `role` oraz `inputs`. |
| `loop.task.progress` | zdarzenie | Zgłoszenie postępu przez Wykonawcę; ładunek niesie `taskId`, `progress` (0–100), `stage` i `note`. |
| `loop.task.result` | zdarzenie | Zwrot wyniku zadania przez Wykonawcę; ładunek niesie `taskId`, `status` (`succeeded`, `failed`), `artifacts` oraz `summary`. |
| `loop.quality.result` | zdarzenie | Wynik kontroli jakości wykonanej przez rolę weryfikującą; ładunek niesie `taskId`, `verdict` (`accepted`, `rejected`), `findings` oraz `reviewerId`. |
| `loop.task.retry` | zdarzenie | Decyzja Koordynatora o ponowieniu zadania po odrzuceniu wyniku; ładunek niesie `taskId`, `attempt`, `reason` oraz skorygowane `inputs`. |
| `loop.order.close` | zdarzenie | Zamknięcie zlecenia; ładunek niesie `orderId`, `outcome` (`completed`, `aborted`), `deliverables` oraz `duration`. |
| `loop.subscribe` | polecenie | Subskrypcja podglądu pętli wykonawczej wskazanej karty sesji; otwiera dostarczanie zdarzeń pętli do okna. |
| `loop.state.get` | polecenie | Odczyt bieżącego stanu pętli: zlecenia, zadań, ich statusów i wskaźników przebiegu. |
| `loop.order.submit` | polecenie | Skierowanie zlecenia do Koordynatora bez pośrednictwa Chat Window; ładunek niesie `sessionId` i `objective`. |
| `loop.control` | polecenie | Sterowanie przebiegiem pętli; ładunek niesie `orderId` oraz `action`: `pause` (wstrzymanie), `resume` (wznowienie), `abort` (przerwanie), `amend` (korekta zlecenia wraz z polem `objective` albo `tasks`). |
| `loop.changed` | zdarzenie | Zmiana stanu pętli wykonawczej — zlecenia, zadania albo przebiegu — rozgłaszana do wszystkich urządzeń konta. |

```json
// Klient → Serwer — subskrypcja okna pętli wykonawczej

{ "type": "loop.subscribe", "id": "c-50", "payload": { "sessionId": "s-620" },
  "timestamp": "2026-08-06T09:15:00Z" }

// Serwer → Klient — przyjęcie zlecenia przez Koordynatora

{ "type": "loop.order.accept", "id": "e-7001", "payload": {
    "orderId": "ord-12", "sessionId": "s-620", "source": "chat",
    "objective": "Przygotuj moduł raportowania wraz z testami." },
  "timestamp": "2026-08-06T09:15:02Z" }

// Serwer → Klient — dekompozycja na zadania

{ "type": "loop.order.decompose", "id": "e-7002", "payload": {
    "orderId": "ord-12", "tasks": [
      { "taskId": "t-1", "title": "Model danych raportu", "dependsOn": [],
        "acceptanceCriteria": "Schemat zgodny z Model danych, rozdz. 4." },
      { "taskId": "t-2", "title": "Testy jednostkowe", "dependsOn": ["t-1"],
        "acceptanceCriteria": "Pokrycie ścieżek granicznych." } ] },
  "timestamp": "2026-08-06T09:15:03Z" }

// Serwer → Klient — przydział zadania Wykonawcy

{ "type": "loop.task.assign", "id": "e-7003", "payload": {
    "taskId": "t-1", "executorId": "agt-14", "role": "executor1",
    "inputs": { "projectId": "prj-3" } },
  "timestamp": "2026-08-06T09:15:04Z" }

// Serwer → Klient — zgłoszenie postępu

{ "type": "loop.task.progress", "id": "e-7004", "payload": {
    "taskId": "t-1", "progress": 60, "stage": "implementacja",
    "note": "Schemat encji ukończony." },
  "timestamp": "2026-08-06T09:16:00Z" }

// Serwer → Klient — zwrot wyniku

{ "type": "loop.task.result", "id": "e-7005", "payload": {
    "taskId": "t-1", "status": "succeeded", "artifacts": ["schema.sql"],
    "summary": "Model danych raportu gotowy." },
  "timestamp": "2026-08-06T09:17:00Z" }

// Serwer → Klient — wynik kontroli jakości

{ "type": "loop.quality.result", "id": "e-7006", "payload": {
    "taskId": "t-1", "verdict": "rejected", "reviewerId": "agt-21",
    "findings": ["Brak indeksu na kolumnie okresu raportowego."] },
  "timestamp": "2026-08-06T09:17:30Z" }

// Serwer → Klient — decyzja o ponowieniu

{ "type": "loop.task.retry", "id": "e-7007", "payload": {
    "taskId": "t-1", "attempt": 2, "reason": "Uwagi kontroli jakości",
    "inputs": { "projectId": "prj-3", "addIndex": "period" } },
  "timestamp": "2026-08-06T09:17:35Z" }

// Klient → Serwer — sterowanie przebiegiem: korekta zlecenia

{ "type": "loop.control", "id": "c-51", "payload": {
    "orderId": "ord-12", "action": "amend",
    "objective": "Przygotuj moduł raportowania wraz z testami i eksportem CSV." },
  "timestamp": "2026-08-06T09:18:00Z" }

// Serwer → Klient — zamknięcie zlecenia

{ "type": "loop.order.close", "id": "e-7008", "payload": {
    "orderId": "ord-12", "outcome": "completed",
    "deliverables": ["schema.sql", "tests/", "export_csv.py"], "duration": "PT42M" },
  "timestamp": "2026-08-06T09:52:00Z" }
```

Poniższy diagram sekwencji przedstawia pełny obieg pętli wykonawczej między Koordynatorem a Wykonawcą, wraz z kontrolą jakości, ponowieniem i sterowaniem przebiegiem ze strony Użytkownika.

```
Użytkownik        Koordynator                         Wykonawca
     │   loop.order.submit / message.send  │                │
     │ ──────────────────────────────────►│                │
     │                 loop.order.accept   │                │
     │ ◄──────────────────────────────────│                │
     │                 loop.order.decompose│                │
     │ ◄──────────────────────────────────│                │
     │                                     │ loop.task.assign
     │                                     │ ──────────────►│
     │                                     │ loop.task.progress
     │                                     │ ◄──────────────│
     │                                     │ loop.task.result
     │                                     │ ◄──────────────│
     │                 loop.quality.result │                │
     │ ◄──────────────────────────────────│                │
     │                                     │ loop.task.retry
     │                                     │ ──────────────►│
     │   loop.control (pause · resume ·    │                │
     │   abort · amend)                    │                │
     │ ──────────────────────────────────►│                │
     │                 loop.order.close    │                │
     │ ◄──────────────────────────────────│                │
```

Stany zadania w pętli wykonawczej tworzą zamknięty zbiór, wspólny dla wskaźników przebiegu i dla ładunku `loop.state.get`.

| Stan zadania | Wyzwalany komunikatem | Znaczenie |
|---|---|---|
| `assigned` | `loop.task.assign` | Zadanie przydzielone Wykonawcy, praca nierozpoczęta. |
| `running` | `loop.task.progress` | Zadanie w toku realizacji. |
| `review` | `loop.task.result` | Wynik zwrócony, kontrola jakości w toku. |
| `accepted` | `loop.quality.result` (`accepted`) | Wynik przyjęty przez rolę weryfikującą. |
| `retrying` | `loop.task.retry` | Zadanie skierowane do ponowienia po odrzuceniu wyniku. |
| `paused` | `loop.control` (`pause`) | Przebieg wstrzymany decyzją Użytkownika. |
| `aborted` | `loop.control` (`abort`) | Przebieg przerwany decyzją Użytkownika. |
| `closed` | `loop.order.close` | Zlecenie zamknięte, zadanie rozliczone. |

Warunki błędu kanału drugiego: `not_found` — wskazany `orderId` albo `taskId` nie istnieje; `conflict` — `loop.control` z akcją `resume` dla przebiegu niewstrzymanego, akcja `pause` dla przebiegu zamkniętego albo `amend` dla zlecenia zamkniętego; `permission_denied` — konto nie ma uprawnienia sterowania przebiegiem wskazanego zlecenia; `validation_failed` — wartość `action` spoza zbioru `pause`, `resume`, `abort`, `amend`, brak `objective` przy akcji `amend` albo puste `objective` w `loop.order.submit`; `rate_limited` — przekroczona częstość poleceń sterujących; `channel_unavailable` — kanał integracji Wykonawcy chwilowo niedostępny, zadanie podlega ponowieniu.

Relacja obu kanałów do panelu orkiestracji środowiska MultitaskingAI opisana jest w rozdziale 11: role Koordynatora i Wykonawców obsadza `role.assign` (rozdz. 11.2), kolejność i zależności zadań — kontrakty kolejek i orkiestracji (rozdz. 11.3, 11.4), a nadzór nad przebiegiem — monitor procesu (rozdz. 11.6).

---

## 10. Moduły — kontrakty

### 10.0. Wzorzec wspólny

Każdy moduł otwiera własny zestaw okien operacyjnych (Załącznik A koncepcji) i własny kontekst obu kanałów komunikacji operacyjnej (rozdz. 9). Kontrakty poniższych piętnastu podrozdziałów obejmują wyłącznie komunikaty specyficzne dla danego modułu, wynikające z jego zestawu funkcjonalności (koncepcja, rozdz. 11); operacje wspólne — otwarcie okna, odczyt jego stanu, akcja bez dedykowanego typu, zmiana stanu — realizuje wzorzec ogólny `window.state.get` / `window.action` / `window.state.changed` opisany w rozdziale 7. Moduły dostępne są w środowiskach zgodnie z macierzą dostępności modułów (koncepcja, rozdz. 10); niniejszy rozdział nie powtarza tej macierzy.

### 10.1. Studio

| Typ komunikatu | Kierunek | Okno | Opis |
|---|---|---|---|
| `studio.document.open` | polecenie | Studio Editor | Wczytanie dokumentu (PDF, DOCX, TXT, Markdown) do edytora. |
| `studio.document.save` | polecenie | Studio Editor | Zapis bieżącej wersji dokumentu. |
| `studio.contextual.op` | polecenie | Tools Panel | Operacja kontekstowa AI na zaznaczeniu: korekta, przepisanie, zmiana stylu, streszczenie, rozwinięcie, tłumaczenie. |
| `studio.diff.compare` | polecenie | Diff/Grep Panel | Porównanie dwóch wersji dokumentu. |
| `studio.repository.list` | polecenie | Session Repository | Wykaz wersji zapisanych w toku pracy nad dokumentem. |
| `studio.repository.restore` | polecenie | Session Repository | Przywrócenie wskazanej wersji. |
| `studio.document.changed` | zdarzenie | Studio Editor | Zmiana treści dokumentu, w tym w wyniku operacji kontekstowej AI. |

### 10.2. Workspace

| Typ komunikatu | Kierunek | Okno | Opis |
|---|---|---|---|
| `workspace.dashboard.get` | polecenie | Project Dashboard | Zbiorczy stan projektu. |
| `workspace.instructions.set` | polecenie | Instructions Panel | Zapis instrukcji systemowych właściwych projektowi. |
| `workspace.context.set` / `workspace.context.get` | polecenie | Context Memory | Zapis i odczyt pamięci kontekstowej projektu (Model danych, rozdz. 5). |
| `workspace.library.list` | polecenie | Project Library | Wykaz materiałów przypisanych do projektu — odpowiednik modułu Library (10.6) w zakresie ograniczonym do projektu. |
| `workspace.agent.assign` | polecenie | Agent Manager | Przypisanie agenta z modułu Agents (10.15) jako wykonawcy zadań w ramach projektu. |
| `workspace.project.changed` | zdarzenie | — | Zmiana dowolnego atrybutu projektu, w tym izolacji, widoczności w środowiskach lub przypisanych agentów. |

### 10.3. Automations

Moduł nie ma własnego okna w bocznej nawigacji żadnego środowiska (koncepcja, rozdz. 10) — komunikaty poniżej obsługują okno konfiguracji otwierane ze strefy 2 strony głównej (rozdz. 6.2) oraz komponent własny wpięty w sesję dowolnego modułu.

| Typ komunikatu | Kierunek | Okno | Opis |
|---|---|---|---|
| `automation.workflow.save` | polecenie | Workflow Builder | Zapis definicji procesu automatycznego. |
| `automation.workflow.list` | polecenie | Workflow Builder | Wykaz zapisanych automatyk użytkownika. |
| `automation.schedule.set` | polecenie | Scheduler | Ustalenie cykliczności uruchomienia automatyki. |
| `automation.queue.action` | polecenie | Queue Manager | Akcja na kolejce automatyki — akcje wspólne z silnikiem kolejek, rozdz. 11.4. |
| `automation.orchestrator.define` | polecenie | Orchestrator | Definicja zależności między krokami automatyki — mechanizm wspólny z orkiestracją, rozdz. 11.5. |
| `automation.execution.subscribe` | polecenie | Execution Monitor | Subskrypcja statusu kolejnych uruchomień automatyki. |
| `automation.execution.status` | zdarzenie | Execution Monitor | Status przebiegu automatyki, w tym błędy wymagające interwencji. |

Powiązanie automatyki ze środowiskiem MultitaskingAI (koncepcja, rozdz. 11.3, 13.6) realizuje `automation.link` (rozdz. 11.6 niniejszego dokumentu) — połączenie nie jest domyślne i wymaga jawnego polecenia.

### 10.4. Browser

| Typ komunikatu | Kierunek | Okno | Opis |
|---|---|---|---|
| `browser.navigate` | polecenie | Browser Window | Przejście pod wskazany adres, widoczny jednocześnie użytkownikowi i AI. |
| `browser.snapshot.get` | polecenie | Browser Window | Odczyt bieżącego stanu wspólnego podglądu strony. |
| `browser.source.add` | polecenie | Sources Panel | Dodanie źródła do listy wykorzystanej w toku pracy. |
| `browser.note.add` | polecenie | Notes Panel | Odnotowanie istotnego fragmentu treści strony. |
| `browser.page.changed` | zdarzenie | Browser Window | Zmiana treści przeglądanej strony w wyniku nawigacji lub interakcji AI. |

### 10.5. Research

| Typ komunikatu | Kierunek | Okno | Opis |
|---|---|---|---|
| `research.source.add` | polecenie | Sources Manager | Dodanie źródła do badania. |
| `research.finding.add` | polecenie | Findings Panel | Odnotowanie ustalenia cząstkowego. |
| `research.report.build` | polecenie | Report Builder | Skompletowanie raportu końcowego z ustaleń i źródeł. |
| `research.report.export` | polecenie | Export Panel | Eksport raportu do formatu wymaganego przez odbiorcę. |
| `research.report.changed` | zdarzenie | Report Builder | Zmiana treści raportu w toku pracy. |

### 10.6. Library

| Typ komunikatu | Kierunek | Okno | Opis |
|---|---|---|---|
| `library.file.upload` | polecenie | Library Explorer | Przesłanie pliku do repozytorium. |
| `library.file.list` | polecenie | Library Explorer | Wykaz plików repozytorium, z filtrami. |
| `library.tag.set` | polecenie | Tags & Collections | Przypisanie tagów lub kolekcji do pliku. |
| `library.file.preview` | polecenie | File Preview | Podgląd zawartości pliku bez opuszczania modułu. |
| `library.version.list` | polecenie | Versioning Panel | Wykaz wersji pliku, w tym wersji wygenerowanych automatycznie przez AI w innych modułach. |
| `library.file.changed` | zdarzenie | Library Explorer | Dodanie, zmiana lub nowa wersja pliku w repozytorium. |

### 10.7. Translate

| Typ komunikatu | Kierunek | Okno | Opis |
|---|---|---|---|
| `translate.source.set` | polecenie | Source Panel | Ustalenie tekstu źródłowego. |
| `translate.target.add` | polecenie | Translation Panels | Dodanie języka docelowego do równoległego tłumaczenia. |
| `translate.glossary.set` | polecenie | Glossary Manager | Zapis preferowanego odpowiednika terminu w słowniku projektu. |
| `translate.translation.changed` | zdarzenie | Translation Panels | Zmiana treści dowolnego z równoległych tłumaczeń. |

### 10.8. Roundtable

| Typ komunikatu | Kierunek | Okno | Opis |
|---|---|---|---|
| `roundtable.model.add` | polecenie | Model Panels | Dodanie modelu do wspólnej dyskusji nad zagadnieniem. |
| `roundtable.debate.start` | polecenie | Debate Panel | Rozpoczęcie wymiany argumentów między modelami. |
| `roundtable.moderator.direct` | polecenie | Moderator Panel | Ukierunkowanie przebiegu dyskusji przez użytkownika. |
| `roundtable.consensus.get` | polecenie | Consensus Panel | Odczyt wypracowanego wspólnego stanowiska. |
| `roundtable.debate.changed` | zdarzenie | Debate Panel | Nowa wymiana argumentów w toku debaty. |

### 10.9. Design

| Typ komunikatu | Kierunek | Okno | Opis |
|---|---|---|---|
| `design.asset.generate` | polecenie | Prompt Builder | Wygenerowanie zasobu graficznego na podstawie polecenia. |
| `design.board.update` | polecenie | Design Board | Zestawienie zasobów w kompozycji roboczej. |
| `design.asset.list` | polecenie | Assets Panel | Wykaz wygenerowanych i zgromadzonych zasobów. |
| `design.asset.changed` | zdarzenie | Assets Panel | Nowy lub zmieniony zasób graficzny. |

### 10.10. Assistant

Moduł jest konfigurowany jako komponent własny (profil asystenta) na stronie głównej (rozdz. 6.2) i wykorzystywany operacyjnie w środowisku TalkIn.

| Typ komunikatu | Kierunek | Okno | Opis |
|---|---|---|---|
| `assistant.voice.command` | polecenie | Voice Console | Polecenie głosowe przekazane do realizacji. |
| `assistant.action.status` | polecenie | Actions Monitor | Odczyt statusu realizacji zlecenia wieloetapowego. |
| `assistant.activity.list` | polecenie | Activity Feed | Chronologiczny zapis wykonanych działań. |
| `assistant.action.changed` | zdarzenie | Actions Monitor | Zmiana statusu realizowanego działania. |

### 10.11. Terminal

| Typ komunikatu | Kierunek | Okno | Opis |
|---|---|---|---|
| `terminal.session.open` | polecenie | Terminal Tabs | Otwarcie nowej karty powłoki (PowerShell, CMD, Bash, Node.js, Python). |
| `terminal.command.exec` | polecenie | Terminal Tabs | Wykonanie polecenia w powłoce. |
| `terminal.process.list` | polecenie | Process Monitor | Wykaz uruchomionych procesów, w tym zainicjowanych przez AI. |
| `terminal.process.kill` | polecenie | Process Monitor | Zakończenie wskazanego procesu. |
| `terminal.output.stream` | strumień | Output Console | Wynik wykonania polecenia, przekazywany na żywo. |
| `terminal.process.changed` | zdarzenie | Process Monitor | Zmiana stanu procesu. |

### 10.12. Developer

| Typ komunikatu | Kierunek | Okno | Opis |
|---|---|---|---|
| `developer.file.open` | polecenie | Code Editor | Otwarcie pliku źródłowego. |
| `developer.file.save` | polecenie | Code Editor | Zapis pliku źródłowego. |
| `developer.tree.get` | polecenie | Project Tree | Odczyt struktury projektu. |
| `developer.git.action` | polecenie | Git Panel | Akcja kontroli wersji (status, commit, push, pull, branch). |
| `developer.build.run` | polecenie | Build Output | Uruchomienie kompilacji lub budowania. |
| `developer.build.changed` | zdarzenie | Build Output | Zmiana statusu lub wyniku budowania. |

### 10.13. Diagnostics

| Typ komunikatu | Kierunek | Okno | Opis |
|---|---|---|---|
| `diagnostics.log.query` | polecenie | Logs Viewer | Zapytanie o logi wg kryteriów. |
| `diagnostics.error.list` | polecenie | Errors Panel | Wykaz zarejestrowanych błędów. |
| `diagnostics.analyze.run` | polecenie | Diagnostics Center | Uruchomienie analizy agregującej logi i błędy w spójny obraz stanu systemu. |
| `diagnostics.recommendation.list` | polecenie | Recommendations Panel | Wykaz sugerowanych kroków naprawczych. |
| `diagnostics.analysis.changed` | zdarzenie | Diagnostics Center | Zakończenie lub aktualizacja analizy. |

### 10.14. Apps

| Typ komunikatu | Kierunek | Okno | Opis |
|---|---|---|---|
| `apps.architecture.define` | polecenie | Architecture Designer | Zaprojektowanie architektury rozwiązania. |
| `apps.workspace.update` | polecenie | Frontend Workspace / Backend Workspace | Aktualizacja pracy nad wskazaną warstwą produktu. |
| `apps.deployment.run` | polecenie | Deployment Panel | Uruchomienie procesu wdrożenia. |
| `apps.build.changed` | zdarzenie | Deployment Panel | Zmiana statusu budowy lub wdrożenia. |

Realizacja rozbudowanych projektów w module Apps jest naturalnym zastosowaniem środowiska MultitaskingAI (koncepcja, rozdz. 11.14) — kontrakt ról Executor 1 i Executor 2 przy równoległej pracy nad backendem i frontendem opisano w rozdziale 11.2.

### 10.15. Agents

| Typ komunikatu | Kierunek | Okno | Opis |
|---|---|---|---|
| `agent.create` | polecenie | Agent Builder | Utworzenie agenta — komponentu własnego (rozdz. 6.2). |
| `agent.update` | polecenie | Agent Builder | Zmiana konfiguracji agenta. |
| `agent.list` | polecenie | Agent Builder | Wykaz agentów użytkownika. |
| `agent.delete` | polecenie | Agent Builder | Usunięcie agenta. |
| `agent.model.set` | polecenie | Model Configuration | Wybór modelu bazowego agenta. |
| `agent.skill.add` | polecenie | Skills Manager | Dodanie umiejętności lub pluginu agentowi. |
| `agent.connector.add` | polecenie | Connectors Manager | Podłączenie integracji zewnętrznej. |
| `agent.permission.set` | polecenie | Permissions Center | Ustalenie zakresu uprawnień agenta. |
| `agent.changed` | zdarzenie | — | Zmiana konfiguracji agenta, w tym gdy pełni rolę w środowisku MultitaskingAI (rozdz. 11.2). |

---

## 11. MultitaskingAI — panel orkiestracji

Środowisko MultitaskingAI nie ma bocznej nawigacji modułów (koncepcja, rozdz. 9.4, 10) — `module.list` dla tego środowiska zwraca zamiast wykazu modułów strukturę panelu orkiestracji, złożonego z sześciu sekcji ustalonych w rozdziale 13.8 koncepcji. Poniższe podrozdziały 11.1–11.6 odpowiadają tym sekcjom w tej samej kolejności, a podrozdział 11.7 obejmuje konfigurację widoczności i kolejności samego panelu.

| Typ komunikatu | Kierunek | Opis |
|---|---|---|
| `panel.sections.get` | polecenie | Odczyt kolejności i widoczności sześciu sekcji panelu orkiestracji dla bieżącej karty sesji. |
| `panel.sections.set` | polecenie | Zmiana kolejności lub widoczności sekcji; zestaw domyślny — Zespoły, Role, Kolejki, Orkiestracja, Harmonogram i automatyki, Monitor procesu — obowiązuje przy braku odmiennego ustawienia (koncepcja, rozdz. 13.8). |

### 11.1. Zespoły

| Typ komunikatu | Kierunek | Opis |
|---|---|---|
| `team.save` | polecenie | Zapis zestawu ról, powiązań i kolejek jako nazwanej, zapisanej konfiguracji zespołu. |
| `team.list` | polecenie | Wykaz zapisanych zespołów. |
| `team.load` | polecenie | Wczytanie zapisanej konfiguracji zespołu do bieżącej karty sesji MultitaskingAI. |
| `team.duplicate` | polecenie | Duplikowanie zapisanej konfiguracji zespołu pod nową nazwą. |
| `team.changed` | zdarzenie | Zmiana zapisanej konfiguracji zespołu. |

### 11.2. Role

Realizuje cztery okna robocze ról ustalone w rozdziale 13.3 koncepcji: Executor 1, Coordinator, Executor 2, Executor 3 / Validator, wraz z mechanizmem Subagent Network.

| Typ komunikatu | Kierunek | Opis |
|---|---|---|
| `role.assign` | polecenie | Przypisanie modelu lub agenta z modułu Agents (rozdz. 10.15) do roli: `executor1`, `executor2`, `coordinator`, `executor3`. |
| `role.list` | polecenie | Wykaz ról obsadzonych w bieżącej konfiguracji zespołu. |
| `role.update` | polecenie | Zmiana konfiguracji roli, w tym trybu współpracy dla Executor 2 (praca niezależna, przekazywanie wyników, praca naprzemienna, praca iteracyjna) oraz wcielenia Executor 3 / Validator (Validator, Reviewer, Security Auditor, Architect, Product Owner, QA Lead, Arbitrator). |
| `role.remove` | polecenie | Usunięcie przypisania roli. |
| `subagent.spawn` | polecenie | Uruchomienie przez Executora podagentów Subagent Network — do 15 jednocześnie, z podaną specjalizacją każdego (na przykład Agent UI, Agent Backend, Agent API, Agent Security, Agent Database, Agent Testing). |
| `subagent.list` | polecenie | Wykaz aktywnych podagentów uruchomionych przez wskazanego Executora. |
| `subagent.result.collect` | polecenie | Agregacja wyników pracy podagentów przez Executora. |
| `role.changed` | zdarzenie | Zmiana obsady lub konfiguracji roli. |
| `subagent.changed` | zdarzenie | Zmiana stanu podagenta Subagent Network. |

### 11.3. Kolejki

Silnik kolejek jest mechanizmem wspólnym dla modułu Automations (rozdz. 10.3) i środowiska MultitaskingAI (koncepcja, rozdz. 13.4) — kontrakt poniżej obowiązuje w obu miejscach.

| Typ komunikatu | Kierunek | Opis |
|---|---|---|
| `queue.create` | polecenie | Utworzenie kolejki w zakresie globalnym, lokalnym, dla modelu, dla agenta lub dla projektu. |
| `queue.list` | polecenie | Wykaz kolejek wraz z zakresem i zawartością. |
| `queue.action` | polecenie | Akcja na kolejce: `enqueue`, `dequeue`, `delay`, `retry`, `pause`, `resume`, `split`, `merge`, `route`, `branch`, `condition`. |
| `queue.link` | polecenie | Powiązanie silnika kolejek MultitaskingAI z silnikiem kolejek modułu Automations — połączenie ustanawiane decyzją użytkownika (koncepcja, rozdz. 6.2, 13.6). |
| `queue.changed` | zdarzenie | Zmiana zawartości lub stanu kolejki. |

### 11.4. Orkiestracja

| Typ komunikatu | Kierunek | Opis |
|---|---|---|
| `orchestration.dependency.set` | polecenie | Definicja zależności między modelami, agentami, zadaniami, kolejkami, automatyzacjami lub projektami. |
| `orchestration.dependency.list` | polecenie | Wykaz zdefiniowanych zależności. |
| `orchestration.dependency.remove` | polecenie | Usunięcie zależności. |
| `orchestration.validate` | polecenie | Weryfikacja spójności zdefiniowanych zależności (na przykład wykrycie cyklu). |
| `orchestration.changed` | zdarzenie | Zmiana zależności lub wynik jej respektowania w toku wykonania procesu. |

### 11.5. Harmonogram i automatyki

Realizuje integrację z modułem Automations opisaną w rozdziale 13.6 koncepcji — połączenie nie jest domyślne i wymaga decyzji podjętej zarówno w oknie konfiguracji, jak i w ustawieniach okna modułu Automations.

| Typ komunikatu | Kierunek | Opis |
|---|---|---|
| `automation.link` | polecenie | Powiązanie automatyki (komponentu własnego, rozdz. 10.3) z harmonogramem pracy ciągłej środowiska MultitaskingAI. |
| `automation.link.remove` | polecenie | Zniesienie powiązania. |
| `schedule.get` | polecenie | Odczyt harmonogramu pracy ciągłej powiązanego z bieżącym zespołem. |
| `automation.link.changed` | zdarzenie | Zmiana powiązania automatyki z harmonogramem pracy ciągłej. |

### 11.6. Monitor procesu

Miejsce nadzoru funkcji globalnej Always On Display (koncepcja, rozdz. 12.2, 13.2) nad przebiegiem pętli.

| Typ komunikatu | Kierunek | Opis |
|---|---|---|
| `monitor.subscribe` | polecenie | Subskrypcja podglądu przebiegu pętli w zakresie procesu, roli albo kolejki. |
| `monitor.status` | zdarzenie | Status przebiegu pętli: aktywny krok, hierarchia decyzji, statusy przebiegów ról i kolejek. |
| `aod.observe.attach` | polecenie | Dołączenie Always On Display do procesu MultitaskingAI w trybie `observer` (obserwacja) albo `operator` (nadzór z możliwością interwencji). |
| `aod.observe.detach` | polecenie | Odłączenie Always On Display od procesu. |

### 11.7. Schemat przepływu kontraktu

Poniższy schemat przedstawia, które sekcje panelu orkiestracji odpowiadają którym komunikatom, w kolejności odpowiadającej przepływowi orkiestracji z rozdziału 13.9 koncepcji.

```
Coordinator                       role.assign(coordinator) · role.update
    │  orchestration.dependency.set
    ▼
Executor 1 / Executor 2           role.assign(executor1|executor2) · role.update
    │  subagent.spawn (≤15) · subagent.result.collect
    ▼
Executor 3 / Validator             role.assign(executor3) · role.update
    │
    ▼
queue.action (enqueue · dequeue · delay · retry · pause · resume
              · split · merge · route · branch · condition)
    │
    ▼
orchestration.validate  →  orchestration.changed
    │
    ▼
automation.link  →  schedule.get   (praca ciągła 24/7/365)
    │
    ▼
monitor.subscribe  →  monitor.status   (nadzór Always On Display: aod.observe.attach)
```

---

## 12. Pamięć

Realizuje mechanizm pamięci opisany w rozdziale 6 dokumentu Model danych, działający na poziomach: globalnym, środowiska, modułu, projektu i sesji, zgodnie z hierarchią zasięgów izolacji (rozdz. 15.1).

| Poziom pamięci | Zasięg widoczności |
|---|---|
| Globalny | Wszystkie środowiska, moduły i sesje konta. |
| Środowisko | Jedno środowisko: TalkIn, WorkSpace, CodeStudio albo MultitaskingAI. |
| Moduł | Jeden moduł w obrębie środowiska. |
| Projekt | Jeden projekt modułu Workspace. |
| Sesja | Pojedyncza karta sesji. |

Poziom wskazuje się w ładunku poleceń `memory.list`, `memory.set`, `memory.toggle` i `memory.delete`.

| Typ komunikatu | Kierunek | Opis |
|---|---|---|
| `memory.list` | polecenie | Wykaz zasobów pamięci widocznych na wskazanym poziomie. |
| `memory.set` | polecenie | Zapis zasobu pamięci na wskazanym poziomie. |
| `memory.toggle` | polecenie | Włączenie lub wyłączenie pamięci na wskazanym poziomie. |
| `memory.detach` | polecenie | Odłączenie dostępu do pamięci w bieżącej sesji — dostępne przed pierwszym promptem sesji. |
| `memory.delete` | polecenie | Usunięcie zasobu pamięci. |
| `memory.changed` | zdarzenie | Zmiana zasobu pamięci na dowolnym poziomie. |

---

## 13. Historia i retencja

| Typ komunikatu | Kierunek | Opis |
|---|---|---|
| `history.load` | polecenie | Odczyt historii sesji lub jej fragmentu. |
| `history.delete` | polecenie | Usunięcie wskazanego wpisu historii albo historii w całości — pełny dostęp do operacji na danych zgodnie z rozdziałem 10 Architektury. |
| `retention.set` | polecenie | Ustalenie okresu przechowywania historii na wskazanym poziomie; brak ustawienia oznacza przechowywanie bez limitu (Model danych, rozdz. 4.3). |
| `history.changed` | zdarzenie | Zmiana historii sesji, rozgłaszana do wszystkich urządzeń. |

---

## 14. Konfiguracja

Realizuje okno konfiguracji opisane w rozdziale 5.2 dokumentu Architektury oraz warstwową strukturę konfiguracji z rozdziału 10 tego dokumentu: warstwa domyślna oraz warstwa sesji, nakładająca się na warstwę domyślną bez jej modyfikacji.

| Typ komunikatu | Kierunek | Opis |
|---|---|---|
| `config.get` | polecenie | Odczyt wartości konfiguracji na wskazanym poziomie (globalny, środowisko, moduł, projekt, sesja) i dla wskazanego klucza. |
| `config.set` | polecenie | Zapis wartości konfiguracji na wskazanym poziomie. |
| `config.reset` | polecenie | Usunięcie nadpisania na wskazanym poziomie — wartość zaczyna być dziedziczona z poziomu szerszego. |
| `config.explain.get` | polecenie | Odczyt objaśnienia kontekstowego (`[?]`) opisującego działanie ustawienia i jego wpływ na aplikację; treść objaśnienia jest częścią definicji ustawienia (Architektura, rozdz. 13). |
| `model.channel.set` | polecenie | Konfiguracja kanału integracji modelu (API, CLI, SSH, HTTP) dla sesji lub roli (Architektura, rozdz. 9). |
| `config.changed` | zdarzenie | Zmiana wartości konfiguracji na dowolnym poziomie. |

---

## 15. Izolacja — okno konfiguracji punktów izolacji

Rozdział realizuje kontrakt komunikacyjny okna trzypanelowego ustalonego w rozdziale 6.4–6.6 koncepcji, łączącego izolację kontekstu (historia, pamięć, kontekst) i izolację techniczną procesu sesji (osiem zakresów z rozdziału 12 Architektury). Zgodnie z zasadą braku twardych blokad okno nigdy nie wymusza izolacji — udostępnia ją jako możliwość; żadne z poniższych poleceń nie ma odpowiednika wymuszającego wykonanie bez jawnego polecenia użytkownika.

### 15.1. Panel selektora zasięgu

| Typ komunikatu | Kierunek | Opis |
|---|---|---|
| `isolation.scope.list` | polecenie | Wykaz siedmiu poziomów zasięgu (koncepcja, rozdz. 6.5) wraz z kolejnością pierwszeństwa: globalny, środowisko, moduł, para modułów, projekt, karta sesji, rola. |

| Poziom zasięgu | Pole `scope` | Pierwszeństwo |
|---|---|---|
| Globalny | `global` | Najniższe |
| Środowisko | `environment` | ↑ |
| Moduł | `module` | ↑ |
| Para modułów | `modulePair` | ↑ |
| Projekt | `project` | ↑ |
| Karta sesji | `sessionTab` | ↑ |
| Rola (MultitaskingAI) | `role` | Najwyższe |

### 15.2. Panel macierzy izolacji

| Typ komunikatu | Kierunek | Opis |
|---|---|---|
| `isolation.context.get` | polecenie | Odczyt ustawień izolacji kontekstu (`history`, `memory`, `context` — każde `shared` albo `separate`) dla wskazanego zasięgu i identyfikatora. |
| `isolation.context.set` | polecenie | Zapis ustawień izolacji kontekstu. |
| `isolation.technical.get` | polecenie | Odczyt ustawień ośmiu zakresów izolacji technicznej dla wskazanego zasięgu. |
| `isolation.technical.set` | polecenie | Zapis ustawień izolacji technicznej. |

Ośmioma zakresami izolacji technicznej, każdy w stanie `on` albo `off`, są: `workingDirectory` (katalog roboczy sesji), `processEnvironment` (środowisko procesu), `modelDataConfig` (katalog danych i konfiguracji modelu), `networkAccess` (dostęp sieciowy), `fileReadWrite` (zakres odczytu i zapisu plików), `accountToken` (konto i token per sesja), `processModel` (model procesu — odrębny lub współdzielony), `executionServer` (serwer wykonania).

### 15.3. Panel profilu i podglądu

| Typ komunikatu | Kierunek | Opis |
|---|---|---|
| `isolation.profile.save` | polecenie | Zapis bieżącego zestawu ustawień jako nazwanego profilu izolacji dla wskazanego zasięgu. |
| `isolation.profile.list` | polecenie | Wykaz zapisanych profili izolacji. |
| `isolation.profile.load` | polecenie | Wczytanie profilu do panelu macierzy izolacji. |
| `isolation.profile.assign` | polecenie | Przypisanie profilu do sesji albo do roli. |
| `isolation.profile.delete` | polecenie | Usunięcie profilu. |
| `isolation.layer.set` | polecenie | Wybór warstwy konfiguracji: `default` (obowiązuje przy każdej nowej sesji) albo `session` (zmiany bieżącej sesji, nakładające się na warstwę domyślną bez jej modyfikacji). |
| `isolation.policy.preview` | polecenie | Odczyt polityki efektywnej dla wskazanego zasięgu — wynikowego zestawu reguł po uwzględnieniu dziedziczenia zgodnie z zasadą pierwszeństwa zasięgu najbardziej szczegółowego (koncepcja, rozdz. 6.5): brak ustawienia na danym poziomie dziedziczy wartość z poziomu bezpośrednio szerszego, aż do poziomu globalnego. |
| `isolation.profile.changed` | zdarzenie | Zmiana profilu izolacji. |
| `isolation.policy.changed` | zdarzenie | Zmiana polityki efektywnej dla wskazanego zasięgu, w wyniku zmiany ustawienia na dowolnym poziomie dziedziczenia. |

Stanem wyjściowym, zwracanym przez `isolation.policy.preview` przy braku jakiegokolwiek profilu lub zmienionego ustawienia, jest stan opisany w rozdziale 6.6 koncepcji: `isolation.context` dla zasięgu `sessionTab` równe `separate` dla wszystkich trzech pozycji (odrębna historia, pamięć i kontekst dla każdej nowej karty sesji), przy jednoczesnym `isolation.technical` równym `off` dla wszystkich ośmiu zakresów na poziomie globalnym.

```json
// Klient → Serwer — podgląd polityki efektywnej dla pary modułów Studio–Translate

{ "type": "isolation.policy.preview", "id": "c-40", "payload": {
    "scope": "modulePair", "scopeId": "studio+translate" },
  "timestamp": "2026-08-06T09:20:00Z" }

// Serwer → Klient

{ "type": "response", "id": "c-40", "status": "ok", "payload": {
    "context": { "history": "separate", "memory": "separate", "context": "shared" },
    "technical": { "workingDirectory": "off", "processEnvironment": "off",
      "modelDataConfig": "off", "networkAccess": "off", "fileReadWrite": "off",
      "accountToken": "off", "processModel": "off", "executionServer": "off" },
    "resolvedFrom": [
      { "field": "context.context", "scope": "modulePair", "scopeId": "studio+translate" },
      { "field": "context.history", "scope": "global", "scopeId": null },
      { "field": "context.memory", "scope": "global", "scopeId": null } ] },
  "timestamp": "2026-08-06T09:20:00Z" }
```

Pole `resolvedFrom` w odpowiedzi wskazuje, z którego poziomu zasięgu pochodzi każda wartość polityki efektywnej — mechanizm ten czyni dziedziczenie jawnym i weryfikowalnym z poziomu klienta, zgodnie z zasadą jawności i konfigurowalności zależności (koncepcja, rozdz. 14, zasada 3).

---

## 16. Funkcje globalne

Funkcje globalne nie tworzą własnej przestrzeni roboczej ani nowych sesji (koncepcja, rozdz. 5.3) — ich kontrakt nie obejmuje `workspace.enter` ani mechanizmu kart sesji.

### 16.1. Mobile

Pełne opracowanie funkcji wraz z protokołem parowania urządzenia i jego komunikatami (`device.pair.*`, `device.credential.refresh`, `device.revoke`, `device.logout`, `device.changed`) zawiera dokument `funkcje-globalne/mobile.md`, rozdz. 9.10.

| Typ komunikatu | Kierunek | Opis |
|---|---|---|
| `mobile.status.get` | polecenie | Stan dostępu mobilnego dla bieżącego urządzenia. |
| `mobile.process.list` | polecenie | Podgląd projektów, automatyzacji, agentów, aplikacji i sesji — monitoring procesów niezależny od aktywnego środowiska. |
| `mobile.process.control` | polecenie | Zarządzanie zadaniem: `start`, `stop`, `approve`, `pause`, `resume`, `modify`. |
| `mobile.process.changed` | zdarzenie | Zmiana stanu monitorowanego procesu, w szczególności procesu długotrwałego uruchomionego przez środowisko MultitaskingAI. |
| `mobile.queue.sync` | polecenie | Przekazanie kolejki działań wykonanych bez połączenia wraz z wersjami stanu do uzgodnienia (`funkcje-globalne/mobile.md`, rozdz. 7.3). |
| `mobile.notification.action` | polecenie | Działanie wykonane z poziomu powiadomienia wypychanego (`funkcje-globalne/mobile.md`, rozdz. 8.3). |

### 16.2. Always On Display

Kontrakt realizuje funkcję globalną opisaną w `funkcje-globalne/always-on-display.md` (rozdz. 9.5).

| Typ komunikatu | Kierunek | Opis |
|---|---|---|
| `aod.status.get` | polecenie | Stan aktywacji i widoczności pływającego avatara. |
| `aod.chat.send` | polecenie | Wiadomość tekstowa skierowana do Always On Display, poza kontekstem konkretnej karty sesji. |
| `aod.voice.command` | polecenie | Polecenie głosowe skierowane do Always On Display. |
| `aod.context.get` | polecenie | Odczyt kontekstu, jakim dysponuje Always On Display: aktualny moduł, środowisko, projekt i historia użytkownika. |
| `aod.suggestion` | zdarzenie | Proaktywna sugestia: propozycja działania, konfiguracji, wskazanie problemu lub rekomendacja kolejnego kroku. |

Dołączenie Always On Display jako obserwatora lub operatora procesu środowiska MultitaskingAI realizuje `aod.observe.attach` (rozdz. 11.6).

---

## 17. Rozszerzenia

Realizuje mechanizm rozszerzeń opisany w rozdziale 12 dokumentu Architektury: wtyczki, umiejętności, konektory i serwery MCP, dzielące się na dwa źródła — Danaco Plugin (zestaw wbudowany) i Personal (instalowane przez Operatora).

| Typ komunikatu | Kierunek | Opis |
|---|---|---|
| `extension.list` | polecenie | Wykaz rozszerzeń, z filtrem po źródle (`danacoPlugin`, `personal`). |
| `extension.install` | polecenie | Instalacja rozszerzenia Personal przesłanego z urządzenia; trafia do katalogu użytkownika na serwerze. |
| `extension.configure` | polecenie | Zmiana konfiguracji rozszerzenia. |
| `extension.toggle` | polecenie | Włączenie lub wyłączenie rozszerzenia. |
| `extension.uninstall` | polecenie | Usunięcie rozszerzenia Personal. |
| `extension.changed` | zdarzenie | Zmiana stanu lub konfiguracji rozszerzenia. |

Rdzeń nie rozróżnia źródła rozszerzenia w toku działania — oba źródła realizują wspólny kontrakt integracji (Architektura, rozdz. 14); pole `source` w `extension.list` ma wyłącznie charakter informacyjny.

---

## 18. Synchronizacja wielourządzeniowa

Zgodnie z rozdziałem 16 Architektury synchronizacja odbywa się na żywo tym samym kanałem WebSocket, którym prowadzona jest cała pozostała komunikacja — nie jest odrębnym protokołem.

| Cecha synchronizacji | Zasada |
|---|---|
| Nośnik zmiany | Zmiana wykonana na jednym urządzeniu jest przekazywana przez serwer do pozostałych urządzeń powiązanych z tym samym kontem w postaci zdarzenia właściwego zmienionemu obszarowi (rozdz. 4.2 — na przykład `session.changed`, `memory.changed`, `config.changed`). |
| Zakres synchronizacji | Sesje, rozmowy, historia, pamięć, konfiguracja, stan procesów. |
| Inicjatywa | Serwer inicjuje przekazanie zmiany samodzielnie, bez zapytania klienta — urządzenia prezentują wyłącznie obraz pochodzący z serwera i nie utrzymują własnego, rozbieżnego stanu. |
| Brak scalania po stronie klienta | Ponieważ serwer jest jedynym źródłem prawdy (rozdz. 1, zasada 3), nie występuje scalanie zmian po stronie klienta — każda zmiana zatwierdzona przez serwer zastępuje poprzedni stan widoczny na wszystkich urządzeniach. |

Poniższy diagram przedstawia rozgłoszenie zmiany: polecenie z jednego urządzenia skutkuje zdarzeniem rozesłanym przez serwer do wszystkich pozostałych urządzeń tego samego konta.

```
Urządzenie A                     Serwer                Urządzenie B      Urządzenie C
     │                             │                        │                │
     │  polecenie                  │                        │                │
     │ ───────────────────────────►│                        │                │
     │  response                   │                        │                │
     │ ◄───────────────────────────│                        │                │
     │                             │  session.changed /      │                │
     │                             │  memory.changed /       │                │
     │                             │  config.changed         │                │
     │                             │ ──────────────────────►│                │
     │                             │ ───────────────────────────────────────►│
     │                             │                        │                │
 (stan źródłowy)          (jedyne źródło prawdy)      (odświeżenie)    (odświeżenie)
```

---

## 19. Obsługa błędów

Odpowiedź na polecenie zawiera wynik pozytywny (`status: "ok"`) albo obiekt błędu (`status: "error"`, pole `error`). Błąd operacji nie przerywa połączenia WebSocket ani nie blokuje wysyłania kolejnych poleceń.

| Pole obiektu `error` | Opis |
|---|---|
| `code` | Kod błędu z tabeli w Załączniku D. |
| `message` | Pełny, czytelny opis przyczyny błędu. |
| `details` | Obiekt z dodatkowymi danymi diagnostycznymi, obecny wtedy, gdy serwer takie dane przekazuje (na przykład wskazanie niepoprawnego pola). |
| `retryable` | `true`, jeżeli ponowienie tego samego polecenia bez zmian ma sens (na przykład chwilowa niedostępność kanału modelu). |

```json
{ "type": "response", "id": "c-77", "status": "error",
  "error": { "code": "validation_failed", "message": "Pole 'moduleId' jest wymagane.",
    "details": { "field": "moduleId" }, "retryable": false },
  "timestamp": "2026-08-06T09:30:00Z" }
```

---

## 20. Wersjonowanie kontraktu

Kontrakt komunikacji jest wersjonowany niezależnie od wersji koncepcji platformy i od wersji serwera. Numer wersji ma postać `MAJOR.MINOR.PATCH`. Klient i serwer uzgadniają wersję protokołu przy nawiązaniu połączenia, w komunikacie `connection.hello` / `connection.hello.ack` (rozdz. 2).

| Segment | Rodzaj zmiany | Zgodność wsteczna | Przykłady |
|---|---|---|---|
| `MAJOR` | Zmiana niezgodna wstecznie | Nie | Usunięcie typu komunikatu, zmiana znaczenia istniejącego pola, zmiana wymaganych pól. |
| `MINOR` | Rozszerzenie zgodne wstecznie | Tak | Dodanie nowego typu komunikatu, nowego pola o obecności warunkowej lub nowej wartości enumeracyjnej. |
| `PATCH` | Poprawka nienaruszająca kontraktu | Tak | Korekty niezmieniające zestawu ani znaczenia komunikatów i pól. |

Serwer wspiera bieżącą wersję `MAJOR` oraz jedną poprzednią przez okres przejściowy, sygnalizując przestarzałe typy komunikatów polem `deprecated: true` w odpowiedzi, bez przerywania ich obsługi w tym okresie.

---

## Załącznik A. Zbiorcza tabela poleceń

Zestawienie porządkuje polecenia (klient → serwer) opisane w rozdziałach 2–17 według obszaru funkcjonalnego.

| Obszar | Polecenia |
|---|---|
| Połączenie (rozdz. 2) | `connection.hello` |
| Uwierzytelnianie (rozdz. 5) | `auth.register`, `auth.confirm`, `auth.login`, `auth.method.add`, `auth.method.remove`, `auth.password.reset`, `auth.token.refresh` |
| Strona główna (rozdz. 6) | `home.enter`, `environment.list`, `environment.enter`, `component.list`, `component.create`, `component.update`, `component.delete`, `component.assign`, `config.window.open`, `mobile.status.get`, `aod.status.get` |
| Nawigacja (rozdz. 7.1) | `module.list`, `workspace.enter`, `window.state.get`, `window.action` |
| Warstwy widoczności (rozdz. 7.2) | `layer.element.get`, `layer.element.set`, `layer.profile.get`, `layer.search.query`, `layer.function.invoke` |
| Karty sesji (rozdz. 8) | `session.create`, `session.list`, `session.open`, `session.focus`, `session.update`, `session.bind`, `session.close`, `session.delete` |
| Chat Window — kanał Użytkownik ↔ Wykonawca (rozdz. 9.1) | `message.send`, `message.approve`, `message.stop`, `message.explain`, `message.list`, `chat.context.set`, `chat.context.get` |
| Execution Loop Window — kanał Koordynator ↔ Wykonawca (rozdz. 9.2) | `loop.subscribe`, `loop.state.get`, `loop.order.submit`, `loop.control` |
| Studio (rozdz. 10.1) | `studio.document.open`, `studio.document.save`, `studio.contextual.op`, `studio.diff.compare`, `studio.repository.list`, `studio.repository.restore` |
| Workspace (rozdz. 10.2) | `workspace.dashboard.get`, `workspace.instructions.set`, `workspace.context.set`, `workspace.context.get`, `workspace.library.list`, `workspace.agent.assign` |
| Automations (rozdz. 10.3) | `automation.workflow.save`, `automation.workflow.list`, `automation.schedule.set`, `automation.queue.action`, `automation.orchestrator.define`, `automation.execution.subscribe` |
| Browser (rozdz. 10.4) | `browser.navigate`, `browser.snapshot.get`, `browser.source.add`, `browser.note.add` |
| Research (rozdz. 10.5) | `research.source.add`, `research.finding.add`, `research.report.build`, `research.report.export` |
| Library (rozdz. 10.6) | `library.file.upload`, `library.file.list`, `library.tag.set`, `library.file.preview`, `library.version.list` |
| Translate (rozdz. 10.7) | `translate.source.set`, `translate.target.add`, `translate.glossary.set` |
| Roundtable (rozdz. 10.8) | `roundtable.model.add`, `roundtable.debate.start`, `roundtable.moderator.direct`, `roundtable.consensus.get` |
| Design (rozdz. 10.9) | `design.asset.generate`, `design.board.update`, `design.asset.list` |
| Assistant (rozdz. 10.10) | `assistant.voice.command`, `assistant.action.status`, `assistant.activity.list` |
| Terminal (rozdz. 10.11) | `terminal.session.open`, `terminal.command.exec`, `terminal.process.list`, `terminal.process.kill` |
| Developer (rozdz. 10.12) | `developer.file.open`, `developer.file.save`, `developer.tree.get`, `developer.git.action`, `developer.build.run` |
| Diagnostics (rozdz. 10.13) | `diagnostics.log.query`, `diagnostics.error.list`, `diagnostics.analyze.run`, `diagnostics.recommendation.list` |
| Apps (rozdz. 10.14) | `apps.architecture.define`, `apps.workspace.update`, `apps.deployment.run` |
| Agents (rozdz. 10.15) | `agent.create`, `agent.update`, `agent.list`, `agent.delete`, `agent.model.set`, `agent.skill.add`, `agent.connector.add`, `agent.permission.set` |
| MultitaskingAI (rozdz. 11) | `panel.sections.get`, `panel.sections.set`, `team.save`, `team.list`, `team.load`, `team.duplicate`, `role.assign`, `role.list`, `role.update`, `role.remove`, `subagent.spawn`, `subagent.list`, `subagent.result.collect`, `queue.create`, `queue.list`, `queue.action`, `queue.link`, `orchestration.dependency.set`, `orchestration.dependency.list`, `orchestration.dependency.remove`, `orchestration.validate`, `automation.link`, `automation.link.remove`, `schedule.get`, `monitor.subscribe`, `aod.observe.attach`, `aod.observe.detach` |
| Pamięć (rozdz. 12) | `memory.list`, `memory.set`, `memory.toggle`, `memory.detach`, `memory.delete` |
| Historia (rozdz. 13) | `history.load`, `history.delete`, `retention.set` |
| Konfiguracja (rozdz. 14) | `config.get`, `config.set`, `config.reset`, `config.explain.get`, `model.channel.set` |
| Izolacja (rozdz. 15) | `isolation.scope.list`, `isolation.context.get`, `isolation.context.set`, `isolation.technical.get`, `isolation.technical.set`, `isolation.profile.save`, `isolation.profile.list`, `isolation.profile.load`, `isolation.profile.assign`, `isolation.profile.delete`, `isolation.layer.set`, `isolation.policy.preview` |
| Mobile (rozdz. 16.1) | `mobile.status.get`, `mobile.process.list`, `mobile.process.control`, `mobile.queue.sync`, `mobile.notification.action` |
| Always On Display (rozdz. 16.2) | `aod.status.get`, `aod.chat.send`, `aod.voice.command`, `aod.context.get` |
| Rozszerzenia (rozdz. 17) | `extension.list`, `extension.install`, `extension.configure`, `extension.toggle`, `extension.uninstall` |

---

## Załącznik B. Zbiorcza tabela zdarzeń

| Obszar | Zdarzenia |
|---|---|
| Uwierzytelnianie | `auth.changed` |
| Nawigacja | `window.state.changed` |
| Warstwy widoczności | `layer.profile.changed` |
| Karty sesji | `session.changed`, `session.focus.changed` |
| Chat Window (kanał Użytkownik ↔ Wykonawca) | `message.approval.request`, `message.changed`, `chat.context.changed` |
| Execution Loop Window (kanał Koordynator ↔ Wykonawca) | `loop.order.accept`, `loop.order.decompose`, `loop.task.assign`, `loop.task.progress`, `loop.task.result`, `loop.quality.result`, `loop.task.retry`, `loop.order.close`, `loop.changed` |
| Komponenty własne | `component.changed` |
| Studio | `studio.document.changed` |
| Workspace | `workspace.project.changed` |
| Automations | `automation.execution.status` |
| Browser | `browser.page.changed` |
| Research | `research.report.changed` |
| Library | `library.file.changed` |
| Translate | `translate.translation.changed` |
| Roundtable | `roundtable.debate.changed` |
| Design | `design.asset.changed` |
| Assistant | `assistant.action.changed` |
| Terminal | `terminal.process.changed` |
| Developer | `developer.build.changed` |
| Diagnostics | `diagnostics.analysis.changed` |
| Apps | `apps.build.changed` |
| Agents | `agent.changed` |
| MultitaskingAI | `team.changed`, `role.changed`, `subagent.changed`, `queue.changed`, `orchestration.changed`, `automation.link.changed`, `monitor.status` |
| Pamięć | `memory.changed` |
| Historia | `history.changed` |
| Konfiguracja | `config.changed` |
| Izolacja | `isolation.profile.changed`, `isolation.policy.changed` |
| Mobile | `mobile.process.changed` |
| Always On Display | `aod.suggestion` |
| Rozszerzenia | `extension.changed` |

---

## Załącznik C. Przykładowe wymiany komunikatów

### C.1. Od strony głównej do sesji w module

```json
// 1. Klient → Serwer — strona główna

{ "type": "home.enter", "id": "c-1", "payload": {}, "timestamp": "2026-08-06T09:00:05Z" }

// 2. Klient → Serwer — wybór karty środowiska w strefie 1

{ "type": "environment.enter", "id": "c-2", "payload": { "environmentId": "codestudio" },
  "timestamp": "2026-08-06T09:00:10Z" }

// 3. Klient → Serwer — wybór modułu z bocznej nawigacji

{ "type": "workspace.enter", "id": "c-3", "payload": {
    "environmentId": "codestudio", "moduleId": "developer" },
  "timestamp": "2026-08-06T09:00:15Z" }

// 4. Serwer → Klient — nowa karta sesji z układem okien

{ "type": "response", "id": "c-3", "status": "ok", "payload": {
    "sessionId": "s-501", "environmentId": "codestudio", "moduleId": "developer",
    "windows": ["chat", "codeEditor", "projectTree", "gitPanel", "buildOutput"] },
  "timestamp": "2026-08-06T09:00:15Z" }

// 5. Serwer → Klient — rozgłoszenie do pozostałych urządzeń konta

{ "type": "session.changed", "id": "e-9001", "payload": {
    "sessionId": "s-501", "action": "created" },
  "timestamp": "2026-08-06T09:00:15Z" }
```

### C.2. Utworzenie komponentu własnego i jego wpięcie w sesję

```json
// Klient → Serwer — strefa 2, utworzenie automatyki

{ "type": "component.create", "id": "c-20", "payload": {
    "kind": "automation", "name": "Cotygodniowy raport sprzedaży",
    "workflow": { "trigger": "schedule", "cron": "0 7 * * MON" } },
  "timestamp": "2026-08-06T09:10:00Z" }

// Serwer → Klient

{ "type": "response", "id": "c-20", "status": "ok",
  "payload": { "componentId": "cmp-77" }, "timestamp": "2026-08-06T09:10:00Z" }

// Klient → Serwer — wpięcie automatyki w sesji modułu Developer

{ "type": "component.assign", "id": "c-21", "payload": {
    "componentId": "cmp-77", "sessionId": "s-501" },
  "timestamp": "2026-08-06T09:10:05Z" }
```

### C.3. Przypisanie roli w MultitaskingAI i uruchomienie Subagent Network

```json
// Klient → Serwer

{ "type": "role.assign", "id": "c-30", "payload": {
    "sessionId": "s-620", "role": "executor1", "agentId": "agt-14" },
  "timestamp": "2026-08-06T09:15:00Z" }

// Klient → Serwer — uruchomienie podagentów

{ "type": "subagent.spawn", "id": "c-31", "payload": {
    "sessionId": "s-620", "parentRole": "executor1",
    "specializations": ["Agent UI", "Agent Backend", "Agent API"] },
  "timestamp": "2026-08-06T09:15:10Z" }

// Serwer → Klient

{ "type": "response", "id": "c-31", "status": "ok", "payload": {
    "subagentIds": ["sub-1", "sub-2", "sub-3"] },
  "timestamp": "2026-08-06T09:15:10Z" }
```

### C.4. Od polecenia w Chat Window do zamknięcia zlecenia w Execution Loop Window

```json
// 1. Klient → Serwer — polecenie Użytkownika w kanale pierwszym

{ "type": "message.send", "id": "c-70", "payload": {
    "sessionId": "s-620", "windowId": "chat",
    "content": "Przygotuj moduł raportowania wraz z testami." },
  "timestamp": "2026-08-06T09:15:00Z" }

// 2. Klient → Serwer — subskrypcja okna pętli wykonawczej

{ "type": "loop.subscribe", "id": "c-71", "payload": { "sessionId": "s-620" },
  "timestamp": "2026-08-06T09:15:01Z" }

// 3. Serwer → Klient — Koordynator przyjmuje zlecenie

{ "type": "loop.order.accept", "id": "e-7001", "payload": {
    "orderId": "ord-12", "sessionId": "s-620", "source": "chat",
    "objective": "Przygotuj moduł raportowania wraz z testami." },
  "timestamp": "2026-08-06T09:15:02Z" }

// 4. Klient → Serwer — wstrzymanie i wznowienie przebiegu

{ "type": "loop.control", "id": "c-72", "payload": {
    "orderId": "ord-12", "action": "pause" },
  "timestamp": "2026-08-06T09:30:00Z" }
{ "type": "loop.control", "id": "c-73", "payload": {
    "orderId": "ord-12", "action": "resume" },
  "timestamp": "2026-08-06T09:34:00Z" }

// 5. Serwer → Klient — zamknięcie zlecenia

{ "type": "loop.order.close", "id": "e-7008", "payload": {
    "orderId": "ord-12", "outcome": "completed",
    "deliverables": ["schema.sql", "tests/"], "duration": "PT42M" },
  "timestamp": "2026-08-06T09:52:00Z" }
```

### C.5. Błąd walidacji przy tworzeniu sesji

```json
// Klient → Serwer

{ "type": "session.create", "id": "c-40", "payload": { "environmentId": "talkin" },
  "timestamp": "2026-08-06T09:18:00Z" }

// Serwer → Klient

{ "type": "response", "id": "c-40", "status": "error",
  "error": { "code": "validation_failed",
    "message": "Pole 'moduleId' jest wymagane do utworzenia karty sesji.",
    "details": { "field": "moduleId" }, "retryable": false },
  "timestamp": "2026-08-06T09:18:00Z" }
```

---

## Załącznik D. Kody błędów

| Kod | Znaczenie | `retryable` |
|---|---|---|
| `validation_failed` | Ładunek polecenia nie spełnia wymagań pola lub typu. | Nie |
| `not_found` | Wskazany zasób (sesja, moduł, komponent, profil izolacji) nie istnieje. | Nie |
| `not_authenticated` | Połączenie nie zostało uwierzytelnione, a faza wdrożenia tego wymaga. | Nie |
| `permission_denied` | Operacja przekracza uprawnienia przypisane kontu, agentowi lub roli. | Nie |
| `conflict` | Operacja koliduje z bieżącym stanem zasobu (na przykład usunięcie sesji już usuniętej). | Nie |
| `channel_unavailable` | Kanał integracji modelu (API, CLI, SSH, HTTP) jest chwilowo niedostępny. | Tak |
| `command_not_understood` | Treść polecenia języka naturalnego nie została dopasowana do żadnej funkcji (rozdz. 7.2). | Nie |
| `rate_limited` | Przekroczono dopuszczalną częstość poleceń danego typu. | Tak |
| `internal_error` | Błąd wewnętrzny serwera, niezwiązany z treścią polecenia. | Tak |

---

## Załącznik E. Słowniczek pojęć kontraktu

**Koperta** — wspólna struktura JSON każdego komunikatu wymienianego kanałem WebSocket, obejmująca pola `type`, `id`, `payload`, `timestamp` oraz pola specyficzne dla kategorii komunikatu (rozdz. 3).

**Polecenie** — komunikat klient → serwer, żądający wykonania operacji i zawsze korelowany z dokładnie jedną odpowiedzią serwera (rozdz. 4.1).

**Zdarzenie** — komunikat serwer → klient, inicjowany samodzielnie przez serwer i rozgłaszany do wszystkich urządzeń powiązanych z kontem, informujący o zmianie stanu (rozdz. 4.2).

**Strumień** — ciąg fragmentów przekazywanych na żywo w ramach jednego `id`, otwierany odpowiedzią na polecenie inicjujące i zamykany fragmentem z `done: true` (rozdz. 4.3).

**Korelacja** — mechanizm wiązania odpowiedzi i fragmentów strumienia z poleceniem, które je zainicjowało, przez zgodność pola `id`.

**Polityka efektywna** — wynikowy zestaw reguł izolacji dla wskazanego zasięgu, otrzymany po zastosowaniu zasady pierwszeństwa zasięgu najbardziej szczegółowego i dziedziczeniu wartości nieustawionych z poziomów szerszych (rozdz. 15.3, koncepcja rozdz. 6.5).

**Zasięg** — poziom, na którym obowiązuje reguła izolacji lub wartość konfiguracji: globalny, środowisko, moduł, para modułów, projekt, karta sesji albo rola (rozdz. 15.1).

**Warstwa domyślna / warstwa sesji** — dwa poziomy struktury warstwowej konfiguracji: warstwa domyślna obowiązuje przy każdej nowej sesji, warstwa sesji obejmuje zmiany dokonane dla sesji bieżącej i nakłada się na warstwę domyślną bez jej modyfikacji (rozdz. 14, Architektura rozdz. 13).

**Chat Window** — główne okno komunikacji kanału Użytkownik ↔ Wykonawca, zajmujące lewą kolumnę obszaru roboczego, stałą i o pełnej wysokości; podstawowy mechanizm sterowania procesami platformy (rozdz. 9.1).

**Execution Loop Window** — okno pętli wykonawczej kanału Koordynator ↔ Wykonawca, otwierane jako kolumna sąsiadująca z Chat Window; prezentuje zlecenie, dekompozycję na zadania, przebieg pętli i sterowanie nim (rozdz. 9.2).

**Koordynator** — komponent orkiestrujący platformy, przyjmujący zlecenie, dekomponujący je na zadania oraz przydzielający i nadzorujący ich realizację (rozdz. 9.2).

**Wykonawca** — AI, agent albo system wykonawczy realizujący polecenia Użytkownika i zadania przydzielone przez Koordynatora (rozdz. 9).

**Zlecenie** — jednostka pracy przyjęta przez Koordynatora, identyfikowana polem `orderId`, dekomponowana na zadania i zamykana komunikatem `loop.order.close` (rozdz. 9.2).

**Zadanie** — najmniejsza jednostka pracy pętli wykonawczej, identyfikowana polem `taskId`, przydzielana Wykonawcy i przechodząca zamknięty zbiór stanów od `assigned` do `closed` (rozdz. 9.2).

**Warstwa widoczności** — jedna z czterech warstw, do których należy każdy element interfejsu; wyznacza, czy element jest widoczny bez interakcji, czy ujawniany wyzwalaczem, oraz sposób jego wywołania (rozdz. 7.2).

**Faza wdrożenia** — wartość `build` albo `production` zwracana w `connection.hello.ack`, określająca, czy uwierzytelnienie jest wymagane do nawiązania połączenia (rozdz. 5, Architektura rozdz. 15).

---

## Załącznik F. Schematy ładunku wybranych komunikatów

Załącznik zestawia w formie tabel strukturę ładunku (`payload`) wybranych komunikatów, których przykłady w notacji JSON podano w rozdziałach 2–19 oraz w Załączniku C. Tabele obejmują pola ładunku właściwego; pola wspólnej koperty (`type`, `id`, `payload`, `timestamp`, a dla odpowiedzi i strumieni także `status`, `seq`, `done`) opisano w rozdziale 3 i nie są tu powtarzane. Schematy ładunku ujęte już przy rozdziałach szczegółowych — odpowiedzi na `home.enter` (rozdz. 6), fragmentu `stream.chunk` (rozdz. 9.1) oraz obiektu `error` (rozdz. 19) — nie są w załączniku powielane.

### F.1. Uzgodnienie połączenia — `connection.hello` i odpowiedź

Ładunek polecenia `connection.hello` (klient → serwer):

| Pole | Obecność | Opis |
|---|---|---|
| `protocolVersion` | zawsze | Wersja protokołu kontraktu w postaci `MAJOR.MINOR.PATCH` (rozdz. 20). |
| `deviceId` | zawsze | Identyfikator urządzenia zgłaszającego połączenie. |
| `token` | warunkowo | Token dostępu uzyskany po uwierzytelnieniu; wymagany w fazie produkcyjnej (rozdz. 5). |

Ładunek odpowiedzi na `connection.hello` (serwer → klient):

| Pole | Obecność | Opis |
|---|---|---|
| `protocolVersion` | zawsze | Uzgodniona wersja protokołu kontraktu. |
| `serverVersion` | zawsze | Wersja serwera wykonawczego. |
| `deploymentPhase` | zawsze | Faza wdrożenia: `build` albo `production` (rozdz. 5). |
| `authenticated` | zawsze | Wynik uwierzytelnienia połączenia. |

### F.2. Wejście do modułu — odpowiedź na `workspace.enter`

| Pole | Obecność | Opis |
|---|---|---|
| `sessionId` | zawsze | Identyfikator utworzonej karty sesji (rozdz. 8). |
| `environmentId` | zawsze | Środowisko, w którym utworzono kartę sesji. |
| `moduleId` | zawsze | Moduł otwarty w karcie sesji. |
| `windows` | zawsze | Zestaw okien operacyjnych właściwy modułowi (Załącznik A koncepcji). |

### F.3. Chat Window — `message.send` i odpowiedź

Ładunek polecenia `message.send` (klient → serwer):

| Pole | Obecność | Opis |
|---|---|---|
| `sessionId` | zawsze | Karta sesji, w której przesłano polecenie. |
| `windowId` | zawsze | Okno operacyjne będące kontekstem polecenia. |
| `content` | zawsze | Treść polecenia Użytkownika. |

Odpowiedź na `message.send` niesie w ładunku pole `streamId`, równe `id` polecenia, otwierające strumień odpowiedzi Wykonawcy; strukturę fragmentu `stream.chunk` opisano w rozdziale 9.1.

### F.4. Podgląd polityki izolacji — odpowiedź na `isolation.policy.preview`

| Pole | Obecność | Opis |
|---|---|---|
| `context` | zawsze | Izolacja kontekstu: `history`, `memory`, `context`, każde `shared` albo `separate` (rozdz. 15.2). |
| `technical` | zawsze | Osiem zakresów izolacji technicznej, każdy `on` albo `off` (rozdz. 15.2). |
| `resolvedFrom` | zawsze | Dla każdej wartości polityki efektywnej — poziom zasięgu i identyfikator, z którego wartość pochodzi (rozdz. 15.3). |

### F.5. Sterowanie przebiegiem pętli — `loop.control`

| Pole | Obecność | Opis |
|---|---|---|
| `orderId` | zawsze | Zlecenie, którego dotyczy sterowanie (rozdz. 9.2). |
| `action` | zawsze | Akcja sterująca: `pause`, `resume`, `abort` albo `amend`. |
| `objective` | warunkowo | Skorygowana treść zlecenia; wymagana przy akcji `amend`. |
| `tasks` | warunkowo | Skorygowany wykaz zadań przy akcji `amend`, gdy korekta obejmuje dekompozycję. |

### F.6. Wywołanie funkcji warstwy 4 — `layer.function.invoke`

| Pole | Obecność | Opis |
|---|---|---|
| `sessionId` | zawsze | Karta sesji, w kontekście której wywoływana jest funkcja. |
| `functionId` | warunkowo | Identyfikator funkcji wybranej ze skrótu klawiszowego albo z wyszukiwarki funkcji. |
| `naturalCommand` | warunkowo | Treść polecenia języka naturalnego, gdy funkcja wywoływana jest z okna komunikacji. |
| `invocation` | zawsze | Sposób wywołania: `naturalLanguage`, `shortcut` albo `search`. |

---

*Koniec dokumentu. Danaco Console — Kontrakty komunikacji, wersja 2.0.*

---
*Danaco Console — AI Workspace OS · v2.0*

*© 2026 Danaco Holding Group Sp. z o.o. — [LICENSE](LICENSE). Kontakt: support@danaco-group.pl*
