/**
 * Umiejscowienie rdzenia: ścieżka gniazda i port nasłuchu.
 *
 * Warstwa połączenia zna jedno: gdzie nasłuchuje rdzeń. Wyboru adresu nie
 * dokonuje — wskazuje go wołający, bo to on wie, czy rdzeń stoi na tej samej
 * maszynie, czy pod adresem podanym przez powłokę.
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
 *
 * Nasłuch bez wskazania Operatora wiąże się z pętlą zwrotną
 * (`adresDomyslny` w `server/internal/transport/ustawienia.go`), więc to jest
 * adres, pod którym rdzeń stoi, dopóki nikt nie wskazał inaczej.
 */
export function adresRdzeniaLokalnego(port: number = PORT_RDZENIA_LOKALNEGO): string {
  return `ws://127.0.0.1:${port}${SCIEZKA_GNIAZDA}`;
}

/**
 * Przekłada adres HTTP rdzenia na adres gniazda: `http://host:port` →
 * `ws://host:port/ws`. Zwraca `null`, gdy adres nie jest adresem HTTP —
 * wołający rozstrzyga wtedy, czy sięgnąć po pętlę zwrotną, czy odmówić.
 *
 * Ścieżka gniazda jest własnością tej warstwy, więc przekład mieszka tutaj,
 * a nie u tego, kto adres HTTP zdobył.
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
