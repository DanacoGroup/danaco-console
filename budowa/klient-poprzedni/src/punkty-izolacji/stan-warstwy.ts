import { IsolationLayer } from '../../../shared/contract';

/**
 * Warstwa izolacji czynna w oknie — jedna dla wszystkich obszarów.
 *
 * Warstwy są dwie, ale nie są dwoma trybami do wyboru. `default` to warstwa
 * bazowa platformy: obowiązuje przy każdej nowej sesji, karcie, roli, projekcie
 * i oknie. `session` nakłada się na bazową bez jej zmiany i wygasa z zamknięciem
 * karty sesji — to podkład i naklejka na nim, nie dwa osobne magazyny.
 *
 * Warstwa jest stanem wspólnym, a nie polem jednego obszaru, bo dotyczy każdego
 * odczytu i każdego zapisu izolacji: kontekstu, zakresu technicznego,
 * przypisania profilu i podglądu polityki. Gdyby siedziała w jednej zakładce,
 * pozostałe pytałyby rdzeń o coś innego, niż widać na pasku — dlatego jej
 * kontrolka jest w narzędziach ramy okna, a stan tutaj.
 *
 * Ten plik nie woła rdzenia. Przełączenie warstwy w rdzeniu
 * (`isolation.layer.set`) należy do `sterowanie-warstwa.ts`; tutaj zostaje
 * wyłącznie to, co okno wie o warstwie po odpowiedzi rdzenia, i powiadomienie
 * obszarów o zmianie.
 */
export interface StanWarstwy {
  /** Warstwa czynna; przed pierwszym przełączeniem — bazowa platformy. */
  warstwa(): IsolationLayer;
  /** Zapisuje warstwę potwierdzoną przez rdzeń i powiadamia obszary. */
  ustaw(warstwa: IsolationLayer): void;
  /** Subskrypcja zmiany warstwy; zwraca odsubskrybowanie. */
  naZmiane(sluchacz: (warstwa: IsolationLayer) => void): () => void;
}

/** Nazwy obu warstw po polsku — jedno źródło dla paska i dla zdań obszarów. */
export const NAZWY_WARSTW: Readonly<Record<IsolationLayer, string>> = {
  [IsolationLayer.Default]: 'Domyślna platformy',
  [IsolationLayer.Session]: 'Karty sesji',
};

/** Zdanie o skutku wyboru warstwy — skutek jest ważniejszy niż nazwa. */
export const OPISY_WARSTW: Readonly<Record<IsolationLayer, string>> = {
  [IsolationLayer.Default]:
    'Warstwa bazowa. Zapis obowiązuje przy każdej nowej sesji, karcie, roli, projekcie i oknie — ' +
    'zostaje po zamknięciu karty.',
  [IsolationLayer.Session]:
    'Warstwa karty sesji. Nakłada się na bazową bez zmiany jej wartości i wygasa z zamknięciem ' +
    'karty — po ponownym otwarciu obowiązuje znowu warstwa domyślna.',
};

export function utworzStanWarstwy(poczatkowa: IsolationLayer = IsolationLayer.Default): StanWarstwy {
  let czynna: IsolationLayer = poczatkowa;
  const sluchacze = new Set<(warstwa: IsolationLayer) => void>();

  return {
    warstwa: () => czynna,

    ustaw(warstwa) {
      if (czynna === warstwa) return;
      czynna = warstwa;
      for (const sluchacz of [...sluchacze]) sluchacz(warstwa);
    },

    naZmiane(sluchacz) {
      sluchacze.add(sluchacz);
      return () => sluchacze.delete(sluchacz);
    },
  };
}
