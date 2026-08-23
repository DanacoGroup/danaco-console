import {
  ClipboardEntryKind,
  DiffHunkKind,
  StudioAuthor,
  StudioOperationScope,
  StudioPageOrientation,
  type StudioAnnotation,
  type StudioComment,
  type StudioDiffHunk,
  type StudioTextMatch,
  type StudioTrackedChange,
} from '../../../../shared/contract';
import { opisOdmowy } from '../../komponenty/odmowa';
import { poleWyboru, utworzWierszOdpowiedzi } from '../../modele/kontrolki-formularza';
import {
  wczytajDokument,
  wstawWMiejsceKursora,
  zapiszDokument,
  type ZapleczeEdytora,
} from './czynnosci-edytora';
import {
  eksportujDokument,
  przekazDoLibrary,
  type ZapleczePodgladu,
} from './czynnosci-podgladu';
import {
  dodajKomentarz,
  odczytajZmiany,
  przestawSledzenie,
  rozstrzygnijPropozycje,
  rozstrzygnijZmiany,
  zlecOperacje,
  type ZapleczePracy,
} from './czynnosci-pracy';
import type { AparatZrodlo } from './aparat-zrodlo';
import { utworzAparatPanel } from './aparat-panel';
import { utworzBlokadaPanel } from './blokada-panel';
import { utworzDziennikPanel } from './dziennik-panel';
import { utworzKopiePanel } from './kopie-panel';
import type { ObiektZrodlo } from './obiekt-zrodlo';
import { utworzObiektPanel } from './obiekt-panel';
import { utworzPostacDokumentu } from './strona-postac-dokumentu';
import type { SzablonZrodlo } from './szablon-zrodlo';
import { utworzSzablonPanel } from './szablon-panel';
import type { TabelaZrodlo } from './tabela-zrodlo';
import { utworzTabelaPanel } from './tabela-panel';
import type { ZrodloKontroliStudio } from './zrodlo-kontroli-studio';
import type { ZrodloPostaciStudio } from './zrodlo-postaci-studio';
import type { ZrodloWstawienStudio } from './zrodlo-wstawien-studio';
import { utworzDymkiKomentarzy } from './dymki-komentarzy';
import { utworzGalerieSzablonow } from './galeria-szablonow';
import { domyslnaStrona, nosnikPoOznaczeniu, stronaZProfilu } from './nastawy-strony';
import { domyslneNastawy, opiszNastawy, utworzNastawyAkapitow } from './nastawy-wizualne';
import { utworzOknoStudio } from './okno-studio';
import { utworzPanelRedaktora } from './panel-redaktora';
import { utworzPasekStatusu } from './pasek-statusu';
import { utworzPolaRoznicy } from './pola-roznicy';
import {
  utworzPowierzchnieDokumentu,
  type KartkaRdzenia,
  type TrybWidoku,
} from './powierzchnia-dokumentu';
import { utworzPowiazanieRozmowy } from './powiazanie-rozmowy';
import { wyrysFragmentow, wyrysTrafien } from './roznica-i-trafienia';
import type { StanStudio } from './stan-studio';
import { opiszStatystyke, policzRoznice, przefiltrujRoznice, type FiltrRoznicy } from './statystyka-roznicy';
import { neutralneNastawy, type NastawySuwakow } from './suwaki-koncepcyjne';
import { utworzWczytanieDokumentu, utworzWniesieniePliku } from './wczytanie-dokumentu';
import { utworzWierszPolecenia } from './wiersz-polecenia';
import { utworzWstazkePracy } from './wstazka-pracy';
import { utworzZnajdzZamien } from './znajdz-zamien';
import { zastosujNarzedzie } from './znaczniki-markdown';
import { opiszZmiany } from './zmiany-modelu';
import { ZNACZNIK_PODZIALU } from './zapis-formatowany';
import { utworzWydaniePanel } from './konwersja-dokumentu';
import { utworzKartyDokumentow } from './karty-dokumentow';
import { przybornikZetonBarwy, utworzPrzybornikZnakowania } from './przybornik-znakowania';
import { utworzKatalogOperacji } from './przybornik-katalog';
import { utworzPlywakOperacji, CZYNNOSCI_NA_WIERZCHU } from './przybornik-plywak';
import { utworzPrzyciskMikrofonu } from './przybornik-mowa';
import { utworzZnacznikiWlasne } from './przybornik-znaczniki';
import { przybornikZlozZnakowania } from './przybornik-wykaz';
import {
  TrybOperacji,
  utworzUzyciePrzybornika,
  KLUCZ_PRZYPIETYCH,
  KLUCZ_TRYBU,
  KLUCZ_UZYCIA,
} from './przybornik-uzycie';
import type { ZapleczePrzybornika } from './przybornik-zaplecze';
import { utworzSchowekHistorie } from './schowek-historia';
import { utworzOsadzeniePanel, OSADZENIE_ZNAKOW_PODGLADU } from './osadzenie-panel';
import {
  ZrodloWniesienia,
  osadzenieWierszPochodzenia,
  utworzWykazPochodzen,
} from './osadzenie-pochodzenie';
import type { ZrodloAkcjiStudio } from './zrodlo-akcji-studio';
import type { ZrodloDokumentuStudio } from './zrodlo-dokumentu-studio';
import type { ZrodloPracyStudio } from './zrodlo-pracy-studio';
import type { ZrodloPrzekazania } from './zrodlo-przekazania';

/**
 * Okno pracy z dokumentem — jedna powierzchnia zamiast czterech.
 *
 * ── Co się w nim zeszło ─────────────────────────────────────────────────────
 * Studio Editor, kanwa tekstowa, Preview Window i Diff/Grep Panel. Cztery okna
 * pracowały nad tą samą treścią w czterech miejscach, a przy dwóch dokumentach
 * dawało to sześć okien. Tutaj treść jest jedna, tryby widoku są jej trybami,
 * a różnica jest warstwą nakładaną na nią, nie drugą kolumną.
 *
 * Kanwa tekstowa nie ma następnika w postaci osobnej powierzchni i mieć go nie
 * miała: wynik operacji wchodzi do dokumentu jako zmiana śledzona autora
 * `model` (`studio.contextual.op` w `adapter_modul_studio_roznice.go`), a jej
 * poprawianie przed decyzją odbywa się tam, gdzie zmiana stoi — w treści.
 * Poprawka po stronie klienta nie jest więc już potrzebna, bo poprawia się
 * dokument, a nie bufor obok niego.
 *
 * ── Układ ───────────────────────────────────────────────────────────────────
 * Wstążka z zakładkami na górze, zakładki dokumentów pod nią, dalej trzy
 * kolumny: powierzchnia z kartkami, margines z dymkami komentarzy i panel
 * Redaktora. Wiersz polecenia otwiera się przy kursorze, wewnątrz powierzchni.
 *
 * ── Czego okno nie robi ─────────────────────────────────────────────────────
 * Nie prowadzi rozmowy — Chat Window jest jedno na produkt i stoi obok
 * (`powiazanie-rozmowy.ts`). Nie zastępuje Tools Panelu jako pełnego wykazu
 * operacji, Session Repository jako historii wersji ani Ingest/OCR Panelu jako
 * cyfryzacji: te trzy nie mają powierzchni tekstowej i zostają osobno.
 */
export interface OknoPracyZDokumentem {
  element: HTMLElement;
  /** Wczytuje to, co okno potrzebuje z rdzenia: szablony i profile wydania. */
  wczytaj(): Promise<void>;
  /**
   * Przestawia tryb wykazu operacji — narzędzia ukryte albo stały panel.
   *
   * Woła to stały panel operacji, gdy Operator przełączy tryb u niego: nastawa
   * jest jedna na moduł i zapisuje ją okno pracy, bo tu stoi pływak.
   */
  ustawTrybOperacji(tryb: TrybOperacji): void;
  odswiez(): void;
}

/** Czynności, które okno zleca poza siebie — do innych okien modułu. */
export interface CzynnosciPracy {
  /** Przenosi ognisko do Tools Panelu — pełnego wykazu operacji. */
  naWiecejOperacji(): void;
  /** Przenosi ognisko do Warsztatu dokumentu — czynności ochrony. */
  naWarsztat(): void;
  /**
   * Przestawia tryb wykazu operacji w całym module.
   *
   * Nastawa dotyczy dwóch okien naraz — pływaka w tym oknie i stałego panelu
   * obok — więc rozstrzyga ją moduł, a nie żadne z nich osobno.
   */
  naTrybOperacji(tryb: TrybOperacji): void;
}

/**
 * Siedem źródeł postaci dokumentu, kontroli pracy i wstawień.
 *
 * Rodziny `studio.page.*`, `style.*`, `format.*`, `table.*`, `object.*`,
 * `apparatus.*`, `field.*`, `journal.*`, `markup.*`, `lock.*`, `backup.*`
 * i `template.*` nie należą do portu pracy — dlatego wchodzą osobnym bytem,
 * a nie dopiskiem do `ZrodloPracyStudio`. Byt jest **nieobowiązkowy**: okno bez
 * niego działa jak dotąd i mówi wprost, że drogi do tych rodzin nie podano —
 * brak jest wtedy po stronie złożenia modułu, nie rdzenia.
 */
export interface ZrodlaPostaciStudia {
  postaci: ZrodloPostaciStudio;
  kontrola: ZrodloKontroliStudio;
  wstawienia: ZrodloWstawienStudio;
  tabele: TabelaZrodlo;
  obiekty: ObiektZrodlo;
  aparat: AparatZrodlo;
  szablony: SzablonZrodlo;
}

/**
 * Górna granica wielkości JEDNEJ kartki przyjmowanej z rdzenia.
 *
 * Kartka A4 w 96 punktach na cal to obraz PNG rzędu setek kilobajtów; sześć
 * megabajtów mieści ją z zapasem także przy nośniku wielkoformatowym. Granica
 * jest podawana rdzeniowi, bo kontrakt każe mu wtedy ODMÓWIĆ nazywając zmierzoną
 * wielkość, zamiast oddać treść uciętą — a obraz ucięty wyglądałby na kartkę
 * zepsutą przez rdzeń.
 */
const GRANICA_BAJTOW_KARTKI = 6 * 1024 * 1024;

/**
 * Górna liczba kartek pobieranych z rdzenia jednym renderem.
 *
 * Pismo dwustustronicowe to dwieście żądań i dwieście obrazów w pamięci karty.
 * Podgląd wydania służy sprawdzeniu składu, a nie czytaniu całego pisma
 * obrazkami, więc pobiera się początek i mówi o tym wprost.
 */
const GRANICA_KARTEK_RDZENIA = 24;

export function utworzOknoPracyZDokumentem(
  stan: StanStudio,
  praca: ZrodloPracyStudio,
  akcje: ZrodloAkcjiStudio,
  dokumenty: ZrodloDokumentuStudio,
  przekazanie: ZrodloPrzekazania,
  przybornikZaplecze: ZapleczePrzybornika,
  zrodlaPostaci: ZrodlaPostaciStudia,
  czynnosciZewnetrzne: CzynnosciPracy,
): OknoPracyZDokumentem {
  const rama = utworzOknoStudio({
    kod: 'studio.praca-z-dokumentem',
    tytul: 'Praca z dokumentem',
    rola: 'wiodące',
    objasnienie:
      'Jedno okno pracy nad dokumentem: treść z formatowaniem na kartce, podgląd wydania ' +
      'i różnica jako tryby tego samego widoku, zmiany modelu w miejscu. Wczytuje komendą ' +
      'studio.document.open, zapisuje studio.document.save, zleca studio.contextual.op, ' +
      'rozstrzyga studio.tracking.decide i studio.proposal.decide.',
  });

  const odpowiedz = utworzWierszOdpowiedzi();
  const zapleczeEdytora: ZapleczeEdytora = { stan, pas: rama.stan, odpowiedz };
  const zapleczePracy: ZapleczePracy = {
    stan,
    praca,
    pas: rama.stan,
    odpowiedz,
    poZmianie: () => void odswiezZRdzenia(),
  };

  /** Zaplecze czynności podglądu — wydanie i przekazanie do Library. */
  function zapleczePodgladu(): ZapleczePodgladu {
    return { stan, dokumenty, pas: rama.stan, odpowiedz, odswiez: () => odswiez() };
  }

  /* ── Nastawy widoku ──────────────────────────────────────────────────────── */

  let strona = domyslnaStrona();
  let nastawyWizualne = domyslneNastawy();
  const akapity = utworzNastawyAkapitow();
  let nastawySuwakow: NastawySuwakow = neutralneNastawy();
  let zmiany: readonly StudioTrackedChange[] = [];
  let komentarze: readonly StudioComment[] = [];
  let fragmentyOdpowiedzi: readonly StudioDiffHunk[] = [];
  let trafieniaOdpowiedzi: readonly StudioTextMatch[] = [];
  let dokumentOdpowiedzi = '';
  let wskazanyKomentarz = -1;
  let adnotacje: readonly StudioAnnotation[] = [];

  /* ── Znakowanie, użycie, pochodzenie ─────────────────────────────────────── */

  const znaczniki = utworzZnacznikiWlasne();
  const uzycie = utworzUzyciePrzybornika();
  const pochodzenia = utworzWykazPochodzen();
  /** Zakładki powrotu żyjące przez sesję okna — aparatu dokumentu rdzeń nie ma. */
  const zakladki: { kod: string; nazwa: string; naZnaku: number }[] = [];

  /* ── Powierzchnia ────────────────────────────────────────────────────────── */

  /* ── Postać dokumentu: strona, marginesy, numeracja, styl ───────────────── */

  // Postać powstaje PRZED powierzchnią, bo powierzchnia bierze od niej magazyn
  // nastaw widoku i zgłasza jej chwyty linijki. Odwrotna kolejność zostawiłaby
  // chwyty bez odbiorcy — a to jest właśnie ta usterka, po której nastawa ginie
  // przy zapisie.
  const postac = utworzPostacDokumentu({
    zrodlo: zrodlaPostaci.postaci,
    idDokumentu: () => stan.dokument()?.id ?? '',
    idOkna: () => stan.idOkna(),
    zaznaczenie: () => stan.zaznaczenie(),
    miejsceKursora: () => stan.zaznaczenie()?.poczatek ?? stan.trescRobocza().length,
    naPostac: () => powierzchnia.pokaz(stan.trescRobocza(), true),
    naZdanie: (tresc, udane) => odpowiedz.pokaz(tresc, udane),
  });

  const powierzchnia = utworzPowierzchnieDokumentu(strona, nastawyWizualne, akapity, {
    naTresc: (tresc) => {
      stan.ustawTresc(tresc);
      // Kartki wyrysowane przez rdzeń przestają być podglądem TEJ treści w chwili,
      // w której Operator dopisze literę. Zdejmuje się je więc od razu i podgląd
      // wraca do kartek liczonych w oknie — pokazywanie starego wyrysu jako
      // podglądu treści bieżącej byłoby kłamstwem o dokumencie, a nie oszczędnością
      // jednego wywołania.
      if (powierzchnia.kartekZRdzenia() > 0) powierzchnia.ustawKartkiRdzenia([]);
    },
    naZaznaczenie: (zakres) => {
      stan.ustawZaznaczenie(zakres);
      pokazWierszPolecenia(zakres === null ? 0 : zakres.koniec - zakres.poczatek);
    },
    naDecyzjeZmiany: (kod, przyjmij) => void rozstrzygnijZmiany(zapleczePracy, [kod], przyjmij),
    // Chwyty linijki dojeżdżają teraz do rdzenia: margines i wcięcie idą
    // `studio.page.setup.set`, tabulatory `studio.ruler.tabstop.set`. Bez tych
    // trzech przewodów chwyt działał, ale nastawa ginęła przy zapisie.
    naStrone: (nowa) => postac.zglosStrone(nowa),
    naWciecieAkapitu: (_numer, wciecia) => postac.zglosWciecia(wciecia),
    naTabulatoryAkapitu: (_numer, tabulatory) => postac.zglosTabulatory(tabulatory),
  }, postac.magazynWidoku);

  /* ── Narzędzia ukryte: katalog, pływak, mikrofon ─────────────────────────── */

  const katalog = utworzKatalogOperacji({
    naOperacje: (idAkcji) => zlecOperacjeZUzyciem(idAkcji, ''),
    naZapisOperacji: (nazwa, kategoria, polecenie) =>
      void zapiszOperacjeWlasna(nazwa, kategoria, polecenie),
    naUsuniecieOperacji: (idOperacji) => void usunOperacjeWlasna(idOperacji),
  });

  const plywak = utworzPlywakOperacji(
    {
      naOperacje: (idAkcji) => zlecOperacjeZUzyciem(idAkcji, ''),
      naPrzypiecie: (idAkcji) => void przypnijCzynnosc(idAkcji),
      naSuwak: (kod, wartosc) => {
        nastawySuwakow = { ...nastawySuwakow, [kod]: wartosc };
        wstazka.odswiezSuwaki(nastawySuwakow);
      },
      naTryb: (tryb) => ustawTrybOperacji(tryb),
    },
    katalog.element,
  );

  const mikrofonPolecenia = utworzPrzyciskMikrofonu({
    zrodlo: przybornikZaplecze,
    // Moduł dostaje z powłoki okno komunikacji, nie kartę sesji, a kontrakt
    // `speech.audio.upload` przyjmuje puste `sessionId` jako „nagranie bez
    // przypisania do karty". Zgadywanie karty z okna byłoby wskazaniem
    // nieprawdziwym; okno jedzie i wystarcza do odnalezienia nagrania.
    idSesji: () => '',
    idOkna: () => stan.idOkna(),
    naTekst: (tekst) => wiersz.ustawTresc(tekst),
    naZdanie: (tresc, udane) => wiersz.pokazZdanie(tresc, udane),
    etykieta: 'Podyktuj polecenie',
  });

  const wiersz = utworzWierszPolecenia(
    {
      naOperacje: (idAkcji, polecenie) => zlecOperacjeZUzyciem(idAkcji, polecenie),
      naKomentarz: (tresc) => {
        const zakres = stan.zaznaczenie();
        void dodajKomentarz(zapleczePracy, tresc, zakres, '');
      },
      naWiecej: () => czynnosciZewnetrzne.naWiecejOperacji(),
      naPrzypiecie: (idAkcji) => void przypnijCzynnosc(idAkcji),
      naSuwak: (kod, wartosc) => {
        nastawySuwakow = { ...nastawySuwakow, [kod]: wartosc };
        wstazka.odswiezSuwaki(nastawySuwakow);
      },
      naTryb: (tryb) => ustawTrybOperacji(tryb),
    },
    plywak,
    mikrofonPolecenia,
  );

  function pokazWierszPolecenia(dlugosc: number): void {
    przybornik.ustawZaznaczenie(dlugosc);
    const polozenie = powierzchnia.polozenieKursora();
    if (polozenie === null) {
      wiersz.ukryj();
      return;
    }
    wiersz.pokaz(polozenie, dlugosc);
  }

  /**
   * Zleca operację i podnosi jej licznik użycia.
   *
   * Licznik jedzie do rdzenia po każdym uruchomieniu, bo z niego bierze się
   * kolejność czynności na wierzchu pływaka — „z użycia, nie z domysłu".
   */
  function zlecOperacjeZUzyciem(idAkcji: string, polecenie: string): void {
    void zlecOperacje(zapleczePracy, idAkcji, polecenie, nastawySuwakow);
    void przybornikZaplecze.przybornikZapisz(KLUCZ_UZYCIA, uzycie.policz(idAkcji));
    odswiezPlywak();
  }

  /* ── Dymki, redaktor, galeria ────────────────────────────────────────────── */

  const dymki = utworzDymkiKomentarzy({
    naOdpowiedz: (idWatku, tresc) => void dodajKomentarz(zapleczePracy, tresc, null, idWatku),
    naRozwiazanie: (idKomentarza, rozwiazany) =>
      void rozstrzygnijKomentarz(idKomentarza, rozwiazany),
    naKotwice: (idKomentarza) => wskazKomentarz(idKomentarza),
    naDecyzjePropozycji: (przyjmij) => void rozstrzygnijPropozycje(zapleczePracy, przyjmij, []),
    naDecyzjeZmiany: (kod, przyjmij) => void rozstrzygnijZmiany(zapleczePracy, [kod], przyjmij),
  });

  /* ── Przybornik znakowania, schowek, osadzenie źródeł ────────────────────── */

  const mikrofonTresci = utworzPrzyciskMikrofonu({
    zrodlo: przybornikZaplecze,
    idSesji: () => '',
    idOkna: () => stan.idOkna(),
    naTekst: (tekst) => wstawWTresc(tekst),
    naZdanie: (tresc, udane) => przybornik.pokazOdpowiedz(tresc, udane),
    etykieta: 'Dyktuj do treści',
  });

  const przybornik = utworzPrzybornikZnakowania(
    {
      naKomentarz: (tresc) => {
        if (tresc.trim() === '') {
          przybornik.pokazOdpowiedz('Komentarz bez treści nie jest komentarzem.', false);
          return;
        }
        void dodajKomentarz(zapleczePracy, tresc, stan.zaznaczenie(), '');
      },
      naPropozycje: (polecenie) => {
        // Propozycja brzmienia idzie operacją kontekstową: rdzeń oddaje treść
        // wyniku jako propozycję, a ta staje NA MARGINESIE — w treści jej nie ma
        // do chwili przyjęcia (studio.proposal.decide).
        zlecOperacjeZUzyciem('studio.styl.rejestr', polecenie);
      },
      naAdnotacje: (numer, tresc) => void dodajAdnotacjePrzybornika(numer, tresc),
      naZnacznik: (nazwa, barwa) => {
        const znacznik = znaczniki.zaloz(nazwa, barwa, stan.zaznaczenie(), StudioAuthor.Uzytkownik);
        przybornik.pokazOdpowiedz(
          znacznik === null
            ? 'Znacznik bez nazwy nie da się później odnaleźć w wykazie — nazwij go.'
            : `Znacznik „${znacznik.nazwa}" założony. Żyje przez sesję okna: kontrakt nie ma pola ` +
                'na znacznik własny, więc zamknięcie karty go zgubi.',
          znacznik !== null,
        );
        odswiez();
      },
      naOdhaczenie: (kod, rodzaj, zamknij) => {
        if (rodzaj === 'komentarz') {
          void rozstrzygnijKomentarz(kod, zamknij);
          return;
        }
        if (rodzaj === 'znacznik') {
          znaczniki.przestaw(kod, !zamknij);
          odswiez();
          return;
        }
        przybornik.pokazOdpowiedz(
          'Tej pozycji nie odhacza się w wykazie: adnotacji kontrakt nie zamyka, a zmianę śledzoną ' +
            'rozstrzyga decyzja o przyjęciu albo odrzuceniu, nie odhaczenie.',
          false,
        );
      },
      naZdjecieZnacznika: (kod) => {
        znaczniki.zdejmij(kod);
        odswiez();
      },
      naWyroznienie: (barwa) => {
        const numer = powierzchnia.blokKursora();
        if (numer < 0) {
          przybornik.pokazOdpowiedz(
            'Wyróżnienie dotyczy akapitu, w którym stoi kursor — postaw go w treści.',
            false,
          );
          return;
        }
        akapity.ustaw(numer, { barwa: `var(${przybornikZetonBarwy(barwa)})` });
        powierzchnia.pokaz(stan.trescRobocza(), true);
        przybornik.pokazOdpowiedz(
          `Akapit oznaczony barwą „${barwa}" w oknie. Nastawa widoku rozporządza barwą PISMA ` +
            'akapitu, nie tłem fragmentu. Wyróżnienie tła, które PRZEŻYJE zapis, nakłada się ' +
            'w panelu formatowania — postać znaku niesie barwę wyróżnienia i dojeżdża nią do ' +
            'rdzenia; ta droga stoi obok, w panelu „Arkusz stylów i formatowanie".',
          true,
        );
      },
      naZakladke: (nazwa) => {
        const oczyszczona = nazwa.trim();
        if (oczyszczona === '') {
          przybornik.pokazOdpowiedz('Zakładka bez nazwy nie da się wskazać — nazwij ją.', false);
          return;
        }
        const zakres = stan.zaznaczenie();
        zakladki.push({
          kod: `zakladka-${zakladki.length + 1}-${Date.now()}`,
          nazwa: oczyszczona,
          naZnaku: zakres?.poczatek ?? 0,
        });
        przybornik.ustawZakladki(zakladki);
        // Zakładka jest elementem APARATU dokumentu (`StudioApparatusKind.Bookmark`),
        // więc zapisuje się w rdzeniu i przeżywa zapis. Wykaz w oknie zostaje, bo
        // powrót ma działać od ręki, nie po odczycie — ale prawdą o dokumencie
        // jest wiersz w rdzeniu, nie ten wykaz.
        const idDokumentu = stan.dokument()?.id ?? '';
        if (idDokumentu === '') {
          przybornik.pokazOdpowiedz(
            `Zakładka „${oczyszczona}" działa w tym oknie, ale nie ma dokumentu, w którym ` +
              'miałaby się zapisać — otwórz dokument.',
            false,
          );
          return;
        }
        void zrodlaPostaci.aparat
          .zaloz({
            documentId: idDokumentu,
            kind: 'bookmark',
            offset: zakres?.poczatek ?? 0,
            label: oczyszczona,
          })
          .then((wynik) => {
            przybornik.pokazOdpowiedz(
              wynik.udany
                ? `Zakładka „${oczyszczona}" założona na znaku ${zakres?.poczatek ?? 0} ` +
                    'i zapisana w dokumencie — przeżywa zapis i wraca po ponownym otwarciu.'
                : `Zakładka „${oczyszczona}" działa w tym oknie, ale rdzeń jej nie przyjął: ` +
                    `${wynik.blad?.message ?? 'powodu nie podał'}.`,
              wynik.udany,
            );
            void aparat.odswiez();
          });
      },
      naSkok: (pozycja) => {
        if (pozycja.zakres === null) {
          przybornik.pokazOdpowiedz(
            'Ta pozycja nie jest zakotwiczona w treści — adnotacja wisi przy numerze fragmentu ' +
              'różnicy, a nie przy znaku dokumentu, więc nie ma dokąd przewinąć.',
            false,
          );
          return;
        }
        stan.ustawZaznaczenie(pozycja.zakres);
        wskazKomentarz(pozycja.kod);
      },
    },
    mikrofonTresci.element,
    // Trzeci argument jest rdzeniem znakowania — bez niego przybornik pokazuje
    // znakowania sesji i mówi wprost, że droga do `studio.markup.*` nie jest
    // wpięta. Z nim znakowanie jest TRWAŁE i przechodzi przez zapis.
    {
      zrodlo: zrodlaPostaci.kontrola,
      idDokumentu: () => stan.dokument()?.id ?? '',
      zaznaczenie: () => stan.zaznaczenie(),
      naSkutek: przyjmijSkutek,
      naMiejsce: (poczatek: number, koniec: number) =>
        stan.ustawZaznaczenie({ poczatek, koniec }),
    },
  );

  const schowek = utworzSchowekHistorie({
    naWklejenieZPostacia: (tresc) => wstawWTresc(tresc),
    naWklejenieCzyste: (tresc) => wstawWTresc(schowekTekstCzysty(tresc)),
    naOdlozenie: () => void odlozDoSchowka(),
    naPrzypiecie: (idWpisu, przypiety) => void przypnijWpisSchowka(idWpisu, przypiety),
    naUsuniecie: (idWpisu) => void usunWpisSchowka(idWpisu),
    naOdczyt: (fraza, tylkoPrzypiete) => void odczytajSchowek(fraza, tylkoPrzypiete),
    naPobraniePostaci: () => pobierzPostacAkapitu(),
    naNalozeniePostaci: () => nalozPostacAkapitu(),
  });

  /* ── Kontrola pracy: dziennik, kopie, blokady ───────────────────────────── */

  // Trzy panele stoją nakładkami przy powierzchni, bo dotyczą TEJ treści:
  // dziennik cofa jej czynności, kopie ją ratują, blokady jej pilnują. Każdy
  // przyjmuje skutek tą samą drogą co reszta okna — treścią roboczą stanu, więc
  // po cofnięciu i po przywróceniu powierzchnia pokazuje to, co oddał rdzeń.
  // Postać po czynności czytamy z rdzenia, a nie przyjmujemy nieopisanym
  // ładunkiem z odpowiedzi: postać jest drzewem o kształcie kontraktu, a `unknown`
  // wstawione w stan okna byłoby drugą prawdą o dokumencie, tym razem
  // niesprawdzoną kompilacją.
  function przyjmijSkutek(tresc: string | undefined, _postacPoZmianie: unknown): void {
    if (tresc !== undefined) {
      stan.ustawTresc(tresc);
      powierzchnia.pokaz(tresc, true);
    }
    postac.odswiez();
  }

  const dziennik = utworzDziennikPanel(zrodlaPostaci.kontrola, {
    idDokumentu: () => stan.dokument()?.id ?? '',
    zaznaczenie: () => stan.zaznaczenie(),
    kursor: () => stan.zaznaczenie()?.poczatek ?? stan.trescRobocza().length,
    naSkutek: przyjmijSkutek,
    // Skok po zmianach modelu przestawia zaznaczenie stanu — powierzchnia idzie
    // za nim tą samą drogą, którą idzie za zaznaczeniem Operatora.
    naMiejsce: (poczatek, koniec) => stan.ustawZaznaczenie({ poczatek, koniec }),
  });

  const kopie = utworzKopiePanel(zrodlaPostaci.kontrola, {
    idDokumentu: () => stan.dokument()?.id ?? '',
    idOkna: () => stan.idOkna(),
    tresc: () => stan.trescRobocza(),
    postac: () => postac.postacZapamietana(),
    zmianyNiezapisane: () => stan.trescRobocza() !== (stan.dokument()?.content ?? ''),
    naSkutek: przyjmijSkutek,
  });

  const blokady = utworzBlokadaPanel(zrodlaPostaci.kontrola, {
    idDokumentu: () => stan.dokument()?.id ?? '',
    zaznaczenie: () => stan.zaznaczenie(),
  });

  /* ── Wstawienia: tabele, obiekty, aparat, szablony ──────────────────────── */

  const tabele = utworzTabelaPanel(stan, zrodlaPostaci.tabele);
  const obiekty = utworzObiektPanel(stan, zrodlaPostaci.obiekty);
  const aparat = utworzAparatPanel(stan, zrodlaPostaci.aparat);
  const szablony = utworzSzablonPanel(stan, zrodlaPostaci.szablony);

  const osadzenie = utworzOsadzeniePanel({
    naSzukanieBiblioteki: (fraza) => void szukajWBibliotece(fraza),
    naPodglad: (idPliku, strona) => void podejrzyjPlik(idPliku, strona),
    naWniesienieZBiblioteki: (plik, podgladPliku, tresc) =>
      wniesZBiblioteki(plik, podgladPliku, tresc),
    naWciagniecieStrony: (adres, zObrazami) => void wciagnijStrone(adres, zObrazami),
    naMigawke: (idOknaPrzegladarki) => void wezMigawke(idOknaPrzegladarki),
    naOtwarcieStrony: (idOknaPrzegladarki, adres) =>
      void otworzStroneWPrzegladarce(idOknaPrzegladarki, adres),
    naWniesienieZeStrony: (tytul, adres, tresc) => wniesZeStrony(tytul, adres, tresc),
  }, { stan, wstawienia: zrodlaPostaci.wstawienia });

  const redaktor = utworzPanelRedaktora({
    naUscislenie: (idAkcji) => void zlecOperacje(zapleczePracy, idAkcji, '', nastawySuwakow),
    naPodobienstwa: (wskazanie) => void zmierzPodobienstwa(wskazanie),
    naWyszukanie: (zapytanie) => void wyszukajZnaczeniowo(zapytanie),
  });

  const galeria = utworzGalerieSzablonow({
    naZalozenie: (idSzablonu, wartosci, tytul) =>
      void zalozZSzablonu(idSzablonu, wartosci, tytul),
    // Galeria zakłada dokument z szablonu; warsztat szablon ZMIENIA. To dwie
    // czynności i dwa miejsca, więc galeria nie udaje warsztatu, tylko do niego
    // prowadzi.
    naWarsztat: (idSzablonu) => szablony.wskaz(idSzablonu),
  });

  /* ── Kawałki przeniesione z okien scalonych ──────────────────────────────── */

  const wczytanie = utworzWczytanieDokumentu();
  // Wniesienie pliku i wydanie do formatu stoją nakładkami, bo obie czynności
  // wychodzą poza treść: pierwsza ją zakłada, druga wypuszcza na zewnątrz wraz
  // z bilansem cech pominiętych.
  const wniesienie = utworzWniesieniePliku(stan, zrodlaPostaci.wstawienia);
  const wydanie = utworzWydaniePanel(stan, zrodlaPostaci.wstawienia);
  const status = utworzPasekStatusu(stan);
  const szukanie = utworzZnajdzZamien(stan);
  const pola = utworzPolaRoznicy();

  const filtr = poleWyboru(
    {
      etykieta: 'Filtr różnic',
      opis: 'Zawęża wykaz po rodzaju fragmentu. Statystyka liczy CAŁĄ różnicę, nie sam wykaz widoczny.',
    },
    [
      { wartosc: 'wszystkie', etykieta: 'Wszystkie fragmenty' },
      { wartosc: DiffHunkKind.Added, etykieta: 'Tylko dodania' },
      { wartosc: DiffHunkKind.Removed, etykieta: 'Tylko usunięcia' },
      { wartosc: DiffHunkKind.Changed, etykieta: 'Tylko zmiany' },
      { wartosc: DiffHunkKind.Context, etykieta: 'Tylko kontekst' },
    ],
  );

  const statystyka = document.createElement('p');
  statystyka.className = 'dn-pole-opis ms-roznica__statystyka';

  const wykazFragmentow = document.createElement('ul');
  wykazFragmentow.className = 'ms-roznica';

  const wskazaneFragmenty = document.createElement('input');
  wskazaneFragmenty.type = 'text';
  wskazaneFragmenty.className = 'dn-pole-kontrolka';
  wskazaneFragmenty.placeholder = 'numery fragmentów do przyjęcia, np. 2, 3';
  wskazaneFragmenty.setAttribute('aria-label', 'Numery fragmentów objętych decyzją');

  const przyjmijFragmenty = document.createElement('button');
  przyjmijFragmenty.type = 'button';
  przyjmijFragmenty.className = 'dn-btn dn-btn--sm dn-btn--sygnal';
  przyjmijFragmenty.textContent = 'Przyjmij wskazane fragmenty';
  przyjmijFragmenty.dataset['czynnosc'] = 'przyjmij-fragmenty';
  przyjmijFragmenty.title =
    'Idzie komendą studio.proposal.decide z polem hunkIndexes — rdzeń składa treść z fragmentów ' +
    'wskazanych, a pozostałe zostawia w postaci z dokumentu.';
  przyjmijFragmenty.addEventListener('click', () => {
    void rozstrzygnijPropozycje(zapleczePracy, true, numeryFragmentow());
  });

  const adnotacja = document.createElement('button');
  adnotacja.type = 'button';
  adnotacja.className = 'dn-btn dn-btn--sm dn-btn--duch';
  adnotacja.textContent = 'Dodaj adnotację do różnicy';
  adnotacja.dataset['czynnosc'] = 'adnotacja';
  adnotacja.addEventListener('click', () => void dodajAdnotacje());

  const wstaw = document.createElement('button');
  wstaw.type = 'button';
  wstaw.className = 'dn-btn dn-btn--sm dn-btn--duch';
  wstaw.textContent = 'Wstaw wynik w miejsce kursora';
  wstaw.dataset['czynnosc'] = 'wstaw';
  wstaw.title =
    'Wstawienie nie jest przyjęciem: treść wyniku wchodzi do bufora w miejscu kursora, ' +
    'a decyzja o zmianie modelu czeka dalej.';

  const pochodzenie = document.createElement('p');
  pochodzenie.className = 'dn-pole-opis ms-praca__pochodzenie';

  /**
   * Wskaźnik zakresu operacji — przeniesiony z paska zaznaczenia Studio Editora.
   *
   * Mówi, co pojedzie do rdzenia: zaznaczenie wraz z jego długością, cały
   * dokument, albo cały dokument z zaznaczeniem POMINIĘTYM wyborem ręcznym
   * z Tools Panelu. Trzeci przypadek jest osobny, bo bez niego wybór ręczny
   * wyglądałby jak usterka: zaznaczenie widać, a operacja idzie na całość.
   */
  const zakresOperacji = document.createElement('span');
  zakresOperacji.className = 'dn-plakietka ms-praca__zakres';

  function opiszZakres(): string {
    const wybrany = stan.zaznaczenie();
    const dlugosc = wybrany === null ? 0 : wybrany.koniec - wybrany.poczatek;
    if (stan.zakresSkuteczny() === StudioOperationScope.Selection && wybrany !== null) {
      return `zakres operacji: zaznaczenie · ${dlugosc} znaków (od ${wybrany.poczatek} do ${wybrany.koniec})`;
    }
    if (wybrany !== null) {
      return (
        `zakres operacji: cały dokument — wybór ręczny w Tools Panelu ma pierwszeństwo, ` +
        `więc zaznaczenie (${dlugosc} znaków) zostaje pominięte do następnej zmiany zaznaczenia`
      );
    }
    if (stan.zakresZadany() === StudioOperationScope.Selection) {
      return (
        'zakres operacji: cały dokument — wybrano „zaznaczenie", ale w treści nic nie jest ' +
        'zaznaczone; wybór nie jest blokowany'
      );
    }
    return 'zakres operacji: cały dokument — w treści nic nie jest zaznaczone';
  }

  const fragmentyGniazdo = document.createElement('div');
  fragmentyGniazdo.className = 'ms-praca__fragmenty';
  fragmentyGniazdo.append(
    filtr.element,
    statystyka,
    wskazaneFragmenty,
    przyjmijFragmenty,
    adnotacja,
    wykazFragmentow,
  );

  /* ── Wstążka ─────────────────────────────────────────────────────────────── */

  const wstazka = utworzWstazkePracy(
    {
      naZapis: () => void zapiszDokument(zapleczeEdytora),
      naGalerie: () => galeria.przestawWidocznosc(),
      // Wydanie i przekazanie idą tą samą drogą, którą szły z Preview Window:
      // `czynnosci-podgladu.ts`. Druga kopia tych dwóch czynności rozjechałaby się
      // z pierwszą przy pierwszej poprawce zdania o wyniku.
      naWydanie: (format) => void eksportujDokument(zapleczePodgladu(), format),
      naPrzekazanie: () => void przekazDoLibrary(zapleczePodgladu(), przekazanie),
      naNarzedzieTekstu: (narzedzie) => {
        const zakres = stan.zaznaczenie();
        const wynik = zastosujNarzedzie(
          {
            tresc: stan.trescRobocza(),
            poczatek: zakres?.poczatek ?? stan.trescRobocza().length,
            koniec: zakres?.koniec ?? stan.trescRobocza().length,
          },
          narzedzie,
        );
        stan.ustawTresc(wynik.tresc);
        stan.ustawZaznaczenie(
          wynik.poczatek === wynik.koniec
            ? null
            : { poczatek: wynik.poczatek, koniec: wynik.koniec },
        );
      },
      naStyl: (rodzaj) => powierzchnia.ustawStylBloku(rodzaj),
      naKrój: (krój) => przestawNastawy({ krój }),
      naStopien: (stopien) => przestawNastawy({ stopien }),
      naInterlinie: (interlinia) => przestawNastawy({ interlinia }),
      naWciecie: (wciecieMm) => przestawNastawy({ wciecieMm }),
      naOdstep: (odstepMm) => przestawNastawy({ odstepMm }),
      naWyrownanie: (wyrownanie) => {
        const numer = powierzchnia.blokKursora();
        if (numer < 0) {
          odpowiedz.pokaz(
            'Wyrównanie dotyczy akapitu, w którym stoi kursor — postaw go w treści.',
            false,
          );
          return;
        }
        akapity.ustaw(numer, { wyrownanie });
        powierzchnia.pokaz(stan.trescRobocza(), true);
      },
      naBarwe: (barwa) => {
        const numer = powierzchnia.blokKursora();
        if (numer < 0) return;
        akapity.ustaw(numer, { barwa });
        powierzchnia.pokaz(stan.trescRobocza(), true);
      },
      naSzukanie: () => {
        szukanie.element.hidden = !szukanie.element.hidden;
        szukanie.odswiez();
      },
      naPodzialStrony: () => {
        const tresc = stan.trescRobocza();
        stan.ustawTresc(`${tresc}\n\n${ZNACZNIK_PODZIALU}\n\n`);
      },
      naProfil: (idProfilu) => void wezProfil(idProfilu),
      naZapisProfilu: (nazwa) => void zapiszProfil(nazwa),
      naNosnik: (oznaczenie) => {
        const nosnik = nosnikPoOznaczeniu(oznaczenie);
        if (nosnik === null) {
          odpowiedz.pokaz(
            `Nośnika ${oznaczenie} okno nie zna z wymiarów — kartka zostaje bez zmiany, ` +
              'zamiast być rysowana w rozmiarze zgadniętym.',
            false,
          );
          return;
        }
        przestawStrone({ nosnik });
      },
      naOrientacje: (pozioma) =>
        przestawStrone({
          orientacja: pozioma ? StudioPageOrientation.Pozioma : StudioPageOrientation.Pionowa,
        }),
      naMargines: (ktory, milimetry) => {
        if (ktory === 'gora') przestawStrone({ marginesGoraMm: milimetry });
        if (ktory === 'dol') przestawStrone({ marginesDolMm: milimetry });
        if (ktory === 'lewy') przestawStrone({ marginesLewyMm: milimetry });
        if (ktory === 'prawy') przestawStrone({ marginesPrawyMm: milimetry });
      },
      naNaglowekStrony: (tresc) => przestawStrone({ naglowek: tresc }),
      naStopkeStrony: (tresc) => przestawStrone({ stopka: tresc }),
      naNumeracje: (czynna) => przestawStrone({ numeracja: czynna }),
      naSkale: (procent) => przestawStrone({ skala: procent }),
      naKolumny: (kolumny) => powierzchnia.ustawKolumny(kolumny),
      naSledzenie: (czynne) => void przestawSledzenieOkna(czynne),
      naDecyzjeWszystkich: (przyjmij) =>
        void rozstrzygnijZmiany(
          zapleczePracy,
          zmiany.map((zmiana) => zmiana.id),
          przyjmij,
        ),
      naSkokZmiany: (wPrzod) => {
        const kod = powierzchnia.skoczDoZmiany(wPrzod);
        odpowiedz.pokaz(
          kod === null
            ? 'Nie ma zmian oczekujących decyzji — nie ma po czym skakać.'
            : `Kursor stoi przy zmianie ${kod}. Decyzję podejmiesz znacznikiem w treści albo ` +
                'paskiem zatwierdzenia pod nią.',
          kod !== null,
        );
      },
      naAdiustacje: (tryb) => powierzchnia.ustawAdiustacje(tryb),
      naNowyKomentarz: () => {
        const polozenie = powierzchnia.polozenieKursora();
        if (polozenie === null) {
          odpowiedz.pokaz('Postaw kursor w treści — komentarz przypina się do miejsca.', false);
          return;
        }
        wiersz.pokaz(polozenie, stan.zaznaczenie() === null ? 0 : 1);
      },
      naSkokKomentarza: (wPrzod) => skoczDoKomentarza(wPrzod),
      naDymki: () => {
        dymki.element.hidden = !dymki.element.hidden;
      },
      naPorownanie: () => void porownajWersje(),
      naRoznicaWygladu: () => void porownajWyglad(),
      naOchrone: () => czynnosciZewnetrzne.naWarsztat(),
      naTryb: (tryb) => ustawTryb(tryb),
      naRedaktora: () => redaktor.przestawWidocznosc(),
      naSuwak: (kod, wartosc) => {
        nastawySuwakow = { ...nastawySuwakow, [kod]: wartosc };
      },
      naOperacje: (idAkcji) => void zlecOperacje(zapleczePracy, idAkcji, '', nastawySuwakow),
      naDecyzjePropozycji: (przyjmij, fragmenty) =>
        void rozstrzygnijPropozycje(zapleczePracy, przyjmij, fragmenty),
    },
    {
      wczytanie: wczytanie.element,
      szukanie: szukanie.element,
      roznica: pola.elementy,
      fragmenty: fragmentyGniazdo,
    },
  );

  /* ── Zakładki dokumentów ─────────────────────────────────────────────────── */

  const karty = utworzKartyDokumentow(stan, () => {
    // Przełączenie zakładki to inny dokument czynny: zmiany, komentarze
    // i odpowiedź różnicy dotyczyły poprzedniego, więc schodzą do czasu odczytu.
    zmiany = [];
    komentarze = [];
    adnotacje = [];
    fragmentyOdpowiedzi = [];
    trafieniaOdpowiedzi = [];
    dokumentOdpowiedzi = '';
    // Znaczniki, zakładki i pochodzenia dotyczą dokumentu, nie okna: przeniesienie
    // ich na dokument drugi wskazywałoby fragmenty, których w nim nie ma.
    znaczniki.wyczysc();
    pochodzenia.wyczysc();
    zakladki.length = 0;
    odswiez();
    void odswiezZRdzenia();
    // Trzeci argument daje zakładkom drogę do `studio.document.create`: nowy
    // dokument staje OBOK, nie zamiast — zakładka bez tej drogi mówiła wprost,
    // że założenia dokumentu nie ma czym zlecić.
  }, zrodlaPostaci.wstawienia);

  /* ── Układ okna ──────────────────────────────────────────────────────────── */

  const powloka = document.createElement('div');
  powloka.className = 'ms-praca__powloka';
  // Panel osadzenia źródeł stoi wewnątrz powłoki treści, bo w położeniu „na całej
  // powierzchni" przykrywa kartkę, a nie całe okno: wstążka i pasek statusu mają
  // zostać widoczne, żeby Operator nie stracił drogi powrotu.
  powloka.append(powierzchnia.element, wiersz.element, osadzenie.element);

  const kolumny = document.createElement('div');
  kolumny.className = 'ms-praca__kolumny';
  // Przybornik znakowania stoi PRZY KRAWĘDZI treści, między marginesem z dymkami
  // a panelem Redaktora: znakowanie fragmentu należy do treści bliżej niż pomiary
  // dokumentu.
  kolumny.append(powloka, dymki.element, przybornik.element, redaktor.element);

  const decyzja = document.createElement('div');
  decyzja.className = 'ms-decyzja';
  decyzja.append(zakresOperacji, wstaw);

  /**
   * Pas otwarć odcinka znakowania.
   *
   * Trzy panele — schowek, źródła i przybornik — otwiera się stąd, a nie ze
   * wstążki: wstążka niesie formatowanie i widok, a te trzy są narzędziami
   * bocznymi treści. Pas jest jednym wierszem, więc zamknięte panele nie zajmują
   * powierzchni.
   */
  const pasPrzybornika = document.createElement('div');
  pasPrzybornika.className = 'ms-przybornik__pas';

  const otworzSchowek = document.createElement('button');
  otworzSchowek.type = 'button';
  otworzSchowek.className = 'dn-btn dn-btn--sm dn-btn--zarys';
  otworzSchowek.textContent = 'Schowek z historią';
  otworzSchowek.dataset['czynnosc'] = 'otworz-schowek';
  otworzSchowek.addEventListener('click', () => schowek.przestawWidocznosc());

  const otworzZrodla = document.createElement('button');
  otworzZrodla.type = 'button';
  otworzZrodla.className = 'dn-btn dn-btn--sm dn-btn--zarys';
  otworzZrodla.textContent = 'Źródła: przeglądarka i Biblioteka';
  otworzZrodla.dataset['czynnosc'] = 'otworz-zrodla';
  otworzZrodla.addEventListener('click', () => {
    osadzenie.przestawWidocznosc();
    if (osadzenie.widoczny()) {
      osadzenie.pokazPochodzenia(pochodzenia.wykaz());
      void szukajWBibliotece('');
    }
  });

  const otworzPrzybornik = document.createElement('button');
  otworzPrzybornik.type = 'button';
  otworzPrzybornik.className = 'dn-btn dn-btn--sm dn-btn--zarys';
  otworzPrzybornik.textContent = 'Przybornik znakowania';
  otworzPrzybornik.dataset['czynnosc'] = 'otworz-przybornik';
  otworzPrzybornik.addEventListener('click', () => przybornik.przestawWidocznosc());

  /**
   * Otwarcia paneli postaci i kontroli pracy.
   *
   * Siedem powierzchni idzie jednym wzorem: przycisk w pasie, treść nakładką.
   * Panel zamknięty nie zajmuje powierzchni, a Operator widzi w pasie, co ma pod
   * ręką — to jest reguła stopniowego ujawniania, nie oszczędność miejsca.
   */
  function otwarcie(nazwa: string, czynnosc: string, przestaw: () => void): HTMLButtonElement {
    const przycisk = document.createElement('button');
    przycisk.type = 'button';
    przycisk.className = 'dn-btn dn-btn--sm dn-btn--zarys';
    przycisk.textContent = nazwa;
    przycisk.dataset['czynnosc'] = czynnosc;
    przycisk.addEventListener('click', przestaw);
    return przycisk;
  }

  pasPrzybornika.append(
    otworzPrzybornik,
    otworzSchowek,
    otworzZrodla,
    otwarcie('Tabele', 'otworz-tabele', () => tabele.przestawWidocznosc()),
    otwarcie('Obrazy i kształty', 'otworz-obiekty', () => obiekty.przestawWidocznosc()),
    otwarcie('Aparat i pola', 'otworz-aparat', () => aparat.przestawWidocznosc()),
    otwarcie('Warsztat szablonów', 'otworz-szablony', () => szablony.przestawWidocznosc()),
    otwarcie('Dziennik i zmiany modelu', 'otworz-dziennik', () => dziennik.przestawWidocznosc()),
    otwarcie('Kopie i autozapis', 'otworz-kopie', () => kopie.przestawWidocznosc()),
    otwarcie('Blokady fragmentów', 'otworz-blokady', () => blokady.przestawWidocznosc()),
  );

  rama.stan.tresc.append(
    wstazka.element,
    karty.element,
    galeria.element,
    pasPrzybornika,
    schowek.element,
    postac.element,
    wniesienie.element,
    wydanie.element,
    tabele.element,
    obiekty.element,
    aparat.element,
    szablony.element,
    dziennik.element,
    kopie.element,
    blokady.element,
    kolumny,
    pochodzenie,
    decyzja,
    status.element,
    utworzPowiazanieRozmowy().element,
    odpowiedz.element,
  );

  szukanie.element.hidden = true;
  wstaw.addEventListener('click', () => {
    const propozycja = stan.propozycja();
    if (propozycja === null || propozycja.tresc === '') {
      odpowiedz.pokaz(
        'Nie ma wyniku operacji do wstawienia — zleć operację wierszem polecenia albo ze wstążki.',
        false,
      );
      return;
    }
    // Wstawienie idzie przez ten sam rachunek co w Studio Editorze, tylko kursor
    // czyta się z powierzchni, a nie z pola tekstowego.
    const atrapa = document.createElement('textarea');
    atrapa.value = stan.trescRobocza();
    const zakres = stan.zaznaczenie();
    atrapa.selectionStart = zakres?.poczatek ?? atrapa.value.length;
    atrapa.selectionEnd = zakres?.koniec ?? atrapa.value.length;
    wstawWMiejsceKursora(stan, atrapa, odpowiedz);
  });

  wczytanie.wczytaj.addEventListener('click', () => {
    void wczytajDokument(zapleczeEdytora, wczytanie.zadanie(stan.idOkna()), wczytanie.brak()).then(
      () => {
        karty.odswiez();
        void odswiezZRdzenia();
      },
    );
  });

  filtr.kontrolka.addEventListener('change', () => przerysujRoznice());

  /* ── Nastawy: przestawianie ──────────────────────────────────────────────── */

  function przestawStrone(zmiana: Partial<typeof strona>): void {
    strona = { ...strona, ...zmiana };
    powierzchnia.ustawStrone(strona);
    odswiez();
  }

  function przestawNastawy(zmiana: Partial<typeof nastawyWizualne>): void {
    nastawyWizualne = { ...nastawyWizualne, ...zmiana };
    powierzchnia.ustawNastawy(nastawyWizualne);
    odswiez();
  }

  function ustawTryb(tryb: TrybWidoku): void {
    powierzchnia.ustawTryb(tryb);
    if (tryb === 'wydanie') {
      // Podgląd wydania czyta treść ZAAKCEPTOWANĄ — pisanie w oknie podglądu nie
      // rusza, rusza nim zapis albo decyzja o zmianie. Tak samo czytał ją Preview
      // Window przed scaleniem.
      powierzchnia.pokaz(stan.trescZaakceptowana(), true);
      void zlecRender();
      return;
    }
    powierzchnia.pokaz(stan.trescRobocza(), true);
  }

  /* ── Rozmowa z rdzeniem ──────────────────────────────────────────────────── */

  async function przestawSledzenieOkna(czynne: boolean): Promise<void> {
    const stanSledzenia = await przestawSledzenie(zapleczePracy, czynne);
    wstazka.odswiezSledzenie(stanSledzenia);
  }

  async function odswiezZRdzenia(): Promise<void> {
    const dokument = stan.dokument();
    if (dokument === null) {
      zmiany = [];
      komentarze = [];
      adnotacje = [];
      powierzchnia.ustawZmiany(zmiany);
      przerysujMargines(() => null);
      return;
    }
    zmiany = await odczytajZmiany(zapleczePracy);
    powierzchnia.ustawZmiany(zmiany);
    const wynik = await praca.komentarze(dokument.id, true);
    komentarze = wynik.udany && wynik.wynik !== undefined ? wynik.wynik.comments : [];
    powierzchnia.ustawKomentarze(komentarze);
    przerysujMargines((kod) => powierzchnia.polozenieKotwicy(kod));
    odswiez();
  }

  /**
   * Przerysowuje margines trzema bytami naraz.
   *
   * Jedno wywołanie, bo komentarz, propozycja i zmiana śledzona stoją w jednej
   * kolumnie i muszą się ustawić wobec siebie — dwa osobne przerysowania
   * kładłyby jedne na drugich.
   */
  function przerysujMargines(polozenie: (kod: string) => number | null): void {
    dymki.pokaz(
      {
        komentarze,
        propozycja: stan.propozycja(),
        zakresPropozycji: stan.zaznaczenie(),
        zmiany,
      },
      polozenie,
    );
  }

  async function rozstrzygnijKomentarz(idKomentarza: string, rozwiazany: boolean): Promise<void> {
    const wynik = await praca.rozstrzygnijKomentarz(idKomentarza, rozwiazany);
    if (!wynik.udany) {
      odpowiedz.pokaz(
        opisOdmowy('Rozstrzygnięcie wątku', wynik.blad?.code, wynik.blad?.message),
        false,
      );
      return;
    }
    odpowiedz.pokaz(
      rozwiazany ? 'Wątek oznaczony jako rozwiązany.' : 'Wątek otwarty ponownie.',
      true,
    );
    await odswiezZRdzenia();
  }

  async function porownajWersje(): Promise<void> {
    const zadanie = pola.zadanie(stan);
    if (zadanie === null) {
      odpowiedz.pokaz(
        'Porównanie wymaga dokumentu wczytanego — komenda studio.diff.compare przyjmuje jego ' +
          'identyfikator.',
        false,
      );
      return;
    }
    rama.stan.ladowanie('Porównanie wersji i wyszukiwanie wzorca w toku…');
    const wynik = await stan.zrodlo.porownaj(zadanie);
    dokumentOdpowiedzi = stan.dokument()?.id ?? '';
    if (!wynik.udany || wynik.wynik === undefined) {
      rama.stan.blad(opisOdmowy('Porównanie wersji', wynik.blad?.code, wynik.blad?.message));
      return;
    }
    fragmentyOdpowiedzi = wynik.wynik.hunks ?? [];
    trafieniaOdpowiedzi = wynik.wynik.matches ?? [];
    przerysujRoznice();
    powierzchnia.ustawFragmenty(fragmentyOdpowiedzi);
    rama.stan.gotowe();
    if (fragmentyOdpowiedzi.length === 0 && trafieniaOdpowiedzi.length === 0) {
      const bezStrony =
        zadanie.baseVersionId === undefined &&
        zadanie.targetVersionId === undefined &&
        zadanie.proposalId === undefined;
      odpowiedz.pokaz(
        bezStrony
          ? 'Nie było czego porównać ani przeszukać: rdzeń czyta treść WERSJI, a ten dokument ' +
              'jeszcze żadnej nie ma. Zapisz go — zapis zakłada pierwszą wersję sesji.'
          : 'Rdzeń porównał wskazane strony i nie oddał ani fragmentu różnicy, ani trafienia ' +
              'wzorca: strony są zgodne, a wzorca w nich nie ma.',
        false,
      );
      return;
    }
    odpowiedz.pokaz(
      `Rdzeń oddał fragmentów ${fragmentyOdpowiedzi.length} i trafień ${trafieniaOdpowiedzi.length}. ` +
        'Przełącz widok na „Różnica na treści", żeby zobaczyć je w miejscu.',
      true,
    );
  }

  async function porownajWyglad(): Promise<void> {
    const dokument = stan.dokument();
    const para = stan.paraPorownania();
    if (dokument === null || para === null || para.odniesienie === '' || para.porownywana === '') {
      odpowiedz.pokaz(
        'Różnica wyglądu wymaga DWÓCH wskazanych wersji — komenda studio.diff.visual przyjmuje ' +
          'baseVersionId i targetVersionId jako obowiązkowe. Wskaż parę w Session Repository.',
        false,
      );
      return;
    }
    const wynik = await praca.roznicaWygladu(
      dokument.id,
      para.odniesienie,
      para.porownywana,
      stan.idOkna(),
    );
    if (!wynik.udany || wynik.wynik === undefined) {
      odpowiedz.pokaz(
        opisOdmowy('Różnica wyglądu', wynik.blad?.code, wynik.blad?.message),
        false,
      );
      return;
    }
    // Nakładki różnicy wyglądu SĄ obrazami stron i pobierają się tą samą drogą,
    // co kartki podglądu wydania: `design.asset.content.get` dotyczy każdego
    // zasobu magazynu, bo magazyn jest jeden dla całej platformy. Zdanie „komendy
    // pobierającej ich bajty kontrakt nie niesie" stało tu do 17.08.2026 i było
    // nieprawdą — komenda była, brakowało wołacza.
    const nakladki = wynik.wynik.overlayAssetIds ?? [];
    const obrazy: KartkaRdzenia[] = [];
    for (const [numer, kod] of nakladki.slice(0, GRANICA_KARTEK_RDZENIA).entries()) {
      const tresc = await praca.trescZasobu(kod, GRANICA_BAJTOW_KARTKI);
      const bajty = tresc.udany ? (tresc.wynik?.contentBase64 ?? '') : '';
      if (bajty === '') continue;
      obrazy.push({
        numer: numer + 1,
        zrodlo: `data:${tresc.wynik?.mediaType ?? 'image/png'};base64,${bajty}`,
      });
    }
    const zdanie =
      obrazy.length > 0
        ? `Nakładki widać w trybie podglądu wydania (${obrazy.length} z ${nakladki.length}).`
        : 'Bajtów nakładek rdzeń nie oddał, więc okno pokazuje same liczby obszarów.';
    if (obrazy.length > 0) powierzchnia.ustawKartkiRdzenia(obrazy);
    odpowiedz.pokaz(
      `Rdzeń wyrysował różnicę wyglądu: obszarów ${wynik.wynik.regions.length}, nakładek ` +
        `${nakladki.length}. ${zdanie}`,
      true,
    );
  }

  async function zlecRender(): Promise<void> {
    const dokument = stan.dokument();
    if (dokument === null) return;
    const wynik = await praca.render(dokument.id, 'pdf', stan.idOkna(), '');
    if (!wynik.udany || wynik.wynik === undefined) {
      powierzchnia.ustawKartkiRdzenia([]);
      odpowiedz.pokaz(
        `Podgląd okna rysuje kartki po stronie klienta; render rdzenia odmówił: ` +
          opisOdmowy('Render podglądu', wynik.blad?.code, wynik.blad?.message),
        false,
      );
      return;
    }

    // Zasoby stron stoją na POCZĄTKU wykazu, a dokument w formacie docelowym —
    // przy formacie `pdf` — dochodzi na jego końcu. Kartek jest więc dokładnie
    // `pages`, i to jest liczba, którą wykaz się przycina; rozpoznawanie kartki
    // po typie treści wymagałoby pobrania także dokumentu, czyli megabajtów
    // pobranych po to, żeby je odrzucić.
    const kartek = Math.min(wynik.wynik.pages, GRANICA_KARTEK_RDZENIA);
    const kody = wynik.wynik.pageAssetIds.slice(0, kartek);
    const kartki: KartkaRdzenia[] = [];
    const nieudane: string[] = [];
    for (const [numer, kod] of kody.entries()) {
      const tresc = await praca.trescZasobu(kod, GRANICA_BAJTOW_KARTKI);
      if (!tresc.udany || tresc.wynik === undefined) {
        nieudane.push(kod);
        continue;
      }
      const bajty = tresc.wynik.contentBase64 ?? '';
      if (bajty === '') {
        nieudane.push(kod);
        continue;
      }
      kartki.push({ numer: numer + 1, zrodlo: `data:${tresc.wynik.mediaType};base64,${bajty}` });
    }
    powierzchnia.ustawKartkiRdzenia(kartki);

    const czesci = [`Rdzeń policzył paginację wydania: stron ${wynik.wynik.pages}.`];
    if (kartki.length > 0) {
      czesci.push(
        `Podgląd wydania pokazuje kartki WYRYSOWANE PRZEZ RDZEŃ (${kartki.length} z ` +
          `${wynik.wynik.pages}) — typografia, paginacja, nagłówek i stopka są tam takie, ` +
          'jakie wyjdą w wydaniu.',
      );
    }
    if (kartek < wynik.wynik.pages) {
      czesci.push(
        `Pobrano pierwsze ${kartek} kartek; dalsze zostają zasobami magazynu rdzenia, bo ` +
          'podgląd wydania służy sprawdzeniu składu, a nie czytaniu pisma obrazkami.',
      );
    }
    if (nieudane.length > 0) {
      czesci.push(
        `Bajtów ${nieudane.length} kartek rdzeń nie oddał — te miejsca pokazują kartki ` +
          'liczone w oknie.',
      );
    }
    if (kartki.length === 0) {
      czesci.push('Kartki w oknie liczy klient, bo z rdzenia nie przyszedł ani jeden obraz.');
    }
    odpowiedz.pokaz(czesci.join(' '), true);
  }

  async function zalozZSzablonu(
    idSzablonu: string,
    wartosci: Record<string, string>,
    tytul: string,
  ): Promise<void> {
    rama.stan.ladowanie('Zakładanie dokumentu z szablonu…');
    const wynik = await praca.zastosujSzablon(stan.idOkna(), idSzablonu, wartosci, tytul);
    if (!wynik.udany || wynik.wynik === undefined) {
      rama.stan.blad(
        opisOdmowy('Zastosowanie szablonu', wynik.blad?.code, wynik.blad?.message),
      );
      return;
    }
    stan.wchlon(wynik.wynik.document);
    rama.stan.gotowe();
    galeria.przestawWidocznosc();
    karty.odswiez();
    odpowiedz.pokaz(
      `Rdzeń założył dokument ${wynik.wynik.document.id} z szablonu ${idSzablonu}.`,
      true,
    );
    await odswiezZRdzenia();
  }

  async function wezProfil(idProfilu: string): Promise<void> {
    if (idProfilu === '') {
      odpowiedz.pokaz('Nastawy strony zostają własne — bez profilu z rdzenia.', true);
      return;
    }
    const wynik = await praca.profile();
    if (!wynik.udany || wynik.wynik === undefined) return;
    const profil = wynik.wynik.profiles.find((pozycja) => pozycja.id === idProfilu);
    if (profil === undefined) return;
    strona = stronaZProfilu(profil.pageSetup);
    powierzchnia.ustawStrone(strona);
    odpowiedz.pokaz(`Nastawy strony wzięte z profilu „${profil.name}".`, true);
    odswiez();
  }

  async function zapiszProfil(nazwa: string): Promise<void> {
    if (nazwa === '') {
      odpowiedz.pokaz('Profil bez nazwy nie da się później wskazać — nazwij go.', false);
      return;
    }
    const { profilZeStrony } = await import('./nastawy-strony');
    const wynik = await praca.zapiszProfil(nazwa, 'pdf', profilZeStrony(strona));
    if (!wynik.udany || wynik.wynik === undefined) {
      odpowiedz.pokaz(opisOdmowy('Zapis profilu', wynik.blad?.code, wynik.blad?.message), false);
      return;
    }
    odpowiedz.pokaz(
      `Nastawy strony zapisane w rdzeniu jako profil „${wynik.wynik.profile.name}".`,
      true,
    );
    await wczytajProfile();
  }

  async function wczytajProfile(): Promise<void> {
    const wynik = await praca.profile();
    if (!wynik.udany || wynik.wynik === undefined) return;
    wstazka.ustawProfile(
      wynik.wynik.profiles.map((profil) => ({ id: profil.id, nazwa: profil.name })),
    );
  }

  async function zmierzPodobienstwa(wskazanie: string): Promise<void> {
    const dokument = stan.dokument();
    if (dokument === null || wskazanie === '') {
      odpowiedz.pokaz(
        'Zestawienie ze źródłem wymaga dokumentu i wskazania materiału wejściowego — ' +
          'komenda studio.diff.source przyjmuje plik Library albo zasób magazynu.',
        false,
      );
      return;
    }
    const wynik = await praca.zestawZeZrodlem(dokument.id, { plikLibrary: wskazanie });
    if (!wynik.udany || wynik.wynik === undefined) {
      odpowiedz.pokaz(
        opisOdmowy('Zestawienie ze źródłem', wynik.blad?.code, wynik.blad?.message),
        false,
      );
      return;
    }
    redaktor.ustawPodobienstwa(wynik.wynik.hunks, wynik.wynik.sourceResolved);
  }

  async function wyszukajZnaczeniowo(zapytanie: string): Promise<void> {
    const dokument = stan.dokument();
    if (dokument === null || zapytanie === '') {
      odpowiedz.pokaz('Wyszukiwanie znaczeniowe wymaga dokumentu i zapytania.', false);
      return;
    }
    const wynik = await praca.wyszukajZnaczeniowo(dokument.id, zapytanie);
    if (!wynik.udany || wynik.wynik === undefined) {
      odpowiedz.pokaz(
        opisOdmowy('Wyszukiwanie znaczeniowe', wynik.blad?.code, wynik.blad?.message),
        false,
      );
      return;
    }
    redaktor.ustawTrafienia(wynik.wynik.matches, wynik.wynik.mode);
  }

  /**
   * Adnotacja przy fragmencie różnicy — droga własnej komendy.
   *
   * Dotąd szła generyczną `window.action` z identyfikatorem `studio.diff.adnotacja`,
   * którego katalog akcji rdzenia nie ma, więc wracała odmowa `not_found` przy
   * każdym naciśnięciu. Kontrakt ma parę komend własnych — `studio.annotation.add`
   * i `studio.annotation.list` — i tędy adnotacja jedzie teraz. Treść bierze się
   * z pola przybornika, bo adnotacja bez treści nie jest adnotacją.
   */
  async function dodajAdnotacje(): Promise<void> {
    przybornik.pokazOdpowiedz(
      'Adnotację zakłada się w przyborniku znakowania: wpisz jej treść i numer fragmentu różnicy, ' +
        'a potem naciśnij „Adnotacja przy fragmencie różnicy". Idzie komendą studio.annotation.add.',
      true,
    );
    if (!przybornik.widoczny()) przybornik.przestawWidocznosc();
  }

  /* ── Wstawianie w miejsce kursora ────────────────────────────────────────── */

  /**
   * Wstawia treść w miejsce kursora — jedna droga dla schowka, dyktowania
   * i wniesienia ze źródła.
   *
   * Zaznaczenie zostaje zastąpione, tak jak w każdym edytorze; brak zaznaczenia
   * znaczy „na końcu treści", bo kursor bez zaznaczenia stoi tam, gdzie stan
   * modułu ostatnio go widział. Druga kopia tego rachunku rozjechałaby się przy
   * pierwszej poprawce.
   */
  function wstawWTresc(tekst: string): void {
    if (tekst === '') {
      odpowiedz.pokaz('Nie ma czego wstawić — treść wstawiana jest pusta.', false);
      return;
    }
    const tresc = stan.trescRobocza();
    const zakres = stan.zaznaczenie();
    const poczatek = zakres?.poczatek ?? tresc.length;
    const koniec = zakres?.koniec ?? tresc.length;
    stan.ustawTresc(`${tresc.slice(0, poczatek)}${tekst}${tresc.slice(koniec)}`);
    stan.ustawZaznaczenie({ poczatek, koniec: poczatek + tekst.length });
    odpowiedz.pokaz(
      `Wstawiono ${tekst.length} znaków na znaku ${poczatek}. Zapis dokumentu zakłada wersję.`,
      true,
    );
  }

  /* ── Schowek ─────────────────────────────────────────────────────────────── */

  async function odczytajSchowek(fraza: string, tylkoPrzypiete: boolean): Promise<void> {
    const wynik = await przybornikZaplecze.schowekWykaz(fraza, tylkoPrzypiete, 50);
    if (!wynik.udany || wynik.wynik === undefined) {
      schowek.odmowa(opisOdmowy('Odczyt schowka', wynik.blad?.code, wynik.blad?.message));
      return;
    }
    schowek.pokaz(wynik.wynik.entries, wynik.wynik.total);
  }

  async function odlozDoSchowka(): Promise<void> {
    const zakres = stan.zaznaczenie();
    const tresc =
      zakres === null
        ? stan.trescRobocza()
        : stan.trescRobocza().slice(zakres.poczatek, zakres.koniec);
    if (tresc === '') {
      schowek.odmowa('Nie ma czego odłożyć — zaznaczenie i treść dokumentu są puste.');
      return;
    }
    const wynik = await przybornikZaplecze.schowekOdloz(
      tresc,
      ClipboardEntryKind.Text,
      stan.idOkna(),
    );
    if (!wynik.udany || wynik.wynik === undefined) {
      schowek.odmowa(opisOdmowy('Odłożenie do schowka', wynik.blad?.code, wynik.blad?.message));
      return;
    }
    odpowiedz.pokaz(
      wynik.wynik.alreadyPresent
        ? `Ta treść stała już w historii — rdzeń podniósł wpis zastany na czoło wykazu, nie ` +
            'założył drugiego. Powtórzenie nie jest błędem.'
        : `Odłożono ${tresc.length} znaków do historii schowka (wpis ${wynik.wynik.entry.id}).`,
      true,
    );
    await odczytajSchowek('', false);
  }

  async function przypnijWpisSchowka(idWpisu: string, przypiety: boolean): Promise<void> {
    const wynik = await przybornikZaplecze.schowekPrzypnij(idWpisu, przypiety);
    if (!wynik.udany) {
      schowek.odmowa(opisOdmowy('Przypięcie wpisu', wynik.blad?.code, wynik.blad?.message));
      return;
    }
    await odczytajSchowek('', false);
  }

  async function usunWpisSchowka(idWpisu: string): Promise<void> {
    const wynik = await przybornikZaplecze.schowekUsun(idWpisu);
    if (!wynik.udany || wynik.wynik === undefined) {
      schowek.odmowa(opisOdmowy('Usunięcie wpisu', wynik.blad?.code, wynik.blad?.message));
      return;
    }
    odpowiedz.pokaz(`Rdzeń usunął wpisów: ${wynik.wynik.deleted}. Przypięte zostają.`, true);
    await odczytajSchowek('', false);
  }

  /**
   * Postać akapitu pobrana malarzem formatów.
   *
   * Malarz kopiuje POSTAĆ, nie treść, więc nie jedzie schowkiem rdzenia:
   * `clipboard.*` niesie treść i rodzaj wpisu, a nie arkusz nastaw akapitu.
   * Postać czyta się z nastaw akapitu okna, bo dziś tam ona mieszka.
   */
  let postacMalarza: ReturnType<typeof akapity.dla> | null = null;

  function pobierzPostacAkapitu(): string | null {
    const numer = powierzchnia.blokKursora();
    if (numer < 0) return null;
    const nastawa = akapity.dla(numer);
    postacMalarza = { ...nastawa };
    return `wyrównanie ${nastawa.wyrownanie}, barwa ${nastawa.barwa === '' ? 'motywu' : nastawa.barwa}`;
  }

  function nalozPostacAkapitu(): void {
    const numer = powierzchnia.blokKursora();
    if (numer < 0 || postacMalarza === null) {
      odpowiedz.pokaz(
        'Malarz nakłada postać na akapit, w którym stoi kursor — postaw go w treści i najpierw ' +
          'pobierz postać z akapitu wzorcowego.',
        false,
      );
      return;
    }
    akapity.ustaw(numer, postacMalarza);
    powierzchnia.pokaz(stan.trescRobocza(), true);
    odpowiedz.pokaz('Postać akapitu nałożona; treści malarz nie ruszył.', true);
  }

  /* ── Operacje własne i przypięcia ────────────────────────────────────────── */

  async function zapiszOperacjeWlasna(
    nazwa: string,
    kategoria: string,
    polecenie: string,
  ): Promise<void> {
    if (nazwa === '' || kategoria === '' || polecenie === '') {
      katalog.pokazOdpowiedz(
        'Operacja własna wymaga nazwy, kategorii i treści polecenia — komenda ' +
          'studio.operation.save przyjmuje wszystkie trzy jako obowiązkowe.',
        false,
      );
      return;
    }
    const wynik = await przybornikZaplecze.przybornikZapiszOperacje('', nazwa, kategoria, polecenie);
    if (!wynik.udany || wynik.wynik === undefined) {
      katalog.pokazOdpowiedz(
        opisOdmowy('Zapis operacji własnej', wynik.blad?.code, wynik.blad?.message),
        false,
      );
      return;
    }
    katalog.pokazOdpowiedz(
      `Rdzeń zapisał operację „${wynik.wynik.operation.name}" (${wynik.wynik.operation.id}). ` +
        'Stoi w tym samym wykazie co fabryczne, w grupie operacji własnych.',
      true,
    );
    await wczytajOperacje();
  }

  async function usunOperacjeWlasna(idOperacji: string): Promise<void> {
    const wynik = await przybornikZaplecze.przybornikUsunOperacje(idOperacji);
    if (!wynik.udany || wynik.wynik === undefined) {
      katalog.pokazOdpowiedz(
        opisOdmowy('Usunięcie operacji', wynik.blad?.code, wynik.blad?.message),
        false,
      );
      return;
    }
    katalog.pokazOdpowiedz(
      wynik.wynik.deleted
        ? `Operacja ${idOperacji} usunięta.`
        : `Rdzeń nie usunął operacji ${idOperacji} i nie podał powodu w błędzie — pole deleted ` +
            'wróciło fałszem.',
      wynik.wynik.deleted,
    );
    await wczytajOperacje();
  }

  async function wczytajOperacje(): Promise<void> {
    const wynik = await przybornikZaplecze.przybornikOperacje();
    if (!wynik.udany || wynik.wynik === undefined) {
      katalog.pokazOdpowiedz(
        opisOdmowy('Odczyt operacji rdzenia', wynik.blad?.code, wynik.blad?.message),
        false,
      );
      return;
    }
    katalog.ustawOperacje(wynik.wynik.operations);
  }

  async function przypnijCzynnosc(idAkcji: string): Promise<void> {
    const przypiete = uzycie.przestawPrzypiecie(idAkcji);
    odswiezPlywak();
    const wynik = await przybornikZaplecze.przybornikZapisz(KLUCZ_PRZYPIETYCH, przypiete);
    if (!wynik.udany) {
      odpowiedz.pokaz(
        opisOdmowy('Zapis przypięcia czynności', wynik.blad?.code, wynik.blad?.message),
        false,
      );
    }
  }

  function ustawTrybOperacji(tryb: TrybOperacji): void {
    uzycie.ustawTryb(tryb);
    plywak.ustawTryb(tryb);
    czynnosciZewnetrzne.naTrybOperacji(tryb);
    void przybornikZaplecze.przybornikZapisz(KLUCZ_TRYBU, tryb);
  }

  function odswiezPlywak(): void {
    plywak.ustawCzynnosci(
      uzycie.naWierzchu(CZYNNOSCI_NA_WIERZCHU),
      (idAkcji) => uzycie.przypieta(idAkcji),
      uzycie.opiszKolejnosc(),
    );
    plywak.ustawSuwaki(nastawySuwakow);
  }

  /* ── Adnotacje przybornika ───────────────────────────────────────────────── */

  /**
   * Zakłada adnotację drogą WŁASNEJ komendy.
   *
   * Kontrakt ostrzega przy `studio.annotation.add`, że czynność ta „idzie dziś
   * drogą generyczną `window.action` i wraca odmowa `not_found`". Sprawdzone:
   * tak było i w tym oknie. Teraz jedzie komendą własną, więc odmowa — jeśli
   * przyjdzie — nazwie brakujący uchwyt komendy, a nie brakujący wiersz katalogu
   * akcji. To różnica, po której poznaje się, czego naprawdę brakuje.
   */
  async function dodajAdnotacjePrzybornika(numer: number, tresc: string): Promise<void> {
    const dokument = stan.dokument();
    if (dokument === null || tresc.trim() === '') {
      przybornik.pokazOdpowiedz(
        'Adnotacja wymaga wczytanego dokumentu i treści — komenda studio.annotation.add przyjmuje ' +
          'oba pola jako obowiązkowe.',
        false,
      );
      return;
    }
    const para = stan.paraPorownania();
    const wynik = await przybornikZaplecze.przybornikDodajAdnotacje(dokument.id, numer, tresc, {
      ...(para === null ? {} : { odniesienie: para.odniesienie, porownywana: para.porownywana }),
    });
    if (!wynik.udany || wynik.wynik === undefined) {
      przybornik.pokazOdpowiedz(
        opisOdmowy('Dodanie adnotacji', wynik.blad?.code, wynik.blad?.message),
        false,
      );
      return;
    }
    przybornik.pokazOdpowiedz(
      `Adnotacja założona przy fragmencie ${numer} (autor ` +
        `${wynik.wynik.annotation.author}).`,
      true,
    );
    await wczytajAdnotacje();
  }

  async function wczytajAdnotacje(): Promise<void> {
    const dokument = stan.dokument();
    if (dokument === null) {
      adnotacje = [];
      return;
    }
    const para = stan.paraPorownania();
    const wynik = await przybornikZaplecze.przybornikAdnotacje(dokument.id, {
      ...(para === null ? {} : { odniesienie: para.odniesienie, porownywana: para.porownywana }),
    });
    adnotacje = wynik.udany && wynik.wynik !== undefined ? wynik.wynik.annotations : [];
    odswiez();
  }

  /* ── Osadzenie: Biblioteka i strona sieciowa ─────────────────────────────── */

  async function szukajWBibliotece(fraza: string): Promise<void> {
    const wynik = await przybornikZaplecze.osadzenieBiblioteka(fraza, 30);
    if (!wynik.udany || wynik.wynik === undefined) {
      osadzenie.pokazOdpowiedz(
        opisOdmowy('Odczyt Biblioteki', wynik.blad?.code, wynik.blad?.message),
        false,
      );
      return;
    }
    osadzenie.pokazPliki(wynik.wynik.files);
    osadzenie.pokazOdpowiedz(
      `Biblioteka oddała plików: ${wynik.wynik.files.length} z ${wynik.wynik.total ?? wynik.wynik.files.length}.`,
      true,
    );
  }

  /** Plik ostatnio odczytany z Biblioteki — do wiersza pochodzenia. */
  let plikiBiblioteki: readonly { id: string; nazwa: string; wersja: string; suma: string }[] = [];

  async function podejrzyjPlik(idPliku: string, strona: number): Promise<void> {
    const wynikPlikow = await przybornikZaplecze.osadzenieBiblioteka('', 30);
    const plik =
      wynikPlikow.udany && wynikPlikow.wynik !== undefined
        ? wynikPlikow.wynik.files.find((pozycja) => pozycja.id === idPliku)
        : undefined;
    if (plik === undefined) {
      osadzenie.pokazOdpowiedz(
        `Biblioteka nie oddała pliku ${idPliku} przy odczycie wykazu, więc podglądu nie ma czego ` +
          'dotyczyć. Odczytaj wykaz ponownie.',
        false,
      );
      return;
    }
    plikiBiblioteki = [
      ...plikiBiblioteki.filter((pozycja) => pozycja.id !== plik.id),
      {
        id: plik.id,
        nazwa: plik.name,
        wersja: plik.versionId ?? '',
        suma: plik.checksum ?? '',
      },
    ];
    const wynik = await przybornikZaplecze.osadzeniePodglad(
      idPliku,
      strona,
      OSADZENIE_ZNAKOW_PODGLADU,
    );
    if (!wynik.udany || wynik.wynik === undefined) {
      osadzenie.pokazOdpowiedz(
        opisOdmowy('Podgląd pliku Biblioteki', wynik.blad?.code, wynik.blad?.message),
        false,
      );
      return;
    }
    osadzenie.pokazPodglad(plik, wynik.wynik.preview);
  }

  function wniesZBiblioteki(
    plik: { id: string; name: string; versionId?: string; checksum?: string },
    podgladPliku: { page?: number },
    tresc: string,
  ): void {
    const zakres = stan.zaznaczenie();
    const naZnaku = zakres?.poczatek ?? stan.trescRobocza().length;
    const zapis = pochodzenia.zapisz({
      zrodlo: ZrodloWniesienia.Biblioteka,
      nazwa: plik.name,
      wskazanie: plik.id,
      wersja: plik.versionId ?? '',
      sumaKontrolna: plik.checksum ?? '',
      strona: podgladPliku.page ?? 0,
      znakow: tresc.length,
      wstawionoNaZnaku: naZnaku,
    });
    wstawWTresc(`${tresc}\n${osadzenieWierszPochodzenia(zapis)}\n`);
    osadzenie.pokazPochodzenia(pochodzenia.wykaz());
    osadzenie.pokazOdpowiedz(
      `Wniesiono ${tresc.length} znaków z pliku „${plik.name}" wraz z wierszem pochodzenia.`,
      true,
    );
  }

  async function wciagnijStrone(adres: string, zObrazami: boolean): Promise<void> {
    if (adres === '') {
      osadzenie.pokazOdpowiedz('Podaj adres strony — bez niego nie ma czego pobrać.', false);
      return;
    }
    const wynik = await przybornikZaplecze.osadzenieWciagnijStrone(
      stan.idOkna(),
      adres,
      zObrazami,
    );
    if (!wynik.udany || wynik.wynik === undefined) {
      osadzenie.pokazOdpowiedz(
        opisOdmowy('Wciągnięcie strony', wynik.blad?.code, wynik.blad?.message),
        false,
      );
      return;
    }
    const pozycja = wynik.wynik.item;
    osadzenie.pokazStrone(
      adres,
      adres,
      pozycja.text ?? '',
      `studio.ingest.url · stan pozycji ${pozycja.state}` +
        (pozycja.usedOcr === true ? ' · tekst z rozpoznania pisma' : '') +
        (pozycja.failureReason === undefined ? '' : ` · rdzeń zgłasza: ${pozycja.failureReason}`),
    );
  }

  async function wezMigawke(idOknaPrzegladarki: string): Promise<void> {
    if (idOknaPrzegladarki === '') {
      osadzenie.pokazOdpowiedz(
        'Migawka wymaga okna modułu Browser — komenda browser.snapshot.get przyjmuje jego ' +
          'identyfikator jako pole obowiązkowe. Bez otwartej przeglądarki weź drogę drugą: ' +
          '„Wciągnij stronę oknem Studia".',
        false,
      );
      return;
    }
    const wynik = await przybornikZaplecze.osadzenieMigawka(idOknaPrzegladarki, false);
    if (!wynik.udany || wynik.wynik === undefined) {
      osadzenie.pokazOdpowiedz(
        opisOdmowy('Migawka strony', wynik.blad?.code, wynik.blad?.message),
        false,
      );
      return;
    }
    const migawkaStrony = wynik.wynik.snapshot;
    osadzenie.pokazStrone(
      migawkaStrony.title ?? migawkaStrony.url,
      migawkaStrony.url,
      migawkaStrony.text ?? '',
      `browser.snapshot.get · pobrano ${new Date(migawkaStrony.capturedAt).toLocaleString('pl-PL')}`,
    );
  }

  async function otworzStroneWPrzegladarce(
    idOknaPrzegladarki: string,
    adres: string,
  ): Promise<void> {
    if (idOknaPrzegladarki === '' || adres === '') {
      osadzenie.pokazOdpowiedz(
        'Otwarcie strony w przeglądarce wymaga okna modułu Browser i adresu — komenda ' +
          'browser.navigate przyjmuje oba jako obowiązkowe.',
        false,
      );
      return;
    }
    const wynik = await przybornikZaplecze.osadzenieOtworz(idOknaPrzegladarki, adres);
    if (!wynik.udany || wynik.wynik === undefined) {
      osadzenie.pokazOdpowiedz(
        opisOdmowy('Otwarcie strony', wynik.blad?.code, wynik.blad?.message),
        false,
      );
      return;
    }
    const migawkaStrony = wynik.wynik.snapshot;
    osadzenie.pokazStrone(
      migawkaStrony.title ?? migawkaStrony.url,
      migawkaStrony.url,
      migawkaStrony.text ?? '',
      'browser.navigate · strona otwarta w module Browser',
    );
  }

  function wniesZeStrony(tytul: string, adres: string, tresc: string): void {
    const zakres = stan.zaznaczenie();
    const naZnaku = zakres?.poczatek ?? stan.trescRobocza().length;
    const zapis = pochodzenia.zapisz({
      zrodlo: adres === tytul ? ZrodloWniesienia.Strona : ZrodloWniesienia.MigawkaPrzegladarki,
      nazwa: tytul,
      wskazanie: adres,
      wersja: '',
      sumaKontrolna: '',
      strona: 0,
      znakow: tresc.length,
      wstawionoNaZnaku: naZnaku,
    });
    wstawWTresc(`${tresc}\n${osadzenieWierszPochodzenia(zapis)}\n`);
    osadzenie.pokazPochodzenia(pochodzenia.wykaz());
    osadzenie.pokazOdpowiedz(
      `Wniesiono ${tresc.length} znaków ze strony „${tytul}" wraz z wierszem pochodzenia.`,
      true,
    );
  }

  /* ── Komentarze: skoki i wskazanie ───────────────────────────────────────── */

  function watki(): readonly StudioComment[] {
    return komentarze.filter((komentarz) => komentarz.parentCommentId === undefined);
  }

  function skoczDoKomentarza(wPrzod: boolean): void {
    const wykaz = watki();
    if (wykaz.length === 0) {
      odpowiedz.pokaz('Komentarzy nie ma — nie ma po czym skakać.', false);
      return;
    }
    wskazanyKomentarz = wPrzod
      ? (wskazanyKomentarz + 1) % wykaz.length
      : (wskazanyKomentarz - 1 + wykaz.length) % wykaz.length;
    const komentarz = wykaz[wskazanyKomentarz];
    if (komentarz !== undefined) wskazKomentarz(komentarz.id);
  }

  function wskazKomentarz(idKomentarza: string): void {
    const dymek = dymki.element.querySelector<HTMLElement>(`[data-komentarz='${idKomentarza}']`);
    dymek?.scrollIntoView({ block: 'nearest' });
    for (const inny of Array.from(
      dymki.element.querySelectorAll<HTMLElement>('[data-komentarz]'),
    )) {
      inny.dataset['wskazany'] = inny === dymek ? 'tak' : 'nie';
    }
  }

  /* ── Różnica: wyrys i decyzja wybiórcza ──────────────────────────────────── */

  function numeryFragmentow(): number[] {
    return wskazaneFragmenty.value
      .split(/[^\d]+/u)
      .filter((czesc) => czesc !== '')
      .map((czesc) => Number(czesc));
  }

  function przerysujRoznice(): void {
    const widoczne = przefiltrujRoznice(
      fragmentyOdpowiedzi,
      filtr.kontrolka.value as FiltrRoznicy,
    );
    wykazFragmentow.replaceChildren(
      ...wyrysFragmentow(widoczne),
      ...wyrysTrafien(trafieniaOdpowiedzi),
    );
    statystyka.textContent =
      fragmentyOdpowiedzi.length === 0
        ? ''
        : `${opiszStatystyke(policzRoznice(fragmentyOdpowiedzi))} Widocznych fragmentów: ${widoczne.length}.`;
  }

  /* ── Odświeżenie ─────────────────────────────────────────────────────────── */

  function odswiez(): void {
    status.odswiez();
    // Siedem powierzchni postaci i kontroli pracy odświeża się razem z oknem:
    // każda patrzy na TEN dokument, więc rozjazd między nimi byłby rozjazdem
    // o jednym bycie.
    postac.odswiez();
    void tabele.odswiez();
    void obiekty.odswiez();
    void aparat.odswiez();
    void szablony.odswiez();
    void dziennik.odswiez();
    void kopie.odswiez();
    void blokady.odswiez();
    zakresOperacji.textContent = opiszZakres();
    szukanie.odswiez();
    karty.odswiez();
    wstazka.odswiezSuwaki(nastawySuwakow);
    redaktor.odswiez(stan.trescRobocza());
    odswiezPlywak();
    przybornik.pokaz(
      przybornikZlozZnakowania({
        komentarze,
        zmiany,
        adnotacje,
        znaczniki: znaczniki.wykaz(),
        propozycja: stan.propozycja(),
        zakresPropozycji: stan.zaznaczenie(),
      }),
    );
    przybornik.ustawZakladki(zakladki);

    const propozycja = stan.propozycja();
    const odmowa = stan.odmowaOperacji();
    pochodzenie.textContent =
      propozycja !== null
        ? `Wynik operacji ${propozycja.idAkcji} · ${propozycja.tresc.length} znaków · propozycja: ` +
          `${propozycja.idPropozycji === '' ? 'bez odwołania w rdzeniu' : propozycja.idPropozycji}` +
          ` · ${opiszZmiany(zmiany)} · wątków otwartych ${dymki.ile()} · ${opiszNastawy(nastawyWizualne)}`
        : odmowa !== null
          ? odmowa
          : `${opiszZmiany(zmiany)} · wątków otwartych ${dymki.ile()} · ${powierzchnia.opis()} · ${opiszNastawy(nastawyWizualne)}`;
    decyzja.dataset['propozycja'] = propozycja === null ? 'brak' : 'oczekuje';

    // Treść powierzchni idzie za stanem, ale wyłącznie wtedy, gdy stan mówi co
    // innego niż powierzchnia — przerysowanie przy każdym naciśnięciu klawisza
    // zabierałoby kursor.
    const tresc = powierzchnia.tryb() === 'wydanie' ? stan.trescZaakceptowana() : stan.trescRobocza();
    powierzchnia.pokaz(tresc);

    const dokument = stan.dokument();
    if (dokument !== null && dokument.id !== dokumentOdpowiedzi && wykazFragmentow.childElementCount > 0) {
      // Odpowiedź różnicy dotyczyła innego dokumentu — jej wyrys i statystyka nie
      // mówią prawdy o tym, który stoi w oknie teraz.
      fragmentyOdpowiedzi = [];
      trafieniaOdpowiedzi = [];
      przerysujRoznice();
      powierzchnia.ustawFragmenty([]);
    }

    if (rama.stan.faza() === 'ladowanie' || rama.stan.faza() === 'blad') return;
    if (dokument === null) {
      rama.stan.puste(
        'Okno pracy bez wczytanego dokumentu',
        'To jedno okno prowadzi cały dokument: treść z formatowaniem na kartce, podgląd wydania ' +
          'i różnicę jako tryby widoku, zmiany modelu w miejscu, komentarze na marginesie. ' +
          'Wskaż plik Library, dokument tej sesji albo ścieżkę na urządzeniu i naciśnij ' +
          '„Wczytaj dokument" — albo załóż nowy z galerii szablonów.',
      );
      return;
    }
    rama.stan.gotowe();
  }

  return {
    element: rama.element,

    async wczytaj() {
      // Postać dokumentu idzie pierwsza: arkusz stylów, nastawy strony i sekcje
      // są tym, wobec czego liczą się wszystkie pozostałe odczyty okna.
      await postac.wczytaj();
      const szablony = await praca.szablony();
      if (szablony.udany && szablony.wynik !== undefined) {
        galeria.ustawSzablony(szablony.wynik.templates);
      }
      await wczytajProfile();

      // Nastawy przybornika idą przed wykazem operacji: kolejność czynności na
      // wierzchu pływaka liczy się z użycia, a użycie leży w nastawach.
      const nastawy = await przybornikZaplecze.przybornikNastawy();
      if (nastawy.udany && nastawy.wynik !== undefined) {
        uzycie.wchlon(nastawy.wynik);
        plywak.ustawTryb(nastawy.wynik.tryb);
        czynnosciZewnetrzne.naTrybOperacji(nastawy.wynik.tryb);
      }
      odswiezPlywak();

      // Rejestr akcji rdzenia zasila katalog narzędzi ukrytych tym samym
      // odczytem, którym żywi się stały panel — wykaz jest jeden, nie dwa.
      const rejestr = await akcje.katalog();
      if (rejestr.udany && rejestr.wynik !== undefined) katalog.ustawRejestr(rejestr.wynik.actions);
      await wczytajOperacje();
      await wczytajAdnotacje();
      await odswiezZRdzenia();
    },

    ustawTrybOperacji,
    odswiez,
  };
}

/**
 * Treść wpisu schowka zdjęta z postaci — droga „wklej jako czysty tekst".
 *
 * Postać w treści tego edytora niosą znaczniki markdown (`znaczniki-markdown.ts`)
 * oraz znacznik podziału strony. Czyszczenie zdejmuje właśnie je, a nie same
 * litery: wklejenie czyste ma dać brzmienie bez formatu, nie brzmienie okrojone.
 * Oba warianty wklejenia są równorzędne i wybór należy do Operatora.
 */
function schowekTekstCzysty(tresc: string): string {
  return tresc
    .split('\n')
    .map((wiersz) =>
      wiersz
        .replace(/^\s{0,3}#{1,6}\s+/u, '')
        .replace(/^\s{0,3}>\s?/u, '')
        .replace(/^\s{0,3}[-*+]\s+/u, '')
        .replace(/\*\*(.+?)\*\*/gu, '$1')
        .replace(/__(.+?)__/gu, '$1')
        .replace(/(?<!\*)\*(?!\*)(.+?)(?<!\*)\*(?!\*)/gu, '$1')
        .replace(/`([^`]+)`/gu, '$1'),
    )
    .filter((wiersz) => wiersz.trim() !== ZNACZNIK_PODZIALU)
    .join('\n');
}
