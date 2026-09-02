// Pasmo kart okien roboczych: karta pasma jest jedną pracą Operatora, nie
// modułem. Znacznik pasma niesie biblioteka Właściciela.
export interface KartaOkna {
  id: string;
  nazwa: string;
  kodModulu?: string;
  idOknaKomunikacji?: string;
}

let wzor: HTMLElement | null = null;
let lista: HTMLElement | null = null;
let naglowekWykazu: HTMLElement | null = null;
let wykaz: HTMLElement | null = null;
let wzorPozycjiWykazu: HTMLElement | null = null;
let wykazOkien: HTMLElement | null = null;
let wzorPozycjiOkna: HTMLElement | null = null;

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
  // Okna robocze stoją w menu powłok, osobnym od menu otwartych kart.
  wykazOkien = document.querySelector<HTMLElement>('#menu-powloki');
  const przelaczniki = [...(wykazOkien?.querySelectorAll<HTMLElement>('[data-okno-akcja="przelacz"]') ?? [])];
  wzorPozycjiOkna = przelaczniki[0]?.cloneNode(true) as HTMLElement | undefined ?? null;
  for (const pozycja of przelaczniki) pozycja.remove();
  return true;
}

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

export function ustawKarty(karty: KartaOkna[], biezaca: string): void {
  if (lista === null) return;
  for (const karta of lista.querySelectorAll('.dn-karta-widoku')) {
    if ((karta as HTMLElement).dataset.kartaRodzaj !== 'centrum') karta.remove();
  }
  for (const opis of karty) {
    const karta = zbudujKarte(opis);
    if (karta !== null) lista.appendChild(karta);
  }
  odswiezWykaz(karty);
  zaznaczKarte(biezaca);
}

export function zaznaczKarte(idKarty: string): void {
  if (lista === null) return;
  for (const karta of lista.querySelectorAll<HTMLElement>('.dn-karta-widoku')) {
    const wskazana = idKarty === ''
      ? karta.dataset.kartaRodzaj === 'centrum'
      : karta.dataset.karta === idKarty;
    karta.setAttribute('aria-selected', String(wskazana));
  }
}

export function zwiazPasmo(
  przelacz: (idKarty: string) => void,
  zamknij: (idKarty: string) => void,
): void {
  if (lista === null) return;
  lista.addEventListener('click', (zdarzenie) => {
    const cel = zdarzenie.target;
    if (!(cel instanceof Element)) return;
    const karta = cel.closest<HTMLElement>('.dn-karta-widoku');
    if (karta === null || karta.dataset.kartaRodzaj === 'centrum') return;
    const idKarty = karta.dataset.karta ?? '';
    if (idKarty === '') return;
    zdarzenie.stopPropagation();
    if (cel.closest('.dn-karta-widoku-zamknij') !== null) {
      zamknij(idKarty);
      return;
    }
    przelacz(idKarty);
  }, true);

  wykaz?.addEventListener('click', (zdarzenie) => {
    const cel = zdarzenie.target;
    if (!(cel instanceof Element)) return;
    const pozycja = cel.closest<HTMLElement>('[data-karta-przelacz]');
    const idKarty = pozycja?.dataset.kartaPrzelacz ?? '';
    if (idKarty === '' || idKarty === 'centrum') return;
    zdarzenie.stopPropagation();
    przelacz(idKarty);
  }, true);
}

export function nazwijKarte(idKarty: string, nazwa: string): void {
  const karta = lista?.querySelector<HTMLElement>(`.dn-karta-widoku[data-karta="${idKarty}"]`);
  if (karta === undefined || karta === null || nazwa === '') return;
  karta.setAttribute('aria-label', nazwa);
  karta.setAttribute('data-etykietka', nazwa);
  const podpis = karta.querySelector('.dn-karta-widoku-nazwa');
  if (podpis !== null) podpis.textContent = nazwa;
  const pozycja = wykaz?.querySelector<HTMLElement>(`[data-karta-przelacz="${idKarty}"]`);
  if (pozycja !== null && pozycja !== undefined) wpiszNazwe(pozycja, nazwa);
}

export function przypnijKarte(idKarty: string): void {
  const karta = lista?.querySelector<HTMLElement>(`.dn-karta-widoku[data-karta="${idKarty}"]`);
  if (karta === undefined || karta === null) return;
  karta.dataset.przypieta = karta.dataset.przypieta === 'tak' ? 'nie' : 'tak';
}

export function zdejmijKarte(idKarty: string): void {
  lista?.querySelector(`.dn-karta-widoku[data-karta="${idKarty}"]`)?.remove();
  wykaz?.querySelector(`[data-karta-przelacz="${idKarty}"]`)?.remove();
}

function zbudujKarte(opis: KartaOkna): HTMLElement | null {
  if (wzor === null) return null;
  const karta = wzor.cloneNode(true) as HTMLElement;
  karta.dataset.karta = opis.id;
  if (opis.kodModulu !== undefined) karta.dataset.modul = opis.kodModulu;
  if (opis.idOknaKomunikacji !== undefined && opis.idOknaKomunikacji !== '') {
    karta.dataset.oknoKomunikacji = opis.idOknaKomunikacji;
  }
  karta.setAttribute('aria-label', opis.nazwa);
  karta.setAttribute('data-etykietka', opis.nazwa);
  const nazwa = karta.querySelector('.dn-karta-widoku-nazwa');
  if (nazwa !== null) nazwa.textContent = opis.nazwa;
  return karta;
}

function odswiezWykaz(karty: KartaOkna[]): void {
  if (naglowekWykazu !== null) {
    naglowekWykazu.textContent = 'Otwarte karty (' + String(karty.length + 1) + ')';
  }
  if (wykaz === null || wzorPozycjiWykazu === null) return;
  for (const pozycja of wykaz.querySelectorAll<HTMLElement>('[data-karta-przelacz]')) {
    if (pozycja.dataset.kartaPrzelacz !== 'centrum') pozycja.remove();
  }
  for (const opis of karty) {
    const pozycja = wzorPozycjiWykazu.cloneNode(true) as HTMLElement;
    pozycja.dataset.kartaPrzelacz = opis.id;
    // Skrót klawiszowy pozycji przykładowej opisuje kolejność, której wydanie nie prowadzi.
    pozycja.querySelector('.skrot')?.remove();
    wpiszNazwe(pozycja, opis.nazwa);
    wykaz.appendChild(pozycja);
  }
}

function wpiszNazwe(pozycja: HTMLElement, nazwa: string): void {
  for (const dziecko of pozycja.childNodes) {
    if (dziecko.nodeType !== Node.TEXT_NODE) continue;
    if ((dziecko.nodeValue ?? '').trim() === '') continue;
    dziecko.nodeValue = nazwa;
    return;
  }
}
