import { TranslationStatus, type TranslationPanel } from '../../../../shared/contract';
import { opisOdmowy } from '../../komponenty/odmowa';
import {
  poleWielowierszowe,
  przycisk,
  utworzWierszOdpowiedzi,
  type WierszOdpowiedzi,
} from '../../modele/kontrolki-formularza';
import { nazwaStanuPanelu, opisTonuPanelu } from './etykiety-translate';
import { utworzNarzedziaPanelu, type NarzedziaPanelu } from './narzedzia-panelu';
import { utworzStanOkna, type StanOkna, type ZdanieStanu } from './stan-okna-translate';
import type { StanTranslate } from './stan-translate';
import { rozbieznoscOdpowiedzi } from './zgodnosc-odpowiedzi';

/**
 * Jedna instancja Translation Panels stanowi panel jednego języka docelowego; błąd zostaje
 * w panelu, którego dotyczy, więc odmowa jednego panelu nie gasi drugiego.
 */
export interface PanelJezyka {
  element: HTMLElement;
  /** Czy panel jest wskazany do widoku porównawczego. */
  doPorownania(): boolean;
  /** Treść tłumaczenia widoczna w panelu — dla widoku porównawczego. */
  tresc(): string;
  /** Język docelowy panelu. */
  jezyk(): string;
  /** Przerysowuje panel z danych rdzenia. */
  odswiez(panel: TranslationPanel): void;
  /** Prowadzi ognisko do pola tłumaczenia tej instancji, drogą przejścia skrótem między panelami. */
  ogniskuj(): void;
  /** Zdejmuje nasłuchy instancji przed jej usunięciem z dokumentu, gdy rdzeń przestaje oddawać jej panel. */
  rozlacz(): void;
}

/** Elementy panelu przerysowywane danymi rdzenia obejmują tytuł, plakietkę stanu, ton, pole tłumaczenia oraz wykaz zastrzeżeń kontroli jakości. */
interface CzesciPanelu {
  element: HTMLElement;
  tytul: HTMLElement;
  plakietka: HTMLElement;
  /** Ton oddany przez rdzeń. */
  ton: HTMLElement;
  pole: HTMLTextAreaElement;
  zastrzezenia: HTMLElement;
  okno: StanOkna;
  narzedzia: NarzedziaPanelu;
}

export function utworzPanelJezyka(
  stan: StanTranslate,
  poczatkowy: TranslationPanel,
  naPorownanie: () => void,
): PanelJezyka {
  let biezacy = poczatkowy;
  const okno: StanOkna = utworzStanOkna(pustyPanel(poczatkowy.language));

  const tytul = document.createElement('h4');
  tytul.className = 'mt-panel__tytul';

  const plakietka = document.createElement('span');
  plakietka.className = 'dn-plakietka dn-plakietka--rola mt-panel__stan';

  // Ton stoi w nagłówku obok stanu: obie rzeczy opisują panel jako całość i przychodzą z rdzenia razem.
  const ton = document.createElement('span');
  ton.className = 'mt-panel__ton';

  const porownanie = document.createElement('input');
  porownanie.type = 'checkbox';
  porownanie.className = 'dn-check mt-panel__porownanie';
  porownanie.setAttribute('aria-label', 'Weź panel do widoku porównawczego');
  porownanie.addEventListener('change', naPorownanie);

  const naglowek = document.createElement('header');
  naglowek.className = 'mt-panel__naglowek';
  naglowek.append(tytul, ton, plakietka, porownanie);

  const tresc = poleWielowierszowe(
    { etykieta: 'Tłumaczenie', opis: 'Korekta ręczna wraca do rdzenia komendą zapisu.' },
    6,
  );

  const zastrzezenia = document.createElement('ul');
  zastrzezenia.className = 'mt-zastrzezenia';

  const odpowiedz = utworzWierszOdpowiedzi();
  const zapisz = przycisk('Zapisz korektę', 'dn-btn dn-btn--sm dn-btn--atrament');

  const narzedzia = utworzNarzedziaPanelu(
    stan.panele,
    {
      idPanelu: () => biezacy.id,
      // Pamięć tłumaczeń pyta o segment źródłowy: pierwszy trafia do niej, a bez podziału cały tekst.
      segment: () => stan.segmenty()[0] ?? stan.tekstZrodlowy(),
      // Treść widziana przez operatora, nie ostatnio wchłonięta z rdzenia: zwrotne porównuje to, co widać.
      tresc: () => tresc.kontrolka.value,
      wchlon: (panel) => stan.wchlonPanel(panel),
      kanaly: stan.kanaly,
    },
    odpowiedz,
    poczatkowy.id,
  );

  okno.tresc.append(tresc.element, zapisz, odpowiedz.element, narzedzia.element, zastrzezenia);

  const element = document.createElement('article');
  element.className = 'mt-panel';
  element.dataset['panel'] = poczatkowy.id;
  element.append(naglowek, okno.element);

  const czesci: CzesciPanelu = {
    element,
    tytul,
    plakietka,
    ton,
    pole: tresc.kontrolka,
    zastrzezenia,
    okno,
    narzedzia,
  };

  zapisz.addEventListener('click', () => {
    void zapiszKorekte(stan, biezacy.id, tresc.kontrolka, okno, odpowiedz);
  });

  function odswiez(panel: TranslationPanel): void {
    biezacy = panel;
    przerysuj(panel, czesci);
  }

  odswiez(poczatkowy);

  return {
    element,
    doPorownania: () => porownanie.checked,
    tresc: () => tresc.kontrolka.value,
    jezyk: () => biezacy.language,
    odswiez,
    ogniskuj: () => tresc.kontrolka.focus(),
    rozlacz: () => narzedzia.rozlacz(),
  };
}

/**
 * Zapis korekty operatora potwierdza treść, która wróciła, odczytaną raz przed wysyłką, bo
 * operator pisze dalej, gdy żądanie jest w drodze.
 */
async function zapiszKorekte(
  stan: StanTranslate,
  idPanelu: string,
  pole: HTMLTextAreaElement,
  okno: StanOkna,
  odpowiedz: WierszOdpowiedzi,
): Promise<void> {
  const wyslana = pole.value;
  odpowiedz.pokaz('Zapisywanie korekty…', true);
  okno.ladowanie('Zapis korekty w rdzeniu…');
  const wynik = await stan.panele.zapiszKorekte(idPanelu, wyslana);
  if (!wynik.udany || wynik.wynik === undefined) {
    const zdanie = opisOdmowy('Zapis korekty', wynik.blad?.code, wynik.blad?.message);
    odpowiedz.pokaz(zdanie, false);
    okno.blad(zdanie);
    return;
  }
  const panel = wynik.wynik.panel;
  stan.wchlonPanel(panel);
  const rozbiezne = rozbieznoscOdpowiedzi('Zapis korekty', [
    { nazwa: 'treść panelu', zamowione: wyslana, oddane: panel.text ?? '' },
  ]);
  if (rozbiezne !== null) {
    odpowiedz.pokaz(rozbiezne, false);
    okno.blad(rozbiezne);
    return;
  }
  odpowiedz.pokaz('Korekta zapisana w rdzeniu — oddana treść panelu zgadza się z wysłaną.', true);
}

/**
 * Przerysowanie panelu danymi rdzenia dobiera stan obowiązkowy; ton widać w nagłówku oraz
 * w pasku narzędzi, oba z tej samej odpowiedzi rdzenia.
 */
function przerysuj(panel: TranslationPanel, czesci: CzesciPanelu): void {
  czesci.tytul.textContent = `Panel ${panel.language}`;
  czesci.plakietka.textContent = nazwaStanuPanelu(panel.status);
  czesci.ton.textContent = opisTonuPanelu(panel.tone);
  czesci.narzedzia.odswiez(panel);
  czesci.element.dataset['stan'] = panel.status;
  // Treść z rdzenia nie nadpisuje korekty: pole pod ogniskiem należy do operatora, dopóki go nie odda.
  if (document.activeElement !== czesci.pole) czesci.pole.value = panel.text ?? '';
  czesci.zastrzezenia.replaceChildren(...(panel.issues ?? []).map(wierszZastrzezenia));

  if (panel.status === TranslationStatus.Translating) {
    czesci.okno.ladowanie(`Rdzeń tłumaczy panel ${panel.language}…`);
    return;
  }
  if (panel.status === TranslationStatus.Error) {
    czesci.okno.blad(`Rdzeń zgłosił błąd tłumaczenia panelu ${panel.language}.`);
    return;
  }
  if ((panel.text ?? '') === '') {
    czesci.okno.puste(pustyPanel(panel.language));
    return;
  }
  czesci.okno.gotowe();
}

/**
 * Pustka panelu języka jest stanem oczekiwanym, nie usterką: panel zakłada się przed przekładem
 * i przez chwilę stoi bez treści.
 */
function pustyPanel(jezyk: string): ZdanieStanu {
  return {
    tytul: `Panel ${jezyk} bez przekładu`,
    opis:
      'Rdzeń nie oddał jeszcze treści tego panelu. Zapisz tekst źródłowy w Source Panel ' +
      'albo wpisz przekład ręcznie i zapisz korektę.',
  };
}

function wierszZastrzezenia(zastrzezenie: {
  kind: string;
  segment?: string;
  detail?: string;
}): HTMLElement {
  const element = document.createElement('li');
  element.className = 'mt-zastrzezenia__wiersz';
  element.dataset['rodzaj'] = zastrzezenie.kind;
  const czesci = [zastrzezenie.kind, zastrzezenie.segment, zastrzezenie.detail].filter(
    (czesc): czesc is string => czesc !== undefined && czesc !== '',
  );
  element.textContent = czesci.join(' — ');
  return element;
}
