import { describe, expect, it } from 'vitest';

import {
  Command,
  ConfigScope,
  StudioAuthor,
  StudioChangeDecision,
  StudioChangeKind,
  type Action,
  type StudioAnnotation,
  type StudioComment,
  type StudioOperation,
  type StudioTrackedChange,
} from '../../../../shared/contract';
import { WSZYSTKIE_OPERACJE, operacjePoFrazie } from './kategorie-operacji';
import { przybornikZlozDrzewo } from './przybornik-katalog';
import {
  RodzajZnakowania,
  przybornikOpiszWykaz,
  przybornikPolicz,
  przybornikPrzefiltruj,
  przybornikZlozZnakowania,
} from './przybornik-wykaz';
import {
  TrybOperacji,
  przybornikNaWierzchu,
  przybornikOdczytajPrzypiete,
  przybornikOdczytajTryb,
  przybornikOdczytajUzycie,
} from './przybornik-uzycie';
import { utworzZnacznikiWlasne } from './przybornik-znaczniki';
import {
  ZrodloWniesienia,
  osadzenieOpiszPochodzenia,
  osadzenieWierszPochodzenia,
  utworzWykazPochodzen,
} from './osadzenie-pochodzenie';

/**
 * Sprawdziany odcinka znakowania, asystenta, schowka i osadzenia.
 *
 * Mierzą rzeczy, które da się zmierzyć bez stawiania okna, i tylko te, w których
 * pomyłka jest cicha: rozróżnienie trzech bytów marginesu, zawężenie wykazu
 * znakowań, kolejność czynności pływaka liczona z użycia, zapora przed
 * pokazaniem tej samej operacji dwa razy w katalogu oraz zapis pochodzenia.
 *
 * Czego tu NIE ma: sprawdzianu, że przycisk wywołuje komendę. Sprawdzian takiego
 * kształtu mierzy atrapę, którą sam stawia, a nie skutek na dokumencie.
 */

/* ── Materiał sprawdzianów ─────────────────────────────────────────────────── */

/** Komentarz o wskazanym autorze, zakresie i stanie wątku. */
function przybornikKomentarz(
  kod: string,
  autor: StudioAuthor,
  poczatek: number,
  rozwiazany: boolean,
): StudioComment {
  return {
    id: kod,
    documentId: 'dokument-1',
    author: autor,
    body: `treść ${kod}`,
    selectionStart: poczatek,
    selectionEnd: poczatek + 10,
    resolved: rozwiazany,
    createdAt: 1_000,
  };
}

/** Zmiana śledzona oczekująca decyzji. */
function przybornikZmiana(kod: string, autor: StudioAuthor, poczatek: number): StudioTrackedChange {
  return {
    id: kod,
    documentId: 'dokument-1',
    kind: StudioChangeKind.Wstawienie,
    author: autor,
    rangeStart: poczatek,
    rangeEnd: poczatek + 5,
    before: 'było',
    after: 'jest teraz',
    decision: StudioChangeDecision.Oczekuje,
    createdAt: 2_000,
  };
}

/** Adnotacja przy fragmencie różnicy. */
function przybornikAdnotacja(kod: string, numerFragmentu: number): StudioAnnotation {
  return {
    id: kod,
    documentId: 'dokument-1',
    hunkIndex: numerFragmentu,
    author: StudioAuthor.Uzytkownik,
    body: 'uwaga do fragmentu',
    createdAt: 3_000,
  };
}

/* ── Wykaz znakowań ────────────────────────────────────────────────────────── */

describe('wykaz znakowań', () => {
  it('rozdziela pięć rodzajów i NIE zlewa trzech bytów marginesu w jeden', () => {
    const pozycje = przybornikZlozZnakowania({
      komentarze: [przybornikKomentarz('komentarz-1', StudioAuthor.Model, 10, false)],
      zmiany: [przybornikZmiana('zmiana-1', StudioAuthor.Model, 20)],
      adnotacje: [przybornikAdnotacja('adnotacja-1', 3)],
      znaczniki: [],
      propozycja: { idPropozycji: 'propozycja-1', tresc: 'nowe brzmienie', idAkcji: 'studio.styl.ton' },
      zakresPropozycji: { poczatek: 30, koniec: 40 },
    });

    expect(przybornikPolicz(pozycje, RodzajZnakowania.Komentarz)).toBe(1);
    expect(przybornikPolicz(pozycje, RodzajZnakowania.Propozycja)).toBe(1);
    expect(przybornikPolicz(pozycje, RodzajZnakowania.Zmiana)).toBe(1);
    expect(przybornikPolicz(pozycje, RodzajZnakowania.Adnotacja)).toBe(1);

    // Propozycja niesie brzmienie i NIE jest w treści; zmiana śledzona JEST.
    // Rozróżnienie musi zostać widoczne w podstawie pozycji, nie tylko w nazwie.
    const propozycja = pozycje.find((pozycja) => pozycja.rodzaj === RodzajZnakowania.Propozycja);
    const zmiana = pozycje.find((pozycja) => pozycja.rodzaj === RodzajZnakowania.Zmiana);
    expect(propozycja?.tresc).toBe('nowe brzmienie');
    expect(propozycja?.podstawa).toContain('NA MARGINESIE');
    expect(zmiana?.podstawa).toContain('JUŻ');
  });

  it('układa wykaz wedle położenia w treści, a niezakotwiczone na końcu', () => {
    const pozycje = przybornikZlozZnakowania({
      komentarze: [
        przybornikKomentarz('daleki', StudioAuthor.Uzytkownik, 500, false),
        przybornikKomentarz('bliski', StudioAuthor.Uzytkownik, 5, false),
      ],
      zmiany: [],
      adnotacje: [przybornikAdnotacja('adnotacja-1', 0)],
      znaczniki: [],
      propozycja: null,
      zakresPropozycji: null,
    });

    expect(pozycje.map((pozycja) => pozycja.kod)).toEqual(['bliski', 'daleki', 'adnotacja-1']);
    // Adnotacja wisi przy numerze fragmentu różnicy, nie przy znaku treści —
    // udawany zakres wskazywałby niewłaściwe miejsce.
    expect(pozycje[2]?.zakres).toBeNull();
  });

  it('nie liczy odpowiedzi w wątku jako osobnego znakowania', () => {
    const watek = przybornikKomentarz('watek', StudioAuthor.Uzytkownik, 10, false);
    const odpowiedz: StudioComment = {
      ...przybornikKomentarz('odpowiedz', StudioAuthor.Model, 10, false),
      parentCommentId: 'watek',
    };
    const pozycje = przybornikZlozZnakowania({
      komentarze: [watek, odpowiedz],
      zmiany: [],
      adnotacje: [],
      znaczniki: [],
      propozycja: null,
      zakresPropozycji: null,
    });

    expect(pozycje).toHaveLength(1);
    expect(pozycje[0]?.podstawa).toContain('odpowiedzi w wątku: 1');
  });

  it('zawęża wedle rodzaju, autora i stanu — trzy osie osobno i razem', () => {
    const pozycje = przybornikZlozZnakowania({
      komentarze: [
        przybornikKomentarz('operatora-otwarty', StudioAuthor.Uzytkownik, 10, false),
        przybornikKomentarz('modelu-rozwiazany', StudioAuthor.Model, 20, true),
      ],
      zmiany: [przybornikZmiana('zmiana-modelu', StudioAuthor.Model, 30)],
      adnotacje: [],
      znaczniki: [],
      propozycja: null,
      zakresPropozycji: null,
    });

    expect(
      przybornikPrzefiltruj(pozycje, {
        rodzaj: RodzajZnakowania.Komentarz,
        autor: 'wszyscy',
        stan: 'wszystkie',
      }),
    ).toHaveLength(2);
    expect(
      przybornikPrzefiltruj(pozycje, {
        rodzaj: 'wszystkie',
        autor: StudioAuthor.Model,
        stan: 'wszystkie',
      }),
    ).toHaveLength(2);
    expect(
      przybornikPrzefiltruj(pozycje, { rodzaj: 'wszystkie', autor: 'wszyscy', stan: 'zamkniete' }),
    ).toHaveLength(1);
    expect(
      przybornikPrzefiltruj(pozycje, {
        rodzaj: RodzajZnakowania.Komentarz,
        autor: StudioAuthor.Model,
        stan: 'otwarte',
      }),
    ).toHaveLength(0);
  });

  it('podaje liczby policzone, a brak znakowań nazywa zdaniem, nie zerem', () => {
    expect(przybornikOpiszWykaz([], [])).toContain('ani jednego znakowania');

    const pozycje = przybornikZlozZnakowania({
      komentarze: [przybornikKomentarz('komentarz-1', StudioAuthor.Model, 10, false)],
      zmiany: [przybornikZmiana('zmiana-1', StudioAuthor.Uzytkownik, 20)],
      adnotacje: [],
      znaczniki: [],
      propozycja: null,
      zakresPropozycji: null,
    });
    const widoczne = przybornikPrzefiltruj(pozycje, {
      rodzaj: 'wszystkie',
      autor: StudioAuthor.Model,
      stan: 'wszystkie',
    });
    const zdanie = przybornikOpiszWykaz(pozycje, widoczne);
    expect(zdanie).toContain('Znakowań 2');
    expect(zdanie).toContain('autora „model" 1');
    expect(zdanie).toContain('widocznych 1');
  });
});

/* ── Znaczniki własne ──────────────────────────────────────────────────────── */

describe('znaczniki własne Operatora i modelu', () => {
  it('odrzuca znacznik bez nazwy, a barwę spoza palety schodzi na pierwszą', () => {
    const znaczniki = utworzZnacznikiWlasne();
    expect(znaczniki.zaloz('   ', 'sygnal', null, StudioAuthor.Uzytkownik)).toBeNull();

    const znacznik = znaczniki.zaloz('wymaga źródła', 'barwa-której-nie-ma', null, StudioAuthor.Model);
    expect(znacznik?.barwa).toBe('sygnal');
    expect(znacznik?.autor).toBe(StudioAuthor.Model);
  });

  it('wchodzi do wykazu znakowań z autorem i stanem, i daje się odhaczyć', () => {
    const znaczniki = utworzZnacznikiWlasne();
    const znacznik = znaczniki.zaloz('do sprawdzenia', 'ostrzezenie', { poczatek: 4, koniec: 9 }, StudioAuthor.Model);
    if (znacznik === null) throw new Error('znacznik nie został założony');

    znaczniki.przestaw(znacznik.kod, false);
    const pozycje = przybornikZlozZnakowania({
      komentarze: [],
      zmiany: [],
      adnotacje: [],
      znaczniki: znaczniki.wykaz(),
      propozycja: null,
      zakresPropozycji: null,
    });

    expect(pozycje).toHaveLength(1);
    expect(pozycje[0]?.autor).toBe(StudioAuthor.Model);
    expect(pozycje[0]?.otwarta).toBe(false);
    expect(pozycje[0]?.zakres).toEqual({ poczatek: 4, koniec: 9 });

    znaczniki.zdejmij(znacznik.kod);
    expect(znaczniki.wykaz()).toHaveLength(0);
  });
});

/* ── Użycie: kolejność czynności pływaka ───────────────────────────────────── */

describe('kolejność czynności na wierzchu pływaka', () => {
  it('stawia przypięte przed najczęstszymi, a najczęstsze przed rzadkimi', () => {
    const kolejnosc = przybornikNaWierzchu(
      [
        { idAkcji: 'studio.styl.ton', razy: 9, ostatnio: 100 },
        { idAkcji: 'studio.korekta.interpunkcja', razy: 2, ostatnio: 200 },
      ],
      ['studio.streszczenie.tezy'],
      3,
    );
    expect(kolejnosc).toEqual([
      'studio.streszczenie.tezy',
      'studio.styl.ton',
      'studio.korekta.interpunkcja',
    ]);
  });

  it('przy równej liczbie użyć wygrywa użyta ostatnio', () => {
    const kolejnosc = przybornikNaWierzchu(
      [
        { idAkcji: 'studio.styl.ton', razy: 3, ostatnio: 100 },
        { idAkcji: 'studio.styl.skrocenie', razy: 3, ostatnio: 900 },
      ],
      [],
      2,
    );
    expect(kolejnosc[0]).toBe('studio.styl.skrocenie');
  });

  it('bez ani jednego użycia dopełnia wykazem początkowym, nie pustką', () => {
    const kolejnosc = przybornikNaWierzchu([], [], 4);
    expect(kolejnosc).toHaveLength(4);
    expect(new Set(kolejnosc).size).toBe(4);
  });

  it('odczytuje nastawy z rdzenia i pomija wpisy o kształcie niezgodnym', () => {
    const uzycie = przybornikOdczytajUzycie(
      JSON.stringify([
        { idAkcji: 'studio.styl.ton', razy: 4, ostatnio: 7 },
        { idAkcji: 'bez-liczby' },
        { razy: 3 },
        'napis',
      ]),
    );
    expect(uzycie).toEqual([{ idAkcji: 'studio.styl.ton', razy: 4, ostatnio: 7 }]);

    expect(przybornikOdczytajPrzypiete(JSON.stringify(['studio.styl.ton', 7]))).toEqual([
      'studio.styl.ton',
    ]);
  });

  it('brak nastawy trybu znaczy NARZĘDZIA UKRYTE, a nie stały panel', () => {
    expect(przybornikOdczytajTryb(undefined)).toBe(TrybOperacji.Ukryte);
    expect(przybornikOdczytajTryb('cokolwiek')).toBe(TrybOperacji.Ukryte);
    expect(przybornikOdczytajTryb(JSON.stringify(TrybOperacji.Panel))).toBe(TrybOperacji.Panel);
  });
});

/* ── Katalog operacji ──────────────────────────────────────────────────────── */

describe('katalog operacji jako jeden wykaz', () => {
  /** Wiersz rejestru akcji wskazujący operację kontekstową. */
  function przybornikAkcja(kod: string, komenda: Command): Action {
    return {
      id: kod,
      name: `akcja ${kod}`,
      command: komenda,
      scope: ConfigScope.Module,
      enabled: true,
    } as Action;
  }

  /** Operacja zapisana w rdzeniu — fabryczna albo własna. */
  function przybornikOperacja(kod: string, fabryczna: boolean): StudioOperation {
    return {
      id: kod,
      name: `operacja ${kod}`,
      category: 'Korekta i jakość',
      prompt: 'polecenie',
      builtin: fabryczna,
      scope: ConfigScope.Global,
    };
  }

  it('zdejmuje z wykazu dokumentacji pozycję przejętą przez rdzeń', () => {
    const przejeta = WSZYSTKIE_OPERACJE[0];
    if (przejeta === undefined) throw new Error('wykaz operacji jest pusty');

    const bezRdzenia = przybornikZlozDrzewo([], []);
    const zRdzeniem = przybornikZlozDrzewo(
      [przybornikAkcja(przejeta.id, Command.StudioContextualOp)],
      [],
    );

    expect(przybornikPoliczLiscie(bezRdzenia)).toContain(przejeta.id);
    // Ta sama operacja nie może stać w menu dwa razy: raz z rejestru, raz
    // z wykazu dokumentacji.
    expect(przybornikPoliczLiscie(zRdzeniem).filter((klucz) => klucz === przejeta.id)).toHaveLength(1);
  });

  it('rozdziela operacje własne od fabrycznych i daje usunięcie tylko własnym', () => {
    const drzewo = przybornikZlozDrzewo(
      [],
      [przybornikOperacja('wlasna-1', false), przybornikOperacja('fabryczna-1', true)],
    );
    const liscie = przybornikPoliczLiscie(drzewo);
    expect(liscie).toContain('wlasna-1');
    expect(liscie).toContain('fabryczna-1');
    expect(liscie).toContain('usun:wlasna-1');
    expect(liscie).not.toContain('usun:fabryczna-1');
  });

  it('pomija w menu wiersze rejestru wskazujące inną komendę', () => {
    const drzewo = przybornikZlozDrzewo(
      [przybornikAkcja('studio.document.save', Command.StudioDocumentSave)],
      [],
    );
    expect(przybornikPoliczLiscie(drzewo)).not.toContain('studio.document.save');
  });

  it('podpowiada czynności po frazie, a fraza pusta oddaje wykaz w całości', () => {
    expect(operacjePoFrazie('')).toHaveLength(WSZYSTKIE_OPERACJE.length);
    expect(operacjePoFrazie('streszczenie').length).toBeGreaterThan(0);
    expect(operacjePoFrazie('czegoś takiego nie ma')).toHaveLength(0);
  });
});

/** Klucze wszystkich liści drzewa menu, na dowolnej głębokości. */
function przybornikPoliczLiscie(drzewo: readonly unknown[]): string[] {
  const klucze: string[] = [];
  for (const pozycja of drzewo) {
    if (typeof pozycja !== 'object' || pozycja === null) continue;
    const pola = pozycja as { rodzaj?: string; klucz?: string; dzieci?: readonly unknown[] };
    if (pola.rodzaj === 'wybor' && typeof pola.klucz === 'string') klucze.push(pola.klucz);
    if (pola.dzieci !== undefined) klucze.push(...przybornikPoliczLiscie(pola.dzieci));
  }
  return klucze;
}

/* ── Pochodzenie wniesionych fragmentów ────────────────────────────────────── */

describe('pochodzenie wniesionego fragmentu', () => {
  it('zapisuje plik Biblioteki wraz z wersją i sumą kontrolną', () => {
    const wykaz = utworzWykazPochodzen();
    const zapis = wykaz.zapisz({
      zrodlo: ZrodloWniesienia.Biblioteka,
      nazwa: 'wzór pisma.docx',
      wskazanie: 'plik-7',
      wersja: 'wersja-2',
      sumaKontrolna: 'suma-abc',
      strona: 3,
      znakow: 120,
      wstawionoNaZnaku: 44,
    });

    const wiersz = osadzenieWierszPochodzenia(zapis);
    expect(wiersz).toContain('wzór pisma.docx');
    expect(wiersz).toContain('plik-7');
    expect(wiersz).toContain('wersja-2');
    expect(wiersz).toContain('suma-abc');
    expect(wiersz).toContain('strona 3');
    // Zakaz numeracji wymyślonej obowiązuje także tutaj: wiersz jest zdaniem,
    // nie kodem odsyłacza.
    expect(wiersz).not.toContain('§');
    expect(wiersz).not.toMatch(/\[\d+\]/u);
  });

  it('pomija pola nieznane, zamiast wpisywać puste wskazanie', () => {
    const wykaz = utworzWykazPochodzen();
    const zapis = wykaz.zapisz({
      zrodlo: ZrodloWniesienia.Strona,
      nazwa: 'Dziennik Ustaw',
      wskazanie: 'https://przyklad.test/akt',
      wersja: '',
      sumaKontrolna: '',
      strona: 0,
      znakow: 80,
      wstawionoNaZnaku: 0,
    });

    const wiersz = osadzenieWierszPochodzenia(zapis);
    expect(wiersz).toContain('studio.ingest.url');
    expect(wiersz).not.toContain('wersja ,');
    expect(wiersz).not.toContain('suma kontrolna ,');
  });

  it('brak wniesień nazywa zdaniem, a wniesienia liczy', () => {
    const wykaz = utworzWykazPochodzen();
    expect(osadzenieOpiszPochodzenia(wykaz.wykaz())).toContain('ani jednego fragmentu');

    wykaz.zapisz({
      zrodlo: ZrodloWniesienia.Biblioteka,
      nazwa: 'plik',
      wskazanie: 'plik-1',
      wersja: '',
      sumaKontrolna: '',
      strona: 0,
      znakow: 50,
      wstawionoNaZnaku: 0,
    });
    wykaz.zapisz({
      zrodlo: ZrodloWniesienia.Strona,
      nazwa: 'strona',
      wskazanie: 'https://przyklad.test',
      wersja: '',
      sumaKontrolna: '',
      strona: 0,
      znakow: 70,
      wstawionoNaZnaku: 50,
    });

    const zdanie = osadzenieOpiszPochodzenia(wykaz.wykaz());
    expect(zdanie).toContain('wniesionych 2');
    expect(zdanie).toContain('z Biblioteki 1');
    expect(zdanie).toContain('120 znaków');
  });
});
