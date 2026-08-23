import type { AutomationWorkflow } from '../../../../shared/contract';
import { przyciskAkcji as przycisk, pozycjaWykazu, wykaz } from '../../modele/kontrolki-formularza';

/**
 * Widok wykazu gotowych pętli — czysta konstrukcja z bytów `AutomationWorkflow`.
 * Plik nie zna ani źródła, ani stanu modułu, ani stanu treści okna: dostaje
 * wykaz, napis szukania i wywołanie zwrotne uruchomienia, a oddaje albo gotowy
 * element, albo zdanie o pustce. Rozstrzygnięcie, co z tym zrobić, należy do
 * okna, bo w tym samym miejscu treści stoi jeszcze wykaz przekazań.
 *
 * Przycisk uruchomienia ma każda pozycja, także pętla wyłączona — jej stan stoi
 * w opisie pozycji, a nie w blokadzie wiersza. Wykaz skraca zawężenie: napis
 * szukania albo przełącznik „tylko czynne”, który zawęża żądanie do rdzenia
 * (`enabledOnly`).
 *
 * Szukanie idzie po stronie klienta, bo `automation.workflow.list` przyjmuje
 * wyłącznie `enabledOnly` i `limit`, a wykaz jest już w oknie. Dopasowanie bierze
 * kolejno nazwę, identyfikator i opis.
 */

/** Zawęża wykaz napisem szukania; pusty napis oddaje wykaz w całości. */
export function dopasowanePetle(
  petle: readonly AutomationWorkflow[],
  szukane: string,
): AutomationWorkflow[] {
  const igla = szukane.trim().toLocaleLowerCase('pl-PL');
  if (igla === '') return [...petle];
  return petle.filter((petla) => {
    const stog = [petla.name, petla.id, petla.description ?? '']
      .join(' ')
      .toLocaleLowerCase('pl-PL');
    return stog.includes(igla);
  });
}

/**
 * Zdanie przy pozycji: liczba kroków, stan i wersja definicji. Liczba kroków
 * bierze się z długości `steps`; pole jest nieobowiązkowe, więc jego brak znaczy
 * zero kroków, a nie wartość nieznaną — definicja bez kroków jest definicją pustą.
 */
export function opisPetli(petla: AutomationWorkflow): string {
  const kroki = petla.steps?.length ?? 0;
  const czesci = [
    `${kroki} ${odmianaKrokow(kroki)}`,
    petla.enabled ? 'czynna' : 'wyłączona',
    `wersja ${petla.version ?? 1}`,
    `zmieniona ${new Date(petla.updatedAt).toLocaleString('pl-PL')}`,
  ];
  const opis = petla.description ?? '';
  return opis === '' ? czesci.join(' · ') : `${czesci.join(' · ')} — ${opis}`;
}

/** Odmiana rzeczownika „krok” przez liczbę, według reguł polskiej liczebności. */
function odmianaKrokow(liczba: number): string {
  const reszta = liczba % 10;
  const setka = liczba % 100;
  if (liczba === 1) return 'krok';
  if (reszta >= 2 && reszta <= 4 && (setka < 12 || setka > 14)) return 'kroki';
  return 'kroków';
}

/** Wykaz gotowy do osadzenia albo zdanie o pustce — rozstrzygnięcie dla okna. */
export type WidokWykazuPetli =
  | { rodzaj: 'wykaz'; element: HTMLElement }
  | { rodzaj: 'pustka'; zdanie: string };

/**
 * Składa wykaz pętli albo zdanie o pustce. Pustka bywa dwojaka i każda ma
 * osobne zdanie: rdzeń bez ani jednej zapisanej pętli to inny brak niż wykaz
 * zawężony napisem, który do niczego nie pasuje.
 */
export function widokWykazuPetli(
  petle: readonly AutomationWorkflow[],
  szukane: string,
  uruchom: (petla: AutomationWorkflow) => void,
): WidokWykazuPetli {
  if (petle.length === 0) {
    return {
      rodzaj: 'pustka',
      zdanie:
        'Rdzeń nie ma ani jednej zapisanej pętli. Wykaz bierze się z zapisanych definicji automatyk, ' +
        'a żadnej jeszcze nie zapisano. Zbuduj kroki w Workflow Builderze i zapisz definicję — ' +
        'pozycja pojawi się tutaj.',
    };
  }
  const widoczne = dopasowanePetle(petle, szukane);
  if (widoczne.length === 0) {
    return {
      rodzaj: 'pustka',
      zdanie:
        `Żadna z ${petle.length} pętli wykazu nie ma w nazwie napisu „${szukane.trim()}”. ` +
        'Szukanie idzie po nazwie, identyfikatorze i opisie definicji. ' +
        'Wyczyść pole szukania albo wpisz krótszy fragment nazwy.',
    };
  }
  return { rodzaj: 'wykaz', element: listaPetli(widoczne, uruchom) };
}

/**
 * Buduje wykaz pozycji wraz z przyciskiem uruchomienia przy każdej. Stan pętli
 * ląduje w `data-stan-petli`, żeby arkusz i sprawdziany widoku mogły się o niego
 * zaczepić bez czytania treści — tak samo jak stan przebiegu w Execution
 * Monitorze.
 */
function listaPetli(
  petle: readonly AutomationWorkflow[],
  uruchom: (petla: AutomationWorkflow) => void,
): HTMLElement {
  const lista = wykaz('Gotowe pętle do uruchomienia', 'da-wykaz');
  for (const petla of petle) {
    const pozycja = pozycjaWykazu(petla.name, opisPetli(petla), 'da');
    pozycja.element.dataset['stanPetli'] = petla.enabled ? 'czynna' : 'wylaczona';
    pozycja.element.dataset['petla'] = petla.id;
    const start = przycisk('Uruchom', 'dn-btn dn-btn--atrament');
    start.addEventListener('click', () => uruchom(petla));
    pozycja.akcje.append(start);
    lista.append(pozycja.element);
  }
  return lista;
}
