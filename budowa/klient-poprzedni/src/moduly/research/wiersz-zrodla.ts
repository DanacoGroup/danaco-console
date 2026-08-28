import { ResearchCredibility, type ResearchSource } from '../../../../shared/contract';
import { przycisk } from '../../modele/kontrolki-formularza';

/** Jedna pozycja wykazu źródeł Sources Manager przekłada byt źródła na wiersz z oceną wiarygodności; wariant plakietki dla tej oceny nie zależy od samej barwy. */
const PLAKIETKI: Readonly<Record<string, string>> = {
  [ResearchCredibility.High]: 'dn-plakietka--sukces',
  [ResearchCredibility.Medium]: 'dn-plakietka--informacja',
  [ResearchCredibility.Low]: 'dn-plakietka--ostrzezenie',
  [ResearchCredibility.Unverified]: 'dn-plakietka--blad',
};

/** Nazwa oceny wiarygodności źródła pokazywana w wierszu wykazu, dobrana do wartości oceny oddanej przez rdzeń. */
const OCENY: Readonly<Record<string, string>> = {
  [ResearchCredibility.High]: 'wiarygodność wysoka',
  [ResearchCredibility.Medium]: 'wiarygodność średnia',
  [ResearchCredibility.Low]: 'wiarygodność niska',
  [ResearchCredibility.Unverified]: 'źródło niezweryfikowane',
};

/** Dwie drogi działania na pozycji źródła: zaznaczenie wielokrotne do eksportu oraz otwarcie materiału w lekturze. */
export interface UchwytyZrodla {
  /** Przestawia zaznaczenie źródła — wybór wielokrotny wspólny dwóm oknom. */
  naZaznaczenie(identyfikator: string): void;
  /** Wskazuje źródło do lektury w Reading View. */
  naLekture(identyfikator: string): void;
}

export function utworzWierszZrodla(
  zrodlo: ResearchSource,
  wybrane: boolean,
  uchwyty: UchwytyZrodla,
): HTMLElement {
  const element = document.createElement('li');
  element.className = 'mr-wykaz__wiersz';
  element.dataset['zrodlo'] = zrodlo.id;

  const wybor = document.createElement('input');
  wybor.type = 'checkbox';
  wybor.className = 'dn-check mr-wykaz__wybor';
  wybor.checked = wybrane;
  wybor.setAttribute('aria-label', `Zaznacz źródło ${zrodlo.title}`);
  wybor.addEventListener('change', () => uchwyty.naZaznaczenie(zrodlo.id));

  const tytul = document.createElement('span');
  tytul.className = 'mr-wykaz__tytul';
  tytul.textContent = zrodlo.title;

  const ocena = document.createElement('span');
  ocena.className = 'dn-plakietka';
  ocena.classList.add(PLAKIETKI[zrodlo.credibility] ?? 'dn-plakietka--informacja');
  ocena.textContent = OCENY[zrodlo.credibility] ?? zrodlo.credibility;

  const czytaj = przycisk('Czytaj', 'dn-btn dn-btn--sm dn-btn--duch');
  czytaj.setAttribute('aria-label', `Czytaj źródło ${zrodlo.title}`);
  czytaj.addEventListener('click', () => uchwyty.naLekture(zrodlo.id));

  element.append(wybor, tytul, ocena, czytaj, metadane(zrodlo));
  return element;
}

/** Metadane wykazu źródła: rodzaj, pochodzenie, adres oraz data pozyskania, wypisane jednym wierszem opisowym. */
function metadane(zrodlo: ResearchSource): HTMLElement {
  const element = document.createElement('span');
  element.className = 'mr-wykaz__meta';
  const czesci = [
    `typ: ${zrodlo.kind}`,
    zrodlo.origin === undefined ? '' : `pochodzenie: ${zrodlo.origin}`,
    zrodlo.url === undefined ? '' : zrodlo.url,
    zrodlo.libraryFileId === undefined ? '' : `plik repozytorium: ${zrodlo.libraryFileId}`,
    `pozyskane: ${new Date(zrodlo.acquiredAt).toLocaleString('pl-PL')}`,
  ];
  element.textContent = czesci.filter((czesc) => czesc !== '').join(' · ');
  return element;
}
