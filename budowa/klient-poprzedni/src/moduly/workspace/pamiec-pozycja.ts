import {
  MemoryEntryOrigin,
  type ConfigScope,
  type WorkspaceMemoryEntry,
} from '../../../../shared/contract';
import { pozycjaWykazu, przyciskAkcji as przycisk } from '../../modele/kontrolki-formularza';

/** Jedna pozycja wykazu pamięci niesie edycję, przypięcie, przyjęcie propozycji i scalenie treści. */

/** Zapis wpisu widziany przez pozycję wykazu niesie treść, zasięg, stan przypięcia, pochodzenie oraz identyfikator wpisu. */
export interface ZapisWpisu {
  tresc: string;
  zasieg: ConfigScope;
  przypiety: boolean;
  pochodzenie: MemoryEntryOrigin;
  wpis?: string;
}

/** Czynności okna dostępne pozycji obejmują zapis, wczytanie i dołączenie do edytora oraz usunięcie wpisu pamięci. */
export interface CzynnosciPozycji {
  zapisz(zadanie: ZapisWpisu): void;
  wczytajDoEdytora(wpis: WorkspaceMemoryEntry): void;
  dolaczDoEdytora(wpis: WorkspaceMemoryEntry): void;
  usun(wpis: WorkspaceMemoryEntry): void;
}

export function pozycjaPamieci(
  wpis: WorkspaceMemoryEntry,
  czynnosci: CzynnosciPozycji,
): HTMLElement {
  const opis = [
    `zasięg: ${wpis.scope}`,
    `pochodzenie: ${wpis.origin}`,
    wpis.pinned === true ? 'przypięty' : 'nieprzypięty',
    new Date(wpis.updatedAt).toLocaleString('pl-PL'),
  ].join(' · ');
  const { element, akcje } = pozycjaWykazu(wpis.content, opis, 'dw');
  element.dataset['pochodzenie'] = wpis.origin;

  const edytuj = przycisk('Edytuj', 'dn-btn dn-btn--zarys');
  edytuj.addEventListener('click', () => czynnosci.wczytajDoEdytora(wpis));

  const przypnij = przycisk(wpis.pinned === true ? 'Odepnij' : 'Przypnij', 'dn-btn dn-btn--zarys');
  przypnij.addEventListener('click', () =>
    czynnosci.zapisz({
      tresc: wpis.content,
      zasieg: wpis.scope,
      przypiety: wpis.pinned !== true,
      pochodzenie: wpis.origin,
      wpis: wpis.id,
    }),
  );
  akcje.append(edytuj, przypnij);

  // Propozycja modelu czeka na decyzję operatora: przyjęcie zapisuje wpis operatora, odrzucenie usuwa.
  if (wpis.origin === MemoryEntryOrigin.Model) {
    const przyjmij = przycisk('Akceptuj propozycję', 'dn-btn dn-btn--atrament');
    przyjmij.addEventListener('click', () =>
      czynnosci.zapisz({
        tresc: wpis.content,
        zasieg: wpis.scope,
        przypiety: wpis.pinned === true,
        pochodzenie: MemoryEntryOrigin.Operator,
        wpis: wpis.id,
      }),
    );
    const odrzuc = przycisk('Odrzuć propozycję', 'dn-btn dn-btn--zarys');
    odrzuc.addEventListener('click', () => czynnosci.usun(wpis));
    akcje.append(przyjmij, odrzuc);
  }

  const scal = przycisk('Scal do edytora', 'dn-btn dn-btn--zarys');
  scal.addEventListener('click', () => czynnosci.dolaczDoEdytora(wpis));
  akcje.append(scal);

  // Usunięcie wpisu jest dostępne każdej pozycji; potwierdzenie mówi, co oddał rdzeń, nie okno.
  const usun = przycisk('Usuń', 'dn-btn dn-btn--zarys');
  usun.addEventListener('click', () => czynnosci.usun(wpis));
  akcje.append(usun);
  return element;
}

/** Pamięć projektu w postaci Markdown stanowi treść pliku eksportu, budowaną z wykazu wpisów pamięci projektu. */
export function eksportPamieci(wpisy: readonly WorkspaceMemoryEntry[]): string {
  const wiersze = wpisy.map(
    (wpis) =>
      `- [${wpis.pinned === true ? 'przypięty' : ' '}] (${wpis.scope}, ${wpis.origin}) ${wpis.content}`,
  );
  return ['# Pamięć projektu', '', ...wiersze, ''].join('\n');
}
