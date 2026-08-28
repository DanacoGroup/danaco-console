import {
  Command,
  EnvelopeStatus,
  KOMENDY,
  PROTOCOL_VERSION,
  ZDARZENIA,
  type ConnectionHelloResponse,
} from '../../../shared/contract';
import { opisOdmowy } from '../komponenty/odmowa';
import type { Kanal } from '../protokol/kanal';
import { czyTablica, sprawdzKsztalt } from '../protokol/ksztalt-odpowiedzi';
import { wywolaj } from '../protokol/wywolanie';

/**
 * Wykaz komend rdzenia ustala dla każdego połączenia jeden fakt i jedno zdanie o nim: rozstrzygnięcie o pojedynczej komendzie ujmuje pięć stanów pokrycia, a wartość ta sama trafia do atrybutu pokrycia znacznika.
 */
export type StanPokrycia =
  | 'nieustalone'
  | 'brak-w-kontrakcie'
  | 'zdarzenie-nie-komenda'
  | 'rdzen-nie-ma'
  | 'rdzen-ma';

/** Komendy i zdarzenia kontraktu — wykazy stałych pochodzące wprost z definicji kontraktu współdzielonego między rdzeniem a klientem. */
const KOMENDY_KONTRAKTU: ReadonlySet<string> = new Set<string>(KOMENDY);
const ZDARZENIA_KONTRAKTU: ReadonlySet<string> = new Set<string>(ZDARZENIA);

/** Odczyt wykazu komend wraz z ewentualną odmową rdzenia; wartość null uchwytów znaczy, że rdzeń jeszcze nie orzekł o żadnej komendzie. */
interface Odczyt {
  uchwyty: ReadonlySet<string> | null;
  /** Zdanie odmowy powitania; puste, dopóki nic nie odmówiło. */
  odmowa: string;
}

/** Pamięć podręczna jednego kanału łącząca odczyt wykazu komend z kontrolkami modułów, które trzeba przerysować po jego zmianie. */
interface WpisPamieci {
  odczyt: Odczyt;
  /** Powitanie w drodze — wszystkie moduły czekają na jedną odpowiedź. */
  wToku: Promise<void> | null;
  zalezni: Set<() => void>;
}

const pamiec = new WeakMap<Kanal, WpisPamieci>();

/** Podpięcie kontrolki modułu pod wykaz komend danego kanału; wywołanie zwrotne oddaje funkcję odpinającą ten nasłuch. */
export function podepnijDoWykazu(kanal: Kanal, przerysuj: () => void): () => void {
  const wspolny = wpisPamieci(kanal);
  wspolny.zalezni.add(przerysuj);
  return () => wspolny.zalezni.delete(przerysuj);
}

/** Pyta rdzeń o wykaz komend tylko wtedy, gdy dla tego kanału nikt jeszcze nie zapytał ani odpowiedź nie jest już znana. */
export async function zapewnijOdczyt(kanal: Kanal): Promise<void> {
  const wspolny = wpisPamieci(kanal);
  if (wspolny.odczyt.uchwyty !== null) return;
  if (wspolny.wToku !== null) return wspolny.wToku;
  const bieg = powitaj(kanal, wspolny);
  wspolny.wToku = bieg;
  await bieg;
  wspolny.wToku = null;
}

/**
 * Zapomnienie wykazu — po zerwaniu połączenia albo wymianie rdzenia.
 *
 * Kontrolki wracają do zdania „odczyt w toku”, a nie do ostatniego orzeczenia:
 * wykaz rdzenia, którego już nie ma, nie jest wiedzą o rdzeniu, który przyjdzie.
 */
export function zapomnijWykazKomend(kanal: Kanal): void {
  zapisz(wpisPamieci(kanal), { uchwyty: null, odmowa: '' });
}

/** Rozstrzygnięcie o pojedynczej komendzie kanału; dwa pierwsze możliwe stany rozstrzyga sam kontrakt, bez pytania rdzenia. */
export function stanKomendy(kanal: Kanal, komenda: string): StanPokrycia {
  if (!KOMENDY_KONTRAKTU.has(komenda)) {
    return ZDARZENIA_KONTRAKTU.has(komenda) ? 'zdarzenie-nie-komenda' : 'brak-w-kontrakcie';
  }
  const { uchwyty } = wpisPamieci(kanal).odczyt;
  if (uchwyty === null) return 'nieustalone';
  return uchwyty.has(komenda) ? 'rdzen-ma' : 'rdzen-nie-ma';
}

/** Waga poszczególnych rozstrzygnięć uporządkowana od najcięższego; tą kolejnością rozstrzyga się stan pozycji złożonej z kilku komend. */
const WAGA_STANU: readonly StanPokrycia[] = [
  'brak-w-kontrakcie',
  'zdarzenie-nie-komenda',
  'rdzen-nie-ma',
  'nieustalone',
  'rdzen-ma',
];

/** Rozstrzygnięcie zbiorcze pozycji wykazu — najcięższe spośród rozstrzygnięć wszystkich komend, które ta pozycja obejmuje. */
export function stanPozycji(kanal: Kanal, komendy: readonly string[]): StanPokrycia {
  // Pozycja bez ani jednej komendy nie ma czego rozstrzygać — i nie udaje, że ma.
  if (komendy.length === 0) return 'nieustalone';
  let wybor: StanPokrycia = 'rdzen-ma';
  for (const komenda of komendy) {
    const stan = stanKomendy(kanal, komenda);
    if (WAGA_STANU.indexOf(stan) < WAGA_STANU.indexOf(wybor)) wybor = stan;
  }
  return wybor;
}

/** Powód rozstrzygnięcia dotyczącego jednej komendy, zapisany zdaniem samodzielnym gotowym do pokazania osobno. */
export function powodKomendy(kanal: Kanal, komenda: string): string {
  return zDuzej(powodMala(kanal, komenda));
}

/** Powód dotyczący całej pozycji wykazu — zdania poszczególnych jej komend, każde poprzedzone nazwą własnej czynności. */
export function zdanieOPozycji(
  kanal: Kanal,
  komendy: readonly string[],
  czynnosc: string,
): string {
  if (komendy.length === 0) return POKRYCIE_W_ODCZYCIE;
  const powody = komendy.map((komenda) => powodMala(kanal, komenda)).join(' ');
  const nazwa = czynnosc.trim();
  return nazwa === '' ? zDuzej(powody) : `${nazwa}: ${powody}`;
}

/** Zdanie opisujące stan nieustalony bez podawania nazwy komendy; używane przy rozstrzygnięciach pozycji zbiorczych. */
export const POKRYCIE_W_ODCZYCIE =
  'Pokrycie tej pozycji w rdzeniu — odczyt w toku. Powód pojawi się po odpowiedzi rdzenia.';

/** Nagłówek całego wykazu komend; braków wynikających z ciszy rdzenia nie liczy, bo liczba podana przed odpowiedzią byłaby przedwczesnym orzeczeniem. */
export function naglowekWykazu(kanal: Kanal, komendy: readonly string[]): string {
  const { uchwyty, odmowa } = wpisPamieci(kanal).odczyt;
  if (odmowa !== '') return 'Pokrycie komend tego okna w rdzeniu — nieustalone';
  if (uchwyty === null) return 'Pokrycie komend tego okna w rdzeniu — odczyt w toku…';
  const bezDrogi = komendy.filter((k) => stanKomendy(kanal, k) !== 'rdzen-ma').length;
  return `Komendy tego okna bez drogi do rdzenia (${bezDrogi} z ${komendy.length})`;
}

/** Wpis pamięci podręcznej jednego kanału wraz z nasłuchem powitań rdzenia, zakładany dokładnie raz na każdy kanał. */
function wpisPamieci(kanal: Kanal): WpisPamieci {
  const znany = pamiec.get(kanal);
  if (znany !== undefined) return znany;
  const swiezy: WpisPamieci = { odczyt: { uchwyty: null, odmowa: '' }, wToku: null, zalezni: new Set() };
  pamiec.set(kanal, swiezy);
  // Powitanie cudze liczy się jak własne — uzgodnienie wita rdzeń przy każdym połączeniu.
  kanal.naDowolny((koperta) => {
    if (koperta.type !== Command.ConnectionHello || koperta.status !== EnvelopeStatus.Ok) return;
    const komendy = (koperta.payload as ConnectionHelloResponse | undefined)?.commands;
    if (!czyTablica(komendy)) return;
    zapisz(swiezy, { uchwyty: new Set(komendy), odmowa: '' });
  });
  return swiezy;
}

/** Obsługa jednego powitania rdzenia i zapis jego wyniku we wpisie kanału — zarówno wyniku udanego, jak i odmowy. */
async function powitaj(kanal: Kanal, wspolny: WpisPamieci): Promise<void> {
  const wynik = sprawdzKsztalt(
    await wywolaj(kanal, Command.ConnectionHello, {
      // Tożsamość jest stała, nie świeża, bo obsługiwacz powitania nie czyta żądania.
      clientId: 'pokrycie-komend',
      clientVersion: '1.0',
      protocolVersion: PROTOCOL_VERSION,
    }),
    Command.ConnectionHello,
    (tresc) => czyTablica(tresc.commands),
  );
  if (!wynik.udany || wynik.wynik === undefined) {
    // Odmowa nie wypełnia `uchwyty`, więc następne `zapewnijOdczyt` ponowi pytanie.
    const odmowa = opisOdmowy('Odczyt wykazu komend rdzenia', wynik.blad?.code, wynik.blad?.message);
    zapisz(wspolny, { uchwyty: null, odmowa });
    return;
  }
  zapisz(wspolny, { uchwyty: new Set(wynik.wynik.commands ?? []), odmowa: '' });
}

/** Zapis wyniku odczytu we wpisie kanału i przerysowanie kontrolek wszystkich modułów, które na ten kanał nasłuchują. */
function zapisz(wspolny: WpisPamieci, odczyt: Odczyt): void {
  wspolny.odczyt = odczyt;
  for (const przerysuj of wspolny.zalezni) przerysuj();
}

/** Powód dotyczący jednej komendy zapisany zdaniem podrzędnym z małej litery, gotowym do złożenia z nazwą czynności. */
function powodMala(kanal: Kanal, komenda: string): string {
  switch (stanKomendy(kanal, komenda)) {
    case 'brak-w-kontrakcie':
      return `kontrakt nie ma komendy ${komenda} — nie ma jej ani rdzeń, ani żaden klient, więc czynność wymaga najpierw pozycji w kontrakcie (dziś ${KOMENDY.length} komend).`;
    case 'zdarzenie-nie-komenda':
      return `${komenda} nie jest komendą kontraktu, a jego ZDARZENIEM — tej pozycji nie wywołuje się wcale, można jej wyłącznie nasłuchiwać.`;
    case 'rdzen-nie-ma':
      return `rdzeń nie ma uchwytu komendy ${komenda} — kontrakt ją niesie, więc brak jest po stronie rdzenia; czynność nie ma drogi do rdzenia i nie zostaje wysłana.`;
    case 'rdzen-ma':
      return `rdzeń ma uchwyt komendy ${komenda}, ale to okno jeszcze go nie wywołuje — brak jest po stronie okna, nie po stronie rdzenia.`;
    default:
      return wpisPamieci(kanal).odczyt.odmowa !== ''
        ? `pokrycie komendy ${komenda} w rdzeniu nieustalone — ${wpisPamieci(kanal).odczyt.odmowa}; okno nie orzeka o braku, którego rdzeń nie orzekł.`
        : `pokrycie komendy ${komenda} w rdzeniu — odczyt w toku; powód pojawi się po odpowiedzi rdzenia.`;
  }
}

/** Zamienia pierwszą literę zdania na wielką, bo powody poszczególnych komend układa się konsekwentnie z litery małej. */
function zDuzej(zdanie: string): string {
  return zdanie.charAt(0).toUpperCase() + zdanie.slice(1);
}
