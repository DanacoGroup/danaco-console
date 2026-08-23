/**
 * Symbole zastępcze i znaczniki formatu — rozpoznanie w tekście źródłowym oraz
 * przekład ich stylu między platformami.
 *
 * Obie czynności liczy okno, nie rdzeń, i jest to rozstrzygnięcie wymuszone
 * kontraktem: kontrola jakości rdzenia zgłasza niezgodność symbolu zastępczego
 * dopiero w gotowym panelu, a komendy zamieniającej styl zmiennych nie ma
 * w ogóle. Okno liczy więc to, co da się policzyć z samego tekstu, i nie
 * przypisuje wyniku rdzeniowi.
 *
 * Wykaz stylów jest zamknięty i wzięty z opracowania modułu: `%s`, `{0}`,
 * `{name}`, `{{name}}`, `${name}` oraz znacznik formatu `<b>`. Styl spoza
 * wykazu nie jest rozpoznawany i okno tego nie ukrywa — wynik mówi, ile
 * wystąpień znalazło, a nie że znalazło wszystkie.
 *
 * Przekład stylu zachowuje kolejność wystąpień i nazwy, gdy styl źródłowy je
 * niesie. Gdy styl źródłowy nazwy nie ma (`%s`, `{0}`), a docelowy jej wymaga,
 * nazwą zostaje numer kolejny wystąpienia — reguła jest jedna, jawna
 * i wypowiedziana w sprawozdaniu, bo nazwa zmiennej nie jest czymś, co wolno
 * zgadnąć.
 */

/** Styl symbolu zastępczego rozpoznawany w tekście. */
export interface StylWzorca {
  /** Nazwa stylu widoczna dla Operatora. */
  readonly nazwa: string;
  /** Przykład zapisu w tym stylu. */
  readonly przyklad: string;
  /** Wzorzec rozpoznający wystąpienie; grupa pierwsza niesie nazwę, gdy styl ją ma. */
  readonly wzorzec: RegExp;
  /** Czy styl niesie nazwę zmiennej. */
  readonly zNazwa: boolean;
}

export const STYLE_WZORCOW: readonly StylWzorca[] = [
  {
    nazwa: 'podwójne nawiasy klamrowe',
    przyklad: '{{nazwa}}',
    wzorzec: /\{\{\s*([A-Za-z_][\w.]*)\s*\}\}/g,
    zNazwa: true,
  },
  {
    nazwa: 'dolar z nawiasem klamrowym',
    przyklad: '${nazwa}',
    wzorzec: /\$\{\s*([A-Za-z_][\w.]*)\s*\}/g,
    zNazwa: true,
  },
  {
    nazwa: 'indeks w nawiasach klamrowych',
    przyklad: '{0}',
    wzorzec: /\{(\d+)\}/g,
    zNazwa: false,
  },
  {
    nazwa: 'nazwa w nawiasach klamrowych',
    przyklad: '{nazwa}',
    wzorzec: /\{([A-Za-z_][\w.]*)\}/g,
    zNazwa: true,
  },
  {
    nazwa: 'wzorzec printf',
    przyklad: '%s',
    wzorzec: /%(?:\d+\$)?[sdifv]/g,
    zNazwa: false,
  },
  {
    nazwa: 'znacznik formatu',
    przyklad: '<b>',
    wzorzec: /<\/?[A-Za-z][\w-]*(?:\s[^<>]*)?>/g,
    zNazwa: false,
  },
];

/** Jedno wystąpienie symbolu zastępczego w tekście. */
export interface Wystapienie {
  /** Dosłowny zapis wystąpienia. */
  readonly zapis: string;
  /** Nazwa stylu, który je rozpoznał. */
  readonly styl: string;
  /** Nazwa zmiennej, gdy styl ją niesie; pusta przy stylach bez nazwy. */
  readonly nazwa: string;
  /** Pozycja początku wystąpienia w tekście. */
  readonly poczatek: number;
  /** Pozycja za końcem wystąpienia. */
  readonly koniec: number;
}

/**
 * Wystąpienia symboli zastępczych w tekście, w kolejności występowania.
 *
 * Style sprawdzane są w kolejności wykazu, a odcinek już zajęty nie jest
 * rozpatrywany drugi raz: `{{nazwa}}` pasuje także do wzorca nazwy w pojedynczych
 * nawiasach, więc bez tego jedno wystąpienie policzyłoby się dwa razy pod dwiema
 * nazwami stylu.
 */
export function znajdzWzorce(tekst: string): readonly Wystapienie[] {
  const zajete: { poczatek: number; koniec: number }[] = [];
  const znalezione: Wystapienie[] = [];

  for (const styl of STYLE_WZORCOW) {
    // Kopia wyrażenia na każde przejście: `lastIndex` wzorca globalnego jest
    // stanem, a wzorce stoją w stałej modułowej wspólnej wszystkim wywołaniom.
    const wzorzec = new RegExp(styl.wzorzec.source, styl.wzorzec.flags);
    let dopasowanie = wzorzec.exec(tekst);
    while (dopasowanie !== null) {
      const poczatek = dopasowanie.index;
      const koniec = poczatek + dopasowanie[0].length;
      if (!nachodzi(zajete, poczatek, koniec)) {
        zajete.push({ poczatek, koniec });
        znalezione.push({
          zapis: dopasowanie[0],
          styl: styl.nazwa,
          nazwa: styl.zNazwa ? (dopasowanie[1] ?? '') : '',
          poczatek,
          koniec,
        });
      }
      dopasowanie = wzorzec.exec(tekst);
    }
  }

  return znalezione.sort((pierwsze, drugie) => pierwsze.poczatek - drugie.poczatek);
}

function nachodzi(
  zajete: readonly { poczatek: number; koniec: number }[],
  poczatek: number,
  koniec: number,
): boolean {
  return zajete.some((odcinek) => poczatek < odcinek.koniec && koniec > odcinek.poczatek);
}

/** Styl docelowy przekładu — trzy zapisy wymienione w opracowaniu modułu. */
export const STYLE_DOCELOWE = {
  printf: 'wzorzec printf',
  indeks: 'indeks w nawiasach klamrowych',
  nazwa: 'podwójne nawiasy klamrowe',
} as const;

export type StylDocelowy = (typeof STYLE_DOCELOWE)[keyof typeof STYLE_DOCELOWE];

/** Wynik przekładu stylu: tekst po zamianie i zapis każdej zamiany. */
export interface PrzekladStylu {
  /** Tekst po zamianie symboli zastępczych. */
  readonly tekst: string;
  /** Pary „przed → po" w kolejności wystąpień. */
  readonly zamiany: readonly string[];
  /** Czy którakolwiek nazwa powstała z numeru kolejnego, a nie z tekstu źródłowego. */
  readonly nazwyZNumeru: boolean;
}

/**
 * Przekłada symbole zastępcze na wskazany styl.
 *
 * Znaczniki formatu (`<b>`) nie są przekładane: to nie zmienne, tylko struktura
 * treści, a zamiana ich na zmienną zniszczyłaby dokument. Zostają w tekście
 * takie, jakie są, i nie liczą się do zamian.
 */
export function przelozStyl(tekst: string, docelowy: StylDocelowy): PrzekladStylu {
  const wystapienia = znajdzWzorce(tekst).filter((wpis) => wpis.styl !== 'znacznik formatu');
  const zamiany: string[] = [];
  let nazwyZNumeru = false;
  let wynik = '';
  let ogon = 0;

  wystapienia.forEach((wpis, kolejny) => {
    const nazwa = wpis.nazwa === '' ? String(kolejny) : wpis.nazwa;
    if (wpis.nazwa === '' && docelowy === STYLE_DOCELOWE.nazwa) nazwyZNumeru = true;
    const zapis = zapisWStylu(docelowy, nazwa, kolejny);
    zamiany.push(`${wpis.zapis} → ${zapis}`);
    wynik += tekst.slice(ogon, wpis.poczatek) + zapis;
    ogon = wpis.koniec;
  });

  return { tekst: wynik + tekst.slice(ogon), zamiany, nazwyZNumeru };
}

function zapisWStylu(docelowy: StylDocelowy, nazwa: string, kolejny: number): string {
  if (docelowy === STYLE_DOCELOWE.printf) return '%s';
  if (docelowy === STYLE_DOCELOWE.indeks) return `{${String(kolejny)}}`;
  return `{{${nazwa}}}`;
}
