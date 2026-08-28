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
 * Moduł Browser — pięć okien operacyjnych osadzonych w jednym układzie. Moduł
 * nie osadza się sam i nie zna powłoki: oddaje element, a warstwa składająca
 * decyduje, gdzie go postawić.
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

  // Kolejność rejestracji jest kolejnością wyzwalaczy w pasku kontekstu, od warstwy drugiej do czwartej.
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

  // Pasek uczciwości bierze zdanie z bytu wspólnego, żeby rozjazd nie rozjechał się po cichu.
  const katalog = utworzKatalogOkien(kanal, KOD_MODULU, [
    KODY_OKIEN.przegladarka,
    KODY_OKIEN.zrodla,
    KODY_OKIEN.notatki,
    KODY_OKIEN.automatyzacja,
    KODY_OKIEN.materialy,
  ]);
  // Element powstaje z bytu i sam przerysowuje się po każdej zmianie katalogu, także cudzej.
  const uczciwosc = katalog.zdanieElement('dn-pole-opis mb-uczciwosc');

  const element = document.createElement('div');
  element.className = 'mb-modul';
  element.dataset['modul'] = 'browser';
  element.setAttribute('aria-label', 'Moduł Browser — okna operacyjne');
  element.append(kontekst.element, przegladarka.element, pasRozszerzen, uczciwosc);

  // Jeden odczyt na kanał — jeśli o katalog zapytał już inny moduł, wywołanie nie powtórzy zapytania.
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
      // Stan okna czytany dopiero po ustaleniu okna — wymaga identyfikatora, którego wcześniej nie ma.
      await przegladarka.wczytajStan();
    },

    rozlacz() {
      odsubskrybuj();
      odepnijSkroty();
      // Pasek odpina się od katalogu, żeby wpis kanału nie przerysowywał elementu zdjętego z drzewa.
      katalog.zamknij();
      przegladarka.rozlacz();
      stan.rozlacz();
    },
  };
}
