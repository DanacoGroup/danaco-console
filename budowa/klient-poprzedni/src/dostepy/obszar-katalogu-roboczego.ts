import type { Kanal } from '../protokol/kanal';
import { KLUCZ_PODSTAWA, KLUCZ_WZORZEC, OPIS_PODSTAWY_DOMYSLNEJ } from './klucze-katalogu';
import { utworzPoleKatalogu, type PoleKatalogu } from './pole-katalogu';
import { utworzStanKataloguRoboczego } from './stan-katalogu-roboczego';

/**
 * Obszar katalogu roboczego — ustawienie osobne od dostępów. Dostęp mówi,
 * do czego model sięga, a katalog roboczy mówi, gdzie model zostawia swoje
 * katalogi sesyjne i pliki robocze. Obszar pokazuje dwa klucze i nie
 * zastępuje okna konfiguracji.
 */
export interface ObszarKataloguRoboczego {
  /** Element osadzany w sekcji dostępów. */
  element: HTMLElement;
  /** Wczytuje definicje i wartości z rdzenia. */
  wczytaj(): void;
  /** Odłącza subskrypcję `config.changed`. */
  rozlacz(): void;
}

export function utworzObszarKataloguRoboczego(kanal: Kanal): ObszarKataloguRoboczego {
  const stan = utworzStanKataloguRoboczego(kanal);

  const pola: PoleKatalogu[] = [
    utworzPoleKatalogu({ klucz: KLUCZ_PODSTAWA, stan, zWyborem: true }),
    utworzPoleKatalogu({ klucz: KLUCZ_WZORZEC, stan, zWyborem: false }),
  ];

  const tytul = document.createElement('h3');
  tytul.className = 'dd-katalog__tytul';
  tytul.textContent = 'Katalog roboczy modelu';

  const rozroznienie = document.createElement('p');
  rozroznienie.className = 'dd-katalog__rozroznienie';
  rozroznienie.textContent =
    'To NIE JEST dostęp. Dostęp mówi, do jakich maszyn i katalogów model sięga; katalog roboczy mówi, gdzie powstają jego katalogi sesyjne i pliki robocze. Dwa niezależne ustawienia.';

  const domyslna = document.createElement('p');
  domyslna.className = 'dd-katalog__domyslna';
  domyslna.textContent = OPIS_PODSTAWY_DOMYSLNEJ;

  const naglowek = document.createElement('header');
  naglowek.className = 'dd-katalog__naglowek';
  naglowek.append(tytul, rozroznienie, domyslna);

  const formularz = document.createElement('div');
  formularz.className = 'dd-katalog__formularz';
  formularz.append(...pola.map((pole) => pole.element));

  const element = document.createElement('section');
  element.className = 'dd-katalog';
  element.append(naglowek, formularz);

  const odswiez = (): void => {
    for (const pole of pola) pole.odswiez();
  };

  stan.naZmiane(odswiez);
  odswiez();

  return {
    element,

    // Wczytanie nie blokuje osadzenia obszaru: pola stoją od razu
    // z wartościami domyślnymi.
    wczytaj: () => void stan.odswiez().then(odswiez),

    rozlacz: stan.rozlacz,
  };
}
