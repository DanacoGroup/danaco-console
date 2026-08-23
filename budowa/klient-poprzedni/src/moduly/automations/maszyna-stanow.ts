import { AutomationExecutionStatus } from '../../../../shared/contract';
import { pozycjaWykazu, wykaz } from '../../modele/kontrolki-formularza';

/**
 * Maszyna stanów przebiegu — zakładka „Stany" okna Orchestratora.
 *
 * Opracowanie modułu żąda formalnego modelu stanów przebiegu z przejściami.
 * Model nie jest tu wymyślony: stany są kompletem wyliczenia
 * `AutomationExecutionStatus`, a przejścia wynikają z działań, które kontrakt
 * na przebiegu dopuszcza — uruchomienia, wstrzymania, wznowienia, zatrzymania
 * i zamknięcia przebiegu powodzeniem albo błędem.
 *
 * Widok jest wykazem, nie rysunkiem: przejść jest siedem, a rysunek siedmiu
 * strzałek nie mówi więcej niż siedem zdań i kosztuje drugą kanwę w module.
 *
 * Plik nie woła rdzenia; stan bieżący przychodzi z zewnątrz.
 */

/** Jedno przejście maszyny stanów wraz z jego przyczyną. */
export interface PrzejscieStanu {
  ze: string;
  na: string;
  /** Co powoduje przejście — działanie Operatora albo zamknięcie przebiegu. */
  przyczyna: string;
}

/** Komplet przejść przebiegu automatyki. */
export const PRZEJSCIA_PRZEBIEGU: readonly PrzejscieStanu[] = [
  {
    ze: AutomationExecutionStatus.Pending,
    na: AutomationExecutionStatus.Running,
    przyczyna: 'kolejka podjęła zlecenie przebiegu',
  },
  {
    ze: AutomationExecutionStatus.Running,
    na: AutomationExecutionStatus.Paused,
    przyczyna: 'wstrzymanie kolejki przebiegu',
  },
  {
    ze: AutomationExecutionStatus.Paused,
    na: AutomationExecutionStatus.Running,
    przyczyna: 'wznowienie kolejki przebiegu',
  },
  {
    ze: AutomationExecutionStatus.Running,
    na: AutomationExecutionStatus.Succeeded,
    przyczyna: 'wszystkie kroki zamknięte powodzeniem',
  },
  {
    ze: AutomationExecutionStatus.Running,
    na: AutomationExecutionStatus.Failed,
    przyczyna: 'krok zakończony błędem po wyczerpaniu obiegów naprawczych',
  },
  {
    ze: AutomationExecutionStatus.Running,
    na: AutomationExecutionStatus.Stopped,
    przyczyna: 'zatrzymanie przebiegu przez Operatora',
  },
  {
    ze: AutomationExecutionStatus.Failed,
    na: AutomationExecutionStatus.Running,
    przyczyna: 'bieg naprawczy podjęty ponownie',
  },
];

/** Nazwy stanów w mowie Operatora; klucze są wartościami kontraktu. */
const NAZWY_STANOW: Readonly<Record<string, string>> = {
  [AutomationExecutionStatus.Pending]: 'oczekuje',
  [AutomationExecutionStatus.Running]: 'w toku',
  [AutomationExecutionStatus.Paused]: 'wstrzymany',
  [AutomationExecutionStatus.Succeeded]: 'zakończony powodzeniem',
  [AutomationExecutionStatus.Failed]: 'zakończony błędem',
  [AutomationExecutionStatus.Stopped]: 'przerwany',
};

/** Nazwa stanu w mowie Operatora; stan nieznany zostaje przy swojej wartości. */
export function nazwaStanu(stan: string): string {
  return NAZWY_STANOW[stan] ?? stan;
}

/**
 * Buduje wykaz przejść. Stan wskazany dostaje znacznik na przejściach, które
 * z niego wychodzą — Operator widzi wtedy, co przebiegowi wolno dalej.
 *
 * Pusty stan bieżący znaczy „żaden przebieg nie jest wskazany"; wykaz stoi
 * wtedy bez wyróżnienia i jest samym modelem.
 */
export function wykazPrzejsc(stanBiezacy: string): HTMLElement {
  const lista = wykaz('Przejścia stanów przebiegu automatyki', 'da-wykaz');
  for (const przejscie of PRZEJSCIA_PRZEBIEGU) {
    const pozycja = pozycjaWykazu(
      `${nazwaStanu(przejscie.ze)} → ${nazwaStanu(przejscie.na)}`,
      przejscie.przyczyna,
      'da',
    );
    pozycja.element.dataset['stanZrodlowy'] = przejscie.ze;
    if (stanBiezacy !== '' && przejscie.ze === stanBiezacy) {
      pozycja.element.dataset['dostepne'] = 'tak';
    }
    lista.append(pozycja.element);
  }
  return lista;
}

/** Zdanie nad wykazem przejść: skąd wzięty jest model i czego dotyczy. */
export function zdanieOMaszynieStanow(stanBiezacy: string): string {
  const model =
    `Model stanów przebiegu: ${Object.keys(NAZWY_STANOW).length} stanów, ` +
    `${PRZEJSCIA_PRZEBIEGU.length} przejść.`;
  if (stanBiezacy === '') {
    return `${model} Żaden przebieg nie jest wskazany, więc wykaz stoi bez wyróżnienia.`;
  }
  return `${model} Przebieg wskazany jest w stanie „${nazwaStanu(stanBiezacy)}”; wyróżnione przejścia są dla niego dostępne.`;
}
