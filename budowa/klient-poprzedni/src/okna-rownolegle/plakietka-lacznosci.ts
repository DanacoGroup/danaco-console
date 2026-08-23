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
 * Plakietka łączności w nagłówku okna roboczego.
 *
 * Wygląd idzie rodziną klas `.dn-kropka--*` i `.dn-plakietka--*`, a obok kropki
 * stoi słowo — odczyt nie zależy od rozróżnienia barw. Przy połączeniu bez
 * kolejki plakietka jest schowana, żeby zielone kropki wielu okien nie
 * zagłuszyły jednej czerwonej.
 *
 * Numer próby i czas do następnej należą do transportu; plakietka odczytuje je
 * w chwili rysowania i odświeża, dopóki trwa ponawianie. Własny zegar
 * odliczający rozjechałby się z transportem przy próbie podjętej wcześniej.
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

/** Ustawienia plakietki; każde ma wartość domyślną. */
export interface OpcjePlakietkiLacznosci {
  /**
   * Dojście do przebiegu ponowienia.
   *
   * Pominięte znaczy „transport nie wystawia numeru próby ani czasu do
   * następnej" — plakietka mówi to wtedy wprost w podpowiedzi i nie pokazuje
   * czynności „Ponów teraz", bo nie miałaby czego wywołać.
   */
  ponowienie?: PortPonawiania;
  /** Krok odświeżania odliczania w milisekundach. */
  krokOdliczaniaMs?: number;
}

/** Odstęp odświeżania odliczania — poniżej sekundy, żeby napis nie skakał o dwa. */
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

  // Spinner zastępuje kropkę w stanach w toku — tak samo jak w pasku górnym
  // (`aplikacja/wskaznik-lacznosci.ts`). Dwa ruchy naraz w plakietce wielkości
  // pigułki spierałyby się o uwagę, a znaczenie niesie i tak napis.
  const wskaznik = document.createElement('span');
  wskaznik.className = 'dn-spinner dn-okna__lacznosc-wskaznik';
  wskaznik.setAttribute('aria-hidden', 'true');
  wskaznik.hidden = true;

  const napis = document.createElement('span');
  napis.className = 'dn-okna__lacznosc-napis';

  element.append(kropka, wskaznik, napis);

  // Czynność „Ponów teraz" istnieje tylko z portem: przycisk bez wywołania za
  // nim byłby atrapą. Nie ma portu — nie ma przycisku, a podpowiedź nazywa powód.
  const ponow: HTMLButtonElement | null = port === null ? null : przyciskPonowienia(port);
  if (ponow !== null) element.append(ponow);

  let odczyt: OdczytLacznosci = { stan: 'rozlaczony', oczekujace: 0 };
  /**
   * Czy plakietka dostała już odczyt. Do pierwszego odczytu milczy, zamiast
   * pokazywać „Rozłączony": scena bez transportu (podgląd układu) łącza nie
   * zna, a „nie wiem" i „nie ma łącza" to dwie różne rzeczy.
   */
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
    // Do pierwszego odczytu napis jest pusty, nie schowany: element z rolą
    // `status` czytnik ekranu ogłasza po treści, a treść „Rozłączony" ukryta
    // atrybutem i tak bywa zapowiadana przy najbliższej zmianie.
    napis.textContent = znany ? napisLacznosci(stan) : '';

    // Ponowić da się wyłącznie to, co nie jest połączone. Przy połączeniu cała
    // plakietka i tak znika, więc przycisk nie zostaje sam na scenie.
    if (ponow !== null) ponow.hidden = stan.stan === 'polaczony' || stan.stan === 'laczenie';

    zarzadzajOdliczaniem(stan);
  }

  /**
   * Odliczanie biegnie wyłącznie przy ponawianiu i wyłącznie z portem: bez
   * portu nie ma czego odliczać, a przy każdym innym stanie liczba się nie
   * zmienia i pętla byłaby pracą bez odbiorcy.
   */
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

/** Przycisk „Ponów teraz" — skrót do przodu, nigdy rezygnacja z ponawiania. */
function przyciskPonowienia(port: PortPonawiania): HTMLButtonElement {
  const przycisk = document.createElement('button');
  przycisk.type = 'button';
  przycisk.className = 'dn-btn-ikona dn-okna__lacznosc-ponow';
  const opis = 'Ponów połączenie teraz — bez czekania na zaplanowany odstęp';
  przycisk.setAttribute('aria-label', opis);
  przycisk.title = opis;
  przycisk.append(elementIkony('odswiez', { rozmiar: 14 }));
  przycisk.addEventListener('click', (zdarzenie) => {
    // Plakietka stoi w nagłówku gniazda, a nagłówek bywa klikalny w całości —
    // ponowienie nie jest wyborem okna.
    zdarzenie.stopPropagation();
    port.ponowTeraz();
  });
  return przycisk;
}
