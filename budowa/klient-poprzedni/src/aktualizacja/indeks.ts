/**
 * Wejście warstwy aktualizacji — jedyna rzecz, którą punkt wejścia aplikacji
 * musi o niej wiedzieć.
 *
 * Wychodzą stąd baner aktualizacji, widok wykazu wydań oraz odczyt i porównanie
 * wydań; reszta pozostaje własnością katalogu.
 */
export { utworzBanerAktualizacji, type BanerAktualizacji } from './baner-aktualizacji';
export {
  utworzWykazWydanWidok,
  type WykazWydanWidok,
  type UstawieniaWykazu,
} from './wykaz-wydan-widok';
export {
  nowszaWersja,
  nowszeWydanie,
  pobierzWykazWydan,
  ADRES_WYDAN,
  type Wydanie,
  type OdczytWykazu,
} from './wykaz-wydan';
