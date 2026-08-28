import { describe, expect, it } from 'vitest';

import { Command, StudioTableStructureOp } from '../../../../shared/contract';
import type { Kanal } from '../../protokol/kanal';
import { utworzStanStudio } from './stan-studio';
import { utworzTabelaPanel } from './tabela-panel';
import { utworzTabeleZrodlo } from './tabela-zrodlo';

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

/** Funkcja stanZDokumentem tworzy stan modułu studio wraz z dokumentem czynnym, ponieważ tabela nie może istnieć bez otwartego dokumentu. */
function stanZDokumentem(kanal: Kanal) {
  const stan = utworzStanStudio(kanal);
  stan.ustawOkno('okno-1');
  stan.wchlon({
    id: 'dokument-1',
    windowId: 'okno-1',
    title: 'Umowa',
    format: 'markdown',
  } as never);
  return stan;
}

/** Funkcja zdania zwraca zdania wypisane w panelu tabeli: wiersz odpowiedzi oraz opisy poszczególnych pozycji tabeli. */
function zdania(element: HTMLElement): string {
  return [...element.querySelectorAll('p')].map((akapit) => akapit.textContent ?? '').join(' ');
}

describe('warsztat tabel', () => {
  it('wskazanie rozmiaru siatką jedzie do rdzenia jako rozmiar wskazany', async () => {
    const zapisy: Zapis[] = [];
    const kanal = atrapaKanalu(zapisy, {
      [Command.StudioTableInsert]: {
        form: {},
        balance: { applied: 1, skippedCount: 0 },
        actionId: 'czynnosc-1',
        table: {
          id: 'tabela-1',
          rows: 3,
          columns: 4,
          columnWidthsMm: [40, 30, 30, 60],
          widthMm: 160,
        },
      },
      [Command.StudioTableList]: { tables: [] },
    });
    const panel = utworzTabelaPanel(stanZDokumentem(kanal), utworzTabeleZrodlo(kanal));
    document.body.append(panel.element);
    panel.przestawWidocznosc();

    const komorka = panel.element.querySelector<HTMLButtonElement>(
      '[data-wiersze="3"][data-kolumny="4"]',
    );
    expect(komorka, 'siatka daje wskazać rozmiar tabeli').not.toBeNull();
    komorka!.click();
    await new Promise((gotowe) => setTimeout(gotowe, 0));

    const wstawienie = zapisy.find((z) => z.komenda === Command.StudioTableInsert);
    expect(wstawienie?.zadanie['rows']).toBe(3);
    expect(wstawienie?.zadanie['columns']).toBe(4);
    expect(wstawienie?.zadanie['documentId']).toBe('dokument-1');

    const tresc = zdania(panel.element);
    // Bilans stoi w zdaniu także przy pełnym powodzeniu.
    expect(tresc).toContain('Zmienionych miejsc: 1');
    // Szerokości kolumn są wypisane liczbami, nie słowem „gotowe".
    expect(tresc).toContain('40, 30, 30, 60');
    // Wpis dziennika jest widoczny, bo bez niego Operator nie wie, co cofnąć.
    expect(tresc).toContain('czynnosc-1');
  });

  it('pominięcie blokadą nazywa blokadę i nie mieni scalenia udanym', async () => {
    const zapisy: Zapis[] = [];
    const kanal = atrapaKanalu(zapisy, {
      [Command.StudioTableInsert]: {
        form: {},
        balance: { applied: 1, skippedCount: 0 },
        table: { id: 'tabela-1', rows: 2, columns: 2, columnWidthsMm: [50, 50] },
      },
      [Command.StudioTableStructureEdit]: {
        form: {},
        balance: {
          applied: 0,
          skippedCount: 1,
          skipped: [
            {
              reason: 'fragment pod blokadą',
              detail: 'komórki wiersza pierwszego',
              lockName: 'podstawa prawna — nie zmieniać',
              rangeStart: 120,
              rangeEnd: 180,
            },
          ],
        },
        table: { id: 'tabela-1', rows: 2, columns: 2, columnWidthsMm: [50, 50] },
      },
      [Command.StudioTableList]: { tables: [] },
    });
    const panel = utworzTabelaPanel(stanZDokumentem(kanal), utworzTabeleZrodlo(kanal));
    document.body.append(panel.element);
    panel.przestawWidocznosc();

    // Tabelę trzeba najpierw wskazać: wstawienie robi to samo.
    panel.element.querySelector<HTMLButtonElement>('[data-wiersze="2"][data-kolumny="2"]')!.click();
    await new Promise((gotowe) => setTimeout(gotowe, 0));

    panel.element
      .querySelector<HTMLButtonElement>(`[data-budowa="${StudioTableStructureOp.MergeCells}"]`)!
      .click();
    await new Promise((gotowe) => setTimeout(gotowe, 0));

    const scalenie = zapisy.find((z) => z.komenda === Command.StudioTableStructureEdit);
    expect(scalenie?.zadanie['tableId']).toBe('tabela-1');
    expect(scalenie?.zadanie['operation']).toBe(StudioTableStructureOp.MergeCells);

    const wiersz = panel.element.querySelector<HTMLElement>('.dm-odpowiedz');
    expect(wiersz).not.toBeNull();
    const tresc = wiersz!.textContent ?? '';
    expect(tresc).toContain('Czynność nie zmieniła ani jednego miejsca');
    expect(tresc).toContain('podstawa prawna — nie zmieniać');
    expect(tresc).toContain('120–180');
    // Czynność zatrzymana blokadą NIE jest powodzeniem.
    expect(wiersz!.dataset['powodzenie']).toBe('false');
  });

  it('zamiana tekstu na tabelę bez zaznaczenia odmawia przed wyjściem do rdzenia', () => {
    const zapisy: Zapis[] = [];
    const kanal = atrapaKanalu(zapisy, { [Command.StudioTableList]: { tables: [] } });
    const panel = utworzTabelaPanel(stanZDokumentem(kanal), utworzTabeleZrodlo(kanal));
    document.body.append(panel.element);
    panel.przestawWidocznosc();

    panel.element.querySelector<HTMLButtonElement>('[data-czynnosc="tekst-na-tabele"]')!.click();

    expect(zapisy.some((z) => z.komenda === Command.StudioTableConvert)).toBe(false);
    const tresc = panel.element.querySelector<HTMLElement>('.dm-odpowiedz')?.textContent ?? '';
    expect(tresc).toContain('zaznaczonym fragmencie');
  });
});
