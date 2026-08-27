import {
  marginesyKartki,
  opiszDlugosc,
  wysokoscKartkiMm,
  type StronaPracy,
} from './nastawy-strony';
import {
  krokPrzyciaganiaLinijki,
  milimetryZPunktowLinijki,
  przyciagnijNaLinijce,
  punktyNaLinijce,
  zlozPodzialkeLinijki,
} from './linijka-podzialka';

/** Czynności linijki pionowej zlecane powierzchni dokumentu: margines przestawiony chwytem górnym albo dolnym. */
export interface CzynnosciLinijkiPionowej {
  /** Margines przestawiony chwytem. */
  naMargines(strona: 'gora' | 'dol', milimetry: number): void;
}

/** Linijka pionowa wraz z jej pełnym sterowaniem: stroną, skalą, kursorem, widocznością i zdaniem opisowym. */
export interface LinijkaPionowa {
  element: HTMLElement;
  ustawStrone(strona: StronaPracy): void;
  ustawSkale(procent: number): void;
  /** Położenie kursora w milimetrach od górnej krawędzi pola; `null` chowa znacznik. */
  ustawKursor(milimetry: number | null): void;
  ustawWidocznosc(widoczna: boolean): void;
  widoczna(): boolean;
  opis(): string;
}

/** Najmniejsze pole pisania na stronie dokumentu, jakie chwyt marginesu na linijce pionowej wolno zostawić. */
const NAJMNIEJSZE_POLE_MM = 10;

export function utworzLinijkePionowa(
  strona: StronaPracy,
  czynnosci: CzynnosciLinijkiPionowej,
): LinijkaPionowa {
  let stronaBiezaca = strona;
  let skala = strona.skala / 100;
  let kursorMm: number | null = null;

  const podzialka = document.createElement('div');
  podzialka.className = 'ms-linijka__podzialka';

  const chwyty = document.createElement('div');
  chwyty.className = 'ms-linijka__chwyty';

  const pas = document.createElement('div');
  pas.className = 'ms-linijka__pas';
  pas.append(podzialka, chwyty);

  const element = document.createElement('div');
  element.className = 'ms-linijka ms-linijka--pionowa';
  element.setAttribute('aria-label', 'Linijka pionowa — margines górny i dolny');
  element.append(pas);

  function punkty(milimetry: number): number {
    return punktyNaLinijce(milimetry, skala);
  }

  function marginesy(): { goraMm: number; dolMm: number } {
    const wszystkie = marginesyKartki(stronaBiezaca, 1);
    return { goraMm: wszystkie.goraMm, dolMm: wszystkie.dolMm };
  }

  function wysokoscPola(): number {
    const { goraMm, dolMm } = marginesy();
    return Math.max(1, wysokoscKartkiMm(stronaBiezaca) - goraMm - dolMm);
  }

  function milimetryZdarzenia(zdarzenie: { clientY: number }): number {
    const prostokat = pas.getBoundingClientRect();
    return milimetryZPunktowLinijki(zdarzenie.clientY - prostokat.top, skala);
  }

  /** Opis chwytu pionowego — dwa chwyty, więc bez wykazu, wprost. */
  interface OpisChwytuPionowego {
    kod: string;
    nazwa: string;
    wartosc(): number;
    polozenie(): number;
    granice(): { dolnaMm: number; gornaMm: number };
    zPolozenia(naKartceMm: number): number;
    ustaw(milimetry: number): void;
  }

  const zlozone: { opis: OpisChwytuPionowego; element: HTMLElement }[] = [];

  function utworzChwyt(opis: OpisChwytuPionowego): HTMLElement {
    const uchwyt = document.createElement('button');
    uchwyt.type = 'button';
    uchwyt.className = 'ms-linijka__chwyt';
    uchwyt.dataset['chwyt'] = opis.kod;
    uchwyt.textContent = '▬';
    uchwyt.setAttribute('role', 'slider');
    uchwyt.setAttribute('aria-orientation', 'vertical');

    let wleczony = false;
    uchwyt.addEventListener('mousedown', (zdarzenie) => {
      zdarzenie.preventDefault();
      wleczony = true;
      uchwyt.dataset['wleczenie'] = 'tak';
    });
    document.addEventListener('mousemove', (zdarzenie) => {
      if (!wleczony) return;
      przestaw(opis, opis.zPolozenia(milimetryZdarzenia(zdarzenie)));
    });
    document.addEventListener('mouseup', () => {
      if (!wleczony) return;
      wleczony = false;
      delete uchwyt.dataset['wleczenie'];
    });
    uchwyt.addEventListener('keydown', (zdarzenie) => {
      const krok = krokPrzyciaganiaLinijki(stronaBiezaca.jednostka);
      const wiekszy = zdarzenie.shiftKey ? krok * 10 : krok;
      if (zdarzenie.key === 'ArrowUp' || zdarzenie.key === 'ArrowLeft') {
        zdarzenie.preventDefault();
        przestaw(opis, opis.wartosc() - wiekszy);
        return;
      }
      if (zdarzenie.key === 'ArrowDown' || zdarzenie.key === 'ArrowRight') {
        zdarzenie.preventDefault();
        przestaw(opis, opis.wartosc() + wiekszy);
        return;
      }
      if (zdarzenie.key === 'Home') {
        zdarzenie.preventDefault();
        przestaw(opis, opis.granice().dolnaMm);
        return;
      }
      if (zdarzenie.key === 'End') {
        zdarzenie.preventDefault();
        przestaw(opis, opis.granice().gornaMm);
      }
    });

    zlozone.push({ opis, element: uchwyt });
    return uchwyt;
  }

  function przestaw(opis: OpisChwytuPionowego, wartosc: number): void {
    const granice = opis.granice();
    opis.ustaw(
      przyciagnijNaLinijce(wartosc, stronaBiezaca.jednostka, granice.dolnaMm, granice.gornaMm),
    );
    przerysuj();
  }

  const chwytGory = utworzChwyt({
    kod: 'margines-gora',
    nazwa: 'margines górny',
    wartosc: () => marginesy().goraMm,
    polozenie: () => marginesy().goraMm,
    granice: () => ({
      dolnaMm: 0,
      gornaMm: Math.max(0, wysokoscKartkiMm(stronaBiezaca) - marginesy().dolMm - NAJMNIEJSZE_POLE_MM),
    }),
    zPolozenia: (naKartce) => naKartce,
    ustaw: (milimetry) => czynnosci.naMargines('gora', milimetry),
  });

  const chwytDolu = utworzChwyt({
    kod: 'margines-dol',
    nazwa: 'margines dolny',
    wartosc: () => marginesy().dolMm,
    polozenie: () => wysokoscKartkiMm(stronaBiezaca) - marginesy().dolMm,
    granice: () => ({
      dolnaMm: 0,
      gornaMm: Math.max(
        0,
        wysokoscKartkiMm(stronaBiezaca) - marginesy().goraMm - NAJMNIEJSZE_POLE_MM,
      ),
    }),
    zPolozenia: (naKartce) => wysokoscKartkiMm(stronaBiezaca) - naKartce,
    ustaw: (milimetry) => czynnosci.naMargines('dol', milimetry),
  });

  const polaMarginesow = document.createElement('div');
  polaMarginesow.className = 'ms-linijka__marginesy';

  const znacznikKursora = document.createElement('div');
  znacznikKursora.className = 'ms-linijka__kursor';
  znacznikKursora.hidden = true;

  chwyty.append(polaMarginesow, znacznikKursora, chwytGory, chwytDolu);

  function przerysuj(): void {
    pas.style.height = `${punkty(wysokoscKartkiMm(stronaBiezaca))}px`;
    const { goraMm, dolMm } = marginesy();

    podzialka.replaceChildren();
    for (const kreska of zlozPodzialkeLinijki(wysokoscPola(), stronaBiezaca.jednostka, skala)) {
      const znak = document.createElement('span');
      znak.className = 'ms-linijka__kreska';
      znak.dataset['kreska'] = kreska.rodzaj;
      znak.style.top = `${punkty(goraMm + kreska.milimetry)}px`;
      if (kreska.napis !== '') znak.textContent = kreska.napis;
      podzialka.append(znak);
    }

    polaMarginesow.style.setProperty('--ms-linijka-gora', `${punkty(goraMm)}px`);
    polaMarginesow.style.setProperty('--ms-linijka-dol', `${punkty(dolMm)}px`);

    for (const wpis of zlozone) {
      const wartosc = wpis.opis.wartosc();
      wpis.element.style.top = `${punkty(wpis.opis.polozenie())}px`;
      wpis.element.setAttribute('aria-valuenow', String(Math.round(wartosc * 10) / 10));
      const granice = wpis.opis.granice();
      wpis.element.setAttribute('aria-valuemin', String(Math.round(granice.dolnaMm * 10) / 10));
      wpis.element.setAttribute('aria-valuemax', String(Math.round(granice.gornaMm * 10) / 10));
      wpis.element.setAttribute(
        'aria-valuetext',
        `${wpis.opis.nazwa}: ${opiszDlugosc(wartosc, stronaBiezaca.jednostka)}`,
      );
      wpis.element.title =
        `${wpis.opis.nazwa} — ${opiszDlugosc(wartosc, stronaBiezaca.jednostka)}. Przeciągnij albo ` +
        'przesuń strzałkami; ze Shiftem krok dziesięciokrotny.';
    }

    if (kursorMm === null) {
      znacznikKursora.hidden = true;
      return;
    }
    znacznikKursora.hidden = false;
    znacznikKursora.style.top = `${punkty(goraMm + kursorMm)}px`;
    znacznikKursora.title = `Kursor na ${opiszDlugosc(kursorMm, stronaBiezaca.jednostka)} od marginesu górnego.`;
  }

  przerysuj();

  return {
    element,

    ustawStrone(nowa) {
      stronaBiezaca = nowa;
      skala = nowa.skala / 100;
      przerysuj();
    },

    ustawSkale(procent) {
      skala = procent / 100;
      przerysuj();
    },

    ustawKursor(milimetry) {
      kursorMm = milimetry;
      przerysuj();
    },

    ustawWidocznosc(widoczna) {
      element.hidden = !widoczna;
    },

    widoczna: () => !element.hidden,

    opis() {
      const { goraMm, dolMm } = marginesy();
      return (
        `margines górny ${opiszDlugosc(goraMm, stronaBiezaca.jednostka)} · dolny ` +
        `${opiszDlugosc(dolMm, stronaBiezaca.jednostka)} · pole pisania ` +
        `${opiszDlugosc(wysokoscPola(), stronaBiezaca.jednostka)}`
      );
    },
  };
}
