import { describe, expect, it } from 'vitest';

import { Command, StudioTemplateFieldKind } from '../../../../shared/contract';
import type { Kanal } from '../../protokol/kanal';
import { utworzStanStudio } from './stan-studio';
import { utworzSzablonPanel } from './szablon-panel';
import { utworzSzablonZrodlo } from './szablon-zrodlo';

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

/** Funkcja polePoEtykiecie odnajduje pole tekstowe panelu po etykiecie widocznej w interfejsie, a nie po kolejności elementów w drzewie. */
function polePoEtykiecie(element: HTMLElement, fragment: string): HTMLInputElement {
  const etykiety = [...element.querySelectorAll('label')];
  const etykieta = etykiety.find((pozycja) => (pozycja.textContent ?? '').includes(fragment));
  if (etykieta === undefined) throw new Error(`brak pola o etykiecie ${fragment}`);
  const kontrolka = element.querySelector<HTMLInputElement>(`#${etykieta.htmlFor}`);
  if (kontrolka === null) throw new Error(`etykieta ${fragment} nie wiąże kontrolki`);
  return kontrolka;
}

describe('warsztat szablonów', () => {
  it('zapis szablonu niesie blokady wzorcowe i nazywa, co szablon naprawdę ma', async () => {
    const zapisy: Zapis[] = [];
    const kanal = atrapaKanalu(zapisy, {
      [Command.StudioTemplateSave]: {
        template: {
          id: 'szablon-1',
          name: 'pismo procesowe',
          format: 'markdown',
          builtin: false,
          fields: [
            {
              name: 'adresat',
              label: 'Adresat',
              kind: StudioTemplateFieldKind.Text,
              required: true,
            },
          ],
          locks: [{ id: 'blokada-1' }],
          form: {},
        },
      },
      [Command.StudioTemplateFieldList]: { fields: [] },
    });
    const panel = utworzSzablonPanel(stanZDokumentem(kanal), utworzSzablonZrodlo(kanal));
    document.body.append(panel.element);
    panel.przestawWidocznosc();

    polePoEtykiecie(panel.element, 'Nazwa szablonu').value = 'pismo procesowe';
    panel.element.querySelector<HTMLButtonElement>('[data-czynnosc="zapisz-szablon"]')!.click();
    await new Promise((gotowe) => setTimeout(gotowe, 0));

    const zapis = zapisy.find((z) => z.komenda === Command.StudioTemplateSave);
    expect(zapis?.zadanie['documentId']).toBe('dokument-1');
    expect(zapis?.zadanie['name']).toBe('pismo procesowe');
    // Blokady wzorcowe idą do rdzenia JAWNIE, a nie liczą na wartość domyślną.
    expect(zapis?.zadanie['includeLocks']).toBe(true);

    const tresc = panel.element.textContent ?? '';
    expect(tresc).toContain('szablon własny Operatora');
    expect(tresc).toContain('blokad wzorcowych 1');
    expect(tresc).toContain('niesie postać wzorcową');
  });

  it('szablonu fabrycznego rdzeń nie usuwa i panel nazywa powód', async () => {
    const zapisy: Zapis[] = [];
    const kanal = atrapaKanalu(zapisy, {
      [Command.StudioTemplateFieldList]: { fields: [] },
      [Command.StudioTemplateDelete]: { deleted: false },
    });
    const panel = utworzSzablonPanel(stanZDokumentem(kanal), utworzSzablonZrodlo(kanal));
    document.body.append(panel.element);
    panel.wskaz('szablon-fabryczny');
    await new Promise((gotowe) => setTimeout(gotowe, 0));

    panel.element.querySelector<HTMLButtonElement>('[data-czynnosc="usun-szablon"]')!.click();
    await new Promise((gotowe) => setTimeout(gotowe, 0));

    const wiersz = panel.element.querySelector<HTMLElement>('.dm-odpowiedz');
    expect(wiersz?.textContent ?? '').toContain('NIE usunął');
    expect(wiersz?.textContent ?? '').toContain('fabrycznego się nie usuwa');
    expect(wiersz?.dataset['powodzenie']).toBe('false');
  });

  it('pominięte pole wymagane jest nazwane, a wypełnienie nie mieni się udanym', async () => {
    const zapisy: Zapis[] = [];
    const kanal = atrapaKanalu(zapisy, {
      [Command.StudioTemplateFieldList]: {
        fields: [
          {
            name: 'adresat',
            label: 'Adresat',
            kind: StudioTemplateFieldKind.Text,
            required: true,
          },
        ],
      },
      [Command.StudioTemplateFill]: {
        document: { id: 'dokument-2', windowId: 'okno-1', format: 'markdown' },
        form: {},
        balance: { applied: 1, skippedCount: 0 },
        missingRequired: ['adresat'],
      },
    });
    const panel = utworzSzablonPanel(stanZDokumentem(kanal), utworzSzablonZrodlo(kanal));
    document.body.append(panel.element);
    panel.wskaz('szablon-1');
    await new Promise((gotowe) => setTimeout(gotowe, 0));

    // Pole wypełnienia jest widoczne i puste — Operator go nie podał.
    const wartosc = panel.element.querySelector<HTMLInputElement>('[data-pole="adresat"]');
    expect(wartosc, 'panel stawia pole na każde pole szablonu').not.toBeNull();

    panel.element.querySelector<HTMLButtonElement>('[data-czynnosc="wypelnij-szablon"]')!.click();
    await new Promise((gotowe) => setTimeout(gotowe, 0));

    const wypelnienie = zapisy.find((z) => z.komenda === Command.StudioTemplateFill);
    expect(wypelnienie?.zadanie['templateId']).toBe('szablon-1');
    expect(wypelnienie?.zadanie['values']).toEqual({ adresat: '' });

    const wiersz = panel.element.querySelector<HTMLElement>('.dm-odpowiedz');
    expect(wiersz?.textContent ?? '').toContain('PÓL WYMAGANYCH NIEPODANYCH: 1');
    expect(wiersz?.textContent ?? '').toContain('adresat');
    expect(wiersz?.dataset['powodzenie']).toBe('false');
  });
});
