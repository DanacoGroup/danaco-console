import { QueueAction, type ErrorInfo } from '../../../../shared/contract';
import { przyciskAkcji as przycisk } from '../../modele/kontrolki-formularza';
import { rozdziel, type Zlecenie } from './tryby-wspolpracy';
import type { StanMultitaskingu } from './stan-multitaskingu';
import type { ZrodloBiegu } from './zrodlo-biegu';

/**
 * Siedem przycisków sterowania Coordinator Chat: Start · Stop · Pauza · Wznów ·
 * Przekaż · Powtórz · Waliduj.
 *
 * Pięć z nich to jeden silnik kolejek: Start, Stop, Pauza, Wznów i Powtórz idą
 * `queue.action` na kolejce etapu — tej samej, którą obsługuje pętla sesyjna
 * i moduł Automations. Koordynator nie ma własnego wykonawcy zleceń i nie
 * prowadzi biegu naprawczego: bieg prowadzi rdzeń, a licznik obiegów przychodzi
 * z `window.state.get`.
 *
 * „Przekaż" wydaje zlecenie wykonawcy. Adresata i treść rozstrzyga tryb
 * współpracy; nośnikiem jest `message.send`, bo okno wykonawcy jest oknem
 * komunikacji. Droga modelu — wywołanie narzędzia platformy — biegnie obok
 * i tej kontrolki nie zastępuje.
 *
 * „Waliduj" sprawdza przebieg, a nie ocenia wynik: kontrakt nie ma komendy
 * oceny, więc przycisk odczytuje `window.state.get` koordynatora i wykonawców
 * i mówi, czy tura trwa, ile obiegów zliczono i czy bieg stoi wraz z powodem.
 * Ocena treści wyniku należy do Results Analyzer.
 */

/** Zestaw sterowania wraz z przyciskami, które da się wygasić. */
export interface SterowanieKoordynatora {
  element: HTMLElement;
  /** Przestawia dostępność przycisków wedle obsady i kolejki. */
  odswiez(): void;
}

/** Zależności sterowania. */
export interface OpcjeSterowania {
  zrodlo: ZrodloBiegu;
  stan: StanMultitaskingu;
  /** Treść zlecenia wpisana przez Operatora w kreatorze promptu. */
  trescZlecenia(): string;
  /** Odczyt stanu okien na żądanie — nośnik walidacji przebiegu. */
  zwaliduj(): void;
  /** Krótkie potwierdzenie czynności widoczne w oknie. */
  potwierdz(zdanie: string, udane: boolean): void;
}

export function utworzSterowanie(opcje: OpcjeSterowania): SterowanieKoordynatora {
  const { zrodlo, stan } = opcje;

  const { element, naKolejce, przekaz, waliduj } = zlozPasSterowaniaKoordynatora({
    kolejka: (dzialanie) => {
      void steruj(dzialanie);
    },
    stop: () => {
      void zatrzymajBiegKoordynatora(zrodlo, stan, opcje.potwierdz);
    },
    przekazanie: () => {
      void przekazZlecenie();
    },
    walidacja: () => opcje.zwaliduj(),
  });

  /**
   * Sterowanie kolejką etapu — jedno wywołanie, jeden komunikat zwrotny.
   *
   * Potwierdzenie mówi, co zrobił rdzeń, a nie co zażądano: odpowiedź
   * `queue.action` niesie całą kolejkę wraz z jej stanem po działaniu, więc
   * zdanie zestawia stan sprzed wywołania ze stanem oddanym.
   */
  async function steruj(dzialanie: QueueAction): Promise<void> {
    const kolejka = stan.kolejkaBiezaca();
    if (kolejka === null) {
      opcje.potwierdz('Nie ma etapu do sterowania — załóż etap w planie.', false);
      return;
    }
    const stanPrzed = kolejka.status;
    const wynik = await zrodlo.sterujKolejka({ queueId: kolejka.id, action: dzialanie });
    if (!wynik.udany || wynik.wynik === undefined) {
      opcje.potwierdz(`Rdzeń odmówił działania ${dzialanie}. ${powod(wynik.blad)}`, false);
      return;
    }
    const po = wynik.wynik;
    stan.zapiszKolejke(po);
    opcje.potwierdz(
      po.status === stanPrzed
        ? `Etap „${po.name ?? po.id}" pozostał w stanie ${po.status} — rdzeń przyjął żądanie, ale stan kolejki się nie zmienił.`
        : `Etap „${po.name ?? po.id}": rdzeń przestawił kolejkę ze stanu ${stanPrzed} na ${po.status}.`,
      true,
    );
  }

  /**
   * Przekazanie zlecenia wykonawcom wedle trybu współpracy.
   *
   * Doręczenie treści i utrwalenie więzi to dwa osobne niepowodzenia, więc
   * wracają osobnymi wykazami i liczone są w oknach, nie w powodach.
   */
  async function przekazZlecenie(): Promise<void> {
    const tresc = opcje.trescZlecenia().trim();
    if (tresc === '') {
      opcje.potwierdz('Zlecenie puste — koordynator nie przekazuje pustej tury.', false);
      return;
    }
    const zlecenia = rozdziel(stan.tryb(), tresc, {
      wykonawcy: stan.obsada().wykonawcy.map((okno) => okno.id),
      poprzedni: stan.poprzedniWykonawca(),
      wynik: (okno) => stan.strumien().wynik(okno),
    });
    if (zlecenia.length === 0) {
      opcje.potwierdz('Obsada nie ma wykonawcy — załóż okno wykonawcy.', false);
      return;
    }
    const koordynator = stan.obsada().koordynator;
    const przebieg = await wyslijZleceniaWykonawcom(zrodlo, zlecenia, koordynator?.id ?? '');
    const ostatnieDoreczone = przebieg.doreczone[przebieg.doreczone.length - 1];
    if (ostatnieDoreczone !== undefined) stan.zapamietajPrzekazanie(ostatnieDoreczone);
    opcje.potwierdz(zdanieOPrzekazaniu(przebieg, zlecenia.length), przebieg.niedoreczone.length === 0);
  }

  function odswiez(): void {
    // Przyciski zostają czynne także bez przesłanki: każda z tych czynności ma
    // własną odmowę z powodem, a naciśnięcie nie dociera do rdzenia. Brakująca
    // przesłanka idzie więc znacznikiem, nie wygaszeniem.
    const brakEtapu = stan.kolejkaBiezaca() === null;
    for (const kontrolka of naKolejce) {
      oznaczCzynnosc(kontrolka, brakEtapu ? 'plan nie ma jeszcze etapu' : '');
    }
    oznaczCzynnosc(
      przekaz,
      stan.obsada().wykonawcy.length === 0 ? 'obsada nie ma wykonawcy' : '',
    );
    oznaczCzynnosc(waliduj, stan.obsada().koordynator === null ? 'obsada nie ma koordynatora' : '');
  }

  odswiez();
  return { element, odswiez };
}

/** Pas sterowania wraz z przyciskami, które wytwórnia gasi wedle stanu. */
interface PasSterowaniaKoordynatora {
  element: HTMLElement;
  /** Cztery przyciski jednego silnika kolejek — gasną, gdy nie ma etapu. */
  naKolejce: readonly HTMLButtonElement[];
  przekaz: HTMLButtonElement;
  waliduj: HTMLButtonElement;
}

/** Czynności pasa; pas ich nie wykonuje, tylko o nie prosi. */
interface CzynnosciSterowania {
  kolejka(dzialanie: QueueAction): void;
  stop(): void;
  przekazanie(): void;
  walidacja(): void;
}

/**
 * Składa siedem kontrolek pasa sterowania.
 *
 * Pas jest czystą konstrukcją: nie zna ani źródła biegu, ani stanu wspólnego —
 * każde kliknięcie oddaje wywołaniu zwrotnemu, które te dwa zna.
 */
function zlozPasSterowaniaKoordynatora(
  czynnosci: CzynnosciSterowania,
): PasSterowaniaKoordynatora {
  const element = document.createElement('div');
  element.className = 'dm-sterowanie';
  element.setAttribute('aria-label', 'Sterowanie koordynatora');

  const opisy = [
    ['Start', 'dn-btn dn-btn--atrament', QueueAction.Start],
    ['Pauza', 'dn-btn', QueueAction.Pause],
    ['Wznów', 'dn-btn', QueueAction.Resume],
    ['Powtórz', 'dn-btn', QueueAction.Retry],
  ] as const;

  const naKolejce = opisy.map(([etykieta, klasa, dzialanie]) => {
    const kontrolka = przycisk(etykieta, klasa);
    kontrolka.dataset['dzialanie'] = dzialanie;
    kontrolka.addEventListener('click', () => {
      czynnosci.kolejka(dzialanie);
    });
    element.append(kontrolka);
    return kontrolka;
  });

  // Zatrzymanie jest czynne także bez etapu i przed pierwszym obiegiem, dlatego
  // stoi poza wykazem sterowania kolejką: gasi turę wykonawcy przez
  // `message.stop` nawet wtedy, gdy żadna kolejka nie powstała.
  const stop = przycisk('Stop', 'dn-btn dn-btn--niebezpieczny');
  stop.dataset['dzialanie'] = QueueAction.Stop;
  stop.addEventListener('click', () => {
    czynnosci.stop();
  });
  element.append(stop);

  const przekaz = przycisk('Przekaż', 'dn-btn dn-btn--sygnal');
  przekaz.dataset['dzialanie'] = 'przekazanie';
  przekaz.addEventListener('click', () => {
    czynnosci.przekazanie();
  });

  const waliduj = przycisk('Waliduj');
  waliduj.dataset['dzialanie'] = 'walidacja';
  waliduj.addEventListener('click', () => {
    czynnosci.walidacja();
  });

  element.append(przekaz, waliduj);

  return { element, naKolejce, przekaz, waliduj };
}

/**
 * Zatrzymanie biegu: kolejka etapu i tury wykonawców naraz.
 *
 * Kolejność ma znaczenie — najpierw gaśnie tura, która właśnie zużywa czas
 * modelu, a dopiero potem kolejka, która mogłaby ją wznowić.
 *
 * Źródło biegu, stan wspólny i potwierdzenie przychodzą parametrami; poza tą
 * trójką funkcja nie sięga do wnętrza sterowania po nic.
 */
async function zatrzymajBiegKoordynatora(
  zrodlo: ZrodloBiegu,
  stan: StanMultitaskingu,
  potwierdz: (zdanie: string, udane: boolean) => void,
): Promise<void> {
  const zgaszone: string[] = [];
  const odmowy: string[] = [];
  for (const okno of stan.obsada().wykonawcy) {
    const wynik = await zrodlo.zatrzymaj({ windowId: okno.id });
    // Odmowa gaszenia idzie do osobnego wykazu, żeby zdanie „nic nie biegło"
    // nie zakryło tury, której rdzeń zatrzymać odmówił.
    if (!wynik.udany) odmowy.push(`okno ${okno.id}: ${powod(wynik.blad)}`);
    else if (wynik.wynik?.stopped === true) zgaszone.push(`tura okna ${okno.id}`);
  }
  const kolejka = stan.kolejkaBiezaca();
  if (kolejka !== null) {
    const wynik = await zrodlo.sterujKolejka({ queueId: kolejka.id, action: QueueAction.Stop });
    if (wynik.udany && wynik.wynik !== undefined) {
      stan.zapiszKolejke(wynik.wynik);
      zgaszone.push(`etap „${kolejka.name ?? kolejka.id}"`);
    } else if (!wynik.udany) {
      potwierdz(`Rdzeń odmówił zatrzymania etapu. ${powod(wynik.blad)}`, false);
      return;
    }
  }
  const zdanie =
    zgaszone.length === 0
      ? 'Rdzeń nie zgłosił zatrzymania niczego — nic nie biegło.'
      : `Rdzeń zatrzymał: ${zgaszone.join(', ')}.`;
  potwierdz(
    odmowy.length === 0 ? zdanie : `${zdanie} Odmowy zatrzymania (${odmowy.length}): ${odmowy.join(' ')}`,
    odmowy.length === 0,
  );
}

/** Przebieg przekazania rozdzielony na doręczenie i utrwalenie. */
interface PrzebiegPrzekazania {
  /** Okna, do których treść dotarła — wedle odpowiedzi `message.send`. */
  doreczone: string[];
  /** Okna, do których treść nie dotarła, wraz z powodem odmowy. */
  niedoreczone: string[];
  /** Okna doręczone, których więzi rdzeń nie zapisał, wraz z powodem. */
  nieutrwalone: string[];
}

/**
 * Przekazanie zlecenia to dwie czynności i obie muszą się odbyć.
 *
 * `message.send` doręcza treść oknu wykonawcy — bez tego wykonawca nie rusza do
 * pracy. `window.handoff` utrwala w bazie rdzenia, kto komu co zlecił, i zakłada
 * pozycję kolejki; bez niego więź ginie z zamknięciem przeglądarki, a Mission
 * Control nie ma czego pokazać po ponownym uruchomieniu.
 *
 * Kolejność ma znaczenie: najpierw zapis, potem doręczenie. Odwrotna zostawiałaby
 * wykonawcę pracującego nad zleceniem, którego nikt nie odnotował.
 *
 * Nieudany zapis nie wstrzymuje doręczenia i nie liczy się jako niedoręczenie.
 * Dwa rodzaje niepowodzenia wracają osobnymi wykazami, bo znaczą co innego:
 * pierwszy mówi „wykonawca nie dostał pracy", drugi „dostał, ale po ponownym
 * uruchomieniu nikt tego nie odtworzy".
 *
 * Źródło biegu przychodzi parametrem — o trybie współpracy, obsadzie i stanie
 * okna funkcja nie wie nic.
 */
async function wyslijZleceniaWykonawcom(
  zrodlo: ZrodloBiegu,
  zlecenia: readonly Zlecenie[],
  idKoordynatora: string,
): Promise<PrzebiegPrzekazania> {
  const przebieg: PrzebiegPrzekazania = { doreczone: [], niedoreczone: [], nieutrwalone: [] };
  for (const zlecenie of zlecenia) {
    if (idKoordynatora !== '') {
      const zapis = await zrodlo.przekaz({
        fromWindowId: idKoordynatora,
        toWindowId: zlecenie.okno,
        instruction: zlecenie.tresc,
      });
      if (!zapis.udany) przebieg.nieutrwalone.push(`${zlecenie.okno}: ${powod(zapis.blad)}`);
    }
    const wynik = await zrodlo.wyslij({ windowId: zlecenie.okno, content: zlecenie.tresc });
    // Dowodem doręczenia jest wiadomość, która wróciła, a nie sama zgoda:
    // odpowiedź bez treści znaczy, że rdzeń nie powiedział, co przyjął.
    if (wynik.udany && wynik.wynik !== undefined) przebieg.doreczone.push(zlecenie.okno);
    else przebieg.niedoreczone.push(`${zlecenie.okno}: ${powod(wynik.blad)}`);
  }
  return przebieg;
}

/**
 * Zdanie o przekazaniu — liczy okna i rozdziela doręczenie od utrwalenia.
 *
 * Czysta zamiana przebiegu na zdanie: nie zna ani obsady, ani rdzenia.
 */
function zdanieOPrzekazaniu(przebieg: PrzebiegPrzekazania, zleconych: number): string {
  const czesci: string[] = [];
  czesci.push(
    przebieg.doreczone.length === 0
      ? `Zlecenie nie dotarło do żadnego z ${zleconych} okien wykonawcy.`
      : `Rdzeń przyjął zlecenie dla ${przebieg.doreczone.length} z ${zleconych} okien wykonawcy: ${przebieg.doreczone.join(', ')}.`,
  );
  if (przebieg.niedoreczone.length > 0) {
    czesci.push(`Bez doręczenia (${przebieg.niedoreczone.length}): ${przebieg.niedoreczone.join(' ')}`);
  }
  if (przebieg.nieutrwalone.length > 0) {
    czesci.push(
      `Doręczone, ale NIEUTRWALONE w rdzeniu (${przebieg.nieutrwalone.length}) — po restarcie więzi nie będzie: ${przebieg.nieutrwalone.join(' ')}`,
    );
  }
  return czesci.join(' ');
}

/**
 * Znakuje czynność, której przesłanka jeszcze nie zaszła, bez wygaszania.
 *
 * Przycisk zostaje klikalny, bo jego własna odmowa niesie powód pełniejszy niż
 * jakikolwiek dymek. Znacznik pokazuje przesłankę przed naciśnięciem: `title`
 * dla wskaźnika, `aria-description` dla czytnika ekranu, `data-przeslanka` dla
 * odczytu automatycznego.
 */
function oznaczCzynnosc(kontrolka: HTMLButtonElement, brakujacaPrzeslanka: string): void {
  if (brakujacaPrzeslanka === '') {
    kontrolka.removeAttribute('title');
    kontrolka.removeAttribute('aria-description');
    delete kontrolka.dataset['przeslanka'];
    return;
  }
  const zdanie = `Czynność nie ruszy, dopóki ${brakujacaPrzeslanka} — naciśnięcie powie to samo wprost.`;
  kontrolka.title = zdanie;
  kontrolka.setAttribute('aria-description', zdanie);
  kontrolka.dataset['przeslanka'] = brakujacaPrzeslanka;
}

/** Treść odmowy wraz z kodem kontraktu — po nim odróżnia się rodzaj odmowy. */
function powod(blad?: ErrorInfo): string {
  if (blad === undefined) return 'Rdzeń nie podał przyczyny.';
  return `Powód: ${blad.message} (kod ${blad.code}).`;
}
