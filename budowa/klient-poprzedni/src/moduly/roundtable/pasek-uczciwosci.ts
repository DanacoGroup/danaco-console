import type { Kanal } from '../../protokol/kanal';
import { utworzKatalogOkien, type KatalogOkien } from '../katalog-okien';
import { bilansFunkcji, KATALOG_FUNKCJI } from './katalog-funkcji';

/** Pas uczciwości modułu Roundtable, wraz z katalogiem okien do domknięcia przy zejściu: trzy rzeczy, których żadne okno nie widzi w całości. */
export interface PasekUczciwosci {
  element: HTMLElement;
  /** Katalog okien rdzenia — do odczytu po montażu i do domknięcia przy zejściu. */
  katalog: KatalogOkien;
}

/** Kody okien operacyjnych, które moduł naprawdę buduje, użyte do rozjazdu wobec katalogu okien rdzenia. */
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
