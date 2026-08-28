/** Sygnał wyboru rozsyła wybór dokonany w widoku do warstwy, która wie, co z nim zrobić, zamiast otwierać środowisko albo wysyłać komendę bezpośrednio. */
export type SluchaczWyboru<T> = (wybor: T) => void;

export interface SygnalWyboru<T> {
  /** Dopisuje odbiorcę. Odbiorców może być wielu — kolejność dopisania. */
  sluchaj(sluchacz: SluchaczWyboru<T>): void;
  /** Rozsyła wybór do wszystkich dopisanych odbiorców. */
  nadaj(wybor: T): void;
}

/** Buduje pusty sygnał wyboru: brak odbiorcy nie jest błędem i niczego nie wstrzymuje, wybór po prostu nie ma jeszcze adresata. */
export function utworzSygnalWyboru<T>(): SygnalWyboru<T> {
  const sluchacze: SluchaczWyboru<T>[] = [];

  return {
    sluchaj(sluchacz) {
      sluchacze.push(sluchacz);
    },
    nadaj(wybor) {
      // Kopia wykazu odbiorców, żeby dopisanie kolejnego w trakcie rozesłania nie zmieniło jego przebiegu.
      for (const sluchacz of [...sluchacze]) {
        sluchacz(wybor);
      }
    },
  };
}
