// Napisy okna wejścia i trwały identyfikator urządzenia. Katalog treści niesie
// biblioteka `design/zasoby/okna/wejscie/tresci.js` — skrypt domknięty, nie
// moduł, więc czyta się go przez obiekt globalny.
import { nowyIdentyfikator } from '../protokol/identyfikator.ts';

interface Katalog {
  tekst(sciezka: string): unknown;
}

interface StykKatalogu {
  DanacoNarzedzia?: { zwiaz(katalog: unknown): Katalog };
  DanacoWejscie?: { tresci?: unknown };
}

/** Klucz zapisu trwałego identyfikatora urządzenia w przeglądarce powłoki. */
const KLUCZ_URZADZENIA = 'dn-urzadzenie';

let urzadzenie = '';

/** Napis z katalogu treści okna wejścia; pustka znaczy klucz spoza katalogu. */
export function napis(sciezka: string): string {
  const styk = globalThis as StykKatalogu;
  const katalog = styk.DanacoNarzedzia?.zwiaz(styk.DanacoWejscie?.tresci ?? {});
  const wartosc = katalog?.tekst(sciezka);
  return typeof wartosc === 'string' ? wartosc : '';
}

/**
 * Identyfikator urządzenia trwały między uruchomieniami. Rdzeń wiąże z nim PIN
 * i sesję bramki (`auth.method.add`, `device.list`), więc identyfikator nadawany
 * na czas jednego uruchomienia zostawiałby PIN przy maszynie, do której nikt
 * już nie wraca. Zapis niedostępny daje identyfikator na czas uruchomienia.
 */
export function urzadzenieTrwale(): string {
  if (urzadzenie !== '') return urzadzenie;
  urzadzenie = odczytajZapis() ?? nowyIdentyfikator('urzadzenie');
  zapisz(urzadzenie);
  return urzadzenie;
}

function odczytajZapis(): string | null {
  try {
    const zapisane = globalThis.localStorage?.getItem(KLUCZ_URZADZENIA) ?? '';
    return zapisane === '' ? null : zapisane;
  } catch {
    return null;
  }
}

function zapisz(wartosc: string): void {
  try {
    globalThis.localStorage?.setItem(KLUCZ_URZADZENIA, wartosc);
  } catch {
    // Zapis odmówiony przez przeglądarkę zostawia identyfikator na czas uruchomienia.
  }
}
