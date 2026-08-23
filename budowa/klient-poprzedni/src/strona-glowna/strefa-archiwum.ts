import type { Session } from '../../../shared/contract';
import type { CzynnosciSesji } from './czynnosci-sesji';
import { utworzKarteSesji } from './karta-sesji';
import type { OdbiorcaMeldunku } from './menu-sesji';
import type { WykazSrodowisk } from './wykaz-srodowisk';
import type { WpisSesji } from './zrodlo-sesji';

/**
 * Archiwum sesji — wgląd rozwijany pod wykazem sesji w tle.
 *
 * Pokazuje sesje odłożone do archiwum i oddaje je do przywrócenia. Wykaz idzie
 * z `session.archive.list`, a przywrócenie z `session.restore`; obie komendy
 * wykonuje wpięcie, ten plik zna wyłącznie ich wynik.
 *
 * Strefa startuje zwinięta i odpytuje rdzeń dopiero przy rozwinięciu, a potem
 * przy każdym kolejnym — archiwum bywa długie, nie dotyczy pracy bieżącej,
 * a sesja mogła w międzyczasie wrócić albo dojść.
 *
 * Liczba wszystkich sesji archiwum przychodzi obok strony wyników, więc
 * nagłówek pokazuje obie liczby, gdy się różnią; sama długość strony nie jest
 * liczbą sesji w archiwum.
 *
 * Odpytywanie, archiwum puste i odmowa rdzenia mają osobne napisy — wspólny
 * zacierałby różnicę między brakiem danych a brakiem odpowiedzi.
 */

const ETYKIETA = 'Archiwum sesji';

/** Odpytanie archiwum; wpięcie podaje wykonanie komendy kontraktu. */
export type OdpytanieArchiwum = () => Promise<WynikArchiwum>;

export interface WynikArchiwum {
  /** Sesje strony wyników; puste przy odmowie. */
  sesje: readonly Session[];
  /** Liczba wszystkich sesji archiwum wg rdzenia. */
  razem: number;
  /** Treść odmowy rdzenia; obecna wyłącznie przy niepowodzeniu. */
  blad?: string;
}

export interface StrefaArchiwum {
  element: HTMLElement;
  /** Ponawia odpytanie, jeśli archiwum jest rozwinięte; zwinięte pomija. */
  odswiez(): void;
}

export function utworzStrefeArchiwum(
  odpytaj: OdpytanieArchiwum,
  czynnosci: CzynnosciSesji,
  meldunek: OdbiorcaMeldunku,
  srodowiska: WykazSrodowisk,
): StrefaArchiwum {
  const element = document.createElement('details');
  element.className = 'dn-strona__archiwum';

  const naglowek = document.createElement('summary');
  naglowek.className = 'dn-etykieta-wersalikowa dn-strona__archiwum-naglowek';
  naglowek.textContent = ETYKIETA;

  const tresc = document.createElement('div');
  tresc.className = 'dn-strona__sesje';
  tresc.setAttribute('aria-live', 'polite');
  element.append(naglowek, tresc);

  /** Numer odpytania — odpowiedź spóźniona nie nadpisuje świeższej. */
  let numer = 0;

  function pokaz(napis: string): void {
    const pustka = document.createElement('p');
    // Ta sama klasa zawężająca, co pustka wykazu sesji: domyślne odstępy stanu
    // pustego zajmują tu więcej pionu niż cała linia kafli ustawień.
    pustka.className = 'dn-pusty-stan dn-strona__archiwum-pustka';
    pustka.textContent = napis;
    tresc.replaceChildren(pustka);
  }

  function pokazWykaz(wynik: WynikArchiwum): void {
    naglowek.textContent =
      wynik.razem > wynik.sesje.length
        ? `${ETYKIETA} — ${wynik.sesje.length} z ${wynik.razem}`
        : `${ETYKIETA} — ${wynik.razem}`;

    const wykaz = document.createElement('ul');
    wykaz.className = 'dn-strona__sesje-wykaz';
    for (const sesja of wynik.sesje) {
      const wpis: WpisSesji = { sesja };
      wykaz.append(
        utworzKarteSesji(wpis, { powrot: null, czynnosci, meldunek, srodowiska }).element,
      );
    }
    tresc.replaceChildren(wykaz);
  }

  async function wczytaj(): Promise<void> {
    const moj = ++numer;
    pokaz('Odpytujemy rdzeń o archiwum…');
    const wynik = await odpytaj();
    if (moj !== numer) return;

    if (wynik.blad !== undefined) {
      pokaz(wynik.blad);
      return;
    }
    if (wynik.sesje.length === 0) {
      naglowek.textContent = `${ETYKIETA} — 0`;
      pokaz('Archiwum jest puste — żadna sesja nie została jeszcze odłożona.');
      return;
    }
    pokazWykaz(wynik);
  }

  element.addEventListener('toggle', () => {
    if (element.open) void wczytaj();
  });

  return {
    element,
    odswiez() {
      if (element.open) void wczytaj();
    },
  };
}
