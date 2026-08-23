/**
 * Wykaz wydań — odczyt kanału wydań `danaco-console.pl`.
 *
 * Czyta stamtąd jeden plik — `wydania.json`, ten sam, z którego składa się
 * strona „Pobierz" (`budowa/witryna/wydania.json`). Jedno źródło prawdy dla
 * witryny i dla aplikacji: gdyby aplikacja pytała o co innego niż to, co widzi
 * człowiek na stronie, dwie prawdy rozjechałyby się pierwszego dnia.
 *
 * Odczyt nie jest warunkiem pracy. Brak sieci, adres nieosiągalny, odpowiedź
 * nieczytelna i wykaz pusty znaczą dla banera to samo: nie wiadomo o żadnym
 * nowszym wydaniu, więc baner się nie pokazuje. Żadna ścieżka nie rzuca
 * wyjątkiem — kanał wydań nie może popsuć uruchomienia produktu.
 *
 * Dla wykazu pokazywanego Operatorowi te stany znaczą co innego i muszą być
 * rozróżnione, inaczej okno mówi „nie ma wydań" wtedy, gdy prawdą jest „nie ma
 * sieci". Dlatego są tu dwie drogi odczytu:
 *   • `pobierzWykazWydan()` — oddaje nazwany stan świata (typ `OdczytWykazu`),
 *   • `nowszeWydanie()`     — nakładka na tamtą, oddająca `Wydanie|null` dla
 *                             banera, któremu wystarczy „jest co zakładać".
 */

/** Adres kanału wydań. */
export const ADRES_WYDAN = 'https://danaco-console.pl/wydania.json';

/** Jedno wydanie w chronologii — kształt pliku `wydania.json`. */
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

/**
 * Porównanie wersji w postaci `1.2.3`.
 *
 * Człony porównuje się liczbowo, nie napisami — inaczej `1.10` byłoby starsze
 * niż `1.9`. Człon nieliczbowy (np. `1.2.0-rc1`) schodzi do zera: wydanie próbne
 * nie ma prawa udawać nowszego niż wydanie właściwe.
 *
 * Dwa miejsca, w których `parseInt` sam z siebie nie wystarcza:
 *  • `Number.parseInt('4-rc1')` oddaje 4, nie NaN — `1.2.4-rc1` wyszłoby równe
 *    `1.2.4`. Dlatego człon musi być liczbą w całości (`/^\d+$/`), inaczej jest
 *    zerem.
 *  • Przedrostek `v` (`v1.9.0`) zbija pierwszy człon do zera. Zdejmujemy go
 *    z napisu, bo `v1.9.0` i `1.9.0` to zapis tej samej wersji.
 */
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

/**
 * Nazwa systemu, na którym biegnie ta kopia — w słownictwie `wydania.json`
 * („Windows albo Linux").
 *
 * Rozpoznajemy po `navigator`, bo interfejs działa w oknie przeglądarkowym
 * także wtedy, gdy siedzi w powłoce natywnej — i to jest jedyna rzecz o systemie,
 * jaką strona wie bez pytania powłoki. Gdy nie wiadomo, oddajemy pusty napis:
 * niewiedza NIE jest podstawą do odrzucenia wydania.
 */
function systemBiezacy(): string {
  const nawigator = (globalThis as { navigator?: { userAgent?: string; platform?: string } })
    .navigator;
  const opis = `${nawigator?.userAgent ?? ''} ${nawigator?.platform ?? ''}`.toLowerCase();
  if (opis.includes('windows') || opis.includes('win32') || opis.includes('win64')) return 'windows';
  if (opis.includes('mac')) return 'macos';
  if (opis.includes('linux') || opis.includes('x11')) return 'linux';
  return '';
}

/**
 * Czy wydanie jest przeznaczone na TEN system.
 *
 * Pole `system` deklaruje `wydania.json` i rozróżnia nim pliki strona „Pobierz".
 * Bez tego sprawdzenia baner na Linuksie podałby plik `.exe` dla Windowsa,
 * a powłoka podstawiłaby go w miejsce AppImage.
 *
 * Wydanie bez pola `system` przechodzi: wykaz jednosystemowy jest zgodny
 * z takim kształtem pliku, a brak deklaracji to brak wiedzy, nie deklaracja
 * obcego systemu.
 */
function wydanieNaTenSystem(wydanie: Wydanie): boolean {
  const zadeklarowany = (wydanie.system ?? '').trim().toLowerCase();
  if (zadeklarowany === '') return true;
  const tutejszy = systemBiezacy();
  if (tutejszy === '') return true;
  return zadeklarowany === tutejszy;
}

/**
 * Wynik odczytu kanału — sześć różnych prawd o świecie, sześć różnych wartości.
 *
 * Rozróżnienie istnieje obok `nowszeWydanie()`, bo okno wykazu musi umieć
 * powiedzieć co innego przy zerwanym łączu, a co innego przy kanale, który
 * wprost deklaruje brak wydań. Zwinięte do jednego `null` obie sytuacje
 * wyglądałyby dla okna tak samo i przy braku sieci pisałoby ono „nie ma jeszcze
 * żadnego wydania".
 *
 * Nazwy stanów (`brak-lacznosci`, `odpowiedz-serwera`) są wspólne z kodami
 * powłoki w `budowa/desktop/src-tauri/src/aktualizacja/pobranie.rs` — ta sama
 * rzecz nazywa się tak samo po obu stronach mostu.
 *
 * `brak-wydania-pod-adresem` stoi osobno od `odpowiedz-serwera`: „serwer
 * odpowiedział, że tego pliku nie ma" to informacja o kanale, a nie o łączu
 * Operatora, i nie wolno mu wtedy kazać sprawdzać połączenia.
 */
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
    // Rzut z `fetch` znaczy „nie doszło do nikogo" — i tylko to. Odpowiedź
    // serwera, choćby najgorsza, tędy nie przechodzi.
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
    // Odpowiedź przyszła, ale nie jest JSON-em. To NIE jest brak łączności
    // i nie wolno o tym mówić Operatorowi zdaniem o sieci.
    return { stan: 'odpowiedz-nieczytelna' };
  }

  const wykaz = (tresc as { wydania?: unknown } | null)?.wydania;
  if (!Array.isArray(wykaz)) return { stan: 'odpowiedz-nieczytelna' };
  // Tablica pusta to deklaracja kanału („nie ma jeszcze żadnego wydania"),
  // a nie brak odpowiedzi. Witryna mówi to samo tym samym plikiem.
  if (wykaz.length === 0) return { stan: 'pusty' };

  // Wykaz z sieci jest treścią obcą — pozycje bez `wersja`/`data` odpadają.
  const wydania = wykaz.filter(czyWydanie);
  // Coś w wykazie było, ale nic z tego nie jest wydaniem: kanał odpowiada
  // czymś, czego nie umiemy przeczytać. To bliżej „nieczytelnej" niż „pustej".
  if (wydania.length === 0) return { stan: 'odpowiedz-nieczytelna' };

  return { stan: 'wykaz', wydania };
}

/**
 * Pobiera wykaz i oddaje wydanie NOWSZE niż zainstalowane albo `null`.
 *
 * Pierwsza pozycja wykazu jest najnowsza (tak składa go witryna), ale nie
 * ufamy kolejności — porównanie idzie po wszystkich pozycjach. Plik wydania
 * musi być wskazany: wydanie bez pliku jest wpisem historycznym, a nie czymś,
 * co da się zainstalować.
 *
 * Nakładka na `pobierzWykazWydan()` gubiąca rozróżnienie stanów: baner pyta
 * o jedno — „czy jest co zakładać". Brak sieci, 404 i pusty wykaz odpowiadają
 * na to tak samo, a pas ma się wtedy nie pokazać. Kto potrzebuje zdania
 * o świecie, woła `pobierzWykazWydan()` wprost.
 */
export async function nowszeWydanie(
  wersjaBiezaca: string,
  adres: string = ADRES_WYDAN,
): Promise<Wydanie | null> {
  const odczyt = await pobierzWykazWydan(adres);
  if (odczyt.stan !== 'wykaz') return null;

  let najnowsze: Wydanie | null = null;
  for (const pozycja of odczyt.wydania) {
    if (!pozycja.plik) continue;
    // Wydanie na inny system nie jest kandydatem w ogóle — pominięcie musi
    // więc iść przed porównaniem wersji.
    if (!wydanieNaTenSystem(pozycja)) continue;
    if (!nowszaWersja(pozycja.wersja, wersjaBiezaca)) continue;
    if (najnowsze === null || nowszaWersja(pozycja.wersja, najnowsze.wersja)) {
      najnowsze = pozycja;
    }
  }
  return najnowsze;
}
