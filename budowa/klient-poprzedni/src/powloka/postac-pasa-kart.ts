import type { MigawkaKart, StanZrodlaKart } from './zrodlo-kart-sesji';

// Postać pasa kart poza samymi kartami: stan pusty oraz komunikat odpowiedzi rdzenia.

/** Napisy pasa kart bez ani jednej karty, po jednym na każdy możliwy stan źródła kart sesji tego rdzenia. */
const PUSTE: Readonly<Record<StanZrodlaKart, string>> = {
  oczekiwanie: 'Oczekiwanie na rdzeń…',
  gotowe: 'Rdzeń nie ma otwartej sesji',
  blad: 'Rdzeń odmówił wykazu sesji',
};

export interface PostacPasa {
  /** Napis stanu pustego, montowany w pasie obok wykazu kart. */
  pustka: HTMLElement;
  /** Treść odpowiedzi rdzenia na czynność Operatora. */
  komunikat: HTMLElement;
  /** Przerysowuje stan pusty; pas z kartami go chowa. */
  ustawPustke(migawka: MigawkaKart, liczbaKart: number): void;
  /** Pokazuje treść odpowiedzi rdzenia. */
  zglos(tekst: string, waga: 'blad' | 'info'): void;
  /** Zdejmuje komunikat — robi to każda kolejna migawka rdzenia. */
  schowajKomunikat(): void;
}

export function utworzPostacPasa(): PostacPasa {
  const pustka = document.createElement('p');
  pustka.className = 'dn-sesje__pustka';
  pustka.hidden = true;

  const komunikat = document.createElement('p');
  komunikat.className = 'dn-plakietka dn-sesje__komunikat';
  komunikat.hidden = true;

  return {
    pustka,
    komunikat,

    ustawPustke(migawka, liczbaKart) {
      const tresc = migawka.stan === 'blad' ? (migawka.blad ?? PUSTE.blad) : PUSTE[migawka.stan];
      pustka.textContent = tresc;
      pustka.hidden = liczbaKart > 0;
    },

    zglos(tekst, waga) {
      komunikat.textContent = tekst;
      // Pełna treść zostaje w dymku przeglądarki, bo pasek przycina napis.
      komunikat.title = tekst;
      komunikat.classList.toggle('dn-plakietka--blad', waga === 'blad');
      komunikat.classList.toggle('dn-plakietka--informacja', waga === 'info');
      komunikat.setAttribute('role', waga === 'blad' ? 'alert' : 'status');
      komunikat.hidden = false;
    },

    schowajKomunikat() {
      komunikat.hidden = true;
    },
  };
}
