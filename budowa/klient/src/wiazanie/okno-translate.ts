/**
 * Wiązanie wnętrza okna modułu Translate z rdzeniem. Znacznik należy do
 * Właściciela — ten plik wypełnia stojące węzły odpowiedziami rodziny
 * translate.*, powiela wzory zdjęte z treści przykładowej, a węzły, dla
 * których kontrakt nie ma pola, zdejmuje.
 */
import {
  Command,
  EventType,
  type GlossaryTerm,
  type TranslationPanel,
} from '../../../shared/contract.ts';
import type { Kanal } from '../protokol/kanal.ts';
import { wywolaj } from '../protokol/wywolanie.ts';

/** Panele prototypu, którym rodzina translate.* nie oddaje ani jednego pola; schodzą razem z przełącznikami menu sesji. */
const PANELE_BEZ_POKRYCIA = ['panel-plan', 'panel-kolejka', 'panel-zadania', 'panel-pliki'];

/** Wzory zdjęte z treści przykładowej, zanim wykazy zostaną wyczyszczone; klon zachowuje ikonę i klasy nadane przez bibliotekę. */
interface Wzory {
  segment: HTMLElement | null;
  segmentRewizji: HTMLElement | null;
  kolumna: HTMLElement | null;
  termin: HTMLElement | null;
  para: HTMLElement | null;
  bezTlumaczenia: HTMLElement | null;
}

/** Stan okna tłumaczenia: tekst źródłowy ustalony w tym oknie, panele języków docelowych i panel wybrany. */
interface Stan {
  zrodlo: string;
  panele: TranslationPanel[];
  panel: string;
}

/**
 * Wiąże wnętrze okna modułu Translate. Brak węzłów wnętrza kończy funkcję bez
 * żadnego działania — okna w ramie nie ma, więc nie ma czego wypełniać.
 */
export function zwiazTranslate(kanal: Kanal, idOkna: string): void {
  const cialo = document.querySelector<HTMLElement>('.sta-cialo');
  if (cialo === null) return;

  const wzory = zdejmijWzory(cialo);
  zdejmijTrescPrzykladowa(cialo);
  zdejmijSterowanieBezPokrycia(cialo);

  const stan: Stan = { zrodlo: '', panele: [], panel: '' };

  /* Odwołanie stoi wartością, nie deklaracją: deklaracja wchodzi na wierzch
     zasięgu, więc kompilator bierze dla węzła wnętrza typ sprzed rozpoznania
     pustki. */
  /** Przerysowuje panele języków docelowych ze stanu zebranego odpowiedziami rdzenia. */
  const odswiez = (): void => {
    wypelnijPanele(cialo, wzory, stan);
  };

  zwiazZrodlo(kanal, idOkna, cialo, wzory, stan, odswiez);
  zwiazPanele(kanal, cialo, wzory, stan, odswiez);
  zwiazGlosariusz(kanal, cialo, wzory, stan);
  kanal.naZdarzenie(EventType.TranslateTranslationChanged, (tresc) => {
    if (tresc.panel.windowId !== idOkna) return;
    stan.panele = [
      ...stan.panele.filter((panel) => panel.id !== tresc.panel.id),
      tresc.panel,
    ];
    odswiez();
  });

  odswiez();
  void wypelnijGlosariusz(kanal, cialo, wzory, '');
}

/** Zdejmuje wzory segmentów, kolumn i terminów, zanim wykazy zostaną wyczyszczone; wzór traci oznaczenia bez pokrycia w kontrakcie. */
function zdejmijWzory(cialo: HTMLElement): Wzory {
  const kolumna = sklonuj(cialo.querySelector('#panel-obszar-tlum .tr-kol'));
  // Miara gotowości panelu nie ma pola: panel niesie stan, ton, niezgodności
  // i czas zmiany, nie licznik segmentów przełożonych.
  kolumna?.querySelector('.dn-plakietka')?.classList.remove('dn-plakietka--sukces');
  const segmenty = [...(kolumna?.querySelectorAll<HTMLElement>('.tr-seg') ?? [])];
  const segment = sklonuj(segmenty[0] ?? null);
  const segmentRewizji = sklonuj(segmenty[3] ?? null);
  kolumna?.querySelector('.tr-kol-tresc')?.replaceChildren();
  const termin = sklonuj(cialo.querySelector('#panel-glossary .tr-glos'));
  termin?.removeAttribute('style');
  termin?.querySelector('div[style]')?.remove();
  const para = sklonuj(termin?.querySelector('.tr-glos-pary span') ?? null);
  termin?.querySelector('.tr-glos-pary')?.replaceChildren();
  const przyciski = [...(termin?.querySelectorAll('button') ?? [])];
  // Wykaz wystąpień terminu nie ma w oknie węzła, w którym mógłby stanąć,
  // a zmiana terminu żąda pól, których okno nie niesie.
  przyciski[1]?.remove();
  przyciski[2]?.remove();
  return {
    segment,
    segmentRewizji,
    kolumna,
    termin,
    para,
    bezTlumaczenia: sklonuj(
      cialo.querySelectorAll('#panel-glossary .tr-glos')[1]?.querySelector('.dn-plakietka') ?? null,
    ),
  };
}

/**
 * Zdejmuje treść przykładową wnętrza. Wykaz tłumaczeń szyny, historia rozmowy
 * i słownik kontekstu wiszą na rodzinach spoza translate.*; liczba słów,
 * pokrycie pamięcią tłumaczeń i nazwa projektu nie mają pola w całym
 * kontrakcie.
 */
function zdejmijTrescPrzykladowa(cialo: HTMLElement): void {
  cialo.querySelector('.dn-szyna-modulu-lista')?.replaceChildren();
  for (const znacznik of cialo.querySelectorAll('#menu-filtr .sta-menu-poz[aria-checked]')) {
    znacznik.removeAttribute('aria-checked');
  }
  cialo.querySelector('.sta-kom-historia')?.replaceChildren();
  cialo.querySelector('.sta-kom-monitor')?.remove();
  zdejmijPoleNaglowka(cialo, 'Wysiłek');
  cialo.querySelector('.sta-kom-naglowek .sta-kom-stan')?.remove();
  cialo.querySelectorAll('.sta-kom-kontekst > .sta-zrodlo')[1]?.remove();
  // Odpięcie znaku kontekstu nie ma komendy w żadnej z rodzin tego okna.
  for (const odepnij of cialo.querySelectorAll('.sta-kom-kontekst .odepnij')) odepnij.remove();
  zdejmijZnakiCzynnosci(cialo);
  zdejmijNaglowekZrodla(cialo);
  cialo.querySelector('#panel-obszar-tlum .tr-siatka')?.replaceChildren();
  cialo.querySelector('#panel-obszar-tlum .sta-okno-tresc')?.replaceChildren();
  for (const stojacy of cialo.querySelectorAll('#panel-glossary .tr-glos')) stojacy.remove();
}

/** Zdejmuje pole nagłówka rozpoznane po podpisie; rodzina translate.* nie oddaje nakładu rozumowania kanału. */
function zdejmijPoleNaglowka(cialo: HTMLElement, podpis: string): void {
  const pola = [...cialo.querySelectorAll('.sta-kom-naglowek .sta-kom-pole')];
  pola.find((pole) => pole.textContent?.startsWith(podpis) === true)?.remove();
}

/**
 * Zostawia w pasie czynności sam licznik segmentów. Zasięg wykonania, nazwa
 * projektu, znak operacji i pokrycie pamięcią tłumaczeń nie mają pola, a
 * zlecenie wsadowe żąda wykazu operacji, którego okno nie niesie.
 */
function zdejmijZnakiCzynnosci(cialo: HTMLElement): void {
  const znaki = [...cialo.querySelectorAll<HTMLElement>('.sta-kontekst-akcji .sta-chip')];
  znaki[0]?.remove();
  znaki[1]?.remove();
  znaki[3]?.remove();
  znaki[4]?.remove();
  cialo.querySelector('.sta-kontekst-akcji button')?.remove();
  wpiszLiczbe(znaki[2] ?? null, 0);
}

/**
 * Porządkuje nagłówek panelu źródłowego. Liczba słów nie ma pola, a wczytanie
 * dokumentu żąda ścieżki i języka docelowego, których okno nie niesie; język
 * źródłowy zostaje oznaczony do wypełnienia odpowiedzią rdzenia.
 */
function zdejmijNaglowekZrodla(cialo: HTMLElement): void {
  const glowa = cialo.querySelector<HTMLElement>('#panel-obszar-tlum .tr-glowa');
  if (glowa === null) return;
  glowa.querySelector('span:not(.tr-kod):not(.rozciag)')?.remove();
  const przyciski = [...glowa.querySelectorAll<HTMLElement>('button')];
  const jezyk = przyciski[0];
  if (jezyk !== undefined) {
    jezyk.dataset.pole = 'jezyk';
    /* Język przykładowy schodzi przed pierwszą odpowiedzią: pole puste jest
       uczciwe, pole z cudzym językiem źródłowym — nie. */
    jezyk.textContent = '';
  }
  przyciski[2]?.remove();
  cialo.querySelector('#panel-obszar-tlum .sta-okno-belka .dn-btn--atrament')?.remove();
  cialo.querySelector('#panel-obszar-tlum .sta-okno-belka .sta-chip')?.remove();
}

/**
 * Zdejmuje sterowanie, którego rodzina translate.* nie obsługuje: panele bez
 * pokrycia wraz z ich przełącznikami, urządzenia wejścia dźwięku, wybór modelu
 * i nakładu, menu dodawania oraz wydania żądające formatu i ścieżki.
 */
function zdejmijSterowanieBezPokrycia(cialo: HTMLElement): void {
  for (const nazwa of PANELE_BEZ_POKRYCIA) {
    cialo.querySelector(`#${nazwa}`)?.remove();
    cialo.querySelector(`[data-panel-toggle="${nazwa}"]`)?.remove();
  }
  for (const przelacznik of cialo.querySelectorAll<HTMLElement>('[data-panel-toggle]')) {
    const panel = cialo.querySelector(`#${przelacznik.dataset.panelToggle ?? ''}`);
    przelacznik.setAttribute('aria-checked', String(panel?.hasAttribute('hidden') === false));
  }
  for (const pozycja of cialo.querySelectorAll('#menu-sesja .sta-menu-poz[role="menuitem"]')) {
    pozycja.remove();
  }
  cialo.querySelector('#menu-sesja .sta-menu-sep')?.remove();
  cialo.querySelector('#okno-czat-1 .sta-okno-belka .dn-btn-ikona')?.remove();
  cialo.querySelector('#pop-mik')?.closest('.sta-nrz')?.remove();
  cialo.querySelector('#pop-model')?.closest('.sta-nrz')?.remove();
  cialo.querySelector('#pop-wysilek')?.closest('.sta-nrz')?.remove();
  cialo.querySelector('#pop-plus')?.closest('.sta-nrz')?.remove();
  // Podpis piątego trybu opisuje domyślność, której rejestr uprawnień nie zna.
  cialo.querySelector('#pop-tryb-upr .sta-popover-wiersz small')?.remove();
  zdejmijCzynnosciBezPokrycia(cialo);
}

/** Zdejmuje czynności żądające wartości, których okno nie niesie: dodanie języka, wydanie panelu i obsługę bazy terminów plikiem. */
function zdejmijCzynnosciBezPokrycia(cialo: HTMLElement): void {
  const stopka = [...cialo.querySelectorAll('#panel-obszar-tlum .dn-wykaz-modulu button')];
  stopka[2]?.remove();
  cialo.querySelector('#panel-glossary .tr-glowa button')?.remove();
  cialo.querySelector('#panel-glossary .dn-wykaz-modulu')?.remove();
}

/** Stawia po jednej kolumnie na panel języka docelowego wraz z jego niezgodnościami; brak paneli zostawia obszar pusty. */
function wypelnijPanele(cialo: HTMLElement, wzory: Wzory, stan: Stan): void {
  const siatka = cialo.querySelector<HTMLElement>('#panel-obszar-tlum .tr-siatka');
  const wzor = wzory.kolumna;
  if (siatka === null || wzor === null) return;
  siatka.replaceChildren();
  if (!stan.panele.some((panel) => panel.id === stan.panel)) {
    stan.panel = stan.panele[0]?.id ?? '';
  }
  for (const panel of stan.panele) {
    const kolumna = wzor.cloneNode(true) as HTMLElement;
    kolumna.dataset.panel = panel.id;
    opiszKolumne(kolumna, panel);
    wypelnijTresc(kolumna, wzory, panel);
    siatka.appendChild(kolumna);
  }
  opiszJezyki(cialo, stan);
}

/** Opisuje głowę kolumny językiem panelu, jego tonem i stanem; ton nieustalony zdejmuje przycisk tonu. */
function opiszKolumne(kolumna: HTMLElement, panel: TranslationPanel): void {
  const kod = kolumna.querySelector('.tr-kod');
  if (kod !== null) kod.textContent = panel.language;
  const ton = kolumna.querySelector<HTMLElement>('.tr-kol-glowa button');
  if (ton !== null) {
    if (panel.tone === undefined || panel.tone === '') ton.remove();
    else ton.textContent = panel.tone;
  }
  const plakietka = kolumna.querySelector('.dn-plakietka');
  if (plakietka !== null) plakietka.textContent = panel.status;
}

/** Wypełnia treść kolumny przekładem panelu i wierszem na każdą niezgodność kontroli jakości; przekład pusty zostawia kolumnę bez wiersza treści. */
function wypelnijTresc(kolumna: HTMLElement, wzory: Wzory, panel: TranslationPanel): void {
  const tresc = kolumna.querySelector<HTMLElement>('.tr-kol-tresc');
  if (tresc === null) return;
  tresc.replaceChildren();
  if (panel.text !== undefined && panel.text !== '' && wzory.segment !== null) {
    const wiersz = wzory.segment.cloneNode(true) as HTMLElement;
    // Przekład panelu jest jedną treścią: kontrakt nie wiąże go z numerem
    // segmentu tekstu źródłowego.
    wiersz.querySelector('.tr-seg-nr')?.remove();
    wpiszOstatni(wiersz, panel.text);
    tresc.appendChild(wiersz);
  }
  if (wzory.segmentRewizji === null) return;
  for (const niezgodnosc of panel.issues ?? []) {
    const wiersz = wzory.segmentRewizji.cloneNode(true) as HTMLElement;
    wiersz.querySelector('.tr-seg-nr')?.remove();
    wpiszOstatni(wiersz, niezgodnosc.detail ?? niezgodnosc.kind);
    tresc.appendChild(wiersz);
  }
}

/**
 * Nanosi języki paneli na znak belki i znak kontekstu okna komunikacji. Oba
 * znaki zostają puste, dopóki panelu nie ma: panel powstaje dopiero w toku
 * pracy okna, więc znak czeka na wartość, a nie jest znakiem bez pokrycia.
 */
function opiszJezyki(cialo: HTMLElement, stan: Stan): void {
  const jezyki = stan.panele.map((panel) => panel.language).join(' · ');
  opiszWezel(cialo.querySelector('#okno-czat-1 .sta-okno-belka > .sta-chip'), jezyki);
  const wybrany = stan.panele.find((panel) => panel.id === stan.panel);
  opiszWezel(cialo.querySelector('.sta-kom-kontekst > .sta-zrodlo'), wybrany?.language ?? '');
}

/** Wypełnia panel źródłowy segmentami rozpoznanymi przez rdzeń i opisuje jego nagłówek językiem oraz liczbą segmentów. */
async function wypelnijZrodlo(
  kanal: Kanal,
  cialo: HTMLElement,
  wzory: Wzory,
  stan: Stan,
): Promise<void> {
  const wynik = await wywolaj(kanal, Command.TranslateSourceSegment, { text: stan.zrodlo });
  const segmenty = wynik.udany ? (wynik.wynik?.segments ?? []) : [];
  const obszar = cialo.querySelector<HTMLElement>('#panel-obszar-tlum .sta-okno-tresc');
  if (obszar === null || wzory.segment === null) return;
  obszar.replaceChildren();
  for (const [numer, segment] of segmenty.entries()) {
    const wiersz = wzory.segment.cloneNode(true) as HTMLElement;
    const kolejny = wiersz.querySelector('.tr-seg-nr');
    if (kolejny !== null) kolejny.textContent = String(numer + 1);
    wpiszOstatni(wiersz, segment);
    obszar.appendChild(wiersz);
  }
  wpiszLiczbe(cialo.querySelector<HTMLElement>('.sta-kontekst-akcji .sta-chip'), segmenty.length);
}

/** Wiąże panel źródłowy: wklejenie ustala tekst źródłowy okna, a ponowna segmentacja przepuszcza go przez rdzeń raz jeszcze. */
function zwiazZrodlo(
  kanal: Kanal,
  idOkna: string,
  cialo: HTMLElement,
  wzory: Wzory,
  stan: Stan,
  odswiez: () => void,
): void {
  const przyciski = [...cialo.querySelectorAll<HTMLElement>('#panel-obszar-tlum .tr-glowa button')];
  const wklej = przyciski.find((przycisk) => przycisk.dataset.pole === undefined);
  wklej?.addEventListener('click', () => {
    void wklejZrodlo(kanal, idOkna, cialo, wzory, stan, odswiez);
  });
  const ponownie = przyciski[przyciski.length - 1];
  if (ponownie === wklej) return;
  ponownie?.addEventListener('click', () => {
    if (stan.zrodlo === '') return;
    void ustalZrodlo(kanal, idOkna, cialo, wzory, stan, stan.zrodlo, true, odswiez);
  });
}

/** Bierze tekst źródłowy ze schowka Operatora i oddaje go rdzeniowi; schowek zamknięty kończy czynność bez zmiany okna. */
async function wklejZrodlo(
  kanal: Kanal,
  idOkna: string,
  cialo: HTMLElement,
  wzory: Wzory,
  stan: Stan,
  odswiez: () => void,
): Promise<void> {
  let tekst = '';
  try {
    tekst = await navigator.clipboard.readText();
  } catch {
    return;
  }
  if (tekst.trim() === '') return;
  await ustalZrodlo(kanal, idOkna, cialo, wzory, stan, tekst, false, odswiez);
}

/** Ustala tekst źródłowy okna; odpowiedź niesie język źródłowy i panele, więc po niej idzie opis nagłówka i przerysowanie kolumn. */
async function ustalZrodlo(
  kanal: Kanal,
  idOkna: string,
  cialo: HTMLElement,
  wzory: Wzory,
  stan: Stan,
  tekst: string,
  ponownie: boolean,
  odswiez: () => void,
): Promise<void> {
  const wynik = await wywolaj(kanal, Command.TranslateSourceSet, {
    windowId: idOkna,
    text: tekst,
    resegment: ponownie,
  });
  if (!wynik.udany || wynik.wynik === undefined) return;
  stan.zrodlo = tekst;
  stan.panele = wynik.wynik.panels;
  const jezyk = cialo.querySelector<HTMLElement>('#panel-obszar-tlum [data-pole="jezyk"]');
  if (jezyk !== null) jezyk.textContent = wynik.wynik.sourceLanguage;
  odswiez();
  await wypelnijZrodlo(kanal, cialo, wzory, stan);
}

/** Wiąże kolumny paneli: wskazanie kolumny czyni panel wybranym, a stopka obszaru zleca tłumaczenie zwrotne i kontrolę jakości panelu wybranego. */
function zwiazPanele(
  kanal: Kanal,
  cialo: HTMLElement,
  wzory: Wzory,
  stan: Stan,
  odswiez: () => void,
): void {
  cialo.querySelector('#panel-obszar-tlum .tr-siatka')?.addEventListener('click', (zdarzenie) => {
    const cel = zdarzenie.target;
    if (!(cel instanceof Element)) return;
    const panel = cel.closest<HTMLElement>('.tr-kol')?.dataset.panel;
    if (panel === undefined) return;
    stan.panel = panel;
    odswiez();
  });
  const stopka = [...cialo.querySelectorAll<HTMLElement>('#panel-obszar-tlum .dn-wykaz-modulu button')];
  stopka[0]?.addEventListener('click', () => {
    void tlumaczZwrotnie(kanal, stan);
  });
  stopka[1]?.addEventListener('click', () => {
    void sprawdzJakosc(kanal, cialo, wzory, stan);
  });
}

/** Zleca tłumaczenie zwrotne panelu wybranego; odpowiedź wraca treścią, którą przyjmuje panel jako korektę. */
async function tlumaczZwrotnie(kanal: Kanal, stan: Stan): Promise<void> {
  if (stan.panel === '') return;
  await wywolaj(kanal, Command.TranslateBacktranslationRun, { panelId: stan.panel });
}

/** Zleca kontrolę jakości panelu wybranego i nanosi jej niezgodności na kolumnę panelu. */
async function sprawdzJakosc(
  kanal: Kanal,
  cialo: HTMLElement,
  wzory: Wzory,
  stan: Stan,
): Promise<void> {
  if (stan.panel === '' || wzory.segmentRewizji === null) return;
  const wynik = await wywolaj(kanal, Command.TranslateQualityCheck, { panelId: stan.panel });
  if (!wynik.udany || wynik.wynik === undefined) return;
  const kolumna = cialo.querySelector<HTMLElement>(`.tr-kol[data-panel="${stan.panel}"]`);
  const tresc = kolumna?.querySelector<HTMLElement>('.tr-kol-tresc');
  if (tresc === null || tresc === undefined) return;
  for (const uwaga of wynik.wynik.issues) {
    const wiersz = wzory.segmentRewizji.cloneNode(true) as HTMLElement;
    wiersz.querySelector('.tr-seg-nr')?.remove();
    wpiszOstatni(wiersz, uwaga);
    tresc.appendChild(wiersz);
  }
}

/** Wypełnia bazę terminologiczną terminami rdzenia; jeden termin źródłowy zbiera pary wszystkich swoich języków w jednej pozycji. */
async function wypelnijGlosariusz(
  kanal: Kanal,
  cialo: HTMLElement,
  wzory: Wzory,
  fraza: string,
): Promise<void> {
  const panel = cialo.querySelector<HTMLElement>('#panel-glossary .sta-okno-tresc');
  const wzor = wzory.termin;
  if (panel === null || wzor === null) return;
  const wynik = await wywolaj(kanal, Command.TranslateGlossaryList, {
    query: fraza === '' ? undefined : fraza,
  });
  const terminy = wynik.udany ? (wynik.wynik?.terms ?? []) : [];
  for (const stojacy of panel.querySelectorAll('.tr-glos')) stojacy.remove();
  for (const [zrodlo, pozycje] of pogrupuj(terminy)) {
    panel.appendChild(zbudujTermin(wzory, wzor, zrodlo, pozycje));
  }
}

/** Grupuje terminy po brzmieniu źródłowym; jedna pozycja bazy niesie po jednym odpowiedniku na język. */
function pogrupuj(terminy: GlossaryTerm[]): Map<string, GlossaryTerm[]> {
  const spis = new Map<string, GlossaryTerm[]>();
  for (const termin of terminy) {
    const pozycje = spis.get(termin.source) ?? [];
    pozycje.push(termin);
    spis.set(termin.source, pozycje);
  }
  return spis;
}

/** Zwraca pozycję bazy terminologicznej opisaną brzmieniem źródłowym, dziedziną, oznaczeniem nietłumaczonego i parami odpowiedników. */
function zbudujTermin(
  wzory: Wzory,
  wzor: HTMLElement,
  zrodlo: string,
  pozycje: GlossaryTerm[],
): HTMLElement {
  const termin = wzor.cloneNode(true) as HTMLElement;
  termin.dataset.termin = zrodlo;
  const brzmienie = termin.querySelector('b');
  if (brzmienie !== null) brzmienie.textContent = zrodlo;
  const dziedzina = termin.querySelector<HTMLElement>('.dn-meta');
  const nazwaDziedziny = pozycje.find((pozycja) => pozycja.domain !== undefined)?.domain ?? '';
  if (dziedzina !== null) {
    if (nazwaDziedziny === '') dziedzina.remove();
    else dziedzina.textContent = nazwaDziedziny;
  }
  if (pozycje.some((pozycja) => pozycja.doNotTranslate === true) && wzory.bezTlumaczenia !== null) {
    brzmienie?.after(wzory.bezTlumaczenia.cloneNode(true));
  }
  const pary = termin.querySelector<HTMLElement>('.tr-glos-pary');
  if (pary === null || wzory.para === null) return termin;
  pary.replaceChildren();
  for (const pozycja of pozycje) {
    if (pozycja.target === undefined) continue;
    const para = wzory.para.cloneNode(true) as HTMLElement;
    wpiszTekst(para, `${pozycja.language}: `);
    const odpowiednik = para.querySelector('b');
    if (odpowiednik !== null) odpowiednik.textContent = pozycja.target;
    pary.appendChild(para);
  }
  return termin;
}

/** Wiąże bazę terminologiczną: fraza zawęża wykaz terminów, a ujednolicenie stosuje je do panelu wybranego. */
function zwiazGlosariusz(kanal: Kanal, cialo: HTMLElement, wzory: Wzory, stan: Stan): void {
  const szukanie = cialo.querySelector<HTMLInputElement>('#panel-glossary input[type="search"]');
  szukanie?.addEventListener('input', () => {
    void wypelnijGlosariusz(kanal, cialo, wzory, szukanie.value.trim());
  });
  cialo.querySelector('#panel-glossary .sta-okno-tresc')?.addEventListener('click', (zdarzenie) => {
    const cel = zdarzenie.target;
    if (!(cel instanceof Element) || cel.closest('button') === null) return;
    void wywolaj(kanal, Command.TranslateGlossaryApply, {
      panelId: stan.panel === '' ? undefined : stan.panel,
    });
  });
}

/**
 * Podmienia ciąg cyfr w treści znaku liczbą odpowiedzi, zostawiając podpis
 * znaku. Znak, który cyfr nie niesie, schodzi — podpis z liczbą przykładową
 * kłamałby.
 */
function wpiszLiczbe(wezel: HTMLElement | null, liczba: number): void {
  if (wezel === null) return;
  for (const dziecko of wezel.childNodes) {
    if (dziecko.nodeType !== Node.TEXT_NODE) continue;
    const tresc = dziecko.nodeValue ?? '';
    if (!/\d/u.test(tresc)) continue;
    dziecko.nodeValue = tresc.replace(/\d+/u, String(liczba));
    return;
  }
  wezel.remove();
}

/** Wpisuje wartość w ostatni niepusty węzeł tekstowy wiersza segmentu; numer segmentu stoi w osobnym węźle i zostaje nietknięty. */
function wpiszOstatni(wezel: Element, tekst: string): void {
  const miejsca = [...wezel.querySelectorAll('span')].filter(
    (span) => !span.classList.contains('tr-seg-nr'),
  );
  const miejsce = miejsca[miejsca.length - 1];
  if (miejsce === undefined) wezel.textContent = tekst;
  else miejsce.textContent = tekst;
}

/**
 * Wpisuje wartość w węzeł tekstowy znacznika. Wartość pusta zostawia węzeł
 * pustym, a nie zdejmuje go: wartość dochodzi dopiero z pracą podjętą w oknie,
 * więc znak na nią czeka. Węzeł raz opróżniony pozostaje celem kolejnych
 * wpisów, bo wybór pada na ostatni węzeł tekstowy, gdy niepustego już nie ma.
 */
function opiszWezel(wezel: HTMLElement | null, wartosc: string): void {
  if (wezel === null) return;
  const wezly = [...wezel.childNodes].filter((dziecko) => dziecko.nodeType === Node.TEXT_NODE);
  const cel = wezly.find((dziecko) => (dziecko.nodeValue ?? '').trim() !== '')
    ?? wezly[wezly.length - 1];
  if (cel !== undefined) cel.nodeValue = wartosc;
}

/** Wpisuje wartość w pierwszy niepusty węzeł tekstowy, zostawiając ikonę i przyciski znacznika; fałsz znaczy węzeł bez miejsca na tekst. */
function wpiszTekst(wezel: Element, tekst: string): boolean {
  for (const dziecko of wezel.childNodes) {
    if (dziecko.nodeType !== Node.TEXT_NODE) continue;
    if ((dziecko.nodeValue ?? '').trim() === '') continue;
    dziecko.nodeValue = tekst;
    return true;
  }
  return false;
}

/** Klon węzła wzorcowego, odporny na jego brak w znaczniku. */
function sklonuj(wezel: Element | null): HTMLElement | null {
  return wezel instanceof HTMLElement ? (wezel.cloneNode(true) as HTMLElement) : null;
}
