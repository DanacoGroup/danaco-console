import { ConfigScope } from '../../../../shared/contract';

/**
 * Poziomy zasięgu modułu Workspace, jedno źródło nazw i kolejności dla wszystkich
 * okien modułu, tak aby ta sama nastawa nazywała się wszędzie tak samo.
 */
interface OpisPoziomu {
  /** Nazwa poziomu widoczna dla Operatora — ta sama we wszystkich oknach. */
  nazwa: string;
  /** Miejsce w kolejności od najszerszego; puste znaczy poza wyborem. */
  rzad: number | null;
  /** Powód pominięcia poziomu w kontrolkach modułu. */
  powodPominiecia?: string;
}

const WYKAZ: Record<ConfigScope, OpisPoziomu> = {
  [ConfigScope.Application]: {
    nazwa: 'aplikacja',
    rzad: null,
    powodPominiecia:
      'nastawa samego programu, nie treści w nim prowadzonej — instrukcje projektu ' +
      'i ustalenia pamięci nie mają na tym poziomie czego zapisać',
  },
  [ConfigScope.Global]: { nazwa: 'globalny', rzad: 1 },
  [ConfigScope.Environment]: { nazwa: 'środowisko', rzad: 2 },
  [ConfigScope.Module]: { nazwa: 'moduł', rzad: 3 },
  [ConfigScope.ModulePair]: { nazwa: 'para modułów', rzad: 4 },
  [ConfigScope.Project]: { nazwa: 'projekt', rzad: 5 },
  [ConfigScope.Session]: { nazwa: 'karta sesji', rzad: 6 },
  [ConfigScope.Role]: { nazwa: 'rola', rzad: 7 },
  [ConfigScope.Window]: { nazwa: 'okno komunikacji', rzad: 8 },
};

/** Poziom zasięgu wraz z jego nazwą widoczną dla Operatora — para gotowa do wstawienia w liście wyboru. */
export type ParaPoziomu = readonly [ConfigScope, string];

/**
 * Osiem poziomów wyboru w kolejności od najszerszego do najwęższego.
 * Poziom najwęższy wygrywa — tę regułę rozstrzyga rdzeń, nie okno.
 */
export const POZIOMY_ZASIEGU: readonly ParaPoziomu[] = Object.entries(WYKAZ)
  .filter((wpis): wpis is [ConfigScope, OpisPoziomu & { rzad: number }] => wpis[1].rzad !== null)
  .sort((a, b) => a[1].rzad - b[1].rzad)
  .map(([poziom, opis]) => [poziom, opis.nazwa] as ParaPoziomu);

/**
 * Nazwa poziomu — także tego spoza wyboru, bo rdzeń może go oddać w odpowiedzi.
 * Wartość nieznana kontraktowi zostaje surowym kodem.
 */
export function nazwaPoziomu(zasieg: ConfigScope): string {
  return WYKAZ[zasieg]?.nazwa ?? zasieg;
}

/** Podaje powód, dla którego dany poziom nie stoi w wyborze; wartość pusta dla poziomów obecnych w wyborze. */
export function powodPominiecia(zasieg: ConfigScope): string {
  return WYKAZ[zasieg]?.powodPominiecia ?? '';
}

/**
 * Ten sam wykaz poziomów zasięgu z jednym wskazanym poziomem przesuniętym na
 * czoło kolejności prezentacji, bez zmiany nazw poziomów.
 */
export function poziomyZPierwszym(pierwszy: ConfigScope): readonly ParaPoziomu[] {
  const wybrany = POZIOMY_ZASIEGU.filter(([poziom]) => poziom === pierwszy);
  return [...wybrany, ...POZIOMY_ZASIEGU.filter(([poziom]) => poziom !== pierwszy)];
}
