import type { Wynik } from '../protokol/kanal';
import type { WpisSesji } from './zrodlo-sesji';

/** Plik niesie dwa zdania wspólne czynnościom historii sesji: nazwanie sesji na potrzeby meldunku oraz nazwanie treści odmowy rdzenia. */
/** Treść odmowy z rdzenia albo zdanie zastępcze nazywające czynność, gdy rdzeń odmówił bez treści odmowy. */
export function odmowa(wynik: Wynik<unknown>, czynnosc: string): string {
  return wynik.blad?.message ?? `Rdzeń odmówił: ${czynnosc}.`;
}

/** Nazwa sesji do zdania meldunku pochodzi z tytułu nadanego przez rdzeń, a bez tytułu — z identyfikatora sesji. */
export function nazwa(wpis: WpisSesji): string {
  return wpis.sesja.title ?? wpis.sesja.id;
}
