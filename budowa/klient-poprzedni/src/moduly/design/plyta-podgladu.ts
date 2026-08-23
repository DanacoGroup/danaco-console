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
 * Płyta podglądu Preview Window — powierzchnia zasobu wraz z tym, co rdzeń
 * o nim wie.
 *
 * Pokazuje jeden zasób takim, jakim rdzeń go zna. Powierzchnia ma trzy stany:
 *   (1) rdzeń podał `uri` — płyta wstawia obraz i mówi, skąd go bierze;
 *   (2) rdzeń podał `uri`, a treść spod niego się nie wczytała — płyta nazywa
 *       przyczynę, zamiast zrzucać ją na przeglądarkę;
 *   (3) rdzeń `uri` nie podał — płyta mówi, czego nie ma, i nie dorabia obrazu.
 *
 * Stan (2) jest stanem każdego zasobu z generowania: rdzeń wypełnia `uri`
 * ścieżką w swoim systemie plików (magazyn oddaje `filepath`), a nie adresem do
 * pobrania. Obraz powstaje poprawnie i mimo to się nie wyświetli, dopóki droga
 * po treść zasobu — opisana już w kontrakcie i wspólna całemu magazynowi — nie
 * dostanie uchwytu w rdzeniu; i tak brzmi zdanie stanu (2)
 * (`slowo-modelu.ts`, `zdanieOSciezceRdzenia`).
 *
 * Wykaz pól mówi także o tym, czego nie ma: wartość pusta to nie pustka
 * w wierszu, tylko zdanie „rdzeń tego nie podał". Format nieznany i format
 * nieistniejący to dwa różne stany.
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

/** Wstawia obraz z `uri` wraz z jego adresem; nieudane wczytanie mówi o sobie. */
function wstawObraz(powierzchnia: HTMLElement, zdanie: HTMLElement, zasob: DesignAsset): void {
  const adres = (zasob.uri ?? '').trim();
  const obraz = document.createElement('img');
  obraz.className = 'md-podglad__obraz';
  obraz.src = adres;
  obraz.alt = slowoModelu(zasob) === '' ? 'Zasób wizualny bez opisu słownego' : slowoModelu(zasob);
  obraz.addEventListener('error', () => {
    powierzchnia.dataset['tresc'] = 'obraz-niewczytany';
    // Napis mówi, gdzie treść jest, a nie że jej nie ma: bajty leżą w magazynie
    // rdzenia i zasób jest w porządku — brakuje wyłącznie drogi po nie.
    powierzchnia.replaceChildren(
      napis('Treść obrazu leży w rdzeniu — droga po nią czeka na dobudowę'),
    );
    zdanie.textContent = zdanieOSciezceRdzenia(adres);
  });
  powierzchnia.dataset['tresc'] = 'obraz';
  powierzchnia.replaceChildren(obraz);
  zdanie.textContent = `Treść obrazu wzięta z pola uri zasobu: ${adres}`;
}

/**
 * Zdanie zastępujące wartość, której rdzeń nie podał.
 *
 * Dla `promptId` zdanie mówi więcej niż „nie podał": klucza nie niesie żadna
 * odpowiedź rdzenia, bo nie ma on przekładu klucza wiersza promptu na kod
 * kontraktu i nie zgaduje go (`adapter_modul_design.go`, `zasobKontraktu`).
 * Samo „rdzeń tego nie podał" kazałoby sądzić, że zasób prompt zgubił.
 */
const BRAK_WARTOSCI: Readonly<Record<string, string>> = {
  'Prompt źródłowy':
    'rdzeń nie podaje promptId w żadnej odpowiedzi — nie ma przekładu klucza wiersza ' +
    'promptu na kod kontraktu i nie zgaduje go (zgłoszone)',
};

/** Wiersz wykazu pól — wartość albo zdanie o jej braku, nigdy pustka. */
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

  // Wartość złożona z samych znaków białych jest brakiem tak samo jak pustka —
  // inaczej wiersz pokazałby pole „wypełnione" i nie wyświetlił ani znaku.
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
