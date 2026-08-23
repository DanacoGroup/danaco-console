import type { StanPolaczenia } from '../polaczenie/stan-polaczenia';
import { nazwaRoli, nazwaSrodowiska } from './etykiety-okna';
import type { OpisOkna } from './opis-okna';

/**
 * Klasa kropki stanu połączenia. Kropka jest komponentem biblioteki
 * (`komponenty/plakietka.css`, sekcja „KROPKA SYGNAŁU”) wraz z kompletem barw
 * stanu — okno wskazuje wariant, a nie powtarza jego reguł.
 *
 * Łączenie i ponawianie dzielą barwę ostrzeżenia: obie mówią „jeszcze nie ma
 * połączenia, praca trwa”, a rozróżnia je zdanie obok kropki, nie kolor.
 */
const KLASA_KROPKI: Record<StanPolaczenia, string> = {
  polaczony: 'dn-kropka dn-kropka--sukces',
  laczenie: 'dn-kropka dn-kropka--ostrzezenie',
  ponawianie: 'dn-kropka dn-kropka--ostrzezenie',
  rozlaczony: 'dn-kropka dn-kropka--blad',
};

/** Nagłówek okna komunikacji. */
export interface Naglowek {
  element: HTMLElement;
  /** Odświeża wskaźnik stanu połączenia z rdzeniem. */
  pokazStan(stan: StanPolaczenia, oczekujace: number): void;
}

export function utworzNaglowek(opis: OpisOkna): Naglowek {
  const element = document.createElement('header');
  element.className = 'dc-naglowek';

  const tozsamosc = document.createElement('div');
  tozsamosc.className = 'dc-naglowek__tozsamosc';
  tozsamosc.append(
    pole('Projekt', opis.projekt),
    pole('Środowisko', nazwaSrodowiska(opis.srodowiskoWykonania)),
    pole('Moduł', opis.modul),
    pole('Model', opis.kanalModelu),
    pole('Rola okna', nazwaRoli(opis.rola)),
  );

  const stan = document.createElement('div');
  stan.className = 'dc-stan';
  const znacznik = document.createElement('span');
  // Stan przed pierwszym odczytem transportu jest nieznany, a nie zły —
  // odmiana neutralna jest jedyną, która tego nie przesądza.
  znacznik.className = 'dn-kropka dn-kropka--neutralna';
  const opisStanu = document.createElement('span');
  opisStanu.className = 'dc-stan__opis';
  stan.append(znacznik, opisStanu);

  element.append(tozsamosc, stan);

  return {
    element,
    pokazStan(nowy, oczekujace) {
      stan.dataset.stan = nowy;
      znacznik.className = KLASA_KROPKI[nowy];
      opisStanu.textContent = opisPolaczenia(nowy, oczekujace);
    },
  };
}

/** Pojedyncza pozycja tożsamości okna: etykieta i wartość. */
function pole(etykieta: string, wartosc: string): HTMLElement {
  const kontener = document.createElement('div');
  kontener.className = 'dc-pole';

  const nazwa = document.createElement('span');
  nazwa.className = 'dc-pole__etykieta';
  nazwa.textContent = etykieta;

  const tresc = document.createElement('span');
  tresc.className = 'dc-pole__wartosc';
  tresc.textContent = wartosc;

  kontener.append(nazwa, tresc);
  return kontener;
}

/** Opis stanu połączenia wraz z liczbą ramek oczekujących w kolejce. */
function opisPolaczenia(stan: StanPolaczenia, oczekujace: number): string {
  const ogon = oczekujace > 0 ? ` · w kolejce: ${oczekujace}` : '';
  switch (stan) {
    case 'polaczony':
      return `Połączony z rdzeniem${ogon}`;
    case 'laczenie':
      return `Łączenie z rdzeniem${ogon}`;
    case 'ponawianie':
      return `Ponawianie połączenia${ogon}`;
    case 'rozlaczony':
      return `Rozłączony${ogon}`;
  }
}
