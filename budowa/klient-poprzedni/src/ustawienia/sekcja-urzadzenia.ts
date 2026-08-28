import type { Device, DeviceChangedEvent } from '../../../shared/contract';
import { elementIkony } from '../ikony/ikony';
import { opisOdmowyBledu } from '../komponenty/odmowa';
import { przycisk } from '../modele/kontrolki-formularza';
import type { Kanal } from '../protokol/kanal';
import type { SekcjaUstawien } from './sekcje';
import { utworzZrodloUrzadzen, type ZrodloUrzadzen } from './zrodlo-urzadzen';

/** Sekcja urządzeń w oknie ustawień pokazuje wykaz maszyn powiązanych z kontem i pozwala unieważnić token dostępu, odświeżając się na żywo zdarzeniem rdzenia. */
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
      // Pustka bez słowa wyglądałaby jak nieudany odczyt; wykaz bez urządzeń jest jednak stanem możliwym.
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
    // Atrybuty niosą stan do arkusza; treść stanu stoi obok w etykiecie, odczyt nie zależy od kreski.
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
      // Przycisku pewnej odmowy tu nie ma: token już zszedł, więc nie ma czego unieważnić.
      const brak = document.createElement('span');
      brak.className = 'du-urzadzenie__bez-czynnosci';
      brak.textContent = 'nie ma czego unieważnić';
      czynnosc.append(brak);
    }
    glowa.append(czynnosc);

    wiersz.append(glowa, szczegol(urzadzenie));

    if (urzadzenie.current) {
      // Ostrzeżenie stoi pod czynnością, nie przed nią: opisuje skutek, a nie pyta o zgodę.
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
    // Rozstrzyga odpowiedź rdzenia, nie samo powodzenie wywołania, bo rdzeń może tokenu nie zdjąć.
    if (!wynik.wynik.revoked) {
      powiedz(
        'Rdzeń przyjął wywołanie, ale tokenu nie unieważnił — dostęp został ' +
          'nietknięty.',
        false,
      );
      return;
    }
    // Wykazu tu nie przerysowujemy: pełny wykaz po zmianie niesie zdarzenie rdzenia.
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
    // Po udanym odczycie odpowiedzią jest sam wykaz — zdanie odczytano powtarzałoby to, co operator widzi.
    powiedz('', true);
  }

  // Zmiana wykonana gdziekolwiek dolatuje tutaj zdarzeniem niosącym komplet urządzeń do każdego gniazda.
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

/** Zdanie o skutku unieważnienia jest jedno na oba głosy, odpowiedź i zdarzenie, żeby operator nie zobaczył dwóch opisów jednej czynności. */
function zdanieUniewaznienia(urzadzenie: Device): string {
  return urzadzenie.current
    ? 'Dostęp odebrany temu urządzeniu: token unieważniony. Przy najbliższym ' +
        'uruchomieniu trzeba zalogować się ponownie.'
    : `Token unieważniony: ${nazwijUrzadzenie(urzadzenie)}. Urządzenie zaloguje ` +
        'się ponownie przy następnym uruchomieniu.';
}

/** Zdanie o zmianie przyniesionej zdarzeniem, bo pole urządzenia jest w kontrakcie opcjonalne i bywa puste przy zmianie hasła. */
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

/** Nazwa urządzenia do zdania i do wiersza; brak nazwy w kontrakcie zastępuje identyfikator, brzydszy, ale jednoznaczny. */
function nazwijUrzadzenie(urzadzenie: Device): string {
  const nazwa = (urzadzenie.name ?? '').trim();
  return nazwa === '' ? urzadzenie.deviceId : nazwa;
}

/** Wiersz szczegółu urządzenia niesie identyfikator i ostatnią widzianą aktywność, w jednym wspólnym elemencie. */
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

/** Ostatnia widziana aktywność słowami; brak znacznika czasu znaczy, że rdzeń jej nie odnotował, nigdy że urządzenie nie było czynne. */
function opiszAktywnosc(znacznik?: number): string {
  if (znacznik === undefined || !Number.isFinite(znacznik)) {
    return 'ostatnia aktywność nieodnotowana';
  }
  return `ostatnio widziane: ${new Date(znacznik).toLocaleString('pl')}`;
}

/** Plakietka stanu niesie ikonę i etykietę, nigdy samą barwę, bo sama barwa bywa nieczytelna dla operatora. */
function plakietka(odmiana: string, ikona: 'uzytkownik' | 'ptaszek' | 'klodka', tresc: string): HTMLElement {
  const element = document.createElement('span');
  element.className = `dn-plakietka ${odmiana}`;
  element.append(elementIkony(ikona, { rozmiar: 12 }), znak(tresc));
  return element;
}

/** Tekst w elemencie własnym, bo ikona i treść muszą stać obok siebie, nie w jednym wspólnym węźle DOM. */
function znak(tresc: string): HTMLElement {
  const element = document.createElement('span');
  element.textContent = tresc;
  return element;
}
