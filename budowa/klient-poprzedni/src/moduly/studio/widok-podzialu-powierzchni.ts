import {
  przytnijUdzialPodzialu,
  type KierunekPodzialu,
  type TrybDwochDokumentow,
} from './widok-nastawy-operatora';

/** Typ NumerPola nazywa numer pola podziału powierzchni dokumentu: pole pierwsze albo pole drugie roboczej powierzchni. */
export type NumerPola = 1 | 2;

/** Interfejs PodzialPowierzchni niesie podział powierzchni na dwa pola wraz z jego sterowaniem: trybem, kierunkiem, udziałem i polem czynnym. */
export interface PodzialPowierzchni {
  element: HTMLElement;
  /** Wstawia zawartość obu pól; wołane raz, przy składaniu okna. */
  ustawPola(pierwsze: HTMLElement, drugie: HTMLElement): void;
  /** Gniazdo pod pasek statusu i wiersz polecenia pola — każde pole ma własne. */
  gniazdoStopki(numer: NumerPola): HTMLElement;
  /** Przestawia tryb; zakładki chowają pole drugie, nie rozbierają go. */
  ustawTryb(tryb: TrybDwochDokumentow): void;
  tryb(): TrybDwochDokumentow;
  ustawKierunek(kierunek: KierunekPodzialu): void;
  kierunek(): KierunekPodzialu;
  /** Udział pola pierwszego, od 0,15 do 0,85. */
  ustawUdzial(udzial: number): void;
  udzial(): number;
  /** Nazwy pól — plakietka i odczyt ekranowy. */
  ustawNazwy(pierwsza: string, druga: string): void;
  /** Pole czynne — to, w którym stoi praca; w zakładkach zawsze pierwsze. */
  ustawCzynne(numer: NumerPola): void;
  czynne(): NumerPola;
  /** Zdanie o podziale — do paska stanu. */
  opis(): string;
}

/** Interfejs CzynnosciPodzialu niesie czynności podziału powierzchni zgłaszane oknu: zmianę udziału granicy i zmianę pola czynnego. */
export interface CzynnosciPodzialu {
  /** Granica przestawiona — okno zapamiętuje udział jako nastawę Operatora. */
  naUdzial(udzial: number): void;
  /** Pole czynne zmienione naciśnięciem w polu. */
  naCzynne(numer: NumerPola): void;
}

/** Stała SKOK_GRANICY wyznacza skok granicy podziału na jedno naciśnięcie strzałki, wyrażony w częściach udziału. */
const SKOK_GRANICY = 0.02;

export function utworzPodzialPowierzchni(
  czynnosci: CzynnosciPodzialu,
  trybPoczatkowy: TrybDwochDokumentow = 'zakladki',
  kierunekPoczatkowy: KierunekPodzialu = 'pionowy',
  udzialPoczatkowy = 0.5,
): PodzialPowierzchni {
  let trybBiezacy = trybPoczatkowy;
  let kierunekBiezacy = kierunekPoczatkowy;
  let udzialBiezacy = przytnijUdzialPodzialu(udzialPoczatkowy);
  let czynneBiezace: NumerPola = 1;

  const element = document.createElement('div');
  element.className = 'ms-podzial';
  element.dataset['tryb'] = trybBiezacy;
  element.dataset['kierunek'] = kierunekBiezacy;

  const pola = new Map<NumerPola, { pole: HTMLElement; tresc: HTMLElement; stopka: HTMLElement; plakietka: HTMLElement }>();

  for (const numer of [1, 2] as const) {
    const plakietka = document.createElement('span');
    plakietka.className = 'dn-plakietka ms-podzial__plakietka';
    plakietka.textContent = numer === 1 ? 'pole pierwsze' : 'pole drugie';

    const tresc = document.createElement('div');
    tresc.className = 'ms-podzial__tresc';

    const stopka = document.createElement('div');
    stopka.className = 'ms-podzial__stopka';

    const pole = document.createElement('section');
    pole.className = 'ms-podzial__pole';
    pole.dataset['pole'] = String(numer);
    pole.dataset['czynne'] = numer === 1 ? 'tak' : 'nie';
    pole.setAttribute('aria-label', numer === 1 ? 'Pole pierwsze podziału' : 'Pole drugie podziału');
    pole.append(plakietka, tresc, stopka);
    // Ognisko w polu czyni je czynnym: kliknięcie w drugie pismo przenosi pracę na drugie pismo.
    pole.addEventListener('focusin', () => ustawCzynne(numer));
    pole.addEventListener('mousedown', () => ustawCzynne(numer));
    pola.set(numer, { pole, tresc, stopka, plakietka });
  }

  const granica = document.createElement('div');
  granica.className = 'ms-podzial__granica';
  granica.setAttribute('role', 'separator');
  granica.tabIndex = 0;
  granica.setAttribute('aria-orientation', kierunekBiezacy === 'pionowy' ? 'vertical' : 'horizontal');
  granica.setAttribute('aria-valuemin', '15');
  granica.setAttribute('aria-valuemax', '85');
  granica.title =
    'Granica podziału. Przeciągnij albo przesuń strzałkami o dwa punkty procentowe; Home i End ' +
    'dosuwają do granic pola pracy, Enter wraca na pół.';

  const pierwsze = pola.get(1);
  const drugie = pola.get(2);
  if (pierwsze !== undefined && drugie !== undefined) {
    element.append(pierwsze.pole, granica, drugie.pole);
  }

  /* ── Granica: chwyt i klawiatura ─────────────────────────────────────────── */

  function udzialZPolozenia(zdarzenie: { clientX: number; clientY: number }): number {
    const prostokat = element.getBoundingClientRect();
    const dlugosc = kierunekBiezacy === 'pionowy' ? prostokat.width : prostokat.height;
    if (dlugosc <= 0) return udzialBiezacy;
    const odsuniecie =
      kierunekBiezacy === 'pionowy'
        ? zdarzenie.clientX - prostokat.left
        : zdarzenie.clientY - prostokat.top;
    return odsuniecie / dlugosc;
  }

  let wleczenie = false;

  granica.addEventListener('mousedown', (zdarzenie) => {
    zdarzenie.preventDefault();
    wleczenie = true;
    granica.dataset['wleczenie'] = 'tak';
  });

  // Nasłuch stoi na dokumencie, nie na granicy, bo mysz wyprzedza kreskę o kilka punktów.
  document.addEventListener('mousemove', (zdarzenie) => {
    if (!wleczenie) return;
    przestawUdzial(udzialZPolozenia(zdarzenie), true);
  });

  document.addEventListener('mouseup', () => {
    if (!wleczenie) return;
    wleczenie = false;
    delete granica.dataset['wleczenie'];
  });

  granica.addEventListener('keydown', (zdarzenie) => {
    const mniej = zdarzenie.key === 'ArrowLeft' || zdarzenie.key === 'ArrowUp';
    const wiecej = zdarzenie.key === 'ArrowRight' || zdarzenie.key === 'ArrowDown';
    if (mniej || wiecej) {
      zdarzenie.preventDefault();
      przestawUdzial(udzialBiezacy + (wiecej ? SKOK_GRANICY : -SKOK_GRANICY), true);
      return;
    }
    if (zdarzenie.key === 'Home') {
      zdarzenie.preventDefault();
      przestawUdzial(0, true);
      return;
    }
    if (zdarzenie.key === 'End') {
      zdarzenie.preventDefault();
      przestawUdzial(1, true);
      return;
    }
    if (zdarzenie.key === 'Enter' || zdarzenie.key === ' ') {
      zdarzenie.preventDefault();
      przestawUdzial(0.5, true);
    }
  });

  function przestawUdzial(udzial: number, zglos: boolean): void {
    const przyciety = przytnijUdzialPodzialu(udzial);
    if (przyciety === udzialBiezacy && !zglos) return;
    udzialBiezacy = przyciety;
    przypnijUdzial();
    if (zglos) czynnosci.naUdzial(udzialBiezacy);
  }

  function przypnijUdzial(): void {
    element.style.setProperty('--ms-podzial-udzial', String(udzialBiezacy));
    granica.setAttribute('aria-valuenow', String(Math.round(udzialBiezacy * 100)));
    granica.setAttribute(
      'aria-valuetext',
      `pole pierwsze zajmuje ${Math.round(udzialBiezacy * 100)} % powierzchni`,
    );
  }

  function ustawCzynne(numer: NumerPola): void {
    const docelowy = trybBiezacy === 'zakladki' ? 1 : numer;
    if (docelowy === czynneBiezace) return;
    czynneBiezace = docelowy;
    for (const [klucz, wpis] of pola) {
      wpis.pole.dataset['czynne'] = klucz === czynneBiezace ? 'tak' : 'nie';
    }
    czynnosci.naCzynne(czynneBiezace);
  }

  przypnijUdzial();
  // W trybie zakładek pole drugie jest schowane od początku, bo daje większe pole pracy.
  if (drugie !== undefined) drugie.pole.hidden = trybBiezacy === 'zakladki';
  granica.hidden = trybBiezacy === 'zakladki';

  return {
    element,

    ustawPola(pierwszeTresc, drugieTresc) {
      pola.get(1)?.tresc.replaceChildren(pierwszeTresc);
      pola.get(2)?.tresc.replaceChildren(drugieTresc);
    },

    gniazdoStopki(numer) {
      const wpis = pola.get(numer);
      if (wpis !== undefined) return wpis.stopka;
      // Numer poza zakresem nie zdarza się przy typie NumerPola; element zapasowy zastępuje wyjątek.
      return document.createElement('div');
    },

    ustawTryb(tryb) {
      trybBiezacy = tryb;
      element.dataset['tryb'] = tryb;
      const wpis = pola.get(2);
      if (wpis !== undefined) wpis.pole.hidden = tryb === 'zakladki';
      granica.hidden = tryb === 'zakladki';
      if (tryb === 'zakladki') ustawCzynne(1);
    },

    tryb: () => trybBiezacy,

    ustawKierunek(kierunek) {
      kierunekBiezacy = kierunek;
      element.dataset['kierunek'] = kierunek;
      granica.setAttribute('aria-orientation', kierunek === 'pionowy' ? 'vertical' : 'horizontal');
    },

    kierunek: () => kierunekBiezacy,

    ustawUdzial(udzial) {
      przestawUdzial(udzial, false);
    },

    udzial: () => udzialBiezacy,

    ustawNazwy(pierwszaNazwa, drugaNazwa) {
      const jedno = pola.get(1);
      const dwa = pola.get(2);
      if (jedno !== undefined) {
        jedno.plakietka.textContent = pierwszaNazwa;
        jedno.pole.setAttribute('aria-label', `Pole pierwsze podziału: ${pierwszaNazwa}`);
      }
      if (dwa !== undefined) {
        dwa.plakietka.textContent = drugaNazwa;
        dwa.pole.setAttribute('aria-label', `Pole drugie podziału: ${drugaNazwa}`);
      }
    },

    ustawCzynne,
    czynne: () => czynneBiezace,

    opis() {
      if (trybBiezacy === 'zakladki') {
        return (
          'tryb zakładek — jeden dokument na całej powierzchni; dokument pola drugiego zostaje ' +
          'otwarty i dostępny modelowi'
        );
      }
      return (
        `podział ${kierunekBiezacy === 'pionowy' ? 'pionowy' : 'poziomy'} · pole pierwsze ` +
        `${Math.round(udzialBiezacy * 100)} % · czynne pole ${czynneBiezace}`
      );
    },
  };
}
