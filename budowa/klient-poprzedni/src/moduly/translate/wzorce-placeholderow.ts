/**
 * Symbole zastępcze i znaczniki formatu — rozpoznanie w tekście źródłowym oraz przekład ich stylu
 * między platformami — liczy okno, nie rdzeń, bo kontrakt takiej komendy nie ma.
 */

/** Styl symbolu zastępczego rozpoznawany w tekście niesie nazwę, przykład zapisu, wzorzec dopasowania oraz informację o nazwie zmiennej. */
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

/** Jedno wystąpienie symbolu zastępczego w tekście niesie dosłowny zapis, nazwę rozpoznającego stylu oraz pozycję początku i końca w tekście. */
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
 * Wystąpienia symboli zastępczych w tekście, w kolejności występowania, sprawdzają style po
 * kolei, a odcinek już zajęty nie jest rozpatrywany drugi raz.
 */
export function znajdzWzorce(tekst: string): readonly Wystapienie[] {
  const zajete: { poczatek: number; koniec: number }[] = [];
  const znalezione: Wystapienie[] = [];

  for (const styl of STYLE_WZORCOW) {
    // Kopia wyrażenia na każde przejście: stan dopasowania wzorca globalnego jest wspólny wywołaniom.
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

/** Styl docelowy przekładu ma trzy zapisy wymienione w opracowaniu modułu: wzorzec printf, indeks i podwójne nawiasy. */
export const STYLE_DOCELOWE = {
  printf: 'wzorzec printf',
  indeks: 'indeks w nawiasach klamrowych',
  nazwa: 'podwójne nawiasy klamrowe',
} as const;

export type StylDocelowy = (typeof STYLE_DOCELOWE)[keyof typeof STYLE_DOCELOWE];

/** Wynik przekładu stylu niesie tekst po zamianie symboli zastępczych oraz zapis każdej wykonanej zamiany z osobna. */
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
