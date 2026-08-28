/**
 * Zgodność zamówienia z odpowiedzią rdzenia jest jednym miejscem na regułę, którą moduł Translate
 * stosuje w zdaniach potwierdzających.
 */

/** Jedno pole porównywane między żądaniem a odpowiedzią niesie nazwę pola, wartość zamówioną oraz wartość oddaną przez rdzeń. */
export interface PoleOdpowiedzi {
  /** Nazwa pola tak, jak nazywa je okno — nie nazwa z kontraktu. */
  nazwa: string;
  /** Wartość, którą okno wysłało w żądaniu. */
  zamowione: string;
  /** Wartość, którą rdzeń oddał; `undefined` znaczy pole nieobecne. */
  oddane: string | undefined;
}

/** Po ilu znakach cytat się urywa — zdanie ma być czytelne, nie kompletne, więc długi cytat kończy się wielokropkiem. */
const DLUGOSC_CYTATU = 60;

/**
 * Zdanie o rozbieżności albo `null`, gdy rdzeń oddał to, co zamówiono.
 *
 * Pole nieobecne w odpowiedzi liczy się jak puste: rdzeń, który nie oddał
 * pola, nie pokazał zapisu tego pola — a milczenie nie jest potwierdzeniem.
 */
export function rozbieznoscOdpowiedzi(
  czynnosc: string,
  pola: readonly PoleOdpowiedzi[],
): string | null {
  const rozbiezne = pola.filter((pole) => (pole.oddane ?? '') !== pole.zamowione);
  if (rozbiezne.length === 0) return null;
  return (
    `${czynnosc}: rdzeń przyjął żądanie, ale oddał co innego, niż zamówiono — ` +
    `${rozbiezne.map(opisPola).join('; ')}. ` +
    'Okno nie potwierdza zapisu, którego rdzeń nie pokazał; poprawne jest to, co wróciło.'
  );
}

function opisPola(pole: PoleOdpowiedzi): string {
  const wrocilo =
    pole.oddane === undefined
      ? 'pola nie ma w odpowiedzi'
      : `wróciło ${cytat(pole.oddane)}`;
  return `${pole.nazwa}: wysłano ${cytat(pole.zamowione)}, ${wrocilo}`;
}

/** Cytat wartości z ujawnionymi znakami niewidocznymi i uciętym ogonem pokazuje różnicę wprost, bez zgadywania. */
function cytat(wartosc: string): string {
  if (wartosc === '') return '(pusto)';
  const widok = ujawnij(wartosc);
  if (widok.length <= DLUGOSC_CYTATU) return `„${widok}"`;
  return `„${widok.slice(0, DLUGOSC_CYTATU)}…" (${String(wartosc.length)} znaków)`;
}

/**
 * Znaki, których Operator w cytacie nie zobaczy, wychodzą jako kod U+XXXX.
 *
 * Zwykła spacja, tabulator i przejście do nowego wiersza zostają sobą — są
 * czytelne w zdaniu i ich ujawnianie zaciemniłoby cytat zamiast go objaśnić.
 */
function ujawnij(wartosc: string): string {
  let wynik = '';
  for (const znak of wartosc) {
    wynik += czyNiewidoczny(znak) ? kodZnaku(znak) : znak;
  }
  return wynik;
}

const NIEWIDOCZNY = /[\p{C}\p{Z}]/u;

function czyNiewidoczny(znak: string): boolean {
  if (znak === ' ' || znak === '\n' || znak === '\t') return false;
  return NIEWIDOCZNY.test(znak);
}

function kodZnaku(znak: string): string {
  const punkt = znak.codePointAt(0) ?? 0;
  return `⟨U+${punkt.toString(16).toUpperCase().padStart(4, '0')}⟩`;
}
