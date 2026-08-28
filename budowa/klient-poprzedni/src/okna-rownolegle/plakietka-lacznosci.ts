import './lacznosc.css';

import { elementIkony } from '../ikony/ikony';
import {
  czyLacznoscWToku,
  czyWidocznaLacznosc,
  napisLacznosci,
  podpowiedzLacznosci,
  wariantKropkiLacznosci,
  wariantPlakietkiLacznosci,
  type OdczytLacznosci,
  type PortPonawiania,
  type StanLacznosciOkna,
} from './lacznosc-okna';

/**
 * Plakietka łączności w nagłówku okna roboczego łączy kropkę z napisem, tak by odczyt nie zależał od barwy, i chowa się przy połączeniu bez kolejki, żeby nie zagłuszać czerwonej kropki wielu okien.
 */
export interface PlakietkaLacznosci {
  /** Element montowany w nagłówku gniazda. */
  element: HTMLElement;
  /** Przyjmuje odczyt transportu i przerysowuje plakietkę. */
  ustaw(odczyt: OdczytLacznosci): void;
  /** Stan, który plakietka pokazuje w tej chwili — jedna prawda do sprawdzenia. */
  pokazywany(): StanLacznosciOkna;
  /** Zatrzymuje odświeżanie odliczania; obowiązkowe przy zdejmowaniu gniazda. */
  rozlacz(): void;
}

/** Ustawienia plakietki łączności; każde pole ma wartość domyślną stosowaną, gdy gospodarz go nie poda. */
export interface OpcjePlakietkiLacznosci {
  /** Dojście do przebiegu ponowienia; pominięte znaczy, że transport nie wystawia numeru próby ani czasu. */
  ponowienie?: PortPonawiania;
  /** Krok odświeżania odliczania w milisekundach. */
  krokOdliczaniaMs?: number;
}

/** Odstęp odświeżania odliczania czasu do następnej próby — poniżej sekundy, żeby napis nie skakał o dwa. */
const KROK_ODLICZANIA_MS = 250;

export function utworzPlakietkeLacznosci(
  opcje: OpcjePlakietkiLacznosci = {},
): PlakietkaLacznosci {
  const port = opcje.ponowienie ?? null;
  const krok = opcje.krokOdliczaniaMs ?? KROK_ODLICZANIA_MS;

  const element = document.createElement('span');
  element.className = 'dn-plakietka dn-okna__lacznosc';
  element.setAttribute('role', 'status');
  element.hidden = true;

  const kropka = document.createElement('span');
  kropka.className = 'dn-kropka';

  // Spinner zastępuje kropkę w stanach w toku jak w pasku górnym, bo dwa ruchy spierałyby się o uwagę.
  const wskaznik = document.createElement('span');
  wskaznik.className = 'dn-spinner dn-okna__lacznosc-wskaznik';
  wskaznik.setAttribute('aria-hidden', 'true');
  wskaznik.hidden = true;

  const napis = document.createElement('span');
  napis.className = 'dn-okna__lacznosc-napis';

  element.append(kropka, wskaznik, napis);

  // Czynność ponowienia istnieje tylko z portem; bez portu nie ma przycisku, a podpowiedź nazywa powód.
  const ponow: HTMLButtonElement | null = port === null ? null : przyciskPonowienia(port);
  if (ponow !== null) element.append(ponow);

  let odczyt: OdczytLacznosci = { stan: 'rozlaczony', oczekujace: 0 };
  /** Czy plakietka dostała już odczyt; do pierwszego odczytu milczy, bo nie wiem to nie brak łącza. */
  let znany = false;
  let pokazywanyStan: StanLacznosciOkna = zloz(odczyt);
  let uchwyt: ReturnType<typeof setInterval> | null = null;
  let czynna = true;

  /** Pełny stan: odczyt transportu wzbogacony o przebieg ponowienia z portu. */
  function zloz(podstawa: OdczytLacznosci): StanLacznosciOkna {
    if (port === null || podstawa.stan !== 'ponawianie') {
      return { ...podstawa, numerProby: null, zaPonowieniemMs: null };
    }
    return {
      ...podstawa,
      numerProby: port.numerProby(),
      zaPonowieniemMs: port.zaPonowieniem(),
    };
  }

  function przerysuj(): void {
    const stan = zloz(odczyt);
    pokazywanyStan = stan;

    element.hidden = !znany || !czyWidocznaLacznosc(stan);
    element.className = `dn-plakietka dn-okna__lacznosc ${wariantPlakietkiLacznosci(stan.stan)}`;
    element.dataset.lacznosc = stan.stan;
    element.dataset.oczekujace = String(stan.oczekujace);
    element.title = podpowiedzLacznosci(stan);

    const wToku = czyLacznoscWToku(stan.stan);
    kropka.className = `dn-kropka ${wariantKropkiLacznosci(stan.stan)}`;
    kropka.hidden = wToku;
    wskaznik.hidden = !wToku;
    // Do pierwszego odczytu napis jest pusty, nie schowany, bo czytnik ogłasza treść po jej zmianie.
    napis.textContent = znany ? napisLacznosci(stan) : '';

    // Ponowić da się wyłącznie to, co nie jest połączone; przy połączeniu cała plakietka i tak znika.
    if (ponow !== null) ponow.hidden = stan.stan === 'polaczony' || stan.stan === 'laczenie';

    zarzadzajOdliczaniem(stan);
  }

  /** Odliczanie biegnie wyłącznie przy ponawianiu i z portem, bo bez portu nie ma czego odliczać. */
  function zarzadzajOdliczaniem(stan: StanLacznosciOkna): void {
    const potrzebne = czynna && znany && port !== null && stan.stan === 'ponawianie';
    if (potrzebne && uchwyt === null) {
      uchwyt = setInterval(przerysuj, krok);
      return;
    }
    if (!potrzebne) zatrzymaj();
  }

  function zatrzymaj(): void {
    if (uchwyt !== null) {
      clearInterval(uchwyt);
      uchwyt = null;
    }
  }

  przerysuj();

  return {
    element,

    ustaw(nowy) {
      odczyt = nowy;
      znany = true;
      przerysuj();
    },

    pokazywany: () => pokazywanyStan,

    rozlacz() {
      czynna = false;
      zatrzymaj();
    },
  };
}

/** Przycisk ponowienia natychmiastowego — skrót do przodu, nigdy rezygnacja z dalszego ponawiania połączenia. */
function przyciskPonowienia(port: PortPonawiania): HTMLButtonElement {
  const przycisk = document.createElement('button');
  przycisk.type = 'button';
  przycisk.className = 'dn-btn-ikona dn-okna__lacznosc-ponow';
  const opis = 'Ponów połączenie teraz — bez czekania na zaplanowany odstęp';
  przycisk.setAttribute('aria-label', opis);
  przycisk.title = opis;
  przycisk.append(elementIkony('odswiez', { rozmiar: 14 }));
  przycisk.addEventListener('click', (zdarzenie) => {
    // Plakietka stoi w nagłówku gniazda, a nagłówek bywa klikalny; ponowienie nie jest wyborem okna.
    zdarzenie.stopPropagation();
    port.ponowTeraz();
  });
  return przycisk;
}
