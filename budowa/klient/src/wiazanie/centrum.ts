/**
 * Wiązanie Centrum dowodzenia z rdzeniem. Znacznik niesie biblioteka
 * Właściciela — ten plik nic nie buduje: wypełnia wykaz sesji odpowiedzią
 * rdzenia, zakłada sesje na żądanie i wprowadza w okno modułu.
 */

import {
  ChangeKind,
  Command,
  ErrorCode,
  EventType,
  ProgressStatus,
  type Component,
  type Environment,
  type Module,
  type Session,
} from '../../../shared/contract.ts';
import type { Kanal, Wynik } from '../protokol/kanal.ts';
import { wywolaj } from '../protokol/wywolanie.ts';
import { zwiazWyborModulu } from './wybor-modulu.ts';
import { oglos } from './ogloszenie.ts';
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
  otworzOkno,
  przelaczOkno,
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
  wskazSrodowisko,
} from './sesja-biezaca.ts';
import { zwiazStudio, zwolnijStudio } from './studio.ts';
import { zglosUchwyt } from './zdarzenia.ts';

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
  /** Karta modułu ze znacznika; stoi zasłonięta jako wzór, z którego powstaje wnętrze każdej karty. */
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
  zdejmijDrogiBezPokrycia();
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
  /* Okna robocze odtwarzają się z odpowiedzi rdzenia, więc przełącznik okien
     i pasmo kart przerysowują się po niej, nie przed nią. */
  void przejmijOgnisko(kanal).then(() => {
    odswiezOknaRobocze();
    odswiezPasmo();
  });
  const katalogSrodowisk = new Map<string, Environment>();
  void wczytajSrodowiska(kanal, katalogSrodowisk);

  /* Drogi do okien platformowych nazywają swoją niegotowość, zanim dojdą do
     biblioteki: `rama.js` prowadzi je do pliku prototypu, a prototyp poza
     zestawem okien nie stoi, więc naciśnięcie kończyłoby się ciszą. */
  document.addEventListener('click', (zdarzenie) => {
    const cel = zdarzenie.target;
    if (!(cel instanceof Element)) return;
    for (const droga of OKNA_PLATFORMOWE) {
      const wezel = cel.closest<HTMLElement>(droga.wybor);
      if (wezel === null) continue;
      zdarzenie.stopPropagation();
      zdarzenie.preventDefault();
      zapowiedzOknoPlatformowe(droga, wezel);
      return;
    }
  }, true);

  /* „Nowa sesja" otwiera okno robocze na karcie Centrum. Sesja powstaje przy
     pierwszej wiadomości, więc naciśnięcie przycisku nie ma czego założyć
     w rejestrze rdzenia. Nasłuch na dokumencie: menu i szyna stoją poza
     obszarem okna. */
  document.addEventListener('click', (zdarzenie) => {
    const cel = zdarzenie.target;
    if (!(cel instanceof Element)) return;
    /* Przycisk panelu bocznego niesie cechę bez wartości; wartość niosą pozycje
       menu, a te prowadzą do okna nowego projektu, nie do nowej sesji. */
    if (cel.closest<HTMLElement>('[data-okno-nowe]')?.dataset.oknoNowe !== '') return;
    otworzOkno();
    pokazOknoRobocze();
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
      // Czynność porzucona przez Operatora nie jest niepowodzeniem rdzenia
      // i wykazu nie rusza.
      if (wynik === null) return;
      // Odmowa rdzenia wychodzi na wierzch: czynność, która milczy po
      // niepowodzeniu, zostawia Operatora przy wykazie sprzed czynności bez
      // słowa, dlaczego się nie zmienił.
      if (!wynik.udany) oglos('Czynność sesji', wynik.blad?.message ?? 'Rdzeń odmówił wykonania.');
      odswiez();
    });
  });

  /* Czynności komponentu własnego. Wykaz wraca odpowiedzią rdzenia, bo kopia
     i usunięcie zmieniają go po stronie magazynu, nie w znaczniku. */
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

  /* Zakładanie komponentu własnego wymaga rodzaju i nazwy; okna, w którym
     Operator je poda, to wydanie nie niesie, więc przycisk nazywa niegotowość
     zamiast zakładać komponent nazwany za niego. */
  document.addEventListener('click', (zdarzenie) => {
    const cel = zdarzenie.target;
    if (!(cel instanceof Element) || cel.closest('#cd-dodaj-komponent') === null) return;
    zdarzenie.stopPropagation();
    oglos('Nowy komponent', 'Okno zakładania komponentu własnego nie wchodzi do tego wydania.');
  }, true);

  /* Wnętrze okna zmienia się w miejscu, więc kolejne wejście podmienia element
     wstawiony poprzednio, nie ten zdjęty przy pierwszym. Szyna stoi poza
     wnętrzem, więc nasłuch obejmuje cały dokument.

     Faza przechwytywania jest konieczna: grot wejścia biblioteki zatrzymuje
     zdarzenie na sobie, więc w fazie bąbelkowania nasłuch nigdy by go nie
     zobaczył. */
  /* Nazwa środowiska przechodzi z Centrum przez przedsionek aż do głowy karty
     modułu; kafel przedsionka nie stoi w szynie, więc sam jej nie niesie. */
  let nazwaSrodowiska = '';
  /* Identyfikator karty pokazywanej w płótnie; pustka znaczy kartę główną. */
  let kartaBiezaca = '';

  /**
   * Stawia kartę okna roboczego: wnętrze, głowa, pasmo i — gdy wnętrze dopiero
   * stanęło — wiązanie z rdzeniem. Wnętrze stojące wraca bez wiązania: jego
   * nasłuchy już stoją, a drugie wiązanie dawałoby podwójne wysyłki.
   */
  const postawKarte = (karta: KartaRobocza, modul: Module): void => {
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
      } else if (idOkna === '') {
        zwiazOkno(kanal, modul.code, nazwaSrodowiska, wnetrze.wezel);
      } else {
        zwiazOknoStojace(kanal, modul.code, nazwaSrodowiska, idOkna, wnetrze.wezel);
      }
    }
    /* Karta nazywa się pracą, którą niesie jej wnętrze — moduł stoi przy niej
       cechą. Nazwę podaje wnętrze po zamontowaniu, więc czyta się ją po nim. */
    const nazwaPracy = nazwaPracyKarty(wnetrze.wezel);
    if (nazwaPracy !== '') {
      karta.nazwa = nazwaPracy;
      nazwijKarte(karta.id, nazwaPracy);
    }
  };

  /** Otwiera moduł kartą okna roboczego: stojącą dla wskazanego okna komunikacji albo nową. */
  const otworzModul = (modul: Module, idOkna: string): void => {
    if (szablonModulu('dn-tresc-' + modul.code) === null) {
      zapowiedzModul(modul);
      return;
    }
    const karta = kartaOkna(zapiszKarte(modul.code, modul.name, idOkna));
    if (karta !== undefined) postawKarte(karta, modul);
  };

  /** Wraca na kartę stojącą w oknie bieżącym; karta modułu spoza rejestru nie ma czego postawić. */
  const pokazKarte = (idKarty: string): void => {
    const karta = wskazKarte(idKarty);
    const modul = karta === undefined ? undefined : katalogModulow.get(karta.kodModulu);
    if (karta === undefined || modul === undefined) return;
    postawKarte(karta, modul);
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
    ustawWskazowke(pozycja.getAttribute('aria-checked') !== 'true');
  }, true);

  /* Zamknięcie wskazówki startowej znakiem przy niej samej. */
  document.addEventListener('click', (zdarzenie) => {
    const cel = zdarzenie.target;
    if (!(cel instanceof Element) || cel.closest('#cd-start-zamknij') === null) return;
    zdarzenie.stopPropagation();
    ustawWskazowke(false);
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
      oznaczStanPanelu(panel);
      return;
    }
    if (cel.closest('[data-zwin-samouczek]') !== null) {
      zdarzenie.stopPropagation();
      panel.hidden = true;
      oznaczStanPanelu(panel);
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

  /** Przerysowuje pasmo kart z rejestru okna bieżącego; zaznaczona jest karta pokazywana w płótnie. */
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

  /** Zdejmuje z płótna wnętrza kart, których rejestr okien roboczych już nie zna, wraz z ich wiązaniem: uchwytami zdarzeń i nasłuchami. */
  const uprzatnijWnetrza = (): void => {
    for (const wnetrze of wezly.plotno.querySelectorAll<HTMLElement>('.cd-tresc--modul[data-karta]')) {
      const idKarty = wnetrze.dataset.karta ?? '';
      if (kartaOkna(idKarty) !== undefined) continue;
      zwolnijStudio(idKarty);
      wnetrze.remove();
    }
  };

  /**
   * Zamyka kartę: zdejmuje ją z okna roboczego — tam idzie `window.close` —
   * a potem z pasma i z płótna. Karta bieżąca wraca najpierw na Centrum.
   */
  const zamknijKarte = (idKarty: string): void => {
    if (kartaBiezaca === idKarty) wrocDoCentrum();
    zdejmijKarteOkna(idKarty);
    zdejmijKarte(idKarty);
    uprzatnijWnetrza();
    odswiezPasmo();
  };

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

  /** Stawia okno robocze na jego karcie bieżącej: Centrum albo karta modułu. */
  const pokazOknoRobocze = (): void => {
    uprzatnijWnetrza();
    const okno = oknoBiezace();
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
    /* Znak zamknięcia zdejmuje kartę z okna roboczego, z pasma i z płótna:
       karta zdjęta z samego pasma wracała przy najbliższym przełączeniu okna. */
    (idKarty) => zamknijKarte(idKarty),
  );

  /* Zdarzenia rdzenia warstwy wspólnej: zmiana zrobiona w innym oknie albo przez
     proces w tle dochodzi do Centrum wyłącznie nimi. */
  zglosUchwyt(EventType.SessionChanged, (tresc) => {
    odswiez();
    void wypelnijProjekty(kanal);
    // Miara sesji na karcie środowiska liczy sesje, nie ich zmiany.
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
    // Postęp w biegu nie zmienia wykazu sesji; zmienia go domknięcie procesu.
    if (tresc.status === ProgressStatus.Running || tresc.status === ProgressStatus.Pending) return;
    odswiez();
  });
  zglosUchwyt(EventType.DeviceChanged, (tresc) => {
    const wlasne = tresc.devices.find((urzadzenie) => urzadzenie.current);
    if (wlasne !== undefined && !wlasne.hasToken) {
      oglos('Urządzenie', 'Dostęp tego urządzenia do konta został unieważniony.', 'ostrzezenie');
    }
  });

  /* Czynności projektu z menu gałęzi drzewa. Nasłuch na dokumencie, bo
     biblioteka menu wynosi treść menu poza panel. */
  document.addEventListener('click', (zdarzenie) => {
    const cel = zdarzenie.target;
    if (!(cel instanceof Element)) return;
    const pozycja = cel.closest<HTMLElement>('[data-projekt-akcja]');
    const idProjektu = pozycja?.dataset.idProjektu;
    if (pozycja === null || idProjektu === undefined) return;
    void wykonajCzynnoscProjektu(kanal, pozycja.dataset.projektAkcja ?? '', idProjektu)
      .then((wynik) => {
        if (wynik === null) return;
        if (!wynik.udany) {
          oglos('Czynność projektu', wynik.blad?.message ?? 'Rdzeń odmówił wykonania.', 'blad');
        }
        void wypelnijProjekty(kanal);
        // Usunięcie projektu zdejmuje przypisanie z jego sesji.
        odswiez();
      });
  });

  /* „Nowy projekt" z menu szyny: nazwę Operator wpisuje w gałęzi drzewa
     projektów, bo okna zakładania projektu to wydanie nie niesie. Środowisko
     z pozycji menu nie ma pola w `project.create` i nie wchodzi w zadanie. */
  document.addEventListener('click', (zdarzenie) => {
    const cel = zdarzenie.target;
    if (!(cel instanceof Element) || cel.closest('[data-nowy-projekt-srodowisko]') === null) return;
    void zalozProjekt(kanal).then((wynik) => {
      if (wynik === null) return;
      if (!wynik.udany) {
        oglos('Nowy projekt', wynik.blad?.message ?? 'Rdzeń odmówił założenia projektu.', 'blad');
      }
      void wypelnijProjekty(kanal);
    });
  });

  /** Wprowadza w przedsionek środowiska; środowisko bez modułów w rejestrze nie ma czego pokazać i mówi to wprost. */
  const wejdzWPrzedsionek = (kodSrodowiska: string, nazwaZKarty: string): void => {
    if (kodSrodowiska === '') return;
    const srodowisko = katalogSrodowisk.get(kodSrodowiska);
    if (srodowisko !== undefined && (srodowisko.moduleCodes?.length ?? 0) === 0) {
      oglos(srodowisko.name, 'Rejestr rdzenia nie wskazuje dla tego środowiska '
        + 'ani jednego modułu, więc przedsionek nie ma czego pokazać.');
      return;
    }
    const widok = widokPrzedsionka(wezly, gniazdoPrzedsionka(kodSrodowiska));
    if (widok === null) return;
    pokazWidok(wezly, widok);
    const nazwa = srodowisko?.name ?? nazwaZKarty;
    if (nazwa !== '') nazwaSrodowiska = nazwa;
    wskazSrodowisko(kodSrodowiska);
    void zwiazWyborModulu(kanal, kodSrodowiska);
  };

  /* Karty środowisk chodzą strzałkami, a Enter i spacja wprowadzają
     w środowisko: karta jest pozycją listy, nie przyciskiem, więc sama
     klawiatury nie obsłuży. */
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
    wejdzWPrzedsionek(karta.dataset.srodowisko ?? '', nazwaKartySrodowiska(karta));
  });

  document.addEventListener('click', (zdarzenie) => {
    const cel = zdarzenie.target;
    if (!(cel instanceof Element)) return;

    // Wpis w miejscu trwa: kliknięcie w polu wpisu prowadzi kursor, nie wskazuje.
    if (cel.closest('[contenteditable]') !== null) return;

    // Powrót na kartę główną okna roboczego.
    if (cel.closest(POWROT_NA_STRONE_GLOWNA) !== null) {
      wrocDoCentrum();
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

    /* „Nowa sesja" wskazana środowiskiem otwiera okno robocze i wprowadza
       w przedsionek; sesji nie zakłada, bo ta powstaje przy pierwszej
       wiadomości. */
    const kodNowejSesji = cel.closest<HTMLElement>('[data-nowa-sesja-srodowisko]')?.dataset
      .nowaSesjaSrodowisko;
    if (kodNowejSesji !== undefined) {
      zdarzenie.stopPropagation();
      otworzOkno();
      odswiezOknaRobocze();
      wejdzWPrzedsionek(kodNowejSesji, '');
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
      zdarzenie.stopPropagation();
      wejdzWPrzedsionek(kodSrodowiska, nazwaKartySrodowiska(cel));
      return;
    }

    /* Komponent własny otwiera się kartą modułu swojego rodzaju: komponent
       automatyki wchodzi w Automations, ekspert w Agents. */
    const komponent = cel.closest<HTMLElement>('[data-otworz-komponent]')?.dataset.otworzKomponent;
    if (komponent !== undefined) {
      // Rodzaj komponentu jest kodem modułu, w którym komponent stoi.
      const modulKomponentu = katalogModulow.get(KOMPONENTY.get(komponent)?.kind ?? '');
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
    // Skrót szybkiego wyboru nazywa rodzaj komponentu, a rodzaj jest kodem modułu.
    cel.closest('.dn-szyna-poz--skrot')?.getAttribute('data-skrot-komponent') ??
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

/**
 * Kod środowiska karty Centrum, w którą Operator nacisnął. Wejście niesie cała
 * karta, nie sam jej grot: grot jest znakiem wejścia, a nie jedynym miejscem,
 * w które da się trafić.
 */
function kodSrodowiskaWejscia(cel: Element): string {
  return cel.closest<HTMLElement>('.dn-karta-srodowiska')?.dataset.srodowisko ?? '';
}

/** Przestawia ognisko na sąsiednią kartę środowiska; poza nią karty stoją w porządku tabulacji jako jedna pozycja. */
function przestawOgnisko(karta: HTMLElement, krok: number): void {
  const karty = [...document.querySelectorAll<HTMLElement>('.dn-karta-srodowiska')];
  if (karty.length === 0) return;
  const numer = karty.indexOf(karta);
  const nastepna = karty[(numer + krok + karty.length) % karty.length];
  for (const kandydat of karty) kandydat.tabIndex = -1;
  nastepna.tabIndex = 0;
  nastepna.focus();
}

/** Stawia albo chowa wskazówkę startową wraz ze znakiem wyboru przy pozycjach menu widoku. */
function ustawWskazowke(widoczna: boolean): void {
  const wskazowka = document.querySelector<HTMLElement>('.cd-wskazowka, #cd-start');
  if (wskazowka !== null) wskazowka.hidden = !widoczna;
  for (const pozycja of document.querySelectorAll('[data-cd-wskazowka]')) {
    pozycja.setAttribute('aria-checked', String(widoczna));
  }
}

/** Uzgadnia znak wyboru pozycji menu ze stanem panelu samouczka; pozycja menu jest przełącznikiem, więc musi mówić, w którym stanie panel stoi. */
function oznaczStanPanelu(panel: HTMLElement): void {
  for (const pozycja of document.querySelectorAll('[data-przelacz-samouczek][role="menuitemcheckbox"]')) {
    pozycja.setAttribute('aria-checked', String(!panel.hidden));
  }
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

/** Blok wnętrza modułu z szablonu; brak znaczy moduł, którego okna to wydanie nie niesie. */
function szablonModulu(gniazdo: string): Element | null {
  const szablon = document.getElementById(gniazdo);
  if (!(szablon instanceof HTMLTemplateElement)) return null;
  return szablon.content.firstElementChild;
}

/** Wnętrze karty stojące w płótnie; brak znaczy kartę jeszcze niepostawioną. */
function wnetrzeStojace(wezly: WezlyCentrum, idKarty: string): HTMLElement | null {
  return wezly.plotno.querySelector<HTMLElement>(`.cd-tresc--modul[data-karta="${idKarty}"]`);
}

/**
 * Wnętrze karty w płótnie: stojące wraca, brakujące powstaje z klonu karty
 * modułu ze znacznika i bloku szablonu. Każda karta ma własny węzeł — dwie
 * karty tego samego modułu nie dzielą wnętrza, więc przełączenie nie czyści
 * rozmowy. Wnętrze nowe staje na czele płótna, bo wiązanie Studia bierze
 * pierwszy węzeł okna w dokumencie.
 */
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
  wezly.plotno.prepend(wezel);
  return { wezel, nowe: true };
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

/** Nazwa pracy, którą niesie wnętrze karty; pustka znaczy wnętrze bez nazwanej pracy. */
function nazwaPracyKarty(wnetrze: Element): string {
  const znacznik = wnetrze.querySelector('.sta-okno-znacznik, .st-wstazka-sesja span');
  return znacznik?.textContent?.trim() ?? '';
}

/** Wpisuje nazwę modułu w pas stanu ramy; karta główna zostawia pole puste. */
function opiszPasekModulu(nazwa: string): void {
  const pole = document.querySelector('[data-pasek-modul-nazwa]');
  if (pole !== null) pole.textContent = nazwa;
}

/** Opisuje głowę wnętrza karty nazwą modułu i nazwą karty sesji. */
function opiszGloweKarty(wnetrze: HTMLElement, nazwaModulu: string, nazwaSesji: string): void {
  const nazwa = wnetrze.querySelector('[data-karta-modul-nazwa]');
  if (nazwa !== null) nazwa.textContent = nazwaModulu;
  const meta = wnetrze.querySelector('.cd-modul-glowa .dn-meta');
  if (meta !== null) meta.textContent = nazwaSesji === '' ? '' : 'sesja: ' + nazwaSesji;
}

/** Komponenty własne po identyfikatorze; kopia bierze z nich definicję, bo kontrakt nie zna kopiowania komponentu. */
const KOMPONENTY = new Map<string, Component>();

/** Wzór pozycji komponentu zdjęty z treści przykładowej; wykaz pustoszeje przed pytaniem rdzenia, więc wzór musi przeżyć pierwsze wypełnienie. */
let wzorKomponentu: HTMLElement | null = null;

/**
 * Wypełnia wykaz komponentów własnych Centrum rejestrem rdzenia. Wykaz pustoszeje
 * przed pytaniem rdzenia: przy odmowie na ekranie ma stać stan pusty ze znacznika,
 * nie komponenty z prototypu.
 */
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

/**
 * Zwraca klon wzoru opisany komponentem rdzenia. Menu czynności dostaje własne
 * odwołanie i identyfikator komponentu przy każdej pozycji, bo biblioteka menu
 * przenosi jego treść poza kafel.
 */
function zbudujKomponent(wzor: HTMLElement, komponent: Component): HTMLElement {
  const pozycja = wzor.cloneNode(true) as HTMLElement;
  for (const otworz of pozycja.querySelectorAll<HTMLElement>('[data-otworz-komponent]')) {
    otworz.dataset.otworzKomponent = komponent.id;
  }
  const nazwa = pozycja.querySelector('.dn-kafel-nazwa');
  if (nazwa !== null) nazwa.textContent = komponent.name;
  /* Wiersz pod nazwą niesie w prototypie rodzaj i skrót konfiguracji; kontrakt
     daje dla niego wyłącznie opis, a bez opisu wiersz znika. */
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

/* Czynności komponentu z pokryciem w kontrakcie. Kopia zakłada komponent od nowa
   z definicji stojącej, bo kontrakt nie zna kopiowania komponentu; wydania
   i wczytania definicji z pliku nie zna wcale. */
const OPERACJE_KOMPONENTU = new Set(['duplikuj', 'usun']);

/** Wykonuje czynność komponentu komendą kontraktu; komponent nieznany Centrum wraca nazwaną odmową, nie cichym niepowodzeniem. */
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

/** Zdejmuje z menu komponentu pozycje bez komendy w kontrakcie; „Otwórz" zostaje, bo prowadzi do okna modułu. */
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

/** Zdejmuje wzór z treści przykładowej przy pierwszym wypełnieniu; przy kolejnych oddaje zapamiętany, bo wykaz pustoszeje przed pytaniem rdzenia. */
function zapamietajWzor(
  zapamietany: HTMLElement | null,
  wykaz: HTMLElement,
  wybor: string,
): HTMLElement | null {
  if (zapamietany !== null) return zapamietany;
  const wezel = wykaz.querySelector<HTMLElement>(wybor);
  return wezel === null ? null : (wezel.cloneNode(true) as HTMLElement);
}

/**
 * Wypełnia drzewo projektów lewego panelu wykazem rdzenia wraz z sesjami
 * każdego projektu. Projekt bez ani jednej sesji też stoi — lewy panel jest
 * jedynym miejscem, w którym widać całość dorobku Operatora.
 */
async function wypelnijProjekty(kanal: Kanal): Promise<void> {
  const drzewo = document.querySelector<HTMLElement>('#panel-projekty .dn-panel-drzewo');
  if (drzewo === null) return;
  wzorGalezi = zapamietajWzor(wzorGalezi, drzewo, '.dn-panel-galaz');
  const wzor = wzorGalezi;
  const wzorSesji = wzor?.querySelector<HTMLElement>('.dn-panel-wiersz') ?? null;
  // Drzewo pustoszeje przed pytaniem rdzenia: przy odmowie ma stać puste,
  // a nie projektami i sesjami z prototypu.
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
    /* Wiersz sesji w drzewie ma ten sam kształt co w wykazie sesji, z menu
       czynności włącznie; przedrostek odwołania menu jest inny, bo ta sama
       sesja stoi w obu wykazach i dwa menu nie mogą dzielić jednego `id`. */
    for (const sesja of wedlugProjektu.get(projekt.id) ?? []) {
      if (wzorSesji === null) break;
      wpis.appendChild(zbudujWiersz(wzorSesji, sesja, 'menu-sesji-projektu-'));
    }
    drzewo.appendChild(wpis);
  }
}

/* Czynności menu gałęzi projektu z pokryciem w kontrakcie, po podpisie pozycji
   prototypu: znacznik nie niesie dla nich własnego uchwytu. Pozostałe pozycje
   schodzą — pozycja bez komendy jest obietnicą bez pokrycia. */
const CZYNNOSCI_PROJEKTU: Record<string, string> = {
  'Zmień nazwę': 'nazwa',
  'Usuń trwale': 'usun',
};

/**
 * Wiąże menu gałęzi z projektem: własne odwołanie menu, czynności z pokryciem
 * i identyfikator projektu przy każdej pozycji, bo biblioteka menu wynosi
 * treść menu poza gałąź. Menu gałęzi stoi wprost pod nią, obok wierszy sesji.
 */
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

/** Wykonuje czynność projektu komendą kontraktu; `null` znaczy czynność porzuconą przez Operatora. */
async function wykonajCzynnoscProjektu(
  kanal: Kanal,
  czynnosc: string,
  idProjektu: string,
): Promise<Wynik<unknown> | null> {
  // Potwierdzenie nieodwracalności niesie sama pozycja menu: nazywa usunięcie trwałym.
  if (czynnosc === 'usun') return wywolaj(kanal, Command.ProjectDelete, { projectId: idProjektu });
  if (czynnosc !== 'nazwa') return null;
  const wezel = wezelNazwyProjektu(idProjektu);
  const nazwa = wezel === null ? null : await zapytajWWezle(wezel, null);
  if (nazwa === null || nazwa === '') return null;
  return wywolaj(kanal, Command.ProjectRename, { projectId: idProjektu, name: nazwa });
}

/** Węzeł nazwy w gałęzi wskazanego projektu; gałąź szuka się w drzewie, nie w przodkach pozycji menu. */
function wezelNazwyProjektu(idProjektu: string): HTMLElement | null {
  for (const galaz of document.querySelectorAll<HTMLElement>('#panel-projekty .dn-panel-galaz')) {
    if (galaz.dataset.projekt === idProjektu) {
      return galaz.querySelector<HTMLElement>('.dn-panel-galaz-nazwa');
    }
  }
  return null;
}

/**
 * Zakłada projekt pod nazwą wpisaną w nowej gałęzi drzewa. Gałąź wpisu jest
 * klonem wzoru bez sesji i menu; schodzi po wpisie, bo drzewo wraca
 * odpowiedzią rdzenia. Wpis pusty i porzucenie zostawiają drzewo bez zmiany.
 */
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
     na wyzwalacz zabrałby ognisko polu wpisu i porzucił wpis. */
  await new Promise((gotowe) => setTimeout(gotowe, 0));
  pokazPanelProjektow();
  drzewo.prepend(galaz);
  const wpis = await zapytajWWezle(nazwa, '');
  galaz.remove();
  if (wpis === null || wpis === '') return null;
  return wywolaj(kanal, Command.ProjectCreate, { name: wpis });
}

/** Odsłania zakładkę projektów panelu bocznego; pole wpisu nazwy przyjmuje ognisko tylko widoczne. */
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

/** Zdejmuje wzór wiersza z treści przykładowej; kształt wiersza bierze się ze znacznika, nie z kodu. */
function zdejmijWzorWiersza(wykaz: HTMLElement): HTMLElement | null {
  const wiersz = wykaz.querySelector<HTMLElement>('.dn-panel-wiersz');
  return wiersz === null ? null : (wiersz.cloneNode(true) as HTMLElement);
}

/**
 * Wczytuje wykaz sesji rdzenia i wstawia go w miejsce treści przykładowej.
 * Wykaz pustoszeje przed pytaniem: przy odmowie panel ma stać pusty, a nie
 * sesjami, których nie ma.
 */
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
 * Przedrostek odwołania rozróżnia wykazy, w których stoi ta sama sesja.
 */
function zbudujWiersz(wzor: HTMLElement, sesja: Session, przedrostek = 'menu-sesji-'): HTMLElement {
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
    const oznaczenie = przedrostek + sesja.id;
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
const CZYNNOSCI_SESJI: Record<
  string,
  (kanal: Kanal, idSesji: string) => Promise<Wynik<unknown> | null>
> = {
  nazwa: (kanal, idSesji) => zmienNazweSesji(kanal, idSesji),
  przenies: (kanal, idSesji) => przeniesSesjeDoProjektu(kanal, idSesji),
  archiwizuj: (kanal, idSesji) =>
    wywolaj(kanal, Command.SessionArchive, { sessionIds: [idSesji] }),
  // Potwierdzenie nieodwracalności niesie sama pozycja menu: nazywa usunięcie
  // trwałym, a rejestr sesji drugiego pytania nie stawia.
  usun: (kanal, idSesji) =>
    wywolaj(kanal, Command.SessionDelete, { sessionIds: [idSesji], confirm: true }),
  wyjmij: (kanal, idSesji) =>
    wywolaj(kanal, Command.SessionProjectClear, { sessionIds: [idSesji] }),
};

/**
 * Zmienia nazwę sesji wpisaną w samym wierszu wykazu. Wpis pusty i porzucenie
 * zostawiają nazwę bez zmiany; pustka nie jest nazwą, a rdzeń przyjąłby ją.
 */
async function zmienNazweSesji(kanal: Kanal, idSesji: string): Promise<Wynik<unknown> | null> {
  const nazwa = await zapytajWWierszu(idSesji, null);
  if (nazwa === null || nazwa === '') return null;
  return wywolaj(kanal, Command.SessionRename, { sessionId: idSesji, title: nazwa });
}

/**
 * Przenosi sesję do projektu o wpisanej nazwie. Kontrakt zakłada projekt nowy,
 * gdy nie wskazano istniejącego — i jest to jedyna droga powstania projektu,
 * bo rodzina `project.*` zakładania nie niesie.
 */
async function przeniesSesjeDoProjektu(
  kanal: Kanal,
  idSesji: string,
): Promise<Wynik<unknown> | null> {
  const nazwa = await zapytajWWierszu(idSesji, '');
  if (nazwa === null || nazwa === '') return null;
  return wywolaj(kanal, Command.SessionProjectSet, { sessionIds: [idSesji], projectName: nazwa });
}

/**
 * Pyta o tekst w samym wierszu wykazu: węzeł nazwy staje się polem wpisu.
 * Wartość pusta zaczyna od czystego pola, pusty wskaźnik zostawia w polu nazwę
 * stojącą. Enter zatwierdza, Escape i wyjście ogniska porzucają; wiersz wraca do
 * treści sprzed pytania, bo rozstrzyga o niej odpowiedź rdzenia.
 *
 * Osobnego okna nazwy rejestr sesji nie ma, a stawianie własnego wyszłoby poza
 * znacznik Właściciela.
 */
function zapytajWWierszu(idSesji: string, wartosc: string | null): Promise<string | null> {
  const wezel = wezelNazwySesji(idSesji);
  return wezel === null ? Promise.resolve(null) : zapytajWWezle(wezel, wartosc);
}

/** Pyta o tekst we wskazanym węźle nazwy; zasady wpisu jak przy wierszu sesji. */
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
      }
    });
    wezel.addEventListener('blur', () => domknij(null));
    /* Wpis jest zwykłym tekstem: treść wklejona z formatowaniem wniosłaby do
       wiersza znacznik, którego nazwa sesji nie niesie. */
    wezel.setAttribute('contenteditable', 'plaintext-only');
    if (wartosc !== null) wezel.textContent = wartosc;
    wezel.focus();
    zaznaczCalosc(wezel);
  });
}

/**
 * Węzeł nazwy w wierszu wskazanej sesji; biblioteka menu wynosi treść menu poza
 * wiersz, więc wiersz szuka się w wykazach, nie w przodkach pozycji menu.
 * Sesja stoi w wykazie sesji i w drzewie projektów naraz; pierwszeństwo ma
 * wiersz widoczny, bo pole wpisu w panelu zasłoniętym nie przyjmie ogniska.
 */
function wezelNazwySesji(idSesji: string): HTMLElement | null {
  const wiersze = [...document.querySelectorAll<HTMLElement>(
    '#wykaz-sesji .dn-panel-wiersz, #panel-projekty .dn-panel-wiersz',
  )].filter((wiersz) => wiersz.dataset.idSesji === idSesji);
  const wiersz = wiersze.find((kandydat) => kandydat.offsetParent !== null) ?? wiersze[0];
  return wiersz?.querySelector<HTMLElement>('.dn-obszar-pozycja-nazwa') ?? null;
}

/** Zaznacza całą treść węzła, żeby wpis ją zastąpił, a nie dopisał się do niej. */
function zaznaczCalosc(wezel: HTMLElement): void {
  const zakres = document.createRange();
  zakres.selectNodeContents(wezel);
  const zaznaczenie = window.getSelection();
  zaznaczenie?.removeAllRanges();
  zaznaczenie?.addRange(zakres);
}

/** Zdejmuje z menu wiersza pozycje bez komendy w kontrakcie wraz z rozdzielnikami, które po nich zostały. */
function zdejmijCzynnosciBezZrodla(wiersz: HTMLElement): void {
  for (const pozycja of wiersz.querySelectorAll('[data-menu-tresc] .sta-menu-poz')) {
    const czynnosc = pozycja.getAttribute('data-poz-akcja') ?? '';
    if (CZYNNOSCI_SESJI[czynnosc] === undefined) pozycja.remove();
  }
  zdejmijRozdzielnikiSieroce(wiersz);
}

/** Zdejmuje rozdzielniki, które zostały bez pozycji przed sobą; rozdzielnik na czele menu dzieli je od niczego. */
function zdejmijRozdzielnikiSieroce(obudowa: HTMLElement): void {
  for (const rozdzielnik of obudowa.querySelectorAll('[data-menu-tresc] .sta-menu-sep')) {
    const przed = rozdzielnik.previousElementSibling;
    if (przed === null || przed.classList.contains('sta-menu-sep')) rozdzielnik.remove();
  }
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

/* Drogi Centrum bez pokrycia w rdzeniu. Listwa ustawień przełączała stan funkcji
   globalnych, których nie ma, a wydanie i wczytanie definicji komponentu z pliku
   nie mają komendy w kontrakcie. */
const DROGI_BEZ_POKRYCIA = [
  '#cd-konfiguracja',
  '#cd-mobile',
  '#cd-aod-przelacz',
  '.cd-sekcja-glowa [data-operacja]',
];

/** Zdejmuje drogi bez pokrycia w rdzeniu; przycisk, który zmienia sam swój wygląd, twierdzi nieprawdę o stanie platformy. */
function zdejmijDrogiBezPokrycia(): void {
  for (const wybor of DROGI_BEZ_POKRYCIA) {
    for (const wezel of document.querySelectorAll(wybor)) wezel.remove();
  }
  /* Listwa ustawień została bez ani jednej czynności, a sam jej nagłówek
     zapowiadałby zakres, do którego okno nie prowadzi. */
  for (const rzad of document.querySelectorAll('.cd-rzad-czynnosci')) {
    if (rzad.querySelector('button') === null) (rzad.closest('.cd-strefa') ?? rzad).remove();
  }
  /* Menu sekcji zostało bez pozycji; sam jego znak otwierałby pustkę. */
  for (const menu of document.querySelectorAll('.cd-sekcja-menu')) {
    if (menu.querySelector('.sta-menu-poz') === null) menu.remove();
  }
}

/** Droga Centrum do okna platformowego, którego to wydanie nie niesie. */
interface DrogaPlatformowa {
  /** Wskazanie węzłów drogi w znaczniku powłoki i Centrum. */
  wybor: string;
  /** Nazwa okna; pustka znaczy, że nazwę niesie podpis samej drogi. */
  nazwa: string;
  /** Zdanie o tym, gdzie ta praca stoi dzisiaj; pustka zostawia samą zapowiedź. */
  zamiast: string;
}

/* Pięć dróg Centrum do okien platformowych. Droga zostaje w oknie i nazywa swoją
   niegotowość — zdjęcie jej zabrałoby Operatorowi ślad, że taki zakres
   w produkcie jest. */
const OKNA_PLATFORMOWE: readonly DrogaPlatformowa[] = [
  {
    wybor: '[data-otwarz-historie]',
    nazwa: 'Historia sesji',
    zamiast: 'Sesje konta stoją w panelu bocznym okna roboczego.',
  },
  { wybor: '[data-skrot-komponent="Mobile"]', nazwa: 'Mobile', zamiast: '' },
  { wybor: '[data-skrot-komponent="Always on Display"]', nazwa: 'Always On Display', zamiast: '' },
  // Pozycje pomocy prowadzą do trzech różnych okien zestawu, więc nazwę niesie ich podpis.
  { wybor: '[data-nawiguj]', nazwa: '', zamiast: '' },
];

/** Zapowiada okno platformowe, którego to wydanie nie niesie. */
function zapowiedzOknoPlatformowe(droga: DrogaPlatformowa, wezel: HTMLElement): void {
  const podpis = droga.nazwa === '' ? (wezel.textContent ?? '').trim() : droga.nazwa;
  const zdanie = droga.zamiast === '' ? '' : droga.zamiast + ' ';
  oglos(
    podpis === '' ? 'Okno platformowe' : podpis,
    zdanie + 'To okno nie wchodzi do tego wydania.',
  );
}

/** Karty środowisk zdjęte ze znacznika, po kodzie środowiska; wykaz rdzenia rozstrzyga, które i w jakiej kolejności wracają. */
const KARTY_SRODOWISK = new Map<string, HTMLElement>();

/** Wzór gałęzi projektu zdjęty z treści przykładowej; drzewo pustoszeje przed pytaniem rdzenia, więc wzór musi przeżyć pierwsze wypełnienie. */
let wzorGalezi: HTMLElement | null = null;

/**
 * Stawia karty środowisk w porządku rejestru rdzenia. Siatka pustoszeje przed
 * pytaniem: przy odmowie na ekranie ma stać pusta strefa, a nie cztery karty
 * z prototypu. Karta bez pokrycia w rejestrze nie wraca — znacznik niesie ich
 * cztery, a rejestr rozstrzyga, ile ich jest.
 */
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
  // Kolejność kart jest własnością rejestru, nie kolejnością odpowiedzi.
  const uporzadkowane = [...wynik.wynik.environments].sort((a, b) => a.order - b.order);
  for (const srodowisko of uporzadkowane) {
    const karta = KARTY_SRODOWISK.get(srodowisko.code.toLowerCase());
    if (karta === undefined) continue;
    karta.dataset.srodowisko = srodowisko.code;
    opiszSrodowisko(karta, srodowisko);
    /* Karty chodzą strzałkami, więc w porządku tabulacji stoi jedna z nich;
       reszta wchodzi w ognisko strzałką, wzorem listy jednego przystanku. */
    karta.tabIndex = pierwsza ? 0 : -1;
    pierwsza = false;
    siatka.appendChild(karta);
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
