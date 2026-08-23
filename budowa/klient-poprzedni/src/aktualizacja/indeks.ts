/**
 * Wejście warstwy aktualizacji — jedyna rzecz, którą punkt wejścia aplikacji
 * musi o niej wiedzieć.
 *
 * Reszta (adres kanału wydań, porównanie wersji, most do powłoki) jest
 * własnością tego katalogu i nie wychodzi na zewnątrz.
 *
 * Widok wykazu wydań wychodzi stąd osobno od banera: baner przywołuje go
 * przyciskiem, ale widok jest gotowy do osadzenia także w stałym miejscu.
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
