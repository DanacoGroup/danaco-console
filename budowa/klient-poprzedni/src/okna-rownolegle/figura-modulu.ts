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
 * Figura rozmowy modułu przekłada profil modułu na to, co układ okien równoległych umie pokazać — liczbę gniazd i ich role, bez własnego zdania o samych modułach.
 */
export interface FiguraModulu extends FiguraDoNapisu {
  /** Kod modułu, dla którego figura powstała; pusty = moduł jeszcze nieustalony. */
  kod: string;
  /** Zdanie dla Operatora — czym ta figura jest; nigdy puste. */
  zdanie: string;
}

/**
 * Moduły, dla których dwa okna rozmowy są ustaloną parą koordynator-wykonawca, nie domyślną figurą silnika, którą Operator mógłby przestawić.
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

  // Kod pusty znaczy, że rdzeń jeszcze nie nazwał modułu okna, więc zdanie mówi to wprost.
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
