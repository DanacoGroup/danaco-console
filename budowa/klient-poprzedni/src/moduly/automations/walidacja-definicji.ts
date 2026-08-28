import { AutomationStepKind, type AutomationStep } from '../../../../shared/contract';

/**
 * Waga zastrzeżenia definicji automatyki wykrytego w oknie Workflow Builder:
 * błąd sygnalizuje wadę uniemożliwiającą poprawne wykonanie, ostrzeżenie
 * sygnalizuje wadę niekrytyczną.
 */
export type WagaZastrzezenia = 'ostrzezenie' | 'blad';

/**
 * Jedno zastrzeżenie do definicji automatyki wraz z miejscem kroku, którego
 * dotyczy, wagą oraz zdaniem opisującym wadę dla operatora.
 */
export interface ZastrzezenieDefinicji {
  /** Miejsce kroku w kolejności wykonania liczy się od 1; wartość zero oznacza całą definicję. */
  miejsce: number;
  /** Identyfikator kroku jest pusty, gdy zastrzeżenie dotyczy całej definicji, a nie kroku. */
  idKroku: string;
  waga: WagaZastrzezenia;
  zdanie: string;
}

/**
 * Zastrzeżenia do wykazu kroków w kolejności wykonania.
 *
 * Wykaz pusty znaczy „definicja bez zastrzeżeń”, nie „nie sprawdzono”.
 */
export function zastrzezeniaDefinicji(
  kroki: readonly AutomationStep[],
): ZastrzezenieDefinicji[] {
  const zastrzezenia: ZastrzezenieDefinicji[] = [];
  const znane = new Set<string>();
  const powielone = new Set<string>();

  for (const krok of kroki) {
    const kod = krok.id.trim();
    if (kod === '') continue;
    if (znane.has(kod)) powielone.add(kod);
    znane.add(kod);
  }

  kroki.forEach((krok, miejsce) => {
    const numer = miejsce + 1;
    const kod = krok.id.trim();
    if (kod === '') {
      zastrzezenia.push({
        miejsce: numer,
        idKroku: '',
        waga: 'blad',
        zdanie: 'Krok bez identyfikatora — Orchestrator nie ma czym ustalić jego zależności.',
      });
    } else if (powielone.has(kod)) {
      zastrzezenia.push({
        miejsce: numer,
        idKroku: kod,
        waga: 'blad',
        zdanie: `Identyfikator ${kod} niesie więcej niż jeden krok — zależność wskaże wtedy krok nieokreślony.`,
      });
    }

    if (krok.kind === AutomationStepKind.Command && (krok.command ?? '').trim() === '') {
      zastrzezenia.push({
        miejsce: numer,
        idKroku: kod,
        waga: 'ostrzezenie',
        zdanie: 'Krok rodzaju „komenda” bez nazwy komendy — nie ma czego wywołać.',
      });
    }
    if (
      (krok.kind === AutomationStepKind.Condition || krok.kind === AutomationStepKind.Branch) &&
      (krok.condition ?? '').trim() === ''
    ) {
      zastrzezenia.push({
        miejsce: numer,
        idKroku: kod,
        waga: 'ostrzezenie',
        zdanie: 'Krok rozstrzygający bez warunku — rozstrzygnięcie nie ma na czym stanąć.',
      });
    }
    if (
      krok.kind === AutomationStepKind.Wait &&
      (krok.condition ?? '').trim() === '' &&
      krok.params === undefined
    ) {
      zastrzezenia.push({
        miejsce: numer,
        idKroku: kod,
        waga: 'ostrzezenie',
        zdanie: 'Krok oczekiwania bez warunku i bez treści żądania — nie wiadomo, na co czeka.',
      });
    }

    for (const poprzednik of krok.dependsOn ?? []) {
      if (poprzednik.trim() === '' || znane.has(poprzednik.trim())) continue;
      zastrzezenia.push({
        miejsce: numer,
        idKroku: kod,
        waga: 'blad',
        zdanie: `Krok zależy od ${poprzednik}, a kroku o tym identyfikatorze nie ma w definicji.`,
      });
    }
  });

  zastrzezenia.push(...zastrzezeniaSpojnosci(kroki, znane));
  return zastrzezenia;
}

/**
 * Zastrzeżenia dotyczące całości definicji: cykl w zależnościach oraz krok
 * bez połączenia. Krok bez połączenia jest zgłaszany tylko wtedy, gdy
 * definicja używa zależności.
 */
function zastrzezeniaSpojnosci(
  kroki: readonly AutomationStep[],
  znane: ReadonlySet<string>,
): ZastrzezenieDefinicji[] {
  const zastrzezenia: ZastrzezenieDefinicji[] = [];
  const poprzednicy = new Map<string, string[]>();
  const nastepnicy = new Set<string>();
  let zaleznosciWDefinicji = 0;

  for (const krok of kroki) {
    const kod = krok.id.trim();
    const wskazane = (krok.dependsOn ?? [])
      .map((poprzednik) => poprzednik.trim())
      .filter((poprzednik) => znane.has(poprzednik));
    zaleznosciWDefinicji += wskazane.length;
    poprzednicy.set(kod, wskazane);
    for (const poprzednik of wskazane) nastepnicy.add(poprzednik);
  }

  const wCyklu = krokiWCyklu(poprzednicy);
  kroki.forEach((krok, miejsce) => {
    const kod = krok.id.trim();
    if (wCyklu.has(kod)) {
      zastrzezenia.push({
        miejsce: miejsce + 1,
        idKroku: kod,
        waga: 'blad',
        zdanie: 'Krok stoi w cyklu zależności — przebieg nie miałby od czego zacząć.',
      });
      return;
    }
    if (
      zaleznosciWDefinicji > 0 &&
      kroki.length > 1 &&
      kod !== '' &&
      (poprzednicy.get(kod) ?? []).length === 0 &&
      !nastepnicy.has(kod)
    ) {
      zastrzezenia.push({
        miejsce: miejsce + 1,
        idKroku: kod,
        waga: 'ostrzezenie',
        zdanie: 'Krok bez połączenia — nie zależy od żadnego kroku i żaden nie zależy od niego.',
      });
    }
  });
  return zastrzezenia;
}

/**
 * Kroki leżące na cyklu zależności, wyznaczone przeglądem w głąb ze
 * znacznikiem odwiedzin. Kroki cyklu wracają kompletem, bo sygnalizacja
 * dotyczy każdego z nich.
 */
function krokiWCyklu(poprzednicy: ReadonlyMap<string, readonly string[]>): Set<string> {
  const wCyklu = new Set<string>();
  const zamkniete = new Set<string>();
  const sciezka: string[] = [];
  const naSciezce = new Set<string>();

  function przejdz(kod: string): void {
    if (naSciezce.has(kod)) {
      const poczatek = sciezka.indexOf(kod);
      for (const wezel of sciezka.slice(poczatek)) wCyklu.add(wezel);
      return;
    }
    if (zamkniete.has(kod)) return;
    naSciezce.add(kod);
    sciezka.push(kod);
    for (const poprzednik of poprzednicy.get(kod) ?? []) przejdz(poprzednik);
    sciezka.pop();
    naSciezce.delete(kod);
    zamkniete.add(kod);
  }

  for (const kod of poprzednicy.keys()) przejdz(kod);
  return wCyklu;
}

/**
 * Zdanie podsumowania walidacji definicji, stanowiące nagłówek wykazu
 * zastrzeżeń w oknie Workflow Builder dla operatora.
 */
export function zdanieWalidacji(zastrzezenia: readonly ZastrzezenieDefinicji[]): string {
  if (zastrzezenia.length === 0) {
    return 'Definicja bez zastrzeżeń okna: każdy krok ma identyfikator, treść właściwą rodzajowi i istniejących poprzedników.';
  }
  const bledy = zastrzezenia.filter((zastrzezenie) => zastrzezenie.waga === 'blad').length;
  const ostrzezenia = zastrzezenia.length - bledy;
  return (
    `Zastrzeżenia okna do definicji: ${bledy} poważnych, ${ostrzezenia} ostrzegawczych. ` +
    'Zapis pozostaje możliwy — walidacja nie jest bramą.'
  );
}
