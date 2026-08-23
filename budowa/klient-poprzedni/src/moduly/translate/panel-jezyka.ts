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
 * Jedna instancja Translation Panels — panel jednego języka docelowego.
 *
 * Okno jest instancją wielokrotną (N = liczba języków), więc panel jest tu tym,
 * czym wiersz w liście: jedną rzeczą na jeden plik. Okno zbiorcze nimi zarządza
 * i nie wie, jak są zbudowane.
 *
 * Błąd zostaje w panelu, którego dotyczy: każda instancja ma własny pas stanu
 * i własny wiersz odpowiedzi, więc odmowa kontroli jakości dla niemieckiego nie
 * gasi panelu francuskiego. Dlatego stan okna powstaje osobno dla każdego panelu,
 * a nie raz na okno zbiorcze.
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
  /**
   * Prowadzi ognisko do pola tłumaczenia tej instancji.
   *
   * Drogą jest przejście skrótem między panelami wymagającymi uwagi: panel
   * wskazany bez ogniska wymagałby jeszcze jednego kliknięcia, żeby zacząć
   * poprawiać, a skrót ma prowadzić do pracy, nie do widoku.
   */
  ogniskuj(): void;
  /**
   * Zdejmuje nasłuchy instancji przed jej usunięciem z dokumentu.
   *
   * Instancja znika, gdy rdzeń przestaje oddawać jej panel — a jej ster kanału
   * bywa wtedy rozwinięty i trzyma nasłuch na dokumencie. Samo `element.remove()`
   * takiego nasłuchu nie zdejmuje.
   */
  rozlacz(): void;
}

/** Elementy panelu przerysowywane danymi rdzenia. */
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

  // Ton stoi w nagłówku obok stanu, bo obie rzeczy opisują panel jako całość
  // i obie przychodzą z rdzenia w tej samej odpowiedzi.
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
      // Pamięć tłumaczeń pyta o segment źródłowy; bierzemy pierwszy segment
      // oddany przez rdzeń, a gdy podziału nie było — cały tekst źródłowy.
      segment: () => stan.segmenty()[0] ?? stan.tekstZrodlowy(),
      // Treść widziana przez Operatora, nie treść ostatnio wchłonięta z rdzenia:
      // tłumaczenie zwrotne ma porównać się z tym, co Operator ma przed oczami.
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
 * Zapis korekty Operatora — `translate.translation.set`.
 *
 * Zdanie potwierdza treść, która wróciła, a nie sam fakt odpowiedzi:
 * `TranslateTranslationSetResponse` niesie cały panel po korekcie razem z jego
 * treścią, więc zgodność da się sprawdzić bez dodatkowego wywołania.
 *
 * Treść odczytujemy raz, przed wysyłką. Operator pisze dalej, gdy żądanie jest
 * w drodze, więc porównanie odpowiedzi z polem odczytanym po powrocie zgłaszałoby
 * rozbieżność za każdym dopisanym znakiem.
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
 * Przerysowanie panelu danymi rdzenia wraz z doborem stanu obowiązkowego.
 *
 * Ton widać w dwóch miejscach i oba biorą z jednego źródła: nagłówek
 * (`mt-panel__ton`) oraz pole „Ton panelu" w pasku narzędzi dostają `panel.tone`
 * z tej samej odpowiedzi rdzenia, w tym samym przebiegu.
 */
function przerysuj(panel: TranslationPanel, czesci: CzesciPanelu): void {
  czesci.tytul.textContent = `Panel ${panel.language}`;
  czesci.plakietka.textContent = nazwaStanuPanelu(panel.status);
  czesci.ton.textContent = opisTonuPanelu(panel.tone);
  czesci.narzedzia.odswiez(panel);
  czesci.element.dataset['stan'] = panel.status;
  // Treść z rdzenia nie nadpisuje wpisywanej korekty: pole pod ogniskiem
  // należy do Operatora, dopóki go nie odda.
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
 * Pustka panelu języka — stan oczekiwany, nie usterka.
 *
 * Panel zakłada się przed przekładem i przez chwilę stoi bez treści; zdanie mówi
 * więc, czym ten panel jest i czym się go zapełnia, zamiast nazywać brak. Forma
 * stanu pustego jest w module jedna (tytuł nad opisem), różnicuje ją treść.
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
