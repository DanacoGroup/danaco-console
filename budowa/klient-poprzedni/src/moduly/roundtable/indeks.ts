import './roundtable.css';
import './debata.css';
import './analiza-debaty.css';

import type { Kanal } from '../../protokol/kanal';
import { utworzRejestrKanalow, type RejestrKanalow } from '../../sterowanie/rejestr-kanalow';
import { widokZOknaSesji, type OpisModulu } from '../rejestracja';
import { utworzOknoArgumentMap } from './okno-argument-map';
import { utworzOknoConsensusPanel } from './okno-consensus-panel';
import { utworzOknoDebatePanel } from './okno-debate-panel';
import { utworzOknoModelPanels } from './okno-model-panels';
import { utworzOknoModeratorPanel } from './okno-moderator-panel';
import { utworzOknoVotingEvaluation } from './okno-voting-evaluation';
import { utworzPasekUczciwosci } from './pasek-uczciwosci';
import { utworzPasRozszerzen } from './rozszerzenia-boczne';
import { utworzStanDebaty, type StanDebaty } from './stan-debaty';
import { utworzStrumienWypowiedzi } from './strumien-wypowiedzi';
import { utworzZrodloArsenaluRoundtable } from './zrodlo-arsenalu';
import { utworzZrodloRoundtable } from './zrodlo-roundtable';
import { utworzZrodloStrumieniaDebaty } from './zrodlo-strumienia-debaty';

/**
 * Moduł Roundtable — złożenie sześciu okien wokół jednej debaty wielu modeli.
 *
 * Jedno okno rozmawia tu z wieloma kanałami naraz (`core/kompozycja.go`, pole
 * `Debata`), po tym samym rejestrze kanałów, którym jedzie okno rozmowy. Moduł
 * nie prowadzi więc własnej listy modeli: uczestnika zakłada się na kanale
 * wziętym z rejestru, a nazwę kanału składa `nazwaKanalu` z biblioteki.
 *
 * Moduł pracuje w oknie, nie w sesji. Każda komenda obszaru, którą moduł wywołuje, wymaga
 * `windowId`, więc złożenie jedzie przez `widokZOknaSesji` wraz z kodem modułu:
 * przejście pyta rdzeń o okna sesji i odracza montaż do chwili, gdy okno tego
 * modułu jest znane. Bez kodu przejście wzięłoby okno pierwsze w wykazie
 * i komendy debaty jechałyby z `windowId` okna cudzego modułu; szczegóły przy
 * `widokZOknaSesji` w `moduly/rejestracja.ts`.
 *
 * Rejestr kanałów zakładany jest tutaj. Wspólny egzemplarz powstaje
 * w `aplikacja/scena-sesji.ts` i schodzi do okna rozmowy przez
 * `wiazanie-gniazda.ts`, ale umowa rejestru modułów daje modułowi wyłącznie
 * `Kanal` (`OpisModulu.utworzWidok(kanal)`), więc egzemplarza rejestru nie da
 * się tą drogą podać. Drugi egzemplarz jest dopuszczalny, bo rejestr to czytająca
 * pamięć podręczna nad `channel.list` — katalog wyboru wspólny dla całego
 * klienta, w którym żadne okno nie zapisuje swojego stanu. Od egzemplarza
 * wspólnego różni się jedynie chwilą odświeżenia.
 *
 * Układ idzie za warstwami widoczności opracowania (rozdz. 3, 3.1). Warstwa
 * pierwsza stoi w pasie górnym: Model Panels (skład debaty, osobny panel na
 * każdy model) obok monitora przebiegu (Debate Panel) — pytanie idzie
 * jednocześnie do wszystkich uczestników i odpowiedzi narastają obok składu.
 * Warstwa druga to cztery rozszerzenia boczne — Argument Map & Analysis,
 * Voting & Evaluation Center, Moderator Panel i Consensus Panel — zwinięte do
 * chwili otwarcia przyciskiem z Debate Panelu. Zwinięcie nie jest blokadą:
 * okno jest zbudowane, zasubskrybowane i o jedno naciśnięcie dalej. Okno
 * rozmowy modułu nie należy do złożenia: jest bytem sesji i składa je warstwa
 * rozmowy.
 *
 * Strumień głosów jest jeden na złożenie, nie jeden na okno. Rdzeń nadaje
 * wypowiedź uczestnika fragment po fragmencie (`stream.chunk`,
 * `adapter_modul_roundtable_glos.go`) i rozgłasza `roundtable.debate.changed`
 * rodzaju `created` z treścią pustą, a `updated` z pełną dopiero po domknięciu
 * strumienia; okno słuchające samego zdarzenia pokazuje pustą wypowiedź przez
 * cały czas mówienia modelu. Gdyby każde z sześciu okien założyło własną
 * subskrypcję i własne gromadzenie, ten sam fragment byłby przyjęty
 * sześciokrotnie, a sześć okien miałoby sześć osobnych obrazów jednego głosu.
 * Subskrypcja stoi zatem tutaj, obok `StanDebaty`, i schodzi do okien gotowym
 * gromadzeniem.
 */
export interface ZamontowanyRoundtable {
  /** Element osadzony w dokumencie. */
  element: HTMLElement;
  /** Stan debaty wspólny sześciu oknom. */
  stan: StanDebaty;
  /** Odczytuje wszystkie okna z rdzenia. */
  odswiez(): void;
  /** Zamyka nasłuch zdarzeń modułu wraz z nasłuchem okien. */
  zamknij(): void;
}

export function zamontujRoundtable(
  gospodarz: HTMLElement,
  kanal: Kanal,
  okno: string,
  rejestrPodany?: RejestrKanalow,
): ZamontowanyRoundtable {
  // Rejestr przyjmowany z zewnątrz, gdy wywołujący go ma — inaczej zakładany
  // tutaj (patrz opis modułu). Parametr istnieje po to, żeby poszerzenie umowy
  // `OpisModulu` nie wymagało przepisywania złożenia.
  const rejestr = rejestrPodany ?? utworzRejestrKanalow(kanal);
  const zrodlo = utworzZrodloRoundtable(kanal);
  // Arsenał obszaru — czterdzieści dwie komendy poza czwórką prowadzącą debatę.
  // Jedno źródło na całe złożenie, bo wszystkie okna pytają o tę samą debatę
  // tego samego okna; drugi egzemplarz nie dałby niczego poza drugim miejscem
  // do zmiany przy zmianie kontraktu.
  const arsenal = utworzZrodloArsenaluRoundtable(kanal);
  const stan = utworzStanDebaty(zrodlo, rejestr, { okno });
  const strumien = utworzStrumienWypowiedzi(stan);
  // Okno czytane w chwili nadejścia fragmentu, nie w chwili subskrypcji:
  // `StanDebaty.ustawOkno` przestawia moduł na inną debatę bez składania okien
  // od nowa, a filtr zapamiętany przy subskrypcji przepuszczałby wtedy
  // fragmenty debaty poprzedniej.
  const zrodloStrumienia = utworzZrodloStrumieniaDebaty(kanal, () => stan.okno());

  // Zmiana tury czyści gromadzenie, i robi to przed oknami. Fragmenty należą do
  // wypowiedzi jednej tury, a `StanDebaty` porzuca wtedy wykaz wypowiedzi;
  // zostawione dokleiłyby zdania tury poprzedniej do mówców tury nowej. Nasłuch
  // stoi przed oknami, bo słuchacze stanu są wołani w kolejności zapisania:
  // gdyby stał po nich, okna zdążyłyby raz narysować głosy tury poprzedniej pod
  // nagłówkiem tury nowej.
  let turaWidziana = stan.tura();
  const odsubskrybujTure = stan.naZmiane(() => {
    if (stan.tura() === turaWidziana) return;
    turaWidziana = stan.tura();
    strumien.wyczysc();
  });

  const obszar = document.createElement('div');
  obszar.className = 'dr-modul';
  obszar.dataset['modul'] = 'roundtable';

  const sklad = utworzOknoModelPanels(zrodlo, arsenal, stan, strumien, okno);
  const moderator = utworzOknoModeratorPanel(zrodlo, arsenal, stan, okno);
  const stanowisko = utworzOknoConsensusPanel(zrodlo, arsenal, stan, okno);
  const analiza = utworzOknoArgumentMap(arsenal, stan, strumien);
  const ocena = utworzOknoVotingEvaluation(arsenal, stan, strumien);

  // Pas rozszerzeń powstaje przed Debate Panelem, bo to on wytwarza przyciski
  // otwarcia, które Debate Panel stawia w swoim pasku akcji.
  const rozszerzenia = utworzPasRozszerzen([
    {
      kod: 'argument-map-analysis',
      nazwa: 'Argument Map & Analysis',
      element: oznacz(analiza.element, 'argument-map-analysis'),
    },
    {
      kod: 'voting-evaluation-center',
      nazwa: 'Voting & Evaluation Center',
      element: oznacz(ocena.element, 'voting-evaluation-center'),
    },
    {
      kod: 'moderator-panel',
      nazwa: 'Moderator Panel',
      element: oznacz(moderator.element, 'moderator-panel'),
    },
    {
      kod: 'consensus-panel',
      nazwa: 'Consensus Panel',
      element: oznacz(stanowisko.element, 'consensus-panel'),
    },
  ]);

  const przebieg = utworzOknoDebatePanel(zrodlo, arsenal, stan, strumien, [
    rozszerzenia.przyciskOtwarcia('argument-map-analysis'),
    rozszerzenia.przyciskOtwarcia('voting-evaluation-center'),
    rozszerzenia.przyciskOtwarcia('moderator-panel'),
    rozszerzenia.przyciskOtwarcia('consensus-panel'),
  ]);

  const gorny = document.createElement('div');
  gorny.className = 'dr-modul__gora';
  gorny.append(
    oznacz(sklad.element, 'model-panels'),
    oznacz(przebieg.element, 'debate-panel'),
  );

  const uczciwosc = utworzPasekUczciwosci(kanal);

  obszar.append(uczciwosc.element, gorny, rozszerzenia.element);
  gospodarz.replaceChildren(obszar);

  function odswiez(): void {
    sklad.odswiez();
    przebieg.odswiez();
    moderator.odswiez();
    stanowisko.odswiez();
    analiza.odswiez();
    ocena.odswiez();
  }

  // Katalog kanałów zamawia złożenie, nie okno: okna modułu czytają ten sam
  // wykaz, więc cztery zapytania o to samo byłyby czynnością bez odbiorcy.
  // Zamówienie idzie przed odczytem okien — inaczej Model Panels, sprawdzając
  // `kanalyOdczytane()` jeszcze przed odpowiedzią, zamawia wykaz drugi raz
  // i na jeden montaż idą dwie koperty `channel.list`.
  rejestr.odswiez();
  odswiez();
  // Katalog okien rdzenia czyta się raz na montaż — pasek uczciwości liczy z niego
  // rozjazd między oknami przypisanymi modułowi a oknami, które moduł buduje.
  void uczciwosc.katalog.odczytaj();

  // Fragment budzi wyłącznie okna, które treść wypowiedzi rysują albo liczą:
  // skład z odpowiedziami, przebieg debaty, analizę struktury i zestawienie
  // udziału. Moderator Panel i Consensus Panel treści wypowiedzi nie pokazują,
  // więc przerysowanie ich kilkanaście razy na sekundę byłoby pracą bez
  // odbiorcy — a że oba trzymają wpisane nastawy tury i zakresu, kasowałoby też
  // wpisywaną treść.
  const odsubskrybujStrumien = zrodloStrumienia.naFragmentWypowiedzi((fragment) => {
    strumien.przyjmij(fragment);
    sklad.odswiezGlosy();
    przebieg.odswiezGlosy();
    analiza.odswiezGlosy();
    ocena.odswiezGlosy();
  });

  return {
    element: obszar,
    stan,
    odswiez,
    zamknij() {
      odsubskrybujStrumien();
      odsubskrybujTure();
      uczciwosc.katalog.zamknij();
      ocena.zamknij();
      analiza.zamknij();
      stanowisko.zamknij();
      moderator.zamknij();
      przebieg.zamknij();
      sklad.zamknij();
      stan.zamknij();
    },
  };
}

/** Znakuje okno kodem rejestru okien operacyjnych — po nim skacze nawigacja. */
function oznacz(element: HTMLElement, kod: string): HTMLElement {
  element.dataset['okno'] = kod;
  element.tabIndex = -1;
  return element;
}

/**
 * Kod modułu dla rejestru powłoki.
 *
 * Kod siedzi w module, a nie w mapie po stronie powłoki: dodanie modułu to jeden
 * wpis, a nie dwa, więc nie da się dodać modułu i zapomnieć o wytwórni.
 */
const KOD_MODULU = 'roundtable';

/**
 * Panel debaty dla stosu paneli pomocniczych — wyłącznie reeksport.
 *
 * Panel nie należy do złożenia okien powyżej i nie jest przez nie
 * stawiany: gospodarzem jest pas okien pomocniczych modułu albo kolumna paneli
 * sceny okien równoległych, a wpina go wytwórnia paneli
 * (`okna-pomocnicze/wytwornia-paneli.ts`).
 *
 * Reeksport jest dla czytającego, nie dla wytwórni: wytwórnia sięga wprost po
 * `moduly/roundtable/panel-debaty`, bo import tego pliku wciągnąłby do gospodarza
 * całe złożenie okien operacyjnych i rejestrację modułu, których
 * gospodarz nie stawia. Wpis tutaj mówi tyle, że panel należy do tego modułu.
 */
export { utworzPanelDebaty, KOD_PANELU_DEBATY } from './panel-debaty';

export const MODUL: OpisModulu = {
  kod: KOD_MODULU,
  utworzWidok: (kanal) =>
    widokZOknaSesji(
      kanal,
      (gospodarz, kanalOkna, okno) => zamontujRoundtable(gospodarz, kanalOkna, okno),
      KOD_MODULU,
    ),
};
