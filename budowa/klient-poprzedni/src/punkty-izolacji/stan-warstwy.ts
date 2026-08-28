import { IsolationLayer } from '../../../shared/contract';

/** Warstwa izolacji czynna w oknie — jedna wspólna dla wszystkich obszarów tego okna Punktów Izolacji klienta. */
export interface StanWarstwy {
  /** Warstwa czynna; przed pierwszym przełączeniem — bazowa platformy. */
  warstwa(): IsolationLayer;
  /** Zapisuje warstwę potwierdzoną przez rdzeń i powiadamia obszary. */
  ustaw(warstwa: IsolationLayer): void;
  /** Subskrypcja zmiany warstwy; zwraca odsubskrybowanie. */
  naZmiane(sluchacz: (warstwa: IsolationLayer) => void): () => void;
}

/** Nazwy obu warstw izolacji po polsku — jedno wspólne źródło nazw dla paska narzędzi i dla zdań obszarów. */
export const NAZWY_WARSTW: Readonly<Record<IsolationLayer, string>> = {
  [IsolationLayer.Default]: 'Domyślna platformy',
  [IsolationLayer.Session]: 'Karty sesji',
};

/** Zdanie o skutku wyboru warstwy izolacji — skutek tego wyboru jest tu ważniejszy niż sama jego nazwa. */
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
