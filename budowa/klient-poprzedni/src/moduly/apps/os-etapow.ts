/**
 * Oś etapów budowy produktu — nawigacja między etapami od architektury po
 * wdrożenie, jedyna funkcja operatora okna Product Builder. Pozycje przychodzą
 * z pliku składającego moduł, a nie z wykazu zapisanego tutaj.
 */
export interface PozycjaOsi {
  /** Kod okna w katalogu rdzenia (`okno_operacyjne.kod`). */
  kod: string;
  /** Nazwa etapu widoczna na osi. */
  tytul: string;
  /** Zdanie mówiące, co w tym etapie robi Operator. */
  opis: string;
}

export interface OsEtapow {
  element: HTMLElement;
}

export function utworzOsEtapow(
  pozycje: readonly PozycjaOsi[],
  przejdz: (kodOkna: string) => void,
): OsEtapow {
  const lista = document.createElement('ol');
  lista.className = 'mp-os';
  lista.append(...pozycje.map((pozycja, numer) => krok(pozycja, numer + 1, przejdz)));

  const element = document.createElement('nav');
  element.className = 'mp-os__powloka';
  element.setAttribute('aria-label', 'Etapy budowy produktu');
  element.append(lista);
  return { element };
}

/**
 * Jeden etap osi: numer porządkowy, nazwa okna oraz zdanie o tym, co się w nim
 * robi, złożone w przycisk prowadzący do okna etapu.
 */
function krok(pozycja: PozycjaOsi, numer: number, przejdz: (kodOkna: string) => void): HTMLElement {
  const znak = document.createElement('span');
  znak.className = 'mp-os__znak';
  znak.textContent = String(numer);

  const tytul = document.createElement('span');
  tytul.className = 'mp-os__tytul';
  tytul.textContent = pozycja.tytul;

  const opis = document.createElement('span');
  opis.className = 'mp-os__opis';
  opis.textContent = pozycja.opis;

  const przycisk = document.createElement('button');
  przycisk.type = 'button';
  przycisk.className = 'mp-os__wybor';
  przycisk.dataset['etap'] = pozycja.kod;
  przycisk.append(znak, tytul, opis);
  przycisk.addEventListener('click', () => przejdz(pozycja.kod));

  const element = document.createElement('li');
  element.className = 'mp-os__krok';
  element.append(przycisk);
  return element;
}
