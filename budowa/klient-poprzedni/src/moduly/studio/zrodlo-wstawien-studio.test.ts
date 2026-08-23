import { describe, expect, it } from 'vitest';

import {
  Command,
  StudioExportFormat,
  StudioImportFormat,
} from '../../../../shared/contract';
import { utworzKartyDokumentow } from './karty-dokumentow';
import { utworzWydaniePanel } from './konwersja-dokumentu';
import { utworzOsadzeniePanel } from './osadzenie-panel';
import type { Kanal } from '../../protokol/kanal';
import { utworzStanStudio } from './stan-studio';
import { utworzWniesieniePliku } from './wczytanie-dokumentu';
import {
  utworzZrodloWstawienStudio,
  wstawieniaOpiszWniesienie,
  wstawieniaOpiszWydanie,
} from './zrodlo-wstawien-studio';

/**
 * Wejście, wydanie i wniesienie ze źródła — BILANS musi być widoczny.
 *
 * To jest sprawdzian pilnujący ciszy, której zlecenie zakazuje wprost:
 *
 *   1. wydanie do formatu uboższego niż dokument wypisuje wykaz cech pominiętych.
 *      Milczące zgubienie tabeli przy wydaniu do tekstu czystego jest dokładnie
 *      tym błędem, którego nie wolno popełnić;
 *   2. PDF bez warstwy tekstowej NIE udaje konwersji: bilans mówi to wprost,
 *      a okno kieruje na rozpoznanie tekstu wraz z numerem pozycji kolejki;
 *   3. fragment wniesiony z Biblioteki niesie zapis pochodzenia oddany przez
 *      RDZEŃ, nie tylko wiersz w treści — bo tamten ginie razem z kartą;
 *   4. pusty dokument zakłada się komendą rdzenia, a odmowa nie znika w ciszy;
 *   5. wydanie do PDF bez profilu wydania jest odmową NAZWANĄ przed próbą,
 *      bo profil niesie paginację i stopkę.
 */

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

describe('bilans widoczny', () => {
  it('wydanie do formatu uboższego wypisuje cechy pominięte, nie samo „zapisano"', () => {
    const zdanie = wstawieniaOpiszWydanie({
      documentId: 'dokument-1',
      format: StudioExportFormat.Txt,
      path: '/wyjscie/pismo.txt',
      bytes: 2048,
      droppedFeatures: [
        { reason: 'tabela', detail: 'trzy tabele sprowadzone do wierszy tekstu' },
        { reason: 'przypisy dolne', detail: 'dwa przypisy przeniesione na koniec' },
      ],
    });
    expect(zdanie).toContain('NIE NIESIE 2 cech dokumentu');
    expect(zdanie).toContain('tabela');
    expect(zdanie).toContain('przypisy dolne');
  });

  it('wydanie bez strat mówi to wprost, a nie milczy', () => {
    const zdanie = wstawieniaOpiszWydanie({
      documentId: 'dokument-1',
      format: StudioExportFormat.Docx,
      assetId: 'zasob-1',
    });
    // Brak strat jest zdaniem osobnym: cisza byłaby nie do odróżnienia od
    // wydania, przy którym nikt bilansu nie policzył.
    expect(zdanie).toContain('nie zgłosił ani jednej cechy');
  });

  it('PDF bez warstwy tekstowej kieruje na rozpoznanie, zamiast udawać konwersję', () => {
    const zdanie = wstawieniaOpiszWniesienie({
      format: StudioImportFormat.Pdf,
      pages: 12,
      pagesWithText: 0,
      pagesWithoutText: 12,
      tablesMissed: 3,
      imagesSkipped: 2,
      needsTextRecognition: true,
    });
    expect(zdanie).toContain('BEZ warstwy tekstowej 12');
    expect(zdanie).toContain('układów tabelarycznych NIEROZPOZNANYCH 3');
    expect(zdanie).toContain('obrazów POMINIĘTYCH 2');
    expect(zdanie).toContain('kieruje');
    expect(zdanie).toContain('rozpoznanie tekstu');
  });
});

describe('wniesienie pliku do edytora', () => {
  it('PDF ze skanów wnosi się z bilansem i z numerem pozycji kolejki rozpoznania', async () => {
    const zapisy: Zapis[] = [];
    const kanal = atrapaKanalu(zapisy, {
      [Command.StudioDocumentImportPdf]: {
        document: { id: 'dokument-2', windowId: 'okno-1', format: 'txt' },
        form: {},
        balance: {
          format: StudioImportFormat.Pdf,
          pages: 4,
          pagesWithText: 0,
          pagesWithoutText: 4,
          needsTextRecognition: true,
        },
        ingestItemId: 'pozycja-9',
      },
    });
    const stan = stanZDokumentem(kanal);
    const panel = utworzWniesieniePliku(stan, utworzZrodloWstawienStudio(kanal));
    document.body.append(panel.element);
    panel.przestawWidocznosc();

    const sciezka = panel.element.querySelector<HTMLInputElement>('input[type="text"]');
    sciezka!.value = '/dane/skan-umowy.pdf';
    panel.element.querySelector<HTMLButtonElement>('[data-czynnosc="wnies-pdf"]')!.click();
    await new Promise((gotowe) => setTimeout(gotowe, 0));

    const wniesienie = zapisy.find((z) => z.komenda === Command.StudioDocumentImportPdf);
    expect(wniesienie?.zadanie['path']).toBe('/dane/skan-umowy.pdf');
    expect(wniesienie?.zadanie['windowId']).toBe('okno-1');

    const tresc = panel.element.querySelector<HTMLElement>('.dm-odpowiedz')?.textContent ?? '';
    expect(tresc).toContain('BEZ warstwy tekstowej 4');
    expect(tresc).toContain('pozycja-9');
    expect(tresc).toContain('narzędziownię cyfryzacji');
    // Dokument wszedł do stanu modułu, więc okno pracy widzi go od razu.
    expect(stan.dokument()?.id).toBe('dokument-2');
  });
});

describe('zapis i wydanie', () => {
  it('PDF bez profilu wydania jest odmową nazwaną PRZED wyjściem do rdzenia', () => {
    const zapisy: Zapis[] = [];
    const kanal = atrapaKanalu(zapisy);
    const panel = utworzWydaniePanel(stanZDokumentem(kanal), utworzZrodloWstawienStudio(kanal));
    document.body.append(panel.element);
    panel.przestawWidocznosc();

    const format = panel.element.querySelector<HTMLSelectElement>('select');
    format!.value = StudioExportFormat.Pdf;
    panel.element.querySelector<HTMLButtonElement>('[data-czynnosc="wydaj-format"]')!.click();

    expect(zapisy.some((z) => z.komenda === Command.StudioDocumentExportFormat)).toBe(false);
    const tresc = panel.element.querySelector<HTMLElement>('.dm-odpowiedz')?.textContent ?? '';
    expect(tresc).toContain('PROFIL WYDANIA');
    expect(tresc).toContain('paginację');
  });

  it('kopia jest osobnym dokumentem i mówi, ile wersji przeniosła', async () => {
    const zapisy: Zapis[] = [];
    const kanal = atrapaKanalu(zapisy, {
      [Command.StudioDocumentCopy]: {
        document: { id: 'dokument-kopia', windowId: 'okno-1', format: 'markdown' },
        form: {},
        versionsCopied: 0,
      },
    });
    const panel = utworzWydaniePanel(stanZDokumentem(kanal), utworzZrodloWstawienStudio(kanal));
    document.body.append(panel.element);
    panel.przestawWidocznosc();

    panel.element.querySelector<HTMLButtonElement>('[data-czynnosc="kopia-dokumentu"]')!.click();
    await new Promise((gotowe) => setTimeout(gotowe, 0));

    const kopia = zapisy.find((z) => z.komenda === Command.StudioDocumentCopy);
    expect(kopia?.zadanie['documentId']).toBe('dokument-1');
    // Trzy wybory Operatora jadą jawnie, a nie liczą na wartości domyślne.
    expect(kopia?.zadanie['includeVersions']).toBe(false);
    expect(kopia?.zadanie['includeLocks']).toBe(true);

    const tresc = panel.element.querySelector<HTMLElement>('.dm-odpowiedz')?.textContent ?? '';
    expect(tresc).toContain('OSOBNY dokument');
    expect(tresc).toContain('nie rusza');
    expect(tresc).toContain('historia nie była przenoszona');
  });
});

describe('wniesienie ze źródła', () => {
  it('fragment z Biblioteki niesie zapis pochodzenia oddany przez rdzeń', async () => {
    const zapisy: Zapis[] = [];
    const kanal = atrapaKanalu(zapisy, {
      [Command.StudioInsertFromLibrary]: {
        form: {},
        balance: { applied: 1, skippedCount: 0 },
        provenance: {
          id: 'pochodzenie-1',
          documentId: 'dokument-1',
          kind: 'libraryFile',
          rangeStart: 400,
          rangeEnd: 900,
          libraryFileId: 'plik-1',
          sourceVersion: 'wersja-3',
        },
      },
    });
    const stan = stanZDokumentem(kanal);
    const panel = utworzOsadzeniePanel(
      {
        naSzukanieBiblioteki: () => undefined,
        naPodglad: () => undefined,
        naWniesienieZBiblioteki: () => {
          throw new Error('droga rdzenia ma pierwszeństwo — czynność okna nie może się odpalić');
        },
        naWciagniecieStrony: () => undefined,
        naMigawke: () => undefined,
        naOtwarcieStrony: () => undefined,
        naWniesienieZeStrony: () => undefined,
      },
      { stan, wstawienia: utworzZrodloWstawienStudio(kanal) },
    );
    document.body.append(panel.element);

    panel.pokazPodglad(
      { id: 'plik-1', name: 'wzór pisma.docx' } as never,
      { kind: 'text', text: 'Treść wzoru pisma.', page: 1, pageCount: 1 } as never,
    );
    panel.element.querySelector<HTMLButtonElement>('[data-czynnosc="wnies-biblioteka"]')!.click();
    await new Promise((gotowe) => setTimeout(gotowe, 0));

    const wniesienie = zapisy.find((z) => z.komenda === Command.StudioInsertFromLibrary);
    expect(wniesienie?.zadanie['libraryFileId']).toBe('plik-1');
    expect(wniesienie?.zadanie['documentId']).toBe('dokument-1');

    const tresc = panel.element.querySelector<HTMLElement>('.ms-osadzenie__odpowiedz')?.textContent ?? '';
    expect(tresc).toContain('pochodzenie-1');
    expect(tresc).toContain('400–900');
    expect(tresc).toContain('wersja-3');
    expect(tresc).toContain('przeżywa zamknięcie okna');
  });
});

describe('nowy dokument', () => {
  it('zakłada pusty dokument komendą rdzenia i nie wypiera dokumentu bieżącego', async () => {
    const zapisy: Zapis[] = [];
    const kanal = atrapaKanalu(zapisy, {
      [Command.StudioDocumentCreate]: {
        document: { id: 'dokument-nowy', windowId: 'okno-1', format: 'markdown' },
        form: {},
      },
    });
    const stan = stanZDokumentem(kanal);
    const karty = utworzKartyDokumentow(stan, () => undefined, utworzZrodloWstawienStudio(kanal));
    document.body.append(karty.element);

    karty.element.querySelector<HTMLButtonElement>('[data-czynnosc="nowy-dokument"]')!.click();
    await new Promise((gotowe) => setTimeout(gotowe, 0));

    expect(zapisy.some((z) => z.komenda === Command.StudioDocumentCreate)).toBe(true);
    // Nowa strona staje w zakładce NOWEJ, a nie zamiast pisma, nad którym Operator
    // pracował — zakładek jest teraz dwie.
    expect(karty.ile()).toBe(2);
    expect(stan.dokument()?.id).toBe('dokument-nowy');
    expect(karty.element.querySelector('.ms-karty__zdanie')?.textContent ?? '').toContain(
      'Pusty dokument dokument-nowy',
    );
  });

  it('odmowa założenia pustego dokumentu nie znika w ciszy', async () => {
    const kanal = atrapaKanalu([]);
    const stan = stanZDokumentem(kanal);
    const karty = utworzKartyDokumentow(stan, () => undefined, utworzZrodloWstawienStudio(kanal));
    document.body.append(karty.element);

    karty.element.querySelector<HTMLButtonElement>('[data-czynnosc="nowy-dokument"]')!.click();
    await new Promise((gotowe) => setTimeout(gotowe, 0));

    const zdanie = karty.element.querySelector<HTMLElement>('.ms-karty__zdanie');
    expect(zdanie?.dataset['udane']).toBe('nie');
    expect(zdanie?.textContent ?? '').toContain('odmówił');
  });
});
