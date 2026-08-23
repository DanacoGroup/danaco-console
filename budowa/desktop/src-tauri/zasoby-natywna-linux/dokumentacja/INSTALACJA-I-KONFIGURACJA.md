# Danaco Console — Instalacja i konfiguracja

**Przewodnik Operatora: od pobrania pakietu klienckiego do skonfigurowanej platformy.**

| | |
|---|---|
| **Produkt** | Danaco Console |
| **Rodzaj** | Platforma AI Workspace OS |
| **Wersja** | v2.0 |
| **Adresat dokumentu** | Operator — właściciel konta pracujący na platformie |
| **Producent** | Danaco Holding Group Sp. z o.o., ul. Gen. J. Hallera 81, 43-400 Cieszyn |
| **Wsparcie** | support@danaco-group.pl |

---

## Spis treści

1. [Zakres dokumentu i model wdrożenia](#1-zakres-dokumentu-i-model-wdrożenia)
2. [Wymagania urządzenia](#2-wymagania-urządzenia)
3. [Instalacja pakietu klienckiego](#3-instalacja-pakietu-klienckiego)
4. [Pierwsze uruchomienie i rejestracja konta](#4-pierwsze-uruchomienie-i-rejestracja-konta)
5. [Logowanie i metody uwierzytelniania](#5-logowanie-i-metody-uwierzytelniania)
6. [Podłączanie kolejnych urządzeń](#6-podłączanie-kolejnych-urządzeń)
7. [Okno Ustawień — dziewięć sekcji](#7-okno-ustawień--dziewięć-sekcji)
8. [Okno konfiguracji — warstwy i zakresy](#8-okno-konfiguracji--warstwy-i-zakresy)
9. [Konfiguracja modeli: kanały i konta](#9-konfiguracja-modeli-kanały-i-konta)
10. [Konfiguracja obu kanałów komunikacji operacyjnej](#10-konfiguracja-obu-kanałów-komunikacji-operacyjnej)
11. [Punkty izolacji i profile izolacji](#11-punkty-izolacji-i-profile-izolacji)
12. [Rozszerzenia: wtyczki, umiejętności, konektory, serwery MCP](#12-rozszerzenia-wtyczki-umiejętności-konektory-serwery-mcp)
13. [Warstwy widoczności funkcji](#13-warstwy-widoczności-funkcji)
14. [Profile konfiguracji](#14-profile-konfiguracji)
15. [Powiadomienia](#15-powiadomienia)
16. [Kolejność zalecana przy uruchamianiu nowej instalacji](#16-kolejność-zalecana-przy-uruchamianiu-nowej-instalacji)
17. [Diagnostyka](#17-diagnostyka)
18. [Aktualizacja i deinstalacja](#18-aktualizacja-i-deinstalacja)

---

## 1. Zakres dokumentu i model wdrożenia

Dokument opisuje wszystko, co Operator wykonuje na **własnym urządzeniu**:
instalację pakietu klienckiego, pierwsze uruchomienie, założenie konta,
podłączenie kolejnych urządzeń oraz pełną konfigurację platformy z okna
Ustawień i okna konfiguracji.

**Danaco Console pracuje w modelu hybrydowym.** Cała logika, praca obliczeniowa,
dostęp do plików, procesy sesji i wywołania modeli wykonują się w aplikacji
serwerowej. Na urządzeniu Operatora działa **cienki klient** — lekki pakiet
uruchamiający natywne okno, które łączy się z serwerem i prezentuje obraz.

```
   Przetwarzanie, logika, stan   ─────►   SERWER PLATFORMY
   Dostęp, prezentacja, obraz    ─────►   CIENKI KLIENT NA URZĄDZENIU
```

Z tego modelu wynikają cztery fakty istotne dla instalacji:

- Operator **nie instaluje żadnych dodatkowych programów** — pakiet kliencki jest
  samowystarczalny. Narzędzia potrzebne modułom platformy stoją na serwerze.
- Stanowisko robocze **nie wymaga układu GPU** ani wysokiej wydajności.
- **Aktualizacja platformy odbywa się centralnie**, po stronie serwera; klient
  aktualizowany jest wyłącznie wtedy, gdy zmienia się samo okno aplikacji.
- **Stan trwały istnieje tylko na serwerze** — utrata urządzenia nie oznacza
  utraty pracy, a każde nowe urządzenie widzi ten sam, aktualny stan platformy.

Instalacja i utrzymanie aplikacji serwerowej należą do wdrożenia serwerowego
i wykraczają poza zakres tego dokumentu. Operator otrzymuje od dostawcy wdrożenia
jedną informację techniczną: **adres serwera platformy**.

---

## 2. Wymagania urządzenia

| Element | Wymaganie |
|---|---|
| Rodzaj urządzenia | Komputer, telefon albo tablet zdolny wyświetlić natywne okno aplikacji |
| System — komputer | Windows, Linux |
| System — urządzenie przenośne | Android, iOS, iPadOS |
| Układ GPU | Niewymagany |
| Pamięć i procesor | Wymagania typowe dla aplikacji okienkowej; obliczenia wykonuje serwer |
| Dodatkowe oprogramowanie | Żadne |
| Sieć | Stałe połączenie z serwerem platformy; kanał komunikacji pozostaje otwarty w czasie pracy |
| Uprawnienia | Prawo instalacji aplikacji na urządzeniu |

Silnik prezentacji osadzany jest przez system urządzenia: na Windows komponent
WebView2, na Androidzie systemowy komponent WebView, na iOS i iPadOS komponent
WKWebView. Klient nie wymaga od Operatora żadnych czynności w tym zakresie.

---

## 3. Instalacja pakietu klienckiego

### 3.1. Pakiety

| Platforma | Pakiet |
|---|---|
| Windows | `Danaco Console_<wersja>_x64-setup.exe` — instalator |
| Linux (dystrybucje z pakietami `deb`) | `Danaco Console_<wersja>_amd64.deb` |
| Linux (uniwersalnie) | `Danaco Console_<wersja>_amd64.AppImage` |
| Android, iOS, iPadOS | Aplikacja mobilna — dostęp funkcją globalną Mobile (rozdz. 6.2) |

### 3.2. Windows

1. Uruchomić instalator `Danaco Console_<wersja>_x64-setup.exe`.
2. Przejść kreator instalacji i wskazać katalog docelowy albo pozostawić domyślny.
3. Zakończyć instalację i uruchomić Danaco Console ze skrótu.

Instalacja nie wymaga uprawnień administratora poza samym zapisem plików
aplikacji i nie dodaje do systemu żadnych usług pracujących w tle.

### 3.3. Linux

Pakiet `deb` instaluje się menedżerem pakietów dystrybucji:

```bash
sudo dpkg -i "Danaco Console_<wersja>_amd64.deb"
```

Wariant `AppImage` nie wymaga instalacji — wystarczy nadać plikowi prawo
wykonywania i uruchomić:

```bash
chmod +x "Danaco Console_<wersja>_amd64.AppImage" && "./Danaco Console_<wersja>_amd64.AppImage"
```

### 3.4. Wskazanie adresu serwera

Klient łączy się z serwerem pod adresem skonfigurowanym w pakiecie klienckim.
Jeżeli pakiet dostarczono bez wpisanego adresu, okno startowe poprosi o jego
podanie przy pierwszym uruchomieniu. Adres zapisuje się na urządzeniu jednorazowo
i przy kolejnych uruchomieniach nie jest wymagany ponownie; można go zmienić
w oknie Ustawień.

---

## 4. Pierwsze uruchomienie i rejestracja konta

### 4.1. Przebieg pierwszego uruchomienia

```
Uruchomienie pakietu klienckiego
        ▼
OKNO STARTOWE — godło marki, wskaźnik „łączenie z serwerem"
   klient otwiera kanał komunikacji i uzgadnia połączenie z serwerem
        ▼
OKNO REJESTRACJI  (konto właściciela jeszcze nie istnieje)
   login · adres e-mail uwierzytelniający · hasło
        ▼
WERYFIKACJA ADRESU E-MAIL
   konto pozostaje niepotwierdzone do chwili potwierdzenia
        ▼
UTWORZENIE KONTA WŁAŚCICIELA + wydanie tokenu dostępu urządzeniu
        ▼
STRONA GŁÓWNA — Centrum dowodzenia
```

Okno startowe jest stanem przejściowym i nie wymaga żadnej decyzji. Jeżeli serwer
jest nieosiągalny, wskaźnik ładowania zostaje zastąpiony komunikatem o braku
połączenia i kontrolką „Spróbuj ponownie" — patrz rozdz. 17.

### 4.2. Dane wymagane przy rejestracji

| Pole | Rola |
|---|---|
| **Login** | Nazwa właściciela konta, używana przy logowaniu razem z hasłem |
| **Adres e-mail uwierzytelniający** | Weryfikowany przy rejestracji; pełni dodatkowo funkcję metody logowania oraz **jedynej drogi odzyskania konta** |
| **Hasło** | Hasło dostępowe; przechowywane w postaci skrótu, poza bazą danych |

Platforma jest systemem jednego konta właściciela. Rejestracja wykonywana jest
**raz** — przy pierwszym uruchomieniu platformy. Każde następne urządzenie dochodzi
do tego samego konta przez logowanie albo parowanie (rozdz. 6).

### 4.3. Stan metod uwierzytelniania po rejestracji

| Metoda | Stan po rejestracji |
|---|---|
| Hasło | Aktywna — metoda bazowa |
| Adres e-mail uwierzytelniający | Aktywna |
| PIN | Nieaktywna — wymaga włączenia przez Operatora |
| Windows Hello | Nieaktywna — wymaga włączenia przez Operatora |

---

## 5. Logowanie i metody uwierzytelniania

### 5.1. Przebieg logowania

Logowanie wykonywane jest przy uruchomieniu platformy na urządzeniu, które nie ma
ważnego tokenu dostępu, oraz przy podłączaniu kolejnego urządzenia.

```
Uruchomienie na urządzeniu bez ważnego tokenu
        ▼
OKNO LOGOWANIA — wybór metody:
   • login + hasło                       (metoda bazowa, zawsze dostępna)
   • PIN · Windows Hello · e-mail        (aktywne metody dodatkowe)
        ▼
Weryfikacja po stronie serwera
        ├── pozytywna ──► wydanie tokenu dostępu → połączenie → strona główna
        └── negatywna ──► ponowna próba albo odzyskiwanie konta
```

Token dostępu jest przypisany do konkretnego urządzenia. Urządzenie z ważnym
tokenem wchodzi przy kolejnych uruchomieniach wprost na stronę główną albo do
ostatnio aktywnej karty sesji — bez okna logowania.

### 5.2. Metody dodatkowe

Metody dodatkowe **rozszerzają**, a nie zastępują logowania hasłem. Żadna nie jest
aktywna z konieczności technicznej — każdą włącza i wyłącza Operator w oknie
Ustawień, w sekcji Uwierzytelnianie.

| Metoda | Przeznaczenie |
|---|---|
| **PIN** | Krótki kod numeryczny do szybkiego potwierdzenia tożsamości na urządzeniu już uwierzytelnionym pełnym loginem i hasłem |
| **Windows Hello** | Uwierzytelnienie mechanizmem systemowym Windows, wiązane z konkretnym urządzeniem |
| **Adres e-mail uwierzytelniający** | Samodzielna metoda logowania, alternatywna wobec loginu i hasła |

Materiał każdej metody — skrót hasła, materiał PIN, materiał Windows Hello —
przechowywany jest poza bazą danych platformy.

### 5.3. Wymóg logowania

Pozycja **„Wymóg logowania"** w oknie konfiguracji, zakres Aplikacja, rozstrzyga,
czy okno rejestracji i logowania musi zostać skutecznie wypełnione, zanim Operator
przejdzie dalej. Ustawienie jest dostępne Operatorowi w każdej chwili;
w środowisku produkcyjnym jego wartością domyślną jest wymóg aktywny.

### 5.4. Odzyskiwanie konta

Jedyną drogą odzyskania dostępu jest adres e-mail uwierzytelniający podany przy
rejestracji. Procedurę uruchamia się z okna logowania. Z tego powodu adres ten
powinien pozostawać pod stałą kontrolą Operatora — jego utrata oznacza utratę
drogi odzyskania konta.

---

## 6. Podłączanie kolejnych urządzeń

Jedno konto obsługuje wiele urządzeń — komputery, telefony, tablety — bez
ograniczenia liczby. Każde urządzenie ma własny token dostępu; **odwołanie
dostępu jednemu urządzeniu nie odcina pozostałych**.

### 6.1. Kolejny komputer

Instalacja pakietu klienckiego (rozdz. 3), uruchomienie, logowanie loginem
i hasłem do istniejącego konta. Urządzenie zostaje wpisane do rejestru urządzeń
i otrzymuje własny token.

### 6.2. Urządzenie przenośne — parowanie

Urządzenie przenośne wchodzi do platformy przez parowanie, zakładane przy
maszynie już uwierzytelnionej, bez wpisywania loginu i hasła na telefonie.

```
Ustawienia → sekcja Urządzenia → [ Sparuj urządzenie ]
        ▼
Modal kodu parowania: kod znakowy + kod graficzny + licznik ważności
        ▼
Odczyt kodu aplikacją na urządzeniu przenośnym
   (kamerą albo wpisaniem kodu znakowego)
        ▼
Modal potwierdzenia zgłoszenia: typ urządzenia, nazwa, czas zgłoszenia
        ▼
Wydanie poświadczenia — urządzenie w rejestrze, tryb Mobile aktywny
```

Kod parowania ma ograniczoną ważność. Po jej upływie wydaje się nowy kod
działaniem „Wydaj nowy kod"; kod niewykorzystany można zamknąć działaniem
„Unieważnij kod".

### 6.3. Zarządzanie urządzeniami

Rejestr urządzeń w sekcji Urządzenia prezentuje przy każdym urządzeniu nazwę, typ,
plakietkę trybu Mobile, datę parowania, ważność poświadczenia i jego stan
(aktywne, zbliża się wygaśnięcie, wygasłe, odwołane). Dostępne operacje:

| Operacja | Skutek |
|---|---|
| Zmiana nazwy urządzenia | Nadanie nazwy rozpoznawalnej w rejestrze |
| Odwołanie dostępu | Unieważnienie poświadczenia i zamknięcie połączeń urządzenia — wylogowanie zdalne |
| Odwołanie dostępu wszystkich urządzeń | Unieważnienie wszystkich wydanych poświadczeń |
| Odświeżenie poświadczenia | Wydanie poświadczenia o nowym terminie ważności |

**Przy utracie urządzenia:** rozpoznać wiersz po nazwie, typie i dacie ostatniego
połączenia, odwołać poświadczenie. Zdarzenie zostaje zapisane w dzienniku audytu.

### 6.4. Synchronizacja

Zmiana dokonana na jednym urządzeniu jest przekazywana przez serwer do pozostałych
na żywo. Synchronizacja jest stałą właściwością platformy, nie ustawieniem do
włączenia — sekcja Urządzenia prezentuje wyłącznie jej bieżący stan.

---

## 7. Okno Ustawień — dziewięć sekcji

Okno Ustawień otwiera się z listwy ustawień na stronie głównej. Dotyczy poziomu
aplikacji: konta, dostępu, wyglądu, urządzeń i powiadomień. Konfiguracja
zachowania platformy w pracy — modele, procesy, izolacja, pamięć — leży w oknie
konfiguracji (rozdz. 8).

| Sekcja | Zawartość |
|---|---|
| **Konto** | Login, adres e-mail uwierzytelniający, zmiana hasła, awatar, data utworzenia konta |
| **Uwierzytelnianie** | Metody logowania: hasło i metody dodatkowe; stan wymogu logowania |
| **Wygląd i język** | Motyw jasny albo ciemny, język interfejsu |
| **Urządzenia** | Rejestr urządzeń połączonych, parowanie urządzenia przenośnego, stan synchronizacji |
| **Powiadomienia** | Przełącznik główny, klasy zdarzeń, kanał dostarczenia |
| **Konta modeli** | Wykaz kont modeli, katalogi profili, przełączanie konta aktywnego |
| **Okna komunikacji** | Ustawienia Chat Window i Execution Loop Window |
| **Warstwy widoczności** | Przypisanie warstw widoczności funkcji do roli użytkownika |
| **Always On Display** | Obecność i tryb funkcji globalnej, reguły wyzwalania sugestii, wyciszanie, tor głosowy |

Sekcja Konto otwiera się domyślnie przy wejściu do okna. Pozostałe sekcje wybiera
się selektorem `☰`. Sekcja Warstwy widoczności należy do funkcji eksperckich
i osiągalna jest z trybu administracyjnego, wyszukiwarki funkcji, skrótu
klawiszowego albo polecenia w języku naturalnym w Chat Window.

Każda pozycja ustawienia ma objaśnienie kontekstowe `[?]` mówiące, **co dane
ustawienie robi i jaki ma wpływ na działanie platformy**.

---

## 8. Okno konfiguracji — warstwy i zakresy

### 8.1. Dwa punkty wejścia

| Punkt wejścia | Zakres zmiany |
|---|---|
| Listwa ustawień na stronie głównej | Pełne okno konfiguracji: warstwa domyślna i warstwa sesji, wszystkie zakresy |
| Menu kontekstowe okna operacyjnego | Szybka zmiana bieżącej sesji — podzbiór ustawień warstwy sesji, bez przerywania pracy |

Oba wejścia prowadzą do **jednego modelu konfiguracji**, nie do dwóch równoległych
mechanizmów.

### 8.2. Cztery warstwy konfiguracji

```
[ GLOBALNA ]     cała platforma; wartość bazowa
      ▼ dziedziczy
[ ŚRODOWISKO ]   TalkIn · WorkSpace · CodeStudio · MultitaskingAI
      ▼ dziedziczy
[ PROJEKT ]      jeden projekt modułu Workspace
      ▼ dziedziczy
[ SESJA ]        jedna karta sesji — warstwa najwęższa
```

Rozstrzyganie wartości idzie od warstwy najwęższej do najszerszej — obowiązuje
pierwsza warstwa, na której ustawienie zostało określone; gdy żadna nie została
ustawiona, obowiązuje wartość domyślna platformy. **Zmiana na warstwie węższej
nakłada się na szerszą, nie zmieniając jej.** Ustawienie retencji historii dla
jednej sesji nie zmienia retencji obowiązującej pozostałe sesje.

| Warstwa | Zakres zmiany dokonanej na tej warstwie |
|---|---|
| Sesja | Tylko bieżąca karta sesji |
| Projekt | Jeden projekt modułu Workspace |
| Środowisko | Jedno z czterech środowisk |
| Globalna | Cała platforma |

### 8.3. Siedemnaście zakresów ustawień

Lewa kolumna okna konfiguracji prowadzi po zakresach; prawa prezentuje pozycje
wybranego zakresu wraz z wartością bieżącą, poziomem warstwy i objaśnieniem `[?]`.

| Zakres | Co konfiguruje |
|---|---|
| **Aplikacja** | Motyw wizualny, język interfejsu, metody uwierzytelniania, urządzenia połączone, powiadomienia, wymóg logowania |
| **Procesy** | Osiem zakresów izolacji technicznej procesu sesji (rozdz. 11) |
| **Akcje** | Akcje kolejki, tryb pracy w tle zadania, zasięg kolejki, obsługa błędów, spięcie silnika kolejek z MultitaskingAI |
| **Zachowanie modeli** | Kanał modelu, model bazowy, parametry wywołania, relacje między wykonawcami, liczba podagentów, dane dostępowe kanału |
| **Tożsamość modeli** | Nazwa i tożsamość agenta, model bazowy agenta, wcielenie roli kontrolnej, uprawnienia agenta |
| **Prompty systemowe** | Instrukcje systemowe projektu i agenta, budowa promptów przez rolę Coordinator |
| **Rozszerzenia** | Wtyczki, umiejętności, konektory, serwery MCP (rozdz. 12) |
| **Integracje** | Jawne powiązania między środowiskami, modułami i komponentami własnymi |
| **Historia** | Przechowywanie i dostępność historii rozmów i sesji, retencja |
| **Pamięć** | Zasoby pamięci kontekstowej na czterech poziomach |
| **Izolacja** | Izolacja kontekstu i izolacja techniczna — okno trzypanelowe (rozdz. 11) |
| **Karty sesji** | Ustawienia karty sesji jako jednostki pracy równoległej |
| **Komponenty własne** | Tworzenie, zapis i przypisanie automatyki, agenta, projektu i profilu asystenta |
| **Chat Window** | Szerokość kolumny, zachowanie strumienia, zatwierdzanie i przerywanie działań, zakres pamięci |
| **Execution Loop Window** | Równoległość zadań, polityka ponowień, progi kontroli jakości, zakres autonomii, punkty zatwierdzeń, sterowanie przebiegiem |
| **Warstwy widoczności funkcji** | Przypisanie funkcji do warstw, profile według roli, wyszukiwarka funkcji, skróty, tryb administracyjny |
| **Always On Display** | Obecność i tryb funkcji, klasy zdarzeń wyzwalających sugestie, progi i częstotliwość, wyciszanie, tor głosowy |

**Zasada nadrzędna:** żaden poziom warstwy nie musi być skonfigurowany. Poziom
nieustawiony nie blokuje pracy — przejmuje wartość odziedziczoną.

---

## 9. Konfiguracja modeli: kanały i konta

### 9.1. Cztery kanały integracji

Kanał określa, w jaki sposób model jest fizycznie osiągalny. Sesja i rola nie
odwołują się do kanału wprost — odwołują się do skonfigurowanego **kanału modelu**,
który dopiero wskazuje drogę dostępu.

| Kanał | Sposób połączenia | Dane uwierzytelniające | Kiedy wybrać |
|---|---|---|---|
| **API** | Żądanie HTTP do punktu końcowego dostawcy | Klucz dostępu | Dostawca udostępnia interfejs programistyczny; pożądana przewidywalność kosztów i stabilność połączenia |
| **CLI** | Wywołanie narzędzia wiersza poleceń dostawcy | Token narzędzia | Model dostępny przede wszystkim jako narzędzie wiersza poleceń |
| **SSH** | Połączenie z hostem zdalnym | Host, konto, klucz albo hasło | Własny model pomocniczy na własnej infrastrukturze; pełna kontrola nad środowiskiem wykonawczym |
| **HTTP** | Interakcja z interfejsem webowym dostawcy | Dane logowania właściwe interfejsowi | Model dostępny wyłącznie przez przeglądarkę |

### 9.2. Ustawienia kanału

Zakres **Zachowanie modeli** okna konfiguracji obejmuje:

| Ustawienie | Warstwy |
|---|---|
| Kanał modelu — API, CLI, SSH albo HTTP | sesja, rola |
| Model bazowy udostępniany przez kanał | sesja, rola, agent |
| Relacje między wykonawcami MultitaskingAI: praca niezależna, przekazywanie wyników, praca naprzemienna, praca iteracyjna | rola |
| Liczba podagentów Subagent Network — do piętnastu | rola |
| Dane dostępowe kanału | globalna, sesja, rola |

### 9.3. Konta modeli

Sekcja **Konta modeli** okna Ustawień prowadzi rejestr kont u dostawców wraz
z katalogami profili. Przełączenie konta aktywnego jest czynnością jawną —
platforma nie zmienia konta w tle. Przy pustym rejestrze sekcja prezentuje stan
pusty z wejściem do dodania pierwszego konta.

**Dane dostępowe kanałów modeli przechowywane są poza bazą danych platformy.**
W konfiguracji zapisywane jest odwołanie do nich, nie sam materiał tajny.

---

## 10. Konfiguracja obu kanałów komunikacji operacyjnej

### 10.1. Chat Window

| Ustawienie | Czego dotyczy |
|---|---|
| Szerokość kolumny | Udział okna w obszarze roboczym |
| Zachowanie strumienia | Sposób prezentacji odpowiedzi przekazywanej na żywo |
| Zatwierdzanie i przerywanie działań | Które działania wymagają potwierdzenia przed wykonaniem |
| Zakres pamięci | Co okno pamięta między turami i między sesjami |

### 10.2. Execution Loop Window

| Ustawienie | Czego dotyczy |
|---|---|
| Równoległość zadań | Ile zadań pętli może biec jednocześnie |
| Polityka ponowień | Zachowanie po niepowodzeniu zadania |
| Progi kontroli jakości | Kiedy wynik uznaje się za przyjęty, a kiedy trafia do ponowienia |
| Zakres autonomii Wykonawcy | Co Wykonawca może zrobić bez pytania |
| Punkty zatwierdzeń | Miejsca obowiązkowej decyzji Operatora |
| Sterowanie przebiegiem | Dostępność wstrzymania, wznowienia, przerwania i korekty zlecenia |

**Działania nieodwracalne wymagają zatwierdzenia Operatora** niezależnie od
ustawionego zakresu autonomii. Przebieg pętli i decyzje zapisywane są w dzienniku
audytu.

---

## 11. Punkty izolacji i profile izolacji

Pozycja **Izolacja** w oknie konfiguracji otwiera osobne, trzypanelowe okno.
Izolacja jest w tym produkcie możliwością, nie wymogiem — rozstrzyga o niej
Operator.

### 11.1. Macierz izolacji

| Grupa | Pozycje | Stany |
|---|---|---|
| **Izolacja kontekstu** | historia, pamięć, kontekst | współdzielone \| odrębne |
| **Izolacja techniczna** | katalog roboczy sesji, środowisko procesu, katalog danych i konfiguracji modelu, dostęp sieciowy, zakres odczytu i zapisu plików, konto i token per sesja, model procesu, serwer wykonania | włączony \| wyłączony |

Stanem wyjściowym każdej pozycji izolacji technicznej jest „wyłączony", a modelu
procesu — model współdzielony.

### 11.2. Siedem poziomów zasięgu

| Poziom | Przykład zastosowania | Pierwszeństwo |
|---|---|---|
| Globalny | Domyślna polityka izolacji całej platformy | najniższe |
| Środowisko | Inna polityka dla CodeStudio niż dla TalkIn | ↑ |
| Moduł | Odrębna pamięć przypisana modułowi Developer | ↑ |
| Para modułów | Współdzielenie operacji kontekstowych między Studio a Translate | ↑ |
| Projekt | Rozdzielenie kontekstu dwóch projektów modułu Workspace | ↑ |
| Karta sesji | Jednorazowe współdzielenie historii dwóch otwartych kart | ↑ |
| Rola MultitaskingAI | Profil izolacji przypisany roli wykonawcy | najwyższe |

Reguła z poziomu węższego wygrywa z regułą z poziomu szerszego; brak ustawienia
oznacza dziedziczenie z poziomu bezpośrednio szerszego.

### 11.3. Profile izolacji

Ustalony układ przełączników zapisuje się jako **profil izolacji** i przypisuje
do sesji, roli albo projektu. Okno udostępnia **podgląd polityki efektywnej** —
zestawienie tego, co faktycznie obowiązuje po nałożeniu wszystkich poziomów.
Zalecana kolejność pracy: ustalić politykę globalną, następnie odstępstwa dla
środowisk i modułów, na końcu zapisać profile dla ról i projektów.

---

## 12. Rozszerzenia: wtyczki, umiejętności, konektory, serwery MCP

Wszystkie rodzaje rozszerzeń działają w jednolitym kontrakcie: tożsamość
i rejestracja, stan włączenia zmienialny w każdej chwili, konfiguracja właściwa
rodzajowi oraz zakres uprawnień przypisany przy podłączeniu do agenta lub modułu.

| Cecha | **Danaco Plugin** | **Personal** |
|---|---|---|
| Pochodzenie | Dostarczane z platformą | Instalowane przez Operatora z urządzenia |
| Miejsce przechowywania | Katalog platformy na serwerze | Katalog Operatora na serwerze |
| Rejestracja | Automatyczna | Po zakończeniu przesyłania |
| Aktualizacja | Wraz z aktualizacją platformy | Ponowne przesłanie przez Operatora |
| Stan wyjściowy | **Włączone** | **Wyłączone** — wymaga świadomego włączenia |

Rozszerzenie własne przesyłane jest z urządzenia kanałem komunikacji i zapisywane
na serwerze; na urządzeniu nie pozostaje trwała kopia. Uprawnienia rozszerzenia
podlegają tym samym zasadom izolacji technicznej co pozostałe elementy platformy —
zakres nadaje się przy podłączeniu, nie przy instalacji.

Rozszerzenia podłącza się do agenta w module Agents albo do modułu w zakresie
Rozszerzenia okna konfiguracji.

---

## 13. Warstwy widoczności funkcji

Interfejs ujawnia możliwości stopniowo. Przypisanie funkcji do warstw
i profile warstw według roli użytkownika są **ustawieniem konfiguracji**.

| Warstwa | Zawartość | Dostęp |
|---|---|---|
| 1 | Chat Window, aktywne okno wiodące, kontekst pracy, nawigacja, wskaźniki stanu | bez interakcji |
| 2 | Wybór modelu, wykonawcy, środowiska, trybu pracy, parametry przepływu | ikona, przycisk, znacznik kontekstowy |
| 3 | Zestawy akcji, ustawienia szybkie, warianty operacji | `⋮`, `☰`, menu kontekstowe, panel |
| 4 | Operacje zaawansowane, tryb administracyjny, diagnostyka niskiego poziomu | polecenie języka naturalnego, skrót klawiszowy, wyszukiwarka funkcji |

Trzy drogi dotarcia do funkcji ukrytej działają zawsze: **wyszukiwarka funkcji**,
**skrót klawiszowy** i **polecenie w języku naturalnym w Chat Window**. Zakres
wyników wyszukiwarki odpowiada uprawnieniom roli. **Tryb administracyjny** odsłania
warstwę czwartą w całości.

---

## 14. Profile konfiguracji

Ustalony zestaw wartości zapisuje się jako profil konfiguracji i przypisuje —
do środowiska, projektu, sesji albo roli. Profile skracają przygotowanie pracy
powtarzalnej: zamiast ustawiać kilkanaście pozycji przy każdym nowym projekcie,
Operator wczytuje profil. W środowisku MultitaskingAI odpowiednikiem profilu na
poziomie zespołu jest zapisana konfiguracja w sekcji **Zespoły** panelu
orkiestracji — presety ról, powiązań i kolejek gotowe do wczytania albo
zduplikowania.

---

## 15. Powiadomienia

| Element sekcji | Zawartość |
|---|---|
| Przełącznik główny | Włączenie i wyłączenie powiadomień |
| Klasy zdarzeń | Rodzaje zdarzeń platformy objęte powiadomieniem |
| Kanał dostarczenia | Droga dotarcia powiadomienia, w tym powiadomienia wypychane na urządzenie przenośne |

Powiadomienia są konfigurowalne na warstwie globalnej i na warstwie sesji; ich
wartością domyślną jest stan aktywny. Powiadomienie na urządzeniu przenośnym jest
wejściem do interwencji — z jego poziomu Operator przechodzi wprost do zatwierdzenia,
wstrzymania albo korekty procesu.

---

## 16. Kolejność zalecana przy uruchamianiu nowej instalacji

| Krok | Czynność | Miejsce |
|---|---|---|
| 1 | Instalacja pakietu klienckiego i wskazanie adresu serwera | rozdz. 3 |
| 2 | Rejestracja konta właściciela i potwierdzenie adresu e-mail | rozdz. 4 |
| 3 | Włączenie metod dodatkowych uwierzytelniania | Ustawienia → Uwierzytelnianie |
| 4 | Motyw i język interfejsu | Ustawienia → Wygląd i język |
| 5 | Dodanie pierwszego konta modelu i wybór kanału | Ustawienia → Konta modeli; konfiguracja → Zachowanie modeli |
| 6 | Ustawienie zakresu autonomii i punktów zatwierdzeń pętli wykonawczej | konfiguracja → Execution Loop Window |
| 7 | Polityka historii i pamięci | konfiguracja → Historia, Pamięć |
| 8 | Polityka izolacji na poziomie globalnym; profile dla ról i projektów | konfiguracja → Izolacja |
| 9 | Włączenie potrzebnych rozszerzeń i nadanie im uprawnień | konfiguracja → Rozszerzenia |
| 10 | Sparowanie urządzenia przenośnego | Ustawienia → Urządzenia |
| 11 | Zakres powiadomień i kanał dostarczenia | Ustawienia → Powiadomienia |
| 12 | Tryb Always On Display i reguły sugestii | Ustawienia → Always On Display |

Po kroku 5 platforma jest gotowa do pracy; kroki 6–12 dostrajają ją do sposobu
pracy Operatora i mogą być wykonane później.

---

## 17. Diagnostyka

| Objaw | Przyczyna | Postępowanie |
|---|---|---|
| Okno startowe pokazuje komunikat o braku połączenia | Serwer nieosiągalny, nieprawidłowy adres, przerwanie sieci | Sprawdzić dostęp sieciowy do serwera i poprawność adresu; użyć kontrolki „Spróbuj ponownie" |
| Okno logowania pojawia się na urządzeniu, które wcześniej wchodziło bez niego | Token dostępu urządzenia wygasł albo został odwołany | Zalogować się ponownie; sprawdzić stan poświadczenia w rejestrze urządzeń |
| Urządzenie przenośne nie kończy parowania | Kod wygasł, został użyty albo zgłoszenie odrzucono | Wydać nowy kod w modalu parowania i powtórzyć odczyt |
| Model nie odpowiada w sesji | Brak konta modelu, brak danych dostępowych kanału albo niedostępność dostawcy | Sprawdzić rejestr kont modeli i ustawienie kanału w zakresie Zachowanie modeli |
| Zadanie kolejki zatrzymuje się w tym samym miejscu | Polityka obsługi błędów kolejki albo próg kontroli jakości | Sprawdzić ustawienia zakresu Akcje oraz progi kontroli jakości pętli wykonawczej |
| Funkcja obecna w dokumentacji nie jest widoczna w interfejsie | Funkcja należy do wyższej warstwy widoczności | Użyć wyszukiwarki funkcji, skrótu klawiszowego albo polecenia w Chat Window; sprawdzić profil warstw przypisany roli |
| Zmiana ustawienia nie działa w bieżącej sesji | Ustawienie określone na warstwie węższej nakłada się na wprowadzoną wartość | Sprawdzić poziom warstwy przy pozycji ustawienia i podgląd polityki efektywnej |
| Zmiana widoczna na jednym urządzeniu, brak na drugim | Utrata połączenia drugiego urządzenia | Sprawdzić stan synchronizacji w sekcji Urządzenia i połączenie tego urządzenia |
| Utrata dostępu do konta | Brak dostępu do metod logowania | Uruchomić odzyskiwanie konta adresem e-mail uwierzytelniającym z okna logowania |

Diagnostyka techniczna prowadzona po stronie platformy — logi, błędy, wydajność —
dostępna jest w module Diagnostics środowiska CodeStudio.

**Zgłoszenia.** Nieprawidłowość, której nie usuwa żadne z powyższych
postępowań, zgłasza się na adres **support@danaco-group.pl**, wskazując wersję
produktu, urządzenie, moduł oraz okoliczności pozwalające odtworzyć zdarzenie.
Tym samym adresem zgłasza się podatności bezpieczeństwa.

---

## 18. Aktualizacja i deinstalacja

**Aktualizacja platformy** odbywa się centralnie, po stronie serwera. Operator nie
wykonuje w tym zakresie żadnej czynności; po ponownym połączeniu klient pracuje
z nową wersją platformy.

**Aktualizacja pakietu klienckiego** wymagana jest wtedy, gdy zmienia się samo
okno aplikacji. Przebiega jak instalacja pierwotna (rozdz. 3) i nie usuwa ani
konta, ani zapisanego adresu serwera, ani tokenu dostępu urządzenia.

**Deinstalacja klienta** usuwa aplikację z urządzenia. Nie usuwa konta, sesji,
historii, projektów ani żadnych innych danych — stan trwały platformy istnieje
wyłącznie na serwerze. Zalecane jest odwołanie dostępu urządzenia w rejestrze
urządzeń przed deinstalacją; urządzenie usunięte bez odwołania poświadczenia
pozostaje w rejestrze do czasu odwołania.

| Platforma | Deinstalacja |
|---|---|
| Windows | Standardowe usuwanie aplikacji z systemu |
| Linux (`deb`) | `sudo apt remove danaco-console` |
| Linux (`AppImage`) | Usunięcie pliku |

---

*Danaco Console — Platforma AI Workspace OS · v2.0*
*© 2026 Danaco Holding Group Sp. z o.o. Wszelkie prawa zastrzeżone — Dariusz Naharnowicz*
