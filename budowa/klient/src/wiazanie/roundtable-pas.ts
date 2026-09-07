// Wspólne części pasów działań okna Roundtable.
import { oglos } from './ogloszenie.ts';

/* Pas stoi nad ciałem panelu, bo ciało jest wymieniane przy każdym odświeżeniu
   wykazu; pas postawiony w nim znikałby razem z wykazem. */
export function postawPas(korzen: Element, idPanelu: string, wezly: HTMLElement[]): void {
  const panel = korzen.querySelector(`#${idPanelu}`);
  const cialo = panel?.querySelector('.sta-okno-tresc');
  if (panel === null || panel === undefined || cialo === null || cialo === undefined) return;
  const pas = korzen.ownerDocument.createElement('div');
  pas.className = 'dn-pas-dzialan';
  pas.append(...wezly);
  panel.insertBefore(pas, cialo);
}

export function polePasa(korzen: Element, klucz: string, opis: string): HTMLInputElement {
  const wezel = korzen.ownerDocument.createElement('input');
  wezel.type = 'text';
  wezel.className = 'dn-pole dn-pole--sm';
  wezel.placeholder = opis;
  wezel.setAttribute('aria-label', opis);
  wezel.dataset[klucz] = '';
  return wezel;
}

export function przyciskPasa(
  korzen: Element,
  klucz: string,
  czynnosc: string,
  etykieta: string,
): HTMLButtonElement {
  const wezel = korzen.ownerDocument.createElement('button');
  wezel.type = 'button';
  wezel.className = 'dn-btn dn-btn--duch dn-btn--sm';
  wezel.dataset[klucz] = czynnosc;
  wezel.textContent = etykieta;
  return wezel;
}

export function wartoscPola(korzen: Element, selektor: string): string {
  const wezel = korzen.querySelector<HTMLInputElement>(selektor);
  return (wezel?.value ?? '').trim();
}

/* Czynność nieodwracalna wymaga dwóch naciśnięć: pierwsze uzbraja przycisk. */
export function potwierdzone(
  korzen: Element,
  selektor: string,
  etykietaUzbrojona: string,
  etykietaZwykla: string,
): boolean {
  const wezel = korzen.querySelector<HTMLElement>(selektor);
  if (wezel !== null && wezel.dataset.uzbrojone !== 'tak') {
    wezel.dataset.uzbrojone = 'tak';
    wezel.textContent = etykietaUzbrojona;
    return false;
  }
  if (wezel !== null) {
    delete wezel.dataset.uzbrojone;
    wezel.textContent = etykietaZwykla;
  }
  return true;
}

/* Przycisk bez obsługi wyglądałby na działający, więc okno odmawia wprost. */
export function odmowaCzynnosci(naglowek: string, czynnosc: string): void {
  oglos(naglowek, `Czynność „${czynnosc}” nie jest prowadzona przez to okno.`, 'ostrzezenie');
}
