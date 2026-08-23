/**
 * Pływający awatar Always On Display — warstwa 1 funkcji globalnej
 * (`docs/funkcje-globalne/always-on-display.md`, rozdz. 2.4 i 8.1).
 *
 * Awatar jest CAŁĄ powierzchnią funkcji w stanie spoczynku: pojedyncze koło
 * przy prawej krawędzi obszaru roboczego, bez etykiety, bez ramki kontenera,
 * ponad całą powłoką aplikacji. Nie znika przy przełączeniu środowiska ani
 * modułu, bo warstwa, w której siedzi, leży poza obszarem podmienianym przez
 * moduł.
 *
 * Siedem stanów rozdz. 2.4 ma tu siedem wartości `data-stan`. Żaden z nich nie
 * jest niesiony samą barwą: przy każdym stoi etykieta dostępności i tytuł, a
 * stany „sugestia oczekująca" i „waga wysoka" niosą dodatkowo plakietkę
 * liczbową. To wymóg `design/KANON.md` — stan nigdy samym kolorem.
 *
 * Awatar nie zna kanału, kolejki ani dymka. Przyjmuje dwa wywołania zwrotne
 * (kliknięcie, kliknięcie podwójne) i cztery czynności nastawcze; co za nimi
 * stoi, rozstrzyga `warstwa-aod.ts`.
 */

/** Stan awatara — rozdz. 2.4, kolumna „Stan awatara". */
export const StanAwatara = {
  /** Brak oczekującej sugestii, gotowość do interakcji. */
  Spoczynek: 'spoczynek',
  /** Jedna lub więcej sugestii oczekuje na zapoznanie. */
  SugestiaOczekujaca: 'sugestia-oczekujaca',
  /** Zdarzenie wymagające decyzji Operatora — plakietka pulsująca. */
  WagaWysoka: 'waga-wysoka',
  /** Trwa synteza mowy albo strumieniowanie odpowiedzi. */
  MowiLubPisze: 'mowi-lub-pisze',
  /** Trwa odbiór polecenia głosowego. */
  Nasluch: 'nasluch',
  /** Wyciszenie sugestii aktywne. */
  Wyciszony: 'wyciszony',
  /** Awatar poza polem widzenia, funkcja czynna w tle. */
  Ukryty: 'ukryty',
} as const;
export type StanAwatara = (typeof StanAwatara)[keyof typeof StanAwatara];

/** Zdanie opisujące stan — czytane przez technologie wspomagające i przy najechaniu. */
const OPISY_STANOW: Readonly<Record<StanAwatara, string>> = {
  [StanAwatara.Spoczynek]: 'Always On Display — spoczynek, żadna sugestia nie czeka',
  [StanAwatara.SugestiaOczekujaca]: 'Always On Display — sugestie oczekują na zapoznanie',
  [StanAwatara.WagaWysoka]: 'Always On Display — sugestia o wadze wysokiej czeka na decyzję',
  [StanAwatara.MowiLubPisze]: 'Always On Display — trwa odpowiedź',
  [StanAwatara.Nasluch]: 'Always On Display — trwa nasłuch polecenia głosowego',
  [StanAwatara.Wyciszony]: 'Always On Display — wyciszony, sugestie gromadzą się bez ujawniania',
  [StanAwatara.Ukryty]: 'Always On Display — awatar ukryty, funkcja czynna w tle',
};

export interface OpisAwatara {
  /** Kliknięcie pojedyncze — otwiera dymek kontekstowy z bieżącą sugestią. */
  naKlikniecie(): void;
  /** Kliknięcie podwójne — otwiera powierzchnię interakcji. */
  naKlikniecePodwojne(): void;
}

export interface AwatarAod {
  element: HTMLElement;
  /** Nadaje stan awatara wraz z jego opisem. */
  ustawStan(stan: StanAwatara): void;
  /**
   * Nadaje liczbę oczekujących sugestii.
   *
   * Zero chowa plakietkę — rozdz. 2.6, stan „ukryta (zero)". Plakietka ukryta
   * ustawieniem wyciszenia idzie osobno przez {@link ustawPlakietkeWidoczna}.
   */
  ustawLiczbe(liczba: number): void;
  /** Chowa albo pokazuje plakietkę niezależnie od liczby (wyciszenie, rozdz. 3.5). */
  ustawPlakietkeWidoczna(widoczna: boolean): void;
  /** Przesuwa awatar o szerokość otwartej kolumny, żeby nie nachodził na jej treść. */
  ustawOdsuniecie(pikseli: number): void;
  /** Zdejmuje nasłuch zdarzeń wskaźnika. */
  rozlacz(): void;
}

/**
 * Odstęp, po którym kliknięcie pojedyncze uznajemy za pojedyncze.
 *
 * Przeglądarka wysyła `click` przed `dblclick`, więc bez tego odstępu każde
 * kliknięcie podwójne otwierałoby najpierw dymek. Odstęp jest krótszy niż
 * czas reakcji na otwarty dymek, więc pojedyncze kliknięcie nadal działa od ręki.
 */
const ODSTEP_PODWOJNEGO_MS = 250;

export function utworzAwatarAod(opis: OpisAwatara): AwatarAod {
  const element = document.createElement('button');
  element.type = 'button';
  element.className = 'ao-awatar';
  element.dataset['stan'] = StanAwatara.Spoczynek;

  // Koło awatara. Treść jest znakiem funkcji, nie obrazem z sieci — warstwa
  // wizualna nie ma dziś rysunku awatara, a kwadrat zastępczy byłby atrapą.
  const znak = document.createElement('span');
  znak.className = 'ao-awatar__znak';
  znak.setAttribute('aria-hidden', 'true');
  znak.textContent = 'AOD';

  const plakietka = document.createElement('span');
  plakietka.className = 'dn-plakietka ao-awatar__plakietka';
  plakietka.hidden = true;

  element.append(znak, plakietka);

  /** Liczba oczekujących sugestii z ostatniego nadania. */
  let liczba = 0;
  /** Czy plakietka wolno się pokazać — wyciszenie ją chowa mimo liczby. */
  let plakietkaDozwolona = true;

  function odswiezPlakietke(): void {
    const widoczna = plakietkaDozwolona && liczba > 0;
    plakietka.hidden = !widoczna;
    plakietka.textContent = widoczna ? String(liczba) : '';
  }

  /** Uchwyt oczekującego kliknięcia pojedynczego; `null`, gdy żadne nie czeka. */
  let oczekujace: ReturnType<typeof setTimeout> | null = null;

  function anulujOczekujace(): void {
    if (oczekujace === null) return;
    clearTimeout(oczekujace);
    oczekujace = null;
  }

  function naKlikniecie(): void {
    anulujOczekujace();
    oczekujace = setTimeout(() => {
      oczekujace = null;
      opis.naKlikniecie();
    }, ODSTEP_PODWOJNEGO_MS);
  }

  function naPodwojne(): void {
    anulujOczekujace();
    opis.naKlikniecePodwojne();
  }

  element.addEventListener('click', naKlikniecie);
  element.addEventListener('dblclick', naPodwojne);

  function ustawStan(stan: StanAwatara): void {
    element.dataset['stan'] = stan;
    element.setAttribute('aria-label', OPISY_STANOW[stan]);
    element.title = OPISY_STANOW[stan];
    // Tryb ukryty zabiera awatar z pola widzenia, ale funkcja pracuje dalej —
    // dostęp zostaje skrótem klawiszowym i listwą ustawień (rozdz. 9.2).
    element.hidden = stan === StanAwatara.Ukryty;
  }

  ustawStan(StanAwatara.Spoczynek);

  return {
    element,

    ustawStan,

    ustawLiczbe(nowa) {
      liczba = Math.max(0, nowa);
      odswiezPlakietke();
    },

    ustawPlakietkeWidoczna(widoczna) {
      plakietkaDozwolona = widoczna;
      odswiezPlakietke();
    },

    ustawOdsuniecie(pikseli) {
      element.style.setProperty('--ao-odsuniecie', `${Math.max(0, pikseli)}px`);
    },

    rozlacz() {
      anulujOczekujace();
      element.removeEventListener('click', naKlikniecie);
      element.removeEventListener('dblclick', naPodwojne);
      element.remove();
    },
  };
}
