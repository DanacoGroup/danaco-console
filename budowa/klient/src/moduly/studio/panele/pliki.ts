/**
 * Panel Pliki okna roboczego Studia. Kształt wzięty z
 * `design/05-okna/moduly/studio.html`, `section#panel-pliki`: belka z nazwą,
 * pole filtrowania `dn-szukaj` i wykaz `st-panel-lista` z wierszami
 * `st-panel-wiersz`. Treść przykładowa prototypu nie przechodzi — każdy
 * wiersz pochodzi z odpowiedzi rdzenia albo stoi nazwanym stanem pustym.
 *
 * Panel prowadzi materiał od wskazania źródła po pozycję w wykazie: dołożenie
 * ze ścieżki, pobranie strony, skan urządzeniem, rozpoznanie pisma, poprawa
 * rozpoznanych słów, przyjęcie wyniku jako dokument, a obok tego wniesienie
 * pliku i PDF wprost do edytora.
 *
 * Kolejka jest własnością okna modułu (`windowId`); bez niego nie ma gdzie
 * osadzić pozycji, więc sekcja kolejki niesie nazwany stan pusty zamiast
 * formularzy. Wykaz urządzeń od okna nie zależy i wczytuje się zawsze.
 *
 * `studio.ingest.changed` odświeża pozycję bez odpytywania w pętli — jedyne
 * odpytanie ponowne idzie po własnej czynności Operatora.
 */

import {
  Command,
  EventType,
  StudioIngestState,
  StudioInputDeviceKind,
  type ErrorInfo,
  type StudioImportBalance,
  type StudioIngestItem,
  type StudioInputDevice,
  type StudioRecognizedWord,
} from '../../../../../shared/contract.ts';
import { wywolaj } from '../../../protokol/wywolanie.ts';
import { el, zeZnacznika } from '../narzedzia.ts';
import { ikony, type NazwaZnaku } from '../ikony.ts';
import { tresciPliki } from './pliki-tresci.ts';
import type { MontazPanelu, ZamontowanyPanel } from './umowa.ts';
import { opisOdmowy } from '../odmowa.ts';
import { zLiczba } from '../liczebnik.ts';

type StanUrzadzen =
  | { rodzaj: 'ladowanie' }
  | { rodzaj: 'wykaz'; urzadzenia: StudioInputDevice[] }
  | { rodzaj: 'odmowa'; blad?: ErrorInfo };

type StanKolejki =
  | { rodzaj: 'brakOkna' }
  | { rodzaj: 'ladowanie' }
  | { rodzaj: 'wykaz'; pozycje: StudioIngestItem[] }
  | { rodzaj: 'odmowa'; blad?: ErrorInfo };

/** Która czynność wniesienia trwa; pusto, gdy żadna. */
type Wniesienie = 'plik' | 'pdf' | null;

/** Kontekst ostatniej czynności nieudanej — klucz sięga po tytuł w `tresciPliki.odmowa`. */
type KontekstOdmowy = keyof typeof tresciPliki.odmowa;

const ETYKIETA_STANU: Record<StudioIngestState, string> = {
  [StudioIngestState.Oczekuje]: tresciPliki.stan.oczekuje,
  [StudioIngestState.Przetwarzanie]: tresciPliki.stan.przetwarzanie,
  [StudioIngestState.Gotowa]: tresciPliki.stan.gotowa,
  [StudioIngestState.Ponowienie]: tresciPliki.stan.ponowienie,
  [StudioIngestState.Odmowa]: tresciPliki.stan.odmowa,
};

const KLASA_STANU: Record<StudioIngestState, string> = {
  [StudioIngestState.Oczekuje]: 'dn-plakietka',
  [StudioIngestState.Przetwarzanie]: 'dn-plakietka dn-plakietka--informacja',
  [StudioIngestState.Gotowa]: 'dn-plakietka dn-plakietka--sukces',
  [StudioIngestState.Ponowienie]: 'dn-plakietka dn-plakietka--ostrzezenie',
  [StudioIngestState.Odmowa]: 'dn-plakietka dn-plakietka--blad',
};

const NAZWA_RODZAJU: Record<StudioInputDeviceKind, string> = {
  [StudioInputDeviceKind.Skaner]: tresciPliki.urzadzenia.rodzaj.skaner,
  [StudioInputDeviceKind.Kamera]: tresciPliki.urzadzenia.rodzaj.kamera,
};

function znak(rysunek: NazwaZnaku): SVGElement {
  const wezel = zeZnacznika(ikony[rysunek]);
  wezel.setAttribute('aria-hidden', 'true');
  return wezel;
}

function pole(atrybuty: Record<string, string | boolean>): HTMLInputElement {
  return el('input', atrybuty) as HTMLInputElement;
}

function przyciskAkcji(etykieta: string, typ: 'button' | 'submit'): HTMLButtonElement {
  return el('button', {
    klasa: 'dn-btn dn-btn--zarys dn-btn--sm dn-na-koniec',
    type: typ,
    tekst: etykieta,
  }) as HTMLButtonElement;
}

function wskaznikLadowania(etykieta: string): HTMLElement {
  return el('div', { klasa: 'st-panel-wiersz' }, [
    el('span', { klasa: 'dn-kropka dn-kropka--sygnal dn-kropka--tetno', 'aria-hidden': 'true' }),
    etykieta,
  ]);
}

function stanPusty(zdanie: string): HTMLElement {
  return el('p', { klasa: 'dn-pusty-stan dn-pusty-stan--zwarty', tekst: zdanie });
}

function alertOdmowy(tytul: string, blad?: ErrorInfo): HTMLElement {
  return el('div', { klasa: 'dn-alert dn-alert--wstega dn-alert--blad', role: 'alert' }, [
    el('span', { klasa: 'dn-alert-tresc' }, [
      el('b', { tekst: tytul }),
      el('span', { tekst: opisOdmowy(blad, 'pliki') }),
    ]),
  ]);
}

export const montujPliki: MontazPanelu = (wezel, { kanal, idOkna }): ZamontowanyPanel => {
  let zdjete = false;

  let stanUrzadzen: StanUrzadzen = { rodzaj: 'ladowanie' };
  let stanKolejki: StanKolejki = idOkna === null ? { rodzaj: 'brakOkna' } : { rodzaj: 'ladowanie' };
  let bladDzialania: { kontekst: KontekstOdmowy; blad?: ErrorInfo } | null = null;
  let bilansWniesienia: StudioImportBalance | null = null;
  let dodawanieWToku = false;
  let adresWToku = false;
  let wniesienieWToku: Wniesienie = null;

  const urzadzeniaWToku = new Set<string>();
  const pozycjeWToku = new Set<string>();
  const slowaPozycji = new Map<string, StudioRecognizedWord[]>();
  const slowaWToku = new Set<string>();
  const slowaZapisane = new Set<string>();
  /* Tekst wpisany w pole słowa przeżywa przebudowę wykazu: zdarzenie rdzenia
     może przyjść w trakcie poprawiania, a wpis Operatora nie ma prawa zginąć. */
  const slowaWpisane = new Map<string, string>();

  const bezOkna = idOkna === null;

  /* ── Belka i treść panelu ─────────────────────────────────────────────── */

  const poleSciezki = pole({
    klasa: 'dn-pole-kontrolka st-pole-rozciagniete',
    type: 'text',
    placeholder: tresciPliki.dodaj.etykietaSciezki,
    'aria-label': tresciPliki.dodaj.etykietaSciezki,
  });
  const przyciskSciezki = przyciskAkcji(tresciPliki.dodaj.doloz, 'submit');
  const formularzSciezki = el('form', { klasa: 'st-panel-wiersz' }, [poleSciezki, przyciskSciezki]);

  const poleAdresu = pole({
    klasa: 'dn-pole-kontrolka st-pole-rozciagniete',
    type: 'url',
    placeholder: tresciPliki.adres.etykieta,
    'aria-label': tresciPliki.adres.etykieta,
  });
  const przyciskAdresu = przyciskAkcji(tresciPliki.adres.dolacz, 'submit');
  const formularzAdresu = el('form', { klasa: 'st-panel-wiersz' }, [poleAdresu, przyciskAdresu]);

  const poleWprost = pole({
    klasa: 'dn-pole-kontrolka st-pole-rozciagniete',
    type: 'text',
    placeholder: tresciPliki.wprost.etykieta,
    'aria-label': tresciPliki.wprost.etykieta,
  });
  const przyciskWprostPlik = przyciskAkcji(tresciPliki.wprost.plik, 'button');
  const przyciskWprostPdf = przyciskAkcji(tresciPliki.wprost.pdf, 'button');
  const wierszWprost = el('div', { klasa: 'st-panel-wiersz' }, [poleWprost, przyciskWprostPlik, przyciskWprostPdf]);

  const bilansWidok = el('div');
  const alertWnoszenia = el('div');

  const blokWnoszenia = el('div', { klasa: 'st-panel-lista st-panel-lista--bez-dolu' }, [
    el('div', { klasa: 'st-szyna-grupa', tekst: tresciPliki.wnoszenie.tytul }),
    formularzSciezki,
    formularzAdresu,
    wierszWprost,
    el('p', { klasa: 'dn-nota', tekst: tresciPliki.wnoszenie.bezOkienka }),
    alertWnoszenia,
    bilansWidok,
  ]);

  const sekcjaUrzadzen = el('div', { klasa: 'st-panel-lista st-panel-lista--bez-dolu' });

  const poleFiltru = pole({
    type: 'search',
    placeholder: tresciPliki.filtr.zacheta,
    'aria-label': tresciPliki.filtr.etykieta,
  });
  const filtr = el('label', { klasa: 'dn-szukaj st-szukaj-odsun' }, [znak('szukaj'), poleFiltru]);

  const poleTytulu = pole({
    klasa: 'dn-pole-kontrolka st-pole-rozciagniete',
    type: 'text',
    placeholder: tresciPliki.tytulDokumentu.etykieta,
    'aria-label': tresciPliki.tytulDokumentu.etykieta,
  });

  const naglowekKolejki = el('div', { klasa: 'st-panel-lista st-panel-lista--bez-dolu' }, [
    el('div', { klasa: 'st-szyna-grupa', tekst: tresciPliki.kolejka.tytul }),
    filtr,
  ]);

  const sekcjaKolejki = el('div', { klasa: 'st-panel-lista' });

  wezel.classList.add('sta-okno');
  wezel.replaceChildren(
    el('header', { klasa: 'sta-okno-belka' }, [
      el('span', { klasa: 'sta-okno-tytul' }, [znak('plik'), el('b', { tekst: tresciPliki.panel.tytul })]),
    ]),
    el('div', { klasa: 'sta-okno-tresc' }, [blokWnoszenia, sekcjaUrzadzen, naglowekKolejki, sekcjaKolejki]),
  );

  formularzSciezki.addEventListener('submit', (zdarzenie) => {
    zdarzenie.preventDefault();
    void dolozPlik(poleSciezki.value);
  });
  formularzAdresu.addEventListener('submit', (zdarzenie) => {
    zdarzenie.preventDefault();
    void dolaczAdres(poleAdresu.value);
  });
  przyciskWprostPlik.addEventListener('click', () => void wniesWprost('plik', poleWprost.value));
  przyciskWprostPdf.addEventListener('click', () => void wniesWprost('pdf', poleWprost.value));
  poleFiltru.addEventListener('input', () => odswiezKolejke());

  const odsubskrybuj = kanal.naZdarzenie(EventType.StudioIngestChanged, (zdarzenie) => {
    if (idOkna === null || zdarzenie.windowId !== idOkna) return;
    zastapPozycje(zdarzenie.item);
    odswiezKolejke();
  });

  odswiez();
  void zaladujUrzadzenia();
  void zaladujKolejke();

  /* ── Odświeżanie widoku ───────────────────────────────────────────────── */

  function odswiez(): void {
    odswiezWnoszenie();
    sekcjaUrzadzen.replaceChildren(...zawartoscUrzadzen());
    odswiezKolejke();
  }

  function odswiezWnoszenie(): void {
    poleSciezki.disabled = bezOkna || dodawanieWToku;
    przyciskSciezki.disabled = bezOkna || dodawanieWToku;
    przyciskSciezki.textContent = dodawanieWToku ? tresciPliki.dodaj.dokladanie : tresciPliki.dodaj.doloz;

    poleAdresu.disabled = bezOkna || adresWToku;
    przyciskAdresu.disabled = bezOkna || adresWToku;
    przyciskAdresu.textContent = adresWToku ? tresciPliki.adres.dolaczanie : tresciPliki.adres.dolacz;

    const wnosi = wniesienieWToku !== null;
    poleWprost.disabled = bezOkna || wnosi;
    przyciskWprostPlik.disabled = bezOkna || wnosi;
    przyciskWprostPdf.disabled = bezOkna || wnosi;
    przyciskWprostPlik.textContent = wniesienieWToku === 'plik' ? tresciPliki.wprost.wnoszenie : tresciPliki.wprost.plik;
    przyciskWprostPdf.textContent = wniesienieWToku === 'pdf' ? tresciPliki.wprost.wnoszenie : tresciPliki.wprost.pdf;

    const odmowaWnoszenia = bladDzialania !== null && bladDzialania.kontekst === 'wniesienie';
    alertWnoszenia.replaceChildren(
      ...(odmowaWnoszenia && bladDzialania !== null
        ? [alertOdmowy(tresciPliki.odmowa.wniesienie, bladDzialania.blad)]
        : []),
    );
    bilansWidok.replaceChildren(...(bilansWniesienia === null ? [] : wierszeBilansu(bilansWniesienia)));
  }

  function odswiezKolejke(): void {
    sekcjaKolejki.replaceChildren(...zawartoscKolejki());
  }

  function zastapPozycje(nowa: StudioIngestItem): void {
    if (stanKolejki.rodzaj !== 'wykaz') return;
    const istnieje = stanKolejki.pozycje.some((p) => p.id === nowa.id);
    stanKolejki = {
      rodzaj: 'wykaz',
      pozycje: istnieje
        ? stanKolejki.pozycje.map((p) => (p.id === nowa.id ? nowa : p))
        : [...stanKolejki.pozycje, nowa],
    };
  }

  /* ── Wywołania rdzenia ────────────────────────────────────────────────── */

  async function zaladujUrzadzenia(): Promise<void> {
    stanUrzadzen = { rodzaj: 'ladowanie' };
    odswiez();
    const wynik = await wywolaj(kanal, Command.StudioIngestDeviceList, {});
    if (zdjete) return;
    if (!wynik.udany || wynik.wynik === undefined) {
      stanUrzadzen = { rodzaj: 'odmowa', blad: wynik.blad };
    } else {
      stanUrzadzen = { rodzaj: 'wykaz', urzadzenia: wynik.wynik.devices };
    }
    odswiez();
  }

  async function zaladujKolejke(): Promise<void> {
    if (idOkna === null) {
      stanKolejki = { rodzaj: 'brakOkna' };
      odswiez();
      return;
    }
    stanKolejki = { rodzaj: 'ladowanie' };
    odswiez();
    const wynik = await wywolaj(kanal, Command.StudioIngestQueueList, { windowId: idOkna });
    if (zdjete) return;
    if (!wynik.udany || wynik.wynik === undefined) {
      stanKolejki = { rodzaj: 'odmowa', blad: wynik.blad };
    } else {
      stanKolejki = { rodzaj: 'wykaz', pozycje: wynik.wynik.items };
    }
    odswiez();
  }

  async function skanuj(urzadzenie: StudioInputDevice): Promise<void> {
    if (idOkna === null) return;
    urzadzeniaWToku.add(urzadzenie.id);
    odswiez();
    const wynik = await wywolaj(kanal, Command.StudioIngestDeviceScan, {
      windowId: idOkna,
      deviceId: urzadzenie.id,
    });
    if (zdjete) return;
    urzadzeniaWToku.delete(urzadzenie.id);
    if (!wynik.udany) {
      bladDzialania = { kontekst: 'skan', blad: wynik.blad };
      odswiez();
      return;
    }
    bladDzialania = null;
    await zaladujKolejke();
  }

  async function dolozPlik(sciezka: string): Promise<void> {
    if (idOkna === null || sciezka.trim() === '') return;
    dodawanieWToku = true;
    odswiez();
    const wynik = await wywolaj(kanal, Command.StudioIngestQueueAdd, {
      windowId: idOkna,
      sourcePaths: [sciezka.trim()],
    });
    if (zdjete) return;
    dodawanieWToku = false;
    if (!wynik.udany) {
      bladDzialania = { kontekst: 'dodanie', blad: wynik.blad };
      odswiez();
      return;
    }
    bladDzialania = null;
    poleSciezki.value = '';
    await zaladujKolejke();
  }

  async function dolaczAdres(adres: string): Promise<void> {
    if (idOkna === null || adres.trim() === '') return;
    adresWToku = true;
    odswiez();
    const wynik = await wywolaj(kanal, Command.StudioIngestUrl, { windowId: idOkna, url: adres.trim() });
    if (zdjete) return;
    adresWToku = false;
    if (!wynik.udany) {
      bladDzialania = { kontekst: 'adres', blad: wynik.blad };
      odswiez();
      return;
    }
    bladDzialania = null;
    poleAdresu.value = '';
    await zaladujKolejke();
  }

  /**
   * Wnosi plik wprost do edytora, z pominięciem kolejki. PDF idzie osobną
   * komendą, bo odzyskanie tekstu z PDF jest odtworzeniem i zwraca bilans —
   * a PDF bez warstwy tekstowej wraca pozycją kolejki do rozpoznania.
   */
  async function wniesWprost(rodzaj: Exclude<Wniesienie, null>, sciezka: string): Promise<void> {
    if (idOkna === null || sciezka.trim() === '') return;
    wniesienieWToku = rodzaj;
    bilansWniesienia = null;
    odswiez();
    const zadanie = { windowId: idOkna, path: sciezka.trim() };
    const wynik =
      rodzaj === 'pdf'
        ? await wywolaj(kanal, Command.StudioDocumentImportPdf, zadanie)
        : await wywolaj(kanal, Command.StudioDocumentImportFile, zadanie);
    if (zdjete) return;
    wniesienieWToku = null;
    if (!wynik.udany || wynik.wynik === undefined) {
      bladDzialania = { kontekst: 'wniesienie', blad: wynik.blad };
      odswiez();
      return;
    }
    bladDzialania = null;
    bilansWniesienia = wynik.wynik.balance;
    poleWprost.value = '';
    /* PDF bez warstwy tekstowej odkłada się pozycją kolejki — wykaz musi to
       pokazać od razu, inaczej Operator nie wie, gdzie szukać materiału. */
    if (bilansWniesienia.needsTextRecognition === true) {
      await zaladujKolejke();
      return;
    }
    odswiez();
  }

  async function rozpoznaj(pozycja: StudioIngestItem): Promise<void> {
    pozycjeWToku.add(pozycja.id);
    odswiez();
    const wynik = await wywolaj(kanal, Command.StudioIngestRecognize, { itemId: pozycja.id });
    if (zdjete) return;
    pozycjeWToku.delete(pozycja.id);
    if (!wynik.udany || wynik.wynik === undefined) {
      bladDzialania = { kontekst: 'rozpoznanie', blad: wynik.blad };
      odswiez();
      return;
    }
    bladDzialania = null;
    zastapPozycje(wynik.wynik.item);
    if (wynik.wynik.words !== undefined) slowaPozycji.set(pozycja.id, wynik.wynik.words);
    odswiez();
  }

  async function zapiszPoprawke(
    pozycja: StudioIngestItem,
    slowo: StudioRecognizedWord,
    tekstPoprawiony: string,
  ): Promise<void> {
    const klucz = `${pozycja.id}:${slowo.index}`;
    slowaWToku.add(klucz);
    odswiezKolejke();
    const wynik = await wywolaj(kanal, Command.StudioIngestCorrectionSet, {
      itemId: pozycja.id,
      wordIndex: slowo.index,
      text: tekstPoprawiony,
    });
    if (zdjete) return;
    slowaWToku.delete(klucz);
    if (!wynik.udany || wynik.wynik === undefined) {
      bladDzialania = { kontekst: 'poprawka', blad: wynik.blad };
      odswiez();
      return;
    }
    bladDzialania = null;
    zastapPozycje(wynik.wynik.item);
    const listaSlow = slowaPozycji.get(pozycja.id);
    if (listaSlow !== undefined) {
      slowaPozycji.set(
        pozycja.id,
        listaSlow.map((w) => (w.index === slowo.index ? { ...w, text: tekstPoprawiony, corrected: true } : w)),
      );
    }
    slowaZapisane.add(klucz);
    slowaWpisane.delete(klucz);
    odswiezKolejke();
  }

  async function przyjmij(pozycja: StudioIngestItem): Promise<void> {
    if (idOkna === null) return;
    pozycjeWToku.add(pozycja.id);
    odswiez();
    const tytul = poleTytulu.value.trim();
    const wynik = await wywolaj(kanal, Command.StudioIngestItemAccept, {
      windowId: idOkna,
      itemIds: [pozycja.id],
      ...(tytul === '' ? {} : { title: tytul }),
    });
    if (zdjete) return;
    pozycjeWToku.delete(pozycja.id);
    if (!wynik.udany) {
      bladDzialania = { kontekst: 'przyjecie', blad: wynik.blad };
      odswiez();
      return;
    }
    bladDzialania = null;
    slowaPozycji.delete(pozycja.id);
    poleTytulu.value = '';
    await zaladujKolejke();
  }

  /* ── Wykaz urządzeń ───────────────────────────────────────────────────── */

  function wierszUrzadzenia(urzadzenie: StudioInputDevice): HTMLElement {
    const wToku = urzadzeniaWToku.has(urzadzenie.id);
    const przycisk = przyciskAkcji(
      wToku ? tresciPliki.urzadzenia.skanowanie : tresciPliki.urzadzenia.skanuj,
      'button',
    );
    przycisk.disabled = wToku || bezOkna;
    przycisk.addEventListener('click', () => void skanuj(urzadzenie));
    return el('div', { klasa: 'st-panel-wiersz' }, [
      el('b', { tekst: urzadzenie.name }),
      el('span', { klasa: 'dn-plakietka', tekst: NAZWA_RODZAJU[urzadzenie.kind] }),
      urzadzenie.hasFeeder === true ? el('span', { klasa: 'dn-meta', tekst: tresciPliki.urzadzenia.podajnik }) : null,
      przycisk,
    ]);
  }

  function zawartoscUrzadzen(): Node[] {
    const naglowek = el('div', { klasa: 'st-szyna-grupa', tekst: tresciPliki.urzadzenia.tytul });
    const odmowaSkanu = bladDzialania !== null && bladDzialania.kontekst === 'skan';
    const alert = odmowaSkanu && bladDzialania !== null ? [alertOdmowy(tresciPliki.odmowa.skan, bladDzialania.blad)] : [];
    switch (stanUrzadzen.rodzaj) {
      case 'ladowanie':
        return [naglowek, ...alert, wskaznikLadowania(tresciPliki.urzadzenia.ladowanie)];
      case 'odmowa':
        return [naglowek, alertOdmowy(tresciPliki.odmowa.urzadzenia, stanUrzadzen.blad)];
      case 'wykaz':
        if (stanUrzadzen.urzadzenia.length === 0) return [naglowek, ...alert, stanPusty(tresciPliki.urzadzenia.brak)];
        return [naglowek, ...alert, ...stanUrzadzen.urzadzenia.map(wierszUrzadzenia)];
    }
  }

  /* ── Bilans wniesienia ────────────────────────────────────────────────── */

  function miara(liczba: number | undefined, nazwa: string): HTMLElement | null {
    if (liczba === undefined) return null;
    return el('span', { klasa: 'dn-meta', tekst: `${liczba} ${nazwa}` });
  }

  function wierszeBilansu(bilans: StudioImportBalance): Node[] {
    const wezly: Node[] = [
      el('div', { klasa: 'st-panel-wiersz' }, [
        el('b', { tekst: tresciPliki.wprost.bilans }),
        miara(bilans.pagesWithText, tresciPliki.wprost.stronyZTekstem),
        miara(bilans.paragraphsRecovered, tresciPliki.wprost.akapity),
        miara(bilans.tablesRecognized, tresciPliki.wprost.tabele),
        miara(bilans.tablesMissed, tresciPliki.wprost.tabeleNierozpoznane),
        miara(bilans.imagesEmbedded, tresciPliki.wprost.obrazy),
        miara(bilans.imagesSkipped, tresciPliki.wprost.obrazyPominiete),
      ]),
    ];
    /* Zdanie o stanie wyniku pisze rdzeń treścią dla Operatora, nie kodem
       odmowy — idzie do widoku wprost, bo tylko ono mówi, czego nie odzyskano. */
    if (bilans.note !== undefined && bilans.note !== '') {
      wezly.push(el('p', { klasa: 'dn-nota', tekst: bilans.note }));
    }
    if (bilans.needsTextRecognition === true) {
      wezly.push(el('p', { klasa: 'dn-nota', tekst: tresciPliki.wprost.bezWarstwyTekstu }));
    }
    return wezly;
  }

  /* ── Kolejka wczytywania ──────────────────────────────────────────────── */

  function akcjaDlaPozycji(pozycja: StudioIngestItem, wToku: boolean): HTMLElement | null {
    if (pozycja.state === StudioIngestState.Oczekuje || pozycja.state === StudioIngestState.Ponowienie) {
      const przycisk = przyciskAkcji(
        wToku ? tresciPliki.dzialania.rozpoznawanie : tresciPliki.dzialania.rozpoznaj,
        'button',
      );
      przycisk.disabled = wToku;
      przycisk.addEventListener('click', () => void rozpoznaj(pozycja));
      return przycisk;
    }
    if (pozycja.state === StudioIngestState.Gotowa) {
      const przycisk = przyciskAkcji(
        wToku ? tresciPliki.dzialania.przyjmowanie : tresciPliki.dzialania.przyjmij,
        'button',
      );
      przycisk.disabled = wToku || bezOkna;
      przycisk.addEventListener('click', () => void przyjmij(pozycja));
      return przycisk;
    }
    return null;
  }

  function wierszSlowa(pozycja: StudioIngestItem, slowo: StudioRecognizedWord): HTMLElement {
    const klucz = `${pozycja.id}:${slowo.index}`;
    const wToku = slowaWToku.has(klucz);
    const poleSlowa = pole({
      klasa: 'dn-pole-kontrolka st-pole-rozciagniete',
      type: 'text',
      'aria-label': tresciPliki.poprawa.etykietaSlowa,
    });
    poleSlowa.value = slowaWpisane.get(klucz) ?? slowo.text;
    poleSlowa.disabled = wToku;
    poleSlowa.addEventListener('input', () => slowaWpisane.set(klucz, poleSlowa.value));
    const przycisk = przyciskAkcji(wToku ? tresciPliki.poprawa.zapisywanie : tresciPliki.poprawa.zapisz, 'button');
    przycisk.disabled = wToku;
    przycisk.addEventListener('click', () => void zapiszPoprawke(pozycja, slowo, poleSlowa.value));
    return el('div', { klasa: 'st-panel-wiersz' }, [
      poleSlowa,
      el('span', { klasa: 'dn-meta', tekst: `${Math.round(slowo.confidence * 100)}%` }),
      slowaZapisane.has(klucz) ? el('span', { klasa: 'dn-meta', tekst: tresciPliki.poprawa.zapisano }) : null,
      przycisk,
    ]);
  }

  function blokPoprawy(pozycja: StudioIngestItem, slowa: StudioRecognizedWord[]): HTMLElement {
    return el('div', { klasa: 'st-panel-lista st-panel-lista--bez-gory st-panel-lista--bez-dolu' }, [
      el('p', { klasa: 'dn-nota', tekst: tresciPliki.poprawa.tytul }),
      ...slowa.map((s) => wierszSlowa(pozycja, s)),
    ]);
  }

  function wierszPozycji(pozycja: StudioIngestItem): HTMLElement[] {
    const wToku = pozycjeWToku.has(pozycja.id);
    const wiersze: HTMLElement[] = [
      el('div', { klasa: 'st-panel-wiersz' }, [
        el('b', { tekst: nazwaZrodla(pozycja) }),
        el('span', { klasa: KLASA_STANU[pozycja.state], tekst: ETYKIETA_STANU[pozycja.state] }),
        pozycja.pages !== undefined
          ? el('span', { klasa: 'dn-meta', tekst: zLiczba(pozycja.pages, tresciPliki.kolejka.stron) })
          : null,
        pozycja.confidence !== undefined
          ? el('span', {
              klasa: 'dn-meta',
              tekst: `${tresciPliki.kolejka.pewnosc} ${Math.round(pozycja.confidence * 100)}%`,
            })
          : null,
        pozycja.usedOcr === true ? el('span', { klasa: 'dn-meta', tekst: tresciPliki.kolejka.zPisma }) : null,
        akcjaDlaPozycji(pozycja, wToku),
      ]),
    ];
    if (pozycja.state === StudioIngestState.Odmowa) {
      /* Powód odmowy zapisany przez serwer idzie do dziennika, nie do widoku:
         opisuje wnętrze wczytywania, a Operatorowi potrzebne jest zdanie
         mówiące, co z tą pozycją zrobić. */
      if (pozycja.failureReason !== undefined) {
        console.warn('[studio] pozycja kolejki odrzucona', pozycja.id, pozycja.failureReason);
      }
      wiersze.push(el('p', { klasa: 'dn-nota', tekst: tresciPliki.kolejka.odrzucona }));
    }
    const slowa = slowaPozycji.get(pozycja.id);
    if (slowa !== undefined && slowa.length > 0) wiersze.push(blokPoprawy(pozycja, slowa));
    return wiersze;
  }

  /** Nazwa pozycji w wykazie: ścieżka, zasób albo — gdy rdzeń nie dał ani jednego — jej identyfikator. */
  function nazwaZrodla(pozycja: StudioIngestItem): string {
    return pozycja.sourcePath ?? pozycja.assetId ?? pozycja.id;
  }

  function zawartoscKolejki(): Node[] {
    const wezly: Node[] = [];
    if (bladDzialania !== null && bladDzialania.kontekst !== 'skan' && bladDzialania.kontekst !== 'wniesienie') {
      wezly.push(alertOdmowy(tresciPliki.odmowa[bladDzialania.kontekst], bladDzialania.blad));
    }
    switch (stanKolejki.rodzaj) {
      case 'brakOkna':
        wezly.push(stanPusty(tresciPliki.panel.brakOkna));
        return wezly;
      case 'ladowanie':
        wezly.push(wskaznikLadowania(tresciPliki.kolejka.ladowanie));
        return wezly;
      case 'odmowa':
        wezly.push(alertOdmowy(tresciPliki.odmowa.kolejka, stanKolejki.blad));
        return wezly;
      case 'wykaz': {
        if (stanKolejki.pozycje.length === 0) {
          wezly.push(stanPusty(tresciPliki.kolejka.brak));
          return wezly;
        }
        const szukane = poleFiltru.value.trim().toLocaleLowerCase('pl');
        const widoczne =
          szukane === ''
            ? stanKolejki.pozycje
            : stanKolejki.pozycje.filter((p) => nazwaZrodla(p).toLocaleLowerCase('pl').includes(szukane));
        if (widoczne.length === 0) {
          wezly.push(stanPusty(tresciPliki.filtr.brakTrafien));
          return wezly;
        }
        wezly.push(el('div', { klasa: 'st-panel-wiersz' }, [poleTytulu]));
        for (const pozycja of widoczne) wezly.push(...wierszPozycji(pozycja));
        return wezly;
      }
    }
  }

  return {
    zdejmij() {
      zdjete = true;
      odsubskrybuj();
      wezel.replaceChildren();
    },
  };
};
