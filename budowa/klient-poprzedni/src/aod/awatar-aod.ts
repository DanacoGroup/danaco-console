/** Siedem stanów awatara jako siedem wartości `data-stan`, z którymi warstwa `warstwa-aod.ts` steruje wyglądem. */
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

/** Zdanie opisujące bieżący stan awatara, czytane przez technologie wspomagające i pokazywane przy najechaniu. */
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
  /** Nadaje liczbę oczekujących sugestii; zero chowa plakietkę niezależnie od wyciszenia. */
  ustawLiczbe(liczba: number): void;
  /** Chowa albo pokazuje plakietkę niezależnie od liczby oczekujących sugestii — sterowane wyciszeniem. */
  ustawPlakietkeWidoczna(widoczna: boolean): void;
  /** Przesuwa awatar o szerokość otwartej kolumny, żeby nie nachodził na jej treść. */
  ustawOdsuniecie(pikseli: number): void;
  /** Zdejmuje nasłuch zdarzeń wskaźnika. */
  rozlacz(): void;
}

/** Odstęp, po którym kliknięcie pojedyncze uznaje się za pojedyncze, a nie za początek kliknięcia podwójnego. */
const ODSTEP_PODWOJNEGO_MS = 250;

export function utworzAwatarAod(opis: OpisAwatara): AwatarAod {
  const element = document.createElement('button');
  element.type = 'button';
  element.className = 'ao-awatar';
  element.dataset['stan'] = StanAwatara.Spoczynek;

  // Koło awatara niesie znak funkcji, nie obraz z sieci — kwadrat zastępczy byłby atrapą.
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
    // Tryb ukryty zabiera awatar z pola widzenia, funkcja pracuje dalej w tle.
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
