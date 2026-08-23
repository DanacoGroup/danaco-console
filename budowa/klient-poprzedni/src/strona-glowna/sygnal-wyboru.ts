/**
 * Sygnał wyboru — nośnik zdarzenia widoku strony głównej.
 *
 * Rozsyła wybór dokonany w widoku do warstwy, która wie, co z nim zrobić.
 * Strona główna nie otwiera środowiska i nie wysyła komendy — skutek wyboru
 * należy do odbiorcy sygnału.
 *
 * Dlaczego własny sygnał, a nie `CustomEvent` na elemencie: zdarzenie DOM
 * niesie ładunek typu `any` i gubi typ wyboru na granicy `detail`. Tu wybór
 * pozostaje w pełni typowany aż do słuchacza.
 */

/** Odbiorca wyboru. */
export type SluchaczWyboru<T> = (wybor: T) => void;

export interface SygnalWyboru<T> {
  /** Dopisuje odbiorcę. Odbiorców może być wielu — kolejność dopisania. */
  sluchaj(sluchacz: SluchaczWyboru<T>): void;
  /** Rozsyła wybór do wszystkich dopisanych odbiorców. */
  nadaj(wybor: T): void;
}

/**
 * Buduje pusty sygnał.
 *
 * Brak odbiorcy nie jest błędem i niczego nie wstrzymuje: element pozostaje
 * klikalny, wybór po prostu nie ma jeszcze adresata.
 */
export function utworzSygnalWyboru<T>(): SygnalWyboru<T> {
  const sluchacze: SluchaczWyboru<T>[] = [];

  return {
    sluchaj(sluchacz) {
      sluchacze.push(sluchacz);
    },
    nadaj(wybor) {
      // Kopia wykazu: odbiorca dopisujący kolejnego odbiorcę nie zmienia
      // przebiegu trwającego rozesłania.
      for (const sluchacz of [...sluchacze]) {
        sluchacz(wybor);
      }
    },
  };
}
