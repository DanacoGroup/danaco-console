import { liczbaOkienRozmowy } from '../okno-komunikacji/profil-modulu';
import { profilModulu } from '../okno-komunikacji/rejestr-profilow';
import {
  zdanieFiguryModulu,
  ZDANIE_MODUL_NIEUSTALONY,
  type FiguraDoNapisu,
  type GniazdoZRola,
} from './etykiety-ukladu';
import { ID_GNIAZD } from './identyfikatory';
import { rolaDomyslna } from './role-domyslne';

/**
 * Figura rozmowy modułu — ile okien scena prowadzi w tym module i w jakich
 * rolach.
 *
 * Jedna odpowiedzialność: przełożenie profilu modułu
 * (`okno-komunikacji/profil-modulu.ts`) na to, co układ okien równoległych umie
 * pokazać — liczbę gniazd i ich role. Własnego zdania o modułach ten plik nie
 * ma: liczba pochodzi z profilu, role z `role-domyslne.ts`, a sufit
 * z `identyfikatory.ts`.
 *
 * Liczbę okien bierzemy funkcją `liczbaOkienRozmowy`, nie polem `granicaOkien`:
 * funkcja oddaje zero dla modułu bez rozmowy, a pole niesie dla niego jedynkę.
 */
export interface FiguraModulu extends FiguraDoNapisu {
  /** Kod modułu, dla którego figura powstała; pusty = moduł jeszcze nieustalony. */
  kod: string;
  /** Zdanie dla Operatora — czym ta figura jest; nigdy puste. */
  zdanie: string;
}

/**
 * Moduły, dla których dwa okna rozmowy są ustaloną parą koordynator–wykonawca.
 *
 * Ról ten wykaz nie nadaje — `rolaDomyslna` i tak daje przy dwóch oknach
 * koordynatora i wykonawcę. Mówi wyłącznie tyle, że dla wymienionego modułu
 * para jest jego właściwością, a nie domyślną figurą silnika, którą Operator
 * może przestawić. Moduł, dla którego role okien nie są ustalone, tu nie stoi.
 *
 * Miejsce docelowe tego wykazu to profil modułu — pole z rolami okien; profil
 * dziś takiego pola nie ma.
 */
const MODULY_PARY: readonly string[] = ['automations'];

/**
 * Figura rozmowy modułu o podanym kodzie.
 *
 * Moduł spoza rejestru profilów nie wywraca układu — dostaje profil
 * wspólny, czyli stan zastany: zwykłe okno i pełny sufit platformy.
 */
export function figuraModulu(kod: string): FiguraModulu {
  const profil = profilModulu(kod);
  const liczbaOkien = liczbaOkienRozmowy(profil);
  const role = rolePrzyLiczbie(liczbaOkien);

  const opis: FiguraDoNapisu = {
    nazwa: profil.nazwa,
    postac: profil.postacRozmowy,
    liczbaOkien,
    role,
    paraKoordynatorWykonawca: MODULY_PARY.includes(profil.kod) && liczbaOkien === 2,
  };

  // Kod pusty znaczy „rdzeń jeszcze nie nazwał modułu okna" — i tak to brzmi.
  // Zdanie o module złożone z profilu wspólnego mówiłoby o module, którego
  // nikt nie wskazał, a pełny sufit podawałoby za jego właściwość.
  const zdanie = profil.kod.length > 0 ? zdanieFiguryModulu(opis) : ZDANIE_MODUL_NIEUSTALONY;

  return { ...opis, kod: profil.kod, zdanie };
}

/**
 * Role gniazd przy zadanej liczbie okien — dokładnie te, które nada
 * `ustawLiczbe`.
 *
 * Dla dwóch okien wychodzi z tego para koordynator–wykonawca, bo tak stanowi
 * `role-domyslne.ts`. Figura modułu nie nadaje ról sama i nie zna ani jednej
 * ich nazwy.
 */
export function rolePrzyLiczbie(liczbaOkien: number): readonly GniazdoZRola[] {
  return ID_GNIAZD.slice(0, Math.max(0, liczbaOkien)).map((id) => ({
    id,
    rola: rolaDomyslna(id, liczbaOkien),
  }));
}
