# Poczta transakcyjna Danaco Console — opracowanie techniczne

| | |
|---|---|
| **Zakres** | trzy szablony wiadomości wychodzących z adresu `noreply@danaco-group.pl` |
| **Odbiorca** | Operator Danaco Console |
| **Źródła projektowe** | `03-marka/ksiega-znaku.md` · `03-marka/typografia-i-glos.md` · `03-marka/zastosowania-marki.md` (rozdz. 4.6) |
| **Pliki kodu** | `html/list-01-aktywacja-konta.html` · `html/list-02-uwierzytelnienie-logowania.html` · `html/list-03-autoresponder-brak-skrzynki.html` |
| **Podgląd** | `podglad.html` |
| **Wariant barwny** | jasny (jedyny) |

---

## 1. Trzy okoliczności wysyłki

| Nr | Plik | Wyzwalacz w rdzeniu | Temat wiadomości |
|---:|---|---|---|
| 1 | `list-01-aktywacja-konta.html` | utworzenie konta Operatora, przed pierwszym logowaniem | `Danaco Console — kod aktywacji konta: {{code}}` |
| 2 | `list-02-uwierzytelnienie-logowania.html` | logowanie wymagające drugiego składnika | `Danaco Console — kod logowania: {{code}}` |
| 3 | `list-03-autoresponder-brak-skrzynki.html` | odebranie wiadomości na `noreply@danaco-group.pl` | `Adres noreply@danaco-group.pl nie przyjmuje korespondencji` |

Kod w temacie skraca ścieżkę Operatora — kod widać na liście wiadomości bez
otwierania listu. Rozwiązanie jest przyjęte świadomie: kod jest jednorazowy,
związany z sesją i wygasa po 60 minutach, więc jego obecność w nagłówku tematu
nie podnosi ryzyka ponad poziom samej dostawy wiadomości.

---

## 2. Decyzje projektowe wobec księgi znaku

### 2.1 Jeden akcent na kadr

Księga znaku, rozdz. 14.4: w jednym kadrze widoczny jest jeden akcent
sygnałowy. W szablonach akcentem pozostaje **kropka znaku** w nagłówku.
Konsekwencje:

- panele informacyjne i panel kodu są **neutralne** — tło `#FAFAFA`, obrys `#E3E3E3`;
  nie użyto rodziny `--dn-sygnal-100/200`, mimo że treść jest informacyjna,
- przycisk działania głównego to **inwersja atramentu** — wypełnienie `#181818`,
  tekst `#FFFFFF` (księga znaku, rozdz. 14.4: „czerń działa, sygnał wskazuje"),
- barwa `#3B6FE0` nie występuje poza plikiem znaku i poza barwą odnośnika tekstowego.

### 2.2 Trzy kroje w trzech rolach

| Krój | Rola w liście | Łańcuch zapasowy |
|---|---|---|
| Space Grotesk | nagłówek listu, nazwa produktu | `'IBM Plex Sans','Segoe UI',Helvetica,Arial,sans-serif` |
| IBM Plex Sans | treść, przycisk, stopka | `'Segoe UI',Helvetica,Arial,sans-serif` |
| IBM Plex Mono | kod, znacznik czasu, adres IP, etykiety metadanych | `'SFMono-Regular',Consolas,'Courier New',monospace` |

Kroje marki **nie są wgrywane** — klienty pocztowe w większości blokują
`@font-face`. Deklaracja `font-family` wymienia krój marki jako pierwszy
(zadziała u Operatora z krojem zainstalowanym w systemie), dalej idą kroje
systemowe. **Inter nie występuje w łańcuchu zapasowym w żadnej roli**
(typografia i głos, rozdz. 2).

### 2.3 Głos

Pięć cech głosu (typografia i głos, rozdz. 14) przełożone na treść listów:

| Cecha | Zastosowanie |
|---|---|
| precyzyjny | podana liczba minut ważności, znacznik czasu żądania, adres źródłowy |
| rzeczowy | zero wykrzykników, zero przeprosin, zero słów wartościujących |
| sprawczy | każdy list kończy się następnym krokiem albo informacją, co zrobi system |
| bez żargonu | „kod logowania", nie „token OTP"; „konto", nie „tenant" |
| bez przechwałek | brak zdań o produkcie; list mówi wyłącznie o zdarzeniu |

Zwrot do odbiorcy: **Operator**, forma bezosobowa w opisach, rozkaźnik w akcjach.

### 2.4 Barwy — żetony użyte w szablonach

| Rola w liście | Wartość | Żeton systemu |
|---|---|---|
| tło strony (poza kartą) | `#F4F4F4` | `--dn-szary-50` |
| powierzchnia karty | `#FFFFFF` | `--dn-szary-0` |
| listwa górna karty | `#181818` | `--dn-szary-900` |
| tekst podstawowy | `#181818` | `--dn-szary-900` |
| tekst drugorzędny | `#4A4A4A` | `--dn-szary-700` |
| tekst stopki i etykiet | `#616161` | `--dn-szary-600` |
| obrysy i linie | `#E3E3E3` | `--dn-szary-150` |
| tło panelu kodu | `#FAFAFA` | `--dn-szary-25` |
| odnośnik tekstowy | `#3B6FE0` | `--dn-sygnal-500` |
| kropka znaku | `#3B6FE0` | `--dn-sygnal-500` |

Kontrast tekstu stopki `#616161` na `#FFFFFF` wynosi **5,99 : 1** — spełnia
WCAG AA dla tekstu zwykłego. Świadomie **nie użyto** `#7C7C7C`
(`--dn-szary-500`, żeton metadanych interfejsu): na białym tle daje 4,06 : 1,
czyli poniżej progu 4,5 : 1 wymaganego dla stopni poniżej 18 px.

---

## 3. Ograniczenia techniczne nośnika

| Ograniczenie | Przyczyna | Rozwiązanie w szablonach |
|---|---|---|
| brak `<style>` w `<head>` | Gmail usuwa arkusz w treści wiadomości | wszystkie style w atrybucie `style` |
| brak flexbox i grid | Outlook renderuje silnikiem Worda | układ na `<table role="presentation">` |
| brak `@font-face` | blokada w większości klientów | łańcuch zapasowy do krojów systemowych |
| wymuszona biel tła | tryb ciemny części klientów odwraca barwy | karta na `#FFFFFF`, brak wariantu ciemnego |
| blokada obrazów domyślnie | polityka klienta pocztowego | znak z `alt`, układ czytelny bez obrazu |
| `data:` w `<img>` bywa odrzucane | Gmail blokuje schemat `data:` w obrazach | znak jako załącznik `cid:` — patrz rozdz. 4 |
| przycięcie treści w Gmailu | limit 102 KB na wiadomość | każdy szablon poniżej 20 KB bez znaku |

### 3.1 Znak w nagłówku — trzy warianty osadzenia

1. **`cid:` (zalecany produkcyjnie).** Plik
   `03-marka/zastosowania/png/znak-poczty-jasny@2x.png` dołączany jako część
   `multipart/related` z nagłówkiem `Content-ID: <danaco-lockup>`,
   w szablonie `src="cid:danaco-lockup"`. Działa we wszystkich badanych klientach,
   nie wymaga serwera obrazów, nie generuje śladu odczytu.
2. **Adres bezwzględny.** `src="https://…/znak-poczty-jasny@2x.png"` — wymaga
   wystawienia zasobu publicznie; w instalacji lokalnej zwykle niedostępne.
3. **`data:` URI.** Użyte **wyłącznie w `podglad.html`**, żeby plik podglądu
   był samowystarczalny. Nie stosować w wysyłce.

Szablony produkcyjne mają w tym miejscu zmienną `{{logo_src}}`.

Wymiary: plik 336 × 145 px (2×), wyświetlanie 168 × 72 px — lockup kompaktowy,
zgodnie z zastosowaniami marki, rozdz. 4.6.

---

## 4. Zmienne szablonu

Identyfikatory po angielsku (standard nazewnictwa Danaco), składnia `{{ }}`.

### 4.1 Wspólne dla wszystkich trzech listów

| Zmienna | Typ | Przykład | Opis |
|---|---|---|---|
| `{{logo_src}}` | tekst | `cid:danaco-lockup` | źródło obrazu znaku |
| `{{recipient_address}}` | tekst | `a.kowalska@…` | adres odbiorcy, powtórzony w stopce |
| `{{year}}` | liczba | `2026` | rok w nocie stopki |
| `{{support_address}}` | tekst | `support@danaco-group.pl` | adres skrzynki obsługiwanej |

### 4.2 List 1 — aktywacja konta

| Zmienna | Typ | Przykład | Opis |
|---|---|---|---|
| `{{code}}` | tekst, 6 cyfr | `418 402` | kod aktywacyjny; w treści rozdzielany spacją co trzy znaki |
| `{{expiry_minutes}}` | liczba | `60` | ważność kodu w minutach |
| `{{requested_at}}` | tekst | `29.08.2026, 17:22 CEST` | znacznik czasu utworzenia konta |
| `{{activation_url}}` | URL | `https://console…/aktywacja` | cel przycisku działania głównego |

### 4.3 List 2 — uwierzytelnienie logowania

| Zmienna | Typ | Przykład | Opis |
|---|---|---|---|
| `{{code}}` | tekst, 6 cyfr | `730 155` | kod logowania |
| `{{expiry_minutes}}` | liczba | `60` | ważność kodu w minutach |
| `{{requested_at}}` | tekst | `29.08.2026, 17:22 CEST` | znacznik czasu żądania |
| `{{ip_address}}` | tekst | `10.14.2.37` | adres źródłowy żądania |
| `{{device}}` | tekst | `Windows 11 · Danaco Console 2.0` | opis urządzenia i wydania powłoki |

### 4.4 List 3 — autoresponder

| Zmienna | Typ | Przykład | Opis |
|---|---|---|---|
| `{{original_subject}}` | tekst | `Prośba o zmianę limitu` | temat wiadomości odrzuconej |
| `{{received_at}}` | tekst | `29.08.2026, 17:22 CEST` | czas odebrania wiadomości |

### 4.5 Reguła sanityzacji [NIENEGOCJOWALNE]

`{{original_subject}}` pochodzi z wiadomości nadesłanej z zewnątrz i jest
jedyną zmienną spoza rdzenia. Przed podstawieniem:

- ucieczka HTML dla `& < > " '`,
- obcięcie do 120 znaków,
- usunięcie znaków sterujących i złamań wiersza (ochrona przed wstrzyknięciem
  nagłówka wiadomości).

Bez ucieczki temat nadesłany z zewnątrz staje się wektorem wstrzyknięcia
znaczników do listu wychodzącego.

---

## 5. Bezpieczeństwo treści

| Zasada | Uzasadnienie |
|---|---|
| kod wyłącznie w treści listu i w temacie, nigdy w logu serwera | dziennik wysyłki zapisuje `message_id` i adres, nie zawartość |
| brak danych ze spraw i akt w liście transakcyjnym | listy dotyczą konta, nie pracy Operatora |
| brak odnośnika śledzącego i piksela odczytu | list nie zbiera zachowania odbiorcy |
| brak nazwiska Operatora w temacie | temat bywa widoczny w powiadomieniu na ekranie blokady |
| autoresponder nie cytuje treści nadesłanej wiadomości | tylko temat, po sanityzacji |
| `{{ip_address}}` widoczny wyłącznie w liście logowania | Operator musi rozpoznać cudze żądanie |

### 5.1 Uwierzytelnienie nadawcy

Adres `noreply@danaco-group.pl` wymaga kompletu rekordów: **SPF**, **DKIM**
(selektor osobny dla poczty transakcyjnej) oraz **DMARC** w polityce
`p=reject`. Bez tego listy z kodem trafiają do folderu wiadomości
niechcianych, a domena jest podatna na podszycie.

### 5.2 Skrzynka nieobsługiwana

Adres `noreply@danaco-group.pl` nie prowadzi skrzynki odbiorczej. Reguła na
serwerze poczty: każda wiadomość przychodząca uruchamia list 3 i zostaje
odrzucona. Ograniczenia autorespondera:

- **jedna odpowiedź na adres na dobę** — ochrona przed pętlą odpowiedzi,
- **brak odpowiedzi**, gdy wiadomość przychodząca ma nagłówek
  `Auto-Submitted` inny niż `no`, `Precedence: bulk` albo pusty adres zwrotny,
- list 3 wysyłany z nagłówkiem `Auto-Submitted: auto-replied`.

### 5.3 Nagłówki wymagane w listach 1 i 2

```
Auto-Submitted: auto-generated
X-Auto-Response-Suppress: All
List-Unsubscribe:            (nieobecny — poczta transakcyjna nie podlega wypisaniu)
```

---

## 6. Struktura pliku szablonu

Wszystkie trzy pliki mają ten sam układ pionowy:

```
┌──────────────────────────────────────────────┐
│ preheader (ukryty — pierwszy wiersz podglądu)│
├──────────────────────────────────────────────┤
│ listwa górna 3 px · #181818                  │
│ znak 168 × 72 px                             │
│ linia 1 px · #E3E3E3                         │
├──────────────────────────────────────────────┤
│ nagłówek listu   Space Grotesk 24 px / 600   │
│ akapit wprowadzający  Plex Sans 15 px        │
├──────────────────────────────────────────────┤
│ panel treści właściwej                       │
│   list 1 i 2 → kod, Plex Mono 32 px          │
│   list 3     → nota o braku skrzynki         │
├──────────────────────────────────────────────┤
│ szyna metadanych  Plex Mono 11 px wersaliki  │
├──────────────────────────────────────────────┤
│ nota bezpieczeństwa / następny krok          │
├──────────────────────────────────────────────┤
│ linia 1 px · stopka nadawcy Plex Sans 12 px  │
└──────────────────────────────────────────────┘
```

Szerokość karty **600 px**, marginesy wewnętrzne **32 px** (24 px poniżej
480 px szerokości okna). Poza kartą tło `#F4F4F4`.

### 6.1 Szyna metadanych

Para etykieta–wartość złożona IBM Plex Mono. Etykieta: 11 px, wersaliki,
odstęp liter 0,14 em, barwa `#616161`. Wartość: 13 px, barwa `#181818`.
Zgodnie z typografią i głosem, rozdz. 14.4: krój maszyny niesie **wartość**,
nie zdanie.

---

## 7. Lista kontrolna przed wydaniem

- [ ] rekordy SPF, DKIM i DMARC dla `danaco-group.pl` odpowiadają poprawnie
- [ ] znak dołączony jako `cid:`, `{{logo_src}}` podstawione
- [ ] `{{original_subject}}` przechodzi przez ucieczkę HTML
- [ ] atrybut `alt` znaku niepusty
- [ ] list czytelny przy zablokowanych obrazach
- [ ] rozmiar wiadomości poniżej 102 KB
- [ ] kontrola w Gmailu (przeglądarka, Android, iOS), Outlooku 2019 i 365,
      Apple Mail, Thunderbird
- [ ] ograniczenie częstotliwości autorespondera aktywne
- [ ] `Auto-Submitted` ustawione we wszystkich trzech listach
- [ ] kod nieobecny w dzienniku serwera poczty

---

## 8. Czego w opracowaniu nie ma

| Poza zakresem | Uwaga |
|---|---|
| wersja `text/plain` | wymagana produkcyjnie jako druga część `multipart/alternative`; do dopracowania |
| wariant ciemny | odrzucony — biel wymuszona przez klienty pocztowe (zastosowania marki, rozdz. 3) |
| listy powiadomień roboczych | inny nośnik, inna częstotliwość, osobne opracowanie |
| wersje językowe inne niż polska | głos marki opisany wyłącznie dla polszczyzny |
