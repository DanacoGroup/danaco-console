import { utworzDymekObjasnienia } from '../komponenty/dymek';
import { czynnosciWiersza, type CzynnosciSesji, type MeldunekCzynnosci } from './czynnosci-sesji';
import type { WpisSesji } from './zrodlo-sesji';

// Menu czynności jednego wiersza historii sesji, prowadzące wybraną czynność do meldunku wyniku.

export interface MenuSesji {
  element: HTMLElement;
}

/** Odbiorca zdania o wyniku czynności menu; strefa komunikatów pokazuje je operatorowi nad wykazem sesji. */
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

  // Zajętość widoczna dla technologii wspomagających przez atrybut zajętości, bez odbierania dostępu.
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
      // Naciśnięcie w trakcie biegu nie jest odrzucane po cichu — operator dostaje zdanie.
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

/** Doprowadza czynność do meldunku wyniku — także wtedy, gdy wykonanie rzuciło nieprzewidzianym wyjątkiem. */
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
