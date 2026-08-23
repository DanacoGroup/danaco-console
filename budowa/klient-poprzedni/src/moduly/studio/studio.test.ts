import { describe, expect, it } from 'vitest';
import {
  Command,
  DiffHunkKind,
  StudioIngestState,
  type StudioIngestItem,
  type StudioVersion,
} from '../../../../shared/contract';
import type { Kanal } from '../../protokol/kanal';
import { przefiltrujHistorie } from './filtr-historii';
import { utworzKolejkeCyfryzacji } from './kolejka-cyfryzacji';
import { formatZrodlowy } from './konwersja-dokumentu';
import { opiszLiczniki, policzTresc } from './liczniki-dokumentu';
import { utworzOknoIngestOcrPanel } from './okno-ingest-ocr-panel';
import { policzRoznice, przefiltrujRoznice } from './statystyka-roznicy';
import { utworzStanStudio } from './stan-studio';
import { znajdzTrafienia, zamienWszystkie } from './wyszukiwanie-tekstu';
import { utworzZrodloDokumentuStudio } from './zrodlo-dokumentu-studio';
import { utworzZrodloMaterialuStudio } from './zrodlo-materialu-studio';
import { utworzZrodloWstawienStudio } from './zrodlo-wstawien-studio';

/**
 * Sprawdziany modułu Studio — czynność, nie kształt pliku.
 *
 * Pilnowane jest to, czego opracowanie żąda wprost, a co da się wykonać bez
 * rdzenia i co przy poprawce łatwo zepsuć po cichu:
 *
 *   1. kolejka cyfryzacji odbija PIĘĆ stanów rdzenia (wraz z ponowieniem poniżej
 *      progu pewności) i nie gubi bilansu, gdy jedna pozycja odmawia (rozdz. 3.8:
 *      wsad idzie dalej mimo błędu);
 *   2. liczniki dokumentu liczą to, co pasek statusu obiecuje (rozdz. 3.3);
 *   3. Znajdź/Zamień traktuje frazę dosłownie, a złą składnię wzorca nazywa,
 *      zamiast oddawać ją jako brak trafień (rozdz. 3.3, funkcja E5);
 *   4. statystyka różnicy liczy fragment zmieniony do obu stron (rozdz. 3.5);
 *   5. filtr historii zawęża po polach, które wersja naprawdę niesie;
 *   6. narzędziownia cyfryzacji prowadzi kolejkę RDZENIA rodziną studio.ingest.*:
 *      dokłada materiał, rozpoznaje go z nastawami, które naprawdę jadą do rdzenia,
 *      daje poprawić rozpoznane słowo PRZED przyjęciem i kończy dokumentem wraz
 *      z pierwszą wersją — dowód, że droga od pliku do dokumentu jest cała;
 *   7. pusty wykaz urządzeń wejściowych jest nazwany jako brak maszyny rdzenia,
 *      a nie jako odmowa produktu — i nazwany PRZED próbą, nie po niej.
 *
 * Rdzeń jest atrapą: sprawdzian pyta o zachowanie modułu, nie serwera.
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

/** Pozycja kolejki w kształcie kontraktu — atrapa odpowiedzi rdzenia. */
function pozycjaKolejki(
  id: string,
  state: StudioIngestState,
  reszta: Partial<StudioIngestItem> = {},
): StudioIngestItem {
  return { id, windowId: 'okno-1', state, createdAt: 1, ...reszta };
}

describe('kolejka cyfryzacji', () => {
  it('odbija pięć stanów rdzenia i liczy bilans mimo odmowy jednej pozycji', () => {
    const kolejka = utworzKolejkeCyfryzacji();
    kolejka.ustawPozycje([
      pozycjaKolejki('pozycja-1', StudioIngestState.Gotowa, {
        sourcePath: '/dane/umowa-strona-1.tif',
        text: 'treść strony',
      }),
      pozycjaKolejki('pozycja-2', StudioIngestState.Odmowa, {
        assetId: 'zasob-umowa-2',
        failureReason: 'plik nieczytelny',
      }),
      pozycjaKolejki('pozycja-3', StudioIngestState.Oczekuje, {
        sourcePath: '/dane/umowa-strona-3.tif',
      }),
      pozycjaKolejki('pozycja-4', StudioIngestState.Ponowienie, {
        sourcePath: '/dane/umowa-strona-4.tif',
        confidence: 0.4,
      }),
    ]);

    const bilans = kolejka.bilans();
    expect(bilans.wszystkie).toBe(4);
    expect(bilans.gotowe).toBe(1);
    expect(bilans.odmowy).toBe(1);
    expect(bilans.oczekujace).toBe(1);
    // Ponowienie jest stanem OSOBNYM: pozycja poniżej progu pewności nie jest ani
    // gotowa, ani odmówiona, a zlanie jej z którymkolwiek z tych dwóch kłamałoby
    // o tym, co Operator ma z nią zrobić.
    expect(bilans.ponowienia).toBe(1);
    expect(bilans.zTekstem).toBe(1);
    // Odmowa jednej pozycji nie zdejmuje pozostałych: następna do rozpoznania to
    // pierwsza oczekująca, a po niej wraca ta z ponowienia.
    expect(kolejka.nastepnaDoRozpoznania()?.id).toBe('pozycja-3');
  });

  it('nie gubi wskazania Operatora przy odświeżeniu kolejki ani słów do poprawki', () => {
    const kolejka = utworzKolejkeCyfryzacji();
    kolejka.ustawPozycje([pozycjaKolejki('pozycja-1', StudioIngestState.Gotowa)]);
    kolejka.ustawWskazana('pozycja-1');
    kolejka.ustawRozpoznanie(
      'pozycja-1',
      [
        {
          index: 0,
          text: 'Umowa',
          page: 1,
          x: 10,
          y: 20,
          width: 40,
          height: 12,
          confidence: 0.42,
        },
      ],
      undefined,
    );

    // Odświeżenie kolejki z rdzenia nie może przestawiać Operatora na inną
    // pozycję ani zdejmować słów, na których stoi poprawianie: `queue.list`
    // słów nie powtarza, więc jedynym ich miejscem jest odbicie okna.
    kolejka.ustawPozycje([
      pozycjaKolejki('pozycja-1', StudioIngestState.Gotowa, { text: 'Umowa' }),
      pozycjaKolejki('pozycja-2', StudioIngestState.Oczekuje),
    ]);
    expect(kolejka.wskazana()).toBe('pozycja-1');
    expect(kolejka.slowa('pozycja-1')).toHaveLength(1);

    // Pozycja, której rdzeń już nie oddaje, przestaje być wskazana — wskazanie na
    // byt nieistniejący byłoby czynnością bez przedmiotu.
    kolejka.ustawPozycje([pozycjaKolejki('pozycja-2', StudioIngestState.Oczekuje)]);
    expect(kolejka.wskazana()).toBe('');
  });
});

describe('liczniki dokumentu', () => {
  it('liczy słowa, zdania i akapity treści redakcyjnej', () => {
    const tresc =
      'Umowa wchodzi w życie z dniem podpisania. Strony ustalają termin dostawy.\n\n' +
      'Załącznik stanowi integralną część umowy.';
    const liczniki = policzTresc(tresc);
    // 7 słów zdania pierwszego + 4 drugiego + 5 akapitu drugiego.
    expect(liczniki.slowa).toBe(16);
    expect(liczniki.zdania).toBe(3);
    expect(liczniki.akapity).toBe(2);
    expect(liczniki.znaki).toBe(tresc.length);
    expect(liczniki.znakiBezOdstepow).toBeLessThan(liczniki.znaki);
    // Tekst krótszy niż minuta czytania nie dostaje zera minut — zero znaczyłoby
    // „nie ma czego czytać", a jest co.
    expect(liczniki.minutyCzytania).toBe(1);
  });

  it('dokument pusty nazywa się pustym, a nie zerami', () => {
    expect(policzTresc('').minutyCzytania).toBe(0);
    expect(opiszLiczniki(policzTresc(''))).toBe('dokument pusty');
  });
});

describe('Znajdź/Zamień', () => {
  it('traktuje frazę dosłownie — kropka we frazie znaczy kropkę', () => {
    const tresc = 'wersja 2.0 oraz wersja 250';
    const wynik = znajdzTrafienia(tresc, {
      wzorzec: 'wersja 2.0',
      regularne: false,
      wielkoscLiter: true,
    });
    expect(wynik.powod).toBe('');
    expect(wynik.trafienia).toHaveLength(1);
    expect(wynik.trafienia[0]?.tekst).toBe('wersja 2.0');
  });

  it('złą składnię wzorca nazywa powodem, a nie brakiem trafień', () => {
    const wynik = znajdzTrafienia('dowolna treść', {
      wzorzec: '(niedomknięty',
      regularne: true,
      wielkoscLiter: false,
    });
    expect(wynik.trafienia).toHaveLength(0);
    expect(wynik.powod).not.toBe('');
    expect(wynik.powod).toContain('wyrażeniem regularnym');
  });

  it('zamienia wszystkie trafienia, nie gubiąc treści przy różnej długości zamiennika', () => {
    const wynik = zamienWszystkie(
      'Wykonawca odpowiada. Wykonawca dostarcza.',
      { wzorzec: 'Wykonawca', regularne: false, wielkoscLiter: true },
      'Zleceniobiorca',
    );
    expect(wynik.liczba).toBe(2);
    expect(wynik.tresc).toBe('Zleceniobiorca odpowiada. Zleceniobiorca dostarcza.');
  });
});

describe('statystyka różnicy', () => {
  const fragmenty = [
    { index: 1, kind: DiffHunkKind.Added, after: 'zdanie dopisane w całości' },
    { index: 2, kind: DiffHunkKind.Removed, before: 'zdanie skreślone' },
    { index: 3, kind: DiffHunkKind.Changed, before: 'stary zapis', after: 'nowy zapis redakcyjny' },
    { index: 4, kind: DiffHunkKind.Context, before: 'bez zmian', after: 'bez zmian' },
  ];

  it('liczy fragment zmieniony do obu stron naraz', () => {
    const statystyka = policzRoznice(fragmenty);
    expect(statystyka.fragmenty).toBe(4);
    expect(statystyka.dodane).toBe(1);
    expect(statystyka.usuniete).toBe(1);
    expect(statystyka.zmienione).toBe(1);
    // Dodane: 4 słowa fragmentu dodanego + 3 słowa strony „po" fragmentu
    // zmienionego. Kontekst nie liczy się do żadnej strony.
    expect(statystyka.slowaDodane).toBe(7);
    // Usunięte: 2 słowa fragmentu usuniętego + 2 słowa strony „przed".
    expect(statystyka.slowaUsuniete).toBe(4);
  });

  it('filtr zawęża wykaz, nie zmieniając odpowiedzi rdzenia', () => {
    expect(przefiltrujRoznice(fragmenty, DiffHunkKind.Added)).toHaveLength(1);
    expect(przefiltrujRoznice(fragmenty, 'wszystkie')).toHaveLength(4);
    // Statystyka liczy się z kompletu, nie z wykazu widocznego.
    expect(policzRoznice(fragmenty).fragmenty).toBe(4);
  });
});

describe('filtr historii', () => {
  const wersje: StudioVersion[] = [
    { id: 'w1', documentId: 'd1', createdAt: 1, label: 'do akceptacji klienta' },
    { id: 'w2', documentId: 'd1', createdAt: 2, summary: 'korekta stylu' },
    { id: 'w3', documentId: 'd1', createdAt: 3 },
  ];

  it('zawęża po polach, które wersja naprawdę niesie', () => {
    expect(przefiltrujHistorie(wersje, 'etykietowane').map((w) => w.id)).toEqual(['w1']);
    expect(przefiltrujHistorie(wersje, 'zopisem').map((w) => w.id)).toEqual(['w2']);
    expect(przefiltrujHistorie(wersje, 'wszystkie')).toHaveLength(3);
  });
});

describe('format źródłowy zamiany', () => {
  it('treść dokumentu prowadzonego jako PDF idzie do zamiany jako tekst czysty', () => {
    // Zamiana dostaje napis, a nie plik: nazwanie go PDF-em kazałoby rdzeniowi
    // szukać struktury, której w napisie nie ma.
    expect(formatZrodlowy('pdf')).toBe('txt');
    expect(formatZrodlowy('docx')).toBe('txt');
    expect(formatZrodlowy('markdown')).toBe('markdown');
  });
});

describe('narzędziownia cyfryzacji', () => {
  it('prowadzi kolejkę rdzenia i rozpoznaje pozycję z pełnym sterowaniem', async () => {
    const zapisy: Zapis[] = [];
    const kanal = atrapaKanalu(zapisy, {
      [Command.StudioIngestQueueAdd]: {
        items: [pozycjaKolejki('pozycja-1', StudioIngestState.Oczekuje, {
          sourcePath: '/dane/umowa-skan.tif',
        })],
      },
      [Command.StudioIngestRecognize]: {
        item: pozycjaKolejki('pozycja-1', StudioIngestState.Gotowa, {
          sourcePath: '/dane/umowa-skan.tif',
          text: 'Treść odczytana ze skanu umowy.',
          usedOcr: true,
          confidence: 0.91,
        }),
        words: [
          { index: 0, text: 'Tresc', page: 1, x: 10, y: 20, width: 40, height: 12, confidence: 0.4 },
        ],
      },
      [Command.StudioIngestCorrectionSet]: {
        item: pozycjaKolejki('pozycja-1', StudioIngestState.Gotowa, {
          text: 'Treść odczytana ze skanu umowy.',
        }),
      },
      [Command.StudioIngestItemAccept]: {
        document: { id: 'dokument-1', windowId: 'okno-1', title: 'Umowa ze skanu', format: 'txt' },
        version: { id: 'wersja-1' },
      },
    });
    const stan = utworzStanStudio(kanal);
    stan.ustawOkno('okno-1');
    const panel = utworzOknoIngestOcrPanel(
      stan,
      utworzZrodloDokumentuStudio(kanal),
      utworzZrodloMaterialuStudio(kanal),
      utworzZrodloWstawienStudio(kanal),
    );
    document.body.append(panel.element);

    const wybory = panel.element.querySelectorAll<HTMLSelectElement>('select');
    const zrodlo = wybory[0];
    expect(zrodlo, 'panel pyta o źródło materiału').not.toBeUndefined();
    zrodlo!.value = 'sciezka';

    // Nastawy rozpoznawania mają skutek, a nie stoją jako ozdoba: silnik, zestaw
    // języków i próg pewności jadą do rdzenia w polu `settings`.
    const jezyki = [...panel.element.querySelectorAll<HTMLInputElement>('input[type="text"]')];
    const wskazanie = jezyki[0];
    wskazanie!.value = '/dane/umowa-skan.tif';
    const poleJezykow = jezyki.find((pole) =>
      (pole.previousElementSibling?.textContent ?? '').includes('Zestaw języków'),
    );
    expect(poleJezykow, 'panel pyta o zestaw języków, nie o jeden język').not.toBeUndefined();
    poleJezykow!.value = 'pol, eng';

    panel.element.querySelector<HTMLButtonElement>('[data-czynnosc="doloz"]')!.click();
    await new Promise((gotowe) => setTimeout(gotowe, 0));

    const dolozenie = zapisy.find((z) => z.komenda === Command.StudioIngestQueueAdd);
    expect(dolozenie, 'kolejka stoi po stronie rdzenia').toBeDefined();
    expect(dolozenie?.zadanie['sourcePaths']).toEqual(['/dane/umowa-skan.tif']);
    expect(dolozenie?.zadanie['settings']).toMatchObject({ languages: ['pol', 'eng'] });

    panel.element.querySelector<HTMLButtonElement>('[data-czynnosc="rozpoznaj"]')!.click();
    await new Promise((gotowe) => setTimeout(gotowe, 0));

    const rozpoznanie = zapisy.find((z) => z.komenda === Command.StudioIngestRecognize);
    expect(rozpoznanie?.zadanie['itemId']).toBe('pozycja-1');

    // Słowo rozpoznane poprawia się PRZED przyjęciem — na warstwie tekstowej
    // pozycji kolejki, a nie w dokumencie, którego jeszcze nie ma.
    const poprawka = panel.element.querySelector<HTMLButtonElement>('button[data-slowo="0"]');
    expect(poprawka, 'panel daje poprawić rozpoznane słowo').not.toBeNull();
    poprawka!.click();
    await new Promise((gotowe) => setTimeout(gotowe, 0));
    const korekta = zapisy.find((z) => z.komenda === Command.StudioIngestCorrectionSet);
    expect(korekta?.zadanie['wordIndex']).toBe(0);

    panel.element.querySelector<HTMLButtonElement>('[data-czynnosc="do-edytora"]')!.click();
    await new Promise((gotowe) => setTimeout(gotowe, 0));

    const przyjecie = zapisy.find((z) => z.komenda === Command.StudioIngestItemAccept);
    expect(przyjecie?.zadanie['itemIds']).toEqual(['pozycja-1']);
    // Przyjęcie kończy się DOKUMENTEM wraz z pierwszą wersją, a nie samą treścią
    // w buforze — to była różnica, której starsza droga nie umiała pokryć.
    expect(stan.dokument()?.id).toBe('dokument-1');
  });

  it('pusty wykaz urządzeń nazywa brak maszyny, a nie odmowę produktu', async () => {
    const zapisy: Zapis[] = [];
    const kanal = atrapaKanalu(zapisy, {
      [Command.StudioIngestDeviceList]: { devices: [] },
    });
    const stan = utworzStanStudio(kanal);
    stan.ustawOkno('okno-1');
    const panel = utworzOknoIngestOcrPanel(
      stan,
      utworzZrodloDokumentuStudio(kanal),
      utworzZrodloMaterialuStudio(kanal),
      utworzZrodloWstawienStudio(kanal),
    );
    document.body.append(panel.element);

    panel.element.querySelector<HTMLButtonElement>('[data-czynnosc="urzadzenia"]')!.click();
    await new Promise((gotowe) => setTimeout(gotowe, 0));

    expect(zapisy.some((z) => z.komenda === Command.StudioIngestDeviceList)).toBe(true);
    const zdania = [...panel.element.querySelectorAll('p')].map((p) => p.textContent ?? '');
    expect(
      zdania.some((zdanie) => zdanie.includes('nie widzi ani jednego skanera')),
      'okno mówi, że urządzenia nie ma na maszynie rdzenia — przed próbą, nie po niej',
    ).toBe(true);
  });
});
