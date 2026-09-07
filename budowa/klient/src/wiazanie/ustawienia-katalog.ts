// Wspólny kształt katalogu w oknie Ustawień: pozycja z nagłówkiem, przyciskiem
// i tabelą, oraz komórki jej wierszy.

export interface OpisGrupy {
  kod: string;
  tytul: string;
  opis: string;
  napisPrzycisku: string;
  kolumny: readonly string[];
  poleWyboru?: string;
  napisWyboru?: string;
}

/* Katalogi Ustawień stoją w sekcjach, których prototyp nie zbudował, więc
   powstają z tego samego słownika, co reszta okna: pozycja, opis, ramka
   z tabelą. */
export function grupaKatalogu(dokument: Document, opis: OpisGrupy): HTMLElement {
  const pozycja = dokument.createElement('div');
  pozycja.className = 'us-pozycja';
  pozycja.append(glowaGrupy(dokument, opis), zdanieGrupy(dokument, opis.opis),
    ramkaTabeli(dokument, opis));
  return pozycja;
}

function glowaGrupy(dokument: Document, opis: OpisGrupy): HTMLElement {
  const glowa = dokument.createElement('div');
  glowa.className = 'us-pozycja-glowa';
  const etykieta = dokument.createElement('span');
  etykieta.className = 'us-etykieta';
  etykieta.textContent = opis.tytul;
  const prawa = dokument.createElement('span');
  prawa.className = 'us-prawa';
  if (opis.poleWyboru !== undefined) {
    const wybor = dokument.createElement('select');
    wybor.className = 'dn-pole-kontrolka';
    wybor.id = `${opis.kod}-${opis.poleWyboru}`;
    const pokaz = dokument.createElement('button');
    pokaz.type = 'button';
    pokaz.className = 'dn-btn dn-btn--duch dn-btn--sm';
    pokaz.id = `${opis.kod}-pokaz`;
    pokaz.textContent = opis.napisWyboru ?? 'Pokaż';
    prawa.append(wybor, pokaz);
  }
  const dodanie = dokument.createElement('button');
  dodanie.type = 'button';
  dodanie.className = 'dn-btn dn-btn--atrament dn-btn--sm';
  dodanie.id = `${opis.kod}-dodaj`;
  dodanie.textContent = opis.napisPrzycisku;
  prawa.appendChild(dodanie);
  glowa.append(etykieta, prawa);
  return glowa;
}

function zdanieGrupy(dokument: Document, tresc: string): HTMLElement {
  const zdanie = dokument.createElement('p');
  zdanie.className = 'us-opis';
  zdanie.textContent = tresc;
  return zdanie;
}

function ramkaTabeli(dokument: Document, opis: OpisGrupy): HTMLElement {
  const ramka = dokument.createElement('div');
  ramka.className = 'us-tabela-ramka';
  const tabela = dokument.createElement('table');
  tabela.className = 'dn-tabela us-tabela';
  const glowica = dokument.createElement('thead');
  const wiersz = dokument.createElement('tr');
  for (const nazwa of opis.kolumny) {
    const komorkaGlowicy = dokument.createElement('th');
    komorkaGlowicy.scope = 'col';
    komorkaGlowicy.textContent = nazwa;
    wiersz.appendChild(komorkaGlowicy);
  }
  glowica.appendChild(wiersz);
  const cialo = dokument.createElement('tbody');
  cialo.id = `${opis.kod}-cialo`;
  tabela.append(glowica, cialo);
  ramka.appendChild(tabela);
  return ramka;
}

export function zdanieWiersza(cialo: Element, ile: number, zdanie: string): HTMLElement {
  const wiersz = cialo.ownerDocument.createElement('tr');
  const komorkaZdania = cialo.ownerDocument.createElement('td');
  komorkaZdania.colSpan = ile;
  komorkaZdania.className = 'us-opis--drobny';
  komorkaZdania.textContent = zdanie;
  wiersz.appendChild(komorkaZdania);
  return wiersz;
}

export function komorka(dokument: Document, tresc: string, klasa = ''): HTMLElement {
  const wezel = dokument.createElement('td');
  if (klasa !== '') wezel.className = klasa;
  wezel.textContent = tresc;
  return wezel;
}

export function komorkaCzynnosci(
  dokument: Document,
  czynnosci: ReadonlyArray<readonly [string, string]>,
  klucz: string,
): HTMLElement {
  const wezel = dokument.createElement('td');
  const gniazdo = dokument.createElement('span');
  gniazdo.className = 'us-akcje-komorka';
  for (const [kod, napis] of czynnosci) {
    const przycisk = dokument.createElement('button');
    przycisk.type = 'button';
    przycisk.className = 'dn-btn dn-btn--duch dn-btn--sm';
    przycisk.setAttribute(klucz, kod);
    przycisk.textContent = napis;
    gniazdo.appendChild(przycisk);
  }
  wezel.appendChild(gniazdo);
  return wezel;
}

export function wypelnijWybor(
  wybor: HTMLSelectElement,
  pozycje: ReadonlyArray<readonly [string, string]>,
  pierwsza: string,
): void {
  const stojaca = wybor.value;
  wybor.replaceChildren();
  for (const [wartosc, nazwa] of [['', pierwsza] as const, ...pozycje]) {
    const pozycja = wybor.ownerDocument.createElement('option');
    pozycja.value = wartosc;
    pozycja.textContent = nazwa;
    wybor.appendChild(pozycja);
  }
  wybor.value = pozycje.some((pozycja) => pozycja[0] === stojaca) ? stojaca : '';
}
