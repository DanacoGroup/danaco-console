import { ResearchFindingStatus, type ResearchFinding } from '../../../../shared/contract';
import { przycisk } from '../../modele/kontrolki-formularza';

/**
 * Jedna pozycja wykazu ustaleń Findings Panel: byt ustalenia przełożony na wiersz z jego
 * stanem, powiązanymi źródłami i dwiema drogami działania.
 */
export interface UchwytyWiersza {
  /** Przestawia zaznaczenie ustalenia do raportu. */
  naZaznaczenie(identyfikator: string): void;
  /** Wciąga ustalenie do formularza edycji. */
  naEdycje(ustalenie: ResearchFinding): void;
}

export function utworzWierszUstalenia(
  ustalenie: ResearchFinding,
  wybrane: boolean,
  uchwyty: UchwytyWiersza,
): HTMLElement {
  const element = document.createElement('li');
  element.className = 'mr-wykaz__wiersz';
  element.dataset['ustalenie'] = ustalenie.id;

  const wybor = document.createElement('input');
  wybor.type = 'checkbox';
  wybor.className = 'dn-check mr-wykaz__wybor';
  wybor.checked = wybrane;
  wybor.setAttribute('aria-label', 'Zaznacz ustalenie do raportu');
  wybor.addEventListener('change', () => uchwyty.naZaznaczenie(ustalenie.id));

  const tresc = document.createElement('span');
  tresc.className = 'mr-wykaz__tytul';
  tresc.textContent = ustalenie.content;

  const stan = document.createElement('span');
  stan.className = 'dn-plakietka';
  const rozstrzygniete = ustalenie.status === ResearchFindingStatus.Resolved;
  stan.classList.add(rozstrzygniete ? 'dn-plakietka--sukces' : 'dn-plakietka--informacja');
  stan.textContent = rozstrzygniete ? 'rozstrzygnięte' : 'otwarte';

  const edytuj = przycisk('Edytuj', 'dn-btn dn-btn--sm dn-btn--duch');
  edytuj.addEventListener('click', () => uchwyty.naEdycje(ustalenie));

  element.append(wybor, tresc, stan, powiazania(ustalenie), edytuj);
  return element;
}

/** Źródła powiązane z ustaleniem, jako druga strona wiązania między źródłem a ustaleniem, wypisane w treści wiersza. */
function powiazania(ustalenie: ResearchFinding): HTMLElement {
  const element = document.createElement('span');
  element.className = 'mr-wykaz__meta';
  const zrodla = ustalenie.sourceIds ?? [];
  element.textContent =
    zrodla.length === 0
      ? 'bez powiązanego źródła'
      : `źródła: ${zrodla.join(' · ')}`;
  return element;
}
