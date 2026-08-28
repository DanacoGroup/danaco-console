import type { Transport } from '../polaczenie/gniazdo';
import type { StanPolaczenia } from '../polaczenie/stan-polaczenia';
import { stanRdzeniaZPowloki, type StanRdzenia } from '../powloka/most-rdzenia';

/** Wskaźnik łączności zamontowany w grupie akcji paska: plakietka montowana wraz ze sposobem jej odłączenia. */
export interface WskaznikLacznosci {
  /** Plakietka montowana na pasku. */
  element: HTMLElement;
  /** Odłącza subskrypcję transportu. */
  rozlacz(): void;
}

/** Nazwa stanu łączności widoczna dla Operatora, osobna dla każdego z czterech stanów samego transportu. */
const NAZWY: Readonly<Record<StanPolaczenia, string>> = {
  rozlaczony: 'Rozłączony',
  laczenie: 'Łączenie',
  polaczony: 'Połączony',
  ponawianie: 'Ponawianie',
};

/** Odmiana kropki stanu ze słownika `komponenty/plakietka.css`; wskaźnik sam nie niesie żadnych barw własnych. */
const KROPKI: Readonly<Record<StanPolaczenia, string>> = {
  rozlaczony: 'dn-kropka--blad',
  laczenie: 'dn-kropka--tetno',
  polaczony: 'dn-kropka--sukces',
  ponawianie: 'dn-kropka--ostrzezenie',
};

/** Stany, w których łączność jest dopiero w toku — nośnikiem jest wtedy `.dn-spinner`, nie kropka stanu. */
const W_TOKU: ReadonlySet<StanPolaczenia> = new Set<StanPolaczenia>(['laczenie', 'ponawianie']);

/** Odstęp dobijania licznika kolejki po połączeniu, wyrażony w milisekundach między kolejnymi odczytami. */
const KROK_ODSWIEZANIA_MS = 150;

/** Stan łączności z rdzeniem pokazany na pasku: przełożenie stanu transportu na plakietkę informacyjną. */
export function utworzWskaznikLacznosci(transport: Transport): WskaznikLacznosci {
  const element = document.createElement('span');
  // Plakietka neutralna — stan dopowiadają kropka i etykieta obok, nigdy sama barwa tła.
  element.className = 'dn-plakietka dn-lacznosc';

  const kropka = document.createElement('span');
  kropka.className = 'dn-kropka';

  const wskaznik = document.createElement('span');
  wskaznik.className = 'dn-spinner dn-lacznosc__wskaznik';
  wskaznik.setAttribute('aria-hidden', 'true');
  wskaznik.hidden = true;

  // Napis stanu nie nosi klasy — całe jego zachowanie niesie już plakietka.
  const opis = document.createElement('span');

  element.append(kropka, wskaznik, opis);

  // Zdanie powłoki natywnej o rdzeniu; puste, dopóki powłoka nie odpowie.
  let diagnoza = '';
  // Pytanie w locie — jedno naraz wystarczy, żeby nie mnożyć wywołań IPC.
  let pytanie = false;
  // Plakietka zdjęta z paska nie przyjmuje spóźnionej odpowiedzi.
  let czynna = true;
  let biezacy: StanPolaczenia = transport.stan();
  // Uchwyt pętli dobijającej licznik kolejki; null, gdy pętla nie biegnie.
  let uchwytOdswiezania: ReturnType<typeof setInterval> | null = null;

  function opiszPodpowiedz(): void {
    const podstawa = `Łączność z rdzeniem: ${opis.textContent}`;
    element.title = diagnoza === '' ? podstawa : `${podstawa}\n${diagnoza}`;
  }

  /** Przerysowuje etykietę i podpowiedź na podstawie bieżącej liczby w kolejce. */
  function odswiezEtykiete(stan: StanPolaczenia): void {
    const oczekujace = transport.oczekujace();
    opis.textContent =
      oczekujace > 0 ? `${NAZWY[stan]} · w kolejce ${oczekujace}` : NAZWY[stan];
    opiszPodpowiedz();
  }

  function zatrzymajOdswiezanie(): void {
    if (uchwytOdswiezania !== null) {
      clearInterval(uchwytOdswiezania);
      uchwytOdswiezania = null;
    }
  }

  // Po połączeniu kolejka opróżnia się po ogłoszeniu stanu — licznik dobijany jest do zera.
  function zarzadzajOdswiezaniem(stan: StanPolaczenia): void {
    if (stan === 'polaczony' && transport.oczekujace() > 0) {
      if (uchwytOdswiezania === null) {
        uchwytOdswiezania = setInterval(() => {
          odswiezEtykiete(biezacy);
          if (transport.oczekujace() === 0) zatrzymajOdswiezanie();
        }, KROK_ODSWIEZANIA_MS);
      }
      return;
    }
    zatrzymajOdswiezanie();
  }

  /** Zdanie dla Operatora złożone z tego, co wie wyłącznie powłoka. */
  function zdanieORdzeniu(stan: StanRdzenia): string {
    const nasluch = stan.pracuje
      ? `Rdzeń odpowiada pod ${stan.adres}.`
      : `Rdzeń nie odpowiada pod ${stan.adres}.`;
    return `${nasluch} ${stan.opis} Dziennik: ${stan.dziennik}`;
  }

  function dopytajPowloke(): void {
    if (pytanie) return;
    pytanie = true;
    void stanRdzeniaZPowloki().then((stan) => {
      pytanie = false;
      // Odpowiedź spóźniona po odzyskaniu łączności jest już nieprawdziwa — milczy się zamiast tłumaczyć.
      if (!czynna || stan === null || biezacy === 'polaczony') return;
      diagnoza = zdanieORdzeniu(stan);
      element.dataset.rdzen = stan.pracuje ? 'nasluchuje' : 'milczy';
      opiszPodpowiedz();
    });
  }

  function pokaz(stan: StanPolaczenia): void {
    biezacy = stan;
    const wToku = W_TOKU.has(stan);
    kropka.className = `dn-kropka ${KROPKI[stan]}`;
    kropka.hidden = wToku;
    wskaznik.hidden = !wToku;
    element.dataset.stan = stan;

    // Łączność odzyskana unieważnia poprzednią diagnozę — powód ciszy przestał istnieć.
    if (stan === 'polaczony') {
      diagnoza = '';
      delete element.dataset.rdzen;
    }
    odswiezEtykiete(stan);
    zarzadzajOdswiezaniem(stan);
    if (stan !== 'polaczony') dopytajPowloke();
  }

  // Transport ogłasza subskrybentowi stan bieżący — plakietka nie czeka na pierwszą zmianę.
  const odsubskrybuj = transport.naStan(pokaz);
  pokaz(transport.stan());

  return {
    element,
    rozlacz: () => {
      // Plakietka zdjęta z paska — ani odpowiedź powłoki, ani pętla licznika nie mają już czego opisywać.
      czynna = false;
      zatrzymajOdswiezanie();
      odsubskrybuj();
    },
  };
}
