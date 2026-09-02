// Polecenia powłoki desktopowej wystawione wiązaniom okien. Nazwy poleceń
// i zdarzenia pochodzą z `desktop/src-tauri/src/polecenia.rs`
// i `dialog_katalogu.rs`; poza powłoką żadne z nich nie odpowiada.
import { invoke, isTauri } from '@tauri-apps/api/core';
import { listen } from '@tauri-apps/api/event';

/** Zdarzenie powłoki niosące katalog wskazany z zasobnika. */
export const ZDARZENIE_KATALOGU = 'powloka:katalog-roboczy';

/** Stan wskazania rdzenia; nazwy pól jak w `wskazanie.rs`, bo serde ich nie przemianowuje. */
export interface WskazanieRdzenia {
  schemat: string;
  host: string | null;
  port: number;
  adres: string | null;
  warstwa: string;
}

/** Opis stanu rdzenia; nazwy pól jak w `rdzen/stan.rs`. */
export interface StanRdzenia {
  pracuje: boolean;
  adres: string | null;
  dziennik: string;
  opis: string;
}

/** Przebieg udanej aktualizacji powłoki; nazwy pól jak w `aktualizacja/mod.rs`. */
export interface PrzebiegAktualizacji {
  droga: string;
  zalozone: string;
  bajtow: number;
  suma_sha256: string;
  restart_za_ms: number;
  zdanie: string;
}

/** Odmowa polecenia powłoki: kod powodu do rozgałęzienia i zdanie dla Operatora. */
export interface OdmowaPowloki {
  powod: string;
  zdanie: string;
}

/** Czy okno stoi w powłoce desktopowej; poza nią poleceń nie ma do kogo skierować. */
export function powlokaStoi(): boolean {
  return isTauri();
}

/** Otwiera natywne okno wyboru katalogu; pustka znaczy rezygnację Operatora. */
export async function wybierzKatalogRoboczy(tytul?: string): Promise<string> {
  if (!isTauri()) return '';
  const wskazany = await invoke<string | null>('wybierz_katalog_roboczy', { tytul: tytul ?? null });
  return wskazany ?? '';
}

/** Stan rdzenia w chwili pytania; pustka znaczy okno poza powłoką. */
export async function stanRdzenia(): Promise<StanRdzenia | null> {
  if (!isTauri()) return null;
  return invoke<StanRdzenia>('stan_rdzenia');
}

/** Adres rdzenia obowiązujący dla powłoki; pustka znaczy brak wskazania. */
export async function adresRdzenia(): Promise<string> {
  if (!isTauri()) return '';
  return (await invoke<string | null>('adres_rdzenia')) ?? '';
}

/** Wskazanie rdzenia obowiązujące wraz z warstwą, z której pochodzi. */
export async function wskazanieRdzenia(): Promise<WskazanieRdzenia | null> {
  if (!isTauri()) return null;
  return invoke<WskazanieRdzenia>('wskazanie_rdzenia');
}

/** Przyjmuje wskazanie Operatora; odmowa wraca wyjątkiem o kształcie `OdmowaPowloki`. */
export async function wskazRdzen(adres: string): Promise<WskazanieRdzenia | null> {
  if (!isTauri()) return null;
  return invoke<WskazanieRdzenia>('wskaz_rdzen', { adres });
}

/** Pobiera wskazane wydanie powłoki, sprawdza sumę i zakłada je; powłoka staje potem na nowo. */
export async function wykonajAktualizacje(
  adres: string,
  sumaSha256: string,
): Promise<PrzebiegAktualizacji | null> {
  if (!isTauri()) return null;
  return invoke<PrzebiegAktualizacji>('wykonaj_aktualizacje', { adres, sumaSha256 });
}

/**
 * Nasłuch katalogu wskazanego z zasobnika powłoki. Zwraca odsubskrybowanie;
 * poza powłoką nasłuch nie powstaje, bo zdarzenia nie ma kto nadać.
 */
export async function naKatalogRoboczy(
  sluchacz: (sciezka: string) => void,
): Promise<() => void> {
  if (!isTauri()) return () => undefined;
  return listen<string>(ZDARZENIE_KATALOGU, (zdarzenie) => sluchacz(zdarzenie.payload));
}
