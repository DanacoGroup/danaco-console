import { nazwaKanalu, type RejestrKanalow } from './rejestr-kanalow';
import { nazwaAgenta, type RejestrAgentow } from './rejestr-agentow';
import { utworzListeWyboru, type ListaWyboru, type OpcjaWyboru } from './lista-wyboru';
import type { StanSterowania } from './stan-sterowania';
import type { ModelKartySesji } from './model-karty-sesji';
import type { ZmianaOkna } from './zmiana-okna';

const NAZWA = 'Model';

/** Nazwy dwóch sekcji jednego menu wyboru. */
const SEKCJA_MODELE = 'Modele';
const SEKCJA_AGENCI = 'Moi agenci';

/**
 * Przedrostek wartości oznaczającej eksperta.
 *
 * Lista wyboru niesie jedną wartość na pozycję, a menu ma dwa rodzaje pozycji.
 * Przedrostek rozstrzyga rodzaj bez drugiego pola i bez zgadywania po
 * identyfikatorze: kod agenta i kod kanału pochodzą z dwóch rejestrów, więc
 * mogą się zderzyć.
 */
export const PRZEDROSTEK_AGENTA = 'agent:';

/**
 * Sterowanie modelem obsługującym okno — kanał surowy albo agent.
 *
 * Kanał jest parametrem okna, nie sesji: dwa okna jednej sesji pracują na
 * dwóch różnych modelach równocześnie. Wykaz pochodzi z rejestru kanałów
 * — nowy model to nowy wiersz rejestru, nie zmiana w tym pliku.
 *
 * Menu ma dwie sekcje i jeden wybór: modele surowe w sekcji „Modele", agenci
 * w „Moi agenci". Okno obsługuje albo model, albo agenta nałożonego na model.
 * Lista płaska zatarłaby tę zależność i przy licznych agentach utopiłaby między
 * nimi modele surowe.
 *
 * Agent nie jest kanałem: kanał to droga do modelu, agent to tożsamość nałożona
 * na tę drogę. Wybór agenta niesie więc jego kod w polu `agentId`, a pole
 * `modelChannelId` bierze kanał bazowy agenta — model, na którym agent stoi.
 *
 * Kanał oznaczony jako nieczynny zostaje na liście i pozostaje wybieralny;
 * o jego stanie mówi nazwa pozycji.
 */
export function utworzSterowanieModelu(
  stan: StanSterowania,
  zmiana: ZmianaOkna,
  rejestr: RejestrKanalow,
  rejestrAgentow: RejestrAgentow,
  kartaSesji?: ModelKartySesji,
): HTMLElement {
  // Zawężanie włączone tylko tutaj: to jedyny wykaz sterowania, który rośnie
  // wraz z liczbą kanałów i agentów. Rola okna czy tryb uprawnień mają pozycji
  // tyle, ile ma ich kontrakt, więc pole zawężania byłoby tam kontrolką nad
  // niczym.
  const lista = utworzListeWyboru(
    NAZWA,
    (wartosc) => {
      const zlecenie = zlecenieWyboru(wartosc, rejestrAgentow);
      zmiana.zastosuj(NAZWA, zlecenie);
      // Ten sam wybór, jeden przekład: rozesłanie na kartę bierze zlecenie
      // złożone tutaj, a nie składa go drugi raz po swojemu.
      kartaSesji?.ustawWybor(zlecenie);
    },
    true,
  );

  function odrysuj(): void {
    const okno = stan.migawka().okno;
    lista.pokaz(opcje(rejestr, rejestrAgentow), wartoscBiezaca(okno.agentId, okno.modelChannelId));
    naniesStanRejestru(lista, rejestr);
  }

  stan.naZmiane(odrysuj);
  rejestr.naZmiane(odrysuj);
  rejestrAgentow.naZmiane(odrysuj);
  rejestrAgentow.odswiez();
  odrysuj();

  if (kartaSesji === undefined) return lista.element;

  // Wybór dotyczy okna; rozesłanie na kartę sesji jest osobną czynnością i stoi
  // pod listą. Jedna wspólna kontrolka nie pokazywałaby, ile okien obejmuje
  // zmiana.
  const blok = document.createElement('div');
  blok.className = 'dn-sterowanie__blok';
  blok.append(lista.element, kartaSesji.element);
  return blok;
}

/**
 * Katalog wyboru: kanały rdzenia, a po nich agenci.
 *
 * Kolejność jest treścią, nie porządkiem alfabetycznym: modele surowe stoją
 * pierwsze, bo agent bez modelu nie istnieje. Sekcja bez pozycji nie powstaje,
 * więc pusty rejestr agentów nie tworzy nagłówka nad pustką.
 */
export function opcje(rejestr: RejestrKanalow, rejestrAgentow?: RejestrAgentow): OpcjaWyboru[] {
  const kanaly: OpcjaWyboru[] = rejestr.kanaly().map((kanal) => ({
    wartosc: kanal.id,
    nazwa: nazwaKanalu(kanal),
    sekcja: SEKCJA_MODELE,
  }));
  const agenci: OpcjaWyboru[] = (rejestrAgentow?.agenci() ?? []).map((agent) => ({
    wartosc: `${PRZEDROSTEK_AGENTA}${agent.id}`,
    nazwa: nazwaAgenta(agent),
    sekcja: SEKCJA_AGENCI,
  }));
  return [...kanaly, ...agenci];
}

/**
 * Wartość zaznaczona na liście: agent wyprzedza kanał.
 *
 * Okno z nałożonym agentem niesie oba pola — `agentId` i kanał bazowy tego
 * agenta — więc zaznaczenie kanału pokazywałoby model surowy tam, gdzie pracuje
 * agent.
 */
export function wartoscBiezaca(idAgenta: string | undefined, idKanalu: string): string {
  return idAgenta !== undefined && idAgenta.length > 0
    ? `${PRZEDROSTEK_AGENTA}${idAgenta}`
    : idKanalu;
}

/**
 * Nazwa modelu obsługującego okno — jeden przekład dla wszystkich czytelników.
 *
 * Kolumna sterowania, pasek zlecenia i wiersz „Model" podsumowania szuflady
 * czytają tę samą funkcję, więc pokazują tę samą nastawę.
 *
 * Rejestr agentów jest nieobowiązkowy. Czytelnik, który go nie ma, nie zobaczy
 * modelu surowego zamiast agenta: dostanie kod agenta z rzeczownikiem, a nie
 * podmieniony kanał.
 *
 * Identyfikator spoza obu rejestrów zostaje pokazany dosłownie, nigdy jako
 * „Bez wskazania" — wskazanie jest, tylko wykaz go nie zna.
 */
export function opisWyboruModelu(
  rejestr: RejestrKanalow,
  rejestrAgentow: RejestrAgentow | undefined,
  idKanalu: string,
  idAgenta: string | undefined,
  gdyPusty: string,
): string {
  const kodAgenta = idAgenta ?? '';
  if (kodAgenta !== '') {
    const agent = rejestrAgentow?.agenci().find((pozycja) => pozycja.id === kodAgenta);
    return agent !== undefined ? nazwaAgenta(agent) : `${kodAgenta} (ekspert)`;
  }
  if (idKanalu === '') return gdyPusty;
  const wykaz = rejestr.kanaly();
  // Wykaz pusty znaczy „odpowiedź rdzenia w drodze", nie „kanał nieznany" —
  // dopisek o braku w wykazie byłby wtedy wyrokiem bez dowodu.
  if (wykaz.length === 0) return idKanalu;
  const kanal = wykaz.find((pozycja) => pozycja.id === idKanalu);
  return kanal !== undefined ? nazwaKanalu(kanal) : `${idKanalu} (spoza wykazu)`;
}

/**
 * Zlecenie zmiany okna wyprowadzone z wybranej pozycji.
 *
 * Wybór modelu surowego zdejmuje agenta pustym `agentId`; bez tego agent
 * zostawałby nałożony na model, który nie jest już wybrany, a wykaz pokazywałby
 * co innego niż stan okna.
 *
 * Agent bez kanału bazowego (`channelId` pusty — agent założony, model jeszcze
 * nie przypisany) nie jest odmową: idzie sam `agentId`, a kanał zostaje ten,
 * który okno miało.
 */
export function zlecenieWyboru(
  wartosc: string,
  rejestrAgentow: RejestrAgentow,
): { modelChannelId?: string; agentId: string } {
  if (!wartosc.startsWith(PRZEDROSTEK_AGENTA)) {
    return { modelChannelId: wartosc, agentId: '' };
  }
  const kod = wartosc.slice(PRZEDROSTEK_AGENTA.length);
  const agent = rejestrAgentow.agenci().find((pozycja) => pozycja.id === kod);
  const kanalBazowy = agent?.channelId ?? '';
  return kanalBazowy === ''
    ? { agentId: kod }
    : { modelChannelId: kanalBazowy, agentId: kod };
}

/**
 * Trzy stany rejestru kanałów naniesione na pole: odczyt, wykaz pusty, wykaz
 * wypełniony.
 *
 * Rozróżnienie „rdzeń jeszcze nie odpowiedział" od „rejestr jest pusty" jest
 * konieczne: bez niego wskaźnik odczytu nigdy nie zgasłby na rdzeniu z pustym
 * rejestrem, zapowiadając wykaz, który nie nadejdzie.
 */
export function naniesStanRejestru(lista: ListaWyboru, rejestr: RejestrKanalow): void {
  const odpowiedziano = rejestr.odpowiedzOtrzymana();
  lista.ustawOdczyt(!odpowiedziano);
  lista.ustawUwage(
    odpowiedziano && rejestr.kanaly().length === 0 ? REJESTR_PUSTY : '',
  );
}

/**
 * Zdanie stanu pustego. Kanał zakłada się komendą `channel.add`; okna do
 * zakładania kanałów w tym wydaniu nie ma, więc zdanie mówi to wprost zamiast
 * odsyłać do nieistniejącego ekranu.
 */
const REJESTR_PUSTY =
  'Rejestr kanałów rdzenia jest pusty — rdzeń odpowiedział, ale nie zna ani jednego kanału. Kanał zakłada się komendą channel.add; okna do jego zakładania w tym wydaniu nie ma.';
