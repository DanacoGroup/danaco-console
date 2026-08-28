/**
 * Ustalenie adresu WebSocket rdzenia, zmiennego umiejscowieniem: lokalnie w trakcie budowy, docelowo w środowisku produkcyjnym, z wartością domyślną wyliczaną ze strony serwującej interfejs.
 */
const SCIEZKA_GNIAZDA = '/ws';

/**
 * Port nasłuchu rdzenia lokalnego musi być równy porcie domyślnym z ustawień rdzenia; rozjazd portów znaczy, że klient w rdzeń nie trafia.
 */
const PORT_RDZENIA_LOKALNEGO = '17870';

/**
 * Adres podany przez powłokę natywną, już przełożony na gniazdo, bo adres wyliczony z lokalizacji dokumentu kłamie przy interfejsie z pakietu osadzonego.
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
 * Przyjmuje adres HTTP rdzenia podany przez powłokę natywną, wywoływane raz przed złożeniem aplikacji, z pierwszeństwem wskazania jawnego z konfiguracji budowania.
 */
export function przyjmijAdresZPowloki(adresHttp: string): boolean {
  const gniazdo = naAdresGniazda(adresHttp);
  adresZPowloki = gniazdo;
  return gniazdo !== null;
}

/**
 * Przekłada adres HTTP rdzenia na adres gniazda; przekład mieszka w tym pliku, a nie w moście do powłoki.
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
 * Adres wyliczony z bieżącej lokalizacji dokumentu: port bierze się ze strony, nie ze stałej, bo rdzeń serwuje pakiet interfejsu tą samą końcówką co gniazdo.
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
