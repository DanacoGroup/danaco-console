import './aplikacja/arkusze-stylow';

import { zlozAplikacje } from './aplikacja/aplikacja';
import { przyjmijAdresZPowloki } from './polaczenie/adres-rdzenia';
import { adresRdzeniaZPowloki } from './powloka/most-rdzenia';
import { utworzBanerAktualizacji } from './aktualizacja/indeks';
import { utworzScenaWejscia } from './ladowanie/indeks';
import { utworzEkranLogowania } from './uwierzytelnienie/indeks';

// Punkt wejścia interfejsu. Wyłącznie kompozycja: logika, typy i obsługiwacze
// mieszkają w modułach. Kolejno:
//   - zapytanie powłoki natywnej o adres rdzenia (`powloka/most-rdzenia`)
//   - wczytanie warstw wizualnych (`aplikacja/arkusze-stylow`)
//   - złożenie aplikacji (`aplikacja/aplikacja`)
//   - uruchomienie: widok, potem łączność
//
// Adres gniazda rdzenia ustala się z lokalizacji dokumentu
// (`polaczenie/adres-rdzenia.ts`), a w oknie powłoki natywnej z pakietem
// osadzonym wychodzi z tego adres `tauri.localhost`, pod którym nikt nie
// nasłuchuje. Powłoka zna adres prawdziwy i podaje go poleceniem
// `adres_rdzenia`. Pytanie musi paść przed złożeniem aplikacji, bo złożenie
// zakłada transport na adresie; zadane później zastałoby gniazdo już otwarte
// pod adresem ślepym. Poza powłoką natywną odpowiedzią jest brak i pytanie
// kosztuje jedno sprawdzenie środowiska.
const zPowloki = await adresRdzeniaZPowloki();
if (zPowloki !== null) przyjmijAdresZPowloki(zPowloki);

// Kolejność dwóch ostatnich wywołań ma znaczenie. Router pokazuje trasę
// zapisaną w adresie dokumentu, zanim ruszy transport — Operator widzi
// Centrum dowodzenia od pierwszej klatki, a nie puste tło czekające na rdzeń.
//
// Uzgodnienia nie rozpoczynamy tutaj. Robi to przepływ komunikatów w chwili,
// gdy transport zgłosi stan „połączony" — wywołanie stąd dałoby drugie
// powitanie i podwojenie historii otwarcia.

const aplikacja = zlozAplikacje();

aplikacja.router.uruchom();
aplikacja.rdzen.transport.polacz();

// Bramka Operatora — jedyna w całej aplikacji. Ekran logowania jest przesłoną,
// nie trasą: kładzie się nad złożoną aplikacją i schodzi po wejściu, zamiast
// zastępować widok i zmuszać router do drugiego stanu.
//
// Uruchomienie idzie po `polacz()`, bo rozpoznanie bramki wysyła komendy
// `auth.*` i ekran potrzebuje kanału, który już rusza; wcześniej pokazywałby
// formularz, którego nie ma jak wysłać.
//
// Ekran montuje się sam (`document.body.append` w `uruchom()`), dlatego punkt
// wejścia go nie osadza — inaczej przesłona wisiałaby w dwóch miejscach naraz.
//
// `naWejscie` wiąże świeżą sesję z żywym połączeniem: `connection.hello`
// przyjmuje token, a rdzeń zapamiętuje, które gniazdo należy do której sesji
// bramki. To powitanie nie dubluje powitania z uzgodnienia — tamto pada przy
// nawiązaniu, zanim Operator poda hasło, więc niesie token sesji poprzedniej
// albo żaden. Bez powtórzenia po wejściu rdzeń znałby gniazdo jako niezwiązane
// przez całe uruchomienie, a `auth.password.reset` rozłączałby tego, kto
// właśnie zmienił hasło. Powitanie jest czystym odczytem: nie zakłada sesji
// pracy ani okna.
//
// Scena wejścia staje w tej samej chwili, w której bramka schodzi. Kolejność
// jest wiążąca: scena idzie do dokumentu PRZED zdjęciem przesłony (bramka woła
// `naWejscie`, zanim usunie swój element), więc między jednym ekranem a drugim
// nie mignie strona główna z pustym jeszcze wykazem środowisk. Scena stoi
// warstwę niżej niż bramka (`ladowanie.css`) i schodzi, gdy pierwszy odczyt
// strony (`home.enter`) się domknie — albo po własnym kresie czekania, żeby
// rdzeń milczący nie zamienił jej w zasłonę nad produktem.
utworzEkranLogowania({
  kanal: aplikacja.rdzen.kanal,
  naWejscie: (sesja) => {
    utworzScenaWejscia({ gotowosc: aplikacja.stronaGlownaGotowa() }).uruchom();
    aplikacja.rdzen.uzgodnienie.zwiazSesjeBramki(sesja.token);
  },
}).uruchom();

// Baner aktualizacji montuje się po ekranie logowania, bo należy do okna pracy;
// pokazany nad przesłoną bramki zapraszałby do restartu kogoś, kto jeszcze nie
// wszedł. Sam pas pojawia się dopiero, gdy kanał wydań ma wersję nowszą niż
// zainstalowana, a przy pierwszym uruchomieniu pyta o to z opóźnieniem, żeby nie
// konkurować z uruchomieniem o łącze.
utworzBanerAktualizacji().uruchom();
