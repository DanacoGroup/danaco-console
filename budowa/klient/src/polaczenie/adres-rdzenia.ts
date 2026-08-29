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
