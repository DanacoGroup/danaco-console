/**
 * Pasmo kart okien roboczych. Karta pasma odpowiada oknu komunikacji stojącemu
 * w karcie sesji — bez tego Operator ma naraz jedno okno i drugie wchodzi na
 * miejsce pierwszego. Znacznik pasma niesie biblioteka Właściciela; ten plik
 * powiela jego kartę wzorcową i wiąże przełączanie oraz zamykanie.
 */

/** Okno pokazywane kartą pasma. */
export interface KartaOkna {
  /** Identyfikator okna komunikacji w rdzeniu. */
  id: string;
  /** Nazwa karty widoczna w paśmie. */
  nazwa: string;
  /** Nazwa modułu okna; pasmo pokazuje ją przy karcie. */
  modul: string;
}

/** Karta wzorcowa zdjęta z pasma przed wyczyszczeniem kart przykładowych. */
let wzor: HTMLElement | null = null;
let lista: HTMLElement | null = null;
let naglowekWykazu: HTMLElement | null = null;
let wykaz: HTMLElement | null = null;
let wzorPozycjiWykazu: HTMLElement | null = null;
let wykazOkien: HTMLElement | null = null;
let wzorPozycjiOkna: HTMLElement | null = null;

/**
 * Zdejmuje z pasma karty przykładowe, zachowując kartę główną i wzór karty
 * modułu. Fałsz znaczy, że pasma nie ma — okno Centrum nie stoi.
 */
export function przygotujPasmo(): boolean {
  lista = document.querySelector<HTMLElement>('.dn-obszar-panel--glowny .dn-karty-lista');
  if (lista === null) return false;
  const karty = [...lista.querySelectorAll<HTMLElement>('.dn-karta-widoku')];
  wzor = karty.find((karta) => karta.dataset.kartaRodzaj === 'modul')?.cloneNode(true) as
    HTMLElement | undefined ?? null;
  for (const karta of karty) {
    if (karta.dataset.kartaRodzaj !== 'centrum') karta.remove();
  }
  wzor?.removeAttribute('data-grupa');
  const menu = document.querySelector<HTMLElement>('#menu-karty-otwarte');
  naglowekWykazu = menu?.querySelector<HTMLElement>('.dn-menu-naglowek') ?? null;
  wykaz = menu;
  const pozycje = [...(menu?.querySelectorAll<HTMLElement>('[data-karta-przelacz]') ?? [])];
  wzorPozycjiWykazu = pozycje.find((pozycja) => pozycja.dataset.kartaPrzelacz !== 'centrum')
    ?.cloneNode(true) as HTMLElement | undefined ?? null;
  for (const pozycja of pozycje) {
    if (pozycja.dataset.kartaPrzelacz !== 'centrum') pozycja.remove();
  }
  /* Okna robocze stoją w menu powłok pasma, osobnym od menu otwartych kart:
     pozycje przełączania powstają z wzoru zdjętego z pozycji przykładowej. */
  wykazOkien = document.querySelector<HTMLElement>('#menu-powloki');
  const przelaczniki = [...(wykazOkien?.querySelectorAll<HTMLElement>('[data-okno-akcja="przelacz"]') ?? [])];
  wzorPozycjiOkna = przelaczniki[0]?.cloneNode(true) as HTMLElement | undefined ?? null;
  for (const pozycja of przelaczniki) pozycja.remove();
  return true;
}

/** Wstawia w menu okna pozycje przełączania okien roboczych i zaznacza bieżące. */
export function ustawOknaRobocze(
  okna: { id: string; nazwa: string }[],
  biezace: string,
): void {
  if (wykazOkien === null || wzorPozycjiOkna === null) return;
  for (const pozycja of wykazOkien.querySelectorAll('[data-okno-akcja="przelacz"]')) pozycja.remove();
  const kotwica = wykazOkien.querySelector('[data-okno-akcja="nowe"]');
  for (const okno of okna) {
    const pozycja = wzorPozycjiOkna.cloneNode(true) as HTMLElement;
    pozycja.dataset.oknoRobocze = okno.id;
    pozycja.dataset.oknoNazwa = okno.nazwa;
    pozycja.setAttribute('aria-current', String(okno.id === biezace));
    // Miara kart pozycji przykładowej opisuje okno, którego nie ma.
    pozycja.querySelector('.dn-meta')?.remove();
    pozycja.querySelector('.skrot')?.remove();
    wpiszNazwe(pozycja, okno.nazwa);
    kotwica?.parentElement?.insertBefore(pozycja, kotwica);
  }
}

/** Wiąże sekcję okien roboczych: przełączenie, nowe okno i zamknięcie bieżącego. */
export function zwiazOknaRobocze(
  przelacz: (id: string) => void,
  nowe: () => void,
  zamknij: () => void,
): void {
  wykazOkien?.addEventListener('click', (zdarzenie) => {
    const cel = zdarzenie.target;
    if (!(cel instanceof Element)) return;
    const pozycja = cel.closest<HTMLElement>('[data-okno-akcja]');
    if (pozycja === null) return;
    zdarzenie.stopPropagation();
    const czynnosc = pozycja.dataset.oknoAkcja ?? '';
    if (czynnosc === 'nowe') nowe();
    else if (czynnosc === 'zamknij') zamknij();
    else if (czynnosc === 'przelacz') przelacz(pozycja.dataset.oknoRobocze ?? '');
  }, true);
}

/** Wstawia w pasmo karty wskazanych okien i zaznacza kartę bieżącą. */
export function ustawKarty(okna: KartaOkna[], biezace: string): void {
  if (lista === null) return;
  for (const karta of lista.querySelectorAll('.dn-karta-widoku')) {
    if ((karta as HTMLElement).dataset.kartaRodzaj !== 'centrum') karta.remove();
  }
  for (const okno of okna) {
    const karta = zbudujKarte(okno);
    if (karta !== null) lista.appendChild(karta);
  }
  odswiezWykaz(okna);
  zaznaczKarte(biezace);
}

/** Zaznacza kartę wskazanego okna; pustka zaznacza kartę główną. */
export function zaznaczKarte(idOkna: string): void {
  if (lista === null) return;
  for (const karta of lista.querySelectorAll<HTMLElement>('.dn-karta-widoku')) {
    const wskazana = idOkna === ''
      ? karta.dataset.kartaRodzaj === 'centrum'
      : karta.dataset.karta === idOkna;
    karta.setAttribute('aria-selected', String(wskazana));
  }
}

/**
 * Wiąże pasmo: naciśnięcie karty przełącza okno, naciśnięcie znaku zamknięcia
 * zamyka je. Nasłuch stoi na paśmie, bo karty powstają i znikają w biegu.
 */
export function zwiazPasmo(
  przelacz: (idOkna: string) => void,
  zamknij: (idOkna: string) => void,
): void {
  if (lista === null) return;
  lista.addEventListener('click', (zdarzenie) => {
    const cel = zdarzenie.target;
    if (!(cel instanceof Element)) return;
    const karta = cel.closest<HTMLElement>('.dn-karta-widoku');
    if (karta === null || karta.dataset.kartaRodzaj === 'centrum') return;
    const idOkna = karta.dataset.karta ?? '';
    if (idOkna === '') return;
    zdarzenie.stopPropagation();
    if (cel.closest('.dn-karta-widoku-zamknij') !== null) {
      zamknij(idOkna);
      return;
    }
    przelacz(idOkna);
  }, true);

  wykaz?.addEventListener('click', (zdarzenie) => {
    const cel = zdarzenie.target;
    if (!(cel instanceof Element)) return;
    const pozycja = cel.closest<HTMLElement>('[data-karta-przelacz]');
    const idOkna = pozycja?.dataset.kartaPrzelacz ?? '';
    if (idOkna === '' || idOkna === 'centrum') return;
    zdarzenie.stopPropagation();
    przelacz(idOkna);
  }, true);
}

/** Zdejmuje kartę zamkniętego okna z pasma. */
export function zdejmijKarte(idOkna: string): void {
  lista?.querySelector(`.dn-karta-widoku[data-karta="${idOkna}"]`)?.remove();
  wykaz?.querySelector(`[data-karta-przelacz="${idOkna}"]`)?.remove();
}

/** Klon karty wzorcowej opisany oknem; brak wzoru znaczy pasmo bez karty modułu. */
function zbudujKarte(okno: KartaOkna): HTMLElement | null {
  if (wzor === null) return null;
  const karta = wzor.cloneNode(true) as HTMLElement;
  karta.dataset.karta = okno.id;
  karta.dataset.modul = okno.modul;
  karta.setAttribute('aria-label', okno.nazwa);
  karta.setAttribute('data-etykietka', okno.nazwa);
  const nazwa = karta.querySelector('.dn-karta-widoku-nazwa');
  if (nazwa !== null) nazwa.textContent = okno.nazwa;
  return karta;
}

/** Odświeża wykaz otwartych kart w menu pasma; liczba w nagłówku liczy karty pasma wraz z główną. */
function odswiezWykaz(okna: KartaOkna[]): void {
  if (naglowekWykazu !== null) {
    naglowekWykazu.textContent = 'Otwarte karty (' + String(okna.length + 1) + ')';
  }
  if (wykaz === null || wzorPozycjiWykazu === null) return;
  for (const pozycja of wykaz.querySelectorAll<HTMLElement>('[data-karta-przelacz]')) {
    if (pozycja.dataset.kartaPrzelacz !== 'centrum') pozycja.remove();
  }
  for (const okno of okna) {
    const pozycja = wzorPozycjiWykazu.cloneNode(true) as HTMLElement;
    pozycja.dataset.kartaPrzelacz = okno.id;
    // Skrót klawiszowy pozycji przykładowej opisuje kolejność, której wydanie nie prowadzi.
    pozycja.querySelector('.skrot')?.remove();
    wpiszNazwe(pozycja, okno.nazwa);
    wykaz.appendChild(pozycja);
  }
}

/** Wpisuje nazwę w węzeł tekstowy pozycji, zostawiając jej ikonę. */
function wpiszNazwe(pozycja: HTMLElement, nazwa: string): void {
  for (const dziecko of pozycja.childNodes) {
    if (dziecko.nodeType !== Node.TEXT_NODE) continue;
    if ((dziecko.nodeValue ?? '').trim() === '') continue;
    dziecko.nodeValue = nazwa;
    return;
  }
}
