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

/** Model Panels jest oknem wiodącym modułu Roundtable: wybór uczestniczących modeli wraz z osobnym panelem odpowiedzi na każdego z nich. */
export interface OknoModelPanels {
  element: HTMLElement;
  odswiez(): void;
  /** Przerysowanie po fragmencie strumienia jest wołane przez złożenie modułu na każdy fragment. */
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
  // Czynności arsenału powstają razem z powierzchnią, częściowo w pasku, częściowo w warstwie trzeciej.
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
  // Nastawy widoku żyją w oknie, bo kontrakt roundtable nie zna anonimizacji ani licznika odpowiedzi.
  let trybSlepy = false;
  let metryki = false;

  // Nota o suficie stoi poza miejscem treści, bo czyszczenie treści zerowałoby ją razem ze stanem błędu.
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
        // Zdanie nie mówi wprost o odmowie rdzenia, bo powód niepowodzenia niesie treść błędu odpowiedzi.
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
    // Etykieta niesie stan słowem, bo technologia wspomagająca nie czyta samej zmiany obrysu przycisku.
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

  // Katalog kanałów zamawia złożenie modułu przed montażem okien, nie każde okno z osobna.
  odswiezWyborKanalow(powierzchnia.wyborKanalu, stan);
  rysuj();

  return { element: rama.element, odswiez: rysuj, odswiezGlosy, zamknij: odsubskrybuj };
}

/** Zdanie potwierdzenia dodania uczestnika składa się z odpowiedzi rdzenia, bo to rdzeń nadaje tożsamość, kanał i miejsce w kolejności głosu. */
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

/** Kontrolki formularza dodania uczestnika i nastaw widoku okna gromadzą pola oraz przyciski powierzchni. */
interface FormularzUczestnika {
  wyborKanalu: HTMLSelectElement;
  polePersony: HTMLInputElement;
  poleSystemowe: HTMLTextAreaElement;
  dodaj: HTMLButtonElement;
  trybSlepy: HTMLButtonElement;
  metryki: HTMLButtonElement;
  /** Nota o suficie jest wymieniana nowym elementem przy każdym rysowaniu składu debaty. */
  notaSufitu: HTMLElement;
  /** Wnętrze warstwy diagnostyki kanałów — przerysowywane przy każdej zmianie. */
  diagnostyka: HTMLElement;
}

/** Funkcja składa formularz dodania uczestnika i pasek akcji ramy jako czystą konstrukcję niezależną od stanu modułu. */
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

/** Zestaw akcji warstwy trzeciej gromadzi czynności uczestnika, które nie mają jeszcze miejsca w pasku akcji okna. */
function zlozZestawUczestnika(czynnosci: HTMLButtonElement[]): HTMLElement {
  return utworzZestawAkcji('Zestaw akcji uczestnika', czynnosci);
}

/** Funkcja przerysowuje warstwę diagnostyki kanałów, wymieniając jej wnętrze przy każdej zmianie składu debaty. */
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

/** Funkcja składa formularz tożsamości uczestnika wraz z przełącznikami widoku oraz ciałem ramy okna Model Panels. */
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

/** Funkcja podpina pasek akcji oraz przełączniki nastaw widoku do czynności okna o nazwie Model Panels. */
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
