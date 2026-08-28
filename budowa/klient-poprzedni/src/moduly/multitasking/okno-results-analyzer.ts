import { ConfigScope, type Message, type MonitorStatus } from '../../../../shared/contract';
import { utworzRameOkna } from '../../komponenty/rama-okna';
import { poleTresci, przyciskAkcji as przycisk, wybor } from '../../modele/kontrolki-formularza';
import type { WykazKomendRdzenia } from './braki-kontraktu';
import { WIERSZE_POLA } from './kontrolki';
import { KOD_ANALITYKA } from './obsada-rol';
import { zestawienie, type WynikWykonawcy } from './porownanie-wynikow';
import { utworzStanTresci, type StanTresci } from './stany-okna';
import {
  KLUCZ_WCIELENIA,
  opisWcielenia,
  pozycjeWcielen,
  trescZgloszenia,
  Wcielenie,
} from './wcielenia-analizy';
import type { StanMultitaskingu } from './stan-multitaskingu';
import type { ZrodloBiegu } from './zrodlo-biegu';
import type { ZrodloOkien } from './zrodlo-okien';

/**
 * Results Analyzer jest oknem monitora zbierającym wyniki obu wykonawców:
 * melduje stan procesów telemetrii i zgłasza werdykt do okna koordynatora, a
 * wcielenie roli analityka utrwala się na poziomie okna.
 */
export interface OknoAnalityka {
  element: HTMLElement;
  /** Odczytuje wyniki wykonawców z rdzenia i przerysowuje zestawienie. */
  odswiez(): void;
}

export interface OpcjeAnalityka {
  okna: ZrodloOkien;
  bieg: ZrodloBiegu;
  stan: StanMultitaskingu;
  /** Wykaz komend rdzenia — stąd bierze się powód nieczynnych kontrolek. */
  komendy: WykazKomendRdzenia;
}

/** Ile ostatnich wiadomości danego okna wystarcza, żeby na pewno znaleźć w nich wynik pracy tego wykonawcy. */
const GLEBOKOSC_ODCZYTU = 20;

export function utworzOknoAnalityka(opcje: OpcjeAnalityka): OknoAnalityka {
  const { stan } = opcje;

  const rama = utworzRameOkna({
    tytul: 'Results Analyzer',
    rola: 'monitor',
    kod: KOD_ANALITYKA,
    przeznaczenie:
      'Ocena wyników pracy wykonawców, zgłoszenie niezgodności i rozstrzygnięcie konfliktu.',
    ogniskowalne: true,
  });

  const tresci = utworzStanTresci();
  const { uzasadnienie, wcielenie, kryteria, zatwierdz, odrzuc, rozstrzygnij, zestaw, stanMonitora } =
    zlozPowierzchnieAnalityka(rama, tresci.element, opcje.komendy);

  stanMonitora.addEventListener('click', () => {
    void odczytajStanMonitora();
  });

  wcielenie.addEventListener('change', () => {
    void zapiszWcielenie();
  });
  zatwierdz.addEventListener('click', () => {
    void zglos('wynik zatwierdzony');
  });
  odrzuc.addEventListener('click', () => {
    void zglos('niezgodność zgłoszona');
  });
  rozstrzygnij.addEventListener('click', () => {
    wcielenie.value = Wcielenie.Arbitrator;
    void zapiszWcielenie().then(() => zglos('konflikt rozstrzygnięty'));
  });

  // Zgłoszenie werdyktu do okna koordynatora; bez wskazanego koordynatora zgłoszenie nie ma adresata.
  async function zglos(werdykt: string): Promise<void> {
    const koordynator = stan.obsada().koordynator;
    if (koordynator === null) {
      tresci.blad('Obsada nie ma koordynatora — zgłoszenie nie ma adresata.');
      return;
    }
    const tresc = trescZgloszenia(opisWcielenia(stan.wcielenie()), werdykt, uzasadnienie.value);
    const wynik = await opcje.bieg.wyslij({ windowId: koordynator.id, content: tresc });
    // Zgłoszenie liczy się za przekazane dopiero, gdy rdzeń odda wiadomość zapisaną w cudzym oknie.
    if (!wynik.udany || wynik.wynik === undefined) {
      tresci.blad('Rdzeń odmówił przyjęcia zgłoszenia.', wynik.blad);
      return;
    }
    const zapisane = wynik.wynik;
    if (zapisane.windowId !== koordynator.id) {
      tresci.potwierdzenie(
        `Rdzeń zapisał zgłoszenie w oknie ${zapisane.windowId}, a szło do koordynatora ${koordynator.id} — werdykt trafił gdzie indziej.`,
        false,
      );
      return;
    }
    uzasadnienie.value = '';
    tresci.potwierdzenie(zdanieOZgloszeniu(zapisane, werdykt, koordynator.id), true);
  }

  async function zapiszWcielenie(): Promise<void> {
    const { zdanie, udane } = await utrwalWcielenie(opcje.okna, stan, wcielenie.value);
    tresci.potwierdzenie(zdanie, udane);
  }

  // Odczyt zbiorczego stanu procesów telemetrii; sitem jest sesja analityka i jego okno, jeśli wskazane.
  async function odczytajStanMonitora(): Promise<void> {
    const analityk = stan.obsada().analityk;
    tresci.ladowanie('Odczyt stanu monitora…');
    const wynik = await opcje.bieg.stanMonitora({
      sessionId: stan.sesja(),
      ...(analityk === null ? {} : { windowId: analityk.id }),
    });
    if (!wynik.udany || wynik.wynik === undefined) {
      tresci.blad('Rdzeń odmówił odczytu stanu monitora.', wynik.blad);
      return;
    }
    const statusy = wynik.wynik.statuses;
    if (statusy.length === 0) {
      tresci.potwierdzenie(
        'Rdzeń oddał stan monitora: brak procesów telemetrii dla tej sesji — nic obecnie nie biegnie.',
        true,
      );
      return;
    }
    tresci.potwierdzenie(
      `Rdzeń oddał stan ${statusy.length} procesów: ${statusy.map(opisProcesuMonitora).join('; ')}.`,
      true,
    );
  }

  function odswiez(): void {
    const opis = opisWcielenia(stan.wcielenie());
    wcielenie.value = opis.wcielenie;
    kryteria.replaceChildren(...pozycjeKryteriow(opis.kryteria));
    const zebrane = wynikiWykonawcow(stan);
    zestaw.replaceChildren(zestawienie(zebrane));
    const gotowe = zebrane.filter((pozycja) => pozycja.wiadomosc !== null).length;
    rama.ustawZnacznik(`${gotowe} z ${zebrane.length} wyników`, wagaKompletuWynikow(gotowe, zebrane.length));
  }

  stan.naZmiane(odswiez);
  odswiez();

  return {
    element: rama.element,
    odswiez() {
      odswiez();
      void wczytajWynikiWykonawcow(opcje.bieg, stan, tresci);
    },
  };
}

/** Kontrolki własne okna analityka wraz z miejscami na kryteria oceny, zestawienie wyników i uzasadnienie. */
interface PowierzchniaAnalityka {
  uzasadnienie: HTMLTextAreaElement;
  wcielenie: HTMLSelectElement;
  kryteria: HTMLElement;
  zatwierdz: HTMLButtonElement;
  odrzuc: HTMLButtonElement;
  rozstrzygnij: HTMLButtonElement;
  zestaw: HTMLElement;
  /** Przycisk czynny — woła `monitor.status` na żywym rdzeniu. */
  stanMonitora: HTMLButtonElement;
}

/**
 * Składa kontrolki, pasek akcji, narzędzia i ciało okna analityka.
 *
 * Sama konstrukcja węzłów — nasłuchy zdarzeń zakłada wytwórnia okna.
 */
function zlozPowierzchnieAnalityka(
  rama: { akcje: HTMLElement; narzedzia: HTMLElement; cialo: HTMLElement },
  stanTresci: HTMLElement,
  komendy: WykazKomendRdzenia,
): PowierzchniaAnalityka {
  const uzasadnienie = poleTresci(
    'Uzasadnienie oceny',
    WIERSZE_POLA,
    'Na czym polega niezgodność albo dlaczego wynik przyjęto…',
  );
  const wcielenie = wybor('Wcielenie roli analityka', pozycjeWcielen());

  const kryteria = document.createElement('ul');
  kryteria.className = 'dm-kryteria';
  kryteria.setAttribute('aria-label', 'Kryteria oceny bieżącego wcielenia');

  const zatwierdz = przycisk('Zatwierdź', 'dn-btn dn-btn--atrament');
  const odrzuc = przycisk('Odrzuć / zgłoś niezgodność', 'dn-btn dn-btn--niebezpieczny');
  const rozstrzygnij = przycisk('Rozstrzygnij (Arbitrator)');
  const stanMonitora = przycisk('Stan monitora', 'dn-btn dn-btn--zarys');

  const zestaw = document.createElement('div');
  zestaw.className = 'dm-zestaw';

  rama.akcje.append(zatwierdz, odrzuc, rozstrzygnij);
  rama.narzedzia.append(wcielenie, stanMonitora);
  // W wykazie braków zostaje tylko zmiana roli: komendy monitora rdzeń rejestruje, więc nie są brakiem.
  rama.cialo.append(kryteria, zestaw, uzasadnienie, stanTresci, komendy.wykazBrakow([
    'role.update',
  ]));

  return { uzasadnienie, wcielenie, kryteria, zatwierdz, odrzuc, rozstrzygnij, zestaw, stanMonitora };
}

/** Jeden proces w meldunku stanu monitora, złożony z tego, co oddał o nim rdzeń: etap, postęp i jego stan. */
function opisProcesuMonitora(status: MonitorStatus): string {
  const etap =
    status.stageIndex !== undefined && status.stageCount !== undefined
      ? `, etap ${status.stageIndex + 1}/${status.stageCount}`
      : status.stage !== undefined
        ? `, etap ${status.stage}`
        : '';
  const ukonczenie = status.completion !== undefined ? `, ${status.completion}%` : '';
  return `${status.label ?? status.processId} (${status.status}${etap}${ukonczenie})`;
}

/** Pozycje wykazu kryteriów bieżącego wcielenia — czysta zamiana zdań kryteriów na wiersze listy widoku. */
function pozycjeKryteriow(kryteria: readonly string[]): HTMLElement[] {
  return kryteria.map((pozycja) => {
    const wiersz = document.createElement('li');
    wiersz.textContent = pozycja;
    return wiersz;
  });
}

/** Waga plakietki kompletu wyników: brak wyników ostrzega, komplet oznacza sukces, reszta jest neutralna. */
function wagaKompletuWynikow(gotowe: number, wszystkich: number): 'neutralna' | 'sukces' | 'ostrzezenie' {
  return gotowe === 0 ? 'ostrzezenie' : gotowe === wszystkich ? 'sukces' : 'neutralna';
}

/** Wyniki obu wykonawców złożone do zestawienia; sam odczyt stanu wspólnego, bez sięgania do rdzenia po nic. */
function wynikiWykonawcow(stan: StanMultitaskingu): readonly WynikWykonawcy[] {
  return stan.obsada().wykonawcy.map((okno, miejsce) => ({
    numer: miejsce + 1,
    okno,
    wiadomosc: stan.wynikWykonawcy(okno.id),
  }));
}

/**
 * Utrwalenie wcielenia na poziomie okna analityka.
 *
 * Wynik oddaje zdaniem dla stanu treści, zamiast sięgać po niego samodzielnie.
 */
async function utrwalWcielenie(
  okna: ZrodloOkien,
  stan: StanMultitaskingu,
  wartosc: string,
): Promise<{ zdanie: string; udane: boolean }> {
  const poprzednie = stan.wcielenie();
  const nowe = opisWcielenia(wartosc).wcielenie;
  stan.ustawWcielenie(nowe);
  const analityk = stan.obsada().analityk;
  if (analityk === null) {
    return {
      zdanie: `Wcielenie „${nowe}" trzyma WYŁĄCZNIE ten widok — okno analityka nie jest wskazane, więc rdzeń nie ma na czym go zapisać.`,
      udane: false,
    };
  }
  const wynik = await okna.zapiszUstawienie({
    key: KLUCZ_WCIELENIA,
    value: nowe,
    scope: ConfigScope.Window,
    scopeId: analityk.id,
  });
  // Odmowa cofa widok: selektor pokazujący niezapisane wcielenie sugerowałby utrwalone kryteria oceny.
  if (!wynik.udany || wynik.wynik === undefined) {
    stan.ustawWcielenie(poprzednie);
    return {
      zdanie: `Rdzeń odmówił zapisu wcielenia, widok wraca do „${poprzednie}". Powód: ${wynik.blad?.message ?? 'rdzeń nie podał przyczyny'} (kod ${wynik.blad?.code ?? 'brak'}).`,
      udane: false,
    };
  }
  const zapisane = opisWcielenia(wynik.wynik.value).wcielenie;
  stan.ustawWcielenie(zapisane);
  return {
    zdanie: `Rdzeń zapisał wcielenie „${zapisane}" na oknie analityka ${analityk.id}.`,
    udane: true,
  };
}

/** Zdanie o przekazanym zgłoszeniu — złożone z wiadomości, którą oddał rdzeń po zapisie w oknie koordynatora. */
function zdanieOZgloszeniu(wiadomosc: Message, werdykt: string, koordynator: string): string {
  return `Rdzeń zapisał zgłoszenie „${werdykt}" w oknie koordynatora ${koordynator} jako wiadomość ${wiadomosc.id} (stan ${wiadomosc.status}).`;
}

/** Pierwsze wypełnienie wyników wykonawców przy otwarciu okna; dalej żywi je zdarzenie zmiany wiadomości. */
async function wczytajWynikiWykonawcow(
  bieg: ZrodloBiegu,
  stan: StanMultitaskingu,
  tresci: StanTresci,
): Promise<void> {
  for (const okno of stan.obsada().wykonawcy) {
    const wynik = await bieg.wiadomosci({ windowId: okno.id, limit: GLEBOKOSC_ODCZYTU });
    if (!wynik.udany || wynik.wynik === undefined) {
      tresci.blad(`Odczyt wiadomości okna ${okno.id} odmówiony.`, wynik.blad);
      continue;
    }
    stan.ustawWiadomosci(okno.id, wynik.wynik);
  }
}
