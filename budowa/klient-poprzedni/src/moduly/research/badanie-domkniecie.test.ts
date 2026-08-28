import { describe, expect, it } from 'vitest';

import {
  Command,
  ResearchAnnotationKind,
  ResearchAttachmentKind,
  type ResearchSource,
} from '../../../../shared/contract';
import type { Kanal, Wynik } from '../../protokol/kanal';
import {
  AKCJE_LEKTURY,
  AKCJE_USTALEN,
  AKCJE_WORKSPACE,
  AKCJE_ZRODEL,
  type AkcjaBadania,
} from './akcje-okien';
import { kodyZTekstu, zlaczKody } from './badanie-ksiazka-kodow';
import { rozbrojZdjecie } from './badanie-zdjecie-adnotacji';
import { utworzStanBadania } from './stan-badania';
import { czyKomendaBadania, wykonajKomendeBadania } from './wywolania-komend';

/** Cztery czynności badania, które miały obsługę w rdzeniu i nie miały drogi z okna. */

/** Kanał próbny: zapamiętuje żądania i oddaje odpowiedź wskazaną per komenda w tym sprawdzianie badania. */
function kanalProbny(odpowiedzi: Record<string, unknown>): {
  kanal: Kanal;
  wyslane: { komenda: string; zadanie: unknown }[];
} {
  const wyslane: { komenda: string; zadanie: unknown }[] = [];
  const kanal = {
    wyslij(komenda: string, zadanie: unknown, przyWyniku?: (wynik: Wynik<unknown>) => void) {
      wyslane.push({ komenda, zadanie });
      const tresc = odpowiedzi[komenda];
      // Odpowiedź wraca po oddaniu identyfikatora żądania, tak jak w kanale prawdziwym.
      queueMicrotask(() =>
        przyWyniku?.(
          tresc === undefined
            ? { udany: false, blad: { code: 'not_found', message: 'brak wykonawcy', retryable: false } }
            : { udany: true, wynik: tresc },
        ),
      );
      return `zadanie-${String(wyslane.length)}`;
    },
    naZdarzenie: () => () => undefined,
    naDowolny: () => () => undefined,
    sesja: () => ({}) as ReturnType<Kanal['sesja']>,
    dziennikNieznanych: () => ({}) as ReturnType<Kanal['dziennikNieznanych']>,
  } as unknown as Kanal;
  return { kanal, wyslane };
}

/** Stan badania z oknem wskazanym przez rdzeń — bez niego komendy badania w ogóle by nigdy nie ruszyły. */
function stanZOknem(odpowiedzi: Record<string, unknown>) {
  const { kanal, wyslane } = kanalProbny(odpowiedzi);
  const stan = utworzStanBadania(kanal);
  stan.ustawOkno('okno-badania-1');
  return { stan, wyslane };
}

/** Źródło w postaci, w której rdzeń je oddaje, gotowe do wywołania wszystkich komend badania w sprawdzianie. */
function zrodlo(zmiany: Partial<ResearchSource> = {}): ResearchSource {
  return {
    id: 'zrodlo-1',
    windowId: 'okno-badania-1',
    title: 'Raport rynkowy',
    kind: 'report',
    ...zmiany,
  } as ResearchSource;
}

/** Akcja o tym kodzie w katalogu okna; jej brak jest brakiem chwytu prowadzącego wprost do rdzenia aplikacji. */
function chwyt(akcje: readonly AkcjaBadania[], kod: string): AkcjaBadania {
  const znaleziona = akcje.find((akcja) => akcja.kod === kod);
  expect(znaleziona, `okno musi mieć jawny chwyt dla ${kod}`).toBeDefined();
  return znaleziona as AkcjaBadania;
}

describe('Domknięcie modułu Research — chwyt w oknie i wywołanie w rozdzielniku', () => {
  const domykane: readonly [string, readonly AkcjaBadania[]][] = [
    [Command.ResearchWorkspaceGet, AKCJE_WORKSPACE],
    [Command.ResearchSourceAttachmentAdd, AKCJE_ZRODEL],
    [Command.ResearchAnnotationRemove, AKCJE_LEKTURY],
    [Command.ResearchCodebookSet, AKCJE_USTALEN],
  ];

  for (const [komenda, akcje] of domykane) {
    it(`${komenda} ma chwyt w oknie i wywołanie w oknie`, () => {
      const akcja = chwyt(akcje, komenda);
      expect(akcja.droga).toBe('komenda');
      expect(czyKomendaBadania(komenda)).toBe(true);
    });
  }
});

describe('research.workspace.get — przestrzeń badania', () => {
  it('wciąga zakres i etapy do pamięci badania i liczy pytania bez pokrycia', async () => {
    const { stan, wyslane } = stanZOknem({
      [Command.ResearchWorkspaceGet]: {
        scope: 'rynek azjatycki',
        stages: ['rozpoznanie', 'zbieranie'],
        questions: [
          { id: 'pytanie-1', text: 'kto dostarcza', sourceIds: ['zrodlo-1'] },
          { id: 'pytanie-2', text: 'jaka cena' },
        ],
        updatedAt: 1_700_000_000_000,
      },
    });
    const wynik = await wykonajKomendeBadania(
      { stan },
      chwyt(AKCJE_WORKSPACE, Command.ResearchWorkspaceGet),
    );

    expect(wyslane.map((wpis) => wpis.komenda)).toEqual([Command.ResearchWorkspaceGet]);
    expect(wynik.udany).toBe(true);
    expect(wynik.opis).toContain('pytań badawczych 2');
    expect(wynik.opis).toContain('bez ani jednego źródła 1');
    // Braki odbiorcy, protokołu i notatki są nazwane, nie przemilczane.
    expect(wynik.opis).toContain('odbiorcy raportu');
    expect(stan.zakres()).toBe('rynek azjatycki');
    expect(stan.etapy()).toEqual(['rozpoznanie', 'zbieranie']);
  });
});

describe('research.source.attachment.add — załącznik pełnego tekstu', () => {
  it('bez wskazanego źródła mówi, czego brakuje, i nie woła rdzenia', async () => {
    const { stan, wyslane } = stanZOknem({});
    const wynik = await wykonajKomendeBadania(
      { stan },
      chwyt(AKCJE_ZRODEL, Command.ResearchSourceAttachmentAdd),
    );

    expect(wynik.udany).toBe(false);
    expect(wynik.opis).toContain('Zaznacz źródło');
    expect(wyslane).toHaveLength(0);
  });

  it('wskazany plik jedzie w polu ścieżki, a źródło z dokumentem repozytorium — w polu dokumentu', async () => {
    const { stan, wyslane } = stanZOknem({
      [Command.ResearchSourceAttachmentAdd]: {
        attachment: {
          id: 'zalacznik-1',
          sourceId: 'zrodlo-1',
          kind: ResearchAttachmentKind.Fulltext,
          createdAt: 1_700_000_000_000,
          sizeBytes: 4096,
        },
      },
    });
    stan.wchlonZrodlo(zrodlo({ libraryFileId: 'plik-1' }));
    stan.lektura.wskaz('zrodlo-1');
    const akcja = chwyt(AKCJE_ZRODEL, Command.ResearchSourceAttachmentAdd);

    const zPliku = await wykonajKomendeBadania({ stan, tekst: '/dom/raport.pdf' }, akcja);
    expect(zPliku.udany).toBe(true);
    expect(zPliku.opis).toContain('bajtów 4096');
    expect(wyslane[0].zadanie).toEqual({
      sourceId: 'zrodlo-1',
      kind: ResearchAttachmentKind.Fulltext,
      sourcePath: '/dom/raport.pdf',
    });

    await wykonajKomendeBadania({ stan, tekst: '' }, akcja);
    expect(wyslane[1].zadanie).toEqual({
      sourceId: 'zrodlo-1',
      kind: ResearchAttachmentKind.Fulltext,
      libraryFileId: 'plik-1',
    });
  });

  it('źródło bez pliku i bez dokumentu repozytorium nie wysyła żądania pustego', async () => {
    const { stan, wyslane } = stanZOknem({});
    stan.wchlonZrodlo(zrodlo());
    stan.lektura.wskaz('zrodlo-1');
    const wynik = await wykonajKomendeBadania(
      { stan, tekst: '' },
      chwyt(AKCJE_ZRODEL, Command.ResearchSourceAttachmentAdd),
    );

    expect(wynik.udany).toBe(false);
    expect(wynik.opis).toContain('załącznik musi mieć treść');
    expect(wyslane).toHaveLength(0);
  });
});

describe('research.annotation.remove — czynność nieodwracalna', () => {
  const adnotacja = {
    id: 'adnotacja-1',
    sourceId: 'zrodlo-1',
    kind: ResearchAnnotationKind.Highlight,
    quote: 'udział rynkowy wzrósł',
    anchor: { kind: 'page' },
    createdAt: 1_700_000_000_000,
  };

  it('pierwsze naciśnięcie ostrzega i NIE zdejmuje, drugie zdejmuje', async () => {
    rozbrojZdjecie();
    const { stan, wyslane } = stanZOknem({
      [Command.ResearchAnnotationList]: { annotations: [adnotacja] },
      [Command.ResearchAnnotationRemove]: { annotationId: 'adnotacja-1' },
    });
    stan.lektura.wskaz('zrodlo-1');
    const akcja = chwyt(AKCJE_LEKTURY, Command.ResearchAnnotationRemove);

    const zapowiedz = await wykonajKomendeBadania({ stan }, akcja);
    expect(zapowiedz.udany).toBe(false);
    expect(zapowiedz.opis).toContain('NIEODWRACALNE');
    // Ostrzeżenie nazywa adnotację imiennie, nie identyfikatorem samym.
    expect(zapowiedz.opis).toContain('udział rynkowy wzrósł');
    expect(wyslane.map((wpis) => wpis.komenda)).toEqual([Command.ResearchAnnotationList]);

    const zdjecie = await wykonajKomendeBadania({ stan }, akcja);
    expect(zdjecie.udany).toBe(true);
    expect(zdjecie.opis).toContain('bezpowrotnie');
    expect(wyslane.map((wpis) => wpis.komenda)).toContain(Command.ResearchAnnotationRemove);

    // Po zdjęciu uzbrojenie znika: następne naciśnięcie ostrzega od nowa.
    const znowu = await wykonajKomendeBadania({ stan }, akcja);
    expect(znowu.opis).toContain('NIEODWRACALNE');
  });

  it('materiał bez adnotacji mówi o pustce, a nie o niepowodzeniu', async () => {
    rozbrojZdjecie();
    const { stan, wyslane } = stanZOknem({
      [Command.ResearchAnnotationList]: { annotations: [] },
    });
    stan.lektura.wskaz('zrodlo-1');
    const wynik = await wykonajKomendeBadania(
      { stan },
      chwyt(AKCJE_LEKTURY, Command.ResearchAnnotationRemove),
    );

    expect(wynik.opis).toContain('nie ma ani jednej adnotacji');
    expect(wyslane.map((wpis) => wpis.komenda)).toEqual([Command.ResearchAnnotationList]);
  });
});

describe('research.codebook.set — książka kodów', () => {
  it('rozkłada tekst na kody i nie gubi kodów zastanych', () => {
    const wpisane = kodyZTekstu('cena | koszt zakupu\n\nlogistyka\n  cena  ');
    expect(wpisane).toEqual([
      { nazwa: 'cena', definicja: 'koszt zakupu' },
      { nazwa: 'logistyka', definicja: '' },
      { nazwa: 'cena', definicja: '' },
    ]);

    const zlaczenie = zlaczKody(
      [{ id: 'kod-1', name: 'Cena', occurrences: 4 }],
      [
        { nazwa: 'cena', definicja: 'koszt zakupu' },
        { nazwa: 'logistyka', definicja: '' },
      ],
    );
    expect(zlaczenie.dopisane).toBe(1);
    expect(zlaczenie.poprawione).toBe(1);
    expect(zlaczenie.kody).toHaveLength(2);
    // Kod zastany zachowuje identyfikator i liczność nadaną przez rdzeń.
    expect(zlaczenie.kody[0]).toEqual({
      id: 'kod-1',
      name: 'Cena',
      occurrences: 4,
      description: 'koszt zakupu',
    });
    expect(zlaczenie.kody[1]).toEqual({ id: '', name: 'logistyka' });
  });

  it('zapisuje książkę kodów po odczycie i mówi, ile kodów zostało zachowanych', async () => {
    const { stan, wyslane } = stanZOknem({
      [Command.ResearchCodebookGet]: { codes: [{ id: 'kod-1', name: 'Cena' }] },
      [Command.ResearchCodebookSet]: {
        codes: [
          { id: 'kod-1', name: 'Cena', description: 'koszt zakupu' },
          { id: 'kod-2', name: 'logistyka' },
        ],
      },
    });
    const wynik = await wykonajKomendeBadania(
      { stan, tekst: 'cena | koszt zakupu\nlogistyka' },
      chwyt(AKCJE_USTALEN, Command.ResearchCodebookSet),
    );

    expect(wyslane.map((wpis) => wpis.komenda)).toEqual([
      Command.ResearchCodebookGet,
      Command.ResearchCodebookSet,
    ]);
    expect(wyslane[1].zadanie).toEqual({
      windowId: 'okno-badania-1',
      codes: [
        { id: 'kod-1', name: 'Cena', description: 'koszt zakupu' },
        { id: '', name: 'logistyka' },
      ],
    });
    expect(wynik.udany).toBe(true);
    expect(wynik.opis).toContain('dopisanych 1');
  });

  it('bez kodów w polu nie wysyła zapisu, który wymazałby książkę kodów', async () => {
    const { stan, wyslane } = stanZOknem({});
    const wynik = await wykonajKomendeBadania(
      { stan, tekst: '   ' },
      chwyt(AKCJE_USTALEN, Command.ResearchCodebookSet),
    );

    expect(wynik.udany).toBe(false);
    expect(wynik.opis).toContain('wymazałby książkę kodów');
    expect(wyslane).toHaveLength(0);
  });
});
