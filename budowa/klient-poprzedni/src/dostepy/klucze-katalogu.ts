import {
  ConfigAxis,
  ConfigScope,
  SettingValueType,
  type SettingDefinition,
} from '../../../shared/contract';

/**
 * Klucze katalogu roboczego modelu.
 *
 * Dostęp to nie katalog roboczy. Model może mieć wgląd w `C:\Projekty\Lex`
 * przez nadanie dostępu, a swoje pliki zostawiać w
 * `<instalacja>/sesje/<identyfikator>`. To dwa niezależne ustawienia i dwa
 * osobne obszary tej sekcji.
 *
 * Katalog roboczy jest zwykłym ustawieniem: idzie komendami `config.get`,
 * `config.set` i `config.reset`, przez ten sam rezolwer zasięgów co reszta
 * konfiguracji. Osobnej komendy dla niego nie ma.
 */

/** Katalog, w którym powstają katalogi sesyjne modelu. */
export const KLUCZ_PODSTAWA = 'katalog.roboczy.podstawa';

/** Wzorzec nazwy katalogu jednej sesji wewnątrz podstawy. */
export const KLUCZ_WZORZEC = 'katalog.roboczy.wzorzec_sesji';

/** Wzorzec domyślny; odpowiednik `core.WzorzecSesjiDomyslny`. */
export const WZORZEC_DOMYSLNY = 'sesje/<identyfikator>';

/**
 * Zdanie o wartości domyślnej podstawy.
 *
 * Katalog roboczy jest domyślnie tworzony w miejscu instalacji aplikacji
 * głównej. Klient tej ścieżki nie zna i jej nie zgaduje: mówi, skąd się bierze,
 * zamiast pokazać zmyśloną ścieżkę jako fakt.
 */
export const OPIS_PODSTAWY_DOMYSLNEJ =
  'Katalog tworzony automatycznie w miejscu instalacji aplikacji głównej. Rdzeń ustala go sam przy pierwszym uruchomieniu; wpisz ścieżkę tylko wtedy, gdy chcesz to nadpisać.';

/** Klucze obsługiwane przez obszar katalogu roboczego, w kolejności pól. */
export const KLUCZE_KATALOGU: readonly string[] = [KLUCZ_PODSTAWA, KLUCZ_WZORZEC];

/**
 * Definicje zastępcze na wypadek katalogu, który nie zna tych kluczy.
 *
 * Brak wiersza w katalogu ustawień nie może zabrać możliwości ustawienia
 * katalogu roboczego. Definicja z katalogu ma pierwszeństwo; ta poniżej wchodzi
 * wyłącznie wtedy, gdy katalog milczy, i jest zbudowana z tych samych wartości,
 * co odpowiadające jej definicje rdzenia (`core.DefinicjeKataloguRoboczego`).
 */
export function definicjaZastepcza(klucz: string): SettingDefinition | null {
  if (klucz === KLUCZ_PODSTAWA) {
    return {
      ...szkielet(klucz, 1),
      name: 'Katalog roboczy — podstawa',
      description: OPIS_PODSTAWY_DOMYSLNEJ,
      valueType: SettingValueType.Path,
      placeholder: 'ścieżka katalogu, np. D:\\Danaco\\robocze',
    };
  }
  if (klucz === KLUCZ_WZORZEC) {
    return {
      ...szkielet(klucz, 2),
      name: 'Wzorzec katalogu sesji',
      description:
        'Nazwa katalogu jednej sesji wewnątrz podstawy. Znacznik <identyfikator> zastępowany jest identyfikatorem sesji.',
      valueType: SettingValueType.String,
      defaultValue: WZORZEC_DOMYSLNY,
      placeholder: WZORZEC_DOMYSLNY,
    };
  }
  return null;
}

/** Człony wspólne definicji zastępczych. */
function szkielet(klucz: string, kolejnosc: number): SettingDefinition {
  return {
    key: klucz,
    categoryId: 'katalog-roboczy',
    name: klucz,
    valueType: SettingValueType.String,
    allowedScopes: [...ZASIEGI_KATALOGU],
    allowedAxes: [ConfigAxis.Platform, ConfigAxis.Model, ConfigAxis.Account],
    required: false,
    order: kolejnosc,
    enabled: true,
  };
}

/**
 * Poziomy, na których katalog roboczy ma sens — nie wszystkie z `ConfigScope`:
 * katalog roboczy jest własnością instalacji i sesji, nie pary modułów.
 */
export const ZASIEGI_KATALOGU: readonly ConfigScope[] = [
  ConfigScope.Global,
  ConfigScope.Environment,
  ConfigScope.Project,
  ConfigScope.Session,
  ConfigScope.Window,
];
