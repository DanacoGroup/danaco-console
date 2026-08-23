import type { HistoryEntry, RetentionPolicy, RetentionSetRequest } from '../../../shared/contract';
import { przyciskAkcji, poleLiczbowe, wybor } from '../modele/kontrolki-formularza-braki';
import { utworzMagistrale, type Odsubskrybuj } from '../polaczenie/magistrala-zdarzen';

/**
 * Nastawa zasady przechowywania historii — progi, zakres i zapowiedź skutku.
 *
 * Rdzeń zasadę egzekwuje, nie tylko zapisuje: przemiata przy swoim starcie
 * (`server/internal/core/trwalosc_kosza.go`, wpięte w `core/montaz.go`), przy
 * każdym `history.load` oraz zaraz po samym `retention.set`
 * (`handlers_historia.go`). Po zapisie pozycje przestają istnieć, więc formularz
 * pokazuje skutek zasady, zanim Operator ją zapisze.
 *
 * Zapowiedź nie jest bramką: nie ma tu pytania „czy na pewno", a przycisk
 * zapisu jest czynny zawsze. Zapowiedź daje wiedzę, nie zgodę — Operator widzi
 * liczbę pozycji, które wypadną, i naciska albo nie.
 *
 * Oba progi puste znaczą „bez ograniczenia" i tak zdejmuje się zasadę
 * z zakresu; kontrakt nie ma osobnej komendy kasującej, bo dwie drogi do
 * jednego skutku byłyby dwiema prawdami o retencji. Formularz mówi to wprost,
 * bo bez tego zdania puste pola czytają się jak „nie ruszaj".
 *
 * Rdzeń przyjmuje zakresy `window`, `session` i `global`, przy czym dwa
 * pierwsze wymagają wskazania bytu (`sprawdzZakresZasady`). Okno panel zna od
 * gospodarza; sesji nie zna — `OpcjePanelu` jej nie niesie. Sesja pada wyłącznie
 * we wczytanych pozycjach (`HistoryEntry.sessionId`), więc zakres sesji wchodzi
 * do wyboru dopiero wtedy, gdy panel tę sesję zobaczył. Pozycja obiecująca
 * zakres, którego nie ma czym wypełnić, byłaby atrapą.
 */
export interface NastawaRetencji {
  /** Formularz osadzany w panelu. */
  element: HTMLElement;
  /** Identyfikatory pozycji, których nastawiona zasada NIE utrzyma. */
  pozaZasada(pozycje: readonly HistoryEntry[]): ReadonlySet<string>;
  /** Podaje sesję odczytaną z pozycji; pusty napis = sesji nie widziano. */
  ustawSesje(kod: string): void;
  /** Pokazuje zasadę oddaną przez rdzeń po zapisie. */
  pokazZasade(zasada: RetentionPolicy): void;
  /** Zgłasza zmianę progów — panel przerysowuje zapowiedź. */
  naZmiane(sluchacz: () => void): Odsubskrybuj;
}

export interface OpcjeNastawy {
  /** Okno gospodarza — byt zakresu `window`. Pusty napis = rdzeń nie dał okna. */
  okno: string;
  /** Wywoływane po naciśnięciu zapisu; żądanie jest już złożone. */
  naZapis(zadanie: RetentionSetRequest): void;
}

/** Doba w milisekundach — jednostka progu `keepDays` kontraktu. */
const DOBA_MS = 24 * 60 * 60 * 1000;

export function utworzNastawaRetencji(opcje: OpcjeNastawy): NastawaRetencji {
  const zmiany = utworzMagistrale<void>();
  let sesja = '';

  const zakres = wybor('Zakres zasady przechowywania', zakresyDostepne(opcje.okno, sesja));
  const dni = poleLiczbowe('Ile dni trzymać historię', 'bez ograniczenia');
  const pozycje = poleLiczbowe('Ile pozycji trzymać', 'bez ograniczenia');

  const zdanie = document.createElement('p');
  zdanie.className = 'dnp-retencja__zdanie';
  zdanie.setAttribute('role', 'note');
  zdanie.textContent = ZDANIE_PUSTEJ;

  const zapisz = przyciskAkcji('Zapisz zasadę', 'dn-btn dn-btn--atrament');
  zapisz.addEventListener('click', () => opcje.naZapis(zlozZadanie(zakres.value, dni.value, pozycje.value, opcje.okno, sesja)));

  for (const kontrolka of [zakres, dni, pozycje]) {
    kontrolka.addEventListener('input', () => zmiany.oglos());
  }

  const rzad = document.createElement('div');
  rzad.className = 'dnp-retencja__rzad';
  rzad.append(zakres, dni, pozycje, zapisz);

  const element = document.createElement('section');
  element.className = 'dnp-retencja';
  const tytul = document.createElement('h4');
  tytul.className = 'dnp-retencja__tytul';
  tytul.textContent = 'Zasada przechowywania';
  element.append(tytul, rzad, zdanie);

  return {
    element,

    pozaZasada: (lista) => policzPozaZasada(lista, dni.value, pozycje.value),

    ustawSesje(kod) {
      if (kod === sesja) return;
      sesja = kod;
      // Wybór przebudowujemy w miejscu, zachowując zaznaczenie: zakres wybrany
      // przez Operatora nie ma się przestawiać dlatego, że doszła pozycja
      // z sesją. Zakres nieobecny w nowym wykazie po prostu nie wróci.
      const wybrany = zakres.value;
      zakres.replaceChildren();
      for (const [wartosc, opis] of zakresyDostepne(opcje.okno, sesja)) {
        const opcja = document.createElement('option');
        opcja.value = wartosc;
        opcja.textContent = opis;
        zakres.append(opcja);
      }
      zakres.value = wybrany;
      // Zakres, który wypadł z wykazu, nie może zostać w polu jako pustka:
      // `select.value` przy wartości nieobecnej daje `selectedIndex === -1`,
      // pole wygląda na niewypełnione, a zapis idzie do rdzenia z `scope: ''`.
      // Zakres `global` jest w wykazie zawsze i bytu nie wymaga, więc to on
      // jest stanem, do którego wybór wraca.
      if (zakres.selectedIndex < 0) zakres.value = 'global';
    },

    pokazZasade(zasada) {
      zdanie.textContent = opisZasady(zasada);
    },

    naZmiane: (sluchacz) => zmiany.subskrybuj(sluchacz),
  };
}

/**
 * Zakresy, które da się wypełnić. Okno bez kodu odpada razem z zakresem
 * `window`, bo rdzeń odmówiłby zapisu bez wskazania bytu — pozycja wiodąca
 * wprost do odmowy jest gorsza niż jej brak.
 */
function zakresyDostepne(okno: string, sesja: string): Array<readonly [string, string]> {
  const lista: Array<readonly [string, string]> = [];
  if (okno.trim() !== '') lista.push(['window', 'To okno']);
  if (sesja.trim() !== '') lista.push(['session', 'Cała sesja tego okna']);
  lista.push(['global', 'Cała instalacja']);
  return lista;
}

/** Żądanie `retention.set`; progu pustego nie wysyła — brak znaczy brak ograniczenia. */
function zlozZadanie(
  zakres: string,
  dni: string,
  pozycje: string,
  okno: string,
  sesja: string,
): RetentionSetRequest {
  const zadanie: RetentionSetRequest = { scope: zakres };
  if (zakres === 'window' && okno !== '') zadanie.scopeId = okno;
  if (zakres === 'session' && sesja !== '') zadanie.scopeId = sesja;
  const prog = (tekst: string): number | null => {
    const liczba = Number.parseInt(tekst, 10);
    return Number.isFinite(liczba) && liczba > 0 ? liczba : null;
  };
  const wDniach = prog(dni);
  const wPozycjach = prog(pozycje);
  if (wDniach !== null) zadanie.keepDays = wDniach;
  if (wPozycjach !== null) zadanie.keepEntries = wPozycjach;
  return zadanie;
}

/**
 * Pozycje, których zasada nie utrzyma — liczone na tym, co panel wczytał.
 *
 * Panel widzi stronę wykazu, nie całą historię okna, a rdzeń liczy próg
 * `keepEntries` na całości. Zapowiedź jest więc dolnym oszacowaniem: pozycji
 * poza zasadą będzie co najmniej tyle, ile tu widać, i tak też brzmi zdanie
 * panelu.
 *
 * Kolejność jest kolejnością `history.load` — od najnowszej — więc numer
 * pozycji w tablicy jest jej numerem od najnowszej i to on wchodzi w próg.
 */
function policzPozaZasada(
  lista: readonly HistoryEntry[],
  dni: string,
  pozycje: string,
): ReadonlySet<string> {
  const poza = new Set<string>();
  const wDniach = Number.parseInt(dni, 10);
  const wPozycjach = Number.parseInt(pozycje, 10);
  const granicaCzasu =
    Number.isFinite(wDniach) && wDniach > 0 ? Date.now() - wDniach * DOBA_MS : null;

  lista.forEach((pozycja, numer) => {
    if (granicaCzasu !== null && pozycja.createdAt < granicaCzasu) poza.add(pozycja.id);
    if (Number.isFinite(wPozycjach) && wPozycjach > 0 && numer >= wPozycjach) poza.add(pozycja.id);
  });
  return poza;
}

/** Zdanie o zasadzie pustej — stan wyjściowy formularza i wynik zdjęcia zasady. */
const ZDANIE_PUSTEJ =
  'Oba progi puste znaczą BRAK OGRANICZENIA — zapisanie takiej zasady zdejmuje ją z zakresu. ' +
  'Zasada nie jest notatką: rdzeń egzekwuje ją od razu po zapisie, przy każdym odczycie ' +
  'historii i przy swoim starcie.';

/** Zasada oddana przez rdzeń, przepisana na zdanie Operatora. */
function opisZasady(zasada: RetentionPolicy): string {
  const progi: string[] = [];
  if (zasada.keepDays !== undefined) progi.push(`${zasada.keepDays} dni`);
  if (zasada.keepEntries !== undefined) progi.push(`${zasada.keepEntries} pozycji`);
  const zakres = zasada.scopeId === undefined || zasada.scopeId === ''
    ? zasada.scope
    : `${zasada.scope} ${zasada.scopeId}`;
  if (progi.length === 0) {
    return `Zasada zakresu ${zakres} jest pusta — historia nie ma ograniczenia i nic z niej nie ubywa.`;
  }
  return `Zasada zakresu ${zakres} obowiązuje: rdzeń trzyma ${progi.join(' i ')}. ` +
    'Egzekucja poszła już na oknach tego zakresu.';
}
