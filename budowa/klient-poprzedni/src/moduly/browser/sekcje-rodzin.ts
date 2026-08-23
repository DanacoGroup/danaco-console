import {
  BrowserArtifactKind,
  Command,
  BrowserDevicePreset,
  BrowserDownloadAction,
  BrowserMacroAction,
  BrowserNoteClassification,
  BrowserScreenshotMode,
} from '../../../../shared/contract';
import {
  liczbaPola,
  odrzuc,
  oknoAlboOdmowa,
  opcjonalne,
  wymagane,
  type SekcjaRodzin,
} from './panel-rodzin';
import type { StanPrzegladania } from './stan-przegladania';

/**
 * Czynności rodzin `browser.*` rozpisane na sekcje panelu — po jednej grupie na
 * okno operacyjne, zgodnie z przypisaniem z opracowania modułu:
 *
 *   Browser Window          → karty, grupy kart, przestrzenie, zakładki,
 *                             przewinięcie, zrzut i narzędzia inspekcyjne
 *   Capture & Monitor Panel → monitory, kanały, kolejka czytania, pobrania,
 *                             wytwory sesji
 *   Notes / Sources Panel   → zestawy źródeł, wątki notatek, zmiana notatki
 *   Automation Studio       → nagrywarka makr i granice Wykonawcy
 *
 * Plik nie buduje elementów i nie trzyma stanu: składa wywołania rdzenia
 * i zdania o ich skutku. Elementy stawia `panel-rodzin.ts`, stan modułu stoi
 * w `stan-przegladania.ts`.
 *
 * Każde zdanie o skutku mówi liczbę albo identyfikator wzięty z odpowiedzi
 * rdzenia — nigdy samego „gotowe". Zdanie bez pokrycia w odpowiedzi byłoby
 * meldunkiem bez skutku.
 */

/** Sekcje wpinane do Browser Window. */
export function sekcjePrzegladania(stan: StanPrzegladania): SekcjaRodzin[] {
  const zrodlo = stan.zrodlo;
  return [
    {
      tytul: 'Karty i przestrzenie robocze',
      opis:
        'Rząd kart okna przeglądarki prowadzi rdzeń: karta otwarta z adresem naprawdę pod ten ' +
        'adres przechodzi, a przestrzeń robocza zapamiętuje skład kart, nie ich kopie.',
      czynnosci: [
        {
          nazwa: 'Otwórz kartę',
          pola: [
            { klucz: 'adres', etykieta: 'Adres karty' },
            { klucz: 'karta', etykieta: 'Identyfikator karty' },
            { klucz: 'nazwa', etykieta: 'Nazwa grupy albo przestrzeni' },
          ],
          async wykonaj(wartosci) {
            const wynik = await zrodlo.otworzKarte({
              windowId: oknoAlboOdmowa(stan),
              url: opcjonalne(wartosci, 'adres'),
            });
            if (!wynik.udany || wynik.wynik === undefined) odrzuc('Otwarcie karty', wynik.blad);
            const karta = wynik.wynik.tab;
            return `Karta ${karta.id} otwarta (stan: ${karta.state}, pozycja ${karta.order}).`;
          },
        },
        {
          nazwa: 'Wykaz kart',
          async wykonaj() {
            const wynik = await zrodlo.wykazKart({ windowId: oknoAlboOdmowa(stan) });
            if (!wynik.udany || wynik.wynik === undefined) odrzuc('Wykaz kart', wynik.blad);
            const grupy = wynik.wynik.groups ?? [];
            return `Kart otwartych: ${wynik.wynik.tabs.length}; grup kart: ${grupy.length}.`;
          },
        },
        {
          nazwa: 'Uczyń kartę czynną',
          async wykonaj(wartosci) {
            const wynik = await zrodlo.zmienKarte({
              tabId: wymagane(wartosci, 'karta', 'identyfikator karty'),
              active: true,
            });
            if (!wynik.udany || wynik.wynik === undefined) odrzuc('Zmiana karty', wynik.blad);
            const migawka = wynik.wynik.snapshot;
            return migawka === undefined
              ? `Karta ${wynik.wynik.tab.id} jest czynna.`
              : `Karta ${wynik.wynik.tab.id} czynna; strona ${migawka.url} wczytana ponownie.`;
          },
        },
        {
          nazwa: 'Zamknij kartę',
          async wykonaj(wartosci) {
            const wynik = await zrodlo.zamknijKarte({
              tabId: wymagane(wartosci, 'karta', 'identyfikator karty'),
            });
            if (!wynik.udany || wynik.wynik === undefined) odrzuc('Zamknięcie karty', wynik.blad);
            return wynik.wynik.closed
              ? 'Karta zamknięta; zostaje w historii przestrzeni roboczej.'
              : 'Karta była już zamknięta — nic się nie zmieniło.';
          },
        },
        {
          nazwa: 'Grupa kart',
          async wykonaj(wartosci) {
            const karta = opcjonalne(wartosci, 'karta');
            const wynik = await zrodlo.ustawGrupeKart({
              windowId: oknoAlboOdmowa(stan),
              name: wymagane(wartosci, 'nazwa', 'nazwę grupy'),
              tabIds: karta === undefined ? undefined : [karta],
            });
            if (!wynik.udany || wynik.wynik === undefined) odrzuc('Grupa kart', wynik.blad);
            const grupa = wynik.wynik.group;
            return `Grupa ${grupa.name} (${grupa.id}) obejmuje ${(grupa.tabIds ?? []).length} kart.`;
          },
        },
        {
          nazwa: 'Zapisz przestrzeń',
          async wykonaj(wartosci) {
            const wynik = await zrodlo.zapiszPrzestrzen({
              windowId: oknoAlboOdmowa(stan),
              name: wymagane(wartosci, 'nazwa', 'nazwę przestrzeni'),
            });
            if (!wynik.udany || wynik.wynik === undefined) odrzuc('Zapis przestrzeni', wynik.blad);
            const przestrzen = wynik.wynik.workspace;
            return `Przestrzeń ${przestrzen.name} (${przestrzen.id}) pamięta ${przestrzen.tabCount} kart.`;
          },
        },
        {
          nazwa: 'Wykaz przestrzeni',
          async wykonaj() {
            const wynik = await zrodlo.wykazPrzestrzeni({ windowId: oknoAlboOdmowa(stan) });
            if (!wynik.udany || wynik.wynik === undefined) odrzuc('Wykaz przestrzeni', wynik.blad);
            const wykaz = wynik.wynik.workspaces
              .map((pozycja) => `${pozycja.name} (${pozycja.id}, kart: ${pozycja.tabCount})`)
              .join('; ');
            return wykaz === '' ? 'Przestrzeni roboczych jeszcze nie ma.' : wykaz;
          },
        },
        {
          nazwa: 'Otwórz przestrzeń',
          pola: [{ klucz: 'przestrzen', etykieta: 'Identyfikator przestrzeni' }],
          async wykonaj(wartosci) {
            const wynik = await zrodlo.otworzPrzestrzen({
              workspaceId: wymagane(wartosci, 'przestrzen', 'identyfikator przestrzeni'),
              windowId: oknoAlboOdmowa(stan),
            });
            if (!wynik.udany || wynik.wynik === undefined) odrzuc('Otwarcie przestrzeni', wynik.blad);
            return `Przywrócono ${wynik.wynik.tabs.length} kart przestrzeni ${wynik.wynik.workspace.name}.`;
          },
        },
        {
          nazwa: 'Usuń przestrzeń',
          async wykonaj(wartosci) {
            const wynik = await zrodlo.usunPrzestrzen({
              workspaceId: wymagane(wartosci, 'przestrzen', 'identyfikator przestrzeni'),
            });
            if (!wynik.udany || wynik.wynik === undefined) odrzuc('Usunięcie przestrzeni', wynik.blad);
            return 'Przestrzeń usunięta; karty zostały otwarte.';
          },
        },
      ],
    },
    {
      tytul: 'Zakładki',
      opis: 'Zakładka jest czynnością świadomą Operatora — ma folder, etykiety i własną notatkę.',
      czynnosci: [
        {
          nazwa: 'Dodaj zakładkę',
          pola: [
            { klucz: 'adresZakladki', etykieta: 'Adres strony' },
            { klucz: 'folder', etykieta: 'Folder' },
            { klucz: 'szukaj', etykieta: 'Szukany tekst' },
            { klucz: 'zakladka', etykieta: 'Identyfikator zakładki' },
          ],
          async wykonaj(wartosci) {
            const migawka = stan.migawka();
            const adres = opcjonalne(wartosci, 'adresZakladki') ?? migawka?.url ?? '';
            if (adres === '') throw new Error('Podaj adres albo przejdź najpierw pod stronę.');
            const wynik = await zrodlo.dodajZakladke({
              windowId: oknoAlboOdmowa(stan),
              url: adres,
              title: migawka?.title,
              folder: opcjonalne(wartosci, 'folder'),
            });
            if (!wynik.udany || wynik.wynik === undefined) odrzuc('Dodanie zakładki', wynik.blad);
            return `Zakładka ${wynik.wynik.bookmark.id} zapisana dla ${wynik.wynik.bookmark.url}.`;
          },
        },
        {
          nazwa: 'Wykaz zakładek',
          async wykonaj(wartosci) {
            const wynik = await zrodlo.wykazZakladek({
              windowId: oknoAlboOdmowa(stan),
              folder: opcjonalne(wartosci, 'folder'),
              query: opcjonalne(wartosci, 'szukaj'),
            });
            if (!wynik.udany || wynik.wynik === undefined) odrzuc('Wykaz zakładek', wynik.blad);
            return `Zakładek w wykazie: ${wynik.wynik.bookmarks.length}.`;
          },
        },
        {
          nazwa: 'Usuń zakładkę',
          async wykonaj(wartosci) {
            const wynik = await zrodlo.usunZakladke({
              bookmarkId: wymagane(wartosci, 'zakladka', 'identyfikator zakładki'),
            });
            if (!wynik.udany || wynik.wynik === undefined) odrzuc('Usunięcie zakładki', wynik.blad);
            return 'Zakładka zdjęta z wykazu okna.';
          },
        },
      ],
    },
    {
      tytul: 'Praca na stronie i narzędzia inspekcyjne',
      opis:
        'Przewinięcie, zrzut, drzewo DOM, konsola, rejestr sieciowy i emulacja urządzenia idą ' +
        'silnikiem przeglądarki na serwerze — to strona uruchomiona, nie samo jej źródło.',
      czynnosci: [
        {
          nazwa: 'Przewiń do końca',
          pola: [
            { klucz: 'selektor', etykieta: 'Selektor elementu' },
            { klucz: 'jakosc', etykieta: 'Jakość zrzutu (1–100)' },
          ],
          async wykonaj(wartosci) {
            const selektor = opcjonalne(wartosci, 'selektor');
            const wynik = await zrodlo.przewin({
              windowId: oknoAlboOdmowa(stan),
              toSelector: selektor,
              toEnd: selektor === undefined ? true : undefined,
            });
            if (!wynik.udany || wynik.wynik === undefined) odrzuc('Przewinięcie strony', wynik.blad);
            stan.wchlonMigawke(wynik.wynik.snapshot);
            return `Strona przewinięta; migawka po przewinięciu ma ${(wynik.wynik.snapshot.text ?? '').length} znaków.`;
          },
        },
        {
          nazwa: 'Zrzut całej strony',
          async wykonaj(wartosci) {
            const wynik = await zrodlo.wykonajZrzut({
              windowId: oknoAlboOdmowa(stan),
              mode: BrowserScreenshotMode.FullPage,
              quality: liczbaPola(wartosci, 'jakosc'),
            });
            if (!wynik.udany || wynik.wynik === undefined) odrzuc('Zrzut strony', wynik.blad);
            const zrzut = wynik.wynik.screenshot;
            return `Zrzut ${zrzut.width}×${zrzut.height} zapisany pod ${zrzut.ref}.`;
          },
        },
        {
          nazwa: 'Drzewo DOM',
          async wykonaj(wartosci) {
            const wynik = await zrodlo.zbadajDrzewo({
              windowId: oknoAlboOdmowa(stan),
              selector: opcjonalne(wartosci, 'selektor'),
              includeStyles: true,
            });
            if (!wynik.udany || wynik.wynik === undefined) odrzuc('Drzewo DOM', wynik.blad);
            return `Węzłów w drzewie: ${wynik.wynik.nodes.length}.`;
          },
        },
        {
          nazwa: 'Konsola strony',
          async wykonaj() {
            const wynik = await zrodlo.odczytajKonsole({ windowId: oknoAlboOdmowa(stan) });
            if (!wynik.udany || wynik.wynik === undefined) odrzuc('Odczyt konsoli', wynik.blad);
            const wpisy = wynik.wynik.entries;
            return wpisy.length === 0
              ? 'Strona nie wypisała ani jednego komunikatu.'
              : `Komunikatów konsoli: ${wpisy.length}; ostatni: ${wpisy[wpisy.length - 1]?.text ?? ''}`;
          },
        },
        {
          nazwa: 'Rejestr sieciowy (HAR)',
          async wykonaj() {
            const wynik = await zrodlo.rejestrSieciowy({ windowId: oknoAlboOdmowa(stan) });
            if (!wynik.udany || wynik.wynik === undefined) odrzuc('Rejestr sieciowy', wynik.blad);
            return `Żądań w rejestrze: ${wynik.wynik.entries.length}; zapis HAR: ${wynik.wynik.harRef ?? 'brak'}.`;
          },
        },
        {
          nazwa: 'Emuluj telefon',
          async wykonaj() {
            const wynik = await zrodlo.emulujUrzadzenie({
              windowId: oknoAlboOdmowa(stan),
              preset: BrowserDevicePreset.Mobile,
            });
            if (!wynik.udany || wynik.wynik === undefined) odrzuc('Emulacja urządzenia', wynik.blad);
            const metryki = wynik.wynik.metrics;
            if (wynik.wynik.snapshot !== undefined) stan.wchlonMigawke(wynik.wynik.snapshot);
            return `Strona wyrenderowana w ${metryki.width}×${metryki.height} (skala ${metryki.deviceScaleFactor}).`;
          },
        },
      ],
    },
  ];
}

/** Sekcje wpinane do Capture & Monitor Panel. */
export function sekcjeMaterialu(stan: StanPrzegladania): SekcjaRodzin[] {
  const zrodlo = stan.zrodlo;
  return [
    {
      tytul: 'Monitory zmian prowadzone przez rdzeń',
      opis:
        'Monitor rdzenia zapamiętuje treść strony jako odniesienie i mierzy różnicę przy ' +
        'sprawdzeniu — bez przestawiania wspólnego podglądu.',
      czynnosci: [
        {
          nazwa: 'Załóż monitor',
          pola: [
            { klucz: 'adresMonitora', etykieta: 'Adres pilnowanej strony' },
            { klucz: 'monitor', etykieta: 'Identyfikator monitora' },
            { klucz: 'cykl', etykieta: 'Cykl sprawdzeń (s)' },
          ],
          async wykonaj(wartosci) {
            const adres = opcjonalne(wartosci, 'adresMonitora') ?? stan.migawka()?.url ?? '';
            if (adres === '') throw new Error('Podaj adres albo przejdź najpierw pod stronę.');
            const wynik = await zrodlo.zalozMonitor({
              windowId: oknoAlboOdmowa(stan),
              url: adres,
              intervalSeconds: liczbaPola(wartosci, 'cykl'),
            });
            if (!wynik.udany || wynik.wynik === undefined) odrzuc('Założenie monitora', wynik.blad);
            const monitor = wynik.wynik.monitor;
            return `Monitor ${monitor.id} pilnuje ${monitor.url}; odniesienie ma ${monitor.baselineLength ?? 0} znaków.`;
          },
        },
        {
          nazwa: 'Wykaz monitorów',
          async wykonaj() {
            const wynik = await zrodlo.wykazMonitorow({ windowId: oknoAlboOdmowa(stan) });
            if (!wynik.udany || wynik.wynik === undefined) odrzuc('Wykaz monitorów', wynik.blad);
            return `Monitorów w oknie: ${wynik.wynik.monitors.length}.`;
          },
        },
        {
          nazwa: 'Sprawdź monitor',
          async wykonaj(wartosci) {
            const wynik = await zrodlo.sprawdzMonitor({
              monitorId: wymagane(wartosci, 'monitor', 'identyfikator monitora'),
              updateBaseline: true,
            });
            if (!wynik.udany || wynik.wynik === undefined) odrzuc('Sprawdzenie monitora', wynik.blad);
            const roznica = wynik.wynik.diff;
            return wynik.wynik.changed
              ? `Zmiana wykryta: +${roznica?.addedLines ?? 0} / −${roznica?.removedLines ?? 0} wierszy.`
              : 'Treść bez zmian wobec odniesienia.';
          },
        },
        {
          nazwa: 'Zdejmij monitor',
          async wykonaj(wartosci) {
            const wynik = await zrodlo.zdejmijMonitor({
              monitorId: wymagane(wartosci, 'monitor', 'identyfikator monitora'),
            });
            if (!wynik.udany || wynik.wynik === undefined) odrzuc('Zdjęcie monitora', wynik.blad);
            return 'Monitor zdjęty wraz z jego odniesieniem.';
          },
        },
      ],
    },
    {
      tytul: 'Kanały i kolejka czytania',
      opis:
        'Subskrypcja odpytuje kanał naprawdę: rdzeń rozpoznaje postać RSS, Atom albo JSON Feed ' +
        'i odkłada jego wpisy.',
      czynnosci: [
        {
          nazwa: 'Subskrybuj kanał',
          pola: [
            { klucz: 'adresKanalu', etykieta: 'Adres kanału' },
            { klucz: 'kanal', etykieta: 'Identyfikator kanału' },
            { klucz: 'adresCzytania', etykieta: 'Adres do przeczytania' },
            { klucz: 'pozycja', etykieta: 'Identyfikator pozycji czytania' },
          ],
          async wykonaj(wartosci) {
            const wynik = await zrodlo.subskrybujKanal({
              windowId: oknoAlboOdmowa(stan),
              url: wymagane(wartosci, 'adresKanalu', 'adres kanału'),
            });
            if (!wynik.udany || wynik.wynik === undefined) odrzuc('Subskrypcja kanału', wynik.blad);
            const kanal = wynik.wynik.feed;
            return `Kanał ${kanal.title ?? kanal.url} (${kanal.format ?? 'postać nierozpoznana'}) ma ${(kanal.entries ?? []).length} wpisów.`;
          },
        },
        {
          nazwa: 'Wykaz kanałów',
          async wykonaj() {
            const wynik = await zrodlo.wykazKanalow({
              windowId: oknoAlboOdmowa(stan),
              includeEntries: true,
            });
            if (!wynik.udany || wynik.wynik === undefined) odrzuc('Wykaz kanałów', wynik.blad);
            return `Kanałów subskrybowanych: ${wynik.wynik.feeds.length}.`;
          },
        },
        {
          nazwa: 'Zdejmij kanał',
          async wykonaj(wartosci) {
            const wynik = await zrodlo.zdejmijKanal({
              feedId: wymagane(wartosci, 'kanal', 'identyfikator kanału'),
            });
            if (!wynik.udany || wynik.wynik === undefined) odrzuc('Zdjęcie kanału', wynik.blad);
            return 'Kanał zdjęty wraz z jego wpisami.';
          },
        },
        {
          nazwa: 'Odłóż do czytania',
          async wykonaj(wartosci) {
            const adres = opcjonalne(wartosci, 'adresCzytania') ?? stan.migawka()?.url ?? '';
            if (adres === '') throw new Error('Podaj adres albo przejdź najpierw pod stronę.');
            const wynik = await zrodlo.odlozDoCzytania({
              windowId: oknoAlboOdmowa(stan),
              url: adres,
              title: stan.migawka()?.title,
            });
            if (!wynik.udany || wynik.wynik === undefined) odrzuc('Kolejka czytania', wynik.blad);
            return `Pozycja ${wynik.wynik.item.id} odłożona do przeczytania.`;
          },
        },
        {
          nazwa: 'Wykaz kolejki',
          async wykonaj() {
            const wynik = await zrodlo.wykazCzytania({ windowId: oknoAlboOdmowa(stan) });
            if (!wynik.udany || wynik.wynik === undefined) odrzuc('Wykaz kolejki', wynik.blad);
            return `Pozycji w kolejce czytania: ${wynik.wynik.items.length}.`;
          },
        },
        {
          nazwa: 'Oznacz przeczytane',
          async wykonaj(wartosci) {
            const wynik = await zrodlo.zdejmijZCzytania({
              itemId: wymagane(wartosci, 'pozycja', 'identyfikator pozycji'),
              markRead: true,
            });
            if (!wynik.udany || wynik.wynik === undefined) odrzuc('Kolejka czytania', wynik.blad);
            return wynik.wynik.item === undefined
              ? 'Pozycja zdjęta z kolejki.'
              : 'Pozycja oznaczona jako przeczytana i zostaje w kolejce jako ślad.';
          },
        },
      ],
    },
    {
      tytul: 'Pobrania i wytwory sesji',
      opis:
        'Pobranie powstaje przez wejście pod adres pliku, którego nie da się pokazać jako ' +
        'strony. Wytwór odkłada bajty w magazynie rdzenia pod sumą kontrolną.',
      czynnosci: [
        {
          nazwa: 'Wykaz pobrań',
          pola: [
            { klucz: 'pobranie', etykieta: 'Identyfikator pobrania' },
            { klucz: 'tytulWytworu', etykieta: 'Tytuł wytworu' },
          ],
          async wykonaj() {
            const wynik = await zrodlo.wykazPobran({ windowId: oknoAlboOdmowa(stan) });
            if (!wynik.udany || wynik.wynik === undefined) odrzuc('Wykaz pobrań', wynik.blad);
            return `Pobrań w wykazie: ${wynik.wynik.downloads.length}.`;
          },
        },
        {
          nazwa: 'Ponów pobranie',
          async wykonaj(wartosci) {
            const wynik = await zrodlo.sterujPobraniem({
              downloadId: wymagane(wartosci, 'pobranie', 'identyfikator pobrania'),
              action: BrowserDownloadAction.Retry,
            });
            if (!wynik.udany || wynik.wynik === undefined) odrzuc('Sterowanie pobraniem', wynik.blad);
            const pobranie = wynik.wynik.download;
            return `Pobranie ${pobranie.id}: stan ${pobranie.status}, odebrano ${pobranie.receivedBytes ?? 0} B.`;
          },
        },
        {
          nazwa: 'Przerwij pobranie',
          async wykonaj(wartosci) {
            const wynik = await zrodlo.sterujPobraniem({
              downloadId: wymagane(wartosci, 'pobranie', 'identyfikator pobrania'),
              action: BrowserDownloadAction.Cancel,
            });
            if (!wynik.udany || wynik.wynik === undefined) odrzuc('Sterowanie pobraniem', wynik.blad);
            return `Pobranie przerwane; stan: ${wynik.wynik.download.status}.`;
          },
        },
        {
          nazwa: 'Zapisz treść strony jako wytwór',
          async wykonaj(wartosci) {
            const migawka = stan.migawka();
            if (migawka === null) throw new Error('Najpierw przejdź pod stronę — nie ma czego zapisać.');
            const tresc = migawka.html ?? migawka.text ?? '';
            if (tresc === '') throw new Error('Migawka nie niesie treści — pobierz ją ze źródłem strony.');
            const wynik = await zrodlo.dodajWytwor({
              windowId: oknoAlboOdmowa(stan),
              kind: BrowserArtifactKind.Archive,
              title: opcjonalne(wartosci, 'tytulWytworu') ?? migawka.title,
              contentBase64: btoa(unescape(encodeURIComponent(tresc))),
              mimeType: 'text/html',
              sourceUrl: migawka.url,
              snapshotId: migawka.id,
            });
            if (!wynik.udany || wynik.wynik === undefined) odrzuc('Zapis wytworu', wynik.blad);
            const wytwor = wynik.wynik.artifact;
            return `Wytwór ${wytwor.id} zapisany (${wytwor.sizeBytes ?? 0} B) pod ${wytwor.contentRef}.`;
          },
        },
        {
          nazwa: 'Odczytaj zrzut migawki',
          async wykonaj() {
            const migawka = stan.migawka();
            if (migawka === null) throw new Error('Najpierw przejdź pod stronę.');
            const wynik = await zrodlo.odczytajZrzut({ snapshotId: migawka.id });
            if (!wynik.udany || wynik.wynik === undefined) odrzuc('Odczyt zrzutu', wynik.blad);
            const zrzut = wynik.wynik.screenshot;
            return `Zrzut ${zrzut.width}×${zrzut.height} (${zrzut.format}) pod ${zrzut.ref}.`;
          },
        },
      ],
    },
  ];
}

/** Sekcje wpinane do Notes Panel i Sources Panel. */
export function sekcjePorzadku(stan: StanPrzegladania): SekcjaRodzin[] {
  const zrodlo = stan.zrodlo;
  return [
    {
      tytul: 'Zestawy źródeł i wątki notatek',
      opis:
        'Oznaczenia idą do rdzenia, więc przeżywają przeładowanie karty: zestaw tematyczny ' +
        'źródła, wątek notatki, jej klasyfikacja i przypięcie.',
      czynnosci: [
        {
          nazwa: 'Zestaw źródeł',
          pola: [
            { klucz: 'nazwaZestawu', etykieta: 'Nazwa zestawu albo wątku' },
            { klucz: 'zrodlo', etykieta: 'Identyfikator źródła' },
            { klucz: 'notatka', etykieta: 'Identyfikator notatki' },
          ],
          async wykonaj(wartosci) {
            const wskazane = opcjonalne(wartosci, 'zrodlo');
            const wynik = await zrodlo.ustawZestawZrodel({
              windowId: oknoAlboOdmowa(stan),
              name: wymagane(wartosci, 'nazwaZestawu', 'nazwę zestawu'),
              sourceIds: wskazane === undefined ? undefined : [wskazane],
            });
            if (!wynik.udany || wynik.wynik === undefined) odrzuc('Zestaw źródeł', wynik.blad);
            const zestaw = wynik.wynik.group;
            return `Zestaw ${zestaw.name} (${zestaw.id}) obejmuje ${(zestaw.sourceIds ?? []).length} źródeł.`;
          },
        },
        {
          nazwa: 'Wykaz zestawów',
          async wykonaj() {
            const wynik = await zrodlo.wykazZestawowZrodel({ windowId: oknoAlboOdmowa(stan) });
            if (!wynik.udany || wynik.wynik === undefined) odrzuc('Wykaz zestawów', wynik.blad);
            return `Zestawów tematycznych: ${wynik.wynik.groups.length}.`;
          },
        },
        {
          nazwa: 'Usuń źródło',
          async wykonaj(wartosci) {
            const wynik = await zrodlo.usunZrodlo({
              sourceId: wymagane(wartosci, 'zrodlo', 'identyfikator źródła'),
            });
            if (!wynik.udany || wynik.wynik === undefined) odrzuc('Usunięcie źródła', wynik.blad);
            await stan.zaciagnijZebrane();
            return 'Źródło zdjęte z wykazu okna.';
          },
        },
        {
          nazwa: 'Wątek notatek',
          async wykonaj(wartosci) {
            const wskazana = opcjonalne(wartosci, 'notatka');
            const wynik = await zrodlo.ustawWatekNotatek({
              windowId: oknoAlboOdmowa(stan),
              name: wymagane(wartosci, 'nazwaZestawu', 'nazwę wątku'),
              noteIds: wskazana === undefined ? undefined : [wskazana],
            });
            if (!wynik.udany || wynik.wynik === undefined) odrzuc('Wątek notatek', wynik.blad);
            const watek = wynik.wynik.thread;
            return `Wątek ${watek.name} (${watek.id}) obejmuje ${(watek.noteIds ?? []).length} notatek.`;
          },
        },
        {
          nazwa: 'Wykaz wątków',
          async wykonaj() {
            const wynik = await zrodlo.wykazWatkowNotatek({ windowId: oknoAlboOdmowa(stan) });
            if (!wynik.udany || wynik.wynik === undefined) odrzuc('Wykaz wątków', wynik.blad);
            return `Wątków tematycznych: ${wynik.wynik.threads.length}.`;
          },
        },
        {
          nazwa: 'Oznacz notatkę jako wniosek',
          async wykonaj(wartosci) {
            const wynik = await zrodlo.zmienNotatke({
              noteId: wymagane(wartosci, 'notatka', 'identyfikator notatki'),
              classification: BrowserNoteClassification.Conclusion,
              pinned: true,
            });
            if (!wynik.udany || wynik.wynik === undefined) odrzuc('Zmiana notatki', wynik.blad);
            await stan.zaciagnijZebrane();
            return `Notatka ${wynik.wynik.note.id} oznaczona jako wniosek i przypięta.`;
          },
        },
      ],
    },
  ];
}

/** Sekcje wpinane do Automation Studio. */
export function sekcjeAutomatyzacji(stan: StanPrzegladania): SekcjaRodzin[] {
  const zrodlo = stan.zrodlo;
  return [
    {
      tytul: 'Nagrywarka makr i granice Wykonawcy',
      opis:
        'Makro zapisuje kroki w kształcie kroków automatyzacji platformy. Granice Wykonawcy ' +
        'obowiązują przebiegi na stronie i są utrwalane w rdzeniu, nie w karcie.',
      czynnosci: [
        {
          nazwa: 'Rozpocznij nagrywanie',
          pola: [
            { klucz: 'nazwaMakra', etykieta: 'Nazwa makra' },
            { klucz: 'makro', etykieta: 'Identyfikator makra' },
            { klucz: 'kroki', etykieta: 'Limit kroków Wykonawcy' },
          ],
          async wykonaj(wartosci) {
            const wynik = await zrodlo.nagrywajMakro({
              windowId: oknoAlboOdmowa(stan),
              action: BrowserMacroAction.Start,
              name: opcjonalne(wartosci, 'nazwaMakra'),
            });
            if (!wynik.udany || wynik.wynik === undefined) odrzuc('Nagrywanie makra', wynik.blad);
            return `Nagrywanie makra ${wynik.wynik.macro.id} rozpoczęte.`;
          },
        },
        {
          nazwa: 'Dopisz krok z bieżącej strony',
          async wykonaj(wartosci) {
            const migawka = stan.migawka();
            if (migawka === null) throw new Error('Najpierw przejdź pod stronę — krok powstaje z migawki.');
            const wynik = await zrodlo.nagrywajMakro({
              windowId: oknoAlboOdmowa(stan),
              action: BrowserMacroAction.Step,
              macroId: wymagane(wartosci, 'makro', 'identyfikator makra'),
              step: {
                id: migawka.id,
                kind: 'command',
                name: migawka.title ?? migawka.url,
                command: Command.BrowserNavigate,
                params: { url: migawka.url },
              },
            });
            if (!wynik.udany || wynik.wynik === undefined) odrzuc('Krok makra', wynik.blad);
            return `Makro ma ${(wynik.wynik.macro.steps ?? []).length} kroków.`;
          },
        },
        {
          nazwa: 'Zakończ nagrywanie',
          async wykonaj(wartosci) {
            const wynik = await zrodlo.nagrywajMakro({
              windowId: oknoAlboOdmowa(stan),
              action: BrowserMacroAction.Stop,
              macroId: wymagane(wartosci, 'makro', 'identyfikator makra'),
            });
            if (!wynik.udany || wynik.wynik === undefined) odrzuc('Nagrywanie makra', wynik.blad);
            return `Makro zamknięte z ${(wynik.wynik.macro.steps ?? []).length} krokami.`;
          },
        },
        {
          nazwa: 'Zapisz granice Wykonawcy',
          async wykonaj(wartosci) {
            const wynik = await zrodlo.ustawGraniceWykonawcy({
              windowId: oknoAlboOdmowa(stan),
              maxSteps: liczbaPola(wartosci, 'kroki'),
            });
            if (!wynik.udany || wynik.wynik === undefined) odrzuc('Granice Wykonawcy', wynik.blad);
            const granice = wynik.wynik.limits;
            return `Granice zasięgu ${granice.scope}: ${granice.maxSteps} kroków, ${granice.maxDurationSeconds} s.`;
          },
        },
        {
          nazwa: 'Odczytaj granice',
          async wykonaj() {
            const wynik = await zrodlo.odczytajGraniceWykonawcy({ windowId: oknoAlboOdmowa(stan) });
            if (!wynik.udany || wynik.wynik === undefined) odrzuc('Granice Wykonawcy', wynik.blad);
            const granice = wynik.wynik.limits;
            return `Obowiązuje: ${granice.maxSteps} kroków, ${granice.maxDurationSeconds} s, ` +
              `potwierdzanie wysyłki: ${granice.confirmBeforeSubmit ? 'tak' : 'nie'}.`;
          },
        },
      ],
    },
  ];
}
