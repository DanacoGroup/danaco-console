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

/** Interfejs OknoMonitora opisuje węzeł okna Process Monitor, które pokazuje na żywo rejestr procesów rdzenia terminala wraz z czynnościami udostępnianymi palecie poleceń. */
export interface OknoMonitora {
  element: HTMLElement;
  odswiez(): void;
  /** Czynności okna oddane palecie poleceń wywołują odświeżenie rejestru i eksport migawki wykazu. */
  czynnosci: readonly CzynnoscOkna[];
}

export function utworzOknoMonitora(
  zrodlo: ZrodloTerminala,
  stan: StanTerminala,
  // Parametr zostaje w podpisie, bo składa go wspólne złożenie modułu wraz z pozostałymi oknami.
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
  // Podglądy wyjścia żyją w oknie, aby przetrwały ponowne rysowanie wykazu procesów.
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
          // Karta procesu bywa nieznana temu widokowi, więc skok na nią musi się nazwać, a nie milczeć.
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
  // Zmiana filtra jest odczytem, bo filtr jedzie do rdzenia parametrem terminal.process.list.
  filtrStanu.naZmiane(odczytaj);
  filtrInicjatora.naZmiane(odczytaj);
  // Grupowanie jest porządkiem wyświetlania, więc zmiana tylko przerysowuje wykaz.
  grupowanie.naZmiane(pokaz);
  migawka.addEventListener('click', () => {
    const procesy = widoczne();
    if (procesy.length === 0) {
      // Migawka pustego wykazu byłaby nie do odróżnienia od migawki, której zapis się nie udał.
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

/** Żądanie wykazu procesów pomija pole filtra, którego wartość jest pusta, ponieważ pusty filtr oznacza brak zawężenia wykazu po tej cesze. */
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

/** Zawężenie i porządek wykazu procesów: przypięte idą na górę wykazu, reszta zachowuje kolejność nadaną przez stan procesu w rdzeniu. */
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

/** Podpis grupy wykazu procesów powstaje z samej nazwy grupy, ponieważ grupowanie jest wyłącznie porządkiem wyświetlania w oknie. */
function podpisGrupyProcesow(grupa: string): HTMLElement {
  const podpis = document.createElement('h3');
  podpis.className = 'dt-grupa';
  podpis.textContent = grupa;
  return podpis;
}

/** Zakończenie procesu w rdzeniu wykonuje się źródłem terminala i stanem treści przekazanymi funkcji jako parametry wywołania. */
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
    // Czasownik bierze się z odpowiedzi rdzenia, nie z samego faktu wysłania żądania.
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

/** Wstrzymanie albo wznowienie procesu w rdzeniu zależy od pola supported odpowiedzi, które odróżnia zmianę biegu od procesu pozostawionego nietkniętym. */
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

/** Stała CZEKANIE_NA_DOMKNIECIE_MS podaje czas w milisekundach, jaki rdzeń czeka na domknięcie procesu, zanim odczyta jego wyjście. */
const CZEKANIE_NA_DOMKNIECIE_MS = 3000;

/** Przełączenie podglądu wyjścia procesu odczytuje terminal.output.read albo zwija treść już pokazaną w wykazie procesów okna. */
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
        // Odmowa rdzenia idzie dosłownie, ponieważ sama nazywa dokładny powód braku wyjścia procesu.
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

/** Zdanie potwierdzenia odczytu wyjścia mierzy treść odpowiedzi rdzenia, a nie sam fakt jej otrzymania od serwera. */
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

/** Powtórzenie polecenia procesu wykonuje się w jego karcie źródłowej wraz ze źródłem terminala i stanem treści przekazanymi parametrem. */
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

/** Kontrolki okna Process Monitor obejmują filtry, grupowanie wykazu oraz przyciski odczytu i eksportu migawki procesów. */
interface PowierzchniaRejestruProcesow {
  filtrStanu: WyborDrzewem;
  filtrInicjatora: WyborDrzewem;
  grupowanie: WyborDrzewem;
  odczyt: HTMLButtonElement;
  migawka: HTMLButtonElement;
}

/** Złożenie powierzchni rejestru procesów tworzy filtry, grupowanie i przyciski okna bez domykania się na stanie okna ani na rdzeniu. */
function zlozPowierzchnieRejestruProcesow(
  rama: { akcje: HTMLElement; narzedzia: HTMLElement; cialo: HTMLElement },
  stanTresci: HTMLElement,
): PowierzchniaRejestruProcesow {
  const filtrStanu = utworzWyborDrzewem({ nastawa: 'Filtr stanu procesu', pozycje: FILTRY });
  const filtrInicjatora = utworzWyborDrzewem({ nastawa: 'Filtr inicjatora', pozycje: INICJATORZY });
  const grupowanie = utworzWyborDrzewem({ nastawa: 'Grupowanie wykazu', pozycje: GRUPOWANIA });
  const odczyt = przyciskAkcji('Odśwież wykaz');
  const migawka = przyciskAkcji('Eksportuj migawkę');

  // Wstrzymanie i wznowienie stoją przy procesie, nie w panelu, bo dotyczą jednego wiersza wykazu.
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
