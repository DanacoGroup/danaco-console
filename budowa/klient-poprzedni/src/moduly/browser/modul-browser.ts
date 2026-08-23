import './browser.css';

import type { Kanal } from '../../protokol/kanal';
import { utworzKatalogOkien } from '../katalog-okien';
import { KOD_MODULU, KODY_OKIEN, KODY_PANELI, WYZWALACZE } from './etykiety-browser';
import { utworzMacierzIzolacji } from './macierz-izolacji';
import { utworzOknoAutomationStudio } from './okno-automation-studio';
import { utworzOknoBrowserWindow } from './okno-browser-window';
import { utworzOknoCaptureMonitor } from './okno-capture-monitor';
import { utworzOknoNotesPanel } from './okno-notes-panel';
import { utworzOknoSourcesPanel } from './okno-sources-panel';
import { utworzPasekKontekstu } from './pasek-kontekstu';
import { utworzStanPrzegladania } from './stan-przegladania';

/**
 * Moduł Browser — pięć okien operacyjnych osadzonych w jednym układzie.
 *
 * Układ wynika z roli okna i z warstwy widoczności. Browser Window jest oknem
 * wiodącym i punktem wejścia modułu, więc stoi w pasie pierwszym na całą
 * szerokość i nie chowa się nigdy. Sources Panel, Notes Panel, Automation
 * Studio, Capture & Monitor Panel oraz macierz izolacji są rozszerzeniami
 * bocznymi: w spoczynku zwinięte, otwierane wyzwalaczem paska kontekstu albo
 * skrótem klawiszowym (`warstwy-widocznosci.ts`).
 *
 * Stan jest jeden na cały moduł: zaznaczenie zrobione w podglądzie strony
 * trafia do formularza notatki, źródło dodane w panelu źródeł pojawia się
 * w wyborze powiązania notatki, a strona przechwycona w Capture & Monitor
 * Panel przestawia wspólny podgląd — bo `stan-przegladania` jest jeden.
 *
 * Moduł nie osadza się sam i nie zna powłoki: oddaje element, a warstwa
 * składająca decyduje, gdzie go postawić.
 */
export interface ModulBrowser {
  /** Element osadzany w obszarze roboczym powłoki. */
  element: HTMLElement;
  /** Ustala okno przeglądarki sesji i odczytuje jego stan. */
  wczytaj(idSesji: string): Promise<void>;
  /** Odłącza subskrypcję zdarzeń rdzenia i skróty warstw widoczności. */
  rozlacz(): void;
}

export function utworzModulBrowser(kanal: Kanal): ModulBrowser {
  const stan = utworzStanPrzegladania(kanal);

  const notatki = utworzOknoNotesPanel(stan);
  const zrodla = utworzOknoSourcesPanel(stan);
  const przegladarka = utworzOknoBrowserWindow(stan, (fragment) =>
    notatki.przygotujZFragmentu(fragment),
  );
  const automatyzacja = utworzOknoAutomationStudio(stan);
  const materialy = utworzOknoCaptureMonitor(stan);
  const izolacja = utworzMacierzIzolacji(stan);

  const pasRozszerzen = document.createElement('div');
  pasRozszerzen.className = 'mb-modul__pas mb-modul__pas--rozszerzenia';
  pasRozszerzen.append(
    zrodla.element,
    notatki.element,
    automatyzacja.element,
    materialy.element,
    izolacja.element,
  );

  // Kolejność rejestracji jest kolejnością wyzwalaczy w pasku kontekstu:
  // najpierw warstwa druga (znaczniki źródeł i notatek), potem trzecia
  // (menu operacji), na końcu czwarta — widoczna dopiero w trybie
  // administracyjnym, ale osiągalna skrótem zawsze.
  stan.warstwy.zarejestruj({
    kod: KODY_OKIEN.zrodla,
    nazwa: WYZWALACZE.zrodla,
    warstwa: 2,
    element: zrodla.element,
    licznik: () => stan.zebrane.zrodla().length,
  });
  stan.warstwy.zarejestruj({
    kod: KODY_OKIEN.notatki,
    nazwa: WYZWALACZE.notatki,
    warstwa: 2,
    element: notatki.element,
    licznik: () => stan.zebrane.notatki().length,
  });
  stan.warstwy.zarejestruj({
    kod: KODY_OKIEN.automatyzacja,
    nazwa: WYZWALACZE.automatyzacja,
    warstwa: 3,
    element: automatyzacja.element,
  });
  stan.warstwy.zarejestruj({
    kod: KODY_OKIEN.materialy,
    nazwa: WYZWALACZE.materialy,
    warstwa: 3,
    element: materialy.element,
    licznik: () => stan.material.przechwycenia().length + stan.material.monitory().length,
  });
  stan.warstwy.zarejestruj({
    kod: KODY_PANELI.macierzIzolacji,
    nazwa: WYZWALACZE.izolacja,
    warstwa: 4,
    element: izolacja.element,
  });
  for (const panel of [...przegladarka.panele, ...automatyzacja.panele]) {
    stan.warstwy.zarejestruj(panel);
  }

  const kontekst = utworzPasekKontekstu(stan);

  // Pasek uczciwości bierze zdanie z bytu wspólnego (`moduly/katalog-okien.ts`),
  // tego samego dla wszystkich modułów, więc zdanie o rozjeździe nie rozjedzie
  // się między modułami po cichu. Byt wspólny wypowiada także okna, które
  // katalog rdzenia modułowi przypisuje, a moduł ich nie buduje — oraz te,
  // które moduł buduje, a rdzeń mu ich nie przypisał.
  const katalog = utworzKatalogOkien(kanal, KOD_MODULU, [
    KODY_OKIEN.przegladarka,
    KODY_OKIEN.zrodla,
    KODY_OKIEN.notatki,
    KODY_OKIEN.automatyzacja,
    KODY_OKIEN.materialy,
  ]);
  // Element powstaje z bytu i sam przerysowuje się po każdej zmianie katalogu —
  // także po odczycie cudzym, na przykład powłoki budującej nawigację. Treść
  // wpisana tutaj byłaby dokładnie tym, co pasek ma tropić: twierdzeniem
  // o stanie rdzenia wypowiedzianym bez zapytania rdzenia.
  const uczciwosc = katalog.zdanieElement('dn-pole-opis mb-uczciwosc');

  const element = document.createElement('div');
  element.className = 'mb-modul';
  element.dataset['modul'] = 'browser';
  element.setAttribute('aria-label', 'Moduł Browser — okna operacyjne');
  element.append(kontekst.element, przegladarka.element, pasRozszerzen, uczciwosc);

  // Jeden odczyt na cały kanał: jeśli o katalog zapytał już inny moduł albo
  // powłoka, to wywołanie nie wyśle drugiego zapytania. Tak samo wykaz komend
  // rdzenia, z którego pozycje modułu biorą powód swojego bezruchu.
  void katalog.odczytaj();
  void stan.pokrycie.odczytaj();

  const odsubskrybuj = stan.obserwuj(() => odswiezWszystkie());
  const odepnijSkroty = stan.warstwy.podepnijSkroty(element);

  function odswiezWszystkie(): void {
    kontekst.odswiez();
    przegladarka.odswiez();
    zrodla.odswiez();
    notatki.odswiez();
    automatyzacja.odswiez();
    materialy.odswiez();
  }

  odswiezWszystkie();

  return {
    element,

    async wczytaj(idSesji) {
      await stan.ustalOkno(idSesji);
      // Stan okna czytany dopiero po ustaleniu okna: `window.state.get`
      // wymaga identyfikatora, którego przed `window.list` nie ma.
      await przegladarka.wczytajStan();
    },

    rozlacz() {
      odsubskrybuj();
      odepnijSkroty();
      // Pasek odpina się od wspólnego katalogu: bez tego wpis kanału trzymałby
      // przerysowanie elementu zdjętego już z drzewa. Warstwa adnotacji odpina
      // obserwatora rozmiaru płótna z tego samego powodu.
      katalog.zamknij();
      przegladarka.rozlacz();
      stan.rozlacz();
    },
  };
}
