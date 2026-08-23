/**
 * Oś etapów budowy produktu — „Nawigacja między etapami budowy produktu od
 * architektury po wdrożenie”, jedyna funkcja operatora okna Product Builder.
 *
 * Oś prowadzi wyłącznie do okien, które moduł zbudował: pozycje przychodzą
 * z pliku składającego moduł, a nie z wykazu zapisanego tutaj, więc dopisanie
 * okna nie wymaga poprawki w drugim miejscu. Kolejność etapów — architektura,
 * warsztaty, wdrożenie — jest kolejnością procesu, nie bramą: naciśnięcie
 * etapu późniejszego jest dozwolone i prowadzi ognisko do jego okna.
 *
 * Oś nie udaje stanu etapu. Stan przychodzi zdarzeniem `apps.build.changed`
 * i rysuje go osobny wykaz okna.
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

/** Jeden etap osi: numer, nazwa okna i zdanie o tym, co się w nim robi. */
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
