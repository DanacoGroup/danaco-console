import type { Device, DeviceChangedEvent } from '../../../shared/contract';
import { elementIkony } from '../ikony/ikony';
import { opisOdmowyBledu } from '../komponenty/odmowa';
import { przycisk } from '../modele/kontrolki-formularza';
import type { Kanal } from '../protokol/kanal';
import type { SekcjaUstawien } from './sekcje';
import { utworzZrodloUrzadzen, type ZrodloUrzadzen } from './zrodlo-urzadzen';

/**
 * Sekcja „Urządzenia" Okna Ustawień — wykaz maszyn powiązanych z kontem
 * i unieważnienie tokenu dostępu.
 *
 * Konto jest jedno, urządzeń dowolnie wiele; każde niesie własny token. Wykaz
 * jest więc jedynym miejscem, w którym Operator widzi, co ma dostęp do
 * platformy, i jedynym, z którego może ten dostęp odebrać.
 *
 * Wykaz wchodzi komendą `device.list` przy wejściu do sekcji, a dalej odświeża
 * się zdarzeniem `device.changed`, które rdzeń rozgłasza do wszystkich
 * połączonych urządzeń. Odpytywania w tle nie ma i nie jest potrzebne:
 * unieważnienie wykonane na drugiej maszynie dolatuje tu samo. Odpowiedź na
 * `device.revoke` niesie samo `revoked`, więc to zdarzenie — nie odpowiedź —
 * przerysowuje wiersze.
 *
 * Wiersz urządzenia bieżącego (`current`) jest oznaczony plakietką i kreską,
 * a przy jego czynności stoi zdanie o skutku. Rdzeń celowo nie blokuje
 * unieważnienia własnego tokenu, więc ostrzeżenie należy do ekranu — ale jako
 * opis obok przycisku, nie jako przesłona przed nim. Okna „czy na pewno" tu nie
 * ma z zasady produktu: czynność jest odwracalna ponownym zalogowaniem, a
 * bramka przed nią uczyłaby wyłącznie odklikiwania.
 *
 * Urządzenie bez ważnego tokenu (`hasToken`) nie dostaje przycisku, tylko
 * zdanie: nie ma czego unieważnić, a przycisk pewnej odmowy byłby gorszy od
 * jego braku — ten sam zamysł, co przy kotwicy bramki w sekcji
 * „Uwierzytelnianie".
 *
 * Stanu nie niesie tu sama barwa: przy każdej plakietce stoi ikona i etykieta,
 * bo kreska w barwie sygnału znika dla Operatora, który barw nie rozróżnia.
 */
export function utworzSekcjeUrzadzenia(kanal: Kanal): SekcjaUstawien {
  const zrodlo: ZrodloUrzadzen = utworzZrodloUrzadzen(kanal);

  const element = document.createElement('div');
  element.className = 'du-sekcja';

  const wykaz = document.createElement('ul');
  wykaz.className = 'du-urzadzenia';

  const odpowiedz = document.createElement('p');
  odpowiedz.className = 'du-odpowiedz';
  odpowiedz.hidden = true;

  const oWykazie = document.createElement('p');
  oWykazie.className = 'dn-pole-opis du-granica';
  oWykazie.textContent =
    'Wykaz wchodzi komendą device.list przy wejściu do sekcji, a potem odświeża ' +
    'się sam zdarzeniem device.changed — także po unieważnieniu wykonanym na ' +
    'innym urządzeniu. Sekcja nie odpytuje rdzenia w tle. Unieważniony token ' +
    'znaczy jedno: przy najbliższym uruchomieniu maszyna loguje się ponownie.';

  element.append(wykaz, odpowiedz, oWykazie);

  function powiedz(tresc: string, powodzenie: boolean): void {
    odpowiedz.textContent = tresc;
    odpowiedz.hidden = tresc === '';
    odpowiedz.dataset['powodzenie'] = String(powodzenie);
  }

  function pokazUrzadzenia(urzadzenia: readonly Device[]): void {
    if (urzadzenia.length === 0) {
      // Pustka bez słowa wyglądałaby jak odczyt, który się nie udał. Wykaz bez
      // ani jednego urządzenia jest jednak stanem możliwym i mówimy to wprost.
      const pusto = document.createElement('li');
      pusto.className = 'du-urzadzenia__pusto';
      pusto.textContent =
        'Rdzeń nie zna ani jednego urządzenia powiązanego z kontem. Powiązanie ' +
        'powstaje przy logowaniu, więc wykaz zapełni się po pierwszym wejściu.';
      wykaz.replaceChildren(pusto);
      return;
    }
    wykaz.replaceChildren(...urzadzenia.map(zbudujWiersz));
  }

  function zbudujWiersz(urzadzenie: Device): HTMLElement {
    const wiersz = document.createElement('li');
    wiersz.className = 'du-urzadzenie';
    // Atrybuty niosą stan do arkusza; treść stanu i tak stoi obok w etykiecie,
    // więc odczyt nie zależy od tego, czy ktoś widzi kreskę.
    wiersz.dataset['biezace'] = String(urzadzenie.current);
    wiersz.dataset['token'] = String(urzadzenie.hasToken);

    const glowa = document.createElement('div');
    glowa.className = 'du-urzadzenie__glowa';

    const nazwa = document.createElement('span');
    nazwa.className = 'du-urzadzenie__nazwa';
    nazwa.textContent = nazwijUrzadzenie(urzadzenie);

    glowa.append(elementIkony('monitor', { rozmiar: 16 }), nazwa);

    if (urzadzenie.current) {
      glowa.append(
        plakietka('dn-plakietka--sygnal', 'uzytkownik', 'to urządzenie'),
      );
    }

    glowa.append(
      urzadzenie.hasToken
        ? plakietka('dn-plakietka--sukces', 'ptaszek', 'dostęp czynny')
        : plakietka('dn-plakietka--ostrzezenie', 'klodka', 'token unieważniony'),
    );

    const czynnosc = document.createElement('div');
    czynnosc.className = 'du-urzadzenie__czynnosc';
    if (urzadzenie.hasToken) {
      const uniewaznij = przycisk(
        'Unieważnij dostęp',
        'dn-btn dn-btn--sm dn-btn--niebezpieczny',
      );
      uniewaznij.addEventListener('click', () => void uniewaznijDostep(urzadzenie));
      czynnosc.append(uniewaznij);
    } else {
      // Przycisku pewnej odmowy tu nie ma: token już zszedł, więc nie ma czego
      // unieważnić. Zamiast kontrolki wyłączonej stoi zdanie, dlaczego jej nie ma.
      const brak = document.createElement('span');
      brak.className = 'du-urzadzenie__bez-czynnosci';
      brak.textContent = 'nie ma czego unieważnić';
      czynnosc.append(brak);
    }
    glowa.append(czynnosc);

    wiersz.append(glowa, szczegol(urzadzenie));

    if (urzadzenie.current) {
      // Ostrzeżenie stoi pod czynnością, nie przed nią: opisuje skutek, a nie
      // pyta o zgodę. Rdzeń tej czynności nie blokuje i ekran też jej nie broni.
      const ostrzezenie = document.createElement('p');
      ostrzezenie.className = 'du-urzadzenie__ostrzezenie';
      ostrzezenie.append(
        elementIkony('ostrzezenie', { rozmiar: 14 }),
        znak(
          'To maszyna, przy której Operator właśnie pracuje. Unieważnienie ' +
            'odbiera dostęp jej samej, w tej samej chwili — rdzeń tej czynności ' +
            'nie wstrzymuje. Powrót wymaga ponownego zalogowania.',
        ),
      );
      wiersz.append(ostrzezenie);
    }

    return wiersz;
  }

  async function uniewaznijDostep(urzadzenie: Device): Promise<void> {
    powiedz(`Unieważnianie tokenu: ${nazwijUrzadzenie(urzadzenie)}…`, true);
    const wynik = await zrodlo.uniewaznij(urzadzenie.deviceId);
    if (!wynik.udany || wynik.wynik === undefined) {
      powiedz(opisOdmowyBledu('Unieważnienie dostępu', wynik.blad), false);
      return;
    }
    // Rozstrzyga odpowiedź rdzenia, nie samo powodzenie wywołania: rdzeń może
    // przyjąć wywołanie i tokenu nie zdjąć, a wtedy „unieważniono" byłoby
    // potwierdzeniem czynności, która się nie odbyła.
    if (!wynik.wynik.revoked) {
      powiedz(
        'Rdzeń przyjął wywołanie, ale tokenu nie unieważnił — dostęp został ' +
          'nietknięty.',
        false,
      );
      return;
    }
    // Wykazu tu nie przerysowujemy: pełny wykaz po zmianie niesie
    // `device.changed`, a druga droga do tych samych wierszy rozjeżdżałaby się
    // z pierwszą przy każdej zmianie kontraktu.
    powiedz(zdanieUniewaznienia(urzadzenie), true);
  }

  async function odczytajWykaz(): Promise<void> {
    powiedz('Odczyt wykazu urządzeń…', true);
    const wynik = await zrodlo.wykaz();
    if (!wynik.udany || wynik.wynik === undefined) {
      powiedz(opisOdmowyBledu('Odczyt wykazu urządzeń', wynik.blad), false);
      return;
    }
    pokazUrzadzenia(wynik.wynik.devices);
    // Po udanym odczycie odpowiedzią jest sam wykaz — zdanie „odczytano"
    // powtarzałoby to, co Operator ma przed oczami.
    powiedz('', true);
  }

  /**
   * Zmiana wykonana gdziekolwiek dolatuje tutaj: `device.changed` niesie komplet
   * urządzeń do każdego gniazda, więc sekcja przerysowuje wykaz bez pytania
   * rdzenia.
   */
  const odsubskrybuj = zrodlo.naZmianeUrzadzen((zmiana) => {
    pokazUrzadzenia(zmiana.devices);
    const zdanie = zdanieZmiany(zmiana);
    if (zdanie !== '') powiedz(zdanie, true);
  });

  return {
    element,
    odswiez: () => void odczytajWykaz(),
    rozlacz: odsubskrybuj,
  };
}

/**
 * Zdanie o skutku unieważnienia — jedno na oba głosy, odpowiedź i zdarzenie.
 *
 * Oba potrafią dojść w dowolnej kolejności, więc gdyby każdy miał własne
 * brzmienie, Operator zobaczyłby dwa różne opisy jednej czynności.
 */
function zdanieUniewaznienia(urzadzenie: Device): string {
  return urzadzenie.current
    ? 'Dostęp odebrany temu urządzeniu: token unieważniony. Przy najbliższym ' +
        'uruchomieniu trzeba zalogować się ponownie.'
    : `Token unieważniony: ${nazwijUrzadzenie(urzadzenie)}. Urządzenie zaloguje ` +
        'się ponownie przy następnym uruchomieniu.';
}

/**
 * Zdanie o zmianie przyniesionej zdarzeniem. Bez niego zmiana wykonana na innej
 * maszynie byłaby niema: wykaz podmieniłby się pod ręką Operatora bez słowa
 * o tym, co zaszło.
 *
 * Pole `deviceId` jest w kontrakcie opcjonalne — tak przychodzi zmiana hasła,
 * która unieważnia tokeny hurtem. Wtedy zdanie nazywa zmianę ogólnie, zamiast
 * zgadywać za rdzeń, którego urządzenia dotyczyła.
 */
function zdanieZmiany(zmiana: DeviceChangedEvent): string {
  const wskazane = (zmiana.deviceId ?? '').trim();
  if (wskazane === '') {
    return 'Rdzeń zgłosił zmianę wykazu urządzeń — wykaz odświeżony zdarzeniem device.changed.';
  }
  const dotkniete = zmiana.devices.find((urzadzenie) => urzadzenie.deviceId === wskazane);
  if (dotkniete === undefined) {
    return `Rdzeń zgłosił zmianę urządzenia ${wskazane}, ale nie ma go w wykazie po zmianie — powiązanie z kontem zniknęło.`;
  }
  if (!dotkniete.hasToken) return zdanieUniewaznienia(dotkniete);
  return `Rdzeń zgłosił zmianę urządzenia ${nazwijUrzadzenie(dotkniete)} — wykaz odświeżony zdarzeniem device.changed.`;
}

/**
 * Nazwa urządzenia do zdania i do wiersza.
 *
 * `name` jest w kontrakcie opcjonalne. Zamiast pustego miejsca idzie wtedy
 * identyfikator: brzydszy, ale jednoznaczny — a wiersz bez nazwy nie daje się
 * odróżnić od sąsiedniego.
 */
function nazwijUrzadzenie(urzadzenie: Device): string {
  const nazwa = (urzadzenie.name ?? '').trim();
  return nazwa === '' ? urzadzenie.deviceId : nazwa;
}

/** Wiersz szczegółu: identyfikator i ostatnia widziana aktywność. */
function szczegol(urzadzenie: Device): HTMLElement {
  const element = document.createElement('p');
  element.className = 'du-urzadzenie__szczegol';

  const identyfikator = document.createElement('span');
  identyfikator.className = 'du-urzadzenie__identyfikator';
  identyfikator.textContent = urzadzenie.deviceId;

  const widziane = document.createElement('span');
  widziane.textContent = opiszAktywnosc(urzadzenie.lastSeenAt);

  element.append(identyfikator, widziane);
  return element;
}

/**
 * Ostatnia widziana aktywność słowami.
 *
 * `lastSeenAt` jest opcjonalne, a jego brak znaczy „rdzeń tego nie odnotował" —
 * nie „urządzenie nigdy nie było czynne". Podstawienie w to miejsce daty
 * zerowej czytałoby się jako aktywność w 1970 roku.
 */
function opiszAktywnosc(znacznik?: number): string {
  if (znacznik === undefined || !Number.isFinite(znacznik)) {
    return 'ostatnia aktywność nieodnotowana';
  }
  return `ostatnio widziane: ${new Date(znacznik).toLocaleString('pl')}`;
}

/** Plakietka stanu: ikona i etykieta, nigdy sama barwa. */
function plakietka(odmiana: string, ikona: 'uzytkownik' | 'ptaszek' | 'klodka', tresc: string): HTMLElement {
  const element = document.createElement('span');
  element.className = `dn-plakietka ${odmiana}`;
  element.append(elementIkony(ikona, { rozmiar: 12 }), znak(tresc));
  return element;
}

/** Tekst w elemencie własnym — ikona i treść muszą stać obok siebie, nie w jednym węźle. */
function znak(tresc: string): HTMLElement {
  const element = document.createElement('span');
  element.textContent = tresc;
  return element;
}
