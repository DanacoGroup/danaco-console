import type { Wynik } from '../protokol/kanal';
import type { WpisSesji } from './zrodlo-sesji';

/**
 * Dwa zdania wspólne wszystkim czynnościom historii sesji: nazwanie sesji
 * i nazwanie odmowy. Wspólne miejsce trzyma jedną odpowiedź na oba pytania.
 *
 * Nazwa sesji pochodzi z rdzenia, a gdy rdzeń tytułu nie nadał — z identyfikatora.
 * Odmowa mówi treścią rdzenia; zdanie zastępcze wchodzi tylko wtedy, gdy odmowa
 * przyszła bez treści, i nazywa wówczas czynność, której dotyczyła.
 */

/** Treść odmowy z rdzenia albo zdanie zastępcze nazywające czynność. */
export function odmowa(wynik: Wynik<unknown>, czynnosc: string): string {
  return wynik.blad?.message ?? `Rdzeń odmówił: ${czynnosc}.`;
}

/** Nazwa sesji do zdania meldunku — tytuł z rdzenia albo identyfikator. */
export function nazwa(wpis: WpisSesji): string {
  return wpis.sesja.title ?? wpis.sesja.id;
}
