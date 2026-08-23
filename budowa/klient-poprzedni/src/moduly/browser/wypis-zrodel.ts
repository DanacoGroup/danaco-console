import type { BrowserSource } from '../../../../shared/contract';

/**
 * Wypis wykazu źródeł w czterech notacjach bibliograficznych, których żąda
 * opracowanie modułu: `BibTeX`, `RIS`, `CSV` i `Markdown`.
 *
 * Wypis powstaje w kliencie, bo cała jego treść jest już w oknie: wykaz przyszedł
 * komendą `browser.source.list`, a notacja jest wyłącznie sposobem jego zapisania.
 * Komendy eksportu bibliografii kontrakt nie niesie i nie musi.
 *
 * Rodzajem wpisu jest zasób elektroniczny — moduł zbiera strony internetowe
 * i tylko o nich może tak zaświadczyć. Datą jest czas dodania źródła do wykazu
 * okna: daty publikacji strony rdzeń nie oddaje, więc wypis jej nie zmyśla.
 */

/** Jedna notacja wypisu: wartość nastawy, nazwa, rozszerzenie i rodzaj treści. */
export interface NotacjaWypisu {
  wartosc: string;
  etykieta: string;
  opis: string;
  rozszerzenie: string;
  rodzaj: string;
}

export const NOTACJE_WYPISU: readonly NotacjaWypisu[] = [
  {
    wartosc: 'markdown',
    etykieta: 'Markdown',
    opis: 'Wykaz odnośników do wklejenia w dokument.',
    rozszerzenie: 'md',
    rodzaj: 'text/markdown',
  },
  {
    wartosc: 'csv',
    etykieta: 'CSV',
    opis: 'Wiersze rozdzielone przecinkiem, z nagłówkiem kolumn.',
    rozszerzenie: 'csv',
    rodzaj: 'text/csv',
  },
  {
    wartosc: 'bibtex',
    etykieta: 'BibTeX',
    opis: 'Wpisy @online do menedżera bibliografii.',
    rozszerzenie: 'bib',
    rodzaj: 'application/x-bibtex',
  },
  {
    wartosc: 'ris',
    etykieta: 'RIS',
    opis: 'Wpisy TY − ELEC przyjmowane przez menedżery bibliografii.',
    rozszerzenie: 'ris',
    rodzaj: 'application/x-research-info-systems',
  },
];

/** Notacja po wartości nastawy; nieznana wartość wraca pierwszą z wykazu. */
export function notacja(wartosc: string): NotacjaWypisu {
  return NOTACJE_WYPISU.find((pozycja) => pozycja.wartosc === wartosc) ?? NOTACJE_WYPISU[0];
}

/** Wykaz źródeł zapisany wskazaną notacją. */
export function wypisz(zrodla: readonly BrowserSource[], wartosc: string): string {
  const wybrana = notacja(wartosc).wartosc;
  if (wybrana === 'csv') return wypiszCsv(zrodla);
  if (wybrana === 'bibtex') return zrodla.map(wpisBibTeX).join('\n\n');
  if (wybrana === 'ris') return zrodla.map(wpisRis).join('\n');
  return zrodla.map(wpisMarkdown).join('\n');
}

/** Tytuł źródła albo jego adres, gdy rdzeń tytułu nie oddał. */
function tytul(zrodlo: BrowserSource): string {
  const nazwa = (zrodlo.title ?? '').trim();
  return nazwa === '' ? zrodlo.url : nazwa;
}

/** Data dodania źródła do wykazu okna w zapisie ISO 8601 (rok-miesiąc-dzień). */
function data(zrodlo: BrowserSource): string {
  return new Date(zrodlo.createdAt).toISOString().slice(0, 10);
}

function wpisMarkdown(zrodlo: BrowserSource): string {
  const znacznik = zrodlo.key === true ? ' — źródło kluczowe' : '';
  return `- [${tytul(zrodlo)}](${zrodlo.url}) · ${data(zrodlo)}${znacznik}`;
}

function wypiszCsv(zrodla: readonly BrowserSource[]): string {
  const naglowek = 'identyfikator,tytul,adres,data_dodania,kluczowe';
  const wiersze = zrodla.map((zrodlo) =>
    [
      polePola(zrodlo.id),
      polePola(tytul(zrodlo)),
      polePola(zrodlo.url),
      polePola(data(zrodlo)),
      zrodlo.key === true ? 'tak' : 'nie',
    ].join(','),
  );
  return [naglowek, ...wiersze].join('\n');
}

/** Pole CSV w cudzysłowie, z cudzysłowem podwojonym wewnątrz (RFC 4180). */
function polePola(wartosc: string): string {
  return `"${wartosc.replaceAll('"', '""')}"`;
}

function wpisBibTeX(zrodlo: BrowserSource): string {
  return [
    `@online{${zrodlo.id},`,
    `  title = {${tytul(zrodlo)}},`,
    `  url = {${zrodlo.url}},`,
    `  urldate = {${data(zrodlo)}}`,
    '}',
  ].join('\n');
}

function wpisRis(zrodlo: BrowserSource): string {
  return [
    'TY  - ELEC',
    `ID  - ${zrodlo.id}`,
    `TI  - ${tytul(zrodlo)}`,
    `UR  - ${zrodlo.url}`,
    `DA  - ${data(zrodlo).replaceAll('-', '/')}`,
    'ER  - ',
    '',
  ].join('\n');
}
