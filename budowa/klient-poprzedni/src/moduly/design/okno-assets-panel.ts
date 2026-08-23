import { Command, type DesignAsset } from '../../../../shared/contract';
import { przycisk } from '../../modele/kontrolki-formularza';
import { OKNO_ASSETS_PANEL } from './etykiety-designu';
import { utworzFiltrZasobow, type FiltrZasobow } from './filtr-zasobow';
import { utworzCzynnosciZasobu, type CzynnosciZasobu } from './czynnosci-zasobu';
import { czyPasujeDoFrazy, utworzKarteZasobu } from './karta-zasobu';
import { utworzNadanieEtykiet, type NadanieEtykiet } from './nadanie-etykiet';
import { utworzWgranieZasobu, type WgranieZasobu } from './wgranie-zasobu';
import { utworzPanelMetadanych, type PanelMetadanych } from './panel-metadanych';
import { utworzWydaniaZasobu, type WydaniaZasobu } from './wydania-zasobu';
import { utworzKolekcjeDesignu, type KolekcjeDesignu } from './kolekcje-designu';
import { naglowekOkna, utworzStanOkna, type StanOkna } from './stan-okna';
import { powodZTorem } from './tor-komendy';
import type { StanDesignu } from './stan-designu';

/**
 * Assets Panel — okno **zarządca** modułu Design (kod katalogu rdzenia
 * `assets-panel`): wykaz zasobów wraz z filtrem i panelem metadanych.
 *
 * Zasób wchodzi do wykazu dwiema drogami. `design.asset.upload` przyjmuje plik
 * wskazany w oknie — klient czyta jego bajty i oddaje je rdzeniowi
 * (`wgranie-zasobu.ts`); kontrolka stoi w oknie pierwsza. `design.asset.generate`
 * oddaje bajty z kanału obrazowego, a powstały zasób trafia do wykazu zdarzeniem
 * `design.asset.changed`, bez czynności w tym oknie.
 *
 * Zawartość wykazu przestawiają jeszcze trzy komendy: `design.asset.tag.set`
 * zasila filtr etykiet, `design.asset.favorite.set` — przełącznik „Tylko
 * ulubione", a `design.asset.remove` jest nadawcą rodzaju zmiany `deleted`;
 * obie ostatnie stoją w `czynnosci-zasobu.ts`. Usunięcie idzie bez pytania
 * „czy na pewno" — uzasadnienie w nagłówku tamtego pliku.
 *
 * Odmowa odczytu nie gasi okna: `design.asset.list` zakończona odmową wchodzi
 * w stan błędu, a okno zostaje czynne. Zasób przysłany zdarzeniem
 * `design.asset.changed` wejdzie do wykazu mimo to, bo wciąga go stan modułu,
 * nie ta odpowiedź.
 */
export interface OknoAssetsPanel {
  element: HTMLElement;
  odswiez(): void;
  /** Zleca odczyt zasobów z warunkami filtra. */
  wczytaj(): Promise<void>;
}

export function utworzOknoAssetsPanel(
  stan: StanDesignu,
  naKanwe: (zasob: DesignAsset) => void,
): OknoAssetsPanel {
  const okno: StanOkna = utworzStanOkna();
  const metadane: PanelMetadanych = utworzPanelMetadanych(stan);
  const etykietowanie: NadanieEtykiet = utworzNadanieEtykiet(stan);
  const wgranie: WgranieZasobu = utworzWgranieZasobu(stan);
  const czynnosci: CzynnosciZasobu = utworzCzynnosciZasobu(stan);
  // Wydania i kolekcje dotyczą zasobu wskazanego w wykazie, więc stoją pod nim,
  // obok pozostałych czynności na jednym zasobie.
  const wydania: WydaniaZasobu = utworzWydaniaZasobu(stan);
  const kolekcje: KolekcjeDesignu = utworzKolekcjeDesignu(stan);

  const wykaz = document.createElement('div');
  wykaz.className = 'md-zasoby';
  wykaz.dataset['widok'] = 'siatka';

  const filtr: FiltrZasobow = utworzFiltrZasobow({
    naOdczyt: () => void wczytaj(),
    naZawezenie: () => odswiez(),
    naWidok: () => odswiez(),
  });

  const naPlansze = przycisk('Przeciągnij na Design Board', 'dn-btn dn-btn--sm dn-btn--zarys');
  naPlansze.addEventListener('click', () => {
    const zasob = stan.wybrany();
    if (zasob === null) return;
    naKanwe(zasob);
  });

  // Wgranie stoi przed filtrem, bo filtr zawęża zbiór, który wgranie dopiero
  // tworzy. Czynności na zasobie wskazanym (ulubiony, usunięcie) idą zaraz pod
  // wykazem, obok nadania etykiet — wszystkie trzy dotyczą jednego zasobu
  // wskazanego kartą.
  okno.tresc.append(
    wgranie.element,
    filtr.element,
    naPlansze,
    wykaz,
    czynnosci.element,
    etykietowanie.element,
    wydania.element,
    kolekcje.element,
    metadane.element,
  );

  const element = document.createElement('section');
  element.className = 'md-okno md-okno--zarzadca';
  element.dataset['okno'] = OKNO_ASSETS_PANEL.kod;
  element.append(naglowekOkna(OKNO_ASSETS_PANEL.nazwa, OKNO_ASSETS_PANEL.rola), okno.element);

  /** Zasoby po zawężeniu miejscowym — po nazwie i po etykietach. */
  function widoczne(): readonly DesignAsset[] {
    const fraza = filtr.fraza();
    if (fraza === '') return stan.zasoby();
    return stan.zasoby().filter((zasob) => czyPasujeDoFrazy(zasob, fraza));
  }

  function odswiez(): void {
    metadane.odswiez();
    etykietowanie.odswiez();
    czynnosci.odswiez();
    wydania.odswiez();
    kolekcje.odswiez();
    wykaz.dataset['widok'] = filtr.widok();
    const zasoby = widoczne();
    const wybrany = stan.wybrany();
    wykaz.replaceChildren(
      ...zasoby.map((zasob) =>
        utworzKarteZasobu(zasob, zasob.id === wybrany?.id, (kod) => stan.wybierz(kod)),
      ),
    );
    pokazStan(zasoby.length);
  }

  function pokazStan(widocznych: number): void {
    if (stan.faza() === 'odczyt') {
      // Zdanie czuwania („rdzeń odpowiada, odczyt trwa N s") wchodzi na miejsce
      // zapowiedzi, gdy jest — pokazuje, że milczenie rdzenia zostało sprawdzone.
      okno.ladowanie(stan.powod() === '' ? 'Odczyt zasobów w toku…' : stan.powod());
      return;
    }
    if (stan.faza() === 'blad') {
      // Zdanie o torze komendy tłumaczy pola żądania i należy się odmowie
      // rdzenia. Zerwanego gniazda nie tłumaczy żadne pole żądania, więc przy
      // braku rozstrzygnięcia dopisek zostaje pominięty.
      okno.blad(
        stan.czyBezRozstrzygniecia()
          ? stan.powod()
          : powodZTorem(stan.powod(), Command.DesignAssetList),
      );
      return;
    }
    if (stan.faza() === 'spoczynek') {
      okno.puste('Zasoby nie były jeszcze odczytywane — naciśnij „Odczytaj zasoby".');
      return;
    }
    if (widocznych === 0) {
      okno.puste(
        stan.zasoby().length === 0
          ? 'Rdzeń nie zna ani jednego zasobu tego okna.'
          : 'Żaden wczytany zasób nie pasuje do frazy zawężającej.',
      );
      return;
    }
    okno.gotowe();
  }

  async function wczytaj(): Promise<void> {
    await stan.odswiez(filtr.warunki());
    // Kolekcje idą tą samą drogą co zasoby: wykaz kolekcji nieodświeżony po
    // odczycie pokazywałby liczniki sprzed zmian, których Operator właśnie
    // dokonał w innym oknie.
    await kolekcje.wczytaj();
  }

  return { element, odswiez, wczytaj };
}
