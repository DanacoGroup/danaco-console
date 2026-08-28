import { KluczUstawieniaOkna } from './klucze-ustawien';

/**
 * Adnotacje wykonania opisują pola konfiguracji zadania, których wykonanie nie czyta wprost.
 */

/**
 * Zdanie o wyborze eksperta jest stanem spoiny, nie obietnicą, bo żadna znana ścieżka rdzenia nie potwierdza nałożenia tożsamości eksperta na wywołanie modelu.
 */
const AGENT_ZAPISYWANY =
  'Wybór eksperta: okno bierze jego model bazowy i zapisuje jego kod (rdzeń utrwala pole agentId), ale nałożenia tożsamości eksperta — warstw promptu, skilli, wtyczek — na wywołanie nie potwierdzono pomiarem.';

/**
 * Zdanie o zasięgu local jest powiedziane wprost: bez toru zwrotnego do urządzenia operatora zasięg startuje na hoście rdzenia i zostawia o tym wpis w dzienniku.
 */
const ZASIEG_LOKALNY =
  'Stan: zasięg zdalny prowadzi proces torem SSH na wskazany host, a brakujące ogniwo drogi jest nazwaną odmową. Zasięg „urządzenie Operatora" startuje dziś na hoście rdzenia — toru zwrotnego do urządzenia nie ma — i zostawia o tym wpis w dzienniku.';

/**
 * Objaśnienia sterowań kompletu, klucz — nazwa sterowania.
 *
 * Każdy element konfiguracji ma objaśnienie kontekstowe — bez wyjątków.
 */
const OBJASNIENIA: Readonly<Record<string, string>> = {
  'Środowisko wykonania': `Zasięg, w którym ma pracować proces modelu tego okna: urządzenie Operatora, host rdzenia albo host zdalny. ${ZASIEG_LOKALNY}`,
  'Host wykonania': `Nazwa serwera używana przy zasięgu zdalnym; pole przyjmuje dowolną nazwę, podpowiedzi tylko skracają drogę. Wartość dociera do wykonania: rdzeń czyta ją z bazy i składa z niej tor SSH. Klient zapisuje ją pod kluczem katalogu „${KluczUstawieniaOkna.HostWykonania}".`,
  'Moduł okna': 'Moduł, w którym okno pracuje. Zmiana idzie komendą window.update i zapisuje się w rdzeniu.',
  Model: `Model obsługujący to okno, w dwóch sekcjach: „Modele" to kanały surowe z rejestru rdzenia, „Moi agenci" to eksperci skomponowani w module Agents. Agent nie zastępuje kanału — stoi na nim jako tożsamość. ${AGENT_ZAPISYWANY}`,
  'Model zapasowy': `Kanał użyty, gdy kanał główny nie odpowiada. Ustawienie dociera do wywołania: rdzeń przekłada kod kanału na identyfikator modelu i podaje go parametrem --fallback-model. Klient zapisuje go pod kluczem katalogu „${KluczUstawieniaOkna.KanalZapasowy}".`,
  'Nakład rozumowania': `Stopień od odpowiedzi szybkiej po pełny namysł; wartość wyliczenia katalogu, nie nazwa modelu, więc zmiana kanału jej nie unieważnia. Ustawienie dociera do wywołania: rdzeń przekłada je na parametr --effort. Klient zapisuje go pod kluczem katalogu „${KluczUstawieniaOkna.NakladRozumowania}".`,
  'Tryb uprawnień': 'Ile wolno modelowi bez pytania. Ustawienie dociera do wywołania: rdzeń przekłada je na parametr --permission-mode programu kanału.',
  'Rola okna': 'Rola okna w scenie okien równoległych: koordynator albo wykonawca. Zmiana idzie komendą window.update.',
  'Katalogi robocze': 'Katalogi, które model widzi w tym oknie. Ustawienie dociera do wywołania: rdzeń przekłada je na parametry --add-dir programu kanału.',
};

/**
 * Objaśnienie sterowania o wskazanej nazwie.
 *
 * Sterowanie bez wpisu nie zostaje bez objaśnienia — dostaje zdanie mówiące,
 * że objaśnienia brak. Milczenie byłoby gorsze niż przyznanie się do luki.
 */
export function objasnienieSterowania(nazwa: string): string {
  return OBJASNIENIA[nazwa] ?? `Sterowanie „${nazwa}" nie ma jeszcze objaśnienia w wykazie.`;
}

/**
 * Wykaz kluczy konfiguracji bez konsumenta w rdzeniu jest pusty, bo droga dopisania ma być ta sama co droga wycofania klucza.
 */
const KLUCZE_BEZ_KONSUMENTA: readonly string[] = [];

/**
 * Przedrostek „harness." obejmuje całą rodzinę kluczy programu kanału
 * (program_claude, plik_ustawien, konfiguracja_mcp): katalog je zna, żadna
 * znana ścieżka rdzenia ich nie czyta.
 */
const PRZEDROSTKI_BEZ_KONSUMENTA: readonly string[] = ['harness.'];

/**
 * Zdanie doklejane do objaśnienia klucza konfiguracji, którego samo wykonanie zadania nie czyta wprost.
 */
const NIE_STERUJE_JESZCZE =
  'Stan: wartość zapisuje się w oknie, ale budowa wywołania modelu jeszcze jej nie czyta — w tym wydaniu nie zmienia sposobu wykonania.';

/**
 * Zdanie o stanie klucza konfiguracji albo `undefined`, gdy klucz dociera do
 * wykonania. Okno konfiguracji dokleja je do objaśnienia [?] pozycji katalogu.
 */
export function adnotacjaKlucza(klucz: string): string | undefined {
  if (KLUCZE_BEZ_KONSUMENTA.includes(klucz)) return NIE_STERUJE_JESZCZE;
  if (PRZEDROSTKI_BEZ_KONSUMENTA.some((przedrostek) => klucz.startsWith(przedrostek))) {
    return 'Stan: klucz jest w katalogu i zapisuje się poprawnie, ale żadna zmierzona ścieżka rdzenia go dziś nie czyta.';
  }
  return undefined;
}
