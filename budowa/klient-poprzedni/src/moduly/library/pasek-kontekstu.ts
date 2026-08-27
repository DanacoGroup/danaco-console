import { opisOdmowy } from '../../komponenty/odmowa';
import type { StanBiblioteki } from './stan-biblioteki';
import { KOD_MODULU, type ZrodloOtoczenia } from './zrodlo-otoczenia';

/**
 * Pasek kontekstu modułu Library, czyli zdanie o oknie komunikacji, w którego
 * imieniu moduł działa. Okno wybierane jest po oznaczeniu modułu, a jego brak
 * pasek opisuje wprost, zamiast podstawiać okno cudzego modułu.
 */
export interface PasekKontekstu {
  element: HTMLElement;
  /** Odczytuje okno sesji i jego stan; wolno wołać wielokrotnie. */
  wczytaj(idSesji: string): Promise<void>;
}

export function utworzPasekKontekstu(
  stan: StanBiblioteki,
  otoczenie: ZrodloOtoczenia,
): PasekKontekstu {
  const element = document.createElement('p');
  element.className = 'ml-kontekst';
  element.dataset['kontekst'] = 'nieznany';
  element.textContent = 'Kontekst okna nieodczytany.';

  return {
    element,

    async wczytaj(idSesji) {
      const okna = await otoczenie.okna(idSesji);
      if (!okna.udany || okna.wynik === undefined) {
        stan.ustawOkno('');
        element.dataset['kontekst'] = 'brak';
        element.textContent = opisOdmowy('Odczyt okien sesji', okna.blad?.code, okna.blad?.message);
        return;
      }
      const wykaz = okna.wynik.windows;
      const wlasne = wykaz.find((okno) => okno.moduleId === KOD_MODULU);
      if (wlasne === undefined) {
        stan.ustawOkno('');
        element.dataset['kontekst'] = 'brak';
        // Liczba okien odróżnia pustą odpowiedź rdzenia od odpowiedzi bez okna modułu.
        element.textContent =
          wykaz.length === 0
            ? 'Rdzeń nie ma ani jednego okna komunikacji tej sesji — przeniesienie kontekstu ' +
              'nie ma okna źródłowego.'
            : `Rdzeń oddał okna sesji (${wykaz.length}), ale żadne nie należy do modułu ` +
              'Library — moduł nie działa w imieniu żadnego z nich i nie ma okna źródłowego ' +
              'przeniesienia kontekstu.';
        return;
      }
      stan.ustawOkno(wlasne.id);
      const wynik = await otoczenie.stanOkna(wlasne.id);
      if (!wynik.udany || wynik.wynik === undefined) {
        element.dataset['kontekst'] = 'blad';
        element.textContent = opisOdmowy('Odczyt stanu okna', wynik.blad?.code, wynik.blad?.message);
        return;
      }
      const okno = wynik.wynik.window;
      element.dataset['kontekst'] = 'gotowy';
      // Moduł okna i tryb uprawnień stoją w panelu sterowania okna, więc pasek ich nie powtarza.
      element.textContent =
        `Okno kontekstu: ${okno.title ?? okno.id} · wiadomości ${wynik.wynik.messageCount}.`;
    },
  };
}
