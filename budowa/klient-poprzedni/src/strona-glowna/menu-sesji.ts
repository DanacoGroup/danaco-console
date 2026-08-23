import { utworzDymekObjasnienia } from '../komponenty/dymek';
import { czynnosciWiersza, type CzynnosciSesji, type MeldunekCzynnosci } from './czynnosci-sesji';
import type { WpisSesji } from './zrodlo-sesji';

/**
 * Menu czynności jednego wiersza historii sesji.
 *
 * Buduje przyciski czynności sensownych dla tego wiersza i prowadzi jedną
 * z nich od naciśnięcia do meldunku. O tym, które czynności są sensowne,
 * rozstrzyga `czynnosci-sesji`; ten plik nie zna stanów sesji.
 *
 * Pusty wykaz czynności daje `null` zamiast pustego pojemnika — wiersz czysto
 * informacyjny nie dostaje ramki po menu, którego nie ma.
 *
 * Na czas wykonania wiersz nie przyjmuje drugiej czynności — także z innego
 * przycisku niż naciśnięty: archiwizacja w trakcie zmiany nazwy dawałaby dwa
 * żądania o tę samą sesję o przypadkowej kolejności skutków.
 *
 * Zajętość nie jest bramą i nie gasi przycisków: atrybutu `disabled` nie
 * stawiamy nigdzie w produkcie. Przycisk zostaje klikalny, a stan niesie napis
 * — przycisk czynny dopisuje „…”, więc zajętość jest widoczna bez samego
 * koloru. Naciśnięcie w trakcie biegu wraca meldunkiem, nie ciszą.
 *
 * Odmowa wraca meldunkiem, nie wyjątkiem: wykonanie oddaje zdanie albo `null`,
 * a ten plik podaje je dalej. Wyjątek nieprzewidziany też kończy się meldunkiem,
 * żeby wiersz nie został z zablokowanymi przyciskami.
 */

export interface MenuSesji {
  element: HTMLElement;
}

/** Odbiorca zdania o wyniku czynności; strefa pokazuje je nad wykazem. */
export type OdbiorcaMeldunku = (tekst: string) => void;

export function utworzMenuSesji(
  wpis: WpisSesji,
  czynnosci: CzynnosciSesji,
  meldunek: OdbiorcaMeldunku,
): MenuSesji | null {
  const pozycje = czynnosciWiersza(wpis, czynnosci);
  if (pozycje.length === 0) return null;

  const element = document.createElement('div');
  element.className = 'dn-strona__sesja-czynnosci';
  element.setAttribute('role', 'group');
  element.setAttribute('aria-label', `Czynności sesji ${wpis.sesja.title ?? wpis.sesja.id}`);

  let zajete = false;

  /**
   * Zajętość widoczna dla technologii wspomagających. `aria-busy` opisuje bieg
   * czynności, a nie odbiera pozycji dostępności — inaczej niż `disabled`.
   */
  function ustawZajetosc(biegnie: boolean): void {
    zajete = biegnie;
    element.setAttribute('aria-busy', biegnie ? 'true' : 'false');
  }

  ustawZajetosc(false);

  for (const pozycja of pozycje) {
    const przycisk = document.createElement('button');
    przycisk.type = 'button';
    przycisk.className = `dn-btn dn-btn--sm${pozycja.postac === undefined ? '' : ` ${pozycja.postac}`}`;
    przycisk.textContent = pozycja.napis;
    przycisk.dataset.czynnosc = pozycja.klucz;

    przycisk.addEventListener('click', () => {
      // Naciśnięcie w trakcie biegu nie jest odrzucane po cichu — Operator
      // dostaje zdanie o tym, dlaczego czynność nie ruszyła teraz.
      if (zajete) {
        meldunek('Czynność tego wiersza jeszcze biegnie — poczekaj na jej meldunek.');
        return;
      }
      ustawZajetosc(true);
      przycisk.textContent = `${pozycja.napis}…`;
      void wykonaj(pozycja.wykonaj(wpis), meldunek).finally(() => {
        przycisk.textContent = pozycja.napis;
        ustawZajetosc(false);
      });
    });

    element.append(przycisk, utworzDymekObjasnienia(pozycja.wyjasnienie));
  }

  return { element };
}

/** Doprowadza czynność do meldunku — także wtedy, gdy rzuciła wyjątkiem. */
async function wykonaj(
  bieg: Promise<MeldunekCzynnosci>,
  meldunek: OdbiorcaMeldunku,
): Promise<void> {
  try {
    const tresc = await bieg;
    if (tresc !== null) meldunek(tresc);
  } catch (powod) {
    meldunek(powod instanceof Error ? powod.message : 'Czynność nie doszła do skutku.');
  }
}
