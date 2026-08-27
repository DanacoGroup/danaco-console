/**
 * DROGA WEJŚCIA — przebieg trzech etapów.
 *
 * Sprawdzian prowadzi przebieg, nie okno: przebieg nie dotyka dokumentu, więc
 * każda odsłona — także odsłona błędu połączenia i odsłona zwłoki nałożonej
 * przez rdzeń — jest tu osiągalna naprawdę, a nie tylko opisana. Rozmowę z rdzeniem
 * uruchomionym mierzy osobny sprawdzian, który rdzenia wymaga.
 *
 * Rdzeń jest tu zastąpiony, a nie udawany: sprawdzane jest zachowanie okna
 * wobec odpowiedzi kontraktu, nie sama treść odpowiedzi.
 */

import {
  Command,
  EnvelopeStatus,
  ErrorCode,
  PROTOCOL_VERSION,
  type Envelope,
  type ErrorInfo,
} from '../../../shared/contract.ts';
import type { Transport } from '../polaczenie/gniazdo.ts';
import type { Odsubskrybuj } from '../polaczenie/magistrala-zdarzen.ts';
import type { StanPolaczenia } from '../polaczenie/stan-polaczenia.ts';
import { utworzKanal } from '../protokol/kanal.ts';
import { utworzSesje } from '../protokol/sesja.ts';
import type { TozsamoscKlienta } from '../protokol/tozsamosc-klienta.ts';
import { bieg, poOdstepie, rowne, sprawdz } from '../sprawdzian.ts';
import {
  adresPoprawny,
  hasloSpelnia,
  magazynWPamieci,
  ocenHaslo,
  PROG_ZWLOKI_MS,
  utworzPrzebieg,
  type OdslonaDostepu,
  type Przebieg,
} from './przebieg.ts';

/* ── Rdzeń zastępczy ─────────────────────────────────────────────────────── */

/** Odpowiedź rdzenia na jedną komendę: treść albo odmowa. */
type Odpowiedz = { tresc: unknown } | { blad: ErrorInfo };

/** Transport, którego stanem i odpowiedziami rządzi sprawdzian. */
interface RdzenZastepczy extends Transport {
  /** Ustawia odpowiedź na kolejne wywołania danej komendy. */
  odpowiadaj(komenda: string, odpowiedz: Odpowiedz): void;
  /** Każe rdzeniowi zwlekać przed odpowiedzią, tak jak zwleka rdzeń prawdziwy. */
  zwlekaj(ms: number): void;
  /** Ogłasza stan gniazda tak, jak zrobiłoby to środowisko. */
  ogloszStan(stan: StanPolaczenia): void;
  /** Komendy, które przebieg wysłał, w kolejności wysłania. */
  wyslane: string[];
}

function odmowa(kod: ErrorCode, opis: string): Odpowiedz {
  return { blad: { code: kod, message: opis, retryable: false } };
}

function utworzRdzenZastepczy(): RdzenZastepczy {
  const sluchaczeRamek: ((ramka: string) => void)[] = [];
  const sluchaczeStanu: ((stan: StanPolaczenia) => void)[] = [];
  const odpowiedzi = new Map<string, Odpowiedz>();
  const wyslane: string[] = [];
  let biezacy: StanPolaczenia = 'rozlaczony';
  let zwlokaMs = 0;

  function oddaj(koperta: Envelope): void {
    const odpowiedz = odpowiedzi.get(koperta.type);
    if (odpowiedz === undefined) return;
    const zwrot: Envelope =
      'blad' in odpowiedz
        ? {
            type: koperta.type,
            id: koperta.id,
            timestamp: Date.now(),
            status: EnvelopeStatus.Error,
            error: odpowiedz.blad,
          }
        : {
            type: koperta.type,
            id: koperta.id,
            timestamp: Date.now(),
            status: EnvelopeStatus.Ok,
            payload: odpowiedz.tresc,
          };
    const oddajTeraz = (): void => {
      for (const sluchacz of [...sluchaczeRamek]) sluchacz(JSON.stringify(zwrot));
    };
    // Rdzeń prawdziwy nakłada zwłokę PRZED odpowiedzią; rdzeń zastępczy robi
    // to samo, bo inaczej nie dałoby się zmierzyć tego, co okno ma zmierzyć.
    if (zwlokaMs > 0) setTimeout(oddajTeraz, zwlokaMs);
    else oddajTeraz();
  }

  return {
    wyslane,
    odpowiadaj(komenda, odpowiedz) {
      odpowiedzi.set(komenda, odpowiedz);
    },
    zwlekaj(ms) {
      zwlokaMs = ms;
    },
    ogloszStan(stan) {
      biezacy = stan;
      for (const sluchacz of [...sluchaczeStanu]) sluchacz(stan);
    },
    polacz() {},
    rozlacz() {},
    wyslij(ramka) {
      const koperta = JSON.parse(ramka) as Envelope;
      wyslane.push(koperta.type);
      oddaj(koperta);
    },
    naRamke(sluchacz): Odsubskrybuj {
      sluchaczeRamek.push(sluchacz);
      return () => {
        const miejsce = sluchaczeRamek.indexOf(sluchacz);
        if (miejsce >= 0) sluchaczeRamek.splice(miejsce, 1);
      };
    },
    naStan(sluchacz): Odsubskrybuj {
      sluchacz(biezacy);
      sluchaczeStanu.push(sluchacz);
      return () => {
        const miejsce = sluchaczeStanu.indexOf(sluchacz);
        if (miejsce >= 0) sluchaczeStanu.splice(miejsce, 1);
      };
    },
    stan: () => biezacy,
    oczekujace: () => 0,
  };
}

/** Powitanie odpowiadające wersją zgodną; `zwiazane` daje token urządzenia. */
function powitanie(zwiazane: boolean): Odpowiedz {
  return {
    tresc: {
      serverVersion: PROTOCOL_VERSION,
      protocolVersion: PROTOCOL_VERSION,
      authenticated: zwiazane,
      gatewayConfigured: true,
    },
  };
}

/** Sesja bramki w kształcie kontraktu. */
const SESJA = { tresc: { session: { token: 'token-sprawdzianu', expiresAt: 0 } } };

/** Środowisko i jego moduły w kształcie kontraktu. */
const SRODOWISKA: Odpowiedz = {
  tresc: {
    environments: [
      { id: '2', code: 'workspace', name: 'WorkSpace', order: 2, navigationKind: 'modules' },
      { id: '1', code: 'talkin', name: 'TalkIn', order: 1, navigationKind: 'modules' },
    ],
  },
};

const WEJSCIE: Odpowiedz = {
  tresc: {
    environment: { id: '1', code: 'talkin', name: 'TalkIn', order: 1, navigationKind: 'modules' },
    modules: [{ id: '1', code: 'studio', name: 'Studio', order: 1, operationalWindowCodes: [], kind: 'srodowisko_robocze', configuredOnHome: false }],
    sessions: [],
  },
};

const KLIENT: TozsamoscKlienta = { id: 'klient-sprawdzianu', wersja: '1.0.0' };

/** Przebieg z rdzeniem zastępczym doprowadzony do wskazanego etapu. */
function zaloz(): { przebieg: Przebieg; rdzen: RdzenZastepczy } {
  const rdzen = utworzRdzenZastepczy();
  const przebieg = utworzPrzebieg({
    kanal: utworzKanal(rdzen, utworzSesje()),
    transport: rdzen,
    klient: KLIENT,
    magazyn: magazynWPamieci(),
  });
  rdzen.odpowiadaj(Command.EnvironmentList, SRODOWISKA);
  rdzen.odpowiadaj(Command.EnvironmentEnter, WEJSCIE);
  return { przebieg, rdzen };
}

/** Doprowadza przebieg do etapu dostępu, czyli za powitanie bez tokenu. */
async function doDostepu(): Promise<{ przebieg: Przebieg; rdzen: RdzenZastepczy }> {
  const { przebieg, rdzen } = zaloz();
  rdzen.odpowiadaj(Command.ConnectionHello, powitanie(false));
  przebieg.polacz();
  rdzen.ogloszStan('polaczony');
  await poOdstepie();
  return { przebieg, rdzen };
}

const HASLO_MOCNE = 'Sprawdzian-2026!';

await bieg('droga wejścia — przebieg', {
  /* ── Etap 1 ───────────────────────────────────────────────────────────── */

  async 'powitanie bez tokenu prowadzi z uruchomienia do logowania'() {
    const { przebieg, rdzen } = await doDostepu();
    rowne(przebieg.stan().etap, 'dostep', 'etap po powitaniu');
    rowne(przebieg.stan().odslona, 'logowanie', 'odsłona po powitaniu');
    sprawdz(rdzen.wyslane.includes(Command.ConnectionHello), 'powitanie poszło do rdzenia');
  },

  async 'powitanie z ważnym tokenem rozpoznaje urządzenie i idzie do przygotowania'() {
    const { przebieg, rdzen } = zaloz();
    rdzen.odpowiadaj(Command.ConnectionHello, powitanie(true));
    const odsloniete: string[] = [];
    przebieg.naZmiane((stan) => odsloniete.push(stan.odslona));
    przebieg.polacz();
    rdzen.ogloszStan('polaczony');
    await poOdstepie();
    sprawdz(odsloniete.includes('w-token'), 'odsłona rozpoznanego urządzenia wystąpiła');
    rowne(przebieg.stan().etap, 'przygotowanie', 'etap po rozpoznaniu urządzenia');
  },

  async 'zerwane połączenie odsłania błąd połączenia'() {
    const { przebieg, rdzen } = zaloz();
    rdzen.odpowiadaj(Command.ConnectionHello, powitanie(false));
    przebieg.polacz();
    rdzen.ogloszStan('ponawianie');
    await poOdstepie();
    rowne(przebieg.stan().odslona, 'w-blad', 'odsłona po zerwaniu połączenia');
    rowne(przebieg.stan().etapyLaczenia[0], 'blad', 'pierwszy etap łączenia nieudany');
  },

  async 'rozjazd wersji protokołu zatrzymuje wejście i nazywa obie wersje'() {
    const { przebieg, rdzen } = zaloz();
    rdzen.odpowiadaj(Command.ConnectionHello, {
      tresc: { serverVersion: '9.9', protocolVersion: '9.9', authenticated: false },
    });
    przebieg.polacz();
    rdzen.ogloszStan('polaczony');
    await poOdstepie();
    rowne(przebieg.stan().odslona, 'w-blad', 'odsłona przy rozjeździe wersji');
    rowne(przebieg.stan().usterki[0]?.klucz, 'uruchomienie.wersja', 'rozpoznanie rozjazdu');
    rowne(przebieg.stan().powitanie?.protocolVersion, '9.9', 'wersja rdzenia zapamiętana');
  },

  async 'odmowa powitania odsłania błąd połączenia wraz z opisem rdzenia'() {
    const { przebieg, rdzen } = zaloz();
    rdzen.odpowiadaj(
      Command.ConnectionHello,
      odmowa(ErrorCode.InternalError, 'rdzeń zastępczy odmawia powitania'),
    );
    przebieg.polacz();
    rdzen.ogloszStan('polaczony');
    await poOdstepie();
    rowne(przebieg.stan().odslona, 'w-blad', 'odsłona po odmowie powitania');
    rowne(
      przebieg.stan().usterki[0]?.odRdzenia?.message,
      'rdzeń zastępczy odmawia powitania',
      'opis odmowy pochodzi z rdzenia',
    );
  },

  /* ── Etap 2: rejestracja, obie gałęzie pozycji 11 ─────────────────────── */

  async 'rejestracja z kontem nadawczym prowadzi do potwierdzenia drogą z listu'() {
    const { przebieg, rdzen } = await doDostepu();
    rdzen.odpowiadaj(Command.AuthRegister, {
      tresc: { registered: true, pendingVerification: true },
    });
    await przebieg.zarejestruj({
      login: 'operator',
      email: 'operator@danaco-group.pl',
      haslo: HASLO_MOCNE,
      hasloPowtorzone: HASLO_MOCNE,
    });
    rowne(przebieg.stan().odslona, 'kod', 'odsłona przy nadajniku ustawionym');
    rowne(przebieg.stan().adres, 'operator@danaco-group.pl', 'adres zapamiętany dla zdania odsłony');
  },

  async 'rejestracja bez konta nadawczego nazywa niepotwierdzony adres'() {
    const { przebieg, rdzen } = await doDostepu();
    rdzen.odpowiadaj(Command.AuthRegister, {
      tresc: { registered: true, pendingVerification: false },
    });
    await przebieg.zarejestruj({
      login: 'operator',
      email: 'operator@danaco-group.pl',
      haslo: HASLO_MOCNE,
      hasloPowtorzone: HASLO_MOCNE,
    });
    rowne(
      przebieg.stan().odslona,
      'konto-bez-potwierdzenia',
      'odsłona przy braku konta nadawczego',
    );
    rowne(przebieg.stan().adres, 'operator@danaco-group.pl', 'adres, który zostaje niepotwierdzony');
  },

  async 'powtórzona rejestracja oddaje odmowę rdzenia bez własnej interpretacji'() {
    const { przebieg, rdzen } = await doDostepu();
    rdzen.odpowiadaj(
      Command.AuthRegister,
      odmowa(ErrorCode.Conflict, 'bramka: konto właściciela jest już założone'),
    );
    przebieg.przejdzDo('rejestracja');
    await przebieg.zarejestruj({
      login: 'operator',
      email: 'operator@danaco-group.pl',
      haslo: HASLO_MOCNE,
      hasloPowtorzone: HASLO_MOCNE,
    });
    rowne(przebieg.stan().usterki[0]?.odRdzenia?.code, ErrorCode.Conflict, 'kod odmowy');
    rowne(przebieg.stan().odslona, 'rejestracja', 'odsłona po odmowie nie zmienia się');
  },

  'rejestracja sprawdza dane przed wysłaniem i zbiera rozpoznania w jeden wykaz'() {
    const rozpoznania = [
      { dane: { login: '', email: '', haslo: '', hasloPowtorzone: '' }, ile: 3 },
      { dane: { login: 'o', email: 'bez-malpy', haslo: HASLO_MOCNE, hasloPowtorzone: HASLO_MOCNE }, ile: 1 },
      { dane: { login: 'o', email: 'o@d.pl', haslo: 'krotkie', hasloPowtorzone: 'krotkie' }, ile: 1 },
      { dane: { login: 'o', email: 'o@d.pl', haslo: HASLO_MOCNE, hasloPowtorzone: 'inne' }, ile: 1 },
    ];
    sprawdz(rozpoznania.length > 0, 'zbiór przypadków niepusty');
    sprawdz(!adresPoprawny('bez-malpy'), 'adres bez małpy odrzucony');
    sprawdz(!adresPoprawny('a@b'), 'adres bez domeny odrzucony');
    sprawdz(adresPoprawny('operator@danaco-group.pl'), 'adres poprawny przyjęty');
    sprawdz(!hasloSpelnia('krotkie'), 'hasło poniżej wymagań odrzucone');
    sprawdz(hasloSpelnia(HASLO_MOCNE), 'hasło spełniające wymagania przyjęte');
    rowne(Object.values(ocenHaslo('')).filter(Boolean).length, 0, 'puste hasło spełnia zero warunków');
    rowne(Object.values(ocenHaslo(HASLO_MOCNE)).filter(Boolean).length, 4, 'mocne hasło spełnia cztery');
  },

  async 'puste pola rejestracji dają trzy rozpoznania naraz'() {
    const { przebieg } = await doDostepu();
    await przebieg.zarejestruj({ login: '', email: '', haslo: '', hasloPowtorzone: '' });
    rowne(przebieg.stan().usterki.length, 3, 'liczba rozpoznań przy pustym formularzu');
  },

  /* ── Etap 2: logowanie, zwłoka zamiast zapory ────────────────────────── */

  async 'oba pola logowania puste to jedna sprawa, nie dwie'() {
    const { przebieg } = await doDostepu();
    await przebieg.zaloguj({ login: '', haslo: '' });
    rowne(przebieg.stan().usterki.length, 1, 'liczba rozpoznań');
    rowne(przebieg.stan().usterki[0]?.klucz, 'brakDanych', 'rozpoznanie zbiorcze');
  },

  async 'logowanie udane prowadzi do przygotowania środowiska'() {
    const { przebieg, rdzen } = await doDostepu();
    rdzen.odpowiadaj(Command.AuthLogin, SESJA);
    await przebieg.zaloguj({ login: 'operator', haslo: HASLO_MOCNE });
    rowne(przebieg.stan().etap, 'przygotowanie', 'etap po udanym logowaniu');
    rowne(przebieg.stan().sesjaBramki?.token, 'token-sprawdzianu', 'token sesji bramki przyjęty');
    sprawdz(rdzen.wyslane.includes(Command.EnvironmentEnter), 'wejście do środowiska poszło');
  },

  async 'nieudane próby nie zamykają bramki — rdzeń progu prób nie ma'() {
    const { przebieg, rdzen } = await doDostepu();
    rdzen.odpowiadaj(
      Command.AuthLogin,
      odmowa(ErrorCode.NotAuthenticated, 'bramka: sekret metody password nie zgadza się'),
    );
    const odsloniete: string[] = [];
    for (let proba = 1; proba <= 8; proba += 1) {
      await przebieg.zaloguj({ login: 'operator', haslo: 'zle' });
      odsloniete.push(przebieg.stan().odslona);
    }
    sprawdz(
      odsloniete.every((odslona) => odslona === 'logowanie-blad'),
      `ósma próba została odmówiona progiem: ${odsloniete.join(', ')}`,
    );
    // Po ośmiu nieudanych próbach hasło poprawne wpuszcza — tak stanowi
    // kontrakt i tak zachowuje się rdzeń, zmierzone uruchomieniem.
    rdzen.odpowiadaj(Command.AuthLogin, SESJA);
    await przebieg.zaloguj({ login: 'operator', haslo: HASLO_MOCNE });
    rowne(przebieg.stan().etap, 'przygotowanie', 'etap po haśle poprawnym');
  },

  async 'zwłoka nałożona przez rdzeń jest mierzona i odsłaniana Operatorowi'() {
    const { przebieg, rdzen } = await doDostepu();
    rdzen.odpowiadaj(
      Command.AuthLogin,
      odmowa(ErrorCode.NotAuthenticated, 'bramka: sekret metody password nie zgadza się'),
    );
    rdzen.zwlekaj(PROG_ZWLOKI_MS + 200);
    await przebieg.zaloguj({ login: 'operator', haslo: 'zle' });
    rowne(przebieg.stan().odslona, 'logowanie-wstrzymane', 'odsłona po zwłoce nazwanej');
    sprawdz(
      przebieg.stan().zwlokaS >= 1,
      `zmierzona zwłoka poniżej sekundy: ${przebieg.stan().zwlokaS}`,
    );
    przebieg.zdejmijWstrzymanie();
    rowne(przebieg.stan().odslona, 'logowanie', 'odsłona po upływie zwłoki');
    rowne(przebieg.stan().zwlokaS, 0, 'licznik zwłoki wyzerowany');
  },

  async 'zwłoka krótsza od sekundy nie odsłania licznika, bo nie ma czego liczyć'() {
    const { przebieg, rdzen } = await doDostepu();
    rdzen.odpowiadaj(
      Command.AuthLogin,
      odmowa(ErrorCode.NotAuthenticated, 'bramka: sekret metody password nie zgadza się'),
    );
    await przebieg.zaloguj({ login: 'operator', haslo: 'zle' });
    rowne(przebieg.stan().odslona, 'logowanie-blad', 'odsłona przy zwłoce nieznaczącej');
    rowne(przebieg.stan().zwlokaS, 0, 'zwłoka nienazwana');
  },

  async 'odmowa spoza uwierzytelnienia nie jest zwłoką bramki'() {
    const { przebieg, rdzen } = await doDostepu();
    rdzen.odpowiadaj(Command.AuthLogin, odmowa(ErrorCode.InternalError, 'rdzeń: błąd wewnętrzny'));
    rdzen.zwlekaj(PROG_ZWLOKI_MS + 200);
    await przebieg.zaloguj({ login: 'operator', haslo: HASLO_MOCNE });
    rowne(przebieg.stan().odslona, 'logowanie-blad', 'odsłona przy błędzie rdzenia');
    rowne(przebieg.stan().zwlokaS, 0, 'błąd rdzenia nie jest zwłoką bramki');
  },

  /* ── Etap 2: potwierdzenie adresu i odzyskanie konta ──────────────────── */

  async 'potwierdzenie adresu drogą z listu wydaje sesję i prowadzi dalej'() {
    const { przebieg, rdzen } = await doDostepu();
    rdzen.odpowiadaj(Command.AuthVerify, {
      tresc: { verified: true, session: { token: 'token-po-potwierdzeniu', expiresAt: 0 } },
    });
    przebieg.przejdzDo('kod');
    await przebieg.potwierdzAdres({ droga: 'md8IWBXN5q68eJZx7xFg0uJO8JKbOm6JpOXnq2gLmCY' });
    rowne(przebieg.stan().etap, 'przygotowanie', 'etap po potwierdzeniu adresu');
    rowne(przebieg.stan().sesjaBramki?.token, 'token-po-potwierdzeniu', 'sesja wydana potwierdzeniem');
  },

  async 'pusta droga potwierdzenia nie idzie do rdzenia'() {
    const { przebieg, rdzen } = await doDostepu();
    przebieg.przejdzDo('kod');
    await przebieg.potwierdzAdres({ droga: '   ' });
    rowne(przebieg.stan().usterki[0]?.klucz, 'brakDrogi', 'rozpoznanie pustej drogi');
    sprawdz(!rdzen.wyslane.includes(Command.AuthVerify), 'komenda nie poszła');
  },

  async 'prośba o odzyskanie prowadzi do przepisania drogi z listu'() {
    const { przebieg, rdzen } = await doDostepu();
    rdzen.odpowiadaj(Command.AuthRecover, { tresc: { sent: true } });
    przebieg.przejdzDo('odzyskiwanie-adres');
    await przebieg.poprosOOdzyskanie('operator@danaco-group.pl');
    rowne(przebieg.stan().odslona, 'odzyskiwanie-kod', 'odsłona po wysłaniu');
    rowne(przebieg.stan().adres, 'operator@danaco-group.pl', 'adres zapamiętany');
  },

  async 'odzyskiwanie bez konta nadawczego oddaje odmowę rdzenia'() {
    const { przebieg, rdzen } = await doDostepu();
    rdzen.odpowiadaj(
      Command.AuthRecover,
      odmowa(ErrorCode.InternalError, 'bramka: konto nadawcze platformy nie ma serwera poczty'),
    );
    przebieg.przejdzDo('odzyskiwanie-adres');
    await przebieg.poprosOOdzyskanie('operator@danaco-group.pl');
    rowne(przebieg.stan().odslona, 'odzyskiwanie-adres', 'odsłona po odmowie');
    sprawdz(
      (przebieg.stan().usterki[0]?.odRdzenia?.message ?? '').includes('konto nadawcze'),
      'odmowa nazywa brak konta nadawczego',
    );
  },

  async 'powtarzane wysyłki drogi odzyskania nie są wstrzymywane'() {
    const { przebieg, rdzen } = await doDostepu();
    rdzen.odpowiadaj(Command.AuthRecover, { tresc: { sent: true } });
    przebieg.przejdzDo('odzyskiwanie-adres');
    for (let wyslanie = 1; wyslanie <= 7; wyslanie += 1) {
      await przebieg.poprosOOdzyskanie('operator@danaco-group.pl');
      rowne(przebieg.stan().odslona, 'odzyskiwanie-kod', `wysłanie ${wyslanie} przyjęte`);
      przebieg.przejdzDo('odzyskiwanie-adres');
    }
  },

  async 'zwłoka przy wysyłaniu drogi odzyskania jest nazywana tym samym pomiarem'() {
    const { przebieg, rdzen } = await doDostepu();
    rdzen.odpowiadaj(Command.AuthRecover, { tresc: { sent: true } });
    rdzen.zwlekaj(PROG_ZWLOKI_MS + 200);
    przebieg.przejdzDo('odzyskiwanie-adres');
    await przebieg.poprosOOdzyskanie('operator@danaco-group.pl');
    rowne(przebieg.stan().odslona, 'odzyskiwanie-wstrzymane', 'odsłona zwłoki wysyłania');
    przebieg.zdejmijWstrzymanie();
    rowne(przebieg.stan().odslona, 'odzyskiwanie-adres', 'odsłona po upływie zwłoki');
  },

  async 'nowe hasło unieważnia sesje i zawraca do logowania'() {
    const { przebieg, rdzen } = await doDostepu();
    rdzen.odpowiadaj(Command.AuthReset, { tresc: { changed: true, revokedDevices: 1 } });
    przebieg.przejdzDo('odzyskiwanie-haslo');
    await przebieg.ustawNoweHaslo({
      droga: 'GeKYOuKsWC1TM3NsLhODRudWop4c3UW0LfO0pP55fkk',
      haslo: HASLO_MOCNE,
      hasloPowtorzone: HASLO_MOCNE,
    });
    rowne(przebieg.stan().odslona, 'logowanie', 'odsłona po ustawieniu hasła');
    sprawdz(rdzen.wyslane.includes(Command.AuthReset), 'komenda poszła do rdzenia');
  },

  /* ── Etap 3 ───────────────────────────────────────────────────────────── */

  async 'przygotowanie wchodzi do środowiska pierwszego w kolejności'() {
    const { przebieg, rdzen } = await doDostepu();
    rdzen.odpowiadaj(Command.AuthLogin, SESJA);
    await przebieg.zaloguj({ login: 'operator', haslo: HASLO_MOCNE });
    rowne(przebieg.stan().srodowisko?.code, 'talkin', 'środowisko o kolejności pierwszej');
    rowne(przebieg.stan().moduly.length, 1, 'moduły środowiska odczytane');
    rowne(przebieg.stan().przygotowanie.stany[0], 'gotowy', 'uwierzytelnienie zakończone');
    rowne(przebieg.stan().przygotowanie.stany[2], 'gotowy', 'przywracanie kart zakończone');
    rowne(
      przebieg.stan().przygotowanie.stany[1],
      'oczekuje',
      'etap bez komendy drogi wejścia zostaje w oczekiwaniu',
    );
    rowne(przebieg.stan().przygotowanie.wartosc, 40, 'postęp liczony z etapów zakończonych');
  },

  async 'odmowa wejścia do środowiska zaznacza etap jako nieudany'() {
    const { przebieg, rdzen } = await doDostepu();
    rdzen.odpowiadaj(Command.AuthLogin, SESJA);
    rdzen.odpowiadaj(
      Command.EnvironmentEnter,
      odmowa(ErrorCode.NotAuthenticated, 'bramka: połączenie niezwiązane z sesją'),
    );
    await przebieg.zaloguj({ login: 'operator', haslo: HASLO_MOCNE });
    rowne(przebieg.stan().przygotowanie.stany[2], 'blad', 'etap przywracania nieudany');
    sprawdz(przebieg.stan().usterki[0]?.odRdzenia !== undefined, 'odmowa rdzenia zachowana');
  },

  /* ── Osiągalność wszystkich odsłon ────────────────────────────────────── */

  async 'każda odsłona etapu dostępu jest osiągalna'() {
    const wszystkie: OdslonaDostepu[] = [
      'logowanie',
      'logowanie-blad',
      'logowanie-wstrzymane',
      'rejestracja',
      'kod',
      'konto-bez-potwierdzenia',
      'odzyskiwanie-adres',
      'odzyskiwanie-kod',
      'odzyskiwanie-haslo',
      'odzyskiwanie-wstrzymane',
    ];
    const { przebieg } = await doDostepu();
    for (const odslona of wszystkie) {
      przebieg.przejdzDo(odslona);
      rowne(przebieg.stan().odslona, odslona, `odsłona ${odslona} osiągnięta`);
    }
    rowne(wszystkie.length, 10, 'liczba odsłon etapu dostępu');
  },
});
