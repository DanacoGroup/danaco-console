import { KluczUstawieniaOkna } from './klucze-ustawien';
import { utworzNaglowekSterowania } from './naglowek-sterowania';
import type { StanSterowania } from './stan-sterowania';
import type { ZmianaUstawienia } from './zmiana-ustawienia';

const NAZWA = 'Host wykonania';

/** Hosty znane w chwili wydania — podpowiedź w polu wpisu, a nie katalog zamknięty ograniczający treść ustawienia. */
const HOSTY_ZNANE = ['danaco-system', 'danaco-data', 'danaco-web'] as const;

/** Sterowanie nazwą hosta wykonania, zapisywane ustawieniem poziomu okna, czynne przy każdym zasięgu wykonania. */
export function utworzSterowanieHostu(
  stan: StanSterowania,
  ustawienia: ZmianaUstawienia,
): HTMLElement {
  // Identyfikator podpowiedzi zawiera identyfikator okna, by dwa komplety się nie dzieliły.
  const idPodpowiedzi = `dc-ster-hosty-${stan.idOkna()}`;

  const identyfikator = `dc-ster-host-${stan.idOkna()}`;

  const element = document.createElement('div');
  element.className = 'dc-ster-pole';

  const naglowek = utworzNaglowekSterowania(NAZWA, identyfikator);

  const pole = document.createElement('input');
  pole.id = identyfikator;
  pole.className = 'dc-ster-pole__wpis';
  pole.type = 'text';
  pole.placeholder = 'nazwa hosta, np. danaco-system';
  pole.setAttribute('list', idPodpowiedzi);
  pole.addEventListener('change', () => {
    ustawienia.zapisz(NAZWA, KluczUstawieniaOkna.HostWykonania, pole.value.trim());
  });

  element.append(naglowek, pole, podpowiedzi(idPodpowiedzi));

  stan.naZmiane((migawka) => {
    pole.value = migawka.ustawienia.hostWykonania;
  });
  pole.value = stan.migawka().ustawienia.hostWykonania;

  return element;
}

/** Lista podpowiedzi hostów w polu wpisu; nie ogranicza treści, którą operator może wpisać ręcznie w polu. */
function podpowiedzi(identyfikator: string): HTMLDataListElement {
  const lista = document.createElement('datalist');
  lista.id = identyfikator;
  for (const host of HOSTY_ZNANE) {
    const pozycja = document.createElement('option');
    pozycja.value = host;
    lista.append(pozycja);
  }
  return lista;
}
