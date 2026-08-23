import type { WorkspaceDashboard } from '../../../../shared/contract';
import type { Odsubskrybuj } from '../../polaczenie/magistrala-zdarzen';
import { POWOD_PRZED_WCZYTANIEM, type OknoRozmowyModulu } from './okno-rozmowy';
import type { ZrodloWorkspace } from './zrodlo-workspace';

/**
 * Projekt bieżący modułu Workspace — jedna prawda dla pięciu okien.
 *
 * Wszystkie okna modułu pracują w tym samym projekcie: pulpit go pokazuje,
 * pamięć zapisuje w nim ustalenia, biblioteka czyta jego pliki, a Agent Manager
 * przypisuje do niego ekspertów. Gdyby każde okno trzymało własne wskazanie,
 * zmiana projektu w jednym oknie rozjechałaby pozostałe cztery. Stan jest więc
 * jeden i to on rozsyła powiadomienie o zmianie.
 *
 * Zdarzenie rdzenia wchodzi tą samą drogą: `workspace.project.changed`
 * przychodzi także z pracy innego okna albo innego urządzenia tego konta —
 * stan przyjmuje je jak zmianę własną i odświeża okna, zamiast czekać na ruch
 * Operatora.
 */
export interface StanProjektu {
  /** Identyfikator projektu bieżącego; pusty znaczy „nie wskazano”. */
  projekt(): string;
  /** Przestawia moduł na inny projekt i powiadamia okna. */
  ustawProjekt(idProjektu: string): void;
  /** Okno rozmowy modułu — nośnik przeniesienia kontekstu. */
  oknoRozmowy(): string;
  /**
   * Zapisuje okno rozmowy odnalezione w karcie sesji albo powód, dla którego
   * go nie ma. Jedno wejście na oba przypadki: identyfikator bez
   * powodu zostawiłby odmowę kontrolek przy zdaniu ogólnym „moduł nie zna jego
   * identyfikatora”, które nie mówi, czy rdzeń odmówił, czy okna po prostu
   * jeszcze nie ma.
   */
  ustawOknoRozmowy(znalezione: OknoRozmowyModulu): void;
  /** Powód braku okna rozmowy; pusty, gdy okno jest znane. */
  powodBrakuOkna(): string;
  /**
   * Karta sesji, w której moduł stoi; pusta, dopóki powłoka nie wczytała modułu.
   *
   * Projekt i karta sesji to dwa różne byty i moduł trzyma je osobno: pamięć
   * projektu jest wspólna wszystkim kartom, a poziomy pamięci (`memory.toggle`)
   * przestawia się karcie. To pole daje panelowi poziomów adresata komendy.
   */
  sesja(): string;
  /** Zapisuje kartę sesji podaną przez powłokę przy wczytaniu modułu. */
  ustawSesje(idSesji: string): void;
  /** Ostatnie zestawienie pulpitu; puste do pierwszego odczytu. */
  pulpit(): WorkspaceDashboard | null;
  /** Zapisuje zestawienie pulpitu i powiadamia okna zależne. */
  ustawPulpit(zestawienie: WorkspaceDashboard): void;
  /** Subskrypcja zmiany projektu albo jego zestawienia. */
  naZmiane(sluchacz: () => void): Odsubskrybuj;
  /** Zamyka nasłuch zdarzenia rdzenia. */
  zamknij(): void;
}

/** Zależności stanu: źródło komend i okno rozmowy modułu. */
export interface OpcjeStanu {
  projekt?: string;
  oknoRozmowy?: string;
}

export function utworzStanProjektu(zrodlo: ZrodloWorkspace, opcje: OpcjeStanu = {}): StanProjektu {
  const sluchacze = new Set<() => void>();
  let idProjektu = opcje.projekt ?? '';
  let idSesji = '';
  // Okno rozmowy podane przy montażu wygrywa z szukaniem: wołający, który je
  // zna, wie o nim więcej niż wykaz okien sesji. Domyślnie pole jest puste,
  // a powód mówi wprost, że nikt jeszcze nie szukał — nie że okna nie ma.
  let idOknaRozmowy = opcje.oknoRozmowy ?? '';
  let powodBraku = idOknaRozmowy === '' ? POWOD_PRZED_WCZYTANIEM : '';
  let zestawienie: WorkspaceDashboard | null = null;

  function powiadom(): void {
    for (const sluchacz of [...sluchacze]) sluchacz();
  }

  const odsubskrybuj = zrodlo.naZmianeProjektu((tresc) => {
    if (tresc.project.id !== idProjektu) return;
    // Zmiana przyszła z rdzenia: zestawienie w oknie jest już nieaktualne,
    // więc kasujemy je i prosimy okna o ponowny odczyt.
    zestawienie = null;
    powiadom();
  });

  return {
    projekt: () => idProjektu,

    ustawProjekt(nowy) {
      const przyciety = nowy.trim();
      if (przyciety === idProjektu) return;
      idProjektu = przyciety;
      zestawienie = null;
      powiadom();
    },

    oknoRozmowy: () => idOknaRozmowy,

    ustawOknoRozmowy(znalezione) {
      const nowe = znalezione.okno.trim();
      if (nowe === idOknaRozmowy && znalezione.powod === powodBraku) return;
      idOknaRozmowy = nowe;
      powodBraku = nowe === '' ? znalezione.powod : '';
      powiadom();
    },

    powodBrakuOkna: () => powodBraku,

    sesja: () => idSesji,

    ustawSesje(nowa) {
      const przycieta = nowa.trim();
      if (przycieta === idSesji) return;
      idSesji = przycieta;
      // Okno rozmowy należy do karty, nie do modułu: karta zmieniona unieważnia
      // identyfikator odnaleziony w poprzedniej. Trzymanie go dalej wysłałoby
      // `context.transfer` z oknem cudzej karty — z odpowiedzią udaną i skutkiem
      // w niewłaściwym miejscu.
      idOknaRozmowy = opcje.oknoRozmowy ?? '';
      powodBraku =
        idOknaRozmowy === '' ? 'Karta sesji zmieniona — moduł szuka w niej okna rozmowy.' : '';
      powiadom();
    },

    pulpit: () => zestawienie,

    ustawPulpit(nowe) {
      zestawienie = nowe;
      powiadom();
    },

    naZmiane(sluchacz) {
      sluchacze.add(sluchacz);
      return () => sluchacze.delete(sluchacz);
    },

    zamknij() {
      sluchacze.clear();
      odsubskrybuj();
    },
  };
}
