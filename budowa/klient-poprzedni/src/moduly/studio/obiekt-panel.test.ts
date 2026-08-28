import { describe, expect, it } from 'vitest';

import { Command, StudioObjectKind, StudioObjectSource } from '../../../../shared/contract';
import type { Kanal } from '../../protokol/kanal';
import { utworzObiektPanel } from './obiekt-panel';
import { utworzObiektZrodlo } from './obiekt-zrodlo';
import { utworzStanStudio } from './stan-studio';

/** Warsztat obiektów sprawdza cztery rzeczy: pole wskazania źródła, nazwaną odmowę wykresu, wynik nieusunięcia i brak tekstu zastępczego. */
interface Zapis {
  komenda: string;
  zadanie: Record<string, unknown>;
}

function atrapaKanalu(zapisy: Zapis[], odpowiedzi: Record<string, unknown> = {}): Kanal {
  return {
    wyslij(
      komenda: Command,
      zadanie: Record<string, unknown>,
      przyWyniku?: (wynik: unknown) => void,
    ): string {
      zapisy.push({ komenda: String(komenda), zadanie });
      const wynik = odpowiedzi[String(komenda)];
      przyWyniku?.(
        wynik === undefined
          ? { udany: false, blad: { code: 'not_found', message: 'atrapa nie zna tej komendy' } }
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

function stanZDokumentem(kanal: Kanal) {
  const stan = utworzStanStudio(kanal);
  stan.ustawOkno('okno-1');
  stan.wchlon({ id: 'dokument-1', windowId: 'okno-1', format: 'markdown' } as never);
  return stan;
}

describe('warsztat obiektów', () => {
  it('węzeł Designu jedzie polem designNodeId, nie ścieżką pliku', async () => {
    const zapisy: Zapis[] = [];
    const kanal = atrapaKanalu(zapisy, {
      [Command.StudioObjectInsert]: {
        form: {},
        balance: { applied: 1, skippedCount: 0 },
        actionId: 'czynnosc-1',
        object: {
          id: 'obiekt-1',
          kind: StudioObjectKind.Shape,
          source: StudioObjectSource.DesignModule,
          designNodeId: 'wezel-7',
          widthMm: 80,
          heightMm: 40,
        },
      },
      [Command.StudioObjectList]: { objects: [] },
    });
    const panel = utworzObiektPanel(stanZDokumentem(kanal), utworzObiektZrodlo(kanal));
    document.body.append(panel.element);
    panel.przestawWidocznosc();

    const wybory = [...panel.element.querySelectorAll<HTMLSelectElement>('select')];
    wybory[0]!.value = StudioObjectKind.Shape;
    wybory[1]!.value = StudioObjectSource.DesignModule;
    const wskazanie = panel.element.querySelector<HTMLInputElement>('input[type="text"]');
    wskazanie!.value = 'wezel-7';

    panel.element.querySelector<HTMLButtonElement>('[data-czynnosc="wstaw-obiekt"]')!.click();
    await new Promise((gotowe) => setTimeout(gotowe, 0));

    const wstawienie = zapisy.find((z) => z.komenda === Command.StudioObjectInsert);
    expect(wstawienie?.zadanie['designNodeId']).toBe('wezel-7');
    expect(wstawienie?.zadanie['path']).toBeUndefined();
    expect(wstawienie?.zadanie['kind']).toBe(StudioObjectKind.Shape);

    const tresc = panel.element.querySelector<HTMLElement>('.dm-odpowiedz')?.textContent ?? '';
    expect(tresc).toContain('Zmienionych miejsc: 1');
    // Brak tekstu zastępczego jest nazwany, a nie przemilczany.
    expect(tresc).toContain('BEZ tekstu zastępczego');
  });

  it('wykres jest odmową nazwaną, widoczną i klikalną, kierującą do Designu', () => {
    const kanal = atrapaKanalu([], { [Command.StudioObjectList]: { objects: [] } });
    const panel = utworzObiektPanel(stanZDokumentem(kanal), utworzObiektZrodlo(kanal));
    document.body.append(panel.element);

    const braki = [...panel.element.querySelectorAll<HTMLButtonElement>('[data-brak-komendy="tak"]')];
    expect(braki.length).toBeGreaterThanOrEqual(1);
    const wykres = braki.find((brak) => (brak.textContent ?? '').includes('wykres'));
    expect(wykres, 'wykres stoi w oknie jako odmowa nazwana, a nie jest ukryty').not.toBeUndefined();
    // Wygaszenie jest tu tak samo niedopuszczalne jak milczenie.
    expect(wykres!.disabled).toBe(false);
    const powod = wykres!.getAttribute('aria-description') ?? '';
    expect(powod).toContain('rdzeń');
    expect(powod).toContain('module Design');
  });

  it('odpowiedź „nie usunąłem" jest wynikiem, nie awarią, i nazywa powód', async () => {
    const zapisy: Zapis[] = [];
    const kanal = atrapaKanalu(zapisy, {
      [Command.StudioObjectList]: {
        objects: [{ id: 'obiekt-1', kind: StudioObjectKind.Image, locked: true }],
      },
      [Command.StudioObjectRemove]: {
        removed: false,
        form: {},
        balance: {
          applied: 0,
          skippedCount: 1,
          skipped: [{ reason: 'obiekt pod blokadą', lockName: 'logo kancelarii' }],
        },
      },
    });
    const panel = utworzObiektPanel(stanZDokumentem(kanal), utworzObiektZrodlo(kanal));
    document.body.append(panel.element);
    panel.przestawWidocznosc();
    await new Promise((gotowe) => setTimeout(gotowe, 0));

    panel.element.querySelector<HTMLButtonElement>('button[data-obiekt="obiekt-1"]')!.click();
    panel.element.querySelector<HTMLButtonElement>('[data-czynnosc="usun-obiekt"]')!.click();
    await new Promise((gotowe) => setTimeout(gotowe, 0));

    expect(zapisy.some((z) => z.komenda === Command.StudioObjectRemove)).toBe(true);
    const wiersz = panel.element.querySelector<HTMLElement>('.dm-odpowiedz');
    const tresc = wiersz?.textContent ?? '';
    expect(tresc).toContain('NIE usunął');
    expect(tresc).toContain('logo kancelarii');
    expect(wiersz?.dataset['powodzenie']).toBe('false');
    // Obiekt zostaje w wykazie, bo nadal jest.
    expect(panel.element.querySelector('li[data-obiekt="obiekt-1"]')).not.toBeNull();
  });
});
