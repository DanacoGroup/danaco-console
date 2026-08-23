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

/** Pięć komend katalogu rozszerzeń wymienionych w kolumnie „Komendy” wykazu. */
const KOMENDY_ROZSZERZEN = [
  'extension.list',
  'extension.install',
  'extension.configure',
  'extension.toggle',
  'extension.uninstall',
] as const;

/**
 * Connectors Manager — okno zarządca modułu Agents.
 *
 * Lista serwerów MCP pochodzi z `access.point.list` ograniczonej do rodzaju
 * `mcpBridge` — z tego samego katalogu punktów dostępu, z którego rdzeń składa
 * plik `mcpServers` procesu modelu. Okno nie zakłada drugiego rejestru mostów
 * i nie zna ani jednego adresu maszyny: podaje kod punktu, resztę wie rdzeń.
 *
 * Katalogiem rozszerzeń zarządza osobne okno modułu
 * (`okno-katalog-rozszerzen.ts`, komplet pięciu komend `extension.*`); tutaj
 * stan tych komend bierze się na żywo z bytu pokrycia (`../pokrycie-komend`),
 * a samo okno robi `agent.connector.add` i wtyczki.
 *
 * Licznik narzędzi stoi także tutaj, bo serwer narzędzi czyta `Agent.skillIds`
 * i `Agent.connectorIds` jako jeden zbiór kodów
 * (`server/internal/narzedzia/ekspert_definicja.go` → `DefinicjaEksperta.Kody`)
 * — podłączenie konektora zmienia tę samą liczbę, którą pokazuje Agent Builder.
 *
 * Wtyczki stoją osobno od konektorów: wykaz wtyczek (`wykaz-wtyczek.ts`) ma
 * własny panel i własne komendy `agent.plugin.add` oraz `agent.plugin.remove`.
 * Konektor jest drogą do usługi (`--mcp-config`), wtyczka katalogiem rozszerzeń
 * powłoki (`--plugin-dir`).
 */
export interface OknoConnectorsManager {
  element: HTMLElement;
  odswiez(): void;
  /** Odczytuje katalog mostów MCP z rdzenia. */
  wczytaj(): Promise<void>;
  /** Odpina wykaz pokrycia od wspólnego bytu. Wołane z `rozlacz()` modułu. */
  zamknij(): void;
}

/** Zależności okna: baner zakresu przenosi ognisko do Permissions Center. */
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
  // Trzy komendy dopełniające okno — odczyt definicji, odłączenie i konfiguracja
  // instancji — jadą źródłem zakresu eksperta, tym samym, którym jedzie
  // Permissions Center.
  const zakres: ZrodloZakresuEksperta = utworzZrodloZakresuEksperta(kanal);
  /** Definicje konektorów wczytane komendą `agent.connector.list`. */
  let definicje = new Map<string, AgentConnector>();
  // Wtyczka nie jest odmianą konektora. Wykaz stoi w tym oknie, bo dotyczy tej
  // samej rzeczy — tego, co ekspert dostaje ponad model — ale jest osobnym
  // panelem z osobnymi komendami `agent.plugin.*`.
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

  // Baner rozdziela dwie decyzje, które opracowanie każe trzymać osobno: TUTAJ
  // rozstrzyga się, czy rozszerzenie jest podłączone do definicji eksperta;
  // W JAKIM ZAKRESIE ekspert może je wykorzystać — w Permissions Center.
  // Bez tego zdania przełącznik podłączenia czytałoby się jako nadanie
  // pełnego dostępu, a to dwie różne rzeczy.
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

  // Pokrycie katalogu rozszerzeń bierze się z bytu pokrycia zamiast z napisu na
  // sztywno — zdanie mówi to, co rdzeń orzekł powitaniem, i przerysuje się samo,
  // gdy rdzeń te komendy doda.
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
    // Formularz czyścimy w całości: konfiguracja zostawiona w polu poszłaby do
    // rdzenia przy następnym podłączeniu, także przy innym ekspercie.
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

  /**
   * Dociąganie definicji konektorów (`agent.connector.list`).
   *
   * Byt eksperta niesie same identyfikatory, więc bez tego odczytu wykaz zna
   * liczbę konektorów i ani jednej nazwy. Odczyt idzie po narysowaniu wierszy
   * i przerysowuje je, gdy wróci — wiersz z kodem widać od razu, a nie po
   * odpowiedzi rdzenia.
   *
   * Nieudany odczyt nie czyści wykazu: wiersze zostają z kodami, a odmowa
   * wraca w wierszu odpowiedzi. Zniknięcie konektorów z ekranu byłoby
   * nieprawdą o definicji eksperta.
   */
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
    // Porównanie po liczbie i kluczach chroni przed pętlą: przerysowanie
    // wywołuje odczyt, więc bez tego warunku okno pytałoby rdzeń bez końca.
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
      // Przełącznik wraca do stanu rdzenia — przestawiony przez przeglądarkę
      // nie ma prawa zostać świadectwem zapisu, którego nie było.
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
      // Powitanie idzie raz na połączenie (pamięć podręczna bytu pokrycia po
      // kanale), więc odczyt wykazu komend można zlecić swobodnie.
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
 * Wiersz podłączonego konektora wraz z dwiema kontrolkami bez pokrycia.
 *
 * Wiersz niesie sam identyfikator, bo tyle oddaje byt eksperta. Odczyt definicji
 * i odłączenie mają komendy w kontrakcie, ale nie mają jeszcze uchwytu
 * w rdzeniu — obie kontrolki zostają więc widoczne i klikalne, i po naciśnięciu
 * nazywają stan komendy. Wiersz z samym napisem wyglądałby na skończony.
 */
function wiersz(
  kod: string,
  definicja: AgentConnector | undefined,
  odlacz: (kod: string) => void,
  przestaw: (kod: string, czynny: boolean) => void,
): HTMLElement {
  const nazwa = document.createElement('span');
  nazwa.className = 'da-konektory__kod';
  // Sam identyfikator nie mówi Operatorowi nic. Nazwa i rodzaj pochodzą
  // z `agent.connector.list`; dopóki odczyt nie wrócił, wiersz niesie kod —
  // i mówi wprost, że definicji jeszcze nie ma, zamiast udawać pustkę.
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
