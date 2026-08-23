/**
 * Kwit decyzji — dowód, że decyzja Operatora doszła do rdzenia, przeżywający
 * zamknięcie telefonu.
 *
 * Telefon jest kanałem interwencji, nie miejscem pracy: Operator podejmuje
 * jedną decyzję i odchodzi od urządzenia. Potwierdzenie wypisane w oknie ginie
 * razem z oknem, a pytanie „czy tamto przeszło?” wraca później, bez dostępu do
 * pulpitu. Kwit leży w pamięci przeglądarki, więc pierwsze, co telefon pokazuje
 * po ponownym otwarciu, to zdanie w rodzaju „12:03 — wstrzymano kolejkę
 * «przekazanie: krok pierwszy»; rdzeń potwierdził stan paused”.
 *
 * Kwit niesie zdanie rdzenia, nie widoku. Każdy krok zapisuje nazwę komendy,
 * rozstrzygnięcie wywołania i stan po zmianie wyjęty z odpowiedzi rdzenia
 * (stan kolejki, `stopped`, źródło konfiguracji). Kwit nieudany zapisuje treść
 * odmowy, bo niepowodzenie także jest wiadomością, na którą Operator czeka.
 *
 * Magazyn jest wstrzykiwany, wzorem `uwierzytelnienie/sesja-bramki.ts`:
 * domyślnie `window.localStorage`, w sprawdzianie — atrapa. Błąd pamięci nie
 * zatrzymuje drogi interwencji: brak magazynu znaczy „kwitu nie zapiszę”,
 * a nie „decyzji nie wykonam”. Stąd kolejność: najpierw rdzeń, potem kwit.
 */

/** Klucz zapisu kwitów w pamięci przeglądarki. */
export const KLUCZ_KWITOW = 'danaco-console.mobile.kwity';

/** Ile kwitów zostaje; starsze wypadają. Telefon czyta ostatnie, nie archiwum. */
export const ILE_KWITOW = 20;

/** Droga interwencji, której kwit dotyczy. */
export const Droga = {
  Zatwierdz: 'zatwierdz-krok',
  Wstrzymaj: 'wstrzymaj',
  Nastaw: 'nastaw-koordynatora',
  Przejmij: 'przejmij-sterowanie',
} as const;
export type Droga = (typeof Droga)[keyof typeof Droga];

/** Jedno wywołanie kontraktu wykonane w ramach drogi. */
export interface KrokKwitu {
  /** Nazwa komendy kontraktu — dokładnie ta, która poszła na gniazdo. */
  komenda: string;
  udany: boolean;
  /** Stan po zmianie odczytany z odpowiedzi rdzenia; puste, gdy odmówił. */
  stanPo?: string;
  /** Treść odmowy rdzenia złożona przez `komponenty/odmowa.ts`. */
  odmowa?: string;
}

/** Kwit jednej drogi interwencji. */
export interface Kwit {
  id: string;
  droga: Droga;
  /** Nagłówek pozycji, której decyzja dotyczyła. */
  pozycja: string;
  kroki: readonly KrokKwitu[];
  udany: boolean;
  /** Zdanie dla Operatora złożone ze stanów oddanych przez rdzeń. */
  zdanie: string;
  /** Chwila w milisekundach epoki. */
  o: number;
}

/** Magazyn kwitów; wstrzykiwany, żeby sprawdzian nie potrzebował przeglądarki. */
export interface MagazynKwitow {
  odczytaj(): Kwit[];
  dopisz(kwit: Kwit): Kwit[];
  wyczysc(): void;
}

/** Pamięć trwała przeglądarki; `null`, gdy jej nie ma (node, tryb prywatny). */
export function pamiecPrzegladarki(): Storage | null {
  try {
    return window.localStorage;
  } catch {
    return null;
  }
}

/**
 * Magazyn kwitów osadzony na dowolnym `Storage`.
 *
 * `null` daje magazyn pusty: odczyt oddaje pustkę, dopisanie oddaje sam
 * dopisany kwit i nic nie utrwala. Ekran pokazuje wtedy potwierdzenie bieżące,
 * a po zamknięciu telefonu go nie ma — taki jest stan faktyczny urządzenia bez
 * pamięci trwałej.
 */
export function utworzMagazynKwitow(pamiec: Storage | null = pamiecPrzegladarki()): MagazynKwitow {
  return {
    odczytaj() {
      if (pamiec === null) return [];
      try {
        const zapis = pamiec.getItem(KLUCZ_KWITOW);
        if (zapis === null) return [];
        const odczyt: unknown = JSON.parse(zapis);
        if (!Array.isArray(odczyt)) return [];
        return odczyt.filter(czyKwit);
      } catch {
        // Zapis nieczytelny znaczy „kwitów nie ma”, nie „aplikacja stoi”.
        return [];
      }
    },

    dopisz(kwit) {
      const wykaz = [kwit, ...this.odczytaj()].slice(0, ILE_KWITOW);
      if (pamiec === null) return wykaz;
      try {
        pamiec.setItem(KLUCZ_KWITOW, JSON.stringify(wykaz));
      } catch {
        // Pamięć pełna albo zamknięta — decyzja i tak została wykonana.
      }
      return wykaz;
    },

    wyczysc() {
      if (pamiec === null) return;
      try {
        pamiec.removeItem(KLUCZ_KWITOW);
      } catch {
        // Nic do zrobienia; brak kwitów jest stanem dopuszczalnym.
      }
    },
  };
}

/** Zdanie kwitu widziane przez Operatora — godzina, droga, stan oddany przez rdzeń. */
export function zdanieKwitu(kwit: Kwit): string {
  const godzina = new Date(kwit.o).toLocaleTimeString('pl-PL', {
    hour: '2-digit',
    minute: '2-digit',
  });
  return `${godzina} — ${kwit.zdanie}`;
}

/** Czy odczyt z pamięci ma kształt kwitu. Kształt spoza wzorca jest pomijany. */
function czyKwit(wartosc: unknown): wartosc is Kwit {
  if (typeof wartosc !== 'object' || wartosc === null) return false;
  const zapis = wartosc as Record<string, unknown>;
  return (
    typeof zapis['id'] === 'string' &&
    typeof zapis['droga'] === 'string' &&
    typeof zapis['pozycja'] === 'string' &&
    typeof zapis['zdanie'] === 'string' &&
    typeof zapis['udany'] === 'boolean' &&
    typeof zapis['o'] === 'number' &&
    Array.isArray(zapis['kroki'])
  );
}
