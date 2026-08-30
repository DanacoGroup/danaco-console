/**
 * Wiązanie Centrum dowodzenia z rdzeniem. Znacznik niesie biblioteka
 * Właściciela — ten plik nic nie buduje: wypełnia wykaz sesji odpowiedzią
 * rdzenia, zakłada sesje na żądanie i wprowadza w okno modułu.
 */

import { Command, type Environment, type Module, type Session } from '../../../shared/contract.ts';
import type { Kanal, Wynik } from '../protokol/kanal.ts';
import { wywolaj } from '../protokol/wywolanie.ts';
import { zwiazWyborModulu } from './wybor-modulu.ts';
import { oglos } from './ogloszenie.ts';
import { zwiazOkno, zwiazOknoStojace } from './okno-modulu.ts';
import {
  przygotujPasmo,
  przypnijKarte,
  ustawKarty,
  ustawOknaRobocze,
  zaznaczKarte,
  zdejmijKarte,
  zwiazOknaRobocze,
  zwiazPasmo,
} from './karty-okien.ts';
import {
  oknaRobocze,
  oknoBiezace,
  otworzOkno,
  przelaczOkno,
  zamknijOkno,
  zapiszKarte,
  zdejmijKarteOkna,
  zdejmijKartyPoPrawej,
  zostawKarte,
} from './okna-robocze.ts';
import {
  otworzSesje,
  przejmijOgnisko,
  sesjaBiezaca,
  wskazSrodowisko,
  zalozSesje,
} from './sesja-biezaca.ts';
import { zwiazStudio } from './studio.ts';

/** Kod modułu, którego wnętrze wchodzi do wydania; pozostałe moduły stoją w szynie, lecz okna w tym wydaniu nie mają. */
const KOD_MODULU_WYDANIA = 'studio';

/* Drogi powrotu na stronę główną, które niesie znacznik Właściciela: przycisk
   pasa narzędzi, pozycja menu aplikacji i karta główna okna. */
const POWROT_NA_STRONE_GLOWNA =
  '[aria-label="Centrum dowodzenia"], .dn-karta-widoku--glowna, [data-wyjscie-modulu]';

/** Węzły Centrum, na których wiązanie pracuje. Brak któregokolwiek znaczy, że okno Centrum nie stoi. */
interface WezlyCentrum {
  obszar: HTMLElement;
  wykazSesji: HTMLElement;
  /** Płótno okna roboczego; widoki kart stoją w nim obok siebie. */
  plotno: HTMLElement;
  /** Widok karty głównej — Centrum dowodzenia. */
  kartaGlowna: HTMLElement;
  /** Widok karty modułu; wnętrze modułu wchodzi pod jego głowę. */
  kartaModulu: HTMLElement;
}

/** Wiązanie stoi raz na dokument: powłoka może wstawić okno ponownie, a podwójny nasłuch dawałby podwójne sesje. */
let zwiazane = false;

/** Wiąże Centrum z rdzeniem; kanał z obiektu globalnego, bo biblioteka nie jest modułem. Prawda znaczy, że znacznik stał i wiązanie stanęło. */
export function zwiazCentrum(kanal: Kanal | undefined = globalThis.DanacoKanal): boolean {
  if (zwiazane || kanal === undefined) return false;
  const wezly = zbierzWezly();
  if (wezly === null) return false;
  zwiazane = true;

  const wzorWiersza = zdejmijWzorWiersza(wezly.wykazSesji);
  /* Pasmo kart i wykaz okien roboczych biorą wzory z treści przykładowej,
     więc przygotowanie pasma stoi przed jej zdjęciem. */
  przygotujPasmo();
  zdejmijTresciPrzykladowe();
  const odswiez = (): void => {
    void odswiezWykaz(kanal, wezly.wykazSesji, wzorWiersza);
  };
  odswiez();
  void wypelnijKomponenty(kanal);
  void wypelnijProjekty(kanal);
  void wypelnijSrodowiska(kanal, wezly.obszar);
  const katalogModulow = new Map<string, Module>();
  const modulyPoId = new Map<string, Module>();
  void wczytajModuly(kanal, katalogModulow, modulyPoId);
  void przejmijOgnisko(kanal);
  const katalogSrodowisk = new Map<string, Environment>();
  void wczytajSrodowiska(kanal, katalogSrodowisk);

  wezly.obszar.addEventListener('click', (zdarzenie) => {
    const cel = zdarzenie.target;
    if (!(cel instanceof Element)) return;
    if (cel.closest('[data-okno-nowe]') !== null) {
      void zalozKarteSesji(kanal, wezly.wykazSesji, wzorWiersza);
      return;
    }
  });

  /* Czynności sesji słuchają na dokumencie, nie na wnętrzu okna: biblioteka
     menu przenosi treść menu poza wiersz, więc zdarzenie nie przechodzi przez
     panel, w którym wiersz stoi. */
  document.addEventListener('click', (zdarzenie) => {
    const cel = zdarzenie.target;
    if (!(cel instanceof Element)) return;
    const czynnosc = cel.closest<HTMLElement>('[data-poz-akcja]');
    const wykonaj = CZYNNOSCI_SESJI[czynnosc?.dataset.pozAkcja ?? ''];
    const idSesji = czynnosc?.dataset.idSesji;
    if (wykonaj === undefined || idSesji === undefined) return;
    void wykonaj(kanal, idSesji).then((wynik) => {
      // Odmowa rdzenia wychodzi na wierzch: czynność, która milczy po
      // niepowodzeniu, zostawia Operatora przy wykazie sprzed czynności bez
      // słowa, dlaczego się nie zmienił.
      if (!wynik.udany) oglos('Czynność sesji', wynik.blad?.message ?? 'Rdzeń odmówił wykonania.');
      odswiez();
    });
  });

  /* Wnętrze okna zmienia się w miejscu, więc kolejne wejście podmienia element
     wstawiony poprzednio, nie ten zdjęty przy pierwszym. Szyna stoi poza
     wnętrzem, więc nasłuch obejmuje cały dokument.

     Faza przechwytywania jest konieczna: grot wejścia biblioteki zatrzymuje
     zdarzenie na sobie, więc w fazie bąbelkowania nasłuch nigdy by go nie
     zobaczył. */
  /* Nazwa środowiska przechodzi z Centrum przez przedsionek aż do głowy karty
     modułu; kafel przedsionka nie stoi w szynie, więc sam jej nie niesie. */
  let nazwaSrodowiska = '';
  /* Moduł karty modułu: po nim wiadomo, czy kolejne wejście otwiera nową kartę,
     czy wraca do karty stojącej. */
  let kartaBiezaca = '';

  /** Otwiera moduł kartą okna roboczego: głowa, wnętrze, pasmo kart i wiązanie. */
  const otworzModul = (modul: Module, idOkna: string): void => {
    if (!wstawTrescModulu(wezly, 'dn-tresc-' + modul.code)) {
      zapowiedzModul(modul);
      return;
    }
    opiszGloweKarty(wezly, modul.name, nazwaSrodowiska);
    opiszPasekModulu(modul.name);
    pokazWidok(wezly, wezly.kartaModulu);
    kartaBiezaca = modul.code;
    zapiszKarte(modul.code, modul.name);
    odswiezPasmo();
    if (modul.code === KOD_MODULU_WYDANIA) zwiazStudio(kanal, nazwaSrodowiska);
    else if (idOkna === '') zwiazOkno(kanal, modul.code, nazwaSrodowiska);
    else zwiazOknoStojace(kanal, modul.code, nazwaSrodowiska, idOkna);
  };

  /* Porządek drzewa projektów; widok „reakcja" nie ma pola w kontrakcie. */
  document.addEventListener('click', (zdarzenie) => {
    const cel = zdarzenie.target;
    if (!(cel instanceof Element)) return;
    const widok = cel.closest<HTMLElement>('[data-projekty-widok]')?.dataset.projektyWidok;
    const porzadek = cel.closest<HTMLElement>('[data-projekty-sort]')?.dataset.projektySort;
    if (widok === undefined && porzadek === undefined) return;
    zdarzenie.stopPropagation();
    if (widok === 'reakcja') {
      oglos('Drzewo projektów', 'Kontrakt nie niesie dla projektu oczekiwania na reakcję — '
        + 'tego wskazania nie da się dziś spełnić.');
      return;
    }
    if (porzadek !== undefined) porzadekProjektow = porzadek;
    void wypelnijProjekty(kanal);
  }, true);

  /* Wydanie zapisu sesji: czynność panelu bocznego oddaje Operatorowi plik
     z sesją bieżącą. Bez wskazanej sesji nie ma czego wydać i mówi to wprost. */
  document.addEventListener('click', (zdarzenie) => {
    const cel = zdarzenie.target;
    if (!(cel instanceof Element) || cel.closest('[data-panel-eksport]') === null) return;
    zdarzenie.stopPropagation();
    const idSesji = sesjaBiezaca();
    if (idSesji === '') {
      oglos('Wydanie sesji', 'Wskaż sesję w wykazie — wydanie obejmuje sesję bieżącą.');
      return;
    }
    void wywolaj(kanal, Command.SessionExport, { sessionId: idSesji }).then((wynik) => {
      if (!wynik.udany || wynik.wynik === undefined) {
        oglos('Wydanie sesji', wynik.blad?.message ?? 'Rdzeń odmówił wydania zapisu.');
        return;
      }
      oddajPlik(wynik.wynik.fileName, wynik.wynik.content);
    });
  }, true);

  /* Wskazówka startowa Centrum: pozycja menu widoku zdejmuje ją i przywraca. */
  document.addEventListener('click', (zdarzenie) => {
    const cel = zdarzenie.target;
    if (!(cel instanceof Element)) return;
    const pozycja = cel.closest<HTMLElement>('[data-cd-wskazowka]');
    if (pozycja === null) return;
    zdarzenie.stopPropagation();
    const wskazowka = document.querySelector<HTMLElement>('.cd-wskazowka, #cd-start');
    if (wskazowka === null) return;
    wskazowka.hidden = !wskazowka.hidden;
    pozycja.setAttribute('aria-checked', String(!wskazowka.hidden));
  }, true);

  /* Zamknięcie wskazówki startowej znakiem przy niej samej. */
  document.addEventListener('click', (zdarzenie) => {
    const cel = zdarzenie.target;
    if (!(cel instanceof Element) || cel.closest('#cd-start-zamknij') === null) return;
    zdarzenie.stopPropagation();
    const wskazowka = document.querySelector<HTMLElement>('.cd-wskazowka, #cd-start');
    if (wskazowka !== null) wskazowka.hidden = true;
  }, true);

  /* Przypięcie panelu bocznego trzyma go rozwiniętym przy zmianie karty. */
  document.addEventListener('click', (zdarzenie) => {
    const cel = zdarzenie.target;
    if (!(cel instanceof Element)) return;
    const przycisk = cel.closest<HTMLElement>('[data-panel-przypnij]');
    if (przycisk === null) return;
    zdarzenie.stopPropagation();
    const przypiety = przycisk.getAttribute('aria-pressed') === 'true';
    przycisk.setAttribute('aria-pressed', String(!przypiety));
  }, true);

  /* Widok i porządek wykazu sesji: panel boczny jest zarządem sesji, więc oba
     wskazania działają na wykazie odpowiedzi rdzenia. Porządek po środowisku
     i widok reakcji nie mają pola w kontrakcie i mówią o tym wprost. */
  document.addEventListener('click', (zdarzenie) => {
    const cel = zdarzenie.target;
    if (!(cel instanceof Element)) return;
    const widok = cel.closest<HTMLElement>('[data-sesje-widok]')?.dataset.sesjeWidok;
    const porzadek = cel.closest<HTMLElement>('[data-sesje-sort]')?.dataset.sesjeSort;
    if (widok === undefined && porzadek === undefined) return;
    zdarzenie.stopPropagation();

    if (widok !== undefined) widokWykazu = widok;
    if (porzadek !== undefined) porzadekWykazu = porzadek;
    odswiez();
  }, true);

  /* Odświeżenie Centrum czyta rejestr rdzenia na nowo: wykaz sesji i karty
     środowisk. Znacznik niesie ten przycisk, nikt go nie wiązał. */
  document.addEventListener('click', (zdarzenie) => {
    const cel = zdarzenie.target;
    if (!(cel instanceof Element) || cel.closest('[data-cd-odswiez]') === null) return;
    zdarzenie.stopPropagation();
    odswiez();
    void wypelnijSrodowiska(kanal, wezly.obszar);
  }, true);

  /* Prawa strona ramy — okno boczne — stoi w stanie domyślnym z samouczkiem.
     Znacznik niesie jej przełącznik i zwinięcie; nikt ich nie wiązał, więc
     strona nie otwierała się wcale. */
  document.addEventListener('click', (zdarzenie) => {
    const cel = zdarzenie.target;
    if (!(cel instanceof Element)) return;
    const panel = document.getElementById('panel-samouczek');
    if (panel === null) return;
    if (cel.closest('[data-przelacz-samouczek]') !== null) {
      zdarzenie.stopPropagation();
      panel.hidden = !panel.hidden;
      return;
    }
    if (cel.closest('[data-zwin-samouczek]') !== null) {
      zdarzenie.stopPropagation();
      panel.hidden = true;
    }
  }, true);

  /* „Zamknij kartę" z menu okna zdejmuje kartę bieżącą i wraca na Centrum —
     tak samo jak znak zamknięcia na samej karcie. */
  /* Czynności na karcie z menu okna roboczego; rozpoznaje je podpis pozycji,
     bo znacznik nie niesie dla nich własnego uchwytu. */
  document.addEventListener('click', (zdarzenie) => {
    const cel = zdarzenie.target;
    if (!(cel instanceof Element)) return;
    /* Czynności czytamy wyłącznie z menu okna roboczego: ten sam podpis stoi
       w innych menu powłoki i tam znaczy co innego. */
    const pozycja = cel.closest('#menu-karty-otwarte .sta-menu-poz');
    if (pozycja === null || kartaBiezaca === '') return;
    const czynnosc = (pozycja.textContent ?? '').trim();
    if (czynnosc === 'Zamknij kartę') {
      zdarzenie.stopPropagation();
      zdejmijKarte(kartaBiezaca);
      zdejmijKarteOkna(kartaBiezaca);
      wrocDoCentrum();
      return;
    }
    if (czynnosc === 'Zamknij pozostałe') {
      zdarzenie.stopPropagation();
      for (const kod of zostawKarte(kartaBiezaca)) zdejmijKarte(kod);
      odswiezPasmo();
      return;
    }
    if (czynnosc === 'Zamknij karty po prawej') {
      zdarzenie.stopPropagation();
      for (const kod of zdejmijKartyPoPrawej(kartaBiezaca)) zdejmijKarte(kod);
      odswiezPasmo();
      return;
    }
    if (czynnosc === 'Przypnij kartę') {
      zdarzenie.stopPropagation();
      przypnijKarte(kartaBiezaca);
    }
  }, true);

  /** Przerysowuje pasmo kart z wykazu okna bieżącego. */
  function odswiezPasmo(): void {
    const okno = oknoBiezace();
    ustawKarty(
      okno.karty.map((kod) => {
        const nazwa = katalogModulow.get(kod)?.name ?? kod;
        return { id: kod, nazwa, modul: nazwa };
      }),
      okno.kartaBiezaca,
    );
    odswiezOknaRobocze();
  }

  /** Wraca do karty głównej okna roboczego — Centrum dowodzenia. */
  const wrocDoCentrum = (): void => {
    pokazWidok(wezly, wezly.kartaGlowna);
    kartaBiezaca = '';
    opiszPasekModulu('');
    zapiszKarte('', '');
    zaznaczKarte('');
    odswiezOknaRobocze();
    odswiez();
  };

  /** Odświeża wykaz okien roboczych w menu okna. */
  function odswiezOknaRobocze(): void {
    ustawOknaRobocze(
      oknaRobocze().map((okno) => ({ id: okno.id, nazwa: okno.nazwa })),
      oknoBiezace().id,
    );
  }

  /** Stawia okno robocze na jego karcie bieżącej: Centrum albo moduł. */
  const pokazOknoRobocze = (): void => {
    const okno = oknoBiezace();
    const modul = okno.kartaBiezaca === '' ? undefined : katalogModulow.get(okno.kartaBiezaca);
    if (modul === undefined) {
      pokazWidok(wezly, wezly.kartaGlowna);
      kartaBiezaca = '';
      odswiezPasmo();
      odswiez();
    } else {
      otworzModul(modul, '');
    }
    odswiezOknaRobocze();
  };

  odswiezOknaRobocze();
  zwiazOknaRobocze(
    (id) => {
      if (przelaczOkno(id) === undefined) return;
      pokazOknoRobocze();
    },
    () => {
      otworzOkno();
      pokazOknoRobocze();
    },
    () => {
      zamknijOkno();
      pokazOknoRobocze();
    },
  );
  zwiazPasmo(
    (kodKarty) => {
      const modul = katalogModulow.get(kodKarty);
      if (modul === undefined) return;
      if (kartaBiezaca === kodKarty) {
        pokazWidok(wezly, wezly.kartaModulu);
        zaznaczKarte(kodKarty);
        return;
      }
      otworzModul(modul, '');
    },
    (kodKarty) => {
      if (kartaBiezaca === kodKarty) wrocDoCentrum();
      zdejmijKarte(kodKarty);
    },
  );

  document.addEventListener('click', (zdarzenie) => {
    const cel = zdarzenie.target;
    if (!(cel instanceof Element)) return;

    // Powrót na kartę główną okna roboczego.
    if (cel.closest(POWROT_NA_STRONE_GLOWNA) !== null) {
      wrocDoCentrum();
      return;
    }

    /* Wykaz sesji stoi w panelu bocznym okna roboczego; osobnego okna rejestru
       sesji to wydanie nie niesie, więc czynność nazywa to wprost. */
    if (cel.closest('[data-otwarz-historie]') !== null) {
      zdarzenie.stopPropagation();
      oglos('Historia sesji', 'Sesje konta stoją w panelu bocznym okna roboczego. '
        + 'Osobne okno rejestru sesji nie wchodzi do tego wydania.');
      return;
    }

    /* Wiersz wykazu sesji jest drogą powrotu do pracy: karta sesji otwiera się
       wraz z oknami, a Operator wraca do okna, w którym był. Wiersz nie
       przechwytuje kliknięć swojego menu — tam stoją czynności karty. */
    const wiersz = cel.closest<HTMLElement>('[data-id-sesji]');
    if (wiersz !== null && cel.closest('[data-menu]') === null
      && cel.closest('[data-poz-akcja]') === null) {
      zdarzenie.stopPropagation();
      void otworzSesje(kanal, wiersz.dataset.idSesji ?? '').then((okna) => {
        const okno = okna.find((kandydat) => kandydat.status !== 'closed') ?? okna[0];
        const modul = okno === undefined
          ? undefined
          : modulyPoId.get(okno.moduleId) ?? katalogModulow.get(okno.moduleId);
        if (okno === undefined || modul === undefined) {
          oglos('Karta sesji', 'Karta jest otwarta i przyjmie okno modułu; '
            + 'okna w niej jeszcze nie ma.');
          return;
        }
        otworzModul(modul, okno.id);
      });
      return;
    }

    const kodNowejSesji = cel.closest<HTMLElement>('[data-nowa-sesja-srodowisko]')?.dataset
      .nowaSesjaSrodowisko;
    if (kodNowejSesji !== undefined) {
      void zalozKarteSesji(kanal, wezly.wykazSesji, wzorWiersza);
      const widok = widokPrzedsionka(wezly, gniazdoPrzedsionka(kodNowejSesji));
      if (widok === null) return;
      pokazWidok(wezly, widok);
      nazwaSrodowiska = katalogSrodowisk.get(kodNowejSesji)?.name ?? '';
      wskazSrodowisko(kodNowejSesji);
      void zwiazWyborModulu(kanal, kodNowejSesji);
      return;
    }

    /* Środowisko bez modułów w rejestrze nie ma czego rozwinąć. Bez tego zdania
       pozycja szyny odsyłałaby do listy, która nigdy nie stanie. */
    const przelacznik = cel.closest<HTMLElement>('.dn-szyna-poz--srodowisko');
    if (przelacznik !== null) {
      const srodowisko = katalogSrodowisk.get(przelacznik.dataset.srodowisko ?? '');
      if (srodowisko !== undefined && (srodowisko.moduleCodes?.length ?? 0) === 0) {
        zdarzenie.stopPropagation();
        oglos(srodowisko.name, 'Rejestr rdzenia nie wskazuje dla tego środowiska '
          + 'ani jednego modułu, więc lista nie ma czego rozwinąć.');
      }
    }

    const kodSrodowiska = kodSrodowiskaWejscia(cel);
    if (kodSrodowiska !== '') {
      const widok = widokPrzedsionka(wezly, gniazdoPrzedsionka(kodSrodowiska));
      if (widok === null) return;
      pokazWidok(wezly, widok);
      nazwaSrodowiska = nazwaKartySrodowiska(cel) || nazwaSrodowiska;
      wskazSrodowisko(kodSrodowiska);
      void zwiazWyborModulu(kanal, kodSrodowiska);
      return;
    }

    /* Komponent własny otwiera się kartą modułu swojego rodzaju: komponent
       automatyki wchodzi w Automations, ekspert w Agents. */
    const komponent = cel.closest<HTMLElement>('[data-otworz-komponent]')?.dataset.otworzKomponent;
    if (komponent !== undefined) {
      const modulKomponentu = katalogModulow.get(RODZAJE_KOMPONENTOW.get(komponent) ?? '');
      if (modulKomponentu === undefined) return;
      zdarzenie.stopPropagation();
      otworzModul(modulKomponentu, '');
      return;
    }

    /* Znak „+" pasma otwiera kartę modułu wskazanego pozycją menu; nazwa
       pozycji jest nazwą modułu z rejestru rdzenia. */
    const nowaKarta = cel.closest<HTMLElement>('[data-nowa-karta-modul]')?.dataset.nowaKartaModul;
    if (nowaKarta !== undefined) {
      const wskazany = [...katalogModulow.values()].find((modul) => modul.name === nowaKarta);
      if (wskazany === undefined) return;
      zdarzenie.stopPropagation();
      otworzModul(wskazany, '');
      return;
    }

    const kod = kodModulu(cel);
    if (kod === '') return;
    const modul = katalogModulow.get(kod);
    if (modul === undefined) return;
    zdarzenie.stopPropagation();
    nazwaSrodowiska = nazwaSrodowiskaWejscia(cel) || nazwaSrodowiska;
    /* Pionowa szyna jest paskiem szybkiego dostępu: jedno naciśnięcie otwiera
       OKNO ROBOCZE z wybranym modułem, z pominięciem drogi przez przedsionek.
       Kafel Centrum i kafel przedsionka otwierają kartę w oknie stojącym. */
    if (cel.closest('.dn-szyna-poz--modul') !== null) otworzOkno();
    otworzModul(modul, '');
  }, true);

  return true;
}

/**
 * Szablon przedsionka wskazanego środowiska. Każde środowisko ma własny
 * prototyp, bo kafle niosą znaki swoich modułów; środowisko bez prototypu
 * wchodzi na przedsionek TalkIn, żeby droga wejścia nie urwała się wcale.
 */
function gniazdoPrzedsionka(kodSrodowiska: string): string {
  return 'dn-tresc-przedsionek-' + kodSrodowiska.toLowerCase();
}

/** Kod modułu wskazanego pozycją szyny, kaflem Centrum albo kaflem przedsionka; zapis sprowadza do małych liter, bo rejestr rdzenia trzyma kody małymi. */
function kodModulu(cel: Element): string {
  const wskazanie =
    cel.closest('.dn-szyna-poz--modul')?.getAttribute('data-modul') ??
    cel.closest('.pd-kafel')?.getAttribute('data-modul') ??
    cel.closest('.cd-kafel, .dn-kafel--modul')?.getAttribute('data-komponent') ??
    '';
  return wskazanie.toLowerCase();
}

/** Wczytuje rejestr środowisk do spisu po kodzie; nazwa i wykaz modułów rozstrzygają, co pozycja szyny może otworzyć. */
async function wczytajSrodowiska(kanal: Kanal, spis: Map<string, Environment>): Promise<void> {
  const wynik = await wywolaj(kanal, Command.EnvironmentList, { includeModules: true });
  if (!wynik.udany || wynik.wynik === undefined) return;
  for (const srodowisko of wynik.wynik.environments) spis.set(srodowisko.code, srodowisko);
}

/** Wczytuje rejestr modułów do spisu po kodzie; opisy modułów są jedynym źródłem zapowiedzi okna, którego wydanie jeszcze nie niesie. */
async function wczytajModuly(
  kanal: Kanal,
  spis: Map<string, Module>,
  poId: Map<string, Module>,
): Promise<void> {
  const wynik = await wywolaj(kanal, Command.ModuleList, {});
  if (!wynik.udany || wynik.wynik === undefined) return;
  for (const modul of wynik.wynik.modules) {
    spis.set(modul.code, modul);
    poId.set(modul.id, modul);
  }
}

/**
 * Zapowiada moduł, którego okna to wydanie nie niesie. Komunikat nazywa
 * niegotowość wprost i podaje opis modułu z rejestru rdzenia; twierdzenie, że
 * okno się otwiera, byłoby nieprawdą, a milczenie zostawiłoby pozycję martwą.
 */
function zapowiedzModul(modul: Module | undefined): void {
  if (modul === undefined) return;
  const opis = modul.description ?? '';
  const zdanie = opis === '' ? '' : opis + ' ';
  oglos(modul.name, zdanie + 'Okno tego modułu nie wchodzi do tego wydania.');
}

/** Powiadomienie biblioteki; jej brak zostawia czynność bez komunikatu, bo dorabianie własnego byłoby stawianiem elementu. */
/** Nazwa środowiska z karty Centrum, wpisana tam wcześniej rejestrem rdzenia. */
function nazwaKartySrodowiska(cel: Element): string {
  const karta = cel.closest('.dn-karta-srodowiska');
  return karta?.querySelector('.dn-karta-srodowiska-tytul')?.textContent?.trim() ?? '';
}

/** Kod środowiska, do którego prowadzi naciśnięty grot karty Centrum; pustka znaczy, że kliknięcie karty nie dotyczyło. */
function kodSrodowiskaWejscia(cel: Element): string {
  const wejscie = cel.closest('[data-wejdz]');
  if (wejscie === null) return '';
  return wejscie.closest<HTMLElement>('.dn-karta-srodowiska')?.dataset.srodowisko ?? '';
}

/** Zbiera węzły Centrum; brak wykazu sesji znaczy, że w ramie stoi inne okno. */
function zbierzWezly(): WezlyCentrum | null {
  const obszar = document.querySelector<HTMLElement>('.dn-rama-prawa .dn-obszar');
  const wykazSesji = document.getElementById('wykaz-sesji');
  const plotno = document.querySelector<HTMLElement>('.cd-plotno');
  const kartaGlowna = document.getElementById('cd-tresc');
  const kartaModulu = document.getElementById('karta-modul');
  if (obszar === null || wykazSesji === null || plotno === null
    || kartaGlowna === null || kartaModulu === null) return null;
  return { obszar, wykazSesji, plotno, kartaGlowna, kartaModulu };
}

/**
 * Pokazuje jeden widok okna roboczego. Widoki kart stoją obok siebie w płótnie
 * i różnią się zasłoną — okno robocze zostaje na miejscu wraz z pasmem kart,
 * panelem bocznym i pasem stanu.
 */
function pokazWidok(wezly: WezlyCentrum, widok: HTMLElement): void {
  for (const kandydat of wezly.plotno.querySelectorAll<HTMLElement>('.cd-tresc')) {
    kandydat.hidden = kandydat !== widok;
  }
  widok.hidden = false;
}

/**
 * Widok przedsionka środowiska. Stoi w płótnie obok karty głównej, wzorem
 * karty modułu; wchodzi raz i wraca przy każdym kolejnym wejściu w środowisko.
 */
function widokPrzedsionka(wezly: WezlyCentrum, gniazdo: string): HTMLElement | null {
  const szablon = document.getElementById(gniazdo);
  if (!(szablon instanceof HTMLTemplateElement)) return null;
  const blok = szablon.content.firstElementChild;
  if (blok === null) return null;
  let widok = wezly.plotno.querySelector<HTMLElement>('.cd-tresc--przedsionek');
  if (widok === null) {
    widok = document.createElement('div');
    widok.className = 'cd-tresc cd-tresc--przedsionek';
    wezly.plotno.appendChild(widok);
  }
  widok.replaceChildren(blok.cloneNode(true));
  return widok;
}

/**
 * Wnętrze modułu wchodzi pod głowę karty modułu. Odnośnik do prototypu jest
 * rusztowaniem prototypu i nie wchodzi do produktu.
 */
function wstawTrescModulu(wezly: WezlyCentrum, gniazdo: string): boolean {
  const szablon = document.getElementById(gniazdo);
  if (!(szablon instanceof HTMLTemplateElement)) return false;
  const blok = szablon.content.firstElementChild;
  if (blok === null) return false;
  wezly.kartaModulu.querySelector('.cd-modul-odnosnik')?.remove();
  for (const stojace of [...wezly.kartaModulu.children]) {
    if (!stojace.classList.contains('cd-modul-glowa')) stojace.remove();
  }
  wezly.kartaModulu.appendChild(blok.cloneNode(true));
  return true;
}

/** Oddaje Operatorowi plik z treścią wydaną przez rdzeń. */
function oddajPlik(nazwa: string, tresc: string): void {
  const adres = URL.createObjectURL(new Blob([tresc], { type: 'text/plain;charset=utf-8' }));
  const odnosnik = document.createElement('a');
  odnosnik.href = adres;
  odnosnik.download = nazwa;
  odnosnik.click();
  URL.revokeObjectURL(adres);
}

/** Wpisuje nazwę modułu w pas stanu ramy; karta główna zostawia pole puste. */
function opiszPasekModulu(nazwa: string): void {
  const pole = document.querySelector('[data-pasek-modul-nazwa]');
  if (pole !== null) pole.textContent = nazwa;
}

/** Opisuje głowę karty modułu nazwą modułu i nazwą karty sesji. */
function opiszGloweKarty(wezly: WezlyCentrum, nazwaModulu: string, nazwaSesji: string): void {
  const nazwa = wezly.kartaModulu.querySelector('[data-karta-modul-nazwa]');
  if (nazwa !== null) nazwa.textContent = nazwaModulu;
  const meta = wezly.kartaModulu.querySelector('.cd-modul-glowa .dn-meta');
  if (meta !== null) meta.textContent = nazwaSesji === '' ? '' : 'sesja: ' + nazwaSesji;
}

/**
 * Wypełnia wykaz komponentów własnych Centrum rejestrem rdzenia. Wzór pozycji
 * zdejmuje się z treści przykładowej, zanim wykaz zostanie wyczyszczony.
 */
const RODZAJE_KOMPONENTOW = new Map<string, string>();

async function wypelnijKomponenty(kanal: Kanal): Promise<void> {
  const wykaz = document.getElementById('cd-wlasne');
  if (wykaz === null) return;
  const wzor = wykaz.querySelector<HTMLElement>('.cd-wlasny');
  const wzorPozycji = wzor === null ? null : (wzor.cloneNode(true) as HTMLElement);
  const wynik = await wywolaj(kanal, Command.ComponentList, {});
  if (!wynik.udany || wynik.wynik === undefined || wzorPozycji === null) return;
  wykaz.replaceChildren();
  for (const komponent of wynik.wynik.components) {
    const pozycja = wzorPozycji.cloneNode(true) as HTMLElement;
    const przycisk = pozycja.querySelector<HTMLElement>('[data-otworz-komponent]');
    if (przycisk !== null) przycisk.dataset.otworzKomponent = komponent.id;
    // Rodzaj komponentu jest kodem modułu, w którym komponent stoi.
    RODZAJE_KOMPONENTOW.set(komponent.id, komponent.kind);
    const nazwa = pozycja.querySelector('.dn-kafel-nazwa');
    if (nazwa !== null) nazwa.textContent = komponent.name;
    const opis = pozycja.querySelector('.dn-kafel-opis');
    if (opis !== null) opis.textContent = komponent.description ?? '';
    wykaz.appendChild(pozycja);
  }
}

/**
 * Wypełnia drzewo projektów lewego panelu wykazem rdzenia wraz z sesjami
 * każdego projektu. Projekt bez ani jednej sesji też stoi — lewy panel jest
 * jedynym miejscem, w którym widać całość dorobku Operatora.
 */
async function wypelnijProjekty(kanal: Kanal): Promise<void> {
  const drzewo = document.querySelector<HTMLElement>('#panel-projekty .dn-panel-drzewo');
  if (drzewo === null) return;
  const galaz = drzewo.querySelector<HTMLElement>('.dn-panel-galaz');
  const wzorGalezi = galaz === null ? null : (galaz.cloneNode(true) as HTMLElement);
  const wzorSesji = wzorGalezi?.querySelector<HTMLElement>('.dn-panel-wiersz') ?? null;
  const projekty = await wywolaj(kanal, Command.ProjectList, {});
  const sesje = await wywolaj(kanal, Command.SessionList, {});
  if (!projekty.udany || projekty.wynik === undefined || wzorGalezi === null) return;
  const wedlugProjektu = new Map<string, Session[]>();
  for (const sesja of sesje.wynik?.sessions ?? []) {
    if (sesja.projectId === undefined || sesja.projectId === '') continue;
    const zebrane = wedlugProjektu.get(sesja.projectId) ?? [];
    zebrane.push(sesja);
    wedlugProjektu.set(sesja.projectId, zebrane);
  }
  const uporzadkowane = [...projekty.wynik.projects].sort((a, b) => (
    porzadekProjektow === 'nazwa' ? a.name.localeCompare(b.name, 'pl') : b.updatedAt - a.updatedAt
  ));
  drzewo.replaceChildren();
  for (const projekt of uporzadkowane) {
    const wpis = wzorGalezi.cloneNode(true) as HTMLElement;
    wpis.dataset.projekt = projekt.id;
    const nazwa = wpis.querySelector('.dn-panel-galaz-nazwa');
    if (nazwa !== null) nazwa.textContent = projekt.name;
    for (const pozycja of wpis.querySelectorAll('.dn-panel-wiersz')) pozycja.remove();
    for (const sesja of wedlugProjektu.get(projekt.id) ?? []) {
      if (wzorSesji === null) break;
      const wiersz = wzorSesji.cloneNode(true) as HTMLElement;
      wiersz.dataset.idSesji = sesja.id;
      const podpis = wiersz.querySelector('.dn-obszar-pozycja-nazwa') ?? wiersz;
      podpis.textContent = sesja.title ?? 'Sesja bez nazwy';
      wpis.appendChild(wiersz);
    }
    drzewo.appendChild(wpis);
  }
}

/** Zdejmuje wzór wiersza z treści przykładowej; kształt wiersza bierze się ze znacznika, nie z kodu. */
function zdejmijWzorWiersza(wykaz: HTMLElement): HTMLElement | null {
  const wiersz = wykaz.querySelector<HTMLElement>('.dn-panel-wiersz');
  return wiersz === null ? null : (wiersz.cloneNode(true) as HTMLElement);
}

/** Wczytuje wykaz sesji rdzenia i wstawia go w miejsce treści przykładowej. */
async function odswiezWykaz(
  kanal: Kanal,
  wykaz: HTMLElement,
  wzor: HTMLElement | null,
): Promise<void> {
  if (wzor === null) return;
  const wynik = await wywolaj(kanal, Command.SessionList, {});
  if (!wynik.udany || wynik.wynik === undefined) return;
  wykaz.replaceChildren();
  for (const sesja of uporzadkuj(przesiej(wynik.wynik.sessions))) {
    wykaz.appendChild(zbudujWiersz(wzor, sesja));
  }
}

/** Widok wykazu sesji wskazany w panelu bocznym; „wszystkie" nie zawęża niczego. */
let widokWykazu = 'wszystkie';
/** Porządek wykazu sesji wskazany w panelu bocznym. */
let porzadekWykazu = 'czynnosc';
/** Porządek drzewa projektów wskazany w panelu bocznym. */
let porzadekProjektow = 'czynnosc';

/** Zawęża wykaz sesji do stanu wskazanego widokiem panelu. */
function przesiej(sesje: Session[]): Session[] {
  if (widokWykazu === 'czynne') return sesje.filter((sesja) => sesja.status === 'active');
  if (widokWykazu === 'zakonczone') return sesje.filter((sesja) => sesja.status !== 'active');
  if (widokWykazu === 'reakcja') return sesje.filter((sesja) => sesja.awaitingReaction === true);
  return sesje;
}

/** Porządkuje wykaz sesji wskazaniem panelu; porządek nieznany zostawia kolejność rdzenia. */
function uporzadkuj(sesje: Session[]): Session[] {
  const wykaz = [...sesje];
  if (porzadekWykazu === 'nazwa') {
    return wykaz.sort((a, b) => (a.title ?? '').localeCompare(b.title ?? '', 'pl'));
  }
  if (porzadekWykazu === 'nazwa-odwrotnie') {
    return wykaz.sort((a, b) => (b.title ?? '').localeCompare(a.title ?? '', 'pl'));
  }
  if (porzadekWykazu === 'najstarsze') return wykaz.sort((a, b) => a.createdAt - b.createdAt);
  if (porzadekWykazu === 'czynnosc') return wykaz.sort((a, b) => b.updatedAt - a.updatedAt);
  if (porzadekWykazu === 'srodowisko') {
    return wykaz.sort((a, b) => (a.environmentCode ?? '').localeCompare(b.environmentCode ?? '', 'pl'));
  }
  return wykaz;
}

/**
 * Zwraca klon wzoru wiersza opisany nazwą sesji. Menu czynności zostaje, bo
 * niesie czynności o pokryciu w kontrakcie; jego odwołanie dostaje
 * identyfikator sesji, żeby dwa wiersze nie wskazywały tego samego menu.
 */
function zbudujWiersz(wzor: HTMLElement, sesja: Session): HTMLElement {
  const wiersz = wzor.cloneNode(true) as HTMLElement;
  wiersz.dataset.idSesji = sesja.id;
  const nazwa = wiersz.querySelector('.dn-obszar-pozycja-nazwa');
  // Sesja bez nazwy dostaje nazwany stan pusty: identyfikator jest oznaczeniem magazynu, nie nazwą pracy Operatora.
  if (nazwa !== null) nazwa.textContent = sesja.title ?? 'Sesja bez nazwy';
  wiersz.querySelector('.dn-obszar-pozycja')?.setAttribute('data-id-sesji', sesja.id);
  opiszStanWiersza(wiersz, sesja);
  const menu = wiersz.querySelector('[data-menu-tresc]');
  const wyzwalacz = wiersz.querySelector('[data-menu]');
  if (menu !== null && wyzwalacz !== null) {
    const oznaczenie = 'menu-sesji-' + sesja.id;
    menu.id = oznaczenie;
    wyzwalacz.setAttribute('data-menu', oznaczenie);
  }
  zdejmijCzynnosciBezZrodla(wiersz);
  /* Identyfikator sesji siada na samej pozycji menu, nie tylko na wierszu:
     biblioteka menu przenosi treść menu poza wiersz, więc szukanie sesji
     w przodkach pozycji nic by nie znalazło. */
  for (const pozycja of wiersz.querySelectorAll<HTMLElement>('[data-poz-akcja]')) {
    pozycja.dataset.idSesji = sesja.id;
  }
  return wiersz;
}

/**
 * Nadaje wierszowi stan, którym znacznik panelu barwi sesję: praca, oczekiwanie
 * na reakcję Operatora albo sesja zakończona. Stany są słownikiem znacznika
 * Właściciela, a rozstrzyga o nich odpowiedź rdzenia.
 */
function opiszStanWiersza(wiersz: HTMLElement, sesja: Session): void {
  const stan = sesja.status !== 'active'
    ? 'zakonczone'
    : (sesja.awaitingReaction === true ? 'reakcja' : 'praca');
  for (const wezel of wiersz.querySelectorAll('[data-stan]')) {
    wezel.setAttribute('data-stan', stan);
  }
}

/* Czynności menu, dla których kontrakt ma komendę. Pozycje spoza tego spisu
   znikają: pozycja menu, która nic nie robi, jest obietnicą bez pokrycia.

   Spis trzyma wywołania, nie same nazwy komend: każda z tych komend bierze
   wykaz sesji, a nie pojedyncze wskazanie, i tylko wywołanie zapisane przy
   swojej komendzie daje się sprawdzić kontraktem przy budowaniu. */
const CZYNNOSCI_SESJI: Record<string, (kanal: Kanal, idSesji: string) => Promise<Wynik<unknown>>> = {
  archiwizuj: (kanal, idSesji) =>
    wywolaj(kanal, Command.SessionArchive, { sessionIds: [idSesji] }),
  // Potwierdzenie nieodwracalności niesie sama pozycja menu: nazywa usunięcie
  // trwałym, a rejestr sesji drugiego pytania nie stawia.
  usun: (kanal, idSesji) =>
    wywolaj(kanal, Command.SessionDelete, { sessionIds: [idSesji], confirm: true }),
  wyjmij: (kanal, idSesji) =>
    wywolaj(kanal, Command.SessionProjectClear, { sessionIds: [idSesji] }),
};

/** Zdejmuje z menu wiersza pozycje bez komendy w kontrakcie wraz z rozdzielnikami, które po nich zostały. */
function zdejmijCzynnosciBezZrodla(wiersz: HTMLElement): void {
  for (const pozycja of wiersz.querySelectorAll('[data-menu-tresc] .sta-menu-poz')) {
    const czynnosc = pozycja.getAttribute('data-poz-akcja') ?? '';
    if (CZYNNOSCI_SESJI[czynnosc] === undefined) pozycja.remove();
  }
  for (const rozdzielnik of wiersz.querySelectorAll('[data-menu-tresc] .sta-menu-sep')) {
    const przed = rozdzielnik.previousElementSibling;
    if (przed === null || przed.classList.contains('sta-menu-sep')) rozdzielnik.remove();
  }
}

/** Zakłada kartę sesji, czyni ją bieżącą i odświeża wykaz odpowiedzią rdzenia. */
async function zalozKarteSesji(
  kanal: Kanal,
  wykaz: HTMLElement,
  wzor: HTMLElement | null,
): Promise<void> {
  if (await zalozSesje(kanal, '') === '') return;
  await odswiezWykaz(kanal, wykaz, wzor);
}

/** Nazwa środowiska, przez które Operator wszedł w moduł; pozycja szyny stoi w grupie środowiska, kafel Centrum nie należy do żadnej. */
function nazwaSrodowiskaWejscia(cel: Element): string {
  const grupa = cel.closest('.dn-szyna-poz--modul')?.closest('.dn-szyna-moduly');
  if (grupa === null || grupa === undefined) return '';
  const przelacznik = document.querySelector(
    `.dn-szyna-poz--srodowisko[aria-controls="${grupa.id}"]`,
  );
  return przelacznik?.getAttribute('aria-label') ?? '';
}


/** Zdejmuje treść przykładową bez pokrycia w rdzeniu: karty okien poza główną i komponenty własne. Pusty wykaz odsłania stan pusty ze znacznika. */
function zdejmijTresciPrzykladowe(): void {
  const karty = document.querySelector('.dn-karty-lista');
  if (karty !== null) {
    for (const karta of [...karty.querySelectorAll('.dn-karta-widoku')].slice(1)) karta.remove();
  }
  // Odsyłacz do pliku prototypu prowadzi poza produkt i w wydaniu nie stoi.
  document.querySelector('.cd-modul-odnosnik')?.remove();
}

/** Wypełnia karty środowisk rejestrem rdzenia; karta bez pokrycia znika, bo znacznik niesie ich cztery, a rejestr rozstrzyga, ile ich jest. */
async function wypelnijSrodowiska(kanal: Kanal, obszar: HTMLElement): Promise<void> {
  const wynik = await wywolaj(kanal, Command.EnvironmentList, { includeModules: true });
  if (!wynik.udany || wynik.wynik === undefined) return;
  const rejestr = new Map(wynik.wynik.environments.map((s) => [s.code, s]));
  for (const karta of obszar.querySelectorAll<HTMLElement>('.dn-karta-srodowiska')) {
    const srodowisko = rejestr.get(karta.dataset.srodowisko ?? '');
    if (srodowisko === undefined) {
      karta.remove();
      continue;
    }
    opiszSrodowisko(karta, srodowisko);
  }
}

/** Wpisuje w kartę nazwę, opis, liczbę modułów i miarę sesji czynnych środowiska. */
function opiszSrodowisko(karta: HTMLElement, srodowisko: Environment): void {
  const tytul = karta.querySelector('.dn-karta-srodowiska-tytul');
  if (tytul !== null) tytul.textContent = srodowisko.name;
  const opis = karta.querySelector('.dn-karta-srodowiska-opis');
  if (opis !== null && srodowisko.description !== undefined) {
    opis.textContent = srodowisko.description;
  }
  const miara = karta.querySelector('.dn-karta-srodowiska-motto .cd-metryka-czlon');
  if (miara !== null) miara.textContent = miaraModulow(srodowisko.moduleCodes?.length ?? 0);
  opiszStopke(karta, srodowisko.sessionCount ?? 0);
}

/* Stopka karty niesie dwie rzeczy naraz: liczbę sesji środowiska i grot wejścia.
   Kontrakt nie wiąże sesji ze środowiskiem, więc liczba znika, a grot zostaje —
   zdjęcie całej stopki zabrałoby Operatorowi drogę do przedsionka. */
function opiszStopke(karta: HTMLElement, sesji: number): void {
  const stopka = karta.querySelector('.cd-karta-meta');
  if (stopka === null) return;
  /* Miara stoi w stopce węzłem tekstowym obok kropki stanu i grotu wejścia;
     podpis dla czytnika ekranu zostaje, bo należy do znacznika Właściciela. */
  for (const wezel of stopka.childNodes) {
    if (wezel.nodeType !== Node.TEXT_NODE) continue;
    if ((wezel.nodeValue ?? '').trim() === '') continue;
    wezel.nodeValue = miaraSesji(sesji);
    return;
  }
}

/** Liczba sesji czynnych środowiska wraz z odmianą rzeczownika. */
function miaraSesji(ile: number): string {
  if (ile === 1) return '1 sesja';
  const reszta = ile % 10;
  const dziesiatki = ile % 100;
  const wiele = reszta >= 2 && reszta <= 4 && (dziesiatki < 12 || dziesiatki > 14);
  return String(ile) + (wiele ? ' sesje' : ' sesji');
}

/** Liczba modułów wraz z odmianą rzeczownika; polszczyzna rozróżnia trzy formy, a karta niesie tę miarę zdaniem, nie samą liczbą. */
function miaraModulow(ile: number): string {
  const reszta = ile % 10;
  const setka = ile % 100;
  if (ile === 1) return '1 moduł';
  if (reszta >= 2 && reszta <= 4 && (setka < 12 || setka > 14)) return `${ile} moduły`;
  return `${ile} modułów`;
}
