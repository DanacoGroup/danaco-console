import {
  Command,
  TerminalProcessStatus,
  TerminalShell,
  type TerminalScriptLintResponse,
  type TerminalSession,
} from '../../../../shared/contract';
import { utworzRameOkna } from '../../komponenty/rama-okna';
import {
  pobierzPlik,
  pole,
  poleTresci,
  pozycjaWykazu,
  przyciskAkcji,
  wykaz,
} from '../../modele/kontrolki-formularza';
import {
  czytajParametry,
  kontrolaWstepna,
  podstawParametry,
  polecenieSprawdzeniaSkladni,
  utworzBiblioteke,
  type Biblioteka,
  type ParametrSkryptu,
  type PozycjaBiblioteki,
  type RodzajPozycji,
} from './biblioteka-skryptow';
import type { PokrycieKomend } from '../pokrycie-komend';
import { oznaczWarstwy, type CzynnoscOkna } from './czynnosci-okna';
import { POWLOKI } from './profil-karty';
import type { StanTerminala } from './stan-terminala';
import { utworzStanTresci } from './stany-okna';
import { utworzWyborDrzewem, type PozycjaWyboru, type WyborDrzewem } from './wybor-drzewem';
import type { ZrodloTerminala } from './zrodlo-terminala';
import { programPowloki, zdanieProgramuPowloki } from './zaleznosci-zewnetrzne';

/**
 * Script Library — okno kreatora modułu Terminal: skrypty i snippety, ich
 * parametry, kontrola wstępna treści i uruchomienie w karcie bieżącej.
 *
 * Do rdzenia idą dwie komendy: `terminal.command.exec` uruchamia treść skryptu,
 * a `terminal.output.read` oddaje wynik sprawdzenia składni. Rdzeń podaje treść
 * programowi powłoki jako pojedynczy argument, więc skrypt wieloliniowy
 * wykonuje się bez zapisywania go do pliku — i dlatego biblioteka nie potrzebuje
 * ani komendy zapisu pliku, ani ścieżki na dysku serwera.
 *
 * Biblioteka jest bytem RDZENIA, a okno jej widokiem: zapis idzie komendą
 * `terminal.script.save` — każdy zapis zakłada kolejną WERSJĘ, więc poprzednie
 * brzmienie treści zostaje — wykaz czyta `terminal.script.list`, a usunięcie
 * `terminal.script.remove` zabiera pozycję wraz ze wszystkimi jej wersjami.
 *
 * Analizę treści prowadzi `terminal.script.lint` programami leżącymi na maszynie
 * rdzenia. Odpowiedź niesie pole mówiące, CZY program analizy był dostępny —
 * i okno je pokazuje, bo pusty wykaz uwag przy braku programu znaczyłby
 * fałszywie „treść bez zastrzeżeń”. Kontrola wstępna okna zostaje obok analizy,
 * nie zamiast niej: sprawdza to, co da się rozstrzygnąć bez żadnego programu,
 * i robi to w chwili pisania.
 */
export interface OknoBibliotekiSkryptow {
  element: HTMLElement;
  odswiez(): void;
  czynnosci: readonly CzynnoscOkna[];
}

/** Rodzaje pozycji biblioteki w kolejności wykazu. */
const RODZAJE: readonly PozycjaWyboru[] = [
  ['skrypt', 'Skrypt', 'Treść wieloliniowa uruchamiana w całości jako jeden proces powłoki.'],
  ['snippet', 'Snippet', 'Krótkie polecenie do powtarzania; uruchamia się tą samą drogą co skrypt.'],
];

export function utworzOknoBibliotekiSkryptow(
  zrodlo: ZrodloTerminala,
  stan: StanTerminala,
  pokrycie: PokrycieKomend,
): OknoBibliotekiSkryptow {
  const rama = utworzRameOkna({
    tytul: 'Script Library',
    rola: 'kreator',
    przeznaczenie:
      'Skrypty i snippety wraz z parametrami. Treść uruchamia się w karcie bieżącej jako jeden proces powłoki.',
    kod: 'script-library',
    przedrostek: 'dt',
    ogniskowalne: true,
  });
  const tresc = utworzStanTresci();
  const biblioteka = utworzBiblioteke();
  const kontrolki = zlozPowierzchnieBiblioteki(rama, tresc.element, pokrycie);
  /** Wartości parametrów wpisane w formularzu, po nazwie parametru. */
  const wartosciParametrow = new Map<string, string>();
  /**
   * Treść po formatowaniu, przywieziona z ostatniej analizy. Nie wchodzi do
   * edytora sama: zamiana treści pisanej przez Operatora bez jego wskazania
   * byłaby zabraniem mu wersji, którą właśnie pisał.
   */
  let sformatowana = '';

  function pokaz(): void {
    const miejsce = tresc.tresc();
    const parametry = czytajParametry(kontrolki.edytor.value);
    miejsce.append(
      formularzParametrow(parametry, wartosciParametrow),
      wykazUwag(kontrolaWstepna(kontrolki.edytor.value)),
      wykazPozycji(biblioteka, {
        wczytaj: (pozycja) => wczytajDoEdytora(kontrolki, wartosciParametrow, pozycja, pokaz),
        uruchom: (pozycja) => uruchomTresc(pozycja),
        usun: (pozycja) => usunPozycje(pozycja),
      }),
    );
  }

  /** Czyta bibliotekę z rdzenia; wykaz na ekranie jest jej odwzorowaniem. */
  function odczytaj(): void {
    void zrodlo.skrypty({}).then((wynik) => {
      if (!wynik.udany || wynik.wynik === undefined) {
        tresc.blad(`Rdzeń nie oddał biblioteki skryptów (${Command.TerminalScriptList}).`, wynik.blad);
        return;
      }
      biblioteka.zastap(
        wynik.wynik.map((pozycja) => ({
          id: pozycja.id,
          nazwa: pozycja.name,
          rodzaj: (pozycja.kind === 'snippet' ? 'snippet' : 'skrypt') as RodzajPozycji,
          powloka: pozycja.shell,
          tresc: pozycja.content,
          tagi: pozycja.tags ?? '',
          wersja: pozycja.version,
          ostatnieUruchomienie: pozycja.lastRunAt ?? 0,
        })),
      );
      pokaz();
    });
  }

  /** Usuwa pozycję z biblioteki RDZENIA wraz ze wszystkimi jej wersjami. */
  function usunPozycje(pozycja: PozycjaBiblioteki): void {
    if (pozycja.id === undefined || pozycja.id === '') {
      biblioteka.usun(pozycja.nazwa);
      pokaz();
      tresc.potwierdzenie(
        `Pozycja ${pozycja.nazwa} zniesiona z wykazu. W rdzeniu jej nie było — nie została zapisana.`,
        true,
      );
      return;
    }
    void zrodlo.usunSkrypt({ scriptId: pozycja.id }).then((wynik) => {
      if (!wynik.udany) {
        tresc.blad(`Rdzeń nie usunął pozycji ${pozycja.nazwa} (${Command.TerminalScriptRemove}).`, wynik.blad);
        return;
      }
      odczytaj();
      tresc.potwierdzenie(
        wynik.wynik === true
          ? `Pozycja ${pozycja.nazwa} usunięta z biblioteki rdzenia wraz ze wszystkimi jej wersjami.`
          : `Pozycji ${pozycja.nazwa} w bibliotece rdzenia nie było — wykaz odświeżony.`,
        true,
      );
    });
  }

  /** Poddaje treść edytora analizie statycznej na maszynie rdzenia. */
  function analizuj(): void {
    const zawartosc = kontrolki.edytor.value;
    if (zawartosc.trim() === '') {
      tresc.potwierdzenie('Pusta treść nie ma czego analizować.', false);
      return;
    }
    const powloka = kontrolki.powloka.wartosc() as TerminalShell;
    tresc.ladowanie(`Analiza treści dla powłoki ${powloka} na maszynie rdzenia…`);
    void zrodlo.sprawdzSkrypt({ content: zawartosc, shell: powloka, format: true }).then((wynik) => {
      if (!wynik.udany || wynik.wynik === undefined) {
        tresc.blad(`Rdzeń nie przeprowadził analizy (${Command.TerminalScriptLint}).`, wynik.blad);
        return;
      }
      const odpowiedz = wynik.wynik;
      const miejsce = tresc.tresc();
      miejsce.append(
        formularzParametrow(czytajParametry(zawartosc), wartosciParametrow),
        wykazUwagAnalizy(odpowiedz),
        wykazPozycji(biblioteka, {
          wczytaj: (pozycja) => wczytajDoEdytora(kontrolki, wartosciParametrow, pozycja, pokaz),
          uruchom: (pozycja) => uruchomTresc(pozycja),
          usun: (pozycja) => usunPozycje(pozycja),
        }),
      );
      if (!odpowiedz.analyzerAvailable) {
        // Brak programu NIE jest tym samym co brak zastrzeżeń i okno tego nie
        // skleja: pusty wykaz uwag pokazany jako „treść czysta” byłby nieprawdą.
        tresc.potwierdzenie(
          `Treści nie sprawdzono: na maszynie rdzenia nie ma programu analizy (${odpowiedz.analyzer}).`,
          false,
        );
        return;
      }
      if (odpowiedz.formatted !== undefined && odpowiedz.formatted !== '') {
        sformatowana = odpowiedz.formatted;
      }
      tresc.potwierdzenie(
        odpowiedz.findings.length === 0
          ? `Analiza (${odpowiedz.analyzer}) nie zgłosiła zastrzeżeń.`
          : `Analiza (${odpowiedz.analyzer}) zgłosiła ${odpowiedz.findings.length} uwag. ` +
              (odpowiedz.formatted === undefined
                ? ''
                : 'Treść sformatowana czeka pod pozycją „Wstaw treść sformatowaną”.'),
        odpowiedz.findings.length === 0,
      );
    });
  }

  /** Pozycja złożona z pól formularza; pusta nazwa albo pusta treść daje brak pozycji. */
  function pozycjaZFormularza(): PozycjaBiblioteki | null {
    const nazwa = kontrolki.nazwa.value.trim();
    const zawartosc = kontrolki.edytor.value;
    if (nazwa === '' || zawartosc.trim() === '') return null;
    const poprzednia = biblioteka.znajdz(nazwa);
    return {
      ...(poprzednia?.id === undefined ? {} : { id: poprzednia.id }),
      nazwa,
      rodzaj: kontrolki.rodzaj.wartosc() as RodzajPozycji,
      powloka: kontrolki.powloka.wartosc() as TerminalShell,
      tresc: zawartosc,
      tagi: kontrolki.tagi.value.trim(),
      wersja: (poprzednia?.wersja ?? 0) + 1,
      ostatnieUruchomienie: poprzednia?.ostatnieUruchomienie ?? 0,
    };
  }

  /** Karta, w której okno uruchamia treść; brak karty ma własne zdanie. */
  function kartaWykonania(czynnosc: string): TerminalSession | null {
    const karta = stan.kartaBiezaca();
    if (karta === null) {
      tresc.blad(`Nie ma karty bieżącej — ${czynnosc} nie ma w czym się wykonać. Otwórz kartę w oknie Terminal Tabs.`);
    }
    return karta;
  }

  function uruchomTresc(pozycja: PozycjaBiblioteki): void {
    const karta = kartaWykonania(`uruchomienie pozycji ${pozycja.nazwa}`);
    if (karta === null) return;
    const podstawienie = podstawParametry(pozycja.tresc, wartosciParametrow);
    if (podstawienie.brakujace.length > 0) {
      tresc.blad(
        `Pozycja ${pozycja.nazwa} ma parametry bez wartości: ${podstawienie.brakujace.join(', ')}. ` +
          'Odwołania zostały w treści nietknięte — wypełnij formularz parametrów i uruchom ponownie.',
      );
      return;
    }
    const rozjazd =
      karta.shell === pozycja.powloka
        ? ''
        : `Uwaga: pozycję napisano dla powłoki ${pozycja.powloka}, a karta bieżąca jest powłoką ${karta.shell}. ` +
          'Treść i tak poszła do rdzenia — rozstrzygnie ją program powłoki karty.';
    tresc.ladowanie(`Uruchamianie pozycji ${pozycja.nazwa} w karcie ${karta.shell}…`);
    void zrodlo
      .wykonaj({ sessionId: karta.id, command: podstawienie.tresc })
      .then((wynik) => {
        if (!wynik.udany || wynik.wynik === undefined) {
          tresc.blad(
            `Rdzeń nie uruchomił pozycji ${pozycja.nazwa} (${Command.TerminalCommandExec}). ` +
              zdanieProgramuPowloki(karta.shell),
            wynik.blad,
          );
          return;
        }
        const proces = wynik.wynik;
        stan.zapiszProces(proces);
        stan.zapamietajPolecenie(karta.id, proces.command);
        biblioteka.odnotujUruchomienie(pozycja.nazwa, Date.now());
        pokaz();
        tresc.potwierdzenie(
          `Rdzeń uruchomił proces ${proces.id} dla pozycji ${pozycja.nazwa} — stan ${proces.status}. ${rozjazd}`.trim(),
          proces.status !== TerminalProcessStatus.Failed,
        );
      });
  }

  function sprawdzSkladnie(): void {
    const karta = kartaWykonania('sprawdzenie składni');
    if (karta === null) return;
    const sprawdzenie = polecenieSprawdzeniaSkladni(karta.shell, kontrolki.edytor.value);
    if (sprawdzenie.polecenie === '') {
      tresc.potwierdzenie(sprawdzenie.powod, false);
      return;
    }
    tresc.ladowanie('Sprawdzanie składni treści bez jej wykonania…');
    void zrodlo
      .wykonaj({ sessionId: karta.id, command: sprawdzenie.polecenie })
      .then(async (wynik) => {
        if (!wynik.udany || wynik.wynik === undefined) {
          tresc.blad(
            `Rdzeń nie uruchomił sprawdzenia składni (${Command.TerminalCommandExec}).`,
            wynik.blad,
          );
          return;
        }
        stan.zapiszProces(wynik.wynik);
        const odczyt = await zrodlo.odczytajWyjscie({ processId: wynik.wynik.id, waitMs: 10000 });
        if (!odczyt.udany || odczyt.wynik === undefined) {
          tresc.blad(
            `Rdzeń nie oddał wyniku sprawdzenia (${Command.TerminalOutputRead}).`,
            odczyt.blad,
          );
          return;
        }
        const odpowiedz = odczyt.wynik;
        const zastrzezenia = odpowiedz.stderr.trim();
        tresc.potwierdzenie(
          zastrzezenia === ''
            ? `Składnia poprawna — powłoka nie zgłosiła zastrzeżeń (stan ${odpowiedz.status}). ${sprawdzenie.powod}`
            : `Powłoka zgłosiła zastrzeżenia do składni: ${zastrzezenia}`,
          zastrzezenia === '',
        );
      });
  }

  kontrolki.zapisz.addEventListener('click', () => {
    const pozycja = pozycjaZFormularza();
    if (pozycja === null) {
      tresc.potwierdzenie('Pozycja bez nazwy albo bez treści nie ma czego zapisać.', false);
      return;
    }
    // Numeru wersji NIE wysyłamy: nadaje go rdzeń w jednej transakcji z zapisem
    // wersji, więc numer policzony w oknie rozjeżdżałby się przy dwóch zapisach
    // naraz.
    void zrodlo
      .zapiszSkrypt({
        script: {
          id: pozycja.id ?? '',
          name: pozycja.nazwa,
          kind: pozycja.rodzaj === 'snippet' ? 'snippet' : 'script',
          shell: pozycja.powloka,
          content: pozycja.tresc,
          version: 0,
          createdAt: 0,
          updatedAt: 0,
          ...(pozycja.tagi === '' ? {} : { tags: pozycja.tagi }),
        },
      })
      .then((wynik) => {
        if (!wynik.udany || wynik.wynik === undefined) {
          tresc.blad(`Rdzeń nie zapisał pozycji ${pozycja.nazwa} (${Command.TerminalScriptSave}).`, wynik.blad);
          return;
        }
        odczytaj();
        tresc.potwierdzenie(
          wynik.wynik.created
            ? `Pozycja ${wynik.wynik.script.name} założona w bibliotece rdzenia jako wersja ${wynik.wynik.script.version}.`
            : `Pozycja ${wynik.wynik.script.name} zapisana jako wersja ${wynik.wynik.script.version}; ` +
                'poprzednie brzmienie treści zostaje w dzienniku wersji.',
          true,
        );
      });
  });

  kontrolki.analizuj.addEventListener('click', analizuj);

  kontrolki.wstawSformatowana.addEventListener('click', () => {
    if (sformatowana === '') {
      tresc.potwierdzenie(
        'Nie ma treści sformatowanej — przeprowadź najpierw analizę treści.',
        false,
      );
      return;
    }
    kontrolki.edytor.value = sformatowana;
    pokaz();
    tresc.potwierdzenie('Treść sformatowana weszła do edytora. Zapis pozycji jest osobną czynnością.', true);
  });

  kontrolki.uruchom.addEventListener('click', () => {
    const pozycja = pozycjaZFormularza();
    if (pozycja === null) {
      tresc.blad('Pozycja bez nazwy albo bez treści nie ma czego uruchomić.');
      return;
    }
    uruchomTresc(pozycja);
  });

  kontrolki.sprawdz.addEventListener('click', sprawdzSkladnie);

  kontrolki.zPolecenia.addEventListener('click', () => {
    const karta = stan.kartaBiezaca();
    if (karta === null) {
      tresc.potwierdzenie('Nie ma karty bieżącej — nie ma czego przepisać do edytora.', false);
      return;
    }
    const polecenie = stan.ostatniePolecenie(karta.id);
    if (polecenie === '') {
      tresc.potwierdzenie('Karta bieżąca nie wykonała jeszcze żadnego polecenia.', false);
      return;
    }
    kontrolki.edytor.value = polecenie;
    pokaz();
    tresc.potwierdzenie(
      `Do edytora weszło ostatnie polecenie karty ${karta.shell}: „${polecenie}". Nadaj nazwę, wybierz rodzaj „Snippet" i zapisz.`,
      true,
    );
  });

  kontrolki.eksport.addEventListener('click', () => {
    const pozycje = biblioteka.pozycje();
    if (pozycje.length === 0) {
      tresc.potwierdzenie('Biblioteka jest pusta — pliku nie zapisano.', false);
      return;
    }
    const nazwa = `biblioteka-skryptow-${Date.now()}.json`;
    pobierzPlik(nazwa, JSON.stringify(pozycje, null, 2), 'text/plain');
    tresc.potwierdzenie(`Zapisano ${nazwa} — ${pozycje.length} pozycji biblioteki widoku.`, true);
  });

  // Przerysowanie idzie dopiero wtedy, gdy zmienia się to, co z treści wynika:
  // deklaracja parametrów albo uwagi kontroli wstępnej. Przerysowywanie przy
  // każdym naciśnięciu klawisza budowałoby formularz parametrów od nowa
  // dziesiątki razy w trakcie pisania jednego wiersza.
  let podpisTresci = '';
  kontrolki.edytor.addEventListener('input', () => {
    const podpis = JSON.stringify([
      czytajParametry(kontrolki.edytor.value),
      kontrolaWstepna(kontrolki.edytor.value),
    ]);
    if (podpis === podpisTresci) return;
    podpisTresci = podpis;
    pokaz();
  });

  const czynnosci: readonly CzynnoscOkna[] = [
    {
      okno: 'Script Library',
      nazwa: 'Uruchom treść edytora',
      opis: 'Wysyła treść z edytora do karty bieżącej jako jeden proces powłoki.',
      warstwa: 'zawsze',
      wykonaj: () => kontrolki.uruchom.click(),
    },
    {
      okno: 'Script Library',
      nazwa: 'Zapisz pozycję biblioteki',
      opis: 'Zapisuje treść edytora jako kolejną wersję pozycji o wpisanej nazwie.',
      warstwa: 'na-zadanie',
      wykonaj: () => kontrolki.zapisz.click(),
    },
    {
      okno: 'Script Library',
      nazwa: 'Sprawdź składnię bez uruchomienia',
      opis: 'Podaje treść powłoce do kontroli składni; żadne polecenie skryptu się nie wykonuje.',
      warstwa: 'na-zadanie',
      wykonaj: sprawdzSkladnie,
    },
    {
      okno: 'Script Library',
      nazwa: 'Przepisz ostatnie polecenie karty',
      opis: 'Wstawia do edytora ostatnie polecenie wykonane w karcie bieżącej.',
      warstwa: 'kontekstowa',
      wykonaj: () => kontrolki.zPolecenia.click(),
    },
    {
      okno: 'Script Library',
      nazwa: 'Analiza statyczna treści',
      opis: 'Poddaje treść programowi analizy na maszynie rdzenia i pokazuje jego uwagi.',
      warstwa: 'na-zadanie',
      wykonaj: () => kontrolki.analizuj.click(),
    },
    {
      okno: 'Script Library',
      nazwa: 'Wstaw treść sformatowaną',
      opis: 'Zamienia treść edytora na wersję sformatowaną przez ostatnią analizę.',
      warstwa: 'kontekstowa',
      wykonaj: () => kontrolki.wstawSformatowana.click(),
    },
    {
      okno: 'Script Library',
      nazwa: 'Odczytaj bibliotekę z rdzenia',
      opis: 'Czyta pozycje biblioteki z dziennika rdzenia i przerysowuje wykaz.',
      warstwa: 'kontekstowa',
      wykonaj: () => odczytaj(),
    },
    {
      okno: 'Script Library',
      nazwa: 'Eksportuj bibliotekę',
      opis: 'Zapisuje pozycje biblioteki do pliku — droga wyniesienia ich poza rdzeń.',
      warstwa: 'kontekstowa',
      wykonaj: () => kontrolki.eksport.click(),
    },
  ];

  return { element: rama.element, odswiez: odczytaj, czynnosci };
}

/** Czynności wiersza biblioteki — wykaz nie zna ani rdzenia, ani stanu okna. */
interface CzynnosciPozycji {
  wczytaj(pozycja: PozycjaBiblioteki): void;
  uruchom(pozycja: PozycjaBiblioteki): void;
  usun(pozycja: PozycjaBiblioteki): void;
}

/** Wykaz pozycji biblioteki; pustka ma własne zdanie, bo jest stanem poprawnym. */
function wykazPozycji(biblioteka: Biblioteka, czynnosci: CzynnosciPozycji): HTMLElement {
  const blok = document.createElement('section');
  blok.className = 'dt-biblioteka';

  const podpis = document.createElement('h4');
  podpis.className = 'dt-biblioteka__podpis';
  podpis.textContent = 'Biblioteka widoku';
  blok.append(podpis);

  const pozycje = biblioteka.pozycje();
  if (pozycje.length === 0) {
    const pusto = document.createElement('p');
    pusto.className = 'dn-pole-opis';
    pusto.textContent =
      'Biblioteka jest pusta. Wpisz nazwę i treść, po czym zapisz pozycję — zapis idzie do dziennika ' +
      'rdzenia i zakłada kolejną wersję pozycji, więc poprzednie brzmienie treści zostaje.';
    blok.append(pusto);
    return blok;
  }

  const lista = wykaz('Pozycje biblioteki skryptów', 'dt-wykaz');
  for (const pozycja of pozycje) {
    const wiersz = pozycjaWykazu(pozycja.nazwa, opisPozycji(pozycja), 'dt');
    wiersz.element.dataset['rodzaj'] = pozycja.rodzaj;

    const uruchom = przyciskAkcji('Uruchom', 'dn-btn dn-btn--atrament');
    uruchom.title = zdanieProgramuPowloki(pozycja.powloka);
    uruchom.addEventListener('click', () => czynnosci.uruchom(pozycja));

    const wczytaj = przyciskAkcji('Wczytaj do edytora');
    wczytaj.addEventListener('click', () => czynnosci.wczytaj(pozycja));

    const usun = przyciskAkcji('Usuń pozycję');
    usun.addEventListener('click', () => czynnosci.usun(pozycja));

    wiersz.akcje.append(uruchom, wczytaj, usun);
    lista.append(wiersz.element);
  }
  blok.append(lista);
  return blok;
}

function opisPozycji(pozycja: PozycjaBiblioteki): string {
  const czesci = [
    pozycja.rodzaj === 'snippet' ? 'snippet' : 'skrypt',
    `powłoka: ${pozycja.powloka} (${programPowloki(pozycja.powloka)?.program ?? 'rdzeń nie ma dla niej programu'})`,
    `wersja: ${pozycja.wersja}`,
  ];
  if (pozycja.tagi !== '') czesci.push(`znaczniki: ${pozycja.tagi}`);
  czesci.push(
    pozycja.ostatnieUruchomienie === 0
      ? 'nie uruchamiano'
      : `ostatnie uruchomienie: ${new Date(pozycja.ostatnieUruchomienie)
          .toISOString()
          .replace('T', ' ')
          .slice(0, 19)}`,
  );
  return czesci.join(' · ');
}

/**
 * Formularz parametrów zbudowany z deklaracji w nagłówku treści.
 *
 * Pola powstają przy każdym przerysowaniu, więc wpisane wartości trzyma mapa
 * przekazana z okna, a nie sam węzeł — inaczej znikałyby przy każdym naciśnięciu
 * klawisza w edytorze.
 */
function formularzParametrow(
  parametry: readonly ParametrSkryptu[],
  wartosci: Map<string, string>,
): HTMLElement {
  const blok = document.createElement('section');
  blok.className = 'dt-parametry';

  const podpis = document.createElement('h4');
  podpis.className = 'dt-parametry__podpis';
  podpis.textContent = 'Parametry pozycji';
  blok.append(podpis);

  if (parametry.length === 0) {
    const pusto = document.createElement('p');
    pusto.className = 'dn-pole-opis';
    pusto.textContent =
      'Treść nie deklaruje parametrów. Deklaracja stoi w komentarzu na początku treści: wiersz „args:”, ' +
      'a pod nim wcięte pary nazwa i opis. W treści parametr występuje jako nazwa w podwójnych nawiasach klamrowych.';
    blok.append(pusto);
    return blok;
  }

  for (const parametr of parametry) {
    const kontrolka = pole(`Parametr ${parametr.nazwa}`, parametr.opis);
    kontrolka.value = wartosci.get(parametr.nazwa) ?? '';
    kontrolka.addEventListener('input', () => wartosci.set(parametr.nazwa, kontrolka.value));
    const wiersz = document.createElement('div');
    wiersz.className = 'dt-parametry__wiersz';
    const nazwa = document.createElement('span');
    nazwa.className = 'dn-pole-etykieta';
    nazwa.textContent = parametr.nazwa;
    wiersz.append(nazwa, kontrolka);
    if (parametr.opis !== '') {
      const opis = document.createElement('span');
      opis.className = 'dn-pole-opis';
      opis.textContent = parametr.opis;
      wiersz.append(opis);
    }
    blok.append(wiersz);
  }
  return blok;
}

/**
 * Wykaz uwag ANALIZY z maszyny rdzenia.
 *
 * Osobny od wykazu kontroli wstępnej, bo mówi o czym innym: kontrola wstępna
 * jest sprawdzeniem, które okno robi samo i natychmiast, a to jest wynik
 * programu analizy. Nagłówek niesie nazwę tego programu — bez niej Operator nie
 * wie, czyje to zastrzeżenia i czego doinstalować, gdy programu zabrakło.
 */
function wykazUwagAnalizy(odpowiedz: TerminalScriptLintResponse): HTMLElement {
  const blok = document.createElement('section');
  blok.className = 'dt-uwagi';

  const podpis = document.createElement('h4');
  podpis.className = 'dt-uwagi__podpis';
  podpis.textContent = `Analiza statyczna — ${odpowiedz.analyzer}`;
  blok.append(podpis);

  if (!odpowiedz.analyzerAvailable) {
    const brak = document.createElement('p');
    brak.className = 'dn-pole-opis';
    brak.dataset['stan'] = 'brak-narzedzia';
    brak.textContent =
      `Treści NIE sprawdzono: programu analizy (${odpowiedz.analyzer}) nie ma na maszynie rdzenia. ` +
      'Pusty wykaz uwag nie znaczy tu „treść bez zastrzeżeń” — znaczy „nie było czym sprawdzić”.';
    blok.append(brak);
    return blok;
  }

  if (odpowiedz.findings.length === 0) {
    const czysto = document.createElement('p');
    czysto.className = 'dn-pole-opis';
    czysto.textContent = 'Program analizy nie zgłosił zastrzeżeń do treści.';
    blok.append(czysto);
    return blok;
  }

  const lista = wykaz('Uwagi analizy statycznej', 'dt-wykaz');
  for (const uwaga of odpowiedz.findings) {
    const miejsce =
      uwaga.column === undefined
        ? `wiersz ${uwaga.line}`
        : `wiersz ${uwaga.line}, kolumna ${uwaga.column}`;
    const regula = uwaga.rule === undefined ? '' : ` [${uwaga.rule}]`;
    lista.append(pozycjaWykazu(`${uwaga.severity}: ${miejsce}${regula}`, uwaga.message, 'dt').element);
  }
  blok.append(lista);
  return blok;
}

/** Uwagi kontroli wstępnej; brak uwag też jest zdaniem, bo cisza znaczyłaby „nie sprawdzono”. */
function wykazUwag(uwagi: readonly string[]): HTMLElement {
  const blok = document.createElement('section');
  blok.className = 'dt-uwagi';

  const podpis = document.createElement('h4');
  podpis.className = 'dt-uwagi__podpis';
  podpis.textContent = 'Kontrola wstępna treści';
  blok.append(podpis);

  if (uwagi.length === 0) {
    const zdanie = document.createElement('p');
    zdanie.className = 'dn-pole-opis';
    zdanie.textContent =
      'Kontrola wstępna nie ma uwag. Sprawdza pustkę, znaki końca wiersza, znacznik kolejności bajtów ' +
      'i zgodność deklaracji parametrów z ich użyciem — nie jest analizą statyczną i nie zastępuje lintera.';
    blok.append(zdanie);
    return blok;
  }

  const lista = wykaz('Uwagi kontroli wstępnej', 'dt-wykaz');
  for (const uwaga of uwagi) {
    const wiersz = document.createElement('li');
    wiersz.className = 'dt-pozycja';
    wiersz.dataset['stan'] = 'failed';
    wiersz.textContent = uwaga;
    lista.append(wiersz);
  }
  blok.append(lista);
  return blok;
}

/** Przepisuje pozycję do formularza wraz z wyczyszczeniem wartości parametrów poprzedniej pozycji. */
function wczytajDoEdytora(
  kontrolki: PowierzchniaBiblioteki,
  wartosci: Map<string, string>,
  pozycja: PozycjaBiblioteki,
  przerysuj: () => void,
): void {
  kontrolki.nazwa.value = pozycja.nazwa;
  kontrolki.tagi.value = pozycja.tagi;
  kontrolki.edytor.value = pozycja.tresc;
  wartosci.clear();
  przerysuj();
  kontrolki.edytor.focus();
}

/** Kontrolki okna Script Library. */
interface PowierzchniaBiblioteki {
  nazwa: HTMLInputElement;
  tagi: HTMLInputElement;
  rodzaj: WyborDrzewem;
  powloka: WyborDrzewem;
  edytor: HTMLTextAreaElement;
  zapisz: HTMLButtonElement;
  uruchom: HTMLButtonElement;
  sprawdz: HTMLButtonElement;
  analizuj: HTMLButtonElement;
  wstawSformatowana: HTMLButtonElement;
  zPolecenia: HTMLButtonElement;
  eksport: HTMLButtonElement;
}

/**
 * Składa kontrolki, pasek akcji, pasek narzędzi i ciało okna.
 *
 * Edytor powstaje raz i wchodzi do ciała przy każdym przerysowaniu — treść
 * pisana przez Operatora nie ma prawa zniknąć przy odświeżeniu wykazu.
 * Pozycje bez pokrycia w kontrakcie stoją jawnie nieczynne wraz z powodem.
 */
function zlozPowierzchnieBiblioteki(
  rama: { akcje: HTMLElement; narzedzia: HTMLElement; cialo: HTMLElement },
  stanTresci: HTMLElement,
  pokrycie: PokrycieKomend,
): PowierzchniaBiblioteki {
  const nazwa = pole('Nazwa pozycji biblioteki', 'np. wdrozenie-srodowiska-probnego');
  const tagi = pole('Znaczniki pozycji', 'np. wdrozenie, baza');
  const rodzaj = utworzWyborDrzewem({ nastawa: 'Rodzaj pozycji', pozycje: RODZAJE });
  const powloka = utworzWyborDrzewem({ nastawa: 'Powłoka pozycji', pozycje: POWLOKI });
  const edytor = poleTresci(
    'Treść pozycji biblioteki',
    12,
    '# args:\n#   srodowisko: nazwa środowiska wdrożenia\nnpm run wdroz -- --env={{srodowisko}}',
    'dt-edytor',
  );

  const zapisz = przyciskAkcji('Zapisz pozycję');
  const uruchom = przyciskAkcji('Uruchom treść', 'dn-btn dn-btn--atrament');
  const sprawdz = przyciskAkcji('Sprawdź składnię');
  const zPolecenia = przyciskAkcji('Przepisz ostatnie polecenie karty');
  const eksport = przyciskAkcji('Eksportuj bibliotekę');

  const analizuj = przyciskAkcji('Analiza statyczna treści');
  analizuj.title =
    'Poddaje treść analizie programem właściwym powłoce (ShellCheck, PSScriptAnalyzer, sam interpreter). ' +
    'Odpowiedź mówi wprost, czy program analizy stoi na maszynie rdzenia — pusty wykaz uwag przy jego ' +
    'braku znaczyłby fałszywie „treść bez zastrzeżeń”.';
  const wstawSformatowana = przyciskAkcji('Wstaw treść sformatowaną');
  wstawSformatowana.title =
    'Zamienia treść w edytorze na wersję sformatowaną przez ostatnią analizę. Czynność jest osobna, ' +
    'bo podmiana treści bez wskazania Operatora zabrałaby mu wersję, którą właśnie pisze.';

  const generowanie = pokrycie.przycisk(
    'Wygeneruj skrypt opisem',
    Command.MessageSend,
    'Ułożenie treści skryptu z opisu w języku naturalnym; robi to model w oknie rozmowy modułu, którego to złożenie nie osadza',
  );
  // Nazwa spoza kontraktu jest tu wskazaniem, nie zapisem stanu: przenoszenia
  // plików konfiguracyjnych powłoki świadomie nie zgłoszono do scalenia, bo
  // wersjonowanie repozytorium prowadzi moduł Developer.
  const dotfiles = pokrycie.przycisk(
    'Pliki konfiguracyjne powłoki',
    'terminal.dotfiles.sync',
    'Wersjonowanie i przenoszenie plików konfiguracyjnych powłoki między hostami',
  );

  oznaczWarstwy([
    [uruchom, 'zawsze'],
    [zapisz, 'na-zadanie'],
    [sprawdz, 'na-zadanie'],
    [analizuj, 'na-zadanie'],
    [zPolecenia, 'kontekstowa'],
    [eksport, 'kontekstowa'],
    [wstawSformatowana, 'kontekstowa'],
    [generowanie, 'ekspercka'],
    [dotfiles, 'ekspercka'],
  ]);

  rama.akcje.append(
    uruchom,
    zapisz,
    sprawdz,
    analizuj,
    wstawSformatowana,
    zPolecenia,
    eksport,
    generowanie,
    dotfiles,
  );
  rama.narzedzia.append(nazwa, tagi, rodzaj.element, powloka.element);
  // Edytor stoi w ciele okna na stałe, a nie w miejscu treści przerysowywanym
  // przy każdej zmianie: element wyjęty z dokumentu traci ognisko, a Operator
  // traciłby je wtedy w połowie pisanego wiersza.
  rama.cialo.append(edytor, stanTresci);

  return {
    nazwa,
    tagi,
    rodzaj,
    powloka,
    edytor,
    zapisz,
    uruchom,
    sprawdz,
    analizuj,
    wstawSformatowana,
    zPolecenia,
    eksport,
  };
}
