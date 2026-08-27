import { utworzRameOkna } from '../../komponenty/rama-okna';
import {
  pobierzPlik,
  przyciskAkcji as przycisk,
  przyciskBezKomendy,
  wiersz,
  wybor,
} from '../../modele/kontrolki-formularza-braki';
import { Command } from '../../../../shared/contract';
import { powodBrakuObslugi } from './braki-kontraktu';
import { czynnosciAnalizy } from './czynnosci-arsenalu';
import { utworzWykazFunkcji, utworzZestawAkcji } from './warstwy-modulu';
import type { ZrodloArsenaluRoundtable } from './zrodlo-arsenalu';
import type { StanDebaty } from './stan-debaty';
import { utworzStanTresci } from './stany-okna';
import type { StrumienWypowiedzi } from './strumien-wypowiedzi';
import {
  strukturaMarkdown,
  zdaniePolRelacji,
  zlozStrukture,
  type StrukturaDebaty,
} from './struktura-argumentow';

/**
 * Argument Map & Analysis — okno monitora otwierane jako rozszerzenie boczne, liczące
 * chronologię i zbieżność wypowiedzi, nie graf argumentów.
 */
export interface OknoArgumentMap {
  element: HTMLElement;
  odswiez(): void;
  /** Przerysowanie po fragmencie strumienia — wołane przez złożenie modułu. */
  odswiezGlosy(): void;
  /** Zamyka nasłuch `stan.naZmiane(...)`. */
  zamknij(): void;
}

/** Widok analizy wybierany przełącznikiem warstwy drugiej: struktura chronologiczna, zbieżność leksykalna albo macierz par mówców. */
type WidokAnalizy = 'struktura' | 'zbieznosc' | 'macierz';

const OPISY_WIDOKU: ReadonlyArray<readonly [string, string]> = [
  ['struktura', 'Struktura chronologiczna tury'],
  ['zbieznosc', 'Zbieżność leksykalna wypowiedzi'],
  ['macierz', 'Macierz zbieżności par mówców'],
];

export function utworzOknoArgumentMap(
  arsenal: ZrodloArsenaluRoundtable,
  stan: StanDebaty,
  strumien: StrumienWypowiedzi,
): OknoArgumentMap {
  const rama = utworzRameOkna({
    tytul: 'Argument Map & Analysis',
    rola: 'monitor',
    kod: 'argument-map-analysis',
    przeznaczenie:
      'Struktura zapisu debaty i zbieżność wypowiedzi tury bieżącej — chronologia mierzona, nie graf argumentów.',
    przedrostek: 'dr',
  });
  const tresc = utworzStanTresci();
  // Czynności analityczne wołają komendy obszaru wprost; rodzaj analizy bierze się z przełączników okna.
  const powierzchnia = zlozPowierzchnieAnalizy(
    rama,
    tresc.element,
    czynnosciAnalizy(
      arsenal,
      stan,
      (zdanie, powodzenie) => tresc.potwierdzenie(zdanie, powodzenie),
      () => 'argumentMining',
      () => 'svg',
    ),
  );

  function rysuj(): void {
    const struktura = zlozStrukture(stan, strumien);
    if (struktura.wezly.length === 0) {
      tresc.pusto(
        'Żadna wypowiedź tej tury nie dotarła jeszcze do tego okna — nie ma czego analizować. ' +
          '(Historię sprzed otwarcia okna oddaje komenda roundtable.debate.get, ale okna modułu jeszcze jej nie wywołują.)',
      );
      rama.ustawZnacznik('bez wypowiedzi', 'neutralna');
      return;
    }
    rama.ustawZnacznik(`wypowiedzi: ${struktura.wezly.length}`, 'neutralna');
    const miejsce = tresc.tresc();
    miejsce.append(naglowekStruktury(stan, struktura));
    const widok = powierzchnia.widok.value as WidokAnalizy;
    if (widok === 'struktura') miejsce.append(widokStruktury(struktura));
    else if (widok === 'zbieznosc') miejsce.append(widokZbieznosci(struktura));
    else miejsce.append(widokMacierzy(struktura));
  }

  function odswiezGlosy(): void {
    const rodzaj = tresc.rodzaj();
    if (rodzaj === 'blad' || rodzaj === 'ladowanie') return;
    rysuj();
  }

  function eksportuj(): void {
    const struktura = zlozStrukture(stan, strumien);
    if (struktura.wezly.length === 0) {
      tresc.potwierdzenie('Nie ma czego wyeksportować — okno nie widziało jeszcze wypowiedzi.', false);
      return;
    }
    pobierzPlik(
      `struktura-debaty-${stan.tura() === '' ? 'bez-tury' : stan.tura()}.md`,
      strukturaMarkdown(struktura, opisTury(stan)),
      'text/markdown',
    );
    tresc.potwierdzenie(
      'Struktura pobrana jako plik Markdown. Obejmuje chronologię, zbieżność leksykalną i pomiar pokrycia pól relacji — grafu argumentów w pliku nie ma, bo okno jeszcze go nie buduje.',
      true,
    );
  }

  powierzchnia.eksport.addEventListener('click', eksportuj);
  powierzchnia.widok.addEventListener('change', rysuj);

  const odsubskrybuj = stan.naZmiane(rysuj);
  rysuj();

  return { element: rama.element, odswiez: rysuj, odswiezGlosy, zamknij: odsubskrybuj };
}

/** Opis tury bieżącej dla nagłówka analizy i eksportu, złożony z numeru, stanu i zagadnienia tury zapisanej w kontrakcie. */
function opisTury(stan: StanDebaty): string {
  const definicja = stan.definicjaTury();
  if (definicja === null) return 'tury nieznanej temu oknu';
  const temat = definicja.topic === undefined || definicja.topic === '' ? '' : ` o zagadnieniu „${definicja.topic}”`;
  return `tury #${definicja.index} (${definicja.status})${temat}`;
}

/** Nagłówek analizy: liczba wypowiedzi, mówców i znaków tury bieżącej, wraz ze zdaniem o polach relacji, których kontrakt nie wypełnia. */
function naglowekStruktury(stan: StanDebaty, struktura: StrukturaDebaty): HTMLElement {
  const blok = document.createElement('div');
  blok.className = 'dr-analiza__naglowek';

  const zakres = document.createElement('p');
  zakres.className = 'dr-analiza__zdanie';
  zakres.textContent =
    `Analiza ${opisTury(stan)}: ${struktura.wezly.length} wypowiedzi, ${struktura.mowcow} mówców, ` +
    `${struktura.znakow} znaków treści.`;

  const granica = document.createElement('p');
  granica.className = 'dr-analiza__zdanie';
  granica.dataset['granica'] = 'tak';
  granica.textContent = zdaniePolRelacji(struktura);

  blok.append(zakres, granica);
  return blok;
}

/** Widok chronologiczny — węzeł na wypowiedź, w kolejności przyjścia, z kolejnością, mówcą, chwilą i długością treści. */
function widokStruktury(struktura: StrukturaDebaty): HTMLElement {
  const lista = document.createElement('ul');
  lista.className = 'dr-wezly';
  lista.setAttribute('aria-label', 'Struktura chronologiczna tury');
  for (const wezel of struktura.wezly) {
    const pozycja = document.createElement('li');
    pozycja.className = 'dr-wezel';
    pozycja.dataset['mowca'] = wezel.idMowcy;

    const znacznik = document.createElement('span');
    znacznik.className = 'dr-wezel__znacznik';
    znacznik.textContent = wezel.znacznik;

    const glowa = document.createElement('span');
    glowa.className = 'dr-wezel__mowca';
    glowa.textContent = `${wezel.kolejnosc}. ${wezel.mowca}`;

    const meta = document.createElement('span');
    meta.className = 'dr-wezel__meta';
    meta.textContent = `${wezel.chwila} · ${wezel.znakow} znaków`;

    pozycja.append(znacznik, glowa, meta);
    lista.append(pozycja);
  }
  return lista;
}

/** Widok zbieżności — zdania powtórzone dosłownie u kilku mówców, po sprowadzeniu do jednej postaci leksykalnej, oraz zdania własne. */
function widokZbieznosci(struktura: StrukturaDebaty): HTMLElement {
  const blok = document.createElement('div');
  blok.className = 'dr-analiza__blok';

  const miara = document.createElement('p');
  miara.className = 'dr-analiza__zdanie';
  miara.dataset['granica'] = 'tak';
  miara.textContent =
    'Miara jest leksykalna: porównywane są zdania sprowadzone do jednej postaci (małe litery, bez znaków przestankowych). ' +
    'Ta sama teza wypowiedziana innymi słowami liczy się tu jako dwa różne zdania, a zdanie przytoczone po to, by je obalić — jako zbieżność. ' +
    `Zdań zbyt krótkich do porównania pominięto: ${struktura.pominietychZdan}.`;
  blok.append(miara);

  const wspolne = document.createElement('ul');
  wspolne.className = 'dr-wezly';
  wspolne.setAttribute('aria-label', 'Zdania powtórzone u kilku mówców');
  if (struktura.wspolne.length === 0) {
    const pusta = document.createElement('li');
    pusta.className = 'dr-wezel';
    pusta.textContent = 'Żadne zdanie nie powtórzyło się dosłownie u więcej niż jednego mówcy.';
    wspolne.append(pusta);
  } else {
    for (const zdanie of struktura.wspolne) {
      const pozycja = document.createElement('li');
      pozycja.className = 'dr-wezel';
      const tresc = document.createElement('span');
      tresc.className = 'dr-wezel__mowca';
      tresc.textContent = zdanie.tresc;
      const mowcy = document.createElement('span');
      mowcy.className = 'dr-wezel__meta';
      mowcy.textContent = `Mówcy: ${zdanie.mowcy.join(' · ')}`;
      pozycja.append(tresc, mowcy);
      wspolne.append(pozycja);
    }
  }
  blok.append(wspolne, wykazWlasnych(struktura));
  return blok;
}

/** Ile zdań każdy mówca powiedział wyłącznie sam, licząc tylko zdania, które przeszły próg długości porównania. */
function wykazWlasnych(struktura: StrukturaDebaty): HTMLElement {
  const lista = document.createElement('ul');
  lista.className = 'dr-wezly';
  lista.setAttribute('aria-label', 'Zdania wyłącznie własne');
  if (struktura.wylacznieWlasne.size === 0) {
    const pusta = document.createElement('li');
    pusta.className = 'dr-wezel';
    pusta.textContent = 'Żadne zdanie nie przeszło progu długości porównania.';
    lista.append(pusta);
    return lista;
  }
  for (const [mowca, ile] of struktura.wylacznieWlasne) {
    const pozycja = document.createElement('li');
    pozycja.className = 'dr-wezel';
    const nazwa = document.createElement('span');
    nazwa.className = 'dr-wezel__mowca';
    nazwa.textContent = mowca;
    const liczba = document.createElement('span');
    liczba.className = 'dr-wezel__meta';
    liczba.textContent = `zdań wyłącznie własnych: ${ile}`;
    pozycja.append(nazwa, liczba);
    lista.append(pozycja);
  }
  return lista;
}

/** Widok macierzy — pary mówców i liczba zdań wspólnych między nimi, licząc wyłącznie zdania powtórzone dosłownie. */
function widokMacierzy(struktura: StrukturaDebaty): HTMLElement {
  const lista = document.createElement('ul');
  lista.className = 'dr-wezly';
  lista.setAttribute('aria-label', 'Macierz zbieżności par mówców');
  if (struktura.pary.length === 0) {
    const pusta = document.createElement('li');
    pusta.className = 'dr-wezel';
    pusta.textContent =
      'Żadna para mówców nie ma zdania wspólnego. Nie znaczy to sporu — znaczy tyle, że mówcy nie użyli tych samych zdań.';
    lista.append(pusta);
    return lista;
  }
  for (const para of struktura.pary) {
    const pozycja = document.createElement('li');
    pozycja.className = 'dr-wezel';
    const nazwa = document.createElement('span');
    nazwa.className = 'dr-wezel__mowca';
    nazwa.textContent = `${para.pierwszy} ↔ ${para.drugi}`;
    const liczba = document.createElement('span');
    liczba.className = 'dr-wezel__meta';
    liczba.textContent = `zdań wspólnych: ${para.wspolnych}`;
    pozycja.append(nazwa, liczba);
    lista.append(pozycja);
  }
  return lista;
}

/** Kontrolki okna: przełącznik widoku analizy oraz przycisk eksportu struktury tury bieżącej do dokumentu Markdown. */
interface PowierzchniaAnalizy {
  widok: HTMLSelectElement;
  eksport: HTMLButtonElement;
}

/**
 * Pasek akcji okna: jeden eksport wykonalny i cztery pozycje bez obsługi, każda nazywająca
 * komendę, która ją wykona po zbudowaniu.
 */
function zlozAkcjeAnalizy(gospodarz: HTMLElement): HTMLButtonElement {
  const eksport = przycisk('Eksportuj strukturę', 'dn-btn dn-btn--atrament');
  gospodarz.append(
    eksport,
    przyciskBezKomendy(
      'Graf argumentów',
      powodBrakuObslugi(
        Command.RoundtableArgumentList,
        'Węzły i krawędzie grafu oddaje ta komenda, a relacje niosą pola replyToId i speechAct wypowiedzi. Do czasu zbudowania obsługi okno pokazuje chronologię i mówi wprost, że grafem ona nie jest.',
      ),
    ),
    przyciskBezKomendy(
      'Wydobycie argumentów',
      powodBrakuObslugi(
        Command.RoundtableAnalysisRun,
        'Wydobycie przesłanek i wniosków zleca ta komenda w odmianie argumentMining.',
      ),
    ),
    przyciskBezKomendy(
      'Wskazanie cruxa',
      powodBrakuObslugi(
        Command.RoundtableCruxGet,
        'Kluczowy punkt sporny wskazuje ta komenda z zależności między węzłami grafu.',
      ),
    ),
    przyciskBezKomendy(
      'Eksport grafu (SVG, PNG, DOT, GraphML, Argdown, AIF)',
      powodBrakuObslugi(
        Command.RoundtableArgumentExport,
        'Graf w formacie wymiany albo jako obraz wydaje ta komenda. Okno eksportuje dziś chronologię w Markdown zamiast składać plik grafu z krawędzi, których nikt nie policzył.',
      ),
    ),
  );
  return eksport;
}

/** Zestaw akcji warstwy trzeciej — operacje analityczne jeszcze niezbudowane, wypisane obok eksportu wykonalnego. */
function zlozZestawAnalizy(czynnosci: HTMLButtonElement[]): HTMLElement {
  return utworzZestawAkcji('Operacje analityczne okna', czynnosci);
}

/** Składa przełącznik widoku, pasek akcji, warstwy analizy oraz ciało ramy okna z gotowymi kontrolkami sterowania. */
function zlozPowierzchnieAnalizy(
  rama: { akcje: HTMLElement; narzedzia: HTMLElement; cialo: HTMLElement },
  stanTresci: HTMLElement,
  czynnosci: HTMLButtonElement[],
): PowierzchniaAnalizy {
  const widok = wybor('Widok analizy', OPISY_WIDOKU);
  const eksport = zlozAkcjeAnalizy(rama.akcje);
  rama.narzedzia.append(
    wiersz('Widok analizy', widok, {
      klasa: 'dr-wiersz',
      objasnienie: 'Trzy widoki tej samej tury: chronologia, zbieżność zdań i macierz par mówców.',
    }),
  );
  rama.cialo.append(zlozZestawAnalizy(czynnosci), stanTresci, utworzWykazFunkcji('argument-map-analysis'));
  return { widok, eksport };
}
