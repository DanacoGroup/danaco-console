import { nazwaKanalu, type RejestrKanalow } from './rejestr-kanalow';
import { nazwaAgenta, type RejestrAgentow } from './rejestr-agentow';
import { utworzListeWyboru, type ListaWyboru, type OpcjaWyboru } from './lista-wyboru';
import type { StanSterowania } from './stan-sterowania';
import type { ModelKartySesji } from './model-karty-sesji';
import type { ZmianaOkna } from './zmiana-okna';

const NAZWA = 'Model';

/** Nazwy dwóch sekcji jednego menu wyboru modelu: sekcji modeli surowych oraz sekcji agentów operatora okna. */
const SEKCJA_MODELE = 'Modele';
const SEKCJA_AGENCI = 'Moi agenci';

/** Przedrostek wartości oznaczającej eksperta, rozstrzygający rodzaj pozycji menu bez drugiego pola i bez zgadywania po identyfikatorze. */
export const PRZEDROSTEK_AGENTA = 'agent:';

/** Sterowanie modelem obsługującym okno: kanał surowy albo agent, wybierane z jednego menu dwusekcyjnego. */
export function utworzSterowanieModelu(
  stan: StanSterowania,
  zmiana: ZmianaOkna,
  rejestr: RejestrKanalow,
  rejestrAgentow: RejestrAgentow,
  kartaSesji?: ModelKartySesji,
): HTMLElement {
  // Zawężanie włączone tylko tutaj: jedyny wykaz sterowania rosnący z liczbą kanałów i agentów.
  const lista = utworzListeWyboru(
    NAZWA,
    (wartosc) => {
      const zlecenie = zlecenieWyboru(wartosc, rejestrAgentow);
      zmiana.zastosuj(NAZWA, zlecenie);
      // Ten sam wybór, jeden przekład: rozesłanie na kartę bierze zlecenie złożone tutaj.
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

  // Wybór dotyczy okna; rozesłanie na kartę sesji jest osobną czynnością pod listą.
  const blok = document.createElement('div');
  blok.className = 'dn-sterowanie__blok';
  blok.append(lista.element, kartaSesji.element);
  return blok;
}

/** Katalog wyboru: kanały rdzenia, a po nich agenci, ułożeni w kolejności treściowej, nie porządku alfabetycznego. */
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

/** Nazwa modelu obsługującego okno, jeden przekład czytany przez kolumnę sterowania, pasek zlecenia i podsumowanie szuflady. */
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
  // Wykaz pusty znaczy odpowiedź rdzenia w drodze, nie kanał nieznany.
  if (wykaz.length === 0) return idKanalu;
  const kanal = wykaz.find((pozycja) => pozycja.id === idKanalu);
  return kanal !== undefined ? nazwaKanalu(kanal) : `${idKanalu} (spoza wykazu)`;
}

/** Zlecenie zmiany okna wyprowadzone z wybranej pozycji menu: modelu surowego albo agenta wraz z jego kanałem bazowym. */
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

/** Trzy stany rejestru kanałów naniesione na pole listy: odczyt w toku, wykaz pusty i wykaz wypełniony pozycjami. */
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
