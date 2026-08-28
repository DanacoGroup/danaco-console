/**
 * Dokument zastępczy dla sprawdzianów ramy: fragment `<body>` prawdziwego
 * `index.html` jako drzewo węzłów — tylko ten wycinek zachowania DOM,
 * którego montaż ramy dotyka, bez zależności zewnętrznej. Węzeł nie
 * implementuje typów `Document`/`HTMLElement`; granicę przenosi
 * rzutowanie w `zbudujDokument`.
 */

const TEKST = '#tekst';
const PUSTE = new Set(['meta', 'link', 'br', 'img', 'input', 'hr']);

/** Zdarzenie zastępcze: tylko rodzaj i cel, jedyne pola, po które sięga rama. */
export interface ZdarzenieZastepcze {
  type: string;
  target: WezelZastepczy;
}

function pasuje(wezel: WezelZastepczy, selektor: string): boolean {
  if (selektor.startsWith('[') && selektor.endsWith(']')) {
    return wezel.atrybuty.has(selektor.slice(1, -1));
  }
  if (selektor.startsWith('.')) {
    return wezel.className.split(/\s+/).includes(selektor.slice(1));
  }
  return false;
}

function doKebab(nazwa: string): string {
  return nazwa.replace(/[A-Z]/g, (litera) => `-${litera.toLowerCase()}`);
}

/** Węzeł drzewa zastępczego: element albo węzeł tekstu, rozróżniane znacznikiem `#tekst`. */
export class WezelZastepczy {
  readonly tag: string;
  readonly atrybuty = new Map<string, string>();
  dzieci: WezelZastepczy[] = [];
  rodzic: WezelZastepczy | null = null;
  wartoscTekstu = '';
  private readonly sluchacze = new Map<string, Set<(zdarzenie: ZdarzenieZastepcze) => void>>();

  constructor(tag: string) {
    this.tag = tag;
  }

  get className(): string {
    return this.atrybuty.get('class') ?? '';
  }

  set className(wartosc: string) {
    this.atrybuty.set('class', wartosc);
  }

  get hidden(): boolean {
    return this.atrybuty.has('hidden');
  }

  set hidden(wartosc: boolean) {
    if (wartosc) this.atrybuty.set('hidden', '');
    else this.atrybuty.delete('hidden');
  }

  get dataset(): Record<string, string | undefined> {
    const wezel = this;
    return new Proxy({} as Record<string, string | undefined>, {
      get(_cel, klucz: string): string | undefined {
        return wezel.atrybuty.get(`data-${doKebab(klucz)}`);
      },
    });
  }

  get textContent(): string {
    return this.dzieci.map((d) => (d.tag === TEKST ? d.wartoscTekstu : d.textContent)).join('');
  }

  set textContent(wartosc: string) {
    this.dzieci = [];
    const wezelTekstu = new WezelZastepczy(TEKST);
    wezelTekstu.wartoscTekstu = wartosc;
    this.appendChild(wezelTekstu);
  }

  set innerHTML(zrodlo: string) {
    const sparsowany = analizujZnacznik(zrodlo);
    this.dzieci = [];
    for (const dziecko of sparsowany.dzieci) this.appendChild(dziecko);
  }

  /** Zastępuje `DocumentFragment` szablonu — sam węzeł niesie swoją zawartość. */
  get content(): WezelZastepczy {
    return this;
  }

  get firstElementChild(): WezelZastepczy | null {
    return this.dzieci.find((d) => d.tag !== TEKST) ?? null;
  }

  setAttribute(nazwa: string, wartosc: string): void {
    this.atrybuty.set(nazwa, wartosc);
  }

  getAttribute(nazwa: string): string | null {
    return this.atrybuty.get(nazwa) ?? null;
  }

  removeAttribute(nazwa: string): void {
    this.atrybuty.delete(nazwa);
  }

  appendChild(dziecko: WezelZastepczy): WezelZastepczy {
    dziecko.rodzic = this;
    this.dzieci.push(dziecko);
    return dziecko;
  }

  replaceChildren(...nowe: WezelZastepczy[]): void {
    this.dzieci = [];
    for (const w of nowe) this.appendChild(w);
  }

  replaceWith(...nowe: WezelZastepczy[]): void {
    if (this.rodzic === null) return;
    const indeks = this.rodzic.dzieci.indexOf(this);
    if (indeks === -1) return;
    this.rodzic.dzieci.splice(indeks, 1, ...nowe);
    for (const w of nowe) w.rodzic = this.rodzic;
  }

  addEventListener(typ: string, sluchacz: (zdarzenie: ZdarzenieZastepcze) => void): void {
    const zbior = this.sluchacze.get(typ) ?? new Set();
    zbior.add(sluchacz);
    this.sluchacze.set(typ, zbior);
  }

  removeEventListener(typ: string, sluchacz: (zdarzenie: ZdarzenieZastepcze) => void): void {
    this.sluchacze.get(typ)?.delete(sluchacz);
  }

  /** Rozgłasza zdarzenie od siebie w górę drzewa — ten sam kierunek, co bąbelkowanie klikniętego przycisku do delegata na przodku. */
  dispatchEvent(zdarzenie: ZdarzenieZastepcze): void {
    let wezel: WezelZastepczy | null = this;
    while (wezel !== null) {
      for (const sluchacz of wezel.sluchacze.get(zdarzenie.type) ?? []) sluchacz(zdarzenie);
      wezel = wezel.rodzic;
    }
  }

  querySelectorAll(selektor: string): WezelZastepczy[] {
    const wynik: WezelZastepczy[] = [];
    const przeszukaj = (wezel: WezelZastepczy): void => {
      for (const dziecko of wezel.dzieci) {
        if (dziecko.tag !== TEKST && pasuje(dziecko, selektor)) wynik.push(dziecko);
        przeszukaj(dziecko);
      }
    };
    przeszukaj(this);
    return wynik;
  }

  querySelector(selektor: string): WezelZastepczy | null {
    return this.querySelectorAll(selektor)[0] ?? null;
  }

  closest(selektor: string): WezelZastepczy | null {
    let wezel: WezelZastepczy | null = this;
    while (wezel !== null) {
      if (pasuje(wezel, selektor)) return wezel;
      wezel = wezel.rodzic;
    }
    return null;
  }
}

function analizujAtrybuty(fragment: string): Map<string, string> {
  const atrybuty = new Map<string, string>();
  let i = 0;
  while (i < fragment.length) {
    while (i < fragment.length && /\s/.test(fragment[i]!)) i += 1;
    if (i >= fragment.length) break;
    let nazwa = '';
    while (i < fragment.length && !/[\s=]/.test(fragment[i]!)) {
      nazwa += fragment[i];
      i += 1;
    }
    if (nazwa.length === 0) {
      i += 1;
      continue;
    }
    while (i < fragment.length && /\s/.test(fragment[i]!)) i += 1;
    if (fragment[i] === '=') {
      i += 1;
      while (i < fragment.length && /\s/.test(fragment[i]!)) i += 1;
      const cudzyslow = fragment[i];
      let wartosc = '';
      if (cudzyslow === '"' || cudzyslow === "'") {
        i += 1;
        while (i < fragment.length && fragment[i] !== cudzyslow) {
          wartosc += fragment[i];
          i += 1;
        }
        i += 1;
      } else {
        while (i < fragment.length && !/\s/.test(fragment[i]!)) {
          wartosc += fragment[i];
          i += 1;
        }
      }
      atrybuty.set(nazwa, wartosc);
    } else {
      atrybuty.set(nazwa, '');
    }
  }
  return atrybuty;
}

/** Parsuje znacznik HTML/SVG na drzewo węzłów zastępczych, zwracając ich wspólnego rodzica. */
function analizujZnacznik(zrodlo: string): WezelZastepczy {
  const pojemnik = new WezelZastepczy('#fragment');
  const stos: WezelZastepczy[] = [pojemnik];
  let i = 0;
  while (i < zrodlo.length) {
    if (zrodlo.startsWith('<!--', i)) {
      const koniec = zrodlo.indexOf('-->', i);
      i = koniec === -1 ? zrodlo.length : koniec + 3;
      continue;
    }
    if (zrodlo[i] === '<' && zrodlo[i + 1] === '!') {
      const koniec = zrodlo.indexOf('>', i);
      i = koniec === -1 ? zrodlo.length : koniec + 1;
      continue;
    }
    if (zrodlo[i] === '<' && zrodlo[i + 1] === '/') {
      const koniec = zrodlo.indexOf('>', i);
      if (stos.length > 1) stos.pop();
      i = koniec === -1 ? zrodlo.length : koniec + 1;
      continue;
    }
    if (zrodlo[i] === '<') {
      const koniec = zrodlo.indexOf('>', i);
      const surowy = zrodlo.slice(i + 1, koniec === -1 ? zrodlo.length : koniec);
      const samozamykajacy = surowy.endsWith('/');
      const tresc = samozamykajacy ? surowy.slice(0, -1) : surowy;
      const dopasowanie = /^[a-zA-Z][a-zA-Z0-9-]*/.exec(tresc);
      const nazwa = (dopasowanie?.[0] ?? '').toLowerCase();
      const wezel = new WezelZastepczy(nazwa);
      for (const [klucz, wartosc] of analizujAtrybuty(tresc.slice(nazwa.length))) {
        wezel.setAttribute(klucz, wartosc);
      }
      stos[stos.length - 1]!.appendChild(wezel);
      if (!samozamykajacy && !PUSTE.has(nazwa)) stos.push(wezel);
      i = koniec === -1 ? zrodlo.length : koniec + 1;
      continue;
    }
    const nastepny = zrodlo.indexOf('<', i);
    const fragmentTekstu = zrodlo.slice(i, nastepny === -1 ? zrodlo.length : nastepny);
    if (fragmentTekstu.trim().length > 0) {
      const wezelTekstu = new WezelZastepczy(TEKST);
      wezelTekstu.wartoscTekstu = fragmentTekstu;
      stos[stos.length - 1]!.appendChild(wezelTekstu);
    }
    i = nastepny === -1 ? zrodlo.length : nastepny;
  }
  return pojemnik;
}

/** Dokument zastępczy: tylko te cztery metody `document`, po które sięga montaż ramy i wejścia do niej. */
export class DokumentZastepczy {
  readonly documentElement = new WezelZastepczy('html');
  private readonly pojemnik: WezelZastepczy;

  constructor(pojemnik: WezelZastepczy) {
    this.pojemnik = pojemnik;
  }

  createElement(znacznik: string): WezelZastepczy {
    return new WezelZastepczy(znacznik);
  }

  createTextNode(wartosc: string): WezelZastepczy {
    const wezel = new WezelZastepczy(TEKST);
    wezel.wartoscTekstu = wartosc;
    return wezel;
  }

  querySelector(selektor: string): WezelZastepczy | null {
    return this.pojemnik.querySelector(selektor);
  }

  querySelectorAll(selektor: string): WezelZastepczy[] {
    return this.pojemnik.querySelectorAll(selektor);
  }
}

/**
 * Buduje dokument zastępczy z zawartości `<body>` prawdziwego `index.html`,
 * wycinając ją znalezieniem znaczników — reszta dokumentu (nagłówek, arkusze)
 * nie niesie nic, czego dotyka montaż ramy. Granicę z typami `Document`
 * przenosi wywołujący, rzutując w miejscu, gdzie typ tego wymaga.
 */
export function zbudujDokument(zrodloHtml: string): DokumentZastepczy {
  const poczatekCiala = zrodloHtml.indexOf('>', zrodloHtml.indexOf('<body')) + 1;
  const koniecCiala = zrodloHtml.indexOf('</body>');
  const pojemnik = analizujZnacznik(zrodloHtml.slice(poczatekCiala, koniecCiala));
  return new DokumentZastepczy(pojemnik);
}
