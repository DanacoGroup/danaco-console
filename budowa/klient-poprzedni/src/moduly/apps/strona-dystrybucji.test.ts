import { describe, expect, it } from 'vitest';

import {
  ChangeKind,
  Command,
  ExtensionKind,
  ExtensionOrigin,
  type Extension,
} from '../../../../shared/contract';
import type { Kanal } from '../../protokol/kanal';
import { BEZ_ZAWEZENIA, zawez, zdanieZawezenia } from './filtr-katalogu';
import { opisStanu } from './karta-rozszerzenia';
import { utworzOknoAppCatalog } from './okno-app-catalog';
import { utworzOknoInstalledAppsManager } from './okno-installed-apps-manager';
import { utworzOknoIntegrationsHub } from './okno-integrations-hub';
import { utworzStanRozszerzen } from './stan-rozszerzen';

/**
 * Sprawdziany strony dystrybucji i konsumpcji modułu Apps.
 *
 * Pilnowane jest zachowanie, nie kształt plików — pięć rzeczy, które opracowanie
 * i zasada zero blokad stawiają wprost, a które łatwo zepsuć po cichu:
 *
 *   1. zawężenie katalogu odsiewa po polach, które pozycja NIESIE, i mówi
 *      o zakresie szukania zamiast pozwalać czytać pustkę jako brak pozycji;
 *   2. odczyt katalogu ZASTĘPUJE zbiór, a nie dokłada do niego — inaczej
 *      w wykazie zostawałaby pozycja, której rdzeń już nie zna;
 *   3. zdarzenie `extension.changed` zmienia zbiór niezależnie od tego, gdzie
 *      zaszła zmiana, a usunięcie zdejmuje też wskazanie panelu bocznego;
 *   4. ani jedna kontrolka okien nie ma atrybutu `disabled` — zasada zero
 *      blokad platformy; kontrolka bez pokrycia ma być klikalna i nazywać brak;
 *   5. instalacja niesie pochodzenie pozycji, bo to jedyne pole rozstrzygające
 *      stan wyjściowy rejestracji.
 *
 * Rdzeń jest atrapą: sprawdzian pyta o zachowanie okna, nie serwera.
 */
interface Zapis {
  komenda: string;
  zadanie: Record<string, unknown>;
}

/** Pozycje przykładowe z domeny produktu — katalog rozszerzeń platformy. */
const SERWER_REPOZYTORIOW: Extension = {
  id: 'ext-1',
  code: 'mcp-repozytoria',
  name: 'Serwer MCP Repozytoria',
  kind: ExtensionKind.Mcp,
  description: 'Dostęp do repozytoriów kodu przez most MCP',
  version: '0.9.0',
  installed: true,
  enabled: false,
  accessPointId: 'punkt-1',
  origin: ExtensionOrigin.Personal,
  updatedAt: 1,
};

const KONEKTOR_POCZTY: Extension = {
  id: 'ext-2',
  code: 'api-poczta',
  name: 'Konektor Poczta',
  kind: ExtensionKind.Api,
  description: 'Integracja skrzynki pocztowej organizacji',
  version: '2.1.0',
  installed: true,
  enabled: true,
  origin: ExtensionOrigin.Danaco,
  updatedAt: 2,
};

const UMIEJETNOSC_FAKTUR: Extension = {
  id: 'ext-3',
  code: 'skill-faktury',
  name: 'Umiejętność Faktury',
  kind: ExtensionKind.Skill,
  installed: false,
  enabled: false,
  origin: ExtensionOrigin.Danaco,
  updatedAt: 3,
};

const KATALOG = [SERWER_REPOZYTORIOW, KONEKTOR_POCZTY, UMIEJETNOSC_FAKTUR];

function atrapaKanalu(zapisy: Zapis[], odpowiedzi: Record<string, unknown> = {}): Kanal {
  const pelne: Record<string, unknown> = {
    [Command.ExtensionList]: { extensions: KATALOG },
    [Command.AccessPointList]: { points: [] },
    ...odpowiedzi,
  };
  return {
    wyslij(
      komenda: Command,
      zadanie: Record<string, unknown>,
      przyWyniku?: (wynik: unknown) => void,
    ): string {
      zapisy.push({ komenda: String(komenda), zadanie });
      const wynik = pelne[String(komenda)];
      przyWyniku?.(
        wynik === undefined
          ? { udany: false, blad: { code: 'not_found', message: 'atrapa nie zna komendy' } }
          : { udany: true, wynik },
      );
      return 'zadanie-1';
    },
    naZdarzenie: () => () => undefined,
    naDowolny: () => () => undefined,
    sesja: () => ({ id: () => 'sesja-1' }) as never,
    // Droga fail-open modułu czyta dziennik nierozpoznanych, więc atrapa musi
    // go mieć — inaczej sprawdzian padłby na warstwie, której nie bada.
    dziennikNieznanych: () => ({ naWpis: () => () => undefined }) as never,
  } as unknown as Kanal;
}

describe('zawężenie katalogu rozszerzeń', () => {
  it('odsiewa po rodzaju, pochodzeniu i stanie', () => {
    expect(zawez(KATALOG, { ...BEZ_ZAWEZENIA, rodzaj: ExtensionKind.Mcp })).toEqual([
      SERWER_REPOZYTORIOW,
    ]);
    expect(zawez(KATALOG, { ...BEZ_ZAWEZENIA, pochodzenie: ExtensionOrigin.Danaco })).toEqual([
      KONEKTOR_POCZTY,
      UMIEJETNOSC_FAKTUR,
    ]);
    expect(zawez(KATALOG, { ...BEZ_ZAWEZENIA, stan: 'wlaczone' })).toEqual([KONEKTOR_POCZTY]);
    expect(zawez(KATALOG, { ...BEZ_ZAWEZENIA, stan: 'dostepne' })).toEqual([UMIEJETNOSC_FAKTUR]);
  });

  it('szuka w nazwie, kodzie i opisie, a nie poza nimi', () => {
    expect(zawez(KATALOG, { ...BEZ_ZAWEZENIA, tekst: 'poczta' })).toEqual([KONEKTOR_POCZTY]);
    expect(zawez(KATALOG, { ...BEZ_ZAWEZENIA, tekst: 'api-poczta' })).toEqual([KONEKTOR_POCZTY]);
    expect(zawez(KATALOG, { ...BEZ_ZAWEZENIA, tekst: 'repozytoriów' })).toEqual([
      SERWER_REPOZYTORIOW,
    ]);
    // Nazwa narzędzia udostępnianego przez pozycję nie jest polem katalogu —
    // brak trafienia jest tu stanem prawdziwym, o którym okno musi powiedzieć.
    expect(zawez(KATALOG, { ...BEZ_ZAWEZENIA, tekst: 'repo.commit' })).toEqual([]);
  });

  it('mówi o zakresie szukania, gdy zawężenie jest czynne', () => {
    const zdanie = zdanieZawezenia(1, 3, { ...BEZ_ZAWEZENIA, tekst: 'poczta' });
    expect(zdanie).toContain('Pokazano 1 z 3');
    expect(zdanie).toContain('nie obejmuje udostępnianych');
    expect(zdanieZawezenia(3, 3, { ...BEZ_ZAWEZENIA })).toBe('Katalog liczy 3 pozycji.');
  });
});

describe('stan katalogu rozszerzeń', () => {
  it('zastępuje zbiór odpowiedzią rdzenia, zamiast do niego dokładać', async () => {
    const zapisy: Zapis[] = [];
    const stan = utworzStanRozszerzen(atrapaKanalu(zapisy));
    await stan.odczytajKatalog();
    expect(stan.rozszerzenia()).toHaveLength(3);
    await stan.odczytajKatalog();
    expect(stan.rozszerzenia()).toHaveLength(3);
    expect(stan.czyKatalogCzytany()).toBe(true);
  });

  it('odróżnia pustkę potwierdzoną od nieodczytanej', async () => {
    const stan = utworzStanRozszerzen(
      atrapaKanalu([], { [Command.ExtensionList]: { extensions: [] } }),
    );
    expect(stan.czyKatalogCzytany()).toBe(false);
    await stan.odczytajKatalog();
    expect(stan.czyKatalogCzytany()).toBe(true);
    expect(stan.rozszerzenia()).toEqual([]);
    expect(stan.powodOdczytu('katalog')).toBe('');
  });

  it('nazywa powód odmowy odczytu zamiast zostawiać pustkę', async () => {
    const stan = utworzStanRozszerzen(atrapaKanalu([], { [Command.ExtensionList]: undefined }));
    await stan.odczytajKatalog();
    expect(stan.powodOdczytu('katalog')).not.toBe('');
    expect(stan.czyKatalogCzytany()).toBe(false);
  });

  it('zdejmuje wskazanie panelu bocznego wraz z usuniętą pozycją', async () => {
    const stan = utworzStanRozszerzen(atrapaKanalu([]));
    await stan.odczytajKatalog();
    stan.wybierz(SERWER_REPOZYTORIOW.id);
    expect(stan.wybrane()?.code).toBe('mcp-repozytoria');
    stan.zdejmijPozycje(SERWER_REPOZYTORIOW.id);
    expect(stan.wybrane()).toBeNull();
    expect(stan.rozszerzenia()).toHaveLength(2);
  });

  it('wchłania pozycję potwierdzoną przez rdzeń na miejsce poprzedniej', async () => {
    const stan = utworzStanRozszerzen(atrapaKanalu([]));
    await stan.odczytajKatalog();
    stan.wchlonPozycje({ ...SERWER_REPOZYTORIOW, enabled: true });
    expect(stan.rozszerzenia()).toHaveLength(3);
    expect(stan.rozszerzenia().find((p) => p.id === 'ext-1')?.enabled).toBe(true);
  });
});

describe('okna strony dystrybucji', () => {
  it('nie wygasza ani jednej kontrolki (zero blokad)', async () => {
    const stan = utworzStanRozszerzen(atrapaKanalu([]));
    await stan.odczytajKatalog();
    const okna = [
      utworzOknoAppCatalog(stan).element,
      utworzOknoInstalledAppsManager(stan).element,
      utworzOknoIntegrationsHub(stan).element,
    ];
    for (const okno of okna) {
      document.body.append(okno);
      expect(okno.querySelectorAll('[disabled]')).toHaveLength(0);
      expect(okno.querySelectorAll('button').length).toBeGreaterThan(0);
    }
  });

  it('App Catalog rysuje kartę na każdą pozycję katalogu', async () => {
    const stan = utworzStanRozszerzen(atrapaKanalu([]));
    await stan.odczytajKatalog();
    const okno = utworzOknoAppCatalog(stan);
    document.body.append(okno.element);
    expect(okno.element.querySelectorAll('.mp-karta')).toHaveLength(3);
  });

  it('Installed Apps Manager pokazuje wyłącznie pozycje zainstalowane', async () => {
    const stan = utworzStanRozszerzen(atrapaKanalu([]));
    await stan.odczytajKatalog();
    const okno = utworzOknoInstalledAppsManager(stan);
    document.body.append(okno.element);
    expect(okno.element.querySelectorAll('tbody tr')).toHaveLength(2);
  });

  it('Integrations Hub bierze serwery MCP i integracje API, bez wtyczek i umiejętności', async () => {
    const stan = utworzStanRozszerzen(atrapaKanalu([]));
    await stan.odczytajKatalog();
    const okno = utworzOknoIntegrationsHub(stan);
    document.body.append(okno.element);
    const wiersze = [...okno.element.querySelectorAll('.mp-integracje__wiersz')];
    expect(wiersze).toHaveLength(2);
    expect(wiersze.map((w) => (w as HTMLElement).dataset['integracja'])).toEqual(['ext-1', 'ext-2']);
  });

  it('instalacja z karty niesie pochodzenie pozycji, bo ono rozstrzyga stan wyjściowy', async () => {
    const zapisy: Zapis[] = [];
    const stan = utworzStanRozszerzen(
      atrapaKanalu(zapisy, {
        [Command.ExtensionInstall]: {
          extension: { ...UMIEJETNOSC_FAKTUR, installed: true, enabled: true },
        },
      }),
    );
    await stan.odczytajKatalog();
    const okno = utworzOknoAppCatalog(stan);
    document.body.append(okno.element);
    const karta = okno.element.querySelector('[data-rozszerzenie="ext-3"]');
    const przyciski = [...(karta?.querySelectorAll('button') ?? [])];
    const instaluj = przyciski.find((p) => p.textContent === 'Zainstaluj');
    expect(instaluj).toBeDefined();
    instaluj?.click();
    await Promise.resolve();
    const zadanie = zapisy.find((z) => z.komenda === String(Command.ExtensionInstall));
    expect(zadanie?.zadanie['origin']).toBe(ExtensionOrigin.Danaco);
    expect(zadanie?.zadanie['code']).toBe('skill-faktury');
  });
});

describe('opis stanu pozycji', () => {
  it('rozróżnia trzy stany, które niosą dwa pola kontraktu', () => {
    expect(opisStanu(UMIEJETNOSC_FAKTUR)).toBe('dostępna');
    expect(opisStanu(SERWER_REPOZYTORIOW)).toBe('zainstalowana, wyłączona');
    expect(opisStanu(KONEKTOR_POCZTY)).toBe('zainstalowana, włączona');
  });
});

describe('zdarzenie katalogu rozszerzeń', () => {
  it('wchłania zmianę i usunięcie niezależnie od miejsca ich powstania', async () => {
    // Słuchacz trzymany w polu obiektu, nie w zmiennej: przypisanie wewnątrz
    // atrapy jest dla kompilatora niewidoczne i zawęziłby typ zmiennej do `null`.
    const gniazdo: { sluchacz: ((tresc: unknown) => void) | null } = { sluchacz: null };
    const kanal = {
      wyslij: (_k: Command, _z: unknown, przy?: (w: unknown) => void) => {
        przy?.({ udany: true, wynik: { extensions: KATALOG } });
        return 'x';
      },
      naZdarzenie: (_typ: string, s: (tresc: unknown) => void) => {
        gniazdo.sluchacz = s;
        return () => undefined;
      },
      naDowolny: () => () => undefined,
      sesja: () => ({ id: () => 'sesja-1' }) as never,
      dziennikNieznanych: () => ({ naWpis: () => () => undefined }) as never,
    } as unknown as Kanal;

    const stan = utworzStanRozszerzen(kanal);
    await stan.odczytajKatalog();
    expect(gniazdo.sluchacz).not.toBeNull();

    gniazdo.sluchacz?.({
      change: ChangeKind.Updated,
      extension: { ...KONEKTOR_POCZTY, enabled: false },
    });
    expect(stan.rozszerzenia().find((p) => p.id === 'ext-2')?.enabled).toBe(false);

    gniazdo.sluchacz?.({ change: ChangeKind.Deleted, extension: KONEKTOR_POCZTY });
    expect(stan.rozszerzenia()).toHaveLength(2);
  });
});
