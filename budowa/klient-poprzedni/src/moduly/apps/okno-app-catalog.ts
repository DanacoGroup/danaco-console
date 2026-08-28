import type { Extension, ExtensionKind } from '../../../../shared/contract';
import { opisOdmowy } from '../../komponenty/odmowa';
import { poleTekstowe, przycisk, utworzWierszOdpowiedzi } from '../../modele/kontrolki-formularza';
import { utworzWykazBrakow } from './braki-kontraktu';
import { opiszPole } from './dymek-objasnienia';
import {
  BEZ_ODCZYTU_KATALOGU,
  BRAKI_APP_CATALOG,
  KODY_OKIEN,
  NAZWY_OKIEN,
  OKNO_SPOZA_KATALOGU,
} from './etykiety-apps';
import {
  BEZ_ZAWEZENIA,
  POCHODZENIA_WIDOKU,
  RODZAJE_WIDOKU,
  STANY_WIDOKU,
  czyZawezone,
  zawez,
  zdanieZawezenia,
  type ZawezenieWidoku,
} from './filtr-katalogu';
import { nazwaPochodzenia, utworzKarteRozszerzenia } from './karta-rozszerzenia';
import { utworzRameApps } from './rama-okna';
import { narzedziaAppCatalog } from './narzedzia-rozszerzen';
import { utworzPrzybornikApps } from './przybornik-apps';
import type { StanRozszerzen } from './stan-rozszerzen';
import { utworzWyborZMenu, wierszWyboru } from './wybor-z-menu';

/**
 * App Catalog jest oknem wiodącym strony dystrybucji i konsumpcji modułu Apps, niezależnym od
 * okna modułu, bo komendy obszaru extension nie niosą identyfikatora okna.
 */
export interface OknoAppCatalog {
  element: HTMLElement;
  odswiez(): void;
}

export function utworzOknoAppCatalog(stan: StanRozszerzen): OknoAppCatalog {
  const kod = KODY_OKIEN.AppCatalog;
  const rama = utworzRameApps(kod, NAZWY_OKIEN[kod] ?? kod, 'wiodące');

  let zawezenie: ZawezenieWidoku = { ...BEZ_ZAWEZENIA };
  let siatka = true;

  const szukaj = poleTekstowe({
    etykieta: 'Szukaj w katalogu',
    podpowiedz: 'nazwa, kod albo opis pozycji',
  });
  const rodzaj = utworzWyborZMenu('Rodzaj pozycji', [...RODZAJE_WIDOKU]);
  const pochodzenie = utworzWyborZMenu('Źródło pozycji', [...POCHODZENIA_WIDOKU]);
  const stanPozycji = utworzWyborZMenu('Stan pozycji', [...STANY_WIDOKU]);

  const odczytaj = przycisk('Odczytaj katalog', 'dn-btn dn-btn--sm dn-btn--zarys');
  const widok = przycisk('Widok: siatka kart', 'dn-btn dn-btn--sm dn-btn--duch');
  widok.setAttribute('aria-pressed', 'true');

  const odpowiedz = utworzWierszOdpowiedzi();

  const zdanieFiltru = document.createElement('p');
  zdanieFiltru.className = 'mp-katalog__zdanie';

  const lista = document.createElement('ul');
  lista.className = 'mp-katalog mp-katalog--siatka';

  const granica = document.createElement('p');
  granica.className = 'dn-pole-opis mp-granica';
  granica.textContent = OKNO_SPOZA_KATALOGU;

  const pasek = document.createElement('div');
  pasek.className = 'mp-pasek';
  pasek.append(odczytaj, widok);

  rama.akcje.append(utworzWykazBrakow('Bez drogi w kontrakcie', BRAKI_APP_CATALOG, 'extension'));
  rama.akcje.append(
    utworzPrzybornikApps('Wyszukiwarka, szczegóły, kolekcje i rejestr', narzedziaAppCatalog(stan))
      .element,
  );
  rama.tresc.append(
    opiszPole(
      szukaj.element,
      'Szukanie obejmuje nazwę, kod i opis pozycji — czyli pola, które rdzeń oddaje ' +
        'w katalogu. Narzędzi udostępnianych przez pozycję katalog nie niesie.',
    ),
    opiszPole(
      wierszWyboru('Rodzaj pozycji', rodzaj),
      'Cztery rodzaje wyliczenia kontraktu: serwer MCP, wtyczka, integracja API, umiejętność.',
    ),
    opiszPole(
      wierszWyboru('Źródło pozycji', pochodzenie),
      'Dwa źródła wyliczenia kontraktu. Źródło rozstrzyga wyłącznie stan wyjściowy ' +
        'przy rejestracji; poza tym jest faktem, nie regułą.',
    ),
    wierszWyboru('Stan pozycji', stanPozycji),
    pasek,
    odpowiedz.element,
    zdanieFiltru,
    lista,
    granica,
  );

  szukaj.kontrolka.addEventListener('input', () => {
    zawezenie = { ...zawezenie, tekst: szukaj.kontrolka.value };
    odswiez();
  });
  rodzaj.naZmiane((klucz) => {
    zawezenie = { ...zawezenie, rodzaj: klucz };
    odswiez();
  });
  pochodzenie.naZmiane((klucz) => {
    zawezenie = { ...zawezenie, pochodzenie: klucz };
    odswiez();
  });
  stanPozycji.naZmiane((klucz) => {
    zawezenie = { ...zawezenie, stan: klucz };
    odswiez();
  });

  odczytaj.addEventListener('click', () => void odczytajKatalog());
  widok.addEventListener('click', () => {
    siatka = !siatka;
    widok.textContent = siatka ? 'Widok: siatka kart' : 'Widok: lista';
    widok.setAttribute('aria-pressed', String(siatka));
    lista.className = siatka ? 'mp-katalog mp-katalog--siatka' : 'mp-katalog mp-katalog--lista';
  });

  async function odczytajKatalog(): Promise<void> {
    rama.ladowanie('Odczyt katalogu rozszerzeń…');
    odpowiedz.pokaz('Odczyt katalogu: żądanie wysłane do rdzenia…', true);
    await stan.odczytajKatalog();
    const powod = stan.powodOdczytu('katalog');
    if (powod !== '') {
      rama.blad(powod);
      odpowiedz.pokaz(powod, false);
      return;
    }
    const ile = stan.rozszerzenia().length;
    odpowiedz.pokaz(
      ile === 0
        ? 'Odczyt katalogu: rdzeń nie zna ani jednej pozycji. To jest pustka POTWIERDZONA, ' +
          'nie brak odczytu.'
        : `Odczyt katalogu: rdzeń oddał ${ile} pozycji.`,
      true,
    );
    rama.gotowe();
    odswiez();
  }

  /** Instalacja pozycji wysyła pochodzenie wprost, bo pominięte pole w kontrakcie znaczy personal. */
  async function zainstaluj(pozycja: Extension): Promise<void> {
    rama.ladowanie(`Instalacja pozycji ${pozycja.name} w toku…`);
    odpowiedz.pokaz(`Instalacja ${pozycja.name}: żądanie wysłane do rdzenia…`, true);
    const wynik = await stan.zrodlo.zainstaluj({
      kod: pozycja.code,
      rodzaj: pozycja.kind as ExtensionKind,
      zrodlo: '',
      idPunktuDostepu: pozycja.accessPointId ?? '',
      pochodzenie: pozycja.origin,
    });
    if (!wynik.udany || wynik.wynik === undefined) {
      const zdanie = opisOdmowy('Instalacja pozycji', wynik.blad?.code, wynik.blad?.message);
      rama.blad(zdanie);
      odpowiedz.pokaz(zdanie, false);
      return;
    }
    const oddana = wynik.wynik.extension;
    stan.wchlonPozycje(oddana);
    rama.gotowe();
    odswiez();
    // Zdanie o stanie wyjściowym pada od razu, zanim Operator zauważy skutek instalacji.
    const oStanie = oddana.enabled
      ? 'Pozycja jest włączona.'
      : `Pozycja pozostaje WYŁĄCZONA — taki jest stan wyjściowy źródła ` +
        `${nazwaPochodzenia(oddana.origin)}. Włącz ją świadomie w Installed Apps Managerze ` +
        'albo przyciskiem karty.';
    odpowiedz.pokaz(
      `Rdzeń zarejestrował pozycję ${oddana.name} (kod ${oddana.code}, wersja ` +
        `${oddana.version ?? 'nienadana'}). ${oStanie}`,
      true,
    );
  }

  async function przelacz(pozycja: Extension): Promise<void> {
    const czynnosc = pozycja.enabled ? 'Wyłączenie' : 'Włączenie';
    rama.ladowanie(`${czynnosc} pozycji ${pozycja.name} w toku…`);
    const wynik = await stan.zrodlo.przelacz(pozycja.id, !pozycja.enabled);
    if (!wynik.udany || wynik.wynik === undefined) {
      const zdanie = opisOdmowy(czynnosc, wynik.blad?.code, wynik.blad?.message);
      rama.blad(zdanie);
      odpowiedz.pokaz(zdanie, false);
      return;
    }
    const oddana = wynik.wynik.extension;
    stan.wchlonPozycje(oddana);
    rama.gotowe();
    odswiez();
    // Stan pochodzi z odpowiedzi rdzenia, nie z zamówienia — rdzeń rozstrzyga przełączenie.
    odpowiedz.pokaz(
      `${czynnosc} pozycji ${oddana.name}: rdzeń oddał stan ` +
        `${oddana.enabled ? 'włączona' : 'wyłączona'}.`,
      oddana.enabled !== pozycja.enabled,
    );
  }

  function odswiez(): void {
    const wszystkie = stan.rozszerzenia();
    const pokazane = zawez(wszystkie, zawezenie);
    lista.replaceChildren(
      ...pokazane.map((pozycja) =>
        utworzKarteRozszerzenia(pozycja, {
          zainstaluj: (cel) => void zainstaluj(cel),
          przelacz: (cel) => void przelacz(cel),
          wskaz: (cel) => stan.wybierz(cel.id),
        }),
      ),
    );
    zdanieFiltru.textContent = zdanieZawezenia(pokazane.length, wszystkie.length, zawezenie);
    zdanieFiltru.hidden = zdanieFiltru.textContent === '';
    if (rama.faza() === 'blad' || rama.faza() === 'ladowanie') return;
    if (pokazane.length > 0) {
      rama.gotowe();
      return;
    }
    // Trzy pustki znaczą co innego: nikt nie pytał, rdzeń nie ma nic, filtr nic nie przepuścił.
    if (!stan.czyKatalogCzytany()) {
      rama.puste(BEZ_ODCZYTU_KATALOGU);
      return;
    }
    if (wszystkie.length === 0) {
      rama.puste(
        'Rdzeń odpowiedział na odczyt i nie zna ani jednej pozycji katalogu. To jest ' +
          'pustka POTWIERDZONA — rejestr rozszerzeń platformy jest pusty.',
      );
      return;
    }
    rama.puste(
      `Żadna z ${wszystkie.length} pozycji katalogu nie spełnia zawężenia. ` +
        (czyZawezone(zawezenie) ? 'Zdejmij filtr, żeby zobaczyć komplet.' : ''),
    );
  }

  odswiez();
  return { element: rama.element, odswiez };
}
