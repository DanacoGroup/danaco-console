# Plan rozbudowy poczty transakcyjnej

Dostawa `design/06-poczta-transakcyjna` niesie **siedem** listów. Rdzeń nadaje
dziś **cztery**: aktywację konta, kod logowania, reset hasła i autoresponder.
Trzy pozostałe mają szablony wpięte w `mail.Load`, ale nie mają przebiegu,
który by je nadał.

[Rozstrzygnięcie Właściciela 25](decyzje.md) stanowi, że przebiegów nie buduje
się teraz: pierwszeństwo ma działająca aplikacja. Ten dokument opisuje zakres,
żeby przy podjęciu pracy nie trzeba było rozpoznawać go od nowa.

Szablony zostają w rdzeniu celowo. Wpięte w `mail.Load` sprawiają, że rozjazd
między szablonem a kodem ujawnia się przy budowie — wyjęte, rozjechałyby się
w ciszy.

---

## Listy 5 i 6 — zmiana adresu konta — WYKONANE 6.09.2026

Przebieg stoi we wszystkich warstwach; rewizja `9e99ee57`.

| Warstwa | Co powstało |
|---|---|
| kontrakt | `auth.email.change.start`, `.confirm`, `.revoke`; powód `emailChanged` |
| rdzeń | `adapter_modul_auth_zmiana_adresu.go`, trzy metody portu, rejestracja komend |
| trwałość | krok 498: tabela `zmiana_adresu_konta`; `UstawAdresKonta` w repozytorium konta |
| trasa | `konfiguracja.AdresWycofaniaZmiany` składa `{{revoke_url}}` z drogą w zapytaniu |
| wysyłka | para listów jest jedną czynnością: niepowodzenie nadania zamyka zamówienie |
| klient | `konto-zmiana-adresu.ts`; wycofanie idzie z zapytania adresu, bez kliknięcia |

**Rozstrzygnięte przy pracy.** `{{revoke_expiry_hours}}` nie miało źródła —
przyjęto **72 godziny**: wycofanie ma przeżyć weekend, bo list na adres
dotychczasowy jest jedyną obroną Operatora, któremu przejęto konto. Kod zostaje
przy godzinie wspólnej wszystkim drogom potwierdzenia.

**Zostaje otwarte.** `{{ip_address}}` listu 6 pokazuje myślnik: rdzeń nie zna
adresu źródłowego żądania — gniazdo WebSocket kończy się na warstwie transportu,
a kontekst żądania adresu nie niesie. Urządzenie idzie z sesji bramki.
Okna Ustawienia konta w kliencie nadal nie ma; wiązanie czeka na elementy
`data-adres-nowy`, `data-adres-haslo`, `data-adres-kod`, `data-adres-zamow`
i `data-adres-potwierdz`.

## Listy 5 i 6 — zakres rozpoznany przed pracą

Oba wychodzą **razem, na dwa różne adresy**: kod na adres nowy, ostrzeżenie
z drogą wycofania na dotychczasowy. Para jest wymaganiem bezpieczeństwa, nie
układu graficznego — bez listu 6 przejęte konto zmienia adres niepostrzeżenie.

Przebiegu nie ma w żadnej warstwie. Brakuje:

| Warstwa | Czego brakuje |
|---|---|
| kontrakt | komendy wszczynającej zmianę i domykającej ją kodem; struktur żądania i wyniku; wartości `AuthChangeReason` znaczącej zmianę adresu |
| rdzeń | obsługiwacza (`handlers_auth.go` rejestruje dziewięć komend, żadnej o zmianie adresu); metody portu `Uwierzytelnianie` i jej wypełnienia |
| trwałość | migracji rozszerzającej `CHECK` na `potwierdzenie_tozsamosci.cel` — dziś dopuszcza wyłącznie `weryfikacja` i `odzyskanie`; miejsca na adres oczekujący; skrótu tokenu wycofania wraz z terminem |
| trasa | obsługi `{{revoke_url}}` — serwer wystawia dziś gniazdo WebSocket i statykę, nie ma czego zawołać odsyłaczem z listu |
| wysyłka | złożenia **pary** listów w jednej czynności; `nadajnik.List` ma jedno pole odbiorcy |
| klient | okna **Ustawienia konta**, do którego odsyła treść obu listów — w kliencie nie ma go wcale |

**Do rozstrzygnięcia przed pracą:** `{{revoke_expiry_hours}}` nie ma źródła.
Rdzeń zna jeden termin — `trwanieDrogiPotwierdzenia`, godzinę, wspólny dla
weryfikacji i odzyskania. Wycofanie zmiany adresu z założenia ma trwać dłużej
niż kod; jak długo, nie mówi żadne źródło.

---

## List 7 — przebieg automatyki zakończony

Punkty zakończenia przebiegu **istnieją, dwa i niezależne**: wyprowadzenie
stanu z kolejki (`adapter_modul_automations_przebiegi.go`) oraz zatrzymanie
biegu po upływie terminu w silniku wybudzeń (`handlers_automatyka_petla.go`).

Brakuje przede wszystkim **adresata**. Wiersz przebiegu nie niesie wskazania
konta ani osoby; łańcuch przebieg → kolejka → sesja → karta → konto jest
w schemacie możliwy, ale kolejka przebiegu powstaje bez wskazania sesji,
a `karta_sesji.konto_id` dopuszcza pustkę. Wskazanie adresata bez decyzji
Właściciela byłoby dopowiedzeniem.

Poza adresatem brakuje źródła dla siedmiu z ośmiu zmiennych listu:

| Zmienna | Czego brakuje |
|---|---|
| `{{environment_name}}` | żaden byt przebiegu nie wskazuje środowiska; migracja 039 stanowi wprost, że automatyka nie jest bytem sesji ani okna |
| `{{run_url}}` | w rdzeniu nie ma adresu bazowego konsoli ani schematu odsyłacza do okna przebiegu — ta sama dziura zostawia dziś pusty `{{activation_url}}` w liście 1 |
| `{{run_status}}` | rdzeń ma trzy stany końcowe (`succeeded`, `failed`, `stopped`), dostawa dopuszcza dwa (`zakończony`, `zatrzymany`) |
| `{{run_duration}}` | brak formatu w rdzeniu, a na drodze przez silnik wybudzeń czas zakończenia w ogóle nie jest zapisywany |
| `{{run_step}}` | postać `3 z 5` nie jest w dostawie opisana jako reguła składania; rdzeń ma same liczby |
| `{{started_at}}` | rdzeń zapisuje czas w UTC i nie zna strefy Operatora, a przykład dostawy niesie skrót strefy |
| `{{run_name}}` | nazwa stoi w tabeli automatyki, ale odczyt przebiegu pobiera sam identyfikator zewnętrzny |

Brakuje też wyzwalacza i kanału: nikt nie ocenia reguł alarmowania przy
zakończeniu przebiegu, wykaz wyzwalaczy kontraktu nie zna tego zdarzenia,
a centrum powiadomień nie ma kanału pocztowego — `kanal_dostarczenia`
dopuszcza `centrum` i `centrum_i_push`. Nie ma również nastawy częstotliwości,
którą stopka listu 7 obiecuje w oknie Ustawienia konta.

**Do rozstrzygnięcia przed pracą:** jak brzmi list dla przebiegu **nieudanego**.
Dostawa dopuszcza dwie wartości stanu, rdzeń ma trzy, a przełożenia trzech na
dwa nie podaje żadne źródło.

---

## Usterka rdzenia znaleziona przy rozpoznaniu — naprawiona

Droga przez silnik wybudzeń ustawiała stan przebiegu na zatrzymany, nie ustawiając
czasu zakończenia. `przerwij` w `handlers_automatyka_petla.go` stempluje go teraz
tak samo jak droga przez kolejkę, więc czas trwania da się policzyć na obu.
