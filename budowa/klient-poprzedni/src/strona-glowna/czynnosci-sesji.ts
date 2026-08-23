import { SessionStatus } from '../../../shared/contract';
import type { WpisSesji } from './zrodlo-sesji';

/**
 * Czynności historii sesji — co strefa umie zlecić rdzeniowi.
 *
 * Nazywa operacje historii w języku widoku i rozstrzyga, które z nich mają sens
 * dla danego wiersza. Plik nie zna kanału ani kontraktu — wykonanie podaje
 * wpięcie.
 *
 * Warunki stoją tutaj, a nie w karcie: karta buduje wiersz, menu buduje
 * przyciski, i żadne z nich nie rozstrzyga, czy wznowienie sesji czynnej ma
 * sens. Jedno miejsce daje jedną odpowiedź i jedno miejsce do poprawienia, gdy
 * rdzeń zmieni stany sesji.
 *
 * Brak czynności zdejmuje przycisk, zamiast go wyszarzać: przycisk widoczny,
 * którego naciśnięcie nic nie robi, jest atrapą. Wiersz bez czynności jest
 * czysto informacyjny.
 */

/** Wynik czynności w postaci zdania na ekran; `null` znaczy „bez meldunku”. */
export type MeldunekCzynnosci = string | null;

/** Jedna czynność historii; obietnica trzyma przycisk zajęty do końca. */
export type Wykonanie = (wpis: WpisSesji) => Promise<MeldunekCzynnosci>;

/**
 * Komplet czynności podawany strefie. Pominięta czynność znika z menu —
 * montaż bez tożsamości klienta może podać wykaz czysto informacyjny.
 */
export interface CzynnosciSesji {
  /** `session.rename` — nowa nazwa sesji w historii. */
  zmienNazwe?: Wykonanie;
  /** `session.copy` — kopia sesji wraz z zapisem. */
  skopiuj?: Wykonanie;
  /** `session.project.set` — przypisanie sesji do projektu. */
  przypiszProjekt?: Wykonanie;
  /** `session.project.clear` — wyjęcie sesji z projektu. */
  wyjmijZProjektu?: Wykonanie;
  /** `session.archive` — odłożenie sesji do archiwum. */
  archiwizuj?: Wykonanie;
  /** `session.restore` — przywrócenie sesji z archiwum. */
  przywroc?: Wykonanie;
  /** `session.resume` — wznowienie sesji wstrzymanej albo zakończonej. */
  wznow?: Wykonanie;
  /** `session.stop` — zatrzymanie tur biegnących w oknach sesji. */
  zatrzymaj?: Wykonanie;
}

/** Pozycja menu gotowa do zbudowania: napis, wykonanie i postać przycisku. */
export interface PozycjaCzynnosci {
  klucz: string;
  napis: string;
  /** Zdanie dymka [?] — co komenda naprawdę robi. */
  wyjasnienie: string;
  wykonaj: Wykonanie;
  /** Postać przycisku z biblioteki komponentów; brak znaczy postać domyślną. */
  postac?: 'dn-btn--zarys' | 'dn-btn--niebezpieczny';
}

/**
 * Czynności sensowne dla tego wiersza, w kolejności wyświetlania.
 *
 * Kolejność idzie od najczęstszej do najrzadszej i kończy się na czynnościach
 * zmieniających miejsce sesji w historii — archiwizacja stoi na końcu, bo
 * zabiera wiersz z wykazu.
 */
export function czynnosciWiersza(
  wpis: WpisSesji,
  czynnosci: CzynnosciSesji,
): readonly PozycjaCzynnosci[] {
  const stan = wpis.sesja.status;
  const zarchiwizowana = stan === SessionStatus.Archived;
  const biegnie = (wpis.obecnosc?.streamingWindowCount ?? 0) > 0;
  const wProjekcie = (wpis.sesja.projectId ?? '') !== '';

  const wykaz: PozycjaCzynnosci[] = [];
  const dodaj = (
    warunek: boolean,
    wykonanie: Wykonanie | undefined,
    pozycja: Omit<PozycjaCzynnosci, 'wykonaj'>,
  ): void => {
    if (warunek && wykonanie !== undefined) wykaz.push({ ...pozycja, wykonaj: wykonanie });
  };

  dodaj(!zarchiwizowana, czynnosci.zmienNazwe, {
    klucz: 'nazwa',
    napis: 'Zmień nazwę',
    wyjasnienie: 'Zmienia nazwę sesji w historii. Zapis rozmowy zostaje nietknięty.',
  });
  dodaj(!zarchiwizowana, czynnosci.skopiuj, {
    klucz: 'kopia',
    napis: 'Kopiuj',
    wyjasnienie: 'Zakłada osobną sesję z kopią zapisu. Kopia nie jest odnośnikiem do źródła.',
  });

  // Zatrzymanie ma sens wyłącznie przy turze w biegu — zero okien
  // strumieniujących znaczy, że komenda nie miałaby czego zatrzymać.
  dodaj(biegnie, czynnosci.zatrzymaj, {
    klucz: 'stop',
    napis: 'Zatrzymaj tury',
    wyjasnienie: 'Przerywa tury biegnące w oknach tej sesji. Sesja zostaje otwarta.',
    postac: 'dn-btn--niebezpieczny',
  });
  // Wznowienie dotyczy sesji, która nie biegnie — wstrzymanej albo zakończonej.
  // Sesja czynna nie ma czego wznawiać.
  dodaj(
    stan === SessionStatus.Paused || stan === SessionStatus.Finished,
    czynnosci.wznow,
    {
      klucz: 'wznow',
      napis: 'Wznów',
      wyjasnienie: 'Otwiera sesję wraz z jej oknami komunikacji.',
    },
  );

  dodaj(!zarchiwizowana && !wProjekcie, czynnosci.przypiszProjekt, {
    klucz: 'projekt',
    napis: 'Do projektu',
    wyjasnienie: 'Przenosi sesję do projektu istniejącego albo zakłada nowy o podanej nazwie.',
  });
  dodaj(!zarchiwizowana && wProjekcie, czynnosci.wyjmijZProjektu, {
    klucz: 'projekt-zdejmij',
    napis: 'Wyjmij z projektu',
    wyjasnienie: 'Zdejmuje przypisanie do projektu. Sesja zostaje w historii bieżącej.',
  });

  dodaj(zarchiwizowana, czynnosci.przywroc, {
    klucz: 'przywroc',
    napis: 'Przywróć',
    wyjasnienie: 'Wraca sesję z archiwum do historii bieżącej.',
  });
  dodaj(!zarchiwizowana, czynnosci.archiwizuj, {
    klucz: 'archiwum',
    napis: 'Archiwizuj',
    wyjasnienie: 'Odkłada sesję do archiwum. To nie jest usunięcie — zapis zostaje w całości.',
    postac: 'dn-btn--zarys',
  });

  return wykaz;
}
