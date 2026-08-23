import type { Kanal } from '../../protokol/kanal';
import { utworzKatalogOkien, type KatalogOkien } from '../katalog-okien';
import { bilansFunkcji, KATALOG_FUNKCJI } from './katalog-funkcji';

/**
 * Pas uczciwości modułu Roundtable — trzy rzeczy, których żadne pojedyncze okno
 * nie widzi w całości.
 *
 * Pierwsza: rozjazd katalogu okien rdzenia z oknami, które moduł buduje. Rdzeń
 * przypisuje modułowi cztery okna operacyjne (`okno_modulu`), a opracowanie
 * modułu wymienia sześć okien własnych — Argument Map & Analysis i Voting &
 * Evaluation Center nie mają w katalogu rdzenia ani jednego wiersza. Moduł je
 * buduje mimo to, bo bez nich okien byłoby cztery zamiast sześciu; zdanie
 * rozjazdu liczy `utworzKatalogOkien` z odczytu `module.list`, nie z napisu.
 *
 * Druga: stosunek liczby narzędzi opracowania do stanu ich wykonania. Bilans
 * mówi, ile moduł wykonuje w całości, ile częściowo, ile ma w kontrakcie
 * komendę bez zbudowanej obsługi i ile nie ma pokrycia. Rozróżnienie dwóch
 * ostatnich liczb jest po scaleniu kontraktu najważniejszą treścią pasa:
 * obszar urósł z czterech komend do czterdziestu sześciu, więc niemal wszystko,
 * co moduł nazywał brakiem, jest dziś pracą do wykonania, a nie brakiem
 * uzgodnienia.
 *
 * Trzecia: odczyt niewywoływany. `roundtable.debate.get` oddaje skład, tury
 * i wypowiedzi, a `roundtable.model.list` sam skład — żadnego z nich moduł
 * jeszcze nie wywołuje, więc okno otwarte w trakcie debaty zna wyłącznie to,
 * co usłyszało od swojego otwarcia. Dotyczy to każdego okna modułu naraz,
 * dlatego zdanie stoi na pasie, a nie w jednym z nich.
 *
 * Pas niczego nie blokuje i niczego nie ocenia — podaje liczby i nazywa granicę.
 */

/** Pas uczciwości wraz z katalogiem, który trzeba domknąć przy zejściu modułu. */
export interface PasekUczciwosci {
  element: HTMLElement;
  /** Katalog okien rdzenia — do odczytu po montażu i do domknięcia przy zejściu. */
  katalog: KatalogOkien;
}

/** Kody okien operacyjnych, które moduł naprawdę buduje. */
export const OKNA_BUDOWANE: readonly string[] = [
  'model-panels',
  'debate-panel',
  'argument-map-analysis',
  'voting-evaluation-center',
  'moderator-panel',
  'consensus-panel',
];

export function utworzPasekUczciwosci(kanal: Kanal): PasekUczciwosci {
  const katalog = utworzKatalogOkien(kanal, 'roundtable', OKNA_BUDOWANE);
  const bilans = bilansFunkcji(KATALOG_FUNKCJI);

  const element = document.createElement('section');
  element.className = 'dr-uczciwosc';
  element.setAttribute('aria-label', 'Czego moduł Roundtable nie umie');

  const narzedzia = document.createElement('p');
  narzedzia.className = 'dn-pole-opis dr-uczciwosc__zdanie';
  narzedzia.textContent =
    `Opracowanie modułu wymienia ${bilans.wszystkich} narzędzi, a moduł wywołuje cztery komendy obszaru: ` +
    'dodanie uczestnika, uruchomienie tury, czynność moderatora i odczyt stanowiska. ' +
    `Wykonane w całości: ${bilans.wykonanych} · wykonane częściowo: ${bilans.czesciowych} · ` +
    `komenda w kontrakcie bez zbudowanej obsługi: ${bilans.bezObslugi} · ` +
    `bez pokrycia w kontrakcie: ${bilans.bezPokrycia} · poza tym modułem: ${bilans.pozaModulem}. ` +
    'Wykaz narzędzi każdego okna wraz z nazwą komendy stoi w tym oknie, na dole.';

  const odczyt = document.createElement('p');
  odczyt.className = 'dn-pole-opis dr-uczciwosc__zdanie';
  odczyt.textContent =
    'Skład i przebieg debaty przychodzą do modułu wyłącznie zdarzeniem roundtable.debate.changed i fragmentami ' +
    'strumienia. Odczyt na żądanie jest w kontrakcie — roundtable.debate.get oddaje skład, tury i wypowiedzi, ' +
    'roundtable.model.list sam skład — ale moduł żadnego z nich jeszcze nie wywołuje. Okno otwarte w trakcie ' +
    'debaty widzi zatem sam ogon: wypowiedzi sprzed jego otwarcia są niewczytane, a nie zgubione.';

  element.append(
    katalog.zdanieElement('dn-pole-opis dr-uczciwosc__zdanie'),
    narzedzia,
    odczyt,
  );
  return { element, katalog };
}
