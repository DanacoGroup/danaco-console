import {
  AgentConnectorKind,
  type AccessPoint,
  type AgentConnector,
} from '../../../../shared/contract';
import { opisOdmowy } from '../../komponenty/odmowa';
import {
  poleTekstowe,
  poleWielowierszowe,
  poleWyboru,
  przycisk,
  ustawPozycje,
  utworzWierszOdpowiedzi,
} from '../../modele/kontrolki-formularza';
import type { Kanal } from '../../protokol/kanal';
import { utworzPokrycieKomend, type PokrycieKomend } from '../pokrycie-komend';
import { utworzZrodloZakresuEksperta, type ZrodloZakresuEksperta } from './zrodlo-zakresu-eksperta';
import { naglowek } from './biblioteka-ekspertow';
import { utworzLicznikNarzedzi, type LicznikNarzedzi } from './licznik-narzedzi';
import { utworzStanOkna, type StanOkna } from './stan-okna';
import type { StanAgentow } from './stan-agentow';
import { utworzWykazWtyczek, type WykazWtyczek } from './wykaz-wtyczek';
import type { ZrodloZaplecza } from './zrodlo-zaplecza';

/** Pięć komend katalogu rozszerzeń wymienionych w kolumnie „Komendy” wykazu: odczyt, instalacja, konfiguracja, przełączenie i usunięcie pozycji. */
const KOMENDY_ROZSZERZEN = [
  'extension.list',
  'extension.install',
  'extension.configure',
  'extension.toggle',
  'extension.uninstall',
] as const;

/**
 * Connectors Manager to okno zarządca modułu Agents: łączy wykaz mostów protokołu MCP
 * podłączonych do eksperta z licznikiem narzędzi oraz osobnym wykazem wtyczek tego samego
 * eksperta.
 */
export interface OknoConnectorsManager {
  element: HTMLElement;
  odswiez(): void;
  /** Odczytuje katalog mostów MCP z rdzenia. */
  wczytaj(): Promise<void>;
  /** Odpina wykaz pokrycia od wspólnego bytu. Wołane z `rozlacz()` modułu. */
  zamknij(): void;
}

/** Zależności okna: baner zakresu przenosi ognisko do okna modułu Permissions Center, wskazanego przekazanym kodem okna. */
export interface OpcjeKonektorow {
  /** Przenosi ognisko do okna modułu wskazanego jego kodem. */
  naOkno(kod: string): void;
}

export function utworzOknoConnectorsManager(
  stan: StanAgentow,
  zaplecze: ZrodloZaplecza,
  kanal: Kanal,
  opcje: OpcjeKonektorow,
): OknoConnectorsManager {
  const okno: StanOkna = utworzStanOkna();
  const pokrycie: PokrycieKomend = utworzPokrycieKomend(kanal);
  // Trzy komendy dopełniające okno jadą źródłem zakresu eksperta, tym samym co Permissions Center.
  const zakres: ZrodloZakresuEksperta = utworzZrodloZakresuEksperta(kanal);
  /** Definicje konektorów wczytane komendą `agent.connector.list`. */
  let definicje = new Map<string, AgentConnector>();
  // Wykaz wtyczek stoi w tym oknie, bo dotyczy tego samego eksperta, ale ma osobne komendy.
  const wtyczki: WykazWtyczek = utworzWykazWtyczek(stan);
  const licznik: LicznikNarzedzi = utworzLicznikNarzedzi();

  const podlaczone = document.createElement('ul');
  podlaczone.className = 'da-konektory';

  const nazwa = poleTekstowe({ etykieta: 'Nazwa konektora', podpowiedz: 'np. Most danaco-system' });
  const rodzaj = poleWyboru(
    { etykieta: 'Rodzaj konektora' },
    Object.values(AgentConnectorKind).map((wartosc) => ({ wartosc, etykieta: wartosc })),
  );
  const most = poleWyboru(
    {
      etykieta: 'Most z katalogu punktów dostępu',
      opis: 'Konektor rodzaju mcp wymaga wskazania mostu — rdzeń odmówi bez niego.',
    },
    [],
  );
  const konfiguracja = poleWielowierszowe(
    { etykieta: 'Konfiguracja integracji (JSON)', podpowiedz: '{"scope": "odczyt"}' },
    3,
  );

  const podlacz = przycisk('Podłącz konektor', 'dn-btn dn-btn--sm dn-btn--atrament');
  const odpowiedz = utworzWierszOdpowiedzi();

  const granica = document.createElement('p');
  granica.className = 'dn-pole-opis da-granica';
  granica.textContent =
    'Wykaz pokazuje nazwę, rodzaj i punkt dostępu każdego konektora — dociąga je ' +
    'agent.connector.list, bo sam byt eksperta niesie wyłącznie identyfikatory. ' +
    'Przy wierszu stoi odłączenie i przełącznik czynności instancji.';

  // Baner rozdziela dwie decyzje: podłączenie rozszerzenia do eksperta i zakres jego wykorzystania.
  const doUprawnien = document.createElement('div');
  doUprawnien.className = 'da-baner da-baner--informacja';
  doUprawnien.setAttribute('role', 'note');

  const trescBaneru = document.createElement('p');
  trescBaneru.className = 'da-baner__tresc';
  trescBaneru.textContent =
    'To okno rozstrzyga wyłącznie, CZY rozszerzenie jest podłączone do definicji tego ' +
    'eksperta. W jakim ZAKRESIE ekspert może je wykorzystać — odczyt, zapis, wywołanie ' +
    'akcji — rozstrzyga Permissions Center.';

  const przejscie = przycisk('Przejdź do Permissions Center', 'dn-btn dn-btn--sm dn-btn--zarys');
  przejscie.addEventListener('click', () => opcje.naOkno('permissions-center'));
  doUprawnien.append(trescBaneru, przejscie);

  // Pokrycie katalogu rozszerzeń bierze się z bytu pokrycia, nie z napisu na sztywno.
  const rozszerzenia = pokrycie.wykaz(KOMENDY_ROZSZERZEN, 'da-granica da-rozszerzenia');

  okno.tresc.append(
    podlaczone,
    nazwa.element,
    rodzaj.element,
    most.element,
    konfiguracja.element,
    podlacz,
    odpowiedz.element,
    granica,
    doUprawnien,
    rozszerzenia,
    wtyczki.element,
  );

  const element = document.createElement('section');
  element.className = 'da-okno da-okno--zarzadca';
  element.dataset['okno'] = 'connectors-manager';
  element.append(naglowek('Connectors Manager'), licznik.element, okno.element);

  let mosty: AccessPoint[] = [];

  async function podlaczenie(): Promise<void> {
    const ekspert = stan.wybrany();
    if (ekspert === null) {
      odpowiedz.pokaz('Wybierz eksperta w Agent Builderze — konektor należy do jednego.', false);
      return;
    }
    if (nazwa.kontrolka.value.trim() === '') {
      odpowiedz.pokaz('Konektor wymaga nazwy — rdzeń odmówi podłączenia bez niej.', false);
      return;
    }
    odpowiedz.pokaz('Podłączanie konektora w toku…', true);
    const wynik = await stan.zrodlo.dodajKonektor({
      idEksperta: ekspert.id,
      nazwa: nazwa.kontrolka.value,
      rodzaj: rodzaj.kontrolka.value as AgentConnectorKind,
      idPunktu: most.kontrolka.value,
      konfiguracja: konfiguracja.kontrolka.value,
    });
    if (!wynik.udany || wynik.wynik === undefined) {
      odpowiedz.pokaz(
        opisOdmowy('Podłączenie konektora', wynik.blad?.code, wynik.blad?.message),
        false,
      );
      return;
    }
    const konektor = wynik.wynik.connector;
    // Formularz czyścimy w całości: zostawiona konfiguracja poszłaby do rdzenia przy podłączeniu.
    nazwa.kontrolka.value = '';
    konfiguracja.kontrolka.value = '';
    await stan.odswiez();
    odpowiedz.pokaz(
      `Konektor „${konektor.name}" (${konektor.kind}) podłączony pod identyfikatorem ${konektor.id}.`,
      true,
    );
  }

  podlacz.addEventListener('click', () => void podlaczenie());

  /** Ekspert, którego konektory stoją w oknie — po nim poznajemy zmianę wyboru. */
  let pokazany = '';

  function odswiez(): void {
    const ekspert = stan.wybrany();
    licznik.ustaw(ekspert);
    if ((ekspert?.id ?? '') !== pokazany) {
      pokazany = ekspert?.id ?? '';
      nazwa.kontrolka.value = '';
      konfiguracja.kontrolka.value = '';
      odpowiedz.wyczysc();
    }
    wtyczki.ustaw(ekspert);
    if (ekspert === null) {
      podlaczone.replaceChildren();
      okno.puste('Wybierz eksperta w Agent Builderze, aby zarządzać jego konektorami.');
      return;
    }
    const kody = ekspert.connectorIds ?? [];
    podlaczone.replaceChildren(
      ...kody.map((kod) =>
        wiersz(
          kod,
          definicje.get(kod),
          (wskazany) => void odlacz(wskazany),
          (wskazany, czynny) => void przestawCzynnosc(wskazany, czynny),
        ),
      ),
    );
    void wczytajDefinicje(ekspert.id);
    if (kody.length === 0) {
      okno.puste(`Ekspert „${ekspert.name}" nie ma jeszcze podłączonego konektora.`);
      return;
    }
    okno.gotowe();
  }

  /** Dociąganie definicji konektorów. Byt eksperta niesie same identyfikatory, bez nazw i rodzajów. */
  async function wczytajDefinicje(idEksperta: string): Promise<void> {
    const wynik = await zakres.konektory(idEksperta);
    if (!wynik.udany || wynik.wynik === undefined) {
      odpowiedz.pokaz(
        opisOdmowy('Odczyt definicji konektorów', wynik.blad?.code, wynik.blad?.message),
        false,
      );
      return;
    }
    const zebrane = new Map<string, AgentConnector>();
    for (const konektor of wynik.wynik.connectors) zebrane.set(konektor.id, konektor);
    // Porównanie po liczbie i kluczach chroni przed pętlą: przerysowanie wywołuje kolejny odczyt.
    const bezZmiany =
      zebrane.size === definicje.size &&
      [...zebrane.keys()].every((kod) => definicje.has(kod));
    definicje = zebrane;
    if (!bezZmiany) odswiez();
  }

  /** Odłączenie konektora od definicji eksperta (`agent.connector.remove`). */
  async function odlacz(idKonektora: string): Promise<void> {
    const ekspert = stan.wybrany();
    if (ekspert === null) return;
    odpowiedz.pokaz(`Odłączanie konektora ${idKonektora}…`, true);
    const wynik = await zakres.odlaczKonektor(ekspert.id, idKonektora);
    if (!wynik.udany || wynik.wynik === undefined) {
      odpowiedz.pokaz(
        opisOdmowy('Odłączenie konektora', wynik.blad?.code, wynik.blad?.message),
        false,
      );
      return;
    }
    definicje.delete(idKonektora);
    stan.wchlon(wynik.wynik.agent);
    odpowiedz.pokaz(
      `Konektor odłączony. Zakres wykorzystania zapisany w uprawnieniach ZOSTAJE — ` +
        'na wypadek ponownego podłączenia tego samego rozszerzenia.',
      true,
    );
  }

  /** Przestawienie czynności instancji konektora (`agent.connector.configure`). */
  async function przestawCzynnosc(idKonektora: string, czynny: boolean): Promise<void> {
    const ekspert = stan.wybrany();
    if (ekspert === null) return;
    odpowiedz.pokaz(`Zapis czynności konektora ${idKonektora}…`, true);
    const wynik = await zakres.skonfigurujKonektor({
      idEksperta: ekspert.id,
      idKonektora,
      czynny,
    });
    if (!wynik.udany || wynik.wynik === undefined) {
      // Przełącznik wraca do stanu rdzenia — przestawiony przez przeglądarkę nie świadczy o zapisie.
      void wczytajDefinicje(ekspert.id);
      odswiez();
      odpowiedz.pokaz(
        opisOdmowy('Zapis konfiguracji konektora', wynik.blad?.code, wynik.blad?.message),
        false,
      );
      return;
    }
    definicje.set(idKonektora, wynik.wynik.connector);
    odswiez();
    odpowiedz.pokaz(
      `Konektor „${wynik.wynik.connector.name}" ${wynik.wynik.connector.enabled ? 'czynny' : 'wstrzymany'}.`,
      true,
    );
  }

  return {
    element,

    odswiez,

    zamknij: () => pokrycie.zamknij(),

    async wczytaj() {
      // Powitanie idzie raz na połączenie, więc odczyt wykazu komend można zlecić swobodnie.
      void pokrycie.odczytaj();
      okno.ladowanie('Odczyt katalogu mostów MCP w toku…');
      const wynik = await zaplecze.mosty();
      if (!wynik.udany || wynik.wynik === undefined) {
        okno.blad(opisOdmowy('Odczyt katalogu mostów', wynik.blad?.code, wynik.blad?.message));
        return;
      }
      mosty = wynik.wynik.points;
      ustawPozycje(most.kontrolka, [
        { wartosc: '', etykieta: mosty.length === 0 ? 'katalog mostów pusty' : 'bez mostu' },
        ...mosty.map((punkt) => ({
          wartosc: punkt.id,
          etykieta: `${punkt.name} · ${punkt.status}`,
        })),
      ]);
      odswiez();
    },
  };
}

/**
 * Wiersz podłączonego konektora wraz z dwiema kontrolkami bez pokrycia w kontrakcie:
 * odłączeniem i przełącznikiem czynności instancji, widocznymi jeszcze przed odczytem
 * definicji.
 */
function wiersz(
  kod: string,
  definicja: AgentConnector | undefined,
  odlacz: (kod: string) => void,
  przestaw: (kod: string, czynny: boolean) => void,
): HTMLElement {
  const nazwa = document.createElement('span');
  nazwa.className = 'da-konektory__kod';
  // Sam identyfikator nic nie mówi Operatorowi; nazwa i rodzaj dochodzą osobnym odczytem definicji.
  nazwa.textContent =
    definicja === undefined ? `${kod} — definicja w odczycie…` : `${definicja.name} (${definicja.kind})`;

  const most = document.createElement('span');
  most.className = 'da-konektory__most';
  most.textContent =
    definicja?.accessPointId === undefined || definicja.accessPointId === ''
      ? 'bez wskazania mostu'
      : `most: ${definicja.accessPointId}`;

  const czynny = document.createElement('input');
  czynny.type = 'checkbox';
  czynny.className = 'dn-przelacznik';
  czynny.checked = definicja?.enabled ?? true;
  czynny.disabled = definicja === undefined;
  czynny.title = `czynność instancji konektora ${kod}`;
  czynny.addEventListener('change', () => przestaw(kod, czynny.checked));

  const odlaczenie = przycisk('Odłącz', 'dn-btn dn-btn--sm dn-btn--zarys');
  odlaczenie.title = `odłączenie konektora ${kod} od definicji tego eksperta`;
  odlaczenie.addEventListener('click', () => odlacz(kod));

  const element = document.createElement('li');
  element.className = 'da-konektory__wiersz';
  element.dataset['konektor'] = kod;
  element.append(nazwa, most, czynny, odlaczenie);
  return element;
}
