/**
 * Panel Pliki okna roboczego Studia. Wykaz urządzeń wejściowych rdzenia
 * (skanery, kamery) i kolejka wczytywania okna: dołożenie materiału ze
 * ścieżki widzianej przez rdzeń albo ze strony sieciowej, skanowanie
 * urządzeniem, rozpoznanie pisma, poprawa rozpoznanych słów i przyjęcie
 * wyniku jako dokument Studio Editora.
 *
 * Kolejka jest własnością okna modułu (`windowId`); bez niego rdzeń nie ma,
 * gdzie osadzić pozycję, więc sekcja kolejki niesie nazwany stan pusty
 * zamiast formularzy. Wykaz urządzeń nie zależy od okna i wczytuje się zawsze.
 *
 * `studio.ingest.changed` odświeża pozycję kolejki bez odpytywania w pętli —
 * jedyne odpytanie ponowne idzie po własnej czynności Operatora (dołożenie,
 * skan, rozpoznanie, przyjęcie), tym samym torem co zapis i odczyt w
 * `dokument.ts`.
 */

import {
  Command,
  EventType,
  StudioIngestState,
  type ErrorInfo,
  type StudioIngestItem,
  type StudioInputDevice,
  type StudioRecognizedWord,
} from '../../../../../shared/contract.ts';
import { wywolaj } from '../../../protokol/wywolanie.ts';
import { el, zeZnacznika } from '../narzedzia.ts';
import { ikony, type NazwaZnaku } from '../ikony.ts';
import { tresci } from '../tresci.ts';
import { tresciPliki } from './pliki-tresci.ts';
import type { MontazPanelu, ZamontowanyPanel } from './umowa.ts';
import { opisOdmowy } from '../odmowa.ts';

type StanUrzadzen =
  | { rodzaj: 'ladowanie' }
  | { rodzaj: 'wykaz'; urzadzenia: StudioInputDevice[] }
  | { rodzaj: 'odmowa'; blad?: ErrorInfo };

type StanKolejki =
  | { rodzaj: 'brakOkna' }
  | { rodzaj: 'ladowanie' }
  | { rodzaj: 'wykaz'; pozycje: StudioIngestItem[] }
  | { rodzaj: 'odmowa'; blad?: ErrorInfo };

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

function znak(rysunek: NazwaZnaku): SVGElement {
  const wezel = zeZnacznika(ikony[rysunek]);
  wezel.setAttribute('aria-hidden', 'true');
  return wezel;
}

function wskaznikLadowania(etykieta: string): HTMLElement {
  return el('div', { klasa: 'dn-wykaz-modulu-poz' }, [
    el('span', { klasa: 'dn-kropka dn-kropka--sygnal dn-kropka--tetno', 'aria-hidden': 'true' }),
    etykieta,
  ]);
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
  let dodawanieWToku = false;
  let adresWToku = false;

  const urzadzeniaWToku = new Set<string>();
  const pozycjeWToku = new Set<string>();
  const slowaPozycji = new Map<string, StudioRecognizedWord[]>();
  const slowaWToku = new Set<string>();
  const slowaZapisane = new Set<string>();

  const sekcjaUrzadzen = el('div', { klasa: 'st-panel-lista st-panel-lista--bez-dolu' });
  const sekcjaKolejki = el('div', { klasa: 'st-panel-lista st-panel-lista--bez-gory' });

  const sekcja = el(
    'section',
    {
      klasa: 'sta-okno',
      'data-nazwa': tresci.karty.pliki,
      hidden: true,
    },
    [
      el('header', { klasa: 'sta-okno-belka' }, [
        el('span', { klasa: 'sta-okno-tytul' }, [znak('plik'), el('b', { tekst: tresci.karty.pliki })]),
      ]),
      el('div', { klasa: 'sta-okno-tresc' }, [sekcjaUrzadzen, sekcjaKolejki]),
    ],
  );
  wezel.appendChild(sekcja);

  const odsubskrybuj = kanal.naZdarzenie(EventType.StudioIngestChanged, (zdarzenie) => {
    if (idOkna === null || zdarzenie.windowId !== idOkna) return;
    zastapPozycje(zdarzenie.item);
    odswiez();
  });

  odswiez();
  void zaladujUrzadzenia();
  void zaladujKolejke();

  function odswiez(): void {
    sekcjaUrzadzen.replaceChildren(...zawartoscUrzadzen());
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
    await zaladujKolejke();
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
    zastapPozycje(wynik.wynik.item);
    if (wynik.wynik.words !== undefined) slowaPozycji.set(pozycja.id, wynik.wynik.words);
    odswiez();
  }

  async function zapiszPoprawke(pozycja: StudioIngestItem, slowo: StudioRecognizedWord, tekstPoprawiony: string): Promise<void> {
    const klucz = `${pozycja.id}:${slowo.index}`;
    slowaWToku.add(klucz);
    odswiez();
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
    zastapPozycje(wynik.wynik.item);
    const listaSlow = slowaPozycji.get(pozycja.id);
    if (listaSlow !== undefined) {
      slowaPozycji.set(
        pozycja.id,
        listaSlow.map((w) => (w.index === slowo.index ? { ...w, text: tekstPoprawiony } : w)),
      );
    }
    slowaZapisane.add(klucz);
    odswiez();
  }

  async function przyjmij(pozycja: StudioIngestItem): Promise<void> {
    if (idOkna === null) return;
    pozycjeWToku.add(pozycja.id);
    odswiez();
    const wynik = await wywolaj(kanal, Command.StudioIngestItemAccept, {
      windowId: idOkna,
      itemIds: [pozycja.id],
    });
    if (zdjete) return;
    pozycjeWToku.delete(pozycja.id);
    if (!wynik.udany) {
      bladDzialania = { kontekst: 'przyjecie', blad: wynik.blad };
      odswiez();
      return;
    }
    slowaPozycji.delete(pozycja.id);
    await zaladujKolejke();
  }

  function wierszUrzadzenia(urzadzenie: StudioInputDevice): HTMLElement {
    const wToku = urzadzeniaWToku.has(urzadzenie.id);
    const przycisk = el('button', {
      klasa: 'dn-btn dn-btn--zarys dn-btn--sm dn-na-koniec',
      type: 'button',
      tekst: wToku ? tresciPliki.urzadzenia.skanowanie : tresciPliki.urzadzenia.skanuj,
      disabled: wToku || idOkna === null,
    });
    przycisk.addEventListener('click', () => void skanuj(urzadzenie));
    return el('div', { klasa: 'dn-wykaz-modulu-poz' }, [
      el('b', { tekst: urzadzenie.name }),
      el('span', { klasa: 'dn-plakietka', tekst: urzadzenie.kind }),
      urzadzenie.hasFeeder === true ? el('span', { klasa: 'dn-meta', tekst: tresciPliki.urzadzenia.podajnik }) : null,
      przycisk,
    ]);
  }

  function zawartoscUrzadzen(): Node[] {
    const naglowek = el('div', { klasa: 'st-szyna-grupa', tekst: tresciPliki.urzadzenia.tytul });
    switch (stanUrzadzen.rodzaj) {
      case 'ladowanie':
        return [naglowek, wskaznikLadowania(tresciPliki.urzadzenia.ladowanie)];
      case 'odmowa':
        return [naglowek, alertOdmowy(tresciPliki.odmowa.urzadzenia, stanUrzadzen.blad)];
      case 'wykaz':
        if (stanUrzadzen.urzadzenia.length === 0) {
          return [naglowek, el('p', { klasa: 'dn-pusty-stan dn-pusty-stan--zwarty', tekst: tresciPliki.urzadzenia.brak })];
        }
        return [naglowek, ...stanUrzadzen.urzadzenia.map(wierszUrzadzenia)];
    }
  }

  function formularzDolozenia(): HTMLElement {
    const pole = el('input', {
      klasa: 'dn-pole-kontrolka st-pole-rozciagniete',
      type: 'text',
      placeholder: tresciPliki.dodaj.etykietaSciezki,
      'aria-label': tresciPliki.dodaj.etykietaSciezki,
      disabled: dodawanieWToku,
    }) as HTMLInputElement;
    const przycisk = el('button', {
      klasa: 'dn-btn dn-btn--zarys dn-btn--sm dn-na-koniec',
      type: 'submit',
      tekst: dodawanieWToku ? tresciPliki.dodaj.dokladanie : tresciPliki.dodaj.doloz,
      disabled: dodawanieWToku,
    });
    przycisk.style.flexShrink = '0';
    const formularz = el('form', { klasa: 'dn-wykaz-modulu-poz' }, [pole, przycisk]);
    formularz.addEventListener('submit', (zdarzenie) => {
      zdarzenie.preventDefault();
      void dolozPlik(pole.value);
    });
    return formularz;
  }

  function formularzAdresu(): HTMLElement {
    const pole = el('input', {
      klasa: 'dn-pole-kontrolka st-pole-rozciagniete',
      type: 'url',
      placeholder: tresciPliki.adres.etykieta,
      'aria-label': tresciPliki.adres.etykieta,
      disabled: adresWToku,
    }) as HTMLInputElement;
    const przycisk = el('button', {
      klasa: 'dn-btn dn-btn--zarys dn-btn--sm dn-na-koniec',
      type: 'submit',
      tekst: adresWToku ? tresciPliki.adres.dolaczanie : tresciPliki.adres.dolacz,
      disabled: adresWToku,
    });
    przycisk.style.flexShrink = '0';
    const formularz = el('form', { klasa: 'dn-wykaz-modulu-poz' }, [pole, przycisk]);
    formularz.addEventListener('submit', (zdarzenie) => {
      zdarzenie.preventDefault();
      void dolaczAdres(pole.value);
    });
    return formularz;
  }

  function akcjaDlaPozycji(pozycja: StudioIngestItem, wToku: boolean): HTMLElement | null {
    if (pozycja.state === StudioIngestState.Oczekuje || pozycja.state === StudioIngestState.Ponowienie) {
      const przycisk = el('button', {
        klasa: 'dn-btn dn-btn--zarys dn-btn--sm dn-na-koniec',
        type: 'button',
        tekst: wToku ? tresciPliki.dzialania.rozpoznawanie : tresciPliki.dzialania.rozpoznaj,
        disabled: wToku,
      });
      przycisk.addEventListener('click', () => void rozpoznaj(pozycja));
      return przycisk;
    }
    if (pozycja.state === StudioIngestState.Gotowa) {
      const przycisk = el('button', {
        klasa: 'dn-btn dn-btn--zarys dn-btn--sm dn-na-koniec',
        type: 'button',
        tekst: wToku ? tresciPliki.dzialania.przyjmowanie : tresciPliki.dzialania.przyjmij,
        disabled: wToku,
      });
      przycisk.addEventListener('click', () => void przyjmij(pozycja));
      return przycisk;
    }
    return null;
  }

  function wierszSlowa(pozycja: StudioIngestItem, slowo: StudioRecognizedWord): HTMLElement {
    const klucz = `${pozycja.id}:${slowo.index}`;
    const wToku = slowaWToku.has(klucz);
    const zapisane = slowaZapisane.has(klucz);
    const pole = el('input', {
      klasa: 'dn-pole-kontrolka st-pole-rozciagniete',
      type: 'text',
      'aria-label': tresciPliki.poprawa.etykietaSlowa,
      disabled: wToku,
    }) as HTMLInputElement;
    pole.value = slowo.text;
    const przycisk = el('button', {
      klasa: 'dn-btn dn-btn--zarys dn-btn--sm dn-na-koniec',
      type: 'button',
      tekst: wToku ? tresciPliki.poprawa.zapisywanie : tresciPliki.poprawa.zapisz,
      disabled: wToku,
    });
    przycisk.style.flexShrink = '0';
    przycisk.addEventListener('click', () => void zapiszPoprawke(pozycja, slowo, pole.value));
    return el('div', { klasa: 'dn-wykaz-modulu-poz' }, [
      el('span', { klasa: 'dn-meta', tekst: `#${slowo.index}` }),
      pole,
      el('span', { klasa: 'dn-meta', tekst: `${Math.round(slowo.confidence * 100)}%` }),
      zapisane ? el('span', { klasa: 'dn-meta', tekst: tresciPliki.poprawa.zapisano }) : null,
      przycisk,
    ]);
  }

  function blokPoprawy(pozycja: StudioIngestItem, slowa: StudioRecognizedWord[]): HTMLElement {
    return el(
      'div',
      { klasa: 'st-panel-lista st-panel-lista--bez-gory st-panel-lista--bez-dolu' },
      [el('p', { klasa: 'dn-nota', tekst: tresciPliki.poprawa.tytul }), ...slowa.map((s) => wierszSlowa(pozycja, s))],
    );
  }

  function wierszPozycji(pozycja: StudioIngestItem): HTMLElement[] {
    const wToku = pozycjeWToku.has(pozycja.id);
    const zrodlo = pozycja.sourcePath ?? pozycja.assetId ?? pozycja.id;
    const wiersze: HTMLElement[] = [
      el('div', { klasa: 'dn-wykaz-modulu-poz' }, [
        el('b', { tekst: zrodlo }),
        el('span', { klasa: KLASA_STANU[pozycja.state], tekst: ETYKIETA_STANU[pozycja.state] }),
        pozycja.pages !== undefined ? el('span', { klasa: 'dn-meta', tekst: `${tresciPliki.kolejka.stron} ${pozycja.pages}` }) : null,
        pozycja.confidence !== undefined
          ? el('span', { klasa: 'dn-meta', tekst: `${tresciPliki.kolejka.pewnosc} ${Math.round(pozycja.confidence * 100)}%` })
          : null,
        pozycja.usedOcr === true ? el('span', { klasa: 'dn-meta', tekst: tresciPliki.kolejka.zPisma }) : null,
        akcjaDlaPozycji(pozycja, wToku),
      ]),
    ];
    if (pozycja.state === StudioIngestState.Odmowa) {
      wiersze.push(el('p', { klasa: 'dn-nota', tekst: pozycja.failureReason ?? tresci.odmowa.brakOpisu }));
    }
    const slowa = slowaPozycji.get(pozycja.id);
    if (slowa !== undefined && slowa.length > 0) wiersze.push(blokPoprawy(pozycja, slowa));
    return wiersze;
  }

  function zawartoscKolejki(): Node[] {
    const naglowek = el('div', { klasa: 'st-szyna-grupa', tekst: tresciPliki.kolejka.tytul });
    const wezly: Node[] = [naglowek];
    if (bladDzialania !== null) {
      wezly.push(alertOdmowy(tresciPliki.odmowa[bladDzialania.kontekst], bladDzialania.blad));
    }
    switch (stanKolejki.rodzaj) {
      case 'brakOkna':
        wezly.push(el('p', { klasa: 'dn-nota', tekst: tresci.dokumentOdmowa.okno }));
        return wezly;
      case 'ladowanie':
        wezly.push(wskaznikLadowania(tresciPliki.kolejka.ladowanie));
        return wezly;
      case 'odmowa':
        wezly.push(alertOdmowy(tresciPliki.odmowa.kolejka, stanKolejki.blad));
        return wezly;
      case 'wykaz':
        wezly.push(formularzDolozenia(), formularzAdresu());
        if (stanKolejki.pozycje.length === 0) {
          wezly.push(el('p', { klasa: 'dn-pusty-stan dn-pusty-stan--zwarty', tekst: tresciPliki.kolejka.brak }));
        } else {
          for (const pozycja of stanKolejki.pozycje) wezly.push(...wierszPozycji(pozycja));
        }
        return wezly;
    }
  }

  return {
    zdejmij() {
      zdjete = true;
      odsubskrybuj();
      if (sekcja.isConnected) sekcja.remove();
    },
  };
};
