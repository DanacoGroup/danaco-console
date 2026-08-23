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
 * Zasoby Designu — panel do stosu paneli pomocniczych obcego modułu.
 *
 * Ma kształt `PanelPomocniczy`, więc staje w pasie paneli dowolnego modułu
 * i w kolumnie paneli sceny okien równoległych. Poza nim zasoby Designu są
 * osiągalne wyłącznie wewnątrz złożenia `modul-design.ts`.
 *
 * Panel jest osobnym, gęstszym widokiem tych samych danych, a nie opakowaniem
 * `okno-assets-panel.ts`. Tamto okno nie ma `zamknij()` — subskrypcję
 * `design.asset.changed` trzyma `stan-designu.ts`, a zdejmuje ją `rozlacz()`
 * wołane przez `modul-design.ts` — i żąda całego `StanDesignu` wraz z komendami
 * `window.list`, `window.state.get`, `channel.list`, `module.list` oraz kanwy
 * Design Board, której poza modułem Design nie ma. Tu mieści się jeden wiersz
 * nagłówka, nie trzy pasy kontrolek okna operacyjnego.
 *
 * Własnych reguł o zasobie panel nie pisze — pożycza komplet od modułu:
 *   `zrodlo-designu.ts`  — jedyna warstwa wywołań `design.asset.list`
 *                          i subskrypcji `design.asset.changed`;
 *   `zapis-designu.ts`   — `pustyZapisDesignu`, `przyjmijOdczytZasobow`,
 *                          `wchlonZasob`, `usunZasob`, czyli całe wciąganie
 *                          zmian wraz z gałęzią `deleted`;
 *   `karta-zasobu.ts`    — karta i predykat frazy (jedna kopia na moduł);
 *   `stan-okna.ts`       — trzy stany obowiązkowe;
 *   `tor-komendy.ts`     — zdanie o torze doklejane do odmowy.
 * Oba widoki zbiegają się na tym samym zdarzeniu rdzenia.
 *
 * Panel czyta zasoby wszystkich okien Designu i nazywa to w stanie pustym.
 * `OpcjePanelu.okno` niesie okno gospodarza (Apps, Developer, Diagnostics), a nie
 * okno modułu Design; pole `windowId` żądania jest opcjonalne
 * (`shared/contract.ts`), a `warunkiFiltruZasobow` dokłada warunek `z.okno = ?`
 * tylko wtedy, gdy pole przyszło (`dane/design_zasoby.go`). Podstawienie okna
 * gospodarza zawęziłoby wykaz do zasobów obcego okna, czyli najczęściej do
 * pustki, więc panel pola nie podstawia.
 *
 * Panel nie oddaje zasobu gospodarzowi: komendy wstawiającej zasób Designu
 * w rozmowę obcego modułu kontrakt nie ma, a `context.transfer` biegnie
 * przeciwnie — z okna źródłowego do modułu docelowego, otwierając tam okno.
 * Nie generuje, nie nadaje etykiet i nie zapisuje kompozycji; te czynności
 * zostają w oknach modułu Design.
 */
export const KOD_PANELU_ZASOBOW = 'zasoby-designu';

export function utworzPanelZasobow(opcje: OpcjePanelu): PanelPomocniczy {
  const zrodlo: ZrodloDesignu = utworzZrodloDesignu(opcje.kanal);
  const zapis: ZapisDesignu = pustyZapisDesignu();
  const okno: StanOkna = utworzStanOkna();

  // Widok listowy, nie siatkowy: `md-zasoby[data-widok='lista']` to jedna
  // kolumna, a siatka po 180px w wąskim pasie dałaby karty urwane w połowie.
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
  // Przedrostek i okno gospodarza idą atrybutami danych, a nie do nazw klas:
  // klasy składanej w locie nie widzi kontrola pokrycia arkuszy, a wygląd
  // panelu należy do arkusza Designu, który wędruje z nim przez oba importy.
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

  /**
   * Ponowny odczyt z rdzenia — czynność `PanelPomocniczy.odswiez()`.
   *
   * `zapis.warunki.idOkna` zostaje pusty, więc `zrodlo-designu` pomija pole
   * `windowId` i rdzeń nie zawęża wykazu do żadnego okna (uzasadnienie
   * w nagłówku pliku).
   */
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

  // Zdarzenie jest drugą drogą odświeżenia: zasób powstały gdziekolwiek —
  // w oknach Designu, w rozmowie, po stronie rdzenia — wchodzi tu bez pytania.
  // Rozdział po rodzaju zmiany idzie z modułu wraz z gałęzią `deleted`, żeby
  // panel nie pokazał jako obecnego zasobu, o którym rdzeń właśnie powiedział,
  // że go nie ma.
  const odsubskrybuj = zrodlo.naZmianeZasobu((tresc) => {
    if (tresc.change === ChangeKind.Deleted) usunZasob(zapis, tresc.asset.id);
    else wchlonZasob(zapis, tresc.asset);
    rysuj();
  });

  // Pierwszego odczytu panel nie robi sam — pas paneli woła `odswiez()` na
  // wszystkich zbudowanych panelach (`pas-pomocniczych.ts`), więc czytanie
  // tutaj zleciłoby `design.asset.list` dwa razy pod rząd.
  rysuj();

  return { element, odswiez, zamknij: odsubskrybuj };
}
