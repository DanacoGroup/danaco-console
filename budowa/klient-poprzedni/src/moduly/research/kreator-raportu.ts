import { poleTekstowe, przycisk } from '../../modele/kontrolki-formularza';
import { zDymkiem } from './dymek-badania';
import { utworzSekcjeRaportu, type SekcjeRaportu } from './sekcje-raportu';
import { utworzStanOknaBadania, type StanOknaBadania } from './stan-okna-badania';

/**
 * Kreator raportu — modal z krokami pionowymi i stanami `zamknięty ·
 * otwierający się · otwarty · próba z brakami · ładowanie · błąd`.
 *
 * Cztery drogi zamknięcia, wszystkie z natywnego `<dialog>`: kontrolka
 * w nagłówku, klawisz Escape, kliknięcie w nakładkę i akcja w stopce. Wzorzec
 * powtórzony za `konfiguracja/okno-konfiguracji.ts`.
 *
 * Przycisk główny jest czynny od otwarcia. Próba z brakami nie gasi kontrolki:
 * kreator zostaje otwarty, a brak wraca komunikatem nad stopką. Ponowne
 * naciśnięcie w trakcie wywołania jest bezpieczne po stronie logiki kreatora,
 * nie przez odebranie klikalności.
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

  // Wnętrze kreatora jest gotowe od pierwszej chwili: pole tytułu, reguła dwóch
  // dróg i redakcja sekcji stoją tu bez żadnego wywołania rdzenia. Pas stanu
  // startuje w fazie `puste`, więc bez tego jednego zdania nad gotowym
  // formularzem wisiałby komunikat o pustce, której nie ma.
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

  // Czwarta droga zamknięcia: kliknięcie w nakładkę. Natywny `<dialog>` kieruje
  // je na sam element, więc odróżnia je od kliknięcia w treść kreatora.
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
 * Reguła, po której rdzeń wybiera drogę budowy — powiedziana Operatorowi przed
 * naciśnięciem, nie dopiero w wyniku.
 *
 * `research.report.build` z polem `sections` składa raport z samych sekcji
 * podanych i pomija `findingIds`; bez tego pola woła kanał modelu po sekcję
 * nadrzędną. Reguła stoi w kreatorze, bo tu Operator rozstrzyga ją każdym
 * wpisanym znakiem.
 *
 * Droga modelu nie gwarantuje odpowiedzi modelu: bez logowania do programu
 * `claude` budowa kończy się powodzeniem, raport zostaje zapisany, a treścią
 * sekcji nadrzędnej jest komunikat procesu kanału. Kreator uprzedza o tym przed
 * naciśnięciem, bo po nim klient nie ma już czym tego rozpoznać — odpowiedź nie
 * niesie znaku pochodzenia treści (`skutek-zlozenia.ts`).
 */
function reguleDwochDrog(): HTMLElement {
  const element = document.createElement('p');
  element.className = 'dn-pole-opis mr-kreator__regula';
  element.textContent =
    'Wiersz wypełniony ma pierwszeństwo: gdy w redakcji stoi choć jedna sekcja z tytułem albo treścią, rdzeń składa raport z NIEJ i nie zamienia zaznaczonych ustaleń na sekcje. Redakcja pusta kieruje budowę na drogę modelu — rdzeń woła kanał modelu po sekcję nadrzędną z ustaleń zaznaczonych w Findings Panel. Uwaga: rdzeń bierze z kanału każdą treść, jaka wróci, i zapisuje ją jako sekcję także wtedy, gdy jest to komunikat procesu kanału (na przykład przy braku logowania). Przeczytaj tę sekcję w podglądzie, zanim wydasz raport — po zapisie klient nie ma jak odróżnić jej od streszczenia.';
  return element;
}

/** Nagłówek kreatora: tytuł i kontrolka zamknięcia — pierwsza z czterech dróg. */
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
