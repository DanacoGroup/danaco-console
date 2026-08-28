import { ConfigScope, type MonitorStatus } from '../../../../shared/contract';
import {
  poleWyboru,
  pozycjaWykazu,
  przyciskAkcji as przycisk,
  wykaz,
} from '../../modele/kontrolki-formularza';
import type { ZrodloSekcjiPaneli } from '../../powloka/zrodlo-sekcji-paneli';
import {
  KLUCZ_TRYBU_AOD,
  KSZTALT_TRYBU_AOD,
  NAZWY_TRYBOW_AOD,
  POZIOMY_DECYZJI,
  TrybNakladki,
  trybNakladkiZWartosci,
  WCIELENIA_WALIDATORA,
} from './hierarchia-decyzji';
import { PRZEDROSTEK, wierszOpisu } from './kontrolki';
import { utworzPowierzchnieSekcji, zdaniePuste, zObjasnieniem } from './powierzchnia-sekcji';
import { utworzWiezAutomations } from './wiez-automations';
import type { ZrodloBiegu } from './zrodlo-biegu';
import type { ZrodloNadzoru } from './zrodlo-nadzoru';
import type { ZrodloOkien } from './zrodlo-okien';

/**
 * Sekcja Monitor procesu pokazuje statusy przebiegów na żywo, hierarchię
 * decyzji i nadzór nakładki Always On Display, której tryb utrwala się
 * ustawieniem karty sesji, bo kontrakt nie niesie pola trybu.
 */
export interface SekcjaMonitora {
  element: HTMLElement;
  odswiez(): void;
  ustawOkno(idOkna: string): void;
  rozlacz(): void;
}

export interface OpcjeSekcjiMonitora {
  bieg: ZrodloBiegu;
  nadzor: ZrodloNadzoru;
  okna: ZrodloOkien;
  sekcje: ZrodloSekcjiPaneli;
  /** Karta sesji — adres telemetrii i miejsce utrwalenia trybu nakładki. */
  sesja(): string;
}

export function utworzSekcjeMonitora(opcje: OpcjeSekcjiMonitora): SekcjaMonitora {
  const { bieg, nadzor, okna, sekcje, sesja } = opcje;
  /** Tryb nakładki widziany przez ten widok; cofany po odmowie zapisu. */
  let tryb: TrybNakladki = TrybNakladki.Obserwator;

  const wyborTrybu = poleWyboru(
    { etykieta: 'Tryb Always On Display' },
    NAZWY_TRYBOW_AOD.map(([wartosc, etykieta]) => ({ wartosc, etykieta })),
  );
  const stanNakladki = document.createElement('div');
  stanNakladki.className = 'dm-orkiestracja__ocena';

  const nadzorAod = document.createElement('div');
  nadzorAod.className = 'dm-orkiestracja__pasek';
  nadzorAod.append(
    zObjasnieniem(
      wyborTrybu.element,
      'Rozstrzyga, co nakładka Always On Display może zrobić z procesem. Obserwator wyłącznie patrzy, podpowiada i ostrzega. Operator dodatkowo zatwierdza i wstrzymuje kroki z dowolnego miejsca — to jest druga pozycja hierarchii decyzji, zaraz po Operatorze przy biurku.',
    ),
    stanNakladki,
  );

  const listaProcesow = wykaz('Przebiegi tej karty sesji', 'dm-wykaz');
  const listaPoziomow = wykaz('Hierarchia decyzji', 'dm-wykaz');
  const listaWcielen = wykaz('Wcielenia Executora 3', 'dm-wykaz');

  const wiez = utworzWiezAutomations({
    okno: 'Execution Monitor (oraz Results Analyzer)',
    rola: 'podgląd przebiegów pracy ciągłej i zestawienie ich wyników.',
  });

  const powierzchnia = utworzPowierzchnieSekcji({
    klucz: 'monitor',
    tytul: 'Monitor procesu',
    zakres:
      'Statusy przebiegów na żywo, hierarchia decyzji i nadzór Always On Display — w tym przełącznik trybu obserwator/operator.',
    sekcje,
    podsekcje: [
      { id: 'przebiegi', tytul: 'Przebiegi tej karty sesji', tresc: listaProcesow },
      { id: 'aod', tytul: 'Nadzór Always On Display', tresc: nadzorAod },
      { id: 'hierarchia', tytul: 'Hierarchia decyzji', tresc: listaPoziomow },
      { id: 'wcielenia', tytul: 'Wcielenia Executora 3', tresc: listaWcielen },
      { id: 'wiez', tytul: 'Powiązanie z modułem Automations', tresc: wiez.element },
    ],
  });

  listaPoziomow.replaceChildren(
    ...POZIOMY_DECYZJI.map(
      (poziom) => pozycjaWykazu(`${poziom.numer}. ${poziom.nazwa}`, poziom.zakres, PRZEDROSTEK).element,
    ),
    zdaniePuste(
      'Hierarchia jest konfigurowalna i możliwa do pominięcia: Operator interweniuje bezpośrednio na dowolnym poziomie, w dowolnej chwili. Ten wykaz opisuje kolejność rozstrzygania, nie odbiera nikomu żadnej kontrolki.',
    ),
  );

  listaWcielen.replaceChildren(
    ...WCIELENIA_WALIDATORA.map(
      (miano) =>
        pozycjaWykazu(miano, 'konfiguracja roli okna Results Analyzer, nie osobny byt', PRZEDROSTEK)
          .element,
    ),
    zdaniePuste(
      'Wcielenie nadaje się w sekcji Role komendą `role.update` na oknie Results Analyzer. Siedem wcieleń to siedem ustawień JEDNEJ roli — rdzeń nie zna siedmiu ról.',
    ),
  );

  wyborTrybu.kontrolka.addEventListener('change', () => {
    void przestawTryb(trybNakladkiZWartosci(wyborTrybu.kontrolka.value));
  });

  // Przestawienie trybu wraz z utrwaleniem na poziomie sesji; nieudany zapis cofa przełącznik.
  async function przestawTryb(nowy: TrybNakladki): Promise<void> {
    const karta = sesja();
    const poprzedni = tryb;
    tryb = nowy;
    if (karta === '') {
      powierzchnia.meldunek(
        `Tryb „${nowy}" trzyma WYŁĄCZNIE ten widok — sesja nie jest uzgodniona z rdzeniem, więc nie ma na czym go utrwalić.`,
        false,
      );
      return;
    }
    const wynik = await okna.zapiszUstawienie({
      key: KLUCZ_TRYBU_AOD,
      value: nowy,
      scope: ConfigScope.Session,
      scopeId: karta,
    });
    if (!wynik.udany || wynik.wynik === undefined) {
      tryb = poprzedni;
      wyborTrybu.kontrolka.value = poprzedni;
      powierzchnia.meldunek(
        `Rdzeń odmówił zapisu trybu nakładki, przełącznik wraca do „${poprzedni}". Powód: ${wynik.blad?.message ?? 'rdzeń nie podał przyczyny'} (kod ${wynik.blad?.code ?? 'brak'}). Kontrakt nie ma pola trybu ani komendy: ${KSZTALT_TRYBU_AOD}`,
        false,
      );
      return;
    }
    const zapisany = trybNakladkiZWartosci(wynik.wynik.value);
    tryb = zapisany;
    wyborTrybu.kontrolka.value = zapisany;
    powierzchnia.meldunek(
      `Tryb nakładki „${zapisany}" utrwalony na poziomie karty sesji. UWAGA: rdzeń trzyma go jako ustawienie, a nie jako stan nakładki — ${KSZTALT_TRYBU_AOD}`,
      true,
    );
  }

  /** Odczyt trybu utrwalonego; brak wpisu znaczy obserwatora, nie usterkę. */
  async function odczytajTryb(): Promise<void> {
    const karta = sesja();
    if (karta === '') return;
    const wpis = await okna.ustawienia({
      key: KLUCZ_TRYBU_AOD,
      scope: ConfigScope.Session,
      scopeId: karta,
    });
    tryb = trybNakladkiZWartosci(wpis.udany ? wpis.wynik?.[0]?.value : undefined);
    wyborTrybu.kontrolka.value = tryb;
  }

  /** Stan nakładki wprost z rdzenia; pola trybu kontrakt w nim nie niesie. */
  async function odczytajNakladke(): Promise<void> {
    const wynik = await nadzor.stanNakladki({});
    if (!wynik.udany || wynik.wynik === undefined) {
      stanNakladki.replaceChildren(
        zdaniePuste(
          `Rdzeń nie oddał stanu nakładki: ${wynik.blad?.message ?? 'bez przyczyny'} (kod ${wynik.blad?.code ?? 'brak'}). Przełącznik trybu działa dalej — tryb utrwala się ustawieniem sesji.`,
        ),
      );
      return;
    }
    const stan = wynik.wynik;
    stanNakladki.replaceChildren(
      wierszOpisu('Procesy w biegu', String(stan.runningProcessCount)),
      wierszOpisu('Procesy przypięte do nakładki', String(stan.attachedProcessIds?.length ?? 0)),
      wierszOpisu('Okno pokazywane w nakładce', stan.activeWindowId ?? 'rdzeń nie wskazał'),
      zdaniePuste(
        'Kontrakt nie niesie trybu nakładki: `AodStatus` mówi, co nakładka pokazuje, ale nie mówi, co jej wolno. Tryb widoczny wyżej jest ustawieniem karty sesji.',
      ),
    );
  }

  /** Telemetria: subskrypcja i stan bieżący jednym żądaniem. */
  function odczytaj(): void {
    powierzchnia.tresci.ladowanie('Odczyt przebiegów i stanu nakładki…');
    void (async () => {
      const automatyki = await nadzor.automatyki({});
      wiez.ustawWykaz(automatyki.udany ? (automatyki.wynik ?? []) : []);
      await odczytajTryb();
      await odczytajNakladke();

      const karta = sesja();
      const wynik = await bieg.subskrybujMonitor(karta === '' ? {} : { sessionId: karta });
      if (!wynik.udany || wynik.wynik === undefined) {
        powierzchnia.tresci.blad('Rdzeń odmówił obserwacji przebiegów.', wynik.blad);
        return;
      }
      powierzchnia.tresci.pusto('');
      if (!wynik.wynik.subscribed) {
        powierzchnia.meldunek(
          'Rdzeń oddał stan przebiegów, ale NIE założył obserwacji — wykaz niżej jest migawką, a nie podglądem na żywo.',
          false,
        );
      }
      pokazPrzebiegi(wynik.wynik.statuses);
    })();
  }

  /** Wykaz przebiegów; pustka jest tu stanem poprawnym. */
  function pokazPrzebiegi(statusy: readonly MonitorStatus[]): void {
    if (statusy.length === 0) {
      listaProcesow.replaceChildren(
        zdaniePuste(
          'Ta karta sesji nie ma w tej chwili ani jednego przebiegu. To nie jest brak monitora — to znaczy, że zespół nie pracuje.',
        ),
      );
      return;
    }
    listaProcesow.replaceChildren(...statusy.map(pozycjaPrzebiegu));
  }

  /** Jeden przebieg wraz z czynnością przypięcia go do nakładki. */
  function pozycjaPrzebiegu(status: MonitorStatus): HTMLElement {
    const { element, akcje } = pozycjaWykazu(
      status.label ?? status.processId,
      `${status.status} · etap ${status.stageIndex ?? 0} z ${status.stageCount ?? 0} · ${status.completion ?? 0}%`,
      PRZEDROSTEK,
    );
    const przypnij = przycisk('Przypnij do nakładki', 'dn-btn dn-btn--sm dn-btn--zarys');
    przypnij.title =
      'Dokłada ten przebieg do nakładki Always On Display: jego postęp będzie widoczny poza tym oknem, a w trybie operatora da się go stamtąd wstrzymać.';
    przypnij.setAttribute('aria-description', przypnij.title);
    przypnij.addEventListener('click', () => {
      void przypnijProces(status);
    });
    const odepnij = przycisk('Odepnij', 'dn-btn dn-btn--sm dn-btn--zarys');
    odepnij.title = 'Zdejmuje przebieg z nakładki. Sam przebieg biegnie dalej — zmienia się tylko to, co widać poza oknem.';
    odepnij.setAttribute('aria-description', odepnij.title);
    odepnij.addEventListener('click', () => {
      void odepnijProces(status);
    });
    akcje.append(przypnij, odepnij);
    element.append(wierszOpisu('Obieg naprawczy', String(status.cycle ?? 0)));
    return element;
  }

  /** Przypięcie przebiegu do nakładki; liczba przypiętych z odpowiedzi. */
  async function przypnijProces(status: MonitorStatus): Promise<void> {
    const wynik = await nadzor.przypnijDoNakladki({ processId: status.processId });
    if (!wynik.udany || wynik.wynik === undefined) {
      powierzchnia.meldunek(
        `Rdzeń odmówił przypięcia przebiegu ${status.processId} do nakładki. Powód: ${wynik.blad?.message ?? 'bez przyczyny'} (kod ${wynik.blad?.code ?? 'brak'}).`,
        false,
      );
      return;
    }
    powierzchnia.meldunek(
      `Przebieg ${status.processId} przypięty; nakładka obserwuje ${wynik.wynik.attachedProcessIds.length} przebiegów.`,
      true,
    );
    void odczytajNakladke();
  }

  /** Odpięcie przebiegu od nakładki. */
  async function odepnijProces(status: MonitorStatus): Promise<void> {
    const wynik = await nadzor.odepnijOdNakladki({ processId: status.processId });
    if (!wynik.udany || wynik.wynik === undefined) {
      powierzchnia.meldunek(
        `Rdzeń odmówił odpięcia przebiegu ${status.processId}. Powód: ${wynik.blad?.message ?? 'bez przyczyny'} (kod ${wynik.blad?.code ?? 'brak'}).`,
        false,
      );
      return;
    }
    powierzchnia.meldunek(
      `Przebieg ${status.processId} odpięty; nakładka obserwuje ${wynik.wynik.attachedProcessIds.length} przebiegów.`,
      true,
    );
    void odczytajNakladke();
  }

  wiez.naZmiane(() => odczytaj());

  // Postęp na żywo odczytuje stan monitora, nie przebudowuje wiersza z samego zdarzenia.
  const odsubskrybuj = nadzor.naPostep(() => {
    const karta = sesja();
    void bieg.stanMonitora(karta === '' ? {} : { sessionId: karta }).then((wynik) => {
      if (wynik.udany && wynik.wynik !== undefined) pokazPrzebiegi(wynik.wynik.statuses);
    });
  });

  return {
    element: powierzchnia.element,
    odswiez: odczytaj,
    ustawOkno: powierzchnia.ustawOkno,
    rozlacz: odsubskrybuj,
  };
}
