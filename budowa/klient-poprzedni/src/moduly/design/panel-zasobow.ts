import './design.css';
import './zasoby.css';

import { ChangeKind, Command, type DesignAsset } from '../../../../shared/contract';
import { pole, przyciskAkcji } from '../../modele/kontrolki-formularza-braki';
import type { OpcjePanelu, PanelPomocniczy } from '../../okna-pomocnicze/panel-pomocniczy';
import { czyPasujeDoFrazy, utworzKarteZasobu } from './karta-zasobu';
import { utworzStanOkna, type StanOkna } from './stan-okna';
import { powodZTorem } from './tor-komendy';
import {
  przyjmijOdczytZasobow,
  pustyZapisDesignu,
  usunZasob,
  wchlonZasob,
  type ZapisDesignu,
} from './zapis-designu';
import { utworzZrodloDesignu, type ZrodloDesignu } from './zrodlo-designu';

/**
 * Zasoby Designu — panel do stosu paneli pomocniczych obcego modułu, osobny gęstszy widok tych
 * samych danych, który pożycza komplet reguł od modułu Design.
 */
export const KOD_PANELU_ZASOBOW = 'zasoby-designu';

export function utworzPanelZasobow(opcje: OpcjePanelu): PanelPomocniczy {
  const zrodlo: ZrodloDesignu = utworzZrodloDesignu(opcje.kanal);
  const zapis: ZapisDesignu = pustyZapisDesignu();
  const okno: StanOkna = utworzStanOkna();

  // Widok listowy, nie siatkowy: jedna kolumna, bo siatka w wąskim pasie dałaby karty urwane w połowie.
  const wykaz = document.createElement('div');
  wykaz.className = 'md-zasoby';
  wykaz.dataset['widok'] = 'lista';
  okno.tresc.append(wykaz);

  const tytul = document.createElement('h3');
  tytul.className = 'md-panel__tytul';
  tytul.textContent = 'Zasoby Designu';

  const licznik = document.createElement('span');
  licznik.className = 'dn-plakietka';
  licznik.classList.add('md-panel__licznik');

  const szukaj = pole('Zawęź zasoby Designu po nazwie lub etykiecie', 'Zawęź…');
  szukaj.classList.add('md-panel__szukaj');
  szukaj.addEventListener('input', () => rysuj());

  // Przycisk jest czynny zawsze; odmowa rdzenia ląduje w stanie panelu, nie
  // w wyszarzeniu kontrolki.
  const odswiezPrzycisk = przyciskAkcji('Odśwież', 'dn-btn dn-btn--sm dn-btn--zarys');
  odswiezPrzycisk.addEventListener('click', () => odswiez());

  const naglowek = document.createElement('header');
  naglowek.className = 'md-panel__naglowek';
  naglowek.append(tytul, licznik, szukaj, odswiezPrzycisk);

  const element = document.createElement('section');
  element.className = 'md-panel';
  element.dataset['panel'] = KOD_PANELU_ZASOBOW;
  // Przedrostek i okno gospodarza idą atrybutami danych, bo wygląd panelu należy do arkusza Designu.
  element.dataset['gospodarz'] = opcje.przedrostek;
  element.dataset['oknoGospodarza'] = opcje.okno;
  element.setAttribute('aria-label', `Zasoby modułu Design w module ${opcje.modul}`);
  element.append(naglowek, okno.element);

  /** Zasoby po zawężeniu miejscowym — fraza nie jest polem żądania. */
  function widoczne(): readonly DesignAsset[] {
    const fraza = szukaj.value.trim().toLowerCase();
    if (fraza === '') return zapis.zbior;
    return zapis.zbior.filter((zasob) => czyPasujeDoFrazy(zasob, fraza));
  }

  /** Wybór jest przestawny: powtórne naciśnięcie karty go zdejmuje. */
  function wybierz(idZasobu: string): void {
    zapis.wybor = zapis.wybor === idZasobu ? null : idZasobu;
    rysuj();
  }

  function rysuj(): void {
    const zasoby = widoczne();
    licznik.textContent =
      zasoby.length === zapis.zbior.length
        ? String(zapis.zbior.length)
        : `${zasoby.length}/${zapis.zbior.length}`;
    wykaz.replaceChildren(
      ...zasoby.map((zasob) => utworzKarteZasobu(zasob, zasob.id === zapis.wybor, wybierz)),
    );
    pokazStan(zasoby.length);
  }

  function pokazStan(widocznych: number): void {
    if (zapis.faza === 'odczyt') {
      okno.ladowanie('Odczyt zasobów Designu w toku…');
      return;
    }
    if (zapis.faza === 'blad') {
      okno.blad(powodZTorem(zapis.powod, Command.DesignAssetList));
      return;
    }
    if (zapis.faza === 'spoczynek') {
      okno.puste('Zasoby Designu nie były jeszcze odczytywane — naciśnij „Odśwież".');
      return;
    }
    if (widocznych === 0) {
      okno.puste(
        zapis.zbior.length === 0
          ? 'Rdzeń nie zna ani jednego zasobu Designu. Panel czyta zasoby WSZYSTKICH okien ' +
            'Designu, nie okna tego modułu — pustka znaczy więc pusty zbiór, a nie złe okno.'
          : 'Żaden wczytany zasób nie pasuje do frazy zawężającej.',
      );
      return;
    }
    okno.gotowe();
  }

  /** Ponowny odczyt z rdzenia — pole okna zostaje puste, więc rdzeń nie zawęża wykazu do żadnego okna. */
  function odswiez(): void {
    zapis.faza = 'odczyt';
    zapis.powod = '';
    zapis.bezRozstrzygniecia = false;
    rysuj();
    void zrodlo.zasoby(zapis.warunki).then((wynik) => {
      przyjmijOdczytZasobow(zapis, wynik);
      rysuj();
    });
  }

  // Zdarzenie jest drugą drogą odświeżenia: zasób powstały gdziekolwiek wchodzi tu bez pytania.
  const odsubskrybuj = zrodlo.naZmianeZasobu((tresc) => {
    if (tresc.change === ChangeKind.Deleted) usunZasob(zapis, tresc.asset.id);
    else wchlonZasob(zapis, tresc.asset);
    rysuj();
  });

  // Pierwszego odczytu panel nie robi sam — pas paneli woła odświeżenie na wszystkich panelach.
  rysuj();

  return { element, odswiez, zamknij: odsubskrybuj };
}
