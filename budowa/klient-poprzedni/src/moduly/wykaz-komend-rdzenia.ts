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
 * Wykaz komend rdzenia — fakt i zdanie o fakcie, jeden na połączenie.
 *
 * Warstwa niższa bytu `pokrycie-komend.ts`: nie zna dokumentu ani kontrolki,
 * zna tylko odpowiedź rdzenia i zdania, które z niej wynikają.
 *
 * Prawda bierze się z powitania `connection.hello`: rdzeń oddaje w polu
 * `commands` wykaz komend, które naprawdę obsługuje, a nie wykaz z kontraktu
 * (`core/handlers_connection.go` → `Rejestr.Nazwy()`). Obsługiwacz powitania
 * nie czyta żądania i niczego nie zapisuje, więc powtórne powitanie jest
 * czystym odczytem — nie drugim uzgodnieniem połączenia i nie rejestracją
 * drugiego klienta. Samych komend nie pytamy: wywołanie
 * `orchestration.dependency.remove`, żeby zobaczyć odmowę, byłoby wykonaniem
 * czynności niszczącej dla samego sprawdzenia.
 *
 * Rozstrzygnięć jest pięć (`StanPokrycia` niżej) i każde znaczy co innego;
 * dwa pierwsze rozstrzyga sam kontrakt, bez pytania rdzenia, więc widać je
 * jeszcze przed odpowiedzią powitania.
 *
 * Jedno powitanie na połączenie, nie jedno na moduł: wykaz leży w pamięci
 * podręcznej przypisanej do kanału (`WeakMap`), więc pierwsze `zapewnijOdczyt`
 * pyta, a pozostałe czekają na tę samą odpowiedź. Odmowa w pamięci nie zostaje
 * — następne wywołanie ponawia pytanie, żeby jedna nieudana chwila nie zatruła
 * zdań do końca sesji.
 *
 * Unieważnienie: wykaz zmienia się z wersją rdzenia, nie w toku sesji, ale
 * transport ponawia połączenie pod tym samym kanałem, więc po zerwaniu można
 * trafić na rdzeń inny niż odczytany. Dlatego odpowiedzi `connection.hello`
 * nasłuchujemy na całym kanale: każde powitanie, także cudze — choćby
 * uzgodnienie powłoki po ponowieniu — odświeża wykaz bez ani jednego
 * zapytania. Powłoka może wymusić zapomnienie wprost przez
 * `zapomnijWykazKomend(kanal)`.
 */

/**
 * Rozstrzygnięcie o jednej komendzie; ta sama wartość idzie w `data-pokrycie`.
 *
 * `nieustalone`           rdzeń nie odpowiedział albo odmówił — okno nie orzeka
 *                         wtedy o braku, bo cisza nie jest orzeczeniem,
 * `brak-w-kontrakcie`     nazwy nie ma w kontrakcie w ogóle (wykazy okien
 *                         obiecują więcej, niż kontrakt niesie),
 * `zdarzenie-nie-komenda` nazwa jest, ale jest zdarzeniem — tego się nie wywołuje,
 * `rdzen-nie-ma`          kontrakt komendę ma, rdzeń nie ma uchwytu,
 * `rdzen-ma`              rdzeń ma uchwyt, a to okno go jeszcze nie wywołuje.
 */
export type StanPokrycia =
  | 'nieustalone'
  | 'brak-w-kontrakcie'
  | 'zdarzenie-nie-komenda'
  | 'rdzen-nie-ma'
  | 'rdzen-ma';

/** Komendy i zdarzenia kontraktu — wykazy z `shared/contract.ts`. */
const KOMENDY_KONTRAKTU: ReadonlySet<string> = new Set<string>(KOMENDY);
const ZDARZENIA_KONTRAKTU: ReadonlySet<string> = new Set<string>(ZDARZENIA);

/** Odczyt wykazu; `uchwyty === null` znaczy „rdzeń jeszcze nie orzekł”. */
interface Odczyt {
  uchwyty: ReadonlySet<string> | null;
  /** Zdanie odmowy powitania; puste, dopóki nic nie odmówiło. */
  odmowa: string;
}

/** Pamięć podręczna jednego kanału wraz z kontrolkami do przerysowania. */
interface WpisPamieci {
  odczyt: Odczyt;
  /** Powitanie w drodze — wszystkie moduły czekają na jedną odpowiedź. */
  wToku: Promise<void> | null;
  zalezni: Set<() => void>;
}

const pamiec = new WeakMap<Kanal, WpisPamieci>();

/** Podpięcie kontrolki pod wykaz kanału; oddaje odpięcie. */
export function podepnijDoWykazu(kanal: Kanal, przerysuj: () => void): () => void {
  const wspolny = wpisPamieci(kanal);
  wspolny.zalezni.add(przerysuj);
  return () => wspolny.zalezni.delete(przerysuj);
}

/** Pyta rdzeń o wykaz komend, jeśli nikt jeszcze nie zapytał ani nie wie. */
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

/** Rozstrzygnięcie o komendzie. Dwa pierwsze stany nie wymagają rdzenia. */
export function stanKomendy(kanal: Kanal, komenda: string): StanPokrycia {
  if (!KOMENDY_KONTRAKTU.has(komenda)) {
    return ZDARZENIA_KONTRAKTU.has(komenda) ? 'zdarzenie-nie-komenda' : 'brak-w-kontrakcie';
  }
  const { uchwyty } = wpisPamieci(kanal).odczyt;
  if (uchwyty === null) return 'nieustalone';
  return uchwyty.has(komenda) ? 'rdzen-ma' : 'rdzen-nie-ma';
}

/** Waga rozstrzygnięć od najcięższego — tak rozstrzyga pozycja o kilku komendach. */
const WAGA_STANU: readonly StanPokrycia[] = [
  'brak-w-kontrakcie',
  'zdarzenie-nie-komenda',
  'rdzen-nie-ma',
  'nieustalone',
  'rdzen-ma',
];

/** Rozstrzygnięcie zbiorcze pozycji: najcięższe z rozstrzygnięć jej komend. */
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

/** Powód dotyczący jednej komendy, zdaniem samodzielnym. */
export function powodKomendy(kanal: Kanal, komenda: string): string {
  return zDuzej(powodMala(kanal, komenda));
}

/** Powód całej pozycji: zdania jej komend poprzedzone nazwą czynności. */
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

/** Zdanie stanu nieustalonego bez nazwy komendy — dla pozycji zbiorczych. */
export const POKRYCIE_W_ODCZYCIE =
  'Pokrycie tej pozycji w rdzeniu — odczyt w toku. Powód pojawi się po odpowiedzi rdzenia.';

/** Nagłówek wykazu. Braków z ciszy nie liczy: liczba przed odpowiedzią byłaby orzeczeniem. */
export function naglowekWykazu(kanal: Kanal, komendy: readonly string[]): string {
  const { uchwyty, odmowa } = wpisPamieci(kanal).odczyt;
  if (odmowa !== '') return 'Pokrycie komend tego okna w rdzeniu — nieustalone';
  if (uchwyty === null) return 'Pokrycie komend tego okna w rdzeniu — odczyt w toku…';
  const bezDrogi = komendy.filter((k) => stanKomendy(kanal, k) !== 'rdzen-ma').length;
  return `Komendy tego okna bez drogi do rdzenia (${bezDrogi} z ${komendy.length})`;
}

/** Wpis kanału wraz z nasłuchem powitań — zakładany raz na kanał. */
function wpisPamieci(kanal: Kanal): WpisPamieci {
  const znany = pamiec.get(kanal);
  if (znany !== undefined) return znany;
  const swiezy: WpisPamieci = { odczyt: { uchwyty: null, odmowa: '' }, wToku: null, zalezni: new Set() };
  pamiec.set(kanal, swiezy);
  // Powitanie cudze jest tą samą odpowiedzią co własne: uzgodnienie powłoki wita
  // rdzeń przy każdym nawiązaniu połączenia, więc wykaz odświeża się sam.
  kanal.naDowolny((koperta) => {
    if (koperta.type !== Command.ConnectionHello || koperta.status !== EnvelopeStatus.Ok) return;
    const komendy = (koperta.payload as ConnectionHelloResponse | undefined)?.commands;
    if (!czyTablica(komendy)) return;
    zapisz(swiezy, { uchwyty: new Set(komendy), odmowa: '' });
  });
  return swiezy;
}

/** Jedno powitanie i zapis jego wyniku — udanego albo odmownego. */
async function powitaj(kanal: Kanal, wspolny: WpisPamieci): Promise<void> {
  const wynik = sprawdzKsztalt(
    await wywolaj(kanal, Command.ConnectionHello, {
      // Tożsamość stała, nie świeża: obsługiwacz powitania nie czyta żądania,
      // a `tozsamoscKlienta()` nadaje przy każdym wywołaniu identyfikator nowy —
      // drugi identyfikator rozdzieliłby ognisko od połączenia, które je zgłosiło
      // (`protokol/uzgodnienie.ts`).
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

/** Zapis odczytu i przerysowanie kontrolek wszystkich modułów tego kanału. */
function zapisz(wspolny: WpisPamieci, odczyt: Odczyt): void {
  wspolny.odczyt = odczyt;
  for (const przerysuj of wspolny.zalezni) przerysuj();
}

/** Powód komendy zdaniem podrzędnym — małą literą, do złożenia z czynnością. */
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

/** Zdanie zaczynające się z dużej litery — powody układa się małą. */
function zDuzej(zdanie: string): string {
  return zdanie.charAt(0).toUpperCase() + zdanie.slice(1);
}
