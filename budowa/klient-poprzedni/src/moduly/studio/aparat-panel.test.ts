import { describe, expect, it } from 'vitest';

import { Command, StudioApparatusKind, StudioFieldKind } from '../../../../shared/contract';
import { utworzAparatPanel } from './aparat-panel';
import { utworzAparatZrodlo } from './aparat-zrodlo';
import type { Kanal } from '../../protokol/kanal';
import { utworzStanStudio } from './stan-studio';

/** Testy panelu aparatu dokumentu i pola sprawdzają, że nieświeżość elementów spisu treści jest widoczna wprost, a nie ukryta jako pustka. */
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

describe('aparat dokumentu', () => {
  it('spis treści zakłada się z poziomów nagłówków, a numer nadaje rdzeń', async () => {
    const zapisy: Zapis[] = [];
    const kanal = atrapaKanalu(zapisy, {
      [Command.StudioApparatusInsert]: {
        form: {},
        balance: { applied: 1, skippedCount: 0 },
        actionId: 'czynnosc-1',
        item: {
          id: 'element-1',
          kind: StudioApparatusKind.Toc,
          levelsFrom: 1,
          levelsTo: 3,
          number: 'I',
          entries: [{}, {}, {}],
        },
      },
      [Command.StudioApparatusList]: { items: [] },
      [Command.StudioFieldList]: { fields: [] },
    });
    const panel = utworzAparatPanel(stanZDokumentem(kanal), utworzAparatZrodlo(kanal));
    document.body.append(panel.element);
    panel.przestawWidocznosc();

    const poziomy = [...panel.element.querySelectorAll<HTMLInputElement>('input[type="number"]')];
    poziomy[0]!.value = '1';
    poziomy[1]!.value = '3';
    panel.element.querySelector<HTMLButtonElement>('[data-czynnosc="zaloz-aparat"]')!.click();
    await new Promise((gotowe) => setTimeout(gotowe, 0));

    const zalozenie = zapisy.find((z) => z.komenda === Command.StudioApparatusInsert);
    expect(zalozenie?.zadanie['kind']).toBe(StudioApparatusKind.Toc);
    expect(zalozenie?.zadanie['levelsFrom']).toBe(1);
    expect(zalozenie?.zadanie['levelsTo']).toBe(3);

    const tresc = panel.element.querySelector<HTMLElement>('.dm-odpowiedz')?.textContent ?? '';
    expect(tresc).toContain('spis treści');
    expect(tresc).toContain('pozycji zebranych 3');
  });

  it('element nieświeży jest liczony i oznaczony, a nie pokazany jako zgodny', async () => {
    const kanal = atrapaKanalu([], {
      [Command.StudioApparatusList]: {
        items: [
          { id: 'element-1', kind: StudioApparatusKind.Toc, stale: true, levelsFrom: 1 },
          { id: 'element-2', kind: StudioApparatusKind.Footnote, number: '1', text: 'zob. wyżej' },
        ],
      },
      [Command.StudioFieldList]: {
        fields: [{ id: 'pole-1', kind: StudioFieldKind.PageCount, stale: true }],
      },
    });
    const panel = utworzAparatPanel(stanZDokumentem(kanal), utworzAparatZrodlo(kanal));
    document.body.append(panel.element);
    panel.przestawWidocznosc();
    await new Promise((gotowe) => setTimeout(gotowe, 0));

    const swiezosc = panel.element.querySelector<HTMLElement>('.ms-aparat__swiezosc');
    expect(swiezosc?.dataset['nieswieze']).toBe('2');
    expect(swiezosc?.textContent ?? '').toContain('WYMAGA ODŚWIEŻENIA');
    expect(
      panel.element.querySelector<HTMLElement>('li[data-element="element-1"]')?.dataset['nieswiezy'],
    ).toBe('true');
    // Licznik stoi też na przycisku, żeby był widoczny przed otwarciem nakładki.
    const wyzwalacz = panel.element.querySelector<HTMLButtonElement>(
      '[data-czynnosc="warsztat-aparatu"]',
    );
    expect(wyzwalacz?.textContent ?? '').toContain('2 do odświeżenia');
    // Pole bez policzonej wartości nazywa to wprost.
    expect(panel.element.querySelector('li[data-pole="pole-1"]')?.textContent ?? '').toContain(
      'wartość NIEPOLICZONA',
    );
  });

  it('usunięcie odmówione nie zdejmuje elementu z wykazu', async () => {
    const kanal = atrapaKanalu([], {
      [Command.StudioApparatusList]: {
        items: [{ id: 'element-1', kind: StudioApparatusKind.Bookmark }],
      },
      [Command.StudioFieldList]: { fields: [] },
      [Command.StudioApparatusRemove]: {
        removed: false,
        form: {},
        balance: {
          applied: 0,
          skippedCount: 1,
          skipped: [{ reason: 'zakładka pod blokadą', lockName: 'numer umowy' }],
        },
      },
    });
    const panel = utworzAparatPanel(stanZDokumentem(kanal), utworzAparatZrodlo(kanal));
    document.body.append(panel.element);
    panel.przestawWidocznosc();
    await new Promise((gotowe) => setTimeout(gotowe, 0));

    panel.element.querySelector<HTMLButtonElement>('button[data-element="element-1"]')!.click();
    panel.element.querySelector<HTMLButtonElement>('[data-czynnosc="usun-element"]')!.click();
    await new Promise((gotowe) => setTimeout(gotowe, 0));

    const wiersz = panel.element.querySelector<HTMLElement>('.dm-odpowiedz');
    expect(wiersz?.textContent ?? '').toContain('NIE usunął');
    expect(wiersz?.textContent ?? '').toContain('numer umowy');
    expect(wiersz?.dataset['powodzenie']).toBe('false');
    expect(panel.element.querySelector('li[data-element="element-1"]')).not.toBeNull();
  });
});
