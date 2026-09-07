// Czynności okien modułów: lista pod przyciskiem belki panelu oraz szuflada
// pytająca o wartości. Czynności stoją w belce, nie nad treścią panelu.

/** Pozycja listy czynności. */
export interface PozycjaCzynnosci {
  kod: string;
  nazwa: string;
  /** Krótkie następstwo czynności, pisane po prawej stronie pozycji. */
  nastepstwo?: string;
}

/** Grupa pozycji nazwana wspólnym nagłówkiem. */
export interface GrupaCzynnosci {
  naglowek: string;
  pozycje: PozycjaCzynnosci[];
}

/** Pole szuflady, o które okno pyta przed wysłaniem czynności. */
export interface PoleSzuflady {
  klucz: string;
  etykieta: string;
  /** Podpowiedź w pustej kontrolce; nie zastępuje etykiety. */
  podpowiedz?: string;
  wartosc?: string;
  /** Wykaz wartości do wyboru; puste daje pole wpisywane. */
  wybor?: ReadonlyArray<readonly [string, string]>;
  obszerne?: boolean;
  /** Czy bez tego pola czynność nie ma sensu. */
  wymagane?: boolean;
}

/** Treść szuflady pytającej o wartości czynności. */
export interface PytanieSzuflady {
  tytul: string;
  /** Zdanie o tym, co czynność zrobi; stoi nad polami. */
  opis?: string;
  pola: PoleSzuflady[];
  /** Napis przycisku wykonującego; mówi, co się stanie. */
  wykonanie: string;
  /** Zdanie o następstwie nieodwracalnym; wymusza drugie naciśnięcie. */
  nieodwracalne?: string;
}

const IKONA_CZYNNOSCI = 'M13 2 4 14h6l-1 8 9-12h-6l1-8Z';

function przestrzenNazw(): string {
  return 'http://www.w3.org/2000/svg';
}

function znak(dokument: Document, sciezka: string): SVGElement {
  const rysunek = dokument.createElementNS(przestrzenNazw(), 'svg') as SVGSVGElement;
  rysunek.setAttribute('viewBox', '0 0 24 24');
  rysunek.setAttribute('fill', 'none');
  rysunek.setAttribute('stroke', 'currentColor');
  rysunek.setAttribute('stroke-width', '1.75');
  rysunek.setAttribute('stroke-linecap', 'round');
  rysunek.setAttribute('stroke-linejoin', 'round');
  rysunek.setAttribute('aria-hidden', 'true');
  const kreska = dokument.createElementNS(przestrzenNazw(), 'path') as SVGPathElement;
  kreska.setAttribute('d', sciezka);
  rysunek.appendChild(kreska);
  return rysunek;
}

function belkaCzynnosci(panel: Element): Element | null {
  const belka = panel.querySelector('.sta-okno-belka');
  if (belka === null) return null;
  const stojace = belka.querySelector('.sta-okno-akcje');
  if (stojace !== null) return stojace;
  const nowe = panel.ownerDocument.createElement('span');
  nowe.className = 'sta-okno-akcje';
  belka.appendChild(nowe);
  return nowe;
}

/*
dolozCzynnosciPanelu wiesza listę czynności na belce wskazanego panelu.

Zwraca zdjęcie listy: mechanika `data-menu` przenosi rozwiniętą listę do
warstwy okna, więc zamknięcie karty musi ją stamtąd zabrać.
*/
export function dolozCzynnosciPanelu(
  korzen: Element,
  panelKod: string,
  etykieta: string,
  grupy: ReadonlyArray<GrupaCzynnosci>,
  wykonaj: (kod: string) => void,
  przy: AddEventListenerOptions,
): (() => void) | null {
  const panel = korzen.querySelector(`#${panelKod}`);
  if (panel === null) return null;
  const akcje = belkaCzynnosci(panel);
  if (akcje === null) return null;
  const dokument = korzen.ownerDocument;
  const oznaczenie = `czynnosci-${korzen.getAttribute('data-karta') ?? 'karta'}-${panelKod}`;
  if (dokument.getElementById(oznaczenie) !== null) return null;

  const gniazdo = dokument.createElement('div');
  gniazdo.className = 'sta-menu';
  const przycisk = dokument.createElement('button');
  przycisk.type = 'button';
  przycisk.className = 'dn-btn-ikona';
  przycisk.dataset.menu = oznaczenie;
  przycisk.dataset.etykietka = etykieta;
  przycisk.setAttribute('aria-expanded', 'false');
  przycisk.setAttribute('aria-label', etykieta);
  przycisk.appendChild(znak(dokument, IKONA_CZYNNOSCI));

  const lista = dokument.createElement('div');
  lista.className = 'sta-menu-tresc';
  lista.id = oznaczenie;
  lista.dataset.menuTresc = '';
  lista.dataset.otwarte = 'nie';
  lista.dataset.kotwica = 'prawo';
  lista.setAttribute('role', 'menu');
  for (const [numer, grupa] of grupy.entries()) {
    if (numer > 0) {
      const przerwa = dokument.createElement('div');
      przerwa.className = 'sta-menu-sep';
      lista.appendChild(przerwa);
    }
    const naglowek = dokument.createElement('div');
    naglowek.className = 'sta-menu-etyk';
    naglowek.textContent = grupa.naglowek;
    lista.appendChild(naglowek);
    for (const pozycja of grupa.pozycje) lista.appendChild(wierszPozycji(dokument, pozycja));
  }
  gniazdo.append(przycisk, lista);

  const zamkniecie = akcje.querySelector('[data-okno-zamknij]');
  akcje.insertBefore(gniazdo, zamkniecie);

  /* Nasłuch stoi na samej liście, nie na korzeniu karty: mechanika menu
     przenosi rozwiniętą listę do warstwy okna, więc kliknięcie w pozycję nie
     przechodzi już przez drzewo karty. */
  lista.addEventListener('click', (zdarzenie) => {
    const cel = zdarzenie.target;
    if (!(cel instanceof Element)) return;
    const pozycja = cel.closest<HTMLElement>('[data-czynnosc]');
    const kod = pozycja?.dataset.czynnosc ?? '';
    if (kod === '') return;
    zdarzenie.stopPropagation();
    lista.dataset.otwarte = 'nie';
    przycisk.setAttribute('aria-expanded', 'false');
    wykonaj(kod);
  }, przy);

  return () => {
    lista.remove();
    gniazdo.remove();
  };
}

function wierszPozycji(dokument: Document, pozycja: PozycjaCzynnosci): HTMLElement {
  const wiersz = dokument.createElement('button');
  wiersz.type = 'button';
  wiersz.className = 'sta-menu-poz';
  wiersz.setAttribute('role', 'menuitem');
  wiersz.dataset.czynnosc = pozycja.kod;
  wiersz.append(dokument.createTextNode(pozycja.nazwa));
  if (pozycja.nastepstwo !== undefined) {
    const dopisek = dokument.createElement('span');
    dopisek.className = 'skrot';
    dopisek.textContent = pozycja.nastepstwo;
    wiersz.appendChild(dopisek);
  }
  return wiersz;
}

/*
zapytajWSzufladzie otwiera szufladę panelu i czeka na wartości od Operatora.

Osobne pole na każdą wartość, bo jedno pole na wszystko byłoby wierszem
poleceń w przebraniu okna — składni nie widać w oknie nigdzie.
*/
export async function zapytajWSzufladzie(
  korzen: Element,
  panelKod: string,
  pytanie: PytanieSzuflady,
): Promise<Record<string, string> | null> {
  const panel = korzen.querySelector(`#${panelKod}`);
  if (panel === null) return null;
  const dokument = korzen.ownerDocument;
  panel.querySelector('.sta-konfig[data-czynnosc-szuflada]')?.remove();

  const szuflada = dokument.createElement('div');
  szuflada.className = 'sta-konfig';
  szuflada.dataset.czynnoscSzuflada = '';
  const belka = dokument.createElement('div');
  belka.className = 'sta-konfig-belka';
  const tytul = dokument.createElement('b');
  tytul.textContent = pytanie.tytul;
  const zamknij = dokument.createElement('button');
  zamknij.type = 'button';
  zamknij.className = 'dn-btn-ikona';
  zamknij.style.marginLeft = 'auto';
  zamknij.setAttribute('aria-label', 'Zamknij bez wykonania');
  zamknij.appendChild(znak(dokument, 'M18 6 6 18M6 6l12 12'));
  belka.append(tytul, zamknij);

  const sekcja = dokument.createElement('details');
  sekcja.className = 'sta-konfig-sekcja';
  sekcja.open = true;
  const podpis = dokument.createElement('summary');
  podpis.textContent = 'Wartości czynności';
  sekcja.appendChild(podpis);
  if (pytanie.opis !== undefined) {
    const zdanie = dokument.createElement('p');
    zdanie.className = 'dn-meta';
    zdanie.textContent = pytanie.opis;
    sekcja.appendChild(zdanie);
  }
  const kontrolki = new Map<string, HTMLInputElement | HTMLSelectElement | HTMLTextAreaElement>();
  for (const pole of pytanie.pola) {
    const opakowanie = dokument.createElement('div');
    opakowanie.className = 'sta-konfig-pole';
    const nazwa = dokument.createElement('label');
    nazwa.textContent = pole.etykieta;
    nazwa.htmlFor = `${panelKod}-${pole.klucz}`;
    const kontrolka = zbudujKontrolke(dokument, pole);
    kontrolka.id = nazwa.htmlFor;
    opakowanie.append(nazwa, kontrolka);
    sekcja.appendChild(opakowanie);
    kontrolki.set(pole.klucz, kontrolka);
  }

  const stopka = dokument.createElement('div');
  stopka.className = 'sta-konfig-sekcja';
  const ostrzezenie = dokument.createElement('p');
  ostrzezenie.className = 'dn-meta';
  ostrzezenie.hidden = true;
  const wykonanie = dokument.createElement('button');
  wykonanie.type = 'button';
  wykonanie.className = 'dn-btn dn-btn--sygnal';
  wykonanie.textContent = pytanie.wykonanie;
  stopka.append(ostrzezenie, wykonanie);

  szuflada.append(belka, sekcja, stopka);
  panel.appendChild(szuflada);
  kontrolki.values().next().value?.focus();

  return new Promise((rozwiaz) => {
    const domknij = (wynik: Record<string, string> | null): void => {
      szuflada.remove();
      rozwiaz(wynik);
    };
    zamknij.addEventListener('click', () => {
      domknij(null);
    });
    szuflada.addEventListener('keydown', (zdarzenie) => {
      if (zdarzenie.key === 'Escape') domknij(null);
    });
    wykonanie.addEventListener('click', () => {
      const brakujace = pytanie.pola
        .filter((pole) => pole.wymagane === true && (kontrolki.get(pole.klucz)?.value ?? '') === '')
        .map((pole) => pole.etykieta);
      for (const pole of pytanie.pola) {
        kontrolki.get(pole.klucz)?.setAttribute('aria-invalid',
          String(brakujace.includes(pole.etykieta)));
      }
      if (brakujace.length > 0) {
        ostrzezenie.hidden = false;
        ostrzezenie.textContent = `Bez wartości nie da się wykonać: ${brakujace.join(', ')}.`;
        return;
      }
      if (pytanie.nieodwracalne !== undefined && ostrzezenie.dataset.odslonione !== 'tak') {
        ostrzezenie.dataset.odslonione = 'tak';
        ostrzezenie.hidden = false;
        ostrzezenie.textContent = pytanie.nieodwracalne;
        wykonanie.textContent = `${pytanie.wykonanie} — potwierdź`;
        return;
      }
      const wartosci: Record<string, string> = {};
      for (const [klucz, kontrolka] of kontrolki) wartosci[klucz] = kontrolka.value.trim();
      domknij(wartosci);
    });
  });
}

function zbudujKontrolke(
  dokument: Document,
  pole: PoleSzuflady,
): HTMLInputElement | HTMLSelectElement | HTMLTextAreaElement {
  if (pole.wybor !== undefined) {
    const wybor = dokument.createElement('select');
    wybor.className = 'dn-pole-kontrolka';
    for (const [wartosc, nazwa] of pole.wybor) {
      const pozycja = dokument.createElement('option');
      pozycja.value = wartosc;
      pozycja.textContent = nazwa;
      wybor.appendChild(pozycja);
    }
    wybor.value = pole.wartosc ?? pole.wybor[0]?.[0] ?? '';
    return wybor;
  }
  if (pole.obszerne === true) {
    const obszar = dokument.createElement('textarea');
    obszar.className = 'dn-pole-kontrolka';
    obszar.rows = 4;
    obszar.value = pole.wartosc ?? '';
    if (pole.podpowiedz !== undefined) obszar.placeholder = pole.podpowiedz;
    return obszar;
  }
  const wpis = dokument.createElement('input');
  wpis.type = 'text';
  wpis.className = 'dn-pole-kontrolka';
  wpis.value = pole.wartosc ?? '';
  if (pole.podpowiedz !== undefined) wpis.placeholder = pole.podpowiedz;
  return wpis;
}
