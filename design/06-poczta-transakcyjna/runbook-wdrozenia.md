# Runbook · wdrożenie poczty transakcyjnej na Stalwart

| | |
|---|---|
| **Zakres** | uruchomienie wysyłki z `noreply@danaco-group.pl` na serwerze Stalwart w instalacji lokalnej |
| **Dotyczy** | Stalwart 0.16.x (najnowsze wydanie w chwili pisania: **0.16.19**, 24 sierpnia 2026) |
| **Odbiorcy** | administrator serwera poczty, osoba wdrażająca rdzeń Danaco Console |
| **Poprzedza** | pierwszą wysyłkę listów z katalogu `html/` i `text/` |
| **Dokument towarzyszący** | `opracowanie-techniczne.md` (szablony), `rdzen/mail/` (składanie wiadomości) |

Bez tych czynności szablony są bezużyteczne: listy z kodem trafią do folderu
wiadomości niechcianych, a domena pozostanie podatna na podszycie.

---

## 0. Zanim zaczniesz — model konfiguracji Stalwart 0.16

W wydaniu **0.16.0 (kwiecień 2026) struktura konfiguracji zmieniła się
zasadniczo**. Jeżeli pracowałeś na Stalwart 0.15 albo starszym, poniższe
różnice są istotne:

| Do 0.15 | Od 0.16 |
|---|---|
| duży plik TOML z sekcjami `[signature."x"]`, `[queue.throttle]` | mały `config.json` zawierający **wyłącznie** lokalizację bazy |
| ustawienia w pliku | ustawienia jako obiekty w bazie, edytowane przez WebUI albo `stalwart-cli` |
| nazwy sekcji z kropkami | nazwy pól obiektów w camelCase (`dkimSignDomain`, `maxRecipients`) |
| czasy jako `10d`, `1h` | liczby całkowite w **milisekundach**; rozmiary w bajtach |
| listy jako tablice TOML | zbiory jako obiekty JSON `{"klucz": true}` |

Plik startowy sprowadza się do wskazania magazynu:

```json
{"@type": "RocksDb", "path": "/var/lib/stalwart/"}
```

**Zmiany ustawień skompilowanych** (listenery, reguły MTA, podpisywanie DKIM,
telemetria) **nie wchodzą w życie same**. Po każdej zmianie:

```bash
stalwart-cli create Action/ReloadSettings
```

Dane katalogowe (konta, domeny, aliasy) działają od razu.

Wersję przypinaj — dla obrazu kontenera taguj gałąź (`v0.16`), nie `latest`.
Dokumentacja zapowiada 1.0.0 z zamrożonym schematem, ale termin już minął
i gałąź 0.16 nadal się rozwija.

---

## 1. Warstwa DNS

### 1.1 SPF

Jeden rekord `TXT` na domenę. Podwójny rekord SPF to błąd trwały — odbiorca
odrzuca wtedy ocenę w całości.

```
danaco-group.pl.  IN TXT  "v=spf1 ip4:<adres-serwera> -all"
```

`-all` (odrzucaj) zamiast `~all` (miękko). Poczta transakcyjna wychodzi
z jednego serwera, więc nie ma powodu na złagodzenie.

### 1.2 DKIM — rekord publiczny

Nazwa rekordu bierze się z selektora:

```
<selektor>._domainkey.danaco-group.pl.  IN TXT  "v=DKIM1; k=rsa; p=<klucz-publiczny>"
```

Dla klucza Ed25519: `k=ed25519`.

### 1.3 DMARC

Wdrażaj **stopniowo**. Natychmiastowe `p=reject` przy błędzie w SPF albo DKIM
kasuje całą pocztę wychodzącą, w tym listy z kodem logowania.

| Etap | Rekord | Czas |
|---|---|---|
| 1 | `v=DMARC1; p=none; rua=mailto:dmarc@danaco-group.pl; pct=100` | 2 tygodnie — zbieranie raportów |
| 2 | `v=DMARC1; p=quarantine; rua=…; pct=25` | tydzień |
| 3 | `v=DMARC1; p=reject; rua=…` | docelowo |

Przejście na kolejny etap dopiero wtedy, gdy raporty nie pokazują wysyłki
niepodpisanej. Rekord: `_dmarc.danaco-group.pl. IN TXT "…"`.

---

## 2. DKIM w Stalwart

### 2.1 Wytworzenie klucza

```bash
openssl genrsa -out /opt/stalwart/etc/private/dkim-transakcyjna.key 2048
openssl rsa -in /opt/stalwart/etc/private/dkim-transakcyjna.key -pubout
```

Prawa dostępu: właściciel procesu Stalwart, tryb `0600`. Klucz prywatny nie
wchodzi do repozytorium ani do kopii zapasowej dostępnej szerzej niż serwer.

### 2.2 Obiekt podpisu

Podpis to obiekt `DkimSignature` (WebUI: Management › Domains › DKIM
Signatures). Wariant obiektu zastępuje dawne pole `algorithm`:
`Dkim1RsaSha256` albo `Dkim1Ed25519Sha256`.

```json
{
  "@type": "Dkim1RsaSha256",
  "domainId": "<identyfikator obiektu Domain>",
  "selector": "transakcyjna",
  "privateKey": {"@type": "File", "filePath": "/opt/stalwart/etc/private/dkim-transakcyjna.key"},
  "headers": {"From": true, "To": true, "Subject": true, "Date": true, "Message-ID": true},
  "canonicalization": "relaxed/relaxed",
  "expire": 864000000,
  "report": true,
  "stage": "active"
}
```

`expire` w **milisekundach** (`864000000` = 10 dni). `domainId` to
identyfikator obiektu `Domain`, nie nazwa domeny.

O tym, co podpisywać, rozstrzyga pole `dkimSignDomain` na singletonie
`SenderAuth` (Settings › MTA › Inbound › Sender Authentication). Wyrażenie
zwraca **nazwę domeny** albo `false`.

### 2.3 Osobny selektor dla poczty transakcyjnej — ograniczenie

W opracowaniu technicznym napisałem wcześniej o „DKIM z osobnym selektorem
dla poczty transakcyjnej". **Stalwart 0.16 tego nie umożliwia w sposób,
w jaki zwykle się to rozumie** i trzeba to powiedzieć wprost:

- `dkimSignDomain` zwraca **domenę**, nie selektor,
- Stalwart stosuje **wszystkie** podpisy przypisane do zwróconej domeny
  jednocześnie; dodanie drugiego obiektu `DkimSignature` nie przełącza
  selektora, tylko dokłada drugi podpis do **każdej** wiadomości tej domeny.

Do wyboru są trzy drogi:

| Droga | Konsekwencja |
|---|---|
| **jeden selektor na domenę** (zalecana) | prościej; rozdział strumieni realizuje osobny obiekt `Domain`, jeżeli kiedyś zajdzie potrzeba |
| **poddomena nadawcza**, np. `noreply@poczta.danaco-group.pl` | pełny rozdział selektorów i reputacji, ale **zmienia adres nadawcy** widoczny dla Operatora — decyzja poza zakresem tego runbooka |
| **dwa podpisy na domenie** | obie sygnatury na każdej wiadomości; nie daje rozdziału strumieni, podnosi rozmiar nagłówków |

Rotacja klucza odbywa się przez pole `stage`
(`active` → `pending` → `retiring` → `retired`) wraz z `nextTransitionAt` —
i to jest właściwe narzędzie do wymiany klucza, nie drugi selektor.

**Ustalenie:** jeden selektor `transakcyjna` na domenie `danaco-group.pl`,
rotacja przez `stage`.

---

## 3. Skrzynka nieobsługiwana i autoresponder

### 3.1 Wybór rozwiązania

Adres ma nie przyjmować korespondencji, ale Operator ma dostać list 3
(`list-03-autoresponder-brak-skrzynki`). To wyklucza odrzucenie na etapie
koperty: po `reject` na `MtaStageRcpt` nadawca dostaje zwrotkę od własnego
serwera, a nasz list nigdy nie powstaje.

| Rozwiązanie | Skutek |
|---|---|
| `reject` na `MtaStageRcpt` | brak listu 3, zwrotka od serwera nadawcy, zero kosztu wysyłki |
| **przyjęcie + `vacation` + `discard`** | Operator dostaje list 3, wiadomość nie trafia do żadnej skrzynki |

Wybrane: **drugie**, bo list 3 jest częścią zamówienia.

### 3.2 Skrypt Sieve

Stalwart obsługuje Sieve z rozszerzeniem `vacation` (RFC 5230) wraz
z `:seconds` (RFC 6131), `reject`/`ereject` (RFC 5429), `discard` (RFC 5228)
i `duplicate` (RFC 7352).

```sieve
require ["vacation", "fileinto"];

# Jedna odpowiedź na adres nadawcy na dobę. Treść bierze się z listu 3
# złożonego przez rdzeń; poniżej wariant awaryjny, gdyby rdzeń nie działał.
vacation
  :days 1
  :subject "Adres noreply@danaco-group.pl nie przyjmuje korespondencji"
  :from "noreply@danaco-group.pl"
  "Adres noreply@danaco-group.pl nie prowadzi skrzynki odbiorczej.
W sprawach konta pisz na support@danaco-group.pl.";

discard;
```

**RFC 5230 sam pilnuje reguł bezpieczeństwa**, które opisałem w opracowaniu
technicznym jako wymagania: brak odpowiedzi przy pustym adresie zwrotnym,
przy `Precedence: bulk` oraz przy `Auto-Submitted` innym niż `no`.
Nie trzeba tego programować ręcznie.

Ustawienia interpretera skryptów użytkownika (singleton `SieveUserInterpreter`):

| Pole | Domyślnie | Uwaga |
|---|---|---|
| `defaultExpiryVacation` | `2592000000` ms (30 dni) | okres pamiętania adresu, gdy skrypt nie poda `:days` |
| `maxOutMessages` | `3` | licznik odpowiedzi wychodzących; obejmuje `vacation` |
| `defaultSubject` | `"Automated reply"` | nadpisujemy przez `:subject` |

Skrypty administratora to obiekty `SieveSystemScript` (Settings › Sieve ›
System Scripts) z polami `name`, `isActive`, `contents` — treść skryptu
wpisuje się w całości, wczytanie z pliku nie jest wspierane. Podpięcie pod
etap SMTP przez pole `script` na singletonie `MtaStageData` albo
`MtaStageRcpt`.

**Niepotwierdzone:** czy istnieje globalny hook uruchamiany przy dostarczeniu
inny niż skrypt aktywny konta. W dokumentacji skrypt dostarczania to skrypt
użytkownika — zakładamy konto techniczne `noreply` z aktywnym skryptem
powyższej treści.

---

## 4. Ograniczenia częstotliwości

Obiekty `MtaOutboundThrottle` (Settings › MTA › Rates & Quotas › Outbound
Rate Limits) i `MtaInboundThrottle`. Pole `rate` to `{count, period}`,
gdzie **`period` jest w milisekundach**. `key` to zbiór zmiennych grupujących
jako obiekt JSON; pusty `key` oznacza limit globalny.

Limit wysyłki na jeden adres odbiorcy — zabezpieczenie przed zalewem kodami
przy pętli w rdzeniu:

```json
{
  "enable": true,
  "description": "Poczta transakcyjna — limit na odbiorcę",
  "key": {"rcptDomain": true},
  "match": {"else": "true"},
  "rate": {"count": 20, "period": 3600000}
}
```

Klucze grupujące wychodzące: `remoteIp`, `localIp`, `mx`, `sender`,
`senderDomain`, `rcptDomain`. Przychodzące dodatkowo: `listener`,
`authenticatedAs`, `heloDomain`, `rcpt`.

Zachowanie przy przekroczeniu:

- **przychodzące** → `451 4.4.5 Rate limit exceeded, try again later.`
- **wychodzące** → wiadomość **zostaje w kolejce** do resetu okna; nie jest
  odrzucana. Dla listu z kodem oznacza to opóźnienie, nie utratę — ale kod
  ma ważność 60 minut, więc limit musi być ustawiony powyżej realnego ruchu.

---

## 5. Dziennik zdarzeń a kody uwierzytelniające

### 5.1 Co da się ustawić

Obiekt `Tracer` (Settings › Telemetry › Tracers), warianty `Log`, `Journal`,
`Stdout`, `OtelHttp`, `OtelGrpc`. Pola: `enable`, `level`
(`error` · `warn` · `info` · `debug` · `trace`), `events`, `eventsPolicy`
(`include` / `exclude`), `lossy`. Wariant `Log` dodatkowo: `path`, `prefix`,
`rotate` (`daily` domyślnie).

Poziom pojedynczego zdarzenia zmienia obiekt `EventTracingLevel`, np.
`{"event": "auth.success", "level": "debug"}`.

### 5.2 Czego ustawić się nie da

**Stalwart nie ma przełącznika „nie zapisuj treści ani nagłówków".** Jedyne
narzędzia to poziom i filtr zdarzeń.

Zbiór atrybutów zdarzeń w kodzie Stalwarta nie zawiera pozycji odpowiadającej
tematowi wiadomości — są `From`, `To`, `MessageId`, `RemoteIp`, `Size`,
`QueueId`. Temat zatem do dziennika nie trafia. **Ustalenie oparte na kodzie
źródłowym, nie na dokumentacji** — potraktuj je jako prawdopodobne, nie pewne,
i sprawdź na własnym dzienniku (rozdz. 6, kontrola 7).

### 5.3 Konsekwencja dla tematu z kodem

Tematy listów 1, 2, 4 i 5 zawierają kod. Dopóki temat nie trafia do dziennika,
sprzeczności nie ma. Ale każdy element trasy, który temat zapisze — filtr
antyspamowy, brama, kopia zapasowa kolejki, historia śladów w wydaniu
Enterprise — rozszerza powierzchnię wycieku.

Zabezpieczenia:

- historia śladów wyłączona: `TracingStore` = `{"@type": "Disabled"}`,
  albo krótkie `holdTracesFor` na obiekcie `DataRetention` (domyślnie `30d`),
- brak zewnętrznego filtru antyspamowego na trasie wychodzącej,
- kontrola z rozdz. 6.

Jeżeli którykolwiek warunek nie jest spełniony, **usuń kod z tematu**: zmiana
dotyczy mapy `subjects` w `rdzen/mail/szablon.go` i jest jednowierszowa dla
każdego listu.

### 5.4 Strona rdzenia

Pakiet `rdzen/mail` nie prowadzi dziennika samodzielnie. Warstwa wysyłkowa
zapisuje wyłącznie to, co zwraca `Message.LogFields()`: rodzaj listu, adres
odbiorcy, identyfikator wiadomości i rozmiar. Kod, temat i treść nie
występują tam z założenia, a test `TestLogFieldsBezTresci` tego pilnuje.

---

## 6. Kontrola po wdrożeniu

Wykonaj **przed** pierwszą wysyłką produkcyjną.

| # | Kontrola | Sposób | Wynik oczekiwany |
|---:|---|---|---|
| 1 | rekord SPF pojedynczy | `dig +short TXT danaco-group.pl` | dokładnie jeden ciąg `v=spf1` |
| 2 | rekord DKIM czytelny | `dig +short TXT transakcyjna._domainkey.danaco-group.pl` | `v=DKIM1;` z kluczem |
| 3 | rekord DMARC | `dig +short TXT _dmarc.danaco-group.pl` | polityka bieżącego etapu |
| 4 | podpis na wiadomości | wysyłka próbna na skrzynkę kontrolną, nagłówki wiadomości | `DKIM-Signature` z `s=transakcyjna`, `Authentication-Results` z `dkim=pass spf=pass dmarc=pass` |
| 5 | obie części obecne | źródło wiadomości | `multipart/related`, w nim `multipart/alternative` z `text/plain` **przed** `text/html` |
| 6 | znak widoczny | podgląd w kliencie z obrazami włączonymi | jeden znak, wariant zgodny z motywem |
| 7 | temat poza dziennikiem | wysyłka próbna, `grep` po temacie w dzienniku serwera | brak trafień |
| 8 | autoresponder odpowiada raz | dwie wiadomości pod rząd na `noreply@` | jedna odpowiedź |
| 9 | brak pętli odpowiedzi | wiadomość z `Auto-Submitted: auto-generated` na `noreply@` | brak odpowiedzi |
| 10 | wiadomość nieprzycięta | rozmiar wysłanej wiadomości | poniżej 102 KB |
| 11 | limit działa | wysyłka ponad próg | wiadomości w kolejce, nie odrzucone |
| 12 | ustawienia wczytane | `stalwart-cli create Action/ReloadSettings` | brak błędu |

Kontrola 4 działa też przez wysyłkę na zewnętrzną skrzynkę kontrolną
i odczytanie nagłówków `Authentication-Results` po stronie odbiorcy — to
jedyna kontrola oceniająca stan tak, jak widzi go odbiorca.

---

## 7. Postępowanie przy usterce

| Objaw | Pierwszy trop |
|---|---|
| listy trafiają do wiadomości niechcianych | brak części `text/plain`, `dkim=fail` w `Authentication-Results`, `p=reject` przy niedziałającym DKIM |
| `dkim=fail` przy poprawnym rekordzie | rozjazd klucza prywatnego z opublikowanym, `expire` minęło, zmiana ustawień bez `ReloadSettings` |
| autoresponder odpowiada w pętli | brak `discard` w skrypcie, `Auto-Submitted` nieustawiony w listach rdzenia |
| Operator nie dostaje kodu, brak śladu odrzucenia | limit wychodzący — wiadomość stoi w kolejce; sprawdź kolejkę, nie dziennik odrzuceń |
| dwa znaki jeden pod drugim | Outlook dla Windows i brak `mso-hide:all` — regresja w szablonie |
| w liście widać `{{code}}` | wysyłka z pominięciem pakietu `rdzen/mail`; `substitute` zwraca błąd zamiast wysłać |

---

## 8. Ustalenia niepotwierdzone

Zebrane w jednym miejscu, żeby nie ginęły w treści:

- czy temat wiadomości trafia do dziennika — ustalone z kodu źródłowego
  Stalwarta, nie z dokumentacji (rozdz. 5.2),
- czy istnieje globalny hook dostarczania inny niż skrypt aktywny konta
  (rozdz. 3.2),
- czy `dkimSignDomain` może zwrócić domenę inną niż domena nadawcy —
  dokumentacja pokazuje wyłącznie przykłady zwracające `sender_domain`
  (rozdz. 2.3).

Sprawdź je na własnej instalacji przed wdrożeniem produkcyjnym.

---

**Źródła:**

- [Stalwart — DKIM Signing](https://stalw.art/docs/mta/authentication/dkim/sign/)
- [Stalwart — Sieve](https://stalw.art/docs/sieve/)
- [Stalwart — Trusted interpreter](https://stalw.art/docs/sieve/interpreter/trusted/)
- [Stalwart — Untrusted interpreter](https://stalw.art/docs/sieve/interpreter/untrusted/)
- [Stalwart — Outbound rate limits](https://stalw.art/docs/mta/outbound/rate-limit/)
- [Stalwart — Tracing](https://stalw.art/docs/telemetry/tracing/)
- [Stalwart — Configuration](https://stalw.art/docs/configuration/)
- [Stalwart — Object encoding](https://stalw.art/docs/configuration/object-encoding/)
- [Stalwart — wydania](https://github.com/stalwartlabs/stalwart/releases)
