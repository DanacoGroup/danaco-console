/** Pamięć zwinięcia stref i sekcji trzyma jedną nastawę widoku w pamięci przeglądarki na tym urządzeniu, wspólną dla obu powierzchni wejściowych. */
const PRZEDROSTEK = 'dn.zwiniecie.';

/** Rozstrzyga, czy strefa o podanym kluczu pamięci ma być rozwinięta, na podstawie zapisanej postaci albo wartości domyślnej. */
export function czyRozwiniete(klucz: string, domyslnie: boolean): boolean {
  try {
    const zapis = globalThis.localStorage?.getItem(PRZEDROSTEK + klucz);
    if (zapis === null || zapis === undefined) return domyslnie;
    return zapis === '1';
  } catch {
    return domyslnie;
  }
}

/** Zapamiętuje postać strefy na tym urządzeniu pod jej kluczem pamięci, bez zgłaszania błędu przy niedostępności pamięci przeglądarki. */
export function zapamietajZwiniecie(klucz: string, rozwiniete: boolean): void {
  try {
    globalThis.localStorage?.setItem(PRZEDROSTEK + klucz, rozwiniete ? '1' : '0');
  } catch {
    // Nastawa widoku nie jest powodem, żeby cokolwiek przerywać.
  }
}
