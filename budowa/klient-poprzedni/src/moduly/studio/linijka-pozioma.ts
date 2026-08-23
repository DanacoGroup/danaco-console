import {
  marginesyKartki,
  opiszDlugosc,
  szerokoscKartkiMm,
  type StronaPracy,
} from './nastawy-strony';
import {
  krokPrzyciaganiaLinijki,
  milimetryZPunktowLinijki,
  nastepnyRodzajTabulatora,
  nazwaRodzajuTabulatora,
  nazwaZnakuWiodacego,
  opiszWciecieAkapitu,
  przyciagnijNaLinijce,
  punktyNaLinijce,
  zerowaWciecieAkapitu,
  zlozPodzialkeLinijki,
  znakTabulatora,
  type RodzajTabulatora,
  type TabulatorAkapitu,
  type WciecieAkapitu,
  type ZnakWiodacy,
} from './linijka-podzialka';

/**
 * Linijka pozioma — podziałka z chwytami, nie sama miarka do patrzenia.
 *
 * ── Co się na niej chwyta ───────────────────────────────────────────────────
 * Margines lewy i prawy, wcięcie pierwszego wiersza, wcięcie lewe i prawe
 * osobnymi znacznikami, tabulatory zakładane naciśnięciem wraz z rodzajem
 * i znakiem wiodącym, oraz szerokości kolumn tabeli. Wszystko chwytem I z
 * klawiatury: każdy chwyt jest `slider` z wartością w milimetrach albo calach, bo
 * przestawianie marginesu pisma urzędowego wyłącznie myszą odcięłoby połowę
 * Operatorów.
 *
 * ── Dlaczego wcięcie pierwszego wiersza liczy się względnie ─────────────────
 * Bo tak zachowuje się linijka pakietu biurowego: przeciągnięcie wcięcia lewego
 * zabiera pierwszy wiersz ze sobą. Gdyby oba wcięcia były liczone od krawędzi
 * pola, Operator przy każdej zmianie wcięcia lewego poprawiałby drugie — a to nie
 * jest praca, to nadrabianie za oknem.
 *
 * ── Czego linijka nie robi ──────────────────────────────────────────────────
 * Nie zmienia dokumentu i nie woła rdzenia. Oddaje nastawę temu, kto ją zbudował
 * (`CzynnosciLinijkiPoziomej`), a nastawy strony i wcięcia trzyma powierzchnia —
 * inaczej byłyby dwa źródła prawdy o marginesie i rozjechałyby się przy pierwszym
 * profilu wydania wziętym z rdzenia.
 */

/** Czynności linijki poziomej zlecane powierzchni. */
export interface CzynnosciLinijkiPoziomej {
  /** Margines przestawiony chwytem. */
  naMargines(strona: 'lewy' | 'prawy', milimetry: number): void;
  /** Wcięcia akapitu przestawione chwytem. */
  naWciecie(wciecie: WciecieAkapitu): void;
  /** Wykaz tabulatorów po zmianie — założeniu, przestawieniu albo zdjęciu. */
  naTabulatory(tabulatory: readonly TabulatorAkapitu[]): void;
  /**
   * Krawędź kolumny tabeli przestawiona chwytem.
   *
   * Numer krawędzi liczony od lewej, od zera; wartość jest odległością od lewej
   * krawędzi pola pisania w milimetrach.
   */
  naKrawedzKolumny(numer: number, milimetry: number): void;
}

/** Linijka pozioma wraz z jej sterowaniem. */
export interface LinijkaPozioma {
  element: HTMLElement;
  ustawStrone(strona: StronaPracy): void;
  /** Skala widoku w procentach — podziałka rozciąga się razem z kartką. */
  ustawSkale(procent: number): void;
  ustawWciecie(wciecie: WciecieAkapitu): void;
  wciecie(): WciecieAkapitu;
  ustawTabulatory(tabulatory: readonly TabulatorAkapitu[]): void;
  tabulatory(): readonly TabulatorAkapitu[];
  /** Krawędzie kolumn tabeli w milimetrach od lewej krawędzi pola pisania. */
  ustawKrawedzieKolumn(krawedzie: readonly number[]): void;
  /** Położenie kursora w milimetrach od lewej krawędzi pola; `null` chowa znacznik. */
  ustawKursor(milimetry: number | null): void;
  /** Granice zaznaczenia w milimetrach; `null` chowa pasek zaznaczenia. */
  ustawZaznaczenie(zakres: { odMm: number; doMm: number } | null): void;
  ustawWidocznosc(widoczna: boolean): void;
  widoczna(): boolean;
  /** Zdanie o linijce — do paska stanu. */
  opis(): string;
}

/** Najmniejsze pole pisania, jakie chwyt marginesu wolno zostawić. */
const NAJMNIEJSZE_POLE_MM = 10;

export function utworzLinijkePozioma(
  strona: StronaPracy,
  czynnosci: CzynnosciLinijkiPoziomej,
): LinijkaPozioma {
  let stronaBiezaca = strona;
  let skala = strona.skala / 100;
  let wciecieBiezace = zerowaWciecieAkapitu();
  let tabulatoryBiezace: TabulatorAkapitu[] = [];
  let krawedzieKolumn: number[] = [];
  let kursorMm: number | null = null;
  let zaznaczenieMm: { odMm: number; doMm: number } | null = null;
  let rodzajNowego: RodzajTabulatora = 'lewy';
  let znakNowego: ZnakWiodacy = 'brak';

  const podzialka = document.createElement('div');
  podzialka.className = 'ms-linijka__podzialka';

  const chwyty = document.createElement('div');
  chwyty.className = 'ms-linijka__chwyty';

  /**
   * Wybieracz rodzaju tabulatora — narożnik linijki, wzorem pakietu biurowego.
   *
   * Jedno naciśnięcie przestawia rodzaj następnego zakładanego tabulatora. Rodzaj
   * jest widoczny na przycisku, żeby Operator wiedział, co założy, ZANIM naciśnie
   * podziałkę.
   */
  const wybieracz = document.createElement('button');
  wybieracz.type = 'button';
  wybieracz.className = 'ms-linijka__wybieracz';
  wybieracz.dataset['czynnosc'] = 'rodzaj-tabulatora';
  wybieracz.addEventListener('click', () => {
    rodzajNowego = nastepnyRodzajTabulatora(rodzajNowego);
    opiszWybieracz();
  });

  const znakWiodacy = document.createElement('select');
  znakWiodacy.className = 'dn-pole-kontrolka ms-linijka__znak';
  znakWiodacy.setAttribute('aria-label', 'Znak wiodący zakładanego tabulatora');
  for (const znak of ['brak', 'kropka', 'kreska', 'podkreslenie'] as const) {
    const opcja = document.createElement('option');
    opcja.value = znak;
    opcja.textContent = nazwaZnakuWiodacego(znak);
    znakWiodacy.append(opcja);
  }
  znakWiodacy.addEventListener('change', () => {
    znakNowego = znakWiodacy.value as ZnakWiodacy;
  });

  function opiszWybieracz(): void {
    wybieracz.textContent = znakTabulatora(rodzajNowego);
    wybieracz.title =
      `Rodzaj zakładanego tabulatora: ${nazwaRodzajuTabulatora(rodzajNowego)}. Naciśnij, żeby ` +
      'przestawić na następny; potem naciśnij podziałkę w miejscu, w którym tabulator ma stanąć.';
    wybieracz.setAttribute('aria-label', wybieracz.title);
  }
  opiszWybieracz();

  const pas = document.createElement('div');
  pas.className = 'ms-linijka__pas';
  pas.append(podzialka, chwyty);

  const element = document.createElement('div');
  element.className = 'ms-linijka ms-linijka--pozioma';
  element.setAttribute('aria-label', 'Linijka pozioma — marginesy, wcięcia, tabulatory, kolumny tabeli');
  element.append(wybieracz, znakWiodacy, pas);

  /* ── Przeliczenia ────────────────────────────────────────────────────────── */

  /** Milimetry od lewej krawędzi KARTKI dla położenia myszy. */
  function milimetryZdarzenia(zdarzenie: { clientX: number }): number {
    const prostokat = pas.getBoundingClientRect();
    return milimetryZPunktowLinijki(zdarzenie.clientX - prostokat.left, skala);
  }

  function punkty(milimetry: number): number {
    return punktyNaLinijce(milimetry, skala);
  }

  function marginesy(): { lewyMm: number; prawyMm: number } {
    const wszystkie = marginesyKartki(stronaBiezaca, 1);
    return { lewyMm: wszystkie.lewyMm, prawyMm: wszystkie.prawyMm };
  }

  function szerokoscPola(): number {
    const { lewyMm, prawyMm } = marginesy();
    return Math.max(1, szerokoscKartkiMm(stronaBiezaca) - lewyMm - prawyMm);
  }

  /* ── Chwyt ───────────────────────────────────────────────────────────────── */

  /** Opis jednego chwytu linijki. */
  interface OpisChwytu {
    kod: string;
    nazwa: string;
    znak: string;
    /** Wartość chwytu w milimetrach — miara jego własna, nie położenie na kartce. */
    wartosc(): number;
    /** Położenie chwytu na kartce w milimetrach od jej lewej krawędzi. */
    polozenie(): number;
    /** Granice wartości chwytu. */
    granice(): { dolnaMm: number; gornaMm: number };
    /** Wartość z położenia myszy na kartce. */
    zPolozenia(naKartceMm: number): number;
    /** Nastawa po przestawieniu chwytu. */
    ustaw(milimetry: number): void;
  }

  const chwytyZlozone: { opis: OpisChwytu; element: HTMLElement }[] = [];

  function utworzChwyt(opis: OpisChwytu): HTMLElement {
    const uchwyt = document.createElement('button');
    uchwyt.type = 'button';
    uchwyt.className = 'ms-linijka__chwyt';
    uchwyt.dataset['chwyt'] = opis.kod;
    uchwyt.textContent = opis.znak;
    uchwyt.setAttribute('role', 'slider');
    uchwyt.setAttribute('aria-orientation', 'horizontal');

    let wleczony = false;

    uchwyt.addEventListener('mousedown', (zdarzenie) => {
      zdarzenie.preventDefault();
      zdarzenie.stopPropagation();
      wleczony = true;
      uchwyt.dataset['wleczenie'] = 'tak';
    });

    document.addEventListener('mousemove', (zdarzenie) => {
      if (!wleczony) return;
      przestawChwyt(opis, opis.zPolozenia(milimetryZdarzenia(zdarzenie)));
    });

    document.addEventListener('mouseup', () => {
      if (!wleczony) return;
      wleczony = false;
      delete uchwyt.dataset['wleczenie'];
    });

    uchwyt.addEventListener('keydown', (zdarzenie) => {
      const krok = krokPrzyciaganiaLinijki(stronaBiezaca.jednostka);
      const wiekszy = zdarzenie.shiftKey ? krok * 10 : krok;
      if (zdarzenie.key === 'ArrowLeft' || zdarzenie.key === 'ArrowDown') {
        zdarzenie.preventDefault();
        przestawChwyt(opis, opis.wartosc() - wiekszy);
        return;
      }
      if (zdarzenie.key === 'ArrowRight' || zdarzenie.key === 'ArrowUp') {
        zdarzenie.preventDefault();
        przestawChwyt(opis, opis.wartosc() + wiekszy);
        return;
      }
      if (zdarzenie.key === 'Home') {
        zdarzenie.preventDefault();
        przestawChwyt(opis, opis.granice().dolnaMm);
        return;
      }
      if (zdarzenie.key === 'End') {
        zdarzenie.preventDefault();
        przestawChwyt(opis, opis.granice().gornaMm);
      }
    });

    chwytyZlozone.push({ opis, element: uchwyt });
    return uchwyt;
  }

  function przestawChwyt(opis: OpisChwytu, wartosc: number): void {
    const granice = opis.granice();
    opis.ustaw(
      przyciagnijNaLinijce(wartosc, stronaBiezaca.jednostka, granice.dolnaMm, granice.gornaMm),
    );
    przerysuj();
  }

  /* ── Chwyty stałe: marginesy i wcięcia ───────────────────────────────────── */

  const chwytMarginesuLewego = utworzChwyt({
    kod: 'margines-lewy',
    nazwa: 'margines lewy',
    znak: '▎',
    wartosc: () => marginesy().lewyMm,
    polozenie: () => marginesy().lewyMm,
    granice: () => ({
      dolnaMm: 0,
      gornaMm: Math.max(
        0,
        szerokoscKartkiMm(stronaBiezaca) - marginesy().prawyMm - NAJMNIEJSZE_POLE_MM,
      ),
    }),
    zPolozenia: (naKartce) => naKartce,
    ustaw: (milimetry) => czynnosci.naMargines('lewy', milimetry),
  });

  const chwytMarginesuPrawego = utworzChwyt({
    kod: 'margines-prawy',
    nazwa: 'margines prawy',
    znak: '▊',
    wartosc: () => marginesy().prawyMm,
    polozenie: () => szerokoscKartkiMm(stronaBiezaca) - marginesy().prawyMm,
    granice: () => ({
      dolnaMm: 0,
      gornaMm: Math.max(
        0,
        szerokoscKartkiMm(stronaBiezaca) - marginesy().lewyMm - NAJMNIEJSZE_POLE_MM,
      ),
    }),
    // Chwyt prawy jedzie w drugą stronę: im dalej w prawo, tym margines mniejszy.
    zPolozenia: (naKartce) => szerokoscKartkiMm(stronaBiezaca) - naKartce,
    ustaw: (milimetry) => czynnosci.naMargines('prawy', milimetry),
  });

  const chwytPierwszegoWiersza = utworzChwyt({
    kod: 'wciecie-pierwszego-wiersza',
    nazwa: 'wcięcie pierwszego wiersza',
    znak: '▽',
    wartosc: () => wciecieBiezace.leweMm + wciecieBiezace.pierwszyWierszMm,
    polozenie: () =>
      marginesy().lewyMm + wciecieBiezace.leweMm + wciecieBiezace.pierwszyWierszMm,
    granice: () => ({ dolnaMm: 0, gornaMm: szerokoscPola() - wciecieBiezace.praweMm }),
    zPolozenia: (naKartce) => naKartce - marginesy().lewyMm,
    ustaw: (milimetry) =>
      czynnosci.naWciecie({
        ...wciecieBiezace,
        pierwszyWierszMm: milimetry - wciecieBiezace.leweMm,
      }),
  });

  const chwytWciecieLewego = utworzChwyt({
    kod: 'wciecie-lewe',
    nazwa: 'wcięcie lewe',
    znak: '△',
    wartosc: () => wciecieBiezace.leweMm,
    polozenie: () => marginesy().lewyMm + wciecieBiezace.leweMm,
    granice: () => ({ dolnaMm: 0, gornaMm: szerokoscPola() - wciecieBiezace.praweMm }),
    zPolozenia: (naKartce) => naKartce - marginesy().lewyMm,
    // Wcięcie pierwszego wiersza jest liczone względem lewego, więc jedzie razem
    // z nim samo — bez dopisywania czegokolwiek do nastawy.
    ustaw: (milimetry) => czynnosci.naWciecie({ ...wciecieBiezace, leweMm: milimetry }),
  });

  const chwytWciecieprawego = utworzChwyt({
    kod: 'wciecie-prawe',
    nazwa: 'wcięcie prawe',
    znak: '△',
    wartosc: () => wciecieBiezace.praweMm,
    polozenie: () =>
      marginesy().lewyMm + szerokoscPola() - wciecieBiezace.praweMm,
    granice: () => ({ dolnaMm: 0, gornaMm: szerokoscPola() - wciecieBiezace.leweMm }),
    zPolozenia: (naKartce) => marginesy().lewyMm + szerokoscPola() - naKartce,
    ustaw: (milimetry) => czynnosci.naWciecie({ ...wciecieBiezace, praweMm: milimetry }),
  });

  /* ── Zakładanie tabulatora naciśnięciem podziałki ────────────────────────── */

  podzialka.addEventListener('click', (zdarzenie) => {
    const naKartce = milimetryZdarzenia(zdarzenie);
    const { lewyMm } = marginesy();
    const wPolu = przyciagnijNaLinijce(
      naKartce - lewyMm,
      stronaBiezaca.jednostka,
      0,
      szerokoscPola(),
    );
    tabulatoryBiezace = [
      ...tabulatoryBiezace,
      { milimetry: wPolu, rodzaj: rodzajNowego, znakWiodacy: znakNowego },
    ].sort((jeden, dwa) => jeden.milimetry - dwa.milimetry);
    czynnosci.naTabulatory(tabulatoryBiezace);
    przerysuj();
  });

  /* ── Wyrys ───────────────────────────────────────────────────────────────── */

  const pasZaznaczenia = document.createElement('div');
  pasZaznaczenia.className = 'ms-linijka__zaznaczenie';
  pasZaznaczenia.hidden = true;

  const znacznikKursora = document.createElement('div');
  znacznikKursora.className = 'ms-linijka__kursor';
  znacznikKursora.hidden = true;

  const polaMarginesow = document.createElement('div');
  polaMarginesow.className = 'ms-linijka__marginesy';

  chwyty.append(
    polaMarginesow,
    pasZaznaczenia,
    znacznikKursora,
    chwytMarginesuLewego,
    chwytMarginesuPrawego,
    chwytPierwszegoWiersza,
    chwytWciecieLewego,
    chwytWciecieprawego,
  );

  /** Znaczniki zakładane i zdejmowane przy każdym wyrysie — tabulatory i kolumny. */
  const znacznikiRuchome: HTMLElement[] = [];

  function przerysuj(): void {
    const szerokoscKartki = szerokoscKartkiMm(stronaBiezaca);
    pas.style.width = `${punkty(szerokoscKartki)}px`;

    // Podziałka liczona od lewej krawędzi POLA PISANIA, nie kartki: Operator mierzy
    // wcięcie od marginesu, a nie od krawędzi papieru.
    const { lewyMm, prawyMm } = marginesy();
    podzialka.replaceChildren();
    for (const kreska of zlozPodzialkeLinijki(szerokoscPola(), stronaBiezaca.jednostka, skala)) {
      const znak = document.createElement('span');
      znak.className = 'ms-linijka__kreska';
      znak.dataset['kreska'] = kreska.rodzaj;
      znak.style.left = `${punkty(lewyMm + kreska.milimetry)}px`;
      if (kreska.napis !== '') znak.textContent = kreska.napis;
      podzialka.append(znak);
    }

    polaMarginesow.style.setProperty('--ms-linijka-lewy', `${punkty(lewyMm)}px`);
    polaMarginesow.style.setProperty('--ms-linijka-prawy', `${punkty(prawyMm)}px`);

    for (const wpis of chwytyZlozone) {
      const wartosc = wpis.opis.wartosc();
      wpis.element.style.left = `${punkty(wpis.opis.polozenie())}px`;
      wpis.element.setAttribute('aria-valuenow', String(Math.round(wartosc * 10) / 10));
      wpis.element.setAttribute(
        'aria-valuetext',
        `${wpis.opis.nazwa}: ${opiszDlugosc(wartosc, stronaBiezaca.jednostka)}`,
      );
      const granice = wpis.opis.granice();
      wpis.element.setAttribute('aria-valuemin', String(Math.round(granice.dolnaMm * 10) / 10));
      wpis.element.setAttribute('aria-valuemax', String(Math.round(granice.gornaMm * 10) / 10));
      wpis.element.title =
        `${wpis.opis.nazwa} — ${opiszDlugosc(wartosc, stronaBiezaca.jednostka)}. Przeciągnij albo ` +
        'przesuń strzałkami; ze Shiftem krok dziesięciokrotny, Home i End dosuwają do granic.';
    }

    for (const znacznik of znacznikiRuchome) znacznik.remove();
    znacznikiRuchome.length = 0;

    tabulatoryBiezace.forEach((tabulator, numer) => {
      const znak = document.createElement('button');
      znak.type = 'button';
      znak.className = 'ms-linijka__tabulator';
      znak.dataset['tabulator'] = String(numer);
      znak.dataset['rodzaj'] = tabulator.rodzaj;
      znak.textContent = znakTabulatora(tabulator.rodzaj);
      znak.style.left = `${punkty(lewyMm + tabulator.milimetry)}px`;
      znak.title =
        `${nazwaRodzajuTabulatora(tabulator.rodzaj)} na ` +
        `${opiszDlugosc(tabulator.milimetry, stronaBiezaca.jednostka)} · ` +
        `${nazwaZnakuWiodacego(tabulator.znakWiodacy)}. Naciśnięcie przestawia rodzaj, ` +
        'naciśnięcie z Shiftem zdejmuje tabulator.';
      znak.setAttribute('aria-label', znak.title);
      znak.addEventListener('click', (zdarzenie) => {
        zdarzenie.stopPropagation();
        if (zdarzenie.shiftKey) {
          tabulatoryBiezace = tabulatoryBiezace.filter((_, inny) => inny !== numer);
        } else {
          tabulatoryBiezace = tabulatoryBiezace.map((wpis, inny) =>
            inny === numer ? { ...wpis, rodzaj: nastepnyRodzajTabulatora(wpis.rodzaj) } : wpis,
          );
        }
        czynnosci.naTabulatory(tabulatoryBiezace);
        przerysuj();
      });
      chwyty.append(znak);
      znacznikiRuchome.push(znak);
    });

    krawedzieKolumn.forEach((krawedz, numer) => {
      const znak = document.createElement('button');
      znak.type = 'button';
      znak.className = 'ms-linijka__kolumna';
      znak.dataset['kolumna'] = String(numer);
      znak.textContent = '⇹';
      znak.style.left = `${punkty(lewyMm + krawedz)}px`;
      znak.setAttribute('role', 'slider');
      znak.setAttribute('aria-valuenow', String(Math.round(krawedz * 10) / 10));
      znak.title =
        `Krawędź kolumny tabeli numer ${numer + 1} na ` +
        `${opiszDlugosc(krawedz, stronaBiezaca.jednostka)}. Przeciągnij albo przesuń strzałkami, ` +
        'żeby zmienić szerokości kolumn sąsiadujących.';
      znak.setAttribute('aria-label', znak.title);

      let wleczona = false;
      znak.addEventListener('mousedown', (zdarzenie) => {
        zdarzenie.preventDefault();
        zdarzenie.stopPropagation();
        wleczona = true;
      });
      document.addEventListener('mousemove', (zdarzenie) => {
        if (!wleczona) return;
        przestawKrawedz(numer, milimetryZdarzenia(zdarzenie) - lewyMm);
      });
      document.addEventListener('mouseup', () => {
        wleczona = false;
      });
      znak.addEventListener('keydown', (zdarzenie) => {
        const krok = krokPrzyciaganiaLinijki(stronaBiezaca.jednostka);
        if (zdarzenie.key === 'ArrowLeft') {
          zdarzenie.preventDefault();
          przestawKrawedz(numer, krawedz - krok);
          return;
        }
        if (zdarzenie.key === 'ArrowRight') {
          zdarzenie.preventDefault();
          przestawKrawedz(numer, krawedz + krok);
        }
      });
      chwyty.append(znak);
      znacznikiRuchome.push(znak);
    });

    if (kursorMm === null) {
      znacznikKursora.hidden = true;
    } else {
      znacznikKursora.hidden = false;
      znacznikKursora.style.left = `${punkty(lewyMm + kursorMm)}px`;
      znacznikKursora.title = `Kursor na ${opiszDlugosc(kursorMm, stronaBiezaca.jednostka)} od marginesu lewego.`;
    }

    if (zaznaczenieMm === null) {
      pasZaznaczenia.hidden = true;
    } else {
      pasZaznaczenia.hidden = false;
      pasZaznaczenia.style.left = `${punkty(lewyMm + zaznaczenieMm.odMm)}px`;
      pasZaznaczenia.style.width = `${punkty(Math.max(0, zaznaczenieMm.doMm - zaznaczenieMm.odMm))}px`;
    }
  }

  /**
   * Przestawia krawędź kolumny, pilnując kolejności krawędzi.
   *
   * Krawędź nie może przejść przez sąsiednią: kolumna o szerokości ujemnej nie
   * jest kolumną. Granicą jest sąsiad odsunięty o najmniejsze pole, a nie sam
   * sąsiad — dwie krawędzie w jednym miejscu dałyby kolumnę zerową.
   */
  function przestawKrawedz(numer: number, milimetry: number): void {
    const poprzednia = krawedzieKolumn[numer - 1];
    const nastepna = krawedzieKolumn[numer + 1];
    const dolna = (poprzednia ?? 0) + NAJMNIEJSZE_POLE_MM / 2;
    const gorna = (nastepna ?? szerokoscPola()) - NAJMNIEJSZE_POLE_MM / 2;
    if (gorna <= dolna) return;
    const wartosc = przyciagnijNaLinijce(milimetry, stronaBiezaca.jednostka, dolna, gorna);
    krawedzieKolumn = krawedzieKolumn.map((inna, indeks) => (indeks === numer ? wartosc : inna));
    czynnosci.naKrawedzKolumny(numer, wartosc);
    przerysuj();
  }

  przerysuj();

  return {
    element,

    ustawStrone(nowa) {
      stronaBiezaca = nowa;
      skala = nowa.skala / 100;
      znakWiodacy.value = znakNowego;
      przerysuj();
    },

    ustawSkale(procent) {
      skala = procent / 100;
      przerysuj();
    },

    ustawWciecie(wciecie) {
      wciecieBiezace = wciecie;
      przerysuj();
    },

    wciecie: () => wciecieBiezace,

    ustawTabulatory(tabulatory) {
      tabulatoryBiezace = [...tabulatory].sort((jeden, dwa) => jeden.milimetry - dwa.milimetry);
      przerysuj();
    },

    tabulatory: () => tabulatoryBiezace,

    ustawKrawedzieKolumn(krawedzie) {
      krawedzieKolumn = [...krawedzie].sort((jeden, dwa) => jeden - dwa);
      przerysuj();
    },

    ustawKursor(milimetry) {
      kursorMm = milimetry;
      przerysuj();
    },

    ustawZaznaczenie(zakres) {
      zaznaczenieMm = zakres;
      przerysuj();
    },

    ustawWidocznosc(widoczna) {
      element.hidden = !widoczna;
    },

    widoczna: () => !element.hidden,

    opis() {
      const tabulatory =
        tabulatoryBiezace.length === 0
          ? 'bez tabulatorów'
          : `tabulatorów ${tabulatoryBiezace.length}`;
      const kolumny =
        krawedzieKolumn.length === 0 ? '' : ` · krawędzi kolumn ${krawedzieKolumn.length}`;
      return `${opiszWciecieAkapitu(wciecieBiezace, stronaBiezaca.jednostka)} · ${tabulatory}${kolumny}`;
    },
  };
}
