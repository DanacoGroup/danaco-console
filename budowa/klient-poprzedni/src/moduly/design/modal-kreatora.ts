/**
 * Powłoka modala kreatora: rama okna bez wiedzy o jego treści, z czterema
 * równorzędnymi drogami zamknięcia — kontrolka nagłówka, klawisz Escape,
 * kliknięcie w nakładkę oraz akcja stopki.
 */
export interface ModalKreatora {
  element: HTMLDialogElement;
  /** Miejsce na treść kreatora. */
  cialo: HTMLElement;
  /** Miejsce na komunikat blokowy nad stopką z akcjami. */
  komunikat: HTMLElement;
  /** Przycisk główny stopki. */
  glowny: HTMLButtonElement;
  otworz(): void;
  zamknij(): void;
  /** Czy modal jest otwarty — do sprawdzianu stanu `zamknięty/otwarty`. */
  otwarty(): boolean;
}

export function utworzModalKreatora(tytul: string, akcja: string): ModalKreatora {
  const element = document.createElement('dialog');
  element.className = 'dn-modal md-kreator';
  element.setAttribute('aria-label', tytul);

  const naglowekTytul = document.createElement('h3');
  naglowekTytul.className = 'dn-modal-tytul';
  naglowekTytul.textContent = tytul;

  const zamkniecie = document.createElement('button');
  zamkniecie.type = 'button';
  zamkniecie.className = 'dn-btn-ikona';
  zamkniecie.textContent = '×';
  zamkniecie.setAttribute('aria-label', 'Zamknij kreator');

  const naglowek = document.createElement('header');
  naglowek.className = 'dn-modal-naglowek';
  naglowek.append(naglowekTytul, zamkniecie);

  const cialo = document.createElement('div');
  cialo.className = 'dn-modal-cialo md-kreator__cialo';

  const komunikat = document.createElement('p');
  komunikat.className = 'md-kreator__komunikat';
  komunikat.hidden = true;

  const glowny = document.createElement('button');
  glowny.type = 'button';
  glowny.className = 'dn-btn dn-btn--sygnal';
  glowny.textContent = akcja;

  const zamknij = document.createElement('button');
  zamknij.type = 'button';
  zamknij.className = 'dn-btn dn-btn--zarys';
  zamknij.textContent = 'Zamknij';

  const stopka = document.createElement('footer');
  stopka.className = 'dn-modal-stopka';
  stopka.append(glowny, zamknij);

  element.append(naglowek, cialo, komunikat, stopka);

  for (const kontrolka of [zamkniecie, zamknij]) {
    kontrolka.addEventListener('click', () => element.close());
  }
  // Nakładka to sam dialog: wnętrze przykrywa go w całości i zatrzymuje
  // kliknięcie na sobie.
  element.addEventListener('click', (zdarzenie) => {
    if (zdarzenie.target === element) element.close();
  });

  return {
    element,
    cialo,
    komunikat,
    glowny,

    otworz() {
      if (element.open) return;
      // `showModal` daje nakładkę, stos okien i klawisz Escape; bez niego
      // zostaje zwykłe otwarcie.
      if (typeof element.showModal === 'function') element.showModal();
      else element.open = true;
    },

    zamknij() {
      element.close();
    },

    otwarty: () => element.open,
  };
}
