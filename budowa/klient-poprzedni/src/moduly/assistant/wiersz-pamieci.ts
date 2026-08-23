import {
  MemoryEntryOrigin,
  type ConfigScope,
  type WorkspaceMemoryEntry,
} from '../../../../shared/contract';
import { pozycjaWykazu, przyciskAkcji } from '../../modele/kontrolki-formularza';
import { chwila, NAZWY_POCHODZEN, NAZWY_ZASIEGOW } from './etykiety-assistant';

/**
 * Jeden wpis pamięci w wykazie Memory & Context Manager wraz z czynnościami,
 * które można na nim wykonać.
 *
 * Plik odpowiada wyłącznie za zamianę wpisu w pozycję wykazu. Wiersz niczego
 * nie wywołuje — zamiar oddaje oknu, które prowadzi wywołania `memory.*`; ta
 * sama granica obowiązuje wiersz zlecenia (`wiersz-zlecenia.ts`).
 *
 * Wpis o pochodzeniu `model` jest propozycją czekającą na decyzję: przyjęcie
 * zapisuje tę samą treść z pochodzeniem `operator`, odrzucenie usuwa wpis.
 * Rozróżnienie należy do kontraktu (`MemoryEntryOrigin`), nie do okna — bez
 * niego pamięć zapisana przez model byłaby nie do odróżnienia od ustalenia
 * Operatora.
 */

/** Zapis wpisu zamówiony z wiersza. */
export interface ZamowienieWpisu {
  tresc: string;
  zasieg: ConfigScope;
  przypiety: boolean;
  pochodzenie: MemoryEntryOrigin;
  wpis: string;
}

/** Czynności okna dostępne wierszowi. */
export interface CzynnosciWpisu {
  zapisz(zamowienie: ZamowienieWpisu): void;
  wczytajDoEdytora(wpis: WorkspaceMemoryEntry): void;
  usun(wpis: WorkspaceMemoryEntry): void;
}

export function wierszPamieci(
  wpis: WorkspaceMemoryEntry,
  czynnosci: CzynnosciWpisu,
): HTMLElement {
  const opis = [
    `zasięg: ${NAZWY_ZASIEGOW[wpis.scope]}`,
    `pochodzenie: ${NAZWY_POCHODZEN[wpis.origin]}`,
    wpis.pinned === true ? 'przypięty' : 'nieprzypięty',
    `zmieniony: ${chwila(wpis.updatedAt)}`,
  ].join(' · ');

  const { element, akcje } = pozycjaWykazu(wpis.content, opis, 'ma');
  element.dataset['wpis'] = wpis.id;
  element.dataset['pochodzenie'] = wpis.origin;

  akcje.append(
    guzik('Edytuj', () => czynnosci.wczytajDoEdytora(wpis)),
    guzik(wpis.pinned === true ? 'Odepnij' : 'Przypnij', () =>
      czynnosci.zapisz({
        tresc: wpis.content,
        zasieg: wpis.scope,
        przypiety: wpis.pinned !== true,
        pochodzenie: wpis.origin,
        wpis: wpis.id,
      }),
    ),
  );

  // Propozycja modelu czeka na decyzję Operatora. Przyjęcie jest zapisem tej
  // samej treści z pochodzeniem `operator`; odrzucenie — usunięciem wpisu.
  if (wpis.origin === MemoryEntryOrigin.Model) {
    akcje.append(
      guzik('Przyjmij propozycję', () =>
        czynnosci.zapisz({
          tresc: wpis.content,
          zasieg: wpis.scope,
          przypiety: wpis.pinned === true,
          pochodzenie: MemoryEntryOrigin.Operator,
          wpis: wpis.id,
        }),
      ),
    );
  }

  akcje.append(guzik('Usuń', () => czynnosci.usun(wpis)));
  return element;
}

function guzik(etykieta: string, naKlik: () => void): HTMLButtonElement {
  const element = przyciskAkcji(etykieta, 'dn-btn dn-btn--sm dn-btn--zarys');
  element.addEventListener('click', naKlik);
  return element;
}
