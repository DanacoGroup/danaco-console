import type { LibraryFile, LibraryVersion } from '../../../../shared/contract';

// Wywóz treści, którą okno już ma, do pliku na urządzenie: opis zasobów wychodzi, bajty plików nie.

/** Nagłówek wspólny wytworów tekstowych — mówi, czym plik jest i czym nie jest, żeby manifest nie uchodził za archiwum z treścią. */
const ZASTRZEZENIE_OPISU =
  'Wytwór zawiera wyłącznie opis zasobów (metryka, etykiety, kolekcje, sumy kontrolne). ' +
  'Treści plików w nim nie ma: kontrakt nie niesie komendy pobierającej bajty zasobu ' +
  'biblioteki na urządzenie.';

/**
 * Raport zmian dokumentu zestawia pełną historię wersji w postaci czytelnej;
 * wariant z treścią każdej wersji nie powstaje, bo kontrakt nie oddaje treści
 * wersji niebieżącej.
 */
export function raportZmian(plik: LibraryFile, wersje: readonly LibraryVersion[]): string {
  const wiersze: string[] = [
    `# Historia wersji — ${plik.name}`,
    '',
    `Zasób: ${plik.id}`,
    `Wersja bieżąca: ${plik.versionId ?? 'nieokreślona'}`,
    `Wersji w raporcie: ${wersje.length}`,
    '',
    ZASTRZEZENIE_OPISU,
    '',
    '| Wersja | Etykieta | Sprawca | Rozmiar | Suma kontrolna | Powstała |',
    '|---|---|---|---|---|---|',
  ];
  for (const wersja of wersje) {
    wiersze.push(
      [
        '',
        wersja.id,
        wersja.label ?? '—',
        wersja.author ?? '—',
        wersja.sizeBytes === undefined ? '—' : `${wersja.sizeBytes} B`,
        wersja.checksum ?? '—',
        new Date(wersja.createdAt).toLocaleString('pl'),
        '',
      ].join(' | '),
    );
  }
  return `${wiersze.join('\n')}\n`;
}

/**
 * Mapa struktury kolekcji dokumentuje porządek repozytorium: kolekcja jest
 * w kontrakcie samym identyfikatorem, więc mapa wypisuje kody i przypisane
 * pliki, bez nazw własnych, których okno nie zna.
 */
export function mapaKolekcji(pliki: readonly LibraryFile[], kolekcje: readonly string[]): string {
  const wiersze: string[] = [
    '# Mapa kolekcji repozytorium',
    '',
    `Kolekcji: ${kolekcje.length} · plików w wykazie: ${pliki.length}`,
    '',
    ZASTRZEZENIE_OPISU,
    '',
    'Kolekcje wypisane są kodami. Kontrakt nie ma komendy katalogu kolekcji, ' +
      'więc nazwa własna kolekcji nie jest oknu znana poza chwilą jej założenia.',
    '',
  ];
  for (const kod of kolekcje) {
    const nalezace = pliki.filter((plik) => (plik.collectionIds ?? []).includes(kod));
    wiersze.push(`## ${kod} — plików ${nalezace.length}`, '');
    for (const plik of nalezace) {
      wiersze.push(`- ${plik.name} (${plik.id})`);
    }
    wiersze.push('');
  }
  const pozaKolekcjami = pliki.filter((plik) => (plik.collectionIds ?? []).length === 0);
  wiersze.push(`## Poza kolekcjami — plików ${pozaKolekcjami.length}`, '');
  for (const plik of pozaKolekcjami) {
    wiersze.push(`- ${plik.name} (${plik.id})`);
  }
  return `${wiersze.join('\n')}\n`;
}

/**
 * Tezaurus etykiet w zapisie SKOS, serializacja Turtle, jest płaski: kontrakt
 * niesie etykietę jako sam napis, więc nie ma z czego zbudować relacji
 * nadrzędności ani pokrewieństwa między pojęciami.
 */
export function tezaurusSkos(etykiety: readonly string[], licznik: (kod: string) => number): string {
  const wiersze: string[] = [
    '@prefix skos: <http://www.w3.org/2004/02/skos/core#> .',
    '@prefix dct: <http://purl.org/dc/terms/> .',
    '@prefix rdf: <http://www.w3.org/1999/02/22-rdf-syntax-ns#> .',
    '',
    '# Tezaurus etykiet repozytorium Danaco Console, moduł Library.',
    '# Wykaz jest płaski: etykieta jest w kontrakcie samym napisem, więc relacje',
    '# skos:broader, skos:narrower i skos:related nie mają skąd powstać.',
    '',
    '<urn:danaco:library:tezaurus> a skos:ConceptScheme ;',
    '    dct:title "Etykiety repozytorium Library"@pl .',
    '',
  ];
  for (const kod of etykiety) {
    wiersze.push(
      `<${adresPojecia(kod)}> a skos:Concept ;`,
      '    skos:inScheme <urn:danaco:library:tezaurus> ;',
      `    skos:prefLabel ${literal(kod)}@pl ;`,
      `    dct:extent "${licznik(kod)}" .`,
      '',
    );
  }
  return wiersze.join('\n');
}

/**
 * Manifest repozytorium opisuje stan zbioru w punkcie czasu bez treści
 * zasobów: pole zawartości mówi to wprost, bo archiwum z bajtami wymagałoby
 * komendy, której kontrakt nie niesie.
 */
export function manifestRepozytorium(pliki: readonly LibraryFile[]): string {
  return `${JSON.stringify(
    {
      wytwor: 'manifest repozytorium Library',
      zawartosc: 'wyłącznie opis zasobów; treści plików manifest nie zawiera',
      zlozony: new Date().toISOString(),
      pozycji: pliki.length,
      zasoby: pliki.map((plik) => ({
        id: plik.id,
        nazwa: plik.name,
        sciezka: plik.path ?? null,
        rodzajTresci: plik.mimeType ?? null,
        rozmiarBajtow: plik.sizeBytes ?? null,
        sumaKontrolna: plik.checksum ?? null,
        wersjaBiezaca: plik.versionId ?? null,
        modulWytworcy: plik.sourceModuleId ?? null,
        projekt: plik.projectId ?? null,
        etykiety: plik.tags ?? [],
        kolekcje: plik.collectionIds ?? [],
        utworzony: new Date(plik.createdAt).toISOString(),
        zmieniony: new Date(plik.updatedAt).toISOString(),
      })),
    },
    null,
    2,
  )}\n`;
}

/** Metadane repozytorium w postaci tabelarycznej, gotowe do wywozu jako arkusz kalkulacyjny z pełnym opisem zasobów. */
export function metadaneCsv(pliki: readonly LibraryFile[]): string {
  const naglowek = [
    'id',
    'nazwa',
    'sciezka',
    'rodzaj_tresci',
    'rozmiar_bajtow',
    'suma_kontrolna',
    'wersja_biezaca',
    'modul_wytworcy',
    'projekt',
    'etykiety',
    'kolekcje',
    'utworzony',
    'zmieniony',
  ];
  const wiersze = pliki.map((plik) =>
    [
      plik.id,
      plik.name,
      plik.path ?? '',
      plik.mimeType ?? '',
      plik.sizeBytes === undefined ? '' : String(plik.sizeBytes),
      plik.checksum ?? '',
      plik.versionId ?? '',
      plik.sourceModuleId ?? '',
      plik.projectId ?? '',
      (plik.tags ?? []).join(' '),
      (plik.collectionIds ?? []).join(' '),
      new Date(plik.createdAt).toISOString(),
      new Date(plik.updatedAt).toISOString(),
    ]
      .map(poleCsv)
      .join(','),
  );
  return `${[naglowek.join(','), ...wiersze].join('\r\n')}\r\n`;
}

/**
 * Pole CSV według RFC 4180: cudzysłów podwaja się, a pole z separatorem,
 * cudzysłowem albo końcem wiersza wchodzi w cudzysłowy.
 */
function poleCsv(wartosc: string): string {
  if (!/[",\r\n]/.test(wartosc)) return wartosc;
  return `"${wartosc.replaceAll('"', '""')}"`;
}

/** Adres pojęcia tezaurusa; etykieta wchodzi w identyfikator zakodowana, tak aby znaki specjalne nie łamały składni. */
function adresPojecia(kod: string): string {
  return `urn:danaco:library:etykieta:${encodeURIComponent(kod)}`;
}

/** Literał Turtle wraz z ucieczkami wymaganymi przez zapis: lewy ukośnik, cudzysłów oraz znaki końca wiersza. */
function literal(wartosc: string): string {
  const uciekniety = wartosc
    .replaceAll('\\', '\\\\')
    .replaceAll('"', '\\"')
    .replaceAll('\n', '\\n')
    .replaceAll('\r', '\\r')
    .replaceAll('\t', '\\t');
  return `"${uciekniety}"`;
}
