/**
 * Droga wejścia — rozmowa z rdzeniem uruchomionym naprawdę. Sprawdzian
 * prowadzi trzy etapy wejścia przez ten sam przebieg, którym idzie okno,
 * wymagając rdzenia nasłuchującego na świeżej bazie.
 */

import { PROTOCOL_VERSION } from '../../../shared/contract.ts';
import { adresGniazdaRdzenia, adresRdzeniaLokalnego } from '../polaczenie/adres-rdzenia.ts';
import { utworzTransport } from '../polaczenie/gniazdo.ts';
import { utworzKanal } from '../protokol/kanal.ts';
import { utworzSesje } from '../protokol/sesja.ts';
import { tozsamoscKlienta } from '../protokol/tozsamosc-klienta.ts';
import { bieg, sprawdz } from '../sprawdzian.ts';
import { magazynWPamieci, utworzPrzebieg, type StanPrzebiegu } from './przebieg.ts';

/** Górna granica oczekiwania na przejście przebiegu do spodziewanego stanu, wyrażona w milisekundach czasu. */
const GRANICA_MS = 10_000;

/** Konto zakładane przez ten sprawdzian na świeżej bazie rdzenia, jednorazowo, przy pierwszym uruchomieniu. */
const KONTO = {
  login: 'operator',
  email: 'operator@danaco-group.pl',
  haslo: 'Sprawdzian-Wejscia-2026!',
};

/**
 * Adres rdzenia. Bez wskazania idzie pętla zwrotna z portem domyślnym; ze
 * wskazaniem — adres HTTP przełożony przez warstwę połączenia, żeby dwa
 * rdzenie dały się zmierzyć jednym sprawdzianem.
 */
function adresRdzenia(): string {
  const srodowisko = (globalThis as { process?: { env?: Record<string, string | undefined> } })
    .process?.env;
  const wskazany = srodowisko?.['DANACO_RDZEN'];
  if (wskazany === undefined || wskazany.length === 0) return adresRdzeniaLokalnego();
  return adresGniazdaRdzenia(wskazany) ?? adresRdzeniaLokalnego();
}

const ADRES = adresRdzenia();

/**
 * Droga potwierdzenia przychodzi listem, więc sprawdzian nie ma jej skąd wziąć
 * sam. Wskazana z zewnątrz domyka gałąź z kontem nadawczym w całości; bez niej
 * sprawdzian mierzy samą odmowę bramki, która tej drogi żąda.
 */
function drogaZListu(): string | undefined {
  const srodowisko = (globalThis as { process?: { env?: Record<string, string | undefined> } })
    .process?.env;
  const wskazana = srodowisko?.['DANACO_DROGA'];
  return wskazana === undefined || wskazana.length === 0 ? undefined : wskazana;
}

/** Czeka, aż przebieg dojdzie do stanu spełniającego podany warunek, sprawdzając go po każdej jego zmianie. */
function azDo(
  przebieg: { naZmiane(s: (stan: StanPrzebiegu) => void): () => void; stan(): StanPrzebiegu },
  warunek: (stan: StanPrzebiegu) => boolean,
  czynnosc: string,
): Promise<StanPrzebiegu> {
  if (warunek(przebieg.stan())) return Promise.resolve(przebieg.stan());
  return new Promise((rozstrzygnij, odrzuc) => {
    const zegar = setTimeout(() => {
      odsubskrybuj();
      odrzuc(
        new Error(
          `${czynnosc}: przebieg nie doszedł do stanu w ciągu ${GRANICA_MS} ms pod adresem ` +
            `${ADRES} — uruchom rdzeń na świeżej bazie przed tym sprawdzianem`,
        ),
      );
    }, GRANICA_MS);
    const odsubskrybuj = przebieg.naZmiane((stan) => {
      if (!warunek(stan)) return;
      clearTimeout(zegar);
      odsubskrybuj();
      rozstrzygnij(stan);
    });
  });
}

function wypisz(nazwa: string, tresc: unknown): void {
  console.log(`           ${nazwa}: ${JSON.stringify(tresc)}`);
}

await bieg('droga wejścia — rozmowa z rdzeniem', {
  async 'trzy etapy wejścia przechodzą przez żywy rdzeń'() {
    const transport = utworzTransport(ADRES);
    const przebieg = utworzPrzebieg({
      kanal: utworzKanal(transport, utworzSesje()),
      transport,
      klient: tozsamoscKlienta(),
      magazyn: magazynWPamieci(),
    });
    try {
      /* ── Etap 1: łączenie i powitanie ───────────────────────────────── */
      przebieg.polacz();
      const poPowitaniu = await azDo(
        przebieg,
        (stan) => stan.powitanie !== undefined || stan.odslona === 'w-blad',
        'powitanie połączenia',
      );
      sprawdz(
        poPowitaniu.powitanie !== undefined,
        `rdzeń odmówił powitania: ${JSON.stringify(poPowitaniu.usterki)}`,
      );
      const powitanie = poPowitaniu.powitanie!;
      wypisz('connection.hello', {
        serverVersion: powitanie.serverVersion,
        protocolVersion: powitanie.protocolVersion,
        authenticated: powitanie.authenticated,
        gatewayConfigured: powitanie.gatewayConfigured,
        commands: powitanie.commands === undefined ? undefined : `${powitanie.commands.length} nazw`,
      });
      sprawdz(
        powitanie.protocolVersion === PROTOCOL_VERSION,
        `wersja protokołu rdzenia ${powitanie.protocolVersion} wobec klienta ${PROTOCOL_VERSION}`,
      );
      sprawdz(
        typeof powitanie.serverVersion === 'string' && powitanie.serverVersion.length > 0,
        'rdzeń nie podał własnej wersji',
      );
      sprawdz(poPowitaniu.etap === 'dostep', `etap po powitaniu: ${poPowitaniu.etap}`);

      /* ── Etap 2a: rejestracja, obie gałęzie ──────────────────────────── */
      const droga = drogaZListu();
      let zNadajnikiem: boolean;
      if (droga === undefined) {
        sprawdz(
          powitanie.gatewayConfigured === false,
          'rdzeń ma już założone konto właściciela — sprawdzian wymaga bazy świeżej',
        );
        przebieg.przejdzDo('rejestracja');
        await przebieg.zarejestruj({
          login: KONTO.login,
          email: KONTO.email,
          haslo: KONTO.haslo,
          hasloPowtorzone: KONTO.haslo,
        });
        const poRejestracji = przebieg.stan();
        sprawdz(
          poRejestracji.usterki.length === 0,
          `rdzeń odmówił rejestracji: ${JSON.stringify(poRejestracji.usterki)}`,
        );
        // Gałąź rozstrzyga rdzeń, nie sprawdzian: z kontem nadawczym prowadzi do przepisania drogi z listu.
        zNadajnikiem = poRejestracji.odslona === 'kod';
        wypisz('auth.register', {
          odslona: poRejestracji.odslona,
          galaz: zNadajnikiem ? 'pendingVerification=true' : 'pendingVerification=false',
        });
        sprawdz(
          poRejestracji.odslona === 'kod' || poRejestracji.odslona === 'konto-bez-potwierdzenia',
          `odsłona po rejestracji poza obiema gałęziami: ${poRejestracji.odslona}`,
        );
        sprawdz(
          poRejestracji.adres === KONTO.email,
          'okno nie zapamiętało adresu, który zostaje niepotwierdzony',
        );
      } else {
        // Bieg kontynuujący: konto stoi z biegu poprzedniego, a droga z listu
        // domyka gałąź z kontem nadawczym.
        sprawdz(
          powitanie.gatewayConfigured === true,
          'wskazano drogę potwierdzenia, a konta właściciela nie ma — bieg nie ma czego kontynuować',
        );
        zNadajnikiem = true;
        wypisz('auth.register', { galaz: 'pendingVerification=true (bieg poprzedni)' });
      }

      // Etap 2b: wejście gałęzią wskazaną przez rdzeń; gałęzie różnią się bramką, nie tylko komunikatem.
      if (zNadajnikiem) {
        if (droga === undefined) {
          przebieg.przejdzDo('logowanie');
          await przebieg.zaloguj({ login: KONTO.login, haslo: KONTO.haslo });
          const odmowa = przebieg.stan().usterki[0]?.odRdzenia;
          wypisz('auth.login', { code: odmowa?.code, message: odmowa?.message });
          sprawdz(
            (odmowa?.message ?? '').includes(KONTO.email),
            'odmowa logowania nie nazywa adresu czekającego na potwierdzenie',
          );
          return;
        }
        przebieg.przejdzDo('kod');
        await przebieg.potwierdzAdres({ droga });
        wypisz('auth.verify', { znakowDrogi: droga.length });
      } else {
        przebieg.przejdzDo('logowanie');
        await przebieg.zaloguj({ login: KONTO.login, haslo: KONTO.haslo });
      }
      const poWejsciuDoKonta = przebieg.stan();
      sprawdz(
        poWejsciuDoKonta.sesjaBramki !== undefined,
        `rdzeń nie wydał sesji: ${JSON.stringify(poWejsciuDoKonta.usterki)}`,
      );
      wypisz(zNadajnikiem ? 'auth.verify — sesja' : 'auth.login', {
        tokenDlugosc: poWejsciuDoKonta.sesjaBramki!.token.length,
        method: poWejsciuDoKonta.sesjaBramki!.method,
        wygasa: poWejsciuDoKonta.sesjaBramki!.expiresAt > Date.now(),
      });

      /* ── Etap 3: przygotowanie środowiska ───────────────────────────── */
      const poWejsciu = await azDo(
        przebieg,
        (stan) => stan.srodowisko !== undefined || stan.przygotowanie.stany[1] === 'blad',
        'wejście do środowiska',
      );
      sprawdz(
        poWejsciu.srodowisko !== undefined,
        `rdzeń odmówił wejścia do środowiska: ${JSON.stringify(poWejsciu.usterki)}`,
      );
      wypisz('environment.enter', {
        environment: poWejsciu.srodowisko!.code,
        order: poWejsciu.srodowisko!.order,
        modulow: poWejsciu.moduly.length,
        sesji: poWejsciu.sesje.length,
      });
      sprawdz(poWejsciu.etap === 'przygotowanie', `etap końcowy: ${poWejsciu.etap}`);
      sprawdz(
        poWejsciu.przygotowanie.stany[0] === 'gotowy',
        'etap uwierzytelnienia nie został domknięty',
      );
    } finally {
      transport.rozlacz();
    }
  },
});
