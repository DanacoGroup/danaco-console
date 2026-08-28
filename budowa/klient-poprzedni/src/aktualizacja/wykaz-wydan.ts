/** Adres kanału wydań, czytanego jako plik `wydania.json` — jedno źródło prawdy dla witryny i aplikacji. */
export const ADRES_WYDAN = 'https://danaco-console.pl/wydania.json';

/** Jedno wydanie w chronologii wykazu, w kształcie zgodnym z zapisem pliku `wydania.json` po stronie kanału. */
export interface Wydanie {
  wersja: string;
  data: string;
  system?: string;
  plik?: string;
  nazwaPliku?: string;
  rozmiar?: string;
  suma?: string;
  zmiany?: string;
}

/** Czy odczytana wartość ma kształt wydania. Wykaz z sieci jest treścią obcą
 *  i sprawdza się go w całości, zanim cokolwiek z niego pokażemy. */
function czyWydanie(wartosc: unknown): wartosc is Wydanie {
  if (typeof wartosc !== 'object' || wartosc === null) return false;
  const zapis = wartosc as Record<string, unknown>;
  return typeof zapis['wersja'] === 'string' && typeof zapis['data'] === 'string';
}

/** Porównuje wersje w postaci `1.2.3` człon po członie liczbowo, z pominięciem przedrostka `v` i członów nieliczbowych. */
export function nowszaWersja(kandydat: string, biezaca: string): boolean {
  const czlony = (w: string): number[] =>
    w
      .trim()
      .replace(/^[vV]/, '')
      .split('.')
      .map((c) => {
        const przyciety = c.trim();
        if (!/^\d+$/.test(przyciety)) return 0;
        const liczba = Number.parseInt(przyciety, 10);
        return Number.isFinite(liczba) ? liczba : 0;
      });
  const a = czlony(kandydat);
  const b = czlony(biezaca);
  const dlugosc = Math.max(a.length, b.length);
  for (let i = 0; i < dlugosc; i += 1) {
    const x = a[i] ?? 0;
    const y = b[i] ?? 0;
    if (x !== y) return x > y;
  }
  return false;
}

/** Nazwa systemu, na którym biegnie ta kopia, w słownictwie pliku `wydania.json` — „windows”, „macos” albo „linux”. */
function systemBiezacy(): string {
  const nawigator = (globalThis as { navigator?: { userAgent?: string; platform?: string } })
    .navigator;
  const opis = `${nawigator?.userAgent ?? ''} ${nawigator?.platform ?? ''}`.toLowerCase();
  if (opis.includes('windows') || opis.includes('win32') || opis.includes('win64')) return 'windows';
  if (opis.includes('mac')) return 'macos';
  if (opis.includes('linux') || opis.includes('x11')) return 'linux';
  return '';
}

/** Rozstrzyga, czy wydanie jest przeznaczone na system bieżącej kopii, na podstawie pola `system` z pliku wydań. */
function wydanieNaTenSystem(wydanie: Wydanie): boolean {
  const zadeklarowany = (wydanie.system ?? '').trim().toLowerCase();
  if (zadeklarowany === '') return true;
  const tutejszy = systemBiezacy();
  if (tutejszy === '') return true;
  return zadeklarowany === tutejszy;
}

/** Wynik odczytu kanału wydań: sześć rozróżnionych stanów świata zamiast jednej zbiorczej wartości null. */
export type OdczytWykazu =
  /** Kanał odpowiedział i niesie co najmniej jedno czytelne wydanie. */
  | { stan: 'wykaz'; wydania: Wydanie[] }
  /** Kanał odpowiedział i mówi wprost: wydania jeszcze nie ma. To PRAWDA, nie brak. */
  | { stan: 'pusty' }
  /** Pod adresem kanału nic nie stoi (404). Sieć działa, wykazu nie ma. */
  | { stan: 'brak-wydania-pod-adresem'; status: number }
  /** Serwer odpowiedział kodem błędu — kanał stoi, ale ma kłopot po swojej stronie. */
  | { stan: 'odpowiedz-serwera'; status: number }
  /** Żądanie nie doszło do nikogo: brak sieci, DNS bez odpowiedzi, TLS odrzucony. */
  | { stan: 'brak-lacznosci' }
  /** Coś przyszło, ale nie jest wykazem wydań (nie-JSON, obcy kształt, same śmieci). */
  | { stan: 'odpowiedz-nieczytelna' };

/**
 * Czyta kanał wydań i NAZYWA to, co zastał.
 *
 * Odczyt nie rzuca nigdy. Każda droga wyjścia jest nazwanym stanem — kanał
 * wydań nie może popsuć uruchomienia produktu ani przerwać obiegu pytania.
 */
export async function pobierzWykazWydan(adres: string = ADRES_WYDAN): Promise<OdczytWykazu> {
  let odpowiedz: Response;
  try {
    odpowiedz = await fetch(adres, { cache: 'no-store' });
  } catch {
    // Rzut z `fetch` znaczy „nie doszło do nikogo" — odpowiedź serwera tędy nie przechodzi.
    return { stan: 'brak-lacznosci' };
  }

  if (!odpowiedz.ok) {
    return odpowiedz.status === 404
      ? { stan: 'brak-wydania-pod-adresem', status: odpowiedz.status }
      : { stan: 'odpowiedz-serwera', status: odpowiedz.status };
  }

  let tresc: unknown;
  try {
    tresc = await odpowiedz.json();
  } catch {
    // Odpowiedź przyszła, ale nie jest JSON-em — to nie jest brak łączności.
    return { stan: 'odpowiedz-nieczytelna' };
  }

  const wykaz = (tresc as { wydania?: unknown } | null)?.wydania;
  if (!Array.isArray(wykaz)) return { stan: 'odpowiedz-nieczytelna' };
  // Tablica pusta to deklaracja kanału, „nie ma jeszcze żadnego wydania", a nie brak odpowiedzi.
  if (wykaz.length === 0) return { stan: 'pusty' };

  // Wykaz z sieci jest treścią obcą — pozycje bez `wersja`/`data` odpadają.
  const wydania = wykaz.filter(czyWydanie);
  // Coś w wykazie było, ale nic z tego nie jest wydaniem — bliżej „nieczytelnej" niż „pustej".
  if (wydania.length === 0) return { stan: 'odpowiedz-nieczytelna' };

  return { stan: 'wykaz', wydania };
}

/** Pobiera wykaz i oddaje wydanie nowsze niż zainstalowane, porównane po wszystkich pozycjach, albo `null`. */
export async function nowszeWydanie(
  wersjaBiezaca: string,
  adres: string = ADRES_WYDAN,
): Promise<Wydanie | null> {
  const odczyt = await pobierzWykazWydan(adres);
  if (odczyt.stan !== 'wykaz') return null;

  let najnowsze: Wydanie | null = null;
  for (const pozycja of odczyt.wydania) {
    if (!pozycja.plik) continue;
    // Wydanie na inny system nie jest kandydatem — pominięcie idzie przed porównaniem wersji.
    if (!wydanieNaTenSystem(pozycja)) continue;
    if (!nowszaWersja(pozycja.wersja, wersjaBiezaca)) continue;
    if (najnowsze === null || nowszaWersja(pozycja.wersja, najnowsze.wersja)) {
      najnowsze = pozycja;
    }
  }
  return najnowsze;
}
