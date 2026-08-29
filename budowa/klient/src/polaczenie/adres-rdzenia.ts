/**
 * Umiejscowienie rdzenia: ścieżka gniazda i port nasłuchu. Warstwa połączenia
 * zna jedno, gdzie nasłuchuje rdzeń — adres wskazuje wołający, zależnie od
 * tego, czy rdzeń stoi na tej maszynie, czy pod adresem podanym przez powłokę.
 */

/**
 * Ścieżka HTTP, pod którą rdzeń wystawia kanał WebSocket. Musi być równa
 * `SciezkaGniazdaDomyslna` z `server/internal/transport/ustawienia.go`; reszta
 * ścieżek nasłuchu oddaje pliki interfejsu, nie gniazdo.
 */
const SCIEZKA_GNIAZDA = '/ws';

/**
 * Port nasłuchu rdzenia lokalnego. Musi być równy `PortDomyslny`
 * z `server/internal/konfiguracja/ustawienia.go`. Rozjazd z portem rdzenia
 * znaczy, że klient w rdzeń nie trafia.
 */
const PORT_RDZENIA_LOKALNEGO = 17870;

/**
 * Adres gniazda rdzenia nasłuchującego na pętli zwrotnej tej samej maszyny.
 * Nasłuch bez wskazania adresu wiąże się z pętlą zwrotną — to adres, pod
 * którym rdzeń stoi, dopóki nikt nie wskazał inaczej.
 */
export function adresRdzeniaLokalnego(port: number = PORT_RDZENIA_LOKALNEGO): string {
  return `ws://127.0.0.1:${port}${SCIEZKA_GNIAZDA}`;
}

/**
 * Przekłada adres HTTP rdzenia na adres gniazda, zamieniając schemat i
 * ścieżkę. Zwraca `null`, gdy adres nie jest adresem HTTP. Przekład mieszka
 * w tej warstwie, bo ścieżka gniazda jest jej własnością.
 */
export function adresGniazdaRdzenia(adresHttp: string): string | null {
  try {
    const adres = new URL(adresHttp);
    if (adres.protocol !== 'http:' && adres.protocol !== 'https:') return null;
    const schemat = adres.protocol === 'https:' ? 'wss:' : 'ws:';
    return `${schemat}//${adres.host}${SCIEZKA_GNIAZDA}`;
  } catch {
    return null;
  }
}
