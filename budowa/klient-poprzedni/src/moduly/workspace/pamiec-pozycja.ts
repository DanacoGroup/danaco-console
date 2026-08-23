import {
  MemoryEntryOrigin,
  type ConfigScope,
  type WorkspaceMemoryEntry,
} from '../../../../shared/contract';
import { pozycjaWykazu, przyciskAkcji as przycisk } from '../../modele/kontrolki-formularza';

/**
 * Jedna pozycja wykazu pamięci projektu wraz z czynnościami, które można na
 * niej wykonać: edycja, przypięcie, przyjęcie propozycji modelu i scalenie
 * treści.
 *
 * Osobny plik od okna, bo to inna odpowiedzialność: okno prowadzi odczyt
 * i stany, pozycja rysuje jeden wpis. Czynności przychodzą wstrzyknięte
 * — pozycja nie zna ani kanału, ani stanu okna.
 *
 * Usunięcie wpisu idzie komendą `memory.delete`: kasuje ona wpis założony przez
 * `workspace.context.set` i oddaje `deleted`, a usunięcie wpisu nieistniejącego
 * kończy odmową. Odrzucenie propozycji modelu jest właśnie takim usunięciem
 * i woła tę samą komendę.
 *
 * `memory.detach` nie jest tu wołany: znaczenie odpięcia nie jest w kontrakcie
 * ustalone, więc okno nie nadaje mu własnego sensu.
 */

/** Zapis wpisu widziany przez pozycję wykazu. */
export interface ZapisWpisu {
  tresc: string;
  zasieg: ConfigScope;
  przypiety: boolean;
  pochodzenie: MemoryEntryOrigin;
  wpis?: string;
}

/** Czynności okna dostępne pozycji. */
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

  // Propozycja modelu czeka na decyzję Operatora: przyjęcie jest zapisem tego
  // samego wpisu z pochodzeniem `operator`, a odrzucenie — usunięciem wpisu
  // komendą `memory.delete`.
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

  // Usunięcie wpisu — dostępne każdej pozycji, nie tylko propozycji modelu:
  // `memory.delete` kasuje wpis wskazany identyfikatorem, a potwierdzenie mówi,
  // co oddał rdzeń (`deleted`), nie co wysłało okno.
  const usun = przycisk('Usuń', 'dn-btn dn-btn--zarys');
  usun.addEventListener('click', () => czynnosci.usun(wpis));
  akcje.append(usun);
  return element;
}

/** Pamięć projektu w postaci Markdown — treść pliku eksportu. */
export function eksportPamieci(wpisy: readonly WorkspaceMemoryEntry[]): string {
  const wiersze = wpisy.map(
    (wpis) =>
      `- [${wpis.pinned === true ? 'przypięty' : ' '}] (${wpis.scope}, ${wpis.origin}) ${wpis.content}`,
  );
  return ['# Pamięć projektu', '', ...wiersze, ''].join('\n');
}
