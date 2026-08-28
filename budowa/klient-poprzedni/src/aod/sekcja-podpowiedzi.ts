import type { AodSuggestion, AodSuggestionRequest } from '../../../shared/contract';
import { opisOdmowyAod } from './odmowy-aod';
import type { ZrodloAod } from './zrodlo-komend';
import { utworzAkapit, utworzPodtytul } from './pola-wykazu';
import { DZIALANIA_RODZAJU, NAZWY_RODZAJOW, POWSTANIE_RODZAJU, RodzajSugestii } from './rodzaje-sugestii';

/** Lista oczekujących sugestii przychodzących z rdzenia, uszeregowana od bytu najwęższego do platformy. */
export interface SekcjaPodpowiedzi {
  element: HTMLElement;
  /** Pyta rdzeń o podpowiedzi dla wskazanego okna, sesji albo procesu. */
  odswiez(zadanie: AodSuggestionRequest): Promise<void>;
}

/** Górna granica liczby podpowiedzi pokazywanych naraz — nakładka jest wąska, więc wykaz musi być krótki. */
const GRANICA = 8;

export function utworzSekcjePodpowiedzi(zrodlo: ZrodloAod): SekcjaPodpowiedzi {
  const element = document.createElement('section');
  element.className = 'ao-sekcja';

  const miejsce = document.createElement('div');
  miejsce.className = 'ao-sekcja__tresc';

  element.append(
    utworzPodtytul('Lista oczekujących sugestii (aod.suggestion)'),
    miejsce,
    utworzKatalogRodzajow(),
  );

  return {
    element,

    async odswiez(zadanie) {
      miejsce.replaceChildren(utworzAkapit('ao-pusto', 'Pytam rdzeń o podpowiedzi…'));

      const wynik = await zrodlo.podpowiedzi({ ...zadanie, limit: zadanie.limit ?? GRANICA });

      if (!wynik.udany || wynik.wynik === undefined) {
        miejsce.replaceChildren(
          utworzAkapit(
            'ao-odmowa',
            opisOdmowyAod('Odczyt podpowiedzi', 'podpowiedzi', wynik.blad),
          ),
        );
        return;
      }

      const podpowiedzi = wynik.wynik.suggestions;
      if (podpowiedzi.length === 0) {
        miejsce.replaceChildren(
          utworzAkapit(
            'ao-pusto',
            'Rdzeń nie ma teraz żadnej podpowiedzi. Pusty wykaz jest stanem poprawnym — ' +
              'podpowiedzi przybywa wraz z postępem pracy.',
          ),
        );
        return;
      }

      miejsce.replaceChildren(zbudujWykaz(podpowiedzi));
    },
  };
}

/** Katalog rodzajów sugestii wypisany pod wykazem, bo pojedynczej pozycji nie da się do rodzaju przypisać. */
function utworzKatalogRodzajow(): HTMLElement {
  const katalog = document.createElement('details');
  katalog.className = 'ao-katalog';

  const naglowek = document.createElement('summary');
  naglowek.className = 'ao-katalog__naglowek';
  naglowek.textContent = 'Katalog rodzajów sugestii (rozdz. 4 opracowania)';
  katalog.append(naglowek);

  katalog.append(
    utworzAkapit(
      'ao-granica',
      'Rodzaju sugestii kontrakt nie niesie: struktura AodSuggestion ma id, treść, ' +
        'proponowaną komendę, okno i czas powstania — pola rodzaju w niej nie ma. ' +
        'Pozycje wykazu powyżej są więc bez rodzaju, a nakładka rodzaju nie zgaduje. ' +
        'Poniżej stoi katalog, którym rodzaj będzie się posługiwał, gdy kontrakt go poniesie.',
    ),
  );

  const wykaz = document.createElement('dl');
  wykaz.className = 'ao-pola';
  for (const rodzaj of Object.values(RodzajSugestii)) {
    const nazwa = document.createElement('dt');
    nazwa.textContent = NAZWY_RODZAJOW[rodzaj];

    const opis = document.createElement('dd');
    opis.textContent =
      `${POWSTANIE_RODZAJU[rodzaj]} Działania: ` +
      `${DZIALANIA_RODZAJU[rodzaj].map((dzialanie) => dzialanie.nazwa).join(' · ')}.`;

    wykaz.append(nazwa, opis);
  }
  katalog.append(wykaz);

  return katalog;
}

function zbudujWykaz(podpowiedzi: readonly AodSuggestion[]): HTMLElement {
  const lista = document.createElement('ol');
  lista.className = 'ao-podpowiedzi';
  for (const podpowiedz of podpowiedzi) {
    lista.append(zbudujPozycje(podpowiedz));
  }
  return lista;
}

function zbudujPozycje(podpowiedz: AodSuggestion): HTMLLIElement {
  const pozycja = document.createElement('li');
  pozycja.className = 'ao-podpowiedz';

  const tresc = document.createElement('p');
  tresc.className = 'ao-podpowiedz__tresc';
  tresc.textContent = podpowiedz.text;
  pozycja.append(tresc);

  const stopka = document.createElement('p');
  stopka.className = 'ao-podpowiedz__opis';
  stopka.textContent = opisPozycji(podpowiedz);
  pozycja.append(stopka);

  // Rodzaj i waga sugestii: pól tych kontrakt nie niesie, pozycja mówi to wprost.
  const rodzaj = document.createElement('p');
  rodzaj.className = 'ao-podpowiedz__opis';
  rodzaj.textContent = 'rodzaj i waga: nie do odczytania — kontrakt nie niesie tych pól';
  pozycja.append(rodzaj);

  return pozycja;
}

/**
 * Podpis pozycji: czego podpowiedź dotyczy i kiedy powstała.
 *
 * Podpowiedź bywa zdaniem bez komendy — wtedy podpis o komendzie milczy.
 */
function opisPozycji(podpowiedz: AodSuggestion): string {
  const czesci: string[] = [];
  const komenda = (podpowiedz.commandType ?? '').trim();
  if (komenda !== '') czesci.push(`proponowana komenda: ${komenda}`);
  const okno = (podpowiedz.windowId ?? '').trim();
  if (okno !== '') czesci.push(`okno: ${okno}`);
  czesci.push(`powstała ${new Date(podpowiedz.createdAt).toLocaleString()}`);
  return czesci.join(' · ');
}
