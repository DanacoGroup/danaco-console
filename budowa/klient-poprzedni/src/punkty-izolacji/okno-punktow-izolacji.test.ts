import { describe, expect, it } from 'vitest';
import { Command, ConfigScope, IsolationLayer } from '../../../shared/contract';
import type { Kanal } from '../protokol/kanal';
import { utworzOknoPunktowIzolacji } from './okno-punktow-izolacji';

/**
 * Sprawdziany okna Punktów Izolacji — czynność, nie kształt pliku.
 *
 * Pilnowane są cztery rzeczy, których opracowanie żąda wprost, a które wcześniej
 * były w tym oknie zaszyte na sztywno albo rozdzielone na zakładki:
 *
 *   1. trzy panele w kolejności selektor → macierz → profil (rozdz. 6.1 Modelu
 *      konfiguracji), a nie pasek zakładek;
 *   2. objaśnienie kontekstowe `[?]` przy każdym przełączniku macierzy
 *      (rozdz. 3.3), z treścią, nie samym znakiem;
 *   3. zasięg i warstwa niesione w żądaniach `isolation.*` zamiast wpisanych
 *      w kod;
 *   4. wskazanie poziomu w selektorze przestawiające odczyt macierzy — dowód,
 *      że lewy panel steruje, a nie tylko wygląda na selektor.
 *
 * Rdzeń jest atrapą: sprawdzian pyta o zachowanie okna, nie o zachowanie
 * serwera. Odpowiedź nieznana atrapie wraca odmową — tak jak wraca z rdzenia,
 * który komendy nie zna.
 */
interface Zapis {
  komenda: string;
  zadanie: Record<string, unknown>;
}

const ODPOWIEDZI: Record<string, unknown> = {
  [Command.IsolationTechnicalGet]: { switches: [] },
  [Command.IsolationContextGet]: { switches: [] },
  [Command.IsolationScopeList]: {
    scopes: [
      { scope: ConfigScope.Global, label: 'Globalny', order: 1 },
      { scope: ConfigScope.Session, label: 'Karta sesji', order: 6, narrowest: false },
    ],
  },
  [Command.IsolationProfileList]: { profiles: [] },
};

function atrapaKanalu(zapisy: Zapis[]): Kanal {
  return {
    wyslij(
      komenda: Command,
      zadanie: Record<string, unknown>,
      przyWyniku?: (wynik: unknown) => void,
    ): string {
      zapisy.push({ komenda: String(komenda), zadanie });
      const wynik = ODPOWIEDZI[String(komenda)];
      przyWyniku?.(
        wynik === undefined
          ? { udany: false, blad: { code: 'not_found', message: 'atrapa' } }
          : { udany: true, wynik },
      );
      return 'x';
    },
    naZdarzenie: () => () => undefined,
    naDowolny: () => () => undefined,
    sesja: () => ({ id: () => 'sesja-1' }) as never,
    dziennikNieznanych: () => ({}) as never,
  } as unknown as Kanal;
}

describe('okno punktów izolacji', () => {
  it('stoi na trzech panelach i sześciu obszarach', () => {
    const okno = utworzOknoPunktowIzolacji(atrapaKanalu([]));
    document.body.append(okno.element);
    const panele = [...okno.element.querySelectorAll('.pi-panel')].map(
      (p) => (p as HTMLElement).dataset['panel'],
    );
    expect(panele).toEqual(['zasieg', 'macierz', 'profil']);
    expect(okno.element.querySelectorAll('.pi-obszar')).toHaveLength(6);
  });

  it('każdy przełącznik macierzy niesie objaśnienie [?]', async () => {
    const okno = utworzOknoPunktowIzolacji(atrapaKanalu([]));
    document.body.append(okno.element);
    (okno.element as unknown as { showModal: () => void }).showModal = () => undefined;
    okno.otworz();
    await new Promise((gotowe) => setTimeout(gotowe, 0));
    const macierz = okno.element.querySelector('[data-panel="macierz"]');
    const dymki = macierz?.querySelectorAll('.dn-tooltip') ?? [];
    // Osiem punktów technicznych rysuje się bez odpowiedzi rdzenia (tabela
    // stała); kontekst dorysowuje trzy po udanym odczycie.
    expect(dymki.length).toBeGreaterThanOrEqual(8);
    for (const dymek of dymki) {
      expect(dymek.querySelector('button')?.getAttribute('aria-label') ?? '').not.toBe('');
    }
  });

  it('zasięg i warstwa jadą w żądaniach, nie są zaszyte', () => {
    const zapisy: Zapis[] = [];
    const okno = utworzOknoPunktowIzolacji(atrapaKanalu(zapisy));
    document.body.append(okno.element);
    (okno.element as unknown as { showModal: () => void }).showModal = () => undefined;
    okno.otworz();

    const odczyty = zapisy.filter(
      (z) => z.komenda === Command.IsolationTechnicalGet || z.komenda === Command.IsolationContextGet,
    );
    expect(odczyty.length).toBeGreaterThan(0);
    for (const odczyt of odczyty) {
      expect(odczyt.zadanie['scope']).toBe(ConfigScope.Global);
    }
    // Odczyt wartości obecnej idzie na warstwie czynnej okna — domyślnej,
    // dopóki pas narzędzi jej nie przełączy.
    expect(
      odczyty.some((o) => o.zadanie['layer'] === IsolationLayer.Default),
    ).toBe(true);
  });

  it('wskazanie poziomu w selektorze przestawia odczyt macierzy', async () => {
    const zapisy: Zapis[] = [];
    const okno = utworzOknoPunktowIzolacji(atrapaKanalu(zapisy));
    document.body.append(okno.element);
    (okno.element as unknown as { showModal: () => void }).showModal = () => undefined;
    okno.otworz();
    await new Promise((gotowe) => setTimeout(gotowe, 0));

    const wybor = okno.element.querySelector<HTMLInputElement>(
      `[data-panel="zasieg"] tr[data-poziom="${ConfigScope.Session}"] input[type="radio"]`,
    );
    expect(wybor).not.toBeNull();

    zapisy.length = 0;
    wybor?.click();
    await new Promise((gotowe) => setTimeout(gotowe, 0));

    const poZmianie = zapisy.filter((z) => z.komenda === Command.IsolationTechnicalGet);
    expect(poZmianie.length).toBeGreaterThan(0);
    for (const odczyt of poZmianie) {
      expect(odczyt.zadanie['scope']).toBe(ConfigScope.Session);
      expect(odczyt.zadanie['scopeId']).toBe('sesja-1');
    }
  });
});
