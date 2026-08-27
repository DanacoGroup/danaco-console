import './uchwyt.css';

/**
 * Uchwyt ręcznego ustawiania szerokości kolumny sąsiadującej, sterowany wskaźnikiem lub
 * klawiaturą, oznaczony rolą separatora i pełnym kompletem atrybutów dostępności ARIA
 * opisujących bieżącą wartość.
 */
export interface UchwytSzerokosci {
  element: HTMLElement;
  /** Ustawia wartość pokazywaną przez uchwyt (bez wołania `naZmiane`). */
  ustawWartosc(px: number): void;
}

/**
 * Zależności uchwytu szerokości: etykieta dostępności, wartość początkowa, dolna
 * granica stała oraz górna granica liczona dynamicznie w chwili ruchu wraz
 * z powiadomieniem o zmianie.
 */
export interface OpcjeUchwytu {
  /** Etykieta dostępności — co ten uchwyt rozsuwa. */
  etykieta: string;
  wartosc: number;
  minimum: number;
  /** Górna granica liczona w chwili ruchu — scena zmienia szerokość. */
  maksimum(): number;
  naZmiane(px: number): void;
}

/**
 * Krok jednego naciśnięcia strzałki w pikselach, ustawiony na wielokrotność czterech
 * pikseli siatki, żeby ruch klawiaturą pozostawał praktyczny i precyzyjny zarazem.
 */
const KROK_PX = 16;

/**
 * Tworzy uchwyt rozsuwający kolumnę sąsiadującą, obsługujący ruch wskaźnika oraz
 * klawiatury, i oddający wołającemu bieżącą szerokość wyłącznie w pikselach.
 */
export function utworzUchwytSzerokosci(opcje: OpcjeUchwytu): UchwytSzerokosci {
  const element = document.createElement('div');
  element.className = 'dn-uchwyt';
  element.setAttribute('role', 'separator');
  element.setAttribute('aria-orientation', 'vertical');
  element.setAttribute('aria-label', opcje.etykieta);
  element.tabIndex = 0;

  let wartosc = opcje.wartosc;

  /** Wartość sprowadzona do granic bieżących; górna liczona w chwili pytania. */
  function wGranicach(px: number): number {
    if (!Number.isFinite(px)) return opcje.minimum;
    const gora = Math.max(opcje.minimum, opcje.maksimum());
    return Math.round(Math.min(gora, Math.max(opcje.minimum, px)));
  }

  /** Zapisuje wartość i opowiada o niej czytnikowi ekranu. */
  function pokaz(px: number): void {
    wartosc = px;
    element.setAttribute('aria-valuenow', String(px));
    element.setAttribute('aria-valuemin', String(opcje.minimum));
    element.setAttribute('aria-valuemax', String(Math.max(opcje.minimum, opcje.maksimum())));
  }

  /** Zmiana wychodząca na zewnątrz — wołana wyłącznie przy ruchu uchwytem. */
  function przestaw(px: number): void {
    const nowa = wGranicach(px);
    if (nowa === wartosc) {
      // Wartość niezmieniona, ale opis odświeżany, bo górna granica mogła się zmienić.
      pokaz(nowa);
      return;
    }
    pokaz(nowa);
    opcje.naZmiane(nowa);
  }

  // Przechwycenie wskaźnika na pointerdown, żeby ruch poza uchwytem nadal do niego trafiał.
  let poczatek: { x: number; wartosc: number } | null = null;

  element.addEventListener('pointerdown', (zdarzenie) => {
    poczatek = { x: zdarzenie.clientX, wartosc };
    element.setPointerCapture(zdarzenie.pointerId);
    element.dataset['przeciagany'] = 'tak';
    // Bez tego przeciąganie zaznacza tekst obu kolumn zamiast rozsuwać.
    zdarzenie.preventDefault();
  });

  element.addEventListener('pointermove', (zdarzenie) => {
    if (poczatek === null) return;
    przestaw(poczatek.wartosc - (zdarzenie.clientX - poczatek.x));
  });

  function koniecRuchu(): void {
    poczatek = null;
    delete element.dataset['przeciagany'];
  }

  element.addEventListener('pointerup', koniecRuchu);
  element.addEventListener('pointercancel', koniecRuchu);

  element.addEventListener('keydown', (zdarzenie) => {
    switch (zdarzenie.key) {
      case 'ArrowLeft':
        przestaw(wartosc + KROK_PX);
        break;
      case 'ArrowRight':
        przestaw(wartosc - KROK_PX);
        break;
      case 'Home':
        przestaw(opcje.minimum);
        break;
      case 'End':
        przestaw(opcje.maksimum());
        break;
      default:
        return;
    }
    // Strzałka na uchwycie rozsuwa kolumnę, nie przewija toru okien.
    zdarzenie.preventDefault();
  });

  pokaz(wGranicach(opcje.wartosc));

  return {
    element,
    ustawWartosc(px: number): void {
      pokaz(wGranicach(px));
    },
  };
}
