import { AppDeployStatus, type AppDeployment } from '../../../../shared/contract';

/**
 * Tabela wdrożeń — „Przegląd statusu publikacji” okna Deployment Panel.
 *
 * Nagłówek kolumn zostaje przy pustym wykazie: pusty stan wchodzi w miejsce
 * ciała tabeli, a nie zamiast całej tabeli — Operator ma widzieć, czego wykaz
 * dotyczy, zanim cokolwiek w nim będzie.
 *
 * Błąd wdrożenia znakuje plakietka w kolumnie stanu, bez zmiany struktury
 * wiersza. Stan nigdy nie opiera się na samej barwie — plakietka niesie słowo.
 *
 * Gdy silnik wdrożeń wpisze do pola `logRef` zdanie o tym, dlaczego przebieg
 * się nie powiódł, tabela pokazuje jego treść tak, jak przyszła, i mówi,
 * z którego pola pochodzi. Samo słowo `failed` zostawiałoby Operatora bez
 * wyjaśnienia, które przyszło w tej samej ramce.
 */
export interface TabelaWdrozen {
  element: HTMLElement;
  /** Nanosi wykaz wdrożeń; pusty wykaz zostawia sam nagłówek kolumn. */
  nanies(wdrozenia: readonly AppDeployment[]): void;
}

/** Wariant plakietki dla stanu wdrożenia; nazwy z biblioteki `komponenty/`. */
const PLAKIETKI: Readonly<Record<string, string>> = {
  [AppDeployStatus.Pending]: 'dn-plakietka',
  [AppDeployStatus.Running]: 'dn-plakietka dn-plakietka--informacja',
  [AppDeployStatus.Succeeded]: 'dn-plakietka dn-plakietka--sukces',
  [AppDeployStatus.Failed]: 'dn-plakietka dn-plakietka--blad',
  [AppDeployStatus.RolledBack]: 'dn-plakietka dn-plakietka--ostrzezenie',
};

const KOLUMNY = ['Wersja', 'Środowisko', 'Strategia', 'Stan', 'Cofnięcie'];

export function utworzTabeleWdrozen(cofnij: (idWdrozenia: string) => void): TabelaWdrozen {
  const naglowek = document.createElement('tr');
  naglowek.append(
    ...KOLUMNY.map((tytul) => {
      const komorka = document.createElement('th');
      komorka.scope = 'col';
      komorka.textContent = tytul;
      return komorka;
    }),
  );

  const glowa = document.createElement('thead');
  glowa.append(naglowek);

  const cialo = document.createElement('tbody');

  const element = document.createElement('table');
  element.className = 'dn-tabela mp-wdrozenia';
  element.append(glowa, cialo);

  return {
    element,
    nanies(wdrozenia) {
      cialo.replaceChildren(...wdrozenia.map((wdrozenie) => wiersz(wdrozenie, cofnij)));
    },
  };
}

/** Jeden wiersz wykazu wdrożeń wraz z przyciskiem cofnięcia do tej wersji. */
function wiersz(wdrozenie: AppDeployment, cofnij: (idWdrozenia: string) => void): HTMLElement {
  const element = document.createElement('tr');
  element.dataset['wdrozenie'] = wdrozenie.id;

  element.append(
    komorka(wdrozenie.version ?? 'bez numeru'),
    komorka(wdrozenie.environment),
    komorka(wdrozenie.strategy),
    komorkaStanu(wdrozenie),
    komorkaCofniecia(wdrozenie, cofnij),
  );
  return element;
}

function komorka(tresc: string): HTMLElement {
  const element = document.createElement('td');
  element.textContent = tresc;
  return element;
}

/**
 * Kolumna stanu: plakietka ze słowem stanu, bez zmiany struktury wiersza,
 * a pod nią — powód z pola `logRef`, gdy rdzeń go przysłał.
 */
function komorkaStanu(wdrozenie: AppDeployment): HTMLElement {
  const plakietka = document.createElement('span');
  plakietka.className = PLAKIETKI[wdrozenie.status] ?? 'dn-plakietka';
  plakietka.textContent = wdrozenie.status;

  const element = document.createElement('td');
  element.append(plakietka);
  if (wdrozenie.logRef !== undefined && wdrozenie.logRef !== '') {
    const powod = document.createElement('span');
    powod.className = 'mp-wdrozenia__powod';
    powod.textContent = `logRef: ${wdrozenie.logRef}`;
    element.append(powod);
  }
  return element;
}

/**
 * Kolumna cofnięcia. Przycisk jest czynny zawsze — także dla wdrożenia
 * nieudanego; o dopuszczalności cofnięcia rozstrzyga rdzeń, nie
 * wygaszona kontrolka.
 */
function komorkaCofniecia(
  wdrozenie: AppDeployment,
  cofnij: (idWdrozenia: string) => void,
): HTMLElement {
  const przycisk = document.createElement('button');
  przycisk.type = 'button';
  przycisk.className = 'dn-btn dn-btn--zarys dn-btn--sm';
  przycisk.textContent = 'Cofnij do tej wersji';
  przycisk.addEventListener('click', () => cofnij(wdrozenie.id));

  const element = document.createElement('td');
  element.append(przycisk);
  return element;
}
