/**
 * Klucze katalogu roboczego modelu wraz z definicjami zastępczymi. Plik nazywa
 * klucz podstawy i klucz wzorca sesji, podaje wzorzec domyślny oraz opis
 * podstawy, składa definicję zastępczą klucza i wylicza zasięgi katalogu.
 */
import {
  ConfigAxis,
  ConfigScope,
  SettingValueType,
  type SettingDefinition,
} from '../../../shared/contract';

/**
 * Klucz ustawienia wskazującego katalog, w którym powstają katalogi sesyjne
 * modelu. Wartość jest ścieżką w systemie plików maszyny, na której stoi rdzeń.
 */
export const KLUCZ_PODSTAWA = 'katalog.roboczy.podstawa';

/**
 * Klucz ustawienia niosącego wzorzec nazwy katalogu jednej sesji wewnątrz
 * podstawy. Wzorzec jest nazwą względną, więc katalog sesji nie wychodzi poza
 * podstawę.
 */
export const KLUCZ_WZORZEC = 'katalog.roboczy.wzorzec_sesji';

/**
 * Wzorzec domyślny nazwy katalogu sesji, równy wartości domyślnej rdzenia.
 * Klient trzyma go u siebie, żeby pole miało podpowiedź także wtedy, gdy katalog
 * ustawień tego klucza nie zna.
 */
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

/**
 * Klucze obsługiwane przez obszar katalogu roboczego, w kolejności pól
 * formularza. Obszar czyta i zapisuje wyłącznie te dwa klucze, więc wykaz jest
 * zarazem jego zakresem.
 */
export const KLUCZE_KATALOGU: readonly string[] = [KLUCZ_PODSTAWA, KLUCZ_WZORZEC];

/**
 * Definicja zastępcza na wypadek katalogu ustawień, który tych kluczy nie zna.
 * Definicja z katalogu ma pierwszeństwo, a ta wchodzi wyłącznie wtedy, gdy
 * katalog milczy, i niesie te same wartości co odpowiadające jej definicje
 * rdzenia.
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

/**
 * Człony wspólne definicji zastępczych: klucz, kategoria, dopuszczone zasięgi
 * i osie oraz kolejność pola. Wywołanie uzupełnia je nazwą, opisem i rodzajem
 * wartości.
 */
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
