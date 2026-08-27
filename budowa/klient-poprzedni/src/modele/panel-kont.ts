import { poleWyboru, przycisk } from './kontrolki-formularza';
import { NAZWY_RODZAJOW, RODZAJE_KONT, rodzajZWyboru } from './rodzaje-kont';
import type { StanKont } from './stan-kont';
import { utworzFormularzKonta } from './formularz-konta';
import { utworzWykazKont } from './wykaz-kont';

/**
 * Panel kont, czyli wykaz kont po lewej stronie i formularz jednego konta po
 * prawej. Panel pamięta wyłącznie, które konto jest już wypełnione
 * w formularzu, i wypełnia formularz dopiero przy zmianie konta czynnego.
 */
export interface PanelKont {
  /** Panel osadzany w sekcji modeli. */
  element: HTMLElement;
  /** Nanosi stan na wykaz i formularz. */
  odswiez(): void;
}

export function utworzPanelKont(stan: StanKont): PanelKont {
  const wykaz = utworzWykazKont(stan);
  const formularz = utworzFormularzKonta(stan);

  const rodzaj = poleWyboru({ etykieta: 'Rodzaj na wykazie' }, [
    { wartosc: '', etykieta: 'wszystkie rodzaje' },
    ...RODZAJE_KONT.map((wartosc) => ({ wartosc, etykieta: NAZWY_RODZAJOW[wartosc] })),
  ]);

  const noweKonto = przycisk('Nowe konto', 'dn-btn dn-btn--sm dn-btn--atrament');
  const odczyt = przycisk('Odczytaj rejestr', 'dn-btn dn-btn--sm dn-btn--zarys');

  const pasek = document.createElement('header');
  pasek.className = 'dm-panel__pasek';
  pasek.append(rodzaj.element, rozpychacz(), noweKonto, odczyt);

  const cialo = document.createElement('div');
  cialo.className = 'dm-panel__cialo';
  cialo.append(wykaz.element, formularz.element);

  const element = document.createElement('section');
  element.className = 'dm-panel';
  element.append(pasek, cialo);

  /** Konto już wypełnione w formularzu; `null` znaczy formularz konta nowego. */
  let wypelnione: string | null = null;

  function odswiez(): void {
    wykaz.odswiez();
    const czynne = stan.wybrane()?.id ?? null;
    if (czynne !== wypelnione) {
      wypelnione = czynne;
      formularz.pokaz(stan.wybrane());
    }
  }

  rodzaj.kontrolka.addEventListener('change', () => {
    void stan.ustawRodzaj(rodzajZWyboru(rodzaj.kontrolka.value));
  });

  noweKonto.addEventListener('click', () => {
    // Wybór zdejmujemy zawsze, także gdy już go nie ma, aby formularz wrócił do stanu czystego.
    wypelnione = null;
    formularz.pokaz(null);
    stan.wybierz(null);
  });

  odczyt.addEventListener('click', () => void stan.odswiez());

  return { element, odswiez };
}

function rozpychacz(): HTMLElement {
  const element = document.createElement('span');
  element.className = 'dm-panel__rozpychacz';
  return element;
}
