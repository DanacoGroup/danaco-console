/**
 * Jeden wpis pamięci w wykazie Memory & Context Manager wraz z czynnościami
 * dostępnymi na nim. Plik odpowiada wyłącznie za zamianę wpisu w pozycję
 * wykazu; wiersz niczego nie wywołuje, a zamiar oddaje oknu prowadzącemu
 * wywołania `memory.*`.
 */
import {
  MemoryEntryOrigin,
  type ConfigScope,
  type WorkspaceMemoryEntry,
} from '../../../../shared/contract';
import { pozycjaWykazu, przyciskAkcji } from '../../modele/kontrolki-formularza';
import { chwila, NAZWY_POCHODZEN, NAZWY_ZASIEGOW } from './etykiety-assistant';

/**
 * Zapis wpisu zamówiony z wiersza: treść, zasięg, znacznik przypięcia,
 * pochodzenie oraz identyfikator wpisu, który okno przekazuje do wywołania
 * zapisu pamięci.
 */
export interface ZamowienieWpisu {
  tresc: string;
  zasieg: ConfigScope;
  przypiety: boolean;
  pochodzenie: MemoryEntryOrigin;
  wpis: string;
}

/**
 * Czynności okna dostępne wierszowi: zapis zamówionego wpisu, wczytanie wpisu
 * do edytora oraz usunięcie wpisu z pamięci obszaru roboczego.
 */
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

  // Propozycja modelu czeka na decyzję Operatora; przyjęcie zapisuje ją
  // z pochodzeniem `operator`.
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
