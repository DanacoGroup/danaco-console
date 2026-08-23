import {
  ConfigScope,
  IsolationContextKind,
  IsolationTechnicalScope,
} from '../../../shared/contract';

/**
 * Słownik nazw punktów izolacji i poziomów zasięgu.
 *
 * Obszar Profili składa i pokazuje ten sam zestaw kluczy, co obszary
 * „Kontekst" i „Zakres techniczny" — profil jest nazwanym zestawem trzech
 * przełączników kontekstu i ośmiu zakresów technicznych. Nazwy stoją w jednym
 * miejscu, żeby profil nie nazywał klucza inaczej niż zakładka, w której
 * Operator ten sam klucz przestawia.
 *
 * Plik trzyma wyłącznie nazewnictwo: stany przełączników przychodzą z rdzenia
 * (`isolation.context.get`, `isolation.technical.get`, `isolation.profile.load`).
 */

/** Jeden przełącznik kontekstu w profilu: klucz kontraktu i nazwa widoczna. */
export interface PozycjaKontekstu {
  kind: IsolationContextKind;
  nazwa: string;
}

/** Jeden zakres techniczny w profilu: klucz kontraktu i nazwa widoczna. */
export interface PozycjaTechniczna {
  scope: IsolationTechnicalScope;
  nazwa: string;
}

/** Trzy przełączniki kontekstu — wartość „odrębna" znaczy `isolated: true`. */
export const POZYCJE_KONTEKSTU: readonly PozycjaKontekstu[] = [
  { kind: IsolationContextKind.History, nazwa: 'Historia wymiany' },
  { kind: IsolationContextKind.Memory, nazwa: 'Pamięć długoterminowa' },
  { kind: IsolationContextKind.Context, nazwa: 'Bieżący stan roboczy' },
];

/** Osiem zakresów technicznych — wartość „włączony" znaczy `isolated: true`. */
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

/** Nazwy poziomów zasięgu kontraktu po polsku. */
export const ETYKIETY_ZASIEGU: Readonly<Record<ConfigScope, string>> = {
  // Zakres „aplikacja" obejmuje całą aplikację, nie pojedyncze okno, sesję ani
  // środowisko. Słownik pokrywa cały kontrakt: brak wartości zatrzymuje kompilację.
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

/** Poziomy w kolejności od najszerszego — porządek listy wyboru celu przypisania. */
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

/** Pary wartość–etykieta dla listy wyboru poziomu zasięgu. */
export const OPCJE_ZASIEGU: ReadonlyArray<readonly [string, string]> = ZASIEGI_OD_NAJSZERSZEGO.map(
  (poziom) => [poziom, ETYKIETY_ZASIEGU[poziom]] as const,
);

/** Etykieta wartości przełącznika kontekstu. */
export function etykietaKontekstu(isolated: boolean): string {
  return isolated ? 'Odrębna' : 'Współdzielona';
}

/** Etykieta wartości zakresu technicznego. */
export function etykietaTechniczna(isolated: boolean): string {
  return isolated ? 'Włączony' : 'Wyłączony';
}
