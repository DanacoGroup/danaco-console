import { ConfigScope } from '../../../../shared/contract';

/**
 * Poziomy zasięgu modułu Workspace — jedno źródło nazw i kolejności dla
 * wszystkich okien modułu, żeby ta sama nastawa nazywała się wszędzie tak samo.
 *
 * Poziom `application` stoi w wykazie, ale poza wyborem: kontrakt przypisuje mu
 * nastawy samego programu (wymóg logowania, adres nasłuchu), a nie treści w nim
 * prowadzonej, więc kontrolka oferująca ten poziom proponowałaby zapis, którego
 * rdzeń nie ma gdzie umieścić. Nazwę zachowuje, bo rdzeń może oddać ten poziom
 * w odpowiedzi (`workspace.instructions.set` niesie warstwę obowiązującą)
 * i okno ma go wtedy nazwać, a nie pokazać surowego kodu.
 *
 * Wykaz jest zapisany jako `Record<ConfigScope, …>`, więc poziom dołożony do
 * kontraktu nie przejdzie kompilacji, dopóki nie zostanie rozstrzygnięte, czy
 * moduł ma go pokazywać; tablica przyjęłaby nowy poziom bez śladu.
 */
interface OpisPoziomu {
  /** Nazwa poziomu widoczna dla Operatora — ta sama we wszystkich oknach. */
  nazwa: string;
  /**
   * Miejsce w kolejności od najszerszego; `null` znaczy „poza
   * wyborem” i wtedy `powodPominiecia` mówi, dlaczego.
   */
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

/** Poziom zasięgu wraz z jego nazwą — para gotowa dla listy wyboru. */
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

/** Powód, dla którego poziom nie stoi w wyborze; pusty dla poziomów w wyborze. */
export function powodPominiecia(zasieg: ConfigScope): string {
  return WYKAZ[zasieg]?.powodPominiecia ?? '';
}

/**
 * Ten sam wykaz z jednym poziomem przesuniętym na czoło.
 *
 * Pamięć projektu otwiera się na poziomie `project`, bo taki jest jej domyślny
 * zasięg. Zmienia się wyłącznie kolejność prezentacji; nazwy zostają te same,
 * więc oba okna pokazują tę samą nastawę pod tą samą nazwą.
 */
export function poziomyZPierwszym(pierwszy: ConfigScope): readonly ParaPoziomu[] {
  const wybrany = POZIOMY_ZASIEGU.filter(([poziom]) => poziom === pierwszy);
  return [...wybrany, ...POZIOMY_ZASIEGU.filter(([poziom]) => poziom !== pierwszy)];
}
