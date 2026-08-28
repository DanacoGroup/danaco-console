import {
  TranslationStatus,
  type TranslateTargetAddRequest,
  type TranslationPanel,
} from '../../../../shared/contract';
import {
  poleTekstowe,
  przycisk,
  przyciskBezKomendy,
  utworzWierszOdpowiedzi,
  type WierszOdpowiedzi,
} from '../../modele/kontrolki-formularza';
import {
  BRAK_OKNA,
  OBJASNIENIA,
  PODPOWIEDZ_JEZYKOW,
  PUSTE,
  nazwaStanuPanelu,
} from './etykiety-translate';
import { dopnijDymek, naglowekOkna, podepnijPodpowiedz } from './kontrolki-translate';
import { utworzMapowanieStylow } from './mapowanie-stylow';
import { zdanieOdmowyModelu } from './odmowa-translate';
import { utworzPanelJezyka, type PanelJezyka } from './panel-jezyka';
import { utworzPorownaniePaneli } from './porownanie-paneli';
import { utworzStanOkna, type StanOkna } from './stan-okna-translate';
import type { StanTranslate } from './stan-translate';
import { utworzSterKanalu, type SterKanalu } from './ster-kanalu';
import { odswiezNoteSufituPaneli, utworzNoteSufituPaneli } from './sufit-paneli';
import { utworzRozwiniecie } from './warstwy-translate';
import { rozbieznoscOdpowiedzi } from './zgodnosc-odpowiedzi';

/**
 * Translation Panels to okno wiodące o instancji wielokrotnej modułu Translate; okno zbiorcze jest
 * zarządcą instancji, języki docelowe wybiera formularz, a przegląd dzieje się w każdym panelu osobno.
 */
export interface OknoTranslationPanels {
  element: HTMLElement;
  odswiez(): void;
  /** Prowadzi ognisko do pola języka docelowego — droga skrótu klawiszowego. */
  ogniskujDodanieJezyka(): void;
  /** Przewija do panelu wymagającego uwagi, poprzedniego albo następnego, ze stanem różnym od gotowego. */
  przejdzDoUwagi(wstecz: boolean): void;
  /** Zwija stery okna i wszystkich instancji — wołane przy rozłączeniu modułu. */
  rozlacz(): void;
}

const BRAK_PIVOTA =
  'Trasy przekładu przez język pośredni nie ma czym zlecić: żądanie dodania panelu niesie okno, ' +
  'język docelowy, ton i kanał modelu — pola na język pośredni w nim nie ma, a odpowiedź nie ' +
  'oddaje trasy, którą przekład poszedł.';

const BRAK_POROWNANIA_SILNIKOW =
  'Porównania wariantów kilku silników dla jednego segmentu nie ma czym zlecić: jedno żądanie ' +
  'przekładu wskazuje jeden kanał modelu i oddaje jeden panel. Warianty da się zestawić, zakładając ' +
  'panele tego samego języka na różnych kanałach.';

export function utworzOknoTranslationPanels(stan: StanTranslate): OknoTranslationPanels {
  const okno: StanOkna = utworzStanOkna(PUSTE.panele);
  const instancje = new Map<string, PanelJezyka>();

  const jezyk = poleTekstowe({ etykieta: 'Język docelowy', podpowiedz: 'np. de' });
  dopnijDymek(jezyk.element, OBJASNIENIA.jezykDocelowy);
  const podpowiedzi = podepnijPodpowiedz(jezyk.kontrolka, 'mt-jezyki-docelowe', PODPOWIEDZ_JEZYKOW);

  const ton = poleTekstowe({
    etykieta: 'Ton nowego panelu',
    podpowiedz: 'puste = ton domyślny rdzenia',
  });
  dopnijDymek(ton.element, OBJASNIENIA.tonPanelu);

  const sterKanalu: SterKanalu = utworzSterKanalu(stan.kanaly, {
    podpis: 'Kanał przekładu nowego panelu',
    czynnosc: 'przekład tekstu źródłowego',
  });

  const odpowiedz = utworzWierszOdpowiedzi();
  const dodaj = przycisk('+ Dodaj język', 'dn-btn dn-btn--sm dn-btn--atrament');
  dodaj.addEventListener('click', () => {
    void dodajJezyk(
      stan,
      { jezyk: jezyk.kontrolka, ton: ton.kontrolka, kanal: sterKanalu },
      okno,
      odpowiedz,
    );
  });

  const formularz = document.createElement('div');
  formularz.className = 'mt-dodaj-jezyk';
  formularz.append(jezyk.element, podpowiedzi, ton.element, sterKanalu.element, dodaj);

  // Nota sufitu stoi przy formularzu dodania panelu i nie jest bramką liczby paneli okna.
  const sufit = utworzNoteSufituPaneli(0);

  const siatka = document.createElement('div');
  siatka.className = 'mt-panele';

  const porownanie = utworzPorownaniePaneli();
  const mapowanie = utworzMapowanieStylow(() => stan.tekstZrodlowy());
  okno.tresc.append(
    formularz,
    sufit,
    odpowiedz.element,
    siatka,
    porownanie.element,
    silnikiIPivot(),
    mapowanie.element,
  );

  const element = document.createElement('section');
  element.className = 'mt-okno mt-okno--wiodace';
  element.dataset['okno'] = 'translation-panels';
  element.append(naglowekOkna('Translation Panels', 'wiodące'), okno.element);

  /** Przerysowanie zawsze odświeża siatkę paneli i nigdy nie zdejmuje fazy trwającej ani fazy błędu. */
  function odswiez(): void {
    const wykaz = stan.panelJezykow();
    siatka.replaceChildren(
      ...zestroj(wykaz, instancje, (panel) =>
        utworzPanelJezyka(stan, panel, odswiezPorownanie),
      ),
    );
    // Nota liczy panele, które okno zna w tej chwili, z tego samego wykazu, z którego powstaje siatka.
    odswiezNoteSufituPaneli(sufit, wykaz.length);
    sterKanalu.odswiez();
    odswiezPorownanie();
    if (okno.faza() === 'ladowanie' || okno.faza() === 'blad') return;
    if (wykaz.length === 0) {
      okno.puste(PUSTE.panele);
      return;
    }
    okno.gotowe();
  }

  function odswiezPorownanie(): void {
    porownanie.odswiez([...instancje.values()], stan.tekstZrodlowy());
  }

  odswiez();

  return {
    element,
    odswiez,

    ogniskujDodanieJezyka() {
      jezyk.kontrolka.focus();
      jezyk.kontrolka.select();
    },

    /** Przejście między panelami wymagającymi uwagi zaczyna od panelu z ogniskiem albo od krańca siatki. */
    przejdzDoUwagi(wstecz) {
      const doUwagi = stan
        .panelJezykow()
        .filter((panel) => panel.status !== TranslationStatus.Ready);
      if (doUwagi.length === 0) {
        odpowiedz.pokaz(
          'Żaden panel nie czeka na uwagę — rdzeń oddaje wszystkie jako gotowe.',
          true,
        );
        return;
      }
      const kolejnosc = wstecz ? [...doUwagi].reverse() : doUwagi;
      const biezacy = kolejnosc.findIndex((panel) =>
        instancje.get(panel.id)?.element.contains(document.activeElement),
      );
      const nastepny = kolejnosc[(biezacy + 1) % kolejnosc.length];
      if (nastepny === undefined) return;
      const instancja = instancje.get(nastepny.id);
      if (instancja === undefined) return;
      instancja.element.scrollIntoView({ block: 'nearest' });
      instancja.ogniskuj();
      odpowiedz.pokaz(
        `Panel ${nastepny.language} — ${nazwaStanuPanelu(nastepny.status)}; paneli czekających ` +
          `na uwagę: ${String(doUwagi.length)}.`,
        true,
      );
    },

    rozlacz() {
      sterKanalu.zwin();
      for (const instancja of instancje.values()) instancja.rozlacz();
    },
  };
}

/**
 * Warstwa czwarta: trasa przekładu przez język pośredni i porównanie silników.
 *
 * Obie pozycje stoją razem, bo obie rozbijają się o ten sam kształt żądania:
 * jedno wywołanie przekładu to jeden język docelowy i jeden kanał modelu.
 */
function silnikiIPivot(): HTMLElement {
  const rozwiniecie = utworzRozwiniecie({
    warstwa: 4,
    nazwa: 'Trasa przekładu i porównanie silników',
    wyjasnienie: 'Tłumaczenie przez język pośredni oraz zestawienie wariantów kilku silników.',
    znacznik: '☰',
  });
  rozwiniecie.tresc.append(
    przyciskBezKomendy('Konfiguracja języka pośredniego', BRAK_PIVOTA),
    przyciskBezKomendy('Porównanie wariantów silników', BRAK_POROWNANIA_SILNIKOW),
  );
  return rozwiniecie.element;
}

/** Pola formularza Dodaj język widziane przez ścieżkę zapisu: kod języka docelowego, ton tłumaczenia oraz ster kanału modelu wykonującego przekład. */
interface PolaJezyka {
  jezyk: HTMLInputElement;
  ton: HTMLInputElement;
  /** Ster kanału modelu; pusty wybór znaczy „kanał czynny okna". */
  kanal: SterKanalu;
}

/**
 * Dodanie języka docelowego zakłada nowy panel komendą modelu; nieudane dodanie stawia okno
 * w fazie błędu, którą zdejmuje dopiero dodanie udane.
 */
async function dodajJezyk(
  stan: StanTranslate,
  pola: PolaJezyka,
  okno: StanOkna,
  odpowiedz: WierszOdpowiedzi,
): Promise<void> {
  if (stan.idOkna() === '') {
    // Odmowa, nie pustka: brak okna sesji jest powodem niewykonania, idzie do odpowiedzi i w stan błędu.
    odpowiedz.pokaz(BRAK_OKNA, false);
    okno.blad(BRAK_OKNA);
    return;
  }
  const kod = pola.jezyk.value.trim();
  if (kod === '') {
    odpowiedz.pokaz('Wskaż język docelowy — rdzeń odmówi założenia panelu bez niego.', false);
    return;
  }
  const zadanie: TranslateTargetAddRequest = { windowId: stan.idOkna(), language: kod };
  const wybranyTon = pola.ton.value.trim();
  if (wybranyTon !== '') zadanie.tone = wybranyTon;

  // Kanał wchodzi do żądania wyłącznie wskazany; kontrakt bierze przy braku pola kanał czynny okna.
  const wybranyKanal = pola.kanal.wybrany();
  if (wybranyKanal !== '') zadanie.channelId = wybranyKanal;

  // Dodanie języka idzie modelem i trwa — okno mówi to stanem ładowania, nie tylko wierszem odpowiedzi.
  okno.ladowanie(`Rdzeń zakłada panel języka ${kod} i przekłada tekst źródłowy modelem.`);
  odpowiedz.pokaz(`Zakładanie panelu języka ${kod}…`, true);
  const wynik = await stan.panele.dodajJezyk(zadanie);
  if (!wynik.udany || wynik.wynik === undefined) {
    const zdanie = zdanieOdmowyModelu('Dodanie języka docelowego', wynik.blad);
    odpowiedz.pokaz(zdanie, false);
    okno.blad(zdanie);
    return;
  }
  const panel = wynik.wynik.panel;
  stan.wchlonPanel(panel);
  pola.jezyk.value = '';

  // Żądanie niosło kod języka i ton, odpowiedź niesie cały panel, więc pola sprawdza się bez wywołania.
  const rozbiezne = rozbieznoscOdpowiedzi('Dodanie języka docelowego', [
    { nazwa: 'język panelu', zamowione: kod, oddane: panel.language },
    ...(wybranyTon === ''
      ? []
      : [{ nazwa: 'ton panelu', zamowione: wybranyTon, oddane: panel.tone }]),
  ]);
  if (rozbiezne !== null) {
    odpowiedz.pokaz(rozbiezne, false);
    okno.blad(rozbiezne);
    return;
  }
  okno.gotowe();
  odpowiedz.pokaz(`Panel języka ${panel.language} założony w rdzeniu.`, true);
}

/**
 * Zestrojenie instancji z wykazem rdzenia usuwa panel, którego rdzeń już nie oddaje, i dodaje
 * instancję nowemu panelowi w kolejności wykazu rdzenia.
 */
function zestroj(
  wykaz: readonly TranslationPanel[],
  instancje: Map<string, PanelJezyka>,
  utworz: (panel: TranslationPanel) => PanelJezyka,
): HTMLElement[] {
  for (const [id, instancja] of [...instancje]) {
    if (wykaz.some((panel) => panel.id === id)) continue;
    // Zdjęcie nasłuchów przed usunięciem, bo ster kanału zostawiłby nasłuch po usunięciu elementu.
    instancja.rozlacz();
    instancja.element.remove();
    instancje.delete(id);
  }
  return wykaz.map((panel) => {
    const istniejaca = instancje.get(panel.id);
    if (istniejaca === undefined) {
      const nowa = utworz(panel);
      instancje.set(panel.id, nowa);
      return nowa.element;
    }
    istniejaca.odswiez(panel);
    return istniejaca.element;
  });
}
