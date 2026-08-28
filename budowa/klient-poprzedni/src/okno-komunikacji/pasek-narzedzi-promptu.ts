import type { NarzedziePromptu, ProfilModulu } from './profil-modulu';

/** Pasek narzędzi promptu przestawiany przy każdej zmianie modułu na komplet przycisków właściwy temu modułowi. */
export interface PasekNarzedziPromptu {
  /** Element montowany nad polem wypowiedzi. */
  element: HTMLElement;
  /** Buduje pasek na nowo z narzędzi wskazanego profilu. */
  ustaw(profil: ProfilModulu): void;
}

/**
 * Pasek narzędzi promptu jednego modułu wstawia gotowe polecenie do pola wypowiedzi po kliknięciu, a wykonanie zleca operator modelowi, który sięga po narzędzie platformy — nic nie dzieje się bez jego wiedzy.
 */
export function utworzPasekNarzedzi(
  naWybor: (narzedzie: NarzedziePromptu) => void,
): PasekNarzedziPromptu {
  const element = document.createElement('div');
  element.className = 'dc-pasek-narzedzi';
  element.setAttribute('role', 'toolbar');
  element.setAttribute('aria-label', 'Narzędzia promptu modułu');

  return {
    element,
    ustaw(profil) {
      element.dataset['modul'] = profil.kod;
      element.replaceChildren(...profil.narzedzia.map((n) => przycisk(n, naWybor)));
    },
  };
}

/** Pojedyncza pozycja paska narzędzi: etykieta widoczna operatorowi oraz zdanie przeznaczenia w podpowiedzi. */
function przycisk(
  narzedzie: NarzedziePromptu,
  naWybor: (narzedzie: NarzedziePromptu) => void,
): HTMLButtonElement {
  const element = document.createElement('button');
  element.type = 'button';
  element.className = 'dn-btn dn-btn--zarys dn-btn--sm dc-pasek-narzedzi__pozycja';
  element.dataset['narzedzie'] = narzedzie.kod;
  element.title = narzedzie.opis;
  element.setAttribute('aria-label', `${narzedzie.nazwa} — ${narzedzie.opis}`);
  element.textContent = narzedzie.nazwa;
  element.addEventListener('click', () => naWybor(narzedzie));
  return element;
}
