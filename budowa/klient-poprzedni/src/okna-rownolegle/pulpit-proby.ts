import { WindowRole } from '../../../shared/contract';
import { elementIkony } from '../ikony/ikony';
import { motywObowiazujacy, przelaczMotyw } from '../motyw/motyw';
import { nazwaRoli } from '../okno-komunikacji/etykiety-okna';
import { ID_GNIAZD, numerGniazda, type IdGniazda } from './identyfikatory';
import { nazwaStanuPary } from './etykiety-ukladu';
import type { StanPary } from './stan-pary';
import type { UkladOkien } from './uklad-okien';

/**
 * Pulpit strony podglądu układu.
 *
 * Wyłącznie na potrzeby podglądu wizualnego: wywołuje publiczne metody układu,
 * żeby dało się obejrzeć każdy stan sceny bez rdzenia. Punkt wejścia aplikacji
 * tego pliku nie importuje — w aplikacji te same wywołania pochodzą od modelu
 * przez narzędzia kontraktu i od pętli kolejek.
 */
export function utworzPulpitProby(uklad: UkladOkien): HTMLElement {
  const element = document.createElement('div');
  element.className = 'dn-listwa dn-proba__pulpit';

  element.append(
    grupa('Rola okna', ...ID_GNIAZD.flatMap((id) => przyciskiRol(uklad, id))),
    grupa(
      'Stan pętli',
      ...stany().map((stan) =>
        przycisk(nazwaStanuPary(stan), () => uklad.ustawStanPary(stan)),
      ),
    ),
    grupa(
      'Przekazanie zlecenia (podgląd układu)',
      przycisk(
        'Okno 1 → Okno 2',
        () => uklad.pokazPrzekazanie('okno-1', 'okno-2'),
        'dn-btn--atrament',
      ),
      przycisk(
        'Okno 1 → Okno 3',
        () => uklad.pokazPrzekazanie('okno-1', 'okno-3'),
        'dn-btn--atrament',
      ),
    ),
    grupa('Motyw', przyciskMotywu()),
  );

  return element;
}

/** Stany pętli osiągalne z pulpitu, wraz z gotowością. */
function stany(): StanPary[] {
  return ['gotowa', 'wykonawca-pracuje', 'koordynator-wybudzony', 'kolejka-wstrzymana'];
}

/** Komplet przycisków ról dla jednego gniazda. */
function przyciskiRol(uklad: UkladOkien, id: IdGniazda): HTMLElement[] {
  const numer = numerGniazda(id);
  return Object.values(WindowRole).map((rola) =>
    przycisk(`Okno ${numer}: ${nazwaRoli(rola)}`, () => uklad.nadajRole(id, rola)),
  );
}

/** Grupa pulpitu: etykieta wersalikowa i przyciski pod nią. */
function grupa(etykieta: string, ...dzieci: HTMLElement[]): HTMLElement {
  const element = document.createElement('div');
  element.className = 'dn-proba__grupa';

  const napis = document.createElement('span');
  napis.className = 'dn-proba__etykieta';
  napis.textContent = etykieta;

  const rzad = document.createElement('div');
  rzad.className = 'dn-proba__rzad';
  rzad.append(...dzieci);

  element.append(napis, rzad);
  return element;
}

/** Przycisk pulpitu — zawsze czynny, bez blokad. */
function przycisk(napis: string, dzialanie: () => void, wariant = 'dn-btn--zarys'): HTMLElement {
  const element = document.createElement('button');
  element.type = 'button';
  element.className = `dn-btn dn-btn--sm ${wariant}`;
  element.textContent = napis;
  element.addEventListener('click', dzialanie);
  return element;
}

/** Przełącznik motywu — oba motywy są równoprawne. */
function przyciskMotywu(): HTMLElement {
  const element = document.createElement('button');
  element.type = 'button';
  element.className = 'dn-btn-ikona';

  function odswiez(): void {
    const motyw = motywObowiazujacy();
    const ikona = motyw === 'dark' ? 'slonce' : 'ksiezyc';
    const opis = motyw === 'dark' ? 'Włącz motyw jasny' : 'Włącz motyw ciemny';
    element.setAttribute('aria-label', opis);
    element.title = opis;
    element.replaceChildren(elementIkony(ikona, { rozmiar: 18 }));
  }

  element.addEventListener('click', () => {
    przelaczMotyw();
    odswiez();
  });

  odswiez();
  return element;
}
