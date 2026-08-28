/** Treść dokumentu jako bloki formatowane: napis markdown, bloki i element widoczny, z przejściem odwracalnym między wszystkimi trzema postaciami. */

/** Styl nazwany bloku — to, co w pakiecie biurowym stoi na liście stylów: nagłówek, tekst, cytat, lista, zadanie, kod, tabela, linia albo podział strony. */
export type RodzajBloku =
  | 'naglowek-1'
  | 'naglowek-2'
  | 'naglowek-3'
  | 'tekst'
  | 'cytat'
  | 'lista'
  | 'lista-numerowana'
  | 'zadanie'
  | 'kod'
  | 'tabela'
  | 'linia'
  | 'podzial-strony';

/** Jeden blok treści: styl nazwany rodzaju bloku oraz jego wiersze zapisane bez znaczników składni markdown. */
export interface BlokTresci {
  rodzaj: RodzajBloku;
  /** Wiersze bloku bez znaczników składni; blok pusty ma jeden wiersz pusty. */
  wiersze: string[];
}

/** Styl nazwany wraz z jego nazwą widoczną Operatorowi — pozycja wykazu pokazywanego na liście stylów okna pracy. */
export interface StylNazwany {
  rodzaj: RodzajBloku;
  nazwa: string;
}

/** Znacznik podziału strony — komentarz HTML, bo markdown własnego znacznika podziału nie ma, a przepuszcza komentarz nietknięty. */
export const ZNACZNIK_PODZIALU = '<!-- podział strony -->';

/** Style nazwane oferowane na liście stylów okna pracy, w kolejności, w jakiej Operator je tam widzi i wybiera. */
export const STYLE_NAZWANE: readonly StylNazwany[] = [
  { rodzaj: 'tekst', nazwa: 'Tekst zasadniczy' },
  { rodzaj: 'naglowek-1', nazwa: 'Nagłówek pierwszego stopnia' },
  { rodzaj: 'naglowek-2', nazwa: 'Nagłówek drugiego stopnia' },
  { rodzaj: 'naglowek-3', nazwa: 'Nagłówek trzeciego stopnia' },
  { rodzaj: 'cytat', nazwa: 'Cytat' },
  { rodzaj: 'lista', nazwa: 'Lista wypunktowana' },
  { rodzaj: 'lista-numerowana', nazwa: 'Lista numerowana' },
  { rodzaj: 'zadanie', nazwa: 'Lista zadań' },
  { rodzaj: 'kod', nazwa: 'Blok kodu' },
];

/** Przedrostek wiersza dla stylów, które w składni markdown mają przedrostek: nagłówki, cytat, lista i zadanie. */
const PRZEDROSTKI: Partial<Record<RodzajBloku, string>> = {
  'naglowek-1': '# ',
  'naglowek-2': '## ',
  'naglowek-3': '### ',
  cytat: '> ',
  lista: '- ',
  zadanie: '- [ ] ',
};

/** Czyta napis dokumentu na bloki treści, rozpoznając składnię markdown wiersz po wierszu aż do końca napisu. */
export function czytajBloki(tresc: string): BlokTresci[] {
  const wiersze = tresc.split('\n');
  const bloki: BlokTresci[] = [];
  let i = 0;

  while (i < wiersze.length) {
    const wiersz = wiersze[i] ?? '';
    const przyciety = wiersz.trim();

    if (przyciety === ZNACZNIK_PODZIALU) {
      bloki.push({ rodzaj: 'podzial-strony', wiersze: [''] });
      i += 1;
      continue;
    }
    if (/^(-{3,}|\*{3,}|_{3,})$/u.test(przyciety)) {
      bloki.push({ rodzaj: 'linia', wiersze: [''] });
      i += 1;
      continue;
    }
    if (przyciety.startsWith('```')) {
      const zebrane: string[] = [];
      i += 1;
      while (i < wiersze.length && !(wiersze[i] ?? '').trim().startsWith('```')) {
        zebrane.push(wiersze[i] ?? '');
        i += 1;
      }
      // Klamra zamykająca bywa nieodkryta w dokumencie urwanym — blok kończy się na końcu treści.
      if (i < wiersze.length) i += 1;
      bloki.push({ rodzaj: 'kod', wiersze: zebrane.length === 0 ? [''] : zebrane });
      continue;
    }
    if (przyciety.startsWith('|')) {
      const zebrane: string[] = [];
      while (i < wiersze.length && (wiersze[i] ?? '').trim().startsWith('|')) {
        zebrane.push((wiersze[i] ?? '').trim());
        i += 1;
      }
      bloki.push({ rodzaj: 'tabela', wiersze: zebrane });
      continue;
    }

    const wykaz = rodzajWykazu(wiersz);
    if (wykaz !== null) {
      const zebrane: string[] = [];
      while (i < wiersze.length && rodzajWykazu(wiersze[i] ?? '') === wykaz) {
        zebrane.push(zdejmijPrzedrostek(wiersze[i] ?? '', wykaz));
        i += 1;
      }
      bloki.push({ rodzaj: wykaz, wiersze: zebrane });
      continue;
    }

    const naglowek = /^(#{1,3})\s+(.*)$/u.exec(wiersz);
    if (naglowek !== null) {
      const stopien = (naglowek[1] ?? '#').length;
      bloki.push({
        rodzaj: `naglowek-${stopien}` as RodzajBloku,
        wiersze: [naglowek[2] ?? ''],
      });
      i += 1;
      continue;
    }
    if (wiersz.startsWith('>')) {
      const zebrane: string[] = [];
      while (i < wiersze.length && (wiersze[i] ?? '').startsWith('>')) {
        zebrane.push((wiersze[i] ?? '').replace(/^>\s?/u, ''));
        i += 1;
      }
      bloki.push({ rodzaj: 'cytat', wiersze: zebrane });
      continue;
    }

    // Akapit: wiersze do najbliższego wiersza pustego; wiersz pusty sam z siebie akapitu nie zakłada.
    if (przyciety === '') {
      i += 1;
      continue;
    }
    const akapit: string[] = [];
    while (i < wiersze.length && (wiersze[i] ?? '').trim() !== '' && !zaczynaBlok(wiersze[i] ?? '')) {
      akapit.push(wiersze[i] ?? '');
      i += 1;
    }
    bloki.push({ rodzaj: 'tekst', wiersze: akapit });
  }

  if (bloki.length === 0) bloki.push({ rodzaj: 'tekst', wiersze: [''] });
  return bloki;
}

/** Czy wiersz zaczyna blok innego rodzaju niż akapit — nagłówek, cytat, tabelę, kod, wykaz, linię albo podział strony. */
function zaczynaBlok(wiersz: string): boolean {
  const przyciety = wiersz.trim();
  if (przyciety.startsWith('#') || przyciety.startsWith('>') || przyciety.startsWith('|')) return true;
  if (przyciety.startsWith('```')) return true;
  if (przyciety === ZNACZNIK_PODZIALU) return true;
  if (/^(-{3,}|\*{3,}|_{3,})$/u.test(przyciety)) return true;
  return rodzajWykazu(wiersz) !== null;
}

/** Rodzaj wykazu, którym wiersz się zaczyna: lista, lista numerowana albo zadanie; wartość pusta dla wiersza spoza wykazów. */
function rodzajWykazu(wiersz: string): RodzajBloku | null {
  if (/^\s*-\s\[[ xX]\]\s/u.test(wiersz)) return 'zadanie';
  if (/^\s*[-*+]\s/u.test(wiersz)) return 'lista';
  if (/^\s*\d+\.\s/u.test(wiersz)) return 'lista-numerowana';
  return null;
}

/** Zdejmuje przedrostek wykazu markdown, zostawiając samą treść pozycji wykazu bez znaku wypunktowania albo numeru. */
function zdejmijPrzedrostek(wiersz: string, rodzaj: RodzajBloku): string {
  if (rodzaj === 'zadanie') return wiersz.replace(/^\s*-\s\[[ xX]\]\s/u, '');
  if (rodzaj === 'lista') return wiersz.replace(/^\s*[-*+]\s/u, '');
  return wiersz.replace(/^\s*\d+\.\s/u, '');
}

/** Zapisuje bloki z powrotem na napis dokumentu, łącząc zapis każdego bloku pustym wierszem oddzielającym. */
export function zapiszBloki(bloki: readonly BlokTresci[]): string {
  const czesci: string[] = [];
  for (const blok of bloki) {
    czesci.push(zapiszBlok(blok));
  }
  return czesci.join('\n\n');
}

/** Zapisuje jeden blok składnią markdown właściwą jego rodzajowi: nagłówek, cytat, wykaz, kod, tabelę albo akapit. */
export function zapiszBlok(blok: BlokTresci): string {
  if (blok.rodzaj === 'podzial-strony') return ZNACZNIK_PODZIALU;
  if (blok.rodzaj === 'linia') return '---';
  if (blok.rodzaj === 'kod') return ['```', ...blok.wiersze, '```'].join('\n');
  if (blok.rodzaj === 'tabela') return blok.wiersze.join('\n');
  if (blok.rodzaj === 'lista-numerowana') {
    return blok.wiersze.map((wiersz, numer) => `${numer + 1}. ${wiersz}`).join('\n');
  }
  const przedrostek = PRZEDROSTKI[blok.rodzaj];
  if (przedrostek !== undefined) {
    return blok.wiersze.map((wiersz) => `${przedrostek}${wiersz}`).join('\n');
  }
  return blok.wiersze.join('\n');
}

/** Nazwa stylu widoczna Operatorowi; styl bez wiersza na liście stylów zostaje pokazany pod swoim kodem wewnętrznym. */
export function nazwaStylu(rodzaj: RodzajBloku): string {
  const wiersz = STYLE_NAZWANE.find((styl) => styl.rodzaj === rodzaj);
  if (wiersz !== undefined) return wiersz.nazwa;
  if (rodzaj === 'linia') return 'Linia pozioma';
  if (rodzaj === 'podzial-strony') return 'Podział strony';
  if (rodzaj === 'tabela') return 'Tabela';
  return rodzaj;
}

/* ── Wyrys bloku ───────────────────────────────────────────────────────────── */

/** Nazwa elementu HTML noszącego styl nazwany bloku przy jego wyrysowaniu na powierzchni edycji dokumentu. */
function elementBloku(rodzaj: RodzajBloku): string {
  if (rodzaj === 'naglowek-1') return 'h1';
  if (rodzaj === 'naglowek-2') return 'h2';
  if (rodzaj === 'naglowek-3') return 'h3';
  if (rodzaj === 'cytat') return 'blockquote';
  if (rodzaj === 'lista' || rodzaj === 'zadanie') return 'ul';
  if (rodzaj === 'lista-numerowana') return 'ol';
  if (rodzaj === 'kod') return 'pre';
  if (rodzaj === 'tabela') return 'table';
  if (rodzaj === 'linia' || rodzaj === 'podzial-strony') return 'div';
  return 'p';
}

/** Buduje element widoczny bloku; rodzaj bloku idzie do atrybutu danych, żeby arkusz i odczyt powierzchni mogły go odczytać z powrotem. */
export function wyrysBloku(blok: BlokTresci): HTMLElement {
  const element = document.createElement(elementBloku(blok.rodzaj));
  element.className = 'ms-blok';
  element.dataset['blok'] = blok.rodzaj;

  if (blok.rodzaj === 'linia') {
    element.classList.add('ms-blok--linia');
    element.setAttribute('aria-label', 'Linia pozioma');
    return element;
  }
  if (blok.rodzaj === 'podzial-strony') {
    element.classList.add('ms-blok--podzial');
    element.textContent = 'podział strony';
    return element;
  }
  if (blok.rodzaj === 'kod') {
    element.textContent = blok.wiersze.join('\n');
    return element;
  }
  if (blok.rodzaj === 'tabela') {
    wypelnijTabele(element as HTMLTableElement, blok.wiersze);
    return element;
  }
  if (blok.rodzaj === 'lista' || blok.rodzaj === 'lista-numerowana' || blok.rodzaj === 'zadanie') {
    for (const wiersz of blok.wiersze) {
      const pozycja = document.createElement('li');
      pozycja.append(wyrysInline(wiersz));
      element.append(pozycja);
    }
    return element;
  }

  // Akapit i cytat wielowierszowy: złamania wiersza zostają złamaniami, nie nowymi akapitami.
  blok.wiersze.forEach((wiersz, numer) => {
    if (numer > 0) element.append(document.createElement('br'));
    element.append(wyrysInline(wiersz));
  });
  return element;
}

/** Wypełnia tabelę HTML wierszami zapisanymi składnią markdown, pomijając wiersz rozdzielający nagłówek od treści. */
function wypelnijTabele(tabela: HTMLTableElement, wiersze: readonly string[]): void {
  const komorki = wiersze
    .filter((wiersz) => !/^\|[\s|:-]+\|$/u.test(wiersz))
    .map((wiersz) =>
      wiersz
        .replace(/^\|/u, '')
        .replace(/\|$/u, '')
        .split('|')
        .map((komorka) => komorka.trim()),
    );
  tabela.className = 'ms-blok dn-tabela';
  tabela.dataset['blok'] = 'tabela';
  komorki.forEach((wiersz, numer) => {
    const rzad = document.createElement('tr');
    for (const komorka of wiersz) {
      const pole = document.createElement(numer === 0 ? 'th' : 'td');
      pole.append(wyrysInline(komorka));
      rzad.append(pole);
    }
    tabela.append(rzad);
  });
}

/** Znacznik składni znakowej wraz z elementem HTML, którym markdown pokazuje pogrubienie, kursywę i pozostałe wyróżnienia. */
const ZNACZNIKI_INLINE: readonly { wzorzec: RegExp; element: string }[] = [
  { wzorzec: /\*\*([^*]+)\*\*/u, element: 'strong' },
  { wzorzec: /(?<!\*)\*([^*]+)\*(?!\*)/u, element: 'em' },
  { wzorzec: /~~([^~]+)~~/u, element: 's' },
  { wzorzec: /<u>([^<]+)<\/u>/u, element: 'u' },
  { wzorzec: /`([^`]+)`/u, element: 'code' },
];

/** Buduje treść wiersza z formatowaniem znakowym: pogrubieniem, kursywą, przekreśleniem, kodem, odnośnikiem i podkreśleniem znacznikiem HTML. */
export function wyrysInline(tekst: string): DocumentFragment {
  const wynik = document.createDocumentFragment();
  if (tekst === '') return wynik;

  const odnosnik = /\[([^\]]+)\]\(([^)]*)\)/u.exec(tekst);
  const pierwszy = najblizszyZnacznik(tekst);
  if (odnosnik !== null && (pierwszy === null || (odnosnik.index ?? 0) < pierwszy.pozycja)) {
    wynik.append(wyrysInline(tekst.slice(0, odnosnik.index)));
    const element = document.createElement('a');
    element.href = odnosnik[2] ?? '';
    element.append(wyrysInline(odnosnik[1] ?? ''));
    wynik.append(element);
    wynik.append(wyrysInline(tekst.slice((odnosnik.index ?? 0) + odnosnik[0].length)));
    return wynik;
  }
  if (pierwszy === null) {
    wynik.append(document.createTextNode(tekst));
    return wynik;
  }

  wynik.append(document.createTextNode(tekst.slice(0, pierwszy.pozycja)));
  const element = document.createElement(pierwszy.element);
  element.append(wyrysInline(pierwszy.tresc));
  wynik.append(element);
  wynik.append(wyrysInline(tekst.slice(pierwszy.pozycja + pierwszy.dlugosc)));
  return wynik;
}

/** Pierwszy znacznik znakowy w wierszu spośród pogrubienia, kursywy, przekreślenia, kodu i podkreślenia; wartość pusta dla czystego tekstu. */
function najblizszyZnacznik(
  tekst: string,
): { pozycja: number; dlugosc: number; tresc: string; element: string } | null {
  let najblizszy: { pozycja: number; dlugosc: number; tresc: string; element: string } | null = null;
  for (const znacznik of ZNACZNIKI_INLINE) {
    const trafienie = znacznik.wzorzec.exec(tekst);
    if (trafienie === null) continue;
    const pozycja = trafienie.index;
    if (najblizszy !== null && pozycja >= najblizszy.pozycja) continue;
    najblizszy = {
      pozycja,
      dlugosc: trafienie[0].length,
      tresc: trafienie[1] ?? '',
      element: znacznik.element,
    };
  }
  return najblizszy;
}

/* ── Odczyt powierzchni z powrotem na napis ────────────────────────────────── */

/** Odczytuje powierzchnię edycji z powrotem na napis dokumentu, znosząc elementy dokładane samodzielnie przez przeglądarkę. */
export function serializujPowierzchnie(korzen: HTMLElement): string {
  const czesci: string[] = [];
  for (const dziecko of Array.from(korzen.children)) {
    const blok = odczytajBlok(dziecko as HTMLElement);
    if (blok === null) continue;
    czesci.push(zapiszBlok(blok));
  }
  if (czesci.length === 0) return korzen.textContent ?? '';
  return czesci.join('\n\n');
}

/** Odczytuje jeden element powierzchni edycji na blok treści, rozpoznając jego rodzaj i wiersze składni markdown. */
function odczytajBlok(element: HTMLElement): BlokTresci | null {
  const rodzaj = rodzajElementu(element);
  if (rodzaj === 'podzial-strony' || rodzaj === 'linia') return { rodzaj, wiersze: [''] };
  if (rodzaj === 'kod') return { rodzaj, wiersze: (element.textContent ?? '').split('\n') };
  if (rodzaj === 'tabela') return { rodzaj, wiersze: odczytajTabele(element) };
  if (rodzaj === 'lista' || rodzaj === 'lista-numerowana' || rodzaj === 'zadanie') {
    const pozycje = Array.from(element.querySelectorAll('li')).map((pozycja) =>
      serializujInline(pozycja),
    );
    return { rodzaj, wiersze: pozycje.length === 0 ? [''] : pozycje };
  }
  const wiersze = serializujInline(element).split('\n');
  return { rodzaj, wiersze };
}

/** Rodzaj bloku: odczytany ze znacznika danych elementu, a bez niego wyprowadzony z nazwy elementu HTML. */
function rodzajElementu(element: HTMLElement): RodzajBloku {
  const opisany = element.dataset['blok'];
  if (opisany !== undefined && opisany !== '') return opisany as RodzajBloku;
  const nazwa = element.tagName.toLowerCase();
  if (nazwa === 'h1') return 'naglowek-1';
  if (nazwa === 'h2') return 'naglowek-2';
  if (nazwa === 'h3') return 'naglowek-3';
  if (nazwa === 'blockquote') return 'cytat';
  if (nazwa === 'ol') return 'lista-numerowana';
  if (nazwa === 'ul') return 'lista';
  if (nazwa === 'pre') return 'kod';
  if (nazwa === 'table') return 'tabela';
  if (nazwa === 'hr') return 'linia';
  return 'tekst';
}

/** Odczytuje tabelę HTML na wiersze składni markdown wraz z wierszem rozdzielającym nagłówek od treści tabeli. */
function odczytajTabele(element: HTMLElement): string[] {
  const rzedy = Array.from(element.querySelectorAll('tr'));
  if (rzedy.length === 0) return ['|  |  |', '| --- | --- |'];
  const wiersze: string[] = [];
  rzedy.forEach((rzad, numer) => {
    const komorki = Array.from(rzad.children).map((komorka) =>
      serializujInline(komorka as HTMLElement).replace(/\|/gu, '\\|'),
    );
    wiersze.push(`| ${komorki.join(' | ')} |`);
    if (numer === 0) wiersze.push(`| ${komorki.map(() => '---').join(' | ')} |`);
  });
  return wiersze;
}

/** Znacznik markdown dla elementu znakowego: pogrubienia, kursywy, przekreślenia i kodu; wartość pusta dla elementu bez składni. */
function znacznikElementu(nazwa: string): string {
  if (nazwa === 'strong' || nazwa === 'b') return '**';
  if (nazwa === 'em' || nazwa === 'i') return '*';
  if (nazwa === 's' || nazwa === 'strike' || nazwa === 'del') return '~~';
  if (nazwa === 'code') return '`';
  return '';
}

/** Długość zapisu markdown treści elementu do wskazanego punktu — służy przełożeniu zaznaczenia widoku na zakres znaków dokumentu. */
export function dlugoscDoPunktu(element: HTMLElement, wezel: Node, offset: number): number | null {
  if (element === wezel) {
    // Punkt wskazany na samym elemencie znaczy „przed dzieckiem o tym numerze".
    let dlugosc = 0;
    for (const dziecko of Array.from(element.childNodes).slice(0, offset)) {
      dlugosc += dlugoscWezla(dziecko);
    }
    return dlugosc;
  }
  if (!element.contains(wezel)) return null;

  let dlugosc = 0;
  for (const dziecko of Array.from(element.childNodes)) {
    if (dziecko === wezel) {
      if (dziecko.nodeType === Node.TEXT_NODE) return dlugosc + offset;
      return dlugosc;
    }
    if (dziecko.nodeType === Node.ELEMENT_NODE && dziecko.contains(wezel)) {
      const wewnatrz = dlugoscDoPunktu(dziecko as HTMLElement, wezel, offset);
      const nazwa = (dziecko as HTMLElement).tagName.toLowerCase();
      return dlugosc + dlugoscOtwarcia(nazwa) + (wewnatrz ?? 0);
    }
    dlugosc += dlugoscWezla(dziecko);
  }
  return dlugosc;
}

/** Długość zapisu markdown całego węzła powierzchni, liczona łącznie ze znacznikami składni, nie tylko tekstem widocznym. */
function dlugoscWezla(wezel: Node): number {
  return serializujWezel(wezel).length;
}

/** Długość znacznika otwierającego elementu znakowego — podkreślenia, odnośnika albo znacznika pogrubienia i pozostałych wyróżnień. */
function dlugoscOtwarcia(nazwa: string): number {
  if (nazwa === 'u') return 3;
  if (nazwa === 'a') return 1;
  return znacznikElementu(nazwa).length;
}

/** Zapis markdown jednego węzła powierzchni edycji: tekstu, złamania wiersza albo elementu znakowego ze znacznikiem właściwym. */
function serializujWezel(wezel: Node): string {
  if (wezel.nodeType === Node.TEXT_NODE) return wezel.textContent ?? '';
  if (wezel.nodeType !== Node.ELEMENT_NODE) return '';
  const element = wezel as HTMLElement;
  const nazwa = element.tagName.toLowerCase();
  if (nazwa === 'br') return '\n';
  if (nazwa === 'u') return `<u>${serializujInline(element)}</u>`;
  if (nazwa === 'a') return `[${serializujInline(element)}](${element.getAttribute('href') ?? ''})`;
  const znacznik = znacznikElementu(nazwa);
  return `${znacznik}${serializujInline(element)}${znacznik}`;
}

/** Odczytuje treść znakową elementu z powrotem na składnię markdown, przechodząc rekurencyjnie przez wszystkie węzły potomne. */
export function serializujInline(element: HTMLElement): string {
  return Array.from(element.childNodes).map(serializujWezel).join('');
}
