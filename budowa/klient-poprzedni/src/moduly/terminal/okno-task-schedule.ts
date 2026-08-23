import {
  AutomationStepKind,
  Command,
  QueueAction,
  QueueStatus,
  TerminalProcessStatus,
  TerminalWatchStatus,
  type AutomationSchedule,
  type AutomationStep,
  type Queue,
  type TerminalProcess,
  type TerminalSession,
  type TerminalWatch,
} from '../../../../shared/contract';
import { utworzRameOkna } from '../../komponenty/rama-okna';
import {
  pobierzPlik,
  pole,
  poleLiczbowe,
  poleTresci,
  pozycjaWykazu,
  przelacznikWidoku,
  przestaw,
  przyciskAkcji,
  wykaz,
} from '../../modele/kontrolki-formularza';
import { oznaczWarstwy, type CzynnoscOkna } from './czynnosci-okna';
import {
  czytajZadania,
  MANIFESTY,
  type RodzajManifestu,
  type ZadanieManifestu,
} from './manifesty-zadan';
import type { PokrycieKomend } from '../pokrycie-komend';
import type { StanTerminala } from './stan-terminala';
import { utworzStanTresci } from './stany-okna';
import { utworzWyborDrzewem, type PozycjaWyboru, type WyborDrzewem } from './wybor-drzewem';
import type { ZrodloTerminala } from './zrodlo-terminala';
import { notaZaleznosci, PROGRAMY_CZYNNOSCI } from './zaleznosci-zewnetrzne';

/**
 * Task & Schedule — okno zarządcy modułu Terminal: zadania projektu, odczyt
 * harmonogramów i kolejek, potok poleceń oraz historia przebiegów tego okna.
 *
 * `terminal.command.exec` uruchamia zadanie i krok potoku, `terminal.output.read`
 * oddaje wynik kroku wraz z kodem wyjścia, `schedule.get` czyta harmonogramy,
 * `automation.workflow.save` i `automation.schedule.set` zakładają plan zadania
 * powłoki, a `queue.list` i `queue.action` pokazują kolejkę silnika pętli
 * obsługującą to okno i posuwają ją sześcioma działaniami słownika kontraktu.
 *
 * ── Plan zadania powłoki: druga POWIERZCHNIA, nie druga rodzina komend ──────
 * Opracowanie modułu opisuje okno Task & Schedule z runnerem zadań i wyrażeniami
 * cron („cron 0 3 * * * kopia zapasowa → jutro 03:00”), a kontrakt wiąże
 * cykliczność z AUTOMATYKĄ, nie z kartą powłoki: zapisuje ją
 * `automation.schedule.set`, czyta `schedule.get`. Rozjazd rozstrzyga się na
 * korzyść kontraktu — rodzina zostaje jedna, a to okno jest jej drugą
 * powierzchnią. Rodziny `terminal.schedule.*` nie ma i nie będzie: dwie rodziny
 * na jeden harmonogram to dwie prawdy o jednym bycie.
 *
 * Zaplanowanie zadania idzie więc dwoma krokami tej samej rodziny:
 * `automation.workflow.save` zapisuje automatykę o jednym kroku rodzaju
 * `command` wołającym `terminal.command.exec`, a `automation.schedule.set`
 * dokłada jej wyrażenie cron. Edytora automatyk to okno nie stawia — pełna praca
 * na krokach, wersjach i zmiennych zostaje w module Automations, a kopia tamtego
 * edytora tutaj rozjechałaby się z pierwowzorem.
 *
 * ── Plan nie jest budzikiem i okno tego nie ukrywa ──────────────────────────
 * Rdzeń wylicza chwilę najbliższego uruchomienia, ale nie ma czym odpalić
 * automatyki samodzielnie: kontrakt nie zna komendy ani zdarzenia, którym
 * harmonogram zgłaszałby wyzwolenie. Zapisany plan jest zapisem obowiązującym
 * wraz z terminem, a uruchomienie prowadzi Operator z okien modułu Automations.
 * Potwierdzenie zapisu mówi to wprost — plan przedstawiony jako budzik byłby
 * obietnicą uruchomienia, którego nikt nie wykona.
 *
 * Kolejki okno nie zakłada: `queue.create` żąda identyfikatora karty sesji.
 * Niesie go od scalenia kontraktu pole `parentSessionId` karty terminala, ale
 * rdzeń jeszcze go nie wypełnia, więc do czasu wpięcia obsługi kolejka powstaje
 * poza tym oknem.
 *
 * Potok jest sekwencjonowaniem po stronie klienta, nie silnikiem w rdzeniu:
 * okno wysyła krok, czeka na jego domknięcie odczytem wyjścia i dopiero wtedy
 * decyduje o następnym. Zamknięcie okna w trakcie przerywa sekwencjonowanie,
 * ale nie przerywa kroku już uruchomionego — ten kończy się w rdzeniu i widać
 * go w Process Monitorze.
 */
export interface OknoZadanIHarmonogramu {
  element: HTMLElement;
  odswiez(): void;
  czynnosci: readonly CzynnoscOkna[];
}

/** Sekcje okna — jedna naraz w ciele, bo każda jest osobnym zadaniem Operatora. */
const SEKCJE: readonly PozycjaWyboru[] = [
  ['zadania', 'Zadania projektu', 'Zadania odczytane z manifestu w katalogu roboczym karty bieżącej.'],
  ['harmonogram', 'Harmonogram i kolejki', 'Odczyt harmonogramów rdzenia i kolejek silnika pętli obsługujących to okno.'],
  ['potok', 'Potok poleceń', 'Sekwencja poleceń wykonywana krok po kroku z oceną kodu wyjścia.'],
  ['historia', 'Historia przebiegów', 'Uruchomienia wydane z tego okna wraz z ich stanem i kodem wyjścia.'],
];

/**
 * Ile milisekund rdzeń ma czekać na domknięcie kroku potoku.
 *
 * Górna granica kontraktu, w odróżnieniu od podglądu wyjścia w Process
 * Monitorze, gdzie czekanie jest krótkie, bo Operator patrzy na okno. Tutaj
 * czekanie jest warunkiem decyzji: bez kodu wyjścia kroku nie ma jak
 * rozstrzygnąć, czy wolno ruszyć z krokiem następnym. Krok, który nie domknie
 * się w tej granicy, zatrzymuje potok wraz ze zdaniem o powodzie — zgadywanie
 * powodzenia byłoby gorsze od zatrzymania.
 */
const CZEKANIE_NA_KROK_MS = 60000;

export function utworzOknoZadanIHarmonogramu(
  zrodlo: ZrodloTerminala,
  stan: StanTerminala,
  pokrycie: PokrycieKomend,
): OknoZadanIHarmonogramu {
  const rama = utworzRameOkna({
    tytul: 'Task & Schedule',
    rola: 'zarządca',
    przeznaczenie:
      'Zadania projektu, harmonogram, kolejka silnika pętli i potok poleceń — wszystko uruchamiane w karcie bieżącej okna wiodącego.',
    kod: 'task-schedule',
    przedrostek: 'dt',
    ogniskowalne: true,
  });
  const tresc = utworzStanTresci();
  const kontrolki = zlozPowierzchnieZadan(rama, tresc.element, pokrycie);

  let zadania: ZadanieManifestu[] = [];
  let harmonogramy: AutomationSchedule[] = [];
  let kolejki: Queue[] = [];
  let obserwacje: readonly TerminalWatch[] = [];
  const wynikiPotoku: WynikKroku[] = [];
  /** Procesy uruchomione z tego okna wraz z nazwą czynności, która je zleciła. */
  const przebiegi = new Map<string, string>();

  function pokaz(): void {
    const sekcja = kontrolki.sekcja.wartosc();
    // Pole kroków potoku znika z pola widzenia atrybutem, a nie wyjęciem
    // z dokumentu: element wyjęty traci ognisko, a przebieg potoku przerysowuje
    // okno po każdym kroku.
    kontrolki.pojemnikPotoku.hidden = sekcja !== 'potok';
    const miejsce = tresc.tresc();
    switch (sekcja) {
      case 'zadania':
        miejsce.append(sekcjaZadan(zadania, (zadanie) => uruchomZadanie(zadanie)));
        break;
      case 'harmonogram':
        miejsce.append(
          sekcjaHarmonogramu(harmonogramy),
          sekcjaKolejek(kolejki, (kolejka, dzialanie) => wykonajNaKolejce(kolejka, dzialanie)),
        );
        break;
      case 'potok':
        miejsce.append(sekcjaWynikowPotoku(wynikiPotoku));
        break;
      case 'obserwacje':
        miejsce.append(sekcjaObserwacji(obserwacje, (kod) => zatrzymajObserwacje(kod)));
        break;
      default:
        miejsce.append(sekcjaHistorii(stan.procesy(), przebiegi, (proces) => powtorz(proces)));
    }
  }

  /** Karta, w której okno uruchamia zadania; brak karty ma własne zdanie. */
  function kartaWykonania(czynnosc: string): TerminalSession | null {
    const karta = stan.kartaBiezaca();
    if (karta === null) {
      tresc.blad(
        `Nie ma karty bieżącej — ${czynnosc} nie ma w czym się wykonać. Otwórz kartę w oknie Terminal Tabs.`,
      );
    }
    return karta;
  }

  /**
   * Wykrycie zadań projektu — odczyt manifestu KOMENDĄ RDZENIA.
   *
   * Do niedawna okno wypisywało manifest poleceniem powłoki, więc wykrycie zadań
   * zależało od tego, czy na maszynie stoi program wypisujący plik, i od składni
   * powłoki karty. `terminal.file.read` czyta plik rdzeniem, więc zależność
   * znika, a odczyt działa jednakowo w każdej powłoce — także tam, gdzie żadnego
   * `cat` ani `Get-Content` nie ma.
   */
  function wykryjZadania(): void {
    const karta = kartaWykonania('wykrycie zadań');
    if (karta === null) return;
    const manifest = kontrolki.manifest.wartosc() as RodzajManifestu;
    tresc.ladowanie(`Odczyt pliku ${manifest} w katalogu roboczym karty ${karta.shell}…`);
    void zrodlo.odczytajPlik({ sessionId: karta.id, path: manifest }).then((wynik) => {
      if (!wynik.udany || wynik.wynik === undefined) {
        tresc.blad(
          `Rdzeń nie oddał treści manifestu ${manifest} (${Command.TerminalFileRead}). ` +
            'Sprawdź, czy plik leży w katalogu roboczym karty bieżącej.',
          wynik.blad,
        );
        return;
      }
      zadania = czytajZadania(manifest, wynik.wynik.content);
      pokaz();
      tresc.potwierdzenie(
        zadania.length === 0
          ? `Plik ${manifest} odczytany, ale nie ma w nim ani jednego zadania do uruchomienia.`
          : `Z pliku ${manifest} weszło ${zadania.length} zadań. Uruchomienie każdego z nich wymaga ` +
              'programu, którego instalka Danaco Console nie niesie.',
        zadania.length > 0,
      );
    });
  }

  /** Czyta obserwacje plików tego okna. */
  function odczytajObserwacje(): void {
    const okno = stan.okno();
    void zrodlo.obserwacje(okno === '' ? {} : { windowId: okno }).then((wynik) => {
      if (!wynik.udany || wynik.wynik === undefined) {
        tresc.blad(`Rdzeń nie oddał wykazu obserwacji (${Command.TerminalWatchList}).`, wynik.blad);
        return;
      }
      obserwacje = wynik.wynik;
      pokaz();
    });
  }

  /** Zakłada obserwację plików wyzwalającą polecenie w karcie bieżącej. */
  function zalozObserwacje(): void {
    const karta = kartaWykonania('założenie obserwacji');
    if (karta === null) return;
    const wzorzec = kontrolki.wzorzec.value.trim();
    const polecenie = kontrolki.polecenieObserwacji.value.trim();
    if (wzorzec === '' || polecenie === '') {
      tresc.potwierdzenie(
        'Obserwacja wymaga wzorca ścieżek i polecenia uruchamianego po zmianie — nie założono.',
        false,
      );
      return;
    }
    const tlumienie = Number.parseInt(kontrolki.tlumienie.value, 10);
    tresc.ladowanie(`Zakładanie obserwacji ${wzorzec} w karcie ${karta.shell}…`);
    void zrodlo
      .zalozObserwacje({
        sessionId: karta.id,
        pattern: wzorzec,
        command: polecenie,
        recursive: kontrolki.rekurencyjnie.dataset['wlaczony'] === 'true',
        ...(Number.isFinite(tlumienie) && tlumienie >= 0 ? { debounceMs: tlumienie } : {}),
      })
      .then((wynik) => {
        if (!wynik.udany || wynik.wynik === undefined) {
          tresc.blad(`Rdzeń nie założył obserwacji (${Command.TerminalWatchStart}).`, wynik.blad);
          return;
        }
        odczytajObserwacje();
        tresc.potwierdzenie(
          `Obserwacja ${wynik.wynik.id} założona — polecenie wykona się w karcie ${karta.shell} ` +
            'przy każdej zmianie pasujących plików. Pierwszy przegląd wyłącznie zapamiętuje stan zastany, ' +
            'więc samo założenie niczego nie uruchamia.',
          true,
        );
      });
  }

  /** Zatrzymuje obserwację; polecenie już uruchomione biegnie dalej. */
  function zatrzymajObserwacje(kod: string): void {
    void zrodlo.zatrzymajObserwacje({ watchId: kod }).then((wynik) => {
      if (!wynik.udany || wynik.wynik === undefined) {
        tresc.blad(`Rdzeń nie zatrzymał obserwacji ${kod} (${Command.TerminalWatchStop}).`, wynik.blad);
        return;
      }
      odczytajObserwacje();
      tresc.potwierdzenie(
        `Obserwacja ${kod} zatrzymana (stan ${wynik.wynik.status}). Polecenie już uruchomione biegnie dalej.`,
        true,
      );
    });
  }

  function uruchomZadanie(zadanie: ZadanieManifestu): void {
    const karta = kartaWykonania(`uruchomienie zadania ${zadanie.nazwa}`);
    if (karta === null) return;
    tresc.ladowanie(`Uruchamianie „${zadanie.polecenie}” w karcie ${karta.shell}…`);
    void zrodlo.wykonaj({ sessionId: karta.id, command: zadanie.polecenie }).then((wynik) => {
      if (!wynik.udany || wynik.wynik === undefined) {
        tresc.blad(
          `Rdzeń nie uruchomił zadania ${zadanie.nazwa} (${Command.TerminalCommandExec}).`,
          wynik.blad,
        );
        return;
      }
      const proces = wynik.wynik;
      stan.zapiszProces(proces);
      stan.zapamietajPolecenie(karta.id, proces.command);
      przebiegi.set(proces.id, `zadanie ${zadanie.nazwa} (${zadanie.zrodlo})`);
      tresc.potwierdzenie(
        `Rdzeń uruchomił proces ${proces.id} dla zadania ${zadanie.nazwa} — stan ${proces.status}. ` +
          'Wynik jedzie do Output Console; wpis stoi w Process Monitorze i w historii przebiegów tego okna.',
        proces.status !== TerminalProcessStatus.Failed,
      );
    });
  }

  function powtorz(proces: TerminalProcess): void {
    const karta = proces.sessionId ?? '';
    if (karta === '') {
      tresc.blad('Ten przebieg nie wskazuje karty źródłowej — nie ma go gdzie powtórzyć.');
      return;
    }
    tresc.ladowanie(`Ponowne uruchomienie „${proces.command}”…`);
    void zrodlo
      .wykonaj({ sessionId: karta, command: proces.command, initiator: proces.initiator })
      .then((wynik) => {
        if (!wynik.udany || wynik.wynik === undefined) {
          tresc.blad(`Rdzeń nie powtórzył przebiegu (${Command.TerminalCommandExec}).`, wynik.blad);
          return;
        }
        stan.zapiszProces(wynik.wynik);
        przebiegi.set(wynik.wynik.id, `ponowienie: ${proces.command}`);
        tresc.potwierdzenie(`Proces ${wynik.wynik.id} uruchomiony ponownie.`, true);
      });
  }

  function odczytajHarmonogramy(): void {
    const automatyka = kontrolki.automatyka.value.trim();
    tresc.ladowanie('Odczyt harmonogramów rdzenia…');
    void zrodlo
      .harmonogramy({
        ...(automatyka === '' ? {} : { workflowId: automatyka }),
        ...(kontrolki.tylkoCzynne.dataset['wlaczony'] === 'true' ? { enabledOnly: true } : {}),
      })
      .then((wynik) => {
        if (!wynik.udany || wynik.wynik === undefined) {
          tresc.blad(`Rdzeń nie oddał harmonogramów (${Command.ScheduleGet}).`, wynik.blad);
          return;
        }
        harmonogramy = [...wynik.wynik];
        pokaz();
        tresc.potwierdzenie(
          harmonogramy.length === 0
            ? 'Rdzeń nie ma harmonogramu spełniającego warunki odczytu. Harmonogram zakłada się w module Automations — kontrakt wiąże go z automatyką, nie z kartą powłoki.'
            : `Rdzeń oddał ${harmonogramy.length} harmonogramów.`,
          true,
        );
      });
  }

  /**
   * Zakłada plan zadania powłoki dwoma krokami JEDNEJ rodziny komend.
   *
   * Kolejność jest wymuszona kontraktem, nie wygodą: cykliczność wskazuje
   * automatykę identyfikatorem, więc automatyka musi istnieć wcześniej. Gdy
   * pierwszy krok przejdzie, a drugi nie, okno mówi o automatyce zapisanej BEZ
   * cykliczności i podaje jej identyfikator — milczenie zostawiłoby Operatorowi
   * automatykę, o której nie wie, a druga próba założyłaby ją po raz drugi.
   *
   * Krok automatyki niesie kartę bieżącą, bo `terminal.command.exec` żąda karty.
   * Karta zamknięta przed uruchomieniem unieważnia plan — i to też jest w zdaniu
   * potwierdzenia, zamiast wyjść dopiero przy uruchomieniu.
   */
  async function zaplanujZadanie(): Promise<void> {
    const karta = kartaWykonania('zaplanowanie zadania powłoki');
    if (karta === null) return;
    const nazwa = kontrolki.nazwaPlanu.value.trim();
    const polecenie = kontrolki.poleceniePlanu.value.trim();
    const cron = kontrolki.wyrazenieCron.value.trim();
    if (nazwa === '' || polecenie === '' || cron === '') {
      tresc.potwierdzenie(
        'Plan wymaga nazwy, polecenia powłoki i wyrażenia cron — nie zapisano żadnego z dwóch kroków.',
        false,
      );
      return;
    }

    tresc.ladowanie(`Zapis automatyki „${nazwa}” dla polecenia „${polecenie}”…`);
    const automatyka = await zrodlo.zapiszAutomatyke({
      name: nazwa,
      description: `Zadanie powłoki zaplanowane z okna Task & Schedule modułu Terminal: ${polecenie}`,
      steps: [krokPoleceniaPowloki(karta.id, polecenie)],
      enabled: true,
    });
    if (!automatyka.udany || automatyka.wynik === undefined) {
      tresc.blad(
        `Rdzeń nie zapisał automatyki planu (${Command.AutomationWorkflowSave}). ` +
          'Bez niej nie ma czego objąć cyklicznością — drugi krok nie poszedł.',
        automatyka.blad,
      );
      return;
    }
    const zapisana = automatyka.wynik;

    const strefa = kontrolki.strefaPlanu.value.trim();
    tresc.ladowanie(`Zapis cykliczności „${cron}” automatyki ${zapisana.id}…`);
    const plan = await zrodlo.ustawHarmonogram({
      workflowId: zapisana.id,
      cron,
      enabled: true,
      ...(strefa === '' ? {} : { timeZone: strefa }),
    });
    if (!plan.udany || plan.wynik === undefined) {
      tresc.blad(
        `Automatyka ${zapisana.id} („${nazwa}") jest zapisana, ale cykliczności nie przyjęła ` +
          `(${Command.AutomationScheduleSet}). Powtórz zapis cykliczności, wpisując ten ` +
          'identyfikator w module Automations — powtórzenie planu z tego okna założyłoby automatykę drugi raz.',
        plan.blad,
      );
      return;
    }

    harmonogramy = [plan.wynik, ...harmonogramy.filter((pozycja) => pozycja.id !== plan.wynik?.id)];
    pokaz();
    tresc.potwierdzenie(
      `Plan zapisany: automatyka ${zapisana.id} („${nazwa}") z cyklicznością ${cron}, ` +
        `najbliższe uruchomienie ${
          plan.wynik.nextRunAt === undefined ? 'nie wyliczone' : chwila(plan.wynik.nextRunAt)
        }. ` +
        'Rdzeń nie ma budzika, który odpali ją sam — uruchomienie prowadzi Operator z okien modułu ' +
        `Automations. Krok wskazuje kartę ${karta.shell} (${karta.id}); karta zamknięta przed ` +
        'uruchomieniem unieważnia plan.',
      true,
    );
  }

  function odczytajKolejki(): void {
    const okno = stan.okno();
    if (okno === '') {
      tresc.blad('Moduł nie zna okna komunikacji — nie ma czego zapytać o kolejki.');
      return;
    }
    tresc.ladowanie('Odczyt kolejek silnika pętli…');
    void zrodlo.kolejki({ windowId: okno }).then((wynik) => {
      if (!wynik.udany || wynik.wynik === undefined) {
        tresc.blad(`Rdzeń nie oddał kolejek (${Command.QueueList}).`, wynik.blad);
        return;
      }
      kolejki = [...wynik.wynik];
      pokaz();
      tresc.potwierdzenie(
        kolejki.length === 0
          ? 'To okno nie ma ani jednej kolejki silnika pętli. Kolejkę zakłada komenda zakładania kolejki, która żąda identyfikatora karty sesji; niesie go pole karty terminala, ale rdzeń jeszcze go nie wypełnia.'
          : `Rdzeń oddał ${kolejki.length} kolejek obsługujących to okno.`,
        true,
      );
    });
  }

  function wykonajNaKolejce(kolejka: Queue, dzialanie: QueueAction): void {
    tresc.ladowanie(`Działanie ${dzialanie} na kolejce ${kolejka.id}…`);
    void zrodlo.dzialanieKolejki({ queueId: kolejka.id, action: dzialanie }).then((wynik) => {
      if (!wynik.udany || wynik.wynik === undefined) {
        tresc.blad(`Rdzeń nie wykonał działania ${dzialanie} (${Command.QueueAction}).`, wynik.blad);
        return;
      }
      const zapisana = wynik.wynik;
      const miejsce = kolejki.findIndex((pozycja) => pozycja.id === zapisana.id);
      if (miejsce < 0) kolejki.push(zapisana);
      else kolejki.splice(miejsce, 1, zapisana);
      pokaz();
      // Stan bierze się z odpowiedzi, nie z żądania: silnik potrafi przyjąć
      // wstrzymanie i oddać kolejkę nadal biegnącą.
      tresc.potwierdzenie(
        `Rdzeń oddał kolejkę ${zapisana.id} w stanie ${zapisana.status} po działaniu ${dzialanie}.`,
        true,
      );
    });
  }

  /** Wstrzymanie albo wznowienie pierwszej kolejki okna — nośnik skrótu klawiszowego. */
  function przestawKolejke(): void {
    const kolejka = kolejki[0];
    if (kolejka === undefined) {
      tresc.potwierdzenie(
        'To okno nie ma odczytanej kolejki — najpierw odczytaj kolejki, potem wstrzymuj albo wznawiaj.',
        false,
      );
      return;
    }
    wykonajNaKolejce(
      kolejka,
      kolejka.status === QueueStatus.Paused ? QueueAction.Resume : QueueAction.Pause,
    );
  }

  async function uruchomPotok(): Promise<void> {
    const karta = kartaWykonania('uruchomienie potoku');
    if (karta === null) return;
    const kroki = kontrolki.potok.value
      .split('\n')
      .map((wiersz) => wiersz.trim())
      .filter((wiersz) => wiersz !== '' && !wiersz.startsWith('#'));
    if (kroki.length === 0) {
      tresc.blad('Potok nie ma ani jednego kroku — wpisz polecenia, po jednym w wierszu.');
      return;
    }
    const przerywaj = kontrolki.przerywanie.dataset['wlaczony'] === 'true';
    wynikiPotoku.length = 0;
    for (const [numer, polecenie] of kroki.entries()) {
      tresc.ladowanie(`Potok — krok ${numer + 1} z ${kroki.length}: „${polecenie}”…`);
      const wynik = await wykonajKrok(zrodlo, stan, karta.id, polecenie);
      wynikiPotoku.push(wynik);
      if (wynik.proces !== '') przebiegi.set(wynik.proces, `krok potoku: ${polecenie}`);
      pokaz();
      if (!wynik.powodzenie && przerywaj) {
        tresc.potwierdzenie(
          `Potok zatrzymany na kroku ${numer + 1} z ${kroki.length}: ${wynik.zdanie} ` +
            'Kroki dalsze nie zostały uruchomione.',
          false,
        );
        return;
      }
    }
    const nieudane = wynikiPotoku.filter((wynik) => !wynik.powodzenie).length;
    pokaz();
    tresc.potwierdzenie(
      nieudane === 0
        ? `Potok przeszedł w całości — ${kroki.length} kroków, każdy z zerowym kodem wyjścia.`
        : `Potok doszedł do końca; ${nieudane} z ${kroki.length} kroków nie powiodło się.`,
      nieudane === 0,
    );
  }

  kontrolki.sekcja.naZmiane(pokaz);
  kontrolki.wykryj.addEventListener('click', wykryjZadania);
  kontrolki.harmonogram.addEventListener('click', odczytajHarmonogramy);
  kontrolki.zaplanuj.addEventListener('click', () => void zaplanujZadanie());
  kontrolki.kolejkiOdczyt.addEventListener('click', odczytajKolejki);
  kontrolki.przestawKolejke.addEventListener('click', przestawKolejke);
  kontrolki.zalozObserwacje.addEventListener('click', zalozObserwacje);
  kontrolki.odczytajObserwacje.addEventListener('click', odczytajObserwacje);
  kontrolki.rekurencyjnie.addEventListener('click', () => przestaw(kontrolki.rekurencyjnie));
  kontrolki.uruchomPotok.addEventListener('click', () => void uruchomPotok());
  kontrolki.przerywanie.addEventListener('click', () => przestaw(kontrolki.przerywanie));
  kontrolki.tylkoCzynne.addEventListener('click', () => przestaw(kontrolki.tylkoCzynne));
  kontrolki.eksportHistorii.addEventListener('click', () => {
    const pozycje = stan.procesy().filter((proces) => przebiegi.has(proces.id));
    if (pozycje.length === 0) {
      tresc.potwierdzenie('To okno nie wydało jeszcze ani jednego uruchomienia — pliku nie zapisano.', false);
      return;
    }
    const nazwa = `historia-przebiegow-${Date.now()}.json`;
    pobierzPlik(
      nazwa,
      JSON.stringify(
        pozycje.map((proces) => ({ proces, czynnosc: przebiegi.get(proces.id) ?? '' })),
        null,
        2,
      ),
      'text/plain',
    );
    tresc.potwierdzenie(`Zapisano ${nazwa} — ${pozycje.length} przebiegów tego okna.`, true);
  });

  stan.naZmiane(() => {
    if (kontrolki.sekcja.wartosc() === 'historia') pokaz();
  });

  const czynnosci: readonly CzynnoscOkna[] = [
    {
      okno: 'Task & Schedule',
      nazwa: 'Załóż obserwację plików',
      opis: 'Uruchamia polecenie w karcie bieżącej przy każdej zmianie pasujących plików.',
      warstwa: 'kontekstowa',
      wykonaj: () => kontrolki.zalozObserwacje.click(),
    },
    {
      okno: 'Task & Schedule',
      nazwa: 'Odczytaj obserwacje plików',
      opis: 'Czyta obserwacje tego okna wraz z licznikiem wyzwoleń i powodem niepowodzenia.',
      warstwa: 'kontekstowa',
      wykonaj: () => kontrolki.odczytajObserwacje.click(),
    },
    {
      okno: 'Task & Schedule',
      nazwa: 'Wykryj zadania projektu',
      opis: 'Wypisuje manifest w karcie bieżącej i składa z jego treści wykaz zadań do uruchomienia.',
      warstwa: 'zawsze',
      wykonaj: wykryjZadania,
    },
    {
      okno: 'Task & Schedule',
      nazwa: 'Uruchom potok poleceń',
      opis: 'Wykonuje wpisane polecenia krok po kroku, oceniając kod wyjścia każdego kroku.',
      warstwa: 'na-zadanie',
      wykonaj: () => void uruchomPotok(),
    },
    {
      okno: 'Task & Schedule',
      nazwa: 'Odczytaj harmonogramy',
      opis: 'Pyta rdzeń o harmonogramy komendą schedule.get.',
      warstwa: 'na-zadanie',
      wykonaj: odczytajHarmonogramy,
    },
    {
      okno: 'Task & Schedule',
      nazwa: 'Zaplanuj zadanie powłoki',
      opis:
        'Zapisuje automatykę o jednym kroku wołającym polecenie powłoki i dokłada jej wyrażenie cron ' +
        '— tą samą rodziną komend, którą prowadzi moduł Automations.',
      warstwa: 'na-zadanie',
      wykonaj: () => void zaplanujZadanie(),
    },
    {
      okno: 'Task & Schedule',
      nazwa: 'Odczytaj kolejki okna',
      opis: 'Pyta rdzeń o kolejki silnika pętli obsługujące to okno komendą queue.list.',
      warstwa: 'na-zadanie',
      wykonaj: odczytajKolejki,
    },
    {
      okno: 'Task & Schedule',
      nazwa: 'Wstrzymaj albo wznów kolejkę',
      opis: 'Przestawia pierwszą odczytaną kolejkę okna między stanem wstrzymanym a biegnącym.',
      warstwa: 'kontekstowa',
      skrot: 'Ctrl/Cmd + .',
      wykonaj: przestawKolejke,
    },
    {
      okno: 'Task & Schedule',
      nazwa: 'Eksportuj historię przebiegów',
      opis: 'Zapisuje uruchomienia wydane z tego okna wraz z ich stanem i kodem wyjścia.',
      warstwa: 'kontekstowa',
      wykonaj: () => kontrolki.eksportHistorii.click(),
    },
  ];

  return { element: rama.element, odswiez: pokaz, czynnosci };
}

/** Wynik jednego kroku potoku — treść wiersza wykazu i rozstrzygnięcie o kontynuacji. */
interface WynikKroku {
  polecenie: string;
  /** Proces rejestru rdzenia; pusty, gdy rdzeń nie przyjął uruchomienia. */
  proces: string;
  powodzenie: boolean;
  zdanie: string;
}

/**
 * Jeden krok potoku: uruchomienie polecenia i odczyt jego wyniku.
 *
 * Rozstrzygnięcie bierze się z kodu wyjścia, nie z faktu, że żądanie poszło.
 * Krok, który nie domknął się w granicy czekania, nie jest ani udany, ani
 * nieudany — okno mówi to wprost i zatrzymuje potok, zamiast zgadywać.
 */
async function wykonajKrok(
  zrodlo: ZrodloTerminala,
  stan: StanTerminala,
  idKarty: string,
  polecenie: string,
): Promise<WynikKroku> {
  const uruchomienie = await zrodlo.wykonaj({ sessionId: idKarty, command: polecenie });
  if (!uruchomienie.udany || uruchomienie.wynik === undefined) {
    const powod = uruchomienie.blad;
    return {
      polecenie,
      proces: '',
      powodzenie: false,
      zdanie:
        powod === undefined
          ? `Rdzeń nie uruchomił kroku (${Command.TerminalCommandExec}).`
          : `Rdzeń nie uruchomił kroku: ${powod.message} (kod ${powod.code}).`,
    };
  }
  const proces = uruchomienie.wynik;
  stan.zapiszProces(proces);

  const odczyt = await zrodlo.odczytajWyjscie({
    processId: proces.id,
    waitMs: CZEKANIE_NA_KROK_MS,
  });
  if (!odczyt.udany || odczyt.wynik === undefined) {
    const powod = odczyt.blad;
    return {
      polecenie,
      proces: proces.id,
      powodzenie: false,
      zdanie:
        powod === undefined
          ? `Rdzeń nie oddał wyniku kroku (${Command.TerminalOutputRead}); proces ${proces.id} biegnie dalej.`
          : `Rdzeń nie oddał wyniku kroku: ${powod.message} (kod ${powod.code}).`,
    };
  }
  const wynik = odczyt.wynik;
  if (wynik.status === TerminalProcessStatus.Running) {
    return {
      polecenie,
      proces: proces.id,
      powodzenie: false,
      zdanie: `Krok nie domknął się w granicy ${CZEKANIE_NA_KROK_MS} ms — proces ${proces.id} biegnie dalej, a kodu wyjścia jeszcze nie ma.`,
    };
  }
  const kod = wynik.exitCode;
  return {
    polecenie,
    proces: proces.id,
    powodzenie: wynik.status === TerminalProcessStatus.Finished && (kod === undefined || kod === 0),
    zdanie: `Stan ${wynik.status}, ${kod === undefined ? 'kodu wyjścia rdzeń nie podał' : `kod wyjścia ${kod}`}.`,
  };
}

/** Wykaz zadań projektu wraz z uruchomieniem każdego z nich. */
function sekcjaZadan(
  zadania: readonly ZadanieManifestu[],
  uruchom: (zadanie: ZadanieManifestu) => void,
): HTMLElement {
  const blok = document.createElement('section');
  blok.className = 'dt-zadania';

  const podpis = document.createElement('h4');
  podpis.className = 'dt-zadania__podpis';
  podpis.textContent = 'Zadania projektu';
  blok.append(podpis);

  if (zadania.length === 0) {
    const pusto = document.createElement('p');
    pusto.className = 'dn-pole-opis';
    pusto.textContent =
      'Wykaz jest pusty. Wskaż manifest i naciśnij „Wykryj zadania” — okno wypisze plik w karcie bieżącej ' +
      'i złoży wykaz z jego treści. Kontrakt nie ma komendy przeglądającej katalog, a odczytu pliku okno ' +
      'jeszcze nie wywołuje, więc na dziś to jedyna droga.';
    blok.append(pusto);
    return blok;
  }

  const lista = wykaz('Zadania wykryte w manifeście', 'dt-wykaz');
  for (const zadanie of zadania) {
    const pozycja = pozycjaWykazu(zadanie.nazwa, `${zadanie.zrodlo} · ${zadanie.polecenie}`, 'dt');
    const uruchomienie = przyciskAkcji('Uruchom', 'dn-btn dn-btn--atrament');
    uruchomienie.title = `Wysyła „${zadanie.polecenie}” do karty bieżącej. Program uruchamiający zadanie musi leżeć na maszynie rdzenia.`;
    uruchomienie.addEventListener('click', () => uruchom(zadanie));
    pozycja.akcje.append(uruchomienie);
    lista.append(pozycja.element);
  }
  blok.append(lista);
  return blok;
}

/** Wykaz harmonogramów oddanych przez rdzeń. */
function sekcjaHarmonogramu(harmonogramy: readonly AutomationSchedule[]): HTMLElement {
  const blok = document.createElement('section');
  blok.className = 'dt-harmonogram';

  const podpis = document.createElement('h4');
  podpis.className = 'dt-harmonogram__podpis';
  podpis.textContent = 'Harmonogramy rdzenia (odczyt)';
  blok.append(podpis);

  const zastrzezenie = document.createElement('p');
  zastrzezenie.className = 'dn-pole-opis';
  zastrzezenie.textContent =
    'Kontrakt wiąże harmonogram z automatyką, a nie z kartą powłoki, więc „Zaplanuj zadanie powłoki” ' +
    'zapisuje automatykę o jednym kroku wołającym polecenie i dokłada jej wyrażenie cron — tą samą ' +
    'rodziną komend, którą prowadzi moduł Automations. Zapis jest planem wraz z terminem, nie budzikiem: ' +
    'rdzeń nie ma czym odpalić automatyki samodzielnie, więc uruchomienie prowadzi Operator z okien ' +
    'modułu Automations. Pełna praca na krokach, wersjach i zmiennych automatyki też należy do tamtego modułu.';
  blok.append(zastrzezenie);

  if (harmonogramy.length === 0) {
    const pusto = document.createElement('p');
    pusto.className = 'dn-pole-opis';
    pusto.textContent = 'Rdzeń nie oddał jeszcze żadnego harmonogramu — naciśnij „Odczytaj harmonogramy”.';
    blok.append(pusto);
    return blok;
  }

  const lista = wykaz('Harmonogramy rdzenia', 'dt-wykaz');
  for (const harmonogram of harmonogramy) {
    lista.append(pozycjaWykazu(harmonogram.id, opisHarmonogramu(harmonogram), 'dt').element);
  }
  blok.append(lista);
  return blok;
}

/**
 * Krok automatyki wołający polecenie powłoki.
 *
 * Rodzaj `command` znaczy „krok wywołuje komendę kontraktu”, więc krok niesie
 * nazwę komendy i treść jej żądania — a nie polecenie powłoki wprost. Dzięki temu
 * plan przechodzi tą samą bramą uprawnień i tym samym egzekutorem izolacji, co
 * polecenie wydane ręcznie z karty; krok omijający komendę byłby drugą drogą do
 * powłoki, bez ani jednego z tych sprawdzeń.
 *
 * Identyfikator kroku jest nazwą czynności, nie liczbą porządkową: automatyka
 * planu ma dokładnie jeden krok, więc numer nie rozróżniałby niczego.
 */
function krokPoleceniaPowloki(idKarty: string, polecenie: string): AutomationStep {
  return {
    id: 'polecenie-powloki',
    name: 'Polecenie powłoki',
    kind: AutomationStepKind.Command,
    command: Command.TerminalCommandExec,
    params: { sessionId: idKarty, command: polecenie },
    order: 1,
  };
}

function opisHarmonogramu(harmonogram: AutomationSchedule): string {
  const czesci = [
    `automatyka: ${harmonogram.workflowId}`,
    `cykliczność: ${harmonogram.cron ?? 'bez zapisu cron'}`,
    `strefa: ${harmonogram.timeZone ?? 'nie podana'}`,
    harmonogram.enabled ? 'obowiązuje' : 'wstrzymany',
    `najbliższe uruchomienie: ${
      harmonogram.nextRunAt === undefined ? 'nie wyliczone' : chwila(harmonogram.nextRunAt)
    }`,
  ];
  const wyzwalacze = harmonogram.triggers ?? [];
  if (wyzwalacze.length > 0) {
    czesci.push(`wyzwalacze: ${wyzwalacze.map((wyzwalacz) => wyzwalacz.kind).join(', ')}`);
  }
  return czesci.join(' · ');
}

/** Wykaz kolejek silnika pętli wraz z sześcioma działaniami słownika kontraktu. */
function sekcjaKolejek(
  kolejki: readonly Queue[],
  wykonaj: (kolejka: Queue, dzialanie: QueueAction) => void,
): HTMLElement {
  const blok = document.createElement('section');
  blok.className = 'dt-kolejki';

  const podpis = document.createElement('h4');
  podpis.className = 'dt-kolejki__podpis';
  podpis.textContent = 'Kolejki silnika pętli obsługujące to okno';
  blok.append(podpis);

  if (kolejki.length === 0) {
    const pusto = document.createElement('p');
    pusto.className = 'dn-pole-opis';
    pusto.textContent =
      'Rdzeń nie oddał jeszcze kolejek — naciśnij „Odczytaj kolejki”. Kolejkę zakłada komenda żądająca ' +
      'identyfikatora karty sesji; niesie go pole karty terminala, ale rdzeń jeszcze go nie wypełnia, ' +
      'więc do tego czasu kolejka powstaje poza tym oknem.';
    blok.append(pusto);
    return blok;
  }

  const lista = wykaz('Kolejki okna', 'dt-wykaz');
  for (const kolejka of kolejki) {
    const pozycja = pozycjaWykazu(kolejka.name ?? kolejka.id, opisKolejki(kolejka), 'dt');
    pozycja.element.dataset['stan'] = kolejka.status;
    for (const dzialanie of Object.values(QueueAction)) {
      const przycisk = przyciskAkcji(nazwaDzialania(dzialanie));
      przycisk.title = `Wysyła działanie ${dzialanie} do silnika kolejek (${Command.QueueAction}).`;
      przycisk.addEventListener('click', () => wykonaj(kolejka, dzialanie));
      pozycja.akcje.append(przycisk);
    }
    lista.append(pozycja.element);
  }
  blok.append(lista);
  return blok;
}

function opisKolejki(kolejka: Queue): string {
  const czesci = [`stan: ${kolejka.status}`, `sesja: ${kolejka.sessionId}`];
  if (kolejka.cycle !== undefined) czesci.push(`obiegi naprawcze: ${kolejka.cycle}`);
  czesci.push(`ostatnia zmiana: ${chwila(kolejka.updatedAt)}`);
  return czesci.join(' · ');
}

/** Nazwa działania kolejki po polsku; wartość kontraktu zostaje w podpowiedzi przycisku. */
function nazwaDzialania(dzialanie: QueueAction): string {
  switch (dzialanie) {
    case QueueAction.Start:
      return 'Uruchom';
    case QueueAction.Pause:
      return 'Wstrzymaj';
    case QueueAction.Resume:
      return 'Wznów';
    case QueueAction.Stop:
      return 'Zatrzymaj';
    case QueueAction.Retry:
      return 'Ponów';
    default:
      return 'Wyczyść';
  }
}

/** Wykaz wyników potoku — jeden wiersz na krok, w kolejności wykonania. */
function sekcjaWynikowPotoku(wyniki: readonly WynikKroku[]): HTMLElement {
  const blok = document.createElement('section');
  blok.className = 'dt-potok';

  const podpis = document.createElement('h4');
  podpis.className = 'dt-potok__podpis';
  podpis.textContent = 'Kroki ostatniego przebiegu potoku';
  blok.append(podpis);

  if (wyniki.length === 0) {
    const pusto = document.createElement('p');
    pusto.className = 'dn-pole-opis';
    pusto.textContent =
      'Potok nie był jeszcze uruchamiany w tym oknie. Wpisz polecenia po jednym w wierszu; wiersz zaczynający ' +
      'się od kratki jest komentarzem i nie trafia do powłoki.';
    blok.append(pusto);
    return blok;
  }

  const lista = wykaz('Kroki potoku', 'dt-wykaz');
  for (const [numer, wynik] of wyniki.entries()) {
    const pozycja = pozycjaWykazu(`Krok ${numer + 1}: ${wynik.polecenie}`, wynik.zdanie, 'dt');
    pozycja.element.dataset['stan'] = wynik.powodzenie ? 'finished' : 'failed';
    lista.append(pozycja.element);
  }
  blok.append(lista);
  return blok;
}

/** Historia przebiegów wydanych z tego okna, złożona z rejestru procesów rdzenia. */
/**
 * Sekcja obserwacji plików.
 *
 * Licznik wyzwoleń stoi w opisie, bo bez niego obserwacja założona
 * i niedziałająca wygląda tak samo jak taka, która nie miała czego złapać.
 * Powód niepowodzenia też: stan `failed` bez powodu mówi wyłącznie, że coś nie
 * wyszło.
 */
function sekcjaObserwacji(
  obserwacje: readonly TerminalWatch[],
  zatrzymaj: (kod: string) => void,
): HTMLElement {
  const blok = document.createElement('section');
  blok.className = 'dt-obserwacje';

  const podpis = document.createElement('h4');
  podpis.className = 'dt-obserwacje__podpis';
  podpis.textContent = 'Obserwacje plików';
  blok.append(podpis);

  if (obserwacje.length === 0) {
    const pusto = document.createElement('p');
    pusto.className = 'dn-pole-opis';
    pusto.textContent =
      'Nie założono ani jednej obserwacji. Wzorzec ścieżek i polecenie stoją w pasie narzędzi okna; ' +
      'polecenie wykona się w karcie bieżącej przy każdej zmianie pasujących plików.';
    blok.append(pusto);
    return blok;
  }

  const lista = wykaz('Obserwacje plików', 'dt-wykaz');
  for (const obserwacja of obserwacje) {
    const czesci = [
      `wzorzec: ${obserwacja.pattern}`,
      `stan: ${obserwacja.status}`,
      `wyzwoleń: ${obserwacja.triggerCount ?? 0}`,
    ];
    if (obserwacja.recursive === true) czesci.push('obejmuje podkatalogi');
    if (obserwacja.debounceMs !== undefined) czesci.push(`tłumienie: ${obserwacja.debounceMs} ms`);
    if (obserwacja.errorMessage !== undefined && obserwacja.errorMessage !== '') {
      czesci.push(`powód: ${obserwacja.errorMessage}`);
    }
    const pozycja = pozycjaWykazu(obserwacja.command, czesci.join(' · '), 'dt');
    if (obserwacja.status === TerminalWatchStatus.Active) {
      const zatrzymaj_ = przyciskAkcji('Zatrzymaj obserwację');
      zatrzymaj_.addEventListener('click', () => zatrzymaj(obserwacja.id));
      pozycja.akcje.append(zatrzymaj_);
    }
    lista.append(pozycja.element);
  }
  blok.append(lista);
  return blok;
}

function sekcjaHistorii(
  procesy: readonly TerminalProcess[],
  przebiegi: ReadonlyMap<string, string>,
  powtorz: (proces: TerminalProcess) => void,
): HTMLElement {
  const blok = document.createElement('section');
  blok.className = 'dt-historia';

  const podpis = document.createElement('h4');
  podpis.className = 'dt-historia__podpis';
  podpis.textContent = 'Przebiegi wydane z tego okna';
  blok.append(podpis);

  const wybrane = procesy.filter((proces) => przebiegi.has(proces.id));
  if (wybrane.length === 0) {
    const pusto = document.createElement('p');
    pusto.className = 'dn-pole-opis';
    pusto.textContent =
      'To okno nie wydało jeszcze ani jednego uruchomienia. Historia liczy przebiegi tego okna; komplet procesów ' +
      'rejestru rdzenia — także uruchomionych gdzie indziej — stoi w Process Monitorze.';
    blok.append(pusto);
    return blok;
  }

  const lista = wykaz('Historia przebiegów okna', 'dt-wykaz');
  for (const proces of wybrane) {
    const pozycja = pozycjaWykazu(proces.command, opisPrzebiegu(proces, przebiegi), 'dt');
    pozycja.element.dataset['stan'] = proces.status;
    const ponowienie = przyciskAkcji('Ponów');
    ponowienie.addEventListener('click', () => powtorz(proces));
    pozycja.akcje.append(ponowienie);
    lista.append(pozycja.element);
  }
  blok.append(lista);
  return blok;
}

function opisPrzebiegu(proces: TerminalProcess, przebiegi: ReadonlyMap<string, string>): string {
  const czesci = [
    przebiegi.get(proces.id) ?? 'uruchomienie tego okna',
    `stan: ${proces.status}`,
    `kod wyjścia: ${proces.exitCode ?? '—'}`,
    `start: ${chwila(proces.startedAt)}`,
  ];
  if (proces.finishedAt !== undefined) czesci.push(`koniec: ${chwila(proces.finishedAt)}`);
  return czesci.join(' · ');
}

function chwila(znacznik: number): string {
  if (znacznik <= 0) return '—';
  return new Date(znacznik).toISOString().replace('T', ' ').slice(0, 19);
}

/** Kontrolki okna Task & Schedule. */
interface PowierzchniaZadan {
  sekcja: WyborDrzewem;
  manifest: WyborDrzewem;
  automatyka: HTMLInputElement;
  tylkoCzynne: HTMLButtonElement;
  nazwaPlanu: HTMLInputElement;
  poleceniePlanu: HTMLInputElement;
  wyrazenieCron: HTMLInputElement;
  strefaPlanu: HTMLInputElement;
  zaplanuj: HTMLButtonElement;
  potok: HTMLTextAreaElement;
  /** Obudowa pola potoku — chowana atrybutem przy sekcjach innych niż potok. */
  pojemnikPotoku: HTMLElement;
  przerywanie: HTMLButtonElement;
  wykryj: HTMLButtonElement;
  harmonogram: HTMLButtonElement;
  kolejkiOdczyt: HTMLButtonElement;
  przestawKolejke: HTMLButtonElement;
  uruchomPotok: HTMLButtonElement;
  eksportHistorii: HTMLButtonElement;
  wzorzec: HTMLInputElement;
  polecenieObserwacji: HTMLInputElement;
  tlumienie: HTMLInputElement;
  rekurencyjnie: HTMLButtonElement;
  zalozObserwacje: HTMLButtonElement;
  odczytajObserwacje: HTMLButtonElement;
}

/**
 * Składa kontrolki, pasek akcji, pasek narzędzi i ciało okna.
 *
 * Pole potoku powstaje raz i wchodzi do ciała dopiero przy swojej sekcji —
 * budowane od nowa przy każdym przerysowaniu gubiłoby wpisaną treść pod ręką
 * Operatora w trakcie przebiegu.
 */
function zlozPowierzchnieZadan(
  rama: { akcje: HTMLElement; narzedzia: HTMLElement; cialo: HTMLElement },
  stanTresci: HTMLElement,
  pokrycie: PokrycieKomend,
): PowierzchniaZadan {
  const sekcja = utworzWyborDrzewem({ nastawa: 'Sekcja okna', pozycje: SEKCJE });
  const manifest = utworzWyborDrzewem({ nastawa: 'Manifest zadań projektu', pozycje: MANIFESTY });
  const automatyka = pole('Automatyka odczytu harmonogramu', 'puste = wszystkie harmonogramy');
  const tylkoCzynne = przelacznikWidoku('Tylko harmonogramy obowiązujące', false);
  const nazwaPlanu = pole('Nazwa planowanego zadania', 'np. kopia zapasowa');
  const poleceniePlanu = pole('Polecenie powłoki objęte planem', 'np. tar -czf kopia.tgz dane/');
  const wyrazenieCron = pole('Wyrażenie cron', 'np. 0 3 * * * — codziennie o 03:00');
  const strefaPlanu = pole('Strefa czasowa planu', 'puste = strefa rdzenia');
  const potok = poleTresci(
    'Kroki potoku — jedno polecenie w wierszu',
    8,
    'npm ci\nnpm run build\n# wiersz z kratką jest komentarzem',
    'dt-potok__pole',
  );
  const przerywanie = przelacznikWidoku('Zatrzymaj na pierwszym błędzie', true);

  const wykryj = przyciskAkcji('Wykryj zadania', 'dn-btn dn-btn--atrament');
  const harmonogram = przyciskAkcji('Odczytaj harmonogramy');
  const zaplanuj = przyciskAkcji('Zaplanuj zadanie powłoki');
  zaplanuj.title =
    'Zapisuje automatykę o jednym kroku wołającym polecenie powłoki i dokłada jej wyrażenie cron — ' +
    'tą samą rodziną komend, którą prowadzi moduł Automations. Zapis jest planem wraz z terminem, ' +
    'nie budzikiem: uruchomienie prowadzi Operator z okien modułu Automations.';
  const kolejkiOdczyt = przyciskAkcji('Odczytaj kolejki');
  const przestawKolejke = przyciskAkcji('Wstrzymaj / wznów kolejkę');
  const uruchomPotok = przyciskAkcji('Uruchom potok', 'dn-btn dn-btn--atrament');
  const eksportHistorii = przyciskAkcji('Eksportuj historię przebiegów');

  const wzorzec = pole('Wzorzec ścieżek obserwacji', 'np. *.go — względem katalogu karty');
  const polecenieObserwacji = pole('Polecenie po zmianie plików', 'np. go build ./...');
  const tlumienie = poleLiczbowe('Tłumienie powtórzeń [ms]', 'puste = wartość rdzenia');
  const rekurencyjnie = przelacznikWidoku('Obserwuj także podkatalogi', false);
  const zalozObserwacje = przyciskAkcji('Załóż obserwację plików');
  zalozObserwacje.title =
    'Uruchamia polecenie w karcie bieżącej przy każdej zmianie pasujących plików. Wyzwolenie idzie tą ' +
    'samą drogą co polecenie wydane ręcznie: przez bramę uprawnień okna i egzekutor izolacji.';
  const odczytajObserwacje = przyciskAkcji('Odczytaj obserwacje');

  // Nazwa spoza kontraktu jest tu wskazaniem, nie zapisem stanu: próg
  // powiadomienia świadomie nie poszedł do scalenia, bo kanał powiadomień jest
  // własnością platformy, a nie modułu.
  const powiadomienia = pokrycie.przycisk(
    'Powiadomienie o długim zadaniu',
    'terminal.notification.set',
    'Próg czasu trwania, po którym zadanie zgłasza się powiadomieniem; czas trwania widać dziś w Process Monitorze',
  );

  oznaczWarstwy([
    [wykryj, 'zawsze'],
    [uruchomPotok, 'na-zadanie'],
    [harmonogram, 'na-zadanie'],
    [zaplanuj, 'na-zadanie'],
    [kolejkiOdczyt, 'na-zadanie'],
    [przerywanie, 'na-zadanie'],
    [przestawKolejke, 'kontekstowa'],
    [eksportHistorii, 'kontekstowa'],
    [automatyka, 'kontekstowa'],
    [tylkoCzynne, 'kontekstowa'],
    [nazwaPlanu, 'kontekstowa'],
    [poleceniePlanu, 'kontekstowa'],
    [wyrazenieCron, 'kontekstowa'],
    [strefaPlanu, 'ekspercka'],
    [zalozObserwacje, 'kontekstowa'],
    [odczytajObserwacje, 'kontekstowa'],
    [wzorzec, 'kontekstowa'],
    [polecenieObserwacji, 'kontekstowa'],
    [tlumienie, 'ekspercka'],
    [rekurencyjnie, 'ekspercka'],
    [powiadomienia, 'ekspercka'],
  ]);

  rama.akcje.append(
    wykryj,
    uruchomPotok,
    harmonogram,
    zaplanuj,
    kolejkiOdczyt,
    przestawKolejke,
    eksportHistorii,
    zalozObserwacje,
    odczytajObserwacje,
    powiadomienia,
  );
  rama.narzedzia.append(
    sekcja.element,
    manifest.element,
    automatyka,
    tylkoCzynne,
    nazwaPlanu,
    poleceniePlanu,
    wyrazenieCron,
    strefaPlanu,
    przerywanie,
    wzorzec,
    polecenieObserwacji,
    tlumienie,
    rekurencyjnie,
  );
  const pojemnikPotoku = document.createElement('section');
  pojemnikPotoku.className = 'dt-potok';
  const podpisPotoku = document.createElement('h4');
  podpisPotoku.className = 'dt-potok__podpis';
  podpisPotoku.textContent = 'Kroki potoku';
  pojemnikPotoku.append(podpisPotoku, potok);
  pojemnikPotoku.hidden = true;

  rama.cialo.append(
    notaZaleznosci(
      'Zadanie projektu i krok potoku wykonuje program leżący na maszynie rdzenia. Sam manifest czyta ' +
        'rdzeń komendą odczytu pliku, więc wykrycie zadań nie zależy już ani od programu wypisującego ' +
        'plik, ani od składni powłoki karty — zależy od niego dopiero uruchomienie zadania.',
      PROGRAMY_CZYNNOSCI,
    ),
    pojemnikPotoku,
    stanTresci,
  );

  return {
    sekcja,
    manifest,
    automatyka,
    tylkoCzynne,
    nazwaPlanu,
    poleceniePlanu,
    wyrazenieCron,
    strefaPlanu,
    zaplanuj,
    potok,
    pojemnikPotoku,
    przerywanie,
    wykryj,
    harmonogram,
    kolejkiOdczyt,
    przestawKolejke,
    uruchomPotok,
    eksportHistorii,
    wzorzec,
    polecenieObserwacji,
    tlumienie,
    rekurencyjnie,
    zalozObserwacje,
    odczytajObserwacje,
  };
}
