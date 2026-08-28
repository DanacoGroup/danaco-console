import { przestawWTablicy } from './przeciaganie-kolejnosci';

/**
 * Kolejność ustawiona ręką jest nakładką miejscową na kolejność przychodzącą z rdzenia, bo kontrakt nie ma komendy porządku ani pola kolejności, więc przeciągnięcie nie jedzie do rdzenia i nie przeżywa czyszczenia magazynu przeglądarki.
 */
export interface KolejnoscMiejscowa<T> {
  /** Ustawia wpisy w kolejności zapamiętanej ręką; wpis spoza nakładki idzie na koniec. */
  uporzadkuj(wpisy: readonly T[]): T[];
  /** Zapamiętuje przestawienie pozycji `z` na miejsce `na` w podanym porządku. */
  przestaw(wpisy: readonly T[], z: number, na: number): T[];
  /** Czy kolejność jest dziś ustawiona ręką, a nie brana wprost z rdzenia. */
  czyWlasna(): boolean;
  /** Zrzeka się kolejności własnej — pas wraca do porządku rdzenia. */
  zapomnij(): void;
  /** Zdanie dla widoku o tym, czym ta kolejność jest. Nigdy puste. */
  zdanie(): string;
}

/** Ustawienia nakładki kolejności miejscowej; pola klucz oraz identyfikator wpisu są wymagane do jej działania. */
export interface OpcjeKolejnosci<T> {
  /** Klucz zapisu — nazwa własna pasa, żeby dwa pasy się nie mieszały. */
  klucz: string;
  /** Identyfikator wpisu; po nim nakładka rozpoznaje pozycje między migawkami. */
  identyfikator(wpis: T): string;
  /** Magazyn trwały; pominięty znaczy pamięć lokalną przeglądarki, wartość pusta brak zapisu. */
  magazyn?: Storage | null;
}

/** Zdanie o kolejności wziętej wprost z rdzenia, pokazywane wtedy, gdy nakładka ręczna nie jest jeszcze czynna. */
export const ZDANIE_KOLEJNOSCI_RDZENIA =
  'Kolejność kart jest ta, którą podaje rdzeń.';

/** Zdanie o kolejności ustawionej ręką — mówi wprost, dokąd nakładka miejscowa sięga, a dokąd już nie sięga. */
export const ZDANIE_KOLEJNOSCI_MIEJSCOWEJ =
  'Kolejność kart ustawiona ręcznie jest MIEJSCOWA: zostaje na tym '
  + 'urządzeniu i w tej przeglądarce. Rdzeń jej nie zna — kontrakt nie ma '
  + 'komendy porządku sesji, a session.list nie zwraca pola kolejności.';

/** Zdanie o kolejności ręcznej bez magazynu trwałego — znika przy najbliższym przeładowaniu całej strony. */
export const ZDANIE_KOLEJNOSCI_ULOTNEJ =
  'Kolejność kart ustawiona ręcznie znika przy przeładowaniu okna: magazyn '
  + 'trwały jest niedostępny, a rdzeń kolejności nie przechowuje.';

export function utworzKolejnoscMiejscowa<T>(
  opcje: OpcjeKolejnosci<T>,
): KolejnoscMiejscowa<T> {
  const magazyn = opcje.magazyn === undefined ? magazynDomyslny() : opcje.magazyn;
  const klucz = `dn.kolejnosc.${opcje.klucz}`;

  let porzadek: string[] = wczytaj();

  function wczytaj(): string[] {
    if (magazyn === null) return [];
    try {
      const zapis = magazyn.getItem(klucz);
      if (zapis === null) return [];
      const odczytane: unknown = JSON.parse(zapis);
      if (!Array.isArray(odczytane)) return [];
      return odczytane.filter((wpis): wpis is string => typeof wpis === 'string');
    } catch {
      // Zapis uszkodzony nie wywraca pasa: kolejność wraca do porządku rdzenia, co jest stanem poprawnym.
      return [];
    }
  }

  function zapisz(): void {
    if (magazyn === null) return;
    try {
      if (porzadek.length === 0) magazyn.removeItem(klucz);
      else magazyn.setItem(klucz, JSON.stringify(porzadek));
    } catch {
      // Magazyn pełny albo zablokowany zostawia kolejność w pamięci strony, przestając być trwały.
    }
  }

  function uporzadkuj(wpisy: readonly T[]): T[] {
    if (porzadek.length === 0) return [...wpisy];

    const wedlugId = new Map(wpisy.map((wpis) => [opcje.identyfikator(wpis), wpis]));
    const znane: T[] = [];
    for (const id of porzadek) {
      const wpis = wedlugId.get(id);
      if (wpis === undefined) continue;
      znane.push(wpis);
      wedlugId.delete(id);
    }
    // Reszta zachowuje kolejność rdzenia i idzie na koniec.
    return [...znane, ...wedlugId.values()];
  }

  return {
    uporzadkuj,

    przestaw(wpisy, z, na) {
      const ustawione = przestawWTablicy(uporzadkuj(wpisy), z, na);
      porzadek = ustawione.map((wpis) => opcje.identyfikator(wpis));
      zapisz();
      return ustawione;
    },

    czyWlasna: () => porzadek.length > 0,

    zapomnij() {
      porzadek = [];
      zapisz();
    },

    zdanie() {
      if (porzadek.length === 0) return ZDANIE_KOLEJNOSCI_RDZENIA;
      return magazyn === null ? ZDANIE_KOLEJNOSCI_ULOTNEJ : ZDANIE_KOLEJNOSCI_MIEJSCOWEJ;
    },
  };
}

/**
 * Magazyn trwały przeglądarki albo `null`.
 *
 * `localStorage` bywa niedostępny (tryb prywatny, zablokowane ciasteczka
 * witryny), a samo sięgnięcie po niego potrafi rzucić wyjątkiem — stąd próba,
 * a nie sprawdzenie obecności.
 */
function magazynDomyslny(): Storage | null {
  try {
    return globalThis.localStorage ?? null;
  } catch {
    return null;
  }
}
