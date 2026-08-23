import type { SettingDefinition } from '../../../shared/contract';

/**
 * Kontrolka pola formularza — jeden kształt dla każdego rodzaju wartości.
 *
 * Formularz okna konfiguracji jest generowany z katalogu, więc nie może znać
 * żadnej kontrolki z osobna. Zna wyłącznie ten kształt: element do osadzenia,
 * odczyt wartości w postaci wysyłanej komendą `config.set`, zapis wartości
 * przyszłej z rdzenia oraz zgłoszenie zamiaru zapisu.
 *
 * Dzięki temu dodanie rodzaju wartości do kontraktu jest jednym przypadkiem
 * w module rozdzielającym (`wybor-kontrolki.ts`) i jedną funkcją budującą —
 * nie nowym ekranem i nie zmianą formularza.
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

/** Zależności budowniczego kontrolki: definicja pozycji katalogu i jej pole. */
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

/** Napis odczytany z wartości nieznanego kształtu; `null` i `undefined` znaczą pustkę. */
export function napis(wartosc: unknown): string {
  if (wartosc === null || wartosc === undefined) return '';
  if (typeof wartosc === 'string') return wartosc;
  if (typeof wartosc === 'number' || typeof wartosc === 'boolean') return String(wartosc);
  return bezpiecznyZapis(wartosc);
}

/** Zapis JSON odporny na strukturę cykliczną. */
export function bezpiecznyZapis(wartosc: unknown): string {
  try {
    return JSON.stringify(wartosc, null, 2) ?? '';
  } catch {
    return String(wartosc);
  }
}

/** Lista napisów odczytana z wartości; wartość pojedyncza daje listę jednoelementową. */
export function lista(wartosc: unknown): string[] {
  if (Array.isArray(wartosc)) return wartosc.map((pozycja) => napis(pozycja));
  const pojedyncza = napis(wartosc);
  return pojedyncza === '' ? [] : [pojedyncza];
}
