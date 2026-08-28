import type { DesignAsset } from '../../../../shared/contract';
import { nazwaRodzaju } from './karta-zasobu';
import { przytnijKod } from './przyciecie-pol';
import {
  czyRdzenZnaObraz,
  slowoModelu,
  zdanieOBrakuObrazu,
  zdanieOSciezceRdzenia,
  zdanieOSlowieModelu,
} from './slowo-modelu';

/**
 * Płyta podglądu Preview Window — powierzchnia zasobu wraz z tym, co rdzeń o nim wie, w trzech
 * stanach: obraz wstawiony, obraz niewczytany nazwany, obraz nieoddany przez rdzeń.
 */
export interface PlytaPodgladu {
  element: HTMLElement;
  /** Pokazuje zasób; `null` znaczy „żaden nie jest wskazany". */
  pokaz(zasob: DesignAsset | null): void;
}

export function utworzPlytePodgladu(): PlytaPodgladu {
  const powierzchnia = document.createElement('div');
  powierzchnia.className = 'md-podglad__plyta';

  const zdanieTresci = document.createElement('p');
  zdanieTresci.className = 'dn-pole-opis md-podglad__zrodlo';

  const slowo = document.createElement('p');
  slowo.className = 'md-podglad__slowo';

  const oSlowie = document.createElement('p');
  oSlowie.className = 'dn-pole-opis md-podglad__o-slowie';

  const wykaz = document.createElement('dl');
  wykaz.className = 'md-podglad__wykaz';

  const element = document.createElement('div');
  element.className = 'md-podglad';
  element.append(powierzchnia, zdanieTresci, slowo, oSlowie, wykaz);

  return {
    element,

    pokaz(zasob) {
      if (zasob === null) {
        powierzchnia.dataset['tresc'] = 'bez-zasobu';
        powierzchnia.replaceChildren(napis('Żaden zasób nie jest wskazany'));
        zdanieTresci.textContent =
          'Podgląd dotyczy zasobu wskazanego w Assets Panel — wskaż go tam, a stanie tutaj.';
        slowo.textContent = '';
        slowo.hidden = true;
        oSlowie.textContent = '';
        wykaz.replaceChildren();
        return;
      }

      if (czyRdzenZnaObraz(zasob)) wstawObraz(powierzchnia, zdanieTresci, zasob);
      else {
        powierzchnia.dataset['tresc'] = 'bez-obrazu';
        powierzchnia.replaceChildren(napis(`${nazwaRodzaju(zasob.kind)} — rdzeń nie zna treści obrazu`));
        zdanieTresci.textContent = zdanieOBrakuObrazu();
      }

      const tresc = slowoModelu(zasob);
      slowo.hidden = tresc === '';
      slowo.textContent = tresc;
      oSlowie.textContent = zdanieOSlowieModelu(zasob);
      wykaz.replaceChildren(...wierszeWykazu(zasob));
    },
  };
}

/** Wstawia obraz z jego adresu wraz z podpisem tekstowym; nieudane wczytanie samo nazywa swoją przyczynę. */
function wstawObraz(powierzchnia: HTMLElement, zdanie: HTMLElement, zasob: DesignAsset): void {
  const adres = (zasob.uri ?? '').trim();
  const obraz = document.createElement('img');
  obraz.className = 'md-podglad__obraz';
  obraz.src = adres;
  obraz.alt = slowoModelu(zasob) === '' ? 'Zasób wizualny bez opisu słownego' : slowoModelu(zasob);
  obraz.addEventListener('error', () => {
    powierzchnia.dataset['tresc'] = 'obraz-niewczytany';
    // Napis mówi, gdzie treść jest, a nie że jej nie ma: bajty leżą w rdzeniu, brakuje drogi po nie.
    powierzchnia.replaceChildren(
      napis('Treść obrazu leży w rdzeniu — droga po nią czeka na dobudowę'),
    );
    zdanie.textContent = zdanieOSciezceRdzenia(adres);
  });
  powierzchnia.dataset['tresc'] = 'obraz';
  powierzchnia.replaceChildren(obraz);
  zdanie.textContent = `Treść obrazu wzięta z pola uri zasobu: ${adres}`;
}

/** Zdanie zastępujące wartość, której rdzeń w odpowiedzi nie podał; dla promptu mówi więcej niż „nie podał". */
const BRAK_WARTOSCI: Readonly<Record<string, string>> = {
  'Prompt źródłowy':
    'rdzeń nie podaje promptId w żadnej odpowiedzi — nie ma przekładu klucza wiersza ' +
    'promptu na kod kontraktu i nie zgaduje go (zgłoszone)',
};

/** Wiersz wykazu pól — wartość albo zdanie o jej braku, nigdy pustka bez wyjaśnienia jej powodu Operatorowi. */
function wierszeWykazu(zasob: DesignAsset): readonly HTMLElement[] {
  const wymiary =
    zasob.width !== undefined && zasob.height !== undefined ? `${zasob.width}×${zasob.height}` : '';
  const etykiety = (zasob.tags ?? []).join(', ');
  const pary: readonly (readonly [string, string])[] = [
    ['Identyfikator', zasob.id],
    ['Rodzaj', nazwaRodzaju(zasob.kind)],
    ['Format', zasob.format ?? ''],
    ['Wymiary', wymiary],
    ['Treść obrazu (uri)', zasob.uri ?? ''],
    ['Etykiety', etykiety],
    ['Ulubiony', zasob.favorite === true ? 'tak' : 'nie'],
    ['Wariant zasobu', zasob.variantOfAssetId ?? ''],
    ['Prompt źródłowy', zasob.promptId ?? ''],
    ['Okno modułu', zasob.windowId],
    ['Powstał', new Date(zasob.createdAt).toLocaleString('pl')],
  ];
  return pary.flatMap(([nazwa, wartosc]) => wiersz(nazwa, wartosc));
}

function wiersz(nazwa: string, wartosc: string): readonly HTMLElement[] {
  const klucz = document.createElement('dt');
  klucz.className = 'md-podglad__klucz';
  klucz.textContent = nazwa;

  // Wartość złożona z samych znaków białych jest brakiem tak samo jak pustka, nie polem wypełnionym.
  const podane = przytnijKod(wartosc) !== '';
  const wartoscElement = document.createElement('dd');
  wartoscElement.className = 'md-podglad__wartosc';
  wartoscElement.dataset['podane'] = String(podane);
  wartoscElement.textContent = podane
    ? wartosc
    : (BRAK_WARTOSCI[nazwa] ?? 'rdzeń tego nie podał');
  return [klucz, wartoscElement];
}

function napis(tresc: string): HTMLElement {
  const element = document.createElement('span');
  element.className = 'md-podglad__napis';
  element.textContent = tresc;
  return element;
}
