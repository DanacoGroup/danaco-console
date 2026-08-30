/**
 * Wiązanie wnętrza okna modułu Roundtable z rdzeniem. Znacznik należy do
 * Właściciela — ten plik wypełnia stojące węzły odpowiedziami rodziny
 * roundtable.*, powiela wzory zdjęte z treści przykładowej, a węzły, dla
 * których kontrakt nie ma pola, zdejmuje.
 */
import {
  Command,
  EventType,
  ModeratorAction,
  RoundtableHandoffTarget,
  type Channel,
  type RoundtableAgreementPoint,
  type RoundtableArgumentNode,
  type RoundtableConsensus,
  type RoundtableDebateSnapshot,
  type RoundtableParticipant,
  type RoundtableStatement,
  type RoundtableTurn,
} from '../../../shared/contract.ts';
import type { Kanal } from '../protokol/kanal.ts';
import { wywolaj } from '../protokol/wywolanie.ts';

/** Panele prototypu, którym rodzina roundtable.* nie oddaje ani jednego pola; schodzą razem z przełącznikami menu sesji. */
const PANELE_BEZ_POKRYCIA = ['panel-plan', 'panel-zadania', 'panel-artefakty', 'panel-subagenci'];

/** Wzory zdjęte z treści przykładowej, zanim wykazy zostaną wyczyszczone; klon zachowuje ikonę i klasy nadane przez bibliotekę. */
interface Wzory {
  panel: HTMLElement | null;
  akapit: HTMLElement | null;
  kluczowy: HTMLElement | null;
  tura: HTMLElement | null;
  wypowiedz: HTMLElement | null;
  spor: HTMLElement | null;
  zdanieOdrebne: HTMLElement | null;
  znakUczestnika: HTMLElement | null;
  wierszKanalu: HTMLElement | null;
}

/** Stan okna debaty: skład, tura bieżąca i widok panelu debaty; wszystkie trzy pochodzą z odpowiedzi rdzenia albo z wyboru Operatora. */
interface Stan {
  uczestnicy: Map<string, RoundtableParticipant>;
  tura: string;
  mapaArgumentow: boolean;
}

/**
 * Wiąże wnętrze okna modułu Roundtable. Brak węzłów wnętrza kończy funkcję bez
 * żadnego działania — okna w ramie nie ma, więc nie ma czego wypełniać.
 */
export function zwiazRoundtable(kanal: Kanal, idOkna: string): void {
  const cialo = document.querySelector<HTMLElement>('.sta-cialo');
  if (cialo === null) return;

  const wzory = zdejmijWzory(cialo);
  zdejmijTrescPrzykladowa(cialo);
  zdejmijSterowanieBezPokrycia(cialo);

  const stan: Stan = { uczestnicy: new Map(), tura: '', mapaArgumentow: false };

  /* Odwołanie stoi wartością, nie deklaracją: deklaracja wchodzi na wierzch
     zasięgu, więc kompilator bierze dla węzła wnętrza typ sprzed rozpoznania
     pustki. */
  /** Powtarza odczyt stanu debaty; wołane przy wejściu, po każdej czynności i po każdym przyroście zgłoszonym zdarzeniem. */
  const odswiez = (): void => {
    void wypelnijStanowisko(kanal, idOkna, cialo, wzory, stan);
  };

  zwiazSklad(kanal, idOkna, cialo, stan, odswiez);
  zwiazTure(kanal, idOkna, cialo, odswiez);
  zwiazModeratora(kanal, idOkna, cialo, odswiez);
  zwiazStanowisko(kanal, idOkna, cialo, odswiez);
  zwiazWidokDebaty(cialo, stan, odswiez);
  kanal.naZdarzenie(EventType.RoundtableDebateChanged, (tresc) => {
    if (tresc.turn.windowId !== idOkna) return;
    odswiez();
  });

  odswiez();
  void wypelnijKanaly(kanal, cialo, wzory);
}

/** Zdejmuje wzory paneli, tur i wierszy, zanim wykazy zostaną wyczyszczone; wzór traci oznaczenia bez pokrycia w kontrakcie. */
function zdejmijWzory(cialo: HTMLElement): Wzory {
  /* Wzór panelu bierze się z drugiego mówcy: opis nagłówka okna sięga po
     pierwszy nagłówek dokumentu, a tym jest nagłówek panelu pierwszego. */
  const panel = sklonuj(cialo.querySelector('#model-panel-b'));
  panel?.removeAttribute('id');
  panel?.querySelector('.sta-kom-stan')?.remove();
  panel?.querySelector('.sta-okno-akcje button:last-of-type')?.remove();
  // Długość wypowiedzi i czas odpowiedzi nie mają pola: wypowiedź niesie treść,
  // akt mowy, pewność i numer redakcji.
  for (const miara of panel?.querySelectorAll('.rt-mp-dol .dn-plakietka') ?? []) miara.remove();
  const strumien = panel?.querySelector('.rt-strumien');
  const akapit = sklonuj(strumien?.querySelector('p') ?? null);
  const kluczowy = sklonuj(strumien?.querySelector('.rt-kluczowy') ?? null);
  strumien?.replaceChildren();
  const tura = sklonuj(cialo.querySelector('#panel-debate .rt-tura'));
  const wypowiedz = sklonuj(tura?.querySelector('.dn-wykaz-modulu-poz') ?? null);
  for (const stojaca of tura?.querySelectorAll('.dn-wykaz-modulu-poz') ?? []) stojaca.remove();
  const wierszKanalu = sklonuj(cialo.querySelector('#pop-model .sta-popover-wiersz'));
  wierszKanalu?.querySelector('small')?.remove();
  wierszKanalu?.querySelector('.pt-mono')?.remove();
  return {
    panel,
    akapit,
    kluczowy,
    tura,
    wypowiedz,
    spor: sklonuj(cialo.querySelector('#panel-debate .dn-alert--ostrzezenie')),
    zdanieOdrebne: sklonuj(cialo.querySelector('#panel-consensus .dn-alert--info')),
    znakUczestnika: sklonuj(cialo.querySelectorAll('.sta-kontekst-akcji .sta-chip')[2] ?? null),
    wierszKanalu,
  };
}

/**
 * Zdejmuje treść przykładową wnętrza. Wykaz debat szyny i historia rozmowy
 * wiszą na rodzinach spoza roundtable.*; zasięg wykonania, znak syntezy
 * i podpis środowiska nie mają pola osiągalnego z tego wiązania.
 */
function zdejmijTrescPrzykladowa(cialo: HTMLElement): void {
  cialo.querySelector('.dn-szyna-modulu-lista')?.replaceChildren();
  for (const znacznik of cialo.querySelectorAll('#menu-filtr .sta-menu-poz[aria-checked]')) {
    znacznik.removeAttribute('aria-checked');
  }
  cialo.querySelector('.sta-kom-historia')?.replaceChildren();
  cialo.querySelector('.sta-kom-monitor')?.remove();
  cialo.querySelector('.rt-modele')?.replaceChildren();
  /* Podpis środowiska zdejmuje się razem z panelami mówców: opis nagłówka
     okna sięga po pierwszy nagłówek dokumentu, a po ich zdjęciu okno
     komunikacji stoi już bez wartości, którą ten opis wnosi. Środowisko
     nazywa pas stanu ramy. */
  zdejmijPoleNaglowka(cialo, 'Środowisko');
  cialo.querySelector('#okno-czat-1 .sta-kom-kontekst > .sta-zrodlo')?.remove();
  // Odpięcie znaku kontekstu nie ma komendy w żadnej z rodzin tego okna.
  for (const odepnij of cialo.querySelectorAll('.sta-kom-kontekst .odepnij')) odepnij.remove();
  cialo.querySelectorAll('.sta-kontekst-akcji .sta-chip')[0]?.remove();
  cialo.querySelector('#panel-consensus .sta-okno-belka .dn-plakietka')?.remove();
  zdejmijCzasTury(cialo);
}

/** Zdejmuje pole nagłówka rozpoznane po podpisie; wartości brakującej nie zastępuje się treścią prototypu. */
function zdejmijPoleNaglowka(cialo: HTMLElement, podpis: string): void {
  const pola = [...cialo.querySelectorAll('#okno-czat-1 .sta-kom-naglowek .sta-kom-pole')];
  pola.find((pole) => pole.textContent?.startsWith(podpis) === true)?.remove();
}

/**
 * Zdejmuje zegar tury wraz z paskiem postępu. Tura niesie granicę czasu
 * i chwilę rozpoczęcia; czasu pozostałego nie niesie ani ona, ani żadna
 * odpowiedź obszaru, więc odliczanie byłoby miarą dorobioną.
 */
function zdejmijCzasTury(cialo: HTMLElement): void {
  const panel = cialo.querySelector<HTMLElement>('#panel-moderator .dn-wykaz-modulu');
  if (panel === null) return;
  panel.querySelector('.dn-postep')?.remove();
  const etykiety = [...panel.querySelectorAll('.pt-etykieta')];
  etykiety[0]?.remove();
  const wiersze = [...panel.querySelectorAll('.dn-wykaz-modulu-poz')];
  wiersze[3]?.remove();
}

/**
 * Zdejmuje sterowanie, którego rodzina roundtable.* nie obsługuje: panele bez
 * pokrycia wraz z ich przełącznikami, urządzenia wejścia dźwięku, nakład
 * rozumowania, menu dodawania oraz wydania żądające formatu, którego okno nie
 * niesie.
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
  cialo.querySelector('#pop-mik')?.closest('.sta-nrz')?.remove();
  cialo.querySelector('#pop-wysilek')?.closest('.sta-nrz')?.remove();
  // Skład debaty wchodzi wyborem kanału, nie pozycją menu bez wykazu kanałów.
  cialo.querySelector('#pop-plus')?.closest('.sta-nrz')?.remove();
  // Podpis piątego trybu opisuje domyślność, której rejestr uprawnień nie zna.
  cialo.querySelector('#pop-tryb-upr .sta-popover-wiersz small')?.remove();
  cialo.querySelector('#panel-debate .sta-okno-akcje button')?.remove();
  zdejmijWydaniaBezFormatu(cialo);
}

/** Zdejmuje przyciski wydania i stylu; zapis debaty, graf i stanowisko wychodzą w formacie, którego okno nie ma czym wskazać. */
function zdejmijWydaniaBezFormatu(cialo: HTMLElement): void {
  const debata = [...cialo.querySelectorAll('#panel-debate .rt-mp-dol button')];
  debata[1]?.remove();
  const stanowisko = [...cialo.querySelectorAll('#panel-consensus .rt-mp-dol button')];
  stanowisko[1]?.remove();
  stanowisko[4]?.remove();
  const moderator = [...cialo.querySelectorAll('#panel-moderator .rt-mp-dol button')];
  moderator[1]?.remove();
}

/** Pyta rdzeń o stan debaty, punkty zgody i graf argumentów, po czym wypełnia nimi panele mówców, panel debaty, moderatora i stanowiska. */
async function wypelnijStanowisko(
  kanal: Kanal,
  idOkna: string,
  cialo: HTMLElement,
  wzory: Wzory,
  stan: Stan,
): Promise<void> {
  const debata = await wywolaj(kanal, Command.RoundtableDebateGet, { windowId: idOkna });
  const zapis = debata.udany ? debata.wynik?.snapshot : undefined;
  const uczestnicy = zapis?.participants ?? [];
  stan.uczestnicy = new Map(uczestnicy.map((uczestnik) => [uczestnik.id, uczestnik]));
  const tura = zapis?.turns[zapis.turns.length - 1];
  stan.tura = tura?.id ?? '';

  const graf = await wywolaj(kanal, Command.RoundtableArgumentList, { windowId: idOkna });
  const wezly = graf.udany ? (graf.wynik?.graph.nodes ?? []) : [];
  const zgoda = await wywolaj(kanal, Command.RoundtableAgreementGet, { windowId: idOkna });
  const punkty = zgoda.udany ? (zgoda.wynik?.points ?? []) : [];

  const kanaly = await wywolaj(kanal, Command.ChannelList, { enabledOnly: true });
  const wykazKanalow = kanaly.udany ? (kanaly.wynik?.channels ?? []) : [];

  wypelnijPanele(cialo, wzory, zapis, wezly, wykazKanalow);
  opiszNaglowekDebaty(cialo, uczestnicy, tura);
  wypelnijZnakiSkladu(cialo, wzory, uczestnicy);
  wypelnijPrzebieg(cialo, wzory, zapis, wezly, punkty, stan);
  wypelnijModeratora(cialo, uczestnicy, tura);
  const kluczowe = wezly.filter((wezel) => wezel.pinned === true).length;
  wypelnijStanowiskoKoncowe(cialo, wzory, zapis?.consensus, punkty, kluczowe);
}

/** Stawia po jednym panelu mówcy na uczestnika debaty; skład pusty zostawia obszar paneli pustym, bo panel bez uczestnika byłby panelem dorobionym. */
function wypelnijPanele(
  cialo: HTMLElement,
  wzory: Wzory,
  zapis: RoundtableDebateSnapshot | undefined,
  wezly: RoundtableArgumentNode[],
  kanaly: Channel[],
): void {
  const obszar = cialo.querySelector<HTMLElement>('.rt-modele');
  const wzor = wzory.panel;
  if (obszar === null || wzor === null) return;
  obszar.replaceChildren();
  const nazwyKanalow = new Map(kanaly.map((kanal) => [kanal.id, kanal.model ?? kanal.name]));
  for (const uczestnik of zapis?.participants ?? []) {
    const panel = wzor.cloneNode(true) as HTMLElement;
    panel.dataset.uczestnik = uczestnik.id;
    opiszMowce(panel, uczestnik, nazwyKanalow.get(uczestnik.channelId) ?? '');
    wypelnijStrumien(panel, wzory, uczestnik, zapis?.statements ?? [], wezly);
    obszar.appendChild(panel);
  }
}

/** Opisuje belkę panelu mówcy tożsamością uczestnika: awatar, nazwa persony, kanał modelu i rola w naradzie; wartość nieobecna węzeł zdejmuje. */
function opiszMowce(panel: HTMLElement, uczestnik: RoundtableParticipant, kanal: string): void {
  const awatar = panel.querySelector<HTMLElement>('.rt-awatar');
  if (awatar !== null) {
    if (uczestnik.avatar === undefined || uczestnik.avatar === '') awatar.remove();
    else awatar.textContent = uczestnik.avatar;
  }
  const nazwa = panel.querySelector('.sta-okno-tytul b');
  if (nazwa !== null) {
    if (uczestnik.personaName === undefined) nazwa.remove();
    else nazwa.textContent = uczestnik.personaName;
  }
  const plakietka = panel.querySelector('.dn-plakietka--rola');
  if (plakietka !== null) {
    if (kanal === '') plakietka.remove();
    else plakietka.textContent = kanal;
  }
  const rola = uczestnik.roleDescription ?? uczestnik.role ?? '';
  const pole = panel.querySelector('.sta-kom-pole b');
  if (pole === null) return;
  if (rola === '') panel.querySelector('.sta-kom-naglowek')?.remove();
  else pole.textContent = rola;
}

/** Wypełnia strumień panelu wypowiedziami uczestnika; węzeł grafu oznaczony jako kluczowy staje osobnym wyróżnieniem pod wypowiedzią, z której go wydobyto. */
function wypelnijStrumien(
  panel: HTMLElement,
  wzory: Wzory,
  uczestnik: RoundtableParticipant,
  wypowiedzi: RoundtableStatement[],
  wezly: RoundtableArgumentNode[],
): void {
  const strumien = panel.querySelector<HTMLElement>('.rt-strumien');
  if (strumien === null) return;
  strumien.replaceChildren();
  for (const wypowiedz of wypowiedzi) {
    if (wypowiedz.participantId !== uczestnik.id) continue;
    if (wzory.akapit !== null) {
      const akapit = wzory.akapit.cloneNode(true) as HTMLElement;
      akapit.textContent = wypowiedz.content;
      strumien.appendChild(akapit);
    }
    if (wzory.kluczowy === null) continue;
    for (const wezel of wezly) {
      if (wezel.statementId !== wypowiedz.id || wezel.pinned !== true) continue;
      const kluczowy = wzory.kluczowy.cloneNode(true) as HTMLElement;
      kluczowy.textContent = wezel.text;
      strumien.appendChild(kluczowy);
    }
  }
}

/** Opisuje nagłówek okna komunikacji formatem i numerem tury oraz liczbą uczestników na belce; brak tury zdejmuje pola tury. */
function opiszNaglowekDebaty(
  cialo: HTMLElement,
  uczestnicy: RoundtableParticipant[],
  tura: RoundtableTurn | undefined,
): void {
  const pola = [...cialo.querySelectorAll<HTMLElement>('#okno-czat-1 .sta-kom-pole')];
  wpiszDane(pola[0] ?? null, tura?.format ?? '');
  const numer = tura === undefined
    ? ''
    : tura.turnLimit === undefined
      ? String(tura.index)
      : `${tura.index} z ${tura.turnLimit}`;
  wpiszDane(pola[1] ?? null, numer);
  const stan = cialo.querySelector<HTMLElement>('#okno-czat-1 .sta-kom-stan .dn-plakietka');
  opiszWezel(stan, tura === undefined ? '' : ` ${tura.status}`);
  wpiszLiczby(
    cialo.querySelector<HTMLElement>('#okno-czat-1 .sta-okno-belka > .sta-chip'),
    [uczestnicy.length],
  );
  opiszWezel(
    cialo.querySelector('#okno-czat-1 .sta-kom-kontekst > .sta-zrodlo'),
    tura?.topic ?? '',
  );
}

/**
 * Wpisuje wartość w wyróżnienie pola nagłówka. Pole zostaje także bez wartości:
 * tura powstaje dopiero w toku debaty, więc pole na nią czeka, a nie jest
 * polem bez pokrycia w kontrakcie.
 */
function wpiszDane(pole: HTMLElement | null, wartosc: string): void {
  if (pole === null) return;
  const dane = pole.querySelector('.dane') ?? pole.querySelector('b');
  if (dane !== null) dane.textContent = wartosc;
}

/** Stawia w pasie czynności liczbę uczestników i po jednym znaku na uczestnika; skład pusty zostawia sam licznik. */
function wypelnijZnakiSkladu(
  cialo: HTMLElement,
  wzory: Wzory,
  uczestnicy: RoundtableParticipant[],
): void {
  const pas = cialo.querySelector<HTMLElement>('.sta-kontekst-akcji');
  if (pas === null) return;
  wpiszLiczby(pas.querySelector<HTMLElement>('.sta-chip'), [uczestnicy.length]);
  for (const stojacy of [...pas.querySelectorAll('.sta-chip')].slice(1)) stojacy.remove();
  if (wzory.znakUczestnika === null) return;
  const przycisk = pas.querySelector('button');
  for (const uczestnik of uczestnicy) {
    if (uczestnik.personaName === undefined) continue;
    const znak = wzory.znakUczestnika.cloneNode(true) as HTMLElement;
    znak.textContent = uczestnik.personaName;
    pas.insertBefore(znak, przycisk);
  }
}

/** Wypełnia panel debaty: chronologia stawia turę po turze, a mapa argumentów węzły grafu; oba widoki biorą treść z tej samej odpowiedzi rdzenia. */
function wypelnijPrzebieg(
  cialo: HTMLElement,
  wzory: Wzory,
  zapis: RoundtableDebateSnapshot | undefined,
  wezly: RoundtableArgumentNode[],
  punkty: RoundtableAgreementPoint[],
  stan: Stan,
): void {
  const mapa = cialo.querySelector<HTMLElement>('#panel-debate .rt-mapa');
  const stopka = mapa?.querySelector('.rt-mp-dol') ?? null;
  const wzorTury = wzory.tura;
  if (mapa === null || wzorTury === null) return;
  mapa.replaceChildren();
  if (stan.mapaArgumentow) wypelnijMape(mapa, wzorTury, wzory, wezly, stan);
  else wypelnijChronologie(mapa, wzorTury, wzory, zapis, punkty, stan);
  if (stopka !== null) mapa.appendChild(stopka);
  const znacznik = cialo.querySelector<HTMLElement>('#panel-debate .dn-plakietka--informacja');
  const tura = zapis?.turns[zapis.turns.length - 1];
  /* Znacznik nazywa turę i granicę tur; bez tury nazywałby turę, której nie
     ma, więc schodzi zamiast stanąć z zerem. */
  if (tura === undefined || tura.turnLimit === undefined) znacznik?.remove();
  else wpiszLiczby(znacznik, [tura.index, tura.turnLimit]);
}

/** Stawia w panelu debaty blok na każdą turę wraz z jej wypowiedziami i punktami spornymi tej tury. */
function wypelnijChronologie(
  mapa: HTMLElement,
  wzorTury: HTMLElement,
  wzory: Wzory,
  zapis: RoundtableDebateSnapshot | undefined,
  punkty: RoundtableAgreementPoint[],
  stan: Stan,
): void {
  for (const tura of zapis?.turns ?? []) {
    const blok = wzorTury.cloneNode(true) as HTMLElement;
    const etykieta = blok.querySelector('.pt-etykieta');
    if (etykieta !== null) etykieta.textContent = String(tura.index);
    for (const wypowiedz of zapis?.statements ?? []) {
      if (wypowiedz.turnId !== tura.id || wzory.wypowiedz === null) continue;
      blok.appendChild(zbudujWypowiedz(wzory.wypowiedz, wypowiedz, stan));
    }
    for (const punkt of punkty) {
      if (punkt.agreed || punkt.turnId !== tura.id || wzory.spor === null) continue;
      const spor = wzory.spor.cloneNode(true) as HTMLElement;
      spor.textContent = punkt.text;
      blok.appendChild(spor);
    }
    mapa.appendChild(blok);
  }
}

/** Stawia w panelu debaty jeden blok z węzłami grafu argumentów; kluczowy punkt sporny wchodzi w etykietę bloku. */
function wypelnijMape(
  mapa: HTMLElement,
  wzorTury: HTMLElement,
  wzory: Wzory,
  wezly: RoundtableArgumentNode[],
  stan: Stan,
): void {
  const blok = wzorTury.cloneNode(true) as HTMLElement;
  blok.querySelector('.pt-etykieta')?.remove();
  for (const wezel of wezly) {
    if (wzory.wypowiedz === null) continue;
    const uczestnik = wezel.participantId === undefined
      ? undefined
      : stan.uczestnicy.get(wezel.participantId);
    const wiersz = wzory.wypowiedz.cloneNode(true) as HTMLElement;
    opiszAwatar(wiersz, uczestnik);
    const akt = wiersz.querySelector('b');
    if (akt !== null) akt.textContent = wezel.speechAct;
    wpiszTeksty(wiersz, [uczestnik?.personaName ?? '', wezel.text]);
    blok.appendChild(wiersz);
  }
  mapa.appendChild(blok);
}

/** Zwraca wiersz przebiegu opisany wypowiedzią: awatar i nazwa uczestnika, akt mowy oraz treść wypowiedzi. */
function zbudujWypowiedz(
  wzor: HTMLElement,
  wypowiedz: RoundtableStatement,
  stan: Stan,
): HTMLElement {
  const wiersz = wzor.cloneNode(true) as HTMLElement;
  const uczestnik = stan.uczestnicy.get(wypowiedz.participantId);
  opiszAwatar(wiersz, uczestnik);
  const akt = wiersz.querySelector('b');
  if (akt !== null) {
    if (wypowiedz.speechAct === undefined) akt.remove();
    else akt.textContent = wypowiedz.speechAct;
  }
  wpiszTeksty(wiersz, [uczestnik?.personaName ?? '', wypowiedz.content]);
  return wiersz;
}

/** Nanosi na wiersz awatar uczestnika; uczestnik bez awatara znaku nie dostaje. */
function opiszAwatar(wiersz: HTMLElement, uczestnik: RoundtableParticipant | undefined): void {
  const awatar = wiersz.querySelector<HTMLElement>('.rt-awatar');
  if (awatar === null) return;
  if (uczestnik?.avatar === undefined || uczestnik.avatar === '') awatar.remove();
  else awatar.textContent = uczestnik.avatar;
}

/** Wypełnia panel moderatora formatem debaty, granicą liczby tur i kolejnością głosu; wartość nieobecna wiersz zdejmuje. */
function wypelnijModeratora(
  cialo: HTMLElement,
  uczestnicy: RoundtableParticipant[],
  tura: RoundtableTurn | undefined,
): void {
  const wiersze = [...cialo.querySelectorAll<HTMLElement>('#panel-moderator .dn-wykaz-modulu-poz')];
  wpiszWartoscWiersza(wiersze[0] ?? null, '.sta-chip', tura?.format ?? '');
  wpiszWartoscWiersza(
    wiersze[1] ?? null,
    '.sta-chip',
    tura?.turnLimit === undefined ? '' : String(tura.turnLimit),
  );
  const kolejnosc = [...uczestnicy]
    .sort((a, b) => (a.order ?? 0) - (b.order ?? 0))
    .map((uczestnik) => uczestnik.personaName ?? '')
    .filter((nazwa) => nazwa !== '')
    .join(' → ');
  wpiszWartoscWiersza(wiersze[2] ?? null, '.dn-meta', kolejnosc);
}

/** Wpisuje wartość we wskazany węzeł wiersza wykazu; wiersz bez miejsca na wartość schodzi, wiersz bez wartości na nią czeka. */
function wpiszWartoscWiersza(wiersz: HTMLElement | null, wybor: string, wartosc: string): void {
  if (wiersz === null) return;
  const miejsce = wiersz.querySelector(wybor);
  if (miejsce === null) wiersz.remove();
  else miejsce.textContent = wartosc;
}

/** Wypełnia panel stanowiska treścią stanowiska końcowego, liczbą punktów zgody i sporu oraz zdaniami odrębnymi uczestników. */
function wypelnijStanowiskoKoncowe(
  cialo: HTMLElement,
  wzory: Wzory,
  stanowisko: RoundtableConsensus | undefined,
  punkty: RoundtableAgreementPoint[],
  kluczowe: number,
): void {
  const panel = cialo.querySelector<HTMLElement>('#panel-consensus .dn-wykaz-modulu');
  if (panel === null) return;
  const zrodlo = panel.querySelector<HTMLElement>('.dn-meta');
  /* Miara mówi, z ilu tur stanowisko złożono; bez stanowiska mówiłaby o
     złożeniu, którego nie było, więc schodzi zamiast stanąć z zerem. */
  if (stanowisko === undefined) zrodlo?.remove();
  else wpiszLiczby(zrodlo, [stanowisko.turnIds?.length ?? 0, kluczowe]);
  const tresc = panel.querySelector<HTMLElement>('.rt-kluczowy');
  if (tresc !== null) {
    tresc.textContent = stanowisko?.content ?? '';
    tresc.setAttribute('contenteditable', 'true');
  }
  const plakietki = [...panel.querySelectorAll<HTMLElement>('.dn-plakietka')];
  wpiszLiczby(plakietki[0] ?? null, [punkty.filter((punkt) => punkt.agreed).length]);
  wpiszLiczby(plakietki[1] ?? null, [punkty.filter((punkt) => !punkt.agreed).length]);
  wypelnijZdaniaOdrebne(panel, wzory, stanowisko);
}

/** Stawia po jednym wyróżnieniu na zdanie odrębne podpisane pod stanowiskiem; brak zdań odrębnych wyróżnienie zdejmuje. */
function wypelnijZdaniaOdrebne(
  panel: HTMLElement,
  wzory: Wzory,
  stanowisko: RoundtableConsensus | undefined,
): void {
  for (const stojace of panel.querySelectorAll('.dn-alert--info')) stojace.remove();
  const zdania = stanowisko?.minority ?? [];
  if (wzory.zdanieOdrebne === null || zdania.length === 0) return;
  const stopka = panel.querySelectorAll('.rt-mp-dol')[1] ?? null;
  for (const zdanie of zdania) {
    const wyroznienie = wzory.zdanieOdrebne.cloneNode(true) as HTMLElement;
    wyroznienie.textContent = zdanie.content;
    if (stopka === null) panel.appendChild(wyroznienie);
    else panel.insertBefore(wyroznienie, stopka);
  }
}

/** Wypełnia wykaz modeli kanałami rejestru rdzenia; wskazanie kanału wnosi uczestnika do debaty tego okna. */
async function wypelnijKanaly(kanal: Kanal, cialo: HTMLElement, wzory: Wzory): Promise<void> {
  const menu = cialo.querySelector<HTMLElement>('#pop-model');
  if (menu === null) return;
  const wynik = await wywolaj(kanal, Command.ChannelList, { enabledOnly: true });
  const kanaly = wynik.udany ? (wynik.wynik?.channels ?? []) : [];
  for (const stojacy of menu.querySelectorAll('.sta-popover-wiersz')) stojacy.remove();
  const chip = cialo.querySelector<HTMLElement>('.sta-chip--model');
  const pierwszy = kanaly[0];
  if (chip !== null) {
    if (pierwszy === undefined) chip.closest('.sta-nrz')?.remove();
    else wpiszTekst(chip, pierwszy.model ?? pierwszy.name);
  }
  if (wzory.wierszKanalu === null) return;
  for (const pozycja of kanaly) {
    const wiersz = wzory.wierszKanalu.cloneNode(true) as HTMLElement;
    wiersz.dataset.kanal = pozycja.id;
    const etykieta = wiersz.querySelector('.sta-popover-etykieta');
    if (etykieta !== null) etykieta.textContent = pozycja.model ?? pozycja.name;
    menu.appendChild(wiersz);
  }
}

/** Wiąże skład debaty: wskazanie kanału wnosi uczestnika, gwiazdka oznacza go jako kluczowego, a wyciszenie idzie czynnością moderatora. */
function zwiazSklad(
  kanal: Kanal,
  idOkna: string,
  cialo: HTMLElement,
  stan: Stan,
  odswiez: () => void,
): void {
  cialo.querySelector('#pop-model')?.addEventListener('click', (zdarzenie) => {
    const cel = zdarzenie.target;
    if (!(cel instanceof Element)) return;
    const kanalModelu = cel.closest<HTMLElement>('.sta-popover-wiersz')?.dataset.kanal;
    if (kanalModelu === undefined) return;
    void dodajUczestnika(kanal, idOkna, kanalModelu, odswiez);
  });
  cialo.querySelector('.rt-modele')?.addEventListener('click', (zdarzenie) => {
    const cel = zdarzenie.target;
    if (!(cel instanceof Element)) return;
    const uczestnik = cel.closest<HTMLElement>('[data-uczestnik]')?.dataset.uczestnik;
    const przycisk = cel.closest<HTMLElement>('button');
    if (uczestnik === undefined || przycisk === null) return;
    if (przycisk.closest('.sta-okno-akcje') !== null) {
      void oznaczKluczowego(kanal, idOkna, stan, uczestnik, odswiez);
      return;
    }
    void przestawGlos(kanal, idOkna, stan, uczestnik, odswiez);
  });
}

/** Wnosi do debaty uczestnika na wskazanym kanale modelu i odświeża skład, gdy rdzeń wniesienie przyjął. */
async function dodajUczestnika(
  kanal: Kanal,
  idOkna: string,
  idKanalu: string,
  odswiez: () => void,
): Promise<void> {
  const wynik = await wywolaj(kanal, Command.RoundtableModelAdd, {
    windowId: idOkna,
    channelId: idKanalu,
  });
  if (!wynik.udany) return;
  odswiez();
}

/** Przestawia oznaczenie uczestnika kluczowego na przeciwne i odświeża skład. */
async function oznaczKluczowego(
  kanal: Kanal,
  idOkna: string,
  stan: Stan,
  idUczestnika: string,
  odswiez: () => void,
): Promise<void> {
  const wynik = await wywolaj(kanal, Command.RoundtableModelUpdate, {
    windowId: idOkna,
    participantId: idUczestnika,
    key: stan.uczestnicy.get(idUczestnika)?.key !== true,
  });
  if (!wynik.udany) return;
  odswiez();
}

/** Wycisza uczestnika w turze albo przywraca mu głos, zależnie od stanu zapisanego w składzie. */
async function przestawGlos(
  kanal: Kanal,
  idOkna: string,
  stan: Stan,
  idUczestnika: string,
  odswiez: () => void,
): Promise<void> {
  const wyciszony = stan.uczestnicy.get(idUczestnika)?.muted === true;
  const wynik = await wywolaj(kanal, Command.RoundtableModeratorDirect, {
    windowId: idOkna,
    action: wyciszony ? ModeratorAction.Unmute : ModeratorAction.Mute,
    participantId: idUczestnika,
  });
  if (!wynik.udany) return;
  odswiez();
}

/**
 * Wiąże uruchomienie tury. Pytanie bierze się z pola polecenia okna
 * komunikacji — to ono jest w tym oknie miejscem, w które Operator wpisuje
 * pytanie kierowane do całego składu.
 */
function zwiazTure(kanal: Kanal, idOkna: string, cialo: HTMLElement, odswiez: () => void): void {
  const pole = cialo.querySelector<HTMLTextAreaElement>('.sta-prompt-obszar');
  if (pole === null) return;
  /** Uruchamia turę pytaniem z pola polecenia; pole puste tury nie uruchamia, bo tura bez pytania nie ma czego rozstrzygać. */
  const uruchom = (): void => {
    const pytanie = pole.value.trim();
    if (pytanie === '') return;
    pole.value = '';
    void uruchomTure(kanal, idOkna, pytanie, odswiez);
  };
  cialo.querySelector('.sta-prompt')?.addEventListener('submit', (zdarzenie) => {
    zdarzenie.preventDefault();
    uruchom();
  });
  for (const przycisk of cialo.querySelectorAll('.sta-prompt-wyslij')) {
    przycisk.addEventListener('click', uruchom);
  }
  cialo.querySelector('.sta-kontekst-akcji button')?.addEventListener('click', uruchom);
  cialo.querySelector('#panel-debate .rt-mp-dol button')?.addEventListener('click', uruchom);
}

/** Uruchamia turę debaty pytaniem Operatora i odświeża stanowisko, gdy rdzeń turę założył. */
async function uruchomTure(
  kanal: Kanal,
  idOkna: string,
  pytanie: string,
  odswiez: () => void,
): Promise<void> {
  const wynik = await wywolaj(kanal, Command.RoundtableDebateStart, {
    windowId: idOkna,
    question: pytanie,
  });
  if (!wynik.udany) return;
  odswiez();
}

/** Wiąże czynności moderatora: zamknięcie tury, ukierunkowanie dyskusji treścią pola interwencji i otwarcie wątku pobocznego. */
function zwiazModeratora(
  kanal: Kanal,
  idOkna: string,
  cialo: HTMLElement,
  odswiez: () => void,
): void {
  const panel = cialo.querySelector<HTMLElement>('#panel-moderator');
  if (panel === null) return;
  panel.querySelector('button.dn-btn--zarys')?.addEventListener('click', () => {
    void czynnoscModeratora(kanal, idOkna, ModeratorAction.CloseTurn, undefined, odswiez);
  });
  const interwencja = panel.querySelector<HTMLInputElement>('.dn-szukaj input');
  interwencja?.addEventListener('keydown', (zdarzenie) => {
    if (zdarzenie.key !== 'Enter') return;
    zdarzenie.preventDefault();
    const tresc = interwencja.value.trim();
    if (tresc === '') return;
    interwencja.value = '';
    void czynnoscModeratora(kanal, idOkna, ModeratorAction.Direct, tresc, odswiez);
  });
  panel.querySelector('.rt-mp-dol button')?.addEventListener('click', () => {
    void czynnoscModeratora(kanal, idOkna, ModeratorAction.OpenSideThread, undefined, odswiez);
  });
}

/** Wykonuje czynność moderatora i odświeża stanowisko, gdy rdzeń ją przyjął. */
async function czynnoscModeratora(
  kanal: Kanal,
  idOkna: string,
  czynnosc: ModeratorAction,
  wiadomosc: string | undefined,
  odswiez: () => void,
): Promise<void> {
  const wynik = await wywolaj(kanal, Command.RoundtableModeratorDirect, {
    windowId: idOkna,
    action: czynnosc,
    message: wiadomosc,
  });
  if (!wynik.udany) return;
  odswiez();
}

/** Wiąże panel stanowiska: korekta zapisuje treść z pola stanowiska, a przekazanie oddaje je modułowi docelowemu. */
function zwiazStanowisko(
  kanal: Kanal,
  idOkna: string,
  cialo: HTMLElement,
  odswiez: () => void,
): void {
  const panel = cialo.querySelector<HTMLElement>('#panel-consensus');
  if (panel === null) return;
  const tresc = panel.querySelector<HTMLElement>('.rt-kluczowy');
  panel.querySelector('.rt-mp-dol button')?.addEventListener('click', () => {
    const zapis = (tresc?.textContent ?? '').trim();
    if (zapis === '') return;
    void zapiszStanowisko(kanal, idOkna, zapis, odswiez);
  });
  const przekazania = [...panel.querySelectorAll<HTMLElement>('.rt-mp-dol')][1];
  przekazania?.addEventListener('click', (zdarzenie) => {
    const cel = zdarzenie.target;
    if (!(cel instanceof Element)) return;
    const przycisk = cel.closest('button');
    if (przycisk === null) return;
    void przekazStanowisko(kanal, idOkna, przycisk.textContent ?? '');
  });
}

/** Zapisuje treść stanowiska końcowego i odświeża panel, gdy rdzeń zapis przyjął. */
async function zapiszStanowisko(
  kanal: Kanal,
  idOkna: string,
  tresc: string,
  odswiez: () => void,
): Promise<void> {
  const wynik = await wywolaj(kanal, Command.RoundtableConsensusSet, {
    windowId: idOkna,
    content: tresc,
  });
  if (!wynik.udany) return;
  odswiez();
}

/**
 * Przekazuje stanowisko końcowe modułowi docelowemu. Nazwa modułu bierze się
 * z podpisu przycisku i musi trafić w rejestr celów kontraktu; podpis spoza
 * rejestru przekazania nie wysyła.
 */
async function przekazStanowisko(
  kanal: Kanal,
  idOkna: string,
  podpis: string,
): Promise<void> {
  const stanowisko = await wywolaj(kanal, Command.RoundtableConsensusGet, { windowId: idOkna });
  const cel = celPrzekazania(podpis);
  if (cel === undefined || !stanowisko.udany || stanowisko.wynik === undefined) return;
  await wywolaj(kanal, Command.RoundtableConsensusHandoff, {
    windowId: idOkna,
    consensusId: stanowisko.wynik.consensus.id,
    target: cel,
  });
}

/** Cel przekazania rozpoznany po podpisie przycisku; podpis spoza rejestru celów kontraktu nie daje żadnego. */
function celPrzekazania(podpis: string): RoundtableHandoffTarget | undefined {
  const nazwa = podpis.replace('→', '').trim().toLowerCase();
  const cele: RoundtableHandoffTarget[] = Object.values(RoundtableHandoffTarget);
  return cele.find((cel) => cel === nazwa);
}

/** Wiąże zakładki panelu debaty: chronologia i mapa argumentów biorą treść z tej samej odpowiedzi, więc przełączenie odświeża panel. */
function zwiazWidokDebaty(cialo: HTMLElement, stan: Stan, odswiez: () => void): void {
  const zakladki = [...cialo.querySelectorAll<HTMLElement>('#panel-debate .dn-zakladka')];
  for (const [numer, zakladka] of zakladki.entries()) {
    zakladka.addEventListener('click', () => {
      for (const inna of zakladki) inna.setAttribute('aria-selected', 'false');
      zakladka.setAttribute('aria-selected', 'true');
      stan.mapaArgumentow = numer === 1;
      odswiez();
    });
  }
}

/**
 * Podmienia kolejne ciągi cyfr w treści węzła kolejnymi liczbami odpowiedzi,
 * zostawiając podpis znacznika. Węzeł, który cyfr nie niesie albo niesie ich
 * więcej niż jest liczb, schodzi — podpis z liczbą przykładową kłamałby.
 */
function wpiszLiczby(wezel: HTMLElement | null, liczby: number[]): void {
  if (wezel === null) return;
  let numer = 0;
  for (const dziecko of wezel.childNodes) {
    if (dziecko.nodeType !== Node.TEXT_NODE) continue;
    dziecko.nodeValue = (dziecko.nodeValue ?? '').replace(/\d+/gu, () => {
      const liczba = liczby[numer];
      numer += 1;
      return liczba === undefined ? '' : String(liczba);
    });
  }
  if (numer !== liczby.length) wezel.remove();
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

/** Wpisuje kolejne wartości w kolejne niepuste węzły tekstowe znacznika; węzeł bez wartości traci treść, bo zostawiony niósłby treść przykładową. */
function wpiszTeksty(wezel: Element, teksty: string[]): void {
  let numer = 0;
  for (const dziecko of wezel.childNodes) {
    if (dziecko.nodeType !== Node.TEXT_NODE) continue;
    if ((dziecko.nodeValue ?? '').trim() === '') continue;
    dziecko.nodeValue = teksty[numer] ?? '';
    numer += 1;
  }
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
