import { przestawWTablicy } from './przeciaganie-kolejnosci';

/**
 * Kolejność ustawiona ręką — nakładka na kolejność przychodzącą z rdzenia.
 *
 * Kontrakt (`shared/contract.json`) nie ma w obszarze `session.*` komendy
 * porządku, a struktura `Session` nie ma pola kolejności. Kolejność ustawiona
 * przeciągnięciem jest więc miejscowa: nie jedzie do rdzenia, nie widzi jej
 * drugie urządzenie i nie przeżyje wyczyszczenia magazynu przeglądarki.
 * `zdanie()` niesie tę informację do widoku, żeby interfejs nie obiecywał
 * trwałości, której nie ma.
 *
 * Porządek trzyma nakładka, a nie zapis w miejscu, bo `ustawMigawke`
 * w `powloka/karty-sesji.ts` przy każdym zdarzeniu `session.changed`
 * bezwarunkowo przepisuje kolejność z rdzenia — przestawienie zapisane
 * w miejscu żyłoby do pierwszego takiego zdarzenia.
 */

/** Nakładka kolejności dla pasa pozycji rozpoznawanych po identyfikatorze. */
export interface KolejnoscMiejscowa<T> {
  /**
   * Ustawia wpisy w kolejności zapamiętanej ręką.
   *
   * Wpis spoza nakładki (nowa sesja) idzie na koniec, żeby nie przeskakiwał
   * pozycji już poukładanych ręką. Wpis, którego rdzeń już nie podaje, wypada —
   * nakładka nie wskrzesza sesji, których nie ma.
   */
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

/** Ustawienia nakładki; `klucz` i `identyfikator` są wymagane. */
export interface OpcjeKolejnosci<T> {
  /** Klucz zapisu — nazwa własna pasa, żeby dwa pasy się nie mieszały. */
  klucz: string;
  /** Identyfikator wpisu; po nim nakładka rozpoznaje pozycje między migawkami. */
  identyfikator(wpis: T): string;
  /**
   * Magazyn trwały. Pominięty znaczy `localStorage`, a `null` — brak zapisu:
   * kolejność żyje wtedy do przeładowania strony i mówi to zdaniem.
   */
  magazyn?: Storage | null;
}

/** Zdanie o kolejności wziętej wprost z rdzenia. */
export const ZDANIE_KOLEJNOSCI_RDZENIA =
  'Kolejność kart jest ta, którą podaje rdzeń.';

/** Zdanie o kolejności ustawionej ręką — mówi, dokąd ona nie sięga. */
export const ZDANIE_KOLEJNOSCI_MIEJSCOWEJ =
  'Kolejność kart ustawiona ręcznie jest MIEJSCOWA: zostaje na tym '
  + 'urządzeniu i w tej przeglądarce. Rdzeń jej nie zna — kontrakt nie ma '
  + 'komendy porządku sesji, a session.list nie zwraca pola kolejności.';

/** Zdanie o kolejności ręcznej bez magazynu — znika przy przeładowaniu. */
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
      // Zapis uszkodzony nie wywraca pasa: kolejność wraca do
      // porządku rdzenia, a to jest stan poprawny, nie awaria.
      return [];
    }
  }

  function zapisz(): void {
    if (magazyn === null) return;
    try {
      if (porzadek.length === 0) magazyn.removeItem(klucz);
      else magazyn.setItem(klucz, JSON.stringify(porzadek));
    } catch {
      // Magazyn pełny albo zablokowany. Kolejność zostaje w pamięci strony
      // i przestaje być trwała; zdanie dla widoku i tak mówi, że jest miejscowa.
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
