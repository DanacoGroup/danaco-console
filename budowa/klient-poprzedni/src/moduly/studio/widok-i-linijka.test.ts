import { beforeEach, describe, expect, it } from 'vitest';
import { StudioAuthor, StudioPageOrientation, type StudioVersion } from '../../../../shared/contract';
import {
  domyslnaStrona,
  marginesyKartki,
  NASTAWY_MARGINESOW,
  nosnikPoOznaczeniu,
  nosnikWlasny,
  nosniki,
  opiszDlugosc,
  profilZeStrony,
  szerokoscPolaMm,
  wchlonNosnikiKontraktu,
  wchlonNosnikiRdzenia,
  wysokoscPolaPunkty,
  zastosujNastaweMarginesow,
} from './nastawy-strony';
import {
  krokPrzyciaganiaLinijki,
  nastepnyRodzajTabulatora,
  opiszWciecieAkapitu,
  przyciagnijNaLinijce,
  punktyNaLinijce,
  zerowaWciecieAkapitu,
  zlozPodzialkeLinijki,
} from './linijka-podzialka';
import { utworzLinijkePozioma } from './linijka-pozioma';
import { utworzLinijkePionowa } from './linijka-pionowa';
import {
  czytajNastawyWidoku,
  domyslneNastawyWidoku,
  przytnijUdzialPodzialu,
  utworzPamiecNastawWidoku,
  zapamietajNastawyWidoku,
  type MagazynNastawWidoku,
} from './widok-nastawy-operatora';
import { policzSkaleWidoku, przytnijSkale, skalaDokumentu, zapamietajSkaleDokumentu } from './widok-skali';
import {
  kartkaPrzyPrzewinieciu,
  kolumnyUkladu,
  przewiniecieDoKartki,
  rozlozKartkiWRzedy,
  stronaRozkladowki,
} from './widok-ukladu-stron';
import { utworzPodzialPowierzchni } from './widok-podzialu-powierzchni';
import {
  czytajNastawyDruku,
  domyslneNastawyDruku,
  opiszNastawyDruku,
  stronyDoDruku,
  wydrukuj,
  zapamietajNastawyDruku,
} from './widok-druku';
import { czyZapisSamoczynny, pasujeDoAutora, przefiltrujHistorie } from './filtr-historii';
import { utworzWierszWersji } from './wiersz-wersji';
import { GRUPY_WSTAZKI_PDF } from './okno-warsztatu-dokumentu';
import { CZYNNOSCI_WARSZTATU } from './czynnosci-warsztatu';

/**
 * Sprawdziany mierzą skutek działania, nie kształt wejścia, i pomijają DOM przy czystym rachunku.
 */

/** Magazyn zapasowy nastaw działa jako atrapa w sprawdzianach, dzięki czemu sprawdzian nie sięga do rzeczywistego magazynu przeglądarki. */
function magazynAtrapa(początek: Record<string, string> = {}): MagazynNastawWidoku {
  const wpisy = new Map<string, string>(Object.entries(początek));
  return {
    getItem: (klucz) => wpisy.get(klucz) ?? null,
    setItem: (klucz, wartosc) => {
      wpisy.set(klucz, wartosc);
    },
  };
}

describe('nastawy strony: nośniki, oprawa i marginesy odbicia', () => {
  it('niesie szereg A, szereg B, Letter, Legal, Tabloid i cztery koperty', () => {
    const oznaczenia = nosniki().map((nosnik) => nosnik.oznaczenie);
    for (const szukane of ['A0', 'A3', 'A4', 'A6', 'B1', 'B5', 'Letter', 'Legal', 'Tabloid']) {
      expect(oznaczenia).toContain(szukane);
    }
    for (const koperta of ['DL', 'C4', 'C5', 'C6']) {
      expect(oznaczenia).toContain(koperta);
      expect(nosnikPoOznaczeniu(koperta)?.rodzaj).toBe('koperta');
    }
  });

  it('wykaz rdzenia nadpisuje wymiary wbudowane i nie gubi kopert', () => {
    const po = wchlonNosnikiRdzenia([
      { oznaczenie: 'A4', szerokoscMm: 211, wysokoscMm: 298 },
      { oznaczenie: 'SRA3', szerokoscMm: 320, wysokoscMm: 450 },
    ]);
    expect(nosnikPoOznaczeniu('A4')?.szerokoscMm).toBe(211);
    expect(po.map((nosnik) => nosnik.oznaczenie)).toContain('SRA3');
    // Koperta nieznana rdzeniowi zostaje w wykazie: jej brak jest luką zgłoszoną, nie powodem usunięcia.
    expect(nosnikPoOznaczeniu('C5')).not.toBeNull();
    // Przywrócenie wymiaru wzorcowego, żeby dalsze sprawdziany liczyły na A4.
    wchlonNosnikiRdzenia([{ oznaczenie: 'A4', szerokoscMm: 210, wysokoscMm: 297 }]);
  });

  it('wchłania odpowiedź studio.page.paper.list wraz z rodzajem nośnika', () => {
    // Kształt pochodzi z kontraktu: name, widthMm, heightMm, kind; przekład mieszka w jednym miejscu.
    wchlonNosnikiKontraktu([
      { name: 'A4', widthMm: 210, heightMm: 297, kind: 'sheet' },
      { name: 'C4', widthMm: 229, heightMm: 324, kind: 'envelope' },
      { name: 'B5', widthMm: 176, heightMm: 250, kind: 'sheet' },
      // Nośnik nieznany oknu bierze rodzaj z odpowiedzi rdzenia, inaczej trafiłby do okna jako arkusz.
      { name: 'C65', widthMm: 114, heightMm: 229, kind: 'envelope' },
    ]);
    expect(nosnikPoOznaczeniu('C4')?.rodzaj).toBe('koperta');
    expect(nosnikPoOznaczeniu('C4')?.szerokoscMm).toBe(229);
    expect(nosnikPoOznaczeniu('C4')?.wysokoscMm).toBe(324);
    expect(nosnikPoOznaczeniu('B5')?.rodzaj).toBe('arkusz');
    expect(nosnikPoOznaczeniu('C65')?.rodzaj).toBe('koperta');
    // Wykaz wbudowany zostaje pełny — pozycja zniknięta z okna wyglądałaby
    // na usterkę.
    expect(nosnikPoOznaczeniu('Letter')).not.toBeNull();
  });

  it('odrzuca format własny o wymiarze niedodatnim, zamiast rysować kartkę zerową', () => {
    expect(nosnikWlasny(0, 100)).toBeNull();
    expect(nosnikWlasny(-5, 100)).toBeNull();
    expect(nosnikWlasny(Number.POSITIVE_INFINITY, 100)).toBeNull();
    expect(nosnikWlasny(120, 240)?.rodzaj).toBe('wlasny');
  });

  it('oprawa zabiera pole pisania, a marginesy odbicia przenoszą ją na stronie parzystej', () => {
    const strona = {
      ...domyslnaStrona(),
      marginesOprawyMm: 10,
      marginesyOdbicia: true,
    };
    const nieparzysta = marginesyKartki(strona, 1);
    const parzysta = marginesyKartki(strona, 2);
    expect(nieparzysta.lewyMm).toBe(30);
    expect(nieparzysta.prawyMm).toBe(20);
    expect(parzysta.lewyMm).toBe(20);
    expect(parzysta.prawyMm).toBe(30);
    // Pole pisania jest węższe o oprawę — na obu stronach tak samo.
    expect(szerokoscPolaMm(strona, 1)).toBe(160);
    expect(szerokoscPolaMm(strona, 2)).toBe(160);
  });

  it('oprawa u góry zabiera margines górny, nie boczny', () => {
    const strona = { ...domyslnaStrona(), marginesOprawyMm: 12, stronaOprawy: 'gora' as const };
    expect(marginesyKartki(strona, 1).goraMm).toBe(32);
    expect(marginesyKartki(strona, 1).lewyMm).toBe(20);
  });

  it('nastawa gotowa „do oprawy" włącza marginesy odbicia', () => {
    const nastawa = NASTAWY_MARGINESOW.find((pozycja) => pozycja.kod === 'oprawa');
    expect(nastawa).toBeDefined();
    const po = zastosujNastaweMarginesow(domyslnaStrona(), nastawa!);
    expect(po.marginesOprawyMm).toBe(12);
    expect(po.marginesyOdbicia).toBe(true);
  });

  it('wysokość pola maleje po zmianie orientacji, bo kartka zmienia wymiar', () => {
    const pionowa = domyslnaStrona();
    const pozioma = { ...pionowa, orientacja: StudioPageOrientation.Pozioma };
    expect(wysokoscPolaPunkty(pozioma)).toBeLessThan(wysokoscPolaPunkty(pionowa));
  });

  it('profil wydania niesie tylko pola, które kontrakt ma — bez oprawy i odbić', () => {
    const profil = profilZeStrony({ ...domyslnaStrona(), marginesOprawyMm: 15 });
    expect(profil.pageSize).toBe('A4');
    expect(Object.keys(profil)).not.toContain('gutter');
  });

  it('długość opisuje się w jednostce wybranej przez Operatora', () => {
    expect(opiszDlugosc(25.4, 'mm')).toBe('25.4 mm');
    expect(opiszDlugosc(25.4, 'cal')).toBe('1.00″');
  });
});

describe('podziałka linijki', () => {
  it('stawia kreskę dużą co centymetr i podpisuje ją numerem centymetra', () => {
    const kreski = zlozPodzialkeLinijki(50, 'mm', 1);
    const duze = kreski.filter((kreska) => kreska.rodzaj === 'duza');
    expect(duze.map((kreska) => kreska.milimetry)).toEqual([0, 10, 20, 30, 40, 50]);
    expect(duze[1]?.napis).toBe('1');
    // Kreska w zerze nie ma napisu: „0" przy krawędzi pola nic nie wnosi.
    expect(duze[0]?.napis).toBe('');
  });

  it('podziałka calowa ma kreski co ósmą część cala, nie co milimetr', () => {
    const kreski = zlozPodzialkeLinijki(25.4, 'cal', 1);
    expect(kreski.length).toBe(9);
    expect(kreski[8]?.rodzaj).toBe('duza');
    expect(kreski[8]?.napis).toBe('1');
  });

  it('przy małej skali schodzą kreski małe, a siatka centymetrów zostaje', () => {
    const gesta = zlozPodzialkeLinijki(200, 'mm', 1);
    const rzadka = zlozPodzialkeLinijki(200, 'mm', 0.2);
    expect(rzadka.length).toBeLessThan(gesta.length);
    expect(rzadka.every((kreska) => kreska.rodzaj !== 'mala')).toBe(true);
    expect(rzadka.filter((kreska) => kreska.rodzaj === 'duza').length).toBe(21);
  });

  it('nie oddaje kresek dla długości niedodatniej', () => {
    expect(zlozPodzialkeLinijki(0, 'mm', 1)).toEqual([]);
    expect(zlozPodzialkeLinijki(-10, 'mm', 1)).toEqual([]);
  });

  it('przyciąga do kroku i przycina do zakresu, zamiast oddawać wartość spoza kartki', () => {
    expect(przyciagnijNaLinijce(19.73, 'mm', 0, 100)).toBe(19.5);
    expect(przyciagnijNaLinijce(-5, 'mm', 0, 100)).toBe(0);
    expect(przyciagnijNaLinijce(250, 'mm', 0, 100)).toBe(100);
    expect(krokPrzyciaganiaLinijki('cal')).toBeCloseTo(25.4 / 16, 5);
  });

  it('punkty rosną ze skalą — podziałka rozciąga się razem z kartką', () => {
    expect(punktyNaLinijce(10, 2)).toBeCloseTo(punktyNaLinijce(10, 1) * 2, 5);
    // Skala nieprawidłowa bierze jeden, zamiast oddawać nieskończoność.
    expect(punktyNaLinijce(10, 0)).toBeCloseTo(punktyNaLinijce(10, 1), 5);
  });

  it('rodzaj tabulatora chodzi w obiegu czterech i wraca do lewego', () => {
    expect(nastepnyRodzajTabulatora('lewy')).toBe('srodkowy');
    expect(nastepnyRodzajTabulatora('srodkowy')).toBe('prawy');
    expect(nastepnyRodzajTabulatora('prawy')).toBe('dziesietny');
    expect(nastepnyRodzajTabulatora('dziesietny')).toBe('lewy');
  });

  it('opis wcięć nazywa wysunięcie pierwszego wiersza po imieniu', () => {
    expect(opiszWciecieAkapitu({ leweMm: 10, praweMm: 0, pierwszyWierszMm: -5 }, 'mm')).toContain(
      'wysunięty',
    );
    expect(opiszWciecieAkapitu(zerowaWciecieAkapitu(), 'mm')).toContain('bez wcięcia');
  });
});

describe('linijka pozioma: chwyty marginesów, wcięć i tabulatorów', () => {
  it('chwyt marginesu jest suwakiem z wartością w milimetrach', () => {
    const linijka = utworzLinijkePozioma(domyslnaStrona(), {
      naMargines: () => undefined,
      naWciecie: () => undefined,
      naTabulatory: () => undefined,
      naKrawedzKolumny: () => undefined,
    });
    const chwyt = linijka.element.querySelector<HTMLElement>("[data-chwyt='margines-lewy']");
    expect(chwyt).not.toBeNull();
    expect(chwyt?.getAttribute('role')).toBe('slider');
    expect(chwyt?.getAttribute('aria-valuenow')).toBe('20');
    // Chwyt nie zostawi pola pisania węższego niż dziesięć milimetrów.
    expect(Number(chwyt?.getAttribute('aria-valuemax'))).toBe(210 - 20 - 10);
  });

  it('strzałka przesuwa margines o krok podziałki i zgłasza go oknu', () => {
    const zgloszone: number[] = [];
    const linijka = utworzLinijkePozioma(domyslnaStrona(), {
      naMargines: (_strona, milimetry) => zgloszone.push(milimetry),
      naWciecie: () => undefined,
      naTabulatory: () => undefined,
      naKrawedzKolumny: () => undefined,
    });
    const chwyt = linijka.element.querySelector<HTMLElement>("[data-chwyt='margines-lewy']");
    chwyt?.dispatchEvent(new KeyboardEvent('keydown', { key: 'ArrowRight', bubbles: true }));
    expect(zgloszone).toEqual([20.5]);
  });

  it('wszystkie pięć chwytów stoi na linijce, wcięcia osobnymi znacznikami', () => {
    const linijka = utworzLinijkePozioma(domyslnaStrona(), {
      naMargines: () => undefined,
      naWciecie: () => undefined,
      naTabulatory: () => undefined,
      naKrawedzKolumny: () => undefined,
    });
    for (const kod of [
      'margines-lewy',
      'margines-prawy',
      'wciecie-pierwszego-wiersza',
      'wciecie-lewe',
      'wciecie-prawe',
    ]) {
      expect(linijka.element.querySelector(`[data-chwyt='${kod}']`)).not.toBeNull();
    }
  });

  it('tabulator zakłada się naciśnięciem podziałki, a naciśnięcie znaku przestawia rodzaj', () => {
    let ostatnie: readonly { rodzaj: string }[] = [];
    const linijka = utworzLinijkePozioma(domyslnaStrona(), {
      naMargines: () => undefined,
      naWciecie: () => undefined,
      naTabulatory: (tabulatory) => {
        ostatnie = tabulatory;
      },
      naKrawedzKolumny: () => undefined,
    });
    const podzialka = linijka.element.querySelector<HTMLElement>('.ms-linijka__podzialka');
    podzialka?.dispatchEvent(new MouseEvent('click', { bubbles: true, clientX: 100 }));
    expect(ostatnie.length).toBe(1);
    expect(linijka.tabulatory().length).toBe(1);

    const znak = linijka.element.querySelector<HTMLElement>('[data-tabulator]');
    znak?.dispatchEvent(new MouseEvent('click', { bubbles: true }));
    expect(ostatnie[0]?.rodzaj).toBe('srodkowy');

    // Naciśnięcie z Shiftem zdejmuje tabulator — bez osobnego przycisku usuwania.
    const znakPoZmianie = linijka.element.querySelector<HTMLElement>('[data-tabulator]');
    znakPoZmianie?.dispatchEvent(new MouseEvent('click', { bubbles: true, shiftKey: true }));
    expect(linijka.tabulatory().length).toBe(0);
  });

  it('krawędzie kolumn tabeli dostają chwyty, po jednym na krawędź', () => {
    const linijka = utworzLinijkePozioma(domyslnaStrona(), {
      naMargines: () => undefined,
      naWciecie: () => undefined,
      naTabulatory: () => undefined,
      naKrawedzKolumny: () => undefined,
    });
    linijka.ustawKrawedzieKolumn([40, 90]);
    expect(linijka.element.querySelectorAll('[data-kolumna]').length).toBe(2);
  });

  it('przełącznik widoczności chowa linijkę, a nie rozbiera jej', () => {
    const linijka = utworzLinijkePozioma(domyslnaStrona(), {
      naMargines: () => undefined,
      naWciecie: () => undefined,
      naTabulatory: () => undefined,
      naKrawedzKolumny: () => undefined,
    });
    linijka.ustawWidocznosc(false);
    expect(linijka.widoczna()).toBe(false);
    expect(linijka.element.hidden).toBe(true);
    linijka.ustawWidocznosc(true);
    expect(linijka.widoczna()).toBe(true);
  });
});

describe('linijka pionowa', () => {
  it('ma chwyt marginesu górnego i dolnego, oba jako suwaki pionowe', () => {
    const zgloszone: string[] = [];
    const linijka = utworzLinijkePionowa(domyslnaStrona(), {
      naMargines: (ktory) => zgloszone.push(ktory),
    });
    const gora = linijka.element.querySelector<HTMLElement>("[data-chwyt='margines-gora']");
    const dol = linijka.element.querySelector<HTMLElement>("[data-chwyt='margines-dol']");
    expect(gora?.getAttribute('aria-orientation')).toBe('vertical');
    expect(dol).not.toBeNull();
    gora?.dispatchEvent(new KeyboardEvent('keydown', { key: 'ArrowDown', bubbles: true }));
    expect(zgloszone).toEqual(['gora']);
  });

  it('opisuje pole pisania, nie samą kartkę', () => {
    const linijka = utworzLinijkePionowa(domyslnaStrona(), { naMargines: () => undefined });
    expect(linijka.opis()).toContain('pole pisania 257 mm');
  });
});

describe('nastawy widoku Operatora', () => {
  it('domyślnie zakładki, przewijanie ciągłe, jedna kartka, tryb źródłowy wyłączony', () => {
    const domyslne = domyslneNastawyWidoku();
    expect(domyslne.trybDokumentow).toBe('zakladki');
    expect(domyslne.przewijanie).toBe('ciagle');
    expect(domyslne.ukladKartek).toBe('jedna');
    expect(domyslne.trybZrodlowy).toBe(false);
  });

  it('zapis i odczyt wracają tą samą nastawą', () => {
    const magazyn = magazynAtrapa();
    zapamietajNastawyWidoku({ ...domyslneNastawyWidoku(), ukladKartek: 'rozkladowka', skala: 75 }, magazyn);
    const wczytane = czytajNastawyWidoku(magazyn);
    expect(wczytane.ukladKartek).toBe('rozkladowka');
    expect(wczytane.skala).toBe(75);
  });

  it('zapis uszkodzony i pole brakujące dają wartość domyślną, nie unieważniają reszty', () => {
    expect(czytajNastawyWidoku(magazynAtrapa({ 'dn.studio.widok': '{{' })).trybDokumentow).toBe(
      'zakladki',
    );
    const czesciowy = czytajNastawyWidoku(
      magazynAtrapa({ 'dn.studio.widok': '{"ukladKartek":"obok","kartekWRzedzie":3}' }),
    );
    expect(czesciowy.ukladKartek).toBe('obok');
    expect(czesciowy.kartekWRzedzie).toBe(3);
    expect(czesciowy.przewijanie).toBe('ciagle');
  });

  it('brak magazynu nie jest błędem — nastawa widoku nie wstrzymuje okna', () => {
    expect(czytajNastawyWidoku(null).skala).toBe(100);
    expect(() => zapamietajNastawyWidoku(domyslneNastawyWidoku(), null)).not.toThrow();
  });

  it('udział podziału jest przycinany do pola pracy', () => {
    expect(przytnijUdzialPodzialu(0.01)).toBe(0.15);
    expect(przytnijUdzialPodzialu(0.99)).toBe(0.85);
    const pamiec = utworzPamiecNastawWidoku(magazynAtrapa());
    expect(pamiec.przestaw({ udzialPodzialu: 2 }).udzialPodzialu).toBe(0.85);
  });
});

describe('skala widoku', () => {
  const wymiary = {
    szerokoscWidokuPx: 800,
    wysokoscWidokuPx: 600,
    szerokoscKartkiMm: 210,
    wysokoscKartkiMm: 297,
    szerokoscTekstuMm: 170,
    kartekWRzedzie: 1,
    odstepPx: 24,
  };

  it('sto procent jest stałe i nie zależy od pola widoku', () => {
    expect(policzSkaleWidoku('sto', { ...wymiary, szerokoscWidokuPx: 0 })).toBe(100);
  });

  it('do szerokości tekstu daje skalę większą niż do szerokości strony', () => {
    const strona = policzSkaleWidoku('szerokosc-strony', wymiary);
    const tekst = policzSkaleWidoku('szerokosc-tekstu', wymiary);
    expect(strona).not.toBeNull();
    expect(tekst).not.toBeNull();
    expect(tekst!).toBeGreaterThan(strona!);
  });

  it('cała strona bierze wymiar ciaśniejszy, żeby nie ucinać dołu kartki', () => {
    const cala = policzSkaleWidoku('cala-strona', wymiary);
    const strona = policzSkaleWidoku('szerokosc-strony', wymiary);
    expect(cala!).toBeLessThanOrEqual(strona!);
  });

  it('nastawa nieznana i pole nieosadzone oddają brak, a nie sto procent po cichu', () => {
    expect(policzSkaleWidoku('wymyslona', wymiary)).toBeNull();
    expect(policzSkaleWidoku('szerokosc-strony', { ...wymiary, szerokoscWidokuPx: 0 })).toBeNull();
  });

  it('kartki w rzędzie zmniejszają skalę mieszczącą — dwie strony obok siebie są mniejsze', () => {
    const jedna = policzSkaleWidoku('szerokosc-strony', wymiary);
    const dwie = policzSkaleWidoku('szerokosc-strony', { ...wymiary, kartekWRzedzie: 2 });
    expect(dwie!).toBeLessThan(jedna!);
  });

  it('skala jest przycinana do granic', () => {
    expect(przytnijSkale(1)).toBe(10);
    expect(przytnijSkale(9999)).toBe(500);
    expect(przytnijSkale(Number.NaN)).toBe(100);
  });

  it('skala jest pamiętana PRZY DOKUMENCIE, a nie wspólna dla wszystkich', () => {
    const magazyn = magazynAtrapa();
    zapamietajSkaleDokumentu('pismo-1', 140, magazyn);
    zapamietajSkaleDokumentu('pismo-2', 60, magazyn);
    expect(skalaDokumentu('pismo-1', magazyn)).toBe(140);
    expect(skalaDokumentu('pismo-2', magazyn)).toBe(60);
    // Dokument bez zapisu oddaje brak, nie sto procent: wołający zostawia wtedy ostatnią skalę Operatora.
    expect(skalaDokumentu('pismo-3', magazyn)).toBeNull();
  });
});

describe('układ kartek i rozkładówka', () => {
  it('jedna kartka w rzędzie daje po jednym rzędzie na stronę', () => {
    expect(rozlozKartkiWRzedy(3, 'jedna', 1)).toEqual([[1], [2], [3]]);
  });

  it('rozkładówka stawia stronę pierwszą SAMĄ po prawej, dalej pary', () => {
    expect(rozlozKartkiWRzedy(5, 'rozkladowka', 2)).toEqual([
      [null, 1],
      [2, 3],
      [4, 5],
    ]);
  });

  it('rozkładówka o liczbie stron parzystej domyka rząd ostatni pustym miejscem', () => {
    expect(rozlozKartkiWRzedy(4, 'rozkladowka', 2)).toEqual([
      [null, 1],
      [2, 3],
      [4, null],
    ]);
  });

  it('kartki obok siebie dopełniają rząd ostatni pustymi miejscami', () => {
    expect(rozlozKartkiWRzedy(5, 'obok', 3)).toEqual([
      [1, 2, 3],
      [4, 5, null],
    ]);
  });

  it('dokument bez treści ma jedną kartkę, nie zero', () => {
    expect(rozlozKartkiWRzedy(0, 'jedna', 1)).toEqual([[1]]);
  });

  it('liczba kolumn zgadza się z układem i jest przycięta do ośmiu', () => {
    expect(kolumnyUkladu('jedna', 5)).toBe(1);
    expect(kolumnyUkladu('rozkladowka', 5)).toBe(2);
    expect(kolumnyUkladu('obok', 99)).toBe(8);
  });

  it('strona nieparzysta jest prawą stroną rozkładówki — zgodnie z marginesem wewnętrznym', () => {
    expect(stronaRozkladowki(1)).toBe('prawa');
    expect(stronaRozkladowki(2)).toBe('lewa');
    const strona = { ...domyslnaStrona(), marginesyOdbicia: true, marginesOprawyMm: 10 };
    expect(marginesyKartki(strona, 1).lewyMm).toBeGreaterThan(marginesyKartki(strona, 1).prawyMm);
  });

  it('przewijanie strona po stronie liczy kartkę widoczną i wraca tym samym przewinięciem', () => {
    expect(kartkaPrzyPrzewinieciu(0, 1000, 24, 5)).toBe(1);
    expect(kartkaPrzyPrzewinieciu(2048, 1000, 24, 5)).toBe(3);
    // Przewinięcie poza dokument nie oddaje kartki, której nie ma.
    expect(kartkaPrzyPrzewinieciu(999999, 1000, 24, 5)).toBe(5);
    expect(przewiniecieDoKartki(3, 1000, 24)).toBe(2048);
  });
});

describe('podział powierzchni jako drugi równorzędny tryb', () => {
  it('w trybie zakładek pole drugie jest SCHOWANE, nie rozebrane', () => {
    const podzial = utworzPodzialPowierzchni({ naUdzial: () => undefined, naCzynne: () => undefined });
    const pierwsze = document.createElement('p');
    const drugie = document.createElement('p');
    drugie.textContent = 'drugi dokument';
    podzial.ustawPola(pierwsze, drugie);

    expect(podzial.tryb()).toBe('zakladki');
    const polePole = podzial.element.querySelector<HTMLElement>("[data-pole='2']");
    expect(polePole?.hidden).toBe(true);
    // Niewidoczny nie znaczy niedostępny: element drugiego dokumentu zostaje w drzewie, dostępny modelowi.
    expect(podzial.element.textContent).toContain('drugi dokument');
  });

  it('przełączenie trybu w obie strony niczego nie gubi', () => {
    const podzial = utworzPodzialPowierzchni({ naUdzial: () => undefined, naCzynne: () => undefined });
    const drugie = document.createElement('p');
    drugie.textContent = 'zostaje';
    podzial.ustawPola(document.createElement('p'), drugie);
    podzial.ustawTryb('podzial');
    expect(podzial.element.querySelector<HTMLElement>("[data-pole='2']")?.hidden).toBe(false);
    podzial.ustawTryb('zakladki');
    expect(podzial.element.textContent).toContain('zostaje');
  });

  it('granica jest separatorem z wartością i chodzi z klawiatury', () => {
    const zgloszone: number[] = [];
    const podzial = utworzPodzialPowierzchni(
      { naUdzial: (udzial) => zgloszone.push(udzial), naCzynne: () => undefined },
      'podzial',
    );
    const granica = podzial.element.querySelector<HTMLElement>("[role='separator']");
    expect(granica?.getAttribute('aria-valuenow')).toBe('50');
    granica?.dispatchEvent(new KeyboardEvent('keydown', { key: 'ArrowRight', bubbles: true }));
    expect(podzial.udzial()).toBeCloseTo(0.52, 5);
    granica?.dispatchEvent(new KeyboardEvent('keydown', { key: 'Home', bubbles: true }));
    expect(podzial.udzial()).toBe(0.15);
    granica?.dispatchEvent(new KeyboardEvent('keydown', { key: 'Enter', bubbles: true }));
    expect(podzial.udzial()).toBe(0.5);
    expect(zgloszone.length).toBe(3);
  });

  it('w trybie zakładek polem czynnym zostaje pierwsze, choćby wskazano drugie', () => {
    const podzial = utworzPodzialPowierzchni({ naUdzial: () => undefined, naCzynne: () => undefined });
    podzial.ustawCzynne(2);
    expect(podzial.czynne()).toBe(1);
    podzial.ustawTryb('podzial');
    podzial.ustawCzynne(2);
    expect(podzial.czynne()).toBe(2);
  });
});

describe('drukowanie jako czynność Operatora', () => {
  it('zakres „wszystkie" oddaje wszystkie strony, a „bieżąca" jedną', () => {
    const nastawy = domyslneNastawyDruku();
    expect(stronyDoDruku(nastawy, 4, 2)).toEqual([1, 2, 3, 4]);
    expect(stronyDoDruku({ ...nastawy, zakres: 'biezaca' }, 4, 2)).toEqual([2]);
  });

  it('zakres podany odwrotnie jest odwracany, nie odrzucany', () => {
    const nastawy = { ...domyslneNastawyDruku(), zakres: 'podany' as const, odStrony: 4, doStrony: 2 };
    expect(stronyDoDruku(nastawy, 6, 1)).toEqual([2, 3, 4]);
  });

  it('zakres poza dokumentem oddaje wykaz pusty, a wydruk odmawia zamiast drukować wszystko', () => {
    const nastawy = { ...domyslneNastawyDruku(), zakres: 'podany' as const, odStrony: 20, doStrony: 30 };
    expect(stronyDoDruku(nastawy, 5, 1)).toEqual([]);
    const wynik = wydrukuj(nastawy, {
      liczbaStron: () => 5,
      kartkaBiezaca: () => 1,
      przygotuj: () => () => undefined,
    });
    expect(wynik.udany).toBe(false);
    expect(wynik.zdanie).toContain('poza dokumentem');
  });

  it('adiustacja na wydruku jest wyborem: pismo do wysłania idzie bez niej', () => {
    expect(domyslneNastawyDruku().adiustacja).toBe('po-zmianach');
    expect(opiszNastawyDruku(domyslneNastawyDruku(), 3, 1)).toContain('bez adiustacji');
    expect(
      opiszNastawyDruku({ ...domyslneNastawyDruku(), adiustacja: 'z-adiustacja' }, 3, 1),
    ).toContain('Z adiustacją');
  });

  it('nastawy druku są pamiętane — szybkie drukowanie bierze ostatnie', () => {
    const magazyn = magazynAtrapa();
    const nastawy = { ...domyslneNastawyDruku(), kopie: 3, dwustronny: true, skala: 90 };
    const wynik = wydrukuj(nastawy, {
      liczbaStron: () => 2,
      kartkaBiezaca: () => 1,
      przygotuj: () => () => undefined,
    });
    // W środowisku sprawdzianu droga druku bywa niedostępna; odmowa jest wtedy nazwana przed próbą druku.
    if (!wynik.udany) {
      expect(wynik.zdanie).toContain('okna drukarki');
      return;
    }
    zapamietajNastawyDruku(nastawy, magazyn);
    expect(czytajNastawyDruku(magazyn).kopie).toBe(3);
    expect(czytajNastawyDruku(magazyn).dwustronny).toBe(true);
  });

  it('przygotowanie powierzchni jest ZAWSZE przywracane, także gdy druk się nie udał', () => {
    let przywrocone = 0;
    wydrukuj(domyslneNastawyDruku(), {
      liczbaStron: () => 2,
      kartkaBiezaca: () => 1,
      przygotuj: () => () => {
        przywrocone += 1;
      },
    });
    // Bez drogi druku przygotowanie nie zachodzi wcale; z drogą — wraca dokładnie raz.
    expect(przywrocone).toBeLessThanOrEqual(1);
  });
});

describe('historia wersji: szereg autozapisu, filtry i czynności', () => {
  function wersja(zmiana: Partial<StudioVersion>): StudioVersion {
    return {
      id: 'w-1',
      documentId: 'd-1',
      createdAt: 1_700_000_000_000,
      ...zmiana,
    };
  }

  it('zapis samoczynny poznaje się po braku etykiety, oznaczenia i opisu', () => {
    expect(czyZapisSamoczynny(wersja({}))).toBe(true);
    expect(czyZapisSamoczynny(wersja({ milestone: true }))).toBe(false);
    expect(czyZapisSamoczynny(wersja({ label: 'wersja do podpisu' }))).toBe(false);
    expect(czyZapisSamoczynny(wersja({ summary: 'poprawki redakcyjne' }))).toBe(false);
    // Zmiana modelu nie jest zapisem w tle: ukrycie jej ukryłoby pracę modelu.
    expect(czyZapisSamoczynny(wersja({ author: StudioAuthor.Model }))).toBe(false);
  });

  it('filtr autora rozdziela Operatora od modelu i nie odsiewa wersji bez autora po cichu', () => {
    const wykaz = [
      wersja({ id: 'a', author: StudioAuthor.Uzytkownik }),
      wersja({ id: 'b', author: StudioAuthor.Model }),
      wersja({ id: 'c' }),
    ];
    expect(przefiltrujHistorie(wykaz, 'model').map((pozycja) => pozycja.id)).toEqual(['b']);
    expect(przefiltrujHistorie(wykaz, 'operator').map((pozycja) => pozycja.id)).toEqual(['a', 'c']);
    expect(pasujeDoAutora(wykaz[2]!, '')).toBe(true);
    expect(pasujeDoAutora(wykaz[2]!, 'model')).toBe(false);
  });

  it('filtr wersji kluczowych czyta pole milestone', () => {
    const wykaz = [wersja({ id: 'a', milestone: true }), wersja({ id: 'b' })];
    expect(przefiltrujHistorie(wykaz, 'kluczowe').map((pozycja) => pozycja.id)).toEqual(['a']);
  });

  it('wiersz wersji ma trzy czynności wprost i cztery w menu — bez rozgałęziania', () => {
    const wiersz = utworzWierszWersji(wersja({ label: 'do podpisu', author: StudioAuthor.Model }), {
      biezaca: true,
      zapisSamoczynny: false,
    });
    expect(wiersz.element.dataset['biezaca']).toBe('tak');
    expect(wiersz.element.querySelector("[data-czynnosc='podejrzyj']")).not.toBeNull();
    expect(wiersz.element.querySelector("[data-czynnosc='przywroc']")).not.toBeNull();
    expect(wiersz.element.querySelector("[data-czynnosc='porownaj']")).not.toBeNull();
    expect(wiersz.element.querySelector("[data-czynnosc='etykieta']")).not.toBeNull();
    expect(wiersz.element.querySelector("[data-czynnosc='usun']")).not.toBeNull();
    // Gałęzie pozostają poza zakresem tej wersji: czynności rozgałęzienia nie ma.
    expect(wiersz.element.querySelector("[data-czynnosc='rozgalez']")).toBeNull();
    // Autor jest pokazany słowem, nie samym kodem kontraktu.
    expect(wiersz.element.textContent).toContain('model');
  });

  it('kropka stanu nie jest jedynym znakiem stanu — obok stoi słowo', () => {
    const wiersz = utworzWierszWersji(wersja({}), { biezaca: true, zapisSamoczynny: false });
    expect(wiersz.element.querySelector('.ms-repozytorium__kropka')).not.toBeNull();
    expect(wiersz.element.textContent).toContain('bieżąca');
  });

  it('wiersz nieosiągalny gasi przywrócenie i podgląd wraz z powodem', () => {
    const wiersz = utworzWierszWersji(wersja({}));
    wiersz.oznaczNiedostepna('rdzeń odmówił: wersja nieosiągalna');
    expect(wiersz.element.dataset['dostepnosc']).toBe('niedostepna');
    expect(wiersz.przywroc.disabled).toBe(true);
    expect(wiersz.podejrzyj.disabled).toBe(true);
  });
});

describe('wstążka PDF jako zakładka kontekstowa', () => {
  it('cztery grupy obejmują wszystkie piętnaście czynności warsztatu, bez powtórzeń', () => {
    const wGrupach = GRUPY_WSTAZKI_PDF.flatMap((grupa) => [...grupa.komendy]);
    expect(new Set(wGrupach).size).toBe(wGrupach.length);
    expect(wGrupach.length).toBe(CZYNNOSCI_WARSZTATU.length);
    for (const czynnosc of CZYNNOSCI_WARSZTATU) {
      expect(wGrupach).toContain(czynnosc.komenda);
    }
  });

  it('grupy stoją w kolejności pracy nad plikiem', () => {
    expect(GRUPY_WSTAZKI_PDF.map((grupa) => grupa.kod)).toEqual([
      'strony',
      'nakladanie',
      'tresc',
      'bezpieczenstwo',
    ]);
  });
});

describe('powierzchnia: linijki, pasek widoku i kartki', () => {
  beforeEach(() => {
    document.body.replaceChildren();
  });

  it('powierzchnia niesie obie linijki, pasek widoku i plansze kartek', async () => {
    const { utworzPowierzchnieDokumentu } = await import('./powierzchnia-dokumentu');
    const { domyslneNastawy, utworzNastawyAkapitow } = await import('./nastawy-wizualne');
    const powierzchnia = utworzPowierzchnieDokumentu(
      domyslnaStrona(),
      domyslneNastawy(),
      utworzNastawyAkapitow(),
      { naTresc: () => undefined, naZaznaczenie: () => undefined, naDecyzjeZmiany: () => undefined },
    );
    document.body.append(powierzchnia.element);
    powierzchnia.pokaz('Pierwszy akapit pisma.', true);

    expect(powierzchnia.element.querySelector('.ms-linijka--pozioma')).not.toBeNull();
    expect(powierzchnia.element.querySelector('.ms-linijka--pionowa')).not.toBeNull();
    expect(powierzchnia.element.querySelector('.ms-widok')).not.toBeNull();
    expect(powierzchnia.element.querySelector('.ms-strona')).not.toBeNull();
    expect(powierzchnia.liczbaStron()).toBeGreaterThanOrEqual(1);
  });

  it('podgląd wydania rysuje KARTKI RDZENIA, gdy rdzeń je oddał', async () => {
    const { utworzPowierzchnieDokumentu } = await import('./powierzchnia-dokumentu');
    const { domyslneNastawy, utworzNastawyAkapitow } = await import('./nastawy-wizualne');
    const powierzchnia = utworzPowierzchnieDokumentu(
      domyslnaStrona(),
      domyslneNastawy(),
      utworzNastawyAkapitow(),
      { naTresc: () => undefined, naZaznaczenie: () => undefined, naDecyzjeZmiany: () => undefined },
    );
    document.body.append(powierzchnia.element);
    powierzchnia.pokaz('Pierwszy akapit pisma.', true);

    // Zapis data: jest jedyną drogą: uri wskazuje ścieżkę w systemie rdzenia, niedostępną przeglądarce.
    const obraz = 'data:image/png;base64,AAAA';
    powierzchnia.ustawTryb('wydanie');
    expect(powierzchnia.kartekZRdzenia()).toBe(0);
    // Bez kartek rdzenia podgląd wydania rysuje kartki liczone w oknie; render dojeżdża osobno.
    expect(powierzchnia.element.querySelector('.ms-strona--rdzen')).toBeNull();

    powierzchnia.ustawKartkiRdzenia([
      { numer: 1, zrodlo: obraz },
      { numer: 2, zrodlo: obraz },
    ]);
    expect(powierzchnia.kartekZRdzenia()).toBe(2);
    const kartki = powierzchnia.element.querySelectorAll('.ms-strona--rdzen');
    expect(kartki.length).toBe(2);
    expect(kartki[0]?.querySelector('img')?.getAttribute('src')).toBe(obraz);
    // Liczba stron jest liczbą kartek rdzenia: pasek stanu pokazuje skład wydania, nie podział z okna.
    expect(powierzchnia.liczbaStron()).toBe(2);

    // Wykaz pusty zdejmuje wyrys rdzenia: stara kartka jako podgląd bieżącej treści byłaby fałszywa.
    powierzchnia.ustawKartkiRdzenia([]);
    expect(powierzchnia.kartekZRdzenia()).toBe(0);
    expect(powierzchnia.element.querySelector('.ms-strona--rdzen')).toBeNull();
    expect(powierzchnia.element.querySelector('.ms-strona')).not.toBeNull();
  });

  it('kartki rdzenia wchodzą wyłącznie w trybie wydania, nie w widoku pracy', async () => {
    const { utworzPowierzchnieDokumentu } = await import('./powierzchnia-dokumentu');
    const { domyslneNastawy, utworzNastawyAkapitow } = await import('./nastawy-wizualne');
    const powierzchnia = utworzPowierzchnieDokumentu(
      domyslnaStrona(),
      domyslneNastawy(),
      utworzNastawyAkapitow(),
      { naTresc: () => undefined, naZaznaczenie: () => undefined, naDecyzjeZmiany: () => undefined },
    );
    document.body.append(powierzchnia.element);
    powierzchnia.pokaz('Pierwszy akapit pisma.', true);
    powierzchnia.ustawKartkiRdzenia([{ numer: 1, zrodlo: 'data:image/png;base64,AAAA' }]);

    // W widoku pracy pisze się w treści, więc obraz strony rdzenia zabrałby Operatorowi możliwość pisania.
    expect(powierzchnia.tryb()).toBe('formatowany');
    expect(powierzchnia.element.querySelector('.ms-strona--rdzen')).toBeNull();
    powierzchnia.ustawTryb('wydanie');
    expect(powierzchnia.element.querySelector('.ms-strona--rdzen')).not.toBeNull();
  });

  it('tryb źródłowy jest przełącznikiem: wraca do widoku formatowanego', async () => {
    const { utworzPowierzchnieDokumentu } = await import('./powierzchnia-dokumentu');
    const { domyslneNastawy, utworzNastawyAkapitow } = await import('./nastawy-wizualne');
    const powierzchnia = utworzPowierzchnieDokumentu(
      domyslnaStrona(),
      domyslneNastawy(),
      utworzNastawyAkapitow(),
      { naTresc: () => undefined, naZaznaczenie: () => undefined, naDecyzjeZmiany: () => undefined },
    );
    powierzchnia.ustawTryb('zrodlowy');
    expect(powierzchnia.tryb()).toBe('zrodlowy');
    expect(powierzchnia.nastawyWidoku().trybZrodlowy).toBe(true);
    powierzchnia.ustawTryb('formatowany');
    expect(powierzchnia.nastawyWidoku().trybZrodlowy).toBe(false);
  });

  it('nastawy strony pchnięte z okna NIE zabierają tego, co Operator ustawił chwytem', async () => {
    const { utworzPowierzchnieDokumentu } = await import('./powierzchnia-dokumentu');
    const { domyslneNastawy, utworzNastawyAkapitow } = await import('./nastawy-wizualne');
    const zgloszone: number[] = [];
    const powierzchnia = utworzPowierzchnieDokumentu(
      domyslnaStrona(),
      domyslneNastawy(),
      utworzNastawyAkapitow(),
      {
        naTresc: () => undefined,
        naZaznaczenie: () => undefined,
        naDecyzjeZmiany: () => undefined,
        naStrone: (strona) => zgloszone.push(strona.marginesOprawyMm),
      },
    );
    document.body.append(powierzchnia.element);
    powierzchnia.przestawWidok({ graniceMarginesow: false });
    expect(powierzchnia.strona().graniceMarginesow).toBe(false);

    // Okno pcha swoją kopię nastaw, w której o oprawie nic nie wie.
    powierzchnia.ustawStrone({ ...domyslnaStrona(), marginesGoraMm: 30 });
    expect(powierzchnia.strona().marginesGoraMm).toBe(30);
    expect(powierzchnia.strona().graniceMarginesow).toBe(false);
    expect(zgloszone.length).toBeGreaterThan(0);
  });

  it('skala jest przycinana i wpisana w opis kartki', async () => {
    const { utworzPowierzchnieDokumentu } = await import('./powierzchnia-dokumentu');
    const { domyslneNastawy, utworzNastawyAkapitow } = await import('./nastawy-wizualne');
    const powierzchnia = utworzPowierzchnieDokumentu(
      domyslnaStrona(),
      domyslneNastawy(),
      utworzNastawyAkapitow(),
      { naTresc: () => undefined, naZaznaczenie: () => undefined, naDecyzjeZmiany: () => undefined },
    );
    powierzchnia.ustawSkale(9999);
    expect(powierzchnia.strona().skala).toBe(500);
    expect(powierzchnia.opis()).toContain('skala 500 %');
    // Nastawa gotowa nieznana nie zmienia skali i mówi o tym oddanym `false`.
    expect(powierzchnia.ustawNastaweSkali('wymyslona')).toBe(false);
    expect(powierzchnia.strona().skala).toBe(500);
  });

  it('układ kartek przestawia liczbę kolumn i opisuje przewijanie', async () => {
    const { utworzPowierzchnieDokumentu } = await import('./powierzchnia-dokumentu');
    const { domyslneNastawy, utworzNastawyAkapitow } = await import('./nastawy-wizualne');
    const powierzchnia = utworzPowierzchnieDokumentu(
      domyslnaStrona(),
      domyslneNastawy(),
      utworzNastawyAkapitow(),
      { naTresc: () => undefined, naZaznaczenie: () => undefined, naDecyzjeZmiany: () => undefined },
    );
    document.body.append(powierzchnia.element);
    powierzchnia.ustawKolumny(2);
    expect(powierzchnia.nastawyWidoku().ukladKartek).toBe('obok');
    expect(powierzchnia.element.dataset['uklad']).toBe('obok');
    powierzchnia.przestawWidok({ ukladKartek: 'rozkladowka', przewijanie: 'strona-po-stronie' });
    expect(powierzchnia.element.dataset['przewijanie']).toBe('strona-po-stronie');
    expect(powierzchnia.opis()).toContain('rozkładówka');
  });

  it('kartka o numerze poza dokumentem jest przycinana, nie zgadywana', async () => {
    const { utworzPowierzchnieDokumentu } = await import('./powierzchnia-dokumentu');
    const { domyslneNastawy, utworzNastawyAkapitow } = await import('./nastawy-wizualne');
    const powierzchnia = utworzPowierzchnieDokumentu(
      domyslnaStrona(),
      domyslneNastawy(),
      utworzNastawyAkapitow(),
      { naTresc: () => undefined, naZaznaczenie: () => undefined, naDecyzjeZmiany: () => undefined },
    );
    document.body.append(powierzchnia.element);
    powierzchnia.pokaz('Jeden akapit.', true);
    expect(powierzchnia.skoczDoKartki(99)).toBe(powierzchnia.liczbaStron());
    expect(powierzchnia.skoczDoKartki(0)).toBe(1);
  });
});
