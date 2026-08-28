import type { RoundtableParticipant, RoundtableStatement } from '../../../../shared/contract';
import { opisMowcy, PARTICIPANT_ID_MODERATORA } from './czytelnosc-glosow';
import {
  ostatniaWypowiedz,
  STAN_GLOSU,
  wyciszony,
  zdanieBezSlow,
  zlozGlosBiezacy,
} from './glos-biezacy';
import type { StanDebaty } from './stan-debaty';
import type { StanTresci } from './stany-okna';
import type { GlosNaZywo, StrumienWypowiedzi } from './strumien-wypowiedzi';

/** Widok panelu debaty — głosy wielu modeli jeden pod drugim: warianty wstęgi przypisywane po kolejności głosu, dalsi dostają wariant wspólny. */
const ROLE_MOWCOW = ['pierwszy', 'drugi', 'trzeci'] as const;

/** Rysuje cały panel: nagłówek tury, głosy uczestników, głosy nieprzypisane, albo stan pusty, gdy panel niczego nie widział. */
export function rysujPanelDebaty(
  stan: StanDebaty,
  strumien: StrumienWypowiedzi,
  tresc: StanTresci,
): void {
  const uczestnicy = stan.uczestnicy();
  const wypowiedzi = stan.wypowiedzi();
  const nieprzypisane = strumien.nieprzypisane();
  const obcy = wypowiedziSpozaSkladu(wypowiedzi, uczestnicy);

  if (uczestnicy.length === 0 && wypowiedzi.length === 0 && nieprzypisane.length === 0) {
    tresc.pusto(zdaniePustego(stan));
    return;
  }

  const widok = tresc.tresc();
  widok.append(naglowekTury(stan, uczestnicy.length));

  const lista = document.createElement('ul');
  lista.className = 'dr-glosy';
  lista.setAttribute('aria-label', 'Głosy uczestników tury bieżącej');
  uczestnicy.forEach((uczestnik, numer) => {
    lista.append(glosUczestnika(uczestnik, numer, stan, strumien, wypowiedzi));
  });
  for (const wypowiedz of obcy) lista.append(glosSpozaSkladu(wypowiedz, stan));
  for (const wpis of nieprzypisane) lista.append(glosNieprzypisany(wpis.identyfikator, wpis.glos));
  widok.append(lista);
}

/**
 * Nagłówek gęsty: tura w jednym wierszu, skład w drugim.
 *
 * Tura bez definicji nie dostaje wiersza udającego turę zerową — dostaje zdanie
 * o tym, że żadnej nie uruchomiono.
 */
function naglowekTury(stan: StanDebaty, ilu: number): HTMLElement {
  const element = document.createElement('div');
  element.className = 'dr-panel-naglowek';

  const tura = document.createElement('span');
  tura.className = 'dr-panel-naglowek__tura';
  const definicja = stan.definicjaTury();
  if (definicja === null) {
    tura.textContent = 'Żadnej tury jeszcze nie uruchomiono w tym oknie.';
  } else {
    tura.dataset['stanTury'] = definicja.status;
    const temat = definicja.topic === undefined || definicja.topic === '' ? '' : ` · ${definicja.topic}`;
    tura.textContent = `Tura ${definicja.index} · ${definicja.format} · ${definicja.status}${temat}`;
  }

  const skladWiersz = document.createElement('span');
  skladWiersz.className = 'dr-panel-naglowek__sklad';
  skladWiersz.textContent =
    ilu === 0
      ? 'Skład debaty nieznany temu panelowi — odczytu roundtable.model.list panel jeszcze nie wywołuje.'
      : `Przy stole: ${ilu}.`;

  element.append(tura, skladWiersz);
  return element;
}

/** Głos jednego uczestnika składu — tożsamość, stan głosu, treść rosnąca albo utrwalona, ze zdaniem o braku słów. */
function glosUczestnika(
  uczestnik: RoundtableParticipant,
  numer: number,
  stan: StanDebaty,
  strumien: StrumienWypowiedzi,
  wypowiedzi: readonly RoundtableStatement[],
): HTMLElement {
  const cisza = wyciszony(uczestnik);
  const biezacy = zlozGlosBiezacy(
    strumien.glos(uczestnik.id),
    ostatniaWypowiedz(wypowiedzi, uczestnik.id),
    cisza,
  );
  const element = pozycja(ROLE_MOWCOW[numer] ?? 'dalszy');
  element.dataset['uczestnik'] = uczestnik.id;
  if (cisza) element.dataset['wyciszony'] = 'tak';

  element.append(mowca(opisMowcy(uczestnik.id, uczestnik, stan)), znacznik(biezacy.stan));
  if (biezacy.tekst === '') {
    element.append(uwaga(zdanieBezSlow(biezacy, cisza)));
  } else {
    element.append(akapit(biezacy.tekst));
  }
  if (biezacy.rozumowanie !== '') element.append(rozumowanie(biezacy.rozumowanie));
  if (biezacy.przyczyna !== '') element.append(uwaga(`Strumień zerwany: ${biezacy.przyczyna}`));
  return element;
}

/**
 * Wypowiedź mówcy spoza składu — moderator albo uczestnik, o którym panel jeszcze nie wie
 * z odczytu składu.
 */
function glosSpozaSkladu(wypowiedz: RoundtableStatement, stan: StanDebaty): HTMLElement {
  const moderator = wypowiedz.participantId === PARTICIPANT_ID_MODERATORA;
  const element = pozycja(moderator ? 'moderator' : 'dalszy');
  element.dataset['uczestnik'] = wypowiedz.participantId;
  element.append(
    mowca(opisMowcy(wypowiedz.participantId, null, stan)),
    znacznik(STAN_GLOSU.utrwalona),
    akapit(wypowiedz.content),
  );
  return element;
}

/**
 * Głos, którego nie dało się przypisać — pokazany, nie porzucony, pod surowym identyfikatorem
 * strumienia.
 */
function glosNieprzypisany(identyfikator: string, glos: GlosNaZywo): HTMLElement {
  const element = pozycja('nieznany');
  element.append(
    mowca(`Mówca nieustalony (identyfikator strumienia: ${identyfikator})`),
    znacznik(glos.domkniety ? 'strumień domknięty' : 'strumień otwarty'),
    uwaga(
      'Ten fragment przyszedł z okna tej debaty, ale jego identyfikatora nie zna ani skład, ' +
        'ani wykaz wypowiedzi tury. Przypisanie nastąpi samo, gdy rdzeń rozgłosi brakującą tożsamość.',
    ),
  );
  if (glos.tekst !== '') element.append(akapit(glos.tekst));
  if (glos.przyczyna !== '') element.append(uwaga(`Strumień zerwany: ${glos.przyczyna}`));
  return element;
}

/** Wypowiedzi, których mówca nie stoi w składzie znanym panelowi — moderator albo uczestnik jeszcze nieznany. */
function wypowiedziSpozaSkladu(
  wypowiedzi: readonly RoundtableStatement[],
  uczestnicy: readonly RoundtableParticipant[],
): RoundtableStatement[] {
  const znani = new Set(uczestnicy.map((uczestnik) => uczestnik.id));
  return wypowiedzi.filter((wypowiedz) => !znani.has(wypowiedz.participantId) && wypowiedz.content !== '');
}

/** Zdanie stanu pustego — rozdziela brak okna od braku debaty, nazywając wprost, czego panelowi brakuje. */
function zdaniePustego(stan: StanDebaty): string {
  if (stan.okno() === '') {
    return (
      'Gospodarz nie podał okna debaty, więc panel nie ma czego słuchać: fragmenty `stream.chunk` ' +
      'jadą wspólną drogą wszystkich kanałów i bez okna nie da się odróżnić głosów tej debaty ' +
      'od strumieni okna rozmowy.'
    );
  }
  return (
    'To okno debaty nie ma jeszcze ani jednego uczestnika i ani jednej wypowiedzi widzianej przez ' +
    'panel. Odczytu składu i przebiegu (roundtable.model.list, roundtable.debate.get) panel jeszcze ' +
    'nie wywołuje, więc pokazuje wyłącznie to, co zdarzenia przyniosły od jego otwarcia.'
  );
}

/** Pozycja wykazu głosów wraz z wariantem wstęgi, wspólnym elementem listy dla trzech rodzajów mówców debaty. */
function pozycja(rola: string): HTMLElement {
  const element = document.createElement('li');
  element.className = 'dr-glos';
  element.dataset['rolaMowcy'] = rola;
  return element;
}

/**
 * Drobne wytwórnie węzłów. Klasy padają napisami stałymi, nie składanymi
 * z członów: kontrola pokrycia klas CSS czyta napisy, a nazwa złożona w czasie
 * działania jest dla niej ślepa.
 */
function mowca(nazwa: string): HTMLElement {
  const element = document.createElement('span');
  element.className = 'dr-glos__mowca';
  element.textContent = nazwa;
  return element;
}

function znacznik(zdanie: string): HTMLElement {
  const element = document.createElement('span');
  element.className = 'dr-glos__znacznik';
  element.dataset['stanGlosu'] = zdanie;
  element.textContent = zdanie;
  return element;
}

function akapit(slowa: string): HTMLElement {
  const element = document.createElement('p');
  element.className = 'dr-glos__tresc';
  element.textContent = slowa;
  return element;
}

function rozumowanie(slowa: string): HTMLElement {
  const element = document.createElement('p');
  element.className = 'dr-glos__rozumowanie';
  element.textContent = slowa;
  return element;
}

function uwaga(zdanie: string): HTMLElement {
  const element = document.createElement('p');
  element.className = 'dr-glos__uwaga';
  element.setAttribute('role', 'note');
  element.textContent = zdanie;
  return element;
}
