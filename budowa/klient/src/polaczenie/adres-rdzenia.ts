// Adres rdzenia: ścieżka WebSocket i port nasłuchu.

// Równa `SciezkaGniazdaDomyslna` z `transport/ustawienia.go` w rdzeniu.
const SCIEZKA_GNIAZDA = '/ws';

// Równy `PortDomyslny` z `konfiguracja/ustawienia.go` w rdzeniu.
const PORT_RDZENIA_LOKALNEGO = 17870;

// Równy `ParametrSekretu` z `transport/ustawienia.go` w rdzeniu.
const PARAMETR_SEKRETU = 'sekret';

export function adresRdzeniaLokalnego(port: number = PORT_RDZENIA_LOKALNEGO): string {
  return `ws://127.0.0.1:${port}${SCIEZKA_GNIAZDA}`;
}

/* Pochodzenia własne powłoki desktopowej. Strona z pakietu Tauri stoi pod
   `tauri://localhost` albo `http://tauri.localhost`, a pod tą nazwą nie
   nasłuchuje żaden rdzeń — adres złożony z takiego pochodzenia wskazywałby na
   samą powłokę. Wskazanie rdzenia powłoka podaje osobno. */
const HOSTY_POWLOKI: readonly string[] = ['tauri.localhost'];
const SCHEMATY_POWLOKI: readonly string[] = ['tauri:'];

/** Czy pochodzenie dokumentu jest pochodzeniem własnym powłoki desktopowej. */
export function czyPochodzeniePowloki(adresHttp: string): boolean {
  try {
    const adres = new URL(adresHttp);
    return SCHEMATY_POWLOKI.includes(adres.protocol) || HOSTY_POWLOKI.includes(adres.hostname);
  } catch {
    return false;
  }
}

export function adresGniazdaRdzenia(adresHttp: string): string | null {
  if (czyPochodzeniePowloki(adresHttp)) return null;
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
 * Wskazanie rdzenia podane przez powłokę desktopową przed wczytaniem strony.
 * Strona z pakietu powłoki ma pochodzenie `tauri://localhost`, z którego nie da
 * się wywieść serwera wdrożenia — powłoka wpisuje je więc skryptem wstępnym.
 * Pustka znaczy: strona stoi poza powłoką albo powłoka wskazania nie ma.
 */
export function adresGniazdaOdPowloki(): string | null {
  const most = globalThis as { DanacoAdresRdzenia?: string };
  const adres = most.DanacoAdresRdzenia ?? '';
  return adres === '' ? null : adresGniazdaRdzenia(adres);
}

/**
 * Sekret nawiązania wpisany przez powłokę skryptem wstępnym obok wskazania
 * rdzenia. Pustka znaczy stronę poza powłoką albo rdzeń bez sekretu.
 */
export function sekretNawiazaniaOdPowloki(): string | null {
  const most = globalThis as { DanacoSekretNawiazania?: unknown };
  const sekret = most.DanacoSekretNawiazania;
  return typeof sekret === 'string' && sekret !== '' ? sekret : null;
}

/**
 * Adres gniazda z sekretem nawiązania w parametrze zapytania. Rdzeń porównuje
 * go przed uaktualnieniem gniazda; bez sekretu adres wraca bez parametru.
 */
export function adresNawiazania(adres: string, sekret: string | null): string {
  if (sekret === null || adres === '') return adres;
  const zlozony = new URL(adres);
  zlozony.searchParams.set(PARAMETR_SEKRETU, sekret);
  return zlozony.toString();
}
