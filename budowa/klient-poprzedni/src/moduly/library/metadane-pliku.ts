import type { LibraryFile, LibraryMetadata } from '../../../../shared/contract';
import { opisOdmowy } from '../../komponenty/odmowa';
import { poleTekstowe, przycisk, utworzWierszOdpowiedzi } from '../../modele/kontrolki-formularza';
import type { StanBiblioteki } from './stan-biblioteki';

/**
 * Zakładka Metadane panelu metadanych pokazuje blok techniczny z pól pliku i, na
 * żądanie, metadane osadzone w bajtach, oraz formularz Dublin Core zapisywany
 * komendą scalającą wartości bez nadpisywania pustką.
 */
export interface MetadanePliku {
  element: HTMLElement;
  /** Przerysowuje bloki z pliku czynnego. */
  odswiez(): void;
}

/** Pola schematu Dublin Core wymienione w formularzu opisu zasobu, każde z własną etykietą i kodem zapisu. */
const POLA_DUBLIN_CORE: ReadonlyArray<{ kod: string; etykieta: string }> = [
  { kod: 'title', etykieta: 'Tytuł (dc:title)' },
  { kod: 'creator', etykieta: 'Twórca (dc:creator)' },
  { kod: 'subject', etykieta: 'Temat (dc:subject)' },
  { kod: 'date', etykieta: 'Data (dc:date)' },
  { kod: 'rights', etykieta: 'Prawa (dc:rights)' },
];

export function utworzMetadanePliku(stan: StanBiblioteki): MetadanePliku {
  const odpowiedz = utworzWierszOdpowiedzi();

  const techniczne = document.createElement('dl');
  techniczne.className = 'ml-metadane__techniczne';

  const naglowekTechnicznych = document.createElement('h4');
  naglowekTechnicznych.className = 'ml-metadane__naglowek';
  naglowekTechnicznych.textContent = 'Techniczne — z odpowiedzi rdzenia';

  const pola = POLA_DUBLIN_CORE.map((opis) => ({
    kod: opis.kod,
    pole: poleTekstowe({ etykieta: opis.etykieta }),
  }));

  const formularz = document.createElement('div');
  formularz.className = 'ml-metadane__formularz';
  const naglowekOpisu = document.createElement('h4');
  naglowekOpisu.className = 'ml-metadane__naglowek';
  naglowekOpisu.textContent = 'Dublin Core — opis zasobu';
  formularz.append(naglowekOpisu, ...pola.map((wpis) => wpis.pole.element));

  const zapisz = przycisk('Zapisz opis', 'dn-btn dn-btn--sm dn-btn--duch');
  zapisz.dataset['czynnosc'] = 'dublin-core-zapis';
  zapisz.addEventListener('click', () => void zapiszOpis());

  const wczytajOsadzone = przycisk(
    'Odczytaj metadane osadzone w pliku',
    'dn-btn dn-btn--sm dn-btn--zarys',
  );
  wczytajOsadzone.dataset['czynnosc'] = 'metadane-osadzone';
  wczytajOsadzone.addEventListener('click', () => void wczytajOpis(true));

  /** Ostatnio odczytane metadane osadzone — pokazywane w bloku technicznym. */
  let osadzone: string[] = [];

  /** Formularz wypełnia się treścią z rdzenia, żeby zapis wcześniejszy nie zniknął pod pustym polem. */
  async function wczytajOpis(zOsadzonymi: boolean): Promise<void> {
    const plik = stan.czynny();
    if (plik === null) {
      odpowiedz.pokaz('Nie wskazano pliku — opis dotyczy zasobu wybranego w Explorerze.', false);
      return;
    }
    if (zOsadzonymi) {
      odpowiedz.pokaz('Rdzeń czyta metadane osadzone w bajtach pliku…', true);
    }
    const wynik = await stan.zrodlo.opis(plik.id, zOsadzonymi);
    if (!wynik.udany || wynik.wynik === undefined) {
      odpowiedz.pokaz(opisOdmowy('Opis zasobu', wynik.blad?.code, wynik.blad?.message), false);
      return;
    }
    const opis = wynik.wynik.metadata as unknown as Record<string, string | undefined>;
    for (const wpis of pola) {
      wpis.pole.kontrolka.value = opis[wpis.kod] ?? '';
    }
    osadzone = wierszeOsadzone(wynik.wynik.metadata);
    if (zOsadzonymi) {
      odpowiedz.pokaz(
        osadzone.length === 0
          ? 'Rdzeń otworzył bajty pliku i nie znalazł w nich metadanych osadzonych.'
          : `Odczytano metadane osadzone: ${osadzone.length} pozycji.`,
        true,
      );
    }
  }

  /** Zapisuje opis Dublin Core wprowadzony w formularzu. */
  async function zapiszOpis(): Promise<void> {
    const plik = stan.czynny();
    if (plik === null) {
      odpowiedz.pokaz('Nie wskazano pliku — opis nie ma do czego przylgnąć.', false);
      return;
    }
    const opis: Record<string, string> = { fileId: plik.id };
    for (const wpis of pola) {
      opis[wpis.kod] = wpis.pole.kontrolka.value.trim();
    }
    odpowiedz.pokaz('Zapisywanie opisu zasobu…', true);
    const wynik = await stan.zrodlo.zapiszOpis(plik.id, opis as unknown as LibraryMetadata, false);
    if (!wynik.udany || wynik.wynik === undefined) {
      odpowiedz.pokaz(opisOdmowy('Zapis opisu', wynik.blad?.code, wynik.blad?.message), false);
      return;
    }
    stan.wchlon(wynik.wynik.file);
    odpowiedz.pokaz(
      `Opis zasobu ${plik.name} zapisany. Pole zostawione puste zostaje bez zmiany, ` +
        'pole wyczyszczone kasuje wartość — tak działa scalanie opisu.',
      true,
    );
  }

  const pasek = document.createElement('div');
  pasek.className = 'ml-metadane__pasek';
  pasek.append(zapisz, wczytajOsadzone);

  const element = document.createElement('div');
  element.className = 'ml-metadane';
  element.append(naglowekTechnicznych, techniczne, formularz, pasek, odpowiedz.element);

  return {
    element,

    odswiez() {
      const plik = stan.czynny();
      if (plik === null) {
        techniczne.replaceChildren();
        osadzone = [];
        return;
      }
      techniczne.replaceChildren(...wierszeTechniczne(plik, stan, osadzone));
      void wczytajOpis(false);
    },
  };
}

/** Wiersze bloku technicznego — jeden wiersz na pole, które rdzeń faktycznie oddał w odpowiedzi o zasobie. */
function wierszeTechniczne(
  plik: LibraryFile,
  stan: StanBiblioteki,
  osadzone: readonly string[],
): HTMLElement[] {
  const tresc = stan.tresc(plik.id);
  const wpisy: Array<[string, string]> = [
    ['Identyfikator', plik.id],
    ['Nazwa', plik.name],
    ['Ścieżka w repozytorium', plik.path ?? 'rdzeń nie podał'],
    ['Rodzaj treści', plik.mimeType ?? 'rdzeń nie podał'],
    ['Rozmiar', plik.sizeBytes === undefined ? 'rdzeń nie podał' : `${plik.sizeBytes} B`],
    ['Suma kontrolna', plik.checksum ?? 'rdzeń nie podał'],
    ['Wersja bieżąca', plik.versionId ?? 'rdzeń nie podał'],
    ['Moduł wytwórcy', plik.sourceModuleId ?? 'rdzeń nie podał'],
    ['Projekt', plik.projectId ?? 'poza projektem'],
    ['Etykiety', (plik.tags ?? []).join(' · ') || 'brak'],
    ['Kolekcje', (plik.collectionIds ?? []).join(' · ') || 'brak'],
    ['Utworzony', new Date(plik.createdAt).toLocaleString('pl')],
    ['Zmieniony', new Date(plik.updatedAt).toLocaleString('pl')],
    ['Treść w repozytorium', zdanieOTresci(tresc.werdykt, tresc.powod)],
    [
      'Metadane osadzone w pliku',
      osadzone.length === 0
        ? 'nieodczytane — sięgają po bajty pliku, więc wchodzą na wyraźne żądanie'
        : osadzone.join(' · '),
    ],
  ];
  return wpisy.flatMap(([nazwa, wartosc]) => {
    const podpis = document.createElement('dt');
    podpis.className = 'ml-metadane__podpis';
    podpis.textContent = nazwa;

    const opis = document.createElement('dd');
    opis.className = 'ml-metadane__wartosc';
    opis.textContent = wartosc;
    return [podpis, opis];
  });
}

/** Zdanie o treści zasobu — jedno zdanie na każdy możliwy werdykt odpowiedzi rdzenia o jej dostępności. */
function zdanieOTresci(werdykt: string, powod: string): string {
  if (werdykt === 'osiagalna') return 'rdzeń oddaje treść pliku';
  if (werdykt === 'brak') return `rdzeń orzekł brak treści — ${powod}`;
  if (werdykt === 'odmowa') return `odczyt treści się nie udał — ${powod}`;
  if (werdykt === 'odwolanie') return `rdzeń wskazał miejsce treści zamiast bajtów — ${powod}`;
  return 'nikt jeszcze nie pytał rdzenia o treść tego pliku';
}

/**
 * Składa czytelny wykaz metadanych osadzonych z odpowiedzi rdzenia.
 *
 * Wymienia wyłącznie pola, które przyszły. Pole puste wypisane jako „brak"
 * mówiłoby o pliku coś, czego nikt nie zmierzył.
 */
function wierszeOsadzone(opis: LibraryMetadata): string[] {
  const techniczne = opis.technical;
  if (techniczne === undefined) return [];
  const wpisy: string[] = [];
  if (techniczne.mimeType !== undefined) wpisy.push(`rodzaj treści: ${techniczne.mimeType}`);
  if (techniczne.sizeBytes !== undefined) wpisy.push(`rozmiar: ${techniczne.sizeBytes} B`);
  if (techniczne.checksum !== undefined) wpisy.push(`suma: ${techniczne.checksum.slice(0, 12)}…`);
  if (techniczne.pageCount !== undefined) wpisy.push(`stron: ${techniczne.pageCount}`);
  if (techniczne.width !== undefined && techniczne.height !== undefined) {
    wpisy.push(`wymiary: ${techniczne.width}×${techniczne.height}`);
  }
  if (techniczne.durationMs !== undefined) {
    wpisy.push(`czas trwania: ${Math.round(techniczne.durationMs / 1000)} s`);
  }
  if (techniczne.gpsLatitude !== undefined && techniczne.gpsLongitude !== undefined) {
    wpisy.push(`położenie: ${techniczne.gpsLatitude.toFixed(5)}, ${techniczne.gpsLongitude.toFixed(5)}`);
  }
  if (techniczne.exif !== undefined) {
    const exif = techniczne.exif as unknown as Record<string, string>;
    for (const [nazwa, wartosc] of Object.entries(exif)) {
      wpisy.push(`${nazwa}: ${wartosc}`);
    }
  }
  return wpisy;
}
