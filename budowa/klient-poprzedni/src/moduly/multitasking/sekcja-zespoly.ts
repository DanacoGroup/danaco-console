import type { Agent, Team } from '../../../../shared/contract';
import {
  poleTekstowe,
  pozycjaWykazu,
  przyciskAkcji as przycisk,
  wykaz,
} from '../../modele/kontrolki-formularza';
import type { ZrodloSekcjiPaneli } from '../../powloka/zrodlo-sekcji-paneli';
import { PRZEDROSTEK, wierszOpisu } from './kontrolki';
import { utworzPowierzchnieSekcji, zdaniePuste } from './powierzchnia-sekcji';
import type { ZrodloZespolow } from './zrodlo-zespolow';

/**
 * Sekcja ZESPOŁY panelu orkiestracji — presety składu, którymi obsadza się role.
 *
 * Sekcja steruje zespołami zapisanymi w rdzeniu: zakładaniem, powielaniem
 * i wczytywaniem. Szablony stoją tu jako gotowe nazwy do zapisania —
 * naciśnięcie szablonu zakłada zespół komendą `team.save`.
 *
 * Preset obejmuje sam skład: `TeamSaveRequest` przyjmuje nazwę, opis
 * i `agentIds`. Ról ani kolejek kontrakt nie niesie, więc sekcja wypisuje
 * kształt brakującego pola zamiast upychać je w opisie. Nadanie ról oknom
 * robi sekcja Role komendą `role.assign`.
 *
 * Skład pokazuje się nazwami ekspertów, bo `Team.agentIds` niesie same
 * identyfikatory, a nazwy dokłada `agent.list`. Ekspert nieznany wykazowi
 * zostaje w składzie — z identyfikatorem i zdaniem, że rdzeń go nie oddał.
 */
export interface SekcjaZespolow {
  element: HTMLElement;
  /** Odczyt zespołów i ekspertów z rdzenia. */
  odswiez(): void;
  /** Wskazanie okna, do którego przypięty jest układ podsekcji. */
  ustawOkno(idOkna: string): void;
  /** Odpięcie nasłuchu `team.changed`. */
  rozlacz(): void;
}

export interface OpcjeSekcjiZespolow {
  zrodlo: ZrodloZespolow;
  sekcje: ZrodloSekcjiPaneli;
}

/**
 * Szablony zespołu — nazwa i przeznaczenie.
 *
 * Skład zostaje pusty: rdzeń nie zna pojęcia roli w zespole, a dobór ekspertów
 * zależy od tego, kogo instalacja ma zarejestrowanych. Szablon zakłada zespół
 * i oddaje go do uzupełnienia.
 */
const SZABLONY: ReadonlyArray<{ nazwa: string; opis: string }> = [
  { nazwa: 'Budowa aplikacji', opis: 'Koordynator planuje, dwaj wykonawcy budują, Validator ocenia.' },
  { nazwa: 'Badanie i redakcja', opis: 'Wykonawcy zbierają materiał, Validator redaguje i rozstrzyga.' },
  { nazwa: 'Pętla ciągła 24/7', opis: 'Praca ciągła sterowana harmonogramem i kolejką.' },
];

/** Kształt pola, którego `team.save` nie ma — wypisany w meldunku wprost. */
const BRAK_PRESETU =
  'team.save { …, roles?: [{ windowRole, agentId, channelId }], queueIds?: string[] } — kontrakt tych pól nie niesie.';

export function utworzSekcjeZespolow(opcje: OpcjeSekcjiZespolow): SekcjaZespolow {
  const { zrodlo, sekcje } = opcje;
  /** Eksperci z ostatniego odczytu — wyłącznie po to, by nazwać skład. */
  let eksperci: readonly Agent[] = [];

  const nazwa = poleTekstowe({
    etykieta: 'Nazwa zespołu',
    podpowiedz: 'np. Budowa aplikacji',
    opis: 'Zespół zakłada się pusty; skład uzupełnia się ekspertami z wykazu.',
  });
  const zaloz = przycisk('Załóż zespół', 'dn-btn dn-btn--sm');

  const pasek = document.createElement('div');
  pasek.className = 'dm-orkiestracja__pasek';
  pasek.append(nazwa.element, zaloz);

  const listaZespolow = wykaz('Zespoły zapisane w rdzeniu', 'dm-wykaz');
  const listaSzablonow = wykaz('Szablony zespołu', 'dm-wykaz');
  const listaSkladu = wykaz('Skład zespołu wczytanego', 'dm-wykaz');

  const powierzchnia = utworzPowierzchnieSekcji({
    klucz: 'zespoly',
    tytul: 'Zespoły',
    zakres:
      'Presety składu, którymi obsadza się role tego środowiska: zakładanie, powielanie i wczytywanie zespołu.',
    sekcje,
    podsekcje: [
      { id: 'zaloz', tytul: 'Nowy zespół', tresc: pasek },
      { id: 'zapisane', tytul: 'Zespoły zapisane w rdzeniu', tresc: listaZespolow },
      { id: 'szablony', tytul: 'Szablony zespołu', tresc: listaSzablonow },
      { id: 'sklad', tytul: 'Skład zespołu wczytanego', tresc: listaSkladu },
    ],
  });

  zaloz.addEventListener('click', () => {
    void zapisz(nazwa.kontrolka.value.trim(), '');
  });

  listaSzablonow.replaceChildren(
    ...SZABLONY.map((szablon) => {
      const { element, akcje } = pozycjaWykazu(szablon.nazwa, szablon.opis, PRZEDROSTEK);
      const kontrolka = przycisk('Załóż z szablonu', 'dn-btn dn-btn--sm dn-btn--zarys');
      kontrolka.addEventListener('click', () => {
        void zapisz(szablon.nazwa, szablon.opis);
      });
      akcje.append(kontrolka);
      return element;
    }),
  );

  /**
   * Zapis zespołu. Pusta nazwa nie jedzie do rdzenia: `TeamSaveRequest.name`
   * jest polem obowiązkowym, więc żądanie bez niej skończyłoby się odmową
   * trudną do powiązania z pustym polem formularza.
   */
  async function zapisz(mianoZespolu: string, opis: string): Promise<void> {
    if (mianoZespolu === '') {
      powierzchnia.meldunek('Zespół bez nazwy nie ma czego zapisać — nazwa jest polem obowiązkowym kontraktu.', false);
      return;
    }
    const wynik = await zrodlo.zapisz({
      name: mianoZespolu,
      agentIds: [],
      ...(opis === '' ? {} : { description: opis }),
    });
    if (!wynik.udany || wynik.wynik === undefined) {
      powierzchnia.meldunek(
        `Rdzeń odmówił zapisu zespołu „${mianoZespolu}". ${powod(wynik.blad?.message, wynik.blad?.code)}`,
        false,
      );
      return;
    }
    // Nazwa i identyfikator w meldunku pochodzą z zespołu oddanego przez rdzeń
    // po zapisie, nie z treści żądania.
    powierzchnia.meldunek(
      `Rdzeń założył zespół „${wynik.wynik.name}" (${wynik.wynik.id}). Preset obejmuje SAM SKŁAD — ${BRAK_PRESETU}`,
      true,
    );
    nazwa.kontrolka.value = '';
    odczytaj();
  }

  /** Powielenie zespołu — kopia idzie do rdzenia, nie do widoku. */
  async function powiel(zespol: Team): Promise<void> {
    const wynik = await zrodlo.powiel({ teamId: zespol.id });
    if (!wynik.udany || wynik.wynik === undefined) {
      powierzchnia.meldunek(
        `Rdzeń odmówił powielenia zespołu „${zespol.name}". ${powod(wynik.blad?.message, wynik.blad?.code)}`,
        false,
      );
      return;
    }
    powierzchnia.meldunek(
      `Rdzeń założył kopię zespołu: „${wynik.wynik.name}" (${wynik.wynik.id}).`,
      true,
    );
    odczytaj();
  }

  /** Wczytanie zespołu — skład idzie do podsekcji składu, nazwami ekspertów. */
  async function wczytaj(zespol: Team): Promise<void> {
    const wynik = await zrodlo.wczytaj({ teamId: zespol.id });
    if (!wynik.udany || wynik.wynik === undefined) {
      powierzchnia.meldunek(
        `Rdzeń odmówił odczytu zespołu „${zespol.name}". ${powod(wynik.blad?.message, wynik.blad?.code)}`,
        false,
      );
      return;
    }
    pokazSklad(wynik.wynik);
    powierzchnia.meldunek(
      `Zespół „${wynik.wynik.name}" wczytany; skład liczy ${wynik.wynik.agentIds.length} ekspertów.`,
      true,
    );
  }

  /** Skład zespołu nazwami ekspertów; ekspert nieznany zostaje z powodem. */
  function pokazSklad(zespol: Team): void {
    if (zespol.agentIds.length === 0) {
      listaSkladu.replaceChildren(
        zdaniePuste(
          `Zespół „${zespol.name}" nie ma jeszcze ani jednego eksperta. Pusty skład jest stanem poprawnym — zespół zakłada się przed obsadzeniem.`,
        ),
      );
      return;
    }
    listaSkladu.replaceChildren(
      ...zespol.agentIds.map((id) => {
        const ekspert = eksperci.find((wpis) => wpis.id === id);
        const { element } = pozycjaWykazu(
          ekspert?.displayName ?? ekspert?.name ?? id,
          ekspert === undefined
            ? 'Rdzeń nie oddał tego eksperta w wykazie — pokazujemy identyfikator.'
            : (ekspert.description ?? `model ${ekspert.model ?? 'nie wskazany'}`),
          PRZEDROSTEK,
        );
        return element;
      }),
    );
  }

  /** Odczyt obu wykazów: eksperci nazywają skład, zespoły są treścią sekcji. */
  function odczytaj(): void {
    powierzchnia.tresci.ladowanie('Odczyt zespołów i ekspertów…');
    void (async () => {
      const wykazEkspertow = await zrodlo.eksperci({});
      eksperci = wykazEkspertow.udany ? (wykazEkspertow.wynik ?? []) : [];

      const wynik = await zrodlo.zespoly({});
      if (!wynik.udany || wynik.wynik === undefined) {
        powierzchnia.tresci.blad('Odczyt zespołów odmówiony.', wynik.blad);
        return;
      }
      powierzchnia.tresci.pusto('');
      if (wynik.wynik.length === 0) {
        listaZespolow.replaceChildren(
          zdaniePuste(
            'Rdzeń nie ma zapisanego ani jednego zespołu. To nie jest usterka — zespół zakłada Operator, a szablony niżej robią to jednym naciśnięciem.',
          ),
        );
        return;
      }
      listaZespolow.replaceChildren(...wynik.wynik.map(pozycjaZespolu));
    })();
  }

  /** Jeden zespół jako pozycja wykazu wraz z dwiema czynnościami. */
  function pozycjaZespolu(zespol: Team): HTMLElement {
    const { element, akcje } = pozycjaWykazu(
      zespol.name,
      zespol.description ?? `skład: ${zespol.agentIds.length} ekspertów`,
      PRZEDROSTEK,
    );
    const wczytajGo = przycisk('Wczytaj skład', 'dn-btn dn-btn--sm dn-btn--zarys');
    wczytajGo.addEventListener('click', () => {
      void wczytaj(zespol);
    });
    const powielGo = przycisk('Powiel', 'dn-btn dn-btn--sm dn-btn--zarys');
    powielGo.addEventListener('click', () => {
      void powiel(zespol);
    });
    akcje.append(wczytajGo, powielGo);
    element.append(wierszOpisu('Założony', new Date(zespol.createdAt).toLocaleString('pl-PL')));
    return element;
  }

  listaSkladu.replaceChildren(
    zdaniePuste('Wczytaj zespół z wykazu wyżej — jego skład stanie tutaj.'),
  );

  // Zespół zmieniony poza tym widokiem przerysowuje wykaz: `team.changed`
  // istnieje po to, żeby drugie okno nie musiało pytać ponownie w pętli.
  const odsubskrybuj = zrodlo.naZespol(() => odczytaj());

  return {
    element: powierzchnia.element,
    odswiez: odczytaj,
    ustawOkno: powierzchnia.ustawOkno,
    rozlacz: odsubskrybuj,
  };
}

/** Treść odmowy wraz z kodem kontraktu — ta sama postać co w oknach ról. */
function powod(zdanie: string | undefined, kod: string | undefined): string {
  return `Powód: ${zdanie ?? 'rdzeń nie podał przyczyny'} (kod ${kod ?? 'brak'}).`;
}
