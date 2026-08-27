import {
  GitActionKind,
  type DeveloperGitActionRequest,
  type GitActionResult,
} from '../../../../shared/contract';
import { utworzRameOkna } from '../../komponenty/rama-okna';
import {
  pole,
  poleTresci,
  przyciskAkcji as przycisk,
  przyciskBezKomendy,
  pozycjaWykazu,
  wiersz,
  wykaz,
} from '../../modele/kontrolki-formularza';
import { powodBezKomendy } from './braki-kontraktu';
import {
  KATALOG_W_ODCZYCIE,
  odczytajKatalogKomend,
  zdanieWiedzyGitPanelu,
} from './katalog-komend';
import { rysujZaleznosci, zaleznosci } from './zaleznosci-zewnetrzne';
import type { StanDevelopera } from './stan-developer';
import { utworzStanTresci, type StanTresci } from './stany-okna';
import { NAZWY_CZYNNOSCI, zdanieCzynnosciRepozytorium } from './zdania-odpowiedzi';
import type { ZrodloDeveloper } from './zrodlo-developer';

/**
 * Git Panel jest oknem zarządcy modułu Developer: zatwierdzanie zmian,
 * przełączanie gałęzi i praca ze zdalnym repozytorium sesji. Okno poznaje
 * repozytorium dopiero z wyniku wykonanej czynności, nie zakłada go z góry.
 */
export interface OknoGitPanelu {
  element: HTMLElement;
  odswiez(): void;
  /** Zamyka nasłuch `stan.naZmiane` założony przy konstrukcji okna. */
  zamknij(): void;
}

export function utworzOknoGitPanelu(zrodlo: ZrodloDeveloper, stan: StanDevelopera): OknoGitPanelu {
  const rama = utworzRameOkna({
    tytul: 'Git Panel',
    rola: 'zarządca',
    kod: 'git-panel',
    przeznaczenie:
      'Zatwierdzenie zmian, przełączanie gałęzi i praca ze zdalnym repozytorium sesji. ' +
      'Wiedza o repozytorium przychodzi dopiero z wyniku wykonanej czynności.',
    modul: 'Developer',
    przedrostek: 'mdev',
  });
  const tresc = utworzStanTresci();
  const powierzchnia = zlozPowierzchnieGitPanelu(rama, tresc.element);
  let zdanieWiedzy = KATALOG_W_ODCZYCIE;

  function odswiez(): void {
    tresc.pusto(zdanieWiedzy);
  }

  // Wykaz komend pochodzi z rdzenia; odczyt odświeża stan pusty tylko, gdy okno nadal go pokazuje.
  void odczytajKatalogKomend(zrodlo).then((katalog) => {
    zdanieWiedzy = zdanieWiedzyGitPanelu(katalog);
    if (tresc.rodzaj() === 'pusto') odswiez();
  });

  function wykonaj(action: GitActionKind, force = false): void {
    const opis = NAZWY_CZYNNOSCI[action];
    const zadanie = zlozZadanie(stan.okno(), action, powierzchnia, force);
    tresc.ladowanie(`${opis}${force ? ' (wymuszone)' : ''}…`);
    void zrodlo.czynnoscRepozytorium(zadanie).then((wynik) => {
      if (!wynik.udany || wynik.wynik === undefined) {
        tresc.blad(`Rdzeń odmówił czynności „${opis}”.`, wynik.blad);
        return;
      }
      przejmijGalaz(powierzchnia, wynik.wynik);
      rysujWynikCzynnosci(wynik.wynik, tresc, (sciezka) => stan.wskazPlik(sciezka));
      // Nazwa w potwierdzeniu pochodzi z pola action wyniku, nie z etykiety przycisku.
      const potwierdzenie = zdanieCzynnosciRepozytorium(action, wynik.wynik);
      tresc.potwierdzenie(potwierdzenie.zdanie, potwierdzenie.udane);
    });
  }

  podepnijAkcjeGitPanelu(powierzchnia, tresc, wykonaj);
  const odsubskrybuj = podepnijWskazanie(powierzchnia, stan);

  // Jeden punkt wejścia dla odczytu: złożenie modułu woła odswiez() samo, dla każdego okna.
  return { element: rama.element, odswiez, zamknij: odsubskrybuj };
}

/**
 * Wiąże panel ze wspólnym wskazaniem pliku modułu; wskazanie i pole Ścieżki
 * to dwa osobne byty, a podpowiedź w polu wchodzi tam tylko, dopóki Operator
 * sam w nim nie pisał.
 */
function podepnijWskazanie(powierzchnia: PowierzchniaGitPanelu, stan: StanDevelopera): () => void {
  let tknietePrzezOperatora = false;
  powierzchnia.sciezki.addEventListener('input', () => {
    tknietePrzezOperatora = powierzchnia.sciezki.value.trim() !== '';
  });

  function odswiezWskazanie(): void {
    const sciezka = stan.sciezka();
    powierzchnia.wskazanie.dataset['wskazana'] = sciezka === '' ? 'nie' : 'tak';
    powierzchnia.wskazanie.textContent =
      sciezka === ''
        ? 'Nie wskazano pliku — wskazanie przychodzi z Project Tree, z Build Output i z wykazu poniżej.'
        : `Plik wskazany we wspólnym stanie modułu: ${sciezka}`;
    if (!tknietePrzezOperatora) powierzchnia.sciezki.value = sciezka;
  }

  odswiezWskazanie();
  return stan.naZmiane(odswiezWskazanie);
}

/**
 * Przejmuje gałąź z wyniku rdzenia; pole Gałąź dostaje ją tylko wtedy, gdy
 * jest puste, a wiersz obok pokazuje ją zawsze, bez nadpisywania wpisu
 * Operatora.
 */
function przejmijGalaz(powierzchnia: PowierzchniaGitPanelu, wynik: GitActionResult): void {
  if (wynik.branch === undefined || wynik.branch === '') return;
  powierzchnia.galazWyniku.textContent = `Gałąź wedle ostatniego wyniku rdzenia: ${wynik.branch}`;
  powierzchnia.galazWyniku.dataset['wskazana'] = 'tak';
  if (powierzchnia.galaz.value.trim() === '') powierzchnia.galaz.value = wynik.branch;
}

/** Pola wspólne wszystkich czynności Git Panelu: ścieżki objęte czynnością, opis, gałąź oraz zdalne repozytorium. */
interface PolaGitPanelu {
  sciezki: HTMLInputElement;
  opisZatwierdzenia: HTMLTextAreaElement;
  galaz: HTMLInputElement;
  zdalne: HTMLInputElement;
}

/** Przyciski panelu akcji Git Panelu, pogrupowane wedle inwentarza czynności GitActionKind dostępnych w kontrakcie. */
interface AkcjeGitPanelu {
  stage: HTMLButtonElement;
  unstage: HTMLButtonElement;
  commit: HTMLButtonElement;
  amend: HTMLButtonElement;
  revert: HTMLButtonElement;
  checkout: HTMLButtonElement;
  merge: HTMLButtonElement;
  rebase: HTMLButtonElement;
  tag: HTMLButtonElement;
  fetch: HTMLButtonElement;
  pull: HTMLButtonElement;
  push: HTMLButtonElement;
  pushWymuszony: HTMLButtonElement;
  stash: HTMLButtonElement;
  stashPop: HTMLButtonElement;
}

/**
 * Składa piętnaście przycisków: czternaście czynności GitActionKind plus
 * wariant wymuszony wysłania, oraz przyciski bez komendy dla czynności,
 * których kontrakt nie niesie.
 */
function zlozAkcjeGitPanelu(gospodarz: HTMLElement): AkcjeGitPanelu {
  const stage = przycisk('Dodaj do indeksu (stage)');
  const unstage = przycisk('Wycofaj z indeksu (unstage)');
  const commit = przycisk('Zatwierdź zmiany', 'dn-btn dn-btn--atrament');
  const amend = przycisk('Popraw ostatnie zatwierdzenie (amend)');
  const revert = przycisk('Odwróć zatwierdzenie (revert)');
  const checkout = przycisk('Przełącz gałąź / wersję (checkout)');
  const merge = przycisk('Scal gałąź (merge)');
  const rebase = przycisk('Przestaw gałąź (rebase)');
  const tag = przycisk('Nadaj etykietę (tag)');
  const fetch = przycisk('Pobierz zmiany zdalne (fetch)');
  const pull = przycisk('Pobierz i scal (pull)');
  const push = przycisk('Wyślij zmiany (push)');
  // Wymuszenie jest niebezpieczne, więc dostaje jedyny przycisk z force: true, widocznie odróżniony.
  const pushWymuszony = przycisk('Wyślij zmiany wymuszone', 'dn-btn dn-btn--niebezpieczny');
  const stash = przycisk('Odłóż zmiany (stash)');
  const stashPop = przycisk('Przywróć odłożone (stash pop)');

  gospodarz.append(
    stage,
    unstage,
    commit,
    amend,
    revert,
    checkout,
    merge,
    rebase,
    tag,
    fetch,
    pull,
    push,
    pushWymuszony,
    stash,
    stashPop,
    przyciskBezKomendy(
      '✨ Generuj opis commitu (AI)',
      powodBezKomendy(
        'Opis zatwierdzenia z rzeczywistej różnicy wymagałby dwóch komend naraz: ' +
          'developer.git.diff, żeby różnicę policzyć, i developer.contextual.op, żeby wywołać nad ' +
          'nią model. Dziś opis wpisuje Operator.',
      ),
    ),
    przyciskBezKomendy(
      'Pokaż różnicę',
      powodBezKomendy(
        'Widok różnicowy per plik i per fragment wymagałby komendy developer.git.diff; ' +
          'developer.git.action oddaje wyłącznie wykaz ścieżek objętych czynnością.',
      ),
    ),
    przyciskBezKomendy(
      'Historia zatwierdzeń',
      powodBezKomendy(
        'Log zatwierdzeń z autorem, datą i skrótem wymagałby komendy developer.git.log.',
      ),
    ),
    przyciskBezKomendy(
      'Wykaz gałęzi',
      powodBezKomendy(
        'Wykaz gałęzi wraz z rozbieżnością wobec gałęzi zdalnej wymagałby komendy ' +
          'developer.git.branch.list; selektor gałęzi ma dziś pole tekstowe, bo nie ma z czego ' +
          'zbudować listy.',
      ),
    ),
    przyciskBezKomendy(
      'Stan repozytorium',
      powodBezKomendy(
        'Stan zmian — niezatwierdzone, przygotowane, nieśledzone — wymagałby komendy ' +
          'developer.git.status. Rdzeń ten stan JUŻ liczy po każdej czynności i wkłada gałąź oraz ' +
          'zatwierdzenie do wyniku; brakuje wyłącznie komendy oddającej go bez wykonywania czynności.',
      ),
    ),
    przyciskBezKomendy(
      'Rozwiąż konflikty (widok trójstronny)',
      powodBezKomendy(
        'Kontrakt zwraca wykaz ścieżek konfliktu, nie treść trzech wersji — widok bieżąca / ' +
          'przychodząca / wynik wymagałby komend developer.git.conflict.get ' +
          'i developer.git.conflict.resolve.',
      ),
    ),
    przyciskBezKomendy(
      'Changelog przez AI',
      powodBezKomendy(
        'Nota wydania dla zakresu zatwierdzeń wymagałaby developer.git.log, żeby zakres zebrać, ' +
          'i developer.contextual.op, żeby go modelowi przekazać.',
      ),
    ),
    przyciskBezKomendy(
      'Przegląd kodu na różnicy (AI)',
      powodBezKomendy(
        'Uwagi przypięte do konkretnych linii różnicy wymagałyby developer.git.diff ' +
          'i developer.contextual.op.',
      ),
    ),
  );
  return {
    stage,
    unstage,
    commit,
    amend,
    revert,
    checkout,
    merge,
    rebase,
    tag,
    fetch,
    pull,
    push,
    pushWymuszony,
    stash,
    stashPop,
  };
}

/** Kontrolki tego okna: pola wspólne czynności wraz z panelem akcji i ciałem ramy tego okna Git Panelu. */
interface PowierzchniaGitPanelu extends AkcjeGitPanelu, PolaGitPanelu {
  /** Lustro wspólnego wskazania pliku — nigdy nie jedzie do rdzenia. */
  wskazanie: HTMLElement;
  /** Lustro gałęzi oddanej w ostatnim wyniku — też tylko do czytania. */
  galazWyniku: HTMLElement;
  /** Licznik znaków pierwszej linii opisu — czynność wyłącznie kliencka. */
  licznikOpisu: HTMLElement;
}

/**
 * Zalecana górna granica pierwszej linii opisu zatwierdzenia, przy której git
 * log --oneline i wykazy hostingów przestają skracać podsumowanie; licznik
 * jest wyłącznie informacyjny.
 */
const ZALECANA_DLUGOSC_PIERWSZEJ_LINII = 50;

/**
 * Odświeża licznik pierwszej linii opisu zatwierdzenia.
 *
 * Liczy się pierwsza linia, nie cały opis: konwencja zatwierdzeń dzieli treść
 * na podsumowanie i akapit uzasadnienia, a granicę przekracza wyłącznie
 * podsumowanie.
 */
function odswiezLicznikOpisu(miejsce: HTMLElement, opis: string): void {
  const pierwsza = opis.split('\n')[0] ?? '';
  const dlugosc = pierwsza.length;
  const zaDlugo = dlugosc > ZALECANA_DLUGOSC_PIERWSZEJ_LINII;
  miejsce.dataset['udane'] = String(!zaDlugo);
  miejsce.hidden = false;
  miejsce.textContent = zaDlugo
    ? `Pierwsza linia opisu: ${dlugosc} znaków — powyżej zalecanych ` +
      `${ZALECANA_DLUGOSC_PIERWSZEJ_LINII}. Wykazy skracają dłuższe podsumowania. ` +
      'Zatwierdzenie mimo to pojedzie do rdzenia — to zalecenie, nie blokada.'
    : `Pierwsza linia opisu: ${dlugosc} z zalecanych ${ZALECANA_DLUGOSC_PIERWSZEJ_LINII} znaków.`;
}

function zlozPowierzchnieGitPanelu(
  rama: { akcje: HTMLElement; narzedzia: HTMLElement; cialo: HTMLElement },
  stanTresci: HTMLElement,
): PowierzchniaGitPanelu {
  const sciezki = pole('Ścieżki objęte czynnością', 'np. internal/core/brama.go, druga/sciezka.ts');
  const opisZatwierdzenia = poleTresci('Opis zatwierdzenia / stash', 3, 'wiadomość commitu albo notatka stash');
  const galaz = pole('Gałąź / etykieta', 'np. main, v1.2.0');
  const zdalne = pole('Repozytorium zdalne', 'np. origin');
  const akcje = zlozAkcjeGitPanelu(rama.akcje);

  const wskazanie = document.createElement('p');
  wskazanie.className = 'mdev-sciezka';
  wskazanie.dataset['wskazana'] = 'nie';

  const galazWyniku = document.createElement('p');
  galazWyniku.className = 'mdev-sciezka';
  galazWyniku.dataset['wskazana'] = 'nie';
  galazWyniku.textContent = 'Gałęzi rdzeń jeszcze nie podał — poda ją wynik czynności.';

  const licznikOpisu = document.createElement('p');
  licznikOpisu.className = 'mdev-potwierdzenie';
  odswiezLicznikOpisu(licznikOpisu, '');
  opisZatwierdzenia.addEventListener('input', () =>
    odswiezLicznikOpisu(licznikOpisu, opisZatwierdzenia.value),
  );

  // Git jest programem spoza instalki, a stoi na nim każda czynność tego okna Git Panelu.
  rama.narzedzia.append(rysujZaleznosci(zaleznosci(['git'])));

  rama.cialo.append(
    wskazanie,
    wiersz('Ścieżki', sciezki, {
      klasa: 'mdev-wiersz',
      objasnienie:
        'Rozdzielone przecinkiem; puste pole obejmuje czynnością cały indeks. Ścieżka wskazana ' +
        'podstawia się sama, dopóki nie wpiszesz własnej.',
    }),
    wiersz('Opis', opisZatwierdzenia, {
      klasa: 'mdev-wiersz',
      objasnienie: 'Wymagany przy zatwierdzeniu — puste pole nie jedzie do rdzenia jako commit.',
    }),
    licznikOpisu,
    wiersz('Gałąź', galaz, { klasa: 'mdev-wiersz' }),
    galazWyniku,
    wiersz('Zdalne', zdalne, { klasa: 'mdev-wiersz' }),
    stanTresci,
  );
  return {
    sciezki,
    opisZatwierdzenia,
    galaz,
    zdalne,
    wskazanie,
    galazWyniku,
    licznikOpisu,
    ...akcje,
  };
}

/**
 * Podpina piętnaście przycisków do jednej czynności parametryzowanej kodem;
 * nazwy czynności mają jedno źródło w zdania-odpowiedzi.ts, żeby dwa wykazy
 * nazw się nie rozjechały.
 */
function podepnijAkcjeGitPanelu(
  powierzchnia: PowierzchniaGitPanelu,
  tresc: StanTresci,
  wykonaj: (action: GitActionKind, force?: boolean) => void,
): void {
  const zwykle: ReadonlyArray<readonly [HTMLButtonElement, GitActionKind]> = [
    [powierzchnia.stage, GitActionKind.Stage],
    [powierzchnia.unstage, GitActionKind.Unstage],
    [powierzchnia.amend, GitActionKind.Amend],
    [powierzchnia.revert, GitActionKind.Revert],
    [powierzchnia.checkout, GitActionKind.Checkout],
    [powierzchnia.merge, GitActionKind.Merge],
    [powierzchnia.rebase, GitActionKind.Rebase],
    [powierzchnia.tag, GitActionKind.Tag],
    [powierzchnia.fetch, GitActionKind.Fetch],
    [powierzchnia.pull, GitActionKind.Pull],
    [powierzchnia.push, GitActionKind.Push],
    [powierzchnia.stash, GitActionKind.Stash],
    [powierzchnia.stashPop, GitActionKind.StashPop],
  ];
  for (const [przycisk, czynnosc] of zwykle) {
    przycisk.addEventListener('click', () => wykonaj(czynnosc));
  }
  powierzchnia.commit.addEventListener('click', () => {
    if (powierzchnia.opisZatwierdzenia.value.trim() === '') {
      wykonajOdmowaBrakuOpisu(powierzchnia, tresc);
      return;
    }
    wykonaj(GitActionKind.Commit);
  });
  // Wymuszenie NIE jest domyślne — to jedyny przycisk, który je włącza.
  powierzchnia.pushWymuszony.addEventListener('click', () => wykonaj(GitActionKind.Push, true));
}

/**
 * Zatwierdzenie bez opisu nie jedzie do rdzenia jako puste, a Operator dostaje
 * o tym zdanie zamiast samego ogniska w polu; to lokalna walidacja, nie
 * odmowa rdzenia.
 */
function wykonajOdmowaBrakuOpisu(powierzchnia: PowierzchniaGitPanelu, tresc: StanTresci): void {
  tresc.potwierdzenie('Zatwierdzenie wymaga opisu — pole jest puste, więc żądanie nie poszło.', false);
  powierzchnia.opisZatwierdzenia.focus();
}

/** Buduje żądanie czynności repozytorium z pól okna; puste pola nie trafiają do żądania wysyłanego rdzeniowi. */
function zlozZadanie(
  idOkna: string,
  action: GitActionKind,
  pola: PolaGitPanelu,
  force: boolean,
): DeveloperGitActionRequest {
  const zadanie: DeveloperGitActionRequest = { windowId: idOkna, action };
  const sciezki = pola.sciezki.value
    .split(',')
    .map((sciezka) => sciezka.trim())
    .filter((sciezka) => sciezka !== '');
  if (sciezki.length > 0) zadanie.paths = sciezki;
  const opis = pola.opisZatwierdzenia.value.trim();
  if (opis !== '') zadanie.message = opis;
  const galaz = pola.galaz.value.trim();
  if (galaz !== '') zadanie.branch = galaz;
  const zdalne = pola.zdalne.value.trim();
  if (zdalne !== '') zadanie.remote = zdalne;
  if (force) zadanie.force = true;
  return zadanie;
}

/**
 * Rysuje GitActionResult: wyjście polecenia, gałąź po czynności, ścieżki
 * objęte czynnością i, wyróżnione osobno, ścieżki konfliktu z regułą wstęgi
 * w arkuszu modułu.
 */
function rysujWynikCzynnosci(
  wynik: GitActionResult,
  tresc: StanTresci,
  naWskazanie: (sciezka: string) => void,
): void {
  const miejsce = tresc.tresc();
  const lista = wykaz('Wynik czynności repozytorium', 'mdev-wykaz');

  const pozycjaOgolna = pozycjaWykazu(
    wynik.action,
    opisWynikuOgolnego(wynik),
    'mdev',
  );
  lista.append(pozycjaOgolna.element);

  if (wynik.conflictPaths !== undefined && wynik.conflictPaths.length > 0) {
    for (const sciezka of wynik.conflictPaths) {
      const pozycja = pozycjaWykazu(sciezka, 'Ścieżka z konfliktem do rozstrzygnięcia.', 'mdev');
      pozycja.element.dataset['konflikt'] = 'tak';
      pozycja.akcje.append(przejscieDoEdytora(sciezka, naWskazanie));
      lista.append(pozycja.element);
    }
  }

  if (wynik.changedPaths !== undefined && wynik.changedPaths.length > 0) {
    for (const sciezka of wynik.changedPaths) {
      const pozycja = pozycjaWykazu(sciezka, 'Ścieżka objęta czynnością.', 'mdev');
      pozycja.akcje.append(przejscieDoEdytora(sciezka, naWskazanie));
      lista.append(pozycja.element);
    }
  }

  miejsce.append(lista);
}

/** Przycisk „Otwórz w edytorze” — wskazuje plik we wspólnym stanie modułu Developer, by otworzył go Code Editor. */
function przejscieDoEdytora(
  sciezka: string,
  naWskazanie: (sciezka: string) => void,
): HTMLButtonElement {
  const otworz = przycisk('Otwórz w edytorze', 'dn-btn dn-btn--zarys');
  otworz.dataset['sciezka'] = sciezka;
  otworz.addEventListener('click', () => naWskazanie(sciezka));
  return otworz;
}

/** Zdanie ogólne wyniku czynności repozytorium: gałąź, zatwierdzenie oraz wyjście wykonanego polecenia gita. */
function opisWynikuOgolnego(wynik: GitActionResult): string {
  const czesci: string[] = [wynik.succeeded ? 'powodzenie' : 'repozytorium odmówiło'];
  if (wynik.branch !== undefined) czesci.push(`gałąź: ${wynik.branch}`);
  if (wynik.commitId !== undefined) czesci.push(`zatwierdzenie: ${wynik.commitId}`);
  if (wynik.output !== undefined && wynik.output !== '') czesci.push(wynik.output);
  return czesci.join(' — ');
}
