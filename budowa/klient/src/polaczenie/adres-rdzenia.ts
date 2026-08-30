// Adres rdzenia: ścieżka WebSocket i port nasłuchu.

// Równa `SciezkaGniazdaDomyslna` z `transport/ustawienia.go` w rdzeniu.
const SCIEZKA_GNIAZDA = '/ws';

// Równy `PortDomyslny` z `konfiguracja/ustawienia.go` w rdzeniu.
const PORT_RDZENIA_LOKALNEGO = 17870;

export function adresRdzeniaLokalnego(port: number = PORT_RDZENIA_LOKALNEGO): string {
  return `ws://127.0.0.1:${port}${SCIEZKA_GNIAZDA}`;
}

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

/**
 * Wskazanie rdzenia od powłoki desktopowej. Strona wczytana z pakietu powłoki
 * ma pochodzenie `tauri://localhost`, z którego nie da się wywieść serwera
 * wdrożenia — jedynym źródłem jest wtedy powłoka, która wskazanie prowadzi
 * (wpisane przy składaniu instalki, zmienną środowiska albo w oknie).
 * Pustka znaczy: strona stoi poza powłoką albo powłoka wskazania nie ma.
 */
export async function adresGniazdaOdPowloki(): Promise<string | null> {
  const most = globalThis as {
    __TAURI__?: { core?: { invoke?: (nazwa: string) => Promise<unknown> } };
  };
  const wywolaj = most.__TAURI__?.core?.invoke;
  if (wywolaj === undefined) return null;
  try {
    const adres = await wywolaj('adres_rdzenia');
    return typeof adres === 'string' ? adresGniazdaRdzenia(adres) : null;
  } catch {
    // Powłoka bez tej komendy nie jest błędem strony: zostaje droga zwykła.
    return null;
  }
}
