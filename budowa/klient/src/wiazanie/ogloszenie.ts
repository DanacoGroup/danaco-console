/**
 * Ogłoszenie dla Operatora. Jedyne miejsce, w którym klient sięga po most
 * `dnToast` biblioteki; rodzaj komunikatu nazywa wagę zdarzenia.
 */

/** Rodzaje komunikatu, które rozpoznaje biblioteka okna. */
export type RodzajOgloszenia = 'informacja' | 'ostrzezenie' | 'blad';

/** Czas stania komunikatu na ekranie w milisekundach. */
const CZAS_STANIA = 4200;

/** Ogłasza rzecz Operatorowi. Okno bez wpiętego mostu milczy — brak komunikatu nie zatrzymuje czynności. */
export function oglos(tytul: string, tresc: string, rodzaj: RodzajOgloszenia = 'informacja'): void {
  const most = globalThis as {
    dnToast?: (tytul: string, tresc: string, rodzaj: string, ms: number) => void;
  };
  most.dnToast?.(tytul, tresc, rodzaj, CZAS_STANIA);
}
