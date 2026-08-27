/**
 * Gniazdo zastępcze zastępuje WebSocket środowiska w sprawdzianach, odwzorowując
 * kolejkowanie, ponawianie oraz kolejność stanów klienta bez udziału rdzenia
 * rzeczywistego.
 */
export class GniazdoZastepcze {
  static readonly OPEN = 1;
  static otwarte: GniazdoZastepcze[] = [];

  readonly adres: string;
  readyState = 0;
  wyslane: string[] = [];
  private sluchacze = new Map<string, ((zdarzenie: unknown) => void)[]>();

  constructor(adres: string) {
    this.adres = adres;
    GniazdoZastepcze.otwarte.push(this);
  }

  addEventListener(nazwa: string, sluchacz: (zdarzenie: unknown) => void): void {
    const wpisani = this.sluchacze.get(nazwa) ?? [];
    wpisani.push(sluchacz);
    this.sluchacze.set(nazwa, wpisani);
  }

  send(ramka: string): void {
    this.wyslane.push(ramka);
  }

  close(): void {
    this.readyState = 3;
    this.oglos('close', {});
  }

  /** Doprowadza gniazdo do stanu otwartego, tak jak zrobiłoby to środowisko. */
  otworz(): void {
    this.readyState = GniazdoZastepcze.OPEN;
    this.oglos('open', {});
  }

  /** Podaje ramkę przychodzącą. */
  przyjmij(dane: unknown): void {
    this.oglos('message', { data: dane });
  }

  private oglos(nazwa: string, zdarzenie: unknown): void {
    for (const sluchacz of this.sluchacze.get(nazwa) ?? []) sluchacz(zdarzenie);
  }
}

/**
 * Podstawia gniazdo zastępcze pod nazwę środowiska i oddaje przywrócenie stanu
 * pierwotnego. Wykaz gniazd założonych zaczyna się przy każdym podstawieniu
 * od nowa.
 */
export function podstawGniazdo(): () => void {
  const srodowisko = globalThis as unknown as Record<string, unknown>;
  const pierwotne = srodowisko['WebSocket'];
  GniazdoZastepcze.otwarte = [];
  srodowisko['WebSocket'] = GniazdoZastepcze;
  return () => {
    srodowisko['WebSocket'] = pierwotne;
  };
}
