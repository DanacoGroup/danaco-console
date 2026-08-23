/**
 * Ustalenie adresu WebSocket rdzenia.
 *
 * Umiejscowienie rdzenia jest zmienne: lokalnie w trakcie budowy, docelowo
 * `danaco-system`. Klient nie zakłada, że działa tam, gdzie rdzeń,
 * dlatego adres pochodzi z konfiguracji budowania, a wartość domyślna jest
 * wyliczana ze strony serwującej interfejs (brak ustawienia =
 * wartość domyślna).
 */
const SCIEZKA_GNIAZDA = '/ws';

/**
 * Port nasłuchu rdzenia lokalnego. Musi być równy `PortDomyslny`
 * z `server/internal/konfiguracja/ustawienia.go` — wartość wybrana spoza
 * portów zajętych na `danaco-system` (8090, 8091, 8093, 8095, 8760, 8766,
 * 8771, 8772, 8790, 9980, 18765). Rozjazd z portem rdzenia znaczy, że klient
 * w rdzeń nie trafia.
 */
const PORT_RDZENIA_LOKALNEGO = '17870';

/**
 * Adres podany przez powłokę natywną, już przełożony na gniazdo.
 *
 * Adres domyślny wylicza się z lokalizacji dokumentu, a ta kłamie w jednym
 * przypadku: gdy okno powłoki dostało interfejs z pakietu osadzonego. Strona ma
 * wtedy pochodzenie `tauri.localhost`, więc wyliczenie dawałoby
 * `ws://tauri.localhost:17870/ws` — adres, pod którym nie nasłuchuje nikt.
 * Ponawianie by nie pomogło, bo milczy nie rdzeń, tylko adres. Powłoka zna port
 * rdzenia ze zmiennej `DANACO_PORT` i wystawia go poleceniem `adres_rdzenia`;
 * ta zmienna trzyma jego odpowiedź.
 */
let adresZPowloki: string | null = null;

export function adresRdzenia(): string {
  const zKonfiguracji = import.meta.env.VITE_ADRES_RDZENIA;
  if (typeof zKonfiguracji === 'string' && zKonfiguracji.length > 0) {
    return zKonfiguracji;
  }
  if (adresZPowloki !== null) {
    return adresZPowloki;
  }
  return adresDomyslny();
}

/**
 * Przyjmuje adres HTTP rdzenia podany przez powłokę natywną.
 *
 * Wywołuje się raz, przed złożeniem aplikacji (`main.ts`) — po to, żeby
 * `adresRdzenia()` pozostało wywołaniem natychmiastowym dla wszystkich swoich
 * wywołujących (złożenie łączności i podglądy modułów), zamiast rozlewać
 * obietnicę po całym drzewie kompozycji.
 *
 * Pierwszeństwo pozostaje przy wskazaniu jawnym z konfiguracji budowania
 * (`VITE_ADRES_RDZENIA`) — powłoka dopowiada wartość domyślną
 * lepszą od wyliczonej, a nie nadpisuje wskazania Operatora.
 *
 * Zwraca informację, czy adres dało się przyjąć. Adres nie do rozłożenia nie
 * wstrzymuje niczego: zostaje przy wyliczeniu z lokalizacji.
 */
export function przyjmijAdresZPowloki(adresHttp: string): boolean {
  const gniazdo = naAdresGniazda(adresHttp);
  adresZPowloki = gniazdo;
  return gniazdo !== null;
}

/**
 * Przekłada adres HTTP rdzenia na adres gniazda: `http://host:port` →
 * `ws://host:port/ws`. Ścieżka gniazda i tak jest własnością tego pliku, więc
 * przekład mieszka tutaj, a nie w moście do powłoki.
 */
function naAdresGniazda(adresHttp: string): string | null {
  try {
    const adres = new URL(adresHttp);
    if (adres.protocol !== 'http:' && adres.protocol !== 'https:') return null;
    const schemat = adres.protocol === 'https:' ? 'wss:' : 'ws:';
    return `${schemat}//${adres.host}${SCIEZKA_GNIAZDA}`;
  } catch {
    return null;
  }
}

/**
 * Adres wyliczony z bieżącej lokalizacji dokumentu.
 *
 * Port bierze się ze strony, nie ze stałej: rdzeń serwuje pakiet interfejsu sam
 * — ta sama nasłuchująca końcówka oddaje `client/dist` i gniazdo `/ws`.
 * Doklejenie zaszytego portu do hosta ze strony kazałoby interfejsowi otwartemu
 * spod rdzenia na porcie innym niż domyślny szukać gniazda tam, gdzie nikt nie
 * słucha. Stała zostaje wyłącznie jako wartość ostatniej szansy.
 *
 * Rozstrzygnięcie idzie po pochodzeniu strony:
 *   http/https z portem  → ten sam host i ten sam port — rdzeń serwuje stronę
 *   http/https bez portu → host ze strony, port domyślny rdzenia (80/443 nie
 *                          jest portem rdzenia, więc byłby ślepy)
 *   file: i pozostałe    → pętla zwrotna i port domyślny; tak wygląda pakiet
 *                          osadzony w powłoce, gdzie strona nie ma źródła sieciowego
 */
function adresDomyslny(): string {
  const lokalizacja = globalThis.location;
  const sieciowa =
    lokalizacja !== undefined &&
    (lokalizacja.protocol === 'http:' || lokalizacja.protocol === 'https:');
  if (!sieciowa) {
    return `ws://127.0.0.1:${PORT_RDZENIA_LOKALNEGO}${SCIEZKA_GNIAZDA}`;
  }
  const schemat = lokalizacja.protocol === 'https:' ? 'wss:' : 'ws:';
  const host = lokalizacja.hostname || '127.0.0.1';
  const port = lokalizacja.port || PORT_RDZENIA_LOKALNEGO;
  return `${schemat}//${host}:${port}${SCIEZKA_GNIAZDA}`;
}
