import { KluczUstawieniaOkna } from './klucze-ustawien';

/**
 * Adnotacje o tym, co ze sterowań okna dociera do wywołania modelu, a co jest
 * wyłącznie zapisywane.
 *
 * Sterowanie, którego wartość nie dociera do wywołania, zostaje czynne —
 * wyszarzenie byłoby blokadą — ale mówi o swoim stanie wprost w objaśnieniu [?].
 * Milczenie kazałoby uznać zapisaną wartość za sterującą modelem.
 *
 * Wykaz, którego nikt nie odświeża, kłamie w dymkach i każe budować drugi raz
 * to, co już działa. Dlatego przy każdym sterowaniu stoi plik i miejsce
 * w rdzeniu, po których da się zdanie sprawdzić zamiast w nie uwierzyć:
 *   · naklad_rozumowania — `core/adapter_rozmowa_wykonanie.go` (funkcja `Ustal`,
 *     pole `Naklad`) → `injection/argumenty.go`, parametr `--effort`;
 *   · kanal_modelu_zapasowy — tamże, pole `ModelZapasowy`, po przekładzie kodu
 *     kanału na identyfikator modelu (`modelKanalu`) → `--fallback-model`;
 *   · host_wykonania — `zdalne/hosty.go` (odczyt z bazy), `zdalne/tor.go`
 *     (złożenie toru SSH), `injection/uruchamiacz_okna.go` (gałąź
 *     `shared.ExecutionEnvRemote` → `zdalne.Przeloz`);
 *   · środowisko wykonania — zasięg `remote` prowadzi jawny tor SSH, a każde
 *     brakujące ogniwo jest nazwaną odmową, nie cichym startem na rdzeniu
 *     (`injection/uruchamiacz_okna.go`, `rozruchZdalny`); zasięg `local` schodzi
 *     na host rdzenia i zostawia o tym wpis w dzienniku;
 *   · model — kanał modelu jest osią całego wywołania.
 *
 * Zdanie o sterowaniu niepotwierdzonym zostaje jedno — patrz
 * `AGENT_ZAPISYWANY`. Adnotację zdejmuje się przez skreślenie wiersza z tego
 * wykazu, w jednym miejscu, wraz ze wskazaniem miejsca w rdzeniu, które to
 * uzasadnia.
 */

/**
 * Zdanie o wyborze eksperta — stan spoiny, nie obietnica.
 *
 * Pole `agentId` się utrwala: `session/okno.go` (pole `Agent`),
 * `core/przeklad.go` (funkcja `zmianaOkna`). Nałożenia tożsamości eksperta —
 * warstw promptu, skilli, wtyczek — na wywołanie modelu nie potwierdza żadna
 * znana ścieżka rdzenia, więc zdanie zostaje: milczenie kazałoby uznać, że okno
 * pracuje tożsamością eksperta.
 *
 * Sekcja „Moi agenci" nie jest przy tym wyszarzana — to byłaby blokada.
 */
const AGENT_ZAPISYWANY =
  'Wybór eksperta: okno bierze jego model bazowy i zapisuje jego kod (rdzeń utrwala pole agentId), ale nałożenia tożsamości eksperta — warstw promptu, skilli, wtyczek — na wywołanie nie potwierdzono pomiarem.';

/**
 * Zdanie o zasięgu `local`, powiedziane wprost zamiast przemilczane.
 *
 * Źródło: `injection/uruchamiacz_okna.go` — toru zwrotnego do urządzenia
 * Operatora w drzewie nie ma, więc `local` startuje na hoście rdzenia
 * i zostawia o tym wpis w dzienniku. Wybór jest honorowany dosłownie dopóty,
 * dopóki rdzeń stoi na urządzeniu Operatora.
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
 * Klucze konfiguracji bez konsumenta w rdzeniu.
 *
 * Wykaz jest pusty: host wykonania, kanał zapasowy i nakład rozumowania mają
 * konsumenta wskazanego w nagłówku pliku. Pusty wykaz zostaje, bo droga
 * dopisania ma być ta sama co droga wycofania — następny klucz bez konsumenta
 * wpisuje się tutaj, a nie w nowym mechanizmie.
 */
const KLUCZE_BEZ_KONSUMENTA: readonly string[] = [];

/**
 * Przedrostek „harness." obejmuje całą rodzinę kluczy programu kanału
 * (program_claude, plik_ustawien, konfiguracja_mcp): katalog je zna, żadna
 * znana ścieżka rdzenia ich nie czyta.
 */
const PRZEDROSTKI_BEZ_KONSUMENTA: readonly string[] = ['harness.'];

/** Zdanie doklejane do objaśnienia klucza, którego wykonanie nie czyta. */
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
