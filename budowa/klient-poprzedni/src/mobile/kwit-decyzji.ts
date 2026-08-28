/**
 * Kwit decyzji jest dowodem, że decyzja Operatora doszła do rdzenia, i przeżywa
 * zamknięcie telefonu w pamięci trwałej przeglądarki. Niesie zdanie rdzenia:
 * nazwę komendy, rozstrzygnięcie wywołania i stan po zmianie.
 */

/**
 * Klucz, pod którym magazyn trzyma wykaz kwitów w pamięci trwałej przeglądarki.
 * Przestrzeń nazw `danaco-console.mobile` oddziela zapis kanału telefonicznego od
 * pozostałych zapisów aplikacji na tym samym urządzeniu.
 */
export const KLUCZ_KWITOW = 'danaco-console.mobile.kwity';

/**
 * Ile kwitów zostaje w pamięci; przy dopisaniu starsze wypadają poza ten próg.
 * Telefon czyta ostatnie interwencje, a nie archiwum, więc zapis nie rośnie
 * w nieskończoność i mieści się w pojemności pamięci przeglądarki.
 */
export const ILE_KWITOW = 20;

/**
 * Droga interwencji, której kwit dotyczy: zatwierdzenie kroku, wstrzymanie,
 * nastawienie koordynatora albo przejęcie sterowania. Wartości są zapisywane
 * do pamięci wprost, więc stanowią część zapisu odczytywanego po ponownym
 * otwarciu telefonu.
 */
export const Droga = {
  Zatwierdz: 'zatwierdz-krok',
  Wstrzymaj: 'wstrzymaj',
  Nastaw: 'nastaw-koordynatora',
  Przejmij: 'przejmij-sterowanie',
} as const;
export type Droga = (typeof Droga)[keyof typeof Droga];

/**
 * Jedno wywołanie kontraktu wykonane w ramach drogi interwencji, wraz z jego
 * rozstrzygnięciem. Droga składa się z kolejnych kroków, a kwit zapisuje je
 * wszystkie, żeby widać było, na którym wywołaniu rdzeń odmówił.
 */
export interface KrokKwitu {
  /** Nazwa komendy kontraktu — dokładnie ta, która poszła na gniazdo. */
  komenda: string;
  udany: boolean;
  /** Stan po zmianie odczytany z odpowiedzi rdzenia; puste, gdy odmówił. */
  stanPo?: string;
  /** Treść odmowy rdzenia złożona przez `komponenty/odmowa.ts`. */
  odmowa?: string;
}

/**
 * Kwit jednej drogi interwencji: własne oznaczenie, droga, nagłówek pozycji,
 * wykaz wykonanych kroków, rozstrzygnięcie całości, zdanie dla Operatora oraz
 * chwila wykonania. Taki zapis idzie do pamięci trwałej i z niej wraca.
 */
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

/**
 * Magazyn kwitów odczytuje wykaz, dopisuje kwit i czyści zapis. Jest wstrzykiwany
 * do drogi interwencji, więc sprawdzian podstawia atrapę i obywa się bez
 * pamięci trwałej przeglądarki.
 */
export interface MagazynKwitow {
  odczytaj(): Kwit[];
  dopisz(kwit: Kwit): Kwit[];
  wyczysc(): void;
}

/**
 * Oddaje pamięć trwałą przeglądarki albo `null`, gdy jej nie ma: pod Node oraz
 * w oknie prywatnym samo sięgnięcie po `window.localStorage` kończy się błędem,
 * więc odczyt idzie w bloku przechwytującym wyjątek.
 */
export function pamiecPrzegladarki(): Storage | null {
  try {
    return window.localStorage;
  } catch {
    return null;
  }
}

/**
 * Składa magazyn kwitów osadzony na dowolnym `Storage`. Wartość `null` daje
 * magazyn pusty: odczyt oddaje pustkę, dopisanie oddaje sam dopisany kwit
 * i niczego nie utrwala.
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

/**
 * Składa zdanie kwitu widziane przez Operatora: godzinę wykonania w formacie
 * lokalnym wraz ze zdaniem złożonym ze stanów oddanych przez rdzeń. Godzina
 * pochodzi z chwili zapisanej w kwicie, nie z chwili odczytu.
 */
export function zdanieKwitu(kwit: Kwit): string {
  const godzina = new Date(kwit.o).toLocaleTimeString('pl-PL', {
    hour: '2-digit',
    minute: '2-digit',
  });
  return `${godzina} — ${kwit.zdanie}`;
}

/**
 * Sprawdza, czy odczyt z pamięci ma kształt kwitu. Zapis mógł zostać podmieniony
 * albo pochodzić z wcześniejszego kształtu danych, więc wartość spoza wzorca
 * jest pomijana zamiast wchodzić do wykazu.
 */
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
