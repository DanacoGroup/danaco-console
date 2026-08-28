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

/** Zależności widoku środowiska: droga do rdzenia, opis okna sceny sesji oraz przejście na inną trasę powłoki. */
export interface ZaleznosciSrodowiska {
  /** Droga do rdzenia. */
  rdzen: PolaczenieZRdzeniem;
  /** Opis okna wspólny dla sceny sesji. */
  opis: OpisOkna;
  /** Przejście na inną trasę. */
  naTrase(trasa: Trasa): void;
}

/** Widok trasy „Środowisko" — powłoka środowiska, scena sesji i ich wzajemne przestawianie nawigacją boczną. */
export interface WidokSrodowiska extends WidokTrasy {
  /** Przestawia powłokę na wskazane środowisko i pozycję jego nawigacji. */
  ustawSrodowisko(klucz: KluczSrodowiska, modul?: string): void;
  /** Środowisko obecnie otwarte. */
  srodowisko(): KluczSrodowiska;
  /** Scena sesji — okna, rozmowy i sterowanie. */
  scena: ScenaSesji;
}

/** Widok trasy środowiska: powłoka środowiska ze sceną pracy w środku, złączone razem z routerem aplikacji. */
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

  // Moduł spoza wykazu środowiska — droga bezpośrednia, dobierana z katalogu `module.list`.
  async function otworzPozaWykazem(kodModulu: string): Promise<void> {
    const wynik = await zadajWykazModulow(rdzen.kanal, {});
    const moduly = wynik.wynik?.modules;
    if (!wynik.udany || moduly === undefined) return;
    const modul = moduly.find((wpis) => wpis.code === kodModulu);
    if (modul === undefined) return;
    powloka.nawigacja.wskazPozaWykazem(pozycjaModulu(modul));
  }

  powloka.nawigacja.naBrakPozycji((kodModulu) => void otworzPozaWykazem(kodModulu));

  // Wejście do przestrzeni odłożone rusza, gdy rdzeń odda sesję i okno rozmowy.
  rdzen.uzgodnienie.naOtwarcieOkna((okno) => {
    oknoRozmowy = okno.id;
    przestrzen.ponow();
  });

  const akcjeTras = utworzAkcjeTras(rdzen.transport, naTrase);
  akcje.prepend(akcjeTras.element);

  // Centrum powiadomień — kolumna boczna wyzwalana plakietką dzwonka w pasku górnym.

  // Kolumna siada w korzeniu powłoki, poza obszarem modułu — centrum obowiązuje w każdym środowisku.

  // Licznik plakietki nadąża zdarzeniami także wtedy, gdy kolumna stoi zamknięta.
  const centrum = utworzKolumnePowiadomien(rdzen.kanal, (nowe) =>
    powloka.pasek.ustawPowiadomienia(nowe),
  );
  powloka.element.append(centrum.element);
  powloka.podlaczZaczepy({ otworzPowiadomienia: () => centrum.przelacz() });

  // Pomocnik aplikacji, nie moduł obok innych — stąd korzeń powłoki.
  powloka.element.append(zaczepAsystenta(rdzen.kanal).element);

  // Always On Display jest funkcją globalną, nie modułem — warstwa siada w korzeniu powłoki.
  powloka.element.append(
    zaczepAod(rdzen.kanal, { klient: { id: rdzen.uzgodnienie.klient.id } }).element,
  );

  // Widoczność posunięć asystenta — jedyne miejsce widzące naraz rdzeń, nawigację modułów i scenę okien.
  const posuniecia = utworzZrodloPosuniec(rdzen.kanal, rdzen.uzgodnienie.klient.id);
  const pas = utworzPasPosuniec(posuniecia);
  powloka.element.append(pas.element);

  // Okno otwarte gdziekolwiek w tej sesji wchodzi na scenę od razu.
  posuniecia.naNoweOkno((okno) => scena.przyjmijOknoZRdzenia(okno));

  // Czy użytkownik jest w środku pisania — wymaga ogniska w polu tekstowym oraz niepustej treści.
  function operatorPisze(): boolean {
    const czynny = document.activeElement;
    if (!(czynny instanceof HTMLTextAreaElement) && !(czynny instanceof HTMLInputElement)) {
      return false;
    }
    return czynny.value.trim().length > 0;
  }

  // Podążanie ekranu za modułem asystenta — przestawia nawigację, tak jakby wybór padł ręcznie.
  function podazajZaModulem(kodModulu: string): void {
    if (kodModulu === '') return;
    // Porównanie po kluczu pozycji, nie po polu `modul` — pole `modul` jest identyfikatorem wiersza.
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
