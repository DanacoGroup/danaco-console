import { pokazKomunikat } from '../../aplikacja/komunikaty';
import { opisOdmowyBledu } from '../../komponenty/odmowa';

/** Przybornik okna modułu Apps prowadzi drogę z okna do komendy obszaru, bez blokowania kontrolek. */

/**
 * Jedno pole wejściowe narzędzia, opisane kluczem, etykietą, wartością startową i tym, czy jest
 * wieloliniowe.
 */
export interface PoleNarzedzia {
  /** Klucz, pod którym wartość trafia do `wykonaj`. */
  klucz: string;
  /** Etykieta widoczna Operatorowi; jest też podpowiedzią pola. */
  etykieta: string;
  /** Wartość startowa pola; pusta znaczy, że Operator wpisze ją sam, a nie wartość wymuszona. */
  wartosc?: string;
  /** Czy pole jest wieloliniowe — manifest i reguły JSON nie mieszczą się w linii. */
  obszerne?: boolean;
}

/**
 * Jedno narzędzie przybornika: przycisk, jego pola wejściowe i czynność wywoływana po jego
 * naciśnięciu przez Operatora.
 */
export interface NarzedzieApps {
  /** Nazwa czynności tak, jak wymienia ją opracowanie okna. */
  etykieta: string;
  /** Nazwa komendy, którą narzędzie woła — wchodzi do podpowiedzi przycisku. */
  komenda: string;
  /** Pola, o które okno pyta przed wysłaniem żądania. */
  pola?: readonly PoleNarzedzia[];
  /** Czynność narzędzia oddaje zdanie o skutku albo rzuca błąd z powodem, którego rdzeń nie podał. */
  wykonaj(wartosci: Readonly<Record<string, string>>): Promise<string>;
}

export interface PrzybornikApps {
  element: HTMLElement;
}

/**
 * Składa przybornik okna.
 *
 * `tytul` nazywa grupę czynności — okno ma ich zwykle kilka rodzajów i wykaz bez
 * nagłówka czytałby się jak lista przypadkowych przycisków.
 */
export function utworzPrzybornikApps(
  tytul: string,
  narzedzia: readonly NarzedzieApps[],
): PrzybornikApps {
  const naglowek = document.createElement('p');
  naglowek.className = 'mp-przybornik__tytul';
  naglowek.textContent = tytul;

  const lista = document.createElement('ul');
  lista.className = 'mp-przybornik';
  lista.append(...narzedzia.map(pozycjaNarzedzia));

  const element = document.createElement('div');
  element.className = 'mp-przybornik__powloka';
  element.dataset['narzedzi'] = String(narzedzia.length);
  element.append(naglowek, lista);

  return { element };
}

/**
 * Jedna pozycja przybornika wraz z polami wejściowymi i miejscem na zdanie o skutku naciśnięcia
 * przycisku.
 */
function pozycjaNarzedzia(narzedzie: NarzedzieApps): HTMLElement {
  const element = document.createElement('li');
  element.className = 'mp-przybornik__wiersz';
  element.dataset['komenda'] = narzedzie.komenda;

  const przycisk = document.createElement('button');
  przycisk.type = 'button';
  przycisk.className = 'dn-btn dn-btn--zarys dn-btn--sm';
  przycisk.textContent = narzedzie.etykieta;
  przycisk.dataset['narzedzie'] = narzedzie.komenda;
  przycisk.title = `Wywołuje komendę ${narzedzie.komenda}.`;
  przycisk.setAttribute('aria-description', `Wywołuje komendę ${narzedzie.komenda}.`);

  const skutek = document.createElement('span');
  skutek.className = 'mp-przybornik__skutek';
  skutek.setAttribute('role', 'status');

  const wejscia = new Map<string, HTMLInputElement | HTMLTextAreaElement>();
  const pola = document.createElement('span');
  pola.className = 'mp-przybornik__pola';
  for (const pole of narzedzie.pola ?? []) {
    const kontrolka = pole.obszerne === true
      ? document.createElement('textarea')
      : document.createElement('input');
    kontrolka.className = 'dn-pole-kontrolka mp-przybornik__pole';
    kontrolka.value = pole.wartosc ?? '';
    kontrolka.placeholder = pole.etykieta;
    kontrolka.setAttribute('aria-label', `${narzedzie.etykieta}: ${pole.etykieta}`);
    kontrolka.dataset['pole'] = pole.klucz;
    wejscia.set(pole.klucz, kontrolka);
    pola.append(kontrolka);
  }

  przycisk.addEventListener('click', () => {
    const wartosci: Record<string, string> = {};
    for (const [klucz, kontrolka] of wejscia) wartosci[klucz] = kontrolka.value;

    skutek.textContent = `Wołam ${narzedzie.komenda}…`;
    void narzedzie
      .wykonaj(wartosci)
      .then((zdanie) => {
        skutek.textContent = zdanie;
      })
      .catch((powod: unknown) => {
        const zdanie = powod instanceof Error ? powod.message : String(powod);
        skutek.textContent = zdanie;
        // Komunikat obok zdania pod przyciskiem, bo Operator patrzy w tej chwili na kontrolkę.
        pokazKomunikat({ tytul: narzedzie.etykieta, tresc: zdanie, waga: 'ostrz' });
      });
  });

  element.append(przycisk, pola, skutek);
  return element;
}

/**
 * Zdanie o odmowie rdzenia, gotowe do rzucenia z `wykonaj`.
 *
 * Wspólne dla wszystkich narzędzi, bo odmowa czyta się tak samo niezależnie od
 * tego, która komenda ją zwróciła — a różne zdania o tym samym byłyby różnymi
 * prawdami o zachowaniu rdzenia.
 */
export function bladOdmowyNarzedzia(czynnosc: string, blad: unknown): Error {
  return new Error(opisOdmowyBledu(czynnosc, blad as never));
}

/**
 * Sprawdzenie okna modułu przed wysłaniem żądania.
 *
 * `windowId` jest polem wymaganym każdej komendy obszaru. Bez niego żądanie
 * poszłoby po dane niczyje, a odmowa rdzenia mówiłaby o brakującym polu zamiast
 * o tym, że sesja nie ma jeszcze okna modułu.
 */
export function wymagajOknaModulu(idOkna: string): string {
  if (idOkna !== '') return idOkna;
  throw new Error(
    'Rdzeń nie wskazał jeszcze okna modułu Apps, a każda komenda obszaru wymaga ' +
      'jego identyfikatora — żądanie nie zostało wysłane.',
  );
}

/** Sprawdzenie pola, którego kontrakt wymaga, a Operator zostawił puste, rzucające błąd z nazwą tego pola. */
export function wymagajPola(wartosc: string, etykieta: string): string {
  const przyciete = wartosc.trim();
  if (przyciete !== '') return przyciete;
  throw new Error(`Pole „${etykieta}" jest wymagane przez kontrakt — żądanie nie poszło.`);
}

/**
 * Odczyt liczby z pola tekstowego. Puste pole znaczy „bez wskazania" i wraca
 * jako liczba ujemna — źródło pomija wtedy pole żądania, zamiast wysyłać zero,
 * które w kontrakcie jest wartością, nie brakiem.
 */
export function liczbaZPola(wartosc: string): number {
  const przyciete = wartosc.trim();
  if (przyciete === '') return -1;
  const liczba = Number(przyciete);
  if (!Number.isFinite(liczba)) {
    throw new Error(`Wartość „${przyciete}" nie jest liczbą — żądanie nie poszło.`);
  }
  return liczba;
}

/**
 * Odczyt JSON-a z pola tekstowego. Puste pole znaczy „bez wskazania".
 *
 * Rozbiór po stronie klienta jest po to, żeby literówka wróciła zdaniem o niej,
 * a nie odmową rdzenia o niepoprawnym ładunku — Operator poprawia wtedy pole,
 * które właśnie wypełniał.
 */
export function jsonZPola(wartosc: string, etykieta: string): unknown {
  const przyciete = wartosc.trim();
  if (przyciete === '') return undefined;
  try {
    return JSON.parse(przyciete);
  } catch (powod) {
    const szczegol = powod instanceof Error ? powod.message : String(powod);
    throw new Error(`Pole „${etykieta}" nie jest poprawnym JSON-em: ${szczegol}`);
  }
}
