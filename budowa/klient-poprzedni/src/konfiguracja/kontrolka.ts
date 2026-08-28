import type { SettingDefinition } from '../../../shared/contract';

/**
 * Kontrolka pola formularza — jeden kształt dla każdego rodzaju wartości:
 * element do osadzenia, odczyt wartości w postaci wysyłanej komendą
 * `config.set`, zapis wartości nadchodzącej z rdzenia oraz zgłoszenie
 * zamiaru zapisu.
 */
export interface Kontrolka {
  /** Element osadzany w wierszu pola. */
  element: HTMLElement;
  /** Wartość w postaci przyjmowanej przez `config.set`. */
  odczytaj(): unknown;
  /** Ustawia wartość obowiązującą — z rdzenia albo z wartości domyślnej. */
  ustaw(wartosc: unknown): void;
  /** Zgłasza zamiar zapisu; wywoływane po zatwierdzeniu wartości. */
  naZatwierdzenie(sluchacz: () => void): void;
  /** Ostrzeżenie kontrolki; pusty napis znaczy brak ostrzeżenia. */
  ostrzezenie?: string;
}

/**
 * Zależności budowniczego kontrolki: pozycja katalogu, z której powstaje
 * kontrolka, oraz identyfikator elementu wiążący etykietę pola z kontrolką.
 */
export interface ZaleznosciKontrolki {
  /** Pozycja katalogu, z której powstaje kontrolka. */
  definicja: SettingDefinition;
  /** Identyfikator elementu wiążący etykietę z kontrolką. */
  identyfikator: string;
}

/**
 * Zbiornik zamiarów zapisu wspólny dla wszystkich kontrolek.
 *
 * Kontrolka nie wysyła komendy sama — zgłasza, że wartość została
 * zatwierdzona. Co z tym zrobić, rozstrzyga wiersz pola: on zna adres zapisu
 * i wskaźnik zasięgu.
 */
export function utworzZbiornikZamiarow(): {
  zglos(): void;
  naZatwierdzenie(sluchacz: () => void): void;
} {
  const sluchacze: Array<() => void> = [];
  return {
    zglos: () => {
      for (const sluchacz of [...sluchacze]) sluchacz();
    },
    naZatwierdzenie: (sluchacz) => void sluchacze.push(sluchacz),
  };
}

/**
 * Napis odczytany z wartości nieznanego kształtu. Wartości `null` oraz
 * `undefined` dają napis pusty, napis wraca bez zmiany, liczba i wartość
 * logiczna przechodzą przez konwersję, a pozostałe kształty przez zapis JSON.
 */
export function napis(wartosc: unknown): string {
  if (wartosc === null || wartosc === undefined) return '';
  if (typeof wartosc === 'string') return wartosc;
  if (typeof wartosc === 'number' || typeof wartosc === 'boolean') return String(wartosc);
  return bezpiecznyZapis(wartosc);
}

/**
 * Zapis JSON odporny na strukturę cykliczną. Wartość zostaje zapisana
 * z wcięciem dwóch spacji, a gdy zapis się nie powiedzie, wynikiem jest
 * tekstowa postać wartości.
 */
export function bezpiecznyZapis(wartosc: unknown): string {
  try {
    return JSON.stringify(wartosc, null, 2) ?? '';
  } catch {
    return String(wartosc);
  }
}

/**
 * Lista napisów odczytana z wartości nieznanego kształtu. Tablica przechodzi
 * element po elemencie, wartość pojedyncza daje listę jednoelementową,
 * a napis pusty daje listę pustą.
 */
export function lista(wartosc: unknown): string[] {
  if (Array.isArray(wartosc)) return wartosc.map((pozycja) => napis(pozycja));
  const pojedyncza = napis(wartosc);
  return pojedyncza === '' ? [] : [pojedyncza];
}
