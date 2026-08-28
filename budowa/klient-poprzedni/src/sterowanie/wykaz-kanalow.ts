import type { Channel } from '../../../shared/contract';
import { oznaczFaze, type FazaOkna } from '../komponenty/faza-okna';
import { przyciskAkcji, wykaz } from '../modele/kontrolki-formularza-braki';
import { nazwaKanalu, type RejestrKanalow } from './rejestr-kanalow';

/** Wykaz wierszy rejestru kanałów wraz z trzema stanami obowiązkowymi: ładowanie, pustka i odmowa rdzenia. */
export interface WykazKanalow {
  /** Element osadzany w ciele okna. */
  element: HTMLElement;
  /** Odrysowuje wykaz ze stanu rejestru. */
  odrysuj(): void;
  /** Kanał wskazany do zmiany albo usunięcia; `undefined` = brak wskazania. */
  wskazany(): Channel | undefined;
  /** Ustawia wskazanie po identyfikatorze; nieznany identyfikator je zdejmuje. */
  wskaz(idKanalu: string): void;
  /** Przesłania wykaz zdaniem odmowy rdzenia — faza `blad`. */
  pokazOdmowe(zdanie: string): void;
}

/** Zdania trzech stanów pasa: jedno miejsce dla tytułu i opisu, bez powtórzeń w poszczególnych gałęziach. */
const ZDANIA: Record<'ladowanie' | 'pusto', { tytul: string; opis: string }> = {
  ladowanie: {
    tytul: 'Odczytuję rejestr kanałów',
    opis: 'Pytanie channel.list poszło do rdzenia; odpowiedzi jeszcze nie ma.',
  },
  pusto: {
    tytul: 'Rejestr kanałów jest pusty',
    opis:
      'Rdzeń odpowiedział, ale nie zna ani jednego kanału. Wypełnij formularz i założ pierwszy ' +
      'wiersz — bez kanału okno komunikacji nie ma czym rozmawiać.',
  },
};

export function utworzWykazKanalow(
  rejestr: RejestrKanalow,
  przyWskazaniu: (kanal: Channel | undefined) => void,
): WykazKanalow {
  let wskazanie = '';
  const { element, pas, lista } = zlozPowloke();

  function kanaly(): Channel[] {
    return rejestr.kanaly();
  }

  function ustawStan(faza: FazaOkna, tytul: string, opis: string): void {
    naniesStan(element, pas, faza, tytul, opis);
  }

  function odrysuj(): void {
    const wykazKanalow = kanaly();
    if (wykazKanalow.every((kanal) => kanal.id !== wskazanie)) wskazanie = '';
    lista.replaceChildren(
      ...wykazKanalow.map((kanal) =>
        pozycja(kanal, kanal.id === wskazanie, () => {
          wskazanie = kanal.id;
          odrysuj();
          przyWskazaniu(kanal);
        }),
      ),
    );
    if (!rejestr.odpowiedzOtrzymana()) {
      ustawStan('ladowanie', ZDANIA.ladowanie.tytul, ZDANIA.ladowanie.opis);
      return;
    }
    if (wykazKanalow.length === 0) {
      ustawStan('puste', ZDANIA.pusto.tytul, ZDANIA.pusto.opis);
      return;
    }
    ustawStan('gotowe', '', '');
  }

  return {
    element,
    odrysuj,
    wskazany: () => kanaly().find((kanal) => kanal.id === wskazanie),
    wskaz(idKanalu) {
      wskazanie = idKanalu;
      odrysuj();
    },
    pokazOdmowe(tresc) {
      ustawStan('blad', 'Rdzeń odmówił', tresc);
    },
  };
}

/** Powłoka wykazu: pas stanu nad listą wierszy, jako czysty fragment konstrukcyjny bez domknięcia na stanie. */
function zlozPowloke(): { element: HTMLElement; pas: HTMLElement; lista: HTMLUListElement } {
  const lista = wykaz('Wiersze rejestru kanałów', 'dc-kanaly__lista');

  const pas = document.createElement('div');
  pas.className = 'dn-pusty-stan dc-kanaly__stan';

  const element = document.createElement('div');
  element.className = 'dc-kanaly__wykaz';
  element.append(pas, lista);
  return { element, pas, lista };
}

/** Nanosi na wykaz fazę wraz ze zdaniami pasa stanu — jedno miejsce wspólne dla wszystkich trzech stanów. */
function naniesStan(
  element: HTMLElement,
  pas: HTMLElement,
  faza: FazaOkna,
  tytul: string,
  opis: string,
): void {
  pas.replaceChildren(zdanie('dn-pusty-stan-tytul', tytul), zdanie('dn-pusty-stan-opis', opis));
  oznaczFaze(element, pas, faza);
}

/** Pozycja wykazu: nazwa wiersza z rejestru, jego identyfikator i jeden przycisk wskazania, bez przycisku usuwania. */
function pozycja(kanal: Channel, wskazany: boolean, przyWskazaniu: () => void): HTMLElement {
  const tytul = document.createElement('strong');
  tytul.className = 'dc-kanaly-pozycja__tytul';
  tytul.textContent = nazwaKanalu(kanal);

  const opis = document.createElement('span');
  opis.className = 'dc-kanaly-pozycja__opis';
  opis.textContent = kanal.id;

  const przycisk = przyciskAkcji(
    wskazany ? 'Wskazany' : 'Wskaż',
    wskazany ? 'dn-btn dn-btn--sm dn-btn--wybrany' : 'dn-btn dn-btn--sm dn-btn--zarys',
  );
  przycisk.setAttribute('aria-pressed', String(wskazany));
  przycisk.addEventListener('click', przyWskazaniu);

  const element = document.createElement('li');
  element.className = 'dc-kanaly-pozycja';
  element.dataset['kanal'] = kanal.id;
  element.dataset['wskazany'] = String(wskazany);
  element.append(tytul, opis, przycisk);
  return element;
}

/** Akapit pasa stanu, budowany jednakowo dla tytułu i dla opisu wykazu, różniący się jedynie klasą stylu. */
function zdanie(klasa: string, tresc: string): HTMLElement {
  const element = document.createElement('p');
  element.className = klasa;
  element.textContent = tresc;
  return element;
}
