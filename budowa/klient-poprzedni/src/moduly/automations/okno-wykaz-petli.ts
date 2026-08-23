import type { AutomationWorkflow } from '../../../../shared/contract';
import { utworzRameOkna } from '../../komponenty/rama-okna';
import {
  pole,
  poleLiczbowe,
  przelacznikWidoku,
  przestaw,
  przyciskAkcji as przycisk,
  wiersz,
} from '../../modele/kontrolki-formularza';
import {
  odczytajPrzekazania,
  zdanieOPrzekazaniach,
  type PrzekazanyScenariusz,
} from './przyjecie-przekazania';
import type { StanAutomatyki } from './stan-automatyki';
import { utworzStanTresci } from './stany-okna';
import { uruchomPetle, zdanieUruchomienia } from './uruchomienie-petli';
import { listaPrzekazan } from './widok-przekazan';
import { widokWykazuPetli } from './wykaz-petli-widok';
import type { ZrodloAutomations } from './zrodlo-automations';

/**
 * Wykaz gotowych pętli — okno wiodące modułu Automations. Bierze pętlę zapisaną
 * wcześniej i puszcza ją w ruch jednym kliknięciem: przycisk „Uruchom” nie pyta
 * o potwierdzenie ani o wskazanie kolejki — zakłada kolejkę i posuwa ją
 * działaniem `start` (`uruchomienie-petli.ts`). Pętla wyłączona ma ten sam
 * przycisk co czynna; jej stan stoi w opisie pozycji.
 *
 * Zawężanie zamiast blokowania wierszy: napis szukania skraca wykaz po stronie
 * klienta (kontrakt nie ma pola zapytania), przełącznik „tylko czynne” zawęża
 * samo żądanie polem `enabledOnly`.
 *
 * Uruchomienie przestawia wspólny stan modułu (`stan-automatyki.ts`) na tę pętlę
 * i jej kolejkę, więc Execution Monitor pokazuje jej przebiegi, a Queue Manager
 * posuwa jej kolejkę bez przepisywania identyfikatorów.
 *
 * Nad wykazem stoi wykaz przekazań: scenariusze przeniesione tu z innych modułów
 * komendą `context.transfer` (`przyjecie-przekazania.ts`). Okno jest wejściem do
 * modułu, więc to tutaj widać, że ktoś coś oddał. Przekazanie zostaje propozycją
 * do chwili, w której Operator naciśnie zapis — nic nie wchodzi do magazynu
 * automatyk samo.
 */
export interface OknoWykazuPetli {
  element: HTMLElement;
  odswiez(): void;
  /** Zamyka nasłuch zmiany powiązania automatyki. */
  zamknij(): void;
}

export function utworzOknoWykazuPetli(
  zrodlo: ZrodloAutomations,
  stan: StanAutomatyki,
): OknoWykazuPetli {
  const rama = utworzRameOkna({
    tytul: 'Wykaz gotowych pętli',
    rola: 'wiodące',
    przeznaczenie:
      'Pętle zapisane wcześniej, uruchamiane jednym kliknięciem — bez budowania definicji za każdym razem.',
    przedrostek: 'da',
  });
  const tresc = utworzStanTresci();
  const powierzchnia = zlozPowierzchnieWykazu(rama, tresc.element);
  const { szukanie, granica, tylkoCzynne, odswiezPrzycisk, przekazaniaPrzycisk } = powierzchnia;

  /** Ostatni wykaz oddany przez rdzeń — po nim idzie zawężanie bez odpytywania. */
  let petle: readonly AutomationWorkflow[] = [];
  /** Scenariusze przeniesione tu z innych modułów, czekające na decyzję. */
  let przekazania: readonly PrzekazanyScenariusz[] = [];

  /**
   * Rysuje oba wykazy w jednym miejscu treści.
   *
   * Miejsce treści czyści się przy każdym wskazaniu (`stan-tresci.ts`), więc
   * dwa osobne rysunki wypchnęłyby się nawzajem. Pustkę ogłaszamy dopiero
   * wtedy, gdy nie ma ani pętli, ani przekazań — samo przekazanie oczekujące
   * jest już treścią.
   */
  function rysuj(): void {
    const widok = widokWykazuPetli(petle, szukanie.value, uruchom);
    const wykazPrzekazan = listaPrzekazan(przekazania, zapiszPrzekazanie);
    if (widok.rodzaj === 'pustka' && wykazPrzekazan === null) {
      tresc.pusto(widok.zdanie);
      return;
    }
    const miejsce = tresc.tresc();
    if (wykazPrzekazan !== null) miejsce.append(wykazPrzekazan);
    if (widok.rodzaj === 'wykaz') miejsce.append(widok.element);
    else miejsce.append(zdanieOBrakuPetli(widok.zdanie));
  }

  /**
   * Odczyt wykazu. Zdanie potwierdzenia podaje liczbę oddanych pozycji, żeby
   * wykaz krótki dał się odróżnić od przyciętego granicą.
   *
   * Oddaje obietnicę, bo po odczycie wykazu idzie odczyt przekazań i oba piszą
   * w to samo miejsce treści — bez ustawienia ich w kolejkę zdanie późniejszej
   * odpowiedzi nadpisywałoby zdanie wcześniejszej w porządku przypadkowym.
   */
  function odczytaj(): Promise<void> {
    tresc.ladowanie('Odczyt wykazu gotowych pętli…');
    return zrodlo.automatyki(zadanieWykazu(tylkoCzynne, granica.value)).then((wynik) => {
      if (!wynik.udany || wynik.wynik === undefined) {
        tresc.blad(
          'Rdzeń nie oddał wykazu gotowych pętli, więc nie ma czego uruchomić. ' +
            'Wykaz bierze się z komendy automation.workflow.list i to ona nie przeszła. ' +
            'Naciśnij „Odśwież wykaz” ponownie.',
          wynik.blad,
        );
        return;
      }
      petle = wynik.wynik;
      rysuj();
      tresc.potwierdzenie(zdanieOdczytu(petle.length, tylkoCzynne), true);
    });
  }

  /**
   * Uruchomienie pozycji. Wykaz zostaje na ekranie — rysunek idzie ponownie
   * z pamięci okna, więc druga pętla rusza bez powtórnego odczytu.
   */
  function uruchom(petla: AutomationWorkflow): void {
    tresc.ladowanie(`Uruchamianie pętli „${petla.name}”…`);
    void uruchomPetle(zrodlo, petla).then((skutek) => {
      if (!skutek.udane) {
        rysuj();
        tresc.blad(skutek.zdanie, skutek.powod);
        return;
      }
      stan.ustawAutomatyke(petla.id, petla);
      stan.ustawKolejke(skutek.kolejka.id);
      rysuj();
      tresc.potwierdzenie(zdanieUruchomienia(petla, skutek.kolejka), true);
    });
  }

  /**
   * Odczyt przekazań oczekujących.
   *
   * Wykaz pętli zostaje na ekranie: odczyt przekazań dokłada sekcję, a nie
   * zastępuje treści. Zdanie mówi, ile okien sprawdzono i ile z nich niesie
   * scenariusz — pustka po sprawdzeniu jest odpowiedzią, nie usterką.
   */
  function odczytajPrzekazaniaOkna(): Promise<void> {
    tresc.ladowanie('Odczyt przekazań z innych modułów…');
    return odczytajPrzekazania(zrodlo).then((skutek) => {
      if (!skutek.odczytane) {
        przekazania = [];
        rysuj();
        tresc.blad(skutek.zdanie);
        return;
      }
      przekazania = skutek.scenariusze;
      rysuj();
      tresc.potwierdzenie(zdanieOPrzekazaniach(skutek), true);
    });
  }

  /**
   * Zapis przekazanego scenariusza do magazynu automatyk.
   *
   * To jedyne miejsce, w którym przekazanie staje się automatyką, i dzieje się
   * to na jawne naciśnięcie. Po zapisie moduł przestawia się na tę automatykę,
   * więc Workflow Builder pokazuje jej kroki od razu — Operator ogląda, co
   * przyszło, dopiero gdy to przyjął.
   */
  function zapiszPrzekazanie(scenariusz: PrzekazanyScenariusz): void {
    tresc.ladowanie(`Zapis przekazanego scenariusza „${scenariusz.nazwa}”…`);
    const zadanie: Parameters<ZrodloAutomations['zapiszAutomatyke']>[0] = {
      name: scenariusz.nazwa,
      steps: scenariusz.kroki,
      enabled: scenariusz.czynny,
    };
    if (scenariusz.opis !== '') zadanie.description = scenariusz.opis;
    void zrodlo.zapiszAutomatyke(zadanie).then((wynik) => {
      if (!wynik.udany || wynik.wynik === undefined) {
        rysuj();
        tresc.blad(
          `Rdzeń nie zapisał przekazanego scenariusza „${scenariusz.nazwa}”. ` +
            'Scenariusz zostaje przekazaniem oczekującym — naciśnij zapis ponownie.',
          wynik.blad,
        );
        return;
      }
      const automatyka = wynik.wynik;
      przekazania = przekazania.filter((pozycja) => pozycja !== scenariusz);
      stan.ustawDefinicje(automatyka);
      void odczytaj().then(() => {
        tresc.potwierdzenie(
          `Rdzeń zapisał przekazany scenariusz jako automatykę ${automatyka.id} ` +
            `(wersja ${automatyka.version ?? 1}, kroków: ${(automatyka.steps ?? []).length}). ` +
            'Jej kroki stoją w Workflow Builderze.',
          true,
        );
      });
    });
  }

  /** Pełny odczyt okna: najpierw magazyn automatyk, potem to, co przekazano. */
  function odswiez(): void {
    void odczytaj().then(() => odczytajPrzekazaniaOkna());
  }

  odswiezPrzycisk.addEventListener('click', () => void odczytaj());
  przekazaniaPrzycisk.addEventListener('click', () => void odczytajPrzekazaniaOkna());
  szukanie.addEventListener('input', rysuj);
  // Przełącznik zawęża żądanie, nie rysunek — `enabledOnly` rozstrzyga rdzeń,
  // więc przestawienie musi pójść po wykaz jeszcze raz.
  tylkoCzynne.addEventListener('click', () => {
    przestaw(tylkoCzynne);
    void odczytaj();
  });

  // Powiązanie automatyki zmienia się także z pracy innego modułu — tak wraca
  // scenariusz wpięty tu z Browsera czy z Assistanta. Wykaz w oknie jest wtedy
  // nieaktualny, więc czyta się ponownie wraz z przekazaniami.
  const odsubskrybuj = zrodlo.naZmianePowiazania(odswiez);

  return { element: rama.element, odswiez, zamknij: odsubskrybuj };
}

/** Zdanie o braku pętli osadzane obok wykazu przekazań, gdy ten nie jest pusty. */
function zdanieOBrakuPetli(zdanie: string): HTMLElement {
  const akapit = document.createElement('p');
  akapit.className = 'dn-pole-opis';
  akapit.textContent = zdanie;
  return akapit;
}

/** Kontrolki okna: szukanie, granica wykazu, przełącznik zawężenia i odczyty. */
interface PowierzchniaWykazu {
  szukanie: HTMLInputElement;
  granica: HTMLInputElement;
  tylkoCzynne: HTMLButtonElement;
  odswiezPrzycisk: HTMLButtonElement;
  przekazaniaPrzycisk: HTMLButtonElement;
}

/**
 * Składa kontrolki, pasek akcji i ciało okna. Konstrukcja czysta — nie domyka
 * się ani na stanie okna, ani na źródle.
 */
function zlozPowierzchnieWykazu(
  rama: { akcje: HTMLElement; narzedzia: HTMLElement; cialo: HTMLElement },
  stanTresci: HTMLElement,
): PowierzchniaWykazu {
  const szukanie = pole('Szukaj pętli po nazwie', 'fragment nazwy, identyfikatora albo opisu');
  szukanie.type = 'search';
  const granica = poleLiczbowe('Górna granica liczby pętli', 'domyślnie wszystkie');
  const tylkoCzynne = przelacznikWidoku('Tylko czynne', false);
  const odswiezPrzycisk = przycisk('Odśwież wykaz', 'dn-btn dn-btn--atrament');
  // Przekazania idą dwiema komendami kontraktu: `window.list` po okna tego
  // modułu i `config.effective.get` po przeniesiony komplet każdego z nich.
  const przekazaniaPrzycisk = przycisk('Odczytaj przekazania');

  rama.akcje.append(odswiezPrzycisk, przekazaniaPrzycisk);
  rama.narzedzia.append(tylkoCzynne);
  rama.cialo.append(
    wiersz('Szukanie', szukanie, {
      klasa: 'da-wiersz',
      objasnienie:
        'Zawęża wykaz już odczytany — komenda wykazu nie ma pola zapytania, więc szukanie idzie po stronie klienta.',
    }),
    wiersz('Granica wykazu', granica, {
      klasa: 'da-wiersz',
      objasnienie: 'Puste pole zwraca wszystkie pętle znane rdzeniowi.',
    }),
    stanTresci,
  );

  return { szukanie, granica, tylkoCzynne, odswiezPrzycisk, przekazaniaPrzycisk };
}

/**
 * Żądanie wykazu — czysta konstrukcja z wartości pól. Granica pusta albo
 * niepoprawna znaczy „wszystkie”, więc nie trafia do żądania wcale; tak samo
 * `enabledOnly` idzie tylko wtedy, gdy przełącznik jest wciśnięty.
 */
function zadanieWykazu(
  tylkoCzynne: HTMLButtonElement,
  granica: string,
): Parameters<ZrodloAutomations['automatyki']>[0] {
  const zadanie: Parameters<ZrodloAutomations['automatyki']>[0] = {};
  if (tylkoCzynne.dataset['wlaczony'] === 'true') zadanie.enabledOnly = true;
  const limit = Number.parseInt(granica, 10);
  if (Number.isInteger(limit) && limit > 0) zadanie.limit = limit;
  return zadanie;
}

/** Zdanie po odczycie: liczba pozycji oddanych przez rdzeń i zakres żądania. */
function zdanieOdczytu(ile: number, tylkoCzynne: HTMLButtonElement): string {
  const zakres = tylkoCzynne.dataset['wlaczony'] === 'true' ? ' czynnych' : '';
  return ile === 0
    ? `Rdzeń oddał pusty wykaz${zakres === '' ? '' : ' pętli czynnych'}.`
    : `Rdzeń oddał ${ile}${zakres} pętli. Uruchom dowolną jednym kliknięciem.`;
}
