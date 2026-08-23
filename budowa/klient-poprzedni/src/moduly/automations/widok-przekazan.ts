import { przyciskAkcji as przycisk, pozycjaWykazu, wykaz } from '../../modele/kontrolki-formularza';
import type { PrzekazanyScenariusz } from './przyjecie-przekazania';

/**
 * Wykaz przekazań oczekujących — czysta konstrukcja z bytów
 * `PrzekazanyScenariusz`. Plik nie zna ani źródła, ani stanu modułu: dostaje
 * wykaz i wywołanie zwrotne zapisu.
 *
 * Pozycja mówi, skąd przyszedł scenariusz i ile ma kroków, bo to jedyne, czym
 * Operator może się kierować przed zapisem — treści kroków nie pokazujemy tu
 * wcale, ta należy do Workflow Buildera po zapisie.
 */

/** Zdanie przy pozycji: liczba kroków, stan po zapisie i polecenie wyjściowe. */
export function opisPrzekazania(scenariusz: PrzekazanyScenariusz): string {
  const czesci = [
    `${scenariusz.kroki.length} ${odmianaKrokow(scenariusz.kroki.length)}`,
    scenariusz.czynny ? 'po zapisie czynna' : 'po zapisie wyłączona',
    `okno ${scenariusz.idOkna}`,
  ];
  const dopowiedzenie = scenariusz.polecenie === '' ? scenariusz.opis : scenariusz.polecenie;
  return dopowiedzenie === '' ? czesci.join(' · ') : `${czesci.join(' · ')} — ${dopowiedzenie}`;
}

/** Odmiana rzeczownika „krok" przez liczbę, według reguł polskiej liczebności. */
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
