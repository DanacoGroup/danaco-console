// Czynności obserwacji i odkryć: nastawienie monitora, jego odświeżenie oraz
// szukanie źródeł i odrzucanie pozycji wyniku wraz z uzasadnieniem.
import {
  Command,
  ResearchDiscoveryMode,
  ResearchMonitorKind,
} from '../../../shared/contract.ts';
import type { Odsubskrybuj } from '../polaczenie/magistrala-zdarzen.ts';
import { wywolaj } from '../protokol/wywolanie.ts';
import { niegotowyPanel } from './okno-modulu.ts';
import {
  brakOkna,
  nasluchCzynnosci,
  nieznanaCzynnosc,
  odmowa,
  powiedz,
  wierszCzynnosci,
  wskazanie,
  wybierz,
  zapytaj,
  type CzynnoscWiersza,
  type KontekstBadania,
} from './research-czynnosci.ts';

const ATRYBUT_ODKRYCIA = 'data-odkrycie-badania';
const ATRYBUT_CZYNNOSCI = 'data-czynnosc-obserwacji';
const GNIAZDO_ODKRYC = '[data-odkrycia-wynik]';

const CZYNNOSCI_ODKRYCIA: readonly CzynnoscWiersza[] = [{ kod: 'odrzuc', etykieta: 'Odrzuć' }];

const RODZAJE_MONITORA: readonly (readonly [string, ResearchMonitorKind])[] = [
  ['Monitor tematu', ResearchMonitorKind.Topic],
  ['Kanał RSS albo Atom', ResearchMonitorKind.Feed],
];

const TRYBY_SZUKANIA: readonly (readonly [string, ResearchDiscoveryMode])[] = [
  ['Sieć', ResearchDiscoveryMode.Web],
  ['Bazy publikacji naukowych', ResearchDiscoveryMode.Scholarly],
];

export function zwiazCzynnosciObserwacji(kontekst: KontekstBadania): Odsubskrybuj {
  return nasluchCzynnosci(kontekst.korzen, ATRYBUT_CZYNNOSCI, async (kod, przycisk) => {
    if (brakOkna(kontekst.idOkna())) return;
    await wykonaj(kontekst, kod, wskazanie(przycisk, ATRYBUT_ODKRYCIA));
  });
}

async function wykonaj(kontekst: KontekstBadania, kod: string, klucz: string): Promise<void> {
  if (kod === 'nastaw') return nastaw(kontekst);
  if (kod === 'odswiez') return odswiezMonitory(kontekst);
  if (kod === 'szukaj') return szukaj(kontekst);
  if (kod === 'odrzuc') return odrzuc(kontekst, klucz);
  nieznanaCzynnosc(kod);
}

async function nastaw(kontekst: KontekstBadania): Promise<void> {
  const numer = wybierz('Rodzaj monitora', RODZAJE_MONITORA.map((pozycja) => pozycja[0]));
  const rodzaj = numer < 0 ? undefined : RODZAJE_MONITORA[numer];
  if (rodzaj === undefined) {
    odmowa(undefined, 'Rodzaj monitora nie został wskazany, więc nastawienie nie poszło do rdzenia.');
    return;
  }
  const kanalowy = rodzaj[1] === ResearchMonitorKind.Feed;
  const wpis = zapytaj(kanalowy ? 'Adres kanału' : 'Zapytanie monitora tematu');
  if (wpis === '') return;
  const wynik = await wywolaj(kontekst.kanal, Command.ResearchMonitorSet, {
    windowId: kontekst.idOkna(),
    kind: rodzaj[1],
    query: kanalowy ? undefined : wpis,
    url: kanalowy ? wpis : undefined,
    enabled: true,
  });
  if (!wynik.udany) {
    odmowa(wynik.blad, 'Rdzeń odmówił nastawienia monitora.');
    return;
  }
  powiedz('Monitor nastawiony.');
  await kontekst.odswiez();
}

/* Monitor, którego rdzeń nie odświeżył, jest nazwany osobno: milcząca cisza
   dostawcy wyglądałaby jak brak nowych pozycji. */
async function odswiezMonitory(kontekst: KontekstBadania): Promise<void> {
  const wynik = await wywolaj(kontekst.kanal, Command.ResearchMonitorRefresh, {
    windowId: kontekst.idOkna(),
  });
  if (!wynik.udany || wynik.wynik === undefined) {
    odmowa(wynik.blad, 'Rdzeń odmówił odświeżenia monitorów.');
    return;
  }
  const nieudane = wynik.wynik.failedMonitorIds ?? [];
  powiedz(`Nowych pozycji: ${String(wynik.wynik.results.length)}.`);
  if (nieudane.length > 0) {
    odmowa(undefined, `Monitorów nieodświeżonych: ${String(nieudane.length)}.`);
  }
  await kontekst.odswiez();
}

async function szukaj(kontekst: KontekstBadania): Promise<void> {
  const zapytanie = zapytaj('Zapytanie wyszukiwawcze');
  if (zapytanie === '') return;
  const numer = wybierz('Tryb wyszukiwania', TRYBY_SZUKANIA.map((pozycja) => pozycja[0]));
  const tryb = numer < 0 ? undefined : TRYBY_SZUKANIA[numer];
  if (tryb === undefined) {
    odmowa(undefined, 'Tryb wyszukiwania nie został wskazany, więc zapytanie nie poszło do rdzenia.');
    return;
  }
  const cialo = kontekst.korzen.querySelector(GNIAZDO_ODKRYC);
  const wynik = await wywolaj(kontekst.kanal, Command.ResearchDiscoverySearch, {
    windowId: kontekst.idOkna(),
    query: zapytanie,
    mode: tryb[1],
  });
  if (!wynik.udany || wynik.wynik === undefined) {
    niegotowyPanel(cialo, wynik.blad?.message ?? 'Wynik wyszukiwania nie doszedł.');
    odmowa(wynik.blad, 'Rdzeń odmówił wyszukania źródeł.');
    return;
  }
  const pozycje = wynik.wynik.results;
  if (cialo === null) return;
  if (pozycje.length === 0) {
    niegotowyPanel(cialo, 'Żaden dostawca nie oddał pozycji na to zapytanie.');
  } else {
    cialo.replaceChildren(...pozycje.map((pozycja) => wierszCzynnosci(cialo,
      pozycja.title, pozycja.provider, ATRYBUT_ODKRYCIA, pozycja.key,
      ATRYBUT_CZYNNOSCI, CZYNNOSCI_ODKRYCIA)));
  }
  const milczacy = wynik.wynik.providersFailed ?? [];
  if (milczacy.length > 0) {
    odmowa(undefined, `Dostawcy, którzy nie odpowiedzieli: ${milczacy.join(', ')}.`);
  }
}

async function odrzuc(kontekst: KontekstBadania, klucz: string): Promise<void> {
  if (klucz === '') {
    odmowa(undefined, 'Przycisk nie stoi przy żadnej pozycji wyniku, więc nic nie zostało odrzucone.');
    return;
  }
  const powod = zapytaj('Uzasadnienie odrzucenia');
  if (powod === '') {
    odmowa(undefined, 'Odrzucenie bez uzasadnienia nie zasila przesiewu, więc nie poszło do rdzenia.');
    return;
  }
  const wynik = await wywolaj(kontekst.kanal, Command.ResearchDiscoveryReject, {
    windowId: kontekst.idOkna(),
    resultKeys: [klucz],
    reason: powod,
  });
  if (!wynik.udany || wynik.wynik === undefined) {
    odmowa(wynik.blad, 'Rdzeń odmówił odrzucenia pozycji.');
    return;
  }
  powiedz(`Odrzucono ${String(wynik.wynik.rejected)}; wykluczonych w przesiewie: `
    + `${String(wynik.wynik.counts.excluded)}.`);
}
