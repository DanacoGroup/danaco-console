import {
  AppDeployEnvironment,
  AppDeployStatus,
  AppDeployStrategy,
  type AppDeployEnvironment as SrodowiskoWdrozenia,
  type AppDeployStrategy as StrategiaWdrozenia,
} from '../../../../shared/contract';
import { opisOdmowy } from '../../komponenty/odmowa';
import {
  poleTekstowe,
  poleWielowierszowe,
  przycisk,
  utworzWierszOdpowiedzi,
} from '../../modele/kontrolki-formularza';
import { utworzWykazBrakow } from './braki-kontraktu';
import { opiszPole } from './dymek-objasnienia';
import { BEZ_OKNA_MODULU, BRAKI_DEPLOYMENT, KODY_OKIEN, NAZWY_OKIEN } from './etykiety-apps';
import { utworzRameApps } from './rama-okna';
import { narzedziaDeploymentPanel } from './narzedzia-apps';
import { utworzPrzybornikApps } from './przybornik-apps';
import type { StanProduktu } from './stan-produktu';
import { utworzTabeleWdrozen } from './tabela-wdrozen';
import { utworzWyborZMenu, wierszWyboru } from './wybor-z-menu';
import {
  czyStanKoncowy,
  rozbieznoscZlecenia,
  warunekWarsztatow,
  zdaniePustkiWdrozen,
  type Zamowienie,
} from './zdania-wdrozen';

/**
 * Deployment Panel jest oknem zarządcą modułu Apps: uruchamia wdrożenie oraz
 * cofnięcie do wersji wcześniejszej, a przegląd statusu publikacji czyta
 * z własnej komendy odczytu, nie wyłącznie ze zdarzeń przejścia.
 */
export interface OknoDeploymentPanel {
  element: HTMLElement;
  odswiez(): void;
}

export function utworzOknoDeploymentPanel(stan: StanProduktu): OknoDeploymentPanel {
  const kod = KODY_OKIEN.DeploymentPanel;
  const rama = utworzRameApps(kod, NAZWY_OKIEN[kod] ?? kod, 'zarządca');

  // Rozwijanie z biblioteki kontrolek: rozwinięcie treści wraca do pola opisu, nie do nazwy pozycji.
  const srodowisko = utworzWyborZMenu('Środowisko wdrożenia', [
    { wartosc: AppDeployEnvironment.Dev, etykieta: 'dev', opis: 'środowisko deweloperskie' },
    { wartosc: AppDeployEnvironment.Staging, etykieta: 'staging', opis: 'przedprodukcyjne' },
    { wartosc: AppDeployEnvironment.Production, etykieta: 'produkcja' },
  ]);
  const strategia = utworzWyborZMenu('Strategia wdrożenia', [
    { wartosc: AppDeployStrategy.Immediate, etykieta: 'natychmiastowa' },
    { wartosc: AppDeployStrategy.Staged, etykieta: 'etapowa' },
    { wartosc: AppDeployStrategy.BlueGreen, etykieta: 'blue-green' },
  ]);
  const wersja = poleTekstowe({ etykieta: 'Wersja wdrażana', podpowiedz: 'np. 1.4.0' });
  const notatki = poleWielowierszowe({ etykieta: 'Notatki wydania' }, 4);

  const wdroz = przycisk('Wdróż', 'dn-btn dn-btn--sygnal dn-btn--sm');
  // Jedyna droga po historię trzymaną w rdzeniu; przycisk stoi tuż nad tabelą, której dotyczy.
  const odczytaj = przycisk('Odczytaj wdrożenia z rdzenia', 'dn-btn dn-btn--zarys dn-btn--sm');
  const odpowiedz = utworzWierszOdpowiedzi();
  const tabela = utworzTabeleWdrozen((identyfikator) => void wyslij(identyfikator));

  rama.akcje.append(utworzWykazBrakow('Bez drogi w kontrakcie', BRAKI_DEPLOYMENT));
  rama.akcje.append(
    utworzPrzybornikApps('Środowiska, domena, skalowanie i kondycja', narzedziaDeploymentPanel(stan))
      .element,
  );
  rama.tresc.append(
    opiszPole(
      wierszWyboru('Środowisko wdrożenia', srodowisko),
      'Trzy wartości wyliczenia kontraktu: dev, staging, produkcja.',
    ),
    opiszPole(
      wierszWyboru('Strategia wdrożenia', strategia),
      'Strategia z wyliczenia kontraktu: natychmiastowa, etapowa, blue-green.',
    ),
    opiszPole(wersja.element, 'Wersja bywa pusta — kontrakt jej nie wymaga; rdzeń nadaje własną.'),
    opiszPole(notatki.element, 'Notatki wydania idą do rdzenia wraz ze zleceniem wdrożenia.'),
    wdroz,
    odpowiedz.element,
    odczytaj,
    tabela.element,
  );

  wdroz.addEventListener('click', () => void wyslij(''));
  odczytaj.addEventListener('click', () => void odczytajWdrozenia());

  /**
   * Odczyt historii wdrożeń okna; zawężenia nie podstawiamy, o granicy rozstrzyga rdzeń.
   */
  async function odczytajWdrozenia(): Promise<void> {
    rama.ladowanie('Odczyt wdrożeń z rdzenia…');
    odpowiedz.pokaz('Odczyt wdrożeń: żądanie wysłane do rdzenia…', true);
    await stan.odczytajWdrozenia('', 0);
    const powod = stan.powodOdczytu('wdrozenia');
    if (powod !== '') {
      rama.blad(powod);
      odpowiedz.pokaz(powod, false);
      return;
    }
    // Liczba pochodzi ze zbioru po wchłonięciu, więc mówi o tym, co stoi w tabeli.
    odpowiedz.pokaz(`Odczyt wdrożeń: wykaz liczy ${stan.wdrozenia().length} pozycji.`, true);
    rama.gotowe();
    odswiez();
  }

  /** Przebieg zlecony z tego okna — o nim mówi wiersz odpowiedzi. */
  let zleconyPrzebieg = '';
  /** Czynność, którą zlecono ostatnio: „Wdrożenie" albo „Cofnięcie wdrożenia". */
  let zleconaCzynnosc = '';
  /** Stan przebiegu opowiedziany już Operatorowi — żeby nie powtarzać zdania. */
  let opowiedzianyStan = '';
  /** Rozbieżność zamówienia z odpowiedzią rdzenia; pusty łańcuch, gdy jej nie ma. */
  let rozbieznoscOdpowiedzi = '';

  async function wyslij(cofnijDo: string): Promise<void> {
    const idOkna = stan.idOkna();
    if (idOkna === '') {
      odpowiedz.pokaz(BEZ_OKNA_MODULU, false);
      return;
    }
    const czynnosc = cofnijDo === '' ? 'Wdrożenie' : 'Cofnięcie wdrożenia';
    const ostrzezenie = warunekWarsztatow(stan.architektura() !== null);
    const zamowienie: Zamowienie = {
      srodowisko: srodowisko.wartosc(),
      strategia: strategia.wartosc(),
      cofnijDo,
    };
    zleconyPrzebieg = '';
    zleconaCzynnosc = czynnosc;
    opowiedzianyStan = '';
    rozbieznoscOdpowiedzi = '';
    rama.ladowanie(`${czynnosc} w toku…`);
    odpowiedz.pokaz(`${czynnosc}: żądanie wysłane do rdzenia.${ostrzezenie}`, true);
    const wynik = await stan.zrodlo.uruchomWdrozenie({
      idOkna,
      srodowisko: zamowienie.srodowisko as SrodowiskoWdrozenia,
      strategia: zamowienie.strategia as StrategiaWdrozenia,
      wersja: wersja.kontrolka.value,
      notatki: notatki.kontrolka.value,
      cofnijDo,
    });
    if (!wynik.udany || wynik.wynik === undefined) {
      const zdanie = opisOdmowy(czynnosc, wynik.blad?.code, wynik.blad?.message);
      rama.blad(`${zdanie}${ostrzezenie}`);
      odpowiedz.pokaz(zdanie, false);
      return;
    }
    // Kolejność tych trzech kroków jest wiążąca: zdanie o przyjęciu musi paść przed wchłonięciem migawki.
    const przebieg = wynik.wynik.deployment;
    zleconyPrzebieg = przebieg.id;
    rozbieznoscOdpowiedzi = rozbieznoscZlecenia(zamowienie, przebieg);
    if (rozbieznoscOdpowiedzi === '') {
      // Wszystkie wartości tego zdania pochodzą z odpowiedzi rdzenia, nie z kontrolek okna.
      odpowiedz.pokaz(
        `${czynnosc}: rdzeń przyjął zlecenie i uruchomił przebieg ${przebieg.id} — ` +
          `środowisko ${przebieg.environment}, strategia ${przebieg.strategy}, ` +
          `wersja ${przebieg.version ?? 'nienadana'}, stan początkowy ${przebieg.status}. ` +
          'Stan końcowy przyjdzie zdarzeniem apps.build.changed.',
        true,
      );
      rama.gotowe();
    } else {
      const zdanie =
        `${czynnosc}: rdzeń założył przebieg ${przebieg.id}, ale ODDAŁ CO INNEGO, ` +
        `NIŻ ZAMÓWIONO — ${rozbieznoscOdpowiedzi}. Zlecenia nie potwierdzam.`;
      odpowiedz.pokaz(zdanie, false);
      rama.blad(zdanie);
    }
    stan.wchlonOdpowiedzWdrozenia(przebieg);
    // Zdarzenia przejścia potrafią wyprzedzić ten ciąg dalszy, więc stan końcowy opowiadamy tu ponownie.
    opowiedzStanPrzebiegu();
  }

  /** Dopisuje do wiersza odpowiedzi stan przebiegu zleconego z tego okna, z pól przysłanych przez rdzeń. */
  function opowiedzStanPrzebiegu(): void {
    if (zleconyPrzebieg === '') return;
    const przebieg = stan.wdrozenia().find((wdrozenie) => wdrozenie.id === zleconyPrzebieg);
    if (przebieg === undefined || przebieg.status === opowiedzianyStan) return;
    opowiedzianyStan = przebieg.status;
    // Rozbieżność raz zauważona jedzie z każdą dalszą wiadomością o tym przebiegu.
    const dopisek =
      rozbieznoscOdpowiedzi === ''
        ? ''
        : ` ZAMÓWIENIE I ODPOWIEDŹ SIĘ ROZESZŁY: ${rozbieznoscOdpowiedzi}.`;
    if (!czyStanKoncowy(przebieg.status)) {
      odpowiedz.pokaz(
        `${zleconaCzynnosc} ${zleconyPrzebieg}: stan ${przebieg.status}.${dopisek}`,
        rozbieznoscOdpowiedzi === '',
      );
      return;
    }
    // Powód dopisujemy, gdy rdzeń go przysłał; jego brak przy niepowodzeniu też jest wiadomością.
    const udany = przebieg.status === AppDeployStatus.Succeeded;
    const maPowod = przebieg.logRef !== undefined && przebieg.logRef !== '';
    const powod = maPowod
      ? ` Rdzeń podał: ${przebieg.logRef ?? ''}.`
      : udany
        ? ''
        : ' Rdzeń nie podał powodu (pole logRef puste).';
    // Środowisko i wersja jadą także tutaj z odpowiedzi rdzenia, nie z zamówienia — z powodu cofnięcia.
    odpowiedz.pokaz(
      `${zleconaCzynnosc} ${zleconyPrzebieg}: stan końcowy ${przebieg.status} — ` +
        `środowisko ${przebieg.environment}, wersja ${przebieg.version ?? 'nienadana'}.` +
        `${powod}${dopisek}`,
      udany && rozbieznoscOdpowiedzi === '',
    );
  }

  function odswiez(): void {
    const wdrozenia = stan.wdrozenia();
    tabela.nanies(wdrozenia);
    opowiedzStanPrzebiegu();
    if (rama.faza() === 'blad' || rama.faza() === 'ladowanie') return;
    if (wdrozenia.length === 0) {
      // Zdanie składamy przy każdym odświeżeniu, a nie raz przy budowie okna:
      rama.puste(zdaniePustkiWdrozen(stan.ramki(), stan.czyWdrozeniaCzytane()));
      return;
    }
    rama.gotowe();
  }

  odswiez();
  return { element: rama.element, odswiez };
}
