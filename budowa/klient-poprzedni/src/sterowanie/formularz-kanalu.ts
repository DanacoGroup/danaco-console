import { KnownChannelKinds, type Channel } from '../../../shared/contract';
import {
  pole,
  poleTresci,
  przelacznik,
  wiersz,
} from '../modele/kontrolki-formularza-braki';

// Formularz wiersza rejestru kanałów: jedna powierzchnia obsługująca założenie i zmianę.

/** Odczyt pól formularza wraz z rozstrzygnięciem ich poprawności: niepusty błąd znaczy, że żądania do rdzenia nie wysyłamy. */
export interface OdczytKanalu {
  /** Przeszkoda po stronie klienta; pusty napis znaczy „pola w porządku". */
  blad: string;
  nazwa: string;
  rodzaj: string;
  /** Identyfikator modelu; pusty napis znaczy „model domyślny programu". */
  model: string;
  czynny: boolean;
  /** Parametry z pola JSON; `undefined` znaczy „nie zmieniaj parametrów". */
  parametry?: unknown;
}

export interface FormularzKanalu {
  /** Wiersze pól osadzane w ciele okna. */
  element: HTMLElement;
  /** Wpisuje w pola wskazany wiersz rejestru; `undefined` czyści formularz. */
  wypelnij(kanal: Channel | undefined): void;
  /** Odczyt pól wraz z rozstrzygnięciem czytelności parametrów. */
  odczytaj(): OdczytKanalu;
  /** Zdanie o niezmienialnym rodzaju — widoczne, gdy wskazany jest wiersz. */
  ostrzezRodzaj(widoczne: boolean): void;
}

/** Podpowiedź rodzaju kanału: rodzaje znane w chwili wydania kontraktu, pokazywane pod polem jako informacja, nie jako brama wpisu. */
const PODPOWIEDZ_RODZAJU = KnownChannelKinds.join(' · ');

const ZDANIE_RODZAJU =
  'Rodzaj kanału ustala się przy zakładaniu. Żądanie channel.update nie ma pola rodzaju, ' +
  'więc zmiana wpisana tutaj nie dojdzie do rdzenia dla kanału już istniejącego.';

export function utworzFormularzKanalu(): FormularzKanalu {
  const nazwa = pole('Nazwa kanału', 'np. Kanał lokalny (claude)');
  const rodzaj = pole('Rodzaj kanału', PODPOWIEDZ_RODZAJU);
  const model = pole('Identyfikator modelu', 'pusty = model domyślny programu');
  const czynny = przelacznik('Kanał czynny');
  const parametry = poleTresci('Parametry kanału (JSON)', 4, '{ "provider": "anthropic" }');

  const uwagaRodzaju = document.createElement('p');
  uwagaRodzaju.className = 'dn-pole-opis dc-kanaly__uwaga';
  uwagaRodzaju.textContent = ZDANIE_RODZAJU;
  uwagaRodzaju.hidden = true;

  const element = document.createElement('div');
  element.className = 'dc-kanaly__formularz';
  element.append(...zlozWiersze({ nazwa, rodzaj, uwagaRodzaju, model, czynny, parametry }));

  return {
    element,

    wypelnij(kanal) {
      nazwa.value = kanal?.name ?? '';
      rodzaj.value = kanal?.kind ?? '';
      model.value = kanal?.model ?? '';
      czynny.checked = kanal?.enabled ?? true;
      parametry.value = zapisParametrow(kanal?.config);
    },

    odczytaj: () => odczyt(nazwa, rodzaj, model, czynny, parametry),

    ostrzezRodzaj(widoczne) {
      uwagaRodzaju.hidden = !widoczne;
    },
  };
}

/** Buduje wiersze formularza jako czysty fragment konstrukcyjny, bez domknięcia nad stanem panelu ani nasłuchu zdarzeń. */
function zlozWiersze(kontrolki: {
  nazwa: HTMLInputElement;
  rodzaj: HTMLInputElement;
  uwagaRodzaju: HTMLElement;
  model: HTMLInputElement;
  czynny: HTMLInputElement;
  parametry: HTMLTextAreaElement;
}): HTMLElement[] {
  return [
    wiersz('Nazwa kanału', kontrolki.nazwa, { klasa: 'dc-kanaly__wiersz' }),
    wiersz('Rodzaj kanału', kontrolki.rodzaj, {
      klasa: 'dc-kanaly__wiersz',
      objasnienie: `Wartość danych, nie typ kodu. Rodzaje znane: ${PODPOWIEDZ_RODZAJU}.`,
    }),
    kontrolki.uwagaRodzaju,
    wiersz('Identyfikator modelu', kontrolki.model, {
      klasa: 'dc-kanaly__wiersz',
      objasnienie: 'Pusty wpis znaczy brak wskazania, nie brak działania.',
    }),
    wiersz('Kanał czynny', kontrolki.czynny, {
      klasa: 'dc-kanaly__wiersz dc-kanaly__wiersz--logiczny',
      objasnienie: 'Kanał nieczynny zostaje w rejestrze i pozostaje wybieralny w oknie.',
    }),
    wiersz('Parametry kanału (JSON)', kontrolki.parametry, {
      klasa: 'dc-kanaly__wiersz',
      objasnienie:
        'Parametry niosą wyłącznie odwołania do danych dostępowych, nigdy ich treść. ' +
        'Pole puste znaczy „nie zmieniaj parametrów".',
    }),
  ];
}

/** Odczytuje pola formularza wraz z rozbiorem parametrów z tekstu pola na wartość struktury żądania do rdzenia. */
function odczyt(
  nazwa: HTMLInputElement,
  rodzaj: HTMLInputElement,
  model: HTMLInputElement,
  czynny: HTMLInputElement,
  parametry: HTMLTextAreaElement,
): OdczytKanalu {
  const wspolne = {
    nazwa: nazwa.value.trim(),
    rodzaj: rodzaj.value.trim(),
    model: model.value.trim(),
    czynny: czynny.checked,
  };
  const surowe = parametry.value.trim();
  if (surowe === '') return { blad: '', ...wspolne };
  try {
    return { blad: '', ...wspolne, parametry: JSON.parse(surowe) as unknown };
  } catch {
    return {
      blad: 'Parametry kanału nie są czytelnym JSON-em — żądania nie wysłano, rdzeń o niczym nie wie.',
      ...wspolne,
    };
  }
}

/** Zapisuje parametry wiersza jako tekst pola formularza do edycji; brak parametrów w wierszu daje pole puste. */
function zapisParametrow(config: unknown): string {
  if (config === undefined || config === null) return '';
  if (typeof config === 'string') return config;
  try {
    return JSON.stringify(config, null, 2);
  } catch {
    return '';
  }
}
