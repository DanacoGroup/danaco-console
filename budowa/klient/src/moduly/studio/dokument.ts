/**
 * Moduł Studio — kanwa dokumentu, strefa „Studio Editor" ze źródła kształtu
 * (`design/05-okna/moduly/studio.html`).
 *
 * Kanwa składa treść jak dokument — tytuł, nagłówki, akapity — i wypełnia całą
 * wysokość oraz szerokość okna: dokument czyta się kolumną, nie prostokątem
 * pośrodku pustki. Edycja idzie wprost w tych blokach, bez pola formularza.
 *
 * Rdzeń dotykany czterema drogami: `studio.text.get` czyta treść,
 * `studio.text.edit` odkłada zmieniony fragment przy odejściu od kanwy,
 * `studio.document.save` zapisuje tytuł i zakłada wersję, a
 * `studio.document.open` czyta dokument z powrotem — dopiero treść wrócona
 * z rdzenia dowodzi zapisu.
 *
 * Okno modułu jest bytem rdzenia, nie widoku — bez niego żadna komenda Studia
 * dotykająca dokumentu nie ma gdzie stanąć, bo wszystkie wymagają `windowId`.
 */
import {
  Command,
  ExecutionEnv,
  PermissionMode,
  StudioDocumentFormat,
  WindowRole,
  type ErrorInfo,
  type Module,
  type Session,
  type StudioDocument,
  type StudioDocumentSaveRequest,
} from '../../../../shared/contract.ts';
import type { Kanal } from '../../protokol/kanal.ts';
import { wywolaj } from '../../protokol/wywolanie.ts';
import { ikony } from './ikony.ts';
import { zLiczba } from './liczebnik.ts';
import { el, tekst, zeZnacznika } from './narzedzia.ts';
import { opisOdmowy, zaloguj } from './odmowa.ts';
import { pasekEdytora } from './skladniki/pasek-edytora.ts';

export interface NastawyDokumentu {
  /** Kanał, którym panel woła komendy rdzenia. */
  kanal: Kanal;
  /** Karty sesji odtworzone przez rdzeń przy wejściu do środowiska. */
  sesje: Session[];
  /** Moduł, dla którego okno powstaje — jego `id` idzie do `window.create`. */
  modul: Module;
}

export interface PanelDokumentu {
  /** Węzeł panelu, gotowy do wstawienia w bryłę okna. */
  wezel: HTMLElement;
  /** Zdejmuje panel; wynik wywołania w locie ląduje w nicości. */
  zdejmij(): void;
  /* Okno i sesja powstają tu, bo bez nich żadna komenda Studia nie ma gdzie stanąć.
     Pozostałe panele czytają je stąd — odczytem, nie kopią, bo w chwili montażu
     okna jeszcze nie ma: zakładanie jest wywołaniem do rdzenia. */
  idOkna(): string | null;
  idSesji(): string | null;
}

type StanZapisu = 'spoczynek' | 'zapisuje' | 'zapisany';

type Stan =
  | { rodzaj: 'zakladanie' }
  | { rodzaj: 'dokument'; dokument: StudioDocument; tresc: string; zapis: StanZapisu }
  | { rodzaj: 'odmowa'; powod: string; blad?: ErrorInfo };

/* Katalogi robocze okna. Kontrakt wymaga pola, nie zawartości, a źródła nie
   podają katalogu należnego oknu modułu przed wskazaniem Operatora: wykaz
   zostaje pusty, a okno bez katalogu nie sięga plików. */
const KATALOGI_ROBOCZE: string[] = [];

/* Bloki dokumentu rozdziela pusty wiersz, a nagłówki w zapisie Markdown niosą
   znacznik kratki. Rozbiór i złożenie idą tą samą miarą, więc treść wraca do
   rdzenia w tej postaci, w jakiej z niego przyszła. */
const ROZDZIELNIK = '\n\n';
const ZNACZNIK_H1 = '# ';
const ZNACZNIK_H2 = '## ';

export function panelDokumentu(w: NastawyDokumentu): PanelDokumentu {
  let zdjete = false;
  let okno = '';
  let sesjaBiezaca: string | null = null;
  let stanBiezacy: Stan = { rodzaj: 'zakladanie' };
  let dokumentBiezacy: StudioDocument | null = null;
  /* Treść znana rdzeniowi. Zmiana idzie do rdzenia jako różnica wobec niej,
     więc bez niej `studio.text.edit` nie miałby czym wskazać fragmentu. */
  let trescZnana = '';
  /* Czynności rdzenia idą jedna po drugiej: odejście od kanwy zapisuje fragment,
     a klik w zapis wersji pada zaraz po nim — równolegle jedna z nich by przepadła. */
  let kolejka: Promise<void> = Promise.resolve();

  /* Szerokość stoi stylem, bo margines automatyczny `dn-kanwa` odbiera
     składnikowi kolumny rozciągnięcie: kanwa kurczyła się do napisu i stawała
     prostokątem pośrodku pustki. Miara przy 2560 px: 625 × 1097 px. */
  const tresc = el('div', {
    klasa: 'sta-okno-tresc dn-kanwa',
    style: 'display: flex; flex-direction: column; width: 100%',
  });

  /* Korpus i tytuł powstają dopiero z odpowiedzią rdzenia, więc pasek
     narzędziowy i zapis czytają je stąd przy każdym wywołaniu, a nie z kopii
     zrobionej w chwili montażu. */
  let korpus: HTMLElement | null = null;
  let poleTytulu: HTMLElement | null = null;

  const uwaga = el('span', { klasa: 'dn-meta', role: 'status' });

  /** Dokłada czynność do kolejki; odmowa jednej nie zatrzymuje następnych. */
  function wKolejce(czynnosc: () => Promise<void>): void {
    kolejka = kolejka.then(czynnosc).catch((powod: unknown) => {
      console.warn('[studio] czynność kanwy przerwana', powod);
    });
  }

  const pasek = pasekEdytora({
    kanal: w.kanal,
    idDokumentu: () => dokumentBiezacy?.id ?? null,
    zaznaczenie,
    powiadom(zdanie) {
      uwaga.textContent = zdanie;
    },
  });

  /** Przycisk znakowy belki okna: znak niesie rysunek, nazwę czynności etykieta. */
  function przyciskBelki(rysunek: string, etykieta: string): HTMLButtonElement {
    const rysunekWezel = zeZnacznika(rysunek);
    rysunekWezel.setAttribute('aria-hidden', 'true');
    return el('button', {
      klasa: 'dn-btn-ikona',
      type: 'button',
      'aria-label': etykieta,
      title: etykieta,
    }, [rysunekWezel]) as HTMLButtonElement;
  }

  const przyciskZapisu = przyciskBelki(ikony.zapisz, tekst('dokument.zapisz'));
  przyciskZapisu.addEventListener('click', () => wKolejce(() => zapiszDokument(true)));

  /* Konfiguracja okna stoi w prototypie i nie ma pokrycia w kontrakcie: żadna
     komenda Studia nie przyjmuje nastaw okna wiodącego. Przycisk zostaje
     i nazywa niegotowość, zamiast zniknąć albo udać działanie. */
  const przyciskKonfiguracji = przyciskBelki(ikony.dostosuj, tekst('dokument.konfiguracja'));
  przyciskKonfiguracji.addEventListener('click', () => {
    uwaga.textContent = tekst('dokument.zapowiedziane');
  });

  const znakTytulu = zeZnacznika(ikony.olowek);
  znakTytulu.setAttribute('aria-hidden', 'true');

  const pasStanu = el('div', { klasa: 'st-status' });

  const wezel = el(
    'section',
    { klasa: 'sta-okno', id: 'panel-editor', role: 'tabpanel', 'aria-labelledby': 'karta-editor' },
    [
      el('header', { klasa: 'sta-okno-belka' }, [
        el('span', { klasa: 'sta-okno-tytul' }, [
          znakTytulu,
          el('b', { tekst: tekst('dokument.tytul') }),
        ]),
        el('span', { klasa: 'dn-plakietka dn-plakietka--rola', tekst: tekst('dokument.rolaOkna') }),
        el('span', { klasa: 'sta-okno-akcje' }, [przyciskZapisu, przyciskKonfiguracji]),
      ]),
      pasek.wezel,
      tresc,
      pasStanu,
    ],
  );

  odswiez({ rodzaj: 'zakladanie' });
  void zaloz();

  /** Sesja pod okno modułu: pierwsza odtworzona przez rdzeń, a przy ich braku nowa. */
  async function sesja(): Promise<string | null> {
    const odtworzona = w.sesje[0];
    if (odtworzona !== undefined) return odtworzona.id;
    const wynik = await wywolaj(w.kanal, Command.SessionCreate, {});
    return wynik.udany ? (wynik.wynik?.session.id ?? null) : null;
  }

  /** Kanał modelu wymagany przez `window.create`: pierwszy czynny z wykazu rdzenia. */
  async function kanalModelu(): Promise<string | null> {
    const wynik = await wywolaj(w.kanal, Command.ChannelList, {});
    if (!wynik.udany) return null;
    const czynny = (wynik.wynik?.channels ?? []).find((k) => k.enabled);
    return czynny?.id ?? null;
  }

  async function zaloz(): Promise<void> {
    const idSesji = await sesja();
    if (zdjete) return;
    if (idSesji === null) return odswiez({ rodzaj: 'odmowa', powod: 'sesja' });
    /* Sesję trzymamy od razu: pozostałe panele okna biorą ją stąd, a wchodzą
       do okna później niż odpowiedź rdzenia na ten odczyt. */
    sesjaBiezaca = idSesji;

    const idKanalu = await kanalModelu();
    if (zdjete) return;
    if (idKanalu === null) return odswiez({ rodzaj: 'odmowa', powod: 'kanal' });

    const oknoWynik = await wywolaj(w.kanal, Command.WindowCreate, {
      sessionId: idSesji,
      moduleId: w.modul.id,
      modelChannelId: idKanalu,
      workingDirs: KATALOGI_ROBOCZE,
      executionEnv: ExecutionEnv.Local,
      permissionMode: PermissionMode.Manual,
      windowRole: WindowRole.Standalone,
    });
    if (zdjete) return;
    if (!oknoWynik.udany || oknoWynik.wynik === undefined) {
      return odswiez({ rodzaj: 'odmowa', powod: 'okno', blad: oknoWynik.blad });
    }
    okno = oknoWynik.wynik.window.id;

    const dokumentWynik = await wywolaj(w.kanal, Command.StudioDocumentCreate, {
      windowId: okno,
      title: tekst('dokument.nazwaNowego'),
    });
    if (zdjete) return;
    if (!dokumentWynik.udany || dokumentWynik.wynik === undefined) {
      return odswiez({ rodzaj: 'odmowa', powod: 'dokument', blad: dokumentWynik.blad });
    }
    await wczytaj(dokumentWynik.wynik.document.id);
  }

  /**
   * Wczytuje dokument do kanwy: `studio.document.open` oddaje sam dokument,
   * `studio.text.get` jego treść. Gdy rdzeń treści nie odda, zostaje treść
   * niesiona przez dokument — obie pochodzą z rdzenia, żadna stąd.
   */
  async function wczytaj(idDokumentu: string): Promise<void> {
    const otwarcie = await wywolaj(w.kanal, Command.StudioDocumentOpen, {
      windowId: okno,
      documentId: idDokumentu,
    });
    if (zdjete) return;
    if (!otwarcie.udany || otwarcie.wynik === undefined) {
      return odswiez({ rodzaj: 'odmowa', powod: 'otwarcie', blad: otwarcie.blad });
    }
    const dokument = otwarcie.wynik.document;

    const odczyt = await wywolaj(w.kanal, Command.StudioTextGet, { documentId: dokument.id });
    if (zdjete) return;
    if (!odczyt.udany) zaloguj(odczyt.blad, 'odczytTresci');
    const trescRdzenia = odczyt.udany ? (odczyt.wynik?.text ?? '') : (dokument.content ?? '');

    odswiez({ rodzaj: 'dokument', dokument, tresc: trescRdzenia, zapis: 'spoczynek' });
    if (trescRdzenia === '') uwaga.textContent = tekst('dokument.pusty');
  }

  /**
   * Odkłada w rdzeniu sam zmieniony fragment (`studio.text.edit`), a potem
   * czyta treść z powrotem. Kanwa nie jest przerysowywana: Operator właśnie
   * w niej pisze, a wymiana węzłów zabrałaby mu miejsce karetki.
   */
  async function zapiszZmiane(): Promise<void> {
    if (dokumentBiezacy === null || korpus === null) return;
    const biezaca = zTresci();
    const zmiana = roznica(trescZnana, biezaca);
    if (zmiana === null) return;

    odswiezPas(biezaca, 'zapisuje');
    const zapis = await wywolaj(w.kanal, Command.StudioTextEdit, {
      documentId: dokumentBiezacy.id,
      rangeStart: zmiana.poczatek,
      rangeEnd: zmiana.koniec,
      text: zmiana.tekst,
      keepFormat: true,
    });
    if (zdjete) return;
    if (!zapis.udany) {
      uwaga.textContent = `${tekst('dokument.odmowaZmiany')} ${opisOdmowy(zapis.blad, 'zmianaTresci')}`;
      return odswiezPas(biezaca, 'spoczynek');
    }

    const odczyt = await wywolaj(w.kanal, Command.StudioTextGet, { documentId: dokumentBiezacy.id });
    if (zdjete) return;
    if (!odczyt.udany) zaloguj(odczyt.blad, 'odczytTresci');
    trescZnana = odczyt.udany ? (odczyt.wynik?.text ?? '') : biezaca;
    uwaga.textContent = '';
    odswiezPas(trescZnana, 'zapisany');
  }

  /**
   * Zapisuje cały dokument wraz z tytułem (`studio.document.save`). Wersja
   * zakładana jest wyłącznie na żądanie Operatora — zapis samego tytułu przy
   * odejściu od niego nie ma po co dokładać wersji do repozytorium sesji.
   */
  async function zapiszDokument(zWersja: boolean): Promise<void> {
    if (dokumentBiezacy === null || korpus === null) return;
    const dokument = dokumentBiezacy;
    const biezaca = zTresci();
    const tytul = poleTytulu?.textContent?.trim() ?? '';
    /* Odejście od tytułu bez jego zmiany nie ma czego zapisywać. */
    if (!zWersja && tytul === (dokument.title ?? '')) return;

    const zadanie: StudioDocumentSaveRequest = {
      documentId: dokument.id,
      content: biezaca,
      createVersion: zWersja,
    };
    if (tytul !== '') zadanie.title = tytul;

    przyciskZapisu.disabled = true;
    odswiezPas(biezaca, 'zapisuje');
    const zapis = await wywolaj(w.kanal, Command.StudioDocumentSave, zadanie);
    przyciskZapisu.disabled = false;
    if (zdjete) return;
    if (!zapis.udany || zapis.wynik === undefined) {
      uwaga.textContent = `${tekst('dokument.odmowaZapisu')} ${opisOdmowy(zapis.blad, 'zapisDokumentu')}`;
      return odswiezPas(biezaca, 'spoczynek');
    }
    uwaga.textContent = '';

    if (!zWersja) {
      /* Zapis tytułu zostawia kanwę nietkniętą: Operator przeszedł z tytułu do
         treści i pisze dalej, a przerysowanie zabrałoby mu miejsce karetki. */
      dokumentBiezacy = zapis.wynik.document;
      trescZnana = biezaca;
      stanBiezacy = { rodzaj: 'dokument', dokument: zapis.wynik.document, tresc: biezaca, zapis: 'zapisany' };
      return odswiezPas(biezaca, 'zapisany');
    }

    /* Odczyt po zapisie: dopiero treść wrócona z rdzenia dowodzi, że przeszła
       przez repozytorium sesji, a nie została w kanwie. */
    const odczyt = await wywolaj(w.kanal, Command.StudioDocumentOpen, {
      windowId: okno,
      documentId: dokument.id,
    });
    if (zdjete) return;
    if (!odczyt.udany || odczyt.wynik === undefined) {
      uwaga.textContent = `${tekst('dokument.odmowaZapisu')} ${opisOdmowy(odczyt.blad, 'odczytPoZapisie')}`;
      return odswiezPas(biezaca, 'spoczynek');
    }
    const potwierdzony = odczyt.wynik.document;
    const trescPotwierdzona = await trescZRdzenia(potwierdzony);
    if (zdjete) return;
    odswiez({ rodzaj: 'dokument', dokument: potwierdzony, tresc: trescPotwierdzona, zapis: 'zapisany' });
  }

  /** Treść dokumentu prosto z rdzenia; przy odmowie odczytu — treść samego dokumentu. */
  async function trescZRdzenia(dokument: StudioDocument): Promise<string> {
    const odczyt = await wywolaj(w.kanal, Command.StudioTextGet, { documentId: dokument.id });
    if (!odczyt.udany) {
      zaloguj(odczyt.blad, 'odczytTresci');
      return dokument.content ?? '';
    }
    return odczyt.wynik?.text ?? '';
  }

  /** Czy dokument jest zapisany w Markdown — wtedy nagłówek niesie znacznik kratki. */
  function czyMarkdown(): boolean {
    return dokumentBiezacy?.format === StudioDocumentFormat.Markdown;
  }

  /** Znacznik nagłówka należny blokowi; akapit nie niesie żadnego. */
  function znacznikBloku(blok: Node): string {
    if (!czyMarkdown()) return '';
    const nazwa = (blok as Element).tagName;
    if (nazwa === 'H1') return ZNACZNIK_H1;
    if (nazwa === 'H2') return ZNACZNIK_H2;
    return '';
  }

  /** Treść kanwy złożona z powrotem w postać, którą rdzeń wydał. */
  function zTresci(): string {
    if (korpus === null) return trescZnana;
    const kawalki: string[] = [];
    for (const blok of Array.from(korpus.childNodes)) {
      if (blok.nodeType === Node.TEXT_NODE) {
        const goly = blok.textContent ?? '';
        if (goly.trim() !== '') kawalki.push(goly);
        continue;
      }
      kawalki.push(znacznikBloku(blok) + (blok.textContent ?? ''));
    }
    return kawalki.join(ROZDZIELNIK);
  }

  /** Bloki dokumentu: nagłówki i akapity kanwy, nie wiersze pola tekstowego. */
  function bloki(trescDokumentu: string): HTMLElement[] {
    const markdown = czyMarkdown();
    return trescDokumentu.split(ROZDZIELNIK).map((kawalek) => {
      if (markdown && kawalek.startsWith(ZNACZNIK_H2)) {
        return el('h2', { tekst: kawalek.slice(ZNACZNIK_H2.length) });
      }
      if (markdown && kawalek.startsWith(ZNACZNIK_H1)) {
        return el('h1', { tekst: kawalek.slice(ZNACZNIK_H1.length) });
      }
      /* Akapit pusty dostaje złamanie wiersza: bez niego nie ma wysokości,
         więc Operator nie ma gdzie postawić karetki w dokumencie pustym. */
      return kawalek === '' ? el('p', {}, [el('br')]) : el('p', { tekst: kawalek });
    });
  }

  /**
   * Przesunięcie punktu kanwy w znakach treści dokumentu. Rdzeń przyjmuje
   * zakresy znakami całej treści, a kanwa dzieli ją na bloki — pusty wiersz
   * między blokami i znacznik nagłówka też są jej znakami.
   */
  function przesuniecie(wezel: Node, wOffsecie: number): number | null {
    if (korpus === null) return null;
    let suma = 0;
    for (const blok of Array.from(korpus.childNodes)) {
      if (blok === wezel || blok.contains(wezel)) {
        const zakres = document.createRange();
        zakres.setStart(blok, 0);
        zakres.setEnd(wezel, wOffsecie);
        return suma + znacznikBloku(blok).length + zakres.toString().length;
      }
      suma += znacznikBloku(blok).length + (blok.textContent ?? '').length + ROZDZIELNIK.length;
    }
    return null;
  }

  /** Zaznaczenie w kanwie podane paskowi narzędziowemu jako zakres znaków rdzenia. */
  function zaznaczenie(): { poczatek: number; koniec: number } | null {
    if (korpus === null) return null;
    const zaznaczone = window.getSelection();
    if (zaznaczone === null || zaznaczone.rangeCount === 0 || zaznaczone.isCollapsed) return null;
    const zakres = zaznaczone.getRangeAt(0);
    if (!korpus.contains(zakres.commonAncestorContainer)) return null;
    const poczatek = przesuniecie(zakres.startContainer, zakres.startOffset);
    const koniec = przesuniecie(zakres.endContainer, zakres.endOffset);
    if (poczatek === null || koniec === null || poczatek === koniec) return null;
    return { poczatek, koniec };
  }

  /**
   * Różnica dwóch postaci treści jako jeden ciągły fragment: wspólny początek
   * i wspólny koniec odpadają, zostaje to, co Operator naprawdę zmienił.
   * Bez tego zapis fragmentu przepisywałby cały dokument.
   */
  function roznica(
    stara: string,
    nowa: string,
  ): { poczatek: number; koniec: number; tekst: string } | null {
    if (stara === nowa) return null;
    const krotsza = Math.min(stara.length, nowa.length);
    let przod = 0;
    while (przod < krotsza && stara[przod] === nowa[przod]) przod += 1;
    let tyl = 0;
    while (tyl < krotsza - przod && stara[stara.length - 1 - tyl] === nowa[nowa.length - 1 - tyl]) {
      tyl += 1;
    }
    return {
      poczatek: przod,
      koniec: stara.length - tyl,
      tekst: nowa.slice(przod, nowa.length - tyl),
    };
  }

  function widokZakladania(): HTMLElement[] {
    return [
      el('div', { klasa: 'dn-wykaz-modulu-poz' }, [
        el('span', {
          klasa: 'dn-kropka dn-kropka--sygnal dn-kropka--tetno',
          'aria-hidden': 'true',
        }),
        tekst('dokument.zakladanie'),
      ]),
    ];
  }

  function widokOdmowy(powod: string, blad?: ErrorInfo): HTMLElement[] {
    return [
      el('div', { klasa: 'dn-alert dn-alert--wstega dn-alert--blad', role: 'alert' }, [
        el('span', { klasa: 'dn-alert-tresc' }, [
          el('b', { tekst: tekst(`dokumentOdmowa.${powod}`) }),
          el('span', { tekst: opisOdmowy(blad, 'dokument') }),
        ]),
      ]),
    ];
  }

  function widokDokumentu(dokument: StudioDocument, trescDokumentu: string): HTMLElement[] {
    const tytul = el('h1', {
      contenteditable: 'true',
      'aria-label': tekst('dokument.etykietaTytulu'),
      tekst: dokument.title ?? tekst('dokument.nazwaNowego'),
    });
    /* Tytuł jest polem rdzenia (`title` w `studio.document.save`), więc odejście
       od niego zapisuje go od razu — Operator nie ma osobnego przycisku, którego
       prototyp w tej strefie nie niesie. */
    tytul.addEventListener('blur', () => wKolejce(() => zapiszDokument(false)));
    tytul.addEventListener('keydown', (zdarzenie) => {
      if (zdarzenie.key === 'Enter') {
        zdarzenie.preventDefault();
        korpus?.focus();
      }
    });
    poleTytulu = tytul;

    /* Korpus bierze całą wolną wysokość kanwy klasą `dn-pole--rosnace`:
       dokument ma wypełniać okno, a nie stać prostokątem pośrodku pustki. */
    const korpusNowy = el('div', {
      klasa: 'dn-pole--rosnace',
      contenteditable: 'true',
      'aria-label': tekst('dokument.etykietaTresci'),
    }, bloki(trescDokumentu));

    korpusNowy.addEventListener('input', () => {
      odswiezPas(zTresci(), 'spoczynek');
    });
    korpusNowy.addEventListener('blur', () => wKolejce(zapiszZmiane));
    /* Wklejenie idzie samym tekstem: znaczniki z obcego źródła wniosłyby do
       kanwy wygląd spoza arkusza kształtu, a do rdzenia — treść nie swoją. */
    korpusNowy.addEventListener('paste', (zdarzenie) => {
      const wklejane = zdarzenie.clipboardData?.getData('text/plain');
      if (wklejane === undefined) return;
      zdarzenie.preventDefault();
      document.execCommand('insertText', false, wklejane);
    });
    korpus = korpusNowy;

    return [tytul, korpusNowy];
  }

  function zawartosc(stan: Stan): HTMLElement[] {
    switch (stan.rodzaj) {
      case 'zakladanie':
        return widokZakladania();
      case 'odmowa':
        return widokOdmowy(stan.powod, stan.blad);
      case 'dokument':
        return widokDokumentu(stan.dokument, stan.tresc);
    }
  }

  /*
  slowa liczy wyrazy treści — miara pasa stanu.

  Liczenie idzie po ciągach niebiałych znaków, nie po spacjach: tekst z dwoma
  odstępami pod rząd dałby wyraz pusty, a tekst pusty dałby jeden wyraz.
  */
  function slowa(trescDokumentu: string): number {
    const wyrazy = trescDokumentu.trim();
    return wyrazy === '' ? 0 : wyrazy.split(/\s+/u).length;
  }

  /** Zdanie o zapisie: pas stanu mówi, czy praca Operatora jest odłożona. */
  function zdanieZapisu(zapis: StanZapisu): string {
    if (zapis === 'zapisuje') return tekst('stanEdytora.zapisywanie');
    return zapis === 'zapisany' ? tekst('stanEdytora.zapisany') : tekst('stanEdytora.niezapisany');
  }

  /*
  odswiezStan wypełnia pas stanu wartościami dokumentu.

  Pas stoi poza `tresc`, więc przeżywa przerysowanie widoku — i musi, bo niesie
  uwagę paska narzędziowego, która powstaje niezależnie od stanu dokumentu.
  */
  function odswiezStan(stan: Stan): void {
    if (stan.rodzaj !== 'dokument') {
      pasStanu.replaceChildren();
      return;
    }
    const wersja = stan.dokument.versionId ?? null;
    pasStanu.replaceChildren(
      el('span', {}, [
        stan.zapis === 'zapisuje'
          ? el('span', { klasa: 'pt-tetno', 'aria-hidden': 'true' })
          : null,
        zdanieZapisu(stan.zapis),
      ]),
      el('span', {
        tekst: zLiczba(slowa(stan.tresc), {
          jedna: tekst('stanEdytora.slowoJedna'),
          kilka: tekst('stanEdytora.slowoKilka'),
          wiele: tekst('stanEdytora.slowoWiele'),
        }),
      }),
      el('span', {
        tekst: wersja === null
          ? tekst('stanEdytora.bezWersji')
          : `${tekst('stanEdytora.wersja')} ${wersja}`,
      }),
      uwaga,
      el('span', { klasa: 'st-status-prawa', tekst: tekst('stanEdytora.postac') }),
    );
  }

  /** Sam pas stanu, bez przerysowania kanwy — pisanie Operatora zostaje nietknięte. */
  function odswiezPas(trescBiezaca: string, zapis: StanZapisu): void {
    if (stanBiezacy.rodzaj !== 'dokument') return;
    stanBiezacy = { rodzaj: 'dokument', dokument: stanBiezacy.dokument, tresc: trescBiezaca, zapis };
    odswiezStan(stanBiezacy);
  }

  function odswiez(stan: Stan): void {
    stanBiezacy = stan;
    dokumentBiezacy = stan.rodzaj === 'dokument' ? stan.dokument : null;
    trescZnana = stan.rodzaj === 'dokument' ? stan.tresc : '';
    if (stan.rodzaj !== 'dokument') {
      korpus = null;
      poleTytulu = null;
    }
    przyciskZapisu.disabled = stan.rodzaj !== 'dokument';
    tresc.replaceChildren(...zawartosc(stan));
    odswiezStan(stan);
  }

  return {
    idOkna: () => okno,
    idSesji: () => sesjaBiezaca,
    wezel,
    zdejmij() {
      zdjete = true;
      if (wezel.isConnected) wezel.remove();
    },
  };
}
