import { zaczepAod } from '../aod/indeks';
import { zaczepAsystenta } from '../asystent-plywajacy/indeks';
import type { OpisOkna } from '../okno-komunikacji/opis-okna';
import { utworzKolumnePowiadomien } from '../powiadomienia/kolumna-powiadomien';
import { utworzPowloke, type KluczSrodowiska } from '../powloka/powloka';
import { pozycjaModulu } from '../powloka/srodowiska';
import { utworzZrodloNawigacji } from '../powloka/zrodlo-nawigacji';
import { zadajWykazModulow } from '../protokol/wykaz-modulow';
import { utworzAkcjeTras } from './akcje-tras';
import { pokazKomunikat } from './komunikaty';
import type { PolaczenieZRdzeniem } from './polaczenie-z-rdzeniem';
import { utworzPrzestrzenModulu } from './przestrzen-modulu';
import type { WidokTrasy } from './router';
import { utworzPasPosuniec } from './pas-posuniec';
import { utworzSceneSesji, type ScenaSesji } from './scena-sesji';
import { Trasa } from './trasy';
import { utworzZrodloPosuniec } from './zrodlo-posuniec';

/** Zależności widoku środowiska. */
export interface ZaleznosciSrodowiska {
  /** Droga do rdzenia. */
  rdzen: PolaczenieZRdzeniem;
  /** Opis okna wspólny dla sceny sesji. */
  opis: OpisOkna;
  /** Przejście na inną trasę. */
  naTrase(trasa: Trasa): void;
}

/** Widok trasy „Środowisko". */
export interface WidokSrodowiska extends WidokTrasy {
  /** Przestawia powłokę na wskazane środowisko i pozycję jego nawigacji. */
  ustawSrodowisko(klucz: KluczSrodowiska, modul?: string): void;
  /** Środowisko obecnie otwarte. */
  srodowisko(): KluczSrodowiska;
  /** Scena sesji — okna, rozmowy i sterowanie. */
  scena: ScenaSesji;
}

/**
 * Widok trasy środowiska: powłoka środowiska ze sceną pracy w środku.
 *
 * Jedna odpowiedzialność — związanie powłoki środowiska ze sceną sesji i z
 * routerem. Ten plik nie buduje ani pasa kart, ani nawigacji, ani okna;
 * wszystko przychodzi gotowe z `powloka/` i ze sceny.
 *
 * Treść obszaru roboczego rozstrzyga `przestrzen-modulu.ts` przy każdym wyborze
 * pozycji bocznej nawigacji. Scena raz osadzona nie jest niszczona — przełączenie
 * modułu wyłącznie ją odsłania albo chowa, bo w środku żyje okno rozmowy
 * z otwartym strumieniem do rdzenia.
 *
 * Powłoka powstaje bez połączenia, więc źródło wykazu nawigacji
 * (`environment.enter`, `module.list`) dokłada się tutaj — w jedynym miejscu,
 * które zna naraz powłokę i drogę do rdzenia.
 *
 * Karta sesji i okno komunikacji idą w parze: pas kart zgłasza założenie karty,
 * scena odpowiada oknem zamówionym w rdzeniu komendą `window.create`, a
 * zamknięcie karty zdejmuje okno ze sceny.
 *
 * Grupa akcji paska górnego dostaje wskaźnik łączności z rdzeniem oraz
 * przełącznik widoków. Przełącznik motywu jest już w pasku powłoki
 * (`powloka/akcje-paska`), więc drugiego się nie dokłada.
 *
 * Asystent pływający wisi w korzeniu powłoki, nie w obszarze roboczym: moduł
 * podmienia wyłącznie zawartość `powloka.obszar`, więc favikon w korzeniu
 * przetrwa każdą zmianę modułu. Asystent nie wchodzi na scenę sesji — nie jest
 * oknem równoległym i nie liczy się do sufitu liczby okien.
 */
export function utworzWidokSrodowiska(zaleznosci: ZaleznosciSrodowiska): WidokSrodowiska {
  const { rdzen, opis, naTrase } = zaleznosci;

  const powloka = utworzPowloke();
  const akcje = powloka.pasek.akcje.element;

  const scena = utworzSceneSesji({ rdzen, opis, akcjePaska: akcje });

  /** Okno rozmowy uzgodnione z rdzeniem — cel rekonfiguracji przy zmianie modułu. */
  let oknoRozmowy: string | undefined;

  const przestrzen = utworzPrzestrzenModulu({
    kanal: rdzen.kanal,
    obszar: powloka.obszar,
    scena: scena.element,
    idSesji: () => rdzen.kanal.sesja().id(),
    idOkna: () => oknoRozmowy,
  });

  powloka.naWyborModulu((pozycja, dane) => przestrzen.wybierz(pozycja, dane));
  powloka.nawigacja.podlaczZrodlo(
    utworzZrodloNawigacji(rdzen.platforma, rdzen.uzgodnienie.klient.id),
  );

  /**
   * Moduł spoza wykazu środowiska — droga bezpośrednia.
   *
   * Boczna nawigacja pokazuje wyłącznie moduły widoczne w macierzy
   * `srodowisko_modul`; widoczność rozstrzyga, czy pozycja stoi na liście, nie
   * czy wolno moduł otworzyć. Moduł bez wiersza macierzy dobiera się tutaj
   * z katalogu `module.list`, który zwraca komplet modułów platformy niezależnie
   * od macierzy.
   *
   * Moduł nieodnaleziony w katalogu nie daje odmowy ani komunikatu: kolumna
   * została już ustawiona na pierwszą pozycję wykazu, więc praca trwa dalej.
   */
  async function otworzPozaWykazem(kodModulu: string): Promise<void> {
    const wynik = await zadajWykazModulow(rdzen.kanal, {});
    const moduly = wynik.wynik?.modules;
    if (!wynik.udany || moduly === undefined) return;
    const modul = moduly.find((wpis) => wpis.code === kodModulu);
    if (modul === undefined) return;
    powloka.nawigacja.wskazPozaWykazem(pozycjaModulu(modul));
  }

  powloka.nawigacja.naBrakPozycji((kodModulu) => void otworzPozaWykazem(kodModulu));

  // Wejście do przestrzeni odłożone na czas uzgodnienia rusza, gdy rdzeń odda
  // sesję i okno rozmowy — wybór modułu sprzed uzgodnienia nie przepada.
  rdzen.uzgodnienie.naOtwarcieOkna((okno) => {
    oknoRozmowy = okno.id;
    przestrzen.ponow();
  });

  const akcjeTras = utworzAkcjeTras(rdzen.transport, naTrase);
  akcje.prepend(akcjeTras.element);

  // Centrum powiadomień — kolumna boczna obok obszaru roboczego, wyzwalana
  // plakietką dzwonka w pasku górnym (katalog komponentów, rozdz. 11.6).
  //
  // Kolumna siada w korzeniu powłoki, poza obszarem podmienianym przez moduł:
  // centrum obowiązuje w każdym środowisku i w każdym module w tej samej formie,
  // więc zejście razem z modułem gubiłoby rejestr przy każdym przełączeniu.
  //
  // Licznik plakietki nadąża zdarzeniami także wtedy, gdy kolumna stoi
  // zamknięta — po to jest plakietka.
  const centrum = utworzKolumnePowiadomien(rdzen.kanal, (nowe) =>
    powloka.pasek.ustawPowiadomienia(nowe),
  );
  powloka.element.append(centrum.element);
  powloka.podlaczZaczepy({ otworzPowiadomienia: () => centrum.przelacz() });

  // Pomocnik aplikacji, nie moduł obok innych — stąd korzeń powłoki.
  powloka.element.append(zaczepAsystenta(rdzen.kanal).element);

  // Always On Display jest funkcją globalną, a nie modułem: awatar ma być
  // widoczny bez interakcji i nie znikać przy przełączeniu środowiska ani
  // modułu (opracowanie funkcji globalnej, rozdz. 2.4). Warstwa siada więc
  // w korzeniu powłoki, poza obszarem podmienianym przez moduł — tak samo jak
  // asystent, i z tego samego powodu.
  powloka.element.append(
    zaczepAod(rdzen.kanal, { klient: { id: rdzen.uzgodnienie.klient.id } }).element,
  );

  // Widoczność posunięć asystenta. To jedyne miejsce, które widzi naraz drogę do
  // rdzenia, nawigację modułów i scenę okien — trzy rzeczy, którymi asystent
  // rusza — więc źródło posunięć i pas montują się tutaj.
  const posuniecia = utworzZrodloPosuniec(rdzen.kanal, rdzen.uzgodnienie.klient.id);
  const pas = utworzPasPosuniec(posuniecia);
  powloka.element.append(pas.element);

  // Okno otwarte gdziekolwiek w tej sesji wchodzi na scenę od razu.
  posuniecia.naNoweOkno((okno) => scena.przyjmijOknoZRdzenia(okno));

  /**
   * Czy użytkownik jest w środku pisania.
   *
   * Wymagane są oba warunki naraz: ognisko stoi w polu tekstowym i pole ma już
   * treść. Samo ognisko nie wystarcza — pole wypowiedzi dostaje je przy otwarciu
   * okna i trzyma bezterminowo, więc podążanie nie ruszyłoby nigdy.
   */
  function operatorPisze(): boolean {
    const czynny = document.activeElement;
    if (!(czynny instanceof HTMLTextAreaElement) && !(czynny instanceof HTMLInputElement)) {
      return false;
    }
    return czynny.value.trim().length > 0;
  }

  /**
   * Podążanie ekranu za modułem, na który przeszedł asystent.
   *
   * Przejście przestawia boczną nawigację i obszar roboczy, tak jakby wybór padł
   * ręcznie; inaczej okno i kolumna nawigacji pokazywałyby dwa różne moduły.
   *
   * Gdy pole wypowiedzi jest w trakcie pisania, wstrzymane zostaje wyłącznie
   * przesunięcie ekranu — posunięcie asystenta już się wykonało w rdzeniu.
   * Przejście nie przepada: ląduje na pasie jako jedno kliknięcie.
   */
  function podazajZaModulem(kodModulu: string): void {
    if (kodModulu === '') return;
    // Porównanie po kluczu pozycji, nie po polu `modul`: klucz jest kodem
    // modułu z rdzenia (`pozycjaModulu`), a `Window.moduleId` niesie ten sam
    // kod. Pole `modul` jest identyfikatorem wiersza i nigdy się z nim nie
    // zrówna — porównanie po nim dawałoby przejście przy każdym zdarzeniu.
    if (powloka.nawigacja.wybrana()?.klucz === kodModulu) return;
    const przejdz = (): void => powloka.nawigacja.wybierz(kodModulu);
    if (operatorPisze()) {
      pas.zaproponujPrzejscie(`Idź za asystentem do modułu ${kodModulu}`, przejdz);
      return;
    }
    przejdz();
  }

  posuniecia.naModulOkna((okno) => podazajZaModulem(okno.moduleId));

  // Karta sesji i okno komunikacji idą w parze.
  powloka.naNowaKarte(() => {
    if (scena.dodajOkno() !== null) return;
    pokazKomunikat({
      tytul: 'Scena jest pełna',
      tresc: 'Obok siebie mieszczą się najwyżej cztery okna komunikacji. Zamknij jedno, aby otworzyć kolejne.',
      waga: 'ostrz',
    });
  });

  powloka.naZamknieciekarty(() => {
    scena.zamknijOstatnie();
  });

  return {
    element: powloka.element,
    scena,

    przyWejsciu: () => akcjeTras.trasy.ustawBiezaca(Trasa.Srodowisko),

    ustawSrodowisko(klucz, modul) {
      powloka.ustawSrodowisko(klucz);
      if (modul !== undefined) powloka.nawigacja.wybierz(modul);
    },

    srodowisko: () => powloka.nawigacja.srodowisko().klucz,
  };
}
