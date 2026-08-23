import {
  Command,
  TerminalProcessStatus,
  type ProcessInitiator,
  type TerminalOutputReadResponse,
  type TerminalProcess,
  type TerminalProcessListRequest,
} from '../../../../shared/contract';
import { utworzRameOkna } from '../../komponenty/rama-okna';
import { pobierzPlik, przyciskAkcji, wykaz } from '../../modele/kontrolki-formularza';
import type { PokrycieKomend } from '../pokrycie-komend';
import { oznaczWarstwy, type CzynnoscOkna } from './czynnosci-okna';
import { FILTRY, GRUPOWANIA, INICJATORZY, pogrupuj } from './grupowanie-procesow';
import { podgladWyjscia } from './podglad-wyjscia';
import { pozycjaProcesu } from './pozycja-procesu';
import type { StanTerminala } from './stan-terminala';
import { utworzStanTresci, type StanTresci } from './stany-okna';
import { utworzWyborDrzewem, type WyborDrzewem } from './wybor-drzewem';
import type { ZrodloTerminala } from './zrodlo-terminala';

/**
 * Process Monitor — okno monitorujące modułu Terminal: podgląd i zakończenie
 * procesu, w tym procesu zainicjowanego poleceniem AI, zgodnie z rejestrem
 * procesów rdzenia serwera.
 *
 * Na żywo znaczy ze zdarzeń. Okno odczytuje `terminal.process.list` przy
 * wejściu i na wyraźne żądanie Operatora; każdą późniejszą zmianę przynosi
 * `terminal.process.changed`. Odpytywanie w pętli dałoby ten sam obraz drożej
 * i z opóźnieniem, a przy stu procesach zalałoby gniazdo.
 *
 * Filtr stanu i inicjatora jedzie do rdzenia parametrem komendy — tak stanowi
 * kontrakt i tak wynik jest spójny z dziennikiem rdzenia. Grupowanie jest
 * wyłącznie porządkiem wyświetlania i zostaje w oknie.
 *
 * Wyjście na żywo idzie `stream.chunk` do wspólnego bufora Output Console, ale
 * bufor żyje jedno połączenie: po rozłączeniu i ponownym podłączeniu ma zero
 * fragmentów, a rdzeń wciąż oddaje pełną treść. Osobny przycisk pozycji pyta
 * o nią `terminal.output.read` i pokazuje ją przy pozycji, a nie w buforze,
 * więc żaden wiersz nie wchodzi do konsoli dwa razy.
 */
export interface OknoMonitora {
  element: HTMLElement;
  odswiez(): void;
  /** Czynności okna oddane palecie poleceń. */
  czynnosci: readonly CzynnoscOkna[];
}

export function utworzOknoMonitora(
  zrodlo: ZrodloTerminala,
  stan: StanTerminala,
  // Okno nie ma dziś ani jednej pozycji bez pokrycia w rdzeniu: wstrzymanie
  // procesu było ostatnią i stoi już przy wierszu wykazu. Parametr zostaje
  // w podpisie, bo składa go wspólne złożenie modułu wraz z pozostałymi oknami.
  _pokrycie: PokrycieKomend,
): OknoMonitora {
  const rama = utworzRameOkna({
    tytul: 'Process Monitor',
    rola: 'monitor',
    przeznaczenie: 'Rejestr procesów rdzenia serwera: co biegnie, z czyjego polecenia i z jakim wynikiem.',
    kod: 'process-monitor',
    przedrostek: 'dt',
    ogniskowalne: true,
  });
  const tresc = utworzStanTresci();

  const { filtrStanu, filtrInicjatora, grupowanie, odczyt, migawka } =
    zlozPowierzchnieRejestruProcesow(rama, tresc.element);

  const przypiete = new Set<string>();
  // Podglądy wyjścia żyją w oknie, nie w węźle pozycji: wykaz przerysowuje się
  // przy każdym `terminal.process.changed`, więc treść trzymana w węźle znikałaby
  // Operatorowi pod ręką. Klucz to identyfikator procesu, wartość — odpowiedź
  // rdzenia w całości.
  const podglady = new Map<string, TerminalOutputReadResponse>();

  function widoczne(): TerminalProcess[] {
    return zawezWykazProcesow(stan.procesy(), filtrStanu.wartosc(), filtrInicjatora.wartosc(), przypiete);
  }

  function pokaz(): void {
    const procesy = widoczne();
    if (procesy.length === 0) {
      tresc.pusto('Rejestr rdzenia nie ma procesu spełniającego warunki wykazu.');
      return;
    }
    const miejsce = tresc.tresc();
    for (const [grupa, pozycje] of pogrupuj(procesy, grupowanie.wartosc())) {
      if (grupa !== '') miejsce.append(podpisGrupyProcesow(grupa));
      const lista = wykaz('Procesy rejestru rdzenia', 'dt-wykaz');
      for (const proces of pozycje) {
        const wezel = pozycjaProcesu(proces, {
          przypiety: przypiete.has(proces.id),
          podgladOtwarty: podglady.has(proces.id),
          zakoncz: (wymuszony) => zakonczProcesRejestru(zrodlo, stan, tresc, proces, wymuszony),
          wstrzymaj: (wznowienie) => wstrzymajProcesRejestru(zrodlo, stan, tresc, proces, wznowienie),
          uruchomPonownie: () => powtorzProcesRejestru(zrodlo, stan, tresc, proces),
          przypnij: () => {
            if (przypiete.has(proces.id)) przypiete.delete(proces.id);
            else przypiete.add(proces.id);
            pokaz();
          },
          eksportuj: () =>
            pobierzPlik(`proces-${proces.id}.json`, JSON.stringify(proces, null, 2), 'text/plain'),
          // Wykaz procesów przychodzi z rdzenia, a karty żyją w widoku: proces
          // bywa więc związany z kartą, której ten widok nie zna (zamkniętą
          // albo z innej sesji). Skok na taki identyfikator musi się nazwać,
          // zamiast gasić wykaz kart w oknie wiodącym bez słowa.
          doKarty: () => {
            const karta = proces.sessionId ?? '';
            if (karta === '') {
              tresc.potwierdzenie(
                `Proces ${proces.id} nie wskazuje karty źródłowej — nie ma dokąd skoczyć.`,
                false,
              );
              return;
            }
            if (!stan.ustawKarte(karta)) {
              tresc.potwierdzenie(
                `Karty ${karta} nie ma w tym widoku — rdzeń zna ją z rejestru, okno wiodące nie.`,
                false,
              );
              return;
            }
            tresc.potwierdzenie(`Okno wiodące przestawione na kartę ${karta}.`, true);
          },
          pokazWyjscie: () => przelaczPodgladWyjscia(zrodlo, tresc, proces, podglady, pokaz),
        });
        const odczytane = podglady.get(proces.id);
        if (odczytane !== undefined) wezel.append(podgladWyjscia(proces.id, odczytane));
        lista.append(wezel);
      }
      miejsce.append(lista);
    }
  }

  function odczytaj(): void {
    tresc.ladowanie('Odczyt rejestru procesów rdzenia…');
    const zadanie = zadanieWykazuProcesow(stan.okno(), filtrStanu.wartosc(), filtrInicjatora.wartosc());
    void zrodlo.procesy(zadanie).then((wynik) => {
      if (!wynik.udany || wynik.wynik === undefined) {
        tresc.blad(`Rdzeń nie oddał wykazu procesów (${Command.TerminalProcessList}).`, wynik.blad);
        return;
      }
      stan.ustawProcesy(wynik.wynik);
      tresc.potwierdzenie(`Rdzeń oddał ${wynik.wynik.length} procesów rejestru.`, true);
    });
  }

  odczyt.addEventListener('click', odczytaj);
  // Filtry jadą do rdzenia parametrem `terminal.process.list`, więc ich zmiana
  // jest odczytem; grupowanie jest porządkiem wyświetlania, więc wystarczy
  // przerysowanie.
  filtrStanu.naZmiane(odczytaj);
  filtrInicjatora.naZmiane(odczytaj);
  grupowanie.naZmiane(pokaz);
  migawka.addEventListener('click', () => {
    const procesy = widoczne();
    if (procesy.length === 0) {
      // Migawka pustego wykazu to plik `[]` — nie do odróżnienia od migawki,
      // której zapis się nie udał.
      tresc.potwierdzenie('Wykaz jest pusty — migawki nie zapisano.', false);
      return;
    }
    const nazwa = `migawka-procesow-${Date.now()}.json`;
    pobierzPlik(nazwa, JSON.stringify(procesy, null, 2), 'text/plain');
    tresc.potwierdzenie(`Zapisano ${nazwa} — ${procesy.length} procesów widocznych w wykazie.`, true);
  });

  stan.naZmiane(pokaz);

  const czynnosci: readonly CzynnoscOkna[] = [
    {
      okno: 'Process Monitor',
      nazwa: 'Odśwież rejestr procesów',
      opis: 'Odczytuje wykaz procesów rdzenia wraz z zawężeniem stanu i inicjatora.',
      warstwa: 'zawsze',
      wykonaj: odczytaj,
    },
    {
      okno: 'Process Monitor',
      nazwa: 'Eksportuj migawkę wykazu',
      opis: 'Zapisuje procesy widoczne w wykazie jako plik.',
      warstwa: 'kontekstowa',
      wykonaj: () => migawka.click(),
    },
  ];

  return { element: rama.element, odswiez: odczytaj, czynnosci };
}

/**
 * Pięć fragmentów wyjętych z wytwórni okna. Trzy pierwsze nie znają stanu okna
 * wcale — składają żądanie, zawężają wykaz i budują węzeł podpisu z samych
 * danych. Dwa ostatnie sięgają po rdzeń i po stan treści, więc biorą je
 * parametrem zamiast domykać się na wytwórni.
 */

/** Żądanie wykazu; puste pole filtra znaczy „bez zawężenia”, więc klucza nie ma wcale. */
function zadanieWykazuProcesow(
  okno: string,
  stanProcesu: string,
  inicjator: string,
): TerminalProcessListRequest {
  return {
    ...(okno === '' ? {} : { windowId: okno }),
    ...(stanProcesu === '' ? {} : { status: stanProcesu as TerminalProcessStatus }),
    ...(inicjator === '' ? {} : { initiator: inicjator as ProcessInitiator }),
  };
}

/** Zawężenie i porządek wykazu: procesy przypięte idą na górę, reszta zachowuje kolejność stanu. */
function zawezWykazProcesow(
  procesy: readonly TerminalProcess[],
  stanProcesu: string,
  inicjator: string,
  przypiete: ReadonlySet<string>,
): TerminalProcess[] {
  return procesy
    .filter((proces) => stanProcesu === '' || proces.status === stanProcesu)
    .filter((proces) => inicjator === '' || proces.initiator === inicjator)
    .sort((a, b) => Number(przypiete.has(b.id)) - Number(przypiete.has(a.id)));
}

/** Podpis grupy wykazu; grupowanie jest porządkiem wyświetlania, więc węzeł powstaje z samej nazwy. */
function podpisGrupyProcesow(grupa: string): HTMLElement {
  const podpis = document.createElement('h3');
  podpis.className = 'dt-grupa';
  podpis.textContent = grupa;
  return podpis;
}

/** Zakończenie procesu w rdzeniu; źródło i stan treści wchodzą parametrem. */
function zakonczProcesRejestru(
  zrodlo: ZrodloTerminala,
  stan: StanTerminala,
  tresc: StanTresci,
  proces: TerminalProcess,
  wymuszony: boolean,
): void {
  tresc.ladowanie(`Kończenie procesu ${proces.id}${wymuszony ? ' (sygnał wymuszony)' : ''}…`);
  void zrodlo.zakoncz({ processId: proces.id, force: wymuszony }).then((wynik) => {
    if (!wynik.udany || wynik.wynik === undefined) {
      tresc.blad(
        `Rdzeń nie zakończył procesu ${proces.id} (${Command.TerminalProcessKill}).`,
        wynik.blad,
      );
      return;
    }
    const zapisany = wynik.wynik;
    stan.zapiszProces(zapisany);
    // Czasownik bierze się z odpowiedzi, nie z żądania: samo wysłanie żądania
    // nie jest zakończeniem procesu.
    if (zapisany.status === TerminalProcessStatus.Running) {
      tresc.potwierdzenie(
        `Rdzeń przyjął zakończenie procesu ${zapisany.id}, ale oddaje go nadal w stanie ${zapisany.status} — proces biegnie.`,
        false,
      );
      return;
    }
    tresc.potwierdzenie(`Rdzeń zakończył proces ${zapisany.id} — stan ${zapisany.status}.`, true);
  });
}

/**
 * Wstrzymanie albo wznowienie procesu w rdzeniu.
 *
 * Pole `supported` odpowiedzi jest tu treścią, nie ozdobą: fałsz znaczy, że
 * proces został NIETKNIĘTY, bo system tego nie umie — co jest czymś innym niż
 * niepowodzenie czynności. Okno mówi to wprost, zamiast pokazywać powodzenie
 * przy procesie, który dalej zajmuje procesor.
 */
function wstrzymajProcesRejestru(
  zrodlo: ZrodloTerminala,
  stan: StanTerminala,
  tresc: StanTresci,
  proces: TerminalProcess,
  wznowienie: boolean,
): void {
  const czynnosc = wznowienie ? 'Wznawianie' : 'Wstrzymywanie';
  tresc.ladowanie(`${czynnosc} procesu ${proces.id}…`);
  void zrodlo
    .wstrzymaj({ processId: proces.id, resume: wznowienie })
    .then((wynik) => {
      if (!wynik.udany || wynik.wynik === undefined) {
        tresc.blad(
          `Rdzeń nie zmienił biegu procesu ${proces.id} (${Command.TerminalProcessSuspend}).`,
          wynik.blad,
        );
        return;
      }
      stan.zapiszProces(wynik.wynik.process);
      if (!wynik.wynik.supported) {
        tresc.potwierdzenie(
          `Proces ${proces.id} został NIETKNIĘTY: system maszyny rdzenia nie zna wstrzymania obcego ` +
            'drzewa procesów. To nie jest niepowodzenie czynności — to granica platformy.',
          false,
        );
        return;
      }
      tresc.potwierdzenie(
        wznowienie
          ? `Proces ${proces.id} wznowiony — praca podjęta w miejscu, w którym stanęła.`
          : `Proces ${proces.id} wstrzymany wraz z potomstwem; wykonana praca nie przepadła.`,
        true,
      );
    });
}

/**
 * Ile milisekund rdzeń ma czekać na domknięcie procesu, zanim odczyta wyjście.
 *
 * Kontrakt dopuszcza 60 000, ale czekanie trzyma zadanie gniazda i minuta bez
 * odpowiedzi wygląda dla Operatora jak zawieszone okno. Trzy sekundy
 * wystarczają, żeby polecenie krótkie zdążyło się domknąć i oddało komplet wraz
 * z kodem wyjścia; polecenie długie i tak odda wyjście dotychczasowe ze stanem
 * `running`, bo czekanie nie jest warunkiem odpowiedzi. Operator, który chce
 * zobaczyć resztę, klika ponownie.
 */
const CZEKANIE_NA_DOMKNIECIE_MS = 3000;

/**
 * Podgląd wyjścia procesu: odczyt `terminal.output.read` albo zwinięcie tego,
 * co już stoi. Wykaz podglądów i przerysowanie wchodzą parametrem, bo rysowanie
 * zostaje w wytwórni okna.
 *
 * Zwinięcie nie pyta rdzenia: drugie kliknięcie zdejmuje treść z widoku i mówi
 * to wprost, inaczej Operator nie odróżniłby zwinięcia od odczytu, który wrócił
 * pusty.
 */
function przelaczPodgladWyjscia(
  zrodlo: ZrodloTerminala,
  tresc: StanTresci,
  proces: TerminalProcess,
  podglady: Map<string, TerminalOutputReadResponse>,
  przerysuj: () => void,
): void {
  if (podglady.delete(proces.id)) {
    przerysuj();
    tresc.potwierdzenie(
      `Podgląd wyjścia procesu ${proces.id} zwinięty. Do rdzenia nic nie poszło — wyjście dalej tam jest.`,
      true,
    );
    return;
  }

  const czynny = proces.status === TerminalProcessStatus.Running;
  tresc.ladowanie(
    czynny
      ? `Odczyt wyjścia procesu ${proces.id} — rdzeń poczeka do ${CZEKANIE_NA_DOMKNIECIE_MS} ms na jego domknięcie…`
      : `Odczyt wyjścia procesu ${proces.id}…`,
  );
  void zrodlo
    .odczytajWyjscie({
      processId: proces.id,
      ...(czynny ? { waitMs: CZEKANIE_NA_DOMKNIECIE_MS } : {}),
    })
    .then((wynik) => {
      if (!wynik.udany || wynik.wynik === undefined) {
        // Odmowa rdzenia idzie dosłownie — `not_found` niesie pełne zdanie
        // o tym, czy proces nigdy nie ruszył, wypadł z historii, czy rdzeń był
        // uruchomiony ponownie. Parafraza zgubiłaby wszystkie trzy powody.
        tresc.blad(
          `Rdzeń nie oddał wyjścia procesu ${proces.id} (${Command.TerminalOutputRead}).`,
          wynik.blad,
        );
        return;
      }
      const odpowiedz = wynik.wynik;
      podglady.set(proces.id, odpowiedz);
      przerysuj();
      tresc.potwierdzenie(zdanieOdczytuWyjscia(proces, odpowiedz), true);
    });
}

/**
 * Zdanie potwierdzenia odczytu — mierzy treść, nie sam fakt odpowiedzi.
 *
 * Wyjście puste jest przebiegiem udanym (polecenie mogło nic nie wypisać), ale
 * różni się od wyjścia niepustego, inaczej „odczytano" znaczyłoby to samo
 * w obu przypadkach. Rozjazd stanu też idzie wprost: wykaz w oknie jest kopią
 * z `terminal.process.list`, a odpowiedź na odczyt przychodzi z tej chwili —
 * gdy się różnią, świeższa jest odpowiedź.
 */
function zdanieOdczytuWyjscia(
  proces: TerminalProcess,
  odpowiedz: TerminalOutputReadResponse,
): string {
  const bajty = odpowiedz.stdout.length + odpowiedz.stderr.length;
  const czesci = [
    bajty === 0
      ? `Rdzeń oddał odczyt procesu ${proces.id}, ale oba strumienie są puste — to polecenie nic nie wypisało.`
      : `Rdzeń oddał wyjście procesu ${proces.id}: ${odpowiedz.stdout.length} znaków stdout, ${odpowiedz.stderr.length} znaków stderr.`,
  ];
  if (odpowiedz.status !== proces.status) {
    czesci.push(
      `Stan wg odczytu to ${odpowiedz.status}, a wykaz okna nosi ${proces.status} — świeższy jest odczyt.`,
    );
  }
  return czesci.join(' ');
}

/** Powtórzenie polecenia procesu w jego karcie źródłowej; źródło i stan treści wchodzą parametrem. */
function powtorzProcesRejestru(
  zrodlo: ZrodloTerminala,
  stan: StanTerminala,
  tresc: StanTresci,
  proces: TerminalProcess,
): void {
  const karta = proces.sessionId ?? '';
  if (karta === '') {
    tresc.blad('Proces nie wskazuje karty źródłowej — nie ma go gdzie powtórzyć.');
    return;
  }
  tresc.ladowanie(`Ponowne uruchomienie „${proces.command}”…`);
  void zrodlo
    .wykonaj({ sessionId: karta, command: proces.command, initiator: proces.initiator })
    .then((wynik) => {
      if (!wynik.udany || wynik.wynik === undefined) {
        tresc.blad(`Rdzeń nie powtórzył polecenia (${Command.TerminalCommandExec}).`, wynik.blad);
        return;
      }
      stan.zapiszProces(wynik.wynik);
      tresc.potwierdzenie(`Proces ${wynik.wynik.id} uruchomiony ponownie.`, true);
    });
}

/** Kontrolki okna Process Monitor. */
interface PowierzchniaRejestruProcesow {
  filtrStanu: WyborDrzewem;
  filtrInicjatora: WyborDrzewem;
  grupowanie: WyborDrzewem;
  odczyt: HTMLButtonElement;
  migawka: HTMLButtonElement;
}

/**
 * Składa kontrolki, pasek akcji, pasek narzędzi i ciało okna.
 *
 * Czysta konstrukcja: nie domyka się na stanie okna ani na rdzeniu, więc dała
 * się wyjąć bez przenoszenia zależności. Pozycja, której okno nie wykonuje, jest
 * jawnie nieczynna wraz z powodem liczonym z odczytu wykazu komend rdzenia.
 */
function zlozPowierzchnieRejestruProcesow(
  rama: { akcje: HTMLElement; narzedzia: HTMLElement; cialo: HTMLElement },
  stanTresci: HTMLElement,
): PowierzchniaRejestruProcesow {
  const filtrStanu = utworzWyborDrzewem({ nastawa: 'Filtr stanu procesu', pozycje: FILTRY });
  const filtrInicjatora = utworzWyborDrzewem({ nastawa: 'Filtr inicjatora', pozycje: INICJATORZY });
  const grupowanie = utworzWyborDrzewem({ nastawa: 'Grupowanie wykazu', pozycje: GRUPOWANIA });
  const odczyt = przyciskAkcji('Odśwież wykaz');
  const migawka = przyciskAkcji('Eksportuj migawkę');

  // Wstrzymanie i wznowienie stoją PRZY PROCESIE, a nie w panelu okna: dotyczą
  // jednego wiersza wykazu, a przycisk panelu musiałby najpierw zapytać, którego.
  oznaczWarstwy([
    [odczyt, 'zawsze'],
    [filtrStanu.element, 'na-zadanie'],
    [filtrInicjatora.element, 'na-zadanie'],
    [grupowanie.element, 'na-zadanie'],
    [migawka, 'kontekstowa'],
  ]);

  rama.akcje.append(odczyt, migawka);
  rama.narzedzia.append(filtrStanu.element, filtrInicjatora.element, grupowanie.element);
  rama.cialo.append(stanTresci);

  return { filtrStanu, filtrInicjatora, grupowanie, odczyt, migawka };
}
