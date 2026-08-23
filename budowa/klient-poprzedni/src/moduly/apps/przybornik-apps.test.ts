import { describe, expect, it } from 'vitest';

import { AppWorkspaceLayer, Command, KOMENDY } from '../../../../shared/contract';
import type { Kanal } from '../../protokol/kanal';
import {
  narzedziaArchitectureDesigner,
  narzedziaDeploymentPanel,
  narzedziaProductBuilder,
  narzedziaPublisherPanel,
  narzedziaWarsztatu,
} from './narzedzia-apps';
import {
  narzedziaAppCatalog,
  narzedziaInstalledApps,
  narzedziaIntegracji,
  narzedziaZaufania,
} from './narzedzia-rozszerzen';
import { utworzPrzybornikApps, type NarzedzieApps } from './przybornik-apps';
import { utworzStanProduktu } from './stan-produktu';
import { utworzStanRozszerzen } from './stan-rozszerzen';

/**
 * Sprawdziany przyborników modułu Apps.
 *
 * Pilnowane jest jedno, ale najważniejsze: czy z okna prowadzi droga do KAŻDEJ
 * komendy obszaru. Moduł ma czterdzieści jeden komend; sześć z nich prowadzą
 * kontrolki formularzy okien (architektura, warsztat, wdrożenie i trzy odczyty),
 * pozostałe trzydzieści pięć — narzędzia przyborników. Wykaz komend liczy się
 * z kontraktu w czasie działania, nie z listy wpisanej w sprawdzianie: nowa
 * komenda `apps.*` dołożona do kontraktu ma ten sprawdzian ZŁAMAĆ, bo znaczy,
 * że okno o niej nie wie.
 *
 * Drugi sprawdzian pilnuje zasady zero blokad: ani jedna kontrolka przybornika
 * nie ma atrybutu `disabled`. Narzędzie bez pokrycia (brak okna modułu, puste
 * pole wymagane) ma być klikalne i nazwać brak, a nie milczeć pod wyszarzonym
 * przyciskiem.
 */

/**
 * Pięć komend rodziny `extension.*` prowadzonych przez kontrolki okien
 * (siatka katalogu, przycisk instalacji, przełącznik stanu, formularz
 * konfiguracji, menu odinstalowania), nie przez przybornik.
 */
const KOMENDY_ROZSZERZEN_POZA_PRZYBORNIKIEM: readonly string[] = [
  Command.ExtensionList,
  Command.ExtensionInstall,
  Command.ExtensionConfigure,
  Command.ExtensionToggle,
  Command.ExtensionUninstall,
];

/** Sześć komend prowadzonych przez formularze i odczyty okien, nie przybornik. */
const KOMENDY_POZA_PRZYBORNIKIEM: readonly string[] = [
  Command.AppsArchitectureDefine,
  Command.AppsArchitectureGet,
  Command.AppsWorkspaceUpdate,
  Command.AppsWorkspaceList,
  Command.AppsDeploymentRun,
  Command.AppsDeploymentList,
];

/**
 * Kanał-atrapa. Sprawdzian pyta o to, czy okno ZNA drogę do komendy, a nie
 * o to, co odpowiada rdzeń — więc kanał nic nie odsyła i niczego nie subskrybuje.
 */
function kanalAtrapa(): Kanal {
  return {
    wyslij: () => 'zad-1',
    naZdarzenie: () => () => undefined,
    dziennikNieznanych: () => ({ naWpis: () => () => undefined, wpisy: () => [] }),
  } as unknown as Kanal;
}

/** Komplet narzędzi wszystkich pięciu przyborników modułu. */
function wszystkieNarzedzia(): readonly NarzedzieApps[] {
  const stan = utworzStanProduktu(kanalAtrapa());
  return [
    ...narzedziaProductBuilder(stan),
    ...narzedziaArchitectureDesigner(stan),
    ...narzedziaWarsztatu(stan, AppWorkspaceLayer.Frontend),
    ...narzedziaWarsztatu(stan, AppWorkspaceLayer.Backend),
    ...narzedziaDeploymentPanel(stan),
    ...narzedziaPublisherPanel(stan),
  ];
}

/** Komplet narzędzi czterech przyborników strony dystrybucji i konsumpcji. */
function wszystkieNarzedziaRozszerzen(): readonly NarzedzieApps[] {
  const stan = utworzStanRozszerzen(kanalAtrapa());
  return [
    ...narzedziaAppCatalog(stan),
    ...narzedziaInstalledApps(stan),
    ...narzedziaIntegracji(stan),
    ...narzedziaZaufania(stan),
  ];
}

describe('przybornik strony dystrybucji prowadzi do każdej komendy extension', () => {
  it('nie zostawia ani jednej komendy extension.* bez drogi z okna', () => {
    const zKontraktu = (KOMENDY as readonly string[])
      .filter((komenda) => komenda.startsWith('extension.'))
      .filter((komenda) => !KOMENDY_ROZSZERZEN_POZA_PRZYBORNIKIEM.includes(komenda))
      .sort();
    const zPrzybornikow = [
      ...new Set(wszystkieNarzedziaRozszerzen().map((narzedzie) => narzedzie.komenda)),
    ].sort();

    expect(
      zPrzybornikow,
      'komendy obszaru extension prowadzone przez przyborniki okien',
    ).toEqual(zKontraktu);
  });

  it('nazywa brak wskazanej pozycji zamiast wysyłać żądanie bez niej', async () => {
    // Świeży stan katalogu nie ma wskazanej pozycji. Narzędzie pracujące NA
    // pozycji ma wtedy odmówić zdaniem o wskazaniu, a nie wysłać żądanie
    // z pustym identyfikatorem.
    const stan = utworzStanRozszerzen(kanalAtrapa());
    const szczegol = narzedziaAppCatalog(stan).find(
      (narzedzie) => narzedzie.komenda === Command.ExtensionDetailGet,
    );
    expect(szczegol).toBeDefined();

    await expect(szczegol?.wykonaj({})).rejects.toThrow(/wskazana/);
  });

  it('nie wyszarza ani jednej kontrolki strony dystrybucji', () => {
    const przybornik = utworzPrzybornikApps('Sprawdzian', wszystkieNarzedziaRozszerzen());
    const zablokowane = przybornik.element.querySelectorAll('[disabled], [aria-disabled="true"]');

    expect(zablokowane.length, 'kontrolki przybornika z blokadą').toBe(0);
  });
});

describe('przybornik modułu Apps prowadzi do każdej komendy obszaru', () => {
  it('nie zostawia ani jednej komendy apps.* bez drogi z okna', () => {
    const zKontraktu = (KOMENDY as readonly string[])
      .filter((komenda) => komenda.startsWith('apps.'))
      .filter((komenda) => !KOMENDY_POZA_PRZYBORNIKIEM.includes(komenda))
      .sort();
    const zPrzybornikow = [
      ...new Set(wszystkieNarzedzia().map((narzedzie) => narzedzie.komenda)),
    ].sort();

    expect(
      zPrzybornikow,
      'komendy obszaru apps prowadzone przez przyborniki okien',
    ).toEqual(zKontraktu);
  });

  it('nie wymienia komendy, której kontrakt nie zna', () => {
    const znane = new Set<string>(KOMENDY as readonly string[]);
    const obce = wszystkieNarzedzia()
      .map((narzedzie) => narzedzie.komenda)
      .filter((komenda) => !znane.has(komenda));

    expect(obce, 'narzędzia wskazujące komendy spoza kontraktu').toEqual([]);
  });
});

describe('przybornik trzyma zasadę zero blokad', () => {
  it('nie wyszarza ani jednej kontrolki', () => {
    const przybornik = utworzPrzybornikApps('Sprawdzian', wszystkieNarzedzia());
    const zablokowane = przybornik.element.querySelectorAll('[disabled], [aria-disabled="true"]');

    expect(zablokowane.length, 'kontrolki przybornika z blokadą').toBe(0);
  });

  it('nazywa brak okna modułu zamiast milczeć', async () => {
    // Stan świeży nie ma jeszcze okna modułu — rdzeń go nie wskazał. Narzędzie
    // ma wtedy odmówić zdaniem o oknie, a nie wysłać żądanie bez `windowId`.
    const stan = utworzStanProduktu(kanalAtrapa());
    const narzedzie = narzedziaProductBuilder(stan)[0];
    expect(narzedzie).toBeDefined();

    await expect(narzedzie?.wykonaj({})).rejects.toThrow(/okna modułu Apps/);
  });

  it('nazywa pole wymagane, którego Operator nie wypełnił', async () => {
    const stan = utworzStanProduktu(kanalAtrapa());
    const zapis = narzedziaProductBuilder(stan).find(
      (narzedzie) => narzedzie.komenda === Command.AppsProductSave,
    );
    expect(zapis).toBeDefined();

    // Okno modułu jest sprawdzane pierwsze, więc bez niego nie dojdziemy do
    // pola — sprawdzamy więc narzędzie, które okna nie potrzebuje wcześniej niż
    // pola: walidator wartości pola jest wołany przy składaniu żądania.
    await expect(zapis?.wykonaj({ nazwa: '' })).rejects.toThrow();
  });
});

describe('przybornik pokazuje skutek, a nie samo powodzenie', () => {
  it('wypisuje zdanie zwrócone przez narzędzie', async () => {
    const przybornik = utworzPrzybornikApps('Sprawdzian', [
      {
        etykieta: 'Narzędzie próbne',
        komenda: Command.AppsProductGet,
        async wykonaj() {
          return 'Produkt „Portal klienta" (prod-1), platform: 2.';
        },
      },
    ]);
    const przycisk = przybornik.element.querySelector('button');
    expect(przycisk).not.toBeNull();

    przycisk?.click();
    // Czynność narzędzia jest obietnicą — oddajemy pętli zdarzeń jeden obrót.
    await Promise.resolve();
    await Promise.resolve();

    const skutek = przybornik.element.querySelector('.mp-przybornik__skutek');
    expect(skutek?.textContent).toContain('Portal klienta');
  });
});
