/**
 * Dziennik przebiegów, czyli przeglądarka logu okna Execution Monitor. Wiersz
 * logu przychodzi w zdarzeniu `automation.execution.status` i przychodzi raz,
 * więc dziennik zbiera te wiersze w chwili, w której się pojawiają.
 */
import type { AutomationExecutionStatusEvent } from '../../../../shared/contract';

/**
 * Ile wierszy dziennik przechowuje. Granica istnieje, ponieważ przebieg
 * wieloetapowy nadaje wiersz na każdy etap i na każdy obieg naprawczy, a okno
 * bywa otwarte godzinami. Wiersze najstarsze odchodzą pierwsze.
 */
const GRANICA_WIERSZY = 500;

/**
 * Jeden wiersz dziennika wraz z tym, co o nim wiadomo z kontraktu: czas
 * przyjęcia w oknie, identyfikator przebiegu, stan przebiegu z chwili wpisu,
 * nazwa etapu oraz treść wiersza.
 */
export interface WpisDziennika {
  /** Czas przyjęcia wiersza w oknie, w milisekundach epoki. */
  czas: number;
  idPrzebiegu: string;
  /** Stan przebiegu w chwili wpisu — po nim idzie zawężanie dziennika. */
  stan: string;
  /** Nazwa etapu bieżącego, gdy zdarzenie ją niosło. */
  etap: string;
  tresc: string;
}

export interface DziennikPrzebiegow {
  /** Przyjmuje zdarzenie; wiersz bez treści logu nie zostawia wpisu. */
  przyjmij(zdarzenie: AutomationExecutionStatusEvent): void;
  /** Wpisy w kolejności przyjęcia, opcjonalnie zawężone do jednego stanu. */
  wpisy(stan?: string): WpisDziennika[];
  /** Ile wierszy odrzucono po przekroczeniu granicy bufora. */
  odrzucone(): number;
  /** Opróżnia dziennik — po zmianie automatyki bieżącej. */
  wyczysc(): void;
}

export function utworzDziennikPrzebiegow(): DziennikPrzebiegow {
  const wpisy: WpisDziennika[] = [];
  let odrzucone = 0;

  return {
    przyjmij(zdarzenie) {
      const tresc = zdarzenie.logLine ?? '';
      if (tresc.trim() === '') return;
      wpisy.push({
        czas: Date.now(),
        idPrzebiegu: zdarzenie.execution.id,
        stan: zdarzenie.execution.status,
        etap: zdarzenie.stepLabel ?? '',
        tresc,
      });
      while (wpisy.length > GRANICA_WIERSZY) {
        wpisy.shift();
        odrzucone += 1;
      }
    },

    wpisy: (stan) =>
      stan === undefined || stan === ''
        ? [...wpisy]
        : wpisy.filter((wpis) => wpis.stan === stan),

    odrzucone: () => odrzucone,

    wyczysc() {
      wpisy.length = 0;
      odrzucone = 0;
    },
  };
}

/**
 * Wiersz dziennika w zapisie tekstowym, będący zarazem treścią pliku wydania
 * i pozycją wykazu w oknie. Zapis niesie czas przyjęcia, nazwę etapu, gdy
 * zdarzenie ją niosło, oraz treść wiersza.
 */
export function zapisWpisu(wpis: WpisDziennika): string {
  const czas = new Date(wpis.czas).toLocaleString('pl-PL');
  const etap = wpis.etap === '' ? '' : ` [${wpis.etap}]`;
  return `${czas} ${wpis.idPrzebiegu} (${wpis.stan})${etap}: ${wpis.tresc}`;
}

/**
 * Zdanie nad dziennikiem: mówi, skąd pochodzą wiersze i czego w dzienniku nie
 * ma, żeby Operator nie wziął zbioru wierszy okna otwartego za pełny zapis
 * przebiegu automatyki.
 */
export function zdanieODzienniku(wpisow: number, odrzuconych: number): string {
  if (wpisow === 0) {
    return (
      'Dziennik jest pusty. Wiersze zbierają się ze zdarzeń stanu przebiegu, więc pojawią się ' +
      'dopiero przy najbliższej zmianie stanu. Zapisu sprzed otwarcia tego okna dziennik nie ma — ' +
      'ten sięga po pozycję „Log przebiegu z rdzenia” w pasku akcji.'
    );
  }
  const dopowiedzenie =
    odrzuconych === 0
      ? ''
      : ` Najstarszych wierszy odrzucono ${odrzuconych} po przekroczeniu granicy bufora okna.`;
  return (
    `Dziennik zebrał ${wpisow} wierszy ze zdarzeń stanu przebiegu odebranych przez to okno.` +
    dopowiedzenie
  );
}
