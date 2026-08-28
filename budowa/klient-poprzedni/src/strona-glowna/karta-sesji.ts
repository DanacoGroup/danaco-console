import { etykietaStanuSesji, opisCzasu, opisOkien } from './etykiety-sesji';
import type { CzynnosciSesji } from './czynnosci-sesji';
import { utworzMenuSesji, type OdbiorcaMeldunku } from './menu-sesji';
import type { WykazSrodowisk } from './wykaz-srodowisk';
import type { WpisSesji } from './zrodlo-sesji';

/** Karta sesji w tle jest jednym wierszem wykazu sekcji sesji o jednej odpowiedzialności. */

/** Czynność powrotu do sesji zwraca obietnicę, która trzyma przycisk w stanie zajętości aż do jej rozstrzygnięcia. */
export type CzynnoscPowrotu = (wpis: WpisSesji) => Promise<void>;

export interface KartaSesji {
  element: HTMLLIElement;
}

/** Otoczenie wiersza karty określa, co wolno zrobić z sesją i gdzie zgłaszać wynik czynności wykonanej. */
export interface OtoczenieKarty {
  powrot: CzynnoscPowrotu | null;
  czynnosci: CzynnosciSesji;
  meldunek: OdbiorcaMeldunku;
  /** Wykaz środowisk strony jest jedynym źródłem nazwy środowiska wyświetlanej w tym wierszu karty sesji. */
  srodowiska: WykazSrodowisk;
}

export function utworzKarteSesji(wpis: WpisSesji, otoczenie: OtoczenieKarty): KartaSesji {
  const { powrot, czynnosci, meldunek } = otoczenie;
  const element = document.createElement('li');
  element.className = 'dn-strona__sesja';
  element.dataset.sesja = wpis.sesja.id;

  const tresc = document.createElement('div');
  tresc.className = 'dn-strona__sesja-tresc';
  tresc.append(zbudujTytul(wpis), zbudujMete(wpis, otoczenie.srodowiska));

  element.append(zbudujSygnal(wpis), tresc, zbudujStan(wpis));
  if (powrot !== null) element.append(zbudujPowrot(wpis, powrot));
  const menu = utworzMenuSesji(wpis, czynnosci, meldunek);
  if (menu !== null) element.append(menu.element);

  return { element };
}

/** Kropka sygnału pokazuje tętno przy sesji strumieniującej, a pełne wypełnienie przy sesji żywej na rdzeniu. */
function zbudujSygnal(wpis: WpisSesji): HTMLElement {
  const kropka = document.createElement('span');
  const strumieniuje = (wpis.obecnosc?.streamingWindowCount ?? 0) > 0;
  const zywa = wpis.obecnosc?.live === true;
  kropka.className = strumieniuje
    ? 'dn-kropka dn-kropka--sukces dn-kropka--tetno'
    : `dn-kropka ${zywa ? 'dn-kropka--sukces' : 'dn-kropka--neutralna'}`;
  kropka.setAttribute('aria-hidden', 'true');
  return kropka;
}

function zbudujTytul(wpis: WpisSesji): HTMLElement {
  const tytul = document.createElement('span');
  tytul.className = 'dn-strona__sesja-tytul';
  tytul.textContent = wpis.sesja.title ?? wpis.sesja.id;
  return tytul;
}

/** Druga linia wiersza karty niesie nazwę środowiska, moduł, liczbę okien i chwilę ostatniej czynności. */
function zbudujMete(wpis: WpisSesji, wykaz: WykazSrodowisk): HTMLElement {
  const meta = document.createElement('span');
  meta.className = 'dn-strona__sesja-meta';

  const czesci: string[] = [];
  const srodowisko = wykaz.nazwa(wpis.obecnosc?.environmentCode);
  if (srodowisko !== undefined) czesci.push(srodowisko);
  const modul = wpis.obecnosc?.moduleCode;
  if (modul !== undefined && modul !== '') czesci.push(`moduł ${modul}`);
  if (wpis.obecnosc !== undefined) czesci.push(opisOkien(wpis.obecnosc));
  czesci.push(opisCzasu(wpis.obecnosc?.lastActivityAt ?? wpis.sesja.updatedAt));

  meta.textContent = czesci.join(' · ');
  return meta;
}

function zbudujStan(wpis: WpisSesji): HTMLElement {
  const etykieta = etykietaStanuSesji(wpis.sesja.status);
  const plakietka = document.createElement('span');
  // Odmiana neutralna to plakietka bazowa — biblioteka nie ma jej wariantu.
  plakietka.className =
    etykieta.odmiana === 'neutralna'
      ? 'dn-plakietka'
      : `dn-plakietka dn-plakietka--${etykieta.odmiana}`;
  plakietka.textContent = etykieta.tekst;
  return plakietka;
}

function zbudujPowrot(wpis: WpisSesji, powrot: CzynnoscPowrotu): HTMLButtonElement {
  const przycisk = document.createElement('button');
  przycisk.type = 'button';
  przycisk.className = 'dn-btn dn-btn--zarys dn-btn--sm dn-strona__sesja-powrot';
  przycisk.textContent = 'Wróć do sesji';
  przycisk.addEventListener('click', () => {
    // Zajętość przycisku zamiast wyszarzenia pokazuje pracę w toku.
    if (przycisk.getAttribute('aria-busy') === 'true') return;
    przycisk.setAttribute('aria-busy', 'true');
    void powrot(wpis).finally(() => przycisk.removeAttribute('aria-busy'));
  });
  return przycisk;
}
