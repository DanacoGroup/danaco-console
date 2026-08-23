import { ConfigScope } from '../../../../shared/contract';
import { powod } from './nadanie-rol';
import type { StanMultitaskingu } from './stan-multitaskingu';
import type { ZrodloOkien } from './zrodlo-okien';

/**
 * Wskazanie okna analityka na poziomie sesji — odczyt i zapis.
 *
 * `WindowRole` nie ma czwartej wartości: Results Analyzer jest oknem roli
 * `standalone`, a odróżnia go wyłącznie wskazanie zapisane w ustawieniach sesji.
 * Dopisanie czwartej roli po stronie klienta rozjechałoby wykaz ról z bazą.
 *
 * Klucz `multitasking.analityk` nie stoi w katalogu ustawień rdzenia
 * (`definicja_ustawienia`), więc `config.set` odpowiada odmową
 * `validation_failed`, a wskazanie nie przeżywa odświeżenia. Zapis jest
 * sprawdzany i odmowa trafia do zdania oddawanego wołającemu.
 */

/** Klucz wskazania okna analityka na poziomie sesji. */
export const KLUCZ_ANALITYKA = 'multitasking.analityk';

/**
 * Utrwala wskazanie analityka na poziomie sesji i oddaje zdanie o wyniku:
 * odmowa rdzenia i wartość oddana inna niż wysłana są w nim nazwane wprost,
 * bo wskazanie trzyma wtedy sam widok. Źródło okien i stan wspólny bierze
 * z parametrów, panelu nie zna.
 */
export async function utrwalWskazanieAnalityka(
  zrodlo: ZrodloOkien,
  stan: StanMultitaskingu,
  idOkna: string,
): Promise<string> {
  stan.ustawAnalityka(idOkna);
  const wynik = await zrodlo.zapiszUstawienie({
    key: KLUCZ_ANALITYKA,
    value: idOkna,
    scope: ConfigScope.Session,
    scopeId: stan.sesja(),
  });
  if (!wynik.udany) {
    return `Wskazanie analityka trzyma WYŁĄCZNIE ten widok — rdzeń odmówił zapisu na poziomie sesji. ${powod(wynik.blad)}`;
  }
  const zapisane = wynik.wynik?.value;
  if (zapisane !== idOkna) {
    return `Rdzeń przyjął zapis wskazania analityka, ale oddał wartość ${String(zapisane)} zamiast ${idOkna} — wskazanie nie jest utrwalone.`;
  }
  return `Wskazanie analityka utrwalone na poziomie sesji: ${idOkna}.`;
}

/** Wskazanie analityka odczytane z ustawień sesji; puste, gdy nie wskazano. */
export async function wskazanieAnalitykaZSesji(
  zrodlo: ZrodloOkien,
  sesja: string,
): Promise<string> {
  const wskazanie = await zrodlo.ustawienia({
    key: KLUCZ_ANALITYKA,
    scope: ConfigScope.Session,
    scopeId: sesja,
  });
  const wpis = wskazanie.udany ? wskazanie.wynik?.[0]?.value : undefined;
  return typeof wpis === 'string' ? wpis : '';
}
