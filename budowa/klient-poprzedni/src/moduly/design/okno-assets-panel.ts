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
 * Assets Panel — okno zarządca modułu Design, wykaz zasobów wraz z filtrem i panelem metadanych,
 * zasilany dwiema drogami wejścia zasobu oraz trzema komendami przestawiającymi jego zawartość.
 */
export interface OknoAssetsPanel {
  element: HTMLElement;
  odswiez(): void;
  /** Zleca odczyt zasobów z warunkami filtra ustawionymi w oknie. */
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
  // Wydania i kolekcje dotyczą zasobu wskazanego w wykazie, więc stoją pod nim, obok innych czynności.
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

  // Wgranie stoi przed filtrem, bo filtr zawęża zbiór, który wgranie dopiero tworzy.
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

  /** Zasoby po zawężeniu miejscowym — po nazwie zasobu i po jego etykietach. */
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
      // Zdanie czuwania wchodzi na miejsce zapowiedzi, gdy jest — pokazuje, że milczenie rdzenia sprawdzono.
      okno.ladowanie(stan.powod() === '' ? 'Odczyt zasobów w toku…' : stan.powod());
      return;
    }
    if (stan.faza() === 'blad') {
      // Zdanie o torze komendy tłumaczy pola żądania i należy się odmowie rdzenia, nie zerwanemu gniazdu.
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
    // Kolekcje idą tą samą drogą co zasoby: nieodświeżony wykaz pokazywałby liczniki sprzed zmian.
    await kolekcje.wczytaj();
  }

  return { element, odswiez, wczytaj };
}
