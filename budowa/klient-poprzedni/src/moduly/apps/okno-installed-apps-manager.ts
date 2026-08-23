import { ExtensionKind, ExtensionOrigin, type Extension } from '../../../../shared/contract';
import { opisOdmowy } from '../../komponenty/odmowa';
import {
  pobierzPlik,
  poleTekstowe,
  poleWielowierszowe,
  przycisk,
  utworzWierszOdpowiedzi,
} from '../../modele/kontrolki-formularza';
import { utworzWykazBrakow } from './braki-kontraktu';
import { opiszPole } from './dymek-objasnienia';
import {
  BEZ_ODCZYTU_KATALOGU,
  BRAKI_INSTALLED_APPS,
  KODY_OKIEN,
  NAZWY_OKIEN,
  OKNO_SPOZA_KATALOGU,
} from './etykiety-apps';
import { RODZAJE_WIDOKU } from './filtr-katalogu';
import { nazwaPochodzenia, opisStanu } from './karta-rozszerzenia';
import { utworzRameApps } from './rama-okna';
import { narzedziaInstalledApps } from './narzedzia-rozszerzen';
import { utworzPrzybornikApps } from './przybornik-apps';
import type { StanRozszerzen } from './stan-rozszerzen';
import { utworzWyborZMenu, wierszWyboru } from './wybor-z-menu';

/**
 * Installed Apps Manager — okno zarządca rejestru rozszerzeń.
 *
 * Okno pokazuje wyłącznie pozycje zainstalowane, bo tym różni się od App
 * Catalogu: katalog jest widokiem oferty, rejestr — widokiem stanu maszyny.
 * Odsiew idzie po polu `installed` nad wspólnym zbiorem, a nie osobnym żądaniem
 * z `installedOnly`: dwa odczyty dałyby dwie migawki z dwóch chwil, a oba okna
 * mają pokazywać ten sam rejestr.
 *
 * Wyłączenie nie usuwa pozycji z rejestru, a odinstalowanie usuwa — to dwie
 * różne czynności i dwa różne przyciski. Znaczenie odinstalowania nie zostało
 * przez Właściciela rozstrzygnięte (kontrakt mówi to wprost przy komendzie),
 * więc okno nie dopowiada, co pozycja traci: pokazuje odpowiedź rdzenia.
 *
 * Instalacja przyjmuje kod, rodzaj, deklarowane źródło i pochodzenie. Wskazania
 * pliku z urządzenia tu nie ma i nie udajemy go polem wyboru pliku: kontrakt nie
 * prowadzi przesłania treści kanałem, więc pole otwierające okno wyboru pliku
 * kończyłoby się treścią, której nie ma czym wysłać.
 *
 * Eksport konfiguracji nie potrzebuje komendy: treść jest już w oknie, więc
 * zapis do pliku robi przeglądarka. Import idzie `extension.configure` — jedną
 * pozycją naraz, bo tyle niesie kontrakt.
 */
export interface OknoInstalledAppsManager {
  element: HTMLElement;
  odswiez(): void;
}

export function utworzOknoInstalledAppsManager(stan: StanRozszerzen): OknoInstalledAppsManager {
  const kod = KODY_OKIEN.InstalledAppsManager;
  const rama = utworzRameApps(kod, NAZWY_OKIEN[kod] ?? kod, 'zarządca');

  const wykaz = document.createElement('table');
  wykaz.className = 'dn-tabela mp-rejestr';
  const glowa = document.createElement('thead');
  const wierszGlowy = document.createElement('tr');
  for (const tytul of ['Nazwa', 'Rodzaj', 'Źródło', 'Wersja', 'Stan', 'Czynności']) {
    const komorka = document.createElement('th');
    komorka.scope = 'col';
    komorka.textContent = tytul;
    wierszGlowy.append(komorka);
  }
  glowa.append(wierszGlowy);
  const cialo = document.createElement('tbody');
  wykaz.append(glowa, cialo);

  const kodPozycji = poleTekstowe({
    etykieta: 'Kod instalowanej pozycji',
    podpowiedz: 'kod stały między wydaniami',
  });
  const rodzaj = utworzWyborZMenu(
    'Rodzaj rozszerzenia',
    RODZAJE_WIDOKU.filter((pozycja) => pozycja.wartosc !== ''),
  );
  const zrodloPaczki = poleTekstowe({
    etykieta: 'Deklarowane źródło',
    podpowiedz: 'ścieżka paczki albo adres serwera',
  });
  const pochodzenie = utworzWyborZMenu('Pochodzenie pozycji', [
    { wartosc: ExtensionOrigin.Personal, etykieta: 'Personal', opis: 'stan wyjściowy: wyłączone' },
    { wartosc: ExtensionOrigin.Danaco, etykieta: 'Danaco Plugin', opis: 'stan wyjściowy: włączone' },
  ]);
  const zainstaluj = przycisk('Zainstaluj pozycję', 'dn-btn dn-btn--sm dn-btn--atrament');

  const konfiguracja = poleWielowierszowe(
    { etykieta: 'Konfiguracja wskazanej pozycji (JSON)', podpowiedz: '{"transport":"stdio"}' },
    4,
  );
  const zapiszKonfiguracje = przycisk('Zapisz konfigurację', 'dn-btn dn-btn--sm dn-btn--atrament');
  const wyeksportuj = przycisk('Wyeksportuj konfigurację rejestru', 'dn-btn dn-btn--sm dn-btn--zarys');

  const odczytaj = przycisk('Odczytaj rejestr', 'dn-btn dn-btn--sm dn-btn--zarys');
  const odpowiedz = utworzWierszOdpowiedzi();

  const wskazana = document.createElement('p');
  wskazana.className = 'mp-rejestr__wskazana';

  const granica = document.createElement('p');
  granica.className = 'dn-pole-opis mp-granica';
  granica.textContent = OKNO_SPOZA_KATALOGU;

  const pasekInstalacji = document.createElement('div');
  pasekInstalacji.className = 'mp-pasek';
  pasekInstalacji.append(zainstaluj);

  const pasekKonfiguracji = document.createElement('div');
  pasekKonfiguracji.className = 'mp-pasek';
  pasekKonfiguracji.append(zapiszKonfiguracje, wyeksportuj);

  rama.akcje.append(
    utworzPrzybornikApps('Wersje, paczki, zestawy i dziennik', narzedziaInstalledApps(stan)).element,
    utworzWykazBrakow('Bez drogi w kontrakcie', BRAKI_INSTALLED_APPS, 'extension'),
  );
  rama.tresc.append(
    odczytaj,
    wykaz,
    naglowekCzesci('Instalacja pozycji'),
    kodPozycji.element,
    wierszWyboru('Rodzaj rozszerzenia', rodzaj),
    opiszPole(
      zrodloPaczki.element,
      'Napis podawany rdzeniowi jako źródło; rdzeń zapisuje go wprost i sam rozstrzyga, ' +
        'co z nim zrobi. Treści pliku kontrakt tą drogą nie przenosi.',
    ),
    opiszPole(
      wierszWyboru('Pochodzenie pozycji', pochodzenie),
      'Pochodzenie rozstrzyga stan wyjściowy rejestracji i wyłącznie to: pozycja Danaco ' +
        'Plugin staje włączona, pozycja Personal — wyłączona.',
    ),
    pasekInstalacji,
    naglowekCzesci('Konfiguracja pozycji wskazanej'),
    wskazana,
    opiszPole(
      konfiguracja.element,
      'Konfiguracja jest dla kontraktu nieprzezroczystym JSON-em: rdzeń sprawdza wyłącznie ' +
        'poprawność zapisu, a znaczenie pól należy do samego rozszerzenia.',
    ),
    pasekKonfiguracji,
    odpowiedz.element,
    granica,
  );

  odczytaj.addEventListener('click', () => void odczytajRejestr());
  zainstaluj.addEventListener('click', () => void wyslijInstalacje());
  zapiszKonfiguracje.addEventListener('click', () => void wyslijKonfiguracje());
  wyeksportuj.addEventListener('click', () => eksportujRejestr());

  /** Pozycje zainstalowane — to, czym rejestr różni się od katalogu. */
  function zainstalowane(): readonly Extension[] {
    return stan.rozszerzenia().filter((pozycja) => pozycja.installed);
  }

  async function odczytajRejestr(): Promise<void> {
    rama.ladowanie('Odczyt rejestru rozszerzeń…');
    odpowiedz.pokaz('Odczyt rejestru: żądanie wysłane do rdzenia…', true);
    await stan.odczytajKatalog();
    const powod = stan.powodOdczytu('katalog');
    if (powod !== '') {
      rama.blad(powod);
      odpowiedz.pokaz(powod, false);
      return;
    }
    rama.gotowe();
    odswiez();
    odpowiedz.pokaz(
      `Odczyt rejestru: zainstalowanych pozycji ${zainstalowane().length} ` +
        `z ${stan.rozszerzenia().length} w katalogu.`,
      true,
    );
  }

  async function wyslijInstalacje(): Promise<void> {
    if (kodPozycji.kontrolka.value === '') {
      odpowiedz.pokaz(
        'Okno nie wysyła pustego kodu pozycji — wpisz go. O tym, czy kod wskazuje coś ' +
          'istniejącego, rozstrzyga rdzeń, nie to okno.',
        false,
      );
      return;
    }
    rama.ladowanie('Instalacja pozycji w toku…');
    const zamowionePochodzenie = pochodzenie.wartosc();
    const wynik = await stan.zrodlo.zainstaluj({
      kod: kodPozycji.kontrolka.value,
      rodzaj: rodzaj.wartosc() as ExtensionKind,
      zrodlo: zrodloPaczki.kontrolka.value,
      idPunktuDostepu: '',
      pochodzenie: zamowionePochodzenie as ExtensionOrigin,
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
    // Rozbieżność pochodzenia jest odmową, nie drobiazgiem: to jedyne pole
    // rozstrzygające stan wyjściowy, więc oddanie innego niż zamówione znaczy
    // inny stan włączenia, niż Operator zamawiał.
    if (oddana.origin !== zamowionePochodzenie) {
      const zdanie =
        `Instalacja pozycji ${oddana.name}: rdzeń ODDAŁ CO INNEGO, NIŻ ZAMÓWIONO — ` +
        `zamówiono pochodzenie ${nazwaPochodzenia(zamowionePochodzenie)}, rdzeń oddał ` +
        `${nazwaPochodzenia(oddana.origin)}, a wraz z nim stan ${opisStanu(oddana)}.`;
      rama.blad(zdanie);
      odpowiedz.pokaz(zdanie, false);
      return;
    }
    odpowiedz.pokaz(
      `Rdzeń zarejestrował pozycję ${oddana.name} (kod ${oddana.code}) — ` +
        `${opisStanu(oddana)}, źródło ${nazwaPochodzenia(oddana.origin)}.`,
      true,
    );
  }

  async function wyslijKonfiguracje(): Promise<void> {
    const pozycja = stan.wybrane();
    if (pozycja === null) {
      odpowiedz.pokaz(
        'Wskaż pozycję w wykazie — konfiguracja należy do jednej pozycji rejestru.',
        false,
      );
      return;
    }
    let tresc: unknown;
    try {
      tresc = JSON.parse(konfiguracja.kontrolka.value === '' ? '{}' : konfiguracja.kontrolka.value);
    } catch (blad) {
      // Zapis niepoprawny zatrzymujemy w oknie i mówimy, gdzie leży usterka:
      // rdzeń odrzuciłby go tak samo, ale bez wskazania miejsca w tekście.
      odpowiedz.pokaz(
        `Konfiguracja nie jest poprawnym zapisem JSON: ${(blad as Error).message}. ` +
          'Żądanie nie zostało wysłane.',
        false,
      );
      return;
    }
    rama.ladowanie('Zapis konfiguracji w toku…');
    const wynik = await stan.zrodlo.skonfiguruj(pozycja.id, tresc, pozycja.accessPointId ?? '');
    if (!wynik.udany || wynik.wynik === undefined) {
      const zdanie = opisOdmowy('Zapis konfiguracji', wynik.blad?.code, wynik.blad?.message);
      rama.blad(zdanie);
      odpowiedz.pokaz(zdanie, false);
      return;
    }
    stan.wchlonPozycje(wynik.wynik.extension);
    rama.gotowe();
    odswiez();
    odpowiedz.pokaz(
      `Rdzeń potwierdził zapis konfiguracji pozycji ${wynik.wynik.extension.name}.`,
      true,
    );
  }

  /**
   * Eksport konfiguracji rejestru do pliku.
   *
   * Bez komendy i bez potrzeby: treść jest już w oknie. Eksport niesie kod,
   * rodzaj, pochodzenie, stan i konfigurację — czyli to, czym pozycję da się
   * odtworzyć komendą instalacji i konfiguracji na innej instancji. Sekretów
   * w nim nie ma, bo pozycja katalogu ich nie niesie.
   */
  function eksportujRejestr(): void {
    const pozycje = zainstalowane();
    if (pozycje.length === 0) {
      odpowiedz.pokaz('Rejestr nie ma ani jednej pozycji zainstalowanej — nie ma czego wyeksportować.', false);
      return;
    }
    const tresc = JSON.stringify(
      pozycje.map((pozycja) => ({
        code: pozycja.code,
        kind: pozycja.kind,
        origin: pozycja.origin,
        version: pozycja.version ?? null,
        enabled: pozycja.enabled,
        config: pozycja.config ?? null,
      })),
      null,
      2,
    );
    pobierzPlik('rejestr-rozszerzen.json', tresc);
    odpowiedz.pokaz(
      `Zapisano ${pozycje.length} pozycji rejestru do pliku. Eksport nie niesie danych ` +
        'dostępowych, bo pozycja katalogu ich nie niesie.',
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
    stan.wchlonPozycje(wynik.wynik.extension);
    rama.gotowe();
    odswiez();
    odpowiedz.pokaz(
      `${czynnosc} pozycji ${wynik.wynik.extension.name}: rdzeń oddał stan ` +
        `${opisStanu(wynik.wynik.extension)}. Zmiana jest rozgłaszana na wszystkie ` +
        'urządzenia sesji.',
      true,
    );
  }

  async function odinstaluj(pozycja: Extension): Promise<void> {
    rama.ladowanie(`Odinstalowanie pozycji ${pozycja.name} w toku…`);
    const wynik = await stan.zrodlo.odinstaluj(pozycja.id);
    if (!wynik.udany || wynik.wynik === undefined) {
      const zdanie = opisOdmowy('Odinstalowanie pozycji', wynik.blad?.code, wynik.blad?.message);
      rama.blad(zdanie);
      odpowiedz.pokaz(zdanie, false);
      return;
    }
    rama.gotowe();
    // Odpowiedź `false` jest odpowiedzią udaną i znaczy „rdzeń pozycji nie
    // zdjął". Zbioru wtedy nie ruszamy, bo pozycja nadal w nim stoi.
    if (!wynik.wynik.uninstalled) {
      odpowiedz.pokaz(
        `Odinstalowanie pozycji ${pozycja.name}: rdzeń odpowiedział, że pozycji NIE zdjął. ` +
          'Wykaz zostaje bez zmiany.',
        false,
      );
      return;
    }
    stan.zdejmijPozycje(pozycja.id);
    odswiez();
    odpowiedz.pokaz(`Rdzeń zdjął pozycję ${pozycja.name} z rejestru.`, true);
  }

  function odswiez(): void {
    const pozycje = zainstalowane();
    cialo.replaceChildren(...pozycje.map(wierszRejestru));
    const wybrana = stan.wybrane();
    wskazana.textContent =
      wybrana === null
        ? 'Nie wskazano pozycji. Naciśnij „Wskaż" w wierszu wykazu, żeby wczytać jej konfigurację.'
        : `Wskazana pozycja: ${wybrana.name} (kod ${wybrana.code}).`;
    if (rama.faza() === 'blad' || rama.faza() === 'ladowanie') return;
    if (pozycje.length > 0) {
      rama.gotowe();
      return;
    }
    if (!stan.czyKatalogCzytany()) {
      rama.puste(BEZ_ODCZYTU_KATALOGU);
      return;
    }
    rama.puste(
      'Rdzeń odpowiedział na odczyt i nie zna ani jednej pozycji ZAINSTALOWANEJ. ' +
        (stan.rozszerzenia().length === 0
          ? 'Katalog jest pusty w całości.'
          : `Katalog liczy ${stan.rozszerzenia().length} pozycji dostępnych — zainstaluj ` +
            'którąś w App Catalogu albo poniżej.'),
    );
  }

  /** Jeden wiersz rejestru wraz z czynnościami wykonywanymi na pozycji. */
  function wierszRejestru(pozycja: Extension): HTMLElement {
    const element = document.createElement('tr');
    element.dataset['rozszerzenie'] = pozycja.id;
    element.append(
      komorka(pozycja.name),
      komorka(pozycja.kind),
      komorka(nazwaPochodzenia(pozycja.origin)),
      komorka(pozycja.version ?? 'nieoddana'),
      komorkaStanu(pozycja),
      komorkaCzynnosci(pozycja),
    );
    return element;
  }

  function komorkaStanu(pozycja: Extension): HTMLElement {
    const plakietka = document.createElement('span');
    plakietka.className = pozycja.enabled
      ? 'dn-plakietka dn-plakietka--sukces'
      : 'dn-plakietka dn-plakietka--ostrzezenie';
    plakietka.textContent = pozycja.enabled ? 'włączone' : 'wyłączone';
    const element = document.createElement('td');
    element.append(plakietka);
    return element;
  }

  function komorkaCzynnosci(pozycja: Extension): HTMLElement {
    const element = document.createElement('td');
    element.className = 'mp-rejestr__czynnosci';
    element.append(
      przyciskWiersza(pozycja.enabled ? 'Wyłącz' : 'Włącz', 'dn-btn dn-btn--zarys dn-btn--sm', () =>
        void przelacz(pozycja),
      ),
      przyciskWiersza('Wskaż', 'dn-btn dn-btn--duch dn-btn--sm', () => {
        stan.wybierz(pozycja.id);
        konfiguracja.kontrolka.value =
          pozycja.config === undefined || pozycja.config === null
            ? ''
            : JSON.stringify(pozycja.config, null, 2);
        odswiez();
      }),
      przyciskWiersza('Odinstaluj', 'dn-btn dn-btn--niebezpieczny dn-btn--sm', () =>
        void odinstaluj(pozycja),
      ),
    );
    return element;
  }

  odswiez();
  return { element: rama.element, odswiez };
}

function komorka(tresc: string): HTMLElement {
  const element = document.createElement('td');
  element.textContent = tresc;
  return element;
}

function przyciskWiersza(etykieta: string, klasa: string, naNacisniecie: () => void): HTMLElement {
  const element = document.createElement('button');
  element.type = 'button';
  element.className = klasa;
  element.textContent = etykieta;
  element.addEventListener('click', naNacisniecie);
  return element;
}

/** Nagłówek części okna. */
function naglowekCzesci(tresc: string): HTMLElement {
  const element = document.createElement('p');
  element.className = 'mp-czesc__tytul';
  element.textContent = tresc;
  return element;
}
