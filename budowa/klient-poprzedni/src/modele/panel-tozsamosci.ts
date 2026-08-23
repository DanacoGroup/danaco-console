import { utworzEdytorTozsamosci } from './edytor-tozsamosci';
import { przycisk } from './kontrolki-formularza';
import type { StanTozsamosci } from './stan-tozsamosci';
import { utworzWykazKategorii } from './wykaz-kategorii';

/**
 * Panel tożsamości — wykaz kategorii po lewej, edytor jednej kategorii po
 * prawej.
 *
 * Panel łączy oba byty i zleca im odświeżenie. Wiedza o tym, co wolno nadpisać,
 * siedzi w edytorze, bo tylko on wie, czy pole jest właśnie wypełniane.
 */
export interface PanelTozsamosci {
  /** Panel osadzany w sekcji modeli. */
  element: HTMLElement;
  /** Nanosi stan na wykaz i edytor. */
  odswiez(): void;
}

export function utworzPanelTozsamosci(stan: StanTozsamosci): PanelTozsamosci {
  const wykaz = utworzWykazKategorii(stan);
  const edytor = utworzEdytorTozsamosci(stan);

  const odczyt = przycisk('Odczytaj katalog ponownie', 'dn-btn dn-btn--sm dn-btn--zarys');
  odczyt.addEventListener('click', () => void stan.odswiez());

  const pasek = document.createElement('header');
  pasek.className = 'dm-panel__pasek';
  pasek.append(wyjasnienie(), odczyt);

  const cialo = document.createElement('div');
  cialo.className = 'dm-panel__cialo';
  cialo.append(wykaz.element, edytor.element);

  const element = document.createElement('section');
  element.className = 'dm-panel';
  element.append(pasek, cialo);

  function odswiez(): void {
    wykaz.odswiez();
    edytor.odswiez();
  }

  return { element, odswiez };
}

/** Zdanie w pasku panelu, objaśniające rolę kategorii zasad. */
function wyjasnienie(): HTMLElement {
  const element = document.createElement('p');
  element.className = 'dm-panel__wyjasnienie';
  element.textContent =
    'Kategorie zasad i tożsamości modelu pochodzą z katalogu rdzenia. Treść zapisana tutaj składa się na nakładkę systemową wskazanej osi; tryb podania rozstrzyga, czy prompt fabryczny zostanie zastąpiony.';
  return element;
}
