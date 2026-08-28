import { poleTekstowe, przycisk } from '../../modele/kontrolki-formularza';
import { zDymkiem } from './dymek-badania';
import { utworzSekcjeRaportu, type SekcjeRaportu } from './sekcje-raportu';
import { utworzStanOknaBadania, type StanOknaBadania } from './stan-okna-badania';

/**
 * Kreator raportu jest modalem z krokami pionowymi, dostępnym czterema drogami zamknięcia, w tym klawiszem Escape i kliknięciem w nakładkę.
 */
export interface KreatorRaportu {
  element: HTMLDialogElement;
  /** Sekcje w redakcji — okno dokłada je do żądania. */
  sekcje: SekcjeRaportu;
  /** Otwiera kreator; wolno wołać wielokrotnie. */
  otworz(tytul: string): void;
  /** Zamyka kreator; treść redakcji zostaje. */
  zamknij(): void;
  /** Tytuł raportu wpisany przez Operatora. */
  tytul(): string;
  /** Stan wnętrza kreatora: ładowanie, błąd, gotowość. */
  stan: StanOknaBadania;
}

export function utworzKreatorRaportu(naZlozenie: () => void): KreatorRaportu {
  const stan = utworzStanOknaBadania();
  const sekcje = utworzSekcjeRaportu();

  const poleTytulu = poleTekstowe({ etykieta: 'Tytuł raportu', podpowiedz: 'nazwa dokumentu końcowego' });

  const dodajSekcje = przycisk('+ Dodaj sekcję', 'dn-btn dn-btn--sm dn-btn--zarys');
  dodajSekcje.addEventListener('click', () => sekcje.dodaj());

  stan.tresc.append(
    zDymkiem(poleTytulu.element, 'Pole title komendy research.report.build; puste nie idzie do rdzenia.'),
    reguleDwochDrog(),
    dodajSekcje,
    sekcje.element,
  );

  // Stan startuje pusty; wołanie ustawia gotowość, by nie pokazać komunikatu o pustym wnętrzu.
  stan.gotowe();

  const element = document.createElement('dialog');
  element.className = 'dn-modal mr-kreator';
  element.setAttribute('aria-label', 'Kreator raportu badania');

  const cialo = document.createElement('div');
  cialo.className = 'dn-modal-cialo mr-kreator__cialo';
  cialo.append(stan.element);

  const zloz = przycisk('Złóż raport', 'dn-btn dn-btn--sm dn-btn--atrament');
  zloz.addEventListener('click', () => naZlozenie());

  const zamknij = przycisk('Zamknij', 'dn-btn dn-btn--sm dn-btn--zarys');
  zamknij.addEventListener('click', () => element.close());

  const stopka = document.createElement('footer');
  stopka.className = 'dn-modal-stopka mr-kreator__stopka';
  stopka.append(zloz, zamknij);

  element.append(naglowek(() => element.close()), cialo, stopka);

  // Kliknięcie w nakładkę zamyka: zdarzenie trafia w element dialogu, a nie w jego treść.
  element.addEventListener('click', (zdarzenie) => {
    if (zdarzenie.target === element) element.close();
  });

  return {
    element,
    sekcje,
    stan,
    tytul: () => poleTytulu.kontrolka.value,

    otworz(tytul) {
      if (!element.isConnected) document.body.append(element);
      if (poleTytulu.kontrolka.value === '') poleTytulu.kontrolka.value = tytul;
      if (!element.open) element.showModal();
      if (sekcje.liczba() === 0) sekcje.dodaj();
    },

    zamknij: () => element.close(),
  };
}

/**
 * Funkcja opisuje regułę wyboru drogi budowy raportu, informując Operatora przed naciśnięciem, że wypełniona redakcja ma pierwszeństwo przed drogą modelu.
 */
function reguleDwochDrog(): HTMLElement {
  const element = document.createElement('p');
  element.className = 'dn-pole-opis mr-kreator__regula';
  element.textContent =
    'Wiersz wypełniony ma pierwszeństwo: gdy w redakcji stoi choć jedna sekcja z tytułem albo treścią, rdzeń składa raport z NIEJ i nie zamienia zaznaczonych ustaleń na sekcje. Redakcja pusta kieruje budowę na drogę modelu — rdzeń woła kanał modelu po sekcję nadrzędną z ustaleń zaznaczonych w Findings Panel. Uwaga: rdzeń bierze z kanału każdą treść, jaka wróci, i zapisuje ją jako sekcję także wtedy, gdy jest to komunikat procesu kanału (na przykład przy braku logowania). Przeczytaj tę sekcję w podglądzie, zanim wydasz raport — po zapisie klient nie ma jak odróżnić jej od streszczenia.';
  return element;
}

/** Funkcja tworzy nagłówek kreatora złożony z tytułu okna oraz kontrolki zamknięcia, pierwszej z czterech dostępnych dróg. */
function naglowek(naZamkniecie: () => void): HTMLElement {
  const element = document.createElement('header');
  element.className = 'dn-modal-naglowek mr-kreator__naglowek';

  const tytul = document.createElement('h4');
  tytul.className = 'dn-modal-tytul';
  tytul.textContent = 'Kreator raportu badania';

  const zamknij = przycisk('×', 'dn-btn-ikona mr-kreator__zamkniecie');
  zamknij.setAttribute('aria-label', 'Zamknij kreator raportu');
  zamknij.addEventListener('click', naZamkniecie);

  element.append(tytul, zamknij);
  return element;
}
