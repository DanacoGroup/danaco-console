import {
  AutomationTriggerKind,
  Command,
  type AutomationExecution,
  type AutomationSchedule,
  type AutomationTrigger,
} from '../../../../shared/contract';
import { utworzRameOkna } from '../../komponenty/rama-okna';
import {
  pole,
  poleLiczbowe,
  przelacznik,
  przyciskAkcji as przycisk,
  pozycjaWykazu,
  wiersz,
  wybor,
  wykaz,
} from '../../modele/kontrolki-formularza';
import type { PokrycieKomend } from '../pokrycie-komend';
import {
  dniKalendarza,
  podpisKalendarza,
  siatkaKalendarza,
  ZAKRESY_KALENDARZA,
  NAZWY_ZAKRESOW,
  type ZakresKalendarza,
} from './kalendarz-uruchomien';
import { nastepneUruchomienia } from './nastepne-uruchomienia';
import type { StanAutomatyki } from './stan-automatyki';
import { utworzStanTresci } from './stany-okna';
import {
  DNI_TYGODNIA,
  NASTAWY_WYJSCIOWE,
  NAZWY_WZORCOW,
  opisCyklicznosci,
  rozpoznajWzorzec,
  WZORCE_CYKLICZNOSCI,
  zapisWzorca,
  type NastawyWzorca,
  type WzorzecCyklicznosci,
} from './wzorce-cyklicznosci';
import type { ZrodloAutomations } from './zrodlo-automations';

/**
 * Scheduler ustala cykliczność i wyzwalacze automatyki, pozwala ją wstrzymać
 * oraz wznowić, a także pokazuje podgląd kolejnych uruchomień i kalendarz
 * przebiegów. Harmonogram obowiązuje po powiązaniu z automatyką.
 */
export interface OknoSchedulera {
  element: HTMLElement;
  odswiez(): void;
}

/**
 * Nazwy rodzajów wyzwalacza na ekranie.
 *
 * Mapa zupełna po wyliczeniu, nie wykaz przepisany ręcznie: gdy
 * `AutomationTriggerKind` urośnie, kompilacja zatrzyma się tutaj, zamiast
 * zostawić wykaz rozwijany uboższy od kontraktu.
 */
const NAZWY_WYZWALACZY: Readonly<Record<AutomationTriggerKind, string>> = {
  [AutomationTriggerKind.Cron]: 'cykliczność (cron)',
  [AutomationTriggerKind.Webhook]: 'webhook',
  [AutomationTriggerKind.File]: 'zmiana pliku',
  [AutomationTriggerKind.Condition]: 'warunek na wyniku',
  [AutomationTriggerKind.ModuleEvent]: 'zdarzenie modułu',
  [AutomationTriggerKind.Chain]: 'zakończenie innej automatyki',
};

/**
 * Rodzaje wyzwalacza w kolejności kontraktu, złożone z wyliczenia
 * AutomationTriggerKind i mapy nazw ekranowych, gotowe do wykazu wyboru w oknie.
 */
export const RODZAJE_WYZWALACZA: ReadonlyArray<[string, string]> = Object.values(
  AutomationTriggerKind,
).map((rodzaj) => [rodzaj, NAZWY_WYZWALACZY[rodzaj]]);

export function utworzOknoSchedulera(
  zrodlo: ZrodloAutomations,
  stan: StanAutomatyki,
  pokrycie: PokrycieKomend,
): OknoSchedulera {
  const rama = utworzRameOkna({
    tytul: 'Scheduler',
    rola: 'zarządca',
    przeznaczenie:
      'Cykliczność i wyzwalacze automatyki. Harmonogram obowiązuje po powiązaniu z automatyką.',
    przedrostek: 'da',
  });
  const tresc = utworzStanTresci();
  const powierzchnia = zlozPowierzchnieHarmonogramu(rama, tresc.element, pokrycie);
  const { cron, strefa, obowiazuje, rodzajWyzwalacza, wyrazenie } = powierzchnia;
  const wyzwalacze: AutomationTrigger[] = [];

  /** Nastawy kreatora; zmieniają się razem z polami wzorca i z zapisem cron. */
  let nastawy: NastawyWzorca = { ...NASTAWY_WYJSCIOWE };

  // Przepisuje nastawy z pól kreatora i składa z nich zapis cron; wzorzec własny zapis zostawia.
  function zlozZapis(): void {
    nastawy = {
      minuta: liczbaZPola(powierzchnia.minuta, nastawy.minuta),
      godzina: liczbaZPola(powierzchnia.godzina, nastawy.godzina),
      dzienTygodnia: Number.parseInt(powierzchnia.dzienTygodnia.value, 10),
      dzienMiesiaca: liczbaZPola(powierzchnia.dzienMiesiaca, nastawy.dzienMiesiaca),
      krok: liczbaZPola(powierzchnia.krok, nastawy.krok),
    };
    const wybrany = powierzchnia.wzorzec.value as WzorzecCyklicznosci;
    if (wybrany !== WZORCE_CYKLICZNOSCI.wlasny) {
      cron.value = zapisWzorca(wybrany, nastawy, cron.value);
    }
    opiszCyklicznosc();
  }

  // Rozpoznaje wzorzec w zapisie wpisanym wprost i przestawia na niego pola kreatora.
  function rozpoznajZapis(): void {
    const rozpoznany = rozpoznajWzorzec(cron.value);
    nastawy = rozpoznany.nastawy;
    powierzchnia.wzorzec.value = rozpoznany.wzorzec;
    powierzchnia.minuta.value = String(nastawy.minuta);
    powierzchnia.godzina.value = String(nastawy.godzina);
    powierzchnia.dzienTygodnia.value = String(nastawy.dzienTygodnia);
    powierzchnia.dzienMiesiaca.value = String(nastawy.dzienMiesiaca);
    powierzchnia.krok.value = String(nastawy.krok);
    opiszCyklicznosc();
  }

  /** Zdanie o cykliczności pod polem zapisu — mowa Operatora obok składni. */
  function opiszCyklicznosc(): void {
    powierzchnia.opisZapisu.textContent = opisCyklicznosci(cron.value);
  }

  function pokazWyzwalacze(harmonogram?: AutomationSchedule): void {
    const miejsce = tresc.tresc();
    const pozycje = harmonogram?.triggers ?? wyzwalacze;
    if (pozycje.length === 0) miejsce.append(pustkaWyzwalaczy());
    miejsce.append(
      listaWyzwalaczy(pozycje, (wyzwalacz) => {
        const miejsceWykazu = wyzwalacze.indexOf(wyzwalacz);
        if (miejsceWykazu >= 0) wyzwalacze.splice(miejsceWykazu, 1);
        pokazWyzwalacze();
      }),
    );
    if (harmonogram !== undefined) miejsce.append(opisHarmonogramu(harmonogram));
  }

  function zapiszHarmonogram(czynny: boolean, czynnosc: string): void {
    const automatyka = stan.automatyka();
    if (automatyka === '') {
      tresc.pusto('Wskaż automatykę w Workflow Builderze — harmonogram obowiązuje po powiązaniu z nią.');
      return;
    }
    tresc.ladowanie('Zapis harmonogramu…');
    const zadanie = zadanieHarmonogramu(automatyka, czynny, wyzwalacze, cron.value, strefa.value);
    void zrodlo.ustawHarmonogram(zadanie).then((wynik) => {
      if (!wynik.udany || wynik.wynik === undefined) {
        tresc.blad('Rdzeń nie zapisał harmonogramu.', wynik.blad);
        return;
      }
      const harmonogram = wynik.wynik.schedule;
      obowiazuje.checked = harmonogram.enabled;
      wyzwalacze.splice(0, wyzwalacze.length, ...(harmonogram.triggers ?? []));
      pokazWyzwalacze(harmonogram);
      tresc.potwierdzenie(
        `${czynnosc} ${zdanieOStanie(harmonogram, czynny)}`,
        harmonogram.enabled === czynny,
      );
    });
  }

  powierzchnia.dodajWyzwalacz.addEventListener('click', () => {
    const wyrazenieWyzwalacza = wyrazenie.value.trim();
    if (wyrazenieWyzwalacza === '') {
      tresc.potwierdzenie('Wyzwalacz bez treści nie ma czego obserwować — nie został dodany.', false);
      return;
    }
    wyzwalacze.push({
      kind: rodzajWyzwalacza.value as AutomationTrigger['kind'],
      expression: wyrazenieWyzwalacza,
      enabled: true,
    });
    wyrazenie.value = '';
    pokazWyzwalacze();
  });

  powierzchnia.zapisz.addEventListener('click', () =>
    zapiszHarmonogram(obowiazuje.checked, 'Harmonogram zapisany.'));
  powierzchnia.wstrzymaj.addEventListener('click', () =>
    zapiszHarmonogram(false, 'Wysłano wstrzymanie harmonogramu.'));
  powierzchnia.wznow.addEventListener('click', () =>
    zapiszHarmonogram(true, 'Wysłano wznowienie harmonogramu.'));

  // Odczyt harmonogramów automatyki komendą schedule.get; wykaz pochodzi z odpowiedzi rdzenia.
  powierzchnia.odczytaj.addEventListener('click', () => {
    const automatyka = stan.automatyka();
    if (automatyka === '') {
      tresc.pusto('Wskaż automatykę w Workflow Builderze — odczyt dotyczy jej harmonogramów.');
      return;
    }
    tresc.ladowanie('Odczyt harmonogramów…');
    void zrodlo.odczytajHarmonogramy({ workflowId: automatyka }).then((wynik) => {
      if (!wynik.udany || wynik.wynik === undefined) {
        tresc.blad('Rdzeń nie oddał harmonogramów.', wynik.blad);
        return;
      }
      const harmonogramy = wynik.wynik;
      if (harmonogramy.length === 0) {
        tresc.potwierdzenie(
          `Rdzeń nie ma jeszcze harmonogramu dla automatyki ${automatyka} — zapisz go, aby powstał.`,
          true,
        );
        return;
      }
      const pierwszy = harmonogramy[0];
      if (pierwszy !== undefined) {
        obowiazuje.checked = pierwszy.enabled;
        wyzwalacze.splice(0, wyzwalacze.length, ...(pierwszy.triggers ?? []));
        // Kreator pokazuje cykliczność zastaną, nie wpisaną przed odczytem, po rozpoznaniu wzorca.
        cron.value = pierwszy.cron ?? '';
        if (pierwszy.timeZone !== undefined) strefa.value = pierwszy.timeZone;
        rozpoznajZapis();
      }
      const miejsce = tresc.tresc();
      for (const harmonogram of harmonogramy) miejsce.append(opisHarmonogramu(harmonogram));
      tresc.potwierdzenie(`Rdzeń oddał ${harmonogramy.length} harmonogramów tej automatyki.`, true);
    });
  });

  powierzchnia.podglad.addEventListener('click', () => {
    const terminy = nastepneUruchomienia(cron.value);
    if (terminy.length === 0) {
      tresc.blad('Nie umiem odczytać tej cykliczności — podaj pięć pól zapisu cron.');
      return;
    }
    tresc.tresc().append(wykazUruchomien(terminy), zastrzezeniePodgladu());
  });

  // Kalendarz uruchomień zestawia terminy zaplanowane z przebiegami, które już się odbyły.
  powierzchnia.kalendarz.addEventListener('click', () => {
    const terminy = nastepneUruchomienia(cron.value, new Date(), 1);
    if (terminy.length === 0 && cron.value.trim() !== '') {
      tresc.blad('Nie umiem odczytać tej cykliczności — podaj pięć pól zapisu cron.');
      return;
    }
    const automatyka = stan.automatyka();
    const zakres = powierzchnia.zakresKalendarza.value as ZakresKalendarza;
    if (automatyka === '') {
      pokazKalendarz([], zakres, 'Automatyka niewskazana, więc kalendarz pokazuje same terminy zaplanowane.');
      return;
    }
    tresc.ladowanie('Odczyt przebiegów do kalendarza…');
    void zrodlo.przebiegi({ workflowId: automatyka }).then((wynik) => {
      if (!wynik.udany || wynik.wynik === undefined) {
        tresc.blad('Rdzeń nie oddał przebiegów, więc kalendarz nie ma historii do pokazania.', wynik.blad);
        return;
      }
      pokazKalendarz(
        wynik.wynik.executions,
        zakres,
        `Historia wzięta z ${wynik.wynik.executions.length} przebiegów oddanych przez rdzeń.`,
      );
    });
  });

  /** Rysuje kalendarz wraz z podpisem o pochodzeniu obu warstw danych. */
  function pokazKalendarz(
    przebiegi: readonly AutomationExecution[],
    zakres: ZakresKalendarza,
    zdanieOHistorii: string,
  ): void {
    const dni = dniKalendarza(cron.value, przebiegi, zakres);
    const miejsce = tresc.tresc();
    miejsce.append(siatkaKalendarza(dni), akapitOpisowy(podpisKalendarza(dni)));
    tresc.potwierdzenie(zdanieOHistorii, true);
  }

  powierzchnia.wzorzec.addEventListener('change', zlozZapis);
  for (const kontrolka of [
    powierzchnia.minuta,
    powierzchnia.godzina,
    powierzchnia.dzienMiesiaca,
    powierzchnia.krok,
  ]) {
    kontrolka.addEventListener('input', zlozZapis);
  }
  powierzchnia.dzienTygodnia.addEventListener('change', zlozZapis);
  cron.addEventListener('input', rozpoznajZapis);
  powierzchnia.strefaTegoKomputera.addEventListener('click', () => {
    strefa.value = Intl.DateTimeFormat().resolvedOptions().timeZone;
    tresc.potwierdzenie(`Strefa ustawiona na strefę tego komputera: ${strefa.value}.`, true);
  });

  rozpoznajZapis();

  function odswiez(): void {
    if (stan.automatyka() === '') {
      tresc.pusto('Wskaż automatykę w Workflow Builderze, aby nadać jej cykliczność.');
      return;
    }
    pokazWyzwalacze();
  }

  return { element: rama.element, odswiez };
}

/**
 * Kontrolki okna Schedulera: pola kreatora cykliczności, zapis cron, strefa
 * czasowa, wyzwalacze oraz przyciski akcji panelu.
 */
interface PowierzchniaHarmonogramu {
  wzorzec: HTMLSelectElement;
  minuta: HTMLInputElement;
  godzina: HTMLInputElement;
  dzienTygodnia: HTMLSelectElement;
  dzienMiesiaca: HTMLInputElement;
  krok: HTMLInputElement;
  opisZapisu: HTMLElement;
  cron: HTMLInputElement;
  strefa: HTMLInputElement;
  strefaTegoKomputera: HTMLButtonElement;
  obowiazuje: HTMLInputElement;
  rodzajWyzwalacza: HTMLSelectElement;
  wyrazenie: HTMLInputElement;
  zakresKalendarza: HTMLSelectElement;
  dodajWyzwalacz: HTMLButtonElement;
  zapisz: HTMLButtonElement;
  wstrzymaj: HTMLButtonElement;
  wznow: HTMLButtonElement;
  podglad: HTMLButtonElement;
  kalendarz: HTMLButtonElement;
  odczytaj: HTMLButtonElement;
}

/**
 * Składa kontrolki, pasek akcji i ciało okna.
 *
 * Czysta konstrukcja: fragment nie domyka się na stanie okna ani na źródle.
 * Kolejność dokładania jest znacząca — po niej idą sprawdziany widoku.
 */
function zlozPowierzchnieHarmonogramu(
  rama: { akcje: HTMLElement; narzedzia: HTMLElement; cialo: HTMLElement },
  stanTresci: HTMLElement,
  pokrycie: PokrycieKomend,
): PowierzchniaHarmonogramu {
  const wzorzec = wybor('Wzorzec cykliczności', NAZWY_WZORCOW);
  const minuta = poleLiczbowe('Minuta', '0–59');
  minuta.value = String(NASTAWY_WYJSCIOWE.minuta);
  const godzina = poleLiczbowe('Godzina', '0–23');
  godzina.value = String(NASTAWY_WYJSCIOWE.godzina);
  const dzienTygodnia = wybor('Dzień tygodnia', DNI_TYGODNIA);
  dzienTygodnia.value = String(NASTAWY_WYJSCIOWE.dzienTygodnia);
  const dzienMiesiaca = poleLiczbowe('Dzień miesiąca', '1–31');
  dzienMiesiaca.value = String(NASTAWY_WYJSCIOWE.dzienMiesiaca);
  const krok = poleLiczbowe('Krok interwału', 'co ile minut albo godzin');
  krok.value = String(NASTAWY_WYJSCIOWE.krok);

  const cron = pole('Cykliczność w notacji cron', 'np. 0 6 * * 1-5');
  const opisZapisu = document.createElement('p');
  opisZapisu.className = 'dn-pole-opis da-opis-cyklicznosci';
  const strefa = pole('Strefa czasowa', 'np. Europe/Warsaw');
  const strefaTegoKomputera = przycisk('Strefa tego komputera', 'dn-btn dn-btn--sm dn-btn--zarys');
  const obowiazuje = przelacznik('Harmonogram obowiązuje');
  obowiazuje.checked = true;

  const rodzajWyzwalacza = wybor('Rodzaj wyzwalacza', RODZAJE_WYZWALACZA);
  const wyrazenie = pole('Treść wyzwalacza', 'adres, ścieżka albo warunek');
  const zakresKalendarza = wybor('Zakres kalendarza', NAZWY_ZAKRESOW);
  zakresKalendarza.value = ZAKRESY_KALENDARZA.miesiac;

  const dodajWyzwalacz = przycisk('+ Dodaj wyzwalacz');
  const zapisz = przycisk('Zapisz harmonogram', 'dn-btn dn-btn--atrament');
  const wstrzymaj = przycisk('Wstrzymaj harmonogram');
  const wznow = przycisk('Wznów harmonogram');
  const podglad = przycisk('Podgląd kolejnych uruchomień');
  const kalendarz = przycisk('Kalendarz uruchomień');
  // Odczytaj harmonogram woła komendę schedule.get rdzenia i oddaje pole schedules.
  const odczytaj = przycisk('Odczytaj harmonogram');

  rama.akcje.append(
    dodajWyzwalacz,
    zapisz,
    wstrzymaj,
    wznow,
    podglad,
    kalendarz,
    odczytaj,
    pokrycie.przycisk(
      'Okna wykonania',
      Command.ScheduleWindowSet,
      'ograniczenie przedziałów czasu, w których uruchomienie następuje',
    ),
    pokrycie.przycisk(
      'Uruchomienie wsteczne',
      Command.ScheduleBackfillRun,
      'wykonanie przebiegów dla przeszłych, pominiętych terminów w zadanym zakresie dat',
    ),
    pokrycie.przycisk(
      'Historia wyzwoleń',
      Command.ScheduleTriggerHistory,
      'zapis rzeczywistych momentów wyzwolenia wraz z przyczyną',
    ),
    pokrycie.przycisk(
      'Nadzór obecności uruchomień',
      Command.ScheduleHeartbeatSet,
      'alarm, gdy oczekiwane uruchomienie nie nastąpiło w oknie tolerancji',
    ),
  );

  rama.narzedzia.append(zakresKalendarza, strefaTegoKomputera);

  rama.cialo.append(
    wiersz('Wzorzec', wzorzec, {
      klasa: 'da-wiersz',
      objasnienie:
        'Wzorzec składa zapis cron z pól poniżej. Zapis własny zostawia pole cykliczności Operatorowi.',
    }),
    kreatorCyklicznosci(minuta, godzina, dzienTygodnia, dzienMiesiaca, krok),
    wiersz('Cykliczność (cron)', cron, {
      klasa: 'da-wiersz',
      objasnienie: 'Pięć pól: minuta, godzina, dzień miesiąca, miesiąc, dzień tygodnia.',
    }),
    opisZapisu,
    wiersz('Strefa czasowa', strefa, {
      klasa: 'da-wiersz',
      objasnienie: 'Zapisywana i oddawana kontraktem; rachunek terminu idzie w UTC.',
    }),
    wiersz('Obowiązuje', obowiazuje, {
      klasa: 'da-wiersz',
      objasnienie: 'Harmonogram wyłączony nie ma najbliższego uruchomienia.',
    }),
    wiersz('Rodzaj wyzwalacza', rodzajWyzwalacza, {
      klasa: 'da-wiersz',
      objasnienie:
        'Wykaz bierze się z wyliczenia kontraktu, nie z zapisu w oknie — rośnie razem z nim.',
    }),
    wiersz('Treść wyzwalacza', wyrazenie, {
      klasa: 'da-wiersz',
      objasnienie: 'Webhook — adres; plik — ścieżka; warunek — wyrażenie.',
    }),
    stanTresci,
  );

  return {
    wzorzec, minuta, godzina, dzienTygodnia, dzienMiesiaca, krok, opisZapisu,
    cron, strefa, strefaTegoKomputera, obowiazuje, rodzajWyzwalacza, wyrazenie, zakresKalendarza,
    dodajWyzwalacz, zapisz, wstrzymaj, wznow, podglad, kalendarz, odczytaj,
  };
}

/**
 * Pola kreatora cykliczności w jednym rzędzie, widoczne zawsze niezależnie od
 * wybranego wzorca: pole nieużywane przez wzorzec jest nieszkodliwe, bo jego
 * wartość nie wchodzi do zapisu.
 */
function kreatorCyklicznosci(
  minuta: HTMLInputElement,
  godzina: HTMLInputElement,
  dzienTygodnia: HTMLSelectElement,
  dzienMiesiaca: HTMLInputElement,
  krok: HTMLInputElement,
): HTMLElement {
  const rzad = document.createElement('div');
  rzad.className = 'da-kreator';
  rzad.append(
    wiersz('Minuta', minuta, { klasa: 'da-wiersz' }),
    wiersz('Godzina', godzina, { klasa: 'da-wiersz' }),
    wiersz('Dzień tygodnia', dzienTygodnia, { klasa: 'da-wiersz' }),
    wiersz('Dzień miesiąca', dzienMiesiaca, { klasa: 'da-wiersz' }),
    wiersz('Krok interwału', krok, { klasa: 'da-wiersz' }),
  );
  return rzad;
}

/**
 * Liczba odczytana z pola kreatora cykliczności; pole puste albo o treści
 * nieliczbowej zostawia wartość zastaną sprzed odczytu.
 */
function liczbaZPola(kontrolka: HTMLInputElement, zastana: number): number {
  const wartosc = Number.parseInt(kontrolka.value, 10);
  return Number.isInteger(wartosc) ? wartosc : zastana;
}

/**
 * Akapit opisowy umieszczany pod treścią okna, wspólny dla podpisu kalendarza
 * uruchomień oraz zastrzeżenia podglądu terminów.
 */
function akapitOpisowy(zdanie: string): HTMLElement {
  const akapit = document.createElement('p');
  akapit.className = 'dn-pole-opis';
  akapit.textContent = zdanie;
  return akapit;
}

/**
 * Stan pusty wykazu wyzwalaczy — treść stała, więc wyszła jako czysta
 * konstrukcja. Harmonogram bez wyzwalaczy jest poprawny, nie wadliwy.
 */
function pustkaWyzwalaczy(): HTMLElement {
  const zdanie = document.createElement('p');
  zdanie.className = 'dn-pole-opis';
  zdanie.textContent = 'Harmonogram nie ma jeszcze wyzwalaczy poza cyklicznością.';
  return zdanie;
}

/**
 * Wykaz wyzwalaczy. Fragment bierze wywołanie zwrotne usunięcia, bo przycisk
 * wiersza musi sięgnąć po wykaz roboczy okna — poza tym nie zna stanu okna.
 */
function listaWyzwalaczy(
  pozycje: readonly AutomationTrigger[],
  naUsuniecie: (wyzwalacz: AutomationTrigger) => void,
): HTMLElement {
  const lista = wykaz('Wyzwalacze harmonogramu', 'da-wykaz');
  for (const wyzwalacz of pozycje) {
    const pozycja = pozycjaWykazu(wyzwalacz.kind, wyzwalacz.expression, 'da');
    const usun = przycisk('Usuń', 'dn-btn dn-btn--zarys');
    usun.addEventListener('click', () => naUsuniecie(wyzwalacz));
    pozycja.akcje.append(usun);
    lista.append(pozycja.element);
  }
  return lista;
}

/**
 * Wykaz kolejnych terminów uruchomienia — czysta konstrukcja złożona z dat
 * wyliczonych przez nastepneUruchomienia.
 */
function wykazUruchomien(terminy: readonly Date[]): HTMLElement {
  const lista = wykaz('Kolejne uruchomienia', 'da-wykaz');
  for (const termin of terminy) {
    lista.append(pozycjaWykazu(termin.toISOString(), termin.toLocaleString('pl-PL'), 'da').element);
  }
  return lista;
}

/**
 * Zastrzeżenie do podglądu — treść stała, więc wyszła jako czysta konstrukcja.
 * Podgląd liczy okno, więc nie wolno go podać jako terminu obowiązującego.
 */
function zastrzezeniePodgladu(): HTMLElement {
  const zastrzezenie = document.createElement('p');
  zastrzezenie.className = 'dn-pole-opis';
  zastrzezenie.textContent =
    'Podgląd wyliczony w oknie z wpisanej cykliczności (UTC). Termin obowiązujący oddaje rdzeń po zapisie.';
  return zastrzezenie;
}

/**
 * Zadanie zapisu harmonogramu — czysta konstrukcja z wartości pól. Pole puste
 * nie trafia do żądania wcale, bo pusty łańcuch nie jest cyklicznością.
 */
function zadanieHarmonogramu(
  automatyka: string,
  czynny: boolean,
  wyzwalacze: readonly AutomationTrigger[],
  cron: string,
  strefa: string,
): Parameters<ZrodloAutomations['ustawHarmonogram']>[0] {
  const zadanie: Parameters<ZrodloAutomations['ustawHarmonogram']>[0] = {
    workflowId: automatyka,
    enabled: czynny,
    triggers: [...wyzwalacze],
  };
  if (cron.trim() !== '') zadanie.cron = cron.trim();
  if (strefa.trim() !== '') zadanie.timeZone = strefa.trim();
  return zadanie;
}

/**
 * Zdanie potwierdzenia harmonogramu opiera się na stanie oddanym przez rdzeń,
 * nie na czynności, o którą prosiło okno, i nazywa wprost rozminięcie się
 * stanu z żądaniem.
 */
function zdanieOStanie(harmonogram: AutomationSchedule, zadany: boolean): string {
  const termin =
    harmonogram.nextRunAt === undefined
      ? 'bez najbliższego uruchomienia'
      : `najbliższe uruchomienie ${new Date(harmonogram.nextRunAt).toLocaleString('pl-PL')}`;
  const stan = harmonogram.enabled ? 'obowiązuje' : 'jest wstrzymany';
  if (harmonogram.enabled !== zadany) {
    return `Rdzeń oddał harmonogram, który ${stan} — inaczej, niż żądało okno (${termin}).`;
  }
  return `Rdzeń oddał harmonogram, który ${stan}, ${termin}.`;
}

/**
 * Zdanie o harmonogramie po zapisie: podaje jego identyfikator, stan
 * obowiązywania oraz termin najbliższego uruchomienia.
 */
function opisHarmonogramu(harmonogram: AutomationSchedule): HTMLElement {
  const zdanie = document.createElement('p');
  zdanie.className = 'dn-pole-opis';
  const termin =
    harmonogram.nextRunAt === undefined
      ? 'brak najbliższego uruchomienia'
      : new Date(harmonogram.nextRunAt).toLocaleString('pl-PL');
  zdanie.textContent =
    `Harmonogram ${harmonogram.id} — ${harmonogram.enabled ? 'obowiązuje' : 'wstrzymany'}; najbliższe uruchomienie: ${termin}.`;
  return zdanie;
}
