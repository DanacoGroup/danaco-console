/**
 * Maszyna stanów przebiegu pokazywana w zakładce stanów okna Orchestratora.
 * Stany są kompletem wyliczenia `AutomationExecutionStatus`, a plik nie woła
 * rdzenia: stan bieżący przychodzi z zewnątrz.
 */
import { AutomationExecutionStatus } from '../../../../shared/contract';
import { pozycjaWykazu, wykaz } from '../../modele/kontrolki-formularza';

/**
 * Jedno przejście maszyny stanów wraz z jego przyczyną: stan wyjściowy, stan
 * docelowy oraz to, co przejście powoduje — działanie Operatora albo zamknięcie
 * przebiegu przez rdzeń.
 */
export interface PrzejscieStanu {
  ze: string;
  na: string;
  /** Co powoduje przejście — działanie Operatora albo zamknięcie przebiegu. */
  przyczyna: string;
}

/**
 * Komplet przejść przebiegu automatyki, wywiedziony z działań, które kontrakt na
 * przebiegu dopuszcza: uruchomienia, wstrzymania, wznowienia, zatrzymania oraz
 * zamknięcia przebiegu powodzeniem albo błędem.
 */
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

/**
 * Nazwy stanów przebiegu w mowie Operatora, przy czym kluczami są wartości
 * kontraktu. Dzięki temu zmiana nazwy pokazywanej na ekranie nie rusza wartości
 * idącej do rdzenia ani z niego wracającej.
 */
const NAZWY_STANOW: Readonly<Record<string, string>> = {
  [AutomationExecutionStatus.Pending]: 'oczekuje',
  [AutomationExecutionStatus.Running]: 'w toku',
  [AutomationExecutionStatus.Paused]: 'wstrzymany',
  [AutomationExecutionStatus.Succeeded]: 'zakończony powodzeniem',
  [AutomationExecutionStatus.Failed]: 'zakończony błędem',
  [AutomationExecutionStatus.Stopped]: 'przerwany',
};

/**
 * Nazwa stanu w mowie Operatora; stan nieznany zostaje przy swojej wartości
 * kontraktowej, ponieważ nazwa zmyślona po stronie okna byłaby słowem, za którym
 * nie stoi żaden zapis rdzenia.
 */
export function nazwaStanu(stan: string): string {
  return NAZWY_STANOW[stan] ?? stan;
}

/**
 * Buduje wykaz przejść. Stan wskazany dostaje znacznik na przejściach, które
 * z niego wychodzą, dzięki czemu Operator widzi, co przebiegowi wolno dalej.
 * Pusty stan bieżący zostawia wykaz bez wyróżnienia.
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

/**
 * Zdanie nad wykazem przejść: mówi, skąd wzięty jest model stanów i czego
 * dotyczy, żeby Operator czytał wykaz jako odwzorowanie kontraktu, a nie jako
 * propozycję okna.
 */
export function zdanieOMaszynieStanow(stanBiezacy: string): string {
  const model =
    `Model stanów przebiegu: ${Object.keys(NAZWY_STANOW).length} stanów, ` +
    `${PRZEJSCIA_PRZEBIEGU.length} przejść.`;
  if (stanBiezacy === '') {
    return `${model} Żaden przebieg nie jest wskazany, więc wykaz stoi bez wyróżnienia.`;
  }
  return `${model} Przebieg wskazany jest w stanie „${nazwaStanu(stanBiezacy)}”; wyróżnione przejścia są dla niego dostępne.`;
}
