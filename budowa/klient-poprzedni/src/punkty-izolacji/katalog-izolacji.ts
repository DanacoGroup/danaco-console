import {
  ConfigScope,
  IsolationContextKind,
  IsolationTechnicalScope,
} from '../../../shared/contract';

// Słownik nazw punktów izolacji i poziomów zasięgu, wspólny dla obszarów okna i dla profili.

/** Jeden przełącznik kontekstu w profilu: klucz kontraktu izolacji wraz z jego nazwą widoczną Operatorowi. */
export interface PozycjaKontekstu {
  kind: IsolationContextKind;
  nazwa: string;
}

/** Jeden zakres techniczny w profilu: klucz kontraktu izolacji wraz z jego nazwą widoczną dla Operatora. */
export interface PozycjaTechniczna {
  scope: IsolationTechnicalScope;
  nazwa: string;
}

/** Trzy przełączniki kontekstu izolacji — wartość odrębna znaczy, że dany kontekst jest tutaj izolowany. */
export const POZYCJE_KONTEKSTU: readonly PozycjaKontekstu[] = [
  { kind: IsolationContextKind.History, nazwa: 'Historia wymiany' },
  { kind: IsolationContextKind.Memory, nazwa: 'Pamięć długoterminowa' },
  { kind: IsolationContextKind.Context, nazwa: 'Bieżący stan roboczy' },
];

/** Osiem zakresów technicznych izolacji — wartość włączony znaczy, że dany zakres jest tutaj izolowany. */
export const POZYCJE_TECHNICZNE: readonly PozycjaTechniczna[] = [
  { scope: IsolationTechnicalScope.WorkingDirectory, nazwa: 'Katalog roboczy sesji' },
  { scope: IsolationTechnicalScope.ProcessEnvironment, nazwa: 'Środowisko procesu' },
  { scope: IsolationTechnicalScope.ModelDataDirectory, nazwa: 'Katalog danych modelu' },
  { scope: IsolationTechnicalScope.NetworkAccess, nazwa: 'Dostęp sieciowy' },
  { scope: IsolationTechnicalScope.FileAccess, nazwa: 'Odczyt/zapis plików' },
  { scope: IsolationTechnicalScope.AccountToken, nazwa: 'Konto i token' },
  { scope: IsolationTechnicalScope.ProcessModel, nazwa: 'Model procesu' },
  { scope: IsolationTechnicalScope.ExecutionServer, nazwa: 'Serwer wykonania' },
];

/** Nazwy poziomów zasięgu kontraktu po polsku, wypisane osobno dla każdego poziomu, bez wyjątku żadnego. */
export const ETYKIETY_ZASIEGU: Readonly<Record<ConfigScope, string>> = {
  // Zakres aplikacja obejmuje całą aplikację, nie pojedyncze okno, sesję ani środowisko klienta.
  [ConfigScope.Application]: 'aplikacja',
  [ConfigScope.Global]: 'Globalny',
  [ConfigScope.Environment]: 'Środowisko',
  [ConfigScope.Module]: 'Moduł',
  [ConfigScope.ModulePair]: 'Para modułów',
  [ConfigScope.Project]: 'Projekt',
  [ConfigScope.Session]: 'Karta sesji',
  [ConfigScope.Role]: 'Rola',
  [ConfigScope.Window]: 'Okno komunikacji',
};

/** Poziomy zasięgu w kolejności od najszerszego — porządek listy wyboru celu przypisania tego profilu tutaj. */
export const ZASIEGI_OD_NAJSZERSZEGO: readonly ConfigScope[] = [
  ConfigScope.Global,
  ConfigScope.Environment,
  ConfigScope.Module,
  ConfigScope.ModulePair,
  ConfigScope.Project,
  ConfigScope.Session,
  ConfigScope.Role,
  ConfigScope.Window,
];

/** Pary wartość-etykieta dla listy wyboru poziomu zasięgu, złożone z tabeli poziomów liczonych od najszerszego. */
export const OPCJE_ZASIEGU: ReadonlyArray<readonly [string, string]> = ZASIEGI_OD_NAJSZERSZEGO.map(
  (poziom) => [poziom, ETYKIETY_ZASIEGU[poziom]] as const,
);

/** Etykieta wartości przełącznika kontekstu, złożona z dwóch możliwych stanów izolacji tego samego kontekstu. */
export function etykietaKontekstu(isolated: boolean): string {
  return isolated ? 'Odrębna' : 'Współdzielona';
}

/** Etykieta wartości zakresu technicznego, złożona z dwóch możliwych stanów izolacji tego zakresu pracy. */
export function etykietaTechniczna(isolated: boolean): string {
  return isolated ? 'Włączony' : 'Wyłączony';
}
