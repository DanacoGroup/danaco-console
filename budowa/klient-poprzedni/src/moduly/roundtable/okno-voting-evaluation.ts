import { Command } from '../../../../shared/contract';
import { utworzRameOkna } from '../../komponenty/rama-okna';
import {
  pobierzPlik,
  przyciskAkcji as przycisk,
  przyciskBezKomendy,
} from '../../modele/kontrolki-formularza-braki';
import { powodBrakuObslugi } from './braki-kontraktu';
import type { StanDebaty } from './stan-debaty';
import { utworzStanTresci } from './stany-okna';
import type { StrumienWypowiedzi } from './strumien-wypowiedzi';
import { czynnosciOceny } from './czynnosci-arsenalu';
import { utworzWykazFunkcji, utworzZestawAkcji } from './warstwy-modulu';
import type { ZrodloArsenaluRoundtable } from './zrodlo-arsenalu';
import {
  zdanieUdzialu,
  zestawienieMarkdown,
  zlozZestawienieUdzialu,
  type UdzialUczestnika,
} from './zestawienie-udzialu';

/**
 * Voting & Evaluation Center — okno monitora otwierane jako rozszerzenie boczne.
 *
 * Okno nie prowadzi głosowania i nie udaje, że je prowadzi — ale powód zmienił
 * się co do istoty. Kontrakt niesie dziś pełną rodzinę głosowania
 * (`roundtable.vote.start`, `.cast`, `.get`) wraz z metodą agregacji, progiem
 * kworum i wynikiem, a także rubryki, sędziów, ranking, macierz decyzyjną
 * i kalibrację. Brakuje ich obsługi: żadnej z tych komend to okno jeszcze nie
 * wywołuje. Suwak, gwiazdki albo przycisk „Uruchom głosowanie” postawione bez
 * obsługi zbierałyby wybór Operatora, który nie dociera nigdzie i znika wraz
 * z odświeżeniem okna — czyli byłyby pozorem czynności.
 *
 * Okno pokazuje zamiast tego jedyną wielkość, którą samo umie zmierzyć: udział
 * uczestników w turze bieżącej — ile razy i jak obszernie każdy się odezwał, kto
 * milczy, kto jest wyciszony i jakie ma miejsce w kolejności głosu. Nazwa
 * wielkości mówi, czym ona jest: udziałem, nie rankingiem i nie kworum.
 *
 * Wyliczenia siedzą w `zestawienie-udzialu.ts`. Subskrypcji strumienia okno nie
 * zakłada — jedna na całe złożenie stoi w `indeks.ts`.
 */
export interface OknoVotingEvaluation {
  element: HTMLElement;
  odswiez(): void;
  /** Przerysowanie po fragmencie strumienia — wołane przez złożenie modułu. */
  odswiezGlosy(): void;
  /** Zamyka nasłuch `stan.naZmiane(...)`. */
  zamknij(): void;
}

export function utworzOknoVotingEvaluation(
  arsenal: ZrodloArsenaluRoundtable,
  stan: StanDebaty,
  strumien: StrumienWypowiedzi,
): OknoVotingEvaluation {
  const rama = utworzRameOkna({
    tytul: 'Voting & Evaluation Center',
    rola: 'monitor',
    kod: 'voting-evaluation-center',
    przeznaczenie:
      'Udział uczestników w turze bieżącej — wielkość zmierzona z wypowiedzi; głosowania i rubryki są w kontrakcie, obsługi jeszcze nie zbudowano.',
    przedrostek: 'dr',
  });
  const tresc = utworzStanTresci();
  // Czynności oceny wołają komendy obszaru wprost: głosowanie, ocenę parami,
  // rubryki, sędziów, ranking i macierz decyzyjną. Warianty głosowania bierze
  // się z wypowiedzi tury bieżącej — to nad nimi głosuje się w tym oknie.
  const powierzchnia = zlozPowierzchnieOceny(
    rama,
    tresc.element,
    czynnosciOceny(
      arsenal,
      stan,
      (zdanie, powodzenie) => tresc.potwierdzenie(zdanie, powodzenie),
      () => stan.wypowiedzi().map((wypowiedz) => wypowiedz.id).join('\n'),
      () => 'approval',
    ),
  );

  function rysuj(): void {
    const zestawienie = zlozZestawienieUdzialu(stan, strumien);
    if (zestawienie.skladu === 0) {
      tresc.pusto(
        'Debata nie ma jeszcze uczestników znanych temu oknu — udziału nie ma z czego policzyć. ' +
          'Skład dodaje się w Model Panels.',
      );
      rama.ustawZnacznik('bez składu', 'neutralna');
      return;
    }
    rama.ustawZnacznik(
      `udział: ${zestawienie.mowiacych} z ${zestawienie.skladu}`,
      zestawienie.mowiacych === 0 ? 'ostrzezenie' : 'neutralna',
    );

    const miejsce = tresc.tresc();

    const udzial = document.createElement('p');
    udzial.className = 'dr-analiza__zdanie';
    udzial.textContent = zdanieUdzialu(zestawienie);
    miejsce.append(udzial, wykazUdzialu(zestawienie.wiersze), zdanieBrakuGlosowania());
  }

  function odswiezGlosy(): void {
    const rodzaj = tresc.rodzaj();
    if (rodzaj === 'blad' || rodzaj === 'ladowanie') return;
    rysuj();
  }

  function eksportuj(): void {
    const zestawienie = zlozZestawienieUdzialu(stan, strumien);
    if (zestawienie.skladu === 0) {
      tresc.potwierdzenie('Nie ma czego wyeksportować — okno nie zna jeszcze składu debaty.', false);
      return;
    }
    pobierzPlik(
      `udzial-w-debacie-${stan.tura() === '' ? 'bez-tury' : stan.tura()}.md`,
      zestawienieMarkdown(zestawienie, opisTury(stan)),
      'text/markdown',
    );
    tresc.potwierdzenie(
      'Zestawienie udziału pobrane jako plik Markdown. Wyników głosowań w pliku nie ma — ich komendy są w kontrakcie, ale okno jeszcze ich nie wywołuje.',
      true,
    );
  }

  powierzchnia.eksport.addEventListener('click', eksportuj);

  const odsubskrybuj = stan.naZmiane(rysuj);
  rysuj();

  return { element: rama.element, odswiez: rysuj, odswiezGlosy, zamknij: odsubskrybuj };
}

/** Opis tury bieżącej dla eksportu i zdań okna. */
function opisTury(stan: StanDebaty): string {
  const definicja = stan.definicjaTury();
  if (definicja === null) return 'debaty bez tury znanej temu oknu';
  return `tury #${definicja.index} (${definicja.status})`;
}

/** Wykaz udziału — wiersz na uczestnika, stany słowem. */
function wykazUdzialu(wiersze: readonly UdzialUczestnika[]): HTMLElement {
  const lista = document.createElement('ul');
  lista.className = 'dr-wezly';
  lista.setAttribute('aria-label', 'Udział uczestników w turze bieżącej');
  for (const udzial of wiersze) {
    const pozycja = document.createElement('li');
    pozycja.className = 'dr-wezel';
    pozycja.dataset['mowca'] = udzial.idUczestnika;
    if (udzial.wypowiedzi === 0) pozycja.dataset['bezGlosu'] = 'tak';

    const znacznik = document.createElement('span');
    znacznik.className = 'dr-wezel__znacznik';
    znacznik.textContent = udzial.znacznik;

    const nazwa = document.createElement('span');
    nazwa.className = 'dr-wezel__mowca';
    nazwa.textContent = udzial.nazwa;

    const meta = document.createElement('span');
    meta.className = 'dr-wezel__meta';
    meta.textContent = zdanieWiersza(udzial);

    pozycja.append(znacznik, nazwa, meta);
    lista.append(pozycja);
  }
  return lista;
}

/** Wszystkie stany wiersza słowem — atrybut danych nie niesie ich nikomu. */
function zdanieWiersza(udzial: UdzialUczestnika): string {
  const czlony = [
    udzial.wypowiedzi === 0
      ? 'bez wypowiedzi widzianej przez to okno'
      : `wypowiedzi: ${udzial.wypowiedzi}`,
    `znaków: ${udzial.znakow}`,
    udzial.kolejnosc === null
      ? 'kolejność głosu nie wskazana przez rdzeń'
      : `kolejność głosu #${udzial.kolejnosc}`,
    udzial.wyciszony ? 'wyciszony w turze' : 'słyszany w turze',
    udzial.stanGlosu,
  ];
  if (udzial.kluczowy) czlony.push('oznaczony przez rdzeń jako kluczowy');
  return czlony.join(' · ');
}

/**
 * Zdanie zamykające wykaz: czego w oknie nie ma i dlaczego.
 *
 * Stoi pod wykazem, nie w dymku przycisku, bo dotyczy całego okna: Operator ma
 * poznać granicę bez naciskania czegokolwiek.
 */
function zdanieBrakuGlosowania(): HTMLElement {
  const zdanie = document.createElement('p');
  zdanie.className = 'dr-analiza__zdanie';
  zdanie.dataset['granica'] = 'tak';
  zdanie.textContent =
    'Głosowania w tym oknie nie ma i nie jest ono ukryte. Komendy głosowania są w kontrakcie — roundtable.vote.start ' +
    'otwiera je wraz z metodą agregacji i progiem kworum, roundtable.vote.cast przyjmuje głos, roundtable.vote.get oddaje wynik — ' +
    'ale obsługi żadnej z nich jeszcze nie zbudowano. Głos oddany tutaj nie dotarłby do rdzenia i zniknąłby przy odświeżeniu okna, ' +
    'dlatego okno mierzy udział zamiast zbierać głosy.';
  return zdanie;
}

/** Kontrolki okna. */
interface PowierzchniaOceny {
  eksport: HTMLButtonElement;
}

/**
 * Pasek akcji okna: jeden raport wykonalny i cztery pozycje bez obsługi, wśród
 * nich obie czynności wiodące opracowania — uruchomienie głosowania i wybór
 * metody agregacji. Wszystkie cztery mają dziś komendę w kontrakcie.
 */
function zlozAkcjeOceny(gospodarz: HTMLElement): HTMLButtonElement {
  const eksport = przycisk('Eksportuj zestawienie udziału', 'dn-btn dn-btn--atrament');
  gospodarz.append(
    eksport,
    przyciskBezKomendy(
      'Uruchom głosowanie',
      powodBrakuObslugi(
        Command.RoundtableVoteStart,
        'Głosowanie otwiera ta komenda; głos przyjmuje roundtable.vote.cast, a wynik po agregacji oddaje roundtable.vote.get.',
      ),
    ),
    przyciskBezKomendy(
      'Metoda głosowania',
      powodBrakuObslugi(
        Command.RoundtableVoteStart,
        'Aprobatę, ranking IRV, Condorceta w wariancie Schulzego, skalę punktową i metodę kwadratową niesie pole method żądania tej komendy.',
      ),
    ),
    przyciskBezKomendy(
      'Próg kworum i konsensusu',
      powodBrakuObslugi(
        Command.RoundtableVoteStart,
        'Próg zgody niesie pole quorum żądania tej komendy, a jego osiągnięcie orzeka wynik z roundtable.vote.get. Do tego czasu okno pokazuje udział zmierzony i nazywa go udziałem, a nie kworum.',
      ),
    ),
    przyciskBezKomendy(
      'Ranking uczestników',
      powodBrakuObslugi(
        Command.RoundtableLeaderboardGet,
        'Ranking akumulowany między sesjami oddaje ta komenda wraz z algorytmem i zakresem akumulacji.',
      ),
    ),
  );
  return eksport;
}

/** Zestaw akcji warstwy trzeciej — operacje oceny jeszcze niezbudowane. */
function zlozZestawOceny(czynnosci: HTMLButtonElement[]): HTMLElement {
  return utworzZestawAkcji('Operacje oceny okna', czynnosci);
}

/** Składa pasek akcji, warstwy i ciało ramy. */
function zlozPowierzchnieOceny(
  rama: { akcje: HTMLElement; cialo: HTMLElement },
  stanTresci: HTMLElement,
  czynnosci: HTMLButtonElement[],
): PowierzchniaOceny {
  const eksport = zlozAkcjeOceny(rama.akcje);
  rama.cialo.append(zlozZestawOceny(czynnosci), stanTresci, utworzWykazFunkcji('voting-evaluation-center'));
  return { eksport };
}
