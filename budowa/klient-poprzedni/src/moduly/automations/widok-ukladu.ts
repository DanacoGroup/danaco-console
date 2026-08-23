import { type AutomationDependency, type AutomationOrchestratorDefineResponse } from '../../../../shared/contract';
import { przyciskAkcji as przycisk, pozycjaWykazu, wykaz } from '../../modele/kontrolki-formularza';

/**
 * Rysowanie układu zależności — czyste fragmenty widoku Orchestratora.
 *
 * Jedyna odpowiedzialność tego pliku to postać układu na ekranie. Żaden fragment
 * nie woła rdzenia i nie zna stanu okna; jeden bierze wywołanie zwrotne
 * usunięcia, bo przycisk wiersza musi sięgnąć po zapis układu.
 */

/** Zależność nazwana jednym napisem — do zdań potwierdzenia i porównań. */
export function opisZaleznosci(zaleznosc: AutomationDependency): string {
  return `${zaleznosc.fromStepId} → ${zaleznosc.toStepId}`;
}

/** Czy rdzeń oddał tę zależność — porównanie po obu krokach i rodzaju. */
export function zawiera(
  oddane: readonly AutomationDependency[],
  szukana: AutomationDependency,
): boolean {
  return oddane.some(
    (wpis) =>
      wpis.fromStepId === szukana.fromStepId &&
      wpis.toStepId === szukana.toStepId &&
      wpis.kind === szukana.kind,
  );
}

/**
 * Ocena poprawności układu — znacznik `data-ocena` dla sprawdzianów widoku.
 *
 * Układ wadliwy bez zastrzeżeń nie zapowiada wykazu, którego nie będzie:
 * zapowiedź „poniżej ich wykaz" przy pustym polu `issues` obiecywałaby treść,
 * której rdzeń nie przysłał.
 */
export function ocenaUkladu(wynik: AutomationOrchestratorDefineResponse): HTMLElement {
  const ocena = document.createElement('p');
  ocena.className = 'dn-pole-opis';
  ocena.dataset['ocena'] = wynik.valid ? 'poprawny' : 'wadliwy';
  const zastrzezen = (wynik.issues ?? []).length;
  if (wynik.valid) {
    ocena.textContent =
      'Układ zależności jest poprawny — kroki dają się ułożyć w kolejność wykonania.';
  } else if (zastrzezen === 0) {
    ocena.textContent =
      'Rdzeń orzekł układ za wadliwy, lecz nie podał ani jednego zastrzeżenia — powodu nie ma skąd wziąć.';
  } else {
    ocena.textContent = `Układ zależności ma zastrzeżenia (${zastrzezen}) — poniżej ich wykaz.`;
  }
  return ocena;
}

/** Wykaz zastrzeżeń; `null` znaczy „układ bez zastrzeżeń", nie „brak danych". */
export function zastrzezeniaUkladu(zdania: readonly string[]): HTMLElement | null {
  if (zdania.length === 0) return null;
  const zastrzezenia = wykaz('Zastrzeżenia do układu', 'da-wykaz');
  for (const zdanie of zdania) {
    zastrzezenia.append(pozycjaWykazu('zastrzeżenie', zdanie, 'da').element);
  }
  return zastrzezenia;
}

/** Stan pusty układu — automatyka bez zależności jest poprawna, nie wadliwa. */
export function pustyUklad(): HTMLElement {
  const puste = document.createElement('p');
  puste.className = 'dn-pole-opis';
  puste.textContent = 'Automatyka nie ma jeszcze zależności — kroki wykonają się w kolejności zapisu.';
  return puste;
}

/** Wykaz zależności; krok na ścieżce krytycznej dostaje `data-sciezka`. */
export function listaZaleznosci(
  uklad: readonly AutomationDependency[],
  sciezka: ReadonlySet<string>,
  naUsuniecie: (zaleznosc: AutomationDependency) => void,
): HTMLElement {
  const lista = wykaz('Zależności układu', 'da-wykaz');
  for (const zaleznosc of uklad) {
    const pozycja = pozycjaWykazu(
      `${zaleznosc.fromStepId} → ${zaleznosc.toStepId}`,
      `${zaleznosc.kind}${zaleznosc.condition === undefined ? '' : ` — ${zaleznosc.condition}`}`,
      'da',
    );
    if (sciezka.has(zaleznosc.fromStepId) && sciezka.has(zaleznosc.toStepId)) {
      pozycja.element.dataset['sciezka'] = 'krytyczna';
    }
    const usun = przycisk('Usuń', 'dn-btn dn-btn--zarys');
    usun.addEventListener('click', () => naUsuniecie(zaleznosc));
    pozycja.akcje.append(usun);
    lista.append(pozycja.element);
  }
  return lista;
}

/**
 * Opis ścieżki krytycznej.
 *
 * Puste `criticalPathStepIds` nie znaczy cyklu — cykl jest jedną z możliwych
 * przyczyn pustki, nie jedyną. Okno nie orzeka przyczyny, której rdzeń nie
 * podał; powód, jeśli istnieje, stoi w wykazie zastrzeżeń.
 */
export function opisSciezkiKrytycznej(sciezka: ReadonlySet<string>): HTMLElement {
  const krytyczna = document.createElement('p');
  krytyczna.className = 'dn-pole-opis';
  krytyczna.dataset['sciezkaKrytyczna'] = 'tak';
  krytyczna.textContent =
    sciezka.size === 0
      ? 'Rdzeń nie wskazał ścieżki krytycznej dla tego układu.'
      : `Ścieżka krytyczna: ${[...sciezka].join(' → ')}.`;
  return krytyczna;
}
