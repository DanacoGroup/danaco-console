import type { AodContextGetResponse, ContextBundle } from '../../../shared/contract';
import type { Wynik } from '../protokol/kanal';
import { opisOdmowyAod } from './odmowy-aod';
import {
  BRAK,
  dodajPole,
  dodajPoleWykazu,
  utworzAkapit,
  utworzPodtytul,
  utworzWykazPol,
} from './pola-wykazu';

/**
 * Sekcja odczytu kontekstu okna ogniskowanego (`aod.context.get`).
 *
 * Każde pole kompletu dostaje własny wiersz z etykietą po polsku, pole puste
 * pisze się jako `(brak)`. Wszystkie pola `ContextBundle` są w kontrakcie
 * opcjonalne, więc komplet pusty jest stanem poprawnym, nie odmową — sekcja
 * mówi to osobnym zdaniem. `executionParams` ma w kontrakcie typ `unknown`,
 * więc sekcja go nie rozbiera i melduje jedynie obecność parametrów.
 */
export interface SekcjaKontekstu {
  element: HTMLElement;
  /** Nanosi odpowiedź rdzenia — wynik udany albo odmowę. */
  pokaz(wynik: Wynik<AodContextGetResponse>): void;
}

export function utworzSekcjeKontekstu(): SekcjaKontekstu {
  const element = document.createElement('section');
  element.className = 'ao-sekcja';

  const miejsce = document.createElement('div');
  miejsce.className = 'ao-sekcja__tresc';

  element.append(utworzPodtytul('Kontekst okna ogniskowanego (aod.context.get)'), miejsce);

  return {
    element,

    pokaz(wynik) {
      if (!wynik.udany || wynik.wynik === undefined) {
        miejsce.replaceChildren(
          utworzAkapit('ao-odmowa', opisOdmowyAod('Odczyt kontekstu okna', 'kontekst', wynik.blad)),
        );
        return;
      }

      miejsce.replaceChildren(zbudujWykaz(wynik.wynik.context));
    },
  };
}

/** Rozpoznaje komplet, w którym rdzeń nie wypełnił ani jednego pola. */
function kompletPusty(komplet: ContextBundle): boolean {
  return (
    (komplet.prompt ?? '') === '' &&
    (komplet.documentIds ?? []).length === 0 &&
    (komplet.projectId ?? '') === '' &&
    (komplet.agentIds ?? []).length === 0 &&
    (komplet.historyMessageIds ?? []).length === 0 &&
    (komplet.knowledgeSourceIds ?? []).length === 0 &&
    komplet.executionParams === undefined
  );
}

function zbudujWykaz(komplet: ContextBundle): HTMLElement {
  const opakowanie = document.createElement('div');
  opakowanie.className = 'ao-kontekst';

  if (kompletPusty(komplet)) {
    opakowanie.append(
      utworzAkapit(
        'ao-pusto',
        'Komplet kontekstu jest pusty — rdzeń nie ma dla tego okna ani polecenia, ani ' +
          'dokumentów, ani historii. Pusty komplet jest stanem poprawnym, nie odmową.',
      ),
    );
    return opakowanie;
  }

  const lista = utworzWykazPol();
  dodajPole(lista, 'Polecenie wyjściowe', komplet.prompt ?? '');
  dodajPole(lista, 'Projekt', komplet.projectId ?? '');
  dodajPoleWykazu(lista, 'Dokumenty', komplet.documentIds);
  dodajPoleWykazu(lista, 'Agenci', komplet.agentIds);
  dodajPoleWykazu(lista, 'Historia rozmowy', komplet.historyMessageIds);
  dodajPoleWykazu(lista, 'Źródła wiedzy', komplet.knowledgeSourceIds);
  dodajPole(
    lista,
    'Parametry wykonania',
    komplet.executionParams === undefined ? BRAK : 'rdzeń przysłał parametry wykonania',
  );
  opakowanie.append(lista);

  return opakowanie;
}
