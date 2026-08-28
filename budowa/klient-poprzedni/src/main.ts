import './aplikacja/arkusze-stylow';

import { zlozAplikacje } from './aplikacja/aplikacja';
import { przyjmijAdresZPowloki } from './polaczenie/adres-rdzenia';
import { adresRdzeniaZPowloki } from './powloka/most-rdzenia';
import { utworzBanerAktualizacji } from './aktualizacja/indeks';
import { utworzScenaWejscia } from './ladowanie/indeks';
import { utworzEkranLogowania } from './uwierzytelnienie/indeks';

// Punkt wejścia interfejsu. Wyłącznie kompozycja: adres rdzenia z powłoki, warstwy
// wizualne, złożenie aplikacji i uruchomienie — najpierw widok, potem łączność.
// Pytanie o adres pada przed złożeniem, bo złożenie zakłada transport na adresie.
const zPowloki = await adresRdzeniaZPowloki();
if (zPowloki !== null) przyjmijAdresZPowloki(zPowloki);

// Kolejność jest wiążąca: router pokazuje trasę z adresu, zanim ruszy transport.

const aplikacja = zlozAplikacje();

aplikacja.router.uruchom();
aplikacja.rdzen.transport.polacz();

// Bramka Operatora — jedyna w całej aplikacji. Ekran logowania jest przesłoną,
// nie trasą: kładzie się nad złożoną aplikacją i schodzi po wejściu, zamiast
// zastępować widok i zmuszać router do drugiego stanu.
utworzEkranLogowania({
  kanal: aplikacja.rdzen.kanal,
  naWejscie: (sesja) => {
    utworzScenaWejscia({ gotowosc: aplikacja.stronaGlownaGotowa() }).uruchom();
    aplikacja.rdzen.uzgodnienie.zwiazSesjeBramki(sesja.token);
  },
}).uruchom();

// Baner aktualizacji montuje się po ekranie logowania, bo należy do okna pracy;
// nad przesłoną bramki zapraszałby do restartu kogoś, kto jeszcze nie wszedł.
utworzBanerAktualizacji().uruchom();
