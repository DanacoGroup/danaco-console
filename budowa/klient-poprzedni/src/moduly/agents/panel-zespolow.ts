import { ChangeKind, type Agent, type Team } from '../../../../shared/contract';
import { opisOdmowyBledu } from '../../komponenty/odmowa';
import {
  poleTekstowe,
  przycisk,
  utworzWierszOdpowiedzi,
  type WierszOdpowiedzi,
} from '../../modele/kontrolki-formularza';
import type { Kanal } from '../../protokol/kanal';
import { utworzSkladZespolu, type SkladZespolu } from './sklad-zespolu';
import {
  utworzZrodloZespolow,
  zdanieOZespole,
  type ZrodloZespolow,
} from './zrodlo-zespolow';

/**
 * Panel „Zespoły ekspertów” Agent Buildera — dostęp do czterech komend obszaru
 * `team.*`.
 *
 * Panel stoi przy bibliotece ekspertów, bo zespół to nazwany skład tej właśnie
 * biblioteki; osobne okno musiałoby wykaz ekspertów powielić albo pokazywać
 * skład samymi identyfikatorami. Zapis bez wskazanego zespołu zakłada nowy —
 * to drugie znaczenie tego samego przycisku, wypisane przy nim wprost. Zespół
 * założony na innym urządzeniu konta dochodzi zdarzeniem `team.changed`, więc
 * panel nie odpytuje rdzenia w pętli.
 */
export interface PanelZespolow {
  element: HTMLElement;
  /** Nanosi bibliotekę ekspertów — z niej składa się zespół. */
  ustawBiblioteke(eksperci: readonly Agent[]): void;
  /** Odczytuje wykaz zespołów z rdzenia. */
  wczytaj(): Promise<void>;
  /** Odłącza subskrypcję zdarzeń rdzenia. */
  rozlacz(): void;
}

export function utworzPanelZespolow(kanal: Kanal): PanelZespolow {
  const zrodlo: ZrodloZespolow = utworzZrodloZespolow(kanal);

  /** Zespół otwarty w formularzu; pusty znaczy „zapis założy nowy”. */
  let otwarty = '';
  let zespoly: Team[] = [];

  const odpowiedz: WierszOdpowiedzi = utworzWierszOdpowiedzi();

  const nazwa = poleTekstowe({
    etykieta: 'Nazwa zespołu',
    podpowiedz: 'np. Redakcja umowy',
    opis: 'Nazwa jest tym, po czym Operator poznaje zespół na wykazie.',
  });

  const opis = poleTekstowe({
    etykieta: 'Przeznaczenie zespołu',
    podpowiedz: 'do czego ten skład służy',
    opis: 'Pole puste przy zmianie zostawia opis zapisany wcześniej nietknięty.',
  });

  const sklad: SkladZespolu = utworzSkladZespolu(() => odpowiedz.wyczysc());

  const zapisz = przycisk('Zapisz zespół', 'dn-btn dn-btn--sm dn-btn--atrament');
  const nowy = przycisk('Zacznij nowy zespół', 'dn-btn dn-btn--sm');

  const zdanieOZapisie = document.createElement('p');
  zdanieOZapisie.className = 'dn-pole-opis';

  const wykaz = document.createElement('ul');
  wykaz.className = 'da-zespoly';

  function opiszZapis(): void {
    const otwartyZespol = zespoly.find((wpis) => wpis.id === otwarty);
    zdanieOZapisie.textContent =
      otwartyZespol === undefined
        ? 'Zapis założy zespół nowy — żaden nie jest otwarty w formularzu.'
        : `Zapis zmieni zespół „${otwartyZespol.name}". Aby założyć nowy, zacznij nowy zespół.`;
  }

  function wchlon(zespol: Team): void {
    const pozycja = zespoly.findIndex((wpis) => wpis.id === zespol.id);
    if (pozycja === -1) zespoly = [...zespoly, zespol];
    else zespoly = zespoly.map((wpis) => (wpis.id === zespol.id ? zespol : wpis));
    przerysujWykaz();
  }

  function otworz(zespol: Team): void {
    otwarty = zespol.id;
    nazwa.kontrolka.value = zespol.name;
    opis.kontrolka.value = zespol.description ?? '';
    sklad.ustawSklad(zespol.agentIds);
    przerysujWykaz();
  }

  async function wczytajZespol(idZespolu: string): Promise<void> {
    odpowiedz.pokaz('Wczytywanie zespołu…', true);
    const wynik = await zrodlo.wczytaj(idZespolu);
    if (!wynik.udany || wynik.wynik === undefined) {
      odpowiedz.pokaz(opisOdmowyBledu('Wczytanie zespołu', wynik.blad), false);
      return;
    }
    const zespol = wynik.wynik.team;
    wchlon(zespol);
    otworz(zespol);
    odpowiedz.pokaz(
      `Zespół „${zespol.name}" stoi w formularzu — ${zdanieOZespole(zespol)}.`,
      true,
    );
  }

  async function powiel(zespol: Team): Promise<void> {
    odpowiedz.pokaz(`Powielanie zespołu „${zespol.name}"…`, true);
    // Nazwa kopii idzie pusta — wtedy rdzeń bierze nazwę źródła z przyrostkiem.
    // Ułożenie nazwy w kliencie dałoby drugą konwencję nazewniczą obok rdzeniowej.
    const wynik = await zrodlo.powiel(zespol.id, '');
    if (!wynik.udany || wynik.wynik === undefined) {
      odpowiedz.pokaz(opisOdmowyBledu('Powielenie zespołu', wynik.blad), false);
      return;
    }
    const kopia = wynik.wynik.team;
    wchlon(kopia);
    odpowiedz.pokaz(`Kopia stoi na wykazie pod nazwą „${kopia.name}" — ${zdanieOZespole(kopia)}.`, true);
  }

  async function zapiszZespol(): Promise<void> {
    const nadana = nazwa.kontrolka.value.trim();
    if (nadana === '') {
      odpowiedz.pokaz('Nadaj zespołowi nazwę — rdzeń odmówi zapisu bez niej.', false);
      return;
    }
    odpowiedz.pokaz(`Zapisywanie zespołu „${nadana}"…`, true);
    const wynik = await zrodlo.zapisz({
      idZespolu: otwarty,
      nazwa: nadana,
      opis: opis.kontrolka.value,
      idEkspertow: sklad.wybrani(),
    });
    if (!wynik.udany || wynik.wynik === undefined) {
      odpowiedz.pokaz(opisOdmowyBledu('Zapis zespołu', wynik.blad), false);
      return;
    }
    const zapisany = wynik.wynik.team;
    const zalozony = otwarty === '';
    wchlon(zapisany);
    otworz(zapisany);
    odpowiedz.pokaz(
      `${zalozony ? 'Zespół założony' : 'Zespół zmieniony'}: „${zapisany.name}" — ${zdanieOZespole(zapisany)}.`,
      true,
    );
  }

  function zacznijNowy(): void {
    otwarty = '';
    nazwa.kontrolka.value = '';
    opis.kontrolka.value = '';
    sklad.ustawSklad([]);
    przerysujWykaz();
    odpowiedz.pokaz('Formularz jest pusty — zapis założy zespół nowy.', true);
  }

  function wierszZespolu(zespol: Team): HTMLElement {
    const tytul = document.createElement('span');
    tytul.className = 'da-zespoly__nazwa';
    tytul.textContent = zespol.name;

    const pod = document.createElement('span');
    pod.className = 'dn-pole-opis';
    pod.textContent = zdanieOZespole(zespol);

    const otworzGo = przycisk('Wczytaj', 'dn-btn dn-btn--sm');
    otworzGo.addEventListener('click', () => void wczytajZespol(zespol.id));

    const powielGo = przycisk('Powiel', 'dn-btn dn-btn--sm');
    powielGo.addEventListener('click', () => void powiel(zespol));

    const pasek = document.createElement('div');
    pasek.className = 'da-panel__pasek';
    pasek.append(otworzGo, powielGo);

    const wiersz = document.createElement('li');
    wiersz.className = 'da-zespoly__wiersz';
    wiersz.dataset['zespol'] = zespol.id;
    if (zespol.id === otwarty) wiersz.dataset['otwarty'] = 'tak';
    wiersz.append(tytul, pod, pasek);
    return wiersz;
  }

  function przerysujWykaz(): void {
    const uporzadkowane = [...zespoly].sort((pierwszy, drugi) =>
      pierwszy.name.localeCompare(drugi.name, 'pl'),
    );
    wykaz.replaceChildren(...uporzadkowane.map(wierszZespolu));
    opiszZapis();
  }

  zapisz.addEventListener('click', () => void zapiszZespol());
  nowy.addEventListener('click', () => zacznijNowy());

  const odsubskrybuj = zrodlo.naZmiane((tresc) => {
    if (tresc.change === ChangeKind.Deleted) {
      zespoly = zespoly.filter((wpis) => wpis.id !== tresc.team.id);
      if (otwarty === tresc.team.id) otwarty = '';
      przerysujWykaz();
      return;
    }
    wchlon(tresc.team);
  });

  const pasek = document.createElement('div');
  pasek.className = 'da-panel__pasek';
  pasek.append(zapisz, nowy);

  const tytul = document.createElement('h4');
  tytul.className = 'da-panel__tytul';
  tytul.textContent = 'Zespoły ekspertów';

  const element = document.createElement('div');
  element.className = 'da-panel da-zespoly__panel';
  element.dataset['panel'] = 'zespoly';
  element.append(tytul, wykaz, nazwa.element, opis.element, sklad.element, pasek, zdanieOZapisie, odpowiedz.element);

  przerysujWykaz();

  return {
    element,

    ustawBiblioteke(eksperci) {
      sklad.ustawBiblioteke(eksperci);
    },

    async wczytaj() {
      // Fraza pusta — panel pobiera cały wykaz. Komenda `team.list` przyjmuje
      // `query`, ale ekran nie ma pola zawężającego.
      const wynik = await zrodlo.wykaz('');
      if (!wynik.udany || wynik.wynik === undefined) {
        odpowiedz.pokaz(opisOdmowyBledu('Odczyt wykazu zespołów', wynik.blad), false);
        return;
      }
      zespoly = [...wynik.wynik.teams];
      przerysujWykaz();
      if (zespoly.length === 0) {
        odpowiedz.pokaz('Żaden zespół nie jest jeszcze zapisany — złóż pierwszy poniżej.', true);
      }
    },

    rozlacz() {
      odsubskrybuj();
    },
  };
}
