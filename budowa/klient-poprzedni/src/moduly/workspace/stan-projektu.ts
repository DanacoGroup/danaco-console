import type { WorkspaceDashboard } from '../../../../shared/contract';
import type { Odsubskrybuj } from '../../polaczenie/magistrala-zdarzen';
import { POWOD_PRZED_WCZYTANIEM, type OknoRozmowyModulu } from './okno-rozmowy';
import type { ZrodloWorkspace } from './zrodlo-workspace';

/**
 * Projekt bieżący modułu Workspace jest jedną prawdą dla pięciu okien: pulpit go pokazuje, pamięć
 * zapisuje ustalenia, biblioteka czyta pliki, a Agent Manager przypisuje ekspertów.
 */
export interface StanProjektu {
  /** Identyfikator projektu bieżącego; pusty znaczy „nie wskazano”. */
  projekt(): string;
  /** Przestawia moduł na inny projekt i powiadamia okna. */
  ustawProjekt(idProjektu: string): void;
  /** Okno rozmowy modułu — nośnik przeniesienia kontekstu. */
  oknoRozmowy(): string;
  /** Zapisuje okno rozmowy odnalezione w karcie sesji albo powód, dla którego go nie ma, jednym wejściem. */
  ustawOknoRozmowy(znalezione: OknoRozmowyModulu): void;
  /** Powód braku okna rozmowy; pusty, gdy okno jest znane. */
  powodBrakuOkna(): string;
  /** Karta sesji, w której moduł stoi, jest pusta do wczytania modułu; daje adresata komendy poziomów. */
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

/** Zależności stanu niosą źródło komend workspace oraz okno rozmowy modułu wskazane przy montażu stanu projektu. */
export interface OpcjeStanu {
  projekt?: string;
  oknoRozmowy?: string;
}

export function utworzStanProjektu(zrodlo: ZrodloWorkspace, opcje: OpcjeStanu = {}): StanProjektu {
  const sluchacze = new Set<() => void>();
  let idProjektu = opcje.projekt ?? '';
  let idSesji = '';
  // Okno rozmowy podane przy montażu wygrywa z szukaniem; pole puste znaczy, że nikt jeszcze nie szukał.
  let idOknaRozmowy = opcje.oknoRozmowy ?? '';
  let powodBraku = idOknaRozmowy === '' ? POWOD_PRZED_WCZYTANIEM : '';
  let zestawienie: WorkspaceDashboard | null = null;

  function powiadom(): void {
    for (const sluchacz of [...sluchacze]) sluchacz();
  }

  const odsubskrybuj = zrodlo.naZmianeProjektu((tresc) => {
    if (tresc.project.id !== idProjektu) return;
    // Zmiana przyszła z rdzenia: zestawienie w oknie jest już nieaktualne, kasujemy je i prosimy o odczyt.
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
      // Okno rozmowy należy do karty, nie do modułu: zmiana karty unieważnia identyfikator z poprzedniej.
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
