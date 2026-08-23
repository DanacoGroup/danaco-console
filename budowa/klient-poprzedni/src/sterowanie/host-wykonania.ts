import { KluczUstawieniaOkna } from './klucze-ustawien';
import { utworzNaglowekSterowania } from './naglowek-sterowania';
import type { StanSterowania } from './stan-sterowania';
import type { ZmianaUstawienia } from './zmiana-ustawienia';

const NAZWA = 'Host wykonania';

/**
 * Hosty znane w chwili wydania — podpowiedź, nie katalog zamknięty.
 *
 * Nazwa hosta jest ustawieniem okna, a nie wartością wyliczenia `ExecutionEnv`,
 * dzięki czemu dopisanie kolejnego hosta nie wymaga zmiany kodu. Pole przyjmuje
 * dowolną nazwę; poniższa lista skraca drogę do trzech używanych dziś.
 */
const HOSTY_ZNANE = ['danaco-system', 'danaco-data', 'danaco-web'] as const;

/**
 * Sterowanie nazwą hosta wykonania.
 *
 * Wartość idzie ustawieniem poziomu okna (`config.set`, zasięg `window`),
 * ponieważ treść `window.update` nie ma dla niej pola. Pole pozostaje czynne
 * przy każdym zasięgu wykonania; przy zasięgu `local` i `core` wpis nie ma
 * zastosowania.
 */
export function utworzSterowanieHostu(
  stan: StanSterowania,
  ustawienia: ZmianaUstawienia,
): HTMLElement {
  // Identyfikator listy podpowiedzi zawiera identyfikator okna, więc dwa
  // komplety otwarte obok siebie nie dzielą jednego elementu dokumentu.
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

/** Lista podpowiedzi hostów; nie ogranicza treści wpisu. */
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
