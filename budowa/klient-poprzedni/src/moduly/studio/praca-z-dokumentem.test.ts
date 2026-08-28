import { describe, expect, it } from 'vitest';
import {
  StudioAuthor,
  StudioChangeDecision,
  StudioChangeKind,
  type Command,
  type StudioTrackedChange,
} from '../../../../shared/contract';
import type { Kanal } from '../../protokol/kanal';
import { utworzModulStudio } from './modul-studio';
import { rozdzielNaStrony } from './nastawy-strony';
import { ocenaCzytelnosci, policzCzytelnosc, pomiaryKorekty } from './ocena-redaktora';
import { NARZEDZIA_TEKSTU } from './znaczniki-markdown';
import { OPERACJE_PASKA, KATEGORIE_OPERACJI } from './kategorie-operacji';
import { WIELKOSCI_CIAGLE, akcjaWielkosci, ladunekOperacji } from './suwaki-koncepcyjne';
import {
  czytajBloki,
  serializujPowierzchnie,
  wyrysBloku,
  zapiszBloki,
} from './zapis-formatowany';
import { rozdzielNaOdcinki } from './zmiany-modelu';

/** Sprawdza skutek scalenia czterech okien tekstowych w jedno okno pracy z dokumentem, nie kształt plików źródłowych. */
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

describe('wykaz czynności po scaleniu czterech okien', () => {
  it('niesie wszystkie czynności Studio Editora, kanwy, podglądu i panelu różnicy', () => {
    const modul = utworzModulStudio(atrapaKanalu([]));
    const okno = modul.element.querySelector<HTMLElement>("[data-okno='studio.praca-z-dokumentem']");
    expect(okno).not.toBeNull();
    const tekst = (okno as HTMLElement).textContent ?? '';
    const czynnosci = new Set(
      Array.from((okno as HTMLElement).querySelectorAll<HTMLElement>('[data-czynnosc]')).map(
        (element) => element.dataset['czynnosc'],
      ),
    );

    // Studio Editor: wczytanie, zapis, decyzja o wyniku, wstawienie wyniku.
    for (const czynnosc of ['wczytaj', 'zapisz', 'przyjmij', 'odrzuc', 'wstaw']) {
      expect(czynnosci.has(czynnosc), `zginęła czynność „${czynnosc}"`).toBe(true);
    }
    // Preview Window: wydanie w formacie i przekazanie do Library.
    for (const czynnosc of ['wydaj', 'do-library']) {
      expect(czynnosci.has(czynnosc), `zginęła czynność „${czynnosc}"`).toBe(true);
    }
    // Diff/Grep Panel: porównanie, adnotacja, decyzja wybiórcza o fragmentach.
    for (const czynnosc of ['porownaj', 'adnotacja', 'przyjmij-fragmenty']) {
      expect(czynnosci.has(czynnosc), `zginęła czynność „${czynnosc}"`).toBe(true);
    }
    // Znajdź/Zamień i galeria szablonów.
    expect(czynnosci.has('zamien')).toBe(true);
    expect(czynnosci.has('galeria')).toBe(true);
    // Recenzja: śledzenie zmian, skoki, komentarze.
    for (const czynnosc of [
      'przyjmij-wszystkie',
      'odrzuc-wszystkie',
      'zmiana-dalej',
      'zmiana-wstecz',
      'komentarz-nowy',
      'komentarze-widok',
    ]) {
      expect(czynnosci.has(czynnosc), `zginęła czynność „${czynnosc}"`).toBe(true);
    }
    // Tryby widoku niosą podgląd wydania i różnicę jako tryby, nie okna.
    for (const czynnosc of ['tryb-formatowany', 'tryb-zrodlowy', 'tryb-wydanie', 'tryb-roznica']) {
      expect(czynnosci.has(czynnosc), `zginął tryb „${czynnosc}"`).toBe(true);
    }

    // Kanwa tekstowa: pochodzenie wyniku operacji zostało, sama kanwa zniknęła.
    expect(modul.element.querySelector("[data-okno='studio.kanwa-tekstowa']")).toBeNull();
    expect(modul.element.querySelectorAll('textarea').length).toBeGreaterThan(0);
    expect(tekst).toContain('Okno komunikacji modułu');

    // Szesnaście narzędzi znacznikowych i cztery operacje paska zaznaczenia.
    const narzedzia = new Set(
      Array.from((okno as HTMLElement).querySelectorAll<HTMLElement>('[data-narzedzie]')).map(
        (element) => element.dataset['narzedzie'],
      ),
    );
    for (const narzedzie of NARZEDZIA_TEKSTU) {
      expect(
        narzedzia.has(narzedzie.kod) || czynnosci.has(narzedzie.kod),
        `zginęło narzędzie „${narzedzie.kod}"`,
      ).toBe(true);
    }
    const operacje = new Set(
      Array.from((okno as HTMLElement).querySelectorAll<HTMLElement>('[data-operacja]')).map(
        (element) => element.dataset['operacja'],
      ),
    );
    for (const operacja of OPERACJE_PASKA) {
      expect(operacje.has(operacja.id), `zginęła operacja paska „${operacja.id}"`).toBe(true);
    }
    expect(operacje.has('wiecej')).toBe(true);

    modul.rozlacz();
  });

  it('zostawia osobno panele bez powierzchni tekstowej i gasi okna scalone', () => {
    const modul = utworzModulStudio(atrapaKanalu([]));
    const kody = Array.from(
      modul.element.querySelectorAll<HTMLElement>('[data-okno]'),
    ).map((element) => element.dataset['okno']);

    expect(kody).toContain('studio.tools-panel');
    expect(kody).toContain('studio.session-repository');
    expect(kody).toContain('studio.ingest-ocr-panel');
    expect(kody).not.toContain('studio.studio-editor');
    expect(kody).not.toContain('studio.diff-grep-panel');
    expect(kody).not.toContain('studio.preview-window');
    expect(kody).not.toContain('studio.kanwa-tekstowa');

    modul.rozlacz();
  });

  it('prowadzi dwa dokumenty w zakładkach z niezależnym stanem', () => {
    const modul = utworzModulStudio(atrapaKanalu([]));
    const pasek = modul.element.querySelector<HTMLElement>('.ms-karty');
    expect(pasek).not.toBeNull();
    const dodaj = (pasek as HTMLElement).querySelector<HTMLButtonElement>(
      "[data-czynnosc='dodaj-karte']",
    );
    expect(dodaj).not.toBeNull();
    (dodaj as HTMLButtonElement).click();
    expect((pasek as HTMLElement).querySelectorAll('[data-karta]').length).toBe(2);
    expect(
      (pasek as HTMLElement).querySelectorAll("[data-karta][data-czynna='tak']").length,
    ).toBe(1);
    modul.rozlacz();
  });
});

describe('widok formatowany i zapis treści', () => {
  it('czyta i zapisuje bloki bez gubienia składni', () => {
    const tresc = [
      '# Umowa',
      '',
      'Strony ustalają **zakres** oraz *termin*.',
      '',
      '- pozycja pierwsza',
      '- pozycja druga',
      '',
      '> Cytat z opracowania.',
      '',
      '| Kolumna A | Kolumna B |',
      '| --- | --- |',
      '| jeden | dwa |',
    ].join('\n');

    const bloki = czytajBloki(tresc);
    expect(bloki.map((blok) => blok.rodzaj)).toEqual([
      'naglowek-1',
      'tekst',
      'lista',
      'cytat',
      'tabela',
    ]);
    expect(zapiszBloki(bloki)).toBe(tresc);
  });

  it('oddaje z powierzchni tę samą składnię, którą na niej wyrysował', () => {
    const bloki = czytajBloki('## Tytuł\n\nZdanie z **wagą** i <u>podkreśleniem</u>.');
    const pole = document.createElement('div');
    for (const blok of bloki) pole.append(wyrysBloku(blok));

    // Wyrys jest formatowany — pogrubienie jest elementem, nie gwiazdkami.
    expect(pole.querySelector('strong')?.textContent).toBe('wagą');
    expect(pole.querySelector('u')?.textContent).toBe('podkreśleniem');
    expect(serializujPowierzchnie(pole)).toBe(
      '## Tytuł\n\nZdanie z **wagą** i <u>podkreśleniem</u>.',
    );
  });

  it('rozdziela kartki wedle wysokości bloków i honoruje podział jawny', () => {
    // Cztery bloki po 300 punktów na kartce o polu 700 punktów: dwa na pierwszą,
    // dwa na drugą.
    const strony = rozdzielNaStrony(
      4,
      700,
      () => 300,
      () => false,
    );
    expect(strony.map((strona) => strona.length)).toEqual([2, 2]);

    // Podział jawny kończy kartkę niezależnie od miejsca, jakie na niej zostało.
    const zPodzialem = rozdzielNaStrony(
      3,
      700,
      () => 10,
      (numer) => numer === 0,
    );
    expect(zPodzialem.map((strona) => strona.length)).toEqual([1, 2]);
  });
});

describe('zmiany modelu w treści', () => {
  function zmiana(od: number, doZnaku: number): StudioTrackedChange {
    return {
      id: 'studio-zm-1',
      documentId: 'studio-dok-1',
      kind: StudioChangeKind.Wstawienie,
      author: StudioAuthor.Model,
      rangeStart: od,
      rangeEnd: doZnaku,
      before: 'współpracy',
      after: 'współpracy handlowej',
      decision: StudioChangeDecision.Oczekuje,
      createdAt: 0,
    };
  }

  it('oznacza fragment modelu w miejscu i nie gubi ani jednego znaku treści', () => {
    const tresc = 'Strony ustalają zakres współpracy handlowej.';
    const od = tresc.indexOf('współpracy');
    const odcinki = rozdzielNaOdcinki(tresc, [zmiana(od, od + 'współpracy handlowej'.length)]);

    expect(odcinki.map((odcinek) => odcinek.tekst).join('')).toBe(tresc);
    const oznaczony = odcinki.find((odcinek) => odcinek.zmiana !== null);
    expect(oznaczony?.tekst).toBe('współpracy handlowej');
    expect(oznaczony?.zmiana?.author).toBe(StudioAuthor.Model);
  });

  it('pomija zakres spoza treści, zamiast oznaczać niewłaściwe miejsce', () => {
    const tresc = 'Krótka treść.';
    const odcinki = rozdzielNaOdcinki(tresc, [zmiana(100, 140)]);
    expect(odcinki.length).toBe(1);
    expect(odcinki[0]?.zmiana).toBeNull();
    expect(odcinki[0]?.tekst).toBe(tresc);
  });

  it('nie oznacza zmian już rozstrzygniętych', () => {
    const tresc = 'Strony ustalają zakres.';
    const rozstrzygnieta: StudioTrackedChange = {
      ...zmiana(0, 6),
      decision: StudioChangeDecision.Przyjeta,
    };
    expect(rozdzielNaOdcinki(tresc, [rozstrzygnieta]).every((odcinek) => odcinek.zmiana === null)).toBe(
      true,
    );
  });
});

describe('pomiary panelu Redaktora', () => {
  it('liczy ocenę z mglistości, a nie z wrażenia', () => {
    const proste = policzCzytelnosc('Kot pił mleko. Pies jadł kość. Dzień był ciepły.');
    const zawile = policzCzytelnosc(
      'Przedmiotem niniejszego postanowienia jest uszczegółowienie obowiązków ' +
        'sprawozdawczych wynikających z wcześniejszych ustaleń międzyinstytucjonalnych.',
    );
    expect(proste.slowNaZdanie).toBeLessThan(zawile.slowNaZdanie);
    expect(ocenaCzytelnosci(proste)).toBeGreaterThan(ocenaCzytelnosci(zawile));
    expect(ocenaCzytelnosci(policzCzytelnosc(''))).toBe(0);
  });

  it('nazywa brak pomiaru zamiast pokazywać liczbę wymyśloną', () => {
    const pomiary = pomiaryKorekty('Zdanie z błędem interpunkcyjnym ,tutaj.');
    const pisownia = pomiary.find((pomiar) => pomiar.kod === 'pisownia');
    expect(pisownia?.wartosc).toBeNull();
    expect(pisownia?.podstawa).toContain('BEZ POMIARU');

    const interpunkcja = pomiary.find((pomiar) => pomiar.kod === 'interpunkcja');
    expect(interpunkcja?.wartosc).toBeGreaterThan(0);
    expect(interpunkcja?.podstawa).toContain('Odstęp przed znakiem przestankowym');
  });
});

describe('suwaki koncepcyjne', () => {
  it('wskazują przeciwne operacje na przeciwnych końcach skali', () => {
    const objetosc = WIELKOSCI_CIAGLE.find((wielkosc) => wielkosc.kod === 'objetosc');
    expect(objetosc).toBeDefined();
    expect(akcjaWielkosci(objetosc as never, -2)).toBe('studio.styl.skrocenie');
    expect(akcjaWielkosci(objetosc as never, 2)).toBe('studio.styl.rozwiniecie');
  });

  it('nie wysyłają nastaw neutralnych, a polecenie własne wysyłają', () => {
    const ladunek = ladunekOperacji({ objetosc: 0, ton: 2 }, 'skróć wstęp');
    expect(ladunek['objetosc']).toBeUndefined();
    expect(ladunek['ton']).toBe(2);
    expect(ladunek['polecenie']).toBe('skróć wstęp');
  });

  it('bierze operacje z wykazu kategorii, a nie z własnej listy', () => {
    const wszystkie = new Set(
      KATEGORIE_OPERACJI.flatMap((kategoria) => kategoria.operacje.map((operacja) => operacja.id)),
    );
    for (const wielkosc of WIELKOSCI_CIAGLE) {
      expect(wszystkie.has(wielkosc.idAkcji), `operacja ${wielkosc.idAkcji} spoza wykazu`).toBe(true);
    }
  });
});
