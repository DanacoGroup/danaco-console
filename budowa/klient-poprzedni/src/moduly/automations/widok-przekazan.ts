/**
 * Wykaz przekazań oczekujących — czysta konstrukcja z bytów
 * `PrzekazanyScenariusz`. Plik nie zna ani źródła, ani stanu modułu: dostaje
 * wykaz i wywołanie zwrotne zapisu, a pozycja mówi, skąd przyszedł scenariusz
 * i ile ma kroków.
 */
import { przyciskAkcji as przycisk, pozycjaWykazu, wykaz } from '../../modele/kontrolki-formularza';
import type { PrzekazanyScenariusz } from './przyjecie-przekazania';



/**
 * Zdanie przy pozycji wykazu: liczba kroków w odmianie właściwej dla
 * liczebnika, stan automatyki po zapisie, okno pochodzenia oraz polecenie
 * wyjściowe scenariusza.
 */
export function opisPrzekazania(scenariusz: PrzekazanyScenariusz): string {
  const czesci = [
    `${scenariusz.kroki.length} ${odmianaKrokow(scenariusz.kroki.length)}`,
    scenariusz.czynny ? 'po zapisie czynna' : 'po zapisie wyłączona',
    `okno ${scenariusz.idOkna}`,
  ];
  const dopowiedzenie = scenariusz.polecenie === '' ? scenariusz.opis : scenariusz.polecenie;
  return dopowiedzenie === '' ? czesci.join(' · ') : `${czesci.join(' · ')} — ${dopowiedzenie}`;
}

/**
 * Odmiana rzeczownika „krok" przez liczbę według reguł polskiej liczebności:
 * forma pojedyncza dla jedności, mnoga dla końcówek od dwóch do czterech poza
 * nastką, dopełniaczowa dla reszty.
 */
function odmianaKrokow(liczba: number): string {
  const reszta = liczba % 10;
  const setka = liczba % 100;
  if (liczba === 1) return 'krok';
  if (reszta >= 2 && reszta <= 4 && (setka < 12 || setka > 14)) return 'kroki';
  return 'kroków';
}

/**
 * Buduje wykaz przekazań wraz z przyciskiem zapisu przy każdym. Wykaz pusty
 * oddaje `null` — sekcja bez ani jednej pozycji nie ma po co stać na ekranie,
 * a zdanie o pustce niesie stan treści okna.
 */
export function listaPrzekazan(
  scenariusze: readonly PrzekazanyScenariusz[],
  zapisz: (scenariusz: PrzekazanyScenariusz) => void,
): HTMLElement | null {
  if (scenariusze.length === 0) return null;
  const lista = wykaz('Przekazania oczekujące na decyzję', 'da-wykaz');
  for (const scenariusz of scenariusze) {
    const pozycja = pozycjaWykazu(scenariusz.nazwa, opisPrzekazania(scenariusz), 'da');
    pozycja.element.dataset['przekazanie'] = scenariusz.idOkna;
    const przyjmij = przycisk('Zapisz jako automatykę', 'dn-btn dn-btn--atrament');
    przyjmij.addEventListener('click', () => zapisz(scenariusz));
    pozycja.akcje.append(przyjmij);
    lista.append(pozycja.element);
  }
  return lista;
}
