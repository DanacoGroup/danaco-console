/* Wiązanie Centrum z rdzeniem. Znacznik niesie biblioteka Właściciela — ten
   plik nic nie stawia, wypełnia stojące węzły odpowiedzią rdzenia. */

import {
  ChangeKind,
  Command,
  ErrorCode,
  EventType,
  ComponentKind,
  NavigationKind,
  ProgressStatus,
  WindowStatus,
  type Component,
  type Environment,
  type Module,
  type Session,
} from '../../../shared/contract.ts';
import type { Kanal, Wynik } from '../protokol/kanal.ts';
import { wywolaj } from '../protokol/wywolanie.ts';
import { miaraModulow, miaraSesji } from '../model/miary.ts';
import { zwiazWyborModulu } from './wybor-modulu.ts';
import { oglos } from './ogloszenie.ts';
import { zwiazPowiadomienia } from './powiadomienia.ts';
import { zwiazBadania, zwolnijBadania } from './research.ts';
import { zwiazTlumaczenie, zwolnijTlumaczenie } from './translate.ts';
import { zwiazDebate, zwolnijDebate } from './roundtable.ts';
import { zwiazOkno, zwiazOknoStojace } from './okno-modulu.ts';
import {
  nazwijKarte,
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
  kartaOkna,
  kartyOkna,
  oknaRobocze,
  oknoBiezace,
  oknoSesji,
  otworzKarte,
  otworzOkno,
  przelaczOkno,
  przypiszSesjeOkna,
  wskazKanalRdzenia,
  wskazKarte,
  zamknijOkno,
  zapiszKarte,
  zdejmijKarteOkna,
  zdejmijKartyPoPrawej,
  zostawKarte,
  zwiazZdarzeniaOkien,
  type KartaRobocza,
} from './okna-robocze.ts';
import {
  otworzSesje,
  przejmijOgnisko,
  sesjaBiezaca,
  uzgodnijSesjeOkna,
  wskazSrodowisko,
  zapomnijSesje,
} from './sesja-biezaca.ts';
import { zwiazStudio, zwolnijStudio } from './studio.ts';
import { zwiazLibrary, zwolnijLibrary } from './library.ts';
import { zwiazWorkspace, zwolnijWorkspace } from './workspace.ts';
import { zwiazAgents, zwolnijAgents } from './agents.ts';
import { zwiazDesign, zwolnijDesign } from './design.ts';
import { zwiazTerminal, zwolnijTerminal } from './terminal.ts';
import { zwiazPrzegladarke, zwolnijPrzegladarke } from './browser.ts';
import { otworzOperacjePlatformy } from './platforma-operacje.ts';
import { zglosUchwyt } from './zdarzenia.ts';

const KOD_MODULU_WYDANIA = 'studio';
const KOD_MODULU_BIBLIOTEKI = 'library';
const KOD_MODULU_WARSZTATU = 'workspace';
const KOD_MODULU_EKSPERTOW = 'agents';
const KOD_MODULU_DESIGN = 'design';
const KOD_MODULU_TERMINALA = 'terminal';
const KOD_MODULU_PRZEGLADARKI = 'browser';

const POWROT_NA_STRONE_GLOWNA =
  '[aria-label="Centrum dowodzenia"], .dn-karta-widoku--glowna, [data-wyjscie-modulu]';

interface WezlyCentrum {
  obszar: HTMLElement;
  wykazSesji: HTMLElement;
  plotno: HTMLElement;
  kartaGlowna: HTMLElement;
  kartaModulu: HTMLElement;
}

let zwiazane = false;

export function zwiazCentrum(kanal: Kanal | undefined = globalThis.DanacoKanal): boolean {
  if (zwiazane || kanal === undefined) return false;
  const wezly = zbierzWezly();
  if (wezly === null) return false;
  zwiazane = true;
  /* Kanał wchodzi do rejestru okien przy wiązaniu, nie przy odtworzeniu
     z `home.enter`: zamknięcie karty ma czym zamknąć okno także po odmowie. */
  wskazKanalRdzenia(kanal);

  const wzorWiersza = zdejmijWzorWiersza(wezly.wykazSesji);
  /* Pasmo kart i wykaz okien biorą wzory z treści przykładowej, więc
     przygotowanie pasma stoi przed jej zdjęciem. */
  przygotujPasmo();
  zdejmijTresciPrzykladowe();
  zdejmijDrogiBezPokrycia();
  const odswiez = (): void => {
    void odswiezWykaz(kanal, wezly.wykazSesji, wzorWiersza);
  };
  odswiez();
  zwiazPowiadomienia(kanal);
  void wypelnijKomponenty(kanal);
  void wypelnijProjekty(kanal);
  void wypelnijSrodowiska(kanal, wezly.obszar);
  const katalogModulow = new Map<string, Module>();
  const modulyPoId = new Map<string, Module>();
  void wczytajModuly(kanal, katalogModulow, modulyPoId);
  const katalogSrodowisk = new Map<string, Environment>();
  void wczytajSrodowiska(kanal, katalogSrodowisk);

  /* `rama.js` prowadzi drogi okien platformowych do pliku prototypu, którego
     wydanie nie niesie; bez tego nasłuchu naciśnięcie kończy się ciszą. */
  document.addEventListener('click', (zdarzenie) => {
    const cel = zdarzenie.target;
    if (!(cel instanceof Element)) return;
    for (const droga of OKNA_PLATFORMOWE) {
      const wezel = cel.closest<HTMLElement>(droga.wybor);
      if (wezel === null) continue;
      zdarzenie.stopPropagation();
      zdarzenie.preventDefault();
      if (droga.wybor === '#cd-konfiguracja' && kanal !== undefined) {
        otworzOperacjePlatformy(kanal);
      } else {
        zapowiedzOknoPlatformowe(droga, wezel);
      }
      return;
    }
  }, true);

  /* Sesja powstaje przy pierwszej wiadomości, więc „Nowa sesja" nie ma czego
     założyć w rejestrze rdzenia. Menu i szyna stoją poza obszarem okna. */
  document.addEventListener('click', (zdarzenie) => {
    const cel = zdarzenie.target;
    if (!(cel instanceof Element)) return;
    if (cel.closest<HTMLElement>('[data-okno-nowe]')?.dataset.oknoNowe !== '') return;
    otworzOkno();
    pokazOknoRobocze();
  });

  /* Faza przechwytywania odcina narrację `rama.js`, która opowiada o pozycji
     menu zdaniem prototypu zamiast wyniku rdzenia. */
  document.addEventListener('click', (zdarzenie) => {
    const cel = zdarzenie.target;
    if (!(cel instanceof Element)) return;
    const czynnosc = cel.closest<HTMLElement>('[data-poz-akcja]');
    const wykonaj = CZYNNOSCI_SESJI[czynnosc?.dataset.pozAkcja ?? ''];
    const idSesji = czynnosc?.dataset.idSesji;
    if (wykonaj === undefined || idSesji === undefined) return;
    zdarzenie.stopPropagation();
    void wykonaj(kanal, idSesji).then((wynik) => {
      if (wynik === null) return;
      if (!wynik.udany) oglos('Czynność sesji', wynik.blad?.message ?? 'Rdzeń odmówił wykonania.');
      odswiez();
    });
  }, true);

  document.addEventListener('click', (zdarzenie) => {
    const cel = zdarzenie.target;
    if (!(cel instanceof Element)) return;
    const czynnosc = cel.closest<HTMLElement>('[data-operacja]');
    const idKomponentu = czynnosc?.dataset.idKomponentu;
    if (czynnosc === null || idKomponentu === undefined) return;
    zdarzenie.stopPropagation();
    const operacja = czynnosc.dataset.operacja ?? '';
    void wykonajOperacjeKomponentu(kanal, operacja, idKomponentu).then((wynik) => {
      if (!wynik.udany) {
        oglos('Czynność komponentu', wynik.blad?.message ?? 'Rdzeń odmówił wykonania.', 'blad');
      }
      void wypelnijKomponenty(kanal);
    });
  }, true);

  // Kontrakt wymaga rodzaju, a niesie go wyłącznie kafel strefy komponentów.
  document.addEventListener('click', (zdarzenie) => {
    const cel = zdarzenie.target;
    if (!(cel instanceof Element)) return;
    if (cel.closest('#cd-dodaj-komponent') !== null) {
      zdarzenie.stopPropagation();
      oglos('Nowy komponent', 'Wskaż rodzaj kaflem strefy „Komponenty aplikacji" — '
        + 'kafel zakłada komponent własny tego rodzaju.');
      return;
    }
    const rodzaj = rodzajKomponentu(cel);
    if (rodzaj === null) return;
    zdarzenie.stopPropagation();
    void zalozKomponent(kanal, rodzaj).then((wynik) => {
      if (wynik === null) return;
      if (!wynik.udany) {
        oglos('Nowy komponent', wynik.blad?.message ?? 'Rdzeń odmówił założenia komponentu.', 'blad');
      }
      void wypelnijKomponenty(kanal);
    });
  }, true);


  /* Faza przechwytywania jest konieczna: grot wejścia biblioteki zatrzymuje
     zdarzenie na sobie i w fazie bąbelkowania nasłuch by go nie zobaczył. */
  let nazwaSrodowiska = '';
  let kartaBiezaca = '';

  const postawKarte = (karta: KartaRobocza, modul: Module): void => {
    if (MODULY_ZAPOWIEDZIANE.has(modul.code)) {
      zapowiedzModul(modul);
      return;
    }
    const wnetrze = wnetrzeKarty(wezly, karta.id, 'dn-tresc-' + modul.code);
    if (wnetrze === null) {
      zapowiedzModul(modul);
      return;
    }
    opiszGloweKarty(wnetrze.wezel, modul.name, nazwaSrodowiska);
    opiszPasekModulu(modul.name);
    pokazWidok(wezly, wnetrze.wezel);
    kartaBiezaca = karta.id;
    odswiezPasmo();
    if (wnetrze.nowe) {
      const idOkna = karta.idOknaKomunikacji;
      if (modul.code === KOD_MODULU_WYDANIA) {
        zwiazStudio(kanal, nazwaSrodowiska, idOkna, wnetrze.wezel);
      } else if (modul.code === KOD_MODULU_BIBLIOTEKI) {
        zwiazLibrary(kanal, nazwaSrodowiska, idOkna, wnetrze.wezel);
      } else if (modul.code === KOD_MODULU_WARSZTATU) {
        zwiazWorkspace(kanal, nazwaSrodowiska, idOkna, wnetrze.wezel);
      } else if (modul.code === KOD_MODULU_EKSPERTOW) {
        zwiazAgents(kanal, nazwaSrodowiska, idOkna, wnetrze.wezel);
      } else if (modul.code === KOD_MODULU_DESIGN) {
        zwiazDesign(kanal, nazwaSrodowiska, idOkna, wnetrze.wezel);
      } else if (modul.code === KOD_MODULU_TERMINALA) {
        zwiazTerminal(kanal, idOkna, wnetrze.wezel);
      } else if (modul.code === 'roundtable') {
        zwiazDebate(kanal, nazwaSrodowiska, idOkna, wnetrze.wezel);
      } else if (modul.code === 'translate') {
        zwiazTlumaczenie(kanal, nazwaSrodowiska, idOkna, wnetrze.wezel);
      } else if (modul.code === 'research') {
        zwiazBadania(kanal, nazwaSrodowiska, idOkna, wnetrze.wezel);
      } else if (modul.code === KOD_MODULU_PRZEGLADARKI) {
        zwiazPrzegladarke(kanal, nazwaSrodowiska, idOkna, wnetrze.wezel);
      } else if (idOkna === '') {
        zwiazOkno(kanal, modul.code, nazwaSrodowiska, wnetrze.wezel);
      } else {
        zwiazOknoStojace(kanal, karta.id, nazwaSrodowiska, idOkna, wnetrze.wezel);
      }
    }
    const nazwaPracy = nazwaPracyKarty(wnetrze.wezel);
    if (nazwaPracy !== '') {
      karta.nazwa = nazwaPracy;
      nazwijKarte(karta.id, nazwaPracy);
    }
  };

  const otworzModul = (modul: Module, idOkna: string): void => {
    if (szablonModulu('dn-tresc-' + modul.code) === null) {
      zapowiedzModul(modul);
      return;
    }
    const karta = kartaOkna(zapiszKarte(modul.code, modul.name, idOkna));
    if (karta !== undefined) postawKarte(karta, modul);
  };

  const otworzNowaKarte = (modul: Module): void => {
    if (szablonModulu('dn-tresc-' + modul.code) === null) {
      zapowiedzModul(modul);
      return;
    }
    const karta = kartaOkna(otworzKarte(modul.code, modul.name));
    if (karta !== undefined) postawKarte(karta, modul);
  };

  const pokazKarte = (idKarty: string): void => {
    const karta = wskazKarte(idKarty);
    const modul = karta === undefined ? undefined : katalogModulow.get(karta.kodModulu);
    if (karta === undefined || modul === undefined) return;
    postawKarte(karta, modul);
  };

  document.addEventListener('click', (zdarzenie) => {
    const cel = zdarzenie.target;
    if (!(cel instanceof Element)) return;
    const widok = pozycjaWyboru(cel, 'data-projekty-widok')?.dataset.projektyWidok;
    const porzadek = pozycjaWyboru(cel, 'data-projekty-sort')?.dataset.projektySort;
    if (widok === undefined && porzadek === undefined) return;
    zdarzenie.stopPropagation();
    if (widok === 'reakcja') {
      oglos('Drzewo projektów', 'Kontrakt nie niesie dla projektu oczekiwania na reakcję — '
        + 'tego wskazania nie da się dziś spełnić.');
      return;
    }
    if (porzadek !== undefined) porzadekProjektow = porzadek;
    oznaczWybor('#panel-projekty .dn-panel-drzewo', 'projektyWidok', widok);
    oznaczWybor('#panel-projekty .dn-panel-drzewo', 'projektySort', porzadek);
    void wypelnijProjekty(kanal);
  }, true);

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

  document.addEventListener('click', (zdarzenie) => {
    const cel = zdarzenie.target;
    if (!(cel instanceof Element)) return;
    const pozycja = cel.closest<HTMLElement>('[data-cd-wskazowka]');
    if (pozycja === null) return;
    zdarzenie.stopPropagation();
    ustawWskazowke(pozycja.getAttribute('aria-checked') !== 'true');
  }, true);

  document.addEventListener('click', (zdarzenie) => {
    const cel = zdarzenie.target;
    if (!(cel instanceof Element) || cel.closest('#cd-start-zamknij') === null) return;
    zdarzenie.stopPropagation();
    ustawWskazowke(false);
  }, true);

  document.addEventListener('click', (zdarzenie) => {
    const cel = zdarzenie.target;
    if (!(cel instanceof Element)) return;
    const przycisk = cel.closest<HTMLElement>('[data-panel-przypnij]');
    if (przycisk === null) return;
    zdarzenie.stopPropagation();
    const przypiety = przycisk.getAttribute('aria-pressed') === 'true';
    przycisk.setAttribute('aria-pressed', String(!przypiety));
  }, true);

  document.addEventListener('click', (zdarzenie) => {
    const cel = zdarzenie.target;
    if (!(cel instanceof Element)) return;
    const widok = pozycjaWyboru(cel, 'data-sesje-widok')?.dataset.sesjeWidok;
    const porzadek = pozycjaWyboru(cel, 'data-sesje-sort')?.dataset.sesjeSort;
    if (widok === undefined && porzadek === undefined) return;
    zdarzenie.stopPropagation();
    if (widok !== undefined) widokWykazu = widok;
    if (porzadek !== undefined) porzadekWykazu = porzadek;
    oznaczWybor('#wykaz-sesji', 'sesjeWidok', widok);
    oznaczWybor('#wykaz-sesji', 'sesjeSort', porzadek);
    odswiez();
  }, true);

  document.addEventListener('click', (zdarzenie) => {
    const cel = zdarzenie.target;
    if (!(cel instanceof Element) || cel.closest('[data-cd-odswiez]') === null) return;
    zdarzenie.stopPropagation();
    odswiez();
    void wypelnijSrodowiska(kanal, wezly.obszar);
  }, true);

  document.addEventListener('click', (zdarzenie) => {
    const cel = zdarzenie.target;
    if (!(cel instanceof Element)) return;
    const panel = document.getElementById('panel-samouczek');
    if (panel === null) return;
    if (cel.closest('[data-przelacz-samouczek]') !== null) {
      zdarzenie.stopPropagation();
      panel.hidden = !panel.hidden;
      oznaczStanPanelu(panel);
      return;
    }
    if (cel.closest('[data-zwin-samouczek]') !== null) {
      zdarzenie.stopPropagation();
      panel.hidden = true;
      oznaczStanPanelu(panel);
    }
  }, true);

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
      zamknijKarte(kartaBiezaca);
      return;
    }
    if (czynnosc === 'Zamknij pozostałe') {
      zdarzenie.stopPropagation();
      for (const idKarty of zostawKarte(kartaBiezaca)) zdejmijKarte(idKarty);
      uprzatnijWnetrza();
      odswiezPasmo();
      return;
    }
    if (czynnosc === 'Zamknij karty po prawej') {
      zdarzenie.stopPropagation();
      for (const idKarty of zdejmijKartyPoPrawej(kartaBiezaca)) zdejmijKarte(idKarty);
      uprzatnijWnetrza();
      odswiezPasmo();
      return;
    }
    if (czynnosc === 'Przypnij kartę') {
      zdarzenie.stopPropagation();
      przypnijKarte(kartaBiezaca);
    }
  }, true);

  function odswiezPasmo(): void {
    ustawKarty(
      kartyOkna().map((karta) => ({
        id: karta.id,
        nazwa: karta.nazwa,
        kodModulu: karta.kodModulu,
        idOknaKomunikacji: karta.idOknaKomunikacji,
      })),
      kartaBiezaca,
    );
    odswiezOknaRobocze();
  }

  const uprzatnijWnetrza = (): void => {
    for (const wnetrze of wezly.plotno.querySelectorAll<HTMLElement>('.cd-tresc--modul[data-karta]')) {
      const idKarty = wnetrze.dataset.karta ?? '';
      if (kartaOkna(idKarty) !== undefined) continue;
      zwolnijStudio(idKarty);
      zwolnijLibrary(idKarty);
      zwolnijWorkspace(idKarty);
      zwolnijAgents(idKarty);
      zwolnijDesign(idKarty);
      zwolnijTerminal(idKarty);
      zwolnijPrzegladarke(idKarty);
      zwolnijBadania(idKarty);
      zwolnijTlumaczenie(idKarty);
      zwolnijDebate(idKarty);
      wnetrze.remove();
    }
  };

  /* Wpis okna schodzi przed kartą — tam idzie `window.close` — a pasmo
     przerysowuje się na końcu. Karta bieżąca wraca najpierw na Centrum. */
  const zamknijKarte = (idKarty: string): void => {
    if (kartaBiezaca === idKarty) wrocDoCentrum();
    zdejmijKarteOkna(idKarty);
    zdejmijKarte(idKarty);
    uprzatnijWnetrza();
    odswiezPasmo();
  };

  const wrocDoCentrum = (): void => {
    pokazWidok(wezly, wezly.kartaGlowna);
    kartaBiezaca = '';
    opiszPasekModulu('');
    zapiszKarte('', '');
    zaznaczKarte('');
    odswiezOknaRobocze();
    odswiez();
  };

  function odswiezOknaRobocze(): void {
    ustawOknaRobocze(
      oknaRobocze().map((okno) => ({ id: okno.id, nazwa: okno.nazwa })),
      oknoBiezace().id,
    );
  }

  const pokazOknoRobocze = (): void => {
    uprzatnijWnetrza();
    const okno = oknoBiezace();
    void uzgodnijSesjeOkna(kanal, okno);
    if (kartaOkna(okno.kartaBiezaca) === undefined) {
      pokazWidok(wezly, wezly.kartaGlowna);
      kartaBiezaca = '';
      opiszPasekModulu('');
      odswiezPasmo();
      odswiez();
    } else {
      pokazKarte(okno.kartaBiezaca);
    }
    odswiezOknaRobocze();
  };

  odswiezOknaRobocze();
  /* Płótno zostaje na Centrum: to stan początkowy okna roboczego, a karta
     ogniskowana ostatnio stoi w paśmie i wraca na kliknięcie. */
  void przejmijOgnisko(kanal).then(() => {
    zapiszKarte('', '');
    odswiezPasmo();
  });
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
    (idKarty) => pokazKarte(idKarty),
    (idKarty) => zamknijKarte(idKarty),
  );

  zglosUchwyt(EventType.SessionChanged, (tresc) => {
    if (tresc.change === ChangeKind.Deleted) {
      zapomnijSesje(tresc.session.id);
      odswiezOknaRobocze();
    }
    odswiez();
    void wypelnijProjekty(kanal);
    if (tresc.change !== ChangeKind.Updated) void wypelnijSrodowiska(kanal, wezly.obszar);
  });
  zwiazZdarzeniaOkien((zdjete) => {
    if (zdjete.includes(kartaBiezaca)) wrocDoCentrum();
    for (const idKarty of zdjete) zdejmijKarte(idKarty);
    uprzatnijWnetrza();
    odswiezPasmo();
  });
  zglosUchwyt(EventType.WindowStateChanged, (tresc) => {
    if (kartyOkna().some((karta) => karta.idOknaKomunikacji === tresc.windowId)) odswiezPasmo();
  });
  zglosUchwyt(EventType.ConfigChanged, () => {
    void wczytajSrodowiska(kanal, katalogSrodowisk);
    void wczytajModuly(kanal, katalogModulow, modulyPoId);
    void wypelnijSrodowiska(kanal, wezly.obszar);
  });
  zglosUchwyt(EventType.QueueChanged, () => {
    odswiez();
  });
  zglosUchwyt(EventType.ProgressChanged, (tresc) => {
    if (tresc.status === ProgressStatus.Running || tresc.status === ProgressStatus.Pending) return;
    odswiez();
  });
  zglosUchwyt(EventType.DeviceChanged, (tresc) => {
    const wlasne = tresc.devices.find((urzadzenie) => urzadzenie.current);
    if (wlasne !== undefined && !wlasne.hasToken) {
      oglos('Urządzenie', 'Dostęp tego urządzenia do konta został unieważniony.', 'ostrzezenie');
    }
  });

  document.addEventListener('click', (zdarzenie) => {
    const cel = zdarzenie.target;
    if (!(cel instanceof Element)) return;
    const pozycja = cel.closest<HTMLElement>('[data-projekt-akcja]');
    const idProjektu = pozycja?.dataset.idProjektu;
    if (pozycja === null || idProjektu === undefined) return;
    zdarzenie.stopPropagation();
    void wykonajCzynnoscProjektu(kanal, pozycja.dataset.projektAkcja ?? '', idProjektu)
      .then((wynik) => {
        if (wynik === null) return;
        if (!wynik.udany) {
          oglos('Czynność projektu', wynik.blad?.message ?? 'Rdzeń odmówił wykonania.', 'blad');
        }
        void wypelnijProjekty(kanal);
        odswiez();
      });
  }, true);

  document.addEventListener('click', (zdarzenie) => {
    const cel = zdarzenie.target;
    if (!(cel instanceof Element) || cel.closest('[data-nowy-projekt-srodowisko]') === null) return;
    zdarzenie.stopPropagation();
    void zalozProjekt(kanal).then((wynik) => {
      if (wynik === null) return;
      if (!wynik.udany) {
        oglos('Nowy projekt', wynik.blad?.message ?? 'Rdzeń odmówił założenia projektu.', 'blad');
      }
      void wypelnijProjekty(kanal);
    });
  }, true);

  const wejdzWPrzedsionek = (kodSrodowiska: string): void => {
    if (kodSrodowiska === '') return;
    const srodowisko = katalogSrodowisk.get(kodSrodowiska);
    if (srodowisko === undefined) {
      oglos('Środowisko', 'Rejestr rdzenia nie zna tego środowiska, '
        + 'więc przedsionek nie ma czego pokazać.');
      return;
    }
    if ((srodowisko.moduleCodes?.length ?? 0) === 0) {
      oglos(srodowisko.name, zdanieBezModulow(srodowisko));
      return;
    }
    const widok = widokPrzedsionka(wezly, gniazdoPrzedsionka(kodSrodowiska));
    if (widok === null) return;
    pokazWidok(wezly, widok);
    nazwaSrodowiska = srodowisko.name;
    wskazSrodowisko(kodSrodowiska);
    void zwiazWyborModulu(kanal, srodowisko);
  };

  document.addEventListener('keydown', (zdarzenie) => {
    const cel = zdarzenie.target;
    if (!(cel instanceof Element)) return;
    const karta = cel.closest<HTMLElement>('.dn-karta-srodowiska');
    if (karta === null) return;
    const krok = zdarzenie.key === 'ArrowRight' ? 1 : (zdarzenie.key === 'ArrowLeft' ? -1 : 0);
    if (krok !== 0) {
      zdarzenie.preventDefault();
      przestawOgnisko(karta, krok);
      return;
    }
    if (zdarzenie.key !== 'Enter' && zdarzenie.key !== ' ') return;
    zdarzenie.preventDefault();
    wejdzWPrzedsionek(karta.dataset.srodowisko ?? '');
  });

  document.addEventListener('click', (zdarzenie) => {
    const cel = zdarzenie.target;
    if (!(cel instanceof Element)) return;

    // Pole wpisu stoi w przycisku wiersza: odstęp w nazwie uruchamiał wiersz.
    if (cel.closest('[contenteditable]') !== null) return;
    if (cel.querySelector('[contenteditable]') !== null) return;

    if (cel.closest(POWROT_NA_STRONE_GLOWNA) !== null) {
      wrocDoCentrum();
      return;
    }

    /* Sesja należy do jednego okna roboczego, więc kliknięcie przełącza na to
       okno, a gdy okno nie stoi — zakłada je z kartami jej okien komunikacji. */
    const wiersz = cel.closest<HTMLElement>('[data-id-sesji]');
    if (wiersz !== null && cel.closest('[data-menu]') === null
      && cel.closest('[data-poz-akcja]') === null) {
      zdarzenie.stopPropagation();
      const idSesji = wiersz.dataset.idSesji ?? '';
      const stojace = oknoSesji(idSesji);
      if (stojace !== undefined) {
        przelaczOkno(stojace.id);
        pokazOknoRobocze();
        return;
      }
      void otworzSesje(kanal, idSesji).then((otwarta) => {
        if (otwarta === null) {
          oglos('Sesja', 'Rdzeń odmówił otwarcia sesji.', 'blad');
          return;
        }
        const okno = otworzOkno();
        przypiszSesjeOkna(okno.id, otwarta.sesja.id, otwarta.sesja.title ?? '');
        let pierwsza = '';
        let pozaWydaniem = 0;
        for (const okienko of otwarta.okna) {
          if (okienko.status === WindowStatus.Closed) continue;
          const modul = modulyPoId.get(okienko.moduleId) ?? katalogModulow.get(okienko.moduleId);
          if (modul === undefined || szablonModulu('dn-tresc-' + modul.code) === null) {
            pozaWydaniem += 1;
            continue;
          }
          const idKarty = otworzKarte(modul.code, okienko.title ?? modul.name, okienko.id);
          if (pierwsza === '') pierwsza = idKarty;
        }
        if (pierwsza !== '') wskazKarte(pierwsza);
        pokazOknoRobocze();
        if (pierwsza === '') {
          oglos('Sesja', pozaWydaniem === 0
            ? 'Sesja stoi w nowym oknie roboczym; okna modułu jeszcze w niej nie ma.'
            : 'Sesja stoi w nowym oknie roboczym; jej okna należą do modułów spoza tego wydania.');
        }
      });
      return;
    }

    const kodNowejSesji = cel.closest<HTMLElement>('[data-nowa-sesja-srodowisko]')?.dataset
      .nowaSesjaSrodowisko;
    if (kodNowejSesji !== undefined) {
      zdarzenie.stopPropagation();
      otworzOkno();
      odswiezOknaRobocze();
      wejdzWPrzedsionek(kodNowejSesji);
      return;
    }

    const przelacznik = cel.closest<HTMLElement>('.dn-szyna-poz--srodowisko');
    if (przelacznik !== null) {
      const srodowisko = katalogSrodowisk.get(przelacznik.dataset.srodowisko ?? '');
      if (srodowisko !== undefined && (srodowisko.moduleCodes?.length ?? 0) === 0) {
        zdarzenie.stopPropagation();
        oglos(srodowisko.name, zdanieBezModulow(srodowisko));
      }
    }

    const kodSrodowiska = kodSrodowiskaWejscia(cel);
    if (kodSrodowiska !== '') {
      zdarzenie.stopPropagation();
      wejdzWPrzedsionek(kodSrodowiska);
      return;
    }

    const komponent = cel.closest<HTMLElement>('[data-otworz-komponent]')?.dataset.otworzKomponent;
    if (komponent !== undefined) {
      const modulKomponentu = katalogModulow.get(KOMPONENTY.get(komponent)?.kind ?? '');
      if (modulKomponentu === undefined) return;
      zdarzenie.stopPropagation();
      otworzModul(modulKomponentu, '');
      return;
    }

    const nowaKarta = cel.closest<HTMLElement>('[data-nowa-karta-modul]')?.dataset.nowaKartaModul;
    if (nowaKarta !== undefined) {
      const wskazany = [...katalogModulow.values()].find((modul) => modul.name === nowaKarta);
      if (wskazany === undefined) return;
      zdarzenie.stopPropagation();
      otworzNowaKarte(wskazany);
      return;
    }

    const kod = kodModulu(cel);
    if (kod === '') return;
    const modul = katalogModulow.get(kod);
    if (modul === undefined) return;
    zdarzenie.stopPropagation();
    nazwaSrodowiska = nazwaSrodowiskaWejscia(cel) || nazwaSrodowiska;
    /* Pionowa szyna jest paskiem szybkiego dostępu: otwiera OKNO ROBOCZE
       z pominięciem przedsionka, a kafle otwierają kartę w oknie stojącym. */
    if (cel.closest('.dn-szyna-poz--modul') !== null) otworzOkno();
    otworzModul(modul, '');
  }, true);

  return true;
}

/* Wykaz i drzewo niosą tę samą cechę stanem bieżącym, więc wskazanie czyta się
   wyłącznie z pozycji menu — inaczej klik w wykazie zjadałby kliknięcia wierszy. */
function pozycjaWyboru(cel: Element, cecha: string): HTMLElement | null {
  return cel.closest<HTMLElement>(`.sta-menu-poz[${cecha}]`);
}

function oznaczWybor(wyborPojemnika: string, cecha: string, wartosc: string | undefined): void {
  if (wartosc === undefined) return;
  const pojemnik = document.querySelector<HTMLElement>(wyborPojemnika);
  if (pojemnik !== null) pojemnik.dataset[cecha] = wartosc;
  const nazwaCechy = 'data-' + cecha.replace(/[A-Z]/g, (znak) => '-' + znak.toLowerCase());
  for (const pozycja of document.querySelectorAll(`.sta-menu-poz[${nazwaCechy}]`)) {
    pozycja.setAttribute('aria-checked', String(pozycja.getAttribute(nazwaCechy) === wartosc));
  }
}

/* Każde środowisko ma własny prototyp przedsionka, bo kafle niosą znaki swoich
   modułów; środowisko bez prototypu wchodzi na przedsionek TalkIn. */
function gniazdoPrzedsionka(kodSrodowiska: string): string {
  return 'dn-tresc-przedsionek-' + kodSrodowiska.toLowerCase();
}

function kodModulu(cel: Element): string {
  const wskazanie =
    cel.closest('.dn-szyna-poz--modul')?.getAttribute('data-modul') ??
    cel.closest('.pd-kafel')?.getAttribute('data-modul') ??
    // Skrót szybkiego wyboru nazywa rodzaj komponentu, a rodzaj jest kodem modułu.
    cel.closest('.dn-szyna-poz--skrot')?.getAttribute('data-skrot-komponent') ??
    '';
  return wskazanie.toLowerCase();
}

async function wczytajSrodowiska(kanal: Kanal, spis: Map<string, Environment>): Promise<void> {
  const wynik = await wywolaj(kanal, Command.EnvironmentList, { includeModules: true });
  if (!wynik.udany || wynik.wynik === undefined) return;
  for (const srodowisko of wynik.wynik.environments) spis.set(srodowisko.code, srodowisko);
}

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

/* Komunikat nazywa niegotowość wprost i podaje opis modułu z rejestru:
   twierdzenie, że okno się otwiera, byłoby nieprawdą. */
/* Moduły, których znacznik prototypu niesie rozmowę, wykaz i liczniki wymyślone
   na pokaz, a wydanie nie ma czym ich zastąpić. Plan etapu 2 zna dwa stany okna
   i tylko dwa: okno działa albo jest zapowiedziane. Operacje tych rodzin zostają
   osiągalne katalogiem „Operacje platformy" w Centrum. */
const MODULY_ZAPOWIEDZIANE = new Set([
  'developer',
  'apps', 'automations', 'assistant', 'diagnostics',
]);

function zapowiedzModul(modul: Module | undefined): void {
  if (modul === undefined) return;
  const opis = modul.description ?? '';
  const zdanie = opis === '' ? '' : opis + ' ';
  oglos(modul.name, zdanie + 'Okno tego modułu nie wchodzi do tego wydania.');
}

function kodSrodowiskaWejscia(cel: Element): string {
  return cel.closest<HTMLElement>('.dn-karta-srodowiska')?.dataset.srodowisko ?? '';
}

function przestawOgnisko(karta: HTMLElement, krok: number): void {
  const karty = [...document.querySelectorAll<HTMLElement>('.dn-karta-srodowiska')];
  if (karty.length === 0) return;
  const numer = karty.indexOf(karta);
  const nastepna = karty[(numer + krok + karty.length) % karty.length];
  for (const kandydat of karty) kandydat.tabIndex = -1;
  nastepna.tabIndex = 0;
  nastepna.focus();
}

function ustawWskazowke(widoczna: boolean): void {
  const wskazowka = document.querySelector<HTMLElement>('.cd-wskazowka, #cd-start');
  if (wskazowka !== null) wskazowka.hidden = !widoczna;
  for (const pozycja of document.querySelectorAll('[data-cd-wskazowka]')) {
    pozycja.setAttribute('aria-checked', String(widoczna));
  }
}

function oznaczStanPanelu(panel: HTMLElement): void {
  for (const pozycja of document.querySelectorAll('[data-przelacz-samouczek][role="menuitemcheckbox"]')) {
    pozycja.setAttribute('aria-checked', String(!panel.hidden));
  }
}

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

/* Widoki kart stoją obok siebie w płótnie i różnią się zasłoną — przełączenie
   nie rusza stanu DOM karty drugiej. */
function pokazWidok(wezly: WezlyCentrum, widok: HTMLElement): void {
  for (const kandydat of wezly.plotno.querySelectorAll<HTMLElement>('.cd-tresc')) {
    kandydat.hidden = kandydat !== widok;
  }
  widok.hidden = false;
}

/* Pusty wykaz modułów znaczy dwie różne rzeczy. Przy nawigacji orkiestracji
   jest kształtem zamierzonym (migracja 072), więc zdanie nazywa brak panelu
   w wydaniu; poza nią zostaje brakiem w rejestrze rdzenia. */
function zdanieBezModulow(srodowisko: Environment): string {
  if (srodowisko.navigationKind === NavigationKind.Orchestration) {
    return 'To środowisko prowadzi się panelem orkiestracji, a nie listą modułów. '
      + 'Panel nie wchodzi do tego wydania; operacje zespołów i przebiegów stoją '
      + 'w katalogu „Operacje platformy".';
  }
  return 'Rejestr rdzenia nie wskazuje dla tego środowiska ani jednego modułu, '
    + 'więc nie ma czego pokazać.';
}

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

function szablonModulu(gniazdo: string): Element | null {
  const szablon = document.getElementById(gniazdo);
  if (!(szablon instanceof HTMLTemplateElement)) return null;
  return szablon.content.firstElementChild;
}

function wnetrzeStojace(wezly: WezlyCentrum, idKarty: string): HTMLElement | null {
  return wezly.plotno.querySelector<HTMLElement>(`.cd-tresc--modul[data-karta="${idKarty}"]`);
}

/* Każda karta ma własny węzeł: dwie karty tego samego modułu nie dzielą
   wnętrza, więc przełączenie nie czyści rozmowy. */
function wnetrzeKarty(
  wezly: WezlyCentrum,
  idKarty: string,
  gniazdo: string,
): { wezel: HTMLElement; nowe: boolean } | null {
  const stojace = wnetrzeStojace(wezly, idKarty);
  if (stojace !== null) return { wezel: stojace, nowe: false };
  const blok = szablonModulu(gniazdo);
  if (blok === null) return null;
  const wezel = wezly.kartaModulu.cloneNode(true) as HTMLElement;
  wezel.removeAttribute('id');
  wezel.dataset.karta = idKarty;
  for (const dziecko of [...wezel.children]) {
    if (!dziecko.classList.contains('cd-modul-glowa')) dziecko.remove();
  }
  wezel.appendChild(blok.cloneNode(true));
  wezly.plotno.appendChild(wezel);
  return { wezel, nowe: true };
}

function oddajPlik(nazwa: string, tresc: string): void {
  const adres = URL.createObjectURL(new Blob([tresc], { type: 'text/plain;charset=utf-8' }));
  const odnosnik = document.createElement('a');
  odnosnik.href = adres;
  odnosnik.download = nazwa;
  odnosnik.click();
  URL.revokeObjectURL(adres);
}

function nazwaPracyKarty(wnetrze: Element): string {
  const znacznik = wnetrze.querySelector('.sta-okno-znacznik, .st-wstazka-sesja span');
  return znacznik?.textContent?.trim() ?? '';
}

function opiszPasekModulu(nazwa: string): void {
  const pole = document.querySelector('[data-pasek-modul-nazwa]');
  if (pole !== null) pole.textContent = nazwa;
}

function opiszGloweKarty(wnetrze: HTMLElement, nazwaModulu: string, nazwaSesji: string): void {
  const nazwa = wnetrze.querySelector('[data-karta-modul-nazwa]');
  if (nazwa !== null) nazwa.textContent = nazwaModulu;
  const meta = wnetrze.querySelector('.cd-modul-glowa .dn-meta');
  if (meta !== null) meta.textContent = nazwaSesji === '' ? '' : 'sesja: ' + nazwaSesji;
  /* Pas kontekstu niesie w prototypie `aria-label` na dzielniku bez roli, czego
     ARIA zabrania. Rola grupy czyni z niego pojemnik nazwany — poprawka stoi tu,
     bo znacznik jest powielony w dziewiętnastu plikach warstwy projektowej. */
  for (const pas of wnetrze.querySelectorAll('.sta-kom-kontekst[aria-label]')) {
    pas.setAttribute('role', 'group');
  }
}

const KOMPONENTY = new Map<string, Component>();

let wzorKomponentu: HTMLElement | null = null;

/* Wykaz pustoszeje przed pytaniem rdzenia: przy odmowie na ekranie ma stać
   stan pusty ze znacznika, nie komponenty z prototypu. */
async function wypelnijKomponenty(kanal: Kanal): Promise<void> {
  const wykaz = document.getElementById('cd-wlasne');
  if (wykaz === null) return;
  wzorKomponentu = zapamietajWzor(wzorKomponentu, wykaz, '.cd-wlasny');
  const wzor = wzorKomponentu;
  wykaz.replaceChildren();
  if (wzor === null) return;
  const wynik = await wywolaj(kanal, Command.ComponentList, {});
  if (!wynik.udany || wynik.wynik === undefined) {
    oglos('Komponenty własne', wynik.blad?.message ?? 'Rdzeń odmówił wykazu komponentów.', 'blad');
    return;
  }
  KOMPONENTY.clear();
  for (const komponent of wynik.wynik.components) {
    KOMPONENTY.set(komponent.id, komponent);
    wykaz.appendChild(zbudujKomponent(wzor, komponent));
  }
}

function zbudujKomponent(wzor: HTMLElement, komponent: Component): HTMLElement {
  const pozycja = wzor.cloneNode(true) as HTMLElement;
  for (const otworz of pozycja.querySelectorAll<HTMLElement>('[data-otworz-komponent]')) {
    otworz.dataset.otworzKomponent = komponent.id;
  }
  const nazwa = pozycja.querySelector('.dn-kafel-nazwa');
  if (nazwa !== null) nazwa.textContent = komponent.name;
  for (const podpis of pozycja.querySelectorAll('.dn-kafel-rodzaj, .dn-kafel-opis')) {
    if (komponent.description === undefined || komponent.description === '') podpis.remove();
    else podpis.textContent = komponent.description;
  }
  const menu = pozycja.querySelector('[data-menu-tresc]');
  const wyzwalacz = pozycja.querySelector('[data-menu]');
  if (menu !== null && wyzwalacz !== null) {
    const oznaczenie = 'menu-komponentu-' + komponent.id;
    menu.id = oznaczenie;
    wyzwalacz.setAttribute('data-menu', oznaczenie);
  }
  zdejmijOperacjeBezZrodla(pozycja);
  for (const czynnosc of pozycja.querySelectorAll<HTMLElement>('[data-operacja]')) {
    czynnosc.dataset.idKomponentu = komponent.id;
  }
  return pozycja;
}

function rodzajKomponentu(cel: Element): ComponentKind | null {
  const kafel = cel.closest('.cd-kafel');
  const wskazanie = (kafel?.getAttribute('data-komponent') ?? '').toLowerCase();
  const rodzaje: readonly string[] = Object.values(ComponentKind);
  return rodzaje.includes(wskazanie) ? (wskazanie as ComponentKind) : null;
}

// Okna zakładania komponentu wydanie nie niesie: nazwa wchodzi w pozycji wykazu.
async function zalozKomponent(kanal: Kanal, rodzaj: ComponentKind): Promise<Wynik<unknown> | null> {
  const wykaz = document.getElementById('cd-wlasne');
  if (wykaz === null || wzorKomponentu === null) return null;
  const pozycja = wzorKomponentu.cloneNode(true) as HTMLElement;
  for (const zbedne of pozycja.querySelectorAll('.cd-wlasny-menu, [data-menu-tresc], .dn-kafel-rodzaj, .dn-kafel-opis')) {
    zbedne.remove();
  }
  const nazwa = pozycja.querySelector<HTMLElement>('.dn-kafel-nazwa');
  if (nazwa === null) return null;
  await new Promise((gotowe) => setTimeout(gotowe, 0));
  wykaz.prepend(pozycja);
  const wpis = await zapytajWWezle(nazwa, '');
  pozycja.remove();
  if (wpis === null || wpis === '') return null;
  return wywolaj(kanal, Command.ComponentCreate, { kind: rodzaj, name: wpis });
}

// Kontrakt nie zna kopiowania komponentu: kopia zakłada go z definicji stojącej.
const OPERACJE_KOMPONENTU = new Set(['duplikuj', 'usun']);

async function wykonajOperacjeKomponentu(
  kanal: Kanal,
  operacja: string,
  idKomponentu: string,
): Promise<Wynik<unknown>> {
  const zrodlo = KOMPONENTY.get(idKomponentu);
  if (zrodlo === undefined) {
    return {
      udany: false,
      blad: {
        code: ErrorCode.NotFound,
        message: 'Komponent zniknął z wykazu Centrum — odśwież widok.',
        retryable: false,
      },
    };
  }
  if (operacja === 'usun') {
    return wywolaj(kanal, Command.ComponentDelete, { componentId: idKomponentu });
  }
  return wywolaj(kanal, Command.ComponentCreate, {
    kind: zrodlo.kind,
    name: `${zrodlo.name} — kopia`,
    description: zrodlo.description,
    config: zrodlo.config,
  });
}

function zdejmijOperacjeBezZrodla(pozycja: HTMLElement): void {
  for (const czynnosc of pozycja.querySelectorAll('[data-menu-tresc] .sta-menu-poz')) {
    const operacja = czynnosc.getAttribute('data-operacja');
    if (operacja === null) {
      if (!czynnosc.hasAttribute('data-otworz-komponent')) czynnosc.remove();
      continue;
    }
    if (!OPERACJE_KOMPONENTU.has(operacja)) czynnosc.remove();
  }
  zdejmijRozdzielnikiSieroce(pozycja);
}

function zapamietajWzor(
  zapamietany: HTMLElement | null,
  wykaz: HTMLElement,
  wybor: string,
): HTMLElement | null {
  if (zapamietany !== null) return zapamietany;
  const wezel = wykaz.querySelector<HTMLElement>(wybor);
  return wezel === null ? null : (wezel.cloneNode(true) as HTMLElement);
}

async function wypelnijProjekty(kanal: Kanal): Promise<void> {
  const drzewo = document.querySelector<HTMLElement>('#panel-projekty .dn-panel-drzewo');
  if (drzewo === null) return;
  wzorGalezi = zapamietajWzor(wzorGalezi, drzewo, '.dn-panel-galaz');
  const wzor = wzorGalezi;
  const wzorSesji = wzor?.querySelector<HTMLElement>('.dn-panel-wiersz') ?? null;
  // Drzewo pustoszeje przed pytaniem: przy odmowie ma stać puste, nie z prototypu.
  drzewo.replaceChildren();
  if (wzor === null) return;
  const projekty = await wywolaj(kanal, Command.ProjectList, {});
  const sesje = await wywolaj(kanal, Command.SessionList, {});
  if (!projekty.udany || projekty.wynik === undefined) {
    oglos('Drzewo projektów', projekty.blad?.message ?? 'Rdzeń odmówił wykazu projektów.', 'blad');
    return;
  }
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
  for (const projekt of uporzadkowane) {
    const wpis = wzor.cloneNode(true) as HTMLElement;
    wpis.dataset.projekt = projekt.id;
    const nazwa = wpis.querySelector('.dn-panel-galaz-nazwa');
    if (nazwa !== null) nazwa.textContent = projekt.name;
    for (const pozycja of wpis.querySelectorAll('.dn-panel-wiersz')) pozycja.remove();
    opiszMenuProjektu(wpis, projekt.id);
    for (const sesja of wedlugProjektu.get(projekt.id) ?? []) {
      if (wzorSesji === null) break;
      wpis.appendChild(zbudujWiersz(wzorSesji, sesja, 'menu-sesji-projektu-'));
    }
    drzewo.appendChild(wpis);
  }
}

const CZYNNOSCI_PROJEKTU: Record<string, string> = {
  'Zmień nazwę': 'nazwa',
  'Usuń trwale': 'usun',
};

function opiszMenuProjektu(wpis: HTMLElement, idProjektu: string): void {
  const menu = wpis.querySelector(':scope > [data-menu-tresc]');
  const wyzwalacz = wpis.querySelector(':scope > summary [data-menu]');
  if (menu !== null && wyzwalacz !== null) {
    const oznaczenie = 'menu-projektu-' + idProjektu;
    menu.id = oznaczenie;
    wyzwalacz.setAttribute('data-menu', oznaczenie);
  }
  for (const pozycja of wpis.querySelectorAll<HTMLElement>(':scope > [data-menu-tresc] .sta-menu-poz')) {
    const czynnosc = CZYNNOSCI_PROJEKTU[(pozycja.textContent ?? '').trim()];
    if (czynnosc === undefined) {
      pozycja.remove();
      continue;
    }
    pozycja.dataset.projektAkcja = czynnosc;
    pozycja.dataset.idProjektu = idProjektu;
  }
  zdejmijRozdzielnikiSieroce(wpis);
}

async function wykonajCzynnoscProjektu(
  kanal: Kanal,
  czynnosc: string,
  idProjektu: string,
): Promise<Wynik<unknown> | null> {
  if (czynnosc === 'usun') return wywolaj(kanal, Command.ProjectDelete, { projectId: idProjektu });
  if (czynnosc !== 'nazwa') return null;
  /* Wpis zaczyna się po domknięciu menu przez bibliotekę: jej powrót ogniska
     na wyzwalacz domknąłby pole wpisu zdarzeniem blur. */
  await new Promise((gotowe) => setTimeout(gotowe, 0));
  pokazPanelProjektow();
  const wezel = wezelNazwyProjektu(idProjektu);
  const nazwa = wezel === null ? null : await zapytajWWezle(wezel, null);
  if (nazwa === null || nazwa === '') return null;
  return wywolaj(kanal, Command.ProjectRename, { projectId: idProjektu, name: nazwa });
}

function wezelNazwyProjektu(idProjektu: string): HTMLElement | null {
  for (const galaz of document.querySelectorAll<HTMLElement>('#panel-projekty .dn-panel-galaz')) {
    if (galaz.dataset.projekt === idProjektu) {
      return galaz.querySelector<HTMLElement>('.dn-panel-galaz-nazwa');
    }
  }
  return null;
}

async function zalozProjekt(kanal: Kanal): Promise<Wynik<unknown> | null> {
  const drzewo = document.querySelector<HTMLElement>('#panel-projekty .dn-panel-drzewo');
  if (drzewo === null || wzorGalezi === null) return null;
  const galaz = wzorGalezi.cloneNode(true) as HTMLElement;
  galaz.removeAttribute('data-projekt');
  for (const zbedne of galaz.querySelectorAll('.dn-panel-wiersz, .dn-obszar-pozycja-menu, [data-menu-tresc]')) {
    zbedne.remove();
  }
  const nazwa = galaz.querySelector<HTMLElement>('.dn-panel-galaz-nazwa');
  if (nazwa === null) return null;
  /* Wpis zaczyna się po domknięciu menu przez bibliotekę: jej powrót ogniska
     na wyzwalacz zabrałby ognisko polu wpisu. */
  await new Promise((gotowe) => setTimeout(gotowe, 0));
  pokazPanelProjektow();
  drzewo.prepend(galaz);
  const wpis = await zapytajWWezle(nazwa, '');
  galaz.remove();
  if (wpis === null || wpis === '') return null;
  return wywolaj(kanal, Command.ProjectCreate, { name: wpis });
}

function pokazPanelProjektow(): void {
  const zakladka = document.querySelector<HTMLElement>('[role="tab"][aria-controls="panel-projekty"]');
  if (zakladka === null || zakladka.getAttribute('aria-selected') === 'true') return;
  zakladka.click();
  if (document.getElementById('panel-projekty')?.hidden !== true) return;
  // Biblioteka nie przełączyła zakładki, więc panele przełącza wiązanie.
  for (const inna of zakladka.parentElement?.querySelectorAll<HTMLElement>('[role="tab"]') ?? []) {
    inna.setAttribute('aria-selected', String(inna === zakladka));
    const panel = document.getElementById(inna.getAttribute('aria-controls') ?? '');
    if (panel !== null) panel.hidden = inna !== zakladka;
  }
}

function zdejmijWzorWiersza(wykaz: HTMLElement): HTMLElement | null {
  const wiersz = wykaz.querySelector<HTMLElement>('.dn-panel-wiersz');
  return wiersz === null ? null : (wiersz.cloneNode(true) as HTMLElement);
}

/* Wykaz pustoszeje przed pytaniem: przy odmowie panel ma stać pusty, a nie
   sesjami, których nie ma. */
async function odswiezWykaz(
  kanal: Kanal,
  wykaz: HTMLElement,
  wzor: HTMLElement | null,
): Promise<void> {
  if (wzor === null) return;
  wykaz.replaceChildren();
  const wynik = await wywolaj(kanal, Command.SessionList, {});
  if (!wynik.udany || wynik.wynik === undefined) {
    oglos('Wykaz sesji', wynik.blad?.message ?? 'Rdzeń odmówił wykazu sesji.', 'blad');
    return;
  }
  for (const sesja of uporzadkuj(przesiej(wynik.wynik.sessions))) {
    wykaz.appendChild(zbudujWiersz(wzor, sesja));
  }
}

let widokWykazu = 'wszystkie';
let porzadekWykazu = 'czynnosc';
let porzadekProjektow = 'czynnosc';

function przesiej(sesje: Session[]): Session[] {
  if (widokWykazu === 'czynne') return sesje.filter((sesja) => sesja.status === 'active');
  if (widokWykazu === 'zakonczone') return sesje.filter((sesja) => sesja.status !== 'active');
  if (widokWykazu === 'reakcja') return sesje.filter((sesja) => sesja.awaitingReaction === true);
  return sesje;
}

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

function zbudujWiersz(wzor: HTMLElement, sesja: Session, przedrostek = 'menu-sesji-'): HTMLElement {
  const wiersz = wzor.cloneNode(true) as HTMLElement;
  wiersz.dataset.idSesji = sesja.id;
  const nazwa = wiersz.querySelector('.dn-obszar-pozycja-nazwa');
  // Identyfikator jest oznaczeniem magazynu, nie nazwą pracy Operatora.
  if (nazwa !== null) nazwa.textContent = sesja.title ?? 'Sesja bez nazwy';
  wiersz.querySelector('.dn-obszar-pozycja')?.setAttribute('data-id-sesji', sesja.id);
  opiszStanWiersza(wiersz, sesja);
  const menu = wiersz.querySelector('[data-menu-tresc]');
  const wyzwalacz = wiersz.querySelector('[data-menu]');
  if (menu !== null && wyzwalacz !== null) {
    const oznaczenie = przedrostek + sesja.id;
    menu.id = oznaczenie;
    wyzwalacz.setAttribute('data-menu', oznaczenie);
  }
  zdejmijCzynnosciBezZrodla(wiersz);
  /* Identyfikator sesji siada na samej pozycji menu: biblioteka menu przenosi
     treść menu poza wiersz, więc szukanie sesji w przodkach nic nie znajdzie. */
  for (const pozycja of wiersz.querySelectorAll<HTMLElement>('[data-poz-akcja]')) {
    pozycja.dataset.idSesji = sesja.id;
  }
  return wiersz;
}

function opiszStanWiersza(wiersz: HTMLElement, sesja: Session): void {
  const stan = sesja.status !== 'active'
    ? 'zakonczone'
    : (sesja.awaitingReaction === true ? 'reakcja' : 'praca');
  for (const wezel of wiersz.querySelectorAll('[data-stan]')) {
    wezel.setAttribute('data-stan', stan);
  }
}

/* Spis trzyma wywołania, nie same nazwy komend: tylko wywołanie zapisane przy
   swojej komendzie daje się sprawdzić kontraktem przy budowaniu. */

const CZYNNOSCI_SESJI: Record<
  string,
  (kanal: Kanal, idSesji: string) => Promise<Wynik<unknown> | null>
> = {
  nazwa: (kanal, idSesji) => zmienNazweSesji(kanal, idSesji),
  przenies: (kanal, idSesji) => przeniesSesjeDoProjektu(kanal, idSesji),
  archiwizuj: (kanal, idSesji) =>
    wywolaj(kanal, Command.SessionArchive, { sessionIds: [idSesji] }),
  usun: (kanal, idSesji) =>
    wywolaj(kanal, Command.SessionDelete, { sessionIds: [idSesji], confirm: true }),
  wyjmij: (kanal, idSesji) =>
    wywolaj(kanal, Command.SessionProjectClear, { sessionIds: [idSesji] }),
};

async function zmienNazweSesji(kanal: Kanal, idSesji: string): Promise<Wynik<unknown> | null> {
  const nazwa = await zapytajWWierszu(idSesji, null);
  if (nazwa === null || nazwa === '') return null;
  return wywolaj(kanal, Command.SessionRename, { sessionId: idSesji, title: nazwa });
}

/* Kontrakt zakłada projekt nowy, gdy nie wskazano istniejącego: nazwa wpisana
   w wierszu sesji wystarcza za wskazanie. */
async function przeniesSesjeDoProjektu(
  kanal: Kanal,
  idSesji: string,
): Promise<Wynik<unknown> | null> {
  const nazwa = await zapytajWWierszu(idSesji, '');
  if (nazwa === null || nazwa === '') return null;
  return wywolaj(kanal, Command.SessionProjectSet, { sessionIds: [idSesji], projectName: nazwa });
}

function zapytajWWierszu(idSesji: string, wartosc: string | null): Promise<string | null> {
  const wezel = wezelNazwySesji(idSesji);
  return wezel === null ? Promise.resolve(null) : zapytajWWezle(wezel, wartosc);
}

function zapytajWWezle(wezel: HTMLElement, wartosc: string | null): Promise<string | null> {
  const przed = wezel.textContent ?? '';
  return new Promise((rozstrzygnij) => {
    let domkniete = false;
    const domknij = (wpis: string | null): void => {
      if (domkniete) return;
      domkniete = true;
      wezel.removeAttribute('contenteditable');
      wezel.textContent = przed;
      rozstrzygnij(wpis);
    };
    wezel.addEventListener('keydown', (zdarzenie) => {
      if (zdarzenie.key === 'Enter') {
        zdarzenie.preventDefault();
        domknij((wezel.textContent ?? '').trim());
        return;
      }
      if (zdarzenie.key === 'Escape') {
        zdarzenie.preventDefault();
        domknij(null);
        return;
      }
      // Odstęp uruchamiał przycisk wiersza, więc znak wchodzi w miejsce kursora.
      if (zdarzenie.key === ' ') {
        zdarzenie.preventDefault();
        wstawZnak(' ');
      }
    });
    wezel.addEventListener('blur', () => domknij(null));
    wezel.setAttribute('contenteditable', 'plaintext-only');
    if (wartosc !== null) wezel.textContent = wartosc;
    wezel.focus();
    zaznaczCalosc(wezel);
  });
}

/* Sesja stoi w wykazie sesji i w drzewie projektów naraz; pierwszeństwo ma
   wiersz widoczny, bo pole wpisu w panelu zasłoniętym nie przyjmie ogniska. */
function wezelNazwySesji(idSesji: string): HTMLElement | null {
  const wiersze = [...document.querySelectorAll<HTMLElement>(
    '#wykaz-sesji .dn-panel-wiersz, #panel-projekty .dn-panel-wiersz',
  )].filter((wiersz) => wiersz.dataset.idSesji === idSesji);
  const wiersz = wiersze.find((kandydat) => kandydat.offsetParent !== null) ?? wiersze[0];
  return wiersz?.querySelector<HTMLElement>('.dn-obszar-pozycja-nazwa') ?? null;
}

function wstawZnak(znak: string): void {
  const zaznaczenie = window.getSelection();
  if (zaznaczenie === null || zaznaczenie.rangeCount === 0) return;
  const zakres = zaznaczenie.getRangeAt(0);
  zakres.deleteContents();
  const wezel = document.createTextNode(znak);
  zakres.insertNode(wezel);
  zakres.setStartAfter(wezel);
  zakres.collapse(true);
  zaznaczenie.removeAllRanges();
  zaznaczenie.addRange(zakres);
}

function zaznaczCalosc(wezel: HTMLElement): void {
  const zakres = document.createRange();
  zakres.selectNodeContents(wezel);
  const zaznaczenie = window.getSelection();
  zaznaczenie?.removeAllRanges();
  zaznaczenie?.addRange(zakres);
}

function zdejmijCzynnosciBezZrodla(wiersz: HTMLElement): void {
  for (const pozycja of wiersz.querySelectorAll('[data-menu-tresc] .sta-menu-poz')) {
    const czynnosc = pozycja.getAttribute('data-poz-akcja') ?? '';
    if (CZYNNOSCI_SESJI[czynnosc] === undefined) pozycja.remove();
  }
  zdejmijRozdzielnikiSieroce(wiersz);
}

function zdejmijRozdzielnikiSieroce(obudowa: HTMLElement): void {
  for (const rozdzielnik of obudowa.querySelectorAll('[data-menu-tresc] .sta-menu-sep')) {
    const przed = rozdzielnik.previousElementSibling;
    if (przed === null || przed.classList.contains('sta-menu-sep')) rozdzielnik.remove();
  }
}

function nazwaSrodowiskaWejscia(cel: Element): string {
  const grupa = cel.closest('.dn-szyna-poz--modul')?.closest('.dn-szyna-moduly');
  if (grupa === null || grupa === undefined) return '';
  const przelacznik = document.querySelector(
    `.dn-szyna-poz--srodowisko[aria-controls="${grupa.id}"]`,
  );
  return przelacznik?.getAttribute('aria-label') ?? '';
}


/* Karty okien poza główną i odsyłacz do pliku prototypu prowadzą poza produkt.
   Pusty wykaz odsłania stan pusty ze znacznika. */
function zdejmijTresciPrzykladowe(): void {
  const karty = document.querySelector('.dn-karty-lista');
  if (karty !== null) {
    for (const karta of [...karty.querySelectorAll('.dn-karta-widoku')].slice(1)) karta.remove();
  }
  document.querySelector('.cd-modul-odnosnik')?.remove();
  zdejmijZapowiedziPrototypu();
}

/* Cecha `data-komunikat` niosła w prototypie opowieść o skutku naciśnięcia;
   produkt tej warstwy nie wczytuje, więc cecha obiecuje komunikat bez mówcy. */
function zdejmijZapowiedziPrototypu(): void {
  const korzenie: ParentNode[] = [document];
  for (const szablon of document.querySelectorAll('template')) korzenie.push(szablon.content);
  for (const korzen of korzenie) {
    for (const wezel of korzen.querySelectorAll('[data-komunikat]')) {
      wezel.removeAttribute('data-komunikat');
      wezel.removeAttribute('data-komunikat-tytul');
      wezel.removeAttribute('data-komunikat-rodzaj');
    }
  }
}

/* Szukania w oknie nie obsługuje ani wiązanie, ani biblioteka warstwy
   projektowej, a kontrakt nie ma komendy, po której zakres tego szukania
   dałoby się poznać. Przycisk schodzi zamiast stać martwy. */
const DROGI_BEZ_POKRYCIA = [
  '.cd-sekcja-glowa [data-operacja]',
  '[data-etykietka="Szukaj w oknie"]',
];

function zdejmijDrogiBezPokrycia(): void {
  for (const wybor of DROGI_BEZ_POKRYCIA) {
    for (const wezel of document.querySelectorAll(wybor)) wezel.remove();
  }
  for (const rzad of document.querySelectorAll('.cd-rzad-czynnosci')) {
    if (rzad.querySelector('button') === null) (rzad.closest('.cd-strefa') ?? rzad).remove();
  }
  for (const menu of document.querySelectorAll('.cd-sekcja-menu')) {
    if (menu.querySelector('.sta-menu-poz') === null) menu.remove();
  }
}

interface DrogaPlatformowa {
  wybor: string;
  nazwa: string;
  zamiast: string;
}

/* Pięć dróg do okien platformowych. Droga zostaje w oknie i nazywa swoją
   niegotowość — zdjęcie jej zabrałoby ślad, że taki zakres w produkcie jest. */
const OKNA_PLATFORMOWE: readonly DrogaPlatformowa[] = [
  {
    wybor: '[data-otwarz-historie]',
    nazwa: 'Historia sesji',
    zamiast: 'Sesje konta stoją w panelu bocznym okna roboczego.',
  },
  {
    wybor: '#cd-konfiguracja',
    nazwa: 'Konfiguracja',
    zamiast: 'Nastawy rdzenia zmienia dziś rodzina komend config.',
  },
  {
    wybor: '#cd-mobile, [data-skrot-komponent="Mobile"]',
    nazwa: 'Mobile',
    zamiast: 'Operacje rodziny mobile stoją w katalogu „Operacje platformy".',
  },
  {
    wybor: '#cd-aod-przelacz, [data-skrot-komponent="Always On Display"]',
    nazwa: 'Always On Display',
    zamiast: 'Operacje nakładki stoją w katalogu „Operacje platformy".',
  },
  // Pozycje pomocy prowadzą do Instrukcji i Instalatora, więc nazwę niesie podpis.
  { wybor: '[data-nawiguj]', nazwa: '', zamiast: '' },
];

function zapowiedzOknoPlatformowe(droga: DrogaPlatformowa, wezel: HTMLElement): void {
  const podpis = droga.nazwa === '' ? (wezel.textContent ?? '').trim() : droga.nazwa;
  const zdanie = droga.zamiast === '' ? '' : droga.zamiast + ' ';
  oglos(
    podpis === '' ? 'Okno platformowe' : podpis,
    zdanie + 'To okno nie wchodzi do tego wydania.',
  );
}

const KARTY_SRODOWISK = new Map<string, HTMLElement>();

let wzorGalezi: HTMLElement | null = null;

/* Siatka pustoszeje przed pytaniem: przy odmowie ma stać pusta strefa, a nie
   cztery karty z prototypu. Karta bez pokrycia w rejestrze nie wraca. */
async function wypelnijSrodowiska(kanal: Kanal, obszar: HTMLElement): Promise<void> {
  const siatka = obszar.querySelector<HTMLElement>('.cd-siatka-srodowisk');
  if (siatka === null) return;
  for (const karta of siatka.querySelectorAll<HTMLElement>('.dn-karta-srodowiska')) {
    const kod = (karta.dataset.srodowisko ?? '').toLowerCase();
    if (kod !== '') KARTY_SRODOWISK.set(kod, karta);
  }
  siatka.replaceChildren();
  const wynik = await wywolaj(kanal, Command.EnvironmentList, { includeModules: true });
  if (!wynik.udany || wynik.wynik === undefined) {
    oglos('Środowiska pracy', wynik.blad?.message ?? 'Rdzeń odmówił wykazu środowisk.', 'blad');
    return;
  }
  let pierwsza = true;
  const uporzadkowane = [...wynik.wynik.environments].sort((a, b) => a.order - b.order);
  for (const srodowisko of uporzadkowane) {
    const karta = KARTY_SRODOWISK.get(srodowisko.code.toLowerCase());
    if (karta === undefined) continue;
    karta.dataset.srodowisko = srodowisko.code;
    opiszSrodowisko(karta, srodowisko);
    karta.tabIndex = pierwsza ? 0 : -1;
    pierwsza = false;
    siatka.appendChild(karta);
  }
}

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

function opiszStopke(karta: HTMLElement, sesji: number): void {
  const stopka = karta.querySelector('.cd-karta-meta');
  if (stopka === null) return;
  for (const wezel of stopka.childNodes) {
    if (wezel.nodeType !== Node.TEXT_NODE) continue;
    if ((wezel.nodeValue ?? '').trim() === '') continue;
    wezel.nodeValue = miaraSesji(sesji);
    return;
  }
}

