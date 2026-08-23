import { Command, type Agent } from '../../../../shared/contract';
import { zdanieOPrzypisaniach, zmienUprawnienieEksperta } from './agenci-czynnosci';
import { pozycjaAgenta } from './agenci-pozycja';
import type { WykazBrakow } from './braki-kontraktu';
import { utworzRameOkna } from '../../komponenty/rama-okna';
import {
  pole,
  przelacznikWidoku,
  przestaw,
  przyciskAkcji as przycisk,
  wiersz,
  wybor,
  wykaz,
} from '../../modele/kontrolki-formularza';
import type { StanProjektu } from './stan-projektu';
import { utworzStanTresci } from './stany-okna';
import type { ZrodloWorkspace } from './zrodlo-workspace';

/**
 * Agent Manager — okno zarządca modułu Workspace: przypisanie agenta jako
 * wykonawcy zadania w projekcie.
 *
 * Stąd dwa wykazy obok siebie: biblioteka ekspertów z modułu Agents
 * (`agent.list`) i eksperci przypisani do projektu — ci drudzy z zestawienia
 * pulpitu, bo kontrakt nie ma komendy `workspace.agent.list`.
 *
 * Uprawnienia są konfiguracją możliwości, nie kontrolą dostępu. Cztery grupy
 * zakresu ustawia komenda modułu Agents, a stan każdej z nich przychodzi w polu
 * `Agent.permissions` — okno go nie zgaduje.
 *
 * O przypisaniu do projektu mówi rdzeń albo nikt. Dopóki pulpit nie został
 * odczytany, o przypisaniach nie wiadomo nic i okno tak właśnie mówi — zamiast
 * przedstawiać nieodczytany pulpit jako projekt bez ekspertów.
 */
export interface OknoAgentow {
  element: HTMLElement;
  odswiez(): void;
}

export function utworzOknoAgentow(
  zrodlo: ZrodloWorkspace,
  stan: StanProjektu,
  braki: WykazBrakow,
): OknoAgentow {
  const rama = utworzRameOkna({
    tytul: 'Agent Manager',
    kod: 'agent-manager',
    rola: 'zarządca',
    przeznaczenie: 'Eksperci pracujący w projekcie: przypisanie, rola i wskazanie domyślnego wykonawcy.',
    przedrostek: 'dw',
  });
  const tresc = utworzStanTresci();

  const ekspert = wybor('Ekspert do przypisania', [['', 'wykaz nieodczytany']]);
  const rola = pole('Rola eksperta w projekcie', 'np. redakcja pism');
  const domyslny = document.createElement('input');
  domyslny.type = 'checkbox';
  domyslny.className = 'dn-przelacznik';

  const przypisz = przycisk('+ Przypisz agenta', 'dn-btn dn-btn--atrament');
  const odswiez = przycisk('Odśwież wykaz ekspertów');
  const doBudowniczego = przycisk('→ Agent Builder (moduł Agents)');

  rama.akcje.append(
    przypisz,
    odswiez,
    doBudowniczego,
    braki.przyciskBraku(
      'Odłącz eksperta',
      'Odłączenie eksperta od projektu',
      Command.WorkspaceAgentUnassign,
    ),
  );

  // Zawężenie wykazu ekspertów jest czynnością wyłącznie kliencką — pracuje na
  // bibliotece już odczytanej — więc stoi w pasku narzędzi kontekstowych okna.
  const fraza = pole('Fraza w nazwie, umiejętności albo konektorze', 'fragment nazwy lub kodu');
  const tylkoPrzypisani = przelacznikWidoku('Tylko przypisani do projektu', false);
  rama.narzedzia.append(fraza, tylkoPrzypisani);

  rama.cialo.append(
    wiersz('Ekspert', ekspert, { klasa: 'dw-wiersz', objasnienie: 'Wykaz pochodzi z biblioteki ekspertów modułu Agents (agent.list).' }),
    wiersz('Rola w projekcie', rola, { klasa: 'dw-wiersz' }),
    wiersz('Domyślny wykonawca', domyslny, { klasa: 'dw-wiersz', objasnienie: 'Wskazanie zdejmuje domyślność z poprzedniego eksperta.' }),
    tresc.element,
  );

  let biblioteka: Agent[] = [];
  /**
   * Eksperci, o których rdzeń powiedział, że pracują w projekcie.
   * `null` znaczy „rdzeń nie był o to pytany” i nie jest tym samym co zero.
   */
  let przypisaniRdzenia: Set<string> | null = null;
  /** Stan przełącznika „Tylko przypisani” — nośnikiem jest przycisk, nie pole. */
  let filtrPrzypisani = false;

  const czynnosci = { zrodlo, tresc, odrysuj: () => pokaz() };

  /**
   * Czy ekspert przechodzi przez nastawę filtra. Fraza obejmuje nazwę,
   * identyfikator, umiejętności i konektory — czyli to, po czym opracowanie każe
   * eksperta odnajdywać, i wyłącznie te pola, które kontrakt w `agent.list` niesie.
   */
  function przechodziFiltr(agent: Agent): boolean {
    const szukana = fraza.value.trim().toLocaleLowerCase('pl-PL');
    if (szukana !== '') {
      const przeszukiwane = [
        agent.name,
        agent.id,
        ...(agent.skillIds ?? []),
        ...(agent.connectorIds ?? []),
      ];
      if (!przeszukiwane.some((czlon) => czlon.toLocaleLowerCase('pl-PL').includes(szukana))) {
        return false;
      }
    }
    if (!filtrPrzypisani) return true;
    return przypisaniRdzenia !== null && przypisaniRdzenia.has(agent.id);
  }

  function pozycjaEksperta(agent: Agent): HTMLElement {
    return pozycjaAgenta(agent, przypisaniRdzenia === null ? null : przypisaniRdzenia.has(agent.id), {
      zmienUprawnienie: (grupa, opis, chciane) =>
        zmienUprawnienieEksperta(czynnosci, agent, grupa, opis, chciane),
    });
  }

  function pokaz(): void {
    if (biblioteka.length === 0) {
      rama.ustawZnacznik('');
      tresc.pusto('Biblioteka ekspertów jest pusta. Eksperta zakłada się w module Agents (Agent Builder).');
      return;
    }
    // Zawężenie do przypisanych ma czym pracować dopiero wtedy, gdy rdzeń
    // o przypisaniach powiedział. Przy niewiedzy wykaz nie udaje pustego —
    // nazywa brakującą przesłankę i drogę do niej.
    if (filtrPrzypisani && przypisaniRdzenia === null) {
      rama.ustawZnacznik('przypisania nieodczytane', 'ostrzezenie');
      tresc.pusto(
        'Zawężenie do przypisanych opiera się na zestawieniu pulpitu — to ono niesie wykaz ' +
          'ekspertów projektu. Odczytaj zestawienie w oknie Project Dashboard albo zdejmij ' +
          'zawężenie w pasku narzędzi.',
      );
      return;
    }
    const widoczni = biblioteka.filter(przechodziFiltr);
    rama.ustawZnacznik(
      widoczni.length === biblioteka.length
        ? `eksperci: ${biblioteka.length}`
        : `eksperci: ${widoczni.length} z ${biblioteka.length}`,
    );
    if (widoczni.length === 0) {
      tresc.pusto('Żaden ekspert biblioteki nie odpowiada nastawie filtra w pasku narzędzi.');
      return;
    }
    const lista = wykaz('Eksperci projektu i biblioteki', 'dw-wykaz');
    for (const agent of widoczni) lista.append(pozycjaEksperta(agent));
    const miejsce = tresc.tresc();
    miejsce.append(lista);
    miejsce.append(zdanieOPrzypisaniach(przypisaniRdzenia));
  }

  function odczytaj(): void {
    // Wiedza o przypisaniach pochodzi z zestawienia pulpitu; bierzemy jego stan
    // bieżący, żeby wykaz nie rysował się z niewiedzą, którą pulpit już rozwiał.
    const zestawienie = stan.pulpit();
    przypisaniRdzenia = zestawienie === null ? null : new Set(zestawienie.assignedAgentIds ?? []);
    tresc.ladowanie('Odczyt biblioteki ekspertów…');
    void zrodlo.agenci().then((wynik) => {
      if (!wynik.udany || wynik.wynik === undefined) {
        // Plakietka gaśnie razem z wykazem: licznik sprzed odmowy mówiłby
        // o bibliotece, której to okno właśnie nie zdołało odczytać.
        biblioteka = [];
        rama.ustawZnacznik('');
        tresc.blad('Rdzeń nie oddał biblioteki ekspertów (agent.list).', wynik.blad);
        return;
      }
      biblioteka = wynik.wynik;
      ekspert.replaceChildren();
      for (const agent of biblioteka) {
        const pozycja = document.createElement('option');
        pozycja.value = agent.id;
        pozycja.textContent = `${agent.name} (${agent.id})`;
        ekspert.append(pozycja);
      }
      pokaz();
    });
  }

  przypisz.addEventListener('click', () => {
    const projekt = stan.projekt();
    if (projekt === '') {
      tresc.blad('Przypisanie bez projektu nie ma zakresu — wskaż projekt na pulpicie.');
      return;
    }
    if (ekspert.value === '') {
      tresc.potwierdzenie('Wskaż eksperta z biblioteki, aby przypisać go do projektu.', false);
      return;
    }
    tresc.ladowanie('Zapis przypisania eksperta…');
    void zrodlo
      .przypiszAgenta({
        projectId: projekt,
        agentId: ekspert.value,
        role: rola.value,
        defaultExecutor: domyslny.checked,
      })
      .then((wynik) => {
        if (!wynik.udany || wynik.wynik === undefined) {
          tresc.blad('Rdzeń nie przyjął przypisania eksperta.', wynik.blad);
          return;
        }
        const przypisanie = wynik.wynik.assignment;
        // Wykaz przypisanych stoi na zestawieniu pulpitu, a to po zapisie jest
        // już nieaktualne. Wchłaniamy to, co oddał rdzeń, żeby wykaz pod
        // potwierdzeniem nie zaprzeczał potwierdzeniu.
        if (przypisaniRdzenia === null) przypisaniRdzenia = new Set<string>();
        przypisaniRdzenia.add(przypisanie.agentId);
        pokaz();
        tresc.potwierdzenie(
          `Rdzeń zapisał przypisanie eksperta ${przypisanie.agentId} do projektu ` +
            `${przypisanie.projectId}` +
            `${przypisanie.role === undefined || przypisanie.role === '' ? '' : `, rola „${przypisanie.role}”`}` +
            `${przypisanie.defaultExecutor === true ? ', jako domyślny wykonawca' : ''}.`,
          true,
        );
      });
  });

  odswiez.addEventListener('click', odczytaj);
  // Filtr zawęża bibliotekę już odczytaną, więc przerysowuje wykaz, nie pyta rdzenia.
  fraza.addEventListener('input', pokaz);
  tylkoPrzypisani.addEventListener('click', () => {
    filtrPrzypisani = przestaw(tylkoPrzypisani);
    pokaz();
  });

  doBudowniczego.addEventListener('click', () => {
    const okno = stan.oknoRozmowy();
    if (okno === '') {
      // Powód z odpowiedzi rdzenia (`okno-rozmowy.ts`), nie z napisu na sztywno
      // — tak samo jak w `biblioteka-czynnosci.ts`.
      tresc.blad(`Skok do modułu Agents wymaga okna rozmowy modułu. ${stan.powodBrakuOkna()}`);
      return;
    }
    tresc.ladowanie('Przenoszenie kontekstu projektu do modułu Agents…');
    void zrodlo.udostepnij(okno, 'agents', stan.projekt(), []).then((wynik) => {
      if (!wynik.udany || wynik.wynik === undefined) {
        tresc.blad('Rdzeń nie przeniósł kontekstu do modułu Agents.', wynik.blad);
        return;
      }
      // Zdanie mówi o oknie, które rdzeń otworzył, i o jego własnym znaczniku
      // przeniesienia, nie o stałej wpisanej w wywołanie.
      const odpowiedz = wynik.wynik;
      if (odpowiedz.transferred !== true) {
        tresc.potwierdzenie(
          'Rdzeń przyjął wywołanie, ale nie potwierdził przeniesienia ' +
            '(context.transfer oddał transferred = fałsz).',
          false,
        );
        return;
      }
      tresc.potwierdzenie(
        `Rdzeń przeniósł kontekst projektu do okna ${odpowiedz.window.id} ` +
          `modułu ${odpowiedz.window.moduleId}.`,
        true,
      );
    });
  });

  // Zestawienie pulpitu niesie wykaz przypisanych — po jego zmianie wykaz
  // ekspertów pokazuje inne przypisania, więc odrysowujemy go bez pytania rdzenia.
  // Zestawienie skasowane (zmiana projektu, zdarzenie rdzenia) wraca do stanu
  // „nie wiadomo”, a nie do zera: o przypisaniach nowego projektu nikt jeszcze
  // nie pytał.
  stan.naZmiane(() => {
    const zestawienie = stan.pulpit();
    przypisaniRdzenia = zestawienie === null ? null : new Set(zestawienie.assignedAgentIds ?? []);
    if (biblioteka.length > 0) pokaz();
  });

  return { element: rama.element, odswiez: odczytaj };
}
