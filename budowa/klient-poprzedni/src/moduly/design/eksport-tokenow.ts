import type { GrupaZetonow, Zeton } from './zetony-systemu';

/**
 * Wydanie żetonów do kodu oraz przewodnik stylu.
 *
 * Wszystko tutaj składa przeglądarka z wartości odczytanych z motywu
 * obowiązującego — rdzeń nie bierze w tym udziału i nie musi. To jest ta część
 * grupy „Tokeny projektowe i system projektowy", która komendy nie potrzebuje,
 * więc nazwanie jej brakiem kontraktu byłoby zmyśleniem długu.
 *
 * Czego tu NIE ma i dlaczego. Opracowanie wymienia przy eksporcie także wydanie
 * dla systemów mobilnych oraz wydanie przewodnika do modułów Library i Studio.
 * Pierwszego nie ma, bo wymagałoby przekładu ról systemu wizualnego na pojęcia
 * dwóch obcych platform — a taki przekład jest rozstrzygnięciem projektowym,
 * nie zapisem pliku. Drugiego nie ma, bo wydanie czegokolwiek do innego modułu
 * wymaga komendy, której kontrakt nie zna; okno nazywa ten brak zamiast wysyłać
 * plik w próżnię.
 *
 * Postacie wydania nie są tu ozdobą: zmienne CSS to postać, w której ten produkt
 * żetony trzyma; pozostałe trzy są postaciami, w których przyjmuje je kod
 * korzystający z systemu.
 */

/** Postać, w której żetony wychodzą z modułu. */
export type PostacWydania = 'zmienne-css' | 'scss' | 'tailwind' | 'modul-js';

/** Nazwy postaci wydania wraz z nazwą pliku, który powstaje. */
export const POSTACIE_WYDANIA: readonly (readonly [PostacWydania, string, string])[] = [
  ['zmienne-css', 'Zmienne CSS', 'zetony.css'],
  ['scss', 'SCSS', 'zetony.scss'],
  ['tailwind', 'Konfiguracja Tailwind', 'zetony.tailwind.js'],
  ['modul-js', 'Moduł JavaScript', 'zetony.js'],
];

/** Nagłówek wydania — mówi, skąd wartości pochodzą i czym NIE są. */
function naglowek(motyw: string): readonly string[] {
  return [
    'Żetony systemu wizualnego Danaco Console.',
    `Wydanie złożone z motywu obowiązującego w chwili wydania: ${motyw}.`,
    'Wartości odczytano z uruchomionego produktu, nie przepisano z dokumentu — drugi motyw',
    'ma własne wartości tych samych ról i wymaga osobnego wydania.',
  ];
}

/** Żetony wszystkich grup w jednym ciągu, bez tych, których motyw nie definiuje. */
function zdefiniowane(grupy: readonly GrupaZetonow[]): readonly Zeton[] {
  return grupy.flatMap((grupa) => grupa.zetony.filter((zeton) => zeton.wartosc !== ''));
}

/**
 * Składa treść wydania.
 *
 * @param motyw nazwa motywu obowiązującego — wchodzi do nagłówka, bo wydanie bez
 *   niej nie mówi, którą połowę systemu niesie.
 */
export function zlozWydanie(
  grupy: readonly GrupaZetonow[],
  postac: PostacWydania,
  motyw: string,
): string {
  switch (postac) {
    case 'zmienne-css':
      return zlozCss(grupy, motyw);
    case 'scss':
      return zlozScss(grupy, motyw);
    case 'tailwind':
      return zlozTailwind(grupy, motyw);
    default:
      return zlozModulJs(grupy, motyw);
  }
}

function zlozCss(grupy: readonly GrupaZetonow[], motyw: string): string {
  const wiersze: string[] = [`/*`, ...naglowek(motyw).map((zdanie) => ` * ${zdanie}`), ` */`, ''];
  wiersze.push(':root {');
  for (const grupa of grupy) {
    const zetony = grupa.zetony.filter((zeton) => zeton.wartosc !== '');
    if (zetony.length === 0) continue;
    wiersze.push(`  /* ${grupa.nazwa} */`);
    for (const zeton of zetony) wiersze.push(`  ${zeton.zmienna}: ${zeton.wartosc};`);
    wiersze.push('');
  }
  wiersze.push('}');
  return `${wiersze.join('\n')}\n`;
}

function zlozScss(grupy: readonly GrupaZetonow[], motyw: string): string {
  const wiersze: string[] = naglowek(motyw).map((zdanie) => `// ${zdanie}`);
  wiersze.push('');
  for (const grupa of grupy) {
    const zetony = grupa.zetony.filter((zeton) => zeton.wartosc !== '');
    if (zetony.length === 0) continue;
    wiersze.push(`// ${grupa.nazwa}`);
    for (const zeton of zetony) wiersze.push(`$${zeton.nazwa}: ${zeton.wartosc};`);
    wiersze.push('');
  }
  return `${wiersze.join('\n')}\n`;
}

/**
 * Konfiguracja Tailwind — wyłącznie żetony barwne w gałęzi barw i odstępy
 * w gałęzi odstępów.
 *
 * Wrzucenie wszystkiego do jednej gałęzi dałoby konfigurację, która się wczyta
 * i nic sensownego nie zrobi; podział wedle roli jest tym, co czyni wydanie
 * użytecznym.
 */
function zlozTailwind(grupy: readonly GrupaZetonow[], motyw: string): string {
  const barwy = zdefiniowane(grupy).filter((zeton) => zeton.rodzaj === 'barwa');
  const odstepy = zdefiniowane(grupy).filter((zeton) => zeton.nazwa.startsWith('od-'));
  const promienie = zdefiniowane(grupy).filter((zeton) => zeton.nazwa.startsWith('r-'));

  const wiersze: string[] = naglowek(motyw).map((zdanie) => `// ${zdanie}`);
  wiersze.push('', 'export default {', '  theme: {', '    extend: {');
  wiersze.push('      colors: {', ...pary(barwy), '      },');
  wiersze.push('      spacing: {', ...pary(odstepy), '      },');
  wiersze.push('      borderRadius: {', ...pary(promienie), '      },');
  wiersze.push('    },', '  },', '};');
  return `${wiersze.join('\n')}\n`;
}

function pary(zetony: readonly Zeton[]): readonly string[] {
  return zetony.map((zeton) => `        '${zeton.nazwa}': '${zeton.wartosc}',`);
}

function zlozModulJs(grupy: readonly GrupaZetonow[], motyw: string): string {
  const wiersze: string[] = naglowek(motyw).map((zdanie) => `// ${zdanie}`);
  wiersze.push('', 'export const zetony = {');
  for (const grupa of grupy) {
    const zetony = grupa.zetony.filter((zeton) => zeton.wartosc !== '');
    if (zetony.length === 0) continue;
    wiersze.push(`  // ${grupa.nazwa}`);
    for (const zeton of zetony) wiersze.push(`  '${zeton.nazwa}': '${zeton.wartosc}',`);
  }
  wiersze.push('};');
  return `${wiersze.join('\n')}\n`;
}

/**
 * Przewodnik stylu — dokumentacja systemu projektowego jako jeden plik.
 *
 * Powstaje z tych samych wartości co wydania kodu, więc nie może się z nimi
 * rozejść. Jest dokumentem do czytania, nie stroną produktu: nie wciąga
 * arkuszy platformy i nie udaje jej interfejsu.
 */
export function zlozPrzewodnikStylu(grupy: readonly GrupaZetonow[], motyw: string): string {
  const sekcje = grupy
    .map((grupa) => {
      const zetony = grupa.zetony.filter((zeton) => zeton.wartosc !== '');
      if (zetony.length === 0) return '';
      const wiersze = zetony
        .map(
          (zeton) =>
            `      <tr><td>${zeton.nazwa}</td><td><code>${zeton.zmienna}</code></td>` +
            `<td>${zeton.wartosc}</td></tr>`,
        )
        .join('\n');
      return (
        `    <h2>${grupa.nazwa}</h2>\n` +
        '    <table>\n' +
        '      <tr><th>Rola</th><th>Żeton</th><th>Wartość</th></tr>\n' +
        `${wiersze}\n` +
        '    </table>'
      );
    })
    .filter((sekcja) => sekcja !== '')
    .join('\n');

  return (
    '<!doctype html>\n<html lang="pl">\n  <head>\n    <meta charset="utf-8">\n' +
    '    <title>Danaco Console — przewodnik systemu projektowego</title>\n' +
    '  </head>\n  <body>\n' +
    '    <h1>Danaco Console — przewodnik systemu projektowego</h1>\n' +
    `    <p>${naglowek(motyw).join(' ')}</p>\n` +
    `${sekcje}\n` +
    '  </body>\n</html>\n'
  );
}
