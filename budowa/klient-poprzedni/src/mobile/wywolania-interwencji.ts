import {
  Command,
  ConfigScope,
  PermissionMode,
  QueueAction,
  ReasoningEffort,
  SessionConfigArea,
  type SessionConfig,
  type SessionConfigOrigin,
} from '../../../shared/contract';
import { opisOdmowyBledu } from '../komponenty/odmowa';
import type { Kanal, Wynik } from '../protokol/kanal';
import { wywolaj } from '../protokol/wywolanie';
import { Droga, type KrokKwitu, type Kwit } from './kwit-decyzji';
import type { PozycjaDecyzji } from './pozycje-decyzji';

/**
 * Cztery drogi interwencji jako cztery wywołania — jedyne miejsce warstwy
 * mobilnej znające nazwy komend zmieniających stan.
 *
 * Każda droga stoi na komendach kontraktu:
 *
 *   1. Zatwierdź krok — `queue.action` resume / retry;
 *   2. Wstrzymaj — `queue.action` pause · `message.stop` · `session.stop`
 *      (ostatnia oddaje wykaz okien, w których turę zatrzymano);
 *   3. Nastaw koordynatora — `config.session.set` na zasięgu `window`
 *      z odczytem zwrotnym `config.effective.get`, który potwierdza źródło
 *      `override`;
 *   4. Przejmij sterowanie — złożenie trzech wywołań: `message.stop`
 *      + `config.session.set{permissions.mode=manual}` + `message.send`
 *      z poleceniem Operatora.
 *
 * Droga czwarta jest nazwana dokładnie tym, czym jest: „zatrzymuję turę,
 * przestawiam okno na pytanie o każdy krok, wpuszczam twoje polecenie”.
 * Dedykowanej komendy przejęcia kontrakt nie niesie; gdyby ją wniósł,
 * złożenie trzech wywołań zastąpi jedno, a ekran zostanie bez zmian.
 *
 * Droga pierwsza ma widoczną granicę. „Zatwierdź krok” jedzie przez kolejkę;
 * pozycja bez kolejki (sama zatrzymana pętla) nie ma czym zatwierdzić kroku,
 * bo komendy zatwierdzenia pojedynczego kroku kontrakt nie niesie.
 * `drogiDostepne` mówi to zdaniem, a ekran nie rysuje wtedy przycisku.
 *
 * Każde wywołanie oddaje kwit, nie `void` i nie `boolean`: nazwa komendy,
 * rozstrzygnięcie, stan po zmianie wyjęty z odpowiedzi rdzenia albo treść
 * odmowy. Kwit jest jedynym dowodem, jaki Operator zabiera z sobą, zamykając
 * telefon (`kwit-decyzji.ts`).
 */

/** Który byt wstrzymujemy — wybór wynika z pozycji, nie z domysłu. */
export type ZasiegWstrzymania = 'kolejka' | 'tura' | 'sesja';

/** Trzy gotowe rozstrzygnięcia nastawy koordynatora dosięgalne w dwóch dotknięciach. */
export type NastawaKoordynatora =
  | 'pytaj-o-kazdy-krok'
  | 'stoj-przy-braku-postepu'
  | 'obniz-naklad';

/** Co Operator robi z krokiem: puszcza dalej czy każe powtórzyć. */
export type WyborKroku = 'pusc-dalej' | 'powtorz';

/** Czy droga jest przejezdna dla tej pozycji; przy „nie” — powód zdaniem. */
export interface DostepnoscDrogi {
  dostepna: boolean;
  powod?: string;
}

/** Przejezdność czterech dróg dla jednej pozycji — podstawa rysowania przycisków. */
export interface DostepneDrogi {
  zatwierdz: DostepnoscDrogi;
  wstrzymaj: DostepnoscDrogi & { zasiegi: readonly ZasiegWstrzymania[] };
  nastaw: DostepnoscDrogi;
  przejmij: DostepnoscDrogi;
}

export interface WywolaniaInterwencji {
  /** Które drogi ta pozycja naprawdę niesie — i dlaczego pozostałe nie. */
  drogiDostepne(pozycja: PozycjaDecyzji): DostepneDrogi;
  zatwierdzKrok(pozycja: PozycjaDecyzji, wybor: WyborKroku): Promise<Kwit>;
  wstrzymaj(pozycja: PozycjaDecyzji, zasieg: ZasiegWstrzymania): Promise<Kwit>;
  nastawKoordynatora(pozycja: PozycjaDecyzji, nastawa: NastawaKoordynatora): Promise<Kwit>;
  przejmij(pozycja: PozycjaDecyzji, polecenie: string): Promise<Kwit>;
}

/** Etykiety nastaw — jedno źródło dla arkusza dróg i dla kwitu. */
export const ETYKIETY_NASTAW: Readonly<Record<NastawaKoordynatora, string>> = {
  'pytaj-o-kazdy-krok': 'pytaj mnie o każdy krok',
  'stoj-przy-braku-postepu': 'zatrzymaj bieg przy braku postępu',
  'obniz-naklad': 'obniż nakład rozumowania',
};

/** Etykiety zasięgów wstrzymania — to samo źródło dla ekranu i kwitu. */
export const ETYKIETY_ZASIEGOW: Readonly<Record<ZasiegWstrzymania, string>> = {
  kolejka: 'wstrzymaj kolejkę',
  tura: 'zatrzymaj bieżącą turę',
  sesja: 'zatrzymaj tury całej sesji',
};

/** Zdanie granicy drogi pierwszej — jedno źródło dla ekranu i dla kwitu. */
export const ZDANIE_BRAKU_ZATWIERDZENIA =
  'Krok zatwierdza się wyłącznie przez kolejkę (queue.action). Ta pozycja ' +
  'kolejki nie ma, a komendy zatwierdzenia pojedynczego kroku kontrakt nie niesie.';

export function utworzWywolaniaInterwencji(kanal: Kanal): WywolaniaInterwencji {
  /** Okno robocze pozycji: to, w którym biegnie tura. */
  const oknoPracy = (pozycja: PozycjaDecyzji): string | undefined => pozycja.windowId;

  /** Okno, na którym zapisujemy nastawę: koordynator, a gdy go nie ma — okno pracy. */
  const oknoNastawy = (pozycja: PozycjaDecyzji): string | undefined =>
    pozycja.koordynatorWindowId ?? pozycja.windowId;

  function zasiegiWstrzymania(pozycja: PozycjaDecyzji): ZasiegWstrzymania[] {
    const zasiegi: ZasiegWstrzymania[] = [];
    if (pozycja.queueId !== undefined) zasiegi.push('kolejka');
    if (oknoPracy(pozycja) !== undefined) zasiegi.push('tura');
    if (pozycja.sessionId !== undefined) zasiegi.push('sesja');
    return zasiegi;
  }

  return {
    drogiDostepne(pozycja) {
      const zasiegi = zasiegiWstrzymania(pozycja);
      return {
        zatwierdz:
          pozycja.queueId !== undefined
            ? { dostepna: true }
            : { dostepna: false, powod: ZDANIE_BRAKU_ZATWIERDZENIA },
        wstrzymaj: {
          dostepna: zasiegi.length > 0,
          zasiegi,
          ...(zasiegi.length > 0
            ? {}
            : {
                powod:
                  'Rdzeń nie podał przy tej pozycji ani kolejki, ani okna, ani sesji — ' +
                  'nie ma czego wstrzymać.',
              }),
        },
        nastaw:
          oknoNastawy(pozycja) !== undefined
            ? { dostepna: true }
            : {
                dostepna: false,
                powod:
                  'Nastawa zapisuje się na zasięgu okna (config.session.set, scope=window); ' +
                  'ta pozycja okna nie wskazuje.',
              },
        przejmij:
          oknoPracy(pozycja) !== undefined
            ? { dostepna: true }
            : {
                dostepna: false,
                powod:
                  'Przejęcie wpuszcza polecenie do okna (message.send); ' +
                  'ta pozycja okna nie wskazuje.',
              },
      };
    },

    async zatwierdzKrok(pozycja, wybor) {
      const queueId = pozycja.queueId;
      if (queueId === undefined) {
        return kwitOdmowy(Droga.Zatwierdz, pozycja, ZDANIE_BRAKU_ZATWIERDZENIA);
      }

      const action = wybor === 'powtorz' ? QueueAction.Retry : QueueAction.Resume;
      const wynik = await wywolaj(kanal, Command.QueueAction, { queueId, action });
      const krok = zapiszKrok(Command.QueueAction, wynik, (t) =>
        `kolejka ${t.queue.status}, obiegów ${t.queue.cycle ?? 0}`,
      );

      return zlozKwit(Droga.Zatwierdz, pozycja, [krok], (stany) =>
        wybor === 'powtorz'
          ? `Kazano powtórzyć krok kolejki. Rdzeń oddał: ${stany}.`
          : `Puszczono krok dalej. Rdzeń oddał: ${stany}.`,
      );
    },

    async wstrzymaj(pozycja, zasieg) {
      if (zasieg === 'kolejka') {
        const queueId = pozycja.queueId;
        if (queueId === undefined) {
          return kwitOdmowy(Droga.Wstrzymaj, pozycja, 'Ta pozycja nie stoi na żadnej kolejce.');
        }
        const wynik = await wywolaj(kanal, Command.QueueAction, {
          queueId,
          action: QueueAction.Pause,
        });
        const krok = zapiszKrok(Command.QueueAction, wynik, (t) => `kolejka ${t.queue.status}`);
        return zlozKwit(Droga.Wstrzymaj, pozycja, [krok], (stany) =>
          `Wstrzymano kolejkę. Rdzeń oddał: ${stany}.`,
        );
      }

      if (zasieg === 'tura') {
        const windowId = oknoPracy(pozycja);
        if (windowId === undefined) {
          return kwitOdmowy(Droga.Wstrzymaj, pozycja, 'Ta pozycja nie wskazuje okna.');
        }
        const wynik = await wywolaj(kanal, Command.MessageStop, { windowId });
        const krok = zapiszKrok(Command.MessageStop, wynik, (t) =>
          t.stopped ? `tura zatrzymana (wiadomość ${t.messageId})` : 'tura NIE została zatrzymana',
        );
        return zlozKwit(Droga.Wstrzymaj, pozycja, [krok], (stany) =>
          `Zatrzymano bieżącą turę. Rdzeń oddał: ${stany}.`,
        );
      }

      const sessionId = pozycja.sessionId;
      if (sessionId === undefined) {
        return kwitOdmowy(Droga.Wstrzymaj, pozycja, 'Ta pozycja nie wskazuje sesji.');
      }
      const wynik = await wywolaj(kanal, Command.SessionStop, { sessionId });
      const krok = zapiszKrok(Command.SessionStop, wynik, (t) =>
        `zatrzymano tur w oknach: ${t.stoppedWindowIds.length}`,
      );
      return zlozKwit(Droga.Wstrzymaj, pozycja, [krok], (stany) =>
        `Zatrzymano tury całej sesji. Rdzeń oddał: ${stany}.`,
      );
    },

    async nastawKoordynatora(pozycja, nastawa) {
      const windowId = oknoNastawy(pozycja);
      if (windowId === undefined) {
        return kwitOdmowy(
          Droga.Nastaw,
          pozycja,
          'Nastawa zapisuje się na zasięgu okna, a ta pozycja okna nie wskazuje.',
        );
      }

      const { areas, config } = ksztaltNastawy(nastawa);
      const zapis = await wywolaj(kanal, Command.ConfigSessionSet, {
        areas,
        config,
        scope: ConfigScope.Window,
        scopeId: windowId,
      });
      const krokZapisu = zapiszKrok(Command.ConfigSessionSet, zapis, (t) =>
        `zapisane obszary: ${t.storedAreas.join(', ')}`,
      );

      // Odczyt zwrotny potwierdza, że zapis wygrał rozstrzygnięcie poziomów
      // zasięgu, a nie tylko wylądował w bazie.
      const odczyt = await wywolaj(kanal, Command.ConfigEffectiveGet, { windowId, areas });
      const krokOdczytu = zapiszKrok(Command.ConfigEffectiveGet, odczyt, (t) =>
        opiszPochodzenie(t.effective.origins, areas),
      );

      return zlozKwit(
        Droga.Nastaw,
        pozycja,
        [krokZapisu, krokOdczytu],
        (stany) => `Nastawiono koordynatora: ${ETYKIETY_NASTAW[nastawa]}. Rdzeń oddał: ${stany}.`,
      );
    },

    async przejmij(pozycja, polecenie) {
      const windowId = oknoPracy(pozycja);
      if (windowId === undefined) {
        return kwitOdmowy(Droga.Przejmij, pozycja, 'Ta pozycja nie wskazuje okna.');
      }
      const tresc = polecenie.trim();
      if (tresc === '') {
        return kwitOdmowy(
          Droga.Przejmij,
          pozycja,
          'Przejęcie nie ruszyło: polecenia dla okna nie podano, a puste polecenie nie jest przejęciem.',
        );
      }

      const kroki: KrokKwitu[] = [];

      const stop = await wywolaj(kanal, Command.MessageStop, { windowId });
      kroki.push(
        zapiszKrok(Command.MessageStop, stop, (t) =>
          t.stopped ? 'tura zatrzymana' : 'tura nie biegła albo nie dała się zatrzymać',
        ),
      );

      const nastawa = await wywolaj(kanal, Command.ConfigSessionSet, {
        areas: [SessionConfigArea.Permissions],
        config: { permissions: { mode: PermissionMode.Manual } },
        scope: ConfigScope.Window,
        scopeId: windowId,
      });
      kroki.push(
        zapiszKrok(Command.ConfigSessionSet, nastawa, (t) =>
          `tryb uprawnień okna: ${t.config.permissions?.mode ?? PermissionMode.Manual}`,
        ),
      );

      const wiadomosc = await wywolaj(kanal, Command.MessageSend, { windowId, content: tresc });
      kroki.push(
        zapiszKrok(Command.MessageSend, wiadomosc, (t) =>
          `polecenie przyjęte jako wiadomość ${t.message.id} (${t.message.status})`,
        ),
      );

      return zlozKwit(
        Droga.Przejmij,
        pozycja,
        kroki,
        (stany) =>
          'Przejęto sterowanie: zatrzymano turę, przestawiono okno na pytanie o każdy krok, ' +
          `wpuszczono polecenie. Rdzeń oddał: ${stany}.`,
      );
    },
  };
}

/** Obszary i treść zapisu dla trzech gotowych rozstrzygnięć nastawy. */
export function ksztaltNastawy(nastawa: NastawaKoordynatora): {
  areas: SessionConfigArea[];
  config: SessionConfig;
} {
  if (nastawa === 'pytaj-o-kazdy-krok') {
    return {
      areas: [SessionConfigArea.Permissions],
      config: { permissions: { mode: PermissionMode.Manual } },
    };
  }
  if (nastawa === 'stoj-przy-braku-postepu') {
    return {
      areas: [SessionConfigArea.SessionLifecycle],
      config: { sessionLifecycle: { stopOnNoProgress: true } },
    };
  }
  return {
    areas: [SessionConfigArea.Model],
    config: { model: { reasoningEffort: ReasoningEffort.Low } },
  };
}

/** Zdanie o pochodzeniu obszarów — dowód, że zapis wygrał rozstrzygnięcie zasięgów. */
function opiszPochodzenie(
  origins: readonly SessionConfigOrigin[],
  areas: readonly SessionConfigArea[],
): string {
  const nasze = origins.filter((o) => areas.includes(o.area));
  if (nasze.length === 0) return 'rdzeń nie podał pochodzenia obszaru';
  return nasze
    .map((o) => `${o.area}: źródło ${o.source}${o.scope === undefined ? '' : `, zasięg ${o.scope}`}`)
    .join('; ');
}

/** Zapisuje jedno wywołanie jako krok kwitu — stan po zmianie albo odmowa. */
function zapiszKrok<T>(
  komenda: string,
  wynik: Wynik<T>,
  stanPo: (tresc: T) => string,
): KrokKwitu {
  if (!wynik.udany || wynik.wynik === undefined) {
    return { komenda, udany: false, odmowa: opisOdmowyBledu(komenda, wynik.blad) };
  }
  return { komenda, udany: true, stanPo: stanPo(wynik.wynik) };
}

/** Składa kwit z kroków; zdanie powstaje ze stanów oddanych przez rdzeń. */
function zlozKwit(
  droga: Droga,
  pozycja: PozycjaDecyzji,
  kroki: readonly KrokKwitu[],
  zdanieUdane: (stany: string) => string,
): Kwit {
  const udany = kroki.every((k) => k.udany);
  const stany = kroki
    .map((k) => (k.udany ? (k.stanPo ?? k.komenda) : (k.odmowa ?? k.komenda)))
    .join('; ');

  return {
    id: `${droga}:${pozycja.id}:${Date.now()}`,
    droga,
    pozycja: pozycja.naglowek,
    kroki,
    udany,
    zdanie: udany
      ? zdanieUdane(stany)
      : `Rdzeń nie wykonał drogi do końca. Odpowiedź rdzenia: ${stany}.`,
    o: Date.now(),
  };
}

/** Kwit drogi, która nie ruszyła — bez ani jednego wywołania do rdzenia. */
function kwitOdmowy(droga: Droga, pozycja: PozycjaDecyzji, zdanie: string): Kwit {
  return {
    id: `${droga}:${pozycja.id}:${Date.now()}`,
    droga,
    pozycja: pozycja.naglowek,
    kroki: [],
    udany: false,
    zdanie,
    o: Date.now(),
  };
}
