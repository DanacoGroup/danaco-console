import { type RoundtableParticipant } from '../../../../shared/contract';
import { utworzRameOkna } from '../../komponenty/rama-okna';
import {
  pole,
  poleTresci,
  przestaw,
  przelacznikWidoku,
  przyciskAkcji as przycisk,
  wiersz,
  wybor,
} from '../../modele/kontrolki-formularza-braki';
import { kanalyDebaty, zdaniaKanalu, ZDANIE_GRANICY_DIAGNOSTYKI } from './diagnostyka-kanalow';
import { czynnosciSkladu } from './czynnosci-arsenalu';
import { rysujSkladModeli } from './sklad-modeli';
import type { StanDebaty } from './stan-debaty';
import { nazwaUczestnika } from './stan-debaty';
import { utworzStanTresci } from './stany-okna';
import type { StrumienWypowiedzi } from './strumien-wypowiedzi';
import { utworzNoteSufitu } from './sufit-uczestnikow';
import type { ZrodloArsenaluRoundtable } from './zrodlo-arsenalu';
import { utworzWarstweTresci, utworzWykazFunkcji, utworzZestawAkcji } from './warstwy-modulu';
import type { ZrodloRoundtable } from './zrodlo-roundtable';

/**
 * Model Panels — okno wiodące modułu Roundtable: wybór uczestniczących modeli
 * i osobny panel odpowiedzi na każdy z nich, aktualizowany jednocześnie po
 * zadaniu pytania.
 *
 * Skład idzie z `stan.uczestnicy()`, nie z własnej listy. Stan subskrybuje
 * `roundtable.debate.changed` (`stan-debaty.ts`) i dokłada uczestników z
 * przyrostu tury oraz z czynności moderatora — to okno tylko czyta i rysuje.
 * Odpowiedzi rosną na żywo z `StrumienWypowiedzi` (fragmenty `stream.chunk`
 * okna debaty), a rysowanie składu wraz z odpowiedziami leży w
 * `sklad-modeli.ts`. Subskrypcji okno nie zakłada: jedna na całe złożenie stoi
 * w `indeks.ts`.
 *
 * Dwie tożsamości jednego kanału dają dwa panele. Klucz panelu to
 * `participantId`, nigdy `channelId` (tak mówi kontrakt `roundtable.model.add`).
 * `stan.wystapieniaKanalu(idKanalu)` mówi, ile razy kanał już wystąpił — od
 * drugiego wystąpienia wzwyż panel niesie `data-kanal-powtorzony='tak'`, żeby
 * dwa identyczne z pozoru panele nie wyglądały na usterkę.
 *
 * Świeżo otwarte okno nie zna składu debaty. Odczyt jest w kontrakcie —
 * `roundtable.model.list` oddaje skład, a `roundtable.debate.get` całą debatę —
 * ale to okno żadnego z nich jeszcze nie wywołuje: obsługi nie zbudowano.
 * Do tego czasu skład narasta z przyrostu zdarzenia i z `participants`
 * zwracanego przez `moderator.direct`. Stan pusty na starcie jest stanem
 * poprawnym, opisanym wprost, nie udawaną pustką.
 *
 * Kontrakt nie niesie górnego limitu instancji Model Panels i okno nie narzuca
 * własnego — skład rośnie tak, jak rośnie w rdzeniu. Ilu uczestników mieści
 * debata, mówi nota nad formularzem, a nie bramka: uczestnik ponad zwyczajowy
 * skład przechodzi przez rdzeń bez odmowy (`sufit-uczestnikow.ts`), więc
 * przycisk „Dodaj uczestnika" nie zna żadnego progu. Nota podaje, ilu
 * uczestników debata ma, czego nie liczy rdzeń ani kontrakt i gdzie leży sufit
 * gniazd sceny okien równoległych.
 */
export interface OknoModelPanels {
  element: HTMLElement;
  odswiez(): void;
  /**
   * Przerysowanie po fragmencie strumienia — wołane przez złożenie modułu.
   *
   * Idzie na każdy fragment, więc omija stan błędu i zapowiedź odczytu:
   * inaczej zdanie „Uczestnik nie został dodany do debaty" znikałoby przy
   * pierwszym słowie następnego mówcy.
   */
  odswiezGlosy(): void;
  /** Zamyka nasłuch `stan.naZmiane(...)`. */
  zamknij(): void;
}

export function utworzOknoModelPanels(
  zrodlo: ZrodloRoundtable,
  arsenal: ZrodloArsenaluRoundtable,
  stan: StanDebaty,
  strumien: StrumienWypowiedzi,
  idOkna: string,
): OknoModelPanels {
  stan.ustawOkno(idOkna);

  const rama = utworzRameOkna({
    tytul: 'Model Panels',
    rola: 'wiodące',
    kod: 'model-panels',
    przeznaczenie:
      'Wybór uczestniczących modeli; osobny panel odpowiedzi na każdy, aktualizowany jednocześnie po zadaniu pytania.',
    przedrostek: 'dr',
  });
  const tresc = utworzStanTresci();
  // Czynności arsenału powstają razem z powierzchnią, bo część z nich siedzi
  // w pasku akcji, a część w zestawie warstwy trzeciej. Meldunek idzie tą samą
  // drogą co reszta odpowiedzi okna — stanem treści, nie osobnym dymkiem.
  const powierzchnia = zlozPowierzchnieModelPanels(rama, tresc.element, (poleZespolu) =>
    czynnosciSkladu(
      arsenal,
      stan,
      (zdanie, powodzenie) => tresc.potwierdzenie(zdanie, powodzenie),
      () => {
        const wybrany = stan.uczestnicy()[0];
        return wybrany === undefined ? '' : wybrany.id;
      },
      () => poleZespolu.value,
    ),
  );
  let notaSufitu = powierzchnia.notaSufitu;
  // Nastawy widoku żyją w oknie, bo kontrakt ich nie zna: ani anonimizacji, ani
  // licznika odpowiedzi nie niesie żadne pole obszaru roundtable.
  let trybSlepy = false;
  let metryki = false;

  // Nota o suficie stoi poza miejscem treści, bo `pusto(...)` i `blad(...)`
  // czyszczą je do zera — a stan „debata nie ma jeszcze nikogo" jest właśnie tym
  // stanem, w którym Operator pyta, ilu uczestników w ogóle wolno mu dodać.
  function odswiezSufit(): void {
    const nowa = utworzNoteSufitu(stan.uczestnicy().length);
    notaSufitu.replaceWith(nowa);
    notaSufitu = nowa;
  }

  function rysuj(): void {
    odswiezSufit();
    odswiezDiagnostyke(powierzchnia.diagnostyka, stan);
    rysujSkladModeli(stan, strumien, tresc, { trybSlepy, metryki });
  }

  function odswiezGlosy(): void {
    const rodzaj = tresc.rodzaj();
    if (rodzaj === 'blad' || rodzaj === 'ladowanie') return;
    rysuj();
  }

  function dodajUczestnika(): void {
    const idKanalu = powierzchnia.wyborKanalu.value;
    if (idKanalu === '') {
      tresc.blad('Wybierz kanał modelu, zanim dodasz uczestnika.');
      return;
    }
    const zadanie: Parameters<ZrodloRoundtable['dodajModel']>[0] = {
      windowId: idOkna,
      channelId: idKanalu,
    };
    const persona = powierzchnia.polePersony.value.trim();
    if (persona !== '') zadanie.personaName = persona;
    const prompt = powierzchnia.poleSystemowe.value.trim();
    if (prompt !== '') zadanie.systemPrompt = prompt;

    tresc.ladowanie('Dodawanie uczestnika…');
    void zrodlo.dodajModel(zadanie).then((wynik) => {
      if (!wynik.udany || wynik.wynik === undefined) {
        // Zdanie nie mówi „rdzeń odmówił": `sprawdzKsztalt` oddaje niepowodzenie
        // także wtedy, gdy rdzeń odpowiedział bez pola obowiązkowego. Powód
        // niesie treść błędu.
        tresc.blad('Uczestnik nie został dodany do debaty.', wynik.blad);
        return;
      }
      stan.dodajUczestnika(wynik.wynik);
      powierzchnia.polePersony.value = '';
      powierzchnia.poleSystemowe.value = '';
      rysuj();
      const potwierdzenie = zdanieDodania(wynik.wynik, persona, stan);
      tresc.potwierdzenie(potwierdzenie.zdanie, potwierdzenie.udane);
    });
  }

  podepnijAkcjeModelPanels(powierzchnia, {
    dodajUczestnika,
    // Etykieta niesie stan słowem, bo `aria-pressed` czyta wyłącznie technologia
    // wspomagająca, a sama zmiana obrysu nie mówi widzącemu, co jest włączone.
    przelaczTrybSlepy: () => {
      trybSlepy = przestaw(powierzchnia.trybSlepy);
      powierzchnia.trybSlepy.textContent = trybSlepy
        ? 'Tryb ślepy: włączony'
        : 'Tryb ślepy: wyłączony';
      rysuj();
    },
    przelaczMetryki: () => {
      metryki = przestaw(powierzchnia.metryki);
      powierzchnia.metryki.textContent = metryki
        ? 'Metryki odpowiedzi: widoczne'
        : 'Metryki odpowiedzi: ukryte';
      rysuj();
    },
  });

  const odsubskrybuj = stan.naZmiane(() => {
    odswiezWyborKanalow(powierzchnia.wyborKanalu, stan);
    rysuj();
  });

  // Katalog kanałów jest wspólny dla okien Roundtable — zamawia go
  // złożenie modułu przed montażem okien, nie każde okno z osobna. Odczyt
  // wystarczy przeczytać ze stanu; drugie zamówienie zdublowałoby kopertę
  // `channel.list` na starcie.
  odswiezWyborKanalow(powierzchnia.wyborKanalu, stan);
  rysuj();

  return { element: rama.element, odswiez: rysuj, odswiezGlosy, zamknij: odsubskrybuj };
}

/**
 * Zdanie potwierdzenia dodania uczestnika, składane z odpowiedzi rdzenia.
 *
 * Tożsamość, kanał i miejsce w kolejności głosu nadaje rdzeń
 * (`RoundtableParticipant.id`, `.channelId`, `.order`), więc to one wchodzą do
 * zdania — nie to, co okno wysłało. Gdy rdzeń nie zapisał podanej tożsamości,
 * okno mówi to wprost zamiast potwierdzać czynność w kształcie, którego nie
 * było.
 */
function zdanieDodania(
  uczestnik: RoundtableParticipant,
  zadanaPersona: string,
  stan: StanDebaty,
): { zdanie: string; udane: boolean } {
  const nazwa = nazwaUczestnika(uczestnik, stan.opisKanalu(uczestnik.channelId));
  const miejsce =
    uczestnik.order === undefined
      ? 'bez miejsca w kolejności głosu'
      : `miejsce #${uczestnik.order} w kolejności głosu`;
  const oddanaPersona = uczestnik.personaName ?? '';
  if (zadanaPersona !== '' && oddanaPersona !== zadanaPersona) {
    return {
      zdanie:
        `Rdzeń dodał uczestnika ${nazwa} (${miejsce}), ale tożsamości „${zadanaPersona}” ` +
        `nie zapisał — w odpowiedzi stoi ${oddanaPersona === '' ? 'brak tożsamości' : `„${oddanaPersona}”`}.`,
      udane: false,
    };
  }
  return { zdanie: `Rdzeń dodał uczestnika: ${nazwa} — ${miejsce}.`, udane: true };
}

/** Kontrolki formularza dodania uczestnika i nastaw widoku okna. */
interface FormularzUczestnika {
  wyborKanalu: HTMLSelectElement;
  polePersony: HTMLInputElement;
  poleSystemowe: HTMLTextAreaElement;
  dodaj: HTMLButtonElement;
  trybSlepy: HTMLButtonElement;
  metryki: HTMLButtonElement;
  /** Nota o suficie — wymieniana przy każdym rysowaniu, patrz `odswiezSufit`. */
  notaSufitu: HTMLElement;
  /** Wnętrze warstwy diagnostyki kanałów — przerysowywane przy każdej zmianie. */
  diagnostyka: HTMLElement;
}

/**
 * Składa formularz dodania uczestnika i pasek akcji ramy.
 *
 * Fragment jest czystą konstrukcją — nie domyka się na stanie modułu.
 *
 * Dwie pozycje paska mają komendę w kontrakcie i nie mają jeszcze obsługi:
 * usunięcie uczestnika (`roundtable.model.remove`) oraz oznaczenie uczestnika
 * jako kluczowego (`roundtable.model.update`, pole `key`). Wyciszenia w pasku
 * nie ma wcale: jest czynnością moderatora (`ModeratorAction`), więc należy do
 * Moderator Panelu, nie do Model Panels. Pozostałe czynności uczestnika bez
 * obsługi siedzą w zestawie akcji warstwy trzeciej (`zlozZestawUczestnika`).
 */
function zlozAkcjeModelPanels(gospodarz: HTMLElement): {
  wyborKanalu: HTMLSelectElement;
  dodaj: HTMLButtonElement;
  poleZespolu: HTMLInputElement;
} {
  const wyborKanalu = wybor('Kanał modelu', []);
  const dodaj = przycisk('Dodaj uczestnika', 'dn-btn dn-btn--atrament');
  const poleZespolu = pole('Nazwa zespołu', 'nazwa zapisywanego składu');

  gospodarz.append(wyborKanalu, dodaj, poleZespolu);
  return { wyborKanalu, dodaj, poleZespolu };
}

/**
 * Zestaw akcji warstwy trzeciej — czynności uczestnika jeszcze niezbudowane.
 *
 * Pozycje siedzą w zestawie zwiniętym, a nie w pasku akcji okna: pasek zostaje
 * przy dodaniu uczestnika, czyli jedynej czynności okna, którą to okno wykonuje.
 * Siedem pozycji ma dziś komendę w kontrakcie; ósma — liczba wariantów jednej
 * persony — ma pole w uczestniku i nie ma komendy, która by je zapisała.
 */
function zlozZestawUczestnika(czynnosci: HTMLButtonElement[]): HTMLElement {
  return utworzZestawAkcji('Zestaw akcji uczestnika', czynnosci);
}

/**
 * Przerysowanie warstwy diagnostyki kanałów.
 *
 * Diagnostyka idzie po składzie debaty, więc zmienia się z każdym dodanym
 * uczestnikiem i z każdym odczytem rejestru kanałów. Wymieniane jest wnętrze,
 * a nie cała warstwa: warstwa otwarta przez Operatora ma zostać otwarta.
 */
function odswiezDiagnostyke(wnetrze: HTMLElement, stan: StanDebaty): void {
  wnetrze.replaceChildren();
  const kanaly = kanalyDebaty(stan);
  if (kanaly.length === 0) {
    const pusto = document.createElement('p');
    pusto.className = 'dr-sklad__nota';
    pusto.textContent = 'Debata nie ma jeszcze uczestników, więc żaden kanał modelu w niej nie mówi.';
    wnetrze.append(pusto);
  } else {
    const wykaz = document.createElement('ul');
    wykaz.className = 'dr-wezly';
    wykaz.setAttribute('aria-label', 'Kanały modeli uczestniczących w debacie');
    for (const kanal of kanaly) {
      const pozycja = document.createElement('li');
      pozycja.className = 'dr-wezel';
      pozycja.dataset['kanalCzynny'] = kanal.wpis === null || !kanal.wpis.enabled ? 'nie' : 'tak';
      for (const zdanie of zdaniaKanalu(kanal, stan.kanalyOdczytane())) {
        const opis = document.createElement('span');
        opis.className = 'dr-wezel__meta';
        opis.textContent = zdanie;
        pozycja.append(opis);
      }
      wykaz.append(pozycja);
    }
    wnetrze.append(wykaz);
  }
  const granica = document.createElement('p');
  granica.className = 'dr-sklad__nota';
  granica.dataset['brakSygnalu'] = 'tak';
  granica.textContent = ZDANIE_GRANICY_DIAGNOSTYKI;
  wnetrze.append(granica);
}

/** Składa formularz tożsamości uczestnika (persona, prompt systemowy) i ciało ramy. */
function zlozPowierzchnieModelPanels(
  rama: { akcje: HTMLElement; narzedzia: HTMLElement; cialo: HTMLElement },
  stanTresci: HTMLElement,
  czynnosci: (poleZespolu: HTMLInputElement) => HTMLButtonElement[],
): FormularzUczestnika {
  const { wyborKanalu, dodaj, poleZespolu } = zlozAkcjeModelPanels(rama.akcje);
  const polePersony = pole('Nazwa tożsamości', 'opcjonalnie — np. Oponent, Obrońca');
  const poleSystemowe = poleTresci('Prompt systemowy tożsamości', 3, 'opcjonalnie');
  const notaSufitu = utworzNoteSufitu(0);
  const trybSlepy = przelacznikWidoku('Tryb ślepy: wyłączony', false);
  const metryki = przelacznikWidoku('Metryki odpowiedzi: ukryte', false);
  const diagnostyka = utworzWarstweTresci(
    'Diagnostyka kanałów modeli uczestniczących w debacie',
    'diagnostyka',
  );

  rama.narzedzia.append(trybSlepy, metryki);
  rama.cialo.append(
    notaSufitu,
    wiersz('Tożsamość nowego uczestnika', polePersony, {
      klasa: 'dr-tozsamosc',
      objasnienie: 'Ten sam kanał może wystąpić dwukrotnie pod odrębnymi tożsamościami.',
    }),
    wiersz('Prompt systemowy', poleSystemowe, { klasa: 'dr-tozsamosc' }),
    zlozZestawUczestnika(czynnosci(poleZespolu)),
    stanTresci,
    diagnostyka.element,
    utworzWykazFunkcji('model-panels'),
  );
  return {
    wyborKanalu,
    polePersony,
    poleSystemowe,
    dodaj,
    trybSlepy,
    metryki,
    notaSufitu,
    diagnostyka: diagnostyka.wnetrze,
  };
}

/** Podpina pasek akcji i nastawy widoku do czynności okna. */
function podepnijAkcjeModelPanels(
  powierzchnia: FormularzUczestnika,
  obsluga: {
    dodajUczestnika: () => void;
    przelaczTrybSlepy: () => void;
    przelaczMetryki: () => void;
  },
): void {
  powierzchnia.dodaj.addEventListener('click', obsluga.dodajUczestnika);
  powierzchnia.trybSlepy.addEventListener('click', obsluga.przelaczTrybSlepy);
  powierzchnia.metryki.addEventListener('click', obsluga.przelaczMetryki);
}

/**
 * Wypełnia listę wyboru kanałów z rejestru rdzenia.
 *
 * `stan.kanalyOdczytane() === false` znaczy „jeszcze nie wiem", nie „rejestr
 * pusty" — dwa stany mają wyglądać różnie, więc pozycja zastępcza mówi to
 * wprost zamiast zostawiać listę milczącą.
 */
function odswiezWyborKanalow(wyborKanalu: HTMLSelectElement, stan: StanDebaty): void {
  const zaznaczony = wyborKanalu.value;
  wyborKanalu.replaceChildren();

  if (!stan.kanalyOdczytane()) {
    const oczekiwanie = document.createElement('option');
    oczekiwanie.value = '';
    oczekiwanie.textContent = 'Odczyt kanałów w toku…';
    wyborKanalu.append(oczekiwanie);
    return;
  }

  const kanaly = stan.kanaly();
  const pusta = document.createElement('option');
  pusta.value = '';
  pusta.textContent = kanaly.length === 0 ? 'Rejestr kanałów jest pusty' : 'Wybierz kanał…';
  wyborKanalu.append(pusta);

  for (const kanal of kanaly) {
    const pozycja = document.createElement('option');
    pozycja.value = kanal.id;
    pozycja.textContent = stan.opisKanalu(kanal.id) === '' ? kanal.id : stan.opisKanalu(kanal.id);
    wyborKanalu.append(pozycja);
  }
  if (kanaly.some((kanal) => kanal.id === zaznaczony)) wyborKanalu.value = zaznaczony;
}
